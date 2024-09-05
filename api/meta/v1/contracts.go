package v1

type Contracts = List[Contract]

type Contract struct {
	Service     ServiceIdentity `json:"service"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Labels      Labels          `json:"labels"`
}

func (c Contract) Copy() *Contract {
	c.Labels = c.Labels.Copy()
	return &c
}
