// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Path-threading tests for routeSegments


func Test_basicRouteSlice(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  3.1 3.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
		)
	)
	(route partRoad
		(routeSegments road wp1 wp2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{3100000, 3200000, 0, 0, 4},{4100000, 4200000, 0, 0, 6}},
	[]any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0}},
			[]any{
				miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
					latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
					{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
					{3100000, 3200000, 4, 0, 0},{4100000, 4200000, 6, 0, 0}},
					[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000,
					4100000, 4200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 0, 2},{3100000, 3200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 0, 2, 0}}, []any{
				miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
					latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
					{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
					[]any{2100000, 2200000, 3100000, 3200000}},
			}},
		miThreadCheck{mitMarker, "wp2", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_basicRouteSliceLiteralWaypoint1(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(marker wp2  3.1 3.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
		)
	)
	(route partRoad
		(routeSegments road 2.1 2.2 wp2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{3100000, 3200000, 0, 0, 4},{4100000, 4200000, 0, 0, 6}},
	[]any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0}}, []any{
				miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
					latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
					{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
					{3100000, 3200000, 4, 0, 0},{4100000, 4200000, 6, 0, 0}},
					[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000,
					4100000, 4200000}},
			}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 0, 2},{3100000, 3200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "$8", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 0, 2, 0}}, []any{
				miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
					latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
					{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
					[]any{2100000, 2200000, 3100000, 3200000}},
			}},
		miThreadCheck{mitMarker, "wp2", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_basicRouteSliceLiteralWaypoint2(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(marker wp1  2.1 2.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
		)
	)
	(route partRoad
		(routeSegments road wp1 3.1 3.2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{3100000, 3200000, 0, 0, 4},{4100000, 4200000, 0, 0, 6}},
	[]any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0}}, []any{
				miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
					latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
					{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
					{3100000, 3200000, 4, 0, 0},{4100000, 4200000, 6, 0, 0}},
					[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000,
					4100000, 4200000}},
			}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 0, 2},{3100000, 3200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitMarker, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 0, 2, 0}}, []any{
				miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
					latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
					{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
					[]any{2100000, 2200000, 3100000, 3200000}},
			}},
		miThreadCheck{mitPoint, "$8", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_basicRouteSliceLiteralBothWaypoints(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features road partRoad)
		)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
		)
	)
	(route partRoad
		(routeSegments road 2.1 2.2  3.1 3.2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{3100000, 3200000, 0, 0, 4},{4100000, 4200000, 0, 0, 6}},
	[]any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0}}, []any{
				miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
					latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
					{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
					{3100000, 3200000, 4, 0, 0},{4100000, 4200000, 6, 0, 0}},
					[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000,
					4100000, 4200000}},
			}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 0, 2},{3100000, 3200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "$6", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 0, 2, 0}}, []any{
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000}},
		}},
		miThreadCheck{mitPoint, "$7", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
	}})
}


func Test_basicRouteSliceFlipRoute(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features road partRoad)
		)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
		)
	)
	(route partRoad
		(routeSegments road 3.1 3.2  2.1 2.2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{3100000, 3200000, 0, 0, 4},{4100000, 4200000, 0, 0, 6}},
	[]any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{4100000, 4200000}, 0, 0, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{3100000, 3200000, 0, 4, 0},{4100000, 4200000, 0, 6, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{3100000, 3200000, 4, 0, 0},{4100000, 4200000, 6, 0, 0}},
				[]any{1100000, 1200000, 2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{3100000, 3200000},
	latlongType{2100000, 2200000}, 0, 2, []latlongRefProto{{3100000, 3200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{3100000, 3200000, 1, 0, 2},{2100000, 2200000, 2, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "$6", latlongType{3100000, 3200000},
			latlongType{3100000, 3200000}, 0, 0, []latlongRefProto{
			{3100000, 3200000, 0, 0, 0}}, []any{3100000, 3200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{3100000, 3200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 0, 2, 0}}, []any{
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{3100000, 3200000}, 0, 2, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{3100000, 3200000, 2, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000}},
		}},
		miThreadCheck{mitPoint, "$7", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
	}})
}


