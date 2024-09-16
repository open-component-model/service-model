package servicehdlr

import (
	"github.com/mandelsoft/goutils/set"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/spf13/pflag"
	"ocm.software/ocm/api/cli"
	common "ocm.software/ocm/api/utils/misc"
	"ocm.software/ocm/cmds/ocm/commands/common/options/closureoption"
	"ocm.software/ocm/cmds/ocm/common/options"
	"ocm.software/ocm/cmds/ocm/common/output"
	"ocm.software/ocm/cmds/ocm/common/processing"
)

func from(o options.OptionSetProvider) *State {
	var opt *State
	o.AsOptionSet().Get(&opt)
	return opt
}

func NewState(r modeldesc.Resolver) *State {
	return &State{r}
}

type State struct {
	Resolver modeldesc.Resolver
}

func (o *State) AddFlags(fs *pflag.FlagSet) {
	// fake option to pass state
}

func ClosureExplode(opts *output.Options, e interface{}) []interface{} {
	return traverse(common.History{}, e.(*Object), opts.Context, from(opts))
}

func traverse(hist common.History, o *Object, octx cli.Context, state *State) []interface{} {
	key := o.Key
	if err := hist.Add(modeldesc.KIND_SERVICEVERSION, key); err != nil {
		return nil
	}
	result := []interface{}{o}
	deps := o.Element.Kind.GetDependencies()
	found := set.Set[string]{}
	for _, d := range deps {
		if len(d.VersionConstraints) != 1 {
			obj := NewConstraintObject(hist.Copy(), d.Service, d.VersionConstraints, d.Variant)
			result = append(result, obj)
			continue // cannot traverse unconcrete deps
		}
		key := identity.NewServiceVersionVariantIdentity(d.Service, d.VersionConstraints[0], d.Variant)
		if key.IsConstraint() {
			obj := NewConstraintObject(hist.Copy(), d.Service, d.VersionConstraints, d.Variant)
			result = append(result, obj)
			continue
		}
		if found.Has(key.String()) {
			continue // skip same ref with different attributes for recursion
		}
		found.Add(key.String())
		// TODO: provide error entry in list
		nested, err := state.Resolver.LookupServiceVersionVariant(key)
		if err != nil {
			result = append(result, NewErrorObject(err, hist, d.Service, key.Version(), key.Variant()))
			continue
		}
		obj := NewObject(hist.Copy(), nested)
		if nested == nil {
			result = append(result, obj)
		} else {
			result = append(result, traverse(hist, obj, octx, state)...)
		}
	}
	return result
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
