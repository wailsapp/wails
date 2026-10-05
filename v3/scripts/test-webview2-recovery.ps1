# Run on an interactive Windows desktop with Go and WebView2 installed.
# A temporary Go overlay injects failures without changing repository files.
$ErrorActionPreference = 'Stop'
$root = Split-Path $PSScriptRoot -Parent
$edge = Join-Path $root 'internal/webview2/pkg/edge'
$temporary = Join-Path ([IO.Path]::GetTempPath()) ('wails-recovery-' + [guid]::NewGuid())
$utf8 = New-Object Text.UTF8Encoding($false)

function Add-FaultInjection([string]$source, [string]$signature, [string]$marker, [string]$injection) {
    $start = $source.IndexOf($signature, [StringComparison]::Ordinal)
    if ($start -lt 0) { throw "Missing function: $signature" }
    $position = $source.IndexOf($marker, $start, [StringComparison]::Ordinal)
    $nextFunction = $source.IndexOf("`nfunc ", $start + $signature.Length, [StringComparison]::Ordinal)
    if ($position -lt 0 -or ($nextFunction -ge 0 -and $position -ge $nextFunction)) {
        throw "Missing injection point in $signature"
    }
    return $source.Insert($position, $injection)
}

$helper = @'

// Test-only injection, present only in the temporary Go overlay.
func recoveryTestHRESULT(kind string, actual uintptr) uintptr {
    if os.Getenv("WAILS_RECOVERY_TEST_FAULT") == kind {
        fmt.Fprintf(os.Stderr, "INJECT_HRESULT %s actual=0x%x injected=0x80004005\n", kind, actual)
        if os.Getenv("WAILS_RECOVERY_TEST_ONCE") == "1" {
            os.Unsetenv("WAILS_RECOVERY_TEST_FAULT")
        }
        return 0x80004005
    }
    return actual
}
'@

$targets = @(
    @('ICoreWebView2Controller2.go', 'ICoreWebView2Controller2', 'PutDefaultBackgroundColor', 'background'),
    @('ICoreWebViewSettings.go', 'ICoreWebViewSettings', 'PutAreDevToolsEnabled', 'devtools'),
    @('corewebview2.go', 'ICoreWebView2', 'AddWebResourceRequestedFilter', 'resource')
)

New-Item -ItemType Directory $temporary | Out-Null
Push-Location $root
try {
    $replacements = @{}
    foreach ($target in $targets) {
        $original = Join-Path $edge $target[0]
        $source = [IO.File]::ReadAllText($original)
        # The real COM method runs first; replace its HRESULT at the error check.
        $patched = Add-FaultInjection $source "func (i *$($target[1])) $($target[2])(" `
            "`tif windows.Handle(hr)" "`thr = recoveryTestHRESULT(`"$($target[3])`", hr)`n"
        $path = Join-Path $temporary $target[0]
        [IO.File]::WriteAllText($path, $patched, $utf8)
        $replacements[$original] = $path
    }
    $original = Join-Path $edge 'chromium.go'
    # Fail before environment creation, after EmbedWithError installs its cleanup.
    $patched = Add-FaultInjection ([IO.File]::ReadAllText($original)) `
        'func (e *Chromium) EmbedWithError(' "`te.hwnd = hwnd" `
        "`tif hr := recoveryTestHRESULT(`"embed`", 0); hr != 0 { return windows.Errno(hr) }`n"
    $path = Join-Path $temporary 'chromium.go'
    [IO.File]::WriteAllText($path, ($patched + $helper), $utf8)
    $replacements[$original] = $path
    $overlay = Join-Path $temporary 'overlay.json'
    [IO.File]::WriteAllText($overlay, (@{Replace = $replacements} | ConvertTo-Json), $utf8)
    foreach ($tags in @('webview2_recovery_test', 'webview2_recovery_test,production')) {
        & go test -overlay $overlay -tags $tags ./pkg/application `
            -run '^TestLiveWebviewRecoveryFailures$' -count=1 -v -timeout=2m
        if ($LASTEXITCODE -ne 0) { throw "Recovery tests failed with tags $tags" }
    }
} finally {
    Pop-Location
    Remove-Item -LiteralPath $temporary -Recurse -Force
}
