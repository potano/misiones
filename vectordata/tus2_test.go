// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

// Unit tests for more special cases in the formation of threaded routes

import "testing"


func Test_gatherEndToEndPathIntoThreadedSegment(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			3000000, 8100000,
			3000010, 8100020,
			3000015, 8100032,
			3000021, 8100041),
		mis.makePathForTests(seg, "path2",
			3000021, 8100041,
			3000028, 8100053,
			3000036, 8100048),
	}
	markedChildren := pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
			[]latlongRefProto{{3000021, 8100041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{3000021, 8100041, 0, 0, 0}},
			[]latlongRefProto{{3000036, 8100048, 4, 0, 0}}}})

	pickedChildren := pickThreadedItems(seg, markedChildren)

	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		3000000, 8100000, 3000036, 8100048, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 3000000, 8100000, 3000021, 8100041, 0, 6, false, nil},
			{mitPath, "path2", 3000021, 8100041, 3000036, 8100048, 0, 4, false, nil}}})

	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading: %s", err)
	}

	checkThreadableMapItem(T, seg,
	miThreadCheck{mitSegment, "seg1", latlongType{3000000, 8100000},
		latlongType{3000036, 8100048}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path1", latlongType{3000000, 8100000},
				latlongType{3000021, 8100041}, 0, 6, nil, []any{3000000, 8100000,
				3000010, 8100020, 3000015, 8100032, 3000021, 8100041}},
			miThreadCheck{mitPath, "path2", latlongType{3000021, 8100041},
				latlongType{3000036, 8100048}, 0, 4, nil,
				[]any{3000021, 8100041, 3000028, 8100053, 3000036, 8100048}},
		}})
	mis.checkDeferredErrors("")
}

func Test_gatherPathAndSpurIntoThreadedSegment(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			3000000, 8100000,
			3000010, 8100020,
			3000015, 8100032,
			3000021, 8100041),
		mis.makePathForTests(seg, "path2",
			3000015, 8100032,
			3000028, 8100053,
			3000036, 8100048),
	}
	markedChildren := pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
			[]latlongRefProto{{3000015, 8100032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{3000015, 8100032, 0, 0, 0}},
			[]latlongRefProto{{3000036, 8100048, 4, 0, 0}}}})

	pickedChildren := pickThreadedItems(seg, markedChildren)

	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		3000000, 8100000, 3000036, 8100048, 0, 1, false, []pickedItemProto{
			{mitPath, "path1", 3000000, 8100000, 3000015, 8100032, 0, 4, true, nil},
			{mitPath, "path2", 3000015, 8100032, 3000036, 8100048, 0, 4, false, nil}}})

	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading: %s", err)
	}

	checkThreadableMapItem(T, seg,
	miThreadCheck{mitSegment, "seg1", latlongType{3000000, 8100000},
		latlongType{3000036, 8100048}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path1:1", latlongType{3000000, 8100000},
				latlongType{3000015, 8100032}, 0, 4, nil, []any{3000000, 8100000,
				3000010, 8100020, 3000015, 8100032}},
			miThreadCheck{mitPath, "path2", latlongType{3000015, 8100032},
				latlongType{3000036, 8100048}, 0, 4, nil,
				[]any{3000015, 8100032, 3000028, 8100053, 3000036, 8100048}},
		}})
	mis.checkDeferredErrors("")
}

func Test_gatherPathsBridgedByShortTraverseOfPathIntoThreadedSegment(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			3000000, 8100000,
			3000010, 8100020,
			3000015, 8100032,
			3000021, 8100041),
		mis.makePathForTests(mis.rootItem, "middlePath",
			3000012, 8100066,
			3000018, 8100053,
			3000021, 8100041,
			3000025, 8100044,
			3000038, 8100049),
		mis.makePathForTests(seg, "path2",
			3000025, 8100044,
			3000028, 8100053,
			3000036, 8100048),
	}
	markedChildren := pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
			[]latlongRefProto{{3000021, 8100041, 6, 0, 0}}},
		{"middlePath", []latlongRefProto{{3000021, 8100041, 4, 0, 0}},
			[]latlongRefProto{{3000025, 8100044, 6, 0, 0}}},
		{"path2", []latlongRefProto{{3000025, 8100044, 0, 0, 0}},
			[]latlongRefProto{{3000036, 8100048, 4, 0, 0}}}})

	pickedChildren := pickThreadedItems(seg, markedChildren)

	checkPickedItem(T, pickedChildren, pickedItemProto{mitSegment, "seg1",
		3000000, 8100000, 3000036, 8100048, 0, 2, false, []pickedItemProto{
			{mitPath, "path1", 3000000, 8100000, 3000021, 8100041, 0, 6, false, nil},
			{mitPath, "middlePath", 3000021, 8100041, 3000025, 8100044, 4, 6, true,
				nil},
			{mitPath, "path2", 3000025, 8100044, 3000036, 8100048, 0, 4, false, nil}}})

	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading: %s", err)
	}

	checkThreadableMapItem(T, seg,
	miThreadCheck{mitSegment, "seg1", latlongType{3000000, 8100000},
		latlongType{3000036, 8100048}, 0, 2, nil, []any{
			miThreadCheck{mitPath, "path1", latlongType{3000000, 8100000},
				latlongType{3000021, 8100041}, 0, 6, nil, []any{3000000, 8100000,
				3000010, 8100020, 3000015, 8100032, 3000021, 8100041}},
			miThreadCheck{mitPath, "middlePath:1", latlongType{3000021, 8100041},
				latlongType{3000025, 8100044}, 0, 2, nil,
				[]any{3000021, 8100041, 3000025, 8100044}},
			miThreadCheck{mitPath, "path2", latlongType{3000025, 8100044},
				latlongType{3000036, 8100048}, 0, 4, nil,
				[]any{3000025, 8100044, 3000028, 8100053, 3000036, 8100048}},
		}})
	mis.checkDeferredErrors("")
}


