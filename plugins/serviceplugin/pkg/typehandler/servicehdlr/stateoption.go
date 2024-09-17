package servicehdlr

import (
	"sync"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/open-component-model/service-model/api/filedb"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/spf13/pflag"
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/api/datacontext/attrs/vfsattr"
	"ocm.software/ocm/cmds/ocm/common/options"
)

func From(o options.OptionSetProvider) *State {
	var opt *State
	o.AsOptionSet().Get(&opt)
	return opt
}

func NewState(r modeldesc.Resolver) *State {
	return &State{Resolver: r}
}

type State struct {
	lock sync.Mutex

	Resolver modeldesc.Resolver

	UpdatePath   string
	DatabasePath string
	Filesystem   vfs.FileSystem
	Database     *filedb.FileDB
}

var _ Option = (*State)(nil)

func (o *State) ApplyToServiceHandlerOptions(opts *Options) {
	opts.state = o
}

func (o *State) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.UpdatePath, "update", "U", "", "update service database file")
	fs.StringVarP(&o.DatabasePath, "database", "D", "", "examine database file")
}

func (o *State) Configure(ctx cli.Context) error {
	if o.Filesystem == nil {
		o.Filesystem = vfsattr.Get(ctx)
	}

	if o.UpdatePath != "" {
		o.Database = filedb.New(o.UpdatePath, o.Filesystem)
		return o.Database.Load()
	}
	return nil
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
