// Code generated for package migrations by go-bindata DO NOT EDIT. (@generated)
// sources:
// migrations/1633685677_init.down.sql
// migrations/1633685677_init.up.sql
// migrations/1668845636_settings.down.sql
// migrations/1668845636_settings.up.sql
// migrations/1677165148_k8s.down.sql
// migrations/1677165148_k8s.up.sql
package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _migrations1633685677_initDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x48\x2f\xcd\x8c\x4f\xce\xcf\x4b\xcb\x4c\xb7\x06\x04\x00\x00\xff\xff\x49\xa7\x32\xcb\x20\x00\x00\x00")

func migrations1633685677_initDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations1633685677_initDownSql,
		"migrations/1633685677_init.down.sql",
	)
}

func migrations1633685677_initDownSql() (*asset, error) {
	bytes, err := migrations1633685677_initDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1633685677_init.down.sql", size: 32, mode: os.FileMode(436), modTime: time.Unix(1665214714, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations1633685677_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x93\x51\x6f\xb2\x30\x14\x86\xef\xf9\x15\xe7\x0e\x48\xbc\xf0\xfb\x32\xb7\x25\xbb\x42\xad\x5b\x23\xc2\x02\x75\xd1\x2b\xd2\xd0\x2a\xcd\x18\x10\xa8\x41\xff\xfd\x52\x50\x04\x67\x74\x8b\xdb\xb9\x20\xf4\xe4\x3d\x6f\xdb\xe7\x9c\x8e\x3c\x64\x11\x04\xc4\x1a\xda\x08\xf0\x04\x1c\x97\x00\x5a\x60\x9f\xf8\xb0\xde\x88\x20\x4c\x93\x95\x58\x6b\x86\x06\x00\x50\x8a\x84\xa5\x65\x50\x0a\x26\x23\x00\xec\x90\x4a\xed\xcc\x6d\x1b\xc6\x68\x62\xcd\x6d\x02\xff\xfa\xff\xef\x7a\x6d\x71\xc4\xc5\x3a\x92\xe7\xc5\x0f\xf7\x8f\x1d\xed\x16\xaa\x38\xab\x1d\xf4\x3b\xd2\xdd\x37\xa4\x61\xce\xa9\xe4\x2c\xa0\x52\xad\xc6\x16\x41\x04\xcf\x50\x55\x77\x90\x1a\x8c\x4a\x2e\xc5\x07\x37\xf4\x24\x2d\xf5\x1e\xe8\x71\x1a\xd2\x58\x65\x74\xd3\x6c\xac\x6b\xbf\x4d\xc6\x7e\xc9\x4f\x33\x9f\x34\x0d\x3b\x3e\xf2\x88\xba\x82\xdb\x42\x0d\x46\x9b\x72\xaf\x8b\xb1\x59\x6e\x9b\xbf\x9d\xa9\xbd\x59\xf6\x1c\xf9\x60\x54\xec\x2b\xa8\x0a\x01\x0c\xfa\x6a\x97\x0b\xfd\x2d\xd3\xfc\xbd\xc8\x68\xc8\xf7\xed\x15\x0c\x0e\x81\x1d\x82\x9e\x91\x07\xaf\x1e\x9e\x59\xde\x12\xa6\x68\x59\x33\xc8\x68\xce\x13\x19\x28\xe9\x41\x73\x31\x8e\xf4\x22\x5a\x04\x61\x24\x62\x06\x30\x74\x5d\xfb\x72\x99\xaa\x3c\x6d\xeb\xc4\xb2\x7d\x54\x7b\xc9\x5d\xc6\xf7\x32\x82\x16\x04\x46\x2f\x68\x34\x05\xa3\x4a\x63\x07\x0c\x7d\xa5\xd0\x17\xea\x93\x7f\xed\xa3\x14\x32\xe6\xc7\xea\x3a\xc9\xa8\xa4\x6d\xcb\xab\xc7\x6b\xec\x8a\x34\x97\x27\xd8\x7e\x76\xb1\xfd\xb8\xf2\x6d\x46\x13\xc6\x55\x13\x6e\x06\xd4\x1a\xfd\xce\x9c\x5e\x8d\x5b\x1f\xc6\x9f\xee\xa6\x9e\xcd\x67\x00\x00\x00\xff\xff\x4d\x4a\x45\x79\xb0\x04\x00\x00")

