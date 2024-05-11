package matchString

import "unicode/utf8"

type kmp struct {
	next    []int
	pattern string
}

var _ Collision[string, int] = (*kmp)(nil)

// NewKmp 单个词匹配
func NewKmp() Collision[string, int] {
	return &kmp{}
}

// Build 构建next数组，即获取每一位的最大公共前后缀长度
// param s 匹配串
// return next数组
func (k *kmp) Build(s string) error {
	k.pattern = s
	// 初始化next数组,每一个下标对应一个字符
	var next = make([]int, utf8.RuneCountInString(s), utf8.RuneCountInString(s))
	// 定义第一个位置的字符的next数值为0，因为第一位不存在公共前后缀
	next[0] = 0
	// j表示与之对应的匹配串字符
	j := 0
	// i表示扫描到匹配串的第几个字符，位置从1开始
	for i := 1; i < utf8.RuneCountInString(s); i++ {
		// 如果j和i指针指向的字符相同，则i所在位置对应的next指针比i-1的值+1
		if s[j] == s[i] {
			next[i] = next[i-1] + 1
			j++ // 当发生匹配时候，i和j指针共同向后扫描
		} else {
			for {
				// 如果j到达了第一个字符，表示没有与之匹配的字符，则next数组对应位置0，同时j指针保持不动，直到匹配到对应字符
				if j == 0 {
					next[i] = 0
					break
				}
				// 这是核心，直接通过next下标跳转（其实直接向前递进也可以，只是会多算几次）
				// 因为如果前一个下标大于0，表示前面是有匹配的字符，我们需要向前递减一个值，确认与当前字符是否有匹配
				// 如果前面字符next下标为0，表示没有匹配字符，则我们直接跳转到对应字符处即可
				j = next[j-1]
				// 这里逻辑与上面if逻辑相同
				if s[j] == s[i] {
					next[i] = next[i-1] + 1
					j++
					break
				}
			}
		}
	}
	k.next = next
	return nil
}

// Scan 匹配字符串
// param s 原串
// param p 匹配串
func (k *kmp) Scan(s string) int {
	if len(k.next) == 0 {
		return -1
	}
	// 获取字符串长度，这里要兼容中文
	var sRune = []rune(s)
	var pRune = []rune(k.pattern)
	var sLen = len(sRune)
	var pLen = len(pRune)
	// 匹配串比原串还大，就别匹配了
	if pLen > sLen {
		return -1
	}
	// i和j分别表示原串和匹配串的指针
	var i, j int
	for i < sLen && j < pLen {
		// 如果当前指针指向字符相同，则直接平滑后移
		if sRune[i] == pRune[j] {
			i++
			j++
			continue
		}
		// 找到第一个j和i相同的字符
		if j == 0 {
			i++
			continue
		}
		// 如果不相等，根据next指针，跳转到对应位置
		j = k.next[j-1]
	}
	// 如果匹配成功，j的值一定等于pLen-1
	if j == pLen {
		return i - j
	}
	return -1
}
