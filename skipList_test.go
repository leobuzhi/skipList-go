package skiplist

import (
	"fmt"
	"testing"
)

func (s *SkipList) print() {
	fmt.Println("header:")
	for i, link := range s.header.forward {
		if link != nil {
			fmt.Println("\t%d : -> %v\n ", i, link.key)
		} else {
			fmt.Println("\t%d : -> END\n", i)
		}
	}

	for node := s.header.next(); node != nil; node = node.next() {
		fmt.Println("%v: %v (level %d)", node.key, node.value, len(node.forward))
		for i, link := range node.forward {
			if link != nil {
				fmt.Println("\t%d: -> %v", i, link.value)
			} else {
				fmt.Println("\t%d: -> END", i)
			}
		}
	}
}

func TestInit(t *testing.T) {
	s := NewCustomMap(func(left, right interface{}) bool {
		return left.(int) < right.(int)
	})

	if !s.lessThan(1, 2) {
		t.Errorf("lessThan error")
	}
}

func TestEmptyNodeNext(t *testing.T) {
	n := new(node)
	if next := n.next(); next != nil {
		t.Errorf("next error")
	}

	if n.hasNext() {
		t.Errorf("hasNext error")
	}
}

func TestEmptyNodePrevious(t *testing.T) {
	n := new(node)
	if prev := n.previous(); prev != nil {
		t.Errorf("previous error")
	}

	if n.hasPrevious() {
		t.Errorf("hasPrevious error")
	}
}

func TestNodeHasNext(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)
	node := s.header.next()
	if node.key != 0 {
		t.Errorf("get key error ")
	}

	if node.hasNext() {
		t.Errorf("hasNext error")
	}
}

func TestNodeHasPrevious(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)
	node := s.header.previous()
	if node != nil {
		t.Errorf("hasPrevious error")
	}
}

func assertEqual(left, right interface{}) bool {
	if left == right {
		return true
	}
	return false
}

func TestGet(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)

	if value, ok := s.Get(0); !(ok && assertEqual(0, value) ) {
		t.Errorf("get error")
	}

	if value, ok := s.Get(1); ok || !assertEqual(nil, value) {
		t.Errorf("get error")
	}
}

func TestGetGreatOrEqual(t *testing.T) {
	s := NewIntMap()

	if _, value, ok := s.GetGreaterOrEqual(5); !(!ok && assertEqual(value, nil)) {
		t.Errorf("GetGreaterOrEqual error")
	}

	s.Set(0, 0)
	if _, value, ok := s.GetGreaterOrEqual(5); !(!ok && assertEqual(value, nil)) {
		t.Errorf("GetGreaterOrEqual error")
	}

	s.Set(100, 100)
	if _, value, ok := s.GetGreaterOrEqual(5); !(ok && assertEqual(value, 100)) {
		t.Errorf("GetGreaterOrEqual error")
	}
}

func TestSet(t *testing.T) {
	s := NewIntMap()
	if !assertEqual(0, s.Len()) {
		t.Errorf("len is not 0 ,set error")
	}

	s.Set(0, 0)
	s.Set(1, 1)
	if !assertEqual(2, s.Len()) {
		t.Errorf("len is not 0 ,set error")
	}

	if value, ok := s.Get(0); !(assertEqual(0, value) && ok) {
		t.Errorf("set error")
	}

	if value, ok := s.Get(1); !(assertEqual(1, value) && ok) {
		t.Errorf("set error")
	}
}
