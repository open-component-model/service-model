package internal

import (
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

func addReferences(c *CrossReferences, holder *ServiceVersionIdentity, refs References) {
	for _, r := range refs {
		c.AddRef(holder, &r.Id, r.Kind)
	}
}

func ServiceModelReferences(d *ServiceModelDescriptor, os ...Origin) *CrossReferences {
	refs := NewCrossReferences()
	AddServiceModelReferences(refs, d.Services, os...)
	return refs
}

func AddServiceModelReferences(refs *CrossReferences, services []ServiceDescriptor, os ...Origin) {
	for _, e := range services {
		refs.AddService(NewServiceVersionIdentity(e.Service, e.Version), os...)
		refs.AddRefs(ServiceReferences(&e))
	}
}

func ServiceReferences(e *ServiceDescriptor) *CrossReferences {
	refs := NewCrossReferences()
	holder := NewServiceVersionIdentity(e.Service, e.Version)

	addReferences(refs, holder, CommonReferences(&e.CommonServiceSpec))
	addReferences(refs, holder, e.Kind.GetReferences())
	return refs
}

func CommonReferences(s *CommonServiceSpec) References {
	return nil
}

func CommonServiceImplementationReferences(s *v1.CommonServiceImplementationSpec) References {
	var refs References

	for _, e := range s.Dependencies {
		refs.Add(DependencyReferences(&e)...)
	}
	for _, e := range s.Contracts {
		refs.Add(ContractReferences(&e)...)
	}
	return refs
}

func CommonConsumerServiceImplementationReferences(s *v1.CommonConsumerServiceImplementationSpec) References {
	var refs References
	refs.Add(CommonServiceImplementationReferences(&s.CommonServiceImplementationSpec)...)
	for _, e := range s.Installers {
		refs.Add(InstallerReferences(&e)...)
	}
	return refs
}

func DependencyReferences(s *v1.Dependency) References {
	var refs References

	AddVersionReferences(&refs, s.Service, DEP_DEPENDENCY, s.VersionConstraints...)
	for _, e := range s.ServiceInstances {
		refs.Add(ServiceInstanceReferences(&e)...)
	}
	return refs
}

func ContractReferences(s *v1.Contract) References {
	var refs References
	refs.Add(*NewReference(s.Service, s.Version, DEP_MEET))
	return refs
}

func ServiceInstanceReferences(s *v1.ServiceInstance) References {
	var refs References
	AddVersionReferences(&refs, s.Service, DEP_DESCRIPTION, s.Versions...)
	return refs
}

func InstallerReferences(s *v1.Installer) References {
	var refs References

	refs.Add(*NewReference(s.Service, s.Version, DEP_INSTALLER))
	return refs
}

func ManagedServiceReferences(s *v1.ManagedService) References {
	var refs References
	AddVersionReferences(&refs, s.Service, DEP_DESCRIPTION, s.Versions...)
	return refs
}

func DependencyResolutionReferences(s *v1.DependencyResolution) References {
	var refs References
	return refs
}
