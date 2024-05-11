package carr

import (
	"fmt"
	"testing"
)

func TestArrMap_SortMap(t *testing.T) {
	var data = []map[string]int64{
		{"a": 1},
		{"a": 2},
		{"a": 3},
		{"a": 102},
		{"a": 72},
	}
	fmt.Printf("%#v\n", NewArrMap(data).ArrayColumn("a"))
}

func TestSortArrMap_SortMap(t *testing.T) {
	var data = []map[string]int64{
		{"a": 1},
		{"a": 2},
		{"a": 3},
		{"a": 102},
		{"a": 72},
	}
	fmt.Printf("%v\n", NewSortArrMap(data).SortMap("a", false))
}
