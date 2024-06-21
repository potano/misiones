// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"testing"
)



type namedStylerCheck struct {
	name string
	properties [][]string
}

type referencedStyleMapCheck struct {
	key []byte
	contentKey string
}

func checkStylerConfig(T *testing.T, sty *styler, base []namedStylerCheck, attStyles [][][][]string,
		referenced []namedStylerCheck, mapcheck []referencedStyleMapCheck) {
	T.Helper()
	checkNamedStylerMap(T, sty.baseStyles, sty.baseStyleMap, base, "base styles")
	if len(sty.attestationStyles) != len(attStyles) {
		T.Fatalf("have %d weighted attestation groups, have %d",
			len(sty.attestationStyles), len(attStyles))
	}
	for groupID, steps := range attStyles {
		if len(sty.attestationStyles[groupID]) != len(steps) {
			T.Fatalf("have %d steps in weighted attestation group %d, want %d",
				len(sty.attestationStyles[groupID]), groupID, len(steps))
		}
		for stepNum, pairs := range steps {
			checkStyleProperties(T, sty.attestationStyles[groupID][stepNum], pairs,
				fmt.Sprintf("step %d of attestation group %d", stepNum, groupID))
		}
	}
	checkNamedStylerMap(T, sty.referencedStyles, sty.referencedStyleMapByContent, referenced,
		"referenced styles")
	if len(sty.referencedStyleMap) != len(mapcheck) {
		T.Fatalf("have %d referenceStyleMap entries, want %d", len(sty.referencedStyleMap),
			len(mapcheck))
	}
	for _, rec := range mapcheck {
		index, exists := sty.referencedStyleMap[string(rec.key)]
		if !exists {
			T.Fatalf("reference key %v does not exist", rec.key)
		}
		byContentIndex, exists := sty.referencedStyleMapByContent[rec.contentKey]
		if !exists && len(rec.contentKey) > 0 {
			T.Fatalf("reference key %v: no map entry for content key %s", rec.key,
				rec.contentKey)
		}
		if index != byContentIndex {
			T.Fatalf("reference key %v: expected index %d, got %d", rec.key,
				byContentIndex, index)
		}
	}
}

func checkNamedStylerMap(T *testing.T, haveVec []cssPropertyMap, haveMap map[string]int,
		want []namedStylerCheck, desc string) {
	T.Helper()
	if len(haveMap) != len(want) {
		T.Fatalf("have %d %s map items, want %d", len(haveMap), desc, len(want))
	}
	if len(haveVec) != len(want) + 1 {
		T.Fatalf("have %d %s, want %d", len(haveVec), desc, len(want) + 1)
	}
	for _, wantCheck := range want {
		styleName := wantCheck.name
		haveProperties := haveVec[haveMap[styleName]]
		checkStyleProperties(T, haveProperties, wantCheck.properties,
			printableString(desc + " " + styleName))
	}
}

func printableString(str string) string {
	bs := []byte(str)
	for i := len(bs) - 1; i >= 0; i-- {
		if bs[i] < ' ' {
			bs = append(bs, []byte{0, 0, 0}...)
			for j := len(bs) - 4; j > i; j-- {
				bs[j+3] = bs[j]
			}
			s := []byte(fmt.Sprintf("\\x%02X", bs[i]))
			copy(bs[i:], s)
		}
	}
	return string(bs)
}

func checkStyleProperties(T *testing.T, have cssPropertyMap, want [][]string, desc string) {
	desc = printableString(desc)
	if have == nil {
		T.Fatalf("%s not found", desc)
	}
	if len(have) != len(want) {
		T.Fatalf("have %d %s, want %d", len(have), desc, len(want))
	}
	for _, pair := range want {
		key := pair[0]
		value :=  pair[1]
		if tstval, exists := have[key]; !exists {
			T.Fatalf("no '%s' %s property found", key, desc)
		} else if string(tstval) != value {
			T.Fatalf("%s property '%s': want '%s', got '%s'", desc, key,
				value, string(tstval))
		}
	}
}



