package features

import "os"

// WakeEnv opts into the experimental HCL build-system CLI.
const WakeEnv = "WAILS_EXP_USE_WAKE"

// WakeEnabled checks presence, not truthiness: even an empty value opts in.
func WakeEnabled() bool {
	_, present := os.LookupEnv(WakeEnv)
	return present
}
