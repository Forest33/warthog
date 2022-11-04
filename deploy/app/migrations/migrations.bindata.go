// Code generated by go-bindata. DO NOT EDIT.
// sources:
// migrations/1633685677_init.down.sql (32B)
// migrations/1633685677_init.up.sql (1.207kB)

package migrations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
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

	info := bindataFileInfo{name: "migrations/1633685677_init.down.sql", size: 32, mode: os.FileMode(0664), modTime: time.Unix(1665214714, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xe8, 0x16, 0x49, 0x52, 0xa1, 0x76, 0xbe, 0xd5, 0x80, 0x60, 0xe1, 0x4e, 0x86, 0x4b, 0xe4, 0xc, 0x6c, 0x14, 0x63, 0x58, 0x6a, 0x87, 0x3d, 0xd3, 0x42, 0xf3, 0xb0, 0xfe, 0x20, 0x2, 0x6f, 0x87}}
	return a, nil
}

var _migrations1633685677_initUpSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x53\x41\x6f\xb2\x40\x10\xbd\xf3\x2b\xe6\x06\x24\x1e\xfc\xbe\xd4\xb6\x89\xa7\x55\xd7\x96\x88\x60\x60\x6d\xea\x89\x6c\xd8\x55\x36\xa5\x40\x60\x0d\xfa\xef\x9b\x15\x41\x68\x1b\xda\xc6\x76\x0e\xc6\x19\xe6\xbd\x61\xde\x1b\xa6\x1e\x46\x04\x03\x41\x13\x1b\xc3\x6e\x2f\x82\x30\x4d\xb6\x62\xa7\x19\x1a\x00\x40\x29\x12\x96\x96\x41\x29\x98\x8c\x00\x2c\x87\x80\xe3\x12\x70\xd6\xb6\x0d\x55\xcc\xf0\x1c\xad\x6d\x02\xff\x86\xff\x6f\x06\x6d\x48\xc4\xc5\x2e\x92\x7d\x90\xbb\xdb\xfb\x0e\xe2\x50\x3d\xee\x41\x8c\x86\x1d\xc0\xf1\xdb\x80\x30\xe7\x54\x72\x16\x50\xa9\x32\x62\x2d\xb1\x4f\xd0\x72\x05\x33\x44\xb0\xca\x1a\x80\xc1\xa8\xe4\x52\xbc\x72\x43\x4f\xd2\x52\x1f\x80\x1e\xa7\x21\x8d\x55\x45\x37\xcd\x66\x4c\xc5\xba\xcf\xd8\xaf\xb2\x6a\xe6\x58\xd3\x2c\xc7\xc7\x1e\x51\x4b\xb9\x2d\x3b\xc0\x68\x3b\x31\xe8\x8a\xdc\xa4\x87\xe6\xdf\xd1\xd4\x9e\x90\xbd\xc6\x3e\x18\x27\x67\x4e\x62\x2b\x39\x60\x34\x54\x53\x3a\xae\x97\x69\xfe\x52\x64\x34\xe4\x67\xd3\x05\x83\x3a\x2c\x87\xe0\x07\xec\xc1\xca\xb3\x96\xc8\xdb\xc0\x02\x6f\xaa\xdd\x33\x9a\xf3\x44\x06\xaa\xb5\xee\xe9\x8d\x8b\x6a\x11\x2d\x82\x30\x12\x31\x03\x98\xb8\xae\xdd\x0f\x53\xc8\xda\xda\x5a\xcd\x39\xb2\x7d\x5c\x71\xc9\x63\xc6\xcf\x6d\x04\x3f\x13\x98\x3e\xe2\xe9\x02\x8c\x53\xd9\x72\xc0\xd0\xb7\x4a\xec\x42\xfd\xe4\x1f\xfd\x93\x42\xc6\xfc\x82\xae\x8a\x8c\x4a\xda\xa6\xfc\xf2\xf5\x1a\xba\x22\xcd\xe5\x3b\xd9\x7e\xb6\xd8\xf9\x58\xf9\x21\xa3\x09\xe3\xca\x84\xab\x05\x6a\x1d\xfe\x27\xf7\xd9\x13\xd7\x7e\x10\x7f\x3a\x4d\x33\xc7\x6f\x01\x00\x00\xff\xff\x1a\x33\x79\x8f\xb7\x04\x00\x00")

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

	info := bindataFileInfo{name: "migrations/1633685677_init.up.sql", size: 1207, mode: os.FileMode(0664), modTime: time.Unix(1667545360, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x7d, 0x37, 0x20, 0xb, 0x9d, 0xc0, 0xf3, 0x84, 0x2a, 0x13, 0x6c, 0x57, 0x36, 0x68, 0x35, 0xa4, 0xeb, 0x6, 0x34, 0x50, 0x1a, 0x81, 0x36, 0xd2, 0x1d, 0x7d, 0x56, 0xfe, 0x42, 0xd6, 0xbd, 0x6}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"migrations/1633685677_init.down.sql": migrations1633685677_initDownSql,
	"migrations/1633685677_init.up.sql":   migrations1633685677_initUpSql,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"migrations": {nil, map[string]*bintree{
		"1633685677_init.down.sql": {migrations1633685677_initDownSql, map[string]*bintree{}},
		"1633685677_init.up.sql":   {migrations1633685677_initUpSql, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}