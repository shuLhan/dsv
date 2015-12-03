package dsv_test

import (
	"os"
)

var DEBUG = bool (os.Getenv ("DEBUG") != "")

var expectation = []string {
	"&[1 A-B AB 1 0.1]\n",
	"&[2 A-B-C BCD 2 0.02]\n",
	"&[3 A;B-C,D A;B C,D 3 0.003]\n",
	"&[6   6 0.000006]\n",
	"&[9 ok ok 9 0.000000009]\n",
	"&[10 test integer 10 0.101]\n",
	"&[12 test real 123456789 0.123456789]\n",
}

var exp_skip = []string {
	"&[A-B AB 1 0.1]\n",
	"&[A-B-C BCD 2 0.02]\n",
	"&[A;B-C,D A;B C,D 3 0.003]\n",
	"&[  6 0.000006]\n",
	"&[ok ok 9 0.000000009]\n",
	"&[test integer 10 0.101]\n",
	"&[test real 123456789 0.123456789]\n",
}

var exp_skip_fields = []string {
	"[[A-B] [AB] [1] [0.1]]",
	"[[A-B-C] [BCD] [2] [0.02]]",
	"[[A;B-C,D] [A;B C,D] [3] [0.003]]",
	"[[] [] [6] [0.000006]]",
	"[[ok] [ok] [9] [0.000000009]]",
	"[[test] [integer] [10] [0.101]]",
	"[[test] [real] [123456789] [0.123456789]]",
}
