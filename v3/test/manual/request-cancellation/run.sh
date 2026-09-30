#!/usr/bin/env bash
# Requires a Linux graphical session (or xvfb-run) and GTK/WebKit development packages.
# Pass -tags gtk3 to exercise the legacy backend. Other go build flags are forwarded.
set -euo pipefail
cd "$(dirname "$0")/../../.."
probe_dir=$(mktemp -d)
trap 'rm -rf "$probe_dir"' EXIT
go build "$@" -o "$probe_dir/probe" ./test/manual/request-cancellation
for scenario in abort stream normal parallel navigate close; do
  timeout 20s "$probe_dir/probe" "$scenario"
done
