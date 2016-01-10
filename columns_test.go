// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"testing"
)

func TestRandomPickColumns(t *testing.T) {
	var dataset dsv.Dataset
	var e error

	dataset.Init(dsv.DatasetModeRows, testColTypes, testColNames)

	dataset.Rows, e = initRows()
	if e != nil {
		t.Fatal(e)
	}

	dataset.TransposeToColumns()

	// random pick with duplicate
	ncols := 6
	dup := true
	excludeIdx := []int{3}

	for i := 0; i < 5; i++ {
		picked, unpicked, pickedIdx, unpickedIdx :=
			dataset.Columns.RandomPick(ncols, dup, excludeIdx)

		// check if unpicked item exist in picked items.
		for _, un := range unpicked {
			for _, pick := range picked {
				assert.NotEqual(t, un, pick)
			}
		}

		fmt.Println("Random pick with duplicate columns")
		fmt.Println("==> picked columns   :", picked)
		fmt.Println("==> picked idx       :", pickedIdx)
		fmt.Println("==> unpicked columns :", unpicked)
		fmt.Println("==> unpicked idx     :", unpickedIdx)
	}

	// random pick without duplicate
	dup = false
	for i := 0; i < 5; i++ {
		picked, unpicked, pickedIdx, unpickedIdx :=
			dataset.Columns.RandomPick(ncols, dup, excludeIdx)

		// check if unpicked item exist in picked items.
		for _, un := range unpicked {
			for _, pick := range picked {
				assert.NotEqual(t, un, pick)
			}
		}

		fmt.Println("Random pick without duplicate columns")
		fmt.Println("==> picked columns   :", picked)
		fmt.Println("==> picked idx       :", pickedIdx)
		fmt.Println("==> unpicked columns :", unpicked)
		fmt.Println("==> unpicked idx     :", unpickedIdx)
	}
}
