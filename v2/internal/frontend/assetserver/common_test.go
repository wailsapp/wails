package assetserver

import (
	"reflect"
	"testing"
)

const realHTML = `<html>

<head>
  <title>test3</title>
  <meta name="wails-options" content="noautoinject">
  <link rel="stylesheet" href="/main.css">
</head>

<body data-wails-drag>
  <div class="logo"></div>
  <div class="result" id="result">Please enter your name below �</div>
  <div class="input-box" id="input" data-wails-no-drag>
    <input class="input" id="name" type="text" autocomplete="off">
    <button class="btn" onclick="greet()">Greet</button>
  </div>

  <script src="/main.js"></script>
</body>

</html>
`

func genMeta(content string) []byte {
	return []byte("<html><head><meta name=\"wails-options\" content=\"" + content + "\"></head><body></body></html>")
}

func genOptions(runtime bool, bindings bool) *Options {
	return &Options{
		disableRuntimeInjection: runtime,
		disableIPCInjection:     bindings,
	}
}

func Test_extractOptions(t *testing.T) {
	tests := []struct {
		name      string
		htmldata  []byte
		want      *Options
		wantError bool
	}{
		{"empty", []byte(""), &Options{}, false},
		{"bad data", []byte("<"), &Options{}, false},
		{"bad options", genMeta("noauto"), genOptions(false, false), false},
		{"realhtml", []byte(realHTML), genOptions(true, true), false},
		{"noautoinject", genMeta("noautoinject"), genOptions(true, true), false},
		{"noautoinjectipc", genMeta("noautoinjectipc"), genOptions(false, true), false},
		{"noautoinjectruntime", genMeta("noautoinjectruntime"), genOptions(true, false), false},
		{"spaces", genMeta("  noautoinjectruntime  "), genOptions(true, false), false},
		{"multiple", genMeta("noautoinjectruntime,noautoinjectipc"), genOptions(true, true), false},
		{"multiple spaces", genMeta(" noautoinjectruntime, noautoinjectipc "), genOptions(true, true), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlNode, err := getHTMLNode(tt.htmldata)
			if !tt.wantError && err != nil {
				t.Errorf("did not want error but got it")
			}
			got, err := extractOptions(htmlNode)
			if !tt.wantError && err != nil {
				t.Errorf("did not want error but got it")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
