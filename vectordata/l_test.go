// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"io"
	"fmt"
	"math"
	"strings"
	"testing"
)


func checkRouteMeasurments(T *testing.T, wantList, gotList []routeMeasurement) {
	T.Helper()
	if len(wantList) != len(gotList) {
		T.Fatalf("wanted %d measurment results, got %d", len(wantList), len(gotList))
	}
	for _, want := range wantList {
		var got routeMeasurement
		wantName := want.Name
		for _, gg := range gotList {
			if gg.Name == wantName {
				got = gg
				break
			}
		}
		if len(got.Name) == 0 {
			T.Fatalf("did not find measurement for route %s", wantName)
		}
		var wantErr, gotErr string
		if want.Err != nil {
			wantErr = want.Err.Error()
		}
		if got.Err != nil {
			gotErr = got.Err.Error()
		}
		if wantErr != gotErr {
			if len(wantErr) == 0 {
				T.Fatalf("%s: got unexpected error %s", wantName, gotErr)
			}
			if len(gotErr) == 0 {
				T.Fatalf("%s: expected error %s", wantName, wantErr)
			}
			T.Fatalf("%s: expected error '%s', got '%s'", wantName, wantErr, gotErr)
		}
		if len(wantErr) > 0 {
			continue
		}
		if got.MeasurementUnit != want.MeasurementUnit {
			T.Fatalf("%s: expected measurement to be in %s, not %s", wantName,
				want.MeasurementUnit, got.MeasurementUnit)
		}
		if math.Abs(got.Meters - want.Meters) > 0.05 {
			T.Fatalf("%s: expected %.2f meters, got %.2f", wantName, want.Meters,
				got.Meters)
		}
		if math.Abs(got.InUnits - want.InUnits) > 0.05 {
			T.Fatalf("%s: expected %.1f %s, got %.1f", wantName, want.InUnits,
				want.MeasurementUnit, got.InUnits)
		}
		if got.LowBound != want.LowBound || got.HighBound != want.HighBound {
			T.Fatalf("%s: expected range [%f, %f] %s, got [%f, %f]", wantName,
				want.LowBound, want.HighBound, want.MeasurementUnit,
				got.LowBound, got.HighBound)
		}
		if got.InRange != want.InRange {
			T.Fatalf("%s: expected inRange==%t, got %t", wantName, want.InRange,
				got.InRange)
		}
	}
}



func Test_measureExpectedRouteLengthMeters(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(lengthRange 1600 1700 meters)
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})
	measurements := vd.MeasureRoutesToMeasure()
	want := []routeMeasurement{
		{"theRoad", "meters", 1634.76, 1634.76, 1600, 1700, nil, true},
	}
	checkRouteMeasurments(T, want, measurements)
}


func Test_measureExpectedRouteLengthUnknownUnit(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(lengthRange 79 80 chains)
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})
	measurements := vd.MeasureRoutesToMeasure()
	wantErr := "infile0:8: measurement unit 'chains' is undefined"
	want := []routeMeasurement{
		{"theRoad", "chains", 0, 0, 0, 0, fmt.Errorf(wantErr), true},
	}
	checkRouteMeasurments(T, want, measurements)
}


func Test_measureExpectedRouteLengthLeagues(T *testing.T) {
	sourceText := `(layers
		(layer one
			(menuitem "Look")
			(features theRoad)
		)
	)
	(route theRoad
		(lengthRange 1 1.5 leagues)
		(segment roadSeg1
			(paths path1 path2)
		)
		(segment roadSeg2
			(paths path3 path4 path5 path6)
		)
	)
	(config
		(lengthUnit league 3 miles)
		(lengthUnit leagues 1 league)
	)
	` + path1 + path2 + path3 + path4 + path5 + path6
	vd := prepareAndParse(T, []io.Reader{strings.NewReader(sourceText)})
	measurements := vd.MeasureRoutesToMeasure()
	want := []routeMeasurement{
		{"theRoad", "leagues", 1634.76, 0.34, 1, 1.5, nil, false},
	}
	checkRouteMeasurments(T, want, measurements)
}

