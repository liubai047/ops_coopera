package deepClone

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// GobClone 使用gob加解码器的方式进行对象克隆
//
// tips：加解码对象都只能在go中进行，如果有不可复制的结构体，会触发报错。此方法性能略高于json.Unmarshal的方式。
func GobClone[T any](src T) (T, error) {
	var buf bytes.Buffer
	var nVal T
	err := gob.NewEncoder(&buf).Encode(src)
	if err != nil {
		return src, err
	}
	err = gob.NewDecoder(&buf).Decode(nVal)
	if err != nil {
		return src, err
	}
	return nVal, nil
}

// JsonClone 使用json序列化的方式进行对象克隆
//
// tips：此方法在跨网络传输过程中的兼容性比gob方法好。如果有不可复制的结构体，会触发报错。
func JsonClone[T any](src T) (T, error) {
	var nVal T
	bts, err := json.Marshal(src)
	if err != nil {
		return nVal, err
	}
	err = json.Unmarshal(bts, &nVal)
	if err != nil {
		return src, err
	}
	return nVal, nil
}
