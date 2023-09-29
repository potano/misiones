// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"strconv"

	"potano.misiones/sexp"
)


type mapLengthRangeType struct {
	mapItemCore
	min, max float64
	units string
}

func newMapLengthRange(doc *VectorData, parent mapItemType, listType, listName string,
		source sexp.ValueSource) (mapItemType, error) {
	ml := &mapLengthRangeType{}
	ml.source = source
	ml.name = parent.Name()
	doc.routesToMeasure = append(doc.routesToMeasure, ml)
	return ml, nil
}

func (ml *mapLengthRangeType) ItemType() int {
	return mitLengthRange
}

func (ml *mapLengthRangeType) addScalars(targetName string, scalars []sexp.LispScalar) error {
	switch targetName {
	case "minAndMaxLength":
		var min, max float64
		min, err := strconv.ParseFloat(scalars[0].String(), 64)
		if err == nil {
			max, err = strconv.ParseFloat(scalars[1].String(), 64)
		}
		if err != nil {
			return ml.Error("%s", err)
		}
		if min < 0 || max < 0 {
			return ml.Error("negative minimum or maximum length")
		}
		ml.min, ml.max = min, max
	case "units":
		ml.units = scalars[0].String()
	}
	return nil
}



type routeMeasurement struct {
	Name, MeasurementUnit string
	Meters, InUnits float64
	LowBound, HighBound float64
	Err error
	InRange bool
}

func (vd *VectorData) MeasureRoutesToMeasure() []routeMeasurement {
	measurements := make([]routeMeasurement, len(vd.routesToMeasure))
	for routeX, setup := range vd.routesToMeasure {
		measurements[routeX] = vd.measureOneRoute(setup)
	}
	return measurements
}

func (vd *VectorData) measureOneRoute(setup *mapLengthRangeType) routeMeasurement {
	units := setup.units
	measurement := routeMeasurement{
		Name: setup.name,
		MeasurementUnit: units,
		LowBound: setup.min,
		HighBound: setup.max,
	}
	metersPerUnit, exists := vd.lengthUnits[units]
	if !exists {
		measurement.Err = setup.Error("measurement unit '%s' is undefined", units)
		return measurement
	}
	meters, err := vd.MeasurePath(setup.name)
	if err != nil {
		measurement.Err = setup.Error("%s when measuring path", err)
		return measurement
	}
	numUnits := meters / metersPerUnit
	measurement.Meters = meters
	measurement.InUnits = numUnits
	measurement.InRange =  numUnits >= setup.min && (setup.max == 0 || numUnits <= setup.max)
	return measurement
}

