// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"testing"
)

var data = []string{"9.987654321", "8.8", "7.7", "6.6", "5.5", "4.4", "3.3"}
var expFloat = []float64{9.987654321, 8.8, 7.7, 6.6, 5.5, 4.4, 3.3}

func TestToFloatSlice(t *testing.T) {
	var col dsv.Column

	for x := range data {
		rec, e := dsv.NewRecord([]byte(data[x]), dsv.TReal)
		if e != nil {
			t.Fatal(e)
		}

		col.PushBack(rec)
	}

	got := col.ToFloatSlice()

	assert.Equal(t, expFloat, got)
}

func TestToStringSlice(t *testing.T) {
	var col dsv.Column

	for x := range data {
		rec, e := dsv.NewRecord([]byte(data[x]), dsv.TString)
		if e != nil {
			t.Fatal(e)
		}

		col.PushBack(rec)
	}

	got := col.ToStringSlice()

	assert.Equal(t, data, got)
}
