// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"os"
	"fmt"
	"strings"
)

type dbgType int
var dbg dbgType

func (d *dbgType) Printf(msg string, args ...any) {
	fmt.Printf(msg, args...)
	if *d > 0 {
		*d--
		if *d <= 0 {
			os.Exit(0)
		}
	}
}

func (d *dbgType) Join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func (d *dbgType) Exit(rc int) {
	os.Exit(rc)
}

