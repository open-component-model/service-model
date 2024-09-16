package internal

import (
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/identity"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
)

func addReferences(c *crossref.CrossReferences, holder identity.ServiceVersionIdentity, variant identity.Variant, refs crossref.References) {
	for _, r := range refs {
		c.AddRef(holder, variant, r.Id, r.Kind)
	}
}

func ServiceModelReferences(d *ServiceModelDescriptor, os ...identity.Origin) *crossref.CrossReferences {
	refs := crossref.NewCrossReferences()
	AddServiceModelReferences(refs, d.Services, os...)
	return refs
}

func AddServiceModelReferences(refs *crossref.CrossReferences, services []ServiceDescriptor, os ...identity.Origin) {
	for _, e := range services {
		desc := e.Copy()
		refs.AddService(crossref.NewServiceVersionIdentity(e.Service, e.Version), e.Kind.GetVariant(), desc, os...)
		refs.AddRefs(ServiceReferences(&e))
	}
}

func ServiceReferences(e *ServiceDescriptor) *crossref.CrossReferences {
	refs := crossref.NewCrossReferences()
	holder := identity.NewServiceVersionId(e.Service, e.Version)

	addReferences(refs, holder, e.Kind.GetVariant(), CommonReferences(&e.CommonServiceSpec))
	addReferences(refs, holder, e.Kind.GetVariant(), e.Kind.GetReferences())
	return refs
}

func CommonReferences(s *CommonServiceSpec) crossref.References {
	return nil
}

func CommonServiceImplementationReferences(s *metav1.CommonServiceImplementationSpec) crossref.References {
	var refs crossref.References

	for _, e := range s.Dependencies {
		refs.Add(DependencyReferences(&e)...)
	}
	for _, e := range s.Contracts {
		refs.Add(ContractReferences(&e)...)
	}
	return refs
}

func CommonConsumerServiceImplementationReferences(s *metav1.CommonConsumerServiceImplementationSpec) crossref.References {
	var refs crossref.References
	refs.Add(CommonServiceImplementationReferences(&s.CommonServiceImplementationSpec)...)
	for _, e := range s.Installers {
		refs.Add(InstallerReferences(&e)...)
	}
	return refs
}

func DependencyReferences(s *metav1.Dependency) crossref.References {
	var refs crossref.References

	crossref.AddVersionReferences(&refs, s.Service, s.Variant, crossref.DEP_DEPENDENCY, s.VersionConstraints...)
	for _, e := range s.ServiceInstances {
		refs.Add(ServiceInstanceReferences(&e)...)
	}
	return refs
}

func ContractReferences(s *metav1.Contract) crossref.References {
	var refs crossref.References
	refs.Add(*crossref.NewReference(s.Service, s.Version, nil, crossref.DEP_MEET))
	return refs
}

func ServiceInstanceReferences(s *metav1.ServiceInstance) crossref.References {
	var refs crossref.References
	crossref.AddVersionReferences(&refs, s.Service, s.Variant, crossref.DEP_DESCRIPTION, s.Versions...)
	return refs
}

func InstallerReferences(s *metav1.Installer) crossref.References {
	var refs crossref.References

	refs.Add(*crossref.NewReference(s.Service, s.Version, nil, crossref.DEP_INSTALLER))
	return refs
}

func ManagedServiceReferences(s *metav1.ManagedService) crossref.References {
	var refs crossref.References
	crossref.AddVersionReferences(&refs, s.Service, s.Variant, crossref.DEP_DESCRIPTION, s.Versions...)
	return refs
}

func DependencyResolutionReferences(s *metav1.DependencyResolution) crossref.References {
	var refs crossref.References
	return refs
}
