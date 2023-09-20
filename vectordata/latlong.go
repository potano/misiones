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

const (
	noPathMatch = iota
	pathMatchForward
	pathMatchReverse
)


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

func (lp locationPairs) pathEndpointMatch(lat, long float64) int {
	if len(lp) < 2 {
		return noPathMatch
	}
	latDiff, longDiff := lat - lp[0], long - lp[1]
	if -samePointDistance < latDiff && latDiff < samePointDistance &&
			-samePointDistance < longDiff && longDiff < samePointDistance {
		return pathMatchForward
	}
	latDiff, longDiff = lat - lp[len(lp)-2], long - lp[len(lp)-1]
	if -samePointDistance < latDiff && latDiff < samePointDistance &&
			-samePointDistance < longDiff && longDiff < samePointDistance {
		return pathMatchReverse
	}
	return noPathMatch
}

func (lp locationPairs) oppositeEndpoint(nearEndpoint int) (float64, float64) {
	if len(lp) < 2 {
		return 0, 0
	}
	if nearEndpoint == pathMatchForward {
		return lp[len(lp)-2], lp[len(lp)-1]
	}
	return lp[0], lp[1]
}

