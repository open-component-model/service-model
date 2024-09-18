package desc2model

import (
	"github.com/open-component-model/service-model/api/model/internal"
	"github.com/open-component-model/service-model/api/modeldesc"
)

func ServiceDescriptorToModelService(model internal.Model, desc *modeldesc.ServiceDescriptor) (internal.ServiceVersionVariant, error) {
	return internal.DefaultServiceKindRegistry.Create(model, desc)
}
