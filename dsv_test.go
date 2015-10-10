package dsv_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"testing"

	"github.com/shuLhan/dsv"
)

/*
doInit create read-write object.
*/
func doInit (testName string, fcfg string, t *testing.T) (rw *dsv.ReadWriter, e error) {
	if DEBUG {
		log.Println (">>> ", testName)
	}

	// Initialize dsv
	rw = dsv.New ()

	e = rw.Open (fcfg)

	if nil != e {
		t.Fatal (e)
	}

	return
}

/*
doReadWriteDSV test reading and writing the DSV data.
*/
func doReadWriteDSV (rw *dsv.ReadWriter, t *testing.T, check bool) {
	i	:= 0
	n 	:= 0
	e	:= error (nil)

	for {
		n, e = rw.Read ()

		if DEBUG {
			log.Println ("n records: ", n)
		}

		if n > 0 {
			if check {
				r := fmt.Sprint (rw.Records)

				if r != expectation[i] {
					t.Fatal ("dsv_test: expecting\n",
						expectation[i],
						" got\n", r)
					break
				}
				i++
			}

			rw.Write (&rw.Reader)
		} else if e == io.EOF {
			// EOF
			break
		}
	}
}

func doCompare (fout *string, t *testing.T) {
	// Compare the ouput from Writer
	out, e := ioutil.ReadFile (*fout)

	if nil != e {
		t.Fatal (e)
	}

	exp, e := ioutil.ReadFile ("expected.dsv")

	if nil != e {
		t.Fatal (e)
	}

	r := bytes.Compare (out, exp)

	if 0 != r {
		t.Fatal ("Output different from expected (", r ,")")
	}
}

/*
TestReadWriter test reading and writing DSV.
*/
func TestReadWriter (t *testing.T) {
	rw, _ := doInit ("TestReadWriter", "config.dsv", t)

	doReadWriteDSV (rw, t, true)

	rw.Close ()

	doCompare (&rw.Output, t)
}

/*
TestReadWriter test reading and writing DSV.
*/
func TestReadWriterAll (t *testing.T) {
	rw, _ := doInit ("TestReadWriterAll", "config.dsv", t)

	rw.MaxRecord = -1;

	doReadWriteDSV (rw, t, false)

	rw.Close ()

	doCompare (&rw.Output, t)
}
