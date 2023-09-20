// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package sexp

import (
	"bufio"
	"io"
	"encoding/base64"
	"unicode"
)

type sourceInfo struct {
	filename string
}


type parser struct {
	source *sourceInfo
	input *bufio.Scanner
	workline []byte
	lineno uint32
}


func Parse(filename string, input io.Reader) (LispList, error) {
	p := &parser{&sourceInfo{filename}, bufio.NewScanner(input), nil, 0}
	var rootItems []LispValue
	var returnErr error
	for {
		list, err := p.next()
		if err != nil {
			if err != io.EOF {
				returnErr = err
			}
			break
		} else if !list.IsList() {
			if _, isEndmark := list.(endList); isEndmark {
				returnErr = p.newError("Unmatched closing parenthesis")
			} else {
				returnErr = p.newError("Top-level value is not a list")
			}
			break;
		}
		rootItems = append(rootItems, list)
	}
	if len(rootItems) == 1 {
		return rootItems[0].(LispList), returnErr
	}
	return newLispList(p.source, 1, "0", rootItems), returnErr
}


func (p *parser) next() (LispValue, error) {
	again:
	for len(p.workline) > 0 && p.workline[0] <= ' ' {
		p.workline = p.workline[1:]
	}
	if len(p.workline) == 0 {
		if !p.input.Scan() {
			err := p.input.Err()
			if err != nil {
				return dummyValue, err
			}
			return dummyValue, io.EOF
		}
		p.workline = p.input.Bytes()
		p.lineno++
		goto again
	}
	c1 := p.workline[0]
	switch c1 {
	case '(':
		p.workline = p.workline[1:]
		return p.parseList()
	case ')':
		p.workline = p.workline[1:]
		return newEndList(p.source, p.lineno), nil
	case '"', '\'':
		return p.parseLiteralString()
	case ';':
		p.workline = nil
		goto again
	case '#':
		return p.parseHexLiteral()
	case '|':
		return p.parseBase64Literal()
	default:
		pos := symbolBreakpos(p.workline)
		value := string(p.workline[:pos])
		p.workline = p.workline[pos:]
		if isNumeric, isLegal, isFloat := isLegalNumeral(value); isNumeric {
			if !isLegal {
				return dummyValue, p.newError("Illegal number %s", value)
			}
			if isFloat {
				return newLispFloat(p.source, p.lineno, value), nil
			}
			return newLispInteger(p.source, p.lineno, value), nil
		}
		if isSymbolic, isLegal := isLegalIdentifier(value); isSymbolic {
			if !isLegal {
				return dummyValue, p.newError("Illegal symbol %s", value)
			}
			return newLispSymbol(p.source, p.lineno, value), nil
		}
		return newLispOperator(p.source, p.lineno, value), nil
	}
}


func (p *parser) parseList() (LispList, error) {
	list := make([]LispValue, 0, 2)
	startline := p.lineno
	for {
		item, err := p.next()
		if err != nil {
			if err == io.EOF {
				return dummyList,
					newSexpError(p.source, startline, "Unterminated list")
			}
			return dummyList, err
		}
		if _, isEnd := item.(endList); isEnd {
			break
		}
		list = append(list, item)
	}
	var head string
	if len(list) > 0 && list[0].MayBeHead() {
		head = list[0].(LispScalar).String()
		list = list[1:]
	}
	return newLispList(p.source, startline, head, list), nil
}


func (p *parser) parseLiteralString() (LispScalar, error) {
	c1 := p.workline[0]
	out := make([]byte, len(p.workline))
	var outpos int
	var escape bool
	for i, c := range p.workline[1:] {
		if escape {
			switch c {
			case 'n':
				c = '\n'
			case 'r':
				c = '\r'
			case 't':
				c = '\t'
			}
			out[outpos] = c
			outpos++
			escape = false
		} else if c == '\\' {
			escape = true
		} else if c == c1 {
			p.workline = p.workline[i+2:]
			return newLispString(p.source, p.lineno, string(out[:outpos])), nil
		} else {
			out[outpos] = c
			outpos++
		}
	}
	return dummyValue, p.newError("Unterminated string")
}


