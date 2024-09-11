package desc2model

import (
	"github.com/open-component-model/service-model/api/model"
	"github.com/open-component-model/service-model/api/modeldesc"
)

func ServiceDescriptorToModelService(desc *modeldesc.ServiceDescriptor) (model.Service, error) {
	return desc, nil
}