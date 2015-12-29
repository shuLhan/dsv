// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"math"
)

const (
	// DatasetModeRows for output mode in rows.
	DatasetModeRows = 1
	// DatasetModeColumns for output mode in columns.
	DatasetModeColumns = 2
	// DatasetModeMatrix will save data in rows and columns.
	DatasetModeMatrix = 3
)

var (
	// ErrColIdxOutOfRange operation on column index is invalid
	ErrColIdxOutOfRange = errors.New("dsv: Column index out of range")
	// ErrInvalidColType operation on column with different type
	ErrInvalidColType = errors.New("dsv: Invalid column type")
	// ErrEmptySet operation on empty data set.
	ErrEmptySet = errors.New("dsv: dataset is empty")
)

/*
Dataset contain the data, mode of saved data, number of columns and rows in
data.
*/
type Dataset struct {
	// Mode define the numeric value of output mode.
	Mode int
	// Columns is input data that has been parsed.
	Columns Columns
	// Rows is input data that has been parsed.
	Rows Rows
}

/*
NewDataset create new dataset, use the mode to initialize the dataset.
*/
func NewDataset(mode int, types []int) (dataset *Dataset, e error) {
	dataset = &Dataset{}

	e = dataset.Init(mode, types)

	return
}

/*
Init will set the dataset using mode and types.
*/
func (dataset *Dataset) Init(mode int, types []int) (e error) {
	if types == nil {
		dataset.Columns = make(Columns, 0)
	} else {
		dataset.Columns = make(Columns, len(types))
		e = dataset.Columns.SetType(types)
		if e != nil {
			return
		}
	}

	dataset.SetMode(mode)

	return
}

/*
Reset all data and attributes.
*/
func (dataset *Dataset) Reset() {
	dataset.Rows = Rows{}
	dataset.Columns.Reset()
}

/*
GetMode return mode of data.
*/
func (dataset *Dataset) GetMode() int {
	return dataset.Mode
}

/*
SetMode of saved data to `mode`.
*/
func (dataset *Dataset) SetMode(mode int) error {
	switch mode {
	case DatasetModeRows:
		dataset.Mode = DatasetModeRows
		dataset.Rows = make(Rows, 0)
	case DatasetModeColumns:
		dataset.Mode = DatasetModeColumns
		dataset.Columns.Reset()
	case DatasetModeMatrix:
		dataset.Mode = DatasetModeMatrix
		dataset.Rows = make(Rows, 0)
		dataset.Columns.Reset()
	default:
		return ErrUnknownDatasetMode
	}
	dataset.Mode = mode

	return nil
}

/*
GetNColumn return the number of column in dataset.
*/
func (dataset *Dataset) GetNColumn() int {
	return len(dataset.Columns)
}

/*
GetNRow return number of rows in dataset.
*/
func (dataset *Dataset) GetNRow() (nrow int) {
	switch dataset.Mode {
	case DatasetModeRows:
		nrow = len(dataset.Rows)
	case DatasetModeColumns:
		if len(dataset.Columns) <= 0 {
			nrow = 0
		} else {
			// get length of record in the first column
			nrow = dataset.Columns[0].GetLength()
		}
	case DatasetModeMatrix:
		// matrix mode could have empty either in rows or column.
		nrow = len(dataset.Rows)
	}
	return
}

/*
SetColumnsType of data in all columns.
*/
func (dataset *Dataset) SetColumnsType(types []int) (e error) {
	dataset.Columns = make(Columns, len(types))
	e = dataset.Columns.SetType(types)
	return
}

/*
SetColumnsName set column name.
*/
func (dataset *Dataset) SetColumnsName(names []string) {
	nameslen := len(names)

	if nameslen <= 0 {
		// empty names, return immediately.
		return
	}

	collen := dataset.GetNColumn()

	if collen <= 0 {
		dataset.Columns = make(Columns, nameslen)
		collen = nameslen
	}

	// find minimum length
	minlen := collen
	if nameslen < collen {
		minlen = nameslen
	}

	for x := 0; x < minlen; x++ {
		dataset.Columns[x].Name = names[x]
	}

	return
}

