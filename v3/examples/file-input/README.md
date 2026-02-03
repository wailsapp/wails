# File Input Example

Tests HTML `<input type="file">` and Wails JS runtime Dialog API functionality (#4862).

## Test Cases

**HTML File Input:**
1. **Single File** - Basic file selection
2. **Multiple Files** - Multi-select with `multiple` attribute
3. **Files or Directories** - Selection with `webkitdirectory` attribute

**Wails JS Runtime Dialog API:**
4. **Directory Only** - Using `wails.Dialogs.OpenFile()` with `CanChooseDirectories: true, CanChooseFiles: false`
5. **Filtered (.txt)** - Using `wails.Dialogs.OpenFile()` with `Filters` option

## Run

```bash
go run .
```
