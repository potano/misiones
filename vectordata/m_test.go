// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"math"
	"testing"
)

const pointStartPath1 = "(point pointStartPath1 30.350075 -83.507595)"
const path1 = `
	(path path1
		30.350075 -83.507595
		30.350177 -83.507918
		30.351014 -83.513659
		30.351541 -83.517636
	)
	`
const path1_length = 978.04902
const path1Reversed = `
	(path path1Reversed
		30.351541 -83.517636
		30.351014 -83.513659
		30.350177 -83.507918
		30.350075 -83.507595
	)
	`
const path1Spur = `
	(path path1Spur
		30.351014 -83.513659
		30.351035 -83.513658
		30.351056 -83.513662
		30.351102 -83.513671
	)
	`
const path1Spur_length = 9.893595
const path1DownFromSpur = `
	(path path1DownFromSpur
		30.351014 -83.513659
		30.351541 -83.517636
	)
	`
const path1DownFromSpurLength = 386.193026
const path1SpurWaypoint = "(point path1SpurWaypoint  30.351014 -83.513659)"
const path1SpurReversed = `
	(path path1SpurReversed
		30.351102 -83.513671
		30.351056 -83.513662
		30.351035 -83.513658
		30.351014 -83.513659
	)
	`
const path2 = `
	(path path2
		30.351541 -83.517636
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
	)
	`
const path2_length = 257.81256
const path2Reversed = `
	(path path2Reversed
		30.351842 -83.520299
		30.351830 -83.520140
		30.351815 -83.519952
		30.351709 -83.519064
		30.351541 -83.517636
	)
	`
const path2Disconnected = `
	(path path2Disconnected
		30.351709 -83.519064
		30.351815 -83.519952
		30.351830 -83.520140
		30.351842 -83.520299
	)
	`
const path3 = `
	(path path3
		30.351842 -83.520299
		30.351850 -83.520426
		30.351861 -83.520554
		30.351870 -83.520668
		30.351879 -83.520762
	)
	`
const path3_length = 44.63347
const path3Reversed = `
	(path path3Reversed
		30.351879 -83.520762
		30.351870 -83.520668
		30.351861 -83.520554
		30.351850 -83.520426
		30.351842 -83.520299
	)
	`
const path4 = `
	(path path4
		30.351879 -83.520762
		30.351881 -83.520799
		30.351879 -83.520838
		30.351872 -83.520932
		30.351848 -83.521227
		30.351827 -83.521477
		30.351801 -83.521739
	)
	`
const path4_length = 94.22527
const path5 = `
	(path path5
		30.351801 -83.521739
		30.351774 -83.521975
		30.351746 -83.522218
		30.351730 -83.522325
		30.351707 -83.522451
		30.351668 -83.522666
		30.351639 -83.522823
	)
	`
const path5_length = 105.66177
const path6 = `
	(path path6
		30.351639 -83.522823
		30.351614 -83.522970
		30.351570 -83.523202
		30.351534 -83.523411
		30.351486 -83.523684
		30.351440 -83.523922
		30.351395 -83.524186
		30.351354 -83.524397
	)
	`
const path6_length = 154.37424
const path6Reversed = `
	(path path6Reversed
		30.351354 -83.524397
		30.351395 -83.524186
		30.351440 -83.523922
		30.351486 -83.523684
		30.351534 -83.523411
		30.351570 -83.523202
		30.351614 -83.522970
		30.351639 -83.522823
	)
	`



func compareTestLengths(T *testing.T, name string, want, got float64) {
	T.Helper()
	diff := want - got
	if diff < -0.05 || diff > 0.05 {
		T.Fatalf("%s: wanted length %.1f, got %.1f", name, want, got)
	}
}

func compareTestUpTo(T *testing.T,
		wantLat, wantLong, wantDistance float64, wantName string, wantIndex int,
		gotLat, gotLong, gotDistance float64, gotName string, gotIndex int) {
	diffLat := math.Abs(wantLat - gotLat)
	diffLong := math.Abs(wantLong - gotLong)
	diffDistance := math.Abs(wantDistance - gotDistance)
	if diffLat > 5E-7 || diffLong > 5E-7 || diffDistance > 5E-2 || gotName != wantName ||
			gotIndex != wantIndex {
		T.Fatalf("wanted %.1f @ [%.6f %.6f] in %s #%d\ngot %.1f @ [%.6f %.6f] in %s #%d",
			wantDistance, wantLat, wantLong, wantName, wantIndex,
			gotDistance, gotLat, gotLong, gotName, gotIndex)
	}
}



