// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"errors"
	"math"
	"strconv"
)

const (
	// TUndefined for undefined type
	TUndefined = -1
	// TString string type.
	TString = 0
	// TInteger integer type (64 bit).
	TInteger = 1
	// TReal float type (64 bit).
	TReal = 2
)

var (
	// ErrRecTypeUndefined show an error if type is unknown.
	ErrRecTypeUndefined = errors.New("dsv: Undefined type")
)

/*
Record represent the smallest building set of data-set.
*/
type Record struct {
	V interface{}
}

/*
NewRecord create new record from string with specific type.
Return record object or error when fail to convert the byte to type.
*/
func NewRecord(v string, t int) (r *Record, e error) {
	r = &Record{}

	e = r.SetValue(v, t)
	if e != nil {
		return nil, e
	}

	return
}

/*
NewRecordReal create new record from float value.
*/
func NewRecordReal(v float64) (r *Record) {
	return &Record{
		V: v,
	}
}

/*
GetType of record.
*/
func (r *Record) GetType() int {
	switch r.V.(type) {
	case int64:
		return TInteger
	case float64:
		return TReal
	}
	return TString
}

/*
SetValue set the record values from string. If value can not be converted
to type, it will return an error.
*/
func (r *Record) SetValue(v string, t int) error {
	switch t {
	case TString:
		r.V = v

	case TInteger:
		i64, e := strconv.ParseInt(v, 10, 64)
		if nil != e {
			return e
		}

		r.V = i64

	case TReal:
		f64, e := strconv.ParseFloat(v, 64)
		if nil != e {
			return e
		}

		r.V = f64
	}
	return nil
}

/*
SetFloat will set the record content with float value and type.
*/
func (r *Record) SetFloat(v float64) {
	r.V = v
}

/*
Value return value of record based on their type.
*/
func (r *Record) Value() interface{} {
	switch r.V.(type) {
	case int64:
		return r.V.(int64)
	case float64:
		return r.V.(float64)
	}

	return r.V.(string)
}

/*
ToByte convert record value to byte.
*/
func (r *Record) ToByte() (b []byte) {
	switch r.V.(type) {
	case string:
		b = []byte(r.V.(string))

	case int64:
		b = []byte(strconv.FormatInt(r.V.(int64), 10))

	case float64:
		b = []byte(strconv.FormatFloat(r.V.(float64), 'f', -1, 64))
	}

	return b
}

/*
IsMissingValue check wether the value is a missing attribute.
If it string the missing value is indicated by character '?'.
If it integer the missing value is indicated by minimum negative integer.
If it real the missing value is indicated by -Inf.
*/
func (r *Record) IsMissingValue() bool {
	switch r.V.(type) {
	case string:
		str := r.V.(string)
		if str == "?" {
			return true
		}

	case int64:
		i64 := r.V.(int64)
		if i64 == math.MinInt64 {
			return true
		}

	case float64:
		f64 := r.V.(float64)
		return math.IsInf(f64, -1)
	}

	return false
}

/*
String convert record value to string.
*/
func (r Record) String() (s string) {
	switch r.V.(type) {
	case string:
		s = r.V.(string)

	case int64:
		s = strconv.FormatInt(r.V.(int64), 10)

	case float64:
		s = strconv.FormatFloat(r.V.(float64), 'f', -1, 64)
	}
	return
}

/*
Float convert given record to float value.
*/
func (r *Record) Float() (f64 float64) {
	var e error

	switch r.V.(type) {
	case string:
		f64, e = strconv.ParseFloat(r.V.(string), 64)

		if nil != e {
			f64 = math.Inf(-1)
		}

	case int64:
		f64 = float64(r.V.(int64))

	case float64:
		f64 = r.V.(float64)
	}

	return
}

/*
IsEqual compare record value with another. Return true if both value are match,
otherwise return false.
*/
func (r *Record) IsEqual(v string) bool {
	if r.String() == v {
		return true
	}
	return false
}
