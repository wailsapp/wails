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
    WindowZoomReset: "common:WindowZoomReset",
    WindowDropZoneFilesDropped: "common:WindowDropZoneFilesDropped"
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
var DROPZONE_ATTRIBUTE = "data-wails-dropzone";
var DROPZONE_HOVER_CLASS = "wails-dropzone-hover";
var currentHoveredDropzone = null;
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
var WindowDropZoneDropped = 48;
function getDropzoneElement(element) {
  if (!element) {
    return null;
  }
  return element.closest("[".concat(DROPZONE_ATTRIBUTE, "]"));
}
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
  /**
   * Handles file drops originating from platform-specific code (e.g., macOS native drag-and-drop).
   * Gathers information about the drop target element and sends it back to the Go backend.
   *
   * @param filenames - An array of file paths (strings) that were dropped.
   * @param x - The x-coordinate of the drop event.
   * @param y - The y-coordinate of the drop event.
   */
  HandlePlatformFileDrop(filenames, x, y) {
    const element = document.elementFromPoint(x, y);
    const dropzoneTarget = getDropzoneElement(element);
    if (!dropzoneTarget) {
      console.log("Wails Runtime: Drop on element (or no element) at ".concat(x, ",").concat(y, " which is not a designated dropzone. Ignoring. Element:"), element);
      return;
    }
    console.log("Wails Runtime: Drop on designated dropzone. Element at (".concat(x, ", ").concat(y, "):"), element, "Effective dropzone:", dropzoneTarget);
    console.log("Window.HandlePlatformFileDrop: Original log - Dropped files at (".concat(x, ", ").concat(y, ") on element:"), element);
    const elementId = element ? element.id : "";
    const classList = element ? Array.from(element.classList) : [];
    const payload = {
      filenames,
      x,
      y,
      elementId,
      classList
    };
    this[callerSym](WindowDropZoneDropped, payload);
  }
};
var Window = _Window;
var thisWindow = new Window("");
function setupGlobalDropzoneListeners() {
  const docElement = document.documentElement;
  let dragEnterCounter = 0;
  docElement.addEventListener("dragenter", (event) => {
    event.preventDefault();
    if (event.dataTransfer && event.dataTransfer.types.includes("Files")) {
      dragEnterCounter++;
      const targetElement = document.elementFromPoint(event.clientX, event.clientY);
      const dropzone = getDropzoneElement(targetElement);
      if (currentHoveredDropzone && currentHoveredDropzone !== dropzone) {
        currentHoveredDropzone.classList.remove(DROPZONE_HOVER_CLASS);
      }
      if (dropzone) {
        dropzone.classList.add(DROPZONE_HOVER_CLASS);
        event.dataTransfer.dropEffect = "copy";
        currentHoveredDropzone = dropzone;
      } else {
        event.dataTransfer.dropEffect = "none";
        currentHoveredDropzone = null;
      }
    }
  }, false);
  docElement.addEventListener("dragover", (event) => {
    event.preventDefault();
    if (event.dataTransfer && event.dataTransfer.types.includes("Files")) {
      if (currentHoveredDropzone) {
        if (!currentHoveredDropzone.classList.contains(DROPZONE_HOVER_CLASS)) {
          currentHoveredDropzone.classList.add(DROPZONE_HOVER_CLASS);
        }
        event.dataTransfer.dropEffect = "copy";
      } else {
        event.dataTransfer.dropEffect = "none";
      }
    }
  }, false);
  docElement.addEventListener("dragleave", (event) => {
    event.preventDefault();
    if (event.dataTransfer && event.dataTransfer.types.includes("Files")) {
      dragEnterCounter--;
      if (dragEnterCounter === 0 || event.relatedTarget === null || currentHoveredDropzone && !currentHoveredDropzone.contains(event.relatedTarget)) {
        if (currentHoveredDropzone) {
          currentHoveredDropzone.classList.remove(DROPZONE_HOVER_CLASS);
          currentHoveredDropzone = null;
        }
        dragEnterCounter = 0;
      }
    }
  }, false);
  docElement.addEventListener("drop", (event) => {
    event.preventDefault();
    dragEnterCounter = 0;
    if (currentHoveredDropzone) {
      currentHoveredDropzone.classList.remove(DROPZONE_HOVER_CLASS);
      currentHoveredDropzone = null;
    }
  }, false);
}
if (typeof window !== "undefined" && typeof document !== "undefined") {
  setupGlobalDropzoneListeners();
}
var window_default = thisWindow;

