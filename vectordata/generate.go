// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"strings"
	"strconv"
)


func (vd *VectorData) GenerateJs() (string, error) {
	json, err := vd.generateJson()
	return "allData=" + json, err
}

func (vd *VectorData) generateJson() (string, error) {
	if !vd.styler.styleCheckRun() {
		err := vd.CheckInStylesAndAttestations()
		if err != nil {
			return "", err
		}
	}
	jsg := jsGenerator{
		vd: vd,
		styles: newGenGroup("styles", len(vd.styler.referencedStyles)),
		menuitems: newGenGroup("menuitems", 10),
		texts: newGenGroup("texts", 30),
		features: newGenGroup("features", len(vd.mapItems)),
		points: newPointsGroup("points", len(vd.mapItems)),
	}
	if vd.styler != nil {
		jsg.styles.addEntry("0")
		vd.styler.serializeStyles(jsg)
	}
	jsg.texts.addEntry("0")
	err := jsg.serializeFromRoot()
	if err != nil {
		return "", err
	}
	blobs := []string{jsg.styles.json(), jsg.menuitems.json(), jsg.texts.json(),
		jsg.features.json(), jsg.points.json()}
	return "{" + strings.Join(blobs, ",") + "}", nil
}



type jsGenerator struct {
	vd *VectorData
	styles, menuitems, texts, features *genGroup
	points *pointsGroup
}

type genGroup struct {
	key string
	blobs []string
	indices map[string]int
}

func newGenGroup(key string, initSize int) *genGroup {
	gg := &genGroup{key, make([]string, 0, initSize + 1), make(map[string]int, initSize)}
	return gg
}

func (gg *genGroup) allocEntry() int {
	gg.blobs = append(gg.blobs, "")
	return len(gg.blobs) - 1
}

func (gg *genGroup) allocEntryWithKey(key string) int {
	index := len(gg.blobs)
	gg.blobs = append(gg.blobs, "")
	gg.indices[key] = index
	return index
}

func (gg *genGroup) index(key string) (int, bool) {
	i, b := gg.indices[key]
	return i, b
}

func (gg *genGroup) addEntry(text string) int {
	gg.blobs = append(gg.blobs, text)
	return len(gg.blobs) - 1
}

func (gg *genGroup) addStringWithKey(key, text string) int {
	index, exists := gg.indices[key]
	if exists {
		return index
	}
	index = len(gg.blobs)
	gg.blobs = append(gg.blobs, strconv.Quote(text))
	gg.indices[key] = index
	return index
}

func (gg *genGroup) setEntry(index int, text string) {
	gg.blobs[index] = text
}

func (gg *genGroup) setJsObject(index int, args ...any) {
	gg.setEntry(index, generateJsObject(args...))
}

func (gg *genGroup) setAll(all []string) {
	gg.blobs = all
}

func (gg *genGroup) json() string {
	return "\"" + gg.key + "\":[" + strings.Join(gg.blobs, ",") + "]"
}


type pointsGroup = genGroup

func newPointsGroup(key string, initSize int) *pointsGroup {
	initSize *= 2 * 20
	return &pointsGroup{key, make([]string, 0, initSize), make(map[string]int, initSize)}
}

func (pg *pointsGroup) addPoints(name string, pairs locationPairs) int {
	index := len(pg.blobs)
	slice := make([]string, len(pairs))
	for i, v := range pairs {
		slice[i] = v.String()
	}
	pg.blobs = append(pg.blobs, slice...)
	pg.indices[name] = index
	return index
}





func (jsg jsGenerator) serializeFromRoot() error {
	for _, item := range jsg.vd.layersRoot.layers {
		index := jsg.menuitems.allocEntry()
		features, err := jsg.resolveFeatures(item.features)
		if err != nil {
			return err
		}
		jsg.menuitems.setJsObject(index, "menuitem", item.menuitem, "f", features)
	}
	return nil
}


