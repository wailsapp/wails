package commands

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestParseIconSizes(t *testing.T) {
	cases := []struct {
		input   string
		want    []int
		invalid bool
	}{
		{"1,16,256", []int{1, 16, 256}, false},
		{"0,16", []int{16}, false},
		{"-16", nil, true},
		{"257", nil, true},
		{"1000000000", nil, true},
		{"0,0", nil, true},
		{"", nil, true},
		{"bad", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := parseSizes(tc.input)
			if tc.invalid {
				if err == nil {
					t.Errorf("accepted invalid sizes %q: %v", tc.input, got)
				}
			} else if err != nil || !slices.Equal(got, tc.want) {
				t.Errorf("sizes=%v, error=%v; want %v", got, err, tc.want)
			}
		})
	}
}

func TestGenerateIconsRejectsInvalidSizesBeforeReadingInput(t *testing.T) {
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	root := t.TempDir()
	output := filepath.Join(root, "existing.ico")
	if err := os.WriteFile(output, []byte("previous icon"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, sizes := range []string{"-16", "257", "0,0"} {
		err := GenerateIcons(&IconsOptions{Input: filepath.Join(root, "missing.png"), WindowsFilename: output, Sizes: sizes})
		if err == nil || !strings.Contains(err.Error(), "size") {
			t.Errorf("sizes %q: error=%v, want size validation failure", sizes, err)
		}
		data, err := os.ReadFile(output)
		if err != nil || string(data) != "previous icon" {
			t.Fatalf("output=%q, error=%v", data, err)
		}
	}
}
