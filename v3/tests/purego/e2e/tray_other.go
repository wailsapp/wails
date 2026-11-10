//go:build !linux

package main

// trayHost is Linux only: elsewhere the tray is not on D-Bus.
func trayHost(*probe) error { return nil }
