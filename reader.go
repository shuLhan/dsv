// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bufio"
	"github.com/shuLhan/tabula"
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
	tabula.Dataset
	// Input file, mandatory.
	Input string `json:"Input"`
	// Skip n lines from the head.
	Skip int `json:"Skip"`
	// TrimSpace or not. If its true, before parsing the line, the white
	// space in the beginning and end of each input line will be removed,
	// otherwise it will leave unmodified.  Default is true.
	TrimSpace bool `json:"TrimSpace"`
	// Rejected is the file name where row that does not fit
	// with metadata will be saved.
	Rejected string `json:"Rejected"`
	// InputMetadata define format for each column in input data.
	InputMetadata []Metadata `json:"InputMetadata"`
	// MaxRows define maximum row that this reader will read and
	// saved in the memory at one read operation.
	// If the value is -1, all rows will read.
	MaxRows int `json:"MaxRows"`
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
	DatasetMode string `json:"DatasetMode"`
	// fRead is read descriptor.
	fRead *os.File
	// fReject is reject descriptor.
	fReject *os.File
	// bufRead is a buffer for working with input file.
	bufRead *bufio.Reader
	// bufReject is a buffer for working with rejected file.
	bufReject *bufio.Writer
}

/*
NewReader create and initialize new instance of DSV Reader with default values.
*/
func NewReader(config string) (reader *Reader, e error) {
	reader = &Reader{
		Input:         "",
		Skip:          0,
		TrimSpace:     true,
		Rejected:      "rejected.dat",
		InputMetadata: nil,
		MaxRows:       DefaultMaxRows,
		DatasetMode:   DefDatasetMode,
		fRead:         nil,
		fReject:       nil,
		bufRead:       nil,
		bufReject:     nil,
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
CopyConfig copy configuration from other reader object not including data
and metadata.
*/
func (reader *Reader) CopyConfig(src *Reader) {
	reader.ConfigPath = src.GetConfigPath()
	reader.Input = src.GetInput()
	reader.Skip = src.GetSkip()
	reader.TrimSpace = src.IsTrimSpace()
	reader.Rejected = src.GetRejected()
	reader.MaxRows = src.GetMaxRows()
	reader.DatasetMode = src.GetDatasetMode()
}

/*
GetInput return the input file.
*/
func (reader *Reader) GetInput() string {
	return reader.Input
}

/*
SetInput file.
*/
func (reader *Reader) SetInput(path string) {
	reader.Input = path
}

/*
GetSkip return number of line that will be skipped.
*/
func (reader *Reader) GetSkip() int {
	return reader.Skip
}

/*
SetSkip set number of lines that will be skipped before reading actual data.
*/
func (reader *Reader) SetSkip(n int) {
	reader.Skip = n
}

/*
IsTrimSpace return value of TrimSpace option.
*/
func (reader *Reader) IsTrimSpace() bool {
	return reader.TrimSpace
}

/*
GetRejected return name of rejected file.
*/
func (reader *Reader) GetRejected() string {
	return reader.Rejected
}

/*
SetRejected file.
*/
func (reader *Reader) SetRejected(path string) {
	reader.Rejected = path
}

/*
AddInputMetadata add new input metadata to reader.
*/
func (reader *Reader) AddInputMetadata(md *Metadata) {
	reader.InputMetadata = append(reader.InputMetadata, *md)
	reader.AddColumn(md.GetType(), md.GetName(), md.GetValueSpace())
}

/*
AppendMetadata will append new metadata `md` to list of reader input metadata.
*/
func (reader *Reader) AppendMetadata(mdi MetadataInterface) {
	md := mdi.(*Metadata)
	reader.InputMetadata = append(reader.InputMetadata, *md)
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
		reader.SetMode(tabula.DatasetModeRows)
	case DatasetModeCOLUMNS:
		reader.SetMode(tabula.DatasetModeColumns)
	case DatasetModeMATRIX:
		reader.SetMode(tabula.DatasetModeMatrix)
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
func (reader *Reader) SetDefault() {
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
func (reader *Reader) OpenInput() (e error) {
	reader.fRead, e = os.OpenFile(reader.Input, os.O_RDONLY, 0600)
	if nil != e {
		return e
	}

	reader.bufRead = bufio.NewReader(reader.fRead)

	return nil
}

/*
OpenRejected open rejected file, for saving unparseable line.
*/
func (reader *Reader) OpenRejected() (e error) {
	reader.fReject, e = os.OpenFile(reader.Rejected,
		os.O_CREATE|os.O_WRONLY, 0600)
	if nil != e {
		return e
	}

	reader.bufReject = bufio.NewWriter(reader.fReject)

	return nil
}

/*
Open input and rejected file.
*/
func (reader *Reader) Open() (e error) {
	// do not let file descriptor leaked
	e = reader.Close()
	if e != nil {
		return
	}

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
func (reader *Reader) SkipLines() (e error) {
	for i := 0; i < reader.Skip; i++ {
		_, e = reader.ReadLine()

		if nil != e {
			log.Print("dsv: ", e)
			return
		}
	}
	return
}

/*
Reset all variables for next read operation. Number of rows will be 0, and
Rows will be empty again.
*/
func (reader *Reader) Reset() (e error) {
	e = reader.Flush()
	if e != nil {
		return
	}
	e = reader.Dataset.Reset()
	return
}

/*
Flush all output buffer.
*/
func (reader *Reader) Flush() error {
	return reader.bufReject.Flush()
}

/*
ReadLine will read one line from input file.
*/
func (reader *Reader) ReadLine() (line []byte, e error) {
	line, e = reader.bufRead.ReadBytes(DefEOL)

	if e == nil {
		// remove EOL
		line = line[:len(line)-1]
	}

	return
}

/*
FetchNextLine read the next line and combine it with the `lastline`.
*/
func (reader *Reader) FetchNextLine(lastline []byte) (line []byte, e error) {
	line, e = reader.ReadLine()

	lastline = append(lastline, DefEOL)
	lastline = append(lastline, line...)

	return lastline, e
}

/*
Reject the line and save it to the reject file.
*/
func (reader *Reader) Reject(line []byte) (int, error) {
	return reader.bufReject.Write(line)
}

/*
Close all open descriptors.
*/
func (reader *Reader) Close() (e error) {
	if nil != reader.bufReject {
		e = reader.bufReject.Flush()
		if e != nil {
			return
		}
	}
	if nil != reader.fReject {
		e = reader.fReject.Close()
		if e != nil {
			return
		}
	}
	if nil != reader.fRead {
		e = reader.fRead.Close()
	}
	return
}

/*
IsEqual compare only the configuration and metadata with other instance.
*/
func (reader *Reader) IsEqual(other *Reader) bool {
	if reader == other {
		return true
	}
	if reader.Input != other.Input {
		return false
	}

	l, r := len(reader.InputMetadata), len(other.InputMetadata)

	if l != r {
		return false
	}

	for a := 0; a < l; a++ {
		if !reader.InputMetadata[a].IsEqual(&other.InputMetadata[a]) {
			return false
		}
	}

	return true
}

/*
GetDataset return reader dataset.
*/
func (reader *Reader) GetDataset() tabula.DatasetInterface {
	return &reader.Dataset
}

/*
MergeColumns append metadata and columns from another reader if not exist in
current metadata set.
*/
func (reader *Reader) MergeColumns(other ReaderInterface) {
	for _, md := range other.GetInputMetadata() {
		if md.GetSkip() {
			continue
		}

		// Check if the same metadata name exist in current dataset.
		found := false
		for _, lmd := range reader.GetInputMetadata() {
			if lmd.GetName() == md.GetName() {
				found = true
				break
			}
		}

		if found {
			continue
		}

		reader.AppendMetadata(md)
	}

	reader.GetDataset().MergeColumns(other.GetDataset())
}

/*
MergeRows append rows from another reader.
*/
func (reader *Reader) MergeRows(other *Reader) {
	reader.Dataset.MergeRows(other.Dataset)
}