func Test_sliceTwoSegmentRoute(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  7.1 7.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(path two
				4.1 4.2
				5.1 5.2
				6.1 6.2
			)
		)
		(segment s2
			(path three
				6.1 6.2
				7.1 7.2
				8.1 8.2
			)
		)
	)
	(route partRoad
		(routeSegments road wp1 wp2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{8100000, 8200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{4100000, 4200000, 0, 0, 6},{4100000, 4200000, 0, 1, 0},
	{6100000, 6200000, 0, 1, 4},{6100000, 6200000, 1, 0, 0},{7100000, 7200000, 1, 0, 2},
	{8100000, 8200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
			{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000,
				2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2", latlongType{6100000, 6200000},
			latlongType{8100000, 8200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0},
			{8100000, 8200000, 0, 4, 0}}, []any{
			miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
				latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0},
				{8100000, 8200000, 4, 0, 0}}, []any{6100000, 6200000,
				7100000, 7200000, 8100000, 8200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{7100000, 7200000}, 0, 3, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 0, 4},{4100000, 4200000, 1, 1, 0},
	{6100000, 6200000, 1, 1, 4},{6100000, 6200000, 2, 0, 0},{7100000, 7200000, 2, 0, 2},
	{7100000, 7200000, 3, 0, 0}}, []any{
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 0, 4, 0},
			{4100000, 4200000, 1, 0, 0},{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2:1", latlongType{6100000, 6200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0}}, []any{
				miThreadCheck{mitPath, "three:1", latlongType{6100000, 6200000},
					latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
					{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0}},
					[]any{6100000, 6200000, 7100000, 7200000}},
			}},
		miThreadCheck{mitMarker, "wp2", latlongType{7100000, 7200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0}}, []any{7100000, 7200000}},
	}})
}


func Test_sliceTwoSegmentRouteBothTwoPaths(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  9.1 9.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(path two
				4.1 4.2
				5.1 5.2
				6.1 6.2
			)
		)
		(segment s2
			(path three
				6.1 6.2
				7.1 7.2
				8.1 8.2
			)
			(path four
				8.1 8.2
				9.1 9.2
				10.1 10.2
			)
		)
	)
	(route partRoad
		(routeSegments road wp1 wp2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{10100000, 10200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{4100000, 4200000, 0, 0, 6},{4100000, 4200000, 0, 1, 0},
	{6100000, 6200000, 0, 1, 4},{6100000, 6200000, 1, 0, 0},{8100000, 8200000, 1, 0, 4},
	{8100000, 8200000, 1, 1, 0},{9100000, 9200000, 1, 1, 2},{10100000, 10200000, 1, 1, 4}},
	[]any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
			{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000,
				2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}}, []any{
				4100000, 4200000, 5100000, 5200000, 6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2", latlongType{6100000, 6200000},
			latlongType{10100000, 10200000}, 0, 1, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 0, 4, 0},
			{8100000, 8200000, 1, 0, 0},{9100000, 9200000, 1, 2, 0},
			{10100000, 10200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
				latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 4, 0, 0}}, []any{
				6100000, 6200000, 7100000, 7200000, 8100000, 8200000}},
			miThreadCheck{mitPath, "four", latlongType{8100000, 8200000},
				latlongType{10100000, 10200000}, 0, 4, []latlongRefProto{
				{8100000, 8200000, 0, 0, 0},{9100000, 9200000, 2, 0, 0},
				{10100000, 10200000, 4, 0, 0}}, []any{
				8100000, 8200000, 9100000, 9200000, 10100000, 10200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{9100000, 9200000}, 0, 3, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 0, 4},{4100000, 4200000, 1, 1, 0},
	{6100000, 6200000, 1, 1, 4},{6100000, 6200000, 2, 0, 0},{8100000, 8200000, 2, 0, 4},
	{8100000, 8200000, 2, 1, 0},{9100000, 9200000, 2, 1, 2},{9100000, 9200000, 3, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 0, 4, 0},
			{4100000, 4200000, 1, 0, 0},{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000, 6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2:1", latlongType{6100000, 6200000},
			latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 0, 4, 0},
			{8100000, 8200000, 1, 0, 0},{9100000, 9200000, 1, 2, 0}}, []any{
			miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
				latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 4, 0, 0}},
				[]any{6100000, 6200000, 7100000, 7200000, 8100000, 8200000}},
			miThreadCheck{mitPath, "four:1", latlongType{8100000, 8200000},
				latlongType{9100000, 9200000}, 0, 2, []latlongRefProto{
				{8100000, 8200000, 0, 0, 0},{9100000, 9200000, 2, 0, 0}},
				[]any{8100000, 8200000, 9100000, 9200000}},
		}},
		miThreadCheck{mitMarker, "wp2", latlongType{9100000, 9200000},
			latlongType{9100000, 9200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0}}, []any{9100000, 9200000}},
	}})
}


