var __defProp = Object.defineProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};

// desktop/@wailsio/runtime/src/index.js
var index_exports = {};
__export(index_exports, {
  Application: () => application_exports,
  Browser: () => browser_exports,
  Call: () => calls_exports,
  Clipboard: () => clipboard_exports,
  Create: () => create_exports,
  Dialogs: () => dialogs_exports,
  Events: () => events_exports,
  Flags: () => flags_exports,
  Screens: () => screens_exports,
  System: () => system_exports,
  WML: () => wml_exports,
  Window: () => window_default,
  init: () => init
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

// desktop/@wailsio/runtime/src/nanoid.js
var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
var nanoid = (size = 21) => {
  let id = "";
  let i = size | 0;
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
    WindowClosing: "windows:WindowClosing",
    WindowSetFocus: "windows:WindowSetFocus",
    WindowKillFocus: "windows:WindowKillFocus",
    WindowDragDrop: "windows:WindowDragDrop",
    WindowDragEnter: "windows:WindowDragEnter",
    WindowDragLeave: "windows:WindowDragLeave",
    WindowDragOver: "windows:WindowDragOver",
    WindowDidMove: "windows:WindowDidMove",
    WindowDidResize: "windows:WindowDidResize",
    WindowShow: "windows:WindowShow",
    WindowHide: "windows:WindowHide",
    WindowStartMove: "windows:WindowStartMove",
    WindowEndMove: "windows:WindowEndMove",
    WindowStartResize: "windows:WindowStartResize",
    WindowEndResize: "windows:WindowEndResize",
    WindowKeyDown: "windows:WindowKeyDown",
    WindowKeyUp: "windows:WindowKeyUp",
    WindowZOrderChanged: "windows:WindowZOrderChanged",
    WindowPaint: "windows:WindowPaint",
    WindowBackgroundErase: "windows:WindowBackgroundErase",
    WindowNonClientHit: "windows:WindowNonClientHit",
    WindowNonClientMouseDown: "windows:WindowNonClientMouseDown",
    WindowNonClientMouseUp: "windows:WindowNonClientMouseUp",
    WindowNonClientMouseMove: "windows:WindowNonClientMouseMove",
    WindowNonClientMouseLeave: "windows:WindowNonClientMouseLeave",
    WindowDPIChanged: "windows:WindowDPIChanged"
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
    WindowDidDeminiaturize: "mac:WindowDidDeminiaturize",
    WindowDidEndSheet: "mac:WindowDidEndSheet",
    WindowDidEnterFullScreen: "mac:WindowDidEnterFullScreen",
    WindowMaximise: "mac:WindowMaximise",
    WindowUnMaximise: "mac:WindowUnMaximise",
    WindowDidZoom: "mac:WindowDidZoom!",
    WindowZoomIn: "mac:WindowZoomIn!",
    WindowZoomOut: "mac:WindowZoomOut!",
    WindowZoomReset: "mac:WindowZoomReset!",
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
    WindowFileDraggingExited: "mac:WindowFileDraggingExited",
    WindowShow: "mac:WindowShow",
    WindowHide: "mac:WindowHide"
  },
  Linux: {
    SystemThemeChanged: "linux:SystemThemeChanged",
    WindowLoadChanged: "linux:WindowLoadChanged",
    WindowDeleteEvent: "linux:WindowDeleteEvent",
    WindowDidMove: "linux:WindowDidMove",
    WindowDidResize: "linux:WindowDidResize",
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
    ThemeChanged: "common:ThemeChanged",
    WindowDidMove: "common:WindowDidMove",
    WindowDidResize: "common:WindowDidResize",
    ApplicationOpenedWithFile: "common:ApplicationOpenedWithFile"
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
      if (this.maxCallbacks === -1) return false;
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
      if (remove) return true;
    });
    if (toRemove.length > 0) {
      listeners = listeners.filter((l) => !toRemove.includes(l));
      if (listeners.length === 0) eventListeners.delete(event.name);
      else eventListeners.set(event.name, listeners);
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
  if (listeners.length === 0) eventListeners.delete(eventName);
  else eventListeners.set(eventName, listeners);
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
var PositionMethod = 0;
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
var SetPositionMethod = 24;
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
   * Returns the absolute position of the window.
   *
   * @public
   * @return {Promise<Position>} - The current absolute position of the window.
   */
  Position() {
    return this[caller](PositionMethod);
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
   * Sets the absolute position of the window.
   *
   * @public
   * @param {number} x - The desired horizontal absolute position of the window
   * @param {number} y - The desired vertical absolute position of the window
   * @return {Promise<void>}
   */
  SetPosition(x, y) {
    return this[caller](SetPositionMethod, { x, y });
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
    const eventType = element.getAttribute("data-wml-event");
    const targetWindow = element.getAttribute("data-wml-target-window") || "";
    const windowMethod = element.getAttribute("data-wml-window");
    const url = element.getAttribute("data-wml-openURL");
    if (eventType !== null)
      sendEvent(eventType);
    if (windowMethod !== null)
      callWindowMethod(targetWindow, windowMethod);
    if (url !== null)
      void OpenURL(url);
  }
  const confirm = element.getAttribute("data-wml-confirm");
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
  const triggerAttr = element.getAttribute("data-wml-trigger") || "click";
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
  document.body.querySelectorAll("[data-wml-event], [data-wml-window], [data-wml-openURL]").forEach(addWMLListeners);
}

// desktop/compiled/main.js
window.wails = index_exports;
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
var _invoke = (() => {
  try {
    if (window?.chrome?.webview) {
      return (msg) => window.chrome.webview.postMessage(msg);
    }
    if (window?.webkit?.messageHandlers?.external) {
      return (msg) => window.webkit.messageHandlers.external.postMessage(msg);
    }
  } catch (e) {
    console.warn(
      "\n%c\u26A0\uFE0F Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n",
      "background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;",
      "background: transparent;",
      "color: #ffffff; font-style: italic; font-weight: bold;"
    );
  }
  return null;
})();
function invoke(msg) {
  if (!_invoke) return;
  return _invoke(msg);
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
    invoke("wails:resize:" + resizeEdge);
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
      invoke("wails:drag");
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
  } else if (rightCorner && bottomCorner) setResize("se-resize");
  else if (leftCorner && bottomCorner) setResize("sw-resize");
  else if (leftCorner && topCorner) setResize("nw-resize");
  else if (topCorner && rightCorner) setResize("ne-resize");
  else if (leftBorder) setResize("w-resize");
  else if (topBorder) setResize("n-resize");
  else if (bottomBorder) setResize("s-resize");
  else if (rightBorder) setResize("e-resize");
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
    return cancelCall(type, { "call-id": id });
  };
  let queuedCancel = false, callRunning = false;
  let p = new Promise((resolve, reject) => {
    options["call-id"] = id;
    callResponses.set(id, { resolve, reject });
    call7(type, options).then((_) => {
      callRunning = true;
      if (queuedCancel) {
        return doCancel();
      }
    }).catch((error) => {
      reject(error);
      callResponses.delete(id);
    });
  });
  p.cancel = () => {
    if (callRunning) {
      return doCancel();
    } else {
      queuedCancel = true;
    }
  };
  return p;
}
function Call(options) {
  return callBinding(CallBinding, options);
}
function ByName(methodName, ...args) {
  return callBinding(CallBinding, {
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

// desktop/@wailsio/runtime/src/create.js
var create_exports = {};
__export(create_exports, {
  Any: () => Any,
  Array: () => Array,
  ByteSlice: () => ByteSlice,
  Map: () => Map2,
  Nullable: () => Nullable,
  Struct: () => Struct
});
function Any(source) {
  return (
    /** @type {T} */
    source
  );
}
function ByteSlice(source) {
  return (
    /** @type {any} */
    source == null ? "" : source
  );
}
function Array(element) {
  if (element === Any) {
    return (source) => source === null ? [] : source;
  }
  return (source) => {
    if (source === null) {
      return [];
    }
    for (let i = 0; i < source.length; i++) {
      source[i] = element(source[i]);
    }
    return source;
  };
}
function Map2(key, value) {
  if (value === Any) {
    return (source) => source === null ? {} : source;
  }
  return (source) => {
    if (source === null) {
      return {};
    }
    for (const key2 in source) {
      source[key2] = value(source[key2]);
    }
    return source;
  };
}
function Nullable(element) {
  if (element === Any) {
    return Any;
  }
  return (source) => source === null ? null : element(source);
}
function Struct(createField) {
  let allAny = true;
  for (const name in createField) {
    if (createField[name] !== Any) {
      allAny = false;
      break;
    }
  }
  if (allAny) {
    return Any;
  }
  return (source) => {
    for (const name in createField) {
      if (name in source) {
        source[name] = createField[name](source[name]);
      }
    }
    return source;
  };
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
var initialised = false;
function init() {
  window._wails.invoke = invoke;
  invoke("wails:runtime:ready");
  initialised = true;
}
window.addEventListener("load", () => {
  if (!initialised) {
    init();
  }
});
export {
  application_exports as Application,
  browser_exports as Browser,
  calls_exports as Call,
  clipboard_exports as Clipboard,
  create_exports as Create,
  dialogs_exports as Dialogs,
  events_exports as Events,
  flags_exports as Flags,
  screens_exports as Screens,
  system_exports as System,
  wml_exports as WML,
  window_default as Window,
  init
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudF90eXBlcy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMuanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG5cclxuaW1wb3J0IFwiLi9jb250ZXh0bWVudVwiO1xyXG5pbXBvcnQgXCIuL2RyYWdcIjtcclxuXHJcbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXHJcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gXCIuL2FwcGxpY2F0aW9uXCI7XHJcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3NlclwiO1xyXG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzXCI7XHJcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmRcIjtcclxuaW1wb3J0ICogYXMgQ3JlYXRlIGZyb20gXCIuL2NyZWF0ZVwiO1xyXG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3NcIjtcclxuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuL2V2ZW50c1wiO1xyXG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi9mbGFnc1wiO1xyXG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnNcIjtcclxuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gXCIuL3N5c3RlbVwiO1xyXG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvd1wiO1xyXG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sXCI7XHJcblxyXG5leHBvcnQge1xyXG4gICAgQXBwbGljYXRpb24sXHJcbiAgICBCcm93c2VyLFxyXG4gICAgQ2FsbCxcclxuICAgIENsaXBib2FyZCxcclxuICAgIENyZWF0ZSxcclxuICAgIERpYWxvZ3MsXHJcbiAgICBFdmVudHMsXHJcbiAgICBGbGFncyxcclxuICAgIFNjcmVlbnMsXHJcbiAgICBTeXN0ZW0sXHJcbiAgICBXaW5kb3csXHJcbiAgICBXTUxcclxufTtcclxuXHJcbmxldCBpbml0aWFsaXNlZCA9IGZhbHNlO1xyXG5leHBvcnQgZnVuY3Rpb24gaW5pdCgpIHtcclxuICAgIHdpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcclxuICAgIFN5c3RlbS5pbnZva2UoXCJ3YWlsczpydW50aW1lOnJlYWR5XCIpO1xyXG4gICAgaW5pdGlhbGlzZWQgPSB0cnVlO1xyXG59XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihcImxvYWRcIiwgKCkgPT4ge1xyXG4gICAgaWYgKCFpbml0aWFsaXNlZCkge1xyXG4gICAgICAgIGluaXQoKTtcclxuICAgIH1cclxufSk7XHJcblxyXG4vLyBOb3RpZnkgYmFja2VuZFxyXG5cclxuIiwgIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7T3BlblVSTH0gZnJvbSBcIi4vYnJvd3NlclwiO1xyXG5pbXBvcnQge1F1ZXN0aW9ufSBmcm9tIFwiLi9kaWFsb2dzXCI7XHJcbmltcG9ydCB7RW1pdCwgV2FpbHNFdmVudH0gZnJvbSBcIi4vZXZlbnRzXCI7XHJcbmltcG9ydCB7Y2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeX0gZnJvbSBcIi4vdXRpbHNcIjtcclxuaW1wb3J0IFdpbmRvdyBmcm9tIFwiLi93aW5kb3dcIjtcclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBvcHRpb25hbCBkYXRhLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXHJcbiAqIEBwYXJhbSB7YW55fSBbZGF0YT1udWxsXSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xyXG4gICAgRW1pdChuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgc3BlY2lmaWVkIHdpbmRvdy5cclxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGNhbGwgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKHdpbmRvd05hbWUsIG1ldGhvZE5hbWUpIHtcclxuICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IFdpbmRvdy5HZXQod2luZG93TmFtZSk7XHJcbiAgICBjb25zdCBtZXRob2QgPSB0YXJnZXRXaW5kb3dbbWV0aG9kTmFtZV07XHJcblxyXG4gICAgaWYgKHR5cGVvZiBtZXRob2QgIT09IFwiZnVuY3Rpb25cIikge1xyXG4gICAgICAgIGNvbnNvbGUuZXJyb3IoYFdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nIG5vdCBmb3VuZGApO1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICB0cnkge1xyXG4gICAgICAgIG1ldGhvZC5jYWxsKHRhcmdldFdpbmRvdyk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc29sZS5lcnJvcihgRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JzogYCwgZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXNwb25kcyB0byBhIHRyaWdnZXJpbmcgZXZlbnQgYnkgcnVubmluZyBhcHByb3ByaWF0ZSBXTUwgYWN0aW9ucyBmb3IgdGhlIGN1cnJlbnQgdGFyZ2V0XHJcbiAqXHJcbiAqIEBwYXJhbSB7RXZlbnR9IGV2XHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBvbldNTFRyaWdnZXJlZChldikge1xyXG4gICAgY29uc3QgZWxlbWVudCA9IGV2LmN1cnJlbnRUYXJnZXQ7XHJcblxyXG4gICAgZnVuY3Rpb24gcnVuRWZmZWN0KGNob2ljZSA9IFwiWWVzXCIpIHtcclxuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xyXG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10YXJnZXQtd2luZG93JykgfHwgXCJcIjtcclxuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XHJcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW5VUkwnKTtcclxuXHJcbiAgICAgICAgaWYgKGV2ZW50VHlwZSAhPT0gbnVsbClcclxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XHJcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcclxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcclxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XHJcblxyXG4gICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICBRdWVzdGlvbih7XHJcbiAgICAgICAgICAgIFRpdGxlOiBcIkNvbmZpcm1cIixcclxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcclxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxyXG4gICAgICAgICAgICBCdXR0b25zOiBbXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XHJcbiAgICAgICAgICAgIF1cclxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJ1bkVmZmVjdCgpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogQHR5cGUge3N5bWJvbH1cclxuICovXHJcbmNvbnN0IGNvbnRyb2xsZXIgPSBTeW1ib2woKTtcclxuXHJcbi8qKlxyXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXHJcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXHJcbiAqL1xyXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XHJcbiAgICBjb25zdHJ1Y3RvcigpIHtcclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTdG9yZXMgdGhlIEFib3J0Q29udHJvbGxlciB0aGF0IGNhbiBiZSB1c2VkIHRvIHJlbW92ZSBhbGwgY3VycmVudGx5IGFjdGl2ZSBsaXN0ZW5lcnMuXHJcbiAgICAgICAgICpcclxuICAgICAgICAgKiBAcHJpdmF0ZVxyXG4gICAgICAgICAqIEBuYW1lIHtAbGluayBjb250cm9sbGVyfVxyXG4gICAgICAgICAqIEBtZW1iZXIge0Fib3J0Q29udHJvbGxlcn1cclxuICAgICAgICAgKi9cclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJdID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyBhbiBvcHRpb25zIG9iamVjdCBmb3IgYWRkRXZlbnRMaXN0ZW5lciB0aGF0IHRpZXMgdGhlIGxpc3RlbmVyXHJcbiAgICAgKiB0byB0aGUgQWJvcnRTaWduYWwgZnJvbSB0aGUgY3VycmVudCBBYm9ydENvbnRyb2xsZXIuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHtIVE1MRWxlbWVudH0gZWxlbWVudCBBbiBIVE1MIGVsZW1lbnRcclxuICAgICAqIEBwYXJhbSB7c3RyaW5nW119IHRyaWdnZXJzIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudHNcclxuICAgICAqIEByZXR1cm5zIHtBZGRFdmVudExpc3RlbmVyT3B0aW9uc31cclxuICAgICAqL1xyXG4gICAgc2V0KGVsZW1lbnQsIHRyaWdnZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuIHsgc2lnbmFsOiB0aGlzW2NvbnRyb2xsZXJdLnNpZ25hbCB9O1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMge3ZvaWR9XHJcbiAgICAgKi9cclxuICAgIHJlc2V0KCkge1xyXG4gICAgICAgIHRoaXNbY29udHJvbGxlcl0uYWJvcnQoKTtcclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJdID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogQHR5cGUge3N5bWJvbH1cclxuICovXHJcbmNvbnN0IHRyaWdnZXJNYXAgPSBTeW1ib2woKTtcclxuXHJcbi8qKlxyXG4gKiBAdHlwZSB7c3ltYm9sfVxyXG4gKi9cclxuY29uc3QgZWxlbWVudENvdW50ID0gU3ltYm9sKCk7XHJcblxyXG4vKipcclxuICogV2Vha01hcFJlZ2lzdHJ5IG1hcHMgYWN0aXZlIHRyaWdnZXIgZXZlbnRzIHRvIGVhY2ggRE9NIGVsZW1lbnQgdGhyb3VnaCBhIFdlYWtNYXAuXHJcbiAqIFRoaXMgZW5zdXJlcyB0aGF0IHRoZSBtYXBwaW5nIHJlbWFpbnMgcHJpdmF0ZSB0byB0aGlzIG1vZHVsZSwgd2hpbGUgc3RpbGwgYWxsb3dpbmcgZ2FyYmFnZVxyXG4gKiBjb2xsZWN0aW9uIG9mIHRoZSBpbnZvbHZlZCBlbGVtZW50cy5cclxuICovXHJcbmNsYXNzIFdlYWtNYXBSZWdpc3RyeSB7XHJcbiAgICBjb25zdHJ1Y3RvcigpIHtcclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTdG9yZXMgdGhlIGN1cnJlbnQgZWxlbWVudC10by10cmlnZ2VyIG1hcHBpbmcuXHJcbiAgICAgICAgICpcclxuICAgICAgICAgKiBAcHJpdmF0ZVxyXG4gICAgICAgICAqIEBuYW1lIHtAbGluayB0cmlnZ2VyTWFwfVxyXG4gICAgICAgICAqIEBtZW1iZXIge1dlYWtNYXA8SFRNTEVsZW1lbnQsIHN0cmluZ1tdPn1cclxuICAgICAgICAgKi9cclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBdID0gbmV3IFdlYWtNYXAoKTtcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLlxyXG4gICAgICAgICAqXHJcbiAgICAgICAgICogQHByaXZhdGVcclxuICAgICAgICAgKiBAbmFtZSB7QGxpbmsgZWxlbWVudENvdW50fVxyXG4gICAgICAgICAqIEBtZW1iZXIge251bWJlcn1cclxuICAgICAgICAgKi9cclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudF0gPSAwO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgYWN0aXZlIHRyaWdnZXJzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHtIVE1MRWxlbWVudH0gZWxlbWVudCBBbiBIVE1MIGVsZW1lbnRcclxuICAgICAqIEBwYXJhbSB7c3RyaW5nW119IHRyaWdnZXJzIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudFxyXG4gICAgICogQHJldHVybnMge0FkZEV2ZW50TGlzdGVuZXJPcHRpb25zfVxyXG4gICAgICovXHJcbiAgICBzZXQoZWxlbWVudCwgdHJpZ2dlcnMpIHtcclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudF0gKz0gIXRoaXNbdHJpZ2dlck1hcF0uaGFzKGVsZW1lbnQpO1xyXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcF0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcclxuICAgICAgICByZXR1cm4ge307XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyB7dm9pZH1cclxuICAgICAqL1xyXG4gICAgcmVzZXQoKSB7XHJcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50XSA8PSAwKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xyXG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRdIDw9IDApXHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuXHJcbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwXS5nZXQoZWxlbWVudCk7XHJcbiAgICAgICAgICAgIHRoaXNbZWxlbWVudENvdW50XSAtPSAodHlwZW9mIHRyaWdnZXJzICE9PSBcInVuZGVmaW5lZFwiKTtcclxuXHJcbiAgICAgICAgICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycyB8fCBbXSlcclxuICAgICAgICAgICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBdID0gbmV3IFdlYWtNYXAoKTtcclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudF0gPSAwO1xyXG4gICAgfVxyXG59XHJcblxyXG5jb25zdCB0cmlnZ2VyUmVnaXN0cnkgPSBjYW5BYm9ydExpc3RlbmVycygpID8gbmV3IEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5KCkgOiBuZXcgV2Vha01hcFJlZ2lzdHJ5KCk7XHJcblxyXG4vKipcclxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gdGhlIHNwZWNpZmllZCBlbGVtZW50LlxyXG4gKlxyXG4gKiBAcGFyYW0ge0hUTUxFbGVtZW50fSBlbGVtZW50XHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxMaXN0ZW5lcnMoZWxlbWVudCkge1xyXG4gICAgY29uc3QgdHJpZ2dlclJlZ0V4cCA9IC9cXFMrL2c7XHJcbiAgICBjb25zdCB0cmlnZ2VyQXR0ciA9IChlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIik7XHJcbiAgICBjb25zdCB0cmlnZ2VycyA9IFtdO1xyXG5cclxuICAgIGxldCBtYXRjaDtcclxuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcclxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcclxuXHJcbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XHJcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVuYWJsZSgpIHtcclxuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpIHtcclxuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xyXG4gICAgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtZXZlbnRdLCBbZGF0YS13bWwtd2luZG93XSwgW2RhdGEtd21sLW9wZW5VUkxdJykuZm9yRWFjaChhZGRXTUxMaXN0ZW5lcnMpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkJyb3dzZXIsICcnKTtcclxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xyXG5cclxuLyoqXHJcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybCkge1xyXG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmx9KTtcclxufVxyXG4iLCAiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vYWkvbmFub2lkXG5cbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxuLy9cbi8vIENvcHlyaWdodCAyMDE3IEFuZHJleSBTaXRuaWsgPGFuZHJleUBzaXRuaWsucnU+XG4vL1xuLy8gUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weSBvZlxuLy8gdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbCBpblxuLy8gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cyB0b1xuLy8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbCBjb3BpZXMgb2Zcbi8vIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbyxcbi8vICAgICBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcbi8vXG4vLyAgICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXG4vLyBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLlxuLy9cbi8vICAgICBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTU1xuLy8gRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBBVVRIT1JTIE9SXG4vLyBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVJcbi8vIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwgT1VUIE9GIE9SIElOXG4vLyBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFIFNPRlRXQVJFLlxuXG4vLyBUaGlzIGFscGhhYmV0IHVzZXMgYEEtWmEtejAtOV8tYCBzeW1ib2xzLlxuLy8gVGhlIG9yZGVyIG9mIGNoYXJhY3RlcnMgaXMgb3B0aW1pemVkIGZvciBiZXR0ZXIgZ3ppcCBhbmQgYnJvdGxpIGNvbXByZXNzaW9uLlxuLy8gUmVmZXJlbmNlcyB0byB0aGUgc2FtZSBmaWxlICh3b3JrcyBib3RoIGZvciBnemlwIGFuZCBicm90bGkpOlxuLy8gYCd1c2VgLCBgYW5kb21gLCBhbmQgYHJpY3QnYFxuLy8gUmVmZXJlbmNlcyB0byB0aGUgYnJvdGxpIGRlZmF1bHQgZGljdGlvbmFyeTpcbi8vIGAtMjZUYCwgYDE5ODNgLCBgNDBweGAsIGA3NXB4YCwgYGJ1c2hgLCBgamFja2AsIGBtaW5kYCwgYHZlcnlgLCBhbmQgYHdvbGZgXG5sZXQgdXJsQWxwaGFiZXQgPVxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuXG5leHBvcnQgbGV0IG5hbm9pZCA9IChzaXplID0gMjEpID0+IHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIC8vIEEgY29tcGFjdCBhbHRlcm5hdGl2ZSBmb3IgYGZvciAodmFyIGkgPSAwOyBpIDwgc3RlcDsgaSsrKWAuXG4gICAgbGV0IGkgPSBzaXplIHwgMFxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgICAgLy8gYHwgMGAgaXMgbW9yZSBjb21wYWN0IGFuZCBmYXN0ZXIgdGhhbiBgTWF0aC5mbG9vcigpYC5cbiAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbn0iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XHJcblxyXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcclxuXHJcbi8vIE9iamVjdCBOYW1lc1xyXG5leHBvcnQgY29uc3Qgb2JqZWN0TmFtZXMgPSB7XHJcbiAgICBDYWxsOiAwLFxyXG4gICAgQ2xpcGJvYXJkOiAxLFxyXG4gICAgQXBwbGljYXRpb246IDIsXHJcbiAgICBFdmVudHM6IDMsXHJcbiAgICBDb250ZXh0TWVudTogNCxcclxuICAgIERpYWxvZzogNSxcclxuICAgIFdpbmRvdzogNixcclxuICAgIFNjcmVlbnM6IDcsXHJcbiAgICBTeXN0ZW06IDgsXHJcbiAgICBCcm93c2VyOiA5LFxyXG4gICAgQ2FuY2VsQ2FsbDogMTAsXHJcbn1cclxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xyXG5cclxuLyoqXHJcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cclxuICpcclxuICogQHBhcmFtIHtPYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3Qgb24gd2hpY2ggdGhlIG1ldGhvZCBpcyB0byBiZSBpbnZva2VkLlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXHJcbiAqIEByZXR1cm5zIHtGdW5jdGlvbn0gQSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IHRha2VzIHRoZSBtZXRob2QgbmFtZSBhbmQgb3B0aW9uYWxseSBhcmd1bWVudHMgYW5kIGludm9rZXMgdGhlIG1ldGhvZCB3aXRoaW4gdGhlIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xyXG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogQ3JlYXRlcyBhIG5ldyBydW50aW1lIGNhbGxlciB3aXRoIHNwZWNpZmllZCBJRC5cclxuICpcclxuICogQHBhcmFtIHtvYmplY3R9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICogQHJldHVybiB7RnVuY3Rpb259IC0gVGhlIG5ldyBydW50aW1lIGNhbGxlciBmdW5jdGlvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdCwgd2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xyXG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XHJcbiAgICB9O1xyXG59XHJcblxyXG5cclxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcclxuICAgIGlmKCBtZXRob2QgKSB7XHJcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIH1cclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG5cclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgZmV0Y2godXJsLCBmZXRjaE9wdGlvbnMpXHJcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcclxuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xyXG4gICAgICAgICAgICAgICAgICAgIC8vIGNoZWNrIGNvbnRlbnQgdHlwZVxyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcclxuICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XHJcbiAgICAgICAgICAgIH0pXHJcbiAgICAgICAgICAgIC50aGVuKGRhdGEgPT4gcmVzb2x2ZShkYXRhKSlcclxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgaWYgKGFyZ3MpIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xyXG4gICAgfVxyXG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IE9wZW5GaWxlRGlhbG9nT3B0aW9uc1xyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNyZWF0ZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbU2hvd0hpZGRlbkZpbGVzXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtSZXNvbHZlc0FsaWFzZXNdIC0gSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbl0gLSBJbmRpY2F0ZXMgaWYgbXVsdGlwbGUgc2VsZWN0aW9uIGlzIGFsbG93ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0hpZGVFeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1RyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c090aGVyRmlsZXR5cGVzXSAtIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLlxyXG4gKiBAcHJvcGVydHkge0ZpbGVGaWx0ZXJbXX0gW0ZpbHRlcnNdIC0gQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1RpdGxlXSAtIFRpdGxlIG9mIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTWVzc2FnZV0gLSBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtCdXR0b25UZXh0XSAtIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0RpcmVjdG9yeV0gLSBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LlxyXG4gKi9cclxuXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gU2F2ZUZpbGVEaWFsb2dPcHRpb25zXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRmlsZW5hbWVdIC0gRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DcmVhdGVEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW1Nob3dIaWRkZW5GaWxlc10gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb25dIC0gSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtIaWRlRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uXHJcbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtBbGxvd3NPdGhlckZpbGV0eXBlc10gLSBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC5cclxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaXRsZSBvZiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtEaXJlY3RvcnldIC0gRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBbRGV0YWNoZWRdIC0gSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy5cclxuICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gTWVzc2FnZURpYWxvZ09wdGlvbnNcclxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTWVzc2FnZV0gLSBUaGUgbWFpbiBtZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cclxuICogQHByb3BlcnR5IHtCdXR0b25bXX0gW0J1dHRvbnNdIC0gQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuXHJcbiAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IEJ1dHRvblxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW0xhYmVsXSAtIFRleHQgdGhhdCBhcHBlYXJzIHdpdGhpbiB0aGUgYnV0dG9uLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtJc0NhbmNlbF0gLSBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtJc0RlZmF1bHRdIC0gVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBiZSB0aGUgZGVmYXVsdCBhY3Rpb24gd2hlbiB0aGUgdXNlciBwcmVzc2VzIGVudGVyLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBGaWxlRmlsdGVyXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlzcGxheU5hbWVdIC0gRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbUGF0dGVybl0gLSBQYXR0ZXJuIHRvIG1hdGNoIGZvciB0aGUgZmlsdGVyLCBlLmcuIFwiKi50eHQ7Ki5tZFwiIGZvciB0ZXh0IG1hcmtkb3duIGZpbGVzLlxyXG4gKi9cclxuXHJcbi8vIHNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmRpYWxvZ0Vycm9yQ2FsbGJhY2sgPSBkaWFsb2dFcnJvckNhbGxiYWNrO1xyXG53aW5kb3cuX3dhaWxzLmRpYWxvZ1Jlc3VsdENhbGxiYWNrID0gZGlhbG9nUmVzdWx0Q2FsbGJhY2s7XHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XHJcblxyXG4vLyBEZWZpbmUgY29uc3RhbnRzIGZyb20gdGhlIGBtZXRob2RzYCBvYmplY3QgaW4gVGl0bGUgQ2FzZVxyXG5jb25zdCBEaWFsb2dJbmZvID0gMDtcclxuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XHJcbmNvbnN0IERpYWxvZ0Vycm9yID0gMjtcclxuY29uc3QgRGlhbG9nUXVlc3Rpb24gPSAzO1xyXG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XHJcbmNvbnN0IERpYWxvZ1NhdmVGaWxlID0gNTtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZywgJycpO1xyXG5jb25zdCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwKCk7XHJcblxyXG4vKipcclxuICogR2VuZXJhdGVzIGEgdW5pcXVlIGlkIHRoYXQgaXMgbm90IHByZXNlbnQgaW4gZGlhbG9nUmVzcG9uc2VzLlxyXG4gKiBAcmV0dXJucyB7c3RyaW5nfSB1bmlxdWUgaWRcclxuICovXHJcbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XHJcbiAgICBsZXQgcmVzdWx0O1xyXG4gICAgZG8ge1xyXG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xyXG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhIGRpYWxvZyBvZiBzcGVjaWZpZWQgdHlwZSB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge251bWJlcn0gdHlwZSAtIHR5cGUgb2YgZGlhbG9nXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN8T3BlbkZpbGVEaWFsb2dPcHRpb25zfFNhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIG9wdGlvbnMgZm9yIHRoZSBkaWFsb2dcclxuICogQHJldHVybnMge1Byb21pc2V9IHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHJlc3VsdCBvZiBkaWFsb2dcclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zID0ge30pIHtcclxuICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XHJcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5jYXRjaCgoZXJyb3IpID0+IHtcclxuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcclxuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIGNhbGxiYWNrIGZyb20gYSBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBJRCBvZiB0aGUgZGlhbG9nIHJlc3BvbnNlLlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZGF0YSAtIFRoZSBkYXRhIHJlY2VpdmVkIGZyb20gdGhlIGRpYWxvZy5cclxuICogQHBhcmFtIHtib29sZWFufSBpc0pTT04gLSBGbGFnIGluZGljYXRpbmcgd2hldGhlciB0aGUgZGF0YSBpcyBpbiBKU09OIGZvcm1hdC5cclxuICpcclxuICogQHJldHVybiB7dW5kZWZpbmVkfVxyXG4gKi9cclxuZnVuY3Rpb24gZGlhbG9nUmVzdWx0Q2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGlmIChwKSB7XHJcbiAgICAgICAgaWYgKGlzSlNPTikge1xyXG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxiYWNrIGZ1bmN0aW9uIGZvciBoYW5kbGluZyBlcnJvcnMgaW4gZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xyXG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgfVxyXG59XHJcblxyXG5cclxuLy8gUmVwbGFjZSBgbWV0aG9kc2Agd2l0aCBjb25zdGFudHMgaW4gVGl0bGUgQ2FzZVxyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBJbmZvID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTtcclxuXHJcbi8qKlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgV2FybmluZyA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7XHJcblxyXG4vKipcclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IEVycm9yID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7XHJcblxyXG4vKipcclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IFF1ZXN0aW9uID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7XHJcblxyXG4vKipcclxuICogQHBhcmFtIHtPcGVuRmlsZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcy4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cclxuICovXHJcbmV4cG9ydCBjb25zdCBPcGVuRmlsZSA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7U2F2ZUZpbGVEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZS4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cclxuICovXHJcbmV4cG9ydCBjb25zdCBTYXZlRmlsZSA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nU2F2ZUZpbGUsIG9wdGlvbnMpO1xyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcclxuICovXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmltcG9ydCB7RXZlbnRUeXBlc30gZnJvbSBcIi4vZXZlbnRfdHlwZXNcIjtcclxuZXhwb3J0IGNvbnN0IFR5cGVzID0gRXZlbnRUeXBlcztcclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkV2ZW50cywgJycpO1xyXG5jb25zdCBFbWl0TWV0aG9kID0gMDtcclxuY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XHJcblxyXG5jbGFzcyBMaXN0ZW5lciB7XHJcbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcclxuICAgICAgICAgICAgY2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcclxuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcclxuICAgIGNvbnN0cnVjdG9yKG5hbWUsIGRhdGEgPSBudWxsKSB7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcclxuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gc2V0dXAoKSB7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudCkge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcclxuICAgIGlmIChsaXN0ZW5lcnMpIHtcclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+IHtcclxuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcclxuICAgICAgICAgICAgaWYgKHJlbW92ZSkgcmV0dXJuIHRydWU7XHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XHJcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XHJcbiAgICAgICAgICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxyXG4gKlxyXG4gQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XHJcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQuIEl0IHRha2VzIG5vIHBhcmFtZXRlcnMuXHJcbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7IHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTsgfVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgZXZlbnQgb2NjdXJzLlxyXG4gKiBAcmV0dXJuIHtmdW5jdGlvbn0gLSBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lLCBjYWxsYmFjaykgeyByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTsgfVxyXG5cclxuLyoqXHJcbiAqIFJlbW92ZXMgdGhlIHNwZWNpZmllZCBsaXN0ZW5lciBmcm9tIHRoZSBldmVudCBsaXN0ZW5lcnMgY29sbGVjdGlvbi5cclxuICogSWYgYWxsIGxpc3RlbmVycyBmb3IgdGhlIGV2ZW50IGFyZSByZW1vdmVkLCB0aGUgZXZlbnQga2V5IGlzIGRlbGV0ZWQgZnJvbSB0aGUgY29sbGVjdGlvbi5cclxuICpcclxuICogQHBhcmFtIHtPYmplY3R9IGxpc3RlbmVyIC0gVGhlIGxpc3RlbmVyIHRvIGJlIHJlbW92ZWQuXHJcbiAqL1xyXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xyXG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcclxuICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cclxuICogQHBhcmFtIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzIC0gQWRkaXRpb25hbCBldmVudCBuYW1lcyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cclxuICogQHJldHVybiB7dW5kZWZpbmVkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XHJcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XHJcbiAgICBldmVudHNUb1JlbW92ZS5mb3JFYWNoKGV2ZW50TmFtZSA9PiBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKSk7XHJcbn1cclxuLyoqXHJcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cclxuICpcclxuICogQGZ1bmN0aW9uIE9mZkFsbFxyXG4gKiBAcmV0dXJucyB7dm9pZH1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7IGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7IH1cclxuXHJcbi8qKlxyXG4gKiBFbWl0cyBhbiBldmVudCB1c2luZyB0aGUgZ2l2ZW4gZXZlbnQgbmFtZS5cclxuICpcclxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxyXG4gKiBAcmV0dXJucyB7YW55fSAtIFRoZSByZXN1bHQgb2YgdGhlIGVtaXR0ZWQgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkgeyByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCBldmVudCk7IH1cclxuIiwgIlxuZXhwb3J0IGNvbnN0IEV2ZW50VHlwZXMgPSB7XG5cdFdpbmRvd3M6IHtcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwid2luZG93czpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVN1c3BlbmRcIixcblx0XHRBUE1SZXN1bWVBdXRvbWF0aWM6IFwid2luZG93czpBUE1SZXN1bWVBdXRvbWF0aWNcIixcblx0XHRBUE1SZXN1bWVTdXNwZW5kOiBcIndpbmRvd3M6QVBNUmVzdW1lU3VzcGVuZFwiLFxuXHRcdEFQTVBvd2VyU2V0dGluZ0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU2V0dGluZ0NoYW5nZVwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93QWN0aXZlOiBcIndpbmRvd3M6V2luZG93QWN0aXZlXCIsXG5cdFx0V2luZG93Q2xpY2tBY3RpdmU6IFwid2luZG93czpXaW5kb3dDbGlja0FjdGl2ZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJ3aW5kb3dzOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJ3aW5kb3dzOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcIndpbmRvd3M6V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwid2luZG93czpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIndpbmRvd3M6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93U2hvdzogXCJ3aW5kb3dzOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dIaWRlOiBcIndpbmRvd3M6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd0VuZE1vdmU6IFwid2luZG93czpXaW5kb3dFbmRNb3ZlXCIsXG5cdFx0V2luZG93U3RhcnRSZXNpemU6IFwid2luZG93czpXaW5kb3dTdGFydFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0tleURvd246IFwid2luZG93czpXaW5kb3dLZXlEb3duXCIsXG5cdFx0V2luZG93S2V5VXA6IFwid2luZG93czpXaW5kb3dLZXlVcFwiLFxuXHRcdFdpbmRvd1pPcmRlckNoYW5nZWQ6IFwid2luZG93czpXaW5kb3daT3JkZXJDaGFuZ2VkXCIsXG5cdFx0V2luZG93UGFpbnQ6IFwid2luZG93czpXaW5kb3dQYWludFwiLFxuXHRcdFdpbmRvd0JhY2tncm91bmRFcmFzZTogXCJ3aW5kb3dzOldpbmRvd0JhY2tncm91bmRFcmFzZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudEhpdDogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudEhpdFwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlRG93bjogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlRG93blwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlVXA6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZVVwXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VNb3ZlOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VNb3ZlXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93RFBJQ2hhbmdlZFwiLFxuXHR9LFxuXHRNYWM6IHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZU5vdGlmaWNhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVOb3RpZmljYXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZSFcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW4hXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwibWFjOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJtYWM6V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb20hXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcIm1hYzpXaW5kb3dab29tSW4hXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJtYWM6V2luZG93Wm9vbU91dCFcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldCFcIixcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxuXHRcdFdpbmRvd0RpZEZvY3VzOiBcIm1hYzpXaW5kb3dEaWRGb2N1c1wiLFxuXHRcdFdpbmRvd0RpZE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWRNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPZmZTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcblx0XHRXaW5kb3dEaWRSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIm1hYzpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93RGlkVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlIVwiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdE1lbnVXaWxsT3BlbjogXCJtYWM6TWVudVdpbGxPcGVuXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWREaXNwbGF5SXRlbTogXCJtYWM6TWVudURpZERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVEaWRSZW1vdmVJdGVtOiBcIm1hYzpNZW51RGlkUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsQmVnaW5UcmFja2luZzogXCJtYWM6TWVudVdpbGxCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51RGlkUG9wVXA6IFwibWFjOk1lbnVEaWRQb3BVcFwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dTaG93OiBcIm1hYzpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93SGlkZTogXCJtYWM6V2luZG93SGlkZVwiLFxuXHR9LFxuXHRMaW51eDoge1xuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRXaW5kb3dMb2FkQ2hhbmdlZDogXCJsaW51eDpXaW5kb3dMb2FkQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0RlbGV0ZUV2ZW50OiBcImxpbnV4OldpbmRvd0RlbGV0ZUV2ZW50XCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJsaW51eDpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImxpbnV4OldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0ZvY3VzSW46IFwibGludXg6V2luZG93Rm9jdXNJblwiLFxuXHRcdFdpbmRvd0ZvY3VzT3V0OiBcImxpbnV4OldpbmRvd0ZvY3VzT3V0XCIsXG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxuXHR9LFxuXHRDb21tb246IHtcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwiY29tbW9uOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcImNvbW1vbjpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJjb21tb246V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwiY29tbW9uOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcImNvbW1vbjpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcblx0XHRXaW5kb3dGb2N1czogXCJjb21tb246V2luZG93Rm9jdXNcIixcblx0XHRXaW5kb3dMb3N0Rm9jdXM6IFwiY29tbW9uOldpbmRvd0xvc3RGb2N1c1wiLFxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dIaWRlOiBcImNvbW1vbjpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RmlsZXNEcm9wcGVkXCIsXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiY29tbW9uOlRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwiY29tbW9uOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwiY29tbW9uOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiY29tbW9uOkFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGVcIixcblx0fSxcbn07XG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIExvZ3MgYSBtZXNzYWdlIHRvIHRoZSBjb25zb2xlIHdpdGggY3VzdG9tIGZvcm1hdHRpbmcuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIG1lc3NhZ2UgdG8gYmUgbG9nZ2VkLlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIGRlYnVnTG9nKG1lc3NhZ2UpIHtcclxuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG4gICAgY29uc29sZS5sb2coXHJcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXHJcbiAgICApO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIHdoZXRoZXIgdGhlIGJyb3dzZXIgc3VwcG9ydHMgcmVtb3ZpbmcgbGlzdGVuZXJzIGJ5IHRyaWdnZXJpbmcgYW4gQWJvcnRTaWduYWxcclxuICogKHNlZSBodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZy9lbi1VUy9kb2NzL1dlYi9BUEkvRXZlbnRUYXJnZXQvYWRkRXZlbnRMaXN0ZW5lciNzaWduYWwpXHJcbiAqXHJcbiAqIEByZXR1cm4ge2Jvb2xlYW59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XHJcbiAgICBpZiAoIUV2ZW50VGFyZ2V0IHx8ICFBYm9ydFNpZ25hbCB8fCAhQWJvcnRDb250cm9sbGVyKVxyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuXHJcbiAgICBsZXQgcmVzdWx0ID0gdHJ1ZTtcclxuXHJcbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcclxuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB0YXJnZXQuYWRkRXZlbnRMaXN0ZW5lcigndGVzdCcsICgpID0+IHsgcmVzdWx0ID0gZmFsc2U7IH0sIHsgc2lnbmFsOiBjb250cm9sbGVyLnNpZ25hbCB9KTtcclxuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcclxuICAgIHRhcmdldC5kaXNwYXRjaEV2ZW50KG5ldyBDdXN0b21FdmVudCgndGVzdCcpKTtcclxuXHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKioqXHJcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcclxuXHJcbiBCU0QgMi1DbGF1c2UgTGljZW5zZVxyXG5cclxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxyXG4gQWxsIHJpZ2h0cyByZXNlcnZlZC5cclxuXHJcbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcclxuIG1vZGlmaWNhdGlvbiwgYXJlIHBlcm1pdHRlZCBwcm92aWRlZCB0aGF0IHRoZSBmb2xsb3dpbmcgY29uZGl0aW9ucyBhcmUgbWV0OlxyXG5cclxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xyXG4gbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIuXHJcblxyXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxyXG4gdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvblxyXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXHJcblxyXG4gVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyBcIkFTIElTXCJcclxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcclxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxyXG4gRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1JTIEJFIExJQUJMRVxyXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcclxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXHJcbiBTRVJWSUNFUzsgTE9TUyBPRiBVU0UsIERBVEEsIE9SIFBST0ZJVFM7IE9SIEJVU0lORVNTIElOVEVSUlVQVElPTikgSE9XRVZFUlxyXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcclxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXHJcbiBPRiBUSElTIFNPRlRXQVJFLCBFVkVOIElGIEFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLlxyXG5cclxuICoqKi9cclxuXHJcbmxldCBpc1JlYWR5ID0gZmFsc2U7XHJcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCAoKSA9PiBpc1JlYWR5ID0gdHJ1ZSk7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gd2hlblJlYWR5KGNhbGxiYWNrKSB7XHJcbiAgICBpZiAoaXNSZWFkeSB8fCBkb2N1bWVudC5yZWFkeVN0YXRlID09PSAnY29tcGxldGUnKSB7XHJcbiAgICAgICAgY2FsbGJhY2soKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGNhbGxiYWNrKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLy8gSW1wb3J0IHNjcmVlbiBqc2RvYyBkZWZpbml0aW9uIGZyb20gLi9zY3JlZW5zLmpzXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9zY3JlZW5zXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuXHJcbi8qKlxyXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBwb3NpdGlvbiBvZiBhIHdpbmRvdy5cclxuICpcclxuICogQHR5cGVkZWYge09iamVjdH0gUG9zaXRpb25cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IHggLSBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSB5IC0gVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcclxuICovXHJcblxyXG5cclxuLyoqXHJcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXHJcbiAqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNpemVcclxuICogQHByb3BlcnR5IHtudW1iZXJ9IHdpZHRoIC0gVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3dcclxuICogQHByb3BlcnR5IHtudW1iZXJ9IGhlaWdodCAtIFRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvd1xyXG4gKi9cclxuXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5jb25zdCBQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAwO1xyXG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xyXG5jb25zdCBDbG9zZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSAyO1xyXG5jb25zdCBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgPSAzO1xyXG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xyXG5jb25zdCBGb2N1c01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1O1xyXG5jb25zdCBGb3JjZVJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgPSA2O1xyXG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xyXG5jb25zdCBHZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA4O1xyXG5jb25zdCBHZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA5O1xyXG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcclxuY29uc3QgSGlkZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTE7XHJcbmNvbnN0IElzRm9jdXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDEyO1xyXG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcclxuY29uc3QgSXNNYXhpbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTQ7XHJcbmNvbnN0IElzTWluaW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE1O1xyXG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcclxuY29uc3QgTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTc7XHJcbmNvbnN0IE5hbWVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDE4O1xyXG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcclxuY29uc3QgUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMjA7XHJcbmNvbnN0IFJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDIxO1xyXG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcclxuY29uc3QgUmVzdG9yZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMjM7XHJcbmNvbnN0IFNldFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICA9IDI0O1xyXG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcclxuY29uc3QgU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCAgICAgICAgID0gMjY7XHJcbmNvbnN0IFNldEZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICAgICA9IDI3O1xyXG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcclxuY29uc3QgU2V0TWF4U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMjk7XHJcbmNvbnN0IFNldE1pblNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDMwO1xyXG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcclxuY29uc3QgU2V0UmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgID0gMzI7XHJcbmNvbnN0IFNldFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDMzO1xyXG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcclxuY29uc3QgU2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzU7XHJcbmNvbnN0IFNob3dNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM2O1xyXG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcclxuY29uc3QgVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgID0gMzg7XHJcbmNvbnN0IFRvZ2dsZU1heGltaXNlTWV0aG9kICAgICAgICAgICAgICA9IDM5O1xyXG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MDtcclxuY29uc3QgVW5NYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDE7XHJcbmNvbnN0IFVuTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQyO1xyXG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0MztcclxuY29uc3QgWm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XHJcbmNvbnN0IFpvb21Jbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDQ1O1xyXG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NjtcclxuY29uc3QgWm9vbVJlc2V0TWV0aG9kICAgICAgICAgICAgICAgICAgID0gNDc7XHJcblxyXG4vKipcclxuICogQHR5cGUge3N5bWJvbH1cclxuICovXHJcbmNvbnN0IGNhbGxlciA9IFN5bWJvbCgpO1xyXG5cclxuZXhwb3J0IGNsYXNzIFdpbmRvdyB7XHJcbiAgICAvKipcclxuICAgICAqIEluaXRpYWxpc2VzIGEgd2luZG93IG9iamVjdCB3aXRoIHRoZSBzcGVjaWZpZWQgbmFtZS5cclxuICAgICAqXHJcbiAgICAgKiBAcHJpdmF0ZVxyXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgdGFyZ2V0IHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IobmFtZSA9ICcnKSB7XHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogQHByaXZhdGVcclxuICAgICAgICAgKiBAbmFtZSB7QGxpbmsgY2FsbGVyfVxyXG4gICAgICAgICAqIEB0eXBlIHsoLi4uYXJnczogYW55W10pID0+IGFueX1cclxuICAgICAgICAgKi9cclxuICAgICAgICB0aGlzW2NhbGxlcl0gPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcclxuXHJcbiAgICAgICAgLy8gYmluZCBpbnN0YW5jZSBtZXRob2QgdG8gbWFrZSB0aGVtIGVhc2lseSB1c2FibGUgaW4gZXZlbnQgaGFuZGxlcnNcclxuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xyXG4gICAgICAgICAgICBpZiAoXHJcbiAgICAgICAgICAgICAgICBtZXRob2QgIT09IFwiY29uc3RydWN0b3JcIlxyXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mIHRoaXNbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXHJcbiAgICAgICAgICAgICkge1xyXG4gICAgICAgICAgICAgICAgdGhpc1ttZXRob2RdID0gdGhpc1ttZXRob2RdLmJpbmQodGhpcyk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBHZXRzIHRoZSBzcGVjaWZpZWQgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBnZXQuXHJcbiAgICAgKiBAcmV0dXJuIHtXaW5kb3d9IC0gVGhlIGNvcnJlc3BvbmRpbmcgd2luZG93IG9iamVjdC5cclxuICAgICAqL1xyXG4gICAgR2V0KG5hbWUpIHtcclxuICAgICAgICByZXR1cm4gbmV3IFdpbmRvdyhuYW1lKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxQb3NpdGlvbj59IC0gVGhlIGN1cnJlbnQgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgUG9zaXRpb24oKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShQb3NpdGlvbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBDZW50ZXIoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShDZW50ZXJNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgQ2xvc2UoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShDbG9zZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEaXNhYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgRGlzYWJsZVNpemVDb25zdHJhaW50cygpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRW5hYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBGb2N1cygpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEZvY3VzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEZvcmNlcyB0aGUgd2luZG93IHRvIHJlbG9hZCB0aGUgcGFnZSBhc3NldHMuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgRm9yY2VSZWxvYWQoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShGb3JjZVJlbG9hZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEb2MuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgRnVsbHNjcmVlbigpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgc2NyZWVuIHRoYXQgdGhlIHdpbmRvdyBpcyBvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPFNjcmVlbj59IC0gVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvblxyXG4gICAgICovXHJcbiAgICBHZXRTY3JlZW4oKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShHZXRTY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgY3VycmVudCB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxudW1iZXI+fSAtIFRoZSBjdXJyZW50IHpvb20gbGV2ZWxcclxuICAgICAqL1xyXG4gICAgR2V0Wm9vbSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEdldFpvb21NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxudW1iZXI+fSAtIFRoZSBjdXJyZW50IGhlaWdodCBvZiB0aGUgd2luZG93XHJcbiAgICAgKi9cclxuICAgIEhlaWdodCgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEhlaWdodE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBIaWRlcyB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIEhpZGUoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShIaWRlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZvY3VzZWQuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gLSBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZvY3VzZWRcclxuICAgICAqL1xyXG4gICAgSXNGb2N1c2VkKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSXNGb2N1c2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gLSBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW5cclxuICAgICAqL1xyXG4gICAgSXNGdWxsc2NyZWVuKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSXNGdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1heGltaXNlZC5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWF4aW1pc2VkXHJcbiAgICAgKi9cclxuICAgIElzTWF4aW1pc2VkKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSXNNYXhpbWlzZWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IC0gV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtaW5pbWlzZWRcclxuICAgICAqL1xyXG4gICAgSXNNaW5pbWlzZWQoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShJc01pbmltaXNlZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBNYXhpbWlzZSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKE1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIE1pbmltaXNlKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oTWluaW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93XHJcbiAgICAgKi9cclxuICAgIE5hbWUoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShOYW1lTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIE9wZW5EZXZUb29scygpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKE9wZW5EZXZUb29sc01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxQb3NpdGlvbj59IC0gVGhlIGN1cnJlbnQgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gICAgICovXHJcbiAgICBSZWxhdGl2ZVBvc2l0aW9uKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZWxvYWRzIHRoZSBwYWdlIGFzc2V0cy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBSZWxvYWQoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShSZWxvYWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IC0gV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGVcclxuICAgICAqL1xyXG4gICAgUmVzaXphYmxlKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oUmVzaXphYmxlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlc3RvcmVzIHRoZSB3aW5kb3cgdG8gaXRzIHByZXZpb3VzIHN0YXRlIGlmIGl0IHdhcyBwcmV2aW91c2x5IG1pbmltaXNlZCwgbWF4aW1pc2VkIG9yIGZ1bGxzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgUmVzdG9yZSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFJlc3RvcmVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTZXRQb3NpdGlvbih4LCB5KSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHBhcmFtIHtib29sZWFufSBhbHdheXNPblRvcCAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgc3RheSBvbiB0b3BcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFNldEFsd2F5c09uVG9wKGFsd2F5c09uVG9wKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRBbHdheXNPblRvcE1ldGhvZCwgeyBhbHdheXNPblRvcCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHIgLSBUaGUgZGVzaXJlZCByZWQgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZFxyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IGcgLSBUaGUgZGVzaXJlZCBncmVlbiBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZFxyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IGEgLSBUaGUgZGVzaXJlZCBhbHBoYSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyKHIsIGcsIGIsIGEpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QsIHsgciwgZywgYiwgYSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbW92ZXMgdGhlIHdpbmRvdyBmcmFtZSBhbmQgdGl0bGUgYmFyLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3NcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFNldEZyYW1lbGVzcyhmcmFtZWxlc3MpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEaXNhYmxlcyB0aGUgc3lzdGVtIGZ1bGxzY3JlZW4gYnV0dG9uLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gZW5hYmxlZCAtIFdoZXRoZXIgdGhlIGZ1bGxzY3JlZW4gYnV0dG9uIHNob3VsZCBiZSBlbmFibGVkXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZChlbmFibGVkKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCwgeyBlbmFibGVkIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIGRlc2lyZWQgbWF4aW11bSB3aWR0aCBvZiB0aGUgd2luZG93XHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvd1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgU2V0TWF4U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRNYXhTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGggLSBUaGUgZGVzaXJlZCBtaW5pbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3dcclxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtaW5pbXVtIGhlaWdodCBvZiB0aGUgd2luZG93XHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTZXRNaW5TaXplKHdpZHRoLCBoZWlnaHQpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldE1pblNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93XHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uKHgsIHkpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcGFyYW0ge2Jvb2xlYW59IHJlc2l6YWJsZSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgcmVzaXphYmxlXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRSZXNpemFibGVNZXRob2QsIHsgcmVzaXphYmxlIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCAtIFRoZSBkZXNpcmVkIHdpZHRoIG9mIHRoZSB3aW5kb3dcclxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvd1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgU2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSB0aXRsZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSB0aXRsZSAtIFRoZSBkZXNpcmVkIHRpdGxlIG9mIHRoZSB3aW5kb3dcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFNldFRpdGxlKHRpdGxlKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWxcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFNldFpvb20oem9vbSkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0Wm9vbU1ldGhvZCwgeyB6b29tIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBTaG93KCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2hvd01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxTaXplPn0gLSBUaGUgY3VycmVudCBzaXplIG9mIHRoZSB3aW5kb3dcclxuICAgICAqL1xyXG4gICAgU2l6ZSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNpemVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnVsbHNjcmVlbiBhbmQgbm9ybWFsLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZUZ1bGxzY3JlZW4oKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShUb2dnbGVGdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZU1heGltaXNlKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oVG9nZ2xlTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tZnVsbHNjcmVlbnMgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gICAgICovXHJcbiAgICBVbkZ1bGxzY3JlZW4oKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgVW5NYXhpbWlzZSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFVuTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgVW5NaW5pbWlzZSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFVuTWluaW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcHVibGljXHJcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPG51bWJlcj59IC0gVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvd1xyXG4gICAgICovXHJcbiAgICBXaWR0aCgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFdpZHRoTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFpvb21zIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgWm9vbSgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFpvb21NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXHJcbiAgICAgKlxyXG4gICAgICogQHB1YmxpY1xyXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICAgICAqL1xyXG4gICAgWm9vbUluKCkge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oWm9vbUluTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIERlY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFpvb21PdXQoKSB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShab29tT3V0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlc2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICpcclxuICAgICAqIEBwdWJsaWNcclxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAgICAgKi9cclxuICAgIFpvb21SZXNldCgpIHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFpvb21SZXNldE1ldGhvZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXHJcbiAqXHJcbiAqIEB0eXBlIHtXaW5kb3d9XHJcbiAqL1xyXG5jb25zdCB0aGlzV2luZG93ID0gbmV3IFdpbmRvdygnJyk7XHJcblxyXG5leHBvcnQgZGVmYXVsdCB0aGlzV2luZG93O1xyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0ICogYXMgUnVudGltZSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmNcIjtcclxuXHJcbi8vIE5PVEU6IHRoZSBmb2xsb3dpbmcgbWV0aG9kcyBNVVNUIGJlIGltcG9ydGVkIGV4cGxpY2l0bHkgYmVjYXVzZSBvZiBob3cgZXNidWlsZCBpbmplY3Rpb24gd29ya3NcclxuaW1wb3J0IHtFbmFibGUgYXMgRW5hYmxlV01MfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XHJcbmltcG9ydCB7ZGVidWdMb2d9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy91dGlsc1wiO1xyXG5cclxud2luZG93LndhaWxzID0gUnVudGltZTtcclxuRW5hYmxlV01MKCk7XHJcblxyXG5pZiAoREVCVUcpIHtcclxuICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIilcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU3lzdGVtLCAnJyk7XHJcbmNvbnN0IHN5c3RlbUlzRGFya01vZGUgPSAwO1xyXG5jb25zdCBlbnZpcm9ubWVudCA9IDE7XHJcblxyXG5jb25zdCBfaW52b2tlID0gKCgpID0+IHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgaWYod2luZG93Py5jaHJvbWU/LndlYnZpZXcpIHtcclxuICAgICAgICAgICAgcmV0dXJuIChtc2cpID0+IHdpbmRvdy5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZShtc2cpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBpZih3aW5kb3c/LndlYmtpdD8ubWVzc2FnZUhhbmRsZXJzPy5leHRlcm5hbCkge1xyXG4gICAgICAgICAgICByZXR1cm4gKG1zZykgPT4gd2luZG93LndlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2UobXNnKTtcclxuICAgICAgICB9XHJcbiAgICB9IGNhdGNoKGUpIHtcclxuICAgICAgICBjb25zb2xlLndhcm4oJ1xcbiVjXHUyNkEwXHVGRTBGIEJyb3dzZXIgRW52aXJvbm1lbnQgRGV0ZWN0ZWQgJWNcXG5cXG4lY09ubHkgVUkgcHJldmlld3MgYXJlIGF2YWlsYWJsZSBpbiB0aGUgYnJvd3Nlci4gRm9yIGZ1bGwgZnVuY3Rpb25hbGl0eSwgcGxlYXNlIHJ1biB0aGUgYXBwbGljYXRpb24gaW4gZGVza3RvcCBtb2RlLlxcbk1vcmUgaW5mb3JtYXRpb24gYXQ6IGh0dHBzOi8vdjMud2FpbHMuaW8vbGVhcm4vYnVpbGQvI3VzaW5nLWEtYnJvd3Nlci1mb3ItZGV2ZWxvcG1lbnRcXG4nLFxyXG4gICAgICAgICAgICAnYmFja2dyb3VuZDogI2ZmZmZmZjsgY29sb3I6ICMwMDAwMDA7IGZvbnQtd2VpZ2h0OiBib2xkOyBwYWRkaW5nOiA0cHggOHB4OyBib3JkZXItcmFkaXVzOiA0cHg7IGJvcmRlcjogMnB4IHNvbGlkICMwMDAwMDA7JyxcclxuICAgICAgICAgICAgJ2JhY2tncm91bmQ6IHRyYW5zcGFyZW50OycsXHJcbiAgICAgICAgICAgICdjb2xvcjogI2ZmZmZmZjsgZm9udC1zdHlsZTogaXRhbGljOyBmb250LXdlaWdodDogYm9sZDsnKTtcclxuICAgIH1cclxuICAgIHJldHVybiBudWxsO1xyXG59KSgpO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2cpIHtcclxuICAgIGlmICghX2ludm9rZSkgcmV0dXJuO1xyXG4gICAgcmV0dXJuIF9pbnZva2UobXNnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEBmdW5jdGlvblxyXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxib29sZWFuPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoc3lzdGVtSXNEYXJrTW9kZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBGZXRjaGVzIHRoZSBjYXBhYmlsaXRpZXMgb2YgdGhlIGFwcGxpY2F0aW9uIGZyb20gdGhlIHNlcnZlci5cclxuICpcclxuICogQGFzeW5jXHJcbiAqIEBmdW5jdGlvbiBDYXBhYmlsaXRpZXNcclxuICogQHJldHVybnMge1Byb21pc2U8T2JqZWN0Pn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKSB7XHJcbiAgICBsZXQgcmVzcG9uc2UgPSBmZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIik7XHJcbiAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG59XHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gT1NJbmZvXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBCcmFuZGluZyAtIFRoZSBicmFuZGluZyBvZiB0aGUgT1MuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJRCAtIFRoZSBJRCBvZiB0aGUgT1MuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIE9TLlxyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gVmVyc2lvbiAtIFRoZSB2ZXJzaW9uIG9mIHRoZSBPUy5cclxuICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gRW52aXJvbm1lbnRJbmZvXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBBcmNoIC0gVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLlxyXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IERlYnVnIC0gVHJ1ZSBpZiB0aGUgYXBwbGljYXRpb24gaXMgcnVubmluZyBpbiBkZWJ1ZyBtb2RlLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBPUyAtIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS5cclxuICogQHByb3BlcnR5IHtPU0luZm99IE9TSW5mbyAtIERldGFpbHMgb2YgdGhlIG9wZXJhdGluZyBzeXN0ZW0uXHJcbiAqIEBwcm9wZXJ0eSB7T2JqZWN0fSBQbGF0Zm9ybUluZm8gLSBBZGRpdGlvbmFsIHBsYXRmb3JtIGluZm9ybWF0aW9uLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAZnVuY3Rpb25cclxuICogUmV0cmlldmVzIGVudmlyb25tZW50IGRldGFpbHMuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPEVudmlyb25tZW50SW5mbz59IC0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChlbnZpcm9ubWVudCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtib29sZWFufSBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MsIG90aGVyd2lzZSBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc1dpbmRvd3MoKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJ3aW5kb3dzXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59IFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNMaW51eCgpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImxpbnV4XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImRhcndpblwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cclxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYW1kNjRcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBUk0oKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFybVwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59IC0gUmV0dXJucyB0cnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUsIG90aGVyd2lzZSByZXR1cm5zIGZhbHNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKSB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFybTY0XCI7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RlYnVnKCkge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuRGVidWcgPT09IHRydWU7XHJcbn1cclxuXHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmltcG9ydCB7SXNEZWJ1Z30gZnJvbSBcIi4vc3lzdGVtXCI7XHJcblxyXG4vLyBzZXR1cFxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUsICcnKTtcclxuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcclxuXHJcbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xyXG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudCkge1xyXG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcclxuICAgIGxldCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xyXG4gICAgbGV0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpO1xyXG4gICAgY3VzdG9tQ29udGV4dE1lbnUgPSBjdXN0b21Db250ZXh0TWVudSA/IGN1c3RvbUNvbnRleHRNZW51LnRyaW0oKSA6IFwiXCI7XHJcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIGxldCBjdXN0b21Db250ZXh0TWVudURhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnUtZGF0YVwiKTtcclxuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGN1c3RvbUNvbnRleHRNZW51RGF0YSk7XHJcbiAgICAgICAgcmV0dXJuXHJcbiAgICB9XHJcblxyXG4gICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCk7XHJcbn1cclxuXHJcblxyXG4vKlxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuXHJcblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXHJcbiovXHJcbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcclxuXHJcbiAgICAvLyBEZWJ1ZyBidWlsZHMgYWx3YXlzIHNob3cgdGhlIG1lbnVcclxuICAgIGlmIChJc0RlYnVnKCkpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG4gICAgY29uc3QgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcclxuICAgIGNvbnN0IGNvbXB1dGVkU3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcclxuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XHJcbiAgICBzd2l0Y2ggKGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbikge1xyXG4gICAgICAgIGNhc2UgXCJzaG93XCI6XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICBjYXNlIFwiaGlkZVwiOlxyXG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgZGVmYXVsdDpcclxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcclxuICAgICAgICAgICAgaWYgKGVsZW1lbnQuaXNDb250ZW50RWRpdGFibGUpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxyXG4gICAgICAgICAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XHJcbiAgICAgICAgICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IChzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKVxyXG4gICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcclxuICAgICAgICAgICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xyXG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcclxuICAgICAgICAgICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbiAgICAgICAgICAgIGlmIChlbGVtZW50LnRhZ05hbWUgPT09IFwiSU5QVVRcIiB8fCBlbGVtZW50LnRhZ05hbWUgPT09IFwiVEVYVEFSRUFcIikge1xyXG4gICAgICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIWVsZW1lbnQucmVhZE9ubHkgJiYgIWVsZW1lbnQuZGlzYWJsZWQpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgdGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleSBmcm9tIHRoZSBmbGFnIG1hcC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGtleVN0cmluZyAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cclxuICogQHJldHVybiB7Kn0gLSBUaGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5U3RyaW5nKSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleVN0cmluZ107XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiVW5hYmxlIHRvIHJldHJpZXZlIGZsYWcgJ1wiICsga2V5U3RyaW5nICsgXCInOiBcIiArIGUpO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0IHtpbnZva2UsIElzV2luZG93c30gZnJvbSBcIi4vc3lzdGVtXCI7XHJcbmltcG9ydCB7R2V0RmxhZ30gZnJvbSBcIi4vZmxhZ3NcIjtcclxuXHJcbi8vIFNldHVwXHJcbmxldCBzaG91bGREcmFnID0gZmFsc2U7XHJcbmxldCByZXNpemFibGUgPSBmYWxzZTtcclxubGV0IHJlc2l6ZUVkZ2UgPSBudWxsO1xyXG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xyXG5cclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9IGZ1bmN0aW9uKHZhbHVlKSB7XHJcbiAgICByZXNpemFibGUgPSB2YWx1ZTtcclxufTtcclxuXHJcbndpbmRvdy5fd2FpbHMuZW5kRHJhZyA9IGZ1bmN0aW9uKCkge1xyXG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XHJcbiAgICBzaG91bGREcmFnID0gZmFsc2U7XHJcbn07XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgb25Nb3VzZURvd24pO1xyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgb25Nb3VzZU1vdmUpO1xyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIG9uTW91c2VVcCk7XHJcblxyXG5cclxuZnVuY3Rpb24gZHJhZ1Rlc3QoZSkge1xyXG4gICAgbGV0IHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13YWlscy1kcmFnZ2FibGVcIik7XHJcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xyXG4gICAgaWYgKCF2YWwgfHwgdmFsID09PSBcIlwiIHx8IHZhbC50cmltKCkgIT09IFwiZHJhZ1wiIHx8IG1vdXNlUHJlc3NlZCA9PT0gMCkge1xyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgIH1cclxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xyXG5cclxuICAgIC8vIENoZWNrIGZvciByZXNpemluZ1xyXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICBpbnZva2UoXCJ3YWlsczpyZXNpemU6XCIgKyByZXNpemVFZGdlKTtcclxuICAgICAgICBlLnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChkcmFnVGVzdChlKSkge1xyXG4gICAgICAgIC8vIFRoaXMgY2hlY2tzIGZvciBjbGlja3Mgb24gdGhlIHNjcm9sbCBiYXJcclxuICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICAgICAgc2hvdWxkRHJhZyA9IHRydWU7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHNob3VsZERyYWcgPSBmYWxzZTtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZVVwKCkge1xyXG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG59XHJcblxyXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XHJcbiAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gY3Vyc29yIHx8IGRlZmF1bHRDdXJzb3I7XHJcbiAgICByZXNpemVFZGdlID0gY3Vyc29yO1xyXG59XHJcblxyXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShlKSB7XHJcbiAgICBpZiAoc2hvdWxkRHJhZykge1xyXG4gICAgICAgIHNob3VsZERyYWcgPSBmYWxzZTtcclxuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xyXG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XHJcbiAgICAgICAgICAgIGludm9rZShcIndhaWxzOmRyYWdcIik7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbiAgICBpZiAoIXJlc2l6YWJsZSB8fCAhSXNXaW5kb3dzKCkpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcbiAgICBpZiAoZGVmYXVsdEN1cnNvciA9PSBudWxsKSB7XHJcbiAgICAgICAgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3I7XHJcbiAgICB9XHJcbiAgICBsZXQgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcclxuICAgIGxldCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcclxuXHJcbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXNcclxuICAgIGxldCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcclxuXHJcbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xyXG5cclxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyc1xyXG4gICAgbGV0IHJpZ2h0Q29ybmVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBsZXQgbGVmdENvcm5lciA9IGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGxldCB0b3BDb3JuZXIgPSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG4gICAgbGV0IGJvdHRvbUNvcm5lciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XHJcblxyXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcclxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xyXG4gICAgICAgIHNldFJlc2l6ZSgpO1xyXG4gICAgfVxyXG4gICAgLy8gQWRqdXN0ZWQgZm9yIGNvcm5lciBhcmVhc1xyXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcclxufSIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24sICcnKTtcclxuXHJcbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xyXG5jb25zdCBTaG93TWV0aG9kID0gMTtcclxuY29uc3QgUXVpdE1ldGhvZCA9IDI7XHJcblxyXG4vKipcclxuICogSGlkZXMgYSBjZXJ0YWluIG1ldGhvZCBieSBjYWxsaW5nIHRoZSBIaWRlTWV0aG9kIGZ1bmN0aW9uLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gKlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XHJcbiAgICByZXR1cm4gY2FsbChIaWRlTWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXHJcbiAqXHJcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcclxuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xyXG5cclxuLy8gU2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbndpbmRvdy5fd2FpbHMuY2FsbFJlc3VsdEhhbmRsZXIgPSByZXN1bHRIYW5kbGVyO1xyXG53aW5kb3cuX3dhaWxzLmNhbGxFcnJvckhhbmRsZXIgPSBlcnJvckhhbmRsZXI7XHJcblxyXG5cclxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYWxsLCAnJyk7XHJcbmNvbnN0IGNhbmNlbENhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNhbmNlbENhbGwsICcnKTtcclxubGV0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwKCk7XHJcblxyXG4vKipcclxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cclxuICpcclxuICogQHJldHVybiB7c3RyaW5nfSAtIEEgdW5pcXVlIElEIHRoYXQgZG9lcyBub3QgZXhpc3QgaW4gdGhlIGNhbGxSZXNwb25zZXMgc2V0LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBjYWxsIHJlcXVlc3QuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBpZCBvZiB0aGUgcmVxdWVzdCB0byBoYW5kbGUgdGhlIHJlc3VsdCBmb3IuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxyXG4gKiBAcGFyYW0ge2Jvb2xlYW59IGlzSlNPTiAtIEluZGljYXRlcyB3aGV0aGVyIHRoZSBkYXRhIGlzIEpTT04gb3Igbm90LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9IC0gVGhpcyBtZXRob2QgZG9lcyBub3QgcmV0dXJuIGFueSB2YWx1ZS5cclxuICovXHJcbmZ1bmN0aW9uIHJlc3VsdEhhbmRsZXIoaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XHJcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcclxuICAgICAgICBwcm9taXNlSGFuZGxlci5yZXNvbHZlKGlzSlNPTiA/IEpTT04ucGFyc2UoZGF0YSkgOiBkYXRhKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIGVycm9yIGZyb20gYSBjYWxsIHJlcXVlc3QuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlIHRvIHJlamVjdCB0aGUgcHJvbWlzZSBoYW5kbGVyIHdpdGguXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBlcnJvckhhbmRsZXIoaWQsIG1lc3NhZ2UpIHtcclxuICAgIGNvbnN0IHByb21pc2VIYW5kbGVyID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xyXG4gICAgaWYgKHByb21pc2VIYW5kbGVyKSB7XHJcbiAgICAgICAgcHJvbWlzZUhhbmRsZXIucmVqZWN0KG1lc3NhZ2UpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmV0cmlldmVzIGFuZCByZW1vdmVzIHRoZSByZXNwb25zZSBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElEIGZyb20gdGhlIGNhbGxSZXNwb25zZXMgbWFwLlxyXG4gKlxyXG4gKiBAcGFyYW0ge2FueX0gaWQgLSBUaGUgSUQgb2YgdGhlIHJlc3BvbnNlIHRvIGJlIHJldHJpZXZlZCBhbmQgcmVtb3ZlZC5cclxuICpcclxuICogQHJldHVybnMge2FueX0gVGhlIHJlc3BvbnNlIG9iamVjdCBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElELlxyXG4gKi9cclxuZnVuY3Rpb24gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpIHtcclxuICAgIGNvbnN0IHJlc3BvbnNlID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgcmV0dXJuIHJlc3BvbnNlO1xyXG59XHJcblxyXG4vKipcclxuICogRXhlY3V0ZXMgYSBjYWxsIHVzaW5nIHRoZSBwcm92aWRlZCB0eXBlIGFuZCBvcHRpb25zLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ3xudW1iZXJ9IHR5cGUgLSBUaGUgdHlwZSBvZiBjYWxsIHRvIGV4ZWN1dGUuXHJcbiAqIEBwYXJhbSB7T2JqZWN0fSBbb3B0aW9ucz17fV0gLSBBZGRpdGlvbmFsIG9wdGlvbnMgZm9yIHRoZSBjYWxsLlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgYmFzZWQgb24gdGhlIHJlc3VsdCBvZiB0aGUgY2FsbC4gSXQgYWxzbyBoYXMgYSBjYW5jZWwgbWV0aG9kIHRvIGNhbmNlbCBhIGxvbmcgcnVubmluZyByZXF1ZXN0LlxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucyA9IHt9KSB7XHJcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuICAgIGNvbnN0IGRvQ2FuY2VsID0gKCkgPT4geyByZXR1cm4gY2FuY2VsQ2FsbCh0eXBlLCB7XCJjYWxsLWlkXCI6IGlkfSkgfTtcclxuICAgIGxldCBxdWV1ZWRDYW5jZWwgPSBmYWxzZSwgY2FsbFJ1bm5pbmcgPSBmYWxzZTtcclxuICAgIGxldCBwID0gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuXHJcbiAgICAgICAgICAgIHRoZW4oKF8pID0+IHtcclxuICAgICAgICAgICAgICAgIGNhbGxSdW5uaW5nID0gdHJ1ZTtcclxuICAgICAgICAgICAgICAgIGlmIChxdWV1ZWRDYW5jZWwpIHtcclxuICAgICAgICAgICAgICAgICAgICByZXR1cm4gZG9DYW5jZWwoKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSkuXHJcbiAgICAgICAgICAgIGNhdGNoKChlcnJvcikgPT4ge1xyXG4gICAgICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcclxuICAgICAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICAgICAgfSk7XHJcbiAgICB9KTtcclxuICAgIHAuY2FuY2VsID0gKCkgPT4ge1xyXG4gICAgICAgIGlmIChjYWxsUnVubmluZykge1xyXG4gICAgICAgICAgICByZXR1cm4gZG9DYW5jZWwoKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBxdWV1ZWRDYW5jZWwgPSB0cnVlO1xyXG4gICAgICAgIH1cclxuICAgIH07XHJcblxyXG4gICAgcmV0dXJuIHA7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsIG1ldGhvZC5cclxuICpcclxuICogQHBhcmFtIHtPYmplY3R9IG9wdGlvbnMgLSBUaGUgb3B0aW9ucyBmb3IgdGhlIG1ldGhvZC5cclxuICogQHJldHVybnMge09iamVjdH0gLSBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEV4ZWN1dGVzIGEgbWV0aG9kIGJ5IG5hbWUuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxyXG4gKiBAcGFyYW0gey4uLip9IGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cclxuICogQHRocm93cyB7RXJyb3J9IElmIHRoZSBuYW1lIGlzIG5vdCBhIHN0cmluZyBvciBpcyBub3QgaW4gdGhlIGNvcnJlY3QgZm9ybWF0LlxyXG4gKiBAcmV0dXJucyB7Kn0gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGV4ZWN1dGlvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeU5hbWUobWV0aG9kTmFtZSwgLi4uYXJncykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgbWV0aG9kTmFtZSxcclxuICAgICAgICBhcmdzXHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIGJ5IGl0cyBJRCB3aXRoIHRoZSBzcGVjaWZpZWQgYXJndW1lbnRzLlxyXG4gKlxyXG4gKiBAcGFyYW0ge251bWJlcn0gbWV0aG9kSUQgLSBUaGUgSUQgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxyXG4gKiBAcGFyYW0gey4uLip9IGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cclxuICogQHJldHVybiB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeUlEKG1ldGhvZElELCAuLi5hcmdzKSB7XHJcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIHtcclxuICAgICAgICBtZXRob2RJRCxcclxuICAgICAgICBhcmdzXHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgcGx1Z2luLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gcGx1Z2luTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBwbHVnaW4uXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxyXG4gKiBAcGFyYW0gey4uLip9IGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cclxuICogQHJldHVybnMgeyp9IC0gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUGx1Z2luKHBsdWdpbk5hbWUsIG1ldGhvZE5hbWUsIC4uLmFyZ3MpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xyXG4gICAgICAgIHBhY2thZ2VOYW1lOiBcIndhaWxzLXBsdWdpbnNcIixcclxuICAgICAgICBzdHJ1Y3ROYW1lOiBwbHVnaW5OYW1lLFxyXG4gICAgICAgIG1ldGhvZE5hbWUsXHJcbiAgICAgICAgYXJnc1xyXG4gICAgfSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCwgJycpO1xyXG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcclxuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQW55IGlzIGEgZHVtbXkgY3JlYXRpb24gZnVuY3Rpb24gZm9yIHNpbXBsZSBvciB1bmtub3duIHR5cGVzLlxyXG4gKiBAdGVtcGxhdGUgVFxyXG4gKiBAcGFyYW0ge2FueX0gc291cmNlXHJcbiAqIEByZXR1cm5zIHtUfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEFueShzb3VyY2UpIHtcclxuICAgIHJldHVybiAvKiogQHR5cGUge1R9ICovKHNvdXJjZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXHJcbiAqIG51bGwgc3RyaW5ncyB3aXRoIGVtcHR5IHN0cmluZ3MuXHJcbiAqIEBwYXJhbSB7YW55fSBzb3VyY2VcclxuICogQHJldHVybnMge3N0cmluZ31cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeXRlU2xpY2Uoc291cmNlKSB7XHJcbiAgICByZXR1cm4gLyoqIEB0eXBlIHthbnl9ICovKChzb3VyY2UgPT0gbnVsbCkgPyBcIlwiIDogc291cmNlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEFycmF5IHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcnJheVxyXG4gKiB3aG9zZSBlbGVtZW50cyBhcmUgb2YgdGhhdCB0eXBlLlxyXG4gKiBAdGVtcGxhdGUgVFxyXG4gKiBAcGFyYW0geyhzb3VyY2U6IGFueSkgPT4gVH0gZWxlbWVudFxyXG4gKiBAcmV0dXJucyB7KHNvdXJjZTogYW55KSA9PiBUW119XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQXJyYXkoZWxlbWVudCkge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm4gW107XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XHJcbiAgICAgICAgICAgIHNvdXJjZVtpXSA9IGVsZW1lbnQoc291cmNlW2ldKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcclxuICogd2hvc2Uga2V5cyBhbmQgdmFsdWVzIGFyZSBvZiB0aG9zZSB0eXBlcy5cclxuICogQHRlbXBsYXRlIEssIFZcclxuICogQHBhcmFtIHsoc291cmNlOiBhbnkpID0+IEt9IGtleVxyXG4gKiBAcGFyYW0geyhzb3VyY2U6IGFueSkgPT4gVn0gdmFsdWVcclxuICogQHJldHVybnMgeyhzb3VyY2U6IGFueSkgPT4geyBbXzogS106IFYgfX1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBNYXAoa2V5LCB2YWx1ZSkge1xyXG4gICAgaWYgKHZhbHVlID09PSBBbnkpIHtcclxuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IHt9IDogc291cmNlKTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xyXG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcclxuICAgICAgICAgICAgcmV0dXJuIHt9O1xyXG4gICAgICAgIH1cclxuICAgICAgICBmb3IgKGNvbnN0IGtleSBpbiBzb3VyY2UpIHtcclxuICAgICAgICAgICAgc291cmNlW2tleV0gPSB2YWx1ZShzb3VyY2Vba2V5XSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBzb3VyY2U7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogTnVsbGFibGUgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcclxuICogYW5kIHJldHVybnMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBudWxsYWJsZSB2YWx1ZSBvZiB0aGF0IHR5cGUuXHJcbiAqIEB0ZW1wbGF0ZSBUXHJcbiAqIEBwYXJhbSB7KHNvdXJjZTogYW55KSA9PiBUfSBlbGVtZW50XHJcbiAqIEByZXR1cm5zIHsoc291cmNlOiBhbnkpID0+IChUIHwgbnVsbCl9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTnVsbGFibGUoZWxlbWVudCkge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiBBbnk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFN0cnVjdCB0YWtlcyBhbiBvYmplY3QgbWFwcGluZyBmaWVsZCBuYW1lcyB0byBjcmVhdGlvbiBmdW5jdGlvbnNcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgc3RydWN0LlxyXG4gKiBAdGVtcGxhdGUge3sgW186IHN0cmluZ106ICgoc291cmNlOiBhbnkpID0+IGFueSkgfX0gVFxyXG4gKiBAdGVtcGxhdGUge3sgW0tleSBpbiBrZXlvZiBUXT86IFJldHVyblR5cGU8VFtLZXldPiB9fSBVXHJcbiAqIEBwYXJhbSB7VH0gY3JlYXRlRmllbGRcclxuICogQHJldHVybnMgeyhzb3VyY2U6IGFueSkgPT4gVX1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQpIHtcclxuICAgIGxldCBhbGxBbnkgPSB0cnVlO1xyXG4gICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XHJcbiAgICAgICAgaWYgKGNyZWF0ZUZpZWxkW25hbWVdICE9PSBBbnkpIHtcclxuICAgICAgICAgICAgYWxsQW55ID0gZmFsc2U7XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuICAgIGlmIChhbGxBbnkpIHtcclxuICAgICAgICByZXR1cm4gQW55O1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XHJcbiAgICAgICAgICAgIGlmIChuYW1lIGluIHNvdXJjZSkge1xyXG4gICAgICAgICAgICAgICAgc291cmNlW25hbWVdID0gY3JlYXRlRmllbGRbbmFtZV0oc291cmNlW25hbWVdKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gc291cmNlO1xyXG4gICAgfTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNpemVcclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFdpZHRoIC0gVGhlIHdpZHRoLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gSGVpZ2h0IC0gVGhlIGhlaWdodC5cclxuICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge09iamVjdH0gUmVjdFxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gWCAtIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uXHJcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBXaWR0aCAtIFRoZSB3aWR0aCBvZiB0aGUgcmVjdGFuZ2xlLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gSGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTY3JlZW5cclxuICogQHByb3BlcnR5IHtzdHJpbmd9IElEIC0gVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBzY3JlZW4uXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBOYW1lIC0gSHVtYW4gcmVhZGFibGUgbmFtZSBvZiB0aGUgc2NyZWVuLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gU2NhbGVGYWN0b3IgLSBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLlxyXG4gKiBAcHJvcGVydHkge251bWJlcn0gWCAtIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uXHJcbiAqIEBwcm9wZXJ0eSB7U2l6ZX0gU2l6ZSAtIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uXHJcbiAqIEBwcm9wZXJ0eSB7UmVjdH0gQm91bmRzIC0gQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LlxyXG4gKiBAcHJvcGVydHkge1JlY3R9IFBoeXNpY2FsQm91bmRzIC0gQ29udGFpbnMgdGhlIHBoeXNpY2FsIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0IChiZWZvcmUgc2NhbGluZykuXHJcbiAqIEBwcm9wZXJ0eSB7UmVjdH0gV29ya0FyZWEgLSBDb250YWlucyB0aGUgYXJlYSBvZiB0aGUgc2NyZWVuIHRoYXQgaXMgYWN0dWFsbHkgdXNhYmxlIChleGNsdWRpbmcgdGFza2JhciBhbmQgb3RoZXIgc3lzdGVtIFVJKS5cclxuICogQHByb3BlcnR5IHtSZWN0fSBQaHlzaWNhbFdvcmtBcmVhIC0gQ29udGFpbnMgdGhlIHBoeXNpY2FsIFdvcmtBcmVhIG9mIHRoZSBzY3JlZW4gKGJlZm9yZSBzY2FsaW5nKS5cclxuICogQHByb3BlcnR5IHtib29sZWFufSBJc1ByaW1hcnkgLSBUcnVlIGlmIHRoaXMgaXMgdGhlIHByaW1hcnkgbW9uaXRvciBzZWxlY3RlZCBieSB0aGUgdXNlciBpbiB0aGUgb3BlcmF0aW5nIHN5c3RlbS5cclxuICogQHByb3BlcnR5IHtudW1iZXJ9IFJvdGF0aW9uIC0gVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uXHJcbiAqL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLlNjcmVlbnMsIFwiXCIpO1xyXG5cclxuY29uc3QgZ2V0QWxsID0gMDtcclxuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XHJcbmNvbnN0IGdldEN1cnJlbnQgPSAyO1xyXG5cclxuLyoqXHJcbiAqIEdldHMgYWxsIHNjcmVlbnMuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gYXJyYXkgb2YgU2NyZWVuIG9iamVjdHMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcclxufVxyXG4vKipcclxuICogR2V0cyB0aGUgcHJpbWFyeSBzY3JlZW4uXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBwcmltYXJ5IHNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XHJcbn1cclxuLyoqXHJcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cclxuICpcclxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcclxufVxyXG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBOzs7QUM2QkEsSUFBSSxjQUNBO0FBRUcsSUFBSSxTQUFTLENBQUMsT0FBTyxPQUFPO0FBQy9CLE1BQUksS0FBSztBQUVULE1BQUksSUFBSSxPQUFPO0FBQ2YsU0FBTyxLQUFLO0FBRVIsVUFBTSxZQUFhLEtBQUssT0FBTyxJQUFJLEtBQU0sQ0FBQztBQUFBLEVBQzlDO0FBQ0EsU0FBTztBQUNYOzs7QUM1QkEsSUFBTSxhQUFhLE9BQU8sU0FBUyxT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Sjs7O0FDNUhBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQTJDQSxJQUFNQyxRQUFPLHVCQUF1QixZQUFZLFNBQVMsRUFBRTtBQUUzRCxJQUFNLFNBQVM7QUFDZixJQUFNLGFBQWE7QUFDbkIsSUFBTSxhQUFhO0FBTVosU0FBUyxTQUFTO0FBQ3JCLFNBQU9BLE1BQUssTUFBTTtBQUN0QjtBQUtPLFNBQVMsYUFBYTtBQUN6QixTQUFPQSxNQUFLLFVBQVU7QUFDMUI7QUFNTyxTQUFTLGFBQWE7QUFDekIsU0FBT0EsTUFBSyxVQUFVO0FBQzFCOzs7QW5CM0RBLE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQWtDbEMsSUFBSSxjQUFjO0FBQ1gsU0FBUyxPQUFPO0FBQ25CLFNBQU8sT0FBTyxTQUFnQjtBQUM5QixFQUFPLE9BQU8scUJBQXFCO0FBQ25DLGdCQUFjO0FBQ2xCO0FBRUEsT0FBTyxpQkFBaUIsUUFBUSxNQUFNO0FBQ2xDLE1BQUksQ0FBQyxhQUFhO0FBQ2QsU0FBSztBQUFBLEVBQ1Q7QUFDSixDQUFDOyIsCiAgIm5hbWVzIjogWyJFcnJvciIsICJjYWxsIiwgIkVycm9yIiwgImNhbGwiLCAiZXZlbnROYW1lIiwgImNvbnRyb2xsZXIiLCAicmVzaXphYmxlIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgIkhpZGVNZXRob2QiLCAiU2hvd01ldGhvZCIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiY2FsbCIsICJNYXAiLCAiTWFwIiwgImtleSIsICJjYWxsIl0KfQo=
