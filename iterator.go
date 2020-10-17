package skiplist

type SkiplistIterator struct {
	node *SkipListNode
	rank uint // 从1开始
}

func newIterator(node *SkipListNode, rank uint) *SkiplistIterator {
	return &SkiplistIterator{node: node, rank: rank}
}

func (iter *SkiplistIterator) Next() bool {
	if iter.node == nil {
		return false
	}

	iter.rank++ // 等价于 iter.rank += iter.node.level[0].span
	iter.node = iter.node.level[0].forward
	return iter.node != nil
}

func (iter *SkiplistIterator) Rank() uint {
	return iter.rank
}

func (iter *SkiplistIterator) Elem() SkipListElem {
	return iter.node.elem
}
