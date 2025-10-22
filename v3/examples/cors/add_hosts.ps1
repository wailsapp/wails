# PowerShell script to add hosts entry
$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
$entry = "127.0.0.1    app-local.wails-awesome.io"

# Check if already exists
$content = Get-Content $hostsPath -ErrorAction SilentlyContinue
if ($content -notcontains $entry -and -not ($content | Select-String -Pattern "app-local.wails-awesome.io" -Quiet)) {
    # Add the entry
    Add-Content -Path $hostsPath -Value "`n$entry"
    Write-Host "✅ Added hosts entry for app-local.wails-awesome.io" -ForegroundColor Green
} else {
    Write-Host "✅ Hosts entry already exists for app-local.wails-awesome.io" -ForegroundColor Yellow
}

# Verify
if (Select-String -Path $hostsPath -Pattern "app-local.wails-awesome.io" -Quiet) {
    Write-Host "✅ Verified: Entry is present in hosts file" -ForegroundColor Green

    # Test resolution
    $ping = Test-Connection -ComputerName "app-local.wails-awesome.io" -Count 1 -ErrorAction SilentlyContinue
    if ($ping) {
        Write-Host "✅ Name resolution working: app-local.wails-awesome.io resolves to $($ping.IPV4Address)" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Name resolution test failed - you may need to flush DNS cache" -ForegroundColor Yellow
        Write-Host "   Run: ipconfig /flushdns" -ForegroundColor Yellow
    }
} else {
    Write-Host "❌ Failed to add hosts entry" -ForegroundColor Red
}