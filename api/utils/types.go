package utils

import (
	"maps"
	"reflect"

	"github.com/mandelsoft/goutils/generics"
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

func Convert[T, S any](in []S) []T {
	if in == nil {
		return nil
	}
	t := generics.TypeOf[T]()
	out := make([]T, len(in))
	for i, v := range in {
		out[i] = reflect.ValueOf(v).Convert(t).Interface().(T)
	}
	return out
}
