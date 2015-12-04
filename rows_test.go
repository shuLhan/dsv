// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/shuLhan/dsv"
)

var exp = []string{
	"0\n",
	"1\n",
	"2\n",
	"3\n",
	"4\n",
}

func TestPushBack(t *testing.T) {
	rows, e := initRecords()
	if e != nil {
		t.Fatal(e)
	}

	exp := strings.Join(rowsExpect, "")
	got := fmt.Sprint(rows)

	assert(t, exp, got)
}

func TestPopFront(t *testing.T) {
	rows, e := initRecords()
	if e != nil {
		t.Fatal(e)
	}

	l := len(rows) - 1
	for i := range rows {
		row := rows.PopFront()

		exp := rowsExpect[i]
		got := fmt.Sprint(row)

		assert(t, exp, got)

		if i < l {
			exp = strings.Join(rowsExpect[i+1:], "")
		} else {
			exp = ""
		}
		got = fmt.Sprint(rows)

		assert(t, exp, got)
	}

	// empty rows
	row := rows.PopFront()

	exp := "[]"
	got := fmt.Sprint(row)

	assert(t, exp, got)
}

func TestPopFrontRow(t *testing.T) {
	rows, e := initRecords()
	if e != nil {
		t.Fatal(e)
	}

	l := len(rows) - 1
	for i := range rows {
		newRows := rows.PopFrontAsRows()

		exp := rowsExpect[i]
		got := fmt.Sprint(newRows)

		assert(t, exp, got)

		if i < l {
			exp = strings.Join(rowsExpect[i+1:], "")
		} else {
			exp = ""
		}
		got = fmt.Sprint(rows)

		assert(t, exp, got)
	}

	// empty rows
	row := rows.PopFrontAsRows()

	exp := ""
	got := fmt.Sprint(row)

	assert(t, exp, got)
}

func TestGroupByValue(t *testing.T) {
	rows, e := initRecords()
	if e != nil {
		t.Fatal(e)
	}

	mapRows := rows.GroupByValue(1)

	got := fmt.Sprint(mapRows)

	assert(t, groupByExpect, got)
}

func TestRandomPick(t *testing.T) {
	for i := 0; i < 5; i++ {
		rows, e := initRecords()
		if e != nil {
			t.Fatal(e)
		}

		shuffledRows := rows.RandomPick(2, true)

		fmt.Println("==> shuffled rows :", shuffledRows)
		fmt.Println("==> remaining rows:", rows)
	}
}
