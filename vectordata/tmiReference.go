// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"potano.misiones/sexp"
)


// Special mapItemType to note points of reference to threadable map items
// to improve error messages


type threadableMapItemReference struct {
	item threadableMapItemType
	source sexp.ValueSource
}


// Methods satisfying the mapItemType interface

func (ti *threadableMapItemReference) Name() string {
	return ti.item.Name()
}

func (ti *threadableMapItemReference) Source() sexp.ValueSource {
	return ti.item.Source()
}

func (ti *threadableMapItemReference) ItemType() int {
	return ti.item.ItemType()
}

func (ti *threadableMapItemReference) ItemTypeString() string {
	return ti.item.ItemTypeString()
}

func (ti *threadableMapItemReference) noteReferrer(name string, ref mapItemType) error {
	return ti.item.noteReferrer(name, ref)
}

func (ti *threadableMapItemReference) Referrers() []string {
	return ti.item.Referrers()
}

func (ti *threadableMapItemReference) addScalars(targetName string,
		scalars []sexp.LispScalar) error {
	return ti.item.addScalars(targetName, scalars)
}

func (ti *threadableMapItemReference) setMenuitem(layer *map_textType) {
	ti.item.setMenuitem(layer)
}

func (ti *threadableMapItemReference) setPopup(popup *mapPopupType) {
	ti.item.setPopup(popup)
}

func (ti *threadableMapItemReference) setStyle(style *mapStyleType) {
	ti.item.setStyle(style)
}

func (ti *threadableMapItemReference) setAttestation(attestation *mapAttestationType) {
	ti.item.setAttestation(attestation)
}

func (ti *threadableMapItemReference) setHtml(html *map_textType) {
	ti.item.setHtml(html)
}

func (ti *threadableMapItemReference) setRadius(radius *mapRadiusType) {
	ti.item.setRadius(radius)
}

func (ti *threadableMapItemReference) addFeature(feature mapItemType) {
	ti.item.addFeature(feature)
}

func (ti *threadableMapItemReference) setConfigurationItem(item mapItemType) error {
	return ti.item.setConfigurationItem(item)
}

func (ti *threadableMapItemReference) styleAndAttestation() (*mapStyleType, *mapAttestationType) {
	return ti.item.styleAndAttestation()
}

func (ti *threadableMapItemReference) Error(msg string, args ...any) error {
	return ti.source.Error(msg, args...)
}

// Methods satisfying the threadableMapItemType interface

func (ti *threadableMapItemReference) getCrosspoints() latlongRefs {
	return ti.item.getCrosspoints()
}

func (ti *threadableMapItemReference) endpointsAndOffsets() (pt1, pt2 latlongType,
		off1, off2 locationIndexType) {
	return ti.item.endpointsAndOffsets()
}

func (ti *threadableMapItemReference) oppositeEndpoint(refPoint latlongType) (latlongType,
		locationIndexType, locationIndexType) {
	return ti.item.oppositeEndpoint(refPoint)
}

func (ti *threadableMapItemReference) isPoint() bool {
	return ti.item.isPoint()
}

func (ti *threadableMapItemReference) resolveReferenceToLocation(ref latlongRef) *map_locationType {
	return ti.item.resolveReferenceToLocation(ref)
}


func logThreadingError(vd *VectorData, child, parent mapItemType, msg string) {
	if tmir, is := child.(*threadableMapItemReference); is {
		msg += " (" + tmir.Name() + " is defined at " +
			tmir.item.Source().SourceDescription() + ")"
	}
	vd.recordDeferredError(child.Error(msg, child.ItemTypeString(), child.Name(),
		parent.ItemTypeString(), parent.Name()))
}

