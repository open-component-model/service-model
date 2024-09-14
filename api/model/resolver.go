package model

import (
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

type Resolver interface {
	LookupServiceVersionVariant(model Model, id v1.ServiceVersionVariantIdentity) (Service, error)
	ListVersions(id v1.ServiceIdentity, variant ...v1.Variant) ([]string, error)
}
