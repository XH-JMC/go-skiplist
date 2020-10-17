package skiplist_test

import (
	"XH-JMC/go-skiplist"
	"testing"
)

var (
	sl *skiplist.SkipList
	n  int
)

func init() {
	sl = skiplist.NewSkipList().WithObjectCompareFunc(func(a, b skiplist.SkipListElem) int {
		return a.(int) - b.(int)
	})

	n = 100
	step := 4
	for i := 1; i <= n; i += step {
		sl.Insert(i)
	}
	for i := 2; i <= n; i += step {
		sl.Insert(i)
	}
	for i := n - n%step; i > 0; i -= step {
		sl.Insert(i)
	}
	for i := n - n%step - 1; i > 0; i -= step {
		sl.Insert(i)
	}
}

func assert(t *testing.T, condition bool, fatalArgs ...interface{}) {
	if !condition {
		// t.Fail()
		t.Fatal(fatalArgs...)
	}
}

func TestSkipList_Delete(t *testing.T) {
	assert(t, sl.Delete(0) == false)
	assert(t, sl.Delete(2))
	assert(t, sl.Delete(1))
	assert(t, sl.Delete(n-1))
	assert(t, sl.Delete(n))
	assert(t, sl.Delete(n+1) == false)

	cnt := 0
	elem := 3
	rank := uint(1)
	iter := sl.Begin()
	for iter.Next() {
		assert(t, iter.Elem().(int) == elem, iter.Elem(), elem)
		assert(t, iter.Rank() == rank, iter.Rank(), rank)
		elem++
		rank++
		cnt++
	}
	assert(t, n-4 == cnt, n, cnt)
}

func TestSkipList_Find(t *testing.T) {
	check := func(elem int, exist bool) {
		elemRes, ok := sl.Find(elem)
		assert(t, ok == exist, elem, ok, exist)
		if ok {
			assert(t, elemRes.(int) == elem, elemRes, elem)
		}
	}

	check(0, false)
	check(1, true)
	check(2, true)
	check(n-1, true)
	check(n, true)
	check(n+1, false)
}

func TestSkipList_Delete_Find(t *testing.T) {
	check := func(elem int, exist bool) {
		_, ok := sl.Find(elem)
		assert(t, ok == exist)
		ok = sl.Delete(elem)
		assert(t, ok == exist)
		_, ok = sl.Find(elem)
		assert(t, ok == false)
	}

	check(0, false)
	check(1, true)
	check(2, true)
	check(n-1, true)
	check(n, true)
	check(n+1, false)

	cnt := 0
	elem := 3
	rank := uint(1)
	iter := sl.Begin()
	for iter.Next() {
		assert(t, iter.Elem().(int) == elem, iter.Elem(), elem)
		assert(t, iter.Rank() == rank, iter.Rank(), rank)
		elem++
		rank++
		cnt++
	}
	assert(t, n-4 == cnt, n, cnt)
}

func TestSkipList_LowerBound(t *testing.T) {
	check := func(queryElem int, elem int, rank uint, num int) {
		iter := sl.LowerBound(queryElem)
		cnt := 0
		for iter.Next() {
			assert(t, iter.Elem().(int) == elem, queryElem, iter.Elem(), elem)
			assert(t, iter.Rank() == rank, queryElem, iter.Rank(), rank)
			elem++
			rank++
			cnt++
		}
		assert(t, cnt == num, queryElem, num, cnt)
	}

	check(0, 1, 1, n)
	check(1, 1, 1, n)
	check(2, 2, 2, n-1)
	check(n-1, n-1, uint(n)-1, 2)
	check(n, n, uint(n), 1)
	check(n+1, n+1, uint(n)+1, 0)
}

func TestSkipList_UpperBound(t *testing.T) {
	check := func(queryElem int, elem int, rank uint, num int) {
		iter := sl.UpperBound(queryElem)
		cnt := 0
		for iter.Next() {
			assert(t, iter.Elem().(int) == elem, queryElem, iter.Elem(), elem)
			assert(t, iter.Rank() == rank, queryElem, iter.Rank(), rank)
			elem++
			rank++
			cnt++
		}
		assert(t, cnt == num, queryElem, num, cnt)
	}

	check(0, 1, 1, n)
	check(1, 2, 2, n-1)
	check(2, 3, 3, n-2)
	check(n-1, n, uint(n), 1)
	check(n, n+1, uint(n)+1, 0)
	check(n+1, n+2, uint(n)+2, 0)
}

func TestSkipList_LowerBoundByRank(t *testing.T) {
	check := func(queryRank uint, elem int, rank uint, num int) {
		iter := sl.LowerBoundByRank(queryRank)
		cnt := 0
		for iter.Next() {
			assert(t, iter.Elem().(int) == elem, queryRank, iter.Elem(), elem)
			assert(t, iter.Rank() == rank, queryRank, iter.Rank(), rank)
			elem++
			rank++
			cnt++
		}
		assert(t, cnt == num, queryRank, num, cnt)
	}

	check(0, 1, 1, n)
	check(1, 1, 1, n)
	check(2, 2, 2, n-1)
	check(uint(n)-1, n-1, uint(n)-1, 2)
	check(uint(n), n, uint(n), 1)
	check(uint(n)+1, n+1, uint(n)+1, 0)
}
