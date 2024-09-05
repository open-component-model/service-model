package ordinary

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	me "github.com/open-component-model/service-model/api/modeldesc/types/ordinary"
	"github.com/open-component-model/service-model/api/modeldesc/versions/v1/types"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
)

const TYPE = "ServiceProvider"

func init() {
	types.RegisterServiceType(vpi.NewServiceKindType[ServiceSpec](me.TYPE, Converter{}))
}

type ServiceSpec struct {
	metav1.CommonServiceImplementationSpec
}

type Converter struct{}

func (c Converter) ConvertFrom(object modeldesc.ServiceKindSpec) (vpi.ServiceKindSpec, error) {
	in := object.(*me.ServiceSpec)
	return &ServiceSpec{*in.Copy()}, nil
}

func (c Converter) ConvertTo(object vpi.ServiceKindSpec) (modeldesc.ServiceKindSpec, error) {
	in := object.(*ServiceSpec)
	return &me.ServiceSpec{*in.Copy()}, nil
}
