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

const (
	// DefSeparator default separator that will be used if its not given
	// in config file.
	DefSeparator = ","
	// DefOutput file.
	DefOutput = "output.dat"
)

/*
Writer write records from reader or slice using format configuration in
metadata.
*/
type Writer struct {
	Config `json:"-"`
	// Output file where the records will be written.
	Output string `json:"Output"`
	// OutputMetadata define format for each column.
	OutputMetadata []Metadata `json:"OutputMetadata"`
	// fWriter as write descriptor.
	fWriter *os.File
	// BufWriter for buffered writer.
	BufWriter *bufio.Writer
}

/*
NewWriter create a writer object.
User must call Open after that to populate the output and metadata.
*/
func NewWriter(config string) (writer *Writer, e error) {
	writer = &Writer{
		Output:         "",
		OutputMetadata: nil,
		fWriter:        nil,
		BufWriter:      nil,
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
func (writer *Writer) OpenOutput(file string) (e error) {
	if file == "" {
		if writer.Output == "" {
			file = DefOutput
		} else {
			file = writer.Output
		}
	}

	writer.fWriter, e = os.OpenFile(file,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if nil != e {
		return e
	}

	writer.BufWriter = bufio.NewWriter(writer.fWriter)

	return nil
}

/*
Flush output buffer to disk.
*/
func (writer *Writer) Flush() error {
	return writer.BufWriter.Flush()
}

/*
Close all open descriptor.
*/
func (writer *Writer) Close() (e error) {
	if nil != writer.BufWriter {
		e = writer.BufWriter.Flush()
		if e != nil {
			return
		}
	}
	if nil != writer.fWriter {
		e = writer.fWriter.Close()
	}
	return
}

/*
WriteRow dump content of Row to file using format in metadata.
*/
func (writer *Writer) WriteRow(row Row, recordMd []MetadataInterface) (e error) {
	var md Metadata
	var inMd MetadataInterface
	var rIdx int
	var nRecord = len(row)
	var recV []byte
	v := []byte{}

	for i := range writer.OutputMetadata {
		md = writer.OutputMetadata[i]

		// find the input index based on name on record metadata.
		rIdx = 0
		for y := range recordMd {
			inMd = recordMd[y]

			if inMd.GetName() == md.GetName() {
				break
			}
			if !inMd.GetSkip() {
				rIdx++
			}
		}

		// If input column is ignored, continue to next record.
		if inMd.GetSkip() {
			continue
		}

		// No input metadata matched? skip it too.
		if rIdx >= nRecord {
			continue
		}

		recV = row[rIdx].ToByte()

		if "" != md.GetLeftQuote() {
			v = append(v, []byte(md.GetLeftQuote())...)
		}

		v = append(v, recV...)

		if "" != md.GetRightQuote() {
			v = append(v, []byte(md.GetRightQuote())...)
		}

		if "" != md.GetSeparator() {
			v = append(v, []byte(md.GetSeparator())...)
		}
	}

	v = append(v, '\n')

	_, e = writer.BufWriter.Write(v)

	return e
}

/*
WriteRows will loop each row in the list of rows and write their content to
output file.
Return n for number of row written, and e if error happened.
*/
func (writer *Writer) WriteRows(rows Rows, recordMd []MetadataInterface) (
	n int,
	e error,
) {
	n = 0

	for i := range rows {
		e = writer.WriteRow(rows[i], recordMd)
		if nil != e {
			if DEBUG {
				log.Println(e)
			}
		}
		n++
	}

	return n, nil
}

/*
WriteColumns will write content of columns to output file.
Return n for number of row written, and e if error happened.
*/
func (writer *Writer) WriteColumns(columns *Columns, md []MetadataInterface) (
	n int,
	e error,
) {
	nColumns := len(*columns)
	if nColumns <= 0 {
		return
	}

	// Get minimum length of all columns.
	// In case one of the column have different length (shorter or longer),
	// we will take the column with minimum length.
	minLen := (*columns)[0].Len()

	for i := 1; i < nColumns; i++ {
		l := (*columns)[i].Len()
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
			row[f] = (*columns)[f].Records[r]
		}

		e = writer.WriteRow(row, md)
		if e != nil {
			break
		}
	}

	return n, e
}

/*
WriteRawRows write rows data using separator `sep` for each record.
*/
func (writer *Writer) WriteRawRows(rows *Rows, sep string) (nrow int, e error) {
	nrow = len(*rows)

	if nrow <= 0 {
		return
	}

	var v []byte

	for x := 0; x < nrow; x++ {
		v = []byte{}
		for y, rec := range (*rows)[x] {
			if y > 0 {
				v = append(v, []byte(sep)...)
			}
			v = append(v, rec.ToByte()...)
		}

		v = append(v, '\n')

		_, e = writer.BufWriter.Write(v)

		if nil != e {
			return 0, e
		}
	}

	e = writer.Flush()
	return
}

/*
WriteRawColumns write raw columns using separator `sep` for each record to
file.
*/
func (writer *Writer) WriteRawColumns(cols *Columns, sep string) (
	nrow int,
	e error,
) {
	ncol := len(*cols)
	if ncol <= 0 {
		return
	}

	nrow = (*cols)[0].Len()
	if nrow <= 0 {
		return
	}

	var v []byte

	for x := 0; x < nrow; x++ {
		v = []byte{}
		for y := 0; y < ncol; y++ {
			if y > 0 {
				v = append(v, []byte(sep)...)
			}
			v = append(v, (*cols)[y].Records[x].ToByte()...)
		}

		v = append(v, '\n')

		_, e = writer.BufWriter.Write(v)

		if nil != e {
			return 0, e
		}
	}

	e = writer.Flush()
	return
}

/*
WriteDataset will write content of dataset to file without metadata but using
separator `sep` for each record.

We use pointer in separator parameter, so we can use empty string as separator.
*/
func (writer *Writer) WriteDataset(dataset *Dataset, sep *string) (int, error) {
	if nil == writer.fWriter {
		return 0, ErrNotOpen
	}
	if nil == dataset {
		return 0, nil
	}
	if sep == nil {
		sep = new(string)
		*sep = DefSeparator
	}

	switch dataset.GetMode() {
	case DatasetModeRows, DatasetModeMatrix:
		return writer.WriteRawRows(&dataset.Rows, *sep)

	case DatasetModeColumns:
		return writer.WriteRawColumns(&dataset.Columns, *sep)
	}

	return 0, ErrUnknownDatasetMode
}

/*
Write rows from Reader to file.
Return n for number of row written, or e if error happened.
*/
func (writer *Writer) Write(reader *Reader) (int, error) {
	if nil == reader {
		return 0, ErrNilReader
	}
	if nil == writer.fWriter {
		return 0, ErrNotOpen
	}

	switch reader.GetMode() {
	case DatasetModeRows:
		return writer.WriteRows(reader.Rows, reader.GetInputMetadata())
	case DatasetModeColumns:
		return writer.WriteColumns(&reader.Columns,
			reader.GetInputMetadata())
	case DatasetModeMatrix:
		return writer.WriteRows(reader.Rows, reader.GetInputMetadata())
	}

	return 0, ErrUnknownDatasetMode
}

/*
String yes, it will print it in JSON like format.
*/
func (writer *Writer) String() string {
	r, e := json.MarshalIndent(writer, "", "\t")

	if nil != e {
		log.Print(e)
	}

	return string(r)
}
