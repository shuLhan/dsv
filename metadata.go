/*
Copyright 2015 Mhd Sulhan <ms@kilabit.info>
All rights reserved.  Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/
package dsv

import (
	"encoding/json"
	"log"
)

const (
	DefaultSeparator = "," // Default separator for field.
)

/*
Metadata represent on how to parse each column in record.
*/
type Metadata struct {
	// Name of the field, optional.
	Name		string	`json:"Name"`
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
SetDefault value in instance.
*/
func (this *Metadata) SetDefault () {
	if "" == this.Separator {
		this.Separator = DefaultSeparator
	}
}

/*
IsEqual return true if this metadata equal with other instance, return false
otherwise.
*/
func (this *Metadata) IsEqual (o *Metadata) bool {
	if this == o {
		return true
	}
	if this.Name != o.Name {
		return false
	}
	if this.Separator != o.Separator {
		return false
	}
	if this.LeftQuote != o.LeftQuote {
		return false
	}
	if this.RightQuote != o.RightQuote {
		return false
	}
	return true
}

/*
String yes, it will print it JSON like format.
*/
func (this *Metadata) String() string {
	r, e := json.MarshalIndent (this, "", "\t")
	if nil != e {
		log.Print (e)
	}
	return string (r)
}
