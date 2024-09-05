package ordinary

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
)

const TYPE = "OrdinaryService"

type ServiceSpec struct {
	metav1.CommonServiceImplementationSpec
}
