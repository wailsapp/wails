package shell

import (
	"fmt"
	"strings"
)

func UpsertEnv(env []string, key string, update func(v string) string) []string {
	newEnv := make([]string, len(env), len(env)+1)
	found := false
	for i := range env {
		if strings.HasPrefix(env[i], key+"=") {
			eqIndex := strings.Index(env[i], "=")
			val := env[i][eqIndex+1:]
			newEnv[i] = fmt.Sprintf("%s=%v", key, update(val))
			found = true
			continue
		}
		newEnv[i] = env[i]
	}
	if !found {
		newEnv = append(newEnv, fmt.Sprintf("%s=%v", key, update("")))
	}
	return newEnv
}

func RemoveEnv(env []string, key string) []string {
	newEnv := make([]string, 0, len(env))
	for _, e := range env {
		if strings.HasPrefix(e, key+"=") {
			continue
		}
		newEnv = append(newEnv, e)
	}
	return newEnv
}

func SetEnv(env []string, key string, value string) []string {
	return UpsertEnv(env, key, func(_ string) string { return value })
}
