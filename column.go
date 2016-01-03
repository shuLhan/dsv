// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
Column represent slice of record. A vertical representation of data.
*/
type Column struct {
	// Name of column. String identifier for the column.
	Name string
	// Type of column. All record in column have the same type.
	Type int
	// Flag additional attribute that can be set to mark some value on this
	// column
	Flag int
	// ValueSpace contain the possible value in records
	ValueSpace []string
	// Records contain column data.
	Records []*Record
}

/*
NewColumn return new column with type and name.
*/
func NewColumn(colType int, colName string) (col *Column) {
	col = &Column{
		Type: colType,
		Name: colName,
		Flag: 0,
	}

	col.Records = make([]*Record, 0)

	return
}

/*
NewColumnString initialize column with type anda data as string.
*/
func NewColumnString(data []string, colType int, colName string) (
	col *Column,
	e error,
) {
	col = NewColumn(colType, colName)

	datalen := len(data)

	if datalen <= 0 {
		return
	}

	col.Records = make([]*Record, datalen)

	for x := 0; x < datalen; x++ {
		rec, e := NewRecord([]byte(data[x]), colType)
		if e != nil {
			return nil, e
		}
		col.Records[x] = rec
	}

	return col, nil
}

/*
NewColumnReal create new column with record type is real.
*/
func NewColumnReal(data []float64, colName string) (col *Column) {
	col = NewColumn(TReal, colName)

	datalen := len(data)

	if datalen <= 0 {
		return
	}

	col.Records = make([]*Record, datalen)

	for x := 0; x < datalen; x++ {
		rec := NewRecordReal(data[x])
		col.Records[x] = rec
	}

	return
}

/*
Reset column data and flag.
*/
func (column *Column) Reset() {
	column.Flag = 0
	column.Records = make([]*Record, 0)
}

/*
GetType return column type.
*/
func (column *Column) GetType() int {
	return column.Type
}

/*
GetName return column name in string.
*/
func (column *Column) GetName() string {
	return column.Name
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
func (column *Column) ToStringSlice() (newcol []string) {
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

/*
SetValues of all column record.
*/
func (column *Column) SetValues(values []string) {
	vallen := len(values)
	reclen := column.GetLength()

	// initialize column record if its empty.
	if reclen <= 0 {
		column.Records = make([]*Record, vallen)
		reclen = vallen
	}

	// pick the least length
	minlen := reclen
	if vallen < reclen {
		minlen = vallen
	}

	for x := 0; x < minlen; x++ {
		column.Records[x].SetValue([]byte(values[x]), column.Type)
	}
}
