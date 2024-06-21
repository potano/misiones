// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"encoding/json"
	"fmt"
	"testing"
)


func checkGeneratedJson(T *testing.T, source string,
		styles, texts, menuitems, features, points []any) {
	var doc interface{}
	err := json.Unmarshal([]byte(source), &doc)
	if err != nil {
		T.Fatal(err.Error())
	}
	docobj, is := doc.(map[string]any)
	if !is {
		T.Fatal("JSON blob does not correspond to an object")
	}
	checkGenByKey(T, docobj, "styles", "style", styles)
	checkGenByKey(T, docobj, "texts", "text", texts)
	checkGenByKey(T, docobj, "menuitems", "menuitem", menuitems)
	checkGenByKey(T, docobj, "features", "feature", features)
	checkGenByKey(T, docobj, "points", "point", points)
}

func checkGenByKey(T *testing.T, doc map[string]any, key, sing string, wantGroup any) {
	if item, exists := doc[key]; !exists {
		T.Fatalf("JSON object has no '%s' member", key)
	} else {
		checkAnyValue(T, item, sing, wantGroup)
	}
}

func checkAnyValue(T *testing.T, have any, desc string, want any) {
	switch want := want.(type) {
	case bool:
		if val, is := have.(bool); !is {
			T.Fatalf("%s: expected boolean, got %v", desc, have)
		} else if val != want {
			T.Fatalf("%s: wanted %t, got %t", desc, want, val)
		}
	case int:
		if val, is := have.(float64); !is {	// silly Javascript! all numbers are floats
			T.Fatalf("%s: expected integer, got %v", desc, have)
		} else if int(val) != want {
			T.Fatalf("%s: wanted %d, got %d", desc, want, int(val))
		}
	case float64:
		if val, is := have.(float64); !is {
			T.Fatalf("%s: expected float, got %v", desc, have)
		} else if val != want {
			T.Fatalf("%s: wanted %.6f, got %.6f", desc, want, val)
		}
	case string:
		if val, is := have.(string); !is {
			T.Fatalf("%s: expected string, got %v", desc, have)
		} else if val != want {
			T.Fatalf("%s: wanted '%s', got '%s'", desc, want, val)
		}
	case []any:
		if val, is := have.([]any); !is {
			T.Fatalf("%s: expected array, got %v", desc, have)
		} else if len(val) != len(want) {
			T.Fatalf("%s: expected %d-element array; got length %d", desc, len(want),
				len(val))
		} else {
			for i, v := range want {
				checkAnyValue(T, val[i], fmt.Sprintf("%s[%d]", desc, i), v)
			}
		}
	case map[string]any:
		if val, is := have.(map[string]any); !is {
			T.Fatalf("%s: expected map, got %v", desc, have)
		} else if len(val) != len(want) {
			T.Fatalf("%s: expected %d-element map; got length %d", desc, len(want),
				len(val))
		} else {
			for k, v := range want {
				if haveval, exists := val[k]; !exists {
					T.Fatalf("%s: expected to find member %s", desc, k)
				} else {
					checkAnyValue(T, haveval, fmt.Sprintf("%s.%s", desc, k), v)
				}
			}
		}
	default:
		panic("unknown type to check")
	}
}


func Test_generateBasic(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hill theRoad)
		)
	)
	(feature hill
		(popup "Trail along")
		(path likely
		      29.50 -83.43  29.50 -83.41
		      29.40 -83.43  29.40 -83.41)
	)
	(route theRoad
		(segment
			(paths path1 path2)
		)
	)
	(path path1
		30.350075 -83.507595
		30.350177 -83.507918
		30.351014 -83.513659
		30.351541 -83.517636
	)
	(path path2
		(attestation modern_name maybe)
		30.351541 -83.517636
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
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
	vd := prepareAndParseStrings(T, sourceText)


	generated, err := vd.generateJson()
	if err != nil {
		T.Fatal(err.Error())
	}
	checkGeneratedJson(T, generated,
		[]any{
			0,
			map[string]any{"opacity": 0.4},
		},
		[]any{
			0,
			"Trail along",
		},
		[]any{
			map[string]any{
				"menuitem": "Look",
				"f": []any{0, 2},
			},
		},
		[]any{
			map[string]any{
				"t": "feature",
				"popup": 1,
				"f": []any{1},
			},
			map[string]any{
				"t": "path",
				"loc": []any{0,8},
			},
			map[string]any{
				"t": "route",
				"f": []any{3},
			},
			map[string]any{
				"t": "segment",
				"f": []any{4,5},
			},
			map[string]any{
				"t": "path",
				"loc": []any{8,8},
			},
			map[string]any{
				"t": "path",
				"style": 1,
				"loc": []any{16,10},
			},
		},
		[]any{29.500000,-83.430000,29.500000,-83.410000,29.400000,-83.430000,
			29.400000,-83.410000,30.350075,-83.507595,30.350177,-83.507918,
			30.351014,-83.513659,30.351541,-83.517636,30.351541,-83.517636,
			30.351709,-83.519064,30.351815,-83.519952,30.351830,-83.520140,
			30.351842,-83.520299})
}



