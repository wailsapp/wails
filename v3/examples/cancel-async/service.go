package main

import (
	"context"
	"time"
)

type Service struct {
}

// LongRunning - A long-running operation of specified duration.
func (*Service) LongRunning(ctx context.Context, milliseconds int) error {
	select {
	case <-time.After(time.Duration(milliseconds) * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
