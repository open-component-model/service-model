package provider

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
)

const TYPE = "ServiceProvider"

type ServiceSpec struct {
	metav1.CommonConsumerServiceImplementationSpec

	ManagedServices metav1.ManagedServices `json:"managedServices"`
}
