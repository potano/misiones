// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Path-threading tests


type gsCheck struct {
	lat1, long1, lat2, long2 float64
	paths []gsPath
}

type gsPath struct {
	name string
	forward bool
	points locationPairs
}

func checkGatheredSegmentPaths(T *testing.T, vd *VectorData, segName string, wantSeg gsCheck) {
	T.Helper()
	item, exists := vd.mapItems[segName]
	if !exists {
		T.Fatal("Can't find segment " + segName)
	}
	seg, is := item.(*mapSegmentType)
	if !is {
		T.Fatalf("%s is not a segment", segName)
	}
	gathered, err := seg.threadPaths()
	if err != nil {
		T.Fatalf("error processing %s: %s", segName, err)
	}
	if gathered.obj.Name() != segName {
		T.Fatalf("expected segment name %s, got %s", segName, gathered.obj.Name())
	}
	baseSegmentCheck(T, gathered, wantSeg)
}

func checkGatheredRouteSegments(T *testing.T, vd *VectorData, routeName string, want []gsCheck) {
	T.Helper()
	item, exists := vd.mapItems[routeName]
	if !exists {
		T.Fatal("Can't find route " + routeName)
	}
	route, is := item.(*mapRouteType)
	if !is {
		T.Fatalf("%s is not a route", routeName)
	}
	gSegments, err := route.threadSegments()
	if err != nil {
		T.Fatalf("error processing %s: %s", routeName, err)
	}
	if len(gSegments) != len(want) {
		T.Fatalf("expected %d segments, got %d", len(want), len(gSegments))
	}
	for segX, gs := range gSegments {
		baseSegmentCheck(T, gs, want[segX])
	}
}

func baseSegmentCheck(T *testing.T, gathered *gatheredSegment, wantSeg gsCheck) {
	lat, long := gathered.lat1, gathered.long1
	if !isSamePoint(wantSeg.lat1, wantSeg.long1, lat, long) {
		T.Fatalf("expected segment start [%f  %f], got [%f %f]", wantSeg.lat1,
			wantSeg.long1, lat, long)
	}
	lat, long = gathered.lat2, gathered.long2
	if !isSamePoint(wantSeg.lat2, wantSeg.long2, lat, long) {
		T.Fatalf("expected segment end [%f  %f], got [%f %f]", wantSeg.lat2,
			wantSeg.long2, lat, long)
	}
	if len(gathered.paths) != len(wantSeg.paths) {
		T.Fatalf("expected %d paths in segment, got %d", len(wantSeg.paths),
			len(gathered.paths))
	}
	for pathX, wantPath := range wantSeg.paths {
		gotPath := gathered.paths[pathX]
		if gotPath.path.Name() != wantPath.name {
			T.Fatalf("path %d: expected name '%s', got '%s'", pathX, wantPath.name,
				gotPath.path.Name())
		}
		points, _, forward := gotPath.points()
		if forward != wantPath.forward {
			T.Fatalf("path %d (%s): expected forward=%t, got %t", pathX, wantPath.name,
				wantPath.forward, forward)
		}
		if len(points) != len(wantPath.points) {
			T.Fatalf("path %d (%s): expected %d point values, got %d", pathX,
				wantPath.name, len(wantPath.points), len(points))
		}
		for i := 0; i < len(wantPath.points); i += 2 {
			wantLat, wantLong := wantPath.points[i], wantPath.points[i+1]
			gotLat, gotLong := points[i], points[i+1]
			if !isSamePoint(wantLat, wantLong, gotLat, gotLong) {
				T.Fatalf("path %d (%s) point %d: expected [%f %f], got [%f %f]",
					pathX, wantPath.name, i >> 1, wantLat, wantLong,
					gotLat, gotLong)
			}
		}
	}
}



