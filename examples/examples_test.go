package examples

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

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

var _ = Describe("Examples", func() {
	Context("description files", func() {
		It("MSPGardener", func() {
			Check(MSPGardener)
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
			Check(MSPHana)
		})

		It("SvcHana", func() {
			Check(SvcHana)
		})

		It("MSPSteampunk", func() {
			Check(MSPSteampunk)
		})

		It("Abap System", func() {
			Check(SvcAbap)
		})
	})
})
