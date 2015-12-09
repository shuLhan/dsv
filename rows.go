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
type Row []*Record

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
If duplicate is true, row that has been picked can be picked up again,
otherwise it will only picked up once.
*/
func (rows *Rows) RandomPick(n int, duplicate bool) (unpicked Rows,
							shuffled Rows,
							pickedIdx []int) {
	rowsLen := len(*rows)

	// since duplication is not allowed, we can only select as many as rows
	// that we have.
	if n > rowsLen && !duplicate {
		n = rowsLen
	}

	rand.Seed(time.Now().UnixNano())

	for ; n >= 1; n-- {
		idx := 0
		for {
			idx = rand.Intn(len(*rows))

			if duplicate {
				// allow duplicate idx
				pickedIdx = append(pickedIdx, idx)
				break
			}

			// check if its already picked
			isPicked := false
			for _, pastIdx := range pickedIdx {
				if idx == pastIdx {
					isPicked = true
					break
				}
			}
			// get another random idx again
			if isPicked {
				continue
			}

			// bingo, we found unique idx that has not been picked.
			pickedIdx = append(pickedIdx, idx)
			break
		}

		row := (*rows)[idx]

		shuffled.PushBack(row)
	}

	// select unpicked rows using picked index.
	for rid := range *rows {
		// check if row index has been picked up
		isPicked := false
		for _, idx := range pickedIdx {
			if rid == idx {
				isPicked = true
				break
			}
		}
		if !isPicked {
			unpicked.PushBack((*rows)[rid])
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
