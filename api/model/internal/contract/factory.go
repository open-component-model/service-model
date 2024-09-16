package contract

import (
	"fmt"
	"github.com/open-component-model/service-model/api/model/internal"
	"github.com/open-component-model/service-model/api/model/internal/common"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/contract"
)

func init() {
	internal.DefaultFactory.Register(contract.TYPE, Factory)
}

type ServiceVersion struct {
	*common.ServiceVersionVariant
	spec *contract.ServiceSpec
}

func Factory(model internal.Model, descriptor *modeldesc.ServiceDescriptor) (internal.ServiceVersionVariant, error) {
	s, ok := descriptor.Kind.(*contract.ServiceSpec)
	if !ok {
		return nil, fmt.Errorf("invalid service spec type: %T", descriptor.Kind)
	}
	return &ServiceVersion{
		ServiceVersionVariant: common.New(descriptor),
		spec:                  s,
	}, nil
}
