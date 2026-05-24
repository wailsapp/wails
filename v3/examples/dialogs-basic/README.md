# Dialog Test Application

This application is designed to test macOS file dialog functionality across different versions of macOS. It provides a comprehensive suite of tests for various dialog features and configurations.

## Features Tested

1. Basic file open dialog
2. Single extension filter
3. Multiple extension filter
4. Multiple file selection
5. Directory selection
6. Save dialog with extension
7. Complex filters
8. Hidden files
9. Default directory
10. Full featured dialog with all options

## Running the Tests

```bash
go run main.go
```

## Test Results

When running tests:
- Each test will show the selected file(s) and their types
- For multiple selections, all selected files will be listed
- Errors will be displayed in an error dialog
- The application logs debug information to help track issues

## Notes

- This test application is primarily for development and testing purposes
- It can be used to verify dialog behavior across different macOS versions
- The tests are designed to not interfere with CI pipelines
