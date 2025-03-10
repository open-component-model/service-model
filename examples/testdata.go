package examples

import (
	"ocm.software/ocm/api/helper/env"
)

const (
	COMP_MSP_GARDENER      = "acme.org/gardener/service"
	VERS_MSP_GARDENER      = "v1.0.0"
	NAME_MSP_GARDENER_PROV = "provider"
	NAME_MSP_GARDENER_INST = "installer"
)

const (
	COMP_MSP_HANA      = "acme.org/hana/service"
	VERS_MSP_HANA      = "v1.0.0"
	NAME_MSP_HANA_PROV = "provider"
	NAME_MSP_HANA_INST = "installer"
)

const (
	COMP_MSP_STEAMPUNK      = "acme.org/steampunk/service"
	VERS_MSP_STEAMPUNK      = "v1.0.0"
	NAME_MSP_STEAMPUNK_PROV = "provider"
	NAME_MSP_STEAMPUNK_INST = "installer"
)

const (
	COMP_CONTRACT_K8S_CLUSTER_22 = "acme.org/kubernetes/apis"
	VERS_CONTRACT_K8S_CLUSTER_22 = "v1.22.0"
	NAME_CONTRACT_K8S_CLUSTER_22 = "cluster"
)

func Descriptors(dest ...string) env.Option {
	return env.ProjectTestDataForCaller("descriptors", dest...)
}