func Test_sliceTwoSegmentRouteReversed(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  7.1 7.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(path two
				4.1 4.2
				5.1 5.2
				6.1 6.2
			)
		)
		(segment s2
			(path three
				6.1 6.2
				7.1 7.2
				8.1 8.2
			)
		)
	)
	(route partRoad
		(routeSegments road wp2 wp1)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{8100000, 8200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{4100000, 4200000, 0, 0, 6},{4100000, 4200000, 0, 1, 0},
	{6100000, 6200000, 0, 1, 4},{6100000, 6200000, 1, 0, 0},{7100000, 7200000, 1, 0, 2},
	{8100000, 8200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
			{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000,
				2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2", latlongType{6100000, 6200000},
			latlongType{8100000, 8200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0},
			{8100000, 8200000, 0, 4, 0}}, []any{
				miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
					latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
					{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0},
					{8100000, 8200000, 4, 0, 0}}, []any{6100000, 6200000,
					7100000, 7200000, 8100000, 8200000}},
			}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{7100000, 7200000},
	latlongType{2100000, 2200000}, 0, 3, []latlongRefProto{{7100000, 7200000, 0, 0, 0},
	{6100000, 6200000, 1, 0, 0},{7100000, 7200000, 1, 0, 2},{4100000, 4200000, 2, 0, 0},
	{6100000, 6200000, 2, 0, 4},{2100000, 2200000, 2, 1, 0},{4100000, 4200000, 2, 1, 4},
	{2100000, 2200000, 3, 0, 0}}, []any{
		miThreadCheck{mitMarker, "wp2", latlongType{7100000, 7200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0}}, []any{7100000, 7200000}},
		miThreadCheck{mitSegment, "s2:1", latlongType{7100000, 7200000},
			latlongType{6100000, 6200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0}}, []any{
				miThreadCheck{mitPath, "three:1", latlongType{6100000, 6200000},
					latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
					{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0}},
					[]any{6100000, 6200000, 7100000, 7200000}},
			}},
		miThreadCheck{mitSegment, "s1:1", latlongType{6100000, 6200000},
			latlongType{2100000, 2200000}, 0, 1, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 0, 4, 0},
			{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			}},
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
	}})
}


