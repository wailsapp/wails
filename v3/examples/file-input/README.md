# File Input Example

Tests HTML `<input type="file">` functionality on macOS (#4862).

## Test Cases

1. **Single File** - Basic file selection
2. **Multiple Files** - Multi-select with `multiple` attribute
3. **Files or Directories** - Selection with `webkitdirectory` attribute

## Run

```bash
go run .
```

## Notes

- For directory-only selection, use the Wails Dialog API with `CanChooseFiles(false)`
- The HTML `accept` attribute is not supported by macOS WKWebView
