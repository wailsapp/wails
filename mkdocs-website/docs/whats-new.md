# What's new in v3?

!!! note
        The features that will be included in the v3 release may change from this list.

## Multiple Windows

It's now possible to create multiple windows and configure each one independently. 

## Systrays

Systrays allow you to add an icon in the system tray area of your desktop and have the following features:
    - Attach window (the window will be centered to the systray icon)
    - Full menu support 
    - Light/Dark mode icons

## Plugins

Plugins allow you to extend the functionality of the Wails system. Not only can plugin methods be used in Go, but also called from Javascript. Included plugins:
  - kvstore - A key/value store
  - browser - open links in a browser
  - log - custom logger
  - oauth - handles oauth authentication and supports 60 providers
  - single_instance - only allow one copy of your app to be run
  - sqlite - add a sqlite db to your app. Uses the modernc pure go library
  - start_at_login - Register/Unregister your application to start at login

## Improved bindings generation

v3 uses a new static analyser to generate bindings. This makes it extremely fast and maintains comments and parameter names in your bindings.

## Improved events

Events are now emitted for a lot of the runtime operations, allowing you to hook into application/system events. Cross-platform (common) events are also emitted where there are common platform events, allowing you to write the same event handling methods cross platform.

Event hooks can also be registered. These are like the `On` method but are synchronous and allow you to cancel the event. An example of this would be to show a confirmation dialog before closing a window. 

## Wails Markup Language (wml)

An experimental feature to call runtime methods using plain html, similar to [htmx](https://htmx.org).
