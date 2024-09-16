package ocm

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/desc2model"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/model"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/ocm"
)

type resolver struct {
	modeldesc.VersionResolver
}

var _ model.Resolver = (*resolver)(nil)

func NewResolver(r ocm.ComponentResolver) model.Resolver {
	return &resolver{
		NewVersionResolver(r),
	}
}

func (r *resolver) LookupServiceVersionVariant(model model.Model, id identity.ServiceVersionVariantIdentity) (model.Service, error) {
	if id.IsConstraint() {
		return nil, errors.ErrInvalid(modeldesc.KIND_SERVICEVERSION, id.String())
	}

	s, err := r.VersionResolver.LookupServiceVersionVariant(id)
	if err != nil {
		return nil, err
	}
	return desc2model.ServiceDescriptorToModelService(model, s)
}
