// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Path-threading tests



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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0}, {3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
	{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 0, 4, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 1, []latlongRefProto{
	{1100000, 1200000, 0, 0, 0},{1100000, 1200000, 1, 0, 0},{3100000, 3200000, 1, 4, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp", latlongType{1100000, 1200000},
			latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 1, []latlongRefProto{
	{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPoint, "wp", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_illegalGatherSegmentSinglePathWaypointBeforeMidPath(T *testing.T) {
	//This construction was legal under the older threading system which assumed the
	//lexical direction of a path implied its endpoints.  Placing a waypoint before
	//a path in this way still sets one end of the following path, but now this
	//leaves the threading mechanism unable to resolve which endpoint to use.
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	vd := prepareAndParseStringsIgnoreThreadingError(T, sourceText)
	checkDeferredErrors(T, vd,
		"infile0:10: cannot determine free endpoint of path one under segment test")
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{2100000, 2200000},
	latlongType{4100000, 4200000}, 0, 1, []latlongRefProto{
	{2100000, 2200000, 0, 0, 0},{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 4, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
			latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
			[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
	}})
}


func Test_gatherSegmentSinglePathWaypointBeforeMidPath(T *testing.T) {
	//Resolves the problem with the above example by adding a trailing waypoint
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
		(segment test
			(point wp 2.1 2.2)
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(point wp2 4.1 4.2)
		)
	)
	`
	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{2100000, 2200000},
	latlongType{4100000, 4200000}, 0, 2, []latlongRefProto{
	{2100000, 2200000, 0, 0, 0},{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 4, 0},
	{4100000, 4200000, 2, 0, 0}}, []any{
		miThreadCheck{mitPoint, "wp", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
			latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
			[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{4100000, 4200000},
			latlongType{4100000, 4200000}, 0, 0,
			[]latlongRefProto{{4100000, 4200000, 0, 0, 0}}, []any{4100000, 4200000}},
	}})
}


func Test_gatherSegmentSinglePathWaypointBeforeAndAfterMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"], 
	miThreadCheck{mitSegment, "test", latlongType{2100000, 2200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 2, 0},{3100000, 3200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
			latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
			[]any{2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitMarker, "wp2", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypoints(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
		(segment test
			(point wp 1.1 1.2)
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
	{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 1, 0, 0},{1100000, 1200000, 1, 4, 0},
	{3100000, 3200000, 2, 0, 0}}, []any{
		miThreadCheck{mitPoint, "wp", latlongType{1100000, 1200000},
			latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointBefore(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 1, 0, 0},{1100000, 1200000, 1, 4, 0}}, []any{
		miThreadCheck{mitPoint, "wp", latlongType{1100000, 1200000},
			latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointAfter(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 1, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{1100000, 1200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPoint, "wp", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointBeforeAndAfter(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
	{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 1, 0, 0},{1100000, 1200000, 1, 4, 0},
	{3100000, 3200000, 2, 0, 0}}, []any{
		miThreadCheck{mitPoint, "wp1", latlongType{1100000, 1200000},
			latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_gatherSegmentSinglePathReversedByWaypointBeforeAndAfterInMiddle(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{2100000, 2200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{1100000, 1200000, 1, 2, 0},{2100000, 2200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp1", latlongType{1100000, 1200000},
			latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
		miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
			latlongType{1100000, 1200000}, 0, 2, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{1100000, 1200000, 2, 0, 0}},
			[]any{2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
	}})
}


func Test_gatherSegmentReversedSinglePathReversedByInteriorWaypoints(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{2100000, 2200000},
	latlongType{4100000, 4200000}, 0, 2, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{4100000, 4200000, 1, 0, 0},{2100000, 2200000, 1, 4, 0},{4100000, 4200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitPath, "one:1", latlongType{4100000, 4200000},
			latlongType{2100000, 2200000}, 0, 4, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0},{2100000, 2200000, 4, 0, 0}},
			[]any{4100000, 4200000, 3100000, 3200000, 2100000, 2200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{4100000, 4200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0}}, []any{4100000, 4200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointBefore(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{1100000, 1200000, 1, 0, 0},{3100000, 3200000, 1, 4, 0},{3100000, 3200000, 2, 0, 0},
	{5100000, 5200000, 2, 4, 0}}, []any{
		miThreadCheck{mitCircle, "wp", latlongType{1100000, 1200000},
			latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointBeforeMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{2100000, 2200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{
	{2100000, 2200000, 0, 0, 0},{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 2, 0},
	{3100000, 3200000, 2, 0, 0},{5100000, 5200000, 2, 4, 0}}, []any{
		miThreadCheck{mitCircle, "wp", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
			latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
			[]any{2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointMiddleMidFirstPath(T *testing.T) {
	// Note that this construct is illegal since it leaves path one without a
	// definite endpoint.
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	vd := prepareAndParseStringsIgnoreThreadingError(T, sourceText)
	checkDeferredErrors(T, vd,
		"infile0:9: cannot determine free endpoint of path one under segment test")
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 2, 0},{2100000, 2200000, 1, 0, 0},{2100000, 2200000, 2, 0, 0},
	{3100000, 3200000, 2, 2, 0},{5100000, 5200000, 2, 6, 0}}, []any{
		miThreadCheck{mitPath, "one:1", latlongType{1100000, 1200000},
			latlongType{2100000, 2200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000}},
		miThreadCheck{mitPoint, "wp", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitPath, "two", latlongType{2100000, 2200000},
			latlongType{5100000, 5200000}, 0, 6, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0},
			{5100000, 5200000, 6, 0, 0}},
			[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000,
			5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointMiddleMidSecondPath(T *testing.T) {
	// Note that this construct is illegal since it leaves path two without a
	// definite endpoint.
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	vd := prepareAndParseStringsIgnoreThreadingError(T, sourceText)
	checkDeferredErrors(T, vd,
		"infile0:16: cannot determine free endpoint of path two under segment test")
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 2, 0},{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0},
	{4100000, 4200000, 1, 0, 0},{4100000, 4200000, 2, 0, 0},{5100000, 5200000, 2, 2, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
			{3100000, 3200000, 4, 0, 0},{4100000, 4200000, 6, 0, 0}}, []any{
			1100000, 1200000, 2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
		miThreadCheck{mitPoint, "wp", latlongType{4100000, 4200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0}}, []any{4100000, 4200000}},
		miThreadCheck{mitPath, "two:1", latlongType{4100000, 4200000},
			latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0},{5100000, 5200000, 2, 0, 0}},
			[]any{4100000, 4200000, 5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointMiddleMidBothPaths(T *testing.T) {
	// Note that this construct is illegal since it leaves neither path with a
	// definite endpoint.
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	vd := prepareAndParseStringsIgnoreThreadingError(T, sourceText)
	checkDeferredErrors(T, vd,
		"infile0:9: cannot determine free endpoint of path one under segment test\n" +
		"infile0:16: cannot determine free endpoint of path two under segment test")
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{3100000, 3200000, 2, 0, 0},
	{4100000, 4200000, 2, 2, 0},{5100000, 5200000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "one:1", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPoint, "wp", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
		miThreadCheck{mitPath, "two:1", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{4100000, 4200000, 2, 0, 0},
			{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointEnd(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
	{5100000, 5200000, 2, 0, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
		miThreadCheck{mitCircle, "wp", latlongType{5100000, 5200000},
			latlongType{5100000, 5200000}, 0, 0, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0}}, []any{5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsWaypointEndMidPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{4100000, 4200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{4100000, 4200000, 1, 2, 0},
	{4100000, 4200000, 2, 0, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two:1", latlongType{3100000, 3200000},
			latlongType{4100000, 4200000}, 0, 2, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{4100000, 4200000, 2, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000}},
		miThreadCheck{mitCircle, "wp", latlongType{4100000, 4200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0}}, []any{4100000, 4200000}},
	}})
}


func Test_gatherSegmentTwoPathsFlipFirstPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 1, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{1100000, 1200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsFlipSecondPath(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{5100000, 5200000, 1, 0, 0},{3100000, 3200000, 1, 4, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{5100000, 5200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{5100000, 5200000, 4100000, 4200000, 3100000, 3200000}},
	}})
}


func Test_gatherSegmentTwoPathsFlipSecondPathWaypointAmbiguitySolved(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 3, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
	{5100000, 5200000, 2, 0, 0},{4100000, 4200000, 2, 2, 0},{5100000, 5200000, 3, 0, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0},
			{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000, 2100000, 2200000,
			3100000, 3200000, 4100000, 4200000}},
		miThreadCheck{mitPoint, "wp", latlongType{4100000, 4200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0}}, []any{4100000, 4200000}},
		miThreadCheck{mitPath, "two:1", latlongType{5100000, 5200000},
			latlongType{4100000, 4200000}, 0, 2, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{4100000, 4200000, 2, 0, 0}},
			[]any{5100000, 5200000, 4100000, 4200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{5100000, 5200000},
			latlongType{5100000, 5200000}, 0, 0, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0}}, []any{5100000, 5200000}},
	}})
}


func Test_gatherSegmentTwoPathsFlipSecondPathWaypointAmbiguitySolvedOtherWay(T *testing.T) {
	// Note that this sets up a cycle at point [3.1 3.2]
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{3100000, 3200000}, 0, 3, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
	{4100000, 4200000, 2, 0, 0},{3100000, 3200000, 2, 2, 0},{3100000, 3200000, 3, 0, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0},
			{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000, 2100000, 2200000,
			3100000, 3200000, 4100000, 4200000}},
		miThreadCheck{mitPoint, "wp", latlongType{4100000, 4200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0}}, []any{4100000, 4200000}},
		miThreadCheck{mitPath, "two:1", latlongType{4100000, 4200000},
			latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
			[]any{4100000, 4200000, 3100000, 3200000}},
		miThreadCheck{mitPoint, "wp2", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_gatherSegmentTwoPathsFlippedByMiddleWaypoint(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 2, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{1100000, 1200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 2, 0, 0},
	{3100000, 3200000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitMarker, "$5", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{5100000, 5200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{5100000, 5200000, 4100000, 4200000, 3100000, 3200000}},
	}})
}


func Test_gatherSegmentTwoPathsFlippedWithoutWaypoint(T *testing.T) {
	// Similar to Test_gatherSegmentTwoPathsFlippedByMiddleWaypoint above but without
	// use of a waypoint between the paths
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
		(segment test
			(path one
				3.1 3.2
				2.1 2.2
				1.1 1.2
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{5100000, 5200000}, 0, 1, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{1100000, 1200000, 0, 4, 0},{5100000, 5200000, 1, 0, 0},{3100000, 3200000, 1, 4, 0}},
	[]any{
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPath, "two", latlongType{5100000, 5200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{5100000, 5200000, 4100000, 4200000, 3100000, 3200000}},
	}})
}


func Test_gatherSegmentThreePaths(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
	{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
		miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
			latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
			[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
	}})
}


func Test_gatherSegmentThreePathsFirstReversed(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{1100000, 1200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
	{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
			latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
			[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
		miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
			latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
			[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
	}})
}


func Test_gatherSegmentThreePathsSecondReversed(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{5100000, 5200000, 1, 0, 0},{3100000, 3200000, 1, 4, 0},
	{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{5100000, 5200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{5100000, 5200000, 4100000, 4200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
			latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
			{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
			[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
	}})
}


func Test_gatherSegmentThreePathsThirdReversed(T *testing.T) {
	sourceText := `(layers
		(layer ll
			(menuitem "here")
			(features road)
		)
	)
	(feature road
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
	checkThreadableMapItem(T, vd.mapItems["test"],
	miThreadCheck{mitSegment, "test", latlongType{1100000, 1200000},
	latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
	{7100000, 7200000, 2, 0, 0},{5100000, 5200000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
			latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
			[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
		miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
		miThreadCheck{mitPath, "three", latlongType{7100000, 7200000},
			latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
			[]any{7100000, 7200000, 6100000, 6200000, 5100000, 5200000}},
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},{5100000, 5200000, 0, 1, 4},
	{5100000, 5200000, 0, 2, 0},{7100000, 7200000, 0, 2, 4},{7100000, 7200000, 1, 0, 0},
	{9100000, 9200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 0, 4, 0},
			{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
			{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{7100000, 7200000},
			latlongType{9100000, 9200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{9100000, 9200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{7100000, 7200000},
				latlongType{9100000, 9200000}, 0, 4, []latlongRefProto{
				{7100000, 7200000, 0, 0, 0},{9100000, 9200000, 4, 0, 0}},
				[]any{7100000, 7200000, 8100000, 8200000, 9100000, 9200000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},{5100000, 5200000, 0, 1, 4},
	{5100000, 5200000, 0, 2, 0},{7100000, 7200000, 0, 2, 4},{9100000, 9200000, 1, 0, 0},
	{7100000, 7200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 0, 4, 0},
			{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
			{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{9100000, 9200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{9100000, 9200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{9100000, 9200000, 8100000, 8200000, 7100000, 7200000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{7100000, 7200000, 0, 0, 0},
	{5100000, 5200000, 0, 0, 4},{5100000, 5200000, 0, 1, 0},{3100000, 3200000, 0, 1, 4},
	{3100000, 3200000, 0, 2, 0},{1100000, 1200000, 0, 2, 4},{7100000, 7200000, 1, 0, 0},
	{9100000, 9200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{7100000, 7200000},
			latlongType{1100000, 1200000}, 0, 2, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{5100000, 5200000, 0, 4, 0},
			{5100000, 5200000, 1, 0, 0},{3100000, 3200000, 1, 4, 0},
			{3100000, 3200000, 2, 0, 0},{1100000, 1200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "three", latlongType{7100000, 7200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{7100000, 7200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{7100000, 7200000, 6100000, 6200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "two", latlongType{5100000, 5200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{5100000, 5200000, 4100000, 4200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
				latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
				[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{7100000, 7200000},
			latlongType{9100000, 9200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{9100000, 9200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{7100000, 7200000},
				latlongType{9100000, 9200000}, 0, 4, []latlongRefProto{
				{7100000, 7200000, 0, 0, 0},{9100000, 9200000, 4, 0, 0}},
				[]any{7100000, 7200000, 8100000, 8200000, 9100000, 9200000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{7100000, 7200000, 0, 0, 0},
	{5100000, 5200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},{5100000, 5200000, 0, 1, 4},
	{3100000, 3200000, 0, 2, 0},{1100000, 1200000, 0, 2, 4},{7100000, 7200000, 1, 0, 0},
	{9100000, 9200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{7100000, 7200000},
			latlongType{1100000, 1200000}, 0, 2, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{5100000, 5200000, 0, 4, 0},
			{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
			{3100000, 3200000, 2, 0, 0},{1100000, 1200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "three", latlongType{7100000, 7200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{7100000, 7200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{7100000, 7200000, 6100000, 6200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
				latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
				[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{7100000, 7200000},
			latlongType{9100000, 9200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{9100000, 9200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{7100000, 7200000},
				latlongType{9100000, 9200000}, 0, 4, []latlongRefProto{
				{7100000, 7200000, 0, 0, 0},{9100000, 9200000, 4, 0, 0}},
				[]any{7100000, 7200000, 8100000, 8200000, 9100000, 9200000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},{5100000, 5200000, 0, 1, 4},
	{5100000, 5200000, 0, 2, 0},{7100000, 7200000, 0, 2, 4},{9100000, 9200000, 1, 0, 0},
	{7100000, 7200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 0, 4, 0},
			{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
			{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{9100000, 9200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{9100000, 9200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{9100000, 9200000, 8100000, 8200000, 7100000, 7200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["sideRoute"],
	miThreadCheck{mitRoute, "sideRoute", latlongType{1500000, 1600000},
	latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{{1500000, 1600000, 0, 0, 0},
	{1100000, 1200000, 0, 0, 4},{1100000, 1200000, 0, 1, 0},{1100000, 1200000, 0, 2, 0},
	{3100000, 3200000, 0, 2, 4},{3100000, 3200000, 0, 3, 0},{5100000, 5200000, 0, 3, 4},
	{5100000, 5200000, 0, 4, 0},{7100000, 7200000, 0, 4, 4}}, []any{
		miThreadCheck{mitSegment, "leadIn", latlongType{1500000, 1600000},
			latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
			{1500000, 1600000, 0, 0, 0},{1100000, 1200000, 0, 4, 0},
			{1100000, 1200000, 1, 0, 0},{1100000, 1200000, 2, 0, 0},
			{3100000, 3200000, 2, 4, 0},{3100000, 3200000, 3, 0, 0},
			{5100000, 5200000, 3, 4, 0},{5100000, 5200000, 4, 0, 0},
			{7100000, 7200000, 4, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "dogleg", latlongType{1500000, 1600000},
				latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
				{1500000, 1600000, 0, 0, 0},{1100000, 1200000, 4, 0, 0}},
				[]any{1500000, 1600000, 1300000, 1400000, 1100000, 1200000}},
			miThreadCheck{mitPoint, "$12", latlongType{1100000, 1200000},
				latlongType{1100000, 1200000}, 0, 0, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0}}, []any{1100000, 1200000}},
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0,0, 2},{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},
	{5100000, 5200000, 0, 1, 4},{5100000, 5200000, 0, 2, 0},{7100000, 7200000, 0, 2, 4},
	{9100000, 9200000, 1, 0, 0},{7100000, 7200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},
			{5100000, 5200000, 1, 4, 0},{5100000, 5200000, 2, 0, 0},
			{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{9100000, 9200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{9100000, 9200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{9100000, 9200000, 8100000, 8200000, 7100000, 7200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["sideRoute"],
	miThreadCheck{mitRoute, "sideRoute", latlongType{1500000, 1600000},
	latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{{1500000, 1600000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 4},{2100000, 2200000, 0, 1, 0},{2100000, 2200000, 0, 2, 0},
	{3100000, 3200000, 0, 2, 2},{3100000, 3200000, 0, 3, 0},{5100000, 5200000, 0, 3, 4},
	{5100000, 5200000, 0, 4, 0},{7100000, 7200000, 0, 4, 4}}, []any{
		miThreadCheck{mitSegment, "leadIn", latlongType{1500000, 1600000},
			latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
			{1500000, 1600000, 0, 0, 0},{2100000, 2200000, 0, 4, 0},
			{2100000, 2200000, 1, 0, 0},{2100000, 2200000, 2, 0, 0},
			{3100000, 3200000, 2, 2, 0},{3100000, 3200000, 3, 0, 0},
			{5100000, 5200000, 3, 4, 0},{5100000, 5200000, 4, 0, 0},
			{7100000, 7200000, 4, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "dogleg", latlongType{1500000, 1600000},
				latlongType{2100000, 2200000}, 0, 4, []latlongRefProto{
				{1500000, 1600000, 0, 0, 0},{2100000, 2200000, 4, 0, 0}},
				[]any{1500000, 1600000, 1300000, 1400000, 2100000, 2200000}},
			miThreadCheck{mitPoint, "$12", latlongType{2100000, 2200000},
				latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
	}})
}


func Test_gatherSideRouteJoiningMainRouteAtEndFirstPath(T *testing.T) {
	// This formulation is illegal under the new threading model, which avoids
	// zero-length paths (i.e. paths which contain only one point).  The difficulty
	// here is that the 'leadIn' segment references path 'one' at a point shared
	// with the next path, 'two'.  This confuses the threading system.
	// The rule to follow is that every path or segment named in a (paths) or (segments)
	// list must contribute at least two points to the item being formed.
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
	vd := prepareAndParseStringsIgnoreThreadingError(T, sourceText)
	checkDeferredErrors(T, vd, "infile0:41: path two does not connect with segment leadIn " +
		"(two is defined at infile0:14)")
	// The case introduces too many pathologies to consider for unit tests
}


func Test_gatherSideRouteJoiningMainRouteAtStartSecondPath(T *testing.T) {
	// This is the legal version of the preceding test case
	// Test_gatherSideRouteJoiningMainRouteAtEndFirstPath.  The current test case
	// omits the reference to path one in the leadIn segment
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},{5100000, 5200000, 0, 1, 4},
	{5100000, 5200000, 0, 2, 0},{7100000, 7200000, 0, 2, 4},{9100000, 9200000, 1, 0, 0},
	{7100000, 7200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 0, 4, 0},
			{3100000, 3200000, 1, 0, 0},{5100000, 5200000, 1, 4, 0},
			{5100000, 5200000, 2, 0, 0},{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{9100000, 9200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{9100000, 9200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{9100000, 9200000, 8100000, 8200000, 7100000, 7200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["sideRoute"],
	miThreadCheck{mitRoute, "sideRoute", latlongType{1500000, 1600000},
	latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{{1500000, 1600000, 0, 0, 0},
	{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},{3100000, 3200000, 0, 2, 0},
	{5100000, 5200000, 0, 2, 4},{5100000, 5200000, 0, 3, 0},{7100000, 7200000, 0, 3, 4}},
	[]any{
		miThreadCheck{mitSegment, "leadIn", latlongType{1500000, 1600000},
			latlongType{7100000, 7200000}, 0, 3, []latlongRefProto{
			{1500000, 1600000, 0, 0, 0},{3100000, 3200000, 0, 4, 0},
			{3100000, 3200000, 1, 0, 0},{3100000, 3200000, 2, 0, 0},
			{5100000, 5200000, 2, 4, 0},{5100000, 5200000, 3, 0, 0},
			{7100000, 7200000, 3, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "dogleg", latlongType{1500000, 1600000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1500000, 1600000, 0, 0, 0},{3100000, 3200000, 4, 0, 0}},
				[]any{1500000, 1600000, 1300000, 1400000, 3100000, 3200000}},
			miThreadCheck{mitPoint, "joinpoint", latlongType{3100000, 3200000},
				latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{3100000, 3200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},
	{5100000, 5200000, 0, 1, 4},{5100000, 5200000, 0, 2, 0},{6100000, 6200000, 0, 2, 2},
	{7100000, 7200000, 0, 2, 4},{9100000, 9200000, 1, 0, 0},{7100000, 7200000, 1, 0, 4}},
	[]any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},
			{5100000, 5200000, 1, 4, 0},{5100000, 5200000, 2, 0, 0},
			{6100000, 6200000, 2, 2, 0},{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{3100000, 3200000}, 0, 4, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{3100000, 3200000, 4, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{6100000, 6200000, 2, 0, 0},
				{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{9100000, 9200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{9100000, 9200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{9100000, 9200000, 8100000, 8200000, 7100000, 7200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["sideRoute"],
	miThreadCheck{mitRoute, "sideRoute", latlongType{1500000, 1600000},
	latlongType{6500000, 6600000}, 0, 0, []latlongRefProto{{1500000, 1600000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 4},{2100000, 2200000, 0, 1, 0},{2100000, 2200000, 0, 2, 0},
	{3100000, 3200000, 0, 2, 2},{3100000, 3200000, 0, 3, 0},{5100000, 5200000, 0, 3, 4},
	{5100000, 5200000, 0, 4, 0},{6100000, 6200000, 0, 4, 2},{6100000, 6200000, 0, 5, 0},
	{6100000, 6200000, 0, 6, 0},{6500000, 6600000, 0, 6, 4},{6500000, 6600000, 0, 7, 0}},
	[]any{
		miThreadCheck{mitSegment, "house1_to_house2", latlongType{1500000, 1600000},
			latlongType{6500000, 6600000}, 0, 7, []latlongRefProto{
			{1500000, 1600000, 0, 0, 0},{2100000, 2200000, 0, 4, 0},
			{2100000, 2200000, 1, 0, 0},{2100000, 2200000, 2, 0, 0},
			{3100000, 3200000, 2, 2, 0},{3100000, 3200000, 3, 0, 0},
			{5100000, 5200000, 3, 4, 0},{5100000, 5200000, 4, 0, 0},
			{6100000, 6200000, 4, 2, 0},{6100000, 6200000, 5, 0, 0},
			{6100000, 6200000, 6, 0, 0},{6500000, 6600000, 6, 4, 0},
			{6500000, 6600000, 7, 0, 0}},
			[]any{
			miThreadCheck{mitPath, "toHouse1", latlongType{1500000, 1600000},
				latlongType{2100000, 2200000}, 0, 4, []latlongRefProto{
				{1500000, 1600000, 0, 0, 0},{2100000, 2200000, 4, 0, 0}},
				[]any{1500000, 1600000, 1300000, 1400000, 2100000, 2200000}},
			miThreadCheck{mitPoint, "turn1", latlongType{2100000, 2200000},
				latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three:1", latlongType{5100000, 5200000},
				latlongType{6100000, 6200000}, 0, 2, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{6100000, 6200000, 2, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000}},
			miThreadCheck{mitPoint, "turn2", latlongType{6100000, 6200000},
				latlongType{6100000, 6200000}, 0, 0, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0}}, []any{6100000, 6200000}},
			miThreadCheck{mitPath, "toHouse2", latlongType{6100000, 6200000},
				latlongType{6500000, 6600000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{6500000, 6600000, 4, 0, 0}},
				[]any{6100000, 6200000, 6300000, 6400000, 6500000, 6600000}},
			miThreadCheck{mitMarker, "house2", latlongType{6500000, 6600000},
				latlongType{6500000, 6600000}, 0, 0, []latlongRefProto{
				{6500000, 6600000, 0, 0, 0}}, []any{6500000, 6600000}},
		}},
	}})
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
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{1100000, 1200000, 0, 0, 4},{3100000, 3200000, 0, 1, 0},
	{5100000, 5200000, 0, 1, 4},{5100000, 5200000, 0, 2, 0},{6100000, 6200000, 0, 2, 2},
	{7100000, 7200000, 0, 2, 4},{9100000, 9200000, 1, 0, 0},{7100000, 7200000, 1, 0, 4}},
	[]any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{1100000, 1200000},
			latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{1100000, 1200000, 0, 4, 0},{3100000, 3200000, 1, 0, 0},
			{5100000, 5200000, 1, 4, 0},{5100000, 5200000, 2, 0, 0},
			{6100000, 6200000, 2, 2, 0},{7100000, 7200000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "one", latlongType{3100000, 3200000},
				latlongType{1100000, 1200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{1100000, 1200000, 4, 0, 0}},
				[]any{3100000, 3200000, 2100000, 2200000, 1100000, 1200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three", latlongType{5100000, 5200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{6100000, 6200000, 2, 0, 0},
				{7100000, 7200000, 4, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{9100000, 9200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "four", latlongType{9100000, 9200000},
				latlongType{7100000, 7200000}, 0, 4, []latlongRefProto{
				{9100000, 9200000, 0, 0, 0},{7100000, 7200000, 4, 0, 0}},
				[]any{9100000, 9200000, 8100000, 8200000, 7100000, 7200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["sideRoute"],
	miThreadCheck{mitRoute, "sideRoute", latlongType{1500000, 1600000},
	latlongType{6500000, 6600000}, 0, 0, []latlongRefProto{{1500000, 1600000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 4},{2100000, 2200000, 0, 1, 0},{3100000, 3200000, 0, 2, 0},
	{2100000, 2200000, 0, 2, 2},{3100000, 3200000, 0, 3, 0},{5100000, 5200000, 0, 3, 4},
	{5100000, 5200000, 0, 4, 0},{6100000, 6200000, 0, 4, 2},{6100000, 6200000, 0, 5, 0},
	{6100000, 6200000, 0, 6, 0},{6500000, 6600000, 0, 6, 4},{6500000, 6600000, 0, 7, 0}},
	[]any{
		miThreadCheck{mitSegment, "house1_to_house2", latlongType{1500000, 1600000},
			latlongType{6500000, 6600000}, 0, 7, []latlongRefProto{
			{1500000, 1600000, 0, 0, 0},{2100000, 2200000, 0, 4, 0},
			{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 2, 0, 0},
			{2100000, 2200000, 2, 2, 0},{3100000, 3200000, 3, 0, 0},
			{5100000, 5200000, 3, 4, 0},{5100000, 5200000, 4, 0, 0},
			{6100000, 6200000, 4, 2, 0},{6100000, 6200000, 5, 0, 0},
			{6100000, 6200000, 6, 0, 0},{6500000, 6600000, 6, 4, 0},
			{6500000, 6600000, 7, 0, 0}},
			[]any{
			miThreadCheck{mitPath, "toHouse1", latlongType{1500000, 1600000},
				latlongType{2100000, 2200000}, 0, 4, []latlongRefProto{
				{1500000, 1600000, 0, 0, 0},{2100000, 2200000, 4, 0, 0}},
				[]any{1500000, 1600000, 1300000, 1400000, 2100000, 2200000}},
			miThreadCheck{mitPoint, "turn1", latlongType{2100000, 2200000},
				latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
			miThreadCheck{mitPath, "one:1", latlongType{3100000, 3200000},
				latlongType{2100000, 2200000}, 0, 2, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0}},
				[]any{3100000, 3200000, 2100000, 2200000}},
			miThreadCheck{mitPath, "two", latlongType{3100000, 3200000},
				latlongType{5100000, 5200000}, 0, 4, []latlongRefProto{
				{3100000, 3200000, 0, 0, 0},{5100000, 5200000, 4, 0, 0}},
				[]any{3100000, 3200000, 4100000, 4200000, 5100000, 5200000}},
			miThreadCheck{mitPath, "three:1", latlongType{5100000, 5200000},
				latlongType{6100000, 6200000}, 0, 2, []latlongRefProto{
				{5100000, 5200000, 0, 0, 0},{6100000, 6200000, 2, 0, 0}},
				[]any{5100000, 5200000, 6100000, 6200000}},
			miThreadCheck{mitPoint, "turn2", latlongType{6100000, 6200000},
				latlongType{6100000, 6200000}, 0, 0, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0}}, []any{6100000, 6200000}},
			miThreadCheck{mitPath, "toHouse2", latlongType{6100000, 6200000},
				latlongType{6500000, 6600000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{6500000, 6600000, 4, 0, 0}},
				[]any{6100000, 6200000, 6300000, 6400000, 6500000, 6600000}},
			miThreadCheck{mitMarker, "house2", latlongType{6500000, 6600000},
				latlongType{6500000, 6600000}, 0, 0, []latlongRefProto{
				{6500000, 6600000, 0, 0, 0}}, []any{6500000, 6600000}},
		}},
	}})
}

