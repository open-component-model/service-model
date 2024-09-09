package provider

import (
	"github.com/mandelsoft/goutils/errors"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc/internal"
	"github.com/open-component-model/service-model/api/utils"
)

const TYPE = "ServiceProvider"

type ServiceSpec struct {
	metav1.CommonConsumerServiceImplementationSpec

	ManagedServices metav1.ManagedServices `json:"managedServices"`
}

func (s *ServiceSpec) ToCanonicalForm(c internal.DescriptionContext) internal.ServiceKindSpec {
	r := &ServiceSpec{
		CommonConsumerServiceImplementationSpec: *internal.CommonConsumerServiceImplementationSpecToCanonicalForm(&s.CommonConsumerServiceImplementationSpec, c),
		ManagedServices:                         utils.InitialSliceFor(s.ManagedServices),
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

func (s *ServiceSpec) GetReferences() internal.References {
	var refs internal.References

	refs.Add(internal.CommonConsumerServiceImplementationReferences(&s.CommonConsumerServiceImplementationSpec)...)
	for _, e := range s.ManagedServices {
		refs.Add(internal.ManagedServiceReferences(&e)...)
	}
	refs.Add(internal.CommonServiceImplementationReferences(&s.CommonServiceImplementationSpec)...)
	return refs
}
