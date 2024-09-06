package v1

import (
	v1 "ocm.software/ocm/api/ocm/compdesc/meta/v1"
)

type TargetEnvironment = v1.StringMap

type InstallerResource v1.ResourceReference

func (r *InstallerResource) String() string {
	return (*v1.ResourceReference)(r).String()
}

func (r *InstallerResource) Copy() *InstallerResource {
	if r == nil {
		return r
	}
	c := *r
	c.Resource = r.Resource.Copy()
	c.ReferencePath = make([]v1.Identity, len(r.ReferencePath), len(r.ReferencePath))
	for i, e := range r.ReferencePath {
		c.ReferencePath[i] = e.Copy()
	}
	return &c
}
