package skiplist

import (
	"math/rand"
)

const (
	p               = 0.25
	DefaultMaxLevel = 32
)

type node struct {
	forward    []*node
	backward   *node
	key, value interface{}
}

func (n *node) next() *node {
	if len(n.forward) == 0 {
		return nil
	}
	return n.forward[0]
}

func (n *node) previous() *node {
	return n.backward
}

func (n *node) hasNext() bool {
	return n.next() != nil
}

func (n *node) hasPrevious() bool {
	return n.previous() != nil
}

type SkipList struct {
	lessThan func(l, r interface{}) bool
	header   *node
	footer   *node
	length   int
	MaxLevel int
}

func (s *SkipList) Len() int {
	return s.length
}

type Iterator interface {
	//如果可向后迭代则指向下一个节点，并返回true，否则返回false
	Next() (ok bool)
	//如果可向前迭代则指向上一个节点，并返回true，否则返回false
	Previous() (ok bool)
	//返回当前节点的key
	Key() interface{}
	//返回当前节点的value
	Value() interface{}

	Seek(key interface{}) (ok bool)
	//关闭迭代器并回收资源，不是必要的，但对GC有好处
	Close()
}

type iter struct {
	current *node
	list    *SkipList
	key     interface{}
	value   interface{}
}

func (i iter) Key() interface{} {
	return i.key
}

func (i iter) Value() interface{} {
	return i.value
}

func (i *iter) Next() (ok bool) {
	if !i.current.hasNext() {
		return false
	}

	i.current = i.current.next()
	i.key = i.current.key
	i.value = i.current.value

	return true
}

func (i *iter) Previous() (ok bool) {
	if !i.current.hasPrevious() {
		return false
	}

	i.current = i.current.previous()
	i.key = i.current.key
	i.value = i.current.value

	return true
}

func (i *iter) Seek(key interface{}) (ok bool) {
	current := i.current
	list := i.list

	if current == nil {
		current = list.header
	}

	if current.key != nil && list.lessThan(key, current.key) {
		current = list.header
	}

	if current.backward == nil {
		current = list.header
	} else {
		current = current.backward
	}

	current = list.getPath(current, nil, key)

	if current == nil {
		return false
	}

	i.current = current
	i.key = current.key
	i.value = current.value

	return true
}

func (i *iter) Close() {
	i.current = nil
	i.list = nil
	i.key = nil
	i.value = nil
}

func (s *SkipList) getPath(current *node, update []*node, key interface{}) *node {
	depth := len(current.forward) - 1

	for i := depth; i >= 0; i-- {
		for current.forward[i] != nil && s.lessThan(current.forward[i].key, key) {
			current = current.forward[i]
		}
		if update != nil {
			update[i] = current
		}
	}

	return current.next()
}

func (s *SkipList) Iterator() Iterator {
	return &iter{
		current: s.header,
		list:    s,
	}
}

func (s *SkipList) Seek(key interface{}) Iterator {
	current := s.getPath(s.header, nil, key)
	if current == nil {
		return nil
	}
	return &iter{
		current: current,
		list:    s,
		key:     current.key,
		value:   current.value,
	}
}

func (s *SkipList) SeekToFirst() Iterator {
	if s.length == 0 {
		return nil
	}
	current := s.header.next()
	return &iter{
		current: current,
		list:    s,
		key:     current.key,
		value:   current.value,
	}
}

func (s *SkipList) SeekToLast() Iterator {
	current := s.footer
	if current == nil {
		return nil
	}
	return &iter{
		current: current,
		list:    s,
		key:     current.key,
		value:   current.value,
	}
}

