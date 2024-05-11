package matchString

import (
	"slices"
)

type phrase struct {
	ac        Collision[[]string, []string] // ac机
	invertIdx map[string][]int              // 倒排索引
	wordGroup [][]string                    // 关键词组
}

// NewPhrase 词组匹配
func NewPhrase() Collision[[][]string, [][]string] {
	return &phrase{}
}

func (m *phrase) Build(words [][]string) error {
	var keywords = make([]string, 0)
	for _, v := range words {
		keywords = append(keywords, v...)
	}
	m.wordGroup = words
	m.ac = buildACMachine(keywords)
	m.invertIdx = buildInvertIdx(words)
	return nil
}

func (m *phrase) Scan(text string) [][]string {
	var res = make([][]string, 0)
	// 用AC机识别所有命中词
	acMatchWords := m.ac.Scan(text)
	// 过滤重复词(AC机识别结果可能会有重复词)
	matchWords := m.filterDuplicate(acMatchWords)
	// 用于存储每个命中的词的索引，对应的词（）
	var hashWords = make(map[int][]string)
	for _, word := range matchWords {
		// 识别出的每个关键词，都需要去倒排索引里面查找是否有命中
		idxList, ok := m.invertIdx[word]
		if !ok { // 这只是加个保险，理论上是不可能出现
			continue
		}
		for _, idx := range idxList {
			if _, ok = hashWords[idx]; !ok {
				hashWords[idx] = []string{word}
				continue
			}
			hashWords[idx] = append(hashWords[idx], word)
		}
	}
	// 将每个词过倒排索引，提取所有倒排索引命中的hash
	for idx, hashWord := range hashWords {
		if len(hashWord) < 2 { // 词组的数量不能小于2
			continue
		}
		// 直接比较相同idx下，长度是否一致，即可判定是否命中
		// 这里不考虑原词组中有重复词的情况由入口处保证；而新词组再开始就做了去重，这里就不可能出现重复
		if len(m.wordGroup[idx]) == len(hashWord) {
			res = append(res, m.wordGroup[idx])
		}
	}
	return res
}

// filterDuplicate 过滤重复词（双指针方法，更节约内存）
func (m *phrase) filterDuplicate(words []string) []string {
	if len(words) == 0 {
		return words
	}
	// 对字符串数组进行排序
	slices.Sort(words)
	// 双指针过滤重复词
	writeIdx := 0
	for readIdx := 1; readIdx < len(words); readIdx++ {
		if words[writeIdx] != words[readIdx] {
			writeIdx++
			words[writeIdx] = words[readIdx]
		}
	}
	return words[:writeIdx+1]
}

// buildACMachine 初始化AC自动机
func buildACMachine(keywords []string) Collision[[]string, []string] {
	m := NewAc()
	m.Build(keywords)
	return m
}

// buildInvertIdx 倒排索引
func buildInvertIdx(wordGroup [][]string) map[string][]int {
	var invertIdx = make(map[string][]int)
	for k, words := range wordGroup {
		for _, word := range words {
			if _, ok := invertIdx[word]; !ok {
				invertIdx[word] = []int{k}
			} else {
				invertIdx[word] = append(invertIdx[word], k)
			}
		}
	}
	return invertIdx
}