func Test_gatherSegmentTypical(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentSinglePath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentSinglePathWaypointBefore(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp 1.1 1.2)
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentSinglePathWaypointAfter(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(point wp 3.1 3.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentSinglePathWaypointBeforeMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp 2.1 2.2)
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{2.1, 2.2, 4.1, 4.2, []gsPath{
		{"one", true, locationPairs{2.1, 2.2, 3.1, 3.2, 4.1, 4.2}},
	}})
}


func Test_gatherSegmentSinglePathWaypointAfterMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(point wp 3.1 3.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentSinglePathWaypointBeforeAndAfterMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp 2.1 2.2)
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(marker wp2 3.1 3.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{2.1, 2.2, 3.1, 3.2, []gsPath{
		{"one", true, locationPairs{2.1, 2.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointBefore(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp 1.1 1.2)
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointAfter(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(point wp 3.1 3.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointBeforeAndAfter(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp1 1.1 1.2)
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(point wp2 3.1 3.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointBeforeAndAfterInMiddle(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp1 1.1 1.2)
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(point wp2 2.1 2.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 2.1, 2.2, []gsPath{
		{"one", false, locationPairs{2.1, 2.2, 1.1, 1.2}},
	}})
}


func Test_gatherSegmentReversedSinglePathReversedByInteriorWaypoints(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(point wp1 2.1 2.2)
			(path one
				5.1 5.2
				4.1 4.2
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(point wp2 4.1 4.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{2.1, 2.2, 4.1, 4.2, []gsPath{
		{"one", false, locationPairs{4.1, 4.2, 3.1, 3.2, 2.1, 2.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointBefore(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(circle wp 1.1 1.2 (pixels 4))
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointBeforeMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(circle wp 2.1 2.2 (pixels 4))
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{2.1, 2.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointMiddleMidFirstPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(point wp 2.1 2.2)
			(path two
				2.1 2.2
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2}},
		{"two", true, locationPairs{2.1, 2.2, 3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointMiddleMidSecondPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(point wp 4.1 4.2)
			(path two
				2.1 2.2
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2, 4.1, 4.2}},
		{"two", true, locationPairs{4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointMiddleMidBothPaths(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(point wp 3.1 3.2)
			(path two
				2.3 2.4
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointEnd(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(circle wp 5.1 5.2 (pixels 4))
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointEndMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(circle wp 4.1 4.2 (pixels 4))
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 4.1, 4.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2}},
	}})
}


func Test_gatherSegmentTwoPathsFlipFirstPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
	}})
}


func Test_gatherSegmentTwoPathsFlipSecondPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				5.1 5.2
				4.1 4.2
				3.1 3.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", false, locationPairs{5.1, 5.2, 4.1, 4.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentTwoPathsFlipSecondPathWaypointAmbiguitySolved(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(point wp 4.1 4.2)
			(path two
				5.1 5.2
				4.1 4.2      ;undefined whether to go forward or back!
				3.1 3.2
			)
			(point wp2 5.1 5.2)  ;the user needs to supply ending waypoint
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2, 4.1, 4.2}},
		{"two", false, locationPairs{5.1, 5.2, 4.1, 4.2}},
	}})
}


func Test_gatherSegmentTwoPathsFlipSecondPathWaypointAmbiguitySolvedOtherWay(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(point wp 4.1 4.2)
			(path two
				5.1 5.2
				4.1 4.2
				3.1 3.2
			)
			(point wp2 3.1 3.2)  ;user chose other disambiguation
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 3.1, 3.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2, 4.1, 4.2}},
		{"two", true, locationPairs{4.1, 4.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentTwoPathsFlippedByMiddleWaypoint(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(marker 3.1 3.2)
			(path two
				5.1 5.2
				4.1 4.2
				3.1 3.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 5.1, 5.2, []gsPath{
		{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
		{"two", false, locationPairs{5.1, 5.2, 4.1, 4.2, 3.1, 3.2}},
	}})
}


func Test_gatherSegmentThreePaths(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 7.1, 7.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
		{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}},
	}})
}


func Test_gatherSegmentThreePathsFirstReversed(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 7.1, 7.2, []gsPath{
		{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
		{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}},
	}})
}


