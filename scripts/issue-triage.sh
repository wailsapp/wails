#!/bin/bash
# issue-triage.sh - Script to help with quick issue triage
# Run this at the start of your GitHub time to quickly process issues

# Set your GitHub username
GITHUB_USERNAME="your-username"

# Get the latest 10 open issues that aren't assigned and aren't labeled as "awaiting feedback"
echo "Fetching recent unprocessed issues..."
gh issue list --repo wailsapp/wails --limit 10 --json number,title,labels,assignees --jq '.[] | select(.assignees | length == 0) | select(any(.labels[]; .name != "awaiting feedback"))' > new_issues.json

# Process each issue
echo -e "\n===== Issues Needing Triage =====\n"
cat new_issues.json | jq -c '.[]' | while read -r issue; do
    number=$(echo $issue | jq -r '.number')
    title=$(echo $issue | jq -r '.title')
    labels=$(echo $issue | jq -r '.labels[] | .name' 2>/dev/null | tr '\n' ', ' | sed 's/,$//')
    
    if [ -z "$labels" ]; then
        labels="none"
    fi
    
    echo -e "Issue #$number: $title"
    echo -e "Labels: $labels\n"
    
    while true; do
        echo "Options:"
        echo "  [v] View issue in browser"
        echo "  [2] Add v2-only label"
        echo "  [3] Add v3-alpha label"
        echo "  [b] Add bug label"
        echo "  [e] Add enhancement label"
        echo "  [d] Add documentation label"
        echo "  [w] Add webview2 label"
        echo "  [f] Request more info (awaiting feedback)"
        echo "  [c] Close issue (duplicate/invalid)"
        echo "  [a] Assign to yourself"
        echo "  [s] Skip to next issue"
        echo "  [q] Quit script"
        read -p "Enter action: " action
        
        case $action in
            v)
                gh issue view $number --repo wailsapp/wails --web
                ;;
            2)
                echo "Adding v2-only label..."
                gh issue edit $number --repo wailsapp/wails --add-label "v2-only"
                ;;
            3)
                echo "Adding v3-alpha label..."
                gh issue edit $number --repo wailsapp/wails --add-label "v3-alpha"
                ;;
            b)
                echo "Adding bug label..."
                gh issue edit $number --repo wailsapp/wails --add-label "Bug"
                ;;
            e)
                echo "Adding enhancement label..."
                gh issue edit $number --repo wailsapp/wails --add-label "Enhancement"
                ;;
            d)
                echo "Adding documentation label..."
                gh issue edit $number --repo wailsapp/wails --add-label "Documentation"
                ;;
            w)
                echo "Adding webview2 label..."
                gh issue edit $number --repo wailsapp/wails --add-label "webview2"
                ;;
            f)
                echo "Requesting more info..."
                gh issue comment $number --repo wailsapp/wails --body "Thank you for reporting this issue. Could you please provide additional information to help us investigate?\n\n- [Specific details needed]\n\nThis will help us address your issue more effectively."
                gh issue edit $number --repo wailsapp/wails --add-label "awaiting feedback"
                ;;
            c)
                read -p "Reason for closing (duplicate/invalid/etc): " reason
                gh issue comment $number --repo wailsapp/wails --body "Closing this issue: $reason"
                gh issue close $number --repo wailsapp/wails
                ;;
            a)
                echo "Assigning to yourself..."
                gh issue edit $number --repo wailsapp/wails --add-assignee "$GITHUB_USERNAME"
                ;;
            s)
                echo "Skipping to next issue..."
                break
                ;;
            q)
                echo "Exiting script."
                exit 0
                ;;
            *)
                echo "Invalid option. Please try again."
                ;;
        esac
        
        echo ""
    done
    
    echo -e "--------------------------------\n"
done

echo "No more issues to triage!"
