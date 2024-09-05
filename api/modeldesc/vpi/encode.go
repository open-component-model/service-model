package vpi

import (
	"encoding/json"

	"github.com/mandelsoft/goutils/errors"
	"ocm.software/ocm/api/utils/runtime"
)

type ServiceKindSchemeProvider interface {
	ServiceKindTypeScheme() ServiceKindTypeScheme
}

type ServiceDescriptor[BASE any, S ServiceKindSchemeProvider] struct {
	Base BASE
	Kind ServiceKindSpec
}

var (
	_ json.Marshaler   = (*ServiceDescriptor[any, ServiceKindSchemeProvider])(nil)
	_ json.Unmarshaler = (*ServiceDescriptor[any, ServiceKindSchemeProvider])(nil)
)

func (s ServiceDescriptor[B, S]) MarshalJSON() ([]byte, error) {
	base, err := json.Marshal(&s.Base)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot marshal service")
	}
	kind, err := json.Marshal(&s.Kind)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}

	err = json.Unmarshal(base, &m)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(kind, &m)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (s *ServiceDescriptor[B, S]) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &s.Base)
	if err != nil {
		return errors.Wrapf(err, "cannot unmarshal service descriptor")
	}

	var scheme S
	s.Kind, err = scheme.ServiceKindTypeScheme().Decode(bytes, runtime.DefaultJSONEncoding)
	return err
}
