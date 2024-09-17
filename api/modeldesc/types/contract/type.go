package contract

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/generics"
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/identity"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
	"ocm.software/ocm/api/utils/runtime"
)

const TYPE = "ServiceContract"

type ServiceSpec struct {
	runtime.ObjectTypedObject `json:",inline"`

	APISpecificationType string                    `json:"apiSpecificationType,omitempty"`
	APISpecVersion       string                    `json:"apiSpecificationVersion,omitempty"`
	Specification        *runtime.RawValue         `json:"specification,omitempty"`
	Artifact             *metav1.ResourceReference `json:"artifact,omitempty"`
}

func (s *ServiceSpec) GetVariant() identity.Variant {
	return nil
}

func (s *ServiceSpec) GetDependencies() []metav1.Dependency {
	return nil
}

func (s *ServiceSpec) Copy() internal.ServiceKindSpec {
	var spec *runtime.RawValue
	if s.Specification != nil {
		spec = generics.Pointer(s.Specification.Copy())
	}
	return &ServiceSpec{
		ObjectTypedObject:    runtime.NewTypedObject(s.GetType()),
		APISpecificationType: s.APISpecificationType,
		APISpecVersion:       s.APISpecVersion,
		Specification:        spec,
		Artifact:             s.Artifact.Copy(),
	}
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	r := *s
	if r.Specification != nil {
		r.Specification = generics.Pointer(s.Specification.Copy())
		r.Artifact = r.Artifact.Copy()
	}
	return &r
}

func (s *ServiceSpec) Validate(c internal.DescriptionContext) error {
	var list errors.ErrorList

	if s.APISpecificationType == "" {
		list.Add(fmt.Errorf("apiSpecificationType must be set"))
	}
	list.Add(c.ValidateResource(s.Artifact.AsResourceRef()))
	return list.Result()
}

func (s *ServiceSpec) GetReferences() crossref.References {
	return nil
}
