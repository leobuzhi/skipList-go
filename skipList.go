package skiplist

type node struct {
	forward    []*node
	backword   *node
	key, value interface{}
}

func (n *node) next() *node {
	if len(n.forward) == 0 {
		return nil
	}
	return n.forward[0]
}

func (n *node) previous() *node {
	return n.backword
}

func (n *node) hasNext() bool {
	return n.next() != nil
}

func (n *node) hasPrevious() bool {
	return n.backword != nil
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
	Next() (ok bool)
	Previous() (ok bool)
	Key() interface{}
	Value() interface{}
	Seek(key interface{}) (ok bool)
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

}

func (i *iter) Previous() (ok bool) {

}

func (i *iter) Seek(key interface{}) (ok bool) {

}

func (i *iter) Close() {

}
