package provider

import (
	"fmt"
	"github.com/open-component-model/service-model/api/model/internal"
	"github.com/open-component-model/service-model/api/model/internal/common"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/provider"
)

func init() {
	internal.DefaultServiceKindRegistry.Register(provider.TYPE, Factory)
}

type ServiceVersionVariant struct {
	*common.ServiceVersionVariant
	spec *provider.ServiceSpec
}

func Factory(model internal.Model, descriptor *modeldesc.ServiceDescriptor) (internal.ServiceVersionVariant, error) {
	s, ok := descriptor.Kind.(*provider.ServiceSpec)
	if !ok {
		return nil, fmt.Errorf("invalid service spec type: %T", descriptor.Kind)
	}
	return &ServiceVersionVariant{
		ServiceVersionVariant: common.New(descriptor),
		spec:                  s,
	}, nil
}

func (s *ServiceVersionVariant) AsOrdinaryService() internal.OrdinaryService {
	return s
}
