package objectPool

import (
	"testing"
)

func TestByteObjPool(t *testing.T) {
	var p = NewObjPool(make([]byte, 10))
	pv := p.Get() // 对象池取对象
	p.Put(pv)     // 返回对象池
}
