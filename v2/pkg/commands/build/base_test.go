package build

import "testing"

func TestUpdateEnv(t *testing.T) {

	env := []string{"one=1", "two=a=b", "three="}
	newEnv := upsertEnv(env, "two", func(v string) string {
		return v + "+added"
	})
	newEnv = upsertEnv(newEnv, "newVar", func(v string) string {
		return "added"
	})
	newEnv = upsertEnv(newEnv, "three", func(v string) string {
		return "3"
	})
	newEnv = upsertEnv(newEnv, "GOARCH", func(v string) string {
		return "amd64"
	})

	if len(newEnv) != 5 {
		t.Errorf("expected: 5, got: %d", len(newEnv))
	}
	if newEnv[1] != "two=a=b+added" {
		t.Errorf("expected: \"two=a=b+added\", got: %q", newEnv[1])
	}
	if newEnv[2] != "three=3" {
		t.Errorf("expected: \"three=3\", got: %q", newEnv[2])
	}
	if newEnv[3] != "newVar=added" {
		t.Errorf("expected: \"newVar=added\", got: %q", newEnv[3])
	}
	if newEnv[4] != "GOARCH=amd64" {
		t.Errorf("expected: \"newVar=added\", got: %q", newEnv[4])
	}

}

func Test_commandPrettifier(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "empty",
			input: []string{},
			want:  "",
		},
		{
			name:  "one arg",
			input: []string{"one"},
			want:  "one",
		},
		{
			name:  "args where one has spaces",
			input: []string{"one", "two three"},
			want:  `one "two three"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandPrettifier(tt.input); got != tt.want {
				t.Errorf("commandPrettifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
