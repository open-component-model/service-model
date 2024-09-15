package v1

import (
	"github.com/mandelsoft/goutils/sliceutils"
	"ocm.software/ocm/api/utils/runtime"
)

type CommonServiceSpec struct {
	Service     ServiceIdentity `json:"service"`
	Version     string          `json:"version,omitempty"`
	ShortName   string          `json:"shortName"`
	Description string          `json:"description,omitempty"`
	Labels      Labels          `json:"labels,omitempty"`
}

func (s *CommonServiceSpec) GetId() ServiceVersionIdentity {
	return NewServiceVersionId(s.Service, s.Version)
}

func (s *CommonServiceSpec) Copy() *CommonServiceSpec {
	c := *s
	c.Labels = s.Labels.Copy()
	return &c
}

type CommonServiceImplementationSpec struct {
	runtime.ObjectTypedObject `json:",inline"`
	Variant                   Variant      `json:"variant,omitempty"`
	Abstract                  bool         `json:"abstract,omitempty"`
	InheritFrom               Variant      `json:"inheritFrom,omitempty"`
	Dependencies              Dependencies `json:"dependencies,omitempty"`
	Contracts                 Contracts    `json:"contracts,omitempty"`
}

func (s *CommonServiceImplementationSpec) GetVariant() Variant {
	return s.Variant
}

func (s *CommonServiceImplementationSpec) GetDependencies() []Dependency {
	return s.Dependencies.Copy()
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

func (c CommonConsumerServiceImplementationSpec) GetDependencies() []Dependency {
	deps := c.CommonServiceImplementationSpec.GetDependencies()
	for _, e := range c.Installers {
		deps = append(deps, Dependency{
			Name:               "",
			Service:            e.Service,
			Variant:            e.Variant.Copy(),
			Kind:               DEPKIND_INSTALLER,
			VersionConstraints: sliceutils.AsSlice(e.Version),
			ServiceInstances:   nil,
			Optional:           false,
			Description:        e.Description,
			Labels:             e.Labels.Copy(),
		})
	}
	return deps
}

func (c CommonConsumerServiceImplementationSpec) Copy() *CommonConsumerServiceImplementationSpec {
	c.CommonServiceImplementationSpec = *c.CommonServiceImplementationSpec.Copy()
	c.Installers = c.Installers.Copy()
	return &c
}
