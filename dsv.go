// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package dsv is a library for working with delimited separated value (DSV).

DSV is a free-style form of Comma Separated Value (CSV) format of text data,
where each row is separated by newline, and each column can be separated by
any string enclosed with left-quote and right-quote.
*/
package dsv

import (
	"errors"
	"os"
)

const (
	// DefaultRejected define the default file which will contain the
	// rejected row.
	DefaultRejected		= "rejected.dsv"
	// DefaultMaxRows define default maximum row that will be saved
	// in memory for each read if input data is too large and can not be
	// consumed in one read operation.
	DefaultMaxRows = 256
)

var (
	// ErrNoInput define an error when no Input file is given to Reader.
	ErrNoInput	= errors.New ("dsv: No input file is given in config")
	// ErrMissRecordsLen define an error when trying to push Row
	// to Field, when their length is not equal.
	// See reader.PushRowToColumns().
	ErrMissRecordsLen = errors.New("dsv: Mismatch between number of record in row and columns length")
	// ErrNoOutput define an error when no output file is given to Writer.
	ErrNoOutput	= errors.New ("dsv: No output file is given in config")
	// ErrNotOpen define an error when output file has not been opened
	// by Writer.
	ErrNotOpen	= errors.New ("dsv: Output file is not opened")
	// ErrNilReader define an error when Reader object is nil when passed
	// to Write function.
	ErrNilReader	= errors.New ("dsv: Reader object is nil")
	// ErrUnknownOutputMode will tell you when output mode is unknown.
	ErrUnknownOutputMode = errors.New ("dsv: Unknown output mode")

	// DEBUG exported from environment to debug the library.
	DEBUG		= bool (os.Getenv ("DEBUG") != "")
)

/*
ReadWriter combine reader and writer.
*/
type ReadWriter struct {
	Reader
	Writer
}

/*
New create a new ReadWriter object.
*/
func New () *ReadWriter {
	return &ReadWriter {}
}

/*
SetConfigPath of input and output file.
*/
func (dsv *ReadWriter) SetConfigPath(dir string) {
	dsv.Reader.SetConfigPath(dir)
	dsv.Writer.SetConfigPath(dir)
}

/*
Open configuration file for reading and writing.
*/
func (dsv *ReadWriter) Open(fcfg string) (e error) {
	e = OpenReader(&dsv.Reader, fcfg)

	if e != nil {
		return
	}

	e = OpenWriter(&dsv.Writer, fcfg)

	return e
}

/*
Close reader and writer.
*/
func (dsv *ReadWriter) Close () {
	dsv.Writer.Close ()
	dsv.Reader.Close ()
}
