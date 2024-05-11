package combinMap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func stringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}
func Test_CurrentMap(t *testing.T) {
	// 测试读写
	tsCurRWMap := NewCurrentMap[string, int](5)
	var tsMap = map[string]int{"dsfasdfasdf": 5253, "fiodnfgin": 63636, "sdfasdf": 4634}
	for k, v := range tsMap {
		tsCurRWMap.Set(k, v)
		if tsCurRWMap.MustGet(k) != v {
			t.Error("测试读写失败，读写结果不一致")
		}
	}

	// 测试并发读写
	tsCursMap := NewCurrentMap[string, any](5000)
	var wg = sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				tsCursMap.Set(randomString(10), rand.Int())
			}
		}()
	}
	wg.Wait()
	println(tsCursMap.Count())
	now := time.Now()
	kLen1 := len(tsCursMap.Keys1())
	afTime := time.Now()
	fmt.Printf("key总数为：%d,key1 方法总耗时：%dms\n", kLen1, afTime.Sub(now).Milliseconds())
	now = time.Now()
	kLen2 := len(tsCursMap.Keys1())
	afTime = time.Now()
	fmt.Printf("key总数为：%d,key2 方法总耗时：%dms\n", kLen2, afTime.Sub(now).Milliseconds())
}
