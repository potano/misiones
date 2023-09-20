// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"io"
	"fmt"
	"strings"
	"testing"

	"potano.misiones/sexp"
)


func prepareReader(T *testing.T) (*VectorData, *VectorDataReader) {
	T.Helper()
	vd := NewVectorData()
	vdReader, err := NewVectorDataReader(vd)
	if err != nil {
		T.Fatalf("Grammar error: %s", err)
	}
	return vd, vdReader
}

func prepareAndParse(T *testing.T, streams []io.Reader) *VectorData {
	T.Helper()
	vd, err := basePrepareAndParse(T, streams, true)
	if err != nil {
		T.Fatal(err.Error())
	}
	return vd
}

func prepareAndParseExpectingError(T *testing.T, streams []io.Reader, errmsg string) {
	T.Helper()
	_, err := basePrepareAndParse(T, streams, false)
	if err == nil {
		T.Fatalf("expected error %s", errmsg)
	}
	if err.Error() != errmsg {
		T.Fatalf("expected error '%s', got '%s'", errmsg, err)
	}
}

func basePrepareAndParse(T *testing.T, streams []io.Reader, notePhase bool) (*VectorData, error) {
	T.Helper()
	vd, vdReader := prepareReader(T)
	for i, stream := range streams {
		sourceList, err := sexp.Parse(fmt.Sprintf("infile%d", i), stream)
		if err != nil {
			if notePhase {
				return vd, fmt.Errorf("error in sexp.Parse: %s", err)
			}
			return vd, err
		}
		err = vdReader.ConsumeList(sourceList)
		if err != nil {
			if notePhase {
				return vd, fmt.Errorf("error in VectorDataReader: %s", err)
			}
			return vd, err
		}
	}
	err := vd.ResolveReferences()
	if err != nil {
		if notePhase {
			return vd, fmt.Errorf("error resolving references: %s", err)
		}
		return vd, err
	}
	return vd, nil
}

func checkParse(T *testing.T, vd *VectorData, blob string) {
	got := vd.DescribeNodes("  ")
	if got != blob {
		wantLines := strings.Split(blob, "\n")
		gotLines := strings.Split(got, "\n")
		for i, l := range wantLines {
			g := gotLines[i]
			if g != l {
				T.Fatalf("output line %d: wanted '%s', got '%s'", i, l, g)
			}
		}
		T.Fatalf("expecting %s\ngot %s", blob, got)
	}
}



func Test_basic(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop
      →feature 'hilltop' @ infile0:7
        →marker '$3' @ infile0:8
            location: 29.450000  -83.420000`)
}


func Test_missingReference(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop valley)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"infile0:4: name 'valley' is not registered")
}


func Test_basicMultipleTargets(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop valley)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop valley
      →feature 'hilltop' @ infile0:7
        →marker '$3' @ infile0:8
            location: 29.450000  -83.420000
      →feature 'valley' @ infile0:10
        →point '$5' @ infile0:11
            location: 29.950000  -83.523000`)
}


func Test_multipleTargetsOneMissing(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop vale)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"infile0:4: name 'vale' is not registered")
}


func Test_oneReferrerMultipleReferencesSameTarget(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop valley hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"infile0:4: multiple references to target node 'hilltop'")
}


func Test_multipleTargetsSameName(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	(feature hilltop
		(point 29.95 -83.523)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"infile0:10: duplicate use of name 'hilltop'")
}


func Test_multipleReferrersSameTarget(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
		(layer two
			(menuitem "There")
			(features hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop
      →feature 'hilltop' @ infile0:11
        →marker '$4' @ infile0:12
            location: 29.450000  -83.420000
  →layer 'two' @ infile0:6
      menuitem: 'There'
    →features '' @ infile0:8
        parent: two
        target names: hilltop
      →feature 'hilltop' @ infile0:11
        →marker '$4' @ infile0:12
            location: 29.450000  -83.420000`)
}


func Test_orphanTarget(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"infile0:10: feature 'valley' is an orphan")
}


func Test_targetIsReferrer(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
		(features valley)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop
      →feature 'hilltop' @ infile0:7
        →marker '$3' @ infile0:8
            location: 29.450000  -83.420000
        →features '' @ infile0:9
            parent: hilltop
            target names: valley
          →feature 'valley' @ infile0:11
            →point '$5' @ infile0:12
                location: 29.950000  -83.523000`)
}


func Test_childTargetMultipleReferrers(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop valley)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
		(features valley)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop valley
      →feature 'hilltop' @ infile0:7
        →marker '$3' @ infile0:8
            location: 29.450000  -83.420000
        →features '' @ infile0:9
            parent: hilltop
            target names: valley
          →feature 'valley' @ infile0:11
            →point '$5' @ infile0:12
                location: 29.950000  -83.523000
      →feature 'valley' @ infile0:11
        →point '$5' @ infile0:12
            location: 29.950000  -83.523000`)
}


func Test_cycleToSelf(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop valley)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
		(features valley hilltop)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"cycle detected in dependency graph; check references to 'hilltop'")
}


func Test_cycleToNodeHigherInChain(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
		(features valley)
	)
	(feature valley
		(point 29.95 -83.523)
		(features hilltop)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"cycle detected in dependency graph; check references to 'hilltop'")
}


func Test_cycleLowInChain(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features MoundPark)
		)
		(layer two
			(menuitem "There")
			(features MoundPkHill)
		)
	)
	(feature MoundPark
		(popup "This area")
		(features MoundPkHill)
	)
	(feature MoundPkHill
		(marker 29.45 -83.42)
		(features MoundPath dale)
	)
	(feature dale
		(point 29.95 -83.523)
		(features MoundPath DaleHead)
	)
	(feature DaleHead
		(features dale)
	)
	(feature MoundPath
		(path 30.1 -83.2  30.2 -83.1)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"cycle detected in dependency graph; check references to 'MoundPkHill' and 'dale'")
}



func Test_featureComponents(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop mound)
		)
	)
	(feature hilltop
		(popup "Here I am")
		(marker 29.45 -83.42)
		(style basicStyle)
		(attestation one two)
		(rectangle 29.50 -83.43  29.50 -83.41
		           29.40 -83.43  29.40 -83.41)
	)
	(feature mound
		(circle 29.95 -83.523 (radius 20))
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop mound
      →feature 'hilltop' @ infile0:7
          popup text: 'Here I am'
          style: basicStyle
          attestation: one two
        →marker '$3' @ infile0:9
            location: 29.450000  -83.420000
        →rectangle '$4' @ infile0:12
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
      →feature 'mound' @ infile0:15
        →circle '$6' @ infile0:16
            radius: 20 meters
            location: 29.950000  -83.523000`)
}


