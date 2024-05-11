package matchString

import (
	"reflect"
	"strings"
	"testing"

	"git.woa.com/kf_cdms/go-public/helper/carr"
)

func TestAcAutomaton(t *testing.T) {
	words := []string{"测试", "自动机", "中文", "匹配", "代码"}

	acMachine := NewAc()
	acMachine.Build(words)

	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "Test1",
			text:     "这是一个测试代码，用于测试AC自动机是否可以正确匹配中文字符。",
			expected: []string{"测试", "代码", "测试", "自动机", "匹配", "中文"},
		},
		{
			name:     "Test2",
			text:     "这段代码没有包含任何关键词。",
			expected: []string{"代码"},
		},
		{
			name:     "Test3",
			text:     "测试自动机中文匹配代码",
			expected: []string{"测试", "自动机", "中文", "匹配", "代码"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := acMachine.Scan(test.text)
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, res)
			}
		})
	}
}

func BenchmarkAcAutomaton(b *testing.B) {
	words := []string{"测试", "自动机", "中文", "匹配", "代码"}
	for i := 0; i < 10000; i++ {
		words = append(words, carr.NewArr(words).RandSlice())
	}
	acMachine := NewAc()

	b.Run("BuildTree", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			acMachine.Build(words)
		}
	})

	text := strings.Repeat("这是一个测试代码，用于测试AC自动机是否可以正确匹配中文字符。", 1000)

	b.Run("Scan", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			acMachine.Scan(text)
		}
	})
}
