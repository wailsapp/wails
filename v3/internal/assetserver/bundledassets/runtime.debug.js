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
    ApplicationLaunchedWithUrl: "common:ApplicationLaunchedWithUrl",
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
  let eventName;
  let eventData;
  if (typeof name === "object" && name !== null && "name" in name && "data" in name) {
    eventName = name["name"];
    eventData = name["data"];
  } else {
    eventName = name;
    eventData = data;
  }
  return call3(EmitMethod, { name: eventName, data: eventData });
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHBsaWNhdGlvbi50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYW5jZWxsYWJsZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG5cclxuaW1wb3J0IFwiLi9jb250ZXh0bWVudS5qc1wiO1xyXG5pbXBvcnQgXCIuL2RyYWcuanNcIjtcclxuXHJcbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXHJcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gXCIuL2FwcGxpY2F0aW9uLmpzXCI7XHJcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xyXG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzLmpzXCI7XHJcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmQuanNcIjtcclxuaW1wb3J0ICogYXMgQ3JlYXRlIGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xyXG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcclxuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuL2V2ZW50cy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBGbGFncyBmcm9tIFwiLi9mbGFncy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnMuanNcIjtcclxuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gXCIuL3N5c3RlbS5qc1wiO1xyXG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sLmpzXCI7XHJcblxyXG5leHBvcnQge1xyXG4gICAgQXBwbGljYXRpb24sXHJcbiAgICBCcm93c2VyLFxyXG4gICAgQ2FsbCxcclxuICAgIENsaXBib2FyZCxcclxuICAgIERpYWxvZ3MsXHJcbiAgICBFdmVudHMsXHJcbiAgICBGbGFncyxcclxuICAgIFNjcmVlbnMsXHJcbiAgICBTeXN0ZW0sXHJcbiAgICBXaW5kb3csXHJcbiAgICBXTUxcclxufTtcclxuXHJcbi8qKlxyXG4gKiBBbiBpbnRlcm5hbCB1dGlsaXR5IGNvbnN1bWVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cclxuICpcclxuICogQGlnbm9yZVxyXG4gKiBAaW50ZXJuYWxcclxuICovXHJcbmV4cG9ydCB7IENyZWF0ZSB9O1xyXG5cclxuZXhwb3J0ICogZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcclxuXHJcbi8vIE5vdGlmeSBiYWNrZW5kXHJcbndpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcclxuU3lzdGVtLmludm9rZShcIndhaWxzOnJ1bnRpbWU6cmVhZHlcIik7XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBPcGVuVVJMIH0gZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xyXG5pbXBvcnQgeyBRdWVzdGlvbiB9IGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcclxuaW1wb3J0IHsgRW1pdCwgV2FpbHNFdmVudCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xyXG5pbXBvcnQgeyBjYW5BYm9ydExpc3RlbmVycywgd2hlblJlYWR5IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcclxuaW1wb3J0IFdpbmRvdyBmcm9tIFwiLi93aW5kb3cuanNcIjtcclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBvcHRpb25hbCBkYXRhLlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gc2VuZC5cclxuICogQHBhcmFtIFtkYXRhPW51bGxdIC0gLSBPcHRpb25hbCBkYXRhIHRvIHNlbmQgYWxvbmcgd2l0aCB0aGUgZXZlbnQuXHJcbiAqL1xyXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpOiB2b2lkIHtcclxuICAgIEVtaXQobmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAqXHJcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lOiBzdHJpbmcsIG1ldGhvZE5hbWU6IHN0cmluZykge1xyXG4gICAgY29uc3QgdGFyZ2V0V2luZG93ID0gV2luZG93LkdldCh3aW5kb3dOYW1lKTtcclxuICAgIGNvbnN0IG1ldGhvZCA9ICh0YXJnZXRXaW5kb3cgYXMgYW55KVttZXRob2ROYW1lXTtcclxuXHJcbiAgICBpZiAodHlwZW9mIG1ldGhvZCAhPT0gXCJmdW5jdGlvblwiKSB7XHJcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWV0aG9kLmNhbGwodGFyZ2V0V2luZG93KTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXQuXHJcbiAqL1xyXG5mdW5jdGlvbiBvbldNTFRyaWdnZXJlZChldjogRXZlbnQpOiB2b2lkIHtcclxuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0IGFzIEVsZW1lbnQ7XHJcblxyXG4gICAgZnVuY3Rpb24gcnVuRWZmZWN0KGNob2ljZSA9IFwiWWVzXCIpIHtcclxuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcclxuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XHJcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XHJcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW51cmwnKTtcclxuXHJcbiAgICAgICAgaWYgKGV2ZW50VHlwZSAhPT0gbnVsbClcclxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XHJcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcclxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcclxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XHJcblxyXG4gICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICBRdWVzdGlvbih7XHJcbiAgICAgICAgICAgIFRpdGxlOiBcIkNvbmZpcm1cIixcclxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcclxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxyXG4gICAgICAgICAgICBCdXR0b25zOiBbXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XHJcbiAgICAgICAgICAgIF1cclxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJ1bkVmZmVjdCgpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBjb250cm9sbGVyU3ltID0gU3ltYm9sKFwiY29udHJvbGxlclwiKTtcclxuY29uc3QgdHJpZ2dlck1hcFN5bSA9IFN5bWJvbChcInRyaWdnZXJNYXBcIik7XHJcbmNvbnN0IGVsZW1lbnRDb3VudFN5bSA9IFN5bWJvbChcImVsZW1lbnRDb3VudFwiKTtcclxuXHJcbi8qKlxyXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXHJcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXHJcbiAqL1xyXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XHJcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cclxuICAgIFtjb250cm9sbGVyU3ltXTogQWJvcnRDb250cm9sbGVyO1xyXG5cclxuICAgIGNvbnN0cnVjdG9yKCkge1xyXG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIGFuIG9wdGlvbnMgb2JqZWN0IGZvciBhZGRFdmVudExpc3RlbmVyIHRoYXQgdGllcyB0aGUgbGlzdGVuZXJcclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxyXG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xyXG4gICAgICovXHJcbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xyXG4gICAgICAgIHJldHVybiB7IHNpZ25hbDogdGhpc1tjb250cm9sbGVyU3ltXS5zaWduYWwgfTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzIGFuZCByZXNldHMgdGhlIHJlZ2lzdHJ5LlxyXG4gICAgICovXHJcbiAgICByZXNldCgpOiB2b2lkIHtcclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dLmFib3J0KCk7XHJcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFdlYWtNYXBSZWdpc3RyeSBtYXBzIGFjdGl2ZSB0cmlnZ2VyIGV2ZW50cyB0byBlYWNoIERPTSBlbGVtZW50IHRocm91Z2ggYSBXZWFrTWFwLlxyXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcclxuICogY29sbGVjdGlvbiBvZiB0aGUgaW52b2x2ZWQgZWxlbWVudHMuXHJcbiAqL1xyXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xyXG4gICAgLyoqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy4gKi9cclxuICAgIFt0cmlnZ2VyTWFwU3ltXTogV2Vha01hcDxFbGVtZW50LCBzdHJpbmdbXT47XHJcbiAgICAvKiogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLiAqL1xyXG4gICAgW2VsZW1lbnRDb3VudFN5bV06IG51bWJlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcigpIHtcclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxyXG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XHJcbiAgICAgKi9cclxuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XHJcbiAgICAgICAgaWYgKCF0aGlzW3RyaWdnZXJNYXBTeW1dLmhhcyhlbGVtZW50KSkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0rKzsgfVxyXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcclxuICAgICAgICByZXR1cm4ge307XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cclxuICAgICAqL1xyXG4gICAgcmVzZXQoKTogdm9pZCB7XHJcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xyXG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuXHJcbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwU3ltXS5nZXQoZWxlbWVudCk7XHJcbiAgICAgICAgICAgIGlmICh0cmlnZ2VycyAhPSBudWxsKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXS0tOyB9XHJcblxyXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXHJcbiAgICAgICAgICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XHJcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcclxuICAgIH1cclxufVxyXG5cclxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xyXG5cclxuLyoqXHJcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cclxuICovXHJcbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50OiBFbGVtZW50KTogdm9pZCB7XHJcbiAgICBjb25zdCB0cmlnZ2VyUmVnRXhwID0gL1xcUysvZztcclxuICAgIGNvbnN0IHRyaWdnZXJBdHRyID0gKGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiKTtcclxuICAgIGNvbnN0IHRyaWdnZXJzOiBzdHJpbmdbXSA9IFtdO1xyXG5cclxuICAgIGxldCBtYXRjaDtcclxuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcclxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcclxuXHJcbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XHJcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVuYWJsZSgpOiB2b2lkIHtcclxuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpOiB2b2lkIHtcclxuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xyXG4gICAgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XSwgW3dtbC13aW5kb3ddLCBbd21sLW9wZW51cmxdLCBbZGF0YS13bWwtZXZlbnRdLCBbZGF0YS13bWwtd2luZG93XSwgW2RhdGEtd21sLW9wZW51cmxdJykuZm9yRWFjaChhZGRXTUxMaXN0ZW5lcnMpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xyXG5cclxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xyXG5cclxuLyoqXHJcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMLlxyXG4gKlxyXG4gKiBAcGFyYW0gdXJsIC0gVGhlIFVSTCB0byBvcGVuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmw6IHN0cmluZyB8IFVSTCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmw6IHVybC50b1N0cmluZygpfSk7XHJcbn1cclxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxyXG5cclxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXHJcbi8vXHJcbi8vIENvcHlyaWdodCAyMDE3IEFuZHJleSBTaXRuaWsgPGFuZHJleUBzaXRuaWsucnU+XHJcbi8vXHJcbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2ZcclxuLy8gdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbCBpblxyXG4vLyB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvXHJcbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXHJcbi8vIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbyxcclxuLy8gICAgIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxyXG4vL1xyXG4vLyAgICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXHJcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXHJcbi8vXHJcbi8vICAgICBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXHJcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTXHJcbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxyXG4vLyBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVJcclxuLy8gSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU5cclxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cclxuXHJcbi8vIFRoaXMgYWxwaGFiZXQgdXNlcyBgQS1aYS16MC05Xy1gIHN5bWJvbHMuXHJcbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cclxuLy8gUmVmZXJlbmNlcyB0byB0aGUgc2FtZSBmaWxlICh3b3JrcyBib3RoIGZvciBnemlwIGFuZCBicm90bGkpOlxyXG4vLyBgJ3VzZWAsIGBhbmRvbWAsIGFuZCBgcmljdCdgXHJcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XHJcbi8vIGAtMjZUYCwgYDE5ODNgLCBgNDBweGAsIGA3NXB4YCwgYGJ1c2hgLCBgamFja2AsIGBtaW5kYCwgYHZlcnlgLCBhbmQgYHdvbGZgXHJcbmNvbnN0IHVybEFscGhhYmV0ID1cclxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5hbm9pZChzaXplOiBudW1iZXIgPSAyMSk6IHN0cmluZyB7XHJcbiAgICBsZXQgaWQgPSAnJ1xyXG4gICAgLy8gQSBjb21wYWN0IGFsdGVybmF0aXZlIGZvciBgZm9yICh2YXIgaSA9IDA7IGkgPCBzdGVwOyBpKyspYC5cclxuICAgIGxldCBpID0gc2l6ZSB8IDBcclxuICAgIHdoaWxlIChpLS0pIHtcclxuICAgICAgICAvLyBgfCAwYCBpcyBtb3JlIGNvbXBhY3QgYW5kIGZhc3RlciB0aGFuIGBNYXRoLmZsb29yKClgLlxyXG4gICAgICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cclxuICAgIH1cclxuICAgIHJldHVybiBpZFxyXG59XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XHJcblxyXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcclxuXHJcbi8vIE9iamVjdCBOYW1lc1xyXG5leHBvcnQgY29uc3Qgb2JqZWN0TmFtZXMgPSBPYmplY3QuZnJlZXplKHtcclxuICAgIENhbGw6IDAsXHJcbiAgICBDbGlwYm9hcmQ6IDEsXHJcbiAgICBBcHBsaWNhdGlvbjogMixcclxuICAgIEV2ZW50czogMyxcclxuICAgIENvbnRleHRNZW51OiA0LFxyXG4gICAgRGlhbG9nOiA1LFxyXG4gICAgV2luZG93OiA2LFxyXG4gICAgU2NyZWVuczogNyxcclxuICAgIFN5c3RlbTogOCxcclxuICAgIEJyb3dzZXI6IDksXHJcbiAgICBDYW5jZWxDYWxsOiAxMCxcclxufSk7XHJcbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcclxuXHJcbi8qKlxyXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxyXG4gKlxyXG4gKiBAcGFyYW0gb2JqZWN0IC0gVGhlIG9iamVjdCB0byBpbnZva2UgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcmV0dXJuIFRoZSBuZXcgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcclxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kOiBudW1iZXIsIGFyZ3M6IGFueSA9IG51bGwpIHtcclxuICAgICAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0LCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xyXG4gICAgfTtcclxufVxyXG5cclxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQudG9TdHJpbmcoKSk7XHJcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QudG9TdHJpbmcoKSk7XHJcbiAgICBpZiAoYXJncykgeyB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpOyB9XHJcblxyXG4gICAgbGV0IGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4gPSB7XHJcbiAgICAgICAgW1wieC13YWlscy1jbGllbnQtaWRcIl06IGNsaWVudElkXHJcbiAgICB9XHJcbiAgICBpZiAod2luZG93TmFtZSkge1xyXG4gICAgICAgIGhlYWRlcnNbXCJ4LXdhaWxzLXdpbmRvdy1uYW1lXCJdID0gd2luZG93TmFtZTtcclxuICAgIH1cclxuXHJcbiAgICBsZXQgcmVzcG9uc2UgPSBhd2FpdCBmZXRjaCh1cmwsIHsgaGVhZGVycyB9KTtcclxuICAgIGlmICghcmVzcG9uc2Uub2spIHtcclxuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcG9uc2UudGV4dCgpKTtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpPy5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSA/PyAtMSkgIT09IC0xKSB7XHJcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlLnRleHQoKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICcuL25hbm9pZC5qcyc7XHJcblxyXG4vLyBzZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5kaWFsb2dFcnJvckNhbGxiYWNrID0gZGlhbG9nRXJyb3JDYWxsYmFjaztcclxud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xyXG5cclxudHlwZSBQcm9taXNlUmVzb2x2ZXJzID0gT21pdDxQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIj47XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5EaWFsb2cpO1xyXG5jb25zdCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwPHN0cmluZywgUHJvbWlzZVJlc29sdmVycz4oKTtcclxuXHJcbi8vIERlZmluZSBjb25zdGFudHMgZnJvbSB0aGUgYG1ldGhvZHNgIG9iamVjdCBpbiBUaXRsZSBDYXNlXHJcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xyXG5jb25zdCBEaWFsb2dXYXJuaW5nID0gMTtcclxuY29uc3QgRGlhbG9nRXJyb3IgPSAyO1xyXG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XHJcbmNvbnN0IERpYWxvZ09wZW5GaWxlID0gNDtcclxuY29uc3QgRGlhbG9nU2F2ZUZpbGUgPSA1O1xyXG5cclxuZXhwb3J0IGludGVyZmFjZSBPcGVuRmlsZURpYWxvZ09wdGlvbnMge1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xyXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xyXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cclxuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cclxuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xyXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgbXVsdGlwbGUgc2VsZWN0aW9uIGlzIGFsbG93ZWQuICovXHJcbiAgICBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cclxuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXHJcbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xyXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXHJcbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XHJcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xyXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcclxuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgVGl0bGU/OiBzdHJpbmc7XHJcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBNZXNzYWdlPzogc3RyaW5nO1xyXG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xyXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcclxuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXHJcbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgU2F2ZUZpbGVEaWFsb2dPcHRpb25zIHtcclxuICAgIC8qKiBEZWZhdWx0IGZpbGVuYW1lIHRvIHVzZSBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgRmlsZW5hbWU/OiBzdHJpbmc7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xyXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xyXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXHJcbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXHJcbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xyXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cclxuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xyXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cclxuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XHJcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cclxuICAgIFRpdGxlPzogc3RyaW5nO1xyXG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgTWVzc2FnZT86IHN0cmluZztcclxuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cclxuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XHJcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIERpcmVjdG9yeT86IHN0cmluZztcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xyXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIE1lc3NhZ2VEaWFsb2dPcHRpb25zIHtcclxuICAgIC8qKiBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuICovXHJcbiAgICBUaXRsZT86IHN0cmluZztcclxuICAgIC8qKiBUaGUgbWFpbiBtZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XHJcbiAgICAvKiogQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgQnV0dG9ucz86IEJ1dHRvbltdO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93IChpZiBhcHBsaWNhYmxlKS4gKi9cclxuICAgIERldGFjaGVkPzogYm9vbGVhbjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBCdXR0b24ge1xyXG4gICAgLyoqIFRleHQgdGhhdCBhcHBlYXJzIHdpdGhpbiB0aGUgYnV0dG9uLiAqL1xyXG4gICAgTGFiZWw/OiBzdHJpbmc7XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC4gKi9cclxuICAgIElzQ2FuY2VsPzogYm9vbGVhbjtcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGJlIHRoZSBkZWZhdWx0IGFjdGlvbiB3aGVuIHRoZSB1c2VyIHByZXNzZXMgZW50ZXIuICovXHJcbiAgICBJc0RlZmF1bHQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIEZpbGVGaWx0ZXIge1xyXG4gICAgLyoqIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLiAqL1xyXG4gICAgRGlzcGxheU5hbWU/OiBzdHJpbmc7XHJcbiAgICAvKiogUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy4gKi9cclxuICAgIFBhdHRlcm4/OiBzdHJpbmc7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIHRoZSByZXN1bHQgb2YgYSBkaWFsb2cgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSByZXF1ZXN0IHRvIGhhbmRsZSB0aGUgcmVzdWx0IGZvci5cclxuICogQHBhcmFtIGRhdGEgLSBUaGUgcmVzdWx0IGRhdGEgb2YgdGhlIHJlcXVlc3QuXHJcbiAqIEBwYXJhbSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZ1Jlc3VsdENhbGxiYWNrKGlkOiBzdHJpbmcsIGRhdGE6IHN0cmluZywgaXNKU09OOiBib29sZWFuKTogdm9pZCB7XHJcbiAgICBsZXQgcmVzb2x2ZXJzID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xyXG4gICAgaWYgKCFyZXNvbHZlcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGlzSlNPTikge1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIHJlc29sdmVycy5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gY2F0Y2ggKGVycjogYW55KSB7XHJcbiAgICAgICAgICAgIHJlc29sdmVycy5yZWplY3QobmV3IFR5cGVFcnJvcihcImNvdWxkIG5vdCBwYXJzZSByZXN1bHQ6IFwiICsgZXJyLm1lc3NhZ2UsIHsgY2F1c2U6IGVyciB9KSk7XHJcbiAgICAgICAgfVxyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShkYXRhKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIGVycm9yIGZyb20gYSBkaWFsb2cgcmVxdWVzdC5cclxuICpcclxuICogQHBhcmFtIGlkIC0gVGhlIGlkIG9mIHRoZSBwcm9taXNlIGhhbmRsZXIuXHJcbiAqIEBwYXJhbSBtZXNzYWdlIC0gQW4gZXJyb3IgbWVzc2FnZS5cclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQ6IHN0cmluZywgbWVzc2FnZTogc3RyaW5nKTogdm9pZCB7XHJcbiAgICBnZXRBbmREZWxldGVSZXNwb25zZShpZCk/LnJlamVjdChuZXcgd2luZG93LkVycm9yKG1lc3NhZ2UpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyBhbmQgcmVtb3ZlcyB0aGUgcmVzcG9uc2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCBmcm9tIHRoZSBkaWFsb2dSZXNwb25zZXMgbWFwLlxyXG4gKlxyXG4gKiBAcGFyYW0gaWQgLSBUaGUgSUQgb2YgdGhlIHJlc3BvbnNlIHRvIGJlIHJldHJpZXZlZCBhbmQgcmVtb3ZlZC5cclxuICogQHJldHVybnMgVGhlIHJlc3BvbnNlIG9iamVjdCBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElELCBpZiBhbnkuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZXRBbmREZWxldGVSZXNwb25zZShpZDogc3RyaW5nKTogUHJvbWlzZVJlc29sdmVycyB8IHVuZGVmaW5lZCB7XHJcbiAgICBjb25zdCByZXNwb25zZSA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICByZXR1cm4gcmVzcG9uc2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBkaWFsb2dSZXNwb25zZXMgc2V0LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICpcclxuICogQHBhcmFtIHR5cGUgLSBEaWFsb2cgdHlwZS5cclxuICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIGZvciB0aGUgZGlhbG9nLlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHJlc3VsdCBvZiBkaWFsb2cuXHJcbiAqL1xyXG5mdW5jdGlvbiBkaWFsb2codHlwZTogbnVtYmVyLCBvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyB8IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB8IFNhdmVGaWxlRGlhbG9nT3B0aW9ucyA9IHt9KTogUHJvbWlzZTxhbnk+IHtcclxuICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmUsIHJlamVjdCB9KTtcclxuICAgICAgICBjYWxsKHR5cGUsIE9iamVjdC5hc3NpZ24oeyBcImRpYWxvZy1pZFwiOiBpZCB9LCBvcHRpb25zKSkuY2F0Y2goKGVycjogYW55KSA9PiB7XHJcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgICAgICByZWplY3QoZXJyKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYW4gaW5mbyBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIHdhcm5pbmcgZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXYXJuaW5nKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dXYXJuaW5nLCBvcHRpb25zKTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGFuIGVycm9yIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRXJyb3Iob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgcXVlc3Rpb24gZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSBmaWxlIHNlbGVjdGlvbiBkaWFsb2cgdG8gcGljayBvbmUgb3IgbW9yZSBmaWxlcyB0byBvcGVuLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxyXG4gKiBAcmV0dXJucyBTZWxlY3RlZCBmaWxlIG9yIGxpc3Qgb2YgZmlsZXMsIG9yIGEgYmxhbmsgc3RyaW5nL2VtcHR5IGxpc3QgaWYgbm8gZmlsZSBoYXMgYmVlbiBzZWxlY3RlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uOiB0cnVlIH0pOiBQcm9taXNlPHN0cmluZ1tdPjtcclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyAmIHsgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24/OiBmYWxzZSB8IHVuZGVmaW5lZCB9KTogUHJvbWlzZTxzdHJpbmc+O1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT47XHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZyB8IHN0cmluZ1tdPiB7IHJldHVybiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpID8/IFtdOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSBmaWxlIHNlbGVjdGlvbiBkaWFsb2cgdG8gcGljayBhIGZpbGUgdG8gc2F2ZS5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSwgb3IgYSBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBoYXMgYmVlbiBzZWxlY3RlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zOiBTYXZlRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTsgfVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB7IGV2ZW50TGlzdGVuZXJzLCBMaXN0ZW5lciwgbGlzdGVuZXJPZmYgfSBmcm9tIFwiLi9saXN0ZW5lci5qc1wiO1xyXG5cclxuLy8gU2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbndpbmRvdy5fd2FpbHMuZGlzcGF0Y2hXYWlsc0V2ZW50ID0gZGlzcGF0Y2hXYWlsc0V2ZW50O1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRXZlbnRzKTtcclxuY29uc3QgRW1pdE1ldGhvZCA9IDA7XHJcblxyXG5leHBvcnQgeyBUeXBlcyB9IGZyb20gXCIuL2V2ZW50X3R5cGVzLmpzXCI7XHJcblxyXG4vKipcclxuICogVGhlIHR5cGUgb2YgaGFuZGxlcnMgZm9yIGEgZ2l2ZW4gZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgdHlwZSBDYWxsYmFjayA9IChldjogV2FpbHNFdmVudCkgPT4gdm9pZDtcclxuXHJcbi8qKlxyXG4gKiBSZXByZXNlbnRzIGEgc3lzdGVtIGV2ZW50IG9yIGEgY3VzdG9tIGV2ZW50IGVtaXR0ZWQgdGhyb3VnaCB3YWlscy1wcm92aWRlZCBmYWNpbGl0aWVzLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQge1xyXG4gICAgLyoqXHJcbiAgICAgKiBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXHJcbiAgICAgKi9cclxuICAgIG5hbWU6IHN0cmluZztcclxuXHJcbiAgICAvKipcclxuICAgICAqIE9wdGlvbmFsIGRhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBlbWl0dGVkIGV2ZW50LlxyXG4gICAgICovXHJcbiAgICBkYXRhOiBhbnk7XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBOYW1lIG9mIHRoZSBvcmlnaW5hdGluZyB3aW5kb3cuIE9taXR0ZWQgZm9yIGFwcGxpY2F0aW9uIGV2ZW50cy5cclxuICAgICAqIFdpbGwgYmUgb3ZlcnJpZGRlbiBpZiBzZXQgbWFudWFsbHkuXHJcbiAgICAgKi9cclxuICAgIHNlbmRlcj86IHN0cmluZztcclxuXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpIHtcclxuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xyXG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudDogYW55KSB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xyXG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgbGV0IHdhaWxzRXZlbnQgPSBuZXcgV2FpbHNFdmVudChldmVudC5uYW1lLCBldmVudC5kYXRhKTtcclxuICAgIGlmICgnc2VuZGVyJyBpbiBldmVudCkge1xyXG4gICAgICAgIHdhaWxzRXZlbnQuc2VuZGVyID0gZXZlbnQuc2VuZGVyO1xyXG4gICAgfVxyXG5cclxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4gIWxpc3RlbmVyLmRpc3BhdGNoKHdhaWxzRXZlbnQpKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXHJcbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxyXG4gKiBAcGFyYW0gbWF4Q2FsbGJhY2tzIC0gVGhlIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRoZSBjYWxsYmFjayBjYW4gYmUgY2FsbGVkIGZvciB0aGUgZXZlbnQuIE9uY2UgdGhlIG1heGltdW0gbnVtYmVyIGlzIHJlYWNoZWQsIHRoZSBjYWxsYmFjayB3aWxsIG5vIGxvbmdlciBiZSBjYWxsZWQuXHJcbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IENhbGxiYWNrLCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XHJcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxyXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IENhbGxiYWNrKTogKCkgPT4gdm9pZCB7XHJcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXHJcbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IENhbGxiYWNrKTogKCkgPT4gdm9pZCB7XHJcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlbW92ZXMgZXZlbnQgbGlzdGVuZXJzIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50IG5hbWVzLlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lcyAtIFRoZSBuYW1lIG9mIHRoZSBldmVudHMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmKC4uLmV2ZW50TmFtZXM6IFtzdHJpbmcsIC4uLnN0cmluZ1tdXSk6IHZvaWQge1xyXG4gICAgZXZlbnROYW1lcy5mb3JFYWNoKGV2ZW50TmFtZSA9PiBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCk6IHZvaWQge1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEVtaXRzIGFuIGV2ZW50IHVzaW5nIHRoZSBuYW1lIGFuZCBkYXRhLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCB3aWxsIGJlIGZ1bGZpbGxlZCBvbmNlIHRoZSBldmVudCBoYXMgYmVlbiBlbWl0dGVkLlxyXG4gKiBAcGFyYW0gbmFtZSAtIHRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxyXG4gKiBAcGFyYW0gZGF0YSAtIHRoZSBkYXRhIHRvIGJlIHNlbnQgd2l0aCB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW1pdChuYW1lOiBzdHJpbmcsIGRhdGE/OiBhbnkpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIGxldCBldmVudE5hbWU6IHN0cmluZztcclxuICAgIGxldCBldmVudERhdGE6IGFueTtcclxuXHJcbiAgICBpZiAodHlwZW9mIG5hbWUgPT09ICdvYmplY3QnICYmIG5hbWUgIT09IG51bGwgJiYgJ25hbWUnIGluIG5hbWUgJiYgJ2RhdGEnIGluIG5hbWUpIHtcclxuICAgICAgICAvLyBJZiBuYW1lIGlzIGFuIG9iamVjdCB3aXRoIGEgbmFtZSBwcm9wZXJ0eSwgdXNlIGl0IGRpcmVjdGx5XHJcbiAgICAgICAgZXZlbnROYW1lID0gbmFtZVsnbmFtZSddO1xyXG4gICAgICAgIGV2ZW50RGF0YSA9IG5hbWVbJ2RhdGEnXTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgLy8gT3RoZXJ3aXNlIHVzZSB0aGUgc3RhbmRhcmQgcGFyYW1ldGVyc1xyXG4gICAgICAgIGV2ZW50TmFtZSA9IG5hbWUgYXMgc3RyaW5nO1xyXG4gICAgICAgIGV2ZW50RGF0YSA9IGRhdGE7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgeyBuYW1lOiBldmVudE5hbWUsIGRhdGE6IGV2ZW50RGF0YSB9KTtcclxufVxyXG5cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8vIFRoZSBmb2xsb3dpbmcgdXRpbGl0aWVzIGhhdmUgYmVlbiBmYWN0b3JlZCBvdXQgb2YgLi9ldmVudHMudHNcclxuLy8gZm9yIHRlc3RpbmcgcHVycG9zZXMuXHJcblxyXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwPHN0cmluZywgTGlzdGVuZXJbXT4oKTtcclxuXHJcbmV4cG9ydCBjbGFzcyBMaXN0ZW5lciB7XHJcbiAgICBldmVudE5hbWU6IHN0cmluZztcclxuICAgIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkO1xyXG4gICAgbWF4Q2FsbGJhY2tzOiBudW1iZXI7XHJcblxyXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkLCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xyXG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xyXG4gICAgICAgIHRoaXMuY2FsbGJhY2sgPSBjYWxsYmFjaztcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgIH1cclxuXHJcbiAgICBkaXNwYXRjaChkYXRhOiBhbnkpOiBib29sZWFuIHtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICB0aGlzLmNhbGxiYWNrKGRhdGEpO1xyXG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGVycik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcclxuICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcjogTGlzdGVuZXIpOiB2b2lkIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQobGlzdGVuZXIuZXZlbnROYW1lKTtcclxuICAgIGlmICghbGlzdGVuZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XHJcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShsaXN0ZW5lci5ldmVudE5hbWUpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBldmVudExpc3RlbmVycy5zZXQobGlzdGVuZXIuZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vLyBDeW5oeXJjaHd5ZCB5IGZmZWlsIGhvbiB5biBhd3RvbWF0aWcuIFBFSURJV0NIIFx1MDBDMiBNT0RJV0xcclxuLy8gVGhpcyBmaWxlIGlzIGF1dG9tYXRpY2FsbHkgZ2VuZXJhdGVkLiBETyBOT1QgRURJVFxyXG5cclxuZXhwb3J0IGNvbnN0IFR5cGVzID0gT2JqZWN0LmZyZWV6ZSh7XHJcblx0V2luZG93czogT2JqZWN0LmZyZWV6ZSh7XHJcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcclxuXHRcdEFQTVBvd2VyU3RhdHVzQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTdGF0dXNDaGFuZ2VcIixcclxuXHRcdEFQTVJlc3VtZUF1dG9tYXRpYzogXCJ3aW5kb3dzOkFQTVJlc3VtZUF1dG9tYXRpY1wiLFxyXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcclxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXHJcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwid2luZG93czpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcclxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxyXG5cdFx0V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWQ6IFwid2luZG93czpXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZFwiLFxyXG5cdFx0V2luZG93QWN0aXZlOiBcIndpbmRvd3M6V2luZG93QWN0aXZlXCIsXHJcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcclxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcclxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwid2luZG93czpXaW5kb3dDbG9zaW5nXCIsXHJcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxyXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIndpbmRvd3M6V2luZG93RGlkUmVzaXplXCIsXHJcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93RFBJQ2hhbmdlZFwiLFxyXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxyXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXHJcblx0XHRXaW5kb3dEcmFnTGVhdmU6IFwid2luZG93czpXaW5kb3dEcmFnTGVhdmVcIixcclxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcclxuXHRcdFdpbmRvd0VuZE1vdmU6IFwid2luZG93czpXaW5kb3dFbmRNb3ZlXCIsXHJcblx0XHRXaW5kb3dFbmRSZXNpemU6IFwid2luZG93czpXaW5kb3dFbmRSZXNpemVcIixcclxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXHJcblx0XHRXaW5kb3dIaWRlOiBcIndpbmRvd3M6V2luZG93SGlkZVwiLFxyXG5cdFx0V2luZG93SW5hY3RpdmU6IFwid2luZG93czpXaW5kb3dJbmFjdGl2ZVwiLFxyXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcclxuXHRcdFdpbmRvd0tleVVwOiBcIndpbmRvd3M6V2luZG93S2V5VXBcIixcclxuXHRcdFdpbmRvd0tpbGxGb2N1czogXCJ3aW5kb3dzOldpbmRvd0tpbGxGb2N1c1wiLFxyXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXHJcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZURvd246IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZURvd25cIixcclxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlTGVhdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZUxlYXZlXCIsXHJcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcclxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlVXA6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZVVwXCIsXHJcblx0XHRXaW5kb3dQYWludDogXCJ3aW5kb3dzOldpbmRvd1BhaW50XCIsXHJcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxyXG5cdFx0V2luZG93U2V0Rm9jdXM6IFwid2luZG93czpXaW5kb3dTZXRGb2N1c1wiLFxyXG5cdFx0V2luZG93U2hvdzogXCJ3aW5kb3dzOldpbmRvd1Nob3dcIixcclxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxyXG5cdFx0V2luZG93U3RhcnRSZXNpemU6IFwid2luZG93czpXaW5kb3dTdGFydFJlc2l6ZVwiLFxyXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcIndpbmRvd3M6V2luZG93VW5GdWxsc2NyZWVuXCIsXHJcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxyXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxyXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWluaW1pc2VcIixcclxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcclxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1heGltaXNlXCIsXHJcblx0fSksXHJcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcclxuXHRcdEFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJPcmllbnRhdGlvblwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkSGlkZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkVXBkYXRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFVwZGF0ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmdcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZVwiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXHJcblx0XHRNZW51RGlkQWRkSXRlbTogXCJtYWM6TWVudURpZEFkZEl0ZW1cIixcclxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxyXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcclxuXHRcdE1lbnVEaWREaXNwbGF5SXRlbTogXCJtYWM6TWVudURpZERpc3BsYXlJdGVtXCIsXHJcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxyXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXHJcblx0XHRNZW51RGlkT3BlbjogXCJtYWM6TWVudURpZE9wZW5cIixcclxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXHJcblx0XHRNZW51RGlkUmVtb3ZlSXRlbTogXCJtYWM6TWVudURpZFJlbW92ZUl0ZW1cIixcclxuXHRcdE1lbnVEaWRTZW5kQWN0aW9uOiBcIm1hYzpNZW51RGlkU2VuZEFjdGlvblwiLFxyXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXHJcblx0XHRNZW51RGlkVXBkYXRlOiBcIm1hYzpNZW51RGlkVXBkYXRlXCIsXHJcblx0XHRNZW51V2lsbEFkZEl0ZW06IFwibWFjOk1lbnVXaWxsQWRkSXRlbVwiLFxyXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcclxuXHRcdE1lbnVXaWxsRGlzcGxheUl0ZW06IFwibWFjOk1lbnVXaWxsRGlzcGxheUl0ZW1cIixcclxuXHRcdE1lbnVXaWxsRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsRW5kVHJhY2tpbmdcIixcclxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXHJcblx0XHRNZW51V2lsbE9wZW46IFwibWFjOk1lbnVXaWxsT3BlblwiLFxyXG5cdFx0TWVudVdpbGxQb3BVcDogXCJtYWM6TWVudVdpbGxQb3BVcFwiLFxyXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcclxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvbjogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uXCIsXHJcblx0XHRNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbVwiLFxyXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXHJcblx0XHRXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb25cIixcclxuXHRcdFdlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcclxuXHRcdFdpbmRvd0RpZEJlY29tZUtleTogXCJtYWM6V2luZG93RGlkQmVjb21lS2V5XCIsXHJcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXHJcblx0XHRXaW5kb3dEaWRCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dEaWRCZWdpblNoZWV0XCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VBbHBoYTogXCJtYWM6V2luZG93RGlkQ2hhbmdlQWxwaGFcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUHJvZmlsZVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllc1wiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlVGl0bGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVRpdGxlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXHJcblx0XHRXaW5kb3dEaWREZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWREZW1pbmlhdHVyaXplXCIsXHJcblx0XHRXaW5kb3dEaWRFbmRTaGVldDogXCJtYWM6V2luZG93RGlkRW5kU2hlZXRcIixcclxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXHJcblx0XHRXaW5kb3dEaWRFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRXhpdEZ1bGxTY3JlZW5cIixcclxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXHJcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxyXG5cdFx0V2luZG93RGlkRm9jdXM6IFwibWFjOldpbmRvd0RpZEZvY3VzXCIsXHJcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcclxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcclxuXHRcdFdpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9mZlNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcclxuXHRcdFdpbmRvd0RpZFJlc2lnbk1haW46IFwibWFjOldpbmRvd0RpZFJlc2lnbk1haW5cIixcclxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxyXG5cdFx0V2luZG93RGlkVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUFscGhhXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxyXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVNoYWRvd1wiLFxyXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXHJcblx0XHRXaW5kb3dEaWRab29tOiBcIm1hYzpXaW5kb3dEaWRab29tXCIsXHJcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXHJcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZFwiLFxyXG5cdFx0V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWRcIixcclxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcclxuXHRcdFdpbmRvd01heGltaXNlOiBcIm1hYzpXaW5kb3dNYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJtYWM6V2luZG93VW5NYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIm1hYzpXaW5kb3dVbk1pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dTaG91bGRDbG9zZTogXCJtYWM6V2luZG93U2hvdWxkQ2xvc2VcIixcclxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcclxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcclxuXHRcdFdpbmRvd1dpbGxCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lTWFpblwiLFxyXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXHJcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXHJcblx0XHRXaW5kb3dXaWxsQ2xvc2U6IFwibWFjOldpbmRvd1dpbGxDbG9zZVwiLFxyXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXHJcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXJcIixcclxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXHJcblx0XHRXaW5kb3dXaWxsRm9jdXM6IFwibWFjOldpbmRvd1dpbGxGb2N1c1wiLFxyXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcclxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxyXG5cdFx0V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW5cIixcclxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxyXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXHJcblx0XHRXaW5kb3dXaWxsUmVzaXplOiBcIm1hYzpXaW5kb3dXaWxsUmVzaXplXCIsXHJcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZUFscGhhXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVNoYWRvd1wiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHk6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5XCIsXHJcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcclxuXHRcdFdpbmRvd1pvb21JbjogXCJtYWM6V2luZG93Wm9vbUluXCIsXHJcblx0XHRXaW5kb3dab29tT3V0OiBcIm1hYzpXaW5kb3dab29tT3V0XCIsXHJcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxyXG5cdH0pLFxyXG5cdExpbnV4OiBPYmplY3QuZnJlZXplKHtcclxuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcclxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcclxuXHRcdFdpbmRvd0RlbGV0ZUV2ZW50OiBcImxpbnV4OldpbmRvd0RlbGV0ZUV2ZW50XCIsXHJcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcclxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJsaW51eDpXaW5kb3dEaWRSZXNpemVcIixcclxuXHRcdFdpbmRvd0ZvY3VzSW46IFwibGludXg6V2luZG93Rm9jdXNJblwiLFxyXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcclxuXHRcdFdpbmRvd0xvYWRDaGFuZ2VkOiBcImxpbnV4OldpbmRvd0xvYWRDaGFuZ2VkXCIsXHJcblx0fSksXHJcblx0Q29tbW9uOiBPYmplY3QuZnJlZXplKHtcclxuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiY29tbW9uOkFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGVcIixcclxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXHJcblx0XHRBcHBsaWNhdGlvbkxhdW5jaGVkV2l0aFVybDogXCJjb21tb246QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmxcIixcclxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXHJcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXHJcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImNvbW1vbjpXaW5kb3dEaWRNb3ZlXCIsXHJcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwiY29tbW9uOldpbmRvd0RpZFJlc2l6ZVwiLFxyXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxyXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcclxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxyXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxyXG5cdFx0V2luZG93SGlkZTogXCJjb21tb246V2luZG93SGlkZVwiLFxyXG5cdFx0V2luZG93TG9zdEZvY3VzOiBcImNvbW1vbjpXaW5kb3dMb3N0Rm9jdXNcIixcclxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dUb2dnbGVGcmFtZWxlc3M6IFwiY29tbW9uOldpbmRvd1RvZ2dsZUZyYW1lbGVzc1wiLFxyXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxyXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcclxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcclxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXHJcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1heGltaXNlXCIsXHJcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXHJcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxyXG5cdFx0V2luZG93Wm9vbU91dDogXCJjb21tb246V2luZG93Wm9vbU91dFwiLFxyXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcclxuXHRcdFdpbmRvd0Ryb3Bab25lRmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dEcm9wWm9uZUZpbGVzRHJvcHBlZFwiLFxyXG5cdH0pLFxyXG59KTtcclxuIiwgIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxyXG4gKlxyXG4gKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBtZXNzYWdlIHRvIGJlIGxvZ2dlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlOiBhbnkpIHtcclxuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxyXG4gICAgY29uc29sZS5sb2coXHJcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXHJcbiAgICApO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIHdoZXRoZXIgdGhlIHdlYnZpZXcgc3VwcG9ydHMgdGhlIHtAbGluayBNb3VzZUV2ZW50I2J1dHRvbnN9IHByb3BlcnR5LlxyXG4gKiBMb29raW5nIGF0IHlvdSBtYWNPUyBIaWdoIFNpZXJyYSFcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBjYW5UcmFja0J1dHRvbnMoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKG5ldyBNb3VzZUV2ZW50KCdtb3VzZWRvd24nKSkuYnV0dG9ucyA9PT0gMDtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXHJcbiAqIChzZWUgaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmcvZW4tVVMvZG9jcy9XZWIvQVBJL0V2ZW50VGFyZ2V0L2FkZEV2ZW50TGlzdGVuZXIjc2lnbmFsKS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBjYW5BYm9ydExpc3RlbmVycygpIHtcclxuICAgIGlmICghRXZlbnRUYXJnZXQgfHwgIUFib3J0U2lnbmFsIHx8ICFBYm9ydENvbnRyb2xsZXIpXHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG5cclxuICAgIGxldCByZXN1bHQgPSB0cnVlO1xyXG5cclxuICAgIGNvbnN0IHRhcmdldCA9IG5ldyBFdmVudFRhcmdldCgpO1xyXG4gICAgY29uc3QgY29udHJvbGxlciA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xyXG4gICAgY29udHJvbGxlci5hYm9ydCgpO1xyXG4gICAgdGFyZ2V0LmRpc3BhdGNoRXZlbnQobmV3IEN1c3RvbUV2ZW50KCd0ZXN0JykpO1xyXG5cclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXNvbHZlcyB0aGUgY2xvc2VzdCBIVE1MRWxlbWVudCBhbmNlc3RvciBvZiBhbiBldmVudCdzIHRhcmdldC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBldmVudFRhcmdldChldmVudDogRXZlbnQpOiBIVE1MRWxlbWVudCB7XHJcbiAgICBpZiAoZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpIHtcclxuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0O1xyXG4gICAgfSBlbHNlIGlmICghKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSAmJiBldmVudC50YXJnZXQgaW5zdGFuY2VvZiBOb2RlKSB7XHJcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldC5wYXJlbnRFbGVtZW50ID8/IGRvY3VtZW50LmJvZHk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiBkb2N1bWVudC5ib2R5O1xyXG4gICAgfVxyXG59XHJcblxyXG4vKioqXHJcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcclxuXHJcbiBCU0QgMi1DbGF1c2UgTGljZW5zZVxyXG5cclxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxyXG4gQWxsIHJpZ2h0cyByZXNlcnZlZC5cclxuXHJcbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcclxuIG1vZGlmaWNhdGlvbiwgYXJlIHBlcm1pdHRlZCBwcm92aWRlZCB0aGF0IHRoZSBmb2xsb3dpbmcgY29uZGl0aW9ucyBhcmUgbWV0OlxyXG5cclxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xyXG4gbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIuXHJcblxyXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxyXG4gdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvblxyXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXHJcblxyXG4gVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyBcIkFTIElTXCJcclxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcclxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxyXG4gRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1JTIEJFIExJQUJMRVxyXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcclxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXHJcbiBTRVJWSUNFUzsgTE9TUyBPRiBVU0UsIERBVEEsIE9SIFBST0ZJVFM7IE9SIEJVU0lORVNTIElOVEVSUlVQVElPTikgSE9XRVZFUlxyXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcclxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXHJcbiBPRiBUSElTIFNPRlRXQVJFLCBFVkVOIElGIEFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLlxyXG5cclxuICoqKi9cclxuXHJcbmxldCBpc1JlYWR5ID0gZmFsc2U7XHJcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCAoKSA9PiB7IGlzUmVhZHkgPSB0cnVlIH0pO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjazogKCkgPT4gdm9pZCkge1xyXG4gICAgaWYgKGlzUmVhZHkgfHwgZG9jdW1lbnQucmVhZHlTdGF0ZSA9PT0gJ2NvbXBsZXRlJykge1xyXG4gICAgICAgIGNhbGxiYWNrKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCBjYWxsYmFjayk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHR5cGUgeyBTY3JlZW4gfSBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XHJcblxyXG4vLyBORVc6IERyb3B6b25lIGNvbnN0YW50c1xyXG5jb25zdCBEUk9QWk9ORV9BVFRSSUJVVEUgPSAnZGF0YS13YWlscy1kcm9wem9uZSc7XHJcbmNvbnN0IERST1BaT05FX0hPVkVSX0NMQVNTID0gJ3dhaWxzLWRyb3B6b25lLWhvdmVyJzsgLy8gVXNlciBjYW4gc3R5bGUgdGhpcyBjbGFzc1xyXG5sZXQgY3VycmVudEhvdmVyZWREcm9wem9uZTogRWxlbWVudCB8IG51bGwgPSBudWxsO1xyXG5cclxuY29uc3QgUG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgICAgID0gMDtcclxuY29uc3QgQ2VudGVyTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMTtcclxuY29uc3QgQ2xvc2VNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gMjtcclxuY29uc3QgRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgID0gMztcclxuY29uc3QgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgID0gNDtcclxuY29uc3QgRm9jdXNNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNTtcclxuY29uc3QgRm9yY2VSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgID0gNjtcclxuY29uc3QgRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgID0gNztcclxuY29uc3QgR2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gODtcclxuY29uc3QgR2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gOTtcclxuY29uc3QgSGVpZ2h0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMTA7XHJcbmNvbnN0IEhpZGVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDExO1xyXG5jb25zdCBJc0ZvY3VzZWRNZXRob2QgICAgICAgICAgICAgICAgICAgPSAxMjtcclxuY29uc3QgSXNGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gMTM7XHJcbmNvbnN0IElzTWF4aW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE0O1xyXG5jb25zdCBJc01pbmltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNTtcclxuY29uc3QgTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTY7XHJcbmNvbnN0IE1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE3O1xyXG5jb25zdCBOYW1lTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxODtcclxuY29uc3QgT3BlbkRldlRvb2xzTWV0aG9kICAgICAgICAgICAgICAgID0gMTk7XHJcbmNvbnN0IFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICAgICA9IDIwO1xyXG5jb25zdCBSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAyMTtcclxuY29uc3QgUmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMjI7XHJcbmNvbnN0IFJlc3RvcmVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDIzO1xyXG5jb25zdCBTZXRQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgPSAyNDtcclxuY29uc3QgU2V0QWx3YXlzT25Ub3BNZXRob2QgICAgICAgICAgICAgID0gMjU7XHJcbmNvbnN0IFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QgICAgICAgICA9IDI2O1xyXG5jb25zdCBTZXRGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgICAgPSAyNztcclxuY29uc3QgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWRNZXRob2QgID0gMjg7XHJcbmNvbnN0IFNldE1heFNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDI5O1xyXG5jb25zdCBTZXRNaW5TaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAzMDtcclxuY29uc3QgU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgID0gMzE7XHJcbmNvbnN0IFNldFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICA9IDMyO1xyXG5jb25zdCBTZXRTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzMztcclxuY29uc3QgU2V0VGl0bGVNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMzQ7XHJcbmNvbnN0IFNldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDM1O1xyXG5jb25zdCBTaG93TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNjtcclxuY29uc3QgU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzc7XHJcbmNvbnN0IFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICA9IDM4O1xyXG5jb25zdCBUb2dnbGVNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgPSAzOTtcclxuY29uc3QgVG9nZ2xlRnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgID0gNDA7IFxyXG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MTtcclxuY29uc3QgVW5NYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDI7XHJcbmNvbnN0IFVuTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQzO1xyXG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0NDtcclxuY29uc3QgWm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gNDU7XHJcbmNvbnN0IFpvb21Jbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDQ2O1xyXG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NztcclxuY29uc3QgWm9vbVJlc2V0TWV0aG9kICAgICAgICAgICAgICAgICAgID0gNDg7XHJcbmNvbnN0IFdpbmRvd0Ryb3Bab25lRHJvcHBlZCAgICAgICAgICAgICA9IDQ5O1xyXG5cclxuZnVuY3Rpb24gZ2V0RHJvcHpvbmVFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQgfCBudWxsKTogRWxlbWVudCB8IG51bGwge1xyXG4gICAgaWYgKCFlbGVtZW50KSB7XHJcbiAgICAgICAgcmV0dXJuIG51bGw7XHJcbiAgICB9XHJcbiAgICAvLyBBbGxvdyBkcm9wem9uZSBhdHRyaWJ1dGUgdG8gYmUgb24gdGhlIGVsZW1lbnQgaXRzZWxmIG9yIGFueSBwYXJlbnRcclxuICAgIHJldHVybiBlbGVtZW50LmNsb3Nlc3QoYFske0RST1BaT05FX0FUVFJJQlVURX1dYCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBwb3NpdGlvbiBvZiBhIHdpbmRvdy5cclxuICovXHJcbmludGVyZmFjZSBQb3NpdGlvbiB7XHJcbiAgICAvKiogVGhlIGhvcml6b250YWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIHg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgdmVydGljYWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIHk6IG51bWJlcjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXHJcbiAqL1xyXG5pbnRlcmZhY2UgU2l6ZSB7XHJcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICB3aWR0aDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIGhlaWdodDogbnVtYmVyO1xyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBjYWxsZXJTeW0gPSBTeW1ib2woXCJjYWxsZXJcIik7XHJcblxyXG5jbGFzcyBXaW5kb3cge1xyXG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXHJcbiAgICBwcml2YXRlIFtjYWxsZXJTeW1dOiAobWVzc2FnZTogbnVtYmVyLCBhcmdzPzogYW55KSA9PiBQcm9taXNlPGFueT47XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBJbml0aWFsaXNlcyBhIHdpbmRvdyBvYmplY3Qgd2l0aCB0aGUgc3BlY2lmaWVkIG5hbWUuXHJcbiAgICAgKlxyXG4gICAgICogQHByaXZhdGVcclxuICAgICAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHRhcmdldCB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IHN0cmluZyA9ICcnKSB7XHJcbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5XaW5kb3csIG5hbWUpXHJcblxyXG4gICAgICAgIC8vIGJpbmQgaW5zdGFuY2UgbWV0aG9kIHRvIG1ha2UgdGhlbSBlYXNpbHkgdXNhYmxlIGluIGV2ZW50IGhhbmRsZXJzXHJcbiAgICAgICAgZm9yIChjb25zdCBtZXRob2Qgb2YgT2JqZWN0LmdldE93blByb3BlcnR5TmFtZXMoV2luZG93LnByb3RvdHlwZSkpIHtcclxuICAgICAgICAgICAgaWYgKFxyXG4gICAgICAgICAgICAgICAgbWV0aG9kICE9PSBcImNvbnN0cnVjdG9yXCJcclxuICAgICAgICAgICAgICAgICYmIHR5cGVvZiAodGhpcyBhcyBhbnkpW21ldGhvZF0gPT09IFwiZnVuY3Rpb25cIlxyXG4gICAgICAgICAgICApIHtcclxuICAgICAgICAgICAgICAgICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9ICh0aGlzIGFzIGFueSlbbWV0aG9kXS5iaW5kKHRoaXMpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogR2V0cyB0aGUgc3BlY2lmaWVkIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gZ2V0LlxyXG4gICAgICogQHJldHVybnMgVGhlIGNvcnJlc3BvbmRpbmcgd2luZG93IG9iamVjdC5cclxuICAgICAqL1xyXG4gICAgR2V0KG5hbWU6IHN0cmluZyk6IFdpbmRvdyB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBXaW5kb3cobmFtZSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFBvc2l0aW9uTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdyBvbiB0aGUgc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBDZW50ZXIoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShDZW50ZXJNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIENsb3NlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2xvc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRGlzYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxyXG4gICAgICovXHJcbiAgICBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBFbmFibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cclxuICAgICAqL1xyXG4gICAgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgRm9jdXMoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb2N1c01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBGb3JjZXMgdGhlIHdpbmRvdyB0byByZWxvYWQgdGhlIHBhZ2UgYXNzZXRzLlxyXG4gICAgICovXHJcbiAgICBGb3JjZVJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvcmNlUmVsb2FkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFN3aXRjaGVzIHRoZSB3aW5kb3cgdG8gZnVsbHNjcmVlbiBtb2RlLlxyXG4gICAgICovXHJcbiAgICBGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBzY3JlZW4gdGhhdCB0aGUgd2luZG93IGlzIG9uLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgb24uXHJcbiAgICAgKi9cclxuICAgIEdldFNjcmVlbigpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0U2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHpvb20gbGV2ZWwuXHJcbiAgICAgKi9cclxuICAgIEdldFpvb20oKTogUHJvbWlzZTxudW1iZXI+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFpvb21NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIEhlaWdodCgpOiBQcm9taXNlPG51bWJlcj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSGVpZ2h0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEhpZGVzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIaWRlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZvY3VzZWQuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkLlxyXG4gICAgICovXHJcbiAgICBJc0ZvY3VzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0ZvY3VzZWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZnVsbHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW4uXHJcbiAgICAgKi9cclxuICAgIElzRnVsbHNjcmVlbigpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtYXhpbWlzZWQuXHJcbiAgICAgKi9cclxuICAgIElzTWF4aW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNYXhpbWlzZWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkLlxyXG4gICAgICovXHJcbiAgICBJc01pbmltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWluaW1pc2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE1heGltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBNaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1pbmltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBOYW1lKCk6IFByb21pc2U8c3RyaW5nPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShOYW1lTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxyXG4gICAgICovXHJcbiAgICBPcGVuRGV2VG9vbHMoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShPcGVuRGV2VG9vbHNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFJlbGF0aXZlUG9zaXRpb24oKTogUHJvbWlzZTxQb3NpdGlvbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZWxvYWRzIHRoZSBwYWdlIGFzc2V0cy5cclxuICAgICAqL1xyXG4gICAgUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsb2FkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IHJlc2l6YWJsZS5cclxuICAgICAqL1xyXG4gICAgUmVzaXphYmxlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzaXphYmxlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlc3RvcmVzIHRoZSB3aW5kb3cgdG8gaXRzIHByZXZpb3VzIHN0YXRlIGlmIGl0IHdhcyBwcmV2aW91c2x5IG1pbmltaXNlZCwgbWF4aW1pc2VkIG9yIGZ1bGxzY3JlZW4uXHJcbiAgICAgKi9cclxuICAgIFJlc3RvcmUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZXN0b3JlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcC5cclxuICAgICAqL1xyXG4gICAgU2V0QWx3YXlzT25Ub3AoYWx3YXlzT25Ub3A6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gciAtIFRoZSBkZXNpcmVkIHJlZCBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICogQHBhcmFtIGcgLSBUaGUgZGVzaXJlZCBncmVlbiBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICogQHBhcmFtIGIgLSBUaGUgZGVzaXJlZCBibHVlIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXHJcbiAgICAgKiBAcGFyYW0gYSAtIFRoZSBkZXNpcmVkIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXHJcbiAgICAgKi9cclxuICAgIFNldEJhY2tncm91bmRDb2xvdXIocjogbnVtYmVyLCBnOiBudW1iZXIsIGI6IG51bWJlciwgYTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIHRoZSB3aW5kb3cgZnJhbWUgYW5kIHRpdGxlIGJhci5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3MuXHJcbiAgICAgKi9cclxuICAgIFNldEZyYW1lbGVzcyhmcmFtZWxlc3M6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEaXNhYmxlcyB0aGUgc3lzdGVtIGZ1bGxzY3JlZW4gYnV0dG9uLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWQuXHJcbiAgICAgKi9cclxuICAgIFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkKGVuYWJsZWQ6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0TWF4U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWF4U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWluaW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldE1pblNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1pblNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0UmVsYXRpdmVQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cclxuICAgICAqL1xyXG4gICAgU2V0UmVzaXphYmxlKHJlc2l6YWJsZTogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSB0aXRsZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB0aXRsZSAtIFRoZSBkZXNpcmVkIHRpdGxlIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFRpdGxlKHRpdGxlOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFRpdGxlTWV0aG9kLCB7IHRpdGxlIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB6b29tIC0gVGhlIGRlc2lyZWQgem9vbSBsZXZlbC5cclxuICAgICAqL1xyXG4gICAgU2V0Wm9vbSh6b29tOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFpvb21NZXRob2QsIHsgem9vbSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNob3dzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTaG93TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNpemUoKTogUHJvbWlzZTxTaXplPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTaXplTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZ1bGxzY3JlZW4gYW5kIG5vcm1hbC5cclxuICAgICAqL1xyXG4gICAgVG9nZ2xlRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gbWF4aW1pc2VkIGFuZCBub3JtYWwuXHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZU1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnJhbWVsZXNzIGFuZCBub3JtYWwuXHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZUZyYW1lbGVzcygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZyYW1lbGVzc01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBVbkZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFVuTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBVbk1pbmltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NaW5pbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFdpZHRoKCk6IFByb21pc2U8bnVtYmVyPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShXaWR0aE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBab29tcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBab29tKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBJbmNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cclxuICAgICAqL1xyXG4gICAgWm9vbUluKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbUluTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIERlY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICovXHJcbiAgICBab29tT3V0KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU91dE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cclxuICAgICAqL1xyXG4gICAgWm9vbVJlc2V0KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbVJlc2V0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEhhbmRsZXMgZmlsZSBkcm9wcyBvcmlnaW5hdGluZyBmcm9tIHBsYXRmb3JtLXNwZWNpZmljIGNvZGUgKGUuZy4sIG1hY09TIG5hdGl2ZSBkcmFnLWFuZC1kcm9wKS5cclxuICAgICAqIEdhdGhlcnMgaW5mb3JtYXRpb24gYWJvdXQgdGhlIGRyb3AgdGFyZ2V0IGVsZW1lbnQgYW5kIHNlbmRzIGl0IGJhY2sgdG8gdGhlIEdvIGJhY2tlbmQuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGZpbGVuYW1lcyAtIEFuIGFycmF5IG9mIGZpbGUgcGF0aHMgKHN0cmluZ3MpIHRoYXQgd2VyZSBkcm9wcGVkLlxyXG4gICAgICogQHBhcmFtIHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxyXG4gICAgICogQHBhcmFtIHkgLSBUaGUgeS1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxyXG4gICAgICovXHJcbiAgICBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wKGZpbGVuYW1lczogc3RyaW5nW10sIHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XHJcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XHJcblxyXG4gICAgICAgIC8vIE5FVzogQ2hlY2sgaWYgdGhlIGRyb3AgdGFyZ2V0IGlzIGEgdmFsaWQgZHJvcHpvbmVcclxuICAgICAgICBjb25zdCBkcm9wem9uZVRhcmdldCA9IGdldERyb3B6b25lRWxlbWVudChlbGVtZW50KTtcclxuXHJcbiAgICAgICAgaWYgKCFkcm9wem9uZVRhcmdldCkge1xyXG4gICAgICAgICAgICBjb25zb2xlLmxvZyhgV2FpbHMgUnVudGltZTogRHJvcCBvbiBlbGVtZW50IChvciBubyBlbGVtZW50KSBhdCAke3h9LCR7eX0gd2hpY2ggaXMgbm90IGEgZGVzaWduYXRlZCBkcm9wem9uZS4gSWdub3JpbmcuIEVsZW1lbnQ6YCwgZWxlbWVudCk7XHJcbiAgICAgICAgICAgIC8vIE5vIG5lZWQgdG8gY2FsbCBiYWNrZW5kIGlmIG5vdCBhIHZhbGlkIGRyb3B6b25lIHRhcmdldFxyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBjb25zb2xlLmxvZyhgV2FpbHMgUnVudGltZTogRHJvcCBvbiBkZXNpZ25hdGVkIGRyb3B6b25lLiBFbGVtZW50IGF0ICgke3h9LCAke3l9KTpgLCBlbGVtZW50LCAnRWZmZWN0aXZlIGRyb3B6b25lOicsIGRyb3B6b25lVGFyZ2V0KTtcclxuICAgICAgICBjb25zdCBlbGVtZW50RGV0YWlscyA9IHtcclxuICAgICAgICAgICAgaWQ6IGRyb3B6b25lVGFyZ2V0LmlkLFxyXG4gICAgICAgICAgICBjbGFzc0xpc3Q6IEFycmF5LmZyb20oZHJvcHpvbmVUYXJnZXQuY2xhc3NMaXN0KSxcclxuICAgICAgICAgICAgYXR0cmlidXRlczoge30gYXMgeyBba2V5OiBzdHJpbmddOiBzdHJpbmcgfSxcclxuICAgICAgICB9O1xyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgZHJvcHpvbmVUYXJnZXQuYXR0cmlidXRlcy5sZW5ndGg7IGkrKykge1xyXG4gICAgICAgICAgICBjb25zdCBhdHRyID0gZHJvcHpvbmVUYXJnZXQuYXR0cmlidXRlc1tpXTtcclxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMuYXR0cmlidXRlc1thdHRyLm5hbWVdID0gYXR0ci52YWx1ZTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XHJcbiAgICAgICAgICAgIGZpbGVuYW1lcyxcclxuICAgICAgICAgICAgeCxcclxuICAgICAgICAgICAgeSxcclxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMsXHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dKFdpbmRvd0Ryb3Bab25lRHJvcHBlZCwgcGF5bG9hZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXHJcbiAqL1xyXG5jb25zdCB0aGlzV2luZG93ID0gbmV3IFdpbmRvdygnJyk7XHJcblxyXG4vLyBORVc6IEdsb2JhbCBEcmFnIEV2ZW50IExpc3RlbmVyc1xyXG5mdW5jdGlvbiBzZXR1cEdsb2JhbERyb3B6b25lTGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZG9jRWxlbWVudCA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudDtcclxuICAgIGxldCBkcmFnRW50ZXJDb3VudGVyID0gMDsgLy8gVG8gaGFuZGxlIGRyYWdlbnRlci9kcmFnbGVhdmUgb24gY2hpbGQgZWxlbWVudHNcclxuXHJcbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdlbnRlcicsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgaWYgKGV2ZW50LmRhdGFUcmFuc2ZlciAmJiBldmVudC5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcclxuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlcisrO1xyXG4gICAgICAgICAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZKTtcclxuICAgICAgICAgICAgY29uc3QgZHJvcHpvbmUgPSBnZXREcm9wem9uZUVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XHJcblxyXG4gICAgICAgICAgICAvLyBDbGVhciBwcmV2aW91cyBob3ZlciByZWdhcmRsZXNzLCB0aGVuIGFwcGx5IG5ldyBpZiB2YWxpZFxyXG4gICAgICAgICAgICBpZiAoY3VycmVudEhvdmVyZWREcm9wem9uZSAmJiBjdXJyZW50SG92ZXJlZERyb3B6b25lICE9PSBkcm9wem9uZSkge1xyXG4gICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QucmVtb3ZlKERST1BaT05FX0hPVkVSX0NMQVNTKTtcclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgaWYgKGRyb3B6b25lKSB7XHJcbiAgICAgICAgICAgICAgICBkcm9wem9uZS5jbGFzc0xpc3QuYWRkKERST1BaT05FX0hPVkVSX0NMQVNTKTtcclxuICAgICAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xyXG4gICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZSA9IGRyb3B6b25lO1xyXG4gICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XHJcbiAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lID0gbnVsbDsgLy8gRW5zdXJlIGl0J3MgY2xlYXJlZCBpZiBubyBkcm9wem9uZSBmb3VuZFxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG5cclxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCAoZXZlbnQpID0+IHtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBOZWNlc3NhcnkgdG8gYWxsb3cgZHJvcFxyXG4gICAgICAgIGlmIChldmVudC5kYXRhVHJhbnNmZXIgJiYgZXZlbnQuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XHJcbiAgICAgICAgICAgIC8vIE5vIG5lZWQgdG8gcXVlcnkgZWxlbWVudEZyb21Qb2ludCBhZ2FpbiBpZiBhbHJlYWR5IGhhbmRsZWQgYnkgZHJhZ2VudGVyIGNvcnJlY3RseVxyXG4gICAgICAgICAgICAvLyBKdXN0IGVuc3VyZSBkcm9wRWZmZWN0IGlzIGNvbnRpbnVvdXNseSBzZXQgYmFzZWQgb24gY3VycmVudEhvdmVyZWREcm9wem9uZVxyXG4gICAgICAgICAgICBpZiAoY3VycmVudEhvdmVyZWREcm9wem9uZSkge1xyXG4gICAgICAgICAgICAgICAgIC8vIFJlLWFwcGx5IGNsYXNzIGp1c3QgaW4gY2FzZSBpdCB3YXMgcmVtb3ZlZCBieSBzb21lIG90aGVyIEpTXHJcbiAgICAgICAgICAgICAgICBpZighY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QuY29udGFpbnMoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZS5jbGFzc0xpc3QuYWRkKERST1BaT05FX0hPVkVSX0NMQVNTKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xyXG4gICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9LCBmYWxzZSk7XHJcblxyXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCAoZXZlbnQpID0+IHtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIGlmIChldmVudC5kYXRhVHJhbnNmZXIgJiYgZXZlbnQuZGF0YVRyYW5zZmVyLnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XHJcbiAgICAgICAgICAgIGRyYWdFbnRlckNvdW50ZXItLTtcclxuICAgICAgICAgICAgLy8gT25seSByZW1vdmUgaG92ZXIgaWYgZHJhZyB0cnVseSBsZWZ0IHRoZSB3aW5kb3cgb3IgdGhlIGxhc3QgZHJvcHpvbmVcclxuICAgICAgICAgICAgaWYgKGRyYWdFbnRlckNvdW50ZXIgPT09IDAgfHwgZXZlbnQucmVsYXRlZFRhcmdldCA9PT0gbnVsbCB8fCAoY3VycmVudEhvdmVyZWREcm9wem9uZSAmJiAhY3VycmVudEhvdmVyZWREcm9wem9uZS5jb250YWlucyhldmVudC5yZWxhdGVkVGFyZ2V0IGFzIE5vZGUpKSkge1xyXG4gICAgICAgICAgICAgICAgaWYgKGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUpIHtcclxuICAgICAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lLmNsYXNzTGlzdC5yZW1vdmUoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpO1xyXG4gICAgICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgPSBudWxsO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7IC8vIFJlc2V0IGNvdW50ZXIgaWYgaXQgd2VudCBuZWdhdGl2ZSBvciBsZWZ0IHdpbmRvd1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG5cclxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJvcCcsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIFByZXZlbnQgZGVmYXVsdCBicm93c2VyIGZpbGUgaGFuZGxpbmdcclxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyID0gMDsgLy8gUmVzZXQgY291bnRlclxyXG4gICAgICAgIGlmIChjdXJyZW50SG92ZXJlZERyb3B6b25lKSB7XHJcbiAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY2xhc3NMaXN0LnJlbW92ZShEUk9QWk9ORV9IT1ZFUl9DTEFTUyk7XHJcbiAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgICAgICAvLyBUaGUgYWN0dWFsIGRyb3AgcHJvY2Vzc2luZyBpcyBpbml0aWF0ZWQgYnkgdGhlIG5hdGl2ZSBzaWRlIGNhbGxpbmcgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcFxyXG4gICAgICAgIC8vIEhhbmRsZVBsYXRmb3JtRmlsZURyb3Agd2lsbCB0aGVuIGNoZWNrIGlmIHRoZSBkcm9wIHdhcyBvbiBhIHZhbGlkIHpvbmUuXHJcbiAgICB9LCBmYWxzZSk7XHJcbn1cclxuXHJcbi8vIEluaXRpYWxpemUgbGlzdGVuZXJzIHdoZW4gdGhlIHNjcmlwdCBsb2Fkc1xyXG5pZiAodHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgZG9jdW1lbnQgIT09IFwidW5kZWZpbmVkXCIpIHtcclxuICAgIHNldHVwR2xvYmFsRHJvcHpvbmVMaXN0ZW5lcnMoKTtcclxufVxyXG5cclxuZXhwb3J0IGRlZmF1bHQgdGhpc1dpbmRvdztcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XHJcblxyXG4vLyBOT1RFOiB0aGUgZm9sbG93aW5nIG1ldGhvZHMgTVVTVCBiZSBpbXBvcnRlZCBleHBsaWNpdGx5IGJlY2F1c2Ugb2YgaG93IGVzYnVpbGQgaW5qZWN0aW9uIHdvcmtzXHJcbmltcG9ydCB7IEVuYWJsZSBhcyBFbmFibGVXTUwgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XHJcbmltcG9ydCB7IGRlYnVnTG9nIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzXCI7XHJcblxyXG53aW5kb3cud2FpbHMgPSBSdW50aW1lO1xyXG5FbmFibGVXTUwoKTtcclxuXHJcbmlmIChERUJVRykge1xyXG4gICAgZGVidWdMb2coXCJXYWlscyBSdW50aW1lIExvYWRlZFwiKVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlN5c3RlbSk7XHJcblxyXG5jb25zdCBTeXN0ZW1Jc0RhcmtNb2RlID0gMDtcclxuY29uc3QgU3lzdGVtRW52aXJvbm1lbnQgPSAxO1xyXG5jb25zdCBBcHBsaWNhdGlvbkZpbGVzRHJvcHBlZFdpdGhDb250ZXh0ID0gMTAwOyAvLyBOZXcgbWV0aG9kIElEIGZvciBlbnJpY2hlZCBkcm9wIGV2ZW50XHJcblxyXG5jb25zdCBfaW52b2tlID0gKGZ1bmN0aW9uICgpIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlKSB7XHJcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcpO1xyXG4gICAgICAgIH0gZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndlYmtpdD8ubWVzc2FnZUhhbmRsZXJzPy5bJ2V4dGVybmFsJ10/LnBvc3RNZXNzYWdlKSB7XHJcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXS5wb3N0TWVzc2FnZS5iaW5kKCh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddKTtcclxuICAgICAgICB9XHJcbiAgICB9IGNhdGNoKGUpIHt9XHJcblxyXG4gICAgY29uc29sZS53YXJuKCdcXG4lY1x1MjZBMFx1RkUwRiBCcm93c2VyIEVudmlyb25tZW50IERldGVjdGVkICVjXFxuXFxuJWNPbmx5IFVJIHByZXZpZXdzIGFyZSBhdmFpbGFibGUgaW4gdGhlIGJyb3dzZXIuIEZvciBmdWxsIGZ1bmN0aW9uYWxpdHksIHBsZWFzZSBydW4gdGhlIGFwcGxpY2F0aW9uIGluIGRlc2t0b3AgbW9kZS5cXG5Nb3JlIGluZm9ybWF0aW9uIGF0OiBodHRwczovL3YzLndhaWxzLmlvL2xlYXJuL2J1aWxkLyN1c2luZy1hLWJyb3dzZXItZm9yLWRldmVsb3BtZW50XFxuJyxcclxuICAgICAgICAnYmFja2dyb3VuZDogI2ZmZmZmZjsgY29sb3I6ICMwMDAwMDA7IGZvbnQtd2VpZ2h0OiBib2xkOyBwYWRkaW5nOiA0cHggOHB4OyBib3JkZXItcmFkaXVzOiA0cHg7IGJvcmRlcjogMnB4IHNvbGlkICMwMDAwMDA7JyxcclxuICAgICAgICAnYmFja2dyb3VuZDogdHJhbnNwYXJlbnQ7JyxcclxuICAgICAgICAnY29sb3I6ICNmZmZmZmY7IGZvbnQtc3R5bGU6IGl0YWxpYzsgZm9udC13ZWlnaHQ6IGJvbGQ7Jyk7XHJcbiAgICByZXR1cm4gbnVsbDtcclxufSkoKTtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBpbnZva2UobXNnOiBhbnkpOiB2b2lkIHtcclxuICAgIF9pbnZva2U/Lihtc2cpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0cmlldmVzIHRoZSBzeXN0ZW0gZGFyayBtb2RlIHN0YXR1cy5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYSBib29sZWFuIHZhbHVlIGluZGljYXRpbmcgaWYgdGhlIHN5c3RlbSBpcyBpbiBkYXJrIG1vZGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNEYXJrTW9kZSgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgIHJldHVybiBjYWxsKFN5c3RlbUlzRGFya01vZGUpO1xyXG59XHJcblxyXG4vKipcclxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIHRoZSBjYXBhYmlsaXRpZXMuXHJcbiAqL1xyXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gQ2FwYWJpbGl0aWVzKCk6IFByb21pc2U8UmVjb3JkPHN0cmluZywgYW55Pj4ge1xyXG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2goXCIvd2FpbHMvY2FwYWJpbGl0aWVzXCIpO1xyXG4gICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiY291bGQgbm90IGZldGNoIGNhcGFiaWxpdGllczogXCIgKyByZXNwb25zZS5zdGF0dXNUZXh0KTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBPU0luZm8ge1xyXG4gICAgLyoqIFRoZSBicmFuZGluZyBvZiB0aGUgT1MuICovXHJcbiAgICBCcmFuZGluZzogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSBJRCBvZiB0aGUgT1MuICovXHJcbiAgICBJRDogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSBuYW1lIG9mIHRoZSBPUy4gKi9cclxuICAgIE5hbWU6IHN0cmluZztcclxuICAgIC8qKiBUaGUgdmVyc2lvbiBvZiB0aGUgT1MuICovXHJcbiAgICBWZXJzaW9uOiBzdHJpbmc7XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgRW52aXJvbm1lbnRJbmZvIHtcclxuICAgIC8qKiBUaGUgYXJjaGl0ZWN0dXJlIG9mIHRoZSBzeXN0ZW0uICovXHJcbiAgICBBcmNoOiBzdHJpbmc7XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYXBwbGljYXRpb24gaXMgcnVubmluZyBpbiBkZWJ1ZyBtb2RlLCBvdGhlcndpc2UgZmFsc2UuICovXHJcbiAgICBEZWJ1ZzogYm9vbGVhbjtcclxuICAgIC8qKiBUaGUgb3BlcmF0aW5nIHN5c3RlbSBpbiB1c2UuICovXHJcbiAgICBPUzogc3RyaW5nO1xyXG4gICAgLyoqIERldGFpbHMgb2YgdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXHJcbiAgICBPU0luZm86IE9TSW5mbztcclxuICAgIC8qKiBBZGRpdGlvbmFsIHBsYXRmb3JtIGluZm9ybWF0aW9uLiAqL1xyXG4gICAgUGxhdGZvcm1JbmZvOiBSZWNvcmQ8c3RyaW5nLCBhbnk+O1xyXG59XHJcblxyXG4vKipcclxuICogUmV0cmlldmVzIGVudmlyb25tZW50IGRldGFpbHMuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIE9TIGFuZCBzeXN0ZW0gYXJjaGl0ZWN0dXJlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCk6IFByb21pc2U8RW52aXJvbm1lbnRJbmZvPiB7XHJcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1FbnZpcm9ubWVudCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxyXG4gKlxyXG4gKiBAcmV0dXJuIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzV2luZG93cygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcIndpbmRvd3NcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBSZXR1cm5zIHRydWUgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJsaW51eFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Lk9TID09PSBcImRhcndpblwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FNRDY0KCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhbWQ2NFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm1cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuQXJjaCA9PT0gXCJhcm02NFwiO1xyXG59XHJcblxyXG4vKipcclxuICogUmVwb3J0cyB3aGV0aGVyIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGFwcCBpcyBiZWluZyBydW4gaW4gZGVidWcgbW9kZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RlYnVnKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIEJvb2xlYW4od2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5EZWJ1Zyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUyBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXHJcbiAqIEdhdGhlcnMgaW5mb3JtYXRpb24gYWJvdXQgdGhlIGRyb3AgdGFyZ2V0IGVsZW1lbnQgYW5kIHNlbmRzIGl0IGJhY2sgdG8gdGhlIEdvIGJhY2tlbmQuXHJcbiAqXHJcbiAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cclxuICogQHBhcmFtIHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxyXG4gKiBAcGFyYW0geSAtIFRoZSB5LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xyXG4gICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XHJcbiAgICBjb25zdCBlbGVtZW50SWQgPSBlbGVtZW50ID8gZWxlbWVudC5pZCA6ICcnO1xyXG4gICAgY29uc3QgY2xhc3NMaXN0ID0gZWxlbWVudCA/IEFycmF5LmZyb20oZWxlbWVudC5jbGFzc0xpc3QpIDogW107XHJcblxyXG4gICAgY29uc3QgcGF5bG9hZCA9IHtcclxuICAgICAgICBmaWxlbmFtZXMsXHJcbiAgICAgICAgeCxcclxuICAgICAgICB5LFxyXG4gICAgICAgIGVsZW1lbnRJZCxcclxuICAgICAgICBjbGFzc0xpc3QsXHJcbiAgICB9O1xyXG5cclxuICAgIGNhbGwoQXBwbGljYXRpb25GaWxlc0Ryb3BwZWRXaXRoQ29udGV4dCwgcGF5bG9hZClcclxuICAgICAgICAudGhlbigoKSA9PiB7XHJcbiAgICAgICAgICAgIC8vIE9wdGlvbmFsOiBMb2cgc3VjY2VzcyBvciBoYW5kbGUgaWYgbmVlZGVkXHJcbiAgICAgICAgICAgIGNvbnNvbGUubG9nKFwiUGxhdGZvcm0gZmlsZSBkcm9wIHByb2Nlc3NlZCBhbmQgc2VudCB0byBHby5cIik7XHJcbiAgICAgICAgfSlcclxuICAgICAgICAuY2F0Y2goZXJyID0+IHtcclxuICAgICAgICAgICAgLy8gT3B0aW9uYWw6IExvZyBlcnJvclxyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igc2VuZGluZyBwbGF0Zm9ybSBmaWxlIGRyb3AgdG8gR286XCIsIGVycik7XHJcbiAgICAgICAgfSk7XHJcbn1cclxuXHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgSXNEZWJ1ZyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xyXG5pbXBvcnQgeyBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzXCI7XHJcblxyXG4vLyBzZXR1cFxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xyXG5cclxuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcclxuXHJcbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZDogc3RyaW5nLCB4OiBudW1iZXIsIHk6IG51bWJlciwgZGF0YTogYW55KTogdm9pZCB7XHJcbiAgICB2b2lkIGNhbGwoQ29udGV4dE1lbnVPcGVuLCB7aWQsIHgsIHksIGRhdGF9KTtcclxufVxyXG5cclxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50OiBNb3VzZUV2ZW50KSB7XHJcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XHJcblxyXG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcclxuICAgIGNvbnN0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIikudHJpbSgpO1xyXG5cclxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgY29uc3QgZGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XHJcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBkYXRhKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCwgdGFyZ2V0KTtcclxuICAgIH1cclxufVxyXG5cclxuXHJcbi8qXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG5cclxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcclxuKi9cclxuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudDogTW91c2VFdmVudCwgdGFyZ2V0OiBIVE1MRWxlbWVudCkge1xyXG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XHJcbiAgICBpZiAoSXNEZWJ1ZygpKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgIHN3aXRjaCAod2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKSkge1xyXG4gICAgICAgIGNhc2UgJ3Nob3cnOlxyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgY2FzZSAnaGlkZSc6XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxyXG4gICAgaWYgKHRhcmdldC5pc0NvbnRlbnRFZGl0YWJsZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXHJcbiAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XHJcbiAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSBzZWxlY3Rpb24gJiYgc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMDtcclxuICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcclxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcclxuICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcclxuICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xyXG4gICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XHJcbiAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XHJcbiAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gdGFyZ2V0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8vIENoZWNrIGlmIHRhZyBpcyBpbnB1dCBvciB0ZXh0YXJlYS5cclxuICAgIGlmICh0YXJnZXQgaW5zdGFuY2VvZiBIVE1MSW5wdXRFbGVtZW50IHx8IHRhcmdldCBpbnN0YW5jZW9mIEhUTUxUZXh0QXJlYUVsZW1lbnQpIHtcclxuICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghdGFyZ2V0LnJlYWRPbmx5ICYmICF0YXJnZXQuZGlzYWJsZWQpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG4gICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyB0aGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5IGZyb20gdGhlIGZsYWcgbWFwLlxyXG4gKlxyXG4gKiBAcGFyYW0ga2V5IC0gVGhlIGtleSB0byByZXRyaWV2ZSB0aGUgdmFsdWUgZm9yLlxyXG4gKiBAcmV0dXJuIFRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXk6IHN0cmluZyk6IGFueSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleV07XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiVW5hYmxlIHRvIHJldHJpZXZlIGZsYWcgJ1wiICsga2V5ICsgXCInOiBcIiArIGUsIHsgY2F1c2U6IGUgfSk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IGludm9rZSwgSXNXaW5kb3dzIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XHJcbmltcG9ydCB7IEdldEZsYWcgfSBmcm9tIFwiLi9mbGFncy5qc1wiO1xyXG5pbXBvcnQgeyBjYW5UcmFja0J1dHRvbnMsIGV2ZW50VGFyZ2V0IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcclxuXHJcbi8vIFNldHVwXHJcbmxldCBjYW5EcmFnID0gZmFsc2U7XHJcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xyXG5cclxubGV0IHJlc2l6YWJsZSA9IGZhbHNlO1xyXG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XHJcbmxldCByZXNpemluZyA9IGZhbHNlO1xyXG5sZXQgcmVzaXplRWRnZTogc3RyaW5nID0gXCJcIjtcclxubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcclxuXHJcbmxldCBidXR0b25zID0gMDtcclxuY29uc3QgYnV0dG9uc1RyYWNrZWQgPSBjYW5UcmFja0J1dHRvbnMoKTtcclxuXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9ICh2YWx1ZTogYm9vbGVhbik6IHZvaWQgPT4ge1xyXG4gICAgcmVzaXphYmxlID0gdmFsdWU7XHJcbiAgICBpZiAoIXJlc2l6YWJsZSkge1xyXG4gICAgICAgIC8vIFN0b3AgcmVzaXppbmcgaWYgaW4gcHJvZ3Jlc3MuXHJcbiAgICAgICAgY2FuUmVzaXplID0gcmVzaXppbmcgPSBmYWxzZTtcclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH1cclxufTtcclxuXHJcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xyXG5mb3IgKGNvbnN0IGV2IG9mIFsnY2xpY2snLCAnY29udGV4dG1lbnUnLCAnZGJsY2xpY2snXSkge1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoZXYsIHN1cHByZXNzRXZlbnQsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcclxufVxyXG5cclxuZnVuY3Rpb24gc3VwcHJlc3NFdmVudChldmVudDogRXZlbnQpIHtcclxuICAgIC8vIFN1cHByZXNzIGNsaWNrIGV2ZW50cyB3aGlsZSByZXNpemluZyBvciBkcmFnZ2luZy5cclxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xyXG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xyXG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIFVzZSBjb25zdGFudHMgdG8gYXZvaWQgY29tcGFyaW5nIHN0cmluZ3MgbXVsdGlwbGUgdGltZXMuXHJcbmNvbnN0IE1vdXNlRG93biA9IDA7XHJcbmNvbnN0IE1vdXNlVXAgICA9IDE7XHJcbmNvbnN0IE1vdXNlTW92ZSA9IDI7XHJcblxyXG5mdW5jdGlvbiB1cGRhdGUoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcclxuICAgIC8vIFdpbmRvd3Mgc3VwcHJlc3NlcyBtb3VzZSBldmVudHMgYXQgdGhlIGVuZCBvZiBkcmFnZ2luZyBvciByZXNpemluZyxcclxuICAgIC8vIHNvIHdlIG5lZWQgdG8gYmUgc21hcnQgYW5kIHN5bnRoZXNpemUgYnV0dG9uIGV2ZW50cy5cclxuXHJcbiAgICBsZXQgZXZlbnRUeXBlOiBudW1iZXIsIGV2ZW50QnV0dG9ucyA9IGV2ZW50LmJ1dHRvbnM7XHJcbiAgICBzd2l0Y2ggKGV2ZW50LnR5cGUpIHtcclxuICAgICAgICBjYXNlICdtb3VzZWRvd24nOlxyXG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZURvd247XHJcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyB8ICgxIDw8IGV2ZW50LmJ1dHRvbik7IH1cclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgY2FzZSAnbW91c2V1cCc6XHJcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlVXA7XHJcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyAmIH4oMSA8PCBldmVudC5idXR0b24pOyB9XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIGRlZmF1bHQ6XHJcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlTW92ZTtcclxuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zOyB9XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgfVxyXG5cclxuICAgIGxldCByZWxlYXNlZCA9IGJ1dHRvbnMgJiB+ZXZlbnRCdXR0b25zO1xyXG4gICAgbGV0IHByZXNzZWQgPSBldmVudEJ1dHRvbnMgJiB+YnV0dG9ucztcclxuXHJcbiAgICBidXR0b25zID0gZXZlbnRCdXR0b25zO1xyXG5cclxuICAgIC8vIFN5bnRoZXNpemUgYSByZWxlYXNlLXByZXNzIHNlcXVlbmNlIGlmIHdlIGRldGVjdCBhIHByZXNzIG9mIGFuIGFscmVhZHkgcHJlc3NlZCBidXR0b24uXHJcbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZURvd24gJiYgIShwcmVzc2VkICYgZXZlbnQuYnV0dG9uKSkge1xyXG4gICAgICAgIHJlbGVhc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XHJcbiAgICAgICAgcHJlc3NlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFN1cHByZXNzIGFsbCBidXR0b24gZXZlbnRzIGR1cmluZyBkcmFnZ2luZyBhbmQgcmVzaXppbmcsXHJcbiAgICAvLyB1bmxlc3MgdGhpcyBpcyBhIG1vdXNldXAgZXZlbnQgdGhhdCBpcyBlbmRpbmcgYSBkcmFnIGFjdGlvbi5cclxuICAgIGlmIChcclxuICAgICAgICBldmVudFR5cGUgIT09IE1vdXNlTW92ZSAvLyBGYXN0IHBhdGggZm9yIG1vdXNlbW92ZVxyXG4gICAgICAgICYmIHJlc2l6aW5nXHJcbiAgICAgICAgfHwgKFxyXG4gICAgICAgICAgICBkcmFnZ2luZ1xyXG4gICAgICAgICAgICAmJiAoXHJcbiAgICAgICAgICAgICAgICBldmVudFR5cGUgPT09IE1vdXNlRG93blxyXG4gICAgICAgICAgICAgICAgfHwgZXZlbnQuYnV0dG9uICE9PSAwXHJcbiAgICAgICAgICAgIClcclxuICAgICAgICApXHJcbiAgICApIHtcclxuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIEhhbmRsZSByZWxlYXNlc1xyXG4gICAgaWYgKHJlbGVhc2VkICYgMSkgeyBwcmltYXJ5VXAoZXZlbnQpOyB9XHJcbiAgICAvLyBIYW5kbGUgcHJlc3Nlc1xyXG4gICAgaWYgKHByZXNzZWQgJiAxKSB7IHByaW1hcnlEb3duKGV2ZW50KTsgfVxyXG5cclxuICAgIC8vIEhhbmRsZSBtb3VzZW1vdmVcclxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlTW92ZSkgeyBvbk1vdXNlTW92ZShldmVudCk7IH07XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHByaW1hcnlEb3duKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XHJcbiAgICAvLyBSZXNldCByZWFkaW5lc3Mgc3RhdGUuXHJcbiAgICBjYW5EcmFnID0gZmFsc2U7XHJcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcclxuXHJcbiAgICAvLyBJZ25vcmUgcmVwZWF0ZWQgY2xpY2tzIG9uIG1hY09TIGFuZCBMaW51eC5cclxuICAgIGlmICghSXNXaW5kb3dzKCkpIHtcclxuICAgICAgICBpZiAoZXZlbnQudHlwZSA9PT0gJ21vdXNlZG93bicgJiYgZXZlbnQuYnV0dG9uID09PSAwICYmIGV2ZW50LmRldGFpbCAhPT0gMSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIGlmIChyZXNpemVFZGdlKSB7XHJcbiAgICAgICAgLy8gUmVhZHkgdG8gcmVzaXplIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUuXHJcbiAgICAgICAgY2FuUmVzaXplID0gdHJ1ZTtcclxuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZyBvcGVyYXRpb25zIHdoZW4gb24gcmVzaXplIGVkZ2VzLlxyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBSZXRyaWV2ZSB0YXJnZXQgZWxlbWVudFxyXG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xyXG5cclxuICAgIC8vIFJlYWR5IHRvIGRyYWcgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZSBvbiBhIGRyYWdnYWJsZSBlbGVtZW50LlxyXG4gICAgLy8gSWdub3JlIGNsaWNrcyBvbiB0aGUgc2Nyb2xsYmFyLlxyXG4gICAgY29uc3Qgc3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpO1xyXG4gICAgY2FuRHJhZyA9IChcclxuICAgICAgICBzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13YWlscy1kcmFnZ2FibGVcIikudHJpbSgpID09PSBcImRyYWdcIlxyXG4gICAgICAgICYmIChcclxuICAgICAgICAgICAgZXZlbnQub2Zmc2V0WCAtIHBhcnNlRmxvYXQoc3R5bGUucGFkZGluZ0xlZnQpIDwgdGFyZ2V0LmNsaWVudFdpZHRoXHJcbiAgICAgICAgICAgICYmIGV2ZW50Lm9mZnNldFkgLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdUb3ApIDwgdGFyZ2V0LmNsaWVudEhlaWdodFxyXG4gICAgICAgIClcclxuICAgICk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHByaW1hcnlVcChldmVudDogTW91c2VFdmVudCkge1xyXG4gICAgLy8gU3RvcCBkcmFnZ2luZyBhbmQgcmVzaXppbmcuXHJcbiAgICBjYW5EcmFnID0gZmFsc2U7XHJcbiAgICBkcmFnZ2luZyA9IGZhbHNlO1xyXG4gICAgY2FuUmVzaXplID0gZmFsc2U7XHJcbiAgICByZXNpemluZyA9IGZhbHNlO1xyXG59XHJcblxyXG5jb25zdCBjdXJzb3JGb3JFZGdlID0gT2JqZWN0LmZyZWV6ZSh7XHJcbiAgICBcInNlLXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXHJcbiAgICBcInN3LXJlc2l6ZVwiOiBcIm5lc3ctcmVzaXplXCIsXHJcbiAgICBcIm53LXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXHJcbiAgICBcIm5lLXJlc2l6ZVwiOiBcIm5lc3ctcmVzaXplXCIsXHJcbiAgICBcInctcmVzaXplXCI6IFwiZXctcmVzaXplXCIsXHJcbiAgICBcIm4tcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXHJcbiAgICBcInMtcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXHJcbiAgICBcImUtcmVzaXplXCI6IFwiZXctcmVzaXplXCIsXHJcbn0pXHJcblxyXG5mdW5jdGlvbiBzZXRSZXNpemUoZWRnZT86IGtleW9mIHR5cGVvZiBjdXJzb3JGb3JFZGdlKTogdm9pZCB7XHJcbiAgICBpZiAoZWRnZSkge1xyXG4gICAgICAgIGlmICghcmVzaXplRWRnZSkgeyBkZWZhdWx0Q3Vyc29yID0gZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3I7IH1cclxuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGN1cnNvckZvckVkZ2VbZWRnZV07XHJcbiAgICB9IGVsc2UgaWYgKCFlZGdlICYmIHJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGRlZmF1bHRDdXJzb3I7XHJcbiAgICB9XHJcblxyXG4gICAgcmVzaXplRWRnZSA9IGVkZ2UgfHwgXCJcIjtcclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZU1vdmUoZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcclxuICAgIGlmIChjYW5SZXNpemUgJiYgcmVzaXplRWRnZSkge1xyXG4gICAgICAgIC8vIFN0YXJ0IHJlc2l6aW5nLlxyXG4gICAgICAgIHJlc2l6aW5nID0gdHJ1ZTtcclxuICAgICAgICBpbnZva2UoXCJ3YWlsczpyZXNpemU6XCIgKyByZXNpemVFZGdlKTtcclxuICAgIH0gZWxzZSBpZiAoY2FuRHJhZykge1xyXG4gICAgICAgIC8vIFN0YXJ0IGRyYWdnaW5nLlxyXG4gICAgICAgIGRyYWdnaW5nID0gdHJ1ZTtcclxuICAgICAgICBpbnZva2UoXCJ3YWlsczpkcmFnXCIpO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xyXG4gICAgICAgIC8vIEVpdGhlciBkcmFnIG9yIHJlc2l6ZSBpcyBvbmdvaW5nLFxyXG4gICAgICAgIC8vIHJlc2V0IHJlYWRpbmVzcyBhbmQgc3RvcCBwcm9jZXNzaW5nLlxyXG4gICAgICAgIGNhbkRyYWcgPSBjYW5SZXNpemUgPSBmYWxzZTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKCFyZXNpemFibGUgfHwgIUlzV2luZG93cygpKSB7XHJcbiAgICAgICAgaWYgKHJlc2l6ZUVkZ2UpIHsgc2V0UmVzaXplKCk7IH1cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcclxuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xyXG5cclxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhcy5cclxuICAgIGNvbnN0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xyXG5cclxuICAgIGNvbnN0IHJpZ2h0Qm9yZGVyID0gKHdpbmRvdy5vdXRlcldpZHRoIC0gZXZlbnQuY2xpZW50WCkgPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGNvbnN0IGxlZnRCb3JkZXIgPSBldmVudC5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XHJcbiAgICBjb25zdCB0b3BCb3JkZXIgPSBldmVudC5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xyXG4gICAgY29uc3QgYm90dG9tQm9yZGVyID0gKHdpbmRvdy5vdXRlckhlaWdodCAtIGV2ZW50LmNsaWVudFkpIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xyXG5cclxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyIGFyZWFzLlxyXG4gICAgY29uc3QgcmlnaHRDb3JuZXIgPSAod2luZG93Lm91dGVyV2lkdGggLSBldmVudC5jbGllbnRYKSA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGNvbnN0IGxlZnRDb3JuZXIgPSBldmVudC5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xyXG4gICAgY29uc3QgdG9wQ29ybmVyID0gZXZlbnQuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBjb25zdCBib3R0b21Db3JuZXIgPSAod2luZG93Lm91dGVySGVpZ2h0IC0gZXZlbnQuY2xpZW50WSkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG5cclxuICAgIGlmICghbGVmdENvcm5lciAmJiAhdG9wQ29ybmVyICYmICFib3R0b21Db3JuZXIgJiYgIXJpZ2h0Q29ybmVyKSB7XHJcbiAgICAgICAgLy8gT3B0aW1pc2F0aW9uOiBvdXQgb2YgYWxsIGNvcm5lciBhcmVhcyBpbXBsaWVzIG91dCBvZiBib3JkZXJzLlxyXG4gICAgICAgIHNldFJlc2l6ZSgpO1xyXG4gICAgfVxyXG4gICAgLy8gRGV0ZWN0IGNvcm5lcnMuXHJcbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIHRvcENvcm5lcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XHJcbiAgICAvLyBEZXRlY3QgYm9yZGVycy5cclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XHJcbiAgICAvLyBPdXQgb2YgYm9yZGVyIGFyZWEuXHJcbiAgICBlbHNlIHNldFJlc2l6ZSgpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQXBwbGljYXRpb24pO1xyXG5cclxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XHJcbmNvbnN0IFNob3dNZXRob2QgPSAxO1xyXG5jb25zdCBRdWl0TWV0aG9kID0gMjtcclxuXHJcbi8qKlxyXG4gKiBIaWRlcyBhIGNlcnRhaW4gbWV0aG9kIGJ5IGNhbGxpbmcgdGhlIEhpZGVNZXRob2QgZnVuY3Rpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKEhpZGVNZXRob2QpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoU2hvd01ldGhvZCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKFF1aXRNZXRob2QpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBDYW5jZWxsYWJsZVByb21pc2UsIHR5cGUgQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVycyB9IGZyb20gXCIuL2NhbmNlbGxhYmxlLmpzXCI7XHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcclxuXHJcbi8vIFNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyID0gcmVzdWx0SGFuZGxlcjtcclxud2luZG93Ll93YWlscy5jYWxsRXJyb3JIYW5kbGVyID0gZXJyb3JIYW5kbGVyO1xyXG5cclxudHlwZSBQcm9taXNlUmVzb2x2ZXJzID0gT21pdDxDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPGFueT4sIFwicHJvbWlzZVwiIHwgXCJvbmNhbmNlbGxlZFwiPlxyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2FsbCk7XHJcbmNvbnN0IGNhbmNlbENhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbmNlbENhbGwpO1xyXG5jb25zdCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcDxzdHJpbmcsIFByb21pc2VSZXNvbHZlcnM+KCk7XHJcblxyXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XHJcbmNvbnN0IENhbmNlbE1ldGhvZCA9IDBcclxuXHJcbi8qKlxyXG4gKiBIb2xkcyBhbGwgcmVxdWlyZWQgaW5mb3JtYXRpb24gZm9yIGEgYmluZGluZyBjYWxsLlxyXG4gKiBNYXkgcHJvdmlkZSBlaXRoZXIgYSBtZXRob2QgSUQgb3IgYSBtZXRob2QgbmFtZSwgYnV0IG5vdCBib3RoLlxyXG4gKi9cclxuZXhwb3J0IHR5cGUgQ2FsbE9wdGlvbnMgPSB7XHJcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xyXG4gICAgbWV0aG9kSUQ6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgZnVsbHkgcXVhbGlmaWVkIG5hbWUgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xyXG4gICAgbWV0aG9kTmFtZT86IG5ldmVyO1xyXG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xyXG4gICAgYXJnczogYW55W107XHJcbn0gfCB7XHJcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xyXG4gICAgbWV0aG9kSUQ/OiBuZXZlcjtcclxuICAgIC8qKiBUaGUgZnVsbHkgcXVhbGlmaWVkIG5hbWUgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xyXG4gICAgbWV0aG9kTmFtZTogc3RyaW5nO1xyXG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xyXG4gICAgYXJnczogYW55W107XHJcbn07XHJcblxyXG4vKipcclxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB0aHJvd24gaW4gY2FzZSB0aGUgYm91bmQgbWV0aG9kIHJldHVybnMgYW4gZXJyb3IuXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIFJ1bnRpbWVFcnJvciNuYW1lfSBwcm9wZXJ0eSBpcyBcIlJ1bnRpbWVFcnJvclwiLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIFJ1bnRpbWVFcnJvciBleHRlbmRzIEVycm9yIHtcclxuICAgIC8qKlxyXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBSdW50aW1lRXJyb3IgaW5zdGFuY2UuXHJcbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxyXG4gICAgICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIHRvIGJlIGZvcndhcmRlZCB0byB0aGUgRXJyb3IgY29uc3RydWN0b3IuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcclxuICAgICAgICBzdXBlcihtZXNzYWdlLCBvcHRpb25zKTtcclxuICAgICAgICB0aGlzLm5hbWUgPSBcIlJ1bnRpbWVFcnJvclwiO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogSGFuZGxlcyB0aGUgcmVzdWx0IG9mIGEgY2FsbCByZXF1ZXN0LlxyXG4gKlxyXG4gKiBAcGFyYW0gaWQgLSBUaGUgaWQgb2YgdGhlIHJlcXVlc3QgdG8gaGFuZGxlIHRoZSByZXN1bHQgZm9yLlxyXG4gKiBAcGFyYW0gZGF0YSAtIFRoZSByZXN1bHQgZGF0YSBvZiB0aGUgcmVxdWVzdC5cclxuICogQHBhcmFtIGlzSlNPTiAtIEluZGljYXRlcyB3aGV0aGVyIHRoZSBkYXRhIGlzIEpTT04gb3Igbm90LlxyXG4gKi9cclxuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZDogc3RyaW5nLCBkYXRhOiBzdHJpbmcsIGlzSlNPTjogYm9vbGVhbik6IHZvaWQge1xyXG4gICAgY29uc3QgcmVzb2x2ZXJzID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xyXG4gICAgaWYgKCFyZXNvbHZlcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKCFkYXRhKSB7XHJcbiAgICAgICAgcmVzb2x2ZXJzLnJlc29sdmUodW5kZWZpbmVkKTtcclxuICAgIH0gZWxzZSBpZiAoIWlzSlNPTikge1xyXG4gICAgICAgIHJlc29sdmVycy5yZXNvbHZlKGRhdGEpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcclxuICAgICAgICB9IGNhdGNoIChlcnI6IGFueSkge1xyXG4gICAgICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJjb3VsZCBub3QgcGFyc2UgcmVzdWx0OiBcIiArIGVyci5tZXNzYWdlLCB7IGNhdXNlOiBlcnIgfSkpO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEhhbmRsZXMgdGhlIGVycm9yIGZyb20gYSBjYWxsIHJlcXVlc3QuXHJcbiAqXHJcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxyXG4gKiBAcGFyYW0gZGF0YSAtIFRoZSBlcnJvciBkYXRhIHRvIHJlamVjdCB0aGUgcHJvbWlzZSBoYW5kbGVyIHdpdGguXHJcbiAqIEBwYXJhbSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cclxuICovXHJcbmZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZDogc3RyaW5nLCBkYXRhOiBzdHJpbmcsIGlzSlNPTjogYm9vbGVhbik6IHZvaWQge1xyXG4gICAgY29uc3QgcmVzb2x2ZXJzID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xyXG4gICAgaWYgKCFyZXNvbHZlcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKCFpc0pTT04pIHtcclxuICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBFcnJvcihkYXRhKSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGxldCBlcnJvcjogYW55O1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIGVycm9yID0gSlNPTi5wYXJzZShkYXRhKTtcclxuICAgICAgICB9IGNhdGNoIChlcnI6IGFueSkge1xyXG4gICAgICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJjb3VsZCBub3QgcGFyc2UgZXJyb3I6IFwiICsgZXJyLm1lc3NhZ2UsIHsgY2F1c2U6IGVyciB9KSk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGxldCBvcHRpb25zOiBFcnJvck9wdGlvbnMgPSB7fTtcclxuICAgICAgICBpZiAoZXJyb3IuY2F1c2UpIHtcclxuICAgICAgICAgICAgb3B0aW9ucy5jYXVzZSA9IGVycm9yLmNhdXNlO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgbGV0IGV4Y2VwdGlvbjtcclxuICAgICAgICBzd2l0Y2ggKGVycm9yLmtpbmQpIHtcclxuICAgICAgICAgICAgY2FzZSBcIlJlZmVyZW5jZUVycm9yXCI6XHJcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgUmVmZXJlbmNlRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuICAgICAgICAgICAgY2FzZSBcIlR5cGVFcnJvclwiOlxyXG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFR5cGVFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgICAgICBjYXNlIFwiUnVudGltZUVycm9yXCI6XHJcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgUnVudGltZUVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xyXG4gICAgICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgICAgIGRlZmF1bHQ6XHJcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIHJlc29sdmVycy5yZWplY3QoZXhjZXB0aW9uKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyBhbmQgcmVtb3ZlcyB0aGUgcmVzcG9uc2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCBmcm9tIHRoZSBjYWxsUmVzcG9uc2VzIG1hcC5cclxuICpcclxuICogQHBhcmFtIGlkIC0gVGhlIElEIG9mIHRoZSByZXNwb25zZSB0byBiZSByZXRyaWV2ZWQgYW5kIHJlbW92ZWQuXHJcbiAqIEByZXR1cm5zIFRoZSByZXNwb25zZSBvYmplY3QgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCwgaWYgYW55LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQ6IHN0cmluZyk6IFByb21pc2VSZXNvbHZlcnMgfCB1bmRlZmluZWQge1xyXG4gICAgY29uc3QgcmVzcG9uc2UgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICByZXR1cm4gcmVzcG9uc2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cclxuICovXHJcbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKTogc3RyaW5nIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsIGEgYm91bmQgbWV0aG9kIGFjY29yZGluZyB0byB0aGUgZ2l2ZW4gY2FsbCBvcHRpb25zLlxyXG4gKlxyXG4gKiBJbiBjYXNlIG9mIGZhaWx1cmUsIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYW4gZXhjZXB0aW9uXHJcbiAqIGFtb25nIFJlZmVyZW5jZUVycm9yICh1bmtub3duIG1ldGhvZCksIFR5cGVFcnJvciAod3JvbmcgYXJndW1lbnQgY291bnQgb3IgdHlwZSksXHJcbiAqIHtAbGluayBSdW50aW1lRXJyb3J9IChtZXRob2QgcmV0dXJuZWQgYW4gZXJyb3IpLCBvciBvdGhlciAobmV0d29yayBvciBpbnRlcm5hbCBlcnJvcnMpLlxyXG4gKiBUaGUgZXhjZXB0aW9uIG1pZ2h0IGhhdmUgYSBcImNhdXNlXCIgZmllbGQgd2l0aCB0aGUgdmFsdWUgcmV0dXJuZWRcclxuICogYnkgdGhlIGFwcGxpY2F0aW9uLSBvciBzZXJ2aWNlLWxldmVsIGVycm9yIG1hcnNoYWxpbmcgZnVuY3Rpb25zLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIEEgbWV0aG9kIGNhbGwgZGVzY3JpcHRvci5cclxuICogQHJldHVybnMgVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnM6IENhbGxPcHRpb25zKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XHJcblxyXG4gICAgY29uc3QgcmVzdWx0ID0gQ2FuY2VsbGFibGVQcm9taXNlLndpdGhSZXNvbHZlcnM8YW55PigpO1xyXG4gICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZTogcmVzdWx0LnJlc29sdmUsIHJlamVjdDogcmVzdWx0LnJlamVjdCB9KTtcclxuXHJcbiAgICBjb25zdCByZXF1ZXN0ID0gY2FsbChDYWxsQmluZGluZywgT2JqZWN0LmFzc2lnbih7IFwiY2FsbC1pZFwiOiBpZCB9LCBvcHRpb25zKSk7XHJcbiAgICBsZXQgcnVubmluZyA9IGZhbHNlO1xyXG5cclxuICAgIHJlcXVlc3QudGhlbigoKSA9PiB7XHJcbiAgICAgICAgcnVubmluZyA9IHRydWU7XHJcbiAgICB9LCAoZXJyKSA9PiB7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIHJlc3VsdC5yZWplY3QoZXJyKTtcclxuICAgIH0pO1xyXG5cclxuICAgIGNvbnN0IGNhbmNlbCA9ICgpID0+IHtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgcmV0dXJuIGNhbmNlbENhbGwoQ2FuY2VsTWV0aG9kLCB7XCJjYWxsLWlkXCI6IGlkfSkuY2F0Y2goKGVycikgPT4ge1xyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igd2hpbGUgcmVxdWVzdGluZyBiaW5kaW5nIGNhbGwgY2FuY2VsbGF0aW9uOlwiLCBlcnIpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfTtcclxuXHJcbiAgICByZXN1bHQub25jYW5jZWxsZWQgPSAoKSA9PiB7XHJcbiAgICAgICAgaWYgKHJ1bm5pbmcpIHtcclxuICAgICAgICAgICAgcmV0dXJuIGNhbmNlbCgpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHJldHVybiByZXF1ZXN0LnRoZW4oY2FuY2VsKTtcclxuICAgICAgICB9XHJcbiAgICB9O1xyXG5cclxuICAgIHJldHVybiByZXN1bHQucHJvbWlzZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgYm91bmQgbWV0aG9kIGJ5IG5hbWUgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cclxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cclxuICpcclxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXHJcbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShtZXRob2ROYW1lOiBzdHJpbmcsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2ROYW1lLCBhcmdzIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIG51bWVyaWMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cclxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cclxuICpcclxuICogQHBhcmFtIG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cclxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cclxuICogQHJldHVybiBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeUlEKG1ldGhvZElEOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2RJRCwgYXJncyB9KTtcclxufVxyXG4iLCAiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vaW5zcGVjdC1qcy9pcy1jYWxsYWJsZVxyXG5cclxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXHJcbi8vXHJcbi8vIENvcHlyaWdodCAoYykgMjAxNSBKb3JkYW4gSGFyYmFuZFxyXG4vL1xyXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5XHJcbi8vIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWxcclxuLy8gaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0c1xyXG4vLyB0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsXHJcbi8vIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpc1xyXG4vLyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxyXG4vL1xyXG4vLyBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcclxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cclxuLy9cclxuLy8gVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxyXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSxcclxuLy8gRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFXHJcbi8vIEFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVJcclxuLy8gTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSxcclxuLy8gT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEVcclxuLy8gU09GVFdBUkUuXHJcblxyXG52YXIgZm5Ub1N0ciA9IEZ1bmN0aW9uLnByb3RvdHlwZS50b1N0cmluZztcclxudmFyIHJlZmxlY3RBcHBseTogdHlwZW9mIFJlZmxlY3QuYXBwbHkgfCBmYWxzZSB8IG51bGwgPSB0eXBlb2YgUmVmbGVjdCA9PT0gJ29iamVjdCcgJiYgUmVmbGVjdCAhPT0gbnVsbCAmJiBSZWZsZWN0LmFwcGx5O1xyXG52YXIgYmFkQXJyYXlMaWtlOiBhbnk7XHJcbnZhciBpc0NhbGxhYmxlTWFya2VyOiBhbnk7XHJcbmlmICh0eXBlb2YgcmVmbGVjdEFwcGx5ID09PSAnZnVuY3Rpb24nICYmIHR5cGVvZiBPYmplY3QuZGVmaW5lUHJvcGVydHkgPT09ICdmdW5jdGlvbicpIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgYmFkQXJyYXlMaWtlID0gT2JqZWN0LmRlZmluZVByb3BlcnR5KHt9LCAnbGVuZ3RoJywge1xyXG4gICAgICAgICAgICBnZXQ6IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgICAgIHRocm93IGlzQ2FsbGFibGVNYXJrZXI7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgICAgICBpc0NhbGxhYmxlTWFya2VyID0ge307XHJcbiAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXRocm93LWxpdGVyYWxcclxuICAgICAgICByZWZsZWN0QXBwbHkoZnVuY3Rpb24gKCkgeyB0aHJvdyA0MjsgfSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcclxuICAgIH0gY2F0Y2ggKF8pIHtcclxuICAgICAgICBpZiAoXyAhPT0gaXNDYWxsYWJsZU1hcmtlcikge1xyXG4gICAgICAgICAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufSBlbHNlIHtcclxuICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XHJcbn1cclxuXHJcbnZhciBjb25zdHJ1Y3RvclJlZ2V4ID0gL15cXHMqY2xhc3NcXGIvO1xyXG52YXIgaXNFUzZDbGFzc0ZuID0gZnVuY3Rpb24gaXNFUzZDbGFzc0Z1bmN0aW9uKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgdmFyIGZuU3RyID0gZm5Ub1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICByZXR1cm4gY29uc3RydWN0b3JSZWdleC50ZXN0KGZuU3RyKTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7IC8vIG5vdCBhIGZ1bmN0aW9uXHJcbiAgICB9XHJcbn07XHJcblxyXG52YXIgdHJ5RnVuY3Rpb25PYmplY3QgPSBmdW5jdGlvbiB0cnlGdW5jdGlvblRvU3RyKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XHJcbiAgICAgICAgZm5Ub1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICByZXR1cm4gdHJ1ZTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcbn07XHJcbnZhciB0b1N0ciA9IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmc7XHJcbnZhciBvYmplY3RDbGFzcyA9ICdbb2JqZWN0IE9iamVjdF0nO1xyXG52YXIgZm5DbGFzcyA9ICdbb2JqZWN0IEZ1bmN0aW9uXSc7XHJcbnZhciBnZW5DbGFzcyA9ICdbb2JqZWN0IEdlbmVyYXRvckZ1bmN0aW9uXSc7XHJcbnZhciBkZGFDbGFzcyA9ICdbb2JqZWN0IEhUTUxBbGxDb2xsZWN0aW9uXSc7IC8vIElFIDExXHJcbnZhciBkZGFDbGFzczIgPSAnW29iamVjdCBIVE1MIGRvY3VtZW50LmFsbCBjbGFzc10nO1xyXG52YXIgZGRhQ2xhc3MzID0gJ1tvYmplY3QgSFRNTENvbGxlY3Rpb25dJzsgLy8gSUUgOS0xMFxyXG52YXIgaGFzVG9TdHJpbmdUYWcgPSB0eXBlb2YgU3ltYm9sID09PSAnZnVuY3Rpb24nICYmICEhU3ltYm9sLnRvU3RyaW5nVGFnOyAvLyBiZXR0ZXI6IHVzZSBgaGFzLXRvc3RyaW5ndGFnYFxyXG5cclxudmFyIGlzSUU2OCA9ICEoMCBpbiBbLF0pOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIG5vLXNwYXJzZS1hcnJheXMsIGNvbW1hLXNwYWNpbmdcclxuXHJcbnZhciBpc0REQTogKHZhbHVlOiBhbnkpID0+IGJvb2xlYW4gPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKCkgeyByZXR1cm4gZmFsc2U7IH07XHJcbmlmICh0eXBlb2YgZG9jdW1lbnQgPT09ICdvYmplY3QnKSB7XHJcbiAgICAvLyBGaXJlZm94IDMgY2Fub25pY2FsaXplcyBEREEgdG8gdW5kZWZpbmVkIHdoZW4gaXQncyBub3QgYWNjZXNzZWQgZGlyZWN0bHlcclxuICAgIHZhciBhbGwgPSBkb2N1bWVudC5hbGw7XHJcbiAgICBpZiAodG9TdHIuY2FsbChhbGwpID09PSB0b1N0ci5jYWxsKGRvY3VtZW50LmFsbCkpIHtcclxuICAgICAgICBpc0REQSA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwodmFsdWUpIHtcclxuICAgICAgICAgICAgLyogZ2xvYmFscyBkb2N1bWVudDogZmFsc2UgKi9cclxuICAgICAgICAgICAgLy8gaW4gSUUgNi04LCB0eXBlb2YgZG9jdW1lbnQuYWxsIGlzIFwib2JqZWN0XCIgYW5kIGl0J3MgdHJ1dGh5XHJcbiAgICAgICAgICAgIGlmICgoaXNJRTY4IHx8ICF2YWx1ZSkgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3VuZGVmaW5lZCcgfHwgdHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JykpIHtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgdmFyIHN0ciA9IHRvU3RyLmNhbGwodmFsdWUpO1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybiAoXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHN0ciA9PT0gZGRhQ2xhc3NcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczJcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczMgLy8gb3BlcmEgMTIuMTZcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBvYmplY3RDbGFzcyAvLyBJRSA2LThcclxuICAgICAgICAgICAgICAgICAgICApICYmIHZhbHVlKCcnKSA9PSBudWxsOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIGVxZXFlcVxyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZSkgeyAvKiovIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gaXNDYWxsYWJsZVJlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55ICB7XHJcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XHJcbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgdHJ5IHtcclxuICAgICAgICAocmVmbGVjdEFwcGx5IGFzIGFueSkodmFsdWUsIG51bGwsIGJhZEFycmF5TGlrZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgaWYgKGUgIT09IGlzQ2FsbGFibGVNYXJrZXIpIHsgcmV0dXJuIGZhbHNlOyB9XHJcbiAgICB9XHJcbiAgICByZXR1cm4gIWlzRVM2Q2xhc3NGbih2YWx1ZSkgJiYgdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBpc0NhbGxhYmxlTm9SZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSB7XHJcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XHJcbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKGhhc1RvU3RyaW5nVGFnKSB7IHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7IH1cclxuICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgdmFyIHN0ckNsYXNzID0gdG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICBpZiAoc3RyQ2xhc3MgIT09IGZuQ2xhc3MgJiYgc3RyQ2xhc3MgIT09IGdlbkNsYXNzICYmICEoL15cXFtvYmplY3QgSFRNTC8pLnRlc3Qoc3RyQ2xhc3MpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcclxufTtcclxuXHJcbmV4cG9ydCBkZWZhdWx0IHJlZmxlY3RBcHBseSA/IGlzQ2FsbGFibGVSZWZBcHBseSA6IGlzQ2FsbGFibGVOb1JlZkFwcGx5O1xyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IGlzQ2FsbGFibGUgZnJvbSBcIi4vY2FsbGFibGUuanNcIjtcclxuXHJcbi8qKlxyXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHVzZWQgYXMgcmVqZWN0aW9uIHJlYXNvblxyXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gaXMgY2FuY2VsbGVkIHN1Y2Nlc3NmdWxseS5cclxuICpcclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxFcnJvclwiYC5cclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSBjYXVzZSBwYXNzZWQgdG8gdGhlIGNhbmNlbCBtZXRob2QsIGlmIGFueS5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBDYW5jZWxFcnJvciBleHRlbmRzIEVycm9yIHtcclxuICAgIC8qKlxyXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsRXJyb3JgIGluc3RhbmNlLlxyXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XHJcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxFcnJvclwiO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uXHJcbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSByZWplY3RzIGFmdGVyIGJlaW5nIGNhbmNlbGxlZCxcclxuICogb3Igd2hlbiB0aGUgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB0aHJvd3Mgb3IgcmVqZWN0cy5cclxuICpcclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiYC5cclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSByZWFzb24gdGhlIHByb21pc2UgcmVqZWN0ZWQgd2l0aC5cclxuICpcclxuICogQmVjYXVzZSB0aGUgb3JpZ2luYWwgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLFxyXG4gKiBhIHdyYXBwZXIgcHJvbWlzZSB3aWxsIGJlIHBhc3NlZCB0byB0aGUgdW5oYW5kbGVkIHJlamVjdGlvbiBsaXN0ZW5lciBpbnN0ZWFkLlxyXG4gKiBUaGUge0BsaW5rIHByb21pc2V9IHByb3BlcnR5IGhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBvcmlnaW5hbCBwcm9taXNlLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIENhbmNlbGxlZFJlamVjdGlvbkVycm9yIGV4dGVuZHMgRXJyb3Ige1xyXG4gICAgLyoqXHJcbiAgICAgKiBIb2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgcHJvbWlzZSB0aGF0IHdhcyBjYW5jZWxsZWQgYW5kIHRoZW4gcmVqZWN0ZWQuXHJcbiAgICAgKi9cclxuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPjtcclxuXHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbGxlZFJlamVjdGlvbkVycm9yYCBpbnN0YW5jZS5cclxuICAgICAqIEBwYXJhbSBwcm9taXNlIC0gVGhlIHByb21pc2UgdGhhdCBjYXVzZWQgdGhlIGVycm9yIG9yaWdpbmFsbHkuXHJcbiAgICAgKiBAcGFyYW0gcmVhc29uIC0gVGhlIHJlamVjdGlvbiByZWFzb24uXHJcbiAgICAgKiBAcGFyYW0gaW5mbyAtIEFuIG9wdGlvbmFsIGluZm9ybWF0aXZlIG1lc3NhZ2Ugc3BlY2lmeWluZyB0aGUgY2lyY3Vtc3RhbmNlcyBpbiB3aGljaCB0aGUgZXJyb3Igd2FzIHRocm93bi5cclxuICAgICAqICAgICAgICAgICAgICAgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiYC5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCByZWFzb24/OiBhbnksIGluZm8/OiBzdHJpbmcpIHtcclxuICAgICAgICBzdXBlcigoaW5mbyA/PyBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCIpICsgXCIgUmVhc29uOiBcIiArIGVycm9yTWVzc2FnZShyZWFzb24pLCB7IGNhdXNlOiByZWFzb24gfSk7XHJcbiAgICAgICAgdGhpcy5wcm9taXNlID0gcHJvbWlzZTtcclxuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCI7XHJcbiAgICB9XHJcbn1cclxuXHJcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4gPSAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XHJcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IgPSAocmVhc29uPzogYW55KSA9PiB2b2lkO1xyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciA9IChjYXVzZT86IGFueSkgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+ID0gKHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+LCByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yKSA9PiB2b2lkO1xyXG5cclxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcclxuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCk6IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj47XHJcbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBXcmFwcyBhIGNhbmNlbGxhYmxlIHByb21pc2UgYWxvbmcgd2l0aCBpdHMgcmVzb2x1dGlvbiBtZXRob2RzLlxyXG4gKiBUaGUgYG9uY2FuY2VsbGVkYCBmaWVsZCB3aWxsIGJlIG51bGwgaW5pdGlhbGx5IGJ1dCBtYXkgYmUgc2V0IHRvIHByb3ZpZGUgYSBjdXN0b20gY2FuY2VsbGF0aW9uIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcclxuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcclxuICAgIHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+O1xyXG4gICAgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcjtcclxuICAgIG9uY2FuY2VsbGVkOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgfCBudWxsO1xyXG59XHJcblxyXG5pbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUge1xyXG4gICAgcmVhZG9ubHkgcm9vdDogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGU7XHJcbiAgICByZXNvbHZpbmc6IGJvb2xlYW47XHJcbiAgICBzZXR0bGVkOiBib29sZWFuO1xyXG4gICAgcmVhc29uPzogQ2FuY2VsRXJyb3I7XHJcbn1cclxuXHJcbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXHJcbmNvbnN0IGJhcnJpZXJTeW0gPSBTeW1ib2woXCJiYXJyaWVyXCIpO1xyXG5jb25zdCBjYW5jZWxJbXBsU3ltID0gU3ltYm9sKFwiY2FuY2VsSW1wbFwiKTtcclxuY29uc3Qgc3BlY2llcyA9IFN5bWJvbC5zcGVjaWVzID8/IFN5bWJvbChcInNwZWNpZXNQb2x5ZmlsbFwiKTtcclxuXHJcbi8qKlxyXG4gKiBBIHByb21pc2Ugd2l0aCBhbiBhdHRhY2hlZCBtZXRob2QgZm9yIGNhbmNlbGxpbmcgbG9uZy1ydW5uaW5nIG9wZXJhdGlvbnMgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbH0pLlxyXG4gKiBDYW5jZWxsYXRpb24gY2FuIG9wdGlvbmFsbHkgYmUgYm91bmQgdG8gYW4ge0BsaW5rIEFib3J0U2lnbmFsfVxyXG4gKiBmb3IgYmV0dGVyIGNvbXBvc2FiaWxpdHkgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbE9ufSkuXHJcbiAqXHJcbiAqIENhbmNlbGxpbmcgYSBwZW5kaW5nIHByb21pc2Ugd2lsbCByZXN1bHQgaW4gYW4gaW1tZWRpYXRlIHJlamVjdGlvblxyXG4gKiB3aXRoIGFuIGluc3RhbmNlIG9mIHtAbGluayBDYW5jZWxFcnJvcn0gYXMgcmVhc29uLFxyXG4gKiBidXQgd2hvZXZlciBzdGFydGVkIHRoZSBwcm9taXNlIHdpbGwgYmUgcmVzcG9uc2libGVcclxuICogZm9yIGFjdHVhbGx5IGFib3J0aW5nIHRoZSB1bmRlcmx5aW5nIG9wZXJhdGlvbi5cclxuICogVG8gdGhpcyBwdXJwb3NlLCB0aGUgY29uc3RydWN0b3IgYW5kIGFsbCBjaGFpbmluZyBtZXRob2RzXHJcbiAqIGFjY2VwdCBvcHRpb25hbCBjYW5jZWxsYXRpb24gY2FsbGJhY2tzLlxyXG4gKlxyXG4gKiBJZiBhIGBDYW5jZWxsYWJsZVByb21pc2VgIHN0aWxsIHJlc29sdmVzIGFmdGVyIGhhdmluZyBiZWVuIGNhbmNlbGxlZCxcclxuICogdGhlIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZC4gSWYgaXQgcmVqZWN0cywgdGhlIHJlYXNvblxyXG4gKiB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24sXHJcbiAqIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGluc3RhbmNlLlxyXG4gKiBUbyBmYWNpbGl0YXRlIHRoZSBoYW5kbGluZyBvZiBjYW5jZWxsYXRpb24gcmVxdWVzdHMsXHJcbiAqIGNhbmNlbGxlZCBgQ2FuY2VsbGFibGVQcm9taXNlYHMgd2lsbCBfbm90XyByZXBvcnQgdW5oYW5kbGVkIGBDYW5jZWxFcnJvcmBzXHJcbiAqIHdob3NlIGBjYXVzZWAgZmllbGQgaXMgdGhlIHNhbWUgYXMgdGhlIG9uZSB3aXRoIHdoaWNoIHRoZSBjdXJyZW50IHByb21pc2Ugd2FzIGNhbmNlbGxlZC5cclxuICpcclxuICogQWxsIHVzdWFsIHByb21pc2UgbWV0aG9kcyBhcmUgZGVmaW5lZCBhbmQgcmV0dXJuIGEgYENhbmNlbGxhYmxlUHJvbWlzZWBcclxuICogd2hvc2UgY2FuY2VsIG1ldGhvZCB3aWxsIGNhbmNlbCB0aGUgcGFyZW50IG9wZXJhdGlvbiBhcyB3ZWxsLCBwcm9wYWdhdGluZyB0aGUgY2FuY2VsbGF0aW9uIHJlYXNvblxyXG4gKiB1cHdhcmRzIHRocm91Z2ggcHJvbWlzZSBjaGFpbnMuXHJcbiAqIENvbnZlcnNlbHksIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpbGwgbm90IGF1dG9tYXRpY2FsbHkgY2FuY2VsIGRlcGVuZGVudCBwcm9taXNlcyBkb3duc3RyZWFtOlxyXG4gKiBgYGB0c1xyXG4gKiBsZXQgcm9vdCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XHJcbiAqIGxldCBjaGlsZDEgPSByb290LnRoZW4oKCkgPT4geyAuLi4gfSk7XHJcbiAqIGxldCBjaGlsZDIgPSBjaGlsZDEudGhlbigoKSA9PiB7IC4uLiB9KTtcclxuICogbGV0IGNoaWxkMyA9IHJvb3QuY2F0Y2goKCkgPT4geyAuLi4gfSk7XHJcbiAqIGNoaWxkMS5jYW5jZWwoKTsgLy8gQ2FuY2VscyBjaGlsZDEgYW5kIHJvb3QsIGJ1dCBub3QgY2hpbGQyIG9yIGNoaWxkM1xyXG4gKiBgYGBcclxuICogQ2FuY2VsbGluZyBhIHByb21pc2UgdGhhdCBoYXMgYWxyZWFkeSBzZXR0bGVkIGlzIHNhZmUgYW5kIGhhcyBubyBjb25zZXF1ZW5jZS5cclxuICpcclxuICogVGhlIGBjYW5jZWxgIG1ldGhvZCByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IF9hbHdheXMgZnVsZmlsbHNfXHJcbiAqIGFmdGVyIHRoZSB3aG9sZSBjaGFpbiBoYXMgcHJvY2Vzc2VkIHRoZSBjYW5jZWwgcmVxdWVzdFxyXG4gKiBhbmQgYWxsIGF0dGFjaGVkIGNhbGxiYWNrcyB1cCB0byB0aGF0IG1vbWVudCBoYXZlIHJ1bi5cclxuICpcclxuICogQWxsIEVTMjAyNCBwcm9taXNlIG1ldGhvZHMgKHN0YXRpYyBhbmQgaW5zdGFuY2UpIGFyZSBkZWZpbmVkIG9uIENhbmNlbGxhYmxlUHJvbWlzZSxcclxuICogYnV0IGFjdHVhbCBhdmFpbGFiaWxpdHkgbWF5IHZhcnkgd2l0aCBPUy93ZWJ2aWV3IHZlcnNpb24uXHJcbiAqXHJcbiAqIEluIGxpbmUgd2l0aCB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZyxcclxuICogYENhbmNlbGxhYmxlUHJvbWlzZWAgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy5cclxuICogRXh0ZW5kZXJzIHNob3VsZCB0YWtlIGNhcmUgdG8gcHJvdmlkZSB0aGVpciBvd24gbWV0aG9kIGltcGxlbWVudGF0aW9ucy5cclxuICogVGhpcyBtaWdodCBiZSByZWNvbnNpZGVyZWQgaW4gY2FzZSB0aGUgcHJvcG9zYWwgaXMgcmV0aXJlZC5cclxuICpcclxuICogQ2FuY2VsbGFibGVQcm9taXNlIGlzIGEgd3JhcHBlciBhcm91bmQgdGhlIERPTSBQcm9taXNlIG9iamVjdFxyXG4gKiBhbmQgaXMgY29tcGxpYW50IHdpdGggdGhlIFtQcm9taXNlcy9BKyBzcGVjaWZpY2F0aW9uXShodHRwczovL3Byb21pc2VzYXBsdXMuY29tLylcclxuICogKGl0IHBhc3NlcyB0aGUgW2NvbXBsaWFuY2Ugc3VpdGVdKGh0dHBzOi8vZ2l0aHViLmNvbS9wcm9taXNlcy1hcGx1cy9wcm9taXNlcy10ZXN0cykpXHJcbiAqIGlmIHNvIGlzIHRoZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIENhbmNlbGxhYmxlUHJvbWlzZTxUPiBleHRlbmRzIFByb21pc2U8VD4gaW1wbGVtZW50cyBQcm9taXNlTGlrZTxUPiwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XHJcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cclxuICAgIC8qKiBAaW50ZXJuYWwgKi9cclxuICAgIHByaXZhdGUgW2JhcnJpZXJTeW1dITogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gfCBudWxsO1xyXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xyXG4gICAgcHJpdmF0ZSByZWFkb25seSBbY2FuY2VsSW1wbFN5bV0hOiAocmVhc29uOiBDYW5jZWxFcnJvcikgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBgQ2FuY2VsbGFibGVQcm9taXNlYC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZXhlY3V0b3IgLSBBIGNhbGxiYWNrIHVzZWQgdG8gaW5pdGlhbGl6ZSB0aGUgcHJvbWlzZS4gVGhpcyBjYWxsYmFjayBpcyBwYXNzZWQgdHdvIGFyZ3VtZW50czpcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGEgYHJlc29sdmVgIGNhbGxiYWNrIHVzZWQgdG8gcmVzb2x2ZSB0aGUgcHJvbWlzZSB3aXRoIGEgdmFsdWVcclxuICAgICAqICAgICAgICAgICAgICAgICAgIG9yIHRoZSByZXN1bHQgb2YgYW5vdGhlciBwcm9taXNlIChwb3NzaWJseSBjYW5jZWxsYWJsZSksXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYSBgcmVqZWN0YCBjYWxsYmFjayB1c2VkIHRvIHJlamVjdCB0aGUgcHJvbWlzZSB3aXRoIGEgcHJvdmlkZWQgcmVhc29uIG9yIGVycm9yLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgdGhlIHZhbHVlIHByb3ZpZGVkIHRvIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqICAgICAgICAgICAgICAgICAgIChpdCBoYXMgYSBgdGhlbmAgX2FuZF8gYSBgY2FuY2VsYCBtZXRob2QpLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZm9yd2FyZGVkIHRvIHRoYXQgb2JqZWN0IGFuZCB0aGUgb25jYW5jZWxsZWQgd2lsbCBub3QgYmUgaW52b2tlZCBhbnltb3JlLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgYW55IG9uZSBvZiB0aGUgdHdvIGNhbGxiYWNrcyBpcyBjYWxsZWQgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGUgcHJvdmlkZWQgdmFsdWVzIHdpbGwgYmUgY2FuY2VsbGVkIGFuZCByZXNvbHZlZCBhcyB1c3VhbCxcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGJ1dCB0aGVpciByZXN1bHRzIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgSG93ZXZlciwgaWYgdGhlIHJlc29sdXRpb24gcHJvY2VzcyB1bHRpbWF0ZWx5IGVuZHMgdXAgaW4gYSByZWplY3Rpb25cclxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoYXQgaXMgbm90IGR1ZSB0byBjYW5jZWxsYXRpb24sIHRoZSByZWplY3Rpb24gcmVhc29uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICB3aWxsIGJlIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICogQHBhcmFtIG9uY2FuY2VsbGVkIC0gSXQgaXMgdGhlIGNhbGxlcidzIHJlc3BvbnNpYmlsaXR5IHRvIGVuc3VyZSB0aGF0IGFueSBvcGVyYXRpb25cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHN0YXJ0ZWQgYnkgdGhlIGV4ZWN1dG9yIGlzIHByb3Blcmx5IGhhbHRlZCB1cG9uIGNhbmNlbGxhdGlvbi5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFRoaXMgb3B0aW9uYWwgY2FsbGJhY2sgY2FuIGJlIHVzZWQgdG8gdGhhdCBwdXJwb3NlLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSXQgd2lsbCBiZSBjYWxsZWQgX3N5bmNocm9ub3VzbHlfIHdpdGggYSBjYW5jZWxsYXRpb24gY2F1c2VcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdoZW4gY2FuY2VsbGF0aW9uIGlzIHJlcXVlc3RlZCwgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSByZWplY3RlZFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0sIGJ1dCBfYmVmb3JlX1xyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW55IHtAbGluayB0aGVufS97QGxpbmsgY2F0Y2h9L3tAbGluayBmaW5hbGx5fSBjYWxsYmFjayBydW5zLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfVxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2lsbCBvbmx5IGZ1bGZpbGwgYWZ0ZXIgdGhlIGZvcm1lciBoYXMgc2V0dGxlZC5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFVuaGFuZGxlZCBleGNlcHRpb25zIG9yIHJlamVjdGlvbnMgZnJvbSB0aGUgY2FsbGJhY2sgd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgdW5oYW5kbGVkIHJlamVjdGlvbnMuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGNhbGxlZCBiZWZvcmUgY2FuY2VsbGF0aW9uIHdpdGggYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoaXMgcHJvbWlzZSB3aWxsIGJlIGRpdmVydGVkIHRvIHRoYXQgcHJvbWlzZSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFuZCB0aGUgb3JpZ2luYWwgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoZXhlY3V0b3I6IENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+LCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcikge1xyXG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XHJcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XHJcbiAgICAgICAgc3VwZXIoKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XHJcblxyXG4gICAgICAgIGlmICgodGhpcy5jb25zdHJ1Y3RvciBhcyBhbnkpW3NwZWNpZXNdICE9PSBQcm9taXNlKSB7XHJcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy4gUGxlYXNlIHJlZnJhaW4gZnJvbSBvdmVycmlkaW5nIHRoZSBbU3ltYm9sLnNwZWNpZXNdIHN0YXRpYyBwcm9wZXJ0eS5cIik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBsZXQgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHtcclxuICAgICAgICAgICAgcHJvbWlzZTogdGhpcyxcclxuICAgICAgICAgICAgcmVzb2x2ZSxcclxuICAgICAgICAgICAgcmVqZWN0LFxyXG4gICAgICAgICAgICBnZXQgb25jYW5jZWxsZWQoKSB7IHJldHVybiBvbmNhbmNlbGxlZCA/PyBudWxsOyB9LFxyXG4gICAgICAgICAgICBzZXQgb25jYW5jZWxsZWQoY2IpIHsgb25jYW5jZWxsZWQgPSBjYiA/PyB1bmRlZmluZWQ7IH1cclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICBjb25zdCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XHJcbiAgICAgICAgICAgIGdldCByb290KCkgeyByZXR1cm4gc3RhdGU7IH0sXHJcbiAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXHJcbiAgICAgICAgICAgIHNldHRsZWQ6IGZhbHNlXHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgLy8gU2V0dXAgY2FuY2VsbGF0aW9uIHN5c3RlbS5cclxuICAgICAgICB2b2lkIE9iamVjdC5kZWZpbmVQcm9wZXJ0aWVzKHRoaXMsIHtcclxuICAgICAgICAgICAgW2JhcnJpZXJTeW1dOiB7XHJcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogdHJ1ZSxcclxuICAgICAgICAgICAgICAgIHZhbHVlOiBudWxsXHJcbiAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgIFtjYW5jZWxJbXBsU3ltXToge1xyXG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgdmFsdWU6IGNhbmNlbGxlckZvcihwcm9taXNlLCBzdGF0ZSlcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pO1xyXG5cclxuICAgICAgICAvLyBSdW4gdGhlIGFjdHVhbCBleGVjdXRvci5cclxuICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIHN0YXRlKTtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBleGVjdXRvcihyZXNvbHZlckZvcihwcm9taXNlLCBzdGF0ZSksIHJlamVjdG9yKTtcclxuICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykge1xyXG4gICAgICAgICAgICAgICAgY29uc29sZS5sb2coXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIENhbmNlbGxhYmxlUHJvbWlzZSBleGVjdXRvci5cIiwgZXJyKTtcclxuICAgICAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDYW5jZWxzIGltbWVkaWF0ZWx5IHRoZSBleGVjdXRpb24gb2YgdGhlIG9wZXJhdGlvbiBhc3NvY2lhdGVkIHdpdGggdGhpcyBwcm9taXNlLlxyXG4gICAgICogVGhlIHByb21pc2UgcmVqZWN0cyB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSBpbnN0YW5jZSBhcyByZWFzb24sXHJcbiAgICAgKiB3aXRoIHRoZSB7QGxpbmsgQ2FuY2VsRXJyb3IjY2F1c2V9IHByb3BlcnR5IHNldCB0byB0aGUgZ2l2ZW4gYXJndW1lbnQsIGlmIGFueS5cclxuICAgICAqXHJcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCBhZnRlciB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkO1xyXG4gICAgICogcmVwZWF0ZWQgY2FsbHMgaW4gcGFydGljdWxhciBhcmUgc2FmZSwgYnV0IG9ubHkgdGhlIGZpcnN0IG9uZVxyXG4gICAgICogd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgYENhbmNlbEVycm9yYCBleGNlcHRpb24gX25lZWQgbm90XyBiZSBoYW5kbGVkIGV4cGxpY2l0bHkgX29uIHRoZSBwcm9taXNlcyB0aGF0IGFyZSBiZWluZyBjYW5jZWxsZWQ6X1xyXG4gICAgICogY2FuY2VsbGluZyBhIHByb21pc2Ugd2l0aCBubyBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlciBkb2VzIG5vdCB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQuXHJcbiAgICAgKiBUaGVyZWZvcmUsIHRoZSBmb2xsb3dpbmcgaWRpb21zIGFyZSBhbGwgZXF1YWxseSBjb3JyZWN0OlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkuY2FuY2VsKCk7XHJcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYW5jZWwoKTtcclxuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhdGNoKC4uLikuY2FuY2VsKCk7XHJcbiAgICAgKiBgYGBcclxuICAgICAqIFdoZW5ldmVyIHNvbWUgY2FuY2VsbGVkIHByb21pc2UgaW4gYSBjaGFpbiByZWplY3RzIHdpdGggYSBgQ2FuY2VsRXJyb3JgXHJcbiAgICAgKiB3aXRoIHRoZSBzYW1lIGNhbmNlbGxhdGlvbiBjYXVzZSBhcyBpdHNlbGYsIHRoZSBlcnJvciB3aWxsIGJlIGRpc2NhcmRlZCBzaWxlbnRseS5cclxuICAgICAqIEhvd2V2ZXIsIHRoZSBgQ2FuY2VsRXJyb3JgIF93aWxsIHN0aWxsIGJlIGRlbGl2ZXJlZF8gdG8gYWxsIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyc1xyXG4gICAgICogYWRkZWQgYnkge0BsaW5rIHRoZW59IGFuZCByZWxhdGVkIG1ldGhvZHM6XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcclxuICAgICAqIGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkuY2F0Y2goY29uc29sZS5sb2cpO1xyXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIEEgQ2FuY2VsRXJyb3IgaXMgcHJpbnRlZCB0byB0aGUgY29uc29sZS5cclxuICAgICAqIGBgYFxyXG4gICAgICogSWYgdGhlIGBDYW5jZWxFcnJvcmAgaXMgbm90IGhhbmRsZWQgZG93bnN0cmVhbSBieSB0aGUgdGltZSBpdCByZWFjaGVzXHJcbiAgICAgKiBhIF9ub24tY2FuY2VsbGVkXyBwcm9taXNlLCBpdCBfd2lsbF8gdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LFxyXG4gICAgICoganVzdCBsaWtlIG5vcm1hbCByZWplY3Rpb25zIHdvdWxkOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XHJcbiAgICAgKiBsZXQgY2hhaW5lZCA9IGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkudGhlbigoKSA9PiB7IC4uLiB9KTsgLy8gTm8gY2F0Y2guLi5cclxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBVbmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50IG9uIGNoYWluZWQhXHJcbiAgICAgKiBgYGBcclxuICAgICAqIFRoZXJlZm9yZSwgaXQgaXMgaW1wb3J0YW50IHRvIGVpdGhlciBjYW5jZWwgd2hvbGUgcHJvbWlzZSBjaGFpbnMgZnJvbSB0aGVpciB0YWlsLFxyXG4gICAgICogYXMgc2hvd24gaW4gdGhlIGNvcnJlY3QgaWRpb21zIGFib3ZlLCBvciB0YWtlIGNhcmUgb2YgaGFuZGxpbmcgZXJyb3JzIGV2ZXJ5d2hlcmUuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgQSBjYW5jZWxsYWJsZSBwcm9taXNlIHRoYXQgX2Z1bGZpbGxzXyBhZnRlciB0aGUgY2FuY2VsIGNhbGxiYWNrIChpZiBhbnkpXHJcbiAgICAgKiBhbmQgYWxsIGhhbmRsZXJzIGF0dGFjaGVkIHVwIHRvIHRoZSBjYWxsIHRvIGNhbmNlbCBoYXZlIHJ1bi5cclxuICAgICAqIElmIHRoZSBjYW5jZWwgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBieSBgY2FuY2VsYFxyXG4gICAgICogd2lsbCBhbHNvIHdhaXQgZm9yIHRoYXQgdGhlbmFibGUgdG8gc2V0dGxlLlxyXG4gICAgICogVGhpcyBlbmFibGVzIGNhbGxlcnMgdG8gd2FpdCBmb3IgdGhlIGNhbmNlbGxlZCBvcGVyYXRpb24gdG8gdGVybWluYXRlXHJcbiAgICAgKiB3aXRob3V0IGJlaW5nIGZvcmNlZCB0byBoYW5kbGUgcG90ZW50aWFsIGVycm9ycyBhdCB0aGUgY2FsbCBzaXRlLlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpLnRoZW4oKCkgPT4ge1xyXG4gICAgICogICAgIC8vIENsZWFudXAgZmluaXNoZWQsIGl0J3Mgc2FmZSB0byBkbyBzb21ldGhpbmcgZWxzZS5cclxuICAgICAqIH0sIChlcnIpID0+IHtcclxuICAgICAqICAgICAvLyBVbnJlYWNoYWJsZTogdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSBjYW5jZWwgd2lsbCBuZXZlciByZWplY3QuXHJcbiAgICAgKiB9KTtcclxuICAgICAqIGBgYFxyXG4gICAgICogTm90ZSB0aGF0IHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgX25vdF8gaGFuZGxlIGltcGxpY2l0bHkgYW55IHJlamVjdGlvblxyXG4gICAgICogdGhhdCBtaWdodCBoYXZlIG9jY3VycmVkIGFscmVhZHkgaW4gdGhlIGNhbmNlbGxlZCBjaGFpbi5cclxuICAgICAqIEl0IHdpbGwganVzdCB0cmFjayB3aGV0aGVyIHJlZ2lzdGVyZWQgaGFuZGxlcnMgaGF2ZSBiZWVuIGV4ZWN1dGVkIG9yIG5vdC5cclxuICAgICAqIFRoZXJlZm9yZSwgdW5oYW5kbGVkIHJlamVjdGlvbnMgd2lsbCBuZXZlciBiZSBzaWxlbnRseSBoYW5kbGVkIGJ5IGNhbGxpbmcgY2FuY2VsLlxyXG4gICAgICovXHJcbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XHJcbiAgICAgICAgICAgIC8vIElOVkFSSUFOVDogdGhlIHJlc3VsdCBvZiB0aGlzW2NhbmNlbEltcGxTeW1dIGFuZCB0aGUgYmFycmllciBkbyBub3QgZXZlciByZWplY3QuXHJcbiAgICAgICAgICAgIC8vIFVuZm9ydHVuYXRlbHkgbWFjT1MgSGlnaCBTaWVycmEgZG9lcyBub3Qgc3VwcG9ydCBQcm9taXNlLmFsbFNldHRsZWQuXHJcbiAgICAgICAgICAgIFByb21pc2UuYWxsKFtcclxuICAgICAgICAgICAgICAgIHRoaXNbY2FuY2VsSW1wbFN5bV0obmV3IENhbmNlbEVycm9yKFwiUHJvbWlzZSBjYW5jZWxsZWQuXCIsIHsgY2F1c2UgfSkpLFxyXG4gICAgICAgICAgICAgICAgY3VycmVudEJhcnJpZXIodGhpcylcclxuICAgICAgICAgICAgXSkudGhlbigoKSA9PiByZXNvbHZlKCksICgpID0+IHJlc29sdmUoKSk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBCaW5kcyBwcm9taXNlIGNhbmNlbGxhdGlvbiB0byB0aGUgYWJvcnQgZXZlbnQgb2YgdGhlIGdpdmVuIHtAbGluayBBYm9ydFNpZ25hbH0uXHJcbiAgICAgKiBJZiB0aGUgc2lnbmFsIGhhcyBhbHJlYWR5IGFib3J0ZWQsIHRoZSBwcm9taXNlIHdpbGwgYmUgY2FuY2VsbGVkIGltbWVkaWF0ZWx5LlxyXG4gICAgICogV2hlbiBlaXRoZXIgY29uZGl0aW9uIGlzIHZlcmlmaWVkLCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlIHdpbGwgYmUgc2V0XHJcbiAgICAgKiB0byB0aGUgc2lnbmFsJ3MgYWJvcnQgcmVhc29uIChzZWUge0BsaW5rIEFib3J0U2lnbmFsI3JlYXNvbn0pLlxyXG4gICAgICpcclxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIChvciBpZiB0aGUgc2lnbmFsIGFib3J0cykgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkLlxyXG4gICAgICogT25seSB0aGUgZmlyc3Qgc2lnbmFsIHRvIGFib3J0IHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogRm9yIG1vcmUgZGV0YWlscyBhYm91dCB0aGUgY2FuY2VsbGF0aW9uIHByb2Nlc3MsXHJcbiAgICAgKiBzZWUge0BsaW5rIGNhbmNlbH0gYW5kIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYCBjb25zdHJ1Y3Rvci5cclxuICAgICAqXHJcbiAgICAgKiBUaGlzIG1ldGhvZCBlbmFibGVzIGBhd2FpdGBpbmcgY2FuY2VsbGFibGUgcHJvbWlzZXMgd2l0aG91dCBoYXZpbmdcclxuICAgICAqIHRvIHN0b3JlIHRoZW0gZm9yIGZ1dHVyZSBjYW5jZWxsYXRpb24sIGUuZy46XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogYXdhaXQgbG9uZ1J1bm5pbmdPcGVyYXRpb24oKS5jYW5jZWxPbihzaWduYWwpO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBpbnN0ZWFkIG9mOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGxldCBwcm9taXNlVG9CZUNhbmNlbGxlZCA9IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCk7XHJcbiAgICAgKiBhd2FpdCBwcm9taXNlVG9CZUNhbmNlbGxlZDtcclxuICAgICAqIGBgYFxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoaXMgcHJvbWlzZSwgZm9yIG1ldGhvZCBjaGFpbmluZy5cclxuICAgICAqL1xyXG4gICAgY2FuY2VsT24oc2lnbmFsOiBBYm9ydFNpZ25hbCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgaWYgKHNpZ25hbC5hYm9ydGVkKSB7XHJcbiAgICAgICAgICAgIHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbilcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBzaWduYWwuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pLCB7Y2FwdHVyZTogdHJ1ZX0pO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgcmV0dXJuIHRoaXM7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBBdHRhY2hlcyBjYWxsYmFja3MgZm9yIHRoZSByZXNvbHV0aW9uIGFuZC9vciByZWplY3Rpb24gb2YgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcclxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cclxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxyXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cclxuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBvbmZ1bGZpbGxlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlc29sdmVkLlxyXG4gICAgICogQHBhcmFtIG9ucmVqZWN0ZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZWplY3RlZC5cclxuICAgICAqIEByZXR1cm5zIEEgYENhbmNlbGxhYmxlUHJvbWlzZWAgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHdoaWNoZXZlciBjYWxsYmFjayBpcyBleGVjdXRlZC5cclxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGhvb2tlZCB1cCB0byBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHVwIHRoZSBjaGFpbiwgYnV0IG5vdCBkb3duOlxyXG4gICAgICpcclxuICAgICAqICAgLSBpZiB0aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBiZSBpbnZva2VkIHdpdGggYSBgQ2FuY2VsRXJyb3JgXHJcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIF93aWxsIHJlc29sdmUgcmVndWxhcmx5XyB3aXRoIGl0cyByZXN1bHQ7XHJcbiAgICAgKiAgIC0gY29udmVyc2VseSwgaWYgdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCBfdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCB0b287X1xyXG4gICAgICogICAgIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIHN0aWxsIGJlIGludm9rZWQgd2l0aCB0aGUgcGFyZW50J3MgYENhbmNlbEVycm9yYCxcclxuICAgICAqICAgICBidXQgaXRzIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZFxyXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGEgYENhbmNlbEVycm9yYCBhcyB3ZWxsLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXHJcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cclxuICAgICAqXHJcbiAgICAgKiBJZiBlaXRoZXIgY2FsbGJhY2sgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXHJcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcclxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKi9cclxuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4ge1xyXG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XHJcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLnRoZW4gY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIE5PVEU6IFR5cGVTY3JpcHQncyBidWlsdC1pbiB0eXBlIGZvciB0aGVuIGlzIGJyb2tlbixcclxuICAgICAgICAvLyBhcyBpdCBhbGxvd3Mgc3BlY2lmeWluZyBhbiBhcmJpdHJhcnkgVFJlc3VsdDEgIT0gVCBldmVuIHdoZW4gb25mdWxmaWxsZWQgaXMgbm90IGEgZnVuY3Rpb24uXHJcbiAgICAgICAgLy8gV2UgY2Fubm90IGZpeCBpdCBpZiB3ZSB3YW50IHRvIENhbmNlbGxhYmxlUHJvbWlzZSB0byBpbXBsZW1lbnQgUHJvbWlzZUxpa2U8VD4uXHJcblxyXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZ1bGZpbGxlZCkpIHsgb25mdWxmaWxsZWQgPSBpZGVudGl0eSBhcyBhbnk7IH1cclxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25yZWplY3RlZCkpIHsgb25yZWplY3RlZCA9IHRocm93ZXI7IH1cclxuXHJcbiAgICAgICAgaWYgKG9uZnVsZmlsbGVkID09PSBpZGVudGl0eSAmJiBvbnJlamVjdGVkID09IHRocm93ZXIpIHtcclxuICAgICAgICAgICAgLy8gU2hvcnRjdXQgZm9yIHRyaXZpYWwgYXJndW1lbnRzLlxyXG4gICAgICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSkgPT4gcmVzb2x2ZSh0aGlzIGFzIGFueSkpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgY29uc3QgYmFycmllcjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSB7fTtcclxuICAgICAgICB0aGlzW2JhcnJpZXJTeW1dID0gYmFycmllcjtcclxuXHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICB2b2lkIHN1cGVyLnRoZW4oXHJcbiAgICAgICAgICAgICAgICAodmFsdWUpID0+IHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxyXG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XHJcblxyXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25mdWxmaWxsZWQhKHZhbHVlKSk7XHJcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgICAgICAocmVhc29uPykgPT4ge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XHJcbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcclxuXHJcbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbnJlamVjdGVkIShyZWFzb24pKTtcclxuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICApO1xyXG4gICAgICAgIH0sIGFzeW5jIChjYXVzZT8pID0+IHtcclxuICAgICAgICAgICAgLy9jYW5jZWxsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuIG9uY2FuY2VsbGVkPy4oY2F1c2UpO1xyXG4gICAgICAgICAgICB9IGZpbmFsbHkge1xyXG4gICAgICAgICAgICAgICAgYXdhaXQgdGhpcy5jYW5jZWwoY2F1c2UpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIGZvciBvbmx5IHRoZSByZWplY3Rpb24gb2YgdGhlIFByb21pc2UuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxyXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXHJcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xyXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxyXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogSXQgaXMgZXF1aXZhbGVudCB0b1xyXG4gICAgICogYGBgdHNcclxuICAgICAqIGNhbmNlbGxhYmxlUHJvbWlzZS50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cclxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxyXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcclxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xyXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXHJcbiAgICAgKlxyXG4gICAgICogSWYgYG9ucmVqZWN0ZWRgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXHJcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxyXG4gICAgICovXHJcbiAgICBjYXRjaDxUUmVzdWx0ID0gbmV2ZXI+KG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiAoUHJvbWlzZUxpa2U8VFJlc3VsdD4gfCBUUmVzdWx0KSkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUIHwgVFJlc3VsdD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIHRoYXQgaXMgaW52b2tlZCB3aGVuIHRoZSBDYW5jZWxsYWJsZVByb21pc2UgaXMgc2V0dGxlZCAoZnVsZmlsbGVkIG9yIHJlamVjdGVkKS4gVGhlXHJcbiAgICAgKiByZXNvbHZlZCB2YWx1ZSBjYW5ub3QgYmUgYWNjZXNzZWQgb3IgbW9kaWZpZWQgZnJvbSB0aGUgY2FsbGJhY2suXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHNldHRsZSBpbiB0aGUgc2FtZSBzdGF0ZSBhcyB0aGUgb3JpZ2luYWwgb25lXHJcbiAgICAgKiBhZnRlciB0aGUgcHJvdmlkZWQgY2FsbGJhY2sgaGFzIGNvbXBsZXRlZCBleGVjdXRpb24sXHJcbiAgICAgKiB1bmxlc3MgdGhlIGNhbGxiYWNrIHRocm93cyBvciByZXR1cm5zIGEgcmVqZWN0aW5nIHByb21pc2UsXHJcbiAgICAgKiBpbiB3aGljaCBjYXNlIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IGFzIHdlbGwuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxyXG4gICAgICogT25jZSB0aGUgcGFyZW50IHByb21pc2Ugc2V0dGxlcywgdGhlIGBvbmZpbmFsbHlgIGNhbGxiYWNrIHdpbGwgcnVuLFxyXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cclxuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIFRoaXMgbWV0aG9kIGlzIGltcGxlbWVudGVkIGluIHRlcm1zIG9mIHtAbGluayB0aGVufSBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cclxuICAgICAqIEl0IGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXHJcbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcclxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXHJcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcclxuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxyXG4gICAgICpcclxuICAgICAqIElmIGBvbmZpbmFsbHlgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXHJcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxyXG4gICAgICovXHJcbiAgICBmaW5hbGx5KG9uZmluYWxseT86ICgoKSA9PiB2b2lkKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS5maW5hbGx5IGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25maW5hbGx5KSkge1xyXG4gICAgICAgICAgICByZXR1cm4gdGhpcy50aGVuKG9uZmluYWxseSwgb25maW5hbGx5LCBvbmNhbmNlbGxlZCk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICByZXR1cm4gdGhpcy50aGVuKFxyXG4gICAgICAgICAgICAodmFsdWUpID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHZhbHVlKSxcclxuICAgICAgICAgICAgKHJlYXNvbj8pID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHsgdGhyb3cgcmVhc29uOyB9KSxcclxuICAgICAgICAgICAgb25jYW5jZWxsZWQsXHJcbiAgICAgICAgKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFdlIHVzZSB0aGUgYFtTeW1ib2wuc3BlY2llc11gIHN0YXRpYyBwcm9wZXJ0eSwgaWYgYXZhaWxhYmxlLFxyXG4gICAgICogdG8gZGlzYWJsZSB0aGUgYnVpbHQtaW4gYXV0b21hdGljIHN1YmNsYXNzaW5nIGZlYXR1cmVzIGZyb20ge0BsaW5rIFByb21pc2V9LlxyXG4gICAgICogSXQgaXMgY3JpdGljYWwgZm9yIHBlcmZvcm1hbmNlIHJlYXNvbnMgdGhhdCBleHRlbmRlcnMgZG8gbm90IG92ZXJyaWRlIHRoaXMuXHJcbiAgICAgKiBPbmNlIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nXHJcbiAgICAgKiBpcyBlaXRoZXIgYWNjZXB0ZWQgb3IgcmV0aXJlZCwgdGhpcyBpbXBsZW1lbnRhdGlvbiB3aWxsIGhhdmUgdG8gYmUgcmV2aXNlZCBhY2NvcmRpbmdseS5cclxuICAgICAqXHJcbiAgICAgKiBAaWdub3JlXHJcbiAgICAgKiBAaW50ZXJuYWxcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGdldCBbc3BlY2llc10oKSB7XHJcbiAgICAgICAgcmV0dXJuIFByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXHJcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSwgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgUHJvbWlzZSBpcyByZWplY3RlZC5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBhbGw8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD5bXT47XHJcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBBd2FpdGVkPFRbUF0+OyB9PjtcclxuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xyXG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcclxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXHJcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xyXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUgb3IgcmVqZWN0LlxyXG4gICAgICpcclxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VD4+W10+O1xyXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VFtQXT4+OyB9PjtcclxuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXHJcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxyXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbFNldHRsZWQoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XHJcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBUaGUgYW55IGZ1bmN0aW9uIHJldHVybnMgYSBwcm9taXNlIHRoYXQgaXMgZnVsZmlsbGVkIGJ5IHRoZSBmaXJzdCBnaXZlbiBwcm9taXNlIHRvIGJlIGZ1bGZpbGxlZCxcclxuICAgICAqIG9yIHJlamVjdGVkIHdpdGggYW4gQWdncmVnYXRlRXJyb3IgY29udGFpbmluZyBhbiBhcnJheSBvZiByZWplY3Rpb24gcmVhc29uc1xyXG4gICAgICogaWYgYWxsIG9mIHRoZSBnaXZlbiBwcm9taXNlcyBhcmUgcmVqZWN0ZWQuXHJcbiAgICAgKiBJdCByZXNvbHZlcyBhbGwgZWxlbWVudHMgb2YgdGhlIHBhc3NlZCBpdGVyYWJsZSB0byBwcm9taXNlcyBhcyBpdCBydW5zIHRoaXMgYWxnb3JpdGhtLlxyXG4gICAgICpcclxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGFueTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XHJcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXHJcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxyXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFueShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIGFyZSByZXNvbHZlZCBvciByZWplY3RlZC5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByYWNlPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XHJcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yYWNlKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGNhbmNlbGxlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGNhbmNlbDxUID0gbmV2ZXI+KGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBjb25zdCBwID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XHJcbiAgICAgICAgcC5jYW5jZWwoY2F1c2UpO1xyXG4gICAgICAgIHJldHVybiBwO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBjYW5jZWxzXHJcbiAgICAgKiBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIElmIHRoZSB7QGxpbmsgQWJvcnRTaWduYWwudGltZW91dH0gZmFjdG9yeSBtZXRob2QgaXMgYXZhaWxhYmxlLFxyXG4gICAgICogaXQgaXMgdXNlZCB0byBiYXNlIHRoZSB0aW1lb3V0IG9uIF9hY3RpdmVfIHRpbWUgcmF0aGVyIHRoYW4gX2VsYXBzZWRfIHRpbWUuXHJcbiAgICAgKiBPdGhlcndpc2UsIGB0aW1lb3V0YCBmYWxscyBiYWNrIHRvIHtAbGluayBzZXRUaW1lb3V0fS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHRpbWVvdXQ8VCA9IG5ldmVyPihtaWxsaXNlY29uZHM6IG51bWJlciwgY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcclxuICAgICAgICBpZiAoQWJvcnRTaWduYWwgJiYgdHlwZW9mIEFib3J0U2lnbmFsID09PSAnZnVuY3Rpb24nICYmIEFib3J0U2lnbmFsLnRpbWVvdXQgJiYgdHlwZW9mIEFib3J0U2lnbmFsLnRpbWVvdXQgPT09ICdmdW5jdGlvbicpIHtcclxuICAgICAgICAgICAgQWJvcnRTaWduYWwudGltZW91dChtaWxsaXNlY29uZHMpLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSksIG1pbGxpc2Vjb25kcyk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQuXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgc2xlZXAobWlsbGlzZWNvbmRzOiBudW1iZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXJcclxuICAgICAqIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgdmFsdWUuXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgc2xlZXA8VD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xyXG4gICAgc3RhdGljIHNsZWVwPFQgPSB2b2lkPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU/OiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSkgPT4ge1xyXG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHJlc29sdmUodmFsdWUhKSwgbWlsbGlzZWNvbmRzKTtcclxuICAgICAgICB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVqZWN0ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgcmVhc29uLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmVqZWN0PFQgPSBuZXZlcj4ocmVhc29uPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoXywgcmVqZWN0KSA9PiByZWplY3QocmVhc29uKSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlc29sdmUoKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIHN0YXRpYyByZXNvbHZlPFQgPSB2b2lkPih2YWx1ZT86IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PiB7XHJcbiAgICAgICAgaWYgKHZhbHVlIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSB7XHJcbiAgICAgICAgICAgIC8vIE9wdGltaXNlIGZvciBjYW5jZWxsYWJsZSBwcm9taXNlcy5cclxuICAgICAgICAgICAgcmV0dXJuIHZhbHVlO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+KChyZXNvbHZlKSA9PiByZXNvbHZlKHZhbHVlKSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSBhbmQgcmV0dXJucyBpdCBpbiBhbiBvYmplY3QsIGFsb25nIHdpdGggaXRzIHJlc29sdmUgYW5kIHJlamVjdCBmdW5jdGlvbnNcclxuICAgICAqIGFuZCBhIGdldHRlci9zZXR0ZXIgZm9yIHRoZSBjYW5jZWxsYXRpb24gY2FsbGJhY2suXHJcbiAgICAgKlxyXG4gICAgICogVGhpcyBtZXRob2QgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHdpdGhSZXNvbHZlcnM8VD4oKTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XHJcbiAgICAgICAgbGV0IHJlc3VsdDogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHsgb25jYW5jZWxsZWQ6IG51bGwgfSBhcyBhbnk7XHJcbiAgICAgICAgcmVzdWx0LnByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgcmVzdWx0LnJlc29sdmUgPSByZXNvbHZlO1xyXG4gICAgICAgICAgICByZXN1bHQucmVqZWN0ID0gcmVqZWN0O1xyXG4gICAgICAgIH0sIChjYXVzZT86IGFueSkgPT4geyByZXN1bHQub25jYW5jZWxsZWQ/LihjYXVzZSk7IH0pO1xyXG4gICAgICAgIHJldHVybiByZXN1bHQ7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSBjYW5jZWxsYXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cclxuICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgcmVzdWx0aW5nIGZ1bmN0aW9uIGRvZXMgbm90IHJlamVjdC5cclxuICovXHJcbmZ1bmN0aW9uIGNhbmNlbGxlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpIHtcclxuICAgIGxldCBjYW5jZWxsYXRpb25Qcm9taXNlOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPSB1bmRlZmluZWQ7XHJcblxyXG4gICAgcmV0dXJuIChyZWFzb246IENhbmNlbEVycm9yKTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0+IHtcclxuICAgICAgICBpZiAoIXN0YXRlLnNldHRsZWQpIHtcclxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgIHN0YXRlLnJlYXNvbiA9IHJlYXNvbjtcclxuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcclxuXHJcbiAgICAgICAgICAgIC8vIEF0dGFjaCBhbiBlcnJvciBoYW5kbGVyIHRoYXQgaWdub3JlcyB0aGlzIHNwZWNpZmljIHJlamVjdGlvbiByZWFzb24gYW5kIG5vdGhpbmcgZWxzZS5cclxuICAgICAgICAgICAgLy8gSW4gdGhlb3J5LCBhIHNhbmUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbiBhdCB0aGlzIHBvaW50XHJcbiAgICAgICAgICAgIC8vIHNob3VsZCBhbHdheXMgcmVqZWN0IHdpdGggb3VyIGNhbmNlbGxhdGlvbiByZWFzb24sXHJcbiAgICAgICAgICAgIC8vIGhlbmNlIHRoZSBoYW5kbGVyIHdpbGwgbmV2ZXIgdGhyb3cuXHJcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5wcm90b3R5cGUudGhlbi5jYWxsKHByb21pc2UucHJvbWlzZSwgdW5kZWZpbmVkLCAoZXJyKSA9PiB7XHJcbiAgICAgICAgICAgICAgICBpZiAoZXJyICE9PSByZWFzb24pIHtcclxuICAgICAgICAgICAgICAgICAgICB0aHJvdyBlcnI7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0pO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gSWYgcmVhc29uIGlzIG5vdCBzZXQsIHRoZSBwcm9taXNlIHJlc29sdmVkIHJlZ3VsYXJseSwgaGVuY2Ugd2UgbXVzdCBub3QgY2FsbCBvbmNhbmNlbGxlZC5cclxuICAgICAgICAvLyBJZiBvbmNhbmNlbGxlZCBpcyB1bnNldCwgbm8gbmVlZCB0byBnbyBhbnkgZnVydGhlci5cclxuICAgICAgICBpZiAoIXN0YXRlLnJlYXNvbiB8fCAhcHJvbWlzZS5vbmNhbmNlbGxlZCkgeyByZXR1cm47IH1cclxuXHJcbiAgICAgICAgY2FuY2VsbGF0aW9uUHJvbWlzZSA9IG5ldyBQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICByZXNvbHZlKHByb21pc2Uub25jYW5jZWxsZWQhKHN0YXRlLnJlYXNvbiEuY2F1c2UpKTtcclxuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcclxuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcclxuICAgICAgICB9KTtcclxuXHJcbiAgICAgICAgLy8gVW5zZXQgb25jYW5jZWxsZWQgdG8gcHJldmVudCByZXBlYXRlZCBjYWxscy5cclxuICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gbnVsbDtcclxuXHJcbiAgICAgICAgcmV0dXJuIGNhbmNlbGxhdGlvblByb21pc2U7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZXNvbHV0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXHJcbiAqL1xyXG5mdW5jdGlvbiByZXNvbHZlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiB7XHJcbiAgICByZXR1cm4gKHZhbHVlKSA9PiB7XHJcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cclxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xyXG5cclxuICAgICAgICBpZiAodmFsdWUgPT09IHByb21pc2UucHJvbWlzZSkge1xyXG4gICAgICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cclxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJBIHByb21pc2UgY2Fubm90IGJlIHJlc29sdmVkIHdpdGggaXRzZWxmLlwiKSk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICh2YWx1ZSAhPSBudWxsICYmICh0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ2Z1bmN0aW9uJykpIHtcclxuICAgICAgICAgICAgbGV0IHRoZW46IGFueTtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIHRoZW4gPSAodmFsdWUgYXMgYW55KS50aGVuO1xyXG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QoZXJyKTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUodGhlbikpIHtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgbGV0IGNhbmNlbCA9ICh2YWx1ZSBhcyBhbnkpLmNhbmNlbDtcclxuICAgICAgICAgICAgICAgICAgICBpZiAoaXNDYWxsYWJsZShjYW5jZWwpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IG9uY2FuY2VsbGVkID0gKGNhdXNlPzogYW55KSA9PiB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9O1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoc3RhdGUucmVhc29uKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBJZiBhbHJlYWR5IGNhbmNlbGxlZCwgcHJvcGFnYXRlIGNhbmNlbGxhdGlvbi5cclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIGNhbmNlbGxlciBhbGdvcml0aG0gZG9lcyBub3QgcmVqZWN0XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBzbyBpdCBjYW4gYmUgZGlzY2FyZGVkIHNhZmVseS5cclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgY2FuY2VsbGVyRm9yKHsgLi4ucHJvbWlzZSwgb25jYW5jZWxsZWQgfSwgc3RhdGUpKHN0YXRlLnJlYXNvbik7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gb25jYW5jZWxsZWQ7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgICAgICAgICAgICAgY29uc3QgbmV3U3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xyXG4gICAgICAgICAgICAgICAgICAgIHJvb3Q6IHN0YXRlLnJvb3QsXHJcbiAgICAgICAgICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgICAgICBnZXQgc2V0dGxlZCgpIHsgcmV0dXJuIHRoaXMucm9vdC5zZXR0bGVkIH0sXHJcbiAgICAgICAgICAgICAgICAgICAgc2V0IHNldHRsZWQodmFsdWUpIHsgdGhpcy5yb290LnNldHRsZWQgPSB2YWx1ZTsgfSxcclxuICAgICAgICAgICAgICAgICAgICBnZXQgcmVhc29uKCkgeyByZXR1cm4gdGhpcy5yb290LnJlYXNvbiB9XHJcbiAgICAgICAgICAgICAgICB9O1xyXG5cclxuICAgICAgICAgICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgbmV3U3RhdGUpO1xyXG4gICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KHRoZW4sIHZhbHVlLCBbcmVzb2x2ZXJGb3IocHJvbWlzZSwgbmV3U3RhdGUpLCByZWplY3Rvcl0pO1xyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIHJldHVybjsgLy8gSU1QT1JUQU5UIVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cclxuICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICBwcm9taXNlLnJlc29sdmUodmFsdWUpO1xyXG4gICAgfTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlamVjdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxyXG4gKi9cclxuZnVuY3Rpb24gcmVqZWN0b3JGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3Ige1xyXG4gICAgcmV0dXJuIChyZWFzb24/KSA9PiB7XHJcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cclxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xyXG5cclxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkge1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIHN0YXRlLnJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIE9iamVjdC5pcyhyZWFzb24uY2F1c2UsIHN0YXRlLnJlYXNvbi5jYXVzZSkpIHtcclxuICAgICAgICAgICAgICAgICAgICAvLyBTd2FsbG93IGxhdGUgcmVqZWN0aW9ucyB0aGF0IGFyZSBDYW5jZWxFcnJvcnMgd2hvc2UgY2FuY2VsbGF0aW9uIGNhdXNlIGlzIHRoZSBzYW1lIGFzIG91cnMuXHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbikpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENhbmNlbHMgYWxsIHZhbHVlcyBpbiBhbiBhcnJheSB0aGF0IGxvb2sgbGlrZSBjYW5jZWxsYWJsZSB0aGVuYWJsZXMuXHJcbiAqIFJldHVybnMgYSBwcm9taXNlIHRoYXQgZnVsZmlsbHMgb25jZSBhbGwgY2FuY2VsbGF0aW9uIHByb2NlZHVyZXMgZm9yIHRoZSBnaXZlbiB2YWx1ZXMgaGF2ZSBzZXR0bGVkLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FuY2VsQWxsKHBhcmVudDogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCB2YWx1ZXM6IGFueVtdLCBjYXVzZT86IGFueSk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgY29uc3QgcmVzdWx0cyA9IFtdO1xyXG5cclxuICAgIGZvciAoY29uc3QgdmFsdWUgb2YgdmFsdWVzKSB7XHJcbiAgICAgICAgbGV0IGNhbmNlbDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyO1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZSh2YWx1ZS50aGVuKSkgeyBjb250aW51ZTsgfVxyXG4gICAgICAgICAgICBjYW5jZWwgPSB2YWx1ZS5jYW5jZWw7XHJcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZShjYW5jZWwpKSB7IGNvbnRpbnVlOyB9XHJcbiAgICAgICAgfSBjYXRjaCB7IGNvbnRpbnVlOyB9XHJcblxyXG4gICAgICAgIGxldCByZXN1bHQ6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICByZXN1bHQgPSBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xyXG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XHJcbiAgICAgICAgICAgIGNvbnRpbnVlO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKCFyZXN1bHQpIHsgY29udGludWU7IH1cclxuICAgICAgICByZXN1bHRzLnB1c2goXHJcbiAgICAgICAgICAgIChyZXN1bHQgaW5zdGFuY2VvZiBQcm9taXNlICA/IHJlc3VsdCA6IFByb21pc2UucmVzb2x2ZShyZXN1bHQpKS5jYXRjaCgocmVhc29uPykgPT4ge1xyXG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xyXG4gICAgICAgICAgICB9KVxyXG4gICAgICAgICk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIFByb21pc2UuYWxsKHJlc3VsdHMpIGFzIGFueTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgaXRzIGFyZ3VtZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gaWRlbnRpdHk8VD4oeDogVCk6IFQge1xyXG4gICAgcmV0dXJuIHg7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUaHJvd3MgaXRzIGFyZ3VtZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gdGhyb3dlcihyZWFzb24/OiBhbnkpOiBuZXZlciB7XHJcbiAgICB0aHJvdyByZWFzb247XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBdHRlbXB0cyB2YXJpb3VzIHN0cmF0ZWdpZXMgdG8gY29udmVydCBhbiBlcnJvciB0byBhIHN0cmluZy5cclxuICovXHJcbmZ1bmN0aW9uIGVycm9yTWVzc2FnZShlcnI6IGFueSk6IHN0cmluZyB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGlmIChlcnIgaW5zdGFuY2VvZiBFcnJvciB8fCB0eXBlb2YgZXJyICE9PSAnb2JqZWN0JyB8fCBlcnIudG9TdHJpbmcgIT09IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcpIHtcclxuICAgICAgICAgICAgcmV0dXJuIFwiXCIgKyBlcnI7XHJcbiAgICAgICAgfVxyXG4gICAgfSBjYXRjaCB7fVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgcmV0dXJuIEpTT04uc3RyaW5naWZ5KGVycik7XHJcbiAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZy5jYWxsKGVycik7XHJcbiAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgcmV0dXJuIFwiPGNvdWxkIG5vdCBjb252ZXJ0IGVycm9yIHRvIHN0cmluZz5cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIGN1cnJlbnQgYmFycmllciBwcm9taXNlIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS4gSWYgbmVjZXNzYXJ5LCBpbml0aWFsaXNlcyB0aGUgYmFycmllci5cclxuICovXHJcbmZ1bmN0aW9uIGN1cnJlbnRCYXJyaWVyPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgbGV0IHB3cjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSBwcm9taXNlW2JhcnJpZXJTeW1dID8/IHt9O1xyXG4gICAgaWYgKCEoJ3Byb21pc2UnIGluIHB3cikpIHtcclxuICAgICAgICBPYmplY3QuYXNzaWduKHB3ciwgcHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4oKSk7XHJcbiAgICB9XHJcbiAgICBpZiAocHJvbWlzZVtiYXJyaWVyU3ltXSA9PSBudWxsKSB7XHJcbiAgICAgICAgcHdyLnJlc29sdmUhKCk7XHJcbiAgICAgICAgcHJvbWlzZVtiYXJyaWVyU3ltXSA9IHB3cjtcclxuICAgIH1cclxuICAgIHJldHVybiBwd3IucHJvbWlzZSE7XHJcbn1cclxuXHJcbi8vIFBvbHlmaWxsIFByb21pc2Uud2l0aFJlc29sdmVycy5cclxubGV0IHByb21pc2VXaXRoUmVzb2x2ZXJzID0gUHJvbWlzZS53aXRoUmVzb2x2ZXJzO1xyXG5pZiAocHJvbWlzZVdpdGhSZXNvbHZlcnMgJiYgdHlwZW9mIHByb21pc2VXaXRoUmVzb2x2ZXJzID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IHByb21pc2VXaXRoUmVzb2x2ZXJzLmJpbmQoUHJvbWlzZSk7XHJcbn0gZWxzZSB7XHJcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IGZ1bmN0aW9uIDxUPigpOiBQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XHJcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcclxuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IFByb21pc2U8VD4oKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XHJcbiAgICAgICAgcmV0dXJuIHsgcHJvbWlzZSwgcmVzb2x2ZSwgcmVqZWN0IH07XHJcbiAgICB9XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2xpcGJvYXJkKTtcclxuXHJcbmNvbnN0IENsaXBib2FyZFNldFRleHQgPSAwO1xyXG5jb25zdCBDbGlwYm9hcmRUZXh0ID0gMTtcclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSB0ZXh0IHRvIHRoZSBDbGlwYm9hcmQuXHJcbiAqXHJcbiAqIEBwYXJhbSB0ZXh0IC0gVGhlIHRleHQgdG8gYmUgc2V0IHRvIHRoZSBDbGlwYm9hcmQuXHJcbiAqIEByZXR1cm4gQSBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2hlbiB0aGUgb3BlcmF0aW9uIGlzIHN1Y2Nlc3NmdWwuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0OiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKTogUHJvbWlzZTxzdHJpbmc+IHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKipcclxuICogQW55IGlzIGEgZHVtbXkgY3JlYXRpb24gZnVuY3Rpb24gZm9yIHNpbXBsZSBvciB1bmtub3duIHR5cGVzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEFueTxUID0gYW55Pihzb3VyY2U6IGFueSk6IFQge1xyXG4gICAgcmV0dXJuIHNvdXJjZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEJ5dGVTbGljZSBpcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIHRoYXQgcmVwbGFjZXNcclxuICogbnVsbCBzdHJpbmdzIHdpdGggZW1wdHkgc3RyaW5ncy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeXRlU2xpY2Uoc291cmNlOiBhbnkpOiBzdHJpbmcge1xyXG4gICAgcmV0dXJuICgoc291cmNlID09IG51bGwpID8gXCJcIiA6IHNvdXJjZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBBcnJheSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxyXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJyYXlcclxuICogd2hvc2UgZWxlbWVudHMgYXJlIG9mIHRoYXQgdHlwZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBBcnJheTxUID0gYW55PihlbGVtZW50OiAoc291cmNlOiBhbnkpID0+IFQpOiAoc291cmNlOiBhbnkpID0+IFRbXSB7XHJcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBbXSA6IHNvdXJjZSk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcclxuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XHJcbiAgICAgICAgICAgIHJldHVybiBbXTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzb3VyY2UubGVuZ3RoOyBpKyspIHtcclxuICAgICAgICAgICAgc291cmNlW2ldID0gZWxlbWVudChzb3VyY2VbaV0pO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gc291cmNlO1xyXG4gICAgfTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE1hcCB0YWtlcyBjcmVhdGlvbiBmdW5jdGlvbnMgZm9yIHR3byBhcmJpdHJhcnkgdHlwZXNcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIG9iamVjdFxyXG4gKiB3aG9zZSBrZXlzIGFuZCB2YWx1ZXMgYXJlIG9mIHRob3NlIHR5cGVzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE1hcDxWID0gYW55PihrZXk6IChzb3VyY2U6IGFueSkgPT4gc3RyaW5nLCB2YWx1ZTogKHNvdXJjZTogYW55KSA9PiBWKTogKHNvdXJjZTogYW55KSA9PiBSZWNvcmQ8c3RyaW5nLCBWPiB7XHJcbiAgICBpZiAodmFsdWUgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8ge30gOiBzb3VyY2UpO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm4ge307XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGZvciAoY29uc3Qga2V5IGluIHNvdXJjZSkge1xyXG4gICAgICAgICAgICBzb3VyY2Vba2V5XSA9IHZhbHVlKHNvdXJjZVtrZXldKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBOdWxsYWJsZSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxyXG4gKiBhbmQgcmV0dXJucyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhIG51bGxhYmxlIHZhbHVlIG9mIHRoYXQgdHlwZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBOdWxsYWJsZTxUID0gYW55PihlbGVtZW50OiAoc291cmNlOiBhbnkpID0+IFQpOiAoc291cmNlOiBhbnkpID0+IChUIHwgbnVsbCkge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiBBbnk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFN0cnVjdCB0YWtlcyBhbiBvYmplY3QgbWFwcGluZyBmaWVsZCBuYW1lcyB0byBjcmVhdGlvbiBmdW5jdGlvbnNcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgc3RydWN0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFN0cnVjdChjcmVhdGVGaWVsZDogUmVjb3JkPHN0cmluZywgKHNvdXJjZTogYW55KSA9PiBhbnk+KTpcclxuICAgIDxVIGV4dGVuZHMgUmVjb3JkPHN0cmluZywgYW55PiA9IGFueT4oc291cmNlOiBhbnkpID0+IFVcclxue1xyXG4gICAgbGV0IGFsbEFueSA9IHRydWU7XHJcbiAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcclxuICAgICAgICBpZiAoY3JlYXRlRmllbGRbbmFtZV0gIT09IEFueSkge1xyXG4gICAgICAgICAgICBhbGxBbnkgPSBmYWxzZTtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG4gICAgaWYgKGFsbEFueSkge1xyXG4gICAgICAgIHJldHVybiBBbnk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcclxuICAgICAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcclxuICAgICAgICAgICAgaWYgKG5hbWUgaW4gc291cmNlKSB7XHJcbiAgICAgICAgICAgICAgICBzb3VyY2VbbmFtZV0gPSBjcmVhdGVGaWVsZFtuYW1lXShzb3VyY2VbbmFtZV0pO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBzb3VyY2U7XHJcbiAgICB9O1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFNpemUge1xyXG4gICAgLyoqIFRoZSB3aWR0aCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXHJcbiAgICBXaWR0aDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgYSByZWN0YW5ndWxhciBhcmVhLiAqL1xyXG4gICAgSGVpZ2h0OiBudW1iZXI7XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgUmVjdCB7XHJcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xyXG4gICAgWDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cclxuICAgIFk6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgd2lkdGggb2YgdGhlIHJlY3RhbmdsZS4gKi9cclxuICAgIFdpZHRoOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xyXG4gICAgSGVpZ2h0OiBudW1iZXI7XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgU2NyZWVuIHtcclxuICAgIC8qKiBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi4gKi9cclxuICAgIElEOiBzdHJpbmc7XHJcbiAgICAvKiogSHVtYW4tcmVhZGFibGUgbmFtZSBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgTmFtZTogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSBzY2FsZSBmYWN0b3Igb2YgdGhlIHNjcmVlbiAoRFBJLzk2KS4gMSA9IHN0YW5kYXJkIERQSSwgMiA9IEhpRFBJIChSZXRpbmEpLCBldGMuICovXHJcbiAgICBTY2FsZUZhY3RvcjogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBZOiBudW1iZXI7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIHdpZHRoIGFuZCBoZWlnaHQgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFNpemU6IFNpemU7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LiAqL1xyXG4gICAgQm91bmRzOiBSZWN0O1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodCAoYmVmb3JlIHNjYWxpbmcpLiAqL1xyXG4gICAgUGh5c2ljYWxCb3VuZHM6IFJlY3Q7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIGFyZWEgb2YgdGhlIHNjcmVlbiB0aGF0IGlzIGFjdHVhbGx5IHVzYWJsZSAoZXhjbHVkaW5nIHRhc2tiYXIgYW5kIG90aGVyIHN5c3RlbSBVSSkuICovXHJcbiAgICBXb3JrQXJlYTogUmVjdDtcclxuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgV29ya0FyZWEgb2YgdGhlIHNjcmVlbiAoYmVmb3JlIHNjYWxpbmcpLiAqL1xyXG4gICAgUGh5c2ljYWxXb3JrQXJlYTogUmVjdDtcclxuICAgIC8qKiBUcnVlIGlmIHRoaXMgaXMgdGhlIHByaW1hcnkgbW9uaXRvciBzZWxlY3RlZCBieSB0aGUgdXNlciBpbiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cclxuICAgIElzUHJpbWFyeTogYm9vbGVhbjtcclxuICAgIC8qKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFJvdGF0aW9uOiBudW1iZXI7XHJcbn1cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TY3JlZW5zKTtcclxuXHJcbmNvbnN0IGdldEFsbCA9IDA7XHJcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xyXG5jb25zdCBnZXRDdXJyZW50ID0gMjtcclxuXHJcbi8qKlxyXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKTogUHJvbWlzZTxTY3JlZW5bXT4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCk6IFByb21pc2U8U2NyZWVuPiB7XHJcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcclxufVxyXG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNBQTtBQUFBO0FBQUE7QUFBQTs7O0FDNkJBLElBQU0sY0FDRjtBQUVHLFNBQVMsT0FBTyxPQUFlLElBQVk7QUFDOUMsTUFBSSxLQUFLO0FBRVQsTUFBSSxJQUFJLE9BQU87QUFDZixTQUFPLEtBQUs7QUFFUixVQUFNLFlBQWEsS0FBSyxPQUFPLElBQUksS0FBTSxDQUFDO0FBQUEsRUFDOUM7QUFDQSxTQUFPO0FBQ1g7OztBQzdCQSxJQUFNLGFBQWEsT0FBTyxTQUFTLFNBQVM7QUFHckMsSUFBTSxjQUFjLE9BQU8sT0FB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SUFDckIsNkJBQTZCO0FBQUEsSUFDN0IsMEJBQTBCO0FBQUEsSUFDMUIsdUJBQXVCO0FBQUEsSUFDdkIsdUJBQXVCO0FBQUEsSUFDdkIsZ0JBQWdCO0FBQUEsSUFDaEIsc0JBQXNCO0FBQUEsSUFDdEIsY0FBYztBQUFBLElBQ2Qsb0JBQW9CO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsc0JBQXNCO0FBQUEsSUFDdEIsYUFBYTtBQUFBLElBQ2IsY0FBYztBQUFBLElBQ2QsbUJBQW1CO0FBQUEsSUFDbkIsbUJBQW1CO0FBQUEsSUFDbkIseUJBQXlCO0FBQUEsSUFDekIsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsSUFDakIsdUJBQXVCO0FBQUEsSUFDdkIscUJBQXFCO0FBQUEsSUFDckIscUJBQXFCO0FBQUEsSUFDckIsdUJBQXVCO0FBQUEsSUFDdkIsY0FBYztBQUFBLElBQ2QsZUFBZTtBQUFBLElBQ2Ysb0JBQW9CO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsMEJBQTBCO0FBQUEsSUFDMUIsZ0JBQWdCO0FBQUEsSUFDaEIsNEJBQTRCO0FBQUEsSUFDNUIsNEJBQTRCO0FBQUEsSUFDNUIseURBQXlEO0FBQUEsSUFDekQsc0NBQXNDO0FBQUEsSUFDdEMsb0JBQW9CO0FBQUEsSUFDcEIscUJBQXFCO0FBQUEsSUFDckIscUJBQXFCO0FBQUEsSUFDckIsc0JBQXNCO0FBQUEsSUFDdEIsZ0NBQWdDO0FBQUEsSUFDaEMsa0NBQWtDO0FBQUEsSUFDbEMsbUNBQW1DO0FBQUEsSUFDbkMsb0NBQW9DO0FBQUEsSUFDcEMsK0JBQStCO0FBQUEsSUFDL0IsNkJBQTZCO0FBQUEsSUFDN0IsdUJBQXVCO0FBQUEsSUFDdkIsaUNBQWlDO0FBQUEsSUFDakMsOEJBQThCO0FBQUEsSUFDOUIsNEJBQTRCO0FBQUEsSUFDNUIsc0NBQXNDO0FBQUEsSUFDdEMsNEJBQTRCO0FBQUEsSUFDNUIsc0JBQXNCO0FBQUEsSUFDdEIsa0NBQWtDO0FBQUEsSUFDbEMsc0JBQXNCO0FBQUEsSUFDdEIsd0JBQXdCO0FBQUEsSUFDeEIsd0JBQXdCO0FBQUEsSUFDeEIsbUJBQW1CO0FBQUEsSUFDbkIsMEJBQTBCO0FBQUEsSUFDMUIsOEJBQThCO0FBQUEsSUFDOUIseUJBQXlCO0FBQUEsSUFDekIsNkJBQTZCO0FBQUEsSUFDN0IsaUJBQWlCO0FBQUEsSUFDakIsZ0JBQWdCO0FBQUEsSUFDaEIsc0JBQXNCO0FBQUEsSUFDdEIsZUFBZTtBQUFBLElBQ2YseUJBQXlCO0FBQUEsSUFDekIsd0JBQXdCO0FBQUEsSUFDeEIsb0JBQW9CO0FBQUEsSUFDcEIscUJBQXFCO0FBQUEsSUFDckIsaUJBQWlCO0FBQUEsSUFDakIsaUJBQWlCO0FBQUEsSUFDakIsc0JBQXNCO0FBQUEsSUFDdEIsbUNBQW1DO0FBQUEsSUFDbkMscUNBQXFDO0FBQUEsSUFDckMsdUJBQXVCO0FBQUEsSUFDdkIsc0JBQXNCO0FBQUEsSUFDdEIsd0JBQXdCO0FBQUEsSUFDeEIsZUFBZTtBQUFBLElBQ2YsMkJBQTJCO0FBQUEsSUFDM0IsMEJBQTBCO0FBQUEsSUFDMUIsNkJBQTZCO0FBQUEsSUFDN0IsWUFBWTtBQUFBLElBQ1osZ0JBQWdCO0FBQUEsSUFDaEIsa0JBQWtCO0FBQUEsSUFDbEIsZ0JBQWdCO0FBQUEsSUFDaEIsa0JBQWtCO0FBQUEsSUFDbEIsbUJBQW1CO0FBQUEsSUFDbkIsWUFBWTtBQUFBLElBQ1oscUJBQXFCO0FBQUEsSUFDckIsc0JBQXNCO0FBQUEsSUFDdEIsc0JBQXNCO0FBQUEsSUFDdEIsOEJBQThCO0FBQUEsSUFDOUIsaUJBQWlCO0FBQUEsSUFDakIseUJBQXlCO0FBQUEsSUFDekIsMkJBQTJCO0FBQUEsSUFDM0IsK0JBQStCO0FBQUEsSUFDL0IsMEJBQTBCO0FBQUEsSUFDMUIsOEJBQThCO0FBQUEsSUFDOUIsaUJBQWlCO0FBQUEsSUFDakIsdUJBQXVCO0FBQUEsSUFDdkIsZ0JBQWdCO0FBQUEsSUFDaEIsMEJBQTBCO0FBQUEsSUFDMUIseUJBQXlCO0FBQUEsSUFDekIsc0JBQXNCO0FBQUEsSUFDdEIsa0JBQWtCO0FBQUEsSUFDbEIsbUJBQW1C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QUFBLEVBQ1g7QUFDSjs7O0FDekdBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQXdEQSxJQUFNQyxRQUFPLGlCQUFpQixZQUFZLE9BQU87QUFFakQsSUFBTSxTQUFTO0FBQ2YsSUFBTSxhQUFhO0FBQ25CLElBQU0sYUFBYTtBQU9aLFNBQVMsU0FBNEI7QUFDeEMsU0FBT0EsTUFBSyxNQUFNO0FBQ3RCO0FBT08sU0FBUyxhQUE4QjtBQUMxQyxTQUFPQSxNQUFLLFVBQVU7QUFDMUI7QUFPTyxTQUFTLGFBQThCO0FBQzFDLFNBQU9BLE1BQUssVUFBVTtBQUMxQjs7O0F0QjVFQSxPQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUE0Q2xDLE9BQU8sT0FBTyxTQUFnQjtBQUN2QixPQUFPLHFCQUFxQjsiLAogICJuYW1lcyI6IFsiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJfYSIsICJFcnJvciIsICJjYWxsIiwgIl9hIiwgInJlc2l6YWJsZSIsICJjYWxsIiwgIl9hIiwgImNhbGwiLCAiY2FsbCIsICJIaWRlTWV0aG9kIiwgIlNob3dNZXRob2QiLCAiaXNEb2N1bWVudERvdEFsbCIsICJfYSIsICJyZWFzb24iLCAidmFsdWUiLCAiY2FsbCIsICJnZXRBbmREZWxldGVSZXNwb25zZSIsICJnZW5lcmF0ZUlEIiwgImNhbGwiLCAiQXJyYXkiLCAiTWFwIiwgIkFycmF5IiwgIk1hcCIsICJrZXkiLCAiY2FsbCJdCn0K
