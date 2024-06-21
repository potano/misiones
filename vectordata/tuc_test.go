// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Unit tests for crosspoint discovery and marking in the threading mechanism



//
// Tests of creation of basic threadable map items (paths, segments, and routes)
// without any further processing
//

func Test_basicMakePath(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	path := mis.makePathForTests(nil, "path1",
		30000000, 81000000,
		30000010, 81000020,
		30000015, 81000032,
		30000021, 81000041)

	checkThreadableMapItem(T, path, miThreadCheck{mitPath, "path1",
		latlongType{30000000, 81000000}, latlongType{30000021, 81000041}, 0, 6, nil,
		[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
		30000021, 81000041}})
}

func Test_basicMakeSegment(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(nil, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
		        30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	checkThreadableMapItem(T, seg,
	miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitPath, "path1",
			latlongType{30000000, 81000000}, latlongType{30000021, 81000041}, 0, 6, nil,
			[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
			30000021, 81000041}},
		miThreadCheck{mitPath, "path2",
			latlongType{30000021, 81000041}, latlongType{30000051, 81000040}, 0, 4, nil,
			[]any{30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
	}})
}

func Test_basicMakeRoute(T *testing.T) {
	mis := newMapItemSynthesizer(T)
	route := mis.makeRouteForTests(nil, "route1")
	seg1 := mis.makeSegmentForTests(route, "seg1")
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
	seg2 := mis.makeSegmentForTests(route, "seg2")
	seg2.children = []mapItemType{
		mis.makePathForTests(seg2, "path3",
			30000051, 81000040,
			30000084, 81000035,
			30000105, 81000030,
			30000126, 81000028,
			30000162, 81000031),
		mis.makePathForTests(seg2, "path4",
			30000162, 81000031,
			30000229, 81000068,
			30000298, 81000104,
			30000341, 81000265),
	}
	route.children = []mapItemType{seg1, seg2}

	checkThreadableMapItem(T, route,
	miThreadCheck{mitRoute, "route1", latlongType{}, latlongType{}, 0, 1, nil, []any{
		miThreadCheck{mitSegment, "seg1", latlongType{}, latlongType{}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path1",
				latlongType{30000000, 81000000}, latlongType{30000021, 81000041},
				0, 6, nil,
				[]any{30000000, 81000000, 30000010, 81000020, 30000015, 81000032,
				30000021, 81000041}},
			miThreadCheck{mitPath, "path2",
				latlongType{30000021, 81000041}, latlongType{30000051, 81000040},
				0, 4, nil,
				[]any{30000021, 81000041, 30000030, 81000043, 30000051, 81000040}},
		}},
		miThreadCheck{mitSegment, "seg2", latlongType{}, latlongType{}, 0, 1, nil, []any{
			miThreadCheck{mitPath, "path3",
				latlongType{30000051, 81000040}, latlongType{30000162, 81000031},
				0, 8, nil,
				[]any{30000051, 81000040, 30000084, 81000035, 30000105, 81000030,
				30000126, 81000028, 30000162, 81000031}},
			miThreadCheck{mitPath, "path4",
				latlongType{30000162, 81000031}, latlongType{30000341, 81000265},
				0, 6, nil,
				[]any{30000162, 81000031, 30000229, 81000068, 30000298, 81000104,
				30000341, 81000265}},
		}},
	}})
}


//
// Tests of the determination of intersection points of pairs of paths within a segment
// Viz: these are tests of markComponentIntersections() of the threader component
// These require that paths have already had their candidate crossing points set, so we
// let the path-crossing finder run to compute these.
//

