package typehandler

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mandelsoft/goutils/sliceutils"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/utils"
	"github.com/open-component-model/service-model/plugins/serviceplugin/pkg/typehdlrutils"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/runtime"
	"ocm.software/ocm/cmds/ocm/common/processing"
	"ocm.software/ocm/cmds/ocm/common/tree"
)

func Elem(e interface{}) *modeldesc.ServiceDescriptor {
	return e.(*Object).Element
}

type Objects = typehdlrutils.Objects[*Object]

type Manifest struct {
	History common.History  `json:"context,omitempty"`
	Element json.RawMessage `json:"element"`
}

type Object struct {
	History   common.History
	Sort      common.History
	Id        v1.ServiceVersionVariantIdentity
	Key       common.NameVersion
	HasNested bool
	Duplicate bool

	Element *modeldesc.ServiceDescriptor
	Node    *common.NameVersion
}

func NewObject(hist common.History, elem *modeldesc.ServiceDescriptor) *Object {
	id := v1.NewServiceVersionVariantIdentity(elem.Service, elem.Version, elem.Kind.GetVariant())
	nv := NewNameVersion(id.ServiceIdentity, id.Version, id.Variant)
	return &Object{
		History: hist,
		Sort:    sliceutils.AsSlice(nv),
		Id:      id,
		Key:     nv,
		Element: elem,
		Node:    &nv,
	}
}

func NewConstraintObject(hist common.History, sid v1.ServiceIdentity, constraints []string, variant ...v1.Variant) *Object {
	vers := strings.Join(constraints, ";")
	id := v1.NewServiceVersionVariantIdentity(sid, vers, variant...)
	nv := NewNameVersion(sid, vers, variant...)
	return &Object{
		History: hist,
		Sort:    sliceutils.AsSlice(nv),
		Id:      id,
		Key:     nv,
		Element: nil,
		Node:    &nv,
	}
}

func (o *Object) String() string {
	return fmt.Sprintf("history: %s, id: %s", o.History, o.Key)
}

func (o *Object) WithHistory(hist ...common.NameVersion) *Object {
	n := *o
	n.History = slices.Clone(hist)
	n.Sort = append(n.History, n.Sort[len(n.Sort)-1])
	return &n
}

var (
	_ common.HistorySource     = (*Object)(nil)
	_ tree.Object              = (*Object)(nil)
	_ typehdlrutils.NormObject = (*Object)(nil)
)

func (o *Object) AsManifest() interface{} {
	desc := &modeldesc.ServiceModelDescriptor{
		DocType:  runtime.NewVersionedObjectType(modeldesc.ABS_TYPE + "/v1"),
		Services: sliceutils.AsSlice(*o.Element),
	}
	data, _ := modeldesc.Encode(desc, runtime.DefaultJSONEncoding)
	return &Manifest{
		History: o.History,
		Element: data,
	}
}

func (o *Object) GetHistory() common.History {
	return o.History
}

func (o *Object) IsNode() *common.NameVersion {
	return o.Node
}

func (o *Object) IsValid() bool {
	return true
}

func (o *Object) GetKey() common.NameVersion {
	return NewNameVersion(o.Id.ServiceIdentity, o.Id.Version, o.Id.Variant)
}

func (o *Object) CreateContinue() typehdlrutils.NormObject {
	dummy := *o
	dummy.Node = nil
	dummy.Element = nil
	dummy.History = append(dummy.History, NewNameVersion(o.Id.ServiceIdentity, o.Id.Version, o.Id.Variant))
	return &dummy
}

func (o *Object) Compare(b *Object) int {
	return o.Sort.Compare(b.Sort)
}

// Sort is a processing chain sorting original objects provided by type handler.
var Sort = processing.Sort(utils.Compare[*Object])

func NewNameVersion(sid v1.ServiceIdentity, version string, variant ...v1.Variant) common.NameVersion {
	return common.NewNameVersion(v1.NewServiceVersionVariantIdentity(sid, version, variant...).GetServiceVariantName(), version)
}
