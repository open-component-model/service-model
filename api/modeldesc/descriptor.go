package modeldesc

import (
	"github.com/open-component-model/service-model/api/modeldesc/internal"
	v1 "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/runtime"
)

const (
	KIND_SERVICEVERSION   = internal.KIND_SERVICEVERSION
	KIND_SERVICE_TYPE     = internal.KIND_SERVICE_TYPE
	KIND_MODELVERSION     = internal.KIND_MODELVERSION
	KIND_DESCRIPTORFORMAT = internal.KIND_DESCRIPTORFORMAT
)

const (
	DEP_DEPENDENCY  = internal.DEP_DEPENDENCY
	DEP_DESCRIPTION = internal.DEP_DESCRIPTION
	DEP_MEET        = internal.DEP_MEET
	DEP_INSTALLER   = internal.DEP_INSTALLER
)

const (
	REL_TYPE = internal.REL_TYPE
	ABS_TYPE = internal.ABS_TYPE
)

type (
	ServiceKindSpec        = internal.ServiceKindSpec
	ServiceDescriptor      = internal.ServiceDescriptor
	ServiceModelDescriptor = internal.ServiceModelDescriptor

	DescriptionContext = internal.DescriptionContext

	CrossReferences = internal.CrossReferences
	Reference       = internal.Reference
	References      = internal.References

	Origin = internal.Origin
)

func NewDescriptionContext(name, vers string, desc *ServiceModelDescriptor) DescriptionContext {
	return internal.NewDescriptionContext(name, vers, desc)
}

func Decode(data []byte, unmarshaller ...runtime.Unmarshaler) (*ServiceModelDescriptor, error) {
	return internal.Decode(data, unmarshaller...)
}

func Encode(desc *ServiceModelDescriptor, marshaller ...runtime.Marshaler) ([]byte, error) {
	return internal.Encode(desc, marshaller...)
}

func CrossReferencesFor(desc *ServiceModelDescriptor, os ...Origin) *CrossReferences {
	return internal.ServiceModelReferences(desc, os...)
}

func NewNewOCMOrigin(comp, vers string, res v1.Identity) Origin {
	return internal.NewOCMOrigin(common.NewNameVersion(comp, vers), res)
}

func NewNewOCMOriginFor(nv common.VersionedElement, res v1.Identity) Origin {
	return internal.NewOCMOrigin(nv, res)
}