func Test_sliceThreeSegmentRouteToTwoSegments(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  9.1 9.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(path two
				4.1 4.2
				5.1 5.2
				6.1 6.2
			)
		)
		(segment s2
			(path three
				6.1 6.2
				7.1 7.2
				8.1 8.2
			)
			(path four
				8.1 8.2
				9.1 9.2
				10.1 10.2
			)
		)
		(segment s3
			(path five
				10.1 10.2
				11.1 11.2
			)
		)
	)
	(route partRoad
		(routeSegments road wp1 wp2)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{11100000, 11200000}, 0, 2, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{4100000, 4200000, 0, 0, 6},{4100000, 4200000, 0, 1, 0},
	{6100000, 6200000, 0, 1, 4},{6100000, 6200000, 1, 0, 0},{8100000, 8200000, 1, 0, 4},
	{8100000, 8200000, 1, 1, 0},{9100000, 9200000, 1, 1, 2},{10100000, 10200000, 1, 1, 4},
	{10100000, 10200000, 2, 0, 0},{11100000, 11200000, 2, 0, 2}}, []any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
			{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000,
				2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000, 6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2", latlongType{6100000, 6200000},
			latlongType{10100000, 10200000}, 0, 1, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 0, 4, 0},
			{8100000, 8200000, 1, 0, 0},{9100000, 9200000, 1, 2, 0},
			{10100000, 10200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
				latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 4, 0, 0}},
				[]any{6100000, 6200000, 7100000, 7200000, 8100000, 8200000}},
			miThreadCheck{mitPath, "four", latlongType{8100000, 8200000},
				latlongType{10100000, 10200000}, 0, 4, []latlongRefProto{
				{8100000, 8200000, 0, 0, 0},{9100000, 9200000, 2, 0, 0},
				{10100000, 10200000, 4, 0, 0}},
				[]any{8100000, 8200000, 9100000, 9200000, 10100000, 10200000}},
		}},
		miThreadCheck{mitSegment, "s3", latlongType{10100000, 10200000},
			latlongType{11100000, 11200000}, 0, 0, []latlongRefProto{
			{10100000, 10200000, 0, 0, 0},{11100000, 11200000, 0, 2, 0}}, []any{
			miThreadCheck{mitPath, "five", latlongType{10100000, 10200000},
				latlongType{11100000, 11200000}, 0, 2, []latlongRefProto{
				{10100000, 10200000, 0, 0, 0},{11100000, 11200000, 2, 0, 0}},
				[]any{10100000, 10200000, 11100000, 11200000}},
		}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2100000, 2200000},
	latlongType{9100000, 9200000}, 0, 3, []latlongRefProto{{2100000, 2200000, 0, 0, 0},
	{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 0, 4},{4100000, 4200000, 1, 1, 0},
	{6100000, 6200000, 1, 1, 4},{6100000, 6200000, 2, 0, 0},{8100000, 8200000, 2, 0, 4},
	{8100000, 8200000, 2, 1, 0},{9100000, 9200000, 2, 1, 2},{9100000, 9200000, 3, 0, 0}},
	[]any{
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 0, 4, 0},
			{4100000, 4200000, 1, 0, 0},{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000, 6100000, 6200000}},
		}},
		miThreadCheck{mitSegment, "s2:1", latlongType{6100000, 6200000},
			latlongType{9100000, 9200000}, 0, 1, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 0, 4, 0},
			{8100000, 8200000, 1, 0, 0},{9100000, 9200000, 1, 2, 0}}, []any{
			miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
				latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{8100000, 8200000, 4, 0, 0}},
				[]any{6100000, 6200000, 7100000, 7200000, 8100000, 8200000}},
			miThreadCheck{mitPath, "four:1", latlongType{8100000, 8200000},
				latlongType{9100000, 9200000}, 0, 2, []latlongRefProto{
				{8100000, 8200000, 0, 0, 0},{9100000, 9200000, 2, 0, 0}},
				[]any{8100000, 8200000, 9100000, 9200000}},
		}},
		miThreadCheck{mitMarker, "wp2", latlongType{9100000, 9200000},
			latlongType{9100000, 9200000}, 0, 0, []latlongRefProto{
			{9100000, 9200000, 0, 0, 0}}, []any{9100000, 9200000}},
	}})
}


