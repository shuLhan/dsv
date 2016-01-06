// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
Row represent slice of record.
*/
type Row []*Record

/*
PushBack will add new record to the end of row.
*/
func (row *Row) PushBack(r *Record) {
	*row = append(*row, r)
}

/*
GetTypes return type of all records.
*/
func (row *Row) GetTypes() (types []int) {
	for _, r := range *row {
		types = append(types, r.GetType())
	}
	return
}
