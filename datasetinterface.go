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
	SetNColumn(n int)
	GetNRow() int
	SetNRow(n int)
	SetColumnType(types []int)
	GetColumnTypeAt(colidx int) (int, error)

	GetData() interface{}
	GetDataAsRows() Rows
	GetDataAsColumns() Columns
	TransposeToColumns()
	TransposeToRows()

	PushRow(r Row)
	PushRowToColumns(r Row) error
}
