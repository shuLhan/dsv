// Copyright 2015-2018, Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"testing"
)

/*
doInit create read-write object.
*/
func doInit(t *testing.T, fcfg string) (rw *dsv.ReadWriter, e error) {
	// Initialize dsv
	rw, e = dsv.New(fcfg, nil)

	if nil != e {
		t.Fatal(e)
	}

	return
}

/*
TestReadWriter test reading and writing DSV.
*/
func TestReadWriter(t *testing.T) {
	rw, _ := doInit(t, "testdata/config.dsv")

	doReadWrite(t, &rw.Reader, &rw.Writer, expectation, true)

	e := rw.Close()
	if e != nil {
		t.Fatal(e)
	}

	assertFile(t, rw.GetOutput(), "testdata/expected.dat", true)
}

/*
TestReadWriter test reading and writing DSV.
*/
func TestReadWriterAll(t *testing.T) {
	rw, _ := doInit(t, "testdata/config.dsv")

	rw.SetMaxRows(-1)

	doReadWrite(t, &rw.Reader, &rw.Writer, expectation, false)

	e := rw.Close()
	if e != nil {
		t.Fatal(e)
	}

	assertFile(t, rw.GetOutput(), "testdata/expected.dat", true)
}

func TestSimpleReadWrite(t *testing.T) {
	fcfg := "testdata/config_simpleread.dsv"

	reader, e := dsv.SimpleRead(fcfg, nil)
	if e != nil {
		t.Fatal(e)
	}

	fout := "testdata/output.dat"
	fexp := "testdata/expected.dat"

	_, e = dsv.SimpleWrite(reader, fcfg)
	if e != nil {
		t.Fatal(e)
	}

	assertFile(t, fexp, fout, true)
}

func TestSimpleMerge(t *testing.T) {
	fcfg1 := "testdata/config_simpleread.dsv"
	fcfg2 := "testdata/config_simpleread.dsv"

	reader, e := dsv.SimpleMerge(fcfg1, fcfg2, nil, nil)
	if e != nil {
		t.Fatal(e)
	}

	_, e = dsv.SimpleWrite(reader, fcfg1)
	if e != nil {
		t.Fatal(e)
	}

	fexp := "testdata/expected_simplemerge.dat"
	fout := "testdata/output.dat"

	assertFile(t, fexp, fout, true)
}
