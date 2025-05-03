package matchString

// 树的节点
type node struct {
	fail  *node          //	失败指针
	isEnd bool           // 是否词组结尾
	child map[rune]*node // 子节点
	depth int 	// 深度
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
	greed bool  // 是否贪婪模式
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
				nodePtr.child[by].depth = nodePtr.depth + 1
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
		queue = queue[1:]
		// 遍历当前节点的子节点
		for word, childNode := range nowNode.child {
			// 将子节点写入队列中
			queue = append(queue, childNode)
			// 如果当前节点为root节点，则其子节点直接指向root节点
			if nowNode == a.root {
				childNode.fail = a.root
				continue
			}
			failNode := newNode.fail
			for failNode != nil {
				if failChild, ok := failNode.child[word]; ok {
					childNode.fail = failChild
					break
				}
				failNode = failNode.fail
			}
			if failNode == nil {
				childNode.fail = a.root
			}
		}
	}
}

// Scan 扫描树
func (a *acTree) Scan(text string) []string {
	var p = a.root
	var res = make([]string, 0)
	var runeText = []rune(text)
	// 每次循环一个字符
	for k, i := range runeText {
		// 循环找p的fail节点直到找到 有当前字符子节点的节点 或是 回到root节点
		for {
			if child, ok := p.child[r];ok {
				p = child
				break
			}
			if p == a.root {
				break
			}
			p = p.fail
		}
		// 检查当前节点以及沿着fail指针回溯的所有节点
		temp := p
		for temp != a.root {
			if temp.isEnd {
				start := k + 1 -temp.depth
				matched := runeText[start : k+1]
				res = append(res, string(matched))
				// 如果是贪婪模式，只记录最长匹配，然后退出
				if a.greed {
					break
				}
			}
			temp = temp.fail
		}
	}
	return res
}
