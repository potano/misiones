// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
        "potano.misiones/sexp"
        "potano.misiones/parser"
)


//go:generate go run mk_enums.go

func prepareGrammar() (parser.PreparedGrammar, error) {
	return parser.PrepareGrammar(parser.Grammar{
		{
			"0", parser.UnnamedList,
			[]parser.SymbolAction{
				{"layers", sexp.TList, "layer"},
				{"feature", sexp.TList, "feature"},
				{"segment", sexp.TList, "segment"},
				{"route", sexp.TList, "route"},
				{"path", sexp.TList, "path"},
				{"config", sexp.TList, "configItem"},
			},
			[]parser.TargetSpec{
				{"layer", 0, 1, 1},
				{"feature", 0, 0, 1},
				{"segment", 0, 0, 1},
				{"route", 0, 0, 1},
				{"path", 0, 0, 1},
				{"configItem", 0, 1, 1},
			},
		},
		{
			"layers", parser.UnnamedList,
			[]parser.SymbolAction{
				{"layer", sexp.TList, "layer"},
			},
			[]parser.TargetSpec{
				{"layer", 1, 0, 1},
			},
		},
		{
			"layer", parser.NameRequired,
			[]parser.SymbolAction{
				{"menuitem", sexp.TList, "menuitem"},
				{"features", sexp.TList, "feature"},
			},
			[]parser.TargetSpec{
				{"menuitem", 1, 1, 1},
				{"feature", 1, 0, 1},
			},
		},
		{
			"menuitem", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TString, "menuitem"},
			},
			[]parser.TargetSpec{
				{"menuitem", 1, 1, 1},
			},
		},
		{
			"features", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TSymbolOrString, "features"},
			},
			[]parser.TargetSpec{
				{"features", 1, 0, 0},
			},
		},
		{
			"feature", parser.NameOptional,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"marker", sexp.TList, "marker"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"point", sexp.TList, "feature"},
				{"path", sexp.TList, "feature"},
				{"route", sexp.TList, "feature"},
				{"polygon", sexp.TList, "feature"},
				{"rectangle", sexp.TList, "feature"},
				{"circle", sexp.TList, "feature"},
				{"feature", sexp.TList, "feature"},
				{"features", sexp.TList, "feature"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 0},
				{"marker", 0, 0, 0},
				{"style", 0, 1, 0},
				{"attestation", 0, 1, 0},
				{"feature", 0, 0, 1},
			},
		},
		{
			"marker", parser.NameOptional,
			[]parser.SymbolAction{
				{"html", sexp.TList, "html"},
				{"popup", sexp.TList, "popup"},
				{"", sexp.TFloat, "coordinates"},
			},
			[]parser.TargetSpec{
				{"html", 0, 1, 0},
				{"popup", 0, 1, 0},
				{"coordinates", 2, 2, 1},
			},
		},
		{
			"html", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TString, "html"},
			},
			[]parser.TargetSpec{
				{"html", 1, 0, 1},
			},
		},
		{
			"popup", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TString, "text"},
			},
			[]parser.TargetSpec{
				{"text", 1, 0, 1},
			},
		},
		{
			"style", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TSymbol, "symbol"},
			},
			[]parser.TargetSpec{
				{"symbol", 1, 1, 1},
			},
		},
		{
			"attestation", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TSymbol, "attestation"},
			},
			[]parser.TargetSpec{
				{"attestation", 1, 0, 1},
			},
		},
		{
			"point", parser.NameOptional,
			[]parser.SymbolAction{
				{"", sexp.TFloat, "coordinates"},
			},
			[]parser.TargetSpec{
				{"coordinates", 2, 2, 1},
			},
		},
		{
			"paths", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TSymbol, "reference"},
			},
			[]parser.TargetSpec{
				{"reference", 1, 0, 1},
			},
		},
		{
			"path", parser.NameOptional,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"", sexp.TFloat, "points"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 1},
				{"style", 0, 1, 1},
				{"attestation", 0, 1, 1},
				{"points", 4, 0, 2},
			},
		},
		{
			"route", parser.NameRequired,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"segment", sexp.TList, "feature"},
				{"path", sexp.TList, "feature"},
				{"segments", sexp.TList, "feature"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 1},
				{"style", 0, 1, 1},
				{"attestation", 0, 1, 1},
				{"feature", 0, 0, 1},
			},
		},
		{
			"rectangle", parser.NameOptional,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"", sexp.TFloat, "points"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 1},
				{"style", 0, 1, 1},
				{"attestation", 0, 1, 1},
				{"points", 8, 8, 0},
			},
		},
		{
			"polygon", parser.NameOptional,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"", sexp.TFloat, "points"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 1},
				{"style", 0, 1, 1},
				{"attestation", 0, 1, 1},
				{"points", 4, 0, 2},
			},
		},
		{
			"circle", parser.NameOptional,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"", sexp.TFloat, "points"},
				{"radius", sexp.TList, "radius"},
				{"pixels", sexp.TList, "radius"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 1},
				{"style", 0, 1, 1},
				{"attestation", 0, 1, 1},
				{"points", 2, 2, 0},
				{"radius", 1, 1, 0},
			},
		},
		{
			"radius", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TInt, "meters"},
			},
			[]parser.TargetSpec{
				{"meters", 1, 1, 0},
			},
		},
		{
			"pixels", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TInt, "pixels"},
			},
			[]parser.TargetSpec{
				{"pixels", 1, 1, 0},
			},
		},
		{
			"segment", parser.NameOptional,
			[]parser.SymbolAction{
				{"popup", sexp.TList, "popup"},
				{"style", sexp.TList, "style"},
				{"attestation", sexp.TList, "attestation"},
				{"path", sexp.TList, "path"},
				{"paths", sexp.TList, "path"},
			},
			[]parser.TargetSpec{
				{"popup", 0, 1, 1},
				{"style", 0, 1, 1},
				{"attestation", 0, 1, 1},
				{"path", 1, 0, 1},
			},
		},
		{
			"segments", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TSymbol, "reference"},
			},
			[]parser.TargetSpec{
				{"reference", 1, 0, 1},
			},
		},
		{
			"config", parser.UnnamedList,
			[]parser.SymbolAction{
				{"baseStyle", sexp.TList, "configItem"},
				{"attestationType", sexp.TList, "configItem"},
			},
			[]parser.TargetSpec{
				{"configItem", 1, 0, 1},
			},
		},
		{
			"baseStyle", parser.NameRequired,
			[]parser.SymbolAction{
				{"", sexp.TString, "baseStyleProperty"},
			},
			[]parser.TargetSpec{
				{"baseStyleProperty", 0, 0, 1},
			},
		},
		{
			"attestationType", parser.NameRequired,
			[]parser.SymbolAction{
				{"", sexp.TSymbol, "typeSymbol"},
				{"attSym", sexp.TList, "configItem"},
				{"modStyle", sexp.TList, "configItem"},
			},
			[]parser.TargetSpec{
				{"typeSymbol", 1, 1, 1},
				{"configItem", 1, 0, 1},
			},
		},
		{
			"attSym", parser.NameRequired,
			[]parser.SymbolAction{
				{"", sexp.TString, "attSymKeyValue"},
				{"modStyle", sexp.TString, "configItem"},
			},
			[]parser.TargetSpec{
				{"attSymKeyValue", 0, 0, 0},
				{"configItem", 0, 1, 1},
			},
		},
		{
			"modStyle", parser.UnnamedList,
			[]parser.SymbolAction{
				{"", sexp.TString, "baseStyleProperty"},
			},
			[]parser.TargetSpec{
				{"baseStyleProperty", 0, 0, 1},
			},
		},
	})
}

