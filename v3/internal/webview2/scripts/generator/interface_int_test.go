package generator

import (
	"bytes"
	_ "embed"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"testing"
	"updater/generator/types"
)

func TestInterfaceInt(t *testing.T) {

	i := is.New(t)

	var buf bytes.Buffer
	buf.Write([]byte(`
[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

	[uuid(4dab9422-46fa-4c3e-a5d2-41d2071d3680), object, pointer_default(unique)]
	interface ICoreWebView2ProcessFailedEventArgs2 : ICoreWebView2ProcessFailedEventArgs {
	
		[propget] HRESULT ExitCode(
		[out, retval] int* exitCode);
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
			FileName: "ICoreWebView2ProcessFailedEventArgs2.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2ProcessFailedEventArgs2.go.txt"),
		},
	}

	require.Equal(t, len(files), len(expected))
	require.ElementsMatch(t, expected, files)

}
