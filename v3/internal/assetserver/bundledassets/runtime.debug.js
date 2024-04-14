var __defProp = Object.defineProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};

// desktop/@wailsio/runtime/src/index.js
var src_exports = {};
__export(src_exports, {
  Application: () => application_exports,
  Browser: () => browser_exports,
  Call: () => calls_exports,
  Clipboard: () => clipboard_exports,
  Dialogs: () => dialogs_exports,
  Events: () => events_exports,
  Flags: () => flags_exports,
  Screens: () => screens_exports,
  System: () => system_exports,
  WML: () => wml_exports,
  Window: () => window_default
});

// desktop/@wailsio/runtime/src/wml.js
var wml_exports = {};
__export(wml_exports, {
  Enable: () => Enable,
  Reload: () => Reload
});

// desktop/@wailsio/runtime/src/browser.js
var browser_exports = {};
__export(browser_exports, {
  OpenURL: () => OpenURL
});

// node_modules/nanoid/non-secure/index.js
var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
var nanoid = (size = 21) => {
  let id = "";
  let i = size;
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

// desktop/@wailsio/runtime/src/browser.js
var call = newRuntimeCallerWithID(objectNames.Browser, "");
var BrowserOpenURL = 0;
function OpenURL(url) {
  return call(BrowserOpenURL, { url });
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
var call2 = newRuntimeCallerWithID(objectNames.Dialog, "");
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
    call2(type, options).catch((error) => {
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
var call3 = newRuntimeCallerWithID(objectNames.Events, "");
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
  return call3(EmitMethod, event);
}

// desktop/@wailsio/runtime/src/utils.js
function debugLog(message) {
  console.log(
    "%c wails3 %c " + message + " ",
    "background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem",
    "background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem"
  );
}
function canAbortListeners() {
  if (!EventTarget || !AbortSignal || !AbortController)
    return false;
  let result = true;
  const target = new EventTarget();
  const controller2 = new AbortController();
  target.addEventListener("test", () => {
    result = false;
  }, { signal: controller2.signal });
  controller2.abort();
  target.dispatchEvent(new CustomEvent("test"));
  return result;
}
var isReady = false;
document.addEventListener("DOMContentLoaded", () => isReady = true);
function whenReady(callback) {
  if (isReady || document.readyState === "complete") {
    callback();
  } else {
    document.addEventListener("DOMContentLoaded", callback);
  }
}

// desktop/@wailsio/runtime/src/window.js
var AbsolutePositionMethod = 0;
var CenterMethod = 1;
var CloseMethod = 2;
var DisableSizeConstraintsMethod = 3;
var EnableSizeConstraintsMethod = 4;
var FocusMethod = 5;
var ForceReloadMethod = 6;
var FullscreenMethod = 7;
var GetScreenMethod = 8;
var GetZoomMethod = 9;
var HeightMethod = 10;
var HideMethod = 11;
var IsFocusedMethod = 12;
var IsFullscreenMethod = 13;
var IsMaximisedMethod = 14;
var IsMinimisedMethod = 15;
var MaximiseMethod = 16;
var MinimiseMethod = 17;
var NameMethod = 18;
var OpenDevToolsMethod = 19;
var RelativePositionMethod = 20;
var ReloadMethod = 21;
var ResizableMethod = 22;
var RestoreMethod = 23;
var SetAbsolutePositionMethod = 24;
var SetAlwaysOnTopMethod = 25;
var SetBackgroundColourMethod = 26;
var SetFramelessMethod = 27;
var SetFullscreenButtonEnabledMethod = 28;
var SetMaxSizeMethod = 29;
var SetMinSizeMethod = 30;
var SetRelativePositionMethod = 31;
var SetResizableMethod = 32;
var SetSizeMethod = 33;
var SetTitleMethod = 34;
var SetZoomMethod = 35;
var ShowMethod = 36;
var SizeMethod = 37;
var ToggleFullscreenMethod = 38;
var ToggleMaximiseMethod = 39;
var UnFullscreenMethod = 40;
var UnMaximiseMethod = 41;
var UnMinimiseMethod = 42;
var WidthMethod = 43;
var ZoomMethod = 44;
var ZoomInMethod = 45;
var ZoomOutMethod = 46;
var ZoomResetMethod = 47;
var caller = Symbol();
var Window = class _Window {
  /**
   * Initialises a window object with the specified name.
   *
   * @private
   * @param {string} name - The name of the target window.
   */
  constructor(name = "") {
    this[caller] = newRuntimeCallerWithID(objectNames.Window, name);
    for (const method of Object.getOwnPropertyNames(_Window.prototype)) {
      if (method !== "constructor" && typeof this[method] === "function") {
        this[method] = this[method].bind(this);
      }
    }
  }
  /**
   * Gets the specified window.
   *
   * @public
   * @param {string} name - The name of the window to get.
   * @return {Window} - The corresponding window object.
   */
  Get(name) {
    return new _Window(name);
  }
  /**
   * Returns the absolute position of the window to the screen.
   *
   * @public
   * @return {Promise<Position>} - The current absolute position of the window.
   */
  AbsolutePosition() {
    return this[caller](AbsolutePositionMethod);
  }
  /**
   * Centers the window on the screen.
   *
   * @public
   * @return {Promise<void>}
   */
  Center() {
    return this[caller](CenterMethod);
  }
  /**
   * Closes the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Close() {
    return this[caller](CloseMethod);
  }
  /**
   * Disables min/max size constraints.
   *
   * @public
   * @return {Promise<void>}
   */
  DisableSizeConstraints() {
    return this[caller](DisableSizeConstraintsMethod);
  }
  /**
   * Enables min/max size constraints.
   *
   * @public
   * @return {Promise<void>}
   */
  EnableSizeConstraints() {
    return this[caller](EnableSizeConstraintsMethod);
  }
  /**
   * Focuses the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Focus() {
    return this[caller](FocusMethod);
  }
  /**
   * Forces the window to reload the page assets.
   *
   * @public
   * @return {Promise<void>}
   */
  ForceReload() {
    return this[caller](ForceReloadMethod);
  }
  /**
   * Doc.
   *
   * @public
   * @return {Promise<void>}
   */
  Fullscreen() {
    return this[caller](FullscreenMethod);
  }
  /**
   * Returns the screen that the window is on.
   *
   * @public
   * @return {Promise<Screen>} - The screen the window is currently on
   */
  GetScreen() {
    return this[caller](GetScreenMethod);
  }
  /**
   * Returns the current zoom level of the window.
   *
   * @public
   * @return {Promise<number>} - The current zoom level
   */
  GetZoom() {
    return this[caller](GetZoomMethod);
  }
  /**
   * Returns the height of the window.
   *
   * @public
   * @return {Promise<number>} - The current height of the window
   */
  Height() {
    return this[caller](HeightMethod);
  }
  /**
   * Hides the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Hide() {
    return this[caller](HideMethod);
  }
  /**
   * Returns true if the window is focused.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently focused
   */
  IsFocused() {
    return this[caller](IsFocusedMethod);
  }
  /**
   * Returns true if the window is fullscreen.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently fullscreen
   */
  IsFullscreen() {
    return this[caller](IsFullscreenMethod);
  }
  /**
   * Returns true if the window is maximised.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently maximised
   */
  IsMaximised() {
    return this[caller](IsMaximisedMethod);
  }
  /**
   * Returns true if the window is minimised.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently minimised
   */
  IsMinimised() {
    return this[caller](IsMinimisedMethod);
  }
  /**
   * Maximises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Maximise() {
    return this[caller](MaximiseMethod);
  }
  /**
   * Minimises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Minimise() {
    return this[caller](MinimiseMethod);
  }
  /**
   * Returns the name of the window.
   *
   * @public
   * @return {Promise<string>} - The name of the window
   */
  Name() {
    return this[caller](NameMethod);
  }
  /**
   * Opens the development tools pane.
   *
   * @public
   * @return {Promise<void>}
   */
  OpenDevTools() {
    return this[caller](OpenDevToolsMethod);
  }
  /**
   * Returns the relative position of the window to the screen.
   *
   * @public
   * @return {Promise<Position>} - The current relative position of the window
   */
  RelativePosition() {
    return this[caller](RelativePositionMethod);
  }
  /**
   * Reloads the page assets.
   *
   * @public
   * @return {Promise<void>}
   */
  Reload() {
    return this[caller](ReloadMethod);
  }
  /**
   * Returns true if the window is resizable.
   *
   * @public
   * @return {Promise<boolean>} - Whether the window is currently resizable
   */
  Resizable() {
    return this[caller](ResizableMethod);
  }
  /**
   * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
   *
   * @public
   * @return {Promise<void>}
   */
  Restore() {
    return this[caller](RestoreMethod);
  }
  /**
   * Sets the absolute position of the window to the screen.
   *
   * @public
   * @param {number} x - The desired horizontal absolute position of the window
   * @param {number} y - The desired vertical absolute position of the window
   * @return {Promise<void>}
   */
  SetAbsolutePosition(x, y) {
    return this[caller](SetAbsolutePositionMethod, { x, y });
  }
  /**
   * Sets the window to be always on top.
   *
   * @public
   * @param {boolean} alwaysOnTop - Whether the window should stay on top
   * @return {Promise<void>}
   */
  SetAlwaysOnTop(alwaysOnTop) {
    return this[caller](SetAlwaysOnTopMethod, { alwaysOnTop });
  }
  /**
   * Sets the background colour of the window.
   *
   * @public
   * @param {number} r - The desired red component of the window background
   * @param {number} g - The desired green component of the window background
   * @param {number} b - The desired blue component of the window background
   * @param {number} a - The desired alpha component of the window background
   * @return {Promise<void>}
   */
  SetBackgroundColour(r, g, b, a) {
    return this[caller](SetBackgroundColourMethod, { r, g, b, a });
  }
  /**
   * Removes the window frame and title bar.
   *
   * @public
   * @param {boolean} frameless - Whether the window should be frameless
   * @return {Promise<void>}
   */
  SetFrameless(frameless) {
    return this[caller](SetFramelessMethod, { frameless });
  }
  /**
   * Disables the system fullscreen button.
   *
   * @public
   * @param {boolean} enabled - Whether the fullscreen button should be enabled
   * @return {Promise<void>}
   */
  SetFullscreenButtonEnabled(enabled) {
    return this[caller](SetFullscreenButtonEnabledMethod, { enabled });
  }
  /**
   * Sets the maximum size of the window.
   *
   * @public
   * @param {number} width - The desired maximum width of the window
   * @param {number} height - The desired maximum height of the window
   * @return {Promise<void>}
   */
  SetMaxSize(width, height) {
    return this[caller](SetMaxSizeMethod, { width, height });
  }
  /**
   * Sets the minimum size of the window.
   *
   * @public
   * @param {number} width - The desired minimum width of the window
   * @param {number} height - The desired minimum height of the window
   * @return {Promise<void>}
   */
  SetMinSize(width, height) {
    return this[caller](SetMinSizeMethod, { width, height });
  }
  /**
   * Sets the relative position of the window to the screen.
   *
   * @public
   * @param {number} x - The desired horizontal relative position of the window
   * @param {number} y - The desired vertical relative position of the window
   * @return {Promise<void>}
   */
  SetRelativePosition(x, y) {
    return this[caller](SetRelativePositionMethod, { x, y });
  }
  /**
   * Sets whether the window is resizable.
   *
   * @public
   * @param {boolean} resizable - Whether the window should be resizable
   * @return {Promise<void>}
   */
  SetResizable(resizable2) {
    return this[caller](SetResizableMethod, { resizable: resizable2 });
  }
  /**
   * Sets the size of the window.
   *
   * @public
   * @param {number} width - The desired width of the window
   * @param {number} height - The desired height of the window
   * @return {Promise<void>}
   */
  SetSize(width, height) {
    return this[caller](SetSizeMethod, { width, height });
  }
  /**
   * Sets the title of the window.
   *
   * @public
   * @param {string} title - The desired title of the window
   * @return {Promise<void>}
   */
  SetTitle(title) {
    return this[caller](SetTitleMethod, { title });
  }
  /**
   * Sets the zoom level of the window.
   *
   * @public
   * @param {number} zoom - The desired zoom level
   * @return {Promise<void>}
   */
  SetZoom(zoom) {
    return this[caller](SetZoomMethod, { zoom });
  }
  /**
   * Shows the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Show() {
    return this[caller](ShowMethod);
  }
  /**
   * Returns the size of the window.
   *
   * @public
   * @return {Promise<Size>} - The current size of the window
   */
  Size() {
    return this[caller](SizeMethod);
  }
  /**
   * Toggles the window between fullscreen and normal.
   *
   * @public
   * @return {Promise<void>}
   */
  ToggleFullscreen() {
    return this[caller](ToggleFullscreenMethod);
  }
  /**
   * Toggles the window between maximised and normal.
   *
   * @public
   * @return {Promise<void>}
   */
  ToggleMaximise() {
    return this[caller](ToggleMaximiseMethod);
  }
  /**
   * Un-fullscreens the window.
   *
   * @public
   * @return {Promise<void>}
   */
  UnFullscreen() {
    return this[caller](UnFullscreenMethod);
  }
  /**
   * Un-maximises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  UnMaximise() {
    return this[caller](UnMaximiseMethod);
  }
  /**
   * Un-minimises the window.
   *
   * @public
   * @return {Promise<void>}
   */
  UnMinimise() {
    return this[caller](UnMinimiseMethod);
  }
  /**
   * Returns the width of the window.
   *
   * @public
   * @return {Promise<number>} - The current width of the window
   */
  Width() {
    return this[caller](WidthMethod);
  }
  /**
   * Zooms the window.
   *
   * @public
   * @return {Promise<void>}
   */
  Zoom() {
    return this[caller](ZoomMethod);
  }
  /**
   * Increases the zoom level of the webview content.
   *
   * @public
   * @return {Promise<void>}
   */
  ZoomIn() {
    return this[caller](ZoomInMethod);
  }
  /**
   * Decreases the zoom level of the webview content.
   *
   * @public
   * @return {Promise<void>}
   */
  ZoomOut() {
    return this[caller](ZoomOutMethod);
  }
  /**
   * Resets the zoom level of the webview content.
   *
   * @public
   * @return {Promise<void>}
   */
  ZoomReset() {
    return this[caller](ZoomResetMethod);
  }
};
var thisWindow = new Window("");
var window_default = thisWindow;

// desktop/@wailsio/runtime/src/wml.js
function sendEvent(eventName, data = null) {
  Emit(new WailsEvent(eventName, data));
}
function callWindowMethod(windowName, methodName) {
  const targetWindow = window_default.Get(windowName);
  const method = targetWindow[methodName];
  if (typeof method !== "function") {
    console.error(`Window method '${methodName}' not found`);
    return;
  }
  try {
    method.call(targetWindow);
  } catch (e) {
    console.error(`Error calling window method '${methodName}': `, e);
  }
}
function onWMLTriggered(ev) {
  const element = ev.currentTarget;
  function runEffect(choice = "Yes") {
    if (choice !== "Yes")
      return;
    const eventType = element.getAttribute("wml-event");
    const targetWindow = element.getAttribute("wml-target-window") || "";
    const windowMethod = element.getAttribute("wml-window");
    const url = element.getAttribute("wml-openurl");
    if (eventType !== null)
      sendEvent(eventType);
    if (windowMethod !== null)
      callWindowMethod(targetWindow, windowMethod);
    if (url !== null)
      void OpenURL(url);
  }
  const confirm = element.getAttribute("wml-confirm");
  if (confirm) {
    Question({
      Title: "Confirm",
      Message: confirm,
      Detached: false,
      Buttons: [
        { Label: "Yes" },
        { Label: "No", IsDefault: true }
      ]
    }).then(runEffect);
  } else {
    runEffect();
  }
}
var controller = Symbol();
var AbortControllerRegistry = class {
  constructor() {
    this[controller] = new AbortController();
  }
  /**
   * Returns an options object for addEventListener that ties the listener
   * to the AbortSignal from the current AbortController.
   *
   * @param {HTMLElement} element An HTML element
   * @param {string[]} triggers The list of active WML trigger events for the specified elements
   * @returns {AddEventListenerOptions}
   */
  set(element, triggers) {
    return { signal: this[controller].signal };
  }
  /**
   * Removes all registered event listeners.
   *
   * @returns {void}
   */
  reset() {
    this[controller].abort();
    this[controller] = new AbortController();
  }
};
var triggerMap = Symbol();
var elementCount = Symbol();
var WeakMapRegistry = class {
  constructor() {
    this[triggerMap] = /* @__PURE__ */ new WeakMap();
    this[elementCount] = 0;
  }
  /**
   * Sets the active triggers for the specified element.
   *
   * @param {HTMLElement} element An HTML element
   * @param {string[]} triggers The list of active WML trigger events for the specified element
   * @returns {AddEventListenerOptions}
   */
  set(element, triggers) {
    this[elementCount] += !this[triggerMap].has(element);
    this[triggerMap].set(element, triggers);
    return {};
  }
  /**
   * Removes all registered event listeners.
   *
   * @returns {void}
   */
  reset() {
    if (this[elementCount] <= 0)
      return;
    for (const element of document.body.querySelectorAll("*")) {
      if (this[elementCount] <= 0)
        break;
      const triggers = this[triggerMap].get(element);
      this[elementCount] -= typeof triggers !== "undefined";
      for (const trigger of triggers || [])
        element.removeEventListener(trigger, onWMLTriggered);
    }
    this[triggerMap] = /* @__PURE__ */ new WeakMap();
    this[elementCount] = 0;
  }
};
var triggerRegistry = canAbortListeners() ? new AbortControllerRegistry() : new WeakMapRegistry();
function addWMLListeners(element) {
  const triggerRegExp = /\S+/g;
  const triggerAttr = element.getAttribute("wml-trigger") || "click";
  const triggers = [];
  let match;
  while ((match = triggerRegExp.exec(triggerAttr)) !== null)
    triggers.push(match[0]);
  const options = triggerRegistry.set(element, triggers);
  for (const trigger of triggers)
    element.addEventListener(trigger, onWMLTriggered, options);
}
function Enable() {
  whenReady(Reload);
}
function Reload() {
  triggerRegistry.reset();
  document.body.querySelectorAll("[wml-event], [wml-window], [wml-openurl]").forEach(addWMLListeners);
}

// desktop/compiled/main.js
window.wails = src_exports;
Enable();
if (true) {
  debugLog("Wails Runtime Loaded");
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
var shouldDrag = false;
var resizable = false;
var resizeEdge = null;
var defaultCursor = "auto";
window._wails = window._wails || {};
window._wails.setResizable = function(value) {
  resizable = value;
};
window._wails.endDrag = function() {
  document.body.style.cursor = "default";
  shouldDrag = false;
};
window.addEventListener("mousedown", onMouseDown);
window.addEventListener("mousemove", onMouseMove);
window.addEventListener("mouseup", onMouseUp);
function dragTest(e) {
  let val = window.getComputedStyle(e.target).getPropertyValue("--wails-draggable");
  let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
  if (!val || val === "" || val.trim() !== "drag" || mousePressed === 0) {
    return false;
  }
  return e.detail === 1;
}
function onMouseDown(e) {
  if (resizeEdge) {
    invoke("resize:" + resizeEdge);
    e.preventDefault();
    return;
  }
  if (dragTest(e)) {
    if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
      return;
    }
    shouldDrag = true;
  } else {
    shouldDrag = false;
  }
}
function onMouseUp() {
  shouldDrag = false;
}
function setResize(cursor) {
  document.documentElement.style.cursor = cursor || defaultCursor;
  resizeEdge = cursor;
}
function onMouseMove(e) {
  if (shouldDrag) {
    shouldDrag = false;
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (mousePressed > 0) {
      invoke("drag");
      return;
    }
  }
  if (!resizable || !IsWindows()) {
    return;
  }
  if (defaultCursor == null) {
    defaultCursor = document.documentElement.style.cursor;
  }
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

// desktop/@wailsio/runtime/src/application.js
var application_exports = {};
__export(application_exports, {
  Hide: () => Hide,
  Quit: () => Quit,
  Show: () => Show
});
var call6 = newRuntimeCallerWithID(objectNames.Application, "");
var HideMethod2 = 0;
var ShowMethod2 = 1;
var QuitMethod = 2;
function Hide() {
  return call6(HideMethod2);
}
function Show() {
  return call6(ShowMethod2);
}
function Quit() {
  return call6(QuitMethod);
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
var call7 = newRuntimeCallerWithID(objectNames.Call, "");
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
    call7(type, options).then((_) => {
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
  if (typeof name !== "string" || name.split(".").length < 3) {
    throw new Error("CallByName requires a string in the format 'package.struct.method'");
  }
  return callBinding(CallBinding, {
    methodName: name,
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

// desktop/@wailsio/runtime/src/clipboard.js
var clipboard_exports = {};
__export(clipboard_exports, {
  SetText: () => SetText,
  Text: () => Text
});
var call8 = newRuntimeCallerWithID(objectNames.Clipboard, "");
var ClipboardSetText = 0;
var ClipboardText = 1;
function SetText(text) {
  return call8(ClipboardSetText, { text });
}
function Text() {
  return call8(ClipboardText);
}

// desktop/@wailsio/runtime/src/screens.js
var screens_exports = {};
__export(screens_exports, {
  GetAll: () => GetAll,
  GetCurrent: () => GetCurrent,
  GetPrimary: () => GetPrimary
});
var call9 = newRuntimeCallerWithID(objectNames.Screens, "");
var getAll = 0;
var getPrimary = 1;
var getCurrent = 2;
function GetAll() {
  return call9(getAll);
}
function GetPrimary() {
  return call9(getPrimary);
}
function GetCurrent() {
  return call9(getCurrent);
}

// desktop/@wailsio/runtime/src/index.js
window._wails = window._wails || {};
window._wails.invoke = invoke;
invoke("wails:runtime:ready");
export {
  application_exports as Application,
  browser_exports as Browser,
  calls_exports as Call,
  clipboard_exports as Clipboard,
  dialogs_exports as Dialogs,
  events_exports as Events,
  flags_exports as Flags,
  screens_exports as Screens,
  system_exports as System,
  wml_exports as WML,
  window_default as Window
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vcnVudGltZS9ub2RlX21vZHVsZXMvbmFub2lkL25vbi1zZWN1cmUvaW5kZXguanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RpYWxvZ3MuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50cy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHBsaWNhdGlvbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHMuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NsaXBib2FyZC5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc2NyZWVucy5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxuaW1wb3J0IFwiLi9jb250ZXh0bWVudVwiO1xuaW1wb3J0IFwiLi9kcmFnXCI7XG5cbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi9hcHBsaWNhdGlvblwiO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyXCI7XG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkXCI7XG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3NcIjtcbmltcG9ydCAqIGFzIEV2ZW50cyBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuL2ZsYWdzXCI7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnNcIjtcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi9zeXN0ZW1cIjtcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sXCI7XG5cbmV4cG9ydCB7XG4gICAgQXBwbGljYXRpb24sXG4gICAgQnJvd3NlcixcbiAgICBDYWxsLFxuICAgIENsaXBib2FyZCxcbiAgICBEaWFsb2dzLFxuICAgIEV2ZW50cyxcbiAgICBGbGFncyxcbiAgICBTY3JlZW5zLFxuICAgIFN5c3RlbSxcbiAgICBXaW5kb3csXG4gICAgV01MXG59O1xuXG4vLyBOb3RpZnkgYmFja2VuZFxud2luZG93Ll93YWlscy5pbnZva2UgPSBTeXN0ZW0uaW52b2tlO1xuU3lzdGVtLmludm9rZShcIndhaWxzOnJ1bnRpbWU6cmVhZHlcIik7XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7T3BlblVSTH0gZnJvbSBcIi4vYnJvd3NlclwiO1xuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7Y2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeX0gZnJvbSBcIi4vdXRpbHNcIjtcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93XCI7XG5cbi8qKlxuICogU2VuZHMgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgb3B0aW9uYWwgZGF0YS5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXG4gKiBAcGFyYW0ge2FueX0gW2RhdGE9bnVsbF0gLSBPcHRpb25hbCBkYXRhIHRvIHNlbmQgYWxvbmcgd2l0aCB0aGUgZXZlbnQuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZSwgZGF0YT1udWxsKSB7XG4gICAgRW1pdChuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpKTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKi9cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TmFtZSwgbWV0aG9kTmFtZSkge1xuICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IFdpbmRvdy5HZXQod2luZG93TmFtZSk7XG4gICAgY29uc3QgbWV0aG9kID0gdGFyZ2V0V2luZG93W21ldGhvZE5hbWVdO1xuXG4gICAgaWYgKHR5cGVvZiBtZXRob2QgIT09IFwiZnVuY3Rpb25cIikge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBXaW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JyBub3QgZm91bmRgKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICAgIG1ldGhvZC5jYWxsKHRhcmdldFdpbmRvdyk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmVzcG9uZHMgdG8gYSB0cmlnZ2VyaW5nIGV2ZW50IGJ5IHJ1bm5pbmcgYXBwcm9wcmlhdGUgV01MIGFjdGlvbnMgZm9yIHRoZSBjdXJyZW50IHRhcmdldFxuICpcbiAqIEBwYXJhbSB7RXZlbnR9IGV2XG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5mdW5jdGlvbiBvbldNTFRyaWdnZXJlZChldikge1xuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0O1xuXG4gICAgZnVuY3Rpb24gcnVuRWZmZWN0KGNob2ljZSA9IFwiWWVzXCIpIHtcbiAgICAgICAgaWYgKGNob2ljZSAhPT0gXCJZZXNcIilcbiAgICAgICAgICAgIHJldHVybjtcblxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50Jyk7XG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtb3BlbnVybCcpO1xuXG4gICAgICAgIGlmIChldmVudFR5cGUgIT09IG51bGwpXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICBpZiAodXJsICE9PSBudWxsKVxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgfVxuXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xuXG4gICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgUXVlc3Rpb24oe1xuICAgICAgICAgICAgVGl0bGU6IFwiQ29uZmlybVwiLFxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcbiAgICAgICAgICAgIERldGFjaGVkOiBmYWxzZSxcbiAgICAgICAgICAgIEJ1dHRvbnM6IFtcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJOb1wiLCBJc0RlZmF1bHQ6IHRydWUgfVxuICAgICAgICAgICAgXVxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcnVuRWZmZWN0KCk7XG4gICAgfVxufVxuXG4vKipcbiAqIEB0eXBlIHtzeW1ib2x9XG4gKi9cbmNvbnN0IGNvbnRyb2xsZXIgPSBTeW1ib2woKTtcblxuLyoqXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXG4gKiBpdCB0aWVzIHRoZW0gdG8gYW4gQWJvcnRTaWduYWwgYW5kIHVzZXMgYW4gQWJvcnRDb250cm9sbGVyIHRvIHJlbW92ZSB0aGVtIGFsbCBhdCBvbmNlLlxuICovXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XG4gICAgY29uc3RydWN0b3IoKSB7XG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTdG9yZXMgdGhlIEFib3J0Q29udHJvbGxlciB0aGF0IGNhbiBiZSB1c2VkIHRvIHJlbW92ZSBhbGwgY3VycmVudGx5IGFjdGl2ZSBsaXN0ZW5lcnMuXG4gICAgICAgICAqXG4gICAgICAgICAqIEBwcml2YXRlXG4gICAgICAgICAqIEBuYW1lIHtAbGluayBjb250cm9sbGVyfVxuICAgICAgICAgKiBAbWVtYmVyIHtBYm9ydENvbnRyb2xsZXJ9XG4gICAgICAgICAqL1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJdID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgYW4gb3B0aW9ucyBvYmplY3QgZm9yIGFkZEV2ZW50TGlzdGVuZXIgdGhhdCB0aWVzIHRoZSBsaXN0ZW5lclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB7SFRNTEVsZW1lbnR9IGVsZW1lbnQgQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHtzdHJpbmdbXX0gdHJpZ2dlcnMgVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xuICAgICAqIEByZXR1cm5zIHtBZGRFdmVudExpc3RlbmVyT3B0aW9uc31cbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudCwgdHJpZ2dlcnMpIHtcbiAgICAgICAgcmV0dXJuIHsgc2lnbmFsOiB0aGlzW2NvbnRyb2xsZXJdLnNpZ25hbCB9O1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzLlxuICAgICAqXG4gICAgICogQHJldHVybnMge3ZvaWR9XG4gICAgICovXG4gICAgcmVzZXQoKSB7XG4gICAgICAgIHRoaXNbY29udHJvbGxlcl0uYWJvcnQoKTtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG59XG5cbi8qKlxuICogQHR5cGUge3N5bWJvbH1cbiAqL1xuY29uc3QgdHJpZ2dlck1hcCA9IFN5bWJvbCgpO1xuXG4vKipcbiAqIEB0eXBlIHtzeW1ib2x9XG4gKi9cbmNvbnN0IGVsZW1lbnRDb3VudCA9IFN5bWJvbCgpO1xuXG4vKipcbiAqIFdlYWtNYXBSZWdpc3RyeSBtYXBzIGFjdGl2ZSB0cmlnZ2VyIGV2ZW50cyB0byBlYWNoIERPTSBlbGVtZW50IHRocm91Z2ggYSBXZWFrTWFwLlxuICogVGhpcyBlbnN1cmVzIHRoYXQgdGhlIG1hcHBpbmcgcmVtYWlucyBwcml2YXRlIHRvIHRoaXMgbW9kdWxlLCB3aGlsZSBzdGlsbCBhbGxvd2luZyBnYXJiYWdlXG4gKiBjb2xsZWN0aW9uIG9mIHRoZSBpbnZvbHZlZCBlbGVtZW50cy5cbiAqL1xuY2xhc3MgV2Vha01hcFJlZ2lzdHJ5IHtcbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgLyoqXG4gICAgICAgICAqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy5cbiAgICAgICAgICpcbiAgICAgICAgICogQHByaXZhdGVcbiAgICAgICAgICogQG5hbWUge0BsaW5rIHRyaWdnZXJNYXB9XG4gICAgICAgICAqIEBtZW1iZXIge1dlYWtNYXA8SFRNTEVsZW1lbnQsIHN0cmluZ1tdPn1cbiAgICAgICAgICovXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcF0gPSBuZXcgV2Vha01hcCgpO1xuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBDb3VudHMgdGhlIG51bWJlciBvZiBlbGVtZW50cyB3aXRoIGFjdGl2ZSBXTUwgdHJpZ2dlcnMuXG4gICAgICAgICAqXG4gICAgICAgICAqIEBwcml2YXRlXG4gICAgICAgICAqIEBuYW1lIHtAbGluayBlbGVtZW50Q291bnR9XG4gICAgICAgICAqIEBtZW1iZXIge251bWJlcn1cbiAgICAgICAgICovXG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50XSA9IDA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgYWN0aXZlIHRyaWdnZXJzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXG4gICAgICpcbiAgICAgKiBAcGFyYW0ge0hUTUxFbGVtZW50fSBlbGVtZW50IEFuIEhUTUwgZWxlbWVudFxuICAgICAqIEBwYXJhbSB7c3RyaW5nW119IHRyaWdnZXJzIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudFxuICAgICAqIEByZXR1cm5zIHtBZGRFdmVudExpc3RlbmVyT3B0aW9uc31cbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudCwgdHJpZ2dlcnMpIHtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRdICs9ICF0aGlzW3RyaWdnZXJNYXBdLmhhcyhlbGVtZW50KTtcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwXS5zZXQoZWxlbWVudCwgdHJpZ2dlcnMpO1xuICAgICAgICByZXR1cm4ge307XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyB7dm9pZH1cbiAgICAgKi9cbiAgICByZXNldCgpIHtcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50XSA8PSAwKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xuICAgICAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50XSA8PSAwKVxuICAgICAgICAgICAgICAgIGJyZWFrO1xuXG4gICAgICAgICAgICBjb25zdCB0cmlnZ2VycyA9IHRoaXNbdHJpZ2dlck1hcF0uZ2V0KGVsZW1lbnQpO1xuICAgICAgICAgICAgdGhpc1tlbGVtZW50Q291bnRdIC09ICh0eXBlb2YgdHJpZ2dlcnMgIT09IFwidW5kZWZpbmVkXCIpO1xuXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXG4gICAgICAgICAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkKTtcbiAgICAgICAgfVxuXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcF0gPSBuZXcgV2Vha01hcCgpO1xuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudF0gPSAwO1xuICAgIH1cbn1cblxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xuXG4vKipcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAqXG4gKiBAcGFyYW0ge0hUTUxFbGVtZW50fSBlbGVtZW50XG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5mdW5jdGlvbiBhZGRXTUxMaXN0ZW5lcnMoZWxlbWVudCkge1xuICAgIGNvbnN0IHRyaWdnZXJSZWdFeHAgPSAvXFxTKy9nO1xuICAgIGNvbnN0IHRyaWdnZXJBdHRyID0gKGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIik7XG4gICAgY29uc3QgdHJpZ2dlcnMgPSBbXTtcblxuICAgIGxldCBtYXRjaDtcbiAgICB3aGlsZSAoKG1hdGNoID0gdHJpZ2dlclJlZ0V4cC5leGVjKHRyaWdnZXJBdHRyKSkgIT09IG51bGwpXG4gICAgICAgIHRyaWdnZXJzLnB1c2gobWF0Y2hbMF0pO1xuXG4gICAgY29uc3Qgb3B0aW9ucyA9IHRyaWdnZXJSZWdpc3RyeS5zZXQoZWxlbWVudCwgdHJpZ2dlcnMpO1xuICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycylcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTY2hlZHVsZXMgYW4gYXV0b21hdGljIHJlbG9hZCBvZiBXTUwgdG8gYmUgcGVyZm9ybWVkIGFzIHNvb24gYXMgdGhlIGRvY3VtZW50IGlzIGZ1bGx5IGxvYWRlZC5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gRW5hYmxlKCkge1xuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xufVxuXG4vKipcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cbiAqXG4gKiBAcmV0dXJuIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCkge1xuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xuICAgIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0sIFt3bWwtd2luZG93XSwgW3dtbC1vcGVudXJsXScpLmZvckVhY2goYWRkV01MTGlzdGVuZXJzKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5Ccm93c2VyLCAnJyk7XG5jb25zdCBCcm93c2VyT3BlblVSTCA9IDA7XG5cbi8qKlxuICogT3BlbiBhIGJyb3dzZXIgd2luZG93IHRvIHRoZSBnaXZlbiBVUkxcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybCkge1xuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsfSk7XG59XG4iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxuICAgIENhbmNlbENhbGw6IDEwLFxufVxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG4vKipcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cbiAqXG4gKiBAcGFyYW0ge09iamVjdH0gb2JqZWN0IC0gVGhlIG9iamVjdCBvbiB3aGljaCB0aGUgbWV0aG9kIGlzIHRvIGJlIGludm9rZWQuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyB7RnVuY3Rpb259IEEgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24gdGhhdCB0YWtlcyB0aGUgbWV0aG9kIG5hbWUgYW5kIG9wdGlvbmFsbHkgYXJndW1lbnRzIGFuZCBpbnZva2VzIHRoZSBtZXRob2Qgd2l0aGluIHRoZSBzcGVjaWZpZWQgd2luZG93IGNvbnRleHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG4vKipcbiAqIENyZWF0ZXMgYSBuZXcgcnVudGltZSBjYWxsZXIgd2l0aCBzcGVjaWZpZWQgSUQuXG4gKlxuICogQHBhcmFtIHtvYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gKiBAcmV0dXJuIHtGdW5jdGlvbn0gLSBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3QsIHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICBpZiggbWV0aG9kICkge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIH1cbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcblxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm9iamVjdFwiLCBvYmplY3RJRCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gT3BlbkZpbGVEaWFsb2dPcHRpb25zXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ3JlYXRlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbU2hvd0hpZGRlbkZpbGVzXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c011bHRpcGxlU2VsZWN0aW9uXSAtIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0hpZGVFeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c090aGVyRmlsZXR5cGVzXSAtIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGl0bGUgb2YgdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTWVzc2FnZV0gLSBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlyZWN0b3J5XSAtIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTYXZlRmlsZURpYWxvZ09wdGlvbnNcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRmlsZW5hbWVdIC0gRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DcmVhdGVEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtTaG93SGlkZGVuRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtSZXNvbHZlc0FsaWFzZXNdIC0gSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb25dIC0gSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSGlkZUV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzT3RoZXJGaWxldHlwZXNdIC0gSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuXG4gKiBAcHJvcGVydHkge0ZpbGVGaWx0ZXJbXX0gW0ZpbHRlcnNdIC0gQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaXRsZSBvZiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtNZXNzYWdlXSAtIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtCdXR0b25UZXh0XSAtIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtEaXJlY3RvcnldIC0gRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0RldGFjaGVkXSAtIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge0J1dHRvbltdfSBbQnV0dG9uc10gLSBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBCdXR0b25cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTGFiZWxdIC0gVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtJc0NhbmNlbF0gLSBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNEZWZhdWx0XSAtIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IEZpbGVGaWx0ZXJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlzcGxheU5hbWVdIC0gRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1BhdHRlcm5dIC0gUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy5cbiAqL1xuXG4vLyBzZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpYWxvZ0Vycm9yQ2FsbGJhY2sgPSBkaWFsb2dFcnJvckNhbGxiYWNrO1xud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZywgJycpO1xuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBpZCB0aGF0IGlzIG5vdCBwcmVzZW50IGluIGRpYWxvZ1Jlc3BvbnNlcy5cbiAqIEByZXR1cm5zIHtzdHJpbmd9IHVuaXF1ZSBpZFxuICovXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogU2hvd3MgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7bnVtYmVyfSB0eXBlIC0gdHlwZSBvZiBkaWFsb2dcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN8T3BlbkZpbGVEaWFsb2dPcHRpb25zfFNhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIG9wdGlvbnMgZm9yIHRoZSBkaWFsb2dcbiAqIEByZXR1cm5zIHtQcm9taXNlfSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zID0ge30pIHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBjYWxsYmFjayBmcm9tIGEgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBJRCBvZiB0aGUgZGlhbG9nIHJlc3BvbnNlLlxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgZGF0YSByZWNlaXZlZCBmcm9tIHRoZSBkaWFsb2cuXG4gKiBAcGFyYW0ge2Jvb2xlYW59IGlzSlNPTiAtIEZsYWcgaW5kaWNhdGluZyB3aGV0aGVyIHRoZSBkYXRhIGlzIGluIEpTT04gZm9ybWF0LlxuICpcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cbiAqL1xuZnVuY3Rpb24gZGlhbG9nUmVzdWx0Q2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBDYWxsYmFjayBmdW5jdGlvbiBmb3IgaGFuZGxpbmcgZXJyb3JzIGluIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuXG4vLyBSZXBsYWNlIGBtZXRob2RzYCB3aXRoIGNvbnN0YW50cyBpbiBUaXRsZSBDYXNlXG5cbi8qKlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGNvbnN0IEluZm8gPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpO1xuXG4vKipcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBjb25zdCBXYXJuaW5nID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dXYXJuaW5nLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgY29uc3QgRXJyb3IgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgY29uc3QgUXVlc3Rpb24gPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge09wZW5GaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcy4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGNvbnN0IE9wZW5GaWxlID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dPcGVuRmlsZSwgb3B0aW9ucyk7XG5cbi8qKlxuICogQHBhcmFtIHtTYXZlRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZS4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGNvbnN0IFNhdmVGaWxlID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi90eXBlc1wiKS5XYWlsc0V2ZW50fSBXYWlsc0V2ZW50XG4gKi9cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHtFdmVudFR5cGVzfSBmcm9tIFwiLi9ldmVudF90eXBlc1wiO1xuZXhwb3J0IGNvbnN0IFR5cGVzID0gRXZlbnRUeXBlcztcblxuLy8gU2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkV2ZW50cywgJycpO1xuY29uc3QgRW1pdE1ldGhvZCA9IDA7XG5jb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXAoKTtcblxuY2xhc3MgTGlzdGVuZXIge1xuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xuICAgICAgICAgICAgY2FsbGJhY2soZGF0YSk7XG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQge1xuICAgIGNvbnN0cnVjdG9yKG5hbWUsIGRhdGEgPSBudWxsKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gc2V0dXAoKSB7XG59XG5cbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudCkge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XG4gICAgaWYgKGxpc3RlbmVycykge1xuICAgICAgICBsZXQgdG9SZW1vdmUgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+IHtcbiAgICAgICAgICAgIGxldCByZW1vdmUgPSBsaXN0ZW5lci5DYWxsYmFjayhldmVudCk7XG4gICAgICAgICAgICBpZiAocmVtb3ZlKSByZXR1cm4gdHJ1ZTtcbiAgICAgICAgfSk7XG4gICAgICAgIGlmICh0b1JlbW92ZS5sZW5ndGggPiAwKSB7XG4gICAgICAgICAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gIXRvUmVtb3ZlLmluY2x1ZGVzKGwpKTtcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgICAgICAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCBsaXN0ZW5lcnMpO1xuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzIC0gVGhlIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRoZSBjYWxsYmFjayBjYW4gYmUgY2FsbGVkIGZvciB0aGUgZXZlbnQuIE9uY2UgdGhlIG1heGltdW0gbnVtYmVyIGlzIHJlYWNoZWQsIHRoZSBjYWxsYmFjayB3aWxsIG5vIGxvbmdlciBiZSBjYWxsZWQuXG4gKlxuIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50LlxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQuIEl0IHRha2VzIG5vIHBhcmFtZXRlcnMuXG4gKiBAcmV0dXJuIHtmdW5jdGlvbn0gLSBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHsgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpOyB9XG5cbi8qKlxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgb25seSBvbmNlIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50LlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFjayAtIFRoZSBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCB3aGVuIHRoZSBldmVudCBvY2N1cnMuXG4gKiBAcmV0dXJuIHtmdW5jdGlvbn0gLSBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7IHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpOyB9XG5cbi8qKlxuICogUmVtb3ZlcyB0aGUgc3BlY2lmaWVkIGxpc3RlbmVyIGZyb20gdGhlIGV2ZW50IGxpc3RlbmVycyBjb2xsZWN0aW9uLlxuICogSWYgYWxsIGxpc3RlbmVycyBmb3IgdGhlIGV2ZW50IGFyZSByZW1vdmVkLCB0aGUgZXZlbnQga2V5IGlzIGRlbGV0ZWQgZnJvbSB0aGUgY29sbGVjdGlvbi5cbiAqXG4gKiBAcGFyYW0ge09iamVjdH0gbGlzdGVuZXIgLSBUaGUgbGlzdGVuZXIgdG8gYmUgcmVtb3ZlZC5cbiAqL1xuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgZWxzZSBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xufVxuXG5cbi8qKlxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cbiAqIEBwYXJhbSB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lcyAtIEFkZGl0aW9uYWwgZXZlbnQgbmFtZXMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIGxldCBldmVudHNUb1JlbW92ZSA9IFtldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzXTtcbiAgICBldmVudHNUb1JlbW92ZS5mb3JFYWNoKGV2ZW50TmFtZSA9PiBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKSk7XG59XG4vKipcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cbiAqXG4gKiBAZnVuY3Rpb24gT2ZmQWxsXG4gKiBAcmV0dXJucyB7dm9pZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpIHsgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTsgfVxuXG4vKipcbiAqIEVtaXRzIGFuIGV2ZW50IHVzaW5nIHRoZSBnaXZlbiBldmVudCBuYW1lLlxuICpcbiAqIEBwYXJhbSB7V2FpbHNFdmVudH0gZXZlbnQgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdC5cbiAqIEByZXR1cm5zIHthbnl9IC0gVGhlIHJlc3VsdCBvZiB0aGUgZW1pdHRlZCBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVtaXQoZXZlbnQpIHsgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgZXZlbnQpOyB9XG4iLCAiXG5leHBvcnQgY29uc3QgRXZlbnRUeXBlcyA9IHtcblx0V2luZG93czoge1xuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdEFQTVBvd2VyU3RhdHVzQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTdGF0dXNDaGFuZ2VcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFQTVJlc3VtZUF1dG9tYXRpYzogXCJ3aW5kb3dzOkFQTVJlc3VtZUF1dG9tYXRpY1wiLFxuXHRcdEFQTVJlc3VtZVN1c3BlbmQ6IFwid2luZG93czpBUE1SZXN1bWVTdXNwZW5kXCIsXG5cdFx0QVBNUG93ZXJTZXR0aW5nQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTZXR0aW5nQ2hhbmdlXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcIndpbmRvd3M6QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWQ6IFwid2luZG93czpXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZFwiLFxuXHRcdFdpbmRvd0luYWN0aXZlOiBcIndpbmRvd3M6V2luZG93SW5hY3RpdmVcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd0Nsb3NlOiBcIndpbmRvd3M6V2luZG93Q2xvc2VcIixcblx0XHRXaW5kb3dTZXRGb2N1czogXCJ3aW5kb3dzOldpbmRvd1NldEZvY3VzXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0fSxcblx0TWFjOiB7XG5cdFx0QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUljb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0JhckZyYW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0JhckZyYW1lXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVOb3RpZmljYXRpb246IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlTm90aWZpY2F0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlVGhlbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlVGhlbWUhXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuIVwiLFxuXHRcdFdpbmRvd0RpZEJlY29tZUtleTogXCJtYWM6V2luZG93RGlkQmVjb21lS2V5XCIsXG5cdFx0V2luZG93RGlkQmVjb21lTWFpbjogXCJtYWM6V2luZG93RGlkQmVjb21lTWFpblwiLFxuXHRcdFdpbmRvd0RpZEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd0RpZEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VBbHBoYTogXCJtYWM6V2luZG93RGlkQ2hhbmdlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nTG9jYXRpb246IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVGl0bGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkQ2hhbmdlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVZpc2liaWxpdHk6IFwibWFjOldpbmRvd0RpZENoYW5nZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dEaWREZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWREZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkRW5kU2hlZXQ6IFwibWFjOldpbmRvd0RpZEVuZFNoZWV0XCIsXG5cdFx0V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxuXHRcdFdpbmRvd0RpZEZvY3VzOiBcIm1hYzpXaW5kb3dEaWRGb2N1c1wiLFxuXHRcdFdpbmRvd0RpZE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWRNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPZmZTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcblx0XHRXaW5kb3dEaWRSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIm1hYzpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93RGlkVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93RGlkVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZSFcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRNZW51V2lsbE9wZW46IFwibWFjOk1lbnVXaWxsT3BlblwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvbjogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25cIixcblx0XHRNZW51V2lsbEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVXaWxsSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVXaWxsRGlzcGxheUl0ZW06IFwibWFjOk1lbnVXaWxsRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRGlzcGxheUl0ZW06IFwibWFjOk1lbnVEaWREaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudURpZEFkZEl0ZW06IFwibWFjOk1lbnVEaWRBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51RGlkUmVtb3ZlSXRlbTogXCJtYWM6TWVudURpZFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsRW5kVHJhY2tpbmdcIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsVXBkYXRlOiBcIm1hYzpNZW51V2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRVcGRhdGU6IFwibWFjOk1lbnVEaWRVcGRhdGVcIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudURpZFBvcFVwOiBcIm1hYzpNZW51RGlkUG9wVXBcIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbVwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblRvSXRlbVwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvblwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRXhpdGVkXCIsXG5cdH0sXG5cdExpbnV4OiB7XG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcImxpbnV4OlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0xvYWRDaGFuZ2VkOiBcImxpbnV4OldpbmRvd0xvYWRDaGFuZ2VkXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dGb2N1c0luOiBcImxpbnV4OldpbmRvd0ZvY3VzSW5cIixcblx0XHRXaW5kb3dGb2N1c091dDogXCJsaW51eDpXaW5kb3dGb2N1c091dFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcblx0fSxcblx0Q29tbW9uOiB7XG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcImNvbW1vbjpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJjb21tb246V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93Wm9vbTogXCJjb21tb246V2luZG93Wm9vbVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJjb21tb246V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJjb21tb246V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJjb21tb246V2luZG93Wm9vbVJlc2V0XCIsXG5cdFx0V2luZG93Rm9jdXM6IFwiY29tbW9uOldpbmRvd0ZvY3VzXCIsXG5cdFx0V2luZG93TG9zdEZvY3VzOiBcImNvbW1vbjpXaW5kb3dMb3N0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcImNvbW1vbjpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93SGlkZTogXCJjb21tb246V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwiY29tbW9uOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dGaWxlc0Ryb3BwZWQ6IFwiY29tbW9uOldpbmRvd0ZpbGVzRHJvcHBlZFwiLFxuXHRcdFdpbmRvd1J1bnRpbWVSZWFkeTogXCJjb21tb246V2luZG93UnVudGltZVJlYWR5XCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImNvbW1vbjpUaGVtZUNoYW5nZWRcIixcblx0fSxcbn07XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogTG9ncyBhIG1lc3NhZ2UgdG8gdGhlIGNvbnNvbGUgd2l0aCBjdXN0b20gZm9ybWF0dGluZy5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIG1lc3NhZ2UgdG8gYmUgbG9nZ2VkLlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGRlYnVnTG9nKG1lc3NhZ2UpIHtcbiAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICBjb25zb2xlLmxvZyhcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNhYTAwMDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAzcHggMHB4IDBweCAzcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nLFxuICAgICAgICAnYmFja2dyb3VuZDogIzAwOTkwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDBweCAzcHggM3B4IDBweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbSdcbiAgICApO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXG4gKiAoc2VlIGh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnL2VuLVVTL2RvY3MvV2ViL0FQSS9FdmVudFRhcmdldC9hZGRFdmVudExpc3RlbmVyI3NpZ25hbClcbiAqXG4gKiBAcmV0dXJuIHtib29sZWFufVxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XG4gICAgaWYgKCFFdmVudFRhcmdldCB8fCAhQWJvcnRTaWduYWwgfHwgIUFib3J0Q29udHJvbGxlcilcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuXG4gICAgbGV0IHJlc3VsdCA9IHRydWU7XG5cbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcbiAgICBjb25zdCBjb250cm9sbGVyID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcbiAgICB0YXJnZXQuZGlzcGF0Y2hFdmVudChuZXcgQ3VzdG9tRXZlbnQoJ3Rlc3QnKSk7XG5cbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKioqXG4gVGhpcyB0ZWNobmlxdWUgZm9yIHByb3BlciBsb2FkIGRldGVjdGlvbiBpcyB0YWtlbiBmcm9tIEhUTVg6XG5cbiBCU0QgMi1DbGF1c2UgTGljZW5zZVxuXG4gQ29weXJpZ2h0IChjKSAyMDIwLCBCaWcgU2t5IFNvZnR3YXJlXG4gQWxsIHJpZ2h0cyByZXNlcnZlZC5cblxuIFJlZGlzdHJpYnV0aW9uIGFuZCB1c2UgaW4gc291cmNlIGFuZCBiaW5hcnkgZm9ybXMsIHdpdGggb3Igd2l0aG91dFxuIG1vZGlmaWNhdGlvbiwgYXJlIHBlcm1pdHRlZCBwcm92aWRlZCB0aGF0IHRoZSBmb2xsb3dpbmcgY29uZGl0aW9ucyBhcmUgbWV0OlxuXG4gMS4gUmVkaXN0cmlidXRpb25zIG9mIHNvdXJjZSBjb2RlIG11c3QgcmV0YWluIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLCB0aGlzXG4gbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIuXG5cbiAyLiBSZWRpc3RyaWJ1dGlvbnMgaW4gYmluYXJ5IGZvcm0gbXVzdCByZXByb2R1Y2UgdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsXG4gdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvblxuIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLlxuXG4gVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyBcIkFTIElTXCJcbiBBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFXG4gSU1QTElFRCBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSBBTkQgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQVJFXG4gRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1JTIEJFIExJQUJMRVxuIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMXG4gREFNQUdFUyAoSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFBST0NVUkVNRU5UIE9GIFNVQlNUSVRVVEUgR09PRFMgT1JcbiBTRVJWSUNFUzsgTE9TUyBPRiBVU0UsIERBVEEsIE9SIFBST0ZJVFM7IE9SIEJVU0lORVNTIElOVEVSUlVQVElPTikgSE9XRVZFUlxuIENBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksXG4gT1IgVE9SVCAoSU5DTFVESU5HIE5FR0xJR0VOQ0UgT1IgT1RIRVJXSVNFKSBBUklTSU5HIElOIEFOWSBXQVkgT1VUIE9GIFRIRSBVU0VcbiBPRiBUSElTIFNPRlRXQVJFLCBFVkVOIElGIEFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLlxuXG4gKioqL1xuXG5sZXQgaXNSZWFkeSA9IGZhbHNlO1xuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsICgpID0+IGlzUmVhZHkgPSB0cnVlKTtcblxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjaykge1xuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcbiAgICAgICAgY2FsbGJhY2soKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgY2FsbGJhY2spO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vLyBJbXBvcnQgc2NyZWVuIGpzZG9jIGRlZmluaXRpb24gZnJvbSAuL3NjcmVlbnMuanNcbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vc2NyZWVuc1wiKS5TY3JlZW59IFNjcmVlblxuICovXG5cblxuLyoqXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBwb3NpdGlvbiBvZiBhIHdpbmRvdy5cbiAqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBQb3NpdGlvblxuICogQHByb3BlcnR5IHtudW1iZXJ9IHggLSBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKiBAcHJvcGVydHkge251bWJlcn0geSAtIFRoZSB2ZXJ0aWNhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93XG4gKi9cblxuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXG4gKlxuICogQHR5cGVkZWYge09iamVjdH0gU2l6ZVxuICogQHByb3BlcnR5IHtudW1iZXJ9IHdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3dcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3dcbiAqL1xuXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuY29uc3QgQWJzb2x1dGVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMDtcbmNvbnN0IENlbnRlck1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDE7XG5jb25zdCBDbG9zZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSAyO1xuY29uc3QgRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgID0gMztcbmNvbnN0IEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgICA9IDQ7XG5jb25zdCBGb2N1c01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1O1xuY29uc3QgRm9yY2VSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgID0gNjtcbmNvbnN0IEZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICA9IDc7XG5jb25zdCBHZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA4O1xuY29uc3QgR2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gOTtcbmNvbnN0IEhlaWdodE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDEwO1xuY29uc3QgSGlkZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTE7XG5jb25zdCBJc0ZvY3VzZWRNZXRob2QgICAgICAgICAgICAgICAgICAgPSAxMjtcbmNvbnN0IElzRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDEzO1xuY29uc3QgSXNNYXhpbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTQ7XG5jb25zdCBJc01pbmltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNTtcbmNvbnN0IE1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE2O1xuY29uc3QgTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTc7XG5jb25zdCBOYW1lTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxODtcbmNvbnN0IE9wZW5EZXZUb29sc01ldGhvZCAgICAgICAgICAgICAgICA9IDE5O1xuY29uc3QgUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMjA7XG5jb25zdCBSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAyMTtcbmNvbnN0IFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDIyO1xuY29uc3QgUmVzdG9yZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMjM7XG5jb25zdCBTZXRBYnNvbHV0ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAyNDtcbmNvbnN0IFNldEFsd2F5c09uVG9wTWV0aG9kICAgICAgICAgICAgICA9IDI1O1xuY29uc3QgU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCAgICAgICAgID0gMjY7XG5jb25zdCBTZXRGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgICAgPSAyNztcbmNvbnN0IFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kICA9IDI4O1xuY29uc3QgU2V0TWF4U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMjk7XG5jb25zdCBTZXRNaW5TaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAzMDtcbmNvbnN0IFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICA9IDMxO1xuY29uc3QgU2V0UmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgID0gMzI7XG5jb25zdCBTZXRTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzMztcbmNvbnN0IFNldFRpdGxlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDM0O1xuY29uc3QgU2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzU7XG5jb25zdCBTaG93TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNjtcbmNvbnN0IFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM3O1xuY29uc3QgVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgID0gMzg7XG5jb25zdCBUb2dnbGVNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgPSAzOTtcbmNvbnN0IFVuRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDQwO1xuY29uc3QgVW5NYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDE7XG5jb25zdCBVbk1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MjtcbmNvbnN0IFdpZHRoTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDQzO1xuY29uc3QgWm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XG5jb25zdCBab29tSW5NZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSA0NTtcbmNvbnN0IFpvb21PdXRNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDQ2O1xuY29uc3QgWm9vbVJlc2V0TWV0aG9kICAgICAgICAgICAgICAgICAgID0gNDc7XG5cbi8qKlxuICogQHR5cGUge3N5bWJvbH1cbiAqL1xuY29uc3QgY2FsbGVyID0gU3ltYm9sKCk7XG5cbmNsYXNzIFdpbmRvdyB7XG4gICAgLyoqXG4gICAgICogSW5pdGlhbGlzZXMgYSB3aW5kb3cgb2JqZWN0IHdpdGggdGhlIHNwZWNpZmllZCBuYW1lLlxuICAgICAqXG4gICAgICogQHByaXZhdGVcbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB0YXJnZXQgd2luZG93LlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG5hbWUgPSAnJykge1xuICAgICAgICAvKipcbiAgICAgICAgICogQHByaXZhdGVcbiAgICAgICAgICogQG5hbWUge0BsaW5rIGNhbGxlcn1cbiAgICAgICAgICogQHR5cGUgeyguLi5hcmdzOiBhbnlbXSkgPT4gYW55fVxuICAgICAgICAgKi9cbiAgICAgICAgdGhpc1tjYWxsZXJdID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5XaW5kb3csIG5hbWUpXG5cbiAgICAgICAgLy8gYmluZCBpbnN0YW5jZSBtZXRob2QgdG8gbWFrZSB0aGVtIGVhc2lseSB1c2FibGUgaW4gZXZlbnQgaGFuZGxlcnNcbiAgICAgICAgZm9yIChjb25zdCBtZXRob2Qgb2YgT2JqZWN0LmdldE93blByb3BlcnR5TmFtZXMoV2luZG93LnByb3RvdHlwZSkpIHtcbiAgICAgICAgICAgIGlmIChcbiAgICAgICAgICAgICAgICBtZXRob2QgIT09IFwiY29uc3RydWN0b3JcIlxuICAgICAgICAgICAgICAgICYmIHR5cGVvZiB0aGlzW21ldGhvZF0gPT09IFwiZnVuY3Rpb25cIlxuICAgICAgICAgICAgKSB7XG4gICAgICAgICAgICAgICAgdGhpc1ttZXRob2RdID0gdGhpc1ttZXRob2RdLmJpbmQodGhpcyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBHZXRzIHRoZSBzcGVjaWZpZWQgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBnZXQuXG4gICAgICogQHJldHVybiB7V2luZG93fSAtIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXG4gICAgICovXG4gICAgR2V0KG5hbWUpIHtcbiAgICAgICAgcmV0dXJuIG5ldyBXaW5kb3cobmFtZSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8UG9zaXRpb24+fSAtIFRoZSBjdXJyZW50IGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgQWJzb2x1dGVQb3NpdGlvbigpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShBYnNvbHV0ZVBvc2l0aW9uTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIENlbnRlcigpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShDZW50ZXJNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENsb3NlcyB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgQ2xvc2UoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oQ2xvc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIERpc2FibGVTaXplQ29uc3RyYWludHMoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRW5hYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBFbmFibGVTaXplQ29uc3RyYWludHMoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBGb2N1c2VzIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBGb2N1cygpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShGb2N1c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIEZvcmNlUmVsb2FkKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEZvcmNlUmVsb2FkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEb2MuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBGdWxsc2NyZWVuKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxTY3JlZW4+fSAtIFRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgb25cbiAgICAgKi9cbiAgICBHZXRTY3JlZW4oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oR2V0U2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBjdXJyZW50IHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPG51bWJlcj59IC0gVGhlIGN1cnJlbnQgem9vbSBsZXZlbFxuICAgICAqL1xuICAgIEdldFpvb20oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oR2V0Wm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxudW1iZXI+fSAtIFRoZSBjdXJyZW50IGhlaWdodCBvZiB0aGUgd2luZG93XG4gICAgICovXG4gICAgSGVpZ2h0KCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEhlaWdodE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSGlkZXMgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIEhpZGUoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSGlkZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZm9jdXNlZC5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZm9jdXNlZFxuICAgICAqL1xuICAgIElzRm9jdXNlZCgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShJc0ZvY3VzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gLSBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW5cbiAgICAgKi9cbiAgICBJc0Z1bGxzY3JlZW4oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSXNGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gLSBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1heGltaXNlZFxuICAgICAqL1xuICAgIElzTWF4aW1pc2VkKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKElzTWF4aW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtaW5pbWlzZWQuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gLSBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1pbmltaXNlZFxuICAgICAqL1xuICAgIElzTWluaW1pc2VkKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKElzTWluaW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIE1heGltaXNlKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKE1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNaW5pbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIE1pbmltaXNlKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKE1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3dcbiAgICAgKi9cbiAgICBOYW1lKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKE5hbWVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgT3BlbkRldlRvb2xzKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKE9wZW5EZXZUb29sc01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8UG9zaXRpb24+fSAtIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKi9cbiAgICBSZWxhdGl2ZVBvc2l0aW9uKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFJlbGF0aXZlUG9zaXRpb25NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbG9hZHMgdGhlIHBhZ2UgYXNzZXRzLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgUmVsb2FkKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IC0gV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGVcbiAgICAgKi9cbiAgICBSZXNpemFibGUoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oUmVzaXphYmxlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXN0b3JlcyB0aGUgd2luZG93IHRvIGl0cyBwcmV2aW91cyBzdGF0ZSBpZiBpdCB3YXMgcHJldmlvdXNseSBtaW5pbWlzZWQsIG1heGltaXNlZCBvciBmdWxsc2NyZWVuLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgUmVzdG9yZSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShSZXN0b3JlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcGFyYW0ge251bWJlcn0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldEFic29sdXRlUG9zaXRpb24oeCwgeSkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldEFic29sdXRlUG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB3aW5kb3cgdG8gYmUgYWx3YXlzIG9uIHRvcC5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcGFyYW0ge2Jvb2xlYW59IGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcFxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0QWx3YXlzT25Ub3AoYWx3YXlzT25Ub3ApIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRBbHdheXNPblRvcE1ldGhvZCwgeyBhbHdheXNPblRvcCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSByIC0gVGhlIGRlc2lyZWQgcmVkIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmRcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gZyAtIFRoZSBkZXNpcmVkIGdyZWVuIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmRcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZFxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBhIC0gVGhlIGRlc2lyZWQgYWxwaGEgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZFxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0QmFja2dyb3VuZENvbG91cihyLCBnLCBiLCBhKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCwgeyByLCBnLCBiLCBhIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgdGhlIHdpbmRvdyBmcmFtZSBhbmQgdGl0bGUgYmFyLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3NcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldEZyYW1lbGVzcyhmcmFtZWxlc3MpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRGcmFtZWxlc3NNZXRob2QsIHsgZnJhbWVsZXNzIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIHRoZSBzeXN0ZW0gZnVsbHNjcmVlbiBidXR0b24uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtib29sZWFufSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWRcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkKGVuYWJsZWQpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCwgeyBlbmFibGVkIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1heGltdW0gd2lkdGggb2YgdGhlIHdpbmRvd1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtYXhpbXVtIGhlaWdodCBvZiB0aGUgd2luZG93XG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBTZXRNYXhTaXplKHdpZHRoLCBoZWlnaHQpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRNYXhTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIGRlc2lyZWQgbWluaW11bSB3aWR0aCBvZiB0aGUgd2luZG93XG4gICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldE1pblNpemUod2lkdGgsIGhlaWdodCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldE1pblNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcGFyYW0ge251bWJlcn0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldFJlbGF0aXZlUG9zaXRpb24oeCwgeSkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtib29sZWFufSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZVxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0UmVzaXphYmxlKHJlc2l6YWJsZSkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldFJlc2l6YWJsZU1ldGhvZCwgeyByZXNpemFibGUgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCAtIFRoZSBkZXNpcmVkIHdpZHRoIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgaGVpZ2h0IG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldFNpemUod2lkdGgsIGhlaWdodCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB0aXRsZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZSAtIFRoZSBkZXNpcmVkIHRpdGxlIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldFRpdGxlKHRpdGxlKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0VGl0bGVNZXRob2QsIHsgdGl0bGUgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSB6b29tIC0gVGhlIGRlc2lyZWQgem9vbSBsZXZlbFxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0Wm9vbSh6b29tKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0Wm9vbU1ldGhvZCwgeyB6b29tIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNob3dzIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBTaG93KCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNob3dNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPFNpemU+fSAtIFRoZSBjdXJyZW50IHNpemUgb2YgdGhlIHdpbmRvd1xuICAgICAqL1xuICAgIFNpemUoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2l6ZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnVsbHNjcmVlbiBhbmQgbm9ybWFsLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgVG9nZ2xlRnVsbHNjcmVlbigpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShUb2dnbGVGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBtYXhpbWlzZWQgYW5kIG5vcm1hbC5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFRvZ2dsZU1heGltaXNlKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgVW5GdWxsc2NyZWVuKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFVuRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBVbk1heGltaXNlKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFVuTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgVW5NaW5pbWlzZSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShVbk1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8bnVtYmVyPn0gLSBUaGUgY3VycmVudCB3aWR0aCBvZiB0aGUgd2luZG93XG4gICAgICovXG4gICAgV2lkdGgoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oV2lkdGhNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFpvb21zIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBab29tKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFpvb21NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEluY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgWm9vbUluKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFpvb21Jbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBab29tT3V0KCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFpvb21PdXRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlc2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgWm9vbVJlc2V0KCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFpvb21SZXNldE1ldGhvZCk7XG4gICAgfVxufVxuXG4vKipcbiAqIFRoZSB3aW5kb3cgd2l0aGluIHdoaWNoIHRoZSBzY3JpcHQgaXMgcnVubmluZy5cbiAqXG4gKiBAdHlwZSB7V2luZG93fVxuICovXG5jb25zdCB0aGlzV2luZG93ID0gbmV3IFdpbmRvdygnJyk7XG5cbmV4cG9ydCBkZWZhdWx0IHRoaXNXaW5kb3c7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XG5cbi8vIE5PVEU6IHRoZSBmb2xsb3dpbmcgbWV0aG9kcyBNVVNUIGJlIGltcG9ydGVkIGV4cGxpY2l0bHkgYmVjYXVzZSBvZiBob3cgZXNidWlsZCBpbmplY3Rpb24gd29ya3NcbmltcG9ydCB7RW5hYmxlIGFzIEVuYWJsZVdNTH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dtbFwiO1xuaW1wb3J0IHtkZWJ1Z0xvZ30gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzXCI7XG5cbndpbmRvdy53YWlscyA9IFJ1bnRpbWU7XG5FbmFibGVXTUwoKTtcblxuaWYgKERFQlVHKSB7XG4gICAgZGVidWdMb2coXCJXYWlscyBSdW50aW1lIExvYWRlZFwiKVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TeXN0ZW0sICcnKTtcbmNvbnN0IHN5c3RlbUlzRGFya01vZGUgPSAwO1xuY29uc3QgZW52aXJvbm1lbnQgPSAxO1xuXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZykge1xuICAgIGlmKHdpbmRvdy5jaHJvbWUpIHtcbiAgICAgICAgcmV0dXJuIHdpbmRvdy5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZShtc2cpO1xuICAgIH1cbiAgICByZXR1cm4gd2luZG93LndlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2UobXNnKTtcbn1cblxuLyoqXG4gKiBAZnVuY3Rpb25cbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxib29sZWFuPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoc3lzdGVtSXNEYXJrTW9kZSk7XG59XG5cbi8qKlxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXG4gKlxuICogQGFzeW5jXG4gKiBAZnVuY3Rpb24gQ2FwYWJpbGl0aWVzXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxPYmplY3Q+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyB0aGUgY2FwYWJpbGl0aWVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FwYWJpbGl0aWVzKCkge1xuICAgIGxldCByZXNwb25zZSA9IGZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKTtcbiAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xufVxuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IE9TSW5mb1xuICogQHByb3BlcnR5IHtzdHJpbmd9IEJyYW5kaW5nIC0gVGhlIGJyYW5kaW5nIG9mIHRoZSBPUy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJRCAtIFRoZSBJRCBvZiB0aGUgT1MuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBPUy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBWZXJzaW9uIC0gVGhlIHZlcnNpb24gb2YgdGhlIE9TLlxuICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gRW52aXJvbm1lbnRJbmZvXG4gKiBAcHJvcGVydHkge3N0cmluZ30gQXJjaCAtIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gRGVidWcgLSBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBPUyAtIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS5cbiAqIEBwcm9wZXJ0eSB7T1NJbmZvfSBPU0luZm8gLSBEZXRhaWxzIG9mIHRoZSBvcGVyYXRpbmcgc3lzdGVtLlxuICogQHByb3BlcnR5IHtPYmplY3R9IFBsYXRmb3JtSW5mbyAtIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uXG4gKi9cblxuLyoqXG4gKiBAZnVuY3Rpb25cbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxuICogQHJldHVybnMge1Byb21pc2U8RW52aXJvbm1lbnRJbmZvPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyBPUyBhbmQgc3lzdGVtIGFyY2hpdGVjdHVyZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xuICAgIHJldHVybiBjYWxsKGVudmlyb25tZW50KTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxuICpcbiAqIEByZXR1cm4ge2Jvb2xlYW59IFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcIndpbmRvd3NcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwibGludXhcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxuICpcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBtYWNPUywgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNNYWMoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwiZGFyd2luXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFtZDY0XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXG4gKlxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBUk0oKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm1cIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxuICpcbiAqIEByZXR1cm5zIHtib29sZWFufSAtIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm02NFwiO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5EZWJ1ZyA9PT0gdHJ1ZTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcbmltcG9ydCB7SXNEZWJ1Z30gZnJvbSBcIi4vc3lzdGVtXCI7XG5cbi8vIHNldHVwXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5Db250ZXh0TWVudSwgJycpO1xuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkLCB4LCB5LCBkYXRhKSB7XG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XG59XG5cbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudCkge1xuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XG4gICAgbGV0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgbGV0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpO1xuICAgIGN1c3RvbUNvbnRleHRNZW51ID0gY3VzdG9tQ29udGV4dE1lbnUgPyBjdXN0b21Db250ZXh0TWVudS50cmltKCkgOiBcIlwiO1xuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnVEYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgY3VzdG9tQ29udGV4dE1lbnVEYXRhKTtcbiAgICAgICAgcmV0dXJuXG4gICAgfVxuXG4gICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCk7XG59XG5cblxuLypcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG5cblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXG4qL1xuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCkge1xuXG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XG4gICAgaWYgKElzRGVidWcoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcbiAgICAgICAgY2FzZSBcInNob3dcIjpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQuaXNDb250ZW50RWRpdGFibGUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWRcbiAgICAgICAgICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICAgICAgICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IChzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKVxuICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbikge1xuICAgICAgICAgICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0cyA9IHJhbmdlLmdldENsaWVudFJlY3RzKCk7XG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHJlY3QubGVmdCwgcmVjdC50b3ApID09PSBlbGVtZW50KSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGFnbmFtZSBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQudGFnTmFtZSA9PT0gXCJJTlBVVFwiIHx8IGVsZW1lbnQudGFnTmFtZSA9PT0gXCJURVhUQVJFQVwiKSB7XG4gICAgICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIWVsZW1lbnQucmVhZE9ubHkgJiYgIWVsZW1lbnQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGtleVN0cmluZyAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cbiAqIEByZXR1cm4geyp9IC0gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5U3RyaW5nKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZmxhZ3Nba2V5U3RyaW5nXTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleVN0cmluZyArIFwiJzogXCIgKyBlKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQge2ludm9rZSwgSXNXaW5kb3dzfSBmcm9tIFwiLi9zeXN0ZW1cIjtcbmltcG9ydCB7R2V0RmxhZ30gZnJvbSBcIi4vZmxhZ3NcIjtcblxuLy8gU2V0dXBcbmxldCBzaG91bGREcmFnID0gZmFsc2U7XG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5sZXQgcmVzaXplRWRnZSA9IG51bGw7XG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xuXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxud2luZG93Ll93YWlscy5zZXRSZXNpemFibGUgPSBmdW5jdGlvbih2YWx1ZSkge1xuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xufTtcblxud2luZG93Ll93YWlscy5lbmREcmFnID0gZnVuY3Rpb24oKSB7XG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufTtcblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIG9uTW91c2VEb3duKTtcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBvbk1vdXNlTW92ZSk7XG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIG9uTW91c2VVcCk7XG5cblxuZnVuY3Rpb24gZHJhZ1Rlc3QoZSkge1xuICAgIGxldCB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpO1xuICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgaWYgKCF2YWwgfHwgdmFsID09PSBcIlwiIHx8IHZhbC50cmltKCkgIT09IFwiZHJhZ1wiIHx8IG1vdXNlUHJlc3NlZCA9PT0gMCkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xuXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgaW52b2tlKFwicmVzaXplOlwiICsgcmVzaXplRWRnZSk7XG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmIChkcmFnVGVzdChlKSkge1xuICAgICAgICAvLyBUaGlzIGNoZWNrcyBmb3IgY2xpY2tzIG9uIHRoZSBzY3JvbGwgYmFyXG4gICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBzaG91bGREcmFnID0gdHJ1ZTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBvbk1vdXNlVXAoKSB7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCBkZWZhdWx0Q3Vyc29yO1xuICAgIHJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGUpIHtcbiAgICBpZiAoc2hvdWxkRHJhZykge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgICAgICBpbnZva2UoXCJkcmFnXCIpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmICghcmVzaXphYmxlIHx8ICFJc1dpbmRvd3MoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGlmIChkZWZhdWx0Q3Vyc29yID09IG51bGwpIHtcbiAgICAgICAgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3I7XG4gICAgfVxuICAgIGxldCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGxldCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcblxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhc1xuICAgIGxldCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcblxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG5cbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lcnNcbiAgICBsZXQgcmlnaHRDb3JuZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgbGVmdENvcm5lciA9IGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgdG9wQ29ybmVyID0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgYm90dG9tQ29ybmVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcblxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgcmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbiAgICAvLyBBZGp1c3RlZCBmb3IgY29ybmVyIGFyZWFzXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIHRvcENvcm5lcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24sICcnKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcbn1cblxuLyoqXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXG4gKlxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcbndpbmRvdy5fd2FpbHMuY2FsbEVycm9ySGFuZGxlciA9IGVycm9ySGFuZGxlcjtcblxuXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYWxsLCAnJyk7XG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYW5jZWxDYWxsLCAnJyk7XG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxuICpcbiAqIEByZXR1cm4ge3N0cmluZ30gLSBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBjYWxsIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxuICogQHBhcmFtIHtib29sZWFufSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cbiAqXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9IC0gVGhpcyBtZXRob2QgZG9lcyBub3QgcmV0dXJuIGFueSB2YWx1ZS5cbiAqL1xuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XG4gICAgaWYgKHByb21pc2VIYW5kbGVyKSB7XG4gICAgICAgIHByb21pc2VIYW5kbGVyLnJlc29sdmUoaXNKU09OID8gSlNPTi5wYXJzZShkYXRhKSA6IGRhdGEpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBlcnJvciBmcm9tIGEgY2FsbCByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZSB0byByZWplY3QgdGhlIHByb21pc2UgaGFuZGxlciB3aXRoLlxuICpcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZCwgbWVzc2FnZSkge1xuICAgIGNvbnN0IHByb21pc2VIYW5kbGVyID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xuICAgIGlmIChwcm9taXNlSGFuZGxlcikge1xuICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobWVzc2FnZSk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHJpZXZlcyBhbmQgcmVtb3ZlcyB0aGUgcmVzcG9uc2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCBmcm9tIHRoZSBjYWxsUmVzcG9uc2VzIG1hcC5cbiAqXG4gKiBAcGFyYW0ge2FueX0gaWQgLSBUaGUgSUQgb2YgdGhlIHJlc3BvbnNlIHRvIGJlIHJldHJpZXZlZCBhbmQgcmVtb3ZlZC5cbiAqXG4gKiBAcmV0dXJucyB7YW55fSBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQuXG4gKi9cbmZ1bmN0aW9uIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKSB7XG4gICAgY29uc3QgcmVzcG9uc2UgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIHJldHVybiByZXNwb25zZTtcbn1cblxuLyoqXG4gKiBFeGVjdXRlcyBhIGNhbGwgdXNpbmcgdGhlIHByb3ZpZGVkIHR5cGUgYW5kIG9wdGlvbnMuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd8bnVtYmVyfSB0eXBlIC0gVGhlIHR5cGUgb2YgY2FsbCB0byBleGVjdXRlLlxuICogQHBhcmFtIHtPYmplY3R9IFtvcHRpb25zPXt9XSAtIEFkZGl0aW9uYWwgb3B0aW9ucyBmb3IgdGhlIGNhbGwuXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgYmFzZWQgb24gdGhlIHJlc3VsdCBvZiB0aGUgY2FsbC4gSXQgYWxzbyBoYXMgYSBjYW5jZWwgbWV0aG9kIHRvIGNhbmNlbCBhIGxvbmcgcnVubmluZyByZXF1ZXN0LlxuICovXG5mdW5jdGlvbiBjYWxsQmluZGluZyh0eXBlLCBvcHRpb25zID0ge30pIHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICBjb25zdCBkb0NhbmNlbCA9ICgpID0+IHsgY2FuY2VsQ2FsbCh0eXBlLCB7XCJjYWxsLWlkXCI6IGlkfSkgfTtcbiAgICB2YXIgcXVldWVkQ2FuY2VsID0gZmFsc2UsIGNhbGxSdW5uaW5nID0gZmFsc2U7XG4gICAgdmFyIHAgPSBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmUsIHJlamVjdCB9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5cbiAgICAgICAgICAgIHRoZW4oKF8pID0+IHtcbiAgICAgICAgICAgICAgICBjYWxsUnVubmluZyA9IHRydWU7XG4gICAgICAgICAgICAgICAgaWYgKHF1ZXVlZENhbmNlbCkge1xuICAgICAgICAgICAgICAgICAgICBkb0NhbmNlbCgpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0pLlxuICAgICAgICAgICAgY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgICAgICB9KTtcbiAgICB9KTtcbiAgICBwLmNhbmNlbCA9ICgpID0+IHtcbiAgICAgICAgaWYgKGNhbGxSdW5uaW5nKSB7XG4gICAgICAgICAgICBkb0NhbmNlbCgpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcXVldWVkQ2FuY2VsID0gdHJ1ZTtcbiAgICAgICAgfVxuICAgIH07XG4gICAgXG4gICAgcmV0dXJuIHA7XG59XG5cbi8qKlxuICogQ2FsbCBtZXRob2QuXG4gKlxuICogQHBhcmFtIHtPYmplY3R9IG9wdGlvbnMgLSBUaGUgb3B0aW9ucyBmb3IgdGhlIG1ldGhvZC5cbiAqIEByZXR1cm5zIHtPYmplY3R9IC0gVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogRXhlY3V0ZXMgYSBtZXRob2QgYnkgbmFtZS5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJy5cbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHRocm93cyB7RXJyb3J9IElmIHRoZSBuYW1lIGlzIG5vdCBhIHN0cmluZyBvciBpcyBub3QgaW4gdGhlIGNvcnJlY3QgZm9ybWF0LlxuICogQHJldHVybnMgeyp9IFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBleGVjdXRpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCeU5hbWUobmFtZSwgLi4uYXJncykge1xuICAgIGlmICh0eXBlb2YgbmFtZSAhPT0gXCJzdHJpbmdcIiB8fCBuYW1lLnNwbGl0KFwiLlwiKS5sZW5ndGggPCAzKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIkNhbGxCeU5hbWUgcmVxdWlyZXMgYSBzdHJpbmcgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJ1wiKTtcbiAgICB9XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XG4gICAgICAgIG1ldGhvZE5hbWU6IG5hbWUsXG4gICAgICAgIGFyZ3NcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBieSBpdHMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqXG4gKiBAcGFyYW0ge251bWJlcn0gbWV0aG9kSUQgLSBUaGUgSUQgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlJRChtZXRob2RJRCwgLi4uYXJncykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xuICAgICAgICBtZXRob2RJRCxcbiAgICAgICAgYXJnc1xuICAgIH0pO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgcGx1Z2luLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHBsdWdpbi5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJucyB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFBsdWdpbihwbHVnaW5OYW1lLCBtZXRob2ROYW1lLCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XG4gICAgICAgIHBhY2thZ2VOYW1lOiBcIndhaWxzLXBsdWdpbnNcIixcbiAgICAgICAgc3RydWN0TmFtZTogcGx1Z2luTmFtZSxcbiAgICAgICAgbWV0aG9kTmFtZSxcbiAgICAgICAgYXJnc1xuICAgIH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2xpcGJvYXJkLCAnJyk7XG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcbmNvbnN0IENsaXBib2FyZFRleHQgPSAxO1xuXG4vKipcbiAqIFNldHMgdGhlIHRleHQgdG8gdGhlIENsaXBib2FyZC5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxuICogQHJldHVybiB7UHJvbWlzZX0gLSBBIFByb21pc2UgdGhhdCByZXNvbHZlcyB3aGVuIHRoZSBvcGVyYXRpb24gaXMgc3VjY2Vzc2Z1bC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBUZXh0KCkge1xuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gUG9zaXRpb25cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZS5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBZIC0gVGhlIFkgY29vcmRpbmF0ZS5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNpemVcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIHdpZHRoLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgaGVpZ2h0LlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBSZWN0XG4gKiBAcHJvcGVydHkge251bWJlcn0gWCAtIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHRvcC1sZWZ0IGNvcm5lci5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBZIC0gVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgdG9wLWxlZnQgY29ybmVyLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuXG4gKiBAcHJvcGVydHkge251bWJlcn0gSGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTY3JlZW5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJZCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IE5hbWUgLSBIdW1hbiByZWFkYWJsZSBuYW1lIG9mIHRoZSBzY3JlZW4uXG4gKiBAcHJvcGVydHkge251bWJlcn0gU2NhbGUgLSBUaGUgcmVzb2x1dGlvbiBzY2FsZSBvZiB0aGUgc2NyZWVuLiAxID0gc3RhbmRhcmQgcmVzb2x1dGlvbiwgMiA9IGhpZ2ggKFJldGluYSksIGV0Yy5cbiAqIEBwcm9wZXJ0eSB7UG9zaXRpb259IFBvc2l0aW9uIC0gQ29udGFpbnMgdGhlIFggYW5kIFkgY29vcmRpbmF0ZXMgb2YgdGhlIHNjcmVlbidzIHBvc2l0aW9uLlxuICogQHByb3BlcnR5IHtTaXplfSBTaXplIC0gQ29udGFpbnMgdGhlIHdpZHRoIGFuZCBoZWlnaHQgb2YgdGhlIHNjcmVlbi5cbiAqIEBwcm9wZXJ0eSB7UmVjdH0gQm91bmRzIC0gQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LlxuICogQHByb3BlcnR5IHtSZWN0fSBXb3JrQXJlYSAtIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLlxuICogQHByb3BlcnR5IHtib29sZWFufSBJc1ByaW1hcnkgLSBUcnVlIGlmIHRoaXMgaXMgdGhlIHByaW1hcnkgbW9uaXRvciBzZWxlY3RlZCBieSB0aGUgdXNlciBpbiB0aGUgb3BlcmF0aW5nIHN5c3RlbS5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBSb3RhdGlvbiAtIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLlxuICovXG5cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU2NyZWVucywgJycpO1xuXG5jb25zdCBnZXRBbGwgPSAwO1xuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XG5jb25zdCBnZXRDdXJyZW50ID0gMjtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChnZXRBbGwpO1xufVxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKGdldEN1cnJlbnQpO1xufVxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBOzs7QUNBQSxJQUFJLGNBQ0Y7QUFXSyxJQUFJLFNBQVMsQ0FBQyxPQUFPLE9BQU87QUFDakMsTUFBSSxLQUFLO0FBQ1QsTUFBSSxJQUFJO0FBQ1IsU0FBTyxLQUFLO0FBQ1YsVUFBTSxZQUFhLEtBQUssT0FBTyxJQUFJLEtBQU0sQ0FBQztBQUFBLEVBQzVDO0FBQ0EsU0FBTztBQUNUOzs7QUNOQSxJQUFNLGFBQWEsT0FBTyxTQUFTLFNBQVM7QUFHckMsSUFBTSxjQUFjO0FBQUEsRUFDdkIsTUFBTTtBQUFBLEVBQ04sV0FBVztBQUFBLEVBQ1gsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsWUFBWTtBQUNoQjtBQUNPLElBQUksV0FBVyxPQUFPO0FBc0J0QixTQUFTLHVCQUF1QixRQUFRLFlBQVk7QUFDdkQsU0FBTyxTQUFVLFFBQVEsT0FBSyxNQUFNO0FBQ2hDLFdBQU8sa0JBQWtCLFFBQVEsUUFBUSxZQUFZLElBQUk7QUFB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QSxJQUNKLDRCQUE0QjtBQUFBLElBQzVCLHVDQUF1QztBQUFBLElBQ3ZDLHlDQUF5QztBQUFBLElBQ3pDLDBCQUEwQjtBQUFBLElBQzFCLG9DQUFvQztBQUFBLElBQ3BDLHNDQUFzQztBQUFBLElBQ3RDLG9DQUFvQztBQUFBLElBQ3BDLDBDQUEwQztBQUFBLElBQzFDLCtCQUErQjtBQUFBLElBQy9CLG9CQUFvQjtBQUFBLElBQ3BCLHdDQUF3QztBQUFBLElBQ3hDLHNCQUFzQjtBQUFBLElBQ3RCLHNCQUFzQjtBQUFBLElBQ3RCLDZCQUE2QjtBQUFBLElBQzdCLGdDQUFnQztBQUFBLElBQ2hDLHFCQUFxQjtBQUFBLElBQ3JCLDZCQUE2QjtBQUFBLElBQzdCLDBCQUEwQjtBQUFBLElBQzFCLHVCQUF1QjtBQUFBLElBQ3ZCLHVCQUF1QjtBQUFBLElBQ3ZCLDJCQUEyQjtBQUFBLElBQzNCLCtCQUErQjtBQUFBLElBQy9CLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLHFCQUFxQjtBQUFBLElBQ3JCLHNCQUFzQjtBQUFBLElBQ3RCLGdDQUFnQztBQUFBLElBQ2hDLGtDQUFrQztBQUFBLElBQ2xDLG1DQUFtQztBQUFBLElBQ25DLG9DQUFvQztBQUFBLElBQ3BDLCtCQUErQjtBQUFBLElBQy9CLDZCQUE2QjtBQUFBLElBQzdCLHVCQUF1QjtBQUFBLElBQ3ZCLGlDQUFpQztBQUFBLElBQ2pDLDhCQUE4QjtBQUFBLElBQzlCLDRCQUE0QjtBQUFBLElBQzVCLHNDQUFzQztBQUFBLElBQ3RDLDRCQUE0QjtBQUFBLElBQzVCLHNCQUFzQjtBQUFBLElBQ3RCLGtDQUFrQztBQUFBLElBQ2xDLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLDJCQUEyQjtBQUFBLElBQzNCLHdCQUF3QjtBQUFBLElBQ3hCLG1CQUFtQjtBQUFBLElBQ25CLDBCQUEwQjtBQUFBLElBQzFCLDhCQUE4QjtBQUFBLElBQzlCLHlCQUF5QjtBQUFBLElBQ3pCLDZCQUE2QjtBQUFBLElBQzdCLGlCQUFpQjtBQUFB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U0FBUyxZQUFZLEtBQUssTUFBTSxHQUFHLEVBQUUsU0FBUyxHQUFHO0FBQ3hELFVBQU0sSUFBSSxNQUFNLG9FQUFvRTtBQUFBLEVBQ3hGO0FBQ0EsU0FBTyxZQUFZLGFBQWE7QUFBQSxJQUM1QixZQUFZO0FBQUEsSUFDWjtBQUFBLEVBQ0osQ0FBQztBQUNMO0FBU08sU0FBUyxLQUFLLGFBQWEsTUFBTTtBQUNwQyxTQUFPLFlBQVksYUFBYTtBQUFBLElBQzVCO0FBQUEsSUFDQTtBQUFBLEVBQ0osQ0FBQztBQUNMO0FBVU8sU0FBUyxPQUFPLFlBQVksZUFBZSxNQUFNO0FBQ3BELFNBQU8sWUFBWSxhQUFhO0FBQUEsSUFDNUIsYUFBYTtBQUFBLElBQ2IsWUFBWTtBQUFBLElBQ1o7QUFBQSxJQUNBO0FBQUEsRUFDSixDQUFDO0FBQ0w7OztBQ2hMQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBY0EsSUFBTUUsUUFBTyx1QkFBdUIsWUFBWSxXQUFXLEVBQUU7QUFDN0QsSUFBTSxtQkFBbUI7QUFDekIsSUFBTSxnQkFBZ0I7QUFRZixTQUFTLFFBQVEsTUFBTTtBQUMxQixTQUFPQSxNQUFLLGtCQUFrQixFQUFDLEtBQUksQ0FBQztBQUN4QztBQU1PLFNBQVMsT0FBTztBQUNuQixTQUFPQSxNQUFLLGFBQWE7QUFDN0I7OztBQ2xDQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFpREEsSUFBTUMsUUFBTyx1QkFBdUIsWUFBWSxTQUFTLEVBQUU7QUFFM0QsSUFBTSxTQUFTO0FBQ2YsSUFBTSxhQUFhO0FBQ25CLElBQU0sYUFBYTtBQU1aLFNBQVMsU0FBUztBQUNyQixTQUFPQSxNQUFLLE1BQU07QUFDdEI7QUFLTyxTQUFTLGFBQWE7QUFDekIsU0FBT0EsTUFBSyxVQUFVO0FBQzFCO0FBTU8sU0FBUyxhQUFhO0FBQ3pCLFNBQU9BLE1BQUssVUFBVTtBQUMxQjs7O0FsQmpFQSxPQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFpQ2xDLE9BQU8sT0FBTyxTQUFnQjtBQUN2QixPQUFPLHFCQUFxQjsiLAogICJuYW1lcyI6IFsiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjb250cm9sbGVyIiwgInJlc2l6YWJsZSIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIaWRlTWV0aG9kIiwgIlNob3dNZXRob2QiLCAiY2FsbCIsICJnZW5lcmF0ZUlEIiwgImNhbGwiLCAiY2FsbCJdCn0K
