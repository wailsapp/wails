package commands

import (
	"net"
	"testing"
	"time"
)

func TestToolCheckPortConnections(t *testing.T) {
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	for _, useURL := range []bool{false, true} {
		name := "host_and_port"
		if useURL {
			name = "url"
		}
		t.Run(name, func(t *testing.T) {
			listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			address := listener.Addr().(*net.TCPAddr)
			options := &ToolCheckPortOptions{Host: address.IP.String(), Port: address.Port}
			if useURL {
				options.URL = "http://" + address.String()
			}
			if err := ToolCheckPort(options); err != nil {
				t.Fatal(err)
			}

			// All probes have completed, so drain the listener's queued connections.
			if err := listener.SetDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
				t.Fatal(err)
			}
			count := 0
			for {
				conn, err := listener.Accept()
				if err != nil {
					if timeout, ok := err.(net.Error); ok && timeout.Timeout() {
						break
					}
					t.Fatal(err)
				}
				count++
				conn.Close()
			}
			if count != 1 {
				t.Errorf("opened %d connections, want 1", count)
			}
		})
	}
}

func TestToolCheckPortValidation(t *testing.T) {
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	for _, options := range []ToolCheckPortOptions{
		{Host: "localhost"},
		{URL: "%"},
		{URL: "http://localhost:invalid"},
	} {
		if err := ToolCheckPort(&options); err == nil {
			t.Errorf("accepted invalid options: %+v", options)
		}
	}
}
