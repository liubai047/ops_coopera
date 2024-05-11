package delayQueue

import (
	"context"
	"sync"
	"time"

	"git.woa.com/kf_cdms/go-public/structure"
)

// 消息体类型
type messageTyps any

// HeapDelayQueue 延迟队列实现(最小堆)
type HeapDelayQueue[T messageTyps] struct {
	m     map[string]*MessageItem[T]           // 用于消息去重/查询，消息key必须是string类型，且唯一
	heap  *structure.HeapArea[*MessageItem[T]] // 最小堆
	lock  sync.RWMutex                         // 加把锁
	timer *time.Timer                          // 最近一个任务的timer
}

// NewDelayQe 初始化延迟队列
func NewDelayQe[T messageTyps]() *HeapDelayQueue[T] {
	tm := time.NewTimer(time.Second)
	tm.Stop() // 队列创建时无任务，所以将计时器置为stop状态
	return &HeapDelayQueue[T]{
		m: make(map[string]*MessageItem[T]),
		heap: structure.NewHeapArea[*MessageItem[T]](false, func(data []*MessageItem[T], i, j int) bool {
			if i == j {
				return true
			}
			return data[i].sec.Before(data[j].sec)
		}),
		lock:  sync.RWMutex{},
		timer: tm,
	}
}

// Add 添加任务
//
// @param key 消息id，需保持唯一
// @param msg 消息体结构
func (hd *HeapDelayQueue[T]) Add(msg *MessageItem[T]) error {
	hd.lock.Lock()
	defer hd.lock.Unlock()
	// 更新全局定时器 第一个任务 或者 新加入的任务执行时间比最快执行的任务时间还要早
	if len(hd.m) < 1 || msg.sec.Before(hd.heap.Get(0).sec) {
		hd.resetTimerWithDelay(msg.sec.Sub(time.Now()))
	}
	if _, ok := hd.m[msg.key]; ok {
		return RepeatError
	}
	hd.m[msg.key] = msg
	hd.heap.Push(msg)
	return nil
}

// Peek 提取最近的一个任务详情，但并不出栈
func (hd *HeapDelayQueue[T]) Peek() *MessageItem[T] {
	hd.lock.Lock()
	defer hd.lock.Unlock()
	val := hd.heap.Peek()
	return val
}

// Search 查询任务
func (hd *HeapDelayQueue[T]) Search(key string) (*MessageItem[T], error) {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	return hd.m[key], nil
}

// Delete 移除任务
func (hd *HeapDelayQueue[T]) Delete(key string) error {
	hd.lock.Lock()
	defer hd.lock.Unlock()
	return hd.delete(key)
}

func (hd *HeapDelayQueue[T]) delete(key string) error {
	idx := hd.heap.Search(hd.m[key])
	if idx < 0 {
		return KeyError
	}
	hd.heap.Delete(idx)
	delete(hd.m, key)
	// 移除前需要判断全局定时器是否会受到影响
	if idx == 0 {
		if len(hd.m) == 0 { // 当任务集为空时，停止任务
			hd.timer.Stop()
		} else { // 当还存在任务集时，重置到下一个任务
			hd.resetTimerWithDelay(hd.heap.Get(0).sec.Sub(time.Now()))
		}
	}
	return nil
}

// Watch 监听延迟队列，该方法会阻塞，直到延迟队列最早的事件触发
//
// @param ctx 通过ctx.Done控制读取延迟的阻塞停止
func (hd *HeapDelayQueue[T]) Watch(ctx context.Context) (*MessageItem[T], error) {
	select {
	case <-ctx.Done():
		return nil, CtxDoneError
	case <-hd.timer.C:
		hd.lock.Lock() // 防止同时触发Delete,出现幻读
		defer hd.lock.Unlock()
		if len(hd.m) == 0 {
			return nil, EmptyQueue
		}
		messageItem := hd.heap.Get(0)
		if err := hd.delete(messageItem.key); err != nil {
			return nil, err
		}
		if len(hd.m) > 0 {
			hd.resetTimerWithDelay(hd.heap.Get(0).sec.Sub(time.Now()))
		}
		return messageItem, nil
	}
}

// 清空通道的值
func (hd *HeapDelayQueue[T]) resetTimerWithDelay(duration time.Duration) {
	// 对于已经关闭的timer,检查是否有尚未消费的timer，有的话直接移除
	if !hd.timer.Stop() {
		select {
		case <-hd.timer.C:
		default:
		}
	}
	hd.timer.Reset(duration)
}
