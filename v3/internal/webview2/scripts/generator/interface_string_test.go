package generator

import (
	"bytes"
	_ "embed"
	"testing"
	"updater/generator/types"
)

func TestInterfaceString(t *testing.T) {

	var buf bytes.Buffer
	buf.Write([]byte(`
[uuid(26d34152-879f-4065-bea2-3daa2cfadfb8), version(1.0)]
library WebView2 {

	[uuid(da86b8a1-bdf3-4f11-9955-528cefa59727), object, pointer_default(unique)]
	interface ICoreWebView2FrameInfo : IUnknown {
	
	  [propget] HRESULT Name([out, retval] LPWSTR* name);
	  [propget] HRESULT Source([out, retval] LPWSTR* source);
	}

}
`))

	files := generateFixture(t, buf.Bytes())

	//for _, file := range files {
	//	os.WriteFile("testfiles/"+file.FileName+".txt", file.Content.Bytes(), 0644)
	//}

	expected := []*types.GeneratedFile{
		{
			FileName: "ICoreWebView2FrameInfo.go",
			Package:  "webview2",
			Content:  testfile("ICoreWebView2FrameInfo.go.txt"),
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