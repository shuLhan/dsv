package dsv

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
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
		OutputMetadata	:nil,
		fWriter		:nil,
		BufWriter	:nil,
	}
}

/*
Open file for writing.
*/
func (writer *Writer) Open (fcfg string) (e error) {
	cfg, e := ioutil.ReadFile (fcfg)

	if nil != e {
		log.Print ("dsv: ", e)
		return e
	}

	e = writer.ParseConfig (cfg)

	return nil
}

/*
Init initialize writer by opening output file.
*/
func (writer *Writer) Init () error {
	return writer.openOutput ()
}

/*
ParseConfig from JSON string.
*/
func (writer *Writer) ParseConfig (cfg []byte) (e error) {
	e = json.Unmarshal ([]byte (cfg), writer)

	if nil != e {
		return
	}

	if "" == writer.Output {
		return ErrNoOutput
	}

	return writer.Init ()
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
func (writer *Writer) WriteRecords (records *[]Record) (e error) {
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
Write records from Reader to file.
Return n for number of records written, and e for error that happened when
writing to file.
*/
func (writer *Writer) Write (reader *Reader) (n int, e error) {
	if nil == reader {
		return 0, ErrNilReader
	}
	if nil == writer.fWriter {
		return 0, ErrNotOpen
	}

	n = 0
	row := reader.Records.Front ()

	for nil != row {
		e = writer.WriteRecords (row.Value.(*[]Record))
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
String yes, it will print it in JSON like format.
*/
func (writer *Writer) String() string {
	r, e := json.MarshalIndent (writer, "", "\t")

	if nil != e {
		log.Print (e)
	}

	return string (r)
}
