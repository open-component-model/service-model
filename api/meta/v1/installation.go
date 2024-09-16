package v1

import (
	"github.com/mandelsoft/goutils/sliceutils"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
)

type TargetEnvironment = ocmmeta.StringMap

type ResourceReference ocmmeta.ResourceReference

func (r *ResourceReference) AsResourceRef() *ocmmeta.ResourceReference {
	return (*ocmmeta.ResourceReference)(r)
}

func (r *ResourceReference) String() string {
	return (*ocmmeta.ResourceReference)(r).String()
}

func (r *ResourceReference) Copy() *ResourceReference {
	if r == nil {
		return r
	}
	c := *r
	c.Resource = r.Resource.Copy()
	c.ReferencePath = sliceutils.InitialSliceFor(r.ReferencePath)
	for i, e := range r.ReferencePath {
		c.ReferencePath[i] = e.Copy()
	}
	return &c
}
