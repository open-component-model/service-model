package model

import (
	"github.com/open-component-model/service-model/api/identity"
)

type Resolver interface {
	LookupServiceVersionVariant(model Model, id identity.ServiceVersionVariantIdentity) (ServiceVersionVariant, error)
	ListVersions(id identity.ServiceIdentity, variant ...identity.Variant) ([]string, error)
}
