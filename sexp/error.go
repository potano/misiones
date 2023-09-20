// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package sexp

import "fmt"

type SexpError struct {
	message string
}

func newSexpError(source *sourceInfo, lineno uint32, msg string, args... any) SexpError {
	msg = fmt.Sprintf(msg, args...)
	prefix := formSourceDescription(source, lineno)
	return SexpError{prefix + ": " + msg}
}

func (e SexpError) Error() string {
	return e.message
}


func formSourceDescription(source *sourceInfo, lineno uint32) string {
	var desc string
	if len(source.filename) > 0 {
		desc = source.filename + ":"
		if lineno > 0 {
			desc += fmt.Sprintf("%d", lineno)
		}
	} else if lineno > 0 {
		desc = fmt.Sprintf("line %d", lineno)
	}
	return desc
}

