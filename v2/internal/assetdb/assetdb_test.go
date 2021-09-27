package assetdb

import "testing"
import "github.com/matryer/is"

func TestExistsAsBytes(t *testing.T) {

	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("hello", helloworld)

	result, err := db.Read("hello")

	is.True(err == nil)
	is.Equal(result, helloworld)
}

func TestNotExistsAsBytes(t *testing.T) {

	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("hello4", helloworld)

	result, err := db.Read("hello")

	is.True(err != nil)
	is.True(result == nil)
}

func TestExistsAsString(t *testing.T) {

	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("hello", helloworld)

	result, err := db.String("hello")

	// Ensure it exists
	is.True(err == nil)

	// Ensure the string is the same as the byte slice
	is.Equal(result, "Hello, World!")
}

func TestNotExistsAsString(t *testing.T) {

	is := is.New(t)

	var helloworld = []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}

	db := NewAssetDB()
	db.AddAsset("hello", helloworld)

	result, err := db.String("help")

	// Ensure it doesn't exist
	is.True(err != nil)

	// Ensure the string is blank
	is.Equal(result, "")
}
