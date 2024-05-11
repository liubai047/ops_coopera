package cmap

import (
	"fmt"
	"testing"
)

func Test_GetDeepMapValue(t *testing.T) {
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

	fmt.Println("a:", GetDeepMapValue(m, "a", 0))
	fmt.Println("a.b:", GetDeepMapValue[any](m, "a.b", nil))
	fmt.Println("a.b.c:", GetDeepMapValue(m, "a.b.c", ""))
	fmt.Println("d:", GetDeepMapValue(m, "d", []string{}))
	fmt.Println("d.a", GetDeepMapValue(m, "d.a", ""))
	fmt.Println("d.0:", GetDeepMapValue(m, "d.0", ""))
	fmt.Println("b.c:", GetDeepMapValue[any](m, "b.c", nil))
	fmt.Println("b.e.f:", GetDeepMapValue[any](m, "b.e.f", nil))
	fmt.Println("h.g:", GetDeepMapValue(m, "h.g", ""))
}
