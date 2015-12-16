// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
Column represent slice of record. A vertical representation of data.
*/
type Column []*Record

/*
Columns represent slice of Column.
*/
type Columns []Column

/*
ToFloatSlice convert slice of record to slice of float64.
*/
func (column *Column) ToFloatSlice() (newcol []float64) {
	newcol = make([]float64, len(*column))

	for i := range *column {
		newcol[i] = (*column)[i].Float()
	}

	return newcol
}

/*
ToStringSlice convert slice of record to slice of string.
*/
func (column *Column)ToStringSlice() (newcol []string) {
	newcol = make([]string, len(*column))

	for i := range *column {
		newcol[i] = (*column)[i].String()
	}

	return newcol
}

/*
ClearValues set all value in column to empty string or zero if column type is
numeric.
*/
func (column *Column) ClearValues() {
	if len(*column) <= 0 {
		return
	}

	var v interface{}

	switch (*column)[0].V.(type) {
	case string:
		v = ""
	case int64:
		v = 0
	case float64:
		v = 0.0
	}

	for i := range *column {
		(*column)[i].V = v
	}
}
