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
	// Name of the column, optional.
	Name string `json:"Name"`
	// Type of the column, default to "string".
	// Valid value are: "string", "integer", "real"
	Type string `json:"Type"`
	// T type of column in integer.
	T int
	// Separator for column in record.
	Separator string `json:"Separator"`
	// LeftQuote define the characters that enclosed the column in the left
	// side.
	LeftQuote string `json:"LeftQuote"`
	// RightQuote define the characters that enclosed the column in the
	// right side.
	RightQuote string `json:"RightQuote"`
	// Skip, if its true this column will be ignored, not saved in reader
	// object. Default to false.
	Skip bool `json:"Skip"`
}

/*
Init initalize metadata column, i.e. check and set column type.
*/
func (md *Metadata) Init() (e error) {
	switch strings.ToUpper(md.Type) {
	case "STRING":
		md.T = TString
	case "INTEGER", "INT":
		md.T = TInteger
	case "REAL":
		md.T = TReal
	case "":
		md.T = TString
	default:
		e = &ErrReader{
			"dsv: Invalid type",
			[]byte(md.Type),
		}
	}

	return
}

/*
GetName return the name of metadata.
*/
func (md *Metadata) GetName() string {
	return md.Name
}

/*
GetType return type of metadata.
*/
func (md *Metadata) GetType() int {
	return md.T
}

/*
GetSeparator return the field separator.
*/
func (md *Metadata) GetSeparator() string {
	return md.Separator
}

/*
GetLeftQuote return the string used in the beginning of record value.
*/
func (md *Metadata) GetLeftQuote() string {
	return md.LeftQuote
}

/*
GetRightQuote return string that end in record value.
*/
func (md *Metadata) GetRightQuote() string {
	return md.RightQuote
}

/*
GetSkip return number of rows that will be skipped when reading data.
*/
func (md *Metadata) GetSkip() bool {
	return md.Skip
}

/*
IsEqual return true if this metadata equal with other instance, return false
otherwise.
*/
func (md *Metadata) IsEqual(o MetadataInterface) bool {
	if md.Name != o.GetName() {
		return false
	}
	if md.Separator != o.GetSeparator() {
		return false
	}
	if md.LeftQuote != o.GetLeftQuote() {
		return false
	}
	if md.RightQuote != o.GetRightQuote() {
		return false
	}
	return true
}

/*
String yes, it will print it JSON like format.
*/
func (md *Metadata) String() string {
	r, e := json.MarshalIndent(md, "", "\t")
	if nil != e {
		log.Print(e)
	}
	return string(r)
}
