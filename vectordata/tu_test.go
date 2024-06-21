// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"strings"
	"testing"
	"potano.misiones/sexp"
)

// Helpers for unit tests of threading mechanism

type miThreadCheck struct {
	itemType int
	name string
	startPoint, endPoint latlongType
	startOffset, endOffset locationIndexType
	crossings []latlongRefProto
	children []any
}

type latlongRefProto struct {
	lat, long locAngleType
	index0, index1, index2 locationIndexType
}

type pendingChildProto struct {
	name string
	startRefs, endRefs []latlongRefProto
}

func checkThreadableMapItem(T *testing.T, have mapItemType, want miThreadCheck) {
	T.Helper()
	if have.ItemType() != want.itemType {
		T.Fatalf("wanted item type %s, got %s", typeMapToName[want.itemType],
			have.ItemTypeString())
	}
	name := have.Name()
	if name != want.name {
		T.Fatalf("wanted item name %s, got %s", want.name, name)
	}
	ti, is := have.(threadableMapItemType)
	if !is {
		T.Fatalf("%s is not a threadable item type", name)
	}
	startPoint, endPoint, startOffset, endOffset := ti.endpointsAndOffsets()
	if !want.startPoint.samePoint(startPoint) {
		T.Fatalf("%s: wanted start point %d, got %d", name, want.startPoint, startPoint)
	}
	if !want.endPoint.samePoint(endPoint) {
		T.Fatalf("%s: wanted end point %d, got %d", name, want.endPoint, endPoint)
	}
	if want.startOffset != startOffset {
		T.Fatalf("%s: wanted start offset %d, got %d", name, want.startOffset, startOffset)
	}
	if want.endOffset != endOffset {
		T.Fatalf("%s: wanted end offset %d, got %d", name, want.endOffset, endOffset)
	}
	checkCrossings(T, name, ti.getCrosspoints(), want.crossings)
	switch tItem := have.(type) {
	case *map_locationType:
		if len(tItem.location) != len(want.children) {
			T.Fatalf("%s: want %d points, have %d", name, len(want.children) >> 1,
				len(tItem.location) >> 1)
		}
		for i, a := range want.children {
			if tItem.location[i] != locAngleType(a.(int)) {
				T.Fatalf("%s point %d: want angle %d, got %d", name, i >> 1,
					a, tItem.location[i])
			}
		}
	case *mapRouteOrSegmentType:
		if len(tItem.children) != len(want.children) {
			T.Fatalf("%s: want %d children, got %d", name, len(want.children),
				len(tItem.children))
		}
		for i, a := range want.children {
			checkThreadableMapItem(T, tItem.children[i], a.(miThreadCheck))
		}
	}
}

func checkCrossings(T *testing.T, name string, have latlongRefs, want []latlongRefProto) {
	T.Helper()
	if len(have) != len(want) {
		T.Fatalf("%s: want %d cross points, got %d", name, len(want), len(have))
	}
	for i, w := range want {
		ref := have[i]
		if ref.point.lat != w.lat || ref.point.long != w.long ||
				ref.indices[0] != w.index0 || ref.indices[1] != w.index1 ||
				ref.indices[2] != w.index2 {
			T.Fatalf("%s crosspoint %d: want (%d %d) [%d,%d,%d], got (%d %d) " +
				"[%d,%d,%d]", name, i, w.lat, w.long, w.index0, w.index1,
				w.index2, ref.point.lat, ref.point.long, ref.indices[0],
				ref.indices[1], ref.indices[2])
		}
	}
}

func checkPendingChildren(T *testing.T, pending []pendingChildInfo, want []pendingChildProto) {
	T.Helper()
	if len(pending) != len(want) {
		T.Fatalf("expecting %d pending children, got %d", len(want), len(pending))
	}
	for i, w := range want {
		childInfo := pending[i]
		if childInfo.child.Name() != w.name {
			T.Fatalf("pending child %d: expected name %s, got %s", i, w.name,
				childInfo.child.Name())
		}
		checkCrossings(T, w.name + " startRefs", childInfo.startRefs, w.startRefs)
		checkCrossings(T, w.name + " endRefs", childInfo.endRefs, w.endRefs)
	}
}

