package installer

import (
	"slices"

	"github.com/open-component-model/service-model/api/identity"
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
	InstalledService  identity.ServiceIdentity `json:"installedService,omitempty"`
	Versions          []string                 `json:"versions,omitempty"`
	InstallerResource metav1.ResourceReference `json:"installerResource"`
	InstallerType     string                   `json:"installerType"`
}

type Converter struct{}

func (c Converter) ConvertFrom(object modeldesc.ServiceKindSpec) (vpi.ServiceKindSpec, error) {
	in := object.(*me.ServiceSpec)
	return &ServiceSpec{
		CommonServiceImplementationSpec: *in.Copy(),
		TargetEnvironment:               in.TargetEnvironment.Copy(),
		InstalledService:                in.InstalledService,
		Versions:                        slices.Clone(in.Versions),
		InstallerResource:               in.InstallerResource,
		InstallerType:                   in.InstallerType,
	}, nil
}

func (c Converter) ConvertTo(object vpi.ServiceKindSpec) (modeldesc.ServiceKindSpec, error) {
	in := object.(*ServiceSpec)
	return &me.ServiceSpec{
		CommonServiceImplementationSpec: *in.Copy(),
		TargetEnvironment:               in.TargetEnvironment.Copy(),
		InstalledService:                in.InstalledService,
		Versions:                        slices.Clone(in.Versions),
		InstallerResource:               in.InstallerResource,
		InstallerType:                   in.InstallerType,
	}, nil
}
