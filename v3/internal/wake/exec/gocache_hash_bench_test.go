package exec

import "testing"

var commandHashResult string

func BenchmarkGoCmdHash(b *testing.B) {
	for _, tc := range []struct {
		name, command, want string
	}{
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"tidy", "go mod tidy", "5aebdf7446ebca99f13ca2b455880d5d55233a2d83498cf097922da02828b02f"},
		{"build", "go build -tags production -o \"bin/my app\" ./cmd/app", "ac7ef54bd368432e1a8f47960fc7859e44a7f43203dff888a3777478ce24887a"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			if got := goCmdHash(tc.command); got != tc.want {
				b.Fatalf("command hash = %q, want %q", got, tc.want)
			}
			b.ReportAllocs()
			for b.Loop() {
				commandHashResult = goCmdHash(tc.command)
			}
		})
	}
}
