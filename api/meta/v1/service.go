package v1

import (
	"github.com/open-component-model/service-model/api/identity"
	"ocm.software/ocm/api/utils/runtime"
)

type CommonServiceSpec struct {
	Service     identity.ServiceIdentity `json:"service"`
	Version     string                   `json:"version,omitempty"`
	ShortName   string                   `json:"shortName"`
	Description string                   `json:"description,omitempty"`
	Labels      Labels                   `json:"labels,omitempty"`
}

func (s *CommonServiceSpec) GetId() identity.ServiceVersionIdentity {
	return identity.NewServiceVersionId(s.Service, s.Version)
}

func (s *CommonServiceSpec) Copy() *CommonServiceSpec {
	c := *s
	c.Labels = s.Labels.Copy()
	return &c
}

type CommonServiceImplementationSpec struct {
	runtime.ObjectTypedObject `json:",inline"`
	Variant                   identity.Variant `json:"variant,omitempty"`
	Abstract                  bool             `json:"abstract,omitempty"`
	InheritFrom               identity.Variant `json:"inheritFrom,omitempty"`
	Dependencies              Dependencies     `json:"dependencies,omitempty"`
	Contracts                 Contracts        `json:"contracts,omitempty"`
}

func (s *CommonServiceImplementationSpec) GetVariant() identity.Variant {
	return s.Variant
}

func (c CommonServiceImplementationSpec) Copy() *CommonServiceImplementationSpec {
	c.Dependencies = c.Dependencies.Copy()
	c.Variant = c.Variant.Copy()
	c.InheritFrom = c.InheritFrom.Copy()
	c.Contracts = c.Contracts.Copy()
	return &c
}

type CommonConsumerServiceImplementationSpec struct {
	CommonServiceImplementationSpec `json:",inline"`
	External                        bool       `json:"external,omitempty"`
	Installers                      Installers `json:"installers,omitempty"`
}

func (c CommonConsumerServiceImplementationSpec) Copy() *CommonConsumerServiceImplementationSpec {
	c.CommonServiceImplementationSpec = *c.CommonServiceImplementationSpec.Copy()
	c.Installers = c.Installers.Copy()
	return &c
}
