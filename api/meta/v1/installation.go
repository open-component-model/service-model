package v1

import (
	"github.com/open-component-model/service-model/api/utils"
	v1 "ocm.software/ocm/api/ocm/compdesc/meta/v1"
)

type TargetEnvironment = v1.StringMap

type ResourceReference v1.ResourceReference

func (r *ResourceReference) AsResourceRef() *v1.ResourceReference {
	return (*v1.ResourceReference)(r)
}

func (r *ResourceReference) String() string {
	return (*v1.ResourceReference)(r).String()
}

func (r *ResourceReference) Copy() *ResourceReference {
	if r == nil {
		return r
	}
	c := *r
	c.Resource = r.Resource.Copy()
	c.ReferencePath = utils.InitialSliceFor(r.ReferencePath)
	for i, e := range r.ReferencePath {
		c.ReferencePath[i] = e.Copy()
	}
	return &c
}
