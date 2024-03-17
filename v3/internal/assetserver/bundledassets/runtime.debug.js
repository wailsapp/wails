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
    return window.webkit.messageHandlers.external.postMessage;
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb24uanMiLCAiLi4vLi4vLi4vcnVudGltZS9ub2RlX21vZHVsZXMvbmFub2lkL25vbi1zZWN1cmUvaW5kZXguanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZC5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvd2luZG93LmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50cy5qcyIsICIuLi8uLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RpYWxvZ3MuanMiLCAiLi4vLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxzLmpzIiwgIi4uLy4uLy4uL3J1bnRpbWUvZGVza3RvcC9jb21waWxlZC9tYWluLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyIvKipcclxuICogTG9ncyBhIG1lc3NhZ2UgdG8gdGhlIGNvbnNvbGUgd2l0aCBjdXN0b20gZm9ybWF0dGluZy5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZSkge1xyXG4gICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXHJcbiAgICBjb25zb2xlLmxvZyhcclxuICAgICAgICAnJWMgd2FpbHMzICVjICcgKyBtZXNzYWdlICsgJyAnLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjYWEwMDAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogM3B4IDBweCAwcHggM3B4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogIzAwOTkwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDBweCAzcHggM3B4IDBweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbSdcclxuICAgICk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uLCAnJyk7XHJcblxyXG5jb25zdCBIaWRlTWV0aG9kID0gMDtcclxuY29uc3QgU2hvd01ldGhvZCA9IDE7XHJcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xyXG5cclxuLyoqXHJcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cclxuICpcclxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICpcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyB0aGUgU2hvd01ldGhvZCBhbmQgcmV0dXJucyB0aGUgcmVzdWx0LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XHJcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBRdWl0TWV0aG9kIHRvIHRlcm1pbmF0ZSB0aGUgcHJvZ3JhbS5cclxuICpcclxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XHJcbn1cclxuIiwgImxldCB1cmxBbHBoYWJldCA9XG4gICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuZXhwb3J0IGxldCBjdXN0b21BbHBoYWJldCA9IChhbHBoYWJldCwgZGVmYXVsdFNpemUgPSAyMSkgPT4ge1xuICByZXR1cm4gKHNpemUgPSBkZWZhdWx0U2l6ZSkgPT4ge1xuICAgIGxldCBpZCA9ICcnXG4gICAgbGV0IGkgPSBzaXplXG4gICAgd2hpbGUgKGktLSkge1xuICAgICAgaWQgKz0gYWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiBhbHBoYWJldC5sZW5ndGgpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG4gIH1cbn1cbmV4cG9ydCBsZXQgbmFub2lkID0gKHNpemUgPSAyMSkgPT4ge1xuICBsZXQgaWQgPSAnJ1xuICBsZXQgaSA9IHNpemVcbiAgd2hpbGUgKGktLSkge1xuICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgfVxuICByZXR1cm4gaWRcbn1cbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xyXG5cclxuLy8gT2JqZWN0IE5hbWVzXHJcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcclxuICAgIENhbGw6IDAsXHJcbiAgICBDbGlwYm9hcmQ6IDEsXHJcbiAgICBBcHBsaWNhdGlvbjogMixcclxuICAgIEV2ZW50czogMyxcclxuICAgIENvbnRleHRNZW51OiA0LFxyXG4gICAgRGlhbG9nOiA1LFxyXG4gICAgV2luZG93OiA2LFxyXG4gICAgU2NyZWVuczogNyxcclxuICAgIFN5c3RlbTogOCxcclxuICAgIEJyb3dzZXI6IDksXHJcbn1cclxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xyXG5cclxuLyoqXHJcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cclxuICpcclxuICogQHBhcmFtIHtPYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3Qgb24gd2hpY2ggdGhlIG1ldGhvZCBpcyB0byBiZSBpbnZva2VkLlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXHJcbiAqIEByZXR1cm5zIHtGdW5jdGlvbn0gQSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IHRha2VzIHRoZSBtZXRob2QgbmFtZSBhbmQgb3B0aW9uYWxseSBhcmd1bWVudHMgYW5kIGludm9rZXMgdGhlIG1ldGhvZCB3aXRoaW4gdGhlIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xyXG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogQ3JlYXRlcyBhIG5ldyBydW50aW1lIGNhbGxlciB3aXRoIHNwZWNpZmllZCBJRC5cclxuICpcclxuICogQHBhcmFtIHtvYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICogQHJldHVybiB7RnVuY3Rpb259IC0gVGhlIG5ldyBydW50aW1lIGNhbGxlciBmdW5jdGlvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdCwgd2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xyXG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9O1xyXG59XHJcblxyXG5cclxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcclxuICAgIGlmKCBtZXRob2QgKSB7XHJcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIH1cclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG5cclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgZmV0Y2godXJsLCBmZXRjaE9wdGlvbnMpXHJcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcclxuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xyXG4gICAgICAgICAgICAgICAgICAgIC8vIGNoZWNrIGNvbnRlbnQgdHlwZVxyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcclxuICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XHJcbiAgICAgICAgICAgIH0pXHJcbiAgICAgICAgICAgIC50aGVuKGRhdGEgPT4gcmVzb2x2ZShkYXRhKSlcclxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5Ccm93c2VyLCAnJyk7XHJcbmNvbnN0IEJyb3dzZXJPcGVuVVJMID0gMDtcclxuXHJcbi8qKlxyXG4gKiBPcGVuIGEgYnJvd3NlciB3aW5kb3cgdG8gdGhlIGdpdmVuIFVSTFxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIC0gVGhlIFVSTCB0byBvcGVuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmwpIHtcclxuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsfSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCwgJycpO1xyXG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcclxuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBrZXlTdHJpbmcgLSBUaGUga2V5IHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBmb3IuXHJcbiAqIEByZXR1cm4geyp9IC0gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleVN0cmluZykge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gd2luZG93Ll93YWlscy5mbGFnc1trZXlTdHJpbmddO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleVN0cmluZyArIFwiJzogXCIgKyBlKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFBvc2l0aW9uXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZS5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgWSBjb29yZGluYXRlLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTaXplXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIHdpZHRoLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gWSAtIFRoZSBoZWlnaHQuXHJcbiAqL1xyXG5cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBSZWN0XHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgdG9wLWxlZnQgY29ybmVyLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gWSAtIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHRvcC1sZWZ0IGNvcm5lci5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBIZWlnaHQgLSBUaGUgaGVpZ2h0IG9mIHRoZSByZWN0YW5nbGUuXHJcbiAqL1xyXG5cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7KCdaZXJvJ3wnTmluZXR5J3wnT25lRWlnaHR5J3wnVHdvU2V2ZW50eScpfSBSb3RhdGlvblxyXG4gKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gQ2FuIGJlIG9uZSBvZiAnWmVybycsICdOaW5ldHknLCAnT25lRWlnaHR5JywgJ1R3b1NldmVudHknLlxyXG4gKi9cclxuXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gU2NyZWVuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gTmFtZSAtIEh1bWFuIHJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFNjYWxlIC0gVGhlIHJlc29sdXRpb24gc2NhbGUgb2YgdGhlIHNjcmVlbi4gMSA9IHN0YW5kYXJkIHJlc29sdXRpb24sIDIgPSBoaWdoIChSZXRpbmEpLCBldGMuXHJcbiAqIEBwcm9wZXJ0eSB7UG9zaXRpb259IFBvc2l0aW9uIC0gQ29udGFpbnMgdGhlIFggYW5kIFkgY29vcmRpbmF0ZXMgb2YgdGhlIHNjcmVlbidzIHBvc2l0aW9uLlxyXG4gKiBAcHJvcGVydHkge1NpemV9IFNpemUgLSBDb250YWlucyB0aGUgd2lkdGggYW5kIGhlaWdodCBvZiB0aGUgc2NyZWVuLlxyXG4gKiBAcHJvcGVydHkge1JlY3R9IEJvdW5kcyAtIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC5cclxuICogQHByb3BlcnR5IHtSZWN0fSBXb3JrQXJlYSAtIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IElzUHJpbWFyeSAtIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLlxyXG4gKiBAcHJvcGVydHkge1JvdGF0aW9ufSBSb3RhdGlvbiAtIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLlxyXG4gKi9cclxuXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLlNjcmVlbnMsICcnKTtcclxuXHJcbmNvbnN0IGdldEFsbCA9IDA7XHJcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xyXG5jb25zdCBnZXRDdXJyZW50ID0gMjtcclxuXHJcbi8qKlxyXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW5bXT59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIGFycmF5IG9mIFNjcmVlbiBvYmplY3RzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcclxuICAgIHJldHVybiBjYWxsKGdldEFsbCk7XHJcbn1cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpIHtcclxuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xyXG59XHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0Q3VycmVudCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU3lzdGVtLCAnJyk7XHJcbmNvbnN0IHN5c3RlbUlzRGFya01vZGUgPSAwO1xyXG5jb25zdCBlbnZpcm9ubWVudCA9IDE7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZykge1xyXG4gICAgaWYod2luZG93LmNocm9tZSkge1xyXG4gICAgICAgIHJldHVybiB3aW5kb3cuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UobXNnKTtcclxuICAgIH1cclxuICAgIHJldHVybiB3aW5kb3cud2Via2l0Lm1lc3NhZ2VIYW5kbGVycy5leHRlcm5hbC5wb3N0TWVzc2FnZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEBmdW5jdGlvblxyXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxib29sZWFuPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoc3lzdGVtSXNEYXJrTW9kZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBGZXRjaGVzIHRoZSBjYXBhYmlsaXRpZXMgb2YgdGhlIGFwcGxpY2F0aW9uIGZyb20gdGhlIHNlcnZlci5cclxuICpcclxuICogQGFzeW5jXHJcbiAqIEBmdW5jdGlvbiBDYXBhYmlsaXRpZXNcclxuICogQHJldHVybnMge1Byb21pc2U8T2JqZWN0Pn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKSB7XHJcbiAgICBsZXQgcmVzcG9uc2UgPSBmZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIik7XHJcbiAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG59XHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge29iamVjdH0gRW52aXJvbm1lbnRJbmZvXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBPUyAtIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IEFyY2ggLSBUaGUgYXJjaGl0ZWN0dXJlIG9mIHRoZSBzeXN0ZW0uXHJcbiAqL1xyXG5cclxuLyoqXHJcbiAqIEBmdW5jdGlvblxyXG4gKiBSZXRyaWV2ZXMgZW52aXJvbm1lbnQgZGV0YWlscy5cclxuICogQHJldHVybnMge1Byb21pc2U8RW52aXJvbm1lbnRJbmZvPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyBPUyBhbmQgc3lzdGVtIGFyY2hpdGVjdHVyZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpIHtcclxuICAgIHJldHVybiBjYWxsKGVudmlyb25tZW50KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MuXHJcbiAqXHJcbiAqIEByZXR1cm4ge2Jvb2xlYW59IFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzV2luZG93cygpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcIndpbmRvd3NcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwibGludXhcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBhIG1hY09TIG9wZXJhdGluZyBzeXN0ZW0uXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBtYWNPUywgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwiZGFyd2luXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LlxyXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FNRDY0KCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhbWQ2NFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTSgpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gLSBSZXR1cm5zIHRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZSwgb3RoZXJ3aXNlIHJldHVybnMgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBUk02NCgpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtNjRcIjtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGVidWcoKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5EZWJ1ZyA9PT0gdHJ1ZTtcclxufVxyXG5cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8vIEltcG9ydCBzY3JlZW4ganNkb2MgZGVmaW5pdGlvbiBmcm9tIC4vc2NyZWVucy5qc1xyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vc2NyZWVuc1wiKS5TY3JlZW59IFNjcmVlblxyXG4gKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNlbnRlciA9IDA7XHJcbmNvbnN0IHNldFRpdGxlID0gMTtcclxuY29uc3QgZnVsbHNjcmVlbiA9IDI7XHJcbmNvbnN0IHVuRnVsbHNjcmVlbiA9IDM7XHJcbmNvbnN0IHNldFNpemUgPSA0O1xyXG5jb25zdCBzaXplID0gNTtcclxuY29uc3Qgc2V0TWF4U2l6ZSA9IDY7XHJcbmNvbnN0IHNldE1pblNpemUgPSA3O1xyXG5jb25zdCBzZXRBbHdheXNPblRvcCA9IDg7XHJcbmNvbnN0IHNldFJlbGF0aXZlUG9zaXRpb24gPSA5O1xyXG5jb25zdCByZWxhdGl2ZVBvc2l0aW9uID0gMTA7XHJcbmNvbnN0IHNjcmVlbiA9IDExO1xyXG5jb25zdCBoaWRlID0gMTI7XHJcbmNvbnN0IG1heGltaXNlID0gMTM7XHJcbmNvbnN0IHVuTWF4aW1pc2UgPSAxNDtcclxuY29uc3QgdG9nZ2xlTWF4aW1pc2UgPSAxNTtcclxuY29uc3QgbWluaW1pc2UgPSAxNjtcclxuY29uc3QgdW5NaW5pbWlzZSA9IDE3O1xyXG5jb25zdCByZXN0b3JlID0gMTg7XHJcbmNvbnN0IHNob3cgPSAxOTtcclxuY29uc3QgY2xvc2UgPSAyMDtcclxuY29uc3Qgc2V0QmFja2dyb3VuZENvbG91ciA9IDIxO1xyXG5jb25zdCBzZXRSZXNpemFibGUgPSAyMjtcclxuY29uc3Qgd2lkdGggPSAyMztcclxuY29uc3QgaGVpZ2h0ID0gMjQ7XHJcbmNvbnN0IHpvb21JbiA9IDI1O1xyXG5jb25zdCB6b29tT3V0ID0gMjY7XHJcbmNvbnN0IHpvb21SZXNldCA9IDI3O1xyXG5jb25zdCBnZXRab29tTGV2ZWwgPSAyODtcclxuY29uc3Qgc2V0Wm9vbUxldmVsID0gMjk7XHJcblxyXG5jb25zdCB0aGlzV2luZG93ID0gR2V0KCcnKTtcclxuXHJcbmZ1bmN0aW9uIGNyZWF0ZVdpbmRvdyhjYWxsKSB7XHJcbiAgICByZXR1cm4ge1xyXG4gICAgICAgIEdldDogKHdpbmRvd05hbWUpID0+IGNyZWF0ZVdpbmRvdyhuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgd2luZG93TmFtZSkpLFxyXG4gICAgICAgIENlbnRlcjogKCkgPT4gY2FsbChjZW50ZXIpLFxyXG4gICAgICAgIFNldFRpdGxlOiAodGl0bGUpID0+IGNhbGwoc2V0VGl0bGUsIHt0aXRsZX0pLFxyXG4gICAgICAgIEZ1bGxzY3JlZW46ICgpID0+IGNhbGwoZnVsbHNjcmVlbiksXHJcbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiBjYWxsKHVuRnVsbHNjcmVlbiksXHJcbiAgICAgICAgU2V0U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcclxuICAgICAgICBTaXplOiAoKSA9PiBjYWxsKHNpemUpLFxyXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKHNldE1heFNpemUsIHt3aWR0aCwgaGVpZ2h0fSksXHJcbiAgICAgICAgU2V0TWluU2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWluU2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcclxuICAgICAgICBTZXRBbHdheXNPblRvcDogKG9uVG9wKSA9PiBjYWxsKHNldEFsd2F5c09uVG9wLCB7YWx3YXlzT25Ub3A6IG9uVG9wfSksXHJcbiAgICAgICAgU2V0UmVsYXRpdmVQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoc2V0UmVsYXRpdmVQb3NpdGlvbiwge3gsIHl9KSxcclxuICAgICAgICBSZWxhdGl2ZVBvc2l0aW9uOiAoKSA9PiBjYWxsKHJlbGF0aXZlUG9zaXRpb24pLFxyXG4gICAgICAgIFNjcmVlbjogKCkgPT4gY2FsbChzY3JlZW4pLFxyXG4gICAgICAgIEhpZGU6ICgpID0+IGNhbGwoaGlkZSksXHJcbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IGNhbGwobWF4aW1pc2UpLFxyXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IGNhbGwodW5NYXhpbWlzZSksXHJcbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IGNhbGwodG9nZ2xlTWF4aW1pc2UpLFxyXG4gICAgICAgIE1pbmltaXNlOiAoKSA9PiBjYWxsKG1pbmltaXNlKSxcclxuICAgICAgICBVbk1pbmltaXNlOiAoKSA9PiBjYWxsKHVuTWluaW1pc2UpLFxyXG4gICAgICAgIFJlc3RvcmU6ICgpID0+IGNhbGwocmVzdG9yZSksXHJcbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbChzaG93KSxcclxuICAgICAgICBDbG9zZTogKCkgPT4gY2FsbChjbG9zZSksXHJcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IGNhbGwoc2V0QmFja2dyb3VuZENvbG91ciwge3IsIGcsIGIsIGF9KSxcclxuICAgICAgICBTZXRSZXNpemFibGU6IChyZXNpemFibGUpID0+IGNhbGwoc2V0UmVzaXphYmxlLCB7cmVzaXphYmxlfSksXHJcbiAgICAgICAgV2lkdGg6ICgpID0+IGNhbGwod2lkdGgpLFxyXG4gICAgICAgIEhlaWdodDogKCkgPT4gY2FsbChoZWlnaHQpLFxyXG4gICAgICAgIFpvb21JbjogKCkgPT4gY2FsbCh6b29tSW4pLFxyXG4gICAgICAgIFpvb21PdXQ6ICgpID0+IGNhbGwoem9vbU91dCksXHJcbiAgICAgICAgWm9vbVJlc2V0OiAoKSA9PiBjYWxsKHpvb21SZXNldCksXHJcbiAgICAgICAgR2V0Wm9vbUxldmVsOiAoKSA9PiBjYWxsKGdldFpvb21MZXZlbCksXHJcbiAgICAgICAgU2V0Wm9vbUxldmVsOiAoem9vbUxldmVsKSA9PiBjYWxsKHNldFpvb21MZXZlbCwge3pvb21MZXZlbH0pLFxyXG4gICAgfTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBnZXQuXHJcbiAqIEByZXR1cm4ge09iamVjdH0gLSBUaGUgc3BlY2lmaWVkIHdpbmRvdyBvYmplY3QuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0KHdpbmRvd05hbWUpIHtcclxuICAgIHJldHVybiBjcmVhdGVXaW5kb3cobmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIHdpbmRvd05hbWUpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENlbnRlcnMgdGhlIHdpbmRvdyBvbiB0aGUgc2NyZWVuLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENlbnRlcigpIHtcclxuICAgIHRoaXNXaW5kb3cuQ2VudGVyKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSB0aXRsZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdGl0bGUgLSBUaGUgdGl0bGUgdG8gc2V0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRpdGxlKHRpdGxlKSB7XHJcbiAgICB0aGlzV2luZG93LlNldFRpdGxlKHRpdGxlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHdpbmRvdyB0byBmdWxsc2NyZWVuLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEZ1bGxzY3JlZW4oKSB7XHJcbiAgICB0aGlzV2luZG93LkZ1bGxzY3JlZW4oKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB0aGlzV2luZG93LlNldFNpemUod2lkdGgsIGhlaWdodCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2l6ZSgpIHtcclxuICAgIHJldHVybiB0aGlzV2luZG93LlNpemUoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0gd2lkdGggLSBUaGUgbWF4aW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICB0aGlzV2luZG93LlNldE1heFNpemUod2lkdGgsIGhlaWdodCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCAtIFRoZSBtaW5pbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldE1pblNpemUod2lkdGgsIGhlaWdodCkge1xyXG4gICAgdGhpc1dpbmRvdy5TZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgd2luZG93IHRvIGFsd2F5cyBiZSBvbiB0b3AuXHJcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gb25Ub3AgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGFsd2F5cyBiZSBvbiB0b3AuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0QWx3YXlzT25Ub3Aob25Ub3ApIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0QWx3YXlzT25Ub3Aob25Ub3ApO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSB3aW5kb3cncyBwb3NpdGlvbi5cclxuICogQHBhcmFtIHtudW1iZXJ9IHkgLSBUaGUgeS1jb29yZGluYXRlIG9mIHRoZSB3aW5kb3cncyBwb3NpdGlvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRSZWxhdGl2ZVBvc2l0aW9uKHgsIHkpIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0UmVsYXRpdmVQb3NpdGlvbih4LCB5KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUmVsYXRpdmVQb3NpdGlvbigpIHtcclxuICAgIHJldHVybiB0aGlzV2luZG93LlJlbGF0aXZlUG9zaXRpb24oKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2NyZWVuKCkge1xyXG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuU2NyZWVuKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIaWRlcyB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XHJcbiAgICB0aGlzV2luZG93LkhpZGUoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1heGltaXNlcyB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE1heGltaXNlKCkge1xyXG4gICAgdGhpc1dpbmRvdy5NYXhpbWlzZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVW5NYXhpbWlzZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuVW5NYXhpbWlzZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogVG9nZ2xlcyB0aGUgbWF4aW1pc2F0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVG9nZ2xlTWF4aW1pc2UoKSB7XHJcbiAgICB0aGlzV2luZG93LlRvZ2dsZU1heGltaXNlKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNaW5pbWlzZXMgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBNaW5pbWlzZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuTWluaW1pc2UoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFVuTWluaW1pc2UoKSB7XHJcbiAgICB0aGlzV2luZG93LlVuTWluaW1pc2UoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlc3RvcmVzIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUmVzdG9yZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuUmVzdG9yZSgpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xyXG4gICAgdGhpc1dpbmRvdy5TaG93KCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDbG9zZXMgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDbG9zZSgpIHtcclxuICAgIHRoaXNXaW5kb3cuQ2xvc2UoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSByIC0gVGhlIHJlZCBjb21wb25lbnQgb2YgdGhlIGNvbG91ci5cclxuICogQHBhcmFtIHtudW1iZXJ9IGcgLSBUaGUgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBiIC0gVGhlIGJsdWUgY29tcG9uZW50IG9mIHRoZSBjb2xvdXIuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBhIC0gVGhlIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgY29sb3VyLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldEJhY2tncm91bmRDb2xvdXIociwgZywgYiwgYSkge1xyXG4gICAgdGhpc1dpbmRvdy5TZXRCYWNrZ3JvdW5kQ29sb3VyKHIsIGcsIGIsIGEpO1xyXG59XHJcblxyXG4vKipcclxuICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxyXG4gKiBAcGFyYW0ge2Jvb2xlYW59IHJlc2l6YWJsZSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgcmVzaXphYmxlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFJlc2l6YWJsZShyZXNpemFibGUpIHtcclxuICAgIHRoaXNXaW5kb3cuU2V0UmVzaXphYmxlKHJlc2l6YWJsZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdpZHRoKCkge1xyXG4gICAgcmV0dXJuIHRoaXNXaW5kb3cuV2lkdGgoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhlaWdodCgpIHtcclxuICAgIHJldHVybiB0aGlzV2luZG93LkhlaWdodCgpO1xyXG59XHJcblxyXG4vKipcclxuICogWm9vbXMgaW4gdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBab29tSW4oKSB7XHJcbiAgICB0aGlzV2luZG93Llpvb21JbigpO1xyXG59XHJcblxyXG4vKipcclxuICogWm9vbXMgb3V0IHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gWm9vbU91dCgpIHtcclxuICAgIHRoaXNXaW5kb3cuWm9vbU91dCgpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVzZXRzIHRoZSB6b29tIG9mIHRoZSB3aW5kb3cuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gWm9vbVJlc2V0KCkge1xyXG4gICAgdGhpc1dpbmRvdy5ab29tUmVzZXQoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRab29tTGV2ZWwoKSB7XHJcbiAgICByZXR1cm4gdGhpc1dpbmRvdy5HZXRab29tTGV2ZWwoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHpvb21MZXZlbCAtIFRoZSB6b29tIGxldmVsIHRvIHNldC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRab29tTGV2ZWwoem9vbUxldmVsKSB7XHJcbiAgICB0aGlzV2luZG93LlNldFpvb21MZXZlbCh6b29tTGV2ZWwpO1xyXG59XHJcbiIsICJcclxuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcclxuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xyXG5pbXBvcnQge0dldH0gZnJvbSBcIi4vd2luZG93XCI7XHJcbmltcG9ydCB7T3BlblVSTH0gZnJvbSBcIi4vYnJvd3NlclwiO1xyXG5pbXBvcnQge2RlYnVnTG9nfSBmcm9tIFwiLi9sb2dcIjtcclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBvcHRpb25hbCBkYXRhLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXHJcbiAqIEBwYXJhbSB7YW55fSBbZGF0YT1udWxsXSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xyXG4gICAgbGV0IGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKTtcclxuICAgIEVtaXQoZXZlbnQpO1xyXG59XHJcblxyXG4vKipcclxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gZWxlbWVudHMgd2l0aCBgd21sLWV2ZW50YCBhdHRyaWJ1dGUuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcclxuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0nKTtcclxuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcclxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50Jyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgRGV0YWNoZWQ6IGZhbHNlLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXHJcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuXHJcbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2QgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKHdpbmRvd05hbWUsIG1ldGhvZCkge1xyXG4gICAgbGV0IHRhcmdldFdpbmRvdyA9IEdldCh3aW5kb3dOYW1lKTtcclxuICAgIGxldCBtZXRob2RNYXAgPSBXaW5kb3dNZXRob2RzKHRhcmdldFdpbmRvdyk7XHJcbiAgICBpZiAoIW1ldGhvZE1hcC5oYXMobWV0aG9kKSkge1xyXG4gICAgICAgIGNvbnNvbGUubG9nKFwiV2luZG93IG1ldGhvZCBcIiArIG1ldGhvZCArIFwiIG5vdCBmb3VuZFwiKTtcclxuICAgIH1cclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWV0aG9kTWFwLmdldChtZXRob2QpKCk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIGNhbGxpbmcgd2luZG93IG1ldGhvZCAnXCIgKyBtZXRob2QgKyBcIic6IFwiICsgZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBZGRzIHdpbmRvdyBsaXN0ZW5lcnMgZm9yIGVsZW1lbnRzIHdpdGggdGhlICd3bWwtd2luZG93JyBhdHRyaWJ1dGUuXHJcbiAqIFJlbW92ZXMgYW55IGV4aXN0aW5nIGxpc3RlbmVycyBiZWZvcmUgYWRkaW5nIG5ldyBvbmVzLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZnVuY3Rpb24gYWRkV01MV2luZG93TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLXdpbmRvd10nKTtcclxuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcclxuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpO1xyXG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcclxuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgJ2NsaWNrJztcclxuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCAnJztcclxuXHJcbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xyXG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xyXG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHRhcmdldFdpbmRvdywgd2luZG93TWV0aG9kKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHRhcmdldFdpbmRvdywgd2luZG93TWV0aG9kKTtcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXHJcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuXHJcbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEFkZHMgYSBsaXN0ZW5lciB0byBlbGVtZW50cyB3aXRoIHRoZSAnd21sLW9wZW51cmwnIGF0dHJpYnV0ZS5cclxuICogV2hlbiB0aGUgc3BlY2lmaWVkIHRyaWdnZXIgZXZlbnQgaXMgZmlyZWQgb24gYW55IG9mIHRoZXNlIGVsZW1lbnRzLFxyXG4gKiB0aGUgbGlzdGVuZXIgd2lsbCBvcGVuIHRoZSBVUkwgc3BlY2lmaWVkIGJ5IHRoZSAnd21sLW9wZW51cmwnIGF0dHJpYnV0ZS5cclxuICogSWYgYSAnd21sLWNvbmZpcm0nIGF0dHJpYnV0ZSBpcyBwcm92aWRlZCwgYSBjb25maXJtYXRpb24gZGlhbG9nIHdpbGwgYmUgZGlzcGxheWVkLFxyXG4gKiBhbmQgdGhlIFVSTCB3aWxsIG9ubHkgYmUgb3BlbmVkIGlmIHRoZSB1c2VyIGNvbmZpcm1zLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZnVuY3Rpb24gYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpIHtcclxuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1vcGVudXJsXScpO1xyXG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xyXG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtb3BlbnVybCcpO1xyXG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcclxuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xyXG5cclxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XHJcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgT3BlblVSTCh1cmwpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0pO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIHZvaWQgT3BlblVSTCh1cmwpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpIHtcclxuICAgIGlmKERFQlVHKSB7XHJcbiAgICAgICAgZGVidWdMb2coXCJSZWxvYWRpbmcgV01MXCIpO1xyXG4gICAgfVxyXG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcclxuICAgIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpO1xyXG4gICAgYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIG1hcCBvZiBhbGwgbWV0aG9kcyBpbiB0aGUgY3VycmVudCB3aW5kb3cuXHJcbiAqIEByZXR1cm5zIHtNYXB9IC0gQSBtYXAgb2Ygd2luZG93IG1ldGhvZHMuXHJcbiAqL1xyXG5mdW5jdGlvbiBXaW5kb3dNZXRob2RzKHRhcmdldFdpbmRvdykge1xyXG4gICAgLy8gQ3JlYXRlIGEgbmV3IG1hcCB0byBzdG9yZSBtZXRob2RzXHJcbiAgICBsZXQgcmVzdWx0ID0gbmV3IE1hcCgpO1xyXG5cclxuICAgIC8vIEl0ZXJhdGUgb3ZlciBhbGwgcHJvcGVydGllcyBvZiB0aGUgd2luZG93IG9iamVjdFxyXG4gICAgZm9yIChsZXQgbWV0aG9kIGluIHRhcmdldFdpbmRvdykge1xyXG4gICAgICAgIC8vIENoZWNrIGlmIHRoZSBwcm9wZXJ0eSBpcyBpbmRlZWQgYSBtZXRob2QgKGZ1bmN0aW9uKVxyXG4gICAgICAgIGlmKHR5cGVvZiB0YXJnZXRXaW5kb3dbbWV0aG9kXSA9PT0gJ2Z1bmN0aW9uJykge1xyXG4gICAgICAgICAgICAvLyBBZGQgdGhlIG1ldGhvZCB0byB0aGUgbWFwXHJcbiAgICAgICAgICAgIHJlc3VsdC5zZXQobWV0aG9kLCB0YXJnZXRXaW5kb3dbbWV0aG9kXSk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgIH1cclxuICAgIC8vIFJldHVybiB0aGUgbWFwIG9mIHdpbmRvdyBtZXRob2RzXHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi90eXBlc1wiKS5XYWlsc0V2ZW50fSBXYWlsc0V2ZW50XHJcbiAqL1xyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5pbXBvcnQge0V2ZW50VHlwZXN9IGZyb20gXCIuL2V2ZW50X3R5cGVzXCI7XHJcbmV4cG9ydCBjb25zdCBUeXBlcyA9IEV2ZW50VHlwZXM7XHJcblxyXG4vLyBTZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5FdmVudHMsICcnKTtcclxuY29uc3QgRW1pdE1ldGhvZCA9IDA7XHJcbmNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xyXG5cclxuY2xhc3MgTGlzdGVuZXIge1xyXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XHJcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XHJcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XHJcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XHJcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xyXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XHJcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcclxuICAgICAgICB9O1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhID0gbnVsbCkge1xyXG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XHJcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHNldHVwKCkge1xyXG59XHJcblxyXG5mdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XHJcbiAgICBpZiAobGlzdGVuZXJzKSB7XHJcbiAgICAgICAgbGV0IHRvUmVtb3ZlID0gbGlzdGVuZXJzLmZpbHRlcihsaXN0ZW5lciA9PiB7XHJcbiAgICAgICAgICAgIGxldCByZW1vdmUgPSBsaXN0ZW5lci5DYWxsYmFjayhldmVudCk7XHJcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHJldHVybiB0cnVlO1xyXG4gICAgICAgIH0pO1xyXG4gICAgICAgIGlmICh0b1JlbW92ZS5sZW5ndGggPiAwKSB7XHJcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xyXG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xyXG4gICAgICAgICAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCBsaXN0ZW5lcnMpO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3MgLSBUaGUgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdGhlIGNhbGxiYWNrIGNhbiBiZSBjYWxsZWQgZm9yIHRoZSBldmVudC4gT25jZSB0aGUgbWF4aW11bSBudW1iZXIgaXMgcmVhY2hlZCwgdGhlIGNhbGxiYWNrIHdpbGwgbm8gbG9uZ2VyIGJlIGNhbGxlZC5cclxuICpcclxuIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcclxuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIHNwZWNpZmllZCBldmVudCBvY2N1cnMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkLiBJdCB0YWtlcyBubyBwYXJhbWV0ZXJzLlxyXG4gKiBAcmV0dXJuIHtmdW5jdGlvbn0gLSBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lLCBjYWxsYmFjaykgeyByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7IH1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXHJcbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIGV2ZW50IG9jY3Vycy5cclxuICogQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHsgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7IH1cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIHRoZSBzcGVjaWZpZWQgbGlzdGVuZXIgZnJvbSB0aGUgZXZlbnQgbGlzdGVuZXJzIGNvbGxlY3Rpb24uXHJcbiAqIElmIGFsbCBsaXN0ZW5lcnMgZm9yIHRoZSBldmVudCBhcmUgcmVtb3ZlZCwgdGhlIGV2ZW50IGtleSBpcyBkZWxldGVkIGZyb20gdGhlIGNvbGxlY3Rpb24uXHJcbiAqXHJcbiAqIEBwYXJhbSB7T2JqZWN0fSBsaXN0ZW5lciAtIFRoZSBsaXN0ZW5lciB0byBiZSByZW1vdmVkLlxyXG4gKi9cclxuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcclxuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XHJcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XHJcbiAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXHJcbiAqIEBwYXJhbSB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lcyAtIEFkZGl0aW9uYWwgZXZlbnQgbmFtZXMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXHJcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xyXG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xyXG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xyXG59XHJcbi8qKlxyXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXHJcbiAqXHJcbiAqIEBmdW5jdGlvbiBPZmZBbGxcclxuICogQHJldHVybnMge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkgeyBldmVudExpc3RlbmVycy5jbGVhcigpOyB9XHJcblxyXG4vKipcclxuICogRW1pdHMgYW4gZXZlbnQgdXNpbmcgdGhlIGdpdmVuIGV2ZW50IG5hbWUuXHJcbiAqXHJcbiAqIEBwYXJhbSB7V2FpbHNFdmVudH0gZXZlbnQgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdC5cclxuICogQHJldHVybnMge2FueX0gLSBUaGUgcmVzdWx0IG9mIHRoZSBlbWl0dGVkIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQoZXZlbnQpIHsgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgZXZlbnQpOyB9XHJcbiIsICJcbmV4cG9ydCBjb25zdCBFdmVudFR5cGVzID0ge1xuXHRXaW5kb3dzOiB7XG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0QVBNUG93ZXJTdGF0dXNDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclN0YXR1c0NoYW5nZVwiLFxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwid2luZG93czpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZDogXCJ3aW5kb3dzOldlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkXCIsXG5cdFx0V2luZG93SW5hY3RpdmU6IFwid2luZG93czpXaW5kb3dJbmFjdGl2ZVwiLFxuXHRcdFdpbmRvd0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Q2xvc2U6IFwid2luZG93czpXaW5kb3dDbG9zZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHR9LFxuXHRNYWM6IHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZU5vdGlmaWNhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVOb3RpZmljYXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZSFcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW4hXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd0RpZERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRFbmRTaGVldDogXCJtYWM6V2luZG93RGlkRW5kU2hlZXRcIixcblx0XHRXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4cG9zZTogXCJtYWM6V2luZG93RGlkRXhwb3NlXCIsXG5cdFx0V2luZG93RGlkRm9jdXM6IFwibWFjOldpbmRvd0RpZEZvY3VzXCIsXG5cdFx0V2luZG93RGlkTWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZE1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJtYWM6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT25TY3JlZW5cIixcblx0XHRXaW5kb3dEaWRSZXNpZ25LZXk6IFwibWFjOldpbmRvd0RpZFJlc2lnbktleVwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbk1haW46IFwibWFjOldpbmRvd0RpZFJlc2lnbk1haW5cIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibWFjOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZTogXCJtYWM6V2luZG93RGlkVXBkYXRlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUFscGhhXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVNoYWRvd1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVWaXNpYmlsaXR5XCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlIVwiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdE1lbnVXaWxsT3BlbjogXCJtYWM6TWVudVdpbGxPcGVuXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWREaXNwbGF5SXRlbTogXCJtYWM6TWVudURpZERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVEaWRSZW1vdmVJdGVtOiBcIm1hYzpNZW51RGlkUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsQmVnaW5UcmFja2luZzogXCJtYWM6TWVudVdpbGxCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51RGlkUG9wVXA6IFwibWFjOk1lbnVEaWRQb3BVcFwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0fSxcblx0TGludXg6IHtcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdH0sXG5cdENvbW1vbjoge1xuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwiY29tbW9uOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJjb21tb246V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJjb21tb246V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd1pvb206IFwiY29tbW9uOldpbmRvd1pvb21cIixcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwiY29tbW9uOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwiY29tbW9uOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJjb21tb246V2luZG93TG9zdEZvY3VzXCIsXG5cdFx0V2luZG93U2hvdzogXCJjb21tb246V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd0hpZGU6IFwiY29tbW9uOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcImNvbW1vbjpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcblx0XHRXaW5kb3dSdW50aW1lUmVhZHk6IFwiY29tbW9uOldpbmRvd1J1bnRpbWVSZWFkeVwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXG5cdH0sXG59O1xuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBPcGVuRmlsZURpYWxvZ09wdGlvbnNcclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DcmVhdGVEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Nob3dIaWRkZW5GaWxlc10gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb25dIC0gSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtIaWRlRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtBbGxvd3NPdGhlckZpbGV0eXBlc10gLSBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC5cclxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaXRsZSBvZiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtEaXJlY3RvcnldIC0gRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy5cclxuICovXHJcblxyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNhdmVGaWxlRGlhbG9nT3B0aW9uc1xyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0ZpbGVuYW1lXSAtIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VGaWxlc10gLSBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ3JlYXRlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtTaG93SGlkZGVuRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Jlc29sdmVzQWxpYXNlc10gLSBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c011bHRpcGxlU2VsZWN0aW9uXSAtIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbSGlkZUV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzT3RoZXJGaWxldHlwZXNdIC0gSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7RmlsZUZpbHRlcltdfSBbRmlsdGVyc10gLSBBcnJheSBvZiBmaWxlIGZpbHRlcnMuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGl0bGUgb2YgdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtNZXNzYWdlXSAtIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0J1dHRvblRleHRdIC0gVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlyZWN0b3J5XSAtIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0RldGFjaGVkXSAtIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuXHJcbiAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IE1lc3NhZ2VEaWFsb2dPcHRpb25zXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGhlIHRpdGxlIG9mIHRoZSBkaWFsb2cgd2luZG93LlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7QnV0dG9uW119IFtCdXR0b25zXSAtIEFycmF5IG9mIGJ1dHRvbiBvcHRpb25zIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBCdXR0b25cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtMYWJlbF0gLSBUZXh0IHRoYXQgYXBwZWFycyB3aXRoaW4gdGhlIGJ1dHRvbi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNDYW5jZWxdIC0gVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNEZWZhdWx0XSAtIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci5cclxuICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gRmlsZUZpbHRlclxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0Rpc3BsYXlOYW1lXSAtIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1BhdHRlcm5dIC0gUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy5cclxuICovXHJcblxyXG4vLyBzZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5kaWFsb2dFcnJvckNhbGxiYWNrID0gZGlhbG9nRXJyb3JDYWxsYmFjaztcclxud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xyXG5cclxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcclxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XHJcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xyXG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XHJcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcclxuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xyXG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5EaWFsb2csICcnKTtcclxuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xyXG5cclxuLyoqXHJcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBpZCB0aGF0IGlzIG5vdCBwcmVzZW50IGluIGRpYWxvZ1Jlc3BvbnNlcy5cclxuICogQHJldHVybnMge3N0cmluZ30gdW5pcXVlIGlkXHJcbiAqL1xyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICogQHBhcmFtIHtudW1iZXJ9IHR5cGUgLSB0eXBlIG9mIGRpYWxvZ1xyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfE9wZW5GaWxlRGlhbG9nT3B0aW9uc3xTYXZlRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBvcHRpb25zIGZvciB0aGUgZGlhbG9nXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlfSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nXHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucyA9IHt9KSB7XHJcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuICAgIG9wdGlvbnNbXCJkaWFsb2ctaWRcIl0gPSBpZDtcclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIHRoZSBjYWxsYmFjayBmcm9tIGEgZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgSUQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cclxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgZGF0YSByZWNlaXZlZCBmcm9tIHRoZSBkaWFsb2cuXHJcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gaXNKU09OIC0gRmxhZyBpbmRpY2F0aW5nIHdoZXRoZXIgdGhlIGRhdGEgaXMgaW4gSlNPTiBmb3JtYXQuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZ1Jlc3VsdENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsYmFjayBmdW5jdGlvbiBmb3IgaGFuZGxpbmcgZXJyb3JzIGluIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSBkaWFsb2cgcmVzcG9uc2UuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XHJcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgaWYgKHApIHtcclxuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5cclxuXHJcbi8vIFJlcGxhY2UgYG1ldGhvZHNgIHdpdGggY29uc3RhbnRzIGluIFRpdGxlIENhc2VcclxuXHJcbi8qKlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgSW5mbyA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7XHJcblxyXG4vKipcclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IFdhcm5pbmcgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBFcnJvciA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBRdWVzdGlvbiA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7T3BlbkZpbGVEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nW118c3RyaW5nPn0gUmV0dXJucyBzZWxlY3RlZCBmaWxlIG9yIGxpc3Qgb2YgZmlsZXMuIFJldHVybnMgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgT3BlbkZpbGUgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKTtcclxuXHJcbi8qKlxyXG4gKiBAcGFyYW0ge1NhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUuIFJldHVybnMgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgU2F2ZUZpbGUgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTtcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcclxud2luZG93Ll93YWlscy5jYWxsRXJyb3JIYW5kbGVyID0gZXJyb3JIYW5kbGVyO1xyXG5cclxuXHJcbmNvbnN0IENhbGxCaW5kaW5nID0gMDtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2FsbCwgJycpO1xyXG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcclxuXHJcbi8qKlxyXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtzdHJpbmd9IC0gQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgY2FsbFJlc3BvbnNlcyBzZXQuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIHJlc3VsdCBvZiBhIGNhbGwgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cclxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgcmVzdWx0IGRhdGEgb2YgdGhlIHJlcXVlc3QuXHJcbiAqIEBwYXJhbSB7Ym9vbGVhbn0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH0gLSBUaGlzIG1ldGhvZCBkb2VzIG5vdCByZXR1cm4gYW55IHZhbHVlLlxyXG4gKi9cclxuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZCwgZGF0YSwgaXNKU09OKSB7XHJcbiAgICBjb25zdCBwcm9taXNlSGFuZGxlciA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcclxuICAgIGlmIChwcm9taXNlSGFuZGxlcikge1xyXG4gICAgICAgIHByb21pc2VIYW5kbGVyLnJlc29sdmUoaXNKU09OID8gSlNPTi5wYXJzZShkYXRhKSA6IGRhdGEpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogSGFuZGxlcyB0aGUgZXJyb3IgZnJvbSBhIGNhbGwgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSBwcm9taXNlIGhhbmRsZXIuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UgdG8gcmVqZWN0IHRoZSBwcm9taXNlIGhhbmRsZXIgd2l0aC5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZCwgbWVzc2FnZSkge1xyXG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XHJcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcclxuICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobWVzc2FnZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgY2FsbFJlc3BvbnNlcyBtYXAuXHJcbiAqXHJcbiAqIEBwYXJhbSB7YW55fSBpZCAtIFRoZSBJRCBvZiB0aGUgcmVzcG9uc2UgdG8gYmUgcmV0cmlldmVkIGFuZCByZW1vdmVkLlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7YW55fSBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZXRBbmREZWxldGVSZXNwb25zZShpZCkge1xyXG4gICAgY29uc3QgcmVzcG9uc2UgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICByZXR1cm4gcmVzcG9uc2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFeGVjdXRlcyBhIGNhbGwgdXNpbmcgdGhlIHByb3ZpZGVkIHR5cGUgYW5kIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfG51bWJlcn0gdHlwZSAtIFRoZSB0eXBlIG9mIGNhbGwgdG8gZXhlY3V0ZS5cclxuICogQHBhcmFtIHtPYmplY3R9IFtvcHRpb25zPXt9XSAtIEFkZGl0aW9uYWwgb3B0aW9ucyBmb3IgdGhlIGNhbGwuXHJcbiAqIEByZXR1cm4ge1Byb21pc2V9IC0gQSBwcm9taXNlIHRoYXQgd2lsbCBiZSByZXNvbHZlZCBvciByZWplY3RlZCBiYXNlZCBvbiB0aGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucyA9IHt9KSB7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbCBtZXRob2QuXHJcbiAqXHJcbiAqIEBwYXJhbSB7T2JqZWN0fSBvcHRpb25zIC0gVGhlIG9wdGlvbnMgZm9yIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm5zIHtPYmplY3R9IC0gVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFeGVjdXRlcyBhIG1ldGhvZCBieSBuYW1lLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJy5cclxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEB0aHJvd3Mge0Vycm9yfSBJZiB0aGUgbmFtZSBpcyBub3QgYSBzdHJpbmcgb3IgaXMgbm90IGluIHRoZSBjb3JyZWN0IGZvcm1hdC5cclxuICogQHJldHVybnMgeyp9IFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBleGVjdXRpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG5hbWUsIC4uLmFyZ3MpIHtcclxuICAgIGlmICh0eXBlb2YgbmFtZSAhPT0gXCJzdHJpbmdcIiB8fCBuYW1lLnNwbGl0KFwiLlwiKS5sZW5ndGggIT09IDMpIHtcclxuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJDYWxsQnlOYW1lIHJlcXVpcmVzIGEgc3RyaW5nIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCdcIik7XHJcbiAgICB9XHJcbiAgICBsZXQgW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXSA9IG5hbWUuc3BsaXQoXCIuXCIpO1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWUsXHJcbiAgICAgICAgc3RydWN0TmFtZSxcclxuICAgICAgICBtZXRob2ROYW1lLFxyXG4gICAgICAgIGFyZ3NcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7bnVtYmVyfSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQsIC4uLmFyZ3MpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xyXG4gICAgICAgIG1ldGhvZElELFxyXG4gICAgICAgIGFyZ3NcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBwbHVnaW4uXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHBsdWdpbi5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJucyB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxyXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXHJcbiAgICAgICAgbWV0aG9kTmFtZSxcclxuICAgICAgICBhcmdzXHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtkZWJ1Z0xvZ30gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2xvZ1wiO1xyXG5cclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvYXBwbGljYXRpb25cIjtcclxuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvYnJvd3NlclwiO1xyXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZFwiO1xyXG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3NcIjtcclxuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvc2NyZWVuc1wiO1xyXG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbVwiO1xyXG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvd1wiO1xyXG5pbXBvcnQgKiBhcyBXTUwgZnJvbSAnLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sJztcclxuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudHNcIjtcclxuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9nc1wiO1xyXG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxsc1wiO1xyXG5pbXBvcnQge2ludm9rZX0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbVwiO1xyXG5cclxuLyoqKlxyXG4gVGhpcyB0ZWNobmlxdWUgZm9yIHByb3BlciBsb2FkIGRldGVjdGlvbiBpcyB0YWtlbiBmcm9tIEhUTVg6XHJcblxyXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcclxuXHJcbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcclxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXHJcblxyXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XHJcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcclxuXHJcbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcclxuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxyXG5cclxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcclxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cclxuIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLlxyXG5cclxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXHJcbiBBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFXHJcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcclxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcclxuIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMXHJcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxyXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcclxuIENBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksXHJcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxyXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cclxuXHJcbiAqKiovXHJcblxyXG53aW5kb3cuX3dhaWxzLmludm9rZT1pbnZva2U7XHJcblxyXG53aW5kb3cud2FpbHMgPSB3aW5kb3cud2FpbHMgfHwge307XHJcbndpbmRvdy53YWlscy5BcHBsaWNhdGlvbiA9IEFwcGxpY2F0aW9uO1xyXG53aW5kb3cud2FpbHMuQnJvd3NlciA9IEJyb3dzZXI7XHJcbndpbmRvdy53YWlscy5DYWxsID0gQ2FsbDtcclxud2luZG93LndhaWxzLkNsaXBib2FyZCA9IENsaXBib2FyZDtcclxud2luZG93LndhaWxzLkRpYWxvZ3MgPSBEaWFsb2dzO1xyXG53aW5kb3cud2FpbHMuRXZlbnRzID0gRXZlbnRzO1xyXG53aW5kb3cud2FpbHMuRmxhZ3MgPSBGbGFncztcclxud2luZG93LndhaWxzLlNjcmVlbnMgPSBTY3JlZW5zO1xyXG53aW5kb3cud2FpbHMuU3lzdGVtID0gU3lzdGVtO1xyXG53aW5kb3cud2FpbHMuV2luZG93ID0gV2luZG93O1xyXG53aW5kb3cud2FpbHMuV01MID0gV01MO1xyXG5cclxuXHJcbmxldCBpc1JlYWR5ID0gZmFsc2VcclxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGZ1bmN0aW9uKCkge1xyXG4gICAgaXNSZWFkeSA9IHRydWVcclxuICAgIHdpbmRvdy5fd2FpbHMuaW52b2tlKCd3YWlsczpydW50aW1lOnJlYWR5Jyk7XHJcbiAgICBpZihERUJVRykge1xyXG4gICAgICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIik7XHJcbiAgICB9XHJcbn0pXHJcblxyXG5mdW5jdGlvbiB3aGVuUmVhZHkoZm4pIHtcclxuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcclxuICAgICAgICBmbigpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgZm4pO1xyXG4gICAgfVxyXG59XHJcblxyXG53aGVuUmVhZHkoKCkgPT4ge1xyXG4gICAgV01MLlJlbG9hZCgpO1xyXG59KTtcclxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7QUFLTyxXQUFTLFNBQVMsU0FBUztBQUU5QixZQUFRO0FBQUEsTUFDSixrQkFBa0IsVUFBVTtBQUFBLE1BQzVCO0FBQUEsTUFDQTtBQUFBLElBQ0o7QUFBQSxFQUNKOzs7QUNaQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBLE1BQUksY0FDRjtBQVdLLE1BQUksU0FBUyxDQUFDQSxRQUFPLE9BQU87QUFDakMsUUFBSSxLQUFLO0FBQ1QsUUFBSSxJQUFJQTtBQUNSLFdBQU8sS0FBSztBQUNWLFlBQU0sWUFBYSxLQUFLLE9BQU8sSUFBSSxLQUFNLENBQUM7QUFBQSxJQUM1QztBQUNBLFdBQU87QUFBQSxFQUNUOzs7QUNOQSxNQUFNLGFBQWEsT0FBTyxTQUFTLFNBQVM7QUFHckMsTUFBTSxjQUFjO0FBQUEsSUFDdkIsTUFBTTtBQUFBLElBQ04sV0FBVztBQUFBLElBQ1gsYUFBYTtBQUFBLElBQ2IsUUFBUTtBQUFBLElBQ1IsYUFBYTtBQUFBLElBQ2IsUUFBUTtBQUFBLElBQ1IsUUFBUTtBQUFBLElBQ1IsU0FBUztBQUFBLElBQ1QsUUFBUTtBQUFBLElBQ1IsU0FBUztBQUFBLEVBQ2I7QUFDTyxNQUFJLFdBQVcsT0FBTztBQXNCdEIsV0FBUyx1QkFBdUIsUUFBUSxZQUFZO0FBQ3ZELFdBQU8sU0FBVSxRQUFRLE9BQUssTUFBTTtBQUNoQyxhQUFPLGtCQUFrQixRQUFRLFFBQVEsWUFBWSxJQUFJO0FBQUEsSUFDN0Q7QUFBQSxFQUNKO0FBcUNBLFdBQVMsa0JBQWtCLFVBQVUsUUFBUSxZQUFZLE1BQU07QUFDM0QsUUFBSSxNQUFNLElBQUksSUFBSSxVQUFVO0FBQzVCLFFBQUksYUFBYSxPQUFPLFVBQVUsUUFBUTtBQUMxQyxRQUFJLGFBQWEsT0FBTyxVQUFVLE1BQU07QUFDeEMsUUFBSSxlQUFlO0FBQUEsTUFDZixTQUFTLENBQUM7QUFBQSxJQUNkO0FBQ0EsUUFBSSxZQUFZO0FBQ1osbUJBQWEsUUFBUSxxQkFBcUIsSUFBSTtBQUFBLElBQ2xEO0FBQ0EsUUFBSSxNQUFNO0FBQ04sVUFBSSxhQUFhLE9BQU8sUUFBUSxLQUFLLFVBQVUsSUFBSSxDQUFDO0FBQUEsSUFDeEQ7QUFDQSxpQkFBYSxRQUFRLG1CQUFtQixJQUFJO0FBQzVDLFdBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLFlBQU0sS0FBSyxZQUFZLEVBQ2xCLEtBQUssY0FBWTtBQUNkLFlBQUksU0FBUyxJQUFJO0FBRWIsY0FBSSxTQUFTLFFBQVEsSUFBSSxjQUFjLEtBQUssU0FBUyxRQUFRLElBQUksY0FBYyxFQUFFLFFBQVEsa0JBQWtCLE1BQU0sSUFBSTtBQUNqSCxtQkFBTyxTQUFTLEtBQUs7QUFBQSxVQUN6QixPQUFPO0FBQ0gsbUJBQU8sU0FBUyxLQUFLO0FBQUEsVUFDekI7QUFBQSxRQUNKO0FBQ0EsZUFBTyxNQUFNLFNBQVMsVUFBVSxDQUFDO0FBQUEsTUFDckMsQ0FBQyxFQUNBLEtBQUssVUFBUSxRQUFRLElBQUksQ0FBQyxFQUMxQixNQUFNLFdBQVMsT0FBTyxLQUFLLENBQUM7QUFBQSxJQUNyQyxDQUFDO0FBQUEsRUFDTDs7O0FGNUdBLE1BQU0sT0FBTyx1QkFBdUIsWUFBWSxhQUFhLEVBQUU7QUFFL0QsTUFBTSxhQUFhO0FBQ25CLE1BQU0sYUFBYTtBQUNuQixNQUFNLGFBQWE7QUFRWixXQUFTLE9BQU87QUFDbkIsV0FBTyxLQUFLLFVBQVU7QUFBQSxFQUMxQjtBQU9PLFdBQVMsT0FBTztBQUNuQixXQUFPLEtBQUssVUFBVTtBQUFBLEVBQzFCO0FBT08sV0FBUyxPQUFPO0FBQ25CLFdBQU8sS0FBSyxVQUFVO0FBQUEsRUFDMUI7OztBRzdDQTtBQUFBO0FBQUE7QUFBQTtBQWFBLE1BQU1DLFFBQU8sdUJBQXVCLFlBQVksU0FBUyxFQUFFO0FBQzNELE1BQU0saUJBQWlCO0FBT2hCLFdBQVMsUUFBUSxLQUFLO0FBQ3pCLFdBQU9BLE1BQUssZ0JBQWdCLEVBQUMsSUFBRyxDQUFDO0FBQUEsRUFDckM7OztBQ3ZCQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBY0EsTUFBTUMsUUFBTyx1QkFBdUIsWUFBWSxXQUFXLEVBQUU7QUFDN0QsTUFBTSxtQkFBbUI7QUFDekIsTUFBTSxnQkFBZ0I7QUFRZixXQUFTLFFBQVEsTUFBTTtBQUMxQixXQUFPQSxNQUFLLGtCQUFrQixFQUFDLEtBQUksQ0FBQztBQUFBLEVBQ3hDO0FBTU8sV0FBUyxPQUFPO0FBQ25CLFdBQU9BLE1BQUssYUFBYTtBQUFBLEVBQzdCOzs7QUNsQ0E7QUFBQTtBQUFBO0FBQUE7QUFrQk8sV0FBUyxRQUFRLFdBQVc7QUFDL0IsUUFBSTtBQUNBLGFBQU8sT0FBTyxPQUFPLE1BQU0sU0FBUztBQUFBLElBQ3hDLFNBQVMsR0FBRztBQUNSLFlBQU0sSUFBSSxNQUFNLDhCQUE4QixZQUFZLFFBQVEsQ0FBQztBQUFBLElBQ3ZFO0FBQUEsRUFDSjs7O0FDeEJBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQXVEQSxNQUFNQyxRQUFPLHVCQUF1QixZQUFZLFNBQVMsRUFBRTtBQUUzRCxNQUFNLFNBQVM7QUFDZixNQUFNLGFBQWE7QUFDbkIsTUFBTSxhQUFhO0FBTVosV0FBUyxTQUFTO0FBQ3JCLFdBQU9BLE1BQUssTUFBTTtBQUFBLEVBQ3RCO0FBS08sV0FBUyxhQUFhO0FBQ3pCLFdBQU9BLE1BQUssVUFBVTtBQUFBLEVBQzFCO0FBTU8sV0FBUyxhQUFhO0FBQ3pCLFdBQU9BLE1BQUssVUFBVTtBQUFBLEVBQzFCOzs7QUNsRkE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWFBLE1BQUlDLFFBQU8sdUJBQXVCLFlBQVksUUFBUSxFQUFFO0FBQ3hELE1BQU0sbUJBQW1CO0FBQ3pCLE1BQU0sY0FBYztBQUViLFdBQVMsT0FBTyxLQUFLO0FBQ3hCLFFBQUcsT0FBTyxRQUFRO0FBQ2QsYUFBTyxPQUFPLE9BQU8sUUFBUSxZQUFZLEdBQUc7QUFBQSxJQUNoRDtBQUNBLFdBQU8sT0FBTyxPQUFPLGdCQUFnQixTQUFTO0FBQUEsRUFDbEQ7QUFPTyxXQUFTLGFBQWE7QUFDekIsV0FBT0EsTUFBSyxnQkFBZ0I7QUFBQSxFQUNoQztBQVNPLFdBQVMsZUFBZTtBQUMzQixRQUFJLFdBQVcsTUFBTSxxQkFBcUI7QUFDMUMsV0FBTyxTQUFTLEtBQUs7QUFBQSxFQUN6QjtBQWFPLFdBQVMsY0FBYztBQUMxQixXQUFPQSxNQUFLLFdBQVc7QUFBQSxFQUMzQjtBQU9PLFdBQVMsWUFBWTtBQUN4QixXQUFPLE9BQU8sT0FBTyxZQUFZLE9BQU87QUFBQSxFQUM1QztBQU9PLFdBQVMsVUFBVTtBQUN0QixXQUFPLE9BQU8sT0FBTyxZQUFZLE9BQU87QUFBQSxFQUM1QztBQU9PLFdBQVMsUUFBUTtBQUNwQixXQUFPLE9BQU8sT0FBTyxZQUFZLE9BQU87QUFBQSxFQUM1QztBQU1PLFdBQVMsVUFBVTtBQUN0QixXQUFPLE9BQU8sT0FBTyxZQUFZLFNBQVM7QUFBQSxFQUM5QztBQU9PLFdBQVMsUUFBUTtBQUNwQixXQUFPLE9BQU8sT0FBTyxZQUFZLFNBQVM7QUFBQSxFQUM5QztBQU9PLFdBQVMsVUFBVTtBQUN0QixXQUFPLE9BQU8sT0FBTyxZQUFZLFNBQVM7QUFBQSxFQUM5QztBQUVPLFdBQVMsVUFBVTtBQUN0QixXQUFPLE9BQU8sT0FBTyxZQUFZLFVBQVU7QUFBQSxFQUMvQzs7O0FDbkhBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxnQkFBQUM7QUFBQSxJQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxnQkFBQUM7QUFBQSxJQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQW1CQSxNQUFNLFNBQVM7QUFDZixNQUFNLFdBQVc7QUFDakIsTUFBTSxhQUFhO0FBQ25CLE1BQU0sZUFBZTtBQUNyQixNQUFNLFVBQVU7QUFDaEIsTUFBTSxPQUFPO0FBQ2IsTUFBTSxhQUFhO0FBQ25CLE1BQU0sYUFBYTtBQUNuQixNQUFNLGlCQUFpQjtBQUN2QixNQUFNLHNCQUFzQjtBQUM1QixNQUFNLG1CQUFtQjtBQUN6QixNQUFNLFNBQVM7QUFDZixNQUFNLE9BQU87QUFDYixNQUFNLFdBQVc7QUFDakIsTUFBTSxhQUFhO0FBQ25CLE1BQU0saUJBQWlCO0FBQ3ZCLE1BQU0sV0FBVztBQUNqQixNQUFNLGFBQWE7QUFDbkIsTUFBTSxVQUFVO0FBQ2hCLE1BQU0sT0FBTztBQUNiLE1BQU0sUUFBUTtBQUNkLE1BQU0sc0JBQXNCO0FBQzVCLE1BQU0sZUFBZTtBQUNyQixNQUFNLFFBQVE7QUFDZCxNQUFNLFNBQVM7QUFDZixNQUFNLFNBQVM7QUFDZixNQUFNLFVBQVU7QUFDaEIsTUFBTSxZQUFZO0FBQ2xCLE1BQU0sZUFBZTtBQUNyQixNQUFNLGVBQWU7QUFFckIsTUFBTSxhQUFhLElBQUksRUFBRTtBQUV6QixXQUFTLGFBQWFDLE9BQU07QUFDeEIsV0FBTztBQUFBLE1BQ0gsS0FBSyxDQUFDLGVBQWUsYUFBYSx1QkFBdUIsWUFBWSxRQUFRLFVBQVUsQ0FBQztBQUFBLE1BQ3hGLFFBQVEsTUFBTUEsTUFBSyxNQUFNO0FBQUEsTUFDekIsVUFBVSxDQUFDLFVBQVVBLE1BQUssVUFBVSxFQUFDLE1BQUssQ0FBQztBQUFBLE1BQzNDLFlBQVksTUFBTUEsTUFBSyxVQUFVO0FBQUEsTUFDakMsY0FBYyxNQUFNQSxNQUFLLFlBQVk7QUFBQSxNQUNyQyxTQUFTLENBQUNDLFFBQU9DLFlBQVdGLE1BQUssU0FBUyxFQUFDLE9BQUFDLFFBQU8sUUFBQUMsUUFBTSxDQUFDO0FBQUEsTUFDekQsTUFBTSxNQUFNRixNQUFLLElBQUk7QUFBQSxNQUNyQixZQUFZLENBQUNDLFFBQU9DLFlBQVdGLE1BQUssWUFBWSxFQUFDLE9BQUFDLFFBQU8sUUFBQUMsUUFBTSxDQUFDO0FBQUEsTUFDL0QsWUFBWSxDQUFDRCxRQUFPQyxZQUFXRixNQUFLLFlBQVksRUFBQyxPQUFBQyxRQUFPLFFBQUFDLFFBQU0sQ0FBQztBQUFBLE1BQy9ELGdCQUFnQixDQUFDLFVBQVVGLE1BQUssZ0JBQWdCLEVBQUMsYUFBYSxNQUFLLENBQUM7QUFBQSxNQUNwRSxxQkFBcUIsQ0FBQyxHQUFHLE1BQU1BLE1BQUsscUJBQXFCLEVBQUMsR0FBRyxFQUFDLENBQUM7QUFBQSxNQUMvRCxrQkFBa0IsTUFBTUEsTUFBSyxnQkFBZ0I7QUFBQSxNQUM3QyxRQUFRLE1BQU1BLE1BQUssTUFBTTtBQUFBLE1BQ3pCLE1BQU0sTUFBTUEsTUFBSyxJQUFJO0FBQUEsTUFDckIsVUFBVSxNQUFNQSxNQUFLLFFBQVE7QUFBQSxNQUM3QixZQUFZLE1BQU1BLE1BQUssVUFBVTtBQUFBLE1BQ2pDLGdCQUFnQixNQUFNQSxNQUFLLGNBQWM7QUFBQSxNQUN6QyxVQUFVLE1BQU1BLE1BQUssUUFBUTtBQUFBLE1BQzdCLFlBQVksTUFBTUEsTUFBSyxVQUFVO0FBQUEsTUFDakMsU0FBUyxNQUFNQSxNQUFLLE9BQU87QUFBQSxNQUMzQixNQUFNLE1BQU1BLE1BQUssSUFBSTtBQUFBLE1BQ3JCLE9BQU8sTUFBTUEsTUFBSyxLQUFLO0FBQUEsTUFDdkIscUJBQXFCLENBQUMsR0FBRyxHQUFHLEdBQUcsTUFBTUEsTUFBSyxxQkFBcUIsRUFBQyxHQUFHLEdBQUcsR0FBRyxFQUFDLENBQUM7QUFBQSxNQUMzRSxjQUFjLENBQUMsY0FBY0EsTUFBSyxjQUFjLEVBQUMsVUFBUyxDQUFDO0FBQUEsTUFDM0QsT0FBTyxNQUFNQSxNQUFLLEtBQUs7QUFBQSxNQUN2QixRQUFRLE1BQU1BLE1BQUssTUFBTTtBQUFBLE1BQ3pCLFFBQVEsTUFBTUEsTUFBSyxNQUFNO0FBQUEsTUFDekIsU0FBUyxNQUFNQSxNQUFLLE9BQU87QUFBQSxNQUMzQixXQUFXLE1BQU1BLE1BQUssU0FBUztBQUFBLE1BQy9CLGNBQWMsTUFBTUEsTUFBSyxZQUFZO0FBQUEsTUFDckMsY0FBYyxDQUFDLGNBQWNBLE1BQUssY0FBYyxFQUFDLFVBQVMsQ0FBQztBQUFBLElBQy9EO0FBQUEsRUFDSjtBQVFPLFdBQVMsSUFBSSxZQUFZO0FBQzVCLFdBQU8sYUFBYSx1QkFBdUIsWUFBWSxRQUFRLFVBQVUsQ0FBQztBQUFBLEVBQzlFO0FBS08sV0FBUyxTQUFTO0FBQ3JCLGVBQVcsT0FBTztBQUFBLEVBQ3RCO0FBTU8sV0FBUyxTQUFTLE9BQU87QUFDNUIsZUFBVyxTQUFTLEtBQUs7QUFBQSxFQUM3QjtBQUtPLFdBQVMsYUFBYTtBQUN6QixlQUFXLFdBQVc7QUFBQSxFQUMxQjtBQU9PLFdBQVMsUUFBUUMsUUFBT0MsU0FBUTtBQUNuQyxlQUFXLFFBQVFELFFBQU9DLE9BQU07QUFBQSxFQUNwQztBQUtPLFdBQVMsT0FBTztBQUNuQixXQUFPLFdBQVcsS0FBSztBQUFBLEVBQzNCO0FBT08sV0FBUyxXQUFXRCxRQUFPQyxTQUFRO0FBQ3RDLGVBQVcsV0FBV0QsUUFBT0MsT0FBTTtBQUFBLEVBQ3ZDO0FBT08sV0FBUyxXQUFXRCxRQUFPQyxTQUFRO0FBQ3RDLGVBQVcsV0FBV0QsUUFBT0MsT0FBTTtBQUFBLEVBQ3ZDO0FBTU8sV0FBUyxlQUFlLE9BQU87QUFDbEMsZUFBVyxlQUFlLEtBQUs7QUFBQSxFQUNuQztBQU9PLFdBQVMsb0JBQW9CLEdBQUcsR0FBRztBQUN0QyxlQUFXLG9CQUFvQixHQUFHLENBQUM7QUFBQSxFQUN2QztBQUtPLFdBQVMsbUJBQW1CO0FBQy9CLFdBQU8sV0FBVyxpQkFBaUI7QUFBQSxFQUN2QztBQUtPLFdBQVMsU0FBUztBQUNyQixXQUFPLFdBQVcsT0FBTztBQUFBLEVBQzdCO0FBS08sV0FBU0MsUUFBTztBQUNuQixlQUFXLEtBQUs7QUFBQSxFQUNwQjtBQUtPLFdBQVMsV0FBVztBQUN2QixlQUFXLFNBQVM7QUFBQSxFQUN4QjtBQUtPLFdBQVMsYUFBYTtBQUN6QixlQUFXLFdBQVc7QUFBQSxFQUMxQjtBQUtPLFdBQVMsaUJBQWlCO0FBQzdCLGVBQVcsZUFBZTtBQUFBLEVBQzlCO0FBS08sV0FBUyxXQUFXO0FBQ3ZCLGVBQVcsU0FBUztBQUFBLEVBQ3hCO0FBS08sV0FBUyxhQUFhO0FBQ3pCLGVBQVcsV0FBVztBQUFBLEVBQzFCO0FBS08sV0FBUyxVQUFVO0FBQ3RCLGVBQVcsUUFBUTtBQUFBLEVBQ3ZCO0FBS08sV0FBU0MsUUFBTztBQUNuQixlQUFXLEtBQUs7QUFBQSxFQUNwQjtBQUtPLFdBQVMsUUFBUTtBQUNwQixlQUFXLE1BQU07QUFBQSxFQUNyQjtBQVNPLFdBQVMsb0JBQW9CLEdBQUcsR0FBRyxHQUFHLEdBQUc7QUFDNUMsZUFBVyxvQkFBb0IsR0FBRyxHQUFHLEdBQUcsQ0FBQztBQUFBLEVBQzdDO0FBTU8sV0FBUyxhQUFhLFdBQVc7QUFDcEMsZUFBVyxhQUFhLFNBQVM7QUFBQSxFQUNyQztBQUtPLFdBQVMsUUFBUTtBQUNwQixXQUFPLFdBQVcsTUFBTTtBQUFBLEVBQzVCO0FBS08sV0FBUyxTQUFTO0FBQ3JCLFdBQU8sV0FBVyxPQUFPO0FBQUEsRUFDN0I7QUFLTyxXQUFTLFNBQVM7QUFDckIsZUFBVyxPQUFPO0FBQUEsRUFDdEI7QUFLTyxXQUFTLFVBQVU7QUFDdEIsZUFBVyxRQUFRO0FBQUEsRUFDdkI7QUFLTyxXQUFTLFlBQVk7QUFDeEIsZUFBVyxVQUFVO0FBQUEsRUFDekI7QUFLTyxXQUFTLGVBQWU7QUFDM0IsV0FBTyxXQUFXLGFBQWE7QUFBQSxFQUNuQztBQU1PLFdBQVMsYUFBYSxXQUFXO0FBQ3BDLGVBQVcsYUFBYSxTQUFTO0FBQUEsRUFDckM7OztBQzNUQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNDTyxNQUFNLGFBQWE7QUFBQSxJQUN6QixTQUFTO0FBQUEsTUFDUixvQkFBb0I7QUFBQSxNQUNwQixzQkFBc0I7QUFBQSxNQUN0QixZQUFZO0FBQUEsTUFDWixvQkFBb0I7QUFBQSxNQUNwQixrQkFBa0I7QUFBQSxNQUNsQix1QkFBdUI7QUFBQSxNQUN2QixvQkFBb0I7QUFBQSxNQUNwQiw0QkFBNEI7QUFBQSxNQUM1QixnQkFBZ0I7QUFBQSxNQUNoQixjQUFjO0FBQUEsTUFDZCxtQkFBbUI7QUFBQSxNQUNuQixnQkFBZ0I7QUFBQSxNQUNoQixrQkFBa0I7QUFBQSxNQUNsQixrQkFBa0I7QUFBQSxNQUNsQixvQkFBb0I7QUFBQSxNQUNwQixlQUFlO0FBQUEsTUFDZixnQkFBZ0I7QUFBQSxNQUNoQixrQkFBa0I7QUFBQSxNQUNsQixhQUFhO0FBQUEsTUFDYixnQkFBZ0I7QUFBQSxNQUNoQixpQkFBaUI7QUFBQSxNQUNqQixnQkFBZ0I7QUFBQSxNQUNoQixpQkFBaUI7QUFBQSxNQUNqQixpQkFBaUI7QUFBQSxNQUNqQixnQkFBZ0I7QUFBQSxJQUNqQjtBQUFBLElBQ0EsS0FBSztBQUFBLE1BQ0osNEJBQTRCO0FBQUEsTUFDNUIsdUNBQXVDO0FBQUEsTUFDdkMseUNBQXlDO0FBQUEsTUFDekMsMEJBQTBCO0FBQUEsTUFDMUIsb0NBQW9DO0FBQUEsTUFDcEMsc0NBQXNDO0FBQUEsTUFDdEMsb0NBQW9DO0FBQUEsTUFDcEMsMENBQTBDO0FBQUEsTUFDMUMsK0JBQStCO0FBQUEsTUFDL0Isb0JBQW9CO0FBQUEsTUFDcEIsd0NBQXdDO0FBQUEsTUFDeEMsc0JBQXNCO0FBQUEsTUFDdEIsc0JBQXNCO0FBQUEsTUFDdEIsNkJBQTZCO0FBQUEsTUFDN0IsZ0NBQWdDO0FBQUEsTUFDaEMscUJBQXFCO0FBQUEsTUFDckIsNkJBQTZCO0FBQUEsTUFDN0IsMEJBQTBCO0FBQUEsTUFDMUIsdUJBQXVCO0FBQUEsTUFDdkIsdUJBQXVCO0FBQUEsTUFDdkIsMkJBQTJCO0FBQUEsTUFDM0IsK0JBQStCO0FBQUEsTUFDL0Isb0JBQW9CO0FBQUEsTUFDcEIscUJBQXFCO0FBQUEsTUFDckIscUJBQXFCO0FBQUEsTUFDckIsc0JBQXNCO0FBQUEsTUFDdEIsZ0NBQWdDO0FBQUEsTUFDaEMsa0NBQWtDO0FBQUEsTUFDbEMsbUNBQW1DO0FBQUEsTUFDbkMsb0NBQW9DO0FBQUEsTUFDcEMsK0JBQStCO0FBQUEsTUFDL0IsNkJBQTZCO0FBQUEsTUFDN0IsdUJBQXVCO0FBQUEsTUFDdkIsaUNBQWlDO0FBQUEsTUFDakMsOEJBQThCO0FBQUEsTUFDOUIsNEJBQTRCO0FBQUEsTUFDNUIsc0NBQXNDO0FBQUEsTUFDdEMsNEJBQTRCO0FBQUEsTUFDNUIsc0JBQXNCO0FBQUEsTUFDdEIsa0NBQWtDO0FBQUEsTUFDbEMsc0JBQXNCO0FBQUEsTUFDdEIsd0JBQXdCO0FBQUEsTUFDeEIsMkJBQTJCO0FBQUEsTUFDM0Isd0JBQXdCO0FBQUEsTUFDeEIsbUJBQW1CO0FBQUEsTUFDbkIsMEJBQTBCO0FBQUEsTUFDMUIsOEJBQThCO0FBQUEsTUFDOUIseUJBQXlCO0FBQUEsTUFDekIsNkJBQTZCO0FBQUEsTUFDN0IsaUJBQWlCO0FBQUEsTUFDakIsZ0JBQWdCO0FBQUEsTUFDaEIsc0JBQXNCO0FBQUEsTUFDdEIsZUFBZTtBQUFBLE1BQ2YseUJBQXlCO0FBQUEsTUFDekIsd0JBQXdCO0FBQUEsTUFDeEIsb0JBQW9CO0FBQUEsTUFDcEIscUJBQXFCO0FBQUEsTUFDckIsaUJBQWlCO0FBQUEsTUFDakIsaUJBQWlCO0FBQUEsTUFDakIsc0JBQXNCO0FBQUEsTUFDdEIsbUNBQW1DO0FBQUEsTUFDbkMscUNBQXFDO0FBQUEsTUFDckMsdUJBQXVCO0FBQUEsTUFDdkIsc0JBQXNCO0FBQUEsTUFDdEIsd0JBQXdCO0FBQUEsTUFDeEIsMkJBQTJCO0FBQUEsTUFDM0IsbUJBQW1CO0FBQUEsTUFDbkIscUJBQXFCO0FBQUEsTUFDckIsc0JBQXNCO0FBQUEsTUFDdEIsc0JBQXNCO0FBQUEsTUFDdEIsOEJBQThCO0FBQUEsTUFDOUIsaUJBQWlCO0FBQUEsTUFDakIseUJBQXlCO0FBQUEsTUFDekIsMkJBQTJCO0FBQUEsTUFDM0IsK0JBQStCO0FBQUEsTUFDL0IsMEJBQTBCO0FBQUEsTUFDMUIsOEJBQThCO0FBQUEsTUFDOUIsaUJBQWlCO0FBQUEsTUFDakIsdUJBQXVCO0FBQUEsTUFDdkIsZ0JBQWdCO0FBQUEsTUFDaEIsMEJBQTBCO0FBQUEsTUFDMUIseUJBQXlCO0FBQUEsTUFDekIsc0JBQXNCO0FBQUEsTUFDdEIsa0JBQWtCO0FBQUEsTUFDbEIsbUJBQW1CO0FBQUEsTUFDbkIsa0JBQWtCO0FBQUEsTUFDbEIsdUJBQXVCO0FBQUEsTUFDdkIsb0NBQW9DO0FBQUEsTUFDcEMsc0NBQXNDO0FBQUEsTUFDdEMsd0JBQXdCO0FBQUEsTUFDeEIsdUJBQXVCO0FBQUEsTUFDdkIseUJBQXlCO0FBQUEsTUFDekIsNEJBQTRCO0FBQUEsTUFDNUIsNEJBQTRCO0FBQUEsTUFDNUIsY0FBYztBQUFBLE1BQ2QsYUFBYTtBQUFBLE1BQ2IsY0FBYztBQUFBLE1BQ2Qsb0JBQW9CO0FBQUEsTUFDcEIsbUJBQW1CO0FBQUEsTUFDbkIsdUJBQXVCO0FBQUEsTUFDdkIsc0JBQXNCO0FBQUEsTUFDdEIscUJBQXFCO0FBQUEsTUFDckIsb0JBQW9CO0FBQUEsTUFDcEIsaUJBQWlCO0FBQUEsTUFDakIsZ0JBQWdCO0FBQUEsTUFDaEIsb0JBQW9CO0FBQUEsTUFDcEIsbUJBQW1CO0FBQUEsTUFDbkIsdUJBQXVCO0FBQUEsTUFDdkIsc0JBQXNCO0FBQUEsTUFDdEIscUJBQXFCO0FBQUEsTUFDckIsb0JBQW9CO0FBQUEsTUFDcEIsZ0JBQWdCO0FBQUEsTUFDaEIsZUFBZTtBQUFBLE1BQ2YsZUFBZTtBQUFBLE1BQ2YsY0FBYztBQUFBLE1BQ2QsMEJBQTBCO0FBQUEsTUFD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TyxZQUFZLGFBQWE7QUFBQSxNQUM1QixhQUFhO0FBQUEsTUFDYixZQUFZO0FBQUEsTUFDWjtBQUFBLE1BQ0E7QUFBQSxJQUNKLENBQUM7QUFBQSxFQUNMOzs7QUNwSkEsU0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBOENsQyxTQUFPLE9BQU8sU0FBTztBQUVyQixTQUFPLFFBQVEsT0FBTyxTQUFTLENBQUM7QUFDaEMsU0FBTyxNQUFNLGNBQWM7QUFDM0IsU0FBTyxNQUFNLFVBQVU7QUFDdkIsU0FBTyxNQUFNLE9BQU87QUFDcEIsU0FBTyxNQUFNLFlBQVk7QUFDekIsU0FBTyxNQUFNLFVBQVU7QUFDdkIsU0FBTyxNQUFNLFNBQVM7QUFDdEIsU0FBTyxNQUFNLFFBQVE7QUFDckIsU0FBTyxNQUFNLFVBQVU7QUFDdkIsU0FBTyxNQUFNLFNBQVM7QUFDdEIsU0FBTyxNQUFNLFNBQVM7QUFDdEIsU0FBTyxNQUFNLE1BQU07QUFHbkIsTUFBSSxVQUFVO0FBQ2QsV0FBUyxpQkFBaUIsb0JBQW9CLFdBQVc7QUFDckQsY0FBVTtBQUNWLFdBQU8sT0FBTyxPQUFPLHFCQUFxQjtBQUMxQyxRQUFHLE1BQU87QUFDTixlQUFTLHNCQUFzQjtBQUFBLElBQ25DO0FBQUEsRUFDSixDQUFDO0FBRUQsV0FBUyxVQUFVLElBQUk7QUFDbkIsUUFBSSxXQUFXLFNBQVMsZUFBZSxZQUFZO0FBQy9DLFNBQUc7QUFBQSxJQUNQLE9BQU87QUFDSCxlQUFTLGlCQUFpQixvQkFBb0IsRUFBRTtBQUFBLElBQ3BEO0FBQUEsRUFDSjtBQUVBLFlBQVUsTUFBTTtBQUNaLElBQUksT0FBTztBQUFBLEVBQ2YsQ0FBQzsiLAogICJuYW1lcyI6IFsic2l6ZSIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgIkhpZGUiLCAiU2hvdyIsICJjYWxsIiwgIndpZHRoIiwgImhlaWdodCIsICJIaWRlIiwgIlNob3ciLCAiY2FsbCIsICJldmVudE5hbWUiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJjYWxsIiwgImdlbmVyYXRlSUQiXQp9Cg==
