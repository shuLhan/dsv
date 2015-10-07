package dsv

/*
Record represent each field in record.
*/
type Record []byte

/*
String return the value of record in string enclosed with double quoted.
*/
func (record Record) String() string {
	return string (record)
}
