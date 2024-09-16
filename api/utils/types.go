package utils

import (
	"maps"

	"github.com/mandelsoft/goutils/sliceutils"
)

type Copyable[T any] interface {
	Copy() *T
}

type CopyableList[T Copyable[T]] []T

func (l CopyableList[T]) Copy() CopyableList[T] {
	c := sliceutils.InitialSliceFor(l)
	for i, e := range l {
		c[i] = *e.Copy()
	}
	return c
}

type StringMap map[string]string

func (m StringMap) Copy() StringMap {
	return maps.Clone(m)
}
