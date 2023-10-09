// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package parser

import (
        "testing"

        "potano.misiones/sexp"
)


func Test_basic(T *testing.T) {
	grammar := Grammar{
		{
			"top", UnnamedList,
			[]SymbolAction{
				{"info", sexp.TList, "info"},
				{"", sexp.TInt, "intvals"},
			},
			[]TargetSpec{
				{"info", 1, 1, 1},
				{"intvals", 0, 0, 1},
			},
		},
		{
			"info", UnnamedList,
			[]SymbolAction{
				{"", sexp.TString, "string"},
			},
			[]TargetSpec{
				{"string", 1, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	if err != nil {
		T.Fatal(err.Error())
	}
}


func Test_redefinedList(T *testing.T) {
	grammar := Grammar{
		{
			"top", UnnamedList,
			[]SymbolAction{
			},
			[]TargetSpec{
			},
		},
		{
			"top", NameOptional,
			[]SymbolAction{
			},
			[]TargetSpec{
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: redefinition of list type top"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_overlappingActions(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"", sexp.TInt, "int"},
				{"", sexp.TNum, "num"},
			},
			[]TargetSpec{
				{"int", 0, 1, 1},
				{"num", 0, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: overlapping token masks in 'worklist'"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_emptyActionName(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"", sexp.TList, "trip"},
			},
			[]TargetSpec{
				{"trip", 0, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: 'worklist' refers to empty list name"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_unknownTargetReference(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"trip", sexp.TList, "nowhere"},
			},
			[]TargetSpec{
				{"trip", 0, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: 'worklist' list has unknown target 'nowhere'"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_listActionScalarTarget(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"trip", sexp.TList, "trip"},
				{"", sexp.TInt, "trip"},
			},
			[]TargetSpec{
				{"trip", 0, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: target 'trip' in list 'worklist' is list target"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_redefinedTarget(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"trip", sexp.TList, "trip"},
			},
			[]TargetSpec{
				{"trip", 0, 1, 1},
				{"trip", 0, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: multiple 'trip' targets in 'worklist'"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_unboundedTarget(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"trip", sexp.TList, "noproblem"},
			},
			[]TargetSpec{
				{"noproblem", 1, 0, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := ""
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}


func Test_derangedTarget(T *testing.T) {
	grammar := Grammar{
		{
			"worklist", UnnamedList,
			[]SymbolAction{
				{"trip", sexp.TList, "bad"},
			},
			[]TargetSpec{
				{"bad", 2, 1, 1},
			},
		},
	}
	_, err := PrepareGrammar(grammar)
	var gotError string
	wantError := "Setup error: bad range in target 'bad' of 'worklist'"
	if err != nil {
		gotError = err.Error()
	}
	if gotError != wantError {
		T.Fatalf("Wanted error %s, got %s", wantError, gotError)
	}
}

