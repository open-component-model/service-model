package contract

import (
	"fmt"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/model/internal"
	"github.com/open-component-model/service-model/api/model/internal/common"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/contract"
	"ocm.software/ocm/api/utils/runtime"
)

func init() {
	internal.DefaultServiceKindRegistry.Register(contract.TYPE, Factory)
}

type ServiceVersion struct {
	*common.ServiceVersionVariant
	model internal.Model
	spec  *contract.ServiceSpec
}

var _ internal.ServiceContract = (*ServiceVersion)(nil)

func Factory(model internal.Model, descriptor *modeldesc.ServiceDescriptor) (internal.ServiceVersionVariant, error) {
	s, ok := descriptor.Kind.(*contract.ServiceSpec)
	if !ok {
		return nil, fmt.Errorf("invalid service spec type: %T", descriptor.Kind)
	}
	return &ServiceVersion{
		ServiceVersionVariant: common.New(descriptor),
		model:                 model,
		spec:                  s,
	}, nil
}

func (s *ServiceVersion) GetAPISpecificationType() string {
	return s.spec.APISpecificationType
}

func (s *ServiceVersion) GetAPISpecVersion() string {
	return s.spec.APISpecVersion
}

func (s *ServiceVersion) GetSpecification() *runtime.RawValue {
	return s.spec.Specification
}

// TODO: Initially, I thought it'd be nice to resolve that resource reference here. But, since we want to allow the
//  service model to be based on different data sources than ocm, we probably do not want that. But, considering that
//  intention, does it make sense to include such an ocm specific type here? Shouldn't we consequently aim for a
//  stricter decoupling here too?
//  Without thinking deeply about it, we instead introduce a ResourceResolver abstraction that could also be part of the
//  model which would also allow the GetArtifact() method to actually resolve that reference.

func (s *ServiceVersion) GetArtifact() *metav1.ResourceReference {
	return s.spec.Artifact
}

func (s *ServiceVersion) AsServiceContract() internal.ServiceContract {
	return s
}
