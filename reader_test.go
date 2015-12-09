// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/shuLhan/dsv"
)

var jsonSample = []string {
	`{}`,
	`{
		"Input"		:"testdata/input.dat"
	}`,
	`{
		"Input"		:"testdata/input.dat"
	}`,
	`{
		"Input"		:"testdata/input.dat"
	,	"InputMetadata"	:
		[{
			"Name"		:"A"
		,	"Separator"	:","
		},{
			"Name"		:"B"
		,	"Separator"	:";"
		}]
	}`,
	`{
		"Input"		:"testdata/input.dat"
	,	"Skip"		:1
	,	"MaxRows"	:1
	,	"InputMetadata"	:
		[{
			"Name"		:"id"
		,	"Separator"	:";"
		,	"Type"		:"integer"
		},{
			"Name"		:"name"
		,	"Separator"	:"-"
		,	"LeftQuote"	:"\""
		,	"RightQuote"	:"\""
		},{
			"Name"		:"value"
		,	"Separator"	:";"
		,	"LeftQuote"	:"[["
		,	"RightQuote"	:"]]"
		},{
			"Name"		:"integer"
		,	"Type"		:"integer"
		,	"Separator"	:";"
		},{
			"Name"		:"real"
		,	"Type"		:"real"
		}]
	}`,
	`{
		"Input"		:"testdata/input.dat"
	,	"Skip"		:1
	,	"MaxRows"	:1
	,	"InputMetadata"	:
		[{
			"Name"		:"id"
		},{
			"Name"		:"editor"
		},{
			"Name"		:"old_rev_id"
		},{
			"Name"		:"new_rev_id"
		},{
			"Name"		:"diff_url"
		},{
			"Name"		:"edit_time"
		},{
			"Name"		:"edit_comment"
		},{
			"Name"		:"article_id"
		},{
			"Name"		:"article_title"
		}]
	}`,
}

var readers = []*dsv.Reader {
	{},
	{
		Input	:"testdata/input.dat",
	},
	{
		Input	:"test-another.dsv",
	},
	{
		Input		:"testdata/input.dat",
		InputMetadata	:[]dsv.Metadata {
			{
				Name		:"A",
				Separator	:",",
			},
			{
				Name		:"B",
				Separator	:";",
			},
		},
	},
}

/*
TestReaderNoInput will print error that the input is not defined.
*/
func TestReaderNoInput (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderNoInput")
	}

	dsvReader := dsv.NewReader ()
	e := dsv.ConfigParse(dsvReader, []byte(jsonSample[0]))

	if nil != e {
		t.Fatal (e)
	}

	e = dsv.InitReader(dsvReader)

	if nil == e {
		t.Fatal ("TestReaderNoInput: failed, should return non nil!")
	}

	if DEBUG {
		log.Print (e)
	}
}

/*
TestConfigParse test parsing metadata.
*/
func TestConfigParse (t *testing.T) {

	if DEBUG {
		log.Println (">>> TestConfigParse")
	}

	cases := []struct {
		in	string
		out	*dsv.Reader
	}{
		{
			jsonSample[1],
			readers[1],
		},
		{
			jsonSample[3],
			readers[3],
		},

	}

	dsvReader := dsv.NewReader ()

	for _, c := range cases {
		e := dsv.ConfigParse(dsvReader, []byte(c.in))

		if e != nil {
			t.Fatal (e)
		}
		if ! dsvReader.IsEqual (c.out) {
			t.Fatal ("Test failed on ", c.in);
		} else {
			if DEBUG {
				log.Print (dsvReader)
			}
		}
	}
}

func TestReaderIsEqual (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderIsEqual")
	}

	cases := []struct {
		in	*dsv.Reader
		out	*dsv.Reader
		result	bool
	}{
		{
			readers[1],
			&dsv.Reader {
				Input :"testdata/input.dat",
			},
			true,
		},
		{
			readers[1],
			readers[2],
			false,
		},

	}

	var r bool

	for _, c := range cases {
		r = c.in.IsEqual (c.out)

		if r != c.result {
			t.Fatal ("Test failed on equality between ", c.in,
				"\n and ", c.out);
		}
	}
}

