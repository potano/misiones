// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"strconv"
	"potano.misiones/sexp"
	"potano.misiones/parser"
)


type readerValet struct {
	doc *VectorData
	parent, curItem mapItemType
}

func (rv readerValet) NewChild(listType, listName string, source sexp.ValueSource,
		) (parser.ListItemType, error) {
	var constructor newMapItemFunc
	switch listType {
	case "0":
		return rv, nil
	case "layers":
		constructor = newMapLayers
	case "layer":
		constructor = newMapLayer
	case "menuitem", "html":
		constructor = newMap_text
	case "features", "paths", "segments":
		constructor = newMap_referenceAggregate
	case "feature":
		constructor = newMapFeature
	case "popup":
		constructor = newMapPopup
	case "style":
		constructor = newMapStyle
	case "attestation":
		constructor = newMapAttestation
	case "point", "path", "rectangle", "polygon", "circle", "marker":
		constructor = newMap_location
	case "route":
		constructor = newMapRoute
	case "lengthRange":
		constructor = newMapLengthRange
	case "radius", "pixels":
		constructor = newMapRadius
	case "segment":
		constructor = newMapSegment
	case "config":
		constructor = newMapConfig
	case "baseStyle", "modStyle":
		constructor = newMapStyleConfig
	case "attestationType":
		constructor = newMapAttestationType
	case "attSym":
		constructor = newAttSym
	case "lengthUnit":
		constructor = newMapLengthUnit
	}
	newItem, err := constructor(rv.doc, rv.curItem, listType, listName, source)
	if err != nil {
		return nil, err
	}
	return readerValet{rv.doc, rv.curItem, newItem}, nil
}

func (rv readerValet) SetScalars(targetName string, scalars []sexp.LispScalar) error {
	return rv.curItem.addScalars(targetName, scalars)
}

func (rv readerValet) SetList(targetName, listType string, source sexp.ValueSource,
		value parser.ListItemType) error {
	curItem := rv.curItem
	newChild := value.(readerValet).curItem
	newChild.noteReferrer(curItem.Name(), curItem)
	switch targetName {
	case "menuitem":
		if asMenuitem, is := newChild.(*map_textType); !is {
			return source.Error("not a menuitem")
		} else {
			curItem.setMenuitem(asMenuitem)
		}
	case "popup":
		if asPopup, is := newChild.(*mapPopupType); !is {
			return source.Error("not a popup")
		} else {
			curItem.setPopup(asPopup)
		}
	case "style":
		if asStyle, is := newChild.(*mapStyleType); !is {
			return source.Error("not a style")
		} else {
			curItem.setStyle(asStyle)
		}
	case "attestation":
		if asAttestation, is := newChild.(*mapAttestationType); !is {
			return source.Error("not an attestation")
		} else {
			curItem.setAttestation(asAttestation)
		}
	case "feature":
		curItem.addFeature(newChild)
	case "html":
		if asHtml, is := newChild.(*map_textType); !is {
			return source.Error("not an html")
		} else {
			curItem.setHtml(asHtml)
		}
	case "radius":
		if asRadius, is := newChild.(*mapRadiusType); !is {
			return source.Error("not a radius")
		} else {
			curItem.setRadius(asRadius)
		}
	case "configItem":
		err := curItem.setConfigurationItem(newChild)
		if err != nil {
			return err
		}
	case "lengthRange":
	default:
		return source.Error("** internal error **: unhandled target type %s", targetName)
	}
	return nil
}





type newMapItemFunc func (doc *VectorData, parent mapItemType, listType, listName string,
	source sexp.ValueSource) (mapItemType, error)



type mapItemCore struct {
	name string
	source sexp.ValueSource
	referrers []string
}

func (mic *mapItemCore) Name() string {
	return mic.name
}

func (mic *mapItemCore) Source() sexp.ValueSource {
	return mic.source
}

func (mic *mapItemCore) ItemType() int {
	return 0
}

func (mic *mapItemCore) noteReferrer(name string, referrer mapItemType) error {
	if len(name) > 0 {
		for _, n := range mic.referrers {
			if n == name {
				return referrer.Error("multiple references to target node '%s'",
					mic.Name())
			}
		}
		mic.referrers = append(mic.referrers, name)
	}
	return nil
}

