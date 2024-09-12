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
	InstalledService  metav1.ServiceIdentity   `json:"installedService,omitempty"`
	Versions          []string                 `json:"versions,omitempty"`
	InstallerResource metav1.ResourceReference `json:"installerResource"`
	InstallerType     string                   `json:"installerType"`
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	r := *s
	r.CommonServiceImplementationSpec = *internal.CommonServiceImplementationSpecToCanonicalForm(&r.CommonServiceImplementationSpec, c)

	if s.InstalledService.Name != "" {
		if r.InstalledService.IsRelative() {
			r.InstalledService.Component = c.GetName()
		}
		if len(r.Versions) == 0 {
			r.Versions = []string{c.GetVersion()}
		}
	}
	return &r
}

func (s *ServiceSpec) Validate(c internal.DescriptionContext) error {
	var list errors.ErrorList

	if s.InstalledService.Name != "" {
		if s.InstalledService.Component == c.GetName() || s.InstalledService.Component == "" {
			if c.LookupService(s.InstalledService.Name) == nil {
				list.Add(fmt.Errorf("local installer service %q not defined", s.InstalledService.Name))
			}
		}
	} else {
		if len(s.Versions) > 0 {
			list.Add(fmt.Errorf("versions must not be set for omitted installedService"))
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
	crossref.AddVersionReferences(&refs, s.InstalledService, crossref.DEP_DESCRIPTION, s.Versions...)
	return refs
}
