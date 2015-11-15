// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package util contain common function to work with data.
*/
package util

const (
	// Threshold for insertion sort. When the data less than Threshold,
	// insertion sort will be used in IndirectSort instead of QuickSort.
	Threshold = 7
)

/*
InsertionSort will sort the data using insertion algorithm.
*/
func InsertionSort(data *[]float64, idx *[]int, l, r int) {
	var tmpIdx int
	var tmpData float64
	var i, j int

	for i = l + 1; i <= r; i++ {
		tmpIdx = (*idx)[i]
		tmpData = (*data)[i]

		for j = i; j > l && tmpData < (*data)[j-1]; j-- {
			(*idx)[j] = (*idx)[j-1]
			(*data)[j] = (*data)[j-1]
		}

		(*idx)[j] = tmpIdx
		(*data)[j] = tmpData
	}
}

/*
SwapInt swap two indices value of integer.
*/
func SwapInt(data *[]int, i, j int) {
	if i == j {
		return
	}

	tmp := (*data)[i]
	(*data)[i] = (*data)[j]
	(*data)[j] = tmp
}

/*
SwapFloat64 swap two indices value of 64bit float.
*/
func SwapFloat64(data *[]float64, i, j int) {
	if i == j {
		return
	}

	tmp := (*data)[i]
	(*data)[i] = (*data)[j]
	(*data)[j] = tmp
}

/*
SwapString swap two indices value of string.
*/
func SwapString(data *[]string, i, j int) {
	if i == j {
		return
	}

	tmp := (*data)[i]
	(*data)[i] = (*data)[j]
	(*data)[j] = tmp
}

/*
Median3 sort the left, center, and right data.
Return pivot.
*/
func Median3(data *[]float64, idx *[]int, l, r int) float64 {
	c := (l + r) / 2

	if (*data)[l] > (*data)[c] {
		SwapInt(idx, l, c)
		SwapFloat64(data, l, c)
	}
	if (*data)[l] > (*data)[r] {
		SwapInt(idx, l, r)
		SwapFloat64(data, l, r)
	}
	if (*data)[c] > (*data)[r] {
		SwapInt(idx, c, r)
		SwapFloat64(data, c, r)
	}

	SwapInt(idx, c, r-1)
	SwapFloat64(data, c, r-1)

	return (*data)[r-1]
}

/*
QuickSort will sort the data from left index `l` to right index `r` and save
the sorted index in idx.
*/
func QuickSort(data *[]float64, idx *[]int, l, r int) {
	var pivot float64
	var i, j int

	if l+Threshold >= r {
		InsertionSort(data, idx, l, r)
	} else {
		pivot = Median3(data, idx, l, r)

		i = l
		j = r - 1
		for {
			i++
			for (*data)[i] < pivot {
				i++
			}
			j--
			for (*data)[j] > pivot {
				j--
			}
			if i < j {
				SwapInt(idx, i, j)
				SwapFloat64(data, i, j)
			} else {
				break
			}
		}
		SwapInt(idx, i, r-1)
		SwapFloat64(data, i, r-1)

		QuickSort(data, idx, l, i-1)
		QuickSort(data, idx, i+1, r)
	}
}

/*
IndirectSort will sort the data and return the sorted index.
*/
func IndirectSort(data *[]float64) *[]int {
	r := make([]int, len(*data))
	for i := range r {
		r[i] = i
	}

	QuickSort(data, &r, 0, len(*data)-1)

	return &r
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
