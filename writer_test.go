// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"testing"

	"github.com/shuLhan/dsv"
)

/*
doReadWrite test reading and writing the DSV data.
*/
func doReadWrite (t *testing.T, dsvReader *dsv.Reader, dsvWriter *dsv.Writer,
						expectation []string) {
	i	:= 0
	e	:= error (nil)

	for {
		_, e = dsv.Read(dsvReader)

		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal(e)
			break
		}

		r := fmt.Sprint(dsvReader.GetOutput())

		if r != expectation[i] {
			t.Error("dsv_test: expecting\n", expectation[i],
				" got\n", r)
		}

		_, e = dsvWriter.Write(dsvReader)

		if e != nil {
			t.Fatal(e)
		}

		i++
	}
}

/*
TestWriter test reading and writing DSV.
*/
func TestWriter (t *testing.T) {
	// Initialize dsv reader
	dsvReader := dsv.NewReader ()

	e := dsv.Open (dsvReader, "testdata/config.dsv")

	if nil != e {
		t.Error (e)
	}

	defer dsvReader.Close ()

	// Initialize dsv writer
	dsvWriter := dsv.NewWriter ()

	e = dsv.Open (dsvWriter, "testdata/config.dsv")

	if nil != e {
		t.Error (e)
	}

	if DEBUG {
		log.Print (dsvWriter)
	}

	doReadWrite (t, dsvReader, dsvWriter, expectation)
	dsvWriter.Close ()

	// Compare the ouput from Writer
	out, e := ioutil.ReadFile (dsvWriter.Output)

	if nil != e {
		t.Error (e)
	}

	exp, e := ioutil.ReadFile ("testdata/expected.dat")

	if nil != e {
		t.Error (e)
	}

	r := bytes.Compare (out, exp)

	if 0 != r {
		t.Error ("Output different from expected (", r ,")")
	}
}

/*
TestWriterWithSkip test reading and writing DSV with some field in input being
skipped.
*/
func TestWriterWithSkip (t *testing.T) {
	// Initialize dsv reader
	dsvReader := dsv.NewReader ()

	e := dsv.Open (dsvReader, "testdata/config_skip.dsv")

	if nil != e {
		t.Error (e)
	}

	defer dsvReader.Close ()

	// Initialize dsv writer
	dsvWriter := dsv.NewWriter ()

	e = dsv.Open (dsvWriter, "testdata/config_skip.dsv")

	if nil != e {
		t.Error (e)
	}

	if DEBUG {
		log.Print (dsvWriter)
	}

	doReadWrite (t, dsvReader, dsvWriter, exp_skip)
	dsvWriter.Close ()

	// Compare the Writer output
	out, e := ioutil.ReadFile (dsvWriter.Output)

	if nil != e {
		t.Error (e)
	}

	exp, e := ioutil.ReadFile ("testdata/expected_skip.dat")

	if nil != e {
		t.Error (e)
	}

	r := bytes.Compare (out, exp)

	if 0 != r {
		t.Error ("Output different from expected (", r ,")")
	}
}

/*
TestWriterWithFields test reading and writing DSV with where each row
is saved in OutputMode = 'fields'.
*/
func TestWriterWithFields(t *testing.T) {
	// Initialize dsv reader
	dsvReader := dsv.NewReader()

	e := dsv.Open(dsvReader, "testdata/config_skip.dsv")
	if nil != e {
		t.Error(e)
	}
	dsvReader.InitOutputMode("fields")
	defer dsvReader.Close()

	// Initialize dsv writer
	dsvWriter := dsv.NewWriter()
	e = dsv.Open(dsvWriter, "testdata/config_skip.dsv")
	if nil != e {
		t.Error(e)
	}

	doReadWrite(t, dsvReader, dsvWriter, exp_skip_fields)
	dsvWriter.Close()

	// Compare the Writer output
	out, e := ioutil.ReadFile(dsvWriter.Output)

	if nil != e {
		t.Error(e)
	}

	exp, e := ioutil.ReadFile("testdata/expected_skip.dat")

	if nil != e {
		t.Error(e)
	}

	r := bytes.Compare(out, exp)

	if 0 != r {
		t.Error("Output different from expected (", r ,")")
	}
}
