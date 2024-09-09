package examples

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"ocm.software/ocm/api/utils/runtime"

	"github.com/open-component-model/service-model/api/modeldesc"
	common "ocm.software/ocm/api/utils/misc"
)

const COMP = "acme.org/test"
const VERS = "v1.0.0"

func Check(data string) *modeldesc.ServiceModelDescriptor {
	desc := MustWithOffset(1, Calling(modeldesc.Decode([]byte(data))))
	MustBeSuccessfulWithOffset(1, desc.Validate(common.NewNameVersion(COMP, VERS)))
	r, err := modeldesc.Encode(desc)
	ExpectWithOffset(1, r, err).To(YAMLEqual(data))
	return desc
}

func CheckRefs(comp, vers string, desc *modeldesc.ServiceModelDescriptor, exp string) {
	ctx := modeldesc.NewDescriptionContext(comp, vers, desc)
	canon := desc.ToCanonicalForm(ctx)
	refs := modeldesc.CrossReferencesFor(canon)

	data := Must(runtime.DefaultYAMLEncoding.Marshal(refs))
	ExpectWithOffset(1, data).To(YAMLEqual(exp))
}

const (
	COMP_MSP_GARDENER = "acme.org/gardener/service"
	VERS_MSP_GARDENER = "v1.0.0"
)

const (
	COMP_MSP_HANA = "acme.org/hana/service"
	VERS_MSP_HANA = "v1.0.0"
)

const (
	COMP_MSP_STEAMPUNK = "acme.org/steampunk/service"
	VERS_MSP_STEAMPUNK = "v1.0.0"
)

var _ = Describe("Examples", func() {
	Context("description files", func() {
		It("MSPGardener", func() {
			desc := Check(MSPGardener)
			CheckRefs(COMP_MSP_GARDENER, VERS_MSP_GARDENER, desc, `
  services:
    acme.org/gardener/service/installer:
      v1.0.0: {}
    acme.org/gardener/service/provider:
      v1.0.0:
        references:
        - acme.org/gardener/service/installer
  usages:
    acme.org/gardener/service/installer:
      v1.0.0:
      - acme.org/gardener/service/provider
`)
		})

		It("Garden Cluster", func() {
			Check(SvcGardenCluster22)
			Check(SvcGardenCluster23)
		})

		It("K8S Cluster", func() {
			Check(ContractK8sCluster22)
			Check(ContractK8sCluster23)
		})

		It("MSPHana", func() {
			desc := Check(MSPHana)
			CheckRefs(COMP_MSP_HANA, VERS_MSP_HANA, desc, `
  services:
    acme.org/hana/service/installer:
      v1.0.0:
        references:
        - acme.org/gardener/service/provider
    acme.org/hana/service/provider:
      v1.0.0:
        references:
        - acme.org/gardener/service/provider
        - acme.org/hana/service/installer
  usages:
    acme.org/gardener/service/provider:
      v1.x.x:
      - acme.org/hana/service/installer
      - acme.org/hana/service/provider
    acme.org/hana/service/installer:
      v1.0.0:
      - acme.org/hana/service/provider
`)
		})

		It("SvcHana", func() {
			Check(SvcHana)
		})

		It("MSPSteampunk", func() {
			desc := Check(MSPSteampunk)
			CheckRefs(COMP_MSP_STEAMPUNK, VERS_MSP_STEAMPUNK, desc, `
  services:
    acme.org/steampunk/service/installer:
      v1.0.0:
        references:
        - acme.org/gardener/service/provider
        - acme.org/hana/service/provider
    acme.org/steampunk/service/provider:
      v1.0.0:
        references:
        - acme.org/gardener/service/provider
        - acme.org/hana/service/provider
        - acme.org/steampunk/service/installer
  usages:
    acme.org/gardener/service/provider:
      v1.x.x:
      - acme.org/steampunk/service/installer
      - acme.org/steampunk/service/provider
    acme.org/hana/service/provider:
      v1.x.x:
      - acme.org/steampunk/service/installer
      - acme.org/steampunk/service/provider
    acme.org/steampunk/service/installer:
      v1.0.0:
      - acme.org/steampunk/service/provider
`)
		})

		It("Abap System", func() {
			Check(SvcAbap)
		})
	})

})
