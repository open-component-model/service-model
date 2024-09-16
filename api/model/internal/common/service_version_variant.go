package common

import (
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
)

type ServiceVersionVariant struct {
	id   identity.ServiceVersionVariantIdentity
	spec *modeldesc.ServiceDescriptor
}

func New(spec *modeldesc.ServiceDescriptor) *ServiceVersionVariant {
	return &ServiceVersionVariant{
		id:   identity.NewServiceVersionVariantIdentity(spec.Service, spec.Version, spec.Kind.GetVariant()),
		spec: spec.Copy(),
	}
}

func (s *ServiceVersionVariant) GetType() string {
	return s.spec.Kind.GetType()
}

func (s *ServiceVersionVariant) GetName() string {
	return s.spec.Service.Name()
}

func (s *ServiceVersionVariant) GetComponent() string {
	return s.spec.Service.Component()
}

func (s *ServiceVersionVariant) GetIdentity() identity.ServiceVersionVariantIdentity {
	return s.id
}

func (s *ServiceVersionVariant) GetVersion() string {
	return s.spec.Version
}

func (s *ServiceVersionVariant) GetVariant() identity.Variant {
	return s.spec.Kind.GetVariant()
}
