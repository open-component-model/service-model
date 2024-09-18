package internal

import (
	metav1 "github.com/open-component-model/service-model/api/identity"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	common "ocm.software/ocm/api/utils/misc"
)

type descriptionContext struct {
	Component string
	Version   string
}

type DescriptionContext interface {
	common.VersionedElement
	IsCanonical() bool
	MatchComponent(metav1.ServiceIdentity) bool

	ValidateResource(reference *ocmmeta.ResourceReference) error
	LookupService(string) *ServiceDescriptor
}

type ContextBuilder interface {
	DescriptionContext
	WithResourceValidator(validator ResourceValidator) ContextBuilder
}

type ResourceValidator func(reference *ocmmeta.ResourceReference) error

type _context struct {
	common.VersionedElement
	descriptor *ServiceModelDescriptor
	validator  ResourceValidator
}

func NewDescriptionContext(name, vers string, desc *ServiceModelDescriptor) ContextBuilder {
	return &_context{
		common.NewNameVersion(name, vers),
		desc,
		nil,
	}
}

func (c *_context) WithResourceValidator(v ResourceValidator) ContextBuilder {
	c.validator = v
	return c
}

func (c *_context) IsCanonical() bool {
	return c.descriptor.DocType.GetKind() == ABS_TYPE
}

func (c *_context) MatchComponent(s metav1.ServiceIdentity) bool {
	return s.Component() == "" || s.Component() == c.GetName()
}

func (c *_context) ValidateResource(r *ocmmeta.ResourceReference) error {
	if c.validator == nil || r == nil {
		return nil
	}
	return c.validator(r)
}

func (c *_context) LookupService(n string) *ServiceDescriptor {
	for i, s := range c.descriptor.Services {
		if s.Service.Name() == n && c.MatchComponent(s.Service) {
			return &c.descriptor.Services[i]
		}
	}
	return nil
}
