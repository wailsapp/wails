package usernssandbox

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func fakeEnv(files map[string]string, env map[string]string, bwrap bool, runErr error, out string, ran *int) Env {
	return Env{
		ReadFile: func(p string) ([]byte, error) {
			if v, ok := files[p]; ok {
				return []byte(v), nil
			}
			return nil, os.ErrNotExist
		},
		Getenv: func(k string) string { return env[k] },
		LookPath: func(string) (string, error) {
			if bwrap {
				return "/usr/bin/bwrap", nil
			}
			return "", errors.New("not found")
		},
		Run: func(string, ...string) ([]byte, error) {
			if ran != nil {
				*ran++
			}
			return []byte(out), runErr
		},
	}
}

func TestCheck(t *testing.T) {
	restricted := map[string]string{procRestrictUserns: "1\n"}
	tests := []struct {
		name    string
		files   map[string]string
		env     map[string]string
		bwrap   bool
		runErr  error
		probe   bool
		want    Status
		wantRun int
		disable bool
	}{
		{"no restrictions skips probe", map[string]string{procRestrictUserns: "0"}, nil, true, nil, false, OK, 0, false},
		{"missing sysctls skips probe", nil, nil, true, errors.New("x"), false, OK, 0, false},
		{"restricted and bwrap fails", restricted, nil, true, errors.New("exit 1"), false, Blocked, 1, false},
		{"restricted but bwrap works (profile allows)", restricted, nil, true, nil, false, OK, 1, false},
		{"restricted and no bwrap", restricted, nil, false, nil, false, Unknown, 0, false},
		{"userns_clone off", map[string]string{procUsernsClone: "0"}, nil, true, errors.New("x"), false, Blocked, 1, false},
		{"max_user_namespaces zero", map[string]string{procMaxUserNS: "0"}, nil, true, errors.New("x"), false, Blocked, 1, false},
		{"doctor probes even when unrestricted", nil, nil, true, errors.New("x"), true, Blocked, 1, false},
		{"doctor, no bwrap, no restrictions", nil, nil, false, nil, true, OK, 0, false},
		{"sandbox disabled by env", restricted, map[string]string{DisableEnvVar: "1"}, true, errors.New("x"), true, OK, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ran := 0
			res := Check(fakeEnv(tt.files, tt.env, tt.bwrap, tt.runErr, "bwrap: setting up uid map: Permission denied", &ran), tt.probe)
			if res.Status != tt.want {
				t.Fatalf("status = %v, want %v", res.Status, tt.want)
			}
			if ran != tt.wantRun {
				t.Fatalf("probe ran %d times, want %d", ran, tt.wantRun)
			}
			if res.SandboxDisabled != tt.disable {
				t.Fatalf("SandboxDisabled = %v", res.SandboxDisabled)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	res := Check(fakeEnv(map[string]string{procRestrictUserns: "1"}, nil, true, errors.New("x"), "bwrap: setting up uid map: Permission denied", nil), false)
	m := res.Message()
	for _, s := range []string{"uid map", "apparmor_restrict_unprivileged_userns"} {
		if !strings.Contains(m, s) {
			t.Errorf("message %q missing %q", m, s)
		}
	}
	if (Result{}).Message() != "" {
		t.Error("OK result should have empty message")
	}
}
