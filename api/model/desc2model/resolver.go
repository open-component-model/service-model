package desc2model

import (
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/model"
	"github.com/open-component-model/service-model/api/modeldesc"
)

type resolver struct {
	resolver modeldesc.VersionResolver
}

func NewResolver(r modeldesc.VersionResolver) model.Resolver {
	return &resolver{r}
}

func (r resolver) LookupServiceVersionVariant(model model.Model, id identity.ServiceVersionVariantIdentity) (model.ServiceVersionVariant, error) {
	s, err := r.resolver.LookupServiceVersionVariant(id)
	if s == nil || err != nil {
		return nil, err
	}
	return ServiceDescriptorToModelService(model, s)
}

func (r resolver) ListVersions(id identity.ServiceIdentity, variant ...identity.Variant) ([]string, error) {
	return r.resolver.ListVersions(id, variant...)
}
