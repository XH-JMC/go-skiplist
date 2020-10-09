package skiplist

import (
	"math"
	"math/rand"
)

const (
	SKIPLIST_MAXLEVEL = 31 // level: [0, MaxLevel]
	SKIPLIST_P        = 4  // 概率: 1/4
)

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
	span    uint
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

func (s *SkipList) Insert(val int, obj interface{}) {
	newLevel := s.randLevel()
	if s.level < newLevel {
		s.head.level = append(s.head.level, make([]SkipListLevel, newLevel-s.level)...)
		s.level = newLevel
	}
	s.size++

	node := &SkipListNode{
		level:    make([]SkipListLevel, newLevel+1),
		backward: nil,
		val:      val,
		obj:      obj,
	}

	level := s.level // 当前遍历的层级
	p := s.head      // 当前遍历的节点
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && forward.val < val {
			p = forward
			forward = p.level[level].forward
		}
		if level <= newLevel {
			p.level[level].forward = node
			node.level[level].forward = forward
			if forward != nil {
				forward.backward = node
			}
		}
		level--
	}
	node.backward = p
}

func (s *SkipList) randLevel() int {
	level := 0
	for (rand.Int()&math.MaxInt32) < (math.MaxInt32/SKIPLIST_P) {
		level++
	}
	if level < SKIPLIST_MAXLEVEL {
		return level
	}
	return SKIPLIST_MAXLEVEL
}

func (s *SkipList) LowerBound(val int) *SkiplistIterator {
	var rank uint
	level := s.level // 当前遍历的层级
	p := s.head      // 当前遍历的节点
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && forward.val < val {
			rank += p.level[level].span
			p = forward
			forward = p.level[level].forward
		}
		level--
	}

	return newIterator(p, rank)
}
