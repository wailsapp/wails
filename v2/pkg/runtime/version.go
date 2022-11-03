package runtime

import "context"

const Version = "2.0.0"

const (
	GtkCompiledVersion string = "gtk3-compiled"
	GtkRuntimeVersion  string = "gtk3-runtime"

	Webkit2GtkCompiledVersion string = "webkit2gtk-compiled"
	Webkit2GtkRuntimeVersion  string = "webkit2gtk-runtime"

	Webview2          string = "webview2"
)

func GetNativeVersions(ctx context.Context) map[string]string {
    versions := make(map[string]string)

    frontend := getFrontend(ctx)
    frontend.PopulateVersionMap(versions)

    return versions
}