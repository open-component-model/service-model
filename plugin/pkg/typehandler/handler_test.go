package typehandler_test

import (
	"github.com/mandelsoft/goutils/sliceutils"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "github.com/open-component-model/service-model/api/meta/v1"
	. "github.com/open-component-model/service-model/examples"
	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/api/ocm/extensions/artifacttypes"
	"ocm.software/ocm/cmds/ocm/common/output"
	"ocm.software/ocm/cmds/ocm/common/utils"
	"ocm.software/ocm/cmds/ocm/testhelper"

	ocmdesc "github.com/open-component-model/service-model/api/ocm"
	"github.com/open-component-model/service-model/plugin/pkg/typehandler"
	v1 "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	"ocm.software/ocm/api/ocm/extensions/repositories/ctf"
	"ocm.software/ocm/api/ocm/resolvers"
	"ocm.software/ocm/api/utils/accessio"
	"ocm.software/ocm/api/utils/mime"
)

const ARCH = "arch.ctf"

var _ = Describe("Handler Test Environment", func() {

	var env *testhelper.TestEnv
	var repo ocm.Repository
	var resolver resolvers.ComponentResolver

	BeforeEach(func() {
		env = testhelper.NewTestEnv()

		env.OCMCommonTransport(ARCH, accessio.FormatDirectory, func() {
			env.ComponentVersion(COMP_MSP_GARDENER, VERS_MSP_GARDENER, func() {
				env.Resource("service", VERS_MSP_GARDENER, ocmdesc.RESOURCE_TYPE, v1.LocalRelation, func() {
					env.BlobStringData(mime.MIME_YAML, MSPGardener)
				})
				env.Resource("installer", VERS_MSP_GARDENER, artifacttypes.PLAIN_TEXT, v1.LocalRelation, func() {
					env.BlobStringData(mime.MIME_TEXT, "some installer description")
				})
			})
		})

		repo = Must(ctf.Open(env, ctf.ACC_READONLY, ARCH, 0, env))
		resolver = resolvers.ComponentResolverForRepository(repo)
	})

	AfterEach(func() {
		Close(repo)
		env.Cleanup()
	})

	Context("service", func() {
		It("dedicated service", func() {
			h := typehandler.ForServices(resolver)

			list := Must(h.Get(utils.StringSpec(metav1.NewServiceVersionVariantIdentity(metav1.NewServiceId(COMP_MSP_GARDENER, "provider"), VERS_MSP_GARDENER).String())))
			Expect(len(list)).To(Equal(1))
			Expect(typehandler.Elem(list[0]).Service).To(Equal(metav1.NewServiceId(COMP_MSP_GARDENER, "provider")))
			Expect(typehandler.Elem(list[0]).Version).To(Equal(VERS_MSP_GARDENER))
		})
	})

	Context("component", func() {
		It("dedicated component", func() {
			sess := ocm.NewSession(nil)
			defer Close(sess, "session")
			h := Must(typehandler.ForComponents(env.OCM(), resolver, &output.Options{}, repo, sess, sliceutils.AsSlice(COMP_MSP_GARDENER)))

			list := Must(h.All())
			Expect(len(list)).To(Equal(2))
			Expect(typehandler.Elem(list[0]).Service).To(Equal(metav1.NewServiceId(COMP_MSP_GARDENER, "provider")))
			Expect(typehandler.Elem(list[0]).Version).To(Equal(VERS_MSP_GARDENER))
			Expect(typehandler.Elem(list[1]).Service).To(Equal(metav1.NewServiceId(COMP_MSP_GARDENER, "installer")))
			Expect(typehandler.Elem(list[1]).Version).To(Equal(VERS_MSP_GARDENER))
		})

		It("dedicated service in component", func() {
			sess := ocm.NewSession(nil)
			defer Close(sess, "session")
			h := Must(typehandler.ForComponents(env.OCM(), resolver, &output.Options{}, repo, sess, sliceutils.AsSlice(COMP_MSP_GARDENER)))

			list := Must(h.Get(utils.StringSpec("provider")))
			Expect(len(list)).To(Equal(1))
			Expect(typehandler.Elem(list[0]).Service).To(Equal(metav1.NewServiceId(COMP_MSP_GARDENER, "provider")))
			Expect(typehandler.Elem(list[0]).Version).To(Equal(VERS_MSP_GARDENER))
		})
	})
})
