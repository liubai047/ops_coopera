package matchString

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"

	"git.woa.com/kf_cdms/go-public/helper/carr"
)

func Test_filterDuplicate(t *testing.T) {
	m := &phrase{}
	tests := []struct {
		name  string
		words []string
		want  []string
	}{
		{
			name:  "No duplicates",
			words: []string{"apple", "banana", "orange"},
			want:  []string{"apple", "banana", "orange"},
		},
		{
			name:  "With duplicates",
			words: []string{"apple", "banana", "orange", "apple", "banana"},
			want:  []string{"apple", "banana", "orange"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.filterDuplicate(tt.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 实验证明，在1W数量级下，双指针去重的性能依然好过map去重
func Benchmark_filterDuplicate(b *testing.B) {
	// m := &phrase{}
	var tests = []string{"apple", "banana", "orange", "apple", "banana", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange", "apple", "banana", "orange"}
	for i := 0; i < 10000; i++ {
		tests = append(tests, carr.NewArr(tests).RandSlice())
	}
	println(len(tests))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		(&phrase{}).filterDuplicate(tests) // 346.8 ns/op(50)   108694 ns/op(1w+)
		// carr.NewArr(tests).ArrayUnique()   // 751.7 ns/op(50)	140485 ns/op(1w+)
	}
}

func TestWgMatch_Scan(t *testing.T) {
	wordGroups := [][]string{
		{"苹果", "香蕉", "橙子"},
		{"西瓜", "葡萄"},
		{"菠萝", "橙子", "葡萄"},
	}
	m := NewPhrase()
	err := m.Build(wordGroups)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		text string
		want [][]string
	}{
		{
			name: "No match",
			text: "这是一个测试文本",
			want: [][]string{},
		},
		{
			name: "Single match",
			text: "这个文本包含苹果、香蕉和橙子",
			want: [][]string{{"苹果", "香蕉", "橙子", "0"}},
		},
		{
			name: "Multiple matches",
			text: "这个文本包含苹果、香蕉、橙子、西瓜和葡萄",
			want: [][]string{{"苹果", "香蕉", "橙子", "0"}, {"西瓜", "葡萄", "1"}},
		},
	}

	sortResults := func(results [][]string) {
		for _, result := range results {
			sort.Strings(result)
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i][0] < results[j][0]
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.Scan(tt.text)
			sortResults(got)
			sortResults(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkWgMatch_Scan(b *testing.B) {
	wordGroups := [][]string{
		{"苹果", "香蕉", "橙子"},
		{"西瓜", "葡萄"},
		{"菠萝", "橙子", "葡萄"},
		{"紫苹果莓", "椰子", "柠檬"},
		{"热情果", "莲雾"},
		{"番石榴", "青柠"},
		{"猕猴桃", "红桑子"},
		{"牛油果", "西瓜", "西柚"},
		{"大树菠萝", "山竹", "杨桃"},
		{"青苹果", "哈密瓜", "荔枝"},
		{"桃驳梨", "枇杷", "猕猴桃"},

		{"苹果", "橙子"},
		{"西瓜", "香蕉", "葡萄"},
		{"菠萝", "葡萄"},
		{"紫苹果莓", "橙子", "椰子", "柠檬"},
		{"热情果", "青柠"},
		{"番石榴", "莲雾"},
		{"猕猴桃", "西瓜", "西柚"},
		{"牛油果", "红桑子"},
		{"大树菠萝", "枇杷", "山竹"},
		{"青苹果", "杨桃", "哈密瓜", "荔枝"},
		{"桃驳梨", "猕猴桃"},

		{"苹果", "香蕉", "葡萄", "橙子"},
		{"菠萝", "橙子", "椰子", "葡萄"},
		{"柠檬", "番石榴", "莲雾"},
		{"紫苹果莓", "青柠"},
		{"猕猴桃", "西瓜", "西柚", "热情果"},
		{"哈密瓜", "荔枝", "橙子", "红桑子"},
		{"大树菠萝", "枇杷", "山竹"},
		{"青苹果", "牛油果", "杨桃", "猕猴桃"},
		{"苹果", "橙子", "桃驳梨", "橙子"},
		{"大树菠萝", "牛油果", "杨桃", "山竹"},
		{"青苹果", "枇杷", "猕猴桃"},
		{"苹果", "杨桃", "橙子", "桃驳梨"},

		{"西瓜", "香蕉", "葡萄"},
		{"热情果", "青柠", "菠萝", "葡萄"},
		{"紫苹果莓", "椰子", "柠檬"},
		{"番石榴", "西瓜", "西柚", "莲雾", "桃驳梨"},
		{"大树菠萝", "枇杷", "山竹", "牛油果", "红桑子"},
		{"青苹果", "猕猴桃", "杨桃", "哈密瓜", "荔枝"},
	}
	// println(len(wordGroups))
	m := NewPhrase()
	m.Build(wordGroups)
	text := strings.Repeat("这个文本包含苹果、香蕉、橙子、西瓜和葡萄", 150)
	// println(len(text))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Scan(text)
	}
}

func Test_Context(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "a1", "a1_v")
	println(ctx.Done() == nil) // 父节点是valueCtx,执行的是emptyCtx的Done方法，返回为nil
	ctx, cancel1 := context.WithCancel(ctx)
	defer cancel1()
	ctx = context.WithValue(ctx, "ca1", "ca1_v")
	println(ctx.Done() == nil) // 执行的是cancelCtx的Done方法，此时不为nil
	ctx2, cancel2 := context.WithCancel(ctx)
	ctx3 := context.WithValue(ctx2, "a2", "a2_v")
	println(ctx.Err() == nil, ctx2.Err() == nil) // 还未cancel，两个都为true
	cancel2()
	println(ctx.Err(), ctx2.Err().Error()) // 父cancelCtx执行cancel后，子cancel（map中）受到影响，开始全部cancel
	println(ctx3.Value("a1").(string))
	println(ctx3.Value("ca1").(string))
	println(ctx3.Value("a2").(string))
}
