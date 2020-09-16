package assetdb

import (
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestOpenLeadingSlash(t *testing.T) {
	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("/hello", helloworld)

	file, err := db.Open("/hello")
	// Ensure it does exist
	is.True(err == nil)

	buff := make([]byte, len(helloworld))
	n, err := file.Read(buff)
	fmt.Printf("Error %v\n", err)
	is.True(err == nil)
	is.Equal(n, len(helloworld))
	result := string(buff)

	// Ensure the string is blank
	is.Equal(result, string(helloworld))
}

func TestOpen(t *testing.T) {
	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("/hello", helloworld)

	file, err := db.Open("/hello")

	// Ensure it does exist
	is.True(err == nil)

	buff := make([]byte, len(helloworld))
	n, err := file.Read(buff)
	is.True(err == nil)
	is.Equal(n, len(helloworld))
	result := string(buff)

	// Ensure the string is blank
	is.Equal(result, string(helloworld))
}

func TestReaddir(t *testing.T) {
	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("/hello", helloworld)
	db.AddAsset("/directory/hello", helloworld)
	db.AddAsset("/directory/subdirectory/hello", helloworld)

	dir, err := db.Open("/doesntexist")
	is.True(err == os.ErrNotExist)
	ents, err := dir.Readdir(-1)
	is.Equal([]os.FileInfo{}, ents)

	dir, err = db.Open("/")
	is.True(dir != nil)
	is.True(err == nil)
	ents, err = dir.Readdir(-1)
	is.True(err == nil)
	is.Equal(3, len(ents))
}

func TestReaddirSubdirectory(t *testing.T) {
	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("/hello", helloworld)
	db.AddAsset("/directory/hello", helloworld)
	db.AddAsset("/directory/subdirectory/hello", helloworld)

	expected := []os.FileInfo{
		FI{name: "hello", dir: false, size: len(helloworld)},
		FI{name: "subdirectory", dir: true, size: -1},
	}

	dir, err := db.Open("/directory")
	is.True(dir != nil)
	is.True(err == nil)
	ents, err := dir.Readdir(-1)
	is.Equal(expected, ents)

	// Check sub-subdirectory
	dir, err = db.Open("/directory/subdirectory")
	is.True(dir != nil)
	is.True(err == nil)
	ents, err = dir.Readdir(-1)
	is.True(err == nil)
	is.Equal([]os.FileInfo{FI{name: "hello", size: len(helloworld)}}, ents)
}
