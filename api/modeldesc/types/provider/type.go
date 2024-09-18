package provider

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/crossref"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

const TYPE = "ServiceProvider"

type ServiceSpec struct {
	metav1.CommonConsumerServiceImplementationSpec

	ManagedServices metav1.ManagedServices `json:"managedServices"`
}

func (s *ServiceSpec) Copy() internal.ServiceKindSpec {
	return &ServiceSpec{
		CommonConsumerServiceImplementationSpec: *s.CommonConsumerServiceImplementationSpec.Copy(),
		ManagedServices:                         s.ManagedServices.Copy(),
	}
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	r := &ServiceSpec{
		CommonConsumerServiceImplementationSpec: *internal.CommonConsumerServiceImplementationSpecToCanonicalForm(&s.CommonConsumerServiceImplementationSpec, c),
		ManagedServices:                         sliceutils.InitialSliceFor(s.ManagedServices),
	}

	for i, e := range s.ManagedServices {
		r.ManagedServices[i] = *internal.ManagedServiceToCanonicalForm(&e, c)
	}
	return r
}

func (s *ServiceSpec) Validate(c internal.DescriptionContext) error {
	var list errors.ErrorList

	list.Add(
		internal.ValidateCommonConsumerServiceImplementation(&s.CommonConsumerServiceImplementationSpec, c),
	)
	for i, e := range s.ManagedServices {
		list.Addf(nil, internal.ValidateManagedService(&e, c), "managed service %d(%s)", i, e.Service.Name)
	}
	return list.Result()
}

func (s *ServiceSpec) GetReferences() crossref.References {
	var refs crossref.References

	refs.Add(internal.CommonConsumerServiceImplementationReferences(&s.CommonConsumerServiceImplementationSpec)...)
	for _, e := range s.ManagedServices {
		refs.Add(internal.ManagedServiceReferences(&e)...)
	}
	return refs
}
