// Copyright 2015-2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/tabula"
	"testing"
)

func TestReaderWithClaset(t *testing.T) {
	fcfg := "testdata/claset.dsv"

	claset := tabula.Claset{}

	_, e := dsv.NewReader(fcfg, &claset)
	if e != nil {
		t.Fatal(e)
	}

	assert(t, 3, claset.GetClassIndex(), true)

	claset.SetMajorityClass("regular")
	claset.SetMinorityClass("vandalism")

	clone := claset.Clone().(tabula.ClasetInterface)

	assert(t, 3, clone.GetClassIndex(), true)
	assert(t, "regular", clone.MajorityClass(), true)
	assert(t, "vandalism", clone.MinorityClass(), true)
}
