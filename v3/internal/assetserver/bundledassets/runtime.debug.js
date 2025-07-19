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
    WindowToggleFrameless: "common:WindowToggleFrameless",
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
var ToggleFramelessMethod = 40;
var UnFullscreenMethod = 41;
var UnMaximiseMethod = 42;
var UnMinimiseMethod = 43;
var WidthMethod = 44;
var ZoomMethod = 45;
var ZoomInMethod = 46;
var ZoomOutMethod = 47;
var ZoomResetMethod = 48;
var WindowDropZoneDropped = 49;
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
   * Toggles the window between frameless and normal.
   */
  ToggleFrameless() {
    return this[callerSym](ToggleFramelessMethod);
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
    const elementDetails = {
      id: dropzoneTarget.id,
      classList: Array.from(dropzoneTarget.classList),
      attributes: {}
    };
    for (let i = 0; i < dropzoneTarget.attributes.length; i++) {
      const attr = dropzoneTarget.attributes[i];
      elementDetails.attributes[attr.name] = attr.value;
    }
    const payload = {
      filenames,
      x,
      y,
      elementDetails
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHBsaWNhdGlvbi50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYW5jZWxsYWJsZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG5cbmltcG9ydCBcIi4vY29udGV4dG1lbnUuanNcIjtcbmltcG9ydCBcIi4vZHJhZy5qc1wiO1xuXG4vLyBSZS1leHBvcnQgcHVibGljIEFQSVxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSBcIi4vYXBwbGljYXRpb24uanNcIjtcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xuaW1wb3J0ICogYXMgQ2FsbCBmcm9tIFwiLi9jYWxscy5qc1wiO1xuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZC5qc1wiO1xuaW1wb3J0ICogYXMgQ3JlYXRlIGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xuaW1wb3J0ICogYXMgRGlhbG9ncyBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XG5pbXBvcnQgKiBhcyBFdmVudHMgZnJvbSBcIi4vZXZlbnRzLmpzXCI7XG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi9mbGFncy5qc1wiO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuaW1wb3J0ICogYXMgV01MIGZyb20gXCIuL3dtbC5qc1wiO1xuXG5leHBvcnQge1xuICAgIEFwcGxpY2F0aW9uLFxuICAgIEJyb3dzZXIsXG4gICAgQ2FsbCxcbiAgICBDbGlwYm9hcmQsXG4gICAgRGlhbG9ncyxcbiAgICBFdmVudHMsXG4gICAgRmxhZ3MsXG4gICAgU2NyZWVucyxcbiAgICBTeXN0ZW0sXG4gICAgV2luZG93LFxuICAgIFdNTFxufTtcblxuLyoqXG4gKiBBbiBpbnRlcm5hbCB1dGlsaXR5IGNvbnN1bWVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cbiAqXG4gKiBAaWdub3JlXG4gKiBAaW50ZXJuYWxcbiAqL1xuZXhwb3J0IHsgQ3JlYXRlIH07XG5cbmV4cG9ydCAqIGZyb20gXCIuL2NhbmNlbGxhYmxlLmpzXCI7XG5cbi8vIE5vdGlmeSBiYWNrZW5kXG53aW5kb3cuX3dhaWxzLmludm9rZSA9IFN5c3RlbS5pbnZva2U7XG5TeXN0ZW0uaW52b2tlKFwid2FpbHM6cnVudGltZTpyZWFkeVwiKTtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCB7IFF1ZXN0aW9uIH0gZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0IHsgRW1pdCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXG4gKiBAcGFyYW0gW2RhdGE9bnVsbF0gLSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cbiAqL1xuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XG4gICAgdm9pZCBFbWl0KGV2ZW50TmFtZSwgZGF0YSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBzcGVjaWZpZWQgd2luZG93LlxuICpcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqL1xuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lOiBzdHJpbmcsIG1ldGhvZE5hbWU6IHN0cmluZykge1xuICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IFdpbmRvdy5HZXQod2luZG93TmFtZSk7XG4gICAgY29uc3QgbWV0aG9kID0gKHRhcmdldFdpbmRvdyBhcyBhbnkpW21ldGhvZE5hbWVdO1xuXG4gICAgaWYgKHR5cGVvZiBtZXRob2QgIT09IFwiZnVuY3Rpb25cIikge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBXaW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JyBub3QgZm91bmRgKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICAgIG1ldGhvZC5jYWxsKHRhcmdldFdpbmRvdyk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmVzcG9uZHMgdG8gYSB0cmlnZ2VyaW5nIGV2ZW50IGJ5IHJ1bm5pbmcgYXBwcm9wcmlhdGUgV01MIGFjdGlvbnMgZm9yIHRoZSBjdXJyZW50IHRhcmdldC5cbiAqL1xuZnVuY3Rpb24gb25XTUxUcmlnZ2VyZWQoZXY6IEV2ZW50KTogdm9pZCB7XG4gICAgY29uc3QgZWxlbWVudCA9IGV2LmN1cnJlbnRUYXJnZXQgYXMgRWxlbWVudDtcblxuICAgIGZ1bmN0aW9uIHJ1bkVmZmVjdChjaG9pY2UgPSBcIlllc1wiKSB7XG4gICAgICAgIGlmIChjaG9pY2UgIT09IFwiWWVzXCIpXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1ldmVudCcpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCB1cmwgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLW9wZW51cmwnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtb3BlbnVybCcpO1xuXG4gICAgICAgIGlmIChldmVudFR5cGUgIT09IG51bGwpXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICBpZiAodXJsICE9PSBudWxsKVxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgfVxuXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG5cbiAgICBpZiAoY29uZmlybSkge1xuICAgICAgICBRdWVzdGlvbih7XG4gICAgICAgICAgICBUaXRsZTogXCJDb25maXJtXCIsXG4gICAgICAgICAgICBNZXNzYWdlOiBjb25maXJtLFxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxuICAgICAgICAgICAgQnV0dG9uczogW1xuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiWWVzXCIgfSxcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XG4gICAgICAgICAgICBdXG4gICAgICAgIH0pLnRoZW4ocnVuRWZmZWN0KTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBydW5FZmZlY3QoKTtcbiAgICB9XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBjb250cm9sbGVyU3ltID0gU3ltYm9sKFwiY29udHJvbGxlclwiKTtcbmNvbnN0IHRyaWdnZXJNYXBTeW0gPSBTeW1ib2woXCJ0cmlnZ2VyTWFwXCIpO1xuY29uc3QgZWxlbWVudENvdW50U3ltID0gU3ltYm9sKFwiZWxlbWVudENvdW50XCIpO1xuXG4vKipcbiAqIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IGRvZXMgbm90IGFjdHVhbGx5IHJlbWVtYmVyIGFjdGl2ZSBldmVudCBsaXN0ZW5lcnM6IGluc3RlYWRcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXG4gKi9cbmNsYXNzIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICBbY29udHJvbGxlclN5bV06IEFib3J0Q29udHJvbGxlcjtcblxuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgYW4gb3B0aW9ucyBvYmplY3QgZm9yIGFkZEV2ZW50TGlzdGVuZXIgdGhhdCB0aWVzIHRoZSBsaXN0ZW5lclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xuICAgICAqL1xuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XG4gICAgICAgIHJldHVybiB7IHNpZ25hbDogdGhpc1tjb250cm9sbGVyU3ltXS5zaWduYWwgfTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycyBhbmQgcmVzZXRzIHRoZSByZWdpc3RyeS5cbiAgICAgKi9cbiAgICByZXNldCgpOiB2b2lkIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXS5hYm9ydCgpO1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBXZWFrTWFwUmVnaXN0cnkgbWFwcyBhY3RpdmUgdHJpZ2dlciBldmVudHMgdG8gZWFjaCBET00gZWxlbWVudCB0aHJvdWdoIGEgV2Vha01hcC5cbiAqIFRoaXMgZW5zdXJlcyB0aGF0IHRoZSBtYXBwaW5nIHJlbWFpbnMgcHJpdmF0ZSB0byB0aGlzIG1vZHVsZSwgd2hpbGUgc3RpbGwgYWxsb3dpbmcgZ2FyYmFnZVxuICogY29sbGVjdGlvbiBvZiB0aGUgaW52b2x2ZWQgZWxlbWVudHMuXG4gKi9cbmNsYXNzIFdlYWtNYXBSZWdpc3RyeSB7XG4gICAgLyoqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy4gKi9cbiAgICBbdHJpZ2dlck1hcFN5bV06IFdlYWtNYXA8RWxlbWVudCwgc3RyaW5nW10+O1xuICAgIC8qKiBDb3VudHMgdGhlIG51bWJlciBvZiBlbGVtZW50cyB3aXRoIGFjdGl2ZSBXTUwgdHJpZ2dlcnMuICovXG4gICAgW2VsZW1lbnRDb3VudFN5bV06IG51bWJlcjtcblxuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIGFjdGl2ZSB0cmlnZ2VycyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50LlxuICAgICAqXG4gICAgICogQHBhcmFtIGVsZW1lbnQgLSBBbiBIVE1MIGVsZW1lbnRcbiAgICAgKiBAcGFyYW0gdHJpZ2dlcnMgLSBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRcbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xuICAgICAgICBpZiAoIXRoaXNbdHJpZ2dlck1hcFN5bV0uaGFzKGVsZW1lbnQpKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXSsrOyB9XG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICAgICAgcmV0dXJuIHt9O1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzLlxuICAgICAqL1xuICAgIHJlc2V0KCk6IHZvaWQge1xuICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnKicpKSB7XG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXG4gICAgICAgICAgICAgICAgYnJlYWs7XG5cbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwU3ltXS5nZXQoZWxlbWVudCk7XG4gICAgICAgICAgICBpZiAodHJpZ2dlcnMgIT0gbnVsbCkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0tLTsgfVxuXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXG4gICAgICAgICAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkKTtcbiAgICAgICAgfVxuXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0gPSBuZXcgV2Vha01hcCgpO1xuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xuICAgIH1cbn1cblxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xuXG4vKipcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAqL1xuZnVuY3Rpb24gYWRkV01MTGlzdGVuZXJzKGVsZW1lbnQ6IEVsZW1lbnQpOiB2b2lkIHtcbiAgICBjb25zdCB0cmlnZ2VyUmVnRXhwID0gL1xcUysvZztcbiAgICBjb25zdCB0cmlnZ2VyQXR0ciA9IChlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIik7XG4gICAgY29uc3QgdHJpZ2dlcnM6IHN0cmluZ1tdID0gW107XG5cbiAgICBsZXQgbWF0Y2g7XG4gICAgd2hpbGUgKChtYXRjaCA9IHRyaWdnZXJSZWdFeHAuZXhlYyh0cmlnZ2VyQXR0cikpICE9PSBudWxsKVxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcblxuICAgIGNvbnN0IG9wdGlvbnMgPSB0cmlnZ2VyUmVnaXN0cnkuc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2NoZWR1bGVzIGFuIGF1dG9tYXRpYyByZWxvYWQgb2YgV01MIHRvIGJlIHBlcmZvcm1lZCBhcyBzb29uIGFzIHRoZSBkb2N1bWVudCBpcyBmdWxseSBsb2FkZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbmFibGUoKTogdm9pZCB7XG4gICAgd2hlblJlYWR5KFJlbG9hZCk7XG59XG5cbi8qKlxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCk6IHZvaWQge1xuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xuICAgIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0sIFt3bWwtd2luZG93XSwgW3dtbC1vcGVudXJsXSwgW2RhdGEtd21sLWV2ZW50XSwgW2RhdGEtd21sLXdpbmRvd10sIFtkYXRhLXdtbC1vcGVudXJsXScpLmZvckVhY2goYWRkV01MTGlzdGVuZXJzKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xuXG5jb25zdCBCcm93c2VyT3BlblVSTCA9IDA7XG5cbi8qKlxuICogT3BlbiBhIGJyb3dzZXIgd2luZG93IHRvIHRoZSBnaXZlbiBVUkwuXG4gKlxuICogQHBhcmFtIHVybCAtIFRoZSBVUkwgdG8gb3BlblxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmw6IHN0cmluZyB8IFVSTCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsOiB1cmwudG9TdHJpbmcoKX0pO1xufVxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxuXG4vLyBUaGUgTUlUIExpY2Vuc2UgKE1JVClcbi8vXG4vLyBDb3B5cmlnaHQgMjAxNyBBbmRyZXkgU2l0bmlrIDxhbmRyZXlAc2l0bmlrLnJ1PlxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2Zcbi8vIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWwgaW5cbi8vIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMgdG9cbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXG4vLyB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sXG4vLyAgICAgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6XG4vL1xuLy8gICAgIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cbi8vXG4vLyAgICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksIEZJVE5FU1Ncbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxuLy8gQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLCBXSEVUSEVSXG4vLyBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sIE9VVCBPRiBPUiBJTlxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cblxuLy8gVGhpcyBhbHBoYWJldCB1c2VzIGBBLVphLXowLTlfLWAgc3ltYm9scy5cbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cbi8vIFJlZmVyZW5jZXMgdG8gdGhlIHNhbWUgZmlsZSAod29ya3MgYm90aCBmb3IgZ3ppcCBhbmQgYnJvdGxpKTpcbi8vIGAndXNlYCwgYGFuZG9tYCwgYW5kIGByaWN0J2Bcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XG4vLyBgLTI2VGAsIGAxOTgzYCwgYDQwcHhgLCBgNzVweGAsIGBidXNoYCwgYGphY2tgLCBgbWluZGAsIGB2ZXJ5YCwgYW5kIGB3b2xmYFxuY29uc3QgdXJsQWxwaGFiZXQgPVxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuXG5leHBvcnQgZnVuY3Rpb24gbmFub2lkKHNpemU6IG51bWJlciA9IDIxKTogc3RyaW5nIHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIC8vIEEgY29tcGFjdCBhbHRlcm5hdGl2ZSBmb3IgYGZvciAodmFyIGkgPSAwOyBpIDwgc3RlcDsgaSsrKWAuXG4gICAgbGV0IGkgPSBzaXplIHwgMFxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgICAgLy8gYHwgMGAgaXMgbW9yZSBjb21wYWN0IGFuZCBmYXN0ZXIgdGhhbiBgTWF0aC5mbG9vcigpYC5cbiAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuLy8gT2JqZWN0IE5hbWVzXG5leHBvcnQgY29uc3Qgb2JqZWN0TmFtZXMgPSBPYmplY3QuZnJlZXplKHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxuICAgIENhbmNlbENhbGw6IDEwLFxufSk7XG5leHBvcnQgbGV0IGNsaWVudElkID0gbmFub2lkKCk7XG5cbi8qKlxuICogQ3JlYXRlcyBhIG5ldyBydW50aW1lIGNhbGxlciB3aXRoIHNwZWNpZmllZCBJRC5cbiAqXG4gKiBAcGFyYW0gb2JqZWN0IC0gVGhlIG9iamVjdCB0byBpbnZva2UgdGhlIG1ldGhvZCBvbi5cbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAqIEByZXR1cm4gVGhlIG5ldyBydW50aW1lIGNhbGxlciBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0OiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZyA9ICcnKSB7XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2Q6IG51bWJlciwgYXJnczogYW55ID0gbnVsbCkge1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0LCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xuICAgIH07XG59XG5cbmFzeW5jIGZ1bmN0aW9uIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55PiB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQudG9TdHJpbmcoKSk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kLnRvU3RyaW5nKCkpO1xuICAgIGlmIChhcmdzKSB7IHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7IH1cblxuICAgIGxldCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+ID0ge1xuICAgICAgICBbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXTogY2xpZW50SWRcbiAgICB9XG4gICAgaWYgKHdpbmRvd05hbWUpIHtcbiAgICAgICAgaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xuICAgIH1cblxuICAgIGxldCByZXNwb25zZSA9IGF3YWl0IGZldGNoKHVybCwgeyBoZWFkZXJzIH0pO1xuICAgIGlmICghcmVzcG9uc2Uub2spIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XG4gICAgfVxuXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XG5cbi8vIHNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbndpbmRvdy5fd2FpbHMuZGlhbG9nRXJyb3JDYWxsYmFjayA9IGRpYWxvZ0Vycm9yQ2FsbGJhY2s7XG53aW5kb3cuX3dhaWxzLmRpYWxvZ1Jlc3VsdENhbGxiYWNrID0gZGlhbG9nUmVzdWx0Q2FsbGJhY2s7XG5cbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8UHJvbWlzZVdpdGhSZXNvbHZlcnM8YW55PiwgXCJwcm9taXNlXCI+O1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5EaWFsb2cpO1xuY29uc3QgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcDxzdHJpbmcsIFByb21pc2VSZXNvbHZlcnM+KCk7XG5cbi8vIERlZmluZSBjb25zdGFudHMgZnJvbSB0aGUgYG1ldGhvZHNgIG9iamVjdCBpbiBUaXRsZSBDYXNlXG5jb25zdCBEaWFsb2dJbmZvID0gMDtcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xuY29uc3QgRGlhbG9nRXJyb3IgPSAyO1xuY29uc3QgRGlhbG9nUXVlc3Rpb24gPSAzO1xuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xuY29uc3QgRGlhbG9nU2F2ZUZpbGUgPSA1O1xuXG5leHBvcnQgaW50ZXJmYWNlIE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLiAqL1xuICAgIEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFNhdmVGaWxlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuICovXG4gICAgRmlsZW5hbWU/OiBzdHJpbmc7XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIE1lc3NhZ2VEaWFsb2dPcHRpb25zIHtcbiAgICAvKiogVGhlIHRpdGxlIG9mIHRoZSBkaWFsb2cgd2luZG93LiAqL1xuICAgIFRpdGxlPzogc3RyaW5nO1xuICAgIC8qKiBUaGUgbWFpbiBtZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBNZXNzYWdlPzogc3RyaW5nO1xuICAgIC8qKiBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgQnV0dG9ucz86IEJ1dHRvbltdO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuICovXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEJ1dHRvbiB7XG4gICAgLyoqIFRleHQgdGhhdCBhcHBlYXJzIHdpdGhpbiB0aGUgYnV0dG9uLiAqL1xuICAgIExhYmVsPzogc3RyaW5nO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLiAqL1xuICAgIElzQ2FuY2VsPzogYm9vbGVhbjtcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBiZSB0aGUgZGVmYXVsdCBhY3Rpb24gd2hlbiB0aGUgdXNlciBwcmVzc2VzIGVudGVyLiAqL1xuICAgIElzRGVmYXVsdD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgRmlsZUZpbHRlciB7XG4gICAgLyoqIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLiAqL1xuICAgIERpc3BsYXlOYW1lPzogc3RyaW5nO1xuICAgIC8qKiBQYXR0ZXJuIHRvIG1hdGNoIGZvciB0aGUgZmlsdGVyLCBlLmcuIFwiKi50eHQ7Ki5tZFwiIGZvciB0ZXh0IG1hcmtkb3duIGZpbGVzLiAqL1xuICAgIFBhdHRlcm4/OiBzdHJpbmc7XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgcmVzdWx0IG9mIGEgZGlhbG9nIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cbiAqIEBwYXJhbSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxuICogQHBhcmFtIGlzSlNPTiAtIEluZGljYXRlcyB3aGV0aGVyIHRoZSBkYXRhIGlzIEpTT04gb3Igbm90LlxuICovXG5mdW5jdGlvbiBkaWFsb2dSZXN1bHRDYWxsYmFjayhpZDogc3RyaW5nLCBkYXRhOiBzdHJpbmcsIGlzSlNPTjogYm9vbGVhbik6IHZvaWQge1xuICAgIGxldCByZXNvbHZlcnMgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XG4gICAgaWYgKCFyZXNvbHZlcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHJlc29sdmVycy5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGNhdGNoIChlcnI6IGFueSkge1xuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgVHlwZUVycm9yKFwiY291bGQgbm90IHBhcnNlIHJlc3VsdDogXCIgKyBlcnIubWVzc2FnZSwgeyBjYXVzZTogZXJyIH0pKTtcbiAgICAgICAgfVxuICAgIH0gZWxzZSB7XG4gICAgICAgIHJlc29sdmVycy5yZXNvbHZlKGRhdGEpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBlcnJvciBmcm9tIGEgZGlhbG9nIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSBwcm9taXNlIGhhbmRsZXIuXG4gKiBAcGFyYW0gbWVzc2FnZSAtIEFuIGVycm9yIG1lc3NhZ2UuXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQ6IHN0cmluZywgbWVzc2FnZTogc3RyaW5nKTogdm9pZCB7XG4gICAgZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpPy5yZWplY3QobmV3IHdpbmRvdy5FcnJvcihtZXNzYWdlKSk7XG59XG5cbi8qKlxuICogUmV0cmlldmVzIGFuZCByZW1vdmVzIHRoZSByZXNwb25zZSBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElEIGZyb20gdGhlIGRpYWxvZ1Jlc3BvbnNlcyBtYXAuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIElEIG9mIHRoZSByZXNwb25zZSB0byBiZSByZXRyaWV2ZWQgYW5kIHJlbW92ZWQuXG4gKiBAcmV0dXJucyBUaGUgcmVzcG9uc2Ugb2JqZWN0IGFzc29jaWF0ZWQgd2l0aCB0aGUgZ2l2ZW4gSUQsIGlmIGFueS5cbiAqL1xuZnVuY3Rpb24gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQ6IHN0cmluZyk6IFByb21pc2VSZXNvbHZlcnMgfCB1bmRlZmluZWQge1xuICAgIGNvbnN0IHJlc3BvbnNlID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgcmV0dXJuIHJlc3BvbnNlO1xufVxuXG4vKipcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBJRCB1c2luZyB0aGUgbmFub2lkIGxpYnJhcnkuXG4gKlxuICogQHJldHVybnMgQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgZGlhbG9nUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogUHJlc2VudHMgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqXG4gKiBAcGFyYW0gdHlwZSAtIERpYWxvZyB0eXBlLlxuICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIGZvciB0aGUgZGlhbG9nLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nLlxuICovXG5mdW5jdGlvbiBkaWFsb2codHlwZTogbnVtYmVyLCBvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyB8IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB8IFNhdmVGaWxlRGlhbG9nT3B0aW9ucyA9IHt9KTogUHJvbWlzZTxhbnk+IHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmUsIHJlamVjdCB9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBPYmplY3QuYXNzaWduKHsgXCJkaWFsb2ctaWRcIjogaWQgfSwgb3B0aW9ucykpLmNhdGNoKChlcnI6IGFueSkgPT4ge1xuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgICAgICByZWplY3QoZXJyKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cbi8qKlxuICogUHJlc2VudHMgYW4gaW5mbyBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgd2FybmluZyBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhbiBlcnJvciBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRXJyb3Iob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgcXVlc3Rpb24gZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIG9uZSBvciBtb3JlIGZpbGVzIHRvIG9wZW4uXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcywgb3IgYSBibGFuayBzdHJpbmcvZW1wdHkgbGlzdCBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbjogdHJ1ZSB9KTogUHJvbWlzZTxzdHJpbmdbXT47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGZhbHNlIHwgdW5kZWZpbmVkIH0pOiBQcm9taXNlPHN0cmluZz47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT4geyByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKSA/PyBbXTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgYSBmaWxlIHRvIHNhdmUuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUsIG9yIGEgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zOiBTYXZlRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTsgfVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IGV2ZW50TGlzdGVuZXJzLCBMaXN0ZW5lciwgbGlzdGVuZXJPZmYgfSBmcm9tIFwiLi9saXN0ZW5lci5qc1wiO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRXZlbnRzKTtcbmNvbnN0IEVtaXRNZXRob2QgPSAwO1xuXG5leHBvcnQgeyBUeXBlcyB9IGZyb20gXCIuL2V2ZW50X3R5cGVzLmpzXCI7XG5cbi8qKlxuICogVGhlIHR5cGUgb2YgaGFuZGxlcnMgZm9yIGEgZ2l2ZW4gZXZlbnQuXG4gKi9cbmV4cG9ydCB0eXBlIENhbGxiYWNrID0gKGV2OiBXYWlsc0V2ZW50KSA9PiB2b2lkO1xuXG4vKipcbiAqIFJlcHJlc2VudHMgYSBzeXN0ZW0gZXZlbnQgb3IgYSBjdXN0b20gZXZlbnQgZW1pdHRlZCB0aHJvdWdoIHdhaWxzLXByb3ZpZGVkIGZhY2lsaXRpZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcbiAgICAvKipcbiAgICAgKiBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXG4gICAgICovXG4gICAgbmFtZTogc3RyaW5nO1xuXG4gICAgLyoqXG4gICAgICogT3B0aW9uYWwgZGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGVtaXR0ZWQgZXZlbnQuXG4gICAgICovXG4gICAgZGF0YTogYW55O1xuXG4gICAgLyoqXG4gICAgICogTmFtZSBvZiB0aGUgb3JpZ2luYXRpbmcgd2luZG93LiBPbWl0dGVkIGZvciBhcHBsaWNhdGlvbiBldmVudHMuXG4gICAgICogV2lsbCBiZSBvdmVycmlkZGVuIGlmIHNldCBtYW51YWxseS5cbiAgICAgKi9cbiAgICBzZW5kZXI/OiBzdHJpbmc7XG5cbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudDogYW55KSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAoIWxpc3RlbmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgbGV0IHdhaWxzRXZlbnQgPSBuZXcgV2FpbHNFdmVudChldmVudC5uYW1lLCBldmVudC5kYXRhKTtcbiAgICBpZiAoJ3NlbmRlcicgaW4gZXZlbnQpIHtcbiAgICAgICAgd2FpbHNFdmVudC5zZW5kZXIgPSBldmVudC5zZW5kZXI7XG4gICAgfVxuXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsaXN0ZW5lciA9PiAhbGlzdGVuZXIuZGlzcGF0Y2god2FpbHNFdmVudCkpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcGFyYW0gbWF4Q2FsbGJhY2tzIC0gVGhlIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRoZSBjYWxsYmFjayBjYW4gYmUgY2FsbGVkIGZvciB0aGUgZXZlbnQuIE9uY2UgdGhlIG1heGltdW0gbnVtYmVyIGlzIHJlYWNoZWQsIHRoZSBjYWxsYmFjayB3aWxsIG5vIGxvbmdlciBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IENhbGxiYWNrLCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IENhbGxiYWNrKTogKCkgPT4gdm9pZCB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZTogc3RyaW5nLCBjYWxsYmFjazogQ2FsbGJhY2spOiAoKSA9PiB2b2lkIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuLyoqXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lcyAtIFRoZSBuYW1lIG9mIHRoZSBldmVudHMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoLi4uZXZlbnROYW1lczogW3N0cmluZywgLi4uc3RyaW5nW11dKTogdm9pZCB7XG4gICAgZXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKSk7XG59XG5cbi8qKlxuICogUmVtb3ZlcyBhbGwgZXZlbnQgbGlzdGVuZXJzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCk6IHZvaWQge1xuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XG59XG5cbi8qKlxuICogRW1pdHMgYW4gZXZlbnQgdXNpbmcgdGhlIG5hbWUgYW5kIGRhdGEuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgd2lsbCBiZSBmdWxmaWxsZWQgb25jZSB0aGUgZXZlbnQgaGFzIGJlZW4gZW1pdHRlZC5cbiAqIEBwYXJhbSBuYW1lIC0gdGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIGVtaXQuXG4gKiBAcGFyYW0gZGF0YSAtIHRoZSBkYXRhIHRvIGJlIHNlbnQgd2l0aCB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbWl0KG5hbWU6IHN0cmluZywgZGF0YT86IGFueSk6IFByb21pc2U8dm9pZD4ge1xuICAgIGxldCBldmVudDogV2FpbHNFdmVudDtcblxuICAgIGlmICh0eXBlb2YgbmFtZSA9PT0gJ29iamVjdCcgJiYgbmFtZSAhPT0gbnVsbCAmJiAnbmFtZScgaW4gbmFtZSAmJiAnZGF0YScgaW4gbmFtZSkge1xuICAgICAgICAvLyBJZiBuYW1lIGlzIGFuIG9iamVjdCB3aXRoIGEgbmFtZSBwcm9wZXJ0eSwgdXNlIGl0IGRpcmVjdGx5XG4gICAgICAgIGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQobmFtZVsnbmFtZSddLCBuYW1lWydkYXRhJ10pO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIC8vIE90aGVyd2lzZSB1c2UgdGhlIHN0YW5kYXJkIHBhcmFtZXRlcnNcbiAgICAgICAgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChuYW1lIGFzIHN0cmluZywgZGF0YSk7XG4gICAgfVxuXG4gICAgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgZXZlbnQpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFRoZSBmb2xsb3dpbmcgdXRpbGl0aWVzIGhhdmUgYmVlbiBmYWN0b3JlZCBvdXQgb2YgLi9ldmVudHMudHNcbi8vIGZvciB0ZXN0aW5nIHB1cnBvc2VzLlxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwPHN0cmluZywgTGlzdGVuZXJbXT4oKTtcblxuZXhwb3J0IGNsYXNzIExpc3RlbmVyIHtcbiAgICBldmVudE5hbWU6IHN0cmluZztcbiAgICBjYWxsYmFjazogKGRhdGE6IGFueSkgPT4gdm9pZDtcbiAgICBtYXhDYWxsYmFja3M6IG51bWJlcjtcblxuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZTogc3RyaW5nLCBjYWxsYmFjazogKGRhdGE6IGFueSkgPT4gdm9pZCwgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIHRoaXMuY2FsbGJhY2sgPSBjYWxsYmFjaztcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgfVxuXG4gICAgZGlzcGF0Y2goZGF0YTogYW55KTogYm9vbGVhbiB7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICB0aGlzLmNhbGxiYWNrKGRhdGEpO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZXJyKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcbiAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyOiBMaXN0ZW5lcik6IHZvaWQge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQobGlzdGVuZXIuZXZlbnROYW1lKTtcbiAgICBpZiAoIWxpc3RlbmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUobGlzdGVuZXIuZXZlbnROYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQobGlzdGVuZXIuZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gQ3luaHlyY2h3eWQgeSBmZmVpbCBob24geW4gYXd0b21hdGlnLiBQRUlESVdDSCBcdTAwQzIgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHR9KSxcblx0TGludXg6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibGludXg6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcblx0XHRXaW5kb3dMb2FkQ2hhbmdlZDogXCJsaW51eDpXaW5kb3dMb2FkQ2hhbmdlZFwiLFxuXHR9KSxcblx0Q29tbW9uOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbk9wZW5lZFdpdGhGaWxlOiBcImNvbW1vbjpBcHBsaWNhdGlvbk9wZW5lZFdpdGhGaWxlXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcImNvbW1vbjpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiY29tbW9uOlRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwiY29tbW9uOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImNvbW1vbjpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImNvbW1vbjpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcImNvbW1vbjpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcblx0XHRXaW5kb3dGb2N1czogXCJjb21tb246V2luZG93Rm9jdXNcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJjb21tb246V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJjb21tb246V2luZG93TG9zdEZvY3VzXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1J1bnRpbWVSZWFkeTogXCJjb21tb246V2luZG93UnVudGltZVJlYWR5XCIsXG5cdFx0V2luZG93U2hvdzogXCJjb21tb246V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcImNvbW1vbjpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcblx0XHRXaW5kb3dEcm9wWm9uZUZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RHJvcFpvbmVGaWxlc0Ryb3BwZWRcIixcblx0fSksXG59KTtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxuICpcbiAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIG1lc3NhZ2UgdG8gYmUgbG9nZ2VkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZTogYW55KSB7XG4gICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgY29uc29sZS5sb2coXG4gICAgICAgICclYyB3YWlsczMgJWMgJyArIG1lc3NhZ2UgKyAnICcsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjYWEwMDAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogM3B4IDBweCAwcHggM3B4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXG4gICAgKTtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgd2VidmlldyBzdXBwb3J0cyB0aGUge0BsaW5rIE1vdXNlRXZlbnQjYnV0dG9uc30gcHJvcGVydHkuXG4gKiBMb29raW5nIGF0IHlvdSBtYWNPUyBIaWdoIFNpZXJyYSFcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGNhblRyYWNrQnV0dG9ucygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKG5ldyBNb3VzZUV2ZW50KCdtb3VzZWRvd24nKSkuYnV0dG9ucyA9PT0gMDtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgYnJvd3NlciBzdXBwb3J0cyByZW1vdmluZyBsaXN0ZW5lcnMgYnkgdHJpZ2dlcmluZyBhbiBBYm9ydFNpZ25hbFxuICogKHNlZSBodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZy9lbi1VUy9kb2NzL1dlYi9BUEkvRXZlbnRUYXJnZXQvYWRkRXZlbnRMaXN0ZW5lciNzaWduYWwpLlxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XG4gICAgaWYgKCFFdmVudFRhcmdldCB8fCAhQWJvcnRTaWduYWwgfHwgIUFib3J0Q29udHJvbGxlcilcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuXG4gICAgbGV0IHJlc3VsdCA9IHRydWU7XG5cbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcbiAgICBjb25zdCBjb250cm9sbGVyID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcbiAgICB0YXJnZXQuZGlzcGF0Y2hFdmVudChuZXcgQ3VzdG9tRXZlbnQoJ3Rlc3QnKSk7XG5cbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIFJlc29sdmVzIHRoZSBjbG9zZXN0IEhUTUxFbGVtZW50IGFuY2VzdG9yIG9mIGFuIGV2ZW50J3MgdGFyZ2V0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gZXZlbnRUYXJnZXQoZXZlbnQ6IEV2ZW50KTogSFRNTEVsZW1lbnQge1xuICAgIGlmIChldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkge1xuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0O1xuICAgIH0gZWxzZSBpZiAoIShldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkgJiYgZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgTm9kZSkge1xuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0LnBhcmVudEVsZW1lbnQgPz8gZG9jdW1lbnQuYm9keTtcbiAgICB9IGVsc2Uge1xuICAgICAgICByZXR1cm4gZG9jdW1lbnQuYm9keTtcbiAgICB9XG59XG5cbi8qKipcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcblxuIEJTRCAyLUNsYXVzZSBMaWNlbnNlXG5cbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcbiBBbGwgcmlnaHRzIHJlc2VydmVkLlxuXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XG4gbW9kaWZpY2F0aW9uLCBhcmUgcGVybWl0dGVkIHByb3ZpZGVkIHRoYXQgdGhlIGZvbGxvd2luZyBjb25kaXRpb25zIGFyZSBtZXQ6XG5cbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcbiBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lci5cblxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcbiB0aGlzIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyIGluIHRoZSBkb2N1bWVudGF0aW9uXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXG5cbiBUSElTIFNPRlRXQVJFIElTIFBST1ZJREVEIEJZIFRIRSBDT1BZUklHSFQgSE9MREVSUyBBTkQgQ09OVFJJQlVUT1JTIFwiQVMgSVNcIlxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcbiBESVNDTEFJTUVELiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQ09QWVJJR0hUIEhPTERFUiBPUiBDT05UUklCVVRPUlMgQkUgTElBQkxFXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxuIFNFUlZJQ0VTOyBMT1NTIE9GIFVTRSwgREFUQSwgT1IgUFJPRklUUzsgT1IgQlVTSU5FU1MgSU5URVJSVVBUSU9OKSBIT1dFVkVSXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxuIE9GIFRISVMgU09GVFdBUkUsIEVWRU4gSUYgQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YgU1VDSCBEQU1BR0UuXG5cbiAqKiovXG5cbmxldCBpc1JlYWR5ID0gZmFsc2U7XG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgKCkgPT4geyBpc1JlYWR5ID0gdHJ1ZSB9KTtcblxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjazogKCkgPT4gdm9pZCkge1xuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcbiAgICAgICAgY2FsbGJhY2soKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgY2FsbGJhY2spO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHR5cGUgeyBTY3JlZW4gfSBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XG5cbi8vIE5FVzogRHJvcHpvbmUgY29uc3RhbnRzXG5jb25zdCBEUk9QWk9ORV9BVFRSSUJVVEUgPSAnZGF0YS13YWlscy1kcm9wem9uZSc7XG5jb25zdCBEUk9QWk9ORV9IT1ZFUl9DTEFTUyA9ICd3YWlscy1kcm9wem9uZS1ob3Zlcic7IC8vIFVzZXIgY2FuIHN0eWxlIHRoaXMgY2xhc3NcbmxldCBjdXJyZW50SG92ZXJlZERyb3B6b25lOiBFbGVtZW50IHwgbnVsbCA9IG51bGw7XG5cbmNvbnN0IFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDA7XG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xuY29uc3QgQ2xvc2VNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gMjtcbmNvbnN0IERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICA9IDM7XG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xuY29uc3QgRm9jdXNNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNTtcbmNvbnN0IEZvcmNlUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICA9IDY7XG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xuY29uc3QgR2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gODtcbmNvbnN0IEdldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDk7XG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcbmNvbnN0IEhpZGVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDExO1xuY29uc3QgSXNGb2N1c2VkTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMTI7XG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcbmNvbnN0IElzTWF4aW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE0O1xuY29uc3QgSXNNaW5pbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTU7XG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcbmNvbnN0IE1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE3O1xuY29uc3QgTmFtZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTg7XG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcbmNvbnN0IFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICAgICA9IDIwO1xuY29uc3QgUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMjE7XG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcbmNvbnN0IFJlc3RvcmVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDIzO1xuY29uc3QgU2V0UG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgID0gMjQ7XG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcbmNvbnN0IFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QgICAgICAgICA9IDI2O1xuY29uc3QgU2V0RnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgICAgID0gMjc7XG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcbmNvbnN0IFNldE1heFNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDI5O1xuY29uc3QgU2V0TWluU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMzA7XG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcbmNvbnN0IFNldFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICA9IDMyO1xuY29uc3QgU2V0U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzM7XG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcbmNvbnN0IFNldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDM1O1xuY29uc3QgU2hvd01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzY7XG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcbmNvbnN0IFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICA9IDM4O1xuY29uc3QgVG9nZ2xlTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgID0gMzk7XG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MDtcbmNvbnN0IFVuTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQxO1xuY29uc3QgVW5NaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDI7XG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0MztcbmNvbnN0IFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDQ0O1xuY29uc3QgWm9vbUluTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gNDU7XG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NjtcbmNvbnN0IFpvb21SZXNldE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDQ3O1xuY29uc3QgV2luZG93RHJvcFpvbmVEcm9wcGVkICAgICAgICAgICAgID0gNDg7XG5cbmZ1bmN0aW9uIGdldERyb3B6b25lRWxlbWVudChlbGVtZW50OiBFbGVtZW50IHwgbnVsbCk6IEVsZW1lbnQgfCBudWxsIHtcbiAgICBpZiAoIWVsZW1lbnQpIHtcbiAgICAgICAgcmV0dXJuIG51bGw7XG4gICAgfVxuICAgIC8vIEFsbG93IGRyb3B6b25lIGF0dHJpYnV0ZSB0byBiZSBvbiB0aGUgZWxlbWVudCBpdHNlbGYgb3IgYW55IHBhcmVudFxuICAgIHJldHVybiBlbGVtZW50LmNsb3Nlc3QoYFske0RST1BaT05FX0FUVFJJQlVURX1dYCk7XG59XG5cbi8qKlxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgcG9zaXRpb24gb2YgYSB3aW5kb3cuXG4gKi9cbmludGVyZmFjZSBQb3NpdGlvbiB7XG4gICAgLyoqIFRoZSBob3Jpem9udGFsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXG4gICAgeDogbnVtYmVyO1xuICAgIC8qKiBUaGUgdmVydGljYWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cbiAgICB5OiBudW1iZXI7XG59XG5cbi8qKlxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgc2l6ZSBvZiBhIHdpbmRvdy5cbiAqL1xuaW50ZXJmYWNlIFNpemUge1xuICAgIC8qKiBUaGUgd2lkdGggb2YgdGhlIHdpbmRvdy4gKi9cbiAgICB3aWR0aDogbnVtYmVyO1xuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuICovXG4gICAgaGVpZ2h0OiBudW1iZXI7XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBjYWxsZXJTeW0gPSBTeW1ib2woXCJjYWxsZXJcIik7XG5cbmNsYXNzIFdpbmRvdyB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgcHJpdmF0ZSBbY2FsbGVyU3ltXTogKG1lc3NhZ2U6IG51bWJlciwgYXJncz86IGFueSkgPT4gUHJvbWlzZTxhbnk+O1xuXG4gICAgLyoqXG4gICAgICogSW5pdGlhbGlzZXMgYSB3aW5kb3cgb2JqZWN0IHdpdGggdGhlIHNwZWNpZmllZCBuYW1lLlxuICAgICAqXG4gICAgICogQHByaXZhdGVcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB0YXJnZXQgd2luZG93LlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG5hbWU6IHN0cmluZyA9ICcnKSB7XG4gICAgICAgIHRoaXNbY2FsbGVyU3ltXSA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuV2luZG93LCBuYW1lKVxuXG4gICAgICAgIC8vIGJpbmQgaW5zdGFuY2UgbWV0aG9kIHRvIG1ha2UgdGhlbSBlYXNpbHkgdXNhYmxlIGluIGV2ZW50IGhhbmRsZXJzXG4gICAgICAgIGZvciAoY29uc3QgbWV0aG9kIG9mIE9iamVjdC5nZXRPd25Qcm9wZXJ0eU5hbWVzKFdpbmRvdy5wcm90b3R5cGUpKSB7XG4gICAgICAgICAgICBpZiAoXG4gICAgICAgICAgICAgICAgbWV0aG9kICE9PSBcImNvbnN0cnVjdG9yXCJcbiAgICAgICAgICAgICAgICAmJiB0eXBlb2YgKHRoaXMgYXMgYW55KVttZXRob2RdID09PSBcImZ1bmN0aW9uXCJcbiAgICAgICAgICAgICkge1xuICAgICAgICAgICAgICAgICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9ICh0aGlzIGFzIGFueSlbbWV0aG9kXS5iaW5kKHRoaXMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogR2V0cyB0aGUgc3BlY2lmaWVkIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBnZXQuXG4gICAgICogQHJldHVybnMgVGhlIGNvcnJlc3BvbmRpbmcgd2luZG93IG9iamVjdC5cbiAgICAgKi9cbiAgICBHZXQobmFtZTogc3RyaW5nKTogV2luZG93IHtcbiAgICAgICAgcmV0dXJuIG5ldyBXaW5kb3cobmFtZSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUG9zaXRpb24oKTogUHJvbWlzZTxQb3NpdGlvbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFBvc2l0aW9uTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cbiAgICAgKi9cbiAgICBDZW50ZXIoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2VudGVyTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDbG9zZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBDbG9zZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShDbG9zZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqL1xuICAgIERpc2FibGVTaXplQ29uc3RyYWludHMoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRW5hYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXG4gICAgICovXG4gICAgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9jdXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIEZvY3VzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvY3VzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBGb3JjZXMgdGhlIHdpbmRvdyB0byByZWxvYWQgdGhlIHBhZ2UgYXNzZXRzLlxuICAgICAqL1xuICAgIEZvcmNlUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvcmNlUmVsb2FkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTd2l0Y2hlcyB0aGUgd2luZG93IHRvIGZ1bGxzY3JlZW4gbW9kZS5cbiAgICAgKi9cbiAgICBGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgc2NyZWVuIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG9uLlxuICAgICAqL1xuICAgIEdldFNjcmVlbigpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgY3VycmVudCB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB6b29tIGxldmVsLlxuICAgICAqL1xuICAgIEdldFpvb20oKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRab29tTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIEhlaWdodCgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhlaWdodE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSGlkZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhpZGVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZvY3VzZWQuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZvY3VzZWQuXG4gICAgICovXG4gICAgSXNGb2N1c2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRm9jdXNlZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZnVsbHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZnVsbHNjcmVlbi5cbiAgICAgKi9cbiAgICBJc0Z1bGxzY3JlZW4oKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1heGltaXNlZC5cbiAgICAgKi9cbiAgICBJc01heGltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc01heGltaXNlZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtaW5pbWlzZWQuXG4gICAgICovXG4gICAgSXNNaW5pbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNaW5pbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIE1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNaW5pbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBNaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNaW5pbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBOYW1lKCk6IFByb21pc2U8c3RyaW5nPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTmFtZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogT3BlbnMgdGhlIGRldmVsb3BtZW50IHRvb2xzIHBhbmUuXG4gICAgICovXG4gICAgT3BlbkRldlRvb2xzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE9wZW5EZXZUb29sc01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBSZWxhdGl2ZVBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZWxvYWRzIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKi9cbiAgICBSZWxvYWQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsb2FkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IHJlc2l6YWJsZS5cbiAgICAgKi9cbiAgICBSZXNpemFibGUoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzaXphYmxlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXN0b3JlcyB0aGUgd2luZG93IHRvIGl0cyBwcmV2aW91cyBzdGF0ZSBpZiBpdCB3YXMgcHJldmlvdXNseSBtaW5pbWlzZWQsIG1heGltaXNlZCBvciBmdWxsc2NyZWVuLlxuICAgICAqL1xuICAgIFJlc3RvcmUoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzdG9yZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFBvc2l0aW9uTWV0aG9kLCB7IHgsIHkgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gYWx3YXlzT25Ub3AgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIHN0YXkgb24gdG9wLlxuICAgICAqL1xuICAgIFNldEFsd2F5c09uVG9wKGFsd2F5c09uVG9wOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0QWx3YXlzT25Ub3BNZXRob2QsIHsgYWx3YXlzT25Ub3AgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSByIC0gVGhlIGRlc2lyZWQgcmVkIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICogQHBhcmFtIGcgLSBUaGUgZGVzaXJlZCBncmVlbiBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBiIC0gVGhlIGRlc2lyZWQgYmx1ZSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBhIC0gVGhlIGRlc2lyZWQgYWxwaGEgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKi9cbiAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyKHI6IG51bWJlciwgZzogbnVtYmVyLCBiOiBudW1iZXIsIGE6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QsIHsgciwgZywgYiwgYSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIHRoZSB3aW5kb3cgZnJhbWUgYW5kIHRpdGxlIGJhci5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBmcmFtZWxlc3MgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGZyYW1lbGVzcy5cbiAgICAgKi9cbiAgICBTZXRGcmFtZWxlc3MoZnJhbWVsZXNzOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0RnJhbWVsZXNzTWV0aG9kLCB7IGZyYW1lbGVzcyB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEaXNhYmxlcyB0aGUgc3lzdGVtIGZ1bGxzY3JlZW4gYnV0dG9uLlxuICAgICAqXG4gICAgICogQHBhcmFtIGVuYWJsZWQgLSBXaGV0aGVyIHRoZSBmdWxsc2NyZWVuIGJ1dHRvbiBzaG91bGQgYmUgZW5hYmxlZC5cbiAgICAgKi9cbiAgICBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZChlbmFibGVkOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWRNZXRob2QsIHsgZW5hYmxlZCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1heGltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRNYXhTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWF4U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWluaW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtaW5pbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldE1pblNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRNaW5TaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFJlbGF0aXZlUG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kLCB7IHgsIHkgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHBhcmFtIHJlc2l6YWJsZSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgcmVzaXphYmxlLlxuICAgICAqL1xuICAgIFNldFJlc2l6YWJsZShyZXNpemFibGU6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRSZXNpemFibGVNZXRob2QsIHsgcmVzaXphYmxlIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgdGl0bGUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB0aXRsZSAtIFRoZSBkZXNpcmVkIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0VGl0bGUodGl0bGU6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFRpdGxlTWV0aG9kLCB7IHRpdGxlIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB6b29tIC0gVGhlIGRlc2lyZWQgem9vbSBsZXZlbC5cbiAgICAgKi9cbiAgICBTZXRab29tKHpvb206IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFpvb21NZXRob2QsIHsgem9vbSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTaG93cyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2hvd01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNpemUoKTogUHJvbWlzZTxTaXplPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2l6ZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnVsbHNjcmVlbiBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZUZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gbWF4aW1pc2VkIGFuZCBub3JtYWwuXG4gICAgICovXG4gICAgVG9nZ2xlTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLWZ1bGxzY3JlZW5zIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5GdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5NYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1taW5pbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBVbk1pbmltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFdpZHRoKCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oV2lkdGhNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFpvb21zIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgWm9vbSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBJbmNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tSW4oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbUluTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEZWNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tT3V0KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21PdXRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlc2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxuICAgICAqL1xuICAgIFpvb21SZXNldCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tUmVzZXRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEhhbmRsZXMgZmlsZSBkcm9wcyBvcmlnaW5hdGluZyBmcm9tIHBsYXRmb3JtLXNwZWNpZmljIGNvZGUgKGUuZy4sIG1hY09TIG5hdGl2ZSBkcmFnLWFuZC1kcm9wKS5cbiAgICAgKiBHYXRoZXJzIGluZm9ybWF0aW9uIGFib3V0IHRoZSBkcm9wIHRhcmdldCBlbGVtZW50IGFuZCBzZW5kcyBpdCBiYWNrIHRvIHRoZSBHbyBiYWNrZW5kLlxuICAgICAqXG4gICAgICogQHBhcmFtIGZpbGVuYW1lcyAtIEFuIGFycmF5IG9mIGZpbGUgcGF0aHMgKHN0cmluZ3MpIHRoYXQgd2VyZSBkcm9wcGVkLlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIHgtY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudC5cbiAgICAgKiBAcGFyYW0geSAtIFRoZSB5LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQuXG4gICAgICovXG4gICAgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgICAgICBjb25zdCBlbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KTtcblxuICAgICAgICAvLyBORVc6IENoZWNrIGlmIHRoZSBkcm9wIHRhcmdldCBpcyBhIHZhbGlkIGRyb3B6b25lXG4gICAgICAgIGNvbnN0IGRyb3B6b25lVGFyZ2V0ID0gZ2V0RHJvcHpvbmVFbGVtZW50KGVsZW1lbnQpO1xuXG4gICAgICAgIGlmICghZHJvcHpvbmVUYXJnZXQpIHtcbiAgICAgICAgICAgIGNvbnNvbGUubG9nKGBXYWlscyBSdW50aW1lOiBEcm9wIG9uIGVsZW1lbnQgKG9yIG5vIGVsZW1lbnQpIGF0ICR7eH0sJHt5fSB3aGljaCBpcyBub3QgYSBkZXNpZ25hdGVkIGRyb3B6b25lLiBJZ25vcmluZy4gRWxlbWVudDpgLCBlbGVtZW50KTtcbiAgICAgICAgICAgIC8vIE5vIG5lZWQgdG8gY2FsbCBiYWNrZW5kIGlmIG5vdCBhIHZhbGlkIGRyb3B6b25lIHRhcmdldFxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc29sZS5sb2coYFdhaWxzIFJ1bnRpbWU6IERyb3Agb24gZGVzaWduYXRlZCBkcm9wem9uZS4gRWxlbWVudCBhdCAoJHt4fSwgJHt5fSk6YCwgZWxlbWVudCwgJ0VmZmVjdGl2ZSBkcm9wem9uZTonLCBkcm9wem9uZVRhcmdldCk7XG4gICAgICAgIGNvbnN0IGVsZW1lbnREZXRhaWxzID0ge1xuICAgICAgICAgICAgaWQ6IGRyb3B6b25lVGFyZ2V0LmlkLFxuICAgICAgICAgICAgY2xhc3NMaXN0OiBBcnJheS5mcm9tKGRyb3B6b25lVGFyZ2V0LmNsYXNzTGlzdCksXG4gICAgICAgICAgICBhdHRyaWJ1dGVzOiB7fSBhcyB7IFtrZXk6IHN0cmluZ106IHN0cmluZyB9LFxuICAgICAgICB9O1xuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IGRyb3B6b25lVGFyZ2V0LmF0dHJpYnV0ZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgICAgIGNvbnN0IGF0dHIgPSBkcm9wem9uZVRhcmdldC5hdHRyaWJ1dGVzW2ldO1xuICAgICAgICAgICAgZWxlbWVudERldGFpbHMuYXR0cmlidXRlc1thdHRyLm5hbWVdID0gYXR0ci52YWx1ZTtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XG4gICAgICAgICAgICBmaWxlbmFtZXMsXG4gICAgICAgICAgICB4LFxuICAgICAgICAgICAgeSxcbiAgICAgICAgICAgIGVsZW1lbnREZXRhaWxzLFxuICAgICAgICB9O1xuXG4gICAgICAgIHRoaXNbY2FsbGVyU3ltXShXaW5kb3dEcm9wWm9uZURyb3BwZWQsIHBheWxvYWQpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXG4gKi9cbmNvbnN0IHRoaXNXaW5kb3cgPSBuZXcgV2luZG93KCcnKTtcblxuLy8gTkVXOiBHbG9iYWwgRHJhZyBFdmVudCBMaXN0ZW5lcnNcbmZ1bmN0aW9uIHNldHVwR2xvYmFsRHJvcHpvbmVMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZG9jRWxlbWVudCA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudDtcbiAgICBsZXQgZHJhZ0VudGVyQ291bnRlciA9IDA7IC8vIFRvIGhhbmRsZSBkcmFnZW50ZXIvZHJhZ2xlYXZlIG9uIGNoaWxkIGVsZW1lbnRzXG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdlbnRlcicsIChldmVudCkgPT4ge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyICYmIGV2ZW50LmRhdGFUcmFuc2Zlci50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlcisrO1xuICAgICAgICAgICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSk7XG4gICAgICAgICAgICBjb25zdCBkcm9wem9uZSA9IGdldERyb3B6b25lRWxlbWVudCh0YXJnZXRFbGVtZW50KTtcblxuICAgICAgICAgICAgLy8gQ2xlYXIgcHJldmlvdXMgaG92ZXIgcmVnYXJkbGVzcywgdGhlbiBhcHBseSBuZXcgaWYgdmFsaWRcbiAgICAgICAgICAgIGlmIChjdXJyZW50SG92ZXJlZERyb3B6b25lICYmIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgIT09IGRyb3B6b25lKSB7XG4gICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QucmVtb3ZlKERST1BaT05FX0hPVkVSX0NMQVNTKTtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgaWYgKGRyb3B6b25lKSB7XG4gICAgICAgICAgICAgICAgZHJvcHpvbmUuY2xhc3NMaXN0LmFkZChEUk9QWk9ORV9IT1ZFUl9DTEFTUyk7XG4gICAgICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG4gICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZSA9IGRyb3B6b25lO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcbiAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lID0gbnVsbDsgLy8gRW5zdXJlIGl0J3MgY2xlYXJlZCBpZiBubyBkcm9wem9uZSBmb3VuZFxuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnb3ZlcicsIChldmVudCkgPT4ge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBOZWNlc3NhcnkgdG8gYWxsb3cgZHJvcFxuICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyICYmIGV2ZW50LmRhdGFUcmFuc2Zlci50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgLy8gTm8gbmVlZCB0byBxdWVyeSBlbGVtZW50RnJvbVBvaW50IGFnYWluIGlmIGFscmVhZHkgaGFuZGxlZCBieSBkcmFnZW50ZXIgY29ycmVjdGx5XG4gICAgICAgICAgICAvLyBKdXN0IGVuc3VyZSBkcm9wRWZmZWN0IGlzIGNvbnRpbnVvdXNseSBzZXQgYmFzZWQgb24gY3VycmVudEhvdmVyZWREcm9wem9uZVxuICAgICAgICAgICAgaWYgKGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUpIHtcbiAgICAgICAgICAgICAgICAgLy8gUmUtYXBwbHkgY2xhc3MganVzdCBpbiBjYXNlIGl0IHdhcyByZW1vdmVkIGJ5IHNvbWUgb3RoZXIgSlNcbiAgICAgICAgICAgICAgICBpZighY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QuY29udGFpbnMoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpKSB7XG4gICAgICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY2xhc3NMaXN0LmFkZChEUk9QWk9ORV9IT1ZFUl9DTEFTUyk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIGlmIChldmVudC5kYXRhVHJhbnNmZXIgJiYgZXZlbnQuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICBkcmFnRW50ZXJDb3VudGVyLS07XG4gICAgICAgICAgICAvLyBPbmx5IHJlbW92ZSBob3ZlciBpZiBkcmFnIHRydWx5IGxlZnQgdGhlIHdpbmRvdyBvciB0aGUgbGFzdCBkcm9wem9uZVxuICAgICAgICAgICAgaWYgKGRyYWdFbnRlckNvdW50ZXIgPT09IDAgfHwgZXZlbnQucmVsYXRlZFRhcmdldCA9PT0gbnVsbCB8fCAoY3VycmVudEhvdmVyZWREcm9wem9uZSAmJiAhY3VycmVudEhvdmVyZWREcm9wem9uZS5jb250YWlucyhldmVudC5yZWxhdGVkVGFyZ2V0IGFzIE5vZGUpKSkge1xuICAgICAgICAgICAgICAgIGlmIChjdXJyZW50SG92ZXJlZERyb3B6b25lKSB7XG4gICAgICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY2xhc3NMaXN0LnJlbW92ZShEUk9QWk9ORV9IT1ZFUl9DTEFTUyk7XG4gICAgICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgPSBudWxsO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICBkcmFnRW50ZXJDb3VudGVyID0gMDsgLy8gUmVzZXQgY291bnRlciBpZiBpdCB3ZW50IG5lZ2F0aXZlIG9yIGxlZnQgd2luZG93XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gUHJldmVudCBkZWZhdWx0IGJyb3dzZXIgZmlsZSBoYW5kbGluZ1xuICAgICAgICBkcmFnRW50ZXJDb3VudGVyID0gMDsgLy8gUmVzZXQgY291bnRlclxuICAgICAgICBpZiAoY3VycmVudEhvdmVyZWREcm9wem9uZSkge1xuICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QucmVtb3ZlKERST1BaT05FX0hPVkVSX0NMQVNTKTtcbiAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgPSBudWxsO1xuICAgICAgICB9XG4gICAgICAgIC8vIFRoZSBhY3R1YWwgZHJvcCBwcm9jZXNzaW5nIGlzIGluaXRpYXRlZCBieSB0aGUgbmF0aXZlIHNpZGUgY2FsbGluZyBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wXG4gICAgICAgIC8vIEhhbmRsZVBsYXRmb3JtRmlsZURyb3Agd2lsbCB0aGVuIGNoZWNrIGlmIHRoZSBkcm9wIHdhcyBvbiBhIHZhbGlkIHpvbmUuXG4gICAgfSwgZmFsc2UpO1xufVxuXG4vLyBJbml0aWFsaXplIGxpc3RlbmVycyB3aGVuIHRoZSBzY3JpcHQgbG9hZHNcbmlmICh0eXBlb2Ygd2luZG93ICE9PSBcInVuZGVmaW5lZFwiICYmIHR5cGVvZiBkb2N1bWVudCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgIHNldHVwR2xvYmFsRHJvcHpvbmVMaXN0ZW5lcnMoKTtcbn1cblxuZXhwb3J0IGRlZmF1bHQgdGhpc1dpbmRvdztcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0ICogYXMgUnVudGltZSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmNcIjtcblxuLy8gTk9URTogdGhlIGZvbGxvd2luZyBtZXRob2RzIE1VU1QgYmUgaW1wb3J0ZWQgZXhwbGljaXRseSBiZWNhdXNlIG9mIGhvdyBlc2J1aWxkIGluamVjdGlvbiB3b3Jrc1xuaW1wb3J0IHsgRW5hYmxlIGFzIEVuYWJsZVdNTCB9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWxcIjtcbmltcG9ydCB7IGRlYnVnTG9nIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzXCI7XG5cbndpbmRvdy53YWlscyA9IFJ1bnRpbWU7XG5FbmFibGVXTUwoKTtcblxuaWYgKERFQlVHKSB7XG4gICAgZGVidWdMb2coXCJXYWlscyBSdW50aW1lIExvYWRlZFwiKVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuU3lzdGVtKTtcblxuY29uc3QgU3lzdGVtSXNEYXJrTW9kZSA9IDA7XG5jb25zdCBTeXN0ZW1FbnZpcm9ubWVudCA9IDE7XG5jb25zdCBBcHBsaWNhdGlvbkZpbGVzRHJvcHBlZFdpdGhDb250ZXh0ID0gMTAwOyAvLyBOZXcgbWV0aG9kIElEIGZvciBlbnJpY2hlZCBkcm9wIGV2ZW50XG5cbmNvbnN0IF9pbnZva2UgPSAoZnVuY3Rpb24gKCkge1xuICAgIHRyeSB7XG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZSkge1xuICAgICAgICAgICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZS5iaW5kKCh3aW5kb3cgYXMgYW55KS5jaHJvbWUud2Vidmlldyk7XG4gICAgICAgIH0gZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndlYmtpdD8ubWVzc2FnZUhhbmRsZXJzPy5bJ2V4dGVybmFsJ10/LnBvc3RNZXNzYWdlKSB7XG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10ucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXSk7XG4gICAgICAgIH1cbiAgICB9IGNhdGNoKGUpIHt9XG5cbiAgICBjb25zb2xlLndhcm4oJ1xcbiVjXHUyNkEwXHVGRTBGIEJyb3dzZXIgRW52aXJvbm1lbnQgRGV0ZWN0ZWQgJWNcXG5cXG4lY09ubHkgVUkgcHJldmlld3MgYXJlIGF2YWlsYWJsZSBpbiB0aGUgYnJvd3Nlci4gRm9yIGZ1bGwgZnVuY3Rpb25hbGl0eSwgcGxlYXNlIHJ1biB0aGUgYXBwbGljYXRpb24gaW4gZGVza3RvcCBtb2RlLlxcbk1vcmUgaW5mb3JtYXRpb24gYXQ6IGh0dHBzOi8vdjMud2FpbHMuaW8vbGVhcm4vYnVpbGQvI3VzaW5nLWEtYnJvd3Nlci1mb3ItZGV2ZWxvcG1lbnRcXG4nLFxuICAgICAgICAnYmFja2dyb3VuZDogI2ZmZmZmZjsgY29sb3I6ICMwMDAwMDA7IGZvbnQtd2VpZ2h0OiBib2xkOyBwYWRkaW5nOiA0cHggOHB4OyBib3JkZXItcmFkaXVzOiA0cHg7IGJvcmRlcjogMnB4IHNvbGlkICMwMDAwMDA7JyxcbiAgICAgICAgJ2JhY2tncm91bmQ6IHRyYW5zcGFyZW50OycsXG4gICAgICAgICdjb2xvcjogI2ZmZmZmZjsgZm9udC1zdHlsZTogaXRhbGljOyBmb250LXdlaWdodDogYm9sZDsnKTtcbiAgICByZXR1cm4gbnVsbDtcbn0pKCk7XG5cbmV4cG9ydCBmdW5jdGlvbiBpbnZva2UobXNnOiBhbnkpOiB2b2lkIHtcbiAgICBfaW52b2tlPy4obXNnKTtcbn1cblxuLyoqXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGEgYm9vbGVhbiB2YWx1ZSBpbmRpY2F0aW5nIGlmIHRoZSBzeXN0ZW0gaXMgaW4gZGFyayBtb2RlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNEYXJrTW9kZSgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1Jc0RhcmtNb2RlKTtcbn1cblxuLyoqXG4gKiBGZXRjaGVzIHRoZSBjYXBhYmlsaXRpZXMgb2YgdGhlIGFwcGxpY2F0aW9uIGZyb20gdGhlIHNlcnZlci5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyB0aGUgY2FwYWJpbGl0aWVzLlxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gQ2FwYWJpbGl0aWVzKCk6IFByb21pc2U8UmVjb3JkPHN0cmluZywgYW55Pj4ge1xuICAgIGxldCByZXNwb25zZSA9IGF3YWl0IGZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKTtcbiAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJjb3VsZCBub3QgZmV0Y2ggY2FwYWJpbGl0aWVzOiBcIiArIHJlc3BvbnNlLnN0YXR1c1RleHQpO1xuICAgIH1cbn1cblxuZXhwb3J0IGludGVyZmFjZSBPU0luZm8ge1xuICAgIC8qKiBUaGUgYnJhbmRpbmcgb2YgdGhlIE9TLiAqL1xuICAgIEJyYW5kaW5nOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBJRCBvZiB0aGUgT1MuICovXG4gICAgSUQ6IHN0cmluZztcbiAgICAvKiogVGhlIG5hbWUgb2YgdGhlIE9TLiAqL1xuICAgIE5hbWU6IHN0cmluZztcbiAgICAvKiogVGhlIHZlcnNpb24gb2YgdGhlIE9TLiAqL1xuICAgIFZlcnNpb246IHN0cmluZztcbn1cblxuZXhwb3J0IGludGVyZmFjZSBFbnZpcm9ubWVudEluZm8ge1xuICAgIC8qKiBUaGUgYXJjaGl0ZWN0dXJlIG9mIHRoZSBzeXN0ZW0uICovXG4gICAgQXJjaDogc3RyaW5nO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS4gKi9cbiAgICBEZWJ1ZzogYm9vbGVhbjtcbiAgICAvKiogVGhlIG9wZXJhdGluZyBzeXN0ZW0gaW4gdXNlLiAqL1xuICAgIE9TOiBzdHJpbmc7XG4gICAgLyoqIERldGFpbHMgb2YgdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXG4gICAgT1NJbmZvOiBPU0luZm87XG4gICAgLyoqIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uICovXG4gICAgUGxhdGZvcm1JbmZvOiBSZWNvcmQ8c3RyaW5nLCBhbnk+O1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIE9TIGFuZCBzeXN0ZW0gYXJjaGl0ZWN0dXJlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKTogUHJvbWlzZTxFbnZpcm9ubWVudEluZm8+IHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1FbnZpcm9ubWVudCk7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cbiAqXG4gKiBAcmV0dXJuIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcIndpbmRvd3NcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cbiAqXG4gKiBAcmV0dXJucyBSZXR1cm5zIHRydWUgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eCwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNMaW51eCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJsaW51eFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBhIG1hY09TIG9wZXJhdGluZyBzeXN0ZW0uXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImRhcndpblwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhbWQ2NFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBUk0oKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm1cIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxuICpcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm02NFwiO1xufVxuXG4vKipcbiAqIFJlcG9ydHMgd2hldGhlciB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGFwcCBpcyBiZWluZyBydW4gaW4gZGVidWcgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGVidWcoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIEJvb2xlYW4od2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5EZWJ1Zyk7XG59XG5cbi8qKlxuICogSGFuZGxlcyBmaWxlIGRyb3BzIG9yaWdpbmF0aW5nIGZyb20gcGxhdGZvcm0tc3BlY2lmaWMgY29kZSAoZS5nLiwgbWFjT1MgbmF0aXZlIGRyYWctYW5kLWRyb3ApLlxuICogR2F0aGVycyBpbmZvcm1hdGlvbiBhYm91dCB0aGUgZHJvcCB0YXJnZXQgZWxlbWVudCBhbmQgc2VuZHMgaXQgYmFjayB0byB0aGUgR28gYmFja2VuZC5cbiAqXG4gKiBAcGFyYW0gZmlsZW5hbWVzIC0gQW4gYXJyYXkgb2YgZmlsZSBwYXRocyAoc3RyaW5ncykgdGhhdCB3ZXJlIGRyb3BwZWQuXG4gKiBAcGFyYW0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQuXG4gKiBAcGFyYW0geSAtIFRoZSB5LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wKGZpbGVuYW1lczogc3RyaW5nW10sIHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XG4gICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgY29uc3QgZWxlbWVudElkID0gZWxlbWVudCA/IGVsZW1lbnQuaWQgOiAnJztcbiAgICBjb25zdCBjbGFzc0xpc3QgPSBlbGVtZW50ID8gQXJyYXkuZnJvbShlbGVtZW50LmNsYXNzTGlzdCkgOiBbXTtcblxuICAgIGNvbnN0IHBheWxvYWQgPSB7XG4gICAgICAgIGZpbGVuYW1lcyxcbiAgICAgICAgeCxcbiAgICAgICAgeSxcbiAgICAgICAgZWxlbWVudElkLFxuICAgICAgICBjbGFzc0xpc3QsXG4gICAgfTtcblxuICAgIGNhbGwoQXBwbGljYXRpb25GaWxlc0Ryb3BwZWRXaXRoQ29udGV4dCwgcGF5bG9hZClcbiAgICAgICAgLnRoZW4oKCkgPT4ge1xuICAgICAgICAgICAgLy8gT3B0aW9uYWw6IExvZyBzdWNjZXNzIG9yIGhhbmRsZSBpZiBuZWVkZWRcbiAgICAgICAgICAgIGNvbnNvbGUubG9nKFwiUGxhdGZvcm0gZmlsZSBkcm9wIHByb2Nlc3NlZCBhbmQgc2VudCB0byBHby5cIik7XG4gICAgICAgIH0pXG4gICAgICAgIC5jYXRjaChlcnIgPT4ge1xuICAgICAgICAgICAgLy8gT3B0aW9uYWw6IExvZyBlcnJvclxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIHNlbmRpbmcgcGxhdGZvcm0gZmlsZSBkcm9wIHRvIEdvOlwiLCBlcnIpO1xuICAgICAgICB9KTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IElzRGVidWcgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCB7IGV2ZW50VGFyZ2V0IH0gZnJvbSBcIi4vdXRpbHNcIjtcblxuLy8gc2V0dXBcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcblxuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkOiBzdHJpbmcsIHg6IG51bWJlciwgeTogbnVtYmVyLCBkYXRhOiBhbnkpOiB2b2lkIHtcbiAgICB2b2lkIGNhbGwoQ29udGV4dE1lbnVPcGVuLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xuXG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcbiAgICBjb25zdCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpLnRyaW0oKTtcblxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBjb25zdCBkYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnUtZGF0YVwiKTtcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBkYXRhKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50LCB0YXJnZXQpO1xuICAgIH1cbn1cblxuXG4vKlxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBhdXRvOyAoZGVmYXVsdCkgd2lsbCBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudSBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZSBPUiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkIE9SIGVsZW1lbnQgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcblxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcbiovXG5mdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50OiBNb3VzZUV2ZW50LCB0YXJnZXQ6IEhUTUxFbGVtZW50KSB7XG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XG4gICAgaWYgKElzRGVidWcoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIHN3aXRjaCAod2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKSkge1xuICAgICAgICBjYXNlICdzaG93JzpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSAnaGlkZSc6XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgaWYgKHRhcmdldC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxuICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSBzZWxlY3Rpb24gJiYgc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMDtcbiAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcbiAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XG4gICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IHRhcmdldCkge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgdGFnIGlzIGlucHV0IG9yIHRleHRhcmVhLlxuICAgIGlmICh0YXJnZXQgaW5zdGFuY2VvZiBIVE1MSW5wdXRFbGVtZW50IHx8IHRhcmdldCBpbnN0YW5jZW9mIEhUTUxUZXh0QXJlYUVsZW1lbnQpIHtcbiAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIXRhcmdldC5yZWFkT25seSAmJiAhdGFyZ2V0LmRpc2FibGVkKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXG4gKlxuICogQHBhcmFtIGtleSAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cbiAqIEByZXR1cm4gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5OiBzdHJpbmcpOiBhbnkge1xuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleV07XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJVbmFibGUgdG8gcmV0cmlldmUgZmxhZyAnXCIgKyBrZXkgKyBcIic6IFwiICsgZSwgeyBjYXVzZTogZSB9KTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IGludm9rZSwgSXNXaW5kb3dzIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyBHZXRGbGFnIH0gZnJvbSBcIi4vZmxhZ3MuanNcIjtcbmltcG9ydCB7IGNhblRyYWNrQnV0dG9ucywgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuXG4vLyBTZXR1cFxubGV0IGNhbkRyYWcgPSBmYWxzZTtcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xuXG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XG5sZXQgcmVzaXppbmcgPSBmYWxzZTtcbmxldCByZXNpemVFZGdlOiBzdHJpbmcgPSBcIlwiO1xubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcblxubGV0IGJ1dHRvbnMgPSAwO1xuY29uc3QgYnV0dG9uc1RyYWNrZWQgPSBjYW5UcmFja0J1dHRvbnMoKTtcblxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9ICh2YWx1ZTogYm9vbGVhbik6IHZvaWQgPT4ge1xuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xuICAgIGlmICghcmVzaXphYmxlKSB7XG4gICAgICAgIC8vIFN0b3AgcmVzaXppbmcgaWYgaW4gcHJvZ3Jlc3MuXG4gICAgICAgIGNhblJlc2l6ZSA9IHJlc2l6aW5nID0gZmFsc2U7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbn07XG5cbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XG5mb3IgKGNvbnN0IGV2IG9mIFsnY2xpY2snLCAnY29udGV4dG1lbnUnLCAnZGJsY2xpY2snXSkge1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKGV2LCBzdXBwcmVzc0V2ZW50LCB7IGNhcHR1cmU6IHRydWUgfSk7XG59XG5cbmZ1bmN0aW9uIHN1cHByZXNzRXZlbnQoZXZlbnQ6IEV2ZW50KSB7XG4gICAgLy8gU3VwcHJlc3MgY2xpY2sgZXZlbnRzIHdoaWxlIHJlc2l6aW5nIG9yIGRyYWdnaW5nLlxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxufVxuXG4vLyBVc2UgY29uc3RhbnRzIHRvIGF2b2lkIGNvbXBhcmluZyBzdHJpbmdzIG11bHRpcGxlIHRpbWVzLlxuY29uc3QgTW91c2VEb3duID0gMDtcbmNvbnN0IE1vdXNlVXAgICA9IDE7XG5jb25zdCBNb3VzZU1vdmUgPSAyO1xuXG5mdW5jdGlvbiB1cGRhdGUoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICAvLyBXaW5kb3dzIHN1cHByZXNzZXMgbW91c2UgZXZlbnRzIGF0IHRoZSBlbmQgb2YgZHJhZ2dpbmcgb3IgcmVzaXppbmcsXG4gICAgLy8gc28gd2UgbmVlZCB0byBiZSBzbWFydCBhbmQgc3ludGhlc2l6ZSBidXR0b24gZXZlbnRzLlxuXG4gICAgbGV0IGV2ZW50VHlwZTogbnVtYmVyLCBldmVudEJ1dHRvbnMgPSBldmVudC5idXR0b25zO1xuICAgIHN3aXRjaCAoZXZlbnQudHlwZSkge1xuICAgICAgICBjYXNlICdtb3VzZWRvd24nOlxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VEb3duO1xuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zIHwgKDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIGNhc2UgJ21vdXNldXAnOlxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VVcDtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyAmIH4oMSA8PCBldmVudC5idXR0b24pOyB9XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlTW92ZTtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9uczsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgfVxuXG4gICAgbGV0IHJlbGVhc2VkID0gYnV0dG9ucyAmIH5ldmVudEJ1dHRvbnM7XG4gICAgbGV0IHByZXNzZWQgPSBldmVudEJ1dHRvbnMgJiB+YnV0dG9ucztcblxuICAgIGJ1dHRvbnMgPSBldmVudEJ1dHRvbnM7XG5cbiAgICAvLyBTeW50aGVzaXplIGEgcmVsZWFzZS1wcmVzcyBzZXF1ZW5jZSBpZiB3ZSBkZXRlY3QgYSBwcmVzcyBvZiBhbiBhbHJlYWR5IHByZXNzZWQgYnV0dG9uLlxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlRG93biAmJiAhKHByZXNzZWQgJiBldmVudC5idXR0b24pKSB7XG4gICAgICAgIHJlbGVhc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XG4gICAgICAgIHByZXNzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcbiAgICB9XG5cbiAgICAvLyBTdXBwcmVzcyBhbGwgYnV0dG9uIGV2ZW50cyBkdXJpbmcgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLFxuICAgIC8vIHVubGVzcyB0aGlzIGlzIGEgbW91c2V1cCBldmVudCB0aGF0IGlzIGVuZGluZyBhIGRyYWcgYWN0aW9uLlxuICAgIGlmIChcbiAgICAgICAgZXZlbnRUeXBlICE9PSBNb3VzZU1vdmUgLy8gRmFzdCBwYXRoIGZvciBtb3VzZW1vdmVcbiAgICAgICAgJiYgcmVzaXppbmdcbiAgICAgICAgfHwgKFxuICAgICAgICAgICAgZHJhZ2dpbmdcbiAgICAgICAgICAgICYmIChcbiAgICAgICAgICAgICAgICBldmVudFR5cGUgPT09IE1vdXNlRG93blxuICAgICAgICAgICAgICAgIHx8IGV2ZW50LmJ1dHRvbiAhPT0gMFxuICAgICAgICAgICAgKVxuICAgICAgICApXG4gICAgKSB7XG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG5cbiAgICAvLyBIYW5kbGUgcmVsZWFzZXNcbiAgICBpZiAocmVsZWFzZWQgJiAxKSB7IHByaW1hcnlVcChldmVudCk7IH1cbiAgICAvLyBIYW5kbGUgcHJlc3Nlc1xuICAgIGlmIChwcmVzc2VkICYgMSkgeyBwcmltYXJ5RG93bihldmVudCk7IH1cblxuICAgIC8vIEhhbmRsZSBtb3VzZW1vdmVcbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZU1vdmUpIHsgb25Nb3VzZU1vdmUoZXZlbnQpOyB9O1xufVxuXG5mdW5jdGlvbiBwcmltYXJ5RG93bihldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xuICAgIC8vIFJlc2V0IHJlYWRpbmVzcyBzdGF0ZS5cbiAgICBjYW5EcmFnID0gZmFsc2U7XG4gICAgY2FuUmVzaXplID0gZmFsc2U7XG5cbiAgICAvLyBJZ25vcmUgcmVwZWF0ZWQgY2xpY2tzIG9uIG1hY09TIGFuZCBMaW51eC5cbiAgICBpZiAoIUlzV2luZG93cygpKSB7XG4gICAgICAgIGlmIChldmVudC50eXBlID09PSAnbW91c2Vkb3duJyAmJiBldmVudC5idXR0b24gPT09IDAgJiYgZXZlbnQuZGV0YWlsICE9PSAxKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICBpZiAocmVzaXplRWRnZSkge1xuICAgICAgICAvLyBSZWFkeSB0byByZXNpemUgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZS5cbiAgICAgICAgY2FuUmVzaXplID0gdHJ1ZTtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWcgb3BlcmF0aW9ucyB3aGVuIG9uIHJlc2l6ZSBlZGdlcy5cbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFJldHJpZXZlIHRhcmdldCBlbGVtZW50XG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xuXG4gICAgLy8gUmVhZHkgdG8gZHJhZyBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lIG9uIGEgZHJhZ2dhYmxlIGVsZW1lbnQuXG4gICAgLy8gSWdub3JlIGNsaWNrcyBvbiB0aGUgc2Nyb2xsYmFyLlxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KTtcbiAgICBjYW5EcmFnID0gKFxuICAgICAgICBzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13YWlscy1kcmFnZ2FibGVcIikudHJpbSgpID09PSBcImRyYWdcIlxuICAgICAgICAmJiAoXG4gICAgICAgICAgICBldmVudC5vZmZzZXRYIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nTGVmdCkgPCB0YXJnZXQuY2xpZW50V2lkdGhcbiAgICAgICAgICAgICYmIGV2ZW50Lm9mZnNldFkgLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdUb3ApIDwgdGFyZ2V0LmNsaWVudEhlaWdodFxuICAgICAgICApXG4gICAgKTtcbn1cblxuZnVuY3Rpb24gcHJpbWFyeVVwKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgLy8gU3RvcCBkcmFnZ2luZyBhbmQgcmVzaXppbmcuXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xuICAgIGRyYWdnaW5nID0gZmFsc2U7XG4gICAgY2FuUmVzaXplID0gZmFsc2U7XG4gICAgcmVzaXppbmcgPSBmYWxzZTtcbn1cblxuY29uc3QgY3Vyc29yRm9yRWRnZSA9IE9iamVjdC5mcmVlemUoe1xuICAgIFwic2UtcmVzaXplXCI6IFwibndzZS1yZXNpemVcIixcbiAgICBcInN3LXJlc2l6ZVwiOiBcIm5lc3ctcmVzaXplXCIsXG4gICAgXCJudy1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxuICAgIFwibmUtcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcbiAgICBcInctcmVzaXplXCI6IFwiZXctcmVzaXplXCIsXG4gICAgXCJuLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxuICAgIFwicy1yZXNpemVcIjogXCJucy1yZXNpemVcIixcbiAgICBcImUtcmVzaXplXCI6IFwiZXctcmVzaXplXCIsXG59KVxuXG5mdW5jdGlvbiBzZXRSZXNpemUoZWRnZT86IGtleW9mIHR5cGVvZiBjdXJzb3JGb3JFZGdlKTogdm9pZCB7XG4gICAgaWYgKGVkZ2UpIHtcbiAgICAgICAgaWYgKCFyZXNpemVFZGdlKSB7IGRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvcjsgfVxuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGN1cnNvckZvckVkZ2VbZWRnZV07XG4gICAgfSBlbHNlIGlmICghZWRnZSAmJiByZXNpemVFZGdlKSB7XG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gZGVmYXVsdEN1cnNvcjtcbiAgICB9XG5cbiAgICByZXNpemVFZGdlID0gZWRnZSB8fCBcIlwiO1xufVxuXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xuICAgIGlmIChjYW5SZXNpemUgJiYgcmVzaXplRWRnZSkge1xuICAgICAgICAvLyBTdGFydCByZXNpemluZy5cbiAgICAgICAgcmVzaXppbmcgPSB0cnVlO1xuICAgICAgICBpbnZva2UoXCJ3YWlsczpyZXNpemU6XCIgKyByZXNpemVFZGdlKTtcbiAgICB9IGVsc2UgaWYgKGNhbkRyYWcpIHtcbiAgICAgICAgLy8gU3RhcnQgZHJhZ2dpbmcuXG4gICAgICAgIGRyYWdnaW5nID0gdHJ1ZTtcbiAgICAgICAgaW52b2tlKFwid2FpbHM6ZHJhZ1wiKTtcbiAgICB9XG5cbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcbiAgICAgICAgLy8gRWl0aGVyIGRyYWcgb3IgcmVzaXplIGlzIG9uZ29pbmcsXG4gICAgICAgIC8vIHJlc2V0IHJlYWRpbmVzcyBhbmQgc3RvcCBwcm9jZXNzaW5nLlxuICAgICAgICBjYW5EcmFnID0gY2FuUmVzaXplID0gZmFsc2U7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIXJlc2l6YWJsZSB8fCAhSXNXaW5kb3dzKCkpIHtcbiAgICAgICAgaWYgKHJlc2l6ZUVkZ2UpIHsgc2V0UmVzaXplKCk7IH1cbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZUhlaWdodCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlSGVpZ2h0XCIpIHx8IDU7XG4gICAgY29uc3QgcmVzaXplSGFuZGxlV2lkdGggPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZVdpZHRoXCIpIHx8IDU7XG5cbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXMuXG4gICAgY29uc3QgY29ybmVyRXh0cmEgPSBHZXRGbGFnKFwicmVzaXplQ29ybmVyRXh0cmFcIikgfHwgMTA7XG5cbiAgICBjb25zdCByaWdodEJvcmRlciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgY29uc3QgbGVmdEJvcmRlciA9IGV2ZW50LmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBjb25zdCB0b3BCb3JkZXIgPSBldmVudC5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuICAgIGNvbnN0IGJvdHRvbUJvcmRlciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcblxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyIGFyZWFzLlxuICAgIGNvbnN0IHJpZ2h0Q29ybmVyID0gKHdpbmRvdy5vdXRlcldpZHRoIC0gZXZlbnQuY2xpZW50WCkgPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgbGVmdENvcm5lciA9IGV2ZW50LmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgdG9wQ29ybmVyID0gZXZlbnQuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgYm90dG9tQ29ybmVyID0gKHdpbmRvdy5vdXRlckhlaWdodCAtIGV2ZW50LmNsaWVudFkpIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcblxuICAgIGlmICghbGVmdENvcm5lciAmJiAhdG9wQ29ybmVyICYmICFib3R0b21Db3JuZXIgJiYgIXJpZ2h0Q29ybmVyKSB7XG4gICAgICAgIC8vIE9wdGltaXNhdGlvbjogb3V0IG9mIGFsbCBjb3JuZXIgYXJlYXMgaW1wbGllcyBvdXQgb2YgYm9yZGVycy5cbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxuICAgIC8vIERldGVjdCBjb3JuZXJzLlxuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICAvLyBEZXRlY3QgYm9yZGVycy5cbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xuICAgIC8vIE91dCBvZiBib3JkZXIgYXJlYS5cbiAgICBlbHNlIHNldFJlc2l6ZSgpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoU2hvd01ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgQ2FuY2VsbGFibGVQcm9taXNlLCB0eXBlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnMgfSBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcblxuLy8gU2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5jYWxsUmVzdWx0SGFuZGxlciA9IHJlc3VsdEhhbmRsZXI7XG53aW5kb3cuX3dhaWxzLmNhbGxFcnJvckhhbmRsZXIgPSBlcnJvckhhbmRsZXI7XG5cbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8Q2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIiB8IFwib25jYW5jZWxsZWRcIj5cblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2FsbCk7XG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYW5jZWxDYWxsKTtcbmNvbnN0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwPHN0cmluZywgUHJvbWlzZVJlc29sdmVycz4oKTtcblxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xuY29uc3QgQ2FuY2VsTWV0aG9kID0gMFxuXG4vKipcbiAqIEhvbGRzIGFsbCByZXF1aXJlZCBpbmZvcm1hdGlvbiBmb3IgYSBiaW5kaW5nIGNhbGwuXG4gKiBNYXkgcHJvdmlkZSBlaXRoZXIgYSBtZXRob2QgSUQgb3IgYSBtZXRob2QgbmFtZSwgYnV0IG5vdCBib3RoLlxuICovXG5leHBvcnQgdHlwZSBDYWxsT3B0aW9ucyA9IHtcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZElEOiBudW1iZXI7XG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kTmFtZT86IG5ldmVyO1xuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cbiAgICBhcmdzOiBhbnlbXTtcbn0gfCB7XG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2RJRD86IG5ldmVyO1xuICAgIC8qKiBUaGUgZnVsbHkgcXVhbGlmaWVkIG5hbWUgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZE5hbWU6IHN0cmluZztcbiAgICAvKiogQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuICovXG4gICAgYXJnczogYW55W107XG59O1xuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdGhyb3duIGluIGNhc2UgdGhlIGJvdW5kIG1ldGhvZCByZXR1cm5zIGFuIGVycm9yLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgUnVudGltZUVycm9yI25hbWV9IHByb3BlcnR5IGlzIFwiUnVudGltZUVycm9yXCIuXG4gKi9cbmV4cG9ydCBjbGFzcyBSdW50aW1lRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBSdW50aW1lRXJyb3IgaW5zdGFuY2UuXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIlJ1bnRpbWVFcnJvclwiO1xuICAgIH1cbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBjYWxsIHJlcXVlc3QuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cbiAqIEBwYXJhbSBkYXRhIC0gVGhlIHJlc3VsdCBkYXRhIG9mIHRoZSByZXF1ZXN0LlxuICogQHBhcmFtIGlzSlNPTiAtIEluZGljYXRlcyB3aGV0aGVyIHRoZSBkYXRhIGlzIEpTT04gb3Igbm90LlxuICovXG5mdW5jdGlvbiByZXN1bHRIYW5kbGVyKGlkOiBzdHJpbmcsIGRhdGE6IHN0cmluZywgaXNKU09OOiBib29sZWFuKTogdm9pZCB7XG4gICAgY29uc3QgcmVzb2x2ZXJzID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xuICAgIGlmICghcmVzb2x2ZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIWRhdGEpIHtcbiAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUodW5kZWZpbmVkKTtcbiAgICB9IGVsc2UgaWYgKCFpc0pTT04pIHtcbiAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUoZGF0YSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHJlc29sdmVycy5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGNhdGNoIChlcnI6IGFueSkge1xuICAgICAgICAgICAgcmVzb2x2ZXJzLnJlamVjdChuZXcgVHlwZUVycm9yKFwiY291bGQgbm90IHBhcnNlIHJlc3VsdDogXCIgKyBlcnIubWVzc2FnZSwgeyBjYXVzZTogZXJyIH0pKTtcbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBIYW5kbGVzIHRoZSBlcnJvciBmcm9tIGEgY2FsbCByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxuICogQHBhcmFtIGRhdGEgLSBUaGUgZXJyb3IgZGF0YSB0byByZWplY3QgdGhlIHByb21pc2UgaGFuZGxlciB3aXRoLlxuICogQHBhcmFtIGlzSlNPTiAtIEluZGljYXRlcyB3aGV0aGVyIHRoZSBkYXRhIGlzIEpTT04gb3Igbm90LlxuICovXG5mdW5jdGlvbiBlcnJvckhhbmRsZXIoaWQ6IHN0cmluZywgZGF0YTogc3RyaW5nLCBpc0pTT046IGJvb2xlYW4pOiB2b2lkIHtcbiAgICBjb25zdCByZXNvbHZlcnMgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XG4gICAgaWYgKCFyZXNvbHZlcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghaXNKU09OKSB7XG4gICAgICAgIHJlc29sdmVycy5yZWplY3QobmV3IEVycm9yKGRhdGEpKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBsZXQgZXJyb3I6IGFueTtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGVycm9yID0gSlNPTi5wYXJzZShkYXRhKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyOiBhbnkpIHtcbiAgICAgICAgICAgIHJlc29sdmVycy5yZWplY3QobmV3IFR5cGVFcnJvcihcImNvdWxkIG5vdCBwYXJzZSBlcnJvcjogXCIgKyBlcnIubWVzc2FnZSwgeyBjYXVzZTogZXJyIH0pKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBvcHRpb25zOiBFcnJvck9wdGlvbnMgPSB7fTtcbiAgICAgICAgaWYgKGVycm9yLmNhdXNlKSB7XG4gICAgICAgICAgICBvcHRpb25zLmNhdXNlID0gZXJyb3IuY2F1c2U7XG4gICAgICAgIH1cblxuICAgICAgICBsZXQgZXhjZXB0aW9uO1xuICAgICAgICBzd2l0Y2ggKGVycm9yLmtpbmQpIHtcbiAgICAgICAgICAgIGNhc2UgXCJSZWZlcmVuY2VFcnJvclwiOlxuICAgICAgICAgICAgICAgIGV4Y2VwdGlvbiA9IG5ldyBSZWZlcmVuY2VFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgIGNhc2UgXCJUeXBlRXJyb3JcIjpcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgVHlwZUVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgY2FzZSBcIlJ1bnRpbWVFcnJvclwiOlxuICAgICAgICAgICAgICAgIGV4Y2VwdGlvbiA9IG5ldyBSdW50aW1lRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgICAgIGV4Y2VwdGlvbiA9IG5ldyBFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuXG4gICAgICAgIHJlc29sdmVycy5yZWplY3QoZXhjZXB0aW9uKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0cmlldmVzIGFuZCByZW1vdmVzIHRoZSByZXNwb25zZSBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElEIGZyb20gdGhlIGNhbGxSZXNwb25zZXMgbWFwLlxuICpcbiAqIEBwYXJhbSBpZCAtIFRoZSBJRCBvZiB0aGUgcmVzcG9uc2UgdG8gYmUgcmV0cmlldmVkIGFuZCByZW1vdmVkLlxuICogQHJldHVybnMgVGhlIHJlc3BvbnNlIG9iamVjdCBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElELCBpZiBhbnkuXG4gKi9cbmZ1bmN0aW9uIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkOiBzdHJpbmcpOiBQcm9taXNlUmVzb2x2ZXJzIHwgdW5kZWZpbmVkIHtcbiAgICBjb25zdCByZXNwb25zZSA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgcmV0dXJuIHJlc3BvbnNlO1xufVxuXG4vKipcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBJRCB1c2luZyB0aGUgbmFub2lkIGxpYnJhcnkuXG4gKlxuICogQHJldHVybnMgQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgY2FsbFJlc3BvbnNlcyBzZXQuXG4gKi9cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKTogc3RyaW5nIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBDYWxsIGEgYm91bmQgbWV0aG9kIGFjY29yZGluZyB0byB0aGUgZ2l2ZW4gY2FsbCBvcHRpb25zLlxuICpcbiAqIEluIGNhc2Ugb2YgZmFpbHVyZSwgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhbiBleGNlcHRpb25cbiAqIGFtb25nIFJlZmVyZW5jZUVycm9yICh1bmtub3duIG1ldGhvZCksIFR5cGVFcnJvciAod3JvbmcgYXJndW1lbnQgY291bnQgb3IgdHlwZSksXG4gKiB7QGxpbmsgUnVudGltZUVycm9yfSAobWV0aG9kIHJldHVybmVkIGFuIGVycm9yKSwgb3Igb3RoZXIgKG5ldHdvcmsgb3IgaW50ZXJuYWwgZXJyb3JzKS5cbiAqIFRoZSBleGNlcHRpb24gbWlnaHQgaGF2ZSBhIFwiY2F1c2VcIiBmaWVsZCB3aXRoIHRoZSB2YWx1ZSByZXR1cm5lZFxuICogYnkgdGhlIGFwcGxpY2F0aW9uLSBvciBzZXJ2aWNlLWxldmVsIGVycm9yIG1hcnNoYWxpbmcgZnVuY3Rpb25zLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gQSBtZXRob2QgY2FsbCBkZXNjcmlwdG9yLlxuICogQHJldHVybnMgVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9uczogQ2FsbE9wdGlvbnMpOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XG5cbiAgICBjb25zdCByZXN1bHQgPSBDYW5jZWxsYWJsZVByb21pc2Uud2l0aFJlc29sdmVyczxhbnk+KCk7XG4gICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZTogcmVzdWx0LnJlc29sdmUsIHJlamVjdDogcmVzdWx0LnJlamVjdCB9KTtcblxuICAgIGNvbnN0IHJlcXVlc3QgPSBjYWxsKENhbGxCaW5kaW5nLCBPYmplY3QuYXNzaWduKHsgXCJjYWxsLWlkXCI6IGlkIH0sIG9wdGlvbnMpKTtcbiAgICBsZXQgcnVubmluZyA9IGZhbHNlO1xuXG4gICAgcmVxdWVzdC50aGVuKCgpID0+IHtcbiAgICAgICAgcnVubmluZyA9IHRydWU7XG4gICAgfSwgKGVycikgPT4ge1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHJlc3VsdC5yZWplY3QoZXJyKTtcbiAgICB9KTtcblxuICAgIGNvbnN0IGNhbmNlbCA9ICgpID0+IHtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXR1cm4gY2FuY2VsQ2FsbChDYW5jZWxNZXRob2QsIHtcImNhbGwtaWRcIjogaWR9KS5jYXRjaCgoZXJyKSA9PiB7XG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igd2hpbGUgcmVxdWVzdGluZyBiaW5kaW5nIGNhbGwgY2FuY2VsbGF0aW9uOlwiLCBlcnIpO1xuICAgICAgICB9KTtcbiAgICB9O1xuXG4gICAgcmVzdWx0Lm9uY2FuY2VsbGVkID0gKCkgPT4ge1xuICAgICAgICBpZiAocnVubmluZykge1xuICAgICAgICAgICAgcmV0dXJuIGNhbmNlbCgpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcmV0dXJuIHJlcXVlc3QudGhlbihjYW5jZWwpO1xuICAgICAgICB9XG4gICAgfTtcblxuICAgIHJldHVybiByZXN1bHQucHJvbWlzZTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIGJvdW5kIG1ldGhvZCBieSBuYW1lIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKiBTZWUge0BsaW5rIENhbGx9IGZvciBkZXRhaWxzLlxuICpcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG1ldGhvZE5hbWU6IHN0cmluZywgLi4uYXJnczogYW55W10pOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2ROYW1lLCBhcmdzIH0pO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIGJ5IGl0cyBudW1lcmljIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKiBTZWUge0BsaW5rIENhbGx9IGZvciBkZXRhaWxzLlxuICpcbiAqIEBwYXJhbSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybiBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQ6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2RJRCwgYXJncyB9KTtcbn1cbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9pbnNwZWN0LWpzL2lzLWNhbGxhYmxlXG5cbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxuLy9cbi8vIENvcHlyaWdodCAoYykgMjAxNSBKb3JkYW4gSGFyYmFuZFxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHlcbi8vIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWxcbi8vIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHNcbi8vIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGxcbi8vIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpc1xuLy8gZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcbi8vXG4vLyBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXG4vL1xuLy8gVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksXG4vLyBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEVcbi8vIEFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVJcbi8vIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sXG4vLyBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRVxuLy8gU09GVFdBUkUuXG5cbnZhciBmblRvU3RyID0gRnVuY3Rpb24ucHJvdG90eXBlLnRvU3RyaW5nO1xudmFyIHJlZmxlY3RBcHBseTogdHlwZW9mIFJlZmxlY3QuYXBwbHkgfCBmYWxzZSB8IG51bGwgPSB0eXBlb2YgUmVmbGVjdCA9PT0gJ29iamVjdCcgJiYgUmVmbGVjdCAhPT0gbnVsbCAmJiBSZWZsZWN0LmFwcGx5O1xudmFyIGJhZEFycmF5TGlrZTogYW55O1xudmFyIGlzQ2FsbGFibGVNYXJrZXI6IGFueTtcbmlmICh0eXBlb2YgcmVmbGVjdEFwcGx5ID09PSAnZnVuY3Rpb24nICYmIHR5cGVvZiBPYmplY3QuZGVmaW5lUHJvcGVydHkgPT09ICdmdW5jdGlvbicpIHtcbiAgICB0cnkge1xuICAgICAgICBiYWRBcnJheUxpa2UgPSBPYmplY3QuZGVmaW5lUHJvcGVydHkoe30sICdsZW5ndGgnLCB7XG4gICAgICAgICAgICBnZXQ6IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgICAgICB0aHJvdyBpc0NhbGxhYmxlTWFya2VyO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgaXNDYWxsYWJsZU1hcmtlciA9IHt9O1xuICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmUgbm8tdGhyb3ctbGl0ZXJhbFxuICAgICAgICByZWZsZWN0QXBwbHkoZnVuY3Rpb24gKCkgeyB0aHJvdyA0MjsgfSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcbiAgICB9IGNhdGNoIChfKSB7XG4gICAgICAgIGlmIChfICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7XG4gICAgICAgICAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xuICAgICAgICB9XG4gICAgfVxufSBlbHNlIHtcbiAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xufVxuXG52YXIgY29uc3RydWN0b3JSZWdleCA9IC9eXFxzKmNsYXNzXFxiLztcbnZhciBpc0VTNkNsYXNzRm4gPSBmdW5jdGlvbiBpc0VTNkNsYXNzRnVuY3Rpb24odmFsdWU6IGFueSk6IGJvb2xlYW4ge1xuICAgIHRyeSB7XG4gICAgICAgIHZhciBmblN0ciA9IGZuVG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgIHJldHVybiBjb25zdHJ1Y3RvclJlZ2V4LnRlc3QoZm5TdHIpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlOyAvLyBub3QgYSBmdW5jdGlvblxuICAgIH1cbn07XG5cbnZhciB0cnlGdW5jdGlvbk9iamVjdCA9IGZ1bmN0aW9uIHRyeUZ1bmN0aW9uVG9TdHIodmFsdWU6IGFueSk6IGJvb2xlYW4ge1xuICAgIHRyeSB7XG4gICAgICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgICAgICBmblRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG59O1xudmFyIHRvU3RyID0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZztcbnZhciBvYmplY3RDbGFzcyA9ICdbb2JqZWN0IE9iamVjdF0nO1xudmFyIGZuQ2xhc3MgPSAnW29iamVjdCBGdW5jdGlvbl0nO1xudmFyIGdlbkNsYXNzID0gJ1tvYmplY3QgR2VuZXJhdG9yRnVuY3Rpb25dJztcbnZhciBkZGFDbGFzcyA9ICdbb2JqZWN0IEhUTUxBbGxDb2xsZWN0aW9uXSc7IC8vIElFIDExXG52YXIgZGRhQ2xhc3MyID0gJ1tvYmplY3QgSFRNTCBkb2N1bWVudC5hbGwgY2xhc3NdJztcbnZhciBkZGFDbGFzczMgPSAnW29iamVjdCBIVE1MQ29sbGVjdGlvbl0nOyAvLyBJRSA5LTEwXG52YXIgaGFzVG9TdHJpbmdUYWcgPSB0eXBlb2YgU3ltYm9sID09PSAnZnVuY3Rpb24nICYmICEhU3ltYm9sLnRvU3RyaW5nVGFnOyAvLyBiZXR0ZXI6IHVzZSBgaGFzLXRvc3RyaW5ndGFnYFxuXG52YXIgaXNJRTY4ID0gISgwIGluIFssXSk7IC8vIGVzbGludC1kaXNhYmxlLWxpbmUgbm8tc3BhcnNlLWFycmF5cywgY29tbWEtc3BhY2luZ1xuXG52YXIgaXNEREE6ICh2YWx1ZTogYW55KSA9PiBib29sZWFuID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCgpIHsgcmV0dXJuIGZhbHNlOyB9O1xuaWYgKHR5cGVvZiBkb2N1bWVudCA9PT0gJ29iamVjdCcpIHtcbiAgICAvLyBGaXJlZm94IDMgY2Fub25pY2FsaXplcyBEREEgdG8gdW5kZWZpbmVkIHdoZW4gaXQncyBub3QgYWNjZXNzZWQgZGlyZWN0bHlcbiAgICB2YXIgYWxsID0gZG9jdW1lbnQuYWxsO1xuICAgIGlmICh0b1N0ci5jYWxsKGFsbCkgPT09IHRvU3RyLmNhbGwoZG9jdW1lbnQuYWxsKSkge1xuICAgICAgICBpc0REQSA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwodmFsdWUpIHtcbiAgICAgICAgICAgIC8qIGdsb2JhbHMgZG9jdW1lbnQ6IGZhbHNlICovXG4gICAgICAgICAgICAvLyBpbiBJRSA2LTgsIHR5cGVvZiBkb2N1bWVudC5hbGwgaXMgXCJvYmplY3RcIiBhbmQgaXQncyB0cnV0aHlcbiAgICAgICAgICAgIGlmICgoaXNJRTY4IHx8ICF2YWx1ZSkgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3VuZGVmaW5lZCcgfHwgdHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JykpIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICB2YXIgc3RyID0gdG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybiAoXG4gICAgICAgICAgICAgICAgICAgICAgICBzdHIgPT09IGRkYUNsYXNzXG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IGRkYUNsYXNzMlxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczMgLy8gb3BlcmEgMTIuMTZcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gb2JqZWN0Q2xhc3MgLy8gSUUgNi04XG4gICAgICAgICAgICAgICAgICAgICkgJiYgdmFsdWUoJycpID09IG51bGw7IC8vIGVzbGludC1kaXNhYmxlLWxpbmUgZXFlcWVxXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZSkgeyAvKiovIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGlzQ2FsbGFibGVSZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSAge1xuICAgIGlmIChpc0REQSh2YWx1ZSkpIHsgcmV0dXJuIHRydWU7IH1cbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB0cnkge1xuICAgICAgICAocmVmbGVjdEFwcGx5IGFzIGFueSkodmFsdWUsIG51bGwsIGJhZEFycmF5TGlrZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBpZiAoZSAhPT0gaXNDYWxsYWJsZU1hcmtlcikgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB9XG4gICAgcmV0dXJuICFpc0VTNkNsYXNzRm4odmFsdWUpICYmIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcbn1cblxuZnVuY3Rpb24gaXNDYWxsYWJsZU5vUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkge1xuICAgIGlmIChpc0REQSh2YWx1ZSkpIHsgcmV0dXJuIHRydWU7IH1cbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAoaGFzVG9TdHJpbmdUYWcpIHsgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTsgfVxuICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHZhciBzdHJDbGFzcyA9IHRvU3RyLmNhbGwodmFsdWUpO1xuICAgIGlmIChzdHJDbGFzcyAhPT0gZm5DbGFzcyAmJiBzdHJDbGFzcyAhPT0gZ2VuQ2xhc3MgJiYgISgvXlxcW29iamVjdCBIVE1MLykudGVzdChzdHJDbGFzcykpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcbn07XG5cbmV4cG9ydCBkZWZhdWx0IHJlZmxlY3RBcHBseSA/IGlzQ2FsbGFibGVSZWZBcHBseSA6IGlzQ2FsbGFibGVOb1JlZkFwcGx5O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgaXNDYWxsYWJsZSBmcm9tIFwiLi9jYWxsYWJsZS5qc1wiO1xuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdXNlZCBhcyByZWplY3Rpb24gcmVhc29uXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gaXMgY2FuY2VsbGVkIHN1Y2Nlc3NmdWxseS5cbiAqXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbEVycm9yXCJgLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSBjYXVzZSBwYXNzZWQgdG8gdGhlIGNhbmNlbCBtZXRob2QsIGlmIGFueS5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbEVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbEVycm9yYCBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsRXJyb3JcIjtcbiAgICB9XG59XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gcmVqZWN0cyBhZnRlciBiZWluZyBjYW5jZWxsZWQsXG4gKiBvciB3aGVuIHRoZSBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHRocm93cyBvciByZWplY3RzLlxuICpcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcImAuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIHJlYXNvbiB0aGUgcHJvbWlzZSByZWplY3RlZCB3aXRoLlxuICpcbiAqIEJlY2F1c2UgdGhlIG9yaWdpbmFsIHByb21pc2Ugd2FzIGNhbmNlbGxlZCxcbiAqIGEgd3JhcHBlciBwcm9taXNlIHdpbGwgYmUgcGFzc2VkIHRvIHRoZSB1bmhhbmRsZWQgcmVqZWN0aW9uIGxpc3RlbmVyIGluc3RlYWQuXG4gKiBUaGUge0BsaW5rIHByb21pc2V9IHByb3BlcnR5IGhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBvcmlnaW5hbCBwcm9taXNlLlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogSG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIHByb21pc2UgdGhhdCB3YXMgY2FuY2VsbGVkIGFuZCB0aGVuIHJlamVjdGVkLlxuICAgICAqL1xuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPjtcblxuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbGxlZFJlamVjdGlvbkVycm9yYCBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gcHJvbWlzZSAtIFRoZSBwcm9taXNlIHRoYXQgY2F1c2VkIHRoZSBlcnJvciBvcmlnaW5hbGx5LlxuICAgICAqIEBwYXJhbSByZWFzb24gLSBUaGUgcmVqZWN0aW9uIHJlYXNvbi5cbiAgICAgKiBAcGFyYW0gaW5mbyAtIEFuIG9wdGlvbmFsIGluZm9ybWF0aXZlIG1lc3NhZ2Ugc3BlY2lmeWluZyB0aGUgY2lyY3Vtc3RhbmNlcyBpbiB3aGljaCB0aGUgZXJyb3Igd2FzIHRocm93bi5cbiAgICAgKiAgICAgICAgICAgICAgIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgYFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cImAuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCByZWFzb24/OiBhbnksIGluZm8/OiBzdHJpbmcpIHtcbiAgICAgICAgc3VwZXIoKGluZm8gPz8gXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiKSArIFwiIFJlYXNvbjogXCIgKyBlcnJvck1lc3NhZ2UocmVhc29uKSwgeyBjYXVzZTogcmVhc29uIH0pO1xuICAgICAgICB0aGlzLnByb21pc2UgPSBwcm9taXNlO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCI7XG4gICAgfVxufVxuXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+ID0gKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciA9IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciA9IChjYXVzZT86IGFueSkgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiA9IChyZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiwgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcikgPT4gdm9pZDtcblxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwpOiBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxIHwgVFJlc3VsdDI+O1xuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbn1cblxuLyoqXG4gKiBXcmFwcyBhIGNhbmNlbGxhYmxlIHByb21pc2UgYWxvbmcgd2l0aCBpdHMgcmVzb2x1dGlvbiBtZXRob2RzLlxuICogVGhlIGBvbmNhbmNlbGxlZGAgZmllbGQgd2lsbCBiZSBudWxsIGluaXRpYWxseSBidXQgbWF5IGJlIHNldCB0byBwcm92aWRlIGEgY3VzdG9tIGNhbmNlbGxhdGlvbiBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XG4gICAgcmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD47XG4gICAgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcjtcbiAgICBvbmNhbmNlbGxlZDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyIHwgbnVsbDtcbn1cblxuaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlIHtcbiAgICByZWFkb25seSByb290OiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZTtcbiAgICByZXNvbHZpbmc6IGJvb2xlYW47XG4gICAgc2V0dGxlZDogYm9vbGVhbjtcbiAgICByZWFzb24/OiBDYW5jZWxFcnJvcjtcbn1cblxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cbmNvbnN0IGJhcnJpZXJTeW0gPSBTeW1ib2woXCJiYXJyaWVyXCIpO1xuY29uc3QgY2FuY2VsSW1wbFN5bSA9IFN5bWJvbChcImNhbmNlbEltcGxcIik7XG5jb25zdCBzcGVjaWVzID0gU3ltYm9sLnNwZWNpZXMgPz8gU3ltYm9sKFwic3BlY2llc1BvbHlmaWxsXCIpO1xuXG4vKipcbiAqIEEgcHJvbWlzZSB3aXRoIGFuIGF0dGFjaGVkIG1ldGhvZCBmb3IgY2FuY2VsbGluZyBsb25nLXJ1bm5pbmcgb3BlcmF0aW9ucyAoc2VlIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2UjY2FuY2VsfSkuXG4gKiBDYW5jZWxsYXRpb24gY2FuIG9wdGlvbmFsbHkgYmUgYm91bmQgdG8gYW4ge0BsaW5rIEFib3J0U2lnbmFsfVxuICogZm9yIGJldHRlciBjb21wb3NhYmlsaXR5IChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWxPbn0pLlxuICpcbiAqIENhbmNlbGxpbmcgYSBwZW5kaW5nIHByb21pc2Ugd2lsbCByZXN1bHQgaW4gYW4gaW1tZWRpYXRlIHJlamVjdGlvblxuICogd2l0aCBhbiBpbnN0YW5jZSBvZiB7QGxpbmsgQ2FuY2VsRXJyb3J9IGFzIHJlYXNvbixcbiAqIGJ1dCB3aG9ldmVyIHN0YXJ0ZWQgdGhlIHByb21pc2Ugd2lsbCBiZSByZXNwb25zaWJsZVxuICogZm9yIGFjdHVhbGx5IGFib3J0aW5nIHRoZSB1bmRlcmx5aW5nIG9wZXJhdGlvbi5cbiAqIFRvIHRoaXMgcHVycG9zZSwgdGhlIGNvbnN0cnVjdG9yIGFuZCBhbGwgY2hhaW5pbmcgbWV0aG9kc1xuICogYWNjZXB0IG9wdGlvbmFsIGNhbmNlbGxhdGlvbiBjYWxsYmFja3MuXG4gKlxuICogSWYgYSBgQ2FuY2VsbGFibGVQcm9taXNlYCBzdGlsbCByZXNvbHZlcyBhZnRlciBoYXZpbmcgYmVlbiBjYW5jZWxsZWQsXG4gKiB0aGUgcmVzdWx0IHdpbGwgYmUgZGlzY2FyZGVkLiBJZiBpdCByZWplY3RzLCB0aGUgcmVhc29uXG4gKiB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24sXG4gKiB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBpbnN0YW5jZS5cbiAqIFRvIGZhY2lsaXRhdGUgdGhlIGhhbmRsaW5nIG9mIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyxcbiAqIGNhbmNlbGxlZCBgQ2FuY2VsbGFibGVQcm9taXNlYHMgd2lsbCBfbm90XyByZXBvcnQgdW5oYW5kbGVkIGBDYW5jZWxFcnJvcmBzXG4gKiB3aG9zZSBgY2F1c2VgIGZpZWxkIGlzIHRoZSBzYW1lIGFzIHRoZSBvbmUgd2l0aCB3aGljaCB0aGUgY3VycmVudCBwcm9taXNlIHdhcyBjYW5jZWxsZWQuXG4gKlxuICogQWxsIHVzdWFsIHByb21pc2UgbWV0aG9kcyBhcmUgZGVmaW5lZCBhbmQgcmV0dXJuIGEgYENhbmNlbGxhYmxlUHJvbWlzZWBcbiAqIHdob3NlIGNhbmNlbCBtZXRob2Qgd2lsbCBjYW5jZWwgdGhlIHBhcmVudCBvcGVyYXRpb24gYXMgd2VsbCwgcHJvcGFnYXRpbmcgdGhlIGNhbmNlbGxhdGlvbiByZWFzb25cbiAqIHVwd2FyZHMgdGhyb3VnaCBwcm9taXNlIGNoYWlucy5cbiAqIENvbnZlcnNlbHksIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpbGwgbm90IGF1dG9tYXRpY2FsbHkgY2FuY2VsIGRlcGVuZGVudCBwcm9taXNlcyBkb3duc3RyZWFtOlxuICogYGBgdHNcbiAqIGxldCByb290ID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDEgPSByb290LnRoZW4oKCkgPT4geyAuLi4gfSk7XG4gKiBsZXQgY2hpbGQyID0gY2hpbGQxLnRoZW4oKCkgPT4geyAuLi4gfSk7XG4gKiBsZXQgY2hpbGQzID0gcm9vdC5jYXRjaCgoKSA9PiB7IC4uLiB9KTtcbiAqIGNoaWxkMS5jYW5jZWwoKTsgLy8gQ2FuY2VscyBjaGlsZDEgYW5kIHJvb3QsIGJ1dCBub3QgY2hpbGQyIG9yIGNoaWxkM1xuICogYGBgXG4gKiBDYW5jZWxsaW5nIGEgcHJvbWlzZSB0aGF0IGhhcyBhbHJlYWR5IHNldHRsZWQgaXMgc2FmZSBhbmQgaGFzIG5vIGNvbnNlcXVlbmNlLlxuICpcbiAqIFRoZSBgY2FuY2VsYCBtZXRob2QgcmV0dXJucyBhIHByb21pc2UgdGhhdCBfYWx3YXlzIGZ1bGZpbGxzX1xuICogYWZ0ZXIgdGhlIHdob2xlIGNoYWluIGhhcyBwcm9jZXNzZWQgdGhlIGNhbmNlbCByZXF1ZXN0XG4gKiBhbmQgYWxsIGF0dGFjaGVkIGNhbGxiYWNrcyB1cCB0byB0aGF0IG1vbWVudCBoYXZlIHJ1bi5cbiAqXG4gKiBBbGwgRVMyMDI0IHByb21pc2UgbWV0aG9kcyAoc3RhdGljIGFuZCBpbnN0YW5jZSkgYXJlIGRlZmluZWQgb24gQ2FuY2VsbGFibGVQcm9taXNlLFxuICogYnV0IGFjdHVhbCBhdmFpbGFiaWxpdHkgbWF5IHZhcnkgd2l0aCBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gKlxuICogSW4gbGluZSB3aXRoIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nLFxuICogYENhbmNlbGxhYmxlUHJvbWlzZWAgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy5cbiAqIEV4dGVuZGVycyBzaG91bGQgdGFrZSBjYXJlIHRvIHByb3ZpZGUgdGhlaXIgb3duIG1ldGhvZCBpbXBsZW1lbnRhdGlvbnMuXG4gKiBUaGlzIG1pZ2h0IGJlIHJlY29uc2lkZXJlZCBpbiBjYXNlIHRoZSBwcm9wb3NhbCBpcyByZXRpcmVkLlxuICpcbiAqIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBhIHdyYXBwZXIgYXJvdW5kIHRoZSBET00gUHJvbWlzZSBvYmplY3RcbiAqIGFuZCBpcyBjb21wbGlhbnQgd2l0aCB0aGUgW1Byb21pc2VzL0ErIHNwZWNpZmljYXRpb25dKGh0dHBzOi8vcHJvbWlzZXNhcGx1cy5jb20vKVxuICogKGl0IHBhc3NlcyB0aGUgW2NvbXBsaWFuY2Ugc3VpdGVdKGh0dHBzOi8vZ2l0aHViLmNvbS9wcm9taXNlcy1hcGx1cy9wcm9taXNlcy10ZXN0cykpXG4gKiBpZiBzbyBpcyB0aGUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbi5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbGxhYmxlUHJvbWlzZTxUPiBleHRlbmRzIFByb21pc2U8VD4gaW1wbGVtZW50cyBQcm9taXNlTGlrZTxUPiwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIHByaXZhdGUgW2JhcnJpZXJTeW1dITogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gfCBudWxsO1xuICAgIC8qKiBAaW50ZXJuYWwgKi9cbiAgICBwcml2YXRlIHJlYWRvbmx5IFtjYW5jZWxJbXBsU3ltXSE6IChyZWFzb246IENhbmNlbEVycm9yKSA9PiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGBDYW5jZWxsYWJsZVByb21pc2VgLlxuICAgICAqXG4gICAgICogQHBhcmFtIGV4ZWN1dG9yIC0gQSBjYWxsYmFjayB1c2VkIHRvIGluaXRpYWxpemUgdGhlIHByb21pc2UuIFRoaXMgY2FsbGJhY2sgaXMgcGFzc2VkIHR3byBhcmd1bWVudHM6XG4gICAgICogICAgICAgICAgICAgICAgICAgYSBgcmVzb2x2ZWAgY2FsbGJhY2sgdXNlZCB0byByZXNvbHZlIHRoZSBwcm9taXNlIHdpdGggYSB2YWx1ZVxuICAgICAqICAgICAgICAgICAgICAgICAgIG9yIHRoZSByZXN1bHQgb2YgYW5vdGhlciBwcm9taXNlIChwb3NzaWJseSBjYW5jZWxsYWJsZSksXG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGEgYHJlamVjdGAgY2FsbGJhY2sgdXNlZCB0byByZWplY3QgdGhlIHByb21pc2Ugd2l0aCBhIHByb3ZpZGVkIHJlYXNvbiBvciBlcnJvci5cbiAgICAgKiAgICAgICAgICAgICAgICAgICBJZiB0aGUgdmFsdWUgcHJvdmlkZWQgdG8gdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqICAgICAgICAgICAgICAgICAgIChpdCBoYXMgYSBgdGhlbmAgX2FuZF8gYSBgY2FuY2VsYCBtZXRob2QpLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGZvcndhcmRlZCB0byB0aGF0IG9iamVjdCBhbmQgdGhlIG9uY2FuY2VsbGVkIHdpbGwgbm90IGJlIGludm9rZWQgYW55bW9yZS5cbiAgICAgKiAgICAgICAgICAgICAgICAgICBJZiBhbnkgb25lIG9mIHRoZSB0d28gY2FsbGJhY2tzIGlzIGNhbGxlZCBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZCxcbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGUgcHJvdmlkZWQgdmFsdWVzIHdpbGwgYmUgY2FuY2VsbGVkIGFuZCByZXNvbHZlZCBhcyB1c3VhbCxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBidXQgdGhlaXIgcmVzdWx0cyB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiAgICAgICAgICAgICAgICAgICBIb3dldmVyLCBpZiB0aGUgcmVzb2x1dGlvbiBwcm9jZXNzIHVsdGltYXRlbHkgZW5kcyB1cCBpbiBhIHJlamVjdGlvblxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoYXQgaXMgbm90IGR1ZSB0byBjYW5jZWxsYXRpb24sIHRoZSByZWplY3Rpb24gcmVhc29uXG4gICAgICogICAgICAgICAgICAgICAgICAgd2lsbCBiZSB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfVxuICAgICAqICAgICAgICAgICAgICAgICAgIGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICogQHBhcmFtIG9uY2FuY2VsbGVkIC0gSXQgaXMgdGhlIGNhbGxlcidzIHJlc3BvbnNpYmlsaXR5IHRvIGVuc3VyZSB0aGF0IGFueSBvcGVyYXRpb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBzdGFydGVkIGJ5IHRoZSBleGVjdXRvciBpcyBwcm9wZXJseSBoYWx0ZWQgdXBvbiBjYW5jZWxsYXRpb24uXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVGhpcyBvcHRpb25hbCBjYWxsYmFjayBjYW4gYmUgdXNlZCB0byB0aGF0IHB1cnBvc2UuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSXQgd2lsbCBiZSBjYWxsZWQgX3N5bmNocm9ub3VzbHlfIHdpdGggYSBjYW5jZWxsYXRpb24gY2F1c2VcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aGVuIGNhbmNlbGxhdGlvbiBpcyByZXF1ZXN0ZWQsIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgcmVqZWN0ZWRcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSwgYnV0IF9iZWZvcmVfXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW55IHtAbGluayB0aGVufS97QGxpbmsgY2F0Y2h9L3tAbGluayBmaW5hbGx5fSBjYWxsYmFjayBydW5zLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH1cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aWxsIG9ubHkgZnVsZmlsbCBhZnRlciB0aGUgZm9ybWVyIGhhcyBzZXR0bGVkLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFVuaGFuZGxlZCBleGNlcHRpb25zIG9yIHJlamVjdGlvbnMgZnJvbSB0aGUgY2FsbGJhY2sgd2lsbCBiZSB3cmFwcGVkXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIHVuaGFuZGxlZCByZWplY3Rpb25zLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgY2FsbGVkIGJlZm9yZSBjYW5jZWxsYXRpb24gd2l0aCBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoaXMgcHJvbWlzZSB3aWxsIGJlIGRpdmVydGVkIHRvIHRoYXQgcHJvbWlzZSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbmQgdGhlIG9yaWdpbmFsIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXhlY3V0b3I6IENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+LCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcikge1xuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcbiAgICAgICAgc3VwZXIoKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XG5cbiAgICAgICAgaWYgKCh0aGlzLmNvbnN0cnVjdG9yIGFzIGFueSlbc3BlY2llc10gIT09IFByb21pc2UpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy4gUGxlYXNlIHJlZnJhaW4gZnJvbSBvdmVycmlkaW5nIHRoZSBbU3ltYm9sLnNwZWNpZXNdIHN0YXRpYyBwcm9wZXJ0eS5cIik7XG4gICAgICAgIH1cblxuICAgICAgICBsZXQgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHtcbiAgICAgICAgICAgIHByb21pc2U6IHRoaXMsXG4gICAgICAgICAgICByZXNvbHZlLFxuICAgICAgICAgICAgcmVqZWN0LFxuICAgICAgICAgICAgZ2V0IG9uY2FuY2VsbGVkKCkgeyByZXR1cm4gb25jYW5jZWxsZWQgPz8gbnVsbDsgfSxcbiAgICAgICAgICAgIHNldCBvbmNhbmNlbGxlZChjYikgeyBvbmNhbmNlbGxlZCA9IGNiID8/IHVuZGVmaW5lZDsgfVxuICAgICAgICB9O1xuXG4gICAgICAgIGNvbnN0IHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcbiAgICAgICAgICAgIGdldCByb290KCkgeyByZXR1cm4gc3RhdGU7IH0sXG4gICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxuICAgICAgICAgICAgc2V0dGxlZDogZmFsc2VcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBTZXR1cCBjYW5jZWxsYXRpb24gc3lzdGVtLlxuICAgICAgICB2b2lkIE9iamVjdC5kZWZpbmVQcm9wZXJ0aWVzKHRoaXMsIHtcbiAgICAgICAgICAgIFtiYXJyaWVyU3ltXToge1xuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IHRydWUsXG4gICAgICAgICAgICAgICAgdmFsdWU6IG51bGxcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICBbY2FuY2VsSW1wbFN5bV06IHtcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB2YWx1ZTogY2FuY2VsbGVyRm9yKHByb21pc2UsIHN0YXRlKVxuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcblxuICAgICAgICAvLyBSdW4gdGhlIGFjdHVhbCBleGVjdXRvci5cbiAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBzdGF0ZSk7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBleGVjdXRvcihyZXNvbHZlckZvcihwcm9taXNlLCBzdGF0ZSksIHJlamVjdG9yKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7XG4gICAgICAgICAgICAgICAgY29uc29sZS5sb2coXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIENhbmNlbGxhYmxlUHJvbWlzZSBleGVjdXRvci5cIiwgZXJyKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENhbmNlbHMgaW1tZWRpYXRlbHkgdGhlIGV4ZWN1dGlvbiBvZiB0aGUgb3BlcmF0aW9uIGFzc29jaWF0ZWQgd2l0aCB0aGlzIHByb21pc2UuXG4gICAgICogVGhlIHByb21pc2UgcmVqZWN0cyB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSBpbnN0YW5jZSBhcyByZWFzb24sXG4gICAgICogd2l0aCB0aGUge0BsaW5rIENhbmNlbEVycm9yI2NhdXNlfSBwcm9wZXJ0eSBzZXQgdG8gdGhlIGdpdmVuIGFyZ3VtZW50LCBpZiBhbnkuXG4gICAgICpcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCBhZnRlciB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkO1xuICAgICAqIHJlcGVhdGVkIGNhbGxzIGluIHBhcnRpY3VsYXIgYXJlIHNhZmUsIGJ1dCBvbmx5IHRoZSBmaXJzdCBvbmVcbiAgICAgKiB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxuICAgICAqXG4gICAgICogVGhlIGBDYW5jZWxFcnJvcmAgZXhjZXB0aW9uIF9uZWVkIG5vdF8gYmUgaGFuZGxlZCBleHBsaWNpdGx5IF9vbiB0aGUgcHJvbWlzZXMgdGhhdCBhcmUgYmVpbmcgY2FuY2VsbGVkOl9cbiAgICAgKiBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aXRoIG5vIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyIGRvZXMgbm90IHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudC5cbiAgICAgKiBUaGVyZWZvcmUsIHRoZSBmb2xsb3dpbmcgaWRpb21zIGFyZSBhbGwgZXF1YWxseSBjb3JyZWN0OlxuICAgICAqIGBgYHRzXG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS5jYW5jZWwoKTtcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYW5jZWwoKTtcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYXRjaCguLi4pLmNhbmNlbCgpO1xuICAgICAqIGBgYFxuICAgICAqIFdoZW5ldmVyIHNvbWUgY2FuY2VsbGVkIHByb21pc2UgaW4gYSBjaGFpbiByZWplY3RzIHdpdGggYSBgQ2FuY2VsRXJyb3JgXG4gICAgICogd2l0aCB0aGUgc2FtZSBjYW5jZWxsYXRpb24gY2F1c2UgYXMgaXRzZWxmLCB0aGUgZXJyb3Igd2lsbCBiZSBkaXNjYXJkZWQgc2lsZW50bHkuXG4gICAgICogSG93ZXZlciwgdGhlIGBDYW5jZWxFcnJvcmAgX3dpbGwgc3RpbGwgYmUgZGVsaXZlcmVkXyB0byBhbGwgYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXJzXG4gICAgICogYWRkZWQgYnkge0BsaW5rIHRoZW59IGFuZCByZWxhdGVkIG1ldGhvZHM6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICAgICAqIGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkuY2F0Y2goY29uc29sZS5sb2cpO1xuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBBIENhbmNlbEVycm9yIGlzIHByaW50ZWQgdG8gdGhlIGNvbnNvbGUuXG4gICAgICogYGBgXG4gICAgICogSWYgdGhlIGBDYW5jZWxFcnJvcmAgaXMgbm90IGhhbmRsZWQgZG93bnN0cmVhbSBieSB0aGUgdGltZSBpdCByZWFjaGVzXG4gICAgICogYSBfbm9uLWNhbmNlbGxlZF8gcHJvbWlzZSwgaXQgX3dpbGxfIHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCxcbiAgICAgKiBqdXN0IGxpa2Ugbm9ybWFsIHJlamVjdGlvbnMgd291bGQ6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICAgICAqIGxldCBjaGFpbmVkID0gY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS50aGVuKCgpID0+IHsgLi4uIH0pOyAvLyBObyBjYXRjaC4uLlxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBVbmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50IG9uIGNoYWluZWQhXG4gICAgICogYGBgXG4gICAgICogVGhlcmVmb3JlLCBpdCBpcyBpbXBvcnRhbnQgdG8gZWl0aGVyIGNhbmNlbCB3aG9sZSBwcm9taXNlIGNoYWlucyBmcm9tIHRoZWlyIHRhaWwsXG4gICAgICogYXMgc2hvd24gaW4gdGhlIGNvcnJlY3QgaWRpb21zIGFib3ZlLCBvciB0YWtlIGNhcmUgb2YgaGFuZGxpbmcgZXJyb3JzIGV2ZXJ5d2hlcmUuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBBIGNhbmNlbGxhYmxlIHByb21pc2UgdGhhdCBfZnVsZmlsbHNfIGFmdGVyIHRoZSBjYW5jZWwgY2FsbGJhY2sgKGlmIGFueSlcbiAgICAgKiBhbmQgYWxsIGhhbmRsZXJzIGF0dGFjaGVkIHVwIHRvIHRoZSBjYWxsIHRvIGNhbmNlbCBoYXZlIHJ1bi5cbiAgICAgKiBJZiB0aGUgY2FuY2VsIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgYnkgYGNhbmNlbGBcbiAgICAgKiB3aWxsIGFsc28gd2FpdCBmb3IgdGhhdCB0aGVuYWJsZSB0byBzZXR0bGUuXG4gICAgICogVGhpcyBlbmFibGVzIGNhbGxlcnMgdG8gd2FpdCBmb3IgdGhlIGNhbmNlbGxlZCBvcGVyYXRpb24gdG8gdGVybWluYXRlXG4gICAgICogd2l0aG91dCBiZWluZyBmb3JjZWQgdG8gaGFuZGxlIHBvdGVudGlhbCBlcnJvcnMgYXQgdGhlIGNhbGwgc2l0ZS5cbiAgICAgKiBgYGB0c1xuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpLnRoZW4oKCkgPT4ge1xuICAgICAqICAgICAvLyBDbGVhbnVwIGZpbmlzaGVkLCBpdCdzIHNhZmUgdG8gZG8gc29tZXRoaW5nIGVsc2UuXG4gICAgICogfSwgKGVycikgPT4ge1xuICAgICAqICAgICAvLyBVbnJlYWNoYWJsZTogdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSBjYW5jZWwgd2lsbCBuZXZlciByZWplY3QuXG4gICAgICogfSk7XG4gICAgICogYGBgXG4gICAgICogTm90ZSB0aGF0IHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgX25vdF8gaGFuZGxlIGltcGxpY2l0bHkgYW55IHJlamVjdGlvblxuICAgICAqIHRoYXQgbWlnaHQgaGF2ZSBvY2N1cnJlZCBhbHJlYWR5IGluIHRoZSBjYW5jZWxsZWQgY2hhaW4uXG4gICAgICogSXQgd2lsbCBqdXN0IHRyYWNrIHdoZXRoZXIgcmVnaXN0ZXJlZCBoYW5kbGVycyBoYXZlIGJlZW4gZXhlY3V0ZWQgb3Igbm90LlxuICAgICAqIFRoZXJlZm9yZSwgdW5oYW5kbGVkIHJlamVjdGlvbnMgd2lsbCBuZXZlciBiZSBzaWxlbnRseSBoYW5kbGVkIGJ5IGNhbGxpbmcgY2FuY2VsLlxuICAgICAqL1xuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICAvLyBJTlZBUklBTlQ6IHRoZSByZXN1bHQgb2YgdGhpc1tjYW5jZWxJbXBsU3ltXSBhbmQgdGhlIGJhcnJpZXIgZG8gbm90IGV2ZXIgcmVqZWN0LlxuICAgICAgICAgICAgLy8gVW5mb3J0dW5hdGVseSBtYWNPUyBIaWdoIFNpZXJyYSBkb2VzIG5vdCBzdXBwb3J0IFByb21pc2UuYWxsU2V0dGxlZC5cbiAgICAgICAgICAgIFByb21pc2UuYWxsKFtcbiAgICAgICAgICAgICAgICB0aGlzW2NhbmNlbEltcGxTeW1dKG5ldyBDYW5jZWxFcnJvcihcIlByb21pc2UgY2FuY2VsbGVkLlwiLCB7IGNhdXNlIH0pKSxcbiAgICAgICAgICAgICAgICBjdXJyZW50QmFycmllcih0aGlzKVxuICAgICAgICAgICAgXSkudGhlbigoKSA9PiByZXNvbHZlKCksICgpID0+IHJlc29sdmUoKSk7XG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEJpbmRzIHByb21pc2UgY2FuY2VsbGF0aW9uIHRvIHRoZSBhYm9ydCBldmVudCBvZiB0aGUgZ2l2ZW4ge0BsaW5rIEFib3J0U2lnbmFsfS5cbiAgICAgKiBJZiB0aGUgc2lnbmFsIGhhcyBhbHJlYWR5IGFib3J0ZWQsIHRoZSBwcm9taXNlIHdpbGwgYmUgY2FuY2VsbGVkIGltbWVkaWF0ZWx5LlxuICAgICAqIFdoZW4gZWl0aGVyIGNvbmRpdGlvbiBpcyB2ZXJpZmllZCwgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZSB3aWxsIGJlIHNldFxuICAgICAqIHRvIHRoZSBzaWduYWwncyBhYm9ydCByZWFzb24gKHNlZSB7QGxpbmsgQWJvcnRTaWduYWwjcmVhc29ufSkuXG4gICAgICpcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCAob3IgaWYgdGhlIHNpZ25hbCBhYm9ydHMpIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZC5cbiAgICAgKiBPbmx5IHRoZSBmaXJzdCBzaWduYWwgdG8gYWJvcnQgd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cbiAgICAgKlxuICAgICAqIEZvciBtb3JlIGRldGFpbHMgYWJvdXQgdGhlIGNhbmNlbGxhdGlvbiBwcm9jZXNzLFxuICAgICAqIHNlZSB7QGxpbmsgY2FuY2VsfSBhbmQgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgIGNvbnN0cnVjdG9yLlxuICAgICAqXG4gICAgICogVGhpcyBtZXRob2QgZW5hYmxlcyBgYXdhaXRgaW5nIGNhbmNlbGxhYmxlIHByb21pc2VzIHdpdGhvdXQgaGF2aW5nXG4gICAgICogdG8gc3RvcmUgdGhlbSBmb3IgZnV0dXJlIGNhbmNlbGxhdGlvbiwgZS5nLjpcbiAgICAgKiBgYGB0c1xuICAgICAqIGF3YWl0IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCkuY2FuY2VsT24oc2lnbmFsKTtcbiAgICAgKiBgYGBcbiAgICAgKiBpbnN0ZWFkIG9mOlxuICAgICAqIGBgYHRzXG4gICAgICogbGV0IHByb21pc2VUb0JlQ2FuY2VsbGVkID0gbG9uZ1J1bm5pbmdPcGVyYXRpb24oKTtcbiAgICAgKiBhd2FpdCBwcm9taXNlVG9CZUNhbmNlbGxlZDtcbiAgICAgKiBgYGBcbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoaXMgcHJvbWlzZSwgZm9yIG1ldGhvZCBjaGFpbmluZy5cbiAgICAgKi9cbiAgICBjYW5jZWxPbihzaWduYWw6IEFib3J0U2lnbmFsKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgaWYgKHNpZ25hbC5hYm9ydGVkKSB7XG4gICAgICAgICAgICB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pXG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzaWduYWwuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pLCB7Y2FwdHVyZTogdHJ1ZX0pO1xuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHRoaXM7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgY2FsbGJhY2tzIGZvciB0aGUgcmVzb2x1dGlvbiBhbmQvb3IgcmVqZWN0aW9uIG9mIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYC5cbiAgICAgKlxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICpcbiAgICAgKiBAcGFyYW0gb25mdWxmaWxsZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZXNvbHZlZC5cbiAgICAgKiBAcGFyYW0gb25yZWplY3RlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlamVjdGVkLlxuICAgICAqIEByZXR1cm5zIEEgYENhbmNlbGxhYmxlUHJvbWlzZWAgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHdoaWNoZXZlciBjYWxsYmFjayBpcyBleGVjdXRlZC5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBob29rZWQgdXAgdG8gcHJvcGFnYXRlIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB1cCB0aGUgY2hhaW4sIGJ1dCBub3QgZG93bjpcbiAgICAgKlxuICAgICAqICAgLSBpZiB0aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBiZSBpbnZva2VkIHdpdGggYSBgQ2FuY2VsRXJyb3JgXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSBfd2lsbCByZXNvbHZlIHJlZ3VsYXJseV8gd2l0aCBpdHMgcmVzdWx0O1xuICAgICAqICAgLSBjb252ZXJzZWx5LCBpZiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIF90aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkIHRvbztfXG4gICAgICogICAgIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIHN0aWxsIGJlIGludm9rZWQgd2l0aCB0aGUgcGFyZW50J3MgYENhbmNlbEVycm9yYCxcbiAgICAgKiAgICAgYnV0IGl0cyByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWRcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYSBgQ2FuY2VsRXJyb3JgIGFzIHdlbGwuXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgZWl0aGVyIGNhbGxiYWNrIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICovXG4gICAgdGhlbjxUUmVzdWx0MSA9IFQsIFRSZXN1bHQyID0gbmV2ZXI+KG9uZnVsZmlsbGVkPzogKCh2YWx1ZTogVCkgPT4gVFJlc3VsdDEgfCBQcm9taXNlTGlrZTxUUmVzdWx0MT4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gVFJlc3VsdDIgfCBQcm9taXNlTGlrZTxUUmVzdWx0Mj4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQyPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUUmVzdWx0MSB8IFRSZXN1bHQyPiB7XG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS50aGVuIGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBOT1RFOiBUeXBlU2NyaXB0J3MgYnVpbHQtaW4gdHlwZSBmb3IgdGhlbiBpcyBicm9rZW4sXG4gICAgICAgIC8vIGFzIGl0IGFsbG93cyBzcGVjaWZ5aW5nIGFuIGFyYml0cmFyeSBUUmVzdWx0MSAhPSBUIGV2ZW4gd2hlbiBvbmZ1bGZpbGxlZCBpcyBub3QgYSBmdW5jdGlvbi5cbiAgICAgICAgLy8gV2UgY2Fubm90IGZpeCBpdCBpZiB3ZSB3YW50IHRvIENhbmNlbGxhYmxlUHJvbWlzZSB0byBpbXBsZW1lbnQgUHJvbWlzZUxpa2U8VD4uXG5cbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZnVsZmlsbGVkKSkgeyBvbmZ1bGZpbGxlZCA9IGlkZW50aXR5IGFzIGFueTsgfVxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25yZWplY3RlZCkpIHsgb25yZWplY3RlZCA9IHRocm93ZXI7IH1cblxuICAgICAgICBpZiAob25mdWxmaWxsZWQgPT09IGlkZW50aXR5ICYmIG9ucmVqZWN0ZWQgPT0gdGhyb3dlcikge1xuICAgICAgICAgICAgLy8gU2hvcnRjdXQgZm9yIHRyaXZpYWwgYXJndW1lbnRzLlxuICAgICAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUpID0+IHJlc29sdmUodGhpcyBhcyBhbnkpKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnN0IGJhcnJpZXI6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0ge307XG4gICAgICAgIHRoaXNbYmFycmllclN5bV0gPSBiYXJyaWVyO1xuXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgIHZvaWQgc3VwZXIudGhlbihcbiAgICAgICAgICAgICAgICAodmFsdWUpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcblxuICAgICAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbmZ1bGZpbGxlZCEodmFsdWUpKTtcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0sXG4gICAgICAgICAgICAgICAgKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcblxuICAgICAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbnJlamVjdGVkIShyZWFzb24pKTtcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICk7XG4gICAgICAgIH0sIGFzeW5jIChjYXVzZT8pID0+IHtcbiAgICAgICAgICAgIC8vY2FuY2VsbGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG9uY2FuY2VsbGVkPy4oY2F1c2UpO1xuICAgICAgICAgICAgfSBmaW5hbGx5IHtcbiAgICAgICAgICAgICAgICBhd2FpdCB0aGlzLmNhbmNlbChjYXVzZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEF0dGFjaGVzIGEgY2FsbGJhY2sgZm9yIG9ubHkgdGhlIHJlamVjdGlvbiBvZiB0aGUgUHJvbWlzZS5cbiAgICAgKlxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICpcbiAgICAgKiBJdCBpcyBlcXVpdmFsZW50IHRvXG4gICAgICogYGBgdHNcbiAgICAgKiBjYW5jZWxsYWJsZVByb21pc2UudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcbiAgICAgKiBgYGBcbiAgICAgKiBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxuICAgICAqXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cbiAgICAgKlxuICAgICAqIElmIGBvbnJlamVjdGVkYCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cbiAgICAgKi9cbiAgICBjYXRjaDxUUmVzdWx0ID0gbmV2ZXI+KG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiAoUHJvbWlzZUxpa2U8VFJlc3VsdD4gfCBUUmVzdWx0KSkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUIHwgVFJlc3VsdD4ge1xuICAgICAgICByZXR1cm4gdGhpcy50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEF0dGFjaGVzIGEgY2FsbGJhY2sgdGhhdCBpcyBpbnZva2VkIHdoZW4gdGhlIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBzZXR0bGVkIChmdWxmaWxsZWQgb3IgcmVqZWN0ZWQpLiBUaGVcbiAgICAgKiByZXNvbHZlZCB2YWx1ZSBjYW5ub3QgYmUgYWNjZXNzZWQgb3IgbW9kaWZpZWQgZnJvbSB0aGUgY2FsbGJhY2suXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBzZXR0bGUgaW4gdGhlIHNhbWUgc3RhdGUgYXMgdGhlIG9yaWdpbmFsIG9uZVxuICAgICAqIGFmdGVyIHRoZSBwcm92aWRlZCBjYWxsYmFjayBoYXMgY29tcGxldGVkIGV4ZWN1dGlvbixcbiAgICAgKiB1bmxlc3MgdGhlIGNhbGxiYWNrIHRocm93cyBvciByZXR1cm5zIGEgcmVqZWN0aW5nIHByb21pc2UsXG4gICAgICogaW4gd2hpY2ggY2FzZSB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCBhcyB3ZWxsLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBPbmNlIHRoZSBwYXJlbnQgcHJvbWlzZSBzZXR0bGVzLCB0aGUgYG9uZmluYWxseWAgY2FsbGJhY2sgd2lsbCBydW4sXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBpbXBsZW1lbnRlZCBpbiB0ZXJtcyBvZiB7QGxpbmsgdGhlbn0gYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXG4gICAgICogSXQgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxuICAgICAqXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cbiAgICAgKlxuICAgICAqIElmIGBvbmZpbmFsbHlgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxuICAgICAqL1xuICAgIGZpbmFsbHkob25maW5hbGx5PzogKCgpID0+IHZvaWQpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUuZmluYWxseSBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZmluYWxseSkpIHtcbiAgICAgICAgICAgIHJldHVybiB0aGlzLnRoZW4ob25maW5hbGx5LCBvbmZpbmFsbHksIG9uY2FuY2VsbGVkKTtcbiAgICAgICAgfVxuXG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4oXG4gICAgICAgICAgICAodmFsdWUpID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHZhbHVlKSxcbiAgICAgICAgICAgIChyZWFzb24/KSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB7IHRocm93IHJlYXNvbjsgfSksXG4gICAgICAgICAgICBvbmNhbmNlbGxlZCxcbiAgICAgICAgKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBXZSB1c2UgdGhlIGBbU3ltYm9sLnNwZWNpZXNdYCBzdGF0aWMgcHJvcGVydHksIGlmIGF2YWlsYWJsZSxcbiAgICAgKiB0byBkaXNhYmxlIHRoZSBidWlsdC1pbiBhdXRvbWF0aWMgc3ViY2xhc3NpbmcgZmVhdHVyZXMgZnJvbSB7QGxpbmsgUHJvbWlzZX0uXG4gICAgICogSXQgaXMgY3JpdGljYWwgZm9yIHBlcmZvcm1hbmNlIHJlYXNvbnMgdGhhdCBleHRlbmRlcnMgZG8gbm90IG92ZXJyaWRlIHRoaXMuXG4gICAgICogT25jZSB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZ1xuICAgICAqIGlzIGVpdGhlciBhY2NlcHRlZCBvciByZXRpcmVkLCB0aGlzIGltcGxlbWVudGF0aW9uIHdpbGwgaGF2ZSB0byBiZSByZXZpc2VkIGFjY29yZGluZ2x5LlxuICAgICAqXG4gICAgICogQGlnbm9yZVxuICAgICAqIEBpbnRlcm5hbFxuICAgICAqL1xuICAgIHN0YXRpYyBnZXQgW3NwZWNpZXNdKCkge1xuICAgICAgICByZXR1cm4gUHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUsIG9yIHJlamVjdGVkIHdoZW4gYW55IFByb21pc2UgaXMgcmVqZWN0ZWQuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFsbDxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPltdPjtcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBBd2FpdGVkPFRbUF0+OyB9PjtcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlIG9yIHJlamVjdC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8UHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUPj5bXT47XG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VFtQXT4+OyB9PjtcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbFNldHRsZWQoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUaGUgYW55IGZ1bmN0aW9uIHJldHVybnMgYSBwcm9taXNlIHRoYXQgaXMgZnVsZmlsbGVkIGJ5IHRoZSBmaXJzdCBnaXZlbiBwcm9taXNlIHRvIGJlIGZ1bGZpbGxlZCxcbiAgICAgKiBvciByZWplY3RlZCB3aXRoIGFuIEFnZ3JlZ2F0ZUVycm9yIGNvbnRhaW5pbmcgYW4gYXJyYXkgb2YgcmVqZWN0aW9uIHJlYXNvbnNcbiAgICAgKiBpZiBhbGwgb2YgdGhlIGdpdmVuIHByb21pc2VzIGFyZSByZWplY3RlZC5cbiAgICAgKiBJdCByZXNvbHZlcyBhbGwgZWxlbWVudHMgb2YgdGhlIHBhc3NlZCBpdGVyYWJsZSB0byBwcm9taXNlcyBhcyBpdCBydW5zIHRoaXMgYWxnb3JpdGhtLlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbnk8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUW251bWJlcl0+PjtcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYW55KGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIFByb21pc2UgdGhhdCBpcyByZXNvbHZlZCBvciByZWplY3RlZCB3aGVuIGFueSBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgYXJlIHJlc29sdmVkIG9yIHJlamVjdGVkLlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByYWNlPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJhY2UoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgY2FuY2VsbGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGNhbmNlbDxUID0gbmV2ZXI+KGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgY29uc3QgcCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xuICAgICAgICBwLmNhbmNlbChjYXVzZSk7XG4gICAgICAgIHJldHVybiBwO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgY2FuY2Vsc1xuICAgICAqIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgY2F1c2UuXG4gICAgICpcbiAgICAgKiBJZiB0aGUge0BsaW5rIEFib3J0U2lnbmFsLnRpbWVvdXR9IGZhY3RvcnkgbWV0aG9kIGlzIGF2YWlsYWJsZSxcbiAgICAgKiBpdCBpcyB1c2VkIHRvIGJhc2UgdGhlIHRpbWVvdXQgb24gX2FjdGl2ZV8gdGltZSByYXRoZXIgdGhhbiBfZWxhcHNlZF8gdGltZS5cbiAgICAgKiBPdGhlcndpc2UsIGB0aW1lb3V0YCBmYWxscyBiYWNrIHRvIHtAbGluayBzZXRUaW1lb3V0fS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyB0aW1lb3V0PFQgPSBuZXZlcj4obWlsbGlzZWNvbmRzOiBudW1iZXIsIGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xuICAgICAgICBpZiAoQWJvcnRTaWduYWwgJiYgdHlwZW9mIEFib3J0U2lnbmFsID09PSAnZnVuY3Rpb24nICYmIEFib3J0U2lnbmFsLnRpbWVvdXQgJiYgdHlwZW9mIEFib3J0U2lnbmFsLnRpbWVvdXQgPT09ICdmdW5jdGlvbicpIHtcbiAgICAgICAgICAgIEFib3J0U2lnbmFsLnRpbWVvdXQobWlsbGlzZWNvbmRzKS5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSksIG1pbGxpc2Vjb25kcyk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQuXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBzbGVlcChtaWxsaXNlY29uZHM6IG51bWJlcik6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyXG4gICAgICogdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCB2YWx1ZS5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHNsZWVwPFQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcbiAgICBzdGF0aWMgc2xlZXA8VCA9IHZvaWQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZT86IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiByZXNvbHZlKHZhbHVlISksIG1pbGxpc2Vjb25kcyk7XG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVqZWN0ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgcmVhc29uLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlamVjdDxUID0gbmV2ZXI+KHJlYXNvbj86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChfLCByZWplY3QpID0+IHJlamVjdChyZWFzb24pKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZXNvbHZlKCk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIHJlc29sdmU8VCA9IHZvaWQ+KHZhbHVlPzogVCB8IFByb21pc2VMaWtlPFQ+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+IHtcbiAgICAgICAgaWYgKHZhbHVlIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSB7XG4gICAgICAgICAgICAvLyBPcHRpbWlzZSBmb3IgY2FuY2VsbGFibGUgcHJvbWlzZXMuXG4gICAgICAgICAgICByZXR1cm4gdmFsdWU7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8YW55PigocmVzb2x2ZSkgPT4gcmVzb2x2ZSh2YWx1ZSkpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIGFuZCByZXR1cm5zIGl0IGluIGFuIG9iamVjdCwgYWxvbmcgd2l0aCBpdHMgcmVzb2x2ZSBhbmQgcmVqZWN0IGZ1bmN0aW9uc1xuICAgICAqIGFuZCBhIGdldHRlci9zZXR0ZXIgZm9yIHRoZSBjYW5jZWxsYXRpb24gY2FsbGJhY2suXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHdpdGhSZXNvbHZlcnM8VD4oKTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XG4gICAgICAgIGxldCByZXN1bHQ6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7IG9uY2FuY2VsbGVkOiBudWxsIH0gYXMgYW55O1xuICAgICAgICByZXN1bHQucHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgcmVzdWx0LnJlc29sdmUgPSByZXNvbHZlO1xuICAgICAgICAgICAgcmVzdWx0LnJlamVjdCA9IHJlamVjdDtcbiAgICAgICAgfSwgKGNhdXNlPzogYW55KSA9PiB7IHJlc3VsdC5vbmNhbmNlbGxlZD8uKGNhdXNlKTsgfSk7XG4gICAgICAgIHJldHVybiByZXN1bHQ7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIGNhbmNlbGxhdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgcmVzdWx0aW5nIGZ1bmN0aW9uIGRvZXMgbm90IHJlamVjdC5cbiAqL1xuZnVuY3Rpb24gY2FuY2VsbGVyRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSkge1xuICAgIGxldCBjYW5jZWxsYXRpb25Qcm9taXNlOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPSB1bmRlZmluZWQ7XG5cbiAgICByZXR1cm4gKHJlYXNvbjogQ2FuY2VsRXJyb3IpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPT4ge1xuICAgICAgICBpZiAoIXN0YXRlLnNldHRsZWQpIHtcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgc3RhdGUucmVhc29uID0gcmVhc29uO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcblxuICAgICAgICAgICAgLy8gQXR0YWNoIGFuIGVycm9yIGhhbmRsZXIgdGhhdCBpZ25vcmVzIHRoaXMgc3BlY2lmaWMgcmVqZWN0aW9uIHJlYXNvbiBhbmQgbm90aGluZyBlbHNlLlxuICAgICAgICAgICAgLy8gSW4gdGhlb3J5LCBhIHNhbmUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbiBhdCB0aGlzIHBvaW50XG4gICAgICAgICAgICAvLyBzaG91bGQgYWx3YXlzIHJlamVjdCB3aXRoIG91ciBjYW5jZWxsYXRpb24gcmVhc29uLFxuICAgICAgICAgICAgLy8gaGVuY2UgdGhlIGhhbmRsZXIgd2lsbCBuZXZlciB0aHJvdy5cbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5wcm90b3R5cGUudGhlbi5jYWxsKHByb21pc2UucHJvbWlzZSwgdW5kZWZpbmVkLCAoZXJyKSA9PiB7XG4gICAgICAgICAgICAgICAgaWYgKGVyciAhPT0gcmVhc29uKSB7XG4gICAgICAgICAgICAgICAgICAgIHRocm93IGVycjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9KTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIElmIHJlYXNvbiBpcyBub3Qgc2V0LCB0aGUgcHJvbWlzZSByZXNvbHZlZCByZWd1bGFybHksIGhlbmNlIHdlIG11c3Qgbm90IGNhbGwgb25jYW5jZWxsZWQuXG4gICAgICAgIC8vIElmIG9uY2FuY2VsbGVkIGlzIHVuc2V0LCBubyBuZWVkIHRvIGdvIGFueSBmdXJ0aGVyLlxuICAgICAgICBpZiAoIXN0YXRlLnJlYXNvbiB8fCAhcHJvbWlzZS5vbmNhbmNlbGxlZCkgeyByZXR1cm47IH1cblxuICAgICAgICBjYW5jZWxsYXRpb25Qcm9taXNlID0gbmV3IFByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgcmVzb2x2ZShwcm9taXNlLm9uY2FuY2VsbGVkIShzdGF0ZS5yZWFzb24hLmNhdXNlKSk7XG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KS5jYXRjaCgocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcbiAgICAgICAgfSk7XG5cbiAgICAgICAgLy8gVW5zZXQgb25jYW5jZWxsZWQgdG8gcHJldmVudCByZXBlYXRlZCBjYWxscy5cbiAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG51bGw7XG5cbiAgICAgICAgcmV0dXJuIGNhbmNlbGxhdGlvblByb21pc2U7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlc29sdXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cbiAqL1xuZnVuY3Rpb24gcmVzb2x2ZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4ge1xuICAgIHJldHVybiAodmFsdWUpID0+IHtcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcblxuICAgICAgICBpZiAodmFsdWUgPT09IHByb21pc2UucHJvbWlzZSkge1xuICAgICAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJBIHByb21pc2UgY2Fubm90IGJlIHJlc29sdmVkIHdpdGggaXRzZWxmLlwiKSk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cblxuICAgICAgICBpZiAodmFsdWUgIT0gbnVsbCAmJiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyB8fCB0eXBlb2YgdmFsdWUgPT09ICdmdW5jdGlvbicpKSB7XG4gICAgICAgICAgICBsZXQgdGhlbjogYW55O1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICB0aGVuID0gKHZhbHVlIGFzIGFueSkudGhlbjtcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgICAgIHByb21pc2UucmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICBpZiAoaXNDYWxsYWJsZSh0aGVuKSkge1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIGxldCBjYW5jZWwgPSAodmFsdWUgYXMgYW55KS5jYW5jZWw7XG4gICAgICAgICAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKGNhbmNlbCkpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IG9uY2FuY2VsbGVkID0gKGNhdXNlPzogYW55KSA9PiB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIH07XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoc3RhdGUucmVhc29uKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gSWYgYWxyZWFkeSBjYW5jZWxsZWQsIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24uXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgY2FuY2VsbGVyIGFsZ29yaXRobSBkb2VzIG5vdCByZWplY3RcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBzbyBpdCBjYW4gYmUgZGlzY2FyZGVkIHNhZmVseS5cbiAgICAgICAgICAgICAgICAgICAgICAgICAgICB2b2lkIGNhbmNlbGxlckZvcih7IC4uLnByb21pc2UsIG9uY2FuY2VsbGVkIH0sIHN0YXRlKShzdGF0ZS5yZWFzb24pO1xuICAgICAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gb25jYW5jZWxsZWQ7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9IGNhdGNoIHt9XG5cbiAgICAgICAgICAgICAgICBjb25zdCBuZXdTdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XG4gICAgICAgICAgICAgICAgICAgIHJvb3Q6IHN0YXRlLnJvb3QsXG4gICAgICAgICAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXG4gICAgICAgICAgICAgICAgICAgIGdldCBzZXR0bGVkKCkgeyByZXR1cm4gdGhpcy5yb290LnNldHRsZWQgfSxcbiAgICAgICAgICAgICAgICAgICAgc2V0IHNldHRsZWQodmFsdWUpIHsgdGhpcy5yb290LnNldHRsZWQgPSB2YWx1ZTsgfSxcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHJlYXNvbigpIHsgcmV0dXJuIHRoaXMucm9vdC5yZWFzb24gfVxuICAgICAgICAgICAgICAgIH07XG5cbiAgICAgICAgICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIG5ld1N0YXRlKTtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KHRoZW4sIHZhbHVlLCBbcmVzb2x2ZXJGb3IocHJvbWlzZSwgbmV3U3RhdGUpLCByZWplY3Rvcl0pO1xuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICByZXR1cm47IC8vIElNUE9SVEFOVCFcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuXG4gICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgcHJvbWlzZS5yZXNvbHZlKHZhbHVlKTtcbiAgICB9O1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlamVjdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICovXG5mdW5jdGlvbiByZWplY3RvckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciB7XG4gICAgcmV0dXJuIChyZWFzb24/KSA9PiB7XG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnJlc29sdmluZyA9IHRydWU7XG5cbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgaWYgKHJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIHN0YXRlLnJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIE9iamVjdC5pcyhyZWFzb24uY2F1c2UsIHN0YXRlLnJlYXNvbi5jYXVzZSkpIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gU3dhbGxvdyBsYXRlIHJlamVjdGlvbnMgdGhhdCBhcmUgQ2FuY2VsRXJyb3JzIHdob3NlIGNhbmNlbGxhdGlvbiBjYXVzZSBpcyB0aGUgc2FtZSBhcyBvdXJzLlxuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSBjYXRjaCB7fVxuXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbikpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIENhbmNlbHMgYWxsIHZhbHVlcyBpbiBhbiBhcnJheSB0aGF0IGxvb2sgbGlrZSBjYW5jZWxsYWJsZSB0aGVuYWJsZXMuXG4gKiBSZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGZ1bGZpbGxzIG9uY2UgYWxsIGNhbmNlbGxhdGlvbiBwcm9jZWR1cmVzIGZvciB0aGUgZ2l2ZW4gdmFsdWVzIGhhdmUgc2V0dGxlZC5cbiAqL1xuZnVuY3Rpb24gY2FuY2VsQWxsKHBhcmVudDogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCB2YWx1ZXM6IGFueVtdLCBjYXVzZT86IGFueSk6IFByb21pc2U8dm9pZD4ge1xuICAgIGNvbnN0IHJlc3VsdHMgPSBbXTtcblxuICAgIGZvciAoY29uc3QgdmFsdWUgb2YgdmFsdWVzKSB7XG4gICAgICAgIGxldCBjYW5jZWw6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcjtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZSh2YWx1ZS50aGVuKSkgeyBjb250aW51ZTsgfVxuICAgICAgICAgICAgY2FuY2VsID0gdmFsdWUuY2FuY2VsO1xuICAgICAgICAgICAgaWYgKCFpc0NhbGxhYmxlKGNhbmNlbCkpIHsgY29udGludWU7IH1cbiAgICAgICAgfSBjYXRjaCB7IGNvbnRpbnVlOyB9XG5cbiAgICAgICAgbGV0IHJlc3VsdDogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgcmVzdWx0ID0gUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XG4gICAgICAgICAgICBjb250aW51ZTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghcmVzdWx0KSB7IGNvbnRpbnVlOyB9XG4gICAgICAgIHJlc3VsdHMucHVzaChcbiAgICAgICAgICAgIChyZXN1bHQgaW5zdGFuY2VvZiBQcm9taXNlICA/IHJlc3VsdCA6IFByb21pc2UucmVzb2x2ZShyZXN1bHQpKS5jYXRjaCgocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICk7XG4gICAgfVxuXG4gICAgcmV0dXJuIFByb21pc2UuYWxsKHJlc3VsdHMpIGFzIGFueTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIGl0cyBhcmd1bWVudC5cbiAqL1xuZnVuY3Rpb24gaWRlbnRpdHk8VD4oeDogVCk6IFQge1xuICAgIHJldHVybiB4O1xufVxuXG4vKipcbiAqIFRocm93cyBpdHMgYXJndW1lbnQuXG4gKi9cbmZ1bmN0aW9uIHRocm93ZXIocmVhc29uPzogYW55KTogbmV2ZXIge1xuICAgIHRocm93IHJlYXNvbjtcbn1cblxuLyoqXG4gKiBBdHRlbXB0cyB2YXJpb3VzIHN0cmF0ZWdpZXMgdG8gY29udmVydCBhbiBlcnJvciB0byBhIHN0cmluZy5cbiAqL1xuZnVuY3Rpb24gZXJyb3JNZXNzYWdlKGVycjogYW55KTogc3RyaW5nIHtcbiAgICB0cnkge1xuICAgICAgICBpZiAoZXJyIGluc3RhbmNlb2YgRXJyb3IgfHwgdHlwZW9mIGVyciAhPT0gJ29iamVjdCcgfHwgZXJyLnRvU3RyaW5nICE9PSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nKSB7XG4gICAgICAgICAgICByZXR1cm4gXCJcIiArIGVycjtcbiAgICAgICAgfVxuICAgIH0gY2F0Y2gge31cblxuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiBKU09OLnN0cmluZ2lmeShlcnIpO1xuICAgIH0gY2F0Y2gge31cblxuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nLmNhbGwoZXJyKTtcbiAgICB9IGNhdGNoIHt9XG5cbiAgICByZXR1cm4gXCI8Y291bGQgbm90IGNvbnZlcnQgZXJyb3IgdG8gc3RyaW5nPlwiO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYmFycmllciBwcm9taXNlIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS4gSWYgbmVjZXNzYXJ5LCBpbml0aWFsaXNlcyB0aGUgYmFycmllci5cbiAqL1xuZnVuY3Rpb24gY3VycmVudEJhcnJpZXI8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+KTogUHJvbWlzZTx2b2lkPiB7XG4gICAgbGV0IHB3cjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSBwcm9taXNlW2JhcnJpZXJTeW1dID8/IHt9O1xuICAgIGlmICghKCdwcm9taXNlJyBpbiBwd3IpKSB7XG4gICAgICAgIE9iamVjdC5hc3NpZ24ocHdyLCBwcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPigpKTtcbiAgICB9XG4gICAgaWYgKHByb21pc2VbYmFycmllclN5bV0gPT0gbnVsbCkge1xuICAgICAgICBwd3IucmVzb2x2ZSEoKTtcbiAgICAgICAgcHJvbWlzZVtiYXJyaWVyU3ltXSA9IHB3cjtcbiAgICB9XG4gICAgcmV0dXJuIHB3ci5wcm9taXNlITtcbn1cblxuLy8gUG9seWZpbGwgUHJvbWlzZS53aXRoUmVzb2x2ZXJzLlxubGV0IHByb21pc2VXaXRoUmVzb2x2ZXJzID0gUHJvbWlzZS53aXRoUmVzb2x2ZXJzO1xuaWYgKHByb21pc2VXaXRoUmVzb2x2ZXJzICYmIHR5cGVvZiBwcm9taXNlV2l0aFJlc29sdmVycyA9PT0gJ2Z1bmN0aW9uJykge1xuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gcHJvbWlzZVdpdGhSZXNvbHZlcnMuYmluZChQcm9taXNlKTtcbn0gZWxzZSB7XG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBmdW5jdGlvbiA8VD4oKTogUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBQcm9taXNlPFQ+KChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xuICAgICAgICByZXR1cm4geyBwcm9taXNlLCByZXNvbHZlLCByZWplY3QgfTtcbiAgICB9XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNsaXBib2FyZCk7XG5cbmNvbnN0IENsaXBib2FyZFNldFRleHQgPSAwO1xuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XG5cbi8qKlxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxuICpcbiAqIEBwYXJhbSB0ZXh0IC0gVGhlIHRleHQgdG8gYmUgc2V0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKiBAcmV0dXJuIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0OiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRTZXRUZXh0LCB7dGV4dH0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSB0ZXh0IGZyb20gdGhlIENsaXBib2FyZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKTogUHJvbWlzZTxzdHJpbmc+IHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBBbnkgaXMgYSBkdW1teSBjcmVhdGlvbiBmdW5jdGlvbiBmb3Igc2ltcGxlIG9yIHVua25vd24gdHlwZXMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBBbnk8VCA9IGFueT4oc291cmNlOiBhbnkpOiBUIHtcbiAgICByZXR1cm4gc291cmNlO1xufVxuXG4vKipcbiAqIEJ5dGVTbGljZSBpcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIHRoYXQgcmVwbGFjZXNcbiAqIG51bGwgc3RyaW5ncyB3aXRoIGVtcHR5IHN0cmluZ3MuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCeXRlU2xpY2Uoc291cmNlOiBhbnkpOiBzdHJpbmcge1xuICAgIHJldHVybiAoKHNvdXJjZSA9PSBudWxsKSA/IFwiXCIgOiBzb3VyY2UpO1xufVxuXG4vKipcbiAqIEFycmF5IHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJyYXlcbiAqIHdob3NlIGVsZW1lbnRzIGFyZSBvZiB0aGF0IHR5cGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBBcnJheTxUID0gYW55PihlbGVtZW50OiAoc291cmNlOiBhbnkpID0+IFQpOiAoc291cmNlOiBhbnkpID0+IFRbXSB7XG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IFtdIDogc291cmNlKTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm4gW107XG4gICAgICAgIH1cbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzb3VyY2UubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgICAgIHNvdXJjZVtpXSA9IGVsZW1lbnQoc291cmNlW2ldKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogTWFwIHRha2VzIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdHdvIGFyYml0cmFyeSB0eXBlc1xuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIG9iamVjdFxuICogd2hvc2Uga2V5cyBhbmQgdmFsdWVzIGFyZSBvZiB0aG9zZSB0eXBlcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE1hcDxWID0gYW55PihrZXk6IChzb3VyY2U6IGFueSkgPT4gc3RyaW5nLCB2YWx1ZTogKHNvdXJjZTogYW55KSA9PiBWKTogKHNvdXJjZTogYW55KSA9PiBSZWNvcmQ8c3RyaW5nLCBWPiB7XG4gICAgaWYgKHZhbHVlID09PSBBbnkpIHtcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyB7fSA6IHNvdXJjZSk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xuICAgICAgICAgICAgcmV0dXJuIHt9O1xuICAgICAgICB9XG4gICAgICAgIGZvciAoY29uc3Qga2V5IGluIHNvdXJjZSkge1xuICAgICAgICAgICAgc291cmNlW2tleV0gPSB2YWx1ZShzb3VyY2Vba2V5XSk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcbiAgICB9O1xufVxuXG4vKipcbiAqIE51bGxhYmxlIHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXG4gKiBhbmQgcmV0dXJucyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhIG51bGxhYmxlIHZhbHVlIG9mIHRoYXQgdHlwZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE51bGxhYmxlPFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gKFQgfCBudWxsKSB7XG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gQW55O1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gbnVsbCA6IGVsZW1lbnQoc291cmNlKSk7XG59XG5cbi8qKlxuICogU3RydWN0IHRha2VzIGFuIG9iamVjdCBtYXBwaW5nIGZpZWxkIG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9uc1xuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgc3RydWN0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU3RydWN0KGNyZWF0ZUZpZWxkOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4pOlxuICAgIDxVIGV4dGVuZHMgUmVjb3JkPHN0cmluZywgYW55PiA9IGFueT4oc291cmNlOiBhbnkpID0+IFVcbntcbiAgICBsZXQgYWxsQW55ID0gdHJ1ZTtcbiAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcbiAgICAgICAgaWYgKGNyZWF0ZUZpZWxkW25hbWVdICE9PSBBbnkpIHtcbiAgICAgICAgICAgIGFsbEFueSA9IGZhbHNlO1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cbiAgICB9XG4gICAgaWYgKGFsbEFueSkge1xuICAgICAgICByZXR1cm4gQW55O1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICAgICAgaWYgKG5hbWUgaW4gc291cmNlKSB7XG4gICAgICAgICAgICAgICAgc291cmNlW25hbWVdID0gY3JlYXRlRmllbGRbbmFtZV0oc291cmNlW25hbWVdKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmV4cG9ydCBpbnRlcmZhY2UgU2l6ZSB7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXG4gICAgV2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXG4gICAgSGVpZ2h0OiBudW1iZXI7XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgUmVjdCB7XG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cbiAgICBYOiBudW1iZXI7XG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cbiAgICBZOiBudW1iZXI7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xuICAgIFdpZHRoOiBudW1iZXI7XG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHJlY3RhbmdsZS4gKi9cbiAgICBIZWlnaHQ6IG51bWJlcjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBTY3JlZW4ge1xuICAgIC8qKiBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi4gKi9cbiAgICBJRDogc3RyaW5nO1xuICAgIC8qKiBIdW1hbi1yZWFkYWJsZSBuYW1lIG9mIHRoZSBzY3JlZW4uICovXG4gICAgTmFtZTogc3RyaW5nO1xuICAgIC8qKiBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLiAqL1xuICAgIFNjYWxlRmFjdG9yOiBudW1iZXI7XG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBYOiBudW1iZXI7XG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBZOiBudW1iZXI7XG4gICAgLyoqIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uICovXG4gICAgU2l6ZTogU2l6ZTtcbiAgICAvKiogQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LiAqL1xuICAgIEJvdW5kczogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0IChiZWZvcmUgc2NhbGluZykuICovXG4gICAgUGh5c2ljYWxCb3VuZHM6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLiAqL1xuICAgIFdvcmtBcmVhOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgV29ya0FyZWEgb2YgdGhlIHNjcmVlbiAoYmVmb3JlIHNjYWxpbmcpLiAqL1xuICAgIFBoeXNpY2FsV29ya0FyZWE6IFJlY3Q7XG4gICAgLyoqIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xuICAgIElzUHJpbWFyeTogYm9vbGVhbjtcbiAgICAvKiogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uICovXG4gICAgUm90YXRpb246IG51bWJlcjtcbn1cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TY3JlZW5zKTtcblxuY29uc3QgZ2V0QWxsID0gMDtcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xuY29uc3QgZ2V0Q3VycmVudCA9IDI7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpOiBQcm9taXNlPFNjcmVlbltdPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcbn1cbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBOzs7QUM2QkEsSUFBTSxjQUNGO0FBRUcsU0FBUyxPQUFPLE9BQWUsSUFBWTtBQUM5QyxNQUFJLEtBQUs7QUFFVCxNQUFJLElBQUksT0FBTztBQUNmLFNBQU8sS0FBSztBQUVSLFVBQU0sWUFBYSxLQUFLLE9BQU8sSUFBSSxLQUFNLENBQUM7QUFBQSxFQUM5QztBQUNBLFNBQU87QUFDWDs7O0FDN0JBLElBQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUdyQyxJQUFNLGNBQWMsT0FBTyxPQUFPO0FBQUEsRUFDckMsTUFBTTtBQUFBLEVBQ04sV0FBVztBQUFBLEVBQ1gsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsWUFBWTtBQUNoQixDQUFDO0FBQ00sSUFBSSxXQUFXLE9BQU87QUFTdEIsU0FBUyxpQkFBaUIsUUFBZ0IsYUFBcUIsSUFBSTtBQUN0RSxTQUFPLFNBQVUsUUFBZ0IsT0FBWSxNQUFNO0FBQy9DLFdBQU8sa0JBQWtCLFFBQVEsUUFBUSxZQUFZLElBQUk7QUFBQSxFQUM3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LEVBQ1g7QUFDSjtBQU1PLFNBQVMsU0FBa0IsU0FBMEQ7QUFDeEYsTUFBSSxZQUFZLEtBQUs7QUFDakIsV0FBTztBQUFBLEVBQ1g7QUFFQSxTQUFPLENBQUMsV0FBWSxXQUFXLE9BQU8sT0FBTyxRQUFRLE1BQU07QUFDL0Q7QUFNTyxTQUFTLE9BQU8sYUFFdkI7QUFDSSxNQUFJLFNBQVM7QUFDYixhQUFXLFFBQVEsYUFBYTtBQUM1QixRQUFJLFlBQVksSUFBSSxNQUFNLEtBQUs7QUFDM0IsZUFBUztBQUNUO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFDQSxNQUFJLFFBQVE7QUFDUixXQUFPO0FBQUEsRUFDWDtBQUVBLFNBQU8sQ0FBQyxXQUFXO0FBQ2YsZUFBVyxRQUFRLGFBQWE7QUFDNUIsVUFBSSxRQUFRLFFBQVE7QUFDaEIsZUFBTyxJQUFJLElBQUksWUFBWSxJQUFJLEVBQUUsT0FBTyxJQUFJLENBQUM7QUFBQSxNQUNqRDtBQUFBLElBQ0o7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUNKOzs7QUN6R0E7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBd0RBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksT0FBTztBQUVqRCxJQUFNLFNBQVM7QUFDZixJQUFNLGFBQWE7QUFDbkIsSUFBTSxhQUFhO0FBT1osU0FBUyxTQUE0QjtBQUN4QyxTQUFPQSxNQUFLLE1BQU07QUFDdEI7QUFPTyxTQUFTLGFBQThCO0FBQzFDLFNBQU9BLE1BQUssVUFBVTtBQUMxQjtBQU9PLFNBQVMsYUFBOEI7QUFDMUMsU0FBT0EsTUFBSyxVQUFVO0FBQzFCOzs7QXRCNUVBLE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQTRDbEMsT0FBTyxPQUFPLFNBQWdCO0FBQ3ZCLE9BQU8scUJBQXFCOyIsCiAgIm5hbWVzIjogWyJfYSIsICJFcnJvciIsICJjYWxsIiwgIl9hIiwgIkVycm9yIiwgImNhbGwiLCAiX2EiLCAicmVzaXphYmxlIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJjYWxsIiwgIkhpZGVNZXRob2QiLCAiU2hvd01ldGhvZCIsICJpc0RvY3VtZW50RG90QWxsIiwgIl9hIiwgInJlYXNvbiIsICJ2YWx1ZSIsICJjYWxsIiwgImdldEFuZERlbGV0ZVJlc3BvbnNlIiwgImdlbmVyYXRlSUQiLCAiY2FsbCIsICJBcnJheSIsICJNYXAiLCAiQXJyYXkiLCAiTWFwIiwgImtleSIsICJjYWxsIl0KfQo=
