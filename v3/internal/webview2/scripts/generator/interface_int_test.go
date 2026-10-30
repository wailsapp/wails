package generator

import (
	"bytes"
	_ "embed"
	"testing"
	"updater/generator/types"
)

func TestInterfaceInt(t *testing.T) {

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

	files := generateFixture(t, buf.Bytes())

	expected := []*types.GeneratedFile{
		{
			FileName: "ICoreWebView2ProcessFailedEventArgs2.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2ProcessFailedEventArgs2.go.txt"),
		},
	}

	checkGeneratedFiles(t, files, expected)

}
