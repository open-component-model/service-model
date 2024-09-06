package v1

type Contracts = CopyableList[Contract]

type Contract struct {
	Service ServiceIdentity `json:"service"`
	Version string          `json:"version,omitempty"`
	Labels  Labels          `json:"labels,omitempty"`
}

func (c Contract) Copy() *Contract {
	c.Labels = c.Labels.Copy()
	return &c
}
