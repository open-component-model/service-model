package v1

import (
	"slices"

	"github.com/open-component-model/service-model/api/utils"
)

type ServiceInstances = utils.CopyableList[ServiceInstance]

type ServiceInstance struct {
	Service  ServiceIdentity `json:"service"'`
	Versions []string        `json:"versions,omitempty"'`
	Dynamic  bool            `json:"dynamic"`
	Static   []StaticInfo    `json:"static,omitempty"`
}

func (i ServiceInstance) Copy() *ServiceInstance {
	i.Versions = slices.Clone(i.Versions)
	i.Static = slices.Clone(i.Static)
	return &i
}

type StaticInfo struct {
	Name string `json:"name"`
}
