package carr

import (
	"cmp"
	"sort"
)

type sortArrMap[T comparable, K cmp.Ordered] []map[T]K

// NewSortArrMap 初始化value值可被比较的map，数组排序使用
func NewSortArrMap[T cmp.Ordered, K cmp.Ordered](data []map[T]K) sortArrMap[T, K] {
	return data
}

// SortMap 对slice map排序 Desc为true时降序，反之为升序
func (m sortArrMap[T, K]) SortMap(column T, Desc bool) []map[T]K {
	if len(m) < 2 {
		return m
	}
	sort.Slice(m, func(i, j int) bool {
		// 防止访问不存在的column，导致异常
		if _, ok := m[i][column]; !ok {
			return false
		}
		if _, ok := m[j][column]; !ok {
			return false
		}
		// 如果出现可比较的K类型，则执行排序比较
		if Desc {
			return cmp.Compare(m[i][column], m[j][column]) > 0
		} else {
			return cmp.Compare(m[i][column], m[j][column]) < 0
		}
	})
	return m
}
