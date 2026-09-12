package commands

import (
	"testing"

	"github.com/atterpac/refresh/process"
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

func TestEnsurePrimaryExitPolicy(t *testing.T) {
	t.Run("defaults primary process to shutdown", func(t *testing.T) {
		executes := []process.Execute{{Type: process.Primary}}
		ensurePrimaryExitPolicy(executes)
		assert.Equal(t, process.ExitPolicyShutdown, executes[0].ExitPolicy)
	})

	t.Run("preserves explicit policy", func(t *testing.T) {
		executes := []process.Execute{{Type: process.Primary, ExitPolicy: process.ExitPolicyIgnore}}
		ensurePrimaryExitPolicy(executes)
		assert.Equal(t, process.ExitPolicyIgnore, executes[0].ExitPolicy)
	})

	t.Run("ignores non-primary process", func(t *testing.T) {
		executes := []process.Execute{{Type: process.Background}}
		ensurePrimaryExitPolicy(executes)
		assert.Empty(t, executes[0].ExitPolicy)
	})
}
