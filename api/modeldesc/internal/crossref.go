package internal

import (
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

func addReferences(c *CrossReferences, holder *Reference, refs References) {
	for _, r := range refs {
		c.AddRef(holder, &r)
	}
}

func ServiceModelReferences(d *ServiceModelDescriptor, os ...Origin) *CrossReferences {
	refs := NewCrossReferences()
	for _, e := range d.Services {
		refs.AddService(NewReference(e.Service, e.Version), os...)
		refs.AddRefs(ServiceReferences(&e))
	}
	return refs
}

func ServiceReferences(e *ServiceDescriptor) *CrossReferences {
	refs := NewCrossReferences()
	holder := NewReference(e.Service, e.Version)

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

	if len(s.VersionConstraints) == 0 {
		refs.Add(*NewReference(s.Service, ""))
	} else {
		for _, e := range s.VersionConstraints {
			refs.Add(*NewReference(s.Service, e))
		}
	}
	return refs
}

func ContractReferences(s *v1.Contract) References {
	var refs References
	refs.Add(*NewReference(s.Service, s.Version))
	return refs
}

func ServiceInstanceReferences(s *v1.ServiceInstance) References {
	var refs References
	return refs
}

func InstallerReferences(s *v1.Installer) References {
	var refs References

	refs.Add(*NewReference(s.Service, s.Version))
	return refs
}

func ManagedServiceReferences(s *v1.ManagedService) References {
	var refs References
	return refs
}

func DependencyResolutionReferences(s *v1.DependencyResolution) References {
	var refs References
	return refs
}
