package carr

import (
	"math"
	"time"
)

type mArr[T comparable] []T

func NewArr[T comparable](arr []T) mArr[T] {
	return arr
}

// Find 查找slice中是否存在val，返回索引
func (m mArr[T]) Find(val T) (int, bool) {
	for i, item := range m {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// InArray 查看某个元素是否在数组中
func (m mArr[T]) InArray(val T) bool {
	for _, item := range m {
		if item == val {
			return true
		}
	}
	return false
}

// ArrayChunk 切片切分
func (m mArr[T]) ArrayChunk(size int) [][]T {
	if size <= 0 {
		return append(make([][]T, 0), m)
	}
	length := len(m)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]T
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, m[i*size:end])
		i++
	}
	return n
}

// ArrayUnique RemoveDuplicate 数组去重
func (m mArr[T]) ArrayUnique() []T {
	set := make(map[T]struct{}, len(m))
	j := 0
	for _, v := range m {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		m[j] = v
		j++
	}
	return m[:j]
}

func (m mArr[T]) RandSlice() T {
	return m[time.Now().UnixMicro()%int64(len(m))]
}
