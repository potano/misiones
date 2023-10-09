// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package parser

import (
	"fmt"
	"strings"

	"potano.misiones/sexp"
)


type Grammar []ListNonterminal

type ListNonterminal struct {
	TypeName string
	NameRequirement byte
	SymbolActions []SymbolAction
	Targets []TargetSpec
}

const (
	UnnamedList = iota
	NameOptional
	NameRequired
)

type ListItemType interface {
	SetScalars(targetName string, scalars []sexp.LispScalar) error
	SetList(targetName, symbol string, source sexp.ValueSource, value ListItemType) error
	NewChild(listType, listName string, source sexp.ValueSource) (ListItemType, error)
}

type SymbolAction struct {
	ListName string
	Mask uint32
	Target string
}

type TargetSpec struct {
	Name string
	MinCount, MaxCount, InMultiplesOf byte
}





type PreparedGrammar map[string]nonterminalPattern

type nonterminalPattern struct {
	nameRequirement byte
	terminalActions []symbolAction
	nonterminalActions map[string]symbolAction
	targets map[string]TargetSpec
	targetOrder []string
}

type symbolAction struct {
	mask uint32
	targets []string
}

const (
	targetAny = iota
	targetList
	targetScalar
)

func PrepareGrammar(grammarSpec Grammar) (PreparedGrammar, error) {
	prepared := map[string]nonterminalPattern{}
	for _, list := range grammarSpec {
		name := list.TypeName
		if _, exists := prepared[name]; exists {
			return prepared, fmt.Errorf("Setup error: redefinition of list type %s",
				name)
		}
		targets := make(map[string]TargetSpec, len(list.Targets))
		targetOrder := make([]string, len(list.Targets))
		targetHolds := make(map[string]int, len(list.Targets))
		for tX, targ := range list.Targets {
			tname := targ.Name
			if _, exists := targets[tname]; exists {
				return prepared, fmt.Errorf(
					"Setup error: multiple '%s' targets in '%s'",
					tname, name)
			}
			if targ.MaxCount > 0 && targ.MaxCount < targ.MinCount {
				return prepared, fmt.Errorf(
					"Setup error: bad range in target '%s' of '%s'",
					tname, name)
			}
			targets[tname] = targ
			targetOrder[tX] = tname
		}
		terminalActions := []symbolAction{}
		nonterminalActions := map[string]symbolAction{}
		var accumulatedMask uint32
		for _, act := range list.SymbolActions {
			var isList bool
			if len(act.ListName) > 0 || (act.Mask & sexp.TList) > 0 {
				if len(act.ListName) == 0 {
					return prepared, fmt.Errorf(
						"Setup error: '%s' refers to empty list name",
						name)
				}
				isList = true
			} else {
				if (act.Mask & accumulatedMask) > 0 {
					return prepared, fmt.Errorf(
						"Setup error: overlapping token masks in '%s'",
						name)
				}
				accumulatedMask |= act.Mask
			}
			targetNames := strings.Fields(act.Target)
			for _, targetName := range targetNames {
				if _, exists := targets[targetName]; !exists {
					return prepared, fmt.Errorf(
						"Setup error: '%s' list has unknown target '%s'",
							name, targetName)
				}
				targetType := targetScalar
				if isList {
					targetType = targetList
				}
				if targetHolds[targetName] == 0 {
					targetHolds[targetName] = targetType
				} else if targetHolds[targetName] != targetType {
					var msg string
					if isList {
						msg = "scalar"
					} else {
						msg = "list"
					}
					return prepared, fmt.Errorf(
						"Setup error: target '%s' in list '%s' is %s target",
						targetName, name, msg)
				}
			}
			action := symbolAction{act.Mask, targetNames}
			if isList {
				nonterminalActions[act.ListName] = action
			} else {
				terminalActions = append(terminalActions, action)
			}
		}
		prepared[name] = nonterminalPattern{
			nameRequirement: list.NameRequirement,
			terminalActions: terminalActions,
			nonterminalActions: nonterminalActions,
			targets: targets,
			targetOrder: targetOrder,
		}
	}
	return prepared, nil
}

