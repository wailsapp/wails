package appimage

import (
	"fmt"
	"strings"
)

// Launcher renders the Wails-owned AppRun for the primary executable.
func Launcher(binary string) []byte {
	return []byte(fmt.Sprintf(launcher, "'"+strings.ReplaceAll(binary, "'", "'\"'\"'")+"'"))
}

// launcher starts the application from $APPDIR/usr/lib, where the path relocateWebKitHelpers writes into the WebKit
// library resolves, and sets the environment AppImageKit's AppRun would. AppImageKit's AppRun changes directory to
// $APPDIR/usr instead, which cannot be made to work with WebKit's sandbox.
// WEBKIT_INJECTED_BUNDLE_PATH, which release builds honour, points WebKit at the
// bundled injected bundle: the path compiled into the library is the build
// machine's, missing on other distros and a different WebKit build on the same one.
const launcher = `#!/bin/sh
set -e

APPDIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
OWD="${OWD:-$PWD}"
export OWD
usr="$APPDIR/usr"
export PATH="$usr/bin:$PATH"
export LD_LIBRARY_PATH="$usr/lib:$usr/lib/x86_64-linux-gnu:$usr/lib/aarch64-linux-gnu:$usr/lib64${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
export XDG_DATA_DIRS="$usr/share:${XDG_DATA_DIRS:-/usr/local/share:/usr/share}"
export GSETTINGS_SCHEMA_DIR="$usr/share/glib-2.0/schemas"
for bundle in "$usr"/lib*/webkit*gtk-*/injected-bundle/*.so "$usr"/lib*/*/webkit*gtk-*/injected-bundle/*.so; do
	if [ -e "$bundle" ]; then
		export WEBKIT_INJECTED_BUNDLE_PATH="${bundle%%/*}/"
		break
	fi
done
cd "$usr/lib"
exec "$usr/bin/"%s "$@"
`
