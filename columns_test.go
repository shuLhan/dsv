package dsv_test

import (
	_ "fmt"
	"testing"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
)

var data = []string{"9.987654321", "8.8", "7.7", "6.6", "5.5", "4.4", "3.3"}
var exp_float = []float64{9.987654321, 8.8, 7.7, 6.6, 5.5, 4.4, 3.3}

func TestToFloatSlice(t *testing.T) {
	col := make(dsv.Column, len(data))

	for x := range data {
		col[x], _ = dsv.NewRecord([]byte(data[x]), dsv.TReal)
	}

	got := col.ToFloatSlice()

	assert.Equal(t, exp_float, got)
}

func TestToStringSlice(t *testing.T) {
	col := make(dsv.Column, len(data))

	for x := range data {
		col[x], _ = dsv.NewRecord([]byte(data[x]), dsv.TString)
	}

	got := col.ToStringSlice()

	assert.Equal(t, data, got)
}
