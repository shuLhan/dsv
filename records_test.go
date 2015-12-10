package dsv_test

import (
	"fmt"
	"testing"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
)

func TestSortByIndex(t *testing.T) {
	data := make([]*dsv.Record, 3)
	data[0],_ = dsv.NewRecord([]byte{'3'}, dsv.TInteger)
	data[1],_ = dsv.NewRecord([]byte{'2'}, dsv.TInteger)
	data[2],_ = dsv.NewRecord([]byte{'1'}, dsv.TInteger)

	sortedIdx := []int{2, 1, 0,}
	expect := []int{1, 2, 3,}

	sorted := dsv.SortRecordsByIndex(data, sortedIdx)

	got := fmt.Sprint(sorted)
	exp := fmt.Sprint(expect)

	assert.Equal(t, exp, got)
}
