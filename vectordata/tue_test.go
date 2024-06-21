// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Unit tests for threading where segments contain waypoints that should be included in routes

func Test_gatherThreadedSegmentWithEndpointIntoRoute(T *testing.T) {
	//Point at end of first segment of route joins with second segment
	//The route should include this point because it is a point, marker, or circle
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route")
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000010,
			30000020, 81000020,
		),
		mis.makePathForTests(seg1, "path2",
			30000020, 81000020,
			30000030, 81000030,
		),
		mis.makePointForTests(seg1, "wp1", 30000030, 81000030),
	}
	seg2 := mis.makeSegmentForTests(route, "seg2",)
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000030, 81000030,
			30000040, 81000040,
			30000050, 81000050,
		),
	}
	route.children = []mapItemType{seg1, seg2}

	//Check segments before threading
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 2, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, nil,
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, nil,
			[]any{30000020, 81000020, 30000030, 81000030}},
		miThreadCheck{mitPoint, "wp1", latlongType{30000030, 81000030},
			latlongType{30000030, 81000030}, 0, 0, nil, []any{30000030, 81000030}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 0, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, nil,
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	//Thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000020, 81000020, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000020, 81000020, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 2, 0, 0}}},
		{"wp1", []latlongRefProto{{30000030, 81000030, 0, 0 ,0}},
			[]latlongRefProto{{30000030, 81000030, 0, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000030, 81000030, 0, 2, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
				false, nil},
			{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
				false, nil},
			{mitPoint, "wp1", 30000030, 81000030, 30000030, 81000030, 0, 0, false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
	latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
		{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0},
		{30000030, 81000030, 2, 0, 0}}, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
			{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
			[]any{30000020, 81000020, 30000030, 81000030}},
		miThreadCheck{mitPoint, "wp1", latlongType{30000030, 81000030},
			latlongType{30000030, 81000030}, 0, 0, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0}}, []any{30000030, 81000030}},
	}})

	//Thread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000030, 81000030, 30000050, 81000050, 0, 0, false, []pickedItemProto{
			{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
				false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000030, 81000030},
	latlongType{30000050, 81000050}, 0, 0, []latlongRefProto{
		{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 0, 4, 0}}, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	//Thread route
	markedChildren = mis.thread_up_to_markComponentIntersections(route)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 2, 0, 0}}},
		{"seg2", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 0, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route",
		30000000, 81000000, 30000050, 81000050, 0, 1, false, []pickedItemProto{
			{mitSegment, "seg1", 30000000, 81000000, 30000030, 81000030, 0, 2, false,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
					false, nil},
				{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
					false, nil},
				{mitPoint, "wp1", 30000030, 81000030, 30000030, 81000030, 0, 0,
					false, nil},
			}},
			{mitSegment, "seg2", 30000030, 81000030, 30000050, 81000050, 0, 0, false,
			[]pickedItemProto{
				{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
					false, nil},
			}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route: %s", err)
	}
	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route", latlongType{30000000, 81000000},
	latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 0, 4},
		{30000020, 81000020, 0, 1, 0},{30000030, 81000030, 0, 1, 2},
		{30000030, 81000030, 0, 2, 0},
		{30000030, 81000030, 1, 0, 0},{30000050, 81000050, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
			latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
			{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0},
			{30000030, 81000030, 2, 0, 0}}, []any{
			miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
				latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
			miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
				latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
				{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
				[]any{30000020, 81000020, 30000030, 81000030}},
			miThreadCheck{mitPoint, "wp1", latlongType{30000030, 81000030},
				latlongType{30000030, 81000030}, 0, 0, []latlongRefProto{
				{30000030, 81000030, 0, 0, 0}}, []any{30000030, 81000030}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 0, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 0, 4, 0}}, []any{
			miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
				latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
				{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
				[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
		}},
	}})
}

