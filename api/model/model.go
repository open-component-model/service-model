package model

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/model/internal"
	"sync"
)

type model struct {
	lock     sync.RWMutex
	registry internal.ServiceKindRegistry
	resolver Resolver
	store    *serviceStore
}

func NewModel(resolver Resolver) Model {
	return &model{
		registry: internal.DefaultServiceKindRegistry,
		resolver: resolver,
		store:    newServiceStore(),
	}
}

func (m *model) GetServiceKindRegistry() internal.ServiceKindRegistry {
	return m.registry
}

func (m *model) GetServiceVersionVariant(id identity.ServiceVersionVariantIdentity) (ServiceVersionVariant, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var err error
	s, ok := m.store.GetService(id)
	if !ok && s == nil {
		s, err = m.resolver.LookupServiceVersionVariant(m, id)
	}
	if err != nil {
		if errors.IsErrNotFound(err) {
			m.store.SetUnknownService(id)
		}
		return nil, err
	}
	m.store.SetService(s)
	return s, nil
}
