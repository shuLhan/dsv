// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

import (
	"encoding/json"
	"log"
	"strings"
)

/*
Metadata represent on how to parse each column in record.
*/
type Metadata struct {
	// Name of the field, optional.
	Name		string	`json:"Name"`
	// Type of the field, default to "string".
	// Valid value are: "string", "integer", "real"
	Type		string	`json:"Type"`
	T		int
	// Separator for field in record.
	Separator	string	`json:"Separator"`
	// LeftQuote define the characters that enclosed the field in the left
	// side.
	LeftQuote	string	`json:"LeftQuote"`
	// RightQuote define the characters that enclosed the field in the
	// right side.
	RightQuote	string	`json:"RightQuote"`
}

/*
Init initalize metadata field, i.e. check and set field type.
*/
func (md *Metadata) Init () (e error) {
	switch strings.ToUpper (md.Type) {
	case "STRING":
		md.T = TString
	case "INTEGER", "INT":
		md.T = TInteger
	case "REAL":
		md.T = TReal
	case "":
		md.T = TString
	default:
		e = &ErrReader {
			"dsv: Invalid type",
			[]byte (md.Type),
		}
	}

	return
}

/*
IsEqual return true if this metadata equal with other instance, return false
otherwise.
*/
func (md *Metadata) IsEqual (o *Metadata) bool {
	if md == o {
		return true
	}
	if md.Name != o.Name {
		return false
	}
	if md.Separator != o.Separator {
		return false
	}
	if md.LeftQuote != o.LeftQuote {
		return false
	}
	if md.RightQuote != o.RightQuote {
		return false
	}
	return true
}

/*
String yes, it will print it JSON like format.
*/
func (md *Metadata) String() string {
	r, e := json.MarshalIndent (md, "", "\t")
	if nil != e {
		log.Print (e)
	}
	return string (r)
}
