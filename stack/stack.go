package stack

type Stack[T any] []T

func New[T any]() Stack[T] {
	return make(Stack[T], 0)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack[T]) Push(item T) {
	*s = append(*s, item)
}

func (s *Stack[T]) Pop() T {
	if s.IsEmpty() {
		var emptyResult T
		return emptyResult
	}
	index := len(*s) - 1
	pop := (*s)[index]
	*s = (*s)[:index]
	return pop

}
