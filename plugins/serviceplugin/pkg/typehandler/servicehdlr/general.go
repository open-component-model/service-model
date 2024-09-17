package servicehdlr

import (
	"github.com/mandelsoft/goutils/optionutils"
	"github.com/open-component-model/service-model/api/modeldesc"
)

func ForVersionResolver(resolver modeldesc.VersionResolver, opts ...Option) *Services {
	t := &Services{
		opts:     optionutils.EvalArbitraryOptions[Option, Options](opts...),
		resolver: resolver,
	}
	return t
}
