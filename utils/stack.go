package engine

type Stack[T any] []T

func New[T any](val ...T) Stack[T] {
	s := make(Stack[T], 0)
	s = append(s, val...)
	return s
}

// Push adds an element to the stack
func (s *Stack[T]) push(v T) {
	*s = append(*s, v)
}

// Pop removes and returns the top element of the stack
func (s *Stack[T]) pop() (int, error) {
	// Check if the stack is empty
	if len(*s) == 0 {
		return 0, errors.New("Stack is empty")
	}

	// Get the last element
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res, nil
}
