// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Tests paralleling route tests in m_test.go but with routes already reformed


func Test_measureThreadedFirstPathSpur(T *testing.T) {
	//No global route threading applied.  See similar test in mt_test.go
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred)
		)
	)
	(route theRoad
		(segment mainSeg1
			(paths path1 path2)
		)
		(segment mainSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route spurred
		(segment spurSeg1
			(paths path1Spur path1SpurWaypoint path1 path2)
		)
		(segments mainSeg2)
	)
	`
	source2 := path1 + path2 + path3 + path4 + path5 + path6 + path1Spur + path1SpurWaypoint
	vd := prepareAndParseStrings(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path1", path1_length},
		{"path2", path2_length},
		{"path1Spur", path1Spur_length},
		{"spurred", path1Spur_length + path1DownFromSpurLength + path2_length +
			path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{5, 30.351056, -83.513662, 5.2, "path1Spur", 2},
		{9, 30.351014, -83.513659, 9.9, "path1Spur", 0},
		{10, 30.351014, -83.513659, 9.9, "path1:1", 0},	//threading shortened path 1
		{300, 30.351541, -83.517636, 396.1, "path1:1", 1},
		{400, 30.351541, -83.517636, 396.1, "path2", 0},
		{500, 30.351709, -83.519064, 534.4, "path2", 1},
		{650, 30.351842, -83.520299, 653.9, "path2", 4},
		{670, 30.351850, -83.520426, 666.1, "path3", 1},
		{700, 30.351879, -83.520762, 698.5, "path4", 0},
		{800, 30.351801, -83.521739, 792.8, "path5", 0},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("spurred", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureThreadedFirstPathSpurSpurPathReversed(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred)
		)
	)
	(route theRoad
		(segment mainSeg1
			(paths path1 path2)
		)
		(segment mainSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route spurred
		(segment spurSeg1
			(paths path1SpurReversed path1SpurWaypoint path1 path2)
		)
		(segments mainSeg2)
	)
	`
	source2 := path1 + path2 + path3 + path4 + path5 + path6 + path1SpurReversed +
		path1SpurWaypoint
	vd := prepareAndParseStrings(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path1", path1_length},
		{"path2", path2_length},
		{"path1SpurReversed", path1Spur_length},
		{"spurred", path1Spur_length + path1DownFromSpurLength + path2_length +
			path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{5, 30.351056, -83.513662, 5.2, "path1SpurReversed", 1},
		{9, 30.351014, -83.513659, 9.9, "path1SpurReversed", 3},
		{10, 30.351014, -83.513659, 9.9, "path1:1", 0},
		{300, 30.351541, -83.517636, 396.1, "path1:1", 1},
		{400, 30.351541, -83.517636, 396.1, "path2", 0},
		{500, 30.351709, -83.519064, 534.4, "path2", 1},
		{650, 30.351842, -83.520299, 653.9, "path2", 4},
		{670, 30.351850, -83.520426, 666.1, "path3", 1},
		{700, 30.351879, -83.520762, 698.5, "path4", 0},
		{800, 30.351801, -83.521739, 792.8, "path5", 0},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("spurred", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureThreadedFirstPathSpurFirstPathReversed(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred)
		)
	)
	(route theRoad
		(segment mainSeg1
			(paths path1Reversed path2)
		)
		(segment mainSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route spurred
		(segment spurSeg1
			(paths path1Spur path1SpurWaypoint path1Reversed path2)
		)
		(segments mainSeg2)
	)
	`
	source2 := path1Reversed + path2 + path3 + path4 + path5 + path6 + path1Spur +
		path1SpurWaypoint
	vd := prepareAndParseStrings(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path1Reversed", path1_length},
		{"path2", path2_length},
		{"path1Spur", path1Spur_length},
		{"spurred", path1Spur_length + path1DownFromSpurLength + path2_length +
			path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{5, 30.351056, -83.513662, 5.2, "path1Spur", 2},
		{9, 30.351014, -83.513659, 9.9, "path1Spur", 0},
		{10, 30.351014, -83.513659, 9.9, "path1Reversed:1", 1},
		{300, 30.351541, -83.517636, 396.1, "path1Reversed:1", 0},
		{400, 30.351541, -83.517636, 396.1, "path2", 0},
		{500, 30.351709, -83.519064, 534.4, "path2", 1},
		{650, 30.351842, -83.520299, 653.9, "path2", 4},
		{670, 30.351850, -83.520426, 666.1, "path3", 1},
		{700, 30.351879, -83.520762, 698.5, "path4", 0},
		{800, 30.351801, -83.521739, 792.8, "path5", 0},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("spurred", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureThreadedSecondPathSpur(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred fragments)
		)
	)
	(route theRoad
		(segment mainSeg1
			(paths path1 path2)
		)
		(segment mainSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route spurred
		(segment spurSeg1
			(path path2Spur
				30.351815 -83.519952
				30.351622 -83.519899
				30.351426 -83.519868
				30.351401 -83.519822
			)
			(point spurWaypoint 30.351815 -83.519952)
			(paths path2)
			(point wpPath2End 30.351842 -83.520299)
		)
		(segments mainSeg2)
	)
	(feature fragments
		(segment path2_DownFromSpur
			(paths spurWaypoint path2 wpPath2End)
		)
	)
	`
	source2 := path1 + path2 + path3 + path4 + path5 + path6
	const path2SpurLength = 49.281260
	const path2DownFromSpurLength = 33.440643
	vd := prepareAndParseStrings(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path2Spur", path2SpurLength},
		{"path2_DownFromSpur", path2DownFromSpurLength},
		{"spurred", path2SpurLength + path2DownFromSpurLength +
			path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{5, 30.351426, -83.519868, 5.2, "path2Spur", 2},
		{30, 30.351622, -83.519899, 27.2, "path2Spur", 1},
		{50, 30.351815, -83.519952, 49.3, "path2:1", 0},
		{70, 30.351830, -83.520140, 67.4, "path2:1", 1},
		{90, 30.351850, -83.520426, 94.9, "path3", 1},
		{300, 30.351707, -83.522451, 290.7, "path5", 4},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("spurred", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureThreadedSecondPathSpurSecondPathReversed(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred fragments)
		)
	)
	(route theRoad
		(segment mainSeg1
			(paths path1 path2Reversed)
		)
		(segment mainSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route spurred
		(segment spurSeg1
			(path path2Spur
				30.351815 -83.519952
				30.351622 -83.519899
				30.351426 -83.519868
				30.351401 -83.519822
			)
			(point spurWaypoint 30.351815 -83.519952)
			(paths path2Reversed)
			(point wpPath2End 30.351842 -83.520299)
		)
		(segments mainSeg2)
	)
	(feature fragments
		(segment path2_DownFromSpur
			(paths spurWaypoint path2Reversed wpPath2End)
		)
	)
	`
	source2 := path1 + path2Reversed + path3 + path4 + path5 + path6
	const path2SpurLength = 49.281260
	const path2DownFromSpurLength = 33.440643
	vd := prepareAndParseStrings(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path2Spur", path2SpurLength},
		{"path2_DownFromSpur", path2DownFromSpurLength},
		{"spurred", path2SpurLength + path2DownFromSpurLength +
			path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{5, 30.351426, -83.519868, 5.2, "path2Spur", 2},
		{30, 30.351622, -83.519899, 27.2, "path2Spur", 1},
		{50, 30.351815, -83.519952, 49.3, "path2Reversed:1", 2},
		{70, 30.351830, -83.520140, 67.4, "path2Reversed:1", 1},
		{90, 30.351850, -83.520426, 94.9, "path3", 1},
		{300, 30.351707, -83.522451, 290.7, "path5", 4},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("spurred", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureThreadedSecondPathSpurFirstMainSegmentReversed(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred fragments)
		)
	)
	(route theRoad
		(segment mainSeg1
			(paths path2 path1)
		)
		(segment mainSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(route spurred
		(segment spurSeg1
			(path path2Spur
				30.351815 -83.519952
				30.351622 -83.519899
				30.351426 -83.519868
				30.351401 -83.519822
			)
			(point spurWaypoint 30.351815 -83.519952)
			(paths path2)
			(point wpPath2End 30.351842 -83.520299)
		)
		(segments mainSeg2)
	)
	(feature fragments
		(segment path2_DownFromSpur
			(paths spurWaypoint path2 wpPath2End)
		)
	)
	`
	source2 := path1 + path2 + path3 + path4 + path5 + path6
	const path2SpurLength = 49.281260
	const path2DownFromSpurLength = 33.440643
	vd := prepareAndParseStrings(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path2Spur", path2SpurLength},
		{"path2_DownFromSpur", path2DownFromSpurLength},
		{"spurred", path2SpurLength + path2DownFromSpurLength +
			path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{5, 30.351426, -83.519868, 5.2, "path2Spur", 2},
		{30, 30.351622, -83.519899, 27.2, "path2Spur", 1},
		{50, 30.351815, -83.519952, 49.3, "path2:1", 0},
		{70, 30.351830, -83.520140, 67.4, "path2:1", 1},
		{90, 30.351850, -83.520426, 94.9, "path3", 1},
		{300, 30.351707, -83.522451, 290.7, "path5", 4},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("spurred", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}

