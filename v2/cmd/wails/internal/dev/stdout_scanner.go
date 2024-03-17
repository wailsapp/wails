package dev

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/wailsapp/wails/v2/cmd/wails/internal/logutils"
	"golang.org/x/mod/semver"
)

// stdoutScanner acts as a stdout target that will scan the incoming
// data to find out the vite server url
type stdoutScanner struct {
	ViteServerURLChan  chan string
	ViteServerVersionC chan string
	versionDetected    bool
}

// NewStdoutScanner creates a new stdoutScanner
func NewStdoutScanner() *stdoutScanner {
	return &stdoutScanner{
		ViteServerURLChan:  make(chan string, 2),
		ViteServerVersionC: make(chan string, 2),
	}
}

// Write bytes to the scanner. Will copy the bytes to stdout
func (s *stdoutScanner) Write(data []byte) (n int, err error) {
	input := stripansi.Strip(string(data))
	if !s.versionDetected {
		v, err := detectViteVersion(input)
		if v != "" || err != nil {
			if err != nil {
				logutils.LogRed("ViteStdoutScanner: %s", err)
				v = "v0.0.0"
			}
			s.ViteServerVersionC <- v
			s.versionDetected = true
		}
	}

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

func detectViteVersion(line string) (string, error) {
	s := strings.Split(strings.TrimSpace(line), " ")
	if strings.ToLower(s[0]) != "vite" {
		return "", nil
	}

	if len(line) < 2 {
		return "", fmt.Errorf("unable to parse vite version")
	}

	v := s[1]
	if !semver.IsValid(v) {
		return "", fmt.Errorf("%s is not a valid vite version string", v)
	}

	return v, nil
}
