// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"errors"
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
	// ColumnType define the type of data in column
	ColumnType []int
	// NRow define number of rows.
	NRow int
	// NColumn define number of columns.
	NColumn int
	// Columns is input data that has been parsed.
	Columns Columns
	// Rows is input data that has been parsed.
	Rows Rows
}

/*
NewDataset create new dataset, use the mode to initialize the dataset.
*/
func NewDataset(mode int, types []int) (dataset *Dataset) {
	dataset = &Dataset{
		Mode:    mode,
		ColumnType: types,
		NRow:    0,
		NColumn: len(types),
		Columns: nil,
		Rows:    nil,
	}

	dataset.SetMode(mode)

	return
}

/*
Init will set the dataset using mode and types.
*/
func (dataset *Dataset) Init(mode int, types []int) {
	dataset.SetMode(mode)
	dataset.NColumn = len(types)
	dataset.ColumnType = types
	dataset.NRow = 0
}

/*
Reset all data and attributes.
*/
func (dataset *Dataset) Reset() {
	dataset.NRow = 0
	dataset.Rows = Rows{}
	dataset.Columns = make(Columns, dataset.NColumn)
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
		dataset.NRow = 0
		dataset.Mode = DatasetModeRows
		dataset.Rows = Rows{}
	case DatasetModeColumns:
		dataset.Mode = DatasetModeColumns
		dataset.Columns = make(Columns, dataset.NColumn)
	case DatasetModeMatrix:
		dataset.Mode = DatasetModeMatrix
		dataset.Rows = Rows{}
		dataset.Columns = make(Columns, dataset.NColumn)
	default:
		return ErrUnknownDatasetMode
	}
	dataset.Mode = mode

	return nil
}

/*
GetNColumn return number of column that will be used in output, excluding
the column with Skip=true.
*/
func (dataset *Dataset) GetNColumn() int {
	return dataset.NColumn
}

/*
SetNColumn set number of output columns.
*/
func (dataset *Dataset) SetNColumn(n int) {
	dataset.NColumn = n
}

/*
GetNRow return number of rows in dataset.
*/
func (dataset *Dataset) GetNRow() int {
	return dataset.NRow
}

/*
SetNRow will set the number of rows in dataset.
*/
func (dataset *Dataset) SetNRow(n int) {
	dataset.NRow = n
}

/*
SetColumnType of data in all columns.
*/
func (dataset *Dataset) SetColumnType(types []int) {
	dataset.ColumnType = types
	dataset.NColumn = len(types)
}

/*
GetColumnTypeAt return type of column in index `colidx` in
dataset.
*/
func (dataset *Dataset) GetColumnTypeAt(colidx int) (int, error) {
	if colidx >= dataset.GetNColumn() {
		return TUndefined, ErrColIdxOutOfRange
	}

	return dataset.ColumnType[colidx], nil
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
func (dataset *Dataset) GetDataAsColumns() Columns {
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
	toutmode := dataset.GetMode()
	if toutmode == DatasetModeColumns || toutmode == DatasetModeMatrix {
		return
	}
	if dataset.NRow <= 0 {
		// do nothing ...
		return
	}

	// double check column length
	collen := len(dataset.Rows[0])

	if collen > dataset.NColumn {
		dataset.NColumn = collen
	}

	dataset.SetMode(DatasetModeColumns)

	for i := range dataset.Rows {
		dataset.PushRowToColumns(dataset.Rows[i])
	}

	// reset the rows
	dataset.Rows = nil
}

/*
TransposeToRows will move all data from columns (vertical) to rows (horizontal)
mode.
*/
func (dataset *Dataset) TransposeToRows() {
	toutmode := dataset.GetMode()
	if toutmode == DatasetModeRows || toutmode == DatasetModeMatrix {
		return
	}

	rowlen := math.MaxInt32
	flen := len(dataset.Columns)

	dataset.SetMode(DatasetModeRows)

	// Get the least length of columns.
	for f := 0; f < flen; f++ {
		l := len(dataset.Columns[f])

		if l < rowlen {
			rowlen = l
		}
	}

	for r := 0; r < rowlen; r++ {
		row := make(Row, flen)

		for f := 0; f < flen; f++ {
			row[f] = dataset.Columns[f][r]
		}

		dataset.PushRow(row)
	}

	// reset the columns
	dataset.Columns = nil
}

/*
PushRow save the data, which is already in row object, to Rows.
*/
func (dataset *Dataset) PushRow(r Row) {
	dataset.Rows = append(dataset.Rows, r)
	dataset.NRow++
}

/*
PushRowToColumns push each data in Row to Columns.
*/
func (dataset *Dataset) PushRowToColumns(row Row) (e error) {
	// check if row length equal with columns length
	if len(row) != len(dataset.Columns) {
		return ErrMissRecordsLen
	}

	for i := range row {
		dataset.Columns[i] = append(dataset.Columns[i], row[i])
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
) {
	if dataset.GetMode() == DatasetModeColumns {
		dataset.TransposeToRows()
	}

	picked.Init(dataset.Mode, dataset.ColumnType)
	unpicked.Init(dataset.Mode, dataset.ColumnType)

	picked.Rows, unpicked.Rows, pickedIdx, unpickedIdx =
		dataset.Rows.RandomPick(n, duplicate)

	picked.NRow = len(picked.Rows)
	unpicked.NRow = len(unpicked.Rows)

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
	if dataset.GetMode() == DatasetModeRows {
		dataset.TransposeToColumns()
	}

	picked.Init(dataset.Mode, dataset.ColumnType)
	unpicked.Init(dataset.Mode, dataset.ColumnType)

	picked.Columns, unpicked.Columns, pickedIdx, unpickedIdx =
		dataset.Columns.RandomPick(n, dup, excludeIdx)

	picked.NRow = dataset.NRow
	unpicked.NRow = dataset.NRow

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
		(*dataset).Columns[i] = SortRecordsByIndex(col, sortedIdx)
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
	modeIsColumns := false

	if dataset.Mode == DatasetModeColumns {
		modeIsColumns = true
		dataset.TransposeToRows()
	}

	splitLess.Init(DatasetModeRows, dataset.ColumnType)
	splitGreater.Init(DatasetModeRows, dataset.ColumnType)

	for _, row := range dataset.Rows {
		if row[colidx].Float() < splitVal {
			splitLess.PushRow(row)
		} else {
			splitGreater.PushRow(row)
		}
	}

	if modeIsColumns {
		dataset.TransposeToColumns()
		splitLess.TransposeToColumns()
		splitGreater.TransposeToColumns()
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
	modeIsColumns := false

	if dataset.Mode == DatasetModeColumns {
		modeIsColumns = true
		dataset.TransposeToRows()
	}

	splitIn.Init(DatasetModeRows, dataset.ColumnType)
	splitEx.Init(DatasetModeRows, dataset.ColumnType)

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

	// transpose original dataset back to columns
	if modeIsColumns {
		dataset.TransposeToColumns()
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
