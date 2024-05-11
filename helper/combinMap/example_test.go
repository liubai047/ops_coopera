package combinMap

import "fmt"

func ExampleNewCurrentMap() {
	curMap := NewCurrentMap[string, int](5)
	curMap.Set("a", 5253)
	curMap.Set("b", 63636)
	curMap.Set("c", 4634)
	fmt.Println(curMap.Get("a"))
	fmt.Println(curMap.MustGet("a"))
	fmt.Println(curMap.Count())
	fmt.Println(curMap.Keys2())
	fmt.Println(curMap.Keys1())
	curMap.Delete("a")
	fmt.Println(curMap.Get("a"))
	fmt.Println(curMap.MustGet("a"))
	// Output:
	// 5253 true
	// 5253
	// 3
	// [c b a]
	// [a c b]
	// 0 false
	// 0
}

func ExampleMySyncMap() {
	var syncMap = MySyncMap[string, int]{}
	syncMap.Store("a", 1)
	syncMap.Store("b", 2)
	syncMap.Store("c", 3)
	fmt.Println(syncMap.Load("a"))
	syncMap.Delete("a")
	fmt.Println(syncMap.Load("a"))
	syncMap.Range(func(key string, value int) bool {
		fmt.Println(key, value)
		return true
	})
}
