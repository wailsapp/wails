package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureIgnored(t *testing.T) {
	t.Run("adds pattern when not present", func(t *testing.T) {
		list := []string{".gitignore", ".DS_Store"}
		ensureIgnored(&list, "*_test.go")
		assert.Contains(t, list, "*_test.go")
		assert.Len(t, list, 3)
	})

	t.Run("does not duplicate pattern when already present", func(t *testing.T) {
		list := []string{".gitignore", "*_test.go"}
		ensureIgnored(&list, "*_test.go")
		assert.Contains(t, list, "*_test.go")
		assert.Len(t, list, 2)
	})

	t.Run("adds to empty list", func(t *testing.T) {
		var list []string
		ensureIgnored(&list, "*_test.go")
		assert.Contains(t, list, "*_test.go")
		assert.Len(t, list, 1)
	})
}
