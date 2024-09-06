package internal

import (
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
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

func Decode(data []byte, unmarshaller ...runtime.Unmarshaler) (*ServiceModelDescriptor, error) {
	return versions.Decode(data, general.OptionalDefaulted[runtime.Unmarshaler](runtime.DefaultYAMLEncoding, unmarshaller...))
}

func Encode(m *ServiceModelDescriptor, marshaller ...runtime.Marshaler) ([]byte, error) {
	t := versions.GetType(m.GetType())
	if t == nil {
		return nil, errors.ErrNotFound(KIND_MODELVERSION, m.GetType())
	}
	return t.Encode(m, general.OptionalDefaulted[runtime.Marshaler](runtime.DefaultYAMLEncoding, marshaller...))
}
