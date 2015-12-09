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
        ],{
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

    // we will make sure all open descriptor is closed.
    defer dsvReader.Close()

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

Easy enough? We can combine the reader and writer using `dsv.New()`, which will
create reader and writer,

    rw, e := dsv.New("config.dsv")

    if nil != e {
        t.Error(e)
    }

    // do usual process like in the last step.

Thats it!
