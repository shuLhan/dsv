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
Reader hold all configuration, metadata and input records.

DSV Reader work like this,

(1) Initialize new dsv reader object

	dsvReader := dsv.NewReader ()

(2) Open configuration file

	e := dsv.Open (dsvReader, configfile)

(2.1) Do not forget to check for error ...

(3) Make sure to close all files after finished

	defer dsvReader.Close ()

(4) Create loop to read and process records.

	for {
		n, e := dsv.Read (dsvReader)

		if e == io.EOF {
			break
		}

(4.1) Iterate through records

		row := dsvReader.Rows.Front ()
		for row != nil {
			// work with records ...

			row = row.Next ()
		}
	}

Thats it.

*/
type Reader struct {
	// Input file, mandatory.
	Input		string		`json:"Input"`
	// InputPath of input file, relative to the path of config file.
	//
	// If the configuration located in other directory, e.g.
	// "../../config.dsv", and the Input option is set with name only, like
	// "input.dat", we assume that its in the same directory where the
	// configuration file belong.
	InputPath	string
	// Skip n lines from the head.
	Skip		int		`json:"Skip"`
	// Rejected is the file where record that does not fit
	// with metadata will be saved.
	Rejected	string		`json:"Rejected"`
	// InputMetadata define format each column in a record.
	InputMetadata	[]Metadata	`json:"InputMetadata"`
	// MaxRecord define maximum record that this reader will read and
	// saved in the memory at one read operation.
	// If the value is -1 all records will read.
	MaxRecord	int		`json:"MaxRecord"`
	// NRecord define number of record readed and saved in Rows.
	NRecord		int
	// NColumnIn define number of input columns.
	NColumnIn	int		`json:"-"`
	// NColumnOut define number of output columns (input - skiped columns)
	NColumnOut	int		`json:"-"`
	// RecordMode define on how do you want the resulting record. There are
	// two options: either in "rows" mode or "columns" mode.
	// For example, input data file,
	//
	//	a,b,c
	//	1,2,3
	//
	// Row mode is where each line saved in linked-list of row, resulting
	// in Rows:
	//
	//	[a b c]->[1 2 3]
	//
	// Column mode is where each line saved by columns, resulting in Columns:
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
	// fRead as read descriptor.
	fRead		*os.File
	// fReject as reject descriptor.
	fReject		*os.File
	// bufRead is for working with input file.
	bufRead		*bufio.Reader
	// bufReject is for rejected records.
	bufReject	*bufio.Writer
}

/*
NewReader create and initialize new instance of DSV Reader with default values.
*/
func NewReader () *Reader {
	return &Reader {
		Input		:"",
		InputPath	:"",
		Skip		:0,
		Rejected	:"rejected.dat",
		InputMetadata	:nil,
		MaxRecord	:DefaultMaxRecord,
		NRecord		:0,
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
GetPath return the base path of configuration file.
*/
func (reader *Reader) GetPath () string {
	return reader.InputPath
}

/*
SetPath for reading input and writing rejected file.
*/
func (reader *Reader) SetPath (dir string) {
	reader.InputPath = dir
}

/*
GetSkip return number of line that will be skipped.
*/
func (reader *Reader) GetSkip () int {
	return reader.Skip
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
GetMaxRecord return number of maximum record for reading.
*/
func (reader *Reader) GetMaxRecord () int {
	return reader.MaxRecord
}

/*
SetMaxRecord will set maximum record that will be read from input file.
*/
func (reader *Reader) SetMaxRecord (max int) {
	reader.MaxRecord = max
}

/*
GetRecordRead return number of record that has been read before.
*/
func (reader *Reader) GetRecordRead () int {
	return reader.NRecord
}

/*
SetRecordRead will set the number of record that has been read.
*/
func (reader *Reader) SetRecordRead (n int) {
	reader.NRecord = n
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
GetNColumnOut return number of column that will be used in output, excluding
the column with Skip=true.
*/
func (reader *Reader) GetNColumnOut() int {
	return reader.NColumnOut;
}

/*
GetOutput return the output records, based on mode (rows or columns based).
*/
func (reader *Reader) GetOutput() interface{} {
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
	if 0 == reader.MaxRecord {
		reader.MaxRecord = DefaultMaxRecord
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
Init initialize reader object by opening input and rejected files and
skip n lines from input.
*/
func (reader *Reader) Init () (e error) {
	// Exit immediately if no input file is defined in config.
	if "" == reader.Input {
		return ErrNoInput
	}

	// Set number of input columns.
	reader.NColumnIn = len(reader.InputMetadata)

	// Check and initialize metadata.
	for i := range reader.InputMetadata {
		e = reader.InputMetadata[i].Init ()

		// Count number of output columns.
		if ! reader.InputMetadata[i].Skip {
			reader.NColumnOut++
		}

		if nil != e {
			return e
		}
	}

	// Set default value
	reader.SetDefault ()

	// Check if output mode is valid and initialize it if valid.
	e = reader.SetOutputMode(reader.OutputMode)

	if nil != e {
		return
	}

	// Check if Input is name only without path, so we can prefix it with
	// config path.
	reader.SetInput (CheckPath (reader, reader.GetInput ()))
	reader.SetRejected (CheckPath (reader, reader.GetRejected ()))

	// Get ready ...
	e = reader.OpenInput ()

	if nil != e {
		return
	}

	e = reader.OpenRejected ()

	if nil != e {
		return
	}

	// Skip lines
	if reader.Skip > 0 {
		e = reader.SkipLines ()

		if nil != e {
			return
		}
	}

	return
}

/*
Reset all variables for next read operation. NRecord will be 0, and Rows
will be nil again.
*/
func (reader *Reader) Reset () {
	reader.NRecord = 0
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
Push record to row.
*/
func (reader *Reader) Push(r Row) {
	reader.Rows.PushBack(r)
}

/*
PushRowToColumns push each record in Row to Columns.
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
TransposeColumnsToRows will move all record in Columns into Rows.
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
String yes, it will print it in JSON like format.
*/
func (reader *Reader) String() string {
	r, e := json.MarshalIndent (reader, "", "\t")

	if nil != e {
		log.Print (e)
	}

	return string (r)
}
