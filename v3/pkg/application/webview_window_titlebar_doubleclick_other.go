//go:build !darwin || ios || server

package application

func platformTitlebarDoubleClickPreference() string {
	return "None"
}
