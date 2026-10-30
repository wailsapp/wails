//go:build !windows

package cmds

import (
	"testing"
)

func TestRouteGoArguments(t *testing.T) {
	t.Setenv("WAKE_TEST_ARGUMENT", "two words")
	for _, tc := range []struct {
		name, command, want string
	}{
		{"count/invalid", "go test -count not-a-number ./...", "test\n-count\nnot-a-number\n./...\n"},
		{"count/zero", "go test -count 0 ./...", "test\n-count\n0\n./...\n"},
		{"count/negative", "go test -count -1 ./...", "test\n-count\n-1\n./...\n"},
		{"count/one", "go test -count 1 ./...", "test\n-count\n1\n./...\n"},
		{"count/hundred", "go test -count 100 ./...", "test\n-count\n100\n./...\n"},
		{"count/thousand", "go test -count 1000 ./...", "test\n-count\n1000\n./...\n"},
		{"test options", `go test -v -race -cover -tags "gtk4 production" -run "^Test\w+$" ./...`, "test\n-v\n-race\n-cover\n-tags\ngtk4 production\n-run\n^Test\\w+$\n./...\n"},
		{"test default", "go test", "test\n"},
		{"count/positive", "go test -count 10 ./...", "test\n-count\n10\n./...\n"},
		{"count/zero then positive", "go test -count 0 -count 2 ./...", "test\n-count\n0\n-count\n2\n./...\n"},
		{"count/positive then zero", "go test -count 2 -count 0 ./...", "test\n-count\n2\n-count\n0\n./...\n"},
		{"missing count", "go test -count", "test\n-count\n"},
		{"missing run", "go test -run", "test\n-run\n"},
		{"missing tags", "go test -tags", "test\n-tags\n"},
		{"package before missing count", "go test ./... -count", "test\n./...\n-count\n"},
		{"package before missing run", "go test ./... -run", "test\n./...\n-run\n"},
		{"package before missing tags", "go test ./... -tags", "test\n./...\n-tags\n"},
		{"flag-like run value", `go test -run "-run"`, "test\n-run\n-run\n"},
		{"test packages", "go test ./first ./second", "test\n./first\n./second\n"},
		{"test separate flag value", "go test -timeout 5s ./first ./second", "test\n-timeout\n5s\n./first\n./second\n"},
		{"run application flags", "go run ./app -tags runtime-tag -ldflags -s -trimpath", "run\n./app\n-tags\nruntime-tag\n-ldflags\n-s\n-trimpath\n"},
		{"run build and application tags", "go run -tags build-tag ./app -tags runtime-tag", "run\n-tags\nbuild-tag\n./app\n-tags\nruntime-tag\n"},
		{"install packages", "go install ./first ./second", "install\n./first\n./second\n"},
		{"install flags", "go install -trimpath -tags=production ./cmd/tool", "install\n-trimpath\n-tags=production\n./cmd/tool\n"},
		{"vet packages", "go vet ./first ./second", "vet\n./first\n./second\n"},
		{"vet analyzer", "go vet -printf=false ./sub", "vet\n-printf=false\n./sub\n"},
		{"fmt dry run", "go fmt -n ./...", "fmt\n-n\n./...\n"},
		{"fmt packages", `go fmt "./path with spaces" ./second`, "fmt\n./path with spaces\n./second\n"},
		{"tidy diff", "go mod tidy -diff", "mod\ntidy\n-diff\n"},
		{"tidy compatibility", "go mod tidy -compat=1.25 -v", "mod\ntidy\n-compat=1.25\n-v\n"},
		{"unquoted variable", "go run $WAKE_TEST_ARGUMENT", "run\ntwo\nwords\n"},
		{"quoted variable", `go run "$WAKE_TEST_ARGUMENT"`, "run\ntwo words\n"},
		{"quoted substitution", `go run "$(printf 'two words')"`, "run\ntwo words\n"},
		{"backtick substitution", "go run \"`printf 'two words'`\"", "run\ntwo words\n"},
		{"single quoted literal", "go run '$WAKE_TEST_ARGUMENT'", "run\n$WAKE_TEST_ARGUMENT\n"},
		{"escaped literal", `go run \$WAKE_TEST_ARGUMENT`, "run\n$WAKE_TEST_ARGUMENT\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			base := installFakeGo(t, `printf '%s\n' "$@"`)
			if got := captureRoutedOutput(t, tc.command, base); got != tc.want {
				t.Errorf("forwarded Go arguments = %q, want %q", got, tc.want)
			}
		})
	}
}