func Test_gatherThreadedSegmentsIntoRoute(T *testing.T) {
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route1")
	//Prepare and thread segment seg1
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			3000000, 8100000,
			3000010, 8100020,
			3000015, 8100032,
			3000021, 8100041),
		mis.makePathForTests(seg1, "path2",
			3000021, 8100041,
			3000028, 8100053,
			3000036, 8100048),
	}
	markedChildren := pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
			[]latlongRefProto{{3000021, 8100041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{3000021, 8100041, 0, 0, 0}},
			[]latlongRefProto{{3000036, 8100048, 4, 0, 0}}}})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}

	//Prepare and thread segment seg2
	seg2 := mis.makeSegmentForTests(route, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			3000036, 8100048,
			3000045, 8100059,
			3000058, 8100063),
		mis.makePathForTests(seg2, "path4",
			3000058, 8100063,
			3000067, 8100062,
			3000073, 8100065,
			3000084, 8100070),
	}
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path3", []latlongRefProto{{3000036, 8100048, 0, 0, 0}},
			[]latlongRefProto{{3000058, 8100063, 4, 0, 0}}},
		{"path4", []latlongRefProto{{3000058, 8100063, 0, 0, 0}},
			[]latlongRefProto{{3000084, 8100070, 6, 0, 0}}}})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}

	route.children = []mapItemType{seg1, seg2}

	// Simulate marking of crossings for route then thread route
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"seg1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
                        []latlongRefProto{{3000036, 8100048, 1, 4, 0}}},
		{"seg2", []latlongRefProto{{3000036, 8100048, 0, 0, 0}},
			[]latlongRefProto{{3000084, 8100070, 1, 6, 0}}}})
	pickedChildren = pickThreadedItems(route, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}

	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route1", latlongType{3000000, 8100000},
	latlongType{3000084, 8100070}, 0, 1, nil, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{3000000, 8100000},
			latlongType{3000036, 8100048}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path1", latlongType{3000000, 8100000},
				latlongType{3000021, 8100041}, 0, 6, nil,
				[]any{3000000, 8100000, 3000010, 8100020, 3000015, 8100032,
				3000021, 8100041}},
			miThreadCheck{mitPath, "path2", latlongType{3000021, 8100041},
				latlongType{3000036, 8100048}, 0, 4, nil,
				[]any{3000021, 8100041, 3000028, 8100053, 3000036, 8100048}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{3000036, 8100048},
		latlongType{3000084, 8100070}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path3", latlongType{3000036, 8100048},
				latlongType{3000058, 8100063}, 0, 4, nil, []any{
				3000036, 8100048, 3000045, 8100059, 3000058, 8100063}},
			miThreadCheck{mitPath, "path4", latlongType{3000058, 8100063},
				latlongType{3000084, 8100070}, 0, 6, nil, []any{
				3000058, 8100063, 3000067, 8100062, 3000073, 8100065,
				3000084, 8100070}},
		}},
	}})
	mis.checkDeferredErrors("")
}

