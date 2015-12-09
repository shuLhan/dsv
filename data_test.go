// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"os"

	"github.com/shuLhan/dsv"
)

var DEBUG = bool(os.Getenv("DEBUG") != "")

var expectation = []string{
	"[1 A-B AB 1 0.1]",
	"[2 A-B-C BCD 2 0.02]",
	"[3 A;B-C,D A;B C,D 3 0.003]",
	"[6   6 0.000006]",
	"[9 ok ok 9 0.000000009]",
	"[10 test integer 10 0.101]",
	"[12 test real 123456789 0.123456789]",
}

var exp_skip = []string{
	"[A-B AB 1 0.1]",
	"[A-B-C BCD 2 0.02]",
	"[A;B-C,D A;B C,D 3 0.003]",
	"[  6 0.000006]",
	"[ok ok 9 0.000000009]",
	"[test integer 10 0.101]",
	"[test real 123456789 0.123456789]",
}

var exp_skip_columns = []string{
	"[[A-B] [AB] [1] [0.1]]",
	"[[A-B-C] [BCD] [2] [0.02]]",
	"[[A;B-C,D] [A;B C,D] [3] [0.003]]",
	"[[] [] [6] [0.000006]]",
	"[[ok] [ok] [9] [0.000000009]]",
	"[[test] [integer] [10] [0.101]]",
	"[[test] [real] [123456789] [0.123456789]]",
}

var exp_skip_columns_all = []string{
	"[A-B A-B-C A;B-C,D  ok test test]",
	"[AB BCD A;B C,D  ok integer real]",
	"[1 2 3 6 9 10 123456789]",
	"[0.1 0.02 0.003 0.000006 0.000000009 0.101 0.123456789]",
}

// Testing data and function for Rows and MapRows
var rowsData = [][]byte{
	{'1', dsv.TInteger, '+', dsv.TString},
	{'2', dsv.TInteger, '-', dsv.TString},
	{'3', dsv.TInteger, '-', dsv.TString},
	{'4', dsv.TInteger, '+', dsv.TString},
}

var rowsExpect = []string{
	"[1 +]",
	"[2 -]",
	"[3 -]",
	"[4 +]",
}

var groupByExpect = "[{+ [1 +][4 +]} {- [2 -][3 -]}]"

func initRows() (rows dsv.Rows, e error) {
	for i := range rowsData {
		l := len(rowsData[i])
		row := make(dsv.Row, 0)

		z := 0
		for j := 0; j < l; j += 2 {
			rec, e := dsv.NewRecord([]byte{rowsData[i][j]}, int(rowsData[i][j+1]))

			if nil != e {
				return nil, e
			}

			row = append(row, *rec)

			z++
		}

		rows.PushBack(row)
	}
	return rows, nil
}
