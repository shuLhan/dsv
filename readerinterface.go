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
	GetInputMetadata () *[]Metadata
	GetInputMetadataAt (idx int) *Metadata
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

	Flush ()
	ReadLine () ([]byte, error)
	Reject (line []byte)
	Close ()
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

	// Check and initialize metadata.
	for i := range (*md) {
		e = (*md)[i].Init()

		// Count number of output columns.
		if ! (*md)[i].Skip {
			nColOut++
		}

		if nil != e {
			return e
		}
	}

	// Set number of output columns.
	reader.SetNColumn(nColOut)

	// Set default value
	reader.SetDefault()

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
func Read (reader ReaderInterface) (n int, e error) {
	n = 0
	reader.Reset ()

	// remember to flush if we have rejected rows.
	defer reader.Flush ()

	// Loop until we reached MaxRows (> 0) or when all rows has been
	// read (= -1)
	for {
		line, e := reader.ReadLine()

		if nil != e {
			if e != io.EOF {
				log.Print ("dsv: ", e)
			}
			reader.SetNRows(n)
			return n, e
		}

		// check for empty line
		line = bytes.TrimSpace (line)

		if len (line) <= 0 {
			continue
		}

		row, e := ParseLine(reader, &line)

		if nil == e {
			switch reader.GetMode() {
			case DatasetModeRows:
				reader.PushRow(row)
			case DatasetModeColumns:
				e = reader.PushRowToColumns(row)
			case DatasetModeMatrix:
				reader.PushRow(row)
				e = reader.PushRowToColumns(row)
			}
		}
		if nil == e {
			n++
			maxrows := reader.GetMaxRows()
			if maxrows > 0 && n >= maxrows {
				break
			}
		} else {
			// If error, save the rejected line.
			log.Println(e)

			reader.Reject (line)
			reader.Reject ([]byte ("\n"))
		}
	}

	reader.SetNRows(n)

	return n, e
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
func ParseLine (reader ReaderInterface, line *[]byte) (
					row Row, e error) {
	var md *Metadata
	var p = 0
	var l = len (*line)
	var rIdx = 0
	var inputMd *[]Metadata;

	inputMd = reader.GetInputMetadata ()

	row = make(Row, reader.GetNColumn())

	for mdIdx := range (*inputMd) {
		v := []byte{}
		md = &(*inputMd)[mdIdx]

		// skip all whitespace in the beginning
		for (*line)[p] == ' ' || (*line)[p] == '\t' {
			p++
		}

		// (2.1)
		if "" != md.LeftQuote {
			lq := []byte (md.LeftQuote)

			if DEBUG {
				fmt.Println (md.LeftQuote)
			}

			for i := range lq {
				if p >= l {
					return nil, &ErrReader {
						"Premature end-of-line",
						(*line),
					}
				}

				if DEBUG {
					fmt.Printf ("%c:%c\n", (*line)[p], lq[i])
				}

				if (*line)[p] != lq[i] {
					return nil, &ErrReader {
						"Invalid left-quote",
						(*line),
					}
				}
				p++
			}
		}

		if "" != md.RightQuote {
			rq := []byte (md.RightQuote)

			// (2.2)
			for p < l && (*line)[p] != rq[0] {
				v = append (v, (*line)[p])
				p++
			}

			if p >= l {
				return nil, &ErrReader {
					"Missing right-quote, premature end-of-line",
					(*line),
				}
			}

			// (2.2.1)
			for i := range rq {
				if p >= l {
					return nil, &ErrReader {
						"Missing right-quote, premature end-of-line",
						(*line),
					}
				}

				if (*line)[p] != rq[i] {
					return nil, &ErrReader {
						"Invalid right-quote",
						(*line),
					}
				}
				p++
			}

			// (2.2.2)
			if "" != md.Separator {
				sep := []byte (md.Separator)

				for p < l && (*line)[p] != sep[0] {
					p++
				}

				if p >= l {
					return nil, &ErrReader {
						"Missing separator, premature end-of-line",
						(*line),
					}
				}

				for i := range sep {
					if p >= l {
						return nil, &ErrReader {
							"Missing separator, premature end-of-line",
							(*line),
						}
					}
					if (*line)[p] != sep[i] {
						return nil, &ErrReader {
							"Invalid separator",
							(*line),
						}
					}
					p++
				}
			}
		} else if "" != md.Separator {
			// (2.3)
			sep := []byte (md.Separator)

			for p < l && (*line)[p] != sep[0] {
				v = append (v, (*line)[p])
				p++
			}

			if p >= l {
				return nil, &ErrReader {
					"Missing separator, premature end-of-line",
					(*line),
				}
			}

			for i := range sep {
				if p >= l {
					return nil, &ErrReader {
						"Missing separator, premature end-of-line",
						(*line),
					}
				}

				if (*line)[p] != sep[i] {
					return nil, &ErrReader {
						"Invalid separator",
						(*line),
					}
				}
				p++
			}
		} else {
			v = append (v, (*line)[p:]...)
		}

		if DEBUG {
			fmt.Println (string (v))
		}

		if md.Skip {
			continue
		}

		v = bytes.TrimSpace (v)
		r, e := NewRecord(v, md.T)

		if nil != e {
			return nil, &ErrReader {
				"Error or invalid type convertion",
				v,
			}
		}

		row[rIdx] = r
		rIdx++
	}

	return row, e
}
