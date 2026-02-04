package doctor

import "testing"

func TestRun(t *testing.T) {
	err := Run(false)
	if err != nil {
		t.Errorf("TestRun failed: %v", err)
	}
}

func TestRunJSON(t *testing.T) {
	err := Run(true)
	if err != nil {
		t.Errorf("TestRunJSON failed: %v", err)
	}
}
