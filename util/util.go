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
StringCountBy count number of occurence of `class` values in data.
Return number of each class based on their index.

For example, if data is "[A,A,B]" and class is "[A,B]", this function will
return "[2,1]".

	idx cls  count
	0 : A -> 2
	1 : B -> 1
*/
func StringCountBy(data []string, class []string) (clsCnt []int) {
	clsCnt = make([]int, len(class))

	for _, r := range data {
		for k, v := range class {
			if r == v {
				clsCnt[k]++
				break
			}
		}
	}

	return
}

/*
IntFindMax given slice of integer, return the maximum value in slice and index
of maximum value.
If data is empty, return -1 in value and index.
*/
func IntFindMax(data []int) (max int, maxidx int) {
	l := len(data)
	if l <= 0 {
		return -1, -1
	}

	i := 0
	max = data[i]
	maxidx = i

	for i = 1; i < l; i++ {
		if data[i] > max {
			max = data[i]
			maxidx = i
		}
	}

	return
}

/*
IntFindMin given slice of integer, return the minimum value in slice and index
of minimum value.
If data is empty, return -1 in value and index.
*/
func IntFindMin(data []int) (min int, minidx int) {
	l := len(data)
	if l <= 0 {
		return -1, -1
	}

	i := 0
	min = data[i]
	minidx = i

	for i = 1; i < l; i++ {
		if data[i] < min {
			min = data[i]
			minidx = i
		}
	}

	return
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
