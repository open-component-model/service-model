package servicehdlr

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/goutils/finalizer"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/goutils/optionutils"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/provider"
	modelocm "github.com/open-component-model/service-model/api/ocm"
	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/api/utils/refmgmt"
	"ocm.software/ocm/api/utils/semverutils"
	"ocm.software/ocm/cmds/ocm/common/output"
	"ocm.software/ocm/cmds/ocm/common/utils"
)

type Services struct {
	opts     *Options
	resolver modeldesc.VersionResolver
}

func ForServices(resolver ocm.ComponentResolver, opts ...Option) utils.TypeHandler {
	return forServices(resolver, opts...)
}

func forServices(resolver ocm.ComponentResolver, opts ...Option) *Services {
	t := &Services{
		opts:     optionutils.EvalArbitraryOptions[Option, Options](opts...),
		resolver: modelocm.NewVersionResolver(resolver),
	}
	return t
}

func (t *Services) All() ([]output.Object, error) {
	if t.opts.repo == nil {
		return nil, nil
	}

	result := []output.Object{}
	lister := t.opts.repo.ComponentLister()
	if lister == nil {
		return nil, fmt.Errorf("repository does not support listing components")
	}

	var finalize finalizer.Finalizer
	defer finalize.Finalize()

	list, err := lister.GetComponents("", true)
	if err != nil {
		return nil, err
	}
	for _, n := range list {
		loop := finalize.Nested()
		c, err := refmgmt.ToLazy(t.opts.repo.LookupComponent(n))
		if err != nil {
			return nil, err
		}
		loop.Close(c)
		versions, err := c.ListVersions()
		if err != nil {
			return nil, err
		}
		for _, v := range versions {
			cv, err := refmgmt.ToLazy(c.LookupVersion(v))
			if err != nil {
				continue
			}

			loop := loop.Nested()
			loop.Close(cv)

			m, _, err := modelocm.GetServiceModelFromCV(cv)
			if err != nil {
				return nil, err
			}
			for _, s := range m.Services {
				t.Add(&result, NewObject(nil, "", generics.Pointer(s)))
			}
			loop.Finalize()
		}
		loop.Finalize()
	}
	return result, nil
}

func (t *Services) Get(spec utils.ElemSpec) ([]output.Object, error) {
	var managed *string

	name := spec.String()

	if idx := strings.Index(name, "@"); idx >= 0 {
		managed = generics.Pointer(name[:idx])
		name = name[idx+1:]
	}

	var id identity.ServiceVersionVariantIdentity
	err := id.Parse(name)
	if err != nil {
		return nil, err
	}

	if id.Version() == "" {
		versions, err := t.resolver.ListVersions(id.ServiceIdentity(), id.Variant())
		if err != nil {
			return nil, err
		}

		versions, err = t.FilterVersions(versions)
		if err != nil {
			return nil, err
		}

		var result []output.Object
		for _, v := range versions {
			l, err := t.get(identity.NewServiceVersionVariantId(id.ServiceIdentity(), v, id.Variant()), managed)
			if err != nil {
				return nil, err
			}
			result = append(result, l...)
		}
		return result, nil
	}
	return t.get(id, managed)
}

func (t *Services) get(id identity.ServiceVersionVariantIdentity, managed *string) ([]output.Object, error) {
	s, err := t.resolver.LookupServiceVersionVariant(id)
	if err != nil {
		return nil, err
	}
	if managed == nil {
		obj := NewObject(nil, "", s)
		if t.opts.state != nil {
			t.opts.state.Add(obj.Element)
		}
		var result []output.Object
		t.Add(&result, obj)
		return result, nil
	}
	if s.GetType() == provider.TYPE {
		var result []output.Object
		p := s.Kind.(*provider.ServiceSpec)
		for _, m := range p.ManagedServices {
			if *managed == "" || m.Name == *managed {
				for _, v := range m.Versions {
					id := identity.NewServiceVersionVariantId(m.Service, v, m.Variant)
					list, err := t.get(id, nil)
					if err != nil {
						result = append(result, NewErrorObject(err, nil, modeldesc.DEP_MANAGES, m.Service, v, m.Variant))
					} else {
						result = append(result, list...)
					}
				}
			}
		}
		return result, nil
	}
	return nil, fmt.Errorf("%s is no service provider", id)
}

func (t *Services) Add(list *[]output.Object, objs ...*Object) {
	for _, o := range objs {
		if t.opts.state != nil && o.Element != nil {
			t.opts.state.Add(o.Element)
		}
		*list = append(*list, o)
	}
}

func (t *Services) FilterVersions(vers []string) ([]string, error) {
	latest := optionutils.AsBool(t.opts.latestOnly)
	if len(t.opts.constraints) == 0 && !latest {
		return vers, nil
	}
	versions, err := semverutils.MatchVersionStrings(vers, t.opts.constraints...)
	if err != nil {
		return nil, fmt.Errorf("invalid constraints: %w", err)
	}
	if latest && len(versions) > 1 {
		versions = versions[len(versions)-1:]
	}
	vers = nil
	for _, v := range versions {
		vers = append(vers, v.Original())
	}
	return vers, nil
}

func (t *Services) Close() error {
	return nil
}
