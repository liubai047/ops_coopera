package structure

// HeapArea 堆。非并发安全类型，并发场景下需自己加锁二次封装
type HeapArea[T comparable] struct {
	data []T
	t    bool                          // true为大顶堆，false为小顶堆
	less func(data []T, i, j int) bool // 比较两个元素谁更小，当i <j时返回true,否则返回false
}

// NewHeapArea 初始化堆数据结构，默认最小堆
func NewHeapArea[T comparable](t bool, less func(data []T, i, j int) bool) *HeapArea[T] {
	return &HeapArea[T]{
		data: make([]T, 0),
		t:    t,
		less: less,
	}
}

// CreateFromSlice 将数组初始化为堆(自顶向下建堆，复杂度为O(log2n))
func (h *HeapArea[T]) CreateFromSlice(d []T, t bool, less func(data []T, i, j int) bool) *HeapArea[T] {
	h.data = d
	h.t = t
	h.less = less
	for i := fatherChild(len(h.data) - 1); i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

// Push 新增。返回新增元素的索引下标
func (h *HeapArea[T]) Push(val T) int {
	h.data = append(h.data, val)
	// 从底至顶堆化
	return h.siftUp(len(h.data) - 1)
}

// Get 根据索引获取元素
func (h *HeapArea[T]) Get(idx int) T {
	return h.data[idx]
}

// Delete 根据索引删除堆中的元素
func (h *HeapArea[T]) Delete(idx int) {
	if idx >= len(h.data) {
		return
	}
	// 交换堆顶和堆尾元素
	h.swap(idx, len(h.data)-1)
	// 删除堆顶元素
	h.data = h.data[:len(h.data)-1]
	// 堆重新排序
	h.siftDown(idx)
}

// 重新堆化 - 从底至顶堆化
func (h *HeapArea[T]) siftUp(i int) int {
	for {
		// 获取节点的父节点
		p := fatherChild(i)
		// 越过根节点 或 节点大小比较不匹配，则结束
		if h.t { // 大顶堆
			if p < 0 || h.less(h.data, i, p) {
				break
			}
		} else { // 小顶端
			if p < 0 || h.less(h.data, p, i) {
				break
			}
		}
		// 交换两节点
		h.swap(i, p)
		// 循环向上堆化
		i = p
	}
	return i
}

// 重新堆化 - 从顶至底堆化
func (h *HeapArea[T]) siftDown(i int) {
	for {
		// 先获取左右子节点
		lIdx := lChild(i)
		rIdx := rChild(i)
		// 如果左节点索引不合法，则直接退出（左节点不合法则右节点必定也不合法）
		if len(h.data) <= lIdx {
			break
		}
		// 拿当前节点的值和左右节点的值进行比较
		if h.t { // 如果很大顶堆
			maxIdx := lIdx
			if rIdx < len(h.data) && h.less(h.data, lIdx, rIdx) {
				maxIdx = rIdx
			}
			if h.less(h.data, maxIdx, i) {
				break
			}
			h.swap(i, maxIdx)
			i = maxIdx
		} else { // 小顶堆
			// 假设最小值为左节点的值
			minIdx := lIdx
			// 如果右节点合法，且右节点的值小于左节点的值，则更新最小值为右节点的值
			if rIdx < len(h.data) && h.less(h.data, rIdx, lIdx) {
				minIdx = rIdx
			}
			// 如果当前索引值比最小值都小，则退出循环
			if h.less(h.data, i, minIdx) {
				break
			}
			// 交换节点值并更新索引
			h.swap(i, minIdx)
			i = minIdx
		}
	}
}

// 交换元素
func (h *HeapArea[T]) swap(o, n int) {
	h.data[o], h.data[n] = h.data[n], h.data[o]
}

// Pop 弹出堆顶元素
func (h *HeapArea[T]) Pop() T {
	// 交换堆顶和堆尾元素
	h.swap(0, len(h.data)-1)
	// 获取堆顶元素
	v := h.data[len(h.data)-1]
	// 删除堆顶元素
	h.data = h.data[:len(h.data)-1]
	// 堆重新排序
	h.siftDown(0)
	// 返回
	return v
}

// Search 遍历数组，寻找
func (h *HeapArea[T]) Search(val T) int {
	for k, v := range h.data {
		if v == val {
			return k
		}
	}
	return -1
}

// Peek 返回堆顶元素
func (h *HeapArea[T]) Peek() T {
	if len(h.data) == 0 {
		return *new(T)
	}
	return h.data[0]
}

// 左孩子节点
func lChild(index int) int {
	return 2*index + 1
}

// 右孩子节点
func rChild(index int) int {
	return 2*index + 2
}

// 获取父节点
func fatherChild(index int) int {
	v := (index - 1) / 2
	if v < 0 {
		return 0
	}
	return v
}
