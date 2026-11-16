package generator

import (
	"bytes"
	_ "embed"
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

	var buf bytes.Buffer
	buf.Write(interfaceTestData)

	files := generateFixture(t, buf.Bytes())

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

	checkGeneratedFiles(t, files, expected)

}
