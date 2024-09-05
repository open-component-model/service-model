package v1

import (
	"slices"
)

type ServiceInstances = List[ServiceInstance]

type ServiceInstance struct {
	Service  string       `json:"service"'`
	Versions []string     `json:"versions,omitempty"'`
	Dynamic  bool         `json:"dynamic"`
	Static   []StaticInfo `json:"static,omitempty"`
}

func (i ServiceInstance) Copy() *ServiceInstance {
	i.Versions = slices.Clone(i.Versions)
	i.Static = slices.Clone(i.Static)
	return &i
}

type StaticInfo struct {
	Name string `json:"name"`
}