func Test_measurePaths(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment
			(paths path1 path2)
		)
		(segment
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"path1", path1_length},
		{"path2", path2_length},
		{"path3", path3_length},
		{"path4", path4_length},
		{"path5", path5_length},
		{"path6", path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureReversePaths(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment
			(paths path1Reversed path2Reversed path3Reversed)
		)
		(segment
			(paths path4 path5 path6Reversed)
		)
	)
	` + path1Reversed + path2Reversed + path3Reversed + path4 + path5 + path6Reversed
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"path1Reversed", path1_length},
		{"path2Reversed", path2_length},
		{"path3Reversed", path3_length},
		{"path6Reversed", path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureSegments(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"roadSeg1", path1_length + path2_length},
		{"roadSeg2", path3_length + path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureSegmentsWithPathReversals(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2Reversed path3)
		)
		(segment roadSeg2
			(paths path4 path5 path6Reversed)
		)
	)
	` + path1 + path2Reversed + path3 + path4 + path5 + path6Reversed
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"roadSeg1", path1_length + path2_length + path3_length},
		{"roadSeg2", path4_length + path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRoute(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithPathReversals(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2Reversed path3)
		)
		(segment roadSeg2
			(paths path4 path5 path6Reversed)
		)
	)
	` + path1 + path2Reversed + path3 + path4 + path5 + path6Reversed
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithFirstSegmentReversal(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path3 path2 path1)
		)
		(segment roadSeg2
			(paths path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithSegment2Reversal(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2 path3)
		)
		(segment roadSeg2
			(paths path6 path5 path4)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithSegmentAndPathReversal(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2 path3Reversed)
		)
		(segment roadSeg2
			(paths path6 path5 path4)
		)
	)
	` + path1 + path2 + path3Reversed + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithStartingPoint(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
	        (marker
			(popup "Start of route")
			30.350075 -83.507595
		)
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithStartingPointByReference(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths pointStartPath1 path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6 + pointStartPath1
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithWaypoint(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
	        (circle
			(radius 5)
			30.351842 -83.520299
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}


func Test_measureRouteWithStartingPointByReferenceInSegment(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths pointStartPath1 path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6 + pointStartPath1
	vd := prepareAndParseStrings(T, sourceText)
	for _, test := range []struct{name string; meters float64} {
		{"theRoad", path1_length + path2_length + path3_length + path4_length +
			path5_length + path6_length},
	} {
		distance, err := vd.MeasurePath(test.name)
		if err != nil {
			T.Fatalf("error measuring %s: %s", test.name, err)
		}
		compareTestLengths(T, test.name, test.meters, distance)
	}
}




func Test_discontinuousSegment(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2Disconnected)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2Disconnected + path3 + path4 + path5 + path6
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	_, err := vd.MeasurePath("roadSeg1")
	want := "infile0:8: path 'path2Disconnected' does not connect with path 'path1' in segment 'roadSeg1'"
	if err == nil || want != err.Error() {
		T.Fatalf("wanted error \"%s\"\ngot \"%s\"", want, err)
	}
}


func Test_segmentWithDisconnectedStartingPoint(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(point errant 30.350074 -83.507591)
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	_, err := vd.MeasurePath("roadSeg1")
	want := "infile0:8: path 'path1' does not connect with point 'errant' in segment 'roadSeg1'"
	if err == nil || want != err.Error() {
		T.Fatalf("wanted error \"%s\"\ngot \"%s\"", want, err)
	}
}


func Test_discontinuousRoute(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path5 + path6
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	_, err := vd.MeasurePath("theRoad")
	want := "infile0:11: path 'path5' does not connect with path 'path3' in segment 'roadSeg2'"
	if err == nil || want != err.Error() {
		T.Fatalf("wanted error \"%s\"\ngot \"%s\"", want, err)
	}
}


func Test_discontinuousRouteUsingThreading(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path5 + path6
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	err := vd.CheckAndReformRoutes()
	want := "infile0:11: path 'path5' does not connect with path 'path3' in segment 'roadSeg2'"
	if err == nil || want != err.Error() {
		T.Fatalf("wanted error \"%s\"\ngot \"%s\"", want, err)
	}
}




func Test_measurePathUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("path1", 600)
	if err != nil {
		T.Fatal(err.Error())
	}
	compareTestUpTo(T, 30.351014, -83.513659, 591.9, "path1", 2,
		lat, long, distance, pathName, index)
}


func Test_measurePath2UpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
	)
	` + path1 + path2
	vd := prepareAndParseStrings(T, sourceText)
	lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("path2", 100)
	if err != nil {
		T.Fatal(err.Error())
	}
	compareTestUpTo(T, 30.351709, -83.519064, 138.3, "path2", 1,
		lat, long, distance, pathName, index)
}


func Test_measurePath2ReversedUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2Reversed)
		)
	)
	` + path1 + path2Reversed
	vd := prepareAndParseStrings(T, sourceText)
	lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("path2Reversed", 100)
	if err != nil {
		T.Fatal(err.Error())
	}
	compareTestUpTo(T, 30.351709, -83.519064, 119.5, "path2Reversed", 3,
		lat, long, distance, pathName, index)
}


