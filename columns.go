// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"github.com/shuLhan/dsv/util"
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
column, otherwise it will used the minimum length between types or columns.
*/
func (cols *Columns) SetType(types []int) {
	typeslen := len(types)
	colslen := len(*cols)
	minlen := typeslen

	if colslen < minlen {
		minlen = colslen
	}

	for x := 0; x < minlen; x++ {
		(*cols)[x].Type = types[x]
	}
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

	for ; n >= 1; n-- {
		idx := util.GetRandomInteger(colsLen, dup, pickedIdx,
			excludeIdx)

		pickedIdx = append(pickedIdx, idx)
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
