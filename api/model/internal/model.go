package internal

import (
	"github.com/open-component-model/service-model/api/identity"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"ocm.software/ocm/api/utils/runtime"
)

type Model interface {
	GetServiceKindRegistry() ServiceKindRegistry
	GetServiceVersionVariant(id identity.ServiceVersionVariantIdentity) (ServiceVersionVariant, error)
}

type Service interface {
}

type ServiceVersion interface{}

type ServiceVersionVariant interface {
	GetType() string
	GetName() string
	GetIdentity() identity.ServiceVersionVariantIdentity
	GetComponent() string
	GetVersion() string
	GetVariant() identity.Variant

	AsServiceContract() ServiceContract
	AsInstallationService() InstallationService
	AsOrdinaryService() OrdinaryService
	AsServiceProvider() ServiceProvider
}

type ServiceContract interface {
	ServiceVersionVariant
	GetAPISpecificationType() string
	GetAPISpecVersion() string
	GetSpecification() *runtime.RawValue
	GetArtifact() *metav1.ResourceReference
}

type InstallationService interface{}

type OrdinaryService interface{}

type ServiceProvider interface{}