type pickedItemProto struct {
	itemType int
	itemName string
	startLat, startLong, endLat, endLong locAngleType
	startOffset, endOffset locationIndexType
	shortened bool
	children []pickedItemProto
}

func checkPickedItem(T *testing.T, have pickedItem, want pickedItemProto) {
	T.Helper()
	if have.item.ItemType() != want.itemType {
		T.Fatalf("picked item %s is a %s (%d); wanted %s (%d)", have.item.Name(),
			have.item.ItemTypeString(), have.item.ItemType(),
			typeMapToName[want.itemType], want.itemType)
	}
	itemTypeDesc := have.item.ItemTypeString()
	itemName := have.item.Name()
	if have.item.Name() != want.itemName {
		T.Fatalf("wanted picked %s named %s, got %s", itemTypeDesc, want.itemName, itemName)
	}
	wantStartPoint := latlongType{want.startLat, want.startLong}
	wantEndPoint := latlongType{want.endLat, want.endLong}
	if !wantStartPoint.samePoint(have.startPoint) {
		T.Fatalf("picked %s %s: wanted start point %d, got %d", itemTypeDesc, itemName,
			wantStartPoint, have.startPoint)
	}
	if !wantEndPoint.samePoint(have.endPoint) {
		T.Fatalf("picked %s %s: wanted end point %d, got %d", itemTypeDesc, itemName,
			wantEndPoint, have.endPoint)
	}
	if have.startOffset != want.startOffset {
		T.Fatalf("picked %s %s: wanted start offset %d, got %d", itemTypeDesc, itemName,
			want.startOffset, have.startOffset)
	}
	if have.endOffset != want.endOffset {
		T.Fatalf("picked %s %s: wanted end offset %d, got %d", itemTypeDesc, itemName,
			want.endOffset, have.endOffset)
	}
	if have.shortened != want.shortened {
		T.Fatalf("picked %s %s: wanted shortened flag %t, got %t", itemTypeDesc, itemName,
		want.shortened, have.shortened)
	}
	if len(have.children) != len(want.children) {
		T.Fatalf("picked %s %s: wanted %d children, have %d", itemTypeDesc, itemName,
		len(want.children), len(have.children))
	}
	for childX, child := range want.children {
		checkPickedItem(T, have.children[childX], child)
	}
}

func latlongRefFromProto(proto latlongRefProto) latlongRef {
	return latlongRef{latlongType{proto.lat, proto.long},
		[3]locationIndexType{proto.index0, proto.index1, proto.index2}}
}

func pendingChildFromProto(vd *VectorData, proto pendingChildProto) pendingChildInfo {
	startRefs := make(latlongRefs, len(proto.startRefs))
	endRefs := make(latlongRefs, len(proto.endRefs))
	for i, ref := range proto.startRefs {
		startRefs[i] = latlongRefFromProto(ref)
	}
	for i, ref := range proto.endRefs {
		endRefs[i] = latlongRefFromProto(ref)
	}
	return pendingChildInfo{
		child: vd.mapItems[proto.name].(threadableMapItemType),
		startRefs: startRefs,
		endRefs: endRefs,
	}
}

func pendingChildrenFromProto(vd *VectorData, proto []pendingChildProto) []pendingChildInfo {
	pending := make([]pendingChildInfo, len(proto))
	for i, p := range proto {
		pending[i] = pendingChildFromProto(vd, p)
	}
	return pending
}



// Helper functions to synthesize map items without having to parse them

type mapItemSynthesizer struct {
	vd *VectorData
	source sexp.ValueSource
	rootItem mapItemType
	T *testing.T
}

func newMapItemSynthesizer(T *testing.T) mapItemSynthesizer {
	l, err := sexp.Parse("testfile", strings.NewReader("(a)"))
	if err != nil {
		T.Fatalf("error setting up parser: %s", err)
	}
	vd := NewVectorData()
	rootItem := &mapItemCore{}
	name, err := vd.registerMapItem(rootItem, "")
	if err != nil {
		T.Fatalf("failed to create root item: %s", err)
	}
	rootItem.name = name
	return mapItemSynthesizer{vd, l.ValueSource.Source(), rootItem, T}
}

