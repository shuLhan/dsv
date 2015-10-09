package dsv

import (
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

/*
RecordNew create new record from byte with specific type.
Return record object or error when fail to convert the byte to type.
*/
func RecordNew (v []byte, t int) (r Record, e error) {
	s := string (v)

	r = Record {}

	switch t {
	case TString:
		r.V = s

	case TInteger:
		i64, e := strconv.ParseInt (s, 10, 64)

		if nil != e {
			return r, e
		}

		r.V = i64

	case TReal:
		f64, e := strconv.ParseFloat (s, 64)

		if nil != e {
			return r, e
		}

		r.V = f64
	}

	r.T = t

	return r, nil
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