func Test_sliceTwoSegmentRoutePlusAdjoiningSegments(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  7.1 7.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(path two
				4.1 4.2
				5.1 5.2
				6.1 6.2
			)
		)
		(segment s2
			(path three
				6.1 6.2
				7.1 7.2
				8.1 8.2
			)
		)
	)
	(route partRoad
		(segment spur1
			(marker 2.5 2.6)
			(path ps1
				2.5 2.6
				2.3 2.4
				2.1 2.2
			)
		)
		(routeSegments road wp1 wp2)
		(segment spur2
			(path ps2
				7.1 7.2
				7.3 7.4
			)
			(circle (radius 30) 7.3 7.4)
		)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{8100000, 8200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{4100000, 4200000, 0, 0, 6},{4100000, 4200000, 0, 1, 0},
	{6100000, 6200000, 0, 1, 4},{6100000, 6200000, 1, 0, 0},{7100000, 7200000, 1, 0, 2},
	{8100000, 8200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
			{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000,
				2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			}},
		miThreadCheck{mitSegment, "s2", latlongType{6100000, 6200000},
			latlongType{8100000, 8200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0},
			{8100000, 8200000, 0, 4, 0}}, []any{
				miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
					latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
					{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0},
					{8100000, 8200000, 4, 0, 0}}, []any{6100000, 6200000,
					7100000, 7200000, 8100000, 8200000}},
			}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{2500000, 2600000},
	latlongType{7300000, 7400000}, 0, 5, []latlongRefProto{
	{2500000, 2600000, 0, 0, 0},{2500000, 2600000, 0, 1, 0},{2100000, 2200000, 0, 1, 4},
	{2100000, 2200000, 1, 0, 0},
	{2100000, 2200000, 2, 0, 0},{4100000, 4200000, 2, 0, 4},{4100000, 4200000, 2, 1, 0},
	{6100000, 6200000, 2, 1, 4},
	{6100000, 6200000, 3, 0, 0},{7100000, 7200000, 3, 0, 2},
	{7100000, 7200000, 4, 0, 0},
	{7100000, 7200000, 5, 0, 0},{7300000, 7400000, 5, 0, 2},{7300000, 7400000, 5, 1, 0}},
	[]any{
		miThreadCheck{mitSegment, "spur1", latlongType{2500000, 2600000},
			latlongType{2100000, 2200000}, 0, 1, []latlongRefProto{
			{2500000, 2600000, 0, 0, 0},{2500000, 2600000, 1, 0, 0},
			{2100000, 2200000, 1, 4, 0}}, []any{
			miThreadCheck{mitMarker, "$13", latlongType{2500000, 2600000},
				latlongType{2500000, 2600000}, 0, 0, []latlongRefProto{
				{2500000, 2600000, 0, 0, 0}}, []any{2500000, 2600000}},
			miThreadCheck{mitPath, "ps1", latlongType{2500000, 2600000},
				latlongType{2100000, 2200000}, 0, 4, []latlongRefProto{
				{2500000, 2600000, 0, 0, 0},{2100000, 2200000, 4, 0, 0}},
				[]any{2500000, 2600000, 2300000, 2400000, 2100000, 2200000}},
		}},
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "s1:1", latlongType{2100000, 2200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 0, 4, 0},
			{4100000, 4200000, 1, 0, 0},{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			}},
		miThreadCheck{mitSegment, "s2:1", latlongType{6100000, 6200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0}}, []any{
			miThreadCheck{mitPath, "three:1", latlongType{6100000, 6200000},
				latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0}},
				[]any{6100000, 6200000, 7100000, 7200000}},
			}},
		miThreadCheck{mitMarker, "wp2", latlongType{7100000, 7200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0}}, []any{7100000, 7200000}},
		miThreadCheck{mitSegment, "spur2", latlongType{7100000, 7200000},
			latlongType{7300000, 7400000}, 0, 1, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0},{7300000, 7400000, 0, 2, 0},
			{7300000, 7400000, 1, 0, 0}}, []any{
			miThreadCheck{mitPath, "ps2", latlongType{7100000, 7200000},
				latlongType{7300000, 7400000}, 0, 2, []latlongRefProto{
				{7100000, 7200000, 0, 0, 0},{7300000, 7400000, 2, 0, 0}},
				[]any{7100000, 7200000, 7300000, 7400000}},
			miThreadCheck{mitCircle, "$17", latlongType{7300000, 7400000},
				latlongType{7300000, 7400000}, 0, 0, []latlongRefProto{
				{7300000, 7400000, 0, 0, 0}}, []any{7300000, 7400000}},
		}},
	}})
}


