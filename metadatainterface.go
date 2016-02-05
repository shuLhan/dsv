// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
MetadataInterface is the interface for field metadata.
This is to make anyone can extend the DSV library including the metadata.
*/
type MetadataInterface interface {
	Init()
	GetName() string
	GetType() int
	GetLeftQuote() string
	GetRightQuote() string
	GetSeparator() string
	GetSkip() bool
	GetValueSpace() []string

	IsEqual(MetadataInterface) bool
}