func Test_gatherThreadedSegmentWithEndpathIntoRoute(T *testing.T) {
	//Scenario similar to the first:  instead of point at the end of the segment is a path
	//which does not join to the second segment.  In this case the route must exclude this
	//path.
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route")
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000010,
			30000020, 81000020,
		),
		mis.makePathForTests(seg1, "path2",
			30000020, 81000020,
			30000030, 81000030,
		),
		mis.makePathForTests(seg1, "extra",
			30000030, 81000030,
			30000032, 81000032,
		),
	}
	seg2 := mis.makeSegmentForTests(route, "seg2",)
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000030, 81000030,
			30000040, 81000040,
			30000050, 81000050,
		),
	}
	route.children = []mapItemType{seg1, seg2}

	//Check segments before threading
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 2, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, nil,
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, nil,
			[]any{30000020, 81000020, 30000030, 81000030}},
		miThreadCheck{mitPath, "extra", latlongType{30000030, 81000030},
			latlongType{30000032, 81000032}, 0, 2, nil,
			[]any{30000030, 81000030, 30000032, 81000032}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 0, nil, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, nil,
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	//Thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000020, 81000020, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000020, 81000020, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 2, 0, 0}}},
		{"extra", []latlongRefProto{{30000030, 81000030, 0, 0 ,0}},
			[]latlongRefProto{{30000032, 81000032, 2, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000032, 81000032, 0, 2, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
				false, nil},
			{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
				false, nil},
			{mitPath, "extra", 30000030, 81000030, 30000032, 81000032, 0, 2,
				false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
	latlongType{30000032, 81000032}, 0, 2, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
		{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0},
		{30000030, 81000030, 2, 0, 0},{30000032, 81000032, 2, 2, 0}}, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
			{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
			[]any{30000020, 81000020, 30000030, 81000030}},
		miThreadCheck{mitPath, "extra", latlongType{30000030, 81000030},
			latlongType{30000032, 81000032}, 0, 2, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000032, 81000032, 2, 0, 0}},
			[]any{30000030, 81000030, 30000032, 81000032}},
	}})

	//Thread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path3", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000030, 81000030, 30000050, 81000050, 0, 0, false, []pickedItemProto{
			{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
				false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000030, 81000030},
	latlongType{30000050, 81000050}, 0, 0, []latlongRefProto{
		{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 0, 4, 0}}, []any{
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	//Thread route
	markedChildren = mis.thread_up_to_markComponentIntersections(route)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 1, 2, 0}}},
		{"seg2", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 0, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route",
		30000000, 81000000, 30000050, 81000050, 0, 1, false, []pickedItemProto{
			{mitSegment, "seg1", 30000000, 81000000, 30000030, 81000030, 0, 1, true,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
					false, nil},
				{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
					false, nil},
			}},
			{mitSegment, "seg2", 30000030, 81000030, 30000050, 81000050, 0, 0, false,
			[]pickedItemProto{
				{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
					false, nil},
			}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route: %s", err)
	}
	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route", latlongType{30000000, 81000000},
	latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 0, 4},
		{30000020, 81000020, 0, 1, 0},{30000030, 81000030, 0, 1, 2},
		{30000030, 81000030, 1, 0, 0},{30000050, 81000050, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "seg1:1", latlongType{30000000, 81000000},
			latlongType{30000030, 81000030}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
			{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0}}, []any{
			miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
				latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
			miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
				latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
				{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
				[]any{30000020, 81000020, 30000030, 81000030}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 0, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 0, 4, 0}}, []any{
			miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
				latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
				{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
				[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
		}},
	}})
}


