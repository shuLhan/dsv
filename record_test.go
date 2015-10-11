package dsv_test

import (
	"fmt"
	"testing"

	"github.com/shuLhan/dsv"
)

/*
TestRecord simply check how the stringer work.
*/
func TestRecord (t *testing.T) {
	var exp = []string {
		"test",
		"1",
		"2",
	}

	var e error

	r := make (dsv.RecordSlice, len (exp))

	r[0], e = dsv.RecordNew ([]byte ("test"), dsv.TString)
	if nil != e {
		t.Error (e)
	}

	r[1], e = dsv.RecordNew ([]byte ("1"), dsv.TInteger)
	if nil != e {
		t.Error (e)
	}

	r[2], e = dsv.RecordNew ([]byte ("02"), dsv.TInteger)
	if nil != e {
		t.Error (e)
	}

	for i := range exp {
		s := fmt.Sprint (r[i])

		if s != exp[i] {
			t.Error ("dsv_test: expecting\n", exp[i], "\n got\n",
				r[i])
		}
	}
}
