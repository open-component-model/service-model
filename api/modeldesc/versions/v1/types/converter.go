package types

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
)

type CommonServiceSpec = metav1.CommonServiceSpec

type ServiceDescriptor = vpi.ServiceDescriptor[CommonServiceSpec, SchemeProvider]

type ServiceDescriptionConverter struct{}

var _ vpi.Converter[*modeldesc.ServiceDescriptor, *ServiceDescriptor] = (*ServiceDescriptionConverter)(nil)

var ServiceListConverter = vpi.NewListConverter[modeldesc.ServiceDescriptor, ServiceDescriptor](ServiceDescriptionConverter{})

func (s ServiceDescriptionConverter) ConvertFrom(in *modeldesc.ServiceDescriptor) (*ServiceDescriptor, error) {
	var err error
	var out ServiceDescriptor

	out.Base = *ConvertBaseFrom(&in.CommonServiceSpec)
	out.Kind, err = DefaultServiceTypeScheme.ConvertFrom(in.Kind)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (s ServiceDescriptionConverter) ConvertTo(in *ServiceDescriptor) (*modeldesc.ServiceDescriptor, error) {
	var err error
	var out modeldesc.ServiceDescriptor

	out.CommonServiceSpec = *ConvertBaseTo(&in.Base)
	out.Kind, err = DefaultServiceTypeScheme.ConvertTo(in.Kind)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
