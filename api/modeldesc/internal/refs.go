package internal

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/goutils/sliceutils"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

type DepKind string

const (
	DEP_DEPENDENCY  DepKind = "dependency"
	DEP_DESCRIPTION DepKind = "description"
	DEP_MEET        DepKind = "meet"
	DEP_INSTALLER   DepKind = "installer"
)

type (
	ServiceVersionIdentity = v1.ServiceVersionIdentity
	ServiceIdentity        = v1.ServiceIdentity
)

func NewServiceVersionIdentity(s v1.ServiceIdentity, vers string) *ServiceVersionIdentity {
	return generics.Pointer(v1.NewServiceVersionId(s, vers))
}

type ServiceVersionIdentities = v1.ServiceVersionIdentities

////////////////////////////////////////////////////////////////////////////////

type Reference struct {
	Kind DepKind
	Id   ServiceVersionIdentity
}

func NewReference(id v1.ServiceIdentity, vers string, kind DepKind) *Reference {
	return &Reference{kind, *NewServiceVersionIdentity(id, vers)}
}

type References = sliceutils.Slice[Reference]

func AddVersionReferences(refs *References, id ServiceIdentity, kind DepKind, versions ...string) {
	if len(versions) == 0 {
		refs.Add(*NewReference(id, "", kind))
	} else {
		for _, e := range versions {
			refs.Add(*NewReference(id, e, kind))
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

type parsable interface {
	Parse(string) error
}

type parseablePointer[P any] interface {
	*P
	parsable
}

type stringable interface {
	comparable
	String() string
}

type marshallableMap[K stringable, V any, P parseablePointer[K]] map[K]V

var (
	_ json.Marshaler   = marshallableMap[v1.ServiceIdentity, int, *v1.ServiceIdentity](nil)
	_ json.Unmarshaler = (*marshallableMap[v1.ServiceIdentity, int, *v1.ServiceIdentity])(nil)
)

func (r marshallableMap[K, V, P]) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	for k, v := range r {
		data, err := json.Marshal(v)
		if err != nil {
			return nil, errors.Wrapf(err, "service %q", k.String())
		}
		m[k.String()] = data
	}
	return json.Marshal(m)
}

func (r *marshallableMap[K, V, P]) UnmarshalJSON(bytes []byte) error {
	var m map[string]json.RawMessage

	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}

	*r = marshallableMap[K, V, P]{}
	for k, v := range m {
		var s V
		var e K

		err := P(&e).Parse(k)
		if err != nil {
			return errors.Wrapf(err, "key %q", k)
		}

		err = json.Unmarshal(v, &s)
		if err != nil {
			return errors.Wrapf(err, "map entry %q", k)
		}
		(*r)[e] = s
	}
	return nil
}

type UsageMap = marshallableMap[v1.ServiceIdentity, map[string]ServiceVersionIdentities, *v1.ServiceIdentity]

type ServiceEntry struct {
	References map[DepKind]ServiceVersionIdentities `json:"references,omitempty"`
	Origin     Origin                               `json:"origin,omitempty"`
}

type ServiceMap = marshallableMap[v1.ServiceIdentity, map[string]*ServiceEntry, *v1.ServiceIdentity]

////////////////////////////////////////////////////////////////////////////////

type CrossReferences struct {
	Services ServiceMap `json:"services"`
	Usages   UsageMap   `json:"usages"`
}

func NewCrossReferences() *CrossReferences {
	return &CrossReferences{ServiceMap{}, UsageMap{}}
}

func (c *CrossReferences) getService(holder *ServiceVersionIdentity) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity]
	if versions == nil {
		versions = map[string]*ServiceEntry{}
		c.Services[holder.ServiceIdentity] = versions
	}
	e := versions[holder.Version]
	if e == nil {
		e = &ServiceEntry{References: map[DepKind]ServiceVersionIdentities{}}
		versions[holder.Version] = e
	}
	return e
}

func (c *CrossReferences) GetService(holder *ServiceVersionIdentity) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity]
	if versions == nil {
		return nil
	}
	return versions[holder.Version]
}

func (c *CrossReferences) AddService(holder *ServiceVersionIdentity, os ...Origin) {
	h := c.getService(holder)
	if h.Origin == nil {
		h.Origin = general.Optional(os...)
	}
}

func (c *CrossReferences) AddRef(holder *ServiceVersionIdentity, ref *ServiceVersionIdentity, kind DepKind) {
	// references
	entry := c.getService(holder)

	kindentry := entry.References[kind]
	kindentry.Add(*ref)
	sort.Sort(kindentry)
	entry.References[kind] = kindentry

	// usages
	versions := c.Usages[ref.ServiceIdentity]
	if versions == nil {
		versions = map[string]ServiceVersionIdentities{}
		c.Usages[ref.ServiceIdentity] = versions
	}

	holders := versions[ref.Version]
	holders.Add(*holder)
	sort.Sort(holders)
	versions[ref.Version] = holders
}

func (c *CrossReferences) AddRefs(a *CrossReferences) {
	for svc, versions := range a.Services {
		for vers, e := range versions {
			s := NewServiceVersionIdentity(svc, vers)
			c.AddService(s, e.Origin)
			for k, l := range e.References {
				for _, r := range l {
					c.AddRef(s, &r, k)
				}
			}
		}
	}
}

func (c *CrossReferences) CheckLocalConsistency() error {
	var errlist errors.ErrorList

	for u, versions := range c.Usages {
		for v, holders := range versions {
			used := NewServiceVersionIdentity(u, v)
			usedentry := c.GetService(used)
			for _, h := range holders {
				entry := c.GetService(&h)
				if entry != nil {
					if u.Component == h.Component && v == h.Version {
						if usedentry == nil {
							errlist.Add(fmt.Errorf("missing service %s used by %s", used, h))
						}
					}
				}
			}
		}
	}
	return errlist.Result()
}
