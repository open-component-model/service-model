package internal

import (
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/model/internal/common"
)

type Model interface {
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
}

type ServiceContract interface {
	common.ServiceVersionVariant
	GetSpecification() string
}
