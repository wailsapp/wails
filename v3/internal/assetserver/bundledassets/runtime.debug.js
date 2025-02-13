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
async function runtimeCallWithID(objectID, method, windowName, args) {
  let url = new URL(runtimeURL);
  if (objectID != null) {
    url.searchParams.append("object", objectID);
  }
  if (method != null) {
    url.searchParams.append("method", method);
  }
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
  let response = await fetch(url, fetchOptions);
  if (!response.ok) {
    throw new Error(await response.text());
  }
  if (response.headers.get("Content-Type") && response.headers.get("Content-Type").indexOf("application/json") !== -1) {
    return response.json();
  } else {
    return response.text();
  }
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
function dialog(type2, options = {}) {
  const id = generateID();
  options["dialog-id"] = id;
  return new Promise((resolve, reject) => {
    dialogResponses.set(id, { resolve, reject });
    call2(type2, options).catch((error) => {
      reject(error);
      dialogResponses.delete(id);
    });
  });
}
function dialogResultCallback(id, data, isJSON) {
  let p = dialogResponses.get(id);
  if (p) {
    dialogResponses.delete(id);
    if (isJSON) {
      p.resolve(JSON.parse(data));
    } else {
      p.resolve(data);
    }
  }
}
function dialogErrorCallback(id, message) {
  let p = dialogResponses.get(id);
  if (p) {
    dialogResponses.delete(id);
    p.reject(new Error2(message));
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
    APMPowerSettingChange: "windows:APMPowerSettingChange",
    APMPowerStatusChange: "windows:APMPowerStatusChange",
    APMResumeAutomatic: "windows:APMResumeAutomatic",
    APMResumeSuspend: "windows:APMResumeSuspend",
    APMSuspend: "windows:APMSuspend",
    ApplicationStarted: "windows:ApplicationStarted",
    SystemThemeChanged: "windows:SystemThemeChanged",
    WebViewNavigationCompleted: "windows:WebViewNavigationCompleted",
    WindowActive: "windows:WindowActive",
    WindowBackgroundErase: "windows:WindowBackgroundErase",
    WindowClickActive: "windows:WindowClickActive",
    WindowClosing: "windows:WindowClosing",
    WindowDidMove: "windows:WindowDidMove",
    WindowDidResize: "windows:WindowDidResize",
    WindowDPIChanged: "windows:WindowDPIChanged",
    WindowDragDrop: "windows:WindowDragDrop",
    WindowDragEnter: "windows:WindowDragEnter",
    WindowDragLeave: "windows:WindowDragLeave",
    WindowDragOver: "windows:WindowDragOver",
    WindowEndMove: "windows:WindowEndMove",
    WindowEndResize: "windows:WindowEndResize",
    WindowFullscreen: "windows:WindowFullscreen",
    WindowHide: "windows:WindowHide",
    WindowInactive: "windows:WindowInactive",
    WindowKeyDown: "windows:WindowKeyDown",
    WindowKeyUp: "windows:WindowKeyUp",
    WindowKillFocus: "windows:WindowKillFocus",
    WindowNonClientHit: "windows:WindowNonClientHit",
    WindowNonClientMouseDown: "windows:WindowNonClientMouseDown",
    WindowNonClientMouseLeave: "windows:WindowNonClientMouseLeave",
    WindowNonClientMouseMove: "windows:WindowNonClientMouseMove",
    WindowNonClientMouseUp: "windows:WindowNonClientMouseUp",
    WindowPaint: "windows:WindowPaint",
    WindowRestore: "windows:WindowRestore",
    WindowSetFocus: "windows:WindowSetFocus",
    WindowShow: "windows:WindowShow",
    WindowStartMove: "windows:WindowStartMove",
    WindowStartResize: "windows:WindowStartResize",
    WindowUnFullscreen: "windows:WindowUnFullscreen",
    WindowZOrderChanged: "windows:WindowZOrderChanged",
    WindowMinimise: "windows:WindowMinimise",
    WindowUnMinimise: "windows:WindowUnMinimise",
    WindowMaximise: "windows:WindowMaximise",
    WindowUnMaximise: "windows:WindowUnMaximise"
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
    ApplicationDidChangeTheme: "mac:ApplicationDidChangeTheme",
    ApplicationDidFinishLaunching: "mac:ApplicationDidFinishLaunching",
    ApplicationDidHide: "mac:ApplicationDidHide",
    ApplicationDidResignActive: "mac:ApplicationDidResignActive",
    ApplicationDidUnhide: "mac:ApplicationDidUnhide",
    ApplicationDidUpdate: "mac:ApplicationDidUpdate",
    ApplicationShouldHandleReopen: "mac:ApplicationShouldHandleReopen",
    ApplicationWillBecomeActive: "mac:ApplicationWillBecomeActive",
    ApplicationWillFinishLaunching: "mac:ApplicationWillFinishLaunching",
    ApplicationWillHide: "mac:ApplicationWillHide",
    ApplicationWillResignActive: "mac:ApplicationWillResignActive",
    ApplicationWillTerminate: "mac:ApplicationWillTerminate",
    ApplicationWillUnhide: "mac:ApplicationWillUnhide",
    ApplicationWillUpdate: "mac:ApplicationWillUpdate",
    MenuDidAddItem: "mac:MenuDidAddItem",
    MenuDidBeginTracking: "mac:MenuDidBeginTracking",
    MenuDidClose: "mac:MenuDidClose",
    MenuDidDisplayItem: "mac:MenuDidDisplayItem",
    MenuDidEndTracking: "mac:MenuDidEndTracking",
    MenuDidHighlightItem: "mac:MenuDidHighlightItem",
    MenuDidOpen: "mac:MenuDidOpen",
    MenuDidPopUp: "mac:MenuDidPopUp",
    MenuDidRemoveItem: "mac:MenuDidRemoveItem",
    MenuDidSendAction: "mac:MenuDidSendAction",
    MenuDidSendActionToItem: "mac:MenuDidSendActionToItem",
    MenuDidUpdate: "mac:MenuDidUpdate",
    MenuWillAddItem: "mac:MenuWillAddItem",
    MenuWillBeginTracking: "mac:MenuWillBeginTracking",
    MenuWillDisplayItem: "mac:MenuWillDisplayItem",
    MenuWillEndTracking: "mac:MenuWillEndTracking",
    MenuWillHighlightItem: "mac:MenuWillHighlightItem",
    MenuWillOpen: "mac:MenuWillOpen",
    MenuWillPopUp: "mac:MenuWillPopUp",
    MenuWillRemoveItem: "mac:MenuWillRemoveItem",
    MenuWillSendAction: "mac:MenuWillSendAction",
    MenuWillSendActionToItem: "mac:MenuWillSendActionToItem",
    MenuWillUpdate: "mac:MenuWillUpdate",
    WebViewDidCommitNavigation: "mac:WebViewDidCommitNavigation",
    WebViewDidFinishNavigation: "mac:WebViewDidFinishNavigation",
    WebViewDidReceiveServerRedirectForProvisionalNavigation: "mac:WebViewDidReceiveServerRedirectForProvisionalNavigation",
    WebViewDidStartProvisionalNavigation: "mac:WebViewDidStartProvisionalNavigation",
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
    WindowDidZoom: "mac:WindowDidZoom",
    WindowFileDraggingEntered: "mac:WindowFileDraggingEntered",
    WindowFileDraggingExited: "mac:WindowFileDraggingExited",
    WindowFileDraggingPerformed: "mac:WindowFileDraggingPerformed",
    WindowHide: "mac:WindowHide",
    WindowMaximise: "mac:WindowMaximise",
    WindowUnMaximise: "mac:WindowUnMaximise",
    WindowMinimise: "mac:WindowMinimise",
    WindowUnMinimise: "mac:WindowUnMinimise",
    WindowShouldClose: "mac:WindowShouldClose",
    WindowShow: "mac:WindowShow",
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
    WindowZoomIn: "mac:WindowZoomIn",
    WindowZoomOut: "mac:WindowZoomOut",
    WindowZoomReset: "mac:WindowZoomReset"
  },
  Linux: {
    ApplicationStartup: "linux:ApplicationStartup",
    SystemThemeChanged: "linux:SystemThemeChanged",
    WindowDeleteEvent: "linux:WindowDeleteEvent",
    WindowDidMove: "linux:WindowDidMove",
    WindowDidResize: "linux:WindowDidResize",
    WindowFocusIn: "linux:WindowFocusIn",
    WindowFocusOut: "linux:WindowFocusOut",
    WindowLoadChanged: "linux:WindowLoadChanged"
  },
  Common: {
    ApplicationOpenedWithFile: "common:ApplicationOpenedWithFile",
    ApplicationStarted: "common:ApplicationStarted",
    ThemeChanged: "common:ThemeChanged",
    WindowClosing: "common:WindowClosing",
    WindowDidMove: "common:WindowDidMove",
    WindowDidResize: "common:WindowDidResize",
    WindowDPIChanged: "common:WindowDPIChanged",
    WindowFilesDropped: "common:WindowFilesDropped",
    WindowFocus: "common:WindowFocus",
    WindowFullscreen: "common:WindowFullscreen",
    WindowHide: "common:WindowHide",
    WindowLostFocus: "common:WindowLostFocus",
    WindowMaximise: "common:WindowMaximise",
    WindowMinimise: "common:WindowMinimise",
    WindowRestore: "common:WindowRestore",
    WindowRuntimeReady: "common:WindowRuntimeReady",
    WindowShow: "common:WindowShow",
    WindowUnFullscreen: "common:WindowUnFullscreen",
    WindowUnMaximise: "common:WindowUnMaximise",
    WindowUnMinimise: "common:WindowUnMinimise",
    WindowZoom: "common:WindowZoom",
    WindowZoomIn: "common:WindowZoomIn",
    WindowZoomOut: "common:WindowZoomOut",
    WindowZoomReset: "common:WindowZoomReset"
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
  RuntimeError: () => RuntimeError
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
    if (!data) {
      promiseHandler.resolve();
    } else if (!isJSON) {
      promiseHandler.resolve(data);
    } else {
      try {
        promiseHandler.resolve(JSON.parse(data));
      } catch (err) {
        promiseHandler.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
      }
    }
  }
}
function errorHandler(id, data, isJSON) {
  const promiseHandler = getAndDeleteResponse(id);
  if (promiseHandler) {
    if (!isJSON) {
      promiseHandler.reject(new Error(data));
    } else {
      let error;
      try {
        error = JSON.parse(data);
      } catch (err) {
        promiseHandler.reject(new TypeError("could not parse error: " + err.message, { cause: err }));
        return;
      }
      let options = {};
      if (error.cause) {
        options.cause = error.cause;
      }
      let exception;
      switch (error.kind) {
        case "ReferenceError":
          exception = new ReferenceError(error.message, options);
          break;
        case "TypeError":
          exception = new TypeError(error.message, options);
          break;
        case "RuntimeError":
          exception = new RuntimeError(error.message, options);
          break;
        default:
          exception = new Error(error.message, options);
          break;
      }
      promiseHandler.reject(exception);
    }
  }
}
function getAndDeleteResponse(id) {
  const response = callResponses.get(id);
  callResponses.delete(id);
  return response;
}
var RuntimeError = class extends Error {
  /**
   * Constructs a new RuntimeError instance.
   *
   * @param {string} message - The error message.
   * @param {any[]} args - Optional arguments for the Error constructor.
   */
  constructor(message, ...args) {
    super(message, ...args);
    this.name = "RuntimeError";
  }
};
function Call(options) {
  const id = generateID2();
  const doCancel = () => {
    return cancelCall(type, { "call-id": id });
  };
  let queuedCancel = false, callRunning = false;
  let p = new Promise((resolve, reject) => {
    options["call-id"] = id;
    callResponses.set(id, { resolve, reject });
    call7(CallBinding, options).then((_) => {
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
function ByName(methodName, ...args) {
  return Call({
    methodName,
    args
  });
}
function ByID(methodID, ...args) {
  return Call({
    methodID,
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudF90eXBlcy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLmpzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUuanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMuanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG5cbmltcG9ydCBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCBcIi4vZHJhZ1wiO1xuXG4vLyBSZS1leHBvcnQgcHVibGljIEFQSVxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSBcIi4vYXBwbGljYXRpb25cIjtcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3NlclwiO1xuaW1wb3J0ICogYXMgQ2FsbCBmcm9tIFwiLi9jYWxsc1wiO1xuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZFwiO1xuaW1wb3J0ICogYXMgQ3JlYXRlIGZyb20gXCIuL2NyZWF0ZVwiO1xuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi9kaWFsb2dzXCI7XG5pbXBvcnQgKiBhcyBFdmVudHMgZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi9mbGFnc1wiO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi9zY3JlZW5zXCI7XG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4vc3lzdGVtXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvd1wiO1xuaW1wb3J0ICogYXMgV01MIGZyb20gXCIuL3dtbFwiO1xuXG5leHBvcnQge1xuICAgIEFwcGxpY2F0aW9uLFxuICAgIEJyb3dzZXIsXG4gICAgQ2FsbCxcbiAgICBDbGlwYm9hcmQsXG4gICAgQ3JlYXRlLFxuICAgIERpYWxvZ3MsXG4gICAgRXZlbnRzLFxuICAgIEZsYWdzLFxuICAgIFNjcmVlbnMsXG4gICAgU3lzdGVtLFxuICAgIFdpbmRvdyxcbiAgICBXTUxcbn07XG5cbmxldCBpbml0aWFsaXNlZCA9IGZhbHNlO1xuZXhwb3J0IGZ1bmN0aW9uIGluaXQoKSB7XG4gICAgd2luZG93Ll93YWlscy5pbnZva2UgPSBTeXN0ZW0uaW52b2tlO1xuICAgIFN5c3RlbS5pbnZva2UoXCJ3YWlsczpydW50aW1lOnJlYWR5XCIpO1xuICAgIGluaXRpYWxpc2VkID0gdHJ1ZTtcbn1cblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoXCJsb2FkXCIsICgpID0+IHtcbiAgICBpZiAoIWluaXRpYWxpc2VkKSB7XG4gICAgICAgIGluaXQoKTtcbiAgICB9XG59KTtcblxuLy8gTm90aWZ5IGJhY2tlbmRcblxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge09wZW5VUkx9IGZyb20gXCIuL2Jyb3dzZXJcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcbmltcG9ydCB7RW1pdCwgV2FpbHNFdmVudH0gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQge2NhbkFib3J0TGlzdGVuZXJzLCB3aGVuUmVhZHl9IGZyb20gXCIuL3V0aWxzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvd1wiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBzZW5kLlxuICogQHBhcmFtIHthbnl9IFtkYXRhPW51bGxdIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxuICpcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xuICAgIEVtaXQobmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBzcGVjaWZpZWQgd2luZG93LlxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGNhbGwgdGhlIG1ldGhvZCBvbi5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICovXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKHdpbmRvd05hbWUsIG1ldGhvZE5hbWUpIHtcbiAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBXaW5kb3cuR2V0KHdpbmRvd05hbWUpO1xuICAgIGNvbnN0IG1ldGhvZCA9IHRhcmdldFdpbmRvd1ttZXRob2ROYW1lXTtcblxuICAgIGlmICh0eXBlb2YgbWV0aG9kICE9PSBcImZ1bmN0aW9uXCIpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICB0cnkge1xuICAgICAgICBtZXRob2QuY2FsbCh0YXJnZXRXaW5kb3cpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JzogYCwgZSk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXRcbiAqXG4gKiBAcGFyYW0ge0V2ZW50fSBldlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gb25XTUxUcmlnZ2VyZWQoZXYpIHtcbiAgICBjb25zdCBlbGVtZW50ID0gZXYuY3VycmVudFRhcmdldDtcblxuICAgIGZ1bmN0aW9uIHJ1bkVmZmVjdChjaG9pY2UgPSBcIlllc1wiKSB7XG4gICAgICAgIGlmIChjaG9pY2UgIT09IFwiWWVzXCIpXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWV2ZW50Jyk7XG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10YXJnZXQtd2luZG93JykgfHwgXCJcIjtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCB1cmwgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtb3BlblVSTCcpO1xuXG4gICAgICAgIGlmIChldmVudFR5cGUgIT09IG51bGwpXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICBpZiAodXJsICE9PSBudWxsKVxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgfVxuXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG5cbiAgICBpZiAoY29uZmlybSkge1xuICAgICAgICBRdWVzdGlvbih7XG4gICAgICAgICAgICBUaXRsZTogXCJDb25maXJtXCIsXG4gICAgICAgICAgICBNZXNzYWdlOiBjb25maXJtLFxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxuICAgICAgICAgICAgQnV0dG9uczogW1xuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiWWVzXCIgfSxcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XG4gICAgICAgICAgICBdXG4gICAgICAgIH0pLnRoZW4ocnVuRWZmZWN0KTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBydW5FZmZlY3QoKTtcbiAgICB9XG59XG5cbi8qKlxuICogQHR5cGUge3N5bWJvbH1cbiAqL1xuY29uc3QgY29udHJvbGxlciA9IFN5bWJvbCgpO1xuXG4vKipcbiAqIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IGRvZXMgbm90IGFjdHVhbGx5IHJlbWVtYmVyIGFjdGl2ZSBldmVudCBsaXN0ZW5lcnM6IGluc3RlYWRcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXG4gKi9cbmNsYXNzIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IHtcbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgLyoqXG4gICAgICAgICAqIFN0b3JlcyB0aGUgQWJvcnRDb250cm9sbGVyIHRoYXQgY2FuIGJlIHVzZWQgdG8gcmVtb3ZlIGFsbCBjdXJyZW50bHkgYWN0aXZlIGxpc3RlbmVycy5cbiAgICAgICAgICpcbiAgICAgICAgICogQHByaXZhdGVcbiAgICAgICAgICogQG5hbWUge0BsaW5rIGNvbnRyb2xsZXJ9XG4gICAgICAgICAqIEBtZW1iZXIge0Fib3J0Q29udHJvbGxlcn1cbiAgICAgICAgICovXG4gICAgICAgIHRoaXNbY29udHJvbGxlcl0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyBhbiBvcHRpb25zIG9iamVjdCBmb3IgYWRkRXZlbnRMaXN0ZW5lciB0aGF0IHRpZXMgdGhlIGxpc3RlbmVyXG4gICAgICogdG8gdGhlIEFib3J0U2lnbmFsIGZyb20gdGhlIGN1cnJlbnQgQWJvcnRDb250cm9sbGVyLlxuICAgICAqXG4gICAgICogQHBhcmFtIHtIVE1MRWxlbWVudH0gZWxlbWVudCBBbiBIVE1MIGVsZW1lbnRcbiAgICAgKiBAcGFyYW0ge3N0cmluZ1tdfSB0cmlnZ2VycyBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRzXG4gICAgICogQHJldHVybnMge0FkZEV2ZW50TGlzdGVuZXJPcHRpb25zfVxuICAgICAqL1xuICAgIHNldChlbGVtZW50LCB0cmlnZ2Vycykge1xuICAgICAgICByZXR1cm4geyBzaWduYWw6IHRoaXNbY29udHJvbGxlcl0uc2lnbmFsIH07XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyB7dm9pZH1cbiAgICAgKi9cbiAgICByZXNldCgpIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyXS5hYm9ydCgpO1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJdID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBAdHlwZSB7c3ltYm9sfVxuICovXG5jb25zdCB0cmlnZ2VyTWFwID0gU3ltYm9sKCk7XG5cbi8qKlxuICogQHR5cGUge3N5bWJvbH1cbiAqL1xuY29uc3QgZWxlbWVudENvdW50ID0gU3ltYm9sKCk7XG5cbi8qKlxuICogV2Vha01hcFJlZ2lzdHJ5IG1hcHMgYWN0aXZlIHRyaWdnZXIgZXZlbnRzIHRvIGVhY2ggRE9NIGVsZW1lbnQgdGhyb3VnaCBhIFdlYWtNYXAuXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcbiAqIGNvbGxlY3Rpb24gb2YgdGhlIGludm9sdmVkIGVsZW1lbnRzLlxuICovXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICAvKipcbiAgICAgICAgICogU3RvcmVzIHRoZSBjdXJyZW50IGVsZW1lbnQtdG8tdHJpZ2dlciBtYXBwaW5nLlxuICAgICAgICAgKlxuICAgICAgICAgKiBAcHJpdmF0ZVxuICAgICAgICAgKiBAbmFtZSB7QGxpbmsgdHJpZ2dlck1hcH1cbiAgICAgICAgICogQG1lbWJlciB7V2Vha01hcDxIVE1MRWxlbWVudCwgc3RyaW5nW10+fVxuICAgICAgICAgKi9cbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwXSA9IG5ldyBXZWFrTWFwKCk7XG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIENvdW50cyB0aGUgbnVtYmVyIG9mIGVsZW1lbnRzIHdpdGggYWN0aXZlIFdNTCB0cmlnZ2Vycy5cbiAgICAgICAgICpcbiAgICAgICAgICogQHByaXZhdGVcbiAgICAgICAgICogQG5hbWUge0BsaW5rIGVsZW1lbnRDb3VudH1cbiAgICAgICAgICogQG1lbWJlciB7bnVtYmVyfVxuICAgICAgICAgKi9cbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRdID0gMDtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB7SFRNTEVsZW1lbnR9IGVsZW1lbnQgQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHtzdHJpbmdbXX0gdHJpZ2dlcnMgVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XG4gICAgICogQHJldHVybnMge0FkZEV2ZW50TGlzdGVuZXJPcHRpb25zfVxuICAgICAqL1xuICAgIHNldChlbGVtZW50LCB0cmlnZ2Vycykge1xuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudF0gKz0gIXRoaXNbdHJpZ2dlck1hcF0uaGFzKGVsZW1lbnQpO1xuICAgICAgICB0aGlzW3RyaWdnZXJNYXBdLnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgICAgIHJldHVybiB7fTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIHt2b2lkfVxuICAgICAqL1xuICAgIHJlc2V0KCkge1xuICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRdIDw9IDApXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnKicpKSB7XG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRdIDw9IDApXG4gICAgICAgICAgICAgICAgYnJlYWs7XG5cbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwXS5nZXQoZWxlbWVudCk7XG4gICAgICAgICAgICB0aGlzW2VsZW1lbnRDb3VudF0gLT0gKHR5cGVvZiB0cmlnZ2VycyAhPT0gXCJ1bmRlZmluZWRcIik7XG5cbiAgICAgICAgICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycyB8fCBbXSlcbiAgICAgICAgICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQpO1xuICAgICAgICB9XG5cbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwXSA9IG5ldyBXZWFrTWFwKCk7XG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50XSA9IDA7XG4gICAgfVxufVxuXG5jb25zdCB0cmlnZ2VyUmVnaXN0cnkgPSBjYW5BYm9ydExpc3RlbmVycygpID8gbmV3IEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5KCkgOiBuZXcgV2Vha01hcFJlZ2lzdHJ5KCk7XG5cbi8qKlxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gdGhlIHNwZWNpZmllZCBlbGVtZW50LlxuICpcbiAqIEBwYXJhbSB7SFRNTEVsZW1lbnR9IGVsZW1lbnRcbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50KSB7XG4gICAgY29uc3QgdHJpZ2dlclJlZ0V4cCA9IC9cXFMrL2c7XG4gICAgY29uc3QgdHJpZ2dlckF0dHIgPSAoZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCIpO1xuICAgIGNvbnN0IHRyaWdnZXJzID0gW107XG5cbiAgICBsZXQgbWF0Y2g7XG4gICAgd2hpbGUgKChtYXRjaCA9IHRyaWdnZXJSZWdFeHAuZXhlYyh0cmlnZ2VyQXR0cikpICE9PSBudWxsKVxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcblxuICAgIGNvbnN0IG9wdGlvbnMgPSB0cmlnZ2VyUmVnaXN0cnkuc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2NoZWR1bGVzIGFuIGF1dG9tYXRpYyByZWxvYWQgb2YgV01MIHRvIGJlIHBlcmZvcm1lZCBhcyBzb29uIGFzIHRoZSBkb2N1bWVudCBpcyBmdWxseSBsb2FkZWQuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVuYWJsZSgpIHtcbiAgICB3aGVuUmVhZHkoUmVsb2FkKTtcbn1cblxuLyoqXG4gKiBSZWxvYWRzIHRoZSBXTUwgcGFnZSBieSBhZGRpbmcgbmVjZXNzYXJ5IGV2ZW50IGxpc3RlbmVycyBhbmQgYnJvd3NlciBsaXN0ZW5lcnMuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpIHtcbiAgICB0cmlnZ2VyUmVnaXN0cnkucmVzZXQoKTtcbiAgICBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC1ldmVudF0sIFtkYXRhLXdtbC13aW5kb3ddLCBbZGF0YS13bWwtb3BlblVSTF0nKS5mb3JFYWNoKGFkZFdNTExpc3RlbmVycyk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQnJvd3NlciwgJycpO1xuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xuXG4vKipcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIC0gVGhlIFVSTCB0byBvcGVuXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmwpIHtcbiAgICByZXR1cm4gY2FsbChCcm93c2VyT3BlblVSTCwge3VybH0pO1xufVxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxuXG4vLyBUaGUgTUlUIExpY2Vuc2UgKE1JVClcbi8vXG4vLyBDb3B5cmlnaHQgMjAxNyBBbmRyZXkgU2l0bmlrIDxhbmRyZXlAc2l0bmlrLnJ1PlxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2Zcbi8vIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWwgaW5cbi8vIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMgdG9cbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXG4vLyB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sXG4vLyAgICAgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6XG4vL1xuLy8gICAgIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cbi8vXG4vLyAgICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksIEZJVE5FU1Ncbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxuLy8gQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLCBXSEVUSEVSXG4vLyBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sIE9VVCBPRiBPUiBJTlxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cblxuLy8gVGhpcyBhbHBoYWJldCB1c2VzIGBBLVphLXowLTlfLWAgc3ltYm9scy5cbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cbi8vIFJlZmVyZW5jZXMgdG8gdGhlIHNhbWUgZmlsZSAod29ya3MgYm90aCBmb3IgZ3ppcCBhbmQgYnJvdGxpKTpcbi8vIGAndXNlYCwgYGFuZG9tYCwgYW5kIGByaWN0J2Bcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XG4vLyBgLTI2VGAsIGAxOTgzYCwgYDQwcHhgLCBgNzVweGAsIGBidXNoYCwgYGphY2tgLCBgbWluZGAsIGB2ZXJ5YCwgYW5kIGB3b2xmYFxubGV0IHVybEFscGhhYmV0ID1cbiAgICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcblxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICAvLyBBIGNvbXBhY3QgYWx0ZXJuYXRpdmUgZm9yIGBmb3IgKHZhciBpID0gMDsgaSA8IHN0ZXA7IGkrKylgLlxuICAgIGxldCBpID0gc2l6ZSB8IDBcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICAgIC8vIGB8IDBgIGlzIG1vcmUgY29tcGFjdCBhbmQgZmFzdGVyIHRoYW4gYE1hdGguZmxvb3IoKWAuXG4gICAgICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG59IiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxuICAgIENhbmNlbENhbGw6IDEwLFxufVxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG4vKipcbiAqIENyZWF0ZXMgYSBydW50aW1lIGNhbGxlciBmdW5jdGlvbiB0aGF0IGludm9rZXMgYSBzcGVjaWZpZWQgbWV0aG9kIG9uIGEgZ2l2ZW4gb2JqZWN0IHdpdGhpbiBhIHNwZWNpZmllZCB3aW5kb3cgY29udGV4dC5cbiAqXG4gKiBAcGFyYW0ge09iamVjdH0gb2JqZWN0IC0gVGhlIG9iamVjdCBvbiB3aGljaCB0aGUgbWV0aG9kIGlzIHRvIGJlIGludm9rZWQuXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgY29udGV4dCBpbiB3aGljaCB0aGUgbWV0aG9kIHNob3VsZCBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyB7RnVuY3Rpb259IEEgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24gdGhhdCB0YWtlcyB0aGUgbWV0aG9kIG5hbWUgYW5kIG9wdGlvbmFsbHkgYXJndW1lbnRzIGFuZCBpbnZva2VzIHRoZSBtZXRob2Qgd2l0aGluIHRoZSBzcGVjaWZpZWQgd2luZG93IGNvbnRleHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG4vKipcbiAqIENyZWF0ZXMgYSBuZXcgcnVudGltZSBjYWxsZXIgd2l0aCBzcGVjaWZpZWQgSUQuXG4gKlxuICogQHBhcmFtIHtudW1iZXJ9IG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0ge3N0cmluZ30gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gKiBAcmV0dXJuIHtGdW5jdGlvbn0gLSBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3QsIHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQobnVsbCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbn1cblxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncykge1xuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xuICAgIGlmIChvYmplY3RJRCAhPSBudWxsKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwib2JqZWN0XCIsIG9iamVjdElEKTtcbiAgICB9XG4gICAgaWYgKG1ldGhvZCAhPSBudWxsKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwibWV0aG9kXCIsIG1ldGhvZCk7XG4gICAgfVxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XG4gICAgICAgIGhlYWRlcnM6IHt9LFxuICAgIH07XG4gICAgaWYgKHdpbmRvd05hbWUpIHtcbiAgICAgICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLXdpbmRvdy1uYW1lXCJdID0gd2luZG93TmFtZTtcbiAgICB9XG4gICAgaWYgKGFyZ3MpIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcbiAgICB9XG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xuXG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2godXJsLCBmZXRjaE9wdGlvbnMpO1xuICAgIGlmICghcmVzcG9uc2Uub2spIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XG4gICAgfVxuXG4gICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gT3BlbkZpbGVEaWFsb2dPcHRpb25zXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DaG9vc2VEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZUZpbGVzXSAtIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ3JlYXRlRGlyZWN0b3JpZXNdIC0gSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbU2hvd0hpZGRlbkZpbGVzXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbUmVzb2x2ZXNBbGlhc2VzXSAtIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c011bHRpcGxlU2VsZWN0aW9uXSAtIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0hpZGVFeHRlbnNpb25dIC0gSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhblNlbGVjdEhpZGRlbkV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0FsbG93c090aGVyRmlsZXR5cGVzXSAtIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtGaWxlRmlsdGVyW119IFtGaWx0ZXJzXSAtIEFycmF5IG9mIGZpbGUgZmlsdGVycy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbVGl0bGVdIC0gVGl0bGUgb2YgdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTWVzc2FnZV0gLSBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbQnV0dG9uVGV4dF0gLSBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlyZWN0b3J5XSAtIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LlxuICovXG5cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTYXZlRmlsZURpYWxvZ09wdGlvbnNcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRmlsZW5hbWVdIC0gRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0NhbkNob29zZURpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuQ2hvb3NlRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtDYW5DcmVhdGVEaXJlY3Rvcmllc10gLSBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtTaG93SGlkZGVuRmlsZXNdIC0gSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtSZXNvbHZlc0FsaWFzZXNdIC0gSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb25dIC0gSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSGlkZUV4dGVuc2lvbl0gLSBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uXSAtIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzXSAtIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbQWxsb3dzT3RoZXJGaWxldHlwZXNdIC0gSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuXG4gKiBAcHJvcGVydHkge0ZpbGVGaWx0ZXJbXX0gW0ZpbHRlcnNdIC0gQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaXRsZSBvZiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtNZXNzYWdlXSAtIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtCdXR0b25UZXh0XSAtIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFtEaXJlY3RvcnldIC0gRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gW0RldGFjaGVkXSAtIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xuICogQHByb3BlcnR5IHtzdHJpbmd9IFtUaXRsZV0gLSBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW01lc3NhZ2VdIC0gVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge0J1dHRvbltdfSBbQnV0dG9uc10gLSBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtEZXRhY2hlZF0gLSBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBCdXR0b25cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbTGFiZWxdIC0gVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IFtJc0NhbmNlbF0gLSBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLlxuICogQHByb3BlcnR5IHtib29sZWFufSBbSXNEZWZhdWx0XSAtIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IEZpbGVGaWx0ZXJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBbRGlzcGxheU5hbWVdIC0gRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gW1BhdHRlcm5dIC0gUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy5cbiAqL1xuXG4vLyBzZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpYWxvZ0Vycm9yQ2FsbGJhY2sgPSBkaWFsb2dFcnJvckNhbGxiYWNrO1xud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJy4vbmFub2lkLmpzJztcblxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZywgJycpO1xuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBpZCB0aGF0IGlzIG5vdCBwcmVzZW50IGluIGRpYWxvZ1Jlc3BvbnNlcy5cbiAqIEByZXR1cm5zIHtzdHJpbmd9IHVuaXF1ZSBpZFxuICovXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogU2hvd3MgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7bnVtYmVyfSB0eXBlIC0gdHlwZSBvZiBkaWFsb2dcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN8T3BlbkZpbGVEaWFsb2dPcHRpb25zfFNhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIG9wdGlvbnMgZm9yIHRoZSBkaWFsb2dcbiAqIEByZXR1cm5zIHtQcm9taXNlfSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zID0ge30pIHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBjYWxsYmFjayBmcm9tIGEgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBpZCAtIFRoZSBJRCBvZiB0aGUgZGlhbG9nIHJlc3BvbnNlLlxuICogQHBhcmFtIHtzdHJpbmd9IGRhdGEgLSBUaGUgZGF0YSByZWNlaXZlZCBmcm9tIHRoZSBkaWFsb2cuXG4gKiBAcGFyYW0ge2Jvb2xlYW59IGlzSlNPTiAtIEZsYWcgaW5kaWNhdGluZyB3aGV0aGVyIHRoZSBkYXRhIGlzIGluIEpTT04gZm9ybWF0LlxuICpcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cbiAqL1xuZnVuY3Rpb24gZGlhbG9nUmVzdWx0Q2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBDYWxsYmFjayBmdW5jdGlvbiBmb3IgaGFuZGxpbmcgZXJyb3JzIGluIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHAucmVqZWN0KG5ldyBFcnJvcihtZXNzYWdlKSk7XG4gICAgfVxufVxuXG5cbi8vIFJlcGxhY2UgYG1ldGhvZHNgIHdpdGggY29uc3RhbnRzIGluIFRpdGxlIENhc2VcblxuLyoqXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgY29uc3QgSW5mbyA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7XG5cbi8qKlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGNvbnN0IFdhcm5pbmcgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpO1xuXG4vKipcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBjb25zdCBFcnJvciA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpO1xuXG4vKipcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBjb25zdCBRdWVzdGlvbiA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpO1xuXG4vKipcbiAqIEBwYXJhbSB7T3BlbkZpbGVEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZ1tdfHN0cmluZz59IFJldHVybnMgc2VsZWN0ZWQgZmlsZSBvciBsaXN0IG9mIGZpbGVzLiBSZXR1cm5zIGJsYW5rIHN0cmluZyBpZiBubyBmaWxlIGlzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgY29uc3QgT3BlbkZpbGUgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKTtcblxuLyoqXG4gKiBAcGFyYW0ge1NhdmVGaWxlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlLiBSZXR1cm5zIGJsYW5rIHN0cmluZyBpZiBubyBmaWxlIGlzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgY29uc3QgU2F2ZUZpbGUgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcbiAqL1xuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5pbXBvcnQge0V2ZW50VHlwZXN9IGZyb20gXCIuL2V2ZW50X3R5cGVzXCI7XG5leHBvcnQgY29uc3QgVHlwZXMgPSBFdmVudFR5cGVzO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuRXZlbnRzLCAnJyk7XG5jb25zdCBFbWl0TWV0aG9kID0gMDtcbmNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgY29uc3RydWN0b3IobmFtZSwgZGF0YSA9IG51bGwpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzZXR1cCgpIHtcbn1cblxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50KSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIGxldCB0b1JlbW92ZSA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHJldHVybiB0cnVlO1xuICAgICAgICB9KTtcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xuICAgICAgICAgICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3MgLSBUaGUgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdGhlIGNhbGxiYWNrIGNhbiBiZSBjYWxsZWQgZm9yIHRoZSBldmVudC4gT25jZSB0aGUgbWF4aW11bSBudW1iZXIgaXMgcmVhY2hlZCwgdGhlIGNhbGxiYWNrIHdpbGwgbm8gbG9uZ2VyIGJlIGNhbGxlZC5cbiAqXG4gQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCB3aGVuIHRoZSBzcGVjaWZpZWQgZXZlbnQgb2NjdXJzLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXG4gKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZC4gSXQgdGFrZXMgbm8gcGFyYW1ldGVycy5cbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuICovXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lLCBjYWxsYmFjaykgeyByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7IH1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cbiAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrIC0gVGhlIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIGV2ZW50IG9jY3Vycy5cbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHsgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7IH1cblxuLyoqXG4gKiBSZW1vdmVzIHRoZSBzcGVjaWZpZWQgbGlzdGVuZXIgZnJvbSB0aGUgZXZlbnQgbGlzdGVuZXJzIGNvbGxlY3Rpb24uXG4gKiBJZiBhbGwgbGlzdGVuZXJzIGZvciB0aGUgZXZlbnQgYXJlIHJlbW92ZWQsIHRoZSBldmVudCBrZXkgaXMgZGVsZXRlZCBmcm9tIHRoZSBjb2xsZWN0aW9uLlxuICpcbiAqIEBwYXJhbSB7T2JqZWN0fSBsaXN0ZW5lciAtIFRoZSBsaXN0ZW5lciB0byBiZSByZW1vdmVkLlxuICovXG5mdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcikge1xuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG59XG5cblxuLyoqXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxuICogQHBhcmFtIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzIC0gQWRkaXRpb25hbCBldmVudCBuYW1lcyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcbn1cbi8qKlxuICogUmVtb3ZlcyBhbGwgZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBmdW5jdGlvbiBPZmZBbGxcbiAqIEByZXR1cm5zIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkgeyBldmVudExpc3RlbmVycy5jbGVhcigpOyB9XG5cbi8qKlxuICogRW1pdHMgYW4gZXZlbnQgdXNpbmcgdGhlIGdpdmVuIGV2ZW50IG5hbWUuXG4gKlxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxuICogQHJldHVybnMge2FueX0gLSBUaGUgcmVzdWx0IG9mIHRoZSBlbWl0dGVkIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkgeyByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCBldmVudCk7IH1cbiIsICJcbmV4cG9ydCBjb25zdCBFdmVudFR5cGVzID0ge1xuXHRXaW5kb3dzOiB7XG5cdFx0QVBNUG93ZXJTZXR0aW5nQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTZXR0aW5nQ2hhbmdlXCIsXG5cdFx0QVBNUG93ZXJTdGF0dXNDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclN0YXR1c0NoYW5nZVwiLFxuXHRcdEFQTVJlc3VtZUF1dG9tYXRpYzogXCJ3aW5kb3dzOkFQTVJlc3VtZUF1dG9tYXRpY1wiLFxuXHRcdEFQTVJlc3VtZVN1c3BlbmQ6IFwid2luZG93czpBUE1SZXN1bWVTdXNwZW5kXCIsXG5cdFx0QVBNU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVN1c3BlbmRcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwid2luZG93czpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwid2luZG93czpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZDogXCJ3aW5kb3dzOldlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkXCIsXG5cdFx0V2luZG93QWN0aXZlOiBcIndpbmRvd3M6V2luZG93QWN0aXZlXCIsXG5cdFx0V2luZG93QmFja2dyb3VuZEVyYXNlOiBcIndpbmRvd3M6V2luZG93QmFja2dyb3VuZEVyYXNlXCIsXG5cdFx0V2luZG93Q2xpY2tBY3RpdmU6IFwid2luZG93czpXaW5kb3dDbGlja0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwid2luZG93czpXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJ3aW5kb3dzOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwid2luZG93czpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0RyYWdEcm9wOiBcIndpbmRvd3M6V2luZG93RHJhZ0Ryb3BcIixcblx0XHRXaW5kb3dEcmFnRW50ZXI6IFwid2luZG93czpXaW5kb3dEcmFnRW50ZXJcIixcblx0XHRXaW5kb3dEcmFnTGVhdmU6IFwid2luZG93czpXaW5kb3dEcmFnTGVhdmVcIixcblx0XHRXaW5kb3dEcmFnT3ZlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdPdmVyXCIsXG5cdFx0V2luZG93RW5kTW92ZTogXCJ3aW5kb3dzOldpbmRvd0VuZE1vdmVcIixcblx0XHRXaW5kb3dFbmRSZXNpemU6IFwid2luZG93czpXaW5kb3dFbmRSZXNpemVcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd0hpZGU6IFwid2luZG93czpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93SW5hY3RpdmU6IFwid2luZG93czpXaW5kb3dJbmFjdGl2ZVwiLFxuXHRcdFdpbmRvd0tleURvd246IFwid2luZG93czpXaW5kb3dLZXlEb3duXCIsXG5cdFx0V2luZG93S2V5VXA6IFwid2luZG93czpXaW5kb3dLZXlVcFwiLFxuXHRcdFdpbmRvd0tpbGxGb2N1czogXCJ3aW5kb3dzOldpbmRvd0tpbGxGb2N1c1wiLFxuXHRcdFdpbmRvd05vbkNsaWVudEhpdDogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudEhpdFwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlRG93bjogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlRG93blwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlTGVhdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZUxlYXZlXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VNb3ZlOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VNb3ZlXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VVcDogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlVXBcIixcblx0XHRXaW5kb3dQYWludDogXCJ3aW5kb3dzOldpbmRvd1BhaW50XCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJ3aW5kb3dzOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dTZXRGb2N1czogXCJ3aW5kb3dzOldpbmRvd1NldEZvY3VzXCIsXG5cdFx0V2luZG93U2hvdzogXCJ3aW5kb3dzOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dTdGFydE1vdmU6IFwid2luZG93czpXaW5kb3dTdGFydE1vdmVcIixcblx0XHRXaW5kb3dTdGFydFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0UmVzaXplXCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93Wk9yZGVyQ2hhbmdlZDogXCJ3aW5kb3dzOldpbmRvd1pPcmRlckNoYW5nZWRcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJ3aW5kb3dzOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWF4aW1pc2VcIixcblx0fSxcblx0TWFjOiB7XG5cdFx0QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUljb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0JhckZyYW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0JhckZyYW1lXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlVGhlbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlVGhlbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmdcIixcblx0XHRBcHBsaWNhdGlvbkRpZEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkVXBkYXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFVwZGF0ZVwiLFxuXHRcdEFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuOiBcIm1hYzpBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlblwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmdcIixcblx0XHRBcHBsaWNhdGlvbldpbGxIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVXBkYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVcGRhdGVcIixcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcblx0XHRNZW51RGlkQmVnaW5UcmFja2luZzogXCJtYWM6TWVudURpZEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51RGlkQ2xvc2U6IFwibWFjOk1lbnVEaWRDbG9zZVwiLFxuXHRcdE1lbnVEaWREaXNwbGF5SXRlbTogXCJtYWM6TWVudURpZERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudURpZEVuZFRyYWNraW5nOiBcIm1hYzpNZW51RGlkRW5kVHJhY2tpbmdcIixcblx0XHRNZW51RGlkSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudURpZEhpZ2hsaWdodEl0ZW1cIixcblx0XHRNZW51RGlkT3BlbjogXCJtYWM6TWVudURpZE9wZW5cIixcblx0XHRNZW51RGlkUG9wVXA6IFwibWFjOk1lbnVEaWRQb3BVcFwiLFxuXHRcdE1lbnVEaWRSZW1vdmVJdGVtOiBcIm1hYzpNZW51RGlkUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblRvSXRlbVwiLFxuXHRcdE1lbnVEaWRVcGRhdGU6IFwibWFjOk1lbnVEaWRVcGRhdGVcIixcblx0XHRNZW51V2lsbEFkZEl0ZW06IFwibWFjOk1lbnVXaWxsQWRkSXRlbVwiLFxuXHRcdE1lbnVXaWxsQmVnaW5UcmFja2luZzogXCJtYWM6TWVudVdpbGxCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVXaWxsRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsRW5kVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVXaWxsSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVXaWxsT3BlbjogXCJtYWM6TWVudVdpbGxPcGVuXCIsXG5cdFx0TWVudVdpbGxQb3BVcDogXCJtYWM6TWVudVdpbGxQb3BVcFwiLFxuXHRcdE1lbnVXaWxsUmVtb3ZlSXRlbTogXCJtYWM6TWVudVdpbGxSZW1vdmVJdGVtXCIsXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbVwiLFxuXHRcdE1lbnVXaWxsVXBkYXRlOiBcIm1hYzpNZW51V2lsbFVwZGF0ZVwiLFxuXHRcdFdlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRSZWNlaXZlU2VydmVyUmVkaXJlY3RGb3JQcm92aXNpb25hbE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRSZWNlaXZlU2VydmVyUmVkaXJlY3RGb3JQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdpbmRvd0RpZEJlY29tZUtleTogXCJtYWM6V2luZG93RGlkQmVjb21lS2V5XCIsXG5cdFx0V2luZG93RGlkQmVjb21lTWFpbjogXCJtYWM6V2luZG93RGlkQmVjb21lTWFpblwiLFxuXHRcdFdpbmRvd0RpZEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd0RpZEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VBbHBoYTogXCJtYWM6V2luZG93RGlkQ2hhbmdlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nTG9jYXRpb246IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVGl0bGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkQ2hhbmdlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRFbmRTaGVldDogXCJtYWM6V2luZG93RGlkRW5kU2hlZXRcIixcblx0XHRXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4cG9zZTogXCJtYWM6V2luZG93RGlkRXhwb3NlXCIsXG5cdFx0V2luZG93RGlkRm9jdXM6IFwibWFjOldpbmRvd0RpZEZvY3VzXCIsXG5cdFx0V2luZG93RGlkTWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZE1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJtYWM6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT25TY3JlZW5cIixcblx0XHRXaW5kb3dEaWRSZXNpZ25LZXk6IFwibWFjOldpbmRvd0RpZFJlc2lnbktleVwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbk1haW46IFwibWFjOldpbmRvd0RpZFJlc2lnbk1haW5cIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibWFjOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZTogXCJtYWM6V2luZG93RGlkVXBkYXRlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUFscGhhXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVNoYWRvd1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWRab29tOiBcIm1hYzpXaW5kb3dEaWRab29tXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRW50ZXJlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRXhpdGVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWRcIixcblx0XHRXaW5kb3dIaWRlOiBcIm1hYzpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwibWFjOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJtYWM6V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIm1hYzpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwibWFjOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dTaG91bGRDbG9zZTogXCJtYWM6V2luZG93U2hvdWxkQ2xvc2VcIixcblx0XHRXaW5kb3dTaG93OiBcIm1hYzpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZUtleTogXCJtYWM6V2luZG93V2lsbEJlY29tZUtleVwiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dXaWxsQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd1dpbGxDaGFuZ2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd1dpbGxDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dXaWxsQ2xvc2U6IFwibWFjOldpbmRvd1dpbGxDbG9zZVwiLFxuXHRcdFdpbmRvd1dpbGxEZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlckZ1bGxTY3JlZW46IFwibWFjOldpbmRvd1dpbGxFbnRlckZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd1dpbGxFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0VmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFeGl0VmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dXaWxsRm9jdXM6IFwibWFjOldpbmRvd1dpbGxGb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbE1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbE1vdmU6IFwibWFjOldpbmRvd1dpbGxNb3ZlXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT25TY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsUmVzaWduTWFpbjogXCJtYWM6V2luZG93V2lsbFJlc2lnbk1haW5cIixcblx0XHRXaW5kb3dXaWxsUmVzaXplOiBcIm1hYzpXaW5kb3dXaWxsUmVzaXplXCIsXG5cdFx0V2luZG93V2lsbFVuZm9jdXM6IFwibWFjOldpbmRvd1dpbGxVbmZvY3VzXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZUFscGhhXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVNoYWRvd1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUb29sYmFyXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHk6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5XCIsXG5cdFx0V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWU6IFwibWFjOldpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcIm1hYzpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcIm1hYzpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcIm1hYzpXaW5kb3dab29tUmVzZXRcIixcblx0fSxcblx0TGludXg6IHtcblx0XHRBcHBsaWNhdGlvblN0YXJ0dXA6IFwibGludXg6QXBwbGljYXRpb25TdGFydHVwXCIsXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcImxpbnV4OlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0RlbGV0ZUV2ZW50OiBcImxpbnV4OldpbmRvd0RlbGV0ZUV2ZW50XCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJsaW51eDpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImxpbnV4OldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0ZvY3VzSW46IFwibGludXg6V2luZG93Rm9jdXNJblwiLFxuXHRcdFdpbmRvd0ZvY3VzT3V0OiBcImxpbnV4OldpbmRvd0ZvY3VzT3V0XCIsXG5cdFx0V2luZG93TG9hZENoYW5nZWQ6IFwibGludXg6V2luZG93TG9hZENoYW5nZWRcIixcblx0fSxcblx0Q29tbW9uOiB7XG5cdFx0QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZTogXCJjb21tb246QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZVwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImNvbW1vbjpUaGVtZUNoYW5nZWRcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJjb21tb246V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJjb21tb246V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RmlsZXNEcm9wcGVkXCIsXG5cdFx0V2luZG93Rm9jdXM6IFwiY29tbW9uOldpbmRvd0ZvY3VzXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd0hpZGU6IFwiY29tbW9uOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dMb3N0Rm9jdXM6IFwiY29tbW9uOldpbmRvd0xvc3RGb2N1c1wiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwiY29tbW9uOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dSdW50aW1lUmVhZHk6IFwiY29tbW9uOldpbmRvd1J1bnRpbWVSZWFkeVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Wm9vbTogXCJjb21tb246V2luZG93Wm9vbVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJjb21tb246V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJjb21tb246V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJjb21tb246V2luZG93Wm9vbVJlc2V0XCIsXG5cdH0sXG59O1xuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIExvZ3MgYSBtZXNzYWdlIHRvIHRoZSBjb25zb2xlIHdpdGggY3VzdG9tIGZvcm1hdHRpbmcuXG4gKiBAcGFyYW0ge3N0cmluZ30gbWVzc2FnZSAtIFRoZSBtZXNzYWdlIHRvIGJlIGxvZ2dlZC5cbiAqIEByZXR1cm4ge3ZvaWR9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlKSB7XG4gICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgY29uc29sZS5sb2coXG4gICAgICAgICclYyB3YWlsczMgJWMgJyArIG1lc3NhZ2UgKyAnICcsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjYWEwMDAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogM3B4IDBweCAwcHggM3B4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXG4gICAgKTtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgYnJvd3NlciBzdXBwb3J0cyByZW1vdmluZyBsaXN0ZW5lcnMgYnkgdHJpZ2dlcmluZyBhbiBBYm9ydFNpZ25hbFxuICogKHNlZSBodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZy9lbi1VUy9kb2NzL1dlYi9BUEkvRXZlbnRUYXJnZXQvYWRkRXZlbnRMaXN0ZW5lciNzaWduYWwpXG4gKlxuICogQHJldHVybiB7Ym9vbGVhbn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGNhbkFib3J0TGlzdGVuZXJzKCkge1xuICAgIGlmICghRXZlbnRUYXJnZXQgfHwgIUFib3J0U2lnbmFsIHx8ICFBYm9ydENvbnRyb2xsZXIpXG4gICAgICAgIHJldHVybiBmYWxzZTtcblxuICAgIGxldCByZXN1bHQgPSB0cnVlO1xuXG4gICAgY29uc3QgdGFyZ2V0ID0gbmV3IEV2ZW50VGFyZ2V0KCk7XG4gICAgY29uc3QgY29udHJvbGxlciA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB0YXJnZXQuYWRkRXZlbnRMaXN0ZW5lcigndGVzdCcsICgpID0+IHsgcmVzdWx0ID0gZmFsc2U7IH0sIHsgc2lnbmFsOiBjb250cm9sbGVyLnNpZ25hbCB9KTtcbiAgICBjb250cm9sbGVyLmFib3J0KCk7XG4gICAgdGFyZ2V0LmRpc3BhdGNoRXZlbnQobmV3IEN1c3RvbUV2ZW50KCd0ZXN0JykpO1xuXG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqKlxuIFRoaXMgdGVjaG5pcXVlIGZvciBwcm9wZXIgbG9hZCBkZXRlY3Rpb24gaXMgdGFrZW4gZnJvbSBIVE1YOlxuXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcblxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXG5cbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcblxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxuXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cbiBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlIGRpc3RyaWJ1dGlvbi5cblxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXG4gQU5EIEFOWSBFWFBSRVNTIE9SIElNUExJRUQgV0FSUkFOVElFUywgSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFRIRVxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcbiBGT1IgQU5ZIERJUkVDVCwgSU5ESVJFQ1QsIElOQ0lERU5UQUwsIFNQRUNJQUwsIEVYRU1QTEFSWSwgT1IgQ09OU0VRVUVOVElBTFxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcbiBDQVVTRUQgQU5EIE9OIEFOWSBUSEVPUlkgT0YgTElBQklMSVRZLCBXSEVUSEVSIElOIENPTlRSQUNULCBTVFJJQ1QgTElBQklMSVRZLFxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cblxuICoqKi9cblxubGV0IGlzUmVhZHkgPSBmYWxzZTtcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCAoKSA9PiBpc1JlYWR5ID0gdHJ1ZSk7XG5cbmV4cG9ydCBmdW5jdGlvbiB3aGVuUmVhZHkoY2FsbGJhY2spIHtcbiAgICBpZiAoaXNSZWFkeSB8fCBkb2N1bWVudC5yZWFkeVN0YXRlID09PSAnY29tcGxldGUnKSB7XG4gICAgICAgIGNhbGxiYWNrKCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGNhbGxiYWNrKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLy8gSW1wb3J0IHNjcmVlbiBqc2RvYyBkZWZpbml0aW9uIGZyb20gLi9zY3JlZW5zLmpzXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3NjcmVlbnNcIikuU2NyZWVufSBTY3JlZW5cbiAqL1xuXG5cbi8qKlxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgcG9zaXRpb24gb2YgYSB3aW5kb3cuXG4gKlxuICogQHR5cGVkZWYge09iamVjdH0gUG9zaXRpb25cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSB4IC0gVGhlIGhvcml6b250YWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICogQHByb3BlcnR5IHtudW1iZXJ9IHkgLSBUaGUgdmVydGljYWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICovXG5cblxuLyoqXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBzaXplIG9mIGEgd2luZG93LlxuICpcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFNpemVcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSB3aWR0aCAtIFRoZSB3aWR0aCBvZiB0aGUgd2luZG93XG4gKiBAcHJvcGVydHkge251bWJlcn0gaGVpZ2h0IC0gVGhlIGhlaWdodCBvZiB0aGUgd2luZG93XG4gKi9cblxuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmNvbnN0IFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDA7XG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xuY29uc3QgQ2xvc2VNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gMjtcbmNvbnN0IERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICA9IDM7XG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xuY29uc3QgRm9jdXNNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNTtcbmNvbnN0IEZvcmNlUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICA9IDY7XG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xuY29uc3QgR2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gODtcbmNvbnN0IEdldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDk7XG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcbmNvbnN0IEhpZGVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDExO1xuY29uc3QgSXNGb2N1c2VkTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMTI7XG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcbmNvbnN0IElzTWF4aW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE0O1xuY29uc3QgSXNNaW5pbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTU7XG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcbmNvbnN0IE1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE3O1xuY29uc3QgTmFtZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTg7XG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcbmNvbnN0IFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICAgICA9IDIwO1xuY29uc3QgUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMjE7XG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcbmNvbnN0IFJlc3RvcmVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDIzO1xuY29uc3QgU2V0UG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgID0gMjQ7XG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcbmNvbnN0IFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QgICAgICAgICA9IDI2O1xuY29uc3QgU2V0RnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgICAgID0gMjc7XG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcbmNvbnN0IFNldE1heFNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDI5O1xuY29uc3QgU2V0TWluU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMzA7XG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcbmNvbnN0IFNldFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICA9IDMyO1xuY29uc3QgU2V0U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzM7XG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcbmNvbnN0IFNldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDM1O1xuY29uc3QgU2hvd01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzY7XG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcbmNvbnN0IFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICA9IDM4O1xuY29uc3QgVG9nZ2xlTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgID0gMzk7XG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MDtcbmNvbnN0IFVuTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQxO1xuY29uc3QgVW5NaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDI7XG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0MztcbmNvbnN0IFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDQ0O1xuY29uc3QgWm9vbUluTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gNDU7XG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NjtcbmNvbnN0IFpvb21SZXNldE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDQ3O1xuXG4vKipcbiAqIEB0eXBlIHtzeW1ib2x9XG4gKi9cbmNvbnN0IGNhbGxlciA9IFN5bWJvbCgpO1xuXG5leHBvcnQgY2xhc3MgV2luZG93IHtcbiAgICAvKipcbiAgICAgKiBJbml0aWFsaXNlcyBhIHdpbmRvdyBvYmplY3Qgd2l0aCB0aGUgc3BlY2lmaWVkIG5hbWUuXG4gICAgICpcbiAgICAgKiBAcHJpdmF0ZVxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHRhcmdldCB3aW5kb3cuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobmFtZSA9ICcnKSB7XG4gICAgICAgIC8qKlxuICAgICAgICAgKiBAcHJpdmF0ZVxuICAgICAgICAgKiBAbmFtZSB7QGxpbmsgY2FsbGVyfVxuICAgICAgICAgKiBAdHlwZSB7KC4uLmFyZ3M6IGFueVtdKSA9PiBhbnl9XG4gICAgICAgICAqL1xuICAgICAgICB0aGlzW2NhbGxlcl0gPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcblxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xuICAgICAgICAgICAgaWYgKFxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mIHRoaXNbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXG4gICAgICAgICAgICApIHtcbiAgICAgICAgICAgICAgICB0aGlzW21ldGhvZF0gPSB0aGlzW21ldGhvZF0uYmluZCh0aGlzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cbiAgICAgKiBAcmV0dXJuIHtXaW5kb3d9IC0gVGhlIGNvcnJlc3BvbmRpbmcgd2luZG93IG9iamVjdC5cbiAgICAgKi9cbiAgICBHZXQobmFtZSkge1xuICAgICAgICByZXR1cm4gbmV3IFdpbmRvdyhuYW1lKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8UG9zaXRpb24+fSAtIFRoZSBjdXJyZW50IGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUG9zaXRpb24oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oUG9zaXRpb25NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdyBvbiB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgQ2VudGVyKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKENlbnRlck1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBDbG9zZSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShDbG9zZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgRGlzYWJsZVNpemVDb25zdHJhaW50cygpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBFbmFibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIEZvY3VzKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKEZvY3VzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBGb3JjZXMgdGhlIHdpbmRvdyB0byByZWxvYWQgdGhlIHBhZ2UgYXNzZXRzLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgRm9yY2VSZWxvYWQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oRm9yY2VSZWxvYWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERvYy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIEZ1bGxzY3JlZW4oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgc2NyZWVuIHRoYXQgdGhlIHdpbmRvdyBpcyBvbi5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPFNjcmVlbj59IC0gVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvblxuICAgICAqL1xuICAgIEdldFNjcmVlbigpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShHZXRTY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8bnVtYmVyPn0gLSBUaGUgY3VycmVudCB6b29tIGxldmVsXG4gICAgICovXG4gICAgR2V0Wm9vbSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShHZXRab29tTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPG51bWJlcj59IC0gVGhlIGN1cnJlbnQgaGVpZ2h0IG9mIHRoZSB3aW5kb3dcbiAgICAgKi9cbiAgICBIZWlnaHQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSGVpZ2h0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBIaWRlcyB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgSGlkZSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShIaWRlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8Ym9vbGVhbj59IC0gV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkXG4gICAgICovXG4gICAgSXNGb2N1c2VkKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKElzRm9jdXNlZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZnVsbHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZnVsbHNjcmVlblxuICAgICAqL1xuICAgIElzRnVsbHNjcmVlbigpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShJc0Z1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1heGltaXNlZC5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWF4aW1pc2VkXG4gICAgICovXG4gICAgSXNNYXhpbWlzZWQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSXNNYXhpbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPGJvb2xlYW4+fSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkXG4gICAgICovXG4gICAgSXNNaW5pbWlzZWQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oSXNNaW5pbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgTWF4aW1pc2UoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgTWluaW1pc2UoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvd1xuICAgICAqL1xuICAgIE5hbWUoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oTmFtZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogT3BlbnMgdGhlIGRldmVsb3BtZW50IHRvb2xzIHBhbmUuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBPcGVuRGV2VG9vbHMoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oT3BlbkRldlRvb2xzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxQb3NpdGlvbj59IC0gVGhlIGN1cnJlbnQgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICAgICAqL1xuICAgIFJlbGF0aXZlUG9zaXRpb24oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBSZWxvYWQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oUmVsb2FkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxib29sZWFuPn0gLSBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IHJlc2l6YWJsZVxuICAgICAqL1xuICAgIFJlc2l6YWJsZSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShSZXNpemFibGVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlc3RvcmVzIHRoZSB3aW5kb3cgdG8gaXRzIHByZXZpb3VzIHN0YXRlIGlmIGl0IHdhcyBwcmV2aW91c2x5IG1pbmltaXNlZCwgbWF4aW1pc2VkIG9yIGZ1bGxzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBSZXN0b3JlKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFJlc3RvcmVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcGFyYW0ge251bWJlcn0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldFBvc2l0aW9uKHgsIHkpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gYWx3YXlzT25Ub3AgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIHN0YXkgb24gdG9wXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBTZXRBbHdheXNPblRvcChhbHdheXNPblRvcCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHIgLSBUaGUgZGVzaXJlZCByZWQgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZFxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZFxuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBiIC0gVGhlIGRlc2lyZWQgYmx1ZSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IGEgLSBUaGUgZGVzaXJlZCBhbHBoYSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyKHIsIGcsIGIsIGEpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtib29sZWFufSBmcmFtZWxlc3MgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGZyYW1lbGVzc1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0RnJhbWVsZXNzKGZyYW1lbGVzcykge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcGFyYW0ge2Jvb2xlYW59IGVuYWJsZWQgLSBXaGV0aGVyIHRoZSBmdWxsc2NyZWVuIGJ1dHRvbiBzaG91bGQgYmUgZW5hYmxlZFxuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWQoZW5hYmxlZCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIGRlc2lyZWQgbWF4aW11bSB3aWR0aCBvZiB0aGUgd2luZG93XG4gICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNldE1heFNpemUod2lkdGgsIGhlaWdodCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGggLSBUaGUgZGVzaXJlZCBtaW5pbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3dcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWluaW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvd1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0TWluU2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0TWluU2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcGFyYW0ge251bWJlcn0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvd1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0UmVsYXRpdmVQb3NpdGlvbih4LCB5KSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcGFyYW0ge2Jvb2xlYW59IHJlc2l6YWJsZSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgcmVzaXphYmxlXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIC0gVGhlIGRlc2lyZWQgd2lkdGggb2YgdGhlIHdpbmRvd1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvd1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0U2l6ZSh3aWR0aCwgaGVpZ2h0KSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IHRpdGxlIC0gVGhlIGRlc2lyZWQgdGl0bGUgb2YgdGhlIHdpbmRvd1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgU2V0VGl0bGUodGl0bGUpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IHpvb20gLSBUaGUgZGVzaXJlZCB6b29tIGxldmVsXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBTZXRab29tKHpvb20pIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFNob3coKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oU2hvd01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8U2l6ZT59IC0gVGhlIGN1cnJlbnQgc2l6ZSBvZiB0aGUgd2luZG93XG4gICAgICovXG4gICAgU2l6ZSgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShTaXplTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmdWxsc2NyZWVuIGFuZCBub3JtYWwuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBUb2dnbGVGdWxsc2NyZWVuKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxuICAgICAqXG4gICAgICogQHB1YmxpY1xuICAgICAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XG4gICAgICovXG4gICAgVG9nZ2xlTWF4aW1pc2UoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oVG9nZ2xlTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLWZ1bGxzY3JlZW5zIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBVbkZ1bGxzY3JlZW4oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oVW5GdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1tYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFVuTWF4aW1pc2UoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oVW5NYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBVbk1pbmltaXNlKCkge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJdKFVuTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTxudW1iZXI+fSAtIFRoZSBjdXJyZW50IHdpZHRoIG9mIHRoZSB3aW5kb3dcbiAgICAgKi9cbiAgICBXaWR0aCgpIHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyXShXaWR0aE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogWm9vbXMgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFpvb20oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oWm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBab29tSW4oKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oWm9vbUluTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEZWNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKlxuICAgICAqIEBwdWJsaWNcbiAgICAgKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICAgICAqL1xuICAgIFpvb21PdXQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oWm9vbU91dE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVzZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICpcbiAgICAgKiBAcHVibGljXG4gICAgICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAgICAgKi9cbiAgICBab29tUmVzZXQoKSB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlcl0oWm9vbVJlc2V0TWV0aG9kKTtcbiAgICB9XG59XG5cbi8qKlxuICogVGhlIHdpbmRvdyB3aXRoaW4gd2hpY2ggdGhlIHNjcmlwdCBpcyBydW5uaW5nLlxuICpcbiAqIEB0eXBlIHtXaW5kb3d9XG4gKi9cbmNvbnN0IHRoaXNXaW5kb3cgPSBuZXcgV2luZG93KCcnKTtcblxuZXhwb3J0IGRlZmF1bHQgdGhpc1dpbmRvdztcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0ICogYXMgUnVudGltZSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmNcIjtcblxuLy8gTk9URTogdGhlIGZvbGxvd2luZyBtZXRob2RzIE1VU1QgYmUgaW1wb3J0ZWQgZXhwbGljaXRseSBiZWNhdXNlIG9mIGhvdyBlc2J1aWxkIGluamVjdGlvbiB3b3Jrc1xuaW1wb3J0IHtFbmFibGUgYXMgRW5hYmxlV01MfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XG5pbXBvcnQge2RlYnVnTG9nfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHNcIjtcblxud2luZG93LndhaWxzID0gUnVudGltZTtcbkVuYWJsZVdNTCgpO1xuXG5pZiAoREVCVUcpIHtcbiAgICBkZWJ1Z0xvZyhcIldhaWxzIFJ1bnRpbWUgTG9hZGVkXCIpXG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLlN5c3RlbSwgJycpO1xuY29uc3Qgc3lzdGVtSXNEYXJrTW9kZSA9IDA7XG5jb25zdCBlbnZpcm9ubWVudCA9IDE7XG5cbmNvbnN0IF9pbnZva2UgPSAoKCkgPT4ge1xuICAgIHRyeSB7XG4gICAgICAgIGlmKHdpbmRvdz8uY2hyb21lPy53ZWJ2aWV3KSB7XG4gICAgICAgICAgICByZXR1cm4gKG1zZykgPT4gd2luZG93LmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlKG1zZyk7XG4gICAgICAgIH1cbiAgICAgICAgaWYod2luZG93Py53ZWJraXQ/Lm1lc3NhZ2VIYW5kbGVycz8uZXh0ZXJuYWwpIHtcbiAgICAgICAgICAgIHJldHVybiAobXNnKSA9PiB3aW5kb3cud2Via2l0Lm1lc3NhZ2VIYW5kbGVycy5leHRlcm5hbC5wb3N0TWVzc2FnZShtc2cpO1xuICAgICAgICB9XG4gICAgfSBjYXRjaChlKSB7XG4gICAgICAgIGNvbnNvbGUud2FybignXFxuJWNcdTI2QTBcdUZFMEYgQnJvd3NlciBFbnZpcm9ubWVudCBEZXRlY3RlZCAlY1xcblxcbiVjT25seSBVSSBwcmV2aWV3cyBhcmUgYXZhaWxhYmxlIGluIHRoZSBicm93c2VyLiBGb3IgZnVsbCBmdW5jdGlvbmFsaXR5LCBwbGVhc2UgcnVuIHRoZSBhcHBsaWNhdGlvbiBpbiBkZXNrdG9wIG1vZGUuXFxuTW9yZSBpbmZvcm1hdGlvbiBhdDogaHR0cHM6Ly92My53YWlscy5pby9sZWFybi9idWlsZC8jdXNpbmctYS1icm93c2VyLWZvci1kZXZlbG9wbWVudFxcbicsXG4gICAgICAgICAgICAnYmFja2dyb3VuZDogI2ZmZmZmZjsgY29sb3I6ICMwMDAwMDA7IGZvbnQtd2VpZ2h0OiBib2xkOyBwYWRkaW5nOiA0cHggOHB4OyBib3JkZXItcmFkaXVzOiA0cHg7IGJvcmRlcjogMnB4IHNvbGlkICMwMDAwMDA7JyxcbiAgICAgICAgICAgICdiYWNrZ3JvdW5kOiB0cmFuc3BhcmVudDsnLFxuICAgICAgICAgICAgJ2NvbG9yOiAjZmZmZmZmOyBmb250LXN0eWxlOiBpdGFsaWM7IGZvbnQtd2VpZ2h0OiBib2xkOycpO1xuICAgIH1cbiAgICByZXR1cm4gbnVsbDtcbn0pKCk7XG5cbmV4cG9ydCBmdW5jdGlvbiBpbnZva2UobXNnKSB7XG4gICAgaWYgKCFfaW52b2tlKSByZXR1cm47XG4gICAgcmV0dXJuIF9pbnZva2UobXNnKTtcbn1cblxuLyoqXG4gKiBAZnVuY3Rpb25cbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxib29sZWFuPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoc3lzdGVtSXNEYXJrTW9kZSk7XG59XG5cbi8qKlxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXG4gKlxuICogQGFzeW5jXG4gKiBAZnVuY3Rpb24gQ2FwYWJpbGl0aWVzXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxPYmplY3Q+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyB0aGUgY2FwYWJpbGl0aWVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FwYWJpbGl0aWVzKCkge1xuICAgIGxldCByZXNwb25zZSA9IGZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKTtcbiAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xufVxuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IE9TSW5mb1xuICogQHByb3BlcnR5IHtzdHJpbmd9IEJyYW5kaW5nIC0gVGhlIGJyYW5kaW5nIG9mIHRoZSBPUy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJRCAtIFRoZSBJRCBvZiB0aGUgT1MuXG4gKiBAcHJvcGVydHkge3N0cmluZ30gTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBPUy5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBWZXJzaW9uIC0gVGhlIHZlcnNpb24gb2YgdGhlIE9TLlxuICovXG5cbi8qKlxuICogQHR5cGVkZWYge09iamVjdH0gRW52aXJvbm1lbnRJbmZvXG4gKiBAcHJvcGVydHkge3N0cmluZ30gQXJjaCAtIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS5cbiAqIEBwcm9wZXJ0eSB7Ym9vbGVhbn0gRGVidWcgLSBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBPUyAtIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS5cbiAqIEBwcm9wZXJ0eSB7T1NJbmZvfSBPU0luZm8gLSBEZXRhaWxzIG9mIHRoZSBvcGVyYXRpbmcgc3lzdGVtLlxuICogQHByb3BlcnR5IHtPYmplY3R9IFBsYXRmb3JtSW5mbyAtIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uXG4gKi9cblxuLyoqXG4gKiBAZnVuY3Rpb25cbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxuICogQHJldHVybnMge1Byb21pc2U8RW52aXJvbm1lbnRJbmZvPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyBPUyBhbmQgc3lzdGVtIGFyY2hpdGVjdHVyZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCkge1xuICAgIHJldHVybiBjYWxsKGVudmlyb25tZW50KTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxuICpcbiAqIEByZXR1cm4ge2Jvb2xlYW59IFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCkge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcIndpbmRvd3NcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cbiAqXG4gKiBAcmV0dXJucyB7Ym9vbGVhbn0gUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwibGludXhcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxuICpcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBtYWNPUywgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNNYWMoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwiZGFyd2luXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFtZDY0XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXG4gKlxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBUk0oKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm1cIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxuICpcbiAqIEByZXR1cm5zIHtib29sZWFufSAtIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKSB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm02NFwiO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5EZWJ1ZyA9PT0gdHJ1ZTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcbmltcG9ydCB7SXNEZWJ1Z30gZnJvbSBcIi4vc3lzdGVtXCI7XG5cbi8vIHNldHVwXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5Db250ZXh0TWVudSwgJycpO1xuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkLCB4LCB5LCBkYXRhKSB7XG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XG59XG5cbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudCkge1xuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XG4gICAgbGV0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgbGV0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpO1xuICAgIGN1c3RvbUNvbnRleHRNZW51ID0gY3VzdG9tQ29udGV4dE1lbnUgPyBjdXN0b21Db250ZXh0TWVudS50cmltKCkgOiBcIlwiO1xuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnVEYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgY3VzdG9tQ29udGV4dE1lbnVEYXRhKTtcbiAgICAgICAgcmV0dXJuXG4gICAgfVxuXG4gICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCk7XG59XG5cblxuLypcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG5cblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXG4qL1xuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCkge1xuXG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XG4gICAgaWYgKElzRGVidWcoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcbiAgICAgICAgY2FzZSBcInNob3dcIjpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQuaXNDb250ZW50RWRpdGFibGUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWRcbiAgICAgICAgICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICAgICAgICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IChzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKVxuICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbikge1xuICAgICAgICAgICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0cyA9IHJhbmdlLmdldENsaWVudFJlY3RzKCk7XG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHJlY3QubGVmdCwgcmVjdC50b3ApID09PSBlbGVtZW50KSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGFnbmFtZSBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQudGFnTmFtZSA9PT0gXCJJTlBVVFwiIHx8IGVsZW1lbnQudGFnTmFtZSA9PT0gXCJURVhUQVJFQVwiKSB7XG4gICAgICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIWVsZW1lbnQucmVhZE9ubHkgJiYgIWVsZW1lbnQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGtleVN0cmluZyAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cbiAqIEByZXR1cm4geyp9IC0gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5U3RyaW5nKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZmxhZ3Nba2V5U3RyaW5nXTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleVN0cmluZyArIFwiJzogXCIgKyBlKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQge2ludm9rZSwgSXNXaW5kb3dzfSBmcm9tIFwiLi9zeXN0ZW1cIjtcbmltcG9ydCB7R2V0RmxhZ30gZnJvbSBcIi4vZmxhZ3NcIjtcblxuLy8gU2V0dXBcbmxldCBzaG91bGREcmFnID0gZmFsc2U7XG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5sZXQgcmVzaXplRWRnZSA9IG51bGw7XG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xuXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxud2luZG93Ll93YWlscy5zZXRSZXNpemFibGUgPSBmdW5jdGlvbih2YWx1ZSkge1xuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xufTtcblxud2luZG93Ll93YWlscy5lbmREcmFnID0gZnVuY3Rpb24oKSB7XG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufTtcblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIG9uTW91c2VEb3duKTtcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBvbk1vdXNlTW92ZSk7XG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIG9uTW91c2VVcCk7XG5cblxuZnVuY3Rpb24gZHJhZ1Rlc3QoZSkge1xuICAgIGxldCB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpO1xuICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgaWYgKCF2YWwgfHwgdmFsID09PSBcIlwiIHx8IHZhbC50cmltKCkgIT09IFwiZHJhZ1wiIHx8IG1vdXNlUHJlc3NlZCA9PT0gMCkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xuXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgaW52b2tlKFwid2FpbHM6cmVzaXplOlwiICsgcmVzaXplRWRnZSk7XG4gICAgICAgIGUucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmIChkcmFnVGVzdChlKSkge1xuICAgICAgICAvLyBUaGlzIGNoZWNrcyBmb3IgY2xpY2tzIG9uIHRoZSBzY3JvbGwgYmFyXG4gICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBzaG91bGREcmFnID0gdHJ1ZTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBvbk1vdXNlVXAoKSB7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCBkZWZhdWx0Q3Vyc29yO1xuICAgIHJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGUpIHtcbiAgICBpZiAoc2hvdWxkRHJhZykge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgICAgICBpbnZva2UoXCJ3YWlsczpkcmFnXCIpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmICghcmVzaXphYmxlIHx8ICFJc1dpbmRvd3MoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuICAgIGlmIChkZWZhdWx0Q3Vyc29yID09IG51bGwpIHtcbiAgICAgICAgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3I7XG4gICAgfVxuICAgIGxldCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGxldCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcblxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhc1xuICAgIGxldCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcblxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG5cbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lcnNcbiAgICBsZXQgcmlnaHRDb3JuZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgbGVmdENvcm5lciA9IGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgdG9wQ29ybmVyID0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgYm90dG9tQ29ybmVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcblxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgcmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbiAgICAvLyBBZGp1c3RlZCBmb3IgY29ybmVyIGFyZWFzXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIHRvcENvcm5lcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24sICcnKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICpcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqXG4gKiBAcmV0dXJuIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcbn1cblxuLyoqXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXG4gKlxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZVwiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcbndpbmRvdy5fd2FpbHMuY2FsbEVycm9ySGFuZGxlciA9IGVycm9ySGFuZGxlcjtcblxuXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYWxsLCAnJyk7XG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DYW5jZWxDYWxsLCAnJyk7XG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxuICpcbiAqIEByZXR1cm4ge3N0cmluZ30gLSBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBjYWxsIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxuICogQHBhcmFtIHtib29sZWFufSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cbiAqXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9IC0gVGhpcyBtZXRob2QgZG9lcyBub3QgcmV0dXJuIGFueSB2YWx1ZS5cbiAqL1xuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XG4gICAgaWYgKHByb21pc2VIYW5kbGVyKSB7XG4gICAgICAgIGlmICghZGF0YSkge1xuICAgICAgICAgICAgcHJvbWlzZUhhbmRsZXIucmVzb2x2ZSgpO1xuICAgICAgICB9IGVsc2UgaWYgKCFpc0pTT04pIHtcbiAgICAgICAgICAgIHByb21pc2VIYW5kbGVyLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHByb21pc2VIYW5kbGVyLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobmV3IFR5cGVFcnJvcihcImNvdWxkIG5vdCBwYXJzZSByZXN1bHQ6IFwiICsgZXJyLm1lc3NhZ2UsIHsgY2F1c2U6IGVyciB9KSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgZXJyb3IgZnJvbSBhIGNhbGwgcmVxdWVzdC5cbiAqXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIHByb21pc2UgaGFuZGxlci5cbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIGVycm9yIGRhdGEgdG8gcmVqZWN0IHRoZSBwcm9taXNlIGhhbmRsZXIgd2l0aC5cbiAqIEBwYXJhbSB7Ym9vbGVhbn0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXG4gKlxuICogQHJldHVybiB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gZXJyb3JIYW5kbGVyKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBjb25zdCBwcm9taXNlSGFuZGxlciA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcbiAgICAgICAgaWYgKCFpc0pTT04pIHtcbiAgICAgICAgICAgIHByb21pc2VIYW5kbGVyLnJlamVjdChuZXcgRXJyb3IoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgbGV0IGVycm9yO1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBlcnJvciA9IEpTT04ucGFyc2UoZGF0YSk7XG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobmV3IFR5cGVFcnJvcihcImNvdWxkIG5vdCBwYXJzZSBlcnJvcjogXCIgKyBlcnIubWVzc2FnZSwgeyBjYXVzZTogZXJyIH0pKTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIGxldCBvcHRpb25zID0ge307XG4gICAgICAgICAgICBpZiAoZXJyb3IuY2F1c2UpIHtcbiAgICAgICAgICAgICAgICBvcHRpb25zLmNhdXNlID0gZXJyb3IuY2F1c2U7XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIGxldCBleGNlcHRpb247XG4gICAgICAgICAgICBzd2l0Y2ggKGVycm9yLmtpbmQpIHtcbiAgICAgICAgICAgICAgICBjYXNlIFwiUmVmZXJlbmNlRXJyb3JcIjpcbiAgICAgICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFJlZmVyZW5jZUVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgICAgICBjYXNlIFwiVHlwZUVycm9yXCI6XG4gICAgICAgICAgICAgICAgICAgIGV4Y2VwdGlvbiA9IG5ldyBUeXBlRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgICAgIGNhc2UgXCJSdW50aW1lRXJyb3JcIjpcbiAgICAgICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFJ1bnRpbWVFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IEVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgcHJvbWlzZUhhbmRsZXIucmVqZWN0KGV4Y2VwdGlvbik7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogUmV0cmlldmVzIGFuZCByZW1vdmVzIHRoZSByZXNwb25zZSBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElEIGZyb20gdGhlIGNhbGxSZXNwb25zZXMgbWFwLlxuICpcbiAqIEBwYXJhbSB7YW55fSBpZCAtIFRoZSBJRCBvZiB0aGUgcmVzcG9uc2UgdG8gYmUgcmV0cmlldmVkIGFuZCByZW1vdmVkLlxuICpcbiAqIEByZXR1cm5zIHthbnl9IFRoZSByZXNwb25zZSBvYmplY3QgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRC5cbiAqL1xuZnVuY3Rpb24gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpIHtcbiAgICBjb25zdCByZXNwb25zZSA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgcmV0dXJuIHJlc3BvbnNlO1xufVxuXG4vKipcbiAqIENvbGxlY3RzIGFsbCByZXF1aXJlZCBpbmZvcm1hdGlvbiBmb3IgYSBiaW5kaW5nIGNhbGwuXG4gKlxuICogQHR5cGVkZWYge09iamVjdH0gQ2FsbE9wdGlvbnNcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBbbWV0aG9kSURdIC0gVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IFttZXRob2ROYW1lXSAtIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuXG4gKiBAcHJvcGVydHkge2FueVtdfSBhcmdzIC0gQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuXG4gKi9cblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHRocm93biBpbiBjYXNlIHRoZSBib3VuZCBtZXRob2QgcmV0dXJucyBhbiBlcnJvci5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIFJ1bnRpbWVFcnJvciNuYW1lfSBwcm9wZXJ0eSBpcyBcIlJ1bnRpbWVFcnJvclwiLlxuICovXG5leHBvcnQgY2xhc3MgUnVudGltZUVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgUnVudGltZUVycm9yIGluc3RhbmNlLlxuICAgICAqXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cbiAgICAgKiBAcGFyYW0ge2FueVtdfSBhcmdzIC0gT3B0aW9uYWwgYXJndW1lbnRzIGZvciB0aGUgRXJyb3IgY29uc3RydWN0b3IuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobWVzc2FnZSwgLi4uYXJncykge1xuICAgICAgICBzdXBlcihtZXNzYWdlLCAuLi5hcmdzKTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJSdW50aW1lRXJyb3JcIjtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbCBhIGJvdW5kIG1ldGhvZCBhY2NvcmRpbmcgdG8gdGhlIGdpdmVuIGNhbGwgb3B0aW9ucy5cbiAqXG4gKiBJbiBjYXNlIG9mIGZhaWx1cmUsIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYW4gZXhjZXB0aW9uXG4gKiBhbW9uZyBSZWZlcmVuY2VFcnJvciAodW5rbm93biBtZXRob2QpLCBUeXBlRXJyb3IgKHdyb25nIGFyZ3VtZW50IGNvdW50IG9yIHR5cGUpLFxuICoge0BsaW5rIFJ1bnRpbWVFcnJvcn0gKG1ldGhvZCByZXR1cm5lZCBhbiBlcnJvciksIG9yIG90aGVyIChuZXR3b3JrIG9yIGludGVybmFsIGVycm9ycykuXG4gKiBUaGUgZXhjZXB0aW9uIG1pZ2h0IGhhdmUgYSBcImNhdXNlXCIgZmllbGQgd2l0aCB0aGUgdmFsdWUgcmV0dXJuZWRcbiAqIGJ5IHRoZSBhcHBsaWNhdGlvbi0gb3Igc2VydmljZS1sZXZlbCBlcnJvciBtYXJzaGFsaW5nIGZ1bmN0aW9ucy5cbiAqXG4gKiBAcGFyYW0ge0NhbGxPcHRpb25zfSBvcHRpb25zIC0gQSBtZXRob2QgY2FsbCBkZXNjcmlwdG9yLlxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zKSB7XG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgY29uc3QgZG9DYW5jZWwgPSAoKSA9PiB7IHJldHVybiBjYW5jZWxDYWxsKHR5cGUsIHtcImNhbGwtaWRcIjogaWR9KSB9O1xuICAgIGxldCBxdWV1ZWRDYW5jZWwgPSBmYWxzZSwgY2FsbFJ1bm5pbmcgPSBmYWxzZTtcbiAgICBsZXQgcCA9IG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgb3B0aW9uc1tcImNhbGwtaWRcIl0gPSBpZDtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xuICAgICAgICBjYWxsKENhbGxCaW5kaW5nLCBvcHRpb25zKS50aGVuKChfKSA9PiB7XG4gICAgICAgICAgICBjYWxsUnVubmluZyA9IHRydWU7XG4gICAgICAgICAgICBpZiAocXVldWVkQ2FuY2VsKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIGRvQ2FuY2VsKCk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG4gICAgcC5jYW5jZWwgPSAoKSA9PiB7XG4gICAgICAgIGlmIChjYWxsUnVubmluZykge1xuICAgICAgICAgICAgcmV0dXJuIGRvQ2FuY2VsKCk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBxdWV1ZWRDYW5jZWwgPSB0cnVlO1xuICAgICAgICB9XG4gICAgfTtcblxuICAgIHJldHVybiBwO1xufVxuXG4vKipcbiAqIENhbGxzIGEgYm91bmQgbWV0aG9kIGJ5IG5hbWUgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXG4gKiBAcGFyYW0ge2FueVtdfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxhbnk+fSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShtZXRob2ROYW1lLCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIENhbGwoe1xuICAgICAgICBtZXRob2ROYW1lLFxuICAgICAgICBhcmdzXG4gICAgfSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIG51bWVyaWMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXG4gKlxuICogQHBhcmFtIHtudW1iZXJ9IG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqIEBwYXJhbSB7YW55W119IGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEByZXR1cm4ge1Byb21pc2U8YW55Pn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQsIC4uLmFyZ3MpIHtcbiAgICByZXR1cm4gQ2FsbCh7XG4gICAgICAgIG1ldGhvZElELFxuICAgICAgICBhcmdzXG4gICAgfSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DbGlwYm9hcmQsICcnKTtcbmNvbnN0IENsaXBib2FyZFNldFRleHQgPSAwO1xuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XG5cbi8qKlxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSB0ZXh0IC0gVGhlIHRleHQgdG8gYmUgc2V0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKiBAcmV0dXJuIHtQcm9taXNlfSAtIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0KSB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSB0ZXh0IGZyb20gdGhlIENsaXBib2FyZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkVGV4dCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBBbnkgaXMgYSBkdW1teSBjcmVhdGlvbiBmdW5jdGlvbiBmb3Igc2ltcGxlIG9yIHVua25vd24gdHlwZXMuXG4gKiBAdGVtcGxhdGUgVFxuICogQHBhcmFtIHthbnl9IHNvdXJjZVxuICogQHJldHVybnMge1R9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBBbnkoc291cmNlKSB7XG4gICAgcmV0dXJuIC8qKiBAdHlwZSB7VH0gKi8oc291cmNlKTtcbn1cblxuLyoqXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXG4gKiBudWxsIHN0cmluZ3Mgd2l0aCBlbXB0eSBzdHJpbmdzLlxuICogQHBhcmFtIHthbnl9IHNvdXJjZVxuICogQHJldHVybnMge3N0cmluZ31cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5dGVTbGljZShzb3VyY2UpIHtcbiAgICByZXR1cm4gLyoqIEB0eXBlIHthbnl9ICovKChzb3VyY2UgPT0gbnVsbCkgPyBcIlwiIDogc291cmNlKTtcbn1cblxuLyoqXG4gKiBBcnJheSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XG4gKiB3aG9zZSBlbGVtZW50cyBhcmUgb2YgdGhhdCB0eXBlLlxuICogQHRlbXBsYXRlIFRcbiAqIEBwYXJhbSB7KHNvdXJjZTogYW55KSA9PiBUfSBlbGVtZW50XG4gKiBAcmV0dXJucyB7KHNvdXJjZTogYW55KSA9PiBUW119XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBBcnJheShlbGVtZW50KSB7XG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IFtdIDogc291cmNlKTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm4gW107XG4gICAgICAgIH1cbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzb3VyY2UubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgICAgIHNvdXJjZVtpXSA9IGVsZW1lbnQoc291cmNlW2ldKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogTWFwIHRha2VzIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdHdvIGFyYml0cmFyeSB0eXBlc1xuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIG9iamVjdFxuICogd2hvc2Uga2V5cyBhbmQgdmFsdWVzIGFyZSBvZiB0aG9zZSB0eXBlcy5cbiAqIEB0ZW1wbGF0ZSBLLCBWXG4gKiBAcGFyYW0geyhzb3VyY2U6IGFueSkgPT4gS30ga2V5XG4gKiBAcGFyYW0geyhzb3VyY2U6IGFueSkgPT4gVn0gdmFsdWVcbiAqIEByZXR1cm5zIHsoc291cmNlOiBhbnkpID0+IHsgW186IEtdOiBWIH19XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBNYXAoa2V5LCB2YWx1ZSkge1xuICAgIGlmICh2YWx1ZSA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8ge30gOiBzb3VyY2UpO1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybiB7fTtcbiAgICAgICAgfVxuICAgICAgICBmb3IgKGNvbnN0IGtleSBpbiBzb3VyY2UpIHtcbiAgICAgICAgICAgIHNvdXJjZVtrZXldID0gdmFsdWUoc291cmNlW2tleV0pO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBOdWxsYWJsZSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxuICogYW5kIHJldHVybnMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBudWxsYWJsZSB2YWx1ZSBvZiB0aGF0IHR5cGUuXG4gKiBAdGVtcGxhdGUgVFxuICogQHBhcmFtIHsoc291cmNlOiBhbnkpID0+IFR9IGVsZW1lbnRcbiAqIEByZXR1cm5zIHsoc291cmNlOiBhbnkpID0+IChUIHwgbnVsbCl9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBOdWxsYWJsZShlbGVtZW50KSB7XG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gQW55O1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gbnVsbCA6IGVsZW1lbnQoc291cmNlKSk7XG59XG5cbi8qKlxuICogU3RydWN0IHRha2VzIGFuIG9iamVjdCBtYXBwaW5nIGZpZWxkIG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9uc1xuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgc3RydWN0LlxuICogQHRlbXBsYXRlIHt7IFtfOiBzdHJpbmddOiAoKHNvdXJjZTogYW55KSA9PiBhbnkpIH19IFRcbiAqIEB0ZW1wbGF0ZSB7eyBbS2V5IGluIGtleW9mIFRdPzogUmV0dXJuVHlwZTxUW0tleV0+IH19IFVcbiAqIEBwYXJhbSB7VH0gY3JlYXRlRmllbGRcbiAqIEByZXR1cm5zIHsoc291cmNlOiBhbnkpID0+IFV9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQpIHtcbiAgICBsZXQgYWxsQW55ID0gdHJ1ZTtcbiAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcbiAgICAgICAgaWYgKGNyZWF0ZUZpZWxkW25hbWVdICE9PSBBbnkpIHtcbiAgICAgICAgICAgIGFsbEFueSA9IGZhbHNlO1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cbiAgICB9XG4gICAgaWYgKGFsbEFueSkge1xuICAgICAgICByZXR1cm4gQW55O1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICAgICAgaWYgKG5hbWUgaW4gc291cmNlKSB7XG4gICAgICAgICAgICAgICAgc291cmNlW25hbWVdID0gY3JlYXRlRmllbGRbbmFtZV0oc291cmNlW25hbWVdKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTaXplXG4gKiBAcHJvcGVydHkge251bWJlcn0gV2lkdGggLSBUaGUgd2lkdGguXG4gKiBAcHJvcGVydHkge251bWJlcn0gSGVpZ2h0IC0gVGhlIGhlaWdodC5cbiAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtPYmplY3R9IFJlY3RcbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBYIC0gVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFkgLSBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uXG4gKiBAcHJvcGVydHkge251bWJlcn0gV2lkdGggLSBUaGUgd2lkdGggb2YgdGhlIHJlY3RhbmdsZS5cbiAqIEBwcm9wZXJ0eSB7bnVtYmVyfSBIZWlnaHQgLSBUaGUgaGVpZ2h0IG9mIHRoZSByZWN0YW5nbGUuXG4gKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7T2JqZWN0fSBTY3JlZW5cbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBJRCAtIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLlxuICogQHByb3BlcnR5IHtzdHJpbmd9IE5hbWUgLSBIdW1hbiByZWFkYWJsZSBuYW1lIG9mIHRoZSBzY3JlZW4uXG4gKiBAcHJvcGVydHkge251bWJlcn0gU2NhbGVGYWN0b3IgLSBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFggLSBUaGUgWCBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uXG4gKiBAcHJvcGVydHkge251bWJlcn0gWSAtIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi5cbiAqIEBwcm9wZXJ0eSB7U2l6ZX0gU2l6ZSAtIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uXG4gKiBAcHJvcGVydHkge1JlY3R9IEJvdW5kcyAtIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC5cbiAqIEBwcm9wZXJ0eSB7UmVjdH0gUGh5c2ljYWxCb3VuZHMgLSBDb250YWlucyB0aGUgcGh5c2ljYWwgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQgKGJlZm9yZSBzY2FsaW5nKS5cbiAqIEBwcm9wZXJ0eSB7UmVjdH0gV29ya0FyZWEgLSBDb250YWlucyB0aGUgYXJlYSBvZiB0aGUgc2NyZWVuIHRoYXQgaXMgYWN0dWFsbHkgdXNhYmxlIChleGNsdWRpbmcgdGFza2JhciBhbmQgb3RoZXIgc3lzdGVtIFVJKS5cbiAqIEBwcm9wZXJ0eSB7UmVjdH0gUGh5c2ljYWxXb3JrQXJlYSAtIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBXb3JrQXJlYSBvZiB0aGUgc2NyZWVuIChiZWZvcmUgc2NhbGluZykuXG4gKiBAcHJvcGVydHkge2Jvb2xlYW59IElzUHJpbWFyeSAtIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLlxuICogQHByb3BlcnR5IHtudW1iZXJ9IFJvdGF0aW9uIC0gVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uXG4gKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TY3JlZW5zLCBcIlwiKTtcblxuY29uc3QgZ2V0QWxsID0gMDtcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xuY29uc3QgZ2V0Q3VycmVudCA9IDI7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gYXJyYXkgb2YgU2NyZWVuIG9iamVjdHMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKSB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcbn1cbi8qKlxuICogR2V0cyB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xufVxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpIHtcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcbn1cbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7OztBQzZCQSxJQUFJLGNBQ0E7QUFFRyxJQUFJLFNBQVMsQ0FBQyxPQUFPLE9BQU87QUFDL0IsTUFBSSxLQUFLO0FBRVQsTUFBSSxJQUFJLE9BQU87QUFDZixTQUFPLEtBQUs7QUFFUixVQUFNLFlBQWEsS0FBSyxPQUFPLElBQUksS0FBTSxDQUFDO0FBQUEsRUFDOUM7QUFDQSxTQUFPO0FBQ1g7OztBQzVCQSxJQUFNLGFBQWEsT0FBTyxTQUFTLFNBQVM7QUFHckMsSUFBTSxjQUFjO0FBQUEsRUFDdkIsTUFBTTtBQUFBLEVBQ04sV0FBVztBQUFBLEVBQ1gsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsWUFBWTtBQUNoQjtBQUNPLElBQUksV0FBVyxPQUFPO0FBc0J0QixTQUFTLHVCQUF1QixRQUFRLFlBQVk7QUFDdkQsU0FBTyxTQUFVLFFBQVEsT0FBSyxNQUFNO0FBQ2hDLFdBQU8sa0JBQWtCLFFBQVEsUUFBUSxZQUFZLElBQUk7QUFBQSxFQUM3RDtBQUNKO0FBT0EsZUFBZSxrQkFBa0IsVUFBVSxRQUFRLFlBQVksTUFBTTtBQUNqRSxNQUFJLE1BQU0sSUFBSSxJQUFJLFVBQVU7QUFDNUIsTUFBSSxZQUFZLE1BQU07QUFDbEIsUUFBSSxhQUFhLE9BQU8sVUFBVSxRQUFRO0FBQUEsRUFDOUM7QUFDQSxNQUFJLFVBQVUsTUFB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QzFCLGdCQUFnQjtBQUFBLElBQ2hCLDRCQUE0QjtBQUFBLElBQzVCLDRCQUE0QjtBQUFBLElBQzVCLHlEQUF5RDtBQUFBLElBQ3pELHNDQUFzQztBQUFBLElBQ3RDLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLHFCQUFxQjtBQUFBLElBQ3JCLHNCQUFzQjtBQUFBLElBQ3RCLGdDQUFnQztBQUFBLElBQ2hDLGtDQUFrQztBQUFBLElBQ2xDLG1DQUFtQztBQUFBLElBQ25DLG9DQUFvQztBQUFBLElBQ3BDLCtCQUErQjtBQUFBLElBQy9CLDZCQUE2QjtBQUFBLElBQzdCLHVCQUF1QjtBQUFBLElBQ3ZCLGlDQUFpQztBQUFBLElBQ2pDLDhCQUE4QjtBQUFBLElBQzlCLDRCQUE0QjtBQUFBLElBQzVCLHNDQUFzQztBQUFBLElBQ3RDLDRCQUE0QjtBQUFBLElBQzVCLHNCQUFzQjtBQUFBLElBQ3RCLGtDQUFrQztBQUFBLElBQ2xDLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLHdCQUF3QjtBQUFBLElBQ3hCLG1CQUFtQjtBQUFBLElBQ25CLDBCQUEwQjtBQUFBLElBQzFCLDhCQUE4QjtBQUFBLElBQzlCLHlCQUF5QjtBQUFBLElBQ3pCLDZCQUE2QjtBQUFBLElBQzdCLGlCQUFpQjtBQUFBLElBQ2pCLGdCQUFnQjtBQUFBLElBQ2hCLHNCQUFzQjtBQUFBLElBQ3RCLGVBQWU7QUFBQSxJQUNmLHlCQUF5QjtBQUFBLElBQ3pCLHdCQUF3QjtBQUFBLElBQ3hCLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLGlCQUFpQjtBQUFBLElBQ2pCLGlCQUFpQjtBQUFBLElBQ2pCLHNCQUFzQjtBQUFBLElBQ3RCLG1DQUFtQztBQUFBLElBQ25DLHFDQUFxQztBQUFBLElBQ3JDLHVCQUF1QjtBQUFBLElBQ3ZCLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLGVBQWU7QUFBQSxJQUNmLDJCQUEyQjtBQUFBLElBQzNCLDBCQUEwQjtBQUFBLElBQzFCLDZC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SSxZQUFZLElBQUksTUFBTSxLQUFLO0FBQzNCLGVBQVM7QUFDVDtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBQ0EsTUFBSSxRQUFRO0FBQ1IsV0FBTztBQUFBLEVBQ1g7QUFFQSxTQUFPLENBQUMsV0FBVztBQUNmLGVBQVcsUUFBUSxhQUFhO0FBQzVCLFVBQUksUUFBUSxRQUFRO0FBQ2hCLGVBQU8sSUFBSSxJQUFJLFlBQVksSUFBSSxFQUFFLE9BQU8sSUFBSSxDQUFDO0FBQUEsTUFDakQ7QUFBQSxJQUNKO0FBQ0EsV0FBTztBQUFBLEVBQ1g7QUFDSjs7O0FDNUhBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQTJDQSxJQUFNQyxRQUFPLHVCQUF1QixZQUFZLFNBQVMsRUFBRTtBQUUzRCxJQUFNLFNBQVM7QUFDZixJQUFNLGFBQWE7QUFDbkIsSUFBTSxhQUFhO0FBTVosU0FBUyxTQUFTO0FBQ3JCLFNBQU9BLE1BQUssTUFBTTtBQUN0QjtBQUtPLFNBQVMsYUFBYTtBQUN6QixTQUFPQSxNQUFLLFVBQVU7QUFDMUI7QUFNTyxTQUFTLGFBQWE7QUFDekIsU0FBT0EsTUFBSyxVQUFVO0FBQzFCOzs7QW5CM0RBLE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQWtDbEMsSUFBSSxjQUFjO0FBQ1gsU0FBUyxPQUFPO0FBQ25CLFNBQU8sT0FBTyxTQUFnQjtBQUM5QixFQUFPLE9BQU8scUJBQXFCO0FBQ25DLGdCQUFjO0FBQ2xCO0FBRUEsT0FBTyxpQkFBaUIsUUFBUSxNQUFNO0FBQ2xDLE1BQUksQ0FBQyxhQUFhO0FBQ2QsU0FBSztBQUFBLEVBQ1Q7QUFDSixDQUFDOyIsCiAgIm5hbWVzIjogWyJFcnJvciIsICJjYWxsIiwgInR5cGUiLCAiRXJyb3IiLCAiY2FsbCIsICJldmVudE5hbWUiLCAiY29udHJvbGxlciIsICJyZXNpemFibGUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiSGlkZU1ldGhvZCIsICJTaG93TWV0aG9kIiwgImNhbGwiLCAiZ2VuZXJhdGVJRCIsICJjYWxsIiwgIk1hcCIsICJNYXAiLCAia2V5IiwgImNhbGwiXQp9Cg==
