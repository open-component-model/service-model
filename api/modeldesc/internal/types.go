package internal

import (
	v1 "github.com/open-component-model/service-model/api/meta/v1"
	common "ocm.software/ocm/api/utils/misc"
)

type descriptionContext struct {
	Component string
	Version   string
}

type DescriptionContext interface {
	common.VersionedElement
	IsCanonical() bool
	MatchComponent(v1.ServiceIdentity) bool
	LookupService(string) *ServiceDescriptor
}

type _context struct {
	common.VersionedElement
	descriptor *ServiceModelDescriptor
}

func NewDescriptionContext(name, vers string, desc *ServiceModelDescriptor) DescriptionContext {
	return &_context{
		common.NewNameVersion(name, vers),
		desc,
	}
}

func (c *_context) IsCanonical() bool {
	return c.descriptor.DocType == ABS_TYPE
}

func (c *_context) MatchComponent(s v1.ServiceIdentity) bool {
	return s.Component == "" || s.Component == c.GetName()
}

func (c *_context) LookupService(n string) *ServiceDescriptor {
	for i, s := range c.descriptor.Services {
		if s.Service.Name == n && c.MatchComponent(s.Service) {
			return &c.descriptor.Services[i]
		}
	}
	return nil
}
