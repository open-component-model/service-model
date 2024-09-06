package utils

func InitialSliceFor[S ~[]E, E any](in S) S {
	return make(S, len(in), len(in))
}

func InitialSliceWithTypeFor[TS ~[]TE, TE any, S ~[]E, E any](in S) TS {
	return make(TS, len(in), len(in))
}

func Slice[T any](elems ...T) []T {
	return elems
}
