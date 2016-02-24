// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/tabula/util/assert"
	"testing"
)

/*
doInit create read-write object.
*/
func doInit(t *testing.T, fcfg string) (rw *dsv.ReadWriter, e error) {
	// Initialize dsv
	rw, e = dsv.New(fcfg)

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

	assert.EqualFileContent(t, rw.GetOutput(), "testdata/expected.dat")
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

	assert.EqualFileContent(t, rw.GetOutput(), "testdata/expected.dat")
}
