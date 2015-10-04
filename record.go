/*
Copyright 2015 Mhd Sulhan <ms@kilabit.info>
All rights reserved.  Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/
package dsv

import (
	"fmt"
)

/*
Record represent each field in record.
*/
type Record []byte

/*
String return the value of record in string enclosed with double quoted.
*/
func (this Record) String() string {
	return "\""+ string (this) + "\","
}

/*
Row represent each row of record in linked list model.
*/
type Row struct {
	V	*[]Record
	Next	*Row
	Last	*Row
}

/*
NewRow create and initialize new row object using 'r' as value for Record.
*/
func NewRow (r *[]Record) *Row {
	return &Row {
		V	:r,
		Next	:nil,
		Last	:nil,
	}
}

/*
String return the string of each row separated by new line.
*/
func (this *Row) String () (s string) {
	for nil != this {
		s += fmt.Sprintln (this.V)
		this = this.Next
	}

	return s
}