func (mic *mapItemCore) Referrers() []string {
	return mic.referrers
}

func (mic *mapItemCore) addScalars(targetName string, scalars []sexp.LispScalar) error {
	return nil
}

func (mic *mapItemCore) setMenuitem(layer *map_textType) {}
func (mic *mapItemCore) setPopup(popup *mapPopupType) {}
func (mic *mapItemCore) setStyle(style *mapStyleType) {}
func (mic *mapItemCore) setAttestation(attestation *mapAttestationType) {}
func (mic *mapItemCore) setHtml(html *map_textType) {}
func (mic *mapItemCore) setRadius(radius *mapRadiusType) {}
func (mic *mapItemCore) addFeature(feature mapItemType) {}
func (mic *mapItemCore) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return nil, nil}
func (mic *mapItemCore) setConfigurationItem(item mapItemType) error {return nil}


func (mic *mapItemCore) Error(msg string, args ...any) error {
	return mic.source.Error(msg, args...)
}

func (mic *mapItemCore) generateJs() string {
	return ""
}






type map_textType struct {
	mapItemCore
	itemType int
	text string
}

func newMap_text(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mt := &map_textType{}
	mt.source = source
	mt.itemType = nameToTypeMap[listType]
	if mt.itemType == 0 {
		return nil, source.Error("unknown object type '%s'", listType)
	}
	return mt, nil
}

func (mt *map_textType) ItemType() int {
	return mt.itemType
}

func (mt *map_textType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	var text string
	for _, scalar := range scalars {
		text += scalar.String()
	}
	mt.text = text
	return nil
}



type mapFeatureType struct {
	mapItemCore
	popup *mapPopupType
	style *mapStyleType
	attestation *mapAttestationType
	features []mapItemType
}

func newMapFeature(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	if parent == nil && len(listName) == 0 {
		return nil, source.Error("no name given for non-embedded feature")
	}
	mf := &mapFeatureType{}
	mf.source = source
	name, err := doc.registerMapItem(mf, listName)
	mf.name = name
	return mf, err
}

func (mf *mapFeatureType) ItemType() int {
	return mitFeature
}

func (mf *mapFeatureType) setPopup(popup *mapPopupType) {
	mf.popup = popup
}

func (mf *mapFeatureType) setStyle(style *mapStyleType) {
	mf.style = style
}

func (mf *mapFeatureType) setAttestation(attestation *mapAttestationType) {
	mf.attestation = attestation
}

func (mf *mapFeatureType) addFeature(feature mapItemType) {
	mf.features = append(mf.features, feature)
}

func (mf *mapFeatureType) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return mf.style, mf.attestation
}



type mapPopupType struct {
	mapItemCore
	text string
}

func newMapPopup(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	mm := &mapPopupType{}
	mm.source = source
	return mm, nil
}

func (mp *mapPopupType) ItemType() int {
	return mitPopup
}

func (mp *mapPopupType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	mp.text = scalars[0].String()
	return nil
}



type map_locationType struct {
	mapItemCore
	itemType int
	popup *mapPopupType
	html string
	style *mapStyleType
	attestation *mapAttestationType
	radius, radiusType int
	location locationPairs
}

func newMap_location(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	if parent == nil && len(listName) == 0 {
		return nil, source.Error("no name given for non-embedded %s", listType)
	}
	ml := &map_locationType{}
	ml.source = source
	ml.itemType = nameToTypeMap[listType]
	if ml.itemType == 0 {
		return nil, source.Error("unknown object type '%s'", listType)
	}
	name, err := doc.registerMapItem(ml, listName)
	ml.name = name
	return ml, err
}

func (ml *map_locationType) ItemType() int {
	return ml.itemType
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
	var err error
	ml.location, err = toLocationPairs(scalars)
	return err
}

func (ml *map_locationType) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return ml.style, ml.attestation
}



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

func (mr *mapRadiusType) ItemType() int {
	return mr.itemType
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



type mapSegmentType struct {
	mapItemCore
	popup *mapPopupType
	style *mapStyleType
	attestation *mapAttestationType
	paths []mapItemType
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