func Test_sliceTwoSegmentRoutePlusAdjoiningSegmentsFlipRoute(T *testing.T) {
	sourceText := `(layers
		(layer l1
			(menuitem "test")
			(features points road partRoad)
		)
	)
	(feature points
		(point  wp1  2.1 2.2)
		(marker wp2  7.1 7.2)
	)
	(route road
		(segment s1
			(path one
				1.1 1.2
				2.1 2.2
				3.1 3.2
				4.1 4.2
			)
			(path two
				4.1 4.2
				5.1 5.2
				6.1 6.2
			)
		)
		(segment s2
			(path three
				6.1 6.2
				7.1 7.2
				8.1 8.2
			)
		)
	)
	(route partRoad
		(segment spur2
			(circle (radius 30) 7.3 7.4)
			(path ps2
				7.3 7.4
				7.1 7.2
			)
		)
		(routeSegments road wp2 wp1)
		(segment spur1
			(path ps1
				2.1 2.2
				2.3 2.4
				2.5 2.6
			)
			(marker 2.5 2.6)
		)
	)
	`

	vd := prepareAndParseStrings(T, sourceText)
	checkThreadableMapItem(T, vd.mapItems["road"],
	miThreadCheck{mitRoute, "road", latlongType{1100000, 1200000},
	latlongType{8100000, 8200000}, 0, 1, []latlongRefProto{{1100000, 1200000, 0, 0, 0},
	{2100000, 2200000, 0, 0, 2},{4100000, 4200000, 0, 0, 6},{4100000, 4200000, 0, 1, 0},
	{6100000, 6200000, 0, 1, 4},{6100000, 6200000, 1, 0, 0},{7100000, 7200000, 1, 0, 2},
	{8100000, 8200000, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "s1", latlongType{1100000, 1200000},
			latlongType{6100000, 6200000}, 0, 1, []latlongRefProto{
			{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 0, 2, 0},
			{4100000, 4200000, 0, 6, 0},{4100000, 4200000, 1, 0, 0},
			{6100000, 6200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "one", latlongType{1100000, 1200000},
				latlongType{4100000, 4200000}, 0, 6, []latlongRefProto{
				{1100000, 1200000, 0, 0, 0},{2100000, 2200000, 2, 0, 0},
				{4100000, 4200000, 6, 0, 0}}, []any{1100000, 1200000,
				2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			}},
		miThreadCheck{mitSegment, "s2", latlongType{6100000, 6200000},
			latlongType{8100000, 8200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0},
			{8100000, 8200000, 0, 4, 0}}, []any{
				miThreadCheck{mitPath, "three", latlongType{6100000, 6200000},
					latlongType{8100000, 8200000}, 0, 4, []latlongRefProto{
					{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0},
					{8100000, 8200000, 4, 0, 0}}, []any{6100000, 6200000,
					7100000, 7200000, 8100000, 8200000}},
			}},
	}})
	checkThreadableMapItem(T, vd.mapItems["partRoad"],
	miThreadCheck{mitRoute, "partRoad", latlongType{7300000, 7400000},
	latlongType{2500000, 2600000}, 0, 5, []latlongRefProto{
	{7300000, 7400000, 0, 0, 0},{7300000, 7400000, 0, 1, 0},{7100000, 7200000, 0, 1, 2},
	{7100000, 7200000, 1, 0, 0},
	{6100000, 6200000, 2, 0, 0},{7100000, 7200000, 2, 0, 2},
	{4100000, 4200000, 3, 0, 0},{6100000, 6200000, 3, 0, 4},{2100000, 2200000, 3, 1, 0},
	{4100000, 4200000, 3, 1, 4},
	{2100000, 2200000, 4, 0, 0},
	{2100000, 2200000, 5, 0, 0},{2500000, 2600000, 5, 0, 4},{2500000, 2600000, 5, 1, 0}},
	[]any{
		miThreadCheck{mitSegment, "spur2", latlongType{7300000, 7400000},
			latlongType{7100000, 7200000}, 0, 1, []latlongRefProto{
			{7300000, 7400000, 0, 0, 0},{7300000, 7400000, 1, 0, 0},
			{7100000, 7200000, 1, 2, 0}}, []any{
			miThreadCheck{mitCircle, "$13", latlongType{7300000, 7400000},
				latlongType{7300000, 7400000}, 0, 0, []latlongRefProto{
				{7300000, 7400000, 0, 0, 0}}, []any{7300000, 7400000}},
			miThreadCheck{mitPath, "ps2", latlongType{7300000, 7400000},
				latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
				{7300000, 7400000, 0, 0, 0},{7100000, 7200000, 2, 0, 0}},
				[]any{7300000, 7400000, 7100000, 7200000}},
		}},
		miThreadCheck{mitMarker, "wp2", latlongType{7100000, 7200000},
			latlongType{7100000, 7200000}, 0, 0, []latlongRefProto{
			{7100000, 7200000, 0, 0, 0}}, []any{7100000, 7200000}},
		miThreadCheck{mitSegment, "s2:1", latlongType{7100000, 7200000},
			latlongType{6100000, 6200000}, 0, 0, []latlongRefProto{
			{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 0, 2, 0}}, []any{
			miThreadCheck{mitPath, "three:1", latlongType{6100000, 6200000},
				latlongType{7100000, 7200000}, 0, 2, []latlongRefProto{
				{6100000, 6200000, 0, 0, 0},{7100000, 7200000, 2, 0, 0}},
				[]any{6100000, 6200000, 7100000, 7200000}},
		}},
		miThreadCheck{mitSegment, "s1:1", latlongType{6100000, 6200000},
			latlongType{2100000, 2200000}, 0, 1, []latlongRefProto{
			{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 0, 4, 0},
			{2100000, 2200000, 1, 0, 0},{4100000, 4200000, 1, 4, 0}}, []any{
			miThreadCheck{mitPath, "two", latlongType{4100000, 4200000},
				latlongType{6100000, 6200000}, 0, 4, []latlongRefProto{
				{4100000, 4200000, 0, 0, 0},{6100000, 6200000, 4, 0, 0}},
				[]any{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			miThreadCheck{mitPath, "one:1", latlongType{2100000, 2200000},
				latlongType{4100000, 4200000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{4100000, 4200000, 4, 0, 0}},
				[]any{2100000, 2200000, 3100000, 3200000, 4100000, 4200000}},
		}},
		miThreadCheck{mitPoint, "wp1", latlongType{2100000, 2200000},
			latlongType{2100000, 2200000}, 0, 0, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0}}, []any{2100000, 2200000}},
		miThreadCheck{mitSegment, "spur1", latlongType{2100000, 2200000},
			latlongType{2500000, 2600000}, 0, 1, []latlongRefProto{
			{2100000, 2200000, 0, 0, 0},{2500000, 2600000, 0, 4, 0},
			{2500000, 2600000, 1, 0, 0}}, []any{
			miThreadCheck{mitPath, "ps1", latlongType{2100000, 2200000},
				latlongType{2500000, 2600000}, 0, 4, []latlongRefProto{
				{2100000, 2200000, 0, 0, 0},{2500000, 2600000, 4, 0, 0}},
				[]any{2100000, 2200000, 2300000, 2400000, 2500000, 2600000}},
			miThreadCheck{mitMarker, "$17", latlongType{2500000, 2600000},
				latlongType{2500000, 2600000}, 0, 0, []latlongRefProto{
				{2500000, 2600000, 0, 0, 0}}, []any{2500000, 2600000}},
		}},
	}})
}

