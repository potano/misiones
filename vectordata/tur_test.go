// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

// Unit tests for recovery actions for crosspoint errors in the threading mechanism


//
// Handling of multiple intersection points between neighbors
//

func Test_findPathIntersectionsAmbig_1f2f(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// Both paths run forward
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
		        30000021, 81000041,
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
		"testfile:1: path path2 connects with path path1 at multiple points\n" +
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1")
}

func Test_findPathIntersectionsAmbig_1f2f_resolve1(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// Both paths run forward
	// Resolve by adding leading and intermediate points to resolve to first crossing point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePointForTests(seg, "point2", 30000015, 81000032),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
		        30000021, 81000041,
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
		{"point2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1f2f_resolve2(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// Both paths run forward
	// Resolve by adding intermediate and trailing points to resolve to second crossing point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
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
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1f2r(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// First path runs forward, second is reversed
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
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: path path2 connects with path path1 at multiple points\n" +
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersectionsAmbig_1f2r_resolve1(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// First path runs forward, second is reversed
	// Resolve by adding leading and intermediate points to resolve to first crossing point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePointForTests(seg, "point2", 30000015, 81000032),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000021, 81000041,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"point2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1f2r_resolve2(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// First path runs forward, second is reversed
	// Resolve by adding intermediate and trailing points to resolve to second crossing point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000021, 81000041,
			30000015, 81000032),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1r2f(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// First path reversed, second runs forward
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
		        30000021, 81000041,
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
		"testfile:1: path path2 connects with path path1 at multiple points\n" +
		"testfile:1: cannot determine free endpoint of path path1 under segment seg1")
}

func Test_findPathIntersectionsAmbig_1r2f_resolve1(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// First path reversed, second runs forward
	// Resolve by adding intermediate and trailing points to resolve to first crossing point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
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
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1r2f_resolve2(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// First path reversed, second runs forward
	// Resolve by adding leading and intermediate points to resolve to second crossing point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePointForTests(seg, "point2", 30000015, 81000032),
		mis.makePathForTests(seg, "path2",
			30000015, 81000032,
		        30000021, 81000041,
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
		{"point2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1r2r(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// Both paths reversed
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
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: path path2 connects with path path1 at multiple points\n" +
		"testfile:1: cannot determine free endpoint of path path2 under segment seg1")
}

func Test_findPathIntersectionsAmbig_1r2r_resolve1(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// Both paths reversed
	// Resolve by adding intermediate and trailing points to resolve to first crosspoint
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePointForTests(seg, "point1", 30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000021, 81000041,
			30000015, 81000032),
		mis.makePointForTests(seg, "point2", 30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"point1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 4, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
		{"point2", []latlongRefProto{{30000051, 81000040, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}

func Test_findPathIntersectionsAmbig_1r2r_resolve2(T *testing.T) {
	// Segment with two paths which share two points causing ambiguity in resolution
	// Both paths reversed
	// Resolve by adding leading and intermediate points to resolve to second crosspoint
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000015, 81000032,
			30000010, 81000020,
			30000000, 81000000),
		mis.makePointForTests(seg, "point2", 30000015, 81000032),
		mis.makePathForTests(seg, "path2",
			30000051, 81000040,
			30000030, 81000043,
		        30000021, 81000041,
			30000015, 81000032),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000000, 81000000, 6, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 2, 0, 0}}},
		{"point2", []latlongRefProto{{30000015, 81000032, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 0, 0, 0}}},
		{"path2", []latlongRefProto{{30000015, 81000032, 6, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 0, 0, 0}}},
	})
	mis.checkDeferredErrors("")
}



//
// Handling of lack of intersection points between neighbors
//

func Test_findPathIntersectionsNoMatch(T *testing.T) {
	// Segment with two discontinuous paths
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "path2",
			30000025, 81000042,
		        30000028, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000025, 81000042, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: path path1 does not connect with segment seg1\n" +
		"testfile:1: path path2 does not connect with segment seg1")
}

func Test_findPathIntersectionsThirdNoMatch1(T *testing.T) {
	// Segment with two continuous paths and a third, disjoint, path
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
			30000022, 81000043,
			30000023, 81000045),
		mis.makePathForTests(seg, "path3",
			30000025, 81000042,
		        30000028, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000023, 81000045, 4, 0, 0}}},
		{"path3", []latlongRefProto{{30000025, 81000042, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: path path3 does not connect with segment seg1")
}

func Test_findPathIntersectionsThirdNoMatch2(T *testing.T) {
	// Segment with two continuous paths preceded by a path that joins neither
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032),
		mis.makePathForTests(seg, "path2",
			30000021, 81000041,
			30000022, 81000043,
			30000023, 81000045,
			30000025, 81000042),
		mis.makePathForTests(seg, "path3",
			30000025, 81000042,
		        30000028, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000015, 81000032, 4, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000025, 81000042, 6, 0, 0}}},
		{"path3", []latlongRefProto{{30000025, 81000042, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: path path1 does not connect with segment seg1\n" +
		"testfile:1: path path2 does not connect with segment seg1")
}

func Test_findPathIntersectionsDiscontinuousPoint1(T *testing.T) {
	// Segment with two continuous paths and a third, disjoint, point
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
			30000022, 81000043,
			30000023, 81000045),
		mis.makePointForTests(seg, "point1", 30000025, 81000042),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000023, 81000045, 4, 0, 0}}},
		{"point1", []latlongRefProto{{30000025, 81000042, 0, 0, 0}},
			[]latlongRefProto{{30000025, 81000042, 0, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: point point1 does not connect with segment seg1")
}

func Test_findPathIntersectionsDiscontinuousPoint2(T *testing.T) {
	// Segment with two continuous paths preceded by a point that joins neither
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePointForTests(seg, "point1", 30000000, 81000000),
		mis.makePathForTests(seg, "path1",
			30000021, 81000041,
			30000022, 81000043,
			30000023, 81000045,
			30000025, 81000042),
		mis.makePathForTests(seg, "path2",
			30000025, 81000042,
		        30000028, 81000041,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"point1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000000, 81000000, 0, 0, 0}}},
		{"path1", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000025, 81000042, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000025, 81000042, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: point point1 does not connect with segment seg1\n" +
		"testfile:1: path path1 does not connect with segment seg1")
}



//
// Handling of glancing single-point intersections to side paths
// These violate the requirement that each non-point component declared to be part of
// a segment or route must contribute at least two points to the item being threaded.
//

func Test_findPathIntersectionsGlancingSidePathAtEnd(T *testing.T) {
	// Segment with two continuous paths with a side path joining at one end
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "side",
			30000021, 81000041,
			30000022, 81000043,
			30000023, 81000045,
			30000025, 81000042),
		mis.makePathForTests(seg, "path2",
			30000021, 81000041,
		        30000028, 81000042,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"side", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000025, 81000042, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: path path2 does not connect with segment seg1")
}

func Test_findPathIntersectionsIntersectingSidePath(T *testing.T) {
	// Segment with two continuous paths crossing a side path at one point
	mis := newMapItemSynthesizer(T)
	seg := mis.makeSegmentForTests(mis.rootItem, "seg1")
	seg.children = []mapItemType{
		mis.makePathForTests(seg, "path1",
			30000000, 81000000,
			30000010, 81000020,
			30000015, 81000032,
			30000021, 81000041),
		mis.makePathForTests(seg, "side",
			30000018, 81000046,
			30000021, 81000041,
			30000023, 81000045,
			30000025, 81000042),
		mis.makePathForTests(seg, "path2",
			30000021, 81000041,
		        30000028, 81000042,
			30000030, 81000043,
			30000051, 81000040),
	}

	mis.resolveReferences_and_setAllPathCrosspoints()
	markedChildren := mis.thread_up_to_markComponentIntersections(seg)
	checkPendingChildren(T, markedChildren, []pendingChildProto{
		{"path1", []latlongRefProto{{30000000, 81000000, 0, 0, 0}},
			[]latlongRefProto{{30000021, 81000041, 6, 0, 0}}},
		{"side", []latlongRefProto{{30000021, 81000041, 2, 0, 0}},
			[]latlongRefProto{{30000025, 81000042, 6, 0, 0}}},
		{"path2", []latlongRefProto{{30000021, 81000041, 0, 0, 0}},
			[]latlongRefProto{{30000051, 81000040, 6, 0, 0}}},
	})
	mis.checkDeferredErrors(
		"testfile:1: cannot determine free endpoint of path side under segment seg1\n" +
		"testfile:1: path path2 does not connect with segment seg1")
}

