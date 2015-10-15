// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"math"
	"strconv"
)

const (
	// TString string type.
	TString		= 0
	// TInteger integer type (64 bit).
	TInteger	= 1
	// TReal float type (64 bit).
	TReal		= 2
)

/*
Record represent each field in record.
*/
type Record struct {
	V interface{}
	T int
}

// RecordSlice record with slice.
type RecordSlice []Record

/*
SetByte set the record values from bytes.
*/
func (r *Record) SetByte (v []byte, t int) (e error) {
	var i64 int64
	var f64 float64
	s := string (v)

	switch t {
	case TString:
		r.V = s

	case TInteger:
		i64, e = strconv.ParseInt (s, 10, 64)

		if nil != e {
			r.V = int64 (math.MinInt64)
		} else {
			r.V = i64
		}

	case TReal:
		f64, e = strconv.ParseFloat (s, 64)

		if nil != e {
			r.V = math.Inf (-1)
		} else {
			r.V = f64
		}
	}

	r.T = t

	// keep returning error in case convert is failed and missing value
	// is set.
	return
}

/*
SetFloat will set the record content with float value and type.
*/
func (r *Record) SetFloat (v float64) {
	r.V = v
	r.T = TReal
}

/*
RecordNew create new record from byte with specific type.
Return record object or error when fail to convert the byte to type.
*/
func RecordNew (v []byte, t int) (r Record, e error) {
	r = Record {}
	e = (&r).SetByte (v, t)

	return
}

/*
Value return value of record based on their type.
*/
func (r *Record) Value () interface{} {
	switch r.T {
	case TInteger:
		return r.V.(int64)
	case TReal:
		return r.V.(float64)
	}

	return r.V.(string)
}

/*
ToByte convert record value to byte.
*/
func (r *Record) ToByte () (b []byte) {
	switch r.T {
	case TString:
		b = []byte (r.V.(string))

	case TInteger:
		b = []byte (strconv.FormatInt (r.V.(int64), 10))

	case TReal:
		b = []byte (strconv.FormatFloat (r.V.(float64), 'f', -1, 64))
	}

	return b
}

/*
IsMissingValue check wether the value is a missing attribute.
If it string the missing value is indicated by character '?'.
If it integer the missing value is indicated by minimum negative integer.
If it real the missing value is indicated by -Inf.
*/
func (r *Record) IsMissingValue () bool {
	switch r.T {
	case TString:
		str := r.V.(string)
		if str == "?" {
			return true
		}

	case TInteger:
		i64 := r.V.(int64)
		if i64 == math.MinInt64 {
			return true
		}

	case TReal:
		f64 := r.V.(float64)
		return math.IsInf (f64, -1)
	}

	return false
}

/*
String convert record value to string.
*/
func (r Record) String () (s string) {
	switch r.T {
	case TString:
		s = r.V.(string)

	case TInteger:
		s = strconv.FormatInt (r.V.(int64), 10)

	case TReal:
		s = strconv.FormatFloat (r.V.(float64), 'f', -1, 64)
	}
	return
}
