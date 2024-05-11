package structure

import (
	"container/list"
	"errors"
)

type Stack[T any] struct {
	*list.List
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		list.New(),
	}
}

func (s *Stack[T]) Len() int {
	return s.Len()
}

func (s *Stack[T]) Push(v T) {
	s.PushBack(v)
}

func (s *Stack[T]) Pop() (v T, err error) {
	backItem := s.Back()
	if backItem != nil {
		return s.Remove(backItem).(T), nil
	}
	return v, errors.New("cannot pop in empty stack")
}
