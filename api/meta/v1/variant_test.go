package v1_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/open-component-model/service-model/api/meta/v1"
)

func CheckVariant(v v1.Variant, res string) {
	key := Must(v.MarshalMapKey())
	ExpectWithOffset(1, key).To(Equal(res))

	var u v1.Variant
	MustBeSuccessfulWithOffset(1, u.UnmarshalMapKey(key))
	ExpectWithOffset(1, u).To(DeepEqual(v))
}

var _ = Describe("Variant Test Environment", func() {
	It("simple", func() {
		CheckVariant(
			v1.Variant{
				"alice": "bob",
			},
			"{alice=bob}",
		)
	})

	It("multi", func() {
		CheckVariant(
			v1.Variant{
				"alice": "bob",
				"petra": "tom",
			},
			"{alice=bob,petra=tom}",
		)
	})

	It("special,", func() {
		CheckVariant(
			v1.Variant{
				"alice,": "bob",
			},
			"{alice\\,=bob}",
		)
	})

	It("special=", func() {
		CheckVariant(
			v1.Variant{
				"alice=": "bob",
			},
			"{alice\\==bob}",
		)
	})

	It("special{}", func() {
		CheckVariant(
			v1.Variant{
				"alice{}": "bob",
			},
			"{alice{\\}=bob}",
		)
	})

	It("special=", func() {
		CheckVariant(
			v1.Variant{
				"alice\\": "bob",
			},
			"{alice\\\\=bob}",
		)
	})
})