/*
GetColumnsType return the type of all columns.
*/
func (dataset *Dataset) GetColumnsType() (types []int) {
	for x := range dataset.Columns {
		types = append(types, dataset.Columns[x].Type)
	}

	return
}

/*
GetColumnsTypeByIdx get column type filtered by column index `colsIdx`.
*/
func (dataset *Dataset) GetColumnsTypeByIdx(colsIdx []int) (
	types []int,
	e error,
) {
	colslen := dataset.GetNColumn()

	for _, v := range colsIdx {
		if v >= colslen {
			return types, ErrMisColLength
		}

		types = append(types, dataset.Columns[v].GetType())
	}

	return
}

/*
GetColumnTypeAt return type of column in index `colidx` in
dataset.
*/
func (dataset *Dataset) GetColumnTypeAt(colidx int) (int, error) {
	if colidx >= dataset.GetNColumn() {
		return TUndefined, ErrColIdxOutOfRange
	}

	return dataset.Columns[colidx].Type, nil
}

/*
GetColumn return pointer to column object at index `idx`.
*/
func (dataset *Dataset) GetColumn(idx int) (col *Column, e error) {
	if idx > dataset.GetNColumn() {
		return
	}

	switch dataset.Mode {
	case DatasetModeRows:
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
	case DatasetModeColumns:
		// do nothing
	case DatasetModeMatrix:
		// do nothing
	}

	return &dataset.Columns[idx], nil
}

/*
GetData return the data, based on mode (rows, columns, or matrix).
*/
func (dataset *Dataset) GetData() interface{} {
	switch dataset.Mode {
	case DatasetModeRows:
		return dataset.Rows
	case DatasetModeColumns:
		return dataset.Columns
	case DatasetModeMatrix:
		return Matrix{
			Columns: &dataset.Columns,
			Rows:    &dataset.Rows,
		}
	}

	return nil
}

/*
GetDataAsRows return data in rows mode.
*/
func (dataset *Dataset) GetDataAsRows() Rows {
	if dataset.Mode == DatasetModeColumns {
		dataset.TransposeToRows()
	}
	return dataset.Rows
}

/*
GetDataAsColumns return data in columns mode.
*/
func (dataset *Dataset) GetDataAsColumns() (columns Columns, e error) {
	if dataset.Mode == DatasetModeRows {
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
	}
	return dataset.Columns, nil
}

/*
TransposeToColumns move all data from rows (horizontal) to columns
(vertical) mode.
*/
func (dataset *Dataset) TransposeToColumns() (e error) {
	if dataset.GetNRow() <= 0 {
		// nothing to transpose
		return
	}
	if dataset.GetNColumn() <= 0 {
		// no column defined
		return ErrMisColLength
	}

	orgmode := dataset.GetMode()

	// check if column records contain data, if its empty transpose to
	// column, otherwise return, which means its already transposed
	if orgmode == DatasetModeColumns || orgmode == DatasetModeMatrix {
		nrow := dataset.Columns[0].GetLength()
		if nrow > 0 {
			// return if column record is not empty, its already
			// transposed
			return
		}
	}

	// double check column length
	collen := len(dataset.Rows[0])

	if collen > dataset.GetNColumn() {
		return ErrMisColLength
	}

	if orgmode == DatasetModeRows {
		dataset.SetMode(DatasetModeColumns)
	}

	for i := range dataset.Rows {
		dataset.PushRowToColumns(dataset.Rows[i])
	}

	// reset the rows data only if original mode is rows
	// this to prevent empty data when mode is matrix.
	if orgmode == DatasetModeRows {
		dataset.Rows = nil
	}

	return
}

