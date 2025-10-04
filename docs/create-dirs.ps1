$dirs = @(
    "src/content/docs/quick-start",
    "src/content/docs/concepts",
    "src/content/docs/features/windows",
    "src/content/docs/features/menus",
    "src/content/docs/features/bindings",
    "src/content/docs/features/events",
    "src/content/docs/features/dialogs",
    "src/content/docs/features/platform",
    "src/content/docs/guides/dev",
    "src/content/docs/guides/build",
    "src/content/docs/guides/distribution",
    "src/content/docs/guides/patterns",
    "src/content/docs/guides/advanced",
    "src/content/docs/reference/application",
    "src/content/docs/reference/window",
    "src/content/docs/reference/menu",
    "src/content/docs/reference/events",
    "src/content/docs/reference/dialogs",
    "src/content/docs/reference/runtime",
    "src/content/docs/reference/cli",
    "src/content/docs/contributing/architecture",
    "src/content/docs/contributing/codebase",
    "src/content/docs/contributing/workflows",
    "src/content/docs/migration"
)

foreach ($dir in $dirs) {
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    Write-Host "Created: $dir"
}

Write-Host "`nAll directories created successfully!"
