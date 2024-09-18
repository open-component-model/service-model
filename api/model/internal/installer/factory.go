package installer

import (
	"fmt"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/model/internal"
	"github.com/open-component-model/service-model/api/model/internal/common"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/installer"
)

func init() {
	internal.DefaultServiceKindRegistry.Register(installer.TYPE, Factory)
}

type ServiceVersionVariant struct {
	*common.ServiceVersionVariant
	spec *installer.ServiceSpec
}

func Factory(model internal.Model, descriptor *modeldesc.ServiceDescriptor) (internal.ServiceVersionVariant, error) {
	s, ok := descriptor.Kind.(*installer.ServiceSpec)
	if !ok {
		return nil, fmt.Errorf("invalid service spec type: %T", descriptor.Kind)
	}
	return &ServiceVersionVariant{
		ServiceVersionVariant: common.New(descriptor),
		spec:                  s,
	}, nil
}

func (s *ServiceVersionVariant) AsIntallationService() internal.InstallationService {
	return s
}

func (s *ServiceVersionVariant) GetTargetEnviroment() metav1.TargetEnvironment {
	return s.spec.TargetEnvironment
}

func (s *ServiceVersionVariant) GetInstallerResource() metav1.ResourceReference {
	return s.spec.InstallerResource
}

func (s *ServiceVersionVariant) GetInstallerType() string {
	return s.spec.InstallerType
}
