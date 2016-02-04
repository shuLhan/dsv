// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"github.com/shuLhan/dsv"
	"github.com/shuLhan/dsv/util/assert"
	"io"
	"strings"
	"testing"
)

var jsonSample = []string{
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

var readers = []*dsv.Reader{
	{},
	{
		Input: "testdata/input.dat",
	},
	{
		Input: "test-another.dsv",
	},
	{
		Input: "testdata/input.dat",
		InputMetadata: []dsv.Metadata{
			{
				Name:      "A",
				Separator: ",",
			},
			{
				Name:      "B",
				Separator: ";",
			},
		},
	},
}

/*
TestReaderNoInput will print error that the input is not defined.
*/
func TestReaderNoInput(t *testing.T) {
	dsvReader, e := dsv.NewReader("")
	if nil != e {
		t.Fatal(e)
	}

	e = dsv.ConfigParse(dsvReader, []byte(jsonSample[0]))

	if nil != e {
		t.Fatal(e)
	}

	e = dsv.InitReader(dsvReader)

	if nil == e {
		t.Fatal("TestReaderNoInput: failed, should return non nil!")
	}
}

/*
TestConfigParse test parsing metadata.
*/
func TestConfigParse(t *testing.T) {
	cases := []struct {
		in  string
		out *dsv.Reader
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

	dsvReader, e := dsv.NewReader("")
	if nil != e {
		t.Fatal(e)
	}

	for _, c := range cases {
		e := dsv.ConfigParse(dsvReader, []byte(c.in))

		if e != nil {
			t.Fatal(e)
		}
		if !dsvReader.IsEqual(c.out) {
			t.Fatal("Test failed on ", c.in)
		}
	}
}

func TestReaderIsEqual(t *testing.T) {
	cases := []struct {
		in     *dsv.Reader
		out    *dsv.Reader
		result bool
	}{
		{
			readers[1],
			&dsv.Reader{
				Input: "testdata/input.dat",
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
		r = c.in.IsEqual(c.out)

		if r != c.result {
			t.Fatal("Test failed on equality between ", c.in,
				"\n and ", c.out)
		}
	}
}

/*
doRead test reading the DSV data.
*/
func doRead(t *testing.T, dsvReader *dsv.Reader, exp []string) {
	i := 0
	n := 0
	e := error(nil)

	for {
		n, e = dsv.Read(dsvReader)

		if n > 0 {
			r := fmt.Sprint(dsvReader.Rows)

			assert.Equal(t, exp[i], r)

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
func TestReaderRead(t *testing.T) {
	dsvReader, e := dsv.NewReader("")
	if nil != e {
		t.Fatal(e)
	}

	e = dsv.ConfigParse(dsvReader, []byte(jsonSample[4]))

	if nil != e {
		t.Fatal(e)
	}

	e = dsv.InitReader(dsvReader)

	if nil != e {
		t.Fatal(e)
	}

	doRead(t, dsvReader, expectation)

	e = dsvReader.Close()
	if e != nil {
		t.Fatal(e)
	}
}

/*
TestReaderOpen real example from the start.
*/
func TestReaderOpen(t *testing.T) {
	dsvReader, e := dsv.NewReader("testdata/config.dsv")
	if nil != e {
		t.Fatal(e)
	}

	doRead(t, dsvReader, expectation)

	e = dsvReader.Close()
	if e != nil {
		t.Fatal(e)
	}
}

func TestDatasetMode(t *testing.T) {
	var e error
	var config = []string{`{
		"Input"		:"testdata/input.dat"
	,	"DatasetMode"	:"row"
	}`, `{
		"Input"		:"testdata/input.dat"
	,	"DatasetMode"	:"rows"
	}`, `{
		"Input"		:"testdata/input.dat"
	,	"DatasetMode"	:"columns"
	}`}

	var exps = []struct {
		status bool
		value  string
	}{{
		false,
		string(config[0]),
	}, {
		true,
		string(config[1]),
	}, {
		true,
		string(config[2]),
	}}

	reader, e := dsv.NewReader("")
	if nil != e {
		t.Fatal(e)
	}

	for k, v := range exps {
		e = dsv.ConfigParse(reader, []byte(config[k]))

		if e != nil {
			t.Fatal(e)
		}

		e = dsv.InitReader(reader)

		if e != nil {
			if v.status == true {
				t.Fatal(e)
			}
		}
	}
}

func TestReaderToColumns(t *testing.T) {
	reader, e := dsv.NewReader("")

	e = dsv.ConfigParse(reader, []byte(jsonSample[4]))

	if nil != e {
		t.Fatal(e)
	}

	e = reader.SetDatasetMode(dsv.DatasetModeCOLUMNS)
	if e != nil {
		t.Fatal(e)
	}

	e = dsv.InitReader(reader)

	if nil != e {
		t.Fatal(e)
	}

	var n, i int
	for {
		n, e = dsv.Read(reader)

		if n > 0 {
			reader.TransposeToRows()

			r := fmt.Sprint(reader.GetData())

			assert.Equal(t, expectation[i], r)

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
	dsvReader, e := dsv.NewReader("testdata/config_skip.dsv")
	if nil != e {
		t.Fatal(e)
	}

	doRead(t, dsvReader, expSkip)

	e = dsvReader.Close()
	if e != nil {
		t.Fatal(e)
	}
}

func TestTransposeToColumns(t *testing.T) {
	reader, e := dsv.NewReader("testdata/config_skip.dsv")
	if nil != e {
		t.Fatal(e)
	}

	reader.SetMaxRows(-1)

	_, e = dsv.Read(reader)
	if e != io.EOF {
		t.Fatal(e)
	}

	reader.TransposeToColumns()

	exp := fmt.Sprint(expSkipColumnsAll)

	columns := reader.GetDataAsColumns()

	got := fmt.Sprint(columns)

	assert.Equal(t, exp, got)

	e = reader.Close()
	if e != nil {
		t.Fatal(e)
	}
}

func TestSortColumnsByIndex(t *testing.T) {
	reader, e := dsv.NewReader("testdata/config_skip.dsv")
	if nil != e {
		t.Fatal(e)
	}

	reader.SetMaxRows(-1)

	_, e = dsv.Read(reader)
	if e != io.EOF {
		t.Fatal(e)
	}

	// reverse the data
	var idxReverse []int
	var expReverse []string

	for x := len(expSkip) - 1; x >= 0; x-- {
		idxReverse = append(idxReverse, x)
		expReverse = append(expReverse, expSkip[x])
	}

	reader.SortColumnsByIndex(idxReverse)

	exp := strings.Join(expReverse, "")
	got := fmt.Sprint(reader.GetDataAsRows())

	assert.Equal(t, exp, got)

	exp = "[" + strings.Join(expSkipColumnsAllRev, " ") + "]"

	columns := reader.GetDataAsColumns()

	got = fmt.Sprint(columns)

	assert.Equal(t, exp, got)

	e = reader.Close()
	if e != nil {
		t.Fatal(e)
	}
}

func TestSplitRowsByValue(t *testing.T) {
	reader, e := dsv.NewReader("testdata/config.dsv")
	if nil != e {
		t.Fatal(e)
	}

	reader.SetMaxRows(256)

	_, e = dsv.Read(reader)

	if e != nil && e != io.EOF {
		t.Fatal(e)
	}

	splitL, splitR, e := reader.SplitRowsByValue(0, 6)

	if e != nil {
		t.Fatal(e)
	}

	// test left split
	exp := ""
	for x := 0; x < 3; x++ {
		exp += expectation[x]
	}

	got := fmt.Sprint(splitL.GetDataAsRows())

	assert.Equal(t, exp, got)

	// test right split
	exp = ""
	for x := 3; x < len(expectation); x++ {
		exp += expectation[x]
	}

	got = fmt.Sprint(splitR.GetDataAsRows())

	assert.Equal(t, exp, got)

	e = reader.Close()
	if e != nil {
		t.Fatal(e)
	}
}

func TestMergeColumns(t *testing.T) {
	reader1, e := dsv.NewReader("testdata/config.dsv")
	if nil != e {
		t.Fatal(e)
	}

	reader2, e := dsv.NewReader("testdata/config_skip.dsv")
	if nil != e {
		t.Fatal(e)
	}

	reader1.SetMaxRows(-1)
	reader2.SetMaxRows(-1)

	_, e = dsv.Read(reader1)
	if e != io.EOF {
		t.Fatal(e)
	}

	_, e = dsv.Read(reader2)
	if e != io.EOF {
		t.Fatal(e)
	}

	reader1.Close()
	reader2.Close()

	reader1.InputMetadata[len(reader1.InputMetadata)-1].Separator = ";"

	reader1.MergeColumns(reader2)

	// write merged reader
	writer, e := dsv.NewWriter("")
	if e != nil {
		t.Fatal(e)
	}

	outfile := "testdata/output_merge.dat"
	expfile := "testdata/expected_merge.dat"

	e = writer.OpenOutput(outfile)

	if e != nil {
		t.Fatal(e)
	}

	sep := "\t"
	writer.WriteDataset(&reader1.Dataset, &sep)

	writer.Close()

	assert.EqualFileContent(t, outfile, expfile)
}
