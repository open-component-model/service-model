package vpi

// this file is similar to contexts oci.

import (
	"github.com/mandelsoft/goutils/errors"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	"ocm.software/ocm/api/utils/runtime"
)

type ServiceKindSpec interface {
	runtime.TypedObject
}

type ServiceKindType interface {
	runtime.TypedObjectType[ServiceKindSpec]
	Converter[modeldesc.ServiceKindSpec, ServiceKindSpec]
}

type ServiceKindTypeScheme interface {
	runtime.TypeScheme[ServiceKindSpec, ServiceKindType]
	Converter[modeldesc.ServiceKindSpec, ServiceKindSpec]
}

type serviceKindTypeScheme struct {
	runtime.TypeScheme[ServiceKindSpec, ServiceKindType]
}

func NewServiceKindTypeScheme() ServiceKindTypeScheme {
	return &serviceKindTypeScheme{runtime.NewTypeScheme[ServiceKindSpec, ServiceKindType]()}
}

func (s *serviceKindTypeScheme) ConvertFrom(spec modeldesc.ServiceKindSpec) (ServiceKindSpec, error) {
	t := s.GetType(spec.GetType())
	if t == nil {
		return nil, errors.ErrUnknown(modeldesc.KIND_SERVICE_TYPE, spec.GetType())
	}
	return t.ConvertFrom(spec)
}

func (s *serviceKindTypeScheme) ConvertTo(spec ServiceKindSpec) (modeldesc.ServiceKindSpec, error) {
	t := s.GetType(spec.GetType())
	if t == nil {
		return nil, errors.ErrUnknown(modeldesc.KIND_SERVICE_TYPE, spec.GetType())
	}
	return t.ConvertTo(spec)
}

////////////////////////////////////////////////////////////////////////////////

type serviceKindType struct {
	runtime.TypedObjectType[ServiceKindSpec]
	Converter[modeldesc.ServiceKindSpec, ServiceKindSpec]
}

type serviceKindPointer[I any] interface {
	*I
	ServiceKindSpec
}

func NewServiceKindType[I any, P serviceKindPointer[I]](name string, converter Converter[modeldesc.ServiceKindSpec, ServiceKindSpec]) ServiceKindType {
	var proto I

	return &serviceKindType{
		TypedObjectType: runtime.NewTypedObjectTypeByProto[ServiceKindSpec](name, P(&proto)),
		Converter:       converter,
	}
}
