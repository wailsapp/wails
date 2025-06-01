var __defProp = Object.defineProperty;
var __defProps = Object.defineProperties;
var __getOwnPropDescs = Object.getOwnPropertyDescriptors;
var __getOwnPropSymbols = Object.getOwnPropertySymbols;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __propIsEnum = Object.prototype.propertyIsEnumerable;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __spreadValues = (a, b) => {
  for (var prop in b || (b = {}))
    if (__hasOwnProp.call(b, prop))
      __defNormalProp(a, prop, b[prop]);
  if (__getOwnPropSymbols)
    for (var prop of __getOwnPropSymbols(b)) {
      if (__propIsEnum.call(b, prop))
        __defNormalProp(a, prop, b[prop]);
    }
  return a;
};
var __spreadProps = (a, b) => __defProps(a, __getOwnPropDescs(b));
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};

// desktop/@wailsio/runtime/src/index.ts
var index_exports = {};
__export(index_exports, {
  Application: () => application_exports,
  Browser: () => browser_exports,
  Call: () => calls_exports,
  CancelError: () => CancelError,
  CancellablePromise: () => CancellablePromise,
  CancelledRejectionError: () => CancelledRejectionError,
  Clipboard: () => clipboard_exports,
  Create: () => create_exports,
  Dialogs: () => dialogs_exports,
  Events: () => events_exports,
  Flags: () => flags_exports,
  Screens: () => screens_exports,
  System: () => system_exports,
  WML: () => wml_exports,
  Window: () => window_default
});

// desktop/@wailsio/runtime/src/wml.ts
var wml_exports = {};
__export(wml_exports, {
  Enable: () => Enable,
  Reload: () => Reload
});

// desktop/@wailsio/runtime/src/browser.ts
var browser_exports = {};
__export(browser_exports, {
  OpenURL: () => OpenURL
});

// desktop/@wailsio/runtime/src/nanoid.ts
var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
function nanoid(size = 21) {
  let id = "";
  let i = size | 0;
  while (i--) {
    id += urlAlphabet[Math.random() * 64 | 0];
  }
  return id;
}

// desktop/@wailsio/runtime/src/runtime.ts
var runtimeURL = window.location.origin + "/wails/runtime";
var objectNames = Object.freeze({
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
});
var clientId = nanoid();
function newRuntimeCaller(object, windowName = "") {
  return function(method, args = null) {
    return runtimeCallWithID(object, method, windowName, args);
  };
}
async function runtimeCallWithID(objectID, method, windowName, args) {
  var _a2, _b;
  let url = new URL(runtimeURL);
  url.searchParams.append("object", objectID.toString());
  url.searchParams.append("method", method.toString());
  if (args) {
    url.searchParams.append("args", JSON.stringify(args));
  }
  let headers = {
    ["x-wails-client-id"]: clientId
  };
  if (windowName) {
    headers["x-wails-window-name"] = windowName;
  }
  let response = await fetch(url, { headers });
  if (!response.ok) {
    throw new Error(await response.text());
  }
  if (((_b = (_a2 = response.headers.get("Content-Type")) == null ? void 0 : _a2.indexOf("application/json")) != null ? _b : -1) !== -1) {
    return response.json();
  } else {
    return response.text();
  }
}

// desktop/@wailsio/runtime/src/browser.ts
var call = newRuntimeCaller(objectNames.Browser);
var BrowserOpenURL = 0;
function OpenURL(url) {
  return call(BrowserOpenURL, { url: url.toString() });
}

// desktop/@wailsio/runtime/src/dialogs.ts
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
var call2 = newRuntimeCaller(objectNames.Dialog);
var dialogResponses = /* @__PURE__ */ new Map();
var DialogInfo = 0;
var DialogWarning = 1;
var DialogError = 2;
var DialogQuestion = 3;
var DialogOpenFile = 4;
var DialogSaveFile = 5;
function dialogResultCallback(id, data, isJSON) {
  let resolvers = getAndDeleteResponse(id);
  if (!resolvers) {
    return;
  }
  if (isJSON) {
    try {
      resolvers.resolve(JSON.parse(data));
    } catch (err) {
      resolvers.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
    }
  } else {
    resolvers.resolve(data);
  }
}
function dialogErrorCallback(id, message) {
  var _a2;
  (_a2 = getAndDeleteResponse(id)) == null ? void 0 : _a2.reject(new window.Error(message));
}
function getAndDeleteResponse(id) {
  const response = dialogResponses.get(id);
  dialogResponses.delete(id);
  return response;
}
function generateID() {
  let result;
  do {
    result = nanoid();
  } while (dialogResponses.has(result));
  return result;
}
function dialog(type, options = {}) {
  const id = generateID();
  return new Promise((resolve, reject) => {
    dialogResponses.set(id, { resolve, reject });
    call2(type, Object.assign({ "dialog-id": id }, options)).catch((err) => {
      dialogResponses.delete(id);
      reject(err);
    });
  });
}
function Info(options) {
  return dialog(DialogInfo, options);
}
function Warning(options) {
  return dialog(DialogWarning, options);
}
function Error2(options) {
  return dialog(DialogError, options);
}
function Question(options) {
  return dialog(DialogQuestion, options);
}
function OpenFile(options) {
  var _a2;
  return (_a2 = dialog(DialogOpenFile, options)) != null ? _a2 : [];
}
function SaveFile(options) {
  return dialog(DialogSaveFile, options);
}

// desktop/@wailsio/runtime/src/events.ts
var events_exports = {};
__export(events_exports, {
  Emit: () => Emit,
  Off: () => Off,
  OffAll: () => OffAll,
  On: () => On,
  OnMultiple: () => OnMultiple,
  Once: () => Once,
  Types: () => Types,
  WailsEvent: () => WailsEvent
});

// desktop/@wailsio/runtime/src/listener.ts
var eventListeners = /* @__PURE__ */ new Map();
var Listener = class {
  constructor(eventName, callback, maxCallbacks) {
    this.eventName = eventName;
    this.callback = callback;
    this.maxCallbacks = maxCallbacks || -1;
  }
  dispatch(data) {
    try {
      this.callback(data);
    } catch (err) {
      console.error(err);
    }
    if (this.maxCallbacks === -1) return false;
    this.maxCallbacks -= 1;
    return this.maxCallbacks === 0;
  }
};
function listenerOff(listener) {
  let listeners = eventListeners.get(listener.eventName);
  if (!listeners) {
    return;
  }
  listeners = listeners.filter((l) => l !== listener);
  if (listeners.length === 0) {
    eventListeners.delete(listener.eventName);
  } else {
    eventListeners.set(listener.eventName, listeners);
  }
}

// desktop/@wailsio/runtime/src/event_types.ts
var Types = Object.freeze({
  Windows: Object.freeze({
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
  }),
  Mac: Object.freeze({
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
  }),
  Linux: Object.freeze({
    ApplicationStartup: "linux:ApplicationStartup",
    SystemThemeChanged: "linux:SystemThemeChanged",
    WindowDeleteEvent: "linux:WindowDeleteEvent",
    WindowDidMove: "linux:WindowDidMove",
    WindowDidResize: "linux:WindowDidResize",
    WindowFocusIn: "linux:WindowFocusIn",
    WindowFocusOut: "linux:WindowFocusOut",
    WindowLoadChanged: "linux:WindowLoadChanged"
  }),
  Common: Object.freeze({
    ApplicationOpenedWithFile: "ApplicationOpenedWithFile",
    ApplicationStarted: "ApplicationStarted",
    ThemeChanged: "ThemeChanged",
    WindowClosing: "WindowClosing",
    WindowDidMove: "WindowDidMove",
    WindowDidResize: "WindowDidResize",
    WindowDPIChanged: "WindowDPIChanged",
    WindowFilesDropped: "WindowFilesDropped",
    WindowFocus: "WindowFocus",
    WindowFullscreen: "WindowFullscreen",
    WindowHide: "WindowHide",
    WindowLostFocus: "WindowLostFocus",
    WindowMaximise: "WindowMaximise",
    WindowMinimise: "WindowMinimise",
    WindowRestore: "WindowRestore",
    WindowRuntimeReady: "WindowRuntimeReady",
    WindowShow: "WindowShow",
    WindowUnFullscreen: "WindowUnFullscreen",
    WindowUnMaximise: "WindowUnMaximise",
    WindowUnMinimise: "WindowUnMinimise",
    WindowZoom: "WindowZoom",
    WindowZoomIn: "WindowZoomIn",
    WindowZoomOut: "WindowZoomOut",
    WindowZoomReset: "WindowZoomReset"
  })
});

// desktop/@wailsio/runtime/src/events.ts
window._wails = window._wails || {};
window._wails.dispatchWailsEvent = dispatchWailsEvent;
var call3 = newRuntimeCaller(objectNames.Events);
var EmitMethod = 0;
var WailsEvent = class {
  constructor(name, data = null) {
    this.name = name;
    this.data = data;
  }
};
function dispatchWailsEvent(event) {
  let listeners = eventListeners.get(event.name);
  if (!listeners) {
    return;
  }
  let wailsEvent = new WailsEvent(event.name, event.data);
  if ("sender" in event) {
    wailsEvent.sender = event.sender;
  }
  listeners = listeners.filter((listener) => !listener.dispatch(wailsEvent));
  if (listeners.length === 0) {
    eventListeners.delete(event.name);
  } else {
    eventListeners.set(event.name, listeners);
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
function Off(...eventNames) {
  eventNames.forEach((eventName) => eventListeners.delete(eventName));
}
function OffAll() {
  eventListeners.clear();
}
function Emit(name, data) {
  let event;
  if (typeof name === "object" && name !== null && "name" in name && "data" in name) {
    event = new WailsEvent(name["name"], name["data"]);
  } else {
    event = new WailsEvent(name, data);
  }
  return call3(EmitMethod, event);
}

// desktop/@wailsio/runtime/src/utils.ts
function debugLog(message) {
  console.log(
    "%c wails3 %c " + message + " ",
    "background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem",
    "background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem"
  );
}
function canTrackButtons() {
  return new MouseEvent("mousedown").buttons === 0;
}
function canAbortListeners() {
  if (!EventTarget || !AbortSignal || !AbortController)
    return false;
  let result = true;
  const target = new EventTarget();
  const controller = new AbortController();
  target.addEventListener("test", () => {
    result = false;
  }, { signal: controller.signal });
  controller.abort();
  target.dispatchEvent(new CustomEvent("test"));
  return result;
}
function eventTarget(event) {
  var _a2;
  if (event.target instanceof HTMLElement) {
    return event.target;
  } else if (!(event.target instanceof HTMLElement) && event.target instanceof Node) {
    return (_a2 = event.target.parentElement) != null ? _a2 : document.body;
  } else {
    return document.body;
  }
}
var isReady = false;
document.addEventListener("DOMContentLoaded", () => {
  isReady = true;
});
function whenReady(callback) {
  if (isReady || document.readyState === "complete") {
    callback();
  } else {
    document.addEventListener("DOMContentLoaded", callback);
  }
}

// desktop/@wailsio/runtime/src/window.ts
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
var callerSym = Symbol("caller");
callerSym;
var _Window = class _Window {
  /**
   * Initialises a window object with the specified name.
   *
   * @private
   * @param name - The name of the target window.
   */
  constructor(name = "") {
    this[callerSym] = newRuntimeCaller(objectNames.Window, name);
    for (const method of Object.getOwnPropertyNames(_Window.prototype)) {
      if (method !== "constructor" && typeof this[method] === "function") {
        this[method] = this[method].bind(this);
      }
    }
  }
  /**
   * Gets the specified window.
   *
   * @param name - The name of the window to get.
   * @returns The corresponding window object.
   */
  Get(name) {
    return new _Window(name);
  }
  /**
   * Returns the absolute position of the window.
   *
   * @returns The current absolute position of the window.
   */
  Position() {
    return this[callerSym](PositionMethod);
  }
  /**
   * Centers the window on the screen.
   */
  Center() {
    return this[callerSym](CenterMethod);
  }
  /**
   * Closes the window.
   */
  Close() {
    return this[callerSym](CloseMethod);
  }
  /**
   * Disables min/max size constraints.
   */
  DisableSizeConstraints() {
    return this[callerSym](DisableSizeConstraintsMethod);
  }
  /**
   * Enables min/max size constraints.
   */
  EnableSizeConstraints() {
    return this[callerSym](EnableSizeConstraintsMethod);
  }
  /**
   * Focuses the window.
   */
  Focus() {
    return this[callerSym](FocusMethod);
  }
  /**
   * Forces the window to reload the page assets.
   */
  ForceReload() {
    return this[callerSym](ForceReloadMethod);
  }
  /**
   * Switches the window to fullscreen mode.
   */
  Fullscreen() {
    return this[callerSym](FullscreenMethod);
  }
  /**
   * Returns the screen that the window is on.
   *
   * @returns The screen the window is currently on.
   */
  GetScreen() {
    return this[callerSym](GetScreenMethod);
  }
  /**
   * Returns the current zoom level of the window.
   *
   * @returns The current zoom level.
   */
  GetZoom() {
    return this[callerSym](GetZoomMethod);
  }
  /**
   * Returns the height of the window.
   *
   * @returns The current height of the window.
   */
  Height() {
    return this[callerSym](HeightMethod);
  }
  /**
   * Hides the window.
   */
  Hide() {
    return this[callerSym](HideMethod);
  }
  /**
   * Returns true if the window is focused.
   *
   * @returns Whether the window is currently focused.
   */
  IsFocused() {
    return this[callerSym](IsFocusedMethod);
  }
  /**
   * Returns true if the window is fullscreen.
   *
   * @returns Whether the window is currently fullscreen.
   */
  IsFullscreen() {
    return this[callerSym](IsFullscreenMethod);
  }
  /**
   * Returns true if the window is maximised.
   *
   * @returns Whether the window is currently maximised.
   */
  IsMaximised() {
    return this[callerSym](IsMaximisedMethod);
  }
  /**
   * Returns true if the window is minimised.
   *
   * @returns Whether the window is currently minimised.
   */
  IsMinimised() {
    return this[callerSym](IsMinimisedMethod);
  }
  /**
   * Maximises the window.
   */
  Maximise() {
    return this[callerSym](MaximiseMethod);
  }
  /**
   * Minimises the window.
   */
  Minimise() {
    return this[callerSym](MinimiseMethod);
  }
  /**
   * Returns the name of the window.
   *
   * @returns The name of the window.
   */
  Name() {
    return this[callerSym](NameMethod);
  }
  /**
   * Opens the development tools pane.
   */
  OpenDevTools() {
    return this[callerSym](OpenDevToolsMethod);
  }
  /**
   * Returns the relative position of the window to the screen.
   *
   * @returns The current relative position of the window.
   */
  RelativePosition() {
    return this[callerSym](RelativePositionMethod);
  }
  /**
   * Reloads the page assets.
   */
  Reload() {
    return this[callerSym](ReloadMethod);
  }
  /**
   * Returns true if the window is resizable.
   *
   * @returns Whether the window is currently resizable.
   */
  Resizable() {
    return this[callerSym](ResizableMethod);
  }
  /**
   * Restores the window to its previous state if it was previously minimised, maximised or fullscreen.
   */
  Restore() {
    return this[callerSym](RestoreMethod);
  }
  /**
   * Sets the absolute position of the window.
   *
   * @param x - The desired horizontal absolute position of the window.
   * @param y - The desired vertical absolute position of the window.
   */
  SetPosition(x, y) {
    return this[callerSym](SetPositionMethod, { x, y });
  }
  /**
   * Sets the window to be always on top.
   *
   * @param alwaysOnTop - Whether the window should stay on top.
   */
  SetAlwaysOnTop(alwaysOnTop) {
    return this[callerSym](SetAlwaysOnTopMethod, { alwaysOnTop });
  }
  /**
   * Sets the background colour of the window.
   *
   * @param r - The desired red component of the window background.
   * @param g - The desired green component of the window background.
   * @param b - The desired blue component of the window background.
   * @param a - The desired alpha component of the window background.
   */
  SetBackgroundColour(r, g, b, a) {
    return this[callerSym](SetBackgroundColourMethod, { r, g, b, a });
  }
  /**
   * Removes the window frame and title bar.
   *
   * @param frameless - Whether the window should be frameless.
   */
  SetFrameless(frameless) {
    return this[callerSym](SetFramelessMethod, { frameless });
  }
  /**
   * Disables the system fullscreen button.
   *
   * @param enabled - Whether the fullscreen button should be enabled.
   */
  SetFullscreenButtonEnabled(enabled) {
    return this[callerSym](SetFullscreenButtonEnabledMethod, { enabled });
  }
  /**
   * Sets the maximum size of the window.
   *
   * @param width - The desired maximum width of the window.
   * @param height - The desired maximum height of the window.
   */
  SetMaxSize(width, height) {
    return this[callerSym](SetMaxSizeMethod, { width, height });
  }
  /**
   * Sets the minimum size of the window.
   *
   * @param width - The desired minimum width of the window.
   * @param height - The desired minimum height of the window.
   */
  SetMinSize(width, height) {
    return this[callerSym](SetMinSizeMethod, { width, height });
  }
  /**
   * Sets the relative position of the window to the screen.
   *
   * @param x - The desired horizontal relative position of the window.
   * @param y - The desired vertical relative position of the window.
   */
  SetRelativePosition(x, y) {
    return this[callerSym](SetRelativePositionMethod, { x, y });
  }
  /**
   * Sets whether the window is resizable.
   *
   * @param resizable - Whether the window should be resizable.
   */
  SetResizable(resizable2) {
    return this[callerSym](SetResizableMethod, { resizable: resizable2 });
  }
  /**
   * Sets the size of the window.
   *
   * @param width - The desired width of the window.
   * @param height - The desired height of the window.
   */
  SetSize(width, height) {
    return this[callerSym](SetSizeMethod, { width, height });
  }
  /**
   * Sets the title of the window.
   *
   * @param title - The desired title of the window.
   */
  SetTitle(title) {
    return this[callerSym](SetTitleMethod, { title });
  }
  /**
   * Sets the zoom level of the window.
   *
   * @param zoom - The desired zoom level.
   */
  SetZoom(zoom) {
    return this[callerSym](SetZoomMethod, { zoom });
  }
  /**
   * Shows the window.
   */
  Show() {
    return this[callerSym](ShowMethod);
  }
  /**
   * Returns the size of the window.
   *
   * @returns The current size of the window.
   */
  Size() {
    return this[callerSym](SizeMethod);
  }
  /**
   * Toggles the window between fullscreen and normal.
   */
  ToggleFullscreen() {
    return this[callerSym](ToggleFullscreenMethod);
  }
  /**
   * Toggles the window between maximised and normal.
   */
  ToggleMaximise() {
    return this[callerSym](ToggleMaximiseMethod);
  }
  /**
   * Un-fullscreens the window.
   */
  UnFullscreen() {
    return this[callerSym](UnFullscreenMethod);
  }
  /**
   * Un-maximises the window.
   */
  UnMaximise() {
    return this[callerSym](UnMaximiseMethod);
  }
  /**
   * Un-minimises the window.
   */
  UnMinimise() {
    return this[callerSym](UnMinimiseMethod);
  }
  /**
   * Returns the width of the window.
   *
   * @returns The current width of the window.
   */
  Width() {
    return this[callerSym](WidthMethod);
  }
  /**
   * Zooms the window.
   */
  Zoom() {
    return this[callerSym](ZoomMethod);
  }
  /**
   * Increases the zoom level of the webview content.
   */
  ZoomIn() {
    return this[callerSym](ZoomInMethod);
  }
  /**
   * Decreases the zoom level of the webview content.
   */
  ZoomOut() {
    return this[callerSym](ZoomOutMethod);
  }
  /**
   * Resets the zoom level of the webview content.
   */
  ZoomReset() {
    return this[callerSym](ZoomResetMethod);
  }
};
var Window = _Window;
var thisWindow = new Window("");
var window_default = thisWindow;

// desktop/@wailsio/runtime/src/wml.ts
function sendEvent(eventName, data = null) {
  void Emit(eventName, data);
}
function callWindowMethod(windowName, methodName) {
  const targetWindow = window_default.Get(windowName);
  const method = targetWindow[methodName];
  if (typeof method !== "function") {
    console.error("Window method '".concat(methodName, "' not found"));
    return;
  }
  try {
    method.call(targetWindow);
  } catch (e) {
    console.error("Error calling window method '".concat(methodName, "': "), e);
  }
}
function onWMLTriggered(ev) {
  const element = ev.currentTarget;
  function runEffect(choice = "Yes") {
    if (choice !== "Yes")
      return;
    const eventType = element.getAttribute("wml-event") || element.getAttribute("data-wml-event");
    const targetWindow = element.getAttribute("wml-target-window") || element.getAttribute("data-wml-target-window") || "";
    const windowMethod = element.getAttribute("wml-window") || element.getAttribute("data-wml-window");
    const url = element.getAttribute("wml-openurl") || element.getAttribute("data-wml-openurl");
    if (eventType !== null)
      sendEvent(eventType);
    if (windowMethod !== null)
      callWindowMethod(targetWindow, windowMethod);
    if (url !== null)
      void OpenURL(url);
  }
  const confirm = element.getAttribute("wml-confirm") || element.getAttribute("data-wml-confirm");
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
var controllerSym = Symbol("controller");
var triggerMapSym = Symbol("triggerMap");
var elementCountSym = Symbol("elementCount");
controllerSym;
var AbortControllerRegistry = class {
  constructor() {
    this[controllerSym] = new AbortController();
  }
  /**
   * Returns an options object for addEventListener that ties the listener
   * to the AbortSignal from the current AbortController.
   *
   * @param element - An HTML element
   * @param triggers - The list of active WML trigger events for the specified elements
   */
  set(element, triggers) {
    return { signal: this[controllerSym].signal };
  }
  /**
   * Removes all registered event listeners and resets the registry.
   */
  reset() {
    this[controllerSym].abort();
    this[controllerSym] = new AbortController();
  }
};
triggerMapSym, elementCountSym;
var WeakMapRegistry = class {
  constructor() {
    this[triggerMapSym] = /* @__PURE__ */ new WeakMap();
    this[elementCountSym] = 0;
  }
  /**
   * Sets active triggers for the specified element.
   *
   * @param element - An HTML element
   * @param triggers - The list of active WML trigger events for the specified element
   */
  set(element, triggers) {
    if (!this[triggerMapSym].has(element)) {
      this[elementCountSym]++;
    }
    this[triggerMapSym].set(element, triggers);
    return {};
  }
  /**
   * Removes all registered event listeners.
   */
  reset() {
    if (this[elementCountSym] <= 0)
      return;
    for (const element of document.body.querySelectorAll("*")) {
      if (this[elementCountSym] <= 0)
        break;
      const triggers = this[triggerMapSym].get(element);
      if (triggers != null) {
        this[elementCountSym]--;
      }
      for (const trigger of triggers || [])
        element.removeEventListener(trigger, onWMLTriggered);
    }
    this[triggerMapSym] = /* @__PURE__ */ new WeakMap();
    this[elementCountSym] = 0;
  }
};
var triggerRegistry = canAbortListeners() ? new AbortControllerRegistry() : new WeakMapRegistry();
function addWMLListeners(element) {
  const triggerRegExp = /\S+/g;
  const triggerAttr = element.getAttribute("wml-trigger") || element.getAttribute("data-wml-trigger") || "click";
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
  document.body.querySelectorAll("[wml-event], [wml-window], [wml-openurl], [data-wml-event], [data-wml-window], [data-wml-openurl]").forEach(addWMLListeners);
}

// desktop/compiled/main.js
window.wails = index_exports;
Enable();
if (true) {
  debugLog("Wails Runtime Loaded");
}

// desktop/@wailsio/runtime/src/system.ts
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
var call4 = newRuntimeCaller(objectNames.System);
var SystemIsDarkMode = 0;
var SystemEnvironment = 1;
var _invoke = function() {
  var _a2, _b, _c, _d, _e;
  try {
    if ((_b = (_a2 = window.chrome) == null ? void 0 : _a2.webview) == null ? void 0 : _b.postMessage) {
      return window.chrome.webview.postMessage.bind(window.chrome.webview);
    } else if ((_e = (_d = (_c = window.webkit) == null ? void 0 : _c.messageHandlers) == null ? void 0 : _d["external"]) == null ? void 0 : _e.postMessage) {
      return window.webkit.messageHandlers["external"].postMessage.bind(window.webkit.messageHandlers["external"]);
    }
  } catch (e) {
  }
  console.warn(
    "\n%c\u26A0\uFE0F Browser Environment Detected %c\n\n%cOnly UI previews are available in the browser. For full functionality, please run the application in desktop mode.\nMore information at: https://v3.wails.io/learn/build/#using-a-browser-for-development\n",
    "background: #ffffff; color: #000000; font-weight: bold; padding: 4px 8px; border-radius: 4px; border: 2px solid #000000;",
    "background: transparent;",
    "color: #ffffff; font-style: italic; font-weight: bold;"
  );
  return null;
}();
function invoke(msg) {
  _invoke == null ? void 0 : _invoke(msg);
}
function IsDarkMode() {
  return call4(SystemIsDarkMode);
}
async function Capabilities() {
  let response = await fetch("/wails/capabilities");
  if (response.ok) {
    return response.json();
  } else {
    throw new Error("could not fetch capabilities: " + response.statusText);
  }
}
function Environment() {
  return call4(SystemEnvironment);
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
  return Boolean(window._wails.environment.Debug);
}

// desktop/@wailsio/runtime/src/contextmenu.ts
window.addEventListener("contextmenu", contextMenuHandler);
var call5 = newRuntimeCaller(objectNames.ContextMenu);
var ContextMenuOpen = 0;
function openContextMenu(id, x, y, data) {
  void call5(ContextMenuOpen, { id, x, y, data });
}
function contextMenuHandler(event) {
  const target = eventTarget(event);
  const customContextMenu = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu").trim();
  if (customContextMenu) {
    event.preventDefault();
    const data = window.getComputedStyle(target).getPropertyValue("--custom-contextmenu-data");
    openContextMenu(customContextMenu, event.clientX, event.clientY, data);
  } else {
    processDefaultContextMenu(event, target);
  }
}
function processDefaultContextMenu(event, target) {
  if (IsDebug()) {
    return;
  }
  switch (window.getComputedStyle(target).getPropertyValue("--default-contextmenu").trim()) {
    case "show":
      return;
    case "hide":
      event.preventDefault();
      return;
  }
  if (target.isContentEditable) {
    return;
  }
  const selection = window.getSelection();
  const hasSelection = selection && selection.toString().length > 0;
  if (hasSelection) {
    for (let i = 0; i < selection.rangeCount; i++) {
      const range = selection.getRangeAt(i);
      const rects = range.getClientRects();
      for (let j = 0; j < rects.length; j++) {
        const rect = rects[j];
        if (document.elementFromPoint(rect.left, rect.top) === target) {
          return;
        }
      }
    }
  }
  if (target instanceof HTMLInputElement || target instanceof HTMLTextAreaElement) {
    if (hasSelection || !target.readOnly && !target.disabled) {
      return;
    }
  }
  event.preventDefault();
}

// desktop/@wailsio/runtime/src/flags.ts
var flags_exports = {};
__export(flags_exports, {
  GetFlag: () => GetFlag
});
function GetFlag(key) {
  try {
    return window._wails.flags[key];
  } catch (e) {
    throw new Error("Unable to retrieve flag '" + key + "': " + e, { cause: e });
  }
}

// desktop/@wailsio/runtime/src/drag.ts
var canDrag = false;
var dragging = false;
var resizable = false;
var canResize = false;
var resizing = false;
var resizeEdge = "";
var defaultCursor = "auto";
var buttons = 0;
var buttonsTracked = canTrackButtons();
window._wails = window._wails || {};
window._wails.setResizable = (value) => {
  resizable = value;
  if (!resizable) {
    canResize = resizing = false;
    setResize();
  }
};
window.addEventListener("mousedown", update, { capture: true });
window.addEventListener("mousemove", update, { capture: true });
window.addEventListener("mouseup", update, { capture: true });
for (const ev of ["click", "contextmenu", "dblclick"]) {
  window.addEventListener(ev, suppressEvent, { capture: true });
}
function suppressEvent(event) {
  if (dragging || resizing) {
    event.stopImmediatePropagation();
    event.stopPropagation();
    event.preventDefault();
  }
}
var MouseDown = 0;
var MouseUp = 1;
var MouseMove = 2;
function update(event) {
  let eventType, eventButtons = event.buttons;
  switch (event.type) {
    case "mousedown":
      eventType = MouseDown;
      if (!buttonsTracked) {
        eventButtons = buttons | 1 << event.button;
      }
      break;
    case "mouseup":
      eventType = MouseUp;
      if (!buttonsTracked) {
        eventButtons = buttons & ~(1 << event.button);
      }
      break;
    default:
      eventType = MouseMove;
      if (!buttonsTracked) {
        eventButtons = buttons;
      }
      break;
  }
  let released = buttons & ~eventButtons;
  let pressed = eventButtons & ~buttons;
  buttons = eventButtons;
  if (eventType === MouseDown && !(pressed & event.button)) {
    released |= 1 << event.button;
    pressed |= 1 << event.button;
  }
  if (eventType !== MouseMove && resizing || dragging && (eventType === MouseDown || event.button !== 0)) {
    event.stopImmediatePropagation();
    event.stopPropagation();
    event.preventDefault();
  }
  if (released & 1) {
    primaryUp(event);
  }
  if (pressed & 1) {
    primaryDown(event);
  }
  if (eventType === MouseMove) {
    onMouseMove(event);
  }
  ;
}
function primaryDown(event) {
  canDrag = false;
  canResize = false;
  if (!IsWindows()) {
    if (event.type === "mousedown" && event.button === 0 && event.detail !== 1) {
      return;
    }
  }
  if (resizeEdge) {
    canResize = true;
    return;
  }
  const target = eventTarget(event);
  const style = window.getComputedStyle(target);
  canDrag = style.getPropertyValue("--wails-draggable").trim() === "drag" && (event.offsetX - parseFloat(style.paddingLeft) < target.clientWidth && event.offsetY - parseFloat(style.paddingTop) < target.clientHeight);
}
function primaryUp(event) {
  canDrag = false;
  dragging = false;
  canResize = false;
  resizing = false;
}
var cursorForEdge = Object.freeze({
  "se-resize": "nwse-resize",
  "sw-resize": "nesw-resize",
  "nw-resize": "nwse-resize",
  "ne-resize": "nesw-resize",
  "w-resize": "ew-resize",
  "n-resize": "ns-resize",
  "s-resize": "ns-resize",
  "e-resize": "ew-resize"
});
function setResize(edge) {
  if (edge) {
    if (!resizeEdge) {
      defaultCursor = document.body.style.cursor;
    }
    document.body.style.cursor = cursorForEdge[edge];
  } else if (!edge && resizeEdge) {
    document.body.style.cursor = defaultCursor;
  }
  resizeEdge = edge || "";
}
function onMouseMove(event) {
  if (canResize && resizeEdge) {
    resizing = true;
    invoke("wails:resize:" + resizeEdge);
  } else if (canDrag) {
    dragging = true;
    invoke("wails:drag");
  }
  if (dragging || resizing) {
    canDrag = canResize = false;
    return;
  }
  if (!resizable || !IsWindows()) {
    if (resizeEdge) {
      setResize();
    }
    return;
  }
  const resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
  const resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
  const cornerExtra = GetFlag("resizeCornerExtra") || 10;
  const rightBorder = window.outerWidth - event.clientX < resizeHandleWidth;
  const leftBorder = event.clientX < resizeHandleWidth;
  const topBorder = event.clientY < resizeHandleHeight;
  const bottomBorder = window.outerHeight - event.clientY < resizeHandleHeight;
  const rightCorner = window.outerWidth - event.clientX < resizeHandleWidth + cornerExtra;
  const leftCorner = event.clientX < resizeHandleWidth + cornerExtra;
  const topCorner = event.clientY < resizeHandleHeight + cornerExtra;
  const bottomCorner = window.outerHeight - event.clientY < resizeHandleHeight + cornerExtra;
  if (!leftCorner && !topCorner && !bottomCorner && !rightCorner) {
    setResize();
  } else if (rightCorner && bottomCorner) setResize("se-resize");
  else if (leftCorner && bottomCorner) setResize("sw-resize");
  else if (leftCorner && topCorner) setResize("nw-resize");
  else if (topCorner && rightCorner) setResize("ne-resize");
  else if (leftBorder) setResize("w-resize");
  else if (topBorder) setResize("n-resize");
  else if (bottomBorder) setResize("s-resize");
  else if (rightBorder) setResize("e-resize");
  else setResize();
}

// desktop/@wailsio/runtime/src/application.ts
var application_exports = {};
__export(application_exports, {
  Hide: () => Hide,
  Quit: () => Quit,
  Show: () => Show
});
var call6 = newRuntimeCaller(objectNames.Application);
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

// desktop/@wailsio/runtime/src/calls.ts
var calls_exports = {};
__export(calls_exports, {
  ByID: () => ByID,
  ByName: () => ByName,
  Call: () => Call,
  RuntimeError: () => RuntimeError
});

// desktop/@wailsio/runtime/src/callable.ts
var fnToStr = Function.prototype.toString;
var reflectApply = typeof Reflect === "object" && Reflect !== null && Reflect.apply;
var badArrayLike;
var isCallableMarker;
if (typeof reflectApply === "function" && typeof Object.defineProperty === "function") {
  try {
    badArrayLike = Object.defineProperty({}, "length", {
      get: function() {
        throw isCallableMarker;
      }
    });
    isCallableMarker = {};
    reflectApply(function() {
      throw 42;
    }, null, badArrayLike);
  } catch (_) {
    if (_ !== isCallableMarker) {
      reflectApply = null;
    }
  }
} else {
  reflectApply = null;
}
var constructorRegex = /^\s*class\b/;
var isES6ClassFn = function isES6ClassFunction(value) {
  try {
    var fnStr = fnToStr.call(value);
    return constructorRegex.test(fnStr);
  } catch (e) {
    return false;
  }
};
var tryFunctionObject = function tryFunctionToStr(value) {
  try {
    if (isES6ClassFn(value)) {
      return false;
    }
    fnToStr.call(value);
    return true;
  } catch (e) {
    return false;
  }
};
var toStr = Object.prototype.toString;
var objectClass = "[object Object]";
var fnClass = "[object Function]";
var genClass = "[object GeneratorFunction]";
var ddaClass = "[object HTMLAllCollection]";
var ddaClass2 = "[object HTML document.all class]";
var ddaClass3 = "[object HTMLCollection]";
var hasToStringTag = typeof Symbol === "function" && !!Symbol.toStringTag;
var isIE68 = !(0 in [,]);
var isDDA = function isDocumentDotAll() {
  return false;
};
if (typeof document === "object") {
  all = document.all;
  if (toStr.call(all) === toStr.call(document.all)) {
    isDDA = function isDocumentDotAll2(value) {
      if ((isIE68 || !value) && (typeof value === "undefined" || typeof value === "object")) {
        try {
          var str = toStr.call(value);
          return (str === ddaClass || str === ddaClass2 || str === ddaClass3 || str === objectClass) && value("") == null;
        } catch (e) {
        }
      }
      return false;
    };
  }
}
var all;
function isCallableRefApply(value) {
  if (isDDA(value)) {
    return true;
  }
  if (!value) {
    return false;
  }
  if (typeof value !== "function" && typeof value !== "object") {
    return false;
  }
  try {
    reflectApply(value, null, badArrayLike);
  } catch (e) {
    if (e !== isCallableMarker) {
      return false;
    }
  }
  return !isES6ClassFn(value) && tryFunctionObject(value);
}
function isCallableNoRefApply(value) {
  if (isDDA(value)) {
    return true;
  }
  if (!value) {
    return false;
  }
  if (typeof value !== "function" && typeof value !== "object") {
    return false;
  }
  if (hasToStringTag) {
    return tryFunctionObject(value);
  }
  if (isES6ClassFn(value)) {
    return false;
  }
  var strClass = toStr.call(value);
  if (strClass !== fnClass && strClass !== genClass && !/^\[object HTML/.test(strClass)) {
    return false;
  }
  return tryFunctionObject(value);
}
var callable_default = reflectApply ? isCallableRefApply : isCallableNoRefApply;

// desktop/@wailsio/runtime/src/cancellable.ts
var CancelError = class extends Error {
  /**
   * Constructs a new `CancelError` instance.
   * @param message - The error message.
   * @param options - Options to be forwarded to the Error constructor.
   */
  constructor(message, options) {
    super(message, options);
    this.name = "CancelError";
  }
};
var CancelledRejectionError = class extends Error {
  /**
   * Constructs a new `CancelledRejectionError` instance.
   * @param promise - The promise that caused the error originally.
   * @param reason - The rejection reason.
   * @param info - An optional informative message specifying the circumstances in which the error was thrown.
   *               Defaults to the string `"Unhandled rejection in cancelled promise."`.
   */
  constructor(promise, reason, info) {
    super((info != null ? info : "Unhandled rejection in cancelled promise.") + " Reason: " + errorMessage(reason), { cause: reason });
    this.promise = promise;
    this.name = "CancelledRejectionError";
  }
};
var barrierSym = Symbol("barrier");
var cancelImplSym = Symbol("cancelImpl");
var _a;
var species = (_a = Symbol.species) != null ? _a : Symbol("speciesPolyfill");
var CancellablePromise = class _CancellablePromise extends Promise {
  /**
   * Creates a new `CancellablePromise`.
   *
   * @param executor - A callback used to initialize the promise. This callback is passed two arguments:
   *                   a `resolve` callback used to resolve the promise with a value
   *                   or the result of another promise (possibly cancellable),
   *                   and a `reject` callback used to reject the promise with a provided reason or error.
   *                   If the value provided to the `resolve` callback is a thenable _and_ cancellable object
   *                   (it has a `then` _and_ a `cancel` method),
   *                   cancellation requests will be forwarded to that object and the oncancelled will not be invoked anymore.
   *                   If any one of the two callbacks is called _after_ the promise has been cancelled,
   *                   the provided values will be cancelled and resolved as usual,
   *                   but their results will be discarded.
   *                   However, if the resolution process ultimately ends up in a rejection
   *                   that is not due to cancellation, the rejection reason
   *                   will be wrapped in a {@link CancelledRejectionError}
   *                   and bubbled up as an unhandled rejection.
   * @param oncancelled - It is the caller's responsibility to ensure that any operation
   *                      started by the executor is properly halted upon cancellation.
   *                      This optional callback can be used to that purpose.
   *                      It will be called _synchronously_ with a cancellation cause
   *                      when cancellation is requested, _after_ the promise has already rejected
   *                      with a {@link CancelError}, but _before_
   *                      any {@link then}/{@link catch}/{@link finally} callback runs.
   *                      If the callback returns a thenable, the promise returned from {@link cancel}
   *                      will only fulfill after the former has settled.
   *                      Unhandled exceptions or rejections from the callback will be wrapped
   *                      in a {@link CancelledRejectionError} and bubbled up as unhandled rejections.
   *                      If the `resolve` callback is called before cancellation with a cancellable promise,
   *                      cancellation requests on this promise will be diverted to that promise,
   *                      and the original `oncancelled` callback will be discarded.
   */
  constructor(executor, oncancelled) {
    let resolve;
    let reject;
    super((res, rej) => {
      resolve = res;
      reject = rej;
    });
    if (this.constructor[species] !== Promise) {
      throw new TypeError("CancellablePromise does not support transparent subclassing. Please refrain from overriding the [Symbol.species] static property.");
    }
    let promise = {
      promise: this,
      resolve,
      reject,
      get oncancelled() {
        return oncancelled != null ? oncancelled : null;
      },
      set oncancelled(cb) {
        oncancelled = cb != null ? cb : void 0;
      }
    };
    const state = {
      get root() {
        return state;
      },
      resolving: false,
      settled: false
    };
    void Object.defineProperties(this, {
      [barrierSym]: {
        configurable: false,
        enumerable: false,
        writable: true,
        value: null
      },
      [cancelImplSym]: {
        configurable: false,
        enumerable: false,
        writable: false,
        value: cancellerFor(promise, state)
      }
    });
    const rejector = rejectorFor(promise, state);
    try {
      executor(resolverFor(promise, state), rejector);
    } catch (err) {
      if (state.resolving) {
        console.log("Unhandled exception in CancellablePromise executor.", err);
      } else {
        rejector(err);
      }
    }
  }
  /**
   * Cancels immediately the execution of the operation associated with this promise.
   * The promise rejects with a {@link CancelError} instance as reason,
   * with the {@link CancelError#cause} property set to the given argument, if any.
   *
   * Has no effect if called after the promise has already settled;
   * repeated calls in particular are safe, but only the first one
   * will set the cancellation cause.
   *
   * The `CancelError` exception _need not_ be handled explicitly _on the promises that are being cancelled:_
   * cancelling a promise with no attached rejection handler does not trigger an unhandled rejection event.
   * Therefore, the following idioms are all equally correct:
   * ```ts
   * new CancellablePromise((resolve, reject) => { ... }).cancel();
   * new CancellablePromise((resolve, reject) => { ... }).then(...).cancel();
   * new CancellablePromise((resolve, reject) => { ... }).then(...).catch(...).cancel();
   * ```
   * Whenever some cancelled promise in a chain rejects with a `CancelError`
   * with the same cancellation cause as itself, the error will be discarded silently.
   * However, the `CancelError` _will still be delivered_ to all attached rejection handlers
   * added by {@link then} and related methods:
   * ```ts
   * let cancellable = new CancellablePromise((resolve, reject) => { ... });
   * cancellable.then(() => { ... }).catch(console.log);
   * cancellable.cancel(); // A CancelError is printed to the console.
   * ```
   * If the `CancelError` is not handled downstream by the time it reaches
   * a _non-cancelled_ promise, it _will_ trigger an unhandled rejection event,
   * just like normal rejections would:
   * ```ts
   * let cancellable = new CancellablePromise((resolve, reject) => { ... });
   * let chained = cancellable.then(() => { ... }).then(() => { ... }); // No catch...
   * cancellable.cancel(); // Unhandled rejection event on chained!
   * ```
   * Therefore, it is important to either cancel whole promise chains from their tail,
   * as shown in the correct idioms above, or take care of handling errors everywhere.
   *
   * @returns A cancellable promise that _fulfills_ after the cancel callback (if any)
   * and all handlers attached up to the call to cancel have run.
   * If the cancel callback returns a thenable, the promise returned by `cancel`
   * will also wait for that thenable to settle.
   * This enables callers to wait for the cancelled operation to terminate
   * without being forced to handle potential errors at the call site.
   * ```ts
   * cancellable.cancel().then(() => {
   *     // Cleanup finished, it's safe to do something else.
   * }, (err) => {
   *     // Unreachable: the promise returned from cancel will never reject.
   * });
   * ```
   * Note that the returned promise will _not_ handle implicitly any rejection
   * that might have occurred already in the cancelled chain.
   * It will just track whether registered handlers have been executed or not.
   * Therefore, unhandled rejections will never be silently handled by calling cancel.
   */
  cancel(cause) {
    return new _CancellablePromise((resolve) => {
      Promise.all([
        this[cancelImplSym](new CancelError("Promise cancelled.", { cause })),
        currentBarrier(this)
      ]).then(() => resolve(), () => resolve());
    });
  }
  /**
   * Binds promise cancellation to the abort event of the given {@link AbortSignal}.
   * If the signal has already aborted, the promise will be cancelled immediately.
   * When either condition is verified, the cancellation cause will be set
   * to the signal's abort reason (see {@link AbortSignal#reason}).
   *
   * Has no effect if called (or if the signal aborts) _after_ the promise has already settled.
   * Only the first signal to abort will set the cancellation cause.
   *
   * For more details about the cancellation process,
   * see {@link cancel} and the `CancellablePromise` constructor.
   *
   * This method enables `await`ing cancellable promises without having
   * to store them for future cancellation, e.g.:
   * ```ts
   * await longRunningOperation().cancelOn(signal);
   * ```
   * instead of:
   * ```ts
   * let promiseToBeCancelled = longRunningOperation();
   * await promiseToBeCancelled;
   * ```
   *
   * @returns This promise, for method chaining.
   */
  cancelOn(signal) {
    if (signal.aborted) {
      void this.cancel(signal.reason);
    } else {
      signal.addEventListener("abort", () => void this.cancel(signal.reason), { capture: true });
    }
    return this;
  }
  /**
   * Attaches callbacks for the resolution and/or rejection of the `CancellablePromise`.
   *
   * The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
   * with the same semantics as the `oncancelled` argument of the constructor.
   * When the parent promise rejects or is cancelled, the `onrejected` callback will run,
   * _even after the returned promise has been cancelled:_
   * in that case, should it reject or throw, the reason will be wrapped
   * in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
   *
   * @param onfulfilled The callback to execute when the Promise is resolved.
   * @param onrejected The callback to execute when the Promise is rejected.
   * @returns A `CancellablePromise` for the completion of whichever callback is executed.
   * The returned promise is hooked up to propagate cancellation requests up the chain, but not down:
   *
   *   - if the parent promise is cancelled, the `onrejected` handler will be invoked with a `CancelError`
   *     and the returned promise _will resolve regularly_ with its result;
   *   - conversely, if the returned promise is cancelled, _the parent promise is cancelled too;_
   *     the `onrejected` handler will still be invoked with the parent's `CancelError`,
   *     but its result will be discarded
   *     and the returned promise will reject with a `CancelError` as well.
   *
   * The promise returned from {@link cancel} will fulfill only after all attached handlers
   * up the entire promise chain have been run.
   *
   * If either callback returns a cancellable promise,
   * cancellation requests will be diverted to it,
   * and the specified `oncancelled` callback will be discarded.
   */
  then(onfulfilled, onrejected, oncancelled) {
    if (!(this instanceof _CancellablePromise)) {
      throw new TypeError("CancellablePromise.prototype.then called on an invalid object.");
    }
    if (!callable_default(onfulfilled)) {
      onfulfilled = identity;
    }
    if (!callable_default(onrejected)) {
      onrejected = thrower;
    }
    if (onfulfilled === identity && onrejected == thrower) {
      return new _CancellablePromise((resolve) => resolve(this));
    }
    const barrier = {};
    this[barrierSym] = barrier;
    return new _CancellablePromise((resolve, reject) => {
      void super.then(
        (value) => {
          var _a2;
          if (this[barrierSym] === barrier) {
            this[barrierSym] = null;
          }
          (_a2 = barrier.resolve) == null ? void 0 : _a2.call(barrier);
          try {
            resolve(onfulfilled(value));
          } catch (err) {
            reject(err);
          }
        },
        (reason) => {
          var _a2;
          if (this[barrierSym] === barrier) {
            this[barrierSym] = null;
          }
          (_a2 = barrier.resolve) == null ? void 0 : _a2.call(barrier);
          try {
            resolve(onrejected(reason));
          } catch (err) {
            reject(err);
          }
        }
      );
    }, async (cause) => {
      try {
        return oncancelled == null ? void 0 : oncancelled(cause);
      } finally {
        await this.cancel(cause);
      }
    });
  }
  /**
   * Attaches a callback for only the rejection of the Promise.
   *
   * The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
   * with the same semantics as the `oncancelled` argument of the constructor.
   * When the parent promise rejects or is cancelled, the `onrejected` callback will run,
   * _even after the returned promise has been cancelled:_
   * in that case, should it reject or throw, the reason will be wrapped
   * in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
   *
   * It is equivalent to
   * ```ts
   * cancellablePromise.then(undefined, onrejected, oncancelled);
   * ```
   * and the same caveats apply.
   *
   * @returns A Promise for the completion of the callback.
   * Cancellation requests on the returned promise
   * will propagate up the chain to the parent promise,
   * but not in the other direction.
   *
   * The promise returned from {@link cancel} will fulfill only after all attached handlers
   * up the entire promise chain have been run.
   *
   * If `onrejected` returns a cancellable promise,
   * cancellation requests will be diverted to it,
   * and the specified `oncancelled` callback will be discarded.
   * See {@link then} for more details.
   */
  catch(onrejected, oncancelled) {
    return this.then(void 0, onrejected, oncancelled);
  }
  /**
   * Attaches a callback that is invoked when the CancellablePromise is settled (fulfilled or rejected). The
   * resolved value cannot be accessed or modified from the callback.
   * The returned promise will settle in the same state as the original one
   * after the provided callback has completed execution,
   * unless the callback throws or returns a rejecting promise,
   * in which case the returned promise will reject as well.
   *
   * The optional `oncancelled` argument will be invoked when the returned promise is cancelled,
   * with the same semantics as the `oncancelled` argument of the constructor.
   * Once the parent promise settles, the `onfinally` callback will run,
   * _even after the returned promise has been cancelled:_
   * in that case, should it reject or throw, the reason will be wrapped
   * in a {@link CancelledRejectionError} and bubbled up as an unhandled rejection.
   *
   * This method is implemented in terms of {@link then} and the same caveats apply.
   * It is polyfilled, hence available in every OS/webview version.
   *
   * @returns A Promise for the completion of the callback.
   * Cancellation requests on the returned promise
   * will propagate up the chain to the parent promise,
   * but not in the other direction.
   *
   * The promise returned from {@link cancel} will fulfill only after all attached handlers
   * up the entire promise chain have been run.
   *
   * If `onfinally` returns a cancellable promise,
   * cancellation requests will be diverted to it,
   * and the specified `oncancelled` callback will be discarded.
   * See {@link then} for more details.
   */
  finally(onfinally, oncancelled) {
    if (!(this instanceof _CancellablePromise)) {
      throw new TypeError("CancellablePromise.prototype.finally called on an invalid object.");
    }
    if (!callable_default(onfinally)) {
      return this.then(onfinally, onfinally, oncancelled);
    }
    return this.then(
      (value) => _CancellablePromise.resolve(onfinally()).then(() => value),
      (reason) => _CancellablePromise.resolve(onfinally()).then(() => {
        throw reason;
      }),
      oncancelled
    );
  }
  /**
   * We use the `[Symbol.species]` static property, if available,
   * to disable the built-in automatic subclassing features from {@link Promise}.
   * It is critical for performance reasons that extenders do not override this.
   * Once the proposal at https://github.com/tc39/proposal-rm-builtin-subclassing
   * is either accepted or retired, this implementation will have to be revised accordingly.
   *
   * @ignore
   * @internal
   */
  static get [(barrierSym, cancelImplSym, species)]() {
    return Promise;
  }
  static all(values) {
    let collected = Array.from(values);
    const promise = collected.length === 0 ? _CancellablePromise.resolve(collected) : new _CancellablePromise((resolve, reject) => {
      void Promise.all(collected).then(resolve, reject);
    }, (cause) => cancelAll(promise, collected, cause));
    return promise;
  }
  static allSettled(values) {
    let collected = Array.from(values);
    const promise = collected.length === 0 ? _CancellablePromise.resolve(collected) : new _CancellablePromise((resolve, reject) => {
      void Promise.allSettled(collected).then(resolve, reject);
    }, (cause) => cancelAll(promise, collected, cause));
    return promise;
  }
  static any(values) {
    let collected = Array.from(values);
    const promise = collected.length === 0 ? _CancellablePromise.resolve(collected) : new _CancellablePromise((resolve, reject) => {
      void Promise.any(collected).then(resolve, reject);
    }, (cause) => cancelAll(promise, collected, cause));
    return promise;
  }
  static race(values) {
    let collected = Array.from(values);
    const promise = new _CancellablePromise((resolve, reject) => {
      void Promise.race(collected).then(resolve, reject);
    }, (cause) => cancelAll(promise, collected, cause));
    return promise;
  }
  /**
   * Creates a new cancelled CancellablePromise for the provided cause.
   *
   * @group Static Methods
   */
  static cancel(cause) {
    const p = new _CancellablePromise(() => {
    });
    p.cancel(cause);
    return p;
  }
  /**
   * Creates a new CancellablePromise that cancels
   * after the specified timeout, with the provided cause.
   *
   * If the {@link AbortSignal.timeout} factory method is available,
   * it is used to base the timeout on _active_ time rather than _elapsed_ time.
   * Otherwise, `timeout` falls back to {@link setTimeout}.
   *
   * @group Static Methods
   */
  static timeout(milliseconds, cause) {
    const promise = new _CancellablePromise(() => {
    });
    if (AbortSignal && typeof AbortSignal === "function" && AbortSignal.timeout && typeof AbortSignal.timeout === "function") {
      AbortSignal.timeout(milliseconds).addEventListener("abort", () => void promise.cancel(cause));
    } else {
      setTimeout(() => void promise.cancel(cause), milliseconds);
    }
    return promise;
  }
  static sleep(milliseconds, value) {
    return new _CancellablePromise((resolve) => {
      setTimeout(() => resolve(value), milliseconds);
    });
  }
  /**
   * Creates a new rejected CancellablePromise for the provided reason.
   *
   * @group Static Methods
   */
  static reject(reason) {
    return new _CancellablePromise((_, reject) => reject(reason));
  }
  static resolve(value) {
    if (value instanceof _CancellablePromise) {
      return value;
    }
    return new _CancellablePromise((resolve) => resolve(value));
  }
  /**
   * Creates a new CancellablePromise and returns it in an object, along with its resolve and reject functions
   * and a getter/setter for the cancellation callback.
   *
   * This method is polyfilled, hence available in every OS/webview version.
   *
   * @group Static Methods
   */
  static withResolvers() {
    let result = { oncancelled: null };
    result.promise = new _CancellablePromise((resolve, reject) => {
      result.resolve = resolve;
      result.reject = reject;
    }, (cause) => {
      var _a2;
      (_a2 = result.oncancelled) == null ? void 0 : _a2.call(result, cause);
    });
    return result;
  }
};
function cancellerFor(promise, state) {
  let cancellationPromise = void 0;
  return (reason) => {
    if (!state.settled) {
      state.settled = true;
      state.reason = reason;
      promise.reject(reason);
      void Promise.prototype.then.call(promise.promise, void 0, (err) => {
        if (err !== reason) {
          throw err;
        }
      });
    }
    if (!state.reason || !promise.oncancelled) {
      return;
    }
    cancellationPromise = new Promise((resolve) => {
      try {
        resolve(promise.oncancelled(state.reason.cause));
      } catch (err) {
        Promise.reject(new CancelledRejectionError(promise.promise, err, "Unhandled exception in oncancelled callback."));
      }
    }).catch((reason2) => {
      Promise.reject(new CancelledRejectionError(promise.promise, reason2, "Unhandled rejection in oncancelled callback."));
    });
    promise.oncancelled = null;
    return cancellationPromise;
  };
}
function resolverFor(promise, state) {
  return (value) => {
    if (state.resolving) {
      return;
    }
    state.resolving = true;
    if (value === promise.promise) {
      if (state.settled) {
        return;
      }
      state.settled = true;
      promise.reject(new TypeError("A promise cannot be resolved with itself."));
      return;
    }
    if (value != null && (typeof value === "object" || typeof value === "function")) {
      let then;
      try {
        then = value.then;
      } catch (err) {
        state.settled = true;
        promise.reject(err);
        return;
      }
      if (callable_default(then)) {
        try {
          let cancel = value.cancel;
          if (callable_default(cancel)) {
            const oncancelled = (cause) => {
              Reflect.apply(cancel, value, [cause]);
            };
            if (state.reason) {
              void cancellerFor(__spreadProps(__spreadValues({}, promise), { oncancelled }), state)(state.reason);
            } else {
              promise.oncancelled = oncancelled;
            }
          }
        } catch (e) {
        }
        const newState = {
          root: state.root,
          resolving: false,
          get settled() {
            return this.root.settled;
          },
          set settled(value2) {
            this.root.settled = value2;
          },
          get reason() {
            return this.root.reason;
          }
        };
        const rejector = rejectorFor(promise, newState);
        try {
          Reflect.apply(then, value, [resolverFor(promise, newState), rejector]);
        } catch (err) {
          rejector(err);
        }
        return;
      }
    }
    if (state.settled) {
      return;
    }
    state.settled = true;
    promise.resolve(value);
  };
}
function rejectorFor(promise, state) {
  return (reason) => {
    if (state.resolving) {
      return;
    }
    state.resolving = true;
    if (state.settled) {
      try {
        if (reason instanceof CancelError && state.reason instanceof CancelError && Object.is(reason.cause, state.reason.cause)) {
          return;
        }
      } catch (e) {
      }
      void Promise.reject(new CancelledRejectionError(promise.promise, reason));
    } else {
      state.settled = true;
      promise.reject(reason);
    }
  };
}
function cancelAll(parent, values, cause) {
  const results = [];
  for (const value of values) {
    let cancel;
    try {
      if (!callable_default(value.then)) {
        continue;
      }
      cancel = value.cancel;
      if (!callable_default(cancel)) {
        continue;
      }
    } catch (e) {
      continue;
    }
    let result;
    try {
      result = Reflect.apply(cancel, value, [cause]);
    } catch (err) {
      Promise.reject(new CancelledRejectionError(parent, err, "Unhandled exception in cancel method."));
      continue;
    }
    if (!result) {
      continue;
    }
    results.push(
      (result instanceof Promise ? result : Promise.resolve(result)).catch((reason) => {
        Promise.reject(new CancelledRejectionError(parent, reason, "Unhandled rejection in cancel method."));
      })
    );
  }
  return Promise.all(results);
}
function identity(x) {
  return x;
}
function thrower(reason) {
  throw reason;
}
function errorMessage(err) {
  try {
    if (err instanceof Error || typeof err !== "object" || err.toString !== Object.prototype.toString) {
      return "" + err;
    }
  } catch (e) {
  }
  try {
    return JSON.stringify(err);
  } catch (e) {
  }
  try {
    return Object.prototype.toString.call(err);
  } catch (e) {
  }
  return "<could not convert error to string>";
}
function currentBarrier(promise) {
  var _a2;
  let pwr = (_a2 = promise[barrierSym]) != null ? _a2 : {};
  if (!("promise" in pwr)) {
    Object.assign(pwr, promiseWithResolvers());
  }
  if (promise[barrierSym] == null) {
    pwr.resolve();
    promise[barrierSym] = pwr;
  }
  return pwr.promise;
}
var promiseWithResolvers = Promise.withResolvers;
if (promiseWithResolvers && typeof promiseWithResolvers === "function") {
  promiseWithResolvers = promiseWithResolvers.bind(Promise);
} else {
  promiseWithResolvers = function() {
    let resolve;
    let reject;
    const promise = new Promise((res, rej) => {
      resolve = res;
      reject = rej;
    });
    return { promise, resolve, reject };
  };
}

// desktop/@wailsio/runtime/src/calls.ts
window._wails = window._wails || {};
window._wails.callResultHandler = resultHandler;
window._wails.callErrorHandler = errorHandler;
var call7 = newRuntimeCaller(objectNames.Call);
var cancelCall = newRuntimeCaller(objectNames.CancelCall);
var callResponses = /* @__PURE__ */ new Map();
var CallBinding = 0;
var CancelMethod = 0;
var RuntimeError = class extends Error {
  /**
   * Constructs a new RuntimeError instance.
   * @param message - The error message.
   * @param options - Options to be forwarded to the Error constructor.
   */
  constructor(message, options) {
    super(message, options);
    this.name = "RuntimeError";
  }
};
function resultHandler(id, data, isJSON) {
  const resolvers = getAndDeleteResponse2(id);
  if (!resolvers) {
    return;
  }
  if (!data) {
    resolvers.resolve(void 0);
  } else if (!isJSON) {
    resolvers.resolve(data);
  } else {
    try {
      resolvers.resolve(JSON.parse(data));
    } catch (err) {
      resolvers.reject(new TypeError("could not parse result: " + err.message, { cause: err }));
    }
  }
}
function errorHandler(id, data, isJSON) {
  const resolvers = getAndDeleteResponse2(id);
  if (!resolvers) {
    return;
  }
  if (!isJSON) {
    resolvers.reject(new Error(data));
  } else {
    let error;
    try {
      error = JSON.parse(data);
    } catch (err) {
      resolvers.reject(new TypeError("could not parse error: " + err.message, { cause: err }));
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
    resolvers.reject(exception);
  }
}
function getAndDeleteResponse2(id) {
  const response = callResponses.get(id);
  callResponses.delete(id);
  return response;
}
function generateID2() {
  let result;
  do {
    result = nanoid();
  } while (callResponses.has(result));
  return result;
}
function Call(options) {
  const id = generateID2();
  const result = CancellablePromise.withResolvers();
  callResponses.set(id, { resolve: result.resolve, reject: result.reject });
  const request = call7(CallBinding, Object.assign({ "call-id": id }, options));
  let running = false;
  request.then(() => {
    running = true;
  }, (err) => {
    callResponses.delete(id);
    result.reject(err);
  });
  const cancel = () => {
    callResponses.delete(id);
    return cancelCall(CancelMethod, { "call-id": id }).catch((err) => {
      console.error("Error while requesting binding call cancellation:", err);
    });
  };
  result.oncancelled = () => {
    if (running) {
      return cancel();
    } else {
      return request.then(cancel);
    }
  };
  return result.promise;
}
function ByName(methodName, ...args) {
  return Call({ methodName, args });
}
function ByID(methodID, ...args) {
  return Call({ methodID, args });
}

// desktop/@wailsio/runtime/src/clipboard.ts
var clipboard_exports = {};
__export(clipboard_exports, {
  SetText: () => SetText,
  Text: () => Text
});
var call8 = newRuntimeCaller(objectNames.Clipboard);
var ClipboardSetText = 0;
var ClipboardText = 1;
function SetText(text) {
  return call8(ClipboardSetText, { text });
}
function Text() {
  return call8(ClipboardText);
}

// desktop/@wailsio/runtime/src/create.ts
var create_exports = {};
__export(create_exports, {
  Any: () => Any,
  Array: () => Array2,
  ByteSlice: () => ByteSlice,
  Map: () => Map2,
  Nullable: () => Nullable,
  Struct: () => Struct
});
function Any(source) {
  return source;
}
function ByteSlice(source) {
  return source == null ? "" : source;
}
function Array2(element) {
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

// desktop/@wailsio/runtime/src/screens.ts
var screens_exports = {};
__export(screens_exports, {
  GetAll: () => GetAll,
  GetCurrent: () => GetCurrent,
  GetPrimary: () => GetPrimary
});
var call9 = newRuntimeCaller(objectNames.Screens);
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

// desktop/@wailsio/runtime/src/index.ts
window._wails = window._wails || {};
window._wails.invoke = invoke;
invoke("wails:runtime:ready");
export {
  application_exports as Application,
  browser_exports as Browser,
  calls_exports as Call,
  CancelError,
  CancellablePromise,
  CancelledRejectionError,
  clipboard_exports as Clipboard,
  create_exports as Create,
  dialogs_exports as Dialogs,
  events_exports as Events,
  flags_exports as Flags,
  screens_exports as Screens,
  system_exports as System,
  wml_exports as WML,
  window_default as Window
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHBsaWNhdGlvbi50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYW5jZWxsYWJsZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG5cclxuaW1wb3J0IFwiLi9jb250ZXh0bWVudS5qc1wiO1xyXG5pbXBvcnQgXCIuL2RyYWcuanNcIjtcclxuXHJcbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXHJcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gXCIuL2FwcGxpY2F0aW9uLmpzXCI7XHJcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xyXG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzLmpzXCI7XHJcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmQuanNcIjtcclxuaW1wb3J0ICogYXMgQ3JlYXRlIGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xyXG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcclxuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuL2V2ZW50cy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi9mbGFncy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnMuanNcIjtcclxuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gXCIuL3N5c3RlbS5qc1wiO1xyXG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sLmpzXCI7XHJcblxyXG5leHBvcnQge1xyXG4gICAgQXBwbGljYXRpb24sXHJcbiAgICBCcm93c2VyLFxyXG4gICAgQ2FsbCxcclxuICAgIENsaXBib2FyZCxcclxuICAgIERpYWxvZ3MsXHJcbiAgICBFdmVudHMsXHJcbiAgICBGbGFncyxcclxuICAgIFNjcmVlbnMsXHJcbiAgICBTeXN0ZW0sXHJcbiAgICBXaW5kb3csXHJcbiAgICBXTUxcclxufTtcclxuXHJcbi8qKlxyXG4gKiBBbiBpbnRlcm5hbCB1dGlsaXR5IGNvbnN1bWVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cclxuICpcclxuICogQGlnbm9yZVxyXG4gKiBAaW50ZXJuYWxcclxuICovXHJcbmV4cG9ydCB7IENyZWF0ZSB9O1xyXG5cclxuZXhwb3J0ICogZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcclxuXHJcbi8vIE5vdGlmeSBiYWNrZW5kXHJcbndpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcclxuU3lzdGVtLmludm9rZShcIndhaWxzOnJ1bnRpbWU6cmVhZHlcIik7XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBPcGVuVVJMIH0gZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xyXG5pbXBvcnQgeyBRdWVzdGlvbiB9IGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcclxuaW1wb3J0IHsgRW1pdCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xyXG5pbXBvcnQgeyBjYW5BYm9ydExpc3RlbmVycywgd2hlblJlYWR5IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcclxuaW1wb3J0IFdpbmRvdyBmcm9tIFwiLi93aW5kb3cuanNcIjtcclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBvcHRpb25hbCBkYXRhLlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gc2VuZC5cclxuICogQHBhcmFtIFtkYXRhPW51bGxdIC0gLSBPcHRpb25hbCBkYXRhIHRvIHNlbmQgYWxvbmcgd2l0aCB0aGUgZXZlbnQuXHJcbiAqL1xyXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpOiB2b2lkIHtcclxuICAgIHZvaWQgRW1pdChldmVudE5hbWUsIGRhdGEpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBzcGVjaWZpZWQgd2luZG93LlxyXG4gKlxyXG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxyXG4gKiBAcGFyYW0gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cclxuICovXHJcbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TmFtZTogc3RyaW5nLCBtZXRob2ROYW1lOiBzdHJpbmcpIHtcclxuICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IFdpbmRvdy5HZXQod2luZG93TmFtZSk7XHJcbiAgICBjb25zdCBtZXRob2QgPSAodGFyZ2V0V2luZG93IGFzIGFueSlbbWV0aG9kTmFtZV07XHJcblxyXG4gICAgaWYgKHR5cGVvZiBtZXRob2QgIT09IFwiZnVuY3Rpb25cIikge1xyXG4gICAgICAgIGNvbnNvbGUuZXJyb3IoYFdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nIG5vdCBmb3VuZGApO1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICB0cnkge1xyXG4gICAgICAgIG1ldGhvZC5jYWxsKHRhcmdldFdpbmRvdyk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgY29uc29sZS5lcnJvcihgRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JzogYCwgZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXNwb25kcyB0byBhIHRyaWdnZXJpbmcgZXZlbnQgYnkgcnVubmluZyBhcHByb3ByaWF0ZSBXTUwgYWN0aW9ucyBmb3IgdGhlIGN1cnJlbnQgdGFyZ2V0LlxyXG4gKi9cclxuZnVuY3Rpb24gb25XTUxUcmlnZ2VyZWQoZXY6IEV2ZW50KTogdm9pZCB7XHJcbiAgICBjb25zdCBlbGVtZW50ID0gZXYuY3VycmVudFRhcmdldCBhcyBFbGVtZW50O1xyXG5cclxuICAgIGZ1bmN0aW9uIHJ1bkVmZmVjdChjaG9pY2UgPSBcIlllc1wiKSB7XHJcbiAgICAgICAgaWYgKGNob2ljZSAhPT0gXCJZZXNcIilcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG5cclxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWV2ZW50Jyk7XHJcbiAgICAgICAgY29uc3QgdGFyZ2V0V2luZG93ID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10YXJnZXQtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRhcmdldC13aW5kb3cnKSB8fCBcIlwiO1xyXG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xyXG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtb3BlbnVybCcpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1vcGVudXJsJyk7XHJcblxyXG4gICAgICAgIGlmIChldmVudFR5cGUgIT09IG51bGwpXHJcbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xyXG4gICAgICAgIGlmICh3aW5kb3dNZXRob2QgIT09IG51bGwpXHJcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xyXG4gICAgICAgIGlmICh1cmwgIT09IG51bGwpXHJcbiAgICAgICAgICAgIHZvaWQgT3BlblVSTCh1cmwpO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xyXG5cclxuICAgIGlmIChjb25maXJtKSB7XHJcbiAgICAgICAgUXVlc3Rpb24oe1xyXG4gICAgICAgICAgICBUaXRsZTogXCJDb25maXJtXCIsXHJcbiAgICAgICAgICAgIE1lc3NhZ2U6IGNvbmZpcm0sXHJcbiAgICAgICAgICAgIERldGFjaGVkOiBmYWxzZSxcclxuICAgICAgICAgICAgQnV0dG9uczogW1xyXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJZZXNcIiB9LFxyXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJOb1wiLCBJc0RlZmF1bHQ6IHRydWUgfVxyXG4gICAgICAgICAgICBdXHJcbiAgICAgICAgfSkudGhlbihydW5FZmZlY3QpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBydW5FZmZlY3QoKTtcclxuICAgIH1cclxufVxyXG5cclxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cclxuY29uc3QgY29udHJvbGxlclN5bSA9IFN5bWJvbChcImNvbnRyb2xsZXJcIik7XHJcbmNvbnN0IHRyaWdnZXJNYXBTeW0gPSBTeW1ib2woXCJ0cmlnZ2VyTWFwXCIpO1xyXG5jb25zdCBlbGVtZW50Q291bnRTeW0gPSBTeW1ib2woXCJlbGVtZW50Q291bnRcIik7XHJcblxyXG4vKipcclxuICogQWJvcnRDb250cm9sbGVyUmVnaXN0cnkgZG9lcyBub3QgYWN0dWFsbHkgcmVtZW1iZXIgYWN0aXZlIGV2ZW50IGxpc3RlbmVyczogaW5zdGVhZFxyXG4gKiBpdCB0aWVzIHRoZW0gdG8gYW4gQWJvcnRTaWduYWwgYW5kIHVzZXMgYW4gQWJvcnRDb250cm9sbGVyIHRvIHJlbW92ZSB0aGVtIGFsbCBhdCBvbmNlLlxyXG4gKi9cclxuY2xhc3MgQWJvcnRDb250cm9sbGVyUmVnaXN0cnkge1xyXG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXHJcbiAgICBbY29udHJvbGxlclN5bV06IEFib3J0Q29udHJvbGxlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcigpIHtcclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyBhbiBvcHRpb25zIG9iamVjdCBmb3IgYWRkRXZlbnRMaXN0ZW5lciB0aGF0IHRpZXMgdGhlIGxpc3RlbmVyXHJcbiAgICAgKiB0byB0aGUgQWJvcnRTaWduYWwgZnJvbSB0aGUgY3VycmVudCBBYm9ydENvbnRyb2xsZXIuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGVsZW1lbnQgLSBBbiBIVE1MIGVsZW1lbnRcclxuICAgICAqIEBwYXJhbSB0cmlnZ2VycyAtIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudHNcclxuICAgICAqL1xyXG4gICAgc2V0KGVsZW1lbnQ6IEVsZW1lbnQsIHRyaWdnZXJzOiBzdHJpbmdbXSk6IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zIHtcclxuICAgICAgICByZXR1cm4geyBzaWduYWw6IHRoaXNbY29udHJvbGxlclN5bV0uc2lnbmFsIH07XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycyBhbmQgcmVzZXRzIHRoZSByZWdpc3RyeS5cclxuICAgICAqL1xyXG4gICAgcmVzZXQoKTogdm9pZCB7XHJcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXS5hYm9ydCgpO1xyXG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBXZWFrTWFwUmVnaXN0cnkgbWFwcyBhY3RpdmUgdHJpZ2dlciBldmVudHMgdG8gZWFjaCBET00gZWxlbWVudCB0aHJvdWdoIGEgV2Vha01hcC5cclxuICogVGhpcyBlbnN1cmVzIHRoYXQgdGhlIG1hcHBpbmcgcmVtYWlucyBwcml2YXRlIHRvIHRoaXMgbW9kdWxlLCB3aGlsZSBzdGlsbCBhbGxvd2luZyBnYXJiYWdlXHJcbiAqIGNvbGxlY3Rpb24gb2YgdGhlIGludm9sdmVkIGVsZW1lbnRzLlxyXG4gKi9cclxuY2xhc3MgV2Vha01hcFJlZ2lzdHJ5IHtcclxuICAgIC8qKiBTdG9yZXMgdGhlIGN1cnJlbnQgZWxlbWVudC10by10cmlnZ2VyIG1hcHBpbmcuICovXHJcbiAgICBbdHJpZ2dlck1hcFN5bV06IFdlYWtNYXA8RWxlbWVudCwgc3RyaW5nW10+O1xyXG4gICAgLyoqIENvdW50cyB0aGUgbnVtYmVyIG9mIGVsZW1lbnRzIHdpdGggYWN0aXZlIFdNTCB0cmlnZ2Vycy4gKi9cclxuICAgIFtlbGVtZW50Q291bnRTeW1dOiBudW1iZXI7XHJcblxyXG4gICAgY29uc3RydWN0b3IoKSB7XHJcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XHJcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgYWN0aXZlIHRyaWdnZXJzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGVsZW1lbnQgLSBBbiBIVE1MIGVsZW1lbnRcclxuICAgICAqIEBwYXJhbSB0cmlnZ2VycyAtIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudFxyXG4gICAgICovXHJcbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xyXG4gICAgICAgIGlmICghdGhpc1t0cmlnZ2VyTWFwU3ltXS5oYXMoZWxlbWVudCkpIHsgdGhpc1tlbGVtZW50Q291bnRTeW1dKys7IH1cclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dLnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XHJcbiAgICAgICAgcmV0dXJuIHt9O1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMuXHJcbiAgICAgKi9cclxuICAgIHJlc2V0KCk6IHZvaWQge1xyXG4gICAgICAgIGlmICh0aGlzW2VsZW1lbnRDb3VudFN5bV0gPD0gMClcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG5cclxuICAgICAgICBmb3IgKGNvbnN0IGVsZW1lbnQgb2YgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCcqJykpIHtcclxuICAgICAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxyXG4gICAgICAgICAgICAgICAgYnJlYWs7XHJcblxyXG4gICAgICAgICAgICBjb25zdCB0cmlnZ2VycyA9IHRoaXNbdHJpZ2dlck1hcFN5bV0uZ2V0KGVsZW1lbnQpO1xyXG4gICAgICAgICAgICBpZiAodHJpZ2dlcnMgIT0gbnVsbCkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0tLTsgfVxyXG5cclxuICAgICAgICAgICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzIHx8IFtdKVxyXG4gICAgICAgICAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0gPSBuZXcgV2Vha01hcCgpO1xyXG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50U3ltXSA9IDA7XHJcbiAgICB9XHJcbn1cclxuXHJcbmNvbnN0IHRyaWdnZXJSZWdpc3RyeSA9IGNhbkFib3J0TGlzdGVuZXJzKCkgPyBuZXcgQWJvcnRDb250cm9sbGVyUmVnaXN0cnkoKSA6IG5ldyBXZWFrTWFwUmVnaXN0cnkoKTtcclxuXHJcbi8qKlxyXG4gKiBBZGRzIGV2ZW50IGxpc3RlbmVycyB0byB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxMaXN0ZW5lcnMoZWxlbWVudDogRWxlbWVudCk6IHZvaWQge1xyXG4gICAgY29uc3QgdHJpZ2dlclJlZ0V4cCA9IC9cXFMrL2c7XHJcbiAgICBjb25zdCB0cmlnZ2VyQXR0ciA9IChlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIik7XHJcbiAgICBjb25zdCB0cmlnZ2Vyczogc3RyaW5nW10gPSBbXTtcclxuXHJcbiAgICBsZXQgbWF0Y2g7XHJcbiAgICB3aGlsZSAoKG1hdGNoID0gdHJpZ2dlclJlZ0V4cC5leGVjKHRyaWdnZXJBdHRyKSkgIT09IG51bGwpXHJcbiAgICAgICAgdHJpZ2dlcnMucHVzaChtYXRjaFswXSk7XHJcblxyXG4gICAgY29uc3Qgb3B0aW9ucyA9IHRyaWdnZXJSZWdpc3RyeS5zZXQoZWxlbWVudCwgdHJpZ2dlcnMpO1xyXG4gICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzKVxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTY2hlZHVsZXMgYW4gYXV0b21hdGljIHJlbG9hZCBvZiBXTUwgdG8gYmUgcGVyZm9ybWVkIGFzIHNvb24gYXMgdGhlIGRvY3VtZW50IGlzIGZ1bGx5IGxvYWRlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbmFibGUoKTogdm9pZCB7XHJcbiAgICB3aGVuUmVhZHkoUmVsb2FkKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBSZWxvYWQoKTogdm9pZCB7XHJcbiAgICB0cmlnZ2VyUmVnaXN0cnkucmVzZXQoKTtcclxuICAgIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0sIFt3bWwtd2luZG93XSwgW3dtbC1vcGVudXJsXSwgW2RhdGEtd21sLWV2ZW50XSwgW2RhdGEtd21sLXdpbmRvd10sIFtkYXRhLXdtbC1vcGVudXJsXScpLmZvckVhY2goYWRkV01MTGlzdGVuZXJzKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5Ccm93c2VyKTtcclxuXHJcbmNvbnN0IEJyb3dzZXJPcGVuVVJMID0gMDtcclxuXHJcbi8qKlxyXG4gKiBPcGVuIGEgYnJvd3NlciB3aW5kb3cgdG8gdGhlIGdpdmVuIFVSTC5cclxuICpcclxuICogQHBhcmFtIHVybCAtIFRoZSBVUkwgdG8gb3BlblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5VUkwodXJsOiBzdHJpbmcgfCBVUkwpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsOiB1cmwudG9TdHJpbmcoKX0pO1xyXG59XHJcbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9haS9uYW5vaWRcclxuXHJcbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxyXG4vL1xyXG4vLyBDb3B5cmlnaHQgMjAxNyBBbmRyZXkgU2l0bmlrIDxhbmRyZXlAc2l0bmlrLnJ1PlxyXG4vL1xyXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mXHJcbi8vIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWwgaW5cclxuLy8gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cyB0b1xyXG4vLyB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsIGNvcGllcyBvZlxyXG4vLyB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sXHJcbi8vICAgICBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcclxuLy9cclxuLy8gICAgIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxyXG4vLyBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLlxyXG4vL1xyXG4vLyAgICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxyXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTU1xyXG4vLyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1JcclxuLy8gQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLCBXSEVUSEVSXHJcbi8vIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwgT1VUIE9GIE9SIElOXHJcbi8vIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuXHJcblxyXG4vLyBUaGlzIGFscGhhYmV0IHVzZXMgYEEtWmEtejAtOV8tYCBzeW1ib2xzLlxyXG4vLyBUaGUgb3JkZXIgb2YgY2hhcmFjdGVycyBpcyBvcHRpbWl6ZWQgZm9yIGJldHRlciBnemlwIGFuZCBicm90bGkgY29tcHJlc3Npb24uXHJcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIHNhbWUgZmlsZSAod29ya3MgYm90aCBmb3IgZ3ppcCBhbmQgYnJvdGxpKTpcclxuLy8gYCd1c2VgLCBgYW5kb21gLCBhbmQgYHJpY3QnYFxyXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBicm90bGkgZGVmYXVsdCBkaWN0aW9uYXJ5OlxyXG4vLyBgLTI2VGAsIGAxOTgzYCwgYDQwcHhgLCBgNzVweGAsIGBidXNoYCwgYGphY2tgLCBgbWluZGAsIGB2ZXJ5YCwgYW5kIGB3b2xmYFxyXG5jb25zdCB1cmxBbHBoYWJldCA9XHJcbiAgICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBuYW5vaWQoc2l6ZTogbnVtYmVyID0gMjEpOiBzdHJpbmcge1xyXG4gICAgbGV0IGlkID0gJydcclxuICAgIC8vIEEgY29tcGFjdCBhbHRlcm5hdGl2ZSBmb3IgYGZvciAodmFyIGkgPSAwOyBpIDwgc3RlcDsgaSsrKWAuXHJcbiAgICBsZXQgaSA9IHNpemUgfCAwXHJcbiAgICB3aGlsZSAoaS0tKSB7XHJcbiAgICAgICAgLy8gYHwgMGAgaXMgbW9yZSBjb21wYWN0IGFuZCBmYXN0ZXIgdGhhbiBgTWF0aC5mbG9vcigpYC5cclxuICAgICAgICBpZCArPSB1cmxBbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIDY0KSB8IDBdXHJcbiAgICB9XHJcbiAgICByZXR1cm4gaWRcclxufVxyXG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xyXG5cclxuY29uc3QgcnVudGltZVVSTCA9IHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XHJcblxyXG4vLyBPYmplY3QgTmFtZXNcclxuZXhwb3J0IGNvbnN0IG9iamVjdE5hbWVzID0gT2JqZWN0LmZyZWV6ZSh7XHJcbiAgICBDYWxsOiAwLFxyXG4gICAgQ2xpcGJvYXJkOiAxLFxyXG4gICAgQXBwbGljYXRpb246IDIsXHJcbiAgICBFdmVudHM6IDMsXHJcbiAgICBDb250ZXh0TWVudTogNCxcclxuICAgIERpYWxvZzogNSxcclxuICAgIFdpbmRvdzogNixcclxuICAgIFNjcmVlbnM6IDcsXHJcbiAgICBTeXN0ZW06IDgsXHJcbiAgICBCcm93c2VyOiA5LFxyXG4gICAgQ2FuY2VsQ2FsbDogMTAsXHJcbn0pO1xyXG5leHBvcnQgbGV0IGNsaWVudElkID0gbmFub2lkKCk7XHJcblxyXG4vKipcclxuICogQ3JlYXRlcyBhIG5ldyBydW50aW1lIGNhbGxlciB3aXRoIHNwZWNpZmllZCBJRC5cclxuICpcclxuICogQHBhcmFtIG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICogQHJldHVybiBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0OiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZyA9ICcnKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZDogbnVtYmVyLCBhcmdzOiBhbnkgPSBudWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn1cclxuXHJcbmFzeW5jIGZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55PiB7XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcclxuICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwib2JqZWN0XCIsIG9iamVjdElELnRvU3RyaW5nKCkpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kLnRvU3RyaW5nKCkpO1xyXG4gICAgaWYgKGFyZ3MpIHsgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTsgfVxyXG5cclxuICAgIGxldCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+ID0ge1xyXG4gICAgICAgIFtcIngtd2FpbHMtY2xpZW50LWlkXCJdOiBjbGllbnRJZFxyXG4gICAgfVxyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBoZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XHJcbiAgICB9XHJcblxyXG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2godXJsLCB7IGhlYWRlcnMgfSk7XHJcbiAgICBpZiAoIXJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xyXG5cclxuLy8gc2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbndpbmRvdy5fd2FpbHMuZGlhbG9nRXJyb3JDYWxsYmFjayA9IGRpYWxvZ0Vycm9yQ2FsbGJhY2s7XHJcbndpbmRvdy5fd2FpbHMuZGlhbG9nUmVzdWx0Q2FsbGJhY2sgPSBkaWFsb2dSZXN1bHRDYWxsYmFjaztcclxuXHJcbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8UHJvbWlzZVdpdGhSZXNvbHZlcnM8YW55PiwgXCJwcm9taXNlXCI+O1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRGlhbG9nKTtcclxuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcDxzdHJpbmcsIFByb21pc2VSZXNvbHZlcnM+KCk7XHJcblxyXG4vLyBEZWZpbmUgY29uc3RhbnRzIGZyb20gdGhlIGBtZXRob2RzYCBvYmplY3QgaW4gVGl0bGUgQ2FzZVxyXG5jb25zdCBEaWFsb2dJbmZvID0gMDtcclxuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XHJcbmNvbnN0IERpYWxvZ0Vycm9yID0gMjtcclxuY29uc3QgRGlhbG9nUXVlc3Rpb24gPSAzO1xyXG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XHJcbmNvbnN0IERpYWxvZ1NhdmVGaWxlID0gNTtcclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgT3BlbkZpbGVEaWFsb2dPcHRpb25zIHtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cclxuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cclxuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXHJcbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXHJcbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cclxuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLiAqL1xyXG4gICAgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXHJcbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xyXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cclxuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xyXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cclxuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XHJcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cclxuICAgIFRpdGxlPzogc3RyaW5nO1xyXG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgTWVzc2FnZT86IHN0cmluZztcclxuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cclxuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XHJcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIERpcmVjdG9yeT86IHN0cmluZztcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xyXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFNhdmVGaWxlRGlhbG9nT3B0aW9ucyB7XHJcbiAgICAvKiogRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIEZpbGVuYW1lPzogc3RyaW5nO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xyXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xyXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cclxuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cclxuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xyXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xyXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cclxuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXHJcbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cclxuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXHJcbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xyXG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXHJcbiAgICBUaXRsZT86IHN0cmluZztcclxuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XHJcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXHJcbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xyXG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cclxuICAgIERldGFjaGVkPzogYm9vbGVhbjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBNZXNzYWdlRGlhbG9nT3B0aW9ucyB7XHJcbiAgICAvKiogVGhlIHRpdGxlIG9mIHRoZSBkaWFsb2cgd2luZG93LiAqL1xyXG4gICAgVGl0bGU/OiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBNZXNzYWdlPzogc3RyaW5nO1xyXG4gICAgLyoqIEFycmF5IG9mIGJ1dHRvbiBvcHRpb25zIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIEJ1dHRvbnM/OiBCdXR0b25bXTtcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuICovXHJcbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgQnV0dG9uIHtcclxuICAgIC8qKiBUZXh0IHRoYXQgYXBwZWFycyB3aXRoaW4gdGhlIGJ1dHRvbi4gKi9cclxuICAgIExhYmVsPzogc3RyaW5nO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgY2FuY2VsIGFuIG9wZXJhdGlvbiB3aGVuIGNsaWNrZWQuICovXHJcbiAgICBJc0NhbmNlbD86IGJvb2xlYW47XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBiZSB0aGUgZGVmYXVsdCBhY3Rpb24gd2hlbiB0aGUgdXNlciBwcmVzc2VzIGVudGVyLiAqL1xyXG4gICAgSXNEZWZhdWx0PzogYm9vbGVhbjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBGaWxlRmlsdGVyIHtcclxuICAgIC8qKiBEaXNwbGF5IG5hbWUgZm9yIHRoZSBmaWx0ZXIsIGl0IGNvdWxkIGJlIFwiVGV4dCBGaWxlc1wiLCBcIkltYWdlc1wiIGV0Yy4gKi9cclxuICAgIERpc3BsYXlOYW1lPzogc3RyaW5nO1xyXG4gICAgLyoqIFBhdHRlcm4gdG8gbWF0Y2ggZm9yIHRoZSBmaWx0ZXIsIGUuZy4gXCIqLnR4dDsqLm1kXCIgZm9yIHRleHQgbWFya2Rvd24gZmlsZXMuICovXHJcbiAgICBQYXR0ZXJuPzogc3RyaW5nO1xyXG59XHJcblxyXG4vKipcclxuICogSGFuZGxlcyB0aGUgcmVzdWx0IG9mIGEgZGlhbG9nIHJlcXVlc3QuXHJcbiAqXHJcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcmVxdWVzdCB0byBoYW5kbGUgdGhlIHJlc3VsdCBmb3IuXHJcbiAqIEBwYXJhbSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxyXG4gKiBAcGFyYW0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2dSZXN1bHRDYWxsYmFjayhpZDogc3RyaW5nLCBkYXRhOiBzdHJpbmcsIGlzSlNPTjogYm9vbGVhbik6IHZvaWQge1xyXG4gICAgbGV0IHJlc29sdmVycyA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcclxuICAgIGlmICghcmVzb2x2ZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcclxuICAgICAgICB9IGNhdGNoIChlcnI6IGFueSkge1xyXG4gICAgICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJjb3VsZCBub3QgcGFyc2UgcmVzdWx0OiBcIiArIGVyci5tZXNzYWdlLCB7IGNhdXNlOiBlcnIgfSkpO1xyXG4gICAgICAgIH1cclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUoZGF0YSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIHRoZSBlcnJvciBmcm9tIGEgZGlhbG9nIHJlcXVlc3QuXHJcbiAqXHJcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxyXG4gKiBAcGFyYW0gbWVzc2FnZSAtIEFuIGVycm9yIG1lc3NhZ2UuXHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkOiBzdHJpbmcsIG1lc3NhZ2U6IHN0cmluZyk6IHZvaWQge1xyXG4gICAgZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpPy5yZWplY3QobmV3IHdpbmRvdy5FcnJvcihtZXNzYWdlKSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgZGlhbG9nUmVzcG9uc2VzIG1hcC5cclxuICpcclxuICogQHBhcmFtIGlkIC0gVGhlIElEIG9mIHRoZSByZXNwb25zZSB0byBiZSByZXRyaWV2ZWQgYW5kIHJlbW92ZWQuXHJcbiAqIEByZXR1cm5zIFRoZSByZXNwb25zZSBvYmplY3QgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCwgaWYgYW55LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQ6IHN0cmluZyk6IFByb21pc2VSZXNvbHZlcnMgfCB1bmRlZmluZWQge1xyXG4gICAgY29uc3QgcmVzcG9uc2UgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgcmV0dXJuIHJlc3BvbnNlO1xyXG59XHJcblxyXG4vKipcclxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cclxuICpcclxuICogQHJldHVybnMgQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgZGlhbG9nUmVzcG9uc2VzIHNldC5cclxuICovXHJcbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKTogc3RyaW5nIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSB0eXBlIC0gRGlhbG9nIHR5cGUuXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyBmb3IgdGhlIGRpYWxvZy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nLlxyXG4gKi9cclxuZnVuY3Rpb24gZGlhbG9nKHR5cGU6IG51bWJlciwgb3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMgfCBPcGVuRmlsZURpYWxvZ09wdGlvbnMgfCBTYXZlRmlsZURpYWxvZ09wdGlvbnMgPSB7fSk6IFByb21pc2U8YW55PiB7XHJcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlLCByZWplY3QgfSk7XHJcbiAgICAgICAgY2FsbCh0eXBlLCBPYmplY3QuYXNzaWduKHsgXCJkaWFsb2ctaWRcIjogaWQgfSwgb3B0aW9ucykpLmNhdGNoKChlcnI6IGFueSkgPT4ge1xyXG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICAgICAgcmVqZWN0KGVycik7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGFuIGluZm8gZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSB3YXJuaW5nIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhbiBlcnJvciBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIHF1ZXN0aW9uIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgb25lIG9yIG1vcmUgZmlsZXMgdG8gb3Blbi5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSBvciBsaXN0IG9mIGZpbGVzLCBvciBhIGJsYW5rIHN0cmluZy9lbXB0eSBsaXN0IGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbjogdHJ1ZSB9KTogUHJvbWlzZTxzdHJpbmdbXT47XHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogZmFsc2UgfCB1bmRlZmluZWQgfSk6IFByb21pc2U8c3RyaW5nPjtcclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+O1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT4geyByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKSA/PyBbXTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgYSBmaWxlIHRvIHNhdmUuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUsIG9yIGEgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9uczogU2F2ZUZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7IH1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5pbXBvcnQgeyBldmVudExpc3RlbmVycywgTGlzdGVuZXIsIGxpc3RlbmVyT2ZmIH0gZnJvbSBcIi4vbGlzdGVuZXIuanNcIjtcclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkV2ZW50cyk7XHJcbmNvbnN0IEVtaXRNZXRob2QgPSAwO1xyXG5cclxuZXhwb3J0IHsgVHlwZXMgfSBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xyXG5cclxuLyoqXHJcbiAqIFRoZSB0eXBlIG9mIGhhbmRsZXJzIGZvciBhIGdpdmVuIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IHR5cGUgQ2FsbGJhY2sgPSAoZXY6IFdhaWxzRXZlbnQpID0+IHZvaWQ7XHJcblxyXG4vKipcclxuICogUmVwcmVzZW50cyBhIHN5c3RlbSBldmVudCBvciBhIGN1c3RvbSBldmVudCBlbWl0dGVkIHRocm91Z2ggd2FpbHMtcHJvdmlkZWQgZmFjaWxpdGllcy5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcclxuICAgIC8qKlxyXG4gICAgICogVGhlIG5hbWUgb2YgdGhlIGV2ZW50LlxyXG4gICAgICovXHJcbiAgICBuYW1lOiBzdHJpbmc7XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBPcHRpb25hbCBkYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZW1pdHRlZCBldmVudC5cclxuICAgICAqL1xyXG4gICAgZGF0YTogYW55O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogTmFtZSBvZiB0aGUgb3JpZ2luYXRpbmcgd2luZG93LiBPbWl0dGVkIGZvciBhcHBsaWNhdGlvbiBldmVudHMuXHJcbiAgICAgKiBXaWxsIGJlIG92ZXJyaWRkZW4gaWYgc2V0IG1hbnVhbGx5LlxyXG4gICAgICovXHJcbiAgICBzZW5kZXI/OiBzdHJpbmc7XHJcblxyXG4gICAgY29uc3RydWN0b3IobmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKSB7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcclxuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQ6IGFueSkge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcclxuICAgIGlmICghbGlzdGVuZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGxldCB3YWlsc0V2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnQubmFtZSwgZXZlbnQuZGF0YSk7XHJcbiAgICBpZiAoJ3NlbmRlcicgaW4gZXZlbnQpIHtcclxuICAgICAgICB3YWlsc0V2ZW50LnNlbmRlciA9IGV2ZW50LnNlbmRlcjtcclxuICAgIH1cclxuXHJcbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+ICFsaXN0ZW5lci5kaXNwYXRjaCh3YWlsc0V2ZW50KSk7XHJcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlciBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCBtdWx0aXBsZSB0aW1lcyBmb3IgYSBzcGVjaWZpYyBldmVudC5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxyXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHBhcmFtIG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxyXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiBDYWxsYmFjaywgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcclxuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIHNwZWNpZmllZCBldmVudCBvY2N1cnMuXHJcbiAqXHJcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXHJcbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiBDYWxsYmFjayk6ICgpID0+IHZvaWQge1xyXG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgb25seSBvbmNlIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50LlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXHJcbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxyXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiBDYWxsYmFjayk6ICgpID0+IHZvaWQge1xyXG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZXMgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnRzIHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZiguLi5ldmVudE5hbWVzOiBbc3RyaW5nLCAuLi5zdHJpbmdbXV0pOiB2b2lkIHtcclxuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVtb3ZlcyBhbGwgZXZlbnQgbGlzdGVuZXJzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpOiB2b2lkIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0cyBhbiBldmVudCB1c2luZyB0aGUgbmFtZSBhbmQgZGF0YS5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgd2lsbCBiZSBmdWxmaWxsZWQgb25jZSB0aGUgZXZlbnQgaGFzIGJlZW4gZW1pdHRlZC5cclxuICogQHBhcmFtIG5hbWUgLSB0aGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdC5cclxuICogQHBhcmFtIGRhdGEgLSB0aGUgZGF0YSB0byBiZSBzZW50IHdpdGggdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQobmFtZTogc3RyaW5nLCBkYXRhPzogYW55KTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICBsZXQgZXZlbnQ6IFdhaWxzRXZlbnQ7XHJcblxyXG4gICAgaWYgKHR5cGVvZiBuYW1lID09PSAnb2JqZWN0JyAmJiBuYW1lICE9PSBudWxsICYmICduYW1lJyBpbiBuYW1lICYmICdkYXRhJyBpbiBuYW1lKSB7XHJcbiAgICAgICAgLy8gSWYgbmFtZSBpcyBhbiBvYmplY3Qgd2l0aCBhIG5hbWUgcHJvcGVydHksIHVzZSBpdCBkaXJlY3RseVxyXG4gICAgICAgIGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQobmFtZVsnbmFtZSddLCBuYW1lWydkYXRhJ10pO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICAvLyBPdGhlcndpc2UgdXNlIHRoZSBzdGFuZGFyZCBwYXJhbWV0ZXJzXHJcbiAgICAgICAgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChuYW1lIGFzIHN0cmluZywgZGF0YSk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgZXZlbnQpO1xyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLy8gVGhlIGZvbGxvd2luZyB1dGlsaXRpZXMgaGF2ZSBiZWVuIGZhY3RvcmVkIG91dCBvZiAuL2V2ZW50cy50c1xyXG4vLyBmb3IgdGVzdGluZyBwdXJwb3Nlcy5cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXA8c3RyaW5nLCBMaXN0ZW5lcltdPigpO1xyXG5cclxuZXhwb3J0IGNsYXNzIExpc3RlbmVyIHtcclxuICAgIGV2ZW50TmFtZTogc3RyaW5nO1xyXG4gICAgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQ7XHJcbiAgICBtYXhDYWxsYmFja3M6IG51bWJlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQsIG1heENhbGxiYWNrczogbnVtYmVyKSB7XHJcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XHJcbiAgICAgICAgdGhpcy5jYWxsYmFjayA9IGNhbGxiYWNrO1xyXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xyXG4gICAgfVxyXG5cclxuICAgIGRpc3BhdGNoKGRhdGE6IGFueSk6IGJvb2xlYW4ge1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIHRoaXMuY2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZXJyKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyOiBMaXN0ZW5lcik6IHZvaWQge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChsaXN0ZW5lci5ldmVudE5hbWUpO1xyXG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGxpc3RlbmVyLmV2ZW50TmFtZSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChsaXN0ZW5lci5ldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBDeW5oeXJjaHd5ZCB5IGZmZWlsIGhvbiB5biBhd3RvbWF0aWcuIFBFSURJV0NIIFx1MDBDMiBNT0RJV0xcbi8vIFRoaXMgZmlsZSBpcyBhdXRvbWF0aWNhbGx5IGdlbmVyYXRlZC4gRE8gTk9UIEVESVRcblxuZXhwb3J0IGNvbnN0IFR5cGVzID0gT2JqZWN0LmZyZWV6ZSh7XG5cdFdpbmRvd3M6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFQTVBvd2VyU2V0dGluZ0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU2V0dGluZ0NoYW5nZVwiLFxuXHRcdEFQTVBvd2VyU3RhdHVzQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTdGF0dXNDaGFuZ2VcIixcblx0XHRBUE1SZXN1bWVBdXRvbWF0aWM6IFwid2luZG93czpBUE1SZXN1bWVBdXRvbWF0aWNcIixcblx0XHRBUE1SZXN1bWVTdXNwZW5kOiBcIndpbmRvd3M6QVBNUmVzdW1lU3VzcGVuZFwiLFxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcIndpbmRvd3M6QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWQ6IFwid2luZG93czpXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZFwiLFxuXHRcdFdpbmRvd0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0JhY2tncm91bmRFcmFzZTogXCJ3aW5kb3dzOldpbmRvd0JhY2tncm91bmRFcmFzZVwiLFxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcIndpbmRvd3M6V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwid2luZG93czpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIndpbmRvd3M6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJ3aW5kb3dzOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHRcdFdpbmRvd0VuZE1vdmU6IFwid2luZG93czpXaW5kb3dFbmRNb3ZlXCIsXG5cdFx0V2luZG93RW5kUmVzaXplOiBcIndpbmRvd3M6V2luZG93RW5kUmVzaXplXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dIaWRlOiBcIndpbmRvd3M6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0luYWN0aXZlOiBcIndpbmRvd3M6V2luZG93SW5hY3RpdmVcIixcblx0XHRXaW5kb3dLZXlEb3duOiBcIndpbmRvd3M6V2luZG93S2V5RG93blwiLFxuXHRcdFdpbmRvd0tleVVwOiBcIndpbmRvd3M6V2luZG93S2V5VXBcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dOb25DbGllbnRIaXQ6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRIaXRcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZURvd246IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZURvd25cIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZUxlYXZlOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlTW92ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTW92ZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlVXA6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZVVwXCIsXG5cdFx0V2luZG93UGFpbnQ6IFwid2luZG93czpXaW5kb3dQYWludFwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93U2V0Rm9jdXM6IFwid2luZG93czpXaW5kb3dTZXRGb2N1c1wiLFxuXHRcdFdpbmRvd1Nob3c6IFwid2luZG93czpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93U3RhcnRNb3ZlOiBcIndpbmRvd3M6V2luZG93U3RhcnRNb3ZlXCIsXG5cdFx0V2luZG93U3RhcnRSZXNpemU6IFwid2luZG93czpXaW5kb3dTdGFydFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1pPcmRlckNoYW5nZWQ6IFwid2luZG93czpXaW5kb3daT3JkZXJDaGFuZ2VkXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1heGltaXNlXCIsXG5cdH0pLFxuXHRNYWM6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyT3JpZW50YXRpb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyT3JpZW50YXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW5cIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZEFkZEl0ZW06IFwibWFjOk1lbnVEaWRBZGRJdGVtXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51RGlkRGlzcGxheUl0ZW06IFwibWFjOk1lbnVEaWREaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZFBvcFVwOiBcIm1hYzpNZW51RGlkUG9wVXBcIixcblx0XHRNZW51RGlkUmVtb3ZlSXRlbTogXCJtYWM6TWVudURpZFJlbW92ZUl0ZW1cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvbjogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW1cIixcblx0XHRNZW51RGlkVXBkYXRlOiBcIm1hYzpNZW51RGlkVXBkYXRlXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsRGlzcGxheUl0ZW06IFwibWFjOk1lbnVXaWxsRGlzcGxheUl0ZW1cIixcblx0XHRNZW51V2lsbEVuZFRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51V2lsbEhpZ2hsaWdodEl0ZW1cIixcblx0XHRNZW51V2lsbE9wZW46IFwibWFjOk1lbnVXaWxsT3BlblwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvbjogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uXCIsXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW1cIixcblx0XHRNZW51V2lsbFVwZGF0ZTogXCJtYWM6TWVudVdpbGxVcGRhdGVcIixcblx0XHRXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXaW5kb3dEaWRCZWNvbWVLZXk6IFwibWFjOldpbmRvd0RpZEJlY29tZUtleVwiLFxuXHRcdFdpbmRvd0RpZEJlY29tZU1haW46IFwibWFjOldpbmRvd0RpZEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dEaWRCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dEaWRCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQWxwaGE6IFwibWFjOldpbmRvd0RpZENoYW5nZUFscGhhXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nTG9jYXRpb25cIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblByb2ZpbGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblByb2ZpbGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZENoYW5nZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWREZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWREZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkRW5kU2hlZXQ6IFwibWFjOldpbmRvd0RpZEVuZFNoZWV0XCIsXG5cdFx0V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxuXHRcdFdpbmRvd0RpZEZvY3VzOiBcIm1hYzpXaW5kb3dEaWRGb2N1c1wiLFxuXHRcdFdpbmRvd0RpZE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWRNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPZmZTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcblx0XHRXaW5kb3dEaWRSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIm1hYzpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93RGlkVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkWm9vbTogXCJtYWM6V2luZG93RGlkWm9vbVwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93SGlkZTogXCJtYWM6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIm1hYzpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwibWFjOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJtYWM6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIm1hYzpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlXCIsXG5cdFx0V2luZG93U2hvdzogXCJtYWM6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJtYWM6V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJtYWM6V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJtYWM6V2luZG93Wm9vbVJlc2V0XCIsXG5cdH0pLFxuXHRMaW51eDogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRXaW5kb3dEZWxldGVFdmVudDogXCJsaW51eDpXaW5kb3dEZWxldGVFdmVudFwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibGludXg6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJsaW51eDpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dGb2N1c0luOiBcImxpbnV4OldpbmRvd0ZvY3VzSW5cIixcblx0XHRXaW5kb3dGb2N1c091dDogXCJsaW51eDpXaW5kb3dGb2N1c091dFwiLFxuXHRcdFdpbmRvd0xvYWRDaGFuZ2VkOiBcImxpbnV4OldpbmRvd0xvYWRDaGFuZ2VkXCIsXG5cdH0pLFxuXHRDb21tb246IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiQXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZVwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwiV2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJXaW5kb3dGaWxlc0Ryb3BwZWRcIixcblx0XHRXaW5kb3dGb2N1czogXCJXaW5kb3dGb2N1c1wiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwiV2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd0hpZGU6IFwiV2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJXaW5kb3dMb3N0Rm9jdXNcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcIldpbmRvd1J1bnRpbWVSZWFkeVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwiV2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dab29tOiBcIldpbmRvd1pvb21cIixcblx0XHRXaW5kb3dab29tSW46IFwiV2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcIldpbmRvd1pvb21SZXNldFwiLFxuXHR9KSxcbn0pO1xuIiwgIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxyXG4gKlxyXG4gKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBtZXNzYWdlIHRvIGJlIGxvZ2dlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlOiBhbnkpIHtcclxuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG4gICAgY29uc29sZS5sb2coXHJcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXHJcbiAgICApO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIHdoZXRoZXIgdGhlIHdlYnZpZXcgc3VwcG9ydHMgdGhlIHtAbGluayBNb3VzZUV2ZW50I2J1dHRvbnN9IHByb3BlcnR5LlxyXG4gKiBMb29raW5nIGF0IHlvdSBtYWNPUyBIaWdoIFNpZXJyYSFcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBjYW5UcmFja0J1dHRvbnMoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKG5ldyBNb3VzZUV2ZW50KCdtb3VzZWRvd24nKSkuYnV0dG9ucyA9PT0gMDtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXHJcbiAqIChzZWUgaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmcvZW4tVVMvZG9jcy9XZWIvQVBJL0V2ZW50VGFyZ2V0L2FkZEV2ZW50TGlzdGVuZXIjc2lnbmFsKS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBjYW5BYm9ydExpc3RlbmVycygpIHtcclxuICAgIGlmICghRXZlbnRUYXJnZXQgfHwgIUFib3J0U2lnbmFsIHx8ICFBYm9ydENvbnRyb2xsZXIpXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG5cclxuICAgIGxldCByZXN1bHQgPSB0cnVlO1xyXG5cclxuICAgIGNvbnN0IHRhcmdldCA9IG5ldyBFdmVudFRhcmdldCgpO1xyXG4gICAgY29uc3QgY29udHJvbGxlciA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xyXG4gICAgY29udHJvbGxlci5hYm9ydCgpO1xyXG4gICAgdGFyZ2V0LmRpc3BhdGNoRXZlbnQobmV3IEN1c3RvbUV2ZW50KCd0ZXN0JykpO1xyXG5cclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXNvbHZlcyB0aGUgY2xvc2VzdCBIVE1MRWxlbWVudCBhbmNlc3RvciBvZiBhbiBldmVudCdzIHRhcmdldC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBldmVudFRhcmdldChldmVudDogRXZlbnQpOiBIVE1MRWxlbWVudCB7XHJcbiAgICBpZiAoZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpIHtcclxuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0O1xyXG4gICAgfSBlbHNlIGlmICghKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSAmJiBldmVudC50YXJnZXQgaW5zdGFuY2VvZiBOb2RlKSB7XHJcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldC5wYXJlbnRFbGVtZW50ID8/IGRvY3VtZW50LmJvZHk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiBkb2N1bWVudC5ib2R5O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKioqXHJcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcclxuXHJcbiBCU0QgMi1DbGF1c2UgTGljZW5zZVxyXG5cclxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxyXG4gQWxsIHJpZ2h0cyByZXNlcnZlZC5cclxuXHJcbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcclxuIG1vZGlmaWNhdGlvbiwgYXJlIHBlcm1pdHRlZCBwcm92aWRlZCB0aGF0IHRoZSBmb2xsb3dpbmcgY29uZGl0aW9ucyBhcmUgbWV0OlxyXG5cclxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xyXG4gbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIuXHJcblxyXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxyXG4gdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvblxyXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXHJcblxyXG4gVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyBcIkFTIElTXCJcclxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcclxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxyXG4gRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1JTIEJFIExJQUJMRVxyXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcclxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXHJcbiBTRVJWSUNFUzsgTE9TUyBPRiBVU0UsIERBVEEsIE9SIFBST0ZJVFM7IE9SIEJVU0lORVNTIElOVEVSUlVQVElPTikgSE9XRVZFUlxyXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcclxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXHJcbiBPRiBUSElTIFNPRlRXQVJFLCBFVkVOIElGIEFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLlxyXG5cclxuICoqKi9cclxuXHJcbmxldCBpc1JlYWR5ID0gZmFsc2U7XHJcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCAoKSA9PiB7IGlzUmVhZHkgPSB0cnVlIH0pO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjazogKCkgPT4gdm9pZCkge1xyXG4gICAgaWYgKGlzUmVhZHkgfHwgZG9jdW1lbnQucmVhZHlTdGF0ZSA9PT0gJ2NvbXBsZXRlJykge1xyXG4gICAgICAgIGNhbGxiYWNrKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCBjYWxsYmFjayk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHR5cGUgeyBTY3JlZW4gfSBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XHJcblxyXG5jb25zdCBQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAwO1xyXG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xyXG5jb25zdCBDbG9zZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSAyO1xyXG5jb25zdCBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgPSAzO1xyXG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xyXG5jb25zdCBGb2N1c01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1O1xyXG5jb25zdCBGb3JjZVJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgPSA2O1xyXG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xyXG5jb25zdCBHZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA4O1xyXG5jb25zdCBHZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA5O1xyXG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcclxuY29uc3QgSGlkZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTE7XHJcbmNvbnN0IElzRm9jdXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDEyO1xyXG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcclxuY29uc3QgSXNNYXhpbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTQ7XHJcbmNvbnN0IElzTWluaW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE1O1xyXG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcclxuY29uc3QgTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTc7XHJcbmNvbnN0IE5hbWVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDE4O1xyXG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcclxuY29uc3QgUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMjA7XHJcbmNvbnN0IFJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDIxO1xyXG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcclxuY29uc3QgUmVzdG9yZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMjM7XHJcbmNvbnN0IFNldFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICA9IDI0O1xyXG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcclxuY29uc3QgU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCAgICAgICAgID0gMjY7XHJcbmNvbnN0IFNldEZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICAgICA9IDI3O1xyXG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcclxuY29uc3QgU2V0TWF4U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMjk7XHJcbmNvbnN0IFNldE1pblNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDMwO1xyXG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcclxuY29uc3QgU2V0UmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgID0gMzI7XHJcbmNvbnN0IFNldFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDMzO1xyXG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcclxuY29uc3QgU2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzU7XHJcbmNvbnN0IFNob3dNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM2O1xyXG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcclxuY29uc3QgVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgID0gMzg7XHJcbmNvbnN0IFRvZ2dsZU1heGltaXNlTWV0aG9kICAgICAgICAgICAgICA9IDM5O1xyXG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MDtcclxuY29uc3QgVW5NYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDE7XHJcbmNvbnN0IFVuTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQyO1xyXG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0MztcclxuY29uc3QgWm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XHJcbmNvbnN0IFpvb21Jbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDQ1O1xyXG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NjtcclxuY29uc3QgWm9vbVJlc2V0TWV0aG9kICAgICAgICAgICAgICAgICAgID0gNDc7XHJcblxyXG4vKipcclxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgcG9zaXRpb24gb2YgYSB3aW5kb3cuXHJcbiAqL1xyXG5pbnRlcmZhY2UgUG9zaXRpb24ge1xyXG4gICAgLyoqIFRoZSBob3Jpem9udGFsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICB4OiBudW1iZXI7XHJcbiAgICAvKiogVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICB5OiBudW1iZXI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBzaXplIG9mIGEgd2luZG93LlxyXG4gKi9cclxuaW50ZXJmYWNlIFNpemUge1xyXG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgd2luZG93LiAqL1xyXG4gICAgd2lkdGg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICBoZWlnaHQ6IG51bWJlcjtcclxufVxyXG5cclxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cclxuY29uc3QgY2FsbGVyU3ltID0gU3ltYm9sKFwiY2FsbGVyXCIpO1xyXG5cclxuY2xhc3MgV2luZG93IHtcclxuICAgIC8vIFByaXZhdGUgZmllbGRzLlxyXG4gICAgcHJpdmF0ZSBbY2FsbGVyU3ltXTogKG1lc3NhZ2U6IG51bWJlciwgYXJncz86IGFueSkgPT4gUHJvbWlzZTxhbnk+O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogSW5pdGlhbGlzZXMgYSB3aW5kb3cgb2JqZWN0IHdpdGggdGhlIHNwZWNpZmllZCBuYW1lLlxyXG4gICAgICpcclxuICAgICAqIEBwcml2YXRlXHJcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB0YXJnZXQgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcgPSAnJykge1xyXG4gICAgICAgIHRoaXNbY2FsbGVyU3ltXSA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuV2luZG93LCBuYW1lKVxyXG5cclxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xyXG4gICAgICAgIGZvciAoY29uc3QgbWV0aG9kIG9mIE9iamVjdC5nZXRPd25Qcm9wZXJ0eU5hbWVzKFdpbmRvdy5wcm90b3R5cGUpKSB7XHJcbiAgICAgICAgICAgIGlmIChcclxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXHJcbiAgICAgICAgICAgICAgICAmJiB0eXBlb2YgKHRoaXMgYXMgYW55KVttZXRob2RdID09PSBcImZ1bmN0aW9uXCJcclxuICAgICAgICAgICAgKSB7XHJcbiAgICAgICAgICAgICAgICAodGhpcyBhcyBhbnkpW21ldGhvZF0gPSAodGhpcyBhcyBhbnkpW21ldGhvZF0uYmluZCh0aGlzKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cclxuICAgICAqIEByZXR1cm5zIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXHJcbiAgICAgKi9cclxuICAgIEdldChuYW1lOiBzdHJpbmcpOiBXaW5kb3cge1xyXG4gICAgICAgIHJldHVybiBuZXcgV2luZG93KG5hbWUpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQb3NpdGlvbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cclxuICAgICAqL1xyXG4gICAgQ2VudGVyKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2VudGVyTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENsb3NlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBDbG9zZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENsb3NlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cclxuICAgICAqL1xyXG4gICAgRGlzYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRW5hYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXHJcbiAgICAgKi9cclxuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBGb2N1c2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIEZvY3VzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRm9jdXNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cclxuICAgICAqL1xyXG4gICAgRm9yY2VSZWxvYWQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb3JjZVJlbG9hZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTd2l0Y2hlcyB0aGUgd2luZG93IHRvIGZ1bGxzY3JlZW4gbW9kZS5cclxuICAgICAqL1xyXG4gICAgRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgc2NyZWVuIHRoYXQgdGhlIHdpbmRvdyBpcyBvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgc2NyZWVuIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG9uLlxyXG4gICAgICovXHJcbiAgICBHZXRTY3JlZW4oKTogUHJvbWlzZTxTY3JlZW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBjdXJyZW50IHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB6b29tIGxldmVsLlxyXG4gICAgICovXHJcbiAgICBHZXRab29tKCk6IFByb21pc2U8bnVtYmVyPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRab29tTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBIZWlnaHQoKTogUHJvbWlzZTxudW1iZXI+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhlaWdodE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBIaWRlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSGlkZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZm9jdXNlZC5cclxuICAgICAqL1xyXG4gICAgSXNGb2N1c2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNGb2N1c2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmdWxsc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBJc0Z1bGxzY3JlZW4oKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0Z1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWF4aW1pc2VkLlxyXG4gICAgICovXHJcbiAgICBJc01heGltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWF4aW1pc2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1pbmltaXNlZC5cclxuICAgICAqL1xyXG4gICAgSXNNaW5pbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc01pbmltaXNlZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNYXhpbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBNaW5pbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNaW5pbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgTmFtZSgpOiBQcm9taXNlPHN0cmluZz4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTmFtZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBPcGVucyB0aGUgZGV2ZWxvcG1lbnQgdG9vbHMgcGFuZS5cclxuICAgICAqL1xyXG4gICAgT3BlbkRldlRvb2xzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oT3BlbkRldlRvb2xzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBSZWxhdGl2ZVBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbGF0aXZlUG9zaXRpb25NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXHJcbiAgICAgKi9cclxuICAgIFJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbG9hZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGUuXHJcbiAgICAgKi9cclxuICAgIFJlc2l6YWJsZSgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc2l6YWJsZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXN0b3JlcyB0aGUgd2luZG93IHRvIGl0cyBwcmV2aW91cyBzdGF0ZSBpZiBpdCB3YXMgcHJldmlvdXNseSBtaW5pbWlzZWQsIG1heGltaXNlZCBvciBmdWxsc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBSZXN0b3JlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzdG9yZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBhbHdheXNPblRvcCAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgc3RheSBvbiB0b3AuXHJcbiAgICAgKi9cclxuICAgIFNldEFsd2F5c09uVG9wKGFsd2F5c09uVG9wOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRBbHdheXNPblRvcE1ldGhvZCwgeyBhbHdheXNPblRvcCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHIgLSBUaGUgZGVzaXJlZCByZWQgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cclxuICAgICAqIEBwYXJhbSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cclxuICAgICAqIEBwYXJhbSBiIC0gVGhlIGRlc2lyZWQgYmx1ZSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICogQHBhcmFtIGEgLSBUaGUgZGVzaXJlZCBhbHBoYSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICovXHJcbiAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyKHI6IG51bWJlciwgZzogbnVtYmVyLCBiOiBudW1iZXIsIGE6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCwgeyByLCBnLCBiLCBhIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGZyYW1lbGVzcyAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgZnJhbWVsZXNzLlxyXG4gICAgICovXHJcbiAgICBTZXRGcmFtZWxlc3MoZnJhbWVsZXNzOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRGcmFtZWxlc3NNZXRob2QsIHsgZnJhbWVsZXNzIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZW5hYmxlZCAtIFdoZXRoZXIgdGhlIGZ1bGxzY3JlZW4gYnV0dG9uIHNob3VsZCBiZSBlbmFibGVkLlxyXG4gICAgICovXHJcbiAgICBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZChlbmFibGVkOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCwgeyBlbmFibGVkIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWF4aW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldE1heFNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtaW5pbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRNaW5TaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRNaW5TaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFJlbGF0aXZlUG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gcmVzaXphYmxlIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSByZXNpemFibGUuXHJcbiAgICAgKi9cclxuICAgIFNldFJlc2l6YWJsZShyZXNpemFibGU6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFJlc2l6YWJsZU1ldGhvZCwgeyByZXNpemFibGUgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgdGl0bGUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gdGl0bGUgLSBUaGUgZGVzaXJlZCB0aXRsZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRUaXRsZSh0aXRsZTogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWwuXHJcbiAgICAgKi9cclxuICAgIFNldFpvb20oem9vbTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTaG93cyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2hvd01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTaXplKCk6IFByb21pc2U8U2l6ZT4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2l6ZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmdWxsc2NyZWVuIGFuZCBub3JtYWwuXHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZUZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVGdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxyXG4gICAgICovXHJcbiAgICBUb2dnbGVNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFVuLWZ1bGxzY3JlZW5zIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFVuRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBVbi1tYXhpbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgVW5NYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFVuTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1pbmltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgV2lkdGgoKTogUHJvbWlzZTxudW1iZXI+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFdpZHRoTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFpvb21zIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFpvb20oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEluY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICovXHJcbiAgICBab29tSW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tSW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXHJcbiAgICAgKi9cclxuICAgIFpvb21PdXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tT3V0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlc2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICovXHJcbiAgICBab29tUmVzZXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tUmVzZXRNZXRob2QpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogVGhlIHdpbmRvdyB3aXRoaW4gd2hpY2ggdGhlIHNjcmlwdCBpcyBydW5uaW5nLlxyXG4gKi9cclxuY29uc3QgdGhpc1dpbmRvdyA9IG5ldyBXaW5kb3coJycpO1xyXG5cclxuZXhwb3J0IGRlZmF1bHQgdGhpc1dpbmRvdztcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XHJcblxyXG4vLyBOT1RFOiB0aGUgZm9sbG93aW5nIG1ldGhvZHMgTVVTVCBiZSBpbXBvcnRlZCBleHBsaWNpdGx5IGJlY2F1c2Ugb2YgaG93IGVzYnVpbGQgaW5qZWN0aW9uIHdvcmtzXHJcbmltcG9ydCB7IEVuYWJsZSBhcyBFbmFibGVXTUwgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XHJcbmltcG9ydCB7IGRlYnVnTG9nIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzXCI7XHJcblxyXG53aW5kb3cud2FpbHMgPSBSdW50aW1lO1xyXG5FbmFibGVXTUwoKTtcclxuXHJcbmlmIChERUJVRykge1xyXG4gICAgZGVidWdMb2coXCJXYWlscyBSdW50aW1lIExvYWRlZFwiKVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlN5c3RlbSk7XHJcblxyXG5jb25zdCBTeXN0ZW1Jc0RhcmtNb2RlID0gMDtcclxuY29uc3QgU3lzdGVtRW52aXJvbm1lbnQgPSAxO1xyXG5cclxuY29uc3QgX2ludm9rZSA9IChmdW5jdGlvbiAoKSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZSkge1xyXG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3KTtcclxuICAgICAgICB9IGVsc2UgaWYgKCh3aW5kb3cgYXMgYW55KS53ZWJraXQ/Lm1lc3NhZ2VIYW5kbGVycz8uWydleHRlcm5hbCddPy5wb3N0TWVzc2FnZSkge1xyXG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10ucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXSk7XHJcbiAgICAgICAgfVxyXG4gICAgfSBjYXRjaChlKSB7fVxyXG5cclxuICAgIGNvbnNvbGUud2FybignXFxuJWNcdTI2QTBcdUZFMEYgQnJvd3NlciBFbnZpcm9ubWVudCBEZXRlY3RlZCAlY1xcblxcbiVjT25seSBVSSBwcmV2aWV3cyBhcmUgYXZhaWxhYmxlIGluIHRoZSBicm93c2VyLiBGb3IgZnVsbCBmdW5jdGlvbmFsaXR5LCBwbGVhc2UgcnVuIHRoZSBhcHBsaWNhdGlvbiBpbiBkZXNrdG9wIG1vZGUuXFxuTW9yZSBpbmZvcm1hdGlvbiBhdDogaHR0cHM6Ly92My53YWlscy5pby9sZWFybi9idWlsZC8jdXNpbmctYS1icm93c2VyLWZvci1kZXZlbG9wbWVudFxcbicsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNmZmZmZmY7IGNvbG9yOiAjMDAwMDAwOyBmb250LXdlaWdodDogYm9sZDsgcGFkZGluZzogNHB4IDhweDsgYm9yZGVyLXJhZGl1czogNHB4OyBib3JkZXI6IDJweCBzb2xpZCAjMDAwMDAwOycsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6IHRyYW5zcGFyZW50OycsXHJcbiAgICAgICAgJ2NvbG9yOiAjZmZmZmZmOyBmb250LXN0eWxlOiBpdGFsaWM7IGZvbnQtd2VpZ2h0OiBib2xkOycpO1xyXG4gICAgcmV0dXJuIG51bGw7XHJcbn0pKCk7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZzogYW55KTogdm9pZCB7XHJcbiAgICBfaW52b2tlPy4obXNnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGEgYm9vbGVhbiB2YWx1ZSBpbmRpY2F0aW5nIGlmIHRoZSBzeXN0ZW0gaXMgaW4gZGFyayBtb2RlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1Jc0RhcmtNb2RlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyB0aGUgY2FwYWJpbGl0aWVzLlxyXG4gKi9cclxuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIENhcGFiaWxpdGllcygpOiBQcm9taXNlPFJlY29yZDxzdHJpbmcsIGFueT4+IHtcclxuICAgIGxldCByZXNwb25zZSA9IGF3YWl0IGZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKTtcclxuICAgIGlmIChyZXNwb25zZS5vaykge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcImNvdWxkIG5vdCBmZXRjaCBjYXBhYmlsaXRpZXM6IFwiICsgcmVzcG9uc2Uuc3RhdHVzVGV4dCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgT1NJbmZvIHtcclxuICAgIC8qKiBUaGUgYnJhbmRpbmcgb2YgdGhlIE9TLiAqL1xyXG4gICAgQnJhbmRpbmc6IHN0cmluZztcclxuICAgIC8qKiBUaGUgSUQgb2YgdGhlIE9TLiAqL1xyXG4gICAgSUQ6IHN0cmluZztcclxuICAgIC8qKiBUaGUgbmFtZSBvZiB0aGUgT1MuICovXHJcbiAgICBOYW1lOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIHZlcnNpb24gb2YgdGhlIE9TLiAqL1xyXG4gICAgVmVyc2lvbjogc3RyaW5nO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIEVudmlyb25tZW50SW5mbyB7XHJcbiAgICAvKiogVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLiAqL1xyXG4gICAgQXJjaDogc3RyaW5nO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGFwcGxpY2F0aW9uIGlzIHJ1bm5pbmcgaW4gZGVidWcgbW9kZSwgb3RoZXJ3aXNlIGZhbHNlLiAqL1xyXG4gICAgRGVidWc6IGJvb2xlYW47XHJcbiAgICAvKiogVGhlIG9wZXJhdGluZyBzeXN0ZW0gaW4gdXNlLiAqL1xyXG4gICAgT1M6IHN0cmluZztcclxuICAgIC8qKiBEZXRhaWxzIG9mIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xyXG4gICAgT1NJbmZvOiBPU0luZm87XHJcbiAgICAvKiogQWRkaXRpb25hbCBwbGF0Zm9ybSBpbmZvcm1hdGlvbi4gKi9cclxuICAgIFBsYXRmb3JtSW5mbzogUmVjb3JkPHN0cmluZywgYW55PjtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyBPUyBhbmQgc3lzdGVtIGFyY2hpdGVjdHVyZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpOiBQcm9taXNlPEVudmlyb25tZW50SW5mbz4ge1xyXG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtRW52aXJvbm1lbnQpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cclxuICpcclxuICogQHJldHVybiBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MsIG90aGVyd2lzZSBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc1dpbmRvd3MoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJ3aW5kb3dzXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cclxuICpcclxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwibGludXhcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBhIG1hY09TIG9wZXJhdGluZyBzeXN0ZW0uXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIG1hY09TLCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNNYWMoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJkYXJ3aW5cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYW1kNjRcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTSgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBSZXR1cm5zIHRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZSwgb3RoZXJ3aXNlIHJldHVybnMgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBUk02NCgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtNjRcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlcG9ydHMgd2hldGhlciB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiBCb29sZWFuKHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuRGVidWcpO1xyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB7IElzRGVidWcgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcclxuaW1wb3J0IHsgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlsc1wiO1xyXG5cclxuLy8gc2V0dXBcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcclxuXHJcbmNvbnN0IENvbnRleHRNZW51T3BlbiA9IDA7XHJcblxyXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQ6IHN0cmluZywgeDogbnVtYmVyLCB5OiBudW1iZXIsIGRhdGE6IGFueSk6IHZvaWQge1xyXG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudDogTW91c2VFdmVudCkge1xyXG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xyXG5cclxuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XHJcbiAgICBjb25zdCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpLnRyaW0oKTtcclxuXHJcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIGNvbnN0IGRhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xyXG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZGF0YSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQsIHRhcmdldCk7XHJcbiAgICB9XHJcbn1cclxuXHJcblxyXG4vKlxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuXHJcblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXHJcbiovXHJcbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQ6IE1vdXNlRXZlbnQsIHRhcmdldDogSFRNTEVsZW1lbnQpIHtcclxuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxyXG4gICAgaWYgKElzRGVidWcoKSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICBzd2l0Y2ggKHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCkpIHtcclxuICAgICAgICBjYXNlICdzaG93JzpcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGNhc2UgJ2hpZGUnOlxyXG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcclxuICAgIGlmICh0YXJnZXQuaXNDb250ZW50RWRpdGFibGUpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxyXG4gICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xyXG4gICAgY29uc3QgaGFzU2VsZWN0aW9uID0gc2VsZWN0aW9uICYmIHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDA7XHJcbiAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XHJcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XHJcbiAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XHJcbiAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcclxuICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xyXG4gICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xyXG4gICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IHRhcmdldCkge1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiB0YWcgaXMgaW5wdXQgb3IgdGV4dGFyZWEuXHJcbiAgICBpZiAodGFyZ2V0IGluc3RhbmNlb2YgSFRNTElucHV0RWxlbWVudCB8fCB0YXJnZXQgaW5zdGFuY2VvZiBIVE1MVGV4dEFyZWFFbGVtZW50KSB7XHJcbiAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIXRhcmdldC5yZWFkT25seSAmJiAhdGFyZ2V0LmRpc2FibGVkKSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgdGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleSBmcm9tIHRoZSBmbGFnIG1hcC5cclxuICpcclxuICogQHBhcmFtIGtleSAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cclxuICogQHJldHVybiBUaGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5OiBzdHJpbmcpOiBhbnkge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gd2luZG93Ll93YWlscy5mbGFnc1trZXldO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleSArIFwiJzogXCIgKyBlLCB7IGNhdXNlOiBlIH0pO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBpbnZva2UsIElzV2luZG93cyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xyXG5pbXBvcnQgeyBHZXRGbGFnIH0gZnJvbSBcIi4vZmxhZ3MuanNcIjtcclxuaW1wb3J0IHsgY2FuVHJhY2tCdXR0b25zLCBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XHJcblxyXG4vLyBTZXR1cFxyXG5sZXQgY2FuRHJhZyA9IGZhbHNlO1xyXG5sZXQgZHJhZ2dpbmcgPSBmYWxzZTtcclxuXHJcbmxldCByZXNpemFibGUgPSBmYWxzZTtcclxubGV0IGNhblJlc2l6ZSA9IGZhbHNlO1xyXG5sZXQgcmVzaXppbmcgPSBmYWxzZTtcclxubGV0IHJlc2l6ZUVkZ2U6IHN0cmluZyA9IFwiXCI7XHJcbmxldCBkZWZhdWx0Q3Vyc29yID0gXCJhdXRvXCI7XHJcblxyXG5sZXQgYnV0dG9ucyA9IDA7XHJcbmNvbnN0IGJ1dHRvbnNUcmFja2VkID0gY2FuVHJhY2tCdXR0b25zKCk7XHJcblxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5zZXRSZXNpemFibGUgPSAodmFsdWU6IGJvb2xlYW4pOiB2b2lkID0+IHtcclxuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xyXG4gICAgaWYgKCFyZXNpemFibGUpIHtcclxuICAgICAgICAvLyBTdG9wIHJlc2l6aW5nIGlmIGluIHByb2dyZXNzLlxyXG4gICAgICAgIGNhblJlc2l6ZSA9IHJlc2l6aW5nID0gZmFsc2U7XHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9XHJcbn07XHJcblxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcclxuZm9yIChjb25zdCBldiBvZiBbJ2NsaWNrJywgJ2NvbnRleHRtZW51JywgJ2RibGNsaWNrJ10pIHtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKGV2LCBzdXBwcmVzc0V2ZW50LCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHN1cHByZXNzRXZlbnQoZXZlbnQ6IEV2ZW50KSB7XHJcbiAgICAvLyBTdXBwcmVzcyBjbGljayBldmVudHMgd2hpbGUgcmVzaXppbmcgb3IgZHJhZ2dpbmcuXHJcbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcclxuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBVc2UgY29uc3RhbnRzIHRvIGF2b2lkIGNvbXBhcmluZyBzdHJpbmdzIG11bHRpcGxlIHRpbWVzLlxyXG5jb25zdCBNb3VzZURvd24gPSAwO1xyXG5jb25zdCBNb3VzZVVwICAgPSAxO1xyXG5jb25zdCBNb3VzZU1vdmUgPSAyO1xyXG5cclxuZnVuY3Rpb24gdXBkYXRlKGV2ZW50OiBNb3VzZUV2ZW50KSB7XHJcbiAgICAvLyBXaW5kb3dzIHN1cHByZXNzZXMgbW91c2UgZXZlbnRzIGF0IHRoZSBlbmQgb2YgZHJhZ2dpbmcgb3IgcmVzaXppbmcsXHJcbiAgICAvLyBzbyB3ZSBuZWVkIHRvIGJlIHNtYXJ0IGFuZCBzeW50aGVzaXplIGJ1dHRvbiBldmVudHMuXHJcblxyXG4gICAgbGV0IGV2ZW50VHlwZTogbnVtYmVyLCBldmVudEJ1dHRvbnMgPSBldmVudC5idXR0b25zO1xyXG4gICAgc3dpdGNoIChldmVudC50eXBlKSB7XHJcbiAgICAgICAgY2FzZSAnbW91c2Vkb3duJzpcclxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VEb3duO1xyXG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgfCAoMSA8PCBldmVudC5idXR0b24pOyB9XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIGNhc2UgJ21vdXNldXAnOlxyXG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZVVwO1xyXG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgJiB+KDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICBkZWZhdWx0OlxyXG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZU1vdmU7XHJcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9uczsgfVxyXG4gICAgICAgICAgICBicmVhaztcclxuICAgIH1cclxuXHJcbiAgICBsZXQgcmVsZWFzZWQgPSBidXR0b25zICYgfmV2ZW50QnV0dG9ucztcclxuICAgIGxldCBwcmVzc2VkID0gZXZlbnRCdXR0b25zICYgfmJ1dHRvbnM7XHJcblxyXG4gICAgYnV0dG9ucyA9IGV2ZW50QnV0dG9ucztcclxuXHJcbiAgICAvLyBTeW50aGVzaXplIGEgcmVsZWFzZS1wcmVzcyBzZXF1ZW5jZSBpZiB3ZSBkZXRlY3QgYSBwcmVzcyBvZiBhbiBhbHJlYWR5IHByZXNzZWQgYnV0dG9uLlxyXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VEb3duICYmICEocHJlc3NlZCAmIGV2ZW50LmJ1dHRvbikpIHtcclxuICAgICAgICByZWxlYXNlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xyXG4gICAgICAgIHByZXNzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBTdXBwcmVzcyBhbGwgYnV0dG9uIGV2ZW50cyBkdXJpbmcgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLFxyXG4gICAgLy8gdW5sZXNzIHRoaXMgaXMgYSBtb3VzZXVwIGV2ZW50IHRoYXQgaXMgZW5kaW5nIGEgZHJhZyBhY3Rpb24uXHJcbiAgICBpZiAoXHJcbiAgICAgICAgZXZlbnRUeXBlICE9PSBNb3VzZU1vdmUgLy8gRmFzdCBwYXRoIGZvciBtb3VzZW1vdmVcclxuICAgICAgICAmJiByZXNpemluZ1xyXG4gICAgICAgIHx8IChcclxuICAgICAgICAgICAgZHJhZ2dpbmdcclxuICAgICAgICAgICAgJiYgKFxyXG4gICAgICAgICAgICAgICAgZXZlbnRUeXBlID09PSBNb3VzZURvd25cclxuICAgICAgICAgICAgICAgIHx8IGV2ZW50LmJ1dHRvbiAhPT0gMFxyXG4gICAgICAgICAgICApXHJcbiAgICAgICAgKVxyXG4gICAgKSB7XHJcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XHJcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBIYW5kbGUgcmVsZWFzZXNcclxuICAgIGlmIChyZWxlYXNlZCAmIDEpIHsgcHJpbWFyeVVwKGV2ZW50KTsgfVxyXG4gICAgLy8gSGFuZGxlIHByZXNzZXNcclxuICAgIGlmIChwcmVzc2VkICYgMSkgeyBwcmltYXJ5RG93bihldmVudCk7IH1cclxuXHJcbiAgICAvLyBIYW5kbGUgbW91c2Vtb3ZlXHJcbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZU1vdmUpIHsgb25Nb3VzZU1vdmUoZXZlbnQpOyB9O1xyXG59XHJcblxyXG5mdW5jdGlvbiBwcmltYXJ5RG93bihldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xyXG4gICAgLy8gUmVzZXQgcmVhZGluZXNzIHN0YXRlLlxyXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xyXG4gICAgY2FuUmVzaXplID0gZmFsc2U7XHJcblxyXG4gICAgLy8gSWdub3JlIHJlcGVhdGVkIGNsaWNrcyBvbiBtYWNPUyBhbmQgTGludXguXHJcbiAgICBpZiAoIUlzV2luZG93cygpKSB7XHJcbiAgICAgICAgaWYgKGV2ZW50LnR5cGUgPT09ICdtb3VzZWRvd24nICYmIGV2ZW50LmJ1dHRvbiA9PT0gMCAmJiBldmVudC5kZXRhaWwgIT09IDEpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICBpZiAocmVzaXplRWRnZSkge1xyXG4gICAgICAgIC8vIFJlYWR5IHRvIHJlc2l6ZSBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lLlxyXG4gICAgICAgIGNhblJlc2l6ZSA9IHRydWU7XHJcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWcgb3BlcmF0aW9ucyB3aGVuIG9uIHJlc2l6ZSBlZGdlcy5cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gUmV0cmlldmUgdGFyZ2V0IGVsZW1lbnRcclxuICAgIGNvbnN0IHRhcmdldCA9IGV2ZW50VGFyZ2V0KGV2ZW50KTtcclxuXHJcbiAgICAvLyBSZWFkeSB0byBkcmFnIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUgb24gYSBkcmFnZ2FibGUgZWxlbWVudC5cclxuICAgIC8vIElnbm9yZSBjbGlja3Mgb24gdGhlIHNjcm9sbGJhci5cclxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KTtcclxuICAgIGNhbkRyYWcgPSAoXHJcbiAgICAgICAgc3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpLnRyaW0oKSA9PT0gXCJkcmFnXCJcclxuICAgICAgICAmJiAoXHJcbiAgICAgICAgICAgIGV2ZW50Lm9mZnNldFggLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdMZWZ0KSA8IHRhcmdldC5jbGllbnRXaWR0aFxyXG4gICAgICAgICAgICAmJiBldmVudC5vZmZzZXRZIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nVG9wKSA8IHRhcmdldC5jbGllbnRIZWlnaHRcclxuICAgICAgICApXHJcbiAgICApO1xyXG59XHJcblxyXG5mdW5jdGlvbiBwcmltYXJ5VXAoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcclxuICAgIC8vIFN0b3AgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLlxyXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xyXG4gICAgZHJhZ2dpbmcgPSBmYWxzZTtcclxuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xyXG4gICAgcmVzaXppbmcgPSBmYWxzZTtcclxufVxyXG5cclxuY29uc3QgY3Vyc29yRm9yRWRnZSA9IE9iamVjdC5mcmVlemUoe1xyXG4gICAgXCJzZS1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxyXG4gICAgXCJzdy1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxyXG4gICAgXCJudy1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxyXG4gICAgXCJuZS1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxyXG4gICAgXCJ3LXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxyXG4gICAgXCJuLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxyXG4gICAgXCJzLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxyXG4gICAgXCJlLXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxyXG59KVxyXG5cclxuZnVuY3Rpb24gc2V0UmVzaXplKGVkZ2U/OiBrZXlvZiB0eXBlb2YgY3Vyc29yRm9yRWRnZSk6IHZvaWQge1xyXG4gICAgaWYgKGVkZ2UpIHtcclxuICAgICAgICBpZiAoIXJlc2l6ZUVkZ2UpIHsgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yOyB9XHJcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBjdXJzb3JGb3JFZGdlW2VkZ2VdO1xyXG4gICAgfSBlbHNlIGlmICghZWRnZSAmJiByZXNpemVFZGdlKSB7XHJcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBkZWZhdWx0Q3Vyc29yO1xyXG4gICAgfVxyXG5cclxuICAgIHJlc2l6ZUVkZ2UgPSBlZGdlIHx8IFwiXCI7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XHJcbiAgICBpZiAoY2FuUmVzaXplICYmIHJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICAvLyBTdGFydCByZXNpemluZy5cclxuICAgICAgICByZXNpemluZyA9IHRydWU7XHJcbiAgICAgICAgaW52b2tlKFwid2FpbHM6cmVzaXplOlwiICsgcmVzaXplRWRnZSk7XHJcbiAgICB9IGVsc2UgaWYgKGNhbkRyYWcpIHtcclxuICAgICAgICAvLyBTdGFydCBkcmFnZ2luZy5cclxuICAgICAgICBkcmFnZ2luZyA9IHRydWU7XHJcbiAgICAgICAgaW52b2tlKFwid2FpbHM6ZHJhZ1wiKTtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcclxuICAgICAgICAvLyBFaXRoZXIgZHJhZyBvciByZXNpemUgaXMgb25nb2luZyxcclxuICAgICAgICAvLyByZXNldCByZWFkaW5lc3MgYW5kIHN0b3AgcHJvY2Vzc2luZy5cclxuICAgICAgICBjYW5EcmFnID0gY2FuUmVzaXplID0gZmFsc2U7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICghcmVzaXphYmxlIHx8ICFJc1dpbmRvd3MoKSkge1xyXG4gICAgICAgIGlmIChyZXNpemVFZGdlKSB7IHNldFJlc2l6ZSgpOyB9XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZUhlaWdodCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlSGVpZ2h0XCIpIHx8IDU7XHJcbiAgICBjb25zdCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcclxuXHJcbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXMuXHJcbiAgICBjb25zdCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcclxuXHJcbiAgICBjb25zdCByaWdodEJvcmRlciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgcmVzaXplSGFuZGxlV2lkdGg7XHJcbiAgICBjb25zdCBsZWZ0Qm9yZGVyID0gZXZlbnQuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgY29uc3QgdG9wQm9yZGVyID0gZXZlbnQuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcclxuICAgIGNvbnN0IGJvdHRvbUJvcmRlciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcclxuXHJcbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lciBhcmVhcy5cclxuICAgIGNvbnN0IHJpZ2h0Q29ybmVyID0gKHdpbmRvdy5vdXRlcldpZHRoIC0gZXZlbnQuY2xpZW50WCkgPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBjb25zdCBsZWZ0Q29ybmVyID0gZXZlbnQuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGNvbnN0IHRvcENvcm5lciA9IGV2ZW50LmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG4gICAgY29uc3QgYm90dG9tQ29ybmVyID0gKHdpbmRvdy5vdXRlckhlaWdodCAtIGV2ZW50LmNsaWVudFkpIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcclxuXHJcbiAgICBpZiAoIWxlZnRDb3JuZXIgJiYgIXRvcENvcm5lciAmJiAhYm90dG9tQ29ybmVyICYmICFyaWdodENvcm5lcikge1xyXG4gICAgICAgIC8vIE9wdGltaXNhdGlvbjogb3V0IG9mIGFsbCBjb3JuZXIgYXJlYXMgaW1wbGllcyBvdXQgb2YgYm9yZGVycy5cclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH1cclxuICAgIC8vIERldGVjdCBjb3JuZXJzLlxyXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgLy8gRGV0ZWN0IGJvcmRlcnMuXHJcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xyXG4gICAgLy8gT3V0IG9mIGJvcmRlciBhcmVhLlxyXG4gICAgZWxzZSBzZXRSZXNpemUoKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcclxuXHJcbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xyXG5jb25zdCBTaG93TWV0aG9kID0gMTtcclxuY29uc3QgUXVpdE1ldGhvZCA9IDI7XHJcblxyXG4vKipcclxuICogSGlkZXMgYSBjZXJ0YWluIG1ldGhvZCBieSBjYWxsaW5nIHRoZSBIaWRlTWV0aG9kIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChIaWRlTWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgQ2FuY2VsbGFibGVQcm9taXNlLCB0eXBlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnMgfSBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XHJcblxyXG4vLyBTZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5jYWxsUmVzdWx0SGFuZGxlciA9IHJlc3VsdEhhbmRsZXI7XHJcbndpbmRvdy5fd2FpbHMuY2FsbEVycm9ySGFuZGxlciA9IGVycm9ySGFuZGxlcjtcclxuXHJcbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8Q2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIiB8IFwib25jYW5jZWxsZWRcIj5cclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbGwpO1xyXG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYW5jZWxDYWxsKTtcclxuY29uc3QgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xyXG5cclxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xyXG5jb25zdCBDYW5jZWxNZXRob2QgPSAwXHJcblxyXG4vKipcclxuICogSG9sZHMgYWxsIHJlcXVpcmVkIGluZm9ybWF0aW9uIGZvciBhIGJpbmRpbmcgY2FsbC5cclxuICogTWF5IHByb3ZpZGUgZWl0aGVyIGEgbWV0aG9kIElEIG9yIGEgbWV0aG9kIG5hbWUsIGJ1dCBub3QgYm90aC5cclxuICovXHJcbmV4cG9ydCB0eXBlIENhbGxPcHRpb25zID0ge1xyXG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZElEOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZE5hbWU/OiBuZXZlcjtcclxuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cclxuICAgIGFyZ3M6IGFueVtdO1xyXG59IHwge1xyXG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZElEPzogbmV2ZXI7XHJcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZE5hbWU6IHN0cmluZztcclxuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cclxuICAgIGFyZ3M6IGFueVtdO1xyXG59O1xyXG5cclxuLyoqXHJcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdGhyb3duIGluIGNhc2UgdGhlIGJvdW5kIG1ldGhvZCByZXR1cm5zIGFuIGVycm9yLlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBSdW50aW1lRXJyb3IjbmFtZX0gcHJvcGVydHkgaXMgXCJSdW50aW1lRXJyb3JcIi5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBSdW50aW1lRXJyb3IgZXh0ZW5kcyBFcnJvciB7XHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgUnVudGltZUVycm9yIGluc3RhbmNlLlxyXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XHJcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJSdW50aW1lRXJyb3JcIjtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIHJlc3VsdCBvZiBhIGNhbGwgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cclxuICogQHBhcmFtIGRhdGEgLSBUaGUgcmVzdWx0IGRhdGEgb2YgdGhlIHJlcXVlc3QuXHJcbiAqIEBwYXJhbSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cclxuICovXHJcbmZ1bmN0aW9uIHJlc3VsdEhhbmRsZXIoaWQ6IHN0cmluZywgZGF0YTogc3RyaW5nLCBpc0pTT046IGJvb2xlYW4pOiB2b2lkIHtcclxuICAgIGNvbnN0IHJlc29sdmVycyA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcclxuICAgIGlmICghcmVzb2x2ZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICghZGF0YSkge1xyXG4gICAgICAgIHJlc29sdmVycy5yZXNvbHZlKHVuZGVmaW5lZCk7XHJcbiAgICB9IGVsc2UgaWYgKCFpc0pTT04pIHtcclxuICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShkYXRhKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyOiBhbnkpIHtcclxuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgVHlwZUVycm9yKFwiY291bGQgbm90IHBhcnNlIHJlc3VsdDogXCIgKyBlcnIubWVzc2FnZSwgeyBjYXVzZTogZXJyIH0pKTtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIHRoZSBlcnJvciBmcm9tIGEgY2FsbCByZXF1ZXN0LlxyXG4gKlxyXG4gKiBAcGFyYW0gaWQgLSBUaGUgaWQgb2YgdGhlIHByb21pc2UgaGFuZGxlci5cclxuICogQHBhcmFtIGRhdGEgLSBUaGUgZXJyb3IgZGF0YSB0byByZWplY3QgdGhlIHByb21pc2UgaGFuZGxlciB3aXRoLlxyXG4gKiBAcGFyYW0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXHJcbiAqL1xyXG5mdW5jdGlvbiBlcnJvckhhbmRsZXIoaWQ6IHN0cmluZywgZGF0YTogc3RyaW5nLCBpc0pTT046IGJvb2xlYW4pOiB2b2lkIHtcclxuICAgIGNvbnN0IHJlc29sdmVycyA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcclxuICAgIGlmICghcmVzb2x2ZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICghaXNKU09OKSB7XHJcbiAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgRXJyb3IoZGF0YSkpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBsZXQgZXJyb3I6IGFueTtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBlcnJvciA9IEpTT04ucGFyc2UoZGF0YSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyOiBhbnkpIHtcclxuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgVHlwZUVycm9yKFwiY291bGQgbm90IHBhcnNlIGVycm9yOiBcIiArIGVyci5tZXNzYWdlLCB7IGNhdXNlOiBlcnIgfSkpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBsZXQgb3B0aW9uczogRXJyb3JPcHRpb25zID0ge307XHJcbiAgICAgICAgaWYgKGVycm9yLmNhdXNlKSB7XHJcbiAgICAgICAgICAgIG9wdGlvbnMuY2F1c2UgPSBlcnJvci5jYXVzZTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGxldCBleGNlcHRpb247XHJcbiAgICAgICAgc3dpdGNoIChlcnJvci5raW5kKSB7XHJcbiAgICAgICAgICAgIGNhc2UgXCJSZWZlcmVuY2VFcnJvclwiOlxyXG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFJlZmVyZW5jZUVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xyXG4gICAgICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgICAgIGNhc2UgXCJUeXBlRXJyb3JcIjpcclxuICAgICAgICAgICAgICAgIGV4Y2VwdGlvbiA9IG5ldyBUeXBlRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuICAgICAgICAgICAgY2FzZSBcIlJ1bnRpbWVFcnJvclwiOlxyXG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFJ1bnRpbWVFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICBkZWZhdWx0OlxyXG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IEVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xyXG4gICAgICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICByZXNvbHZlcnMucmVqZWN0KGV4Y2VwdGlvbik7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgY2FsbFJlc3BvbnNlcyBtYXAuXHJcbiAqXHJcbiAqIEBwYXJhbSBpZCAtIFRoZSBJRCBvZiB0aGUgcmVzcG9uc2UgdG8gYmUgcmV0cmlldmVkIGFuZCByZW1vdmVkLlxyXG4gKiBAcmV0dXJucyBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQsIGlmIGFueS5cclxuICovXHJcbmZ1bmN0aW9uIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkOiBzdHJpbmcpOiBQcm9taXNlUmVzb2x2ZXJzIHwgdW5kZWZpbmVkIHtcclxuICAgIGNvbnN0IHJlc3BvbnNlID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgcmV0dXJuIHJlc3BvbnNlO1xyXG59XHJcblxyXG4vKipcclxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cclxuICpcclxuICogQHJldHVybnMgQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgY2FsbFJlc3BvbnNlcyBzZXQuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCk6IHN0cmluZyB7XHJcbiAgICBsZXQgcmVzdWx0O1xyXG4gICAgZG8ge1xyXG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xyXG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbCBhIGJvdW5kIG1ldGhvZCBhY2NvcmRpbmcgdG8gdGhlIGdpdmVuIGNhbGwgb3B0aW9ucy5cclxuICpcclxuICogSW4gY2FzZSBvZiBmYWlsdXJlLCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGFuIGV4Y2VwdGlvblxyXG4gKiBhbW9uZyBSZWZlcmVuY2VFcnJvciAodW5rbm93biBtZXRob2QpLCBUeXBlRXJyb3IgKHdyb25nIGFyZ3VtZW50IGNvdW50IG9yIHR5cGUpLFxyXG4gKiB7QGxpbmsgUnVudGltZUVycm9yfSAobWV0aG9kIHJldHVybmVkIGFuIGVycm9yKSwgb3Igb3RoZXIgKG5ldHdvcmsgb3IgaW50ZXJuYWwgZXJyb3JzKS5cclxuICogVGhlIGV4Y2VwdGlvbiBtaWdodCBoYXZlIGEgXCJjYXVzZVwiIGZpZWxkIHdpdGggdGhlIHZhbHVlIHJldHVybmVkXHJcbiAqIGJ5IHRoZSBhcHBsaWNhdGlvbi0gb3Igc2VydmljZS1sZXZlbCBlcnJvciBtYXJzaGFsaW5nIGZ1bmN0aW9ucy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBBIG1ldGhvZCBjYWxsIGRlc2NyaXB0b3IuXHJcbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIGNhbGwuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zOiBDYWxsT3B0aW9ucyk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcclxuICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG5cclxuICAgIGNvbnN0IHJlc3VsdCA9IENhbmNlbGxhYmxlUHJvbWlzZS53aXRoUmVzb2x2ZXJzPGFueT4oKTtcclxuICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmU6IHJlc3VsdC5yZXNvbHZlLCByZWplY3Q6IHJlc3VsdC5yZWplY3QgfSk7XHJcblxyXG4gICAgY29uc3QgcmVxdWVzdCA9IGNhbGwoQ2FsbEJpbmRpbmcsIE9iamVjdC5hc3NpZ24oeyBcImNhbGwtaWRcIjogaWQgfSwgb3B0aW9ucykpO1xyXG4gICAgbGV0IHJ1bm5pbmcgPSBmYWxzZTtcclxuXHJcbiAgICByZXF1ZXN0LnRoZW4oKCkgPT4ge1xyXG4gICAgICAgIHJ1bm5pbmcgPSB0cnVlO1xyXG4gICAgfSwgKGVycikgPT4ge1xyXG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICByZXN1bHQucmVqZWN0KGVycik7XHJcbiAgICB9KTtcclxuXHJcbiAgICBjb25zdCBjYW5jZWwgPSAoKSA9PiB7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIHJldHVybiBjYW5jZWxDYWxsKENhbmNlbE1ldGhvZCwge1wiY2FsbC1pZFwiOiBpZH0pLmNhdGNoKChlcnIpID0+IHtcclxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIHdoaWxlIHJlcXVlc3RpbmcgYmluZGluZyBjYWxsIGNhbmNlbGxhdGlvbjpcIiwgZXJyKTtcclxuICAgICAgICB9KTtcclxuICAgIH07XHJcblxyXG4gICAgcmVzdWx0Lm9uY2FuY2VsbGVkID0gKCkgPT4ge1xyXG4gICAgICAgIGlmIChydW5uaW5nKSB7XHJcbiAgICAgICAgICAgIHJldHVybiBjYW5jZWwoKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICByZXR1cm4gcmVxdWVzdC50aGVuKGNhbmNlbCk7XHJcbiAgICAgICAgfVxyXG4gICAgfTtcclxuXHJcbiAgICByZXR1cm4gcmVzdWx0LnByb21pc2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyBhIGJvdW5kIG1ldGhvZCBieSBuYW1lIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXHJcbiAqXHJcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxyXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeU5hbWUobWV0aG9kTmFtZTogc3RyaW5nLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcclxuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kTmFtZSwgYXJncyB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIGJ5IGl0cyBudW1lcmljIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXHJcbiAqXHJcbiAqIEBwYXJhbSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm4gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnlJRChtZXRob2RJRDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcclxuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kSUQsIGFyZ3MgfSk7XHJcbn1cclxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2luc3BlY3QtanMvaXMtY2FsbGFibGVcclxuXHJcbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxyXG4vL1xyXG4vLyBDb3B5cmlnaHQgKGMpIDIwMTUgSm9yZGFuIEhhcmJhbmRcclxuLy9cclxuLy8gUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weVxyXG4vLyBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsXHJcbi8vIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHNcclxuLy8gdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbFxyXG4vLyBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXNcclxuLy8gZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcclxuLy9cclxuLy8gVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXHJcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXHJcbi8vXHJcbi8vIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1JcclxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksXHJcbi8vIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRVxyXG4vLyBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSXHJcbi8vIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sXHJcbi8vIE9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFXHJcbi8vIFNPRlRXQVJFLlxyXG5cclxudmFyIGZuVG9TdHIgPSBGdW5jdGlvbi5wcm90b3R5cGUudG9TdHJpbmc7XHJcbnZhciByZWZsZWN0QXBwbHk6IHR5cGVvZiBSZWZsZWN0LmFwcGx5IHwgZmFsc2UgfCBudWxsID0gdHlwZW9mIFJlZmxlY3QgPT09ICdvYmplY3QnICYmIFJlZmxlY3QgIT09IG51bGwgJiYgUmVmbGVjdC5hcHBseTtcclxudmFyIGJhZEFycmF5TGlrZTogYW55O1xyXG52YXIgaXNDYWxsYWJsZU1hcmtlcjogYW55O1xyXG5pZiAodHlwZW9mIHJlZmxlY3RBcHBseSA9PT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgT2JqZWN0LmRlZmluZVByb3BlcnR5ID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGJhZEFycmF5TGlrZSA9IE9iamVjdC5kZWZpbmVQcm9wZXJ0eSh7fSwgJ2xlbmd0aCcsIHtcclxuICAgICAgICAgICAgZ2V0OiBmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgICAgICB0aHJvdyBpc0NhbGxhYmxlTWFya2VyO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgaXNDYWxsYWJsZU1hcmtlciA9IHt9O1xyXG4gICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZSBuby10aHJvdy1saXRlcmFsXHJcbiAgICAgICAgcmVmbGVjdEFwcGx5KGZ1bmN0aW9uICgpIHsgdGhyb3cgNDI7IH0sIG51bGwsIGJhZEFycmF5TGlrZSk7XHJcbiAgICB9IGNhdGNoIChfKSB7XHJcbiAgICAgICAgaWYgKF8gIT09IGlzQ2FsbGFibGVNYXJrZXIpIHtcclxuICAgICAgICAgICAgcmVmbGVjdEFwcGx5ID0gbnVsbDtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbn0gZWxzZSB7XHJcbiAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xyXG59XHJcblxyXG52YXIgY29uc3RydWN0b3JSZWdleCA9IC9eXFxzKmNsYXNzXFxiLztcclxudmFyIGlzRVM2Q2xhc3NGbiA9IGZ1bmN0aW9uIGlzRVM2Q2xhc3NGdW5jdGlvbih2YWx1ZTogYW55KTogYm9vbGVhbiB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIHZhciBmblN0ciA9IGZuVG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICAgICAgcmV0dXJuIGNvbnN0cnVjdG9yUmVnZXgudGVzdChmblN0cik7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlOyAvLyBub3QgYSBmdW5jdGlvblxyXG4gICAgfVxyXG59O1xyXG5cclxudmFyIHRyeUZ1bmN0aW9uT2JqZWN0ID0gZnVuY3Rpb24gdHJ5RnVuY3Rpb25Ub1N0cih2YWx1ZTogYW55KTogYm9vbGVhbiB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgICAgIGZuVG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICAgICAgcmV0dXJuIHRydWU7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG59O1xyXG52YXIgdG9TdHIgPSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nO1xyXG52YXIgb2JqZWN0Q2xhc3MgPSAnW29iamVjdCBPYmplY3RdJztcclxudmFyIGZuQ2xhc3MgPSAnW29iamVjdCBGdW5jdGlvbl0nO1xyXG52YXIgZ2VuQ2xhc3MgPSAnW29iamVjdCBHZW5lcmF0b3JGdW5jdGlvbl0nO1xyXG52YXIgZGRhQ2xhc3MgPSAnW29iamVjdCBIVE1MQWxsQ29sbGVjdGlvbl0nOyAvLyBJRSAxMVxyXG52YXIgZGRhQ2xhc3MyID0gJ1tvYmplY3QgSFRNTCBkb2N1bWVudC5hbGwgY2xhc3NdJztcclxudmFyIGRkYUNsYXNzMyA9ICdbb2JqZWN0IEhUTUxDb2xsZWN0aW9uXSc7IC8vIElFIDktMTBcclxudmFyIGhhc1RvU3RyaW5nVGFnID0gdHlwZW9mIFN5bWJvbCA9PT0gJ2Z1bmN0aW9uJyAmJiAhIVN5bWJvbC50b1N0cmluZ1RhZzsgLy8gYmV0dGVyOiB1c2UgYGhhcy10b3N0cmluZ3RhZ2BcclxuXHJcbnZhciBpc0lFNjggPSAhKDAgaW4gWyxdKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBuby1zcGFyc2UtYXJyYXlzLCBjb21tYS1zcGFjaW5nXHJcblxyXG52YXIgaXNEREE6ICh2YWx1ZTogYW55KSA9PiBib29sZWFuID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCgpIHsgcmV0dXJuIGZhbHNlOyB9O1xyXG5pZiAodHlwZW9mIGRvY3VtZW50ID09PSAnb2JqZWN0Jykge1xyXG4gICAgLy8gRmlyZWZveCAzIGNhbm9uaWNhbGl6ZXMgRERBIHRvIHVuZGVmaW5lZCB3aGVuIGl0J3Mgbm90IGFjY2Vzc2VkIGRpcmVjdGx5XHJcbiAgICB2YXIgYWxsID0gZG9jdW1lbnQuYWxsO1xyXG4gICAgaWYgKHRvU3RyLmNhbGwoYWxsKSA9PT0gdG9TdHIuY2FsbChkb2N1bWVudC5hbGwpKSB7XHJcbiAgICAgICAgaXNEREEgPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKHZhbHVlKSB7XHJcbiAgICAgICAgICAgIC8qIGdsb2JhbHMgZG9jdW1lbnQ6IGZhbHNlICovXHJcbiAgICAgICAgICAgIC8vIGluIElFIDYtOCwgdHlwZW9mIGRvY3VtZW50LmFsbCBpcyBcIm9iamVjdFwiIGFuZCBpdCdzIHRydXRoeVxyXG4gICAgICAgICAgICBpZiAoKGlzSUU2OCB8fCAhdmFsdWUpICYmICh0eXBlb2YgdmFsdWUgPT09ICd1bmRlZmluZWQnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcpKSB7XHJcbiAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgIHZhciBzdHIgPSB0b1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICAgICAgICAgICAgICByZXR1cm4gKFxyXG4gICAgICAgICAgICAgICAgICAgICAgICBzdHIgPT09IGRkYUNsYXNzXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MyXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MzIC8vIG9wZXJhIDEyLjE2XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gb2JqZWN0Q2xhc3MgLy8gSUUgNi04XHJcbiAgICAgICAgICAgICAgICAgICAgKSAmJiB2YWx1ZSgnJykgPT0gbnVsbDsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBlcWVxZXFcclxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGUpIHsgLyoqLyB9XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGlzQ2FsbGFibGVSZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSAge1xyXG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxyXG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIHRyeSB7XHJcbiAgICAgICAgKHJlZmxlY3RBcHBseSBhcyBhbnkpKHZhbHVlLCBudWxsLCBiYWRBcnJheUxpa2UpO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIGlmIChlICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgfVxyXG4gICAgcmV0dXJuICFpc0VTNkNsYXNzRm4odmFsdWUpICYmIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcclxufVxyXG5cclxuZnVuY3Rpb24gaXNDYWxsYWJsZU5vUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkge1xyXG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxyXG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIGlmIChoYXNUb1N0cmluZ1RhZykgeyByZXR1cm4gdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpOyB9XHJcbiAgICBpZiAoaXNFUzZDbGFzc0ZuKHZhbHVlKSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIHZhciBzdHJDbGFzcyA9IHRvU3RyLmNhbGwodmFsdWUpO1xyXG4gICAgaWYgKHN0ckNsYXNzICE9PSBmbkNsYXNzICYmIHN0ckNsYXNzICE9PSBnZW5DbGFzcyAmJiAhKC9eXFxbb2JqZWN0IEhUTUwvKS50ZXN0KHN0ckNsYXNzKSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7XHJcbn07XHJcblxyXG5leHBvcnQgZGVmYXVsdCByZWZsZWN0QXBwbHkgPyBpc0NhbGxhYmxlUmVmQXBwbHkgOiBpc0NhbGxhYmxlTm9SZWZBcHBseTtcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCBpc0NhbGxhYmxlIGZyb20gXCIuL2NhbGxhYmxlLmpzXCI7XHJcblxyXG4vKipcclxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB1c2VkIGFzIHJlamVjdGlvbiByZWFzb25cclxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IGlzIGNhbmNlbGxlZCBzdWNjZXNzZnVsbHkuXHJcbiAqXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsRXJyb3JcImAuXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgY2F1c2UgcGFzc2VkIHRvIHRoZSBjYW5jZWwgbWV0aG9kLCBpZiBhbnkuXHJcbiAqL1xyXG5leHBvcnQgY2xhc3MgQ2FuY2VsRXJyb3IgZXh0ZW5kcyBFcnJvciB7XHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbEVycm9yYCBpbnN0YW5jZS5cclxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXHJcbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xyXG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xyXG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsRXJyb3JcIjtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvblxyXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gcmVqZWN0cyBhZnRlciBiZWluZyBjYW5jZWxsZWQsXHJcbiAqIG9yIHdoZW4gdGhlIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgdGhyb3dzIG9yIHJlamVjdHMuXHJcbiAqXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcImAuXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgcmVhc29uIHRoZSBwcm9taXNlIHJlamVjdGVkIHdpdGguXHJcbiAqXHJcbiAqIEJlY2F1c2UgdGhlIG9yaWdpbmFsIHByb21pc2Ugd2FzIGNhbmNlbGxlZCxcclxuICogYSB3cmFwcGVyIHByb21pc2Ugd2lsbCBiZSBwYXNzZWQgdG8gdGhlIHVuaGFuZGxlZCByZWplY3Rpb24gbGlzdGVuZXIgaW5zdGVhZC5cclxuICogVGhlIHtAbGluayBwcm9taXNlfSBwcm9wZXJ0eSBob2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgb3JpZ2luYWwgcHJvbWlzZS5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvciBleHRlbmRzIEVycm9yIHtcclxuICAgIC8qKlxyXG4gICAgICogSG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIHByb21pc2UgdGhhdCB3YXMgY2FuY2VsbGVkIGFuZCB0aGVuIHJlamVjdGVkLlxyXG4gICAgICovXHJcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj47XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcmAgaW5zdGFuY2UuXHJcbiAgICAgKiBAcGFyYW0gcHJvbWlzZSAtIFRoZSBwcm9taXNlIHRoYXQgY2F1c2VkIHRoZSBlcnJvciBvcmlnaW5hbGx5LlxyXG4gICAgICogQHBhcmFtIHJlYXNvbiAtIFRoZSByZWplY3Rpb24gcmVhc29uLlxyXG4gICAgICogQHBhcmFtIGluZm8gLSBBbiBvcHRpb25hbCBpbmZvcm1hdGl2ZSBtZXNzYWdlIHNwZWNpZnlpbmcgdGhlIGNpcmN1bXN0YW5jZXMgaW4gd2hpY2ggdGhlIGVycm9yIHdhcyB0aHJvd24uXHJcbiAgICAgKiAgICAgICAgICAgICAgIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgYFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cImAuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgcmVhc29uPzogYW55LCBpbmZvPzogc3RyaW5nKSB7XHJcbiAgICAgICAgc3VwZXIoKGluZm8gPz8gXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiKSArIFwiIFJlYXNvbjogXCIgKyBlcnJvck1lc3NhZ2UocmVhc29uKSwgeyBjYXVzZTogcmVhc29uIH0pO1xyXG4gICAgICAgIHRoaXMucHJvbWlzZSA9IHByb21pc2U7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiO1xyXG4gICAgfVxyXG59XHJcblxyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+ID0gKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yID0gKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcclxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgPSAoY2F1c2U/OiBhbnkpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiA9IChyZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiwgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcikgPT4gdm9pZDtcclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XHJcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwpOiBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxIHwgVFJlc3VsdDI+O1xyXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG59XHJcblxyXG4vKipcclxuICogV3JhcHMgYSBjYW5jZWxsYWJsZSBwcm9taXNlIGFsb25nIHdpdGggaXRzIHJlc29sdXRpb24gbWV0aG9kcy5cclxuICogVGhlIGBvbmNhbmNlbGxlZGAgZmllbGQgd2lsbCBiZSBudWxsIGluaXRpYWxseSBidXQgbWF5IGJlIHNldCB0byBwcm92aWRlIGEgY3VzdG9tIGNhbmNlbGxhdGlvbiBmdW5jdGlvbi5cclxuICovXHJcbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XHJcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XHJcbiAgICByZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPjtcclxuICAgIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3I7XHJcbiAgICBvbmNhbmNlbGxlZDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyIHwgbnVsbDtcclxufVxyXG5cclxuaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlIHtcclxuICAgIHJlYWRvbmx5IHJvb3Q6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlO1xyXG4gICAgcmVzb2x2aW5nOiBib29sZWFuO1xyXG4gICAgc2V0dGxlZDogYm9vbGVhbjtcclxuICAgIHJlYXNvbj86IENhbmNlbEVycm9yO1xyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBiYXJyaWVyU3ltID0gU3ltYm9sKFwiYmFycmllclwiKTtcclxuY29uc3QgY2FuY2VsSW1wbFN5bSA9IFN5bWJvbChcImNhbmNlbEltcGxcIik7XHJcbmNvbnN0IHNwZWNpZXMgPSBTeW1ib2wuc3BlY2llcyA/PyBTeW1ib2woXCJzcGVjaWVzUG9seWZpbGxcIik7XHJcblxyXG4vKipcclxuICogQSBwcm9taXNlIHdpdGggYW4gYXR0YWNoZWQgbWV0aG9kIGZvciBjYW5jZWxsaW5nIGxvbmctcnVubmluZyBvcGVyYXRpb25zIChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWx9KS5cclxuICogQ2FuY2VsbGF0aW9uIGNhbiBvcHRpb25hbGx5IGJlIGJvdW5kIHRvIGFuIHtAbGluayBBYm9ydFNpZ25hbH1cclxuICogZm9yIGJldHRlciBjb21wb3NhYmlsaXR5IChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWxPbn0pLlxyXG4gKlxyXG4gKiBDYW5jZWxsaW5nIGEgcGVuZGluZyBwcm9taXNlIHdpbGwgcmVzdWx0IGluIGFuIGltbWVkaWF0ZSByZWplY3Rpb25cclxuICogd2l0aCBhbiBpbnN0YW5jZSBvZiB7QGxpbmsgQ2FuY2VsRXJyb3J9IGFzIHJlYXNvbixcclxuICogYnV0IHdob2V2ZXIgc3RhcnRlZCB0aGUgcHJvbWlzZSB3aWxsIGJlIHJlc3BvbnNpYmxlXHJcbiAqIGZvciBhY3R1YWxseSBhYm9ydGluZyB0aGUgdW5kZXJseWluZyBvcGVyYXRpb24uXHJcbiAqIFRvIHRoaXMgcHVycG9zZSwgdGhlIGNvbnN0cnVjdG9yIGFuZCBhbGwgY2hhaW5pbmcgbWV0aG9kc1xyXG4gKiBhY2NlcHQgb3B0aW9uYWwgY2FuY2VsbGF0aW9uIGNhbGxiYWNrcy5cclxuICpcclxuICogSWYgYSBgQ2FuY2VsbGFibGVQcm9taXNlYCBzdGlsbCByZXNvbHZlcyBhZnRlciBoYXZpbmcgYmVlbiBjYW5jZWxsZWQsXHJcbiAqIHRoZSByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWQuIElmIGl0IHJlamVjdHMsIHRoZSByZWFzb25cclxuICogd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLFxyXG4gKiB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBpbnN0YW5jZS5cclxuICogVG8gZmFjaWxpdGF0ZSB0aGUgaGFuZGxpbmcgb2YgY2FuY2VsbGF0aW9uIHJlcXVlc3RzLFxyXG4gKiBjYW5jZWxsZWQgYENhbmNlbGxhYmxlUHJvbWlzZWBzIHdpbGwgX25vdF8gcmVwb3J0IHVuaGFuZGxlZCBgQ2FuY2VsRXJyb3Jgc1xyXG4gKiB3aG9zZSBgY2F1c2VgIGZpZWxkIGlzIHRoZSBzYW1lIGFzIHRoZSBvbmUgd2l0aCB3aGljaCB0aGUgY3VycmVudCBwcm9taXNlIHdhcyBjYW5jZWxsZWQuXHJcbiAqXHJcbiAqIEFsbCB1c3VhbCBwcm9taXNlIG1ldGhvZHMgYXJlIGRlZmluZWQgYW5kIHJldHVybiBhIGBDYW5jZWxsYWJsZVByb21pc2VgXHJcbiAqIHdob3NlIGNhbmNlbCBtZXRob2Qgd2lsbCBjYW5jZWwgdGhlIHBhcmVudCBvcGVyYXRpb24gYXMgd2VsbCwgcHJvcGFnYXRpbmcgdGhlIGNhbmNlbGxhdGlvbiByZWFzb25cclxuICogdXB3YXJkcyB0aHJvdWdoIHByb21pc2UgY2hhaW5zLlxyXG4gKiBDb252ZXJzZWx5LCBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aWxsIG5vdCBhdXRvbWF0aWNhbGx5IGNhbmNlbCBkZXBlbmRlbnQgcHJvbWlzZXMgZG93bnN0cmVhbTpcclxuICogYGBgdHNcclxuICogbGV0IHJvb3QgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xyXG4gKiBsZXQgY2hpbGQxID0gcm9vdC50aGVuKCgpID0+IHsgLi4uIH0pO1xyXG4gKiBsZXQgY2hpbGQyID0gY2hpbGQxLnRoZW4oKCkgPT4geyAuLi4gfSk7XHJcbiAqIGxldCBjaGlsZDMgPSByb290LmNhdGNoKCgpID0+IHsgLi4uIH0pO1xyXG4gKiBjaGlsZDEuY2FuY2VsKCk7IC8vIENhbmNlbHMgY2hpbGQxIGFuZCByb290LCBidXQgbm90IGNoaWxkMiBvciBjaGlsZDNcclxuICogYGBgXHJcbiAqIENhbmNlbGxpbmcgYSBwcm9taXNlIHRoYXQgaGFzIGFscmVhZHkgc2V0dGxlZCBpcyBzYWZlIGFuZCBoYXMgbm8gY29uc2VxdWVuY2UuXHJcbiAqXHJcbiAqIFRoZSBgY2FuY2VsYCBtZXRob2QgcmV0dXJucyBhIHByb21pc2UgdGhhdCBfYWx3YXlzIGZ1bGZpbGxzX1xyXG4gKiBhZnRlciB0aGUgd2hvbGUgY2hhaW4gaGFzIHByb2Nlc3NlZCB0aGUgY2FuY2VsIHJlcXVlc3RcclxuICogYW5kIGFsbCBhdHRhY2hlZCBjYWxsYmFja3MgdXAgdG8gdGhhdCBtb21lbnQgaGF2ZSBydW4uXHJcbiAqXHJcbiAqIEFsbCBFUzIwMjQgcHJvbWlzZSBtZXRob2RzIChzdGF0aWMgYW5kIGluc3RhbmNlKSBhcmUgZGVmaW5lZCBvbiBDYW5jZWxsYWJsZVByb21pc2UsXHJcbiAqIGJ1dCBhY3R1YWwgYXZhaWxhYmlsaXR5IG1heSB2YXJ5IHdpdGggT1Mvd2VidmlldyB2ZXJzaW9uLlxyXG4gKlxyXG4gKiBJbiBsaW5lIHdpdGggdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmcsXHJcbiAqIGBDYW5jZWxsYWJsZVByb21pc2VgIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuXHJcbiAqIEV4dGVuZGVycyBzaG91bGQgdGFrZSBjYXJlIHRvIHByb3ZpZGUgdGhlaXIgb3duIG1ldGhvZCBpbXBsZW1lbnRhdGlvbnMuXHJcbiAqIFRoaXMgbWlnaHQgYmUgcmVjb25zaWRlcmVkIGluIGNhc2UgdGhlIHByb3Bvc2FsIGlzIHJldGlyZWQuXHJcbiAqXHJcbiAqIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBhIHdyYXBwZXIgYXJvdW5kIHRoZSBET00gUHJvbWlzZSBvYmplY3RcclxuICogYW5kIGlzIGNvbXBsaWFudCB3aXRoIHRoZSBbUHJvbWlzZXMvQSsgc3BlY2lmaWNhdGlvbl0oaHR0cHM6Ly9wcm9taXNlc2FwbHVzLmNvbS8pXHJcbiAqIChpdCBwYXNzZXMgdGhlIFtjb21wbGlhbmNlIHN1aXRlXShodHRwczovL2dpdGh1Yi5jb20vcHJvbWlzZXMtYXBsdXMvcHJvbWlzZXMtdGVzdHMpKVxyXG4gKiBpZiBzbyBpcyB0aGUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbi5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBDYW5jZWxsYWJsZVByb21pc2U8VD4gZXh0ZW5kcyBQcm9taXNlPFQ+IGltcGxlbWVudHMgUHJvbWlzZUxpa2U8VD4sIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xyXG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXHJcbiAgICAvKiogQGludGVybmFsICovXHJcbiAgICBwcml2YXRlIFtiYXJyaWVyU3ltXSE6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+IHwgbnVsbDtcclxuICAgIC8qKiBAaW50ZXJuYWwgKi9cclxuICAgIHByaXZhdGUgcmVhZG9ubHkgW2NhbmNlbEltcGxTeW1dITogKHJlYXNvbjogQ2FuY2VsRXJyb3IpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgYENhbmNlbGxhYmxlUHJvbWlzZWAuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGV4ZWN1dG9yIC0gQSBjYWxsYmFjayB1c2VkIHRvIGluaXRpYWxpemUgdGhlIHByb21pc2UuIFRoaXMgY2FsbGJhY2sgaXMgcGFzc2VkIHR3byBhcmd1bWVudHM6XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhIGByZXNvbHZlYCBjYWxsYmFjayB1c2VkIHRvIHJlc29sdmUgdGhlIHByb21pc2Ugd2l0aCBhIHZhbHVlXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBvciB0aGUgcmVzdWx0IG9mIGFub3RoZXIgcHJvbWlzZSAocG9zc2libHkgY2FuY2VsbGFibGUpLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGEgYHJlamVjdGAgY2FsbGJhY2sgdXNlZCB0byByZWplY3QgdGhlIHByb21pc2Ugd2l0aCBhIHByb3ZpZGVkIHJlYXNvbiBvciBlcnJvci5cclxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIHRoZSB2YWx1ZSBwcm92aWRlZCB0byB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAoaXQgaGFzIGEgYHRoZW5gIF9hbmRfIGEgYGNhbmNlbGAgbWV0aG9kKSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGZvcndhcmRlZCB0byB0aGF0IG9iamVjdCBhbmQgdGhlIG9uY2FuY2VsbGVkIHdpbGwgbm90IGJlIGludm9rZWQgYW55bW9yZS5cclxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIGFueSBvbmUgb2YgdGhlIHR3byBjYWxsYmFja3MgaXMgY2FsbGVkIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhlIHByb3ZpZGVkIHZhbHVlcyB3aWxsIGJlIGNhbmNlbGxlZCBhbmQgcmVzb2x2ZWQgYXMgdXN1YWwsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBidXQgdGhlaXIgcmVzdWx0cyB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqICAgICAgICAgICAgICAgICAgIEhvd2V2ZXIsIGlmIHRoZSByZXNvbHV0aW9uIHByb2Nlc3MgdWx0aW1hdGVseSBlbmRzIHVwIGluIGEgcmVqZWN0aW9uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGF0IGlzIG5vdCBkdWUgdG8gY2FuY2VsbGF0aW9uLCB0aGUgcmVqZWN0aW9uIHJlYXNvblxyXG4gICAgICogICAgICAgICAgICAgICAgICAgd2lsbCBiZSB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfVxyXG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqIEBwYXJhbSBvbmNhbmNlbGxlZCAtIEl0IGlzIHRoZSBjYWxsZXIncyByZXNwb25zaWJpbGl0eSB0byBlbnN1cmUgdGhhdCBhbnkgb3BlcmF0aW9uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBzdGFydGVkIGJ5IHRoZSBleGVjdXRvciBpcyBwcm9wZXJseSBoYWx0ZWQgdXBvbiBjYW5jZWxsYXRpb24uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBUaGlzIG9wdGlvbmFsIGNhbGxiYWNrIGNhbiBiZSB1c2VkIHRvIHRoYXQgcHVycG9zZS5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIEl0IHdpbGwgYmUgY2FsbGVkIF9zeW5jaHJvbm91c2x5XyB3aXRoIGEgY2FuY2VsbGF0aW9uIGNhdXNlXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aGVuIGNhbmNlbGxhdGlvbiBpcyByZXF1ZXN0ZWQsIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgcmVqZWN0ZWRcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9LCBidXQgX2JlZm9yZV9cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFueSB7QGxpbmsgdGhlbn0ve0BsaW5rIGNhdGNofS97QGxpbmsgZmluYWxseX0gY2FsbGJhY2sgcnVucy5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH1cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpbGwgb25seSBmdWxmaWxsIGFmdGVyIHRoZSBmb3JtZXIgaGFzIHNldHRsZWQuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBVbmhhbmRsZWQgZXhjZXB0aW9ucyBvciByZWplY3Rpb25zIGZyb20gdGhlIGNhbGxiYWNrIHdpbGwgYmUgd3JhcHBlZFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIHVuaGFuZGxlZCByZWplY3Rpb25zLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBjYWxsZWQgYmVmb3JlIGNhbmNlbGxhdGlvbiB3aXRoIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGlzIHByb21pc2Ugd2lsbCBiZSBkaXZlcnRlZCB0byB0aGF0IHByb21pc2UsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbmQgdGhlIG9yaWdpbmFsIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKGV4ZWN1dG9yOiBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpIHtcclxuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xyXG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xyXG4gICAgICAgIHN1cGVyKChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xyXG5cclxuICAgICAgICBpZiAoKHRoaXMuY29uc3RydWN0b3IgYXMgYW55KVtzcGVjaWVzXSAhPT0gUHJvbWlzZSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuIFBsZWFzZSByZWZyYWluIGZyb20gb3ZlcnJpZGluZyB0aGUgW1N5bWJvbC5zcGVjaWVzXSBzdGF0aWMgcHJvcGVydHkuXCIpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgbGV0IHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7XHJcbiAgICAgICAgICAgIHByb21pc2U6IHRoaXMsXHJcbiAgICAgICAgICAgIHJlc29sdmUsXHJcbiAgICAgICAgICAgIHJlamVjdCxcclxuICAgICAgICAgICAgZ2V0IG9uY2FuY2VsbGVkKCkgeyByZXR1cm4gb25jYW5jZWxsZWQgPz8gbnVsbDsgfSxcclxuICAgICAgICAgICAgc2V0IG9uY2FuY2VsbGVkKGNiKSB7IG9uY2FuY2VsbGVkID0gY2IgPz8gdW5kZWZpbmVkOyB9XHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgY29uc3Qgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xyXG4gICAgICAgICAgICBnZXQgcm9vdCgpIHsgcmV0dXJuIHN0YXRlOyB9LFxyXG4gICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxyXG4gICAgICAgICAgICBzZXR0bGVkOiBmYWxzZVxyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFNldHVwIGNhbmNlbGxhdGlvbiBzeXN0ZW0uXHJcbiAgICAgICAgdm9pZCBPYmplY3QuZGVmaW5lUHJvcGVydGllcyh0aGlzLCB7XHJcbiAgICAgICAgICAgIFtiYXJyaWVyU3ltXToge1xyXG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IHRydWUsXHJcbiAgICAgICAgICAgICAgICB2YWx1ZTogbnVsbFxyXG4gICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICBbY2FuY2VsSW1wbFN5bV06IHtcclxuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIHZhbHVlOiBjYW5jZWxsZXJGb3IocHJvbWlzZSwgc3RhdGUpXHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuXHJcbiAgICAgICAgLy8gUnVuIHRoZSBhY3R1YWwgZXhlY3V0b3IuXHJcbiAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBzdGF0ZSk7XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgZXhlY3V0b3IocmVzb2x2ZXJGb3IocHJvbWlzZSwgc3RhdGUpLCByZWplY3Rvcik7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHtcclxuICAgICAgICAgICAgICAgIGNvbnNvbGUubG9nKFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBDYW5jZWxsYWJsZVByb21pc2UgZXhlY3V0b3IuXCIsIGVycik7XHJcbiAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ2FuY2VscyBpbW1lZGlhdGVseSB0aGUgZXhlY3V0aW9uIG9mIHRoZSBvcGVyYXRpb24gYXNzb2NpYXRlZCB3aXRoIHRoaXMgcHJvbWlzZS5cclxuICAgICAqIFRoZSBwcm9taXNlIHJlamVjdHMgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0gaW5zdGFuY2UgYXMgcmVhc29uLFxyXG4gICAgICogd2l0aCB0aGUge0BsaW5rIENhbmNlbEVycm9yI2NhdXNlfSBwcm9wZXJ0eSBzZXQgdG8gdGhlIGdpdmVuIGFyZ3VtZW50LCBpZiBhbnkuXHJcbiAgICAgKlxyXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgYWZ0ZXIgdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZDtcclxuICAgICAqIHJlcGVhdGVkIGNhbGxzIGluIHBhcnRpY3VsYXIgYXJlIHNhZmUsIGJ1dCBvbmx5IHRoZSBmaXJzdCBvbmVcclxuICAgICAqIHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIGBDYW5jZWxFcnJvcmAgZXhjZXB0aW9uIF9uZWVkIG5vdF8gYmUgaGFuZGxlZCBleHBsaWNpdGx5IF9vbiB0aGUgcHJvbWlzZXMgdGhhdCBhcmUgYmVpbmcgY2FuY2VsbGVkOl9cclxuICAgICAqIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpdGggbm8gYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXIgZG9lcyBub3QgdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LlxyXG4gICAgICogVGhlcmVmb3JlLCB0aGUgZm9sbG93aW5nIGlkaW9tcyBhcmUgYWxsIGVxdWFsbHkgY29ycmVjdDpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLmNhbmNlbCgpO1xyXG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2FuY2VsKCk7XHJcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYXRjaCguLi4pLmNhbmNlbCgpO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBXaGVuZXZlciBzb21lIGNhbmNlbGxlZCBwcm9taXNlIGluIGEgY2hhaW4gcmVqZWN0cyB3aXRoIGEgYENhbmNlbEVycm9yYFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBjYW5jZWxsYXRpb24gY2F1c2UgYXMgaXRzZWxmLCB0aGUgZXJyb3Igd2lsbCBiZSBkaXNjYXJkZWQgc2lsZW50bHkuXHJcbiAgICAgKiBIb3dldmVyLCB0aGUgYENhbmNlbEVycm9yYCBfd2lsbCBzdGlsbCBiZSBkZWxpdmVyZWRfIHRvIGFsbCBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlcnNcclxuICAgICAqIGFkZGVkIGJ5IHtAbGluayB0aGVufSBhbmQgcmVsYXRlZCBtZXRob2RzOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XHJcbiAgICAgKiBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLmNhdGNoKGNvbnNvbGUubG9nKTtcclxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBBIENhbmNlbEVycm9yIGlzIHByaW50ZWQgdG8gdGhlIGNvbnNvbGUuXHJcbiAgICAgKiBgYGBcclxuICAgICAqIElmIHRoZSBgQ2FuY2VsRXJyb3JgIGlzIG5vdCBoYW5kbGVkIGRvd25zdHJlYW0gYnkgdGhlIHRpbWUgaXQgcmVhY2hlc1xyXG4gICAgICogYSBfbm9uLWNhbmNlbGxlZF8gcHJvbWlzZSwgaXQgX3dpbGxfIHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCxcclxuICAgICAqIGp1c3QgbGlrZSBub3JtYWwgcmVqZWN0aW9ucyB3b3VsZDpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xyXG4gICAgICogbGV0IGNoYWluZWQgPSBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLnRoZW4oKCkgPT4geyAuLi4gfSk7IC8vIE5vIGNhdGNoLi4uXHJcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gVW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCBvbiBjaGFpbmVkIVxyXG4gICAgICogYGBgXHJcbiAgICAgKiBUaGVyZWZvcmUsIGl0IGlzIGltcG9ydGFudCB0byBlaXRoZXIgY2FuY2VsIHdob2xlIHByb21pc2UgY2hhaW5zIGZyb20gdGhlaXIgdGFpbCxcclxuICAgICAqIGFzIHNob3duIGluIHRoZSBjb3JyZWN0IGlkaW9tcyBhYm92ZSwgb3IgdGFrZSBjYXJlIG9mIGhhbmRsaW5nIGVycm9ycyBldmVyeXdoZXJlLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIEEgY2FuY2VsbGFibGUgcHJvbWlzZSB0aGF0IF9mdWxmaWxsc18gYWZ0ZXIgdGhlIGNhbmNlbCBjYWxsYmFjayAoaWYgYW55KVxyXG4gICAgICogYW5kIGFsbCBoYW5kbGVycyBhdHRhY2hlZCB1cCB0byB0aGUgY2FsbCB0byBjYW5jZWwgaGF2ZSBydW4uXHJcbiAgICAgKiBJZiB0aGUgY2FuY2VsIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgYnkgYGNhbmNlbGBcclxuICAgICAqIHdpbGwgYWxzbyB3YWl0IGZvciB0aGF0IHRoZW5hYmxlIHRvIHNldHRsZS5cclxuICAgICAqIFRoaXMgZW5hYmxlcyBjYWxsZXJzIHRvIHdhaXQgZm9yIHRoZSBjYW5jZWxsZWQgb3BlcmF0aW9uIHRvIHRlcm1pbmF0ZVxyXG4gICAgICogd2l0aG91dCBiZWluZyBmb3JjZWQgdG8gaGFuZGxlIHBvdGVudGlhbCBlcnJvcnMgYXQgdGhlIGNhbGwgc2l0ZS5cclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKS50aGVuKCgpID0+IHtcclxuICAgICAqICAgICAvLyBDbGVhbnVwIGZpbmlzaGVkLCBpdCdzIHNhZmUgdG8gZG8gc29tZXRoaW5nIGVsc2UuXHJcbiAgICAgKiB9LCAoZXJyKSA9PiB7XHJcbiAgICAgKiAgICAgLy8gVW5yZWFjaGFibGU6IHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gY2FuY2VsIHdpbGwgbmV2ZXIgcmVqZWN0LlxyXG4gICAgICogfSk7XHJcbiAgICAgKiBgYGBcclxuICAgICAqIE5vdGUgdGhhdCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIF9ub3RfIGhhbmRsZSBpbXBsaWNpdGx5IGFueSByZWplY3Rpb25cclxuICAgICAqIHRoYXQgbWlnaHQgaGF2ZSBvY2N1cnJlZCBhbHJlYWR5IGluIHRoZSBjYW5jZWxsZWQgY2hhaW4uXHJcbiAgICAgKiBJdCB3aWxsIGp1c3QgdHJhY2sgd2hldGhlciByZWdpc3RlcmVkIGhhbmRsZXJzIGhhdmUgYmVlbiBleGVjdXRlZCBvciBub3QuXHJcbiAgICAgKiBUaGVyZWZvcmUsIHVuaGFuZGxlZCByZWplY3Rpb25zIHdpbGwgbmV2ZXIgYmUgc2lsZW50bHkgaGFuZGxlZCBieSBjYWxsaW5nIGNhbmNlbC5cclxuICAgICAqL1xyXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xyXG4gICAgICAgICAgICAvLyBJTlZBUklBTlQ6IHRoZSByZXN1bHQgb2YgdGhpc1tjYW5jZWxJbXBsU3ltXSBhbmQgdGhlIGJhcnJpZXIgZG8gbm90IGV2ZXIgcmVqZWN0LlxyXG4gICAgICAgICAgICAvLyBVbmZvcnR1bmF0ZWx5IG1hY09TIEhpZ2ggU2llcnJhIGRvZXMgbm90IHN1cHBvcnQgUHJvbWlzZS5hbGxTZXR0bGVkLlxyXG4gICAgICAgICAgICBQcm9taXNlLmFsbChbXHJcbiAgICAgICAgICAgICAgICB0aGlzW2NhbmNlbEltcGxTeW1dKG5ldyBDYW5jZWxFcnJvcihcIlByb21pc2UgY2FuY2VsbGVkLlwiLCB7IGNhdXNlIH0pKSxcclxuICAgICAgICAgICAgICAgIGN1cnJlbnRCYXJyaWVyKHRoaXMpXHJcbiAgICAgICAgICAgIF0pLnRoZW4oKCkgPT4gcmVzb2x2ZSgpLCAoKSA9PiByZXNvbHZlKCkpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQmluZHMgcHJvbWlzZSBjYW5jZWxsYXRpb24gdG8gdGhlIGFib3J0IGV2ZW50IG9mIHRoZSBnaXZlbiB7QGxpbmsgQWJvcnRTaWduYWx9LlxyXG4gICAgICogSWYgdGhlIHNpZ25hbCBoYXMgYWxyZWFkeSBhYm9ydGVkLCB0aGUgcHJvbWlzZSB3aWxsIGJlIGNhbmNlbGxlZCBpbW1lZGlhdGVseS5cclxuICAgICAqIFdoZW4gZWl0aGVyIGNvbmRpdGlvbiBpcyB2ZXJpZmllZCwgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZSB3aWxsIGJlIHNldFxyXG4gICAgICogdG8gdGhlIHNpZ25hbCdzIGFib3J0IHJlYXNvbiAoc2VlIHtAbGluayBBYm9ydFNpZ25hbCNyZWFzb259KS5cclxuICAgICAqXHJcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCAob3IgaWYgdGhlIHNpZ25hbCBhYm9ydHMpIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZC5cclxuICAgICAqIE9ubHkgdGhlIGZpcnN0IHNpZ25hbCB0byBhYm9ydCB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEZvciBtb3JlIGRldGFpbHMgYWJvdXQgdGhlIGNhbmNlbGxhdGlvbiBwcm9jZXNzLFxyXG4gICAgICogc2VlIHtAbGluayBjYW5jZWx9IGFuZCB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAgY29uc3RydWN0b3IuXHJcbiAgICAgKlxyXG4gICAgICogVGhpcyBtZXRob2QgZW5hYmxlcyBgYXdhaXRgaW5nIGNhbmNlbGxhYmxlIHByb21pc2VzIHdpdGhvdXQgaGF2aW5nXHJcbiAgICAgKiB0byBzdG9yZSB0aGVtIGZvciBmdXR1cmUgY2FuY2VsbGF0aW9uLCBlLmcuOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGF3YWl0IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCkuY2FuY2VsT24oc2lnbmFsKTtcclxuICAgICAqIGBgYFxyXG4gICAgICogaW5zdGVhZCBvZjpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBsZXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQgPSBsb25nUnVubmluZ09wZXJhdGlvbigpO1xyXG4gICAgICogYXdhaXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQ7XHJcbiAgICAgKiBgYGBcclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGlzIHByb21pc2UsIGZvciBtZXRob2QgY2hhaW5pbmcuXHJcbiAgICAgKi9cclxuICAgIGNhbmNlbE9uKHNpZ25hbDogQWJvcnRTaWduYWwpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIGlmIChzaWduYWwuYWJvcnRlZCkge1xyXG4gICAgICAgICAgICB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pXHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgc2lnbmFsLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKSwge2NhcHR1cmU6IHRydWV9KTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIHJldHVybiB0aGlzO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQXR0YWNoZXMgY2FsbGJhY2tzIGZvciB0aGUgcmVzb2x1dGlvbiBhbmQvb3IgcmVqZWN0aW9uIG9mIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYC5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXHJcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXHJcbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcclxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXHJcbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gb25mdWxmaWxsZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZXNvbHZlZC5cclxuICAgICAqIEBwYXJhbSBvbnJlamVjdGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVqZWN0ZWQuXHJcbiAgICAgKiBAcmV0dXJucyBBIGBDYW5jZWxsYWJsZVByb21pc2VgIGZvciB0aGUgY29tcGxldGlvbiBvZiB3aGljaGV2ZXIgY2FsbGJhY2sgaXMgZXhlY3V0ZWQuXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBob29rZWQgdXAgdG8gcHJvcGFnYXRlIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB1cCB0aGUgY2hhaW4sIGJ1dCBub3QgZG93bjpcclxuICAgICAqXHJcbiAgICAgKiAgIC0gaWYgdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgYmUgaW52b2tlZCB3aXRoIGEgYENhbmNlbEVycm9yYFxyXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSBfd2lsbCByZXNvbHZlIHJlZ3VsYXJseV8gd2l0aCBpdHMgcmVzdWx0O1xyXG4gICAgICogICAtIGNvbnZlcnNlbHksIGlmIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgX3RoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQgdG9vO19cclxuICAgICAqICAgICB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBzdGlsbCBiZSBpbnZva2VkIHdpdGggdGhlIHBhcmVudCdzIGBDYW5jZWxFcnJvcmAsXHJcbiAgICAgKiAgICAgYnV0IGl0cyByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWRcclxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhIGBDYW5jZWxFcnJvcmAgYXMgd2VsbC5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xyXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXHJcbiAgICAgKlxyXG4gICAgICogSWYgZWl0aGVyIGNhbGxiYWNrIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXHJcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICovXHJcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+IHtcclxuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS50aGVuIGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBOT1RFOiBUeXBlU2NyaXB0J3MgYnVpbHQtaW4gdHlwZSBmb3IgdGhlbiBpcyBicm9rZW4sXHJcbiAgICAgICAgLy8gYXMgaXQgYWxsb3dzIHNwZWNpZnlpbmcgYW4gYXJiaXRyYXJ5IFRSZXN1bHQxICE9IFQgZXZlbiB3aGVuIG9uZnVsZmlsbGVkIGlzIG5vdCBhIGZ1bmN0aW9uLlxyXG4gICAgICAgIC8vIFdlIGNhbm5vdCBmaXggaXQgaWYgd2Ugd2FudCB0byBDYW5jZWxsYWJsZVByb21pc2UgdG8gaW1wbGVtZW50IFByb21pc2VMaWtlPFQ+LlxyXG5cclxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25mdWxmaWxsZWQpKSB7IG9uZnVsZmlsbGVkID0gaWRlbnRpdHkgYXMgYW55OyB9XHJcbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9ucmVqZWN0ZWQpKSB7IG9ucmVqZWN0ZWQgPSB0aHJvd2VyOyB9XHJcblxyXG4gICAgICAgIGlmIChvbmZ1bGZpbGxlZCA9PT0gaWRlbnRpdHkgJiYgb25yZWplY3RlZCA9PSB0aHJvd2VyKSB7XHJcbiAgICAgICAgICAgIC8vIFNob3J0Y3V0IGZvciB0cml2aWFsIGFyZ3VtZW50cy5cclxuICAgICAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUpID0+IHJlc29sdmUodGhpcyBhcyBhbnkpKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGNvbnN0IGJhcnJpZXI6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0ge307XHJcbiAgICAgICAgdGhpc1tiYXJyaWVyU3ltXSA9IGJhcnJpZXI7XHJcblxyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgdm9pZCBzdXBlci50aGVuKFxyXG4gICAgICAgICAgICAgICAgKHZhbHVlKSA9PiB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cclxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xyXG5cclxuICAgICAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9uZnVsZmlsbGVkISh2YWx1ZSkpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9LFxyXG4gICAgICAgICAgICAgICAgKHJlYXNvbj8pID0+IHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxyXG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XHJcblxyXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25yZWplY3RlZCEocmVhc29uKSk7XHJcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgKTtcclxuICAgICAgICB9LCBhc3luYyAoY2F1c2U/KSA9PiB7XHJcbiAgICAgICAgICAgIC8vY2FuY2VsbGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBvbmNhbmNlbGxlZD8uKGNhdXNlKTtcclxuICAgICAgICAgICAgfSBmaW5hbGx5IHtcclxuICAgICAgICAgICAgICAgIGF3YWl0IHRoaXMuY2FuY2VsKGNhdXNlKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayBmb3Igb25seSB0aGUgcmVqZWN0aW9uIG9mIHRoZSBQcm9taXNlLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcclxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cclxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxyXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cclxuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIEl0IGlzIGVxdWl2YWxlbnQgdG9cclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBjYW5jZWxsYWJsZVByb21pc2UudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcclxuICAgICAqIGBgYFxyXG4gICAgICogYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXHJcbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcclxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXHJcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcclxuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxyXG4gICAgICpcclxuICAgICAqIElmIGBvbnJlamVjdGVkYCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxyXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cclxuICAgICAqL1xyXG4gICAgY2F0Y2g8VFJlc3VsdCA9IG5ldmVyPihvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gKFByb21pc2VMaWtlPFRSZXN1bHQ+IHwgVFJlc3VsdCkpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VCB8IFRSZXN1bHQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpcy50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayB0aGF0IGlzIGludm9rZWQgd2hlbiB0aGUgQ2FuY2VsbGFibGVQcm9taXNlIGlzIHNldHRsZWQgKGZ1bGZpbGxlZCBvciByZWplY3RlZCkuIFRoZVxyXG4gICAgICogcmVzb2x2ZWQgdmFsdWUgY2Fubm90IGJlIGFjY2Vzc2VkIG9yIG1vZGlmaWVkIGZyb20gdGhlIGNhbGxiYWNrLlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBzZXR0bGUgaW4gdGhlIHNhbWUgc3RhdGUgYXMgdGhlIG9yaWdpbmFsIG9uZVxyXG4gICAgICogYWZ0ZXIgdGhlIHByb3ZpZGVkIGNhbGxiYWNrIGhhcyBjb21wbGV0ZWQgZXhlY3V0aW9uLFxyXG4gICAgICogdW5sZXNzIHRoZSBjYWxsYmFjayB0aHJvd3Mgb3IgcmV0dXJucyBhIHJlamVjdGluZyBwcm9taXNlLFxyXG4gICAgICogaW4gd2hpY2ggY2FzZSB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCBhcyB3ZWxsLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcclxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cclxuICAgICAqIE9uY2UgdGhlIHBhcmVudCBwcm9taXNlIHNldHRsZXMsIHRoZSBgb25maW5hbGx5YCBjYWxsYmFjayB3aWxsIHJ1bixcclxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXHJcbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBpbXBsZW1lbnRlZCBpbiB0ZXJtcyBvZiB7QGxpbmsgdGhlbn0gYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXHJcbiAgICAgKiBJdCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxyXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXHJcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxyXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXHJcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cclxuICAgICAqXHJcbiAgICAgKiBJZiBgb25maW5hbGx5YCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxyXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cclxuICAgICAqL1xyXG4gICAgZmluYWxseShvbmZpbmFsbHk/OiAoKCkgPT4gdm9pZCkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcclxuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUuZmluYWxseSBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZmluYWxseSkpIHtcclxuICAgICAgICAgICAgcmV0dXJuIHRoaXMudGhlbihvbmZpbmFsbHksIG9uZmluYWxseSwgb25jYW5jZWxsZWQpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbihcclxuICAgICAgICAgICAgKHZhbHVlKSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB2YWx1ZSksXHJcbiAgICAgICAgICAgIChyZWFzb24/KSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB7IHRocm93IHJlYXNvbjsgfSksXHJcbiAgICAgICAgICAgIG9uY2FuY2VsbGVkLFxyXG4gICAgICAgICk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBXZSB1c2UgdGhlIGBbU3ltYm9sLnNwZWNpZXNdYCBzdGF0aWMgcHJvcGVydHksIGlmIGF2YWlsYWJsZSxcclxuICAgICAqIHRvIGRpc2FibGUgdGhlIGJ1aWx0LWluIGF1dG9tYXRpYyBzdWJjbGFzc2luZyBmZWF0dXJlcyBmcm9tIHtAbGluayBQcm9taXNlfS5cclxuICAgICAqIEl0IGlzIGNyaXRpY2FsIGZvciBwZXJmb3JtYW5jZSByZWFzb25zIHRoYXQgZXh0ZW5kZXJzIGRvIG5vdCBvdmVycmlkZSB0aGlzLlxyXG4gICAgICogT25jZSB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZ1xyXG4gICAgICogaXMgZWl0aGVyIGFjY2VwdGVkIG9yIHJldGlyZWQsIHRoaXMgaW1wbGVtZW50YXRpb24gd2lsbCBoYXZlIHRvIGJlIHJldmlzZWQgYWNjb3JkaW5nbHkuXHJcbiAgICAgKlxyXG4gICAgICogQGlnbm9yZVxyXG4gICAgICogQGludGVybmFsXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBnZXQgW3NwZWNpZXNdKCkge1xyXG4gICAgICAgIHJldHVybiBQcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xyXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUsIG9yIHJlamVjdGVkIHdoZW4gYW55IFByb21pc2UgaXMgcmVqZWN0ZWQuXHJcbiAgICAgKlxyXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgYWxsPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+W10+O1xyXG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogQXdhaXRlZDxUW1BdPjsgfT47XHJcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXHJcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxyXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcclxuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlIG9yIHJlamVjdC5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFQ+PltdPjtcclxuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFRbUF0+PjsgfT47XHJcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxyXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcclxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGxTZXR0bGVkKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVGhlIGFueSBmdW5jdGlvbiByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGlzIGZ1bGZpbGxlZCBieSB0aGUgZmlyc3QgZ2l2ZW4gcHJvbWlzZSB0byBiZSBmdWxmaWxsZWQsXHJcbiAgICAgKiBvciByZWplY3RlZCB3aXRoIGFuIEFnZ3JlZ2F0ZUVycm9yIGNvbnRhaW5pbmcgYW4gYXJyYXkgb2YgcmVqZWN0aW9uIHJlYXNvbnNcclxuICAgICAqIGlmIGFsbCBvZiB0aGUgZ2l2ZW4gcHJvbWlzZXMgYXJlIHJlamVjdGVkLlxyXG4gICAgICogSXQgcmVzb2x2ZXMgYWxsIGVsZW1lbnRzIG9mIHRoZSBwYXNzZWQgaXRlcmFibGUgdG8gcHJvbWlzZXMgYXMgaXQgcnVucyB0aGlzIGFsZ29yaXRobS5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBhbnk8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xyXG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xyXG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxyXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcclxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbnkoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XHJcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIG9yIHJlamVjdGVkIHdoZW4gYW55IG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyBhcmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQuXHJcbiAgICAgKlxyXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmFjZTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xyXG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xyXG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmFjZShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBjYW5jZWxsZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBjYW5jZWw8VCA9IG5ldmVyPihjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgY29uc3QgcCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xyXG4gICAgICAgIHAuY2FuY2VsKGNhdXNlKTtcclxuICAgICAgICByZXR1cm4gcDtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgY2FuY2Vsc1xyXG4gICAgICogYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBJZiB0aGUge0BsaW5rIEFib3J0U2lnbmFsLnRpbWVvdXR9IGZhY3RvcnkgbWV0aG9kIGlzIGF2YWlsYWJsZSxcclxuICAgICAqIGl0IGlzIHVzZWQgdG8gYmFzZSB0aGUgdGltZW91dCBvbiBfYWN0aXZlXyB0aW1lIHJhdGhlciB0aGFuIF9lbGFwc2VkXyB0aW1lLlxyXG4gICAgICogT3RoZXJ3aXNlLCBgdGltZW91dGAgZmFsbHMgYmFjayB0byB7QGxpbmsgc2V0VGltZW91dH0uXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyB0aW1lb3V0PFQgPSBuZXZlcj4obWlsbGlzZWNvbmRzOiBudW1iZXIsIGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XHJcbiAgICAgICAgaWYgKEFib3J0U2lnbmFsICYmIHR5cGVvZiBBYm9ydFNpZ25hbCA9PT0gJ2Z1bmN0aW9uJyAmJiBBYm9ydFNpZ25hbC50aW1lb3V0ICYmIHR5cGVvZiBBYm9ydFNpZ25hbC50aW1lb3V0ID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICAgICAgICAgIEFib3J0U2lnbmFsLnRpbWVvdXQobWlsbGlzZWNvbmRzKS5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpLCBtaWxsaXNlY29uZHMpO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHNsZWVwKG1pbGxpc2Vjb25kczogbnVtYmVyKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyXHJcbiAgICAgKiB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIHZhbHVlLlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHNsZWVwPFQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcclxuICAgIHN0YXRpYyBzbGVlcDxUID0gdm9pZD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlPzogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUpID0+IHtcclxuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiByZXNvbHZlKHZhbHVlISksIG1pbGxpc2Vjb25kcyk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlamVjdGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHJlYXNvbi5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlamVjdDxUID0gbmV2ZXI+KHJlYXNvbj86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKF8sIHJlamVjdCkgPT4gcmVqZWN0KHJlYXNvbikpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByZXNvbHZlKCk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICBzdGF0aWMgcmVzb2x2ZTxUID0gdm9pZD4odmFsdWU/OiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj4ge1xyXG4gICAgICAgIGlmICh2YWx1ZSBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkge1xyXG4gICAgICAgICAgICAvLyBPcHRpbWlzZSBmb3IgY2FuY2VsbGFibGUgcHJvbWlzZXMuXHJcbiAgICAgICAgICAgIHJldHVybiB2YWx1ZTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8YW55PigocmVzb2x2ZSkgPT4gcmVzb2x2ZSh2YWx1ZSkpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgYW5kIHJldHVybnMgaXQgaW4gYW4gb2JqZWN0LCBhbG9uZyB3aXRoIGl0cyByZXNvbHZlIGFuZCByZWplY3QgZnVuY3Rpb25zXHJcbiAgICAgKiBhbmQgYSBnZXR0ZXIvc2V0dGVyIGZvciB0aGUgY2FuY2VsbGF0aW9uIGNhbGxiYWNrLlxyXG4gICAgICpcclxuICAgICAqIFRoaXMgbWV0aG9kIGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyB3aXRoUmVzb2x2ZXJzPFQ+KCk6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xyXG4gICAgICAgIGxldCByZXN1bHQ6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7IG9uY2FuY2VsbGVkOiBudWxsIH0gYXMgYW55O1xyXG4gICAgICAgIHJlc3VsdC5wcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgIHJlc3VsdC5yZXNvbHZlID0gcmVzb2x2ZTtcclxuICAgICAgICAgICAgcmVzdWx0LnJlamVjdCA9IHJlamVjdDtcclxuICAgICAgICB9LCAoY2F1c2U/OiBhbnkpID0+IHsgcmVzdWx0Lm9uY2FuY2VsbGVkPy4oY2F1c2UpOyB9KTtcclxuICAgICAgICByZXR1cm4gcmVzdWx0O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgY2FuY2VsbGF0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXHJcbiAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIHJlc3VsdGluZyBmdW5jdGlvbiBkb2VzIG5vdCByZWplY3QuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYW5jZWxsZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKSB7XHJcbiAgICBsZXQgY2FuY2VsbGF0aW9uUHJvbWlzZTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0gdW5kZWZpbmVkO1xyXG5cclxuICAgIHJldHVybiAocmVhc29uOiBDYW5jZWxFcnJvcik6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9PiB7XHJcbiAgICAgICAgaWYgKCFzdGF0ZS5zZXR0bGVkKSB7XHJcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICBzdGF0ZS5yZWFzb24gPSByZWFzb247XHJcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XHJcblxyXG4gICAgICAgICAgICAvLyBBdHRhY2ggYW4gZXJyb3IgaGFuZGxlciB0aGF0IGlnbm9yZXMgdGhpcyBzcGVjaWZpYyByZWplY3Rpb24gcmVhc29uIGFuZCBub3RoaW5nIGVsc2UuXHJcbiAgICAgICAgICAgIC8vIEluIHRoZW9yeSwgYSBzYW5lIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24gYXQgdGhpcyBwb2ludFxyXG4gICAgICAgICAgICAvLyBzaG91bGQgYWx3YXlzIHJlamVjdCB3aXRoIG91ciBjYW5jZWxsYXRpb24gcmVhc29uLFxyXG4gICAgICAgICAgICAvLyBoZW5jZSB0aGUgaGFuZGxlciB3aWxsIG5ldmVyIHRocm93LlxyXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucHJvdG90eXBlLnRoZW4uY2FsbChwcm9taXNlLnByb21pc2UsIHVuZGVmaW5lZCwgKGVycikgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKGVyciAhPT0gcmVhc29uKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgdGhyb3cgZXJyO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9KTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIElmIHJlYXNvbiBpcyBub3Qgc2V0LCB0aGUgcHJvbWlzZSByZXNvbHZlZCByZWd1bGFybHksIGhlbmNlIHdlIG11c3Qgbm90IGNhbGwgb25jYW5jZWxsZWQuXHJcbiAgICAgICAgLy8gSWYgb25jYW5jZWxsZWQgaXMgdW5zZXQsIG5vIG5lZWQgdG8gZ28gYW55IGZ1cnRoZXIuXHJcbiAgICAgICAgaWYgKCFzdGF0ZS5yZWFzb24gfHwgIXByb21pc2Uub25jYW5jZWxsZWQpIHsgcmV0dXJuOyB9XHJcblxyXG4gICAgICAgIGNhbmNlbGxhdGlvblByb21pc2UgPSBuZXcgUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgcmVzb2x2ZShwcm9taXNlLm9uY2FuY2VsbGVkIShzdGF0ZS5yZWFzb24hLmNhdXNlKSk7XHJcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pLmNhdGNoKChyZWFzb24/KSA9PiB7XHJcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XHJcbiAgICAgICAgfSk7XHJcblxyXG4gICAgICAgIC8vIFVuc2V0IG9uY2FuY2VsbGVkIHRvIHByZXZlbnQgcmVwZWF0ZWQgY2FsbHMuXHJcbiAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG51bGw7XHJcblxyXG4gICAgICAgIHJldHVybiBjYW5jZWxsYXRpb25Qcm9taXNlO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVzb2x1dGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxyXG4gKi9cclxuZnVuY3Rpb24gcmVzb2x2ZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4ge1xyXG4gICAgcmV0dXJuICh2YWx1ZSkgPT4ge1xyXG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcclxuXHJcbiAgICAgICAgaWYgKHZhbHVlID09PSBwcm9taXNlLnByb21pc2UpIHtcclxuICAgICAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgVHlwZUVycm9yKFwiQSBwcm9taXNlIGNhbm5vdCBiZSByZXNvbHZlZCB3aXRoIGl0c2VsZi5cIikpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAodmFsdWUgIT0gbnVsbCAmJiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyB8fCB0eXBlb2YgdmFsdWUgPT09ICdmdW5jdGlvbicpKSB7XHJcbiAgICAgICAgICAgIGxldCB0aGVuOiBhbnk7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICB0aGVuID0gKHZhbHVlIGFzIGFueSkudGhlbjtcclxuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgICAgIHByb21pc2UucmVqZWN0KGVycik7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKHRoZW4pKSB7XHJcbiAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgIGxldCBjYW5jZWwgPSAodmFsdWUgYXMgYW55KS5jYW5jZWw7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUoY2FuY2VsKSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCBvbmNhbmNlbGxlZCA9IChjYXVzZT86IGFueSkgPT4ge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKHN0YXRlLnJlYXNvbikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gSWYgYWxyZWFkeSBjYW5jZWxsZWQsIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24uXHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSBjYW5jZWxsZXIgYWxnb3JpdGhtIGRvZXMgbm90IHJlamVjdFxyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gc28gaXQgY2FuIGJlIGRpc2NhcmRlZCBzYWZlbHkuXHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICB2b2lkIGNhbmNlbGxlckZvcih7IC4uLnByb21pc2UsIG9uY2FuY2VsbGVkIH0sIHN0YXRlKShzdGF0ZS5yZWFzb24pO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG9uY2FuY2VsbGVkO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCB7fVxyXG5cclxuICAgICAgICAgICAgICAgIGNvbnN0IG5ld1N0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcclxuICAgICAgICAgICAgICAgICAgICByb290OiBzdGF0ZS5yb290LFxyXG4gICAgICAgICAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHNldHRsZWQoKSB7IHJldHVybiB0aGlzLnJvb3Quc2V0dGxlZCB9LFxyXG4gICAgICAgICAgICAgICAgICAgIHNldCBzZXR0bGVkKHZhbHVlKSB7IHRoaXMucm9vdC5zZXR0bGVkID0gdmFsdWU7IH0sXHJcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHJlYXNvbigpIHsgcmV0dXJuIHRoaXMucm9vdC5yZWFzb24gfVxyXG4gICAgICAgICAgICAgICAgfTtcclxuXHJcbiAgICAgICAgICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIG5ld1N0YXRlKTtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseSh0aGVuLCB2YWx1ZSwgW3Jlc29sdmVyRm9yKHByb21pc2UsIG5ld1N0YXRlKSwgcmVqZWN0b3JdKTtcclxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICByZXR1cm47IC8vIElNUE9SVEFOVCFcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgcHJvbWlzZS5yZXNvbHZlKHZhbHVlKTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZWplY3Rpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cclxuICovXHJcbmZ1bmN0aW9uIHJlamVjdG9yRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yIHtcclxuICAgIHJldHVybiAocmVhc29uPykgPT4ge1xyXG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XHJcbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcclxuXHJcbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIGlmIChyZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBzdGF0ZS5yZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBPYmplY3QuaXMocmVhc29uLmNhdXNlLCBzdGF0ZS5yZWFzb24uY2F1c2UpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gU3dhbGxvdyBsYXRlIHJlamVjdGlvbnMgdGhhdCBhcmUgQ2FuY2VsRXJyb3JzIHdob3NlIGNhbmNlbGxhdGlvbiBjYXVzZSBpcyB0aGUgc2FtZSBhcyBvdXJzLlxyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSBjYXRjaCB7fVxyXG5cclxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24pKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYW5jZWxzIGFsbCB2YWx1ZXMgaW4gYW4gYXJyYXkgdGhhdCBsb29rIGxpa2UgY2FuY2VsbGFibGUgdGhlbmFibGVzLlxyXG4gKiBSZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGZ1bGZpbGxzIG9uY2UgYWxsIGNhbmNlbGxhdGlvbiBwcm9jZWR1cmVzIGZvciB0aGUgZ2l2ZW4gdmFsdWVzIGhhdmUgc2V0dGxlZC5cclxuICovXHJcbmZ1bmN0aW9uIGNhbmNlbEFsbChwYXJlbnQ6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgdmFsdWVzOiBhbnlbXSwgY2F1c2U/OiBhbnkpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIGNvbnN0IHJlc3VsdHMgPSBbXTtcclxuXHJcbiAgICBmb3IgKGNvbnN0IHZhbHVlIG9mIHZhbHVlcykge1xyXG4gICAgICAgIGxldCBjYW5jZWw6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcjtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUodmFsdWUudGhlbikpIHsgY29udGludWU7IH1cclxuICAgICAgICAgICAgY2FuY2VsID0gdmFsdWUuY2FuY2VsO1xyXG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUoY2FuY2VsKSkgeyBjb250aW51ZTsgfVxyXG4gICAgICAgIH0gY2F0Y2ggeyBjb250aW51ZTsgfVxyXG5cclxuICAgICAgICBsZXQgcmVzdWx0OiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgcmVzdWx0ID0gUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcclxuICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xyXG4gICAgICAgICAgICBjb250aW51ZTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICghcmVzdWx0KSB7IGNvbnRpbnVlOyB9XHJcbiAgICAgICAgcmVzdWx0cy5wdXNoKFxyXG4gICAgICAgICAgICAocmVzdWx0IGluc3RhbmNlb2YgUHJvbWlzZSAgPyByZXN1bHQgOiBQcm9taXNlLnJlc29sdmUocmVzdWx0KSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcclxuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICApO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBQcm9taXNlLmFsbChyZXN1bHRzKSBhcyBhbnk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGl0cyBhcmd1bWVudC5cclxuICovXHJcbmZ1bmN0aW9uIGlkZW50aXR5PFQ+KHg6IFQpOiBUIHtcclxuICAgIHJldHVybiB4O1xyXG59XHJcblxyXG4vKipcclxuICogVGhyb3dzIGl0cyBhcmd1bWVudC5cclxuICovXHJcbmZ1bmN0aW9uIHRocm93ZXIocmVhc29uPzogYW55KTogbmV2ZXIge1xyXG4gICAgdGhyb3cgcmVhc29uO1xyXG59XHJcblxyXG4vKipcclxuICogQXR0ZW1wdHMgdmFyaW91cyBzdHJhdGVnaWVzIHRvIGNvbnZlcnQgYW4gZXJyb3IgdG8gYSBzdHJpbmcuXHJcbiAqL1xyXG5mdW5jdGlvbiBlcnJvck1lc3NhZ2UoZXJyOiBhbnkpOiBzdHJpbmcge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICBpZiAoZXJyIGluc3RhbmNlb2YgRXJyb3IgfHwgdHlwZW9mIGVyciAhPT0gJ29iamVjdCcgfHwgZXJyLnRvU3RyaW5nICE9PSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nKSB7XHJcbiAgICAgICAgICAgIHJldHVybiBcIlwiICsgZXJyO1xyXG4gICAgICAgIH1cclxuICAgIH0gY2F0Y2gge31cclxuXHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiBKU09OLnN0cmluZ2lmeShlcnIpO1xyXG4gICAgfSBjYXRjaCB7fVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgcmV0dXJuIE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcuY2FsbChlcnIpO1xyXG4gICAgfSBjYXRjaCB7fVxyXG5cclxuICAgIHJldHVybiBcIjxjb3VsZCBub3QgY29udmVydCBlcnJvciB0byBzdHJpbmc+XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGJhcnJpZXIgcHJvbWlzZSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuIElmIG5lY2Vzc2FyeSwgaW5pdGlhbGlzZXMgdGhlIGJhcnJpZXIuXHJcbiAqL1xyXG5mdW5jdGlvbiBjdXJyZW50QmFycmllcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIGxldCBwd3I6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0gcHJvbWlzZVtiYXJyaWVyU3ltXSA/PyB7fTtcclxuICAgIGlmICghKCdwcm9taXNlJyBpbiBwd3IpKSB7XHJcbiAgICAgICAgT2JqZWN0LmFzc2lnbihwd3IsIHByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+KCkpO1xyXG4gICAgfVxyXG4gICAgaWYgKHByb21pc2VbYmFycmllclN5bV0gPT0gbnVsbCkge1xyXG4gICAgICAgIHB3ci5yZXNvbHZlISgpO1xyXG4gICAgICAgIHByb21pc2VbYmFycmllclN5bV0gPSBwd3I7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gcHdyLnByb21pc2UhO1xyXG59XHJcblxyXG4vLyBQb2x5ZmlsbCBQcm9taXNlLndpdGhSZXNvbHZlcnMuXHJcbmxldCBwcm9taXNlV2l0aFJlc29sdmVycyA9IFByb21pc2Uud2l0aFJlc29sdmVycztcclxuaWYgKHByb21pc2VXaXRoUmVzb2x2ZXJzICYmIHR5cGVvZiBwcm9taXNlV2l0aFJlc29sdmVycyA9PT0gJ2Z1bmN0aW9uJykge1xyXG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBwcm9taXNlV2l0aFJlc29sdmVycy5iaW5kKFByb21pc2UpO1xyXG59IGVsc2Uge1xyXG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBmdW5jdGlvbiA8VD4oKTogUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xyXG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XHJcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBQcm9taXNlPFQ+KChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xyXG4gICAgICAgIHJldHVybiB7IHByb21pc2UsIHJlc29sdmUsIHJlamVjdCB9O1xyXG4gICAgfVxyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNsaXBib2FyZCk7XHJcblxyXG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcclxuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XHJcblxyXG4vKipcclxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKlxyXG4gKiBAcGFyYW0gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxyXG4gKiBAcmV0dXJuIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dDogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRTZXRUZXh0LCB7dGV4dH0pO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSB0ZXh0IGZyb20gdGhlIENsaXBib2FyZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBUZXh0KCk6IFByb21pc2U8c3RyaW5nPiB7XHJcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIEFueSBpcyBhIGR1bW15IGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBzaW1wbGUgb3IgdW5rbm93biB0eXBlcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBBbnk8VCA9IGFueT4oc291cmNlOiBhbnkpOiBUIHtcclxuICAgIHJldHVybiBzb3VyY2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXHJcbiAqIG51bGwgc3RyaW5ncyB3aXRoIGVtcHR5IHN0cmluZ3MuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnl0ZVNsaWNlKHNvdXJjZTogYW55KTogc3RyaW5nIHtcclxuICAgIHJldHVybiAoKHNvdXJjZSA9PSBudWxsKSA/IFwiXCIgOiBzb3VyY2UpO1xyXG59XHJcblxyXG4vKipcclxuICogQXJyYXkgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XHJcbiAqIHdob3NlIGVsZW1lbnRzIGFyZSBvZiB0aGF0IHR5cGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQXJyYXk8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiBUW10ge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm4gW107XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XHJcbiAgICAgICAgICAgIHNvdXJjZVtpXSA9IGVsZW1lbnQoc291cmNlW2ldKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcclxuICogd2hvc2Uga2V5cyBhbmQgdmFsdWVzIGFyZSBvZiB0aG9zZSB0eXBlcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBNYXA8ViA9IGFueT4oa2V5OiAoc291cmNlOiBhbnkpID0+IHN0cmluZywgdmFsdWU6IChzb3VyY2U6IGFueSkgPT4gVik6IChzb3VyY2U6IGFueSkgPT4gUmVjb3JkPHN0cmluZywgVj4ge1xyXG4gICAgaWYgKHZhbHVlID09PSBBbnkpIHtcclxuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IHt9IDogc291cmNlKTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xyXG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcclxuICAgICAgICAgICAgcmV0dXJuIHt9O1xyXG4gICAgICAgIH1cclxuICAgICAgICBmb3IgKGNvbnN0IGtleSBpbiBzb3VyY2UpIHtcclxuICAgICAgICAgICAgc291cmNlW2tleV0gPSB2YWx1ZShzb3VyY2Vba2V5XSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBzb3VyY2U7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogTnVsbGFibGUgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcclxuICogYW5kIHJldHVybnMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBudWxsYWJsZSB2YWx1ZSBvZiB0aGF0IHR5cGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTnVsbGFibGU8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiAoVCB8IG51bGwpIHtcclxuICAgIGlmIChlbGVtZW50ID09PSBBbnkpIHtcclxuICAgICAgICByZXR1cm4gQW55O1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gbnVsbCA6IGVsZW1lbnQoc291cmNlKSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTdHJ1Y3QgdGFrZXMgYW4gb2JqZWN0IG1hcHBpbmcgZmllbGQgbmFtZXMgdG8gY3JlYXRpb24gZnVuY3Rpb25zXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhIHN0cnVjdC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQ6IFJlY29yZDxzdHJpbmcsIChzb3VyY2U6IGFueSkgPT4gYW55Pik6XHJcbiAgICA8VSBleHRlbmRzIFJlY29yZDxzdHJpbmcsIGFueT4gPSBhbnk+KHNvdXJjZTogYW55KSA9PiBVXHJcbntcclxuICAgIGxldCBhbGxBbnkgPSB0cnVlO1xyXG4gICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XHJcbiAgICAgICAgaWYgKGNyZWF0ZUZpZWxkW25hbWVdICE9PSBBbnkpIHtcclxuICAgICAgICAgICAgYWxsQW55ID0gZmFsc2U7XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuICAgIGlmIChhbGxBbnkpIHtcclxuICAgICAgICByZXR1cm4gQW55O1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XHJcbiAgICAgICAgICAgIGlmIChuYW1lIGluIHNvdXJjZSkge1xyXG4gICAgICAgICAgICAgICAgc291cmNlW25hbWVdID0gY3JlYXRlRmllbGRbbmFtZV0oc291cmNlW25hbWVdKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gc291cmNlO1xyXG4gICAgfTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuZXhwb3J0IGludGVyZmFjZSBTaXplIHtcclxuICAgIC8qKiBUaGUgd2lkdGggb2YgYSByZWN0YW5ndWxhciBhcmVhLiAqL1xyXG4gICAgV2lkdGg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cclxuICAgIEhlaWdodDogbnVtYmVyO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFJlY3Qge1xyXG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cclxuICAgIFg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uICovXHJcbiAgICBZOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuICovXHJcbiAgICBXaWR0aDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHJlY3RhbmdsZS4gKi9cclxuICAgIEhlaWdodDogbnVtYmVyO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFNjcmVlbiB7XHJcbiAgICAvKiogVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBzY3JlZW4uICovXHJcbiAgICBJRDogc3RyaW5nO1xyXG4gICAgLyoqIEh1bWFuLXJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIE5hbWU6IHN0cmluZztcclxuICAgIC8qKiBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLiAqL1xyXG4gICAgU2NhbGVGYWN0b3I6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgWCBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBYOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgWTogbnVtYmVyO1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBTaXplOiBTaXplO1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC4gKi9cclxuICAgIEJvdW5kczogUmVjdDtcclxuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQgKGJlZm9yZSBzY2FsaW5nKS4gKi9cclxuICAgIFBoeXNpY2FsQm91bmRzOiBSZWN0O1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLiAqL1xyXG4gICAgV29ya0FyZWE6IFJlY3Q7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIFdvcmtBcmVhIG9mIHRoZSBzY3JlZW4gKGJlZm9yZSBzY2FsaW5nKS4gKi9cclxuICAgIFBoeXNpY2FsV29ya0FyZWE6IFJlY3Q7XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXHJcbiAgICBJc1ByaW1hcnk6IGJvb2xlYW47XHJcbiAgICAvKiogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBSb3RhdGlvbjogbnVtYmVyO1xyXG59XHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuU2NyZWVucyk7XHJcblxyXG5jb25zdCBnZXRBbGwgPSAwO1xyXG5jb25zdCBnZXRQcmltYXJ5ID0gMTtcclxuY29uc3QgZ2V0Q3VycmVudCA9IDI7XHJcblxyXG4vKipcclxuICogR2V0cyBhbGwgc2NyZWVucy5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gYXJyYXkgb2YgU2NyZWVuIG9iamVjdHMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCk6IFByb21pc2U8U2NyZWVuW10+IHtcclxuICAgIHJldHVybiBjYWxsKGdldEFsbCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKTogUHJvbWlzZTxTY3JlZW4+IHtcclxuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0Q3VycmVudCk7XHJcbn1cclxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7OztBQzZCQSxJQUFNLGNBQ0Y7QUFFRyxTQUFTLE9BQU8sT0FBZSxJQUFZO0FBQzlDLE1BQUksS0FBSztBQUVULE1BQUksSUFBSSxPQUFPO0FBQ2YsU0FBTyxLQUFLO0FBRVIsVUFBTSxZQUFhLEtBQUssT0FBTyxJQUFJLEtBQU0sQ0FBQztBQUFBLEVBQzlDO0FBQ0EsU0FBTztBQUNYOzs7QUM3QkEsSUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBR3JDLElBQU0sY0FBYyxPQUFPLE9BQU87QUFBQSxFQUNyQyxNQUFNO0FBQUEsRUFDTixXQUFXO0FBQUEsRUFDWCxhQUFhO0FBQUEsRUFDYixRQUFRO0FBQUEsRUFDUixhQUFhO0FBQUEsRUFDYixRQUFRO0FBQUEsRUFDUixRQUFRO0FBQUEsRUFDUixTQUFTO0FBQUEsRUFDVCxRQUFRO0FBQUEsRUFDUixTQUFTO0FBQUEsRUFDVCxZQUFZO0FBQ2hCLENBQUM7QUFDTSxJQUFJLFdBQVcsT0FBTztBQVN0QixTQUFTLGlCQUFpQixRQUFnQixhQUFxQixJQUFJO0FBQ3RFLFNBQU8sU0FBVSxRQUFnQixPQUFZLE1BQU07QUFDL0MsV0FBTyxrQkFBa0IsUUFBUSxRQUFRLFlBQVksSUFBSTtBQUFBLEVBQzdEO0FBQ0o7QUFFQSxlQUFlLGtCQUFrQixVQUFrQixRQUFnQixZQUFvQixNQUF5QjtBQTNDaEgsTUFBQUEsS0FBQTtBQTRDSSxNQUFJLE1BQU0sSUFBSSxJQUFJLFVBQVU7QUFDNUIsTUFBSSxhQUFhLE9BQU8sVUFBVSxTQUFTLFNBQVMsQ0FBQztBQUNyRCxNQUFJLGFBQWEsT0FBTyxVQUFVLE9BQU8sU0FBUyxDQUFDO0FBQ25ELE1BQUksTUFBTTtBQUFFLFFBQUksYUFBYSxPQUFPLFFBQVEsS0FBSyxVQUFVLElBQUksQ0FBQztBQUFBLEVBQUc7QUFFbkUsTUFBSSxVQUFrQztBQUFBLElBQ2xDLENBQUMsbUJBQW1CLEdBQUc7QUFBQSxFQUMzQjtBQUNBLE1BQUksWUFBWTtBQUNaLFlBQVEscUJBQXFCLElBQUk7QUFBQSxFQUNyQztBQUVBLE1BQUksV0FBVyxNQUFNLE1BQU0sS0FBSyxFQUFFLFFBQVEsQ0FBQztBQUMzQyxNQUFJLENBQUMsU0FBUyxJQUFJO0FBQ2QsVUFBTSxJQUFJLE1BQU0sTUFBTSxTQUFTLEtBQUssQ0FBQztBQUFBLEVBQ3pDO0FBRUEsUUFBSyxNQUFBQSxNQUFBLFNBQVMsUUFBUSxJQUFJLGNBQWMsTUFBbkMsZ0JBQUFBLElBQXNDLFFBQVEsd0JBQTlDLFlBQXFFLFFBQVEsSUFBSTtBQUNsRixXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCLE9BQU87QUFDSCxXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCO0FBQ0o7OztBRnREQSxJQUFNLE9BQU8saUJBQWlCLFlBQVksT0FBTztBQUVqRCxJQUFNLGlCQUFpQjtBQU9oQixTQUFTLFFBQVEsS0FBa0M7QUFDdEQsU0FBTyxLQUFLLGdCQUFnQixFQUFDLEtBQUssSUFBSSxTQUFTLEVBQUMsQ0FBQztBQUNyRDs7O0FHdkJBO0FBQUE7QUFBQSxlQUFBQztBQUFBLEVBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBY0EsT0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBQ2xDLE9BQU8sT0FBTyxzQkFBc0I7QUFDcEMsT0FBTyxPQUFPLHVCQUF1QjtBQUlyQyxJQUFNQyxRQUFPLGlCQUFpQixZQUFZLE1BQU07QUFDaEQsSUFBTSxrQkFBa0Isb0JBQUksSUFBOEI7QUFHMUQsSUFBTSxhQUFhO0FBQ25CLElBQU0sZ0JBQWdCO0FBQ3RCLElBQU0sY0FBYztBQUNwQixJQUFNLGlCQUFpQjtBQUN2QixJQUFNLGlCQUFpQjtBQUN2QixJQUFNLGlCQUFpQjtBQTBHdkIsU0FBUyxxQkFBcUIsSUFBWSxNQUFjLFFBQXVCO0FBQzNFLE1BQUksWUFBWSxxQkFBcUIsRUFBRTtBQUN2QyxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLE1BQUksUUFBUTtBQUNSLFFBQUk7QUFDQSxnQkFBVSxRQUFRLEtBQUssTUFBTSxJQUFJLENBQUM7QUFBQSxJQUN0QyxTQUFTLEtBQVU7QUFDZixnQkFBVSxPQUFPLElBQUksVUFBVSw2QkFBNkIsSUFBSSxTQUFTLEVBQUUsT0FBTyxJQUFJLENBQUMsQ0FBQztBQUFBLElBQzVGO0FBQUEsRUFDSixPQUFPO0FBQ0gsY0FBVSxRQUFRLElBQUk7QUFBQSxFQUMxQjtBQUNKO0FBUUEsU0FBUyxvQkFBb0IsSUFBWSxTQUF1QjtBQTlKaEUsTUFBQUM7QUErSkksR0FBQUEsTUFBQSxxQkFBcUIsRUFBRSxNQUF2QixnQkFBQUEsSUFBMEIsT0FBTyxJQUFJLE9BQU8sTUFBTSxPQUFPO0FBQzdEO0FBUUEsU0FBUyxxQkFBcUIsSUFBMEM7QUFDcEUsUUFBTSxXQUFXLGdCQUFnQixJQUFJLEVBQUU7QUFDdkMsa0JBQWdCLE9BQU8sRUFBRTtBQUN6QixTQUFPO0FBQ1g7QUFPQSxTQUFTLGFBQXFCO0FBQzFCLE1BQUk7QUFDSixLQUFHO0FBQ0MsYUFBUyxPQUFPO0FBQUEsRUFDcEIsU0FBUyxnQkFBZ0IsSUFBSSxNQUFNO0FBQ25DLFNBQU87QUFDWDtBQVNBLFNBQVMsT0FBTyxNQUFjLFVBQWdGLENBQUMsR0FBaUI7QUFDNUgsUUFBTSxLQUFLLFdBQVc7QUFDdEIsU0FBTyxJQUFJLFFBQVEsQ0FBQyxTQUFTLFdBQVc7QUFDcEMsb0JBQWdCLElBQUksSUFBSSxFQUFFLFNBQVMsT0FBTyxDQUFDO0FBQzNDLElBQUFELE1BQUssTUFBTSxPQUFPLE9BQU8sRUFBRSxhQUFhLEdBQUcsR0FBRyxPQUFPLENBQUMsRUFBRSxNQUFNLENBQUMsUUFBYTtBQUN4RSxzQkFBZ0IsT0FBTyxFQUFFO0FBQ3pCLGFBQU8sR0FBRztBQUFBLElBQ2QsQ0FBQztBQUFBLEVBQ0wsQ0FBQztBQUNMO0FBUU8sU0FBUyxLQUFLLFNBQWdEO0FBQUUsU0FBTyxPQUFPLFlBQVksT0FBTztBQUFHO0FBUXBHLFNBQVMsUUFBUSxTQUFnRDtBQUFFLFNBQU8sT0FBTyxlQUFlLE9BQU87QUFBRztBQVExRyxTQUFTRSxPQUFNLFNBQWdEO0FBQUUsU0FBTyxPQUFPLGFBQWEsT0FBTztBQUFHO0FBUXRHLFNBQVMsU0FBUyxTQUFnRDtBQUFFLFNBQU8sT0FBTyxnQkFBZ0IsT0FBTztBQUFHO0FBVzVHLFNBQVMsU0FBUyxTQUE0RDtBQXRQckYsTUFBQUQ7QUFzUHVGLFVBQU9BLE1BQUEsT0FBTyxnQkFBZ0IsT0FBTyxNQUE5QixPQUFBQSxNQUFtQyxDQUFDO0FBQUc7QUFROUgsU0FBUyxTQUFTLFNBQWlEO0FBQUUsU0FBTyxPQUFPLGdCQUFnQixPQUFPO0FBQUc7OztBQzlQcEg7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDYU8sSUFBTSxpQkFBaUIsb0JBQUksSUFBd0I7QUFFbkQsSUFBTSxXQUFOLE1BQWU7QUFBQSxFQUtsQixZQUFZLFdBQW1CLFVBQStCLGNBQXNCO0FBQ2hGLFNBQUssWUFBWTtBQUNqQixTQUFLLFdBQVc7QUFDaEIsU0FBSyxlQUFlLGdCQUFnQjtBQUFBLEVBQ3hDO0FBQUEsRUFFQSxTQUFTLE1BQW9CO0FBQ3pCLFFBQUk7QUFDQSxXQUFLLFNBQVMsSUFBSTtBQUFBLElBQ3RCLFNBQVMsS0FBSztBQUNWLGNBQVEsTUFBTSxHQUFHO0FBQUEsSUFDckI7QUFFQSxRQUFJLEtBQUssaUJBQWlCLEdBQUksUUFBTztBQUNyQyxTQUFLLGdCQUFnQjtBQUNyQixXQUFPLEtBQUssaUJBQWlCO0FBQUEsRUFDakM7QUFDSjtBQUVPLFNBQVMsWUFBWSxVQUEwQjtBQUNsRCxNQUFJLFlBQVksZUFBZSxJQUFJLFNBQVMsU0FBUztBQUNyRCxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLGNBQVksVUFBVSxPQUFPLE9BQUssTUFBTSxRQUFRO0FBQ2hELE1BQUksVUFBVSxXQUFXLEdBQUc7QUFDeEIsbUJBQWUsT0FBTyxTQUFTLFNBQVM7QUFBQSxFQUM1QyxPQUFPO0FBQ0gsbUJBQWUsSUFBSSxTQUFTLFdBQVcsU0FBUztBQUFBLEVBQ3BEO0FBQ0o7OztBQ3RDTyxJQUFNLFFBQVEsT0FBTyxPQUFPO0FBQUEsRUFDbEMsU0FBUyxPQUFPLE9BQU87QUFBQSxJQUN0Qix1QkFBdUI7QUFBQSxJQUN2QixzQkFBc0I7QUFBQSxJQUN0QixvQkFBb0I7QUFBQSxJQUNwQixrQkFBa0I7QUFBQSxJQUNsQixZQUFZO0FBQUEsSUFDWixvQkFBb0I7QUFBQSxJQUNwQixvQkFBb0I7QUFBQSxJQUNwQiw0QkFBNEI7QUFBQSxJQUM1QixjQUFjO0FBQUEsSUFDZCx1QkFBdUI7QUFBQSxJQUN2QixtQkFBbUI7QUFBQSxJQUNuQixlQUFlO0FBQUEsSUFDZixlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxJQUNqQixrQkFBa0I7QUFBQSxJQUNsQixnQkFBZ0I7QUFBQSxJQUNoQixpQkFBaUI7QUFBQSxJQUNqQixpQkFBaUI7QUFBQSxJQUNqQixnQkFBZ0I7QUFBQSxJQUNoQixlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxJQUNqQixrQkFBa0I7QUFBQSxJQUNsQixZQUFZO0FBQUEsSUFDWixnQkFBZ0I7QUFBQSxJQUNoQixlQUFlO0FBQUEsSUFDZixhQUFhO0FBQUEsSUFDYixpQkFBaUI7QUFBQSxJQUNqQixvQkFBb0I7QUFBQSxJQUNwQiwwQkFBMEI7QUFBQSxJQUMxQiwyQkFBMkI7QUFBQSxJQUMzQiwwQkFBMEI7QUFBQSxJQUMxQix3QkFBd0I7QUFBQSxJQUN4QixhQUFhO0FBQUEsSUFDYixlQUFlO0FBQUEsSUFDZixnQkFBZ0I7QUFBQSxJQUNoQixZQUFZO0FBQUEsSUFDWixpQkFBaUI7QUFBQSxJQUNqQixtQkFBbUI7QUFBQSxJQUNuQixvQkFBb0I7QUFBQSxJQUNwQixxQkFBcUI7QUFBQSxJQUNyQixnQkFBZ0I7QUFBQSxJQUNoQixrQkFBa0I7QUFBQSxJQUNsQixnQkFBZ0I7QUFBQSxJQUNoQixrQkFBa0I7QUFBQSxFQUNuQixDQUFDO0FBQUEsRUFDRCxLQUFLLE9BQU8sT0FBTztBQUFBLElBQ2xCLDRCQUE0QjtBQUFBLElBQzVCLHVDQUF1QztBQUFBLElBQ3ZDLHlDQUF5QztBQUFBLElBQ3pDLDBCQUEwQjtBQUFBLElBQzFCLG9DQUFvQztBQUFBLElBQ3BDLHNDQUFzQztBQUFBLElBQ3RDLG9DQUFvQztBQUFBLElBQ3BDLDBDQUEwQztBQUFBLElBQzFDLDJCQUEyQjtBQUFBLElBQzNCLCtCQUErQjtBQUFBLElBQy9CLG9CQUFvQjtBQUFBLElBQ3BCLDRCQUE0QjtBQUFBLElBQzVCLHNCQUFzQjtBQUFBLElBQ3RCLHNCQUFzQjtBQUFBLElBQ3RCLCtCQUErQjtBQUFBLElBQy9CLDZCQUE2QjtBQUFBLElBQzdCLGdDQUFnQztBQUFBLElBQ2hDLHFCQUFxQjtBQUFBLElBQ3JCLDZCQUE2QjtBQUFBLElBQzdCLDBCQUEwQjtBQUFBLElBQzFCLHVCQUF1QjtBQUFBLElBQ3ZCLHVCQUF1QjtBQUFBLElBQ3ZCLGdCQUFnQjtBQUFBLElBQ2hCLHNCQUFzQjtBQUFBLElBQ3RCLGNBQWM7QUFBQSxJQUNkLG9CQUFvQjtBQUFBLElBQ3BCLG9CQUFvQjtBQUFBLElBQ3BCLHNCQUFzQjtBQUFBLElBQ3RCLGFBQWE7QUFBQSxJQUNiLGNBQWM7QUFBQSxJQUNkLG1CQUFtQjtBQUFBLElBQ25CLG1CQUFtQjtBQUFBLElBQ25CLHlCQUF5QjtBQUFBLElBQ3pCLGVBQWU7QUFBQSxJQUNmLGlCQUFpQjtBQUFBLElBQ2pCLHVCQUF1QjtBQUFBLElBQ3ZCLHFCQUFxQjtBQUFBLElBQ3JCLHFCQUFxQjtBQUFBLElBQ3JCLHVCQUF1QjtBQUFBLElBQ3ZCLGNBQWM7QUFBQSxJQUNkLGVBQWU7QUFBQSxJQUNmLG9CQUFvQjtBQUFBLElBQ3BCLG9CQUFvQjtBQUFBLElBQ3BCLDBCQUEwQjtBQUFBLElBQzFCLGdCQUFnQjtBQUFBLElBQ2hCLDRCQUE0QjtBQUFBLElBQzVCLDRCQUE0QjtBQUFBLElBQzVCLHlEQUF5RDtBQUFBLElBQ3pELHNDQUFzQztBQUFBLElBQ3RDLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLHFCQUFxQjtBQUFBLElBQ3JCLHNCQUFzQjtBQUFBLElBQ3RCLGdDQUFnQztBQUFBLElBQ2hDLGtDQUFrQztBQUFBLElBQ2xDLG1DQUFtQztBQUFBLElBQ25DLG9DQUFvQztBQUFBLElBQ3BDLCtCQUErQjtBQUFBLElBQy9CLDZCQUE2QjtBQUFBLElBQzdCLHVCQUF1QjtBQUFBLElBQ3ZCLGlDQUFpQztBQUFBLElBQ2pDLDhCQUE4QjtBQUFBLElBQzlCLDRCQUE0QjtBQUFBLElBQzVCLHNDQUFzQztBQUFBLElBQ3RDLDRCQUE0QjtBQUFBLElBQzVCLHNCQUFzQjtBQUFBLElBQ3RCLGtDQUFrQztBQUFBLElBQ2xDLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLHdCQUF3QjtBQUFBLElBQ3hCLG1CQUFtQjtBQUFBLElBQ25CLDBCQUEwQjtBQUFBLElBQzFCLDhCQUE4QjtBQUFBLElBQzlCLHlCQUF5QjtBQUFBLElBQ3pCLDZCQUE2QjtBQUFBLElBQzdCLGlCQUFpQjtBQUFBLElBQ2pCLGdCQUFnQjtBQUFBLElBQ2hCLHNCQUFzQjtBQUFBLElBQ3RCLGVBQWU7QUFBQSxJQUNmLHlCQUF5QjtBQUFBLElBQ3pCLHdCQUF3QjtBQUFBLElBQ3hCLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLGlCQUFpQjtBQUFBLElBQ2pCLGlCQUFpQjtBQUFBLElBQ2pCLHNCQUFzQjtBQUFBLElBQ3RCLG1DQUFtQztBQUFBLElBQ25DLHFDQUFxQztBQUFBLElBQ3JDLHVCQUF1QjtBQUFBLElBQ3ZCLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLGVBQWU7QUFBQSxJQUNmLDJCQUEyQjtBQUFBLElBQzNCLDBCQUEwQjtBQUFBLElBQzFCLDZCQUE2QjtBQUFBLElBQzdCLFlBQVk7QUFBQSxJQUNaLGdCQUFnQjtBQUFBLElBQ2hCLGtCQUFrQjtBQUFBLElBQ2xCLGdCQUFnQjtBQUFBLElBQ2hCLGtCQUFrQjtBQUFBLElBQ2xCLG1CQUFtQjtBQUFBLElBQ25CLFlBQVk7QUFBQSxJQUNaLHFCQUFxQjtBQUFBLElBQ3JCLHNCQUFzQjtBQUFBLElBQ3RCLHNCQUFzQjtBQUFBLElBQ3RCLDhCQUE4QjtBQUFBLElBQzlCLGlCQUFpQjtBQUFBLElBQ2pCLHlCQUF5QjtBQUFBLElBQ3pCLDJCQUEyQjtBQUFBLElBQzNCLCtCQUErQjtBQUFBLElBQy9CLDBCQUEwQjtBQUFBLElBQzFCLDhCQUE4QjtBQUFBLElBQzlCLGlCQUFpQjtBQUFBLElBQ2pCLHVCQUF1QjtBQUFBLElBQ3ZCLGdCQUFnQjtBQUFBLElBQ2hCLDBCQUEwQjtBQUFBLElBQzFCLHlCQUF5QjtBQUFBLElBQ3pCLHNCQUFzQjtBQUFBLElBQ3RCLGtCQUFrQjtBQUFBLElBQ2xCLG1CQUFtQjtBQUFBLElBQ25CLGtCQUFrQjtBQUFBLElBQ2xCLHVCQUF1QjtBQUFBLElBQ3ZCLG9DQUFvQztBQUFBLElBQ3BDLHNDQUFzQztBQUFBLElBQ3RDLHdCQUF3QjtBQUFBLElBQ3hCLHVCQUF1QjtBQUFBLElBQ3ZCLHlCQUF5QjtBQUFBLElBQ3pCLDRCQUE0QjtBQUFBLElBQzVCLDRCQUE0QjtBQUFBLElBQzVCLGNBQWM7QUFBQSxJQUNkLGVBQWU7QUFBQSxJQUNmLGlCQUFpQjtBQUFBLEVBQ2xCLENBQUM7QUFBQSxFQUNELE9BQU8sT0FBTyxPQUFPO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsbUJBQW1CO0FBQUEsSUFDbkIsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsSUFDakIsZUFBZTtBQUFBLElBQ2YsZ0JBQWdCO0FBQUEsSUFDaEIsbUJBQW1CO0FBQUEsRUFDcEIsQ0FBQztBQUFBLEVBQ0QsUUFBUSxPQUFPLE9BQU87QUFBQSxJQUNyQiwyQkFBMkI7QUFBQSxJQUMzQixvQkFBb0I7QUFBQSxJQUNwQixjQUFjO0FBQUEsSUFDZCxlQUFlO0FBQUEsSUFDZixlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxJQUNqQixrQkFBa0I7QUFBQSxJQUNsQixvQkFBb0I7QUFBQSxJQUNwQixhQUFhO0FBQUEsSUFDYixrQkFBa0I7QUFBQSxJQUNsQixZQUFZO0FBQUEsSUFDWixpQkFBaUI7QUFBQSxJQUNqQixnQkFBZ0I7QUFBQSxJQUNoQixnQkFBZ0I7QUFBQSxJQUNoQixlQUFlO0FBQUEsSUFDZixvQkFBb0I7QUFBQSxJQUNwQixZQUFZO0FBQUEsSUFDWixvQkFBb0I7QUFBQSxJQUNwQixrQkFBa0I7QUFBQSxJQUNsQixrQkFBa0I7QUFBQSxJQUNsQixZQUFZO0FBQUEsSUFDWixjQUFjO0FBQUEsSUFDZCxlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxFQUNsQixDQUFDO0FBQ0YsQ0FBQzs7O0FGeE5ELE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQUNsQyxPQUFPLE9BQU8scUJBQXFCO0FBRW5DLElBQU1FLFFBQU8saUJBQWlCLFlBQVksTUFBTTtBQUNoRCxJQUFNLGFBQWE7QUFZWixJQUFNLGFBQU4sTUFBaUI7QUFBQSxFQWlCcEIsWUFBWSxNQUFjLE9BQVksTUFBTTtBQUN4QyxTQUFLLE9BQU87QUFDWixTQUFLLE9BQU87QUFBQSxFQUNoQjtBQUNKO0FBRUEsU0FBUyxtQkFBbUIsT0FBWTtBQUNwQyxNQUFJLFlBQVksZUFBZSxJQUFJLE1BQU0sSUFBSTtBQUM3QyxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLE1BQUksYUFBYSxJQUFJLFdBQVcsTUFBTSxNQUFNLE1BQU0sSUFBSTtBQUN0RCxNQUFJLFlBQVksT0FBTztBQUNuQixlQUFXLFNBQVMsTUFBTTtBQUFBLEVBQzlCO0FBRUEsY0FBWSxVQUFVLE9BQU8sY0FBWSxDQUFDLFNBQVMsU0FBUyxVQUFVLENBQUM7QUFDdkUsTUFBSSxVQUFVLFdBQVcsR0FBRztBQUN4QixtQkFBZSxPQUFPLE1BQU0sSUFBSTtBQUFBLEVBQ3BDLE9BQU87QUFDSCxtQkFBZSxJQUFJLE1BQU0sTUFBTSxTQUFTO0FBQUEsRUFDNUM7QUFDSjtBQVVPLFNBQVMsV0FBVyxXQUFtQixVQUFvQixjQUFzQjtBQUNwRixNQUFJLFlBQVksZUFBZSxJQUFJLFNBQVMsS0FBSyxDQUFDO0FBQ2xELFFBQU0sZUFBZSxJQUFJLFNBQVMsV0FBVyxVQUFVLFlBQVk7QUFDbkUsWUFBVSxLQUFLLFlBQVk7QUFDM0IsaUJBQWUsSUFBSSxXQUFXLFNBQVM7QUFDdkMsU0FBTyxNQUFNLFlBQVksWUFBWTtBQUN6QztBQVNPLFNBQVMsR0FBRyxXQUFtQixVQUFnQztBQUNsRSxTQUFPLFdBQVcsV0FBVyxVQUFVLEVBQUU7QUFDN0M7QUFTTyxTQUFTLEtBQUssV0FBbUIsVUFBZ0M7QUFDcEUsU0FBTyxXQUFXLFdBQVcsVUFBVSxDQUFDO0FBQzVDO0FBT08sU0FBUyxPQUFPLFlBQXlDO0FBQzVELGFBQVcsUUFBUSxlQUFhLGVBQWUsT0FBTyxTQUFTLENBQUM7QUFDcEU7QUFLTyxTQUFTLFNBQWU7QUFDM0IsaUJBQWUsTUFBTTtBQUN6QjtBQVNPLFNBQVMsS0FBSyxNQUFjLE1BQTJCO0FBQzFELE1BQUk7QUFFSixNQUFJLE9BQU8sU0FBUyxZQUFZLFNBQVMsUUFBUSxVQUFVLFFBQVEsVUFBVSxNQUFNO0FBRS9FLFlBQVEsSUFBSSxXQUFXLEtBQUssTUFBTSxHQUFHLEtBQUssTUFBTSxDQUFDO0FBQUEsRUFDckQsT0FBTztBQUVILFlBQVEsSUFBSSxXQUFXLE1BQWdCLElBQUk7QUFBQSxFQUMvQztBQUVBLFNBQU9BLE1BQUssWUFBWSxLQUFLO0FBQ2pDOzs7QUdsSU8sU0FBUyxTQUFTLFNBQWM7QUFFbkMsVUFBUTtBQUFBLElBQ0osa0JBQWtCLFVBQVU7QUFBQSxJQUM1QjtBQUFBLElBQ0E7QUFBQSxFQUNKO0FBQ0o7QUFNTyxTQUFTLGtCQUEyQjtBQUN2QyxTQUFRLElBQUksV0FBVyxXQUFXLEVBQUcsWUFBWTtBQUNyRDtBQU1PLFNBQVMsb0JBQW9CO0FBQ2hDLE1BQUksQ0FBQyxlQUFlLENBQUMsZUFBZSxDQUFDO0FBQ2pDLFdBQU87QUFFWCxNQUFJLFNBQVM7QUFFYixRQUFNLFNBQVMsSUFBSSxZQUFZO0FBQy9CLFFBQU0sYUFBYSxJQUFJLGdCQUFnQjtBQUN2QyxTQUFPLGlCQUFpQixRQUFRLE1BQU07QUFBRSxhQUFTO0FBQUEsRUFBTyxHQUFHLEVBQUUsUUFBUSxXQUFXLE9BQU8sQ0FBQztBQUN4RixhQUFXLE1BQU07QUFDakIsU0FBTyxjQUFjLElBQUksWUFBWSxNQUFNLENBQUM7QUFFNUMsU0FBTztBQUNYO0FBS08sU0FBUyxZQUFZLE9BQTJCO0FBdER2RCxNQUFBQztBQXVESSxNQUFJLE1BQU0sa0JBQWtCLGFBQWE7QUFDckMsV0FBTyxNQUFNO0FBQUEsRUFDakIsV0FBVyxFQUFFLE1BQU0sa0JBQWtCLGdCQUFnQixNQUFNLGtCQUFrQixNQUFNO0FBQy9FLFlBQU9BLE1BQUEsTUFBTSxPQUFPLGtCQUFiLE9BQUFBLE1BQThCLFNBQVM7QUFBQSxFQUNsRCxPQUFPO0FBQ0gsV0FBTyxTQUFTO0FBQUEsRUFDcEI7QUFDSjtBQWlDQSxJQUFJLFVBQVU7QUFDZCxTQUFTLGlCQUFpQixvQkFBb0IsTUFBTTtBQUFFLFlBQVU7QUFBSyxDQUFDO0FBRS9ELFNBQVMsVUFBVSxVQUFzQjtBQUM1QyxNQUFJLFdBQVcsU0FBUyxlQUFlLFlBQVk7QUFDL0MsYUFBUztBQUFBLEVBQ2IsT0FBTztBQUNILGFBQVMsaUJBQWlCLG9CQUFvQixRQUFRO0FBQUEsRUFDMUQ7QUFDSjs7O0FDM0ZBLElBQU0saUJBQW9DO0FBQzFDLElBQU0sZUFBb0M7QUFDMUMsSUFBTSxjQUFvQztBQUMxQyxJQUFNLCtCQUFvQztBQUMxQyxJQUFNLDhCQUFvQztBQUMxQyxJQUFNLGNBQW9DO0FBQzFDLElBQU0sb0JBQW9DO0FBQzFDLElBQU0sbUJBQW9DO0FBQzFDLElBQU0sa0JBQW9DO0FBQzFDLElBQU0sZ0JBQW9DO0FBQzFDLElBQU0sZUFBb0M7QUFDMUMsSUFBTSxhQUFvQztBQUMxQyxJQUFNLGtCQUFvQztBQUMxQyxJQUFNLHFCQUFvQztBQUMxQyxJQUFNLG9CQUFvQztBQUMxQyxJQUFNLG9CQUFvQztBQUMxQyxJQUFNLGlCQUFvQztBQUMxQyxJQUFNLGlCQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0scUJBQW9DO0FBQzFDLElBQU0seUJBQW9DO0FBQzFDLElBQU0sZUFBb0M7QUFDMUMsSUFBTSxrQkFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxvQkFBb0M7QUFDMUMsSUFBTSx1QkFBb0M7QUFDMUMsSUFBTSw0QkFBb0M7QUFDMUMsSUFBTSxxQkFBb0M7QUFDMUMsSUFBTSxtQ0FBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSw0QkFBb0M7QUFDMUMsSUFBTSxxQkFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxpQkFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxhQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0seUJBQW9DO0FBQzFDLElBQU0sdUJBQW9DO0FBQzFDLElBQU0scUJBQW9DO0FBQzFDLElBQU0sbUJBQW9DO0FBQzFDLElBQU0sbUJBQW9DO0FBQzFDLElBQU0sY0FBb0M7QUFDMUMsSUFBTSxhQUFvQztBQUMxQyxJQUFNLGVBQW9DO0FBQzFDLElBQU0sZ0JBQW9DO0FBQzFDLElBQU0sa0JBQW9DO0FBdUIxQyxJQUFNLFlBQVksT0FBTyxRQUFRO0FBSXBCO0FBRmIsSUFBTSxVQUFOLE1BQU0sUUFBTztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBVVQsWUFBWSxPQUFlLElBQUk7QUFDM0IsU0FBSyxTQUFTLElBQUksaUJBQWlCLFlBQVksUUFBUSxJQUFJO0FBRzNELGVBQVcsVUFBVSxPQUFPLG9CQUFvQixRQUFPLFNBQVMsR0FBRztBQUMvRCxVQUNJLFdBQVcsaUJBQ1IsT0FBUSxLQUFhLE1BQU0sTUFBTSxZQUN0QztBQUNFLFFBQUMsS0FBYSxNQUFNLElBQUssS0FBYSxNQUFNLEVBQUUsS0FBSyxJQUFJO0FBQUEsTUFDM0Q7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBUUEsSUFBSSxNQUFzQjtBQUN0QixXQUFPLElBQUksUUFBTyxJQUFJO0FBQUEsRUFDMUI7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxXQUE4QjtBQUMxQixXQUFPLEtBQUssU0FBUyxFQUFFLGNBQWM7QUFBQSxFQUN6QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsU0FBd0I7QUFDcEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxZQUFZO0FBQUEsRUFDdkM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFFBQXVCO0FBQ25CLFdBQU8sS0FBSyxTQUFTLEVBQUUsV0FBVztBQUFBLEVBQ3RDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSx5QkFBd0M7QUFDcEMsV0FBTyxLQUFLLFNBQVMsRUFBRSw0QkFBNEI7QUFBQSxFQUN2RDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0Esd0JBQXVDO0FBQ25DLFdBQU8sS0FBSyxTQUFTLEVBQUUsMkJBQTJCO0FBQUEsRUFDdEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFFBQXVCO0FBQ25CLFdBQU8sS0FBSyxTQUFTLEVBQUUsV0FBVztBQUFBLEVBQ3RDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxjQUE2QjtBQUN6QixXQUFPLEtBQUssU0FBUyxFQUFFLGlCQUFpQjtBQUFBLEVBQzVDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxhQUE0QjtBQUN4QixXQUFPLEtBQUssU0FBUyxFQUFFLGdCQUFnQjtBQUFBLEVBQzNDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsWUFBNkI7QUFDekIsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlO0FBQUEsRUFDMUM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxVQUEyQjtBQUN2QixXQUFPLEtBQUssU0FBUyxFQUFFLGFBQWE7QUFBQSxFQUN4QztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFNBQTBCO0FBQ3RCLFdBQU8sS0FBSyxTQUFTLEVBQUUsWUFBWTtBQUFBLEVBQ3ZDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxPQUFzQjtBQUNsQixXQUFPLEtBQUssU0FBUyxFQUFFLFVBQVU7QUFBQSxFQUNyQztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFlBQThCO0FBQzFCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZUFBZTtBQUFBLEVBQzFDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsZUFBaUM7QUFDN0IsV0FBTyxLQUFLLFNBQVMsRUFBRSxrQkFBa0I7QUFBQSxFQUM3QztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGNBQWdDO0FBQzVCLFdBQU8sS0FBSyxTQUFTLEVBQUUsaUJBQWlCO0FBQUEsRUFDNUM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxjQUFnQztBQUM1QixXQUFPLEtBQUssU0FBUyxFQUFFLGlCQUFpQjtBQUFBLEVBQzVDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxXQUEwQjtBQUN0QixXQUFPLEtBQUssU0FBUyxFQUFFLGNBQWM7QUFBQSxFQUN6QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsV0FBMEI7QUFDdEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxjQUFjO0FBQUEsRUFDekM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxPQUF3QjtBQUNwQixXQUFPLEtBQUssU0FBUyxFQUFFLFVBQVU7QUFBQSxFQUNyQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsZUFBOEI7QUFDMUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxrQkFBa0I7QUFBQSxFQUM3QztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLG1CQUFzQztBQUNsQyxXQUFPLEtBQUssU0FBUyxFQUFFLHNCQUFzQjtBQUFBLEVBQ2pEO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxTQUF3QjtBQUNwQixXQUFPLEtBQUssU0FBUyxFQUFFLFlBQVk7QUFBQSxFQUN2QztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFlBQThCO0FBQzFCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZUFBZTtBQUFBLEVBQzFDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxVQUF5QjtBQUNyQixXQUFPLEtBQUssU0FBUyxFQUFFLGFBQWE7QUFBQSxFQUN4QztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBUUEsWUFBWSxHQUFXLEdBQTBCO0FBQzdDLFdBQU8sS0FBSyxTQUFTLEVBQUUsbUJBQW1CLEVBQUUsR0FBRyxFQUFFLENBQUM7QUFBQSxFQUN0RDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGVBQWUsYUFBcUM7QUFDaEQsV0FBTyxLQUFLLFNBQVMsRUFBRSxzQkFBc0IsRUFBRSxZQUFZLENBQUM7QUFBQSxFQUNoRTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVVBLG9CQUFvQixHQUFXLEdBQVcsR0FBVyxHQUEwQjtBQUMzRSxXQUFPLEtBQUssU0FBUyxFQUFFLDJCQUEyQixFQUFFLEdBQUcsR0FBRyxHQUFHLEVBQUUsQ0FBQztBQUFBLEVBQ3BFO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsYUFBYSxXQUFtQztBQUM1QyxXQUFPLEtBQUssU0FBUyxFQUFFLG9CQUFvQixFQUFFLFVBQVUsQ0FBQztBQUFBLEVBQzVEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsMkJBQTJCLFNBQWlDO0FBQ3hELFdBQU8sS0FBSyxTQUFTLEVBQUUsa0NBQWtDLEVBQUUsUUFBUSxDQUFDO0FBQUEsRUFDeEU7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLFdBQVcsT0FBZSxRQUErQjtBQUNyRCxXQUFPLEtBQUssU0FBUyxFQUFFLGtCQUFrQixFQUFFLE9BQU8sT0FBTyxDQUFDO0FBQUEsRUFDOUQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLFdBQVcsT0FBZSxRQUErQjtBQUNyRCxXQUFPLEtBQUssU0FBUyxFQUFFLGtCQUFrQixFQUFFLE9BQU8sT0FBTyxDQUFDO0FBQUEsRUFDOUQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLG9CQUFvQixHQUFXLEdBQTBCO0FBQ3JELFdBQU8sS0FBSyxTQUFTLEVBQUUsMkJBQTJCLEVBQUUsR0FBRyxFQUFFLENBQUM7QUFBQSxFQUM5RDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGFBQWFDLFlBQW1DO0FBQzVDLFdBQU8sS0FBSyxTQUFTLEVBQUUsb0JBQW9CLEVBQUUsV0FBQUEsV0FBVSxDQUFDO0FBQUEsRUFDNUQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLFFBQVEsT0FBZSxRQUErQjtBQUNsRCxXQUFPLEtBQUssU0FBUyxFQUFFLGVBQWUsRUFBRSxPQUFPLE9BQU8sQ0FBQztBQUFBLEVBQzNEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsU0FBUyxPQUE4QjtBQUNuQyxXQUFPLEtBQUssU0FBUyxFQUFFLGdCQUFnQixFQUFFLE1BQU0sQ0FBQztBQUFBLEVBQ3BEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsUUFBUSxNQUE2QjtBQUNqQyxXQUFPLEtBQUssU0FBUyxFQUFFLGVBQWUsRUFBRSxLQUFLLENBQUM7QUFBQSxFQUNsRDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxPQUFzQjtBQUNsQixXQUFPLEtBQUssU0FBUyxFQUFFLFVBQVU7QUFBQSxFQUNyQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsbUJBQWtDO0FBQzlCLFdBQU8sS0FBSyxTQUFTLEVBQUUsc0JBQXNCO0FBQUEsRUFDakQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGlCQUFnQztBQUM1QixXQUFPLEtBQUssU0FBUyxFQUFFLG9CQUFvQjtBQUFBLEVBQy9DO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxlQUE4QjtBQUMxQixXQUFPLEtBQUssU0FBUyxFQUFFLGtCQUFrQjtBQUFBLEVBQzdDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxhQUE0QjtBQUN4QixXQUFPLEtBQUssU0FBUyxFQUFFLGdCQUFnQjtBQUFBLEVBQzNDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxhQUE0QjtBQUN4QixXQUFPLEtBQUssU0FBUyxFQUFFLGdCQUFnQjtBQUFBLEVBQzNDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsUUFBeUI7QUFDckIsV0FBTyxLQUFLLFNBQVMsRUFBRSxXQUFXO0FBQUEsRUFDdEM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLE9BQXNCO0FBQ2xCLFdBQU8sS0FBSyxTQUFTLEVBQUUsVUFBVTtBQUFBLEVBQ3JDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxTQUF3QjtBQUNwQixXQUFPLEtBQUssU0FBUyxFQUFFLFlBQVk7QUFBQSxFQUN2QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsVUFBeUI7QUFDckIsV0FBTyxLQUFLLFNBQVMsRUFBRSxhQUFhO0FBQUEsRUFDeEM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFlBQTJCO0FBQ3ZCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZUFBZTtBQUFBLEVBQzFDO0FBQ0o7QUEzYUEsSUFBTSxTQUFOO0FBZ2JBLElBQU0sYUFBYSxJQUFJLE9BQU8sRUFBRTtBQUVoQyxJQUFPLGlCQUFROzs7QVRqZmYsU0FBUyxVQUFVLFdBQW1CLE9BQVksTUFBWTtBQUMxRCxPQUFLLEtBQUssV0FBVyxJQUFJO0FBQzdCO0FBUUEsU0FBUyxpQkFBaUIsWUFBb0IsWUFBb0I7QUFDOUQsUUFBTSxlQUFlLGVBQU8sSUFBSSxVQUFVO0FBQzFDLFFBQU0sU0FBVSxhQUFxQixVQUFVO0FBRS9DLE1BQUksT0FBTyxXQUFXLFlBQVk7QUFDOUIsWUFBUSxNQUFNLGtCQUFrQixtQkFBVSxjQUFhO0FBQ3ZEO0FBQUEsRUFDSjtBQUVBLE1BQUk7QUFDQSxXQUFPLEtBQUssWUFBWTtBQUFBLEVBQzVCLFNBQVMsR0FBRztBQUNSLFlBQVEsTUFBTSxnQ0FBZ0MsbUJBQVUsUUFBTyxDQUFDO0FBQUEsRUFDcEU7QUFDSjtBQUtBLFNBQVMsZUFBZSxJQUFpQjtBQUNyQyxRQUFNLFVBQVUsR0FBRztBQUVuQixXQUFTLFVBQVUsU0FBUyxPQUFPO0FBQy9CLFFBQUksV0FBVztBQUNYO0FBRUosVUFBTSxZQUFZLFFBQVEsYUFBYSxXQUFXLEtBQUssUUFBUSxhQUFhLGdCQUFnQjtBQUM1RixVQUFNLGVBQWUsUUFBUSxhQUFhLG1CQUFtQixLQUFLLFFBQVEsYUFBYSx3QkFBd0IsS0FBSztBQUNwSCxVQUFNLGVBQWUsUUFBUSxhQUFhLFlBQVksS0FBSyxRQUFRLGFBQWEsaUJBQWlCO0FBQ2pHLFVBQU0sTUFBTSxRQUFRLGFBQWEsYUFBYSxLQUFLLFFBQVEsYUFBYSxrQkFBa0I7QUFFMUYsUUFBSSxjQUFjO0FBQ2QsZ0JBQVUsU0FBUztBQUN2QixRQUFJLGlCQUFpQjtBQUNqQix1QkFBaUIsY0FBYyxZQUFZO0FBQy9DLFFBQUksUUFBUTtBQUNSLFdBQUssUUFBUSxHQUFHO0FBQUEsRUFDeEI7QUFFQSxRQUFNLFVBQVUsUUFBUSxhQUFhLGFBQWEsS0FBSyxRQUFRLGFBQWEsa0JBQWtCO0FBRTlGLE1BQUksU0FBUztBQUNULGFBQVM7QUFBQSxNQUNMLE9BQU87QUFBQSxNQUNQLFNBQVM7QUFBQSxNQUNULFVBQVU7QUFBQSxNQUNWLFNBQVM7QUFBQSxRQUNMLEVBQUUsT0FBTyxNQUFNO0FBQUEsUUFDZixFQUFFLE9BQU8sTUFBTSxXQUFXLEtBQUs7QUFBQSxNQUNuQztBQUFBLElBQ0osQ0FBQyxFQUFFLEtBQUssU0FBUztBQUFBLEVBQ3JCLE9BQU87QUFDSCxjQUFVO0FBQUEsRUFDZDtBQUNKO0FBR0EsSUFBTSxnQkFBZ0IsT0FBTyxZQUFZO0FBQ3pDLElBQU0sZ0JBQWdCLE9BQU8sWUFBWTtBQUN6QyxJQUFNLGtCQUFrQixPQUFPLGNBQWM7QUFReEM7QUFGTCxJQUFNLDBCQUFOLE1BQThCO0FBQUEsRUFJMUIsY0FBYztBQUNWLFNBQUssYUFBYSxJQUFJLElBQUksZ0JBQWdCO0FBQUEsRUFDOUM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBU0EsSUFBSSxTQUFrQixVQUE2QztBQUMvRCxXQUFPLEVBQUUsUUFBUSxLQUFLLGFBQWEsRUFBRSxPQUFPO0FBQUEsRUFDaEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFFBQWM7QUFDVixTQUFLLGFBQWEsRUFBRSxNQUFNO0FBQzFCLFNBQUssYUFBYSxJQUFJLElBQUksZ0JBQWdCO0FBQUEsRUFDOUM7QUFDSjtBQVNLLGVBRUE7QUFKTCxJQUFNLGtCQUFOLE1BQXNCO0FBQUEsRUFNbEIsY0FBYztBQUNWLFNBQUssYUFBYSxJQUFJLG9CQUFJLFFBQVE7QUFDbEMsU0FBSyxlQUFlLElBQUk7QUFBQSxFQUM1QjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBUUEsSUFBSSxTQUFrQixVQUE2QztBQUMvRCxRQUFJLENBQUMsS0FBSyxhQUFhLEVBQUUsSUFBSSxPQUFPLEdBQUc7QUFBRSxXQUFLLGVBQWU7QUFBQSxJQUFLO0FBQ2xFLFNBQUssYUFBYSxFQUFFLElBQUksU0FBUyxRQUFRO0FBQ3pDLFdBQU8sQ0FBQztBQUFBLEVBQ1o7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFFBQWM7QUFDVixRQUFJLEtBQUssZUFBZSxLQUFLO0FBQ3pCO0FBRUosZUFBVyxXQUFXLFNBQVMsS0FBSyxpQkFBaUIsR0FBRyxHQUFHO0FBQ3ZELFVBQUksS0FBSyxlQUFlLEtBQUs7QUFDekI7QUFFSixZQUFNLFdBQVcsS0FBSyxhQUFhLEVBQUUsSUFBSSxPQUFPO0FBQ2hELFVBQUksWUFBWSxNQUFNO0FBQUUsYUFBSyxlQUFlO0FBQUEsTUFBSztBQUVqRCxpQkFBVyxXQUFXLFlBQVksQ0FBQztBQUMvQixnQkFBUSxvQkFBb0IsU0FBUyxjQUFjO0FBQUEsSUFDM0Q7QUFFQSxTQUFLLGFBQWEsSUFBSSxvQkFBSSxRQUFRO0FBQ2xDLFNBQUssZUFBZSxJQUFJO0FBQUEsRUFDNUI7QUFDSjtBQUVBLElBQU0sa0JBQWtCLGtCQUFrQixJQUFJLElBQUksd0JBQXdCLElBQUksSUFBSSxnQkFBZ0I7QUFLbEcsU0FBUyxnQkFBZ0IsU0FBd0I7QUFDN0MsUUFBTSxnQkFBZ0I7QUFDdEIsUUFBTSxjQUFlLFFBQVEsYUFBYSxhQUFhLEtBQUssUUFBUSxhQUFhLGtCQUFrQixLQUFLO0FBQ3hHLFFBQU0sV0FBcUIsQ0FBQztBQUU1QixNQUFJO0FBQ0osVUFBUSxRQUFRLGNBQWMsS0FBSyxXQUFXLE9BQU87QUFDakQsYUFBUyxLQUFLLE1BQU0sQ0FBQyxDQUFDO0FBRTFCLFFBQU0sVUFBVSxnQkFBZ0IsSUFBSSxTQUFTLFFBQVE7QUFDckQsYUFBVyxXQUFXO0FBQ2xCLFlBQVEsaUJBQWlCLFNBQVMsZ0JBQWdCLE9BQU87QUFDakU7QUFLTyxTQUFTLFNBQWU7QUFDM0IsWUFBVSxNQUFNO0FBQ3BCO0FBS08sU0FBUyxTQUFlO0FBQzNCLGtCQUFnQixNQUFNO0FBQ3RCLFdBQVMsS0FBSyxpQkFBaUIsbUdBQW1HLEVBQUUsUUFBUSxlQUFlO0FBQy9KOzs7QVVoTUEsT0FBTyxRQUFRO0FBQ2YsT0FBVTtBQUVWLElBQUksTUFBTztBQUNQLFdBQVMsc0JBQXNCO0FBQ25DOzs7QUNyQkE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQVlBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksTUFBTTtBQUVoRCxJQUFNLG1CQUFtQjtBQUN6QixJQUFNLG9CQUFvQjtBQUUxQixJQUFNLFVBQVcsV0FBWTtBQWpCN0IsTUFBQUMsS0FBQTtBQWtCSSxNQUFJO0FBQ0EsU0FBSyxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsWUFBdkIsbUJBQWdDLGFBQWE7QUFDOUMsYUFBUSxPQUFlLE9BQU8sUUFBUSxZQUFZLEtBQU0sT0FBZSxPQUFPLE9BQU87QUFBQSxJQUN6RixZQUFZLHdCQUFlLFdBQWYsbUJBQXVCLG9CQUF2QixtQkFBeUMsZ0JBQXpDLG1CQUFzRCxhQUFhO0FBQzNFLGFBQVEsT0FBZSxPQUFPLGdCQUFnQixVQUFVLEVBQUUsWUFBWSxLQUFNLE9BQWUsT0FBTyxnQkFBZ0IsVUFBVSxDQUFDO0FBQUEsSUFDakk7QUFBQSxFQUNKLFNBQVEsR0FBRztBQUFBLEVBQUM7QUFFWixVQUFRO0FBQUEsSUFBSztBQUFBLElBQ1Q7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLEVBQXdEO0FBQzVELFNBQU87QUFDWCxFQUFHO0FBRUksU0FBUyxPQUFPLEtBQWdCO0FBQ25DLHFDQUFVO0FBQ2Q7QUFPTyxTQUFTLGFBQStCO0FBQzNDLFNBQU9ELE1BQUssZ0JBQWdCO0FBQ2hDO0FBT0EsZUFBc0IsZUFBNkM7QUFDL0QsTUFBSSxXQUFXLE1BQU0sTUFBTSxxQkFBcUI7QUFDaEQsTUFBSSxTQUFTLElBQUk7QUFDYixXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCLE9BQU87QUFDSCxVQUFNLElBQUksTUFBTSxtQ0FBbUMsU0FBUyxVQUFVO0FBQUEsRUFDMUU7QUFDSjtBQStCTyxTQUFTLGNBQXdDO0FBQ3BELFNBQU9BLE1BQUssaUJBQWlCO0FBQ2pDO0FBT08sU0FBUyxZQUFxQjtBQUNqQyxTQUFPLE9BQU8sT0FBTyxZQUFZLE9BQU87QUFDNUM7QUFPTyxTQUFTLFVBQW1CO0FBQy9CLFNBQU8sT0FBTyxPQUFPLFlBQVksT0FBTztBQUM1QztBQU9PLFNBQVMsUUFBaUI7QUFDN0IsU0FBTyxPQUFPLE9BQU8sWUFBWSxPQUFPO0FBQzVDO0FBT08sU0FBUyxVQUFtQjtBQUMvQixTQUFPLE9BQU8sT0FBTyxZQUFZLFNBQVM7QUFDOUM7QUFPTyxTQUFTLFFBQWlCO0FBQzdCLFNBQU8sT0FBTyxPQUFPLFlBQVksU0FBUztBQUM5QztBQU9PLFNBQVMsVUFBbUI7QUFDL0IsU0FBTyxPQUFPLE9BQU8sWUFBWSxTQUFTO0FBQzlDO0FBT08sU0FBUyxVQUFtQjtBQUMvQixTQUFPLFFBQVEsT0FBTyxPQUFPLFlBQVksS0FBSztBQUNsRDs7O0FDM0lBLE9BQU8saUJBQWlCLGVBQWUsa0JBQWtCO0FBRXpELElBQU1FLFFBQU8saUJBQWlCLFlBQVksV0FBVztBQUVyRCxJQUFNLGtCQUFrQjtBQUV4QixTQUFTLGdCQUFnQixJQUFZLEdBQVcsR0FBVyxNQUFpQjtBQUN4RSxPQUFLQSxNQUFLLGlCQUFpQixFQUFDLElBQUksR0FBRyxHQUFHLEtBQUksQ0FBQztBQUMvQztBQUVBLFNBQVMsbUJBQW1CLE9BQW1CO0FBQzNDLFFBQU0sU0FBUyxZQUFZLEtBQUs7QUFHaEMsUUFBTSxvQkFBb0IsT0FBTyxpQkFBaUIsTUFBTSxFQUFFLGlCQUFpQixzQkFBc0IsRUFBRSxLQUFLO0FBRXhHLE1BQUksbUJBQW1CO0FBQ25CLFVBQU0sZUFBZTtBQUNyQixVQUFNLE9BQU8sT0FBTyxpQkFBaUIsTUFBTSxFQUFFLGlCQUFpQiwyQkFBMkI7QUFDekYsb0JBQWdCLG1CQUFtQixNQUFNLFNBQVMsTUFBTSxTQUFTLElBQUk7QUFBQSxFQUN6RSxPQUFPO0FBQ0gsOEJBQTBCLE9BQU8sTUFBTTtBQUFBLEVBQzNDO0FBQ0o7QUFVQSxTQUFTLDBCQUEwQixPQUFtQixRQUFxQjtBQUV2RSxNQUFJLFFBQVEsR0FBRztBQUNYO0FBQUEsRUFDSjtBQUdBLFVBQVEsT0FBTyxpQkFBaUIsTUFBTSxFQUFFLGlCQUFpQix1QkFBdUIsRUFBRSxLQUFLLEdBQUc7QUFBQSxJQUN0RixLQUFLO0FBQ0Q7QUFBQSxJQUNKLEtBQUs7QUFDRCxZQUFNLGVBQWU7QUFDckI7QUFBQSxFQUNSO0FBR0EsTUFBSSxPQUFPLG1CQUFtQjtBQUMxQjtBQUFBLEVBQ0o7QUFHQSxRQUFNLFlBQVksT0FBTyxhQUFhO0FBQ3RDLFFBQU0sZUFBZSxhQUFhLFVBQVUsU0FBUyxFQUFFLFNBQVM7QUFDaEUsTUFBSSxjQUFjO0FBQ2QsYUFBUyxJQUFJLEdBQUcsSUFBSSxVQUFVLFlBQVksS0FBSztBQUMzQyxZQUFNLFFBQVEsVUFBVSxXQUFXLENBQUM7QUFDcEMsWUFBTSxRQUFRLE1BQU0sZUFBZTtBQUNuQyxlQUFTLElBQUksR0FBRyxJQUFJLE1BQU0sUUFBUSxLQUFLO0FBQ25DLGNBQU0sT0FBTyxNQUFNLENBQUM7QUFDcEIsWUFBSSxTQUFTLGlCQUFpQixLQUFLLE1BQU0sS0FBSyxHQUFHLE1BQU0sUUFBUTtBQUMzRDtBQUFBLFFBQ0o7QUFBQSxNQUNKO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFHQSxNQUFJLGtCQUFrQixvQkFBb0Isa0JBQWtCLHFCQUFxQjtBQUM3RSxRQUFJLGdCQUFpQixDQUFDLE9BQU8sWUFBWSxDQUFDLE9BQU8sVUFBVztBQUN4RDtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBR0EsUUFBTSxlQUFlO0FBQ3pCOzs7QUM3RkE7QUFBQTtBQUFBO0FBQUE7QUFnQk8sU0FBUyxRQUFRLEtBQWtCO0FBQ3RDLE1BQUk7QUFDQSxXQUFPLE9BQU8sT0FBTyxNQUFNLEdBQUc7QUFBQSxFQUNsQyxTQUFTLEdBQUc7QUFDUixVQUFNLElBQUksTUFBTSw4QkFBOEIsTUFBTSxRQUFRLEdBQUcsRUFBRSxPQUFPLEVBQUUsQ0FBQztBQUFBLEVBQy9FO0FBQ0o7OztBQ1BBLElBQUksVUFBVTtBQUNkLElBQUksV0FBVztBQUVmLElBQUksWUFBWTtBQUNoQixJQUFJLFlBQVk7QUFDaEIsSUFBSSxXQUFXO0FBQ2YsSUFBSSxhQUFxQjtBQUN6QixJQUFJLGdCQUFnQjtBQUVwQixJQUFJLFVBQVU7QUFDZCxJQUFNLGlCQUFpQixnQkFBZ0I7QUFFdkMsT0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBQ2xDLE9BQU8sT0FBTyxlQUFlLENBQUMsVUFBeUI7QUFDbkQsY0FBWTtBQUNaLE1BQUksQ0FBQyxXQUFXO0FBRVosZ0JBQVksV0FBVztBQUN2QixjQUFVO0FBQUEsRUFDZDtBQUNKO0FBRUEsT0FBTyxpQkFBaUIsYUFBYSxRQUFRLEVBQUUsU0FBUyxLQUFLLENBQUM7QUFDOUQsT0FBTyxpQkFBaUIsYUFBYSxRQUFRLEVBQUUsU0FBUyxLQUFLLENBQUM7QUFDOUQsT0FBTyxpQkFBaUIsV0FBVyxRQUFRLEVBQUUsU0FBUyxLQUFLLENBQUM7QUFDNUQsV0FBVyxNQUFNLENBQUMsU0FBUyxlQUFlLFVBQVUsR0FBRztBQUNuRCxTQUFPLGlCQUFpQixJQUFJLGVBQWUsRUFBRSxTQUFTLEtBQUssQ0FBQztBQUNoRTtBQUVBLFNBQVMsY0FBYyxPQUFjO0FBRWpDLE1BQUksWUFBWSxVQUFVO0FBQ3RCLFVBQU0seUJBQXlCO0FBQy9CLFVBQU0sZ0JBQWdCO0FBQ3RCLFVBQU0sZUFBZTtBQUFBLEVBQ3pCO0FBQ0o7QUFHQSxJQUFNLFlBQVk7QUFDbEIsSUFBTSxVQUFZO0FBQ2xCLElBQU0sWUFBWTtBQUVsQixTQUFTLE9BQU8sT0FBbUI7QUFJL0IsTUFBSSxXQUFtQixlQUFlLE1BQU07QUFDNUMsVUFBUSxNQUFNLE1BQU07QUFBQSxJQUNoQixLQUFLO0FBQ0Qsa0JBQVk7QUFDWixVQUFJLENBQUMsZ0JBQWdCO0FBQUUsdUJBQWUsVUFBVyxLQUFLLE1BQU07QUFBQSxNQUFTO0FBQ3JFO0FBQUEsSUFDSixLQUFLO0FBQ0Qsa0JBQVk7QUFDWixVQUFJLENBQUMsZ0JBQWdCO0FBQUUsdUJBQWUsVUFBVSxFQUFFLEtBQUssTUFBTTtBQUFBLE1BQVM7QUFDdEU7QUFBQSxJQUNKO0FBQ0ksa0JBQVk7QUFDWixVQUFJLENBQUMsZ0JBQWdCO0FBQUUsdUJBQWU7QUFBQSxNQUFTO0FBQy9DO0FBQUEsRUFDUjtBQUVBLE1BQUksV0FBVyxVQUFVLENBQUM7QUFDMUIsTUFBSSxVQUFVLGVBQWUsQ0FBQztBQUU5QixZQUFVO0FBR1YsTUFBSSxjQUFjLGFBQWEsRUFBRSxVQUFVLE1BQU0sU0FBUztBQUN0RCxnQkFBYSxLQUFLLE1BQU07QUFDeEIsZUFBWSxLQUFLLE1BQU07QUFBQSxFQUMzQjtBQUlBLE1BQ0ksY0FBYyxhQUNYLFlBRUMsYUFFSSxjQUFjLGFBQ1gsTUFBTSxXQUFXLElBRzlCO0FBQ0UsVUFBTSx5QkFBeUI7QUFDL0IsVUFBTSxnQkFBZ0I7QUFDdEIsVUFBTSxlQUFlO0FBQUEsRUFDekI7QUFHQSxNQUFJLFdBQVcsR0FBRztBQUFFLGNBQVUsS0FBSztBQUFBLEVBQUc7QUFFdEMsTUFBSSxVQUFVLEdBQUc7QUFBRSxnQkFBWSxLQUFLO0FBQUEsRUFBRztBQUd2QyxNQUFJLGNBQWMsV0FBVztBQUFFLGdCQUFZLEtBQUs7QUFBQSxFQUFHO0FBQUM7QUFDeEQ7QUFFQSxTQUFTLFlBQVksT0FBeUI7QUFFMUMsWUFBVTtBQUNWLGNBQVk7QUFHWixNQUFJLENBQUMsVUFBVSxHQUFHO0FBQ2QsUUFBSSxNQUFNLFNBQVMsZUFBZSxNQUFNLFdBQVcsS0FBSyxNQUFNLFdBQVcsR0FBRztBQUN4RTtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBRUEsTUFBSSxZQUFZO0FBRVosZ0JBQVk7QUFFWjtBQUFBLEVBQ0o7QUFHQSxRQUFNLFNBQVMsWUFBWSxLQUFLO0FBSWhDLFFBQU0sUUFBUSxPQUFPLGlCQUFpQixNQUFNO0FBQzVDLFlBQ0ksTUFBTSxpQkFBaUIsbUJBQW1CLEVBQUUsS0FBSyxNQUFNLFdBRW5ELE1BQU0sVUFBVSxXQUFXLE1BQU0sV0FBVyxJQUFJLE9BQU8sZUFDcEQsTUFBTSxVQUFVLFdBQVcsTUFBTSxVQUFVLElBQUksT0FBTztBQUdyRTtBQUVBLFNBQVMsVUFBVSxPQUFtQjtBQUVsQyxZQUFVO0FBQ1YsYUFBVztBQUNYLGNBQVk7QUFDWixhQUFXO0FBQ2Y7QUFFQSxJQUFNLGdCQUFnQixPQUFPLE9BQU87QUFBQSxFQUNoQyxhQUFhO0FBQUEsRUFDYixhQUFhO0FBQUEsRUFDYixhQUFhO0FBQUEsRUFDYixhQUFhO0FBQUEsRUFDYixZQUFZO0FBQUEsRUFDWixZQUFZO0FBQUEsRUFDWixZQUFZO0FBQUEsRUFDWixZQUFZO0FBQ2hCLENBQUM7QUFFRCxTQUFTLFVBQVUsTUFBeUM7QUFDeEQsTUFBSSxNQUFNO0FBQ04sUUFBSSxDQUFDLFlBQVk7QUFBRSxzQkFBZ0IsU0FBUyxLQUFLLE1BQU07QUFBQSxJQUFRO0FBQy9ELGFBQVMsS0FBSyxNQUFNLFNBQVMsY0FBYyxJQUFJO0FBQUEsRUFDbkQsV0FBVyxDQUFDLFFBQVEsWUFBWTtBQUM1QixhQUFTLEtBQUssTUFBTSxTQUFTO0FBQUEsRUFDakM7QUFFQSxlQUFhLFFBQVE7QUFDekI7QUFFQSxTQUFTLFlBQVksT0FBeUI7QUFDMUMsTUFBSSxhQUFhLFlBQVk7QUFFekIsZUFBVztBQUNYLFdBQU8sa0JBQWtCLFVBQVU7QUFBQSxFQUN2QyxXQUFXLFNBQVM7QUFFaEIsZUFBVztBQUNYLFdBQU8sWUFBWTtBQUFBLEVBQ3ZCO0FBRUEsTUFBSSxZQUFZLFVBQVU7QUFHdEIsY0FBVSxZQUFZO0FBQ3RCO0FBQUEsRUFDSjtBQUVBLE1BQUksQ0FBQyxhQUFhLENBQUMsVUFBVSxHQUFHO0FBQzVCLFFBQUksWUFBWTtBQUFFLGdCQUFVO0FBQUEsSUFBRztBQUMvQjtBQUFBLEVBQ0o7QUFFQSxRQUFNLHFCQUFxQixRQUFRLDJCQUEyQixLQUFLO0FBQ25FLFFBQU0sb0JBQW9CLFFBQVEsMEJBQTBCLEtBQUs7QUFHakUsUUFBTSxjQUFjLFFBQVEsbUJBQW1CLEtBQUs7QUFFcEQsUUFBTSxjQUFlLE9BQU8sYUFBYSxNQUFNLFVBQVc7QUFDMUQsUUFBTSxhQUFhLE1BQU0sVUFBVTtBQUNuQyxRQUFNLFlBQVksTUFBTSxVQUFVO0FBQ2xDLFFBQU0sZUFBZ0IsT0FBTyxjQUFjLE1BQU0sVUFBVztBQUc1RCxRQUFNLGNBQWUsT0FBTyxhQUFhLE1BQU0sVUFBWSxvQkFBb0I7QUFDL0UsUUFBTSxhQUFhLE1BQU0sVUFBVyxvQkFBb0I7QUFDeEQsUUFBTSxZQUFZLE1BQU0sVUFBVyxxQkFBcUI7QUFDeEQsUUFBTSxlQUFnQixPQUFPLGNBQWMsTUFBTSxVQUFZLHFCQUFxQjtBQUVsRixNQUFJLENBQUMsY0FBYyxDQUFDLGFBQWEsQ0FBQyxnQkFBZ0IsQ0FBQyxhQUFhO0FBRTVELGNBQVU7QUFBQSxFQUNkLFdBRVMsZUFBZSxhQUFjLFdBQVUsV0FBVztBQUFBLFdBQ2xELGNBQWMsYUFBYyxXQUFVLFdBQVc7QUFBQSxXQUNqRCxjQUFjLFVBQVcsV0FBVSxXQUFXO0FBQUEsV0FDOUMsYUFBYSxZQUFhLFdBQVUsV0FBVztBQUFBLFdBRS9DLFdBQVksV0FBVSxVQUFVO0FBQUEsV0FDaEMsVUFBVyxXQUFVLFVBQVU7QUFBQSxXQUMvQixhQUFjLFdBQVUsVUFBVTtBQUFBLFdBQ2xDLFlBQWEsV0FBVSxVQUFVO0FBQUEsTUFFckMsV0FBVTtBQUNuQjs7O0FDNU9BO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQVdBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksV0FBVztBQUVyRCxJQUFNQyxjQUFhO0FBQ25CLElBQU1DLGNBQWE7QUFDbkIsSUFBTSxhQUFhO0FBS1osU0FBUyxPQUFzQjtBQUNsQyxTQUFPRixNQUFLQyxXQUFVO0FBQzFCO0FBS08sU0FBUyxPQUFzQjtBQUNsQyxTQUFPRCxNQUFLRSxXQUFVO0FBQzFCO0FBS08sU0FBUyxPQUFzQjtBQUNsQyxTQUFPRixNQUFLLFVBQVU7QUFDMUI7OztBQ3BDQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDd0JBLElBQUksVUFBVSxTQUFTLFVBQVU7QUFDakMsSUFBSSxlQUFvRCxPQUFPLFlBQVksWUFBWSxZQUFZLFFBQVEsUUFBUTtBQUNuSCxJQUFJO0FBQ0osSUFBSTtBQUNKLElBQUksT0FBTyxpQkFBaUIsY0FBYyxPQUFPLE9BQU8sbUJBQW1CLFlBQVk7QUFDbkYsTUFBSTtBQUNBLG1CQUFlLE9BQU8sZUFBZSxDQUFDLEdBQUcsVUFBVTtBQUFBLE1BQy9DLEtBQUssV0FBWTtBQUNiLGNBQU07QUFBQSxNQUNWO0FBQUEsSUFDSixDQUFDO0FBQ0QsdUJBQW1CLENBQUM7QUFFcEIsaUJBQWEsV0FBWTtBQUFFLFlBQU07QUFBQSxJQUFJLEdBQUcsTUFBTSxZQUFZO0FBQUEsRUFDOUQsU0FBUyxHQUFHO0FBQ1IsUUFBSSxNQUFNLGtCQUFrQjtBQUN4QixxQkFBZTtBQUFBLElBQ25CO0FBQUEsRUFDSjtBQUNKLE9BQU87QUFDSCxpQkFBZTtBQUNuQjtBQUVBLElBQUksbUJBQW1CO0FBQ3ZCLElBQUksZUFBZSxTQUFTLG1CQUFtQixPQUFxQjtBQUNoRSxNQUFJO0FBQ0EsUUFBSSxRQUFRLFFBQVEsS0FBSyxLQUFLO0FBQzlCLFdBQU8saUJBQWlCLEtBQUssS0FBSztBQUFBLEVBQ3RDLFNBQVMsR0FBRztBQUNSLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFFQSxJQUFJLG9CQUFvQixTQUFTLGlCQUFpQixPQUFxQjtBQUNuRSxNQUFJO0FBQ0EsUUFBSSxhQUFhLEtBQUssR0FBRztBQUFFLGFBQU87QUFBQSxJQUFPO0FBQ3pDLFlBQVEsS0FBSyxLQUFLO0FBQ2xCLFdBQU87QUFBQSxFQUNYLFNBQVMsR0FBRztBQUNSLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFDQSxJQUFJLFFBQVEsT0FBTyxVQUFVO0FBQzdCLElBQUksY0FBYztBQUNsQixJQUFJLFVBQVU7QUFDZCxJQUFJLFdBQVc7QUFDZixJQUFJLFdBQVc7QUFDZixJQUFJLFlBQVk7QUFDaEIsSUFBSSxZQUFZO0FBQ2hCLElBQUksaUJBQWlCLE9BQU8sV0FBVyxjQUFjLENBQUMsQ0FBQyxPQUFPO0FBRTlELElBQUksU0FBUyxFQUFFLEtBQUssQ0FBQyxDQUFDO0FBRXRCLElBQUksUUFBaUMsU0FBUyxtQkFBbUI7QUFBRSxTQUFPO0FBQU87QUFDakYsSUFBSSxPQUFPLGFBQWEsVUFBVTtBQUUxQixRQUFNLFNBQVM7QUFDbkIsTUFBSSxNQUFNLEtBQUssR0FBRyxNQUFNLE1BQU0sS0FBSyxTQUFTLEdBQUcsR0FBRztBQUM5QyxZQUFRLFNBQVNHLGtCQUFpQixPQUFPO0FBR3JDLFdBQUssVUFBVSxDQUFDLFdBQVcsT0FBTyxVQUFVLGVBQWUsT0FBTyxVQUFVLFdBQVc7QUFDbkYsWUFBSTtBQUNBLGNBQUksTUFBTSxNQUFNLEtBQUssS0FBSztBQUMxQixrQkFDSSxRQUFRLFlBQ0wsUUFBUSxhQUNSLFFBQVEsYUFDUixRQUFRLGdCQUNWLE1BQU0sRUFBRSxLQUFLO0FBQUEsUUFDdEIsU0FBUyxHQUFHO0FBQUEsUUFBTztBQUFBLE1BQ3ZCO0FBQ0EsYUFBTztBQUFBLElBQ1g7QUFBQSxFQUNKO0FBQ0o7QUFuQlE7QUFxQlIsU0FBUyxtQkFBc0IsT0FBdUQ7QUFDbEYsTUFBSSxNQUFNLEtBQUssR0FBRztBQUFFLFdBQU87QUFBQSxFQUFNO0FBQ2pDLE1BQUksQ0FBQyxPQUFPO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDNUIsTUFBSSxPQUFPLFVBQVUsY0FBYyxPQUFPLFVBQVUsVUFBVTtBQUFFLFdBQU87QUFBQSxFQUFPO0FBQzlFLE1BQUk7QUFDQSxJQUFDLGFBQXFCLE9BQU8sTUFBTSxZQUFZO0FBQUEsRUFDbkQsU0FBUyxHQUFHO0FBQ1IsUUFBSSxNQUFNLGtCQUFrQjtBQUFFLGFBQU87QUFBQSxJQUFPO0FBQUEsRUFDaEQ7QUFDQSxTQUFPLENBQUMsYUFBYSxLQUFLLEtBQUssa0JBQWtCLEtBQUs7QUFDMUQ7QUFFQSxTQUFTLHFCQUF3QixPQUFzRDtBQUNuRixNQUFJLE1BQU0sS0FBSyxHQUFHO0FBQUUsV0FBTztBQUFBLEVBQU07QUFDakMsTUFBSSxDQUFDLE9BQU87QUFBRSxXQUFPO0FBQUEsRUFBTztBQUM1QixNQUFJLE9BQU8sVUFBVSxjQUFjLE9BQU8sVUFBVSxVQUFVO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDOUUsTUFBSSxnQkFBZ0I7QUFBRSxXQUFPLGtCQUFrQixLQUFLO0FBQUEsRUFBRztBQUN2RCxNQUFJLGFBQWEsS0FBSyxHQUFHO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDekMsTUFBSSxXQUFXLE1BQU0sS0FBSyxLQUFLO0FBQy9CLE1BQUksYUFBYSxXQUFXLGFBQWEsWUFBWSxDQUFFLGlCQUFrQixLQUFLLFFBQVEsR0FBRztBQUFFLFdBQU87QUFBQSxFQUFPO0FBQ3pHLFNBQU8sa0JBQWtCLEtBQUs7QUFDbEM7QUFFQSxJQUFPLG1CQUFRLGVBQWUscUJBQXFCOzs7QUN6RzVDLElBQU0sY0FBTixjQUEwQixNQUFNO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBTW5DLFlBQVksU0FBa0IsU0FBd0I7QUFDbEQsVUFBTSxTQUFTLE9BQU87QUFDdEIsU0FBSyxPQUFPO0FBQUEsRUFDaEI7QUFDSjtBQWNPLElBQU0sMEJBQU4sY0FBc0MsTUFBTTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFhL0MsWUFBWSxTQUFzQyxRQUFjLE1BQWU7QUFDM0UsV0FBTyxzQkFBUSwrQ0FBK0MsY0FBYyxhQUFhLE1BQU0sR0FBRyxFQUFFLE9BQU8sT0FBTyxDQUFDO0FBQ25ILFNBQUssVUFBVTtBQUNmLFNBQUssT0FBTztBQUFBLEVBQ2hCO0FBQ0o7QUErQkEsSUFBTSxhQUFhLE9BQU8sU0FBUztBQUNuQyxJQUFNLGdCQUFnQixPQUFPLFlBQVk7QUE3RnpDO0FBOEZBLElBQU0sV0FBVSxZQUFPLFlBQVAsWUFBa0IsT0FBTyxpQkFBaUI7QUFvRG5ELElBQU0scUJBQU4sTUFBTSw0QkFBOEIsUUFBZ0U7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUF1Q3ZHLFlBQVksVUFBeUMsYUFBMkM7QUFDNUYsUUFBSTtBQUNKLFFBQUk7QUFDSixVQUFNLENBQUMsS0FBSyxRQUFRO0FBQUUsZ0JBQVU7QUFBSyxlQUFTO0FBQUEsSUFBSyxDQUFDO0FBRXBELFFBQUssS0FBSyxZQUFvQixPQUFPLE1BQU0sU0FBUztBQUNoRCxZQUFNLElBQUksVUFBVSxtSUFBbUk7QUFBQSxJQUMzSjtBQUVBLFFBQUksVUFBOEM7QUFBQSxNQUM5QyxTQUFTO0FBQUEsTUFDVDtBQUFBLE1BQ0E7QUFBQSxNQUNBLElBQUksY0FBYztBQUFFLGVBQU8sb0NBQWU7QUFBQSxNQUFNO0FBQUEsTUFDaEQsSUFBSSxZQUFZLElBQUk7QUFBRSxzQkFBYyxrQkFBTTtBQUFBLE1BQVc7QUFBQSxJQUN6RDtBQUVBLFVBQU0sUUFBaUM7QUFBQSxNQUNuQyxJQUFJLE9BQU87QUFBRSxlQUFPO0FBQUEsTUFBTztBQUFBLE1BQzNCLFdBQVc7QUFBQSxNQUNYLFNBQVM7QUFBQSxJQUNiO0FBR0EsU0FBSyxPQUFPLGlCQUFpQixNQUFNO0FBQUEsTUFDL0IsQ0FBQyxVQUFVLEdBQUc7QUFBQSxRQUNWLGNBQWM7QUFBQSxRQUNkLFlBQVk7QUFBQSxRQUNaLFVBQVU7QUFBQSxRQUNWLE9BQU87QUFBQSxNQUNYO0FBQUEsTUFDQSxDQUFDLGFBQWEsR0FBRztBQUFBLFFBQ2IsY0FBYztBQUFBLFFBQ2QsWUFBWTtBQUFBLFFBQ1osVUFBVTtBQUFBLFFBQ1YsT0FBTyxhQUFhLFNBQVMsS0FBSztBQUFBLE1BQ3RDO0FBQUEsSUFDSixDQUFDO0FBR0QsVUFBTSxXQUFXLFlBQVksU0FBUyxLQUFLO0FBQzNDLFFBQUk7QUFDQSxlQUFTLFlBQVksU0FBUyxLQUFLLEdBQUcsUUFBUTtBQUFBLElBQ2xELFNBQVMsS0FBSztBQUNWLFVBQUksTUFBTSxXQUFXO0FBQ2pCLGdCQUFRLElBQUksdURBQXVELEdBQUc7QUFBQSxNQUMxRSxPQUFPO0FBQ0gsaUJBQVMsR0FBRztBQUFBLE1BQ2hCO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBeURBLE9BQU8sT0FBdUM7QUFDMUMsV0FBTyxJQUFJLG9CQUF5QixDQUFDLFlBQVk7QUFHN0MsY0FBUSxJQUFJO0FBQUEsUUFDUixLQUFLLGFBQWEsRUFBRSxJQUFJLFlBQVksc0JBQXNCLEVBQUUsTUFBTSxDQUFDLENBQUM7QUFBQSxRQUNwRSxlQUFlLElBQUk7QUFBQSxNQUN2QixDQUFDLEVBQUUsS0FBSyxNQUFNLFFBQVEsR0FBRyxNQUFNLFFBQVEsQ0FBQztBQUFBLElBQzVDLENBQUM7QUFBQSxFQUNMO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQTJCQSxTQUFTLFFBQTRDO0FBQ2pELFFBQUksT0FBTyxTQUFTO0FBQ2hCLFdBQUssS0FBSyxPQUFPLE9BQU8sTUFBTTtBQUFBLElBQ2xDLE9BQU87QUFDSCxhQUFPLGlCQUFpQixTQUFTLE1BQU0sS0FBSyxLQUFLLE9BQU8sT0FBTyxNQUFNLEdBQUcsRUFBQyxTQUFTLEtBQUksQ0FBQztBQUFBLElBQzNGO0FBRUEsV0FBTztBQUFBLEVBQ1g7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUErQkEsS0FBcUMsYUFBc0gsWUFBd0gsYUFBb0Y7QUFDblcsUUFBSSxFQUFFLGdCQUFnQixzQkFBcUI7QUFDdkMsWUFBTSxJQUFJLFVBQVUsZ0VBQWdFO0FBQUEsSUFDeEY7QUFNQSxRQUFJLENBQUMsaUJBQVcsV0FBVyxHQUFHO0FBQUUsb0JBQWM7QUFBQSxJQUFpQjtBQUMvRCxRQUFJLENBQUMsaUJBQVcsVUFBVSxHQUFHO0FBQUUsbUJBQWE7QUFBQSxJQUFTO0FBRXJELFFBQUksZ0JBQWdCLFlBQVksY0FBYyxTQUFTO0FBRW5ELGFBQU8sSUFBSSxvQkFBbUIsQ0FBQyxZQUFZLFFBQVEsSUFBVyxDQUFDO0FBQUEsSUFDbkU7QUFFQSxVQUFNLFVBQStDLENBQUM7QUFDdEQsU0FBSyxVQUFVLElBQUk7QUFFbkIsV0FBTyxJQUFJLG9CQUF3QyxDQUFDLFNBQVMsV0FBVztBQUNwRSxXQUFLLE1BQU07QUFBQSxRQUNQLENBQUMsVUFBVTtBQXJZM0IsY0FBQUM7QUFzWW9CLGNBQUksS0FBSyxVQUFVLE1BQU0sU0FBUztBQUFFLGlCQUFLLFVBQVUsSUFBSTtBQUFBLFVBQU07QUFDN0QsV0FBQUEsTUFBQSxRQUFRLFlBQVIsZ0JBQUFBLElBQUE7QUFFQSxjQUFJO0FBQ0Esb0JBQVEsWUFBYSxLQUFLLENBQUM7QUFBQSxVQUMvQixTQUFTLEtBQUs7QUFDVixtQkFBTyxHQUFHO0FBQUEsVUFDZDtBQUFBLFFBQ0o7QUFBQSxRQUNBLENBQUMsV0FBWTtBQS9ZN0IsY0FBQUE7QUFnWm9CLGNBQUksS0FBSyxVQUFVLE1BQU0sU0FBUztBQUFFLGlCQUFLLFVBQVUsSUFBSTtBQUFBLFVBQU07QUFDN0QsV0FBQUEsTUFBQSxRQUFRLFlBQVIsZ0JBQUFBLElBQUE7QUFFQSxjQUFJO0FBQ0Esb0JBQVEsV0FBWSxNQUFNLENBQUM7QUFBQSxVQUMvQixTQUFTLEtBQUs7QUFDVixtQkFBTyxHQUFHO0FBQUEsVUFDZDtBQUFBLFFBQ0o7QUFBQSxNQUNKO0FBQUEsSUFDSixHQUFHLE9BQU8sVUFBVztBQUVqQixVQUFJO0FBQ0EsZUFBTywyQ0FBYztBQUFBLE1BQ3pCLFVBQUU7QUFDRSxjQUFNLEtBQUssT0FBTyxLQUFLO0FBQUEsTUFDM0I7QUFBQSxJQUNKLENBQUM7QUFBQSxFQUNMO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBK0JBLE1BQXVCLFlBQXFGLGFBQTRFO0FBQ3BMLFdBQU8sS0FBSyxLQUFLLFFBQVcsWUFBWSxXQUFXO0FBQUEsRUFDdkQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBaUNBLFFBQVEsV0FBNkMsYUFBa0U7QUFDbkgsUUFBSSxFQUFFLGdCQUFnQixzQkFBcUI7QUFDdkMsWUFBTSxJQUFJLFVBQVUsbUVBQW1FO0FBQUEsSUFDM0Y7QUFFQSxRQUFJLENBQUMsaUJBQVcsU0FBUyxHQUFHO0FBQ3hCLGFBQU8sS0FBSyxLQUFLLFdBQVcsV0FBVyxXQUFXO0FBQUEsSUFDdEQ7QUFFQSxXQUFPLEtBQUs7QUFBQSxNQUNSLENBQUMsVUFBVSxvQkFBbUIsUUFBUSxVQUFVLENBQUMsRUFBRSxLQUFLLE1BQU0sS0FBSztBQUFBLE1BQ25FLENBQUMsV0FBWSxvQkFBbUIsUUFBUSxVQUFVLENBQUMsRUFBRSxLQUFLLE1BQU07QUFBRSxjQUFNO0FBQUEsTUFBUSxDQUFDO0FBQUEsTUFDakY7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFZQSxhQXpXUyxZQUVTLGVBdVdOLFFBQU8sSUFBSTtBQUNuQixXQUFPO0FBQUEsRUFDWDtBQUFBLEVBYUEsT0FBTyxJQUFzRCxRQUF3QztBQUNqRyxRQUFJLFlBQVksTUFBTSxLQUFLLE1BQU07QUFDakMsVUFBTSxVQUFVLFVBQVUsV0FBVyxJQUMvQixvQkFBbUIsUUFBUSxTQUFTLElBQ3BDLElBQUksb0JBQTRCLENBQUMsU0FBUyxXQUFXO0FBQ25ELFdBQUssUUFBUSxJQUFJLFNBQVMsRUFBRSxLQUFLLFNBQVMsTUFBTTtBQUFBLElBQ3BELEdBQUcsQ0FBQyxVQUEwQixVQUFVLFNBQVMsV0FBVyxLQUFLLENBQUM7QUFDdEUsV0FBTztBQUFBLEVBQ1g7QUFBQSxFQWFBLE9BQU8sV0FBNkQsUUFBd0M7QUFDeEcsUUFBSSxZQUFZLE1BQU0sS0FBSyxNQUFNO0FBQ2pDLFVBQU0sVUFBVSxVQUFVLFdBQVcsSUFDL0Isb0JBQW1CLFFBQVEsU0FBUyxJQUNwQyxJQUFJLG9CQUE0QixDQUFDLFNBQVMsV0FBVztBQUNuRCxXQUFLLFFBQVEsV0FBVyxTQUFTLEVBQUUsS0FBSyxTQUFTLE1BQU07QUFBQSxJQUMzRCxHQUFHLENBQUMsVUFBMEIsVUFBVSxTQUFTLFdBQVcsS0FBSyxDQUFDO0FBQ3RFLFdBQU87QUFBQSxFQUNYO0FBQUEsRUFlQSxPQUFPLElBQXNELFFBQXdDO0FBQ2pHLFFBQUksWUFBWSxNQUFNLEtBQUssTUFBTTtBQUNqQyxVQUFNLFVBQVUsVUFBVSxXQUFXLElBQy9CLG9CQUFtQixRQUFRLFNBQVMsSUFDcEMsSUFBSSxvQkFBNEIsQ0FBQyxTQUFTLFdBQVc7QUFDbkQsV0FBSyxRQUFRLElBQUksU0FBUyxFQUFFLEtBQUssU0FBUyxNQUFNO0FBQUEsSUFDcEQsR0FBRyxDQUFDLFVBQTBCLFVBQVUsU0FBUyxXQUFXLEtBQUssQ0FBQztBQUN0RSxXQUFPO0FBQUEsRUFDWDtBQUFBLEVBWUEsT0FBTyxLQUF1RCxRQUF3QztBQUNsRyxRQUFJLFlBQVksTUFBTSxLQUFLLE1BQU07QUFDakMsVUFBTSxVQUFVLElBQUksb0JBQTRCLENBQUMsU0FBUyxXQUFXO0FBQ2pFLFdBQUssUUFBUSxLQUFLLFNBQVMsRUFBRSxLQUFLLFNBQVMsTUFBTTtBQUFBLElBQ3JELEdBQUcsQ0FBQyxVQUEwQixVQUFVLFNBQVMsV0FBVyxLQUFLLENBQUM7QUFDbEUsV0FBTztBQUFBLEVBQ1g7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxPQUFPLE9BQWtCLE9BQW9DO0FBQ3pELFVBQU0sSUFBSSxJQUFJLG9CQUFzQixNQUFNO0FBQUEsSUFBQyxDQUFDO0FBQzVDLE1BQUUsT0FBTyxLQUFLO0FBQ2QsV0FBTztBQUFBLEVBQ1g7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBWUEsT0FBTyxRQUFtQixjQUFzQixPQUFvQztBQUNoRixVQUFNLFVBQVUsSUFBSSxvQkFBc0IsTUFBTTtBQUFBLElBQUMsQ0FBQztBQUNsRCxRQUFJLGVBQWUsT0FBTyxnQkFBZ0IsY0FBYyxZQUFZLFdBQVcsT0FBTyxZQUFZLFlBQVksWUFBWTtBQUN0SCxrQkFBWSxRQUFRLFlBQVksRUFBRSxpQkFBaUIsU0FBUyxNQUFNLEtBQUssUUFBUSxPQUFPLEtBQUssQ0FBQztBQUFBLElBQ2hHLE9BQU87QUFDSCxpQkFBVyxNQUFNLEtBQUssUUFBUSxPQUFPLEtBQUssR0FBRyxZQUFZO0FBQUEsSUFDN0Q7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUFBLEVBaUJBLE9BQU8sTUFBZ0IsY0FBc0IsT0FBa0M7QUFDM0UsV0FBTyxJQUFJLG9CQUFzQixDQUFDLFlBQVk7QUFDMUMsaUJBQVcsTUFBTSxRQUFRLEtBQU0sR0FBRyxZQUFZO0FBQUEsSUFDbEQsQ0FBQztBQUFBLEVBQ0w7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxPQUFPLE9BQWtCLFFBQXFDO0FBQzFELFdBQU8sSUFBSSxvQkFBc0IsQ0FBQyxHQUFHLFdBQVcsT0FBTyxNQUFNLENBQUM7QUFBQSxFQUNsRTtBQUFBLEVBb0JBLE9BQU8sUUFBa0IsT0FBNEQ7QUFDakYsUUFBSSxpQkFBaUIscUJBQW9CO0FBRXJDLGFBQU87QUFBQSxJQUNYO0FBQ0EsV0FBTyxJQUFJLG9CQUF3QixDQUFDLFlBQVksUUFBUSxLQUFLLENBQUM7QUFBQSxFQUNsRTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVVBLE9BQU8sZ0JBQXVEO0FBQzFELFFBQUksU0FBNkMsRUFBRSxhQUFhLEtBQUs7QUFDckUsV0FBTyxVQUFVLElBQUksb0JBQXNCLENBQUMsU0FBUyxXQUFXO0FBQzVELGFBQU8sVUFBVTtBQUNqQixhQUFPLFNBQVM7QUFBQSxJQUNwQixHQUFHLENBQUMsVUFBZ0I7QUF6ckI1QixVQUFBQTtBQXlyQjhCLE9BQUFBLE1BQUEsT0FBTyxnQkFBUCxnQkFBQUEsSUFBQSxhQUFxQjtBQUFBLElBQVEsQ0FBQztBQUNwRCxXQUFPO0FBQUEsRUFDWDtBQUNKO0FBTUEsU0FBUyxhQUFnQixTQUE2QyxPQUFnQztBQUNsRyxNQUFJLHNCQUFnRDtBQUVwRCxTQUFPLENBQUMsV0FBa0Q7QUFDdEQsUUFBSSxDQUFDLE1BQU0sU0FBUztBQUNoQixZQUFNLFVBQVU7QUFDaEIsWUFBTSxTQUFTO0FBQ2YsY0FBUSxPQUFPLE1BQU07QUFNckIsV0FBSyxRQUFRLFVBQVUsS0FBSyxLQUFLLFFBQVEsU0FBUyxRQUFXLENBQUMsUUFBUTtBQUNsRSxZQUFJLFFBQVEsUUFBUTtBQUNoQixnQkFBTTtBQUFBLFFBQ1Y7QUFBQSxNQUNKLENBQUM7QUFBQSxJQUNMO0FBSUEsUUFBSSxDQUFDLE1BQU0sVUFBVSxDQUFDLFFBQVEsYUFBYTtBQUFFO0FBQUEsSUFBUTtBQUVyRCwwQkFBc0IsSUFBSSxRQUFjLENBQUMsWUFBWTtBQUNqRCxVQUFJO0FBQ0EsZ0JBQVEsUUFBUSxZQUFhLE1BQU0sT0FBUSxLQUFLLENBQUM7QUFBQSxNQUNyRCxTQUFTLEtBQUs7QUFDVixnQkFBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsU0FBUyxLQUFLLDhDQUE4QyxDQUFDO0FBQUEsTUFDcEg7QUFBQSxJQUNKLENBQUMsRUFBRSxNQUFNLENBQUNDLFlBQVk7QUFDbEIsY0FBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsU0FBU0EsU0FBUSw4Q0FBOEMsQ0FBQztBQUFBLElBQ3ZILENBQUM7QUFHRCxZQUFRLGNBQWM7QUFFdEIsV0FBTztBQUFBLEVBQ1g7QUFDSjtBQUtBLFNBQVMsWUFBZSxTQUE2QyxPQUErRDtBQUNoSSxTQUFPLENBQUMsVUFBVTtBQUNkLFFBQUksTUFBTSxXQUFXO0FBQUU7QUFBQSxJQUFRO0FBQy9CLFVBQU0sWUFBWTtBQUVsQixRQUFJLFVBQVUsUUFBUSxTQUFTO0FBQzNCLFVBQUksTUFBTSxTQUFTO0FBQUU7QUFBQSxNQUFRO0FBQzdCLFlBQU0sVUFBVTtBQUNoQixjQUFRLE9BQU8sSUFBSSxVQUFVLDJDQUEyQyxDQUFDO0FBQ3pFO0FBQUEsSUFDSjtBQUVBLFFBQUksU0FBUyxTQUFTLE9BQU8sVUFBVSxZQUFZLE9BQU8sVUFBVSxhQUFhO0FBQzdFLFVBQUk7QUFDSixVQUFJO0FBQ0EsZUFBUSxNQUFjO0FBQUEsTUFDMUIsU0FBUyxLQUFLO0FBQ1YsY0FBTSxVQUFVO0FBQ2hCLGdCQUFRLE9BQU8sR0FBRztBQUNsQjtBQUFBLE1BQ0o7QUFFQSxVQUFJLGlCQUFXLElBQUksR0FBRztBQUNsQixZQUFJO0FBQ0EsY0FBSSxTQUFVLE1BQWM7QUFDNUIsY0FBSSxpQkFBVyxNQUFNLEdBQUc7QUFDcEIsa0JBQU0sY0FBYyxDQUFDLFVBQWdCO0FBQ2pDLHNCQUFRLE1BQU0sUUFBUSxPQUFPLENBQUMsS0FBSyxDQUFDO0FBQUEsWUFDeEM7QUFDQSxnQkFBSSxNQUFNLFFBQVE7QUFJZCxtQkFBSyxhQUFhLGlDQUFLLFVBQUwsRUFBYyxZQUFZLElBQUcsS0FBSyxFQUFFLE1BQU0sTUFBTTtBQUFBLFlBQ3RFLE9BQU87QUFDSCxzQkFBUSxjQUFjO0FBQUEsWUFDMUI7QUFBQSxVQUNKO0FBQUEsUUFDSixTQUFRO0FBQUEsUUFBQztBQUVULGNBQU0sV0FBb0M7QUFBQSxVQUN0QyxNQUFNLE1BQU07QUFBQSxVQUNaLFdBQVc7QUFBQSxVQUNYLElBQUksVUFBVTtBQUFFLG1CQUFPLEtBQUssS0FBSztBQUFBLFVBQVE7QUFBQSxVQUN6QyxJQUFJLFFBQVFDLFFBQU87QUFBRSxpQkFBSyxLQUFLLFVBQVVBO0FBQUEsVUFBTztBQUFBLFVBQ2hELElBQUksU0FBUztBQUFFLG1CQUFPLEtBQUssS0FBSztBQUFBLFVBQU87QUFBQSxRQUMzQztBQUVBLGNBQU0sV0FBVyxZQUFZLFNBQVMsUUFBUTtBQUM5QyxZQUFJO0FBQ0Esa0JBQVEsTUFBTSxNQUFNLE9BQU8sQ0FBQyxZQUFZLFNBQVMsUUFBUSxHQUFHLFFBQVEsQ0FBQztBQUFBLFFBQ3pFLFNBQVMsS0FBSztBQUNWLG1CQUFTLEdBQUc7QUFBQSxRQUNoQjtBQUNBO0FBQUEsTUFDSjtBQUFBLElBQ0o7QUFFQSxRQUFJLE1BQU0sU0FBUztBQUFFO0FBQUEsSUFBUTtBQUM3QixVQUFNLFVBQVU7QUFDaEIsWUFBUSxRQUFRLEtBQUs7QUFBQSxFQUN6QjtBQUNKO0FBS0EsU0FBUyxZQUFlLFNBQTZDLE9BQTREO0FBQzdILFNBQU8sQ0FBQyxXQUFZO0FBQ2hCLFFBQUksTUFBTSxXQUFXO0FBQUU7QUFBQSxJQUFRO0FBQy9CLFVBQU0sWUFBWTtBQUVsQixRQUFJLE1BQU0sU0FBUztBQUNmLFVBQUk7QUFDQSxZQUFJLGtCQUFrQixlQUFlLE1BQU0sa0JBQWtCLGVBQWUsT0FBTyxHQUFHLE9BQU8sT0FBTyxNQUFNLE9BQU8sS0FBSyxHQUFHO0FBRXJIO0FBQUEsUUFDSjtBQUFBLE1BQ0osU0FBUTtBQUFBLE1BQUM7QUFFVCxXQUFLLFFBQVEsT0FBTyxJQUFJLHdCQUF3QixRQUFRLFNBQVMsTUFBTSxDQUFDO0FBQUEsSUFDNUUsT0FBTztBQUNILFlBQU0sVUFBVTtBQUNoQixjQUFRLE9BQU8sTUFBTTtBQUFBLElBQ3pCO0FBQUEsRUFDSjtBQUNKO0FBTUEsU0FBUyxVQUFVLFFBQXFDLFFBQWUsT0FBNEI7QUFDL0YsUUFBTSxVQUFVLENBQUM7QUFFakIsYUFBVyxTQUFTLFFBQVE7QUFDeEIsUUFBSTtBQUNKLFFBQUk7QUFDQSxVQUFJLENBQUMsaUJBQVcsTUFBTSxJQUFJLEdBQUc7QUFBRTtBQUFBLE1BQVU7QUFDekMsZUFBUyxNQUFNO0FBQ2YsVUFBSSxDQUFDLGlCQUFXLE1BQU0sR0FBRztBQUFFO0FBQUEsTUFBVTtBQUFBLElBQ3pDLFNBQVE7QUFBRTtBQUFBLElBQVU7QUFFcEIsUUFBSTtBQUNKLFFBQUk7QUFDQSxlQUFTLFFBQVEsTUFBTSxRQUFRLE9BQU8sQ0FBQyxLQUFLLENBQUM7QUFBQSxJQUNqRCxTQUFTLEtBQUs7QUFDVixjQUFRLE9BQU8sSUFBSSx3QkFBd0IsUUFBUSxLQUFLLHVDQUF1QyxDQUFDO0FBQ2hHO0FBQUEsSUFDSjtBQUVBLFFBQUksQ0FBQyxRQUFRO0FBQUU7QUFBQSxJQUFVO0FBQ3pCLFlBQVE7QUFBQSxPQUNILGtCQUFrQixVQUFXLFNBQVMsUUFBUSxRQUFRLE1BQU0sR0FBRyxNQUFNLENBQUMsV0FBWTtBQUMvRSxnQkFBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsUUFBUSx1Q0FBdUMsQ0FBQztBQUFBLE1BQ3ZHLENBQUM7QUFBQSxJQUNMO0FBQUEsRUFDSjtBQUVBLFNBQU8sUUFBUSxJQUFJLE9BQU87QUFDOUI7QUFLQSxTQUFTLFNBQVksR0FBUztBQUMxQixTQUFPO0FBQ1g7QUFLQSxTQUFTLFFBQVEsUUFBcUI7QUFDbEMsUUFBTTtBQUNWO0FBS0EsU0FBUyxhQUFhLEtBQWtCO0FBQ3BDLE1BQUk7QUFDQSxRQUFJLGVBQWUsU0FBUyxPQUFPLFFBQVEsWUFBWSxJQUFJLGFBQWEsT0FBTyxVQUFVLFVBQVU7QUFDL0YsYUFBTyxLQUFLO0FBQUEsSUFDaEI7QUFBQSxFQUNKLFNBQVE7QUFBQSxFQUFDO0FBRVQsTUFBSTtBQUNBLFdBQU8sS0FBSyxVQUFVLEdBQUc7QUFBQSxFQUM3QixTQUFRO0FBQUEsRUFBQztBQUVULE1BQUk7QUFDQSxXQUFPLE9BQU8sVUFBVSxTQUFTLEtBQUssR0FBRztBQUFBLEVBQzdDLFNBQVE7QUFBQSxFQUFDO0FBRVQsU0FBTztBQUNYO0FBS0EsU0FBUyxlQUFrQixTQUErQztBQTk0QjFFLE1BQUFGO0FBKzRCSSxNQUFJLE9BQTJDQSxNQUFBLFFBQVEsVUFBVSxNQUFsQixPQUFBQSxNQUF1QixDQUFDO0FBQ3ZFLE1BQUksRUFBRSxhQUFhLE1BQU07QUFDckIsV0FBTyxPQUFPLEtBQUsscUJBQTJCLENBQUM7QUFBQSxFQUNuRDtBQUNBLE1BQUksUUFBUSxVQUFVLEtBQUssTUFBTTtBQUM3QixRQUFJLFFBQVM7QUFDYixZQUFRLFVBQVUsSUFBSTtBQUFBLEVBQzFCO0FBQ0EsU0FBTyxJQUFJO0FBQ2Y7QUFHQSxJQUFJLHVCQUF1QixRQUFRO0FBQ25DLElBQUksd0JBQXdCLE9BQU8seUJBQXlCLFlBQVk7QUFDcEUseUJBQXVCLHFCQUFxQixLQUFLLE9BQU87QUFDNUQsT0FBTztBQUNILHlCQUF1QixXQUF3QztBQUMzRCxRQUFJO0FBQ0osUUFBSTtBQUNKLFVBQU0sVUFBVSxJQUFJLFFBQVcsQ0FBQyxLQUFLLFFBQVE7QUFBRSxnQkFBVTtBQUFLLGVBQVM7QUFBQSxJQUFLLENBQUM7QUFDN0UsV0FBTyxFQUFFLFNBQVMsU0FBUyxPQUFPO0FBQUEsRUFDdEM7QUFDSjs7O0FGdDVCQSxPQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFDbEMsT0FBTyxPQUFPLG9CQUFvQjtBQUNsQyxPQUFPLE9BQU8sbUJBQW1CO0FBSWpDLElBQU1HLFFBQU8saUJBQWlCLFlBQVksSUFBSTtBQUM5QyxJQUFNLGFBQWEsaUJBQWlCLFlBQVksVUFBVTtBQUMxRCxJQUFNLGdCQUFnQixvQkFBSSxJQUE4QjtBQUV4RCxJQUFNLGNBQWM7QUFDcEIsSUFBTSxlQUFlO0FBMEJkLElBQU0sZUFBTixjQUEyQixNQUFNO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBTXBDLFlBQVksU0FBa0IsU0FBd0I7QUFDbEQsVUFBTSxTQUFTLE9BQU87QUFDdEIsU0FBSyxPQUFPO0FBQUEsRUFDaEI7QUFDSjtBQVNBLFNBQVMsY0FBYyxJQUFZLE1BQWMsUUFBdUI7QUFDcEUsUUFBTSxZQUFZQyxzQkFBcUIsRUFBRTtBQUN6QyxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLE1BQUksQ0FBQyxNQUFNO0FBQ1AsY0FBVSxRQUFRLE1BQVM7QUFBQSxFQUMvQixXQUFXLENBQUMsUUFBUTtBQUNoQixjQUFVLFFBQVEsSUFBSTtBQUFBLEVBQzFCLE9BQU87QUFDSCxRQUFJO0FBQ0EsZ0JBQVUsUUFBUSxLQUFLLE1BQU0sSUFBSSxDQUFDO0FBQUEsSUFDdEMsU0FBUyxLQUFVO0FBQ2YsZ0JBQVUsT0FBTyxJQUFJLFVBQVUsNkJBQTZCLElBQUksU0FBUyxFQUFFLE9BQU8sSUFBSSxDQUFDLENBQUM7QUFBQSxJQUM1RjtBQUFBLEVBQ0o7QUFDSjtBQVNBLFNBQVMsYUFBYSxJQUFZLE1BQWMsUUFBdUI7QUFDbkUsUUFBTSxZQUFZQSxzQkFBcUIsRUFBRTtBQUN6QyxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLE1BQUksQ0FBQyxRQUFRO0FBQ1QsY0FBVSxPQUFPLElBQUksTUFBTSxJQUFJLENBQUM7QUFBQSxFQUNwQyxPQUFPO0FBQ0gsUUFBSTtBQUNKLFFBQUk7QUFDQSxjQUFRLEtBQUssTUFBTSxJQUFJO0FBQUEsSUFDM0IsU0FBUyxLQUFVO0FBQ2YsZ0JBQVUsT0FBTyxJQUFJLFVBQVUsNEJBQTRCLElBQUksU0FBUyxFQUFFLE9BQU8sSUFBSSxDQUFDLENBQUM7QUFDdkY7QUFBQSxJQUNKO0FBRUEsUUFBSSxVQUF3QixDQUFDO0FBQzdCLFFBQUksTUFBTSxPQUFPO0FBQ2IsY0FBUSxRQUFRLE1BQU07QUFBQSxJQUMxQjtBQUVBLFFBQUk7QUFDSixZQUFRLE1BQU0sTUFBTTtBQUFBLE1BQ2hCLEtBQUs7QUFDRCxvQkFBWSxJQUFJLGVBQWUsTUFBTSxTQUFTLE9BQU87QUFDckQ7QUFBQSxNQUNKLEtBQUs7QUFDRCxvQkFBWSxJQUFJLFVBQVUsTUFBTSxTQUFTLE9BQU87QUFDaEQ7QUFBQSxNQUNKLEtBQUs7QUFDRCxvQkFBWSxJQUFJLGFBQWEsTUFBTSxTQUFTLE9BQU87QUFDbkQ7QUFBQSxNQUNKO0FBQ0ksb0JBQVksSUFBSSxNQUFNLE1BQU0sU0FBUyxPQUFPO0FBQzVDO0FBQUEsSUFDUjtBQUVBLGNBQVUsT0FBTyxTQUFTO0FBQUEsRUFDOUI7QUFDSjtBQVFBLFNBQVNBLHNCQUFxQixJQUEwQztBQUNwRSxRQUFNLFdBQVcsY0FBYyxJQUFJLEVBQUU7QUFDckMsZ0JBQWMsT0FBTyxFQUFFO0FBQ3ZCLFNBQU87QUFDWDtBQU9BLFNBQVNDLGNBQXFCO0FBQzFCLE1BQUk7QUFDSixLQUFHO0FBQ0MsYUFBUyxPQUFPO0FBQUEsRUFDcEIsU0FBUyxjQUFjLElBQUksTUFBTTtBQUNqQyxTQUFPO0FBQ1g7QUFjTyxTQUFTLEtBQUssU0FBK0M7QUFDaEUsUUFBTSxLQUFLQSxZQUFXO0FBRXRCLFFBQU0sU0FBUyxtQkFBbUIsY0FBbUI7QUFDckQsZ0JBQWMsSUFBSSxJQUFJLEVBQUUsU0FBUyxPQUFPLFNBQVMsUUFBUSxPQUFPLE9BQU8sQ0FBQztBQUV4RSxRQUFNLFVBQVVGLE1BQUssYUFBYSxPQUFPLE9BQU8sRUFBRSxXQUFXLEdBQUcsR0FBRyxPQUFPLENBQUM7QUFDM0UsTUFBSSxVQUFVO0FBRWQsVUFBUSxLQUFLLE1BQU07QUFDZixjQUFVO0FBQUEsRUFDZCxHQUFHLENBQUMsUUFBUTtBQUNSLGtCQUFjLE9BQU8sRUFBRTtBQUN2QixXQUFPLE9BQU8sR0FBRztBQUFBLEVBQ3JCLENBQUM7QUFFRCxRQUFNLFNBQVMsTUFBTTtBQUNqQixrQkFBYyxPQUFPLEVBQUU7QUFDdkIsV0FBTyxXQUFXLGNBQWMsRUFBQyxXQUFXLEdBQUUsQ0FBQyxFQUFFLE1BQU0sQ0FBQyxRQUFRO0FBQzVELGNBQVEsTUFBTSxxREFBcUQsR0FBRztBQUFBLElBQzFFLENBQUM7QUFBQSxFQUNMO0FBRUEsU0FBTyxjQUFjLE1BQU07QUFDdkIsUUFBSSxTQUFTO0FBQ1QsYUFBTyxPQUFPO0FBQUEsSUFDbEIsT0FBTztBQUNILGFBQU8sUUFBUSxLQUFLLE1BQU07QUFBQSxJQUM5QjtBQUFBLEVBQ0o7QUFFQSxTQUFPLE9BQU87QUFDbEI7QUFVTyxTQUFTLE9BQU8sZUFBdUIsTUFBc0M7QUFDaEYsU0FBTyxLQUFLLEVBQUUsWUFBWSxLQUFLLENBQUM7QUFDcEM7QUFVTyxTQUFTLEtBQUssYUFBcUIsTUFBc0M7QUFDNUUsU0FBTyxLQUFLLEVBQUUsVUFBVSxLQUFLLENBQUM7QUFDbEM7OztBR3hPQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBWUEsSUFBTUcsUUFBTyxpQkFBaUIsWUFBWSxTQUFTO0FBRW5ELElBQU0sbUJBQW1CO0FBQ3pCLElBQU0sZ0JBQWdCO0FBUWYsU0FBUyxRQUFRLE1BQTZCO0FBQ2pELFNBQU9BLE1BQUssa0JBQWtCLEVBQUMsS0FBSSxDQUFDO0FBQ3hDO0FBT08sU0FBUyxPQUF3QjtBQUNwQyxTQUFPQSxNQUFLLGFBQWE7QUFDN0I7OztBQ2xDQTtBQUFBO0FBQUE7QUFBQSxlQUFBQztBQUFBLEVBQUE7QUFBQSxhQUFBQztBQUFBLEVBQUE7QUFBQTtBQUFBO0FBYU8sU0FBUyxJQUFhLFFBQWdCO0FBQ3pDLFNBQU87QUFDWDtBQU1PLFNBQVMsVUFBVSxRQUFxQjtBQUMzQyxTQUFTLFVBQVUsT0FBUSxLQUFLO0FBQ3BDO0FBT08sU0FBU0MsT0FBZSxTQUFtRDtBQUM5RSxNQUFJLFlBQVksS0FBSztBQUNqQixXQUFPLENBQUMsV0FBWSxXQUFXLE9BQU8sQ0FBQyxJQUFJO0FBQUEsRUFDL0M7QUFFQSxTQUFPLENBQUMsV0FBVztBQUNmLFFBQUksV0FBVyxNQUFNO0FBQ2pCLGFBQU8sQ0FBQztBQUFBLElBQ1o7QUFDQSxhQUFTLElBQUksR0FBRyxJQUFJLE9BQU8sUUFBUSxLQUFLO0FBQ3BDLGFBQU8sQ0FBQyxJQUFJLFFBQVEsT0FBTyxDQUFDLENBQUM7QUFBQSxJQUNqQztBQUNBLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFPTyxTQUFTQyxLQUFhLEtBQThCLE9BQStEO0FBQ3RILE1BQUksVUFBVSxLQUFLO0FBQ2YsV0FBTyxDQUFDLFdBQVksV0FBVyxPQUFPLENBQUMsSUFBSTtBQUFBLEVBQy9DO0FBRUEsU0FBTyxDQUFDLFdBQVc7QUFDZixRQUFJLFdBQVcsTUFBTTtBQUNqQixhQUFPLENBQUM7QUFBQSxJQUNaO0FBQ0EsZUFBV0MsUUFBTyxRQUFRO0FBQ3RCLGFBQU9BLElBQUcsSUFBSSxNQUFNLE9BQU9BLElBQUcsQ0FBQztBQUFBLElBQ25DO0FBQ0EsV0FBTztBQUFBLEVBQ1g7QUFDSjtBQU1PLFNBQVMsU0FBa0IsU0FBMEQ7QUFDeEYsTUFBSSxZQUFZLEtBQUs7QUFDakIsV0FBTztBQUFBLEVBQ1g7QUFFQSxTQUFPLENBQUMsV0FBWSxXQUFXLE9BQU8sT0FBTyxRQUFRLE1BQU07QUFDL0Q7QUFNTyxTQUFTLE9BQU8sYUFFdkI7QUFDSSxNQUFJLFNBQVM7QUFDYixhQUFXLFFBQVEsYUFBYTtBQUM1QixRQUFJLFlBQVksSUFBSSxNQUFNLEtBQUs7QUFDM0IsZUFBUztBQUNUO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFDQSxNQUFJLFFBQVE7QUFDUixXQUFPO0FBQUEsRUFDWDtBQUVBLFNBQU8sQ0FBQyxXQUFXO0FBQ2YsZUFBVyxRQUFRLGFBQWE7QUFDNUIsVUFBSSxRQUFRLFFBQVE7QUFDaEIsZUFBTyxJQUFJLElBQUksWUFBWSxJQUFJLEVBQUUsT0FBTyxJQUFJLENBQUM7QUFBQSxNQUNqRDtBQUFBLElBQ0o7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUNKOzs7QUN6R0E7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBd0RBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksT0FBTztBQUVqRCxJQUFNLFNBQVM7QUFDZixJQUFNLGFBQWE7QUFDbkIsSUFBTSxhQUFhO0FBT1osU0FBUyxTQUE0QjtBQUN4QyxTQUFPQSxNQUFLLE1BQU07QUFDdEI7QUFPTyxTQUFTLGFBQThCO0FBQzFDLFNBQU9BLE1BQUssVUFBVTtBQUMxQjtBQU9PLFNBQVMsYUFBOEI7QUFDMUMsU0FBT0EsTUFBSyxVQUFVO0FBQzFCOzs7QXRCNUVBLE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQTRDbEMsT0FBTyxPQUFPLFNBQWdCO0FBQ3ZCLE9BQU8scUJBQXFCOyIsCiAgIm5hbWVzIjogWyJfYSIsICJFcnJvciIsICJjYWxsIiwgIl9hIiwgIkVycm9yIiwgImNhbGwiLCAiX2EiLCAicmVzaXphYmxlIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJjYWxsIiwgIkhpZGVNZXRob2QiLCAiU2hvd01ldGhvZCIsICJpc0RvY3VtZW50RG90QWxsIiwgIl9hIiwgInJlYXNvbiIsICJ2YWx1ZSIsICJjYWxsIiwgImdldEFuZERlbGV0ZVJlc3BvbnNlIiwgImdlbmVyYXRlSUQiLCAiY2FsbCIsICJBcnJheSIsICJNYXAiLCAiQXJyYXkiLCAiTWFwIiwgImtleSIsICJjYWxsIl0KfQo=
