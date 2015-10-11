package dsv

import (
	"container/list"
	"fmt"
	"math"
	"reflect"
)

/*
Row represent each row of record in linked list model.
*/
type Row struct {
	list.List
}

/*
MapStringRow represent mapping between string (key) with rows (value).
*/
type MapStringRow map[string]*Row

/*
NewRow create new row object.
*/
func NewRow (r interface{}) (row *Row) {
	row = &Row {}
	row.PushBack (r)

	return
}

/*
String return the string of each row separated by new line.
*/
func (rows *Row) String () (s string) {
	row := rows.Front ()
	for nil != row {
		s += fmt.Sprintln (reflect.ValueOf (row.Value))
		row = row.Next ()
	}

	return s
}

/*
PopFront remove the head, return the element value.
*/
func (rows *Row) PopFront () interface{} {
	el := rows.Front ()
	if nil == el {
		return nil
	}

	record := rows.Remove (el)

	return record
}

/*
PopFrontRow cut the head, set the new head to the next element of head, and return
last head.
*/
func (rows *Row) PopFrontRow () *Row {
	record := rows.PopFront ()
	if nil == record {
		return nil
	}

	return NewRow (record)
}

/*
GroupByValue will group each row based on record value in index recGroupIdx
into map of string -> *Row.

WARNING: rows will be modified and will be an empty list.

For example, given rows with target group in field index 1,

	[1 +]
	[2 -]
	[3 -]
	[4 +]

this function will create a map with key is string of target and value is
pointer to sub-rows,

	+ -> [1 +]
	     [4 +]
	- -> [2 -]
	     [3 -]

*/
func (rows *Row) GroupByValue (recGroupIdx int) MapStringRow {
	var row *Row
	var records *[]Record
	var key string
	var v *Row
	var ok bool

	class := make (MapStringRow)

	for {
		row = rows.PopFrontRow ()
		if nil == row {
			break
		}

		records = row.Front ().Value.(*[]Record)
		key = fmt.Sprint ((*records)[recGroupIdx])

		// check if key already mapped.
		v, ok = class[key]

		if ok {
			// push row to the list
			v.PushBackList (&row.List)
		} else {
			// map new key
			class[key] = row
		}
	}

	return class
}

/*
GetMinority return group in groups with minimum rows.
*/
func (groups *MapStringRow) GetMinority () (minorGroup *Row) {
	var min = math.MaxInt32

	for k, v := range *groups {

		if (*groups)[k].Len () < min {
			minorGroup = v
		}
	}

	return
}
