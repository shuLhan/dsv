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
	// DatasetNoMode default to matrix.
	DatasetNoMode = 0
	// DatasetModeRows for output mode in rows.
	DatasetModeRows = 1
	// DatasetModeColumns for output mode in columns.
	DatasetModeColumns = 2
	// DatasetModeMatrix will save data in rows and columns.
	DatasetModeMatrix = 4
)

var (
	// ErrColIdxOutOfRange operation on column index is invalid
	ErrColIdxOutOfRange = errors.New("dsv: Column index out of range")
	// ErrInvalidColType operation on column with different type
	ErrInvalidColType = errors.New("dsv: Invalid column type")
	// ErrEmptySet operation on empty data set.
	ErrEmptySet = errors.New("dsv: dataset is empty")
	// ErrMisColLength returned when operation on columns does not match
	// between parameter and their length
	ErrMisColLength = errors.New("dsv: mismatch on column length")
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
func NewDataset(mode int, types []int, names []string) (
	dataset *Dataset,
) {
	dataset = &Dataset{}

	dataset.Init(mode, types, names)

	return
}

/*
Clone return a copy of current dataset.
*/
func (dataset *Dataset) Clone() (clone Dataset) {
	clone.SetMode(dataset.GetMode())

	for _, col := range dataset.Columns {
		newcol := Column{
			Type:       col.Type,
			Name:       col.Name,
			ValueSpace: col.ValueSpace,
		}
		clone.PushColumn(newcol)
	}
	return
}

/*
Init will set the dataset using mode and types.
*/
func (dataset *Dataset) Init(mode int, types []int, names []string) {
	if types == nil {
		dataset.Columns = make(Columns, 0)
	} else {
		dataset.Columns = make(Columns, len(types))
		dataset.Columns.SetType(types)
	}

	dataset.SetColumnsName(names)
	dataset.SetMode(mode)
}

/*
Reset all data and attributes.
*/
func (dataset *Dataset) Reset() error {
	dataset.Rows = Rows{}
	dataset.Columns.Reset()
	return nil
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
func (dataset *Dataset) SetMode(mode int) {
	switch mode {
	case DatasetModeRows:
		dataset.Mode = DatasetModeRows
		dataset.Rows = make(Rows, 0)
	case DatasetModeColumns:
		dataset.Mode = DatasetModeColumns
		dataset.Columns.Reset()
	case DatasetModeMatrix, DatasetNoMode:
		dataset.Mode = DatasetModeMatrix
		dataset.Rows = make(Rows, 0)
		dataset.Columns.Reset()
	}
	dataset.Mode = mode
}

/*
GetNColumn return the number of column in dataset.
*/
func (dataset *Dataset) GetNColumn() (ncol int) {
	ncol = len(dataset.Columns)

	if ncol > 0 {
		return
	}

	switch dataset.Mode {
	case DatasetModeRows:
		if len(dataset.Rows) <= 0 {
			return 0
		}
		return len(dataset.Rows[0])
	}

	return
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
			nrow = dataset.Columns[0].Len()
		}
	case DatasetModeMatrix, DatasetNoMode:
		// matrix mode could have empty either in rows or column.
		nrow = len(dataset.Rows)
	}
	return
}

/*
Len return number of row in dataset.
*/
func (dataset *Dataset) Len() int {
	return dataset.GetNRow()
}

