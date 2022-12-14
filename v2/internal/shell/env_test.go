package shell

import "testing"

func TestUpdateEnv(t *testing.T) {

	env := []string{"one=1", "two=a=b", "three="}
	newEnv := UpsertEnv(env, "two", func(v string) string {
		return v + "+added"
	})
	newEnv = UpsertEnv(newEnv, "newVar", func(v string) string {
		return "added"
	})
	newEnv = UpsertEnv(newEnv, "three", func(v string) string {
		return "3"
	})
	newEnv = UpsertEnv(newEnv, "GOARCH", func(v string) string {
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

func TestSetEnv(t *testing.T) {
	env := []string{"one=1", "two=a=b", "three="}
	newEnv := SetEnv(env, "two", "set")
	newEnv = SetEnv(newEnv, "newVar", "added")

	if len(newEnv) != 4 {
		t.Errorf("expected: 4, got: %d", len(newEnv))
	}
	if newEnv[1] != "two=set" {
		t.Errorf("expected: \"two=set\", got: %q", newEnv[1])
	}
	if newEnv[3] != "newVar=added" {
		t.Errorf("expected: \"newVar=added\", got: %q", newEnv[3])
	}
}

func TestRemoveEnv(t *testing.T) {
	env := []string{"one=1", "two=a=b", "three=3"}
	newEnv := RemoveEnv(env, "two")

	if len(newEnv) != 2 {
		t.Errorf("expected: 2, got: %d", len(newEnv))
	}
	if newEnv[0] != "one=1" {
		t.Errorf("expected: \"one=1\", got: %q", newEnv[1])
	}
	if newEnv[1] != "three=3" {
		t.Errorf("expected: \"three=3\", got: %q", newEnv[3])
	}
}
