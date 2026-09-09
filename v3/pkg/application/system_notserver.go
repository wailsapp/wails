//go:build !server

package application

// isServerBuild is true when the application is built with the "server" tag.
// In every non-server build it is false.
const isServerBuild = false
