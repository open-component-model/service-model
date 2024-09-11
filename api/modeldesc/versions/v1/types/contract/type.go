package contract

import (
	"github.com/mandelsoft/goutils/generics"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	me "github.com/open-component-model/service-model/api/modeldesc/types/contract"
	"github.com/open-component-model/service-model/api/modeldesc/versions/v1/types"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
	"ocm.software/ocm/api/utils/runtime"
)

func init() {
	types.RegisterServiceType(vpi.NewServiceKindType[ServiceSpec](me.TYPE, Converter{}))
}

type ServiceSpec struct {
	runtime.ObjectTypedObject `json:",inline"`

	APISpecificationType string                    `json:"apiSpecificationType,omitempty"`
	APISpecVersion       string                    `json:"apiSpecificationVersion,omitempty"`
	Specification        *runtime.RawValue         `json:"specification,omitempty"`
	Artifact             *metav1.ResourceReference `json:"artifact,omitempty"`
}

type Converter struct{}

func (c Converter) ConvertFrom(object modeldesc.ServiceKindSpec) (vpi.ServiceKindSpec, error) {
	in := object.(*me.ServiceSpec)
	var spec *runtime.RawValue
	if in.Specification != nil {
		spec = generics.Pointer(in.Specification.Copy())
	}
	return &ServiceSpec{
		ObjectTypedObject:    in.ObjectTypedObject,
		APISpecificationType: in.APISpecificationType,
		APISpecVersion:       in.APISpecVersion,
		Specification:        spec,
		Artifact:             in.Artifact.Copy(),
	}, nil
}

func (c Converter) ConvertTo(object vpi.ServiceKindSpec) (modeldesc.ServiceKindSpec, error) {
	in := object.(*ServiceSpec)
	var spec *runtime.RawValue
	if in.Specification != nil {
		spec = generics.Pointer(in.Specification.Copy())
	}
	return &me.ServiceSpec{
		ObjectTypedObject:    in.ObjectTypedObject,
		APISpecificationType: in.APISpecificationType,
		APISpecVersion:       in.APISpecVersion,
		Specification:        spec,
		Artifact:             in.Artifact.Copy(),
	}, nil
}
