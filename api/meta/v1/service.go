package v1

import (
	"ocm.software/ocm/api/utils/runtime"
)

type BaseServiceSpec struct {
	Service     string  `json:"service"`
	ShortName   string  `json:"shortName"`
	Description string  `json:"description"`
	Labels      Labels  `json:"labels"`
	Variant     Variant `json:"variant"`
	Abstract    bool    `json:"abstract"`
	InheritFrom Variant `json:"inheritFrom"`
}

func (s *BaseServiceSpec) Copy() *BaseServiceSpec {
	c := *s
	c.Labels = s.Labels.Copy()
	c.Variant = s.Variant.Copy()
	c.InheritFrom = s.InheritFrom.Copy()
	return &c
}

type CommonServiceImplementationSpec struct {
	runtime.ObjectTypedObject
	External     bool         `json:"external"`
	Dependencies Dependencies `json:"dependencies,omitempty"`
	Contracts    Contracts    `json:"contracts,omitempty"`
	Installers   Installers   `json:"installers,omitempty"`
}

func (s *CommonServiceImplementationSpec) Copy() *CommonServiceImplementationSpec {
	c := *s
	c.Dependencies = s.Dependencies.Copy()
	c.Contracts = s.Contracts.Copy()
	c.Installers = s.Installers.Copy()
	return &c
}
