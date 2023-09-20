// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package sexp

import (
	"fmt"
	"strings"
	"testing"
)


type tstValue interface {
	mask() uint32
	isList() bool
	String() string
}

type tstList struct {
	lineno uint32
	head string
	list []tstValue
}

func (tl tstList) mask() uint32 {
	return TList
}

func (tl tstList) isList() bool {
	return true
}

func (tl tstList) String() string {
	return fmt.Sprintf("list %s (line %d)", tl.head, tl.lineno)
}

type tstScalar struct {
	lineno uint32
	tp uint32
	value string
}

func (ts tstScalar) mask() uint32 {
	return ts.tp
}

func (ts tstScalar) isList() bool {
	return false
}

func (ts tstScalar) String() string {
	return fmt.Sprintf("%s (line %d)", tTag(ts.tp), ts.lineno)
}

func tTag(tp uint32) string {
	switch tp {
	case TList:
		return "TList"
	case TSymbol:
		return "TSymbol"
	case TOperator:
		return "TOperator"
	case TString:
		return "TString"
	case TInt:
		return "TInt"
	case TFloat:
		return "TFloat"
	default:
		return "(unknown)"
	}
}

func valueOK(T *testing.T, tst tstValue, obj LispValue) {
	if tst.isList() {
		if !obj.IsList() {
			T.Fatalf("%s: expected list, got %s", tst, tTag(obj.TypeMask()))
		}
		tl := tst.(tstList)
		ol := obj.(LispList)
		if ol.lineno != tl.lineno {
			T.Fatalf("%s: wanted line %d, got %d", tl, tl.lineno, ol.lineno)
		}
		if ol.head != tl.head {
			T.Fatalf("%s: wanted head '%s', got '%s'", tl, tl.head, ol.head)
		}
		if len(ol.list) != len(tl.list) {
			T.Fatalf("%s: expected %d elements, got %d", tl, len(tl.list),
				len(ol.list))
		}
		for i, t := range tl.list {
			valueOK(T, t, ol.list[i])
		}
	} else {
		if obj.TypeMask() != tst.mask() {
			T.Fatalf("%s: expected type %s, got %s", tst, tTag(tst.mask()),
				tTag(obj.TypeMask()))
		}
		ts := tst.(tstScalar)
		os := obj.(LispScalar)
		if os.ValueSource.lineno != ts.lineno {
			T.Fatalf("%s: expected line number %d, got %d", ts, ts.lineno,
				os.ValueSource.lineno)
		}
		if os.String() != ts.value {
			T.Fatalf("%s: expected value '%s', got '%s'", ts, ts.value, os.String())
		}
	}
}


