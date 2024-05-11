package cmap

import (
	"fmt"
)

func ExampleGetDeepMapValue() {
	m := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": "2",
			"e": map[string]interface{}{
				"f": "6",
			},
		},
		"d": []string{"1", "2", "3"},
		"h": map[string]interface{}{
			"g": nil,
		},
	}
	a := GetDeepMapValue(m, "a", 0)
	b := GetDeepMapValue[any](m, "a.b", nil)
	fmt.Println(a, b)
	// Output:
	// 1, nil
}
