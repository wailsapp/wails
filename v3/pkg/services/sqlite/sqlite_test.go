package sqlite

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStatementHandleParsing(t *testing.T) {
	service := New()
	if err := service.Open(); err != nil {
		t.Fatal(err)
	}
	defer service.Close()
	stmt, err := service.Prepare("SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	encoded, err := stmt.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	var valid Stmt
	if err := valid.UnmarshalText(encoded); err != nil {
		t.Fatal(err)
	}
	if rows, err := service.QueryPrepared(context.Background(), &valid); err != nil || len(rows) != 1 {
		t.Fatalf("round-trip handle unusable: %v, %v", rows, err)
	}
	for _, text := range []string{
		"", "not-hex", strings.Repeat("f", 17),
		string(encoded) + "junk", string(encoded) + " more",
		" " + string(encoded), string(encoded) + "\n",
	} {
		t.Run(text, func(t *testing.T) {
			var decoded Stmt
			if err := decoded.UnmarshalText([]byte(text)); err == nil {
				t.Errorf("malformed handle %q resolved to statement %d", text, decoded.id)
			}
		})
	}
}

func TestConnectionCleanupDropsStatementHandles(t *testing.T) {
	for _, operation := range []string{"close", "reopen"} {
		t.Run(operation, func(t *testing.T) {
			service := New()
			if err := service.Open(); err != nil {
				t.Fatal(err)
			}
			defer service.Close()
			other := New()
			if err := other.Open(); err != nil {
				t.Fatal(err)
			}
			defer other.Close()
			live, err := other.Prepare("SELECT 1")
			if err != nil {
				t.Fatal(err)
			}
			defer live.Close()
			var prepared []*Stmt
			for range 32 {
				stmt, err := service.Prepare("SELECT 1")
				if err != nil {
					t.Fatal(err)
				}
				prepared = append(prepared, stmt)
			}
			if operation == "close" {
				err = service.Close()
			} else {
				err = service.Open()
			}
			if err != nil {
				t.Fatal(err)
			}
			retained := 0
			for _, stmt := range prepared {
				if _, ok := stmts.Load(stmt.id); ok {
					retained++
				}
			}
			if retained != 0 {
				t.Errorf("cleanup retained %d of 32 statement handles", retained)
			}
			encoded, err := prepared[0].MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			var decoded Stmt
			if err := decoded.UnmarshalText(encoded); err != nil {
				t.Fatal(err)
			}
			if decoded.sqlStmt != nil {
				t.Error("closed statement handle still resolves")
			}
			if _, ok := stmts.Load(live.id); !ok {
				t.Error("cleanup removed another service's live handle")
			}
			if rows, err := other.QueryPrepared(context.Background(), live); err != nil || len(rows) != 1 {
				t.Errorf("other service statement unusable: %v, %v", rows, err)
			}
		})
	}
}

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
