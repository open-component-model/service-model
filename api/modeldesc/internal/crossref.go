package internal

import (
	"github.com/open-component-model/service-model/api/common"
	"github.com/open-component-model/service-model/api/crossref"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

func addReferences(c *crossref.CrossReferences, holder *crossref.ServiceVersionIdentity, refs crossref.References) {
	for _, r := range refs {
		c.AddRef(holder, &r.Id, r.Kind)
	}
}

func ServiceModelReferences(d *ServiceModelDescriptor, os ...common.Origin) *crossref.CrossReferences {
	refs := crossref.NewCrossReferences()
	AddServiceModelReferences(refs, d.Services, os...)
	return refs
}

func AddServiceModelReferences(refs *crossref.CrossReferences, services []ServiceDescriptor, os ...common.Origin) {
	for _, e := range services {
		refs.AddService(crossref.NewServiceVersionIdentity(e.Service, e.Version), os...)
		refs.AddRefs(ServiceReferences(&e))
	}
}

func ServiceReferences(e *ServiceDescriptor) *crossref.CrossReferences {
	refs := crossref.NewCrossReferences()
	holder := crossref.NewServiceVersionIdentity(e.Service, e.Version)

	addReferences(refs, holder, CommonReferences(&e.CommonServiceSpec))
	addReferences(refs, holder, e.Kind.GetReferences())
	return refs
}

func CommonReferences(s *CommonServiceSpec) crossref.References {
	return nil
}

func CommonServiceImplementationReferences(s *v1.CommonServiceImplementationSpec) crossref.References {
	var refs crossref.References

	for _, e := range s.Dependencies {
		refs.Add(DependencyReferences(&e)...)
	}
	for _, e := range s.Contracts {
		refs.Add(ContractReferences(&e)...)
	}
	return refs
}

func CommonConsumerServiceImplementationReferences(s *v1.CommonConsumerServiceImplementationSpec) crossref.References {
	var refs crossref.References
	refs.Add(CommonServiceImplementationReferences(&s.CommonServiceImplementationSpec)...)
	for _, e := range s.Installers {
		refs.Add(InstallerReferences(&e)...)
	}
	return refs
}

func DependencyReferences(s *v1.Dependency) crossref.References {
	var refs crossref.References

	crossref.AddVersionReferences(&refs, s.Service, crossref.DEP_DEPENDENCY, s.VersionConstraints...)
	for _, e := range s.ServiceInstances {
		refs.Add(ServiceInstanceReferences(&e)...)
	}
	return refs
}

func ContractReferences(s *v1.Contract) crossref.References {
	var refs crossref.References
	refs.Add(*crossref.NewReference(s.Service, s.Version, crossref.DEP_MEET))
	return refs
}

func ServiceInstanceReferences(s *v1.ServiceInstance) crossref.References {
	var refs crossref.References
	crossref.AddVersionReferences(&refs, s.Service, crossref.DEP_DESCRIPTION, s.Versions...)
	return refs
}

func InstallerReferences(s *v1.Installer) crossref.References {
	var refs crossref.References

	refs.Add(*crossref.NewReference(s.Service, s.Version, crossref.DEP_INSTALLER))
	return refs
}

func ManagedServiceReferences(s *v1.ManagedService) crossref.References {
	var refs crossref.References
	crossref.AddVersionReferences(&refs, s.Service, crossref.DEP_DESCRIPTION, s.Versions...)
	return refs
}

func DependencyResolutionReferences(s *v1.DependencyResolution) crossref.References {
	var refs crossref.References
	return refs
}
