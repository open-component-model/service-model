package services

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/sliceutils"
	ocmdesc "github.com/open-component-model/service-model/api/ocm"
	"github.com/open-component-model/service-model/plugins/serviceplugin/pkg/typehandler"
	"github.com/open-component-model/service-model/plugins/serviceplugin/pkg/typehdlrutils"
	"github.com/spf13/pflag"
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/api/ocm/resolvers"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/handlers/comphdlr"
	"ocm.software/ocm/cmds/ocm/common/processing"
	"ocm.software/ocm/cmds/ocm/common/utils"

	// bind OCM configuration.
	_ "ocm.software/ocm/api/ocm/plugin/ppi/config"
	"ocm.software/ocm/cmds/ocm/commands/common/options/closureoption"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/options/lookupoption"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/options/repooption"
	"ocm.software/ocm/cmds/ocm/commands/ocmcmds/common/options/versionconstraintsoption"
	"ocm.software/ocm/cmds/ocm/common/options"
	"ocm.software/ocm/cmds/ocm/common/output"

	"github.com/mandelsoft/logging"
	"github.com/spf13/cobra"

	"ocm.software/ocm/api/ocm"
)

const Name = "services"

var log = logging.DynamicLogger(logging.DefaultContext(), logging.NewRealm("cliplugin/service-model"))

func New() *cobra.Command {
	cmd := &command{
		OptionSet: options.OptionSet{
			versionconstraintsoption.New(), repooption.New(),
			output.OutputOptions(outputs,
				closureoption.New("service", options.Not(output.Selected("tree"))),
				lookupoption.New(),
			),
		},
	}
	c := &cobra.Command{
		Use:   Name + " <options> <elems>",
		Short: "get service definitions",
		Long: `Resolve services given by ids or show services defined by component versions
together with the option -C.
`,
		RunE: cmd.Run,
	}

	cmd.AddFlags(c.Flags())
	return c
}

type command struct {
	options.OptionSet
	useComps bool
}

func (c *command) AddFlags(set *pflag.FlagSet) {
	c.OptionSet.AddFlags(set)
	set.BoolVarP(&c.useComps, "components", "C", false, "use component versions instead of service ids")
}

func (c *command) Run(cmd *cobra.Command, args []string) error {
	ctx := ocm.FromContext(cmd.Context())

	clictx := cli.WithOCM(ctx).WithOutput(cmd.OutOrStdout()).WithErrorOutput(cmd.ErrOrStderr()).New()
	session := ocm.NewSession(nil)

	err := c.ProcessOnOptions(common.CompleteOptionsWithSession(clictx, session))
	if err != nil {
		return err
	}
	oopts := output.From(c)
	oopts.Context = clictx

	var mainargs []string
	var h utils.TypeHandler

	var resolver resolvers.ComponentResolver

	repo := repooption.From(c).Repository
	if repo != nil {
		resolver = resolvers.ComponentResolverForRepository(repo)
	} else {
		r := ctx.GetResolver()
		if r == nil {
			return fmt.Errorf("no component resolver configured")
		}
		resolver = r.(resolvers.ComponentResolver)
	}

	if c.useComps {
		var comps []string
		for _, a := range args {
			if strings.Index(a, "/") >= 0 {
				comps = append(comps, a)
			} else {
				mainargs = append(mainargs, a)
			}
		}
		h, err = typehandler.ForComponents(NewOCM(ctx), resolver, output.From(c), repo, session, comps, typehandler.OptionsFor(c))
		if err != nil {
			return err
		}
	} else {
		mainargs = args
		h = typehandler.ForServices(resolver, typehandler.OptionsFor(c))
	}

	oopts.OptionSet = append(oopts.OptionSet, typehandler.NewState(ocmdesc.NewServiceResolver(resolvers.ComponentVersionResolverForComponentResolver(resolver))))

	return utils.HandleArgs(oopts, h, mainargs...)
}

func TableOutput(opts *output.Options, mapping processing.MappingFunction, wide ...string) *output.TableOutput {
	return &output.TableOutput{
		Headers: output.Fields("COMPONENT", "NAME", "VERSION", "VARIANT", "KIND", "SHORTNAME", wide),
		Options: opts,
		Chain:   typehandler.Sort,
		Mapping: mapping,
	}
}

var outputs = output.NewOutputs(getRegular, output.Outputs{
	"tree": getTree,
}).AddChainedManifestOutputs(output.ComposeChain(closureoption.OutputChainFunction(typehandler.ClosureExplode, comphdlr.Sort.Transform(typehdlrutils.NormalizeFunction))))

func getRegular(opts *output.Options) output.Output {
	return NormalizedTableOutput(closureoption.TableOutput(TableOutput(opts, mapGetRegularOutput), typehandler.ClosureExplode), typehdlrutils.NormalizeFunction).New()
}

func getTree(opts *output.Options) output.Output {
	return output.TreeOutput(NormalizedTableOutput(closureoption.TableOutput(TableOutput(opts, mapGetRegularOutput), typehandler.ClosureExplode), typehdlrutils.NormalizeFunction), "NESTING").New()
}

func NormalizedTableOutput(in *output.TableOutput, norm ...typehandler.NormalizeFunction) *output.TableOutput {
	f := general.Optional(norm...)
	out := *in
	out.Chain = processing.Append(in.Chain, processing.Transform(f.Normalizer(in.Options)))
	return &out
}

func mapGetRegularOutput(e interface{}) interface{} {
	obj := e.(*typehandler.Object)
	if obj.Node == nil {
		return sliceutils.AsSlice("...", "", "", "", "", "")
	}
	r := obj.Element
	if r == nil {
		return sliceutils.AsSlice(obj.Id.Component, obj.Id.Name, obj.Id.Version, obj.Id.Variant.String(), "", "")
	}
	return sliceutils.AsSlice(r.Service.Component, r.Service.Name, r.Version, r.Kind.GetVariant().String(), r.Kind.GetType(), r.ShortName)
}
