package servicehdlr

import (
	"sync"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/open-component-model/service-model/api/modeldesc"
	"github.com/spf13/pflag"
	"ocm.software/ocm/api/cli"
	"ocm.software/ocm/api/datacontext/attrs/vfsattr"
	"ocm.software/ocm/api/utils"
	"ocm.software/ocm/api/utils/runtime"
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

	Path       string
	Filesystem vfs.FileSystem
	Database   *modeldesc.ServiceModelDescriptor
}

var _ Option = (*State)(nil)

func (o *State) ApplyToServiceHandlerOptions(opts *Options) {
	opts.state = o
}

func (o *State) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Path, "database", "D", "", "update service database file")
}

func (o *State) Configure(ctx cli.Context) error {
	o.Filesystem = vfsattr.Get(ctx)
	return nil
}

func (o *State) Load() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.Path != "" {
		data, err := vfs.ReadFile(utils.FileSystem(o.Filesystem), o.Path)
		if err != nil {
			if !vfs.IsErrNotExist(err) {
				return errors.Wrapf(err, "database file")
			}
			o.Database = &modeldesc.ServiceModelDescriptor{
				DocType: runtime.NewVersionedObjectType(modeldesc.ABS_TYPE),
			}
		} else {
			o.Database, err = modeldesc.Decode(data)
		}
	}
	return nil
}

func (o *State) Save() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.Path != "" {
		data, err := modeldesc.Encode(o.Database, runtime.DefaultYAMLEncoding)
		if err != nil {
			return err
		}
		return vfs.WriteFile(o.Filesystem, o.Path, data, 0o664)
	}
	return nil
}

func (o *State) Add(s *modeldesc.ServiceDescriptor) {
	o.lock.Lock()
	defer o.lock.Unlock()

	key := s.GetId()
	if o.Database == nil {
		return
	}

	for i, e := range o.Database.Services {
		if key.Equals(e.GetId()) {
			o.Database.Services[i] = *s
			return
		}
	}
	o.Database.Services = append(o.Database.Services, *s)
}
