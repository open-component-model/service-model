package installer

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
)

const TYPE = "InstallationService"

type ServiceSpec struct {
	metav1.CommonServiceImplementationSpec

	TargetEnvironment metav1.TargetEnvironment `json:"targetEnvironment,omitempty"`
	InstalledService  metav1.ServiceIdentity   `json:"installedService,omitempty"`
	Versions          []string                 `json:"versions,omitempty"`
	InstallerResource metav1.InstallerResource `json:"installerResource"`
	InstallerType     string                   `json:"installerType"`
}
