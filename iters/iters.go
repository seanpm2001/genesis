package iters

import c "github.com/life4/genesis/constraints"

// Iter is an interface describing a lazy iterator.
//
// It's very similar to channels except that it's not thread-safe
// and so should be faster when you don't need concurrency.
type Iter[T any] interface {
	// Next returns the next element from the iterator.
	//
	// The second return value indicates if there are more values to pull.
	// If the iterator is exhausted, the first value is the default value
	// of the type and second is false. When the iterator is exhausted,
	// repeated attempts to get the next value should produce the same
	// default+false result.
	//
	// In other words, it should behave like pulling from a (closed) channel.
	//
	// The code using an iterator doesn't guarantee to exhaust it.
	// For example, [Take] only takes the number of elements it needs
	// and never calls Next again. Hence you shouldn't rely on Next
	// for closing connections and cleaning up unused resources.
	// If your iterator needs to provide logic like this, you should
	// implement a Close method and defer it.
	//
	// An iterator is allowed to be infinite and never return false.
	Next() (T, bool)
}

// Filter returns an iterator of elements from the given iterator for which the function returns true.
func Filter[T any](it Iter[T], f func(T) bool) Iter[T] {
	return FromFunc(func() (T, bool) {
		for {
			val, more := it.Next()
			if !more {
				return val, false
			}
			if f(val) {
				return val, true
			}
		}
	})
}

// FromChannel produces an iterator returning elements from the given channel.
//
// Each call to Iter.Next will pull from the channel, which means
// you have to make sure it won't block forever. It's a good idea
// to make the channel cancelable by using channels.WithContext.
func FromChannel[T any](ch <-chan T) Iter[T] {
	return FromFunc(func() (T, bool) {
		v, ok := <-ch
		return v, ok
	})
}

// FromFunc makes an iterator from the given function.
//
// Each time a new iterator element is requested, the function is called.
func FromFunc[T any](f func() (T, bool)) Iter[T] {
	return iFunc[T]{f}
}

type iFunc[T any] struct {
	f func() (T, bool)
}

func (i iFunc[T]) Next() (T, bool) {
	return i.f()
}

// FromSlice produces an iterator returning elements from the given slice.
func FromSlice[S ~[]T, T any](slice S) Iter[T] {
	next := 0
	return FromFunc(func() (T, bool) {
		if next >= len(slice) {
			return *new(T), false
		}
		v := slice[next]
		next += 1
		return v, true
	})
}

// Map returns an iterator of results of applying the function to each element of the given iterator.
func Map[T, R any](it Iter[T], f func(T) R) Iter[R] {
	return FromFunc(func() (R, bool) {
		val, more := it.Next()
		if !more {
			var res R
			return res, false
		}
		res := f(val)
		return res, true
	})
}

// Reduce applies the function to acc and every iterator element and returns the acc.
func Reduce[T, R any](it Iter[T], acc R, f func(T, R) R) R {
	for {
		val, more := it.Next()
		if !more {
			return acc
		}
		acc = f(val, acc)
	}
}

// Take returns an iterator returning only the first n items from the given iterator.
//
// When n items are consumed, Take will not call Next on the input iterator again.
// So, it's possible for the input iterator to not be fully exhausted.
//
// If the input iterator returns fewer than n items, Take will just stop and
// not generate additional items.
func Take[T any, I c.Integer](it Iter[T], n I) Iter[T] {
	return FromFunc(func() (T, bool) {
		if n <= 0 {
			var val T
			return val, false
		}
		n -= 1
		return it.Next()
	})
}

// ToSlice converts the given iterator to a slice.
//
// The function returns only when there are no more elements to consume
// from the iterator. It's a good idea to use [Take] to limit the number
// of elements if it's possible for the iterator to be infinite or just too big.
//
// Also, you should make sure that the iterator doesn't block forever.
// In particular, when creating an iterator from a channel using [FromChannel],
// you may want to use channels.WithContext and set a deadline or cancelation
// on that context.
func ToSlice[T any](it Iter[T]) []T {
	res := make([]T, 0)
	for {
		val, more := it.Next()
		if !more {
			return res
		}
		res = append(res, val)
	}
}
