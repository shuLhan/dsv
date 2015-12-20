// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// DatasetModeROWS is a string representation of output mode rows.
	DatasetModeROWS = "ROWS"
	// DatasetModeCOLUMNS is a string representation of output mode columns.
	DatasetModeCOLUMNS = "COLUMNS"
	// DatasetModeMATRIX will save data in rows and columns. This mode will
	// consume more memory that "rows" and "columns" but give greater
	// flexibility when working with data.
	DatasetModeMATRIX = "MATRIX"
)

/*
ErrReader to handle error data and message.
*/
type ErrReader struct {
	// What cause the error?
	What		string
	// InputLine define the line which cause error
	InputLine	[]byte
}

/*
Error to string.
*/
func (e *ErrReader) Error () string {
	return fmt.Sprintf ("dsv: %s '%s'", e.What, e.InputLine)
}

/*
Reader hold all configuration, metadata and input data.

DSV Reader work like this,

(1) Initialize new dsv reader object

	dsvReader, e := dsv.NewReader(configfile)

(2) Do not forget to check for error ...

	if e != nil {
		// handle error
	}

(3) Make sure to close all files after finished

	defer dsvReader.Close ()

(4) Create loop to read input data

	for {
		n, e := dsv.Read (dsvReader)

		if e == io.EOF {
			break
		}

(4.1) Iterate through rows

		for row := range dsvReader.GetData() {
			// work with row ...
		}
	}

Thats it.

*/
type Reader struct {
	// Config define path of configuration file.
	//
	// If the configuration located in other directory, e.g.
	// "../../config.dsv", and the Input option is set with name only, like
	// "input.dat", we assume that its in the same directory where the
	// configuration file belong.
	Config
	// Dataset contains the content of input file after read.
	Dataset
	// Input file, mandatory.
	Input		string		`json:"Input"`
	// Skip n lines from the head.
	Skip		int		`json:"Skip"`
	// Rejected is the file name where row that does not fit
	// with metadata will be saved.
	Rejected	string		`json:"Rejected"`
	// InputMetadata define format for each column in input data.
	InputMetadata	[]Metadata	`json:"InputMetadata"`
	// MaxRows define maximum row that this reader will read and
	// saved in the memory at one read operation.
	// If the value is -1, all rows will read.
	MaxRows	int		`json:"MaxRows"`
	// DatasetMode define on how do you want the result is saved. There are
	// three options: either in "rows", "columns", or "matrix" mode.
	// For example, input data file,
	//
	//	a,b,c
	//	1,2,3
	//
	// "rows" mode is where each line saved in its own slice, resulting
	// in Rows:
	//
	//	[a b c]
	//	[1 2 3]
	//
	// "columns" mode is where each line saved by columns, resulting in
	// Columns:
	//
	//	[a 1]
	//	[b 2]
	//	[c 3]
	//
	// "matrix" mode is where each record saved in their own row and column.
	//
	DatasetMode	string		`json:"DatasetMode"`
	// fRead is read descriptor.
	fRead		*os.File
	// fReject is reject descriptor.
	fReject		*os.File
	// bufRead is a buffer for working with input file.
	bufRead		*bufio.Reader
	// bufReject is a buffer for working with rejected file.
	bufReject	*bufio.Writer
}

/*
NewReader create and initialize new instance of DSV Reader with default values.
*/
func NewReader(config string) (reader *Reader, e error) {
	reader = &Reader {
		Input		:"",
		Skip		:0,
		Rejected	:"rejected.dat",
		InputMetadata	:nil,
		MaxRows		:DefaultMaxRows,
		DatasetMode	:DefDatasetMode,
		fRead		:nil,
		fReject		:nil,
		bufRead		:nil,
		bufReject	:nil,
	}

	if config == "" {
		return
	}

	e = OpenReader(reader, config)
	if e != nil {
		return nil, e
	}

	return
}

/*
GetInput return the input file.
*/
func (reader *Reader) GetInput () string {
	return reader.Input
}

/*
SetInput file.
*/
func (reader *Reader) SetInput (path string) {
	reader.Input = path
}

/*
GetSkip return number of line that will be skipped.
*/
func (reader *Reader) GetSkip () int {
	return reader.Skip
}

/*
SetSkip set number of lines that will be skipped before reading actual data.
*/
func (reader *Reader) SetSkip(n int) {
	reader.Skip = n
}

/*
GetRejected return name of rejected file.
*/
func (reader *Reader) GetRejected () string {
	return reader.Rejected
}

/*
SetRejected file.
*/
func (reader *Reader) SetRejected (path string) {
	reader.Rejected = path
}

/*
GetInputMetadata return pointer to slice of metadata.
*/
func (reader *Reader) GetInputMetadata() []MetadataInterface {
	md := make([]MetadataInterface, len(reader.InputMetadata))
	for i := range reader.InputMetadata {
		md[i] = &reader.InputMetadata[i]
	}

	return md
}

