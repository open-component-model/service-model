package v1

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/utils"
	"ocm.software/ocm/api/utils/runtime"
)

type Labels []Label

func (l Labels) Copy() Labels {
	c := utils.InitialSliceFor(l)
	for i, e := range l {
		c[i] = *e.Copy()
	}
	return c
}

func (l Labels) Validate() error {
	var err errors.ErrorList
	for i, e := range l {
		err.Addf(nil, e.Validate(), "label %d(%s)", i, e.Name)
	}
	return err.Result()
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
func (l Label) Validate() error {
	return utils.CheckNonEmpty(l.Name, "label name")
}
