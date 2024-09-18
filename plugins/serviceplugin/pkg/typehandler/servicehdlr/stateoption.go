package servicehdlr

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mandelsoft/goutils/maputils"
	"github.com/mandelsoft/goutils/set"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/open-component-model/service-model/api/crossref"
	"github.com/open-component-model/service-model/api/filedb"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/open-component-model/service-model/api/utils"
	"github.com/spf13/pflag"
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/api/datacontext/attrs/vfsattr"
	"ocm.software/ocm/cmds/ocm/common/options"
)

var defaultRelations = []crossref.DepKind{crossref.DEP_DEPENDS, crossref.DEP_INSTALLEDBY}

var relOpts = map[string][]crossref.DepKind{
	string(crossref.DEP_DEPENDS):     {crossref.DEP_DEPENDS},
	string(crossref.DEP_INSTALLEDBY): {crossref.DEP_INSTALLEDBY},
	string(crossref.DEP_MANAGES):     {crossref.DEP_MANAGES},
	string(crossref.DEP_INSTANCE):    {crossref.DEP_INSTANCE},
	string(crossref.DEP_SATISFIES):   {crossref.DEP_SATISFIES},

	"all":      {crossref.DEP_DEPENDS, crossref.DEP_INSTALLEDBY, crossref.DEP_MANAGES, crossref.DEP_INSTANCE, crossref.DEP_SATISFIES},
	"services": {crossref.DEP_DEPENDS, crossref.DEP_INSTALLEDBY, crossref.DEP_MANAGES, crossref.DEP_INSTANCE},
}

func From(o options.OptionSetProvider) *State {
	var opt *State
	o.AsOptionSet().Get(&opt)
	return opt
}

func NewState(r modeldesc.VersionResolver) *State {
	return &State{Resolver: r, Relations: set.New(defaultRelations...)}
}

type State struct {
	lock sync.Mutex

	relations []string

	Resolver        modeldesc.VersionResolver
	ResolveToLatest bool

	UpdatePath   string
	DatabasePath string

	Relations set.Set[crossref.DepKind]

	Filesystem vfs.FileSystem
	Database   *filedb.FileDB
}

var _ Option = (*State)(nil)

func (o *State) WithLatestResolution() *State {
	o.ResolveToLatest = true
	return o
}

func (o *State) WithRelations(rels ...crossref.DepKind) *State {
	o.Relations = set.New[crossref.DepKind]()
	for _, r := range rels {
		o.Relations.Add(relOpts[string(r)]...)
	}
	return o
}

func (o *State) ApplyToServiceHandlerOptions(opts *Options) {
	opts.state = o
}

func (o *State) AddFlags(fs *pflag.FlagSet) {
	possible := strings.Join(maputils.OrderedKeys(relOpts), ",")
	fs.StringVarP(&o.UpdatePath, "update", "U", "", "update service database file")
	fs.StringVarP(&o.DatabasePath, "database", "D", "", "examine database file")
	fs.BoolVarP(&o.ResolveToLatest, "constraintsToLatest", "R", false, "resolve version constraints to latest")
	fs.StringSliceVarP(&o.relations, "relations", "", utils.Convert[string](defaultRelations), "relations to follow ("+possible+")")
}

func (o *State) Configure(ctx cli.Context) error {
	if o.Filesystem == nil {
		o.Filesystem = vfsattr.Get(ctx)
	}

	if o.UpdatePath != "" {
		o.Database = filedb.New(o.UpdatePath, o.Filesystem)
		return o.Database.Load()
	}

	if len(o.relations) > 0 {
		o.Relations = set.New[crossref.DepKind]()
		for _, r := range o.relations {
			f := relOpts[r]
			if len(f) == 0 {
				return fmt.Errorf("invalid relation %q", r)
			}
			o.Relations.Add(f...)
		}
	}
	return nil
}

func (o *State) IsStandardRelations() bool {
	if len(o.Relations) != len(defaultRelations) {
		return false
	}
	for _, r := range defaultRelations {
		if !o.Relations.Has(r) {
			return false
		}
	}
	return true
}

func (o *State) Save() error {
	if o.Database != nil {
		return o.Database.Save()
	}
	return nil
}

func (o *State) Add(s *modeldesc.ServiceDescriptor) {
	if o.Database == nil {
		return
	}
	o.Database.Add(s)
}
