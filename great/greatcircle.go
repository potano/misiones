// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package great

//Great-circle computations

import (
	"math"
)

const (
	DEG_TO_RADIANS = math.Pi / 180
	METERS_PER_MILE = 1609.344
	EARTH_RADIUS = 6372768
//Earth radius is computed according to the WGS 84 datum at 30.174861°N, the latitude midway
// between St. Augustine (29.894722°N) and Tallahassee (30.455000°N)
// Calculator URL: https://planetcalc.com/7721/
)

/**
 * Computes the distance in meters between two points given by their latitude and longitude
 * in radians.  Uses the Haversine Formula with an earth's radius at north Florida.
 */
func MetersBetweenPoints(p1Lat, p1Long, p2Lat, p2Long float64) float64 {
	sinDLat := math.Sin((p1Lat - p2Lat) / 2)
	sinDLong := math.Sin((p1Long - p2Long) / 2)
	a := sinDLat * sinDLat + math.Cos(p1Lat) * math.Cos(p2Lat) * sinDLong * sinDLong
	return EARTH_RADIUS * 2 * math.Asin(math.Sqrt(a))
}


/** Computes length in meters of path of latitude/longitude pairs in degrees
 */
func MetersInPath(path []float64) float64 {
	if len(path) < 4 {
		return 0
	}
	var accum float64
	prevLat, prevLong := path[0] * DEG_TO_RADIANS, path[1] * DEG_TO_RADIANS
	for i := 2; i < len(path) - 1; i += 2 {
		lat, long := path[i] * DEG_TO_RADIANS, path[i+1] * DEG_TO_RADIANS
		accum += MetersBetweenPoints(prevLat, prevLong, lat, long)
		prevLat, prevLong = lat, long
	}
	return accum
}

