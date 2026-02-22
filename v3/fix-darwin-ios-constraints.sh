#!/bin/bash

# Fix build constraints for Darwin files to exclude iOS
echo "Fixing build constraints to exclude iOS from Darwin builds..."

# List of files that need updating
files=(
    "pkg/application/webview_window_darwin.m"
    "pkg/application/webview_window_darwin.h"
    "pkg/application/webview_window_darwin.go"
    "pkg/application/webview_window_darwin_drag.m"
    "pkg/application/webview_window_darwin_drag.h"
    "pkg/application/webview_window_close_darwin.go"
    "pkg/application/systemtray_darwin.m"
    "pkg/application/systemtray_darwin.h"
    "pkg/application/systemtray_darwin.go"
    "pkg/application/single_instance_darwin.go"
    "pkg/application/screen_darwin.go"
    "pkg/application/menuitem_selectors_darwin.go"
    "pkg/application/menuitem_darwin.m"
    "pkg/application/menuitem_darwin.go"
    "pkg/application/menu_darwin.go"
    "pkg/application/mainthread_darwin.go"
    "pkg/application/keys_darwin.go"
    "pkg/application/events_common_darwin.go"
    "pkg/application/dialogs_darwin_delegate.m"
    "pkg/application/dialogs_darwin_delegate.h"
    "pkg/application/dialogs_darwin.go"
    "pkg/application/clipboard_darwin.go"
    "pkg/application/application_darwin_delegate.m"
    "pkg/application/application_darwin_delegate.h"
    "pkg/application/application_darwin.h"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        # Check if file has the build constraint
        if grep -q "^//go:build darwin$" "$file"; then
            echo "Updating $file"
            sed -i '' 's|^//go:build darwin$|//go:build darwin \&\& !ios|' "$file"
        fi
    fi
done

# Also check for other darwin-specific files
echo "Checking for other darwin-specific build constraints..."
find . -name "*_darwin*.go" -o -name "*_darwin*.m" -o -name "*_darwin*.h" | while read -r file; do
    if grep -q "^//go:build darwin$" "$file"; then
        echo "Also updating: $file"
        sed -i '' 's|^//go:build darwin$|//go:build darwin \&\& !ios|' "$file"
    fi
done

echo "Done! Build constraints updated."