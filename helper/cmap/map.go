package cmap

import (
	"maps"
	"net/url"
	"strings"
)

// GetMapValue 获取map中的key
func GetMapValue[T comparable, A any](data map[T]A, property T, def A) A {
	if _, ok := data[property]; ok {
		return data[property]
	}
	return def
}

// MergeMap 多个map合并,后面的map同名key数据将覆盖前面的map。
func MergeMap[T comparable, A any](data map[T]A, d ...map[T]A) map[T]A {
	if len(d) > 0 {
		for _, val := range d {
			maps.Copy(data, val)
		}
	}
	return data
}

// MapDecode 将map中的数据urlDecode
func MapDecode[T comparable](data map[T]string) map[T]string {
	for kl, vl := range data {
		tmp, err := url.QueryUnescape(string(vl))
		if err != nil {
			tmp = vl
		}
		data[kl] = tmp
	}
	return data
}

// MapValues 提取map中的值，组成arr
func MapValues[T comparable, A any](data map[T]A) []A {
	var res = make([]A, 0)
	for _, v := range data {
		res = append(res, v)
	}
	return res
}

// MapKeys 提取map中的key，组成arr
func MapKeys[T comparable, A any](data map[T]A) []T {
	var res = make([]T, 0)
	for k := range data {
		res = append(res, k)
	}
	return res
}

// GetDeepMapValue 根据path获取map中的值
func GetDeepMapValue[A any](m map[string]interface{}, path string, defaultValue A) A {
	if m == nil || path == "" {
		return defaultValue
	}
	keys := strings.Split(path, ".")
	curMap := m
	for _, key := range keys[:len(keys)-1] {
		if curVal, ok := curMap[key]; ok {
			if nextMap, ok := curVal.(map[string]interface{}); ok {
				curMap = nextMap
			} else {
				return defaultValue
			}
		} else {
			return defaultValue
		}
	}
	if val, ok := curMap[keys[len(keys)-1]]; ok {
		if finalVal, ok := val.(A); ok {
			return finalVal
		}
	}
	return defaultValue
}
