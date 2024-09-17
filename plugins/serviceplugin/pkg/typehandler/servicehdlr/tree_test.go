package servicehdlr_test

import (
	"bytes"
	"fmt"

	"github.com/mandelsoft/goutils/generics"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/goutils/sliceutils"
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/identity"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/modeldesc/types/contract"
	"github.com/open-component-model/service-model/api/modeldesc/types/installer"
	"github.com/open-component-model/service-model/api/modeldesc/types/ordinary"
	"github.com/open-component-model/service-model/api/modeldesc/types/provider"
	"github.com/open-component-model/service-model/api/modeldesc/vpi"
	mutils "github.com/open-component-model/service-model/api/utils"
	"github.com/open-component-model/service-model/plugins/serviceplugin/pkg/typehandler"
	"github.com/open-component-model/service-model/plugins/serviceplugin/pkg/typehandler/servicehdlr"
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/cmds/ocm/common/output"
	"ocm.software/ocm/cmds/ocm/common/processing"
	"ocm.software/ocm/cmds/ocm/common/utils"
)

var COMP_S1 = "acme.com/s1"
var COMP_S2 = "acme.com/s2"
var COMP_S3 = "acme.com/s3"
var VERS = "v1.0.0"

var s11 = identity.NewServiceId(COMP_S1, "provider")
var s12 = identity.NewServiceId(COMP_S1, "installer")
var s2 = identity.NewServiceId(COMP_S2, "service")
var s3 = identity.NewServiceId(COMP_S3, "service")

type TypeHandler struct {
	objs []*servicehdlr.Object
}

func NewTypeHandler(objs []*servicehdlr.Object) (utils.TypeHandler, error) {
	t := &TypeHandler{
		objs: objs,
	}
	return t, nil
}

func (t *TypeHandler) All() ([]output.Object, error) {
	return sliceutils.AsAny(t.objs), nil
}

func (t *TypeHandler) Get(name utils.ElemSpec) ([]output.Object, error) {
	return nil, nil
}

func (t *TypeHandler) Close() error {
	return nil
}

type Kind struct {
	kind string
}

var _ modeldesc.ServiceKindSpec = (*Kind)(nil)

func (k *Kind) GetType() string {
	return k.kind
}

func (k *Kind) GetVariant() identity.Variant {
	return nil
}

func (k *Kind) GetDependencies() []metav1.Dependency {
	return nil
}

func (k *Kind) Copy() modeldesc.ServiceKindSpec {
	return generics.Pointer(*k)
}

func (k *Kind) ToCanonicalForm(c modeldesc.DescriptionContext) modeldesc.ServiceKindSpec {
	panic("implement me")
}

func (k *Kind) Validate(c modeldesc.DescriptionContext) error {
	panic("implement me")
}

func (k *Kind) GetReferences() crossref.References {
	panic("implement me")
}

func NewObject(sid identity.ServiceIdentity, vers string, kind string, short string) *servicehdlr.Object {
	return servicehdlr.NewObject(nil, &modeldesc.ServiceDescriptor{
		CommonServiceSpec: vpi.CommonServiceSpec{
			Service:   sid,
			Version:   vers,
			ShortName: short,
		},
		Kind: &Kind{kind},
	})
}

func TableOutput(opts *output.Options, mapping processing.MappingFunction, wide ...string) *output.TableOutput {
	return &output.TableOutput{
		Headers: output.Fields("COMPONENT", "NAME", "VERSION", "VARIANT", "KIND", "SHORTNAME", wide),
		Options: opts,
		Chain:   processing.Append(servicehdlr.Sort, typehandler.Normalize),
		Mapping: mapping,
	}
}

func getRegular(opts *output.Options) output.Output {
	return TableOutput(opts, mapGetRegularOutput).New()
}

func getTree(opts *output.Options) output.Output {
	return output.TreeOutput(TableOutput(opts, mapGetRegularOutput), "NESTING").New()
}

func mapGetRegularOutput(e interface{}) interface{} {
	obj := e.(*servicehdlr.Object)
	if obj.Node == nil {
		return sliceutils.AsSlice("...", "", "", "", "", "")
	}
	r := obj.Element
	if r == nil {
		res := ""
		if obj.Resolved != "" {
			res = fmt.Sprintf("(resolved to %s)", obj.Resolved)
		}
		return sliceutils.AsSlice(obj.Id.Component(), obj.Id.Name(), obj.Id.Version(), obj.Id.Variant().String(), "", res)
	}
	return sliceutils.AsSlice(r.Service.Component(), r.Service.Name(), r.Version, r.Kind.GetVariant().String(), r.Kind.GetType(), r.ShortName)
}

