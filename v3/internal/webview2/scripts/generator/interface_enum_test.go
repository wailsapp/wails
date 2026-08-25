package generator

import (
	"bytes"
	_ "embed"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"testing"
	"updater/generator/types"
)

var interfaceEnumTestData = []byte(`

[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

[v1_enum]
typedef enum COREWEBVIEW2_KEY_EVENT_KIND {
  COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN,
  COREWEBVIEW2_KEY_EVENT_KIND_KEY_UP,
  COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN,
  COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_UP,
} COREWEBVIEW2_KEY_EVENT_KIND;

[uuid(9f760f8a-fb79-42be-9990-7b56900fa9c7), object, pointer_default(unique)]
interface ICoreWebView2AcceleratorKeyPressedEventArgs : IUnknown {
  [propget] HRESULT KeyEventKind([out, retval] COREWEBVIEW2_KEY_EVENT_KIND* keyEventKind);
}
}`)

func TestInterfaceEnum(t *testing.T) {

	i := is.New(t)

	var buf bytes.Buffer
	buf.Write(interfaceEnumTestData)

	idl, err := Parser.Parse("", &buf)
	i.NoErr(err)

	err = idl.Process()
	i.NoErr(err)

	files, err := idl.Generate()
	i.NoErr(err)

	// Remove the `com.go` filename
	files = files[1:]

	//for _, file := range files {
	//	os.WriteFile("testfiles/"+file.FileName+".txt", file.Content.Bytes(), 0644)
	//}

	expected := []*types.GeneratedFile{
		{
			FileName: "COREWEBVIEW2_KEY_EVENT_KIND.go",
			Package:  "webview2",
			Content:  testfile("COREWEBVIEW2_KEY_EVENT_KIND.go.txt"),
		},
		{
			FileName: "ICoreWebView2AcceleratorKeyPressedEventArgs.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2AcceleratorKeyPressedEventArgs.go.txt"),
		},
	}

	require.Equal(t, len(files), len(expected))
	require.ElementsMatch(t, expected, files)

}
