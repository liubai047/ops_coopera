package deepClone

import "fmt"

func ExampleClone() {
	type AStruct struct {
		A int
		B string
		C float64
		D bool
	}

	type NestedStruct struct {
		NestedField1 int
	}
	type MyStruct struct {
		Field1 int
		Field2 *NestedStruct
		Field3 []string
		Field4 [3]int
		Field5 map[string]interface{}
		Field6 interface{}
		Field7 chan string
		Field8 AStruct
		x      float64
	}
	ch := make(chan string, 10)
	ch <- "milo1"
	ch <- "milo2"
	// 示例结构体
	original := MyStruct{
		Field1: 10,
		Field2: &NestedStruct{NestedField1: 20},
		Field3: []string{"a", "b", "c"},
		Field4: [3]int{9, 3, 5},
		Field5: map[string]interface{}{"name": "milopeng", "age": 21, "isHandsome": true},
		Field6: "这是一个很棒的结构体",
		Field7: ch,
		Field8: AStruct{
			A: 9527,
			B: "good job",
			C: 3.88,
			D: false,
		},
		x: 21.1,
	}

	// 执行深拷贝
	copied, err := Clone(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", original)
	fmt.Printf("%#v\n", copied)

	// 验证深拷贝结果
	println(copied.Field1 == original.Field1)                           // true
	println(copied.Field2 != original.Field2)                           // true
	println(copied.Field2.NestedField1 == original.Field2.NestedField1) // true
	println(copied.Field3[0] == original.Field3[0])                     // true

	// Output:
	// deepClone.MyStruct{Field1:10, Field2:(*deepClone.NestedStruct)(0xc00000a380), Field3:[]string{"a", "b", "c"}, Field4:[3]int{9, 3, 5}, Field5:map[string]interface {}{"age":21, "isHandsome":true, "name":"milopeng"}, Field6:"这是一个很棒的结构体", Field7:(chan string)(0xc00005e0c0), Field8:deepCl          one.AStruct{A:9527, B:"good job", C:3.88, D:false}, x:21.1}
	// deepClone.MyStruct{Field1:10, Field2:(*deepClone.NestedStruct)(0xc00000a388), Field3:[]string{"a", "b", "c"}, Field4:[3]int{9, 3, 5}, Field5:map[string]interface {}{"age":21, "isHandsome":true, "name":"milopeng"}, Field6:"这是一个很棒的结构体", Field7:(chan string)(0xc00005e120), Field8:deepC          lone.AStruct{A:9527, B:"good job", C:3.88, D:false}, x:0}
	// true
	//
	//
	//
}
