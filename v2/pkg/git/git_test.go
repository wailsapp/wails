package git

import (
	"testing"
)

func TestEscapeName1(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Escape Apostrophe",
			args: args{
				str: `John O'Keefe`,
			},
			want: `John O'Keefe`,
		},
		{
			name: "Escape backslash",
			args: args{
				str: `MYDOMAIN\USER`,
			},
			want: `MYDOMAIN\\USER`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EscapeName(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("EscapeName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EscapeName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
