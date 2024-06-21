// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Unit tests for special cases in the formation of threaded routes
// Calls each unit in turn noting expected values at each step

func Test_doglegRouteFromIntersectingSegments1(T *testing.T) {
	//Route follows first segment from beginning then turns at intersection with seg2
	//Intersection falls in the interior of paths in both segments
	//This two-segment route requires a point at either end to resolve ambiguity
	mis := newMapItemSynthesizer(T)
	route1 := mis.makeRouteForTests(nil, "route1")
	point1 := mis.makePointForTests(route1, "point1", 30000000, 81000000)
	seg1 := mis.makeSegmentForTests(route1, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg1, "path2",
			30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}
	seg2 := mis.makeSegmentForTests(route1, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000002, 81000010,
			30000006, 81000025,
			30000015, 81000032,
			30000019, 81000031),
		mis.makePathForTests(seg2, "path4",
			30000019, 81000031,
			30000024, 81000036,
			30000036, 81000040,
			30000041, 81000045),
	}
	point2 := mis.makePointForTests(route1, "point2", 30000041, 81000045)
	route1.children = []mapItemType{point1, seg1, seg2, point2}

	//Check segments before threading
	checkThreadableMapItem(T, point1,
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0, nil, []any{30000000, 81000000}})
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 6, nil, []any{
			30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2", latlongType{30000021, 81000041},
			latlongType{30000051, 81000040}, 0, 4, nil, []any{
			30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000019, 81000031}, 0, 6, nil, []any{
			30000002, 81000010, 30000006, 81000025, 30000015, 81000032,
			30000019, 81000031}},
		miThreadCheck{mitPath, "path4", latlongType{30000019, 81000031},
			latlongType{30000041, 81000045}, 0, 6, nil, []any{
			30000019, 81000031, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})
	checkThreadableMapItem(T, point2,
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0, nil, []any{30000041, 81000045}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Mark seg1 crosspoints, pick threaded children, and thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 4, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000051, 81000040, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000021, 81000041, 0, 6,
				false, nil},
			{mitPath, "path2", 30000021, 81000041, 30000051, 81000040, 0, 4,
				false, nil}}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	// Do final check of threaded seg1.  Note that now endpoints and crosspoints are set.
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
		latlongType{30000051, 81000040}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 0, 4, 0},
			{30000021, 81000041, 0, 6, 0},
			{30000021, 81000041, 1, 0, 0},{30000051, 81000040, 1, 4, 0}},
			[]any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 6, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 4, 0, 0},
			{30000021, 81000041, 6, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2", latlongType{30000021, 81000041},
			latlongType{30000051, 81000040}, 0, 4, []latlongRefProto{
			{30000021, 81000041, 0, 0, 0},{30000051, 81000040, 4, 0, 0}},
			[]any{30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})

	// Mark seg2 crosspoints, pick threaded children, and tread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000002, 81000010, 0, 0, 0}},
			[]latlongRefProto{{30000019, 81000031, 6, 0, 0}}},
		{"path4", []latlongRefProto{{30000019, 81000031, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 6, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000002, 81000010, 30000041, 81000045, 0, 1, false, []pickedItemProto{
			{mitPath, "path3", 30000002, 81000010, 30000019, 81000031, 0, 6,
				false, nil},
			{mitPath, "path4", 30000019, 81000031, 30000041, 81000045, 0, 6,
				false, nil}}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	// Do final check of threaded seg2.
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000002, 81000010},
		latlongType{30000041, 81000045}, 0, 1, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000015, 81000032, 0, 4, 0},
			{30000019, 81000031, 0, 6, 0},{30000019, 81000031, 1, 0, 0},
			{30000041, 81000045, 1, 6, 0}},
			[]any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000019, 81000031}, 0, 6, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000015, 81000032, 4, 0, 0},
			{30000019, 81000031, 6, 0, 0}},
			[]any{30000002, 81000010, 30000006, 81000025, 30000015, 81000032,
			30000019, 81000031}},
		miThreadCheck{mitPath, "path4", latlongType{30000019, 81000031},
			latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
			{30000019, 81000031, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
			[]any{30000019, 81000031, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})

	// Finally we thread the route
	markedChildren = mis.thread_up_to_markComponentIntersections(route1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 4, 0}}},
		{"seg2", []latlongRefProto{{30000015, 81000032, 0, 4, 0}},
			[]latlongRefProto{{30000041, 81000045, 1, 6, 0}}},
		{"point2", []latlongRefProto{{30000041, 81000045, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(route1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route1",
		30000000, 81000000, 30000041, 81000045, 0, 3, false, []pickedItemProto{
			{mitPoint, "point1", 30000000, 81000000, 30000000, 81000000, 0, 0,
				false, nil},
			{mitSegment, "seg1", 30000000, 81000000, 30000015, 81000032, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000015, 81000032, 0, 4,
					true, nil},
			}},
			{mitSegment, "seg2", 30000015, 81000032, 30000041, 81000045, 0, 1, true,
			[]pickedItemProto{
				{mitPath, "path3", 30000015, 81000032, 30000019, 81000031, 4, 6,
					true, nil},
				{mitPath, "path4", 30000019, 81000031, 30000041, 81000045, 0, 6,
					false, nil},
			}},
			{mitPoint, "point2", 30000041, 81000045, 30000041, 81000045, 0, 0,
					false, nil},
		}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}
	checkThreadableMapItem(T, route1,
	miThreadCheck{mitRoute, "route1", latlongType{30000000, 81000000},
	latlongType{30000041, 81000045}, 0, 3, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},
		{30000000, 81000000, 1, 0, 0},{30000015, 81000032, 1, 0, 4},
		{30000015, 81000032, 2, 0, 0},{30000019, 81000031, 2, 0, 2},
		{30000019, 81000031, 2, 1, 0},{30000041, 81000045, 2, 1, 6},
		{30000041, 81000045, 3, 0, 0}}, []any{
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0,
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}},[]any{30000000, 81000000}},
		miThreadCheck{mitSegment, "seg1:1", latlongType{30000000, 81000000},
			latlongType{30000015, 81000032}, 0, 0, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path1:1", latlongType{30000000, 81000000},
				latlongType{30000015, 81000032}, 0, 4, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 4, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032}},
		}},
		miThreadCheck{mitSegment, "seg2:1", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 1, []latlongRefProto{
			{30000015, 81000032, 0, 0, 0},{30000019, 81000031, 0, 2, 0},
			{30000019, 81000031, 1, 0, 0},{30000041, 81000045, 1, 6, 0}},
			[]any{
			miThreadCheck{mitPath, "path3:1", latlongType{30000015, 81000032},
				latlongType{30000019, 81000031}, 0, 2, []latlongRefProto{
				{30000015, 81000032, 0, 0, 0},{30000019, 81000031, 2, 0, 0}},
				[]any{30000015, 81000032, 30000019, 81000031}},
			miThreadCheck{mitPath, "path4", latlongType{30000019, 81000031},
				latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
				{30000019, 81000031, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
				[]any{30000019, 81000031, 30000024, 81000036, 30000036, 81000040,
				30000041, 81000045}},
		}},
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0,
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}},[]any{30000041, 81000045}},
	}})
	mis.checkDeferredErrors("")
}