func (mis mapItemSynthesizer) makePathForTests(parent mapItemType, name string,
		pairs ...locAngleType) mapItemType {
	return mis.makeLocationForTests(parent, "path", name, pairs)
}

func (mis mapItemSynthesizer) makePointForTests(parent mapItemType, name string,
		pair ...locAngleType) mapItemType {
	return mis.makeLocationForTests(parent, "point", name, pair)
}

func (mis mapItemSynthesizer) makeLocationForTests(parent mapItemType, itemType, name string,
		pairs locationPairs) mapItemType {
	ml, err := newMap_location(mis.vd, parent, itemType, name, mis.source)
	if err != nil {
		panic("makeLocationForTests returned error: " + err.Error())
	}
	mis.fixReferrer(ml, parent)
	ml.(*map_locationType).appendPoints(pairs)
	return ml
}

func (mis mapItemSynthesizer) makeSegmentForTests(parent mapItemType,
		name string) *mapRouteOrSegmentType {
	ms, err := newMapSegment(mis.vd, parent, "segment", name, mis.source)
	if err != nil {
		panic("makeSegmentForTests returned error: " + err.Error())
	}
	mis.fixReferrer(ms, parent)
	return ms.(*mapRouteOrSegmentType)
}

func (mis mapItemSynthesizer) makeRouteForTests(parent mapItemType,
		name string) *mapRouteOrSegmentType {
	ms, err := newMapRoute(mis.vd, parent, "route", name, mis.source)
	if err != nil {
		panic("makeRouteForTests returned error: " + err.Error())
	}
	mis.fixReferrer(ms, parent)
	return ms.(*mapRouteOrSegmentType)
}

func (mis mapItemSynthesizer) fixReferrer(item mapItemType, parent mapItemType) {
	if parent == nil {
		parent = mis.rootItem
	}
	referrers := []string{parent.Name()}
	switch mi := item.(type) {
	case *map_locationType:
		mi.referrers = referrers
	case *mapRouteOrSegmentType:
		mi.referrers = referrers
	case *mapFeatureType:
		mi.referrers = referrers
	}
}

func (mis mapItemSynthesizer) resolveReferences_and_setAllPathCrosspoints() {
	mis.T.Helper()
	err := mis.vd.ResolveReferences()
	if err != nil {
		mis.T.Fatalf("ResolveReferences: %s", err)
	}
	allCrosspoints := mis.vd.crossingFinder.getAllCrosspoints()
	for _, name := range mis.vd.inDependencyOrder {
		obj := mis.vd.mapItems[name]
		if loc, is := obj.(*map_locationType); is {
			loc.setCrossings(allCrosspoints[loc.locIndex])
		}
	}
}

func (mis mapItemSynthesizer) thread_up_to_markComponentIntersections(item *mapRouteOrSegmentType,
		) []pendingChildInfo {
	mis.T.Helper()
	threadableChildren, err := gatherThreadedItemList(item)
	if err != nil {
		mis.T.Fatalf("gatherThreadedItemList: %s", err)
	}
	markedChildren, err := mis.vd.markComponentIntersections(item, threadableChildren)
	if err != nil {
		mis.T.Fatalf("markComponentIntersections: %s", err)
	}
	return markedChildren
}

func (mis mapItemSynthesizer) checkDeferredErrors(want string) {
	checkDeferredErrors(mis.T, mis.vd, want)
}

func getDeferredErrorsAsString(vd *VectorData) string {
	errs := vd.DeferredErrors()
	if len(errs) == 0 {
		return ""
	}
	strs := make([]string, len(errs))
	for i, e := range errs {
		strs[i] = e.Error()
	}
	return dbg.Join(strs, "\n")
}

func checkDeferredErrors(T *testing.T, vd *VectorData, want string) {
	T.Helper()
	errs := getDeferredErrorsAsString(vd)
	if errs != want {
		if len(want) > 0 {
			T.Fatalf("expected deferred error(s) %s, got %s", want, errs)
		} else {
			T.Fatalf("deferred error(s) %s", errs)
		}
	}
}

