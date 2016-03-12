// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/tabula"
	"github.com/shuLhan/tabula/util/assert"
	"testing"
)

func TestReaderWithClaset(t *testing.T) {
	fcfg := "testdata/claset.dsv"

	claset := tabula.Claset{}

	_, e := dsv.NewReader(fcfg, &claset)
	if e != nil {
		t.Fatal(e)
	}

	assert.Equal(t, 3, claset.GetClassIndex())

	claset.SetMajorityClass("regular")
	claset.SetMinorityClass("vandalism")

	clone := claset.Clone().(tabula.ClasetInterface)

	assert.Equal(t, 3, clone.GetClassIndex())
	assert.Equal(t, "regular", clone.MajorityClass())
	assert.Equal(t, "vandalism", clone.MinorityClass())
}