type allowedAttestationCheck struct {
	name string
	groupNum, weight int
}

func checkAttesterConfig(T *testing.T, att *attester, groups []attestationGroup,
		allowedAttestations []allowedAttestationCheck) {
	if len(att.groups) != len(groups) {
		T.Fatalf("expected %d attestation groups, got %d", len(groups), len(att.groups))
	}
	for groupNum, group := range groups {
		attGroup := att.groups[groupNum]
		if attGroup.name != group.name {
			T.Fatalf("expected attestation group %d to have name '%s', got '%s'",
				groupNum, group.name, attGroup.name)
		}
		if attGroup.groupType != group.groupType {
			T.Fatalf("expected attestation group %d to have type %d, got %d",
				groupNum, group.groupType, attGroup.groupType)
		}
		if attGroup.groupID != group.groupID {
			T.Fatalf("expected attestation group %d to have ID %d, got %d",
				groupNum, group.groupID, attGroup.groupID)
		}
		if attGroup.sumWeights != group.sumWeights {
			T.Fatalf("expected attestation group %d to have Σ weights %d, got %d",
				groupNum, group.sumWeights, attGroup.sumWeights)
		}
		if attGroup.millsPerStep != group.millsPerStep {
			T.Fatalf("expected attestation group %d to have %d mills/step, got %d",
				groupNum, group.millsPerStep, attGroup.millsPerStep)
		}
	}
	if len(att.allowedAttestations) != len(allowedAttestations) {
		T.Fatalf("expected %d enumerated attestations, got %d", len(allowedAttestations),
			len(att.allowedAttestations))
	}
	for _, check := range allowedAttestations {
		name := check.name
		if def, exists := att.allowedAttestations[name]; !exists {
			T.Fatalf("enumerated attestation '%s' does not exists", name)
		} else {
			if def.groupNum != check.groupNum {
				T.Fatalf("attestation '%s' expected to be group %d, not %d",
					name, check.groupNum, def.groupNum)
			}
			if def.weight != check.weight {
				T.Fatalf("attestation '%s' expected to have weight %d, got %d",
					name, check.weight, def.weight)
			}
		}
	}
}



