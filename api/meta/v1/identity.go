package v1

import (
	"encoding/json"
	"strings"
)

type ServiceIdentity struct {
	Component string
	Service   string
}

func (id ServiceIdentity) IsRelative() bool {
	return id.Component == ""
}

func (id ServiceIdentity) String() string {
	if id.Component == "" {
		return id.Service
	}
	return id.Service + "@" + id.Component
}

func (id *ServiceIdentity) Parse(s string) {
	idx := strings.Index(s, "@")
	if idx >= 0 {
		id.Service = s[:idx]
		id.Component = s[idx+1:]
	} else {
		id.Service = s
		id.Component = ""
	}
}

func (id ServiceIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ServiceIdentity) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	id.Parse(s)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type ServiceVersionIdentity struct {
	ServiceIdentity
	Version string
}

func (id ServiceVersionIdentity) IsRelative() bool {
	return id.Version == "" && id.ServiceIdentity.IsRelative()
}

func (id ServiceVersionIdentity) String() string {
	if id.Version == "" {
		return id.ServiceIdentity.String()
	}
	return id.ServiceIdentity.String() + ":" + id.Version
}

func (id *ServiceVersionIdentity) Parse(s string) {
	idx := strings.Index(s, ":")
	if idx >= 0 {
		id.ServiceIdentity.Parse(s[:idx])
		id.Version = s[idx+1:]
	} else {
		id.ServiceIdentity.Parse(s[:idx])
		id.Version = ""
	}
}

func (id ServiceVersionIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ServiceVersionIdentity) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	id.Parse(s)
	return nil
}
