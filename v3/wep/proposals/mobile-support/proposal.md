# Wails Enhancement Proposal (WEP)

## Title
Mobile Platform Support for Wails 3 (Android First)

**Author(s)**: [Your Name]  
**Created**: 2024-02-16

## Summary

Add support for mobile platforms (initially Android) to Wails 3, allowing developers to create mobile applications using the same codebase and API patterns as desktop applications. This proposal outlines the approach to extend Wails' existing architecture to support mobile platforms while maintaining API compatibility where possible, using gomobile as the primary build tool.

## Motivation

Currently, Wails only supports desktop platforms (Windows, macOS, Linux). Many developers want to target mobile platforms without maintaining separate codebases or learning different frameworks. By adding mobile support to Wails 3:

1. Developers can leverage their existing Wails knowledge and codebase
2. Applications can share business logic between desktop and mobile versions
3. The Go ecosystem gains another viable option for mobile app development
4. Wails becomes a more comprehensive cross-platform solution
5. CI/CD pipelines can be simplified through gomobile integration

## Detailed Design

### Architecture Overview

The mobile implementation will follow a single-window architecture for several reasons:
1. Aligns with mobile platform conventions
2. Simplifies lifecycle management
3. Reduces complexity in platform-specific code
4. Better user experience on mobile devices

The implementation will be conditional based on build tags:

```go
//go:build android
```

#### Design Principles

1. **Maintain API Compatibility**
    - All existing APIs should continue to work
    - Mobile-specific features should be no-ops on desktop
    - Use interfaces for platform-specific implementations
    - Leverage existing patterns where possible

2. **Leverage Existing Systems**
    - Use the current event system for lifecycle management
    - Keep the asset server architecture
    - Maintain service routing functionality
    - Use existing dialog API

3. **Security First**
    - Use WebView message handlers over JavascriptInterface
    - Maintain existing security patterns
    - Safe defaults for mobile permissions

#### Core Components

1. **Application Structure**
   ```go
   type androidPlatform struct {
       context  *android.Context
       webview  *android.WebView 
       activity *android.Activity
       parent   *App
       
       msgProcessor *MessageProcessor
       
       // Atomic flags for state management
       isPaused    atomic.Bool
       isDestroyed atomic.Bool
   }
   ```

2. **Message Processing**
    - WebView message handlers preferred over JavascriptInterface for:
        * Better security (no reflection-based attacks)
        * Better performance (no reflection overhead)
        * Better encapsulation
        * Consistent with existing message processor pattern

3. **Event System Integration**
   ```go
   // Core lifecycle events
   const (
       // Application creation and initialization
       MobileApplicationCreate  = "mobile:ApplicationCreate"
       MobileApplicationStart   = "mobile:ApplicationStart"
       
       // Foreground/Background transitions
       MobileApplicationResume  = "mobile:ApplicationResume"
       MobileApplicationPause   = "mobile:ApplicationPause"
       
       // Stopping and destruction
       MobileApplicationStop    = "mobile:ApplicationStop"
       MobileApplicationDestroy = "mobile:ApplicationDestroy"
       
       // User interaction and system events
       MobileBackPressed       = "mobile:BackPressed"
       MobileLowMemory        = "mobile:LowMemory"
       MobileNewIntent        = "mobile:NewIntent"
       
       // Configuration changes
       MobileConfigChange     = "mobile:ConfigChange"
   )

   // Event context data structures
   type ApplicationCreateData struct {
       SavedInstanceState map[string]interface{} // Restored state if app was killed
       Intent            *Intent                 // Launch intent
       IsFirstLaunch     bool                   // True if first time launch
   }

   type ApplicationResumeData struct {
       PauseDuration    time.Duration // How long app was paused
       IsConfigChange   bool          // True if resuming from config change
       PreviousState    string        // Previous app state
       NetworkStatus    string        // Current network connectivity
   }

   type ApplicationPauseData struct {
       IsFinishing    bool      // True if app is being terminated
       IsConfigChange bool      // True if pausing for config change
       PauseTime      time.Time
       RemainingMemory int64    // Available system memory
   }

   type ApplicationStopData struct {
       StopReason     string    // Why the app is stopping
       IsFinishing    bool      // True if app is being terminated
       Duration       time.Duration // How long app was running
       MemoryUsage    int64    // App's memory consumption
   }

   type ApplicationDestroyData struct {
       IsFinishing     bool     // True if normal termination
       IsConfigChange  bool     // True if due to config change
       StateToSave     map[string]interface{} // State to persist
       RunDuration     time.Duration // Total runtime
   }

   type BackPressedData struct {
       CanGoBack       bool     // True if WebView can go back
       WebViewHistory  int      // Number of history entries
       TimeSinceLastPress time.Duration // Time since last back press
   }

   type LowMemoryData struct {
       AvailableMemory int64   // Remaining system memory
       Severity        string  // Warning level: "moderate", "critical"
       RecommendedAction string // Suggested action
   }
   ```

