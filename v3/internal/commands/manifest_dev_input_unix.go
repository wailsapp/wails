//go:build !windows

package commands

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

// Poll instead of leaving a scanner blocked on the process-wide stdin after
// the development session ends. Canonical terminal editing/echo stay intact.
func readManifestDevInput(ctx context.Context, input *os.File, send func(string) bool) error {
	poll := []unix.PollFd{{Fd: int32(input.Fd()), Events: unix.POLLIN}}
	var buffered string
	var data [4096]byte
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		count, err := unix.Poll(poll, 100)
		if errors.Is(err, unix.EINTR) {
			continue
		}
		if err != nil {
			return err
		}
		if count == 0 {
			continue
		}
		if poll[0].Revents&unix.POLLNVAL != 0 {
			return os.ErrClosed
		}
		n, err := unix.Read(int(input.Fd()), data[:])
		if err != nil {
			return err
		}
		if n == 0 {
			return io.EOF
		}
		buffered += string(data[:n])
		for {
			line, rest, ok := strings.Cut(buffered, "\n")
			if !ok {
				break
			}
			buffered = rest
			if !send(strings.TrimSuffix(line, "\r")) {
				return ctx.Err()
			}
		}
		if len(buffered) > 4096 {
			buffered = ""
		}
	}
}
