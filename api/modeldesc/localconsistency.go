package modeldesc

import (
	"github.com/open-component-model/service-model/api/crossref"
)

func NewCrossReferences() *CrossReferences {
	return crossref.NewCrossReferences()
}

func CheckLocalConsistency(desc *ServiceModelDescriptor, os ...Origin) (*CrossReferences, error) {
	refs := CrossReferencesFor(desc, os...)
	return refs, refs.CheckLocalConsistency()
}
