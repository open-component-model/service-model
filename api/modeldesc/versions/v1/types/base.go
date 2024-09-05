package types

import (
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

func ConvertBaseFrom(in *internal.BaseServiceSpec) *BaseServiceSpec {
	return in.Copy()
}

func ConvertBaseTo(in *BaseServiceSpec) *internal.BaseServiceSpec {
	return in.Copy()
}