func Test_measureBeyondPath(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
	)
	` + path1 + path2
	vd := prepareAndParseStrings(T, sourceText)
	_, _, _, _, _, err := vd.MeasurePathUpTo("path1", 1000)
	want := "path 'path1' is only 978.0 meters (0.61 miles) long"
	if err == nil || err.Error() != want {
		T.Fatalf("expected error \"%s\"\n got \"%s\"", want, err)
	}
}


func Test_measureSegmentUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
	)
	` + path1 + path2
	vd := prepareAndParseStrings(T, sourceText)
	lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("roadSeg1", 1100)
	if err != nil {
		T.Fatal(err.Error())
	}
	compareTestUpTo(T, 30.351709, -83.519064, 1116.4, "path2", 1,
		lat, long, distance, pathName, index)
}


func Test_measureSegmentUpToReversedPath(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2Reversed)
		)
	)
	` + path1 + path2Reversed
	vd := prepareAndParseStrings(T, sourceText)
	lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("roadSeg1", 1100)
	if err != nil {
		T.Fatal(err.Error())
	}
	compareTestUpTo(T, 30.351709, -83.519064, 1116.4, "path2Reversed", 3,
		lat, long, distance, pathName, index)
}


func Test_measureBeyondSegment(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
	)
	` + path1 + path2
	vd := prepareAndParseStrings(T, sourceText)
	_, _, _, _, _, err := vd.MeasurePathUpTo("roadSeg1", 2000)
	want := "segment 'roadSeg1' is only 1235.9 meters (0.77 miles) long"
	if err == nil || err.Error() != want {
		T.Fatalf("expected error \"%s\"\n got \"%s\"", want, err)
	}
}


func Test_measureRouteUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{700, 30.351014, -83.513659, 591.9, "path1", 2},
		{1100, 30.351709, -83.519064, 1116.4, "path2", 1},
		{1255, 30.351861, -83.520554, 1260.4, "path3", 2},
		{1300, 30.351872, -83.520932, 1296.9, "path4", 3},
		{1450, 30.351707, -83.522451, 1443.9, "path5", 4},
		{1500, 30.351614, -83.522970, 1494.8, "path6", 1},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureRouteUpToReversedPathSeg0(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2Reversed)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2Reversed + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", 1100)
	if err != nil {
		T.Fatal(err.Error())
	}
	compareTestUpTo(T, 30.351709, -83.519064, 1116.4, "path2Reversed", 3,
		lat, long, distance, pathName, index)
}


