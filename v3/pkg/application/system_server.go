//go:build server

package application

// isServerBuild is true when the application is built with the "server" tag
// (headless HTTP server mode, no native GUI).
const isServerBuild = true
