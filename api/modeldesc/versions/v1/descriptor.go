package v1

import (
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	"github.com/open-component-model/service-model/api/modeldesc/versions/v1/types"
)

type ServiceModelDescriptor struct {
	DocType  string                    `json:"type"`
	Services []types.ServiceDescriptor `json:"services"`
}

func (in *ServiceModelDescriptor) GetType() string {
	return in.DocType
}

type Converter struct{}

func (c Converter) ConvertFrom(in *modeldesc.ServiceModelDescriptor) (*ServiceModelDescriptor, error) {
	var out ServiceModelDescriptor
	var err error
	out.DocType = in.DocType
	out.Services, err = types.ServiceListConverter.ConvertFrom(in.Services)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c Converter) ConvertTo(in *ServiceModelDescriptor) (*modeldesc.ServiceModelDescriptor, error) {
	var out modeldesc.ServiceModelDescriptor

	var err error
	out.DocType = in.DocType
	out.Services, err = types.ServiceListConverter.ConvertTo(in.Services)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