/*
doRead test reading the DSV data.
*/
func doRead (t *testing.T, dsvReader *dsv.Reader, exp []string) {
	i	:= 0
	n 	:= 0
	e	:= error (nil)

	for {
		n, e = dsv.Read (dsvReader)

		if n > 0 {
			r := fmt.Sprint (dsvReader.Rows)

			if r != exp[i] {
				t.Fatal ("dsv_test: expecting\n",
					exp[i],
					" got\n", r)
			}

			i++
		} else if e == io.EOF {
			// EOF
			break
		}
	}
}

/*
TestReader test reading.
*/
func TestReaderRead (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderRead")
	}

	dsvReader := dsv.NewReader ()

	defer dsvReader.Close ()

	e := dsv.ConfigParse(dsvReader, []byte(jsonSample[4]))

	if nil != e {
		t.Fatal (e)
	}

	e = dsv.InitReader(dsvReader)

	if nil != e {
		t.Fatal (e)
	}

	if DEBUG {
		log.Println (dsvReader)
	}

	doRead (t, dsvReader, expectation)
}

/*
TestReaderOpen real example from the start.
*/
func TestReaderOpen (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderOpen")
	}

	dsvReader := dsv.NewReader ()

	e := dsv.OpenReader(dsvReader, "testdata/config.dsv")

	if nil != e {
		t.Fatal (e)
	}

	defer dsvReader.Close ()

	doRead (t, dsvReader, expectation)
}

func TestOutputMode (t *testing.T) {
	var e error
	var config = []string {`{
		"Input"		:"testdata/input.dat"
	,	"OutputMode"	:"row"
	}`,`{
		"Input"		:"testdata/input.dat"
	,	"OutputMode"	:"rows"
	}`,`{
		"Input"		:"testdata/input.dat"
	,	"OutputMode"	:"columns"
	}`}

	var exps = []struct {
		status bool
		value string
	}{{
		false,
		string (config[0]),
	},{
		true,
		string (config[1]),
	},{
		true,
		string (config[2]),
	}}

	reader := dsv.NewReader ()

	for k,v := range exps {
		e = dsv.ConfigParse(reader, []byte(config[k]))

		if e != nil {
			t.Fatal (e)
		}

		e = dsv.InitReader(reader)

		if e != nil {
			if v.status == true {
				t.Fatal (e)
			}
		}
	}
}

func TestReaderToColumns(t *testing.T) {
	var e error

	reader := dsv.NewReader ()

	e = dsv.ConfigParse(reader, []byte(jsonSample[4]))

	if nil != e {
		t.Fatal (e)
	}

	reader.SetOutputMode(dsv.OutputModeColumns)

	e = dsv.InitReader(reader)

	if nil != e {
		t.Fatal (e)
	}

	var n,i int
	for {
		n, e = dsv.Read (reader)

		if n > 0 {
			reader.TransposeColumnsToRows()

			r := fmt.Sprint(reader.GetData())

			if r != expectation[i] {
				t.Fatal ("dsv_test: expecting\n",
					expectation[i],
					" got\n", r)
			}

			i++
		} else if e == io.EOF {
			// EOF
			break
		}
	}
}

/*
TestReaderSkip will test the 'Skip' option in Metadata.
*/
func TestReaderSkip(t *testing.T) {
	var e error

	fmt.Println("==> TestReaderSkip")

	dsvReader := dsv.NewReader ()

	e = dsv.OpenReader(dsvReader, "testdata/config_skip.dsv")

	if nil != e {
		t.Fatal (e)
	}

	defer dsvReader.Close ()

	doRead (t, dsvReader, exp_skip)
}
