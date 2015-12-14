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
func NewDataset(mode int) (dataset *Dataset) {
	dataset = &Dataset{
		Mode:    mode,
		NRow:    0,
		NColumn: 0,
		Columns: nil,
		Rows:    nil,
	}

	dataset.SetMode(mode)

	return
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
GetColumnType return type of column in index `colidx` in
dataset.
*/
func (dataset *Dataset) GetColumnType(colidx int) (int, error) {
	if dataset.NRow <= 0 {
		return TUndefined, ErrEmptySet
	}
	if colidx >= dataset.GetNColumn() {
		return TUndefined, ErrColIdxOutOfRange
	}

	var record *Record

	if dataset.Mode == DatasetModeRows || dataset.Mode == DatasetModeMatrix {
		record = dataset.Rows[0][colidx]
	} else if dataset.Mode == DatasetModeColumns {
		record = dataset.Columns[colidx][0]
	} else {
		return TUndefined, ErrUnknownDatasetMode
	}

	return record.GetType()
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
dataset.Rows. The ids of rows that has been picked is saved id `rowsIdx`.

If duplicate is true, the row that has been picked can be picked up again,
otherwise it only allow one pick. This is also called as random selection with
or without replacement in some machine learning domain.

If output mode is columns, it will be transposed to rows.
*/
func (dataset *Dataset) RandomPickRows(n int, duplicate bool) (
	unpicked Rows,
	shuffled Rows,
	pickedIdx []int,
) {
	if dataset.GetMode() == DatasetModeColumns {
		dataset.TransposeToRows()
	}
	return dataset.Rows.RandomPick(n, duplicate)
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
	splitLess *Dataset,
	splitGreater *Dataset,
	e error,
) {
	// check type of column
	coltype, e := dataset.GetColumnType(colidx)
	if e != nil {
		return nil, nil, e
	}

	if !(coltype == TInteger || coltype == TReal) {
		return nil, nil, ErrInvalidColType
	}

	// should we convert the data mode back?
	modeIsColumns := false

	if dataset.Mode == DatasetModeColumns {
		modeIsColumns = true
		dataset.TransposeToRows()
	}

	splitLess = NewDataset(DatasetModeRows)
	splitGreater = NewDataset(DatasetModeRows)

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
	splitIn *Dataset,
	splitEx *Dataset,
	e error,
) {
	// check type of column
	coltype, e := dataset.GetColumnType(colidx)
	if e != nil {
		return nil, nil, e
	}

	if coltype != TString {
		return nil, nil, ErrInvalidColType
	}

	// should we convert the data mode back?
	modeIsColumns := false

	if dataset.Mode == DatasetModeColumns {
		modeIsColumns = true
		dataset.TransposeToRows()
	}

	splitIn = NewDataset(DatasetModeRows)
	splitEx = NewDataset(DatasetModeRows)
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
