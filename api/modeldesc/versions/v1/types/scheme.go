package types

import (
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
)

var DefaultServiceTypeScheme = vpi.NewServiceKindTypeScheme()

type SchemeProvider struct{}

func (p SchemeProvider) ServiceKindTypeScheme() vpi.ServiceKindTypeScheme {
	return DefaultServiceTypeScheme
}

func RegisterServiceType(rtype vpi.ServiceKindType) {
	DefaultServiceTypeScheme.Register(rtype)
}

func RegisterServiceTypes(s vpi.ServiceKindTypeScheme) {
	DefaultServiceTypeScheme.AddKnownTypes(s)
}