/*
TransposeToRows will move all data from columns (vertical) to rows (horizontal)
mode.
*/
func (dataset *Dataset) TransposeToRows() {
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeRows {
		// already transposed
		return
	}

	if orgmode == DatasetModeColumns {
		// only set mode if transposing from columns to rows
		dataset.SetMode(DatasetModeRows)
	}

	// Get the least length of columns.
	rowlen := math.MaxInt32
	flen := len(dataset.Columns)

	for f := 0; f < flen; f++ {
		l := dataset.Columns[f].GetLength()

		if l < rowlen {
			rowlen = l
		}
	}

	dataset.Rows = make(Rows, 0)

	// transpose record from row to column
	for r := 0; r < rowlen; r++ {
		row := make(Row, flen)

		for f := 0; f < flen; f++ {
			row[f] = dataset.Columns[f].Records[r]
		}

		dataset.Rows = append(dataset.Rows, row)
	}

	// Only reset the columns if original dataset mode is "columns".
	// This to prevent empty data when mode is matrix.
	if orgmode == DatasetModeColumns {
		dataset.Columns.Reset()
	}
}

/*
PushRow save the data, which is already in row object, to Rows.
*/
func (dataset *Dataset) PushRow(row Row) (e error) {
	switch dataset.GetMode() {
	case DatasetModeRows:
		dataset.Rows = append(dataset.Rows, row)
	case DatasetModeColumns:
		e = dataset.PushRowToColumns(row)
	case DatasetModeMatrix:
		dataset.Rows = append(dataset.Rows, row)
		e = dataset.PushRowToColumns(row)
	}
	return
}

/*
PushRowToColumns push each data in Row to Columns.
*/
func (dataset *Dataset) PushRowToColumns(row Row) (e error) {
	rowlen := len(row)
	if rowlen <= 0 {
		// return immediately if no data in row.
		return
	}

	// check if columns is initialize.
	collen := len(dataset.Columns)
	if collen <= 0 {
		dataset.Columns = make(Columns, rowlen)
		collen = rowlen
	}

	// pick the minimum length.
	min := rowlen
	if collen < rowlen {
		min = collen
	}

	for x := 0; x < min; x++ {
		dataset.Columns[x].PushBack(row[x])
	}

	return
}

/*
PushColumn will append new column to the end of slice.
*/
func (dataset *Dataset) PushColumn(col Column) (e error) {
	switch dataset.GetMode() {
	case DatasetModeRows:
		e = dataset.PushColumnToRows(col)
	case DatasetModeColumns:
		dataset.Columns = append(dataset.Columns, col)
	case DatasetModeMatrix:
		dataset.Columns = append(dataset.Columns, col)
		e = dataset.PushColumnToRows(col)
	}
	return
}

/*
PushColumnToRows add each record in column to each rows, from top to bottom.
*/
func (dataset *Dataset) PushColumnToRows(col Column) (e error) {
	colsize := col.GetLength()
	if colsize <= 0 {
		// Do nothing if column is empty.
		return
	}

	nrow := dataset.GetNRow()
	if nrow <= 0 {
		// If no existing rows in dataset, initialize the rows slice.
		dataset.Rows = make(Rows, colsize)

		for nrow = 0; nrow < colsize; nrow++ {
			dataset.Rows[nrow] = make(Row, 0)
		}
	}

	// Pick the minimum length between column or current row length.
	minrow := nrow

	if colsize < nrow {
		minrow = colsize
	}

	// Push each record in column to each rows
	var row *Row
	var rec *Record

	for x := 0; x < minrow; x++ {
		row = &dataset.Rows[x]
		rec = col.Records[x]

		row.PushBack(rec)
	}

	return
}

/*
RandomPickRows return `n` item of row that has been selected randomly from
dataset.Rows. The ids of rows that has been picked is saved id `pickedIdx`.

If duplicate is true, the row that has been picked can be picked up again,
otherwise it only allow one pick. This is also called as random selection with
or without replacement in machine learning domain.

If output mode is columns, it will be transposed to rows.
*/
func (dataset *Dataset) RandomPickRows(n int, duplicate bool) (
	picked Dataset,
	unpicked Dataset,
	pickedIdx []int,
	unpickedIdx []int,
	e error,
) {
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeColumns {
		dataset.TransposeToRows()
	}

	picked.Init(dataset.Mode, dataset.GetColumnsType())
	unpicked.Init(dataset.Mode, dataset.GetColumnsType())

	picked.Rows, unpicked.Rows, pickedIdx, unpickedIdx =
		dataset.Rows.RandomPick(n, duplicate)

	// switch the dataset based on original mode
	switch orgmode {
	case DatasetModeColumns:
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
		// transform the picked and unpicked set.
		e = picked.TransposeToColumns()
		if e != nil {
			return
		}
		e = unpicked.TransposeToColumns()
		if e != nil {
			return
		}

	case DatasetModeMatrix:
		// transform the picked and unpicked set.
		e = picked.TransposeToColumns()
		if e != nil {
			return
		}
		e = unpicked.TransposeToColumns()
		if e != nil {
			return
		}
	}

	return
}

