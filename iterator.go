package skiplist

type SkiplistIterator struct {
	node *SkipListNode
}

func newIterator(node *SkipListNode) *SkiplistIterator {
	return &SkiplistIterator{node: node}
}

func (iter *SkiplistIterator) Next() bool {
	if iter.node == nil {
		return false
	}

	iter.node = iter.node.level[0].forward
	return iter.node != nil
}

func (iter *SkiplistIterator) Value() int {
	return iter.node.val
}

func (iter *SkiplistIterator) Object() interface{} {
	return iter.node.obj
}