func migrations1633685677_initUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations1633685677_initUpSql,
		"migrations/1633685677_init.up.sql",
	)
}

func migrations1633685677_initUpSql() (*asset, error) {
	bytes, err := migrations1633685677_initUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1633685677_init.up.sql", size: 1200, mode: os.FileMode(436), modTime: time.Unix(1669023217, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations1668845636_settingsDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\xcf\x4d\xce\x82\x40\x0c\xc6\xf1\x3d\xa7\xe8\x3d\x58\xf1\xbe\xce\x8e\x0f\x43\x70\xdd\xc0\x50\x87\x46\xa6\x45\xa6\x24\x7a\x7b\x63\xe2\xd2\xc5\x78\x80\x5f\x9e\xe7\x5f\xd5\x83\xeb\x61\xa8\xfe\x6a\x07\x89\xcc\x58\x42\x2a\x00\x00\x4e\x7d\x77\x86\xff\xae\xbe\x34\x2d\x24\x96\xb0\x12\xb2\x24\x1b\xc5\x53\x59\x64\x29\xaf\x22\xe4\x0d\x8d\x23\xe9\x61\x99\x6a\xa7\xfb\x41\xe9\x57\x25\x2a\x38\xad\xea\x6f\x2c\x01\x3f\xc3\xac\x92\xa9\x93\xee\x86\x91\x6c\xd1\x39\xe1\xf4\x44\x19\x63\x6e\x64\x1c\x1f\xb8\xaa\x6e\x38\xd3\x66\x4b\x59\x7c\x55\xd0\xbb\xb6\x6a\x1c\x0c\x1d\x84\x83\xdf\xff\xae\x1c\xca\x57\x00\x00\x00\xff\xff\x8c\x3b\x19\xa0\x7b\x01\x00\x00")

func migrations1668845636_settingsDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations1668845636_settingsDownSql,
		"migrations/1668845636_settings.down.sql",
	)
}

func migrations1668845636_settingsDownSql() (*asset, error) {
	bytes, err := migrations1668845636_settingsDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1668845636_settings.down.sql", size: 379, mode: os.FileMode(436), modTime: time.Unix(1669023217, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations1668845636_settingsUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\xd0\xc1\x4a\x03\x31\x10\x06\xe0\xfb\x3e\xc5\x3c\x82\xe2\xb1\xa7\xd4\x8d\x20\xa4\x09\x2c\xd9\xf3\xb0\x4d\xc7\x34\xb8\x99\xa9\xcd\x2c\xe8\xdb\x0b\xa2\x20\x88\x74\x61\xaf\xc3\xcf\xc7\xfc\xbf\x71\xd1\x0e\x10\xcd\xde\x59\xc8\x4b\xc1\x24\xfc\x52\x72\x07\x00\x30\x58\x6f\x0e\x16\x62\x80\x46\xaa\x85\x73\xdb\x75\xdd\xef\xfc\xcf\xf9\x2b\x6d\xfa\x1e\x1e\x83\x1b\x0f\x1e\x5a\xe1\x3c\x13\x16\x6e\x3a\x71\x22\xd8\x87\xe0\xc0\x87\x08\x7e\x74\x0e\x7a\xfb\x64\x46\x17\x21\x0e\xa3\xdd\xad\xf2\x92\x30\x53\x52\xd4\x52\x49\x16\x85\x67\x1f\xff\x72\xf7\x77\xeb\xb0\x2b\xbd\x2d\xd4\x6e\x60\x0f\x2b\x31\x16\xc6\xe3\x2c\xe9\xb5\x70\xc6\xef\x37\x8b\xf0\xe6\xc6\x4d\xae\x8a\x95\xf4\x2c\xa7\x86\xc7\x0f\xe4\xa9\x6e\x9f\xb1\x4e\xef\x38\x8b\x5c\xf0\x44\x17\x3d\xff\xbb\xe2\x67\x00\x00\x00\xff\xff\x8b\x20\x22\x7f\x12\x02\x00\x00")

