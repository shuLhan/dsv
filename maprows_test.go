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

func TestAddRow(t *testing.T) {
	mapRows := dsv.MapRows{}
	rows, e := initRows()

	if e != nil {
		t.Fatal(e)
	}

	for r := range rows {
		key := fmt.Sprint(rows[r][testClassIdx].Value())
		mapRows.AddRow(key, rows[r])
	}

	got := fmt.Sprint(mapRows)

	assert.Equal(t, groupByExpect, got)
}

func TestGetMinority(t *testing.T) {
	mapRows := dsv.MapRows{}
	rows, e := initRows()

	if e != nil {
		t.Fatal(e)
	}

	for r := range rows {
		key := fmt.Sprint(rows[r][testClassIdx].Value())
		mapRows.AddRow(key, rows[r])
	}

	// remove the first row in the first key, so we can make it minority.
	mapRows[0].Value.PopFront()

	_, minRows := mapRows.GetMinority()

	exp := rowsExpect[3]
	got := fmt.Sprint(minRows)

	assert.Equal(t, exp, got)
}