func Test_measureRouteUpToReversedSeg1(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path6 path5 path4 path3)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{700, 30.351014, -83.513659, 591.9, "path1", 2},
		{1100, 30.351709, -83.519064, 1116.4, "path2", 1},
		{1255, 30.351861, -83.520554, 1260.4, "path3", 2},
		{1300, 30.351872, -83.520932, 1296.9, "path4", 3},
		{1450, 30.351707, -83.522451, 1443.9, "path5", 4},
		{1500, 30.351614, -83.522970, 1494.8, "path6", 1},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureRouteUpToReversedSeg1ReversedPath3(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path6 path5 path4 path3Reversed)
		)
	)
	` + path1 + path2 + path3Reversed + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{700, 30.351014, -83.513659, 591.9, "path1", 2},
		{1100, 30.351709, -83.519064, 1116.4, "path2", 1},
		{1255, 30.351861, -83.520554, 1260.4, "path3Reversed", 2},
		{1300, 30.351872, -83.520932, 1296.9, "path4", 3},
		{1450, 30.351707, -83.522451, 1443.9, "path5", 4},
		{1500, 30.351614, -83.522970, 1494.8, "path6", 1},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureRouteReversedFirstSegmentUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path2 path1)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{700, 30.351014, -83.513659, 591.9, "path1", 2},
		{1100, 30.351709, -83.519064, 1116.4, "path2", 1},
		{1255, 30.351861, -83.520554, 1260.4, "path3", 2},
		{1300, 30.351872, -83.520932, 1296.9, "path4", 3},
		{1450, 30.351707, -83.522451, 1443.9, "path5", 4},
		{1500, 30.351614, -83.522970, 1494.8, "path6", 1},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureRouteReversedSecondSegmentUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path6 path5 path4 path3)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{700, 30.351014, -83.513659, 591.9, "path1", 2},
		{1100, 30.351709, -83.519064, 1116.4, "path2", 1},
		{1255, 30.351861, -83.520554, 1260.4, "path3", 2},
		{1300, 30.351872, -83.520932, 1296.9, "path4", 3},
		{1450, 30.351707, -83.522451, 1443.9, "path5", 4},
		{1500, 30.351614, -83.522970, 1494.8, "path6", 1},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureRouteReversedBothSegmentsUpTo(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path2 path1)
		)
		(segment roadSeg2
			(paths path6 path5 path4 path3)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	for _, tst := range []struct{dist, lat, long, expect float64; name string; index int} {
		{700, 30.351014, -83.513659, 591.9, "path1", 2},
		{1100, 30.351709, -83.519064, 1116.4, "path2", 1},
		{1255, 30.351861, -83.520554, 1260.4, "path3", 2},
		{1300, 30.351872, -83.520932, 1296.9, "path4", 3},
		{1450, 30.351707, -83.522451, 1443.9, "path5", 4},
		{1500, 30.351614, -83.522970, 1494.8, "path6", 1},
	} {
		lat, long, distance, pathName, index, err := vd.MeasurePathUpTo("theRoad", tst.dist)
		if err != nil {
			T.Fatal(err.Error())
		}
		compareTestUpTo(T, tst.lat, tst.long, tst.expect, tst.name, tst.index,
			lat, long, distance, pathName, index)
	}
}


func Test_measureBeyondRoute(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParseStrings(T, sourceText)
	_, _, _, _, _, err := vd.MeasurePathUpTo("theRoad", 2000)
	want := "route 'theRoad' is only 1634.8 meters (1.02 miles) long"
	if err == nil || err.Error() != want {
		T.Fatalf("expected error \"%s\"\n got \"%s\"", want, err)
	}
}



// In these route-measurement tests no global route reformation is applied.
// See similar tests reflective of reformed paths in mt_test.go

func Test_measureFirstPathSpur(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad spurred testDownFromSpurLength)
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
	(feature testDownFromSpurLength ` + path1DownFromSpur + `)
	`
	source2 := path1 + path2 + path3 + path4 + path5 + path6 + path1Spur + path1SpurWaypoint
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText, source2)
	for _, test := range []struct{name string; meters float64} {
		{"path1", path1_length},
		{"path2", path2_length},
		{"path1Spur", path1Spur_length},
		{"path1DownFromSpur", path1DownFromSpurLength},
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
		{10, 30.351014, -83.513659, 9.9, "path1", 2},
		{300, 30.351541, -83.517636, 396.1, "path1", 3},
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


func Test_measureFirstPathSpurSpurPathReversed(T *testing.T) {
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
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText, source2)
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
		{10, 30.351014, -83.513659, 9.9, "path1", 2},
		{300, 30.351541, -83.517636, 396.1, "path1", 3},
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


func Test_measureFirstPathSpurFirstPathReversed(T *testing.T) {
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
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText, source2)
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
		{10, 30.351014, -83.513659, 9.9, "path1Reversed", 1},
		{300, 30.351541, -83.517636, 396.1, "path1Reversed", 0},
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


func Test_measureSecondPathSpur(T *testing.T) {
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
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText, source2)
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
		{50, 30.351815, -83.519952, 49.3, "path2", 2},
		{70, 30.351830, -83.520140, 67.4, "path2", 3},
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


func Test_measureSecondPathSpurSecondPathReversed(T *testing.T) {
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
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText, source2)
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
		{50, 30.351815, -83.519952, 49.3, "path2Reversed", 2},
		{70, 30.351830, -83.520140, 67.4, "path2Reversed", 1},
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


func Test_measureSecondPathSpurFirstMainSegmentReversed(T *testing.T) {
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
	vd := prepareAndParseStringsNoRouteThreading(T, sourceText, source2)
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
		{50, 30.351815, -83.519952, 49.3, "path2", 2},
		{70, 30.351830, -83.520140, 67.4, "path2", 3},
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

