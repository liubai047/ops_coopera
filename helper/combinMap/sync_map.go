package combinMap

import "sync"

type MySyncMap[k comparable, v any] struct {
	m sync.Map
}

func (my *MySyncMap[k, v]) Store(key k, value v) {
	my.m.Store(key, value)
}

func (my *MySyncMap[k, v]) Load(key k) (v, bool) {
	val, ok := my.m.Load(key)
	if !ok {
		var res v
		return res, ok
	}
	return val.(v), ok
}

func (my *MySyncMap[k, v]) Delete(key k) {
	my.m.Delete(key)
}

// Range 遍历map，当函数返回false时，停止后续遍历
func (my *MySyncMap[k, v]) Range(f func(key k, value v) bool) {
	my.m.Range(func(key, value any) bool {
		return f(key.(k), value.(v))
	})
}
