// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
Column represent slice of record. A vertical representation of data.
*/
type Column struct {
	// Type of column. All record in column have the same type.
	Type int
	// Flag additional attribute that can be set to mark some value on this
	// column
	Flag int
	// Records contain column data.
	Records []*Record
}

/*
Reset column data and flag.
*/
func (column *Column) Reset() {
	column.Flag = 0
	column.Records = make([]*Record, 0)
}

/*
GetLength return number of record.
*/
func (column *Column) GetLength() int {
	return len(column.Records)
}

/*
PushBack push record the end of column.
*/
func (column *Column) PushBack(r *Record) {
	column.Records = append(column.Records, r)
}

/*
ToFloatSlice convert slice of record to slice of float64.
*/
func (column *Column) ToFloatSlice() (newcol []float64) {
	newcol = make([]float64, column.GetLength())

	for i := range column.Records {
		newcol[i] = column.Records[i].Float()
	}

	return
}

/*
ToStringSlice convert slice of record to slice of string.
*/
func (column *Column)ToStringSlice() (newcol []string) {
	newcol = make([]string, column.GetLength())

	for i := range column.Records {
		newcol[i] = column.Records[i].String()
	}

	return
}

/*
ClearValues set all value in column to empty string or zero if column type is
numeric.
*/
func (column *Column) ClearValues() {
	if column.GetLength() <= 0 {
		return
	}

	var v interface{}

	switch column.Type {
	case TString:
		v = ""
	case TInteger:
		v = 0
	case TReal:
		v = 0.0
	}

	for i := range column.Records {
		column.Records[i].V = v
	}
}
