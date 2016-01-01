// Copyright 2015 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assert_test

import (
	"github.com/shuLhan/dsv/util/assert"
	"testing"
)

type Map struct {
	K int
	V string
}

var data = []Map{
	{1, "1"},
	{2, "2"},
	{3, "3"},
}

func TestEqual(t *testing.T) {
	assert.Equal(t, 1, 1)
	assert.Equal(t, 1.234, 1.234)
	assert.Equal(t, "1", "1")

	var datacmp []Map
	datacmp = make([]Map, len(data))
	copy(datacmp, data)

	assert.Equal(t, data, datacmp)
}

func TestNotEqual(t *testing.T) {
	assert.NotEqual(t, 1, 2)
	assert.NotEqual(t, 1.234, 1.2345)
	assert.NotEqual(t, "1", "1 ")

	var datacmp []Map
	datacmp = make([]Map, len(data))
	copy(datacmp, data)

	datacmp[0].K = 0
	assert.NotEqual(t, data, datacmp)
}

func TestEqualFileContent(t *testing.T) {
	assert.EqualFileContent(t, "assert.go", "assert.go")
}
