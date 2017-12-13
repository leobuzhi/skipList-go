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
	i.value = i.value
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

	if current.backword == nil {
		current = list.header
	} else {
		current = current.backword
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
	return nil
}