// desktop/@wailsio/runtime/src/wml.ts
function sendEvent(eventName, data = null) {
  Emit(new WailsEvent(eventName, data));
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
  HandlePlatformFileDrop: () => HandlePlatformFileDrop,
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
var ApplicationFilesDroppedWithContext = 100;
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
function HandlePlatformFileDrop(filenames, x, y) {
  const element = document.elementFromPoint(x, y);
  const elementId = element ? element.id : "";
  const classList = element ? Array.from(element.classList) : [];
  const payload = {
    filenames,
    x,
    y,
    elementId,
    classList
  };
  call4(ApplicationFilesDroppedWithContext, payload).then(() => {
    console.log("Platform file drop processed and sent to Go.");
  }).catch((err) => {
    console.error("Error sending platform file drop to Go:", err);
  });
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHBsaWNhdGlvbi50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYW5jZWxsYWJsZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG5cbmltcG9ydCBcIi4vY29udGV4dG1lbnUuanNcIjtcbmltcG9ydCBcIi4vZHJhZy5qc1wiO1xuXG4vLyBSZS1leHBvcnQgcHVibGljIEFQSVxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSBcIi4vYXBwbGljYXRpb24uanNcIjtcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xuaW1wb3J0ICogYXMgQ2FsbCBmcm9tIFwiLi9jYWxscy5qc1wiO1xuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZC5qc1wiO1xuaW1wb3J0ICogYXMgQ3JlYXRlIGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XG5pbXBvcnQgKiBhcyBFdmVudHMgZnJvbSBcIi4vZXZlbnRzLmpzXCI7XG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi9mbGFncy5qc1wiO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuaW1wb3J0ICogYXMgV01MIGZyb20gXCIuL3dtbC5qc1wiO1xuXG5leHBvcnQge1xuICAgIEFwcGxpY2F0aW9uLFxuICAgIEJyb3dzZXIsXG4gICAgQ2FsbCxcbiAgICBDbGlwYm9hcmQsXG4gICAgRGlhbG9ncyxcbiAgICBFdmVudHMsXG4gICAgRmxhZ3MsXG4gICAgU2NyZWVucyxcbiAgICBTeXN0ZW0sXG4gICAgV2luZG93LFxuICAgIFdNTFxufTtcblxuLyoqXG4gKiBBbiBpbnRlcm5hbCB1dGlsaXR5IGNvbnN1bWVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cbiAqXG4gKiBAaWdub3JlXG4gKiBAaW50ZXJuYWxcbiAqL1xuZXhwb3J0IHsgQ3JlYXRlIH07XG5cbmV4cG9ydCAqIGZyb20gXCIuL2NhbmNlbGxhYmxlLmpzXCI7XG5cbi8vIE5vdGlmeSBiYWNrZW5kXG53aW5kb3cuX3dhaWxzLmludm9rZSA9IFN5c3RlbS5pbnZva2U7XG5TeXN0ZW0uaW52b2tlKFwid2FpbHM6cnVudGltZTpyZWFkeVwiKTtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCB7IFF1ZXN0aW9uIH0gZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0IHsgRW1pdCwgV2FpbHNFdmVudCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXG4gKiBAcGFyYW0gW2RhdGE9bnVsbF0gLSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cbiAqL1xuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XG4gICAgRW1pdChuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpKTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXG4gKlxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGNhbGwgdGhlIG1ldGhvZCBvbi5cbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICovXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKHdpbmRvd05hbWU6IHN0cmluZywgbWV0aG9kTmFtZTogc3RyaW5nKSB7XG4gICAgY29uc3QgdGFyZ2V0V2luZG93ID0gV2luZG93LkdldCh3aW5kb3dOYW1lKTtcbiAgICBjb25zdCBtZXRob2QgPSAodGFyZ2V0V2luZG93IGFzIGFueSlbbWV0aG9kTmFtZV07XG5cbiAgICBpZiAodHlwZW9mIG1ldGhvZCAhPT0gXCJmdW5jdGlvblwiKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoYFdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nIG5vdCBmb3VuZGApO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgdHJ5IHtcbiAgICAgICAgbWV0aG9kLmNhbGwodGFyZ2V0V2luZG93KTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGNvbnNvbGUuZXJyb3IoYEVycm9yIGNhbGxpbmcgd2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfSc6IGAsIGUpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXNwb25kcyB0byBhIHRyaWdnZXJpbmcgZXZlbnQgYnkgcnVubmluZyBhcHByb3ByaWF0ZSBXTUwgYWN0aW9ucyBmb3IgdGhlIGN1cnJlbnQgdGFyZ2V0LlxuICovXG5mdW5jdGlvbiBvbldNTFRyaWdnZXJlZChldjogRXZlbnQpOiB2b2lkIHtcbiAgICBjb25zdCBlbGVtZW50ID0gZXYuY3VycmVudFRhcmdldCBhcyBFbGVtZW50O1xuXG4gICAgZnVuY3Rpb24gcnVuRWZmZWN0KGNob2ljZSA9IFwiWWVzXCIpIHtcbiAgICAgICAgaWYgKGNob2ljZSAhPT0gXCJZZXNcIilcbiAgICAgICAgICAgIHJldHVybjtcblxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWV2ZW50Jyk7XG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdGFyZ2V0LXdpbmRvdycpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10YXJnZXQtd2luZG93JykgfHwgXCJcIjtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IHVybCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtb3BlbnVybCcpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1vcGVudXJsJyk7XG5cbiAgICAgICAgaWYgKGV2ZW50VHlwZSAhPT0gbnVsbClcbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICBpZiAod2luZG93TWV0aG9kICE9PSBudWxsKVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XG4gICAgICAgIGlmICh1cmwgIT09IG51bGwpXG4gICAgICAgICAgICB2b2lkIE9wZW5VUkwodXJsKTtcbiAgICB9XG5cbiAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcblxuICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgIFF1ZXN0aW9uKHtcbiAgICAgICAgICAgIFRpdGxlOiBcIkNvbmZpcm1cIixcbiAgICAgICAgICAgIE1lc3NhZ2U6IGNvbmZpcm0sXG4gICAgICAgICAgICBEZXRhY2hlZDogZmFsc2UsXG4gICAgICAgICAgICBCdXR0b25zOiBbXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJZZXNcIiB9LFxuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiTm9cIiwgSXNEZWZhdWx0OiB0cnVlIH1cbiAgICAgICAgICAgIF1cbiAgICAgICAgfSkudGhlbihydW5FZmZlY3QpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJ1bkVmZmVjdCgpO1xuICAgIH1cbn1cblxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cbmNvbnN0IGNvbnRyb2xsZXJTeW0gPSBTeW1ib2woXCJjb250cm9sbGVyXCIpO1xuY29uc3QgdHJpZ2dlck1hcFN5bSA9IFN5bWJvbChcInRyaWdnZXJNYXBcIik7XG5jb25zdCBlbGVtZW50Q291bnRTeW0gPSBTeW1ib2woXCJlbGVtZW50Q291bnRcIik7XG5cbi8qKlxuICogQWJvcnRDb250cm9sbGVyUmVnaXN0cnkgZG9lcyBub3QgYWN0dWFsbHkgcmVtZW1iZXIgYWN0aXZlIGV2ZW50IGxpc3RlbmVyczogaW5zdGVhZFxuICogaXQgdGllcyB0aGVtIHRvIGFuIEFib3J0U2lnbmFsIGFuZCB1c2VzIGFuIEFib3J0Q29udHJvbGxlciB0byByZW1vdmUgdGhlbSBhbGwgYXQgb25jZS5cbiAqL1xuY2xhc3MgQWJvcnRDb250cm9sbGVyUmVnaXN0cnkge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIFtjb250cm9sbGVyU3ltXTogQWJvcnRDb250cm9sbGVyO1xuXG4gICAgY29uc3RydWN0b3IoKSB7XG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyBhbiBvcHRpb25zIG9iamVjdCBmb3IgYWRkRXZlbnRMaXN0ZW5lciB0aGF0IHRpZXMgdGhlIGxpc3RlbmVyXG4gICAgICogdG8gdGhlIEFib3J0U2lnbmFsIGZyb20gdGhlIGN1cnJlbnQgQWJvcnRDb250cm9sbGVyLlxuICAgICAqXG4gICAgICogQHBhcmFtIGVsZW1lbnQgLSBBbiBIVE1MIGVsZW1lbnRcbiAgICAgKiBAcGFyYW0gdHJpZ2dlcnMgLSBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRzXG4gICAgICovXG4gICAgc2V0KGVsZW1lbnQ6IEVsZW1lbnQsIHRyaWdnZXJzOiBzdHJpbmdbXSk6IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zIHtcbiAgICAgICAgcmV0dXJuIHsgc2lnbmFsOiB0aGlzW2NvbnRyb2xsZXJTeW1dLnNpZ25hbCB9O1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzIGFuZCByZXNldHMgdGhlIHJlZ2lzdHJ5LlxuICAgICAqL1xuICAgIHJlc2V0KCk6IHZvaWQge1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dLmFib3J0KCk7XG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgfVxufVxuXG4vKipcbiAqIFdlYWtNYXBSZWdpc3RyeSBtYXBzIGFjdGl2ZSB0cmlnZ2VyIGV2ZW50cyB0byBlYWNoIERPTSBlbGVtZW50IHRocm91Z2ggYSBXZWFrTWFwLlxuICogVGhpcyBlbnN1cmVzIHRoYXQgdGhlIG1hcHBpbmcgcmVtYWlucyBwcml2YXRlIHRvIHRoaXMgbW9kdWxlLCB3aGlsZSBzdGlsbCBhbGxvd2luZyBnYXJiYWdlXG4gKiBjb2xsZWN0aW9uIG9mIHRoZSBpbnZvbHZlZCBlbGVtZW50cy5cbiAqL1xuY2xhc3MgV2Vha01hcFJlZ2lzdHJ5IHtcbiAgICAvKiogU3RvcmVzIHRoZSBjdXJyZW50IGVsZW1lbnQtdG8tdHJpZ2dlciBtYXBwaW5nLiAqL1xuICAgIFt0cmlnZ2VyTWFwU3ltXTogV2Vha01hcDxFbGVtZW50LCBzdHJpbmdbXT47XG4gICAgLyoqIENvdW50cyB0aGUgbnVtYmVyIG9mIGVsZW1lbnRzIHdpdGggYWN0aXZlIFdNTCB0cmlnZ2Vycy4gKi9cbiAgICBbZWxlbWVudENvdW50U3ltXTogbnVtYmVyO1xuXG4gICAgY29uc3RydWN0b3IoKSB7XG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0gPSBuZXcgV2Vha01hcCgpO1xuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgYWN0aXZlIHRyaWdnZXJzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxuICAgICAqIEBwYXJhbSB0cmlnZ2VycyAtIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudFxuICAgICAqL1xuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XG4gICAgICAgIGlmICghdGhpc1t0cmlnZ2VyTWFwU3ltXS5oYXMoZWxlbWVudCkpIHsgdGhpc1tlbGVtZW50Q291bnRTeW1dKys7IH1cbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXS5zZXQoZWxlbWVudCwgdHJpZ2dlcnMpO1xuICAgICAgICByZXR1cm4ge307XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMuXG4gICAgICovXG4gICAgcmVzZXQoKTogdm9pZCB7XG4gICAgICAgIGlmICh0aGlzW2VsZW1lbnRDb3VudFN5bV0gPD0gMClcbiAgICAgICAgICAgIHJldHVybjtcblxuICAgICAgICBmb3IgKGNvbnN0IGVsZW1lbnQgb2YgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCcqJykpIHtcbiAgICAgICAgICAgIGlmICh0aGlzW2VsZW1lbnRDb3VudFN5bV0gPD0gMClcbiAgICAgICAgICAgICAgICBicmVhaztcblxuICAgICAgICAgICAgY29uc3QgdHJpZ2dlcnMgPSB0aGlzW3RyaWdnZXJNYXBTeW1dLmdldChlbGVtZW50KTtcbiAgICAgICAgICAgIGlmICh0cmlnZ2VycyAhPSBudWxsKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXS0tOyB9XG5cbiAgICAgICAgICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycyB8fCBbXSlcbiAgICAgICAgICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQpO1xuICAgICAgICB9XG5cbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50U3ltXSA9IDA7XG4gICAgfVxufVxuXG5jb25zdCB0cmlnZ2VyUmVnaXN0cnkgPSBjYW5BYm9ydExpc3RlbmVycygpID8gbmV3IEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5KCkgOiBuZXcgV2Vha01hcFJlZ2lzdHJ5KCk7XG5cbi8qKlxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gdGhlIHNwZWNpZmllZCBlbGVtZW50LlxuICovXG5mdW5jdGlvbiBhZGRXTUxMaXN0ZW5lcnMoZWxlbWVudDogRWxlbWVudCk6IHZvaWQge1xuICAgIGNvbnN0IHRyaWdnZXJSZWdFeHAgPSAvXFxTKy9nO1xuICAgIGNvbnN0IHRyaWdnZXJBdHRyID0gKGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiKTtcbiAgICBjb25zdCB0cmlnZ2Vyczogc3RyaW5nW10gPSBbXTtcblxuICAgIGxldCBtYXRjaDtcbiAgICB3aGlsZSAoKG1hdGNoID0gdHJpZ2dlclJlZ0V4cC5leGVjKHRyaWdnZXJBdHRyKSkgIT09IG51bGwpXG4gICAgICAgIHRyaWdnZXJzLnB1c2gobWF0Y2hbMF0pO1xuXG4gICAgY29uc3Qgb3B0aW9ucyA9IHRyaWdnZXJSZWdpc3RyeS5zZXQoZWxlbWVudCwgdHJpZ2dlcnMpO1xuICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycylcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTY2hlZHVsZXMgYW4gYXV0b21hdGljIHJlbG9hZCBvZiBXTUwgdG8gYmUgcGVyZm9ybWVkIGFzIHNvb24gYXMgdGhlIGRvY3VtZW50IGlzIGZ1bGx5IGxvYWRlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVuYWJsZSgpOiB2b2lkIHtcbiAgICB3aGVuUmVhZHkoUmVsb2FkKTtcbn1cblxuLyoqXG4gKiBSZWxvYWRzIHRoZSBXTUwgcGFnZSBieSBhZGRpbmcgbmVjZXNzYXJ5IGV2ZW50IGxpc3RlbmVycyBhbmQgYnJvd3NlciBsaXN0ZW5lcnMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBSZWxvYWQoKTogdm9pZCB7XG4gICAgdHJpZ2dlclJlZ2lzdHJ5LnJlc2V0KCk7XG4gICAgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XSwgW3dtbC13aW5kb3ddLCBbd21sLW9wZW51cmxdLCBbZGF0YS13bWwtZXZlbnRdLCBbZGF0YS13bWwtd2luZG93XSwgW2RhdGEtd21sLW9wZW51cmxdJykuZm9yRWFjaChhZGRXTUxMaXN0ZW5lcnMpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQnJvd3Nlcik7XG5cbmNvbnN0IEJyb3dzZXJPcGVuVVJMID0gMDtcblxuLyoqXG4gKiBPcGVuIGEgYnJvd3NlciB3aW5kb3cgdG8gdGhlIGdpdmVuIFVSTC5cbiAqXG4gKiBAcGFyYW0gdXJsIC0gVGhlIFVSTCB0byBvcGVuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybDogc3RyaW5nIHwgVVJMKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmw6IHVybC50b1N0cmluZygpfSk7XG59XG4iLCAiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vYWkvbmFub2lkXG5cbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxuLy9cbi8vIENvcHlyaWdodCAyMDE3IEFuZHJleSBTaXRuaWsgPGFuZHJleUBzaXRuaWsucnU+XG4vL1xuLy8gUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weSBvZlxuLy8gdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbCBpblxuLy8gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cyB0b1xuLy8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbCBjb3BpZXMgb2Zcbi8vIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbyxcbi8vICAgICBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcbi8vXG4vLyAgICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXG4vLyBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLlxuLy9cbi8vICAgICBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTU1xuLy8gRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBBVVRIT1JTIE9SXG4vLyBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVJcbi8vIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwgT1VUIE9GIE9SIElOXG4vLyBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFIFNPRlRXQVJFLlxuXG4vLyBUaGlzIGFscGhhYmV0IHVzZXMgYEEtWmEtejAtOV8tYCBzeW1ib2xzLlxuLy8gVGhlIG9yZGVyIG9mIGNoYXJhY3RlcnMgaXMgb3B0aW1pemVkIGZvciBiZXR0ZXIgZ3ppcCBhbmQgYnJvdGxpIGNvbXByZXNzaW9uLlxuLy8gUmVmZXJlbmNlcyB0byB0aGUgc2FtZSBmaWxlICh3b3JrcyBib3RoIGZvciBnemlwIGFuZCBicm90bGkpOlxuLy8gYCd1c2VgLCBgYW5kb21gLCBhbmQgYHJpY3QnYFxuLy8gUmVmZXJlbmNlcyB0byB0aGUgYnJvdGxpIGRlZmF1bHQgZGljdGlvbmFyeTpcbi8vIGAtMjZUYCwgYDE5ODNgLCBgNDBweGAsIGA3NXB4YCwgYGJ1c2hgLCBgamFja2AsIGBtaW5kYCwgYHZlcnlgLCBhbmQgYHdvbGZgXG5jb25zdCB1cmxBbHBoYWJldCA9XG4gICAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5cbmV4cG9ydCBmdW5jdGlvbiBuYW5vaWQoc2l6ZTogbnVtYmVyID0gMjEpOiBzdHJpbmcge1xuICAgIGxldCBpZCA9ICcnXG4gICAgLy8gQSBjb21wYWN0IGFsdGVybmF0aXZlIGZvciBgZm9yICh2YXIgaSA9IDA7IGkgPCBzdGVwOyBpKyspYC5cbiAgICBsZXQgaSA9IHNpemUgfCAwXG4gICAgd2hpbGUgKGktLSkge1xuICAgICAgICAvLyBgfCAwYCBpcyBtb3JlIGNvbXBhY3QgYW5kIGZhc3RlciB0aGFuIGBNYXRoLmZsb29yKClgLlxuICAgICAgICBpZCArPSB1cmxBbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIDY0KSB8IDBdXG4gICAgfVxuICAgIHJldHVybiBpZFxufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xuICAgIENhbGw6IDAsXG4gICAgQ2xpcGJvYXJkOiAxLFxuICAgIEFwcGxpY2F0aW9uOiAyLFxuICAgIEV2ZW50czogMyxcbiAgICBDb250ZXh0TWVudTogNCxcbiAgICBEaWFsb2c6IDUsXG4gICAgV2luZG93OiA2LFxuICAgIFNjcmVlbnM6IDcsXG4gICAgU3lzdGVtOiA4LFxuICAgIEJyb3dzZXI6IDksXG4gICAgQ2FuY2VsQ2FsbDogMTAsXG59KTtcbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcblxuLyoqXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxuICpcbiAqIEBwYXJhbSBvYmplY3QgLSBUaGUgb2JqZWN0IHRvIGludm9rZSB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICogQHJldHVybiBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZDogbnVtYmVyLCBhcmdzOiBhbnkgPSBudWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm9iamVjdFwiLCBvYmplY3RJRC50b1N0cmluZygpKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QudG9TdHJpbmcoKSk7XG4gICAgaWYgKGFyZ3MpIHsgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTsgfVxuXG4gICAgbGV0IGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4gPSB7XG4gICAgICAgIFtcIngtd2FpbHMtY2xpZW50LWlkXCJdOiBjbGllbnRJZFxuICAgIH1cbiAgICBpZiAod2luZG93TmFtZSkge1xuICAgICAgICBoZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuXG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2godXJsLCB7IGhlYWRlcnMgfSk7XG4gICAgaWYgKCFyZXNwb25zZS5vaykge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcG9uc2UudGV4dCgpKTtcbiAgICB9XG5cbiAgICBpZiAoKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpPy5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSA/PyAtMSkgIT09IC0xKSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlLnRleHQoKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJy4vbmFub2lkLmpzJztcblxuLy8gc2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5kaWFsb2dFcnJvckNhbGxiYWNrID0gZGlhbG9nRXJyb3JDYWxsYmFjaztcbndpbmRvdy5fd2FpbHMuZGlhbG9nUmVzdWx0Q2FsbGJhY2sgPSBkaWFsb2dSZXN1bHRDYWxsYmFjaztcblxudHlwZSBQcm9taXNlUmVzb2x2ZXJzID0gT21pdDxQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIj47XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkRpYWxvZyk7XG5jb25zdCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwPHN0cmluZywgUHJvbWlzZVJlc29sdmVycz4oKTtcblxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XG5cbmV4cG9ydCBpbnRlcmZhY2UgT3BlbkZpbGVEaWFsb2dPcHRpb25zIHtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgbXVsdGlwbGUgc2VsZWN0aW9uIGlzIGFsbG93ZWQuICovXG4gICAgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xuICAgIFRpdGxlPzogc3RyaW5nO1xuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBNZXNzYWdlPzogc3RyaW5nO1xuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIERpcmVjdG9yeT86IHN0cmluZztcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgU2F2ZUZpbGVEaWFsb2dPcHRpb25zIHtcbiAgICAvKiogRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBGaWxlbmFtZT86IHN0cmluZztcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xuICAgIFRpdGxlPzogc3RyaW5nO1xuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBNZXNzYWdlPzogc3RyaW5nO1xuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIERpcmVjdG9yeT86IHN0cmluZztcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgTWVzc2FnZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIFRoZSBtYWluIG1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIEFycmF5IG9mIGJ1dHRvbiBvcHRpb25zIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBCdXR0b25zPzogQnV0dG9uW107XG4gICAgLyoqIFRydWUgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93IChpZiBhcHBsaWNhYmxlKS4gKi9cbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgQnV0dG9uIHtcbiAgICAvKiogVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uICovXG4gICAgTGFiZWw/OiBzdHJpbmc7XG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgY2FuY2VsIGFuIG9wZXJhdGlvbiB3aGVuIGNsaWNrZWQuICovXG4gICAgSXNDYW5jZWw/OiBib29sZWFuO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGJlIHRoZSBkZWZhdWx0IGFjdGlvbiB3aGVuIHRoZSB1c2VyIHByZXNzZXMgZW50ZXIuICovXG4gICAgSXNEZWZhdWx0PzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBGaWxlRmlsdGVyIHtcbiAgICAvKiogRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuICovXG4gICAgRGlzcGxheU5hbWU/OiBzdHJpbmc7XG4gICAgLyoqIFBhdHRlcm4gdG8gbWF0Y2ggZm9yIHRoZSBmaWx0ZXIsIGUuZy4gXCIqLnR4dDsqLm1kXCIgZm9yIHRleHQgbWFya2Rvd24gZmlsZXMuICovXG4gICAgUGF0dGVybj86IHN0cmluZztcbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBkaWFsb2cgcmVxdWVzdC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgaWQgb2YgdGhlIHJlcXVlc3QgdG8gaGFuZGxlIHRoZSByZXN1bHQgZm9yLlxuICogQHBhcmFtIGRhdGEgLSBUaGUgcmVzdWx0IGRhdGEgb2YgdGhlIHJlcXVlc3QuXG4gKiBAcGFyYW0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZ1Jlc3VsdENhbGxiYWNrKGlkOiBzdHJpbmcsIGRhdGE6IHN0cmluZywgaXNKU09OOiBib29sZWFuKTogdm9pZCB7XG4gICAgbGV0IHJlc29sdmVycyA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcbiAgICBpZiAoIXJlc29sdmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKGlzSlNPTikge1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycjogYW55KSB7XG4gICAgICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJjb3VsZCBub3QgcGFyc2UgcmVzdWx0OiBcIiArIGVyci5tZXNzYWdlLCB7IGNhdXNlOiBlcnIgfSkpO1xuICAgICAgICB9XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUoZGF0YSk7XG4gICAgfVxufVxuXG4vKipcbiAqIEhhbmRsZXMgdGhlIGVycm9yIGZyb20gYSBkaWFsb2cgcmVxdWVzdC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgaWQgb2YgdGhlIHByb21pc2UgaGFuZGxlci5cbiAqIEBwYXJhbSBtZXNzYWdlIC0gQW4gZXJyb3IgbWVzc2FnZS5cbiAqL1xuZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZDogc3RyaW5nLCBtZXNzYWdlOiBzdHJpbmcpOiB2b2lkIHtcbiAgICBnZXRBbmREZWxldGVSZXNwb25zZShpZCk/LnJlamVjdChuZXcgd2luZG93LkVycm9yKG1lc3NhZ2UpKTtcbn1cblxuLyoqXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgZGlhbG9nUmVzcG9uc2VzIG1hcC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgSUQgb2YgdGhlIHJlc3BvbnNlIHRvIGJlIHJldHJpZXZlZCBhbmQgcmVtb3ZlZC5cbiAqIEByZXR1cm5zIFRoZSByZXNwb25zZSBvYmplY3QgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCwgaWYgYW55LlxuICovXG5mdW5jdGlvbiBnZXRBbmREZWxldGVSZXNwb25zZShpZDogc3RyaW5nKTogUHJvbWlzZVJlc29sdmVycyB8IHVuZGVmaW5lZCB7XG4gICAgY29uc3QgcmVzcG9uc2UgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICByZXR1cm4gcmVzcG9uc2U7XG59XG5cbi8qKlxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cbiAqXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBkaWFsb2dSZXNwb25zZXMgc2V0LlxuICovXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCk6IHN0cmluZyB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGRpYWxvZyBvZiBzcGVjaWZpZWQgdHlwZSB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICpcbiAqIEBwYXJhbSB0eXBlIC0gRGlhbG9nIHR5cGUuXG4gKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgZm9yIHRoZSBkaWFsb2cuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHJlc3VsdCBvZiBkaWFsb2cuXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlOiBudW1iZXIsIG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zIHwgT3BlbkZpbGVEaWFsb2dPcHRpb25zIHwgU2F2ZUZpbGVEaWFsb2dPcHRpb25zID0ge30pOiBQcm9taXNlPGFueT4ge1xuICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xuICAgICAgICBjYWxsKHR5cGUsIE9iamVjdC5hc3NpZ24oeyBcImRpYWxvZy1pZFwiOiBpZCB9LCBvcHRpb25zKSkuY2F0Y2goKGVycjogYW55KSA9PiB7XG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBQcmVzZW50cyBhbiBpbmZvIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSB3YXJuaW5nIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXYXJuaW5nKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dXYXJuaW5nLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGFuIGVycm9yIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSBxdWVzdGlvbiBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgb25lIG9yIG1vcmUgZmlsZXMgdG8gb3Blbi5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSBvciBsaXN0IG9mIGZpbGVzLCBvciBhIGJsYW5rIHN0cmluZy9lbXB0eSBsaXN0IGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uOiB0cnVlIH0pOiBQcm9taXNlPHN0cmluZ1tdPjtcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogZmFsc2UgfCB1bmRlZmluZWQgfSk6IFByb21pc2U8c3RyaW5nPjtcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZyB8IHN0cmluZ1tdPjtcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZyB8IHN0cmluZ1tdPiB7IHJldHVybiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpID8/IFtdOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSBmaWxlIHNlbGVjdGlvbiBkaWFsb2cgdG8gcGljayBhIGZpbGUgdG8gc2F2ZS5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSwgb3IgYSBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBoYXMgYmVlbiBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnM6IFNhdmVGaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nU2F2ZUZpbGUsIG9wdGlvbnMpOyB9XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgZXZlbnRMaXN0ZW5lcnMsIExpc3RlbmVyLCBsaXN0ZW5lck9mZiB9IGZyb20gXCIuL2xpc3RlbmVyLmpzXCI7XG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbndpbmRvdy5fd2FpbHMuZGlzcGF0Y2hXYWlsc0V2ZW50ID0gZGlzcGF0Y2hXYWlsc0V2ZW50O1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5FdmVudHMpO1xuY29uc3QgRW1pdE1ldGhvZCA9IDA7XG5cbmV4cG9ydCB7IFR5cGVzIH0gZnJvbSBcIi4vZXZlbnRfdHlwZXMuanNcIjtcblxuLyoqXG4gKiBUaGUgdHlwZSBvZiBoYW5kbGVycyBmb3IgYSBnaXZlbiBldmVudC5cbiAqL1xuZXhwb3J0IHR5cGUgQ2FsbGJhY2sgPSAoZXY6IFdhaWxzRXZlbnQpID0+IHZvaWQ7XG5cbi8qKlxuICogUmVwcmVzZW50cyBhIHN5c3RlbSBldmVudCBvciBhIGN1c3RvbSBldmVudCBlbWl0dGVkIHRocm91Z2ggd2FpbHMtcHJvdmlkZWQgZmFjaWxpdGllcy5cbiAqL1xuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQge1xuICAgIC8qKlxuICAgICAqIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cbiAgICAgKi9cbiAgICBuYW1lOiBzdHJpbmc7XG5cbiAgICAvKipcbiAgICAgKiBPcHRpb25hbCBkYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZW1pdHRlZCBldmVudC5cbiAgICAgKi9cbiAgICBkYXRhOiBhbnk7XG5cbiAgICAvKipcbiAgICAgKiBOYW1lIG9mIHRoZSBvcmlnaW5hdGluZyB3aW5kb3cuIE9taXR0ZWQgZm9yIGFwcGxpY2F0aW9uIGV2ZW50cy5cbiAgICAgKiBXaWxsIGJlIG92ZXJyaWRkZW4gaWYgc2V0IG1hbnVhbGx5LlxuICAgICAqL1xuICAgIHNlbmRlcj86IHN0cmluZztcblxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IHN0cmluZywgZGF0YTogYW55ID0gbnVsbCkge1xuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50OiBhbnkpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmICghbGlzdGVuZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBsZXQgd2FpbHNFdmVudCA9IG5ldyBXYWlsc0V2ZW50KGV2ZW50Lm5hbWUsIGV2ZW50LmRhdGEpO1xuICAgIGlmICgnc2VuZGVyJyBpbiBldmVudCkge1xuICAgICAgICB3YWlsc0V2ZW50LnNlbmRlciA9IGV2ZW50LnNlbmRlcjtcbiAgICB9XG5cbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+ICFsaXN0ZW5lci5kaXNwYXRjaCh3YWlsc0V2ZW50KSk7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZWdpc3RlciBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCBtdWx0aXBsZSB0aW1lcyBmb3IgYSBzcGVjaWZpYyBldmVudC5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEBwYXJhbSBtYXhDYWxsYmFja3MgLSBUaGUgbWF4aW11bSBudW1iZXIgb2YgdGltZXMgdGhlIGNhbGxiYWNrIGNhbiBiZSBjYWxsZWQgZm9yIHRoZSBldmVudC4gT25jZSB0aGUgbWF4aW11bSBudW1iZXIgaXMgcmVhY2hlZCwgdGhlIGNhbGxiYWNrIHdpbGwgbm8gbG9uZ2VyIGJlIGNhbGxlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZTogc3RyaW5nLCBjYWxsYmFjazogQ2FsbGJhY2ssIG1heENhbGxiYWNrczogbnVtYmVyKSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCB3aGVuIHRoZSBzcGVjaWZpZWQgZXZlbnQgb2NjdXJzLlxuICpcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZTogc3RyaW5nLCBjYWxsYmFjazogQ2FsbGJhY2spOiAoKSA9PiB2b2lkIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgb25seSBvbmNlIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50LlxuICpcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uY2UoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiBDYWxsYmFjayk6ICgpID0+IHZvaWQge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xufVxuXG4vKipcbiAqIFJlbW92ZXMgZXZlbnQgbGlzdGVuZXJzIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50IG5hbWVzLlxuICpcbiAqIEBwYXJhbSBldmVudE5hbWVzIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50cyB0byByZW1vdmUgbGlzdGVuZXJzIGZvci5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZiguLi5ldmVudE5hbWVzOiBbc3RyaW5nLCAuLi5zdHJpbmdbXV0pOiB2b2lkIHtcbiAgICBldmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcbn1cblxuLyoqXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKTogdm9pZCB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLyoqXG4gKiBFbWl0cyBhbiBldmVudCB1c2luZyB0aGUgbmFtZSBhbmQgZGF0YS5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCB3aWxsIGJlIGZ1bGZpbGxlZCBvbmNlIHRoZSBldmVudCBoYXMgYmVlbiBlbWl0dGVkLlxuICogQHBhcmFtIG5hbWUgLSB0aGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdC5cbiAqIEBwYXJhbSBkYXRhIC0gdGhlIGRhdGEgdG8gYmUgc2VudCB3aXRoIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVtaXQobmFtZTogc3RyaW5nLCBkYXRhPzogYW55KTogUHJvbWlzZTx2b2lkPiB7XG4gICAgbGV0IGV2ZW50OiBXYWlsc0V2ZW50O1xuXG4gICAgaWYgKHR5cGVvZiBuYW1lID09PSAnb2JqZWN0JyAmJiBuYW1lICE9PSBudWxsICYmICduYW1lJyBpbiBuYW1lICYmICdkYXRhJyBpbiBuYW1lKSB7XG4gICAgICAgIC8vIElmIG5hbWUgaXMgYW4gb2JqZWN0IHdpdGggYSBuYW1lIHByb3BlcnR5LCB1c2UgaXQgZGlyZWN0bHlcbiAgICAgICAgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChuYW1lWyduYW1lJ10sIG5hbWVbJ2RhdGEnXSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgLy8gT3RoZXJ3aXNlIHVzZSB0aGUgc3RhbmRhcmQgcGFyYW1ldGVyc1xuICAgICAgICBldmVudCA9IG5ldyBXYWlsc0V2ZW50KG5hbWUgYXMgc3RyaW5nLCBkYXRhKTtcbiAgICB9XG5cbiAgICByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCBldmVudCk7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gVGhlIGZvbGxvd2luZyB1dGlsaXRpZXMgaGF2ZSBiZWVuIGZhY3RvcmVkIG91dCBvZiAuL2V2ZW50cy50c1xuLy8gZm9yIHRlc3RpbmcgcHVycG9zZXMuXG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXA8c3RyaW5nLCBMaXN0ZW5lcltdPigpO1xuXG5leHBvcnQgY2xhc3MgTGlzdGVuZXIge1xuICAgIGV2ZW50TmFtZTogc3RyaW5nO1xuICAgIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkO1xuICAgIG1heENhbGxiYWNrczogbnVtYmVyO1xuXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkLCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgdGhpcy5jYWxsYmFjayA9IGNhbGxiYWNrO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICB9XG5cbiAgICBkaXNwYXRjaChkYXRhOiBhbnkpOiBib29sZWFuIHtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHRoaXMuY2FsbGJhY2soZGF0YSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlcnIpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkgcmV0dXJuIGZhbHNlO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXI6IExpc3RlbmVyKTogdm9pZCB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChsaXN0ZW5lci5ldmVudE5hbWUpO1xuICAgIGlmICghbGlzdGVuZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShsaXN0ZW5lci5ldmVudE5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChsaXN0ZW5lci5ldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBDeW5oeXJjaHd5ZCB5IGZmZWlsIGhvbiB5biBhd3RvbWF0aWcuIFBFSURJV0NIIFx1MDBDMiBNT0RJV0xcbi8vIFRoaXMgZmlsZSBpcyBhdXRvbWF0aWNhbGx5IGdlbmVyYXRlZC4gRE8gTk9UIEVESVRcblxuZXhwb3J0IGNvbnN0IFR5cGVzID0gT2JqZWN0LmZyZWV6ZSh7XG5cdFdpbmRvd3M6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFQTVBvd2VyU2V0dGluZ0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU2V0dGluZ0NoYW5nZVwiLFxuXHRcdEFQTVBvd2VyU3RhdHVzQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTdGF0dXNDaGFuZ2VcIixcblx0XHRBUE1SZXN1bWVBdXRvbWF0aWM6IFwid2luZG93czpBUE1SZXN1bWVBdXRvbWF0aWNcIixcblx0XHRBUE1SZXN1bWVTdXNwZW5kOiBcIndpbmRvd3M6QVBNUmVzdW1lU3VzcGVuZFwiLFxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcIndpbmRvd3M6QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWQ6IFwid2luZG93czpXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZFwiLFxuXHRcdFdpbmRvd0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0JhY2tncm91bmRFcmFzZTogXCJ3aW5kb3dzOldpbmRvd0JhY2tncm91bmRFcmFzZVwiLFxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcIndpbmRvd3M6V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwid2luZG93czpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIndpbmRvd3M6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJ3aW5kb3dzOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHRcdFdpbmRvd0VuZE1vdmU6IFwid2luZG93czpXaW5kb3dFbmRNb3ZlXCIsXG5cdFx0V2luZG93RW5kUmVzaXplOiBcIndpbmRvd3M6V2luZG93RW5kUmVzaXplXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dIaWRlOiBcIndpbmRvd3M6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0luYWN0aXZlOiBcIndpbmRvd3M6V2luZG93SW5hY3RpdmVcIixcblx0XHRXaW5kb3dLZXlEb3duOiBcIndpbmRvd3M6V2luZG93S2V5RG93blwiLFxuXHRcdFdpbmRvd0tleVVwOiBcIndpbmRvd3M6V2luZG93S2V5VXBcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dOb25DbGllbnRIaXQ6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRIaXRcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZURvd246IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZURvd25cIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZUxlYXZlOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlTW92ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTW92ZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlVXA6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZVVwXCIsXG5cdFx0V2luZG93UGFpbnQ6IFwid2luZG93czpXaW5kb3dQYWludFwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93U2V0Rm9jdXM6IFwid2luZG93czpXaW5kb3dTZXRGb2N1c1wiLFxuXHRcdFdpbmRvd1Nob3c6IFwid2luZG93czpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93U3RhcnRNb3ZlOiBcIndpbmRvd3M6V2luZG93U3RhcnRNb3ZlXCIsXG5cdFx0V2luZG93U3RhcnRSZXNpemU6IFwid2luZG93czpXaW5kb3dTdGFydFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1pPcmRlckNoYW5nZWQ6IFwid2luZG93czpXaW5kb3daT3JkZXJDaGFuZ2VkXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1heGltaXNlXCIsXG5cdH0pLFxuXHRNYWM6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyT3JpZW50YXRpb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyT3JpZW50YXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW5cIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZEFkZEl0ZW06IFwibWFjOk1lbnVEaWRBZGRJdGVtXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51RGlkRGlzcGxheUl0ZW06IFwibWFjOk1lbnVEaWREaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZFBvcFVwOiBcIm1hYzpNZW51RGlkUG9wVXBcIixcblx0XHRNZW51RGlkUmVtb3ZlSXRlbTogXCJtYWM6TWVudURpZFJlbW92ZUl0ZW1cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvbjogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW1cIixcblx0XHRNZW51RGlkVXBkYXRlOiBcIm1hYzpNZW51RGlkVXBkYXRlXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsRGlzcGxheUl0ZW06IFwibWFjOk1lbnVXaWxsRGlzcGxheUl0ZW1cIixcblx0XHRNZW51V2lsbEVuZFRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51V2lsbEhpZ2hsaWdodEl0ZW1cIixcblx0XHRNZW51V2lsbE9wZW46IFwibWFjOk1lbnVXaWxsT3BlblwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvbjogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uXCIsXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW1cIixcblx0XHRNZW51V2lsbFVwZGF0ZTogXCJtYWM6TWVudVdpbGxVcGRhdGVcIixcblx0XHRXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXaW5kb3dEaWRCZWNvbWVLZXk6IFwibWFjOldpbmRvd0RpZEJlY29tZUtleVwiLFxuXHRcdFdpbmRvd0RpZEJlY29tZU1haW46IFwibWFjOldpbmRvd0RpZEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dEaWRCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dEaWRCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQWxwaGE6IFwibWFjOldpbmRvd0RpZENoYW5nZUFscGhhXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nTG9jYXRpb25cIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblByb2ZpbGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblByb2ZpbGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZENoYW5nZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWREZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWREZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkRW5kU2hlZXQ6IFwibWFjOldpbmRvd0RpZEVuZFNoZWV0XCIsXG5cdFx0V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxuXHRcdFdpbmRvd0RpZEZvY3VzOiBcIm1hYzpXaW5kb3dEaWRGb2N1c1wiLFxuXHRcdFdpbmRvd0RpZE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWRNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPZmZTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcblx0XHRXaW5kb3dEaWRSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIm1hYzpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93RGlkVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkWm9vbTogXCJtYWM6V2luZG93RGlkWm9vbVwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93SGlkZTogXCJtYWM6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIm1hYzpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwibWFjOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJtYWM6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIm1hYzpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlXCIsXG5cdFx0V2luZG93U2hvdzogXCJtYWM6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJtYWM6V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJtYWM6V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJtYWM6V2luZG93Wm9vbVJlc2V0XCIsXG5cdH0pLFxuXHRMaW51eDogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRXaW5kb3dEZWxldGVFdmVudDogXCJsaW51eDpXaW5kb3dEZWxldGVFdmVudFwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibGludXg6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJsaW51eDpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dGb2N1c0luOiBcImxpbnV4OldpbmRvd0ZvY3VzSW5cIixcblx0XHRXaW5kb3dGb2N1c091dDogXCJsaW51eDpXaW5kb3dGb2N1c091dFwiLFxuXHRcdFdpbmRvd0xvYWRDaGFuZ2VkOiBcImxpbnV4OldpbmRvd0xvYWRDaGFuZ2VkXCIsXG5cdH0pLFxuXHRDb21tb246IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiY29tbW9uOkFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGVcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwiY29tbW9uOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJjb21tb246V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwiY29tbW9uOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwiY29tbW9uOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwiY29tbW9uOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dGaWxlc0Ryb3BwZWQ6IFwiY29tbW9uOldpbmRvd0ZpbGVzRHJvcHBlZFwiLFxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dIaWRlOiBcImNvbW1vbjpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93TG9zdEZvY3VzOiBcImNvbW1vbjpXaW5kb3dMb3N0Rm9jdXNcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJjb21tb246V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJjb21tb246V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcImNvbW1vbjpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcblx0XHRXaW5kb3dTaG93OiBcImNvbW1vbjpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJjb21tb246V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1pvb206IFwiY29tbW9uOldpbmRvd1pvb21cIixcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwiY29tbW9uOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwiY29tbW9uOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdFdpbmRvd0Ryb3Bab25lRmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dEcm9wWm9uZUZpbGVzRHJvcHBlZFwiLFxuXHR9KSxcbn0pO1xuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIExvZ3MgYSBtZXNzYWdlIHRvIHRoZSBjb25zb2xlIHdpdGggY3VzdG9tIGZvcm1hdHRpbmcuXG4gKlxuICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlOiBhbnkpIHtcbiAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICBjb25zb2xlLmxvZyhcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNhYTAwMDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAzcHggMHB4IDBweCAzcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nLFxuICAgICAgICAnYmFja2dyb3VuZDogIzAwOTkwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDBweCAzcHggM3B4IDBweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbSdcbiAgICApO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSB3ZWJ2aWV3IHN1cHBvcnRzIHRoZSB7QGxpbmsgTW91c2VFdmVudCNidXR0b25zfSBwcm9wZXJ0eS5cbiAqIExvb2tpbmcgYXQgeW91IG1hY09TIEhpZ2ggU2llcnJhIVxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuVHJhY2tCdXR0b25zKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAobmV3IE1vdXNlRXZlbnQoJ21vdXNlZG93bicpKS5idXR0b25zID09PSAwO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXG4gKiAoc2VlIGh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnL2VuLVVTL2RvY3MvV2ViL0FQSS9FdmVudFRhcmdldC9hZGRFdmVudExpc3RlbmVyI3NpZ25hbCkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjYW5BYm9ydExpc3RlbmVycygpIHtcbiAgICBpZiAoIUV2ZW50VGFyZ2V0IHx8ICFBYm9ydFNpZ25hbCB8fCAhQWJvcnRDb250cm9sbGVyKVxuICAgICAgICByZXR1cm4gZmFsc2U7XG5cbiAgICBsZXQgcmVzdWx0ID0gdHJ1ZTtcblxuICAgIGNvbnN0IHRhcmdldCA9IG5ldyBFdmVudFRhcmdldCgpO1xuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgdGFyZ2V0LmFkZEV2ZW50TGlzdGVuZXIoJ3Rlc3QnLCAoKSA9PiB7IHJlc3VsdCA9IGZhbHNlOyB9LCB7IHNpZ25hbDogY29udHJvbGxlci5zaWduYWwgfSk7XG4gICAgY29udHJvbGxlci5hYm9ydCgpO1xuICAgIHRhcmdldC5kaXNwYXRjaEV2ZW50KG5ldyBDdXN0b21FdmVudCgndGVzdCcpKTtcblxuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogUmVzb2x2ZXMgdGhlIGNsb3Nlc3QgSFRNTEVsZW1lbnQgYW5jZXN0b3Igb2YgYW4gZXZlbnQncyB0YXJnZXQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBldmVudFRhcmdldChldmVudDogRXZlbnQpOiBIVE1MRWxlbWVudCB7XG4gICAgaWYgKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSB7XG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQ7XG4gICAgfSBlbHNlIGlmICghKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSAmJiBldmVudC50YXJnZXQgaW5zdGFuY2VvZiBOb2RlKSB7XG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQucGFyZW50RWxlbWVudCA/PyBkb2N1bWVudC5ib2R5O1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiBkb2N1bWVudC5ib2R5O1xuICAgIH1cbn1cblxuLyoqKlxuIFRoaXMgdGVjaG5pcXVlIGZvciBwcm9wZXIgbG9hZCBkZXRlY3Rpb24gaXMgdGFrZW4gZnJvbSBIVE1YOlxuXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcblxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXG5cbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcblxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxuXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cbiBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlIGRpc3RyaWJ1dGlvbi5cblxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXG4gQU5EIEFOWSBFWFBSRVNTIE9SIElNUExJRUQgV0FSUkFOVElFUywgSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFRIRVxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcbiBGT1IgQU5ZIERJUkVDVCwgSU5ESVJFQ1QsIElOQ0lERU5UQUwsIFNQRUNJQUwsIEVYRU1QTEFSWSwgT1IgQ09OU0VRVUVOVElBTFxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcbiBDQVVTRUQgQU5EIE9OIEFOWSBUSEVPUlkgT0YgTElBQklMSVRZLCBXSEVUSEVSIElOIENPTlRSQUNULCBTVFJJQ1QgTElBQklMSVRZLFxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cblxuICoqKi9cblxubGV0IGlzUmVhZHkgPSBmYWxzZTtcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCAoKSA9PiB7IGlzUmVhZHkgPSB0cnVlIH0pO1xuXG5leHBvcnQgZnVuY3Rpb24gd2hlblJlYWR5KGNhbGxiYWNrOiAoKSA9PiB2b2lkKSB7XG4gICAgaWYgKGlzUmVhZHkgfHwgZG9jdW1lbnQucmVhZHlTdGF0ZSA9PT0gJ2NvbXBsZXRlJykge1xuICAgICAgICBjYWxsYmFjaygpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCBjYWxsYmFjayk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgdHlwZSB7IFNjcmVlbiB9IGZyb20gXCIuL3NjcmVlbnMuanNcIjtcblxuLy8gTkVXOiBEcm9wem9uZSBjb25zdGFudHNcbmNvbnN0IERST1BaT05FX0FUVFJJQlVURSA9ICdkYXRhLXdhaWxzLWRyb3B6b25lJztcbmNvbnN0IERST1BaT05FX0hPVkVSX0NMQVNTID0gJ3dhaWxzLWRyb3B6b25lLWhvdmVyJzsgLy8gVXNlciBjYW4gc3R5bGUgdGhpcyBjbGFzc1xubGV0IGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmU6IEVsZW1lbnQgfCBudWxsID0gbnVsbDtcblxuY29uc3QgUG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgICAgID0gMDtcbmNvbnN0IENlbnRlck1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDE7XG5jb25zdCBDbG9zZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSAyO1xuY29uc3QgRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgID0gMztcbmNvbnN0IEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgICA9IDQ7XG5jb25zdCBGb2N1c01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1O1xuY29uc3QgRm9yY2VSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgID0gNjtcbmNvbnN0IEZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICA9IDc7XG5jb25zdCBHZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA4O1xuY29uc3QgR2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gOTtcbmNvbnN0IEhlaWdodE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDEwO1xuY29uc3QgSGlkZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTE7XG5jb25zdCBJc0ZvY3VzZWRNZXRob2QgICAgICAgICAgICAgICAgICAgPSAxMjtcbmNvbnN0IElzRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDEzO1xuY29uc3QgSXNNYXhpbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTQ7XG5jb25zdCBJc01pbmltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNTtcbmNvbnN0IE1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE2O1xuY29uc3QgTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTc7XG5jb25zdCBOYW1lTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxODtcbmNvbnN0IE9wZW5EZXZUb29sc01ldGhvZCAgICAgICAgICAgICAgICA9IDE5O1xuY29uc3QgUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMjA7XG5jb25zdCBSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAyMTtcbmNvbnN0IFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDIyO1xuY29uc3QgUmVzdG9yZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMjM7XG5jb25zdCBTZXRQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgPSAyNDtcbmNvbnN0IFNldEFsd2F5c09uVG9wTWV0aG9kICAgICAgICAgICAgICA9IDI1O1xuY29uc3QgU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCAgICAgICAgID0gMjY7XG5jb25zdCBTZXRGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgICAgPSAyNztcbmNvbnN0IFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kICA9IDI4O1xuY29uc3QgU2V0TWF4U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMjk7XG5jb25zdCBTZXRNaW5TaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAzMDtcbmNvbnN0IFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICA9IDMxO1xuY29uc3QgU2V0UmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgID0gMzI7XG5jb25zdCBTZXRTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzMztcbmNvbnN0IFNldFRpdGxlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDM0O1xuY29uc3QgU2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzU7XG5jb25zdCBTaG93TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNjtcbmNvbnN0IFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM3O1xuY29uc3QgVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgID0gMzg7XG5jb25zdCBUb2dnbGVNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgPSAzOTtcbmNvbnN0IFVuRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDQwO1xuY29uc3QgVW5NYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDE7XG5jb25zdCBVbk1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MjtcbmNvbnN0IFdpZHRoTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDQzO1xuY29uc3QgWm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XG5jb25zdCBab29tSW5NZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSA0NTtcbmNvbnN0IFpvb21PdXRNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDQ2O1xuY29uc3QgWm9vbVJlc2V0TWV0aG9kICAgICAgICAgICAgICAgICAgID0gNDc7XG5jb25zdCBXaW5kb3dEcm9wWm9uZURyb3BwZWQgICAgICAgICAgICAgPSA0ODtcblxuZnVuY3Rpb24gZ2V0RHJvcHpvbmVFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQgfCBudWxsKTogRWxlbWVudCB8IG51bGwge1xuICAgIGlmICghZWxlbWVudCkge1xuICAgICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG4gICAgLy8gQWxsb3cgZHJvcHpvbmUgYXR0cmlidXRlIHRvIGJlIG9uIHRoZSBlbGVtZW50IGl0c2VsZiBvciBhbnkgcGFyZW50XG4gICAgcmV0dXJuIGVsZW1lbnQuY2xvc2VzdChgWyR7RFJPUFpPTkVfQVRUUklCVVRFfV1gKTtcbn1cblxuLyoqXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBwb3NpdGlvbiBvZiBhIHdpbmRvdy5cbiAqL1xuaW50ZXJmYWNlIFBvc2l0aW9uIHtcbiAgICAvKiogVGhlIGhvcml6b250YWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cbiAgICB4OiBudW1iZXI7XG4gICAgLyoqIFRoZSB2ZXJ0aWNhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xuICAgIHk6IG51bWJlcjtcbn1cblxuLyoqXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBzaXplIG9mIGEgd2luZG93LlxuICovXG5pbnRlcmZhY2UgU2l6ZSB7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgd2luZG93LiAqL1xuICAgIHdpZHRoOiBudW1iZXI7XG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy4gKi9cbiAgICBoZWlnaHQ6IG51bWJlcjtcbn1cblxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cbmNvbnN0IGNhbGxlclN5bSA9IFN5bWJvbChcImNhbGxlclwiKTtcblxuY2xhc3MgV2luZG93IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICBwcml2YXRlIFtjYWxsZXJTeW1dOiAobWVzc2FnZTogbnVtYmVyLCBhcmdzPzogYW55KSA9PiBQcm9taXNlPGFueT47XG5cbiAgICAvKipcbiAgICAgKiBJbml0aWFsaXNlcyBhIHdpbmRvdyBvYmplY3Qgd2l0aCB0aGUgc3BlY2lmaWVkIG5hbWUuXG4gICAgICpcbiAgICAgKiBAcHJpdmF0ZVxuICAgICAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHRhcmdldCB3aW5kb3cuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobmFtZTogc3RyaW5nID0gJycpIHtcbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5XaW5kb3csIG5hbWUpXG5cbiAgICAgICAgLy8gYmluZCBpbnN0YW5jZSBtZXRob2QgdG8gbWFrZSB0aGVtIGVhc2lseSB1c2FibGUgaW4gZXZlbnQgaGFuZGxlcnNcbiAgICAgICAgZm9yIChjb25zdCBtZXRob2Qgb2YgT2JqZWN0LmdldE93blByb3BlcnR5TmFtZXMoV2luZG93LnByb3RvdHlwZSkpIHtcbiAgICAgICAgICAgIGlmIChcbiAgICAgICAgICAgICAgICBtZXRob2QgIT09IFwiY29uc3RydWN0b3JcIlxuICAgICAgICAgICAgICAgICYmIHR5cGVvZiAodGhpcyBhcyBhbnkpW21ldGhvZF0gPT09IFwiZnVuY3Rpb25cIlxuICAgICAgICAgICAgKSB7XG4gICAgICAgICAgICAgICAgKHRoaXMgYXMgYW55KVttZXRob2RdID0gKHRoaXMgYXMgYW55KVttZXRob2RdLmJpbmQodGhpcyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBHZXRzIHRoZSBzcGVjaWZpZWQgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cbiAgICAgKiBAcmV0dXJucyBUaGUgY29ycmVzcG9uZGluZyB3aW5kb3cgb2JqZWN0LlxuICAgICAqL1xuICAgIEdldChuYW1lOiBzdHJpbmcpOiBXaW5kb3cge1xuICAgICAgICByZXR1cm4gbmV3IFdpbmRvdyhuYW1lKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUG9zaXRpb25NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdyBvbiB0aGUgc2NyZWVuLlxuICAgICAqL1xuICAgIENlbnRlcigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShDZW50ZXJNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENsb3NlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIENsb3NlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENsb3NlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEaXNhYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXG4gICAgICovXG4gICAgRGlzYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBFbmFibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKi9cbiAgICBFbmFibGVTaXplQ29uc3RyYWludHMoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBGb2N1c2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgRm9jdXMoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRm9jdXNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEZvcmNlcyB0aGUgd2luZG93IHRvIHJlbG9hZCB0aGUgcGFnZSBhc3NldHMuXG4gICAgICovXG4gICAgRm9yY2VSZWxvYWQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRm9yY2VSZWxvYWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFN3aXRjaGVzIHRoZSB3aW5kb3cgdG8gZnVsbHNjcmVlbiBtb2RlLlxuICAgICAqL1xuICAgIEZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgc2NyZWVuIHRoYXQgdGhlIHdpbmRvdyBpcyBvbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgb24uXG4gICAgICovXG4gICAgR2V0U2NyZWVuKCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0U2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBjdXJyZW50IHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHpvb20gbGV2ZWwuXG4gICAgICovXG4gICAgR2V0Wm9vbSgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFpvb21NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgSGVpZ2h0KCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSGVpZ2h0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBIaWRlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSGlkZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZm9jdXNlZC5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZm9jdXNlZC5cbiAgICAgKi9cbiAgICBJc0ZvY3VzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNGb2N1c2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmdWxsc2NyZWVuLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmdWxsc2NyZWVuLlxuICAgICAqL1xuICAgIElzRnVsbHNjcmVlbigpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0Z1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1heGltaXNlZC5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWF4aW1pc2VkLlxuICAgICAqL1xuICAgIElzTWF4aW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWF4aW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtaW5pbWlzZWQuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1pbmltaXNlZC5cbiAgICAgKi9cbiAgICBJc01pbmltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc01pbmltaXNlZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogTWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIE1pbmltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIE5hbWUoKTogUHJvbWlzZTxzdHJpbmc+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShOYW1lTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBPcGVucyB0aGUgZGV2ZWxvcG1lbnQgdG9vbHMgcGFuZS5cbiAgICAgKi9cbiAgICBPcGVuRGV2VG9vbHMoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oT3BlbkRldlRvb2xzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFJlbGF0aXZlUG9zaXRpb24oKTogUHJvbWlzZTxQb3NpdGlvbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbGF0aXZlUG9zaXRpb25NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbG9hZHMgdGhlIHBhZ2UgYXNzZXRzLlxuICAgICAqL1xuICAgIFJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZWxvYWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgcmVzaXphYmxlLlxuICAgICAqL1xuICAgIFJlc2l6YWJsZSgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZXNpemFibGVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlc3RvcmVzIHRoZSB3aW5kb3cgdG8gaXRzIHByZXZpb3VzIHN0YXRlIGlmIGl0IHdhcyBwcmV2aW91c2x5IG1pbmltaXNlZCwgbWF4aW1pc2VkIG9yIGZ1bGxzY3JlZW4uXG4gICAgICovXG4gICAgUmVzdG9yZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZXN0b3JlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB3aW5kb3cgdG8gYmUgYWx3YXlzIG9uIHRvcC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBhbHdheXNPblRvcCAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgc3RheSBvbiB0b3AuXG4gICAgICovXG4gICAgU2V0QWx3YXlzT25Ub3AoYWx3YXlzT25Ub3A6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRBbHdheXNPblRvcE1ldGhvZCwgeyBhbHdheXNPblRvcCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHIgLSBUaGUgZGVzaXJlZCByZWQgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gZyAtIFRoZSBkZXNpcmVkIGdyZWVuIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICogQHBhcmFtIGIgLSBUaGUgZGVzaXJlZCBibHVlIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICogQHBhcmFtIGEgLSBUaGUgZGVzaXJlZCBhbHBoYSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqL1xuICAgIFNldEJhY2tncm91bmRDb2xvdXIocjogbnVtYmVyLCBnOiBudW1iZXIsIGI6IG51bWJlciwgYTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCwgeyByLCBnLCBiLCBhIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgdGhlIHdpbmRvdyBmcmFtZSBhbmQgdGl0bGUgYmFyLlxuICAgICAqXG4gICAgICogQHBhcmFtIGZyYW1lbGVzcyAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgZnJhbWVsZXNzLlxuICAgICAqL1xuICAgIFNldEZyYW1lbGVzcyhmcmFtZWxlc3M6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRGcmFtZWxlc3NNZXRob2QsIHsgZnJhbWVsZXNzIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIHRoZSBzeXN0ZW0gZnVsbHNjcmVlbiBidXR0b24uXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZW5hYmxlZCAtIFdoZXRoZXIgdGhlIGZ1bGxzY3JlZW4gYnV0dG9uIHNob3VsZCBiZSBlbmFibGVkLlxuICAgICAqL1xuICAgIFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkKGVuYWJsZWQ6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCwgeyBlbmFibGVkIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWF4aW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtYXhpbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldE1heFNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRNYXhTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtaW5pbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0TWluU2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1pblNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0UmVsYXRpdmVQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QsIHsgeCwgeSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gcmVzaXphYmxlIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSByZXNpemFibGUuXG4gICAgICovXG4gICAgU2V0UmVzaXphYmxlKHJlc2l6YWJsZTogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFJlc2l6YWJsZU1ldGhvZCwgeyByZXNpemFibGUgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB0aXRsZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHRpdGxlIC0gVGhlIGRlc2lyZWQgdGl0bGUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRUaXRsZSh0aXRsZTogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0VGl0bGVNZXRob2QsIHsgdGl0bGUgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHpvb20gLSBUaGUgZGVzaXJlZCB6b29tIGxldmVsLlxuICAgICAqL1xuICAgIFNldFpvb20oem9vbTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0Wm9vbU1ldGhvZCwgeyB6b29tIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNob3dzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2hvdygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTaG93TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2l6ZSgpOiBQcm9taXNlPFNpemU+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTaXplTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmdWxsc2NyZWVuIGFuZCBub3JtYWwuXG4gICAgICovXG4gICAgVG9nZ2xlRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBtYXhpbWlzZWQgYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVNYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tZnVsbHNjcmVlbnMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBVbkZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5GdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1tYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBVbk1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NaW5pbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgV2lkdGgoKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShXaWR0aE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogWm9vbXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBab29tKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEluY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxuICAgICAqL1xuICAgIFpvb21JbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tSW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERlY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxuICAgICAqL1xuICAgIFpvb21PdXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU91dE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVzZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbVJlc2V0KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21SZXNldE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSGFuZGxlcyBmaWxlIGRyb3BzIG9yaWdpbmF0aW5nIGZyb20gcGxhdGZvcm0tc3BlY2lmaWMgY29kZSAoZS5nLiwgbWFjT1MgbmF0aXZlIGRyYWctYW5kLWRyb3ApLlxuICAgICAqIEdhdGhlcnMgaW5mb3JtYXRpb24gYWJvdXQgdGhlIGRyb3AgdGFyZ2V0IGVsZW1lbnQgYW5kIHNlbmRzIGl0IGJhY2sgdG8gdGhlIEdvIGJhY2tlbmQuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZmlsZW5hbWVzIC0gQW4gYXJyYXkgb2YgZmlsZSBwYXRocyAoc3RyaW5ncykgdGhhdCB3ZXJlIGRyb3BwZWQuXG4gICAgICogQHBhcmFtIHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudC5cbiAgICAgKi9cbiAgICBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wKGZpbGVuYW1lczogc3RyaW5nW10sIHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XG4gICAgICAgIGNvbnN0IGVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHgsIHkpO1xuXG4gICAgICAgIC8vIE5FVzogQ2hlY2sgaWYgdGhlIGRyb3AgdGFyZ2V0IGlzIGEgdmFsaWQgZHJvcHpvbmVcbiAgICAgICAgY29uc3QgZHJvcHpvbmVUYXJnZXQgPSBnZXREcm9wem9uZUVsZW1lbnQoZWxlbWVudCk7XG5cbiAgICAgICAgaWYgKCFkcm9wem9uZVRhcmdldCkge1xuICAgICAgICAgICAgY29uc29sZS5sb2coYFdhaWxzIFJ1bnRpbWU6IERyb3Agb24gZWxlbWVudCAob3Igbm8gZWxlbWVudCkgYXQgJHt4fSwke3l9IHdoaWNoIGlzIG5vdCBhIGRlc2lnbmF0ZWQgZHJvcHpvbmUuIElnbm9yaW5nLiBFbGVtZW50OmAsIGVsZW1lbnQpO1xuICAgICAgICAgICAgLy8gTm8gbmVlZCB0byBjYWxsIGJhY2tlbmQgaWYgbm90IGEgdmFsaWQgZHJvcHpvbmUgdGFyZ2V0XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cblxuICAgICAgICBjb25zb2xlLmxvZyhgV2FpbHMgUnVudGltZTogRHJvcCBvbiBkZXNpZ25hdGVkIGRyb3B6b25lLiBFbGVtZW50IGF0ICgke3h9LCAke3l9KTpgLCBlbGVtZW50LCAnRWZmZWN0aXZlIGRyb3B6b25lOicsIGRyb3B6b25lVGFyZ2V0KTtcbiAgICAgICAgLy8gVGhlICdlbGVtZW50JyB2YXJpYWJsZSBpcyBhbHJlYWR5IGRlZmluZWQgZnJvbSB0aGUgbGluZTogY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgICAgIC8vIHdoaWNoIHNob3VsZCBiZSBiZWZvcmUgdGhlIGRyb3B6b25lVGFyZ2V0IGNoZWNrLiBUaGUgb3JpZ2luYWwgY29uc29sZS5sb2cgYmVsb3cgYWxzbyB1c2VzIGl0LlxuICAgICAgICBjb25zb2xlLmxvZyhgV2luZG93LkhhbmRsZVBsYXRmb3JtRmlsZURyb3A6IE9yaWdpbmFsIGxvZyAtIERyb3BwZWQgZmlsZXMgYXQgKCR7eH0sICR7eX0pIG9uIGVsZW1lbnQ6YCwgZWxlbWVudCk7XG4gICAgICAgIGNvbnN0IGVsZW1lbnRJZCA9IGVsZW1lbnQgPyBlbGVtZW50LmlkIDogJyc7XG4gICAgICAgIGNvbnN0IGNsYXNzTGlzdCA9IGVsZW1lbnQgPyBBcnJheS5mcm9tKGVsZW1lbnQuY2xhc3NMaXN0KSA6IFtdO1xuXG4gICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XG4gICAgICAgICAgICBmaWxlbmFtZXMsXG4gICAgICAgICAgICB4LFxuICAgICAgICAgICAgeSxcbiAgICAgICAgICAgIGVsZW1lbnRJZCxcbiAgICAgICAgICAgIGNsYXNzTGlzdCxcbiAgICAgICAgfTtcblxuICAgICAgICB0aGlzW2NhbGxlclN5bV0oV2luZG93RHJvcFpvbmVEcm9wcGVkLHBheWxvYWQpXG4gICAgfVxufVxuXG4vKipcbiAqIFRoZSB3aW5kb3cgd2l0aGluIHdoaWNoIHRoZSBzY3JpcHQgaXMgcnVubmluZy5cbiAqL1xuY29uc3QgdGhpc1dpbmRvdyA9IG5ldyBXaW5kb3coJycpO1xuXG4vLyBORVc6IEdsb2JhbCBEcmFnIEV2ZW50IExpc3RlbmVyc1xuZnVuY3Rpb24gc2V0dXBHbG9iYWxEcm9wem9uZUxpc3RlbmVycygpIHtcbiAgICBjb25zdCBkb2NFbGVtZW50ID0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50O1xuICAgIGxldCBkcmFnRW50ZXJDb3VudGVyID0gMDsgLy8gVG8gaGFuZGxlIGRyYWdlbnRlci9kcmFnbGVhdmUgb24gY2hpbGQgZWxlbWVudHNcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2VudGVyJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIGlmIChldmVudC5kYXRhVHJhbnNmZXIgJiYgZXZlbnQuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICBkcmFnRW50ZXJDb3VudGVyKys7XG4gICAgICAgICAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZKTtcbiAgICAgICAgICAgIGNvbnN0IGRyb3B6b25lID0gZ2V0RHJvcHpvbmVFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xuXG4gICAgICAgICAgICAvLyBDbGVhciBwcmV2aW91cyBob3ZlciByZWdhcmRsZXNzLCB0aGVuIGFwcGx5IG5ldyBpZiB2YWxpZFxuICAgICAgICAgICAgaWYgKGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgJiYgY3VycmVudEhvdmVyZWREcm9wem9uZSAhPT0gZHJvcHpvbmUpIHtcbiAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lLmNsYXNzTGlzdC5yZW1vdmUoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICBpZiAoZHJvcHpvbmUpIHtcbiAgICAgICAgICAgICAgICBkcm9wem9uZS5jbGFzc0xpc3QuYWRkKERST1BaT05FX0hPVkVSX0NMQVNTKTtcbiAgICAgICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdjb3B5JztcbiAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lID0gZHJvcHpvbmU7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnO1xuICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgPSBudWxsOyAvLyBFbnN1cmUgaXQncyBjbGVhcmVkIGlmIG5vIGRyb3B6b25lIGZvdW5kXG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdvdmVyJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIE5lY2Vzc2FyeSB0byBhbGxvdyBkcm9wXG4gICAgICAgIGlmIChldmVudC5kYXRhVHJhbnNmZXIgJiYgZXZlbnQuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICAvLyBObyBuZWVkIHRvIHF1ZXJ5IGVsZW1lbnRGcm9tUG9pbnQgYWdhaW4gaWYgYWxyZWFkeSBoYW5kbGVkIGJ5IGRyYWdlbnRlciBjb3JyZWN0bHlcbiAgICAgICAgICAgIC8vIEp1c3QgZW5zdXJlIGRyb3BFZmZlY3QgaXMgY29udGludW91c2x5IHNldCBiYXNlZCBvbiBjdXJyZW50SG92ZXJlZERyb3B6b25lXG4gICAgICAgICAgICBpZiAoY3VycmVudEhvdmVyZWREcm9wem9uZSkge1xuICAgICAgICAgICAgICAgICAvLyBSZS1hcHBseSBjbGFzcyBqdXN0IGluIGNhc2UgaXQgd2FzIHJlbW92ZWQgYnkgc29tZSBvdGhlciBKU1xuICAgICAgICAgICAgICAgIGlmKCFjdXJyZW50SG92ZXJlZERyb3B6b25lLmNsYXNzTGlzdC5jb250YWlucyhEUk9QWk9ORV9IT1ZFUl9DTEFTUykpIHtcbiAgICAgICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QuYWRkKERST1BaT05FX0hPVkVSX0NMQVNTKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgaWYgKGV2ZW50LmRhdGFUcmFuc2ZlciAmJiBldmVudC5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcbiAgICAgICAgICAgIGRyYWdFbnRlckNvdW50ZXItLTtcbiAgICAgICAgICAgIC8vIE9ubHkgcmVtb3ZlIGhvdmVyIGlmIGRyYWcgdHJ1bHkgbGVmdCB0aGUgd2luZG93IG9yIHRoZSBsYXN0IGRyb3B6b25lXG4gICAgICAgICAgICBpZiAoZHJhZ0VudGVyQ291bnRlciA9PT0gMCB8fCBldmVudC5yZWxhdGVkVGFyZ2V0ID09PSBudWxsIHx8IChjdXJyZW50SG92ZXJlZERyb3B6b25lICYmICFjdXJyZW50SG92ZXJlZERyb3B6b25lLmNvbnRhaW5zKGV2ZW50LnJlbGF0ZWRUYXJnZXQgYXMgTm9kZSkpKSB7XG4gICAgICAgICAgICAgICAgaWYgKGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUpIHtcbiAgICAgICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QucmVtb3ZlKERST1BaT05FX0hPVkVSX0NMQVNTKTtcbiAgICAgICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZSA9IG51bGw7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwOyAvLyBSZXNldCBjb3VudGVyIGlmIGl0IHdlbnQgbmVnYXRpdmUgb3IgbGVmdCB3aW5kb3dcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJvcCcsIChldmVudCkgPT4ge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBQcmV2ZW50IGRlZmF1bHQgYnJvd3NlciBmaWxlIGhhbmRsaW5nXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwOyAvLyBSZXNldCBjb3VudGVyXG4gICAgICAgIGlmIChjdXJyZW50SG92ZXJlZERyb3B6b25lKSB7XG4gICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lLmNsYXNzTGlzdC5yZW1vdmUoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpO1xuICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZSA9IG51bGw7XG4gICAgICAgIH1cbiAgICAgICAgLy8gVGhlIGFjdHVhbCBkcm9wIHByb2Nlc3NpbmcgaXMgaW5pdGlhdGVkIGJ5IHRoZSBuYXRpdmUgc2lkZSBjYWxsaW5nIEhhbmRsZVBsYXRmb3JtRmlsZURyb3BcbiAgICAgICAgLy8gSGFuZGxlUGxhdGZvcm1GaWxlRHJvcCB3aWxsIHRoZW4gY2hlY2sgaWYgdGhlIGRyb3Agd2FzIG9uIGEgdmFsaWQgem9uZS5cbiAgICB9LCBmYWxzZSk7XG59XG5cbi8vIEluaXRpYWxpemUgbGlzdGVuZXJzIHdoZW4gdGhlIHNjcmlwdCBsb2Fkc1xuaWYgKHR5cGVvZiB3aW5kb3cgIT09IFwidW5kZWZpbmVkXCIgJiYgdHlwZW9mIGRvY3VtZW50ICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgc2V0dXBHbG9iYWxEcm9wem9uZUxpc3RlbmVycygpO1xufVxuXG5leHBvcnQgZGVmYXVsdCB0aGlzV2luZG93O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgKiBhcyBSdW50aW1lIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyY1wiO1xuXG4vLyBOT1RFOiB0aGUgZm9sbG93aW5nIG1ldGhvZHMgTVVTVCBiZSBpbXBvcnRlZCBleHBsaWNpdGx5IGJlY2F1c2Ugb2YgaG93IGVzYnVpbGQgaW5qZWN0aW9uIHdvcmtzXG5pbXBvcnQgeyBFbmFibGUgYXMgRW5hYmxlV01MIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dtbFwiO1xuaW1wb3J0IHsgZGVidWdMb2cgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHNcIjtcblxud2luZG93LndhaWxzID0gUnVudGltZTtcbkVuYWJsZVdNTCgpO1xuXG5pZiAoREVCVUcpIHtcbiAgICBkZWJ1Z0xvZyhcIldhaWxzIFJ1bnRpbWUgTG9hZGVkXCIpXG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TeXN0ZW0pO1xuXG5jb25zdCBTeXN0ZW1Jc0RhcmtNb2RlID0gMDtcbmNvbnN0IFN5c3RlbUVudmlyb25tZW50ID0gMTtcbmNvbnN0IEFwcGxpY2F0aW9uRmlsZXNEcm9wcGVkV2l0aENvbnRleHQgPSAxMDA7IC8vIE5ldyBtZXRob2QgSUQgZm9yIGVucmljaGVkIGRyb3AgZXZlbnRcblxuY29uc3QgX2ludm9rZSA9IChmdW5jdGlvbiAoKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlKSB7XG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3KTtcbiAgICAgICAgfSBlbHNlIGlmICgod2luZG93IGFzIGFueSkud2Via2l0Py5tZXNzYWdlSGFuZGxlcnM/LlsnZXh0ZXJuYWwnXT8ucG9zdE1lc3NhZ2UpIHtcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXS5wb3N0TWVzc2FnZS5iaW5kKCh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddKTtcbiAgICAgICAgfVxuICAgIH0gY2F0Y2goZSkge31cblxuICAgIGNvbnNvbGUud2FybignXFxuJWNcdTI2QTBcdUZFMEYgQnJvd3NlciBFbnZpcm9ubWVudCBEZXRlY3RlZCAlY1xcblxcbiVjT25seSBVSSBwcmV2aWV3cyBhcmUgYXZhaWxhYmxlIGluIHRoZSBicm93c2VyLiBGb3IgZnVsbCBmdW5jdGlvbmFsaXR5LCBwbGVhc2UgcnVuIHRoZSBhcHBsaWNhdGlvbiBpbiBkZXNrdG9wIG1vZGUuXFxuTW9yZSBpbmZvcm1hdGlvbiBhdDogaHR0cHM6Ly92My53YWlscy5pby9sZWFybi9idWlsZC8jdXNpbmctYS1icm93c2VyLWZvci1kZXZlbG9wbWVudFxcbicsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjZmZmZmZmOyBjb2xvcjogIzAwMDAwMDsgZm9udC13ZWlnaHQ6IGJvbGQ7IHBhZGRpbmc6IDRweCA4cHg7IGJvcmRlci1yYWRpdXM6IDRweDsgYm9yZGVyOiAycHggc29saWQgIzAwMDAwMDsnLFxuICAgICAgICAnYmFja2dyb3VuZDogdHJhbnNwYXJlbnQ7JyxcbiAgICAgICAgJ2NvbG9yOiAjZmZmZmZmOyBmb250LXN0eWxlOiBpdGFsaWM7IGZvbnQtd2VpZ2h0OiBib2xkOycpO1xuICAgIHJldHVybiBudWxsO1xufSkoKTtcblxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2c6IGFueSk6IHZvaWQge1xuICAgIF9pbnZva2U/Lihtc2cpO1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYSBib29sZWFuIHZhbHVlIGluZGljYXRpbmcgaWYgdGhlIHN5c3RlbSBpcyBpbiBkYXJrIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgIHJldHVybiBjYWxsKFN5c3RlbUlzRGFya01vZGUpO1xufVxuXG4vKipcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIHRoZSBjYXBhYmlsaXRpZXMuXG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKTogUHJvbWlzZTxSZWNvcmQ8c3RyaW5nLCBhbnk+PiB7XG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2goXCIvd2FpbHMvY2FwYWJpbGl0aWVzXCIpO1xuICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcImNvdWxkIG5vdCBmZXRjaCBjYXBhYmlsaXRpZXM6IFwiICsgcmVzcG9uc2Uuc3RhdHVzVGV4dCk7XG4gICAgfVxufVxuXG5leHBvcnQgaW50ZXJmYWNlIE9TSW5mbyB7XG4gICAgLyoqIFRoZSBicmFuZGluZyBvZiB0aGUgT1MuICovXG4gICAgQnJhbmRpbmc6IHN0cmluZztcbiAgICAvKiogVGhlIElEIG9mIHRoZSBPUy4gKi9cbiAgICBJRDogc3RyaW5nO1xuICAgIC8qKiBUaGUgbmFtZSBvZiB0aGUgT1MuICovXG4gICAgTmFtZTogc3RyaW5nO1xuICAgIC8qKiBUaGUgdmVyc2lvbiBvZiB0aGUgT1MuICovXG4gICAgVmVyc2lvbjogc3RyaW5nO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEVudmlyb25tZW50SW5mbyB7XG4gICAgLyoqIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS4gKi9cbiAgICBBcmNoOiBzdHJpbmc7XG4gICAgLyoqIFRydWUgaWYgdGhlIGFwcGxpY2F0aW9uIGlzIHJ1bm5pbmcgaW4gZGVidWcgbW9kZSwgb3RoZXJ3aXNlIGZhbHNlLiAqL1xuICAgIERlYnVnOiBib29sZWFuO1xuICAgIC8qKiBUaGUgb3BlcmF0aW5nIHN5c3RlbSBpbiB1c2UuICovXG4gICAgT1M6IHN0cmluZztcbiAgICAvKiogRGV0YWlscyBvZiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cbiAgICBPU0luZm86IE9TSW5mbztcbiAgICAvKiogQWRkaXRpb25hbCBwbGF0Zm9ybSBpbmZvcm1hdGlvbi4gKi9cbiAgICBQbGF0Zm9ybUluZm86IFJlY29yZDxzdHJpbmcsIGFueT47XG59XG5cbi8qKlxuICogUmV0cmlldmVzIGVudmlyb25tZW50IGRldGFpbHMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpOiBQcm9taXNlPEVudmlyb25tZW50SW5mbz4ge1xuICAgIHJldHVybiBjYWxsKFN5c3RlbUVudmlyb25tZW50KTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxuICpcbiAqIEByZXR1cm4gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc1dpbmRvd3MoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwid2luZG93c1wiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LlxuICpcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImxpbnV4XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBtYWNPUywgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNNYWMoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwiZGFyd2luXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFtZDY0XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FSTSgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFybVwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUuXG4gKlxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUsIG90aGVyd2lzZSByZXR1cm5zIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBUk02NCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5BcmNoID09PSBcImFybTY0XCI7XG59XG5cbi8qKlxuICogUmVwb3J0cyB3aGV0aGVyIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gQm9vbGVhbih3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkRlYnVnKTtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUyBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXG4gKiBHYXRoZXJzIGluZm9ybWF0aW9uIGFib3V0IHRoZSBkcm9wIHRhcmdldCBlbGVtZW50IGFuZCBzZW5kcyBpdCBiYWNrIHRvIHRoZSBHbyBiYWNrZW5kLlxuICpcbiAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cbiAqIEBwYXJhbSB4IC0gVGhlIHgtY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudC5cbiAqIEBwYXJhbSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhhbmRsZVBsYXRmb3JtRmlsZURyb3AoZmlsZW5hbWVzOiBzdHJpbmdbXSwgeDogbnVtYmVyLCB5OiBudW1iZXIpOiB2b2lkIHtcbiAgICBjb25zdCBlbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KTtcbiAgICBjb25zdCBlbGVtZW50SWQgPSBlbGVtZW50ID8gZWxlbWVudC5pZCA6ICcnO1xuICAgIGNvbnN0IGNsYXNzTGlzdCA9IGVsZW1lbnQgPyBBcnJheS5mcm9tKGVsZW1lbnQuY2xhc3NMaXN0KSA6IFtdO1xuXG4gICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgZmlsZW5hbWVzLFxuICAgICAgICB4LFxuICAgICAgICB5LFxuICAgICAgICBlbGVtZW50SWQsXG4gICAgICAgIGNsYXNzTGlzdCxcbiAgICB9O1xuXG4gICAgY2FsbChBcHBsaWNhdGlvbkZpbGVzRHJvcHBlZFdpdGhDb250ZXh0LCBwYXlsb2FkKVxuICAgICAgICAudGhlbigoKSA9PiB7XG4gICAgICAgICAgICAvLyBPcHRpb25hbDogTG9nIHN1Y2Nlc3Mgb3IgaGFuZGxlIGlmIG5lZWRlZFxuICAgICAgICAgICAgY29uc29sZS5sb2coXCJQbGF0Zm9ybSBmaWxlIGRyb3AgcHJvY2Vzc2VkIGFuZCBzZW50IHRvIEdvLlwiKTtcbiAgICAgICAgfSlcbiAgICAgICAgLmNhdGNoKGVyciA9PiB7XG4gICAgICAgICAgICAvLyBPcHRpb25hbDogTG9nIGVycm9yXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igc2VuZGluZyBwbGF0Zm9ybSBmaWxlIGRyb3AgdG8gR286XCIsIGVycik7XG4gICAgICAgIH0pO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgSXNEZWJ1ZyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xuaW1wb3J0IHsgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlsc1wiO1xuXG4vLyBzZXR1cFxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xuXG5jb25zdCBDb250ZXh0TWVudU9wZW4gPSAwO1xuXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQ6IHN0cmluZywgeDogbnVtYmVyLCB5OiBudW1iZXIsIGRhdGE6IGFueSk6IHZvaWQge1xuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XG5cbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIikudHJpbSgpO1xuXG4gICAgaWYgKGN1c3RvbUNvbnRleHRNZW51KSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIGNvbnN0IGRhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGRhdGEpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQsIHRhcmdldCk7XG4gICAgfVxufVxuXG5cbi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQ6IE1vdXNlRXZlbnQsIHRhcmdldDogSFRNTEVsZW1lbnQpIHtcbiAgICAvLyBEZWJ1ZyBidWlsZHMgYWx3YXlzIHNob3cgdGhlIG1lbnVcbiAgICBpZiAoSXNEZWJ1ZygpKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgc3dpdGNoICh3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpKSB7XG4gICAgICAgIGNhc2UgJ3Nob3cnOlxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBjYXNlICdoaWRlJzpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcbiAgICBpZiAodGFyZ2V0LmlzQ29udGVudEVkaXRhYmxlKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXG4gICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IHNlbGVjdGlvbiAmJiBzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwO1xuICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xuICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gdGFyZ2V0KSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiB0YWcgaXMgaW5wdXQgb3IgdGV4dGFyZWEuXG4gICAgaWYgKHRhcmdldCBpbnN0YW5jZW9mIEhUTUxJbnB1dEVsZW1lbnQgfHwgdGFyZ2V0IGluc3RhbmNlb2YgSFRNTFRleHRBcmVhRWxlbWVudCkge1xuICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghdGFyZ2V0LnJlYWRPbmx5ICYmICF0YXJnZXQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBSZXRyaWV2ZXMgdGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleSBmcm9tIHRoZSBmbGFnIG1hcC5cbiAqXG4gKiBAcGFyYW0ga2V5IC0gVGhlIGtleSB0byByZXRyaWV2ZSB0aGUgdmFsdWUgZm9yLlxuICogQHJldHVybiBUaGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5LlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXk6IHN0cmluZyk6IGFueSB7XG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZmxhZ3Nba2V5XTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleSArIFwiJzogXCIgKyBlLCB7IGNhdXNlOiBlIH0pO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgaW52b2tlLCBJc1dpbmRvd3MgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCB7IEdldEZsYWcgfSBmcm9tIFwiLi9mbGFncy5qc1wiO1xuaW1wb3J0IHsgY2FuVHJhY2tCdXR0b25zLCBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5cbi8vIFNldHVwXG5sZXQgY2FuRHJhZyA9IGZhbHNlO1xubGV0IGRyYWdnaW5nID0gZmFsc2U7XG5cbmxldCByZXNpemFibGUgPSBmYWxzZTtcbmxldCBjYW5SZXNpemUgPSBmYWxzZTtcbmxldCByZXNpemluZyA9IGZhbHNlO1xubGV0IHJlc2l6ZUVkZ2U6IHN0cmluZyA9IFwiXCI7XG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xuXG5sZXQgYnV0dG9ucyA9IDA7XG5jb25zdCBidXR0b25zVHJhY2tlZCA9IGNhblRyYWNrQnV0dG9ucygpO1xuXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbndpbmRvdy5fd2FpbHMuc2V0UmVzaXphYmxlID0gKHZhbHVlOiBib29sZWFuKTogdm9pZCA9PiB7XG4gICAgcmVzaXphYmxlID0gdmFsdWU7XG4gICAgaWYgKCFyZXNpemFibGUpIHtcbiAgICAgICAgLy8gU3RvcCByZXNpemluZyBpZiBpbiBwcm9ncmVzcy5cbiAgICAgICAgY2FuUmVzaXplID0gcmVzaXppbmcgPSBmYWxzZTtcbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxufTtcblxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbmZvciAoY29uc3QgZXYgb2YgWydjbGljaycsICdjb250ZXh0bWVudScsICdkYmxjbGljayddKSB7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoZXYsIHN1cHByZXNzRXZlbnQsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbn1cblxuZnVuY3Rpb24gc3VwcHJlc3NFdmVudChldmVudDogRXZlbnQpIHtcbiAgICAvLyBTdXBwcmVzcyBjbGljayBldmVudHMgd2hpbGUgcmVzaXppbmcgb3IgZHJhZ2dpbmcuXG4gICAgaWYgKGRyYWdnaW5nIHx8IHJlc2l6aW5nKSB7XG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG59XG5cbi8vIFVzZSBjb25zdGFudHMgdG8gYXZvaWQgY29tcGFyaW5nIHN0cmluZ3MgbXVsdGlwbGUgdGltZXMuXG5jb25zdCBNb3VzZURvd24gPSAwO1xuY29uc3QgTW91c2VVcCAgID0gMTtcbmNvbnN0IE1vdXNlTW92ZSA9IDI7XG5cbmZ1bmN0aW9uIHVwZGF0ZShldmVudDogTW91c2VFdmVudCkge1xuICAgIC8vIFdpbmRvd3Mgc3VwcHJlc3NlcyBtb3VzZSBldmVudHMgYXQgdGhlIGVuZCBvZiBkcmFnZ2luZyBvciByZXNpemluZyxcbiAgICAvLyBzbyB3ZSBuZWVkIHRvIGJlIHNtYXJ0IGFuZCBzeW50aGVzaXplIGJ1dHRvbiBldmVudHMuXG5cbiAgICBsZXQgZXZlbnRUeXBlOiBudW1iZXIsIGV2ZW50QnV0dG9ucyA9IGV2ZW50LmJ1dHRvbnM7XG4gICAgc3dpdGNoIChldmVudC50eXBlKSB7XG4gICAgICAgIGNhc2UgJ21vdXNlZG93bic6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZURvd247XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgfCAoMSA8PCBldmVudC5idXR0b24pOyB9XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgY2FzZSAnbW91c2V1cCc6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZVVwO1xuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zICYgfigxIDw8IGV2ZW50LmJ1dHRvbik7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VNb3ZlO1xuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zOyB9XG4gICAgICAgICAgICBicmVhaztcbiAgICB9XG5cbiAgICBsZXQgcmVsZWFzZWQgPSBidXR0b25zICYgfmV2ZW50QnV0dG9ucztcbiAgICBsZXQgcHJlc3NlZCA9IGV2ZW50QnV0dG9ucyAmIH5idXR0b25zO1xuXG4gICAgYnV0dG9ucyA9IGV2ZW50QnV0dG9ucztcblxuICAgIC8vIFN5bnRoZXNpemUgYSByZWxlYXNlLXByZXNzIHNlcXVlbmNlIGlmIHdlIGRldGVjdCBhIHByZXNzIG9mIGFuIGFscmVhZHkgcHJlc3NlZCBidXR0b24uXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VEb3duICYmICEocHJlc3NlZCAmIGV2ZW50LmJ1dHRvbikpIHtcbiAgICAgICAgcmVsZWFzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcbiAgICAgICAgcHJlc3NlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xuICAgIH1cblxuICAgIC8vIFN1cHByZXNzIGFsbCBidXR0b24gZXZlbnRzIGR1cmluZyBkcmFnZ2luZyBhbmQgcmVzaXppbmcsXG4gICAgLy8gdW5sZXNzIHRoaXMgaXMgYSBtb3VzZXVwIGV2ZW50IHRoYXQgaXMgZW5kaW5nIGEgZHJhZyBhY3Rpb24uXG4gICAgaWYgKFxuICAgICAgICBldmVudFR5cGUgIT09IE1vdXNlTW92ZSAvLyBGYXN0IHBhdGggZm9yIG1vdXNlbW92ZVxuICAgICAgICAmJiByZXNpemluZ1xuICAgICAgICB8fCAoXG4gICAgICAgICAgICBkcmFnZ2luZ1xuICAgICAgICAgICAgJiYgKFxuICAgICAgICAgICAgICAgIGV2ZW50VHlwZSA9PT0gTW91c2VEb3duXG4gICAgICAgICAgICAgICAgfHwgZXZlbnQuYnV0dG9uICE9PSAwXG4gICAgICAgICAgICApXG4gICAgICAgIClcbiAgICApIHtcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cblxuICAgIC8vIEhhbmRsZSByZWxlYXNlc1xuICAgIGlmIChyZWxlYXNlZCAmIDEpIHsgcHJpbWFyeVVwKGV2ZW50KTsgfVxuICAgIC8vIEhhbmRsZSBwcmVzc2VzXG4gICAgaWYgKHByZXNzZWQgJiAxKSB7IHByaW1hcnlEb3duKGV2ZW50KTsgfVxuXG4gICAgLy8gSGFuZGxlIG1vdXNlbW92ZVxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlTW92ZSkgeyBvbk1vdXNlTW92ZShldmVudCk7IH07XG59XG5cbmZ1bmN0aW9uIHByaW1hcnlEb3duKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XG4gICAgLy8gUmVzZXQgcmVhZGluZXNzIHN0YXRlLlxuICAgIGNhbkRyYWcgPSBmYWxzZTtcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcblxuICAgIC8vIElnbm9yZSByZXBlYXRlZCBjbGlja3Mgb24gbWFjT1MgYW5kIExpbnV4LlxuICAgIGlmICghSXNXaW5kb3dzKCkpIHtcbiAgICAgICAgaWYgKGV2ZW50LnR5cGUgPT09ICdtb3VzZWRvd24nICYmIGV2ZW50LmJ1dHRvbiA9PT0gMCAmJiBldmVudC5kZXRhaWwgIT09IDEpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIGlmIChyZXNpemVFZGdlKSB7XG4gICAgICAgIC8vIFJlYWR5IHRvIHJlc2l6ZSBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lLlxuICAgICAgICBjYW5SZXNpemUgPSB0cnVlO1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZyBvcGVyYXRpb25zIHdoZW4gb24gcmVzaXplIGVkZ2VzLlxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUmV0cmlldmUgdGFyZ2V0IGVsZW1lbnRcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XG5cbiAgICAvLyBSZWFkeSB0byBkcmFnIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUgb24gYSBkcmFnZ2FibGUgZWxlbWVudC5cbiAgICAvLyBJZ25vcmUgY2xpY2tzIG9uIHRoZSBzY3JvbGxiYXIuXG4gICAgY29uc3Qgc3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpO1xuICAgIGNhbkRyYWcgPSAoXG4gICAgICAgIHN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLXdhaWxzLWRyYWdnYWJsZVwiKS50cmltKCkgPT09IFwiZHJhZ1wiXG4gICAgICAgICYmIChcbiAgICAgICAgICAgIGV2ZW50Lm9mZnNldFggLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdMZWZ0KSA8IHRhcmdldC5jbGllbnRXaWR0aFxuICAgICAgICAgICAgJiYgZXZlbnQub2Zmc2V0WSAtIHBhcnNlRmxvYXQoc3R5bGUucGFkZGluZ1RvcCkgPCB0YXJnZXQuY2xpZW50SGVpZ2h0XG4gICAgICAgIClcbiAgICApO1xufVxuXG5mdW5jdGlvbiBwcmltYXJ5VXAoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICAvLyBTdG9wIGRyYWdnaW5nIGFuZCByZXNpemluZy5cbiAgICBjYW5EcmFnID0gZmFsc2U7XG4gICAgZHJhZ2dpbmcgPSBmYWxzZTtcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcbiAgICByZXNpemluZyA9IGZhbHNlO1xufVxuXG5jb25zdCBjdXJzb3JGb3JFZGdlID0gT2JqZWN0LmZyZWV6ZSh7XG4gICAgXCJzZS1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxuICAgIFwic3ctcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcbiAgICBcIm53LXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXG4gICAgXCJuZS1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxuICAgIFwidy1yZXNpemVcIjogXCJldy1yZXNpemVcIixcbiAgICBcIm4tcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXG4gICAgXCJzLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxuICAgIFwiZS1yZXNpemVcIjogXCJldy1yZXNpemVcIixcbn0pXG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShlZGdlPzoga2V5b2YgdHlwZW9mIGN1cnNvckZvckVkZ2UpOiB2b2lkIHtcbiAgICBpZiAoZWRnZSkge1xuICAgICAgICBpZiAoIXJlc2l6ZUVkZ2UpIHsgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yOyB9XG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gY3Vyc29yRm9yRWRnZVtlZGdlXTtcbiAgICB9IGVsc2UgaWYgKCFlZGdlICYmIHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBkZWZhdWx0Q3Vyc29yO1xuICAgIH1cblxuICAgIHJlc2l6ZUVkZ2UgPSBlZGdlIHx8IFwiXCI7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XG4gICAgaWYgKGNhblJlc2l6ZSAmJiByZXNpemVFZGdlKSB7XG4gICAgICAgIC8vIFN0YXJ0IHJlc2l6aW5nLlxuICAgICAgICByZXNpemluZyA9IHRydWU7XG4gICAgICAgIGludm9rZShcIndhaWxzOnJlc2l6ZTpcIiArIHJlc2l6ZUVkZ2UpO1xuICAgIH0gZWxzZSBpZiAoY2FuRHJhZykge1xuICAgICAgICAvLyBTdGFydCBkcmFnZ2luZy5cbiAgICAgICAgZHJhZ2dpbmcgPSB0cnVlO1xuICAgICAgICBpbnZva2UoXCJ3YWlsczpkcmFnXCIpO1xuICAgIH1cblxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xuICAgICAgICAvLyBFaXRoZXIgZHJhZyBvciByZXNpemUgaXMgb25nb2luZyxcbiAgICAgICAgLy8gcmVzZXQgcmVhZGluZXNzIGFuZCBzdG9wIHByb2Nlc3NpbmcuXG4gICAgICAgIGNhbkRyYWcgPSBjYW5SZXNpemUgPSBmYWxzZTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghcmVzaXphYmxlIHx8ICFJc1dpbmRvd3MoKSkge1xuICAgICAgICBpZiAocmVzaXplRWRnZSkgeyBzZXRSZXNpemUoKTsgfVxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgY29uc3QgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcbiAgICBjb25zdCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcblxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhcy5cbiAgICBjb25zdCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcblxuICAgIGNvbnN0IHJpZ2h0Qm9yZGVyID0gKHdpbmRvdy5vdXRlcldpZHRoIC0gZXZlbnQuY2xpZW50WCkgPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBjb25zdCBsZWZ0Qm9yZGVyID0gZXZlbnQuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGNvbnN0IHRvcEJvcmRlciA9IGV2ZW50LmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG4gICAgY29uc3QgYm90dG9tQm9yZGVyID0gKHdpbmRvdy5vdXRlckhlaWdodCAtIGV2ZW50LmNsaWVudFkpIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXIgYXJlYXMuXG4gICAgY29uc3QgcmlnaHRDb3JuZXIgPSAod2luZG93Lm91dGVyV2lkdGggLSBldmVudC5jbGllbnRYKSA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBjb25zdCBsZWZ0Q29ybmVyID0gZXZlbnQuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBjb25zdCB0b3BDb3JuZXIgPSBldmVudC5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcbiAgICBjb25zdCBib3R0b21Db3JuZXIgPSAod2luZG93Lm91dGVySGVpZ2h0IC0gZXZlbnQuY2xpZW50WSkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuXG4gICAgaWYgKCFsZWZ0Q29ybmVyICYmICF0b3BDb3JuZXIgJiYgIWJvdHRvbUNvcm5lciAmJiAhcmlnaHRDb3JuZXIpIHtcbiAgICAgICAgLy8gT3B0aW1pc2F0aW9uOiBvdXQgb2YgYWxsIGNvcm5lciBhcmVhcyBpbXBsaWVzIG91dCBvZiBib3JkZXJzLlxuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9XG4gICAgLy8gRGV0ZWN0IGNvcm5lcnMuXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIHRvcENvcm5lcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xuICAgIC8vIERldGVjdCBib3JkZXJzLlxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XG4gICAgLy8gT3V0IG9mIGJvcmRlciBhcmVhLlxuICAgIGVsc2Ugc2V0UmVzaXplKCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQXBwbGljYXRpb24pO1xuXG5jb25zdCBIaWRlTWV0aG9kID0gMDtcbmNvbnN0IFNob3dNZXRob2QgPSAxO1xuY29uc3QgUXVpdE1ldGhvZCA9IDI7XG5cbi8qKlxuICogSGlkZXMgYSBjZXJ0YWluIG1ldGhvZCBieSBjYWxsaW5nIHRoZSBIaWRlTWV0aG9kIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChIaWRlTWV0aG9kKTtcbn1cblxuLyoqXG4gKiBDYWxscyB0aGUgU2hvd01ldGhvZCBhbmQgcmV0dXJucyB0aGUgcmVzdWx0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcbn1cblxuLyoqXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKFF1aXRNZXRob2QpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBDYW5jZWxsYWJsZVByb21pc2UsIHR5cGUgQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVycyB9IGZyb20gXCIuL2NhbmNlbGxhYmxlLmpzXCI7XG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gXCIuL25hbm9pZC5qc1wiO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcbndpbmRvdy5fd2FpbHMuY2FsbEVycm9ySGFuZGxlciA9IGVycm9ySGFuZGxlcjtcblxudHlwZSBQcm9taXNlUmVzb2x2ZXJzID0gT21pdDxDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPGFueT4sIFwicHJvbWlzZVwiIHwgXCJvbmNhbmNlbGxlZFwiPlxuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYWxsKTtcbmNvbnN0IGNhbmNlbENhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbmNlbENhbGwpO1xuY29uc3QgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xuXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XG5jb25zdCBDYW5jZWxNZXRob2QgPSAwXG5cbi8qKlxuICogSG9sZHMgYWxsIHJlcXVpcmVkIGluZm9ybWF0aW9uIGZvciBhIGJpbmRpbmcgY2FsbC5cbiAqIE1heSBwcm92aWRlIGVpdGhlciBhIG1ldGhvZCBJRCBvciBhIG1ldGhvZCBuYW1lLCBidXQgbm90IGJvdGguXG4gKi9cbmV4cG9ydCB0eXBlIENhbGxPcHRpb25zID0ge1xuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kSUQ6IG51bWJlcjtcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2ROYW1lPzogbmV2ZXI7XG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xuICAgIGFyZ3M6IGFueVtdO1xufSB8IHtcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZElEPzogbmV2ZXI7XG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kTmFtZTogc3RyaW5nO1xuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cbiAgICBhcmdzOiBhbnlbXTtcbn07XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB0aHJvd24gaW4gY2FzZSB0aGUgYm91bmQgbWV0aG9kIHJldHVybnMgYW4gZXJyb3IuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBSdW50aW1lRXJyb3IjbmFtZX0gcHJvcGVydHkgaXMgXCJSdW50aW1lRXJyb3JcIi5cbiAqL1xuZXhwb3J0IGNsYXNzIFJ1bnRpbWVFcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IFJ1bnRpbWVFcnJvciBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiUnVudGltZUVycm9yXCI7XG4gICAgfVxufVxuXG4vKipcbiAqIEhhbmRsZXMgdGhlIHJlc3VsdCBvZiBhIGNhbGwgcmVxdWVzdC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgaWQgb2YgdGhlIHJlcXVlc3QgdG8gaGFuZGxlIHRoZSByZXN1bHQgZm9yLlxuICogQHBhcmFtIGRhdGEgLSBUaGUgcmVzdWx0IGRhdGEgb2YgdGhlIHJlcXVlc3QuXG4gKiBAcGFyYW0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXG4gKi9cbmZ1bmN0aW9uIHJlc3VsdEhhbmRsZXIoaWQ6IHN0cmluZywgZGF0YTogc3RyaW5nLCBpc0pTT046IGJvb2xlYW4pOiB2b2lkIHtcbiAgICBjb25zdCByZXNvbHZlcnMgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XG4gICAgaWYgKCFyZXNvbHZlcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghZGF0YSkge1xuICAgICAgICByZXNvbHZlcnMucmVzb2x2ZSh1bmRlZmluZWQpO1xuICAgIH0gZWxzZSBpZiAoIWlzSlNPTikge1xuICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShkYXRhKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycjogYW55KSB7XG4gICAgICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJjb3VsZCBub3QgcGFyc2UgcmVzdWx0OiBcIiArIGVyci5tZXNzYWdlLCB7IGNhdXNlOiBlcnIgfSkpO1xuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIEhhbmRsZXMgdGhlIGVycm9yIGZyb20gYSBjYWxsIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSBwcm9taXNlIGhhbmRsZXIuXG4gKiBAcGFyYW0gZGF0YSAtIFRoZSBlcnJvciBkYXRhIHRvIHJlamVjdCB0aGUgcHJvbWlzZSBoYW5kbGVyIHdpdGguXG4gKiBAcGFyYW0gaXNKU09OIC0gSW5kaWNhdGVzIHdoZXRoZXIgdGhlIGRhdGEgaXMgSlNPTiBvciBub3QuXG4gKi9cbmZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZDogc3RyaW5nLCBkYXRhOiBzdHJpbmcsIGlzSlNPTjogYm9vbGVhbik6IHZvaWQge1xuICAgIGNvbnN0IHJlc29sdmVycyA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcbiAgICBpZiAoIXJlc29sdmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKCFpc0pTT04pIHtcbiAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgRXJyb3IoZGF0YSkpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGxldCBlcnJvcjogYW55O1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgZXJyb3IgPSBKU09OLnBhcnNlKGRhdGEpO1xuICAgICAgICB9IGNhdGNoIChlcnI6IGFueSkge1xuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgVHlwZUVycm9yKFwiY291bGQgbm90IHBhcnNlIGVycm9yOiBcIiArIGVyci5tZXNzYWdlLCB7IGNhdXNlOiBlcnIgfSkpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG5cbiAgICAgICAgbGV0IG9wdGlvbnM6IEVycm9yT3B0aW9ucyA9IHt9O1xuICAgICAgICBpZiAoZXJyb3IuY2F1c2UpIHtcbiAgICAgICAgICAgIG9wdGlvbnMuY2F1c2UgPSBlcnJvci5jYXVzZTtcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBleGNlcHRpb247XG4gICAgICAgIHN3aXRjaCAoZXJyb3Iua2luZCkge1xuICAgICAgICAgICAgY2FzZSBcIlJlZmVyZW5jZUVycm9yXCI6XG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFJlZmVyZW5jZUVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgY2FzZSBcIlR5cGVFcnJvclwiOlxuICAgICAgICAgICAgICAgIGV4Y2VwdGlvbiA9IG5ldyBUeXBlRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICBjYXNlIFwiUnVudGltZUVycm9yXCI6XG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFJ1bnRpbWVFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IEVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICB9XG5cbiAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChleGNlcHRpb24pO1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXRyaWV2ZXMgYW5kIHJlbW92ZXMgdGhlIHJlc3BvbnNlIGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQgZnJvbSB0aGUgY2FsbFJlc3BvbnNlcyBtYXAuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIElEIG9mIHRoZSByZXNwb25zZSB0byBiZSByZXRyaWV2ZWQgYW5kIHJlbW92ZWQuXG4gKiBAcmV0dXJucyBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQsIGlmIGFueS5cbiAqL1xuZnVuY3Rpb24gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQ6IHN0cmluZyk6IFByb21pc2VSZXNvbHZlcnMgfCB1bmRlZmluZWQge1xuICAgIGNvbnN0IHJlc3BvbnNlID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICByZXR1cm4gcmVzcG9uc2U7XG59XG5cbi8qKlxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cbiAqXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIENhbGwgYSBib3VuZCBtZXRob2QgYWNjb3JkaW5nIHRvIHRoZSBnaXZlbiBjYWxsIG9wdGlvbnMuXG4gKlxuICogSW4gY2FzZSBvZiBmYWlsdXJlLCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGFuIGV4Y2VwdGlvblxuICogYW1vbmcgUmVmZXJlbmNlRXJyb3IgKHVua25vd24gbWV0aG9kKSwgVHlwZUVycm9yICh3cm9uZyBhcmd1bWVudCBjb3VudCBvciB0eXBlKSxcbiAqIHtAbGluayBSdW50aW1lRXJyb3J9IChtZXRob2QgcmV0dXJuZWQgYW4gZXJyb3IpLCBvciBvdGhlciAobmV0d29yayBvciBpbnRlcm5hbCBlcnJvcnMpLlxuICogVGhlIGV4Y2VwdGlvbiBtaWdodCBoYXZlIGEgXCJjYXVzZVwiIGZpZWxkIHdpdGggdGhlIHZhbHVlIHJldHVybmVkXG4gKiBieSB0aGUgYXBwbGljYXRpb24tIG9yIHNlcnZpY2UtbGV2ZWwgZXJyb3IgbWFyc2hhbGluZyBmdW5jdGlvbnMuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBBIG1ldGhvZCBjYWxsIGRlc2NyaXB0b3IuXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zOiBDYWxsT3B0aW9ucyk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcblxuICAgIGNvbnN0IHJlc3VsdCA9IENhbmNlbGxhYmxlUHJvbWlzZS53aXRoUmVzb2x2ZXJzPGFueT4oKTtcbiAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlOiByZXN1bHQucmVzb2x2ZSwgcmVqZWN0OiByZXN1bHQucmVqZWN0IH0pO1xuXG4gICAgY29uc3QgcmVxdWVzdCA9IGNhbGwoQ2FsbEJpbmRpbmcsIE9iamVjdC5hc3NpZ24oeyBcImNhbGwtaWRcIjogaWQgfSwgb3B0aW9ucykpO1xuICAgIGxldCBydW5uaW5nID0gZmFsc2U7XG5cbiAgICByZXF1ZXN0LnRoZW4oKCkgPT4ge1xuICAgICAgICBydW5uaW5nID0gdHJ1ZTtcbiAgICB9LCAoZXJyKSA9PiB7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgcmVzdWx0LnJlamVjdChlcnIpO1xuICAgIH0pO1xuXG4gICAgY29uc3QgY2FuY2VsID0gKCkgPT4ge1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHJldHVybiBjYW5jZWxDYWxsKENhbmNlbE1ldGhvZCwge1wiY2FsbC1pZFwiOiBpZH0pLmNhdGNoKChlcnIpID0+IHtcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoXCJFcnJvciB3aGlsZSByZXF1ZXN0aW5nIGJpbmRpbmcgY2FsbCBjYW5jZWxsYXRpb246XCIsIGVycik7XG4gICAgICAgIH0pO1xuICAgIH07XG5cbiAgICByZXN1bHQub25jYW5jZWxsZWQgPSAoKSA9PiB7XG4gICAgICAgIGlmIChydW5uaW5nKSB7XG4gICAgICAgICAgICByZXR1cm4gY2FuY2VsKCk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICByZXR1cm4gcmVxdWVzdC50aGVuKGNhbmNlbCk7XG4gICAgICAgIH1cbiAgICB9O1xuXG4gICAgcmV0dXJuIHJlc3VsdC5wcm9taXNlO1xufVxuXG4vKipcbiAqIENhbGxzIGEgYm91bmQgbWV0aG9kIGJ5IG5hbWUgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXG4gKlxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybnMgVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCeU5hbWUobWV0aG9kTmFtZTogc3RyaW5nLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gQ2FsbCh7IG1ldGhvZE5hbWUsIGFyZ3MgfSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIG51bWVyaWMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXG4gKlxuICogQHBhcmFtIG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJuIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlJRChtZXRob2RJRDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gQ2FsbCh7IG1ldGhvZElELCBhcmdzIH0pO1xufVxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2luc3BlY3QtanMvaXMtY2FsbGFibGVcblxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXG4vL1xuLy8gQ29weXJpZ2h0IChjKSAyMDE1IEpvcmRhbiBIYXJiYW5kXG4vL1xuLy8gUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weVxuLy8gb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbFxuLy8gaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0c1xuLy8gdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbFxuLy8gY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzXG4vLyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxuLy9cbi8vIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cbi8vXG4vLyBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSxcbi8vIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRVxuLy8gQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUlxuLy8gTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSxcbi8vIE9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFXG4vLyBTT0ZUV0FSRS5cblxudmFyIGZuVG9TdHIgPSBGdW5jdGlvbi5wcm90b3R5cGUudG9TdHJpbmc7XG52YXIgcmVmbGVjdEFwcGx5OiB0eXBlb2YgUmVmbGVjdC5hcHBseSB8IGZhbHNlIHwgbnVsbCA9IHR5cGVvZiBSZWZsZWN0ID09PSAnb2JqZWN0JyAmJiBSZWZsZWN0ICE9PSBudWxsICYmIFJlZmxlY3QuYXBwbHk7XG52YXIgYmFkQXJyYXlMaWtlOiBhbnk7XG52YXIgaXNDYWxsYWJsZU1hcmtlcjogYW55O1xuaWYgKHR5cGVvZiByZWZsZWN0QXBwbHkgPT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIE9iamVjdC5kZWZpbmVQcm9wZXJ0eSA9PT0gJ2Z1bmN0aW9uJykge1xuICAgIHRyeSB7XG4gICAgICAgIGJhZEFycmF5TGlrZSA9IE9iamVjdC5kZWZpbmVQcm9wZXJ0eSh7fSwgJ2xlbmd0aCcsIHtcbiAgICAgICAgICAgIGdldDogZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgICAgIHRocm93IGlzQ2FsbGFibGVNYXJrZXI7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgICAgICBpc0NhbGxhYmxlTWFya2VyID0ge307XG4gICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZSBuby10aHJvdy1saXRlcmFsXG4gICAgICAgIHJlZmxlY3RBcHBseShmdW5jdGlvbiAoKSB7IHRocm93IDQyOyB9LCBudWxsLCBiYWRBcnJheUxpa2UpO1xuICAgIH0gY2F0Y2ggKF8pIHtcbiAgICAgICAgaWYgKF8gIT09IGlzQ2FsbGFibGVNYXJrZXIpIHtcbiAgICAgICAgICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XG4gICAgICAgIH1cbiAgICB9XG59IGVsc2Uge1xuICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XG59XG5cbnZhciBjb25zdHJ1Y3RvclJlZ2V4ID0gL15cXHMqY2xhc3NcXGIvO1xudmFyIGlzRVM2Q2xhc3NGbiA9IGZ1bmN0aW9uIGlzRVM2Q2xhc3NGdW5jdGlvbih2YWx1ZTogYW55KTogYm9vbGVhbiB7XG4gICAgdHJ5IHtcbiAgICAgICAgdmFyIGZuU3RyID0gZm5Ub1N0ci5jYWxsKHZhbHVlKTtcbiAgICAgICAgcmV0dXJuIGNvbnN0cnVjdG9yUmVnZXgudGVzdChmblN0cik7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7IC8vIG5vdCBhIGZ1bmN0aW9uXG4gICAgfVxufTtcblxudmFyIHRyeUZ1bmN0aW9uT2JqZWN0ID0gZnVuY3Rpb24gdHJ5RnVuY3Rpb25Ub1N0cih2YWx1ZTogYW55KTogYm9vbGVhbiB7XG4gICAgdHJ5IHtcbiAgICAgICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgICAgIGZuVG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgIHJldHVybiB0cnVlO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbn07XG52YXIgdG9TdHIgPSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nO1xudmFyIG9iamVjdENsYXNzID0gJ1tvYmplY3QgT2JqZWN0XSc7XG52YXIgZm5DbGFzcyA9ICdbb2JqZWN0IEZ1bmN0aW9uXSc7XG52YXIgZ2VuQ2xhc3MgPSAnW29iamVjdCBHZW5lcmF0b3JGdW5jdGlvbl0nO1xudmFyIGRkYUNsYXNzID0gJ1tvYmplY3QgSFRNTEFsbENvbGxlY3Rpb25dJzsgLy8gSUUgMTFcbnZhciBkZGFDbGFzczIgPSAnW29iamVjdCBIVE1MIGRvY3VtZW50LmFsbCBjbGFzc10nO1xudmFyIGRkYUNsYXNzMyA9ICdbb2JqZWN0IEhUTUxDb2xsZWN0aW9uXSc7IC8vIElFIDktMTBcbnZhciBoYXNUb1N0cmluZ1RhZyA9IHR5cGVvZiBTeW1ib2wgPT09ICdmdW5jdGlvbicgJiYgISFTeW1ib2wudG9TdHJpbmdUYWc7IC8vIGJldHRlcjogdXNlIGBoYXMtdG9zdHJpbmd0YWdgXG5cbnZhciBpc0lFNjggPSAhKDAgaW4gWyxdKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBuby1zcGFyc2UtYXJyYXlzLCBjb21tYS1zcGFjaW5nXG5cbnZhciBpc0REQTogKHZhbHVlOiBhbnkpID0+IGJvb2xlYW4gPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKCkgeyByZXR1cm4gZmFsc2U7IH07XG5pZiAodHlwZW9mIGRvY3VtZW50ID09PSAnb2JqZWN0Jykge1xuICAgIC8vIEZpcmVmb3ggMyBjYW5vbmljYWxpemVzIEREQSB0byB1bmRlZmluZWQgd2hlbiBpdCdzIG5vdCBhY2Nlc3NlZCBkaXJlY3RseVxuICAgIHZhciBhbGwgPSBkb2N1bWVudC5hbGw7XG4gICAgaWYgKHRvU3RyLmNhbGwoYWxsKSA9PT0gdG9TdHIuY2FsbChkb2N1bWVudC5hbGwpKSB7XG4gICAgICAgIGlzRERBID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCh2YWx1ZSkge1xuICAgICAgICAgICAgLyogZ2xvYmFscyBkb2N1bWVudDogZmFsc2UgKi9cbiAgICAgICAgICAgIC8vIGluIElFIDYtOCwgdHlwZW9mIGRvY3VtZW50LmFsbCBpcyBcIm9iamVjdFwiIGFuZCBpdCdzIHRydXRoeVxuICAgICAgICAgICAgaWYgKChpc0lFNjggfHwgIXZhbHVlKSAmJiAodHlwZW9mIHZhbHVlID09PSAndW5kZWZpbmVkJyB8fCB0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnKSkge1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIHZhciBzdHIgPSB0b1N0ci5jYWxsKHZhbHVlKTtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuIChcbiAgICAgICAgICAgICAgICAgICAgICAgIHN0ciA9PT0gZGRhQ2xhc3NcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MyXG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IGRkYUNsYXNzMyAvLyBvcGVyYSAxMi4xNlxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBvYmplY3RDbGFzcyAvLyBJRSA2LThcbiAgICAgICAgICAgICAgICAgICAgKSAmJiB2YWx1ZSgnJykgPT0gbnVsbDsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBlcWVxZXFcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChlKSB7IC8qKi8gfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZnVuY3Rpb24gaXNDYWxsYWJsZVJlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55ICB7XG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxuICAgIGlmICghdmFsdWUpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHRyeSB7XG4gICAgICAgIChyZWZsZWN0QXBwbHkgYXMgYW55KSh2YWx1ZSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGlmIChlICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIH1cbiAgICByZXR1cm4gIWlzRVM2Q2xhc3NGbih2YWx1ZSkgJiYgdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xufVxuXG5mdW5jdGlvbiBpc0NhbGxhYmxlTm9SZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSB7XG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxuICAgIGlmICghdmFsdWUpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmIChoYXNUb1N0cmluZ1RhZykgeyByZXR1cm4gdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpOyB9XG4gICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgdmFyIHN0ckNsYXNzID0gdG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgaWYgKHN0ckNsYXNzICE9PSBmbkNsYXNzICYmIHN0ckNsYXNzICE9PSBnZW5DbGFzcyAmJiAhKC9eXFxbb2JqZWN0IEhUTUwvKS50ZXN0KHN0ckNsYXNzKSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICByZXR1cm4gdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xufTtcblxuZXhwb3J0IGRlZmF1bHQgcmVmbGVjdEFwcGx5ID8gaXNDYWxsYWJsZVJlZkFwcGx5IDogaXNDYWxsYWJsZU5vUmVmQXBwbHk7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCBpc0NhbGxhYmxlIGZyb20gXCIuL2NhbGxhYmxlLmpzXCI7XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB1c2VkIGFzIHJlamVjdGlvbiByZWFzb25cbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSBpcyBjYW5jZWxsZWQgc3VjY2Vzc2Z1bGx5LlxuICpcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsRXJyb3JcImAuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIGNhdXNlIHBhc3NlZCB0byB0aGUgY2FuY2VsIG1ldGhvZCwgaWYgYW55LlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsRXJyb3JgIGluc3RhbmNlLlxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gICAgICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIHRvIGJlIGZvcndhcmRlZCB0byB0aGUgRXJyb3IgY29uc3RydWN0b3IuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xuICAgICAgICBzdXBlcihtZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxFcnJvclwiO1xuICAgIH1cbn1cblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb25cbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSByZWplY3RzIGFmdGVyIGJlaW5nIGNhbmNlbGxlZCxcbiAqIG9yIHdoZW4gdGhlIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgdGhyb3dzIG9yIHJlamVjdHMuXG4gKlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiYC5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgcmVhc29uIHRoZSBwcm9taXNlIHJlamVjdGVkIHdpdGguXG4gKlxuICogQmVjYXVzZSB0aGUgb3JpZ2luYWwgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLFxuICogYSB3cmFwcGVyIHByb21pc2Ugd2lsbCBiZSBwYXNzZWQgdG8gdGhlIHVuaGFuZGxlZCByZWplY3Rpb24gbGlzdGVuZXIgaW5zdGVhZC5cbiAqIFRoZSB7QGxpbmsgcHJvbWlzZX0gcHJvcGVydHkgaG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIG9yaWdpbmFsIHByb21pc2UuXG4gKi9cbmV4cG9ydCBjbGFzcyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBIb2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgcHJvbWlzZSB0aGF0IHdhcyBjYW5jZWxsZWQgYW5kIHRoZW4gcmVqZWN0ZWQuXG4gICAgICovXG4gICAgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+O1xuXG4gICAgLyoqXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JgIGluc3RhbmNlLlxuICAgICAqIEBwYXJhbSBwcm9taXNlIC0gVGhlIHByb21pc2UgdGhhdCBjYXVzZWQgdGhlIGVycm9yIG9yaWdpbmFsbHkuXG4gICAgICogQHBhcmFtIHJlYXNvbiAtIFRoZSByZWplY3Rpb24gcmVhc29uLlxuICAgICAqIEBwYXJhbSBpbmZvIC0gQW4gb3B0aW9uYWwgaW5mb3JtYXRpdmUgbWVzc2FnZSBzcGVjaWZ5aW5nIHRoZSBjaXJjdW1zdGFuY2VzIGluIHdoaWNoIHRoZSBlcnJvciB3YXMgdGhyb3duLlxuICAgICAqICAgICAgICAgICAgICAgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiYC5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3Rvcihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHJlYXNvbj86IGFueSwgaW5mbz86IHN0cmluZykge1xuICAgICAgICBzdXBlcigoaW5mbyA/PyBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCIpICsgXCIgUmVhc29uOiBcIiArIGVycm9yTWVzc2FnZShyZWFzb24pLCB7IGNhdXNlOiByZWFzb24gfSk7XG4gICAgICAgIHRoaXMucHJvbWlzZSA9IHByb21pc2U7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcIjtcbiAgICB9XG59XG5cbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4gPSAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yID0gKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyID0gKGNhdXNlPzogYW55KSA9PiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+ID0gKHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+LCByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yKSA9PiB2b2lkO1xuXG5leHBvcnQgaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCk6IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj47XG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xufVxuXG4vKipcbiAqIFdyYXBzIGEgY2FuY2VsbGFibGUgcHJvbWlzZSBhbG9uZyB3aXRoIGl0cyByZXNvbHV0aW9uIG1ldGhvZHMuXG4gKiBUaGUgYG9uY2FuY2VsbGVkYCBmaWVsZCB3aWxsIGJlIG51bGwgaW5pdGlhbGx5IGJ1dCBtYXkgYmUgc2V0IHRvIHByb3ZpZGUgYSBjdXN0b20gY2FuY2VsbGF0aW9uIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcbiAgICByZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPjtcbiAgICByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yO1xuICAgIG9uY2FuY2VsbGVkOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgfCBudWxsO1xufVxuXG5pbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUge1xuICAgIHJlYWRvbmx5IHJvb3Q6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlO1xuICAgIHJlc29sdmluZzogYm9vbGVhbjtcbiAgICBzZXR0bGVkOiBib29sZWFuO1xuICAgIHJlYXNvbj86IENhbmNlbEVycm9yO1xufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgYmFycmllclN5bSA9IFN5bWJvbChcImJhcnJpZXJcIik7XG5jb25zdCBjYW5jZWxJbXBsU3ltID0gU3ltYm9sKFwiY2FuY2VsSW1wbFwiKTtcbmNvbnN0IHNwZWNpZXMgPSBTeW1ib2wuc3BlY2llcyA/PyBTeW1ib2woXCJzcGVjaWVzUG9seWZpbGxcIik7XG5cbi8qKlxuICogQSBwcm9taXNlIHdpdGggYW4gYXR0YWNoZWQgbWV0aG9kIGZvciBjYW5jZWxsaW5nIGxvbmctcnVubmluZyBvcGVyYXRpb25zIChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWx9KS5cbiAqIENhbmNlbGxhdGlvbiBjYW4gb3B0aW9uYWxseSBiZSBib3VuZCB0byBhbiB7QGxpbmsgQWJvcnRTaWduYWx9XG4gKiBmb3IgYmV0dGVyIGNvbXBvc2FiaWxpdHkgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbE9ufSkuXG4gKlxuICogQ2FuY2VsbGluZyBhIHBlbmRpbmcgcHJvbWlzZSB3aWxsIHJlc3VsdCBpbiBhbiBpbW1lZGlhdGUgcmVqZWN0aW9uXG4gKiB3aXRoIGFuIGluc3RhbmNlIG9mIHtAbGluayBDYW5jZWxFcnJvcn0gYXMgcmVhc29uLFxuICogYnV0IHdob2V2ZXIgc3RhcnRlZCB0aGUgcHJvbWlzZSB3aWxsIGJlIHJlc3BvbnNpYmxlXG4gKiBmb3IgYWN0dWFsbHkgYWJvcnRpbmcgdGhlIHVuZGVybHlpbmcgb3BlcmF0aW9uLlxuICogVG8gdGhpcyBwdXJwb3NlLCB0aGUgY29uc3RydWN0b3IgYW5kIGFsbCBjaGFpbmluZyBtZXRob2RzXG4gKiBhY2NlcHQgb3B0aW9uYWwgY2FuY2VsbGF0aW9uIGNhbGxiYWNrcy5cbiAqXG4gKiBJZiBhIGBDYW5jZWxsYWJsZVByb21pc2VgIHN0aWxsIHJlc29sdmVzIGFmdGVyIGhhdmluZyBiZWVuIGNhbmNlbGxlZCxcbiAqIHRoZSByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWQuIElmIGl0IHJlamVjdHMsIHRoZSByZWFzb25cbiAqIHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbixcbiAqIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGluc3RhbmNlLlxuICogVG8gZmFjaWxpdGF0ZSB0aGUgaGFuZGxpbmcgb2YgY2FuY2VsbGF0aW9uIHJlcXVlc3RzLFxuICogY2FuY2VsbGVkIGBDYW5jZWxsYWJsZVByb21pc2VgcyB3aWxsIF9ub3RfIHJlcG9ydCB1bmhhbmRsZWQgYENhbmNlbEVycm9yYHNcbiAqIHdob3NlIGBjYXVzZWAgZmllbGQgaXMgdGhlIHNhbWUgYXMgdGhlIG9uZSB3aXRoIHdoaWNoIHRoZSBjdXJyZW50IHByb21pc2Ugd2FzIGNhbmNlbGxlZC5cbiAqXG4gKiBBbGwgdXN1YWwgcHJvbWlzZSBtZXRob2RzIGFyZSBkZWZpbmVkIGFuZCByZXR1cm4gYSBgQ2FuY2VsbGFibGVQcm9taXNlYFxuICogd2hvc2UgY2FuY2VsIG1ldGhvZCB3aWxsIGNhbmNlbCB0aGUgcGFyZW50IG9wZXJhdGlvbiBhcyB3ZWxsLCBwcm9wYWdhdGluZyB0aGUgY2FuY2VsbGF0aW9uIHJlYXNvblxuICogdXB3YXJkcyB0aHJvdWdoIHByb21pc2UgY2hhaW5zLlxuICogQ29udmVyc2VseSwgY2FuY2VsbGluZyBhIHByb21pc2Ugd2lsbCBub3QgYXV0b21hdGljYWxseSBjYW5jZWwgZGVwZW5kZW50IHByb21pc2VzIGRvd25zdHJlYW06XG4gKiBgYGB0c1xuICogbGV0IHJvb3QgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMSA9IHJvb3QudGhlbigoKSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDIgPSBjaGlsZDEudGhlbigoKSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDMgPSByb290LmNhdGNoKCgpID0+IHsgLi4uIH0pO1xuICogY2hpbGQxLmNhbmNlbCgpOyAvLyBDYW5jZWxzIGNoaWxkMSBhbmQgcm9vdCwgYnV0IG5vdCBjaGlsZDIgb3IgY2hpbGQzXG4gKiBgYGBcbiAqIENhbmNlbGxpbmcgYSBwcm9taXNlIHRoYXQgaGFzIGFscmVhZHkgc2V0dGxlZCBpcyBzYWZlIGFuZCBoYXMgbm8gY29uc2VxdWVuY2UuXG4gKlxuICogVGhlIGBjYW5jZWxgIG1ldGhvZCByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IF9hbHdheXMgZnVsZmlsbHNfXG4gKiBhZnRlciB0aGUgd2hvbGUgY2hhaW4gaGFzIHByb2Nlc3NlZCB0aGUgY2FuY2VsIHJlcXVlc3RcbiAqIGFuZCBhbGwgYXR0YWNoZWQgY2FsbGJhY2tzIHVwIHRvIHRoYXQgbW9tZW50IGhhdmUgcnVuLlxuICpcbiAqIEFsbCBFUzIwMjQgcHJvbWlzZSBtZXRob2RzIChzdGF0aWMgYW5kIGluc3RhbmNlKSBhcmUgZGVmaW5lZCBvbiBDYW5jZWxsYWJsZVByb21pc2UsXG4gKiBidXQgYWN0dWFsIGF2YWlsYWJpbGl0eSBtYXkgdmFyeSB3aXRoIE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAqXG4gKiBJbiBsaW5lIHdpdGggdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmcsXG4gKiBgQ2FuY2VsbGFibGVQcm9taXNlYCBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLlxuICogRXh0ZW5kZXJzIHNob3VsZCB0YWtlIGNhcmUgdG8gcHJvdmlkZSB0aGVpciBvd24gbWV0aG9kIGltcGxlbWVudGF0aW9ucy5cbiAqIFRoaXMgbWlnaHQgYmUgcmVjb25zaWRlcmVkIGluIGNhc2UgdGhlIHByb3Bvc2FsIGlzIHJldGlyZWQuXG4gKlxuICogQ2FuY2VsbGFibGVQcm9taXNlIGlzIGEgd3JhcHBlciBhcm91bmQgdGhlIERPTSBQcm9taXNlIG9iamVjdFxuICogYW5kIGlzIGNvbXBsaWFudCB3aXRoIHRoZSBbUHJvbWlzZXMvQSsgc3BlY2lmaWNhdGlvbl0oaHR0cHM6Ly9wcm9taXNlc2FwbHVzLmNvbS8pXG4gKiAoaXQgcGFzc2VzIHRoZSBbY29tcGxpYW5jZSBzdWl0ZV0oaHR0cHM6Ly9naXRodWIuY29tL3Byb21pc2VzLWFwbHVzL3Byb21pc2VzLXRlc3RzKSlcbiAqIGlmIHNvIGlzIHRoZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uLlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGFibGVQcm9taXNlPFQ+IGV4dGVuZHMgUHJvbWlzZTxUPiBpbXBsZW1lbnRzIFByb21pc2VMaWtlPFQ+LCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICAvKiogQGludGVybmFsICovXG4gICAgcHJpdmF0ZSBbYmFycmllclN5bV0hOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiB8IG51bGw7XG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIHByaXZhdGUgcmVhZG9ubHkgW2NhbmNlbEltcGxTeW1dITogKHJlYXNvbjogQ2FuY2VsRXJyb3IpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgYENhbmNlbGxhYmxlUHJvbWlzZWAuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZXhlY3V0b3IgLSBBIGNhbGxiYWNrIHVzZWQgdG8gaW5pdGlhbGl6ZSB0aGUgcHJvbWlzZS4gVGhpcyBjYWxsYmFjayBpcyBwYXNzZWQgdHdvIGFyZ3VtZW50czpcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhIGByZXNvbHZlYCBjYWxsYmFjayB1c2VkIHRvIHJlc29sdmUgdGhlIHByb21pc2Ugd2l0aCBhIHZhbHVlXG4gICAgICogICAgICAgICAgICAgICAgICAgb3IgdGhlIHJlc3VsdCBvZiBhbm90aGVyIHByb21pc2UgKHBvc3NpYmx5IGNhbmNlbGxhYmxlKSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYSBgcmVqZWN0YCBjYWxsYmFjayB1c2VkIHRvIHJlamVjdCB0aGUgcHJvbWlzZSB3aXRoIGEgcHJvdmlkZWQgcmVhc29uIG9yIGVycm9yLlxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIHRoZSB2YWx1ZSBwcm92aWRlZCB0byB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogICAgICAgICAgICAgICAgICAgKGl0IGhhcyBhIGB0aGVuYCBfYW5kXyBhIGBjYW5jZWxgIG1ldGhvZCksXG4gICAgICogICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZm9yd2FyZGVkIHRvIHRoYXQgb2JqZWN0IGFuZCB0aGUgb25jYW5jZWxsZWQgd2lsbCBub3QgYmUgaW52b2tlZCBhbnltb3JlLlxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIGFueSBvbmUgb2YgdGhlIHR3byBjYWxsYmFja3MgaXMgY2FsbGVkIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkLFxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoZSBwcm92aWRlZCB2YWx1ZXMgd2lsbCBiZSBjYW5jZWxsZWQgYW5kIHJlc29sdmVkIGFzIHVzdWFsLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGJ1dCB0aGVpciByZXN1bHRzIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqICAgICAgICAgICAgICAgICAgIEhvd2V2ZXIsIGlmIHRoZSByZXNvbHV0aW9uIHByb2Nlc3MgdWx0aW1hdGVseSBlbmRzIHVwIGluIGEgcmVqZWN0aW9uXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhhdCBpcyBub3QgZHVlIHRvIGNhbmNlbGxhdGlvbiwgdGhlIHJlamVjdGlvbiByZWFzb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICB3aWxsIGJlIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9XG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKiBAcGFyYW0gb25jYW5jZWxsZWQgLSBJdCBpcyB0aGUgY2FsbGVyJ3MgcmVzcG9uc2liaWxpdHkgdG8gZW5zdXJlIHRoYXQgYW55IG9wZXJhdGlvblxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHN0YXJ0ZWQgYnkgdGhlIGV4ZWN1dG9yIGlzIHByb3Blcmx5IGhhbHRlZCB1cG9uIGNhbmNlbGxhdGlvbi5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBUaGlzIG9wdGlvbmFsIGNhbGxiYWNrIGNhbiBiZSB1c2VkIHRvIHRoYXQgcHVycG9zZS5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJdCB3aWxsIGJlIGNhbGxlZCBfc3luY2hyb25vdXNseV8gd2l0aCBhIGNhbmNlbGxhdGlvbiBjYXVzZVxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdoZW4gY2FuY2VsbGF0aW9uIGlzIHJlcXVlc3RlZCwgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSByZWplY3RlZFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9LCBidXQgX2JlZm9yZV9cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbnkge0BsaW5rIHRoZW59L3tAbGluayBjYXRjaH0ve0BsaW5rIGZpbmFsbHl9IGNhbGxiYWNrIHJ1bnMuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfVxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpbGwgb25seSBmdWxmaWxsIGFmdGVyIHRoZSBmb3JtZXIgaGFzIHNldHRsZWQuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVW5oYW5kbGVkIGV4Y2VwdGlvbnMgb3IgcmVqZWN0aW9ucyBmcm9tIHRoZSBjYWxsYmFjayB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgdW5oYW5kbGVkIHJlamVjdGlvbnMuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBjYWxsZWQgYmVmb3JlIGNhbmNlbGxhdGlvbiB3aXRoIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhpcyBwcm9taXNlIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gdGhhdCBwcm9taXNlLFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFuZCB0aGUgb3JpZ2luYWwgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihleGVjdXRvcjogQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4sIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKSB7XG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xuICAgICAgICBzdXBlcigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcblxuICAgICAgICBpZiAoKHRoaXMuY29uc3RydWN0b3IgYXMgYW55KVtzcGVjaWVzXSAhPT0gUHJvbWlzZSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZSBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLiBQbGVhc2UgcmVmcmFpbiBmcm9tIG92ZXJyaWRpbmcgdGhlIFtTeW1ib2wuc3BlY2llc10gc3RhdGljIHByb3BlcnR5LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0ge1xuICAgICAgICAgICAgcHJvbWlzZTogdGhpcyxcbiAgICAgICAgICAgIHJlc29sdmUsXG4gICAgICAgICAgICByZWplY3QsXG4gICAgICAgICAgICBnZXQgb25jYW5jZWxsZWQoKSB7IHJldHVybiBvbmNhbmNlbGxlZCA/PyBudWxsOyB9LFxuICAgICAgICAgICAgc2V0IG9uY2FuY2VsbGVkKGNiKSB7IG9uY2FuY2VsbGVkID0gY2IgPz8gdW5kZWZpbmVkOyB9XG4gICAgICAgIH07XG5cbiAgICAgICAgY29uc3Qgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xuICAgICAgICAgICAgZ2V0IHJvb3QoKSB7IHJldHVybiBzdGF0ZTsgfSxcbiAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXG4gICAgICAgICAgICBzZXR0bGVkOiBmYWxzZVxuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFNldHVwIGNhbmNlbGxhdGlvbiBzeXN0ZW0uXG4gICAgICAgIHZvaWQgT2JqZWN0LmRlZmluZVByb3BlcnRpZXModGhpcywge1xuICAgICAgICAgICAgW2JhcnJpZXJTeW1dOiB7XG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogdHJ1ZSxcbiAgICAgICAgICAgICAgICB2YWx1ZTogbnVsbFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIFtjYW5jZWxJbXBsU3ltXToge1xuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHZhbHVlOiBjYW5jZWxsZXJGb3IocHJvbWlzZSwgc3RhdGUpXG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuXG4gICAgICAgIC8vIFJ1biB0aGUgYWN0dWFsIGV4ZWN1dG9yLlxuICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIHN0YXRlKTtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGV4ZWN1dG9yKHJlc29sdmVyRm9yKHByb21pc2UsIHN0YXRlKSwgcmVqZWN0b3IpO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHtcbiAgICAgICAgICAgICAgICBjb25zb2xlLmxvZyhcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gQ2FuY2VsbGFibGVQcm9taXNlIGV4ZWN1dG9yLlwiLCBlcnIpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2FuY2VscyBpbW1lZGlhdGVseSB0aGUgZXhlY3V0aW9uIG9mIHRoZSBvcGVyYXRpb24gYXNzb2NpYXRlZCB3aXRoIHRoaXMgcHJvbWlzZS5cbiAgICAgKiBUaGUgcHJvbWlzZSByZWplY3RzIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9IGluc3RhbmNlIGFzIHJlYXNvbixcbiAgICAgKiB3aXRoIHRoZSB7QGxpbmsgQ2FuY2VsRXJyb3IjY2F1c2V9IHByb3BlcnR5IHNldCB0byB0aGUgZ2l2ZW4gYXJndW1lbnQsIGlmIGFueS5cbiAgICAgKlxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIGFmdGVyIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQ7XG4gICAgICogcmVwZWF0ZWQgY2FsbHMgaW4gcGFydGljdWxhciBhcmUgc2FmZSwgYnV0IG9ubHkgdGhlIGZpcnN0IG9uZVxuICAgICAqIHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXG4gICAgICpcbiAgICAgKiBUaGUgYENhbmNlbEVycm9yYCBleGNlcHRpb24gX25lZWQgbm90XyBiZSBoYW5kbGVkIGV4cGxpY2l0bHkgX29uIHRoZSBwcm9taXNlcyB0aGF0IGFyZSBiZWluZyBjYW5jZWxsZWQ6X1xuICAgICAqIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpdGggbm8gYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXIgZG9lcyBub3QgdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LlxuICAgICAqIFRoZXJlZm9yZSwgdGhlIGZvbGxvd2luZyBpZGlvbXMgYXJlIGFsbCBlcXVhbGx5IGNvcnJlY3Q6XG4gICAgICogYGBgdHNcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLmNhbmNlbCgpO1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhbmNlbCgpO1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhdGNoKC4uLikuY2FuY2VsKCk7XG4gICAgICogYGBgXG4gICAgICogV2hlbmV2ZXIgc29tZSBjYW5jZWxsZWQgcHJvbWlzZSBpbiBhIGNoYWluIHJlamVjdHMgd2l0aCBhIGBDYW5jZWxFcnJvcmBcbiAgICAgKiB3aXRoIHRoZSBzYW1lIGNhbmNlbGxhdGlvbiBjYXVzZSBhcyBpdHNlbGYsIHRoZSBlcnJvciB3aWxsIGJlIGRpc2NhcmRlZCBzaWxlbnRseS5cbiAgICAgKiBIb3dldmVyLCB0aGUgYENhbmNlbEVycm9yYCBfd2lsbCBzdGlsbCBiZSBkZWxpdmVyZWRfIHRvIGFsbCBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlcnNcbiAgICAgKiBhZGRlZCBieSB7QGxpbmsgdGhlbn0gYW5kIHJlbGF0ZWQgbWV0aG9kczpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gICAgICogY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS5jYXRjaChjb25zb2xlLmxvZyk7XG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIEEgQ2FuY2VsRXJyb3IgaXMgcHJpbnRlZCB0byB0aGUgY29uc29sZS5cbiAgICAgKiBgYGBcbiAgICAgKiBJZiB0aGUgYENhbmNlbEVycm9yYCBpcyBub3QgaGFuZGxlZCBkb3duc3RyZWFtIGJ5IHRoZSB0aW1lIGl0IHJlYWNoZXNcbiAgICAgKiBhIF9ub24tY2FuY2VsbGVkXyBwcm9taXNlLCBpdCBfd2lsbF8gdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LFxuICAgICAqIGp1c3QgbGlrZSBub3JtYWwgcmVqZWN0aW9ucyB3b3VsZDpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gICAgICogbGV0IGNoYWluZWQgPSBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLnRoZW4oKCkgPT4geyAuLi4gfSk7IC8vIE5vIGNhdGNoLi4uXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIFVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQgb24gY2hhaW5lZCFcbiAgICAgKiBgYGBcbiAgICAgKiBUaGVyZWZvcmUsIGl0IGlzIGltcG9ydGFudCB0byBlaXRoZXIgY2FuY2VsIHdob2xlIHByb21pc2UgY2hhaW5zIGZyb20gdGhlaXIgdGFpbCxcbiAgICAgKiBhcyBzaG93biBpbiB0aGUgY29ycmVjdCBpZGlvbXMgYWJvdmUsIG9yIHRha2UgY2FyZSBvZiBoYW5kbGluZyBlcnJvcnMgZXZlcnl3aGVyZS5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgY2FuY2VsbGFibGUgcHJvbWlzZSB0aGF0IF9mdWxmaWxsc18gYWZ0ZXIgdGhlIGNhbmNlbCBjYWxsYmFjayAoaWYgYW55KVxuICAgICAqIGFuZCBhbGwgaGFuZGxlcnMgYXR0YWNoZWQgdXAgdG8gdGhlIGNhbGwgdG8gY2FuY2VsIGhhdmUgcnVuLlxuICAgICAqIElmIHRoZSBjYW5jZWwgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBieSBgY2FuY2VsYFxuICAgICAqIHdpbGwgYWxzbyB3YWl0IGZvciB0aGF0IHRoZW5hYmxlIHRvIHNldHRsZS5cbiAgICAgKiBUaGlzIGVuYWJsZXMgY2FsbGVycyB0byB3YWl0IGZvciB0aGUgY2FuY2VsbGVkIG9wZXJhdGlvbiB0byB0ZXJtaW5hdGVcbiAgICAgKiB3aXRob3V0IGJlaW5nIGZvcmNlZCB0byBoYW5kbGUgcG90ZW50aWFsIGVycm9ycyBhdCB0aGUgY2FsbCBzaXRlLlxuICAgICAqIGBgYHRzXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCkudGhlbigoKSA9PiB7XG4gICAgICogICAgIC8vIENsZWFudXAgZmluaXNoZWQsIGl0J3Mgc2FmZSB0byBkbyBzb21ldGhpbmcgZWxzZS5cbiAgICAgKiB9LCAoZXJyKSA9PiB7XG4gICAgICogICAgIC8vIFVucmVhY2hhYmxlOiB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIGNhbmNlbCB3aWxsIG5ldmVyIHJlamVjdC5cbiAgICAgKiB9KTtcbiAgICAgKiBgYGBcbiAgICAgKiBOb3RlIHRoYXQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBfbm90XyBoYW5kbGUgaW1wbGljaXRseSBhbnkgcmVqZWN0aW9uXG4gICAgICogdGhhdCBtaWdodCBoYXZlIG9jY3VycmVkIGFscmVhZHkgaW4gdGhlIGNhbmNlbGxlZCBjaGFpbi5cbiAgICAgKiBJdCB3aWxsIGp1c3QgdHJhY2sgd2hldGhlciByZWdpc3RlcmVkIGhhbmRsZXJzIGhhdmUgYmVlbiBleGVjdXRlZCBvciBub3QuXG4gICAgICogVGhlcmVmb3JlLCB1bmhhbmRsZWQgcmVqZWN0aW9ucyB3aWxsIG5ldmVyIGJlIHNpbGVudGx5IGhhbmRsZWQgYnkgY2FsbGluZyBjYW5jZWwuXG4gICAgICovXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIC8vIElOVkFSSUFOVDogdGhlIHJlc3VsdCBvZiB0aGlzW2NhbmNlbEltcGxTeW1dIGFuZCB0aGUgYmFycmllciBkbyBub3QgZXZlciByZWplY3QuXG4gICAgICAgICAgICAvLyBVbmZvcnR1bmF0ZWx5IG1hY09TIEhpZ2ggU2llcnJhIGRvZXMgbm90IHN1cHBvcnQgUHJvbWlzZS5hbGxTZXR0bGVkLlxuICAgICAgICAgICAgUHJvbWlzZS5hbGwoW1xuICAgICAgICAgICAgICAgIHRoaXNbY2FuY2VsSW1wbFN5bV0obmV3IENhbmNlbEVycm9yKFwiUHJvbWlzZSBjYW5jZWxsZWQuXCIsIHsgY2F1c2UgfSkpLFxuICAgICAgICAgICAgICAgIGN1cnJlbnRCYXJyaWVyKHRoaXMpXG4gICAgICAgICAgICBdKS50aGVuKCgpID0+IHJlc29sdmUoKSwgKCkgPT4gcmVzb2x2ZSgpKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQmluZHMgcHJvbWlzZSBjYW5jZWxsYXRpb24gdG8gdGhlIGFib3J0IGV2ZW50IG9mIHRoZSBnaXZlbiB7QGxpbmsgQWJvcnRTaWduYWx9LlxuICAgICAqIElmIHRoZSBzaWduYWwgaGFzIGFscmVhZHkgYWJvcnRlZCwgdGhlIHByb21pc2Ugd2lsbCBiZSBjYW5jZWxsZWQgaW1tZWRpYXRlbHkuXG4gICAgICogV2hlbiBlaXRoZXIgY29uZGl0aW9uIGlzIHZlcmlmaWVkLCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlIHdpbGwgYmUgc2V0XG4gICAgICogdG8gdGhlIHNpZ25hbCdzIGFib3J0IHJlYXNvbiAoc2VlIHtAbGluayBBYm9ydFNpZ25hbCNyZWFzb259KS5cbiAgICAgKlxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIChvciBpZiB0aGUgc2lnbmFsIGFib3J0cykgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkLlxuICAgICAqIE9ubHkgdGhlIGZpcnN0IHNpZ25hbCB0byBhYm9ydCB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxuICAgICAqXG4gICAgICogRm9yIG1vcmUgZGV0YWlscyBhYm91dCB0aGUgY2FuY2VsbGF0aW9uIHByb2Nlc3MsXG4gICAgICogc2VlIHtAbGluayBjYW5jZWx9IGFuZCB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAgY29uc3RydWN0b3IuXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBlbmFibGVzIGBhd2FpdGBpbmcgY2FuY2VsbGFibGUgcHJvbWlzZXMgd2l0aG91dCBoYXZpbmdcbiAgICAgKiB0byBzdG9yZSB0aGVtIGZvciBmdXR1cmUgY2FuY2VsbGF0aW9uLCBlLmcuOlxuICAgICAqIGBgYHRzXG4gICAgICogYXdhaXQgbG9uZ1J1bm5pbmdPcGVyYXRpb24oKS5jYW5jZWxPbihzaWduYWwpO1xuICAgICAqIGBgYFxuICAgICAqIGluc3RlYWQgb2Y6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQgPSBsb25nUnVubmluZ09wZXJhdGlvbigpO1xuICAgICAqIGF3YWl0IHByb21pc2VUb0JlQ2FuY2VsbGVkO1xuICAgICAqIGBgYFxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhpcyBwcm9taXNlLCBmb3IgbWV0aG9kIGNoYWluaW5nLlxuICAgICAqL1xuICAgIGNhbmNlbE9uKHNpZ25hbDogQWJvcnRTaWduYWwpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBpZiAoc2lnbmFsLmFib3J0ZWQpIHtcbiAgICAgICAgICAgIHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbilcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHNpZ25hbC5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbiksIHtjYXB0dXJlOiB0cnVlfSk7XG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gdGhpcztcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBjYWxsYmFja3MgZm9yIHRoZSByZXNvbHV0aW9uIGFuZC9vciByZWplY3Rpb24gb2YgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBvbmZ1bGZpbGxlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlc29sdmVkLlxuICAgICAqIEBwYXJhbSBvbnJlamVjdGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVqZWN0ZWQuXG4gICAgICogQHJldHVybnMgQSBgQ2FuY2VsbGFibGVQcm9taXNlYCBmb3IgdGhlIGNvbXBsZXRpb24gb2Ygd2hpY2hldmVyIGNhbGxiYWNrIGlzIGV4ZWN1dGVkLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGhvb2tlZCB1cCB0byBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHVwIHRoZSBjaGFpbiwgYnV0IG5vdCBkb3duOlxuICAgICAqXG4gICAgICogICAtIGlmIHRoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIGJlIGludm9rZWQgd2l0aCBhIGBDYW5jZWxFcnJvcmBcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIF93aWxsIHJlc29sdmUgcmVndWxhcmx5XyB3aXRoIGl0cyByZXN1bHQ7XG4gICAgICogICAtIGNvbnZlcnNlbHksIGlmIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgX3RoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQgdG9vO19cbiAgICAgKiAgICAgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgc3RpbGwgYmUgaW52b2tlZCB3aXRoIHRoZSBwYXJlbnQncyBgQ2FuY2VsRXJyb3JgLFxuICAgICAqICAgICBidXQgaXRzIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZFxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhIGBDYW5jZWxFcnJvcmAgYXMgd2VsbC5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBlaXRoZXIgY2FsbGJhY2sgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKi9cbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+IHtcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLnRoZW4gY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIE5PVEU6IFR5cGVTY3JpcHQncyBidWlsdC1pbiB0eXBlIGZvciB0aGVuIGlzIGJyb2tlbixcbiAgICAgICAgLy8gYXMgaXQgYWxsb3dzIHNwZWNpZnlpbmcgYW4gYXJiaXRyYXJ5IFRSZXN1bHQxICE9IFQgZXZlbiB3aGVuIG9uZnVsZmlsbGVkIGlzIG5vdCBhIGZ1bmN0aW9uLlxuICAgICAgICAvLyBXZSBjYW5ub3QgZml4IGl0IGlmIHdlIHdhbnQgdG8gQ2FuY2VsbGFibGVQcm9taXNlIHRvIGltcGxlbWVudCBQcm9taXNlTGlrZTxUPi5cblxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25mdWxmaWxsZWQpKSB7IG9uZnVsZmlsbGVkID0gaWRlbnRpdHkgYXMgYW55OyB9XG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbnJlamVjdGVkKSkgeyBvbnJlamVjdGVkID0gdGhyb3dlcjsgfVxuXG4gICAgICAgIGlmIChvbmZ1bGZpbGxlZCA9PT0gaWRlbnRpdHkgJiYgb25yZWplY3RlZCA9PSB0aHJvd2VyKSB7XG4gICAgICAgICAgICAvLyBTaG9ydGN1dCBmb3IgdHJpdmlhbCBhcmd1bWVudHMuXG4gICAgICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSkgPT4gcmVzb2x2ZSh0aGlzIGFzIGFueSkpO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgYmFycmllcjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSB7fTtcbiAgICAgICAgdGhpc1tiYXJyaWVyU3ltXSA9IGJhcnJpZXI7XG5cbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgdm9pZCBzdXBlci50aGVuKFxuICAgICAgICAgICAgICAgICh2YWx1ZSkgPT4ge1xuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xuXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9uZnVsZmlsbGVkISh2YWx1ZSkpO1xuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICAocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xuXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9ucmVqZWN0ZWQhKHJlYXNvbikpO1xuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgKTtcbiAgICAgICAgfSwgYXN5bmMgKGNhdXNlPykgPT4ge1xuICAgICAgICAgICAgLy9jYW5jZWxsZWQgPSB0cnVlO1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICByZXR1cm4gb25jYW5jZWxsZWQ/LihjYXVzZSk7XG4gICAgICAgICAgICB9IGZpbmFsbHkge1xuICAgICAgICAgICAgICAgIGF3YWl0IHRoaXMuY2FuY2VsKGNhdXNlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayBmb3Igb25seSB0aGUgcmVqZWN0aW9uIG9mIHRoZSBQcm9taXNlLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIEl0IGlzIGVxdWl2YWxlbnQgdG9cbiAgICAgKiBgYGB0c1xuICAgICAqIGNhbmNlbGxhYmxlUHJvbWlzZS50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xuICAgICAqIGBgYFxuICAgICAqIGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgYG9ucmVqZWN0ZWRgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxuICAgICAqL1xuICAgIGNhdGNoPFRSZXN1bHQgPSBuZXZlcj4ob25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IChQcm9taXNlTGlrZTxUUmVzdWx0PiB8IFRSZXN1bHQpKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQgfCBUUmVzdWx0PiB7XG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayB0aGF0IGlzIGludm9rZWQgd2hlbiB0aGUgQ2FuY2VsbGFibGVQcm9taXNlIGlzIHNldHRsZWQgKGZ1bGZpbGxlZCBvciByZWplY3RlZCkuIFRoZVxuICAgICAqIHJlc29sdmVkIHZhbHVlIGNhbm5vdCBiZSBhY2Nlc3NlZCBvciBtb2RpZmllZCBmcm9tIHRoZSBjYWxsYmFjay5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHNldHRsZSBpbiB0aGUgc2FtZSBzdGF0ZSBhcyB0aGUgb3JpZ2luYWwgb25lXG4gICAgICogYWZ0ZXIgdGhlIHByb3ZpZGVkIGNhbGxiYWNrIGhhcyBjb21wbGV0ZWQgZXhlY3V0aW9uLFxuICAgICAqIHVubGVzcyB0aGUgY2FsbGJhY2sgdGhyb3dzIG9yIHJldHVybnMgYSByZWplY3RpbmcgcHJvbWlzZSxcbiAgICAgKiBpbiB3aGljaCBjYXNlIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IGFzIHdlbGwuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIE9uY2UgdGhlIHBhcmVudCBwcm9taXNlIHNldHRsZXMsIHRoZSBgb25maW5hbGx5YCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGlzIGltcGxlbWVudGVkIGluIHRlcm1zIG9mIHtAbGluayB0aGVufSBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cbiAgICAgKiBJdCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgYG9uZmluYWxseWAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXG4gICAgICovXG4gICAgZmluYWxseShvbmZpbmFsbHk/OiAoKCkgPT4gdm9pZCkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS5maW5hbGx5IGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25maW5hbGx5KSkge1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMudGhlbihvbmZpbmFsbHksIG9uZmluYWxseSwgb25jYW5jZWxsZWQpO1xuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbihcbiAgICAgICAgICAgICh2YWx1ZSkgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4gdmFsdWUpLFxuICAgICAgICAgICAgKHJlYXNvbj8pID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHsgdGhyb3cgcmVhc29uOyB9KSxcbiAgICAgICAgICAgIG9uY2FuY2VsbGVkLFxuICAgICAgICApO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFdlIHVzZSB0aGUgYFtTeW1ib2wuc3BlY2llc11gIHN0YXRpYyBwcm9wZXJ0eSwgaWYgYXZhaWxhYmxlLFxuICAgICAqIHRvIGRpc2FibGUgdGhlIGJ1aWx0LWluIGF1dG9tYXRpYyBzdWJjbGFzc2luZyBmZWF0dXJlcyBmcm9tIHtAbGluayBQcm9taXNlfS5cbiAgICAgKiBJdCBpcyBjcml0aWNhbCBmb3IgcGVyZm9ybWFuY2UgcmVhc29ucyB0aGF0IGV4dGVuZGVycyBkbyBub3Qgb3ZlcnJpZGUgdGhpcy5cbiAgICAgKiBPbmNlIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nXG4gICAgICogaXMgZWl0aGVyIGFjY2VwdGVkIG9yIHJldGlyZWQsIHRoaXMgaW1wbGVtZW50YXRpb24gd2lsbCBoYXZlIHRvIGJlIHJldmlzZWQgYWNjb3JkaW5nbHkuXG4gICAgICpcbiAgICAgKiBAaWdub3JlXG4gICAgICogQGludGVybmFsXG4gICAgICovXG4gICAgc3RhdGljIGdldCBbc3BlY2llc10oKSB7XG4gICAgICAgIHJldHVybiBQcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSwgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgUHJvbWlzZSBpcyByZWplY3RlZC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYWxsPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+W10+O1xuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IEF3YWl0ZWQ8VFtQXT47IH0+O1xuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGwoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUgb3IgcmVqZWN0LlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFQ+PltdPjtcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogUHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUW1BdPj47IH0+O1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsU2V0dGxlZChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRoZSBhbnkgZnVuY3Rpb24gcmV0dXJucyBhIHByb21pc2UgdGhhdCBpcyBmdWxmaWxsZWQgYnkgdGhlIGZpcnN0IGdpdmVuIHByb21pc2UgdG8gYmUgZnVsZmlsbGVkLFxuICAgICAqIG9yIHJlamVjdGVkIHdpdGggYW4gQWdncmVnYXRlRXJyb3IgY29udGFpbmluZyBhbiBhcnJheSBvZiByZWplY3Rpb24gcmVhc29uc1xuICAgICAqIGlmIGFsbCBvZiB0aGUgZ2l2ZW4gcHJvbWlzZXMgYXJlIHJlamVjdGVkLlxuICAgICAqIEl0IHJlc29sdmVzIGFsbCBlbGVtZW50cyBvZiB0aGUgcGFzc2VkIGl0ZXJhYmxlIHRvIHByb21pc2VzIGFzIGl0IHJ1bnMgdGhpcyBhbGdvcml0aG0uXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFueTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbnkoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIG9yIHJlamVjdGVkIHdoZW4gYW55IG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyBhcmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJhY2U8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmFjZShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBjYW5jZWxsZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgY2FuY2VsPFQgPSBuZXZlcj4oY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBjb25zdCBwID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XG4gICAgICAgIHAuY2FuY2VsKGNhdXNlKTtcbiAgICAgICAgcmV0dXJuIHA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBjYW5jZWxzXG4gICAgICogYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCBjYXVzZS5cbiAgICAgKlxuICAgICAqIElmIHRoZSB7QGxpbmsgQWJvcnRTaWduYWwudGltZW91dH0gZmFjdG9yeSBtZXRob2QgaXMgYXZhaWxhYmxlLFxuICAgICAqIGl0IGlzIHVzZWQgdG8gYmFzZSB0aGUgdGltZW91dCBvbiBfYWN0aXZlXyB0aW1lIHJhdGhlciB0aGFuIF9lbGFwc2VkXyB0aW1lLlxuICAgICAqIE90aGVyd2lzZSwgYHRpbWVvdXRgIGZhbGxzIGJhY2sgdG8ge0BsaW5rIHNldFRpbWVvdXR9LlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHRpbWVvdXQ8VCA9IG5ldmVyPihtaWxsaXNlY29uZHM6IG51bWJlciwgY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XG4gICAgICAgIGlmIChBYm9ydFNpZ25hbCAmJiB0eXBlb2YgQWJvcnRTaWduYWwgPT09ICdmdW5jdGlvbicgJiYgQWJvcnRTaWduYWwudGltZW91dCAmJiB0eXBlb2YgQWJvcnRTaWduYWwudGltZW91dCA9PT0gJ2Z1bmN0aW9uJykge1xuICAgICAgICAgICAgQWJvcnRTaWduYWwudGltZW91dChtaWxsaXNlY29uZHMpLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSwgbWlsbGlzZWNvbmRzKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dC5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHNsZWVwKG1pbGxpc2Vjb25kczogbnVtYmVyKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXJcbiAgICAgKiB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgc2xlZXA8VD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xuICAgIHN0YXRpYyBzbGVlcDxUID0gdm9pZD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlPzogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHJlc29sdmUodmFsdWUhKSwgbWlsbGlzZWNvbmRzKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZWplY3RlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCByZWFzb24uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVqZWN0PFQgPSBuZXZlcj4ocmVhc29uPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKF8sIHJlamVjdCkgPT4gcmVqZWN0KHJlYXNvbikpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmUoKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgcmVzb2x2ZTxUID0gdm9pZD4odmFsdWU/OiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj4ge1xuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpIHtcbiAgICAgICAgICAgIC8vIE9wdGltaXNlIGZvciBjYW5jZWxsYWJsZSBwcm9taXNlcy5cbiAgICAgICAgICAgIHJldHVybiB2YWx1ZTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+KChyZXNvbHZlKSA9PiByZXNvbHZlKHZhbHVlKSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgYW5kIHJldHVybnMgaXQgaW4gYW4gb2JqZWN0LCBhbG9uZyB3aXRoIGl0cyByZXNvbHZlIGFuZCByZWplY3QgZnVuY3Rpb25zXG4gICAgICogYW5kIGEgZ2V0dGVyL3NldHRlciBmb3IgdGhlIGNhbmNlbGxhdGlvbiBjYWxsYmFjay5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgd2l0aFJlc29sdmVyczxUPigpOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICAgICAgbGV0IHJlc3VsdDogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHsgb25jYW5jZWxsZWQ6IG51bGwgfSBhcyBhbnk7XG4gICAgICAgIHJlc3VsdC5wcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICByZXN1bHQucmVzb2x2ZSA9IHJlc29sdmU7XG4gICAgICAgICAgICByZXN1bHQucmVqZWN0ID0gcmVqZWN0O1xuICAgICAgICB9LCAoY2F1c2U/OiBhbnkpID0+IHsgcmVzdWx0Lm9uY2FuY2VsbGVkPy4oY2F1c2UpOyB9KTtcbiAgICAgICAgcmV0dXJuIHJlc3VsdDtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgY2FuY2VsbGF0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSByZXN1bHRpbmcgZnVuY3Rpb24gZG9lcyBub3QgcmVqZWN0LlxuICovXG5mdW5jdGlvbiBjYW5jZWxsZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKSB7XG4gICAgbGV0IGNhbmNlbGxhdGlvblByb21pc2U6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9IHVuZGVmaW5lZDtcblxuICAgIHJldHVybiAocmVhc29uOiBDYW5jZWxFcnJvcik6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9PiB7XG4gICAgICAgIGlmICghc3RhdGUuc2V0dGxlZCkge1xuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBzdGF0ZS5yZWFzb24gPSByZWFzb247XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xuXG4gICAgICAgICAgICAvLyBBdHRhY2ggYW4gZXJyb3IgaGFuZGxlciB0aGF0IGlnbm9yZXMgdGhpcyBzcGVjaWZpYyByZWplY3Rpb24gcmVhc29uIGFuZCBub3RoaW5nIGVsc2UuXG4gICAgICAgICAgICAvLyBJbiB0aGVvcnksIGEgc2FuZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uIGF0IHRoaXMgcG9pbnRcbiAgICAgICAgICAgIC8vIHNob3VsZCBhbHdheXMgcmVqZWN0IHdpdGggb3VyIGNhbmNlbGxhdGlvbiByZWFzb24sXG4gICAgICAgICAgICAvLyBoZW5jZSB0aGUgaGFuZGxlciB3aWxsIG5ldmVyIHRocm93LlxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnByb3RvdHlwZS50aGVuLmNhbGwocHJvbWlzZS5wcm9taXNlLCB1bmRlZmluZWQsIChlcnIpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoZXJyICE9PSByZWFzb24pIHtcbiAgICAgICAgICAgICAgICAgICAgdGhyb3cgZXJyO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0pO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gSWYgcmVhc29uIGlzIG5vdCBzZXQsIHRoZSBwcm9taXNlIHJlc29sdmVkIHJlZ3VsYXJseSwgaGVuY2Ugd2UgbXVzdCBub3QgY2FsbCBvbmNhbmNlbGxlZC5cbiAgICAgICAgLy8gSWYgb25jYW5jZWxsZWQgaXMgdW5zZXQsIG5vIG5lZWQgdG8gZ28gYW55IGZ1cnRoZXIuXG4gICAgICAgIGlmICghc3RhdGUucmVhc29uIHx8ICFwcm9taXNlLm9uY2FuY2VsbGVkKSB7IHJldHVybjsgfVxuXG4gICAgICAgIGNhbmNlbGxhdGlvblByb21pc2UgPSBuZXcgUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICByZXNvbHZlKHByb21pc2Uub25jYW5jZWxsZWQhKHN0YXRlLnJlYXNvbiEuY2F1c2UpKTtcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pLmNhdGNoKChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xuICAgICAgICB9KTtcblxuICAgICAgICAvLyBVbnNldCBvbmNhbmNlbGxlZCB0byBwcmV2ZW50IHJlcGVhdGVkIGNhbGxzLlxuICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gbnVsbDtcblxuICAgICAgICByZXR1cm4gY2FuY2VsbGF0aW9uUHJvbWlzZTtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVzb2x1dGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICovXG5mdW5jdGlvbiByZXNvbHZlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiB7XG4gICAgcmV0dXJuICh2YWx1ZSkgPT4ge1xuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xuXG4gICAgICAgIGlmICh2YWx1ZSA9PT0gcHJvbWlzZS5wcm9taXNlKSB7XG4gICAgICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QobmV3IFR5cGVFcnJvcihcIkEgcHJvbWlzZSBjYW5ub3QgYmUgcmVzb2x2ZWQgd2l0aCBpdHNlbGYuXCIpKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh2YWx1ZSAhPSBudWxsICYmICh0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ2Z1bmN0aW9uJykpIHtcbiAgICAgICAgICAgIGxldCB0aGVuOiBhbnk7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHRoZW4gPSAodmFsdWUgYXMgYW55KS50aGVuO1xuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKHRoZW4pKSB7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgbGV0IGNhbmNlbCA9ICh2YWx1ZSBhcyBhbnkpLmNhbmNlbDtcbiAgICAgICAgICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUoY2FuY2VsKSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3Qgb25jYW5jZWxsZWQgPSAoY2F1c2U/OiBhbnkpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xuICAgICAgICAgICAgICAgICAgICAgICAgfTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChzdGF0ZS5yZWFzb24pIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBJZiBhbHJlYWR5IGNhbmNlbGxlZCwgcHJvcGFnYXRlIGNhbmNlbGxhdGlvbi5cbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSBjYW5jZWxsZXIgYWxnb3JpdGhtIGRvZXMgbm90IHJlamVjdFxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIHNvIGl0IGNhbiBiZSBkaXNjYXJkZWQgc2FmZWx5LlxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgY2FuY2VsbGVyRm9yKHsgLi4ucHJvbWlzZSwgb25jYW5jZWxsZWQgfSwgc3RhdGUpKHN0YXRlLnJlYXNvbik7XG4gICAgICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBvbmNhbmNlbGxlZDtcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0gY2F0Y2gge31cblxuICAgICAgICAgICAgICAgIGNvbnN0IG5ld1N0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcbiAgICAgICAgICAgICAgICAgICAgcm9vdDogc3RhdGUucm9vdCxcbiAgICAgICAgICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHNldHRsZWQoKSB7IHJldHVybiB0aGlzLnJvb3Quc2V0dGxlZCB9LFxuICAgICAgICAgICAgICAgICAgICBzZXQgc2V0dGxlZCh2YWx1ZSkgeyB0aGlzLnJvb3Quc2V0dGxlZCA9IHZhbHVlOyB9LFxuICAgICAgICAgICAgICAgICAgICBnZXQgcmVhc29uKCkgeyByZXR1cm4gdGhpcy5yb290LnJlYXNvbiB9XG4gICAgICAgICAgICAgICAgfTtcblxuICAgICAgICAgICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgbmV3U3RhdGUpO1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkodGhlbiwgdmFsdWUsIFtyZXNvbHZlckZvcihwcm9taXNlLCBuZXdTdGF0ZSksIHJlamVjdG9yXSk7XG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJldHVybjsgLy8gSU1QT1JUQU5UIVxuICAgICAgICAgICAgfVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICBwcm9taXNlLnJlc29sdmUodmFsdWUpO1xuICAgIH07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVqZWN0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKi9cbmZ1bmN0aW9uIHJlamVjdG9yRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yIHtcbiAgICByZXR1cm4gKHJlYXNvbj8pID0+IHtcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcblxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkge1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBpZiAocmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgc3RhdGUucmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgT2JqZWN0LmlzKHJlYXNvbi5jYXVzZSwgc3RhdGUucmVhc29uLmNhdXNlKSkge1xuICAgICAgICAgICAgICAgICAgICAvLyBTd2FsbG93IGxhdGUgcmVqZWN0aW9ucyB0aGF0IGFyZSBDYW5jZWxFcnJvcnMgd2hvc2UgY2FuY2VsbGF0aW9uIGNhdXNlIGlzIHRoZSBzYW1lIGFzIG91cnMuXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9IGNhdGNoIHt9XG5cbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogQ2FuY2VscyBhbGwgdmFsdWVzIGluIGFuIGFycmF5IHRoYXQgbG9vayBsaWtlIGNhbmNlbGxhYmxlIHRoZW5hYmxlcy5cbiAqIFJldHVybnMgYSBwcm9taXNlIHRoYXQgZnVsZmlsbHMgb25jZSBhbGwgY2FuY2VsbGF0aW9uIHByb2NlZHVyZXMgZm9yIHRoZSBnaXZlbiB2YWx1ZXMgaGF2ZSBzZXR0bGVkLlxuICovXG5mdW5jdGlvbiBjYW5jZWxBbGwocGFyZW50OiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHZhbHVlczogYW55W10sIGNhdXNlPzogYW55KTogUHJvbWlzZTx2b2lkPiB7XG4gICAgY29uc3QgcmVzdWx0cyA9IFtdO1xuXG4gICAgZm9yIChjb25zdCB2YWx1ZSBvZiB2YWx1ZXMpIHtcbiAgICAgICAgbGV0IGNhbmNlbDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyO1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgaWYgKCFpc0NhbGxhYmxlKHZhbHVlLnRoZW4pKSB7IGNvbnRpbnVlOyB9XG4gICAgICAgICAgICBjYW5jZWwgPSB2YWx1ZS5jYW5jZWw7XG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUoY2FuY2VsKSkgeyBjb250aW51ZTsgfVxuICAgICAgICB9IGNhdGNoIHsgY29udGludWU7IH1cblxuICAgICAgICBsZXQgcmVzdWx0OiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICByZXN1bHQgPSBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcbiAgICAgICAgICAgIGNvbnRpbnVlO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFyZXN1bHQpIHsgY29udGludWU7IH1cbiAgICAgICAgcmVzdWx0cy5wdXNoKFxuICAgICAgICAgICAgKHJlc3VsdCBpbnN0YW5jZW9mIFByb21pc2UgID8gcmVzdWx0IDogUHJvbWlzZS5yZXNvbHZlKHJlc3VsdCkpLmNhdGNoKChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xuICAgICAgICAgICAgfSlcbiAgICAgICAgKTtcbiAgICB9XG5cbiAgICByZXR1cm4gUHJvbWlzZS5hbGwocmVzdWx0cykgYXMgYW55O1xufVxuXG4vKipcbiAqIFJldHVybnMgaXRzIGFyZ3VtZW50LlxuICovXG5mdW5jdGlvbiBpZGVudGl0eTxUPih4OiBUKTogVCB7XG4gICAgcmV0dXJuIHg7XG59XG5cbi8qKlxuICogVGhyb3dzIGl0cyBhcmd1bWVudC5cbiAqL1xuZnVuY3Rpb24gdGhyb3dlcihyZWFzb24/OiBhbnkpOiBuZXZlciB7XG4gICAgdGhyb3cgcmVhc29uO1xufVxuXG4vKipcbiAqIEF0dGVtcHRzIHZhcmlvdXMgc3RyYXRlZ2llcyB0byBjb252ZXJ0IGFuIGVycm9yIHRvIGEgc3RyaW5nLlxuICovXG5mdW5jdGlvbiBlcnJvck1lc3NhZ2UoZXJyOiBhbnkpOiBzdHJpbmcge1xuICAgIHRyeSB7XG4gICAgICAgIGlmIChlcnIgaW5zdGFuY2VvZiBFcnJvciB8fCB0eXBlb2YgZXJyICE9PSAnb2JqZWN0JyB8fCBlcnIudG9TdHJpbmcgIT09IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcpIHtcbiAgICAgICAgICAgIHJldHVybiBcIlwiICsgZXJyO1xuICAgICAgICB9XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIEpTT04uc3RyaW5naWZ5KGVycik7XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcuY2FsbChlcnIpO1xuICAgIH0gY2F0Y2gge31cblxuICAgIHJldHVybiBcIjxjb3VsZCBub3QgY29udmVydCBlcnJvciB0byBzdHJpbmc+XCI7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgY3VycmVudCBiYXJyaWVyIHByb21pc2UgZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLiBJZiBuZWNlc3NhcnksIGluaXRpYWxpc2VzIHRoZSBiYXJyaWVyLlxuICovXG5mdW5jdGlvbiBjdXJyZW50QmFycmllcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBsZXQgcHdyOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiA9IHByb21pc2VbYmFycmllclN5bV0gPz8ge307XG4gICAgaWYgKCEoJ3Byb21pc2UnIGluIHB3cikpIHtcbiAgICAgICAgT2JqZWN0LmFzc2lnbihwd3IsIHByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+KCkpO1xuICAgIH1cbiAgICBpZiAocHJvbWlzZVtiYXJyaWVyU3ltXSA9PSBudWxsKSB7XG4gICAgICAgIHB3ci5yZXNvbHZlISgpO1xuICAgICAgICBwcm9taXNlW2JhcnJpZXJTeW1dID0gcHdyO1xuICAgIH1cbiAgICByZXR1cm4gcHdyLnByb21pc2UhO1xufVxuXG4vLyBQb2x5ZmlsbCBQcm9taXNlLndpdGhSZXNvbHZlcnMuXG5sZXQgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBQcm9taXNlLndpdGhSZXNvbHZlcnM7XG5pZiAocHJvbWlzZVdpdGhSZXNvbHZlcnMgJiYgdHlwZW9mIHByb21pc2VXaXRoUmVzb2x2ZXJzID09PSAnZnVuY3Rpb24nKSB7XG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBwcm9taXNlV2l0aFJlc29sdmVycy5iaW5kKFByb21pc2UpO1xufSBlbHNlIHtcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IGZ1bmN0aW9uIDxUPigpOiBQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IFByb21pc2U8VD4oKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XG4gICAgICAgIHJldHVybiB7IHByb21pc2UsIHJlc29sdmUsIHJlamVjdCB9O1xuICAgIH1cbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2xpcGJvYXJkKTtcblxuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XG5jb25zdCBDbGlwYm9hcmRUZXh0ID0gMTtcblxuLyoqXG4gKiBTZXRzIHRoZSB0ZXh0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKlxuICogQHBhcmFtIHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cbiAqIEByZXR1cm4gQSBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2hlbiB0aGUgb3BlcmF0aW9uIGlzIHN1Y2Nlc3NmdWwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQ6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpOiBQcm9taXNlPHN0cmluZz4ge1xuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIEFueSBpcyBhIGR1bW15IGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBzaW1wbGUgb3IgdW5rbm93biB0eXBlcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEFueTxUID0gYW55Pihzb3VyY2U6IGFueSk6IFQge1xuICAgIHJldHVybiBzb3VyY2U7XG59XG5cbi8qKlxuICogQnl0ZVNsaWNlIGlzIGEgY3JlYXRpb24gZnVuY3Rpb24gdGhhdCByZXBsYWNlc1xuICogbnVsbCBzdHJpbmdzIHdpdGggZW1wdHkgc3RyaW5ncy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5dGVTbGljZShzb3VyY2U6IGFueSk6IHN0cmluZyB7XG4gICAgcmV0dXJuICgoc291cmNlID09IG51bGwpID8gXCJcIiA6IHNvdXJjZSk7XG59XG5cbi8qKlxuICogQXJyYXkgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcnJheVxuICogd2hvc2UgZWxlbWVudHMgYXJlIG9mIHRoYXQgdHlwZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEFycmF5PFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gVFtdIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybiBbXTtcbiAgICAgICAgfVxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNvdXJjZS5sZW5ndGg7IGkrKykge1xuICAgICAgICAgICAgc291cmNlW2ldID0gZWxlbWVudChzb3VyY2VbaV0pO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gb2JqZWN0XG4gKiB3aG9zZSBrZXlzIGFuZCB2YWx1ZXMgYXJlIG9mIHRob3NlIHR5cGVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTWFwPFYgPSBhbnk+KGtleTogKHNvdXJjZTogYW55KSA9PiBzdHJpbmcsIHZhbHVlOiAoc291cmNlOiBhbnkpID0+IFYpOiAoc291cmNlOiBhbnkpID0+IFJlY29yZDxzdHJpbmcsIFY+IHtcbiAgICBpZiAodmFsdWUgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IHt9IDogc291cmNlKTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm4ge307XG4gICAgICAgIH1cbiAgICAgICAgZm9yIChjb25zdCBrZXkgaW4gc291cmNlKSB7XG4gICAgICAgICAgICBzb3VyY2Vba2V5XSA9IHZhbHVlKHNvdXJjZVtrZXldKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogTnVsbGFibGUgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgbnVsbGFibGUgdmFsdWUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTnVsbGFibGU8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiAoVCB8IG51bGwpIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcbn1cblxuLyoqXG4gKiBTdHJ1Y3QgdGFrZXMgYW4gb2JqZWN0IG1hcHBpbmcgZmllbGQgbmFtZXMgdG8gY3JlYXRpb24gZnVuY3Rpb25zXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBzdHJ1Y3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQ6IFJlY29yZDxzdHJpbmcsIChzb3VyY2U6IGFueSkgPT4gYW55Pik6XG4gICAgPFUgZXh0ZW5kcyBSZWNvcmQ8c3RyaW5nLCBhbnk+ID0gYW55Pihzb3VyY2U6IGFueSkgPT4gVVxue1xuICAgIGxldCBhbGxBbnkgPSB0cnVlO1xuICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICBpZiAoY3JlYXRlRmllbGRbbmFtZV0gIT09IEFueSkge1xuICAgICAgICAgICAgYWxsQW55ID0gZmFsc2U7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoYWxsQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XG4gICAgICAgICAgICBpZiAobmFtZSBpbiBzb3VyY2UpIHtcbiAgICAgICAgICAgICAgICBzb3VyY2VbbmFtZV0gPSBjcmVhdGVGaWVsZFtuYW1lXShzb3VyY2VbbmFtZV0pO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuZXhwb3J0IGludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cbiAgICBXaWR0aDogbnVtYmVyO1xuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cbiAgICBIZWlnaHQ6IG51bWJlcjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBSZWN0IHtcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xuICAgIFg6IG51bWJlcjtcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xuICAgIFk6IG51bWJlcjtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuICovXG4gICAgV2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xuICAgIEhlaWdodDogbnVtYmVyO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFNjcmVlbiB7XG4gICAgLyoqIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLiAqL1xuICAgIElEOiBzdHJpbmc7XG4gICAgLyoqIEh1bWFuLXJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBOYW1lOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBzY2FsZSBmYWN0b3Igb2YgdGhlIHNjcmVlbiAoRFBJLzk2KS4gMSA9IHN0YW5kYXJkIERQSSwgMiA9IEhpRFBJIChSZXRpbmEpLCBldGMuICovXG4gICAgU2NhbGVGYWN0b3I6IG51bWJlcjtcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFg6IG51bWJlcjtcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFk6IG51bWJlcjtcbiAgICAvKiogQ29udGFpbnMgdGhlIHdpZHRoIGFuZCBoZWlnaHQgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBTaXplOiBTaXplO1xuICAgIC8qKiBDb250YWlucyB0aGUgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQuICovXG4gICAgQm91bmRzOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQgKGJlZm9yZSBzY2FsaW5nKS4gKi9cbiAgICBQaHlzaWNhbEJvdW5kczogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIGFyZWEgb2YgdGhlIHNjcmVlbiB0aGF0IGlzIGFjdHVhbGx5IHVzYWJsZSAoZXhjbHVkaW5nIHRhc2tiYXIgYW5kIG90aGVyIHN5c3RlbSBVSSkuICovXG4gICAgV29ya0FyZWE6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBXb3JrQXJlYSBvZiB0aGUgc2NyZWVuIChiZWZvcmUgc2NhbGluZykuICovXG4gICAgUGh5c2ljYWxXb3JrQXJlYTogUmVjdDtcbiAgICAvKiogVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXG4gICAgSXNQcmltYXJ5OiBib29sZWFuO1xuICAgIC8qKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBSb3RhdGlvbjogbnVtYmVyO1xufVxuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlNjcmVlbnMpO1xuXG5jb25zdCBnZXRBbGwgPSAwO1xuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XG5jb25zdCBnZXRDdXJyZW50ID0gMjtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIGFycmF5IG9mIFNjcmVlbiBvYmplY3RzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCk6IFByb21pc2U8U2NyZWVuW10+IHtcbiAgICByZXR1cm4gY2FsbChnZXRBbGwpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRQcmltYXJ5KTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgIHJldHVybiBjYWxsKGdldEN1cnJlbnQpO1xufVxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUE7QUFBQTtBQUFBO0FBQUE7OztBQzZCQSxJQUFNLGNBQ0Y7QUFFRyxTQUFTLE9BQU8sT0FBZSxJQUFZO0FBQzlDLE1BQUksS0FBSztBQUVULE1BQUksSUFBSSxPQUFPO0FBQ2YsU0FBTyxLQUFLO0FBRVIsVUFBTSxZQUFhLEtBQUssT0FBTyxJQUFJLEtBQU0sQ0FBQztBQUFBLEVBQzlDO0FBQ0EsU0FBTztBQUNYOzs7QUM3QkEsSUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBR3JDLElBQU0sY0FBYyxPQUFPLE9BQU87QUFBQSxFQUNyQyxNQUFNO0FBQUEsRUFDTixXQUFXO0FBQUEsRUFDWCxhQUFhO0FBQUEsRUFDYixRQUFRO0FBQUEsRUFDUixhQUFhO0FBQUEsRUFDYixRQUFRO0FBQUEsRUFDUixRQUFRO0FBQUEsRUFDUixTQUFTO0FBQUEsRUFDVCxRQUFRO0FBQUEsRUFDUixTQUFTO0FBQUEsRUFDVCxZQUFZO0FBQ2hCLENBQUM7QUFDTSxJQUFJLFdBQVcsT0FBTztBQVN0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RyxDQUFDO0FBQUEsSUFDbkM7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUNKO0FBTU8sU0FBUyxTQUFrQixTQUEwRDtBQUN4RixNQUFJLFlBQVksS0FBSztBQUNqQixXQUFPO0FBQUEsRUFDWDtBQUVBLFNBQU8sQ0FBQyxXQUFZLFdBQVcsT0FBTyxPQUFPLFFBQVEsTUFBTTtBQUMvRDtBQU1PLFNBQVMsT0FBTyxhQUV2QjtBQUNJLE1BQUksU0FBUztBQUNiLGFBQVcsUUFBUSxhQUFhO0FBQzVCLFFBQUksWUFBWSxJQUFJLE1BQU0sS0FBSztBQUMzQixlQUFTO0FBQ1Q7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUNBLE1BQUksUUFBUTtBQUNSLFdBQU87QUFBQSxFQUNYO0FBRUEsU0FBTyxDQUFDLFdBQVc7QUFDZixlQUFXLFFBQVEsYUFBYTtBQUM1QixVQUFJLFFBQVEsUUFBUTtBQUNoQixlQUFPLElBQUksSUFBSSxZQUFZLElBQUksRUFBRSxPQUFPLElBQUksQ0FBQztBQUFBLE1BQ2pEO0FBQUEsSUFDSjtBQUNBLFdBQU87QUFBQSxFQUNYO0FBQ0o7OztBQ3pHQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUF3REEsSUFBTUMsUUFBTyxpQkFBaUIsWUFBWSxPQUFPO0FBRWpELElBQU0sU0FBUztBQUNmLElBQU0sYUFBYTtBQUNuQixJQUFNLGFBQWE7QUFPWixTQUFTLFNBQTRCO0FBQ3hDLFNBQU9BLE1BQUssTUFBTTtBQUN0QjtBQU9PLFNBQVMsYUFBOEI7QUFDMUMsU0FBT0EsTUFBSyxVQUFVO0FBQzFCO0FBT08sU0FBUyxhQUE4QjtBQUMxQyxTQUFPQSxNQUFLLFVBQVU7QUFDMUI7OztBdEI1RUEsT0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBNENsQyxPQUFPLE9BQU8sU0FBZ0I7QUFDdkIsT0FBTyxxQkFBcUI7IiwKICAibmFtZXMiOiBbIl9hIiwgIkVycm9yIiwgImNhbGwiLCAiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJfYSIsICJyZXNpemFibGUiLCAiY2FsbCIsICJfYSIsICJjYWxsIiwgImNhbGwiLCAiSGlkZU1ldGhvZCIsICJTaG93TWV0aG9kIiwgImlzRG9jdW1lbnREb3RBbGwiLCAiX2EiLCAicmVhc29uIiwgInZhbHVlIiwgImNhbGwiLCAiZ2V0QW5kRGVsZXRlUmVzcG9uc2UiLCAiZ2VuZXJhdGVJRCIsICJjYWxsIiwgIkFycmF5IiwgIk1hcCIsICJBcnJheSIsICJNYXAiLCAia2V5IiwgImNhbGwiXQp9Cg==
