package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	sqlitedriver "modernc.org/sqlite"
)

type delayedPrepareConnector struct {
	ready, resume chan struct{}
}

func (c *delayedPrepareConnector) Driver() driver.Driver { return &sqlitedriver.Driver{} }

func (c *delayedPrepareConnector) Connect(context.Context) (driver.Conn, error) {
	conn, err := c.Driver().Open(":memory:")
	if err != nil {
		return nil, err
	}
	return &delayedPrepareConn{Conn: conn, ready: c.ready, resume: c.resume}, nil
}

type delayedPrepareConn struct {
	driver.Conn
	ready, resume chan struct{}
}

func (c *delayedPrepareConn) Prepare(query string) (driver.Stmt, error) {
	close(c.ready)
	<-c.resume
	return c.Conn.Prepare(query)
}

func TestPrepareDuringConnectionShutdown(t *testing.T) {
	for _, operation := range []string{"close", "reopen"} {
		t.Run(operation, func(t *testing.T) {
			connector := &delayedPrepareConnector{ready: make(chan struct{}), resume: make(chan struct{})}
			service := New()
			service.conn = sql.OpenDB(connector)
			service.stmts = make(map[uint64]struct{})
			defer service.Close()
			type result struct {
				stmt     *Stmt
				err      error
				panicked any
			}
			done := make(chan result, 1)
			go func() {
				defer func() {
					if p := recover(); p != nil {
						done <- result{panicked: p}
					}
				}()
				stmt, err := service.Prepare("SELECT 1")
				done <- result{stmt: stmt, err: err}
			}()
			select {
			case <-connector.ready:
			case <-time.After(5 * time.Second):
				t.Fatal("prepare did not reach driver")
			}
			var err error
			if operation == "close" {
				err = service.Close()
			} else {
				err = service.Open()
			}
			close(connector.resume)
			if err != nil {
				t.Fatal(err)
			}
			select {
			case result := <-done:
				if result.panicked != nil {
					t.Fatalf("prepare panicked after shutdown: %v", result.panicked)
				}
				if result.stmt != nil {
					defer result.stmt.Close()
				}
				if result.err == nil || result.stmt != nil {
					t.Errorf("stale prepare returned statement=%v, error=%v", result.stmt, result.err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("prepare did not finish after shutdown")
			}
			if len(service.stmts) != 0 {
				t.Errorf("registered %d stale statements", len(service.stmts))
			}
		})
	}
}