4. **Service Support**
    - Maintain existing service interfaces:
        * ServiceStartup
        * ServiceShutdown
        * ServiceName
        * http.Handler
    - Services can hook into mobile events:
   ```go
   func (s *MyService) ServiceStartup(ctx context.Context, options ServiceOptions) error {
       app := application.Get()
       app.RegisterHook("mobile:ApplicationPause", func(event *CustomEvent) {
           // Handle pause
       })
       return nil
   }
   ```

5. **Asset Handling**
    - Maintain current asset server architecture
    - Keep embed.FS support for consistency
    - Use standard URL-based asset requests
    - Implement mobile-specific optimizations in assetserver_mobile.go

6. **Dialog System**
    - Use existing dialog API
    - Map to Android system dialogs:
        * MessageDialog → AlertDialog
        * OpenFileDialog → Intent.ACTION_OPEN_DOCUMENT
        * SaveFileDialog → Intent.ACTION_CREATE_DOCUMENT

#### Mobile-Specific Features

1. **Deep Linking**
   ```go
   type DeepLinkingOptions struct {
       Schemes        []string        
       Hosts         []string        
       IntentFilters []IntentFilter 
       Handlers      DeepLinkHandlers
   }
   ```

2. **Lifecycle Management**
    - Emit platform events for all major lifecycle changes
    - Allow services to hook into lifecycle events
    - Maintain state consistency across the application

## Pros/Cons

### Pros
1. Maintains API compatibility with existing applications
2. Leverages existing Wails patterns and systems
3. Security-focused design choices
4. Minimal learning curve for existing Wails developers
5. Flexible service architecture maintained
6. CI/CD friendly through gomobile

### Cons
1. Single window limitation on mobile
2. Some desktop features may not be available
3. Additional complexity in platform-specific code
4. Performance overhead of WebView vs native UI

## Alternatives Considered

The main design decisions that were considered:

1. **JS Bridge Implementation**
    - JavascriptInterface: Rejected due to security concerns and performance
    - WebView message handlers: Chosen for security and consistency
    - Custom protocol: Rejected as unnecessary complexity

2. **Asset Handling**
    - Android Assets: Considered but adds complexity
    - Current embed.FS: Chosen for consistency and simplicity
    - Custom asset provider interface: Potential future enhancement

3. **Dialog Implementation**
    - DialogFragment: Rejected as too complex
    - System dialogs: Chosen for native feel and simplicity
    - Custom dialogs: Unnecessary given requirements

## Backwards Compatibility

This proposal maintains backward compatibility by:
1. Not modifying existing APIs
2. Making mobile features no-ops on desktop
3. Using existing patterns and interfaces
4. Maintaining service architecture
5. Keeping current asset handling

## Test Plan

Testing focuses on automation with clear goals:

### Testing Goals
1. Ensure core Wails functionality works correctly
2. Verify mobile-specific features
3. Validate backward compatibility
4. Confirm native integration points

### Testing Layers

1. **Unit Tests**
    - Go unit tests for non-UI capabilities
    - Test mobile message processing
    - Test lifecycle management
    - Test deep linking handlers

2. **Integration Tests**
    - Automated Espresso tests for Android
    - Test communication between Go and WebView
    - Test asset loading
    - Test application lifecycle

3. **Automated UI Tests**
    - Appium-based testing
    - Test navigation flows
    - Test platform features
    - Test dialog integration

4. **CI Pipeline Tests**
    - Firebase Test Lab integration
    - Automated build verification
    - Cross-platform compatibility

## Reference Implementation

The reference implementation will be done in phases:

1. Core Platform Integration
    - Android Activity setup
    - WebView configuration
    - Message handler integration
    - Basic lifecycle management

2. Event System Integration
    - Mobile event implementation
    - Lifecycle event mapping
    - Service hooks integration

3. Asset Server Adaptation
    - Mobile-specific optimizations
    - URL handling
    - Service routing

4. Dialog Implementation
    - System dialog mapping
    - Intent handling
    - Result processing

## Maintenance Plan

1. **Version Support**
    - Target Android API 21+ initially
    - Regular updates for new Android versions
    - Automated testing for version compatibility

2. **Performance Monitoring**
    - Regular performance benchmarking
    - Memory usage monitoring
    - startup time optimization

3. **Security Updates**
    - Regular security audits
    - WebView security patches
    - Permission system updates

## Conclusion

This proposal provides a path to mobile support while maintaining Wails' core strengths:
1. API compatibility
2. Familiar patterns
3. Strong security
4. Flexible architecture

The implementation focuses on Android first, with a design that can extend to other mobile platforms in the future.