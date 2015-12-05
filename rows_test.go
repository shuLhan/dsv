// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/shuLhan/dsv/util/assert"
)

var exp = []string{
	"0\n",
	"1\n",
	"2\n",
	"3\n",
	"4\n",
}

func TestPushBack(t *testing.T) {
	rows, e := initRows()
	if e != nil {
		t.Fatal(e)
	}

	exp := strings.Join(rowsExpect, "")
	got := fmt.Sprint(rows)

	assert.Equal(t, exp, got)
}

func TestPopFront(t *testing.T) {
	rows, e := initRows()
	if e != nil {
		t.Fatal(e)
	}

	l := len(rows) - 1
	for i := range rows {
		row := rows.PopFront()

		exp := rowsExpect[i]
		got := fmt.Sprint(row)

		assert.Equal(t, exp, got)

		if i < l {
			exp = strings.Join(rowsExpect[i+1:], "")
		} else {
			exp = ""
		}
		got = fmt.Sprint(rows)

		assert.Equal(t, exp, got)
	}

	// empty rows
	row := rows.PopFront()

	exp := "[]"
	got := fmt.Sprint(row)

	assert.Equal(t, exp, got)
}

func TestPopFrontRow(t *testing.T) {
	rows, e := initRows()
	if e != nil {
		t.Fatal(e)
	}

	l := len(rows) - 1
	for i := range rows {
		newRows := rows.PopFrontAsRows()

		exp := rowsExpect[i]
		got := fmt.Sprint(newRows)

		assert.Equal(t, exp, got)

		if i < l {
			exp = strings.Join(rowsExpect[i+1:], "")
		} else {
			exp = ""
		}
		got = fmt.Sprint(rows)

		assert.Equal(t, exp, got)
	}

	// empty rows
	row := rows.PopFrontAsRows()

	exp := ""
	got := fmt.Sprint(row)

	assert.Equal(t, exp, got)
}

func TestGroupByValue(t *testing.T) {
	rows, e := initRows()
	if e != nil {
		t.Fatal(e)
	}

	mapRows := rows.GroupByValue(1)

	got := fmt.Sprint(mapRows)

	assert.Equal(t, groupByExpect, got)
}

func TestRandomPick(t *testing.T) {
	rows, e := initRows()
	if e != nil {
		t.Fatal(e)
	}


	// random pick with duplicate
	for i := 0; i < 5; i++ {
		unpicked, shuffled,idx := rows.RandomPick(6, true)

		// check if unpicked item exist in shuffled items.
		for _, un := range *unpicked {
			for _, pick := range *shuffled {
				if reflect.DeepEqual(un, pick) {
					t.Fatal("random pick: unpicked is false")
				}
			}
		}

		fmt.Println("Random pick with duplicate rows")
		fmt.Println("==> shuffled rows :", shuffled)
		fmt.Println("==> shuffled idx  :", idx)
		fmt.Println("==> unpicked rows :", unpicked)
		fmt.Println("==> original rows :", rows)
	}

	// random pick without duplication
	for i := 0; i < 5; i++ {
		unpicked, shuffled,idx := rows.RandomPick(3, false)

		// check if picked rows is duplicate
		if reflect.DeepEqual((*shuffled)[0], (*shuffled)[1]) {
			t.Fatal("random pick: duplicate rows found.")
		}

		// check if unpicked item exist in shuffled items.
		for _, un := range *unpicked {
			for _, pick := range *shuffled {
				if reflect.DeepEqual(un, pick) {
					t.Fatal("random pick: unpicked is false")
				}
			}
		}

		fmt.Println("Random pick with no duplicate rows")
		fmt.Println("==> shuffled rows :", shuffled)
		fmt.Println("==> shuffled idx  :", idx)
		fmt.Println("==> unpicked rows :", unpicked)
		fmt.Println("==> original rows :", rows)

	}
}
