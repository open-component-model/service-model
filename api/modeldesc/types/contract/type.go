package installer

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"ocm.software/ocm/api/utils/runtime"
)

const TYPE = "ServiceContract"

type ServiceSpec struct {
	runtime.ObjectTypedObject `json:",inline"`

	APISpecificationType string                    `json:"apiSpecificationType,omitempty"`
	APISpecVersion       string                    `json:"apiSpecificationVersion,omitempty"`
	Specification        *runtime.RawValue         `json:"specification,omitempty"`
	Artifact             *metav1.InstallerResource `json:"artifact,omitempty"`
}
