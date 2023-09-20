// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
        "potano.misiones/sexp"
)

type map_referenceAggregateType struct {
	mapItemCore
	itemType int
	parentName string
	names []sexp.LispScalar
	targets []mapItemType
}

func newMap_referenceAggregate(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mr := &map_referenceAggregateType{}
	mr.source = source
	mr.itemType = nameToTypeMap[listType]
	if mr.itemType == 0 {
		return nil, source.Error("unknown type '%s'", listType)
	}
	if parent != nil {
		mr.parentName = parent.Name()
	}
	doc.registerReferenceItem(mr)
	return mr, nil
}

func (mr *map_referenceAggregateType) ItemType() int {
	return mr.itemType
}

func (mr *map_referenceAggregateType) addScalars(targetName string,
		scalars []sexp.LispScalar) error {
	mr.names = scalars
	return nil
}

func (mr *map_referenceAggregateType) resolveTargets(doc *VectorData) error {
	var acceptable []int
	switch (mr.itemType) {
	case mitLayers:
		acceptable = []int{mitLayer}
	case mitFeatures:
		acceptable = []int{mitFeature, mitMarker, mitPoint, mitPath, mitPolygon,
			mitRectangle, mitCircle, mitRoute}
	case mitPaths:
		acceptable = []int{mitPath}
	case mitSegments:
		acceptable = []int{mitSegment}
	}
	mr.targets = make([]mapItemType, len(mr.names))
	for i, scalar := range mr.names {
		name := scalar.String()
		target := doc.mapItems[name]
		if target == nil {
			return scalar.Error("name '%s' is not registered", name)
		}
		targetType := target.ItemType()
		ok := false
		for _, tp := range acceptable {
			if tp == targetType {
				ok = true
				break
			}
		}
		if !ok {
			return scalar.Error("referenced item '%s' is a %s type; only %s allowed",
				name, typeMapToName[targetType], andMapItemList(acceptable))
		}
		err := target.noteReferrer(mr.parentName, mr)
		if err != nil {
			return err
		}
		mr.targets[i] = target
	}
	return nil
}


func andMapItemList(list []int) string {
	sList := make([]string, len(list))
	for i, val := range list {
		sList[i] = typeMapToName[val]
	}
	return andList(sList)
}


func andList(list []string) string {
	var out string
	numWords := len(list)
	for i, word := range list {
		if i < numWords - 1 {
			if numWords > 2 {
				out += word + ", "
			} else {
				out += word + " "
			}
		} else if numWords > 1 {
			out += "and " + word
		} else {
			out = word
		}
	}
	return out
}