func (jsg jsGenerator) resolveFeatures(list []mapItemType) ([]int, error) {
	resolved := make([]int, 0, len(list))
	for _, child := range list {
		if ref, is := child.(*map_referenceAggregateType); is {
			group, err := jsg.resolveFeatures(ref.targets)
			if err != nil {
				return resolved, err
			}
			resolved = append(resolved, group...)
		} else if child.ItemType() != mitPoint {
			index, _ := jsg.features.index(child.Name())
			if index == 0 {
				index = jsg.features.allocEntryWithKey(child.Name())
				serialized, err := jsg.serializeMapItem(child)
				if err != nil {
					return resolved, err
				}
				jsg.features.setEntry(index, serialized)
			}
			resolved = append(resolved, index)
		}
	}
	return resolved, nil
}


func (jsg jsGenerator) serializeMapItem(item mapItemType) (string, error) {
	t := item.ItemTypeString()
	var popup nonZeroInt
	var features []mapItemType
	style := nonZeroInt(jsg.vd.styler.styleIndex(item))
	switch item := item.(type) {
	case *mapFeatureType:
		popup = item.popup.textIndex(jsg)
		features = item.features
	case *mapRouteOrSegmentType:
		popup = item.popup.textIndex(jsg)
		features = item.children
	case *map_locationType:
		popup = item.popup.textIndex(jsg)
		protoLocation := item.prototypePath
		if protoLocation == nil {
			protoLocation = item
		}
		bigOffset, exists := jsg.points.index(protoLocation.Name())
		if !exists {
			bigOffset = jsg.points.addPoints(protoLocation.Name(),
				protoLocation.location)
		}
		if item.itemType == mitCircle {
			asPixels := item.radiusType == mitPixels
			return generateJsObject(
				"t", t,
				"popup", popup,
				"style", style,
				"asPixels", asPixels,
				"radius", item.radius,
				"loc", []int{bigOffset, 2}), nil
		}
		return generateJsObject(
			"t", t,
			"popup", popup,
			"style", style,
			"html", nonEmptyString(item.html),
			"loc", []int{bigOffset + int(item.offsetInPrototype), len(item.location)},
		), nil
	default:
		return "", fmt.Errorf("unhandled item type %s", t)
	}
	indices, err := jsg.resolveFeatures(features)
	if err != nil {
		return "", err
	}
	return generateJsObject(
		"t", t,
		"popup", popup,
		"style", style,
		"f", indices,
	), nil
}


func (mp *mapPopupType) textIndex(jsg jsGenerator) nonZeroInt {
	if mp == nil {
		return nonZeroInt(0)
	}
	return nonZeroInt(jsg.texts.addStringWithKey(mp.text, mp.text))
}



type nonEmptyString string
type nonZeroInt int



func generateJsObject(args ...any) string {
	entries := []string{}
	for i := 0; i < len(args) - 1; i += 2 {
		key := args[i]
		val := args[i+1]
		var str string
		switch v := val.(type) {
		case int:
			str = strconv.Itoa(v)
		case float64:
			str = strconv.FormatFloat(v, 'f', 6, 64)
		case locAngleType:
			str = v.String()
		case string:
			str = strconv.Quote(v)
		case nonEmptyString:
			if len(v) == 0 {
				continue
			}
			str = strconv.Quote(string(v))
		case nonZeroInt:
			if v == 0 {
				continue
			}
			str = strconv.Itoa(int(v))
		case []int:
			if len(v) == 0 {
				continue
			}
			items := make([]string, len(v))
			for i, val := range v {
				items[i] = strconv.Itoa(val)
			}
			str = "[" + strings.Join(items, ",") + "]"
		case bool:
			str = strconv.FormatBool(v)
		default:
			continue
		}
		entries = append(entries, strconv.Quote(key.(string)) + ":" + str)
	}
	return "{" + strings.Join(entries, ",") + "}"
}

