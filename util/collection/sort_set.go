package collection

import "sync"

type SortSet[T comparable] struct {
	data   map[T]struct{}
	tp     int
	elem   []T
	locker sync.Mutex
}

func NewSet[T comparable]() *SortSet[T] {
	return &SortSet[T]{
		data: make(map[T]struct{}),
	}
}
func (s *SortSet[T]) Add(e T) {
	s.locker.Lock()
	defer s.locker.Unlock()

	if _, ok := s.data[e]; !ok {
		s.elem = append(s.elem, e)
		s.data[e] = struct{}{}
		s.tp++
	}
}

func (s *SortSet[T]) Count() int {
	return s.tp
}

func (s *SortSet[T]) Elems() []T {
	return s.elem[:]
}

func (s *SortSet[T]) Contains(e T) bool {
	s.locker.Lock()
	defer s.locker.Unlock()
	_, ok := s.data[e]
	return ok

}
