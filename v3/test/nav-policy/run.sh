#!/bin/sh
# Builds the check for each Linux backend named on the command line (default: gtk4) and runs it on a headless mutter
# under its own D-Bus session, so nothing appears on the desktop. Backends: gtk4, gtk3, purego, purego-gtk3.
set -eu
cd "$(dirname "$0")/../.."
if [ $# -eq 0 ]; then set -- gtk4; fi
for backend in "$@"; do
  case $backend in
    gtk4) tags=""; cgo=1 ;;
    gtk3) tags="gtk3"; cgo=1 ;;
    purego) tags="purego"; cgo=0 ;;
    purego-gtk3) tags="purego,gtk3"; cgo=0 ;;
    *) echo "unknown backend: $backend" >&2; exit 2 ;;
  esac
  bin=$(mktemp -d /var/tmp/navpolicy.XXXXXX)/nav-policy
  CGO_ENABLED=$cgo go build -tags "$tags" -o "$bin" ./test/nav-policy
  echo "== $backend"
  env -u DISPLAY -u WAYLAND_DISPLAY dbus-run-session -- timeout 60 mutter --wayland --headless --virtual-monitor 1280x720 \
    --wayland-display nav-policy-$$ -- "$bin"
  rm -rf "$(dirname "$bin")"
done
