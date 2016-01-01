// Copyright 2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

/*
ReaderInterface is the interface for reading DSV file.
*/
type ReaderInterface interface {
	ConfigInterface
	DatasetInterface
	GetInputMetadata() []MetadataInterface
	GetInputMetadataAt(idx int) MetadataInterface
	GetMaxRows() int
	SetMaxRows(max int)
	GetDatasetMode() string
	SetDatasetMode(mode string) error
	GetNColumnIn() int
	GetInput() string
	SetInput(path string)
	GetRejected() string
	SetRejected(path string)
	GetSkip() int
	SetSkip(n int)
	SetDefault()
	OpenInput() error
	OpenRejected() error
	SkipLines() error

	Flush()
	ReadLine() ([]byte, error)
	Reject(line []byte)
	Close()
}

/*
OpenReader configuration file and initialize the attributes.
*/
func OpenReader(reader ReaderInterface, fcfg string) (e error) {
	e = ConfigOpen(reader, fcfg)
	if e != nil {
		return e
	}

	return InitReader(reader)
}

/*
InitReader initialize reader object by opening input and rejected files and
skip n lines from input.
*/
func InitReader(reader ReaderInterface) (e error) {
	// Exit immediately if no input file is defined in config.
	if "" == reader.GetInput() {
		return ErrNoInput
	}

	md := reader.GetInputMetadata()
	nColOut := 0
	var types []int
	var names []string

	// Check and initialize metadata.
	for i := range md {
		e = md[i].Init()

		if nil != e {
			return e
		}

		// Count number of output columns.
		if !md[i].GetSkip() {
			nColOut++
			// add type of metadata to list of type
			types = append(types, md[i].GetType())
			names = append(names, md[i].GetName())
		}
	}

	// Set default value
	reader.SetDefault()

	// Set column type in dataset
	e = reader.SetColumnsType(types)
	if e != nil {
		return
	}

	reader.SetColumnsName(names)

	// Check if output mode is valid and initialize it if valid.
	e = reader.SetDatasetMode(reader.GetDatasetMode())

	if nil != e {
		return
	}

	// Check if Input is name only without path, so we can prefix it with
	// config path.
	reader.SetInput(ConfigCheckPath(reader, reader.GetInput()))
	reader.SetRejected(ConfigCheckPath(reader, reader.GetRejected()))

	// Get ready ...
	e = reader.OpenInput()
	if nil != e {
		return
	}

	e = reader.OpenRejected()
	if nil != e {
		return
	}

	// Skip lines
	if reader.GetSkip() > 0 {
		e = reader.SkipLines()

		if nil != e {
			return
		}
	}

	return
}

/*
Read row from input file.
*/
func Read(reader ReaderInterface) (n int, e error) {
	maxrows := reader.GetMaxRows()
	reader.Reset()

	// remember to flush if we have rejected rows.
	defer reader.Flush()

	// Loop until we reached MaxRows (> 0) or when all rows has been
	// read (= -1)
	for {
		line, e := reader.ReadLine()

		if nil != e {
			if e != io.EOF {
				log.Print("dsv: ", e)
			}
			return n, e
		}

		// check for empty line
		line = bytes.TrimSpace(line)

		if len(line) <= 0 {
			continue
		}

		row, e := ParseLine(reader, &line)

		if nil == e {
			e = reader.PushRow(row)
		}
		if nil == e {
			n++
			if maxrows > 0 && n >= maxrows {
				break
			}
		} else {
			// If error, save the rejected line.
			log.Println(e)

			reader.Reject(line)
			reader.Reject([]byte("\n"))
		}
	}

	return n, e
}

/*
parsingLeftQuote parse the left-quote string from line.
*/
func parsingLeftQuote(md MetadataInterface, line *[]byte, p int) (int, error) {
	if "" == md.GetLeftQuote() {
		return p, nil
	}

	linelen := len(*line)
	lq := []byte(md.GetLeftQuote())

	if DEBUG {
		fmt.Println(md.GetLeftQuote())
	}

	for i := range lq {
		if p >= linelen {
			return p, &ErrReader{
				"Premature end-of-line",
				(*line),
			}
		}

		if DEBUG {
			fmt.Printf("%c:%c\n", (*line)[p], lq[i])
		}

		if (*line)[p] != lq[i] {
			return p, &ErrReader{
				"Invalid left-quote",
				(*line),
			}
		}
		p++
	}
	return p, nil
}