func Test_baseline (T *testing.T) {
	for _, tst := range []struct{input string; tval tstValue} {
		{"(testlist)", tstList{1, "testlist", nil}},
		{"(testlist )", tstList{1, "testlist", nil}},
		{"(testlist\n)", tstList{1, "testlist", nil}},
		{"\n(testlist)", tstList{2, "testlist", nil}},
		{"(testlist abc)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc" }}}},
		{"(testlist\nabc)", tstList{1, "testlist", []tstValue{
			tstScalar{2, TSymbol, "abc" }}}},
		{"(testlist \"abc\")", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }}}},
		{"(testlist 'abc')", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }}}},
		{"(testlist 'abc def')", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc def" }}}},
		{"(testlist \"O'Malley\")", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "O'Malley" }}}},
		{"(testlist 'O\\'Malley')", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "O'Malley" }}}},
		{"(testlist \"abc\"\ndef)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }, tstScalar{2, TSymbol, "def"}}}},
		{"(testlist abc\"def\")", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc" }, tstScalar{1, TString, "def"}}}},
		{"(testlist \"abc\"'def')", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }, tstScalar{1, TString, "def"}}}},
		{"(testlist #616263)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }}}},
		{"(testlist #616263 )", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }}}},
		{"(testlist #616263 \"def\")", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc" }, tstScalar{1, TString, "def"}}}},
		{"(testlist |Y29kZWQ=)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "coded" }}}},
		{"(testlist |Y29kZWQ=|Ym9v)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "coded" }, tstScalar{1, TString, "boo"}}}},
		{"(testlist 123)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TInt, "123" }}}},
		{"(testlist +123)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TInt, "+123" }}}},
		{"(testlist -123)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TInt, "-123" }}}},
		{"(testlist 12.3)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TFloat, "12.3" }}}},
		{"(testlist -12.3)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TFloat, "-12.3" }}}},
		{"(testlist -)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, "-" }}}},
		{"(testlist --)", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, "--" }}}},
		{"(+ 123 456)", tstList{1, "+", []tstValue{
			tstScalar{1, TInt, "123" }, tstScalar{1, TInt, "456"}}}},
	} {
		l, err := Parse("testfile", strings.NewReader(tst.input))
		if err != nil {
			T.Fatal(err.Error())
		}
		valueOK(T, tst.tval, l)
	}
}


func Test_list (T *testing.T) {
	for tstnum, tst := range []struct{input string; errmsg string; tval tstValue} {
		{"(testlist)", "", tstList{1, "testlist", nil}},
		{"(testlist", "testfile:1: Unterminated list", nil},
		{"(testlist))", "testfile:1: Unmatched closing parenthesis", nil},
		{"\n(testlist\n1\n'abc'\n)\n", "", tstList{2, "testlist", []tstValue{
			tstScalar{3, TInt, "1"}, tstScalar{4, TString, "abc"}}}},
		{"\n(testlist\n1\n'abc'\n\n", "testfile:2: Unterminated list", nil},
		{"(route here\n(popup \"message\")\n30.1 30.2)", "",
			tstList{1, "route", []tstValue{
				tstScalar{1, TSymbol, "here"},
				tstList{2, "popup", []tstValue{
					tstScalar{2, TString, "message"}}},
				tstScalar{3, TFloat, "30.1"},
				tstScalar{3, TFloat, "30.2"}}}},
		{"(route here\n(popup \"message\"\n30.1 30.2)", "testfile:1: Unterminated list",
			nil},
		{"(list1 abc)\n(list2 (list3 42))", "", tstList{1, "0", []tstValue{
			tstList{1, "list1", []tstValue{
				tstScalar{1, TSymbol, "abc"}}},
			tstList{2, "list2", []tstValue{
				tstList{2, "list3", []tstValue{
					tstScalar{2, TInt, "42"}}}}}}}},
		{"(list1 abc)\n(list2 (list3 42)))", "testfile:2: Unmatched closing parenthesis",
			nil},
		{"(top\n\t(info \"info string\")\n\t1 2 3)", "", tstList{1, "top", []tstValue{
			tstList{2, "info", []tstValue{
				tstScalar{2, TString, "info string"}}},
			tstScalar{3, TInt, "1"},
			tstScalar{3, TInt, "2"},
			tstScalar{3, TInt, "3"}}}},

	} {
		l, err := Parse("testfile", strings.NewReader(tst.input))
		if err != nil {
			if err.Error() != tst.errmsg {
				T.Fatalf("test %d: expected error '%s', got '%s'", tstnum,
					tst.errmsg, err)
			}
		} else if len(tst.errmsg) > 0 {
			T.Fatalf("test %d: expected error '%s'", tstnum, tst.errmsg)
		} else {
			valueOK(T, tst.tval, l)
		}
	}
}


func Test_string (T *testing.T) {
	for tstnum, tst := range []struct{input string; errmsg string; tval tstValue} {
		{"(testlist 'abc')", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "abc"}}}},
		{"(testlist 'abc)", "testfile:1: Unterminated string", nil},
		{"(testlist 'abc\n)", "testfile:1: Unterminated string", nil},
		{"(testlist 'abc\n", "testfile:1: Unterminated string", nil},
		{"(testlist 'a\\nc')", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TString, "a\nc"}}}},
	} {
		l, err := Parse("testfile", strings.NewReader(tst.input))
		if err != nil {
			if err.Error() != tst.errmsg {
				T.Fatalf("test %d: expected error '%s', got '%s'", tstnum,
					tst.errmsg, err)
			}
		} else if len(tst.errmsg) > 0 {
			T.Fatalf("test %d: expected error '%s'", tstnum, tst.errmsg)
		} else {
			valueOK(T, tst.tval, l)
		}
	}
}


func Test_legalNumber(T *testing.T) {
	for _, tst := range []struct{input string; isNumeric, isLegal, isFloat bool} {
		{"123", true, true, false},
		{"3.14", true, true, true},
		{"3.1.4", true, false, true},
		{"+123", true, true, false},
		{"+.123", true, true, true},
		{".+123", true, false, false},
		{"abc", false, false, false},
		{"a1bc", false, false, false},
		{"abc1", false, false, false},
		{"abc1.", false, false, false},
		{"abc$", false, false, false},
		{"123abc", true, false, false},
		{"123$", true, false, false},
		{"$123", false, false, false},
		{"+$123", false, false, false},
		{"+a123", false, false, false},
	} {
		isNumeric, isLegal, isFloat := isLegalNumeral(tst.input)
		if isNumeric != tst.isNumeric {
			T.Fatalf("%s: expected isNumeric=%t, got %t", tst.input, tst.isNumeric,
				isNumeric)
		}
		if isNumeric && isLegal != tst.isLegal {
			T.Fatalf("%s: expected isLegal=%t, got %t", tst.input, tst.isLegal,
				isLegal)
		}
		if isLegal && isFloat != tst.isFloat {
			T.Fatalf("%s: expected isFloat=%t, got %t", tst.input, tst.isFloat,
				isFloat)
		}
	}
}