func Test_findPathIntersections1_1f2f(T *testing.T) {
	// Two paths crossing at their endpoints, both running forward
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
		        30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 4, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections1_1f2r(T *testing.T) {
	// Two paths crossing at their endpoints, first forward, second reversed
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000021, 81000041),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections1_1r2f(T *testing.T) {
	// Two paths crossing at their endpoints, first reversed, second forward
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
		        30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 4, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections1_1r2r(T *testing.T) {
	// Two paths crossing at their endpoints, both reversed
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}


func Test_findPathIntersections2_1f2f_ambig(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// Both paths run forward
	// Sets up ambiguity in choice of free endpoint in first path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
		        30000022, 81000035,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1")
}

func Test_findPathIntersections2_1f2f_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// Both paths run forward
	// Resolves ambiguity by adding leading point to segment to select default endpoint
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
		        30000022, 81000035,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1f2f_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// Both paths run forward
	// Resolves ambiguity by adding leading point to segment to select alternate endpoint
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
		        30000022, 81000035,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000021, 81000041, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1f2r_ambig(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// First path runs forward; second is reversed
	// Sets up ambiguity in choice of free endpoint in first path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000022, 81000035,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1")
}

func Test_findPathIntersections2_1f2r_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// First path runs forward; second is reversed
	// Resolves ambiguity by adding leading point to segment to confirm default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000022, 81000035,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1f2r_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// First path runs forward; second is reversed
	// Resolves ambiguity by adding leading point to segment to confirm reversed resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000022, 81000035,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000021, 81000041, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1r2f_ambig(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// First path is reversed, second runs forward
	// Sets up ambiguity in choice of free endpoint in first path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
		        30000015, 81000032,
			30000022, 81000035,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1")
}

func Test_findPathIntersections2_1r2f_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// First path is reversed, second runs forward
	// Resolves ambiguity by adding leading point to segment by using default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
		        30000015, 81000032,
			30000022, 81000035,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1r2f_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of first path and endpoint of second path
	// First path is reversed, second runs forward
	// Resolves ambiguity by adding leading point to segment by using reversed resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
		        30000015, 81000032,
			30000022, 81000035,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1r2r_ambig(T *testing.T) {
	// Intersection of paths fals in interior of first path and endpoint of second path
	// Both paths are reversed
	// Sets up ambiguity in choice of free endpoint in first path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000022, 81000035,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1")
}

func Test_findPathIntersections2_1r2r_resolve1(T *testing.T) {
	// Intersection of paths fals in interior of first path and endpoint of second path
	// Both paths are reversed
	// Resolves ambiguity by adding leading point to segment to select default endpoint
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000022, 81000035,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections2_1r2r_resolve2(T *testing.T) {
	// Intersection of paths fals in interior of first path and endpoint of second path
	// Both paths are reversed
	// Resolves ambiguity by adding leading point to segment to select alternate endpoint
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000022, 81000035,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}


func Test_findPathIntersections3_1f2f_ambig(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// Both paths run forward
	// Sets up ambiguity in choice of free endpoint in second path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000017, 81000033,
		        30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections3_1f2f_resolve1(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// Both paths run forward
	// Resolves ambiguity by adding trailing point to segment to confirm default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000017, 81000033,
		        30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1f2f_resolve2(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// Both paths run forward
	// Resolves ambiguity by adding trailing point to segment to confirm alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000017, 81000033,
		        30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000017, 81000033),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000033, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1f2r_ambig(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// First path runs forward; second is reversed
	// Sets up ambiguity in choice of free endpoint in second path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041,
			30000017, 81000033),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections3_1f2r_resolve1(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// First path runs forward; second is reversed
	// Resolves ambiguity by adding trailing point to segment to confirm default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041,
			30000017, 81000033),
		mis.makePointForTests(seg, "point2", 30000017, 81000033),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000033, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1f2r_resolve2(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// First path runs forward; second is reversed
	// Resolves ambiguity by adding trailing point to segment to confirm alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041,
			30000017, 81000033),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1r2f_ambig(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// First path is reversed; second runs forward
	// Sets up ambiguity in choice of free endpoint in second path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000017, 81000033,
			30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections3_1r2f_resolve1(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// First path is reversed; second runs forward
	// Resolves ambiguity by adding trailing point to segment to confirm default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000017, 81000033,
			30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1r2f_resolve2(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// First path is reversed; second runs forward
	// Resolves ambiguity by adding trailing point to segment to confirm alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000017, 81000033,
			30000021, 81000041,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000017, 81000033),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000033, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1r2r_ambig(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// Both paths are reversed
	// Sets up ambiguity in choice of free endpoint in second path
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041,
			30000017, 81000033),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections3_1r2r_resolve1(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// Both paths are reversed
	// Resolves ambiguity by adding trailing point to segment to confirm default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041,
			30000017, 81000033),
		mis.makePointForTests(seg, "point2", 30000017, 81000033),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000033, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000033, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections3_1r2r_resolve2(T *testing.T) {
	// Intersection of paths falls endpoint of first path and interior of second path
	// Both paths are reversed
	// Resolves ambiguity by adding trailing point to segment to confirm alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
			30000021, 81000041,
			30000017, 81000033),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1f2f_ambig(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// Both paths run forward
	// Sets up ambiguities in free endpoints of both paths
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000026, 81000037),
		mis.makePathForTests(seg, "path2",
			30000017, 81000028,
		        30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1\n" +
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections4_1f2f_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// Both paths run forward
	// Resolves ambiguities by setting leading and trailing points to match default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000026, 81000037),
		mis.makePathForTests(seg, "path2",
			30000017, 81000028,
		        30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1f2f_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// Both paths run forward
	// Resolves ambiguities by setting leading and trailing points to match alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000026, 81000037),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000026, 81000037),
		mis.makePathForTests(seg, "path2",
			30000017, 81000028,
		        30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000017, 81000028),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000026, 81000037, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000026, 81000037, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 2, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000028, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1f2r_ambig(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// First path runs forward; second is reversed
	// Sets up ambiguities in free endpoints of both paths
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000026, 81000037),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000015, 81000032,
			30000017, 81000028),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1\n" +
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections4_1f2r_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// First path runs forward; second is reversed
	// Resolves ambiguities by setting leading and trailing points to match default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000026, 81000037),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000015, 81000032,
			30000017, 81000028),
		mis.makePointForTests(seg, "point2", 30000017, 81000028),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000028, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1f2r_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// First path runs forward; second is reversed
	// Resolves ambiguities by setting leading and trailing points to match alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000026, 81000037),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000026, 81000037),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000015, 81000032,
			30000017, 81000028),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000026, 81000037, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000026, 81000037, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1r2f_ambig(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// First path reversed; second runs forward
	// Sets up ambiguities in free endpoints of both paths
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000026, 81000037,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000017, 81000028,
		        30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1\n" +
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections4_1r2f_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// First path reversed; second runs forward
	// Resolves ambiguities by adding leading and trailing points to match default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000026, 81000037),
		mis.makePathForTests(seg, "path1",
			30000026, 81000037,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000017, 81000028,
		        30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000026, 81000037, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1r2f_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// First path reversed; second runs forward
	// Resolves ambiguities by adding leading and trailing points to match alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000026, 81000037,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000017, 81000028,
		        30000015, 81000032,
			30000030, 81000043,
			30000051, 81000040),
		mis.makePointForTests(seg, "point2", 30000017, 81000028),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 2, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000028, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1r2r_ambig(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// Both paths reversed
	// Sets up ambiguities in free endpoints of both paths
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000026, 81000037,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000015, 81000032,
			30000017, 81000028),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1\n" +
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersections4_1r2r_resolve1(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// Both paths reversed
	// Resolves ambiguities by adding leading and trailing points to match default resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000026, 81000037),
		mis.makePathForTests(seg, "path1",
			30000026, 81000037,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000015, 81000032,
			30000017, 81000028),
		mis.makePointForTests(seg, "point2", 30000017, 81000028),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000026, 81000037, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000026, 81000037, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 4, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000017, 81000028, 0, 0, 0}},
			[]latlongRefProto{{30000017, 81000028, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersections4_1r2r_resolve2(T *testing.T) {
	// Intersection of paths falls in interior of both paths
	// Both paths reversed
	// Resolves ambiguities by adding leading and trailing points to match alternate resolution
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000026, 81000037,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000015, 81000032,
			30000017, 81000028),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