/*
parsingSeparator parsing the line until we found the separator.

Return the data and index of last parsed line, or error if separator is not
found or not match with specification.
*/
func parsingSeparator(md MetadataInterface, line *[]byte, p int) (
	v []byte, pRet int, e error,
) {
	if "" == md.GetSeparator() {
		v = append(v, (*line)[p:]...)
		return v, p, nil
	}

	linelen := len(*line)
	sep := []byte(md.GetSeparator())

	for p < linelen && (*line)[p] != sep[0] {
		v = append(v, (*line)[p])
		p++
	}

	if p >= linelen {
		return v, p, &ErrReader{
			"Missing separator, premature end-of-line",
			(*line),
		}
	}

	for i := range sep {
		if p >= linelen {
			return v, p, &ErrReader{
				"Missing separator, premature end-of-line",
				(*line),
			}
		}

		if (*line)[p] != sep[i] {
			return v, p, &ErrReader{
				"Invalid separator",
				(*line),
			}
		}
		p++
	}

	return v, p, nil
}

/*
parsingRightQuote parsing the line until we found the right quote.

Return the data and index of last parsed line, or error if right-quote is not
found or not match with specification.
*/
func parsingRightQuote(md MetadataInterface, line *[]byte, p int) (
	v []byte, pRet int, e error,
) {
	if "" == md.GetRightQuote() {
		return parsingSeparator(md, line, p)
	}

	linelen := len(*line)
	rq := []byte(md.GetRightQuote())

	// (2.2)
	for p < linelen && (*line)[p] != rq[0] {
		v = append(v, (*line)[p])
		p++
	}

	if p >= linelen {
		return v, p, &ErrReader{
			"Missing right-quote, premature end-of-line",
			(*line),
		}
	}

	// (2.2.1)
	for i := range rq {
		if p >= linelen {
			return v, p, &ErrReader{
				"Missing right-quote, premature end-of-line",
				(*line),
			}
		}

		if (*line)[p] != rq[i] {
			return v, p, &ErrReader{
				"Invalid right-quote",
				(*line),
			}
		}
		p++
	}

	// (2.2.2)
	if "" == md.GetSeparator() {
		return v, p, nil
	}

	// Skip all character until we found separator.
	sep := []byte(md.GetSeparator())

	for p < linelen && (*line)[p] != sep[0] {
		p++
	}

	if p >= linelen {
		return v, p, &ErrReader{
			"Missing separator, premature end-of-line",
			(*line),
		}
	}

	for i := range sep {
		if p >= linelen {
			return v, p, &ErrReader{
				"Missing separator, premature end-of-line",
				(*line),
			}
		}
		if (*line)[p] != sep[i] {
			return v, p, &ErrReader{
				"Invalid separator",
				(*line),
			}
		}
		p++
	}

	return v, p, nil
}

/*
ParseLine parse a line containing records. The output is array of record
(or single row).

This is how the algorithm works
(1) create n slice of record, where n is number of column metadata
(2) for each metadata
	(2.1) If using left quote, skip it
	(2.2) If using right quote, append byte to buffer until right-quote
		(2.2.1) Skip until the end of right quote
		(2.2.2) If using separator, skip until separator
	(2.3) If using separator, append byte to buffer until separator
	(2.4) else append all byte to buffer.
(3) save buffer to record
*/
func ParseLine(reader ReaderInterface, line *[]byte) (
	row Row, e error,
) {
	var md MetadataInterface
	var p = 0
	var rIdx = 0
	var inputMd []MetadataInterface

	inputMd = reader.GetInputMetadata()

	row = make(Row, reader.GetNColumn())

	for mdIdx := range inputMd {
		v := []byte{}
		md = inputMd[mdIdx]

		// skip all whitespace in the beginning
		for (*line)[p] == ' ' || (*line)[p] == '\t' {
			p++
		}

		// (2.1)
		p, e = parsingLeftQuote(md, line, p)

		if e != nil {
			return
		}

		v, p, e = parsingRightQuote(md, line, p)

		if e != nil {
			return
		}

		if DEBUG {
			fmt.Println(string(v))
		}

		if md.GetSkip() {
			continue
		}

		v = bytes.TrimSpace(v)
		r, e := NewRecord(v, md.GetType())

		if nil != e {
			return nil, &ErrReader{
				"Error or invalid type convertion",
				v,
			}
		}

		row[rIdx] = r
		rIdx++
	}

	return row, e
}
