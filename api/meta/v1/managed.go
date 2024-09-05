package v1

import (
	"slices"
)

type ManagedServices = List[ManagedService]

type ManagedService struct {
	Service               ServiceIdentity       `json:"service"`
	Versions              []string              `json:"versions,omitempty"`
	Labels                Labels                `json:"labels,omitempty"`
	DependencyResolutions DependencyResolutions `json:"dependencyResolutions,omitempty"`
}

func (s ManagedService) Copy() *ManagedService {
	s.Labels = s.Labels.Copy()
	s.Versions = slices.Clone(s.Versions)
	s.DependencyResolutions = s.DependencyResolutions.Copy()
	return &s
}

type DependencyResolutions = List[DependencyResolution]

type DependencyResolution struct {
	Managed    bool `json:"managed"`
	Configured bool `json:"configured"`
}

func (d DependencyResolution) Copy() *DependencyResolution {
	return &d
}
