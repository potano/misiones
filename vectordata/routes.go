// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
        "potano.misiones/sexp"
)



type mapRouteType struct {
	mapItemCore
	popup *mapPopupType
	style *mapStyleType
	attestation *mapAttestationType
	segments []mapItemType
}

func newMapRoute(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mr := &mapRouteType{}
	mr.source = source
	name, err := doc.registerMapItem(mr, listName)
	mr.name = name
	return mr, err
}

func (mr *mapRouteType) ItemType() int {
	return mitRoute
}

func (mr *mapRouteType) setPopup(popup *mapPopupType) {
	mr.popup = popup
}

func (mr *mapRouteType) setStyle(style *mapStyleType) {
	mr.style = style
}

func (mr *mapRouteType) setAttestation(attestation *mapAttestationType) {
	mr.attestation = attestation
}

func (mr *mapRouteType) addFeature(segment mapItemType) {
	mr.segments = append(mr.segments, segment)
}

func (mr *mapRouteType) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return mr.style, mr.attestation
}



type mapSegmentType struct {
	mapItemCore
	popup *mapPopupType
	style *mapStyleType
	attestation *mapAttestationType
	paths []mapItemType
	reformed bool
}

func newMapSegment(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	if parent == nil && len(listName) == 0 {
		return nil, source.Error("no name given for non-embedded segment")
	}
	ms := &mapSegmentType{}
	ms.source = source
	name, err := doc.registerMapItem(ms, listName)
	ms.name = name
	return ms, err
}

func (ms *mapSegmentType) ItemType() int {
	return mitSegment
}

func (ms *mapSegmentType) setPoup(popup *mapPopupType) {
	ms.popup = popup
}

func (ms *mapSegmentType) setStyle(style *mapStyleType) {
	ms.style = style
}

func (ms *mapSegmentType) setAttestation(attestation *mapAttestationType) {
	ms.attestation = attestation
}

func (ms *mapSegmentType) addFeature(path mapItemType) {
	ms.paths = append(ms.paths, path)
}

func (ms *mapSegmentType) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return ms.style, ms.attestation
}





func (mr *mapRouteType) checkAndReformSegments(vd *VectorData) error {
	tSegs, err := mr.threadSegments()
	if err != nil {
		return err
	}
	newSegments := make([]mapItemType, len(tSegs))
	for segX, tSeg := range tSegs {
		if tSeg.splitSegment {
			tSeg = mr.generateSplitSegment(vd, tSeg)
		}
		newSegments[segX] = mr.reformSegment(vd, tSeg)
	}
	mr.segments = newSegments
	return nil
}


func (mr *mapRouteType) reformSegment(vd *VectorData, gathered *gatheredSegment) *mapSegmentType {
	seg := gathered.obj.(*mapSegmentType)
	if seg.reformed {
		seg.referrers = append(seg.referrers, mr.Name())
		return seg
	}
	paths := make([]mapItemType, 0, len(gathered.paths))
	for _, gpath := range gathered.paths {
		if gpath.waypointBefore != nil {
			paths = append(paths, gpath.waypointBefore)
		}
		loc, _, _ := gpath.points()
		path := gpath.path
		if len(loc) != len(path.location) {
			path = &map_locationType{
				itemType: mitPath,
				popup: path.popup,
				html: path.html,
				style: path.style,
				attestation: path.attestation,
				location: loc,
			}
			path.source = gpath.path.Source()
			path.name = registerSplitName(vd, gpath.path)
			path.referrers = []string{seg.Name()}
		}
		paths = append(paths, path)
		if gpath.waypointAfter != nil {
			paths = append(paths, gpath.waypointAfter)
		}
	}
	seg.paths = paths
	seg.reformed = true
	return seg
}


func (mr *mapRouteType) generateSplitSegment(vd *VectorData, tSeg *gatheredSegment,
		) *gatheredSegment {
	seg := tSeg.obj.(*mapSegmentType)
	newSeg := &mapSegmentType{
		popup: seg.popup,
		style: seg.style,
		attestation: seg.attestation,
	}
	newSeg.source = seg.Source()
	newSeg.name = registerSplitName(vd, newSeg)
	newSeg.referrers = []string{mr.Name()}
	tSeg.obj = newSeg
	return tSeg
}


func registerSplitName(vd *VectorData, item mapItemType) string {
	baseName := item.Name()
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

