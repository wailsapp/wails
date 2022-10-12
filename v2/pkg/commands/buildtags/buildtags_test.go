package buildtags

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		tags    string
		want    []string
		wantErr bool
	}{
		{
			name:    "should support single tags",
			tags:    "test",
			want:    []string{"test"},
			wantErr: false,
		},
		{
			name:    "should support space delimited tags",
			tags:    "test test2",
			want:    []string{"test", "test2"},
			wantErr: false,
		},
		{
			name:    "should support comma delimited tags",
			tags:    "test,test2",
			want:    []string{"test", "test2"},
			wantErr: false,
		},
		{
			name:    "should error if mixed tags",
			tags:    "test,test2 test3",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringify(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{
			name: "should support single tags",
			tags: []string{"test"},
			want: "test",
		},
		{
			name: "should support multiple tags",
			tags: []string{"test", "test2"},
			want: "test,test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Stringify(tt.tags); got != tt.want {
				t.Errorf("Stringify() = %v, want %v", got, tt.want)
			}
		})
	}
}
