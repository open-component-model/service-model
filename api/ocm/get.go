package ocm

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/ocm"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	"ocm.software/ocm/api/ocm/ocmutils"
	"ocm.software/ocm/api/ocm/resourcerefs"
	"ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/runtime"
)

const RESOURCE_TYPE = modeldesc.ABS_TYPE

func GetServiceModel(comp, vers string, resolver ocm.ComponentVersionResolver) (*modeldesc.ServiceModelDescriptor, *modeldesc.CrossReferences, error) {

	cv, err := resolver.LookupComponentVersion(comp, vers)
	if err != nil {
		return nil, nil, err
	}
	if cv == nil {
		return nil, nil, errors.ErrNotFound(ocm.KIND_COMPONENTVERSION, misc.NewNameVersion(comp, vers).String())
	}
	defer cv.Close()
	return GetServiceModelFromCV(cv, resolver)
}

func GetServiceModelFromCV(cv ocm.ComponentVersionAccess, resolver ...ocm.ComponentVersionResolver) (*modeldesc.ServiceModelDescriptor, *modeldesc.CrossReferences, error) {
	complete := &modeldesc.ServiceModelDescriptor{
		DocType: runtime.NewVersionedObjectType(modeldesc.REL_TYPE, "v1"),
	}
	for i, r := range cv.GetDescriptor().Resources {
		if r.Type == RESOURCE_TYPE {
			origin := modeldesc.NewNewOCMOrigin(cv.GetName(), cv.GetVersion(), r.GetIdentity(cv.GetDescriptor().Resources))
			res, err := cv.GetResourceByIndex(i)
			if err != nil {
				return nil, nil, err
			}
			data, err := ocmutils.GetResourceData(res)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "resource %d[%s]", i, r.Name)
			}
			desc, err := modeldesc.Decode(data)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "resource %d[%s]", i, r.Name)
			}
			if desc.DocType.GetKind() != modeldesc.REL_TYPE {
				return nil, nil, errors.ErrInvalid(modeldesc.KIND_DESCRIPTORFORMAT, desc.DocType.GetKind())
			}
			for _, s := range desc.Services {
				s.Origin = origin
				complete.Services = append(complete.Services, s)
			}
		}
	}
	complete = complete.ToCanonicalForm(modeldesc.NewDescriptionContext(cv.GetName(), cv.GetVersion(), complete))
	err := complete.Validate(cv, ResourceValidator(cv, resolver...))
	if err != nil {
		return complete, nil, err
	}
	refs, err := modeldesc.CheckLocalConsistency(complete)
	return complete, refs, err
}

func ResourceValidator(cv ocm.ComponentVersionAccess, resolver ...ocm.ComponentVersionResolver) modeldesc.ResourceValidator {
	var cvr ocm.ComponentVersionResolver
	if general.Optional(resolver...) != nil {
		cvr = general.Optional(resolver...)
	}
	return func(ref *ocmmeta.ResourceReference) error {
		if ref == nil {
			return nil
		}
		_, rcv, err := resourcerefs.ResolveResourceReference(cv, *ref, cvr)
		if rcv != nil {
			rcv.Close()
		}
		return err
	}
}
