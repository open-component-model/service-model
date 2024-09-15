package services

import (
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/api/ocm"
)

type OCM struct {
	ctx ocm.Context
}

func (o *OCM) Context() ocm.Context {
	return o.ctx
}

func (o *OCM) OpenCTF(path string) (ocm.Repository, error) {
	panic("not supported")
}

var _ cli.OCM = (*OCM)(nil)

func NewOCM(ctx ocm.Context) cli.OCM {
	return &OCM{ctx}
}
