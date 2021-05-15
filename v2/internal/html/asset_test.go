package html

import "testing"

func TestAsset_minifiedData(t *testing.T) {
	type fields struct {
		Type string
		Path string
		Data string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "multi-line tag",
			fields: fields{
				Type: AssetTypes.HTML,
				Path: "foo.html",
				Data: "<link\n  rel=\"stylesheet\"\n href=\"src/foo.css\"\n>\n",
			},
			want: "data:text/html;charset=utf-8,%3Clink%20rel=%22stylesheet%22%20href=%22src%2ffoo.css%22%20%3E%20",
		},
		{
			name: "multi-line tag no spaces",
			fields: fields{
				Type: AssetTypes.HTML,
				Path: "foo.html",
				Data: "<link\nrel=\"stylesheet\"\nhref=\"src/foo.css\"\n>\n",
			},
			want: "data:text/html;charset=utf-8,%3Clink%20rel=%22stylesheet%22%20href=%22src%2ffoo.css%22%20%3E%20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Asset{
				Type: tt.fields.Type,
				Path: tt.fields.Path,
				Data: tt.fields.Data,
			}
			got, err := a.minifiedData()
			if (err != nil) != tt.wantErr {
				t.Errorf("Asset.minifiedData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Asset.minifiedData() = %v, want %v", got, tt.want)
			}
		})
	}
}
