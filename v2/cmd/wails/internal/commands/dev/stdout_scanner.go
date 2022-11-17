package dev

import (
	"bufio"
	"net/url"
	"os"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/logutils"
)

// stdoutScanner acts as a stdout target that will scan the incoming
// data to find out the vite server url
type stdoutScanner struct {
	ViteServerURLChan chan string
}

// NewStdoutScanner creates a new stdoutScanner
func NewStdoutScanner() *stdoutScanner {
	return &stdoutScanner{
		ViteServerURLChan: make(chan string, 2),
	}
}

// Write bytes to the scanner. Will copy the bytes to stdout
func (s *stdoutScanner) Write(data []byte) (n int, err error) {
	input := stripansi.Strip(string(data))
	match := strings.Index(input, "Local:")
	if match != -1 {
		sc := bufio.NewScanner(strings.NewReader(input))
		for sc.Scan() {
			line := sc.Text()
			index := strings.Index(line, "Local:")
			if index == -1 || len(line) < 7 {
				continue
			}
			viteServerURL := strings.TrimSpace(line[index+6:])
			logutils.LogGreen("Vite Server URL: %s", viteServerURL)
			_, err := url.Parse(viteServerURL)
			if err != nil {
				logutils.LogRed(err.Error())
			} else {
				s.ViteServerURLChan <- viteServerURL
			}
		}
	}
	return os.Stdout.Write(data)
}
