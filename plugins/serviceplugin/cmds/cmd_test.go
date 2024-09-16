//go:build unix

package cmds_test

import (
	"bytes"

	. "github.com/mandelsoft/goutils/testutils"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/projectionfs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ocmdesc "github.com/open-component-model/service-model/api/ocm"
	mutils "github.com/open-component-model/service-model/api/utils"
	. "github.com/open-component-model/service-model/examples"
	v1 "ocm.software/ocm/api/ocm/compdesc/meta/v1"
	"ocm.software/ocm/api/ocm/extensions/artifacttypes"
	. "ocm.software/ocm/api/ocm/plugin/testutils"
	"ocm.software/ocm/api/utils/accessio"
	"ocm.software/ocm/api/utils/mime"
	. "ocm.software/ocm/cmds/ocm/testhelper"

	"ocm.software/ocm/api/ocm/extensions/attrs/plugincacheattr"
)

const ARCH = "arch.ctf"

var _ = Describe("cliplugin", func() {
	Context("lib", func() {
		var env *TestEnv
		var plugins TempPluginDir
		var basepath string

		BeforeEach(func() {
			tmpfs := Must(osfs.NewTempFileSystem())
			basepath = projectionfs.Root(tmpfs)
			env = NewTestEnv(FileSystem(tmpfs))
			plugins = Must(ConfigureTestPlugins(env, "testdata/plugins"))

			registry := plugincacheattr.Get(env)
			//	Expect(registration.RegisterExtensions(env)).To(Succeed())
			p := registry.Get("serviceplugin")
			Expect(p).NotTo(BeNil())
			Expect(p.Error()).To(Equal(""))

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
		})

		AfterEach(func() {
			plugins.Cleanup()
			env.Cleanup()
		})

		It("run plugin based ocm command", func() {
			var buf bytes.Buffer

			MustBeSuccessful(env.CatchOutput(&buf).Execute("get", "services", "--repo", basepath+"/"+ARCH, COMP_MSP_GARDENER+"/provider"))

			Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  COMPONENT                 NAME     VERSION VARIANT KIND            SHORTNAME
  acme.org/gardener/service provider v1.0.0          ServiceProvider Gardener Kubernetes as a Service Management

`, 2)))
		})

		It("run plugin based ocm command with closure", func() {
			var buf bytes.Buffer

			MustBeSuccessful(env.CatchOutput(&buf).Execute("get", "services", "-r", "--repo", basepath+"/"+ARCH, COMP_MSP_GARDENER+"/provider"))

			Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  REFERENCEPATH                             COMPONENT                 NAME      VERSION VARIANT KIND                SHORTNAME
                                            acme.org/gardener/service provider  v1.0.0          ServiceProvider     Gardener Kubernetes as a Service Management
  acme.org/gardener/service/provider:v1.0.0 acme.org/gardener/service installer v1.0.0          InstallationService Installer for Gardener
`, 2)))
		})

		It("run plugin based ocm command with closure tree", func() {
			var buf bytes.Buffer

			MustBeSuccessful(env.CatchOutput(&buf).Execute("get", "services", "-otree", "-r", "--repo", basepath+"/"+ARCH, COMP_MSP_GARDENER+"/provider"))

			Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  NESTING COMPONENT                 NAME      VERSION VARIANT KIND                SHORTNAME
  └─ ⊗    acme.org/gardener/service provider  v1.0.0          ServiceProvider     Gardener Kubernetes as a Service Management
     └─   acme.org/gardener/service installer v1.0.0          InstallationService Installer for Gardener
`, 2)))
		})

		It("run plugin based ocm command with closure yaml", func() {
			var buf bytes.Buffer

			MustBeSuccessful(env.CatchOutput(&buf).Execute("get", "services", "-oyaml", "-r", "--repo", basepath+"/"+ARCH, COMP_MSP_GARDENER+"/provider"))

			Expect(buf.String()).To(StringEqualTrimmedWithContext(mutils.Crop(`
  ---
  element:
    services:
    - installers:
      - service: acme.org/gardener/service/installer
        version: v1.0.0
      managedServices:
      - service: acme.org/gardener/apis/cluster
        versions:
        - v1.22.0
        - v1.23.0
      service: acme.org/gardener/service/provider
      shortName: Gardener Kubernetes as a Service Management
      type: ServiceProvider
      version: v1.0.0
    type: serviceModelDescription/v1
  ---
  context:
  - acme.org/gardener/service/provider:v1.0.0
  element:
    services:
    - installedService: acme.org/gardener/service/provider
      installerResource:
        resource:
          name: installer
      installerType: Deplomat
      service: acme.org/gardener/service/installer
      shortName: Installer for Gardener
      targetEnvironment:
        type: KubernetesCluster
      type: InstallationService
      version: v1.0.0
      versions:
      - v1.0.0
    type: serviceModelDescription/v1
`, 2)))
		})
	})
})
