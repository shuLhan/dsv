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
	}

	r := &dsv.Record {'t','e','s','t'}

	s := fmt.Sprint (r)

	if s != exp[0] {
		t.Error ("dsv_test: expecting\n", exp[0], "\n got\n", r)
	}
}
