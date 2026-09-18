package generator

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"updater/generator/types"
)

//go:embed testfiles/*
var testfiles embed.FS

func testfile(path string) *bytes.Buffer {
	f, err := testfiles.ReadFile("testfiles/" + path)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(f)
}

func makeOutput(input string) *bytes.Buffer {
	var buf bytes.Buffer
	// Normalise newlines
	input = strings.ReplaceAll(input, "\r\n", "\n")
	buf.Write([]byte(input))
	return &buf
}

var testData = []byte(`

[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

[v1_enum]
typedef enum COREWEBVIEW2_PREFERRED_COLOR_SCHEME {
    /// Auto color scheme.
    COREWEBVIEW2_PREFERRED_COLOR_SCHEME_AUTO,

    /// Light color scheme.
    COREWEBVIEW2_PREFERRED_COLOR_SCHEME_LIGHT,

    /// Dark color scheme.
    COREWEBVIEW2_PREFERRED_COLOR_SCHEME_DARK
} COREWEBVIEW2_PREFERRED_COLOR_SCHEME;


[v1_enum]
typedef enum COREWEBVIEW2_PREFERRED_COLOR_SCHEME1 {
    /// Auto color scheme.
    COREWEBVIEW2_PREFERRED_COLOR_SCHEME_AUTO1 = 1,

    /// Light color scheme.
    COREWEBVIEW2_PREFERRED_COLOR_SCHEME_LIGHT1 = 2,

    /// Dark color scheme.
    COREWEBVIEW2_PREFERRED_COLOR_SCHEME_DARK1 = 3,
} COREWEBVIEW2_PREFERRED_COLOR_SCHEME;


[v1_enum]
typedef enum COREWEBVIEW2_PREFERRED_COLOR_SCHEME2 {
   /// Auto color scheme.
   COREWEBVIEW2_PREFERRED_COLOR_SCHEME_AUTO2 = 1 << 1,

   /// Light color scheme.
   COREWEBVIEW2_PREFERRED_COLOR_SCHEME_LIGHT2 = 1 << 2,

   /// Dark color scheme.
   COREWEBVIEW2_PREFERRED_COLOR_SCHEME_DARK2 = 1 << 3
} COREWEBVIEW2_PREFERRED_COLOR_SCHEME;

}`)

func TestEnum(t *testing.T) {

	i := is.New(t)

	var buf bytes.Buffer
	buf.Write(testData)

	idl, err := Parser.Parse("", &buf)
	i.NoErr(err)

	err = idl.Process()
	i.NoErr(err)

	files, err := idl.Generate()
	i.NoErr(err)

	// Remove the `com.go` filename
	files = files[1:]

	expected := []*types.GeneratedFile{
		{
			FileName: "COREWEBVIEW2_PREFERRED_COLOR_SCHEME.go",
			Package:  "webview2",
			Content:  testfile("COREWEBVIEW2_PREFERRED_COLOR_SCHEME.go.txt"),
		},
		{
			FileName: "COREWEBVIEW2_PREFERRED_COLOR_SCHEME1.go",
			Package:  "webview2",
			Content:  testfile("COREWEBVIEW2_PREFERRED_COLOR_SCHEME1.go.txt"),
		},
		{
			FileName: "COREWEBVIEW2_PREFERRED_COLOR_SCHEME2.go",
			Package:  "webview2",
			Content:  testfile("COREWEBVIEW2_PREFERRED_COLOR_SCHEME2.go.txt"),
		},
	}

	require.ElementsMatch(t, files, expected)

}
