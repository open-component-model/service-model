package examples

import (
	"ocm.software/ocm/api/helper/env"
)

func TestData(dest ...string) env.Option {
	return env.ProjectTestDataForCaller("descriptors", dest...)
}

func ModifiableTestData(dest ...string) env.Option {
	return env.ModifiableProjectTestDataForCaller("descriptors", dest...)
}
