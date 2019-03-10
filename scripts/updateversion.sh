#!/usr/bin/env bash
TAG=$(git describe --abbrev=0 --tags)
cat << EOF > cmd/version.go
package cmd

// Version - Wails version
const Version = "${TAG}"
EOF