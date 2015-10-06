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
func doReadWriteDSV (dsv *dsv.ReadWriter, t *testing.T) {
	exp	:= ""
	i	:= 0
	n 	:= 0
	e	:= error (nil)

	for {
		n, e = dsv.Read ()

		if n > 0 {
			r := fmt.Sprint (dsv.Records)

			if r != expectation[i] {
				t.Error ("dsv_test: expecting\n", exp,
					" got\n", r)
			}

			dsv.Write (&dsv.Reader)

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
	dsv := dsv.New ()

	e := dsv.Open ("config.dsv")

	if nil != e {
		t.Error (e)
	}

	doReadWriteDSV (dsv, t)

	dsv.Close ()

	// Compare the ouput from Writer
	out, e := ioutil.ReadFile (dsv.Output)

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
