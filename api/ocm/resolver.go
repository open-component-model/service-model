package ocm

import (
	"sync"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/desc2model"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/model"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/ocm"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/semverutils"
)

type Resolver struct {
	lock       sync.Mutex
	resolver   ocm.ComponentResolver
	services   map[metav1.ServiceVersionIdentity]*modeldesc.ServiceDescriptor
	components map[string][]string
}

var _ model.Resolver = (*Resolver)(nil)

func NewResolver(r ocm.ComponentResolver) model.Resolver {
	return &Resolver{
		resolver:   r,
		services:   map[metav1.ServiceVersionIdentity]*modeldesc.ServiceDescriptor{},
		components: map[string][]string{},
	}
}

func (r *Resolver) LookupService(id metav1.ServiceVersionIdentity) (model.Service, error) {
	if id.IsConstraint() {
		return nil, errors.ErrInvalid(modeldesc.KIND_SERVICEVERSION, id.String())
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	s := r.services[id]
	if s != nil {
		return s, nil
	}

	err := r.addCV(id.ComponentVersionId())
	if err != nil {
		return nil, err
	}
	s = r.services[id]
	if s == nil {
		return nil, errors.ErrNotFound(modeldesc.KIND_SERVICEVERSION, id.String())
	}
	return desc2model.ServiceDescriptorToModelService(s)
}

func (r *Resolver) addCV(id common.NameVersion) error {
	desc, _, err := GetServiceModel(id.GetName(), id.GetVersion(), r.resolver)
	if err != nil {
		return err
	}
	for _, s := range desc.Services {
		r.services[s.GetId()] = generics.Pointer(s)
	}
	return nil
}

func (r *Resolver) ListVersions(id metav1.ServiceIdentity) ([]string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var versions []string
	if list, ok := r.components[id.Component]; ok {
		return list, nil
	}
	for _, res := range r.resolver.LookupComponentProviders(id.Component) {
		c, err := res.LookupComponent(id.Component)
		if err != nil || c == nil {
			continue
		}
		list, err := c.ListVersions()
		if err != nil {
			continue
		}
		sliceutils.AppendUnique(versions, list...)
	}
	semverutils.SortVersions(versions)
	return versions, nil
}