func (s *SkipList) level() int {
	return len(s.header.forward) - 1
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (s *SkipList) effectiveMaxLevel() int {
	return maxInt(s.level(), s.MaxLevel)
}

func (s SkipList) randomLevel() (n int) {
	for n = 0; n < s.effectiveMaxLevel() && rand.Float64() < p; n++ {
	}
	return n
}

func (s *SkipList) GetGreaterOrEqual(min interface{}) (actualKey, value interface{}, ok bool) {
	candidate := s.getPath(s.header, nil, min)

	if candidate == nil {
		return nil, nil, false
	}
	return candidate.key, candidate.value, true
}

func (s *SkipList) Get(key interface{}) (value interface{}, ok bool) {
	candidate := s.getPath(s.header, nil, key)
	if candidate == nil || candidate.key != key {
		return nil, false
	}
	return candidate.value, true
}

func (s *SkipList) Set(key, value interface{}) {
	if key == nil {
		panic("nil key is not supported")
	}

	update := make([]*node, s.level()+1, s.effectiveMaxLevel()+1)
	candidate := s.getPath(s.header, update, key)

	if candidate != nil && candidate.key == key {
		candidate.value = value
		return
	}

	newLevel := s.randomLevel()

	if currentLevel := s.level(); newLevel > currentLevel {
		for i := currentLevel + 1; i <= newLevel; i++ {
			update = append(update, s.header)
			s.header.forward = append(s.header.forward, nil)
		}
	}

	newNode := &node{
		forward: make([]*node, newLevel+1, s.effectiveMaxLevel()+1),
		key:     key,
		value:   value,
	}

	if previous := update[0]; previous.key != nil {
		newNode.backward = previous
	}

	for i := 0; i <= newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	s.length++

	if newNode.forward[0] != nil {
		if newNode.forward[0].backward != newNode {
			newNode.forward[0].backward = newNode
		}
	}

	if s.footer == nil || s.lessThan(s.footer.key, key) {
		s.footer = newNode
	}
}

func (s *SkipList) Delete(key interface{}) (value interface{}, ok bool) {
	if key == nil {
		panic("nil key is not supported")
	}

	update := make([]*node, s.level()+1, s.effectiveMaxLevel())
	candidate := s.getPath(s.header, update, key)
	if candidate == nil || candidate.key != key {
		return nil, false
	}

	previous := candidate.backward
	if s.footer == candidate {
		s.footer = previous
	}

	next := candidate.next()
	if next != nil {
		next.backward = previous
	}

	for i := 0; i <= s.level() && update[i].forward[i] == candidate; i++ {
		update[i].forward[i] = candidate.forward[i]
	}

	for s.level() > 0 && s.header.forward[s.level()] == nil {
		s.header.forward = s.header.forward[:s.level()]
	}
	s.length--

	return candidate.value, true
}

func NewCustomMap(lessthan func(l, r interface{}) bool) *SkipList {
	return &SkipList{
		lessThan: lessthan,
		header: &node{
			forward: []*node{nil},
		},
		MaxLevel: DefaultMaxLevel,
	}
}

type Ordered interface {
	LessThan(other Ordered) bool
}

func New() *SkipList {
	comparator := func(left, right interface{}) bool {
		return left.(Ordered).LessThan(right.(Ordered))
	}
	return NewCustomMap(comparator)
}

func NewIntMap() *SkipList {
	return NewCustomMap(func(left, right interface{}) bool {
		return left.(int) < right.(int)
	})
}

func MewStringMap() *SkipList {
	return NewCustomMap(func(left, right interface{}) bool {
		return left.(string) < right.(string)
	})
}

type rangeIterator struct {
	iter
	upperLimit interface{}
	lowerLimit interface{}
}

func (i *rangeIterator) Next() bool {
	if !i.current.hasNext() {
		return false
	}

	next := i.current.next()

	if !i.list.lessThan(next.key, i.upperLimit) {
		return false
	}

	i.current = i.current.next()
	i.key = i.current.key
	i.value = i.current.value

	return true
}

func (i *rangeIterator) Previous() bool {
	if !i.current.hasPrevious() {
		return false
	}

	previous := i.current.previous()
	if i.list.lessThan(previous.key, i.lowerLimit) {
		return false
	}

	i.current = i.current.previous()
	i.key = i.current.key
	i.value = i.current.value

	return true
}

func (i *rangeIterator) Seek(key interface{}) (ok bool) {
	if i.list.lessThan(key, i.lowerLimit) || !i.list.lessThan(key, i.upperLimit) {
		return false
	}
	return i.iter.Seek(key)
}

func (i *rangeIterator) Close() {
	i.iter.Close()
	i.upperLimit = nil
	i.lowerLimit = nil
}

func (s *SkipList) Range(from, to interface{}) Iterator {
	start := s.getPath(s.header, nil, from)
	return &rangeIterator{
		iter: iter{
			current: &node{
				forward:  []*node{start},
				backward: start,
			},
			list: s,
		},
		upperLimit: to,
		lowerLimit: from,
	}
}
