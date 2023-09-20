// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "potano.misiones/sexp"


type mapConfigType struct {
	mapItemCore
	doc *VectorData
}

func newMapConfig(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	if doc.styler != nil || doc.attester != nil {
		return nil, source.Error("duplicate config section")
	}
	mc := &mapConfigType{doc: doc}
	mc.source = source
	doc.styler = newStyler(doc)
	doc.attester = newAttester(doc)
	return mc, nil
}

func (mc *mapConfigType) ItemType() int {
	return mitConfig
}

func (mc *mapConfigType) setConfigurationItem(newChild mapItemType) error {
	switch item := newChild.(type) {
	case *mapStyleConfigType:
		return mc.doc.styler.setBaseStyle(item)
	case *mapAttestationTypeType:
		return mc.doc.attester.setAttestationType(item)
	default:
		return newChild.Error("unknown config target name")
	}
}







type mapStyleConfigType struct {
	mapItemCore
	itemType int
	properties cssPropertyMap
}

func newMapStyleConfig(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ms := &mapStyleConfigType{}
	ms.source = source
	ms.name = listName
	ms.itemType = nameToTypeMap[listType]
	return ms, nil
}

func (mc *mapStyleConfigType) ItemType() int {
	return mc.itemType
}

func (mc *mapStyleConfigType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	if mc.properties == nil {
		mc.properties = cssPropertyMap{}
	}
	for _, scalar := range scalars {
		key, value, err := decomposeKeyValueScalar(scalar)
		if err != nil {
			return err
		}
		mc.properties[key] = value
	}
	return nil
}







func (att *attester) setAttestationType(atype *mapAttestationTypeType) error {
	groupName := atype.name
	for _, grp := range att.groups {
		if grp.name == groupName {
			return atype.Error("duplicate attribute type %s", groupName)
		}
	}
	groupNum := len(att.groups)
	groupType := atype.groupType
	sumWeights := 0
	var millsPerStep int
	styles := atype.styles
	for _, item := range atype.defs {
		name := item.name
		hasWeight := item.hasWeight
		weight := item.weight
		properties := item.properties
		if _, exists := att.allowedAttestations[name]; exists {
			return item.Error("duplicate definition of attestation %s", name)
		}
		if groupType == 0 {
			if hasWeight {
				groupType = weightedAttestationGroup
			} else {
				groupType = singleValuedAttestationGroup
			}
		}
		if groupType == weightedAttestationGroup {
			if !hasWeight {
				return item.Error("item in weighted attribute group has no weight")
			}
			if len(properties) > 0 {
				return item.Error("weighted attribute item has style")
			}
			sumWeights += item.weight
		} else if groupType == singleValuedAttestationGroup {
			if hasWeight {
				return item.Error(
					"item in single-value attribute group has weight")
			}
			if properties == nil {
				properties = cssPropertyMap{}
			}
			weight = len(styles)
			styles = append(styles, properties)
		}
		att.allowedAttestations[name] = attestationDef{groupNum, weight}
	}
	if groupType == singleValuedAttestationGroup {
		if len(atype.styles) > 0 {
			return atype.Error("single-value attribute group has weighted styles")
		}
	} else if groupType == weightedAttestationGroup {
		if len(styles) == 0 {
			return atype.Error("weighted-value attribute group has no styles")
		}
		if sumWeights == 0 {
			return atype.Error("sum of weights in group is zero")
		}
		millsPerStep = ((1000 * sumWeights) / len(styles)) + 1
		for i, j := 0, len(styles) - 1; i < j; i, j = i+1, j-1 {
			styles[i], styles[j] = styles[j], styles[i]
		}
	}
	groupID := att.doc.styler.setAttestationGroupStyles(styles)
	att.groups = append(att.groups, attestationGroup{
		name: groupName,
		groupType: groupType,
		groupID: groupID,
		sumWeights: sumWeights,
		millsPerStep: millsPerStep,
	})
	return nil
}







type mapAttestationTypeType struct {
	mapItemCore
	groupType int
	defs []*mapAttSymType
	styles []cssPropertyMap
}

func newMapAttestationType(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ma := &mapAttestationTypeType{}
	ma.source = source
	ma.name = listName
	return ma, nil
}

func (ma *mapAttestationTypeType) ItemType() int {
	return mitAttestationType
}

func (ma *mapAttestationTypeType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	sym := scalars[0].String()
	switch sym {
	case "weighted":
		ma.groupType = weightedAttestationGroup
	case "limit1":
		ma.groupType = singleValuedAttestationGroup
	default:
		return scalars[0].Error("unknown group type '%s'", sym)
	}
	return nil
}

func (ma *mapAttestationTypeType) setConfigurationItem(newChild mapItemType) error {
	switch item := newChild.(type) {
	case *mapAttSymType:
		ma.defs = append(ma.defs, item)
	case *mapStyleConfigType:
		ma.styles = append(ma.styles, item.properties)
	}
	return nil
}







type mapAttSymType struct {
	mapItemCore
	hasWeight bool
	weight int
	properties cssPropertyMap
}

func newAttSym(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ma := &mapAttSymType{}
	ma.source = source
	ma.name = listName
	return ma, nil
}

func (ma *mapAttSymType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	for _, scalar := range scalars {
		key, value, err := decomposeKeyValueScalar(scalar)
		if err != nil {
			return err
		}
		if key != "weight" {
			return scalar.Error("item key is not a weight")
		}
		weight, err := value.asInt()
		if err != nil {
			return scalar.Error("weight value is not an integer")
		}
		if weight < 0 {
			return scalar.Error("weight %d is negative", weight)
		}
		if ma.hasWeight {
			return scalar.Error("weight is already set for this attribution type")
		}
		ma.weight = weight
		ma.hasWeight = true
	}
	return nil
}

func (ma *mapAttSymType) setConfigurationItem(newChild mapItemType) error {
	if item, is := newChild.(*mapStyleConfigType); is {
		ma.properties = item.properties
	}
	return nil
}

