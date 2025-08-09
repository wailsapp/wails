# pr-review-helper.ps1 - Script to help with efficient PR reviews
# Run this during your PR review time

# Set your GitHub username
$GITHUB_USERNAME = "your-username"

# Get open PRs that are ready for review
Write-Host "Fetching PRs ready for review..."
gh pr list --repo wailsapp/wails --json number,title,author,labels,reviewDecision,additions,deletions,baseRefName,headRefName --limit 10 | Out-File -Encoding utf8 -FilePath "prs_temp.json"
$prs = Get-Content -Raw -Path "prs_temp.json" | ConvertFrom-Json

# Process each PR
Write-Host "`n===== PRs Needing Review =====`n"
foreach ($pr in $prs) {
    $number = $pr.number
    $title = $pr.title
    $author = $pr.author.login
    $labels = if ($pr.labels) { $pr.labels | ForEach-Object { $_.name } | Join-String -Separator ", " } else { "none" }
    $reviewState = if ($pr.reviewDecision) { $pr.reviewDecision } else { "PENDING" }
    $baseRef = $pr.baseRefName
    $headRef = $pr.headRefName
    $changes = $pr.additions + $pr.deletions
    
    Write-Host "PR #$number`: $title"
    Write-Host "Author: $author"
    Write-Host "Labels: $labels"
    Write-Host "Branch: $headRef -> $baseRef"
    Write-Host "Changes: +$($pr.additions)/-$($pr.deletions) lines"
    Write-Host "Review state: $reviewState`n"
    
    # Determine complexity based on size
    $complexity = if ($changes -lt 50) {
        "Quick review"
    } elseif ($changes -lt 300) {
        "Moderate review"
    } else {
        "Extensive review"
    }
    
    Write-Host "Complexity: $complexity"
    
    $continue = $true
    while ($continue) {
        Write-Host "`nOptions:"
        Write-Host "  [v] View PR in browser"
        Write-Host "  [d] View diff in browser"
        Write-Host "  [c] Generate review checklist"
        Write-Host "  [a] Approve PR"
        Write-Host "  [r] Request changes"
        Write-Host "  [m] Add comment"
        Write-Host "  [l] Add labels"
        Write-Host "  [s] Skip to next PR"
        Write-Host "  [q] Quit script"
        $action = Read-Host "Enter action"
        
        switch ($action) {
            "v" {
                gh pr view $number --repo wailsapp/wails --web
            }
            "d" {
                gh pr diff $number --repo wailsapp/wails --web
            }
            "c" {
                # Generate review checklist
                $checklist = @"
## PR Review: $title

### Basic Checks:
- [ ] PR title is descriptive
- [ ] PR description explains the changes
- [ ] Related issues are linked

### Technical Checks:
- [ ] Code follows project style
- [ ] No unnecessary commented code
- [ ] Error handling is appropriate
- [ ] Documentation updated (if needed)
- [ ] Tests included (if needed)

### Impact Assessment:
- [ ] Changes are backward compatible (if applicable)
- [ ] No breaking changes to public APIs
- [ ] Performance impact considered

### Version Specific:
"@

                if ($baseRef -eq "master") {
                    $checklist += @"

- [ ] Appropriate for v2 maintenance
- [ ] No features that should be v3-only
"@
                } elseif ($baseRef -eq "v3-alpha") {
                    $checklist += @"

- [ ] Appropriate for v3 development
- [ ] Aligns with v3 roadmap
"@
                }

                # Write to clipboard
                $checklist | Set-Clipboard
                Write-Host "`nReview checklist copied to clipboard!`n"
            }
            "a" {
                $comment = Read-Host "Approval comment (blank for none)"
                if ($comment) {
                    gh pr review $number --repo wailsapp/wails --approve --body $comment
                } else {
                    gh pr review $number --repo wailsapp/wails --approve
                }
            }
            "r" {
                $comment = Read-Host "Feedback for changes requested"
                gh pr review $number --repo wailsapp/wails --request-changes --body $comment
            }
            "m" {
                $comment = Read-Host "Comment text"
                gh pr comment $number --repo wailsapp/wails --body $comment
            }
            "l" {
                $labels = Read-Host "Labels to add (comma-separated)"
                $labelArray = $labels -split ","
                foreach ($label in $labelArray) {
                    $labelTrimmed = $label.Trim()
                    if ($labelTrimmed) {
                        gh pr edit $number --repo wailsapp/wails --add-label $labelTrimmed
                    }
                }
            }
            "s" {
                Write-Host "Skipping to next PR..."
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
    }
    
    Write-Host "--------------------------------`n"
}

Write-Host "No more PRs to review!"

# Clean up temp file
Remove-Item -Path "prs_temp.json"
