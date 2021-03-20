
typedef struct {
} Application;

struct Application *NewApplication(const char *title, int width, int height, int resizable, int devtools, int fullscreen, int startHidden, int logLevel, int hideWindowOnClose) {
}
void SetMinWindowSize(struct Application* app, int minWidth, int minHeight) {
}
void SetMaxWindowSize(struct Application* app, int maxWidth, int maxHeight) {
}
void Run(struct Application* app, int argc, char **argv) {
}
void DestroyApplication(struct Application* app) {
}
void SetDebug(struct Application* app, int flag) {
}
void SetBindings(struct Application* app, const char *bindings) {
}
void ExecJS(struct Application* app, const char *script) {
}
void Hide(struct Application* app) {
}
void Show(struct Application* app) {
}
void Center(struct Application* app) {
}
void Maximise(struct Application* app) {
}
void Unmaximise(struct Application* app) {
}
void ToggleMaximise(struct Application* app) {
}
void Minimise(struct Application* app) {
}
void Unminimise(struct Application* app) {
}
void ToggleMinimise(struct Application* app) {
}
void SetColour(struct Application* app, int red, int green, int blue, int alpha) {
}
void SetSize(struct Application* app, int width, int height) {
}
void SetPosition(struct Application* app, int x, int y) {
}
void Quit(struct Application* app) {
}
void SetTitle(struct Application* app, const char *title) {
}
void Fullscreen(struct Application* app) {
}
void UnFullscreen(struct Application* app) {
}
void ToggleFullscreen(struct Application* app) {
}
void DisableFrame(struct Application* app) {
}
void OpenDialog(struct Application* app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int allowFiles, int allowDirs, int allowMultiple, int showHiddenFiles, int canCreateDirectories, int resolvesAliases, int treatPackagesAsDirectories) {
}
void SaveDialog(struct Application* app, char *callbackID, char *title, char *filters, char *defaultFilename, char *defaultDir, int showHiddenFiles, int canCreateDirectories, int treatPackagesAsDirectories) {
}
void MessageDialog(struct Application* app, char *callbackID, char *type, char *title, char *message, char *icon, char *button1, char *button2, char *button3, char *button4, char *defaultButton, char *cancelButton) {
}
void DarkModeEnabled(struct Application* app, char *callbackID) {
}
void SetApplicationMenu(struct Application* app, const char *applicationMenuJSON) {
}
void AddTrayMenu(struct Application* app, const char *menuTrayJSON) {
}
void SetTrayMenu(struct Application* app, const char *menuTrayJSON) {
}
void DeleteTrayMenuByID(struct Application* app, const char *id) {
}
void UpdateTrayMenuLabel(struct Application* app, const char* JSON) {
}
void AddContextMenu(struct Application* app, char *contextMenuJSON) {
}
void UpdateContextMenu(struct Application* app, char *contextMenuJSON) {
}