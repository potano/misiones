// Copyright © 2024 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import "testing"

//Tests conversions to and from locAngleType

func Test_parseLocAngleSuccess(T *testing.T) {
	for _, tst := range []struct{str string; want locAngleType} {
		{"0.0", 0},
		{"3.141", 3141000},
		{"29.8", 29800000},
		{"29.82", 29820000},
		{"29.821", 29821000},
		{"29.8216", 29821600},
		{"29.82164", 29821640},
		{"29.821643", 29821643},
		{"29.8216434", 29821643},
		{"29.8216435", 29821644},
		{"-81.0", -81000000},
		{"-81.01", -81010000},
		{"-81.015", -81015000},
		{"-81.0157", -81015700},
		{"-81.01572", -81015720},
		{"-81.015729", -81015729},
		{"-81.0157294", -81015729},
		{"-81.0157295", -81015730},
		{"0.02", 20000},
		{"-0.02", -20000},
		{".123", 123000},
		{"30.", 30000000},
		{"-0.0", 0},
		{"-.0", 0},
		{"-0.", 0},
		{"180.000000", 180000000},
		{"-180.000000", -180000000},
	} {
		v, err := parseToLocAngle(tst.str, maxLong)
		if err != nil {
			T.Fatalf("parsing %s: %s", tst.str, err)
		}
		if v != tst.want {
			T.Fatalf("parsing %s: wanted %d, got %d", tst.str, tst.want, v)
		}
	}
}

func Test_parseLocAngleFailure(T *testing.T) {
	for _, tst := range []struct{str string; msg string} {
		{" ", "' ' is not a floating-point number"},
		{"4", "'4' is not a floating-point number"},
		{".", "'.' is not a floating-point number"},
		{".-0", "'.-0' is not a floating-point number"},
		{"4-0", "'4-0' is not a floating-point number"},
		{"1.23.45", "'1.23.45' is not a floating-point number"},
		{"9999.", "9999. is out of range"},
		{"180.000001", "180.000001 is out of range"},
		{"-180.000001", "-180.000001 is out of range"},
	} {
		v, err := parseToLocAngle(tst.str, maxLong)
		if err == nil {
			T.Fatalf("parsing %s: expected error, got value %d", tst.str, v)
		}
		errmsg := err.Error()
		if errmsg != tst.msg {
			T.Fatalf("parsing %s: expected error %s, got %s", tst.str, tst.msg, errmsg)
		}
	}
}

func Test_parseLocAngleLatitudeFailure(T *testing.T) {
	for _, tst := range []struct{str string; msg string} {
		{"90.000001", "90.000001 is out of range"},
		{"-90.000001", "-90.000001 is out of range"},
	} {
		v, err := parseToLocAngle(tst.str, maxLat)
		if err == nil {
			T.Fatalf("parsing %s: expected error, got value %d", tst.str, v)
		}
		errmsg := err.Error()
		if errmsg != tst.msg {
			T.Fatalf("parsing %s: expected error %s, got %s", tst.str, tst.msg, errmsg)
		}
	}
}

func Test_locAngleToString(T *testing.T) {
	for _, tst := range []struct{locAngle locAngleType; str string} {
		{0, "0.000000"},
		{3141000, "3.141000"},
		{29800000, "29.800000"},
		{29820000, "29.820000"},
		{29821000, "29.821000"},
		{29821600, "29.821600"},
		{29821640, "29.821640"},
		{29821643, "29.821643"},
		{-81000000, "-81.000000"},
		{-81010000, "-81.010000"},
		{-81015000, "-81.015000"},
		{-81015700, "-81.015700"},
		{-81015720, "-81.015720"},
		{-81015729, "-81.015729"},
		{2, "0.000002"},
		{20, "0.000020"},
		{200, "0.000200"},
		{20000, "0.020000"},
		{200000, "0.200000"},
		{2000000, "2.000000"},
		{-2, "-0.000002"},
		{-20, "-0.000020"},
		{-200, "-0.000200"},
		{-2000, "-0.002000"},
		{-20000, "-0.020000"},
		{-200000, "-0.200000"},
		{-2000000, "-2.000000"},
		{123000, "0.123000"},
		{30000000, "30.000000"},
		{180000000, "180.000000"},
		{-180000000, "-180.000000"},
	} {
		str := tst.locAngle.String()
		if str != tst.str {
			T.Fatalf("stringing %d: wanted %s, got %s", tst.locAngle, tst.str, str)
		}
	}
}

