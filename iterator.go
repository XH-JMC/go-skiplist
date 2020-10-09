package skiplist

type SkiplistIterator struct {
	node *SkipListNode
	rank uint
}

func newIterator(node *SkipListNode, rank uint) *SkiplistIterator {
	return &SkiplistIterator{node: node, rank: rank}
}

func (iter *SkiplistIterator) Next() bool {
	if iter.node == nil {
		return false
	}

	iter.rank += iter.node.level[0].span
	iter.node = iter.node.level[0].forward
	return iter.node != nil
}

func (iter *SkiplistIterator) Rank() uint {
	return iter.rank
}

func (iter *SkiplistIterator) Value() int {
	return iter.node.val
}

func (iter *SkiplistIterator) Object() interface{} {
	return iter.node.obj
}
