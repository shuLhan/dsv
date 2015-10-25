// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shuLhan/dsv"
)

var exp = []string {
	"0\n",
	"1\n",
	"2\n",
	"3\n",
	"4\n",
}

var records []dsv.RecordSlice;

func initRecords () (e error) {
	var i = 0
	var j = 0
	var z = 0
	var l = 0
	var r dsv.Record
	var data = [][]byte{
		{ '1', dsv.TInteger, '+', dsv.TString },
		{ '2', dsv.TInteger, '-', dsv.TString },
		{ '3', dsv.TInteger, '-', dsv.TString },
		{ '4', dsv.TInteger, '+', dsv.TString },
	}

	records = make([]dsv.RecordSlice, len (data))

	for i = range data {
		l = len(data[i])
		records[i] = make(dsv.RecordSlice, l / 2)

		z = 0
		for j = 0; j < l; j += 2 {
			r, e = dsv.RecordNew ([]byte{data[i][j]}, int(data[i][j+1]))

			if nil != e {
				return e
			}

			records[i][z] = r

			z++
		}
	}

	return e
}

func TestRowPopFrontRow (t *testing.T) {
	if DEBUG {
		fmt.Println (">>> TestRowPopFrontRow")
	}

	rows := &dsv.Row {}

	for i := 0; i < 5; i++ {
		rows.PushBack (i)
	}

	exps := strings.Join (exp, "")
	got := fmt.Sprint (rows)

	if got != exps {
		t.Fatal ("Expecting:\n", exps, "\n Got:\n", got)
	}

	row := rows.PopFrontRow ()

	exps = strings.Join (exp[:1], "")
	got = fmt.Sprint (row)

	if got != exps {
		t.Fatal ("Expecting:\n", exps, "\n Got:\n", got)
	}

	exps = strings.Join (exp[1:], "")
	got = fmt.Sprint (rows)

	if got != exps {
		t.Fatal ("Expecting:\n", exps, "\n Got:\n", got)
	}

	rows.PopFrontRow ()
	rows.PopFrontRow ()
	rows.PopFrontRow ()
	rows.PopFrontRow ()

	row = rows.PopFrontRow ()

	if nil != row {
		t.Fatal ("Expecting:\n", nil, "\n Got:\n", row)
	}
}

func generateRow (r *[]dsv.RecordSlice) (rows *dsv.Row) {
	rows = &dsv.Row{}

	for i := range records {
		rows.PushBack (&records[i])
	}
	return
}

/*
TestRowGroupByValue test grouping in rows.
*/
func TestRowGroupByValue (t *testing.T) {
	var exps = []string {
		`[[1 +] [2 -] [3 -] [4 +]]`,
		`&[1 +]
&[2 -]
&[3 -]
&[4 +]
`,
		`map[+:&[1 +]
&[4 +]
 -:&[2 -]
&[3 -]
]`,
	}

	initRecords ()

	// test records
	got := fmt.Sprint (records)

	i := 0
	if got != exps[i] {
		t.Fatal ("Expecting:\n", exps[i], "\n Got:\n", got)
	}

	// test rows
	rows := generateRow (&records)

	got = fmt.Sprint (rows)

	i = 1
	if got != exps[i] {
		t.Fatal ("Expecting:\n", exps[i], "\n Got:\n", got)
	}

	// test rows grouping
	classes := rows.GroupByValue (1)

	got = fmt.Sprint (classes)

	i = 2
	if got != exps[i] {
		t.Fatal ("Expecting:\n", exps[i], "\n Got:\n", got)
	}

	if DEBUG {
		for k,v := range classes {
			fmt.Println (k, v.Len (), v)
		}
	}
}

func TestRowRandomPick (t *testing.T) {
	var shuffled *dsv.Row

	initRecords ()

	for i := 0; i < 5; i++ {
		rows := generateRow (&records)

		shuffled = rows.RandomPick (2)
		fmt.Println ("Shuffled rows:\n", shuffled)
	}
}
