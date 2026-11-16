package installerbuild

import (
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestCompilerEnvironment(t *testing.T) {
	for _, test := range []struct {
		name, host  string
		environment []string
		fallback    bool
	}{
		{"empty", "darwin", nil, true},
		{"C", "darwin", []string{"LC_ALL=C", "LANG=en_US.UTF-8"}, true},
		{"POSIX", "darwin", []string{"LANG=POSIX"}, true},
		{"other encoding", "darwin", []string{"LANG=en_US.ISO8859-1"}, true},
		{"empty higher categories", "darwin", []string{"LC_ALL=", "LC_CTYPE=", "LANG=de_DE.UTF-8"}, false},
		{"all precedence", "darwin", []string{"LC_ALL=ja_JP.UTF-8", "LC_CTYPE=C", "LANG=C"}, false},
		{"ctype precedence", "darwin", []string{"LC_CTYPE=en_GB.utf8", "LANG=C"}, false},
		{"ctype nonUTF8 precedence", "darwin", []string{"LC_CTYPE=C", "LANG=en_US.UTF-8"}, true},
		{"bare UTF8", "darwin", []string{"LC_CTYPE=UTF-8"}, false},
		{"locale modifier", "darwin", []string{"LANG=en_US.UTF-8@custom"}, false},
		{"other category preserved", "darwin", []string{"LANG=C", "LC_MESSAGES=fr_FR.UTF-8", "PATH=/tools", "WAILS_SETTING=kept"}, true},
		{"Linux unchanged", "linux", []string{"LC_ALL=C", "LANG=C"}, false},
		{"Windows unchanged", "windows", []string{"LC_ALL=C", "LANG=C"}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := append([]string(nil), test.environment...)
			result := CompilerEnvironment(test.host, test.environment)
			if !reflect.DeepEqual(original, test.environment) {
				t.Fatal("mutated caller's environment")
			}
			if !test.fallback {
				if !reflect.DeepEqual(result, test.environment) {
					t.Fatalf("changed valid/other-host environment: %v", result)
				}
				return
			}
			for _, key := range []string{"LC_ALL", "LANG"} {
				if !slices.Contains(result, key+"=en_US.UTF-8") {
					t.Errorf("missing fallback %s: %v", key, result)
				}
			}
			for _, entry := range original {
				if strings.HasPrefix(entry, "LC_ALL=") || strings.HasPrefix(entry, "LANG=") {
					continue
				}
				if !slices.Contains(result, entry) {
					t.Errorf("lost unrelated entry %q", entry)
				}
			}
		})
	}
}

func TestCompilerEnvironmentDoesNotMutateProcessOrBackingSlice(t *testing.T) {
	t.Setenv("LC_ALL", "C")
	environment := []string{"LC_ALL=C", "PATH=/tools"}
	for _, host := range []string{"darwin", "linux", "windows"} {
		result := CompilerEnvironment(host, environment)
		result[0] = "LC_ALL=changed"
		if environment[0] != "LC_ALL=C" || os.Getenv("LC_ALL") != "C" {
			t.Fatalf("%s mutated caller/process", host)
		}
	}
}
