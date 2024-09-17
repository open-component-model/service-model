package internal

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/set"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/identity"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/runtime"
)

const KIND_SERVICEVERSION = "service version"
const KIND_SERVICEIDENTITY = "service identity"
const KIND_SERVICE_TYPE = "service type"
const KIND_MODELVERSION = "service model version"
const KIND_DESCRIPTORFORMAT = "descriptor format"

const REL_TYPE = "relativeServiceModelDescription"
const ABS_TYPE = "serviceModelDescription"

type CommonServiceSpec = metav1.CommonServiceSpec

type ServiceKindSpec interface {
	runtime.TypedObject
	GetVariant() identity.Variant
	GetDependencies() []metav1.Dependency

	Copy() ServiceKindSpec
	ToCanonicalForm(c DescriptionContext) ServiceKindSpec
	Validate(c DescriptionContext) error
	GetReferences() crossref.References
}

type ServiceDescriptor struct {
	CommonServiceSpec
	Kind   ServiceKindSpec
	Origin identity.Origin
}

func (d *ServiceDescriptor) GetId() identity.ServiceVersionVariantIdentity {
	return identity.NewServiceVersionVariantIdentityFor(d.CommonServiceSpec.GetId(), d.Kind.GetVariant())
}

func (d *ServiceDescriptor) GetVariant() identity.Variant {
	return d.Kind.GetVariant()
}

func (d *ServiceDescriptor) Copy() *ServiceDescriptor {
	return &ServiceDescriptor{
		CommonServiceSpec: *d.CommonServiceSpec.Copy(),
		Kind:              d.Kind.Copy(),
		Origin:            d.Origin.Copy(),
	}
}

type ServiceModelDescriptor struct {
	DocType  runtime.VersionedObjectType `json:"type"`
	Services []ServiceDescriptor
}

func (d *ServiceModelDescriptor) GetType() string {
	return d.DocType.GetType()
}

func (d *ServiceModelDescriptor) GetKind() string {
	return d.DocType.GetKind()
}

func (d *ServiceModelDescriptor) GetVersion() string {
	return d.DocType.GetVersion()
}

func (d *ServiceModelDescriptor) ToCanonicalForm(c DescriptionContext) *ServiceModelDescriptor {
	if runtime.GetKind(d) == ABS_TYPE {
		return d
	}
	r := &ServiceModelDescriptor{
		DocType:  runtime.NewVersionedObjectType(ABS_TYPE, runtime.GetVersion(d)),
		Services: sliceutils.InitialSliceFor(d.Services),
	}
	for i, e := range d.Services {
		r.Services[i] = *ServiceToCanonicalForm(&e, c)
	}
	return r
}

func (d *ServiceModelDescriptor) Validate(ve common.VersionedElement, rv ...ResourceValidator) error {
	c := NewDescriptionContext(ve.GetName(), ve.GetVersion(), d).WithResourceValidator(general.Optional(rv...))
	list := errors.ErrListf("validation errors for component %s version %s", c.GetName(), c.GetVersion())
	if runtime.GetKind(d) == ABS_TYPE {
		// return fmt.Errorf("cannot validate absolute descriptor")
	}
	found := set.Set[string]{}
	for i, e := range d.Services {
		if c.MatchComponent(e.Service) {
			if e.Service.Name() != "" {
				if found.Contains(e.Service.Name()) {
					list.Add(fmt.Errorf("duplicate service definition %d(%s)", i, e.Service.Name()))
				} else {
					found.Add(e.Service.Name())
				}
			}
		}
		list.Addf(nil, ValidateService(&e, c), "service %d(%s)", i, e.Service.Name())
	}
	return list.Result()
}
