package main

import (
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// TypeMappingService is a service for testing type mappings.
type TypeMappingService struct{}

type Number interface {
	~int8 | ~int16
}

type GenericType[T Number] struct {
	ID   int
	Data T
	List []T
}

// AllTypes holds a field for every fundamental Go-to-TS type mapping.
type AllTypes struct {
	Bool    bool
	String  string
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	Byte    byte
	Rune    rune

	ByteSlice []byte // maps to string (base64)
	Bytes2D   [][]byte
	SliceStr  []string
	ArrInt4   [4]int

	MapStrInt   map[string]int
	MapByteBool map[byte]bool

	PtrStr    *string // pointer means nullable
	PtrInt    *int
	PtrStruct *struct {
		MapIntBool map[int]bool
		Number     int
	}

	Struct struct {
		Any         any
		MapByteBool map[byte]bool
	}

	GenericType GenericType[int8]
	Time        time.Time // maps to any
	Err         error
}

// GetTypes returns an AllTypes value to exercise all type mappings.
func (*TypeMappingService) GetTypes() AllTypes {
	return AllTypes{}
}

func main() {
	app := application.New(application.Options{
		Services: []application.Service{
			application.NewService(&TypeMappingService{}),
		},
	})

	app.Window.New()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
