package migration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportIsVersionedEphemeralJSON(t *testing.T) {
	want := Report{Version: 1, CompletedBy: "v3", Complete: false, Sources: map[string]string{"Taskfile.yml": "digest"}}
	data, err := json.Marshal(want)
	require.NoError(t, err)
	var got Report
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, want, got)
}
