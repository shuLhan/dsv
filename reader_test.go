/*
Copyright 2015 Mhd Sulhan <ms@kilabit.info>
All rights reserved.  Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/
package dsv_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/shuLhan/dsv"
)

var DEBUG = bool (os.Getenv ("DEBUG") != "")

var json_sample = []string {
	`{}`,
	`{
		"Input"		:"test.dsv"
	}`,
	`{
		"Input"		:"test.dsv"
	}`,
	`{
		"Input"		:"test.dsv"
	,	"FieldMetadata"	:
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
	,	"FieldMetadata"	:
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
		,	"RightQuote"	:"-"
		}]
	}`,
	`{
		"Input"		:"test.dsv"
	,	"Skip"		:1
	,	"MaxRecord"	:1
	,	"FieldMetadata"	:
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
		FieldMetadata	:[]dsv.Metadata {
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

	dsv := dsv.NewReader ()
	e := dsv.ParseFieldMetadata (json_sample[0])
	if nil == e {
		t.Error ("TestReaderNoInput: failed, should return non nil!")
	}

	if DEBUG {
		log.Print (e)
	}
}

/*
TestParseFieldMetadata test parsing metadata.
*/
func TestParseFieldMetadata (t *testing.T) {

	if DEBUG {
		log.Println (">>> TestParseFieldMetadata")
	}

	cases := []struct {
		in	string
		out	*dsv.Reader
	}{
		{
			json_sample[1],
			readers[1],
		},
		{
			json_sample[3],
			readers[3],
		},

	}

	dsv := dsv.NewReader ()

	for _, c := range cases {
		e := dsv.ParseFieldMetadata (c.in)

		if e != nil {
			t.Error (e)
		}
		if ! dsv.IsEqual (c.out) {
			t.Error ("Test failed on ", c.in);
		} else {
			if DEBUG {
				log.Print (dsv)
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

/*
TestReader test reading using 
*/
func TestReader (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReaderSkip")
	}

	dsv := dsv.NewReader ()

	dsv.ParseFieldMetadata (json_sample[4])

	if DEBUG {
		log.Println (dsv)
	}

	expectation := []string {
		"&[\"1\", \"A-B\", \"AB\",]\n",
		"&[\"2\", \"A-B-C\", \"BCD\",]\n",
		"&[\"3\", \"A;B-C,D\", \"A;B C,D\",]\n",
	}

	exp	:= ""
	n	:= 0
	e	:= error(nil)

	for i := range expectation {
		n = 0
		e = nil

		for n <= 0 || e != nil {
			n, e = dsv.Read ()
		}

		r := fmt.Sprint (dsv.Records)
		exp += expectation[i]

		if r != exp {
			t.Error ("dsv_test: expecting\n", exp, " got\n", r)
		}
	}
}
