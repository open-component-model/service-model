package v1

import (
	"github.com/mandelsoft/goutils/general"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
	"ocm.software/ocm/api/utils/runtime"
)

const VERSION = "v1"

func init() {
	modeldesc.RegisterVersion(NewVersion(modeldesc.REL_TYPE))
	modeldesc.RegisterVersion(NewVersion(modeldesc.ABS_TYPE))

	modeldesc.RegisterVersion(NewVersion(modeldesc.REL_TYPE, ""))
	modeldesc.RegisterVersion(NewVersion(modeldesc.ABS_TYPE, ""))
}

func NewVersion(base string, vers ...string) modeldesc.Version {
	version := VERSION
	if len(vers) > 0 {
		version = general.Optional(vers...)
	}
	return vpi.NewVersion[ServiceModelDescriptor](runtime.TypeName(base, version), Converter{})
}
