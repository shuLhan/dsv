package dsv

import (
	"container/list"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

/*
Row represent each row of record in linked list model.
*/
type Row struct {
	list.List
}

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
	var records *RecordSlice
	var key string
	var v *Row
	var ok bool

	class := make (MapStringRow)

	for {
		row = rows.PopFrontRow ()
		if nil == row {
			break
		}

		records = row.Front ().Value.(*RecordSlice)
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
RandomPick row in rows until n item and return it like has been shuffled.
Row that has been picked will be removed from original rows.
*/
func (rows *Row) RandomPick (n int) (shuffled *Row) {
	var picked int
	var el *list.Element
	var r interface{}
	var i int
	var rowsL = rows.Len ()

	if n > rowsL {
		n = rowsL
	}

	rand.Seed (time.Now ().UnixNano ())

	shuffled = &Row {}

	for ; n >= 1; n-- {
		picked = rand.Intn (rows.Len ())

		el = rows.Front ()
		for i = 1; i < picked; i++ {
			el = el.Next ()
		}

		if el != nil {
			r = rows.Remove (el)
			shuffled.PushBack (r)
		}
	}

	return
}
