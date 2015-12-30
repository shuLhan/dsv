// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
DatasetInterface is the interface for working with DSV data.
*/
type DatasetInterface interface {
	Reset()
	GetMode() int
	SetMode(mode int) error
	GetNColumn() int
	GetNRow() int
	SetColumnsType(types []int) error
	GetColumnsType() []int
	GetColumnTypeAt(colidx int) (int, error)
	SetColumnsName(names []string)
	GetColumnsName() []string

	GetColumn(idx int) (col *Column, e error)
	GetData() interface{}
	GetDataAsRows() Rows
	GetDataAsColumns() (Columns, error)
	TransposeToColumns() error
	TransposeToRows()

	PushRow(r Row) error
	PushRowToColumns(r Row) error
	PushColumn(col Column) error
}
