package installer

import (
	"fmt"
	"github.com/open-component-model/service-model/api/identity"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/model"
	"github.com/open-component-model/service-model/api/model/internal"
	"github.com/open-component-model/service-model/api/model/internal/common"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/ordinary"
)

func init() {
	internal.DefaultServiceKindRegistry.Register(ordinary.TYPE, Factory)
}

type ServiceVersionVariant struct {
	model model.Model
	*common.ServiceVersionVariant
	spec *ordinary.ServiceSpec
}

// TODO: since this is always the same, might be able to extract this logic into the registry or at least provide a
//	generic Factory() function to do this.

func Factory(model internal.Model, descriptor *modeldesc.ServiceDescriptor) (internal.ServiceVersionVariant, error) {
	s, ok := descriptor.Kind.(*ordinary.ServiceSpec)
	if !ok {
		return nil, fmt.Errorf("invalid service spec type: %T", descriptor.Kind)
	}
	return &ServiceVersionVariant{
		model:                 model,
		ServiceVersionVariant: common.New(descriptor),
		spec:                  s,
	}, nil
}

func (s *ServiceVersionVariant) AsOrdinaryService() internal.OrdinaryService {
	return s
}

func (s *ServiceVersionVariant) GetVariant() identity.Variant {
	return s.spec.Variant
}

func (s *ServiceVersionVariant) IsAbstract() bool {
	return s.spec.Abstract
}

func (s *ServiceVersionVariant) GetInheritFrom() identity.Variant {
	return s.spec.InheritFrom
}
func (s *ServiceVersionVariant) GetInheritFromService() (internal.ServiceVersionVariant, error) {
	if s.spec.InheritFrom == nil {
		return nil, fmt.Errorf("service does not inherit and therefore has no parent")
	}
	return s.model.GetServiceVersionVariant(identity.NewServiceVersionVariantIdentityFor(s.GetIdentity().ServiceVersionIdentity(), s.spec.InheritFrom))
}

func (s *ServiceVersionVariant) GetDependencies() v1.Dependencies {
	return s.spec.Dependencies[0]
}
