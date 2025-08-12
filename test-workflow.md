# Testing the Nightly Release Workflow

## Method 1: Fork Testing (Recommended)

1. **Create a fork** of the Wails repository
2. **Push the workflow** to your fork
3. **Test manually** using `workflow_dispatch`
4. **Verify behavior** without affecting main repo

```bash
# In your fork
git remote add upstream https://github.com/wailsapp/wails.git
git push origin master  # Push workflow to your fork
```

## Method 2: Local Script Testing

Create local test scripts to validate the logic:

```bash
# Test changelog parsing
./test-changelog-extraction.sh

# Test version increment logic  
./test-version-logic.sh

# Test commit analysis
./test-commit-detection.sh
```

## Method 3: Dry Run Workflow

Add a `dry_run` input parameter to test without creating releases:

```yaml
workflow_dispatch:
  inputs:
    dry_run:
      description: 'Run in dry-run mode (no releases created)'
      default: true
      type: boolean
```

## Method 4: Act (GitHub Actions Local Runner)

Use `act` to run GitHub Actions locally:

```bash
brew install act
act workflow_dispatch -W .github/workflows/nightly-releases.yml
```

## Testing Checklist

- [ ] Changelog parsing works correctly
- [ ] Version increment logic is accurate  
- [ ] Conventional commit detection works
- [ ] Release notes format properly
- [ ] Authorization checks function
- [ ] Branch handling (master vs v3-alpha)
- [ ] Error handling and fallbacks