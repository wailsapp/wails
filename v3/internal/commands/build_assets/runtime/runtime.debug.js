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

  // desktop/@wailsio/runtime/src/screens.js
  var screens_exports = {};
  __export(screens_exports, {
    GetAll: () => GetAll,
    GetCurrent: () => GetCurrent,
    GetPrimary: () => GetPrimary
  });
  var call4 = newRuntimeCallerWithID(objectNames.Screens, "");
  var getAll = 0;
  var getPrimary = 1;
  var getCurrent = 2;
  function GetAll() {
    return call4(getAll);
  }
  function GetPrimary() {
    return call4(getPrimary);
  }
  function GetCurrent() {
    return call4(getCurrent);
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
  var call5 = newRuntimeCallerWithID(objectNames.System, "");
  var systemIsDarkMode = 0;
  var environment = 1;
  function invoke(msg) {
    if (window.chrome) {
      return window.chrome.webview.postMessage(msg);
    }
    return window.webkit.messageHandlers.external.postMessage(msg);
  }
  function IsDarkMode() {
    return call5(systemIsDarkMode);
  }
  function Capabilities() {
    let response = fetch("/wails/capabilities");
    return response.json();
  }
  function Environment() {
    return call5(environment);
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
  var setResizable = 22;
  var width = 23;
  var height = 24;
  var zoomIn = 25;
  var zoomOut = 26;
  var zoomReset = 27;
  var getZoomLevel = 28;
  var setZoomLevel = 29;
  var thisWindow = Get("");
  function createWindow(call9) {
    return {
      Get: (windowName) => createWindow(newRuntimeCallerWithID(objectNames.Window, windowName)),
      Center: () => call9(center),
      SetTitle: (title) => call9(setTitle, { title }),
      Fullscreen: () => call9(fullscreen),
      UnFullscreen: () => call9(unFullscreen),
      SetSize: (width2, height2) => call9(setSize, { width: width2, height: height2 }),
      Size: () => call9(size),
      SetMaxSize: (width2, height2) => call9(setMaxSize, { width: width2, height: height2 }),
      SetMinSize: (width2, height2) => call9(setMinSize, { width: width2, height: height2 }),
      SetAlwaysOnTop: (onTop) => call9(setAlwaysOnTop, { alwaysOnTop: onTop }),
      SetRelativePosition: (x, y) => call9(setRelativePosition, { x, y }),
      RelativePosition: () => call9(relativePosition),
      Screen: () => call9(screen),
      Hide: () => call9(hide),
      Maximise: () => call9(maximise),
      UnMaximise: () => call9(unMaximise),
      ToggleMaximise: () => call9(toggleMaximise),
      Minimise: () => call9(minimise),
      UnMinimise: () => call9(unMinimise),
      Restore: () => call9(restore),
      Show: () => call9(show),
      Close: () => call9(close),
      SetBackgroundColour: (r, g, b, a) => call9(setBackgroundColour, { r, g, b, a }),
      SetResizable: (resizable) => call9(setResizable, { resizable }),
      Width: () => call9(width),
      Height: () => call9(height),
      ZoomIn: () => call9(zoomIn),
      ZoomOut: () => call9(zoomOut),
      ZoomReset: () => call9(zoomReset),
      GetZoomLevel: () => call9(getZoomLevel),
      SetZoomLevel: (zoomLevel) => call9(setZoomLevel, { zoomLevel })
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
  function SetResizable(resizable) {
    thisWindow.SetResizable(resizable);
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
  var call6 = newRuntimeCallerWithID(objectNames.Events, "");
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
    return call6(EmitMethod, event);
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
  var call7 = newRuntimeCallerWithID(objectNames.Dialog, "");
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
      call7(type, options).catch((error) => {
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
    if (true) {
      debugLog("Reloading WML");
    }
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
  var call8 = newRuntimeCallerWithID(objectNames.Call, "");
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
      call8(type, options).catch((error) => {
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb24uanMiLCAiLi4vLi4vLi4vcnVudGltZS9ub2RlX21vZHVsZXMvbmFub2lkL25vbi1zZWN1cmUvaW5kZXguanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZC5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvd2luZG93LmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50cy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RpYWxvZ3MuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxzLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9jb21waWxlZC9tYWluLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyIvKipcclxuICogTG9ncyBhIG1lc3NhZ2UgdG8gdGhlIGNvbnNvbGUgd2l0aCBjdXN0b20gZm9ybWF0dGluZy5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZSkge1xyXG4gICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXHJcbiAgICBjb25zb2xlLmxvZyhcclxuICAgICAgICAnJWMgd2FpbHMzICVjICcgKyBtZXNzYWdlICsgJyAnLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjYWEwMDAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogM3B4IDBweCAwcHggM3B4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogIzAwOTkwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDBweCAzcHggM3B4IDBweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbSdcclxuICAgICk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uLCAnJyk7XHJcblxyXG5jb25zdCBIaWRlTWV0aG9kID0gMDtcclxuY29uc3QgU2hvd01ldGhvZCA9IDE7XHJcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xyXG5cclxuLyoqXHJcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cclxuICpcclxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICpcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyB0aGUgU2hvd01ldGhvZCBhbmQgcmV0dXJucyB0aGUgcmVzdWx0LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XHJcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBRdWl0TWV0aG9kIHRvIHRlcm1pbmF0ZSB0aGUgcHJvZ3JhbS5cclxuICpcclxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XHJcbn1cclxuIiwgImxldCB1cmxBbHBoYWJldCA9XG4gICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuZXhwb3J0IGxldCBjdXN0b21BbHBoYWJldCA9IChhbHBoYWJldCwgZGVmYXVsdFNpemUgPSAyMSkgPT4ge1xuICByZXR1cm4gKHNpemUgPSBkZWZhdWx0U2l6ZSkgPT4ge1xuICAgIGxldCBpZCA9ICcnXG4gICAgbGV0IGkgPSBzaXplXG4gICAgd2hpbGUgKGktLSkge1xuICAgICAgaWQgKz0gYWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiBhbHBoYWJldC5sZW5ndGgpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG4gIH1cbn1cbmV4cG9ydCBsZXQgbmFub2lkID0gKHNpemUgPSAyMSkgPT4ge1xuICBsZXQgaWQgPSAnJ1xuICBsZXQgaSA9IHNpemVcbiAgd2hpbGUgKGktLSkge1xuICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgfVxuICByZXR1cm4gaWRcbn1cbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xyXG5cclxuLy8gT2JqZWN0IE5hbWVzXHJcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcclxuICAgIENhbGw6IDAsXHJcbiAgICBDbGlwYm9hcmQ6IDEsXHJcbiAgICBBcHBsaWNhdGlvbjogMixcclxuICAgIEV2ZW50czogMyxcclxuICAgIENvbnRleHRNZW51OiA0LFxyXG4gICAgRGlhbG9nOiA1LFxyXG4gICAgV2luZG93OiA2LFxyXG4gICAgU2NyZWVuczogNyxcclxuICAgIFN5c3RlbTogOCxcclxuICAgIEJyb3dzZXI6IDksXHJcbn1cclxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xyXG5cclxuLyoqXHJcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cclxuICpcclxuICogQHBhcmFtIHtPYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3Qgb24gd2hpY2ggdGhlIG1ldGhvZCBpcyB0byBiZSBpbnZva2VkLlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXHJcbiAqIEByZXR1cm5zIHtGdW5jdGlvbn0gQSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IHRha2VzIHRoZSBtZXRob2QgbmFtZSBhbmQgb3B0aW9uYWxseSBhcmd1bWVudHMgYW5kIGludm9rZXMgdGhlIG1ldGhvZCB3aXRoaW4gdGhlIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xyXG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogQ3JlYXRlcyBhIG5ldyBydW50aW1lIGNhbGxlciB3aXRoIHNwZWNpZmllZCBJRC5cclxuICpcclxuICogQHBhcmFtIHtvYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICogQHJldHVybiB7RnVuY3Rpb259IC0gVGhlIG5ldyBydW50aW1lIGNhbGxlciBmdW5jdGlvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdCwgd2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xyXG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9O1xyXG59XHJcblxyXG5cclxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcclxuICAgIGlmKCBtZXRob2QgKSB7XHJcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIH1cclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG5cclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgZmV0Y2godXJsLCBmZXRjaE9wdGlvbnMpXHJcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcclxuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xyXG4gICAgICAgICAgICAgICAgICAgIC8vIGNoZWNrIGNvbnRlbnQgdHlwZVxyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcclxuICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XHJcbiAgICAgICAgICAgIH0pXHJcbiAgICAgICAgICAgIC50aGVuKGRhdGEgPT4gcmVzb2x2ZShkYXRhKSlcclxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5Ccm93c2VyLCAnJyk7XHJcbmNvbnN0IEJyb3dzZXJPcGVuVVJMID0gMDtcclxuXHJcbi8qKlxyXG4gKiBPcGVuIGEgYnJvd3NlciB3aW5kb3cgdG8gdGhlIGdpdmVuIFVSTFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIC0gVGhlIFVSTCB0byBvcGVuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmwpIHtcclxuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsfSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCwgJycpO1xyXG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcclxuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBrZXlTdHJpbmcgLSBUaGUga2V5IHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBmb3IuXHJcbiAqIEByZXR1cm4geyp9IC0gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleVN0cmluZykge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gd2luZG93Ll93YWlscy5mbGFnc1trZXlTdHJpbmddO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleVN0cmluZyArIFwiJzogXCIgKyBlKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFBvc2l0aW9uXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZS5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgWSBjb29yZGluYXRlLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTaXplXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIHdpZHRoLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gWSAtIFRoZSBoZWlnaHQuXHJcbiAqL1xyXG5cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBSZWN0XHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgdG9wLWxlZnQgY29ybmVyLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gWSAtIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHRvcC1sZWZ0IGNvcm5lci5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBIZWlnaHQgLSBUaGUgaGVpZ2h0IG9mIHRoZSByZWN0YW5nbGUuXHJcbiAqL1xyXG5cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7KCdaZXJvJ3wnTmluZXR5J3wnT25lRWlnaHR5J3wnVHdvU2V2ZW50eScpfSBSb3RhdGlvblxyXG4gKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gQ2FuIGJlIG9uZSBvZiAnWmVybycsICdOaW5ldHknLCAnT25lRWlnaHR5JywgJ1R3b1NldmVudHknLlxyXG4gKi9cclxuXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gU2NyZWVuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gTmFtZSAtIEh1bWFuIHJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFNjYWxlIC0gVGhlIHJlc29sdXRpb24gc2NhbGUgb2YgdGhlIHNjcmVlbi4gMSA9IHN0YW5kYXJkIHJlc29sdXRpb24sIDIgPSBoaWdoIChSZXRpbmEpLCBldGMuXHJcbiAqIEBwcm9wZXJ0eSB7UG9zaXRpb259IFBvc2l0aW9uIC0gQ29udGFpbnMgdGhlIFggYW5kIFkgY29vcmRpbmF0ZXMgb2YgdGhlIHNjcmVlbidzIHBvc2l0aW9uLlxyXG4gKiBAcHJvcGVydHkge1NpemV9IFNpemUgLSBDb250YWlucyB0aGUgd2lkdGggYW5kIGhlaWdodCBvZiB0aGUgc2NyZWVuLlxyXG4gKiBAcHJvcGVydHkge1JlY3R9IEJvdW5kcyAtIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC5cclxuICogQHByb3BlcnR5IHtSZWN0fSBXb3JrQXJlYSAtIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IElzUHJpbWFyeSAtIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLlxyXG4gKiBAcHJvcGVydHkge1JvdGF0aW9ufSBSb3RhdGlvbiAtIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLlxyXG4gKi9cclxuXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLlNjcmVlbnMsICcnKTtcclxuXHJcbmNvbnN0IGdldEFsbCA9IDA7XHJcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xyXG5jb25zdCBnZXRDdXJyZW50ID0gMjtcclxuXHJcbi8qKlxyXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW5bXT59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIGFycmF5IG9mIFNjcmVlbiBvYmplY3RzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcclxuICAgIHJldHVybiBjYWxsKGdldEFsbCk7XHJcbn1cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpIHtcclxuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xyXG59XHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0Q3VycmVudCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU3lzdGVtLCAnJyk7XHJcbmNvbnN0IHN5c3RlbUlzRGFya01vZGUgPSAwO1xyXG5jb25zdCBlbnZpcm9ubWVudCA9IDE7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZykge1xyXG4gICAgaWYod2luZG93LmNocm9tZSkge1xyXG4gICAgICAgIHJldHVybiB3aW5kb3cuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UobXNnKTtcclxuICAgIH1cclxuICAgIHJldHVybiB3aW5kb3cud2Via2l0Lm1lc3NhZ2VIYW5kbGVycy5leHRlcm5hbC5wb3N0TWVzc2FnZShtc2cpO1xyXG59XHJcblxyXG4vKipcclxuICogQGZ1bmN0aW9uXHJcbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPGJvb2xlYW4+fSAtIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGEgYm9vbGVhbiB2YWx1ZSBpbmRpY2F0aW5nIGlmIHRoZSBzeXN0ZW0gaXMgaW4gZGFyayBtb2RlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKSB7XHJcbiAgICByZXR1cm4gY2FsbChzeXN0ZW1Jc0RhcmtNb2RlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxyXG4gKlxyXG4gKiBAYXN5bmNcclxuICogQGZ1bmN0aW9uIENhcGFiaWxpdGllc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxPYmplY3Q+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyB0aGUgY2FwYWJpbGl0aWVzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhcGFiaWxpdGllcygpIHtcclxuICAgIGxldCByZXNwb25zZSA9IGZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKTtcclxuICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7b2JqZWN0fSBFbnZpcm9ubWVudEluZm9cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IE9TIC0gVGhlIG9wZXJhdGluZyBzeXN0ZW0gaW4gdXNlLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gQXJjaCAtIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS5cclxuICovXHJcblxyXG4vKipcclxuICogQGZ1bmN0aW9uXHJcbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxFbnZpcm9ubWVudEluZm8+fSAtIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIE9TIGFuZCBzeXN0ZW0gYXJjaGl0ZWN0dXJlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoZW52aXJvbm1lbnQpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cclxuICpcclxuICogQHJldHVybiB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwid2luZG93c1wiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSBSZXR1cm5zIHRydWUgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJsaW51eFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIG1hY09TLCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNNYWMoKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJkYXJ3aW5cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFtZDY0XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTS5cclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm1cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUuXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSAtIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm02NFwiO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkRlYnVnID09PSB0cnVlO1xyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLy8gSW1wb3J0IHNjcmVlbiBqc2RvYyBkZWZpbml0aW9uIGZyb20gLi9zY3JlZW5zLmpzXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9zY3JlZW5zXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuY29uc3QgY2VudGVyID0gMDtcclxuY29uc3Qgc2V0VGl0bGUgPSAxO1xyXG5jb25zdCBmdWxsc2NyZWVuID0gMjtcclxuY29uc3QgdW5GdWxsc2NyZWVuID0gMztcclxuY29uc3Qgc2V0U2l6ZSA9IDQ7XHJcbmNvbnN0IHNpemUgPSA1O1xyXG5jb25zdCBzZXRNYXhTaXplID0gNjtcclxuY29uc3Qgc2V0TWluU2l6ZSA9IDc7XHJcbmNvbnN0IHNldEFsd2F5c09uVG9wID0gODtcclxuY29uc3Qgc2V0UmVsYXRpdmVQb3NpdGlvbiA9IDk7XHJcbmNvbnN0IHJlbGF0aXZlUG9zaXRpb24gPSAxMDtcclxuY29uc3Qgc2NyZWVuID0gMTE7XHJcbmNvbnN0IGhpZGUgPSAxMjtcclxuY29uc3QgbWF4aW1pc2UgPSAxMztcclxuY29uc3QgdW5NYXhpbWlzZSA9IDE0O1xyXG5jb25zdCB0b2dnbGVNYXhpbWlzZSA9IDE1O1xyXG5jb25zdCBtaW5pbWlzZSA9IDE2O1xyXG5jb25zdCB1bk1pbmltaXNlID0gMTc7XHJcbmNvbnN0IHJlc3RvcmUgPSAxODtcclxuY29uc3Qgc2hvdyA9IDE5O1xyXG5jb25zdCBjbG9zZSA9IDIwO1xyXG5jb25zdCBzZXRCYWNrZ3JvdW5kQ29sb3VyID0gMjE7XHJcbmNvbnN0IHNldFJlc2l6YWJsZSA9IDIyO1xyXG5jb25zdCB3aWR0aCA9IDIzO1xyXG5jb25zdCBoZWlnaHQgPSAyNDtcclxuY29uc3Qgem9vbUluID0gMjU7XHJcbmNvbnN0IHpvb21PdXQgPSAyNjtcclxuY29uc3Qgem9vbVJlc2V0ID0gMjc7XHJcbmNvbnN0IGdldFpvb21MZXZlbCA9IDI4O1xyXG5jb25zdCBzZXRab29tTGV2ZWwgPSAyOTtcclxuXHJcbmNvbnN0IHRoaXNXaW5kb3cgPSBHZXQoJycpO1xyXG5cclxuZnVuY3Rpb24gY3JlYXRlV2luZG93KGNhbGwpIHtcclxuICAgIHJldHVybiB7XHJcbiAgICAgICAgR2V0OiAod2luZG93TmFtZSkgPT4gY3JlYXRlV2luZG93KG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuV2luZG93LCB3aW5kb3dOYW1lKSksXHJcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKGNlbnRlciksXHJcbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbChzZXRUaXRsZSwge3RpdGxlfSksXHJcbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gY2FsbChmdWxsc2NyZWVuKSxcclxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwodW5GdWxsc2NyZWVuKSxcclxuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbChzZXRTaXplLCB7d2lkdGgsIGhlaWdodH0pLFxyXG4gICAgICAgIFNpemU6ICgpID0+IGNhbGwoc2l6ZSksXHJcbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWF4U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcclxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbChzZXRNaW5TaXplLCB7d2lkdGgsIGhlaWdodH0pLFxyXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAob25Ub3ApID0+IGNhbGwoc2V0QWx3YXlzT25Ub3AsIHthbHdheXNPblRvcDogb25Ub3B9KSxcclxuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbChzZXRSZWxhdGl2ZVBvc2l0aW9uLCB7eCwgeX0pLFxyXG4gICAgICAgIFJlbGF0aXZlUG9zaXRpb246ICgpID0+IGNhbGwocmVsYXRpdmVQb3NpdGlvbiksXHJcbiAgICAgICAgU2NyZWVuOiAoKSA9PiBjYWxsKHNjcmVlbiksXHJcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbChoaWRlKSxcclxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gY2FsbChtYXhpbWlzZSksXHJcbiAgICAgICAgVW5NYXhpbWlzZTogKCkgPT4gY2FsbCh1bk1heGltaXNlKSxcclxuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gY2FsbCh0b2dnbGVNYXhpbWlzZSksXHJcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IGNhbGwobWluaW1pc2UpLFxyXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwodW5NaW5pbWlzZSksXHJcbiAgICAgICAgUmVzdG9yZTogKCkgPT4gY2FsbChyZXN0b3JlKSxcclxuICAgICAgICBTaG93OiAoKSA9PiBjYWxsKHNob3cpLFxyXG4gICAgICAgIENsb3NlOiAoKSA9PiBjYWxsKGNsb3NlKSxcclxuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gY2FsbChzZXRCYWNrZ3JvdW5kQ29sb3VyLCB7ciwgZywgYiwgYX0pLFxyXG4gICAgICAgIFNldFJlc2l6YWJsZTogKHJlc2l6YWJsZSkgPT4gY2FsbChzZXRSZXNpemFibGUsIHtyZXNpemFibGV9KSxcclxuICAgICAgICBXaWR0aDogKCkgPT4gY2FsbCh3aWR0aCksXHJcbiAgICAgICAgSGVpZ2h0OiAoKSA9PiBjYWxsKGhlaWdodCksXHJcbiAgICAgICAgWm9vbUluOiAoKSA9PiBjYWxsKHpvb21JbiksXHJcbiAgICAgICAgWm9vbU91dDogKCkgPT4gY2FsbCh6b29tT3V0KSxcclxuICAgICAgICBab29tUmVzZXQ6ICgpID0+IGNhbGwoem9vbVJlc2V0KSxcclxuICAgICAgICBHZXRab29tTGV2ZWw6ICgpID0+IGNhbGwoZ2V0Wm9vbUxldmVsKSxcclxuICAgICAgICBTZXRab29tTGV2ZWw6ICh6b29tTGV2ZWwpID0+IGNhbGwoc2V0Wm9vbUxldmVsLCB7em9vbUxldmVsfSksXHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgc3BlY2lmaWVkIHdpbmRvdy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cclxuICogQHJldHVybiB7T2JqZWN0fSAtIFRoZSBzcGVjaWZpZWQgd2luZG93IG9iamVjdC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXQod2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGNyZWF0ZVdpbmRvdyhuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgd2luZG93TmFtZSkpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2VudGVycyB0aGUgd2luZG93IG9uIHRoZSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2VudGVyKCkge1xyXG4gICAgdGhpc1dpbmRvdy5DZW50ZXIoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZSAtIFRoZSB0aXRsZSB0byBzZXQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0VGl0bGUodGl0bGUpIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0VGl0bGUodGl0bGUpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgd2luZG93IHRvIGZ1bGxzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRnVsbHNjcmVlbigpIHtcclxuICAgIHRoaXNXaW5kb3cuRnVsbHNjcmVlbigpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGggLSBUaGUgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCAtIFRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRTaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaXplKCkge1xyXG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuU2l6ZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCAtIFRoZSBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRNYXhTaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGggLSBUaGUgbWluaW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB0aGlzV2luZG93LlNldE1pblNpemUod2lkdGgsIGhlaWdodCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSB3aW5kb3cgdG8gYWx3YXlzIGJlIG9uIHRvcC5cclxuICogQHBhcmFtIHtib29sZWFufSBvblRvcCAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYWx3YXlzIGJlIG9uIHRvcC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRBbHdheXNPblRvcChvblRvcCkge1xyXG4gICAgdGhpc1dpbmRvdy5TZXRBbHdheXNPblRvcChvblRvcCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIHdpbmRvdydzIHBvc2l0aW9uLlxyXG4gKiBAcGFyYW0ge251bWJlcn0geSAtIFRoZSB5LWNvb3JkaW5hdGUgb2YgdGhlIHdpbmRvdydzIHBvc2l0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFJlbGF0aXZlUG9zaXRpb24oeCwgeSkge1xyXG4gICAgdGhpc1dpbmRvdy5TZXRSZWxhdGl2ZVBvc2l0aW9uKHgsIHkpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBSZWxhdGl2ZVBvc2l0aW9uKCkge1xyXG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuUmVsYXRpdmVQb3NpdGlvbigpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgc2NyZWVuIHRoYXQgdGhlIHdpbmRvdyBpcyBvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTY3JlZW4oKSB7XHJcbiAgICByZXR1cm4gdGhpc1dpbmRvdy5TY3JlZW4oKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEhpZGVzIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuSGlkZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogTWF4aW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTWF4aW1pc2UoKSB7XHJcbiAgICB0aGlzV2luZG93Lk1heGltaXNlKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBVbi1tYXhpbWlzZXMgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBVbk1heGltaXNlKCkge1xyXG4gICAgdGhpc1dpbmRvdy5Vbk1heGltaXNlKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUb2dnbGVzIHRoZSBtYXhpbWlzYXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBUb2dnbGVNYXhpbWlzZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuVG9nZ2xlTWF4aW1pc2UoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE1pbmltaXNlKCkge1xyXG4gICAgdGhpc1dpbmRvdy5NaW5pbWlzZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVW5NaW5pbWlzZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuVW5NaW5pbWlzZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVzdG9yZXMgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBSZXN0b3JlKCkge1xyXG4gICAgdGhpc1dpbmRvdy5SZXN0b3JlKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XHJcbiAgICB0aGlzV2luZG93LlNob3coKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENsb3NlcyB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENsb3NlKCkge1xyXG4gICAgdGhpc1dpbmRvdy5DbG9zZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHIgLSBUaGUgcmVkIGNvbXBvbmVudCBvZiB0aGUgY29sb3VyLlxyXG4gKiBAcGFyYW0ge251bWJlcn0gZyAtIFRoZSBncmVlbiBjb21wb25lbnQgb2YgdGhlIGNvbG91ci5cclxuICogQHBhcmFtIHtudW1iZXJ9IGIgLSBUaGUgYmx1ZSBjb21wb25lbnQgb2YgdGhlIGNvbG91ci5cclxuICogQHBhcmFtIHtudW1iZXJ9IGEgLSBUaGUgYWxwaGEgY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0QmFja2dyb3VuZENvbG91cihyLCBnLCBiLCBhKSB7XHJcbiAgICB0aGlzV2luZG93LlNldEJhY2tncm91bmRDb2xvdXIociwgZywgYiwgYSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXHJcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gcmVzaXphYmxlIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSByZXNpemFibGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0UmVzaXphYmxlKHJlc2l6YWJsZSkge1xyXG4gICAgdGhpc1dpbmRvdy5TZXRSZXNpemFibGUocmVzaXphYmxlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2lkdGgoKSB7XHJcbiAgICByZXR1cm4gdGhpc1dpbmRvdy5XaWR0aCgpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSGVpZ2h0KCkge1xyXG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuSGVpZ2h0KCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBab29tcyBpbiB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFpvb21JbigpIHtcclxuICAgIHRoaXNXaW5kb3cuWm9vbUluKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBab29tcyBvdXQgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBab29tT3V0KCkge1xyXG4gICAgdGhpc1dpbmRvdy5ab29tT3V0KCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXNldHMgdGhlIHpvb20gb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBab29tUmVzZXQoKSB7XHJcbiAgICB0aGlzV2luZG93Llpvb21SZXNldCgpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldFpvb21MZXZlbCgpIHtcclxuICAgIHJldHVybiB0aGlzV2luZG93LkdldFpvb21MZXZlbCgpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0gem9vbUxldmVsIC0gVGhlIHpvb20gbGV2ZWwgdG8gc2V0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFpvb21MZXZlbCh6b29tTGV2ZWwpIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0Wm9vbUxldmVsKHpvb21MZXZlbCk7XHJcbn1cclxuIiwgIlxyXG5pbXBvcnQge0VtaXQsIFdhaWxzRXZlbnR9IGZyb20gXCIuL2V2ZW50c1wiO1xyXG5pbXBvcnQge1F1ZXN0aW9ufSBmcm9tIFwiLi9kaWFsb2dzXCI7XHJcbmltcG9ydCB7R2V0fSBmcm9tIFwiLi93aW5kb3dcIjtcclxuaW1wb3J0IHtPcGVuVVJMfSBmcm9tIFwiLi9icm93c2VyXCI7XHJcbmltcG9ydCB7ZGVidWdMb2d9IGZyb20gXCIuL2xvZ1wiO1xyXG5cclxuLyoqXHJcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gc2VuZC5cclxuICogQHBhcmFtIHthbnl9IFtkYXRhPW51bGxdIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZSwgZGF0YT1udWxsKSB7XHJcbiAgICBsZXQgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpO1xyXG4gICAgRW1pdChldmVudCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBZGRzIGV2ZW50IGxpc3RlbmVycyB0byBlbGVtZW50cyB3aXRoIGB3bWwtZXZlbnRgIGF0dHJpYnV0ZS5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XScpO1xyXG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xyXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKTtcclxuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XHJcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcclxuXHJcbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xyXG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xyXG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBEZXRhY2hlZDogZmFsc2UsIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0pO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgc3BlY2lmaWVkIHdpbmRvdy5cclxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGNhbGwgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZCAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cclxuICovXHJcbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TmFtZSwgbWV0aG9kKSB7XHJcbiAgICBsZXQgdGFyZ2V0V2luZG93ID0gR2V0KHdpbmRvd05hbWUpO1xyXG4gICAgbGV0IG1ldGhvZE1hcCA9IFdpbmRvd01ldGhvZHModGFyZ2V0V2luZG93KTtcclxuICAgIGlmICghbWV0aG9kTWFwLmhhcyhtZXRob2QpKSB7XHJcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xyXG4gICAgfVxyXG4gICAgdHJ5IHtcclxuICAgICAgICBtZXRob2RNYXAuZ2V0KG1ldGhvZCkoKTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICdcIiArIG1ldGhvZCArIFwiJzogXCIgKyBlKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEFkZHMgd2luZG93IGxpc3RlbmVycyBmb3IgZWxlbWVudHMgd2l0aCB0aGUgJ3dtbC13aW5kb3cnIGF0dHJpYnV0ZS5cclxuICogUmVtb3ZlcyBhbnkgZXhpc3RpbmcgbGlzdGVuZXJzIGJlZm9yZSBhZGRpbmcgbmV3IG9uZXMuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XHJcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtd2luZG93XScpO1xyXG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xyXG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93Jyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCAnY2xpY2snO1xyXG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdGFyZ2V0LXdpbmRvdycpIHx8ICcnO1xyXG5cclxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XHJcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0pO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQWRkcyBhIGxpc3RlbmVyIHRvIGVsZW1lbnRzIHdpdGggdGhlICd3bWwtb3BlbnVybCcgYXR0cmlidXRlLlxyXG4gKiBXaGVuIHRoZSBzcGVjaWZpZWQgdHJpZ2dlciBldmVudCBpcyBmaXJlZCBvbiBhbnkgb2YgdGhlc2UgZWxlbWVudHMsXHJcbiAqIHRoZSBsaXN0ZW5lciB3aWxsIG9wZW4gdGhlIFVSTCBzcGVjaWZpZWQgYnkgdGhlICd3bWwtb3BlbnVybCcgYXR0cmlidXRlLlxyXG4gKiBJZiBhICd3bWwtY29uZmlybScgYXR0cmlidXRlIGlzIHByb3ZpZGVkLCBhIGNvbmZpcm1hdGlvbiBkaWFsb2cgd2lsbCBiZSBkaXNwbGF5ZWQsXHJcbiAqIGFuZCB0aGUgVVJMIHdpbGwgb25seSBiZSBvcGVuZWQgaWYgdGhlIHVzZXIgY29uZmlybXMuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxPcGVuQnJvd3Nlckxpc3RlbmVyKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLW9wZW51cmxdJyk7XHJcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XHJcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfSk7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xyXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcblxyXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWxvYWRzIHRoZSBXTUwgcGFnZSBieSBhZGRpbmcgbmVjZXNzYXJ5IGV2ZW50IGxpc3RlbmVycyBhbmQgYnJvd3NlciBsaXN0ZW5lcnMuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCkge1xyXG4gICAgaWYoREVCVUcpIHtcclxuICAgICAgICBkZWJ1Z0xvZyhcIlJlbG9hZGluZyBXTUxcIik7XHJcbiAgICB9XHJcbiAgICBhZGRXTUxFdmVudExpc3RlbmVycygpO1xyXG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XHJcbiAgICBhZGRXTUxPcGVuQnJvd3Nlckxpc3RlbmVyKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgbWFwIG9mIGFsbCBtZXRob2RzIGluIHRoZSBjdXJyZW50IHdpbmRvdy5cclxuICogQHJldHVybnMge01hcH0gLSBBIG1hcCBvZiB3aW5kb3cgbWV0aG9kcy5cclxuICovXHJcbmZ1bmN0aW9uIFdpbmRvd01ldGhvZHModGFyZ2V0V2luZG93KSB7XHJcbiAgICAvLyBDcmVhdGUgYSBuZXcgbWFwIHRvIHN0b3JlIG1ldGhvZHNcclxuICAgIGxldCByZXN1bHQgPSBuZXcgTWFwKCk7XHJcblxyXG4gICAgLy8gSXRlcmF0ZSBvdmVyIGFsbCBwcm9wZXJ0aWVzIG9mIHRoZSB3aW5kb3cgb2JqZWN0XHJcbiAgICBmb3IgKGxldCBtZXRob2QgaW4gdGFyZ2V0V2luZG93KSB7XHJcbiAgICAgICAgLy8gQ2hlY2sgaWYgdGhlIHByb3BlcnR5IGlzIGluZGVlZCBhIG1ldGhvZCAoZnVuY3Rpb24pXHJcbiAgICAgICAgaWYodHlwZW9mIHRhcmdldFdpbmRvd1ttZXRob2RdID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICAgICAgICAgIC8vIEFkZCB0aGUgbWV0aG9kIHRvIHRoZSBtYXBcclxuICAgICAgICAgICAgcmVzdWx0LnNldChtZXRob2QsIHRhcmdldFdpbmRvd1ttZXRob2RdKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgfVxyXG4gICAgLy8gUmV0dXJuIHRoZSBtYXAgb2Ygd2luZG93IG1ldGhvZHNcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcclxuICovXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmltcG9ydCB7RXZlbnRUeXBlc30gZnJvbSBcIi4vZXZlbnRfdHlwZXNcIjtcclxuZXhwb3J0IGNvbnN0IFR5cGVzID0gRXZlbnRUeXBlcztcclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkV2ZW50cywgJycpO1xyXG5jb25zdCBFbWl0TWV0aG9kID0gMDtcclxuY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XHJcblxyXG5jbGFzcyBMaXN0ZW5lciB7XHJcbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcclxuICAgICAgICAgICAgY2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcclxuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcclxuICAgIGNvbnN0cnVjdG9yKG5hbWUsIGRhdGEgPSBudWxsKSB7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcclxuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gc2V0dXAoKSB7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudCkge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcclxuICAgIGlmIChsaXN0ZW5lcnMpIHtcclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+IHtcclxuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcclxuICAgICAgICAgICAgaWYgKHJlbW92ZSkgcmV0dXJuIHRydWU7XHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XHJcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XHJcbiAgICAgICAgICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxyXG4gKlxyXG4gQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XHJcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQuIEl0IHRha2VzIG5vIHBhcmFtZXRlcnMuXHJcbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7IHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTsgfVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgZXZlbnQgb2NjdXJzLlxyXG4gKiBAcmV0dXJuIHtmdW5jdGlvbn0gLSBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykgeyByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTsgfVxyXG5cclxuLyoqXHJcbiAqIFJlbW92ZXMgdGhlIHNwZWNpZmllZCBsaXN0ZW5lciBmcm9tIHRoZSBldmVudCBsaXN0ZW5lcnMgY29sbGVjdGlvbi5cclxuICogSWYgYWxsIGxpc3RlbmVycyBmb3IgdGhlIGV2ZW50IGFyZSByZW1vdmVkLCB0aGUgZXZlbnQga2V5IGlzIGRlbGV0ZWQgZnJvbSB0aGUgY29sbGVjdGlvbi5cclxuICpcclxuICogQHBhcmFtIHtPYmplY3R9IGxpc3RlbmVyIC0gVGhlIGxpc3RlbmVyIHRvIGJlIHJlbW92ZWQuXHJcbiAqL1xyXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xyXG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcclxuICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cclxuICogQHBhcmFtIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzIC0gQWRkaXRpb25hbCBldmVudCBuYW1lcyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cclxuICogQHJldHVybiB7dW5kZWZpbmVkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XHJcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XHJcbiAgICBldmVudHNUb1JlbW92ZS5mb3JFYWNoKGV2ZW50TmFtZSA9PiBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKSk7XHJcbn1cclxuLyoqXHJcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cclxuICpcclxuICogQGZ1bmN0aW9uIE9mZkFsbFxyXG4gKiBAcmV0dXJucyB7dm9pZH1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7IGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7IH1cclxuXHJcbi8qKlxyXG4gKiBFbWl0cyBhbiBldmVudCB1c2luZyB0aGUgZ2l2ZW4gZXZlbnQgbmFtZS5cclxuICpcclxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxyXG4gKiBAcmV0dXJucyB7YW55fSAtIFRoZSByZXN1bHQgb2YgdGhlIGVtaXR0ZWQgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkgeyByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCBldmVudCk7IH1cclxuIiwgIlxyXG5leHBvcnQgY29uc3QgRXZlbnRUeXBlcyA9IHtcclxuXHRXaW5kb3dzOiB7XHJcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwid2luZG93czpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcclxuXHRcdEFQTVBvd2VyU3RhdHVzQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTdGF0dXNDaGFuZ2VcIixcclxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXHJcblx0XHRBUE1SZXN1bWVBdXRvbWF0aWM6IFwid2luZG93czpBUE1SZXN1bWVBdXRvbWF0aWNcIixcclxuXHRcdEFQTVJlc3VtZVN1c3BlbmQ6IFwid2luZG93czpBUE1SZXN1bWVTdXNwZW5kXCIsXHJcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcclxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxyXG5cdFx0V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWQ6IFwid2luZG93czpXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZFwiLFxyXG5cdFx0V2luZG93SW5hY3RpdmU6IFwid2luZG93czpXaW5kb3dJbmFjdGl2ZVwiLFxyXG5cdFx0V2luZG93QWN0aXZlOiBcIndpbmRvd3M6V2luZG93QWN0aXZlXCIsXHJcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXHJcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd01heGltaXNlXCIsXHJcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd0Z1bGxzY3JlZW5cIixcclxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxyXG5cdFx0V2luZG93UmVzdG9yZTogXCJ3aW5kb3dzOldpbmRvd1Jlc3RvcmVcIixcclxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcclxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dDbG9zZTogXCJ3aW5kb3dzOldpbmRvd0Nsb3NlXCIsXHJcblx0XHRXaW5kb3dTZXRGb2N1czogXCJ3aW5kb3dzOldpbmRvd1NldEZvY3VzXCIsXHJcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcclxuXHRcdFdpbmRvd0RyYWdEcm9wOiBcIndpbmRvd3M6V2luZG93RHJhZ0Ryb3BcIixcclxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxyXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXHJcblx0XHRXaW5kb3dEcmFnT3ZlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdPdmVyXCIsXHJcblx0fSxcclxuXHRNYWM6IHtcclxuXHRcdEFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvblwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkSGlkZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVOb3RpZmljYXRpb246IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlTm90aWZpY2F0aW9uXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkVXBkYXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFVwZGF0ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxIaWRlXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVW5oaWRlXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZSFcIixcclxuXHRcdEFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuOiBcIm1hYzpBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbiFcIixcclxuXHRcdFdpbmRvd0RpZEJlY29tZUtleTogXCJtYWM6V2luZG93RGlkQmVjb21lS2V5XCIsXHJcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXHJcblx0XHRXaW5kb3dEaWRCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dEaWRCZWdpblNoZWV0XCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VBbHBoYTogXCJtYWM6V2luZG93RGlkQ2hhbmdlQWxwaGFcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllc1wiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlVGl0bGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVRpdGxlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VWaXNpYmlsaXR5XCIsXHJcblx0XHRXaW5kb3dEaWREZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWREZW1pbmlhdHVyaXplXCIsXHJcblx0XHRXaW5kb3dEaWRFbmRTaGVldDogXCJtYWM6V2luZG93RGlkRW5kU2hlZXRcIixcclxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXHJcblx0XHRXaW5kb3dEaWRFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRXhpdEZ1bGxTY3JlZW5cIixcclxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXHJcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxyXG5cdFx0V2luZG93RGlkRm9jdXM6IFwibWFjOldpbmRvd0RpZEZvY3VzXCIsXHJcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcclxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcclxuXHRcdFdpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9mZlNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcclxuXHRcdFdpbmRvd0RpZFJlc2lnbk1haW46IFwibWFjOldpbmRvd0RpZFJlc2lnbk1haW5cIixcclxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxyXG5cdFx0V2luZG93RGlkVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUFscGhhXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxyXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVNoYWRvd1wiLFxyXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5XCIsXHJcblx0XHRXaW5kb3dTaG91bGRDbG9zZTogXCJtYWM6V2luZG93U2hvdWxkQ2xvc2UhXCIsXHJcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXHJcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcclxuXHRcdFdpbmRvd1dpbGxCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dXaWxsQmVnaW5TaGVldFwiLFxyXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxyXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcclxuXHRcdFdpbmRvd1dpbGxEZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZVwiLFxyXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxyXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXHJcblx0XHRXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd1dpbGxFeGl0RnVsbFNjcmVlblwiLFxyXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxyXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcclxuXHRcdFdpbmRvd1dpbGxNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbE1pbmlhdHVyaXplXCIsXHJcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcclxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dXaWxsT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT25TY3JlZW5cIixcclxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxyXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxyXG5cdFx0V2luZG93V2lsbFVuZm9jdXM6IFwibWFjOldpbmRvd1dpbGxVbmZvY3VzXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRpdGxlXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxyXG5cdFx0V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWU6IFwibWFjOldpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lXCIsXHJcblx0XHRNZW51V2lsbE9wZW46IFwibWFjOk1lbnVXaWxsT3BlblwiLFxyXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXHJcblx0XHRNZW51RGlkQ2xvc2U6IFwibWFjOk1lbnVEaWRDbG9zZVwiLFxyXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25cIixcclxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxyXG5cdFx0TWVudVdpbGxIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51V2lsbEhpZ2hsaWdodEl0ZW1cIixcclxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxyXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxyXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcclxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXHJcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcclxuXHRcdE1lbnVXaWxsUmVtb3ZlSXRlbTogXCJtYWM6TWVudVdpbGxSZW1vdmVJdGVtXCIsXHJcblx0XHRNZW51RGlkUmVtb3ZlSXRlbTogXCJtYWM6TWVudURpZFJlbW92ZUl0ZW1cIixcclxuXHRcdE1lbnVXaWxsQmVnaW5UcmFja2luZzogXCJtYWM6TWVudVdpbGxCZWdpblRyYWNraW5nXCIsXHJcblx0XHRNZW51RGlkQmVnaW5UcmFja2luZzogXCJtYWM6TWVudURpZEJlZ2luVHJhY2tpbmdcIixcclxuXHRcdE1lbnVXaWxsRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsRW5kVHJhY2tpbmdcIixcclxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXHJcblx0XHRNZW51V2lsbFVwZGF0ZTogXCJtYWM6TWVudVdpbGxVcGRhdGVcIixcclxuXHRcdE1lbnVEaWRVcGRhdGU6IFwibWFjOk1lbnVEaWRVcGRhdGVcIixcclxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcclxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXHJcblx0XHRNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbVwiLFxyXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXHJcblx0XHRXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXHJcblx0XHRXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb25cIixcclxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWRcIixcclxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXHJcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZFwiLFxyXG5cdH0sXHJcblx0TGludXg6IHtcclxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcclxuXHR9LFxyXG5cdENvbW1vbjoge1xyXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcImNvbW1vbjpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcclxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxyXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcclxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwiY29tbW9uOldpbmRvd1Jlc3RvcmVcIixcclxuXHRcdFdpbmRvd01pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dNaW5pbWlzZVwiLFxyXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJjb21tb246V2luZG93VW5NaW5pbWlzZVwiLFxyXG5cdFx0V2luZG93Q2xvc2luZzogXCJjb21tb246V2luZG93Q2xvc2luZ1wiLFxyXG5cdFx0V2luZG93Wm9vbTogXCJjb21tb246V2luZG93Wm9vbVwiLFxyXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcclxuXHRcdFdpbmRvd1pvb21PdXQ6IFwiY29tbW9uOldpbmRvd1pvb21PdXRcIixcclxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJjb21tb246V2luZG93Wm9vbVJlc2V0XCIsXHJcblx0XHRXaW5kb3dGb2N1czogXCJjb21tb246V2luZG93Rm9jdXNcIixcclxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJjb21tb246V2luZG93TG9zdEZvY3VzXCIsXHJcblx0XHRXaW5kb3dTaG93OiBcImNvbW1vbjpXaW5kb3dTaG93XCIsXHJcblx0XHRXaW5kb3dIaWRlOiBcImNvbW1vbjpXaW5kb3dIaWRlXCIsXHJcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcImNvbW1vbjpXaW5kb3dEUElDaGFuZ2VkXCIsXHJcblx0XHRXaW5kb3dGaWxlc0Ryb3BwZWQ6IFwiY29tbW9uOldpbmRvd0ZpbGVzRHJvcHBlZFwiLFxyXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcclxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXHJcblx0fSxcclxufTtcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBPcGVuRmlsZURpYWxvZ09wdGlvbnNcclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DcmVhdGVEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Nob3dIaWRkZW5GaWxlc10gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb25dIC0gSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtIaWRlRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtBbGxvd3NPdGhlckZpbGV0eXBlc10gLSBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC5cclxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaXRsZSBvZiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtEaXJlY3RvcnldIC0gRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy5cclxuICovXHJcblxyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNhdmVGaWxlRGlhbG9nT3B0aW9uc1xyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0ZpbGVuYW1lXSAtIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VGaWxlc10gLSBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ3JlYXRlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtTaG93SGlkZGVuRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Jlc29sdmVzQWxpYXNlc10gLSBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c011bHRpcGxlU2VsZWN0aW9uXSAtIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbSGlkZUV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzT3RoZXJGaWxldHlwZXNdIC0gSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7RmlsZUZpbHRlcltdfSBbRmlsdGVyc10gLSBBcnJheSBvZiBmaWxlIGZpbHRlcnMuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGl0bGUgb2YgdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtNZXNzYWdlXSAtIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0J1dHRvblRleHRdIC0gVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlyZWN0b3J5XSAtIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0RldGFjaGVkXSAtIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuXHJcbiAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IE1lc3NhZ2VEaWFsb2dPcHRpb25zXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGhlIHRpdGxlIG9mIHRoZSBkaWFsb2cgd2luZG93LlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7QnV0dG9uW119IFtCdXR0b25zXSAtIEFycmF5IG9mIGJ1dHRvbiBvcHRpb25zIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBCdXR0b25cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtMYWJlbF0gLSBUZXh0IHRoYXQgYXBwZWFycyB3aXRoaW4gdGhlIGJ1dHRvbi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNDYW5jZWxdIC0gVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNEZWZhdWx0XSAtIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci5cclxuICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gRmlsZUZpbHRlclxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0Rpc3BsYXlOYW1lXSAtIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1BhdHRlcm5dIC0gUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy5cclxuICovXHJcblxyXG4vLyBzZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5kaWFsb2dFcnJvckNhbGxiYWNrID0gZGlhbG9nRXJyb3JDYWxsYmFjaztcclxud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xyXG5cclxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcclxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XHJcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xyXG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XHJcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcclxuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xyXG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5EaWFsb2csICcnKTtcclxuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xyXG5cclxuLyoqXHJcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBpZCB0aGF0IGlzIG5vdCBwcmVzZW50IGluIGRpYWxvZ1Jlc3BvbnNlcy5cclxuICogQHJldHVybnMge3N0cmluZ30gdW5pcXVlIGlkXHJcbiAqL1xyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHR5cGUgLSB0eXBlIG9mIGRpYWxvZ1xyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfE9wZW5GaWxlRGlhbG9nT3B0aW9uc3xTYXZlRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBvcHRpb25zIGZvciB0aGUgZGlhbG9nXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlfSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nXHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucyA9IHt9KSB7XHJcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuICAgIG9wdGlvbnNbXCJkaWFsb2ctaWRcIl0gPSBpZDtcclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIHRoZSBjYWxsYmFjayBmcm9tIGEgZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgSUQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cclxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgZGF0YSByZWNlaXZlZCBmcm9tIHRoZSBkaWFsb2cuXHJcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gaXNKU09OIC0gRmxhZyBpbmRpY2F0aW5nIHdoZXRoZXIgdGhlIGRhdGEgaXMgaW4gSlNPTiBmb3JtYXQuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZ1Jlc3VsdENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsYmFjayBmdW5jdGlvbiBmb3IgaGFuZGxpbmcgZXJyb3JzIGluIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSBkaWFsb2cgcmVzcG9uc2UuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XHJcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgaWYgKHApIHtcclxuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5cclxuXHJcbi8vIFJlcGxhY2UgYG1ldGhvZHNgIHdpdGggY29uc3RhbnRzIGluIFRpdGxlIENhc2VcclxuXHJcbi8qKlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgSW5mbyA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7XHJcblxyXG4vKipcclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IFdhcm5pbmcgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBFcnJvciA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBRdWVzdGlvbiA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7T3BlbkZpbGVEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nW118c3RyaW5nPn0gUmV0dXJucyBzZWxlY3RlZCBmaWxlIG9yIGxpc3Qgb2YgZmlsZXMuIFJldHVybnMgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgT3BlbkZpbGUgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKTtcclxuXHJcbi8qKlxyXG4gKiBAcGFyYW0ge1NhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUuIFJldHVybnMgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgU2F2ZUZpbGUgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTtcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcclxud2luZG93Ll93YWlscy5jYWxsRXJyb3JIYW5kbGVyID0gZXJyb3JIYW5kbGVyO1xyXG5cclxuXHJcbmNvbnN0IENhbGxCaW5kaW5nID0gMDtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2FsbCwgJycpO1xyXG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcclxuXHJcbi8qKlxyXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtzdHJpbmd9IC0gQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgY2FsbFJlc3BvbnNlcyBzZXQuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIHJlc3VsdCBvZiBhIGNhbGwgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cclxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgcmVzdWx0IGRhdGEgb2YgdGhlIHJlcXVlc3QuXHJcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH0gLSBUaGlzIG1ldGhvZCBkb2VzIG5vdCByZXR1cm4gYW55IHZhbHVlLlxyXG4gKi9cclxuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZCwgZGF0YSwgaXNKU09OKSB7XHJcbiAgICBjb25zdCBwcm9taXNlSGFuZGxlciA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcclxuICAgIGlmIChwcm9taXNlSGFuZGxlcikge1xyXG4gICAgICAgIHByb21pc2VIYW5kbGVyLnJlc29sdmUoaXNKU09OID8gSlNPTi5wYXJzZShkYXRhKSA6IGRhdGEpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogSGFuZGxlcyB0aGUgZXJyb3IgZnJvbSBhIGNhbGwgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSBwcm9taXNlIGhhbmRsZXIuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UgdG8gcmVqZWN0IHRoZSBwcm9taXNlIGhhbmRsZXIgd2l0aC5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZCwgbWVzc2FnZSkge1xyXG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XHJcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcclxuICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobWVzc2FnZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgY2FsbFJlc3BvbnNlcyBtYXAuXHJcbiAqXHJcbiAqIEBwYXJhbSB7YW55fSBpZCAtIFRoZSBJRCBvZiB0aGUgcmVzcG9uc2UgdG8gYmUgcmV0cmlldmVkIGFuZCByZW1vdmVkLlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7YW55fSBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZXRBbmREZWxldGVSZXNwb25zZShpZCkge1xyXG4gICAgY29uc3QgcmVzcG9uc2UgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICByZXR1cm4gcmVzcG9uc2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFeGVjdXRlcyBhIGNhbGwgdXNpbmcgdGhlIHByb3ZpZGVkIHR5cGUgYW5kIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfG51bWJlcn0gdHlwZSAtIFRoZSB0eXBlIG9mIGNhbGwgdG8gZXhlY3V0ZS5cclxuICogQHBhcmFtIHtPYmplY3R9IFtvcHRpb25zPXt9XSAtIEFkZGl0aW9uYWwgb3B0aW9ucyBmb3IgdGhlIGNhbGwuXHJcbiAqIEByZXR1cm4ge1Byb21pc2V9IC0gQSBwcm9taXNlIHRoYXQgd2lsbCBiZSByZXNvbHZlZCBvciByZWplY3RlZCBiYXNlZCBvbiB0aGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucyA9IHt9KSB7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbCBtZXRob2QuXHJcbiAqXHJcbiAqIEBwYXJhbSB7T2JqZWN0fSBvcHRpb25zIC0gVGhlIG9wdGlvbnMgZm9yIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm5zIHtPYmplY3R9IC0gVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFeGVjdXRlcyBhIG1ldGhvZCBieSBuYW1lLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJy5cclxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEB0aHJvd3Mge0Vycm9yfSBJZiB0aGUgbmFtZSBpcyBub3QgYSBzdHJpbmcgb3IgaXMgbm90IGluIHRoZSBjb3JyZWN0IGZvcm1hdC5cclxuICogQHJldHVybnMgeyp9IFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBleGVjdXRpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG5hbWUsIC4uLmFyZ3MpIHtcclxuICAgIGlmICh0eXBlb2YgbmFtZSAhPT0gXCJzdHJpbmdcIiB8fCBuYW1lLnNwbGl0KFwiLlwiKS5sZW5ndGggIT09IDMpIHtcclxuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJDYWxsQnlOYW1lIHJlcXVpcmVzIGEgc3RyaW5nIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCdcIik7XHJcbiAgICB9XHJcbiAgICBsZXQgW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXSA9IG5hbWUuc3BsaXQoXCIuXCIpO1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWUsXHJcbiAgICAgICAgc3RydWN0TmFtZSxcclxuICAgICAgICBtZXRob2ROYW1lLFxyXG4gICAgICAgIGFyZ3NcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQsIC4uLmFyZ3MpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xyXG4gICAgICAgIG1ldGhvZElELFxyXG4gICAgICAgIGFyZ3NcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBwbHVnaW4uXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHBsdWdpbi5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJucyB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxyXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXHJcbiAgICAgICAgbWV0aG9kTmFtZSxcclxuICAgICAgICBhcmdzXHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtkZWJ1Z0xvZ30gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZ1wiO1xyXG5cclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb25cIjtcclxuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvYnJvd3NlclwiO1xyXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZFwiO1xyXG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3NcIjtcclxuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvc2NyZWVuc1wiO1xyXG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbVwiO1xyXG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvd1wiO1xyXG5pbXBvcnQgKiBhcyBXTUwgZnJvbSAnLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sJztcclxuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudHNcIjtcclxuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9nc1wiO1xyXG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxsc1wiO1xyXG5pbXBvcnQge2ludm9rZX0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbVwiO1xyXG5cclxuLyoqKlxyXG4gVGhpcyB0ZWNobmlxdWUgZm9yIHByb3BlciBsb2FkIGRldGVjdGlvbiBpcyB0YWtlbiBmcm9tIEhUTVg6XHJcblxyXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcclxuXHJcbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcclxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXHJcblxyXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XHJcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcclxuXHJcbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcclxuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxyXG5cclxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcclxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cclxuIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLlxyXG5cclxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXHJcbiBBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFXHJcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcclxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcclxuIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMXHJcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxyXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcclxuIENBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksXHJcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxyXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cclxuXHJcbiAqKiovXHJcblxyXG53aW5kb3cuX3dhaWxzLmludm9rZT1pbnZva2U7XHJcblxyXG53aW5kb3cud2FpbHMgPSB3aW5kb3cud2FpbHMgfHwge307XHJcbndpbmRvdy53YWlscy5BcHBsaWNhdGlvbiA9IEFwcGxpY2F0aW9uO1xyXG53aW5kb3cud2FpbHMuQnJvd3NlciA9IEJyb3dzZXI7XHJcbndpbmRvdy53YWlscy5DYWxsID0gQ2FsbDtcclxud2luZG93LndhaWxzLkNsaXBib2FyZCA9IENsaXBib2FyZDtcclxud2luZG93LndhaWxzLkRpYWxvZ3MgPSBEaWFsb2dzO1xyXG53aW5kb3cud2FpbHMuRXZlbnRzID0gRXZlbnRzO1xyXG53aW5kb3cud2FpbHMuRmxhZ3MgPSBGbGFncztcclxud2luZG93LndhaWxzLlNjcmVlbnMgPSBTY3JlZW5zO1xyXG53aW5kb3cud2FpbHMuU3lzdGVtID0gU3lzdGVtO1xyXG53aW5kb3cud2FpbHMuV2luZG93ID0gV2luZG93O1xyXG53aW5kb3cud2FpbHMuV01MID0gV01MO1xyXG5cclxuXHJcbmxldCBpc1JlYWR5ID0gZmFsc2VcclxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGZ1bmN0aW9uKCkge1xyXG4gICAgaXNSZWFkeSA9IHRydWVcclxuICAgIHdpbmRvdy5fd2FpbHMuaW52b2tlKCd3YWlsczpydW50aW1lOnJlYWR5Jyk7XHJcbiAgICBpZihERUJVRykge1xyXG4gICAgICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIik7XHJcbiAgICB9XHJcbn0pXHJcblxyXG5mdW5jdGlvbiB3aGVuUmVhZHkoZm4pIHtcclxuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcclxuICAgICAgICBmbigpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgZm4pO1xyXG4gICAgfVxyXG59XHJcblxyXG53aGVuUmVhZHkoKCkgPT4ge1xyXG4gICAgV01MLlJlbG9hZCgpO1xyXG59KTtcclxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7QUFLTyxXQUFTLFNBQVMsU0FBUztBQUU5QixZQUFRO0FBQUEsTUFDSixrQkFBa0IsVUFBVTtBQUFBLE1BQzVCO0FBQUEsTUFDQTtBQUFBLElBQ0o7QUFBQSxFQUNKOzs7QUNaQTtBQUFBO0FBQUE7QUFBQTtB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LE1BQ3hCLG1CQUFtQjtBQUFBLE1BQ25CLDBCQUEwQjtBQUFBLE1BQzFCLDhCQUE4QjtBQUFBLE1BQzlCLHlCQUF5QjtBQUFBLE1BQ3pCLDZCQUE2QjtBQUFBLE1BQzdCLGlCQUFpQjtBQUFBLE1BQ2pCLGdCQUFnQjtBQUFBLE1BQ2hCLHNCQUFzQjtBQUFBLE1BQ3RCLGVBQWU7QUFBQSxNQUNmLHlCQUF5QjtBQUFBLE1BQ3pCLHdCQUF3QjtBQUFBLE1BQ3hCLG9CQUFvQjtBQUFBLE1BQ3BCLHFCQUFxQjtBQUFBLE1BQ3JCLGlCQUFpQjtBQUFBLE1BQ2pCLGlCQUFpQjtBQUFBLE1BQ2pCLHNCQUFzQjtBQUFBLE1BQ3RCLG1DQUFtQztBQUFBLE1BQ25DLHFDQUFxQztBQUFBLE1BQ3JDLHVCQUF1QjtBQUFBLE1BQ3ZCLHNCQUFzQjtBQUFBLE1BQ3RCLHdCQUF3QjtBQUFBLE1BQ3hCLDJCQUEyQjtBQUFBLE1BQzNCLG1CQUFtQjtBQUFBLE1BQ25CLHFCQUFxQjtBQUFBLE1BQ3JCLHNCQUFzQjtBQUFBLE1BQ3RCLHNCQUFzQjtBQUFBLE1BQ3RCLDhCQUE4QjtBQUFBLE1BQzlCLGlCQUFpQjtBQUFBLE1BQ2pCLHlCQUF5QjtBQUFBLE1BQ3pCLDJCQUEyQjtBQUFBLE1BQzNCLCtCQUErQjtBQUFBLE1BQy9CLDBCQUEwQjtBQUFBLE1BQzFCLDhCQUE4QjtBQUFBLE1BQzlCLGlCQUFpQjtBQUFBLE1BQ2pCLHVCQUF1QjtBQUFBLE1BQ3ZCLGdCQUFnQjtBQUFBLE1BQ2hCLDBCQUEwQjtBQUFBLE1BQzFCLHlCQUF5QjtBQUFBLE1BQ3pCLHNCQUFzQjtBQUFBLE1BQ3RCLGtCQUFrQjtBQUFBLE1BQ2xCLG1CQUFtQjtBQUFBLE1BQ25CLGtCQUFrQjtBQUFBLE1BQ2xCLHVCQUF1QjtBQUFBLE1BQ3ZCLG9DQUFvQztBQUFBLE1BQ3BDLHNDQUFzQztBQUFBLE1BQ3RDLHdCQUF3QjtBQUFBLE1BQ3hCLHVCQUF1QjtBQUFBLE1BQ3ZCLHlCQUF5QjtBQUFBLE1BQ3pCLDRCQUE0QjtBQUFBLE1BQzVCLDRCQUE0QjtBQUFBLE1B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LE9BQU8sU0FBUyxNQUFNO0FBQ2xDLFFBQUksT0FBTyxTQUFTLFlBQVksS0FBSyxNQUFNLEdBQUcsRUFBRSxXQUFXLEdBQUc7QUFDMUQsWUFBTSxJQUFJLE1BQU0sb0VBQW9FO0FBQUEsSUFDeEY7QUFDQSxRQUFJLENBQUMsYUFBYSxZQUFZLFVBQVUsSUFBSSxLQUFLLE1BQU0sR0FBRztBQUMxRCxXQUFPLFlBQVksYUFBYTtBQUFBLE1BQzVCO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsSUFDSixDQUFDO0FBQUEsRUFDTDtBQVNPLFdBQVMsS0FBSyxhQUFhLE1BQU07QUFDcEMsV0FBTyxZQUFZLGFBQWE7QUFBQSxNQUM1QjtBQUFBLE1BQ0E7QUFBQSxJQUNKLENBQUM7QUFBQSxFQUNMO0FBVU8sV0FBUyxPQUFPLFlBQVksZUFBZSxNQUFNO0FBQ3BELFdBQU8sWUFBWSxhQUFhO0FBQUEsTUFDNUIsYUFBYTtBQUFBLE1BQ2IsWUFBWTtBQUFBLE1BQ1o7QUFBQSxNQUNBO0FBQUEsSUFDSixDQUFDO0FBQUEsRUFDTDs7O0FDcEpBLFNBQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQThDbEMsU0FBTyxPQUFPLFNBQU87QUFFckIsU0FBTyxRQUFRLE9BQU8sU0FBUyxDQUFDO0FBQ2hDLFNBQU8sTUFBTSxjQUFjO0FBQzNCLFNBQU8sTUFBTSxVQUFVO0FBQ3ZCLFNBQU8sTUFBTSxPQUFPO0FBQ3BCLFNBQU8sTUFBTSxZQUFZO0FBQ3pCLFNBQU8sTUFBTSxVQUFVO0FBQ3ZCLFNBQU8sTUFBTSxTQUFTO0FBQ3RCLFNBQU8sTUFBTSxRQUFRO0FBQ3JCLFNBQU8sTUFBTSxVQUFVO0FBQ3ZCLFNBQU8sTUFBTSxTQUFTO0FBQ3RCLFNBQU8sTUFBTSxTQUFTO0FBQ3RCLFNBQU8sTUFBTSxNQUFNO0FBR25CLE1BQUksVUFBVTtBQUNkLFdBQVMsaUJBQWlCLG9CQUFvQixXQUFXO0FBQ3JELGNBQVU7QUFDVixXQUFPLE9BQU8sT0FBTyxxQkFBcUI7QUFDMUMsUUFBRyxNQUFPO0FBQ04sZUFBUyxzQkFBc0I7QUFBQSxJQUNuQztBQUFBLEVBQ0osQ0FBQztBQUVELFdBQVMsVUFBVSxJQUFJO0FBQ25CLFFBQUksV0FBVyxTQUFTLGVBQWUsWUFBWTtBQUMvQyxTQUFHO0FBQUEsSUFDUCxPQUFPO0FBQ0gsZUFBUyxpQkFBaUIsb0JBQW9CLEVBQUU7QUFBQSxJQUNwRDtBQUFBLEVBQ0o7QUFFQSxZQUFVLE1BQU07QUFDWixJQUFJLE9BQU87QUFBQSxFQUNmLENBQUM7IiwKICAibmFtZXMiOiBbInNpemUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIaWRlIiwgIlNob3ciLCAiY2FsbCIsICJ3aWR0aCIsICJoZWlnaHQiLCAiSGlkZSIsICJTaG93IiwgImNhbGwiLCAiZXZlbnROYW1lIiwgIkVycm9yIiwgImNhbGwiLCAiRXJyb3IiLCAiY2FsbCIsICJnZW5lcmF0ZUlEIl0KfQo=
