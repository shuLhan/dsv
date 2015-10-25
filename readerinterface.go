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
	GetInputMetadata () *[]Metadata
	GetInputMetadataAt (idx int) *Metadata
	GetMaxRecord () int
	SetMaxRecord (max int)
	GetRecordRead () int
	SetRecordRead (n int)

	Reset ()
	Flush ()
	ReadLine () (line []byte, e error)
	Push (r *RecordSlice)
	Reject (line []byte)
	Close ()
}

/*
Read records from input file.
*/
func Read (reader ReaderInterface) (n int, e error) {
	var records *RecordSlice
	var line []byte

	reader.Reset ()

	// remember to flush if we have rejected record.
	defer reader.Flush ()

	// Loop until we reached MaxRecord (> 0) or when all record has been
	// read (= -1)
	for {
		line, e = reader.ReadLine ()

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

		records, e = ParseLine (reader, &line)

		if nil == e {
			reader.Push (records)
			n++

			if reader.GetMaxRecord () > 0 &&
			n >= reader.GetMaxRecord () {
				break
			}
		} else {
			// If error, save the rejected line.
			if DEBUG {
				fmt.Println (e)
			}

			reader.Reject (line)
			reader.Reject ([]byte ("\n"))
		}
	}

	reader.SetRecordRead (n)

	return n, e
}

/*
ParseLine parse a line containing record. The output is array of fields added
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
func ParseLine (reader ReaderInterface, line *[]byte) (
					precords *RecordSlice, e error) {
	var md *Metadata
	var p = 0
	var l = len (*line)
	var inputMd *[]Metadata;

	inputMd = reader.GetInputMetadata ()

	records := make (RecordSlice, len ((*inputMd)))

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
