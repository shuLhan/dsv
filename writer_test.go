// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"testing"

	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
)

/*
TestWriter test reading and writing DSV.
*/
func TestWriter(t *testing.T) {
	rw, e := dsv.New("testdata/config.dsv")
	if e != nil {
		t.Fatal(e)
	}

	doReadWrite(t, &rw.Reader, &rw.Writer, expectation, true)
	rw.Close()

	assert.EqualFileContent(t, rw.GetOutput(), "testdata/expected.dat")
}

/*
TestWriterWithSkip test reading and writing DSV with some column in input being
skipped.
*/
func TestWriterWithSkip(t *testing.T) {
	rw, e := dsv.New("testdata/config_skip.dsv")
	if e != nil {
		t.Fatal(e)
	}

	doReadWrite(t, &rw.Reader, &rw.Writer, exp_skip, true)
	rw.Close()

	assert.EqualFileContent(t, rw.GetOutput(), "testdata/expected_skip.dat")
}

/*
TestWriterWithColumns test reading and writing DSV with where each row
is saved in DatasetMode = 'columns'.
*/
func TestWriterWithColumns(t *testing.T) {
	rw, e := dsv.New("testdata/config_skip.dsv")
	if e != nil {
		t.Fatal(e)
	}

	rw.SetDatasetMode(dsv.DatasetModeCOLUMNS)

	doReadWrite(t, &rw.Reader, &rw.Writer, exp_skip_columns, true)
	rw.Close()

	assert.EqualFileContent(t, "testdata/expected_skip.dat", rw.GetOutput())
}
