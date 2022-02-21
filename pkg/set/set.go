package set

type Set map[interface{}]struct{}

func New() *Set {
	return &Set{}
}

func (s *Set) Has(e interface{}) bool {
	_, exist := (*s)[e]
	return exist
}

func (s *Set) Add(e interface{}) bool {
	if s.Has(e) {
		return false
	}
	(*s)[e] = struct{}{}
	return true
}

func (s *Set) Remove(e interface{}) bool {
	if s.Has(e) {
		return false
	}
	delete(*s, e)
	return true
}

func (s *Set) Union(t *Set) *Set {
	union := s.Copy()
	for ele := range *t {
		union.Add(ele)
	}

	return union
}

func (s *Set) Diff(t *Set) *Set {
	diff := s.Copy()
	for ele := range *t {
		diff.Remove(ele)
	}

	return diff
}

func (s *Set) Copy() *Set {
	c := New()
	for ele := range *s {
		(*c)[ele] = struct{}{}
	}

	return c
}

func (s *Set) ToSlice() []interface{} {
	res := make([]interface{}, len(*s))
	for ele := range *s {
		res = append(res, ele)
	}

	return res
}