func Test_gatherThreeEndToEndThreadedSegmentsIntoRoute(T *testing.T) {
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route1")
	//Prepare and thread segment seg1
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			3000000, 8100000,
			3000010, 8100020,
			3000015, 8100032,
			3000021, 8100041),
		mis.makePathForTests(seg1, "path2",
			3000021, 8100041,
			3000028, 8100053,
			3000036, 8100048),
	}
	markedChildren := pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
			[]latlongRefProto{{3000021, 8100041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{3000021, 8100041, 0, 0, 0}},
			[]latlongRefProto{{3000036, 8100048, 4, 0, 0}}}})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}

	//Prepare and thread segment seg2
	seg2 := mis.makeSegmentForTests(route, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			3000036, 8100048,
			3000045, 8100059,
			3000058, 8100063),
		mis.makePathForTests(seg2, "path4",
			3000058, 8100063,
			3000067, 8100062,
			3000073, 8100065,
			3000084, 8100070),
	}
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path3", []latlongRefProto{{3000036, 8100048, 0, 0, 0}},
			[]latlongRefProto{{3000058, 8100063, 4, 0, 0}}},
		{"path4", []latlongRefProto{{3000058, 8100063, 0, 0, 0}},
			[]latlongRefProto{{3000084, 8100070, 6, 0, 0}}}})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}

	//Prepare and thread segment seg3
	seg3 := mis.makeSegmentForTests(route, "seg3")
	seg3.children = []mapItemType{
		mis.makePathForTests(seg3, "path5",
			3000084, 8100070,
			3000096, 8100084,
			3000109, 8100075),
		mis.makePathForTests(seg3, "path6",
			3000109, 8100075,
			3000126, 8100084,
			3000141, 8100092,
			3000138, 8100096),
	}
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path5", []latlongRefProto{{3000084, 8100070, 0, 0, 0}},
			[]latlongRefProto{{3000109, 8100075, 4, 0, 0}}},
		{"path6", []latlongRefProto{{3000109, 8100075, 0, 0, 0}},
			[]latlongRefProto{{3000138, 8100096, 6, 0, 0}}}})
	pickedChildren = pickThreadedItems(seg3, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}

	route.children = []mapItemType{seg1, seg2, seg3}

	// Simulate marking of crossings for route then thread route
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"seg1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
                        []latlongRefProto{{3000036, 8100048, 1, 4, 0}}},
		{"seg2", []latlongRefProto{{3000036, 8100048, 0, 0, 0}},
			[]latlongRefProto{{3000084, 8100070, 1, 6, 0}}},
		{"seg3", []latlongRefProto{{3000084, 8100070, 0, 0, 0}},
			[]latlongRefProto{{3000138, 8100096, 1, 6, 0}}}})
	pickedChildren = pickThreadedItems(route, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}

	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route1", latlongType{3000000, 8100000},
	latlongType{3000138, 8100096}, 0, 2, nil, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{3000000, 8100000},
			latlongType{3000036, 8100048}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path1", latlongType{3000000, 8100000},
				latlongType{3000021, 8100041}, 0, 6, nil,
				[]any{3000000, 8100000, 3000010, 8100020, 3000015, 8100032,
				3000021, 8100041}},
			miThreadCheck{mitPath, "path2", latlongType{3000021, 8100041},
				latlongType{3000036, 8100048}, 0, 4, nil,
				[]any{3000021, 8100041, 3000028, 8100053, 3000036, 8100048}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{3000036, 8100048},
		latlongType{3000084, 8100070}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path3", latlongType{3000036, 8100048},
				latlongType{3000058, 8100063}, 0, 4, nil, []any{
				3000036, 8100048, 3000045, 8100059, 3000058, 8100063}},
			miThreadCheck{mitPath, "path4", latlongType{3000058, 8100063},
				latlongType{3000084, 8100070}, 0, 6, nil, []any{
				3000058, 8100063, 3000067, 8100062, 3000073, 8100065,
				3000084, 8100070}},
		}},
		miThreadCheck{mitSegment, "seg3", latlongType{3000084, 8100070},
		latlongType{3000138, 8100096}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path5", latlongType{3000084, 8100070},
				latlongType{3000109, 8100075}, 0, 4, nil, []any{
				3000084, 8100070, 3000096, 8100084, 3000109, 8100075}},
			miThreadCheck{mitPath, "path6", latlongType{3000109, 8100075},
				latlongType{3000138, 8100096}, 0, 6, nil, []any{
				3000109, 8100075, 3000126, 8100084, 3000141, 8100092,
				3000138, 8100096}},
		}},
	}})
	mis.checkDeferredErrors("")
}