func Test_routeSegmentPathsConfig(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hill theRoad)
		)
	)
	(feature hill
		(popup "Trail along")
		(style baseStyle)
		(path likely
		      29.50 -83.43  29.50 -83.41
		      29.40 -83.43  29.40 -83.41)
	)
	(route theRoad
		(style roadStyle)
		(segment
			(attestation book maybe)
			(paths path1 path2)
		)
		(segment
			(attestation book magazine forSure)
			(paths path3)
		)
	)
	(path path1
		(attestation modern_path)
		30.350075 -83.507595
		30.350177 -83.507918
		30.351014 -83.513659
		30.351541 -83.517636
	)
	(path path2
		(attestation modern_name)
		30.351541 -83.517636
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
	)
	(path path3
		(attestation modern_path)
		30.351842 -83.520299
		30.351850 -83.520426
		30.351861 -83.520554
		30.351870 -83.520668
		30.351879 -83.520762
	)
	(config
		(baseStyle baseStyle
			"color=#1f78b4"
                        "opacity=0.9"
                        "fill=true"
                        "fillColor=#1f78b4"
                        "fillOpacity=0.1"
		)
		(baseStyle roadStyle
		        "color=#AA3333"
                        "opacity=0.6"
                        "width=5"
		)
		(attestationType plusgood weighted
			(attSym book  "weight=2")
			(attSym magazine "weight=1")
			(modStyle "opacity=0.8" "width=3")
			(modStyle "opacity=0.8" "width=2")
		)
		(attestationType manifestation limit1
			(attSym modern_name)
			(attSym modern_path (modStyle "dashArray=4 4"))
		)
		(attestationType confidence limit1
			(attSym forSure)
			(attSym maybe (modStyle "opacity=0.4"))
		)
	)
	`
	vd := prepareAndParseStringsOnly(T, sourceText)

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hill theRoad
      →feature 'hill' @ infile0:7
          popup text: 'Trail along'
          style: baseStyle
        →path 'likely' @ infile0:10
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
      →route 'theRoad' @ infile0:14
          style: roadStyle
        →segment '$5' @ infile0:16
            attestation: book maybe
          →paths '' @ infile0:18
              parent: $5
              target names: path1 path2
            →path 'path1' @ infile0:25
                attestation: modern_path
                location: 30.350075  -83.507595
                          30.350177  -83.507918
                          30.351014  -83.513659
                          30.351541  -83.517636
            →path 'path2' @ infile0:32
                attestation: modern_name
                location: 30.351541  -83.517636
                          30.351709  -83.519064
                          30.351815  -83.519952
                          30.351830  -83.520140
                          30.351842  -83.520299
        →segment '$6' @ infile0:20
            attestation: book magazine forSure
          →paths '' @ infile0:22
              parent: $6
              target names: path3
            →path 'path3' @ infile0:40
                attestation: modern_path
                location: 30.351842  -83.520299
                          30.351850  -83.520426
                          30.351861  -83.520554
                          30.351870  -83.520668
                          30.351879  -83.520762`)
	checkStylerConfig(T, vd.styler,
		[]namedStylerCheck{
			{"baseStyle", [][]string{
				[]string{"color", "#1f78b4"},
				[]string{"opacity", "0.9"},
				[]string{"fill", "true"},
				[]string{"fillColor", "#1f78b4"},
				[]string{"fillOpacity", "0.1"},
			}},
			{"roadStyle", [][]string{
				[]string{"color", "#AA3333"},
				[]string{"opacity", "0.6"},
				[]string{"width", "5"},
			}},
		},
		[][][][]string{
			[][][]string{
				[][]string{
					[]string{"opacity", "0.8"},
					[]string{"width", "2"},
				},
				[][]string{
					[]string{"opacity", "0.8"},
					[]string{"width", "3"},
				},
			},
			[][][]string{
				[][]string{
				},
				[][]string{
					[]string{"dashArray", "4 4"},
				},
			},
			[][][]string{
				[][]string{
				},
				[][]string{
					[]string{"opacity", "0.4"},
				},
			},
		},
		nil, nil)
	checkAttesterConfig(T, vd.attester,
		[]attestationGroup{
			{ "plusgood", weightedAttestationGroup, 0, 3, 1501 },
			{ "manifestation", singleValuedAttestationGroup, 1, 0, 0 },
			{ "confidence", singleValuedAttestationGroup, 2, 0, 0 },
		},
		[]allowedAttestationCheck{
			{ "book", 0, 2 },
			{ "magazine", 0, 1 },
			{ "modern_name", 1, 0 },
			{ "modern_path", 1, 1 },
			{ "forSure", 2, 0 },
			{ "maybe", 2, 1 },
		},
	)
}


