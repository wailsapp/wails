package sqlite

import (
	"context"
	"testing"
	"time"
)

func TestMemoryDatabaseUsesSameConnection(t *testing.T) {
	service := New()
	if err := service.Open(); err != nil {
		t.Fatal(err)
	}
	defer service.Close()
	if err := service.Execute("CREATE TABLE example (value INTEGER); INSERT INTO example VALUES (42)"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	held, err := service.conn.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()
	type result struct {
		rows Rows
		err  error
	}
	done := make(chan result, 1)
	go func() {
		rows, err := service.QueryContext(ctx, "SELECT value FROM example")
		done <- result{rows, err}
	}()
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	waiting := false
	for !waiting {
		select {
		case result := <-done:
			t.Fatalf("query opened a separate memory connection: rows=%v, error=%v", result.rows, result.err)
		case <-ticker.C:
			waiting = service.conn.Stats().WaitCount > 0
		case <-ctx.Done():
			t.Fatal("query did not wait for the existing memory connection")
		}
	}
	if err := held.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case result := <-done:
		if result.err != nil || len(result.rows) != 1 || result.rows[0]["value"] != int64(42) {
			t.Errorf("query lost in-memory contents: rows=%v, error=%v", result.rows, result.err)
		}
	case <-ctx.Done():
		t.Fatal("query did not resume after the connection was returned")
	}
}
