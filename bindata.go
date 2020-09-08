// Code generated by go-bindata.
// sources:
// udfs/postgres.sql
// udfs/snowflake.sql
// DO NOT EDIT!

package difference_digest

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

var _udfsPostgresSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\xff\x6f\x93\x40\x18\xc6\x7f\xbf\xbf\xe2\xc9\xac\x06\x26\x31\x40\x2b\xec\xb6\xb0\x04\xeb\xad\x21\x61\xb4\xf2\x25\x99\x59\x66\xd3\xb5\xd7\xf6\x62\x7b\x6d\x28\x28\xfa\xd7\x9b\x03\x4a\xa7\x66\x46\x13\xf9\x85\xbb\x97\x87\xf7\x7d\x9e\xcf\xdd\x30\x66\x7e\xca\xc0\xee\x52\x16\x25\xc1\x38\x42\x70\x83\x68\x9c\x82\xdd\x05\x49\x9a\x60\xbf\xd9\x7f\x2b\xd6\x3b\x59\x5e\x11\xd2\x4a\xc7\x31\x62\x36\x09\xfd\x21\xc3\x4d\x16\x0d\x53\xf5\xd3\x7e\x35\x2d\xf8\x76\xff\x66\x39\x5d\xcf\x0e\x6b\x4d\x2c\x2a\x08\x59\x18\xe0\x1b\xbe\xe5\xb2\xc0\xa3\x58\x09\x59\xe8\x88\x59\x9a\xc5\x51\xd2\xee\x09\x80\xae\x14\x65\x61\x88\x71\xd4\xbc\x83\x68\x92\xa5\xf5\xe7\xe0\xf6\x36\x4b\xfd\x77\x21\xab\x77\xa1\x1f\x8d\x32\x7f\xc4\x4e\xbe\x88\x9f\xa0\xd7\x23\x80\x1a\x9c\x70\xbe\x38\xc0\xc3\xbd\x75\x61\x3b\x94\x5e\x50\xea\x1a\xa0\xd6\xa0\x4f\x4d\xab\x4f\x0d\xd8\xa6\x6b\xbb\xd4\x72\xa8\x81\x81\xe5\x52\xd3\xec\xbb\x94\x1a\xb0\x9c\xc1\x05\x75\xfa\xa6\x6d\x3d\x10\x60\x3b\x3b\x7c\x86\xa7\xfc\x6b\x67\x66\xf5\xf6\xc9\x73\xa6\xa4\x3a\x21\x80\x90\x92\xe7\xf0\xa0\xd5\xe2\x73\x68\xc7\xa0\x9f\x4e\xcb\xeb\x6b\xf4\x6d\x5d\xd7\xf1\x12\x03\x9b\x0e\xa8\xe3\xda\xd4\x21\x40\xce\x8b\x32\x97\xd0\x3a\xc3\xf7\x62\x51\x3d\xe0\xbc\x69\xfa\x8b\xbc\xd7\xbb\xc2\x5f\x92\x2f\xf2\x99\xd8\x08\xb9\x9a\x7e\xe7\xf9\xee\xd0\xd9\x90\xe5\x96\xe7\x62\x7e\x62\xff\xbf\xc1\x57\x2d\xac\x76\xa0\x4e\x80\xf9\xae\x94\x05\x3c\x98\x0a\xd5\x52\x29\x3c\x98\x97\x75\xa3\x36\x7d\xad\x50\x24\xbf\xae\xc5\x86\x43\xd3\x2a\xbc\x82\xa5\xd7\x42\xfd\x12\xb5\xb4\x69\x5d\x29\x90\x56\x5b\x68\x1a\xbf\xf6\xba\x02\xf9\xb9\x23\xfe\x05\xd8\xa3\x28\xa6\xd5\x2e\x9f\x1e\x96\xa5\x9c\x6b\xb3\xd5\xaa\xbd\x96\x06\xbe\xcc\x36\x25\x7f\xfe\xd2\xfe\x91\xdc\x53\x6e\x1d\xb5\xe4\x43\xd8\xf2\x4a\x58\xc8\x86\x29\xd4\xb4\x17\xcd\x1c\x65\xb9\x73\xec\x8f\x46\x31\x1b\xa9\xd5\x6f\x3e\xa1\x1d\x0d\x69\x75\xf8\x44\x45\x82\xf7\x5c\x20\xa3\x11\xa5\x1f\x27\x0c\xde\x31\x59\x73\xb2\x51\x90\x0e\xc7\xd1\xfb\xfa\x80\xf4\x2b\xf2\x23\x00\x00\xff\xff\xce\x62\xde\x36\x03\x04\x00\x00")

