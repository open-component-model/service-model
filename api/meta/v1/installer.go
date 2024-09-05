package v1

type Installers = List[Installer]

type Installer struct {
	Service     ServiceIdentity `json:"service"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Labels      Labels          `json:"labels"`
}

func (i Installer) Copy() *Installer {
	i.Labels = i.Labels.Copy()
	return &i
}
