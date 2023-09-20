// Copyright © 2023 Michael Thompson
// SPDX-License-Identifier: GPL-2.0-or-later

package vectordata

import (
	"potano.misiones/sexp"
	"potano.misiones/parser"
)

type VectorDataReader struct {
	grammar parser.PreparedGrammar
	data *VectorData
	fileRootItem parser.ListItemType
}

func NewVectorDataReader(data *VectorData) (*VectorDataReader, error) {
	grammar, err := prepareGrammar()
	return &VectorDataReader{grammar, data, readerValet{data, nil, &mapItemCore{}}}, err
}

func (vdr *VectorDataReader) ConsumeList(list sexp.LispList) error {
	_, err := vdr.grammar.ParseList(vdr.fileRootItem, list)
	return err
}

