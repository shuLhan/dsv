// Copyright 2015-2016 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"io"
)

/*
SimpleRead provide a shortcut to read data from file using configuration file
from `fcfg`.
Return the reader contained data or error if failed.
Reader object upon returned has been closed, so if one need to read all
data in it simply set the `MaxRows` to `-1` in config file.
*/
func SimpleRead(fcfg string) (reader ReaderInterface, e error) {
	reader, e = NewReader(fcfg)

	if e != nil {
		return
	}

	_, e = Read(reader)
	if e != nil && e != io.EOF {
		return nil, e
	}

	e = reader.Close()

	return
}

/*
SimpleWrite provide a shortcut to write data from reader using output metadata
format and output file defined in file `fcfg`.
*/
func SimpleWrite(reader ReaderInterface, fcfg string) (e error) {
	writer, e := NewWriter(fcfg)
	if e != nil {
		return
	}

	_, e = writer.Write(reader)
	if e != nil {
		return
	}

	e = writer.Close()

	return
}
