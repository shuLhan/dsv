// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package assert provided common functions for testing.
*/
package assert

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"runtime/debug"
	"testing"
)

/*
Assert print fatal message when our expectation `exp` is different with test
result `got`.
*/
func Equal(t *testing.T, exp, got interface{}) {
	if !reflect.DeepEqual(exp, got) {
		debug.PrintStack()
		t.Fatal("Expecting", exp, "got", got)
	}
}

/*
EqualFileContent compare content of two file, print error message and exit
when both are different.
*/
func EqualFileContent(t *testing.T, a, b string) {
	out, e := ioutil.ReadFile(a)

	if nil != e {
		debug.PrintStack()
		t.Error(e)
	}

	exp, e := ioutil.ReadFile(b)

	if nil != e {
		debug.PrintStack()
		t.Error(e)
	}

	r := bytes.Compare(out, exp)

	if 0 != r {
		debug.PrintStack()
		t.Error("Comparing", a," with ", b,": result different (", r ,")")
		t.Error("a:\n", out)
		t.Error("b:\n", exp)
	}
}
