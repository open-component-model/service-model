package v1

import (
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
	"ocm.software/ocm/api/utils/runtime"
)

const VERSION = "v1"

func init() {
	modeldesc.RegisterVersion(NewVersion(modeldesc.REL_TYPE))
	modeldesc.RegisterVersion(NewVersion(modeldesc.ABS_TYPE))
}

func NewVersion(base string) modeldesc.Version {
	return vpi.NewVersion[ServiceModelDescriptor](runtime.TypeName(base, VERSION), Converter{})
}
