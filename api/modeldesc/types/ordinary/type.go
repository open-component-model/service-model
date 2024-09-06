package ordinary

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

const TYPE = "OrdinaryService"

type ServiceSpec struct {
	metav1.CommonConsumerServiceImplementationSpec
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	return &ServiceSpec{*internal.CommonConsumerServiceImplementationSpecToCanonicalForm(&s.CommonConsumerServiceImplementationSpec, c)}
}

func (s *ServiceSpec) Validate(c internal.DescriptionContext) error {
	return internal.ValidateCommonConsumerImplementation(&s.CommonConsumerServiceImplementationSpec, c)
}
