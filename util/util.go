// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package util contain common function to work with data.
*/
package util

const (
	// SortThreshold for sort. When the data less than SortThreshold,
	// insertion sort will be used to replace the sort.
	SortThreshold = 7
)

/*
InsertionSortFloat64 will sort the data using insertion algorithm.
*/
func InsertionSortFloat64(data []float64, idx []int, l, r int) {
	for x := l; x <= r; x++ {
		for y := x + 1; y <= r; y++ {
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
MergeSortSliceFloat64 sort the slice of float from `l` to `r` using mergesort
algorithm, return the sorted index.
*/
func MergeSortSliceFloat64(data []float64, sortedIdx []int, l, r int) {
	if l + SortThreshold >= r {
		InsertionSortFloat64(data, sortedIdx, l, r)
		return
	}

	res := (r + l) % 2
	c := (r + l) / 2
	if res == 0 {
		c--
	}

	MergeSortSliceFloat64(data, sortedIdx, l, c)
	MergeSortSliceFloat64(data, sortedIdx, c + 1, r)

	// merging
	x := l
	y := c + 1
	for x <= c && y <= r {
		if data[x] <= data[y] {
			x++
		} else {
			SwapInt(sortedIdx, x, y)
			SwapFloat64(data, x, y)
			y++
		}
	}

	for x < y && y <= r {
		if data[x] > data[y] {
			SwapInt(sortedIdx, x, y)
			SwapFloat64(data, x, y)
			x++
		} else {
			y++
		}
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

	MergeSortSliceFloat64(data, sortedIdx, 0, datalen - 1)

	return
}

/*
SortFloatSliceByIndex will sort the slice of float `data` using sorted index
`sortedIdx`.
*/
func SortFloatSliceByIndex(data *[]float64, sortedIdx *[]int) {
	newdata := make([]float64, len(*data))

	for i := range (*sortedIdx) {
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

	for i := range (*sortedIdx) {
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
