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
func (record Record) String() string {
	return "\""+ string (record) + "\","
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
func (row *Row) String () (s string) {
	for nil != row {
		s += fmt.Sprintln (row.V)
		row = row.Next
	}

	return s
}
