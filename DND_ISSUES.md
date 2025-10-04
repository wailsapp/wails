# Wails Drag and Drop Issues Report

## Summary
This report compiles all drag and drop related issues found in Wails v3 GitHub repository. The issues span across multiple platforms (Windows, macOS, Linux) and various versions of Wails.

---

## Issue #1: [V3] Drag-and-drop broken on Windows since alpha 19

**Title:** [V3] Drag-and-drop broken on windows since alpha 19  
**Platform:** Windows 10  
**Reported:** 2024-08-09  
**Location:** GitHub  
**Reporter:** xob0t  
**Link:** https://github.com/wailsapp/wails/issues/4489  

**Summary:** No drag and drop events are triggered on the Go side in Wails v3 alpha 20+. The issue appears to have been introduced by PR #4318.

**Reproduction:**
1. Run dnd example on Windows
2. Observe that no DnD events are triggered

**System:** Windows 10 Pro, Wails v3.0.0-alpha.20, AMD Ryzen 5 5600X, NVIDIA GeForce RTX 3070

---

## Issue #2: Windows 10 doesn't receive drag and drop events on the backend

**Title:** Windows 10 doesn't receive drag and drop events on the backend  
**Platform:** Windows 10  
**Reported:** 2025-01-07  
**Location:** GitHub  
**Reporter:** makew0rld  
**Link:** https://github.com/wailsapp/wails/issues/3985  

**Summary:** Backend drop handler registered with `runtime.OnFileDrop` never gets executed when dropping files on Windows 10, though JavaScript drop events do fire. The event is not being propagated to the backend properly.

**Reproduction:**
1. Clone https://github.com/makew0rld/wails-issue-3985
2. Build for Windows: `wails build -platform windows/amd64 -windowsconsole -debug`
3. Open the .exe in Windows
4. Open web console using right-click menu
5. Drag in a file
6. Observe nothing appears in backend console, but drop event appears in web console

**System:** Built on Arch Linux, tested on Windows 10 VM, Wails v2.9.2

---

## Issue #3: [v3] [mac] drag n drop does not work as expected after resizing window

**Title:** [v3] [mac] drag n drop does not work as expected after resizing window  
**Platform:** macOS  
**Reported:** 2024-09-12  
**Location:** GitHub  
**Reporter:** Etesam913  
**Link:** https://github.com/wailsapp/wails/issues/3743  

**Summary:** On macOS, after resizing the window, drag and drop usually gets rejected when trying to drag a file into the window.

**Reproduction:**
1. `git checkout v3-alpha`
2. `cd v3/examples/drag-n-drop`
3. `go run main.go`
4. Resize window
5. Drag file into window
6. File gets rejected

**System:** macOS 15.0, Apple M2, Wails v3.0.0-alpha.6

---

## Issue #4: Webview failed to set AllowExternalDrag to false on Win 10

**Title:** Webview failed to set AllowExternalDrag to false on Win 10 (after cross compile on macOS)  
**Platform:** Windows 10  
**Reported:** 2024-09-26  
**Location:** GitHub  
**Reporter:** barats  
**Link:** https://github.com/wailsapp/wails/issues/3782  

**Summary:** Two main issues: 1) WebView fails to set AllowExternalDrag to false, causing WebView to automatically open dragged files. 2) Backend `runtime.OnFileDrop` doesn't work - no events are triggered.

**Reproduction:**
1. Build windows/amd64 app on macOS
2. Move .exe files to Win10 and run
3. Drag and drop files into app window
4. WebView DnD not disabled, backend events not called

**System:** Built on macOS 15.0, tested on Windows 10, Wails v2.9.2

---

## Issue #5: [v2] Files DragAndDrop bugs

**Title:** [v2] Files DragAndDrop bugs  
**Platform:** Linux (Ubuntu)  
**Reported:** 2024-06-23  
**Location:** GitHub  
**Reporter:** Vovan-VE  
**Link:** https://github.com/wailsapp/wails/issues/3563  

**Summary:** Two issues: 1) When only `EnableFileDrop: true`, every D'n'D triggers an `OnDomReady` event. 2) When `EnableFileDrop: true, DisableWebViewDrop: true`, D'n'D doesn't work at all for files.

**Reproduction:**
```go
func (a *App) domReady(ctx context.Context) {
    runtime.OnFileDropOff(ctx)
    runtime.LogInfof(ctx, "DOM READY -------------------------------\n%+v", errors.New("WTF?"))
    runtime.OnFileDrop(ctx, a.onFileDrop)
}

func (a *App) onFileDrop(x, y int, paths []string) {
    runtime.LogInfof(a.ctx, "drop files: %#v\n%+v", paths, errors.New("WTF?"))
}
```

**System:** Ubuntu 22.04, AMD Ryzen 9 7950X, Wails v2.9.1 and v2.9.2

---

## Issue #6: File drag and drop displays file in Linux WebKit

**Title:** File drag and drop displays file in Linux WebKit  
**Platform:** Linux  
**Reported:** 2024-08-17  
**Location:** GitHub  
**Reporter:** makew0rld  
**Link:** https://github.com/wailsapp/wails/issues/3686  

