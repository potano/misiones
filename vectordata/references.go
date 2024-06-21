// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
        "potano.misiones/sexp"
)

type map_referenceAggregateType struct {
	mapItemCore
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
			mitRectangle, mitCircle, mitRoute, mitSegment}
	case mitPaths:
		acceptable = []int{mitPath, mitPoint, mitMarker, mitCircle, mitSegment, mitRoute}
	case mitSegments:
		acceptable = []int{mitSegment, mitPoint, mitMarker, mitCircle, mitRoute}
	case mitRouteSegments:
		acceptable = []int{mitRoute, mitPoint, mitMarker, mitCircle}
	}
	mr.targets = make([]mapItemType, len(mr.names))
	for i, scalar := range mr.names {
		name := scalar.String()
		var target mapItemType
		if scalar.IsSymbolOrString() {
			target = doc.mapItems[name]
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
				return scalar.Error(
					"referenced item '%s' is a %s type; only %s allowed",
					name, typeMapToName[targetType],
					andMapItemList(acceptable))
			}
			err := target.noteReferrer(mr.parentName, mr)
			if err != nil {
				return err
			}
		} else {
			target = &mapItemCore{name, 0, scalar.Source(), nil}
		}
		mr.targets[i] = target
	}
	if mr.itemType == mitRouteSegments {
		return mr.postcheckRouteSegments(doc)
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


func (mr *map_referenceAggregateType) postcheckRouteSegments(doc *VectorData) error {
	newTargets := make([]mapItemType, 0, len(mr.targets))
	if len(mr.targets) < 1 {
		return mr.Error("no route name given")
	}
	routeName := mr.names[0]
	item := mr.targets[0]
	if item.ItemType() != mitRoute {
		return mr.Error("%s is not a route", routeName)
	}
	newTargets = append(newTargets, item)
	for i := 1; i < len(mr.targets); {
		item := mr.targets[i]
		if item.ItemType() > 0 {
			newTargets = append(newTargets, item)
			i++
		} else if i + 1 < len(mr.targets) && mr.targets[i+1].ItemType() == 0 {
			point, err := newMap_location(doc, mr, "point", "", item.Source())
			if err != nil {
				return err
			}
			err = point.addScalars("", []sexp.LispScalar{mr.names[i], mr.names[i+1]})
			if err != nil {
				return err
			}
			newTargets = append(newTargets, point)
			i += 2
		} else {
			return item.Error("expected a latitude/longitude pair")
		}
	}
	if len(newTargets) != 3 {
		return mr.Error("expected exactly two route endpoints, got %d", len(newTargets))
	}
	// Place the endpoints before and after the route
	newTargets[0], newTargets[1] = newTargets[1], newTargets[0]
	mr.targets = newTargets
	return nil
}

