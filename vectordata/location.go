// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"strconv"
	"potano.misiones/sexp"
)


type map_locationType struct {
	mapItemCore
	popup *mapPopupType
	html string
	style *mapStyleType
	attestation *mapAttestationType
	radius, radiusType int
	location locationPairs
	vd *VectorData
	prototypePath *map_locationType
	offsetInPrototype, locIndex locationIndexType
	crossings latlongRefs
	isPointType, isRouteComponent bool
}

func newMap_location(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	if parent == nil && len(listName) == 0 {
		return nil, source.Error("no name given for non-embedded %s", listType)
	}
	ml := &map_locationType{locIndex: countRegisteredLocations(doc) + 1, vd: doc}
	ml.source = source
	itemType := nameToTypeMap[listType]
	ml.itemType = itemType
	if itemType == 0 {
		return nil, source.Error("unknown object type '%s'", listType)
	}
	ml.isPointType = itemType == mitPoint || itemType == mitCircle || itemType == mitMarker
	ml.isRouteComponent = itemType == mitPath || ml.isPointType
	name, err := doc.registerMapItem(ml, listName)
	ml.name = name
	return ml, err
}

func (ml *map_locationType) makeSubpath(vd *VectorData, parent mapItemType,
		startOffset, endOffset locationIndexType) *map_locationType {
	if endOffset < startOffset {
		startOffset, endOffset = endOffset, startOffset
	}
	location := ml.location[startOffset:endOffset + 2]
	offsetInPrototype := ml.offsetInPrototype + startOffset
	for _, item := range ml.vd.mapItems {
		if loc, is := item.(*map_locationType); is && loc.locIndex == ml.locIndex &&
				loc.offsetInPrototype == offsetInPrototype &&
				len(loc.location) == len(location) {
			same := true
			for i, angle := range location {
				if loc.location[i] != angle {
					same = false
					break
				}
			}
			if same {
				// Reuse existing item
				loc.referrers = append(loc.referrers, parent.Name())
				return loc
			}
		}
	}
	var newCrossings latlongRefs
	for _, ref := range ml.crossings {
		if ref.indices[0] >= startOffset && ref.indices[0] <= endOffset {
			newRef := ref.clone()
			newRef.indices[0] -= startOffset
			newCrossings = append(newCrossings, newRef)
		}
	}
	newML := &map_locationType{
		location: location,
		vd: ml.vd,
		prototypePath: ml,
		offsetInPrototype: offsetInPrototype,
		locIndex: ml.locIndex,
		crossings: newCrossings,
	}
	newML.source = ml.Source()
	newML.name = registerSplitName(vd, newML, ml.Name())
	newML.itemType = ml.itemType
	newML.referrers = []string{parent.Name()}
	return newML
}

func (ml *map_locationType) setPopup(popup *mapPopupType) {
	ml.popup = popup
}

func (ml *map_locationType) setHtml(html *map_textType) {
	ml.html = html.text
}

func (ml *map_locationType) setStyle(style *mapStyleType) {
	ml.style = style
}

func (ml *map_locationType) setAttestation(attestation *mapAttestationType) {
	ml.attestation = attestation
}

func (ml *map_locationType) setRadius(radius *mapRadiusType) {
	ml.radius = radius.radius
	ml.radiusType = radius.ItemType()
}

func (ml *map_locationType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	location, err := toLocationPairs(scalars)
	if err != nil {
		return err
	}
	ml.appendPoints(location)
	return nil
}

func (ml *map_locationType) appendPoints(location locationPairs) {
	if ml.isRouteComponent {
		ml.vd.crossingFinder.addLocation(ml.locIndex, locationIndexType(len(ml.location)),
			location)
	}
	ml.location = append(ml.location, location...)
}

func (ml *map_locationType) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return ml.style, ml.attestation
}

func (ml *map_locationType) setCrossings(crossings latlongRefs) {
	ml.crossings = crossings
}

func (ml *map_locationType) getCrosspoints() latlongRefs {
	return ml.crossings
}

func (ml *map_locationType) endpointsAndOffsets() (latlongType, latlongType,
		locationIndexType, locationIndexType) {
	off1, off2 := locationIndexType(0), locationIndexType(len(ml.location) - 2)
	return ml.pointAtOffset(off1), ml.pointAtOffset(off2), off1, off2
}

func (ml *map_locationType) oppositeEndpoint(startPoint latlongType,
		) (latlongType, locationIndexType, locationIndexType) {
	off1, off2 := locationIndexType(0), locationIndexType(len(ml.location) - 2)
	pt1, pt2 := ml.pointAtOffset(off1), ml.pointAtOffset(off2)
	if pt1.samePoint(startPoint) {
		off1, off2 = off2, off1
		pt1, pt2 = pt2, pt1
	}
	return pt1, off1, off2
}

func (ml *map_locationType) isPoint() bool {
	return ml.isPointType
}

func (ml *map_locationType) resolveReferenceToLocation(ref latlongRef) *map_locationType {
	return ml
}

func (ml *map_locationType) pointAtOffset(offset locationIndexType) latlongType {
	return ml.location.latlongPair(int(offset))
}

func countRegisteredLocations(vd *VectorData) locationIndexType {
	var count locationIndexType
	for _, item := range vd.mapItems {
		if _, is := item.(*map_locationType); is {
			count++
		}
	}
	return count
}


type mapRadiusType struct {
	mapItemCore
	itemType int
	radius int
}

func newMapRadius(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mr := &mapRadiusType{}
	mr.source = source
	mr.itemType = nameToTypeMap[listType]
	if mr.itemType == 0 {
		return nil, source.Error("unknown object type '%s'", listType)
	}
	return mr, nil
}

func (mr *mapRadiusType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	scalar := scalars[0]
	if !scalar.IsInt() {
		return scalar.Error("radius is not an integer")
	}
	i, err := strconv.Atoi(scalar.String())
	if err != nil {
		return scalar.Error("error convering radius %s: %s", scalar.String(), err)
	}
	mr.radius = i
	return nil
}

