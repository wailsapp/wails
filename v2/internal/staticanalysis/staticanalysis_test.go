package staticanalysis

import (
	"reflect"
	"testing"
)

func TestGetEmbedDetails(t *testing.T) {
	type args struct {
		sourcePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []*EmbedDetails
		wantErr bool
	}{
		{
			name: "GetEmbedDetails",
			args: args{
				sourcePath: "test/standard",
			},
			want: []*EmbedDetails{
				{
					SourcePath: "frontend/dist",
					All:        true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEmbedDetails(tt.args.sourcePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEmbedDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEmbedDetails() got = %v, want %v", got, tt.want)
			}
		})
	}
}
