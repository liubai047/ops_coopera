package matchString

import (
	"reflect"
	"testing"
)

func Test_KMP(t *testing.T) {
	word := "自动机"

	nKmp := NewKmp()
	nKmp.Build(word)

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Test1",
			text:     "这是一个测试代码，用于测试AC自动机是否可以正确匹配中文字符。",
			expected: 15,
		},
		{
			name:     "Test2",
			text:     "这段代码没有包含任何关键词。",
			expected: -1,
		},
		{
			name:     "Test3",
			text:     "测试自动机中文匹配代码",
			expected: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := nKmp.Scan(test.text)
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, res)
			}
		})
	}
}
