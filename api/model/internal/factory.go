package internal

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/modeldesc"
	"sync"
)

var DefaultFactory = NewFactory()

type FactoryFunc func(m Model, s *modeldesc.ServiceDescriptor) (ServiceVersionVariant, error)

type Factory struct {
	lock      sync.Mutex
	factories map[string]FactoryFunc
}

func NewFactory() *Factory {
	return &Factory{
		factories: make(map[string]FactoryFunc),
	}
}

func (f *Factory) Register(kind string, factory FactoryFunc) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.factories[kind] = factory
}

func (f *Factory) Create(m Model, descriptor *modeldesc.ServiceDescriptor) (ServiceVersionVariant, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	factory := f.factories[descriptor.Kind.GetType()]
	if factory == nil {
		return nil, errors.ErrUnknown(modeldesc.KIND_SERVICE_TYPE, descriptor.Kind.GetType())
	}
	return factory(m, descriptor)
}
