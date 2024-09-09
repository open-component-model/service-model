package internal

import (
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/utils"
)

func CommonToCanonicalForm(in *metav1.CommonServiceSpec, c DescriptionContext) *metav1.CommonServiceSpec {
	out := in.Copy()
	if c.MatchComponent(out.Service) {
		out.Service.Component = c.GetName()
	}
	if out.Version == "" {
		out.Version = c.GetVersion()
	}
	return out
}

func CommonServiceImplementationSpecToCanonicalForm(in *metav1.CommonServiceImplementationSpec, c DescriptionContext) *metav1.CommonServiceImplementationSpec {
	out := in.Copy()
	for i, e := range out.Dependencies {
		out.Dependencies[i] = *DependencyToCanonicalForm(&e, c)
	}
	for i, e := range out.Contracts {
		out.Contracts[i] = *ContractToCanonicalForm(&e, c)
	}
	return out
}

func ContractToCanonicalForm(in *metav1.Contract, c DescriptionContext) *metav1.Contract {
	out := in.Copy()
	if c.MatchComponent(out.Service) {
		out.Service.Component = c.GetName()
	}
	if out.Version == "" {
		out.Version = c.GetVersion()
	}
	return out
}

func DependencyToCanonicalForm(in *metav1.Dependency, c DescriptionContext) *metav1.Dependency {
	out := in.Copy()
	if c.MatchComponent(in.Service) {
		out.Service.Component = c.GetName()
	}
	if len(in.VersionConstraints) == 0 {
		out.VersionConstraints = []string{c.GetVersion()}
	}

	for i, e := range in.ServiceInstances {
		out.ServiceInstances[i] = *ServiceInstanceToCanonicalForm(&e, c)
	}
	return out
}

func ServiceInstanceToCanonicalForm(in *metav1.ServiceInstance, c DescriptionContext) *metav1.ServiceInstance {
	out := in.Copy()
	if c.MatchComponent(in.Service) {
		out.Service.Component = c.GetName()
	}
	// cannot default version list, because of constraints
	// and no fixed version) in dependency.
	return out
}

func CommonConsumerServiceImplementationSpecToCanonicalForm(in *metav1.CommonConsumerServiceImplementationSpec, c DescriptionContext) *metav1.CommonConsumerServiceImplementationSpec {
	out := &metav1.CommonConsumerServiceImplementationSpec{
		CommonServiceImplementationSpec: *CommonServiceImplementationSpecToCanonicalForm(&in.CommonServiceImplementationSpec, c),
	}
	out.Installers = utils.InitialSliceFor(in.Installers)
	for i, e := range in.Installers {
		out.Installers[i] = *InstallerToCanonicalForm(&e, c)
	}
	return out
}

func InstallerToCanonicalForm(in *metav1.Installer, c DescriptionContext) *metav1.Installer {
	out := in.Copy()
	if c.MatchComponent(out.Service) {
		out.Service.Component = c.GetName()
	}
	if out.Version == "" {
		out.Version = c.GetVersion()
	}
	return out
}

func ManagedServiceToCanonicalForm(in *metav1.ManagedService, c DescriptionContext) *metav1.ManagedService {
	out := in.Copy()
	if c.MatchComponent(out.Service) {
		out.Service.Component = c.GetName()
	}
	if len(out.Versions) == 0 {
		out.Versions = utils.Slice(c.GetVersion())
	}
	return out
}

func ServiceToCanonicalForm(in *ServiceDescriptor, c DescriptionContext) *ServiceDescriptor {
	return &ServiceDescriptor{
		CommonServiceSpec: *CommonToCanonicalForm(&in.CommonServiceSpec, c),
		Kind:              in.Kind.ToCanonicalForm(c),
	}
}
