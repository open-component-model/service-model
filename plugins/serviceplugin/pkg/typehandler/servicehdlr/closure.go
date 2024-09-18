package servicehdlr

import (
	"github.com/Masterminds/semver/v3"
	"github.com/mandelsoft/goutils/set"
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/cli"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/api/utils/semverutils"
	"ocm.software/ocm/cmds/ocm/commands/common/options/closureoption"
	"ocm.software/ocm/cmds/ocm/common/output"
	"ocm.software/ocm/cmds/ocm/common/processing"
	"ocm.software/ocm/cmds/ocm/common/tree"
)

func DependencyLabel(o *tree.TreeObject) string {
	return string(o.Object.(*Object).Relation)
}

func ClosureExplode(opts *output.Options, e interface{}) []interface{} {
	list := traverse(common.History{}, e.(*Object), opts.Context, From(opts))
	return list
}

func traverse(hist common.History, o *Object, octx cli.Context, state *State) []interface{} {
	key := o.Key
	if err := hist.Add(modeldesc.KIND_SERVICEVERSION, key); err != nil {
		return nil
	}
	result := []interface{}{o}
	deps := crossref.UniqueReferences(o.Element.GetReferences())
	found := set.Set[string]{}
	for _, d := range deps {
		if !state.Relations.Has(d.Kind) {
			continue
		}
		key := resolve(&result, hist, d.Kind, d.Id.ServiceIdentity(), d.Id.Variant(), d.Constaraints, state)
		if key == nil {
			continue
		}

		if found.Has(key.String()) {
			continue // skip same ref with different attributes for recursion
		}
		found.Add(key.String())
		// TODO: provide error entry in list
		nested, err := state.Resolver.LookupServiceVersionVariant(*key)
		if err != nil {
			result = append(result, NewErrorObject(err, hist, d.Kind, d.Id.ServiceIdentity(), key.Version(), key.Variant()))
			continue
		}
		obj := NewObject(hist.Copy(), d.Kind, nested)
		if nested == nil {
			result = append(result, obj)
		} else {
			state.Add(obj.Element)
			result = append(result, traverse(hist, obj, octx, state)...)
		}
	}
	return result
}

func resolve(result *[]interface{}, hist common.History, label crossref.DepKind, s identity.ServiceIdentity, variant identity.Variant, versions []string, state *State) *identity.ServiceVersionVariantIdentity {
	var obj *Object

	if len(versions) != 1 {
		obj = NewConstraintObject(hist.Copy(), label, s, versions, variant)
	} else {
		key := identity.NewServiceVersionVariantId(s, versions[0], variant)
		if key.IsConstraint() {
			obj = NewConstraintObject(hist.Copy(), label, s, versions, variant)
		} else {
			return &key
		}
	}

	*result = append(*result, obj)
	if !state.ResolveToLatest {
		return nil
	}

	var constraints []*semver.Constraints
	for _, v := range versions {
		c, err := semver.NewConstraint(v)
		if err == nil {
			constraints = append(constraints, c)
		}
	}
	if len(constraints) == 0 {
		return nil
	}
	found, err := state.Resolver.ListVersions(s, variant)
	if err != nil {
		return nil
	}
	err = semverutils.SortVersions(found)
	if err != nil {
		return nil
	}
outer:
	for i := range found {
		v, err := semver.NewVersion(found[len(found)-1-i])
		if err == nil {
			for _, c := range constraints {
				if !c.Check(v) {
					continue outer
				}
			}
			obj.Resolved = found[len(found)-1-i]
			key := identity.NewServiceVersionVariantId(s, obj.Resolved, variant)
			return &key
		}
	}
	return nil
}

type NormalizeFunction processing.TransformFunction

func (c NormalizeFunction) Normalizer(opts *output.Options) processing.TransformFunction {
	if c != nil {
		copts := closureoption.From(opts)
		if copts != nil && copts.Closure {
			return c
		}
	}
	return nil
}
