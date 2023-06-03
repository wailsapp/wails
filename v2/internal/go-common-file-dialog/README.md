# Common File Dialog bindings for Golang

[Project Home](https://github.com/harry1453/go-common-file-dialog)

This library contains bindings for Windows Vista and newer's
[Common File Dialogs](https://docs.microsoft.com/en-us/windows/win32/shell/common-file-dialog),
which is the standard system dialog for selecting files or folders to open or
save.

The Common File Dialogs have to be accessed via the
[COM Interface](https://en.wikipedia.org/wiki/Component_Object_Model), normally
via C++ or via bindings (like in C#) .

This library contains bindings for Golang. **It does not require CGO**, and
contains empty stubs for non-windows platforms (so is safe to compile and run on
platforms other than windows, but will just return errors at runtime).

This can be very useful if you want to quickly get a file selector in your
Golang application. The `cfdutil` package contains utility functions with a
single call to open and configure a dialog, and then get the result from it.
Examples for this are in [`_examples/usingutil`](_examples/usingutil). Or, if
you want finer control over the dialog's operation, you can use the base
package. Examples for this are in
[`_examples/notusingutil`](_examples/notusingutil).

This library is available under the MIT license.

Currently supported features:

- Open File Dialog (to open a single file)
- Open Multiple Files Dialog (to open multiple files)
- Open Folder Dialog
- Save File Dialog
- Dialog "roles" to allow Windows to remember different "last locations" for
  different types of dialog
- Set dialog Title, Default Folder and Initial Folder
- Set dialog File Filters
