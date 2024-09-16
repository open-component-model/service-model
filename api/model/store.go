package model

import (
	"github.com/mandelsoft/goutils/maputils"
	"github.com/open-component-model/service-model/api/identity"
)

type serviceStore struct {
	Services map[identity.ServiceIdentity]map[string]map[string]ServiceVersionVariant
}

func newServiceStore() *serviceStore {
	return &serviceStore{Services: make(map[identity.ServiceIdentity]map[string]map[string]ServiceVersionVariant)}
}

func (c *serviceStore) SetService(s ServiceVersionVariant) {
	c.setService(s.GetIdentity(), s)
}

func (c *serviceStore) SetUnknownService(id identity.ServiceVersionVariantIdentity) {
	c.setService(id, nil)
}
func (c *serviceStore) setService(svvid identity.ServiceVersionVariantIdentity, s ServiceVersionVariant) {
	id := svvid
	versions := c.Services[id.ServiceIdentity()]
	if versions == nil {
		versions = map[string]map[string]ServiceVersionVariant{}
		c.Services[id.ServiceIdentity()] = versions
	}
	variants := versions[id.Version()]
	if variants == nil {
		variants = map[string]ServiceVersionVariant{}
		versions[id.Version()] = variants
	}
	variants[id.Variant().String()] = s
}

func (c *serviceStore) GetService(id identity.ServiceVersionVariantIdentity) (ServiceVersionVariant, bool) {
	versions := c.Services[id.ServiceIdentity()]
	if versions == nil {
		return nil, false
	}
	variants := versions[id.Version()]
	if variants == nil {
		return nil, false
	}
	s, ok := variants[id.Variant().String()]
	return s, ok
}

func (c *serviceStore) GetServiceVariants(id identity.ServiceVersionIdentity) []ServiceVersionVariant {
	versions := c.Services[id.ServiceIdentity()]
	if versions == nil {
		return nil
	}
	variants := versions[id.Version()]
	if variants == nil {
		return nil
	}
	return maputils.Values(variants)
}
