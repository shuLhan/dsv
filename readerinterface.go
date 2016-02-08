// Copyright 2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

/*
ReaderInterface is the interface for reading DSV file.
*/
type ReaderInterface interface {
	ConfigInterface
	DatasetInterface
	AddInputMetadata(*Metadata)
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

	Flush() error
	ReadLine() ([]byte, error)
	Reject(line []byte) (int, error)
	Close() error
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

	// Set default value
	reader.SetDefault()

	// Check if output mode is valid and initialize it if valid.
	e = reader.SetDatasetMode(reader.GetDatasetMode())
	if nil != e {
		return
	}

	// Check and initialize metadata and columns attributes.
	md := reader.GetInputMetadata()
	for i := range md {
		md[i].Init()

		if nil != e {
			return e
		}

		// Count number of output columns.
		if !md[i].GetSkip() {
			// add type of metadata to list of type
			col := Column{
				Type:       md[i].GetType(),
				Name:       md[i].GetName(),
				ValueSpace: md[i].GetValueSpace(),
			}
			reader.PushColumn(col)
		}
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
	linenum := 0
	maxrows := reader.GetMaxRows()

	e = reader.Reset()
	if e != nil {
		return
	}

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

		row, errReader := ParseLine(reader, &line)

		if nil == errReader {
			reader.PushRow(row)

			n++
			if maxrows > 0 && n >= maxrows {
				break
			}
		} else {
			errReader.N = linenum
			fmt.Fprintf(os.Stderr, "%s\n", errReader)

			// If error, save the rejected line.
			line = append(line, "\n"...)

			_, e = reader.Reject(line)
			if e != nil {
				break
			}
		}
		linenum++
	}

	// remember to flush if we have rejected rows.
	e = reader.Flush()

	return n, e
}

/*
parsingLeftQuote parse the left-quote string from line.
*/
func parsingLeftQuote(md MetadataInterface, line *[]byte, p int) (
	int, *ReaderError,
) {
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
			goto Err
		}

		if DEBUG {
			fmt.Printf("%c:%c\n", (*line)[p], lq[i])
		}

		if (*line)[p] != lq[i] {
			goto Err
		}
		p++
	}
	return p, nil
Err:
	return p, &ReaderError{
		"parsingLeftQuote",
		"Missing left-quote '" + string(lq) + "'",
		string(*line), p, 0,
	}
}

/*
isForwardMatch return true if current line at index p match with token,
otherwise return false.
*/
func isForwardMatch(line *[]byte, p int, token []byte) (bool, *ReaderError) {
	linelen := len(*line)
	tokenlen := len(token)

	if p+tokenlen > linelen {
		return false, &ReaderError{
			"isForwardMatch",
			"Missing token '" + string(token) + "'",
			string(*line), p, 0,
		}
	}

	for _, v := range token {
		if v != (*line)[p] {
			return false, nil
		}
		p++
	}
	return true, nil
}

/*
parsingUntil we found token. Token that is prefixed with escaped character
'\' will be ignored.
*/
func parsingUntil(line *[]byte, p int, token []byte) (
	v []byte, pRet int, e *ReaderError,
) {
	linelen := len(*line)

	escaped := false
	for p < linelen {
		// Assume the escape character always used to escaped the
		// token ...
		if (*line)[p] == '\\' {
			escaped = true
			p++
			continue
		}
		if (*line)[p] != token[0] {
			if escaped {
				// ... turn out its not escaping token.
				v = append(v, '\\')
				escaped = false
			}

			v = append(v, (*line)[p])
			p++
			continue
		}

		// We found the first token character.
		// Lets check if its match with all content of token.
		match, e := isForwardMatch(line, p, token)

		if e != nil {
			return v, p, e
		}

		// false alarm ...
		if !match {
			if escaped {
				v = append(v, '\\')
				escaped = false
			}

			v = append(v, (*line)[p])
			p++
			continue
		}

		// Its matched, but if its prefixed with escaped char '\', then
		// we assumed it as non breaking token.
		if escaped {
			v = append(v, (*line)[p])
			p++
			escaped = false
			continue
		}

		// Its matched with no escape character.
		break
	}

	if p >= linelen {
		return v, p, &ReaderError{
			"parsingUntil",
			"Missing token '" + string(token) + "'",
			string(*line), p, 0,
		}
	}

	return v, p + len(token), e
}

/*
skipUntil skip all characters until matched token is found.
Return index of line with matched token or error if line end before finding
the token.
*/
func skipUntil(line *[]byte, p int, token []byte) (
	pRet int, e *ReaderError,
) {
	linelen := len(*line)

	for p < linelen {
		if (*line)[p] != token[0] {
			p++
			continue
		}

		// We found the first token character.
		// Lets check if its match with all content of token.
		match, e := isForwardMatch(line, p, token)

		if e != nil {
			return p, e
		}

		// false alarm ...
		if !match {
			p++
			continue
		}

		// Its matched.
		break
	}

	if p >= linelen {
		return p, &ReaderError{
			"skipUntil",
			"Missing token '" + string(token) + "'",
			string(*line), p, 0,
		}
	}

	return p + len(token), e
}

/*
parsingSeparator parsing the line until we found the separator.

Return the data and index of last parsed line, or error if separator is not
found or not match with specification.
*/
func parsingSeparator(md MetadataInterface, line *[]byte, p int) (
	v []byte, pRet int, e *ReaderError,
) {
	if "" == md.GetSeparator() {
		v = append(v, (*line)[p:]...)
		return v, p, nil
	}

	sep := []byte(md.GetSeparator())

	v, p, e = parsingUntil(line, p, sep)

	if e != nil {
		e.Func = "parsingSeparator"
	}

	return v, p, e
}

/*
parsingRightQuote parsing the line until we found the right quote.

Return the data and index of last parsed line, or error if right-quote is not
found or not match with specification.
*/
func parsingRightQuote(md MetadataInterface, line *[]byte, p int) (
	v []byte, pRet int, e *ReaderError,
) {
	if "" == md.GetRightQuote() {
		return parsingSeparator(md, line, p)
	}

	rq := []byte(md.GetRightQuote())

	// (2.2.1)
	v, p, e = parsingUntil(line, p, rq)

	if e != nil {
		e.Func = "parsingRightQuote"
		return v, p, e
	}

	if "" == md.GetSeparator() {
		return v, p, nil
	}

	// (2.2.2)
	// Skip all character until we found separator.
	sep := []byte(md.GetSeparator())

	p, e = skipUntil(line, p, sep)

	if e != nil {
		e.Func = "parsingRightQuote"
	}

	return v, p, e
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
	row Row, e *ReaderError,
) {
	var md MetadataInterface
	var p = 0
	var rIdx = 0
	var inputMd []MetadataInterface
	linelen := len(*line)

	inputMd = reader.GetInputMetadata()

	row = make(Row, reader.GetNColumn())

	for mdIdx := range inputMd {
		v := []byte{}
		md = inputMd[mdIdx]

		// skip all whitespace in the beginning
		for p < linelen && ((*line)[p] == ' ' || (*line)[p] == '\t') {
			p++
		}

		// (2.1)
		p, e = parsingLeftQuote(md, line, p)

		if e != nil {
			return
		}

		// (2.2)
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
		r, e := NewRecord(string(v), md.GetType())

		if nil != e {
			return nil, &ReaderError{
				"ParseLine",
				"Type convertion error '" + string(v) + "'",
				string(*line), p, 0,
			}
		}

		row[rIdx] = r
		rIdx++
	}

	return row, e
}
