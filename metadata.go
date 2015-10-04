package dsv

import (
	"encoding/json"
	"log"
)

const (
	// DefaultSeparator for field.
	DefaultSeparator = ","
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
func (md *Metadata) SetDefault () {
	if "" == md.Separator {
		md.Separator = DefaultSeparator
	}
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