func Test_routeSegmentPathsConfigStyled(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hill theRoad)
		)
	)
	(feature hill
		(popup "Trail along")
		(style baseStyle)
		(path likely
		      29.50 -83.43  29.50 -83.41
		      29.40 -83.43  29.40 -83.41)
	)
	(route theRoad
		(style roadStyle)
		(segment
			(attestation book maybe)
			(paths path1 path2)
		)
		(segment
			(attestation book magazine forSure)
			(paths path3)
		)
	)
	(path path1
		(attestation modern_path)
		30.350075 -83.507595
		30.350177 -83.507918
		30.351014 -83.513659
		30.351541 -83.517636
	)
	(path path2
		(attestation modern_name)
		30.351541 -83.517636
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
	)
	(path path3
		(attestation modern_path)
		30.351842 -83.520299
		30.351850 -83.520426
		30.351861 -83.520554
		30.351870 -83.520668
		30.351879 -83.520762
	)
	(config
		(baseStyle baseStyle
			"color=#1f78b4"
                        "opacity=0.9"
                        "fill=true"
                        "fillColor=#1f78b4"
                        "fillOpacity=0.1"
		)
		(baseStyle roadStyle
		        "color=#AA3333"
                        "opacity=0.6"
                        "width=5"
		)
		(attestationType plusgood weighted
			(attSym book  "weight=2")
			(attSym magazine "weight=1")
			(modStyle "opacity=0.8" "width=3")
			(modStyle "opacity=0.8" "width=2")
		)
		(attestationType manifestation limit1
			(attSym modern_name)
			(attSym modern_path (modStyle "dashArray=4 4"))
		)
		(attestationType confidence limit1
			(attSym forSure)
			(attSym maybe (modStyle "opacity=0.4"))
		)
	)
	`
	vd := prepareAndParseStringsOnly(T, sourceText)

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hill theRoad
      →feature 'hill' @ infile0:7
          popup text: 'Trail along'
          style: baseStyle
        →path 'likely' @ infile0:10
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
      →route 'theRoad' @ infile0:14
          style: roadStyle
        →segment '$5' @ infile0:16
            attestation: book maybe
          →paths '' @ infile0:18
              parent: $5
              target names: path1 path2
            →path 'path1' @ infile0:25
                attestation: modern_path
                location: 30.350075  -83.507595
                          30.350177  -83.507918
                          30.351014  -83.513659
                          30.351541  -83.517636
            →path 'path2' @ infile0:32
                attestation: modern_name
                location: 30.351541  -83.517636
                          30.351709  -83.519064
                          30.351815  -83.519952
                          30.351830  -83.520140
                          30.351842  -83.520299
        →segment '$6' @ infile0:20
            attestation: book magazine forSure
          →paths '' @ infile0:22
              parent: $6
              target names: path3
            →path 'path3' @ infile0:40
                attestation: modern_path
                location: 30.351842  -83.520299
                          30.351850  -83.520426
                          30.351861  -83.520554
                          30.351870  -83.520668
                          30.351879  -83.520762`)

	err := vd.CheckInStylesAndAttestations()
	if err != nil {
		T.Fatal(err.Error())
	}

	checkStylerConfig(T, vd.styler,
		[]namedStylerCheck{
			{"baseStyle", [][]string{
				[]string{"color", "#1f78b4"},
				[]string{"opacity", "0.9"},
				[]string{"fill", "true"},
				[]string{"fillColor", "#1f78b4"},
				[]string{"fillOpacity", "0.1"},
			}},
			{"roadStyle", [][]string{
				[]string{"color", "#AA3333"},
				[]string{"opacity", "0.6"},
				[]string{"width", "5"},
			}},
		},
		[][][][]string{
			[][][]string{
				[][]string{
					[]string{"opacity", "0.8"},
					[]string{"width", "2"},
				},
				[][]string{
					[]string{"opacity", "0.8"},
					[]string{"width", "3"},
				},
			},
			[][][]string{
				[][]string{
				},
				[][]string{
					[]string{"dashArray", "4 4"},
				},
			},
			[][][]string{
				[][]string{
				},
				[][]string{
					[]string{"opacity", "0.4"},
				},
			},
		},
		[]namedStylerCheck{
			{"color:\"#1f78b4\"fill:truefillColor:\"#1f78b4\"fillOpacity:0.1opacity:0.9",
				[][]string{
				[]string{"color", "#1f78b4"},
				[]string{"opacity", "0.9"},
				[]string{"fill", "true"},
				[]string{"fillColor", "#1f78b4"},
				[]string{"fillOpacity", "0.1"},
			}},
			{"color:\"#AA3333\"opacity:0.6width:5", [][]string{
				[]string{"color", "#AA3333"},
				[]string{"opacity", "0.6"},
				[]string{"width", "5"},
			}},
			{"opacity:0.8width:3", [][]string{
				[]string{"opacity", "0.8"},
				[]string{"width", "3"},
			}},
			{"opacity:0.4width:3", [][]string{
				[]string{"opacity", "0.4"},
				[]string{"width", "3"},
			}},
			{"dashArray:\"4 4\"", [][]string{
				[]string{"dashArray", "4 4"},
			}},
		},
		[]referencedStyleMapCheck {
			{[]byte{1}, "color:\"#1f78b4\"fill:truefillColor:\"#1f78b4\"fillOpacity:0.1opacity:0.9"},
			{[]byte{2}, "color:\"#AA3333\"opacity:0.6width:5"},
			{[]byte{0, 0, 1, 0}, ""},
			{[]byte{0, 0, 2, 0}, "dashArray:\"4 4\""},
			{[]byte{0, 2, 0, 1}, "opacity:0.8width:3"},
			{[]byte{0, 2, 0, 2}, "opacity:0.4width:3"},
		})

	checkAttesterConfig(T, vd.attester,
		[]attestationGroup{
			{ "plusgood", weightedAttestationGroup, 0, 3, 1501 },
			{ "manifestation", singleValuedAttestationGroup, 1, 0, 0 },
			{ "confidence", singleValuedAttestationGroup, 2, 0, 0 },
		},
		[]allowedAttestationCheck{
			{ "book", 0, 2 },
			{ "magazine", 0, 1 },
			{ "modern_name", 1, 0 },
			{ "modern_path", 1, 1 },
			{ "forSure", 2, 0 },
			{ "maybe", 2, 1 },
		},
	)
}


