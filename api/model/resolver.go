package model

import (
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

type Resolver interface {
	LookupServiceVersion(model Model, id v1.ServiceVersionIdentity) (Service, error)
	ListVersions(id v1.ServiceIdentity) ([]string, error)
}
