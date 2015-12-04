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
	// Output file where the records will be written.
	Output		string		`json:"Output"`
	// OutputPath define where the output file directory belong.
	OutputPath	string
	// OutputMetadata define format for each field.
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
func NewWriter () *Writer {
	return &Writer {
		Output		:"",
		OutputPath	:"",
		OutputMetadata	:nil,
		fWriter		:nil,
		BufWriter	:nil,
	}
}

/*
SetOutput will set the output file to path.
*/
func (writer *Writer) SetOutput(path string) {
	writer.Output = path
}

/*
GetPath of directory where output file reside.
*/
func (writer *Writer) GetPath () string {
	return writer.OutputPath
}

/*
SetPath where output file will be saved.
*/
func (writer *Writer) SetPath (dir string) {
	writer.OutputPath = dir
}

/*
Init initialize writer by opening output file.
*/
func (writer *Writer) Init () error {
	// Exit immediately if no output file is defined in config.
	if "" == writer.Output {
		return ErrNoOutput
	}

	writer.SetOutput(CheckPath(writer, writer.Output))

	return writer.openOutput ()
}

/*
openOutput file and buffered writer.
*/
func (writer *Writer) openOutput () (e error) {
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
WriteRecords dump content of slice to file using metadata format.
*/
func (writer *Writer) WriteRecords(records RecordSlice, recordMd *[]Metadata) (
								e error) {
	var md *Metadata
	var inMd *Metadata
	var rIdx int
	var nRecord = len(records)
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

		// If input field is ignored, continue to next record.
		if inMd.Skip {
			continue
		}

		// No input metadata matched? skip it too.
		if rIdx >= nRecord {
			continue
		}

		recV = records[rIdx].ToByte()

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
Return n for number of records written, and e if error happened when
writing to file.
*/
func (writer *Writer) WriteRows(rows Rows, recordMd *[]Metadata) (n int, e error) {
	n = 0

	for i := range(rows) {
		e = writer.WriteRecords(rows[i], recordMd)
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
WriteFields will write content of fields to output file.
Return n for number of records written, and e if error happened when
writing to file.
*/
func (writer *Writer) WriteFields(fields *[]RecordSlice, md *[]Metadata) (
							n int, e error) {
	nFields := len(*fields)
	if nFields <= 0 {
		return
	}

	// Get minimum length of all fields.
	// In case one of the field have different length (shorter or longer),
	// we will take the field with minimum length.
	minLen := len((*fields)[0])

	for i := 1; i < nFields; i++ {
		l := len((*fields)[i])
		if minLen > l {
			minLen = l
		}
	}

	lenField := minLen

	// First loop, iterate over the field length.
	var f int
	records := make(RecordSlice, nFields)

	for r := 0; r < lenField; r++ {
		// Second loop, convert fields to record.
		for f = 0; f < nFields; f++ {
			records[f] = (*fields)[f][r]
		}

		writer.WriteRecords(records, md)
	}

	return n,e
}

/*
Write records from Reader to file.
Return n for number of records written, and e for error that happened when
writing to file.
*/
func (writer *Writer) Write (reader *Reader) (int, error) {
	if nil == reader {
		return 0, ErrNilReader
	}
	if nil == writer.fWriter {
		return 0, ErrNotOpen
	}

	switch reader.GetOutputMode() {
	case "ROWS":
		return writer.WriteRows(reader.Rows, &reader.InputMetadata)
	case "FIELDS":
		return writer.WriteFields(&reader.Fields, &reader.InputMetadata)
	}

	return 0, ErrUnknownOutputMode
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
