package skiplist

import "math/rand"

const MaxLevel = 63 // level: [0, MaxLevel]

type SkipList struct {
	head  *SkipListNode
	size  uint
	level int
}

type SkipListNode struct {
	level    []SkipListLevel
	backward *SkipListNode

	val int
	obj interface{}
}

type SkipListLevel struct {
	forward *SkipListNode
	//span    uint
}

func NewSkipList() *SkipList {
	return &SkipList{
		head: &SkipListNode{
			level:    make([]SkipListLevel, 1),
			backward: nil,
			val:      0,
			obj:      nil,
		},
		size:  0,
		level: 0,
	}
}

func (s *SkipList) Size() uint {
	return s.size
}

func (s *SkipList) Add(val int, obj interface{}) {
	level := s.randLevel()
	if s.level < level {
		s.head.level = append(s.head.level, make([]SkipListLevel, level-s.level)...)
		s.level = level
	}
	s.size++

	node := &SkipListNode{
		level:    make([]SkipListLevel, level+1),
		backward: nil,
		val:      val,
		obj:      obj,
	}

	p := s.head
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && forward.val < val {
			p = forward
			forward = p.level[level].forward
		}
		p.level[level].forward = node
		node.level[level].forward = forward
		if forward != nil {
			forward.backward = node
		}
		level--
	}
	node.backward = p
}

func (s *SkipList) randLevel() int {
	for i := 0; i < MaxLevel; i++ {
		if (rand.Int() & 1) == 0 {
			return i
		}
	}
	return MaxLevel
}

func (s *SkipList) LowerBound(val int) *SkiplistIterator {
	level := s.level // 当前遍历的层级
	p := s.head      // 当前遍历的节点
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && forward.val < val {
			p = forward
			forward = p.level[level].forward
		}
		level--
	}

	return newIterator(p)
}