func Test_gatherSegmentThreePathsSecondReversed(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				5.1 5.2
				4.1 4.2
				3.1 3.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 7.1, 7.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", false, locationPairs{5.1, 5.2, 4.1, 4.2, 3.1, 3.2}},
		{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}},
	}})
}


func Test_gatherSegmentThreePathsThirdReversed(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment test
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				7.1 7.2
				6.1 6.2
				5.1 5.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredSegmentPaths(T, vd, "test", gsCheck{1.1, 1.2, 7.1, 7.2, []gsPath{
		{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
		{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
		{"three", false, locationPairs{7.1, 7.2, 6.1, 6.2, 5.1, 5.2}},
	}})
}







func Test_gatherRouteContinuityOfSegments(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				7.1 7.2
				8.1 8.2
				9.1 9.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", true, locationPairs{7.1, 7.2, 8.1, 8.2, 9.1, 9.2}}}},
	})
}


func Test_gatherRouteFirstSegmentFlipsSecondSegment(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
}


func Test_gatherRouteSecondSegmentFlipsFirstSegment(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment roadSeg1
			(path three
				7.1 7.2
				6.1 6.2
				5.1 5.2
			)
			(path two
				5.1 5.2
				4.1 4.2
				3.1 3.2
			)
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
		)
		(segment roadSeg2
			(path four
				7.1 7.2
				8.1 8.2
				9.1 9.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
			{"two", false, locationPairs{5.1, 5.2, 4.1, 4.2, 3.1, 3.2}},
			{"three", false, locationPairs{7.1, 7.2, 6.1, 6.2, 5.1, 5.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", true, locationPairs{7.1, 7.2, 8.1, 8.2, 9.1, 9.2}}}},
	})
}


func Test_gatherRouteSecondSegmentFlipsFirstSegmentWithFlippedPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(route road
		(segment roadSeg1
			(path three
				7.1 7.2
				6.1 6.2
				5.1 5.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
		)
		(segment roadSeg2
			(path four
				7.1 7.2
				8.1 8.2
				9.1 9.2
			)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", false, locationPairs{7.1, 7.2, 6.1, 6.2, 5.1, 5.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", true, locationPairs{7.1, 7.2, 8.1, 8.2, 9.1, 9.2}}}},
	})
}


