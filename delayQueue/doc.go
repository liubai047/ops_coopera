// Package delayQueue 延迟队列实现
//
// 使用最小堆 + 哈希表实现。这种延迟队列的实现对事件总数量敏感（堆排序复杂度是Nlog2n），对事件延迟时间不敏感。
// 如果有大量、高频、重复性事件，可以考虑使用时间轮算法实现.
package delayQueue
