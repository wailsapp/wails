package gomod

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/matryer/is"
)

func TestGetWailsVersion(t *testing.T) {
	tests := []struct {
		name      string
		goModText []byte
		want      *semver.Version
		wantErr   bool
	}{
		{"basic", []byte(basic), semver.MustParse("v2.0.0-beta.7"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWailsVersionFromModFile(tt.goModText)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWailsVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWailsVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateGoModVersion(t *testing.T) {
	is2 := is.New(t)

	type args struct {
		goModText      []byte
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"basic", args{[]byte(basic), "v2.0.0-beta.20"}, []byte(basicUpdated), false},
		{"basicmultiline", args{[]byte(multilineRequire), "v2.0.0-beta.20"}, []byte(multilineRequireUpdated), false},
		{"basicmultilinereplace", args{[]byte(multilineReplace), "v2.0.0-beta.20"}, []byte(multilineReplaceUpdated), false},
		{"basicmultilinereplaceblock", args{[]byte(multilineReplaceBlock), "v2.0.0-beta.20"}, []byte(multilineReplaceBlockUpdated), false},
		{"basicmultilinereplacenoversion", args{[]byte(multilineReplaceNoVersion), "v2.0.0-beta.20"}, []byte(multilineReplaceNoVersionUpdated), false},
		{"basicmultilinereplacenoversionblock", args{[]byte(multilineReplaceNoVersionBlock), "v2.0.0-beta.20"}, []byte(multilineReplaceNoVersionBlockUpdated), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateGoModVersion(tt.args.goModText, tt.args.currentVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGoModVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			is2.Equal(string(got), string(tt.want))
		})
	}
}

func TestGoModOutOfSync(t *testing.T) {
	is2 := is.New(t)

	type args struct {
		goModData      []byte
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"basic", args{[]byte(basic), "v2.0.0-beta.20"}, true, false},
		{"basicmultiline", args{[]byte(multilineRequire), "v2.0.0-beta.20"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoModOutOfSync(tt.args.goModData, tt.args.currentVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoModOutOfSync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			is2.Equal(got, tt.want)
		})
	}
}

const basicGo118 string = `module changeme

go 1.18

require github.com/wailsapp/wails/v2 v2.0.0-beta.7
`

const basicGo119 string = `module changeme

go 1.19

require github.com/wailsapp/wails/v2 v2.0.0-beta.7
`

func TestUpdateGoModGoVersion(t *testing.T) {
	is2 := is.New(t)

	type args struct {
		goModText      []byte
		currentVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		updated bool
	}{
		{"basic1.18", args{[]byte(basicGo118), "1.18"}, []byte(basicGo118), false},
		{"basic1.19", args{[]byte(basicGo119), "1.17"}, []byte(basicGo119), false},
		{"basic1.19", args{[]byte(basicGo119), "1.18"}, []byte(basicGo119), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, updated, err := SyncGoVersion(tt.args.goModText, tt.args.currentVersion)
			if err != nil {
				t.Errorf("UpdateGoModVersion() error = %v", err)
				return
			}
			if updated != tt.updated {
				t.Errorf("UpdateGoModVersion() updated = %t, want = %t", updated, tt.updated)
				return
			}
			is2.Equal(got, tt.want)
		})
	}
}
