# End-to-end regression for #6114. Requires Go, Node.js and the Windows SDK.
param(
    [ValidateSet('amd64', 'arm64')][string]$Arch = 'amd64',
    [string]$WorkDir = (Join-Path ([IO.Path]::GetTempPath()) ('wails-msix-smoke-' + [guid]::NewGuid()))
)
$ErrorActionPreference = 'Stop'
$repo = (Resolve-Path (Join-Path $PSScriptRoot '../..')).Path
New-Item -ItemType Directory -Force $WorkDir | Out-Null
$WorkDir = (Resolve-Path $WorkDir).Path
$oldPath = $env:PATH
Push-Location (Join-Path $repo 'v3')
try {
    go build -o (Join-Path $WorkDir 'wails3.exe') ./cmd/wails3
    if ($LASTEXITCODE -ne 0) { throw 'CLI build failed' }
    # Deliberately do not add the SDK to PATH: exercise SDK discovery too.
    $env:PATH = "$WorkDir;$oldPath"
    Set-Location $WorkDir
    wails3 init -n msix-smoke -t vanilla -skipgomodtidy
    if ($LASTEXITCODE -ne 0) { throw 'Project initialization failed' }
    Set-Location msix-smoke
    go mod edit "-replace=github.com/wailsapp/wails/v3=$repo/v3"
    if ($LASTEXITCODE -ne 0) { throw 'Local module replacement failed' }
    # Repeat to ensure replacing an existing output never prompts in CI.
    foreach ($attempt in 1..2) {
        # Windows PowerShell 5 treats redirected native stderr as error records.
        # Task logs to stderr on success; use its exit code as the failure signal.
        $ErrorActionPreference = 'Continue'
        wails3 task package FORMAT=msix "ARCH=$Arch"
        $packageExit = $LASTEXITCODE
        $ErrorActionPreference = 'Stop'
        if ($packageExit -ne 0) { throw "MSIX packaging failed on attempt $attempt" }
    }
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    $package = Join-Path (Get-Location) "bin/msix-smoke-$Arch.msix"
    $zip = [IO.Compression.ZipFile]::OpenRead($package)
    try {
        $reader = [IO.StreamReader]::new($zip.GetEntry('AppxManifest.xml').Open())
        try { [xml]$manifest = $reader.ReadToEnd() } finally { $reader.Dispose() }
        $expectedArch = if ($Arch -eq 'amd64') { 'x64' } else { 'arm64' }
        if ($manifest.Package.Identity.ProcessorArchitecture -ne $expectedArch) { throw 'Wrong package architecture' }
        $exe = $manifest.Package.Applications.Application.Executable
        $entry = $zip.GetEntry($exe)
        if (!$entry -or $entry.Length -eq 0) { throw 'Manifest executable missing from package' }
        # Verify the PE machine type, not just the manifest's architecture label.
        $stream = $entry.Open()
        $memory = [IO.MemoryStream]::new()
        try {
            $stream.CopyTo($memory)
            $bytes = $memory.ToArray()
            $pe = [BitConverter]::ToInt32($bytes, 0x3c)
            $machine = [BitConverter]::ToUInt16($bytes, $pe + 4)
            $expectedMachine = if ($Arch -eq 'amd64') { 0x8664 } else { 0xaa64 }
            if ($machine -ne $expectedMachine) { throw 'Wrong executable architecture' }
        } finally { $stream.Dispose(); $memory.Dispose() }
        foreach ($asset in 'StoreLogo', 'Square150x150Logo', 'Square44x44Logo', 'Wide310x150Logo', 'SplashScreen') {
            if (!$zip.GetEntry("Assets/$asset.png")) { throw "Missing asset: $asset" }
        }
    } finally { $zip.Dispose() }
    Write-Host "MSIX smoke test passed: $package"
} finally {
    $env:PATH = $oldPath
    Pop-Location
}
