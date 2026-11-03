#!/usr/bin/env bash
# run-linux.sh LANE [TAGS]: build the e2e probe with TAGS (default purego) and
# run it under a fresh Xvfb and D-Bus session in the named configuration.
# Run from v3/. The lanes are the Linux set-ups users have that change how
# input reaches the app:
#
#   default    plain X server, no window manager, US keyboard
#   scale2     GDK_SCALE=2 (HiDPI): pointer and page coordinates differ
#   layout-de  German QWERTZ: Y/Z swapped, punctuation on other keys
#   layout-fr  French AZERTY: digits need Shift
#   layout-ru  Russian first, US second: shortcuts arrive with Cyrillic keysyms
#   numlock    Num Lock on: an extra modifier on every key event
#   capslock   Caps Lock on: shortcuts with the lock modifier set
#   wm         a window manager (openbox) owns focus and decorations
#   a11y       the accessibility bus is on, as on a desktop session
#
# Lanes that change what a typed letter produces skip the scenarios that type
# text. Set GOCOVERDIR to collect coverage, E2E_EVIDENCE_DIR for screenshots.
set -euo pipefail

if [ "${1-}" = "--inside" ]; then
	# Inside the X and D-Bus session: set the keyboard and window manager up,
	# then run the probe.
	if [ -n "${E2E_LAYOUT-}" ]; then
		setxkbmap -layout "$E2E_LAYOUT"
	fi
	if [ -n "${E2E_WM-}" ]; then
		"$E2E_WM" &
		# Kept: the window manager has no readiness signal the probe can wait
		# on, and windows mapped before it runs are not managed.
		sleep 2
	fi
	exec "$E2E_PROBE"
fi

lane=${1:?usage: run-linux.sh LANE [TAGS]}
tags=${2-purego}
screen=1280x800x24
# GSK_RENDERER may be set by the caller: cairo is the default as Xvfb has no
# GPU, but a lane can ask for GTK's GL renderers (on Mesa's llvmpipe).
export GDK_BACKEND=x11 GSK_RENDERER=${GSK_RENDERER:-cairo} LIBGL_ALWAYS_SOFTWARE=1 WEBKIT_DISABLE_DMABUF_RENDERER=1
export GTK_A11Y=none NO_AT_BRIDGE=1
unset E2E_LAYOUT E2E_WM
case "$lane" in
default) ;;
scale2)
	export GDK_SCALE=2
	screen=2560x1600x24
	;;
layout-de) export E2E_LAYOUT=de ;;
layout-fr) export E2E_LAYOUT=fr ;;
layout-ru) export E2E_LAYOUT=ru,us E2E_SKIP=type-text,open-file ;;
numlock) export E2E_LOCKS=numlock ;;
capslock) export E2E_LOCKS=capslock E2E_SKIP=type-text,open-file ;;
wm) export E2E_WM=openbox ;;
a11y) unset GTK_A11Y NO_AT_BRIDGE ;;
*)
	echo "run-linux.sh: unknown lane $lane" >&2
	exit 2
	;;
esac

case "$tags" in
*purego*) export CGO_ENABLED=0 ;;
*) export CGO_ENABLED=1 ;;
esac
probe=$(mktemp -d)/e2e
cover=()
if [ -n "${GOCOVERDIR-}" ]; then
	mkdir -p "$GOCOVERDIR"
	cover=(-cover -covermode=atomic)
fi
# CI containers check the repository out as another user, so git refuses to
# report VCS status: the probe needs no VCS stamp.
go build -buildvcs=false "${cover[@]}" -tags "$tags" -o "$probe" ./tests/purego/e2e
export E2E_PROBE=$probe
echo "e2e lane $lane, tags '$tags', screen $screen"
exec xvfb-run -a -s "-screen 0 $screen" dbus-run-session -- "$0" --inside
