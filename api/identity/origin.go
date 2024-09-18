package identity

import (
	"maps"
	"reflect"

	"github.com/mandelsoft/goutils/maputils"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	"ocm.software/ocm/api/utils/misc"
)

type Origin map[string]map[string]string

func (o Origin) Copy() Origin {
	if o == nil {
		return nil
	}
	result := Origin{}
	for k, v := range o {
		result[k] = maps.Clone(v)
	}
	return result
}

func (o Origin) String() string {
	var s string
	for i, k := range maputils.OrderedKeys(o) {
		if i > 0 {
			s += ", "
		}
		s += k + "={"
		for j, a := range maputils.OrderedKeys(o[k]) {
			if j > 0 {
				s += ", "
			}
			s += a + "=" + o[k][a]
		}
		s += "}"
	}
	return s
}

func (o Origin) Equals(a Origin) bool {
	return reflect.DeepEqual(o, a)
}

const (
	ORIG_COMP     = "component"
	ORIG_VERS     = "version"
	ORIG_RESOURCE = "resource"

	ORIG_VALUE = "value"
)

func NewOCMOrigin(nv misc.VersionedElement, id ocmmeta.Identity) Origin {
	return Origin{
		ORIG_COMP:     {ORIG_VALUE: nv.GetName()},
		ORIG_VERS:     {ORIG_VALUE: nv.GetVersion()},
		ORIG_RESOURCE: id,
	}
}

func OriginEquals(a, b Origin) bool {
	return a.Equals(b)
}
