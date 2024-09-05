package modeldesc

import (
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

const KIND_SERVICE_TYPE = internal.KIND_SERVICE_TYPE
const KIND_MODELVERSION = internal.KIND_SERVICE_TYPE

const REL_TYPE = internal.REL_TYPE
const ABS_TYPE = internal.ABS_TYPE

type (
	ServiceKindSpec        = internal.ServiceKindSpec
	ServiceDescriptor      = internal.ServiceDescriptor
	ServiceModelDescriptor = internal.ServiceModelDescriptor
)
