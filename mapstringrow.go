package dsv

import (
	"math"
)

/*
MapStringRow represent mapping between string (key) with rows (value).
*/
type MapStringRow map[string]*Row

/*
GetMinority return group in groups with minimum rows.
*/
func (groups *MapStringRow) GetMinority () (minorGroup *Row) {
	var min = math.MaxInt32

	for k, v := range *groups {

		if (*groups)[k].Len () < min {
			minorGroup = v
		}
	}

	return
}
