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

/*
Clone create and return a clone of row.
*/
func (row *Row) Clone() (clone Row) {
	for _, rec := range *row {
		newrec := &Record{
			V: rec.V,
		}

		clone.PushBack(newrec)
	}
	return
}

/*
IsNilAt return true if there is no record value in row at `idx`, otherwise
return false.
*/
func (row *Row) IsNilAt(idx int) bool {
	if len(*row) <= idx {
		return true
	}
	if (*row)[idx] == nil {
		return true
	}
	if (*row)[idx].V == nil {
		return true
	}
	return false
}

/*
SetValueAt will set the value of row at cell index `idx` with record `rec`.
*/
func (row *Row) SetValueAt(idx int, rec *Record) {
	(*row)[idx] = rec
}
