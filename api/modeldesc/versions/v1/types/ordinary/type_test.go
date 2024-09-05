package ordinary_test

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "ocm.software/ocm/api/utils/testutils"

	"ocm.software/ocm/api/utils/runtime"

	modeldesc "github.com/open-component-model/service-model/api/modeldesc/internal"
	me "github.com/open-component-model/service-model/api/modeldesc/types/ordinary"
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
  variant:
     iaas: AWS
  inheritFrom:
     iaas: generic
  abstract: true
  labels:
    - name: dummy
      version: v1
      value: service
  external: true
  dependencies:
  - name: dep
    description: optional dependency to reporter instance creation
    service: reporter
    kind: implementation
    versionConstraints:
    - 1.x.x
    serviceInstances:
    - service: used
      versions:
      - 2.1.2
      dynamic: true
      static:
      - name: reporter
    optional: true
    labels:
    - name: test
      version: v1
      value: only for testing
`
	Context("serialization", func() {
		It("", func() {

			desc := Must(version.Decode([]byte(data1), runtime.DefaultYAMLEncoding))
			Expect(len(desc.Services)).To(Equal(1))
			Expect(reflect.TypeOf(desc.Services[0].Kind)).To(Equal(reflect.TypeOf(&me.ServiceSpec{})))

			data := Must(modeldesc.Encode(desc))

			Expect(data).To(YAMLEqual(data1))
		})
	})
})
