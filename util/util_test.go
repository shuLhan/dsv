// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util_test

import (
	"fmt"
	"testing"

	"github.com/shuLhan/dsv/util"
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

func TestIndirectSort(t *testing.T) {
	var res, exp string

	for i := range input {
		util.IndirectSort(&input[i])

		res = fmt.Sprint(input[i])
		exp = fmt.Sprint(output[i])

		if res != exp {
			t.Fatal("Expecting:\n", exp, "\n Got:", res)
		}
	}
}

func TestSortFloatSliceByIndex(t *testing.T) {
	in1 := []float64{9.0, 8.0, 7.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1.0, 0.0}
	in2 := []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0}

	exp := fmt.Sprint(in1)

	fmt.Println("input 1:", in1)
	fmt.Println("input 2:", in2)

	sortedIdx := util.IndirectSort(&in1)

	fmt.Println("sorted idx:", sortedIdx)

	util.SortFloatSliceByIndex(&in2, sortedIdx)

	fmt.Println("input 1:", in1)
	fmt.Println("input 2:", in2)

	got := fmt.Sprint(in2)

	if got != exp {
		t.Fatal("Expecting ", exp,", got ", got)
	}
}
