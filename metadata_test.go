// Copyright 2015-2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"testing"
)

func TestMetadataIsEqual(t *testing.T) {
	cases := []struct {
		in     dsv.Metadata
		out    dsv.Metadata
		result bool
	}{
		{
			dsv.Metadata{
				Name:      "A",
				Separator: ",",
			},
			dsv.Metadata{
				Name:      "A",
				Separator: ",",
			},
			true,
		},
		{
			dsv.Metadata{
				Name:      "A",
				Separator: ",",
			},
			dsv.Metadata{
				Name:      "A",
				Separator: ";",
			},
			false,
		},
	}

	for _, c := range cases {
		r := c.in.IsEqual(&c.out)

		if r != c.result {
			t.Error("Test failed on ", c.in, c.out)
		}
	}
}