**Summary:** When dragging and dropping a file onto the Window, the `OnFileDrop` handler runs successfully, but the entire GUI is replaced with a view of the file if the browser supports it (e.g., images, PDFs). This makes the dropped file unusable as all UI is gone.

**Reproduction:**
1. Enable file dropping and add a simple handler
2. Try and drop an image file with WebKit on Linux
3. Observe the UI get replaced with a view of the image

**System:** Arch Linux, Intel Core i7-1165G7, Wails v2.9.1

**Note:** Issue was closed on 2025-08-02 with a workaround involving preventDefault() on drop and dragover events.

---

## Issue #7: Drag And Drop "wails-drop-target-active" class issue

**Title:** Drag And Drop "wails-drop-target-active" class is only applied if cssDropProperty is set via raw `style`  
**Platform:** Windows  
**Reported:** 2025-07-13  
**Location:** GitHub  
**Reporter:** riannucci  
**Link:** https://github.com/wailsapp/wails/issues/4419  

**Summary:** `OnFileDrop` does not apply the `wails-drop-target-active` class unless `--wails-drop-target: drop` is set directly via the element's `style` attribute, even though the `OnFileDrop` callback does fire.

**Reproduction:**
1. Set up a new wails project
2. Set EnableFileDrop in Application Options
3. Add CSS with `--wails-drop-target: drop` in stylesheet
4. Run app and drag file to target element
5. Base path gets set but element doesn't change styles

**System:** Windows 10/11 Pro, AMD Ryzen 7 5800X3D, Wails v2.10.2

**Note:** Issue was closed on 2025-08-02. Problem was in draganddrop.js using element.style instead of getComputedStyle.

---

## Issue #8: Use a user-agent that is more recognizable by front-end code

**Title:** Use a user-agent that is more recognizable by front-end code  
**Platform:** macOS  
**Reported:** 2024-01-31  
**Location:** GitHub  
**Reporter:** hsiafan  
**Link:** https://github.com/wailsapp/wails/issues/3226  

**Summary:** SortableJS (drag-and-drop library) doesn't work properly in Wails' webview because the userAgent doesn't contain standard browser identifiers like 'Safari/605.1.15', causing library detection to fail.

**Reproduction:**
Libraries like SortableJS check userAgent for specific strings ('Edge', 'Firefox', 'Chrome', 'Safari') and fail to work properly when these are not present.

**System:** macOS, affects drag-and-drop libraries

---

## Issue #9: Application soft crashes when using vue-virtual-scroller on Linux

**Title:** Application soft crashes when using vue-virtual-scroller on Linux  
**Platform:** Linux (Ubuntu)  
**Reported:** 2024-10-22  
**Location:** GitHub  
**Reporter:** ysmilda  
**Link:** https://github.com/wailsapp/wails/issues/3849  

**Summary:** Using vue-virtual-scroller component causes the application to crash on Linux builds, rendering a grey screen and hanging the inspector window. Works fine on Windows and in dev view.

**Reproduction:**
1. Create new wails instance with `wails init -n myproject -t vue-ts`
2. Add vue-virtual-scroller component
3. Create Linux build (`wails build -tags webkit2_41`)
4. Observe crash while rendering

**System:** Ubuntu 24.04, Intel Core i7-12700, Wails v2.9.2

---

## Issue #10: [DragAndDrop] panic when dropping non-file objects

**Title:** [DragAndDrop] panic when dropping non-file objects  
**Platform:** Windows  
**Reported:** 2024-07-09  
**Location:** GitHub  
**Reporter:** Alpa-1  
**Link:** https://github.com/wailsapp/wails/issues/3596  

**Summary:** A panic occurs when dropping non-file objects (like strings or browser tabs) into the webview. The application crashes with a nil pointer dereference.

**Reproduction:**
1. Enable DragAndDrop
2. Add a drop target on frontend with `--wails-drop-target`
3. Register a callback with `OnFileDrop((x,y,f) => { console.log(f) }, true)`
4. Drop a Browser Tab onto the element
5. Application panics

**System:** Windows 10 Enterprise, AMD Ryzen 7 PRO 4750U, Wails v2.9.1

**Note:** Issue was closed on 2024-08-18. A simple nil check fix was provided.

---

## Common Patterns Identified:

1. **Platform-specific issues:** Windows has the most reported issues, particularly with events not reaching the backend
2. **Version regression:** Multiple reports indicate v3 alpha 19+ broke drag and drop on Windows
3. **WebView configuration:** Several issues relate to DisableWebViewDrop not working correctly
4. **Cross-platform builds:** Issues when building for Windows on other platforms
5. **Frontend/Backend disconnect:** Common pattern of JavaScript events firing but Go handlers not receiving events
6. **UI replacement on Linux:** Specific issue where dropped files replace the entire UI in WebKit

## Recommendations:

1. Focus on Windows platform issues first as they are most prevalent
2. Investigate the changes in PR #4318 which allegedly broke Windows DnD
3. Review the event propagation mechanism between frontend and backend
4. Test cross-platform builds more thoroughly
5. Improve documentation on DragAndDrop configuration options