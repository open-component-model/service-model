package vpi

import (
	"github.com/mandelsoft/goutils/general"
	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	"ocm.software/ocm/api/utils/runtime"
)

type Version[V runtime.TypedObject] struct {
	decode    runtime.TypedObjectType[V]
	converter Converter[*modeldesc.ServiceModelDescriptor, V]
}

func NewVersion[I any, P versionPointer[I]](name string, converter Converter[*modeldesc.ServiceModelDescriptor, P]) modeldesc.Version {
	var proto I

	return &Version[P]{
		decode:    runtime.NewTypedObjectTypeByProto[P](name, P(&proto)),
		converter: converter,
	}
}

func (v *Version[V]) GetType() string {
	return v.decode.GetType()
}

func (v *Version[V]) Decode(data []byte, unmarshaler runtime.Unmarshaler) (*modeldesc.ServiceModelDescriptor, error) {
	spec, err := v.decode.Decode(data, unmarshaler)
	if err != nil {
		return nil, err
	}
	return v.converter.ConvertTo(spec)
}

func (v *Version[V]) Encode(descriptor *modeldesc.ServiceModelDescriptor, marshaller runtime.Marshaler) ([]byte, error) {
	spec, err := v.converter.ConvertFrom(descriptor)
	if err != nil {
		return nil, err
	}
	marshaller = general.OptionalDefaulted(runtime.DefaultYAMLEncoding.Marshaler, marshaller)
	return marshaller.Marshal(spec)
}

var _ modeldesc.Version = (*Version[runtime.TypedObject])(nil)

type versionPointer[I any] interface {
	*I
	runtime.TypedObject
}
