package dsv_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"testing"
)

var datasetRows = [][]string{
	{"0", "1", "A"},
	{"1", "1.1", "B"},
	{"2", "1.2", "A"},
	{"3", "1.3", "B"},
	{"4", "1.4", "C"},
	{"5", "1.5", "D"},
	{"6", "1.6", "C"},
	{"7", "1.7", "D"},
	{"8", "1.8", "E"},
	{"9", "1.9", "F"},
}

var datasetCols = [][]string{
	{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
	{"1", "1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9"},
	{"A", "B", "A", "B", "C", "D", "C", "D", "E", "F"},
}

var datasetTypes = []int{
	dsv.TInteger,
	dsv.TReal,
	dsv.TString,
}

var datasetNames = []string{"int", "real", "string"}

func PopulateWithRows(t *testing.T, dataset *dsv.Dataset) {
	for _, rowin := range datasetRows {
		row := make(dsv.Row, len(rowin))

		for x, recin := range rowin {
			rec, e := dsv.NewRecord(recin, datasetTypes[x])
			if e != nil {
				t.Fatal(e)
			}

			row[x] = rec
		}

		dataset.PushRow(row)
	}
}

func PopulateWithColumns(t *testing.T, dataset *dsv.Dataset) {
	for x := range datasetCols {
		col, e := dsv.NewColumnString(datasetCols[x], datasetTypes[x],
			datasetNames[x])
		if e != nil {
			t.Fatal(e)
		}

		dataset.PushColumn(*col)
	}
}

func CreateDataset(t *testing.T) (dataset *dsv.Dataset) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeRows, datasetTypes,
		datasetNames)
	if e != nil {
		t.Fatal(e)
	}

	PopulateWithRows(t, dataset)

	return
}

func DatasetStringJoinByIndex(t *testing.T, dataset [][]string, indis []int) (res string) {
	for x := range indis {
		res += fmt.Sprint(dataset[indis[x]])
	}
	return res
}

func DatasetRowsJoin(t *testing.T) (s string) {
	for x := range datasetRows {
		s += fmt.Sprint(datasetRows[x])
	}
	return
}

func DatasetColumnsJoin(t *testing.T) (s string) {
	for x := range datasetCols {
		s += fmt.Sprint(datasetCols[x])
	}
	return
}

func TestSplitRowsByNumeric(t *testing.T) {
	dataset := CreateDataset(t)

	// Split integer by float
	splitL, splitR, e := dataset.SplitRowsByNumeric(0, 4.5)
	if e != nil {
		t.Fatal(e)
	}

	expIdx := []int{0, 1, 2, 3, 4}
	exp := DatasetStringJoinByIndex(t, datasetRows, expIdx)
	got := fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	expIdx = []int{5, 6, 7, 8, 9}
	exp = DatasetStringJoinByIndex(t, datasetRows, expIdx)
	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)

	// Split by float
	splitL, splitR, e = dataset.SplitRowsByNumeric(1, 1.8)
	if e != nil {
		t.Fatal(e)
	}

	expIdx = []int{0, 1, 2, 3, 4, 5, 6, 7}
	exp = DatasetStringJoinByIndex(t, datasetRows, expIdx)
	got = fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	expIdx = []int{8, 9}
	exp = DatasetStringJoinByIndex(t, datasetRows, expIdx)
	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)
}

func TestSplitRowsByCategorical(t *testing.T) {
	dataset := CreateDataset(t)
	splitval := []string{"A", "D"}

	splitL, splitR, e := dataset.SplitRowsByCategorical(2, splitval)
	if e != nil {
		t.Fatal(e)
	}

	expIdx := []int{0, 2, 5, 7}
	exp := DatasetStringJoinByIndex(t, datasetRows, expIdx)
	got := fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	expIdx = []int{1, 3, 4, 6, 8, 9}
	exp = DatasetStringJoinByIndex(t, datasetRows, expIdx)
	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)
}

func TestModeColumnsPushColumn(t *testing.T) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeColumns, nil, nil)

	if e != nil {
		t.Fatal(e)
	}

	exp := ""
	got := ""
	for x := range datasetCols {
		col, e := dsv.NewColumnString(datasetCols[x], datasetTypes[x],
			datasetNames[x])
		if e != nil {
			t.Fatal(e)
		}

		dataset.PushColumn(*col)

		exp += fmt.Sprint(datasetCols[x])
		got += fmt.Sprint(dataset.Columns[x].Records)
	}

	assert.Equal(t, exp, got)

	// Check rows
	exp = ""
	got = fmt.Sprint(dataset.Rows)
	assert.Equal(t, exp, got)
}

func TestModeRowsPushColumn(t *testing.T) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeRows, nil, nil)

	if e != nil {
		t.Fatal(e)
	}

	PopulateWithColumns(t, dataset)

	// Check rows
	exp := DatasetRowsJoin(t)
	got := fmt.Sprint(dataset.Rows)

	assert.Equal(t, exp, got)

	// Check columns
	exp = "[{int 1 0 [] []} {real 2 0 [] []} {string 0 0 [] []}]"
	got = fmt.Sprint(dataset.Columns)

	assert.Equal(t, exp, got)
}

func TestModeMatrixPushColumn(t *testing.T) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeMatrix, nil, nil)

	if e != nil {
		t.Fatal(e)
	}

	exp := ""
	got := ""
	for x := range datasetCols {
		col, e := dsv.NewColumnString(datasetCols[x], datasetTypes[x],
			datasetNames[x])
		if e != nil {
			t.Fatal(e)
		}

		dataset.PushColumn(*col)

		exp += fmt.Sprint(datasetCols[x])
		got += fmt.Sprint(dataset.Columns[x].Records)
	}

	assert.Equal(t, exp, got)

	// Check rows
	exp = DatasetRowsJoin(t)
	got = fmt.Sprint(dataset.Rows)

	assert.Equal(t, exp, got)
}

func TestModeRowsPushRows(t *testing.T) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeRows, nil, nil)

	if e != nil {
		t.Fatal(e)
	}

	PopulateWithRows(t, dataset)

	exp := DatasetRowsJoin(t)
	got := fmt.Sprint(dataset.Rows)

	assert.Equal(t, exp, got)
}

func TestModeColumnsPushRows(t *testing.T) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeColumns, nil, nil)

	if e != nil {
		t.Fatal(e)
	}

	PopulateWithRows(t, dataset)

	// check rows
	exp := ""
	got := fmt.Sprint(dataset.Rows)

	assert.Equal(t, exp, got)

	// check columns
	exp = DatasetColumnsJoin(t)
	got = ""
	for x := range dataset.Columns {
		got += fmt.Sprint(dataset.Columns[x].Records)
	}

	assert.Equal(t, exp, got)
}

func TestModeMatrixPushRows(t *testing.T) {
	dataset, e := dsv.NewDataset(dsv.DatasetModeMatrix, nil, nil)

	if e != nil {
		t.Fatal(e)
	}

	PopulateWithRows(t, dataset)

	exp := DatasetRowsJoin(t)
	got := fmt.Sprint(dataset.Rows)

	assert.Equal(t, exp, got)

	// check columns
	exp = DatasetColumnsJoin(t)
	got = ""
	for x := range dataset.Columns {
		got += fmt.Sprint(dataset.Columns[x].Records)
	}

	assert.Equal(t, exp, got)
}