func Test_featureOptionalFeatureName(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
	)
	(feature hilltop
		(popup "Here I am")
		(marker 29.45 -83.42)
		(style basicStyle)
		(attestation one two)
		(rectangle 29.50 -83.43  29.50 -83.41
		           29.40 -83.43  29.40 -83.41)
		(feature mound
			(circle 29.95 -83.523 (radius 20))
		)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop
      →feature 'hilltop' @ infile0:7
          popup text: 'Here I am'
          style: basicStyle
          attestation: one two
        →marker '$3' @ infile0:9
            location: 29.450000  -83.420000
        →rectangle '$4' @ infile0:12
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
        →feature 'mound' @ infile0:14
          →circle '$6' @ infile0:15
              radius: 20 meters
              location: 29.950000  -83.523000`)
}


func Test_featureWithReferencedInternalFeature(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop mound)
		)
	)
	(feature hilltop
		(popup "Here I am")
		(marker 29.45 -83.42)
		(style basicStyle)
		(attestation one two)
		(rectangle 29.50 -83.43  29.50 -83.41
		           29.40 -83.43  29.40 -83.41)
		(feature mound
			(circle 29.95 -83.523 (radius 20))
		)
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: hilltop mound
      →feature 'hilltop' @ infile0:7
          popup text: 'Here I am'
          style: basicStyle
          attestation: one two
        →marker '$3' @ infile0:9
            location: 29.450000  -83.420000
        →rectangle '$4' @ infile0:12
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
        →feature 'mound' @ infile0:14
          →circle '$6' @ infile0:15
              radius: 20 meters
              location: 29.950000  -83.523000
      →feature 'mound' @ infile0:14
        →circle '$6' @ infile0:15
            radius: 20 meters
            location: 29.950000  -83.523000`)
}


func Test_featureWithPath(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features tinytrail mound)
		)
	)
	(feature tinytrail
		(popup "Trail along")
		(path 29.50 -83.43  29.50 -83.41
		      29.40 -83.43  29.40 -83.41)
	)
	(feature mound
		(circle 29.95 -83.523 (radius 20))
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: tinytrail mound
      →feature 'tinytrail' @ infile0:7
          popup text: 'Trail along'
        →path '$3' @ infile0:9
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
      →feature 'mound' @ infile0:12
        →circle '$5' @ infile0:13
            radius: 20 meters
            location: 29.950000  -83.523000`)
}


func Test_featureWithNamedEmbeddedPath(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features tinytrail mound)
		)
	)
	(feature tinytrail
		(popup "Trail along")
		(path likely
		      29.50 -83.43  29.50 -83.41
		      29.40 -83.43  29.40 -83.41)
	)
	(feature mound
		(circle 29.95 -83.523 (radius 20))
	)
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

	checkParse(T, vd,
`→layers '$0' @ infile0:1
  →layer 'one' @ infile0:2
      menuitem: 'Look'
    →features '' @ infile0:4
        parent: one
        target names: tinytrail mound
      →feature 'tinytrail' @ infile0:7
          popup text: 'Trail along'
        →path 'likely' @ infile0:9
            location: 29.500000  -83.430000
                      29.500000  -83.410000
                      29.400000  -83.430000
                      29.400000  -83.410000
      →feature 'mound' @ infile0:13
        →circle '$5' @ infile0:14
            radius: 20 meters
            location: 29.950000  -83.523000`)
}


func Test_illegalChildFeature(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features hilltop)
		)
		(layer two
			(menuitem "There")
			(features valley)
		)
	)
	(feature hilltop
		(marker 29.45 -83.42)
		(features valley two)
	)
	(feature valley
		(point 29.95 -83.523)
	)
	`
	prepareAndParseExpectingError(T, []io.Reader{strings.NewReader(sourceText)},
		"infile0:13: referenced item 'two' is a layer type; only feature, marker, point, path, polygon, rectangle, circle, and route allowed")
}




func Test_routeSegmentPaths(T *testing.T) {
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
			(attestation maybe)
			(paths path1 path2)
		)
		(segment
			(attestation forSure)
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
	`
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})

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
            attestation: maybe
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
            attestation: forSure
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
}

