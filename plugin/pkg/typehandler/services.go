package typehandler

import (
	"github.com/mandelsoft/goutils/sliceutils"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc"
	modelocm "github.com/open-component-model/service-model/api/ocm"
	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/api/ocm/resolvers"
	"ocm.software/ocm/cmds/ocm/common/output"
	"ocm.software/ocm/cmds/ocm/common/utils"
)

type Services struct {
	resolver ocm.ComponentResolver
	services modeldesc.Resolver
}

func ForServices(resolver ocm.ComponentResolver) (utils.TypeHandler, error) {
	t := &Services{
		resolver: resolver,
		services: modelocm.NewServiceResolver(resolvers.ComponentVersionResolverForComponentResolver(resolver)),
	}
	return t, nil
}

func (t *Services) All() ([]output.Object, error) {
	return nil, nil
}

func (t *Services) Get(name utils.ElemSpec) ([]output.Object, error) {
	var id v1.ServiceVersionVariantIdentity
	err := id.Parse(name.String())
	if err != nil {
		return nil, err
	}
	s, err := t.services.LookupServiceVersionVariant(id)
	if err != nil {
		return nil, err
	}
	return sliceutils.AsSlice(output.Object(NewObject(nil, s))), nil
}

func (t *Services) Close() error {
	return nil
}
