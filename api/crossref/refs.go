package crossref

import (
	"fmt"
	"slices"
	"sort"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/jsonutils"
	"github.com/mandelsoft/goutils/maputils"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/identity"
)

type DepKind string

const (
	DEP_DEPENDS     DepKind = "depends"
	DEP_INSTANCE    DepKind = "instance"
	DEP_SATISFIES   DepKind = "satisfies"
	DEP_MANAGES     DepKind = "manages"
	DEP_INSTALLEDBY DepKind = "installedby"
	DEP_INSTALLS    DepKind = "installs" // beware of resulting dep cycles
	DEP_OTHER       DepKind = "other"
)

type (
	ServiceVersionVariantIdentity = identity.ServiceVersionVariantIdentity
	ServiceVersionIdentity        = identity.ServiceVersionIdentity
	ServiceVariantIdentity        = identity.ServiceVariantIdentity
	ServiceIdentity               = identity.ServiceIdentity
	Variant                       = identity.Variant
)

func NewServiceVersionIdentity(s identity.ServiceIdentity, vers string) ServiceVersionIdentity {
	return identity.NewServiceVersionId(s, vers)
}

func NewServiceVariantIdentity(s identity.ServiceIdentity, variant ...Variant) ServiceVariantIdentity {
	return identity.NewServiceVariantId(s, variant...)
}

func NewServiceVersionVariantIdentity(s identity.ServiceIdentity, vers string, variant ...identity.Variant) ServiceVersionVariantIdentity {
	return identity.NewServiceVersionVariantId(s, vers, variant...)
}

func NewServiceVersionVariantIdentityFor(svi ServiceVersionIdentity, variant ...identity.Variant) ServiceVersionVariantIdentity {
	return identity.NewServiceVersionVariantIdFor(svi, variant...)
}

type ServiceVersionVariantIdentities = identity.ServiceVersionVariantIdentities

////////////////////////////////////////////////////////////////////////////////

type Reference struct {
	Name         string
	Kind         DepKind
	Id           ServiceVariantIdentity
	Constaraints []string
}

func NewReference(name string, id identity.ServiceIdentity, vers string, variant identity.Variant, kind DepKind) *Reference {
	return &Reference{name, kind, NewServiceVariantIdentity(id, variant), sliceutils.AsSlice(vers)}
}

func (r Reference) Equals(o Reference) bool {
	return r.Name == o.Name && r.Kind == o.Kind &&
		r.Id.Equals(o.Id) &&
		slices.Equal(r.Constaraints, o.Constaraints)
}

type References = sliceutils.Slice[Reference]

func AddVersionReferences(refs *References, name string, id ServiceIdentity, variant identity.Variant, kind DepKind, versions ...string) {
	if len(versions) == 0 {
		refs.Add(*NewReference(name, id, "", variant, kind))
	} else {
		refs.Add(Reference{name, kind, NewServiceVariantIdentity(id, variant), versions})
	}
}

type UsageMap = jsonutils.MarshalableMap[identity.ServiceIdentity, map[string]map[string]ServiceVersionVariantIdentities, *identity.ServiceIdentity]

type ServiceEntry struct {
	References map[DepKind]ServiceVersionVariantIdentities `json:"references,omitempty"`
	Origin     identity.Origin                             `json:"origin,omitempty"`
	Variant    identity.Variant                            `json:"-"`
	Descriptor interface{}                                 `json:"-"`
}

type ServiceMap = jsonutils.MarshalableMap[identity.ServiceIdentity, map[string]map[string]*ServiceEntry, *identity.ServiceIdentity]

////////////////////////////////////////////////////////////////////////////////

type CrossReferences struct {
	Services ServiceMap `json:"services"`
	Usages   UsageMap   `json:"usages"`
}

func NewCrossReferences() *CrossReferences {
	return &CrossReferences{ServiceMap{}, UsageMap{}}
}

