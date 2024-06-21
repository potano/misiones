// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package sexp

type LispValue interface {
	IsList() bool
	IsSymbol() bool
	IsOperator() bool
	MayBeHead() bool
	IsString() bool
	IsSymbolOrString() bool
	IsNumeric() bool
	IsInt() bool
	IsFloat() bool
	TypeMask() uint32
	Source() ValueSource
	Error(msg string, args... any) SexpError
	Desc() string
}


const (
	TList = 1 << iota
	TSymbol
	TOperator
	TString
	TInt
	TFloat
)

const (
	TMayBeHead = TSymbol | TOperator
	TSymbolOrString = TSymbol | TString
	TNum = TInt | TFloat
)


type ValueSource struct {
	source *sourceInfo
	lineno, typeMask uint32
}

func newValueSource(source *sourceInfo, lineno, typeMask uint32) ValueSource {
	return ValueSource{source, lineno, typeMask}
}

func (b ValueSource) IsList() bool {
	return (b.typeMask & TList) > 0
}

func (b ValueSource) IsSymbol() bool {
	return (b.typeMask & TSymbol) > 0
}

func (b ValueSource) IsOperator() bool {
	return (b.typeMask & TOperator) > 0
}

func (b ValueSource) MayBeHead() bool {
	return (b.typeMask & TMayBeHead) > 0
}

func (b ValueSource) IsString() bool {
	return (b.typeMask & TString) > 0
}

func (b ValueSource) IsSymbolOrString() bool {
	return (b.typeMask & TSymbolOrString) > 0
}

func (b ValueSource) IsNumeric() bool {
	return (b.typeMask & TNum) > 0
}

func (b ValueSource) IsInt() bool {
	return (b.typeMask & TInt) > 0
}

func (b ValueSource) IsFloat() bool {
	return (b.typeMask & TFloat) > 0
}

func (b ValueSource) TypeMask() uint32 {
	return b.typeMask
}

func (b ValueSource) Source() ValueSource {
	return b
}

func (b ValueSource) String() string {
	return formSourceDescription(b.source, b.lineno)
}

func (b ValueSource) Error(msg string, args... any) SexpError {
	return newSexpError(b.source, b.lineno, msg, args...)
}

func (b ValueSource) SourceDescription() string {
	return formSourceDescription(b.source, b.lineno)
}


type LispList struct {
	ValueSource
	head string
	list []LispValue
}

func newLispList(source *sourceInfo, lineno uint32, head string, body []LispValue) LispList {
	return LispList{
		newValueSource(source, lineno, TList),
		head,
		body,
	}
}

func (l LispList) Head() string {
	return l.head
}

func (l LispList) List() []LispValue {
	return l.list
}

func (l LispList) Desc() string {
	return "'" + l.head + "' list"
}


type LispScalar struct {
	ValueSource
	value string
}

func newLispScalar(source *sourceInfo, lineno uint32, tp uint32, value string) LispScalar {
	return LispScalar{newValueSource(source, lineno, tp), value}
}

func (lv LispScalar) String() string {
	return lv.value
}

func (lv LispScalar) Desc() string {
	switch lv.typeMask {
	case TSymbol, TOperator:
		return "'" + lv.value + "'"
	case TString:
		return "string"
	case TInt:
		return "integer"
	case TFloat:
		return "float"
	}
	return "item"
}


func newLispSymbol(source *sourceInfo, lineno uint32, value string) LispScalar {
	return newLispScalar(source, lineno, TSymbol, value)
}

func newLispOperator(source *sourceInfo, lineno uint32, value string) LispScalar {
	return newLispScalar(source, lineno, TOperator, value)
}

func newLispString(source *sourceInfo, lineno uint32, value string) LispScalar {
	return newLispScalar(source, lineno, TString, value)
}

func newLispNumber(source *sourceInfo, lineno uint32, value string) LispScalar {
	return newLispScalar(source, lineno, TNum, value)
}

func newLispInteger(source *sourceInfo, lineno uint32, value string) LispScalar {
	return newLispScalar(source, lineno, TInt, value)
}

func newLispFloat(source *sourceInfo, lineno uint32, value string) LispScalar {
	return newLispScalar(source, lineno, TFloat, value)
}


type endList struct {
	LispScalar
}

func newEndList(source *sourceInfo, lineno uint32) endList {
	return endList{newLispScalar(source, lineno, 0, ")")}
}


var dummyValue LispScalar
var dummyList LispList