var _ = Describe("TreeTest Environment", func() {
	var buf *bytes.Buffer
	var ctx cli.Context

	BeforeEach(func() {
		buf = bytes.NewBuffer(nil)
		ctx = cli.New().WithStdIO(nil, buf, buf)
	})

	It("simple list", func() {
		objs := []*servicehdlr.Object{
			NewObject(s11, VERS, provider.TYPE, "gardener"),
			NewObject(s12, VERS, installer.TYPE, "gardener installer"),
		}
		hdlr := Must(NewTypeHandler(objs))

		opts := &output.Options{Context: ctx}

		opts.Output = getRegular(opts)
		MustBeSuccessful(utils.HandleOutput(opts.Output, hdlr))
		Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  COMPONENT   NAME      VERSION VARIANT KIND                SHORTNAME
  acme.com/s1 installer v1.0.0          InstallationService gardener installer
  acme.com/s1 provider  v1.0.0          ServiceProvider     gardener
`, 2)))
	})

	It("simple usage tree", func() {
		objs := []*servicehdlr.Object{
			NewObject(s11, VERS, provider.TYPE, "gardener"),
			NewObject(s12, VERS, installer.TYPE, "gardener installer"),
			NewObject(s12, VERS, installer.TYPE, "gardener installer").WithHistory(servicehdlr.NewNameVersion(s11, VERS)),
		}
		hdlr := Must(NewTypeHandler(objs))

		opts := &output.Options{Context: ctx}

		opts.Output = getTree(opts)
		MustBeSuccessful(utils.HandleOutput(opts.Output, hdlr))
		Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  NESTING COMPONENT   NAME      VERSION VARIANT KIND                SHORTNAME
  └─ ⊗    acme.com/s1 provider  v1.0.0          ServiceProvider     gardener
     └─   acme.com/s1 installer v1.0.0          InstallationService gardener installer
`, 2)))
	})

	It("complex usage tree", func() {
		objs := []*servicehdlr.Object{
			NewObject(s12, VERS, installer.TYPE, "gardener installer"),
			NewObject(s2, VERS, ordinary.TYPE, "common").WithHistory(servicehdlr.NewNameVersion(s12, VERS)),
			NewObject(s3, VERS, contract.TYPE, "common contract").WithHistory(servicehdlr.NewNameVersion(s12, VERS), servicehdlr.NewNameVersion(s2, VERS)),

			NewObject(s11, VERS, provider.TYPE, "gardener"),
			NewObject(s12, VERS, installer.TYPE, "gardener installer").WithHistory(servicehdlr.NewNameVersion(s11, VERS)),
			NewObject(s2, VERS, ordinary.TYPE, "common").WithHistory(servicehdlr.NewNameVersion(s11, VERS), servicehdlr.NewNameVersion(s12, VERS)),
			NewObject(s3, VERS, contract.TYPE, "common contract").WithHistory(servicehdlr.NewNameVersion(s11, VERS), servicehdlr.NewNameVersion(s12, VERS), servicehdlr.NewNameVersion(s2, VERS)),

			NewObject(s2, VERS, ordinary.TYPE, "commom").WithHistory(servicehdlr.NewNameVersion(s11, VERS)),
			NewObject(s3, VERS, contract.TYPE, "common contract").WithHistory(servicehdlr.NewNameVersion(s11, VERS), servicehdlr.NewNameVersion(s2, VERS)),

			NewObject(s2, VERS, ordinary.TYPE, "commom"),
			NewObject(s3, VERS, contract.TYPE, "common contract").WithHistory(servicehdlr.NewNameVersion(s2, VERS)),

			NewObject(s3, VERS, contract.TYPE, "common contract"),
		}

		hdlr := Must(NewTypeHandler(objs))

		opts := &output.Options{Context: ctx}

		opts.Output = getTree(opts)
		MustBeSuccessful(utils.HandleOutput(opts.Output, hdlr))
		Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  NESTING     COMPONENT   NAME      VERSION VARIANT KIND                SHORTNAME
  └─ ⊗        acme.com/s1 provider  v1.0.0          ServiceProvider     gardener
     ├─ ⊗     acme.com/s1 installer v1.0.0          InstallationService gardener installer
     │  └─ ⊗  acme.com/s2 service   v1.0.0          OrdinaryService     common
     │     └─ ...                                                       
     └─ ⊗     acme.com/s2 service   v1.0.0          OrdinaryService     commom
        └─    acme.com/s3 service   v1.0.0          ServiceContract     common contract
`, 2)))
	})

})