func Test_doglegRouteFromIntersectingSegments2(T *testing.T) {
	//Route follows first segment from beginning then turns at intersection with seg2
	//Intersection falls at the junction of two paths paths in the first segment and
	// the interior of a path in the second segment
	//This two-segment route requires a point at either end to resolve ambiguity
	mis := newMapItemSynthesizer(T)
	route1 := mis.makeRouteForTests(nil, "route1")
	point1 := mis.makePointForTests(route1, "point1", 30000000, 81000000)
	seg1 := mis.makeSegmentForTests(route1, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),   // junction point
		mis.makePathForTests(seg1, "path2",
			30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}
	seg2 := mis.makeSegmentForTests(route1, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000002, 81000010,
			30000006, 81000025,
			30000021, 81000041,    // junction point
			30000019, 81000031),
		mis.makePathForTests(seg2, "path4",
			30000019, 81000031,
			30000024, 81000036,
			30000036, 81000040,
			30000041, 81000045),
	}
	point2 := mis.makePointForTests(route1, "point2", 30000041, 81000045)
	route1.children = []mapItemType{point1, seg1, seg2, point2}

	//Check segments before threading
	checkThreadableMapItem(T, point1,
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0, nil, []any{30000000, 81000000}})
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 6, nil, []any{
			30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2", latlongType{30000021, 81000041},
			latlongType{30000051, 81000040}, 0, 4, nil, []any{
			30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000019, 81000031}, 0, 6, nil, []any{
			30000002, 81000010, 30000006, 81000025, 30000021, 81000041,
			30000019, 81000031}},
		miThreadCheck{mitPath, "path4", latlongType{30000019, 81000031},
			latlongType{30000041, 81000045}, 0, 6, nil, []any{
			30000019, 81000031, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})
	checkThreadableMapItem(T, point2,
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0, nil, []any{30000041, 81000045}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Mark seg1 crosspoints, pick threaded children, and thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 4, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000051, 81000040, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000021, 81000041, 0, 6,
				false, nil},
			{mitPath, "path2", 30000021, 81000041, 30000051, 81000040, 0, 4,
				false, nil}}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	// Do final check of threaded seg1.  Note that now endpoints and crosspoints are set.
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
		latlongType{30000051, 81000040}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000021, 81000041, 0, 6, 0},
			{30000021, 81000041, 1, 0, 0},{30000051, 81000040, 1, 4, 0}},
			[]any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 6, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000021, 81000041, 6, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2", latlongType{30000021, 81000041},
			latlongType{30000051, 81000040}, 0, 4, []latlongRefProto{
			{30000021, 81000041, 0, 0, 0},{30000051, 81000040, 4, 0, 0}},
			[]any{30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})

	// Mark seg2 crosspoints, pick threaded children, and tread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000002, 81000010, 0, 0, 0}},
			[]latlongRefProto{{30000019, 81000031, 6, 0, 0}}},
		{"path4", []latlongRefProto{{30000019, 81000031, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 6, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000002, 81000010, 30000041, 81000045, 0, 1, false, []pickedItemProto{
			{mitPath, "path3", 30000002, 81000010, 30000019, 81000031, 0, 6,
				false, nil},
			{mitPath, "path4", 30000019, 81000031, 30000041, 81000045, 0, 6,
				false, nil}}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	// Do final check of threaded seg2.
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000002, 81000010},
		latlongType{30000041, 81000045}, 0, 1, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000021, 81000041, 0, 4, 0},
			{30000019, 81000031, 0, 6, 0},{30000019, 81000031, 1, 0, 0},
			{30000041, 81000045, 1, 6, 0}},
			[]any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000019, 81000031}, 0, 6, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000021, 81000041, 4, 0, 0},
			{30000019, 81000031, 6, 0, 0}},
			[]any{30000002, 81000010, 30000006, 81000025, 30000021, 81000041,
			30000019, 81000031}},
		miThreadCheck{mitPath, "path4", latlongType{30000019, 81000031},
			latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
			{30000019, 81000031, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
			[]any{30000019, 81000031, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})

	// Finally we thread the route
	markedChildren = mis.thread_up_to_markComponentIntersections(route1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 6, 0}}},
		{"seg2", []latlongRefProto{{30000021, 81000041, 0, 4, 0}},
			[]latlongRefProto{{30000041, 81000045, 1, 6, 0}}},
		{"point2", []latlongRefProto{{30000041, 81000045, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(route1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route1",
		30000000, 81000000, 30000041, 81000045, 0, 3, false, []pickedItemProto{
			{mitPoint, "point1", 30000000, 81000000, 30000000, 81000000, 0, 0,
				false, nil},
			{mitSegment, "seg1", 30000000, 81000000, 30000021, 81000041, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000021, 81000041, 0, 6,
					false, nil},
			}},
			{mitSegment, "seg2", 30000021, 81000041, 30000041, 81000045, 0, 1, true,
			[]pickedItemProto{
				{mitPath, "path3", 30000021, 81000041, 30000019, 81000031, 4, 6,
					true, nil},
				{mitPath, "path4", 30000019, 81000031, 30000041, 81000045, 0, 6,
					false, nil},
			}},
			{mitPoint, "point2", 30000041, 81000045, 30000041, 81000045, 0, 0,
					false, nil},
		}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}
	checkThreadableMapItem(T, route1,
	miThreadCheck{mitRoute, "route1", latlongType{30000000, 81000000},
	latlongType{30000041, 81000045}, 0, 3, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},
		{30000000, 81000000, 1, 0, 0},{30000021, 81000041, 1, 0, 6},
		{30000021, 81000041, 2, 0, 0},{30000019, 81000031, 2, 0, 2},
		{30000019, 81000031, 2, 1, 0},{30000041, 81000045, 2, 1, 6},
		{30000041, 81000045, 3, 0, 0}}, []any{
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0,
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}},[]any{30000000, 81000000}},
		miThreadCheck{mitSegment, "seg1:1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 0, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000021, 81000041, 0, 6, 0}},
			[]any{
			miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
				latlongType{30000021, 81000041}, 0, 6, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000021, 81000041, 6, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
				30000021, 81000041}},
		}},
		miThreadCheck{mitSegment, "seg2:1", latlongType{30000021, 81000041},
			latlongType{30000041, 81000045}, 0, 1, []latlongRefProto{
			{30000021, 81000041, 0, 0, 0},{30000019, 81000031, 0, 2, 0},
			{30000019, 81000031, 1, 0, 0},{30000041, 81000045, 1, 6, 0}},
			[]any{
			miThreadCheck{mitPath, "path3:1", latlongType{30000021, 81000041},
				latlongType{30000019, 81000031}, 0, 2, []latlongRefProto{
				{30000021, 81000041, 0, 0, 0},{30000019, 81000031, 2, 0, 0}},
				[]any{30000021, 81000041, 30000019, 81000031}},
			miThreadCheck{mitPath, "path4", latlongType{30000019, 81000031},
				latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
				{30000019, 81000031, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
				[]any{30000019, 81000031, 30000024, 81000036, 30000036, 81000040,
				30000041, 81000045}},
		}},
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0,
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}},[]any{30000041, 81000045}},
	}})
	mis.checkDeferredErrors("")
}

