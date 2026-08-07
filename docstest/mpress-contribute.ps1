$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$target = if ($args.Count -gt 0 -and $args[0]) { [string] $args[0] } else { 'https://github.com/wailsapp/wails.git' }
$draftFile = if ($args.Count -gt 1 -and $args[1]) { [string] $args[1] } else { "" }

function Start-Contribution([string] $Binary) {
  $contributionArguments = @("contribute", "--branch", 'master', $target)
  if ($draftFile) { $contributionArguments += @("--draft-file", $draftFile) }
  & $Binary @contributionArguments
  exit $LASTEXITCODE
}

$installedMPress = Get-Command mpress -CommandType Application -ErrorAction SilentlyContinue
if ($installedMPress) {
  Start-Contribution $installedMPress.Source
}

if ([Net.ServicePointManager]::SecurityProtocol -band [Net.SecurityProtocolType]::Tls12 -eq 0) {
  [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
}

$releaseBase = if ($env:MPRESS_RELEASE_BASE_URL) { $env:MPRESS_RELEASE_BASE_URL.TrimEnd('/') } else { "https://github.com/leaanthony/mpress/releases/latest/download" }
if (([uri]$releaseBase).Scheme -ne "https") { throw "MPRESS_RELEASE_BASE_URL must use HTTPS." }
$architecture = switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString()) {
  "X64" { "amd64" }
  "Arm64" { "arm64" }
  default { throw "M-Press does not provide a Windows release for architecture $($_)." }
}

function Download-ReleaseAsset([string] $Source, [string] $Destination) {
  $lastError = $null
  for ($attempt = 1; $attempt -le 3; $attempt++) {
    try {
      Invoke-WebRequest -UseBasicParsing -Uri $Source -OutFile $Destination
      return
    } catch {
      $lastError = $_
      if ($attempt -lt 3) { Start-Sleep -Seconds $attempt }
    }
  }
  throw "Could not download $Source after three attempts. $lastError"
}

$temporaryDirectory = Join-Path ([System.IO.Path]::GetTempPath()) ("mpress-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $temporaryDirectory | Out-Null
try {
  $asset = "mpress-windows-$architecture.zip"
  $archive = Join-Path $temporaryDirectory $asset
  $checksums = Join-Path $temporaryDirectory "checksums.txt"
  Download-ReleaseAsset "$releaseBase/$asset" $archive
  Download-ReleaseAsset "$releaseBase/checksums.txt" $checksums

  $escapedAsset = [regex]::Escape($asset)
  $checksumLine = Get-Content -LiteralPath $checksums | Where-Object { $_ -match "^([A-Fa-f0-9]{64})\s+\*?$escapedAsset$" } | Select-Object -First 1
  if (-not $checksumLine) { throw "The M-Press release does not contain a checksum for $asset." }
  $expected = ($checksumLine -split "\s+")[0].ToLowerInvariant()
  $actual = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
  if ($actual -ne $expected) { throw "Checksum verification failed for $asset." }

  Expand-Archive -LiteralPath $archive -DestinationPath $temporaryDirectory -Force
  $binary = Join-Path $temporaryDirectory "mpress.exe"
  if (-not (Test-Path -LiteralPath $binary -PathType Leaf)) {
    throw "The M-Press release archive does not contain mpress.exe."
  }
  Start-Contribution $binary
} finally {
  Remove-Item -LiteralPath $temporaryDirectory -Recurse -Force -ErrorAction SilentlyContinue
}
