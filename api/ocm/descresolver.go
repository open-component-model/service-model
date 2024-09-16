package ocm

import (
	"sync"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/generics"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/api/ocm/resolvers"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/semverutils"
)

type serviceResolver struct {
	lock     sync.Mutex
	resolver ocm.ComponentVersionResolver
	services map[string]*modeldesc.ServiceDescriptor
	compvers map[common.NameVersion]error
}

func NewServiceResolver(r ocm.ComponentVersionResolver) modeldesc.Resolver {
	return &serviceResolver{
		resolver: r,
		services: map[string]*modeldesc.ServiceDescriptor{},
		compvers: map[common.NameVersion]error{},
	}
}

func (r *serviceResolver) LookupServiceVersionVariant(id identity.ServiceVersionVariantIdentity) (*modeldesc.ServiceDescriptor, error) {
	if id.IsConstraint() {
		return nil, errors.ErrInvalid(modeldesc.KIND_SERVICEVERSION, id.String())
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	cvid := common.NewNameVersion(id.Component(), id.Version())
	if err, ok := r.compvers[cvid]; !ok {
		err := r.addCV(id.ComponentVersionId())
		r.compvers[cvid] = err
		if err != nil {
			return nil, err
		}
	} else {
		if err != nil {
			return nil, err
		}
	}

	s := r.services[id.String()]
	if s == nil {
		return nil, errors.ErrNotFound(modeldesc.KIND_SERVICEVERSION, id.String())
	}
	return s, nil
}

func (r *serviceResolver) addCV(id common.NameVersion) error {
	desc, _, err := GetServiceModel(id.GetName(), id.GetVersion(), r.resolver)
	if err != nil {
		return err
	}
	for _, s := range desc.Services {
		r.services[s.GetId().String()] = generics.Pointer(s)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type compEntry struct {
	versions []string
	services map[string][]string
	err      error
}

func newComEntry() *compEntry {
	return &compEntry{
		services: map[string][]string{},
	}
}

type versionResolver struct {
	lock sync.Mutex
	modeldesc.Resolver
	resolver resolvers.ComponentResolver
	versions map[string]*compEntry
}

func NewVersionResolver(resolver resolvers.ComponentResolver) modeldesc.VersionResolver {
	return &versionResolver{
		Resolver: NewServiceResolver(resolvers.ComponentVersionResolverForComponentResolver(resolver)),
		resolver: resolver,
		versions: map[string]*compEntry{},
	}
}

func (r *versionResolver) ListVersions(id identity.ServiceIdentity, variant ...identity.Variant) ([]string, error) {
	if id.Component() == "" || id.Name() == "" {
		return nil, errors.ErrInvalid(modeldesc.KIND_SERVICEIDENTITY, id.String())
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	services := r.versions[id.Component()]
	if services == nil {
		services = newComEntry()
		r.versions[id.Component()] = services
		services.versions, services.err = resolvers.ListComponentVersions(id.Component(), r.resolver)
	}

	if services.err != nil {
		return nil, services.err
	}

	versions := services.services[id.Name()]
	if versions == nil {
		versions = []string{}
		for _, v := range services.versions {
			cand := identity.NewServiceVersionVariantIdentity(id, v, variant...)
			s, err := r.Resolver.LookupServiceVersionVariant(cand)
			if err != nil && !errors.IsErrNotFound(err) {
				return nil, err
			}
			if s != nil {
				versions = append(versions, v)
			}
		}
		semverutils.SortVersions(versions)
		services.services[id.Name()] = versions
	}

	return versions, nil
}
