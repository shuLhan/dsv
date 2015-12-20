// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"github.com/shuLhan/dsv"
	"os"
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
	"[{0 0 [A-B]} {0 0 [AB]} {1 0 [1]} {2 0 [0.1]}]",
	"[{0 0 [A-B-C]} {0 0 [BCD]} {1 0 [2]} {2 0 [0.02]}]",
	"[{0 0 [A;B-C,D]} {0 0 [A;B C,D]} {1 0 [3]} {2 0 [0.003]}]",
	"[{0 0 []} {0 0 []} {1 0 [6]} {2 0 [0.000006]}]",
	"[{0 0 [ok]} {0 0 [ok]} {1 0 [9]} {2 0 [0.000000009]}]",
	"[{0 0 [test]} {0 0 [integer]} {1 0 [10]} {2 0 [0.101]}]",
	"[{0 0 [test]} {0 0 [real]} {1 0 [123456789]} {2 0 [0.123456789]}]",
}

var exp_skip_columns_all = []string{
	"{0 0 [A-B A-B-C A;B-C,D  ok test test]}",
	"{0 0 [AB BCD A;B C,D  ok integer real]}",
	"{1 0 [1 2 3 6 9 10 123456789]}",
	"{2 0 [0.1 0.02 0.003 0.000006 0.000000009 0.101 0.123456789]}",
}

var exp_skip_columns_all_rev = []string{
	"{0 0 [test test ok  A;B-C,D A-B-C A-B]}",
	"{0 0 [real integer ok  A;B C,D BCD AB]}",
	"{1 0 [123456789 10 9 6 3 2 1]}",
	"{2 0 [0.123456789 0.101 0.000000009 0.000006 0.003 0.02 0.1]}",
}

// Testing data and function for Rows and MapRows
var rowsData = [][]byte{
	{'1', '5', '9', '+'},
	{'2', '6', '0', '-'},
	{'3', '7', '1', '-'},
	{'4', '8', '2', '+'},
}

var testColTypes = []int{dsv.TInteger, dsv.TInteger, dsv.TInteger, dsv.TString}

var rowsExpect = []string{
	"[1 5 9 +]",
	"[2 6 0 -]",
	"[3 7 1 -]",
	"[4 8 2 +]",
}

var testClassIdx = 3

var groupByExpect = "[{+ [1 5 9 +][4 8 2 +]} {- [2 6 0 -][3 7 1 -]}]"

func initRows() (rows dsv.Rows, e error) {
	for i := range rowsData {
		l := len(rowsData[i])
		row := make(dsv.Row, 0)

		for j := 0; j < l; j++ {
			rec, e := dsv.NewRecord([]byte{rowsData[i][j]},
				testColTypes[j])

			if nil != e {
				return nil, e
			}

			row = append(row, rec)
		}

		rows.PushBack(row)
	}
	return rows, nil
}
