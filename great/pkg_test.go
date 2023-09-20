// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package great

import (
	"fmt"
	"testing"
)

func Test_basic(T *testing.T) {
	for i, test := range []struct{lat1, long1, lat2, long2 float64; want float64}{
		{29.98313, -81.31244, 30.43812, -84.28132, 289809.6},
		{29.98313, -81.31244, 29.98314, -81.31244, 1.1},
		{29.98313, -81.31244, 29.98313, -81.31245, 1.0},
		{29.98313, -81.31244, 29.98314, -81.31245, 1.5},
	} {
		lat1 := test.lat1 * DEG_TO_RADIANS
		long1 := test.long1 * DEG_TO_RADIANS
		lat2 := test.lat2 * DEG_TO_RADIANS
		long2 := test.long2 * DEG_TO_RADIANS
		meters := MetersBetweenPoints(lat1, long1, lat2, long2)
		diff := meters - test.want
		if diff > 0.4 || diff < -0.4 {
			T.Fatalf("test %d: expected %.1f meters, got %.1f", i, test.want, meters)
		}
	}
}

func Test_measurePath(T *testing.T) {
	for i, test := range []struct{path []float64; want float64}{
		{[]float64{
			30.390055, -83.869655,
			30.390169, -83.869931,
			30.390382, -83.870403,
			30.390716, -83.870864,
			30.390904, -83.871219}, 138.2},
		{[]float64{
			30.272119, -84.052327,
			30.273955, -84.052246,
			30.274067, -84.052225,
			30.274182, -84.052209,
			30.275661, -84.051851,
			30.278244, -84.051262,
			30.280109, -84.050817,
			30.282247, -84.050354,
			30.286363, -84.049530,
			30.290595, -84.048719,
			30.291669, -84.048505,
			30.291951, -84.048440,
			30.292559, -84.048376,
			30.293159, -84.048314,
			30.296914, -84.048205,
			30.300107, -84.048090,
			30.304414, -84.047902,
			30.307402, -84.047790,
			30.308195, -84.047702,
			30.308802, -84.047627,
			30.309345, -84.047487,
			30.309873, -84.047342,
			30.310399, -84.047080,
			30.311009, -84.046733,
			30.311686, -84.046269,
			30.312557, -84.045669,
			30.313433, -84.045068}, 0},

	} {
		for i, v := range test.path {
			test.path[i] = DEG_TO_RADIANS * v
		}
		meters := MetersInPath(test.path)
		_ = i
		fmt.Printf("%.1f meters (%.3f miles)\n", meters, meters / METERS_PER_MILE)
	}
}


