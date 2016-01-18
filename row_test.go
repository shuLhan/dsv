// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"testing"
)

var dataFloat64 = []float64{0.1, 0.2, 0.3, 0.4, 0.5}

func createRow() (row dsv.Row) {
	for _, v := range dataFloat64 {
		rec := &dsv.Record{v}
		row.PushBack(rec)
	}
	return
}

func TestClone(t *testing.T) {
	row := createRow()
	rowClone := row.Clone()
	rowClone2 := row.Clone()

	assert.Equal(t, row, rowClone)

	// changing the clone value should not change the original copy.
	rowClone2[0].SetFloat(0)
	assert.Equal(t, row, rowClone)
	assert.NotEqual(t, row, rowClone2)
}
