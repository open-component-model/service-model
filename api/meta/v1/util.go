package v1

type Copyable[T any] interface {
	Copy() *T
}

type List[T Copyable[T]] []T

func (l List[T]) Copy() List[T] {
	c := make([]T, len(l), len(l))
	for i, e := range l {
		c[i] = *e.Copy()
	}
	return c
}
