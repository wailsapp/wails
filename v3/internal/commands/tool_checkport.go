package commands

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

type ToolCheckPortOptions struct {
	URL  string `name:"u" description:"URL to check"`
	Host string `name:"h" description:"Host to check" default:"localhost"`
	Port int    `name:"p" description:"Port to check"`
}

func isPortOpen(ip string, port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, fmt.Sprintf("%d", port)), timeout)
	if err != nil {
		return false
	}

	// If there's no error, close the connection and return true.
	if conn != nil {
		_ = conn.Close()
	}
	return true
}

func ToolCheckPort(options *ToolCheckPortOptions) error {
	DisableFooter = true

	if options.URL != "" {
		// Parse URL
		u, err := url.Parse(options.URL)
		if err != nil {
			return err
		}
		options.Host = u.Hostname()
		options.Port, err = strconv.Atoi(u.Port())
		if err != nil {
			return err
		}
	} else {
		if options.Port == 0 {
			return fmt.Errorf("please use the -p flag to specify a port")
		}
		if !isPortOpen(options.Host, options.Port) {
			return fmt.Errorf("port %d is not open on %s", options.Port, options.Host)
		}
	}
	if !isPortOpen(options.Host, options.Port) {
		return fmt.Errorf("port %d is not open on %s", options.Port, options.Host)
	}
	return nil
}
