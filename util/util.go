// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package util contain common function to work with data.
*/
package util

import (
	"math/rand"
	"time"
)

const (
	// SortThreshold for sort. When the data less than SortThreshold,
	// insertion sort will be used to replace the sort.
	SortThreshold = 7
)

/*
InsertionSortFloat64 will sort the data using insertion algorithm.
*/
func InsertionSortFloat64(data []float64, idx []int, l, r int) {
	for x := l; x < r; x++ {
		for y := x + 1; y < r; y++ {
			if data[x] > data[y] {
				SwapInt(idx, x, y)
				SwapFloat64(data, x, y)
			}
		}
	}
}

/*
SwapInt swap two indices value of integer.
*/
func SwapInt(data []int, i, j int) {
	if i == j {
		return
	}

	tmp := data[i]
	data[i] = data[j]
	data[j] = tmp
}

/*
SwapFloat64 swap two indices value of 64bit float.
*/
func SwapFloat64(data []float64, i, j int) {
	if i == j {
		return
	}

	tmp := data[i]
	data[i] = data[j]
	data[j] = tmp
}

/*
SwapString swap two indices value of string.
*/
func SwapString(data []string, i, j int) {
	if i == j {
		return
	}

	tmp := data[i]
	data[i] = data[j]
	data[j] = tmp
}

/*
MergesortFloat64 sort the slice of float from `l` to `r` using mergesort
algorithm, return the sorted index.
*/
func MergesortFloat64(data []float64, sortedIdx []int, l, r int) {
	if l+SortThreshold >= r {
		InsertionSortFloat64(data, sortedIdx, l, r)
		return
	}

	res := (r + l) % 2
	c := (r + l) / 2
	if res == 1 {
		c++
	}

	MergesortFloat64(data, sortedIdx, l, c)
	MergesortFloat64(data, sortedIdx, c, r)

	// merging
	if data[c-1] < data[c] {
		// the last element of the left is lower then the first element
		// of the right, i.e. [1 2] [3 4].
		return
	}

	datalen := r - l
	newdata := make([]float64, datalen)
	newidx := make([]int, datalen)

	x := l
	y := c
	z := 0
	for ; x < c && y < r; z++ {
		if data[x] <= data[y] {
			newdata[z] = data[x]
			newidx[z] = sortedIdx[x]
			x++
		} else {
			newdata[z] = data[y]
			newidx[z] = sortedIdx[y]
			y++
		}
	}
	for ; x < c; x++ {
		newdata[z] = data[x]
		newidx[z] = sortedIdx[x]
		z++
	}
	for ; y < r; y++ {
		newdata[z] = data[y]
		newidx[z] = sortedIdx[y]
		z++
	}

	x = l
	z = 0
	for ; z < datalen; z++ {
		data[x] = newdata[z]
		sortedIdx[x] = newidx[z]
		x++
	}
}

/*
IndirectSortFloat64 will sort the data and return the sorted index.
*/
func IndirectSortFloat64(data []float64) (sortedIdx []int) {
	datalen := len(data)

	sortedIdx = make([]int, datalen)
	for i := 0; i < datalen; i++ {
		sortedIdx[i] = i
	}

	MergesortFloat64(data, sortedIdx, 0, datalen)

	return
}

/*
SortFloatSliceByIndex will sort the slice of float `data` using sorted index
`sortedIdx`.
*/
func SortFloatSliceByIndex(data *[]float64, sortedIdx *[]int) {
	newdata := make([]float64, len(*data))

	for i := range *sortedIdx {
		newdata[i] = (*data)[(*sortedIdx)[i]]
	}

	(*data) = newdata
}

/*
SortStringSliceByIndex will sort the slice of string `data` using sorted index
`sortedIdx`.
*/
func SortStringSliceByIndex(data *[]string, sortedIdx *[]int) {
	newdata := make([]string, len(*data))

	for i := range *sortedIdx {
		newdata[i] = (*data)[(*sortedIdx)[i]]
	}

	(*data) = newdata
}

/*
IntIsExist check if integer value exist in list of integer, return true if
exist, false otherwise.
*/
func IntIsExist(data []int, val int) bool {
	for _, v := range data {
		if val == v {
			return true
		}
	}
	return false
}

/*
GetRandomInteger return random integer value from 0 to maximum value `maxVal`.
The random value is checked with already picked index, `pickedIdx`.
If `dup` is true, allow duplicate value in `pickedIdx`, otherwise only single
unique value allowed in `pickedIdx`.
If excluding index `excIdx` is not empty, do not pick the integer value listed
in there.
*/
func GetRandomInteger(maxVal int, dup bool, pickedIdx []int, excIdx []int) (
	idx int,
) {
	rand.Seed(time.Now().UnixNano())

	for {
		idx = rand.Intn(maxVal)

		// check if its must not be selected
		excluded := false
		for _, excIdx := range excIdx {
			if idx == excIdx {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		if dup {
			// allow duplicate idx
			return
		}

		// check if its already picked
		isPicked := false
		for _, pastIdx := range pickedIdx {
			if idx == pastIdx {
				isPicked = true
				break
			}
		}
		// get another random idx again
		if isPicked {
			continue
		}

		// bingo, we found unique idx that has not been picked.
		return
	}
}

/*
CountMissRate given two slice of string, count number of string that is not
equal with each other, and return the miss rate
(number of not equal / number of data), miss number, and length of slice.
*/
func CountMissRate(src []string, target []string) (
	missrate float64,
	nmiss, length int,
) {
	// find minimum length
	length = len(src)
	targetlen := len(target)
	if targetlen < length {
		length = targetlen
	}

	for x := 0; x < length; x++ {
		if src[x] != target[x] {
			nmiss++
		}
	}

	return float64(nmiss) / float64(length), nmiss, length
}
