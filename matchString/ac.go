package matchString

// 树的节点
type node struct {
	fail  *node          //	失败指针
	isEnd bool           // 是否词组结尾
	child map[rune]*node // 子节点
}

// 初始化一个节点
func newNode() *node {
	return &node{
		fail:  nil,
		isEnd: false,
		child: make(map[rune]*node),
	}
}

// ac自动机树
type acTree struct {
	root *node // root节点
}

// NewAc AC自动机，词匹配
func NewAc() Collision[[]string, []string] {
	return &acTree{root: newNode()}
}

// Build 构建树
func (a *acTree) Build(words []string) error {
	for _, word := range words {
		// 当前扫描树的指针，每次插入词从根节点开始扫描
		var nodePtr = a.root
		// 将词拆为单个字符循环
		for _, by := range []rune(word) {
			// 判断该字符是否存在于树中，不存在则添加到树中
			if _, ok := nodePtr.child[by]; !ok {
				nodePtr.child[by] = newNode()
			}
			// 将扫描指针移动到当前字符节点的子节点
			nodePtr = nodePtr.child[by]
		}
		// 循环完毕一个词之后，nodePrt指针指向的是最后一个字符的位置，将其词尾标记置为true
		nodePtr.isEnd = true
	}
	// 构建fail指针
	a.BuildFail()
	return nil
}

// BuildFail 构建树的fail指针
func (a *acTree) BuildFail() {
	// 开始广度遍历树
	var queue = make([]*node, 0)
	queue = append(queue, a.root)
	for len(queue) > 0 {
		var nowNode = queue[0]
		// 弹出第一个字符
		queue = append(queue[:0], queue[1:]...)
		// 遍历当前节点的子节点
		for word, childNode := range nowNode.child {
			// 将子节点写入队列中
			queue = append(queue, childNode)
			// 如果当前节点为root节点，则其子节点直接指向root节点
			if nowNode == a.root {
				childNode.fail = a.root
				continue
			}
			// 如果当前节点的fail指针指向的节点的子节点存在该字符,则当前节点的fail指针指向该节点
			if failNode, ok := nowNode.fail.child[word]; ok {
				childNode.fail = failNode
			} else {
				// 否则将当前节点的fail指针指向父节点的fail指针
				childNode.fail = nowNode.fail
			}
		}
	}
}

// Scan 扫描树
func (a *acTree) Scan(text string) []string {
	var p = a.root
	var res = make([]string, 0)
	var sIdx int
	var runeText = []rune(text)
	// 每次循环一个字符
	for k, i := range runeText {
		// 循环找p的fail节点直到找到 子节点相同 或是 回到root节点
		// 如果找到root节点，表示该字符不存在，则跳过该字符，开循环下个字符
		_, ok := p.child[i]
		for !ok && p != a.root {
			p = p.fail
		}
		// 进到这里只存在两个情况，p指向root节点，或者当前节点与子节点相同（上面跳出循环的关键）
		if _, ok := p.child[i]; ok {
			// 如果p指向root节点，说明是新词第一个字符匹配的开始
			if p == a.root {
				sIdx = k
			}
			// 否则说明是已匹配的后续词
			p = p.child[i] // 子节点递归比较，直到回归到root节点(表示没匹配到)或是顺利匹配到节点末尾
			if p.isEnd {
				end := k
				runes := runeText[sIdx : end+1]
				res = append(res, string(runes))
			}
		}
	}
	return res
}
