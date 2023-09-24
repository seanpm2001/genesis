package iters_test

import (
	"testing"

	"github.com/life4/genesis/iters"
	"github.com/matryer/is"
)

var ts = iters.ToSlice[int]

func new[T any](vs ...T) iters.Iter[T] {
	return iters.FromSlice(vs)
}

func TestFromChannel(t *testing.T) {
	is := is.New(t)
	ch := make(chan int)
	go func() {
		ch <- 3
		ch <- 4
		ch <- 5
		close(ch)
	}()
	is.Equal(ts(iters.FromChannel(ch)), []int{3, 4, 5})
}

func TestFromSlice(t *testing.T) {
	is := is.NewRelaxed(t)
	is.Equal(ts(iters.FromSlice([]int{3, 4, 5})), []int{3, 4, 5})
	is.Equal(ts(iters.FromSlice([]int{3})), []int{3})
	is.Equal(ts(iters.FromSlice([]int{})), []int{})
	is.Equal(ts(iters.FromSlice([]int(nil))), []int{})
}

func TestMap(t *testing.T) {
	is := is.NewRelaxed(t)
	double := func(x int) int { return x * 2 }
	is.Equal(ts(iters.Map(new(3, 4, 5), double)), []int{6, 8, 10})
	is.Equal(ts(iters.Map(new[int](), double)), []int{})
}

func TestTake(t *testing.T) {
	is := is.NewRelaxed(t)
	is.Equal(ts(iters.Take(new(3, 4, 5), 2)), []int{3, 4})
	is.Equal(ts(iters.Take(new(3, 4, 5), 1)), []int{3})
	is.Equal(ts(iters.Take(new(3, 4, 5), 10)), []int{3, 4, 5})
	is.Equal(ts(iters.Take(new(3, 4, 5), 0)), []int{})
	is.Equal(ts(iters.Take(new(3, 4, 5), -1)), []int{})
	is.Equal(ts(iters.Take(new(3, 4, 5), -10)), []int{})
}

func TestToSlice(t *testing.T) {
	is := is.NewRelaxed(t)
	is.Equal(iters.ToSlice(new(3, 4, 5)), []int{3, 4, 5})
	is.Equal(iters.ToSlice(new[int]()), []int{})
}