/*
RandomPickColumns will select `n` column randomly from dataset and return
new dataset with picked and unpicked columns, and their column index.

If duplicate is true, column that has been pick up can be pick up again.

If dataset output mode is rows, it will transposed to columns.
*/
func (dataset *Dataset) RandomPickColumns(n int, dup bool, excludeIdx []int) (
	picked Dataset,
	unpicked Dataset,
	pickedIdx []int,
	unpickedIdx []int,
	e error,
) {
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeRows {
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
	}

	picked.Init(dataset.GetMode(), nil)
	unpicked.Init(dataset.GetMode(), nil)

	picked.Columns, unpicked.Columns, pickedIdx, unpickedIdx =
		dataset.Columns.RandomPick(n, dup, excludeIdx)

	// transpose picked and unpicked dataset based on original mode
	switch orgmode {
	case DatasetModeRows:
		dataset.TransposeToRows()
		picked.TransposeToRows()
		unpicked.TransposeToRows()
	case DatasetModeMatrix:
		picked.TransposeToRows()
		unpicked.TransposeToRows()
	}

	return
}

/*
SortColumnsByIndex will sort all columns using sorted index.
*/
func (dataset *Dataset) SortColumnsByIndex(sortedIdx []int) (e error) {
	if dataset.Mode == DatasetModeRows {
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
	}

	for i, col := range (*dataset).Columns {
		(*dataset).Columns[i].Records = SortRecordsByIndex(col.Records,
			sortedIdx)
	}

	return
}

/*
SplitRowsByNumeric will split the data using splitVal in column `colidx`.

For example, given two continuous attribute,

	A: {1,2,3,4}
	B: {5,6,7,8}

if colidx is (1) B and splitVal is 7, the data will splitted into left set

	A': {1,2}
	B': {5,6}

and right set

	A'': {3,4}
	B'': {7,8}
*/
func (dataset *Dataset) SplitRowsByNumeric(colidx int, splitVal float64) (
	splitLess Dataset,
	splitGreater Dataset,
	e error,
) {
	// check type of column
	coltype, e := dataset.GetColumnTypeAt(colidx)
	if e != nil {
		return
	}

	if !(coltype == TInteger || coltype == TReal) {
		return splitLess, splitGreater, ErrInvalidColType
	}

	// should we convert the data mode back?
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeColumns {
		dataset.TransposeToRows()
	}

	glog.V(2).Infoln("dataset:", dataset)

	splitLess.Init(dataset.GetMode(), dataset.GetColumnsType())
	splitGreater.Init(dataset.GetMode(), dataset.GetColumnsType())

	for _, row := range dataset.Rows {
		if row[colidx].Float() < splitVal {
			splitLess.PushRow(row)
		} else {
			splitGreater.PushRow(row)
		}
	}

	glog.V(2).Infoln(">>> split less:", splitLess)
	glog.V(2).Infoln(">>> split greater:", splitGreater)

	switch orgmode {
	case DatasetModeColumns:
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
		e = splitLess.TransposeToColumns()
		if e != nil {
			return
		}
		e = splitGreater.TransposeToColumns()
		if e != nil {
			return
		}
	case DatasetModeMatrix:
		// do nothing, since its already filled when pushing new row.
	}

	return
}

