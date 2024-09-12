package v1_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

func CheckVariantId(id v1.ServiceVersionVariantIdentity, res string) {
	Expect(id.String()).To(Equal(res))

	var u v1.ServiceVersionVariantIdentity
	MustBeSuccessful(u.Parse(id.String()))
	Expect(u).To(DeepEqual(id))
}

var _ = Describe("Identity Test Environment", func() {
	It("service variant version", func() {
		CheckVariantId(v1.NewServiceVersionVariantIdentity(v1.NewServiceId("acme.org/test", "provider"), "v1",
			v1.Variant{"iaas": "AWS"}),
			"acme.org/test/provider:v1{iaas=AWS}",
		)
	})

	It("service variant version without variant", func() {
		CheckVariantId(v1.NewServiceVersionVariantIdentity(v1.NewServiceId("acme.org/test", "provider"), "v1"),
			"acme.org/test/provider:v1",
		)
	})

	It("service variant version without version", func() {
		CheckVariantId(v1.NewServiceVersionVariantIdentity(v1.NewServiceId("acme.org/test", "provider"), "",
			v1.Variant{"iaas": "AWS"}),
			"acme.org/test/provider{iaas=AWS}",
		)
	})

	It("service variant version without version and variant", func() {
		CheckVariantId(v1.NewServiceVersionVariantIdentity(v1.NewServiceId("acme.org/test", "provider"), ""),
			"acme.org/test/provider",
		)
	})
})