func (c *CrossReferences) getService(holder ServiceVersionIdentity, variant identity.Variant) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity()]
	if versions == nil {
		versions = map[string]map[string]*ServiceEntry{}
		c.Services[holder.ServiceIdentity()] = versions
	}
	variants := versions[holder.Version()]
	if variants == nil {
		variants = map[string]*ServiceEntry{}
		versions[holder.Version()] = variants
	}
	e := variants[variant.String()]
	if e == nil {
		e = &ServiceEntry{References: map[DepKind]ServiceVersionVariantIdentities{}}
		variants[variant.String()] = e
	}
	return e
}

func (c *CrossReferences) GetService(holder ServiceVersionIdentity, variant identity.Variant) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity()]
	if versions == nil {
		return nil
	}
	variants := versions[holder.Version()]
	if variants == nil {
		return nil
	}
	return variants[variant.String()]
}

func (c *CrossReferences) GetServiceVariants(holder ServiceVersionIdentity) []*ServiceEntry {
	versions := c.Services[holder.ServiceIdentity()]
	if versions == nil {
		return nil
	}
	variants := versions[holder.Version()]
	if variants == nil {
		return nil
	}
	return maputils.Values(variants)
}

func (c *CrossReferences) AddService(holder ServiceVersionIdentity, variant identity.Variant, desc any, os ...identity.Origin) {
	h := c.getService(holder, variant)
	if h.Origin == nil {
		h.Origin = general.Optional(os...)
	}
	if h.Variant == nil {
		h.Variant = variant.Copy()
	}
	if desc != nil {
		h.Descriptor = desc
	}
}

func (c *CrossReferences) AddRef(holder ServiceVersionIdentity, variant identity.Variant, ref ServiceVersionVariantIdentity, kind DepKind) {
	// references
	entry := c.getService(holder, variant)

	kindentry := entry.References[kind]
	kindentry.Add(ref)
	sort.Sort(kindentry)
	entry.References[kind] = kindentry

	// usages
	versions := c.Usages[ref.ServiceIdentity()]
	if versions == nil {
		versions = map[string]map[string]ServiceVersionVariantIdentities{}
		c.Usages[ref.ServiceIdentity()] = versions
	}

	variants := versions[ref.Version()]
	if variants == nil {
		variants = map[string]ServiceVersionVariantIdentities{}
		versions[ref.Version()] = variants
	}

	holders := variants[ref.Variant().String()]
	holders.Add(NewServiceVersionVariantIdentityFor(holder, variant))
	sort.Sort(holders)
	variants[ref.Variant().String()] = holders
}

func (c *CrossReferences) AddRefs(a *CrossReferences) {
	for svc, versions := range a.Services {
		for vers, variants := range versions {
			s := NewServiceVersionIdentity(svc, vers)
			for _, e := range variants {
				c.AddService(s, e.Variant, e.Origin)
				for k, l := range e.References {
					for _, r := range l {
						c.AddRef(s, e.Variant, r, k)
					}
				}
			}
		}
	}
}

func (c *CrossReferences) CheckLocalConsistency() error {
	var errlist errors.ErrorList

	for u, versions := range c.Usages {
		for v, variants := range versions {
			for vn, holders := range variants {
				var variant identity.Variant
				variant.Parse(vn)
				used := NewServiceVersionIdentity(u, v)
				usedentry := c.GetService(used, variant)
				for _, h := range holders {
					entry := c.GetService(h.ServiceVersionIdentity(), h.Variant())
					if entry != nil {
						if u.Component() == h.Component() && v == h.Version() {
							if usedentry == nil {
								errlist.Add(fmt.Errorf("missing service %s used by %s", used, h))
							}
						}
					}
				}
			}
		}
	}
	return errlist.Result()
}

func UniqueReferences(refs References) References {
	var result References
	for _, r := range refs {
		result = sliceutils.AppendUniqueFunc(result, general.EqualsFuncFor[Reference](), r)
	}
	return refs
}
