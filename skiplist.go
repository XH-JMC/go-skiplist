package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

const (
	SKIPLIST_MAXLEVEL = 31 // level: [0, MaxLevel]
	SKIPLIST_P        = 4  // 概率: 1/4
)

type ElemCompareFunc func(a, b SkipListElem) int // 比较函数: a<b返回负数，a>b返回正数，a==b返回0

type SkipList struct {
	head    *SkipListNode
	tail    *SkipListNode
	size    uint
	level   int
	elemCmp ElemCompareFunc
}

type SkipListNode struct {
	level    []SkipListLevel
	backward *SkipListNode

	elem SkipListElem
}

type SkipListLevel struct {
	forward *SkipListNode
	span    uint
}

type SkipListElem interface{}

func NewSkipList() *SkipList {
	node := &SkipListNode{
		level:    make([]SkipListLevel, SKIPLIST_MAXLEVEL+1),
		backward: nil,
	}

	return &SkipList{
		head:  node,
		tail:  node,
		size:  0,
		level: -1,
	}
}

func (sl *SkipList) WithObjectCompareFunc(objCmp ElemCompareFunc) *SkipList {
	sl.elemCmp = objCmp
	return sl
}

func DefaultObjectCompare(a, b SkipListElem) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return +1
	}
	return strings.Compare(fmt.Sprint(a), fmt.Sprint(b))
}

func (sl *SkipList) Size() uint {
	return sl.size
}

func (sl *SkipList) cmpElem(a, b SkipListElem) int {
	if sl.elemCmp == nil {
		return DefaultObjectCompare(a, b)
	}

	return sl.elemCmp(a, b)
}

func (sl *SkipList) Insert(elem SkipListElem) {
	newLevel := sl.randLevel()
	if sl.level < newLevel {
		sl.level = newLevel
	}
	sl.size++

	node := &SkipListNode{
		level:    make([]SkipListLevel, newLevel+1),
		backward: nil,
		elem:     elem,
	}

	type NodeRank struct {
		node *SkipListNode
		rank uint
	}

	preNodeRanks := make([]NodeRank, sl.level+1) // 新节点每层的前置节点及排名
	rank := uint(0)                              // 当前遍历的节点的排名
	level := sl.level                            // 当前遍历的层级
	p := sl.head                                 // 当前遍历的节点
	// 在每一层插入新节点，同时维护新节点每层的前置节点及排名
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && sl.cmpElem(forward.elem, elem) <= 0 {
			rank += p.level[level].span
			p = forward
			forward = p.level[level].forward
		}

		preNodeRanks[level] = NodeRank{
			node: p,
			rank: rank,
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

	if p != sl.head {
		node.backward = p
	}
	if node.level[0].forward == nil {
		sl.tail = node
	}

	// 根据新节点每层的前置节点及排名和旧span，求新节点与每一层的前后节点的新span
	newNodeRank := preNodeRanks[0].rank + 1
	for i := 0; i <= sl.level; i++ {
		if i <= newLevel {
			preNodeSpan := newNodeRank - preNodeRanks[i].rank
			if node.level[i].forward != nil {
				node.level[i].span = preNodeRanks[i].node.level[i].span + 1 - preNodeSpan
			}
			preNodeRanks[i].node.level[i].span = preNodeSpan
		} else if preNodeRanks[i].node.level[i].forward != nil {
			preNodeRanks[i].node.level[i].span++
		}
	}
}

// 返回[0, SKIPLIST_MAXLEVEL]的随机数
func (sl *SkipList) randLevel() int {
	level := 0
	for (rand.Int() & math.MaxInt32) < (math.MaxInt32 / SKIPLIST_P) {
		level++
	}
	if level < SKIPLIST_MAXLEVEL {
		return level
	}
	return SKIPLIST_MAXLEVEL
}

// 删除特定元素，返回原表中是否存在该元素
func (sl *SkipList) Delete(elem SkipListElem) bool {
	preNodes := make([]*SkipListNode, sl.level+1) // 记录新节点每层的前置节点
	level := sl.level                             // 当前遍历的层级
	p := sl.head                                  // 当前遍历的节点
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && sl.cmpElem(forward.elem, elem) < 0 {
			p = forward
			forward = p.level[level].forward
		}
		preNodes[level] = p
		level--
	}

	p = p.level[0].forward
	if p != nil && sl.cmpElem(p.elem, elem) == 0 {
		sl.deleteNode(p, preNodes)
		return true
	}
	return false
}

func (sl *SkipList) deleteNode(p *SkipListNode, preNodes []*SkipListNode) {
	for i, preNode := range preNodes {
		if preNode.level[i].forward == p {
			preNode.level[i].span += p.level[i].span - 1
			preNode.level[i].forward = p.level[i].forward
		} else {
			preNode.level[i].span--
		}
	}

	if p.level[0].forward != nil {
		p.level[0].forward.backward = p.backward
	} else {
		sl.tail = p.backward
	}

	for sl.level >= 0 && sl.head.level[sl.level].forward == nil {
		sl.level--
	}
	sl.size--
}

func (sl *SkipList) Find(elem SkipListElem) (SkipListElem, bool) {
	iter := sl.LowerBound(elem)
	if iter.Next() && sl.cmpElem(iter.Elem(), elem) == 0 {
		return iter.Elem(), true
	}
	return nil, false
}

func (sl *SkipList) LowerBound(elem SkipListElem) *SkiplistIterator {
	return sl.findWithLessFunc(func(p, forward *SkipListNode, level int, rank uint) bool {
		return sl.cmpElem(forward.elem, elem) < 0
	})
}

func (sl *SkipList) UpperBound(elem SkipListElem) *SkiplistIterator {
	return sl.findWithLessFunc(func(p, forward *SkipListNode, level int, rank uint) bool {
		return sl.cmpElem(forward.elem, elem) <= 0
	})
}

func (sl *SkipList) LowerBoundByRank(targetRank uint) *SkiplistIterator {
	return sl.findWithLessFunc(func(p, forward *SkipListNode, level int, rank uint) bool {
		return rank+p.level[level].span < targetRank
	})
}

// less 返回当前遍历的节点p是否在forward前，level是当前遍历的层级，rank是当前遍历的节点的排名
func (sl *SkipList) findWithLessFunc(less func(p, forward *SkipListNode, level int, rank uint) bool) *SkiplistIterator {
	var rank uint     // 当前遍历的节点的排名
	level := sl.level // 当前遍历的层级
	p := sl.head      // 当前遍历的节点
	for level >= 0 {
		forward := p.level[level].forward
		for forward != nil && less(p, forward, level, rank) {
			rank += p.level[level].span
			p = forward
			forward = p.level[level].forward
		}
		level--
	}

	return newIterator(p, rank)
}

func (sl *SkipList) Begin() *SkiplistIterator {
	return newIterator(sl.head, 0)
}

func (sl *SkipList) End() *SkiplistIterator {
	if sl.size < 1 {
		return newIterator(sl.head, 0)
	}
	return newIterator(sl.tail.backward, sl.size-1)
}
