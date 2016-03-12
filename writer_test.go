// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"testing"

	"github.com/shuLhan/dsv"
	"github.com/shuLhan/tabula"
	"github.com/shuLhan/tabula/util/assert"
)

/*
TestWriter test reading and writing DSV.
*/
func TestWriter(t *testing.T) {
	rw, e := dsv.New("testdata/config.dsv", nil)
	if e != nil {
		t.Fatal(e)
	}

	doReadWrite(t, &rw.Reader, &rw.Writer, expectation, true)

	e = rw.Close()
	if e != nil {
		t.Fatal(e)
	}

	assert.EqualFileContent(t, rw.GetOutput(), "testdata/expected.dat")
}

/*
TestWriterWithSkip test reading and writing DSV with some column in input being
skipped.
*/
func TestWriterWithSkip(t *testing.T) {
	rw, e := dsv.New("testdata/config_skip.dsv", nil)
	if e != nil {
		t.Fatal(e)
	}

	doReadWrite(t, &rw.Reader, &rw.Writer, expSkip, true)

	e = rw.Close()
	if e != nil {
		t.Fatal(e)
	}

	assert.EqualFileContent(t, rw.GetOutput(), "testdata/expected_skip.dat")
}

/*
TestWriterWithColumns test reading and writing DSV with where each row
is saved in DatasetMode = 'columns'.
*/
func TestWriterWithColumns(t *testing.T) {
	rw, e := dsv.New("testdata/config_skip.dsv", nil)
	if e != nil {
		t.Fatal(e)
	}

	e = rw.SetDatasetMode(dsv.DatasetModeCOLUMNS)
	if e != nil {
		t.Fatal(e)
	}

	doReadWrite(t, &rw.Reader, &rw.Writer, expSkipColumns, true)

	e = rw.Close()
	if e != nil {
		t.Fatal(e)
	}

	assert.EqualFileContent(t, "testdata/expected_skip.dat", rw.GetOutput())
}

func TestWriteRawRows(t *testing.T) {
	dataset := tabula.NewDataset(tabula.DatasetModeRows, nil, nil)

	populateWithRows(t, dataset)

	writer, e := dsv.NewWriter("")
	if e != nil {
		t.Fatal(e)
	}

	outfile := "testdata/writerawrows.out"
	expfile := "testdata/writeraw.exp"

	e = writer.OpenOutput(outfile)

	_, e = writer.WriteRawDataset(dataset, nil)

	if e != nil {
		t.Fatal(e)
	}

	assert.EqualFileContent(t, outfile, expfile)
}

func TestWriteRawColumns(t *testing.T) {
	dataset := tabula.NewDataset(tabula.DatasetModeColumns, nil, nil)

	populateWithColumns(t, dataset)

	writer, e := dsv.NewWriter("")
	if e != nil {
		t.Fatal(e)
	}

	outfile := "testdata/writerawcolumns.out"
	expfile := "testdata/writeraw.exp"

	e = writer.OpenOutput(outfile)

	_, e = writer.WriteRawDataset(dataset, nil)

	if e != nil {
		t.Fatal(e)
	}

	assert.EqualFileContent(t, outfile, expfile)
}
