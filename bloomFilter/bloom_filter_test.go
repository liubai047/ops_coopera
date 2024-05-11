package bloomFilter

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

var bloomTest = NewBloom(10000000, 0.001)
var testTexts = make([][]byte, 0)

func init() {
	for i := 0; i < 100000; i++ {
		testTexts = append(testTexts, []byte(string(rand.Int31())+strconv.Itoa(i)))
	}
	println(bloomTest.bitNum, "位数量")
	println(len(bloomTest.hashFunc), "哈希函数数量")
}

// 1.3M内存 0.83S
func TestBloom_Add(b *testing.T) {
	var hashConflict = 0 // bloom冲突数量
	var bloomErr = 0     // bloom过滤器异常次数
	for _, testText := range testTexts {
		if bloomTest.Check(testText) {
			hashConflict += 1
			// fmt.Printf("%d %s is in bloom before\n", i, testText)
		}
		_ = bloomTest.Add(testText)
		if !bloomTest.Check(testText) {
			bloomErr += 1
			// b.Errorf("%d %s this is not in bloom after testText??", i, testText)
		}
	}
	fmt.Println(hashConflict, "冲突次数")
	fmt.Println(float64(hashConflict)/float64(len(testTexts)), "冲突率")
	fmt.Println(bloomErr, "异常次数")
	fmt.Println("finish")
}

// 10w数据，检测+写入 花费0.8秒（单检测花费0.03秒）
func Benchmark_Add(b *testing.B) {
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		for _, testText := range testTexts {
			if bloomTest.Check(testText) {
				// 	// fmt.Printf("%d %s is in bloom before\n", i, testText)
			}
			_ = bloomTest.Add(testText)
			// if !bloomTest.Check(testText) {
			// 	b.Errorf("%d %s this is not in bloom after testText??", i, testText)
			// }
		}
		// fmt.Println("finish")
	}
}

// func TestBloom_Add(t *testing.T) {
// 	bloom := NewBloom(10000, 0.01)
// 	var testTexts = make([]string, 0)
// 	for i := 0; i < 1; i++ {
// 		testTexts = append(testTexts, grand.S(10)+strconv.Itoa(i))
// 	}
// 	for i, testText := range testTexts {
// 		if bloom.Check(testText) {
// 			fmt.Printf("%d %s is in bloom before\n", i, testText)
// 		}
// 		_ = bloom.Add(testText)
// 		if !bloom.Check(testText) {
// 			t.Errorf("%d %s this is not in bloom after testText??", i, testText)
// 		}
// 		fmt.Println("success")
// 	}
// }
