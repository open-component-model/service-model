package contract_test

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "ocm.software/ocm/api/utils/testutils"

	"ocm.software/ocm/api/utils/runtime"

	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	me "github.com/open-component-model/service-model/api/modeldesc/types/contract"
	v1 "github.com/open-component-model/service-model/api/modeldesc/versions/v1"
)

var _ = Describe("V1 Test Environment", func() {
	version := v1.NewVersion(modeldesc.ABS_TYPE)

	data1 := `
type: ` + runtime.TypeName(modeldesc.REL_TYPE, "v1") + `
services:
- type: ` + me.TYPE + `
  service: service1
  shortName: a test service
  description: this service does nothing
  apiSpecificationType: dummy
  apiSpecificationVersion: v1
  specification: bla
`
	data2 := `
type: ` + runtime.TypeName(modeldesc.REL_TYPE, "v1") + `
services:
- type: ` + me.TYPE + `
  service: service1
  shortName: a test service
  description: this service does nothing
  apiSpecificationType: dummy
  apiSpecificationVersion: v1
  artifact:
    resource:
      name: dummy
`

	Context("serialization", func() {
		It("with spec", func() {

			desc := Must(version.Decode([]byte(data1), runtime.DefaultYAMLEncoding))
			Expect(len(desc.Services)).To(Equal(1))
			Expect(reflect.TypeOf(desc.Services[0].Kind)).To(Equal(reflect.TypeOf(&me.ServiceSpec{})))

			data := Must(modeldesc.Encode(desc))

			Expect(data).To(YAMLEqual(data1))
		})

		It("with artifact", func() {

			desc := Must(version.Decode([]byte(data2), runtime.DefaultYAMLEncoding))
			Expect(len(desc.Services)).To(Equal(1))
			Expect(reflect.TypeOf(desc.Services[0].Kind)).To(Equal(reflect.TypeOf(&me.ServiceSpec{})))

			data := Must(modeldesc.Encode(desc))

			Expect(data).To(YAMLEqual(data2))
		})
	})
})
