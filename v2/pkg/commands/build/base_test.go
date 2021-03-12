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

	if len(newEnv) != 4 {
		t.Errorf("expected: 4, got: %d", len(newEnv))
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

}
