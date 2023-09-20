// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"io"
	"strings"
	"testing"
)



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
		(attestation modern_name)
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
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})


	generated, err := vd.GenerateJs()
	if err != nil {
		T.Fatal(err.Error())
	}
	want := `var $s0={"opacity":0.8,"width":2}
allVectors=[{"menuitem":"Look","features":[{"t":"feature","popup":"Trail along","features":[{"t":"path","coords":[29.500000,-83.430000,29.500000,-83.410000,29.400000,-83.430000,29.400000,-83.410000]}]},{"t":"route","features":[{"t":"segment","paths":[{"t":"path","coords":[30.350075,-83.507595,30.350177,-83.507918,30.351014,-83.513659,30.351541,-83.517636]},{"t":"path","style":$s1,"coords":[30.351541,-83.517636,30.351709,-83.519064,30.351815,-83.519952,30.351830,-83.520140,30.351842,-83.520299]}]}]}]}]`
	if generated != want {
		T.Fatalf("Failed to generate\n%s\ngot\n%s", want, generated)
	}
}

