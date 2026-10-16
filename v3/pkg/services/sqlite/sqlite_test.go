package sqlite

import (
	"sync"
	"testing"
	"time"
)

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
