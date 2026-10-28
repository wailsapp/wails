package sqlite

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestQueryReportsRowIterationError(t *testing.T) {
	for _, prepared := range []bool{false, true} {
		name := "direct"
		if prepared {
			name = "prepared"
		}
		t.Run(name, func(t *testing.T) {
			service := New()
			if err := service.Open(); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = service.Close() })
			const query = "SELECT 1 AS value UNION ALL SELECT abs(-9223372036854775808)"
			var rows Rows
			var err error
			if prepared {
				stmt, prepareErr := service.Prepare(query)
				if prepareErr != nil {
					t.Fatal(prepareErr)
				}
				defer stmt.Close()
				rows, err = service.QueryPrepared(context.Background(), stmt)
			} else {
				rows, err = service.Query(query)
			}
			if err == nil || !strings.Contains(err.Error(), "integer overflow") {
				t.Errorf("query error=%v, want integer overflow", err)
			}
			if len(rows) != 0 {
				t.Errorf("failed query returned partial rows as success: %v", rows)
			}
		})
	}
}

func TestParseRowsCancellation(t *testing.T) {
	service := New()
	if err := service.Open(); err != nil {
		t.Fatal(err)
	}
	defer service.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rows, err := service.conn.QueryContext(ctx, "SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if result, err := parseRows(ctx, rows); err != nil || len(result) != 0 {
		t.Errorf("cancelled query result=%v, error=%v", result, err)
	}
}

func TestPrepareConcurrent(t *testing.T) {
	s := New()
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	const workers = 32
	const rounds = 100

	ids := make(chan uint64, workers*rounds)
	errs := make(chan error, workers)
	start := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < rounds; j++ {
				stmt, err := s.Prepare("SELECT 1")
				if err != nil {
					errs <- err
					return
				}
				ids <- stmt.id
				if err := stmt.Close(); err != nil {
					errs <- err
					return
				}
			}
		}()
	}
	close(start)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("Prepare deadlocked under concurrent use")
	}

	close(errs)
	for err := range errs {
		t.Error(err)
	}

	close(ids)
	seen := make(map[uint64]bool)
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate prepared statement id %d", id)
		}
		seen[id] = true
	}
}