func (p *parser) parseHexLiteral() (LispScalar, error) {
	out := make([]byte, len(p.workline) >> 1)
	var inpos, outpos int
	var c, highByte byte
	for inpos, c = range p.workline[1:] {
		if c >= '0' && c <= '9' {
			c -= '0'
		} else if c >= 'A' && c <= 'F' {
			c -= 'A' - 10
		} else if c >= 'a' && c <= 'f' {
			c -= 'a' - 10
		} else {
			break
		}
		if (inpos & 1) > 0 {
			out[outpos] = (highByte << 4) | c
			outpos++
		} else {
			highByte = c
		}
	}
	if outpos == 0 {
		return dummyValue, p.newError("Expected at least one hex digit")
	}
	p.workline = p.workline[inpos+1:]
	return newLispString(p.source, p.lineno, string(out[:outpos])), nil
}


func (p *parser) parseBase64Literal() (LispScalar, error) {
	pos := symbolBreakpos(p.workline[1:])
	if pos < 2 {
		return dummyValue, p.newError("Expected at least one base-64 character")
	}
	out := make([]byte, base64.StdEncoding.DecodedLen(len(p.workline) - 1))
	n, err := base64.StdEncoding.Decode(out, p.workline[1:pos+1])
	if err != nil {
		return dummyValue, p.newError("%s decoding base-64 literal", err)
	}
	p.workline = p.workline[pos+1:]
	return newLispString(p.source, p.lineno, string(out[:n])), nil
}


func (p *parser) newError(msg string, args... any) SexpError {
	return newSexpError(p.source, p.lineno, msg, args...)
}





func symbolBreakpos(line []byte) int {
	for i, c := range line {
		if c <= ' ' || c == '"' || c == '\'' || c == '(' || c == ')' ||
				c == ';' || c == '#' || c == '|' || c == '\\' {
			return i
		}
	}
	return len(line)
}


/**
 *  Tests if value is to be classed as numeric, is a legal numeral, and whether an integer or float
 *  A value is classed as numeric if the first character that is not a '+', '-', or '.' is a digit.
 *  Values like "123", "+123.4", ".1", "1way", and "+1$" are thus considered to be numeric.
 *  A numeric value is considered to be a legal numeral if it contains at most one decimal point;
 *  a legal numeral may begin with a plus or minus sign.
 *  A legal numeral is classed as a float if it has a decimal point; it is integer if not.
 *
 *  The first return value indicates if the value is classed as numeric, the second value
 *  indicates whether it is a legal numeral, and the third is true if the number is a float.
 */
func isLegalNumeral(candidate string) (bool, bool, bool) {
	var haveDot, haveDigit, haveOther, isFloat bool
	isLegal := true
	for i, c := range candidate {
		if c >= '0' && c <= '9' {
			if haveOther {
				return false, false, false
			} else {
				haveDigit = true
			}
		} else if c == '.' {
			if haveDot {
				isLegal = false
			}
			isFloat = true
			haveDot = true
		} else if c == '+' || c == '-' {
			if i > 0 {
				isLegal = false
			}
		} else {
			haveOther = true
		}
	}
	return haveDigit, isLegal && !haveOther, isFloat
}


/**
 *  Tests is a value is to be classed as symbol and, if so, if it is a legal symbol
 *  A value is classed as symbolic if it begins with a letter or an underscore
 *  A legal symbol consists of one or more segments, each beginning with a Unicode letter
 *  or underscore followed by zero or more letters, digits, or underscores.  Multiple
 *  segments are joined with periods.
 *  Examples of legal symbols:  "abc", "_123", and "abc.def".
 */
func isLegalIdentifier(candidate string) (bool, bool) {
	var isSymbolClass, needLetterAfterDot bool
	isLegal := true
	for i, c := range candidate {
		if unicode.IsLetter(c) || c == '_' {
			if i == 0 {
				isSymbolClass = true
			}
			needLetterAfterDot = false
		} else if c >= '0' && c <= '9' {
			if needLetterAfterDot {
				isLegal = false
			}
		} else if c == '.' {
			if needLetterAfterDot {
				isLegal = false
			}
			needLetterAfterDot = true
		} else {
			isLegal = false
		}
	}
	return isSymbolClass, isLegal && !needLetterAfterDot
}