func Test_gatherThreadedSegmentWithStartpointIntoRoute(T *testing.T) {
	//Point at start of second segment of route joins with first segment
	//The route should include this point because it is a point, marker, or circle
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route")
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000010,
			30000020, 81000020,
		),
		mis.makePathForTests(seg1, "path2",
			30000020, 81000020,
			30000030, 81000030,
		),
	}
	seg2 := mis.makeSegmentForTests(route, "seg2",)
	seg2.children = []mapItemType{
		mis.makePointForTests(seg2, "wp1", 30000030, 81000030),
		mis.makePathForTests(seg2, "path3",
			30000030, 81000030,
			30000040, 81000040,
			30000050, 81000050,
		),
	}
	route.children = []mapItemType{seg1, seg2}

	//Check segments before threading
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, nil,
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, nil,
			[]any{30000020, 81000020, 30000030, 81000030}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPoint, "wp1", latlongType{30000030, 81000030},
			latlongType{30000030, 81000030}, 0, 0, nil, []any{30000030, 81000030}},
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, nil,
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	//Thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000020, 81000020, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000020, 81000020, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 2, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000030, 81000030, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
				false, nil},
			{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
				false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
	latlongType{30000030, 81000030}, 0, 1, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
		{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0}}, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
			{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
			[]any{30000020, 81000020, 30000030, 81000030}},
	}})

	//Thread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"wp1", []latlongRefProto{{30000030, 81000030, 0, 0 ,0}},
			[]latlongRefProto{{30000030, 81000030, 0, 0, 0}}},
		{"path3", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000030, 81000030, 30000050, 81000050, 0, 1, false, []pickedItemProto{
			{mitPoint, "wp1", 30000030, 81000030, 30000030, 81000030, 0, 0, false, nil},
			{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
				false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000030, 81000030},
	latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{{30000030, 81000030, 0, 0, 0},
		{30000030, 81000030, 1, 0, 0},{30000050, 81000050, 1, 4, 0}}, []any{
		miThreadCheck{mitPoint, "wp1", latlongType{30000030, 81000030},
			latlongType{30000030, 81000030}, 0, 0, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0}}, []any{30000030, 81000030}},
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	//Thread route
	markedChildren = mis.thread_up_to_markComponentIntersections(route)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 1, 2, 0}}},
		{"seg2", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 1, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route",
		30000000, 81000000, 30000050, 81000050, 0, 1, false, []pickedItemProto{
			{mitSegment, "seg1", 30000000, 81000000, 30000030, 81000030, 0, 1, false,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
					false, nil},
				{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
					false, nil},
			}},
			{mitSegment, "seg2", 30000030, 81000030, 30000050, 81000050, 0, 1, false,
			[]pickedItemProto{
				{mitPoint, "wp1", 30000030, 81000030, 30000030, 81000030, 0, 0,
					false, nil},
				{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
					false, nil},
			}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route: %s", err)
	}
	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route", latlongType{30000000, 81000000},
	latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 0, 4},
		{30000020, 81000020, 0, 1, 0},{30000030, 81000030, 0, 1, 2},
		{30000030, 81000030, 1, 0, 0},
		{30000030, 81000030, 1, 1, 0},{30000050, 81000050, 1, 1, 4}}, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
			latlongType{30000030, 81000030}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
			{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0}}, []any{
			miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
				latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
			miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
				latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
				{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
				[]any{30000020, 81000020, 30000030, 81000030}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},
			{30000030, 81000030, 1, 0, 0},{30000050, 81000050, 1, 4, 0}}, []any{
			miThreadCheck{mitPoint, "wp1", latlongType{30000030, 81000030},
				latlongType{30000030, 81000030}, 0, 0, []latlongRefProto{
				{30000030, 81000030, 0, 0, 0}}, []any{30000030, 81000030}},
			miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
				latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
				{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
				[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
		}},
	}})
}

