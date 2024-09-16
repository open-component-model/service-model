package identity

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/utils"
	common "ocm.software/ocm/api/utils/misc"
)

type ServiceIdentity struct {
	component string
	name      string
}

func NewServiceId(component, name string) ServiceIdentity {
	return ServiceIdentity{component, name}
}

func (id ServiceIdentity) Name() string {
	return id.name
}

func (id ServiceIdentity) Component() string {
	return id.component
}

func (id ServiceIdentity) ForComponent(c string) ServiceIdentity {
	id.component = c
	return id
}

func (id ServiceIdentity) Validate() error {
	return utils.CheckFlatName(id.name, "service name")
}

func (id ServiceIdentity) IsRelative() bool {
	return id.component == ""
}

func (id ServiceIdentity) String() string {
	if id.component == "" {
		return id.name
	}
	return id.component + "/" + id.name
}

func (id *ServiceIdentity) Parse(s string) error {
	idx := strings.LastIndex(s, "/")
	if idx >= 0 {
		id.component = s[:idx]
		id.name = s[idx+1:]
	} else {
		id.name = s
		id.component = ""
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

type ServiceVariantIdentity struct {
	ServiceIdentity
	Variant Variant
}

func (id ServiceVariantIdentity) String() string {
	return id.ServiceIdentity.String() + id.Variant.String()
}

func (id *ServiceVariantIdentity) Parse(s string) error {
	var errlist errors.ErrorList
	idx := strings.LastIndex(s, "{")
	if idx >= 0 {
		errlist.Add(id.ServiceIdentity.Parse(s[:idx]))
		errlist.Add(id.Variant.Parse(s[idx:]))
	} else {
		errlist.Add(id.ServiceIdentity.Parse(s))
		id.Variant = nil
	}
	return errlist.Result()
}

////////////////////////////////////////////////////////////////////////////////

type _ServiceIdentity = ServiceIdentity

type ServiceVersionIdentity struct {
	_ServiceIdentity
	version string
}

func NewServiceVersionId(id ServiceIdentity, vers string) ServiceVersionIdentity {
	return ServiceVersionIdentity{id, vers}
}

func (id ServiceVersionIdentity) Version() string {
	return id.version
}

func (id ServiceVersionIdentity) ServiceIdentity() ServiceIdentity {
	return id._ServiceIdentity
}

func (id ServiceVersionIdentity) ComponentVersionId() common.NameVersion {
	return common.NewNameVersion(id.component, id.version)
}

func (id ServiceVersionIdentity) IsRelative() bool {
	return id.version == "" && id._ServiceIdentity.IsRelative()
}

func (id ServiceVersionIdentity) IsConstraint() bool {
	_, err := semver.NewVersion(id.version)
	return err != nil
}

func (r ServiceVersionIdentity) Equals(o ServiceVersionIdentity) bool {
	return r == o
}

func (id ServiceVersionIdentity) String() string {
	if id.version == "" {
		return id._ServiceIdentity.String()
	}
	return id._ServiceIdentity.String() + ":" + id.version
}

func (id *ServiceVersionIdentity) Parse(s string) error {
	var errlist errors.ErrorList
	idx := strings.LastIndex(s, ":")
	if idx >= 0 {
		errlist.Add(id._ServiceIdentity.Parse(s[:idx]))
		id.version = s[idx+1:]
	} else {
		errlist.Add(id._ServiceIdentity.Parse(s))
		id.version = ""
	}
	return errlist.Result()
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
	c := strings.Compare(a.component, b.component)
	if c == 0 {
		c = strings.Compare(a.name, b.name)
	}
	if c == 0 {
		c = strings.Compare(a.version, b.version)
	}
	return c
}

////////////////////////////////////////////////////////////////////////////////

type _ServiceVersionIdentity = ServiceVersionIdentity

type ServiceVersionVariantIdentity struct {
	_ServiceVersionIdentity
	variant Variant
}

func NewServiceVersionVariantIdentity(si ServiceIdentity, vers string, variant ...Variant) ServiceVersionVariantIdentity {
	return ServiceVersionVariantIdentity{NewServiceVersionId(si, vers), general.Optional(variant...)}
}

func NewServiceVersionVariantIdentityFor(svi ServiceVersionIdentity, variant ...Variant) ServiceVersionVariantIdentity {
	return ServiceVersionVariantIdentity{svi, general.Optional(variant...)}
}

func (id ServiceVersionVariantIdentity) Variant() Variant {
	return id.variant.Copy()
}

func (id ServiceVersionVariantIdentity) ServiceVersionIdentity() ServiceVersionIdentity {
	return id._ServiceVersionIdentity
}

func (id ServiceVersionVariantIdentity) GetServiceVariantName() string {
	if len(id.variant) == 0 {
		return id._ServiceIdentity.String()
	}
	return id._ServiceIdentity.String() + id.variant.String()
}

func (id ServiceVersionVariantIdentity) String() string {
	if len(id.variant) == 0 {
		return id._ServiceVersionIdentity.String()
	}
	return id._ServiceVersionIdentity.String() + id.variant.String()
}

func (id ServiceVersionVariantIdentity) Equals(o ServiceVersionVariantIdentity) bool {
	return id._ServiceVersionIdentity.Equals(o._ServiceVersionIdentity) &&
		id.variant.Equals(o.variant)
}

func (id *ServiceVersionVariantIdentity) Parse(s string) error {
	var errlist errors.ErrorList
	if strings.HasSuffix(s, "}") {
		i := strings.Index(s, "{")
		if i < 0 {
			return fmt.Errorf("invalid service varaint version %q", s)
		}
		errlist.Add(id._ServiceVersionIdentity.Parse(s[:i]))
		errlist.Add(id.variant.Parse(s[i:]))
	} else {
		errlist.Add(id._ServiceVersionIdentity.Parse(s))
	}
	return errlist.Result()
}

func (id ServiceVersionVariantIdentity) MarshalMapKey() (string, error) {
	return id.String(), nil
}

func (id *ServiceVersionVariantIdentity) UnmarshalMapKey(key string) error {
	return id.Parse(key)
}

func (id ServiceVersionVariantIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *ServiceVersionVariantIdentity) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	id.Parse(s)
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type ServiceVersionVariantIdentities sliceutils.Slice[ServiceVersionVariantIdentity]

func (r *ServiceVersionVariantIdentities) Add(refs ...ServiceVersionVariantIdentity) {
	*r = sliceutils.AppendUniqueFunc(*r, ServiceVersionVariantIdentityEquals, refs...)
}

func (r ServiceVersionVariantIdentities) Len() int {
	return len(r)
}

func (r ServiceVersionVariantIdentities) Less(i, j int) bool {
	return ServiceVersionVariantIdentityCompare(r[i], r[j]) < 0
}

func (r ServiceVersionVariantIdentities) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func ServiceVersionVariantIdentityEquals(a, b ServiceVersionVariantIdentity) bool {
	return a.Equals(b)
}

func ServiceVersionVariantIdentityCompare(a, b ServiceVersionVariantIdentity) int {
	c := ServiceVersionIdentityCompare(a._ServiceVersionIdentity, b._ServiceVersionIdentity)
	if c == 0 {
		c = strings.Compare(a.variant.String(), b.variant.String())
	}
	return c
}
