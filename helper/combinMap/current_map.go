package combinMap

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"sync"
)

// mCurrentMap 并发安全map
type mCurrentMap[k comparable, v any] struct {
	arr   []*currentMapShard[k, v] // 分片数组
	total int64                    // key总量
}

// currentMapShard 子map分片
type currentMapShard[k comparable, v any] struct {
	m  map[k]v
	mu sync.RWMutex
}

var (
	maxShardNum = 500
	minShardNum = 4
)

// NewCurrentMap 初始化并发安全map，最大不能超过500(值太大，意义偏小),最小不能超过4
func NewCurrentMap[k comparable, v any](num int) *mCurrentMap[k, v] {
	var shardNum = num
	if num > maxShardNum {
		shardNum = maxShardNum
	}
	if num < minShardNum {
		shardNum = minShardNum
	}
	currentMap := new(mCurrentMap[k, v])
	mp := make([]*currentMapShard[k, v], shardNum)
	for i := 0; i < shardNum; i++ {
		mp[i] = &currentMapShard[k, v]{
			m: make(map[k]v),
		}
	}
	currentMap.arr = mp
	return currentMap
}

// getSharedMap 获取key对应的map分片
func (m *mCurrentMap[k, v]) getSharedMap(key k) *currentMapShard[k, v] {
	return m.arr[uint(hash(key))%uint(len(m.arr))]
}

// hash函数,一旦出现hash错误，将key默认置为第一个map中
func hash[kk comparable](key kk) uint32 {
	// 创建一个新的FNV哈希
	h := fnv.New32a()
	// 使用gob来处理任意类型的编码
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return 0
	}
	// 将编码后的字节写入哈希
	_, err = h.Write(buf.Bytes())
	if err != nil {
		return 0
	}
	// 返回哈希值
	return h.Sum32()
}

// Set 设置key,value
func (m *mCurrentMap[k, v]) Set(key k, value v) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	sharedMap.mu.Lock()              // 加锁(全锁定)
	sharedMap.m[key] = value         // 赋值
	m.total++                        // 总量+1
	sharedMap.mu.Unlock()            // 解锁

}

// Get 获取key对应的value
func (m *mCurrentMap[k, v]) Get(key k) (value v, ok bool) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	sharedMap.mu.RLock()             // 加锁(读锁定)
	value, ok = sharedMap.m[key]     // 取值
	sharedMap.mu.RUnlock()           // 解锁
	return value, ok
}

// MustGet 获取key失败时，返回该类型默认值
func (m *mCurrentMap[k, v]) MustGet(key k) v {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	sharedMap.mu.RLock()             // 加锁(读锁定)
	defer sharedMap.mu.RUnlock()     // 解锁
	return sharedMap.m[key]          // 取值
}

// Count 统计key个数
func (m *mCurrentMap[k, v]) Count() int64 {
	return m.total
}

// Delete 根据key删除元素
func (m *mCurrentMap[k, v]) Delete(key k) {
	sharedMap := m.getSharedMap(key) // 找到对应的分片map
	sharedMap.mu.Lock()
	defer sharedMap.mu.Unlock()
	delete(sharedMap.m, key)
}

// Keys1 所有的key方法1(方法:遍历每个分片map,读取key;缺点:量大时,阻塞时间较长)
func (m *mCurrentMap[k, v]) Keys1() []k {
	var keys = make([]k, m.Count())
	// 遍历所有的分片map
	var keyNums = 0
	for i := 0; i < len(m.arr); i++ {
		m.arr[i].mu.RLock() // 加锁(读锁定)
		for key := range m.arr[i].m {
			keys[keyNums] = key
			keyNums++
		}
		m.arr[i].mu.RUnlock() // 解锁
	}
	return keys
}

// Keys2 所有的key方法2(方法:开多个协程分别对分片map做统计再汇总 优点:量大时,阻塞时间较短)
func (m *mCurrentMap[k, v]) Keys2() []k {
	var keys = make([]k, m.Count())
	var ch = make(chan k, m.Count())
	// 开始遍历
	go func() {
		wg := sync.WaitGroup{}
		for i := 0; i < len(m.arr); i++ {
			wg.Add(1)
			// 每个分片map,单独起一个协程进行统计
			go func(ms *currentMapShard[k, v]) {
				defer wg.Done()
				ms.mu.RLock()         // 加锁(读锁定)
				defer ms.mu.RUnlock() // 解锁
				for ky := range ms.m {
					ch <- ky // 压入通道
				}
			}(m.arr[i])
		}
		wg.Wait()
		close(ch) // 关闭通道，结束for range遍历
	}()
	// 遍历通道,压入所有的key
	idx := 0
	for ky := range ch {
		keys[idx] = ky
		idx++
	}
	return keys
}
