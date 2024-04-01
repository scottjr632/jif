package stack

import (
	"fmt"
)

type Stack[T any] []T

func New[T any](val ...T) Stack[T] {
	s := make(Stack[T], 0)
	s = append(s, val...)
	return s
}

func (s *Stack[T]) PushMany(vals []T) {
	*s = append(*s, vals...)
}

func (s *Stack[T]) PopMany(count int) ([]T, error) {
	if len(*s) < count {
		return nil, fmt.Errorf("Stack is empty")
	}

	res := (*s)[len(*s)-count:]
	*s = (*s)[:len(*s)-count]
	return res, nil
}

func (s *Stack[T]) Peek() (*T, error) {
	if len(*s) == 0 {
		return nil, fmt.Errorf("Stack is empty")
	}
	return &(*s)[len(*s)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Pop() (*T, error) {
	if len(*s) == 0 {
		return nil, fmt.Errorf("Stack is empty")
	}

	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return &res, nil
}
