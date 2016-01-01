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

/*
TestRecord simply check how the stringer work.
*/
func TestRecord(t *testing.T) {
	expec := []string{"test", "1", "2"}
	expec_type := []int{dsv.TString, dsv.TInteger, dsv.TInteger}

	row := make(dsv.Row, 0)

	for i := range expec {
		r, e := dsv.NewRecord([]byte(expec[i]), expec_type[i])
		if nil != e {
			t.Error(e)
		}

		row = append(row, r)
	}

	exp := fmt.Sprint(expec)
	got := fmt.Sprint(row)
	assert.Equal(t, exp, got)
}
