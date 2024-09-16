package ocm_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "ocm.software/ocm/api/helper/builder"

	"github.com/open-component-model/service-model/api/ocm"
	"github.com/open-component-model/service-model/examples"
	ocmmeta "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	"ocm.software/ocm/api/ocm/extensions/artifacttypes"
	"ocm.software/ocm/api/ocm/extensions/repositories/ctf"
	"ocm.software/ocm/api/ocm/resolvers"
	"ocm.software/ocm/api/utils/mime"

	"ocm.software/ocm/api/utils/accessio"
)

const ARCH = "ctf"

var _ = Describe("Test Environment", func() {
	var env *Builder

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
		})
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Context("get from ocm", func() {
		It("read gardener MSP", func() {
			repo := Must(ctf.Open(env, ctf.ACC_READONLY, ARCH, 0, env))
			defer Close(repo, "repo")

			resolver := resolvers.ComponentResolverForRepository(repo)
			desc, refs := Must2(ocm.GetServiceModel(examples.COMP_MSP_GARDENER, examples.VERS_MSP_GARDENER, resolvers.ComponentVersionResolverForComponentResolver(resolver)))
			Expect(desc).NotTo(BeNil())
			Expect(refs).NotTo(BeNil())
		})
	})
})