/*
SplitRowsByCategorical will split the data using a set of split value in column
`colidx`.

For example, given two attributes,

	X: [A,B,A,B,C,D,C,D]
	Y: [1,2,3,4,5,6,7,8]

if colidx is (0) or A and split value is a set `[A,C]`, the data will splitted
into left set which contain all rows that have A or C,

	X': [A,A,C,C]
	Y': [1,3,5,7]

and the right set, excluded set, will contain all rows which is not A or C,

	X'': [B,B,D,D]
	Y'': [2,4,6,8]
*/
func (dataset *Dataset) SplitRowsByCategorical(colidx int, splitVal []string) (
	splitIn Dataset,
	splitEx Dataset,
	e error,
) {
	// check type of column
	coltype, e := dataset.GetColumnTypeAt(colidx)
	if e != nil {
		return
	}

	if coltype != TString {
		return splitIn, splitEx, ErrInvalidColType
	}

	// should we convert the data mode back?
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeColumns {
		dataset.TransposeToRows()
	}

	splitIn.Init(dataset.GetMode(), dataset.GetColumnsType())
	splitEx.Init(dataset.GetMode(), dataset.GetColumnsType())

	found := false

	for _, row := range dataset.Rows {
		found = false
		for _, val := range splitVal {
			if row[colidx].String() == val {
				splitIn.PushRow(row)
				found = true
				break
			}
		}
		if !found {
			splitEx.PushRow(row)
		}
	}

	// convert all dataset based on original
	switch orgmode {
	case DatasetModeColumns:
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
		e = splitIn.TransposeToColumns()
		if e != nil {
			return
		}
		e = splitEx.TransposeToColumns()
		if e != nil {
			return
		}
	case DatasetModeMatrix:
		e = splitIn.TransposeToColumns()
		if e != nil {
			return
		}
		e = splitEx.TransposeToColumns()
		if e != nil {
			return
		}
	}

	return
}

/*
SplitRowsByValue generic function to split data by value. This function will
split data using value in column `colidx`. If value is numeric it will return
any rows that have column value less than `value` in `splitL`, and any column
value greater or equal to `value` in `splitR`.
*/
func (dataset *Dataset) SplitRowsByValue(colidx int, value interface{}) (
	splitL Dataset,
	splitR Dataset,
	e error,
) {
	coltype, e := dataset.GetColumnTypeAt(colidx)
	if e != nil {
		return
	}

	if coltype == TString {
		splitL, splitR, e = dataset.SplitRowsByCategorical(colidx,
							value.([]string))
	} else {
		var splitval float64

		switch value.(type) {
		case int:
			splitval = float64(value.(int))
		case int64:
			splitval = float64(value.(int64))
		case float32:
			splitval = float64(value.(float32))
		case float64:
			splitval = value.(float64)
		}

		splitL, splitR, e = dataset.SplitRowsByNumeric(colidx,
								splitval)
	}

	if e != nil {
		return Dataset{}, Dataset{}, e
	}

	return
}

/*
SelectColumnsByIdx return new dataset with selected column index.
*/
func (dataset *Dataset) SelectColumnsByIdx(colsIdx []int) (
	newset Dataset,
	e error,
) {
	var col *Column

	orgmode := dataset.GetMode()

	if orgmode == DatasetModeRows {
		e = dataset.TransposeToColumns()
		if e != nil {
			return
		}
	}

	newset.Init(dataset.GetMode(), nil)

	for _, idx := range colsIdx {
		col, e = dataset.GetColumn(idx)
		if e != nil {
			return
		}

		e = newset.PushColumn(*col)
		if e != nil {
			return
		}
	}

	// revert the mode back
	switch orgmode {
	case DatasetModeRows:
		dataset.TransposeToRows()
		newset.TransposeToRows()
	case DatasetModeColumns:
		// do nothing
	case DatasetModeMatrix:
		// do nothing
	}

	return
}

/*
String pretty print the data.
*/
func (dataset Dataset) String() (s string) {
	s = fmt.Sprintf("{\n"+
		"\tMode   : %v\n"+
		"\tRows   : %v\n"+
		"\tColumns: %v\n"+
	"}", dataset.Mode, dataset.Rows, dataset.Columns)
	return
}
