// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package dsv is a library for working with delimited separated value (DSV).

DSV is a free-style form of Comma Separated Value (CSV) format of text data,
where each record is separated by newline, and each field can be separated by
any string enclosed with left-quote and right-quote.
*/
package dsv

import (
	"errors"
	"os"
)

const (
	// DefaultRejected define the default file which will contain the
	// rejected record.
	DefaultRejected		= "rejected.dsv"
	// DefaultMaxRecord define default maximum record that will be saved
	// in memory.
	DefaultMaxRecord	= 256
)

var (
	// ErrNoInput define an error when no Input file is given to Reader.
	ErrNoInput	= errors.New ("dsv: No input file is given in config")
	// ErrNoOutput define an error when no output file is given to Writer.
	ErrNoOutput	= errors.New ("dsv: No output file is given in config")
	// ErrNotOpen define an error when output file has not been opened
	// by Writer.
	ErrNotOpen	= errors.New ("dsv: Output file is not opened")
	// ErrNilReader define an error when Reader object is nil when passed
	// to Write function.
	ErrNilReader	= errors.New ("dsv: Reader object is nil")

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
SetPath of input and output file.
*/
func (dsv *ReadWriter) SetPath (dir string) {
	dsv.Reader.SetPath (dir)
	dsv.Writer.SetPath (dir)
}

/*
Open configuration file for reading and writing.
*/
func (dsv *ReadWriter) Open (fcfg string) (e error) {
	e = Open (&dsv.Reader, fcfg)

	if e != nil {
		return
	}

	e = Open (&dsv.Writer, fcfg)

	return e
}

/*
Close reader and writer.
*/
func (dsv *ReadWriter) Close () {
	dsv.Writer.Close ()
	dsv.Reader.Close ()
}
