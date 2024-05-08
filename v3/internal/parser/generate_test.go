package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func ExampleGenerator() {
	generator := NewGenerator(&flags.GenerateBindingsOptions{
		UseBundledRuntime: true,
		TS:                true,
	}, FileCreatorFunc(dummyBufferCreator))

	err := generator.Generate("github.com/wailsapp/wails/v3/internal/parser/testdata/complex_json")
	if err != nil {
		pterm.Error.Println(err)
	}
	// Output: abc
}

var mu sync.Mutex

type dummyBuffer struct {
	bytes.Buffer
	path string
}

func (buffer *dummyBuffer) Close() error {
	mu.Lock()
	defer mu.Unlock()

	_, err := fmt.Fprintln(os.Stdout, "// ====", buffer.path, "====")
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(buffer.Bytes())
	return err
}

func dummyBufferCreator(path string) (io.WriteCloser, error) {
	return &dummyBuffer{
		path: path,
	}, nil
}
