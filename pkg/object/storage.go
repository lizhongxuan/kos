package object

import (
	"context"
	"errors"
	"fmt"
	"io"
	"kos/pkg/utils"
	"os"
	"sync"
	"time"
)

var ctx = context.Background()

var UserAgent = "KOS"

type MtimeChanger interface {
	Chtimes(path string, mtime time.Time) error
}

type SupportSymlink interface {
	// Symlink create a symbolic link
	Symlink(oldName, newName string) error
	// Readlink read a symbolic link
	Readlink(name string) (string, error)
}

type File interface {
	Object
	Owner() string
	Group() string
	Mode() os.FileMode
}

type onlyWriter struct {
	io.Writer
}

type file struct {
	obj
	owner     string
	group     string
	mode      os.FileMode
	isSymlink bool
}

func (f *file) Owner() string     { return f.owner }
func (f *file) Group() string     { return f.group }
func (f *file) Mode() os.FileMode { return f.mode }
func (f *file) IsSymlink() bool   { return f.isSymlink }

func MarshalObject(o Object) map[string]interface{} {
	m := make(map[string]interface{})
	m["key"] = o.Key()
	m["size"] = o.Size()
	m["mtime"] = o.Mtime().UnixNano()
	m["isdir"] = o.IsDir()
	if f, ok := o.(File); ok {
		m["mode"] = f.Mode()
		m["owner"] = f.Owner()
		m["group"] = f.Group()
		m["isSymlink"] = f.IsSymlink()
	}
	return m
}

func UnmarshalObject(m map[string]interface{}) Object {
	mtime := time.Unix(0, int64(m["mtime"].(float64)))
	o := obj{
		key:   m["key"].(string),
		size:  int64(m["size"].(float64)),
		mtime: mtime,
		isDir: m["isdir"].(bool)}
	if _, ok := m["mode"]; ok {
		f := file{o, m["owner"].(string), m["group"].(string), os.FileMode(m["mode"].(float64)), m["isSymlink"].(bool)}
		return &f
	}
	return &o
}

type FileSystem interface {
	MtimeChanger
	Chmod(path string, mode os.FileMode) error
	Chown(path string, owner, group string) error
}

var notSupported = utils.ENOTSUP
var notSupportedDelimiter = errors.New("not supported delimiter")

type DefaultObjectStorage struct{}

var _ ObjectStorage = &DefaultObjectStorage{}

func (s DefaultObjectStorage) Create() error {
	return notSupported
}
func (s *DefaultObjectStorage) Delete(key string) error {
	return notSupported
}

func (s *DefaultObjectStorage) String() string {
	return ""
}

func (s DefaultObjectStorage) Head(key string) (Object, error) {
	return nil, notSupported
}

func (s DefaultObjectStorage) CreateMultipartUpload(key string) (*MultipartUpload, error) {
	return nil, notSupported
}

func (s DefaultObjectStorage) UploadPart(key string, uploadID string, num int, body []byte) (*Part, error) {
	return nil, notSupported
}

func (s DefaultObjectStorage) AbortUpload(key string, uploadID string) {}

func (s DefaultObjectStorage) CompleteUpload(key string, uploadID string, parts []*Part) error {
	return notSupported
}

func (s DefaultObjectStorage) ListUploads(marker string) ([]*PendingPart, string, error) {
	return nil, "", nil
}

func (s DefaultObjectStorage) List(prefix, marker, delimiter string, limit int64) ([]Object, error) {
	return nil, notSupported
}

func (s DefaultObjectStorage) ListAll(prefix, marker string) (<-chan Object, error) {
	return nil, notSupported
}

func (s *DefaultObjectStorage) Get(key string, off, limit int64) (io.ReadCloser, error) {
	return nil, notSupported
}

func (u *DefaultObjectStorage) Put(key string, body io.Reader) error {
	return notSupported
}

type Creator func(bucket, accessKey, secretKey, token string) (ObjectStorage, error)

var storages = make(map[string]Creator)

func Register(name string, register Creator) {
	storages[name] = register
}

func CreateStorage(name, endpoint, accessKey, secretKey, token string) (ObjectStorage, error) {
	f, ok := storages[name]
	if ok {
		logger.Debugf("Creating %s storage at endpoint %s", name, endpoint)
		return f(endpoint, accessKey, secretKey, token)
	}
	return nil, fmt.Errorf("invalid storage: %s", name)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 32<<10)
		return &buf
	},
}