func Test_legalSymbol(T *testing.T) {
	for _, tst := range []struct{input string; isSymbolic, isLegal bool} {
		{"123", false, false},
		{"abc", true, true},
		{"abc1", true, true},
		{"abc$", true, false},
		{"$abc", false, false},
		{"abc.", true, false},
		{"abc.d", true, true},
		{"abc.1", true, false},
		{"abc.d1", true, true},
		{"abc.d1.", true, false},
		{".abc", false, true},
		{"Ángel", true, true},
		{"caña", true, true},
	} {
		isSymbolic, isLegal := isLegalIdentifier(tst.input)
		if isSymbolic != tst.isSymbolic {
			T.Fatalf("%s: expected isSymbolic=%t, got %t", tst.input, tst.isSymbolic,
				isSymbolic)
		}
		if isSymbolic && isLegal != tst.isLegal {
			T.Fatalf("%s: expected isLegal=%t, got %t", tst.input, tst.isLegal,
				isLegal)
		}
	}
}


func Test_number (T *testing.T) {
	for tstnum, tst := range []struct{input string; errmsg string; tval tstValue} {
		{"(testlist 123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TInt, "123"}}}},
		{"(testlist 3.14)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TFloat, "3.14"}}}},
		{"(testlist 3.1.4)", "testfile:1: Illegal number 3.1.4", nil},
		{"(testlist 3.14)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TFloat, "3.14"}}}},
		{"(testlist +123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TInt, "+123"}}}},
		{"(testlist +.123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TFloat, "+.123"}}}},
		{"(testlist .+123)", "testfile:1: Illegal number .+123", nil},
		{"(testlist abc)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc"}}}},
		{"(testlist a1bc)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "a1bc"}}}},
		{"(testlist abc1)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc1"}}}},
		{"(testlist abc1.)", "testfile:1: Illegal symbol abc1.", nil},
		{"(testlist abc$)", "testfile:1: Illegal symbol abc$", nil},
		{"(testlist 123abc)", "testfile:1: Illegal number 123abc", nil},
		{"(testlist 123$)", "testfile:1: Illegal number 123$", nil},
		{"(testlist $123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, "$123"}}}},
		{"(testlist +$123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, "+$123"}}}},
		{"(testlist +a123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, "+a123"}}}},
	} {
		l, err := Parse("testfile", strings.NewReader(tst.input))
		if err != nil {
			if err.Error() != tst.errmsg {
				T.Fatalf("test %d: expected error '%s', got '%s'", tstnum,
					tst.errmsg, err)
			}
		} else if len(tst.errmsg) > 0 {
			T.Fatalf("test %d: expected error '%s'", tstnum, tst.errmsg)
		} else {
			valueOK(T, tst.tval, l)
		}
	}
}


func Test_symbols (T *testing.T) {
	for tstnum, tst := range []struct{input string; errmsg string; tval tstValue} {
		{"(testlist 123)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TInt, "123"}}}},
		{"(testlist abc)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc"}}}},
		{"(testlist abc1)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc1"}}}},
		{"(testlist abc$)", "testfile:1: Illegal symbol abc$", nil},
		{"(testlist $abc)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, "$abc"}}}},
		{"(testlist abc.)", "testfile:1: Illegal symbol abc.", nil},
		{"(testlist abc.d)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc.d"}}}},
		{"(testlist abc.1)", "testfile:1: Illegal symbol abc.1", nil},
		{"(testlist abc.d1)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "abc.d1"}}}},
		{"(testlist abc.d1.)", "testfile:1: Illegal symbol abc.d1.", nil},
		{"(testlist .abc)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TOperator, ".abc"}}}},
		{"(testlist Ángel)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "Ángel"}}}},
		{"(testlist caña)", "", tstList{1, "testlist", []tstValue{
			tstScalar{1, TSymbol, "caña"}}}},
	} {
		l, err := Parse("testfile", strings.NewReader(tst.input))
		if err != nil {
			if err.Error() != tst.errmsg {
				T.Fatalf("test %d: expected error '%s', got '%s'", tstnum,
					tst.errmsg, err)
			}
		} else if len(tst.errmsg) > 0 {
			T.Fatalf("test %d: expected error '%s'", tstnum, tst.errmsg)
		} else {
			valueOK(T, tst.tval, l)
		}
	}
}

