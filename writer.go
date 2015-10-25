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
func (writer *Writer) WriteRecords (records *RecordSlice) (e error) {
	var md *Metadata
	var r *Record
	v := []byte{}

	for i := range writer.OutputMetadata {
		md = &writer.OutputMetadata[i]
		r = &(*records)[i]

		// no more record?
		if nil == r {
			break
		}

		if "" != md.LeftQuote {
			v = append (v, []byte (md.LeftQuote)...)
		}

		v = append (v, r.ToByte ()...)

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
WriteRows will loop each row in the list of rows.
Return n for number of records written, and e for error that happened when
writing to file.
*/
func (writer *Writer) WriteRows (rows *Row) (n int, e error) {
	n = 0
	row := rows.Front ()

	for nil != row {
		e = writer.WriteRecords (row.Value.(*RecordSlice))
		if nil != e {
			if DEBUG {
				log.Println (e)
			}
		}
		row = row.Next ()
		n++
	}

	return n,nil
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

	return writer.WriteRows (reader.Records)
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
