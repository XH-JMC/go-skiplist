package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test(t *testing.T) {
	rand.Seed(0)

	s := NewSkipList()
	n := 100000
	for i := 0; i < n; i++ {
		s.Add(i, nil)
	}

	m := make([]int, MaxLevel+1)
	p := s.head.level[0].forward
	for p != nil {
		m[len(p.level)-1]++
		p = p.level[0].forward
	}
	for i, v := range m {
		fmt.Println(i, v)
	}

	//f := func(val int) {
	//	cnt := 0
	//
	//	iter := s.LowerBound(val)
	//	_ = iter
	//	for iter.Next() {
	//		_, _ = iter.Value(), iter.Object()
	//		cnt++
	//	}
	//
	//	if n != val+cnt {
	//		fmt.Println(val, cnt)
	//	}
	//}
	//
	//for i := 0; i < n; i++ {
	//	f(rand.Intn(n + 1))
	//}
}
