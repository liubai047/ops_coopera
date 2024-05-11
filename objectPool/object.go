package objectPool

import (
	"bytes"
	"reflect"
	"sync"
)

type objTyps interface {
	~[]byte | bytes.Buffer // 限定类型，可以按照需求进行增加
}

type ObjPool[A objTyps] struct {
	pool sync.Pool
}

// NewObjPool 初始化对象池时，请将想要初始化的类型格式传入，将按照传入的格式进行类型初始化
func NewObjPool[A objTyps](v A) *ObjPool[A] {
	var valTyps = reflect.ValueOf(v)
	var initFunc func() any
	switch valTyps.Kind() {
	// 利用反射，根据不同格式数据按对应方式进行初始化(将需要make初始化的写完，就只剩下new初始化的了)
	case reflect.Slice:
		initFunc = func() any {
			valTyps.Type()
			return reflect.MakeSlice(valTyps.Type(), valTyps.Len(), valTyps.Cap())
		}
	case reflect.Map:
		initFunc = func() any {
			return reflect.MakeMapWithSize(valTyps.Type(), valTyps.Len())
		}
	case reflect.Chan:
		initFunc = func() any {
			return reflect.MakeChan(valTyps.Type(), valTyps.Cap())
		}
	default:
		initFunc = func() any {
			return new(A)
		}
	}
	return &ObjPool[A]{
		pool: sync.Pool{
			New: initFunc,
		},
	}
}

func (obj *ObjPool[A]) Get() A {
	return *obj.pool.Get().(*A)
}

func (obj *ObjPool[A]) Put(v A) {
	obj.pool.Put(v)
}
