package generator

import (
	"bytes"
	_ "embed"
	"testing"
	"updater/generator/types"
)

func TestInterfaceDoublePointer(t *testing.T) {

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

	files := generateFixture(t, buf.Bytes())

	expected := []*types.GeneratedFile{
		{
			FileName: "ICoreWebView2.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2.go.txt"),
		},
	}

	checkGeneratedFiles(t, files, expected)

}
