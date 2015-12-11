// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

/*
Writer write records from reader or slice using format configuration in
metadata.
*/
type Writer struct {
	Config				`json:"="`
	// Output file where the records will be written.
	Output		string		`json:"Output"`
	// OutputMetadata define format for each column.
	OutputMetadata	[]Metadata	`json:"OutputMetadata"`
	// fWriter as write descriptor.
	fWriter		*os.File
	// BufWriter for buffered writer.
	BufWriter	*bufio.Writer
}

/*
NewWriter create a writer object.
User must call Open after that to populate the output and metadata.
*/
func NewWriter(config string) (writer *Writer, e error) {
	writer = &Writer {
		Output		:"",
		OutputMetadata	:nil,
		fWriter		:nil,
		BufWriter	:nil,
	}

	if config == "" {
		return
	}

	e = OpenWriter(writer, config)
	if e != nil {
		return nil, e
	}

	return
}

/*
GetOutput return output filename.
*/
func (writer *Writer) GetOutput() string {
	return writer.Output
}

/*
SetOutput will set the output file to path.
*/
func (writer *Writer) SetOutput(path string) {
	writer.Output = path
}

/*
OpenOutput file and buffered writer.
*/
func (writer *Writer) OpenOutput() (e error) {
	writer.fWriter, e = os.OpenFile (writer.Output,
					os.O_CREATE | os.O_TRUNC | os.O_WRONLY,
					0600)
	if nil != e {
		return e
	}

	writer.BufWriter = bufio.NewWriter (writer.fWriter)

	return nil
}

/*
Flush output buffer to disk.
*/
func (writer *Writer) Flush() {
	writer.BufWriter.Flush()
}

/*
Close all open descriptor.
*/
func (writer *Writer) Close () {
	if nil != writer.BufWriter {
		writer.BufWriter.Flush ()
	}
	if nil != writer.fWriter {
		writer.fWriter.Close ()
	}
}

/*
WriteRow dump content of Row to file using format in metadata.
*/
func (writer *Writer) WriteRow(row Row, recordMd *[]Metadata) (e error) {
	var md *Metadata
	var inMd *Metadata
	var rIdx int
	var nRecord = len(row)
	var recV []byte
	v := []byte{}

	for i := range writer.OutputMetadata {
		md = &writer.OutputMetadata[i]

		// find the input index based on name on record metadata.
		rIdx = 0
		for y := range (*recordMd) {
			inMd = &(*recordMd)[y]

			if inMd.Name == md.Name {
				break
			}
			if ! (*recordMd)[y].Skip {
				rIdx++
			}
		}

		// If input column is ignored, continue to next record.
		if inMd.Skip {
			continue
		}

		// No input metadata matched? skip it too.
		if rIdx >= nRecord {
			continue
		}

		recV = row[rIdx].ToByte()

		if "" != md.LeftQuote {
			v = append (v, []byte (md.LeftQuote)...)
		}

		v = append (v, recV...)

		if "" != md.RightQuote {
			v = append (v, []byte (md.RightQuote)...)
		}

		if "" != md.Separator {
			v = append (v, []byte (md.Separator)...)
		}
	}

	v = append (v, '\n')

	_, e = writer.BufWriter.Write (v)

	if nil != e {
		return e
	}

	return nil
}

/*
WriteRows will loop each row in the list of rows and write their content to
output file.
Return n for number of row written, and e if error happened.
*/
func (writer *Writer) WriteRows(rows Rows, recordMd *[]Metadata) (n int, e error) {
	n = 0

	for i := range(rows) {
		e = writer.WriteRow(rows[i], recordMd)
		if nil != e {
			if DEBUG {
				log.Println (e)
			}
		}
		n++
	}

	return n,nil
}

/*
WriteColumns will write content of columns to output file.
Return n for number of row written, and e if error happened.
*/
func (writer *Writer) WriteColumns(columns *Columns, md *[]Metadata) (
							n int, e error) {
	nColumns := len(*columns)
	if nColumns <= 0 {
		return
	}

	// Get minimum length of all columns.
	// In case one of the column have different length (shorter or longer),
	// we will take the column with minimum length.
	minLen := len((*columns)[0])

	for i := 1; i < nColumns; i++ {
		l := len((*columns)[i])
		if minLen > l {
			minLen = l
		}
	}

	lenColumn := minLen

	// First loop, iterate over the column length.
	var f int
	row := make(Row, nColumns)

	for r := 0; r < lenColumn; r++ {
		// Second loop, convert columns to record.
		for f = 0; f < nColumns; f++ {
			row[f] = (*columns)[f][r]
		}

		writer.WriteRow(row, md)
	}

	return n,e
}

/*
Write rows from Reader to file.
Return n for number of row written, or e if error happened.
*/
func (writer *Writer) Write (reader *Reader) (int, error) {
	if nil == reader {
		return 0, ErrNilReader
	}
	if nil == writer.fWriter {
		return 0, ErrNotOpen
	}

	switch reader.GetMode() {
	case DatasetModeRows:
		return writer.WriteRows(reader.Rows, &reader.InputMetadata)
	case DatasetModeColumns:
		return writer.WriteColumns(&reader.Columns, &reader.InputMetadata)
	case DatasetModeMatrix:
		return writer.WriteRows(reader.Rows, &reader.InputMetadata)
	}

	return 0, ErrUnknownDatasetMode
}

/*
String yes, it will print it in JSON like format.
*/
func (writer *Writer) String() string {
	r, e := json.MarshalIndent (writer, "", "\t")

	if nil != e {
		log.Print (e)
	}

	return string (r)
}
