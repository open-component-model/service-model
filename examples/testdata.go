package examples

import (
	"ocm.software/ocm/api/helper/env"
)

const (
	COMP_MSP_GARDENER = "acme.org/gardener/service"
	VERS_MSP_GARDENER = "v1.0.0"
)

const (
	COMP_MSP_HANA = "acme.org/hana/service"
	VERS_MSP_HANA = "v1.0.0"
)

const (
	COMP_MSP_STEAMPUNK = "acme.org/steampunk/service"
	VERS_MSP_STEAMPUNK = "v1.0.0"
)

func Descriptors(dest ...string) env.Option {
	return env.ProjectTestDataForCaller("descriptors", dest...)
}
