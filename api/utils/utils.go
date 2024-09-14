package utils

type Parseable interface {
	Parse(string) error
}

type Comparable[T any] interface {
	Compare(T) int
}

func Compare[T Comparable[T]](a, b any) int {
	return a.(T).Compare(b.(T))
}
