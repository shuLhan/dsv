/*
Package dsv is a library for working with delimited separated value (DSV).

DSV is a free-style form of CSV format of text data, where each record is
separated by newline, and each field can be separated by any string.
*/
package dsv

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	// ErrNoInput define an error when no Input file is given in JSON.
	ErrNoInput = errors.New ("dsv: No input file is given")
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

const (
	// DefaultRejected define the default file which will contain the
	// rejected record when not defined in JSON config.
	DefaultRejected		= "rejected.dsv"
	// DefaultMaxRecord define default maximum record that will be saved
	// in memory when not defined in JSON config.
	DefaultMaxRecord	= 10
)

/*
Reader hold all configuration, metadata and input records.
*/
type Reader struct {
	// Input file, mandatory.
	Input		string		`json:"Input"`
	// Skip n lines from the head.
	Skip		int		`json:"Skip"`
	// Rejected is the file where record that does not fit
	// with metadata will be saved.
	Rejected	string		`json:"Rejected"`
	// FieldMetadata define format each field in a record.
	FieldMetadata	[]Metadata	`json:"FieldMetadata"`
	// MaxRecord define maximum record that this reader will read and
	// saved in the memory at one read operation.
	MaxRecord	int		`json:"MaxRecord"`
	// NRecord define number of record readed and saved in Records.
	NRecord		int		`json:"-"`
	// Records is input record that has been parsed.
	Records		*Row		`json:"-"`
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
		Rejected	:"rejected.dsv",
		FieldMetadata	:nil,
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
CloseReader will close all open descriptors.
*/
func (reader *Reader) CloseReader () {
	reader.bufReject.Flush ()
	reader.fReject.Close ()
	reader.fRead.Close ()
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
	for i := range reader.FieldMetadata {
		reader.FieldMetadata[i].SetDefault ()
	}
}

/*
push record to row.
*/
func (reader *Reader) push (r *[]Record) {
	var row = NewRow (r)

	if nil == reader.Records {
		reader.Records =  row
	} else {
		reader.Records.Last.Next = row
	}

	reader.Records.Last = row
}

/*
openReader open the input file, metadata must have been initialize.
*/
func (reader *Reader) openReader () error {
	fRead, e := os.OpenFile (reader.Input, os.O_RDONLY, 0600)
	if nil != e {
		return e
	}

	reader.bufRead = bufio.NewReader (fRead)

	return nil
}

/*
openRejected open rejected file, for saving unparseable line.
*/
func (reader *Reader) openRejected () error {
	fReject, e := os.OpenFile (reader.Rejected, os.O_CREATE | os.O_WRONLY,
					0600)
	if nil != e {
		return e
	}

	reader.bufReject = bufio.NewWriter (fReject)

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
ParseFieldMetadata from JSON string.
*/
func (reader *Reader) ParseFieldMetadata (md string) (e error) {
	e = json.Unmarshal ([]byte (md), reader)

	if nil != e {
		return
	}

	// Exit immediately if no input file is defined in config.
	if "" == reader.Input {
		return ErrNoInput
	}

	// Set default value for metadata.
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

	return nil
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
		(2.2.2) Skip until separator
	(2.3) else, append byte to buffer until separator
(3) save buffer to record
*/
func (reader *Reader) parseLine (line *[]byte) (records []Record, e error) {
	var md *Metadata
	var p = 0
	var l = len (*line)

	records = make ([]Record, len (reader.FieldMetadata))

	for mdIdx := range reader.FieldMetadata {
		v := []byte{}
		md = &reader.FieldMetadata[mdIdx]

		// (2.1)
		if "" != md.LeftQuote {
			lq := []byte (md.LeftQuote)

			for i := range lq {
				if (*line)[p] != lq[i] {
					return nil, &ErrReader {
						"Invalid left quote",
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
					"Invalid right quote",
					(*line),
				}
			}

			// (2.2.1)
			for i := range rq {
				if (*line)[p] != rq[i] {
					return nil, &ErrReader {
						"Invalid right quote",
						(*line),
					}
				}
				p++
			}

			// (2.2.2)
			sep := []byte (md.Separator)

			for p < l && (*line)[p] != sep[0] {
				p++
			}

			if p < l {
				for i := range sep {
					if (*line)[p] != sep[i] {
						return nil, &ErrReader {
							"Invalid separator",
							(*line),
						}
					}
					p++
				}
			}
		} else {
			// (2.3)
			sep := []byte (md.Separator)

			for p < l && (*line)[p] != sep[0] {
				v = append (v, (*line)[p])
				p++
			}

			if p >= l {
				return nil, &ErrReader {
					"Invalid line",
					(*line),
				}
			}

			for i := range sep {
				if (*line)[p] != sep[i] {
					return nil, &ErrReader {
						"Invalid separator",
						(*line),
					}
				}
				p++
			}
		}

		records[mdIdx] = v
	}

	return records, e
}

/*
Read maximum 'MaxRecord' record from file.
*/
func (reader *Reader) Read () (n int, e error) {
	var records []Record
	var line []byte

	defer reader.bufReject.Flush ()

	for n = 0; n < reader.MaxRecord; n++ {
		line, e = reader.readLine ()

		if nil != e {
			log.Print ("dsv: ", e)
			return n, e
		}

		records, e = reader.parseLine (&line)

		// If error, save the rejected line.
		if nil == e {
			reader.push (&records)
		} else {
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

	l,r := len (reader.FieldMetadata), len (other.FieldMetadata)

	if (l != r) {
		return false
	}

	for a := 0; a < l; a++ {
		if ! reader.FieldMetadata[a].IsEqual (&other.FieldMetadata[a]) {
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
