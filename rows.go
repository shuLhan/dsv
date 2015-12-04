// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Row represent slice of record.
*/
type Row []Record

/*
Rows represent slice of Row.
*/
type Rows []Row

/*
PushBack append record r to the end of rows.
*/
func (rows *Rows) PushBack(r Row) {
	if r != nil {
		(*rows) = append((*rows), r)
	}
}

/*
PopFront remove the head, return the record value.
*/
func (rows *Rows) PopFront() (row Row) {
	l := len(*rows)
	if l > 0 {
		row = (*rows)[0]
		(*rows) = (*rows)[1:]
	}
	return
}

/*
PopFrontAsRows remove the head and return ex-head as new rows.
*/
func (rows *Rows) PopFrontAsRows() (newRows Rows) {
	row := rows.PopFront()
	if nil == row {
		return
	}
	newRows.PushBack(row)
	return
}

/*
GroupByValue will group each row based on record value in index recGroupIdx
into map of string -> *Row.

WARNING: returned rows will be empty!

For example, given rows with target group in column index 1,

	[1 +]
	[2 -]
	[3 -]
	[4 +]

this function will create a map with key is string of target and value is
pointer to sub-rows,

	+ -> [1 +]
             [4 +]
	- -> [2 -]
             [3 -]

*/
func (rows *Rows) GroupByValue(GroupIdx int) (mapRows MapRows) {
	for {
		row := rows.PopFront()
		if nil == row {
			break
		}

		key := fmt.Sprint(row[GroupIdx])

		mapRows.AddRow(key, row)
	}
	return
}

/*
RandomPick row in rows until n item and return it like its has been shuffled.
If remove is true, row that has been picked will be removed from rows,
otherwise it will stay there and can be picked up again.
*/
func (rows *Rows) RandomPick(n int, remove bool) (shuffled Rows) {
	rowsLen := len(*rows)

	if n > rowsLen {
		n = rowsLen
	}

	rand.Seed(time.Now().UnixNano())

	for ; n >= 1; n-- {
		picked := rand.Intn(len(*rows))

		row := (*rows)[picked]

		shuffled.PushBack(row)

		if remove {
			(*rows) = append((*rows)[:picked], (*rows)[picked+1:]...)
		}
	}
	return
}

/*
String return the string representation of each row separated by new line.
*/
func (rows Rows) String() (s string) {
	for x := range rows {
		s += fmt.Sprint(rows[x])
	}
	return
}
