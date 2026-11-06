package commands

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"
)

func manifestDevInput(ctx context.Context) <-chan string {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil
	}
	input := make(chan string, 1)
	go func() {
		defer close(input)
		err := readManifestDevInput(ctx, os.Stdin, func(line string) bool {
			select {
			case input <- line:
				return true
			case <-ctx.Done():
				return false
			}
		})
		if err != nil && ctx.Err() == nil {
			fmt.Fprintln(os.Stderr, "Wake keyboard controls disabled:", err)
		}
	}()
	return input
}
