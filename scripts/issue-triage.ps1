# issue-triage.ps1 - Script to help with quick issue triage
# Run this at the start of your GitHub time to quickly process issues

# Set your GitHub username
$GITHUB_USERNAME = "your-username"

# Get the latest 10 open issues that aren't assigned and aren't labeled as "awaiting feedback"
Write-Host "Fetching recent unprocessed issues..."
gh issue list --repo wailsapp/wails --limit 10 --json number,title,labels,assignees | Out-File -Encoding utf8 -FilePath "issues_temp.json"
$issues = Get-Content -Raw -Path "issues_temp.json" | ConvertFrom-Json
$newIssues = $issues | Where-Object { 
    $_.assignees.Count -eq 0 -and 
    ($_.labels.Count -eq 0 -or -not ($_.labels | Where-Object { $_.name -eq "awaiting feedback" }))
}

# Process each issue
Write-Host "`n===== Issues Needing Triage =====`n"
foreach ($issue in $newIssues) {
    $number = $issue.number
    $title = $issue.title
    $labelNames = $issue.labels | ForEach-Object { $_.name }
    $labelsStr = if ($labelNames) { $labelNames -join ", " } else { "none" }
    
    Write-Host "Issue #$number`: $title"
    Write-Host "Labels: $labelsStr`n"
    
    $continue = $true
    while ($continue) {
        Write-Host "Options:"
        Write-Host "  [v] View issue in browser"
        Write-Host "  [2] Add v2-only label"
        Write-Host "  [3] Add v3-alpha label"
        Write-Host "  [b] Add bug label"
        Write-Host "  [e] Add enhancement label"
        Write-Host "  [d] Add documentation label"
        Write-Host "  [w] Add webview2 label"
        Write-Host "  [f] Request more info (awaiting feedback)"
        Write-Host "  [c] Close issue (duplicate/invalid)"
        Write-Host "  [a] Assign to yourself"
        Write-Host "  [s] Skip to next issue"
        Write-Host "  [q] Quit script"
        $action = Read-Host "Enter action"
        
        switch ($action) {
            "v" {
                gh issue view $number --repo wailsapp/wails --web
            }
            "2" {
                Write-Host "Adding v2-only label..."
                gh issue edit $number --repo wailsapp/wails --add-label "v2-only"
            }
            "3" {
                Write-Host "Adding v3-alpha label..."
                gh issue edit $number --repo wailsapp/wails --add-label "v3-alpha"
            }
            "b" {
                Write-Host "Adding bug label..."
                gh issue edit $number --repo wailsapp/wails --add-label "Bug"
            }
            "e" {
                Write-Host "Adding enhancement label..."
                gh issue edit $number --repo wailsapp/wails --add-label "Enhancement"
            }
            "d" {
                Write-Host "Adding documentation label..."
                gh issue edit $number --repo wailsapp/wails --add-label "Documentation"
            }
            "w" {
                Write-Host "Adding webview2 label..."
                gh issue edit $number --repo wailsapp/wails --add-label "webview2"
            }
            "f" {
                Write-Host "Requesting more info..."
                gh issue comment $number --repo wailsapp/wails --body "Thank you for reporting this issue. Could you please provide additional information to help us investigate?`n`n- [Specific details needed]`n`nThis will help us address your issue more effectively."
                gh issue edit $number --repo wailsapp/wails --add-label "awaiting feedback"
            }
            "c" {
                $reason = Read-Host "Reason for closing (duplicate/invalid/etc)"
                gh issue comment $number --repo wailsapp/wails --body "Closing this issue: $reason"
                gh issue close $number --repo wailsapp/wails
            }
            "a" {
                Write-Host "Assigning to yourself..."
                gh issue edit $number --repo wailsapp/wails --add-assignee "$GITHUB_USERNAME"
            }
            "s" {
                Write-Host "Skipping to next issue..."
                $continue = $false
            }
            "q" {
                Write-Host "Exiting script."
                exit
            }
            default {
                Write-Host "Invalid option. Please try again."
            }
        }
        
        Write-Host ""
    }
    
    Write-Host "--------------------------------`n"
}

Write-Host "No more issues to triage!"

# Clean up temp file
Remove-Item -Path "issues_temp.json"
