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

type latlongType struct {
	lat, long locAngleType
}

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

func (lp locationPairs) latlongPair(i int) latlongType {
	return latlongType{lp[i], lp[i+1]}
}

func (lp locationPairs) asFloatSlice() []float64 {
	out := make([]float64, len(lp))
	for i, v := range lp {
		out[i] = float64(v) * latLongFixedToFloatMultiplier
	}
	return out
}

func (lp locationPairs) asReverseFloatSlice() []float64 {
	lastLP := len(lp) - 1
	out := make([]float64, len(lp))
	for i, v := range lp {
		if i & 1 > 0 {
			i = (lastLP - i) + 1
		} else {
			i = (lastLP - i) - 1
		}
		out[i] = float64(v) * latLongFixedToFloatMultiplier
	}
	return out
}


func (ll latlongType) samePoint(ll2 latlongType) bool {
	return ll.lat == ll2.lat && ll.long == ll2.long
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

