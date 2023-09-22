// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"sort"
	"strings"

	"potano.misiones/sexp"
)


type styler struct {
	doc *VectorData
	baseStyles []cssPropertyMap
	baseStyleMap map[string]int
	attestationStyles [][]cssPropertyMap
	referencedStyles []cssPropertyMap
	referencedStyleMap map[string]int
}


func newStyler(doc *VectorData) *styler {
	return &styler{
		doc: doc,
		baseStyles: []cssPropertyMap{nil},
		baseStyleMap: map[string]int{},
		attestationStyles: [][]cssPropertyMap{},
		referencedStyles: []cssPropertyMap{nil},
		referencedStyleMap: map[string]int{},
	}
}

func (sty *styler) setBaseStyle(baseStyle *mapStyleConfigType) error {
	if _, exists := sty.baseStyleMap[baseStyle.name]; exists {
		return baseStyle.Error("redefinition of base style %s", baseStyle.name)
	} else {
		sty.baseStyles = append(sty.baseStyles, baseStyle.properties)
		sty.baseStyleMap[baseStyle.name] = len(sty.baseStyles) - 1
	}
	return nil
}

func (sty *styler) setAttestationGroupStyles(styles []cssPropertyMap) int {
	groupID := len(sty.attestationStyles)
	sty.attestationStyles = append(sty.attestationStyles, styles)
	return groupID
}

func (sty *styler) styleCheckRun() bool {
	return len(sty.referencedStyles) > 1
}

func (sty *styler) checkConfiguration() error {
	return nil
}

func (sty *styler) baseStyleIndex(style *mapStyleType) (int, error) {
	name := style.name
	if styX, exists := sty.baseStyleMap[name]; !exists {
		return 0, style.Error("unknown style '%s'", name)
	} else {
		return styX, nil
	}
}

func (sty *styler) resolveStyle(style *mapStyleType) error {
	styX, err := sty.baseStyleIndex(style)
	if err != nil {
		return err
	}
	key := string([]byte{byte(styX)})
	rsX, exists := sty.referencedStyleMap[key]
	if !exists {
		rsX = len(sty.referencedStyles)
		sty.referencedStyles = append(sty.referencedStyles, sty.baseStyles[styX])
		sty.referencedStyleMap[key] = rsX
	}
	style.resolvedStyleIndex = rsX
	return nil
}

func (sty *styler) findAttestationStyle(styX int, atypeVector []int) int {
	keyBytes := make([]byte, len(atypeVector) + 1)
	keyBytes[0] = byte(styX)
	for i, v := range atypeVector {
		keyBytes[i+1] = byte(v)
	}
	key := string(keyBytes)
	rsX, exists := sty.referencedStyleMap[key]
	if !exists {
		props := cssPropertyMap{}
		if styX > 0 {
			for k, v := range sty.baseStyles[styX] {
				props[k] = v
			}
		}
		for groupID, step := range atypeVector {
			for k, v := range sty.attestationStyles[groupID][step] {
				props[k] = v
			}
		}
		rsX = len(sty.referencedStyles)
		sty.referencedStyles = append(sty.referencedStyles, props)
		sty.referencedStyleMap[key] = rsX
	}
	return rsX
}

func (sty *styler) generateJs() string {
	if len(sty.referencedStyles) < 2 {
		return ""
	}
	blobs := make([]string, len(sty.referencedStyles) - 1)
	for i, props := range sty.referencedStyles[1:] {
		parts := make([]string, 0, len(props))
		for k, v := range props {
			parts = append(parts, "\"" + k + "\":" + v.jsonForm())
		}
		sort.Strings(parts)
		blobs[i] = formStyleName(i+1) + "={" + strings.Join(parts, ",") + "}"
	}
	return "var " + strings.Join(blobs, ",")
}

func formStyleName(styleIndex int) string {
	return fmt.Sprintf("$s%d", styleIndex)
}




type mapStyleType struct {
	mapItemCore
	resolvedStyleIndex int
}

func newMapStyle(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ms := &mapStyleType{}
	ms.source = source
	return ms, nil
}

func (ms *mapStyleType) ItemType() int {
	return mitStyle
}

func (ms *mapStyleType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	name := scalars[0].String()
	ms.name = name
	return nil
}

