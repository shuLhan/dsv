// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"errors"
	"math/rand"
	"time"
)

var (
	// ErrMisColLength returned when operation on columns does not match
	// between parameter and their length
	ErrMisColLength = errors.New ("dsv: mismatch on column length")
)

/*
Columns represent slice of Column.
*/
type Columns []Column

/*
Reset each data and attribute in all columns.
*/
func (cols *Columns) Reset() {
	for x := range *cols {
		(*cols)[x].Reset()
	}
}

/*
SetType of each column. The length of type must be equal with the number of
column, otherwise it will return error.
*/
func (cols *Columns) SetType(types []int) error {
	if len(types) != len(*cols) {
		return ErrMisColLength
	}
	for x := range *cols {
		(*cols)[x].Type = types[x]
	}
	return nil
}

/*
RandomPick column in columns until n item and return it like its has been
shuffled.  If duplicate is true, column that has been picked can be picked up
again, otherwise it will only picked up once.

This function return picked and unpicked column and index of them.
*/
func (cols *Columns) RandomPick(n int, dup bool, excludeIdx []int) (
	picked Columns,
	unpicked Columns,
	pickedIdx []int,
	unpickedIdx []int,
) {
	excLen := len(excludeIdx)
	colsLen := len(*cols)
	allowedLen := colsLen - excLen

	// if duplication is not allowed, limit the number of selected
	// column.
	if n > allowedLen && !dup {
		n = allowedLen
	}

	rand.Seed(time.Now().UnixNano())

	for ; n >= 1; n-- {
		idx := 0
		for {
			idx = rand.Intn(colsLen)

			// check if its must not be selected
			excluded := false
			for _, excIdx := range excludeIdx {
				if idx == excIdx {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}

			if dup {
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

		picked = append(picked, (*cols)[idx])
	}

	// select unpicked columns using picked index.
	for cid := range *cols {
		// check if column index has been picked up
		isPicked := false
		for _, idx := range pickedIdx {
			if cid == idx {
				isPicked = true
				break
			}
		}
		if !isPicked {
			unpicked = append(unpicked, (*cols)[cid])
			unpickedIdx = append(unpickedIdx, cid)
		}
	}

	return
}
