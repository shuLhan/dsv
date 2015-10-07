package dsv

import (
	"container/list"
	"fmt"
)

/*
Row represent each row of record in linked list model.
*/
type Row struct {
	list.List
}

/*
String return the string of each row separated by new line.
*/
func (rows *Row) String () (s string) {
	row := rows.Front ()
	for nil != row {
		s += fmt.Sprintln (row.Value.(*[]Record))
		row = row.Next ()
	}

	return s
}
