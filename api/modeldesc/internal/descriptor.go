package internal

import (
	v1 "github.com/open-component-model/service-model/api/meta/v1"
	"ocm.software/ocm/api/utils/runtime"
)

const KIND_SERVICE_TYPE = "service type"
const KIND_MODELVERSION = "service model version"

const REL_TYPE = "sap.com/relativeServiceModelDescriptor"
const ABS_TYPE = "sap.com/serviceModelDescriptor"

type BaseServiceSpec = v1.BaseServiceSpec

type ServiceKindSpec interface {
	runtime.TypedObject
}

type ServiceDescriptor struct {
	BaseServiceSpec
	Kind ServiceKindSpec
}

type ServiceModelDescriptor struct {
	DocType  string `json:"type"`
	Services []ServiceDescriptor
}

func (d *ServiceModelDescriptor) GetType() string {
	return d.DocType
}
