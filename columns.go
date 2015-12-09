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
func (column *Column) ToFloatSlice() (*[]float64) {
	newf := make([]float64, len(*column))

	for i := range *column {
		newf[i] = (*column)[i].Float()
	}

	return &newf
}

/*
ToStringSlice convert slice of record to slice of string.
*/
func (column *Column)ToStringSlice() (*[]string) {
	newf := make([]string, len(*column))

	for i := range *column {
		newf[i] = (*column)[i].String()
	}

	return &newf
}
