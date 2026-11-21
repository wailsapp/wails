#!/bin/sh
# ci-deps.sh [cgo]: install what run-linux.sh needs on this distro: the GTK and
# WebKitGTK runtime libraries of both stacks, Xvfb, D-Bus, XTEST, the keyboard
# and window-manager tools the lanes use, and ImageMagick for the evidence
# screenshots. With "cgo", also the headers and compiler the cgo reference
# backend builds with. Runs as root (CI containers). Minimal images may have
# no fonts at all, and GTK then crashes drawing text (Pango asserts), so a
# font is installed explicitly.
set -eu
cgo=${1-}
. /etc/os-release
case "$ID" in
fedora)
	pkgs="git tar gzip findutils gtk4 webkitgtk6.0 gtk3 webkit2gtk4.1 libXtst xorg-x11-server-Xvfb xorg-x11-xauth
		dbus-daemon at-spi2-core ImageMagick setxkbmap xkbcomp openbox"
	[ "$cgo" = cgo ] && pkgs="$pkgs gcc pkgconf-pkg-config gtk4-devel webkitgtk6.0-devel gtk3-devel webkit2gtk4.1-devel"
	dnf install -y $pkgs
	;;
ubuntu | debian)
	export DEBIAN_FRONTEND=noninteractive
	pkgs="git ca-certificates curl libgtk-3-0 libwebkit2gtk-4.1-0 libxtst6 xvfb xauth dbus at-spi2-core
		imagemagick x11-xkb-utils openbox"
	apt-get update -q
	# The GTK4 stack, where the distro has WebKitGTK 6.0 (not Ubuntu 22.04).
	if apt-cache show libwebkitgtk-6.0-4 >/dev/null 2>&1; then
		pkgs="$pkgs libgtk-4-1 libwebkitgtk-6.0-4"
		[ "$cgo" = cgo ] && pkgs="$pkgs libgtk-4-dev libwebkitgtk-6.0-dev"
	fi
	[ "$cgo" = cgo ] && pkgs="$pkgs gcc pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev"
	apt-get install -y -q --no-install-recommends $pkgs
	;;
arch)
	pkgs="git gtk4 webkitgtk-6.0 gtk3 webkit2gtk-4.1 libxtst xorg-server-xvfb xorg-xauth dbus at-spi2-core
		imagemagick xorg-setxkbmap xorg-xkbcomp openbox ttf-dejavu"
	[ "$cgo" = cgo ] && pkgs="$pkgs gcc pkgconf"
	pacman -Syu --noconfirm --needed $pkgs
	;;
opensuse-tumbleweed)
	pkgs="git-core tar gzip libgtk-4-1 libwebkitgtk-6_0-4 libgtk-3-0 libwebkit2gtk-4_1-0 libXtst6 xvfb-run
		xorg-x11-server-Xvfb xauth dbus-1 dbus-1-daemon at-spi2-core ImageMagick setxkbmap xkbcomp openbox dejavu-fonts"
	zypper --non-interactive install --no-recommends $pkgs
	;;
*)
	echo "ci-deps.sh: no package list for $ID" >&2
	exit 2
	;;
esac
