package servicehdlr

import (
	"github.com/Masterminds/semver/v3"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"ocm.software/ocm/api/ocm/resolvers"

	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/handlers/comphdlr"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/options/lookupoption"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/options/versionconstraintsoption"
	"ocm.software/ocm/cmds/ocm/common/options"
)

type Option interface {
	ApplyToServiceHandlerOptions(handler *Options)
}

type Options struct {
	forceEmpty bool

	state *State

	constraints []*semver.Constraints
	latestOnly  *bool
	resolver    resolvers.ComponentVersionResolver
	repo        ocm.Repository
}

func (o *Options) ApplyToServiceHandlerOptions(opts *Options) {
	if o.constraints != nil {
		opts.constraints = o.constraints
	}
	if o.latestOnly != nil {
		opts.latestOnly = o.latestOnly
	}
	if o.resolver != nil {
		opts.resolver = o.resolver
	}
}

////////////////////////////////////////////////////////////////////////////////

type OptionList []Option

func (o OptionList) ApplyToServiceHandlerOptions(handler *Options) {
	for _, e := range o {
		e.ApplyToServiceHandlerOptions(handler)
	}
}

func OptionsFor(o options.OptionSetProvider) OptionList {
	var hopts []Option
	constr := versionconstraintsoption.From(o)
	if len(constr.Constraints) > 0 {
		hopts = append(hopts, WithVersionConstraints(constr.Constraints))
	}
	if constr.Latest {
		hopts = append(hopts, LatestOnly())
	}
	if lookup := lookupoption.From(o); lookup != nil {
		hopts = append(hopts, Resolver(lookup))
	}
	if state := From(o); state != nil {
		hopts = append(hopts, state)
	}
	return hopts
}

////////////////////////////////////////////////////////////////////////////////

type repo struct {
	repo ocm.Repository
}

func (o repo) ApplyToServiceHandlerOptions(opts *Options) {
	opts.repo = o.repo
}

func Repository(r ocm.Repository) Option {
	return repo{r}
}

////////////////////////////////////////////////////////////////////////////////

type forceEmpty struct {
	flag bool
}

func (o forceEmpty) ApplyToServiceHandlerOptions(opts *Options) {
	opts.forceEmpty = o.flag
}

func ForceEmpty(b bool) Option {
	return forceEmpty{b}
}

////////////////////////////////////////////////////////////////////////////////

type compoption = comphdlr.Option

type compoptionwrapper struct {
	compoption

	constraints []*semver.Constraints
	latestOnly  *bool
	resolver    resolvers.ComponentVersionResolver
}

func (o *compoptionwrapper) ApplyToServiceHandlerOptions(opts *Options) {
	if o.constraints != nil {
		opts.constraints = o.constraints
	}
	if o.latestOnly != nil {
		opts.latestOnly = o.latestOnly
	}
	if o.resolver != nil {
		opts.resolver = o.resolver
	}
}

func WithVersionConstraints(c []*semver.Constraints) Option {
	return &compoptionwrapper{compoption: comphdlr.WithVersionConstraints(c), constraints: c}
}

func LatestOnly(b ...bool) Option {
	return &compoptionwrapper{compoption: comphdlr.LatestOnly(b...), latestOnly: generics.Pointer(general.OptionalDefaultedBool(true, b...))}
}

func Resolver(r ocm.ComponentVersionResolver) Option {
	return &compoptionwrapper{compoption: comphdlr.Resolver(r), resolver: r}
}

////////////////////////////////////////////////////////////////////////////////

func MapToCompHandlerOptions(opts ...Option) comphdlr.Options {
	var copts []comphdlr.Option
	for _, o := range opts {
		if c, ok := o.(comphdlr.Option); ok {
			copts = append(copts, c)
		} else {
			if c, ok := o.(OptionList); ok {
				copts = append(copts, MapToCompHandlerOptions(c...))
			}
		}
	}
	return copts
}
