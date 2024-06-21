// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
        "potano.misiones/sexp"
)



type mapRouteOrSegmentType struct {
	mapItemCore
	popup *mapPopupType
	style *mapStyleType
	attestation *mapAttestationType
	children []mapItemType
	startPoint, endPoint latlongType
	crossings latlongRefs
}

func newMapRoute(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mr := &mapRouteOrSegmentType{}
	mr.source = source
	name, err := doc.registerMapItem(mr, listName)
	mr.name = name
	mr.itemType = mitRoute
	return mr, err
}

func newMapSegment(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mr := &mapRouteOrSegmentType{}
	mr.source = source
	name, err := doc.registerMapItem(mr, listName)
	mr.name = name
	mr.itemType = mitSegment
	return mr, err
}

func (mr *mapRouteOrSegmentType) clone(vd *VectorData, parent mapItemType) *mapRouteOrSegmentType {
	newMR := &mapRouteOrSegmentType{
		popup: mr.popup,
		style: mr.style,
		attestation: mr.attestation,
	}
	newMR.source = mr.source
	newMR.name = registerSplitName(vd, newMR, mr.Name())
	newMR.itemType = mr.itemType
	newMR.referrers = []string{mr.Name()}
	return newMR
}

func (mr *mapRouteOrSegmentType) setPopup(popup *mapPopupType) {
	mr.popup = popup
}

func (mr *mapRouteOrSegmentType) setStyle(style *mapStyleType) {
	mr.style = style
}

func (mr *mapRouteOrSegmentType) setAttestation(attestation *mapAttestationType) {
	mr.attestation = attestation
}

func (mr *mapRouteOrSegmentType) addFeature(feature mapItemType) {
	mr.children = append(mr.children, feature)
}

func (mr *mapRouteOrSegmentType) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return mr.style, mr.attestation
}


func (mr *mapRouteOrSegmentType) routeComponents() []mapItemType {
	return mr.children
}

func (mr *mapRouteOrSegmentType) getCrosspoints() latlongRefs {
	return mr.crossings
}

func (mr *mapRouteOrSegmentType) endpointsAndOffsets() (latlongType, latlongType,
		locationIndexType, locationIndexType) {
	off1, off2 := locationIndexType(0), locationIndexType(len(mr.children) - 1)
	return mr.startPoint, mr.endPoint, off1, off2
}

func (mr *mapRouteOrSegmentType) oppositeEndpoint(startPoint latlongType,
		) (latlongType, locationIndexType, locationIndexType) {
	off1, off2 := locationIndexType(0), locationIndexType(len(mr.children) - 1)
	pt1, pt2 := mr.startPoint, mr.endPoint
	if pt1.samePoint(startPoint) {
		off1, off2 = off2, off1
		pt1, pt2 = pt2, pt1
	}
	return pt1, off1, off2
}

func (mr *mapRouteOrSegmentType) isPoint() bool {
	return false
}

func (mr *mapRouteOrSegmentType) resolveReferenceToLocation(ref latlongRef) *map_locationType {
	children := mr.children
	for _, i := range ref.indices {
		child := children[i]
		if child, is := child.(*map_locationType); is {
			return child
		}
		children = child.(*mapRouteOrSegmentType).children
	}
	return nil
}

func (mr *mapRouteOrSegmentType) setEndpointsAndChildren(startPoint, endPoint latlongType,
		children []mapItemType) {
	mr.startPoint = startPoint
	mr.endPoint = endPoint
	mr.children = children
	var crosspoints latlongRefs
	for childX, child := range children {
		for _, ref := range child.(threadableMapItemType).getCrosspoints() {
			ref = ref.cloneAndPushLevel(locationIndexType(childX))
			crosspoints = append(crosspoints, ref)
		}
	}
	mr.crossings = crosspoints
}


func registerSplitName(vd *VectorData, item mapItemType, baseName string) string {
	if len(baseName) == 0 {
		name, _ := vd.registerMapItem(item, "")
		return name
	}
	for i := 1; i < 100; i++ {
		name := fmt.Sprintf("%s:%d", baseName, i)
		if _, err := vd.registerMapItem(item, name); err == nil {
			return name
		}
	}
	return baseName
}