func Test_doglegRouteFromIntersectingSegments3(T *testing.T) {
	//Route follows first segment from beginning then turns at intersection with seg2
	//Intersection falls in the interior of a path in the first segment and at the
	// intersection of two paths in the second
	//This two-segment route requires a point at either end to resolve ambiguity
	mis := newMapItemSynthesizer(T)
	route1 := mis.makeRouteForTests(nil, "route1")
	point1 := mis.makePointForTests(route1, "point1", 30000000, 81000000)
	seg1 := mis.makeSegmentForTests(route1, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,   // junction point
			30000021, 81000041),
		mis.makePathForTests(seg1, "path2",
			30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}
	seg2 := mis.makeSegmentForTests(route1, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000002, 81000010,
			30000006, 81000025,
			30000010, 81000027,
			30000015, 81000032),
		mis.makePathForTests(seg2, "path4",
			30000015, 81000032,    // junction point
			30000024, 81000036,
			30000036, 81000040,
			30000041, 81000045),
	}
	point2 := mis.makePointForTests(route1, "point2", 30000041, 81000045)
	route1.children = []mapItemType{point1, seg1, seg2, point2}

	//Check segments before threading
	checkThreadableMapItem(T, point1,
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0, nil, []any{30000000, 81000000}})
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 6, nil, []any{
			30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2", latlongType{30000021, 81000041},
			latlongType{30000051, 81000040}, 0, 4, nil, []any{
			30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000015, 81000032}, 0, 6, nil, []any{
			30000002, 81000010, 30000006, 81000025, 30000010, 81000027,
			30000015, 81000032}},
		miThreadCheck{mitPath, "path4", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 6, nil, []any{
			30000015, 81000032, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})
	checkThreadableMapItem(T, point2,
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0, nil, []any{30000041, 81000045}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Mark seg1 crosspoints, pick threaded children, and thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 4, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000051, 81000040, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000021, 81000041, 0, 6,
				false, nil},
			{mitPath, "path2", 30000021, 81000041, 30000051, 81000040, 0, 4,
				false, nil}}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	// Do final check of threaded seg1.  Note that now endpoints and crosspoints are set.
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
		latlongType{30000051, 81000040}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 0, 4, 0},
			{30000021, 81000041, 0, 6, 0},
			{30000021, 81000041, 1, 0, 0},{30000051, 81000040, 1, 4, 0}},
			[]any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000021, 81000041}, 0, 6, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 4, 0, 0},
			{30000021, 81000041, 6, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2", latlongType{30000021, 81000041},
			latlongType{30000051, 81000040}, 0, 4, []latlongRefProto{
			{30000021, 81000041, 0, 0, 0},{30000051, 81000040, 4, 0, 0}},
			[]any{30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})

	// Mark seg2 crosspoints, pick threaded children, and tread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000002, 81000010, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 6, 0, 0}}},
		{"path4", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 6, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000002, 81000010, 30000041, 81000045, 0, 1, false, []pickedItemProto{
			{mitPath, "path3", 30000002, 81000010, 30000015, 81000032, 0, 6,
				false, nil},
			{mitPath, "path4", 30000015, 81000032, 30000041, 81000045, 0, 6,
				false, nil}}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	// Do final check of threaded seg2.
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000002, 81000010},
		latlongType{30000041, 81000045}, 0, 1, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000015, 81000032, 0, 6, 0},
			{30000015, 81000032, 1, 0, 0},{30000041, 81000045, 1, 6, 0}},
			[]any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000015, 81000032}, 0, 6, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000015, 81000032, 6, 0, 0}},
			[]any{30000002, 81000010, 30000006, 81000025, 30000010, 81000027,
			30000015, 81000032}},
		miThreadCheck{mitPath, "path4", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
			{30000015, 81000032, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
			[]any{30000015, 81000032, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})

	// Finally we thread the route
	markedChildren = mis.thread_up_to_markComponentIntersections(route1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 4, 0}}},
		{"seg2", []latlongRefProto{{30000015, 81000032, 1, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 1, 6, 0}}},
		{"point2", []latlongRefProto{{30000041, 81000045, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(route1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route1",
		30000000, 81000000, 30000041, 81000045, 0, 3, false, []pickedItemProto{
			{mitPoint, "point1", 30000000, 81000000, 30000000, 81000000, 0, 0,
				false, nil},
			{mitSegment, "seg1", 30000000, 81000000, 30000015, 81000032, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000015, 81000032, 0, 4,
					true, nil},
			}},
			{mitSegment, "seg2", 30000015, 81000032, 30000041, 81000045, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path4", 30000015, 81000032, 30000041, 81000045, 0, 6,
					false, nil},
			}},
			{mitPoint, "point2", 30000041, 81000045, 30000041, 81000045, 0, 0,
					false, nil},
		}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}
	checkThreadableMapItem(T, route1,
	miThreadCheck{mitRoute, "route1", latlongType{30000000, 81000000},
	latlongType{30000041, 81000045}, 0, 3, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},
		{30000000, 81000000, 1, 0, 0},{30000015, 81000032, 1, 0, 4},
		{30000015, 81000032, 2, 0, 0},{30000041, 81000045, 2, 0, 6},
		{30000041, 81000045, 3, 0, 0}}, []any{
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0,
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}},[]any{30000000, 81000000}},
		miThreadCheck{mitSegment, "seg1:1", latlongType{30000000, 81000000},
			latlongType{30000015, 81000032}, 0, 0, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path1:1", latlongType{30000000, 81000000},
				latlongType{30000015, 81000032}, 0, 4, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 4, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032}},
		}},
		miThreadCheck{mitSegment, "seg2:1", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 0, []latlongRefProto{
			{30000015, 81000032, 0, 0, 0},{30000041, 81000045, 0, 6, 0}},
			[]any{
			miThreadCheck{mitPath, "path4", latlongType{30000015, 81000032},
				latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
				{30000015, 81000032, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
				[]any{30000015, 81000032, 30000024, 81000036, 30000036, 81000040,
				30000041, 81000045}},
		}},
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0,
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}},[]any{30000041, 81000045}},
	}})
	mis.checkDeferredErrors("")
}

