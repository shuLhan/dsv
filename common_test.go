// Copyright 2015-2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/tabula"
	"github.com/shuLhan/tabula/util/assert"
	"io"
	"testing"
)

/*
doReadWrite test reading and writing the DSV data.
*/
func doReadWrite(t *testing.T, dsvReader *dsv.Reader, dsvWriter *dsv.Writer,
	expectation []string, check bool) {
	i := 0

	for {
		n, e := dsv.Read(dsvReader)

		if e == io.EOF {
			_, e = dsvWriter.Write(dsvReader)
			if e != nil {
				t.Fatal(e)
			}

			break
		}

		if e != nil {
			continue
		}

		if n > 0 {
			if check {
				r := fmt.Sprint(dsvReader.GetData())
				assert.Equal(t, expectation[i], r)
				i++
			}

			_, e = dsvWriter.Write(dsvReader)
			if e != nil {
				t.Fatal(e)
			}
		}
	}

	e := dsvWriter.Flush()
	if e != nil {
		t.Fatal(e)
	}
}

var datasetRows = [][]string{
	{"0", "1", "A"},
	{"1", "1.1", "B"},
	{"2", "1.2", "A"},
	{"3", "1.3", "B"},
	{"4", "1.4", "C"},
	{"5", "1.5", "D"},
	{"6", "1.6", "C"},
	{"7", "1.7", "D"},
	{"8", "1.8", "E"},
	{"9", "1.9", "F"},
}

var datasetCols = [][]string{
	{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
	{"1", "1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9"},
	{"A", "B", "A", "B", "C", "D", "C", "D", "E", "F"},
}

var datasetTypes = []int{
	tabula.TInteger,
	tabula.TReal,
	tabula.TString,
}

var datasetNames = []string{"int", "real", "string"}

func populateWithRows(t *testing.T, dataset *tabula.Dataset) {
	for _, rowin := range datasetRows {
		row := make(tabula.Row, len(rowin))

		for x, recin := range rowin {
			rec, e := tabula.NewRecord(recin, datasetTypes[x])
			if e != nil {
				t.Fatal(e)
			}

			row[x] = rec
		}

		dataset.PushRow(row)
	}
}

func populateWithColumns(t *testing.T, dataset *tabula.Dataset) {
	for x := range datasetCols {
		col, e := tabula.NewColumnString(datasetCols[x], datasetTypes[x],
			datasetNames[x])
		if e != nil {
			t.Fatal(e)
		}

		dataset.PushColumn(*col)
	}
}
