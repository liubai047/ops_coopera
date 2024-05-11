package delayQueue

import (
	"time"
)

// MessageItem 每条消息的结构
type MessageItem[T any] struct {
	key     string    // 消息key标识,唯一不可重复
	content T         // 消息内容
	sec     time.Time // 过期时间
}

func NewMessageItem[T any](key string, content T, sec time.Time) *MessageItem[T] {
	return &MessageItem[T]{
		key:     key,
		content: content,
		sec:     sec,
	}
}

func (m *MessageItem[T]) Content() T {
	return m.content
}
