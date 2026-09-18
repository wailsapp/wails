package generator

import (
	"bytes"
	_ "embed"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"testing"
	"updater/generator/types"
)

func TestInterfaceDoublePointer(t *testing.T) {

	i := is.New(t)

	var buf bytes.Buffer
	buf.Write([]byte(`
[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

	[uuid(76eceacb-0462-4d94-ac83-423a6793775e), object, pointer_default(unique)]
	interface ICoreWebView2 : IUnknown {

	  HRESULT add_NavigationStarting(
      [in] ICoreWebView2NavigationStartingEventHandler* eventHandler,
      [out] EventRegistrationToken* token);

	  [propget] HRESULT Settings([out, retval] ICoreWebView2Settings** settings);

	}

}
`))

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
			FileName: "ICoreWebView2.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2.go.txt"),
		},
	}

	require.Equal(t, len(files), len(expected))
	require.Equal(t, files[0].Content.String(), expected[0].Content.String())
	require.ElementsMatch(t, expected, files)

}
