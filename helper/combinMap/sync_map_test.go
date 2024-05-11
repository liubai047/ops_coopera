package combinMap

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_sync_map(t *testing.T) {
	var syncMap = MySyncMap[string, int]{}
	var tsMap = map[string]int{"dsfasdfasdf": 5253, "fiodnfgin": 63636, "sdfasdf": 4634}
	for k, v := range tsMap {
		syncMap.Store(k, v)
		if vv, ok := syncMap.Load(k); !ok || vv != v {
			t.Error("测试读写失败，读写结果不一致")
		}
	}
}

func Test_a(t *testing.T) {
	a := make([]map[int]int, 0)
	for i := 0; i < 1000; i++ {
		b := make(map[int]int)
		for j := 0; j < 10000; j++ {
			b[j] = j
		}
		a = append(a, b)
	}
	s1time := time.Now()
	goTest(a)
	e1time := time.Now()
	fmt.Println("go", e1time.Sub(s1time).Milliseconds())
}

func goTest(data []map[int]int) {
	i := 0
	var ch = make(chan int, 10000000)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(data))
		for _, item := range data {
			go func(item map[int]int) {
				defer wg.Done()
				for k, _ := range item {
					ch <- k
				}
			}(item)
		}
		wg.Wait()
		close(ch) // 关闭通道，结束for range遍历
	}()
	for _ = range ch {
		i++
	}
	fmt.Println(i)
}
