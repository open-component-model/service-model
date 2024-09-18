package installer

import (
	"fmt"
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/crossref"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

const TYPE = "InstallationService"

type ServiceSpec struct {
	metav1.CommonServiceImplementationSpec

	TargetEnvironment metav1.TargetEnvironment `json:"targetEnvironment,omitempty"`
	InstalledServices metav1.InstalledServices `json:"installedServices,omitempty"`
	InstallerResource metav1.ResourceReference `json:"installerResource"`
	InstallerType     string                   `json:"installerType"`
}

func (s *ServiceSpec) Copy() internal.ServiceKindSpec {
	return &ServiceSpec{
		CommonServiceImplementationSpec: *s.CommonServiceImplementationSpec.Copy(),
		TargetEnvironment:               s.TargetEnvironment.Copy(),
		InstalledServices:               s.InstalledServices.Copy(),
		InstallerResource:               *s.InstallerResource.Copy(),
		InstallerType:                   s.InstallerType,
	}
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	r := *s
	r.CommonServiceImplementationSpec = *internal.CommonServiceImplementationSpecToCanonicalForm(&r.CommonServiceImplementationSpec, c)

	for i, is := range r.InstalledServices {
		if is.Service.Name() != "" {
			if is.Service.IsRelative() {
				is.Service = is.Service.ForComponent(c.GetName())
			}
			if len(is.Versions) == 0 {
				is.Versions = []string{c.GetVersion()}
			}
			r.InstalledServices[i] = is
		}
	}
	return &r
}

func (s *ServiceSpec) Validate(c internal.DescriptionContext) error {
	var list errors.ErrorList

	for i, is := range s.InstalledServices {
		if is.Service.Name() == "" {
			list.Add(fmt.Errorf("InstalledService #%d must have a name", i))
		}
		if is.Service.Component() == c.GetName() || is.Service.Component() == "" {
			if c.LookupService(is.Service.Name()) == nil {
				list.Add(fmt.Errorf("local installer service %q not defined", is.Service.Name()))
			}
		}
	}
	if s.InstallerType == "" {
		list.Add(fmt.Errorf("installerType must be set"))
	}
	list.Add(errors.Wrapf(c.ValidateResource(s.InstallerResource.AsResourceRef()), "installer resource %s", s.InstallerResource.AsResourceRef()))
	return list.Result()
}

func (s *ServiceSpec) GetReferences() crossref.References {
	var refs crossref.References

	refs.Add(internal.CommonServiceImplementationReferences(&s.CommonServiceImplementationSpec)...)
	for _, is := range s.InstalledServices {
		crossref.AddVersionReferences(&refs, "", is.Service, is.Variant, crossref.DEP_INSTALLS, is.Versions...)
	}
	return refs
}
