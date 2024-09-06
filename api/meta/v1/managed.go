package v1

import (
	"slices"

	"github.com/open-component-model/service-model/api/utils"
)

const (
	DEPRES_MANGED     = "managed"
	DEPRES_CONFIGURED = "configured"

	DEPUSE_EXCLUSIVE  = "exclusive"
	DEPUSE_SHARED     = "shared"
	DEPUSE_CONFIGURED = "configured"
)

type ManagedServices = utils.CopyableList[ManagedService]

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

type DependencyResolutions = utils.CopyableList[DependencyResolution]

type DependencyResolution struct {
	Name       string `json:"name"`
	Resolution string `json:"resolution"`
	Usage      string `json:"usage"`
	Labels     Labels `json:"labels,omitempty"`
}

func (d DependencyResolution) Copy() *DependencyResolution {
	d.Labels = d.Labels.Copy()
	return &d
}
