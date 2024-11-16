package buildinfo

import (
	"testing"
)

func TestGet(t *testing.T) {
	result, err := Get()
	if err != nil {
		t.Error(err)
	}
	_ = result
}
