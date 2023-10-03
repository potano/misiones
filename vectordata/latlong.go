// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"strconv"
	"strings"
	"potano.misiones/sexp"
)


type locationPairs []float64

const samePointDistance = 4E-7


func toLocationPairs(scalars []sexp.LispScalar) (locationPairs, error) {
	out := make([]float64, len(scalars))
	for i, v := range scalars {
		if !v.IsFloat() {
			return out, v.Error("'%s' is not a float", v.String())
		}
		f, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return out, err
		}
		out[i] = f
	}
	return out, nil
}

func (lp locationPairs) generateJs() string {
	out := make([]string, len(lp))
	for i, v := range lp {
		out[i] = strconv.FormatFloat(v, 'f', 6, 64)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (lp locationPairs) indexOfPoint(lat, long float64) int {
	for i := 0; i < len(lp); i += 2 {
		if isSamePoint(lat, long, lp[i], lp[i+1]) {
			return i
		}
	}
	return -1
}

func (lp locationPairs) haveMatchingEndpoint(lat, long float64) bool {
	return isSamePoint(lat, long, lp[0], lp[1]) ||
		isSamePoint(lat, long, lp[len(lp)-2], lp[len(lp)-1])
}

func isSamePoint(lat1, long1, lat2, long2 float64) bool {
	latDiff := lat1 - lat2
	longDiff := long1 - long2
	return -samePointDistance < latDiff && latDiff < samePointDistance &&
			-samePointDistance < longDiff && longDiff < samePointDistance
}

