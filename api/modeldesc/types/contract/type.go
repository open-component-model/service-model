package installer

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/generics"
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

func (s *ServiceSpec) GetReferences() internal.References {
	return nil
}
