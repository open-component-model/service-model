package internal

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/utils"
)

func ValidateService(e *ServiceDescriptor, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(ValidateCommon(&e.CommonServiceSpec, c))
	list.Add(e.Kind.Validate(c))
	return list.Result()
}

func ValidateCommon(s *CommonServiceSpec, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		s.Service.Validate(),
		utils.CheckNonEmpty(s.ShortName, "short name"),
		s.Labels.Validate(),
	)
	if s.Version != "" {
		list.Add(utils.CheckVersion(s.Version, "service version"))
	}
	if !c.MatchComponent(s.Service) {
		list.Add(fmt.Errorf("non-local service definitions are not possible (%s)", s.Service.Component))
	}
	return list.Result()
}

func ValidateCommonServiceImplementation(s *v1.CommonServiceImplementationSpec, c DescriptionContext) error {
	var list errors.ErrorList

	for i, e := range s.Dependencies {
		list.Add(errors.Wrapf(ValidateDependency(&e, c), "dependency %d(%s)", i, e.Name))
	}
	for i, e := range s.Contracts {
		list.Add(errors.Wrapf(ValidateContract(&e, c), "contract %d", i))
	}
	return list.Result()
}

func ValidateCommonConsumerServiceImplementation(s *v1.CommonConsumerServiceImplementationSpec, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(ValidateCommonServiceImplementation(&s.CommonServiceImplementationSpec, c))
	for i, e := range s.Installers {
		list.Add(errors.Wrapf(ValidateInstaller(&e, c), "installer %d(%s)", i, e.Service.Name))
	}
	return list.Result()
}

func ValidateDependency(s *v1.Dependency, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		utils.CheckFlatName(s.Name, "dependency name"),
		s.Service.Validate(),
		utils.CheckValues(s.Kind, "dependency kind", v1.DEPKIND_IMPLEMENTATION, v1.DEPKIND_ORCHESTRATION),
		s.Labels.Validate(),
	)
	for i, e := range s.VersionConstraints {
		list.Add(utils.CheckVersionConstraint(e, "version conraint %d", i))
	}
	for i, e := range s.ServiceInstances {
		list.Addf(nil, ValidateServiceInstance(&e, c), "service instance %d", i)
	}
	return list.Result()
}

func ValidateContract(s *v1.Contract, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		s.Service.Validate(),
		s.Labels.Validate(),
	)
	if s.Version != "" {
		list.Add(utils.CheckVersion(s.Version, "contract version"))
	}
	return list.Result()
}

func ValidateServiceInstance(s *v1.ServiceInstance, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		s.Service.Validate(),
	)
	for i, e := range s.Static {
		list.Addf(nil, utils.CheckFlatName(e.Name, "name"), "static instance %d", i)
	}
	for i, e := range s.Versions {
		list.Add(utils.CheckVersion(e, "version %d", i))
	}
	return list.Result()
}

func ValidateInstaller(s *v1.Installer, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		s.Service.Validate(),
		s.Labels.Validate(),
	)

	if s.Version != "" {
		list.Add(utils.CheckVersion(s.Version, "version"))
	}
	return list.Result()
}

func ValidateManagedService(s *v1.ManagedService, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		s.Service.Validate(),
		s.Labels.Validate(),
	)
	for i, e := range s.DependencyResolutions {
		list.Addf(nil, ValidateDependencyResolution(&e, c), "resolution %d", i)
	}
	return list.Result()
}

func ValidateDependencyResolution(s *v1.DependencyResolution, c DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		utils.CheckFlatName(s.Name, "name"),
		utils.CheckValues(s.Resolution, "resolution", v1.DEPRES_MANGED, v1.DEPRES_CONFIGURED),
		utils.CheckValues(s.Usage, "usage", v1.DEPUSE_SHARED, v1.DEPUSE_CONFIGURED, v1.DEPUSE_EXCLUSIVE),
		s.Labels.Validate(),
	)
	return list.Result()
}
