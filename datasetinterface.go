// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
DatasetInterface is the interface for working with DSV data.
*/
type DatasetInterface interface {
	Reset() error
	GetMode() int
	SetMode(mode int)
	GetNColumn() int
	GetNRow() int
	SetColumnsType(types []int) error
	GetColumnsType() []int
	GetColumnTypeAt(colidx int) (int, error)
	SetColumnsName(names []string)
	GetColumnsName() []string

	GetColumn(idx int) (col *Column, e error)
	GetRow(idx int) *Row
	GetData() interface{}
	GetDataAsRows() Rows
	GetDataAsColumns() (Columns, error)
	TransposeToColumns() error
	TransposeToRows()

	PushRow(r Row)
	PushRowToColumns(r Row)
	PushColumn(col Column)
}
