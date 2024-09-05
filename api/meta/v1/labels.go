package v1

import (
	"ocm.software/ocm/api/utils/runtime"
)

type Labels []Label

func (l Labels) Copy() Labels {
	c := make([]Label, len(l), len(l))
	for i, e := range l {
		c[i] = *e.Copy()
	}
	return c
}

type Label struct {
	Name    string           `json:"name"`
	Version string           `json:"version"`
	Value   runtime.RawValue `json:"value"`
}

func (l *Label) Copy() *Label {
	return &Label{
		Name:    l.Name,
		Version: l.Version,
		Value:   l.Value.Copy(),
	}
}
