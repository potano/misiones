// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"regexp"
	"strings"
	"strconv"

	"potano.misiones/sexp"
)


type cssPropertyValue string
type cssPropertyMap map[string]cssPropertyValue

func (c cssPropertyValue) asInt() (int, error) {
	return strconv.Atoi(string(c))
}

func (c cssPropertyValue) asFloat() (float64, error) {
	return strconv.ParseFloat(string(c), 64)
}

func (c cssPropertyValue) asBool() (bool, error) {
	return strconv.ParseBool(string(c))
}

func (c cssPropertyValue) jsonForm() string {
	if len(parseCssValueRegex.FindString(string(c))) > 0 {
		return string(c)
	}
	return strconv.Quote(string(c))
}

var parseCssValueRegex *regexp.Regexp = regexp.MustCompile(
	"^(?:(?:[+-]?(?:[0-9]+(?:\\.[0-9]*)?|\\.[0-9]+))|true|false)$")




func decomposeKeyValueScalar(scalar sexp.LispScalar) (string, cssPropertyValue, error) {
	str := scalar.String()
	ind := strings.IndexByte(str, '=')
	if ind < 0 {
		return "", "", scalar.Error("malformed property '%s'", str)
	}
	return str[:ind], cssPropertyValue(str[ind+1:]), nil
}

