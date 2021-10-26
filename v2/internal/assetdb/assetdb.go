package assetdb

import (
	"fmt"
	"strings"
	"unsafe"
)

// AssetDB is a database for assets encoded as byte slices
type AssetDB struct {
	db map[string][]byte
}

// NewAssetDB creates a new AssetDB and initialises a blank db
func NewAssetDB() *AssetDB {
	return &AssetDB{
		db: make(map[string][]byte),
	}
}

// AddAsset saves the given byte slice under the given name
func (a *AssetDB) AddAsset(name string, data []byte) {
	a.db[name] = data
}

// Remove removes the named asset
func (a *AssetDB) Remove(name string) {
	delete(a.db, name)
}

// Asset retrieves the byte slice for the given name
func (a *AssetDB) Read(name string) ([]byte, error) {
	result := a.db[name]
	if result == nil {
		return nil, fmt.Errorf("asset '%s' not found", name)
	}
	return result, nil
}

// AssetAsString returns the asset as a string.
// It also returns a boolean indicating whether the asset existed or not.
func (a *AssetDB) String(name string) (string, error) {
	asset, err := a.Read(name)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&asset)), nil
}

func (a *AssetDB) Dump() {
	fmt.Printf("Assets:\n")
	for k, _ := range a.db {
		fmt.Println(k)
	}
}

// Serialize converts the entire database to a file that when compiled will
// reconstruct the AssetDB during init()
// name: name of the asset.AssetDB instance
// pkg:  package name placed at the top of the file
func (a *AssetDB) Serialize(name, pkg string) string {
	var cdata strings.Builder
	// Set buffer size to 4k
	cdata.Grow(4096)

	// Write header
	header := `// DO NOT EDIT - Generated automatically
package %s

import "github.com/wailsapp/wails/v2/internal/assetdb"

var (
	%s *assetdb.AssetDB = assetdb.NewAssetDB()
)

// AssetsDB is a clean interface to the assetdb.AssetDB struct
type AssetsDB interface {
	Read(string) ([]byte, error)
	String(string) (string, error)
}

// Assets returns the asset database
func Assets() AssetsDB {
	return %s
}

func init() {
`
	cdata.WriteString(fmt.Sprintf(header, pkg, name, name))

	for aname, bytes := range a.db {
		cdata.WriteString(fmt.Sprintf("\t%s.AddAsset(\"%s\", []byte{",
			name,
			aname))

		l := len(bytes)
		if l == 0 {
			cdata.WriteString("0x00})\n")
			continue
		}

		// Convert each byte to hex
		for _, b := range bytes[:l-1] {
			cdata.WriteString(fmt.Sprintf("0x%x, ", b))
		}
		cdata.WriteString(fmt.Sprintf("0x%x})\n", bytes[l-1]))
	}

	cdata.WriteString(`}`)

	return cdata.String()
}
