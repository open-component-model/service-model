package ordinary

import (
	"github.com/open-component-model/service-model/api/crossref"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

const TYPE = "OrdinaryService"

type ServiceSpec struct {
	metav1.CommonConsumerServiceImplementationSpec
}

func (s *ServiceSpec) Copy() internal.ServiceKindSpec {
	return &ServiceSpec{
		CommonConsumerServiceImplementationSpec: *s.CommonConsumerServiceImplementationSpec.Copy(),
	}
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	return &ServiceSpec{*internal.CommonConsumerServiceImplementationSpecToCanonicalForm(&s.CommonConsumerServiceImplementationSpec, c)}
}

func (s *ServiceSpec) Validate(c internal.DescriptionContext) error {
	return internal.ValidateCommonConsumerServiceImplementation(&s.CommonConsumerServiceImplementationSpec, c)
}

func (s *ServiceSpec) GetReferences() crossref.References {
	return internal.CommonConsumerServiceImplementationReferences(&s.CommonConsumerServiceImplementationSpec)
}
