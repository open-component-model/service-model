package v1

import (
	"github.com/open-component-model/service-model/api/utils"
)

type Installers = utils.CopyableList[Installer]

type Installer struct {
	Service     ServiceIdentity `json:"service"`
	Variant     Variant         `json:"variant,omitempty"`
	Version     string          `json:"version,omitempty"`
	Description string          `json:"description,omitempty"`
	Labels      Labels          `json:"labels,omitempty"`
}

func (i Installer) Copy() *Installer {
	i.Variant = i.Variant.Copy()
	i.Labels = i.Labels.Copy()
	return &i
}