func Test_gatherSideRouteJoiningMainRouteAtStart(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road sideRoute)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	(route sideRoute
		(segment leadIn
			(path dogleg
				1.5 1.6
				1.3 1.4
				1.1 1.2
			)
			(point 1.1 1.2)
			(paths one two three)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
	checkGatheredRouteSegments(T, vd, "sideRoute", []gsCheck{
		{1.5, 1.6, 7.1, 7.2, []gsPath{
			{"dogleg", true, locationPairs{1.5, 1.6, 1.3, 1.4, 1.1, 1.2}},
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
	})
}


func Test_gatherSideRouteJoiningMainRouteAtMidFirstPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road sideRoute)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	(route sideRoute
		(segment leadIn
			(path dogleg
				1.5 1.6
				1.3 1.4
				2.1 2.2
			)
			(point 2.1 2.2)
			(paths one two three)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
	checkGatheredRouteSegments(T, vd, "sideRoute", []gsCheck{
		{1.5, 1.6, 7.1, 7.2, []gsPath{
			{"dogleg", true, locationPairs{1.5, 1.6, 1.3, 1.4, 2.1, 2.2}},
			{"one", true, locationPairs{2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
	})
}


func Test_gatherSideRouteJoiningMainRouteAtEndFirstPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road sideRoute)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	(route sideRoute
		(segment leadIn
			(path dogleg
				1.5 1.6
				1.3 1.4
				3.1 3.2
			)
			(point joinpoint 3.1 3.2)
			(paths one two three)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
	checkGatheredRouteSegments(T, vd, "sideRoute", []gsCheck{
		{1.5, 1.6, 7.1, 7.2, []gsPath{
			{"dogleg", true, locationPairs{1.5, 1.6, 1.3, 1.4, 3.1, 3.2}},
			{"one", true, locationPairs{3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
	})
}


func Test_gatherSideRouteJoiningMainRouteAtStartSecondPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road sideRoute)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	(route sideRoute
		(segment leadIn
			(path dogleg
				1.5 1.6
				1.3 1.4
				3.1 3.2
			)
			(point joinpoint 3.1 3.2)
			(paths two three)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
	checkGatheredRouteSegments(T, vd, "sideRoute", []gsCheck{
		{1.5, 1.6, 7.1, 7.2, []gsPath{
			{"dogleg", true, locationPairs{1.5, 1.6, 1.3, 1.4, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
	})
}


func Test_gatherSideRouteJoiningTwoDoglegs(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road sideRoute)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	(route sideRoute
		(segment house1_to_house2
			(path toHouse1
				1.5 1.6
				1.3 1.4
				2.1 2.2
			)
			(point turn1 2.1 2.2)
			(paths one two three)
			(point turn2 6.1 6.2)
			(path toHouse2
				6.1 6.2
				6.3 6.4
				6.5 6.6
			)
			(marker house2 6.5 6.6)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", true, locationPairs{1.1, 1.2, 2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
	checkGatheredRouteSegments(T, vd, "sideRoute", []gsCheck{
		{1.5, 1.6, 6.5, 6.6, []gsPath{
			{"toHouse1", true, locationPairs{1.5, 1.6, 1.3, 1.4, 2.1, 2.2}},
			{"one", true, locationPairs{2.1, 2.2, 3.1, 3.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2}},
			{"toHouse2", true, locationPairs{6.1, 6.2, 6.3, 6.4, 6.5, 6.6}}}},
	})
}


func Test_gatherSideRouteJoiningTwoDoglegsReversedMainPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road sideRoute)
		)
	)
	(route road
		(segment roadSeg1
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
			)
			(path two
				3.1 3.2
				4.1 4.2
				5.1 5.2
			)
			(path three
				5.1 5.2
				6.1 6.2
				7.1 7.2
			)
		)
		(segment roadSeg2
			(path four
				9.1 9.2
				8.1 8.2
				7.1 7.2
			)
		)
	)
	(route sideRoute
		(segment house1_to_house2
			(path toHouse1
				1.5 1.6
				1.3 1.4
				2.1 2.2
			)
			(point turn1 2.1 2.2)
			(paths one two three)
			(point turn2 6.1 6.2)
			(path toHouse2
				6.1 6.2
				6.3 6.4
				6.5 6.6
			)
			(marker house2 6.5 6.6)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1.1, 1.2, 7.1, 7.2, []gsPath{
			{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2, 1.1, 1.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2, 7.1, 7.2}}}},
		{7.1, 7.2, 9.1, 9.2, []gsPath{
			{"four", false, locationPairs{9.1, 9.2, 8.1, 8.2, 7.1, 7.2}}}},
	})
	checkGatheredRouteSegments(T, vd, "sideRoute", []gsCheck{
		{1.5, 1.6, 6.5, 6.6, []gsPath{
			{"toHouse1", true, locationPairs{1.5, 1.6, 1.3, 1.4, 2.1, 2.2}},
			{"one", false, locationPairs{3.1, 3.2, 2.1, 2.2}},
			{"two", true, locationPairs{3.1, 3.2, 4.1, 4.2, 5.1, 5.2}},
			{"three", true, locationPairs{5.1, 5.2, 6.1, 6.2}},
			{"toHouse2", true, locationPairs{6.1, 6.2, 6.3, 6.4, 6.5, 6.6}}}},
	})
}

