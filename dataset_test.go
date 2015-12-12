package dsv_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"testing"
)

var dataset_test = [][]string{
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

var dataset_type = []int{
	dsv.TInteger,
	dsv.TReal,
	dsv.TString,
}

func CreateDataset(t *testing.T) (dataset *dsv.Dataset) {
	dataset = dsv.NewDataset(dsv.DatasetModeRows)

	for _, rowin := range dataset_test {
		row := make(dsv.Row, len(rowin))

		for x, recin := range rowin {
			rec, e := dsv.NewRecord([]byte(recin), dataset_type[x])
			if e != nil {
				t.Fatal(e)
			}

			row[x] = rec
		}

		dataset.PushRow(row)
	}

	dataset.SetNColumn(len(dataset_type))

	return
}

func SplitNumericCompare(t *testing.T, splitL, splitR *dsv.Dataset, maxidx int) {
	exp := ""
	x := 0
	for ; x <= maxidx; x++ {
		exp += fmt.Sprint(dataset_test[x])
	}
	got := fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	exp = ""
	for ; x < len(dataset_test); x++ {
		exp += fmt.Sprint(dataset_test[x])
	}
	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)
}

func TestSplitRowsByNumeric(t *testing.T) {
	dataset := CreateDataset(t)

	// Split integer by float
	splitL, splitR, e := dataset.SplitRowsByNumeric(0, 4.5)
	if e != nil {
		t.Fatal(e)
	}

	SplitNumericCompare(t, splitL, splitR, 4)

	// Split by float
	splitL, splitR, e = dataset.SplitRowsByNumeric(1, 1.8)
	if e != nil {
		t.Fatal(e)
	}

	SplitNumericCompare(t, splitL, splitR, 7)
}
