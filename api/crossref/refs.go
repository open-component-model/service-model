package crossref

import (
	"fmt"
	"sort"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/goutils/jsonutils"
	"github.com/mandelsoft/goutils/maputils"
	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/common"
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
	ServiceVersionVariantIdentity = v1.ServiceVersionVariantIdentity
	ServiceVersionIdentity        = v1.ServiceVersionIdentity
	ServiceIdentity               = v1.ServiceIdentity
)

func NewServiceVersionIdentity(s v1.ServiceIdentity, vers string) *ServiceVersionIdentity {
	return generics.Pointer(v1.NewServiceVersionId(s, vers))
}

func NewServiceVersionVariantIdentity(s v1.ServiceIdentity, vers string, variant ...v1.Variant) ServiceVersionVariantIdentity {
	return v1.NewServiceVersionVariantIdentity(s, vers, variant...)
}

func NewServiceVersionVariantIdentityFor(svi ServiceVersionIdentity, variant ...v1.Variant) ServiceVersionVariantIdentity {
	return v1.NewServiceVersionVariantIdentityFor(svi, variant...)
}

type ServiceVersionVariantIdentities = v1.ServiceVersionVariantIdentities

////////////////////////////////////////////////////////////////////////////////

type Reference struct {
	Kind DepKind
	Id   ServiceVersionVariantIdentity
}

func NewReference(id v1.ServiceIdentity, vers string, variant v1.Variant, kind DepKind) *Reference {
	return &Reference{kind, NewServiceVersionVariantIdentity(id, vers, variant)}
}

type References = sliceutils.Slice[Reference]

func AddVersionReferences(refs *References, id ServiceIdentity, variant v1.Variant, kind DepKind, versions ...string) {
	if len(versions) == 0 {
		refs.Add(*NewReference(id, "", variant, kind))
	} else {
		for _, e := range versions {
			refs.Add(*NewReference(id, e, variant, kind))
		}
	}
}

type UsageMap = jsonutils.MarshalableMap[v1.ServiceIdentity, map[string]map[string]ServiceVersionVariantIdentities, *v1.ServiceIdentity]

type ServiceEntry struct {
	References map[DepKind]ServiceVersionVariantIdentities `json:"references,omitempty"`
	Origin     common.Origin                               `json:"origin,omitempty"`
	Variant    v1.Variant                                  `json:"-"`
	Descriptor interface{}                                 `json:"-"`
}

type ServiceMap = jsonutils.MarshalableMap[v1.ServiceIdentity, map[string]map[string]*ServiceEntry, *v1.ServiceIdentity]

////////////////////////////////////////////////////////////////////////////////

type CrossReferences struct {
	Services ServiceMap `json:"services"`
	Usages   UsageMap   `json:"usages"`
}

func NewCrossReferences() *CrossReferences {
	return &CrossReferences{ServiceMap{}, UsageMap{}}
}

func (c *CrossReferences) getService(holder *ServiceVersionIdentity, variant v1.Variant) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity]
	if versions == nil {
		versions = map[string]map[string]*ServiceEntry{}
		c.Services[holder.ServiceIdentity] = versions
	}
	variants := versions[holder.Version]
	if variants == nil {
		variants = map[string]*ServiceEntry{}
		versions[holder.Version] = variants
	}
	e := variants[variant.String()]
	if e == nil {
		e = &ServiceEntry{References: map[DepKind]ServiceVersionVariantIdentities{}}
		variants[variant.String()] = e
	}
	return e
}

func (c *CrossReferences) GetService(holder *ServiceVersionIdentity, variant v1.Variant) *ServiceEntry {
	versions := c.Services[holder.ServiceIdentity]
	if versions == nil {
		return nil
	}
	variants := versions[holder.Version]
	if variants == nil {
		return nil
	}
	return variants[variant.String()]
}

func (c *CrossReferences) GetServiceVariants(holder *ServiceVersionIdentity) []*ServiceEntry {
	versions := c.Services[holder.ServiceIdentity]
	if versions == nil {
		return nil
	}
	variants := versions[holder.Version]
	if variants == nil {
		return nil
	}
	return maputils.Values(variants)
}

func (c *CrossReferences) AddService(holder *ServiceVersionIdentity, variant v1.Variant, desc any, os ...common.Origin) {
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

func (c *CrossReferences) AddRef(holder *ServiceVersionIdentity, variant v1.Variant, ref *ServiceVersionVariantIdentity, kind DepKind) {
	// references
	entry := c.getService(holder, variant)

	kindentry := entry.References[kind]
	kindentry.Add(*ref)
	sort.Sort(kindentry)
	entry.References[kind] = kindentry

	// usages
	versions := c.Usages[ref.ServiceIdentity]
	if versions == nil {
		versions = map[string]map[string]ServiceVersionVariantIdentities{}
		c.Usages[ref.ServiceIdentity] = versions
	}

	variants := versions[ref.Version]
	if variants == nil {
		variants = map[string]ServiceVersionVariantIdentities{}
		versions[ref.Version] = variants
	}

	holders := variants[ref.Variant.String()]
	holders.Add(NewServiceVersionVariantIdentityFor(*holder, variant))
	sort.Sort(holders)
	variants[ref.Variant.String()] = holders
}

func (c *CrossReferences) AddRefs(a *CrossReferences) {
	for svc, versions := range a.Services {
		for vers, variants := range versions {
			s := NewServiceVersionIdentity(svc, vers)
			for _, e := range variants {
				c.AddService(s, e.Variant, e.Origin)
				for k, l := range e.References {
					for _, r := range l {
						c.AddRef(s, e.Variant, &r, k)
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
				var variant v1.Variant
				variant.Parse(vn)
				used := NewServiceVersionIdentity(u, v)
				usedentry := c.GetService(used, variant)
				for _, h := range holders {
					entry := c.GetService(&h.ServiceVersionIdentity, h.Variant)
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
	}
	return errlist.Result()
}
