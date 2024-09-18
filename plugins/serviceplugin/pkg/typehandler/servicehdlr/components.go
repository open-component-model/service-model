package servicehdlr

import (
	"strings"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/generics"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	modelocm "github.com/open-component-model/service-model/api/ocm"
	"github.com/open-component-model/service-model/api/utils"
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/api/ocm/resolvers"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/handlers/comphdlr"
	"ocm.software/ocm/cmds/ocm/common/output"
	cmdutils "ocm.software/ocm/cmds/ocm/common/utils"
)

type Components struct {
	cvresolver resolvers.ComponentVersionResolver
	resolver   modeldesc.VersionResolver
	*Services

	components []*comphdlr.Object
}

var _ cmdutils.TypeHandler = (*Components)(nil)

func ForComponents(octx cli.OCM, resolver ocm.ComponentResolver, oopts *output.Options, repobase ocm.Repository, session ocm.Session, compspecs []string, hopts ...Option) (*Components, error) {
	components, err := comphdlr.Evaluate(octx, session, repobase, compspecs, oopts, MapToCompHandlerOptions(hopts...)...)
	if err != nil {
		return nil, err
	}

	t := &Components{
		cvresolver: resolvers.ComponentVersionResolverForComponentResolver(resolver),
		resolver:   modelocm.NewVersionResolver(resolver),
		Services:   forServices(resolver, hopts...),
		components: components,
	}
	return t, nil
}

func (t *Components) GetResolver() modeldesc.VersionResolver {
	return t.resolver
}

func (t *Components) All() ([]output.Object, error) {
	var result []output.Object

	for _, c := range t.components {
		m, _, err := modelocm.GetServiceModelFromCV(c.ComponentVersion, t.cvresolver)
		if err != nil {
			return nil, err
		}
		for _, s := range m.Services {
			t.Add(&result, NewObject(nil, "", generics.Pointer(s)))
		}
	}
	return result, nil
}

func (t *Components) Get(spec cmdutils.ElemSpec) ([]output.Object, error) {
	name := spec.String()

	idx := strings.Index(spec.String(), ":")
	if idx >= 0 {
		return nil, errors.ErrInvalid(modeldesc.KIND_SERVICEIDENTITY, name)
	}

	id, err := utils.Parse[identity.ServiceVariantIdentity](name)
	if err != nil {
		return nil, errors.Wrapf(err, "service variant %q", name)
	}

	var result []output.Object
	for _, c := range t.components {
		if id.Component() != "" && c.ComponentVersion.GetName() != id.Component() {
			continue
		}
		svid := identity.NewServiceVersionVariantId(identity.NewServiceId(c.ComponentVersion.GetName(), id.Name()), c.ComponentVersion.GetVersion(), id.Variant())
		s, err := t.Services.get(svid)
		if err != nil && !errors.IsErrNotFound(err) {
			return nil, err
		}
		result = append(result, s...)
	}
	return result, nil
}
