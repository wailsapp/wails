package assetserver

import "testing"

func TestGetMimetype(t *testing.T) {
	type args struct {
		filename string
		data     []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"css", args{"test.css", []byte("body{margin:0;padding:0;background-color:#d579b2}#app{font-family:Avenir,Helvetica,Arial,sans-serif;-webkit-font-smoothing:antialiased;-moz-osx-font-smoothing:grayscale;text-align:center;color:#2c3e50;background-color:#ededed}#nav{padding:30px}#nav a{font-weight:700;color:#2c\n3e50}#nav a.router-link-exact-active{color:#42b983}.hello[data-v-4e26ad49]{margin:10px 0}")}, "text/css; charset=utf-8"},
		{"js", args{"test.js", []byte("let foo = 'bar'; console.log(foo);")}, "text/javascript; charset=utf-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMimetype(tt.args.filename, tt.args.data); got != tt.want {
				t.Errorf("GetMimetype() = %v, want %v", got, tt.want)
			}
		})
	}
}
