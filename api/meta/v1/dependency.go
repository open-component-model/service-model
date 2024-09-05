package v1

import (
	"slices"
)

type Dependencies = List[Dependency]

type Dependency struct {
	Name               string           `json:"name"`
	Service            ServiceIdentity  `json:"service"`
	Kind               string           `json:"kind"`
	VersionConstraints []string         `json:"versionConstraints,omitempty"`
	ServiceInstances   ServiceInstances `json:"serviceInstances,omitempty"`
	Optional           bool             `json:"optional"`
	Description        string           `json:"description"`
	Labels             Labels           `json:"labels,omitempty"`
}

func (d Dependency) Copy() *Dependency {
	d.VersionConstraints = slices.Clone(d.VersionConstraints)
	d.Labels = d.Labels.Copy()
	return &d
}
