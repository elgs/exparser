package exparser

type Lifo struct {
	top  *Element
	size int
}

type Element struct {
	value interface{}
	next  *Element
}

func (s *Lifo) Len() int {
	return s.size
}

func (s *Lifo) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

func (s *Lifo) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

func (s *Lifo) Peep() (value interface{}) {
	if s.size > 0 {
		value = s.top.value
		return
	}
	return nil
}
