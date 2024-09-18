package internal

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/modeldesc"
	"sync"
)

type ServiceKindRegistry interface {
	Register(kind string, factory FactoryFunc)
	Create(m Model, descriptor *modeldesc.ServiceDescriptor) (ServiceVersionVariant, error)
}

var DefaultServiceKindRegistry = NewKindRegistry()

type FactoryFunc func(m Model, s *modeldesc.ServiceDescriptor) (ServiceVersionVariant, error)

type registry struct {
	lock      sync.Mutex
	factories map[string]FactoryFunc
}

func NewKindRegistry() ServiceKindRegistry {
	return &registry{
		factories: make(map[string]FactoryFunc),
	}
}

func (f *registry) Register(kind string, factory FactoryFunc) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.factories[kind] = factory
}

func (f *registry) Create(m Model, descriptor *modeldesc.ServiceDescriptor) (ServiceVersionVariant, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	fac := f.factories[descriptor.Kind.GetType()]
	if fac == nil {
		return nil, errors.ErrUnknown(modeldesc.KIND_SERVICE_TYPE, descriptor.Kind.GetType())
	}
	return fac(m, descriptor)
}
