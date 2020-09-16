# Parse

Parse will attempt to parse your Wails project to perform a number of tasks:
  * Verify that you have bound struct pointers
  * Generate JS helper files/docs

It currently checks bindings correctly if your code binds using one of the following methods:
  * Literal Binding: `app.Bind(&MyStruct{})`
  * Variable Binding: `app.Bind(m)` - m can be `m := &MyStruct{}` or `m := newMyStruct()`
  * Function Binding: `app.Bind(newMyStruct())`
