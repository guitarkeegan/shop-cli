package shopifyv2

type Set[C comparable] struct {
	data map[C]struct{}
}

func NewSet[C comparable](items ...C) *Set[C] {
	set := Set[C]{
		data: map[C]struct{}{},
	}
	for _, item := range items {
		set.data[item] = struct{}{}
	}
	return &set
}

func (s *Set[C]) Has(c C) bool {
	_, has := s.data[c]
	return has
}

func (s *Set[C]) Remove(c C) {
	delete(s.data, c)
}

func (s *Set[C]) Add(c C) {
	s.data[c] = struct{}{}
}
