// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"fmt"
	"strconv"
	"strings"
	"potano.misiones/sexp"
)


type locAngleType int32
type locationPairs []locAngleType

const numFractionalDigits = 6
const latLongFixedToFloatMultiplier = 1e-6

// Latitudes are at even offsets; longitudes ar at off offsets
const maxLat  =  90000000
const maxLong = 180000000
var maxLatLong [2]locAngleType = [2]locAngleType{maxLat, maxLong}


func toLocationPairs(scalars []sexp.LispScalar) (locationPairs, error) {
	out := make([]locAngleType, len(scalars))
	for i, scalar := range scalars {
		angle, err := parseToLocAngle(scalar.String(), maxLatLong[i & 1])
		if err != nil {
			return out, scalar.Error(err.Error())
		}
		out[i] = angle
	}
	return out, nil
}

func (lp locationPairs) generateJs() string {
	out := make([]string, len(lp))
	for i, v := range lp {
		out[i] = v.String()
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (lp locationPairs) asFloatSlice() []float64 {
	out := make([]float64, len(lp))
	for i, v := range lp {
		out[i] = float64(v) * latLongFixedToFloatMultiplier
	}
	return out
}

func (lp locationPairs) indexOfPoint(lat, long locAngleType) int {
	for i := 0; i < len(lp); i += 2 {
		if isSamePoint(lat, long, lp[i], lp[i+1]) {
			return i
		}
	}
	return -1
}

func (lp locationPairs) haveMatchingEndpoint(lat, long locAngleType) bool {
	return isSamePoint(lat, long, lp[0], lp[1]) ||
		isSamePoint(lat, long, lp[len(lp)-2], lp[len(lp)-1])
}

func isSamePoint(lat1, long1, lat2, long2 locAngleType) bool {
	return lat1 == lat2 && long1 == long2
}


func parseToLocAngle(strval string, maxval locAngleType) (locAngleType, error) {
	var asInt locAngleType
	var haveDigits, isNegative, rounded bool
	fracDigits := -1
	for pos, c := range strval {
		switch c {
		case '-':
			if pos > 0 {
				goto notFloat
			}
			isNegative = true
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			haveDigits = true
			if fracDigits < numFractionalDigits {
				asInt *= 10
				asInt += locAngleType(c - '0')
				if fracDigits < 0 {
					//forbid an over-long integer part
					if pos > 4 {
						goto outOfRange
					}
				} else {
					fracDigits++
				}
			} else if !rounded {
				if c - '0' > 4 {
					asInt++
				}
				rounded = true
			}
		case '.':
			if fracDigits >= 0 {
				goto notFloat
			}
			fracDigits = 0
		default:
			goto notFloat
		}
	}
	if !haveDigits || fracDigits < 0 {
		goto notFloat
	}
	for fracDigits < numFractionalDigits {
		asInt *= 10
		fracDigits++
	}
	if asInt > maxval {
		goto outOfRange
	}
	if isNegative {
		asInt = -asInt
	}
	return asInt, nil

	notFloat:
	return 0, fmt.Errorf("'%s' is not a floating-point number", strval)

	outOfRange:
	return 0, fmt.Errorf("%s is out of range", strval)
}

func (v locAngleType) String() string {
	strval := strconv.FormatInt(int64(v), 10)
	var sign string
	if strval[0] == '-' {
		sign = "-"
		strval = strval[1:]
	}
	splitpt := len(strval) - numFractionalDigits
	if splitpt < 1 {
		strval = strings.Repeat("0", 1 - splitpt) + strval
		splitpt = 1
	}
	return sign + strval[:splitpt] + "." + strval[splitpt:]
}

