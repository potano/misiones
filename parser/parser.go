// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package parser

import (
	"potano.misiones/sexp"
)



func (g PreparedGrammar) ParseList(parent ListItemType, lispList sexp.LispList,
		) (ListItemType, error) {
	var symbol = lispList.Head()
	guide, exists := g[symbol]
	if !exists {
		return nil, lispList.Error("unrecognized list type %s", symbol)
	}
	var listName string
	list := lispList.List()
	if guide.nameRequirement != UnnamedList {
		if len(list) > 0 && list[0].MayBeHead() {
			listName = list[0].(sexp.LispScalar).String()
			list = list[1:]
		}
	}
	if guide.nameRequirement == NameRequired && len(listName) == 0 {
		return nil, lispList.Error("list type '%s' requires a name", symbol)
	}
	targetMap := map[string][]sexp.LispValue{}
	for _, item := range list {
		var action symbolAction
		var exists bool
		if l, isList := item.(sexp.LispList); isList {
			listHead := l.Head()
			if action, exists = guide.nonterminalActions[listHead]; !exists {
				if _, exists = g[listHead]; !exists {
					return nil, l.Error("unrecognized list type %s", listHead)
				}
				return nil, l.Error("%s list is not allowed in a %s list",
					listHead, symbol)
			}
		} else {
			mask := item.TypeMask()
			for _, action = range guide.terminalActions {
				if (mask & action.mask) > 0 {
					exists = true
					break
				}
			}
			if !exists {
				return nil, item.Error("%s value is not allowed in list type %s",
					item.Desc(), symbol)
			}
		}
		consumed := false
		for _, targetName := range action.targets {
			var mapped []sexp.LispValue
			if mapped, exists = targetMap[targetName]; exists {
				maxCount := guide.targets[targetName].MaxCount
				if maxCount == 0 || byte(len(mapped)) < maxCount {
					mapped = append(mapped, item)
					targetMap[targetName] = mapped
					consumed = true
					break
				}
			} else {
				targetMap[targetName] = []sexp.LispValue{item}
				consumed = true
				break
			}
		}
		if !consumed {
			return nil, item.Error("%s is illegal in this context", item.Desc())
		}
	}
	for _, targetSpec := range guide.targets {
		targetEntries := targetMap[targetSpec.Name]
		numEntries := len(targetEntries)
		if numEntries < int(targetSpec.MinCount) {
			if targetSpec.MinCount == 1 {
				return nil, lispList.Error("%s list lacks %s entry", symbol,
					targetSpec.Name)
			}
			return nil, lispList.Error("list requires at least %d %s entries, got %d",
				targetSpec.MinCount, targetSpec.Name, numEntries)
		}
		if targetSpec.InMultiplesOf > 1 &&
				(numEntries % int(targetSpec.InMultiplesOf)) > 0 {
			errorSource := lispList.Source()
			if numEntries > 0 {
				errorSource = targetEntries[numEntries - 1].Source()
			}
			return nil, errorSource.Error(
				"number of %s entries must be a multiple of %d",
				targetSpec.Name, targetSpec.InMultiplesOf)
		}
	}
	listItem, err := parent.NewChild(symbol, listName, lispList.Source())
	if err != nil {
		return nil, err
	}
	for _, targName := range guide.targetOrder {
		target, exists := targetMap[targName]
		if !exists {
			continue
		}
		if target[0].IsList() {
			for _, item := range target {
				source := item.Source()
				targList := item.(sexp.LispList)
				childItem, err := g.ParseList(listItem, targList)
				if err != nil {
					return nil, err
				}
				err = listItem.SetList(targName, targList.Head(), source, childItem)
				if err != nil {
					return nil, err
				}
			}
		} else {
			lispScalars := make([]sexp.LispScalar, len(target))
			for i, item := range target {
				lispScalars[i] = item.(sexp.LispScalar)
			}
			err = listItem.SetScalars(targName, lispScalars)
			if err != nil {
				return nil, err
			}
		}
	}
	return listItem, nil
}

