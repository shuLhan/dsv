// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dsv

/*
SortRecordsByIndex sort record in column by index.
*/
func SortRecordsByIndex(data []*Record, sortedIdx []int) (sorted []*Record) {
	sorted = make([]*Record, len(data))

	for i := range sortedIdx {
		sorted[i] = data[sortedIdx[i]]
	}
	return
}
