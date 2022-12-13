package build

import "testing"

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
