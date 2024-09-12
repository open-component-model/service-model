package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/utils"
	common "ocm.software/ocm/api/utils/misc"
)

type ServiceIdentity struct {
	Component string
	Name      string
}

func NewServiceId(component, name string) ServiceIdentity {
	return ServiceIdentity{component, name}
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

func (id ServiceIdentity) MarshalMapKey() (string, error) {
	return id.String(), nil
}

func (id *ServiceIdentity) UnmarshalMapKey(key string) error {
	return id.Parse(key)
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
	ServiceIdentity `json:",inline"`
	Version         string `json:"version"`
}

func NewServiceVersionId(id ServiceIdentity, vers string) ServiceVersionIdentity {
	return ServiceVersionIdentity{id, vers}
}

func (id ServiceVersionIdentity) ComponentVersionId() common.NameVersion {
	return common.NewNameVersion(id.Component, id.Version)
}

func (id ServiceVersionIdentity) IsRelative() bool {
	return id.Version == "" && id.ServiceIdentity.IsRelative()
}

func (id ServiceVersionIdentity) IsConstraint() bool {
	_, err := semver.NewVersion(id.Version)
	return err != nil
}

func (r ServiceVersionIdentity) Equals(o ServiceVersionIdentity) bool {
	return r == o
}

func (id ServiceVersionIdentity) String() string {
	if id.Version == "" {
		return id.ServiceIdentity.String()
	}
	return id.ServiceIdentity.String() + ":" + id.Version
}

func (id *ServiceVersionIdentity) Parse(s string) error {
	idx := strings.LastIndex(s, ":")
	if idx >= 0 {
		id.ServiceIdentity.Parse(s[:idx])
		id.Version = s[idx+1:]
	} else {
		id.ServiceIdentity.Parse(s)
		id.Version = ""
	}
	return nil
}

func (id ServiceVersionIdentity) MarshalMapKey() (string, error) {
	return id.String(), nil
}

func (id *ServiceVersionIdentity) UnmarshalMapKey(key string) error {
	return id.Parse(key)
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

////////////////////////////////////////////////////////////////////////////////

type ServiceVersionIdentities sliceutils.Slice[ServiceVersionIdentity]

func (r *ServiceVersionIdentities) Add(refs ...ServiceVersionIdentity) {
	*r = sliceutils.AppendUnique(*r, refs...)
}

func (r ServiceVersionIdentities) Len() int {
	return len(r)
}

func (r ServiceVersionIdentities) Less(i, j int) bool {
	return ServiceVersionIdentityCompare(r[i], r[j]) < 0
}

func (r ServiceVersionIdentities) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func ServiceVersionIdentityEquals(a, b ServiceVersionIdentity) bool {
	return a.Equals(b)
}

func ServiceVersionIdentityCompare(a, b ServiceVersionIdentity) int {
	c := strings.Compare(a.Component, b.Component)
	if c == 0 {
		c = strings.Compare(a.Name, b.Name)
	}
	if c == 0 {
		c = strings.Compare(a.Version, b.Version)
	}
	return c
}

////////////////////////////////////////////////////////////////////////////////

type ServiceVariantVersionIdentity struct {
	ServiceVersionIdentity
	Variant Variant
}

func NewServiceVariantVersionIdentity(si ServiceIdentity, vers string, variant ...Variant) *ServiceVariantVersionIdentity {
	return &ServiceVariantVersionIdentity{NewServiceVersionId(si, vers), general.Optional(variant...)}
}

func (id ServiceVariantVersionIdentity) String() string {
	if len(id.Variant) == 0 {
		return id.ServiceVersionIdentity.String()
	}
	return id.ServiceVersionIdentity.String() + id.Variant.String()
}

func (id *ServiceVariantVersionIdentity) Parse(s string) error {
	if strings.HasSuffix(s, "}") {
		i := strings.Index(s, "{")
		if i < 0 {
			return fmt.Errorf("invalid service varaint version %q", s)
		}
		err := id.ServiceVersionIdentity.Parse(s[:i])
		if err != nil {
			return err
		}
		return id.Variant.Parse(s[i:])
	} else {
		return id.ServiceVersionIdentity.Parse(s)
	}
}

func (id ServiceVariantVersionIdentity) MarshalMapKey() (string, error) {
	return id.String(), nil
}

func (id *ServiceVariantVersionIdentity) UnmarshalMapKey(key string) error {
	return id.Parse(key)
}

func (id ServiceVariantVersionIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ServiceVariantVersionIdentity) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	id.Parse(s)
	return nil
}
