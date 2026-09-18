# File Input Example

Tests HTML `<input type="file">` and Wails Dialog API functionality (#4862).

## Test Cases

**HTML File Input:**
1. **Single File** - Basic file selection
2. **Multiple Files** - Multi-select with `multiple` attribute
3. **Files or Directories** - Selection with `webkitdirectory` attribute

**Wails Dialog API (via generated bindings):**
4. **Directory Only** - `CanChooseDirectories(true)` + `CanChooseFiles(false)`
5. **Filtered (.txt)** - Using `AddFilter()` for file type filtering

## Run

```bash
# Generate bindings (already included)
wails3 generate bindings -d assets/bindings

# Run the app
go run .
```
