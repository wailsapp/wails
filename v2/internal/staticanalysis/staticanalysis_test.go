package staticanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
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
					EmbedPath: "frontend/dist",
					All:       true,
				},
				{
					EmbedPath: "frontend/static",
					All:       false,
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
			require.Equal(t, len(tt.want), len(got))
			for index, g := range got {
				require.Equal(t, tt.want[index].EmbedPath, g.EmbedPath)
				require.Equal(t, tt.want[index].All, g.All)
			}
		})
	}
}
