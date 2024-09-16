package v1

import (
	"slices"

	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/utils"
)

const (
	DEPKIND_IMPLEMENTATION = "implementation"
	DEPKIND_ORCHESTRATION  = "orchestration"

	// DEPKIND_INSTALLER is the dependency to an installer.
	// It is maintained as explicit attribute and not part of regular
	// dependencies, because they feature a restricted attribute set.
	DEPKIND_INSTALLER = "installer"
)

type Dependencies = utils.CopyableList[Dependency]

type Dependency struct {
	Name               string                   `json:"name"`
	Service            identity.ServiceIdentity `json:"service"`
	Variant            identity.Variant         `json:"variant,omitempty"`
	Kind               string                   `json:"kind"`
	VersionConstraints []string                 `json:"versionConstraints,omitempty"`
	ServiceInstances   ServiceInstances         `json:"serviceInstances,omitempty"`
	Optional           bool                     `json:"optional,omitempty"`
	Description        string                   `json:"description,omitempty"`
	Labels             Labels                   `json:"labels,omitempty"`
}

func (d Dependency) Copy() *Dependency {
	d.VersionConstraints = slices.Clone(d.VersionConstraints)
	d.Labels = d.Labels.Copy()
	d.Variant = d.Variant.Copy()
	return &d
}
