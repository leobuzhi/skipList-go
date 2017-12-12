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