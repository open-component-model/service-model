package v1

import (
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/utils"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	"slices"
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

////////////////////////////////////////////////////////////////////////////////

type InstalledService struct {
	Service  identity.ServiceIdentity `json:"service,omitempty"`
	Versions []string                 `json:"versions,omitempty"`
	Variant  identity.Variant         `json:"variant,omitempty"`
}

func (s InstalledService) Copy() *InstalledService {
	s.Variant = s.Variant.Copy()
	s.Versions = slices.Clone(s.Versions)
	return &s
}

type InstalledServices = utils.CopyableList[InstalledService]
