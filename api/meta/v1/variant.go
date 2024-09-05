package v1

import (
	"maps"
)

type Variant map[string]string

func (v Variant) Copy() Variant {
	return maps.Clone(v)
}
