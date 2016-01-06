[![GoDoc](https://godoc.org/github.com/shuLhan/dsv?status.svg)](https://godoc.org/github.com/shuLhan/dsv)

Package `dsv` is a Go library for working with delimited separated value (DSV).

DSV is a free-style form of CSV format of text data, where each record is
separated by newline, and each column can be separated by any string, not just
comma.

## Example

Lets process this input file `input.dat`,

    Mon Dt HH MM SS Process
    Nov 29 23:14:36 process-1
    Nov 29 23:14:37 process-2
    Nov 29 23:14:38 process-3

and generate output file `output.dat` which format like this,

    "process_1","29-Nov"
    "process_2","29-Nov"
    "process_3","29-Nov"

How do we do it?

First, create file metadata for input and output, name it `config.dsv`,

    {
        "Input"         :"input.dat"
    ,   "Skip"          :1
    ,   "InputMetadata" :
        [{
            "Name"      :"month"
        ,   "Separator" :" "
        },{
            "Name"      :"date"
        ,   "Separator" :" "
        ,   "Type"      :"integer"
        },{
            "Name"      :"hour"
        ,   "Separator" :":"
        ,   "Type"      :"integer"
        },{
            "Name"      :"minute"
        ,   "Separator" :":"
        ,   "Type"      :"integer"
        },{
            "Name"      :"second"
        ,   "Separator" :" "
        ,   "Type"      :"integer"
        },{
            "Name"      :"process_name"
        ,   "Separator" :"-"
        },{
            "Name"      :"process_id"
        }]
    ,   "Output"        :"output.dat"
    ,   "OutputMetadata":
        [{
            "Name"      :"process_name"
        ,   "LeftQuote" :"\""
        ,   "Separator" :"_"
        },{
            "Name"      :"process_id"
        ,   "RightQuote":"\""
        ,   "Separator" :","
        },{
            "Name"      :"date"
        ,   "LeftQuote" :"\""
        ,   "Separator" :"-"
        },{
            "Name"      :"month"
        ,   "RightQuote":"\""
        }]
    }

The metadata is using JSON format. For more information see `metadata.go`
and `reader.go`.

Second, we create a reader to read the input file.

    dsvReader, e := dsv.NewReader("config.dsv")

    if nil != e {
        t.Fatal(e)
    }

Third, we create a writer to write our output data,

    dsvWriter, e := dsv.NewWriter("config.dsv")

    if nil != e {
        t.Error(e)
    }

Last action, we process them: read input records and pass them to writer.

    for {
        n, e := dsv.Read(dsvReader)

        if n > 0 {
            dsvWriter.Write(dsvReader)

        // EOF, no more record.
        } else if e == io.EOF {
            break
        }
    }

    // we will make sure all open descriptor is closed.
    _ = dsvReader.Close()

Easy enough? We can combine the reader and writer using `dsv.New()`, which will
create reader and writer,

    rw, e := dsv.New("config.dsv")

    if nil != e {
        t.Error(e)
    }

    // do usual process like in the last step.

Thats it!

## Terminology

Here are some terminology that we used in developing this library, which may
help reader understanding the configuration and API.

* Dataset: is a content of file
* Record: a single cell in row or column, or the smallest building block of
dataset
* Row: is a horizontal representation of records in dataset
* Column: is a vertical representation of records in dataset

```
       COL-0  COL-1  ... COL-x
ROW-0: record record ... record
ROW-1: record record ... record
...
ROW-y: record record ... record
```

## Configuration

We choose and use JSON for configuration because,

1. No additional source to test.
2. Easy to extended. User can embed the current metadata, add additional
configuration, and create another reader to work with it.

### Metadata

Metadata contain information about each column when reading input file and
writing to output file,

* `Name`: mandatory, the name of column
* `Type`: optional, type of record when reading input file. Valid value are
"integer", "real", or "string" (default)
* `Separator`: optional, default to `"\n"`. Separator is a string that
separate the current record with the next record.
* `LeftQuote`: optional, default is empty `""`. LeftQuote is a string that
start at the beginning of record.
* `RightQuote`: optional, default is empty `""`. RightQuote is a string at the
end of record.
* `Skip`: optional, boolean, default is `false`. If true the column will be
saved in dataset when reading input file, otherwise it will be ignored.
* `ValueSpace`: optional, slice of string, default is empty. This contain the
string representation of all possible value in column.

### Input

Input configuration contain information about input file.

* `Input`: mandatory, the name of input file, could use relative or absolute
path. If no path is given then it assumed that the input file is in the same
directory with configuration file.
* `InputMetadata`: mandatory, list of metadata.
* `Skip`: optional, number, default 0. Skip define the number of line that will
be skipped when first input file is opened.
* `Rejected`: optional, default to `rejected.dat`. Rejected is file where
data that does not match with metadata will be saved. One can inspect the
rejected file fix it for re-process or ignore it.
* `MaxRows`: optional, default to `256`. Maximum number of rows for one read
operation that will be saved in memory. If its negative, i.e. `-1`, all data
in input file will be processed.
* `DatasetMode`: optional, default to "rows". Mode of dataset in memory.
Valid values are "rows", "columns", or "matrix". Matrix mode is combination of
rows and columns, it give more flexibility when processing the dataset but
will require additional memory.

#### `DatasetMode` explained

For example, given input data file,

    a,b,c
    1,2,3

"rows" mode is where each line saved in its own slice, resulting in Rows:

    Rows[0]: [a b c]
    Rows[1]: [1 2 3]

"columns" mode is where each line saved by columns, resulting in Columns:

    Columns[0]: [a 1]
    Columns[1]: [b 2]
    Columns[1]: [c 3]

"matrix" mode is where each record saved both in row and column.

### Output

Output configuration contain information about output file when writing the
dataset.

* `Output`: mandatory, the name of output file, could use relative or absolute
path. If no path is given then it assumed that the output file is in the same
directory with configuration file.
* `OutputMetadata`: mandatory, list of metadata.

## Working with DSV

### Process each Rows/Columns

After opening the input file, we can process the dataset based on rows/columns
mode using simple `for` loop. Example,

```
for {
	n, e := dsv.Read(dsvReader)

	if e == io.EOF {
		break
	}
	if e != nil {
		// handle error
	}
	if n > 0 {
		for i, row := dsvReader.GetDataAsRows() {
			// process each row
		}

		for i, column := dsvReader.GetDataAsColumns() {
			// process each column
		}

		// Write the dataset to file after processed
		dsvWriter.Write(dsvReader)
	}
}
```

### Builtin Functions for Dataset

For more information see `dataset.go`.

* Random pick rows with or without replacement
* Random pick columns with or without replacement
* Sort all columns using index from indirect-sort
* Splitting dataset using numeric values
* Splitting dataset using categorical values
