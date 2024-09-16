package modeldesc

import (
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/runtime"
)

const (
	KIND_SERVICEVERSION   = internal.KIND_SERVICEVERSION
	KIND_SERVICEIDENTITY  = internal.KIND_SERVICEIDENTITY
	KIND_SERVICE_TYPE     = internal.KIND_SERVICE_TYPE
	KIND_MODELVERSION     = internal.KIND_MODELVERSION
	KIND_DESCRIPTORFORMAT = internal.KIND_DESCRIPTORFORMAT
)

const (
	DEP_DEPENDENCY  = crossref.DEP_DEPENDENCY
	DEP_DESCRIPTION = crossref.DEP_DESCRIPTION
	DEP_MEET        = crossref.DEP_MEET
	DEP_INSTALLER   = crossref.DEP_INSTALLER
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
	ResourceValidator  = internal.ResourceValidator

	CrossReferences = crossref.CrossReferences
	Reference       = crossref.Reference
	References      = crossref.References

	Origin = identity.Origin
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

func AddServiceModelReferences(refs *CrossReferences, services []ServiceDescriptor, os ...Origin) {
	internal.AddServiceModelReferences(refs, services, os...)
}

func NewNewOCMOrigin(comp, vers string, res ocmmeta.Identity) Origin {
	return identity.NewOCMOrigin(common.NewNameVersion(comp, vers), res)
}

func NewNewOCMOriginFor(nv common.VersionedElement, res ocmmeta.Identity) Origin {
	return identity.NewOCMOrigin(nv, res)
}

type Resolver interface {
	LookupServiceVersionVariant(id identity.ServiceVersionVariantIdentity) (*ServiceDescriptor, error)
}

type VersionResolver interface {
	Resolver
	ListVersions(id identity.ServiceIdentity, variant ...identity.Variant) ([]string, error)
}
