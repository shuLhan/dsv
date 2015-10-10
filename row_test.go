package dsv_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shuLhan/dsv"
)

var exp = []string {
	"0\n",
	"1\n",
	"2\n",
	"3\n",
	"4\n",
}

func TestRowPop (t *testing.T) {
	if DEBUG {
		fmt.Println (">>> TestRowPop")
	}

	rows := &dsv.Row {}

	for i := 0; i < 5; i++ {
		rows.PushBack (i)
	}

	exps := strings.Join (exp, "")
	got := fmt.Sprint (rows)

	if got != exps {
		t.Fatal ("Expecting:\n", exps, "\n Got:\n", got)
	}

	row := rows.Pop ()

	exps = strings.Join (exp[:1], "")
	got = fmt.Sprint (row)

	if got != exps {
		t.Fatal ("Expecting:\n", exps, "\n Got:\n", got)
	}

	exps = strings.Join (exp[1:], "")
	got = fmt.Sprint (rows)

	if got != exps {
		t.Fatal ("Expecting:\n", exps, "\n Got:\n", got)
	}

	rows.Pop ()
	rows.Pop ()
	rows.Pop ()
	rows.Pop ()

	row = rows.Pop ()

	if nil != row {
		t.Fatal ("Expecting:\n", nil, "\n Got:\n", row)
	}
}

/*
TestRowGroupByValue test grouping in rows.
*/
func TestRowGroupByValue (t *testing.T) {
	var exps = []string {
		`[[1 +] [2 -] [3 -] [4 +]]`,
		`&[1 +]
&[2 -]
&[3 -]
&[4 +]
`,
		`map[+:&[1 +]
&[4 +]
 -:&[2 -]
&[3 -]
]`,
	}

	var records = [][]dsv.Record {
		{{ int64(1), dsv.TInteger }, { "+", dsv.TString }},
		{{ int64(2), dsv.TInteger }, { "-", dsv.TString }},
		{{ int64(3), dsv.TInteger }, { "-", dsv.TString }},
		{{ int64(4), dsv.TInteger }, { "+", dsv.TString }},
	}
	var rows = &dsv.Row{}

	// test records
	got := fmt.Sprint (records)

	i := 0
	if got != exps[i] {
		t.Fatal ("Expecting:\n", exps[i], "\n Got:\n", got)
	}

	// test rows
	for i = range records {
		rows.PushBack (&records[i])
	}

	got = fmt.Sprint (rows)

	i = 1
	if got != exps[i] {
		t.Fatal ("Expecting:\n", exps[i], "\n Got:\n", got)
	}

	// test rows grouping
	classes := rows.GroupByValue (1)

	got = fmt.Sprint (classes)

	i = 2
	if got != exps[i] {
		t.Fatal ("Expecting:\n", exps[i], "\n Got:\n", got)
	}

	if DEBUG {
		for k,v := range classes {
			fmt.Println (k, v.Len (), v)
		}
	}
}
