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
	dataset, e := dsv.NewDataset(dsv.DatasetModeRows, dataset_type)
	if e != nil {
		t.Fatal(e)
	}

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

	return
}

func DatasetStringJoinByIndex(t *testing.T, dataset [][]string, indis []int) (res string) {
	for x := range indis {
		res += fmt.Sprint(dataset[indis[x]])
	}
	return res
}

func TestSplitRowsByNumeric(t *testing.T) {
	dataset := CreateDataset(t)

	// Split integer by float
	splitL, splitR, e := dataset.SplitRowsByNumeric(0, 4.5)
	if e != nil {
		t.Fatal(e)
	}

	exp_idx := []int{0, 1, 2, 3, 4}
	exp := DatasetStringJoinByIndex(t, dataset_test, exp_idx)
	got := fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	exp_idx = []int{5, 6, 7, 8, 9}
	exp = DatasetStringJoinByIndex(t, dataset_test, exp_idx)
	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)

	// Split by float
	splitL, splitR, e = dataset.SplitRowsByNumeric(1, 1.8)
	if e != nil {
		t.Fatal(e)
	}

	exp_idx = []int{0, 1, 2, 3, 4, 5, 6, 7}
	exp = DatasetStringJoinByIndex(t, dataset_test, exp_idx)
	got = fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	exp_idx = []int{8, 9}
	exp = DatasetStringJoinByIndex(t, dataset_test, exp_idx)
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

	exp_idx := []int{0, 2, 5, 7}
	exp := DatasetStringJoinByIndex(t, dataset_test, exp_idx)
	got := fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	exp_idx = []int{1, 3, 4, 6, 8, 9}
	exp = DatasetStringJoinByIndex(t, dataset_test, exp_idx)
	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)
}
