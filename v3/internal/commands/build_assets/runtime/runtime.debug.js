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
    Browser: 9
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
      SystemThemeChanged: "linux:SystemThemeChanged"
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
    return new Promise((resolve, reject) => {
      const id = generateID2();
      options["call-id"] = id;
      callResponses.set(id, { resolve, reject });
      call9(type, options).catch((error) => {
        reject(error);
        callResponses.delete(id);
      });
    });
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb24uanMiLCAiLi4vLi4vLi4vcnVudGltZS9ub2RlX21vZHVsZXMvbmFub2lkL25vbi1zZWN1cmUvaW5kZXguanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZC5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvd21sLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudHMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50X3R5cGVzLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kaWFsb2dzLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLyoqXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZSkge1xuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgIGNvbnNvbGUubG9nKFxuICAgICAgICAnJWMgd2FpbHMzICVjICcgKyBtZXNzYWdlICsgJyAnLFxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjMDA5OTAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogMHB4IDNweCAzcHggMHB4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJ1xuICAgICk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24sICcnKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcbn1cblxuLyoqXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXG4gKlxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxufVxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG4vKipcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cbiAqXG4gKiBAcGFyYW0ge09iamVjdH0gb2JqZWN0IC0gVGhlIG9iamVjdCBvbiB3aGljaCB0aGUgbWV0aG9kIGlzIHRvIGJlIGludm9rZWQuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyB7RnVuY3Rpb259IEEgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24gdGhhdCB0YWtlcyB0aGUgbWV0aG9kIG5hbWUgYW5kIG9wdGlvbmFsbHkgYXJndW1lbnRzIGFuZCBpbnZva2VzIHRoZSBtZXRob2Qgd2l0aGluIHRoZSBzcGVjaWZpZWQgd2luZG93IGNvbnRleHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG4vKipcbiAqIENyZWF0ZXMgYSBuZXcgcnVudGltZSBjYWxsZXIgd2l0aCBzcGVjaWZpZWQgSUQuXG4gKlxuICogQHBhcmFtIHtvYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gKiBAcmV0dXJuIHtGdW5jdGlvbn0gLSBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3QsIHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICBpZiggbWV0aG9kICkge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIH1cbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcblxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm9iamVjdFwiLCBvYmplY3RJRCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkJyb3dzZXIsICcnKTtcbmNvbnN0IEJyb3dzZXJPcGVuVVJMID0gMDtcblxuLyoqXG4gKiBPcGVuIGEgYnJvd3NlciB3aW5kb3cgdG8gdGhlIGdpdmVuIFVSTFxuICogQHBhcmFtIHtzdHJpbmd9IHVybCAtIFRoZSBVUkwgdG8gb3BlblxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5VUkwodXJsKSB7XG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmx9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCwgJycpO1xuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XG5jb25zdCBDbGlwYm9hcmRUZXh0ID0gMTtcblxuLyoqXG4gKiBTZXRzIHRoZSB0ZXh0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cbiAqIEByZXR1cm4ge1Byb21pc2V9IC0gQSBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2hlbiB0aGUgb3BlcmF0aW9uIGlzIHN1Y2Nlc3NmdWwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRTZXRUZXh0LCB7dGV4dH0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU3lzdGVtLCAnJyk7XG5jb25zdCBzeXN0ZW1Jc0RhcmtNb2RlID0gMDtcbmNvbnN0IGVudmlyb25tZW50ID0gMTtcblxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2cpIHtcbiAgICBpZih3aW5kb3cuY2hyb21lKSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UobXNnKTtcbiAgICB9XG4gICAgcmV0dXJuIHdpbmRvdy53ZWJraXQubWVzc2FnZUhhbmRsZXJzLmV4dGVybmFsLnBvc3RNZXNzYWdlKG1zZyk7XG59XG5cbi8qKlxuICogQGZ1bmN0aW9uXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxuICogQHJldHVybnMge1Byb21pc2U8Ym9vbGVhbj59IC0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYSBib29sZWFuIHZhbHVlIGluZGljYXRpbmcgaWYgdGhlIHN5c3RlbSBpcyBpbiBkYXJrIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCkge1xuICAgIHJldHVybiBjYWxsKHN5c3RlbUlzRGFya01vZGUpO1xufVxuXG4vKipcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxuICpcbiAqIEBhc3luY1xuICogQGZ1bmN0aW9uIENhcGFiaWxpdGllc1xuICogQHJldHVybnMge1Byb21pc2U8T2JqZWN0Pn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhcGFiaWxpdGllcygpIHtcbiAgICBsZXQgcmVzcG9uc2UgPSBmZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIik7XG4gICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbn1cblxuLyoqXG4gKiBAdHlwZWRlZiB7b2JqZWN0fSBFbnZpcm9ubWVudEluZm9cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBPUyAtIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBBcmNoIC0gVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLlxuICovXG5cbi8qKlxuICogQGZ1bmN0aW9uXG4gKiBSZXRyaWV2ZXMgZW52aXJvbm1lbnQgZGV0YWlscy5cbiAqIEByZXR1cm5zIHtQcm9taXNlPEVudmlyb25tZW50SW5mbz59IC0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpIHtcbiAgICByZXR1cm4gY2FsbChlbnZpcm9ubWVudCk7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cbiAqXG4gKiBAcmV0dXJuIHtib29sZWFufSBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MsIG90aGVyd2lzZSBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzV2luZG93cygpIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJ3aW5kb3dzXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXG4gKlxuICogQHJldHVybnMge2Jvb2xlYW59IFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImxpbnV4XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImRhcndpblwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhbWQ2NFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxuICpcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gLSBSZXR1cm5zIHRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZSwgb3RoZXJ3aXNlIHJldHVybnMgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtNjRcIjtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIElzRGVidWcoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuRGVidWcgPT09IHRydWU7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5pbXBvcnQge0lzRGVidWd9IGZyb20gXCIuL3N5c3RlbVwiO1xuXG4vLyBzZXR1cFxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUsICcnKTtcbmNvbnN0IENvbnRleHRNZW51T3BlbiA9IDA7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxuICAgIGxldCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xuICAgIGxldCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudVwiKTtcbiAgICBjdXN0b21Db250ZXh0TWVudSA9IGN1c3RvbUNvbnRleHRNZW51ID8gY3VzdG9tQ29udGV4dE1lbnUudHJpbSgpIDogXCJcIjtcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgbGV0IGN1c3RvbUNvbnRleHRNZW51RGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGN1c3RvbUNvbnRleHRNZW51RGF0YSk7XG4gICAgICAgIHJldHVyblxuICAgIH1cblxuICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpO1xufVxuXG5cbi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcblxuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxuICAgIGlmIChJc0RlYnVnKCkpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICBjb25zdCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xuICAgIGNvbnN0IGNvbXB1dGVkU3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcbiAgICBjb25zdCBkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24gPSBjb21wdXRlZFN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpO1xuICAgIHN3aXRjaCAoZGVmYXVsdENvbnRleHRNZW51QWN0aW9uKSB7XG4gICAgICAgIGNhc2UgXCJzaG93XCI6XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGNhc2UgXCJoaWRlXCI6XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcbiAgICAgICAgICAgIGlmIChlbGVtZW50LmlzQ29udGVudEVkaXRhYmxlKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXG4gICAgICAgICAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XG4gICAgICAgICAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSAoc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMClcbiAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcbiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xuICAgICAgICAgICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gZWxlbWVudCkge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRhZ25hbWUgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbiAgICAgICAgICAgIGlmIChlbGVtZW50LnRhZ05hbWUgPT09IFwiSU5QVVRcIiB8fCBlbGVtZW50LnRhZ05hbWUgPT09IFwiVEVYVEFSRUFcIikge1xuICAgICAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24gfHwgKCFlbGVtZW50LnJlYWRPbmx5ICYmICFlbGVtZW50LmRpc2FibGVkKSkge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIFJldHJpZXZlcyB0aGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5IGZyb20gdGhlIGZsYWcgbWFwLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBrZXlTdHJpbmcgLSBUaGUga2V5IHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBmb3IuXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleVN0cmluZykge1xuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleVN0cmluZ107XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJVbmFibGUgdG8gcmV0cmlldmUgZmxhZyAnXCIgKyBrZXlTdHJpbmcgKyBcIic6IFwiICsgZSk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7aW52b2tlLCBJc1dpbmRvd3N9IGZyb20gXCIuL3N5c3RlbVwiO1xuaW1wb3J0IHtHZXRGbGFnfSBmcm9tIFwiLi9mbGFnc1wiO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9IHNldFJlc2l6YWJsZTtcbndpbmRvdy5fd2FpbHMuZW5kRHJhZyA9IGVuZERyYWc7XG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgb25Nb3VzZURvd24pO1xud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIG9uTW91c2VNb3ZlKTtcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgb25Nb3VzZVVwKTtcblxuXG5sZXQgc2hvdWxkRHJhZyA9IGZhbHNlO1xubGV0IHJlc2l6ZUVkZ2UgPSBudWxsO1xubGV0IHJlc2l6YWJsZSA9IGZhbHNlO1xubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcblxuZnVuY3Rpb24gZHJhZ1Rlc3QoZSkge1xuICAgIGxldCB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2Via2l0LWFwcC1yZWdpb25cIik7XG4gICAgaWYgKCF2YWwgfHwgdmFsID09PSBcIlwiIHx8IHZhbC50cmltKCkgIT09IFwiZHJhZ1wiIHx8IGUuYnV0dG9ucyAhPT0gMSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcbn1cblxuZnVuY3Rpb24gc2V0UmVzaXphYmxlKHZhbHVlKSB7XG4gICAgcmVzaXphYmxlID0gdmFsdWU7XG59XG5cbmZ1bmN0aW9uIGVuZERyYWcoKSB7XG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiB0ZXN0UmVzaXplKCkge1xuICAgIGlmKCByZXNpemVFZGdlICkge1xuICAgICAgICBpbnZva2UoYHJlc2l6ZToke3Jlc2l6ZUVkZ2V9YCk7XG4gICAgICAgIHJldHVybiB0cnVlXG4gICAgfVxuICAgIHJldHVybiBmYWxzZTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xuICAgIGlmKElzV2luZG93cygpICYmIHRlc3RSZXNpemUoKSB8fCBkcmFnVGVzdChlKSkge1xuICAgICAgICBzaG91bGREcmFnID0gISFpc1ZhbGlkRHJhZyhlKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGlzVmFsaWREcmFnKGUpIHtcbiAgICAvLyBJZ25vcmUgZHJhZyBvbiBzY3JvbGxiYXJzXG4gICAgcmV0dXJuICEoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZVVwKGUpIHtcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgIGVuZERyYWcoKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IgPSBkZWZhdWx0Q3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvcjtcbiAgICByZXNpemVFZGdlID0gY3Vyc29yO1xufVxuXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShlKSB7XG4gICAgc2hvdWxkRHJhZyA9IGNoZWNrRHJhZyhlKTtcbiAgICBpZiAoSXNXaW5kb3dzKCkgJiYgcmVzaXphYmxlKSB7XG4gICAgICAgIGhhbmRsZVJlc2l6ZShlKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNoZWNrRHJhZyhlKSB7XG4gICAgbGV0IG1vdXNlUHJlc3NlZCA9IGUuYnV0dG9ucyAhPT0gdW5kZWZpbmVkID8gZS5idXR0b25zIDogZS53aGljaDtcbiAgICBpZihzaG91bGREcmFnICYmIG1vdXNlUHJlc3NlZCA+IDApIHtcbiAgICAgICAgaW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbiAgICByZXR1cm4gc2hvdWxkRHJhZztcbn1cblxuZnVuY3Rpb24gaGFuZGxlUmVzaXplKGUpIHtcbiAgICBsZXQgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcbiAgICBsZXQgcmVzaXplSGFuZGxlV2lkdGggPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZVdpZHRoXCIpIHx8IDU7XG5cbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXNcbiAgICBsZXQgY29ybmVyRXh0cmEgPSBHZXRGbGFnKFwicmVzaXplQ29ybmVyRXh0cmFcIikgfHwgMTA7XG5cbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgbGV0IHRvcEJvcmRlciA9IGUuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXJzXG4gICAgbGV0IHJpZ2h0Q29ybmVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgbGV0IGxlZnRDb3JuZXIgPSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgbGV0IHRvcENvcm5lciA9IGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG4gICAgbGV0IGJvdHRvbUNvcm5lciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG5cbiAgICAvLyBJZiB3ZSBhcmVuJ3Qgb24gYW4gZWRnZSwgYnV0IHdlcmUsIHJlc2V0IHRoZSBjdXJzb3IgdG8gZGVmYXVsdFxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9XG4gICAgLy8gQWRqdXN0ZWQgZm9yIGNvcm5lciBhcmVhc1xuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gUG9zaXRpb25cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZS5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBZIC0gVGhlIFkgY29vcmRpbmF0ZS5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNpemVcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIHdpZHRoLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgaGVpZ2h0LlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBSZWN0XG4gKiBAcHJvcGVydHkge251bWJlcn0gWCAtIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHRvcC1sZWZ0IGNvcm5lci5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBZIC0gVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgdG9wLWxlZnQgY29ybmVyLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuXG4gKiBAcHJvcGVydHkge251bWJlcn0gSGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7KCdaZXJvJ3wnTmluZXR5J3wnT25lRWlnaHR5J3wnVHdvU2V2ZW50eScpfSBSb3RhdGlvblxuICogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uIENhbiBiZSBvbmUgb2YgJ1plcm8nLCAnTmluZXR5JywgJ09uZUVpZ2h0eScsICdUd29TZXZlbnR5Jy5cbiAqL1xuXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gU2NyZWVuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gSWQgLSBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBOYW1lIC0gSHVtYW4gcmVhZGFibGUgbmFtZSBvZiB0aGUgc2NyZWVuLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFNjYWxlIC0gVGhlIHJlc29sdXRpb24gc2NhbGUgb2YgdGhlIHNjcmVlbi4gMSA9IHN0YW5kYXJkIHJlc29sdXRpb24sIDIgPSBoaWdoIChSZXRpbmEpLCBldGMuXG4gKiBAcHJvcGVydHkge1Bvc2l0aW9ufSBQb3NpdGlvbiAtIENvbnRhaW5zIHRoZSBYIGFuZCBZIGNvb3JkaW5hdGVzIG9mIHRoZSBzY3JlZW4ncyBwb3NpdGlvbi5cbiAqIEBwcm9wZXJ0eSB7U2l6ZX0gU2l6ZSAtIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uXG4gKiBAcHJvcGVydHkge1JlY3R9IEJvdW5kcyAtIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC5cbiAqIEBwcm9wZXJ0eSB7UmVjdH0gV29ya0FyZWEgLSBDb250YWlucyB0aGUgYXJlYSBvZiB0aGUgc2NyZWVuIHRoYXQgaXMgYWN0dWFsbHkgdXNhYmxlIChleGNsdWRpbmcgdGFza2JhciBhbmQgb3RoZXIgc3lzdGVtIFVJKS5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gSXNQcmltYXJ5IC0gVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uXG4gKiBAcHJvcGVydHkge1JvdGF0aW9ufSBSb3RhdGlvbiAtIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLlxuICovXG5cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU2NyZWVucywgJycpO1xuXG5jb25zdCBnZXRBbGwgPSAwO1xuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XG5jb25zdCBnZXRDdXJyZW50ID0gMjtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChnZXRBbGwpO1xufVxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKGdldEN1cnJlbnQpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vLyBJbXBvcnQgc2NyZWVuIGpzZG9jIGRlZmluaXRpb24gZnJvbSAuL3NjcmVlbnMuanNcbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vc2NyZWVuc1wiKS5TY3JlZW59IFNjcmVlblxuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuY29uc3QgY2VudGVyID0gMDtcbmNvbnN0IHNldFRpdGxlID0gMTtcbmNvbnN0IGZ1bGxzY3JlZW4gPSAyO1xuY29uc3QgdW5GdWxsc2NyZWVuID0gMztcbmNvbnN0IHNldFNpemUgPSA0O1xuY29uc3Qgc2l6ZSA9IDU7XG5jb25zdCBzZXRNYXhTaXplID0gNjtcbmNvbnN0IHNldE1pblNpemUgPSA3O1xuY29uc3Qgc2V0QWx3YXlzT25Ub3AgPSA4O1xuY29uc3Qgc2V0UmVsYXRpdmVQb3NpdGlvbiA9IDk7XG5jb25zdCByZWxhdGl2ZVBvc2l0aW9uID0gMTA7XG5jb25zdCBzY3JlZW4gPSAxMTtcbmNvbnN0IGhpZGUgPSAxMjtcbmNvbnN0IG1heGltaXNlID0gMTM7XG5jb25zdCB1bk1heGltaXNlID0gMTQ7XG5jb25zdCB0b2dnbGVNYXhpbWlzZSA9IDE1O1xuY29uc3QgbWluaW1pc2UgPSAxNjtcbmNvbnN0IHVuTWluaW1pc2UgPSAxNztcbmNvbnN0IHJlc3RvcmUgPSAxODtcbmNvbnN0IHNob3cgPSAxOTtcbmNvbnN0IGNsb3NlID0gMjA7XG5jb25zdCBzZXRCYWNrZ3JvdW5kQ29sb3VyID0gMjE7XG5jb25zdCBzZXRSZXNpemFibGUgPSAyMjtcbmNvbnN0IHdpZHRoID0gMjM7XG5jb25zdCBoZWlnaHQgPSAyNDtcbmNvbnN0IHpvb21JbiA9IDI1O1xuY29uc3Qgem9vbU91dCA9IDI2O1xuY29uc3Qgem9vbVJlc2V0ID0gMjc7XG5jb25zdCBnZXRab29tTGV2ZWwgPSAyODtcbmNvbnN0IHNldFpvb21MZXZlbCA9IDI5O1xuXG5jb25zdCB0aGlzV2luZG93ID0gR2V0KCcnKTtcblxuZnVuY3Rpb24gY3JlYXRlV2luZG93KGNhbGwpIHtcbiAgICByZXR1cm4ge1xuICAgICAgICBHZXQ6ICh3aW5kb3dOYW1lKSA9PiBjcmVhdGVXaW5kb3cobmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIHdpbmRvd05hbWUpKSxcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKGNlbnRlciksXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IGNhbGwoc2V0VGl0bGUsIHt0aXRsZX0pLFxuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKGZ1bGxzY3JlZW4pLFxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwodW5GdWxsc2NyZWVuKSxcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcbiAgICAgICAgU2l6ZTogKCkgPT4gY2FsbChzaXplKSxcbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWF4U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcbiAgICAgICAgU2V0TWluU2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWluU2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcbiAgICAgICAgU2V0QWx3YXlzT25Ub3A6IChvblRvcCkgPT4gY2FsbChzZXRBbHdheXNPblRvcCwge2Fsd2F5c09uVG9wOiBvblRvcH0pLFxuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbChzZXRSZWxhdGl2ZVBvc2l0aW9uLCB7eCwgeX0pLFxuICAgICAgICBSZWxhdGl2ZVBvc2l0aW9uOiAoKSA9PiBjYWxsKHJlbGF0aXZlUG9zaXRpb24pLFxuICAgICAgICBTY3JlZW46ICgpID0+IGNhbGwoc2NyZWVuKSxcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbChoaWRlKSxcbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IGNhbGwobWF4aW1pc2UpLFxuICAgICAgICBVbk1heGltaXNlOiAoKSA9PiBjYWxsKHVuTWF4aW1pc2UpLFxuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gY2FsbCh0b2dnbGVNYXhpbWlzZSksXG4gICAgICAgIE1pbmltaXNlOiAoKSA9PiBjYWxsKG1pbmltaXNlKSxcbiAgICAgICAgVW5NaW5pbWlzZTogKCkgPT4gY2FsbCh1bk1pbmltaXNlKSxcbiAgICAgICAgUmVzdG9yZTogKCkgPT4gY2FsbChyZXN0b3JlKSxcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbChzaG93KSxcbiAgICAgICAgQ2xvc2U6ICgpID0+IGNhbGwoY2xvc2UpLFxuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gY2FsbChzZXRCYWNrZ3JvdW5kQ29sb3VyLCB7ciwgZywgYiwgYX0pLFxuICAgICAgICBTZXRSZXNpemFibGU6IChyZXNpemFibGUpID0+IGNhbGwoc2V0UmVzaXphYmxlLCB7cmVzaXphYmxlfSksXG4gICAgICAgIFdpZHRoOiAoKSA9PiBjYWxsKHdpZHRoKSxcbiAgICAgICAgSGVpZ2h0OiAoKSA9PiBjYWxsKGhlaWdodCksXG4gICAgICAgIFpvb21JbjogKCkgPT4gY2FsbCh6b29tSW4pLFxuICAgICAgICBab29tT3V0OiAoKSA9PiBjYWxsKHpvb21PdXQpLFxuICAgICAgICBab29tUmVzZXQ6ICgpID0+IGNhbGwoem9vbVJlc2V0KSxcbiAgICAgICAgR2V0Wm9vbUxldmVsOiAoKSA9PiBjYWxsKGdldFpvb21MZXZlbCksXG4gICAgICAgIFNldFpvb21MZXZlbDogKHpvb21MZXZlbCkgPT4gY2FsbChzZXRab29tTGV2ZWwsIHt6b29tTGV2ZWx9KSxcbiAgICB9O1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cbiAqIEByZXR1cm4ge09iamVjdH0gLSBUaGUgc3BlY2lmaWVkIHdpbmRvdyBvYmplY3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXQod2luZG93TmFtZSkge1xuICAgIHJldHVybiBjcmVhdGVXaW5kb3cobmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIHdpbmRvd05hbWUpKTtcbn1cblxuLyoqXG4gKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENlbnRlcigpIHtcbiAgICB0aGlzV2luZG93LkNlbnRlcigpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gKiBAcGFyYW0ge3N0cmluZ30gdGl0bGUgLSBUaGUgdGl0bGUgdG8gc2V0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGl0bGUodGl0bGUpIHtcbiAgICB0aGlzV2luZG93LlNldFRpdGxlKHRpdGxlKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSB3aW5kb3cgdG8gZnVsbHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEZ1bGxzY3JlZW4oKSB7XG4gICAgdGhpc1dpbmRvdy5GdWxsc2NyZWVuKCk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgdGhpc1dpbmRvdy5TZXRTaXplKHdpZHRoLCBoZWlnaHQpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNpemUoKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuU2l6ZSgpO1xufVxuXG4vKipcbiAqIFNldHMgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIG1heGltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldE1heFNpemUod2lkdGgsIGhlaWdodCkge1xuICAgIHRoaXNXaW5kb3cuU2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCAtIFRoZSBtaW5pbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICB0aGlzV2luZG93LlNldE1pblNpemUod2lkdGgsIGhlaWdodCk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgd2luZG93IHRvIGFsd2F5cyBiZSBvbiB0b3AuXG4gKiBAcGFyYW0ge2Jvb2xlYW59IG9uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBhbHdheXMgYmUgb24gdG9wLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0QWx3YXlzT25Ub3Aob25Ub3ApIHtcbiAgICB0aGlzV2luZG93LlNldEFsd2F5c09uVG9wKG9uVG9wKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSB3aW5kb3cncyBwb3NpdGlvbi5cbiAqIEBwYXJhbSB7bnVtYmVyfSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgd2luZG93J3MgcG9zaXRpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRSZWxhdGl2ZVBvc2l0aW9uKHgsIHkpIHtcbiAgICB0aGlzV2luZG93LlNldFJlbGF0aXZlUG9zaXRpb24oeCwgeSk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbGF0aXZlUG9zaXRpb24oKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuUmVsYXRpdmVQb3NpdGlvbigpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTY3JlZW4oKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuU2NyZWVuKCk7XG59XG5cbi8qKlxuICogSGlkZXMgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgdGhpc1dpbmRvdy5IaWRlKCk7XG59XG5cbi8qKlxuICogTWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBNYXhpbWlzZSgpIHtcbiAgICB0aGlzV2luZG93Lk1heGltaXNlKCk7XG59XG5cbi8qKlxuICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBVbk1heGltaXNlKCkge1xuICAgIHRoaXNXaW5kb3cuVW5NYXhpbWlzZSgpO1xufVxuXG4vKipcbiAqIFRvZ2dsZXMgdGhlIG1heGltaXNhdGlvbiBvZiB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gVG9nZ2xlTWF4aW1pc2UoKSB7XG4gICAgdGhpc1dpbmRvdy5Ub2dnbGVNYXhpbWlzZSgpO1xufVxuXG4vKipcbiAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gTWluaW1pc2UoKSB7XG4gICAgdGhpc1dpbmRvdy5NaW5pbWlzZSgpO1xufVxuXG4vKipcbiAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gVW5NaW5pbWlzZSgpIHtcbiAgICB0aGlzV2luZG93LlVuTWluaW1pc2UoKTtcbn1cblxuLyoqXG4gKiBSZXN0b3JlcyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVzdG9yZSgpIHtcbiAgICB0aGlzV2luZG93LlJlc3RvcmUoKTtcbn1cblxuLyoqXG4gKiBTaG93cyB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICB0aGlzV2luZG93LlNob3coKTtcbn1cblxuLyoqXG4gKiBDbG9zZXMgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENsb3NlKCkge1xuICAgIHRoaXNXaW5kb3cuQ2xvc2UoKTtcbn1cblxuLyoqXG4gKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHIgLSBUaGUgcmVkIGNvbXBvbmVudCBvZiB0aGUgY29sb3VyLlxuICogQHBhcmFtIHtudW1iZXJ9IGcgLSBUaGUgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXG4gKiBAcGFyYW0ge251bWJlcn0gYiAtIFRoZSBibHVlIGNvbXBvbmVudCBvZiB0aGUgY29sb3VyLlxuICogQHBhcmFtIHtudW1iZXJ9IGEgLSBUaGUgYWxwaGEgY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRCYWNrZ3JvdW5kQ29sb3VyKHIsIGcsIGIsIGEpIHtcbiAgICB0aGlzV2luZG93LlNldEJhY2tncm91bmRDb2xvdXIociwgZywgYiwgYSk7XG59XG5cbi8qKlxuICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICogQHBhcmFtIHtib29sZWFufSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFJlc2l6YWJsZShyZXNpemFibGUpIHtcbiAgICB0aGlzV2luZG93LlNldFJlc2l6YWJsZShyZXNpemFibGUpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXaWR0aCgpIHtcbiAgICByZXR1cm4gdGhpc1dpbmRvdy5XaWR0aCgpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICovXG5leHBvcnQgZnVuY3Rpb24gSGVpZ2h0KCkge1xuICAgIHJldHVybiB0aGlzV2luZG93LkhlaWdodCgpO1xufVxuXG4vKipcbiAqIFpvb21zIGluIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBab29tSW4oKSB7XG4gICAgdGhpc1dpbmRvdy5ab29tSW4oKTtcbn1cblxuLyoqXG4gKiBab29tcyBvdXQgdGhlIHdpbmRvdy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFpvb21PdXQoKSB7XG4gICAgdGhpc1dpbmRvdy5ab29tT3V0KCk7XG59XG5cbi8qKlxuICogUmVzZXRzIHRoZSB6b29tIG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBab29tUmVzZXQoKSB7XG4gICAgdGhpc1dpbmRvdy5ab29tUmVzZXQoKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRab29tTGV2ZWwoKSB7XG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuR2V0Wm9vbUxldmVsKCk7XG59XG5cbi8qKlxuICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICogQHBhcmFtIHtudW1iZXJ9IHpvb21MZXZlbCAtIFRoZSB6b29tIGxldmVsIHRvIHNldC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFpvb21MZXZlbCh6b29tTGV2ZWwpIHtcbiAgICB0aGlzV2luZG93LlNldFpvb21MZXZlbCh6b29tTGV2ZWwpO1xufVxuIiwgIlxuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcbmltcG9ydCB7R2V0fSBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCB7T3BlblVSTH0gZnJvbSBcIi4vYnJvd3NlclwiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBzZW5kLlxuICogQHBhcmFtIHthbnl9IFtkYXRhPW51bGxdIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxuICpcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xuICAgIGxldCBldmVudCA9IG5ldyBXYWlsc0V2ZW50KGV2ZW50TmFtZSwgZGF0YSk7XG4gICAgRW1pdChldmVudCk7XG59XG5cbi8qKlxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gZWxlbWVudHMgd2l0aCBgd21sLWV2ZW50YCBhdHRyaWJ1dGUuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XScpO1xuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIERldGFjaGVkOiBmYWxzZSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZCAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqL1xuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lLCBtZXRob2QpIHtcbiAgICBsZXQgdGFyZ2V0V2luZG93ID0gR2V0KHdpbmRvd05hbWUpO1xuICAgIGxldCBtZXRob2RNYXAgPSBXaW5kb3dNZXRob2RzKHRhcmdldFdpbmRvdyk7XG4gICAgaWYgKCFtZXRob2RNYXAuaGFzKG1ldGhvZCkpIHtcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xuICAgIH1cbiAgICB0cnkge1xuICAgICAgICBtZXRob2RNYXAuZ2V0KG1ldGhvZCkoKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoXCJFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJ1wiICsgbWV0aG9kICsgXCInOiBcIiArIGUpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBBZGRzIHdpbmRvdyBsaXN0ZW5lcnMgZm9yIGVsZW1lbnRzIHdpdGggdGhlICd3bWwtd2luZG93JyBhdHRyaWJ1dGUuXG4gKiBSZW1vdmVzIGFueSBleGlzdGluZyBsaXN0ZW5lcnMgYmVmb3JlIGFkZGluZyBuZXcgb25lcy5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLXdpbmRvd10nKTtcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8ICdjbGljayc7XG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdGFyZ2V0LXdpbmRvdycpIHx8ICcnO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XG4gICAgICAgIH07XG5cbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBBZGRzIGEgbGlzdGVuZXIgdG8gZWxlbWVudHMgd2l0aCB0aGUgJ3dtbC1vcGVudXJsJyBhdHRyaWJ1dGUuXG4gKiBXaGVuIHRoZSBzcGVjaWZpZWQgdHJpZ2dlciBldmVudCBpcyBmaXJlZCBvbiBhbnkgb2YgdGhlc2UgZWxlbWVudHMsXG4gKiB0aGUgbGlzdGVuZXIgd2lsbCBvcGVuIHRoZSBVUkwgc3BlY2lmaWVkIGJ5IHRoZSAnd21sLW9wZW51cmwnIGF0dHJpYnV0ZS5cbiAqIElmIGEgJ3dtbC1jb25maXJtJyBhdHRyaWJ1dGUgaXMgcHJvdmlkZWQsIGEgY29uZmlybWF0aW9uIGRpYWxvZyB3aWxsIGJlIGRpc3BsYXllZCxcbiAqIGFuZCB0aGUgVVJMIHdpbGwgb25seSBiZSBvcGVuZWQgaWYgdGhlIHVzZXIgY29uZmlybXMuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtb3BlbnVybF0nKTtcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtb3BlbnVybCcpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICB2b2lkIE9wZW5VUkwodXJsKTtcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG4vKipcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCkge1xuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpO1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBtYXAgb2YgYWxsIG1ldGhvZHMgaW4gdGhlIGN1cnJlbnQgd2luZG93LlxuICogQHJldHVybnMge01hcH0gLSBBIG1hcCBvZiB3aW5kb3cgbWV0aG9kcy5cbiAqL1xuZnVuY3Rpb24gV2luZG93TWV0aG9kcyh0YXJnZXRXaW5kb3cpIHtcbiAgICAvLyBDcmVhdGUgYSBuZXcgbWFwIHRvIHN0b3JlIG1ldGhvZHNcbiAgICBsZXQgcmVzdWx0ID0gbmV3IE1hcCgpO1xuXG4gICAgLy8gSXRlcmF0ZSBvdmVyIGFsbCBwcm9wZXJ0aWVzIG9mIHRoZSB3aW5kb3cgb2JqZWN0XG4gICAgZm9yIChsZXQgbWV0aG9kIGluIHRhcmdldFdpbmRvdykge1xuICAgICAgICAvLyBDaGVjayBpZiB0aGUgcHJvcGVydHkgaXMgaW5kZWVkIGEgbWV0aG9kIChmdW5jdGlvbilcbiAgICAgICAgaWYodHlwZW9mIHRhcmdldFdpbmRvd1ttZXRob2RdID09PSAnZnVuY3Rpb24nKSB7XG4gICAgICAgICAgICAvLyBBZGQgdGhlIG1ldGhvZCB0byB0aGUgbWFwXG4gICAgICAgICAgICByZXN1bHQuc2V0KG1ldGhvZCwgdGFyZ2V0V2luZG93W21ldGhvZF0pO1xuICAgICAgICB9XG5cbiAgICB9XG4gICAgLy8gUmV0dXJuIHRoZSBtYXAgb2Ygd2luZG93IG1ldGhvZHNcbiAgICByZXR1cm4gcmVzdWx0O1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcbiAqL1xuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQge0V2ZW50VHlwZXN9IGZyb20gXCIuL2V2ZW50X3R5cGVzXCI7XG5leHBvcnQgY29uc3QgVHlwZXMgPSBFdmVudFR5cGVzO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuRXZlbnRzLCAnJyk7XG5jb25zdCBFbWl0TWV0aG9kID0gMDtcbmNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgY29uc3RydWN0b3IobmFtZSwgZGF0YSA9IG51bGwpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzZXR1cCgpIHtcbn1cblxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50KSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIGxldCB0b1JlbW92ZSA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHJldHVybiB0cnVlO1xuICAgICAgICB9KTtcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xuICAgICAgICAgICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3MgLSBUaGUgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdGhlIGNhbGxiYWNrIGNhbiBiZSBjYWxsZWQgZm9yIHRoZSBldmVudC4gT25jZSB0aGUgbWF4aW11bSBudW1iZXIgaXMgcmVhY2hlZCwgdGhlIGNhbGxiYWNrIHdpbGwgbm8gbG9uZ2VyIGJlIGNhbGxlZC5cbiAqXG4gQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCB3aGVuIHRoZSBzcGVjaWZpZWQgZXZlbnQgb2NjdXJzLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZC4gSXQgdGFrZXMgbm8gcGFyYW1ldGVycy5cbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuICovXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lLCBjYWxsYmFjaykgeyByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7IH1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIGV2ZW50IG9jY3Vycy5cbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHsgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7IH1cblxuLyoqXG4gKiBSZW1vdmVzIHRoZSBzcGVjaWZpZWQgbGlzdGVuZXIgZnJvbSB0aGUgZXZlbnQgbGlzdGVuZXJzIGNvbGxlY3Rpb24uXG4gKiBJZiBhbGwgbGlzdGVuZXJzIGZvciB0aGUgZXZlbnQgYXJlIHJlbW92ZWQsIHRoZSBldmVudCBrZXkgaXMgZGVsZXRlZCBmcm9tIHRoZSBjb2xsZWN0aW9uLlxuICpcbiAqIEBwYXJhbSB7T2JqZWN0fSBsaXN0ZW5lciAtIFRoZSBsaXN0ZW5lciB0byBiZSByZW1vdmVkLlxuICovXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG59XG5cblxuLyoqXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxuICogQHBhcmFtIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzIC0gQWRkaXRpb25hbCBldmVudCBuYW1lcyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcbn1cbi8qKlxuICogUmVtb3ZlcyBhbGwgZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBmdW5jdGlvbiBPZmZBbGxcbiAqIEByZXR1cm5zIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkgeyBldmVudExpc3RlbmVycy5jbGVhcigpOyB9XG5cbi8qKlxuICogRW1pdHMgYW4gZXZlbnQgdXNpbmcgdGhlIGdpdmVuIGV2ZW50IG5hbWUuXG4gKlxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxuICogQHJldHVybnMge2FueX0gLSBUaGUgcmVzdWx0IG9mIHRoZSBlbWl0dGVkIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkgeyByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCBldmVudCk7IH1cbiIsICJcbmV4cG9ydCBjb25zdCBFdmVudFR5cGVzID0ge1xuXHRXaW5kb3dzOiB7XG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0QVBNUG93ZXJTdGF0dXNDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclN0YXR1c0NoYW5nZVwiLFxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwid2luZG93czpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZDogXCJ3aW5kb3dzOldlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkXCIsXG5cdFx0V2luZG93SW5hY3RpdmU6IFwid2luZG93czpXaW5kb3dJbmFjdGl2ZVwiLFxuXHRcdFdpbmRvd0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Q2xvc2U6IFwid2luZG93czpXaW5kb3dDbG9zZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHR9LFxuXHRNYWM6IHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZU5vdGlmaWNhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVOb3RpZmljYXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZSFcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW4hXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd0RpZERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRFbmRTaGVldDogXCJtYWM6V2luZG93RGlkRW5kU2hlZXRcIixcblx0XHRXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4cG9zZTogXCJtYWM6V2luZG93RGlkRXhwb3NlXCIsXG5cdFx0V2luZG93RGlkRm9jdXM6IFwibWFjOldpbmRvd0RpZEZvY3VzXCIsXG5cdFx0V2luZG93RGlkTWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZE1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJtYWM6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT25TY3JlZW5cIixcblx0XHRXaW5kb3dEaWRSZXNpZ25LZXk6IFwibWFjOldpbmRvd0RpZFJlc2lnbktleVwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbk1haW46IFwibWFjOldpbmRvd0RpZFJlc2lnbk1haW5cIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibWFjOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZTogXCJtYWM6V2luZG93RGlkVXBkYXRlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUFscGhhXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVNoYWRvd1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5XCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlIVwiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdE1lbnVXaWxsT3BlbjogXCJtYWM6TWVudVdpbGxPcGVuXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWREaXNwbGF5SXRlbTogXCJtYWM6TWVudURpZERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVEaWRSZW1vdmVJdGVtOiBcIm1hYzpNZW51RGlkUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsQmVnaW5UcmFja2luZzogXCJtYWM6TWVudVdpbGxCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51RGlkUG9wVXA6IFwibWFjOk1lbnVEaWRQb3BVcFwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0fSxcblx0TGludXg6IHtcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdH0sXG5cdENvbW1vbjoge1xuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwiY29tbW9uOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJjb21tb246V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJjb21tb246V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd1pvb206IFwiY29tbW9uOldpbmRvd1pvb21cIixcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwiY29tbW9uOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwiY29tbW9uOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJjb21tb246V2luZG93TG9zdEZvY3VzXCIsXG5cdFx0V2luZG93U2hvdzogXCJjb21tb246V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd0hpZGU6IFwiY29tbW9uOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcImNvbW1vbjpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcblx0XHRXaW5kb3dSdW50aW1lUmVhZHk6IFwiY29tbW9uOldpbmRvd1J1bnRpbWVSZWFkeVwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXG5cdH0sXG59O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gT3BlbkZpbGVEaWFsb2dPcHRpb25zXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ3JlYXRlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbU2hvd0hpZGRlbkZpbGVzXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c011bHRpcGxlU2VsZWN0aW9uXSAtIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0hpZGVFeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c090aGVyRmlsZXR5cGVzXSAtIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGl0bGUgb2YgdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTWVzc2FnZV0gLSBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlyZWN0b3J5XSAtIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTYXZlRmlsZURpYWxvZ09wdGlvbnNcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRmlsZW5hbWVdIC0gRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DcmVhdGVEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtTaG93SGlkZGVuRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtSZXNvbHZlc0FsaWFzZXNdIC0gSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb25dIC0gSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSGlkZUV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzT3RoZXJGaWxldHlwZXNdIC0gSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuXG4gKiBAcHJvcGVydHkge0ZpbGVGaWx0ZXJbXX0gW0ZpbHRlcnNdIC0gQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaXRsZSBvZiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtNZXNzYWdlXSAtIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtCdXR0b25UZXh0XSAtIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtEaXJlY3RvcnldIC0gRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0RldGFjaGVkXSAtIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge0J1dHRvbltdfSBbQnV0dG9uc10gLSBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBCdXR0b25cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTGFiZWxdIC0gVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtJc0NhbmNlbF0gLSBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNEZWZhdWx0XSAtIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IEZpbGVGaWx0ZXJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlzcGxheU5hbWVdIC0gRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1BhdHRlcm5dIC0gUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy5cbiAqL1xuXG4vLyBzZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpYWxvZ0Vycm9yQ2FsbGJhY2sgPSBkaWFsb2dFcnJvckNhbGxiYWNrO1xud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZywgJycpO1xuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBpZCB0aGF0IGlzIG5vdCBwcmVzZW50IGluIGRpYWxvZ1Jlc3BvbnNlcy5cbiAqIEByZXR1cm5zIHtzdHJpbmd9IHVuaXF1ZSBpZFxuICovXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogU2hvd3MgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7bnVtYmVyfSB0eXBlIC0gdHlwZSBvZiBkaWFsb2dcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN8T3BlbkZpbGVEaWFsb2dPcHRpb25zfFNhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIG9wdGlvbnMgZm9yIHRoZSBkaWFsb2dcbiAqIEByZXR1cm5zIHtQcm9taXNlfSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zID0ge30pIHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBjYWxsYmFjayBmcm9tIGEgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBJRCBvZiB0aGUgZGlhbG9nIHJlc3BvbnNlLlxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgZGF0YSByZWNlaXZlZCBmcm9tIHRoZSBkaWFsb2cuXG4gKiBAcGFyYW0ge2Jvb2xlYW59IGlzSlNPTiAtIEZsYWcgaW5kaWNhdGluZyB3aGV0aGVyIHRoZSBkYXRhIGlzIGluIEpTT04gZm9ybWF0LlxuICpcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cbiAqL1xuZnVuY3Rpb24gZGlhbG9nUmVzdWx0Q2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBDYWxsYmFjayBmdW5jdGlvbiBmb3IgaGFuZGxpbmcgZXJyb3JzIGluIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuXG4vLyBSZXBsYWNlIGBtZXRob2RzYCB3aXRoIGNvbnN0YW50cyBpbiBUaXRsZSBDYXNlXG5cbi8qKlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGNvbnN0IEluZm8gPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpO1xuXG4vKipcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBjb25zdCBXYXJuaW5nID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dXYXJuaW5nLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgY29uc3QgRXJyb3IgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgY29uc3QgUXVlc3Rpb24gPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge09wZW5GaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcy4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGNvbnN0IE9wZW5GaWxlID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dPcGVuRmlsZSwgb3B0aW9ucyk7XG5cbi8qKlxuICogQHBhcmFtIHtTYXZlRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZS4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGNvbnN0IFNhdmVGaWxlID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcbndpbmRvdy5fd2FpbHMuY2FsbEVycm9ySGFuZGxlciA9IGVycm9ySGFuZGxlcjtcblxuXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYWxsLCAnJyk7XG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxuICpcbiAqIEByZXR1cm4ge3N0cmluZ30gLSBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBjYWxsIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxuICogQHBhcmFtIHtib29sZWFufSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cbiAqXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9IC0gVGhpcyBtZXRob2QgZG9lcyBub3QgcmV0dXJuIGFueSB2YWx1ZS5cbiAqL1xuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XG4gICAgaWYgKHByb21pc2VIYW5kbGVyKSB7XG4gICAgICAgIHByb21pc2VIYW5kbGVyLnJlc29sdmUoaXNKU09OID8gSlNPTi5wYXJzZShkYXRhKSA6IGRhdGEpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBlcnJvciBmcm9tIGEgY2FsbCByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZSB0byByZWplY3QgdGhlIHByb21pc2UgaGFuZGxlciB3aXRoLlxuICpcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZCwgbWVzc2FnZSkge1xuICAgIGNvbnN0IHByb21pc2VIYW5kbGVyID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xuICAgIGlmIChwcm9taXNlSGFuZGxlcikge1xuICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobWVzc2FnZSk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHJpZXZlcyBhbmQgcmVtb3ZlcyB0aGUgcmVzcG9uc2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCBmcm9tIHRoZSBjYWxsUmVzcG9uc2VzIG1hcC5cbiAqXG4gKiBAcGFyYW0ge2FueX0gaWQgLSBUaGUgSUQgb2YgdGhlIHJlc3BvbnNlIHRvIGJlIHJldHJpZXZlZCBhbmQgcmVtb3ZlZC5cbiAqXG4gKiBAcmV0dXJucyB7YW55fSBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQuXG4gKi9cbmZ1bmN0aW9uIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKSB7XG4gICAgY29uc3QgcmVzcG9uc2UgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIHJldHVybiByZXNwb25zZTtcbn1cblxuLyoqXG4gKiBFeGVjdXRlcyBhIGNhbGwgdXNpbmcgdGhlIHByb3ZpZGVkIHR5cGUgYW5kIG9wdGlvbnMuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd8bnVtYmVyfSB0eXBlIC0gVGhlIHR5cGUgb2YgY2FsbCB0byBleGVjdXRlLlxuICogQHBhcmFtIHtPYmplY3R9IFtvcHRpb25zPXt9XSAtIEFkZGl0aW9uYWwgb3B0aW9ucyBmb3IgdGhlIGNhbGwuXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgYmFzZWQgb24gdGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cbiAqL1xuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucyA9IHt9KSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmUsIHJlamVjdCB9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5jYXRjaCgoZXJyb3IpID0+IHtcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XG4gICAgICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pO1xuICAgIH0pO1xufVxuXG4vKipcbiAqIENhbGwgbWV0aG9kLlxuICpcbiAqIEBwYXJhbSB7T2JqZWN0fSBvcHRpb25zIC0gVGhlIG9wdGlvbnMgZm9yIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJucyB7T2JqZWN0fSAtIFRoZSByZXN1bHQgb2YgdGhlIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIEV4ZWN1dGVzIGEgbWV0aG9kIGJ5IG5hbWUuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXG4gKiBAcGFyYW0gey4uLip9IGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEB0aHJvd3Mge0Vycm9yfSBJZiB0aGUgbmFtZSBpcyBub3QgYSBzdHJpbmcgb3IgaXMgbm90IGluIHRoZSBjb3JyZWN0IGZvcm1hdC5cbiAqIEByZXR1cm5zIHsqfSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgZXhlY3V0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG5hbWUsIC4uLmFyZ3MpIHtcbiAgICBpZiAodHlwZW9mIG5hbWUgIT09IFwic3RyaW5nXCIgfHwgbmFtZS5zcGxpdChcIi5cIikubGVuZ3RoICE9PSAzKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIkNhbGxCeU5hbWUgcmVxdWlyZXMgYSBzdHJpbmcgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJ1wiKTtcbiAgICB9XG4gICAgbGV0IFtwYWNrYWdlTmFtZSwgc3RydWN0TmFtZSwgbWV0aG9kTmFtZV0gPSBuYW1lLnNwbGl0KFwiLlwiKTtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIHtcbiAgICAgICAgcGFja2FnZU5hbWUsXG4gICAgICAgIHN0cnVjdE5hbWUsXG4gICAgICAgIG1ldGhvZE5hbWUsXG4gICAgICAgIGFyZ3NcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBieSBpdHMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqXG4gKiBAcGFyYW0ge251bWJlcn0gbWV0aG9kSUQgLSBUaGUgSUQgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlJRChtZXRob2RJRCwgLi4uYXJncykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xuICAgICAgICBtZXRob2RJRCxcbiAgICAgICAgYXJnc1xuICAgIH0pO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgcGx1Z2luLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHBsdWdpbi5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJucyB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFBsdWdpbihwbHVnaW5OYW1lLCBtZXRob2ROYW1lLCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XG4gICAgICAgIHBhY2thZ2VOYW1lOiBcIndhaWxzLXBsdWdpbnNcIixcbiAgICAgICAgc3RydWN0TmFtZTogcGx1Z2luTmFtZSxcbiAgICAgICAgbWV0aG9kTmFtZSxcbiAgICAgICAgYXJnc1xuICAgIH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge2RlYnVnTG9nfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvbG9nXCI7XG5cbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xuXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb25cIjtcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXJcIjtcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkXCI7XG5pbXBvcnQgKiBhcyBDb250ZXh0TWVudSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnVcIjtcbmltcG9ydCAqIGFzIERyYWcgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWdcIjtcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9mbGFnc1wiO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvc2NyZWVuc1wiO1xuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9zeXN0ZW1cIjtcbmltcG9ydCAqIGFzIFdpbmRvdyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd2luZG93XCI7XG5pbXBvcnQgKiBhcyBXTUwgZnJvbSAnLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sJztcbmltcG9ydCAqIGFzIEV2ZW50cyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzXCI7XG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9kaWFsb2dzXCI7XG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxsc1wiO1xuaW1wb3J0IHtpbnZva2V9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9zeXN0ZW1cIjtcblxuLyoqKlxuIFRoaXMgdGVjaG5pcXVlIGZvciBwcm9wZXIgbG9hZCBkZXRlY3Rpb24gaXMgdGFrZW4gZnJvbSBIVE1YOlxuXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcblxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXG5cbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcblxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxuXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cbiBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlIGRpc3RyaWJ1dGlvbi5cblxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXG4gQU5EIEFOWSBFWFBSRVNTIE9SIElNUExJRUQgV0FSUkFOVElFUywgSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFRIRVxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcbiBGT1IgQU5ZIERJUkVDVCwgSU5ESVJFQ1QsIElOQ0lERU5UQUwsIFNQRUNJQUwsIEVYRU1QTEFSWSwgT1IgQ09OU0VRVUVOVElBTFxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcbiBDQVVTRUQgQU5EIE9OIEFOWSBUSEVPUlkgT0YgTElBQklMSVRZLCBXSEVUSEVSIElOIENPTlRSQUNULCBTVFJJQ1QgTElBQklMSVRZLFxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cblxuICoqKi9cblxud2luZG93Ll93YWlscy5pbnZva2U9aW52b2tlO1xuXG53aW5kb3cud2FpbHMgPSB3aW5kb3cud2FpbHMgfHwge307XG53aW5kb3cud2FpbHMuQXBwbGljYXRpb24gPSBBcHBsaWNhdGlvbjtcbndpbmRvdy53YWlscy5Ccm93c2VyID0gQnJvd3NlcjtcbndpbmRvdy53YWlscy5DYWxsID0gQ2FsbDtcbndpbmRvdy53YWlscy5DbGlwYm9hcmQgPSBDbGlwYm9hcmQ7XG53aW5kb3cud2FpbHMuRGlhbG9ncyA9IERpYWxvZ3M7XG53aW5kb3cud2FpbHMuRXZlbnRzID0gRXZlbnRzO1xud2luZG93LndhaWxzLkZsYWdzID0gRmxhZ3M7XG53aW5kb3cud2FpbHMuU2NyZWVucyA9IFNjcmVlbnM7XG53aW5kb3cud2FpbHMuU3lzdGVtID0gU3lzdGVtO1xud2luZG93LndhaWxzLldpbmRvdyA9IFdpbmRvdztcbndpbmRvdy53YWlscy5XTUwgPSBXTUw7XG5cblxubGV0IGlzUmVhZHkgPSBmYWxzZVxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGZ1bmN0aW9uKCkge1xuICAgIGlzUmVhZHkgPSB0cnVlXG4gICAgd2luZG93Ll93YWlscy5pbnZva2UoJ3dhaWxzOnJ1bnRpbWU6cmVhZHknKTtcbiAgICBpZihERUJVRykge1xuICAgICAgICBkZWJ1Z0xvZyhcIldhaWxzIFJ1bnRpbWUgTG9hZGVkXCIpO1xuICAgIH1cbn0pXG5cbmZ1bmN0aW9uIHdoZW5SZWFkeShmbikge1xuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcbiAgICAgICAgZm4oKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgZm4pO1xuICAgIH1cbn1cblxud2hlblJlYWR5KCgpID0+IHtcbiAgICBXTUwuUmVsb2FkKCk7XG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBS08sV0FBUyxTQUFTLFNBQVM7QUFFOUIsWUFBUTtBQUFBLE1BQ0osa0JBQWtCLFVBQVU7QUFBQSxNQUM1QjtBQUFBLE1BQ0E7QUFBQSxJQUNKO0FBQUEsRUFDSjs7O0FDWkE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNBQSxNQUFJLGNBQ0Y7QUFXSyxNQUFJLFNBQVMsQ0FBQ0EsUUFBTyxPQUFPO0FBQ2pDLFFBQUksS0FBSztBQUNULFFBQUksSUFBSUE7QUFDUixXQUFPLEtBQUs7QUFDVixZQUFNLFlBQWEsS0FBSyxPQUFPLElBQUksS0FBTSxDQUFDO0FBQUEsSUFDNUM7QUFDQSxXQUFPO0FBQUEsRUFDVDs7O0FDTkEsTUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBR3JDLE1BQU0sY0FBYztBQUFBLElBQ3ZCLE1BQU07QUFBQSxJQUNOLFdBQVc7QUFBQSxJQUNYLGFBQWE7QUFBQSxJQUNiLFFBQVE7QUFBQSxJQUNSLGFBQWE7QUFBQSxJQUNiLFFBQVE7QUFBQSxJQUNSLFFBQVE7QUFBQSxJQUNSLFNBQVM7QUFBQSxJQUNULFFBQVE7QUFBQSxJQUNSLFNBQVM7QUFBQSxFQUNiO0FBQ08sTUFBSSxXQUFXLE9BQU87QUFzQnRCLFdBQVMsdUJBQXVCLFFBQVEsWUFBWTtBQUN2RCxXQUFPLFNBQVUsUUFBUSxPQUFLLE1BQU07QUFDaEMsYUFBTyxrQkFBa0IsUUFBUSxRQUFRLFlBQVksSUFBSTtBQUFBLElBQzdEO0FBQUEsRUFDSjtBQXFDQSxXQUFTLGtCQUFrQixVQUFVLFFBQVEsWUFBWSxNQUFNO0FBQzNELFFBQUksTUFBTSxJQUFJLElBQUksVUFBVTtBQUM1QixRQUFJLGFBQWEsT0FBTyxVQUFVLFFBQVE7QUFDMUMsUUFBSSxhQUFhLE9BQU8sVUFBVSxNQUFNO0FBQ3hDLFFBQUksZUFBZTtBQUFBLE1BQ2YsU0FBUyxDQUFDO0FBQUEsSUFDZDtBQUNBLFFBQUksWUFBWTtBQUNaLG1CQUFhLFFBQVEscUJBQXFCLElBQUk7QUFBQSxJQUNsRDtBQUNBLFFBQUksTUFBTTtBQUNOLFVBQUksYUFBYSxPQUFPLFFBQVEsS0FBSyxVQUFVLElBQUksQ0FBQztBQUFBLElBQ3hEO0FBQ0EsaUJBQWEsUUFBUSxtQkFBbUIsSUFBSTtBQUM1QyxXQUFPLElBQUksUUFBUSxDQUFDLFNBQVMsV0FBVztBQUNwQyxZQUFNLEtBQUssWUFBWSxFQUNsQixLQUFLLGNBQVk7QUFDZCxZQUFJLFNBQVMsSUFBSTtBQUViLGNBQUksU0FBUyxRQUFRLElBQUksY0FBYyxLQUFLLFNBQVMsUUFBUSxJQUFJLGNBQWMsRUFBRSxRQUFRLGtCQUFrQixNQUFNLElBQUk7QUFDakgsbUJBQU8sU0FBUyxLQUFLO0FBQUEsVUFDekIsT0FBTztBQUNILG1CQUFPLFNBQVMsS0FBSztBQUFBLFVBQ3pCO0FBQUEsUUFDSjtBQUNBLGVBQU8sTUFBTSxTQUFTLFVBQVUsQ0FBQztBQUFBLE1BQ3JDLENBQUMsRUFDQSxLQUFLLFVBQVEsUUFBUSxJQUFJLENBQUMsRUFDMUIsTUFBTSxXQUFTLE9BQU8sS0FBSyxDQUFDO0FBQUEsSUFDckMsQ0FBQztBQUFBLEVBQ0w7OztBRjVHQSxNQUFNLE9BQU8sdUJBQXVCLFlBQVksYUFBYSxFQUFFO0FBRS9ELE1BQU0sYUFBYTtBQUNuQixNQUFNLGFBQWE7QUFDbkIsTUFBTSxhQUFhO0FBUVosV0FBUyxPQUFPO0FBQ25CLFdBQU8sS0FBSyxVQUFVO0FBQUEsRUFDMUI7QUFPTyxXQUFTLE9BQU87QUFDbkIsV0FBTyxLQUFLLFVBQVU7QUFBQSxFQUMxQjtBQU9PLFdBQVMsT0FBTztBQUNuQixXQUFPLEtBQUssVUFBVTtBQUFBLEVBQzFCOzs7QUc3Q0E7QUFBQTtBQUFBO0FBQUE7QUFhQSxNQUFNQyxRQUFPLHVCQUF1QixZQUFZLFNBQVMsRUFBRTtBQUMzRCxNQUFNLGlCQUFpQjtBQU9oQixXQUFTLFFBQVEsS0FBSztBQUN6QixXQUFPQSxNQUFLLGdCQUFnQixFQUFDLElBQUcsQ0FBQztBQUFBLEVBQ3JDOzs7QUN2QkE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWNBLE1BQU1DLFFBQU8sdUJBQXVCLFlBQVksV0FBVyxFQUFFO0FBQzdELE1BQU0sbUJBQW1CO0FBQ3pCLE1BQU0sZ0JBQWdCO0FBUWYsV0FBUyxRQUFRLE1BQU07QUFDMUIsV0FBT0EsTUFBSyxrQkFBa0IsRUFBQyxLQUFJLENBQUM7QUFBQSxFQUN4QztBQU1PLFdBQVMsT0FBTztBQUNuQixXQUFPQSxNQUFLLGFBQWE7QUFBQSxFQUM3Qjs7O0FDbENBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFhQSxNQUFJQyxRQUFPLHVCQUF1QixZQUFZLFFBQVEsRUFBRTtBQUN4RCxNQUFNLG1CQUFtQjtBQUN6QixNQUFNLGNBQWM7QUFFYixXQUFTLE9BQU8sS0FBSztBQUN4QixRQUFHLE9BQU8sUUFBUTtBQUNkLGFBQU8sT0FBTyxPQUFPLFFBQVEsWUFBWSxHQUFHO0FBQUEsSUFDaEQ7QUFDQSxXQUFPLE9BQU8sT0FBTyxnQkFBZ0IsU0FBUyxZQUFZLEdBQUc7QUFBQSxFQUNqRTtBQU9PLFdBQVMsYUFBYTtBQUN6QixXQUFPQSxNQUFLLGdCQUFnQjtBQUFBLEVBQ2hDO0FBU08sV0FBUyxlQUFlO0FBQzNCLFFBQUksV0FBVyxNQUFNLHFCQUFxQjtBQUMxQyxXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCO0FBYU8sV0FBUyxjQUFjO0FBQzFCLFdBQU9BLE1BQUssV0FBVztBQUFBLEVBQzNCO0FBT08sV0FBUyxZQUFZO0FBQ3hCLFdBQU8sT0FBTyxPQUFPLFlBQVksT0FBTztBQUFBLEVBQzVDO0FBT08sV0FBUyxVQUFVO0FBQ3RCLFdBQU8sT0FBTyxPQUFPLFlBQVksT0FBTztBQUFBLEVBQzVDO0FBT08sV0FBUyxRQUFRO0FBQ3BCLFdBQU8sT0FBTyxPQUFPLFlBQVksT0FBTztBQUFBLEVBQzVDO0FBTU8sV0FBUyxVQUFVO0FBQ3RCLFdBQU8sT0FBTyxPQUFPLFlBQVksU0FBUztBQUFBLEVBQzlDO0FBT08sV0FBUyxRQUFRO0FBQ3BCLFdBQU8sT0FBTyxPQUFPLFlBQVksU0FBUztBQUFBLEVBQzlDO0FBT08sV0FBUyxVQUFVO0FBQ3RCLFdBQU8sT0FBTyxPQUFPLFlBQVksU0FBUztBQUFBLEVBQzlDO0FBRU8sV0FBUyxVQUFVO0FBQ3RCLFdBQU8sT0FBTyxPQUFPLFlBQVksVUFBVTtBQUFBLEVBQy9DOzs7QUNuR0EsU0FBTyxpQkFBaUIsZUFBZSxrQkFBa0I7QUFFekQsTUFBTUMsUUFBTyx1QkFBdUIsWUFBWSxhQUFhLEVBQUU7QUFDL0QsTUFBTSxrQkFBa0I7QUFFeEIsV0FBUyxnQkFBZ0IsSUFBSSxHQUFHLEdBQUcsTUFBTTtBQUNyQyxTQUFLQSxNQUFLLGlCQUFpQixFQUFDLElBQUksR0FBRyxHQUFHLEtBQUksQ0FBQztBQUFBLEVBQy9DO0FBRUEsV0FBUyxtQkFBbUIsT0FBTztBQUUvQixRQUFJLFVBQVUsTUFBTTtBQUNwQixRQUFJLG9CQUFvQixPQUFPLGlCQUFpQixPQUFPLEVBQUUsaUJBQWlCLHNCQUFzQjtBQUNoRyx3QkFBb0Isb0JBQW9CLGtCQUFrQixLQUFLLElBQUk7QUFDbkUsUUFBSSxtQkFBbUI7QUFDbkIsWUFBTSxlQUFlO0FBQ3JCLFVBQUksd0JBQXdCLE9BQU8saUJBQWlCLE9BQU8sRUFBRSxpQkFBaUIsMkJBQTJCO0FBQ3pHLHNCQUFnQixtQkFBbUIsTUFBTSxTQUFTLE1BQU0sU0FBUyxxQkFBcUI7QUFDdEY7QUFBQSxJQUNKO0FBRUEsOEJBQTBCLEtBQUs7QUFBQSxFQUNuQztBQVVBLFdBQVMsMEJBQTBCLE9BQU87QUFHdEMsUUFBSSxRQUFRLEdBQUc7QUFDWDtBQUFBLElBQ0o7QUFHQSxVQUFNLFVBQVUsTUFBTTtBQUN0QixVQUFNLGdCQUFnQixPQUFPLGlCQUFpQixPQUFPO0FBQ3JELFVBQU0sMkJBQTJCLGNBQWMsaUJBQWlCLHVCQUF1QixFQUFFLEtBQUs7QUFDOUYsWUFBUSwwQkFBMEI7QUFBQSxNQUM5QixLQUFLO0FBQ0Q7QUFBQSxNQUNKLEtBQUs7QUFDRCxjQUFNLGVBQWU7QUFDckI7QUFBQSxNQUNKO0FBRUksWUFBSSxRQUFRLG1CQUFtQjtBQUMzQjtBQUFBLFFBQ0o7QUFHQSxjQUFNLFlBQVksT0FBTyxhQUFhO0FBQ3RDLGNBQU0sZUFBZ0IsVUFBVSxTQUFTLEVBQUUsU0FBUztBQUNwRCxZQUFJLGNBQWM7QUFDZCxtQkFBUyxJQUFJLEdBQUcsSUFBSSxVQUFVLFlBQVksS0FBSztBQUMzQyxrQkFBTSxRQUFRLFVBQVUsV0FBVyxDQUFDO0FBQ3BDLGtCQUFNLFFBQVEsTUFBTSxlQUFlO0FBQ25DLHFCQUFTLElBQUksR0FBRyxJQUFJLE1BQU0sUUFBUSxLQUFLO0FBQ25DLG9CQUFNLE9BQU8sTUFBTSxDQUFDO0FBQ3BCLGtCQUFJLFNBQVMsaUJBQWlCLEtBQUssTUFBTSxLQUFLLEdBQUcsTUFBTSxTQUFTO0FBQzVEO0FBQUEsY0FDSjtBQUFBLFlBQ0o7QUFBQSxVQUNKO0FBQUEsUUFDSjtBQUVBLFlBQUksUUFBUSxZQUFZLFdBQVcsUUFBUSxZQUFZLFlBQVk7QUFDL0QsY0FBSSxnQkFBaUIsQ0FBQyxRQUFRLFlBQVksQ0FBQyxRQUFRLFVBQVc7QUFDMUQ7QUFBQSxVQUNKO0FBQUEsUUFDSjtBQUdBLGNBQU0sZUFBZTtBQUFBLElBQzdCO0FBQUEsRUFDSjs7O0FDaEdBO0FBQUE7QUFBQTtBQUFBO0FBa0JPLFdBQVMsUUFBUSxXQUFXO0FBQy9CLFFBQUk7QUFDQSxhQUFPLE9BQU8sT0FBTyxNQUFNLFNBQVM7QUFBQSxJQUN4QyxTQUFTLEdBQUc7QUFDUixZQUFNLElBQUksTUFBTSw4QkFBOEIsWUFBWSxRQUFRLENBQUM7QUFBQSxJQUN2RTtBQUFBLEVBQ0o7OztBQ1JBLFNBQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQUNsQyxTQUFPLE9BQU8sZUFBZTtBQUM3QixTQUFPLE9BQU8sVUFBVTtBQUN4QixTQUFPLGlCQUFpQixhQUFhLFdBQVc7QUFDaEQsU0FBTyxpQkFBaUIsYUFBYSxXQUFXO0FBQ2hELFNBQU8saUJBQWlCLFdBQVcsU0FBUztBQUc1QyxNQUFJLGFBQWE7QUFDakIsTUFBSSxhQUFhO0FBQ2pCLE1BQUksWUFBWTtBQUNoQixNQUFJLGdCQUFnQjtBQUVwQixXQUFTLFNBQVMsR0FBRztBQUNqQixRQUFJLE1BQU0sT0FBTyxpQkFBaUIsRUFBRSxNQUFNLEVBQUUsaUJBQWlCLHFCQUFxQjtBQUNsRixRQUFJLENBQUMsT0FBTyxRQUFRLE1BQU0sSUFBSSxLQUFLLE1BQU0sVUFBVSxFQUFFLFlBQVksR0FBRztBQUNoRSxhQUFPO0FBQUEsSUFDWDtBQUNBLFdBQU8sRUFBRSxXQUFXO0FBQUEsRUFDeEI7QUFFQSxXQUFTLGFBQWEsT0FBTztBQUN6QixnQkFBWTtBQUFBLEVBQ2hCO0FBRUEsV0FBUyxVQUFVO0FBQ2YsYUFBUyxLQUFLLE1BQU0sU0FBUztBQUM3QixpQkFBYTtBQUFBLEVBQ2pCO0FBRUEsV0FBUyxhQUFhO0FBQ2xCLFFBQUksWUFBYTtBQUNiLGFBQU8sVUFBVSxVQUFVLEVBQUU7QUFDN0IsYUFBTztBQUFBLElBQ1g7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUVBLFdBQVMsWUFBWSxHQUFHO0FBQ3BCLFFBQUcsVUFBVSxLQUFLLFdBQVcsS0FBSyxTQUFTLENBQUMsR0FBRztBQUMzQyxtQkFBYSxDQUFDLENBQUMsWUFBWSxDQUFDO0FBQUEsSUFDaEM7QUFBQSxFQUNKO0FBRUEsV0FBUyxZQUFZLEdBQUc7QUFFcEIsV0FBTyxFQUFFLEVBQUUsVUFBVSxFQUFFLE9BQU8sZUFBZSxFQUFFLFVBQVUsRUFBRSxPQUFPO0FBQUEsRUFDdEU7QUFFQSxXQUFTLFVBQVUsR0FBRztBQUNsQixRQUFJLGVBQWUsRUFBRSxZQUFZLFNBQVksRUFBRSxVQUFVLEVBQUU7QUFDM0QsUUFBSSxlQUFlLEdBQUc7QUFDbEIsY0FBUTtBQUFBLElBQ1o7QUFBQSxFQUNKO0FBRUEsV0FBUyxVQUFVLFNBQVMsZUFBZTtBQUN2QyxhQUFTLGdCQUFnQixNQUFNLFNBQVM7QUFDeEMsaUJBQWE7QUFBQSxFQUNqQjtBQUVBLFdBQVMsWUFBWSxHQUFHO0FBQ3BCLGlCQUFhLFVBQVUsQ0FBQztBQUN4QixRQUFJLFVBQVUsS0FBSyxXQUFXO0FBQzFCLG1CQUFhLENBQUM7QUFBQSxJQUNsQjtBQUFBLEVBQ0o7QUFFQSxXQUFTLFVBQVUsR0FBRztBQUNsQixRQUFJLGVBQWUsRUFBRSxZQUFZLFNBQVksRUFBRSxVQUFVLEVBQUU7QUFDM0QsUUFBRyxjQUFjLGVBQWUsR0FBRztBQUMvQixhQUFPLE1BQU07QUFDYixhQUFPO0FBQUEsSUFDWDtBQUNBLFdBQU87QUFBQSxFQUNYO0FBRUEsV0FBUyxhQUFhLEdBQUc7QUFDckIsUUFBSSxxQkFBcUIsUUFBUSwyQkFBMkIsS0FBSztBQUNqRSxRQUFJLG9CQUFvQixRQUFRLDBCQUEwQixLQUFLO0FBRy9ELFFBQUksY0FBYyxRQUFRLG1CQUFtQixLQUFLO0FBRWxELFFBQUksY0FBYyxPQUFPLGFBQWEsRUFBRSxVQUFVO0FBQ2xELFFBQUksYUFBYSxFQUFFLFVBQVU7QUFDN0IsUUFBSSxZQUFZLEVBQUUsVUFBVTtBQUM1QixRQUFJLGVBQWUsT0FBTyxjQUFjLEVBQUUsVUFBVTtBQUdwRCxRQUFJLGNBQWMsT0FBTyxhQUFhLEVBQUUsVUFBVyxvQkFBb0I7QUFDdkUsUUFBSSxhQUFhLEVBQUUsVUFBVyxvQkFBb0I7QUFDbEQsUUFBSSxZQUFZLEVBQUUsVUFBVyxxQkFBcUI7QUFDbEQsUUFBSSxlQUFlLE9BQU8sY0FBYyxFQUFFLFVBQVcscUJBQXFCO0FBRzFFLFFBQUksQ0FBQyxjQUFjLENBQUMsZUFBZSxDQUFDLGFBQWEsQ0FBQyxnQkFBZ0IsZUFBZSxRQUFXO0FBQ3hGLGdCQUFVO0FBQUEsSUFDZCxXQUVTLGVBQWU7QUFBYyxnQkFBVSxXQUFXO0FBQUEsYUFDbEQsY0FBYztBQUFjLGdCQUFVLFdBQVc7QUFBQSxhQUNqRCxjQUFjO0FBQVcsZ0JBQVUsV0FBVztBQUFBLGFBQzlDLGFBQWE7QUFBYSxnQkFBVSxXQUFXO0FBQUEsYUFDL0M7QUFBWSxnQkFBVSxVQUFVO0FBQUEsYUFDaEM7QUFBVyxnQkFBVSxVQUFVO0FBQUEsYUFDL0I7QUFBYyxnQkFBVSxVQUFVO0FBQUEsYUFDbEM7QUFBYSxnQkFBVSxVQUFVO0FBQUEsRUFDOUM7OztBQzVIQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUF1REEsTUFBTUMsUUFBTyx1QkFBdUIsWUFBWSxTQUFTLEVBQUU7QUFFM0QsTUFBTSxTQUFTO0FBQ2YsTUFBTSxhQUFhO0FBQ25CLE1BQU0sYUFBYTtBQU1aLFdBQVMsU0FBUztBQUNyQixXQUFPQSxNQUFLLE1BQU07QUFBQSxFQUN0QjtBQUtPLFdBQVMsYUFBYTtBQUN6QixXQUFPQSxNQUFLLFVBQVU7QUFBQSxFQUMxQjtBQU1PLFdBQVMsYUFBYTtBQUN6QixXQUFPQSxNQUFLLFVBQVU7QUFBQSxFQUMxQjs7O0FDbEZBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxnQkFBQUM7QUFBQSxJQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxnQkFBQUM7QUFBQSxJQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQW1CQSxNQUFNLFNBQVM7QUFDZixNQUFNLFdBQVc7QUFDakIsTUFBTSxhQUFhO0FBQ25CLE1BQU0sZUFBZTtBQUNyQixNQUFNLFVBQVU7QUFDaEIsTUFBTSxPQUFPO0FBQ2IsTUFBTSxhQUFhO0FBQ25CLE1BQU0sYUFBYTtBQUNuQixNQUFNLGlCQUFpQjtBQUN2QixNQUFNLHNCQUFzQjtBQUM1QixNQUFNLG1CQUFtQjtBQUN6QixNQUFNLFNBQVM7QUFDZixNQUFNLE9BQU87QUFDYixNQUFNLFdBQVc7QUFDakIsTUFBTSxhQUFhO0FBQ25CLE1BQU0saUJBQWlCO0FBQ3ZCLE1BQU0sV0FBVztBQUNqQixNQUFNLGFBQWE7QUFDbkIsTUFBTSxVQUFVO0FBQ2hCLE1BQU0sT0FBTztBQUNiLE1BQU0sUUFBUTtBQUNkLE1BQU0sc0JBQXNCO0FBQzVCLE1BQU1DLGdCQUFlO0FBQ3JCLE1BQU0sUUFBUTtBQUNkLE1BQU0sU0FBUztBQUNmLE1BQU0sU0FBUztBQUNmLE1BQU0sVUFBVTtBQUNoQixNQUFNLFlBQVk7QUFDbEIsTUFBTSxlQUFlO0FBQ3JCLE1BQU0sZUFBZTtBQUVyQixNQUFNLGFBQWEsSUFBSSxFQUFFO0FBRXpCLFdBQVMsYUFBYUMsUUFBTTtBQUN4QixXQUFPO0FBQUEsTUFDSCxLQUFLLENBQUMsZUFBZSxhQUFhLHVCQUF1QixZQUFZLFFBQVEsVUFBVSxDQUFDO0FBQUEsTUFDeEYsUUFBUSxNQUFNQSxPQUFLLE1BQU07QUFBQSxNQUN6QixVQUFVLENBQUMsVUFBVUEsT0FBSyxVQUFVLEVBQUMsTUFBSyxDQUFDO0FBQUEsTUFDM0MsWUFBWSxNQUFNQSxPQUFLLFVBQVU7QUFBQSxNQUNqQyxjQUFjLE1BQU1BLE9BQUssWUFBWTtBQUFBLE1BQ3JDLFNBQVMsQ0FBQ0MsUUFBT0MsWUFBV0YsT0FBSyxTQUFTLEVBQUMsT0FBQUMsUUFBTyxRQUFBQyxRQUFNLENBQUM7QUFBQSxNQUN6RCxNQUFNLE1BQU1GLE9BQUssSUFBSTtBQUFBLE1BQ3JCLFlBQVksQ0FBQ0MsUUFBT0MsWUFBV0YsT0FBSyxZQUFZLEVBQUMsT0FBQUMsUUFBTyxRQUFBQyxRQUFNLENBQUM7QUFBQSxNQUMvRCxZQUFZLENBQUNELFFBQU9DLFlBQVdGLE9BQUssWUFBWSxFQUFDLE9BQUFDLFFBQU8sUUFBQUMsUUFBTSxDQUFDO0FBQUEsTUFDL0QsZ0JBQWdCLENBQUMsVUFBVUYsT0FBSyxnQkFBZ0IsRUFBQyxhQUFhLE1BQUssQ0FBQztBQUFBLE1BQ3BFLHFCQUFxQixDQUFDLEdBQUcsTUFBTUEsT0FBSyxxQkFBcUIsRUFBQyxHQUFHLEVBQUMsQ0FBQztBQUFBLE1BQy9ELGtCQUFrQixNQUFNQSxPQUFLLGdCQUFnQjtBQUFBLE1BQzdDLFFBQVEsTUFBTUEsT0FBSyxNQUFNO0FBQUEsTUFDekIsTUFBTSxNQUFNQSxPQUFLLElBQUk7QUFBQSxNQUNyQixVQUFVLE1BQU1BLE9BQUssUUFBUTtBQUFBLE1BQzdCLFlBQVksTUFBTUEsT0FBSyxVQUFVO0FBQUEsTUFDakMsZ0JBQWdCLE1BQU1BLE9BQUssY0FBYztBQUFBLE1BQ3pDLFVBQVUsTUFBTUEsT0FBSyxRQUFRO0FBQUEsTUFDN0IsWUFBWSxNQUFNQSxPQUFLLFVBQVU7QUFBQSxNQUNqQyxTQUFTLE1BQU1BLE9BQUssT0FBTztBQUFBLE1BQzNCLE1BQU0sTUFBTUEsT0FBSyxJQUFJO0FBQUEsTUFDckIsT0FBTyxNQUFNQSxPQUFLLEtBQUs7QUFBQSxNQUN2QixxQkFBcUIsQ0FBQyxHQUFHLEdBQUcsR0FBRyxNQUFNQSxPQUFLLHFCQUFxQixFQUFDLEdBQUcsR0FBRyxHQUFHLEVBQUMsQ0FBQztBQUFBLE1BQzNFLGNBQWMsQ0FBQ0csZUFBY0gsT0FBS0QsZUFBYyxFQUFDLFdBQUFJLFdBQVMsQ0FBQztBQUFBLE1BQzNELE9BQU8sTUFBTUgsT0FBSyxLQUFLO0FBQUEsTUFDdkIsUUFBUSxNQUFNQSxPQUFLLE1BQU07QUFBQSxNQUN6QixRQUFRLE1BQU1BLE9BQUssTUFBTTtBQUFBLE1BQ3pCLFNBQVMsTUFBTUEsT0FBSyxPQUFPO0FBQUEsTUFDM0IsV0FBVyxNQUFNQSxPQUFLLFNBQVM7QUFBQSxNQUMvQixjQUFjLE1BQU1BLE9BQUssWUFBWTtBQUFBLE1BQ3JDLGNBQWMsQ0FBQyxjQUFjQSxPQUFLLGNBQWMsRUFBQyxVQUFTLENBQUM7QUFBQSxJQUMvRDtBQUFBLEVBQ0o7QUFRTyxXQUFTLElBQUksWUFBWTtBQUM1QixXQUFPLGFBQWEsdUJBQXVCLFlBQVksUUFBUSxVQUFVLENBQUM7QUFBQSxFQUM5RTtBQUtPLFdBQVMsU0FBUztBQUNyQixlQUFXLE9BQU87QUFBQSxFQUN0QjtBQU1PLFdBQVMsU0FBUyxPQUFPO0FBQzVCLGVBQVcsU0FBUyxLQUFLO0FBQUEsRUFDN0I7QUFLTyxXQUFTLGFBQWE7QUFDekIsZUFBVyxXQUFXO0FBQUEsRUFDMUI7QUFPTyxXQUFTLFFBQVFDLFFBQU9DLFNBQVE7QUFDbkMsZUFBVyxRQUFRRCxRQUFPQyxPQUFNO0FBQUEsRUFDcEM7QUFLTyxXQUFTLE9BQU87QUFDbkIsV0FBTyxXQUFXLEtBQUs7QUFBQSxFQUMzQjtBQU9PLFdBQVMsV0FBV0QsUUFBT0MsU0FBUTtBQUN0QyxlQUFXLFdBQVdELFFBQU9DLE9BQU07QUFBQSxFQUN2QztBQU9PLFdBQVMsV0FBV0QsUUFBT0MsU0FBUTtBQUN0QyxlQUFXLFdBQVdELFFBQU9DLE9BQU07QUFBQSxFQUN2QztBQU1PLFdBQVMsZUFBZSxPQUFPO0FBQ2xDLGVBQVcsZUFBZSxLQUFLO0FBQUEsRUFDbkM7QUFPTyxXQUFTLG9CQUFvQixHQUFHLEdBQUc7QUFDdEMsZUFBVyxvQkFBb0IsR0FBRyxDQUFDO0FBQUEsRUFDdkM7QUFLTyxXQUFTLG1CQUFtQjtBQUMvQixXQUFPLFdBQVcsaUJBQWlCO0FBQUEsRUFDdkM7QUFLTyxXQUFTLFNBQVM7QUFDckIsV0FBTyxXQUFXLE9BQU87QUFBQSxFQUM3QjtBQUtPLFdBQVNFLFFBQU87QUFDbkIsZUFBVyxLQUFLO0FBQUEsRUFDcEI7QUFLTyxXQUFTLFdBQVc7QUFDdkIsZUFBVyxTQUFTO0FBQUEsRUFDeEI7QUFLTyxXQUFTLGFBQWE7QUFDekIsZUFBVyxXQUFXO0FBQUEsRUFDMUI7QUFLTyxXQUFTLGlCQUFpQjtBQUM3QixlQUFXLGVBQWU7QUFBQSxFQUM5QjtBQUtPLFdBQVMsV0FBVztBQUN2QixlQUFXLFNBQVM7QUFBQSxFQUN4QjtBQUtPLFdBQVMsYUFBYTtBQUN6QixlQUFXLFdBQVc7QUFBQSxFQUMxQjtBQUtPLFdBQVMsVUFBVTtBQUN0QixlQUFXLFFBQVE7QUFBQSxFQUN2QjtBQUtPLFdBQVNDLFFBQU87QUFDbkIsZUFBVyxLQUFLO0FBQUEsRUFDcEI7QUFLTyxXQUFTLFFBQVE7QUFDcEIsZUFBVyxNQUFNO0FBQUEsRUFDckI7QUFTTyxXQUFTLG9CQUFvQixHQUFHLEdBQUcsR0FBRyxHQUFHO0FBQzVDLGVBQVcsb0JBQW9CLEdBQUcsR0FBRyxHQUFHLENBQUM7QUFBQSxFQUM3QztBQU1PLFdBQVMsYUFBYUYsWUFBVztBQUNwQyxlQUFXLGFBQWFBLFVBQVM7QUFBQSxFQUNyQztBQUtPLFdBQVMsUUFBUTtBQUNwQixXQUFPLFdBQVcsTUFBTTtBQUFBLEVBQzVCO0FBS08sV0FBUyxTQUFTO0FBQ3JCLFdBQU8sV0FBVyxPQUFPO0FBQUEsRUFDN0I7QUFLTyxXQUFTLFNBQVM7QUFDckIsZUFBVyxPQUFPO0FBQUEsRUFDdEI7QUFLTyxXQUFTLFVBQVU7QUFDdEIsZUFBVyxRQUFRO0FBQUEsRUFDdkI7QUFLTyxXQUFTLFlBQVk7QUFDeEIsZUFBVyxVQUFVO0FBQUEsRUFDekI7QUFLTyxXQUFTLGVBQWU7QUFDM0IsV0FBTyxXQUFXLGFBQWE7QUFBQSxFQUNuQztBQU1PLFdBQVMsYUFBYSxXQUFXO0FBQ3BDLGVBQVcsYUFBYSxTQUFTO0FBQUEsRUFDckM7OztBQzNUQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNDTyxNQUFNLGFBQWE7QUFBQSxJQUN6QixTQUFTO0FBQUEsTUFDUixvQkFBb0I7QUFBQSxNQUNwQixzQkFBc0I7QUFBQSxNQUN0QixZQUFZO0FBQUEsTUFDWixvQkFBb0I7QUFBQSxNQUNwQixrQkFBa0I7QUFBQSxNQUNsQix1QkFBdUI7QUFBQSxNQUN2QixvQkFBb0I7QUFBQSxNQUNwQiw0QkFBNEI7QUFBQSxNQUM1QixnQkFBZ0I7QUFBQSxNQUNoQixjQUFjO0FBQUEsTUFDZCxtQkFBbUI7QUFBQSxNQUNuQixnQkFBZ0I7QUFBQSxNQUNoQixrQkFBa0I7QUFBQSxNQUNsQixrQkFBa0I7QUFBQSxNQUNsQixvQkFBb0I7QUFBQSxNQUNwQixlQUFlO0FBQUEsTUFDZixnQkFBZ0I7QUFBQSxNQUNoQixrQkFBa0I7QUFBQSxNQUNsQixhQUFhO0FBQUEsTUFDYixnQkFBZ0I7QUFBQSxNQUNoQixpQkFBaUI7QUFBQSxNQUNqQixnQkFBZ0I7QUFBQSxNQUNoQixpQkFBaUI7QUFBQSxNQUNqQixpQkFBaUI7QUFBQSxNQUNqQixnQkFBZ0I7QUFBQSxJQUNqQjtBQUFBLElBQ0EsS0FBSztBQUFBLE1BQ0osNEJBQTRCO0FBQUEsTUFDNUIsdUNBQXVDO0FBQUEsTUFDdkMseUNBQXlDO0FBQUEsTUFDekMsMEJBQTBCO0FBQUEsTUFDMUIsb0NBQW9DO0FBQUEsTUFDcEMsc0NBQXNDO0FBQUEsTUFDdEMsb0NBQW9DO0FBQUEsTUFDcEMsMENBQTBDO0FBQUEsTUFDMUMsK0JBQStCO0FBQUEsTUFDL0Isb0JBQW9CO0FBQUEsTUFDcEIsd0NBQXdDO0FBQUEsTUFDeEMsc0JBQXNCO0FBQUEsTUFDdEIsc0JBQXNCO0FBQUEsTUFDdEIsNkJBQTZCO0FBQUEsTUFDN0IsZ0NBQWdDO0FBQUEsTUFDaEMscUJBQXFCO0FBQUEsTUFDckIsNkJBQTZCO0FBQUEsTUFDN0IsMEJBQTBCO0FBQUEsTUFDMUIsdUJBQXVCO0FBQUEsTUFDdkIsdUJBQXVCO0FBQUEsTUFDdkIsMkJBQTJCO0FBQUEsTUFDM0IsK0JBQStCO0FBQUEsTUFDL0Isb0JBQW9CO0FBQUEsTUFDcEIscUJBQXFCO0FBQUEsTUFDckIscUJBQXFCO0FBQUEsTUFDckIsc0JBQXNCO0FBQUEsTUFDdEIsZ0NBQWdDO0FBQUEsTUFDaEMsa0NBQWtDO0FBQUEsTUFDbEMsbUNBQW1DO0FBQUEsTUFDbkMsb0NBQW9DO0FBQUEsTUFDcEMsK0JBQStCO0FBQUEsTUFDL0IsNkJBQTZCO0FBQUEsTUFDN0IsdUJBQXVCO0FBQUEsTUFDdkIsaUNBQWlDO0FBQUEsTUFDakMsOEJBQThCO0FBQUEsTUFDOUIsNEJBQTRCO0FBQUEsTUFDNUIsc0NBQXNDO0FBQUEsTUFDdEMsNEJBQTRCO0FBQUEsTUFDNUIsc0JBQXNCO0FBQUEsTUFDdEIsa0NBQWtDO0FBQUEsTUFDbEMsc0JBQXNCO0FBQUEsTUFDdEIsd0JBQXdCO0FBQUEsTUFDeEIsMkJBQTJCO0FBQUEsTUFDM0Isd0JBQXdCO0FBQUEsTUFDeEIsbUJBQW1CO0FBQUEsTUFDbkIsMEJBQTBCO0FBQUEsTUFDMUIsOEJBQThCO0FBQUEsTUFDOUIseUJBQXlCO0FBQUEsTUFDekIsNkJBQTZCO0FBQUEsTUFDN0IsaUJBQWlCO0FBQUEsTUFDakIsZ0JBQWdCO0FBQUEsTUFDaEIsc0JBQXNCO0FBQUEsTUFDdEIsZUFBZTtBQUFBLE1BQ2YseUJBQXlCO0FBQUEsTUFDekIsd0JBQXdCO0FBQUEsTUFDeEIsb0JBQW9CO0FBQUEsTUFDcEIscUJBQXFCO0FBQUEsTUFDckIsaUJBQWlCO0FBQUEsTUFDakIsaUJBQWlCO0FBQUEsTUFDakIsc0JBQXNCO0FBQUEsTUFDdEIsbUNBQW1DO0FBQUEsTUFDbkMscUNBQXFDO0FBQUEsTUFDckMsdUJBQXVCO0FBQUEsTUFDdkIsc0JBQXNCO0FBQUEsTUFDdEIsd0JBQXdCO0FBQUEsTUFDeEIsMkJBQTJCO0FBQUEsTUFDM0IsbUJBQW1CO0FBQUEsTUFDbkIscUJBQXFCO0FBQUEsTUFDckIsc0JBQXNCO0FBQUEsTUFDdEIsc0JBQXNCO0FBQUEsTUFDdEIsOEJBQThCO0FBQUEsTUFDOUIsaUJBQWlCO0FBQUEsTUFDakIseUJBQXlCO0FBQUEsTUFDekIsMkJBQTJCO0FBQUEsTUFDM0IsK0JBQStCO0FBQUEsTUFDL0IsMEJBQTBCO0FBQUEsTUFDMUIsOEJBQThCO0FBQUEsTUFDOUIsaUJBQWlCO0FBQUEsTUFDakIsdUJBQXVCO0FBQUEsTUFDdkIsZ0JBQWdCO0FBQUEsTUFDaEIsMEJBQTBCO0FBQUEsTUFDMUIseUJBQXlCO0FBQUEsTUFDekIsc0JBQXNCO0FBQUEsTUFDdEIsa0JBQWtCO0FBQUEsTUFDbEIsbUJBQW1CO0FBQUEsTUFDbkIsa0JBQWtCO0FBQUEsTUFDbEIsdUJBQXVCO0FBQUEsTUFDdkIsb0NBQW9DO0FBQUEsTUFDcEMsc0NBQXNDO0FBQUEsTUFDdEMsd0JBQXdCO0FBQUEsTUFDeEIsdUJBQXVCO0FBQUEsTUFDdkIseUJBQXlCO0FBQUEsTUFDekIsNEJBQTRCO0FBQUEsTUFDNUIsNEJBQTRC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QUFTLFlBQVksS0FBSyxNQUFNLEdBQUcsRUFBRSxXQUFXLEdBQUc7QUFDMUQsWUFBTSxJQUFJLE1BQU0sb0VBQW9FO0FBQUEsSUFDeEY7QUFDQSxRQUFJLENBQUMsYUFBYSxZQUFZLFVBQVUsSUFBSSxLQUFLLE1BQU0sR0FBRztBQUMxRCxXQUFPLFlBQVksYUFBYTtBQUFBLE1BQzVCO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsSUFDSixDQUFDO0FBQUEsRUFDTDtBQVNPLFdBQVMsS0FBSyxhQUFhLE1BQU07QUFDcEMsV0FBTyxZQUFZLGFBQWE7QUFBQSxNQUM1QjtBQUFBLE1BQ0E7QUFBQSxJQUNKLENBQUM7QUFBQSxFQUNMO0FBVU8sV0FBUyxPQUFPLFlBQVksZUFBZSxNQUFNO0FBQ3BELFdBQU8sWUFBWSxhQUFhO0FBQUEsTUFDNUIsYUFBYTtBQUFBLE1BQ2IsWUFBWTtBQUFBLE1BQ1o7QUFBQSxNQUNBO0FBQUEsSUFDSixDQUFDO0FBQUEsRUFDTDs7O0FDcEpBLFNBQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQWdEbEMsU0FBTyxPQUFPLFNBQU87QUFFckIsU0FBTyxRQUFRLE9BQU8sU0FBUyxDQUFDO0FBQ2hDLFNBQU8sTUFBTSxjQUFjO0FBQzNCLFNBQU8sTUFBTSxVQUFVO0FBQ3ZCLFNBQU8sTUFBTSxPQUFPO0FBQ3BCLFNBQU8sTUFBTSxZQUFZO0FBQ3pCLFNBQU8sTUFBTSxVQUFVO0FBQ3ZCLFNBQU8sTUFBTSxTQUFTO0FBQ3RCLFNBQU8sTUFBTSxRQUFRO0FBQ3JCLFNBQU8sTUFBTSxVQUFVO0FBQ3ZCLFNBQU8sTUFBTSxTQUFTO0FBQ3RCLFNBQU8sTUFBTSxTQUFTO0FBQ3RCLFNBQU8sTUFBTSxNQUFNO0FBR25CLE1BQUksVUFBVTtBQUNkLFdBQVMsaUJBQWlCLG9CQUFvQixXQUFXO0FBQ3JELGNBQVU7QUFDVixXQUFPLE9BQU8sT0FBTyxxQkFBcUI7QUFDMUMsUUFBRyxNQUFPO0FBQ04sZUFBUyxzQkFBc0I7QUFBQSxJQUNuQztBQUFBLEVBQ0osQ0FBQztBQUVELFdBQVMsVUFBVSxJQUFJO0FBQ25CLFFBQUksV0FBVyxTQUFTLGVBQWUsWUFBWTtBQUMvQyxTQUFHO0FBQUEsSUFDUCxPQUFPO0FBQ0gsZUFBUyxpQkFBaUIsb0JBQW9CLEVBQUU7QUFBQSxJQUNwRDtBQUFBLEVBQ0o7QUFFQSxZQUFVLE1BQU07QUFDWixJQUFJLE9BQU87QUFBQSxFQUNmLENBQUM7IiwKICAibmFtZXMiOiBbInNpemUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgIkhpZGUiLCAiU2hvdyIsICJzZXRSZXNpemFibGUiLCAiY2FsbCIsICJ3aWR0aCIsICJoZWlnaHQiLCAicmVzaXphYmxlIiwgIkhpZGUiLCAiU2hvdyIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJFcnJvciIsICJjYWxsIiwgIkVycm9yIiwgImNhbGwiLCAiZ2VuZXJhdGVJRCJdCn0K
