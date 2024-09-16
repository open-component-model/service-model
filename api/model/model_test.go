package model_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/model"
	"github.com/open-component-model/service-model/api/modeldesc/types/contract"
	. "ocm.software/ocm/api/helper/builder"
	"ocm.software/ocm/api/ocm/extensions/artifacttypes"
	"ocm.software/ocm/api/utils/accessobj"

	"ocm.software/ocm/api/ocm/extensions/repositories/ctf"
	"ocm.software/ocm/api/ocm/resolvers"
	"ocm.software/ocm/api/utils/mime"

	"github.com/open-component-model/service-model/api/ocm"
	"github.com/open-component-model/service-model/examples"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"

	"ocm.software/ocm/api/utils/accessio"
)

const ARCH = "ctf"

var _ = Describe("Test Environment", func() {
	var env *Builder
	var m model.Model

	BeforeEach(func() {
		env = NewBuilder()

		env.OCMCommonTransport(ARCH, accessio.FormatDirectory, func() {
			env.Component(examples.COMP_MSP_GARDENER, func() {
				env.Version(examples.VERS_MSP_GARDENER, func() {
					env.Resource("servicedesc", examples.VERS_MSP_GARDENER, ocm.RESOURCE_TYPE, ocmmeta.LocalRelation, func() {
						env.BlobStringData(mime.MIME_YAML, examples.MSPGardener)
					})

					env.Resource("installer", examples.VERS_MSP_GARDENER, artifacttypes.PLAIN_TEXT, ocmmeta.LocalRelation, func() {
						env.BlobStringData(mime.MIME_TEXT, "some installer definition")
					})
				})
			})
			env.Component(examples.COMP_CONTRACT_K8S_CLUSTER_22, func() {
				env.Version(examples.VERS_CONTRACT_K8S_CLUSTER_22, func() {
					env.Resource("servicedesc", examples.VERS_CONTRACT_K8S_CLUSTER_22, ocm.RESOURCE_TYPE, ocmmeta.LocalRelation, func() {
						env.BlobStringData(mime.MIME_YAML, examples.ContractK8sCluster22)
					})
				})
			})
		})
		repo := Must(ctf.Open(env, accessobj.ACC_READONLY, ARCH, 0, env))
		DeferCleanup(repo.Close)
		r := ocm.NewResolver(resolvers.ComponentResolverForRepository(repo))
		m = model.NewModel(r)
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Context("get from ocm", func() {
		It("read k8s cluster contract", func() {
			s := Must(m.GetServiceVersionVariant(identity.NewServiceVersionVariantIdentity(identity.NewServiceId(examples.COMP_CONTRACT_K8S_CLUSTER_22, examples.NAME_CONTRACT_K8S_CLUSTER_22), examples.VERS_CONTRACT_K8S_CLUSTER_22)))
			Expect(s).NotTo(BeNil())
			Expect(s.GetType()).To(Equal(contract.TYPE))
			Expect(s.GetComponent()).To(Equal(examples.COMP_CONTRACT_K8S_CLUSTER_22))
			Expect(s.GetVersion()).To(Equal(examples.VERS_CONTRACT_K8S_CLUSTER_22))
			Expect(s.GetVariant()).To(BeNil())
			Expect(s.GetName()).To(Equal(examples.NAME_CONTRACT_K8S_CLUSTER_22))
		})
	})
})