/*
SetColumnsType of data in all columns.
*/
func (dataset *Dataset) SetColumnsType(types []int) {
	dataset.Columns = make(Columns, len(types))
	dataset.Columns.SetType(types)
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
GetColumnsName return name of all columns.
*/
func (dataset *Dataset) GetColumnsName() (names []string) {
	for x := range dataset.Columns {
		names = append(names, dataset.Columns[x].GetName())
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
If `idx` is out of range return nil.
*/
func (dataset *Dataset) GetColumn(idx int) (col *Column) {
	if idx > dataset.GetNColumn() {
		return
	}

	switch dataset.Mode {
	case DatasetModeRows:
		dataset.TransposeToColumns()
	case DatasetModeColumns:
		// do nothing
	case DatasetModeMatrix:
		// do nothing
	}

	return &dataset.Columns[idx]
}

/*
GetRow return row at index `idx`.
*/
func (dataset *Dataset) GetRow(idx int) *Row {
	return &dataset.Rows[idx]
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
	case DatasetModeMatrix, DatasetNoMode:
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
func (dataset *Dataset) GetDataAsColumns() (columns Columns) {
	if dataset.Mode == DatasetModeRows {
		dataset.TransposeToColumns()
	}
	return dataset.Columns
}

/*
TransposeToColumns move all data from rows (horizontal) to columns
(vertical) mode.
*/
func (dataset *Dataset) TransposeToColumns() {
	if dataset.GetNRow() <= 0 {
		// nothing to transpose
		return
	}

	ncol := dataset.GetNColumn()
	if ncol <= 0 {
		// if no columns defined, initialize it using record type
		// in the first row.
		types := dataset.GetRow(0).GetTypes()
		dataset.SetColumnsType(types)
		ncol = len(types)
	}

	orgmode := dataset.GetMode()

	switch orgmode {
	case DatasetModeRows:
		// do nothing.
	case DatasetModeColumns, DatasetModeMatrix, DatasetNoMode:
		// check if column records contain data.
		nrow := dataset.Columns[0].Len()
		if nrow > 0 {
			// return if column record is not empty, its already
			// transposed
			return
		}
	}

	// use the least length
	minlen := len(*dataset.GetRow(0))

	if minlen > ncol {
		minlen = ncol
	}

	switch orgmode {
	case DatasetModeRows, DatasetNoMode:
		dataset.SetMode(DatasetModeColumns)
	}

	for _, row := range dataset.Rows {
		for y := 0; y < minlen; y++ {
			dataset.Columns[y].PushBack(row[y])
		}
	}

	// reset the rows data only if original mode is rows
	// this to prevent empty data when mode is matrix.
	switch orgmode {
	case DatasetModeRows:
		dataset.Rows = nil
	}
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
		l := dataset.Columns[f].Len()

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
func (dataset *Dataset) PushRow(row Row) {
	switch dataset.GetMode() {
	case DatasetModeRows:
		dataset.Rows = append(dataset.Rows, row)
	case DatasetModeColumns:
		dataset.PushRowToColumns(row)
	case DatasetModeMatrix, DatasetNoMode:
		dataset.Rows = append(dataset.Rows, row)
		dataset.PushRowToColumns(row)
	}
}

/*
PushRowToColumns push each data in Row to Columns.
*/
func (dataset *Dataset) PushRowToColumns(row Row) {
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
}

/*
PushColumn will append new column to the end of slice.
*/
func (dataset *Dataset) PushColumn(col Column) {
	switch dataset.GetMode() {
	case DatasetModeRows:
		dataset.Columns = append(dataset.Columns, col)
		dataset.PushColumnToRows(col)
		// Remove records in column
		dataset.Columns[dataset.GetNColumn()-1].Reset()
	case DatasetModeColumns:
		dataset.Columns = append(dataset.Columns, col)
	case DatasetModeMatrix, DatasetNoMode:
		dataset.Columns = append(dataset.Columns, col)
		dataset.PushColumnToRows(col)
	}
}

/*
PushColumnToRows add each record in column to each rows, from top to bottom.
*/
func (dataset *Dataset) PushColumnToRows(col Column) {
	colsize := col.Len()
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
) {
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeColumns {
		dataset.TransposeToRows()
	}

	picked = dataset.Clone()
	unpicked = dataset.Clone()

	picked.Rows, unpicked.Rows, pickedIdx, unpickedIdx =
		dataset.Rows.RandomPick(n, duplicate)

	// switch the dataset based on original mode
	switch orgmode {
	case DatasetModeColumns:
		dataset.TransposeToColumns()
		// transform the picked and unpicked set.
		picked.TransposeToColumns()
		unpicked.TransposeToColumns()

	case DatasetModeMatrix, DatasetNoMode:
		// transform the picked and unpicked set.
		picked.TransposeToColumns()
		unpicked.TransposeToColumns()
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
) {
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeRows {
		dataset.TransposeToColumns()
	}

	picked.Init(dataset.GetMode(), nil, nil)
	unpicked.Init(dataset.GetMode(), nil, nil)

	picked.Columns, unpicked.Columns, pickedIdx, unpickedIdx =
		dataset.Columns.RandomPick(n, dup, excludeIdx)

	// transpose picked and unpicked dataset based on original mode
	switch orgmode {
	case DatasetModeRows:
		dataset.TransposeToRows()
		picked.TransposeToRows()
		unpicked.TransposeToRows()
	case DatasetModeMatrix, DatasetNoMode:
		picked.TransposeToRows()
		unpicked.TransposeToRows()
	}

	return
}

/*
SortColumnsByIndex will sort all columns using sorted index.
*/
func (dataset *Dataset) SortColumnsByIndex(sortedIdx []int) {
	if dataset.Mode == DatasetModeRows {
		dataset.TransposeToColumns()
	}

	for i, col := range (*dataset).Columns {
		(*dataset).Columns[i].Records = SortRecordsByIndex(col.Records,
			sortedIdx)
	}
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

	splitLess = dataset.Clone()
	splitGreater = dataset.Clone()

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
		dataset.TransposeToColumns()
		splitLess.TransposeToColumns()
		splitGreater.TransposeToColumns()
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

	splitIn = dataset.Clone()
	splitEx = dataset.Clone()

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
		dataset.TransposeToColumns()
		splitIn.TransposeToColumns()
		splitEx.TransposeToColumns()
	case DatasetModeMatrix, DatasetNoMode:
		splitIn.TransposeToColumns()
		splitEx.TransposeToColumns()
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
) {
	var col *Column

	orgmode := dataset.GetMode()

	if orgmode == DatasetModeRows {
		dataset.TransposeToColumns()
	}

	newset.Init(dataset.GetMode(), nil, nil)

	for _, idx := range colsIdx {
		col = dataset.GetColumn(idx)
		if col == nil {
			continue
		}

		newset.PushColumn(*col)
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
SelectRowsWhere return all rows which column value in `colidx` is equal
to `colval`.
*/
func (dataset *Dataset) SelectRowsWhere(colidx int, colval string) (
	selected Dataset,
) {
	orgmode := dataset.GetMode()

	if orgmode == DatasetModeColumns {
		dataset.TransposeToRows()
	}

	selected.Init(dataset.GetMode(), nil, nil)

	selected.Rows = dataset.Rows.SelectWhere(colidx, colval)

	switch orgmode {
	case DatasetModeColumns:
		dataset.TransposeToColumns()
		selected.TransposeToColumns()
	case DatasetModeMatrix, DatasetNoMode:
		selected.TransposeToColumns()
	}

	return
}

/*
MergeColumns append columns from other dataset into current dataset.
*/
func (dataset *Dataset) MergeColumns(other Dataset) {
	cols := other.GetDataAsColumns()
	for _, col := range cols {
		dataset.PushColumn(col)
	}
}

/*
MergeRows append rows from other dataset into current dataset.
*/
func (dataset *Dataset) MergeRows(other Dataset) {
	rows := other.GetDataAsRows()
	for _, row := range rows {
		dataset.PushRow(row)
	}
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
