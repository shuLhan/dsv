// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

// Fields represent slice of record.
type Field []Record

/*
ToFloatSlice convert slice of record to slice of float64.
*/
func (field *Field) ToFloatSlice() (*[]float64) {
	newf := make([]float64, len(*field))

	for i := range *field {
		newf[i] = (*field)[i].Float()
	}

	return &newf
}

/*
ToStringSlice convert slice of record to slice of string.
*/
func (field *Field)ToStringSlice() (*[]string) {
	newf := make([]string, len(*field))

	for i := range *field {
		newf[i] = (*field)[i].String()
	}

	return &newf
}
