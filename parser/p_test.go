// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package parser

import (
	"strings"
        "testing"

        "potano.misiones/sexp"
)


type testDocType struct {
	objects []*testWorkItem
}

func (td *testDocType) rootWorkItem() testGuideItem {
	return testGuideItem{td, nil, nil}
}

func (td *testDocType) register(workItem *testWorkItem) {
	td.objects = append(td.objects, workItem)
}

type testGuideItem struct {
	doc *testDocType
	parent, curItem *testWorkItem
}

func (tgi testGuideItem) NewChild(listType, listName string, source sexp.ValueSource,
		) (ListItemType, error) {
	newItem, err := newWorkItem(tgi.curItem, source, listType, listName)
	if err != nil {
		return nil, err
	}
	tgi.doc.register(newItem)
	return testGuideItem{tgi.doc, tgi.curItem, newItem}, nil
}

func (wi testGuideItem) SetList(targetName, symbol string, source sexp.ValueSource,
		value ListItemType) error {
	wi.curItem.setList(targetName, symbol, source, value.(testGuideItem).curItem)
	return nil
}

func (wi testGuideItem) SetScalars(targetName string, scalars []sexp.LispScalar) error {
	wi.curItem.setScalars(targetName, scalars)
	return nil
}

type testWorkItem struct {
	listType, name string
	parent *testWorkItem
	childObjects []testWorkItemChild
	values []testWorkItemScalarSet
}

type testWorkItemChild struct {
	targetName, symbol string
	child *testWorkItem
}

type testWorkItemScalarSet struct {
	targetName string
	scalars []sexp.LispScalar
}

func newWorkItem(parent *testWorkItem, source sexp.ValueSource,
		listType, name string) (*testWorkItem, error) {
	par := parent
	wi := &testWorkItem{listType, name, par, nil, nil}
	return wi, nil
}

func (wi *testWorkItem) setList(targetName, symbol string, source sexp.ValueSource,
		value *testWorkItem) error {
	wi.childObjects = append(wi.childObjects, testWorkItemChild{targetName, symbol, value})
	return nil
}

func (wi *testWorkItem) setScalars(targetName string, scalars []sexp.LispScalar) error {
	wi.values = append(wi.values, testWorkItemScalarSet{targetName, scalars})
	return nil
}


type checkWorkItem struct {
	listType, name string
	parent int
	childObjects []checkWorkItemChild
	values []checkWorkItemScalars
}

type checkWorkItemChild struct {
	targetName, symbol string
	wi int
}

type checkWorkItemScalars struct {
	targetName string
	values []string
}

func (td *testDocType) verify(T *testing.T, wantSpec []checkWorkItem) {
	if len(wantSpec) != len(td.objects) {
		T.Fatalf("expected %d generated objects, got %d", len(wantSpec), len(td.objects))
	}
	for wiX, wi := range td.objects {
		want := wantSpec[wiX]
		if wi.listType != want.listType {
			T.Fatalf("object %d: expected list type %s, got %s", wiX,
				want.listType, wi.listType)
		}
		if wi.name != want.name {
			T.Fatalf("object %d: expected list name %s, got %s", wiX,
				want.name, wi.name)
		}
		if want.parent >= 0 &&  wi.parent != td.objects[want.parent] {
			T.Fatalf("object %d: wrong parent", wiX)
		}
		if len(wi.childObjects) != len(want.childObjects) {
			T.Fatalf("object %d: expected %d child lists, got %d", wiX,
				len(want.childObjects), len(wi.childObjects))
		}
		for cX, child := range wi.childObjects {
			wantChild := want.childObjects[cX]
			if child.targetName != wantChild.targetName {
				T.Fatalf("object %d, child %d: expected target %s, got %s", wiX,
					cX, wantChild.targetName, child.targetName)
			}
			if child.symbol != wantChild.symbol {
				T.Fatalf("object %d, child %d: expected symbol %s, got %s", wiX,
					cX, wantChild.symbol, child.symbol)
			}
			if child.child != td.objects[wantChild.wi] {
				T.Fatalf("object %d, child %d: wrong parent", wiX, cX)
			}
		}
		for vX, scalars := range wi.values {
			wantScalars := want.values[vX]
			if scalars.targetName != wantScalars.targetName {
				T.Fatalf("object %d, set %d: expected target %s, got %s", wiX, vX,
					wantScalars.targetName, scalars.targetName)
			}
			if len(scalars.scalars) != len(wantScalars.values) {
				T.Fatalf("object %d, set %d: expected %d scalars, got %d", wiX, vX,
					len(wantScalars.values), len(scalars.scalars))
			}
			for sX, scalar := range scalars.scalars {
				if scalar.String() != wantScalars.values[sX] {
					T.Fatalf("object %d, set %d, value %d: want '%s', got '%s'",
						wiX, vX, sX, wantScalars.values[sX],
						scalar.String())
				}
			}
		}
	}
}



