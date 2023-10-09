// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"


func Test_measureRoutePortion(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad pois roadPart)
		)
	)
	(feature pois
		(marker wp2 30.351354 -83.524397)   ;end of path6
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route roadPart
		(segment
			(paths path1Spur path1SpurWaypoint)
		)
		(routeSegments theRoad path1SpurWaypoint wp2)
	)
	`
	sourceText2 := path1 + path2 + path3 + path4 + path5 + path6 + path1SpurWaypoint +
		path1Spur
	vd := prepareAndParseStrings(T, sourceText, sourceText2)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
		{"roadPart", path1Spur_length + path1DownFromSpurLength + path2_length + path3_length +
			path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRoutePortionReverse(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad pois roadPart)
		)
	)
	(feature pois
		(marker wp2 30.351354 -83.524397)   ;end of path6
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route roadPart
		(routeSegments theRoad wp2 path1SpurWaypoint)
		(segment
			(paths path1SpurWaypoint path1Spur)
		)
	)
	`
	sourceText2 := path1 + path2 + path3 + path4 + path5 + path6 + path1SpurWaypoint +
		path1Spur
	vd := prepareAndParseStrings(T, sourceText, sourceText2)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
		{"roadPart", path1Spur_length + path1DownFromSpurLength + path2_length + path3_length +
			path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureThreeSegmentRoutePortion(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad pois roadPart)
		)
	)
	(feature pois
		(marker wp2 30.351354 -83.524397)   ;end of path6
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4)
		)
		(segment roadSeg3
			(paths path5 path6)
		)
	)
	(route roadPart
		(segment
			(paths path1Spur path1SpurWaypoint)
		)
		(routeSegments theRoad path1SpurWaypoint wp2)
	)
	`
	sourceText2 := path1 + path2 + path3 + path4 + path5 + path6 + path1SpurWaypoint +
		path1Spur
	vd := prepareAndParseStrings(T, sourceText, sourceText2)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
		{"roadPart", path1Spur_length + path1DownFromSpurLength + path2_length + path3_length +
			path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}

