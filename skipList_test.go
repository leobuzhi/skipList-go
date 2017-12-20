package skiplist

import (
	"fmt"
	"testing"
	"math/rand"
)

func (s *SkipList) print() {
	fmt.Println("header:")
	for i, link := range s.header.forward {
		if link != nil {
			fmt.Printf("\t%d: -> %v\n", i, link.key)
		} else {
			fmt.Printf("\t%d : -> END\n", i)
		}
	}

	for node := s.header.next(); node != nil; node = node.next() {
		fmt.Printf("%v: %v (level %d)\n", node.key, node.value, len(node.forward))
		for i, link := range node.forward {
			if link != nil {
				fmt.Printf("\t%d: -> %v\n", i, link.value)
			} else {
				fmt.Printf("\t%d: -> END\n", i)
			}
		}
	}
	fmt.Println()
}

func TestPrint(t *testing.T) {
	intMap := NewIntMap()

	for i := 0; i != 50; i++ {
		v := rand.Int()
		intMap.Set(v, v)
	}
	intMap.print()
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

func TestSetAgain(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)
	s.Set(1, 1)
	s.Set(2, 2)

	s.Set(0, 6)
	if value, ok := s.Get(0); !(assertEqual(6, value) && ok) {
		t.Errorf("set error")
	}

	s.Set(1, 9)
	if value, ok := s.Get(1); !(assertEqual(9, value) && ok) {
		t.Errorf("set error")
	}
}

func TestDelete(t *testing.T) {
	s := NewIntMap()
	for i := 0; i != 100; i++ {
		s.Set(i, i)
	}
	for i := 0; i != 100; i += 2 {
		s.Delete(i)
	}

	for i := 0; i != 100; i += 2 {
		if _, ok := s.Get(i); ok {
			t.Errorf("delete error")
		}
	}

	if value, ok := s.Delete(1000); value != nil || ok {
		t.Errorf("delete error")
	}
}

func TestLen(t *testing.T) {
	s := NewIntMap()

	for i := 0; i != 10; i++ {
		s.Set(i, i)
	}

	if !assertEqual(10, s.Len()) {
		t.Errorf("len error")
	}

	for i := 0; i != 5; i++ {
		s.Delete(i)
	}

	if !assertEqual(5, s.Len()) {
		t.Errorf("len error")
	}

	s.Delete(1000)

	if !assertEqual(5, s.Len()) {
		t.Errorf("len error")
	}
}

func TestIteration(t *testing.T) {
	s := NewIntMap()
	for i := 0; i != 20; i++ {
		s.Set(i, i)
	}

	seen := 0
	var lastKey int

	i := s.Iterator()
	defer i.Close()

	for i.Next() {
		seen++
		lastKey = i.Key().(int)
		if i.Key() != i.Value() {
			t.Errorf("iterator error")
		}
	}

	if seen != s.Len() {
		t.Errorf("iterator len error")
	}

	for i.Previous() {
		if i.Key() != i.Value() {
			t.Errorf("iterator wrong key value")
		}

		if i.Key().(int) >= lastKey {
			t.Errorf("iterator wrong key value")
		}

		lastKey = i.Key().(int)
	}

	if lastKey != 0 {
		t.Errorf("iterator wrong key value")
	}
}

func TestRangeIteration(t *testing.T) {
	s := NewIntMap()
	for i := 0; i != 20; i++ {
		s.Set(i, i)
	}

	max, min := 0, 1000

	var lastKey, seen int

	i := s.Range(5, 10)
	defer i.Close()

	for i.Next() {
		seen++
		lastKey = i.Key().(int)
		if lastKey > max {
			max = lastKey
		}

		if lastKey < min {
			min = lastKey
		}

		if i.Key() != i.Value() {
			t.Errorf("range iteraotr error")
		}
	}

	if seen != 5 {
		t.Errorf("range iteraotr error,error seen")
	}

	if min != 5 {
		t.Errorf("range iteraotr error,error min")
	}

	if max != 9 {
		t.Errorf("range iteraotr error,error max")
	}

	if i.Seek(4) {
		t.Errorf("range iteraotr error,invalid range")
	}

	if !i.Seek(5) {
		t.Errorf("range iteraotr error,not seek")
	}

	if i.Key().(int) != 5 || i.Value().(int) != 5 {
		t.Errorf("range iteraotr error,key value error")
	}

	if i.Seek(10) {
		t.Errorf("range iteraotr error,invalid range")
	}
}

func TestInsert(t *testing.T) {
	s := NewIntMap()
	inertions := []int{4, 3, 5, 1, 9, 2, 7}

	for _, v := range inertions {
		s.Set(v, v)
	}

	for _, v := range inertions {
		if k, v := s.Get(v); assertEqual(k, v) {
			t.Errorf("insert error")
		}
	}
}

func makeList(n int) *SkipList {
	s := NewIntMap()
	for i := 0; i != n; i++ {
		k := rand.Int()
		s.Set(k, k)
	}
	return s
}

func TestOrder(t *testing.T) {
	s := NewIntMap()

	for i := 0; i < 10000; i++ {
		r := rand.Int()
		s.Set(r, r)
	}

	last := 0

	i := s.Iterator()
	defer i.Close()

	for i.Next() {
		if last != 0 && i.Key().(int) < last {
			t.Errorf("error order")
		}
		last = i.Key().(int)
	}

	for i.Previous() {
		if i.Key().(int) > last {
			t.Errorf("error order")
		}
		last = i.Key().(int)
	}
}

func LookupBenchmark(b *testing.B) {
	b.StopTimer()
	s := makeList(1000)
	b.StartTimer()
	for i := 0; i != b.N; i++ {
		s.Get(i)
	}
}

func SetBenchmark(b *testing.B) {
	b.StopTimer()
	var values []int
	for i := 0; i != b.N; i++ {
		values = append(values, rand.Int())
	}

	s := NewIntMap()
	b.StartTimer()
	for i := 0; i != b.N; i++ {
		s.Set(values[i], values[i])
	}

}

func BenchmarkSet(b *testing.B) {
	SetBenchmark(b)
}

func BenchmarkGet(b *testing.B) {
	LookupBenchmark(b)
}
