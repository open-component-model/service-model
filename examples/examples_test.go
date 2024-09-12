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
	MustBeSuccessfulWithOffset(1, refs.CheckLocalConsistency())
}

var _ = Describe("Examples", func() {
	Context("description files", func() {
		It("MSPGardener", func() {
			desc := Check(MSPGardener)
			CheckRefs(COMP_MSP_GARDENER, VERS_MSP_GARDENER, desc, `
  services:
    acme.org/gardener/service/installer:
      v1.0.0:
        "":
          references:
            description:
            - acme.org/gardener/service/provider:v1.0.0
    acme.org/gardener/service/provider:
      v1.0.0:
        "":
          references:
            description:
            - acme.org/gardener/apis/cluster:v1.22.0
            - acme.org/gardener/apis/cluster:v1.23.0
            installer:
            - acme.org/gardener/service/installer:v1.0.0
  usages:
    acme.org/gardener/apis/cluster:
      v1.22.0:
        "":
        - acme.org/gardener/service/provider:v1.0.0
      v1.23.0:
        "":
        - acme.org/gardener/service/provider:v1.0.0
    acme.org/gardener/service/installer:
      v1.0.0:
        "":
        - acme.org/gardener/service/provider:v1.0.0
    acme.org/gardener/service/provider:
      v1.0.0:
        "":
        - acme.org/gardener/service/installer:v1.0.0
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
        "":
          references:
            dependency:
            - acme.org/gardener/service/provider:v1.x.x
            description:
            - acme.org/hana/service/provider:v1.0.0
    acme.org/hana/service/provider:
      v1.0.0:
        "":
          references:
            dependency:
            - acme.org/gardener/service/provider:v1.x.x
            description:
            - acme.org/hana/apis/database:v1.5.0
            installer:
            - acme.org/hana/service/installer:v1.0.0
  usages:
    acme.org/gardener/service/provider:
      v1.x.x:
        "":
        - acme.org/hana/service/installer:v1.0.0
        - acme.org/hana/service/provider:v1.0.0
    acme.org/hana/apis/database:
      v1.5.0:
        "":
        - acme.org/hana/service/provider:v1.0.0
    acme.org/hana/service/installer:
      v1.0.0:
        "":
        - acme.org/hana/service/provider:v1.0.0
    acme.org/hana/service/provider:
      v1.0.0:
        "":
        - acme.org/hana/service/installer:v1.0.0
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
        "":
          references:
            dependency:
            - acme.org/gardener/service/provider:v1.x.x
            - acme.org/hana/service/provider:v1.x.x
            description:
            - acme.org/steampunk/service/provider:v1.0.0
    acme.org/steampunk/service/provider:
      v1.0.0:
        "":
          references:
            dependency:
            - acme.org/gardener/service/provider:v1.x.x
            - acme.org/hana/service/provider:v1.x.x
            description:
            - acme.org/gardener/apis/cluster
            - acme.org/hana/apis/database
            - acme.org/steampunk/apis/abap:v8.0.0
            installer:
            - acme.org/steampunk/service/installer:v1.0.0
  usages:
    acme.org/gardener/apis/cluster:
      "":
        "":
        - acme.org/steampunk/service/provider:v1.0.0
    acme.org/gardener/service/provider:
      v1.x.x:
        "":
        - acme.org/steampunk/service/installer:v1.0.0
        - acme.org/steampunk/service/provider:v1.0.0
    acme.org/hana/apis/database:
      "":
        "":
        - acme.org/steampunk/service/provider:v1.0.0
    acme.org/hana/service/provider:
      v1.x.x:
        "":
        - acme.org/steampunk/service/installer:v1.0.0
        - acme.org/steampunk/service/provider:v1.0.0
    acme.org/steampunk/apis/abap:
      v8.0.0:
        "":
        - acme.org/steampunk/service/provider:v1.0.0
    acme.org/steampunk/service/installer:
      v1.0.0:
        "":
        - acme.org/steampunk/service/provider:v1.0.0
    acme.org/steampunk/service/provider:
      v1.0.0:
        "":
        - acme.org/steampunk/service/installer:v1.0.0
`)
		})

		It("Abap System", func() {
			Check(SvcAbap)
		})
	})

})
