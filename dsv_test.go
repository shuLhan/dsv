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
doReadWrite test reading and writing the DSV data.
*/
func doReadWriteDSV (rw *dsv.ReadWriter, t *testing.T) {
	i	:= 0
	n 	:= 0
	e	:= error (nil)

	for {
		n, e = rw.Read ()

		if n > 0 {
			r := fmt.Sprint (rw.Records)

			if r != expectation[i] {
				t.Error ("dsv_test: expecting\n",
					expectation[i],
					" got\n", r)
			}

			rw.Write (&rw.Reader)

			i++
		} else if e == io.EOF {
			// EOF
			break
		}
	}
}

/*
TestReadWriter test reading and writing DSV.
*/
func TestReadWriter (t *testing.T) {
	if DEBUG {
		log.Println (">>> TestReadWriter")
	}

	// Initialize dsv
	rw := dsv.New ()

	e := rw.Open ("config.dsv")

	if nil != e {
		t.Error (e)
	}

	doReadWriteDSV (rw, t)

	rw.Close ()

	// Compare the ouput from Writer
	out, e := ioutil.ReadFile (rw.Output)

	if nil != e {
		t.Error (e)
	}

	exp, e := ioutil.ReadFile ("expected.dsv")

	if nil != e {
		t.Error (e)
	}

	r := bytes.Compare (out, exp)

	if 0 != r {
		t.Error ("Output different from expected (", r ,")")
	}
}