func Test_doglegRouteFromIntersectingSegments4(T *testing.T) {
	//Route follows first segment from beginning then turns at intersection with seg2
	//Intersection falls in the junction of paths in both segments
	//This two-segment route requires a point at either end to resolve ambiguity
	mis := newMapItemSynthesizer(T)
	route1 := mis.makeRouteForTests(nil, "route1")
	point1 := mis.makePointForTests(route1, "point1", 30000000, 81000000)
	seg1 := mis.makeSegmentForTests(route1, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000012, 81000028,
			30000015, 81000032),  // junction point
		mis.makePathForTests(seg1, "path2",
			30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
	}
	seg2 := mis.makeSegmentForTests(route1, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000002, 81000010,
			30000006, 81000025,
			30000010, 81000027,
			30000015, 81000032),
		mis.makePathForTests(seg2, "path4",
			30000015, 81000032,    // junction point
			30000024, 81000036,
			30000036, 81000040,
			30000041, 81000045),
	}
	point2 := mis.makePointForTests(route1, "point2", 30000041, 81000045)
	route1.children = []mapItemType{point1, seg1, seg2, point2}

	//Check segments before threading
	checkThreadableMapItem(T, point1,
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0, nil, []any{30000000, 81000000}})
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000015, 81000032}, 0, 6, nil, []any{
			30000000, 81000000, 30000010, 81000020, 30000012, 81000028,
			30000015, 81000032}},
		miThreadCheck{mitPath, "path2", latlongType{30000015, 81000032},
			latlongType{30000051, 81000040}, 0, 4, nil, []any{
			30000015, 81000032, 30000030, 81000043, 30000051, 81000040}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000015, 81000032}, 0, 6, nil, []any{
			30000002, 81000010, 30000006, 81000025, 30000010, 81000027,
			30000015, 81000032}},
		miThreadCheck{mitPath, "path4", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 6, nil, []any{
			30000015, 81000032, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})
	checkThreadableMapItem(T, point2,
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0, nil, []any{30000041, 81000045}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Mark seg1 crosspoints, pick threaded children, and thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 4, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000051, 81000040, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000015, 81000032, 0, 6,
				false, nil},
			{mitPath, "path2", 30000015, 81000032, 30000051, 81000040, 0, 4,
				false, nil}}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	// Do final check of threaded seg1.  Note that now endpoints and crosspoints are set.
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
		latlongType{30000051, 81000040}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 0, 6, 0},
			{30000015, 81000032, 1, 0, 0},{30000051, 81000040, 1, 4, 0}},
			[]any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000015, 81000032}, 0, 6, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 6, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000020, 30000012, 81000028,
			30000015, 81000032}},
		miThreadCheck{mitPath, "path2", latlongType{30000015, 81000032},
			latlongType{30000051, 81000040}, 0, 4, []latlongRefProto{
			{30000015, 81000032, 0, 0, 0},{30000051, 81000040, 4, 0, 0}},
			[]any{30000015, 81000032, 30000030, 81000043, 30000051, 81000040}},
	}})

	// Mark seg2 crosspoints, pick threaded children, and tread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000002, 81000010, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 6, 0, 0}}},
		{"path4", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 6, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000002, 81000010, 30000041, 81000045, 0, 1, false, []pickedItemProto{
			{mitPath, "path3", 30000002, 81000010, 30000015, 81000032, 0, 6,
				false, nil},
			{mitPath, "path4", 30000015, 81000032, 30000041, 81000045, 0, 6,
				false, nil}}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	// Do final check of threaded seg2.
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000002, 81000010},
		latlongType{30000041, 81000045}, 0, 1, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000015, 81000032, 0, 6, 0},
			{30000015, 81000032, 1, 0, 0},{30000041, 81000045, 1, 6, 0}},
			[]any{
		miThreadCheck{mitPath, "path3", latlongType{30000002, 81000010},
			latlongType{30000015, 81000032}, 0, 6, []latlongRefProto{
			{30000002, 81000010, 0, 0, 0},{30000015, 81000032, 6, 0, 0}},
			[]any{30000002, 81000010, 30000006, 81000025, 30000010, 81000027,
			30000015, 81000032}},
		miThreadCheck{mitPath, "path4", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
			{30000015, 81000032, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
			[]any{30000015, 81000032, 30000024, 81000036, 30000036, 81000040,
			30000041, 81000045}},
	}})

	// Finally we thread the route
	markedChildren = mis.thread_up_to_markComponentIntersections(route1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 6, 0}}},
		{"seg2", []latlongRefProto{{30000015, 81000032, 1, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 1, 6, 0}}},
		{"point2", []latlongRefProto{{30000041, 81000045, 0, 0, 0}},
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(route1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route1",
		30000000, 81000000, 30000041, 81000045, 0, 3, false, []pickedItemProto{
			{mitPoint, "point1", 30000000, 81000000, 30000000, 81000000, 0, 0,
				false, nil},
			{mitSegment, "seg1", 30000000, 81000000, 30000015, 81000032, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000015, 81000032, 0, 6,
					false, nil},
			}},
			{mitSegment, "seg2", 30000015, 81000032, 30000041, 81000045, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path4", 30000015, 81000032, 30000041, 81000045, 0, 6,
					false, nil},
			}},
			{mitPoint, "point2", 30000041, 81000045, 30000041, 81000045, 0, 0,
					false, nil},
		}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}
	checkThreadableMapItem(T, route1,
	miThreadCheck{mitRoute, "route1", latlongType{30000000, 81000000},
	latlongType{30000041, 81000045}, 0, 3, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},
		{30000000, 81000000, 1, 0, 0},{30000015, 81000032, 1, 0, 6},
		{30000015, 81000032, 2, 0, 0},{30000041, 81000045, 2, 0, 6},
		{30000041, 81000045, 3, 0, 0}}, []any{
		miThreadCheck{mitPoint, "point1", latlongType{30000000, 81000000},
			latlongType{30000000, 81000000}, 0, 0,
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}},[]any{30000000, 81000000}},
		miThreadCheck{mitSegment, "seg1:1", latlongType{30000000, 81000000},
			latlongType{30000015, 81000032}, 0, 0, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 0, 6, 0}},
			[]any{
			miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
				latlongType{30000015, 81000032}, 0, 6, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000015, 81000032, 6, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000020, 30000012, 81000028,
				30000015, 81000032}},
		}},
		miThreadCheck{mitSegment, "seg2:1", latlongType{30000015, 81000032},
			latlongType{30000041, 81000045}, 0, 0, []latlongRefProto{
			{30000015, 81000032, 0, 0, 0},{30000041, 81000045, 0, 6, 0}},
			[]any{
			miThreadCheck{mitPath, "path4", latlongType{30000015, 81000032},
				latlongType{30000041, 81000045}, 0, 6, []latlongRefProto{
				{30000015, 81000032, 0, 0, 0},{30000041, 81000045, 6, 0, 0}},
				[]any{30000015, 81000032, 30000024, 81000036, 30000036, 81000040,
				30000041, 81000045}},
		}},
		miThreadCheck{mitPoint, "point2", latlongType{30000041, 81000045},
			latlongType{30000041, 81000045}, 0, 0,
			[]latlongRefProto{{30000041, 81000045, 0, 0, 0}},[]any{30000041, 81000045}},
	}})
	mis.checkDeferredErrors("")
}



