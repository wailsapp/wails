//go:build !windows

package commands

func manifestHookInterpreter(string) (string, error) { return "", nil }

func manifestHookCommand(script string) (string, []string, error) {
	return script, nil, nil
}
