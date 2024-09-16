package identity_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-component-model/service-model/api/identity"
)

func CheckVariantId(id identity.ServiceVersionVariantIdentity, res string) {
	Expect(id.String()).To(Equal(res))

	var u identity.ServiceVersionVariantIdentity
	MustBeSuccessful(u.Parse(id.String()))
	Expect(u).To(DeepEqual(id))
}

var _ = Describe("Identity Test Environment", func() {
	It("service variant version", func() {
		CheckVariantId(identity.NewServiceVersionVariantIdentity(identity.NewServiceId("acme.org/test", "provider"), "v1",
			identity.Variant{"iaas": "AWS"}),
			"acme.org/test/provider:v1{iaas=AWS}",
		)
	})

	It("service variant version without variant", func() {
		CheckVariantId(identity.NewServiceVersionVariantIdentity(identity.NewServiceId("acme.org/test", "provider"), "v1"),
			"acme.org/test/provider:v1",
		)
	})

	It("service variant version without version", func() {
		CheckVariantId(identity.NewServiceVersionVariantIdentity(identity.NewServiceId("acme.org/test", "provider"), "",
			identity.Variant{"iaas": "AWS"}),
			"acme.org/test/provider{iaas=AWS}",
		)
	})

	It("service variant version without version and variant", func() {
		CheckVariantId(identity.NewServiceVersionVariantIdentity(identity.NewServiceId("acme.org/test", "provider"), ""),
			"acme.org/test/provider",
		)
	})
})
