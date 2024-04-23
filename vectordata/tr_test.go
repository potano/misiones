// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Path-threading tests into routes


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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
}


func Test_basicReformedRouteSlice(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
}


func Test_basicReformedRouteSliceLiteralWaypoint1(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
}


func Test_basicReformedRouteSliceLiteralWaypoint2(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
}


func Test_basicReformedRouteSliceLiteralBothWaypoints(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{3100000, 3200000, 2100000, 2200000, []gsPath{
			{"one", false, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
}


func Test_basicReformedRouteSliceFlipRoute(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 4100000, 4200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 3100000, 3200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 7100000, 7200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
	})
}


func Test_sliceTwoSegmentReformedRoute(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 7100000, 7200000, []gsPath{
			{"three:1", true, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 10100000, 10200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four", true, locationPairs{8100000, 8200000, 9100000, 9200000,
				10100000, 10200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 9100000, 9200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four", true, locationPairs{8100000, 8200000, 9100000, 9200000}}}},
	})
}


func Test_sliceTwoSegmentReformedRouteBothTwoPaths(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 10100000, 10200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four", true, locationPairs{8100000, 8200000, 9100000, 9200000,
				10100000, 10200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 9100000, 9200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four:1", true, locationPairs{8100000, 8200000, 9100000, 9200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{7100000, 7200000, 6100000, 6200000, []gsPath{
			{"three", false, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
		{6100000, 6200000, 2100000, 2200000, []gsPath{
			{"two", false, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			{"one", false, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}}}},
	})
}


func Test_sliceTwoSegmentReformedRouteReversed(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{7100000, 7200000, 6100000, 6200000, []gsPath{
			{"three:1", false, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
		{6100000, 6200000, 2100000, 2200000, []gsPath{
			{"two", false, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			{"one:1", false, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 10100000, 10200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four", true, locationPairs{8100000, 8200000, 9100000, 9200000,
				10100000, 10200000}}}},
		{10100000, 10200000, 11100000, 11200000, []gsPath{
			{"five", true, locationPairs{10100000, 10200000, 11100000, 11200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 9100000, 9200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four", true, locationPairs{8100000, 8200000, 9100000, 9200000}}}},
	})
}


func Test_sliceThreeSegmentReformedRouteToTwoSegments(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 10100000, 10200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four", true, locationPairs{8100000, 8200000, 9100000, 9200000,
				10100000, 10200000}}}},
		{10100000, 10200000, 11100000, 11200000, []gsPath{
			{"five", true, locationPairs{10100000, 10200000, 11100000, 11200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 9100000, 9200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}},
			{"four:1", true, locationPairs{8100000, 8200000, 9100000, 9200000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2500000, 2600000, 2100000, 2200000, []gsPath{
			{"ps1", true, locationPairs{2500000, 2600000, 2300000, 2400000,
				2100000, 2200000}}}},
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 7100000, 7200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
		{7100000, 7200000, 7300000, 7400000, []gsPath{
			{"ps2", true, locationPairs{7100000, 7200000, 7300000, 7400000}}}},
	})
}


func Test_sliceTwoSegmentReformedRoutePlusAdjoiningSegments(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{2500000, 2600000, 2100000, 2200000, []gsPath{
			{"ps1", true, locationPairs{2500000, 2600000, 2300000, 2400000,
				2100000, 2200000}}}},
		{2100000, 2200000, 6100000, 6200000, []gsPath{
			{"one:1", true, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 7100000, 7200000, []gsPath{
			{"three:1", true, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
		{7100000, 7200000, 7300000, 7400000, []gsPath{
			{"ps2", true, locationPairs{7100000, 7200000, 7300000, 7400000}}}},
	})
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

	vd := prepareAndParseStringsNoRouteThreading(T, sourceText)
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{7300000, 7400000, 7100000, 7200000, []gsPath{
			{"ps2", true, locationPairs{7300000, 7400000, 7100000, 7200000}}}},
		{7100000, 7200000, 6100000, 6200000, []gsPath{
			{"three", false, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
		{6100000, 6200000, 2100000, 2200000, []gsPath{
			{"two", false, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			{"one", false, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}}}},
		{2100000, 2200000, 2500000, 2600000, []gsPath{
			{"ps1", true, locationPairs{2100000, 2200000, 2300000, 2400000,
				2500000, 2600000}}}},
	})
}


func Test_sliceTwoSegmentReformedRoutePlusAdjoiningSegmentsFlipRoute(T *testing.T) {
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
	checkGatheredRouteSegments(T, vd, "road", []gsCheck{
		{1100000, 1200000, 6100000, 6200000, []gsPath{
			{"one", true, locationPairs{1100000, 1200000, 2100000, 2200000,
				3100000, 3200000, 4100000, 4200000}},
			{"two", true, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}}}},
		{6100000, 6200000, 8100000, 8200000, []gsPath{
			{"three", true, locationPairs{6100000, 6200000, 7100000, 7200000,
				8100000, 8200000}}}},
	})
	checkGatheredRouteSegments(T, vd, "partRoad", []gsCheck{
		{7300000, 7400000, 7100000, 7200000, []gsPath{
			{"ps2", true, locationPairs{7300000, 7400000, 7100000, 7200000}}}},
		{7100000, 7200000, 6100000, 6200000, []gsPath{
			{"three:1", false, locationPairs{6100000, 6200000, 7100000, 7200000}}}},
		{6100000, 6200000, 2100000, 2200000, []gsPath{
			{"two", false, locationPairs{4100000, 4200000, 5100000, 5200000,
				6100000, 6200000}},
			{"one:1", false, locationPairs{2100000, 2200000, 3100000, 3200000,
				4100000, 4200000}}}},
		{2100000, 2200000, 2500000, 2600000, []gsPath{
			{"ps1", true, locationPairs{2100000, 2200000, 2300000, 2400000,
				2500000, 2600000}}}},
	})
}

