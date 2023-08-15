package commands

import (
	"fmt"
	"net"
	"time"
)

type ToolCheckPortOptions struct {
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
	if options.Port == 0 {
		return fmt.Errorf("please use the -p flag to specify a port")
	}
	if !isPortOpen(options.Host, options.Port) {
		return fmt.Errorf("port %d is not open on %s", options.Port, options.Host)
	}
	return nil
}
