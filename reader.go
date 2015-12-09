// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

const (
	// TOutputModeRows for output mode in rows.
	TOutputModeRows = 0
	// OutputModeRows is a string representation of output mode rows.
	OutputModeRows = "ROWS"
	// TOutputModeColumns for output mode in columns.
	TOutputModeColumns = 1
	// OutputModeColumns is a string representation of output mode columns.
	OutputModeColumns = "COLUMNS"
	// DefOutputMode default output mode in string.
	DefOutputMode = OutputModeRows
	// DefTOutputMode default output mode.
	DefTOutputMode = TOutputModeRows
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

	dsvReader := dsv.NewReader ()

(2) Open configuration file

	e := dsv.Open (dsvReader, configfile)

(2.1) Do not forget to check for error ...

(3) Make sure to close all files after finished

	defer dsvReader.Close ()

(4) Create loop to read input data

	for {
		n, e := dsv.Read (dsvReader)

		if e == io.EOF {
			break
		}

(4.1) Iterate through rows

		row := dsvReader.Rows.Front ()
		for row != nil {
			// work with row ...

			row = row.Next ()
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
	// NRows define number of rows has been readed and saved.
	NRows		int
	// NColumnIn define number of input columns.
	NColumnIn	int
	// NColumnOut define number of output columns (input - skiped columns)
	NColumnOut	int
	// OutputMode define on how do you want the result is saved. There are
	// two options: either in "rows" mode or "columns" mode.
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
	OutputMode	string		`json:"OutputMode"`
	// TOutputMode define the numeric value of output mode.
	TOutputMode	int		`json:"-"`
	// Columns is input data that has been parsed.
	Columns		Columns		`json:"-"`
	// Rows is input data that has been parsed.
	Rows		Rows		`json:"-"`
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
func NewReader () *Reader {
	return &Reader {
		Input		:"",
		Skip		:0,
		Rejected	:"rejected.dat",
		InputMetadata	:nil,
		MaxRows	:DefaultMaxRows,
		NRows		:0,
		NColumnIn	:0,
		NColumnOut	:0,
		OutputMode	:DefOutputMode,
		TOutputMode	:DefTOutputMode,
		Rows		:Rows{},
		fRead		:nil,
		fReject		:nil,
		bufRead		:nil,
		bufReject	:nil,
	}
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
func (reader *Reader) GetInputMetadata () *[]Metadata {
	return &reader.InputMetadata
}

/*
GetInputMetadataAt return pointer to metadata at index 'idx'.
*/
func (reader *Reader) GetInputMetadataAt (idx int) *Metadata {
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
GetNRows return number of rows that has been read before.
*/
func (reader *Reader) GetNRows() int {
	return reader.NRows
}

/*
SetNRows will set the number of row that has been read.
*/
func (reader *Reader) SetNRows(n int) {
	reader.NRows = n
}

/*
GetOutputMode return output mode of data.
*/
func (reader *Reader) GetOutputMode() string {
	return reader.OutputMode
}

/*
GetTOutputMode return mode of output in integer, so we does not need to
convert it to uppercase to compare it with string.
*/
func (reader *Reader) GetTOutputMode() int {
	return reader.TOutputMode
}

/*
SetOutputMode to `mode`.
*/
func (reader *Reader) SetOutputMode(mode string) error {
	switch strings.ToUpper(mode) {
	case OutputModeRows:
		reader.TOutputMode = TOutputModeRows
		reader.Rows = Rows{}
	case OutputModeColumns:
		reader.TOutputMode = TOutputModeColumns
		reader.Columns = make(Columns, reader.NColumnOut)
	default:
		return ErrUnknownOutputMode
	}
	reader.OutputMode = mode

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
GetNColumnOut return number of column that will be used in output, excluding
the column with Skip=true.
*/
func (reader *Reader) GetNColumnOut() int {
	return reader.NColumnOut;
}

/*
SetNColumnOut set number of output columns.
*/
func (reader *Reader) SetNColumnOut(n int) {
	reader.NColumnOut = n
}

/*
GetData return the output data, based on mode (rows or columns based).
*/
func (reader *Reader) GetData() interface{} {
	switch reader.TOutputMode {
	case TOutputModeRows:
		return reader.Rows
	case TOutputModeColumns:
		return reader.Columns
	}

	return nil
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
	if "" == reader.OutputMode {
		reader.SetOutputMode(DefOutputMode)
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
Reset all variables for next read operation. NRows will be 0, and Rows
will be nil again.
*/
func (reader *Reader) Reset () {
	reader.NRows = 0
	reader.Rows = Rows{}
	reader.Columns = make(Columns, reader.NColumnOut)
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
Push data to rows.
*/
func (reader *Reader) Push(r Row) {
	reader.Rows.PushBack(r)
}

/*
PushRowToColumns push each data in Row to Columns.
*/
func (reader *Reader) PushRowToColumns(row Row) (e error) {
	// check if row length equal with columns length
	if len(row) != len(reader.Columns) {
		return ErrMissRecordsLen
	}

	for i := range (row) {
		reader.Columns[i] = append(reader.Columns[i], row[i])
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
TransposeColumnsToRows will move all data in Columns into Rows mode.
*/
func (reader *Reader) TransposeColumnsToRows () {
	if reader.GetTOutputMode() != TOutputModeColumns {
		return
	}

	rowlen := math.MaxInt32
	flen := len(reader.Columns)

	reader.SetOutputMode(OutputModeRows)

	// Get the least length of columns.
	for f := 0; f < flen; f++ {
		l := len(reader.Columns[f])

		if l < rowlen {
			rowlen = l
		}
	}

	for r := 0; r < rowlen; r++ {
		row := make(Row, flen)

		for f := 0; f < flen; f++ {
			row[f] = reader.Columns[f][r]
		}

		reader.Push(row)
	}

	// reset the columns
	reader.Columns = nil
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
		if ! reader.InputMetadata[a].IsEqual (&other.InputMetadata[a]) {
			return false
		}
	}

	return true
}

/*
RandomPickRows return `n` item of row that has been selected randomly from
reader.Rows. The ids of rows that has been picked is saved id `rowsIdx`.

If duplicate is true, the row that has been picked can be picked up again,
otherwise it only allow one pick. This is also called as random selection with
or without replacement in some machine learning domain.

If output mode is columns, it will be transposed to rows.
*/
func (reader *Reader) RandomPickRows(n int, duplicate bool) (unpicked Rows,
							shuffled Rows,
							pickedIdx []int) {
	if reader.GetTOutputMode() == TOutputModeColumns {
		reader.TransposeColumnsToRows()
	}
	return reader.Rows.RandomPick(n, duplicate)
}

/*
String yes, it will print it in JSON like format.
*/
func (reader *Reader) String() string {
	r, e := json.MarshalIndent (reader, "", "\t")

	if nil != e {
		log.Print (e)
	}

	return string (r)
}
