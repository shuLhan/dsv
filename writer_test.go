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
func doReadWrite (dsvReader *dsv.Reader, dsvWriter *dsv.Writer, t *testing.T) {
	i	:= 0
	n 	:= 0
	e	:= error (nil)

	for {
		n, e = dsv.Read (dsvReader)

		if n > 0 {
			r := fmt.Sprint (dsvReader.Records)

			if r != expectation[i] {
				t.Error ("dsv_test: expecting\n",
					expectation[i],
					" got\n", r)
			}

			dsvWriter.Write (dsvReader)

			i++
		} else if e == io.EOF {
			// EOF
			break
		}
	}
}

/*
TestWriter test reading and writing DSV.
*/
func TestWriter (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestWriter")
	}

	// Initialize dsv reader
	dsvReader := dsv.NewReader ()

	e := dsv.Open (dsvReader, "testdata/config.dsv")

	if nil != e {
		t.Error (e)
	}

	defer dsvReader.Close ()

	// Initialize dsv writer
	dsvWriter := dsv.NewWriter ()

	e = dsvWriter.Open ("testdata/config.dsv", dsvWriter)

	if nil != e {
		t.Error (e)
	}

	if DEBUG {
		log.Print (dsvWriter)
	}

	doReadWrite (dsvReader, dsvWriter, t)
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
