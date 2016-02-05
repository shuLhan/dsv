// Copyright 2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"fmt"
)

/*
readerError to handle error data and message.
*/
type readerError struct {
	// Func where error happened
	Func string
	// What cause the error?
	What string
	// Line define the line which cause error
	Line string
	// Pos character position which cause error
	Pos int
	// N line number
	N int
}

/*
Error to string.
*/
func (e *readerError) Error() string {
	return fmt.Sprintf("dsv.Reader.%s [%d:%d]: %s |%s|", e.Func, e.N,
		e.Pos, e.What, e.Line)
}
