package v1

import (
	"maps"
)

type Copyable[T any] interface {
	Copy() *T
}

type CopyableList[T Copyable[T]] []T

func (l CopyableList[T]) Copy() CopyableList[T] {
	c := make([]T, len(l), len(l))
	for i, e := range l {
		c[i] = *e.Copy()
	}
	return c
}

type StringMap map[string]string

func (m StringMap) Copy() StringMap {
	return maps.Clone(m)
}
