package v1

type Installers = CopyableList[Installer]

type Installer struct {
	Service     ServiceIdentity `json:"service"`
	Version     string          `json:"version,omitempty"`
	Description string          `json:"description,omitempty"`
	Labels      Labels          `json:"labels,omitempty"`
}

func (i Installer) Copy() *Installer {
	i.Labels = i.Labels.Copy()
	return &i
}
