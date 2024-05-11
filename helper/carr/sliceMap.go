package carr

type arrMap[T comparable, K any] []map[T]K

// NewArrMap 初始化一个数组map
func NewArrMap[T comparable, K any](data []map[T]K) arrMap[T, K] {
	return data
}

// ArrayColumn 查找sliceMap中，某一个列的值
func (m arrMap[T, K]) ArrayColumn(column T) []K {
	var tmp = make([]K, 0)
	for _, v := range m {
		if val, ok := v[column]; ok {
			tmp = append(tmp, val)
		}
	}
	return tmp
}

// ArrayToMap 转换sliceMap,以某个column作为key，值为value
// func (m arrMap[T, K]) ArrayToMap(column T) map[K]map[T]K {
// 	var tmp = make(map[K]map[T]K)
// 	for _, v := range m {
// 		val, ok := v[column]
// 		if !ok {
// 			return nil
// 		}
// 		tmp[val] = v
// 	}
// 	return tmp
// }
