/*
Package dsv is a library for working with delimited separated value (DSV).

DSV is a free-style form of CSV format of text data, where each record is
separated by newline, and each field can be separated by any string.
*/
package dsv

import (
	"errors"
	"io/ioutil"
	"encoding/json"
	"log"
	"os"
)

const (
	// DefaultRejected define the default file which will contain the
	// rejected record when not defined in JSON config.
	DefaultRejected		= "rejected.dsv"
	// DefaultMaxRecord define default maximum record that will be saved
	// in memory when not defined in JSON config.
	DefaultMaxRecord	= 256
)

var (
	// ErrNoInput define an error when no Input file is given to Reader.
	ErrNoInput	= errors.New ("dsv: No input file is given in config")
	// ErrNoOutput define an error when no output file is given to Writer.
	ErrNoOutput	= errors.New ("dsv: No output file is given in config")
	// ErrNotOpen define an error when output file has not been opened
	// by Writer.
	ErrNotOpen	= errors.New ("dsv: Output file is not opened")
	// ErrNilReader define an error when Reader object is nil when passed
	// to Write function.
	ErrNilReader	= errors.New ("dsv: Reader object is nil")

	// DEBUG exported from environment to debug the library.
	DEBUG		= bool (os.Getenv ("DEBUG") != "")
)

/*
DSV combine reader and writer.
*/
type DSV struct {
	Reader
	Writer
}

/*
New create a new DSV object.
*/
func New () *DSV {
	return &DSV {}
}
/*
Open configuration file for reading and writing.
*/
func (dsv *DSV) Open (fcfg string) (e error) {
	cfg, e := ioutil.ReadFile (fcfg)
	if nil != e {
		log.Print ("dsv: ", e)
		return
	}

	e = json.Unmarshal ([]byte (cfg), dsv)

	if nil != e {
		return
	}

	e = dsv.Reader.Init ()

	if nil != e {
		return
	}

	return dsv.Writer.Init ()
}

/*
Close reader and writer.
*/
func (dsv *DSV) Close () {
	dsv.Writer.Close ()
	dsv.Reader.Close ()
}
