// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"fmt"
	"testing"

	"github.com/shuLhan/dsv/util"
	"github.com/shuLhan/dsv/util/assert"
)

var input = [][]float64{
	{9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0},
	{9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0, 0.0},
	{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
	{0.0, 6.0, 7.0, 8.0, 5.0, 1.0, 2.0, 3.0, 4.0, 9.0},
	{9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0},
}

var output = [][]float64{
	{3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
	{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
	{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
	{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
	{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0},
}

var expSortedIdx = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var expReverseSortedIdx = []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

func TestIndirectSortFloat64(t *testing.T) {
	var res, exp string

	for i := range input {
		util.IndirectSortFloat64(input[i])

		res = fmt.Sprint(input[i])
		exp = fmt.Sprint(output[i])

		assert.Equal(t, exp, res)
	}
}

func TestIndirectSortFloat642(t *testing.T) {
	in := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	got := util.IndirectSortFloat64(in)

	assert.Equal(t, expSortedIdx, got)
}

func TestIndirectSortFloat643(t *testing.T) {
	in := []float64{5.1, 5, 5.6, 5.5, 5.5, 5.8, 5.5, 5.5, 5.8, 5.6,
		5.7, 5, 5.6, 5.9, 6.2, 6, 4.9, 6.3, 6.1, 5.6,
		5.8, 6.7, 6.1, 5.9, 6, 4.9, 5.6, 5.2, 6.1, 6.4,
		7, 5.7, 6.5, 6.9, 5.7, 6.4, 6.2, 6.6, 6.3, 6.2,
		5.4, 6.7, 6.1, 5.7, 5.5, 6, 3, 6.6, 5.7, 6,
		6.8, 6, 6.1, 6.3, 5.8, 5.8, 5.6, 5.7, 6, 6.9,
		6.9, 6.4, 6.3, 6.3, 6.7, 6.5, 5.8, 6.3, 6.4, 6.7,
		5.9, 7.2, 6.3, 6.3, 6.5, 7.1, 6.7, 7.6, 7.3, 6.4,
		6.7, 7.4, 6, 6.8, 6.5, 6.4, 6.7, 6.4, 6.5, 6.9,
		7.7, 6.7, 7.2, 7.7, 7.2, 7.7, 6.1, 7.9, 7.7, 6.8,
		6.2}

	exp := []float64{3, 4.9, 4.9, 5, 5, 5.1, 5.2, 5.4, 5.5, 5.5,
		5.5, 5.5, 5.5, 5.6, 5.6, 5.6, 5.6, 5.6, 5.6, 5.7,
		5.7, 5.7, 5.7, 5.7, 5.7, 5.8, 5.8, 5.8, 5.8, 5.8,
		5.8, 5.9, 5.9, 5.9, 6, 6, 6, 6, 6, 6,
		6, 6.1, 6.1, 6.1, 6.1, 6.1, 6.1, 6.2, 6.2, 6.2,
		6.2, 6.3, 6.3, 6.3, 6.3, 6.3, 6.3, 6.3, 6.3, 6.4,
		6.4, 6.4, 6.4, 6.4, 6.4, 6.4, 6.5, 6.5, 6.5, 6.5,
		6.5, 6.6, 6.6, 6.7, 6.7, 6.7, 6.7, 6.7, 6.7, 6.7,
		6.7, 6.8, 6.8, 6.8, 6.9, 6.9, 6.9, 6.9, 7, 7.1,
		7.2, 7.2, 7.2, 7.3, 7.4, 7.6, 7.7, 7.7, 7.7, 7.7,
		7.9}

	util.IndirectSortFloat64(in)

	assert.Equal(t, exp, in)
}

func TestSortFloatSliceByIndex(t *testing.T) {
	in1 := []float64{9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0, 0.0}
	in2 := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0}

	exp := fmt.Sprint(in1)

	sortedIdx := util.IndirectSortFloat64(in1)

	assert.Equal(t, expReverseSortedIdx, sortedIdx)

	util.SortFloatSliceByIndex(&in2, &sortedIdx)

	got := fmt.Sprint(in2)

	assert.Equal(t, exp, got)
}

func TestCountMissRate(t *testing.T) {
	data := []string{"A", "B", "C", "D"}
	test := []string{"A", "B", "C", "D"}
	repl := []string{"B", "C", "D", "E"}
	exps := []float64{0.25, 0.5, 0.75, 1.0}

	got, _, _ := util.CountMissRate(data, test)
	assert.Equal(t, 0.0, got)

	for x, exp := range exps {
		test[x] = repl[x]
		got, _, _ = util.CountMissRate(data, test)
		assert.Equal(t, exp, got)
	}
}
