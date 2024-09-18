package installer

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	me "github.com/open-component-model/service-model/api/modeldesc/types/installer"
	"github.com/open-component-model/service-model/api/modeldesc/versions/v1/types"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
)

func init() {
	types.RegisterServiceType(vpi.NewServiceKindType[ServiceSpec](me.TYPE, Converter{}))
}

type ServiceSpec struct {
	metav1.CommonServiceImplementationSpec

	TargetEnvironment metav1.TargetEnvironment `json:"targetEnvironment,omitempty"`
	InstalledServices metav1.InstalledServices `json:"installedServices,omitempty"`
	InstallerResource metav1.ResourceReference `json:"installerResource"`
	InstallerType     string                   `json:"installerType"`
}

type Converter struct{}

func (c Converter) ConvertFrom(object modeldesc.ServiceKindSpec) (vpi.ServiceKindSpec, error) {
	in := object.(*me.ServiceSpec)
	return &ServiceSpec{
		CommonServiceImplementationSpec: *in.CommonServiceImplementationSpec.Copy(),
		TargetEnvironment:               in.TargetEnvironment.Copy(),
		InstalledServices:               in.InstalledServices.Copy(),
		InstallerResource:               in.InstallerResource,
		InstallerType:                   in.InstallerType,
	}, nil
}

func (c Converter) ConvertTo(object vpi.ServiceKindSpec) (modeldesc.ServiceKindSpec, error) {
	in := object.(*ServiceSpec)
	return &me.ServiceSpec{
		CommonServiceImplementationSpec: *in.CommonServiceImplementationSpec.Copy(),
		TargetEnvironment:               in.TargetEnvironment.Copy(),
		InstalledServices:               in.InstalledServices.Copy(),
		InstallerResource:               in.InstallerResource,
		InstallerType:                   in.InstallerType,
	}, nil
}