func Test_gatherThreadedSegmentWithStartpathIntoRoute(T *testing.T) {
	//Scenario similar to the preceding:  instead of point at the start of the segment is a path
	//which does not join to the first segment.  In this case the route must exclude this
	//path.
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route")
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			30000000, 81000000,
			30000010, 81000010,
			30000020, 81000020,
		),
		mis.makePathForTests(seg1, "path2",
			30000020, 81000020,
			30000030, 81000030,
		),
	}
	seg2 := mis.makeSegmentForTests(route, "seg2",)
	seg2.children = []mapItemType{
		mis.makePathForTests(seg1, "extra",
			30000030, 81000030,
			30000032, 81000032,
		),
		mis.makePathForTests(seg2, "path3",
			30000030, 81000030,
			30000040, 81000040,
			30000050, 81000050,
		),
	}
	route.children = []mapItemType{seg1, seg2}

	//Check segments before threading
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, nil,
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, nil,
			[]any{30000020, 81000020, 30000030, 81000030}},
	}})
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "extra", latlongType{30000030, 81000030},
			latlongType{30000032, 81000032}, 0, 2, nil,
			[]any{30000030, 81000030, 30000032, 81000032}},
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, nil,
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	mis.resolveReferences_and_setAllPathCrosspoints()

	//Thread seg1
	markedChildren := mis.thread_up_to_markComponentIntersections(seg1)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000020, 81000020, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000020, 81000020, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 2, 0, 0}}},
	})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		30000000, 81000000, 30000030, 81000030, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
				false, nil},
			{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
				false, nil},
	}})
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}
	checkThreadableMapItem(T, seg1,
	miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
	latlongType{30000030, 81000030}, 0, 1, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
		{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0}}, []any{
		miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
			latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
			[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
		miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
			latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
			{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
			[]any{30000020, 81000020, 30000030, 81000030}},
	}})

	//Thread seg2
	markedChildren = mis.thread_up_to_markComponentIntersections(seg2)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"extra", []latlongRefProto{{30000032, 81000032, 2, 0 ,0}},
			[]latlongRefProto{{30000030, 81000030, 0, 0, 0}}},
		{"path3", []latlongRefProto{{30000030, 81000030, 0, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 4, 0, 0}}},
	})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg2",
		30000032, 81000032, 30000050, 81000050, 0, 1, false, []pickedItemProto{
			{mitPath, "extra", 30000032, 81000032, 30000030, 81000030, 2, 0,
				false, nil},
			{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
				false, nil},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}
	checkThreadableMapItem(T, seg2,
	miThreadCheck{mitSegment, "seg2", latlongType{30000032, 81000032},
	latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{
		{30000030, 81000030, 0, 0, 0},{30000032, 81000032, 0, 2, 0},
		{30000030, 81000030, 1, 0, 0},{30000050, 81000050, 1, 4, 0}}, []any{
		miThreadCheck{mitPath, "extra", latlongType{30000030, 81000030},
			latlongType{30000032, 81000032}, 0, 2, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000032, 81000032, 2, 0, 0}},
			[]any{30000030, 81000030, 30000032, 81000032}},
		miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
			[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
	}})

	//Thread route
	markedChildren = mis.thread_up_to_markComponentIntersections(route)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"seg1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000030, 81000030, 1, 2, 0}}},
		{"seg2", []latlongRefProto{{30000030, 81000030, 1, 0, 0}},
			[]latlongRefProto{{30000050, 81000050, 1, 4, 0}}},
	})
	pickedChildren = pickThreadedItems(route, markedChildren)
	checkPickedItem(T, pickedChildren, pickedItemProto{mitRoute, "route",
		30000000, 81000000, 30000050, 81000050, 0, 1, false, []pickedItemProto{
			{mitSegment, "seg1", 30000000, 81000000, 30000030, 81000030, 0, 1, false,
			[]pickedItemProto{
				{mitPath, "path1", 30000000, 81000000, 30000020, 81000020, 0, 4,
					false, nil},
				{mitPath, "path2", 30000020, 81000020, 30000030, 81000030, 0, 2,
					false, nil},
			}},
			{mitSegment, "seg2", 30000030, 81000030, 30000050, 81000050, 0, 0, true,
			[]pickedItemProto{
				{mitPath, "path3", 30000030, 81000030, 30000050, 81000050, 0, 4,
					false, nil},
			}},
	}})
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route: %s", err)
	}
	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route", latlongType{30000000, 81000000},
	latlongType{30000050, 81000050}, 0, 1, []latlongRefProto{
		{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 0, 4},
		{30000020, 81000020, 0, 1, 0},{30000030, 81000030, 0, 1, 2},
		{30000030, 81000030, 1, 0, 0},{30000050, 81000050, 1, 0, 4}}, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{30000000, 81000000},
			latlongType{30000030, 81000030}, 0, 1, []latlongRefProto{
			{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 0, 4, 0},
			{30000020, 81000020, 1, 0, 0},{30000030, 81000030, 1, 2, 0}}, []any{
			miThreadCheck{mitPath, "path1", latlongType{30000000, 81000000},
				latlongType{30000020, 81000020}, 0, 4, []latlongRefProto{
				{30000000, 81000000, 0, 0, 0},{30000020, 81000020, 4, 0, 0}},
				[]any{30000000, 81000000, 30000010, 81000010, 30000020, 81000020}},
			miThreadCheck{mitPath, "path2", latlongType{30000020, 81000020},
				latlongType{30000030, 81000030}, 0, 2, []latlongRefProto{
				{30000020, 81000020, 0, 0, 0},{30000030, 81000030, 2, 0, 0}},
				[]any{30000020, 81000020, 30000030, 81000030}},
		}},
		miThreadCheck{mitSegment, "seg2:1", latlongType{30000030, 81000030},
			latlongType{30000050, 81000050}, 0, 0, []latlongRefProto{
			{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 0, 4, 0}}, []any{
			miThreadCheck{mitPath, "path3", latlongType{30000030, 81000030},
				latlongType{30000050, 81000050}, 0, 4, []latlongRefProto{
				{30000030, 81000030, 0, 0, 0},{30000050, 81000050, 4, 0, 0}},
				[]any{30000030, 81000030, 30000040, 81000040, 30000050, 81000050}},
		}},
	}})
}

