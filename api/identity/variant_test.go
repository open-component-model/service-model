package identity_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-component-model/service-model/api/identity"
)

func CheckVariant(v identity.Variant, res string) {
	key := Must(v.MarshalMapKey())
	ExpectWithOffset(1, key).To(Equal(res))

	var u identity.Variant
	MustBeSuccessfulWithOffset(1, u.UnmarshalMapKey(key))
	ExpectWithOffset(1, u).To(DeepEqual(v))
}

var _ = Describe("Variant Test Environment", func() {
	It("simple", func() {
		CheckVariant(
			identity.Variant{
				"alice": "bob",
			},
			"{alice=bob}",
		)
	})

	It("multi", func() {
		CheckVariant(
			identity.Variant{
				"alice": "bob",
				"petra": "tom",
			},
			"{alice=bob,petra=tom}",
		)
	})

	It("special,", func() {
		CheckVariant(
			identity.Variant{
				"alice,": "bob",
			},
			"{alice\\,=bob}",
		)
	})

	It("special=", func() {
		CheckVariant(
			identity.Variant{
				"alice=": "bob",
			},
			"{alice\\==bob}",
		)
	})

	It("special{}", func() {
		CheckVariant(
			identity.Variant{
				"alice{}": "bob",
			},
			"{alice{\\}=bob}",
		)
	})

	It("special=", func() {
		CheckVariant(
			identity.Variant{
				"alice\\": "bob",
			},
			"{alice\\\\=bob}",
		)
	})
})
