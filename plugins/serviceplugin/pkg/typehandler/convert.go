package typehandler

import (
	"ocm.software/ocm/cmds/ocm/common/data"
)

type Objects[T any] []T

func ObjectSlice[T any](s data.Iterable) Objects[T] {
	var a Objects[T]
	i := s.Iterator()
	for i.HasNext() {
		a = append(a, i.Next().(T))
	}
	return a
}

var (
	_ data.IndexedAccess = Objects[int]{}
	_ data.Iterable      = Objects[int]{}
)

func (o Objects[T]) Len() int {
	return len(o)
}

func (o Objects[T]) Get(i int) interface{} {
	return o[i]
}

func (o Objects[T]) Iterator() data.Iterator {
	return data.NewIndexedIterator(o)
}
