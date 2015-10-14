// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	e := dsvReader.Open (configfile)

(2.1) Do not forget to check for error ...

(3) Make sure to close all files after finished

	defer dsvReader.Close ()

(4) Create loop to read and process records.

	for {
		n, e := dsvReader.Read ()

		if e == io.EOF {
			break
		}

(4.1) Iterate through records

		row := dsvReader.Records.Front ()
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
	// Skip n lines from the head.
	Skip		int		`json:"Skip"`
	// Rejected is the file where record that does not fit
	// with metadata will be saved.
	Rejected	string		`json:"Rejected"`
	// InputMetadata define format each field in a record.
	InputMetadata	[]Metadata	`json:"InputMetadata"`
	// MaxRecord define maximum record that this reader will read and
	// saved in the memory at one read operation.
	// If the value is -1 all records will read.
	MaxRecord	int		`json:"MaxRecord"`
	// NRecord define number of record readed and saved in Records.
	NRecord		int
	// Records is input record that has been parsed.
	Records		*Row
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
		Skip		:0,
		Rejected	:"rejected.dat",
		InputMetadata	:nil,
		MaxRecord	:DefaultMaxRecord,
		NRecord		:0,
		Records		:nil,
		fRead		:nil,
		fReject		:nil,
		bufRead		:nil,
		bufReject	:nil,
	}
}

/*
Open configuration file.
*/
func (reader *Reader) Open (fcfg string) error {
	cfg, e := ioutil.ReadFile (fcfg)
	if nil != e {
		log.Print ("dsv: ", e)
		return e
	}

	e = reader.ParseConfig (cfg)

	return e
}

/*
Close will close all open descriptors.
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
setDefault options for global config and each metadata.
*/
func (reader *Reader) setDefault () {
	if "" == reader.Rejected {
		reader.Rejected = DefaultRejected
	}
	if 0 == reader.MaxRecord {
		reader.MaxRecord = DefaultMaxRecord
	}
}

/*
push record to row.
*/
func (reader *Reader) push (r *RecordSlice) {
	if nil == reader.Records {
		reader.Records = &Row {}
	}

	reader.Records.PushBack (r)
}

/*
openReader open the input file, metadata must have been initialize.
*/
func (reader *Reader) openReader () (e error) {
	reader.fRead, e = os.OpenFile (reader.Input, os.O_RDONLY, 0600)
	if nil != e {
		return e
	}

	reader.bufRead = bufio.NewReader (reader.fRead)

	return nil
}

/*
openRejected open rejected file, for saving unparseable line.
*/
func (reader *Reader) openRejected () (e error) {
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
skipLines skip parsing n lines from input file.
The n is defined in the attribute "Skip"
*/
func (reader *Reader) skipLines () (e error) {
	for i := 0; i < reader.Skip; i++ {
		_, e = reader.readLine ()

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
	// Check and initialize metadata.
	for i := range reader.InputMetadata {
		e = reader.InputMetadata[i].Init ()
		if nil != e {
			return e
		}
	}

	// Set default value
	reader.setDefault ()

	// Get ready ...
	e = reader.openReader ()
	if nil != e {
		return
	}

	e = reader.openRejected ()
	if nil != e {
		return
	}

	// Skip lines
	if reader.Skip > 0 {
		e = reader.skipLines ()
		if nil != e {
			return
		}
	}

	return
}

/*
ParseConfig from JSON string.
*/
func (reader *Reader) ParseConfig (cfg []byte) (e error) {
	e = json.Unmarshal ([]byte (cfg), reader)

	if nil != e {
		return
	}

	// Exit immediately if no input file is defined in config.
	if "" == reader.Input {
		return ErrNoInput
	}

	return reader.Init ()
}

/*
readLine will read one line from input file.
*/
func (reader *Reader) readLine () (line []byte, e error) {
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
parseLine parse a line containing record. The output is array of fields added
to list of Reader's Records.

This is how the algorithm works
(1) create n slice of records, where n is number of field metadata
(2) for each metadata
	(2.1) If using left quote, skip it
	(2.2) If using right quote, append byte to buffer until right-quote
		(2.2.1) Skip until the end of right quote
		(2.2.2) If using separator, skip until separator
	(2.3) If using separator, append byte to buffer until separator
	(2.4) else append all byte to buffer.
(3) save buffer to record
*/
func (reader *Reader) parseLine (line *[]byte) (precords *RecordSlice, e error) {
	var md *Metadata
	var p = 0
	var l = len (*line)

	records := make (RecordSlice, len (reader.InputMetadata))

	for mdIdx := range reader.InputMetadata {
		v := []byte{}
		md = &reader.InputMetadata[mdIdx]

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

		v = bytes.TrimSpace (v)

		e = records[mdIdx].SetByte (v, md.T)

		if nil != e {
			return nil, &ErrReader {
				"Error or invalid type convertion",
				v,
			}
		}
	}

	return &records, e
}

/*
Reset all variables for next read operation. NRecord will be 0, and Records
will be nil again.
*/
func (reader *Reader) Reset () {
	reader.NRecord = 0
	reader.Records = nil
}

/*
Read maximum 'MaxRecord' record from file.
*/
func (reader *Reader) Read () (n int, e error) {
	var records *RecordSlice
	var line []byte

	reader.Reset ()

	// remember to flush if we have rejected record.
	defer reader.bufReject.Flush ()

	// Loop until we reached MaxRecord (> 0) or when all record has been
	// read (= -1)
	for {
		line, e = reader.readLine ()

		if nil != e {
			if DEBUG && e != io.EOF {
				log.Print ("dsv: ", e)
			}
			return n, e
		}

		// check for empty line
		line = bytes.TrimSpace (line)
		if len (line) <= 0 {
			continue
		}

		records, e = reader.parseLine (&line)

		if nil == e {
			reader.push (records)
			n++

			if reader.MaxRecord > 0 && n >= reader.MaxRecord {
				break
			}
		} else {
			// If error, save the rejected line.
			if DEBUG {
				fmt.Println (e)
			}
			reader.bufReject.Write (line)
			reader.bufReject.WriteString ("\n")
		}
	}

	reader.NRecord = n

	return n, e
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
