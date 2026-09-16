//go:build android && !production

package application

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

var androidDevOutputOnce sync.Once

// Android has no inherited terminal. In a CLI-owned debug launch, forward Go
// standard streams into logcat so the session can attribute and display them.
func configureDevOutput() {
	androidDevOutputOnce.Do(func() {
		redirect := func(destination *os.File, level, label string) {
			reader, writer, err := os.Pipe()
			if err != nil {
				return
			}
			// Redirect the descriptor itself: log.Default and user loggers may
			// already hold os.Stderr, and Go panic output writes directly to fd 2.
			err = unix.Dup3(int(writer.Fd()), int(destination.Fd()), 0)
			writer.Close()
			if err != nil {
				reader.Close()
				return
			}
			go func() {
				defer reader.Close()
				buffer := bufio.NewReaderSize(reader, 8192)
				for {
					line, err := buffer.ReadSlice('\n')
					if len(line) > 0 {
						androidLogf(level, "[Go %s] %s", label, strings.TrimSuffix(string(line), "\n"))
					}
					if err != nil && err != bufio.ErrBufferFull {
						if err != io.EOF {
							androidLogf("error", "read Go %s: %v", label, err)
						}
						return
					}
				}
			}()
		}
		redirect(os.Stdout, "info", "stdout")
		redirect(os.Stderr, "error", "stderr")
	})
}
