(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

  // desktop/@wailsio/runtime/src/log.js
  function debugLog(message) {
    console.log(
      "%c wails3 %c " + message + " ",
      "background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem",
      "background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem"
    );
  }

  // desktop/@wailsio/runtime/src/application.js
  var application_exports = {};
  __export(application_exports, {
    Hide: () => Hide,
    Quit: () => Quit,
    Show: () => Show
  });

  // node_modules/nanoid/non-secure/index.js
  var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
  var nanoid = (size2 = 21) => {
    let id = "";
    let i = size2;
    while (i--) {
      id += urlAlphabet[Math.random() * 64 | 0];
    }
    return id;
  };

  // desktop/@wailsio/runtime/src/runtime.js
  var runtimeURL = window.location.origin + "/wails/runtime";
  var objectNames = {
    Call: 0,
    Clipboard: 1,
    Application: 2,
    Events: 3,
    ContextMenu: 4,
    Dialog: 5,
    Window: 6,
    Screens: 7,
    System: 8,
    Browser: 9,
    CancelCall: 10
  };
  var clientId = nanoid();
  function newRuntimeCallerWithID(object, windowName) {
    return function(method, args = null) {
      return runtimeCallWithID(object, method, windowName, args);
    };
  }
  function runtimeCallWithID(objectID, method, windowName, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID);
    url.searchParams.append("method", method);
    let fetchOptions = {
      headers: {}
    };
    if (windowName) {
      fetchOptions.headers["x-wails-window-name"] = windowName;
    }
    if (args) {
      url.searchParams.append("args", JSON.stringify(args));
    }
    fetchOptions.headers["x-wails-client-id"] = clientId;
    return new Promise((resolve, reject) => {
      fetch(url, fetchOptions).then((response) => {
        if (response.ok) {
          if (response.headers.get("Content-Type") && response.headers.get("Content-Type").indexOf("application/json") !== -1) {
            return response.json();
          } else {
            return response.text();
          }
        }
        reject(Error(response.statusText));
      }).then((data) => resolve(data)).catch((error) => reject(error));
    });
  }

  // desktop/@wailsio/runtime/src/application.js
  var call = newRuntimeCallerWithID(objectNames.Application, "");
  var HideMethod = 0;
  var ShowMethod = 1;
  var QuitMethod = 2;
  function Hide() {
    return call(HideMethod);
  }
  function Show() {
    return call(ShowMethod);
  }
  function Quit() {
    return call(QuitMethod);
  }

  // desktop/@wailsio/runtime/src/browser.js
  var browser_exports = {};
  __export(browser_exports, {
    OpenURL: () => OpenURL
  });
  var call2 = newRuntimeCallerWithID(objectNames.Browser, "");
  var BrowserOpenURL = 0;
  function OpenURL(url) {
    return call2(BrowserOpenURL, { url });
  }

  // desktop/@wailsio/runtime/src/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    SetText: () => SetText,
    Text: () => Text
  });
  var call3 = newRuntimeCallerWithID(objectNames.Clipboard, "");
  var ClipboardSetText = 0;
  var ClipboardText = 1;
  function SetText(text) {
    return call3(ClipboardSetText, { text });
  }
  function Text() {
    return call3(ClipboardText);
  }

  // desktop/@wailsio/runtime/src/system.js
  var system_exports = {};
  __export(system_exports, {
    Capabilities: () => Capabilities,
    Environment: () => Environment,
    IsAMD64: () => IsAMD64,
    IsARM: () => IsARM,
    IsARM64: () => IsARM64,
    IsDarkMode: () => IsDarkMode,
    IsDebug: () => IsDebug,
    IsLinux: () => IsLinux,
    IsMac: () => IsMac,
    IsWindows: () => IsWindows,
    invoke: () => invoke
  });
  var call4 = newRuntimeCallerWithID(objectNames.System, "");
  var systemIsDarkMode = 0;
  var environment = 1;
  function invoke(msg) {
    if (window.chrome) {
      return window.chrome.webview.postMessage(msg);
    }
    return window.webkit.messageHandlers.external.postMessage(msg);
  }
  function IsDarkMode() {
    return call4(systemIsDarkMode);
  }
  function Capabilities() {
    let response = fetch("/wails/capabilities");
    return response.json();
  }
  function Environment() {
    return call4(environment);
  }
  function IsWindows() {
    return window._wails.environment.OS === "windows";
  }
  function IsLinux() {
    return window._wails.environment.OS === "linux";
  }
  function IsMac() {
    return window._wails.environment.OS === "darwin";
  }
  function IsAMD64() {
    return window._wails.environment.Arch === "amd64";
  }
  function IsARM() {
    return window._wails.environment.Arch === "arm";
  }
  function IsARM64() {
    return window._wails.environment.Arch === "arm64";
  }
  function IsDebug() {
    return window._wails.environment.Debug === true;
  }

  // desktop/@wailsio/runtime/src/contextmenu.js
  window.addEventListener("contextmenu", contextMenuHandler);
  var call5 = newRuntimeCallerWithID(objectNames.ContextMenu, "");
  var ContextMenuOpen = 0;
  function openContextMenu(id, x, y, data) {
    void call5(ContextMenuOpen, { id, x, y, data });
  }
  function contextMenuHandler(event) {
    let element = event.target;
    let customContextMenu = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu");
    customContextMenu = customContextMenu ? customContextMenu.trim() : "";
    if (customContextMenu) {
      event.preventDefault();
      let customContextMenuData = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu-data");
      openContextMenu(customContextMenu, event.clientX, event.clientY, customContextMenuData);
      return;
    }
    processDefaultContextMenu(event);
  }
  function processDefaultContextMenu(event) {
    if (IsDebug()) {
      return;
    }
    const element = event.target;
    const computedStyle = window.getComputedStyle(element);
    const defaultContextMenuAction = computedStyle.getPropertyValue("--default-contextmenu").trim();
    switch (defaultContextMenuAction) {
      case "show":
        return;
      case "hide":
        event.preventDefault();
        return;
      default:
        if (element.isContentEditable) {
          return;
        }
        const selection = window.getSelection();
        const hasSelection = selection.toString().length > 0;
        if (hasSelection) {
          for (let i = 0; i < selection.rangeCount; i++) {
            const range = selection.getRangeAt(i);
            const rects = range.getClientRects();
            for (let j = 0; j < rects.length; j++) {
              const rect = rects[j];
              if (document.elementFromPoint(rect.left, rect.top) === element) {
                return;
              }
            }
          }
        }
        if (element.tagName === "INPUT" || element.tagName === "TEXTAREA") {
          if (hasSelection || !element.readOnly && !element.disabled) {
            return;
          }
        }
        event.preventDefault();
    }
  }

  // desktop/@wailsio/runtime/src/flags.js
  var flags_exports = {};
  __export(flags_exports, {
    GetFlag: () => GetFlag
  });
  function GetFlag(keyString) {
    try {
      return window._wails.flags[keyString];
    } catch (e) {
      throw new Error("Unable to retrieve flag '" + keyString + "': " + e);
    }
  }

  // desktop/@wailsio/runtime/src/drag.js
  window._wails = window._wails || {};
  window._wails.setResizable = setResizable;
  window._wails.endDrag = endDrag;
  window.addEventListener("mousedown", onMouseDown);
  window.addEventListener("mousemove", onMouseMove);
  window.addEventListener("mouseup", onMouseUp);
  var shouldDrag = false;
  var resizeEdge = null;
  var resizable = false;
  var defaultCursor = "auto";
  function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--webkit-app-region");
    if (!val || val === "" || val.trim() !== "drag" || e.buttons !== 1) {
      return false;
    }
    return e.detail === 1;
  }
  function setResizable(value) {
    resizable = value;
  }
  function endDrag() {
    document.body.style.cursor = "default";
    shouldDrag = false;
  }
  function testResize() {
    if (resizeEdge) {
      invoke(`resize:${resizeEdge}`);
      return true;
    }
    return false;
  }
  function onMouseDown(e) {
    if (IsWindows() && testResize() || dragTest(e)) {
      shouldDrag = !!isValidDrag(e);
    }
  }
  function isValidDrag(e) {
    return !(e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight);
  }
  function onMouseUp(e) {
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (mousePressed > 0) {
      endDrag();
    }
  }
  function setResize(cursor = defaultCursor) {
    document.documentElement.style.cursor = cursor;
    resizeEdge = cursor;
  }
  function onMouseMove(e) {
    shouldDrag = checkDrag(e);
    if (IsWindows() && resizable) {
      handleResize(e);
    }
  }
  function checkDrag(e) {
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (shouldDrag && mousePressed > 0) {
      invoke("drag");
      return false;
    }
    return shouldDrag;
  }
  function handleResize(e) {
    let resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
    let resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
    let cornerExtra = GetFlag("resizeCornerExtra") || 10;
    let rightBorder = window.outerWidth - e.clientX < resizeHandleWidth;
    let leftBorder = e.clientX < resizeHandleWidth;
    let topBorder = e.clientY < resizeHandleHeight;
    let bottomBorder = window.outerHeight - e.clientY < resizeHandleHeight;
    let rightCorner = window.outerWidth - e.clientX < resizeHandleWidth + cornerExtra;
    let leftCorner = e.clientX < resizeHandleWidth + cornerExtra;
    let topCorner = e.clientY < resizeHandleHeight + cornerExtra;
    let bottomCorner = window.outerHeight - e.clientY < resizeHandleHeight + cornerExtra;
    if (!leftBorder && !rightBorder && !topBorder && !bottomBorder && resizeEdge !== void 0) {
      setResize();
    } else if (rightCorner && bottomCorner)
      setResize("se-resize");
    else if (leftCorner && bottomCorner)
      setResize("sw-resize");
    else if (leftCorner && topCorner)
      setResize("nw-resize");
    else if (topCorner && rightCorner)
      setResize("ne-resize");
    else if (leftBorder)
      setResize("w-resize");
    else if (topBorder)
      setResize("n-resize");
    else if (bottomBorder)
      setResize("s-resize");
    else if (rightBorder)
      setResize("e-resize");
  }

  // desktop/@wailsio/runtime/src/screens.js
  var screens_exports = {};
  __export(screens_exports, {
    GetAll: () => GetAll,
    GetCurrent: () => GetCurrent,
    GetPrimary: () => GetPrimary
  });
  var call6 = newRuntimeCallerWithID(objectNames.Screens, "");
  var getAll = 0;
  var getPrimary = 1;
  var getCurrent = 2;
  function GetAll() {
    return call6(getAll);
  }
  function GetPrimary() {
    return call6(getPrimary);
  }
  function GetCurrent() {
    return call6(getCurrent);
  }

  // desktop/@wailsio/runtime/src/window.js
  var window_exports = {};
  __export(window_exports, {
    Center: () => Center,
    Close: () => Close,
    Fullscreen: () => Fullscreen,
    Get: () => Get,
    GetZoomLevel: () => GetZoomLevel,
    Height: () => Height,
    Hide: () => Hide2,
    Maximise: () => Maximise,
    Minimise: () => Minimise,
    RelativePosition: () => RelativePosition,
    Restore: () => Restore,
    Screen: () => Screen,
    SetAlwaysOnTop: () => SetAlwaysOnTop,
    SetBackgroundColour: () => SetBackgroundColour,
    SetMaxSize: () => SetMaxSize,
    SetMinSize: () => SetMinSize,
    SetRelativePosition: () => SetRelativePosition,
    SetResizable: () => SetResizable,
    SetSize: () => SetSize,
    SetTitle: () => SetTitle,
    SetZoomLevel: () => SetZoomLevel,
    Show: () => Show2,
    Size: () => Size,
    ToggleMaximise: () => ToggleMaximise,
    UnMaximise: () => UnMaximise,
    UnMinimise: () => UnMinimise,
    Width: () => Width,
    ZoomIn: () => ZoomIn,
    ZoomOut: () => ZoomOut,
    ZoomReset: () => ZoomReset
  });
  var center = 0;
  var setTitle = 1;
  var fullscreen = 2;
  var unFullscreen = 3;
  var setSize = 4;
  var size = 5;
  var setMaxSize = 6;
  var setMinSize = 7;
  var setAlwaysOnTop = 8;
  var setRelativePosition = 9;
  var relativePosition = 10;
  var screen = 11;
  var hide = 12;
  var maximise = 13;
  var unMaximise = 14;
  var toggleMaximise = 15;
  var minimise = 16;
  var unMinimise = 17;
  var restore = 18;
  var show = 19;
  var close = 20;
  var setBackgroundColour = 21;
  var setResizable2 = 22;
  var width = 23;
  var height = 24;
  var zoomIn = 25;
  var zoomOut = 26;
  var zoomReset = 27;
  var getZoomLevel = 28;
  var setZoomLevel = 29;
  var thisWindow = Get("");
  function createWindow(call10) {
    return {
      Get: (windowName) => createWindow(newRuntimeCallerWithID(objectNames.Window, windowName)),
      Center: () => call10(center),
      SetTitle: (title) => call10(setTitle, { title }),
      Fullscreen: () => call10(fullscreen),
      UnFullscreen: () => call10(unFullscreen),
      SetSize: (width2, height2) => call10(setSize, { width: width2, height: height2 }),
      Size: () => call10(size),
      SetMaxSize: (width2, height2) => call10(setMaxSize, { width: width2, height: height2 }),
      SetMinSize: (width2, height2) => call10(setMinSize, { width: width2, height: height2 }),
      SetAlwaysOnTop: (onTop) => call10(setAlwaysOnTop, { alwaysOnTop: onTop }),
      SetRelativePosition: (x, y) => call10(setRelativePosition, { x, y }),
      RelativePosition: () => call10(relativePosition),
      Screen: () => call10(screen),
      Hide: () => call10(hide),
      Maximise: () => call10(maximise),
      UnMaximise: () => call10(unMaximise),
      ToggleMaximise: () => call10(toggleMaximise),
      Minimise: () => call10(minimise),
      UnMinimise: () => call10(unMinimise),
      Restore: () => call10(restore),
      Show: () => call10(show),
      Close: () => call10(close),
      SetBackgroundColour: (r, g, b, a) => call10(setBackgroundColour, { r, g, b, a }),
      SetResizable: (resizable2) => call10(setResizable2, { resizable: resizable2 }),
      Width: () => call10(width),
      Height: () => call10(height),
      ZoomIn: () => call10(zoomIn),
      ZoomOut: () => call10(zoomOut),
      ZoomReset: () => call10(zoomReset),
      GetZoomLevel: () => call10(getZoomLevel),
      SetZoomLevel: (zoomLevel) => call10(setZoomLevel, { zoomLevel })
    };
  }
  function Get(windowName) {
    return createWindow(newRuntimeCallerWithID(objectNames.Window, windowName));
  }
  function Center() {
    thisWindow.Center();
  }
  function SetTitle(title) {
    thisWindow.SetTitle(title);
  }
  function Fullscreen() {
    thisWindow.Fullscreen();
  }
  function SetSize(width2, height2) {
    thisWindow.SetSize(width2, height2);
  }
  function Size() {
    return thisWindow.Size();
  }
  function SetMaxSize(width2, height2) {
    thisWindow.SetMaxSize(width2, height2);
  }
  function SetMinSize(width2, height2) {
    thisWindow.SetMinSize(width2, height2);
  }
  function SetAlwaysOnTop(onTop) {
    thisWindow.SetAlwaysOnTop(onTop);
  }
  function SetRelativePosition(x, y) {
    thisWindow.SetRelativePosition(x, y);
  }
  function RelativePosition() {
    return thisWindow.RelativePosition();
  }
  function Screen() {
    return thisWindow.Screen();
  }
  function Hide2() {
    thisWindow.Hide();
  }
  function Maximise() {
    thisWindow.Maximise();
  }
  function UnMaximise() {
    thisWindow.UnMaximise();
  }
  function ToggleMaximise() {
    thisWindow.ToggleMaximise();
  }
  function Minimise() {
    thisWindow.Minimise();
  }
  function UnMinimise() {
    thisWindow.UnMinimise();
  }
  function Restore() {
    thisWindow.Restore();
  }
  function Show2() {
    thisWindow.Show();
  }
  function Close() {
    thisWindow.Close();
  }
  function SetBackgroundColour(r, g, b, a) {
    thisWindow.SetBackgroundColour(r, g, b, a);
  }
  function SetResizable(resizable2) {
    thisWindow.SetResizable(resizable2);
  }
  function Width() {
    return thisWindow.Width();
  }
  function Height() {
    return thisWindow.Height();
  }
  function ZoomIn() {
    thisWindow.ZoomIn();
  }
  function ZoomOut() {
    thisWindow.ZoomOut();
  }
  function ZoomReset() {
    thisWindow.ZoomReset();
  }
  function GetZoomLevel() {
    return thisWindow.GetZoomLevel();
  }
  function SetZoomLevel(zoomLevel) {
    thisWindow.SetZoomLevel(zoomLevel);
  }

  // desktop/@wailsio/runtime/src/wml.js
  var wml_exports = {};
  __export(wml_exports, {
    Reload: () => Reload
  });

  // desktop/@wailsio/runtime/src/events.js
  var events_exports = {};
  __export(events_exports, {
    Emit: () => Emit,
    Off: () => Off,
    OffAll: () => OffAll,
    On: () => On,
    OnMultiple: () => OnMultiple,
    Once: () => Once,
    Types: () => Types,
    WailsEvent: () => WailsEvent,
    setup: () => setup
  });

  // desktop/@wailsio/runtime/src/event_types.js
  var EventTypes = {
    Windows: {
      SystemThemeChanged: "windows:SystemThemeChanged",
      APMPowerStatusChange: "windows:APMPowerStatusChange",
      APMSuspend: "windows:APMSuspend",
      APMResumeAutomatic: "windows:APMResumeAutomatic",
      APMResumeSuspend: "windows:APMResumeSuspend",
      APMPowerSettingChange: "windows:APMPowerSettingChange",
      ApplicationStarted: "windows:ApplicationStarted",
      WebViewNavigationCompleted: "windows:WebViewNavigationCompleted",
      WindowInactive: "windows:WindowInactive",
      WindowActive: "windows:WindowActive",
      WindowClickActive: "windows:WindowClickActive",
      WindowMaximise: "windows:WindowMaximise",
      WindowUnMaximise: "windows:WindowUnMaximise",
      WindowFullscreen: "windows:WindowFullscreen",
      WindowUnFullscreen: "windows:WindowUnFullscreen",
      WindowRestore: "windows:WindowRestore",
      WindowMinimise: "windows:WindowMinimise",
      WindowUnMinimise: "windows:WindowUnMinimise",
      WindowClose: "windows:WindowClose",
      WindowSetFocus: "windows:WindowSetFocus",
      WindowKillFocus: "windows:WindowKillFocus",
      WindowDragDrop: "windows:WindowDragDrop",
      WindowDragEnter: "windows:WindowDragEnter",
      WindowDragLeave: "windows:WindowDragLeave",
      WindowDragOver: "windows:WindowDragOver"
    },
    Mac: {
      ApplicationDidBecomeActive: "mac:ApplicationDidBecomeActive",
      ApplicationDidChangeBackingProperties: "mac:ApplicationDidChangeBackingProperties",
      ApplicationDidChangeEffectiveAppearance: "mac:ApplicationDidChangeEffectiveAppearance",
      ApplicationDidChangeIcon: "mac:ApplicationDidChangeIcon",
      ApplicationDidChangeOcclusionState: "mac:ApplicationDidChangeOcclusionState",
      ApplicationDidChangeScreenParameters: "mac:ApplicationDidChangeScreenParameters",
      ApplicationDidChangeStatusBarFrame: "mac:ApplicationDidChangeStatusBarFrame",
      ApplicationDidChangeStatusBarOrientation: "mac:ApplicationDidChangeStatusBarOrientation",
      ApplicationDidFinishLaunching: "mac:ApplicationDidFinishLaunching",
      ApplicationDidHide: "mac:ApplicationDidHide",
      ApplicationDidResignActiveNotification: "mac:ApplicationDidResignActiveNotification",
      ApplicationDidUnhide: "mac:ApplicationDidUnhide",
      ApplicationDidUpdate: "mac:ApplicationDidUpdate",
      ApplicationWillBecomeActive: "mac:ApplicationWillBecomeActive",
      ApplicationWillFinishLaunching: "mac:ApplicationWillFinishLaunching",
      ApplicationWillHide: "mac:ApplicationWillHide",
      ApplicationWillResignActive: "mac:ApplicationWillResignActive",
      ApplicationWillTerminate: "mac:ApplicationWillTerminate",
      ApplicationWillUnhide: "mac:ApplicationWillUnhide",
      ApplicationWillUpdate: "mac:ApplicationWillUpdate",
      ApplicationDidChangeTheme: "mac:ApplicationDidChangeTheme!",
      ApplicationShouldHandleReopen: "mac:ApplicationShouldHandleReopen!",
      WindowDidBecomeKey: "mac:WindowDidBecomeKey",
      WindowDidBecomeMain: "mac:WindowDidBecomeMain",
      WindowDidBeginSheet: "mac:WindowDidBeginSheet",
      WindowDidChangeAlpha: "mac:WindowDidChangeAlpha",
      WindowDidChangeBackingLocation: "mac:WindowDidChangeBackingLocation",
      WindowDidChangeBackingProperties: "mac:WindowDidChangeBackingProperties",
      WindowDidChangeCollectionBehavior: "mac:WindowDidChangeCollectionBehavior",
      WindowDidChangeEffectiveAppearance: "mac:WindowDidChangeEffectiveAppearance",
      WindowDidChangeOcclusionState: "mac:WindowDidChangeOcclusionState",
      WindowDidChangeOrderingMode: "mac:WindowDidChangeOrderingMode",
      WindowDidChangeScreen: "mac:WindowDidChangeScreen",
      WindowDidChangeScreenParameters: "mac:WindowDidChangeScreenParameters",
      WindowDidChangeScreenProfile: "mac:WindowDidChangeScreenProfile",
      WindowDidChangeScreenSpace: "mac:WindowDidChangeScreenSpace",
      WindowDidChangeScreenSpaceProperties: "mac:WindowDidChangeScreenSpaceProperties",
      WindowDidChangeSharingType: "mac:WindowDidChangeSharingType",
      WindowDidChangeSpace: "mac:WindowDidChangeSpace",
      WindowDidChangeSpaceOrderingMode: "mac:WindowDidChangeSpaceOrderingMode",
      WindowDidChangeTitle: "mac:WindowDidChangeTitle",
      WindowDidChangeToolbar: "mac:WindowDidChangeToolbar",
      WindowDidChangeVisibility: "mac:WindowDidChangeVisibility",
      WindowDidDeminiaturize: "mac:WindowDidDeminiaturize",
      WindowDidEndSheet: "mac:WindowDidEndSheet",
      WindowDidEnterFullScreen: "mac:WindowDidEnterFullScreen",
      WindowDidEnterVersionBrowser: "mac:WindowDidEnterVersionBrowser",
      WindowDidExitFullScreen: "mac:WindowDidExitFullScreen",
      WindowDidExitVersionBrowser: "mac:WindowDidExitVersionBrowser",
      WindowDidExpose: "mac:WindowDidExpose",
      WindowDidFocus: "mac:WindowDidFocus",
      WindowDidMiniaturize: "mac:WindowDidMiniaturize",
      WindowDidMove: "mac:WindowDidMove",
      WindowDidOrderOffScreen: "mac:WindowDidOrderOffScreen",
      WindowDidOrderOnScreen: "mac:WindowDidOrderOnScreen",
      WindowDidResignKey: "mac:WindowDidResignKey",
      WindowDidResignMain: "mac:WindowDidResignMain",
      WindowDidResize: "mac:WindowDidResize",
      WindowDidUpdate: "mac:WindowDidUpdate",
      WindowDidUpdateAlpha: "mac:WindowDidUpdateAlpha",
      WindowDidUpdateCollectionBehavior: "mac:WindowDidUpdateCollectionBehavior",
      WindowDidUpdateCollectionProperties: "mac:WindowDidUpdateCollectionProperties",
      WindowDidUpdateShadow: "mac:WindowDidUpdateShadow",
      WindowDidUpdateTitle: "mac:WindowDidUpdateTitle",
      WindowDidUpdateToolbar: "mac:WindowDidUpdateToolbar",
      WindowDidUpdateVisibility: "mac:WindowDidUpdateVisibility",
      WindowShouldClose: "mac:WindowShouldClose!",
      WindowWillBecomeKey: "mac:WindowWillBecomeKey",
      WindowWillBecomeMain: "mac:WindowWillBecomeMain",
      WindowWillBeginSheet: "mac:WindowWillBeginSheet",
      WindowWillChangeOrderingMode: "mac:WindowWillChangeOrderingMode",
      WindowWillClose: "mac:WindowWillClose",
      WindowWillDeminiaturize: "mac:WindowWillDeminiaturize",
      WindowWillEnterFullScreen: "mac:WindowWillEnterFullScreen",
      WindowWillEnterVersionBrowser: "mac:WindowWillEnterVersionBrowser",
      WindowWillExitFullScreen: "mac:WindowWillExitFullScreen",
      WindowWillExitVersionBrowser: "mac:WindowWillExitVersionBrowser",
      WindowWillFocus: "mac:WindowWillFocus",
      WindowWillMiniaturize: "mac:WindowWillMiniaturize",
      WindowWillMove: "mac:WindowWillMove",
      WindowWillOrderOffScreen: "mac:WindowWillOrderOffScreen",
      WindowWillOrderOnScreen: "mac:WindowWillOrderOnScreen",
      WindowWillResignMain: "mac:WindowWillResignMain",
      WindowWillResize: "mac:WindowWillResize",
      WindowWillUnfocus: "mac:WindowWillUnfocus",
      WindowWillUpdate: "mac:WindowWillUpdate",
      WindowWillUpdateAlpha: "mac:WindowWillUpdateAlpha",
      WindowWillUpdateCollectionBehavior: "mac:WindowWillUpdateCollectionBehavior",
      WindowWillUpdateCollectionProperties: "mac:WindowWillUpdateCollectionProperties",
      WindowWillUpdateShadow: "mac:WindowWillUpdateShadow",
      WindowWillUpdateTitle: "mac:WindowWillUpdateTitle",
      WindowWillUpdateToolbar: "mac:WindowWillUpdateToolbar",
      WindowWillUpdateVisibility: "mac:WindowWillUpdateVisibility",
      WindowWillUseStandardFrame: "mac:WindowWillUseStandardFrame",
      MenuWillOpen: "mac:MenuWillOpen",
      MenuDidOpen: "mac:MenuDidOpen",
      MenuDidClose: "mac:MenuDidClose",
      MenuWillSendAction: "mac:MenuWillSendAction",
      MenuDidSendAction: "mac:MenuDidSendAction",
      MenuWillHighlightItem: "mac:MenuWillHighlightItem",
      MenuDidHighlightItem: "mac:MenuDidHighlightItem",
      MenuWillDisplayItem: "mac:MenuWillDisplayItem",
      MenuDidDisplayItem: "mac:MenuDidDisplayItem",
      MenuWillAddItem: "mac:MenuWillAddItem",
      MenuDidAddItem: "mac:MenuDidAddItem",
      MenuWillRemoveItem: "mac:MenuWillRemoveItem",
      MenuDidRemoveItem: "mac:MenuDidRemoveItem",
      MenuWillBeginTracking: "mac:MenuWillBeginTracking",
      MenuDidBeginTracking: "mac:MenuDidBeginTracking",
      MenuWillEndTracking: "mac:MenuWillEndTracking",
      MenuDidEndTracking: "mac:MenuDidEndTracking",
      MenuWillUpdate: "mac:MenuWillUpdate",
      MenuDidUpdate: "mac:MenuDidUpdate",
      MenuWillPopUp: "mac:MenuWillPopUp",
      MenuDidPopUp: "mac:MenuDidPopUp",
      MenuWillSendActionToItem: "mac:MenuWillSendActionToItem",
      MenuDidSendActionToItem: "mac:MenuDidSendActionToItem",
      WebViewDidStartProvisionalNavigation: "mac:WebViewDidStartProvisionalNavigation",
      WebViewDidReceiveServerRedirectForProvisionalNavigation: "mac:WebViewDidReceiveServerRedirectForProvisionalNavigation",
      WebViewDidFinishNavigation: "mac:WebViewDidFinishNavigation",
      WebViewDidCommitNavigation: "mac:WebViewDidCommitNavigation",
      WindowFileDraggingEntered: "mac:WindowFileDraggingEntered",
      WindowFileDraggingPerformed: "mac:WindowFileDraggingPerformed",
      WindowFileDraggingExited: "mac:WindowFileDraggingExited"
    },
    Linux: {
      SystemThemeChanged: "linux:SystemThemeChanged",
      WindowLoadChanged: "linux:WindowLoadChanged",
      WindowDeleteEvent: "linux:WindowDeleteEvent",
      WindowFocusIn: "linux:WindowFocusIn",
      WindowFocusOut: "linux:WindowFocusOut",
      ApplicationStartup: "linux:ApplicationStartup"
    },
    Common: {
      ApplicationStarted: "common:ApplicationStarted",
      WindowMaximise: "common:WindowMaximise",
      WindowUnMaximise: "common:WindowUnMaximise",
      WindowFullscreen: "common:WindowFullscreen",
      WindowUnFullscreen: "common:WindowUnFullscreen",
      WindowRestore: "common:WindowRestore",
      WindowMinimise: "common:WindowMinimise",
      WindowUnMinimise: "common:WindowUnMinimise",
      WindowClosing: "common:WindowClosing",
      WindowZoom: "common:WindowZoom",
      WindowZoomIn: "common:WindowZoomIn",
      WindowZoomOut: "common:WindowZoomOut",
      WindowZoomReset: "common:WindowZoomReset",
      WindowFocus: "common:WindowFocus",
      WindowLostFocus: "common:WindowLostFocus",
      WindowShow: "common:WindowShow",
      WindowHide: "common:WindowHide",
      WindowDPIChanged: "common:WindowDPIChanged",
      WindowFilesDropped: "common:WindowFilesDropped",
      WindowRuntimeReady: "common:WindowRuntimeReady",
      ThemeChanged: "common:ThemeChanged"
    }
  };

  // desktop/@wailsio/runtime/src/events.js
  var Types = EventTypes;
  window._wails = window._wails || {};
  window._wails.dispatchWailsEvent = dispatchWailsEvent;
  var call7 = newRuntimeCallerWithID(objectNames.Events, "");
  var EmitMethod = 0;
  var eventListeners = /* @__PURE__ */ new Map();
  var Listener = class {
    constructor(eventName, callback, maxCallbacks) {
      this.eventName = eventName;
      this.maxCallbacks = maxCallbacks || -1;
      this.Callback = (data) => {
        callback(data);
        if (this.maxCallbacks === -1)
          return false;
        this.maxCallbacks -= 1;
        return this.maxCallbacks === 0;
      };
    }
  };
  var WailsEvent = class {
    constructor(name, data = null) {
      this.name = name;
      this.data = data;
    }
  };
  function setup() {
  }
  function dispatchWailsEvent(event) {
    let listeners = eventListeners.get(event.name);
    if (listeners) {
      let toRemove = listeners.filter((listener) => {
        let remove = listener.Callback(event);
        if (remove)
          return true;
      });
      if (toRemove.length > 0) {
        listeners = listeners.filter((l) => !toRemove.includes(l));
        if (listeners.length === 0)
          eventListeners.delete(event.name);
        else
          eventListeners.set(event.name, listeners);
      }
    }
  }
  function OnMultiple(eventName, callback, maxCallbacks) {
    let listeners = eventListeners.get(eventName) || [];
    const thisListener = new Listener(eventName, callback, maxCallbacks);
    listeners.push(thisListener);
    eventListeners.set(eventName, listeners);
    return () => listenerOff(thisListener);
  }
  function On(eventName, callback) {
    return OnMultiple(eventName, callback, -1);
  }
  function Once(eventName, callback) {
    return OnMultiple(eventName, callback, 1);
  }
  function listenerOff(listener) {
    const eventName = listener.eventName;
    let listeners = eventListeners.get(eventName).filter((l) => l !== listener);
    if (listeners.length === 0)
      eventListeners.delete(eventName);
    else
      eventListeners.set(eventName, listeners);
  }
  function Off(eventName, ...additionalEventNames) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach((eventName2) => eventListeners.delete(eventName2));
  }
  function OffAll() {
    eventListeners.clear();
  }
  function Emit(event) {
    return call7(EmitMethod, event);
  }

  // desktop/@wailsio/runtime/src/dialogs.js
  var dialogs_exports = {};
  __export(dialogs_exports, {
    Error: () => Error2,
    Info: () => Info,
    OpenFile: () => OpenFile,
    Question: () => Question,
    SaveFile: () => SaveFile,
    Warning: () => Warning
  });
  window._wails = window._wails || {};
  window._wails.dialogErrorCallback = dialogErrorCallback;
  window._wails.dialogResultCallback = dialogResultCallback;
  var DialogInfo = 0;
  var DialogWarning = 1;
  var DialogError = 2;
  var DialogQuestion = 3;
  var DialogOpenFile = 4;
  var DialogSaveFile = 5;
  var call8 = newRuntimeCallerWithID(objectNames.Dialog, "");
  var dialogResponses = /* @__PURE__ */ new Map();
  function generateID() {
    let result;
    do {
      result = nanoid();
    } while (dialogResponses.has(result));
    return result;
  }
  function dialog(type, options = {}) {
    const id = generateID();
    options["dialog-id"] = id;
    return new Promise((resolve, reject) => {
      dialogResponses.set(id, { resolve, reject });
      call8(type, options).catch((error) => {
        reject(error);
        dialogResponses.delete(id);
      });
    });
  }
  function dialogResultCallback(id, data, isJSON) {
    let p = dialogResponses.get(id);
    if (p) {
      if (isJSON) {
        p.resolve(JSON.parse(data));
      } else {
        p.resolve(data);
      }
      dialogResponses.delete(id);
    }
  }
  function dialogErrorCallback(id, message) {
    let p = dialogResponses.get(id);
    if (p) {
      p.reject(message);
      dialogResponses.delete(id);
    }
  }
  var Info = (options) => dialog(DialogInfo, options);
  var Warning = (options) => dialog(DialogWarning, options);
  var Error2 = (options) => dialog(DialogError, options);
  var Question = (options) => dialog(DialogQuestion, options);
  var OpenFile = (options) => dialog(DialogOpenFile, options);
  var SaveFile = (options) => dialog(DialogSaveFile, options);

  // desktop/@wailsio/runtime/src/wml.js
  function sendEvent(eventName, data = null) {
    let event = new WailsEvent(eventName, data);
    Emit(event);
  }
  function addWMLEventListeners() {
    const elements = document.querySelectorAll("[wml-event]");
    elements.forEach(function(element) {
      const eventType = element.getAttribute("wml-event");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Detached: false, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              sendEvent(eventType);
            }
          });
          return;
        }
        sendEvent(eventType);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function callWindowMethod(windowName, method) {
    let targetWindow = Get(windowName);
    let methodMap = WindowMethods(targetWindow);
    if (!methodMap.has(method)) {
      console.log("Window method " + method + " not found");
    }
    try {
      methodMap.get(method)();
    } catch (e) {
      console.error("Error calling window method '" + method + "': " + e);
    }
  }
  function addWMLWindowListeners() {
    const elements = document.querySelectorAll("[wml-window]");
    elements.forEach(function(element) {
      const windowMethod = element.getAttribute("wml-window");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      const targetWindow = element.getAttribute("wml-target-window") || "";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              callWindowMethod(targetWindow, windowMethod);
            }
          });
          return;
        }
        callWindowMethod(targetWindow, windowMethod);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function addWMLOpenBrowserListener() {
    const elements = document.querySelectorAll("[wml-openurl]");
    elements.forEach(function(element) {
      const url = element.getAttribute("wml-openurl");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              void OpenURL(url);
            }
          });
          return;
        }
        void OpenURL(url);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function Reload() {
    addWMLEventListeners();
    addWMLWindowListeners();
    addWMLOpenBrowserListener();
  }
  function WindowMethods(targetWindow) {
    let result = /* @__PURE__ */ new Map();
    for (let method in targetWindow) {
      if (typeof targetWindow[method] === "function") {
        result.set(method, targetWindow[method]);
      }
    }
    return result;
  }

  // desktop/@wailsio/runtime/src/calls.js
  var calls_exports = {};
  __export(calls_exports, {
    ByID: () => ByID,
    ByName: () => ByName,
    Call: () => Call,
    Plugin: () => Plugin
  });
  window._wails = window._wails || {};
  window._wails.callResultHandler = resultHandler;
  window._wails.callErrorHandler = errorHandler;
  var CallBinding = 0;
  var call9 = newRuntimeCallerWithID(objectNames.Call, "");
  var cancelCall = newRuntimeCallerWithID(objectNames.CancelCall, "");
  var callResponses = /* @__PURE__ */ new Map();
  function generateID2() {
    let result;
    do {
      result = nanoid();
    } while (callResponses.has(result));
    return result;
  }
  function resultHandler(id, data, isJSON) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
      promiseHandler.resolve(isJSON ? JSON.parse(data) : data);
    }
  }
  function errorHandler(id, message) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
      promiseHandler.reject(message);
    }
  }
  function getAndDeleteResponse(id) {
    const response = callResponses.get(id);
    callResponses.delete(id);
    return response;
  }
  function callBinding(type, options = {}) {
    const id = generateID2();
    const doCancel = () => {
      cancelCall(type, { "call-id": id });
    };
    var queuedCancel = false, callRunning = false;
    var p = new Promise((resolve, reject) => {
      options["call-id"] = id;
      callResponses.set(id, { resolve, reject });
      call9(type, options).then((_) => {
        callRunning = true;
        if (queuedCancel) {
          doCancel();
        }
      }).catch((error) => {
        reject(error);
        callResponses.delete(id);
      });
    });
    p.cancel = () => {
      if (callRunning) {
        doCancel();
      } else {
        queuedCancel = true;
      }
    };
    return p;
  }
  function Call(options) {
    return callBinding(CallBinding, options);
  }
  function ByName(name, ...args) {
    if (typeof name !== "string" || name.split(".").length !== 3) {
      throw new Error("CallByName requires a string in the format 'package.struct.method'");
    }
    let [packageName, structName, methodName] = name.split(".");
    return callBinding(CallBinding, {
      packageName,
      structName,
      methodName,
      args
    });
  }
  function ByID(methodID, ...args) {
    return callBinding(CallBinding, {
      methodID,
      args
    });
  }
  function Plugin(pluginName, methodName, ...args) {
    return callBinding(CallBinding, {
      packageName: "wails-plugins",
      structName: pluginName,
      methodName,
      args
    });
  }

  // desktop/compiled/main.js
  window._wails = window._wails || {};
  window._wails.invoke = invoke;
  window.wails = window.wails || {};
  window.wails.Application = application_exports;
  window.wails.Browser = browser_exports;
  window.wails.Call = calls_exports;
  window.wails.Clipboard = clipboard_exports;
  window.wails.Dialogs = dialogs_exports;
  window.wails.Events = events_exports;
  window.wails.Flags = flags_exports;
  window.wails.Screens = screens_exports;
  window.wails.System = system_exports;
  window.wails.Window = window_exports;
  window.wails.WML = wml_exports;
  var isReady = false;
  document.addEventListener("DOMContentLoaded", function() {
    isReady = true;
    window._wails.invoke("wails:runtime:ready");
    if (true) {
      debugLog("Wails Runtime Loaded");
    }
  });
  function whenReady(fn) {
    if (isReady || document.readyState === "complete") {
      fn();
    } else {
      document.addEventListener("DOMContentLoaded", fn);
    }
  }
  whenReady(() => {
    Reload();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb24uanMiLCAiLi4vLi4vLi4vcnVudGltZS9ub2RlX21vZHVsZXMvbmFub2lkL25vbi1zZWN1cmUvaW5kZXguanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZC5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvd21sLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudHMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50X3R5cGVzLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kaWFsb2dzLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLyoqXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZSkge1xuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgIGNvbnNvbGUubG9nKFxuICAgICAgICAnJWMgd2FpbHMzICVjICcgKyBtZXNzYWdlICsgJyAnLFxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjMDA5OTAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogMHB4IDNweCAzcHggMHB4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJ1xuICAgICk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24sICcnKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcbn1cblxuLyoqXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXG4gKlxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxufVxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG4vKipcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cbiAqXG4gKiBAcGFyYW0ge09iamVjdH0gb2JqZWN0IC0gVGhlIG9iamVjdCBvbiB3aGljaCB0aGUgbWV0aG9kIGlzIHRvIGJlIGludm9rZWQuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyB7RnVuY3Rpb259IEEgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24gdGhhdCB0YWtlcyB0aGUgbWV0aG9kIG5hbWUgYW5kIG9wdGlvbmFsbHkgYXJndW1lbnRzIGFuZCBpbnZva2VzIHRoZSBtZXRob2Qgd2l0aGluIHRoZSBzcGVjaWZpZWQgd2luZG93IGNvbnRleHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG4vKipcbiAqIENyZWF0ZXMgYSBuZXcgcnVudGltZSBjYWxsZXIgd2l0aCBzcGVjaWZpZWQgSUQuXG4gKlxuICogQHBhcmFtIHtvYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gKiBAcmV0dXJuIHtGdW5jdGlvbn0gLSBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3QsIHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICBpZiggbWV0aG9kICkge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIH1cbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcblxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm9iamVjdFwiLCBvYmplY3RJRCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkJyb3dzZXIsICcnKTtcbmNvbnN0IEJyb3dzZXJPcGVuVVJMID0gMDtcblxuLyoqXG4gKiBPcGVuIGEgYnJvd3NlciB3aW5kb3cgdG8gdGhlIGdpdmVuIFVSTFxuICogQHBhcmFtIHtzdHJpbmd9IHVybCAtIFRoZSBVUkwgdG8gb3BlblxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5VUkwodXJsKSB7XG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmx9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCwgJycpO1xuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XG5jb25zdCBDbGlwYm9hcmRUZXh0ID0gMTtcblxuLyoqXG4gKiBTZXRzIHRoZSB0ZXh0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cbiAqIEByZXR1cm4ge1Byb21pc2V9IC0gQSBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2hlbiB0aGUgb3BlcmF0aW9uIGlzIHN1Y2Nlc3NmdWwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRTZXRUZXh0LCB7dGV4dH0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU3lzdGVtLCAnJyk7XG5jb25zdCBzeXN0ZW1Jc0RhcmtNb2RlID0gMDtcbmNvbnN0IGVudmlyb25tZW50ID0gMTtcblxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2cpIHtcbiAgICBpZih3aW5kb3cuY2hyb21lKSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UobXNnKTtcbiAgICB9XG4gICAgcmV0dXJuIHdpbmRvdy53ZWJraXQubWVzc2FnZUhhbmRsZXJzLmV4dGVybmFsLnBvc3RNZXNzYWdlKG1zZyk7XG59XG5cbi8qKlxuICogQGZ1bmN0aW9uXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxuICogQHJldHVybnMge1Byb21pc2U8Ym9vbGVhbj59IC0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYSBib29sZWFuIHZhbHVlIGluZGljYXRpbmcgaWYgdGhlIHN5c3RlbSBpcyBpbiBkYXJrIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCkge1xuICAgIHJldHVybiBjYWxsKHN5c3RlbUlzRGFya01vZGUpO1xufVxuXG4vKipcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxuICpcbiAqIEBhc3luY1xuICogQGZ1bmN0aW9uIENhcGFiaWxpdGllc1xuICogQHJldHVybnMge1Byb21pc2U8T2JqZWN0Pn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhcGFiaWxpdGllcygpIHtcbiAgICBsZXQgcmVzcG9uc2UgPSBmZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIik7XG4gICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbn1cblxuLyoqXG4gKiBAdHlwZWRlZiB7b2JqZWN0fSBFbnZpcm9ubWVudEluZm9cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBPUyAtIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBBcmNoIC0gVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLlxuICovXG5cbi8qKlxuICogQGZ1bmN0aW9uXG4gKiBSZXRyaWV2ZXMgZW52aXJvbm1lbnQgZGV0YWlscy5cbiAqIEByZXR1cm5zIHtQcm9taXNlPEVudmlyb25tZW50SW5mbz59IC0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpIHtcbiAgICByZXR1cm4gY2FsbChlbnZpcm9ubWVudCk7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cbiAqXG4gKiBAcmV0dXJuIHtib29sZWFufSBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MsIG90aGVyd2lzZSBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzV2luZG93cygpIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJ3aW5kb3dzXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXG4gKlxuICogQHJldHVybnMge2Jvb2xlYW59IFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImxpbnV4XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImRhcndpblwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhbWQ2NFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxuICpcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gLSBSZXR1cm5zIHRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZSwgb3RoZXJ3aXNlIHJldHVybnMgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtNjRcIjtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIElzRGVidWcoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuRGVidWcgPT09IHRydWU7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5pbXBvcnQge0lzRGVidWd9IGZyb20gXCIuL3N5c3RlbVwiO1xuXG4vLyBzZXR1cFxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUsICcnKTtcbmNvbnN0IENvbnRleHRNZW51T3BlbiA9IDA7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxuICAgIGxldCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xuICAgIGxldCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudVwiKTtcbiAgICBjdXN0b21Db250ZXh0TWVudSA9IGN1c3RvbUNvbnRleHRNZW51ID8gY3VzdG9tQ29udGV4dE1lbnUudHJpbSgpIDogXCJcIjtcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgbGV0IGN1c3RvbUNvbnRleHRNZW51RGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGN1c3RvbUNvbnRleHRNZW51RGF0YSk7XG4gICAgICAgIHJldHVyblxuICAgIH1cblxuICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpO1xufVxuXG5cbi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcblxuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxuICAgIGlmIChJc0RlYnVnKCkpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICBjb25zdCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xuICAgIGNvbnN0IGNvbXB1dGVkU3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcbiAgICBjb25zdCBkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24gPSBjb21wdXRlZFN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpO1xuICAgIHN3aXRjaCAoZGVmYXVsdENvbnRleHRNZW51QWN0aW9uKSB7XG4gICAgICAgIGNhc2UgXCJzaG93XCI6XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGNhc2UgXCJoaWRlXCI6XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcbiAgICAgICAgICAgIGlmIChlbGVtZW50LmlzQ29udGVudEVkaXRhYmxlKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXG4gICAgICAgICAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XG4gICAgICAgICAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSAoc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMClcbiAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcbiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xuICAgICAgICAgICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gZWxlbWVudCkge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRhZ25hbWUgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbiAgICAgICAgICAgIGlmIChlbGVtZW50LnRhZ05hbWUgPT09IFwiSU5QVVRcIiB8fCBlbGVtZW50LnRhZ05hbWUgPT09IFwiVEVYVEFSRUFcIikge1xuICAgICAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24gfHwgKCFlbGVtZW50LnJlYWRPbmx5ICYmICFlbGVtZW50LmRpc2FibGVkKSkge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIFJldHJpZXZlcyB0aGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5IGZyb20gdGhlIGZsYWcgbWFwLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBrZXlTdHJpbmcgLSBUaGUga2V5IHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBmb3IuXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleVN0cmluZykge1xuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleVN0cmluZ107XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJVbmFibGUgdG8gcmV0cmlldmUgZmxhZyAnXCIgKyBrZXlTdHJpbmcgKyBcIic6IFwiICsgZSk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7aW52b2tlLCBJc1dpbmRvd3N9IGZyb20gXCIuL3N5c3RlbVwiO1xuaW1wb3J0IHtHZXRGbGFnfSBmcm9tIFwiLi9mbGFnc1wiO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9IHNldFJlc2l6YWJsZTtcbndpbmRvdy5fd2FpbHMuZW5kRHJhZyA9IGVuZERyYWc7XG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgb25Nb3VzZURvd24pO1xud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIG9uTW91c2VNb3ZlKTtcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgb25Nb3VzZVVwKTtcblxuXG5sZXQgc2hvdWxkRHJhZyA9IGZhbHNlO1xubGV0IHJlc2l6ZUVkZ2UgPSBudWxsO1xubGV0IHJlc2l6YWJsZSA9IGZhbHNlO1xubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcblxuZnVuY3Rpb24gZHJhZ1Rlc3QoZSkge1xuICAgIGxldCB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2Via2l0LWFwcC1yZWdpb25cIik7XG4gICAgaWYgKCF2YWwgfHwgdmFsID09PSBcIlwiIHx8IHZhbC50cmltKCkgIT09IFwiZHJhZ1wiIHx8IGUuYnV0dG9ucyAhPT0gMSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcbn1cblxuZnVuY3Rpb24gc2V0UmVzaXphYmxlKHZhbHVlKSB7XG4gICAgcmVzaXphYmxlID0gdmFsdWU7XG59XG5cbmZ1bmN0aW9uIGVuZERyYWcoKSB7XG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiB0ZXN0UmVzaXplKCkge1xuICAgIGlmKCByZXNpemVFZGdlICkge1xuICAgICAgICBpbnZva2UoYHJlc2l6ZToke3Jlc2l6ZUVkZ2V9YCk7XG4gICAgICAgIHJldHVybiB0cnVlXG4gICAgfVxuICAgIHJldHVybiBmYWxzZTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xuICAgIGlmKElzV2luZG93cygpICYmIHRlc3RSZXNpemUoKSB8fCBkcmFnVGVzdChlKSkge1xuICAgICAgICBzaG91bGREcmFnID0gISFpc1ZhbGlkRHJhZyhlKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGlzVmFsaWREcmFnKGUpIHtcbiAgICAvLyBJZ25vcmUgZHJhZyBvbiBzY3JvbGxiYXJzXG4gICAgcmV0dXJuICEoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZVVwKGUpIHtcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgIGVuZERyYWcoKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IgPSBkZWZhdWx0Q3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvcjtcbiAgICByZXNpemVFZGdlID0gY3Vyc29yO1xufVxuXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShlKSB7XG4gICAgc2hvdWxkRHJhZyA9IGNoZWNrRHJhZyhlKTtcbiAgICBpZiAoSXNXaW5kb3dzKCkgJiYgcmVzaXphYmxlKSB7XG4gICAgICAgIGhhbmRsZVJlc2l6ZShlKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNoZWNrRHJhZyhlKSB7XG4gICAgbGV0IG1vdXNlUHJlc3NlZCA9IGUuYnV0dG9ucyAhPT0gdW5kZWZpbmVkID8gZS5idXR0b25zIDogZS53aGljaDtcbiAgICBpZihzaG91bGREcmFnICYmIG1vdXNlUHJlc3NlZCA+IDApIHtcbiAgICAgICAgaW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbiAgICByZXR1cm4gc2hvdWxkRHJhZztcbn1cblxuZnVuY3Rpb24gaGFuZGxlUmVzaXplKGUpIHtcbiAgICBsZXQgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcbiAgICBsZXQgcmVzaXplSGFuZGxlV2lkdGggPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZVdpZHRoXCIpIHx8IDU7XG5cbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXNcbiAgICBsZXQgY29ybmVyRXh0cmEgPSBHZXRGbGFnKFwicmVzaXplQ29ybmVyRXh0cmFcIikgfHwgMTA7XG5cbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgbGV0IHRvcEJvcmRlciA9IGUuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXJzXG4gICAgbGV0IHJpZ2h0Q29ybmVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgbGV0IGxlZnRDb3JuZXIgPSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgbGV0IHRvcENvcm5lciA9IGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG4gICAgbGV0IGJvdHRvbUNvcm5lciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG5cbiAgICAvLyBJZiB3ZSBhcmVuJ3Qgb24gYW4gZWRnZSwgYnV0IHdlcmUsIHJlc2V0IHRoZSBjdXJzb3IgdG8gZGVmYXVsdFxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9XG4gICAgLy8gQWRqdXN0ZWQgZm9yIGNvcm5lciBhcmVhc1xuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gUG9zaXRpb25cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZS5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBZIC0gVGhlIFkgY29vcmRpbmF0ZS5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNpemVcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIHdpZHRoLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgaGVpZ2h0LlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBSZWN0XG4gKiBAcHJvcGVydHkge251bWJlcn0gWCAtIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHRvcC1sZWZ0IGNvcm5lci5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBZIC0gVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgdG9wLWxlZnQgY29ybmVyLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuXG4gKiBAcHJvcGVydHkge251bWJlcn0gSGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7KCdaZXJvJ3wnTmluZXR5J3wnT25lRWlnaHR5J3wnVHdvU2V2ZW50eScpfSBSb3RhdGlvblxuICogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uIENhbiBiZSBvbmUgb2YgJ1plcm8nLCAnTmluZXR5JywgJ09uZUVpZ2h0eScsICdUd29TZXZlbnR5Jy5cbiAqL1xuXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gU2NyZWVuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gSWQgLSBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBOYW1lIC0gSHVtYW4gcmVhZGFibGUgbmFtZSBvZiB0aGUgc2NyZWVuLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFNjYWxlIC0gVGhlIHJlc29sdXRpb24gc2NhbGUgb2YgdGhlIHNjcmVlbi4gMSA9IHN0YW5kYXJkIHJlc29sdXRpb24sIDIgPSBoaWdoIChSZXRpbmEpLCBldGMuXG4gKiBAcHJvcGVydHkge1Bvc2l0aW9ufSBQb3NpdGlvbiAtIENvbnRhaW5zIHRoZSBYIGFuZCBZIGNvb3JkaW5hdGVzIG9mIHRoZSBzY3JlZW4ncyBwb3NpdGlvbi5cbiAqIEBwcm9wZXJ0eSB7U2l6ZX0gU2l6ZSAtIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uXG4gKiBAcHJvcGVydHkge1JlY3R9IEJvdW5kcyAtIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC5cbiAqIEBwcm9wZXJ0eSB7UmVjdH0gV29ya0FyZWEgLSBDb250YWlucyB0aGUgYXJlYSBvZiB0aGUgc2NyZWVuIHRoYXQgaXMgYWN0dWFsbHkgdXNhYmxlIChleGNsdWRpbmcgdGFza2JhciBhbmQgb3RoZXIgc3lzdGVtIFVJKS5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gSXNQcmltYXJ5IC0gVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uXG4gKiBAcHJvcGVydHkge1JvdGF0aW9ufSBSb3RhdGlvbiAtIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLlxuICovXG5cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU2NyZWVucywgJycpO1xuXG5jb25zdCBnZXRBbGwgPSAwO1xuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XG5jb25zdCBnZXRDdXJyZW50ID0gMjtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChnZXRBbGwpO1xufVxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKGdldEN1cnJlbnQpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vLyBJbXBvcnQgc2NyZWVuIGpzZG9jIGRlZmluaXRpb24gZnJvbSAuL3NjcmVlbnMuanNcbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vc2NyZWVuc1wiKS5TY3JlZW59IFNjcmVlblxuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuY29uc3QgY2VudGVyID0gMDtcbmNvbnN0IHNldFRpdGxlID0gMTtcbmNvbnN0IGZ1bGxzY3JlZW4gPSAyO1xuY29uc3QgdW5GdWxsc2NyZWVuID0gMztcbmNvbnN0IHNldFNpemUgPSA0O1xuY29uc3Qgc2l6ZSA9IDU7XG5jb25zdCBzZXRNYXhTaXplID0gNjtcbmNvbnN0IHNldE1pblNpemUgPSA3O1xuY29uc3Qgc2V0QWx3YXlzT25Ub3AgPSA4O1xuY29uc3Qgc2V0UmVsYXRpdmVQb3NpdGlvbiA9IDk7XG5jb25zdCByZWxhdGl2ZVBvc2l0aW9uID0gMTA7XG5jb25zdCBzY3JlZW4gPSAxMTtcbmNvbnN0IGhpZGUgPSAxMjtcbmNvbnN0IG1heGltaXNlID0gMTM7XG5jb25zdCB1bk1heGltaXNlID0gMTQ7XG5jb25zdCB0b2dnbGVNYXhpbWlzZSA9IDE1O1xuY29uc3QgbWluaW1pc2UgPSAxNjtcbmNvbnN0IHVuTWluaW1pc2UgPSAxNztcbmNvbnN0IHJlc3RvcmUgPSAxODtcbmNvbnN0IHNob3cgPSAxOTtcbmNvbnN0IGNsb3NlID0gMjA7XG5jb25zdCBzZXRCYWNrZ3JvdW5kQ29sb3VyID0gMjE7XG5jb25zdCBzZXRSZXNpemFibGUgPSAyMjtcbmNvbnN0IHdpZHRoID0gMjM7XG5jb25zdCBoZWlnaHQgPSAyNDtcbmNvbnN0IHpvb21JbiA9IDI1O1xuY29uc3Qgem9vbU91dCA9IDI2O1xuY29uc3Qgem9vbVJlc2V0ID0gMjc7XG5jb25zdCBnZXRab29tTGV2ZWwgPSAyODtcbmNvbnN0IHNldFpvb21MZXZlbCA9IDI5O1xuXG5jb25zdCB0aGlzV2luZG93ID0gR2V0KCcnKTtcblxuZnVuY3Rpb24gY3JlYXRlV2luZG93KGNhbGwpIHtcbiAgICByZXR1cm4ge1xuICAgICAgICBHZXQ6ICh3aW5kb3dOYW1lKSA9PiBjcmVhdGVXaW5kb3cobmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIHdpbmRvd05hbWUpKSxcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKGNlbnRlciksXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IGNhbGwoc2V0VGl0bGUsIHt0aXRsZX0pLFxuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKGZ1bGxzY3JlZW4pLFxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwodW5GdWxsc2NyZWVuKSxcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcbiAgICAgICAgU2l6ZTogKCkgPT4gY2FsbChzaXplKSxcbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWF4U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcbiAgICAgICAgU2V0TWluU2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWluU2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcbiAgICAgICAgU2V0QWx3YXlzT25Ub3A6IChvblRvcCkgPT4gY2FsbChzZXRBbHdheXNPblRvcCwge2Fsd2F5c09uVG9wOiBvblRvcH0pLFxuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbChzZXRSZWxhdGl2ZVBvc2l0aW9uLCB7eCwgeX0pLFxuICAgICAgICBSZWxhdGl2ZVBvc2l0aW9uOiAoKSA9PiBjYWxsKHJlbGF0aXZlUG9zaXRpb24pLFxuICAgICAgICBTY3JlZW46ICgpID0+IGNhbGwoc2NyZWVuKSxcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbChoaWRlKSxcbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IGNhbGwobWF4aW1pc2UpLFxuICAgICAgICBVbk1heGltaXNlOiAoKSA9PiBjYWxsKHVuTWF4aW1pc2UpLFxuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gY2FsbCh0b2dnbGVNYXhpbWlzZSksXG4gICAgICAgIE1pbmltaXNlOiAoKSA9PiBjYWxsKG1pbmltaXNlKSxcbiAgICAgICAgVW5NaW5pbWlzZTogKCkgPT4gY2FsbCh1bk1pbmltaXNlKSxcbiAgICAgICAgUmVzdG9yZTogKCkgPT4gY2FsbChyZXN0b3JlKSxcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbChzaG93KSxcbiAgICAgICAgQ2xvc2U6ICgpID0+IGNhbGwoY2xvc2UpLFxuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gY2FsbChzZXRCYWNrZ3JvdW5kQ29sb3VyLCB7ciwgZywgYiwgYX0pLFxuICAgICAgICBTZXRSZXNpemFibGU6IChyZXNpemFibGUpID0+IGNhbGwoc2V0UmVzaXphYmxlLCB7cmVzaXphYmxlfSksXG4gICAgICAgIFdpZHRoOiAoKSA9PiBjYWxsKHdpZHRoKSxcbiAgICAgICAgSGVpZ2h0OiAoKSA9PiBjYWxsKGhlaWdodCksXG4gICAgICAgIFpvb21JbjogKCkgPT4gY2FsbCh6b29tSW4pLFxuICAgICAgICBab29tT3V0OiAoKSA9PiBjYWxsKHpvb21PdXQpLFxuICAgICAgICBab29tUmVzZXQ6ICgpID0+IGNhbGwoem9vbVJlc2V0KSxcbiAgICAgICAgR2V0Wm9vbUxldmVsOiAoKSA9PiBjYWxsKGdldFpvb21MZXZlbCksXG4gICAgICAgIFNldFpvb21MZXZlbDogKHpvb21MZXZlbCkgPT4gY2FsbChzZXRab29tTGV2ZWwsIHt6b29tTGV2ZWx9KSxcbiAgICB9O1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cbiAqIEByZXR1cm4ge09iamVjdH0gLSBUaGUgc3BlY2lmaWVkIHdpbmRvdyBvYmplY3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXQod2luZG93TmFtZSkge1xuICAgIHJldHVybiBjcmVhdGVXaW5kb3cobmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIHdpbmRvd05hbWUpKTtcbn1cblxuLyoqXG4gKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENlbnRlcigpIHtcbiAgICB0aGlzV2luZG93LkNlbnRlcigpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gKiBAcGFyYW0ge3N0cmluZ30gdGl0bGUgLSBUaGUgdGl0bGUgdG8gc2V0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGl0bGUodGl0bGUpIHtcbiAgICB0aGlzV2luZG93LlNldFRpdGxlKHRpdGxlKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSB3aW5kb3cgdG8gZnVsbHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEZ1bGxzY3JlZW4oKSB7XG4gICAgdGhpc1dpbmRvdy5GdWxsc2NyZWVuKCk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgdGhpc1dpbmRvdy5TZXRTaXplKHdpZHRoLCBoZWlnaHQpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNpemUoKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuU2l6ZSgpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIG1heGltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldE1heFNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHRoaXNXaW5kb3cuU2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCAtIFRoZSBtaW5pbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB0aGlzV2luZG93LlNldE1pblNpemUod2lkdGgsIGhlaWdodCk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgd2luZG93IHRvIGFsd2F5cyBiZSBvbiB0b3AuXG4gKiBAcGFyYW0ge2Jvb2xlYW59IG9uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBhbHdheXMgYmUgb24gdG9wLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0QWx3YXlzT25Ub3Aob25Ub3ApIHtcbiAgICB0aGlzV2luZG93LlNldEFsd2F5c09uVG9wKG9uVG9wKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSB3aW5kb3cncyBwb3NpdGlvbi5cbiAqIEBwYXJhbSB7bnVtYmVyfSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgd2luZG93J3MgcG9zaXRpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRSZWxhdGl2ZVBvc2l0aW9uKHgsIHkpIHtcbiAgICB0aGlzV2luZG93LlNldFJlbGF0aXZlUG9zaXRpb24oeCwgeSk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbGF0aXZlUG9zaXRpb24oKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuUmVsYXRpdmVQb3NpdGlvbigpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTY3JlZW4oKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuU2NyZWVuKCk7XG59XG5cbi8qKlxuICogSGlkZXMgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgdGhpc1dpbmRvdy5IaWRlKCk7XG59XG5cbi8qKlxuICogTWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBNYXhpbWlzZSgpIHtcbiAgICB0aGlzV2luZG93Lk1heGltaXNlKCk7XG59XG5cbi8qKlxuICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBVbk1heGltaXNlKCkge1xuICAgIHRoaXNXaW5kb3cuVW5NYXhpbWlzZSgpO1xufVxuXG4vKipcbiAqIFRvZ2dsZXMgdGhlIG1heGltaXNhdGlvbiBvZiB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gVG9nZ2xlTWF4aW1pc2UoKSB7XG4gICAgdGhpc1dpbmRvdy5Ub2dnbGVNYXhpbWlzZSgpO1xufVxuXG4vKipcbiAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gTWluaW1pc2UoKSB7XG4gICAgdGhpc1dpbmRvdy5NaW5pbWlzZSgpO1xufVxuXG4vKipcbiAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gVW5NaW5pbWlzZSgpIHtcbiAgICB0aGlzV2luZG93LlVuTWluaW1pc2UoKTtcbn1cblxuLyoqXG4gKiBSZXN0b3JlcyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVzdG9yZSgpIHtcbiAgICB0aGlzV2luZG93LlJlc3RvcmUoKTtcbn1cblxuLyoqXG4gKiBTaG93cyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICB0aGlzV2luZG93LlNob3coKTtcbn1cblxuLyoqXG4gKiBDbG9zZXMgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENsb3NlKCkge1xuICAgIHRoaXNXaW5kb3cuQ2xvc2UoKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHIgLSBUaGUgcmVkIGNvbXBvbmVudCBvZiB0aGUgY29sb3VyLlxuICogQHBhcmFtIHtudW1iZXJ9IGcgLSBUaGUgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXG4gKiBAcGFyYW0ge251bWJlcn0gYiAtIFRoZSBibHVlIGNvbXBvbmVudCBvZiB0aGUgY29sb3VyLlxuICogQHBhcmFtIHtudW1iZXJ9IGEgLSBUaGUgYWxwaGEgY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRCYWNrZ3JvdW5kQ29sb3VyKHIsIGcsIGIsIGEpIHtcbiAgICB0aGlzV2luZG93LlNldEJhY2tncm91bmRDb2xvdXIociwgZywgYiwgYSk7XG59XG5cbi8qKlxuICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICogQHBhcmFtIHtib29sZWFufSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFJlc2l6YWJsZShyZXNpemFibGUpIHtcbiAgICB0aGlzV2luZG93LlNldFJlc2l6YWJsZShyZXNpemFibGUpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaWR0aCgpIHtcbiAgICByZXR1cm4gdGhpc1dpbmRvdy5XaWR0aCgpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gSGVpZ2h0KCkge1xuICAgIHJldHVybiB0aGlzV2luZG93LkhlaWdodCgpO1xufVxuXG4vKipcbiAqIFpvb21zIGluIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBab29tSW4oKSB7XG4gICAgdGhpc1dpbmRvdy5ab29tSW4oKTtcbn1cblxuLyoqXG4gKiBab29tcyBvdXQgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFpvb21PdXQoKSB7XG4gICAgdGhpc1dpbmRvdy5ab29tT3V0KCk7XG59XG5cbi8qKlxuICogUmVzZXRzIHRoZSB6b29tIG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBab29tUmVzZXQoKSB7XG4gICAgdGhpc1dpbmRvdy5ab29tUmVzZXQoKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRab29tTGV2ZWwoKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuR2V0Wm9vbUxldmVsKCk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHpvb21MZXZlbCAtIFRoZSB6b29tIGxldmVsIHRvIHNldC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFpvb21MZXZlbCh6b29tTGV2ZWwpIHtcbiAgICB0aGlzV2luZG93LlNldFpvb21MZXZlbCh6b29tTGV2ZWwpO1xufVxuIiwgIlxuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcbmltcG9ydCB7R2V0fSBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCB7T3BlblVSTH0gZnJvbSBcIi4vYnJvd3NlclwiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBzZW5kLlxuICogQHBhcmFtIHthbnl9IFtkYXRhPW51bGxdIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxuICpcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xuICAgIGxldCBldmVudCA9IG5ldyBXYWlsc0V2ZW50KGV2ZW50TmFtZSwgZGF0YSk7XG4gICAgRW1pdChldmVudCk7XG59XG5cbi8qKlxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gZWxlbWVudHMgd2l0aCBgd21sLWV2ZW50YCBhdHRyaWJ1dGUuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XScpO1xuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIERldGFjaGVkOiBmYWxzZSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZCAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqL1xuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lLCBtZXRob2QpIHtcbiAgICBsZXQgdGFyZ2V0V2luZG93ID0gR2V0KHdpbmRvd05hbWUpO1xuICAgIGxldCBtZXRob2RNYXAgPSBXaW5kb3dNZXRob2RzKHRhcmdldFdpbmRvdyk7XG4gICAgaWYgKCFtZXRob2RNYXAuaGFzKG1ldGhvZCkpIHtcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xuICAgIH1cbiAgICB0cnkge1xuICAgICAgICBtZXRob2RNYXAuZ2V0KG1ldGhvZCkoKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoXCJFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJ1wiICsgbWV0aG9kICsgXCInOiBcIiArIGUpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBBZGRzIHdpbmRvdyBsaXN0ZW5lcnMgZm9yIGVsZW1lbnRzIHdpdGggdGhlICd3bWwtd2luZG93JyBhdHRyaWJ1dGUuXG4gKiBSZW1vdmVzIGFueSBleGlzdGluZyBsaXN0ZW5lcnMgYmVmb3JlIGFkZGluZyBuZXcgb25lcy5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLXdpbmRvd10nKTtcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8ICdjbGljayc7XG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdGFyZ2V0LXdpbmRvdycpIHx8ICcnO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XG4gICAgICAgIH07XG5cbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBBZGRzIGEgbGlzdGVuZXIgdG8gZWxlbWVudHMgd2l0aCB0aGUgJ3dtbC1vcGVudXJsJyBhdHRyaWJ1dGUuXG4gKiBXaGVuIHRoZSBzcGVjaWZpZWQgdHJpZ2dlciBldmVudCBpcyBmaXJlZCBvbiBhbnkgb2YgdGhlc2UgZWxlbWVudHMsXG4gKiB0aGUgbGlzdGVuZXIgd2lsbCBvcGVuIHRoZSBVUkwgc3BlY2lmaWVkIGJ5IHRoZSAnd21sLW9wZW51cmwnIGF0dHJpYnV0ZS5cbiAqIElmIGEgJ3dtbC1jb25maXJtJyBhdHRyaWJ1dGUgaXMgcHJvdmlkZWQsIGEgY29uZmlybWF0aW9uIGRpYWxvZyB3aWxsIGJlIGRpc3BsYXllZCxcbiAqIGFuZCB0aGUgVVJMIHdpbGwgb25seSBiZSBvcGVuZWQgaWYgdGhlIHVzZXIgY29uZmlybXMuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtb3BlbnVybF0nKTtcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtb3BlbnVybCcpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICB2b2lkIE9wZW5VUkwodXJsKTtcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG4vKipcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCkge1xuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpO1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBtYXAgb2YgYWxsIG1ldGhvZHMgaW4gdGhlIGN1cnJlbnQgd2luZG93LlxuICogQHJldHVybnMge01hcH0gLSBBIG1hcCBvZiB3aW5kb3cgbWV0aG9kcy5cbiAqL1xuZnVuY3Rpb24gV2luZG93TWV0aG9kcyh0YXJnZXRXaW5kb3cpIHtcbiAgICAvLyBDcmVhdGUgYSBuZXcgbWFwIHRvIHN0b3JlIG1ldGhvZHNcbiAgICBsZXQgcmVzdWx0ID0gbmV3IE1hcCgpO1xuXG4gICAgLy8gSXRlcmF0ZSBvdmVyIGFsbCBwcm9wZXJ0aWVzIG9mIHRoZSB3aW5kb3cgb2JqZWN0XG4gICAgZm9yIChsZXQgbWV0aG9kIGluIHRhcmdldFdpbmRvdykge1xuICAgICAgICAvLyBDaGVjayBpZiB0aGUgcHJvcGVydHkgaXMgaW5kZWVkIGEgbWV0aG9kIChmdW5jdGlvbilcbiAgICAgICAgaWYodHlwZW9mIHRhcmdldFdpbmRvd1ttZXRob2RdID09PSAnZnVuY3Rpb24nKSB7XG4gICAgICAgICAgICAvLyBBZGQgdGhlIG1ldGhvZCB0byB0aGUgbWFwXG4gICAgICAgICAgICByZXN1bHQuc2V0KG1ldGhvZCwgdGFyZ2V0V2luZG93W21ldGhvZF0pO1xuICAgICAgICB9XG5cbiAgICB9XG4gICAgLy8gUmV0dXJuIHRoZSBtYXAgb2Ygd2luZG93IG1ldGhvZHNcbiAgICByZXR1cm4gcmVzdWx0O1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcbiAqL1xuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQge0V2ZW50VHlwZXN9IGZyb20gXCIuL2V2ZW50X3R5cGVzXCI7XG5leHBvcnQgY29uc3QgVHlwZXMgPSBFdmVudFR5cGVzO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuRXZlbnRzLCAnJyk7XG5jb25zdCBFbWl0TWV0aG9kID0gMDtcbmNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgY29uc3RydWN0b3IobmFtZSwgZGF0YSA9IG51bGwpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzZXR1cCgpIHtcbn1cblxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50KSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIGxldCB0b1JlbW92ZSA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHJldHVybiB0cnVlO1xuICAgICAgICB9KTtcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xuICAgICAgICAgICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3MgLSBUaGUgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdGhlIGNhbGxiYWNrIGNhbiBiZSBjYWxsZWQgZm9yIHRoZSBldmVudC4gT25jZSB0aGUgbWF4aW11bSBudW1iZXIgaXMgcmVhY2hlZCwgdGhlIGNhbGxiYWNrIHdpbGwgbm8gbG9uZ2VyIGJlIGNhbGxlZC5cbiAqXG4gQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCB3aGVuIHRoZSBzcGVjaWZpZWQgZXZlbnQgb2NjdXJzLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZC4gSXQgdGFrZXMgbm8gcGFyYW1ldGVycy5cbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuICovXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lLCBjYWxsYmFjaykgeyByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7IH1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIGV2ZW50IG9jY3Vycy5cbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHsgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7IH1cblxuLyoqXG4gKiBSZW1vdmVzIHRoZSBzcGVjaWZpZWQgbGlzdGVuZXIgZnJvbSB0aGUgZXZlbnQgbGlzdGVuZXJzIGNvbGxlY3Rpb24uXG4gKiBJZiBhbGwgbGlzdGVuZXJzIGZvciB0aGUgZXZlbnQgYXJlIHJlbW92ZWQsIHRoZSBldmVudCBrZXkgaXMgZGVsZXRlZCBmcm9tIHRoZSBjb2xsZWN0aW9uLlxuICpcbiAqIEBwYXJhbSB7T2JqZWN0fSBsaXN0ZW5lciAtIFRoZSBsaXN0ZW5lciB0byBiZSByZW1vdmVkLlxuICovXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG59XG5cblxuLyoqXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxuICogQHBhcmFtIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzIC0gQWRkaXRpb25hbCBldmVudCBuYW1lcyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcbn1cbi8qKlxuICogUmVtb3ZlcyBhbGwgZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBmdW5jdGlvbiBPZmZBbGxcbiAqIEByZXR1cm5zIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkgeyBldmVudExpc3RlbmVycy5jbGVhcigpOyB9XG5cbi8qKlxuICogRW1pdHMgYW4gZXZlbnQgdXNpbmcgdGhlIGdpdmVuIGV2ZW50IG5hbWUuXG4gKlxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxuICogQHJldHVybnMge2FueX0gLSBUaGUgcmVzdWx0IG9mIHRoZSBlbWl0dGVkIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkgeyByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCBldmVudCk7IH1cbiIsICJcbmV4cG9ydCBjb25zdCBFdmVudFR5cGVzID0ge1xuXHRXaW5kb3dzOiB7XG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0QVBNUG93ZXJTdGF0dXNDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclN0YXR1c0NoYW5nZVwiLFxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwid2luZG93czpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZDogXCJ3aW5kb3dzOldlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkXCIsXG5cdFx0V2luZG93SW5hY3RpdmU6IFwid2luZG93czpXaW5kb3dJbmFjdGl2ZVwiLFxuXHRcdFdpbmRvd0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Q2xvc2U6IFwid2luZG93czpXaW5kb3dDbG9zZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHR9LFxuXHRNYWM6IHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZU5vdGlmaWNhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVOb3RpZmljYXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZSFcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW4hXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd0RpZERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRFbmRTaGVldDogXCJtYWM6V2luZG93RGlkRW5kU2hlZXRcIixcblx0XHRXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4cG9zZTogXCJtYWM6V2luZG93RGlkRXhwb3NlXCIsXG5cdFx0V2luZG93RGlkRm9jdXM6IFwibWFjOldpbmRvd0RpZEZvY3VzXCIsXG5cdFx0V2luZG93RGlkTWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZE1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJtYWM6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT25TY3JlZW5cIixcblx0XHRXaW5kb3dEaWRSZXNpZ25LZXk6IFwibWFjOldpbmRvd0RpZFJlc2lnbktleVwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbk1haW46IFwibWFjOldpbmRvd0RpZFJlc2lnbk1haW5cIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibWFjOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZTogXCJtYWM6V2luZG93RGlkVXBkYXRlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUFscGhhXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVNoYWRvd1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5XCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlIVwiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdE1lbnVXaWxsT3BlbjogXCJtYWM6TWVudVdpbGxPcGVuXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWREaXNwbGF5SXRlbTogXCJtYWM6TWVudURpZERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVEaWRSZW1vdmVJdGVtOiBcIm1hYzpNZW51RGlkUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsQmVnaW5UcmFja2luZzogXCJtYWM6TWVudVdpbGxCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51RGlkUG9wVXA6IFwibWFjOk1lbnVEaWRQb3BVcFwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0fSxcblx0TGludXg6IHtcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93TG9hZENoYW5nZWQ6IFwibGludXg6V2luZG93TG9hZENoYW5nZWRcIixcblx0XHRXaW5kb3dEZWxldGVFdmVudDogXCJsaW51eDpXaW5kb3dEZWxldGVFdmVudFwiLFxuXHRcdFdpbmRvd0ZvY3VzSW46IFwibGludXg6V2luZG93Rm9jdXNJblwiLFxuXHRcdFdpbmRvd0ZvY3VzT3V0OiBcImxpbnV4OldpbmRvd0ZvY3VzT3V0XCIsXG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxuXHR9LFxuXHRDb21tb246IHtcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwiY29tbW9uOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcImNvbW1vbjpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJjb21tb246V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwiY29tbW9uOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcImNvbW1vbjpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcblx0XHRXaW5kb3dGb2N1czogXCJjb21tb246V2luZG93Rm9jdXNcIixcblx0XHRXaW5kb3dMb3N0Rm9jdXM6IFwiY29tbW9uOldpbmRvd0xvc3RGb2N1c1wiLFxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dIaWRlOiBcImNvbW1vbjpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RmlsZXNEcm9wcGVkXCIsXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiY29tbW9uOlRoZW1lQ2hhbmdlZFwiLFxuXHR9LFxufTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IE9wZW5GaWxlRGlhbG9nT3B0aW9uc1xuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VGaWxlc10gLSBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNyZWF0ZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Nob3dIaWRkZW5GaWxlc10gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Jlc29sdmVzQWxpYXNlc10gLSBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbl0gLSBJbmRpY2F0ZXMgaWYgbXVsdGlwbGUgc2VsZWN0aW9uIGlzIGFsbG93ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtIaWRlRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1RyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtBbGxvd3NPdGhlckZpbGV0eXBlc10gLSBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC5cbiAqIEBwcm9wZXJ0eSB7RmlsZUZpbHRlcltdfSBbRmlsdGVyc10gLSBBcnJheSBvZiBmaWxlIGZpbHRlcnMuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1RpdGxlXSAtIFRpdGxlIG9mIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0J1dHRvblRleHRdIC0gVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0RpcmVjdG9yeV0gLSBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy5cbiAqL1xuXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gU2F2ZUZpbGVEaWFsb2dPcHRpb25zXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0ZpbGVuYW1lXSAtIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ3JlYXRlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbU2hvd0hpZGRlbkZpbGVzXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c011bHRpcGxlU2VsZWN0aW9uXSAtIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0hpZGVFeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c090aGVyRmlsZXR5cGVzXSAtIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGl0bGUgb2YgdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTWVzc2FnZV0gLSBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlyZWN0b3J5XSAtIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LlxuICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gTWVzc2FnZURpYWxvZ09wdGlvbnNcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGhlIHRpdGxlIG9mIHRoZSBkaWFsb2cgd2luZG93LlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtNZXNzYWdlXSAtIFRoZSBtYWluIG1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtCdXR0b25bXX0gW0J1dHRvbnNdIC0gQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLlxuICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gQnV0dG9uXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0xhYmVsXSAtIFRleHQgdGhhdCBhcHBlYXJzIHdpdGhpbiB0aGUgYnV0dG9uLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNDYW5jZWxdIC0gVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0lzRGVmYXVsdF0gLSBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGJlIHRoZSBkZWZhdWx0IGFjdGlvbiB3aGVuIHRoZSB1c2VyIHByZXNzZXMgZW50ZXIuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBGaWxlRmlsdGVyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0Rpc3BsYXlOYW1lXSAtIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtQYXR0ZXJuXSAtIFBhdHRlcm4gdG8gbWF0Y2ggZm9yIHRoZSBmaWx0ZXIsIGUuZy4gXCIqLnR4dDsqLm1kXCIgZm9yIHRleHQgbWFya2Rvd24gZmlsZXMuXG4gKi9cblxuLy8gc2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5kaWFsb2dFcnJvckNhbGxiYWNrID0gZGlhbG9nRXJyb3JDYWxsYmFjaztcbndpbmRvdy5fd2FpbHMuZGlhbG9nUmVzdWx0Q2FsbGJhY2sgPSBkaWFsb2dSZXN1bHRDYWxsYmFjaztcblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbi8vIERlZmluZSBjb25zdGFudHMgZnJvbSB0aGUgYG1ldGhvZHNgIG9iamVjdCBpbiBUaXRsZSBDYXNlXG5jb25zdCBEaWFsb2dJbmZvID0gMDtcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xuY29uc3QgRGlhbG9nRXJyb3IgPSAyO1xuY29uc3QgRGlhbG9nUXVlc3Rpb24gPSAzO1xuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xuY29uc3QgRGlhbG9nU2F2ZUZpbGUgPSA1O1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5EaWFsb2csICcnKTtcbmNvbnN0IGRpYWxvZ1Jlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgaWQgdGhhdCBpcyBub3QgcHJlc2VudCBpbiBkaWFsb2dSZXNwb25zZXMuXG4gKiBAcmV0dXJucyB7c3RyaW5nfSB1bmlxdWUgaWRcbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIFNob3dzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge251bWJlcn0gdHlwZSAtIHR5cGUgb2YgZGlhbG9nXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfE9wZW5GaWxlRGlhbG9nT3B0aW9uc3xTYXZlRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBvcHRpb25zIGZvciB0aGUgZGlhbG9nXG4gKiBAcmV0dXJucyB7UHJvbWlzZX0gcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggcmVzdWx0IG9mIGRpYWxvZ1xuICovXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucyA9IHt9KSB7XG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5jYXRjaCgoZXJyb3IpID0+IHtcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgY2FsbGJhY2sgZnJvbSBhIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgSUQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIGRhdGEgcmVjZWl2ZWQgZnJvbSB0aGUgZGlhbG9nLlxuICogQHBhcmFtIHtib29sZWFufSBpc0pTT04gLSBGbGFnIGluZGljYXRpbmcgd2hldGhlciB0aGUgZGF0YSBpcyBpbiBKU09OIGZvcm1hdC5cbiAqXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9XG4gKi9cbmZ1bmN0aW9uIGRpYWxvZ1Jlc3VsdENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcbiAgICAgICAgfVxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbGJhY2sgZnVuY3Rpb24gZm9yIGhhbmRsaW5nIGVycm9ycyBpbiBkaWFsb2cuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSBkaWFsb2cgcmVzcG9uc2UuXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICpcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cblxuLy8gUmVwbGFjZSBgbWV0aG9kc2Agd2l0aCBjb25zdGFudHMgaW4gVGl0bGUgQ2FzZVxuXG4vKipcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBjb25zdCBJbmZvID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgY29uc3QgV2FybmluZyA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7XG5cbi8qKlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGNvbnN0IEVycm9yID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7XG5cbi8qKlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGNvbnN0IFF1ZXN0aW9uID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7XG5cbi8qKlxuICogQHBhcmFtIHtPcGVuRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nW118c3RyaW5nPn0gUmV0dXJucyBzZWxlY3RlZCBmaWxlIG9yIGxpc3Qgb2YgZmlsZXMuIFJldHVybnMgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaXMgc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBjb25zdCBPcGVuRmlsZSA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpO1xuXG4vKipcbiAqIEBwYXJhbSB7U2F2ZUZpbGVEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUuIFJldHVybnMgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaXMgc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBjb25zdCBTYXZlRmlsZSA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nU2F2ZUZpbGUsIG9wdGlvbnMpO1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWVcIjtcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxuLy8gU2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5jYWxsUmVzdWx0SGFuZGxlciA9IHJlc3VsdEhhbmRsZXI7XG53aW5kb3cuX3dhaWxzLmNhbGxFcnJvckhhbmRsZXIgPSBlcnJvckhhbmRsZXI7XG5cblxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2FsbCwgJycpO1xubGV0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwKCk7XG5cbi8qKlxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cbiAqXG4gKiBAcmV0dXJuIHtzdHJpbmd9IC0gQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgY2FsbFJlc3BvbnNlcyBzZXQuXG4gKi9cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgcmVzdWx0IG9mIGEgY2FsbCByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBpZCBvZiB0aGUgcmVxdWVzdCB0byBoYW5kbGUgdGhlIHJlc3VsdCBmb3IuXG4gKiBAcGFyYW0ge3N0cmluZ30gZGF0YSAtIFRoZSByZXN1bHQgZGF0YSBvZiB0aGUgcmVxdWVzdC5cbiAqIEBwYXJhbSB7Ym9vbGVhbn0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXG4gKlxuICogQHJldHVybiB7dW5kZWZpbmVkfSAtIFRoaXMgbWV0aG9kIGRvZXMgbm90IHJldHVybiBhbnkgdmFsdWUuXG4gKi9cbmZ1bmN0aW9uIHJlc3VsdEhhbmRsZXIoaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGNvbnN0IHByb21pc2VIYW5kbGVyID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xuICAgIGlmIChwcm9taXNlSGFuZGxlcikge1xuICAgICAgICBwcm9taXNlSGFuZGxlci5yZXNvbHZlKGlzSlNPTiA/IEpTT04ucGFyc2UoZGF0YSkgOiBkYXRhKTtcbiAgICB9XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgZXJyb3IgZnJvbSBhIGNhbGwgcmVxdWVzdC5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIHByb21pc2UgaGFuZGxlci5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UgdG8gcmVqZWN0IHRoZSBwcm9taXNlIGhhbmRsZXIgd2l0aC5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5mdW5jdGlvbiBlcnJvckhhbmRsZXIoaWQsIG1lc3NhZ2UpIHtcbiAgICBjb25zdCBwcm9taXNlSGFuZGxlciA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcbiAgICAgICAgcHJvbWlzZUhhbmRsZXIucmVqZWN0KG1lc3NhZ2UpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgY2FsbFJlc3BvbnNlcyBtYXAuXG4gKlxuICogQHBhcmFtIHthbnl9IGlkIC0gVGhlIElEIG9mIHRoZSByZXNwb25zZSB0byBiZSByZXRyaWV2ZWQgYW5kIHJlbW92ZWQuXG4gKlxuICogQHJldHVybnMge2FueX0gVGhlIHJlc3BvbnNlIG9iamVjdCBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElELlxuICovXG5mdW5jdGlvbiBnZXRBbmREZWxldGVSZXNwb25zZShpZCkge1xuICAgIGNvbnN0IHJlc3BvbnNlID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICByZXR1cm4gcmVzcG9uc2U7XG59XG5cbi8qKlxuICogRXhlY3V0ZXMgYSBjYWxsIHVzaW5nIHRoZSBwcm92aWRlZCB0eXBlIGFuZCBvcHRpb25zLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfG51bWJlcn0gdHlwZSAtIFRoZSB0eXBlIG9mIGNhbGwgdG8gZXhlY3V0ZS5cbiAqIEBwYXJhbSB7T2JqZWN0fSBbb3B0aW9ucz17fV0gLSBBZGRpdGlvbmFsIG9wdGlvbnMgZm9yIHRoZSBjYWxsLlxuICogQHJldHVybiB7UHJvbWlzZX0gLSBBIHByb21pc2UgdGhhdCB3aWxsIGJlIHJlc29sdmVkIG9yIHJlamVjdGVkIGJhc2VkIG9uIHRoZSByZXN1bHQgb2YgdGhlIGNhbGwuXG4gKi9cbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMgPSB7fSkge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlLCByZWplY3QgfSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBDYWxsIG1ldGhvZC5cbiAqXG4gKiBAcGFyYW0ge09iamVjdH0gb3B0aW9ucyAtIFRoZSBvcHRpb25zIGZvciB0aGUgbWV0aG9kLlxuICogQHJldHVybnMge09iamVjdH0gLSBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBFeGVjdXRlcyBhIG1ldGhvZCBieSBuYW1lLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAdGhyb3dzIHtFcnJvcn0gSWYgdGhlIG5hbWUgaXMgbm90IGEgc3RyaW5nIG9yIGlzIG5vdCBpbiB0aGUgY29ycmVjdCBmb3JtYXQuXG4gKiBAcmV0dXJucyB7Kn0gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGV4ZWN1dGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShuYW1lLCAuLi5hcmdzKSB7XG4gICAgaWYgKHR5cGVvZiBuYW1lICE9PSBcInN0cmluZ1wiIHx8IG5hbWUuc3BsaXQoXCIuXCIpLmxlbmd0aCAhPT0gMykge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJDYWxsQnlOYW1lIHJlcXVpcmVzIGEgc3RyaW5nIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCdcIik7XG4gICAgfVxuICAgIGxldCBbcGFja2FnZU5hbWUsIHN0cnVjdE5hbWUsIG1ldGhvZE5hbWVdID0gbmFtZS5zcGxpdChcIi5cIik7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XG4gICAgICAgIHBhY2thZ2VOYW1lLFxuICAgICAgICBzdHJ1Y3ROYW1lLFxuICAgICAgICBtZXRob2ROYW1lLFxuICAgICAgICBhcmdzXG4gICAgfSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKlxuICogQHBhcmFtIHtudW1iZXJ9IG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybiB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQsIC4uLmFyZ3MpIHtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIHtcbiAgICAgICAgbWV0aG9kSUQsXG4gICAgICAgIGFyZ3NcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHBsdWdpbi5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gcGx1Z2luTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBwbHVnaW4uXG4gKiBAcGFyYW0ge3N0cmluZ30gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybnMgeyp9IC0gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xuICAgICAgICBwYWNrYWdlTmFtZTogXCJ3YWlscy1wbHVnaW5zXCIsXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXG4gICAgICAgIG1ldGhvZE5hbWUsXG4gICAgICAgIGFyZ3NcbiAgICB9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtkZWJ1Z0xvZ30gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZ1wiO1xuXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uXCI7XG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9icm93c2VyXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZFwiO1xuaW1wb3J0ICogYXMgQ29udGV4dE1lbnUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NvbnRleHRtZW51XCI7XG5pbXBvcnQgKiBhcyBEcmFnIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnXCI7XG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3NcIjtcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnNcIjtcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtXCI7XG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvd1wiO1xuaW1wb3J0ICogYXMgV01MIGZyb20gJy4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dtbCc7XG5pbXBvcnQgKiBhcyBFdmVudHMgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50c1wiO1xuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9nc1wiO1xuaW1wb3J0ICogYXMgQ2FsbCBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHNcIjtcbmltcG9ydCB7aW52b2tlfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtXCI7XG5cbi8qKipcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcblxuIEJTRCAyLUNsYXVzZSBMaWNlbnNlXG5cbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcbiBBbGwgcmlnaHRzIHJlc2VydmVkLlxuXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XG4gbW9kaWZpY2F0aW9uLCBhcmUgcGVybWl0dGVkIHByb3ZpZGVkIHRoYXQgdGhlIGZvbGxvd2luZyBjb25kaXRpb25zIGFyZSBtZXQ6XG5cbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcbiBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lci5cblxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcbiB0aGlzIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyIGluIHRoZSBkb2N1bWVudGF0aW9uXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXG5cbiBUSElTIFNPRlRXQVJFIElTIFBST1ZJREVEIEJZIFRIRSBDT1BZUklHSFQgSE9MREVSUyBBTkQgQ09OVFJJQlVUT1JTIFwiQVMgSVNcIlxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcbiBESVNDTEFJTUVELiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQ09QWVJJR0hUIEhPTERFUiBPUiBDT05UUklCVVRPUlMgQkUgTElBQkxFXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxuIFNFUlZJQ0VTOyBMT1NTIE9GIFVTRSwgREFUQSwgT1IgUFJPRklUUzsgT1IgQlVTSU5FU1MgSU5URVJSVVBUSU9OKSBIT1dFVkVSXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxuIE9GIFRISVMgU09GVFdBUkUsIEVWRU4gSUYgQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YgU1VDSCBEQU1BR0UuXG5cbiAqKiovXG5cbndpbmRvdy5fd2FpbHMuaW52b2tlPWludm9rZTtcblxud2luZG93LndhaWxzID0gd2luZG93LndhaWxzIHx8IHt9O1xud2luZG93LndhaWxzLkFwcGxpY2F0aW9uID0gQXBwbGljYXRpb247XG53aW5kb3cud2FpbHMuQnJvd3NlciA9IEJyb3dzZXI7XG53aW5kb3cud2FpbHMuQ2FsbCA9IENhbGw7XG53aW5kb3cud2FpbHMuQ2xpcGJvYXJkID0gQ2xpcGJvYXJkO1xud2luZG93LndhaWxzLkRpYWxvZ3MgPSBEaWFsb2dzO1xud2luZG93LndhaWxzLkV2ZW50cyA9IEV2ZW50cztcbndpbmRvdy53YWlscy5GbGFncyA9IEZsYWdzO1xud2luZG93LndhaWxzLlNjcmVlbnMgPSBTY3JlZW5zO1xud2luZG93LndhaWxzLlN5c3RlbSA9IFN5c3RlbTtcbndpbmRvdy53YWlscy5XaW5kb3cgPSBXaW5kb3c7XG53aW5kb3cud2FpbHMuV01MID0gV01MO1xuXG5cbmxldCBpc1JlYWR5ID0gZmFsc2VcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCBmdW5jdGlvbigpIHtcbiAgICBpc1JlYWR5ID0gdHJ1ZVxuICAgIHdpbmRvdy5fd2FpbHMuaW52b2tlKCd3YWlsczpydW50aW1lOnJlYWR5Jyk7XG4gICAgaWYoREVCVUcpIHtcbiAgICAgICAgZGVidWdMb2coXCJXYWlscyBSdW50aW1lIExvYWRlZFwiKTtcbiAgICB9XG59KVxuXG5mdW5jdGlvbiB3aGVuUmVhZHkoZm4pIHtcbiAgICBpZiAoaXNSZWFkeSB8fCBkb2N1bWVudC5yZWFkeVN0YXRlID09PSAnY29tcGxldGUnKSB7XG4gICAgICAgIGZuKCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGZuKTtcbiAgICB9XG59XG5cbndoZW5SZWFkeSgoKSA9PiB7XG4gICAgV01MLlJlbG9hZCgpO1xufSk7XG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUtPLFdBQVMsU0FBUyxTQUFTO0FBRTlCLFlBQVE7QUFBQSxNQUNKLGtCQUFrQixVQUFVO0FBQUEsTUFDNUI7QUFBQSxNQUNBO0FBQUEsSUFDSjtBQUFBLEVBQ0o7OztBQ1pBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUEsTUFBSSxjQUNGO0FBV0ssTUFBSSxTQUFTLENBQUNBLFFBQU8sT0FBTztBQUNqQyxRQUFJLEtBQUs7QUFDVCxRQUFJLElBQUlBO0FBQ1IsV0FBTyxLQUFLO0FBQ1YsWUFBTSxZQUFhLEtBQUssT0FBTyxJQUFJLEtBQU0sQ0FBQztBQUFBLElBQzVDO0FBQ0EsV0FBTztBQUFBLEVBQ1Q7OztBQ05BLE1BQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUdyQyxNQUFNLGNBQWM7QUFBQSxJQUN2QixNQUFNO0FBQUEsSUFDTixXQUFXO0FBQUEsSUFDWCxhQUFhO0FBQUEsSUFDYixRQUFRO0FBQUEsSUFDUixhQUFhO0FBQUEsSUFDYixRQUFRO0FBQUEsSUFDUixRQUFRO0FBQUEsSUFDUixTQUFT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QUFvQztBQUFBLE1BQ3BDLCtCQUErQjtBQUFBLE1BQy9CLDZCQUE2QjtBQUFBLE1BQzdCLHVCQUF1QjtBQUFBLE1BQ3ZCLGlDQUFpQztBQUFBLE1BQ2pDLDhCQUE4QjtBQUFBLE1BQzlCLDRCQUE0QjtBQUFBLE1BQzVCLHNDQUFzQztBQUFBLE1BQ3RDLDRCQUE0QjtBQUFBLE1BQzVCLHNCQUFzQjtBQUFBLE1BQ3RCLGtDQUFrQztBQUFBLE1BQ2xDLHNCQUFzQjtBQUFBLE1BQ3RCLHdCQUF3QjtBQUFBLE1BQ3hCLDJCQUEyQjtBQUFBLE1BQzNCLHdCQUF3QjtBQUFBLE1BQ3hCLG1CQUFtQjtBQUFBLE1BQ25CLDBCQUEwQjtBQUFBLE1BQzFCLDhCQUE4QjtBQUFBLE1BQzlCLHlCQUF5QjtBQUFBLE1BQ3pCLDZCQUE2QjtBQUFBLE1BQzdCLGlCQUFpQjtBQUFBLE1BQ2pCLGdCQUFnQjtBQUFBLE1BQ2hCLHNCQUFzQjtBQUFBLE1BQ3RCLGVBQWU7QUFBQSxNQUNmLHlCQUF5QjtBQUFBLE1BQ3pCLHdCQUF3QjtBQUFBLE1BQ3hCLG9CQUFvQjtBQUFBLE1BQ3BCLHFCQUFxQjtBQUFBLE1BQ3JCLGlCQUFpQjtBQUFBLE1BQ2pCLGlCQUFpQjtBQUFBLE1BQ2pCLHNCQUFzQjtBQUFBLE1BQ3RCLG1DQUFtQztBQUFBLE1BQ25DLHFDQUFxQztBQUFBLE1BQ3JDLHVCQUF1QjtBQUFBLE1BQ3ZCLHNCQUFzQjtBQUFBLE1BQ3RCLHdCQUF3QjtBQUFBLE1BQ3hCLDJCQUEyQjtBQUFBLE1BQzNCLG1CQUFtQjtBQUFBLE1BQ25CLHFCQUFxQjtBQUFBLE1BQ3JCLHNCQUFzQjtBQUFBLE1BQ3RCLHNCQUFzQjtBQUFBLE1BQ3RCLDhCQUE4QjtBQUFBLE1BQzlCLGlCQUFpQjtBQUFBLE1BQ2pCLHlCQUF5QjtBQUFBLE1BQ3pCLDJCQUEyQjtBQUFBLE1BQzNCLCtCQUErQjtBQUFBLE1BQy9CLDBCQUEwQjtBQUFBLE1BQzFCLDhCQUE4QjtBQUFBLE1BQzlCLGlCQUFpQjtBQUFBLE1BQ2pCLHVCQUF1QjtBQUFBLE1BQ3ZCLGdCQUFn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VTtBQUN2QixTQUFPLE1BQU0sT0FBTztBQUNwQixTQUFPLE1BQU0sWUFBWTtBQUN6QixTQUFPLE1BQU0sVUFBVTtBQUN2QixTQUFPLE1BQU0sU0FBUztBQUN0QixTQUFPLE1BQU0sUUFBUTtBQUNyQixTQUFPLE1BQU0sVUFBVTtBQUN2QixTQUFPLE1BQU0sU0FBUztBQUN0QixTQUFPLE1BQU0sU0FBUztBQUN0QixTQUFPLE1BQU0sTUFBTTtBQUduQixNQUFJLFVBQVU7QUFDZCxXQUFTLGlCQUFpQixvQkFBb0IsV0FBVztBQUNyRCxjQUFVO0FBQ1YsV0FBTyxPQUFPLE9BQU8scUJBQXFCO0FBQzFDLFFBQUcsTUFBTztBQUNOLGVBQVMsc0JBQXNCO0FBQUEsSUFDbkM7QUFBQSxFQUNKLENBQUM7QUFFRCxXQUFTLFVBQVUsSUFBSTtBQUNuQixRQUFJLFdBQVcsU0FBUyxlQUFlLFlBQVk7QUFDL0MsU0FBRztBQUFBLElBQ1AsT0FBTztBQUNILGVBQVMsaUJBQWlCLG9CQUFvQixFQUFFO0FBQUEsSUFDcEQ7QUFBQSxFQUNKO0FBRUEsWUFBVSxNQUFNO0FBQ1osSUFBSSxPQUFPO0FBQUEsRUFDZixDQUFDOyIsCiAgIm5hbWVzIjogWyJzaXplIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIaWRlIiwgIlNob3ciLCAic2V0UmVzaXphYmxlIiwgImNhbGwiLCAid2lkdGgiLCAiaGVpZ2h0IiwgInJlc2l6YWJsZSIsICJIaWRlIiwgIlNob3ciLCAiY2FsbCIsICJldmVudE5hbWUiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJjYWxsIiwgImdlbmVyYXRlSUQiXQp9Cg==