func udfsPostgresSqlBytes() ([]byte, error) {
	return bindataRead(
		_udfsPostgresSql,
		"udfs/postgres.sql",
	)
}

func udfsPostgresSql() (*asset, error) {
	bytes, err := udfsPostgresSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "udfs/postgres.sql", size: 1027, mode: os.FileMode(420), modTime: time.Unix(1599539025, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _udfsSnowflakeSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x91\x5f\x6b\xd4\x40\x14\xc5\x9f\x9d\x4f\x71\x58\xaa\x4d\xda\x3e\xcc\x9f\x98\xe4\xb2\xec\xc2\x34\x3b\x9b\x8c\xc4\x49\x99\xcc\xaa\xf8\xb2\x0f\x36\xe2\x42\x5d\xa1\x6c\xb1\x20\x7e\x77\xc9\xbf\x8a\xb4\xe0\xdb\xbd\x87\x5f\x4e\xe6\x9c\x5b\x78\xa3\x83\x41\xe3\xe1\xcd\x4d\xad\x0b\x83\xed\xce\x15\xc1\x36\x0e\xdb\xfd\x66\xb3\xaf\x74\x5b\x45\x8b\xc3\xed\xe3\x02\xd6\x85\x2b\x2c\xba\xbb\xee\x7b\x77\x3c\x2d\x70\x6d\x4b\xeb\x42\xcc\x5e\x79\x13\x76\xde\xb5\x93\xc0\x00\xdd\x82\x9d\x9d\x31\x00\x88\x22\x52\xb9\x4c\x88\xa4\x54\x69\x9e\xbf\xc5\x05\xae\x6d\xf8\xd4\xf8\x68\xf2\xb9\xea\xf7\xb6\xb2\xdb\xe0\x6d\x59\x85\xbf\xb2\x92\x71\x8c\xd7\x48\x24\x25\x94\x66\x92\xd2\x18\x17\x18\x3c\x0b\xdd\x9a\x61\x00\x3e\x56\xc6\xe1\x70\xfb\x88\x15\x38\x42\xbf\x88\x5c\xa6\x44\x39\x51\xf6\x1c\x11\x23\x42\x22\x51\xc4\x85\xa2\xe7\x84\x1c\x09\xc9\x33\x99\x91\x48\x5f\x20\xd4\x48\x24\x22\x23\xce\x55\x46\x2f\x20\xc9\xf4\x92\x34\xc9\x29\x55\x5c\x8a\x01\x31\x6e\xf3\x6f\x9e\xbe\xa2\xe5\x92\xb1\xff\x5d\x20\x78\x6d\x6b\xeb\xca\xfd\x67\xe3\x9b\x76\x2e\x08\xdb\xba\xd1\x21\x66\xc0\x5c\xff\x20\x30\xa0\xd6\xae\xdc\xe9\xd2\xe0\x9d\xfe\xa0\xdb\xc2\xdb\x9b\xf9\x26\xe7\x0c\x38\x7c\x8d\x4c\x6d\xde\x1b\x17\xb0\x5a\x81\xc7\xb8\xef\x4e\x0f\xf7\x47\x70\x06\x7c\xf9\xf1\x70\x3c\xf5\x55\x32\xe0\xe7\xb7\xc3\x5d\x17\x3d\xb1\x6f\x20\xe2\xe9\x83\x5f\x63\x9c\xd9\xe4\x69\x5a\xaf\xd7\x18\xa3\x8e\x3e\x97\x97\x0c\xf8\xcd\x30\xff\x61\x50\xd9\xf9\xf2\x4f\x00\x00\x00\xff\xff\x7a\x2d\x3c\x12\x72\x02\x00\x00")

func udfsSnowflakeSqlBytes() ([]byte, error) {
	return bindataRead(
		_udfsSnowflakeSql,
		"udfs/snowflake.sql",
	)
}

func udfsSnowflakeSql() (*asset, error) {
	bytes, err := udfsSnowflakeSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "udfs/snowflake.sql", size: 626, mode: os.FileMode(420), modTime: time.Unix(1599539542, 0)}
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
	"udfs/postgres.sql": udfsPostgresSql,
	"udfs/snowflake.sql": udfsSnowflakeSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
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
	"udfs": &bintree{nil, map[string]*bintree{
		"postgres.sql": &bintree{udfsPostgresSql, map[string]*bintree{}},
		"snowflake.sql": &bintree{udfsSnowflakeSql, map[string]*bintree{}},
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

