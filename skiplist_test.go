package skiplist_test

import (
	skiplist "go-skiplist"
	"testing"
)

var (
	s *skiplist.SkipList
	n int
)

func init() {
	s = skiplist.NewSkipList().WithObjectCompareFunc(func(a, b skiplist.SkipListElem) int {
		return a.(int) - b.(int)
	})

	n = 100
	step := 4
	for i := 1; i <= n; i += step {
		s.Insert(i)
	}
	for i := 2; i <= n; i += step {
		s.Insert(i)
	}
	for i := n - n%step; i > 0; i -= step {
		s.Insert(i)
	}
	for i := n - n%step - 1; i > 0; i -= step {
		s.Insert(i)
	}
}

func assert(t *testing.T, condition bool, fatalArgs ...interface{}) {
	if !condition {
		// t.Fail()
		t.Fatal(fatalArgs...)
	}
}

func TestSkipList_LowerBound(t *testing.T) {
	check := func(queryElem int, elem int, rank uint, exist bool) {
		iter := s.LowerBound(queryElem)
		queryExist := false
		for iter.Next() {
			queryExist = true
			assert(t, iter.Elem().(int) == elem, queryElem, iter.Elem(), elem)
			assert(t, iter.Rank() == rank, queryElem, iter.Rank(), rank)
			elem++
			rank++
		}
		assert(t, queryExist == exist, queryElem, queryExist, exist)
	}

	check(0, 1, 1, true)
	check(1, 1, 1, true)
	check(2, 2, 2, true)
	check(n-1, n-1, uint(n)-1, true)
	check(n, n, uint(n), true)
	check(n+1, n+1, uint(n)+1, false)
}

func TestSkipList_UpperBound(t *testing.T) {
	// todo
}

func TestSkipList_LowerBoundByRank(t *testing.T) {
	// todo
}

func TestSkipList_Delete(t *testing.T) {
	// todo
}
