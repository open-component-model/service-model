package modeldesc

import (
	"github.com/open-component-model/service-model/api/modeldesc/internal"
	"ocm.software/ocm/api/utils/runtime"
)

const KIND_SERVICE_TYPE = internal.KIND_SERVICE_TYPE
const KIND_MODELVERSION = internal.KIND_MODELVERSION

const REL_TYPE = internal.REL_TYPE
const ABS_TYPE = internal.ABS_TYPE

type (
	ServiceKindSpec        = internal.ServiceKindSpec
	ServiceDescriptor      = internal.ServiceDescriptor
	ServiceModelDescriptor = internal.ServiceModelDescriptor
)

func Decode(data []byte, unmarshaller ...runtime.Unmarshaler) (*ServiceModelDescriptor, error) {
	return internal.Decode(data, unmarshaller...)
}

func Encode(desc *ServiceModelDescriptor, marshaller ...runtime.Marshaler) ([]byte, error) {
	return internal.Encode(desc, marshaller...)
}
