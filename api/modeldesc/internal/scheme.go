package internal

import (
	"github.com/mandelsoft/goutils/errors"
	"ocm.software/ocm/api/utils/runtime"
)

type Version interface {
	runtime.TypedObjectType[*ServiceModelDescriptor]
	Encode(*ServiceModelDescriptor, runtime.Marshaler) ([]byte, error)
}

type VersionScheme = runtime.TypeScheme[*ServiceModelDescriptor, Version]

var versions = runtime.NewTypeScheme[*ServiceModelDescriptor, Version]()

func RegisterVersion(v Version) {
	versions.Register(v)
}

func Decode(data []byte) (*ServiceModelDescriptor, error) {
	return versions.Decode(data, runtime.DefaultYAMLEncoding)
}

func Encode(m *ServiceModelDescriptor) ([]byte, error) {
	t := versions.GetType(m.GetType())
	if t == nil {
		return nil, errors.ErrNotFound(KIND_MODELVERSION, m.GetType())
	}
	return t.Encode(m, runtime.DefaultYAMLEncoding.Marshaler)
}