func Test_routeSegmentPathsConfigStyledAttestationModifiesStyle(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hill theRoad)
		)
	)
	(feature hill
		(popup "Trail along")
		(style baseStyle)
		(path likely
		      29.50 -83.43  29.50 -83.41
		      29.40 -83.43  29.40 -83.41)
	)
	(route theRoad
		(segment
			(style roadStyle)
			(attestation book maybe)
			(paths path1 path2)
		)
		(segment
			(attestation book magazine forSure)
			(paths path3)
		)
	)
	(path path1
		(attestation modern_path)
		30.350075 -83.507595
		30.350177 -83.507918
		30.351014 -83.513659
		30.351541 -83.517636
	)
	(path path2
		(attestation modern_name)
		30.351541 -83.517636
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
	)
	(path path3
		(attestation modern_path)
		30.351842 -83.520299
		30.351850 -83.520426
		30.351861 -83.520554
		30.351870 -83.520668
		30.351879 -83.520762
	)
	(config
		(baseStyle baseStyle
			"color=#1f78b4"
                        "opacity=0.9"
                        "fill=true"
                        "fillColor=#1f78b4"
                        "fillOpacity=0.1"
		)
		(baseStyle roadStyle
		        "color=#AA3333"
                        "opacity=0.6"
                        "width=5"
		)
		(attestationType plusgood weighted
			(attSym book  "weight=2")
			(attSym magazine "weight=1")
			(modStyle "opacity=0.8" "width=3")
			(modStyle "opacity=0.8" "width=2")
		)
		(attestationType manifestation limit1
			(attSym modern_name)
			(attSym modern_path (modStyle "dashArray=4 4"))
		)
		(attestationType confidence limit1
			(attSym forSure)
			(attSym maybe (modStyle "opacity=0.4"))
		)
	)
	`
	vd := prepareAndParseStringsOnly(T, sourceText)

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hill theRoad
      →feature 'hill' @ infile0:7
          popup text: 'Trail along'
          style: baseStyle
        →path 'likely' @ infile0:10
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
      →route 'theRoad' @ infile0:14
        →segment '$5' @ infile0:15
            style: roadStyle
            attestation: book maybe
          →paths '' @ infile0:18
              parent: $5
              target names: path1 path2
            →path 'path1' @ infile0:25
                attestation: modern_path
                location: 30.350075  -83.507595
                          30.350177  -83.507918
                          30.351014  -83.513659
                          30.351541  -83.517636
            →path 'path2' @ infile0:32
                attestation: modern_name
                location: 30.351541  -83.517636
                          30.351709  -83.519064
                          30.351815  -83.519952
                          30.351830  -83.520140
                          30.351842  -83.520299
        →segment '$6' @ infile0:20
            attestation: book magazine forSure
          →paths '' @ infile0:22
              parent: $6
              target names: path3
            →path 'path3' @ infile0:40
                attestation: modern_path
                location: 30.351842  -83.520299
                          30.351850  -83.520426
                          30.351861  -83.520554
                          30.351870  -83.520668
                          30.351879  -83.520762`)

	err := vd.CheckInStylesAndAttestations()
	if err != nil {
		T.Fatal(err.Error())
	}

	checkStylerConfig(T, vd.styler,
		[]namedStylerCheck{
			{"baseStyle", [][]string{
				[]string{"color", "#1f78b4"},
				[]string{"opacity", "0.9"},
				[]string{"fill", "true"},
				[]string{"fillColor", "#1f78b4"},
				[]string{"fillOpacity", "0.1"},
			}},
			{"roadStyle", [][]string{
				[]string{"color", "#AA3333"},
				[]string{"opacity", "0.6"},
				[]string{"width", "5"},
			}},
		},
		[][][][]string{
			[][][]string{
				[][]string{
					[]string{"opacity", "0.8"},
					[]string{"width", "2"},
				},
				[][]string{
					[]string{"opacity", "0.8"},
					[]string{"width", "3"},
				},
			},
			[][][]string{
				[][]string{
				},
				[][]string{
					[]string{"dashArray", "4 4"},
				},
			},
			[][][]string{
				[][]string{
				},
				[][]string{
					[]string{"opacity", "0.4"},
				},
			},
		},
		[]namedStylerCheck{
			{"color:\"#1f78b4\"fill:truefillColor:\"#1f78b4\"fillOpacity:0.1opacity:0.9",
				[][]string{
				[]string{"color", "#1f78b4"},
				[]string{"opacity", "0.9"},
				[]string{"fill", "true"},
				[]string{"fillColor", "#1f78b4"},
				[]string{"fillOpacity", "0.1"},
			}},
			{"color:\"#AA3333\"opacity:0.4width:3", [][]string{
				[]string{"color", "#AA3333"},
				[]string{"opacity", "0.4"},
				[]string{"width", "3"},
			}},
			{"opacity:0.8width:3", [][]string{
				[]string{"opacity", "0.8"},
				[]string{"width", "3"},
			}},
			{"dashArray:\"4 4\"", [][]string{
				[]string{"dashArray", "4 4"},
			}},
		},
		[]referencedStyleMapCheck {
			{[]byte{1}, "color:\"#1f78b4\"fill:truefillColor:\"#1f78b4\"fillOpacity:0.1opacity:0.9"},
			{[]byte{2, 2, 0, 2}, "color:\"#AA3333\"opacity:0.4width:3"},
			{[]byte{0, 2, 0, 1}, "opacity:0.8width:3"},
			{[]byte{0, 0, 1, 0}, ""},
			{[]byte{0, 0, 2, 0}, "dashArray:\"4 4\""},
		})

	checkAttesterConfig(T, vd.attester,
		[]attestationGroup{
			{ "plusgood", weightedAttestationGroup, 0, 3, 1501 },
			{ "manifestation", singleValuedAttestationGroup, 1, 0, 0 },
			{ "confidence", singleValuedAttestationGroup, 2, 0, 0 },
		},
		[]allowedAttestationCheck{
			{ "book", 0, 2 },
			{ "magazine", 0, 1 },
			{ "modern_name", 1, 0 },
			{ "modern_path", 1, 1 },
			{ "forSure", 2, 0 },
			{ "maybe", 2, 1 },
		},
	)
}

