// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"potano.misiones/great"
)


func initialLengthUnitMap() map[string]float64 {
	return map[string]float64{
		"meter": 1,
		"meters": 1,
		"miles": great.METERS_PER_MILE,
	}
}


func (vd *VectorData) setLengthUnit(item *mapLengthUnitType) error {
	name := item.name
	if _, exists := vd.lengthUnits[name]; exists {
		return item.Error("length unit %s already set", name)
	}
	if item.numUnits < 0.001 {
		if item.numUnits < 0 {
			return item.Error("%.2f is less than zero")
		}
		return item.Error("%.5f is too close to zero to be usable")
	}
	inTermsOf := item.baseUnit
	meters, exists := vd.lengthUnits[inTermsOf]
	if !exists {
		return item.Error("unknown length unit %s", inTermsOf)
	}
	vd.lengthUnits[name] = meters * item.numUnits
	return nil
}

