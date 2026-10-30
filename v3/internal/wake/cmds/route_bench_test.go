package cmds

import "testing"

var routedCommand Executor

func BenchmarkRoute(b *testing.B) {
	for _, command := range []string{"go build -o bin/app ./cmd/app", "echo ordinary shell command"} {
		b.Run(command, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				routedCommand = Route(command, RouteOptions{Dir: "/project"})
			}
		})
	}
}
