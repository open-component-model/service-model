package v1

import (
	"encoding/json"
	"strings"

	"github.com/open-component-model/service-model/api/utils"
)

type ServiceIdentity struct {
	Component string
	Name      string
}

func (id ServiceIdentity) Validate() error {
	return utils.CheckFlatName(id.Name, "service name")
}

func (id ServiceIdentity) IsRelative() bool {
	return id.Component == ""
}

func (id ServiceIdentity) String() string {
	if id.Component == "" {
		return id.Name
	}
	return id.Component + "/" + id.Name
}

func (id *ServiceIdentity) Parse(s string) error {
	idx := strings.LastIndex(s, "/")
	if idx >= 0 {
		id.Component = s[:idx]
		id.Name = s[idx+1:]
	} else {
		id.Name = s
		id.Component = ""
	}
	return nil
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