func Test_generateConcurrentSegment(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features mainRoad)
		)
		(layer two
			(menuitem "Over here")
			(features sideRoad)
		)
	)
	(route mainRoad
		(popup "The Main Way")
		(style roadStyle)
		(segment
			(attestation maybe)
			(paths path1 path2)
		)
		(segment
			(attestation book forSure)
			(path
				30.351842 -83.520299
				30.352932 -83.610901
				30.364123 -83.619623
			)
			(paths path3 mark1)
		)
		(segment
			(attestation maybe)
			(path
				30.372420 -83.632399
				30.375218 -83.629924
				30.374321 -82.623921
			)
		)
	)
	(path path1
		30.350075 -83.507595
		30.350177 -83.507918
		30.351014 -83.513659
		30.351541 -83.517636
	)
	(path path2
		(attestation modern_name maybe)
		30.351541 -83.517636
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
	)
	(path path3
		30.364123 -83.619623
		30.369239 -82.629353
		30.372420 -83.632399
	)
	(marker mark1
		(html ". here it is")
		30.372420 -83.632399
	)

	(route sideRoad
		(style baseStyle)
		(segment
			(attestation maybe)
			(circle
				(popup "side site")
				30.382195 -83.629487
				(radius 300)
			)
			(path
				30.382195 -83.629487
				30.378231 -83.692392
				30.369239 -82.629353
			)
		)
		(segments mainRoad)
		(segment
			(attestation book)
			(path
				30.351842 -83.520299
				30.342397 -83.509359
			)
		)
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
	vd := prepareAndParseStrings(T, sourceText)

	generated, err := vd.generateJson()
	if err != nil {
		T.Fatal(err.Error())
	}
	checkGeneratedJson(T, generated,
		[]any{
			0,
			map[string]any{			//1: baseStyle
				"color": "#1f78b4",
				"fill": true,
				"fillColor": "#1f78b4",
				"fillOpacity": 0.1,
				"opacity": 0.9,
			},
			map[string]any{			//2: roadStyle
				"color": "#AA3333",
				"opacity": 0.6,
				"width": 5,
			},
			map[string]any{			//3: maybe
				"opacity": 0.4,
			},
			map[string]any{			//4: book
				"opacity": 0.8,
				"width": 3,
			},
		},
		[]any{
			0,
			"The Main Way",
			"side site",
		},
		[]any{
			map[string]any{
				"menuitem": "Look",
				"f": []any{0},
			},
			map[string]any{
				"menuitem": "Over here",
				"f": []any{10},
			},
		},
		[]any{
			map[string]any{		//0: mainRoad
				"t": "route",
				"popup": 1,
				"style": 2,
				"f": []any{1, 4, 8},
			},
			map[string]any{
				"t": "segment",
				"style": 3,
				"f": []any{2, 3},
			},
			map[string]any{		//2: path1
				"t": "path",
				"loc": []any{0,8},
			},
			map[string]any{		//3: path2
				"t": "path",
				"style": 3,
				"loc": []any{8,10},
			},
			map[string]any{		//4: (anon segment)
				"t": "segment",
				"style": 4,
				"f": []any{5, 6, 7},
			},
			map[string]any{		//5: (anon path)
				"t": "path",
				"loc": []any{18, 6},
			},
			map[string]any{		//6: path3
				"t": "path",
				"loc": []any{24, 6},
			},
			map[string]any{		//7: mark1
				"t": "marker",
				"html": ". here it is",
				"loc": []any{30, 2},
			},
			map[string]any{		//8: (anon segment)
				"t": "segment",
				"style": 3,
				"f": []any{9},
			},
			map[string]any{		//9: (anon path)
				"t": "path",
				"loc": []any{32, 6},
			},
			map[string]any{		//10: sideRoad
				"t": "route",
				"style": 1,
				"f": []any{11, 14, 16},
			},
			map[string]any{		//11: (anon segment)
				"t": "segment",
				"style": 3,
				"f": []any{12, 13},
			},
			map[string]any{		//12: (anon circle)
				"t": "circle",
				"popup": 2,
				"radius": 300,
				"asPixels": false,
				"loc": []any{38, 2},
			},
			map[string]any{		//13: (anon path)
				"t": "path",
				"loc": []any{40, 6},
			},
			map[string]any{		//14: extract from anon segment #4
				"t": "segment",
				"style": 4,
				"f": []any{15, 5},
			},
			map[string]any{		//15: extract from path3
				"t": "path",
				"loc": []any{24, 4},	// use only the first two points
			},
			map[string]any{		//16: (anon segment)
				"t": "segment",
				"style": 4,
				"f": []any{17},
			},
			map[string]any{		//17: (anon segment)
				"t": "path",
				"loc": []any{46, 4},
			},
		},
		[]any{
			//path1
			30.350075,-83.507595,30.350177,-83.507918,30.351014,-83.513659,
			30.351541,-83.517636,

			//path2
			30.351541,-83.517636,30.351709,-83.519064,30.351815,-83.519952,
			30.351830,-83.520140,30.351842,-83.520299,

			//anon path under mainRoad second segment
			30.351842,-83.520299,30.352932,-83.610901,30.364123,-83.619623,

			//path3
			30.364123,-83.619623,30.369239,-82.629353,30.372420,-83.632399,

			//mark1
			30.372420,-83.632399,

			//anon path under mainRoad third segment
			30.372420,-83.632399,30.375218,-83.629924,30.374321,-82.623921,

			//anon circle under sideRoad first segment
			30.382195,-83.629487,

			//anon path under sideRoad first segment
			30.382195,-83.629487,30.378231,-83.692392,30.369239,-82.629353,

			//anon path under sideRoad third segment
			30.351842,-83.520299,30.342397,-83.509359,
		})
}
