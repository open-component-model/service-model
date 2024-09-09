package internal

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/sliceutils"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

type Reference struct {
	v1.ServiceIdentity `json:",inline"`
	Version            string `json:"version"`
}

func NewReference(id v1.ServiceIdentity, vers string) *Reference {
	return &Reference{id, vers}
}

func (r Reference) String() string {
	if r.Version == "" {
		return r.ServiceIdentity.String()
	}
	return r.ServiceIdentity.String() + ":" + r.Version
}

func (r *Reference) Parse(s string) error {
	i := strings.LastIndex(s, ":")
	if i > 0 {
		r.ServiceIdentity.Parse(s[:i])
		r.Version = s[i+1:]
	} else {
		r.ServiceIdentity.Parse(s)
		r.Version = ""
	}
	return nil
}

func (r Reference) Equals(o Reference) bool {
	return r == o
}

////////////////////////////////////////////////////////////////////////////////

type References sliceutils.Slice[Reference]

func (r *References) Add(refs ...Reference) {
	*r = sliceutils.AppendUnique(*r, refs...)
}

func (r References) Len() int {
	return len(r)
}

func (r References) Less(i, j int) bool {
	return ReferenceCompare(r[i], r[j]) < 0
}

func (r References) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func ReferenceEquals(a, b Reference) bool {
	return a.Equals(b)
}

func ReferenceCompare(a, b Reference) int {
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

type UsageMap = marshallableMap[v1.ServiceIdentity, map[string]References, *v1.ServiceIdentity]

type ServiceEntry struct {
	References References `json:"references,omitempty"`
	Origin     Origin     `json:"origin,omitempty"`
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

func (c *CrossReferences) getService(holder *Reference) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity]
	if versions == nil {
		versions = map[string]*ServiceEntry{}
		c.Services[holder.ServiceIdentity] = versions
	}
	e := versions[holder.Version]
	if e == nil {
		e = &ServiceEntry{}
		versions[holder.Version] = e
	}
	return e
}

func (c *CrossReferences) AddService(holder *Reference, os ...Origin) {
	c.getService(holder).Origin = general.Optional(os...)
}

func (c *CrossReferences) AddRef(holder *Reference, ref *Reference) {
	// references
	entry := c.getService(holder)
	entry.References.Add(*ref)
	sort.Sort(entry.References)

	// usages
	versions := c.Usages[ref.ServiceIdentity]
	if versions == nil {
		versions = map[string]References{}
		c.Usages[ref.ServiceIdentity] = versions
	}

	holders := versions[ref.Version]
	holders.Add(*holder)
	sort.Sort(holders)
	versions[ref.Version] = holders
}

func (c *CrossReferences) AddRefs(a *CrossReferences) {
	for svc, versions := range a.Usages {
		for vers, e := range versions {
			for _, h := range e {
				c.AddService(&h, a.Services[h.ServiceIdentity][h.Version].Origin)
				c.AddRef(&h, NewReference(svc, vers))
			}
		}
	}
}
