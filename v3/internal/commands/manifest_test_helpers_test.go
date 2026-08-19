package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(data)
}