func Test_basicParse(T *testing.T) {
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
	prepared, err := PrepareGrammar(grammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
		(info "info string")
		1 2 3)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	if err != nil {
		T.Fatalf("ParseList error: %s", err)
	}
	testDoc.verify(T, []checkWorkItem{
		{
			"top", "", -1,
			[]checkWorkItemChild{
				{"info", "info", 1},
			},
			[]checkWorkItemScalars{
				{"intvals", []string{"1", "2", "3"}},
			},
		},
		{
			"info", "", 0, nil,
			[]checkWorkItemScalars{
				{"string", []string{"info string"}},
			},
		},
	})
}



var workingGrammar Grammar = Grammar{
	{
		"top", UnnamedList,
		[]SymbolAction{
			{"info", sexp.TList, "info"},
			{"feature", sexp.TList, "feature"},
			{"", sexp.TInt, "intvals"},
		},
		[]TargetSpec{
			{"info", 1, 1, 1},
			{"feature", 0, 0, 1},
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
	{
		"feature", NameRequired,
		[]SymbolAction{
			{"path", sexp.TList, "pointset"},
			{"circle", sexp.TList, "pointset"},
			{"point", sexp.TList, "pointset"},
		},
		[]TargetSpec{
			{"pointset", 1, 0, 1},
		},
	},
	{
		"path", NameOptional,
		[]SymbolAction{
			{"", sexp.TFloat, "coords"},
		},
		[]TargetSpec{
			{"coords", 4, 0, 2},
		},
	},
	{
		"circle", NameOptional,
		[]SymbolAction{
			{"", sexp.TFloat, "coords"},
			{"pixels", sexp.TList, "radius"},
			{"radius", sexp.TList, "radius"},
		},
		[]TargetSpec{
			{"coords", 2, 2, 1},
			{"radius", 1, 1, 1},
		},
	},
	{
		"pixels", UnnamedList,
		[]SymbolAction{
			{"", sexp.TInt, "pixels"},
		},
		[]TargetSpec{
			{"pixels", 1, 1, 1},
		},
	},
	{
		"radius", UnnamedList,
		[]SymbolAction{
			{"", sexp.TNum, "radius"},
		},
		[]TargetSpec{
			{"radius", 1, 1, 1},
		},
	},
	{
		"point", NameOptional,
		[]SymbolAction{
			{"", sexp.TFloat, "latitude longitude"},
		},
		[]TargetSpec{
			{"latitude", 1, 1, 1},
			{"longitude", 1, 1, 1},
		},
	},
}



func Test_namedLists(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(path
					29.48 -83.29
					29.49 -83.29)
				(circle roundy   29.54 -83.40
					(pixels 30)
				)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	if err != nil {
		T.Fatalf("ParseList error: %s", err)
	}
	testDoc.verify(T, []checkWorkItem{
		{
			"top", "", -1,
			[]checkWorkItemChild{
				{"info", "info", 1},
				{"feature", "feature", 2},
			},
			nil,
		},
		{
			"info", "", 0, nil,
			[]checkWorkItemScalars{
				{"string", []string{"info string"}},
			},
		},
		{
			"feature", "shiny", 0,
			[]checkWorkItemChild{
				{"pointset", "path", 3},
				{"pointset", "circle", 4},
			},
			nil,
		},
		{
			"path", "", 2, nil,
			[]checkWorkItemScalars{
				{"coords", []string{"29.48", "-83.29", "29.49", "-83.29"}},
			},
		},
		{
			"circle", "roundy", 2,
			[]checkWorkItemChild{
				{"radius", "pixels", 5},
			},
			[]checkWorkItemScalars{
				{"coords", []string{"29.54", "-83.40"}},
			},
		},
		{
			"pixels", "", 4, nil,
			[]checkWorkItemScalars{
				{"pixels", []string{"30"}},
			},
		},
	})
}



func Test_omittedListName(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature
				(path
					29.48 -83.29
					29.49 -83.29)
				(circle roundy   29.54 -83.40
					(pixels 30)
				)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:3: list type 'feature' requires a name"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_unknownListType(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(features shiny
				(path
					29.48 -83.29
					29.49 -83.29)
				(circle roundy   29.54 -83.40
					(pixels 30)
				)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:3: unrecognized list type features"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_disallowedListType(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(path
					29.48 -83.29
					29.49 -83.29)
			)
			(circle roundy   29.54 -83.40
				(pixels 30)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:8: circle list is not allowed in a top list"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_disallowedValueInList(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(path
					29.48 -83.29
					29.49 -83.29)
				(circle roundy   29.54 -83.40
					(pixels 30)
				)
				boo
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:10: 'boo' value is not allowed in list type feature"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_multitarget(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(point  29.48 -83.29)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	if err != nil {
		T.Fatalf("ParseList error: %s", err)
	}
	testDoc.verify(T, []checkWorkItem{
		{
			"top", "", -1,
			[]checkWorkItemChild{
				{"info", "info", 1},
				{"feature", "feature", 2},
			},
			nil,
		},
		{
			"info", "", 0, nil,
			[]checkWorkItemScalars{
				{"string", []string{"info string"}},
			},
		},
		{
			"feature", "shiny", 0,
			[]checkWorkItemChild{
				{"pointset", "point", 3},
			},
			nil,
		},
		{
			"point", "", 2, nil,
			[]checkWorkItemScalars{
				{"latitude", []string{"29.48"}},
				{"longitude", []string{"-83.29"}},
			},
		},
	})
}



