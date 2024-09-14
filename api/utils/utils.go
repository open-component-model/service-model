package utils

type Parseable interface {
	Parse(string) error
}

type parsablePointer[T any] interface {
	*T
	Parseable
}

func Parse[T any, P parsablePointer[T]](s string) (P, error) {
	var eff T
	err := (P(&eff)).Parse(s)
	if err != nil {
		return nil, err
	}
	return &eff, nil
}

type Comparable[T any] interface {
	Compare(T) int
}

func Compare[T Comparable[T]](a, b any) int {
	return a.(T).Compare(b.(T))
}
