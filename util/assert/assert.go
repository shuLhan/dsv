// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package assert provided common functions for testing.
*/
package assert

import (
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
