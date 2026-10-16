package application

import "testing"

func TestParseMacTitlebarDoubleClickAction(t *testing.T) {
	tests := []struct {
		name       string
		preference string
		want       titlebarDoubleClickAction
	}{
		{name: "missing preference uses native default", want: titlebarDoubleClickToggleMaximise},
		{name: "maximize", preference: "Maximize", want: titlebarDoubleClickToggleMaximise},
		{name: "minimize", preference: "Minimize", want: titlebarDoubleClickMinimise},
		{name: "none", preference: "None", want: titlebarDoubleClickNone},
		{name: "unknown preference is ignored", preference: "FutureValue", want: titlebarDoubleClickNone},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := parseMacTitlebarDoubleClickAction(test.preference); got != test.want {
				t.Fatalf("parseMacTitlebarDoubleClickAction(%q) = %d, want %d", test.preference, got, test.want)
			}
		})
	}
}
