package vpi

import (
	"github.com/mandelsoft/goutils/sliceutils"
)

type Converter[I any, E any] interface {
	// ConvertFrom converts from an internal version into an external format.
	ConvertFrom(object I) (E, error)
	// ConvertTo converts from an external format into an internal version.
	ConvertTo(object E) (I, error)
}

////////////////////////////////////////////////////////////////////////////////

type ListConverter[I, E any] struct {
	element Converter[*I, *E]
}

func NewListConverter[I, E any](element Converter[*I, *E]) *ListConverter[I, E] {
	return &ListConverter[I, E]{element}
}

func (l *ListConverter[I, E]) ConvertFrom(in []I) ([]E, error) {
	out := sliceutils.InitialSliceWithTypeFor[[]E](in)
	for i, e := range in {
		r, err := l.element.ConvertFrom(&e)
		if err != nil {
			return nil, err
		}
		out[i] = *r
	}
	return out, nil
}

func (l *ListConverter[I, E]) ConvertTo(in []E) ([]I, error) {
	out := sliceutils.InitialSliceWithTypeFor[[]I](in)
	for i, e := range in {
		r, err := l.element.ConvertTo(&e)
		if err != nil {
			return nil, err
		}
		out[i] = *r
	}
	return out, nil
}
