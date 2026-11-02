package features

import (
	"os"
	"testing"
)

func TestWakeUsesEnvironmentPresence(t *testing.T) {
	t.Setenv(WakeEnv, "restore")
	if err := os.Unsetenv(WakeEnv); err != nil {
		t.Fatal(err)
	}
	if WakeEnabled() {
		t.Fatal("unset experiment must be disabled")
	}
	for _, value := range []string{"", "false", "0", "true", "anything"} {
		t.Setenv(WakeEnv, value)
		if !WakeEnabled() {
			t.Fatalf("set value %q must enable the experiment", value)
		}
	}
}
