package filedb

import (
	"sync"

	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/open-component-model/service-model/api/identity"
	"github.com/open-component-model/service-model/api/modeldesc"
	"ocm.software/ocm/api/utils"
	"ocm.software/ocm/api/utils/runtime"

	"github.com/mandelsoft/goutils/errors"
)

type FileDB struct {
	lock       sync.Mutex
	path       string
	filesystem vfs.FileSystem
	database   *modeldesc.ServiceModelDescriptor
}

var _ modeldesc.VersionResolver = (*FileDB)(nil)

func New(path string, fss ...vfs.FileSystem) *FileDB {
	return &FileDB{path: path, filesystem: utils.FileSystem(fss...)}
}

func (o *FileDB) Load() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.path != "" {
		data, err := vfs.ReadFile(utils.FileSystem(o.filesystem), o.path)
		if err != nil {
			if !vfs.IsErrNotExist(err) {
				return errors.Wrapf(err, "database file")
			}
			o.database = &modeldesc.ServiceModelDescriptor{
				DocType: runtime.NewVersionedObjectType(modeldesc.ABS_TYPE),
			}
		} else {
			o.database, err = modeldesc.Decode(data)
		}
	}
	return nil
}

func (o *FileDB) Save() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	if o.path != "" {
		data, err := modeldesc.Encode(o.database, runtime.DefaultYAMLEncoding)
		if err != nil {
			return err
		}
		return vfs.WriteFile(o.filesystem, o.path, data, 0o664)
	}
	return nil
}

func (o *FileDB) Add(s *modeldesc.ServiceDescriptor) {
	o.lock.Lock()
	defer o.lock.Unlock()

	key := s.GetId()
	if o.database == nil {
		return
	}

	for i, e := range o.database.Services {
		if key.Equals(e.GetId()) {
			o.database.Services[i] = *s
			return
		}
	}
	o.database.Services = append(o.database.Services, *s)
}

func (o *FileDB) LookupServiceVersionVariant(id identity.ServiceVersionVariantIdentity) (*modeldesc.ServiceDescriptor, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	for _, s := range o.database.Services {
		if id.Equals(s.GetId()) {
			return s.Copy(), nil
		}
	}
	return nil, errors.ErrNotFound(modeldesc.KIND_SERVICEVERSION, id.String())
}

func (o *FileDB) ListVersions(id identity.ServiceIdentity, variant ...identity.Variant) ([]string, error) {
	o.lock.Lock()
	defer o.lock.Unlock()

	var vari = utils.Optional(variant...)
	var result []string
	for _, s := range o.database.Services {
		if !vari.Equals(s.GetVariant()) {
			continue
		}
		if id != s.Service {
			continue
		}
		result = append(result, s.Version)
	}
	return result, nil
}
