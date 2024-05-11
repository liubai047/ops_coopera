package deepClone

import (
	"errors"
	"reflect"
)

/*Clone
* 深拷贝一个对象，能够正确处理：结构体,数组,slice,channel,map,interface,以及这些类型的指针类型
* 如果结构体中存在通道，则只会初始化同类型通道，但不会转移其已存在的值
* 已处理循环引用
 */
func Clone[T any](src T) (T, error) {
	var dst T
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.New(srcValue.Type()).Elem()
	visited := make(map[uintptr]reflect.Value)
	if err := deepCloneValue(srcValue, dstValue, visited); err != nil {
		return dst, err
	}
	dst, ok := dstValue.Interface().(T)
	if !ok {
		return *new(T), errors.New("转换异常")
	}
	return dst, nil
}

// visited 用于处理循环引用
func deepCloneValue(src, dst reflect.Value, visited map[uintptr]reflect.Value) error {
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			if err := deepCloneValue(src.Field(i), dst.Field(i), visited); err != nil {
				return err
			}
		}
	case reflect.Array:
		if !dst.CanSet() {
			return errors.New("deepClone cannot set array")
		}
		dst.Set(reflect.New(src.Type()).Elem())
		for i := 0; i < src.Len(); i++ {
			if err := deepCloneValue(src.Index(i), dst.Index(i), visited); err != nil {
				return err
			}
		}
	case reflect.Slice:
		if src.IsNil() {
			return nil
		}
		if !dst.CanSet() {
			return errors.New("deepClone cannot set slice")
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Len()))
		for i := 0; i < src.Len(); i++ {
			if err := deepCloneValue(src.Index(i), dst.Index(i), visited); err != nil {
				return err
			}
		}
	case reflect.Map:
		if src.IsNil() {
			return nil
		}
		if !dst.CanSet() {
			return errors.New("deepClone cannot set map")
		}
		dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
		for _, key := range src.MapKeys() {
			srcVal := src.MapIndex(key)
			dstVal := reflect.New(srcVal.Type()).Elem()
			if err := deepCloneValue(srcVal, dstVal, visited); err != nil {
				return err
			}
			dst.SetMapIndex(key, dstVal)
		}
	case reflect.Pointer:
		if src.IsNil() {
			return nil
		}
		if !dst.CanSet() {
			return errors.New("deepClone cannot set pointer")
		}
		// 检查是否已经访问过这个指针 -- 解决递归调用
		if dstPointer, ok := visited[src.Pointer()]; ok {
			// 如果已经访问过，直接使用映射中的拷贝
			dst.Set(dstPointer)
			return nil
		}
		dst.Set(reflect.New(src.Type().Elem()))
		if err := deepCloneValue(src.Elem(), dst.Elem(), visited); err != nil {
			return err
		}
	case reflect.Interface:
		if src.IsNil() {
			return nil
		}
		srcValue := src.Elem()
		dstValue := reflect.New(srcValue.Type()).Elem()
		if err := deepCloneValue(srcValue, dstValue, visited); err != nil {
			return err
		}
		dst.Set(dstValue)
	case reflect.Chan:
		dst.Set(reflect.MakeChan(src.Type(), src.Cap()))
	case reflect.UnsafePointer, reflect.Func: // 这两个类型不处理
		return nil
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
	return nil
}
