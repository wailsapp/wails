# Wails Issue Management Automation

This directory contains automation workflows and scripts to help manage the Wails project with minimal time investment.

## GitHub Workflow Files

### 1. Auto-Label Issues (`auto-label-issues.yml`)
- Automatically labels issues and PRs based on their content and modified files
- Labels are defined in `issue-labeler.yml` and `file-labeler.yml`
- Activates when issues are opened, edited, or reopened

### 2. Issue Triage Automation (`issue-triage-automation.yml`)
- Performs automated actions for issue triage
- Requests more info for incomplete bug reports
- Prioritizes security issues
- Adds issues to appropriate project boards

## Configuration Files

### 1. Issue Content Labeler (`issue-labeler.yml`)
- Defines patterns to match in issue title/body
- Categorizes by version (v2/v3), component, type, and priority
- Customize patterns as needed for your project

### 2. File Path Labeler (`file-labeler.yml`)
- Labels PRs based on which files they modify
- Helps identify which areas of the codebase are affected
- Customize file patterns as needed

### 3. Stale Issues Config (`stale.yml`)
- Marks issues as stale after 45 days of inactivity
- Closes stale issues after an additional 10 days
- Exempts issues with important labels

## Helper Scripts

### 1. Issue Triage Script (`scripts/issue-triage.ps1`)
- PowerShell script to quickly triage issues
- Lists recent issues needing attention
- Provides easy keyboard shortcuts for common actions
- Run during your dedicated issue triage time

### 2. PR Review Helper (`scripts/pr-review-helper.ps1`)
- PowerShell script to efficiently review PRs
- Generates review checklists
- Provides easy shortcuts for common review actions
- Run during your dedicated PR review time

## How to Use This System

### Daily Workflow (2 hours max)

**Monday (120 min):**
1. Run `scripts/issue-triage.ps1` (30 min)
2. Run `scripts/pr-review-helper.ps1` (30 min)
3. Check Discord for critical discussions (30 min)
4. Plan your week (30 min)

**Tuesday-Wednesday (120 min/day):**
1. Quick check for urgent issues (10 min)
2. v3 development (110 min)

**Thursday (120 min):**
1. v2 maintenance (90 min)
2. Documentation updates (30 min)

**Friday (120 min):**
1. Run `scripts/pr-review-helper.ps1` (60 min)
2. Discord updates/newsletter (30 min)
3. Weekly reflection (30 min)

## Installation

1. The GitHub workflow files should be placed in `.github/workflows/`
2. Configuration files should be placed in `.github/`
3. Helper scripts should be placed in `scripts/`
4. Make sure you have GitHub CLI (`gh`) installed and authenticated

## Customization

Feel free to modify the configuration files and scripts to better suit your project's needs:

1. **Adding New Label Categories**:
   - Add new patterns to `issue-labeler.yml` for additional components or types
   - Update `file-labeler.yml` if you add new directories or file types

2. **Adjusting Automation Thresholds**:
   - Modify `stale.yml` to change how long issues remain active
   - Update `issue-triage-automation.yml` to change conditions for automated actions

3. **Customizing Scripts**:
   - Update the scripts with your specific GitHub username
   - Add additional actions based on your workflow preferences
   - Adjust time allocations based on which tasks need more attention

## Benefits

This automated issue management system will:

1. **Save Time**: Reduce manual triage of most common issues
2. **Improve Consistency**: Apply the same categorization rules every time
3. **Increase Visibility**: Clear categorization helps community members find issues
4. **Focus Development**: Clearer separation of v2 and v3 work
5. **Reduce Backlog**: Better management of stale issues
6. **Streamline Reviews**: Faster PR processing with guided workflows

## Requirements

- GitHub CLI (`gh`) installed and authenticated
- PowerShell 5.1+ for Windows scripts
- GitHub Actions enabled on your repository
- Appropriate permissions to modify workflows

## Maintenance

This system requires minimal maintenance:

- Periodically review and update label patterns as your project evolves
- Adjust time allocations based on where you need to focus
- Update scripts if GitHub CLI commands change
- Customize the workflow as you find pain points in your process

Remember that the goal is to maximize your limited time (2 hours per day) by automating repetitive tasks and streamlining essential ones.