func Test_extraMultitarget(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(point  29.48 -83.29 29.45)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:4: float is illegal in this context"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_missingEntry(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(feature shiny
				(path
					29.48 -83.29
					29.49 -83.29)
				(circle roundy   29.54 -83.40
					(pixels 30)
				)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:1: top list lacks info entry"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_tooFewEntries(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(path
					29.48 -83.29
					29.49)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:4: list requires at least 4 coords entries, got 3"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}



func Test_wrongMultipleEntries(T *testing.T) {
	prepared, err := PrepareGrammar(workingGrammar)
	if err != nil {
		T.Fatal(err.Error())
	}
	sourceDoc := `(top
			(info "info string")
			(feature shiny
				(path
					29.48 -83.29
					29.49 -83.30
					29.50)
			)
		)`
	input, err := sexp.Parse("testfile", strings.NewReader(sourceDoc))
	if err != nil {
		T.Fatalf("sexp.Parse error: %s", err)
	}
	testDoc := &testDocType{}
	_, err = prepared.ParseList(testDoc.rootWorkItem(), input)
	wantErr := "testfile:7: number of coords entries must be a multiple of 2"
	gotErr := ""
	if err != nil {
		gotErr = err.Error()
	}
	if gotErr != wantErr {
		T.Fatalf("ParseList expected error %s, got %s", wantErr, gotErr)
	}
}

