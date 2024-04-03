package doctor

import "testing"

func TestRun(t *testing.T) {
	err := Run()
	if err != nil {
		t.Errorf("TestRun failed: %v", err)
	}
}
