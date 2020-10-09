package skiplist

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := NewSkipList()
	n := 1000000
	for i := 0; i < n; i++ {
		s.Insert(i, nil)
	}

	list := make([]int, SKIPLIST_MAXLEVEL+1)
	p := s.head.level[0].forward
	for p != nil {
		list[len(p.level)-1]++
		p = p.level[0].forward
	}
	for i := SKIPLIST_MAXLEVEL; i > 0; i-- {
		list[i-1] += list[i]
	}
	for i, num := range list {
		fmt.Println(i, num)
	}

	f := func(val int) {
		iter := s.LowerBound(val)
		_ = iter

		cnt := 0
		if iter.Next() {
			fmt.Println(val, iter.Rank(), iter.Value(), iter.Object())
			cnt++
		} else {
			fmt.Println(val, nil)
		}

		for iter.Next() {
			cnt++
		}

		if n != val+cnt {
			fmt.Println(val, cnt)
		}
	}

	f(0)
	f(1)
	f(2)
	f(n - 2)
	f(n - 1)
	f(n)
}

func TestRandLevel(t *testing.T) {
	s := NewSkipList()
	n := 1000000
	list := make([]int, SKIPLIST_MAXLEVEL+1)
	for i := 0; i < n; i++ {
		level := s.randLevel()
		list[level]++
	}
	for i := SKIPLIST_MAXLEVEL; i > 0; i-- {
		list[i-1] += list[i]
	}
	for i, num := range list {
		fmt.Println(i, num)
	}
}
