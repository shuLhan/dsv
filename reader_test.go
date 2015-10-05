/*
Copyright 2015 Mhd Sulhan <ms@kilabit.info>
All rights reserved.  Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/
package dsv_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/shuLhan/dsv"
)

var DEBUG = bool (os.Getenv ("DEBUG") != "")

var jsonSample = []string {
	`{}`,
	`{
		"Input"		:"test.dsv"
	}`,
	`{
		"Input"		:"test.dsv"
	}`,
	`{
		"Input"		:"test.dsv"
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
		"Input"		:"test.dsv"
	,	"Skip"		:1
	,	"MaxRecord"	:1
	,	"InputMetadata"	:
		[{
			"Name"		:"id"
		,	"Separator"	:";"
		},{
			"Name"		:"name"
		,	"Separator"	:"-"
		,	"LeftQuote"	:"\""
		,	"RightQuote"	:"\""
		},{
			"Name"		:"value"
		,	"LeftQuote"	:"[["
		,	"RightQuote"	:"]]"
		}]
	}`,
	`{
		"Input"		:"test.dsv"
	,	"Skip"		:1
	,	"MaxRecord"	:1
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
		Input	:"test.dsv",
	},
	{
		Input	:"test-another.dsv",
	},
	{
		Input		:"test.dsv",
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
	e := dsvReader.ParseConfig ([]byte (jsonSample[0]))
	if nil == e {
		t.Error ("TestReaderNoInput: failed, should return non nil!")
	}

	if DEBUG {
		log.Print (e)
	}
}

/*
TestParseConfig test parsing metadata.
*/
func TestParseConfig (t *testing.T) {

	if DEBUG {
		log.Println (">>> TestParseConfig")
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
		e := dsvReader.ParseConfig ([]byte (c.in))

		if e != nil {
			t.Error (e)
		}
		if ! dsvReader.IsEqual (c.out) {
			t.Error ("Test failed on ", c.in);
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
				Input :"test.dsv",
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
			t.Error ("Test failed on equality between ", c.in,
				"\n and ",
					c.out);
		}
	}
}

var expectation = []string {
	"&[\"1\", \"A-B\", \"AB\",]\n",
	"&[\"2\", \"A-B-C\", \"BCD\",]\n",
	"&[\"3\", \"A;B-C,D\", \"A;B C,D\",]\n",
	"&[\"6\", \"\", \"\",]\n",
	"&[\"9\", \"ok\", \"ok\",]\n",
}

/*
doRead test reading the DSV data.
*/
func doRead (dsvReader *dsv.Reader, t *testing.T) {
	exp	:= ""
	i	:= 0
	n 	:= 0
	e	:= error (nil)

	for {
		n, e = dsvReader.Read ()

		if n > 0 {
			r := fmt.Sprint (dsvReader.Records)

			if r != expectation[i] {
				t.Error ("dsv_test: expecting\n", exp,
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
TestReader test reading using 
*/
func TestReaderRead (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderRead")
	}

	dsvReader := dsv.NewReader ()

	dsvReader.ParseConfig ([] byte (jsonSample[4]))

	if DEBUG {
		log.Println (dsvReader)
	}

	defer dsvReader.Close ()

	doRead (dsvReader, t)
}

/*
TestReaderOpen real example from the start.
*/
func TestReaderOpen (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderOpen")
	}

	dsvReader := dsv.NewReader ()

	e := dsvReader.Open ("config.dsv")

	if nil != e {
		t.Error (e)
	}

	defer dsvReader.Close ()

	doRead (dsvReader, t)
}