/*
GetInputMetadataAt return pointer to metadata at index 'idx'.
*/
func (reader *Reader) GetInputMetadataAt(idx int) MetadataInterface {
	return &reader.InputMetadata[idx]
}

/*
GetMaxRows return number of maximum rows for reading.
*/
func (reader *Reader) GetMaxRows() int {
	return reader.MaxRows
}

/*
SetMaxRows will set maximum rows that will be read from input file.
*/
func (reader *Reader) SetMaxRows(max int) {
	reader.MaxRows = max
}

/*
GetDatasetMode return output mode of data.
*/
func (reader *Reader) GetDatasetMode() string {
	return reader.DatasetMode
}

/*
SetDatasetMode to `mode`.
*/
func (reader *Reader) SetDatasetMode(mode string) error {
	switch strings.ToUpper(mode) {
	case DatasetModeROWS:
		reader.SetMode(DatasetModeRows)
	case DatasetModeCOLUMNS:
		reader.SetMode(DatasetModeColumns)
	case DatasetModeMATRIX:
		reader.SetMode(DatasetModeMatrix)
	default:
		return ErrUnknownDatasetMode
	}
	reader.DatasetMode = mode

	return nil
}

/*
GetNColumnIn return number of input columns, or number of metadata, including
column with Skip=true.
*/
func (reader *Reader) GetNColumnIn() int {
	return len(reader.InputMetadata)
}

/*
SetDefault options for global config and each metadata.
*/
func (reader *Reader) SetDefault () {
	if "" == reader.Rejected {
		reader.Rejected = DefaultRejected
	}
	if 0 == reader.MaxRows {
		reader.MaxRows = DefaultMaxRows
	}
	if "" == reader.DatasetMode {
		reader.DatasetMode = DefDatasetMode
	}
}

/*
OpenInput open the input file, metadata must have been initialize.
*/
func (reader *Reader) OpenInput () (e error) {
	reader.fRead, e = os.OpenFile (reader.Input, os.O_RDONLY, 0600)
	if nil != e {
		return e
	}

	reader.bufRead = bufio.NewReader (reader.fRead)

	return nil
}

/*
OpenRejected open rejected file, for saving unparseable line.
*/
func (reader *Reader) OpenRejected () (e error) {
	reader.fReject, e = os.OpenFile (reader.Rejected,
					os.O_CREATE | os.O_WRONLY,
					0600)
	if nil != e {
		return e
	}

	reader.bufReject = bufio.NewWriter (reader.fReject)

	return nil
}

/*
Open input and rejected file.
*/
func (reader *Reader) Open() (e error) {
	// do not let file descriptor leaked
	reader.Close()

	e = reader.OpenInput()
	if e != nil {
		return
	}

	e = reader.OpenRejected()

	return
}

/*
SkipLines skip parsing n lines from input file.
The n is defined in the attribute "Skip"
*/
func (reader *Reader) SkipLines () (e error) {
	for i := 0; i < reader.Skip; i++ {
		_, e = reader.ReadLine ()

		if nil != e {
			log.Print ("dsv: ", e)
			return
		}
	}
	return
}

/*
Reset all variables for next read operation. Number of rows will be 0, and
Rows will be empty again.
*/
func (reader *Reader) Reset() {
	reader.Flush()
	reader.Dataset.Reset()
}

/*
Flush all output buffer.
*/
func (reader *Reader) Flush () {
	reader.bufReject.Flush ()
}

/*
ReadLine will read one line from input file.
*/
func (reader *Reader) ReadLine () (line []byte, e error) {
	var read []byte
	stub := true

	// repeat until one full line is read.
	for stub {
		read, stub, e = reader.bufRead.ReadLine ()

		if nil != e {
			return
		}

		line = append (line, read...)
	}

	return
}

/*
Reject the line and save it to the reject file.
*/
func (reader *Reader) Reject (line []byte) {
	reader.bufReject.Write (line)
}

/*
Close all open descriptors.
*/
func (reader *Reader) Close () {
	if nil != reader.bufReject {
		reader.bufReject.Flush ()
	}
	if nil != reader.fReject {
		reader.fReject.Close ()
	}
	if nil != reader.fRead {
		reader.fRead.Close ()
	}
}

/*
IsEqual compare only the configuration and metadata with other instance.
*/
func (reader *Reader) IsEqual (other *Reader) bool {
	if (reader == other) {
		return true
	}
	if (reader.Input != other.Input) {
		return false
	}

	l,r := len (reader.InputMetadata), len (other.InputMetadata)

	if (l != r) {
		return false
	}

	for a := 0; a < l; a++ {
		if ! reader.InputMetadata[a].IsEqual(&other.InputMetadata[a]) {
			return false
		}
	}

	return true
}
