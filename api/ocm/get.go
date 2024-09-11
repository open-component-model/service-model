package ocm

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/ocm"
	ocmutils "ocm.software/ocm/api/ocm/ocmutils"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/runtime"
)

const RESOURCE_TYPE = modeldesc.ABS_TYPE

func GetServiceModel(comp, vers string, resolver ocm.ComponentResolver) (*modeldesc.ServiceModelDescriptor, error) {
	for _, r := range resolver.LookupComponentProviders(comp) {
		c, err := r.LookupComponent(comp)
		if err != nil || c == nil {
			continue
		}
		cv, err := c.LookupVersion(vers)
		if err != nil || cv == nil {
			continue
		}
		defer cv.Close()
		return GetServiceModelFromCV(cv)
	}
	return nil, errors.ErrNotFound(ocm.KIND_COMPONENTVERSION, common.NewNameVersion(comp, vers).String())
}

func GetServiceModelFromCV(cv ocm.ComponentVersionAccess) (*modeldesc.ServiceModelDescriptor, error) {
	complete := &modeldesc.ServiceModelDescriptor{
		DocType: runtime.NewVersionedObjectType(modeldesc.REL_TYPE, "v1"),
	}
	for i, r := range cv.GetDescriptor().Resources {
		if r.Type == RESOURCE_TYPE {
			origin := modeldesc.NewNewOCMOrigin(cv.GetName(), cv.GetVersion(), r.GetIdentity(cv.GetDescriptor().Resources))
			res, err := cv.GetResourceByIndex(i)
			if err != nil {
				return nil, err
			}
			data, err := ocmutils.GetResourceData(res)
			if err != nil {
				return nil, err
			}
			desc, err := modeldesc.Decode(data)
			if err != nil {
				return nil, err
			}
			if desc.DocType.GetKind() != modeldesc.REL_TYPE {
				return nil, errors.ErrInvalid(modeldesc.KIND_DESCRIPTORFORMAT, desc.DocType.GetKind())
			}
			for _, s := range desc.Services {
				s.Origin = origin
				complete.Services = append(complete.Services, s)
			}
		}
	}
	complete = complete.ToCanonicalForm(modeldesc.NewDescriptionContext(cv.GetName(), cv.GetVersion(), complete))
	err := complete.Validate(cv)
	return complete, err
}
