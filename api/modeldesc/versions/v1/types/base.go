package types

import (
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

func ConvertBaseFrom(in *internal.CommonServiceSpec) *CommonServiceSpec {
	return in.Copy()
}

func ConvertBaseTo(in *CommonServiceSpec) *internal.CommonServiceSpec {
	return in.Copy()
}
