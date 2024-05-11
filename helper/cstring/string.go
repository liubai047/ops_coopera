package cstring

import (
	"math"
	"net/url"
	"strings"
	"unicode"
	"unsafe"
)

// SnakeCase 将驼峰字符串转换为下划线写法
func SnakeCase(str string) string {
	var snakeCase strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			snakeCase.WriteRune('_')
		}
		snakeCase.WriteRune(unicode.ToLower(r))
	}
	return snakeCase.String()
}

// CamelCase 将下划线字符串转换为驼峰写法
func CamelCase(str string) string {
	// 将字符串分割为单词
	words := strings.Split(str, "_")
	for i, word := range words {
		// 将每个单词的首字母大写
		words[i] = strings.Title(word)
	}
	// 将单词连接成一个字符串
	return strings.Join(words, "")
}

// HasSuffixes 字符串只能给是否包含后缀，在给定的数组中，有的话返回true并返回该后缀
func HasSuffixes(str string, subStrs []string) (bool, string) {
	for _, subStr := range subStrs {
		if strings.HasSuffix(str, subStr) {
			return true, subStr
		}
	}
	return false, ""
}

func QueryUnescape(str string) string {
	ss, err := url.QueryUnescape(str)
	if err != nil {
		return str
	}
	return ss
}

// Str2Bytes 慎用！！！将字符串高效的转换为[]byte类型。注：不可修改转换后的byte数组，会导致进程直接崩溃，无法recover()
func Str2Bytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

// StrUnique 字符串去重，如：test;test1;test; 去重：test;test1;
func StrUnique(str string, sep ...string) string {
	if str == "" {
		return str
	}
	if len(sep) == 0 {
		// 设置默认分割符
		sep = append(sep, ";")
	}
	strs := strings.Split(str, sep[0])
	result := make([]string, 0, len(strs))
	// 空struct不占内存空间
	temp := map[string]struct{}{}
	for _, item := range strs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	str = strings.Join(result, sep[0])
	return str
}

/*
StrPad

	*字符串长度不足，用指定的字符前后补全
	*Example:
	*input := "Codes";
	*StrPad(input, 10, " ", "RIGHT")        // produces "Codes     "
	*StrPad(input, 10, "-=", "LEFT")        // produces "=-=-=Codes"
	*StrPad(input, 10, "_", "BOTH")         // produces "__Codes___"
	*StrPad(input, 6, "___", "RIGHT")       // produces "Codes_"
	*StrPad(input, 3, "*", "RIGHT")         // produces "Codes"
*/
func StrPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input[0:padLength]
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(length))] + input +
			strings.Repeat(padString, int(repeat))[:int(math.Ceil(length))]
	}

	return output
}
