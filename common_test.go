// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"io"
	"testing"
)

/*
doReadWrite test reading and writing the DSV data.
*/
func doReadWrite(t *testing.T, dsvReader *dsv.Reader, dsvWriter *dsv.Writer,
	expectation []string, check bool) {
	i := 0

	for {
		n, e := dsv.Read(dsvReader)

		if e == io.EOF {
			_, e = dsvWriter.Write(dsvReader)
			if e != nil {
				t.Fatal(e)
			}

			break
		}

		if e != nil {
			continue
		}

		if n > 0 {
			if check {
				r := fmt.Sprint(dsvReader.GetData())
				assert.Equal(t, expectation[i], r)
				i++
			}

			_, e = dsvWriter.Write(dsvReader)
			if e != nil {
				t.Fatal(e)
			}
		}
	}

	e := dsvWriter.Flush()
	if e != nil {
		t.Fatal(e)
	}
}
