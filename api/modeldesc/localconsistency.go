package modeldesc

import (
	"github.com/open-component-model/service-model/api/modeldesc/internal"
)

func NewCrossReferences() *CrossReferences {
	return internal.NewCrossReferences()
}

func CheckLocalConsistency(desc *ServiceModelDescriptor, os ...Origin) (*CrossReferences, error) {
	refs := CrossReferencesFor(desc, os...)
	return refs, refs.CheckLocalConsistency()
}