func Test_threadedRouteFromSegmentWithFlippedPathsAndOtherSegment(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	route1 := mis.makeRouteForTests(nil, "route1")
	seg1 := mis.makeSegmentForTests(route1, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000010, 8100032,
			30000007, 8100027,
			30000000, 8100000),
		mis.makePathForTests(seg1, "path2",
			30000020, 8100043,
			30000015, 8100038,
			30000010, 8100032),
		mis.makePathForTests(seg1, "path3",
			30000030, 8100056,
			30000025, 8100049,
			30000020, 8100043),
	}
	seg2 := mis.makeSegmentForTests(route1, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg1, "path4",
			30000030, 8100056,
			30000035, 8100062,
			30000040, 8100068),
	}
	route1.children = []mapItemType{seg1, seg2}

	//Check segments before threading
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 2, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000010, 8100032},
			latlongType{30000000, 8100000}, 0, 4, nil,
			[]any{30000010, 8100032, 30000007, 8100027, 30000000, 8100000}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 8100043},
			latlongType{30000010, 8100032}, 0, 4, nil,
			[]any{30000020, 8100043, 30000015, 8100038, 30000010, 8100032}},
		miThreadCheck{mitPath, "path3", latlongType{30000030, 8100056},
			latlongType{30000020, 8100043}, 0, 4, nil,
			[]any{30000030, 8100056, 30000025, 8100049, 30000020, 8100043}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 0, nil, []any{
		miThreadCheck{mitPath, "path4", latlongType{30000030, 8100056},
			latlongType{30000040, 8100068}, 0, 4, nil,
			[]any{30000030, 8100056, 30000035, 8100062, 30000040, 8100068}},
	}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 8100000, 4, 0, 0}},
			[]latlongRefProto{{30000010, 8100032, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000010, 8100032, 4, 0, 0}},
			[]latlongRefProto{{30000020, 8100043, 0, 0, 0}}},
		{"path3", []latlongRefProto{{30000020, 8100043, 4, 0, 0}},
			[]latlongRefProto{{30000030, 8100056, 0, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 8100000, 30000030, 8100056, 0, 2, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 8100000, 30000010, 8100032, 4, 0, false, nil},
			{mitPath, "path2", 30000010, 8100032, 30000020, 8100043, 4, 0, false, nil},
			{mitPath, "path3", 30000020, 8100043, 30000030, 8100056, 4, 0, false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 8100000},
	latlongType{30000030, 8100056}, 0, 2, []latlongRefProto{{30000010, 8100032, 0, 0, 0},
	{30000000, 8100000, 0, 4, 0},{30000020, 8100043, 1, 0, 0},{30000010, 8100032, 1, 4, 0},
	{30000030, 8100056, 2, 0, 0},{30000020, 8100043, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000010, 8100032},
			latlongType{30000000, 8100000}, 0, 4, []latlongRefProto{
			{30000010, 8100032, 0, 0, 0},{30000000, 8100000, 4, 0, 0}},
			[]any{30000010, 8100032, 30000007, 8100027, 30000000, 8100000}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 8100043},
			latlongType{30000010, 8100032}, 0, 4, []latlongRefProto{
			{30000020, 8100043, 0, 0, 0},{30000010, 8100032, 4, 0, 0}},
			[]any{30000020, 8100043, 30000015, 8100038, 30000010, 8100032}},
		miThreadCheck{mitPath, "path3", latlongType{30000030, 8100056},
			latlongType{30000020, 8100043}, 0, 4, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000020, 8100043, 4, 0, 0}},
			[]any{30000030, 8100056, 30000025, 8100049, 30000020, 8100043}},
	}})

	// Thread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path4", []latlongRefProto{{30000030, 8100056, 0, 0, 0}},
			[]latlongRefProto{{30000040, 8100068, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000030, 8100056, 30000040, 8100068, 0, 0, false, []pickedItemProto{
			{mitPath, "path4", 30000030, 8100056, 30000040, 8100068, 0, 4, false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000030, 8100056},
	latlongType{30000040, 8100068}, 0, 0, []latlongRefProto{
	{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 0, 4, 0}}, []any{
		miThreadCheck{mitPath, "path4", latlongType{30000030, 8100056},
			latlongType{30000040, 8100068}, 0, 4, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 4, 0, 0}},
			[]any{30000030, 8100056, 30000035, 8100062, 30000040, 8100068}},
	}})

	// Thread route1
	markedChildren = mis.thread_up_to_markComponentIntersections(route1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"seg1", []latlongRefProto{{30000000, 8100000, 0, 4, 0}},
			[]latlongRefProto{{30000030, 8100056, 2, 0, 0}}},
		{"seg2", []latlongRefProto{{30000030, 8100056, 0, 0, 0}},
			[]latlongRefProto{{30000040, 8100068, 0, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route1",
	30000000, 8100000, 30000040, 8100068, 0, 1, false, []pickedItemProto{
		{mitSegment, "seg1", 30000000, 8100000, 30000030, 8100056, 0, 2, false,
			[]pickedItemProto{
			{mitPath, "path1", 30000000, 8100000, 30000010, 8100032, 4, 0, false, nil},
			{mitPath, "path2", 30000010, 8100032, 30000020, 8100043, 4, 0, false, nil},
			{mitPath, "path3", 30000020, 8100043, 30000030, 8100056, 4, 0, false, nil},
		}},
		{mitSegment, "seg2", 30000030, 8100056, 30000040, 8100068, 0, 0, false,
			[]pickedItemProto{
			{mitPath, "path4", 30000030, 8100056, 30000040, 8100068, 0, 4, false, nil},
		}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}
	checkThreadableMapItem(T, route1,
	miThreadCheck{mitRoute, "route1", latlongType{30000000, 8100000},
	latlongType{30000040, 8100068}, 0, 1, []latlongRefProto{{30000010, 8100032, 0, 0, 0},
	{30000000, 8100000, 0, 0, 4},{30000020, 8100043, 0, 1, 0},{30000010, 8100032, 0, 1, 4},
	{30000030, 8100056, 0, 2, 0},{30000020, 8100043, 0, 2, 4},{30000030, 8100056, 1, 0, 0},
	{30000040, 8100068, 1, 0, 4}},[]any{
		miThreadCheck{mitSegment, "seg1", latlongType{30000000, 8100000},
			latlongType{30000030, 8100056}, 0, 2, []latlongRefProto{
			{30000010, 8100032, 0, 0, 0},{30000000, 8100000, 0, 4, 0},
			{30000020, 8100043, 1, 0, 0},{30000010, 8100032, 1, 4, 0},
			{30000030, 8100056, 2, 0, 0},{30000020, 8100043, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path1", latlongType{30000010, 8100032},
				latlongType{30000000, 8100000},	0, 4, []latlongRefProto{
				{30000010, 8100032, 0, 0, 0},{30000000, 8100000, 4, 0, 0}},
				[]any{30000010, 8100032, 30000007, 8100027, 30000000, 8100000}},
			miThreadCheck{mitPath, "path2", latlongType{30000020, 8100043},
				latlongType{30000010, 8100032}, 0, 4, []latlongRefProto{
				{30000020, 8100043, 0, 0, 0},{30000010, 8100032, 4, 0, 0}},
				[]any{30000020, 8100043, 30000015, 8100038, 30000010, 8100032}},
			miThreadCheck{mitPath, "path3", latlongType{30000030, 8100056},
				latlongType{30000020, 8100043}, 0, 4, []latlongRefProto{
				{30000030, 8100056, 0, 0, 0},{30000020, 8100043, 4, 0, 0}},
				[]any{30000030, 8100056, 30000025, 8100049, 30000020, 8100043}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{30000030, 8100056},
			latlongType{30000040, 8100068}, 0, 0, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path4", latlongType{30000030, 8100056},
				latlongType{30000040, 8100068}, 0, 4, []latlongRefProto{
				{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 4, 0, 0}},
				[]any{30000030, 8100056, 30000035, 8100062, 30000040, 8100068}},
		}},
	}})
	mis.checkDeferredErrors("")
}

func Test_threadedRouteFromFlippedSegmentWithFlippedPathsAndOtherSegment(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	route1 := mis.makeRouteForTests(nil, "route1")
	seg1 := mis.makeSegmentForTests(route1, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path3",
			30000030, 8100056,
			30000025, 8100049,
			30000020, 8100043),
		mis.makePathForTests(seg1, "path2",
			30000020, 8100043,
			30000015, 8100038,
			30000010, 8100032),
		mis.makePathForTests(seg1, "path1",
			30000010, 8100032,
			30000007, 8100027,
			30000000, 8100000),
	}
	seg2 := mis.makeSegmentForTests(route1, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg1, "path4",
			30000030, 8100056,
			30000035, 8100062,
			30000040, 8100068),
	}
	route1.children = []mapItemType{seg1, seg2}

	//Check segments before threading
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 2, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000030, 8100056},
			latlongType{30000020, 8100043}, 0, 4, nil,
			[]any{30000030, 8100056, 30000025, 8100049, 30000020, 8100043}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 8100043},
			latlongType{30000010, 8100032}, 0, 4, nil,
			[]any{30000020, 8100043, 30000015, 8100038, 30000010, 8100032}},
		miThreadCheck{mitPath, "path1", latlongType{30000010, 8100032},
			latlongType{30000000, 8100000}, 0, 4, nil,
			[]any{30000010, 8100032, 30000007, 8100027, 30000000, 8100000}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 0, nil, []any{
		miThreadCheck{mitPath, "path4", latlongType{30000030, 8100056},
			latlongType{30000040, 8100068}, 0, 4, nil,
			[]any{30000030, 8100056, 30000035, 8100062, 30000040, 8100068}},
	}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000030, 8100056, 0, 0, 0}},
			[]latlongRefProto{{30000020, 8100043, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000020, 8100043, 0, 0, 0}},
			[]latlongRefProto{{30000010, 8100032, 4, 0, 0}}},
		{"path1", []latlongRefProto{{30000010, 8100032, 0, 0, 0}},
			[]latlongRefProto{{30000000, 8100000, 4, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000030, 8100056, 30000000, 8100000, 0, 2, false, []pickedItemProto{
			{mitPath, "path3", 30000030, 8100056, 30000020, 8100043, 0, 4, false, nil},
			{mitPath, "path2", 30000020, 8100043, 30000010, 8100032, 0, 4, false, nil},
			{mitPath, "path1", 30000010, 8100032, 30000000, 8100000, 0, 4, false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000030, 8100056},
	latlongType{30000000, 8100000}, 0, 2, []latlongRefProto{{30000030, 8100056, 0, 0, 0},
	{30000020, 8100043, 0, 4, 0},{30000020, 8100043, 1, 0, 0},{30000010, 8100032, 1, 4, 0},
	{30000010, 8100032, 2, 0, 0},{30000000, 8100000, 2, 4, 0}}, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000030, 8100056},
			latlongType{30000020, 8100043}, 0, 4, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000020, 8100043, 4, 0, 0}},
			[]any{30000030, 8100056, 30000025, 8100049, 30000020, 8100043}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 8100043},
			latlongType{30000010, 8100032}, 0, 4, []latlongRefProto{
			{30000020, 8100043, 0, 0, 0},{30000010, 8100032, 4, 0, 0}},
			[]any{30000020, 8100043, 30000015, 8100038, 30000010, 8100032}},
		miThreadCheck{mitPath, "path1", latlongType{30000010, 8100032},
			latlongType{30000000, 8100000}, 0, 4, []latlongRefProto{
			{30000010, 8100032, 0, 0, 0},{30000000, 8100000, 4, 0, 0}},
			[]any{30000010, 8100032, 30000007, 8100027, 30000000, 8100000}},
	}})

	// Thread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path4", []latlongRefProto{{30000030, 8100056, 0, 0, 0}},
			[]latlongRefProto{{30000040, 8100068, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000030, 8100056, 30000040, 8100068, 0, 0, false, []pickedItemProto{
			{mitPath, "path4", 30000030, 8100056, 30000040, 8100068, 0, 4, false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000030, 8100056},
	latlongType{30000040, 8100068}, 0, 0, []latlongRefProto{
	{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 0, 4, 0}}, []any{
		miThreadCheck{mitPath, "path4", latlongType{30000030, 8100056},
			latlongType{30000040, 8100068}, 0, 4, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 4, 0, 0}},
			[]any{30000030, 8100056, 30000035, 8100062, 30000040, 8100068}},
	}})

	// Thread route1
	markedChildren = mis.thread_up_to_markComponentIntersections(route1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"seg1", []latlongRefProto{{30000000, 8100000, 2, 4, 0}},
			[]latlongRefProto{{30000030, 8100056, 0, 0, 0}}},
		{"seg2", []latlongRefProto{{30000030, 8100056, 0, 0, 0}},
			[]latlongRefProto{{30000040, 8100068, 0, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route1",
	30000000, 8100000, 30000040, 8100068, 0, 1, false, []pickedItemProto{
		{mitSegment, "seg1", 30000000, 8100000, 30000030, 8100056, 0, 2, false,
			[]pickedItemProto{
			{mitPath, "path1", 30000000, 8100000, 30000010, 8100032, 4, 0, false, nil},
			{mitPath, "path2", 30000010, 8100032, 30000020, 8100043, 4, 0, false, nil},
			{mitPath, "path3", 30000020, 8100043, 30000030, 8100056, 4, 0, false, nil},
		}},
		{mitSegment, "seg2", 30000030, 8100056, 30000040, 8100068, 0, 0, false,
			[]pickedItemProto{
			{mitPath, "path4", 30000030, 8100056, 30000040, 8100068, 0, 4, false, nil},
		}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}
	checkThreadableMapItem(T, route1,
	miThreadCheck{mitRoute, "route1", latlongType{30000000, 8100000},
	latlongType{30000040, 8100068}, 0, 1, []latlongRefProto{{30000030, 8100056, 0, 0, 0},
	{30000020, 8100043, 0, 0, 4},{30000020, 8100043, 0, 1, 0},{30000010, 8100032, 0, 1, 4},
	{30000010, 8100032, 0, 2, 0},{30000000, 8100000, 0, 2, 4},{30000030, 8100056, 1, 0, 0},
	{30000040, 8100068, 1, 0, 4}},[]any{
		miThreadCheck{mitSegment, "seg1", latlongType{30000030, 8100056},
			latlongType{30000000, 8100000}, 0, 2, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000020, 8100043, 0, 4, 0},
			{30000020, 8100043, 1, 0, 0},{30000010, 8100032, 1, 4, 0},
			{30000010, 8100032, 2, 0, 0},{30000000, 8100000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path3", latlongType{30000030, 8100056},
				latlongType{30000020, 8100043}, 0, 4, []latlongRefProto{
				{30000030, 8100056, 0, 0, 0},{30000020, 8100043, 4, 0, 0}},
				[]any{30000030, 8100056, 30000025, 8100049, 30000020, 8100043}},
			miThreadCheck{mitPath, "path2", latlongType{30000020, 8100043},
				latlongType{30000010, 8100032}, 0, 4, []latlongRefProto{
				{30000020, 8100043, 0, 0, 0},{30000010, 8100032, 4, 0, 0}},
				[]any{30000020, 8100043, 30000015, 8100038, 30000010, 8100032}},
			miThreadCheck{mitPath, "path1", latlongType{30000010, 8100032},
				latlongType{30000000, 8100000},	0, 4, []latlongRefProto{
				{30000010, 8100032, 0, 0, 0},{30000000, 8100000, 4, 0, 0}},
				[]any{30000010, 8100032, 30000007, 8100027, 30000000, 8100000}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{30000030, 8100056},
			latlongType{30000040, 8100068}, 0, 0, []latlongRefProto{
			{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 0, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path4", latlongType{30000030, 8100056},
				latlongType{30000040, 8100068}, 0, 4, []latlongRefProto{
				{30000030, 8100056, 0, 0, 0},{30000040, 8100068, 4, 0, 0}},
				[]any{30000030, 8100056, 30000035, 8100062, 30000040, 8100068}},
		}},
	}})
	mis.checkDeferredErrors("")
}

func Test_threadedTwoSegmentRouteFirstSegmentFlippedPath(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	route := mis.makeRouteForTests(mis.rootItem, "theRoad")
	seg1 := mis.makeSegmentForTests(route, "roadSeg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path3",
			31003000, 32003000,
			31003001, 32003001,
			31004000, 32004000),
		mis.makePathForTests(seg1, "path2",
			31002000, 32002000,
			31002001, 32002001,
			31003000, 32003000),
		mis.makePathForTests(seg1, "path1",
			31001000, 32001000,
			31001001, 32001001,
			31002000, 32002000),
	}
	seg2 := mis.makeSegmentForTests(route, "roadSeg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path4",
			31004000, 32004000,
			31004001, 32004001,
			31005000, 32005000),
		mis.makePathForTests(seg2, "path5",
			31005000, 32005000,
			31005001, 32005001,
			31006000, 32006000),
		mis.makePathForTests(seg2, "path6",
			31006000, 32006000,
			31006001, 32006001,
			31006002, 32006002),
	}
	route.children = []mapItemType{seg1, seg2}

	mis.resolveReferences_and_setAllPathCrosspoints()

	// Thread roadSeg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{31004000, 32004000, 4, 0, 0}},
			[]latlongRefProto{{31003000, 32003000, 0, 0, 0}}},
		{"path2", []latlongRefProto{{31003000, 32003000, 4, 0, 0}},
			[]latlongRefProto{{31002000, 32002000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{31002000, 32002000, 4, 0, 0}},
			[]latlongRefProto{{31001000, 32001000, 0, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "roadSeg1",
		31004000, 32004000, 31001000, 32001000, 0, 2, false, []pickedItemProto{
			{mitPath, "path3", 31004000, 32004000, 31003000, 32003000, 4, 0,
				false, nil},
			{mitPath, "path2", 31003000, 32003000, 31002000, 32002000, 4, 0,
				false, nil},
			{mitPath, "path1", 31002000, 32002000, 31001000, 32001000, 4, 0,
				false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading roadSeg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "roadSeg1", latlongType{31004000, 32004000},
		latlongType{31001000, 32001000}, 0, 2, []latlongRefProto{
		{31003000, 32003000, 0, 0, 0},{31004000, 32004000, 0, 4, 0},
		{31002000, 32002000, 1, 0, 0},{31003000, 32003000, 1, 4, 0},
		{31001000, 32001000, 2, 0, 0},{31002000, 32002000, 2, 4, 0}},
		[]any{
		miThreadCheck{mitPath, "path3", latlongType{31003000, 32003000},
			latlongType{31004000, 32004000}, 0, 4, []latlongRefProto{
			{31003000, 32003000, 0, 0, 0},{31004000, 32004000, 4, 0, 0}},
			[]any{31003000, 32003000, 31003001, 32003001, 31004000, 32004000}},
		miThreadCheck{mitPath, "path2", latlongType{31002000, 32002000},
			latlongType{31003000, 32003000}, 0, 4, []latlongRefProto{
			{31002000, 32002000, 0, 0, 0},{31003000, 32003000, 4, 0, 0}},
			[]any{31002000, 32002000, 31002001, 32002001, 31003000, 32003000}},
		miThreadCheck{mitPath, "path1", latlongType{31001000, 32001000},
			latlongType{31002000, 32002000}, 0, 4, []latlongRefProto{
			{31001000, 32001000, 0, 0, 0},{31002000, 32002000, 4, 0, 0}},
			[]any{31001000, 32001000, 31001001, 32001001, 31002000, 32002000}},
	}})

	// Thread roadSeg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path4", []latlongRefProto{{31004000, 32004000, 0, 0, 0}},
			[]latlongRefProto{{31005000, 32005000, 4, 0, 0}}},
		{"path5", []latlongRefProto{{31005000, 32005000, 0, 0, 0}},
			[]latlongRefProto{{31006000, 32006000, 4, 0, 0}}},
		{"path6", []latlongRefProto{{31006000, 32006000, 0, 0, 0}},
			[]latlongRefProto{{31006002, 32006002, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "roadSeg2",
		31004000, 32004000, 31006002, 32006002, 0, 2, false, []pickedItemProto{
			{mitPath, "path4", 31004000, 32004000, 31005000, 32005000, 0, 4,
				false, nil},
			{mitPath, "path5", 31005000, 32005000, 31006000, 32006000, 0, 4,
				false, nil},
			{mitPath, "path6", 31006000, 32006000, 31006002, 32006002, 0, 4,
				false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading roadSeg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "roadSeg2", latlongType{31004000, 32004000},
		latlongType{31006002, 32006002}, 0, 2, []latlongRefProto{
		{31004000, 32004000, 0, 0, 0},{31005000, 32005000, 0, 4, 0},
		{31005000, 32005000, 1, 0, 0},{31006000, 32006000, 1, 4, 0},
		{31006000, 32006000, 2, 0, 0},{31006002, 32006002, 2, 4, 0}},
		[]any{
		miThreadCheck{mitPath, "path4", latlongType{31004000, 32004000},
			latlongType{31005000, 32005000}, 0, 4, []latlongRefProto{
			{31004000, 32004000, 0, 0, 0},{31005000, 32005000, 4, 0, 0}},
			[]any{31004000, 32004000, 31004001, 32004001, 31005000, 32005000}},
		miThreadCheck{mitPath, "path5", latlongType{31005000, 32005000},
			latlongType{31006000, 32006000}, 0, 4, []latlongRefProto{
			{31005000, 32005000, 0, 0, 0},{31006000, 32006000, 4, 0, 0}},
			[]any{31005000, 32005000, 31005001, 32005001, 31006000, 32006000}},
		miThreadCheck{mitPath, "path6", latlongType{31006000, 32006000},
			latlongType{31006002, 32006002}, 0, 4, []latlongRefProto{
			{31006000, 32006000, 0, 0, 0},{31006002, 32006002, 4, 0, 0}},
			[]any{31006000, 32006000, 31006001, 32006001, 31006002, 32006002}},
	}})

	// Thread the route
	markedChildren = mis.thread_up_to_markComponentIntersections(route)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"roadSeg1", []latlongRefProto{{31001000, 32001000, 2, 0, 0}},
			[]latlongRefProto{{31004000, 32004000, 0, 4, 0}}},
		{"roadSeg2", []latlongRefProto{{31004000, 32004000, 0, 0, 0}},
			[]latlongRefProto{{31006002, 32006002, 2, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "theRoad",
		31001000, 32001000, 31006002, 32006002, 0, 1, false, []pickedItemProto{
			{mitSegment, "roadSeg1", 31001000, 32001000, 31004000, 32004000, 0, 2,
			false, []pickedItemProto{
				{mitPath, "path1", 31001000, 32001000, 31002000, 32002000, 0, 4,
					false, nil},
				{mitPath, "path2", 31002000, 32002000, 31003000, 32003000, 0, 4,
					false, nil},
				{mitPath, "path3", 31003000, 32003000, 31004000, 32004000, 0, 4,
					false, nil},
			}},
			{mitSegment, "roadSeg2", 31004000, 32004000, 31006002, 32006002, 0, 2,
			false, []pickedItemProto{
				{mitPath, "path4", 31004000, 32004000, 31005000, 32005000, 0, 4,
					false, nil},
				{mitPath, "path5", 31005000, 32005000, 31006000, 32006000, 0, 4,
					false, nil},
				{mitPath, "path6", 31006000, 32006000, 31006002, 32006002, 0, 4,
					false, nil},
			}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading theRoad: %s", err)
	}
	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "theRoad", latlongType{31001000, 32001000},
	latlongType{31006002, 32006002}, 0, 1, []latlongRefProto{
	{31003000, 32003000, 0, 0, 0},{31004000, 32004000, 0, 0, 4},
	{31002000, 32002000, 0, 1, 0},{31003000, 32003000, 0, 1, 4},
	{31001000, 32001000, 0, 2, 0},{31002000, 32002000, 0, 2, 4},
	{31004000, 32004000, 1, 0, 0},{31005000, 32005000, 1, 0, 4},
	{31005000, 32005000, 1, 1, 0},{31006000, 32006000, 1, 1, 4},
	{31006000, 32006000, 1, 2, 0},{31006002, 32006002, 1, 2, 4}},
	[]any{
		miThreadCheck{mitSegment, "roadSeg1", latlongType{31004000, 32004000},
			latlongType{31001000, 32001000}, 0, 2, []latlongRefProto{
			{31003000, 32003000, 0, 0, 0},{31004000, 32004000, 0, 4, 0},
			{31002000, 32002000, 1, 0, 0},{31003000, 32003000, 1, 4, 0},
			{31001000, 32001000, 2, 0, 0},{31002000, 32002000, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path3", latlongType{31003000, 32003000},
				latlongType{31004000, 32004000}, 0, 4, []latlongRefProto{
				{31003000, 32003000, 0, 0, 0},{31004000, 32004000, 4, 0, 0}},
				[]any{31003000, 32003000, 31003001, 32003001, 31004000, 32004000}},
			miThreadCheck{mitPath, "path2", latlongType{31002000, 32002000},
				latlongType{31003000, 32003000}, 0, 4, []latlongRefProto{
				{31002000, 32002000, 0, 0, 0},{31003000, 32003000, 4, 0, 0}},
				[]any{31002000, 32002000, 31002001, 32002001, 31003000, 32003000}},
			miThreadCheck{mitPath, "path1", latlongType{31001000, 32001000},
				latlongType{31002000, 32002000}, 0, 4, []latlongRefProto{
				{31001000, 32001000, 0, 0, 0},{31002000, 32002000, 4, 0, 0}},
				[]any{31001000, 32001000, 31001001, 32001001, 31002000, 32002000}},
		}},
		miThreadCheck{mitSegment, "roadSeg2", latlongType{31004000, 32004000},
			latlongType{31006002, 32006002}, 0, 2, []latlongRefProto{
			{31004000, 32004000, 0, 0, 0},{31005000, 32005000, 0, 4, 0},
			{31005000, 32005000, 1, 0, 0},{31006000, 32006000, 1, 4, 0},
			{31006000, 32006000, 2, 0, 0},{31006002, 32006002, 2, 4, 0}},
			[]any{
			miThreadCheck{mitPath, "path4", latlongType{31004000, 32004000},
				latlongType{31005000, 32005000}, 0, 4, []latlongRefProto{
				{31004000, 32004000, 0, 0, 0},{31005000, 32005000, 4, 0, 0}},
				[]any{31004000, 32004000, 31004001, 32004001, 31005000, 32005000}},
			miThreadCheck{mitPath, "path5", latlongType{31005000, 32005000},
				latlongType{31006000, 32006000}, 0, 4, []latlongRefProto{
				{31005000, 32005000, 0, 0, 0},{31006000, 32006000, 4, 0, 0}},
				[]any{31005000, 32005000, 31005001, 32005001, 31006000, 32006000}},
			miThreadCheck{mitPath, "path6", latlongType{31006000, 32006000},
				latlongType{31006002, 32006002}, 0, 4, []latlongRefProto{
				{31006000, 32006000, 0, 0, 0},{31006002, 32006002, 4, 0, 0}},
				[]any{31006000, 32006000, 31006001, 32006001, 31006002, 32006002}},
		}},
	}})
}

