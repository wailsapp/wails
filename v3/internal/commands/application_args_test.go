package commands

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplicationArgumentsCompatibilityAndLiteralVectors(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		compat, vector, want []string
		failure              bool
	}{
		{name: "inherit"},
		{name: "clear vector", vector: []string{}, want: []string{}},
		{name: "clear alias", compat: []string{""}, want: []string{}},
		{name: "quoted values", compat: []string{`--config-path 'profile with spaces.yaml' --label "quoted value" ''`}, want: []string{"--config-path", "profile with spaces.yaml", "--label", "quoted value", ""}},
		{name: "literal Windows path", compat: []string{`--config-path 'C:\Users\Example User\testing.yaml'`}, want: []string{"--config-path", `C:\Users\Example User\testing.yaml`}},
		{name: "no expansion", compat: []string{`'$HOME' '$(touch marker)' ';'`}, want: []string{"$HOME", "$(touch marker)", ";"}},
		{name: "literal vector", vector: []string{`C:\Users\Example User\testing.yaml`, "--help", "", "a'b\"c", "$HOME", ";"}, want: []string{`C:\Users\Example User\testing.yaml`, "--help", "", "a'b\"c", "$HOME", ";"}},
		{name: "mixed", compat: []string{""}, vector: []string{}, failure: true},
		{name: "repeated alias", compat: []string{"one", "two"}, failure: true},
		{name: "invalid quotes", compat: []string{`"unfinished`}, failure: true},
		{name: "NUL", vector: []string{"\x00"}, failure: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveApplicationArguments(tc.compat, tc.vector)
			if tc.failure {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