func migrations1668845636_settingsUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations1668845636_settingsUpSql,
		"migrations/1668845636_settings.up.sql",
	)
}

func migrations1668845636_settingsUpSql() (*asset, error) {
	bytes, err := migrations1668845636_settingsUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1668845636_settings.up.sql", size: 530, mode: os.FileMode(436), modTime: time.Unix(1669023217, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations1677165148_k8sDownSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\x28\x4e\x2d\x29\xc9\xcc\x4b\x2f\xe6\x52\x50\x50\x50\x70\x09\xf2\x0f\x50\x70\xf6\xf7\x09\xf5\xf5\x53\xc8\xb6\x28\x8e\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x89\x2f\xc9\xcc\x4d\xcd\x2f\x2d\xb1\xe6\x02\x04\x00\x00\xff\xff\x60\x26\x95\x49\x3a\x00\x00\x00")

func migrations1677165148_k8sDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations1677165148_k8sDownSql,
		"migrations/1677165148_k8s.down.sql",
	)
}

func migrations1677165148_k8sDownSql() (*asset, error) {
	bytes, err := migrations1677165148_k8sDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1677165148_k8s.down.sql", size: 58, mode: os.FileMode(436), modTime: time.Unix(1677165237, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrations1677165148_k8sUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x04\xc0\xc1\x0a\xc2\x30\x0c\x06\xe0\x7b\x9f\xe2\x7f\x04\xc1\x8b\xe0\x29\xda\x0a\x42\x4c\x41\xd2\x73\x4f\x41\x8a\x6c\x63\x4b\xfa\xfe\xfb\x88\xb5\x7c\xa1\xf4\xe0\x02\xb7\x88\xb1\xfe\x3c\x01\x00\xe5\x8c\x67\xe5\xf6\x11\xfc\x6f\xde\x0f\xdb\xa7\x79\xf4\x18\x8b\x6d\x33\xf0\x16\x85\x54\x85\x34\x66\xe4\xf2\xa2\xc6\x8a\xeb\xe5\x9e\xce\x00\x00\x00\xff\xff\x70\xdf\x52\x1c\x51\x00\x00\x00")

func migrations1677165148_k8sUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrations1677165148_k8sUpSql,
		"migrations/1677165148_k8s.up.sql",
	)
}

func migrations1677165148_k8sUpSql() (*asset, error) {
	bytes, err := migrations1677165148_k8sUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/1677165148_k8s.up.sql", size: 81, mode: os.FileMode(436), modTime: time.Unix(1677165237, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"migrations/1633685677_init.down.sql":     migrations1633685677_initDownSql,
	"migrations/1633685677_init.up.sql":       migrations1633685677_initUpSql,
	"migrations/1668845636_settings.down.sql": migrations1668845636_settingsDownSql,
	"migrations/1668845636_settings.up.sql":   migrations1668845636_settingsUpSql,
	"migrations/1677165148_k8s.down.sql":      migrations1677165148_k8sDownSql,
	"migrations/1677165148_k8s.up.sql":        migrations1677165148_k8sUpSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"migrations": &bintree{nil, map[string]*bintree{
		"1633685677_init.down.sql":     &bintree{migrations1633685677_initDownSql, map[string]*bintree{}},
		"1633685677_init.up.sql":       &bintree{migrations1633685677_initUpSql, map[string]*bintree{}},
		"1668845636_settings.down.sql": &bintree{migrations1668845636_settingsDownSql, map[string]*bintree{}},
		"1668845636_settings.up.sql":   &bintree{migrations1668845636_settingsUpSql, map[string]*bintree{}},
		"1677165148_k8s.down.sql":      &bintree{migrations1677165148_k8sDownSql, map[string]*bintree{}},
		"1677165148_k8s.up.sql":        &bintree{migrations1677165148_k8sUpSql, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
