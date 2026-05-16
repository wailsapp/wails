package generator

import (
	"bytes"
	_ "embed"
	"github.com/matryer/is"
	"github.com/stretchr/testify/require"
	"testing"
	"updater/generator/types"
)

var interfaceTestData = []byte(`

[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

[uuid(A0D6DF20-3B92-416D-AA0C-437A9C727857), object, pointer_default(unique)]
interface ICoreWebView2_3 : ICoreWebView2_2 {
[propget] HRESULT IsSuspended([out, retval] BOOL* isSuspended);
}
}`)

func TestInterfaceBool(t *testing.T) {

	i := is.New(t)

	var buf bytes.Buffer
	buf.Write(interfaceTestData)

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
			FileName: "ICoreWebView2_3.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2_3.go.txt"),
		},
	}

	require.Equal(t, len(files), len(expected))
	require.ElementsMatch(t, expected, files)

}
