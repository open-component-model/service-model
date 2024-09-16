package v1

import (
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/utils"
)

type Contracts = utils.CopyableList[Contract]

type Contract struct {
	Service     identity.ServiceIdentity `json:"service"`
	Version     string                   `json:"version,omitempty"`
	Description string                   `json:"description,omitempty"`
	Labels      Labels                   `json:"labels,omitempty"`
}

func (c Contract) Copy() *Contract {
	c.Labels = c.Labels.Copy()
	return &c
}
