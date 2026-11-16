package generator

import (
	"bytes"
	_ "embed"
	"testing"
	"updater/generator/types"
)

func TestInterfaceStringPointer(t *testing.T) {

	var buf bytes.Buffer
	buf.Write([]byte(`
[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

	[uuid(d60ac92c-37a6-4b26-a39e-95cfe59047bb), object, pointer_default(unique)]
	interface ICoreWebView2CustomSchemeRegistration : IUnknown {
		HRESULT SetAllowedOrigins(
			[in] UINT32 allowedOriginsCount,
			[in] LPCWSTR* allowedOrigins);

	}
	
}
`))

	files := generateFixture(t, buf.Bytes())

	//for _, file := range files {
	//	os.WriteFile("testfiles/"+file.FileName+".txt", file.Content.Bytes(), 0644)
	//}

	expected := []*types.GeneratedFile{
		{
			FileName: "ICoreWebView2CustomSchemeRegistration.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2CustomSchemeRegistration.go.txt"),
		},
	}

	checkGeneratedFiles(t, files, expected)

}

/*
[uuid(76eceacb-0462-4d94-ac83-423a6793775e), object, pointer_default(unique)]
interface ICoreWebView2 : IUnknown {

  /// The `ICoreWebView2Settings` object contains various modifiable settings
  /// for the running WebView.

  [propget] HRESULT Settings([out, retval] ICoreWebView2Settings** settings);

*/