func Test_gatherThreeThreadedSegmentsIntoRouteMiddleMiddle(T *testing.T) {
	mis := newMapItemSynthesizer(T)

	route := mis.makeRouteForTests(mis.rootItem, "route1")
	//Prepare and thread segment seg1
	seg1 := mis.makeSegmentForTests(route, "seg1")
	seg1.children = []mapItemType{
		mis.makePathForTests(seg1, "path1",
			3000000, 8100000,
			3000010, 8100020,
			3000015, 8100032,
			3000021, 8100041),
		mis.makePathForTests(seg1, "path2",
			3000021, 8100041,
			3000028, 8100053,
			3000036, 8100048),
	}
	markedChildren := pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
			[]latlongRefProto{{3000021, 8100041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{3000021, 8100041, 0, 0, 0}},
			[]latlongRefProto{{3000036, 8100048, 4, 0, 0}}}})
	pickedChildren := pickThreadedItems(seg1, markedChildren)
	err := mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg1: %s", err)
	}

	//Prepare and thread segment seg2
	seg2 := mis.makeSegmentForTests(route, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			3000029, 8100046,
			3000036, 8100048,
			3000045, 8100059,
			3000058, 8100063),
		mis.makePathForTests(seg2, "path4",
			3000058, 8100063,
			3000067, 8100062,
			3000073, 8100065,
			3000084, 8100070),
	}
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path3", []latlongRefProto{{3000029, 8100046, 0, 0, 0}},
			[]latlongRefProto{{3000058, 8100063, 6, 0, 0}}},
		{"path4", []latlongRefProto{{3000058, 8100063, 0, 0, 0}},
			[]latlongRefProto{{3000084, 8100070, 6, 0, 0}}}})
	pickedChildren = pickThreadedItems(seg2, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}

	//Prepare and thread segment seg3
	seg3 := mis.makeSegmentForTests(route, "seg3")
	seg3.children = []mapItemType{
		mis.makePathForTests(seg3, "path5",
			3000073, 8100065,
			3000096, 8100084,
			3000109, 8100075),
		mis.makePathForTests(seg3, "path6",
			3000109, 8100075,
			3000126, 8100084,
			3000141, 8100092,
			3000138, 8100096),
	}
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"path5", []latlongRefProto{{3000073, 8100065, 0, 0, 0}},
			[]latlongRefProto{{3000109, 8100075, 4, 0, 0}}},
		{"path6", []latlongRefProto{{3000109, 8100075, 0, 0, 0}},
			[]latlongRefProto{{3000138, 8100096, 6, 0, 0}}}})
	pickedChildren = pickThreadedItems(seg3, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading seg2: %s", err)
	}

	route.children = []mapItemType{seg1, seg2, seg3}

	// Simulate marking of crossings for route then thread route
	markedChildren = pendingChildrenFromProto(mis.vd, []pendingChildProto{
		{"seg1", []latlongRefProto{{3000000, 8100000, 0, 0, 0}},
                        []latlongRefProto{{3000036, 8100048, 1, 4, 0}}},
		{"seg2", []latlongRefProto{{3000036, 8100048, 0, 2, 0}},
			[]latlongRefProto{{3000073, 8100065, 1, 4, 0}}},
		{"seg3", []latlongRefProto{{3000073, 8100065, 0, 0, 0}},
			[]latlongRefProto{{3000138, 8100096, 1, 6, 0}}}})
	pickedChildren = pickThreadedItems(route, markedChildren)
	err = mis.vd.finishThreading(pickedChildren)
	if err != nil {
		T.Fatalf("finishThreading route1: %s", err)
	}

	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route1", latlongType{3000000, 8100000},
	latlongType{3000138, 8100096}, 0, 2, nil, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{3000000, 8100000},
			latlongType{3000036, 8100048}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path1", latlongType{3000000, 8100000},
				latlongType{3000021, 8100041}, 0, 6, nil,
				[]any{3000000, 8100000, 3000010, 8100020, 3000015, 8100032,
				3000021, 8100041}},
			miThreadCheck{mitPath, "path2", latlongType{3000021, 8100041},
				latlongType{3000036, 8100048}, 0, 4, nil,
				[]any{3000021, 8100041, 3000028, 8100053, 3000036, 8100048}},
		}},
		miThreadCheck{mitSegment, "seg2:1", latlongType{3000036, 8100048},
		latlongType{3000073, 8100065}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path3:1", latlongType{3000036, 8100048},
				latlongType{3000058, 8100063}, 0, 4, nil, []any{
				3000036, 8100048, 3000045, 8100059, 3000058, 8100063}},
			miThreadCheck{mitPath, "path4:1", latlongType{3000058, 8100063},
				latlongType{3000073, 8100065}, 0, 4, nil, []any{
				3000058, 8100063, 3000067, 8100062, 3000073, 8100065}},
		}},
		miThreadCheck{mitSegment, "seg3", latlongType{3000073, 8100065},
		latlongType{3000138, 8100096}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path5", latlongType{3000073, 8100065},
				latlongType{3000109, 8100075}, 0, 4, nil, []any{
				3000073, 8100065, 3000096, 8100084, 3000109, 8100075}},
			miThreadCheck{mitPath, "path6", latlongType{3000109, 8100075},
				latlongType{3000138, 8100096}, 0, 6, nil, []any{
				3000109, 8100075, 3000126, 8100084, 3000141, 8100092,
				3000138, 8100096}},
		}},
	}})
	mis.checkDeferredErrors("")
}

