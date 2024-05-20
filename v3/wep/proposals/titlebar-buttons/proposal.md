# Wails Enhancement Proposal (WEP)

## Title

**Author**: Lea Anthony 
**Created**: 2024-05-20

## Summary

Create a better API for controlling the appearance and functionality of window buttons in Wails.

## Motivation

Currently, it's only possible to control the visibility of the window buttons (minimize, maximize, close) in Wails. 
This proposal aims to enhance the control over titlebar buttons for applications on Windows and macOS by introducing a 
tri-state system (active, inactive, hidden) for each button. This feature will allow developers to have granular 
control over the appearance and functionality of window buttons, enhancing the customization capabilities of desktop 
applications built with Wails

## Detailed Design

1. A new enum will be added:

```go
      type ButtonState int

      const (
      ButtonActive ButtonState   = 0
      ButtonInactive ButtonState = 1
      ButtonHidden ButtonState   = 2
      )
```

2. The following options will be added to the `WebviewWindowOptions` option struct:

```go
   MinimiseButtonState ButtonState
   MaximiseButtonState ButtonState
   CloseButtonState    ButtonState
```
The following options will be removed from the `WindowsWindow` options:

- DisableMinimiseButton
- DisableMaximiseButton

The following options will be removed from the `MacWindow` options:

- DisableMinimiseButton
- DisableMaximiseButton
- DisableCloseButton

3. The following methods will be added to the `Window` interface:

```go
   SetMinimizeButtonState(state ButtonState)
   SetMaximizeButtonState(state ButtonState)
   SetCloseButtonState(state ButtonState)
```


## Pros/Cons

### Pros

1. Enhanced Customization
Provides developers with the ability to tailor the window controls (minimize, maximize, close) to the specific needs 
of their application. This can improve the user interface by aligning the window's behavior with the application's 
design and functional requirements.

2. Improved User Experience
By controlling the state of window buttons, applications can prevent user actions that might not make sense in certain 
contexts, such as disabling the minimize button during a critical operation, thereby enhancing the overall user experience.

3. Cross-Platform Consistency
The proposal includes handling for both Windows and macOS, ensuring that applications behave consistently across these 
platforms. This reduces the complexity for developers aiming to provide a uniform experience on different operating 
systems.

### Cons

1. Platform-Specific Bugs

The proposal includes handling for both Windows and macOS. There may well be subtle differences in how the controls
are handled on each platform. 

2. Top Level Options

Having the options at the top level of the `WebviewWindowOptions` option struct may be confusing for developers
as it suggests it works on Linux. However, the options are only relevant to macOS and Windows. The reason for this
is to provide a better API for calling these methods at runtime.

## Alternatives Considered

Original design used booleans but we'd like to give the option for disabling as well as hiding.

## Backwards Compatibility

The new features are optional and do not affect existing applications unless they opt to use them. This ensures that
current applications built with Wails will continue to function as before without any modifications.

## Test Plan

As part of the implementation, an example will be created in `v3/examples` to test the functionality.

## Reference Implementation

There is a reference implementation as part of this proposal.

## Maintenance Plan

This feature will be maintained and supported by the Wails developers.

## Conclusion

This API would be a leap forward in giving developers greater control over their application window appearances.



### Wails Enhancement Proposal (WEP)

#### Title
Implementing Tri-State Titlebar Button Control in Wails

#### Abstract
This proposal aims to enhance the control over titlebar buttons (minimize, maximize, close) in the Wails framework for applications on Windows and macOS by introducing a tri-state system (active, inactive, hidden) for each button. This feature will allow developers to have granular control over the appearance and functionality of window buttons, enhancing the customization capabilities of desktop applications built with Wails.

#### Specification
1. **Platform Detection**:
    - Detect the operating system at runtime using Go's `runtime` package to apply platform-specific logic for window controls.

2. **Window Struct Extension**:
    - Extend the `Window` struct in Wails to include properties for the state of the titlebar buttons using an enum to represent the three states:
      BACKTICKSgo
      type ButtonState int

      const (
      ButtonActive ButtonState = iota
      ButtonInactive
      ButtonHidden
      )

      type Window struct {
      ...
      BACKTICKS
    - This extension allows the window's titlebar buttons to be dynamically controlled through the Wails application.

3. **Interfacing with Native APIs**:
    - **Windows**: Use the Windows API to manipulate window styles (`WS_MINIMIZEBOX`, `WS_MAXIMIZEBOX`, `WS_SYSMENU`) and system commands to reflect the desired button states.
    - **macOS**: Utilize Cocoa APIs through Objective-C bridging to manipulate properties of `NSWindow` buttons, including their visibility and enabled states.

4. **API Functions**:
    - Implement Go functions that interface with native code to set the state of the window buttons. Example functions could be:
      BACKTICKSgo
      func (w *Window) SetMinimizeButtonState(state ButtonState) {
      // Platform-specific logic to set the minimize button state
      }

      func (w *Window) SetMaximizeButtonState(state ButtonState) {
      // Platform-specific logic to set the maximize button state
      }

      func (w *Window) SetCloseButtonState(state ButtonState) {
      // Platform-specific logic to set the close button state
      }
      BACKTICKS

5. **Binding to Frontend**:
    - Expose these functions to the frontend via Wails' bindings so they can be called from the JavaScript context. This allows frontend code to control the appearance of the window based on application state or user interactions.

#### Motivation
The ability to control the state of window titlebar buttons (active, inactive, hidden) is crucial for creating polished, user-friendly desktop applications. This feature will provide Wails developers with the tools needed to tailor the window controls to fit the design and functionality requirements of their applications, improving both aesthetics and usability.

#### Rationale
Implementing this feature as described will leverage existing platform-specific APIs to provide a consistent and reliable way to manage window button states across different operating systems, thereby enhancing the Wails framework's capabilities and appeal to developers looking for advanced UI customization options.

#### Backwards Compatibility
This proposal introduces new functionality that does not affect existing applications unless they opt to use the new features. Therefore, it is fully backward compatible.

#### Reference Implementation
A reference implementation will be provided as part of the development process, demonstrating the use of the new window button state controls in a sample Wails application.

#### Open Issues
- Detailed error handling strategies need to be discussed and implemented to ensure robust application behavior.
- Consideration of how these controls interact with other window management features in Wails.

#### Future Possibilities
Further enhancements could include more detailed customization options for the appearance of the buttons themselves, potentially allowing for completely custom designs and interactions.
