package stack

import "testing"

func TestStack(t *testing.T) {
	s := New[int]()
	if !s.IsEmpty() {
		t.Fatal("new stack should be empty")
	}
	if got := s.Pop(); got != 0 {
		t.Fatalf("pop on empty: got %d, want zero value", got)
	}
	for i := 1; i <= 3; i++ {
		s.Push(i)
	}
	if s.IsEmpty() {
		t.Fatal("stack should not be empty after push")
	}
	for want := 3; want >= 1; want-- {
		if got := s.Pop(); got != want {
			t.Fatalf("pop: got %d, want %d", got, want)
		}
	}
	if !s.IsEmpty() {
		t.Fatal("stack should be empty after draining")
	}
}
