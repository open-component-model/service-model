package servicehdlr

import (
	"fmt"

	"github.com/mandelsoft/goutils/optionutils"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	modelocm "github.com/open-component-model/service-model/api/ocm"
	"ocm.software/ocm/api/ocm"
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
	return nil, nil
}

func (t *Services) Get(name utils.ElemSpec) ([]output.Object, error) {
	var id identity.ServiceVersionVariantIdentity
	err := id.Parse(name.String())
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
			l, err := t.get(identity.NewServiceVersionVariantIdentity(id.ServiceIdentity(), v, id.Variant()))
			if err != nil {
				return nil, err
			}
			result = append(result, l...)
		}
		return result, nil
	}
	return t.get(id)
}

func (t *Services) get(id identity.ServiceVersionVariantIdentity) ([]output.Object, error) {
	s, err := t.resolver.LookupServiceVersionVariant(id)
	if err != nil {
		return nil, err
	}
	return sliceutils.AsSlice(output.Object(NewObject(nil, s))), nil
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
