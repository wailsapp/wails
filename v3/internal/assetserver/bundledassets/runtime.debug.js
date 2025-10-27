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
  Base64JSONCodec: () => Base64JSONCodec,
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
  RawJSONCodec: () => RawJSONCodec,
  RawStringCodec: () => RawStringCodec,
  Screens: () => screens_exports,
  System: () => system_exports,
  WML: () => wml_exports,
  Window: () => window_default,
  buildTransportRequest: () => buildTransportRequest,
  clientId: () => clientId,
  dispatchWailsEvent: () => dispatchWailsEvent,
  generateMessageID: () => generateMessageID,
  getTransport: () => getTransport,
  handleWailsCallback: () => handleWailsCallback,
  objectNames: () => objectNames,
  setTransport: () => setTransport
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
var customTransport = null;
function setTransport(transport) {
  customTransport = transport;
}
function getTransport() {
  return customTransport;
}
function newRuntimeCaller(object, windowName = "") {
  return function(method, args = null) {
    return runtimeCallWithID(object, method, windowName, args);
  };
}
async function runtimeCallWithID(objectID, method, windowName, args) {
  var _a2, _b;
  if (customTransport) {
    return customTransport.call(objectID, method, windowName, args);
  }
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
function generateMessageID() {
  return nanoid();
}
function buildTransportRequest(id, objectID, method, windowName, args) {
  return {
    id,
    type: "request",
    request: {
      object: objectID,
      method,
      args: args ? JSON.stringify(args) : void 0,
      windowName: windowName || void 0,
      clientId
    }
  };
}
function handleWailsCallback(pending, responseData, contentType) {
  var _a2, _b, _c;
  if (!responseData || !(contentType == null ? void 0 : contentType.includes("application/json"))) {
    return false;
  }
  const callId = (_b = (_a2 = pending.request) == null ? void 0 : _a2.args) == null ? void 0 : _b["call-id"];
  if (callId && ((_c = window._wails) == null ? void 0 : _c.callResultHandler)) {
    window._wails.callResultHandler(callId, responseData, true);
    return true;
  }
  return false;
}
function dispatchWailsEvent(event) {
  var _a2;
  if ((_a2 = window._wails) == null ? void 0 : _a2.dispatchWailsEvent) {
    window._wails.dispatchWailsEvent(event);
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
window._wails.dispatchWailsEvent = dispatchWailsEvent2;
var call3 = newRuntimeCaller(objectNames.Events);
var EmitMethod = 0;
var WailsEvent = class {
  constructor(name, data = null) {
    this.name = name;
    this.data = data;
  }
};
function dispatchWailsEvent2(event) {
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
var SnapAssistMethod = 49;
var WindowDropZoneDropped = 50;
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
  /* Triggers Windows 11 Snap Assist feature (Windows only).
   * This is equivalent to pressing Win+Z and shows snap layout options.
   */
  SnapAssist() {
    return this[callerSym](SnapAssistMethod);
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
  Emit(eventName, data);
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
var _invoke = (function() {
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
})();
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

// desktop/@wailsio/runtime/src/transport-codec.ts
var Base64JSONCodec = class {
  /**
   * Decode base64-encoded response data to UTF-8 string
   */
  decodeResponse(data, contentType) {
    if (!data) return "";
    try {
      const binaryString = atob(data);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      return new TextDecoder().decode(bytes);
    } catch (err) {
      throw new Error("Failed to decode response: ".concat(err instanceof Error ? err.message : String(err)));
    }
  }
  /**
   * Decode base64-encoded error data to UTF-8 string
   */
  decodeError(data) {
    if (!data) return "Unknown error";
    try {
      return this.decodeResponse(data, "text/plain");
    } catch (e) {
      return "Failed to decode error message";
    }
  }
};
var RawStringCodec = class {
  decodeResponse(data, contentType) {
    return data ? String(data) : "";
  }
  decodeError(data) {
    return data ? String(data) : "Unknown error";
  }
};
var RawJSONCodec = class {
  decodeResponse(data, contentType) {
    if (!data) return "";
    return typeof data === "string" ? data : JSON.stringify(data);
  }
  decodeError(data) {
    if (!data) return "Unknown error";
    return typeof data === "string" ? data : JSON.stringify(data);
  }
};

// desktop/@wailsio/runtime/src/index.ts
window._wails = window._wails || {};
window._wails.invoke = invoke;
window._wails.handlePlatformFileDrop = window_default.HandlePlatformFileDrop.bind(window_default);
invoke("wails:runtime:ready");
export {
  application_exports as Application,
  Base64JSONCodec,
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
  RawJSONCodec,
  RawStringCodec,
  screens_exports as Screens,
  system_exports as System,
  wml_exports as WML,
  window_default as Window,
  buildTransportRequest,
  clientId,
  dispatchWailsEvent,
  generateMessageID,
  getTransport,
  handleWailsCallback,
  objectNames,
  setTransport
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHBsaWNhdGlvbi50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYW5jZWxsYWJsZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2xpcGJvYXJkLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jcmVhdGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3RyYW5zcG9ydC1jb2RlYy50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxuaW1wb3J0IFwiLi9jb250ZXh0bWVudS5qc1wiO1xuaW1wb3J0IFwiLi9kcmFnLmpzXCI7XG5cbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi9hcHBsaWNhdGlvbi5qc1wiO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyLmpzXCI7XG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzLmpzXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkLmpzXCI7XG5pbXBvcnQgKiBhcyBDcmVhdGUgZnJvbSBcIi4vY3JlYXRlLmpzXCI7XG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcbmltcG9ydCAqIGFzIEV2ZW50cyBmcm9tIFwiLi9ldmVudHMuanNcIjtcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuL2ZsYWdzLmpzXCI7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnMuanNcIjtcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93LmpzXCI7XG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sLmpzXCI7XG5cbmV4cG9ydCB7XG4gICAgQXBwbGljYXRpb24sXG4gICAgQnJvd3NlcixcbiAgICBDYWxsLFxuICAgIENsaXBib2FyZCxcbiAgICBEaWFsb2dzLFxuICAgIEV2ZW50cyxcbiAgICBGbGFncyxcbiAgICBTY3JlZW5zLFxuICAgIFN5c3RlbSxcbiAgICBXaW5kb3csXG4gICAgV01MXG59O1xuXG4vKipcbiAqIEFuIGludGVybmFsIHV0aWxpdHkgY29uc3VtZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxuICpcbiAqIEBpZ25vcmVcbiAqIEBpbnRlcm5hbFxuICovXG5leHBvcnQgeyBDcmVhdGUgfTtcblxuZXhwb3J0ICogZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcblxuLy8gRXhwb3J0IHRyYW5zcG9ydCBpbnRlcmZhY2VzIGFuZCB1dGlsaXRpZXNcbmV4cG9ydCB7XG4gICAgc2V0VHJhbnNwb3J0LFxuICAgIGdldFRyYW5zcG9ydCxcbiAgICB0eXBlIFJ1bnRpbWVUcmFuc3BvcnQsXG4gICAgb2JqZWN0TmFtZXMsXG4gICAgY2xpZW50SWQsXG4gICAgLy8gSGVscGVyIHV0aWxpdGllcyBmb3IgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnNcbiAgICBnZW5lcmF0ZU1lc3NhZ2VJRCxcbiAgICBidWlsZFRyYW5zcG9ydFJlcXVlc3QsXG4gICAgaGFuZGxlV2FpbHNDYWxsYmFjayxcbiAgICBkaXNwYXRjaFdhaWxzRXZlbnRcbn0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG4vLyBFeHBvcnQgY29kZWMgaW50ZXJmYWNlcyBhbmQgaW1wbGVtZW50YXRpb25zXG5leHBvcnQge1xuICAgIHR5cGUgVHJhbnNwb3J0Q29kZWMsXG4gICAgQmFzZTY0SlNPTkNvZGVjLFxuICAgIFJhd1N0cmluZ0NvZGVjLFxuICAgIFJhd0pTT05Db2RlY1xufSBmcm9tIFwiLi90cmFuc3BvcnQtY29kZWMuanNcIjtcblxuLy8gTm90aWZ5IGJhY2tlbmRcbndpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcblxuLy8gUmVnaXN0ZXIgcGxhdGZvcm0gaGFuZGxlcnMgKGludGVybmFsIEFQSSlcbi8vIE5vdGU6IFdpbmRvdyBpcyB0aGUgdGhpc1dpbmRvdyBpbnN0YW5jZSAoZGVmYXVsdCBleHBvcnQgZnJvbSB3aW5kb3cudHMpXG4vLyBCaW5kaW5nIGVuc3VyZXMgJ3RoaXMnIGNvcnJlY3RseSByZWZlcnMgdG8gdGhlIGN1cnJlbnQgd2luZG93IGluc3RhbmNlXG53aW5kb3cuX3dhaWxzLmhhbmRsZVBsYXRmb3JtRmlsZURyb3AgPSBXaW5kb3cuSGFuZGxlUGxhdGZvcm1GaWxlRHJvcC5iaW5kKFdpbmRvdyk7XG5cblN5c3RlbS5pbnZva2UoXCJ3YWlsczpydW50aW1lOnJlYWR5XCIpO1xuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBPcGVuVVJMIH0gZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xuaW1wb3J0IHsgUXVlc3Rpb24gfSBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XG5pbXBvcnQgeyBFbWl0IH0gZnJvbSBcIi4vZXZlbnRzLmpzXCI7XG5pbXBvcnQgeyBjYW5BYm9ydExpc3RlbmVycywgd2hlblJlYWR5IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93LmpzXCI7XG5cbi8qKlxuICogU2VuZHMgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgb3B0aW9uYWwgZGF0YS5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gc2VuZC5cbiAqIEBwYXJhbSBbZGF0YT1udWxsXSAtIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxuICovXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpOiB2b2lkIHtcbiAgICBFbWl0KGV2ZW50TmFtZSwgZGF0YSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBzcGVjaWZpZWQgd2luZG93LlxuICpcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqL1xuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lOiBzdHJpbmcsIG1ldGhvZE5hbWU6IHN0cmluZykge1xuICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IFdpbmRvdy5HZXQod2luZG93TmFtZSk7XG4gICAgY29uc3QgbWV0aG9kID0gKHRhcmdldFdpbmRvdyBhcyBhbnkpW21ldGhvZE5hbWVdO1xuXG4gICAgaWYgKHR5cGVvZiBtZXRob2QgIT09IFwiZnVuY3Rpb25cIikge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBXaW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JyBub3QgZm91bmRgKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICAgIG1ldGhvZC5jYWxsKHRhcmdldFdpbmRvdyk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmVzcG9uZHMgdG8gYSB0cmlnZ2VyaW5nIGV2ZW50IGJ5IHJ1bm5pbmcgYXBwcm9wcmlhdGUgV01MIGFjdGlvbnMgZm9yIHRoZSBjdXJyZW50IHRhcmdldC5cbiAqL1xuZnVuY3Rpb24gb25XTUxUcmlnZ2VyZWQoZXY6IEV2ZW50KTogdm9pZCB7XG4gICAgY29uc3QgZWxlbWVudCA9IGV2LmN1cnJlbnRUYXJnZXQgYXMgRWxlbWVudDtcblxuICAgIGZ1bmN0aW9uIHJ1bkVmZmVjdChjaG9pY2UgPSBcIlllc1wiKSB7XG4gICAgICAgIGlmIChjaG9pY2UgIT09IFwiWWVzXCIpXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1ldmVudCcpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCB1cmwgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLW9wZW51cmwnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtb3BlbnVybCcpO1xuXG4gICAgICAgIGlmIChldmVudFR5cGUgIT09IG51bGwpXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICBpZiAodXJsICE9PSBudWxsKVxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgfVxuXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG5cbiAgICBpZiAoY29uZmlybSkge1xuICAgICAgICBRdWVzdGlvbih7XG4gICAgICAgICAgICBUaXRsZTogXCJDb25maXJtXCIsXG4gICAgICAgICAgICBNZXNzYWdlOiBjb25maXJtLFxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxuICAgICAgICAgICAgQnV0dG9uczogW1xuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiWWVzXCIgfSxcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XG4gICAgICAgICAgICBdXG4gICAgICAgIH0pLnRoZW4ocnVuRWZmZWN0KTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBydW5FZmZlY3QoKTtcbiAgICB9XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBjb250cm9sbGVyU3ltID0gU3ltYm9sKFwiY29udHJvbGxlclwiKTtcbmNvbnN0IHRyaWdnZXJNYXBTeW0gPSBTeW1ib2woXCJ0cmlnZ2VyTWFwXCIpO1xuY29uc3QgZWxlbWVudENvdW50U3ltID0gU3ltYm9sKFwiZWxlbWVudENvdW50XCIpO1xuXG4vKipcbiAqIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IGRvZXMgbm90IGFjdHVhbGx5IHJlbWVtYmVyIGFjdGl2ZSBldmVudCBsaXN0ZW5lcnM6IGluc3RlYWRcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXG4gKi9cbmNsYXNzIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICBbY29udHJvbGxlclN5bV06IEFib3J0Q29udHJvbGxlcjtcblxuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgYW4gb3B0aW9ucyBvYmplY3QgZm9yIGFkZEV2ZW50TGlzdGVuZXIgdGhhdCB0aWVzIHRoZSBsaXN0ZW5lclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xuICAgICAqL1xuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XG4gICAgICAgIHJldHVybiB7IHNpZ25hbDogdGhpc1tjb250cm9sbGVyU3ltXS5zaWduYWwgfTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycyBhbmQgcmVzZXRzIHRoZSByZWdpc3RyeS5cbiAgICAgKi9cbiAgICByZXNldCgpOiB2b2lkIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXS5hYm9ydCgpO1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBXZWFrTWFwUmVnaXN0cnkgbWFwcyBhY3RpdmUgdHJpZ2dlciBldmVudHMgdG8gZWFjaCBET00gZWxlbWVudCB0aHJvdWdoIGEgV2Vha01hcC5cbiAqIFRoaXMgZW5zdXJlcyB0aGF0IHRoZSBtYXBwaW5nIHJlbWFpbnMgcHJpdmF0ZSB0byB0aGlzIG1vZHVsZSwgd2hpbGUgc3RpbGwgYWxsb3dpbmcgZ2FyYmFnZVxuICogY29sbGVjdGlvbiBvZiB0aGUgaW52b2x2ZWQgZWxlbWVudHMuXG4gKi9cbmNsYXNzIFdlYWtNYXBSZWdpc3RyeSB7XG4gICAgLyoqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy4gKi9cbiAgICBbdHJpZ2dlck1hcFN5bV06IFdlYWtNYXA8RWxlbWVudCwgc3RyaW5nW10+O1xuICAgIC8qKiBDb3VudHMgdGhlIG51bWJlciBvZiBlbGVtZW50cyB3aXRoIGFjdGl2ZSBXTUwgdHJpZ2dlcnMuICovXG4gICAgW2VsZW1lbnRDb3VudFN5bV06IG51bWJlcjtcblxuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIGFjdGl2ZSB0cmlnZ2VycyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50LlxuICAgICAqXG4gICAgICogQHBhcmFtIGVsZW1lbnQgLSBBbiBIVE1MIGVsZW1lbnRcbiAgICAgKiBAcGFyYW0gdHJpZ2dlcnMgLSBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRcbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xuICAgICAgICBpZiAoIXRoaXNbdHJpZ2dlck1hcFN5bV0uaGFzKGVsZW1lbnQpKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXSsrOyB9XG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICAgICAgcmV0dXJuIHt9O1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzLlxuICAgICAqL1xuICAgIHJlc2V0KCk6IHZvaWQge1xuICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnKicpKSB7XG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXG4gICAgICAgICAgICAgICAgYnJlYWs7XG5cbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwU3ltXS5nZXQoZWxlbWVudCk7XG4gICAgICAgICAgICBpZiAodHJpZ2dlcnMgIT0gbnVsbCkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0tLTsgfVxuXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXG4gICAgICAgICAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkKTtcbiAgICAgICAgfVxuXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0gPSBuZXcgV2Vha01hcCgpO1xuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xuICAgIH1cbn1cblxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xuXG4vKipcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAqL1xuZnVuY3Rpb24gYWRkV01MTGlzdGVuZXJzKGVsZW1lbnQ6IEVsZW1lbnQpOiB2b2lkIHtcbiAgICBjb25zdCB0cmlnZ2VyUmVnRXhwID0gL1xcUysvZztcbiAgICBjb25zdCB0cmlnZ2VyQXR0ciA9IChlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIik7XG4gICAgY29uc3QgdHJpZ2dlcnM6IHN0cmluZ1tdID0gW107XG5cbiAgICBsZXQgbWF0Y2g7XG4gICAgd2hpbGUgKChtYXRjaCA9IHRyaWdnZXJSZWdFeHAuZXhlYyh0cmlnZ2VyQXR0cikpICE9PSBudWxsKVxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcblxuICAgIGNvbnN0IG9wdGlvbnMgPSB0cmlnZ2VyUmVnaXN0cnkuc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2NoZWR1bGVzIGFuIGF1dG9tYXRpYyByZWxvYWQgb2YgV01MIHRvIGJlIHBlcmZvcm1lZCBhcyBzb29uIGFzIHRoZSBkb2N1bWVudCBpcyBmdWxseSBsb2FkZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbmFibGUoKTogdm9pZCB7XG4gICAgd2hlblJlYWR5KFJlbG9hZCk7XG59XG5cbi8qKlxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCk6IHZvaWQge1xuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xuICAgIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0sIFt3bWwtd2luZG93XSwgW3dtbC1vcGVudXJsXSwgW2RhdGEtd21sLWV2ZW50XSwgW2RhdGEtd21sLXdpbmRvd10sIFtkYXRhLXdtbC1vcGVudXJsXScpLmZvckVhY2goYWRkV01MTGlzdGVuZXJzKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xuXG5jb25zdCBCcm93c2VyT3BlblVSTCA9IDA7XG5cbi8qKlxuICogT3BlbiBhIGJyb3dzZXIgd2luZG93IHRvIHRoZSBnaXZlbiBVUkwuXG4gKlxuICogQHBhcmFtIHVybCAtIFRoZSBVUkwgdG8gb3BlblxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmw6IHN0cmluZyB8IFVSTCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsOiB1cmwudG9TdHJpbmcoKX0pO1xufVxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxuXG4vLyBUaGUgTUlUIExpY2Vuc2UgKE1JVClcbi8vXG4vLyBDb3B5cmlnaHQgMjAxNyBBbmRyZXkgU2l0bmlrIDxhbmRyZXlAc2l0bmlrLnJ1PlxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2Zcbi8vIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWwgaW5cbi8vIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMgdG9cbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXG4vLyB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sXG4vLyAgICAgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6XG4vL1xuLy8gICAgIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cbi8vXG4vLyAgICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksIEZJVE5FU1Ncbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxuLy8gQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLCBXSEVUSEVSXG4vLyBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sIE9VVCBPRiBPUiBJTlxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cblxuLy8gVGhpcyBhbHBoYWJldCB1c2VzIGBBLVphLXowLTlfLWAgc3ltYm9scy5cbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cbi8vIFJlZmVyZW5jZXMgdG8gdGhlIHNhbWUgZmlsZSAod29ya3MgYm90aCBmb3IgZ3ppcCBhbmQgYnJvdGxpKTpcbi8vIGAndXNlYCwgYGFuZG9tYCwgYW5kIGByaWN0J2Bcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XG4vLyBgLTI2VGAsIGAxOTgzYCwgYDQwcHhgLCBgNzVweGAsIGBidXNoYCwgYGphY2tgLCBgbWluZGAsIGB2ZXJ5YCwgYW5kIGB3b2xmYFxuY29uc3QgdXJsQWxwaGFiZXQgPVxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuXG5leHBvcnQgZnVuY3Rpb24gbmFub2lkKHNpemU6IG51bWJlciA9IDIxKTogc3RyaW5nIHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIC8vIEEgY29tcGFjdCBhbHRlcm5hdGl2ZSBmb3IgYGZvciAodmFyIGkgPSAwOyBpIDwgc3RlcDsgaSsrKWAuXG4gICAgbGV0IGkgPSBzaXplIHwgMFxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgICAgLy8gYHwgMGAgaXMgbW9yZSBjb21wYWN0IGFuZCBmYXN0ZXIgdGhhbiBgTWF0aC5mbG9vcigpYC5cbiAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuLy8gUmUtZXhwb3J0IG5hbm9pZCBmb3IgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnNcbmV4cG9ydCB7IG5hbm9pZCB9O1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xuICAgIENhbGw6IDAsXG4gICAgQ2xpcGJvYXJkOiAxLFxuICAgIEFwcGxpY2F0aW9uOiAyLFxuICAgIEV2ZW50czogMyxcbiAgICBDb250ZXh0TWVudTogNCxcbiAgICBEaWFsb2c6IDUsXG4gICAgV2luZG93OiA2LFxuICAgIFNjcmVlbnM6IDcsXG4gICAgU3lzdGVtOiA4LFxuICAgIEJyb3dzZXI6IDksXG4gICAgQ2FuY2VsQ2FsbDogMTAsXG59KTtcbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcblxuLyoqXG4gKiBSdW50aW1lVHJhbnNwb3J0IGRlZmluZXMgdGhlIGludGVyZmFjZSBmb3IgY3VzdG9tIElQQyB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zLlxuICogSW1wbGVtZW50IHRoaXMgaW50ZXJmYWNlIHRvIHVzZSBXZWJTb2NrZXRzLCBjdXN0b20gcHJvdG9jb2xzLCBvciBhbnkgb3RoZXJcbiAqIHRyYW5zcG9ydCBtZWNoYW5pc20gaW5zdGVhZCBvZiB0aGUgZGVmYXVsdCBIVFRQIGZldGNoLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIFJ1bnRpbWVUcmFuc3BvcnQge1xuICAgIC8qKlxuICAgICAqIFNlbmQgYSBydW50aW1lIGNhbGwgYW5kIHJldHVybiB0aGUgcmVzcG9uc2UuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gb2JqZWN0SUQgLSBUaGUgV2FpbHMgb2JqZWN0IElEICgwPUNhbGwsIDE9Q2xpcGJvYXJkLCBldGMuKVxuICAgICAqIEBwYXJhbSBtZXRob2QgLSBUaGUgbWV0aG9kIElEIHRvIGNhbGxcbiAgICAgKiBAcGFyYW0gd2luZG93TmFtZSAtIE9wdGlvbmFsIHdpbmRvdyBuYW1lXG4gICAgICogQHBhcmFtIGFyZ3MgLSBBcmd1bWVudHMgdG8gcGFzcyAod2lsbCBiZSBKU09OIHN0cmluZ2lmaWVkIGlmIHByZXNlbnQpXG4gICAgICogQHJldHVybnMgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHJlc3BvbnNlIGRhdGFcbiAgICAgKi9cbiAgICBjYWxsKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55Pjtcbn1cblxuLyoqXG4gKiBDdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uIChjYW4gYmUgc2V0IGJ5IHVzZXIpXG4gKi9cbmxldCBjdXN0b21UcmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsID0gbnVsbDtcblxuLyoqXG4gKiBTZXQgYSBjdXN0b20gdHJhbnNwb3J0IGZvciBhbGwgV2FpbHMgcnVudGltZSBjYWxscy5cbiAqIFRoaXMgYWxsb3dzIHlvdSB0byByZXBsYWNlIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0IHdpdGhcbiAqIFdlYlNvY2tldHMsIGN1c3RvbSBwcm90b2NvbHMsIG9yIGFueSBvdGhlciBtZWNoYW5pc20uXG4gKlxuICogQHBhcmFtIHRyYW5zcG9ydCAtIFlvdXIgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvblxuICpcbiAqIEBleGFtcGxlXG4gKiBgYGB0eXBlc2NyaXB0XG4gKiBpbXBvcnQgeyBzZXRUcmFuc3BvcnQgfSBmcm9tICcvd2FpbHMvcnVudGltZS5qcyc7XG4gKlxuICogY29uc3Qgd3NUcmFuc3BvcnQgPSB7XG4gKiAgIGNhbGw6IGFzeW5jIChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSA9PiB7XG4gKiAgICAgLy8gWW91ciBXZWJTb2NrZXQgaW1wbGVtZW50YXRpb25cbiAqICAgfVxuICogfTtcbiAqXG4gKiBzZXRUcmFuc3BvcnQod3NUcmFuc3BvcnQpO1xuICogYGBgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzZXRUcmFuc3BvcnQodHJhbnNwb3J0OiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCk6IHZvaWQge1xuICAgIGN1c3RvbVRyYW5zcG9ydCA9IHRyYW5zcG9ydDtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIGN1cnJlbnQgdHJhbnNwb3J0ICh1c2VmdWwgZm9yIGV4dGVuZGluZy93cmFwcGluZylcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGdldFRyYW5zcG9ydCgpOiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCB7XG4gICAgcmV0dXJuIGN1c3RvbVRyYW5zcG9ydDtcbn1cblxuLyoqXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxuICpcbiAqIEBwYXJhbSBvYmplY3QgLSBUaGUgb2JqZWN0IHRvIGludm9rZSB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICogQHJldHVybiBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZDogbnVtYmVyLCBhcmdzOiBhbnkgPSBudWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcbiAgICAvLyBVc2UgY3VzdG9tIHRyYW5zcG9ydCBpZiBhdmFpbGFibGVcbiAgICBpZiAoY3VzdG9tVHJhbnNwb3J0KSB7XG4gICAgICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQuY2FsbChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9XG5cbiAgICAvLyBEZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQudG9TdHJpbmcoKSk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kLnRvU3RyaW5nKCkpO1xuICAgIGlmIChhcmdzKSB7IHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7IH1cblxuICAgIGxldCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+ID0ge1xuICAgICAgICBbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXTogY2xpZW50SWRcbiAgICB9XG4gICAgaWYgKHdpbmRvd05hbWUpIHtcbiAgICAgICAgaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xuICAgIH1cblxuICAgIGxldCByZXNwb25zZSA9IGF3YWl0IGZldGNoKHVybCwgeyBoZWFkZXJzIH0pO1xuICAgIGlmICghcmVzcG9uc2Uub2spIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XG4gICAgfVxuXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgfVxufVxuXG4vKipcbiAqIEhlbHBlciB1dGlsaXRpZXMgZm9yIGN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zXG4gKi9cblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgbWVzc2FnZSBJRCBmb3IgdHJhbnNwb3J0IHJlcXVlc3RzLlxuICogVXNlcyB0aGUgc2FtZSBuYW5vaWQgaW1wbGVtZW50YXRpb24gYXMgdGhlIFdhaWxzIHJ1bnRpbWUuXG4gKlxuICogQHJldHVybnMgQSB1bmlxdWUgc3RyaW5nIGlkZW50aWZpZXIgKDIxIGNoYXJhY3RlcnMpXG4gKlxuICogQGV4YW1wbGVcbiAqIGBgYHR5cGVzY3JpcHRcbiAqIGltcG9ydCB7IGdlbmVyYXRlTWVzc2FnZUlEIH0gZnJvbSAnL3dhaWxzL3J1bnRpbWUuanMnO1xuICogY29uc3QgbXNnSUQgPSBnZW5lcmF0ZU1lc3NhZ2VJRCgpO1xuICogYGBgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBnZW5lcmF0ZU1lc3NhZ2VJRCgpOiBzdHJpbmcge1xuICAgIHJldHVybiBuYW5vaWQoKTtcbn1cblxuLyoqXG4gKiBCdWlsZHMgYSB0cmFuc3BvcnQgcmVxdWVzdCBtZXNzYWdlIGluIHRoZSBmb3JtYXQgZXhwZWN0ZWQgYnkgV2FpbHMuXG4gKiBIYW5kbGVzIEpTT04gc3RyaW5naWZpY2F0aW9uIG9mIGFyZ3VtZW50cyBhbmQgcHJvcGVyIGZpZWxkIG5hbWluZy5cbiAqXG4gKiBAcGFyYW0gaWQgLSBVbmlxdWUgbWVzc2FnZSBJRCAodXNlIGdlbmVyYXRlTWVzc2FnZUlEKCkpXG4gKiBAcGFyYW0gb2JqZWN0SUQgLSBXYWlscyBvYmplY3QgSUQgKDA9Q2FsbCwgMT1DbGlwYm9hcmQsIGV0Yy4gLSBzZWUgb2JqZWN0TmFtZXMpXG4gKiBAcGFyYW0gbWV0aG9kIC0gTWV0aG9kIElEIHdpdGhpbiB0aGUgb2JqZWN0XG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFNvdXJjZSB3aW5kb3cgbmFtZSAob3B0aW9uYWwpXG4gKiBAcGFyYW0gYXJncyAtIE1ldGhvZCBhcmd1bWVudHMgKHdpbGwgYmUgSlNPTiBzdHJpbmdpZmllZClcbiAqIEByZXR1cm5zIEZvcm1hdHRlZCB0cmFuc3BvcnQgcmVxdWVzdCBtZXNzYWdlXG4gKlxuICogQGV4YW1wbGVcbiAqIGBgYHR5cGVzY3JpcHRcbiAqIGltcG9ydCB7IGJ1aWxkVHJhbnNwb3J0UmVxdWVzdCwgZ2VuZXJhdGVNZXNzYWdlSUQsIG9iamVjdE5hbWVzIH0gZnJvbSAnL3dhaWxzL3J1bnRpbWUuanMnO1xuICpcbiAqIGNvbnN0IG1lc3NhZ2UgPSBidWlsZFRyYW5zcG9ydFJlcXVlc3QoXG4gKiAgICAgZ2VuZXJhdGVNZXNzYWdlSUQoKSxcbiAqICAgICBvYmplY3ROYW1lcy5DYWxsLFxuICogICAgIDAsXG4gKiAgICAgJycsXG4gKiAgICAgeyBtZXRob2ROYW1lOiAnR3JlZXQnLCBhcmdzOiBbJ1dvcmxkJ10gfVxuICogKTtcbiAqIGBgYFxuICovXG5leHBvcnQgZnVuY3Rpb24gYnVpbGRUcmFuc3BvcnRSZXF1ZXN0KFxuICAgIGlkOiBzdHJpbmcsXG4gICAgb2JqZWN0SUQ6IG51bWJlcixcbiAgICBtZXRob2Q6IG51bWJlcixcbiAgICB3aW5kb3dOYW1lOiBzdHJpbmcsXG4gICAgYXJnczogYW55XG4pOiBhbnkge1xuICAgIHJldHVybiB7XG4gICAgICAgIGlkOiBpZCxcbiAgICAgICAgdHlwZTogJ3JlcXVlc3QnLFxuICAgICAgICByZXF1ZXN0OiB7XG4gICAgICAgICAgICBvYmplY3Q6IG9iamVjdElELFxuICAgICAgICAgICAgbWV0aG9kOiBtZXRob2QsXG4gICAgICAgICAgICBhcmdzOiBhcmdzID8gSlNPTi5zdHJpbmdpZnkoYXJncykgOiB1bmRlZmluZWQsXG4gICAgICAgICAgICB3aW5kb3dOYW1lOiB3aW5kb3dOYW1lIHx8IHVuZGVmaW5lZCxcbiAgICAgICAgICAgIGNsaWVudElkOiBjbGllbnRJZFxuICAgICAgICB9XG4gICAgfTtcbn1cblxuLyoqXG4gKiBIYW5kbGVzIFdhaWxzIGNhbGxiYWNrIGludm9jYXRpb24gZm9yIGJpbmRpbmcgY2FsbHMuXG4gKiBUaGlzIGFic3RyYWN0cyB0aGUgaW50ZXJuYWwgbWVjaGFuaXNtIG9mIGhvdyBXYWlscyBwcm9jZXNzZXMgbWV0aG9kIGNhbGwgcmVzdWx0cy5cbiAqXG4gKiBGb3IgYmluZGluZyBjYWxscyAob2JqZWN0PTAsIG1ldGhvZD0wKSwgV2FpbHMgZXhwZWN0cyB0aGUgcmVzdWx0IHRvIGJlIGRlbGl2ZXJlZFxuICogdmlhIHdpbmRvdy5fd2FpbHMuY2FsbFJlc3VsdEhhbmRsZXIoKSByYXRoZXIgdGhhbiByZXNvbHZpbmcgdGhlIHRyYW5zcG9ydCBwcm9taXNlIGRpcmVjdGx5LlxuICpcbiAqIEBwYXJhbSBwZW5kaW5nIC0gUGVuZGluZyByZXF1ZXN0IG9iamVjdCB3aXRoIHN0b3JlZCBhcmdzXG4gKiBAcGFyYW0gcmVzcG9uc2VEYXRhIC0gRGVjb2RlZCByZXNwb25zZSBkYXRhIChzdHJpbmcpXG4gKiBAcGFyYW0gY29udGVudFR5cGUgLSBSZXNwb25zZSBjb250ZW50IHR5cGVcbiAqIEByZXR1cm5zIHRydWUgaWYgaGFuZGxlZCBhcyBiaW5kaW5nIGNhbGwsIGZhbHNlIG90aGVyd2lzZVxuICpcbiAqIEBleGFtcGxlXG4gKiBgYGB0eXBlc2NyaXB0XG4gKiBpbXBvcnQgeyBoYW5kbGVXYWlsc0NhbGxiYWNrIH0gZnJvbSAnL3dhaWxzL3J1bnRpbWUuanMnO1xuICpcbiAqIC8vIEluIHlvdXIgdHJhbnNwb3J0J3MgbWVzc2FnZSBoYW5kbGVyOlxuICogaWYgKGhhbmRsZVdhaWxzQ2FsbGJhY2socGVuZGluZywgcmVzcG9uc2VEYXRhLCByZXNwb25zZS5jb250ZW50VHlwZSkpIHtcbiAqICAgICBwZW5kaW5nLnJlc29sdmUoKTsgLy8gV2FpbHMgY2FsbGJhY2sgaGFuZGxlZCBpdFxuICogfSBlbHNlIHtcbiAqICAgICBwZW5kaW5nLnJlc29sdmUocmVzcG9uc2VEYXRhKTsgLy8gRGlyZWN0IHJlc29sdmUgZm9yIG5vbi1iaW5kaW5nIGNhbGxzXG4gKiB9XG4gKiBgYGBcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGhhbmRsZVdhaWxzQ2FsbGJhY2soXG4gICAgcGVuZGluZzogYW55LFxuICAgIHJlc3BvbnNlRGF0YTogc3RyaW5nLFxuICAgIGNvbnRlbnRUeXBlOiBzdHJpbmdcbik6IGJvb2xlYW4ge1xuICAgIC8vIE9ubHkgZm9yIEpTT04gcmVzcG9uc2VzIChiaW5kaW5nIGNhbGxzKVxuICAgIGlmICghcmVzcG9uc2VEYXRhIHx8ICFjb250ZW50VHlwZT8uaW5jbHVkZXMoJ2FwcGxpY2F0aW9uL2pzb24nKSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgLy8gRXh0cmFjdCBjYWxsLWlkIGZyb20gc3RvcmVkIHJlcXVlc3RcbiAgICBjb25zdCBjYWxsSWQgPSBwZW5kaW5nLnJlcXVlc3Q/LmFyZ3M/LlsnY2FsbC1pZCddO1xuXG4gICAgLy8gSW52b2tlIFdhaWxzIGNhbGxiYWNrIGhhbmRsZXIgaWYgYXZhaWxhYmxlXG4gICAgaWYgKGNhbGxJZCAmJiB3aW5kb3cuX3dhaWxzPy5jYWxsUmVzdWx0SGFuZGxlcikge1xuICAgICAgICB3aW5kb3cuX3dhaWxzLmNhbGxSZXN1bHRIYW5kbGVyKGNhbGxJZCwgcmVzcG9uc2VEYXRhLCB0cnVlKTtcbiAgICAgICAgcmV0dXJuIHRydWU7XG4gICAgfVxuXG4gICAgcmV0dXJuIGZhbHNlO1xufVxuXG4vKipcbiAqIERpc3BhdGNoZXMgYW4gZXZlbnQgdG8gdGhlIFdhaWxzIGV2ZW50IHN5c3RlbS5cbiAqIFVzZWQgYnkgY3VzdG9tIHRyYW5zcG9ydHMgdG8gZGVsaXZlciBzZXJ2ZXItc2VudCBldmVudHMgdG8gZnJvbnRlbmQgbGlzdGVuZXJzLlxuICpcbiAqIEBwYXJhbSBldmVudCAtIEV2ZW50IG9iamVjdCB3aXRoIG5hbWUsIGRhdGEsIGFuZCBvcHRpb25hbCBzZW5kZXJcbiAqXG4gKiBAZXhhbXBsZVxuICogYGBgdHlwZXNjcmlwdFxuICogaW1wb3J0IHsgZGlzcGF0Y2hXYWlsc0V2ZW50IH0gZnJvbSAnL3dhaWxzL3J1bnRpbWUuanMnO1xuICpcbiAqIC8vIEluIHlvdXIgV2ViU29ja2V0IHRyYW5zcG9ydCdzIG1lc3NhZ2UgaGFuZGxlcjpcbiAqIGlmIChtc2cudHlwZSA9PT0gJ2V2ZW50Jykge1xuICogICAgIGRpc3BhdGNoV2FpbHNFdmVudChtc2cuZXZlbnQpO1xuICogfVxuICogYGBgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQ6IGFueSk6IHZvaWQge1xuICAgIGlmICh3aW5kb3cuX3dhaWxzPy5kaXNwYXRjaFdhaWxzRXZlbnQpIHtcbiAgICAgICAgd2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnLi9uYW5vaWQuanMnO1xuXG4vLyBzZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpYWxvZ0Vycm9yQ2FsbGJhY2sgPSBkaWFsb2dFcnJvckNhbGxiYWNrO1xud2luZG93Ll93YWlscy5kaWFsb2dSZXN1bHRDYWxsYmFjayA9IGRpYWxvZ1Jlc3VsdENhbGxiYWNrO1xuXG50eXBlIFByb21pc2VSZXNvbHZlcnMgPSBPbWl0PFByb21pc2VXaXRoUmVzb2x2ZXJzPGFueT4sIFwicHJvbWlzZVwiPjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRGlhbG9nKTtcbmNvbnN0IGRpYWxvZ1Jlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xuXG4vLyBEZWZpbmUgY29uc3RhbnRzIGZyb20gdGhlIGBtZXRob2RzYCBvYmplY3QgaW4gVGl0bGUgQ2FzZVxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XG5jb25zdCBEaWFsb2dXYXJuaW5nID0gMTtcbmNvbnN0IERpYWxvZ0Vycm9yID0gMjtcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcbmNvbnN0IERpYWxvZ09wZW5GaWxlID0gNDtcbmNvbnN0IERpYWxvZ1NhdmVGaWxlID0gNTtcblxuZXhwb3J0IGludGVyZmFjZSBPcGVuRmlsZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBTYXZlRmlsZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBEZWZhdWx0IGZpbGVuYW1lIHRvIHVzZSBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIEZpbGVuYW1lPzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBNZXNzYWdlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIFRoZSB0aXRsZSBvZiB0aGUgZGlhbG9nIHdpbmRvdy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIEJ1dHRvbnM/OiBCdXR0b25bXTtcbiAgICAvKiogVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBCdXR0b24ge1xuICAgIC8qKiBUZXh0IHRoYXQgYXBwZWFycyB3aXRoaW4gdGhlIGJ1dHRvbi4gKi9cbiAgICBMYWJlbD86IHN0cmluZztcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC4gKi9cbiAgICBJc0NhbmNlbD86IGJvb2xlYW47XG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci4gKi9cbiAgICBJc0RlZmF1bHQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEZpbGVGaWx0ZXIge1xuICAgIC8qKiBEaXNwbGF5IG5hbWUgZm9yIHRoZSBmaWx0ZXIsIGl0IGNvdWxkIGJlIFwiVGV4dCBGaWxlc1wiLCBcIkltYWdlc1wiIGV0Yy4gKi9cbiAgICBEaXNwbGF5TmFtZT86IHN0cmluZztcbiAgICAvKiogUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy4gKi9cbiAgICBQYXR0ZXJuPzogc3RyaW5nO1xufVxuXG4vKipcbiAqIEhhbmRsZXMgdGhlIHJlc3VsdCBvZiBhIGRpYWxvZyByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcmVxdWVzdCB0byBoYW5kbGUgdGhlIHJlc3VsdCBmb3IuXG4gKiBAcGFyYW0gZGF0YSAtIFRoZSByZXN1bHQgZGF0YSBvZiB0aGUgcmVxdWVzdC5cbiAqIEBwYXJhbSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cbiAqL1xuZnVuY3Rpb24gZGlhbG9nUmVzdWx0Q2FsbGJhY2soaWQ6IHN0cmluZywgZGF0YTogc3RyaW5nLCBpc0pTT046IGJvb2xlYW4pOiB2b2lkIHtcbiAgICBsZXQgcmVzb2x2ZXJzID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xuICAgIGlmICghcmVzb2x2ZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyOiBhbnkpIHtcbiAgICAgICAgICAgIHJlc29sdmVycy5yZWplY3QobmV3IFR5cGVFcnJvcihcImNvdWxkIG5vdCBwYXJzZSByZXN1bHQ6IFwiICsgZXJyLm1lc3NhZ2UsIHsgY2F1c2U6IGVyciB9KSk7XG4gICAgICAgIH1cbiAgICB9IGVsc2Uge1xuICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShkYXRhKTtcbiAgICB9XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgZXJyb3IgZnJvbSBhIGRpYWxvZyByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcHJvbWlzZSBoYW5kbGVyLlxuICogQHBhcmFtIG1lc3NhZ2UgLSBBbiBlcnJvciBtZXNzYWdlLlxuICovXG5mdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkOiBzdHJpbmcsIG1lc3NhZ2U6IHN0cmluZyk6IHZvaWQge1xuICAgIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKT8ucmVqZWN0KG5ldyB3aW5kb3cuRXJyb3IobWVzc2FnZSkpO1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyBhbmQgcmVtb3ZlcyB0aGUgcmVzcG9uc2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCBmcm9tIHRoZSBkaWFsb2dSZXNwb25zZXMgbWFwLlxuICpcbiAqIEBwYXJhbSBpZCAtIFRoZSBJRCBvZiB0aGUgcmVzcG9uc2UgdG8gYmUgcmV0cmlldmVkIGFuZCByZW1vdmVkLlxuICogQHJldHVybnMgVGhlIHJlc3BvbnNlIG9iamVjdCBhc3NvY2lhdGVkIHdpdGggdGhlIGdpdmVuIElELCBpZiBhbnkuXG4gKi9cbmZ1bmN0aW9uIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkOiBzdHJpbmcpOiBQcm9taXNlUmVzb2x2ZXJzIHwgdW5kZWZpbmVkIHtcbiAgICBjb25zdCByZXNwb25zZSA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIHJldHVybiByZXNwb25zZTtcbn1cblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxuICpcbiAqIEByZXR1cm5zIEEgdW5pcXVlIElEIHRoYXQgZG9lcyBub3QgZXhpc3QgaW4gdGhlIGRpYWxvZ1Jlc3BvbnNlcyBzZXQuXG4gKi9cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKTogc3RyaW5nIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKlxuICogQHBhcmFtIHR5cGUgLSBEaWFsb2cgdHlwZS5cbiAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyBmb3IgdGhlIGRpYWxvZy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggcmVzdWx0IG9mIGRpYWxvZy5cbiAqL1xuZnVuY3Rpb24gZGlhbG9nKHR5cGU6IG51bWJlciwgb3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMgfCBPcGVuRmlsZURpYWxvZ09wdGlvbnMgfCBTYXZlRmlsZURpYWxvZ09wdGlvbnMgPSB7fSk6IFByb21pc2U8YW55PiB7XG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlLCByZWplY3QgfSk7XG4gICAgICAgIGNhbGwodHlwZSwgT2JqZWN0LmFzc2lnbih7IFwiZGlhbG9nLWlkXCI6IGlkIH0sIG9wdGlvbnMpKS5jYXRjaCgoZXJyOiBhbnkpID0+IHtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICAgICAgcmVqZWN0KGVycik7XG4gICAgICAgIH0pO1xuICAgIH0pO1xufVxuXG4vKipcbiAqIFByZXNlbnRzIGFuIGluZm8gZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIHdhcm5pbmcgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdhcm5pbmcob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYW4gZXJyb3IgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIHF1ZXN0aW9uIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSBmaWxlIHNlbGVjdGlvbiBkaWFsb2cgdG8gcGljayBvbmUgb3IgbW9yZSBmaWxlcyB0byBvcGVuLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBTZWxlY3RlZCBmaWxlIG9yIGxpc3Qgb2YgZmlsZXMsIG9yIGEgYmxhbmsgc3RyaW5nL2VtcHR5IGxpc3QgaWYgbm8gZmlsZSBoYXMgYmVlbiBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyAmIHsgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb246IHRydWUgfSk6IFByb21pc2U8c3RyaW5nW10+O1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyAmIHsgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24/OiBmYWxzZSB8IHVuZGVmaW5lZCB9KTogUHJvbWlzZTxzdHJpbmc+O1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+O1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dPcGVuRmlsZSwgb3B0aW9ucykgPz8gW107IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIGEgZmlsZSB0byBzYXZlLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBTZWxlY3RlZCBmaWxlLCBvciBhIGJsYW5rIHN0cmluZyBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9uczogU2F2ZUZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7IH1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBldmVudExpc3RlbmVycywgTGlzdGVuZXIsIGxpc3RlbmVyT2ZmIH0gZnJvbSBcIi4vbGlzdGVuZXIuanNcIjtcblxuLy8gU2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkV2ZW50cyk7XG5jb25zdCBFbWl0TWV0aG9kID0gMDtcblxuZXhwb3J0IHsgVHlwZXMgfSBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xuXG4vKipcbiAqIFRoZSB0eXBlIG9mIGhhbmRsZXJzIGZvciBhIGdpdmVuIGV2ZW50LlxuICovXG5leHBvcnQgdHlwZSBDYWxsYmFjayA9IChldjogV2FpbHNFdmVudCkgPT4gdm9pZDtcblxuLyoqXG4gKiBSZXByZXNlbnRzIGEgc3lzdGVtIGV2ZW50IG9yIGEgY3VzdG9tIGV2ZW50IGVtaXR0ZWQgdGhyb3VnaCB3YWlscy1wcm92aWRlZCBmYWNpbGl0aWVzLlxuICovXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgLyoqXG4gICAgICogVGhlIG5hbWUgb2YgdGhlIGV2ZW50LlxuICAgICAqL1xuICAgIG5hbWU6IHN0cmluZztcblxuICAgIC8qKlxuICAgICAqIE9wdGlvbmFsIGRhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBlbWl0dGVkIGV2ZW50LlxuICAgICAqL1xuICAgIGRhdGE6IGFueTtcblxuICAgIC8qKlxuICAgICAqIE5hbWUgb2YgdGhlIG9yaWdpbmF0aW5nIHdpbmRvdy4gT21pdHRlZCBmb3IgYXBwbGljYXRpb24gZXZlbnRzLlxuICAgICAqIFdpbGwgYmUgb3ZlcnJpZGRlbiBpZiBzZXQgbWFudWFsbHkuXG4gICAgICovXG4gICAgc2VuZGVyPzogc3RyaW5nO1xuXG4gICAgY29uc3RydWN0b3IobmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQ6IGFueSkge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCB3YWlsc0V2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnQubmFtZSwgZXZlbnQuZGF0YSk7XG4gICAgaWYgKCdzZW5kZXInIGluIGV2ZW50KSB7XG4gICAgICAgIHdhaWxzRXZlbnQuc2VuZGVyID0gZXZlbnQuc2VuZGVyO1xuICAgIH1cblxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4gIWxpc3RlbmVyLmRpc3BhdGNoKHdhaWxzRXZlbnQpKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxuICpcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxuICogQHBhcmFtIG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiBDYWxsYmFjaywgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIHNwZWNpZmllZCBldmVudCBvY2N1cnMuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT24oZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiBDYWxsYmFjayk6ICgpID0+IHZvaWQge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IENhbGxiYWNrKTogKCkgPT4gdm9pZCB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZXMgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnRzIHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKC4uLmV2ZW50TmFtZXM6IFtzdHJpbmcsIC4uLnN0cmluZ1tdXSk6IHZvaWQge1xuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xufVxuXG4vKipcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpOiB2b2lkIHtcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xufVxuXG4vKipcbiAqIEVtaXRzIGFuIGV2ZW50IHVzaW5nIHRoZSBuYW1lIGFuZCBkYXRhLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgZnVsZmlsbGVkIG9uY2UgdGhlIGV2ZW50IGhhcyBiZWVuIGVtaXR0ZWQuXG4gKiBAcGFyYW0gbmFtZSAtIHRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0LlxuICogQHBhcmFtIGRhdGEgLSB0aGUgZGF0YSB0byBiZSBzZW50IHdpdGggdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChuYW1lOiBzdHJpbmcsIGRhdGE/OiBhbnkpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBsZXQgZXZlbnROYW1lOiBzdHJpbmc7XG4gICAgbGV0IGV2ZW50RGF0YTogYW55O1xuXG4gICAgaWYgKHR5cGVvZiBuYW1lID09PSAnb2JqZWN0JyAmJiBuYW1lICE9PSBudWxsICYmICduYW1lJyBpbiBuYW1lICYmICdkYXRhJyBpbiBuYW1lKSB7XG4gICAgICAgIC8vIElmIG5hbWUgaXMgYW4gb2JqZWN0IHdpdGggYSBuYW1lIHByb3BlcnR5LCB1c2UgaXQgZGlyZWN0bHlcbiAgICAgICAgZXZlbnROYW1lID0gbmFtZVsnbmFtZSddO1xuICAgICAgICBldmVudERhdGEgPSBuYW1lWydkYXRhJ107XG4gICAgfSBlbHNlIHtcbiAgICAgICAgLy8gT3RoZXJ3aXNlIHVzZSB0aGUgc3RhbmRhcmQgcGFyYW1ldGVyc1xuICAgICAgICBldmVudE5hbWUgPSBuYW1lIGFzIHN0cmluZztcbiAgICAgICAgZXZlbnREYXRhID0gZGF0YTtcbiAgICB9XG5cbiAgICByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCB7IG5hbWU6IGV2ZW50TmFtZSwgZGF0YTogZXZlbnREYXRhIH0pO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFRoZSBmb2xsb3dpbmcgdXRpbGl0aWVzIGhhdmUgYmVlbiBmYWN0b3JlZCBvdXQgb2YgLi9ldmVudHMudHNcbi8vIGZvciB0ZXN0aW5nIHB1cnBvc2VzLlxuXG5leHBvcnQgY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwPHN0cmluZywgTGlzdGVuZXJbXT4oKTtcblxuZXhwb3J0IGNsYXNzIExpc3RlbmVyIHtcbiAgICBldmVudE5hbWU6IHN0cmluZztcbiAgICBjYWxsYmFjazogKGRhdGE6IGFueSkgPT4gdm9pZDtcbiAgICBtYXhDYWxsYmFja3M6IG51bWJlcjtcblxuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZTogc3RyaW5nLCBjYWxsYmFjazogKGRhdGE6IGFueSkgPT4gdm9pZCwgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIHRoaXMuY2FsbGJhY2sgPSBjYWxsYmFjaztcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgfVxuXG4gICAgZGlzcGF0Y2goZGF0YTogYW55KTogYm9vbGVhbiB7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICB0aGlzLmNhbGxiYWNrKGRhdGEpO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZXJyKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcbiAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyOiBMaXN0ZW5lcik6IHZvaWQge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQobGlzdGVuZXIuZXZlbnROYW1lKTtcbiAgICBpZiAoIWxpc3RlbmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUobGlzdGVuZXIuZXZlbnROYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQobGlzdGVuZXIuZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gQ3luaHlyY2h3eWQgeSBmZmVpbCBob24geW4gYXd0b21hdGlnLiBQRUlESVdDSCBcdTAwQzIgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHR9KSxcblx0TGludXg6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibGludXg6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcblx0XHRXaW5kb3dMb2FkQ2hhbmdlZDogXCJsaW51eDpXaW5kb3dMb2FkQ2hhbmdlZFwiLFxuXHR9KSxcblx0Q29tbW9uOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbk9wZW5lZFdpdGhGaWxlOiBcImNvbW1vbjpBcHBsaWNhdGlvbk9wZW5lZFdpdGhGaWxlXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcImNvbW1vbjpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRBcHBsaWNhdGlvbkxhdW5jaGVkV2l0aFVybDogXCJjb21tb246QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmxcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiY29tbW9uOlRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwiY29tbW9uOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImNvbW1vbjpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImNvbW1vbjpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcImNvbW1vbjpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcblx0XHRXaW5kb3dGb2N1czogXCJjb21tb246V2luZG93Rm9jdXNcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJjb21tb246V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJjb21tb246V2luZG93TG9zdEZvY3VzXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VG9nZ2xlRnJhbWVsZXNzOiBcImNvbW1vbjpXaW5kb3dUb2dnbGVGcmFtZWxlc3NcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcImNvbW1vbjpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcblx0XHRXaW5kb3dTaG93OiBcImNvbW1vbjpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJjb21tb246V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1pvb206IFwiY29tbW9uOldpbmRvd1pvb21cIixcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwiY29tbW9uOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwiY29tbW9uOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdFdpbmRvd0Ryb3Bab25lRmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dEcm9wWm9uZUZpbGVzRHJvcHBlZFwiLFxuXHR9KSxcbn0pO1xuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIExvZ3MgYSBtZXNzYWdlIHRvIHRoZSBjb25zb2xlIHdpdGggY3VzdG9tIGZvcm1hdHRpbmcuXG4gKlxuICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlOiBhbnkpIHtcbiAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICBjb25zb2xlLmxvZyhcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNhYTAwMDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAzcHggMHB4IDBweCAzcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nLFxuICAgICAgICAnYmFja2dyb3VuZDogIzAwOTkwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDBweCAzcHggM3B4IDBweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbSdcbiAgICApO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSB3ZWJ2aWV3IHN1cHBvcnRzIHRoZSB7QGxpbmsgTW91c2VFdmVudCNidXR0b25zfSBwcm9wZXJ0eS5cbiAqIExvb2tpbmcgYXQgeW91IG1hY09TIEhpZ2ggU2llcnJhIVxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuVHJhY2tCdXR0b25zKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAobmV3IE1vdXNlRXZlbnQoJ21vdXNlZG93bicpKS5idXR0b25zID09PSAwO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXG4gKiAoc2VlIGh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnL2VuLVVTL2RvY3MvV2ViL0FQSS9FdmVudFRhcmdldC9hZGRFdmVudExpc3RlbmVyI3NpZ25hbCkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjYW5BYm9ydExpc3RlbmVycygpIHtcbiAgICBpZiAoIUV2ZW50VGFyZ2V0IHx8ICFBYm9ydFNpZ25hbCB8fCAhQWJvcnRDb250cm9sbGVyKVxuICAgICAgICByZXR1cm4gZmFsc2U7XG5cbiAgICBsZXQgcmVzdWx0ID0gdHJ1ZTtcblxuICAgIGNvbnN0IHRhcmdldCA9IG5ldyBFdmVudFRhcmdldCgpO1xuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgdGFyZ2V0LmFkZEV2ZW50TGlzdGVuZXIoJ3Rlc3QnLCAoKSA9PiB7IHJlc3VsdCA9IGZhbHNlOyB9LCB7IHNpZ25hbDogY29udHJvbGxlci5zaWduYWwgfSk7XG4gICAgY29udHJvbGxlci5hYm9ydCgpO1xuICAgIHRhcmdldC5kaXNwYXRjaEV2ZW50KG5ldyBDdXN0b21FdmVudCgndGVzdCcpKTtcblxuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogUmVzb2x2ZXMgdGhlIGNsb3Nlc3QgSFRNTEVsZW1lbnQgYW5jZXN0b3Igb2YgYW4gZXZlbnQncyB0YXJnZXQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBldmVudFRhcmdldChldmVudDogRXZlbnQpOiBIVE1MRWxlbWVudCB7XG4gICAgaWYgKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSB7XG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQ7XG4gICAgfSBlbHNlIGlmICghKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSAmJiBldmVudC50YXJnZXQgaW5zdGFuY2VvZiBOb2RlKSB7XG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQucGFyZW50RWxlbWVudCA/PyBkb2N1bWVudC5ib2R5O1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiBkb2N1bWVudC5ib2R5O1xuICAgIH1cbn1cblxuLyoqKlxuIFRoaXMgdGVjaG5pcXVlIGZvciBwcm9wZXIgbG9hZCBkZXRlY3Rpb24gaXMgdGFrZW4gZnJvbSBIVE1YOlxuXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcblxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXG5cbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcblxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxuXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cbiBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlIGRpc3RyaWJ1dGlvbi5cblxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXG4gQU5EIEFOWSBFWFBSRVNTIE9SIElNUExJRUQgV0FSUkFOVElFUywgSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFRIRVxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcbiBGT1IgQU5ZIERJUkVDVCwgSU5ESVJFQ1QsIElOQ0lERU5UQUwsIFNQRUNJQUwsIEVYRU1QTEFSWSwgT1IgQ09OU0VRVUVOVElBTFxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcbiBDQVVTRUQgQU5EIE9OIEFOWSBUSEVPUlkgT0YgTElBQklMSVRZLCBXSEVUSEVSIElOIENPTlRSQUNULCBTVFJJQ1QgTElBQklMSVRZLFxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cblxuICoqKi9cblxubGV0IGlzUmVhZHkgPSBmYWxzZTtcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCAoKSA9PiB7IGlzUmVhZHkgPSB0cnVlIH0pO1xuXG5leHBvcnQgZnVuY3Rpb24gd2hlblJlYWR5KGNhbGxiYWNrOiAoKSA9PiB2b2lkKSB7XG4gICAgaWYgKGlzUmVhZHkgfHwgZG9jdW1lbnQucmVhZHlTdGF0ZSA9PT0gJ2NvbXBsZXRlJykge1xuICAgICAgICBjYWxsYmFjaygpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCBjYWxsYmFjayk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgdHlwZSB7IFNjcmVlbiB9IGZyb20gXCIuL3NjcmVlbnMuanNcIjtcblxuLy8gTkVXOiBEcm9wem9uZSBjb25zdGFudHNcbmNvbnN0IERST1BaT05FX0FUVFJJQlVURSA9ICdkYXRhLXdhaWxzLWRyb3B6b25lJztcbmNvbnN0IERST1BaT05FX0hPVkVSX0NMQVNTID0gJ3dhaWxzLWRyb3B6b25lLWhvdmVyJzsgLy8gVXNlciBjYW4gc3R5bGUgdGhpcyBjbGFzc1xubGV0IGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmU6IEVsZW1lbnQgfCBudWxsID0gbnVsbDtcblxuY29uc3QgUG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgICAgID0gMDtcbmNvbnN0IENlbnRlck1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDE7XG5jb25zdCBDbG9zZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSAyO1xuY29uc3QgRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgID0gMztcbmNvbnN0IEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgICA9IDQ7XG5jb25zdCBGb2N1c01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1O1xuY29uc3QgRm9yY2VSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgID0gNjtcbmNvbnN0IEZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICA9IDc7XG5jb25zdCBHZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA4O1xuY29uc3QgR2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gOTtcbmNvbnN0IEhlaWdodE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDEwO1xuY29uc3QgSGlkZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTE7XG5jb25zdCBJc0ZvY3VzZWRNZXRob2QgICAgICAgICAgICAgICAgICAgPSAxMjtcbmNvbnN0IElzRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDEzO1xuY29uc3QgSXNNYXhpbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTQ7XG5jb25zdCBJc01pbmltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNTtcbmNvbnN0IE1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE2O1xuY29uc3QgTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTc7XG5jb25zdCBOYW1lTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxODtcbmNvbnN0IE9wZW5EZXZUb29sc01ldGhvZCAgICAgICAgICAgICAgICA9IDE5O1xuY29uc3QgUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgID0gMjA7XG5jb25zdCBSZWxvYWRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAyMTtcbmNvbnN0IFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDIyO1xuY29uc3QgUmVzdG9yZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMjM7XG5jb25zdCBTZXRQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgPSAyNDtcbmNvbnN0IFNldEFsd2F5c09uVG9wTWV0aG9kICAgICAgICAgICAgICA9IDI1O1xuY29uc3QgU2V0QmFja2dyb3VuZENvbG91ck1ldGhvZCAgICAgICAgID0gMjY7XG5jb25zdCBTZXRGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgICAgPSAyNztcbmNvbnN0IFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kICA9IDI4O1xuY29uc3QgU2V0TWF4U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMjk7XG5jb25zdCBTZXRNaW5TaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAzMDtcbmNvbnN0IFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICA9IDMxO1xuY29uc3QgU2V0UmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgID0gMzI7XG5jb25zdCBTZXRTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzMztcbmNvbnN0IFNldFRpdGxlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDM0O1xuY29uc3QgU2V0Wm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzU7XG5jb25zdCBTaG93TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNjtcbmNvbnN0IFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM3O1xuY29uc3QgVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgID0gMzg7XG5jb25zdCBUb2dnbGVNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgPSAzOTtcbmNvbnN0IFRvZ2dsZUZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICA9IDQwOyBcbmNvbnN0IFVuRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDQxO1xuY29uc3QgVW5NYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDI7XG5jb25zdCBVbk1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MztcbmNvbnN0IFdpZHRoTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDQ0O1xuY29uc3QgWm9vbU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gNDU7XG5jb25zdCBab29tSW5NZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSA0NjtcbmNvbnN0IFpvb21PdXRNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDQ3O1xuY29uc3QgWm9vbVJlc2V0TWV0aG9kICAgICAgICAgICAgICAgICAgID0gNDg7XG5jb25zdCBTbmFwQXNzaXN0TWV0aG9kICAgICAgICAgICAgICAgICAgPSA0OTtcbmNvbnN0IFdpbmRvd0Ryb3Bab25lRHJvcHBlZCAgICAgICAgICAgICA9IDUwO1xuXG5mdW5jdGlvbiBnZXREcm9wem9uZUVsZW1lbnQoZWxlbWVudDogRWxlbWVudCB8IG51bGwpOiBFbGVtZW50IHwgbnVsbCB7XG4gICAgaWYgKCFlbGVtZW50KSB7XG4gICAgICAgIHJldHVybiBudWxsO1xuICAgIH1cbiAgICAvLyBBbGxvdyBkcm9wem9uZSBhdHRyaWJ1dGUgdG8gYmUgb24gdGhlIGVsZW1lbnQgaXRzZWxmIG9yIGFueSBwYXJlbnRcbiAgICByZXR1cm4gZWxlbWVudC5jbG9zZXN0KGBbJHtEUk9QWk9ORV9BVFRSSUJVVEV9XWApO1xufVxuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHBvc2l0aW9uIG9mIGEgd2luZG93LlxuICovXG5pbnRlcmZhY2UgUG9zaXRpb24ge1xuICAgIC8qKiBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xuICAgIHg6IG51bWJlcjtcbiAgICAvKiogVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXG4gICAgeTogbnVtYmVyO1xufVxuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXG4gKi9cbmludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuICovXG4gICAgd2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LiAqL1xuICAgIGhlaWdodDogbnVtYmVyO1xufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY2FsbGVyU3ltID0gU3ltYm9sKFwiY2FsbGVyXCIpO1xuXG5jbGFzcyBXaW5kb3cge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIHByaXZhdGUgW2NhbGxlclN5bV06IChtZXNzYWdlOiBudW1iZXIsIGFyZ3M/OiBhbnkpID0+IFByb21pc2U8YW55PjtcblxuICAgIC8qKlxuICAgICAqIEluaXRpYWxpc2VzIGEgd2luZG93IG9iamVjdCB3aXRoIHRoZSBzcGVjaWZpZWQgbmFtZS5cbiAgICAgKlxuICAgICAqIEBwcml2YXRlXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgdGFyZ2V0IHdpbmRvdy5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcgPSAnJykge1xuICAgICAgICB0aGlzW2NhbGxlclN5bV0gPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcblxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xuICAgICAgICAgICAgaWYgKFxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXG4gICAgICAgICAgICApIHtcbiAgICAgICAgICAgICAgICAodGhpcyBhcyBhbnkpW21ldGhvZF0gPSAodGhpcyBhcyBhbnkpW21ldGhvZF0uYmluZCh0aGlzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gZ2V0LlxuICAgICAqIEByZXR1cm5zIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXG4gICAgICovXG4gICAgR2V0KG5hbWU6IHN0cmluZyk6IFdpbmRvdyB7XG4gICAgICAgIHJldHVybiBuZXcgV2luZG93KG5hbWUpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2VudGVycyB0aGUgd2luZG93IG9uIHRoZSBzY3JlZW4uXG4gICAgICovXG4gICAgQ2VudGVyKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENlbnRlck1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgQ2xvc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2xvc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKi9cbiAgICBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEVuYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqL1xuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBGb2N1cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb2N1c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKi9cbiAgICBGb3JjZVJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb3JjZVJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU3dpdGNoZXMgdGhlIHdpbmRvdyB0byBmdWxsc2NyZWVuIG1vZGUuXG4gICAgICovXG4gICAgRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBzY3JlZW4gdGhhdCB0aGUgd2luZG93IGlzIG9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvbi5cbiAgICAgKi9cbiAgICBHZXRTY3JlZW4oKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRTY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgem9vbSBsZXZlbC5cbiAgICAgKi9cbiAgICBHZXRab29tKCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0Wm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBIZWlnaHQoKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIZWlnaHRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEhpZGVzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIaWRlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkLlxuICAgICAqL1xuICAgIElzRm9jdXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0ZvY3VzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW4uXG4gICAgICovXG4gICAgSXNGdWxsc2NyZWVuKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtYXhpbWlzZWQuXG4gICAgICovXG4gICAgSXNNYXhpbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNYXhpbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkLlxuICAgICAqL1xuICAgIElzTWluaW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWluaW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogTWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTmFtZSgpOiBQcm9taXNlPHN0cmluZz4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE5hbWVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxuICAgICAqL1xuICAgIE9wZW5EZXZUb29scygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShPcGVuRGV2VG9vbHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUmVsYXRpdmVQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXG4gICAgICovXG4gICAgUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGUuXG4gICAgICovXG4gICAgUmVzaXphYmxlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc2l6YWJsZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVzdG9yZXMgdGhlIHdpbmRvdyB0byBpdHMgcHJldmlvdXMgc3RhdGUgaWYgaXQgd2FzIHByZXZpb3VzbHkgbWluaW1pc2VkLCBtYXhpbWlzZWQgb3IgZnVsbHNjcmVlbi5cbiAgICAgKi9cbiAgICBSZXN0b3JlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc3RvcmVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0UG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxuICAgICAqXG4gICAgICogQHBhcmFtIGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcC5cbiAgICAgKi9cbiAgICBTZXRBbHdheXNPblRvcChhbHdheXNPblRvcDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gciAtIFRoZSBkZXNpcmVkIHJlZCBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYSAtIFRoZSBkZXNpcmVkIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICovXG4gICAgU2V0QmFja2dyb3VuZENvbG91cihyOiBudW1iZXIsIGc6IG51bWJlciwgYjogbnVtYmVyLCBhOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3MuXG4gICAgICovXG4gICAgU2V0RnJhbWVsZXNzKGZyYW1lbGVzczogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWQuXG4gICAgICovXG4gICAgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0TWF4U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWluaW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRNaW5TaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWluU2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cbiAgICAgKi9cbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gdGl0bGUgLSBUaGUgZGVzaXJlZCB0aXRsZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFRpdGxlKHRpdGxlOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWwuXG4gICAgICovXG4gICAgU2V0Wm9vbSh6b29tOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNob3dNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaXplKCk6IFByb21pc2U8U2l6ZT4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNpemVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZ1bGxzY3JlZW4gYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZU1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmcmFtZWxlc3MgYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGcmFtZWxlc3MoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnJhbWVsZXNzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLW1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5NaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBXaWR0aCgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFdpZHRoTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBab29tcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFpvb20oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbUluKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21Jbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbU91dCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tT3V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tUmVzZXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbVJlc2V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUyBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXG4gICAgICogR2F0aGVycyBpbmZvcm1hdGlvbiBhYm91dCB0aGUgZHJvcCB0YXJnZXQgZWxlbWVudCBhbmQgc2VuZHMgaXQgYmFjayB0byB0aGUgR28gYmFja2VuZC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cbiAgICAgKiBAcGFyYW0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQuXG4gICAgICogQHBhcmFtIHkgLSBUaGUgeS1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxuICAgICAqL1xuICAgIEhhbmRsZVBsYXRmb3JtRmlsZURyb3AoZmlsZW5hbWVzOiBzdHJpbmdbXSwgeDogbnVtYmVyLCB5OiBudW1iZXIpOiB2b2lkIHtcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG5cbiAgICAgICAgLy8gTkVXOiBDaGVjayBpZiB0aGUgZHJvcCB0YXJnZXQgaXMgYSB2YWxpZCBkcm9wem9uZVxuICAgICAgICBjb25zdCBkcm9wem9uZVRhcmdldCA9IGdldERyb3B6b25lRWxlbWVudChlbGVtZW50KTtcblxuICAgICAgICBpZiAoIWRyb3B6b25lVGFyZ2V0KSB7XG4gICAgICAgICAgICBjb25zb2xlLmxvZyhgV2FpbHMgUnVudGltZTogRHJvcCBvbiBlbGVtZW50IChvciBubyBlbGVtZW50KSBhdCAke3h9LCR7eX0gd2hpY2ggaXMgbm90IGEgZGVzaWduYXRlZCBkcm9wem9uZS4gSWdub3JpbmcuIEVsZW1lbnQ6YCwgZWxlbWVudCk7XG4gICAgICAgICAgICAvLyBObyBuZWVkIHRvIGNhbGwgYmFja2VuZCBpZiBub3QgYSB2YWxpZCBkcm9wem9uZSB0YXJnZXRcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnNvbGUubG9nKGBXYWlscyBSdW50aW1lOiBEcm9wIG9uIGRlc2lnbmF0ZWQgZHJvcHpvbmUuIEVsZW1lbnQgYXQgKCR7eH0sICR7eX0pOmAsIGVsZW1lbnQsICdFZmZlY3RpdmUgZHJvcHpvbmU6JywgZHJvcHpvbmVUYXJnZXQpO1xuICAgICAgICBjb25zdCBlbGVtZW50RGV0YWlscyA9IHtcbiAgICAgICAgICAgIGlkOiBkcm9wem9uZVRhcmdldC5pZCxcbiAgICAgICAgICAgIGNsYXNzTGlzdDogQXJyYXkuZnJvbShkcm9wem9uZVRhcmdldC5jbGFzc0xpc3QpLFxuICAgICAgICAgICAgYXR0cmlidXRlczoge30gYXMgeyBba2V5OiBzdHJpbmddOiBzdHJpbmcgfSxcbiAgICAgICAgfTtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBkcm9wem9uZVRhcmdldC5hdHRyaWJ1dGVzLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgICAgICBjb25zdCBhdHRyID0gZHJvcHpvbmVUYXJnZXQuYXR0cmlidXRlc1tpXTtcbiAgICAgICAgICAgIGVsZW1lbnREZXRhaWxzLmF0dHJpYnV0ZXNbYXR0ci5uYW1lXSA9IGF0dHIudmFsdWU7XG4gICAgICAgIH1cblxuICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xuICAgICAgICAgICAgZmlsZW5hbWVzLFxuICAgICAgICAgICAgeCxcbiAgICAgICAgICAgIHksXG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscyxcbiAgICAgICAgfTtcblxuICAgICAgICB0aGlzW2NhbGxlclN5bV0oV2luZG93RHJvcFpvbmVEcm9wcGVkLCBwYXlsb2FkKTtcbiAgICB9XG4gIFxuICAgIC8qIFRyaWdnZXJzIFdpbmRvd3MgMTEgU25hcCBBc3Npc3QgZmVhdHVyZSAoV2luZG93cyBvbmx5KS5cbiAgICAgKiBUaGlzIGlzIGVxdWl2YWxlbnQgdG8gcHJlc3NpbmcgV2luK1ogYW5kIHNob3dzIHNuYXAgbGF5b3V0IG9wdGlvbnMuXG4gICAgICovXG4gICAgU25hcEFzc2lzdCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTbmFwQXNzaXN0TWV0aG9kKTtcbiAgICB9XG59XG5cbi8qKlxuICogVGhlIHdpbmRvdyB3aXRoaW4gd2hpY2ggdGhlIHNjcmlwdCBpcyBydW5uaW5nLlxuICovXG5jb25zdCB0aGlzV2luZG93ID0gbmV3IFdpbmRvdygnJyk7XG5cbi8vIE5FVzogR2xvYmFsIERyYWcgRXZlbnQgTGlzdGVuZXJzXG5mdW5jdGlvbiBzZXR1cEdsb2JhbERyb3B6b25lTGlzdGVuZXJzKCkge1xuICAgIGNvbnN0IGRvY0VsZW1lbnQgPSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQ7XG4gICAgbGV0IGRyYWdFbnRlckNvdW50ZXIgPSAwOyAvLyBUbyBoYW5kbGUgZHJhZ2VudGVyL2RyYWdsZWF2ZSBvbiBjaGlsZCBlbGVtZW50c1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnZW50ZXInLCAoZXZlbnQpID0+IHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgaWYgKGV2ZW50LmRhdGFUcmFuc2ZlciAmJiBldmVudC5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcbiAgICAgICAgICAgIGRyYWdFbnRlckNvdW50ZXIrKztcbiAgICAgICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xuICAgICAgICAgICAgY29uc3QgZHJvcHpvbmUgPSBnZXREcm9wem9uZUVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XG5cbiAgICAgICAgICAgIC8vIENsZWFyIHByZXZpb3VzIGhvdmVyIHJlZ2FyZGxlc3MsIHRoZW4gYXBwbHkgbmV3IGlmIHZhbGlkXG4gICAgICAgICAgICBpZiAoY3VycmVudEhvdmVyZWREcm9wem9uZSAmJiBjdXJyZW50SG92ZXJlZERyb3B6b25lICE9PSBkcm9wem9uZSkge1xuICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY2xhc3NMaXN0LnJlbW92ZShEUk9QWk9ORV9IT1ZFUl9DTEFTUyk7XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIGlmIChkcm9wem9uZSkge1xuICAgICAgICAgICAgICAgIGRyb3B6b25lLmNsYXNzTGlzdC5hZGQoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpO1xuICAgICAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xuICAgICAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgPSBkcm9wem9uZTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XG4gICAgICAgICAgICAgICAgY3VycmVudEhvdmVyZWREcm9wem9uZSA9IG51bGw7IC8vIEVuc3VyZSBpdCdzIGNsZWFyZWQgaWYgbm8gZHJvcHpvbmUgZm91bmRcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCAoZXZlbnQpID0+IHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gTmVjZXNzYXJ5IHRvIGFsbG93IGRyb3BcbiAgICAgICAgaWYgKGV2ZW50LmRhdGFUcmFuc2ZlciAmJiBldmVudC5kYXRhVHJhbnNmZXIudHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcbiAgICAgICAgICAgIC8vIE5vIG5lZWQgdG8gcXVlcnkgZWxlbWVudEZyb21Qb2ludCBhZ2FpbiBpZiBhbHJlYWR5IGhhbmRsZWQgYnkgZHJhZ2VudGVyIGNvcnJlY3RseVxuICAgICAgICAgICAgLy8gSnVzdCBlbnN1cmUgZHJvcEVmZmVjdCBpcyBjb250aW51b3VzbHkgc2V0IGJhc2VkIG9uIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmVcbiAgICAgICAgICAgIGlmIChjdXJyZW50SG92ZXJlZERyb3B6b25lKSB7XG4gICAgICAgICAgICAgICAgIC8vIFJlLWFwcGx5IGNsYXNzIGp1c3QgaW4gY2FzZSBpdCB3YXMgcmVtb3ZlZCBieSBzb21lIG90aGVyIEpTXG4gICAgICAgICAgICAgICAgaWYoIWN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY2xhc3NMaXN0LmNvbnRhaW5zKERST1BaT05FX0hPVkVSX0NMQVNTKSkge1xuICAgICAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lLmNsYXNzTGlzdC5hZGQoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdjb3B5JztcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdsZWF2ZScsIChldmVudCkgPT4ge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyICYmIGV2ZW50LmRhdGFUcmFuc2Zlci50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlci0tO1xuICAgICAgICAgICAgLy8gT25seSByZW1vdmUgaG92ZXIgaWYgZHJhZyB0cnVseSBsZWZ0IHRoZSB3aW5kb3cgb3IgdGhlIGxhc3QgZHJvcHpvbmVcbiAgICAgICAgICAgIGlmIChkcmFnRW50ZXJDb3VudGVyID09PSAwIHx8IGV2ZW50LnJlbGF0ZWRUYXJnZXQgPT09IG51bGwgfHwgKGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUgJiYgIWN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY29udGFpbnMoZXZlbnQucmVsYXRlZFRhcmdldCBhcyBOb2RlKSkpIHtcbiAgICAgICAgICAgICAgICBpZiAoY3VycmVudEhvdmVyZWREcm9wem9uZSkge1xuICAgICAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lLmNsYXNzTGlzdC5yZW1vdmUoRFJPUFpPTkVfSE9WRVJfQ0xBU1MpO1xuICAgICAgICAgICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lID0gbnVsbDtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7IC8vIFJlc2V0IGNvdW50ZXIgaWYgaXQgd2VudCBuZWdhdGl2ZSBvciBsZWZ0IHdpbmRvd1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIFByZXZlbnQgZGVmYXVsdCBicm93c2VyIGZpbGUgaGFuZGxpbmdcbiAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7IC8vIFJlc2V0IGNvdW50ZXJcbiAgICAgICAgaWYgKGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUpIHtcbiAgICAgICAgICAgIGN1cnJlbnRIb3ZlcmVkRHJvcHpvbmUuY2xhc3NMaXN0LnJlbW92ZShEUk9QWk9ORV9IT1ZFUl9DTEFTUyk7XG4gICAgICAgICAgICBjdXJyZW50SG92ZXJlZERyb3B6b25lID0gbnVsbDtcbiAgICAgICAgfVxuICAgICAgICAvLyBUaGUgYWN0dWFsIGRyb3AgcHJvY2Vzc2luZyBpcyBpbml0aWF0ZWQgYnkgdGhlIG5hdGl2ZSBzaWRlIGNhbGxpbmcgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcFxuICAgICAgICAvLyBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wIHdpbGwgdGhlbiBjaGVjayBpZiB0aGUgZHJvcCB3YXMgb24gYSB2YWxpZCB6b25lLlxuICAgIH0sIGZhbHNlKTtcbn1cblxuLy8gSW5pdGlhbGl6ZSBsaXN0ZW5lcnMgd2hlbiB0aGUgc2NyaXB0IGxvYWRzXG5pZiAodHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgZG9jdW1lbnQgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICBzZXR1cEdsb2JhbERyb3B6b25lTGlzdGVuZXJzKCk7XG59XG5cbmV4cG9ydCBkZWZhdWx0IHRoaXNXaW5kb3c7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XG5cbi8vIE5PVEU6IHRoZSBmb2xsb3dpbmcgbWV0aG9kcyBNVVNUIGJlIGltcG9ydGVkIGV4cGxpY2l0bHkgYmVjYXVzZSBvZiBob3cgZXNidWlsZCBpbmplY3Rpb24gd29ya3NcbmltcG9ydCB7IEVuYWJsZSBhcyBFbmFibGVXTUwgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XG5pbXBvcnQgeyBkZWJ1Z0xvZyB9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy91dGlsc1wiO1xuXG53aW5kb3cud2FpbHMgPSBSdW50aW1lO1xuRW5hYmxlV01MKCk7XG5cbmlmIChERUJVRykge1xuICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIilcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlN5c3RlbSk7XG5cbmNvbnN0IFN5c3RlbUlzRGFya01vZGUgPSAwO1xuY29uc3QgU3lzdGVtRW52aXJvbm1lbnQgPSAxO1xuY29uc3QgQXBwbGljYXRpb25GaWxlc0Ryb3BwZWRXaXRoQ29udGV4dCA9IDEwMDsgLy8gTmV3IG1ldGhvZCBJRCBmb3IgZW5yaWNoZWQgZHJvcCBldmVudFxuXG5jb25zdCBfaW52b2tlID0gKGZ1bmN0aW9uICgpIHtcbiAgICB0cnkge1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2UpIHtcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcpO1xuICAgICAgICB9IGVsc2UgaWYgKCh3aW5kb3cgYXMgYW55KS53ZWJraXQ/Lm1lc3NhZ2VIYW5kbGVycz8uWydleHRlcm5hbCddPy5wb3N0TWVzc2FnZSkge1xuICAgICAgICAgICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddLnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10pO1xuICAgICAgICB9XG4gICAgfSBjYXRjaChlKSB7fVxuXG4gICAgY29uc29sZS53YXJuKCdcXG4lY1x1MjZBMFx1RkUwRiBCcm93c2VyIEVudmlyb25tZW50IERldGVjdGVkICVjXFxuXFxuJWNPbmx5IFVJIHByZXZpZXdzIGFyZSBhdmFpbGFibGUgaW4gdGhlIGJyb3dzZXIuIEZvciBmdWxsIGZ1bmN0aW9uYWxpdHksIHBsZWFzZSBydW4gdGhlIGFwcGxpY2F0aW9uIGluIGRlc2t0b3AgbW9kZS5cXG5Nb3JlIGluZm9ybWF0aW9uIGF0OiBodHRwczovL3YzLndhaWxzLmlvL2xlYXJuL2J1aWxkLyN1c2luZy1hLWJyb3dzZXItZm9yLWRldmVsb3BtZW50XFxuJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNmZmZmZmY7IGNvbG9yOiAjMDAwMDAwOyBmb250LXdlaWdodDogYm9sZDsgcGFkZGluZzogNHB4IDhweDsgYm9yZGVyLXJhZGl1czogNHB4OyBib3JkZXI6IDJweCBzb2xpZCAjMDAwMDAwOycsXG4gICAgICAgICdiYWNrZ3JvdW5kOiB0cmFuc3BhcmVudDsnLFxuICAgICAgICAnY29sb3I6ICNmZmZmZmY7IGZvbnQtc3R5bGU6IGl0YWxpYzsgZm9udC13ZWlnaHQ6IGJvbGQ7Jyk7XG4gICAgcmV0dXJuIG51bGw7XG59KSgpO1xuXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZzogYW55KTogdm9pZCB7XG4gICAgX2ludm9rZT8uKG1zZyk7XG59XG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSBzeXN0ZW0gZGFyayBtb2RlIHN0YXR1cy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtSXNEYXJrTW9kZSk7XG59XG5cbi8qKlxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIENhcGFiaWxpdGllcygpOiBQcm9taXNlPFJlY29yZDxzdHJpbmcsIGFueT4+IHtcbiAgICBsZXQgcmVzcG9uc2UgPSBhd2FpdCBmZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIik7XG4gICAgaWYgKHJlc3BvbnNlLm9rKSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiY291bGQgbm90IGZldGNoIGNhcGFiaWxpdGllczogXCIgKyByZXNwb25zZS5zdGF0dXNUZXh0KTtcbiAgICB9XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgT1NJbmZvIHtcbiAgICAvKiogVGhlIGJyYW5kaW5nIG9mIHRoZSBPUy4gKi9cbiAgICBCcmFuZGluZzogc3RyaW5nO1xuICAgIC8qKiBUaGUgSUQgb2YgdGhlIE9TLiAqL1xuICAgIElEOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBuYW1lIG9mIHRoZSBPUy4gKi9cbiAgICBOYW1lOiBzdHJpbmc7XG4gICAgLyoqIFRoZSB2ZXJzaW9uIG9mIHRoZSBPUy4gKi9cbiAgICBWZXJzaW9uOiBzdHJpbmc7XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgRW52aXJvbm1lbnRJbmZvIHtcbiAgICAvKiogVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLiAqL1xuICAgIEFyY2g6IHN0cmluZztcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYXBwbGljYXRpb24gaXMgcnVubmluZyBpbiBkZWJ1ZyBtb2RlLCBvdGhlcndpc2UgZmFsc2UuICovXG4gICAgRGVidWc6IGJvb2xlYW47XG4gICAgLyoqIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS4gKi9cbiAgICBPUzogc3RyaW5nO1xuICAgIC8qKiBEZXRhaWxzIG9mIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xuICAgIE9TSW5mbzogT1NJbmZvO1xuICAgIC8qKiBBZGRpdGlvbmFsIHBsYXRmb3JtIGluZm9ybWF0aW9uLiAqL1xuICAgIFBsYXRmb3JtSW5mbzogUmVjb3JkPHN0cmluZywgYW55Pjtcbn1cblxuLyoqXG4gKiBSZXRyaWV2ZXMgZW52aXJvbm1lbnQgZGV0YWlscy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyBPUyBhbmQgc3lzdGVtIGFyY2hpdGVjdHVyZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVudmlyb25tZW50KCk6IFByb21pc2U8RW52aXJvbm1lbnRJbmZvPiB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtRW52aXJvbm1lbnQpO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MuXG4gKlxuICogQHJldHVybiBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MsIG90aGVyd2lzZSBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzV2luZG93cygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJ3aW5kb3dzXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXG4gKlxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuT1MgPT09IFwibGludXhcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIG1hY09TLCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gd2luZG93Ll93YWlscy5lbnZpcm9ubWVudC5PUyA9PT0gXCJkYXJ3aW5cIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FNRDY0KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYW1kNjRcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cbiAqXG4gKiBAcmV0dXJucyBSZXR1cm5zIHRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZSwgb3RoZXJ3aXNlIHJldHVybnMgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50LkFyY2ggPT09IFwiYXJtNjRcIjtcbn1cblxuLyoqXG4gKiBSZXBvcnRzIHdoZXRoZXIgdGhlIGFwcCBpcyBiZWluZyBydW4gaW4gZGVidWcgbW9kZS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RlYnVnKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBCb29sZWFuKHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQuRGVidWcpO1xufVxuXG4vKipcbiAqIEhhbmRsZXMgZmlsZSBkcm9wcyBvcmlnaW5hdGluZyBmcm9tIHBsYXRmb3JtLXNwZWNpZmljIGNvZGUgKGUuZy4sIG1hY09TIG5hdGl2ZSBkcmFnLWFuZC1kcm9wKS5cbiAqIEdhdGhlcnMgaW5mb3JtYXRpb24gYWJvdXQgdGhlIGRyb3AgdGFyZ2V0IGVsZW1lbnQgYW5kIHNlbmRzIGl0IGJhY2sgdG8gdGhlIEdvIGJhY2tlbmQuXG4gKlxuICogQHBhcmFtIGZpbGVuYW1lcyAtIEFuIGFycmF5IG9mIGZpbGUgcGF0aHMgKHN0cmluZ3MpIHRoYXQgd2VyZSBkcm9wcGVkLlxuICogQHBhcmFtIHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxuICogQHBhcmFtIHkgLSBUaGUgeS1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgIGNvbnN0IGVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHgsIHkpO1xuICAgIGNvbnN0IGVsZW1lbnRJZCA9IGVsZW1lbnQgPyBlbGVtZW50LmlkIDogJyc7XG4gICAgY29uc3QgY2xhc3NMaXN0ID0gZWxlbWVudCA/IEFycmF5LmZyb20oZWxlbWVudC5jbGFzc0xpc3QpIDogW107XG5cbiAgICBjb25zdCBwYXlsb2FkID0ge1xuICAgICAgICBmaWxlbmFtZXMsXG4gICAgICAgIHgsXG4gICAgICAgIHksXG4gICAgICAgIGVsZW1lbnRJZCxcbiAgICAgICAgY2xhc3NMaXN0LFxuICAgIH07XG5cbiAgICBjYWxsKEFwcGxpY2F0aW9uRmlsZXNEcm9wcGVkV2l0aENvbnRleHQsIHBheWxvYWQpXG4gICAgICAgIC50aGVuKCgpID0+IHtcbiAgICAgICAgICAgIC8vIE9wdGlvbmFsOiBMb2cgc3VjY2VzcyBvciBoYW5kbGUgaWYgbmVlZGVkXG4gICAgICAgICAgICBjb25zb2xlLmxvZyhcIlBsYXRmb3JtIGZpbGUgZHJvcCBwcm9jZXNzZWQgYW5kIHNlbnQgdG8gR28uXCIpO1xuICAgICAgICB9KVxuICAgICAgICAuY2F0Y2goZXJyID0+IHtcbiAgICAgICAgICAgIC8vIE9wdGlvbmFsOiBMb2cgZXJyb3JcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoXCJFcnJvciBzZW5kaW5nIHBsYXRmb3JtIGZpbGUgZHJvcCB0byBHbzpcIiwgZXJyKTtcbiAgICAgICAgfSk7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBJc0RlYnVnIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5cbi8vIHNldHVwXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5Db250ZXh0TWVudSk7XG5cbmNvbnN0IENvbnRleHRNZW51T3BlbiA9IDA7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZDogc3RyaW5nLCB4OiBudW1iZXIsIHk6IG51bWJlciwgZGF0YTogYW55KTogdm9pZCB7XG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XG59XG5cbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudDogTW91c2VFdmVudCkge1xuICAgIGNvbnN0IHRhcmdldCA9IGV2ZW50VGFyZ2V0KGV2ZW50KTtcblxuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XG4gICAgY29uc3QgY3VzdG9tQ29udGV4dE1lbnUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudVwiKS50cmltKCk7XG5cbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgY29uc3QgZGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZGF0YSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCwgdGFyZ2V0KTtcbiAgICB9XG59XG5cblxuLypcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG5cblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXG4qL1xuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudDogTW91c2VFdmVudCwgdGFyZ2V0OiBIVE1MRWxlbWVudCkge1xuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxuICAgIGlmIChJc0RlYnVnKCkpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICBzd2l0Y2ggKHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCkpIHtcbiAgICAgICAgY2FzZSAnc2hvdyc6XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGNhc2UgJ2hpZGUnOlxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxuICAgIGlmICh0YXJnZXQuaXNDb250ZW50RWRpdGFibGUpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWRcbiAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XG4gICAgY29uc3QgaGFzU2VsZWN0aW9uID0gc2VsZWN0aW9uICYmIHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDA7XG4gICAgaWYgKGhhc1NlbGVjdGlvbikge1xuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcbiAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XG4gICAgICAgICAgICBjb25zdCByZWN0cyA9IHJhbmdlLmdldENsaWVudFJlY3RzKCk7XG4gICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XG4gICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xuICAgICAgICAgICAgICAgIGlmIChkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHJlY3QubGVmdCwgcmVjdC50b3ApID09PSB0YXJnZXQpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8vIENoZWNrIGlmIHRhZyBpcyBpbnB1dCBvciB0ZXh0YXJlYS5cbiAgICBpZiAodGFyZ2V0IGluc3RhbmNlb2YgSFRNTElucHV0RWxlbWVudCB8fCB0YXJnZXQgaW5zdGFuY2VvZiBIVE1MVGV4dEFyZWFFbGVtZW50KSB7XG4gICAgICAgIGlmIChoYXNTZWxlY3Rpb24gfHwgKCF0YXJnZXQucmVhZE9ubHkgJiYgIXRhcmdldC5kaXNhYmxlZCkpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIFJldHJpZXZlcyB0aGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5IGZyb20gdGhlIGZsYWcgbWFwLlxuICpcbiAqIEBwYXJhbSBrZXkgLSBUaGUga2V5IHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBmb3IuXG4gKiBAcmV0dXJuIFRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleTogc3RyaW5nKTogYW55IHtcbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gd2luZG93Ll93YWlscy5mbGFnc1trZXldO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiVW5hYmxlIHRvIHJldHJpZXZlIGZsYWcgJ1wiICsga2V5ICsgXCInOiBcIiArIGUsIHsgY2F1c2U6IGUgfSk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBpbnZva2UsIElzV2luZG93cyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xuaW1wb3J0IHsgR2V0RmxhZyB9IGZyb20gXCIuL2ZsYWdzLmpzXCI7XG5pbXBvcnQgeyBjYW5UcmFja0J1dHRvbnMsIGV2ZW50VGFyZ2V0IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcblxuLy8gU2V0dXBcbmxldCBjYW5EcmFnID0gZmFsc2U7XG5sZXQgZHJhZ2dpbmcgPSBmYWxzZTtcblxubGV0IHJlc2l6YWJsZSA9IGZhbHNlO1xubGV0IGNhblJlc2l6ZSA9IGZhbHNlO1xubGV0IHJlc2l6aW5nID0gZmFsc2U7XG5sZXQgcmVzaXplRWRnZTogc3RyaW5nID0gXCJcIjtcbmxldCBkZWZhdWx0Q3Vyc29yID0gXCJhdXRvXCI7XG5cbmxldCBidXR0b25zID0gMDtcbmNvbnN0IGJ1dHRvbnNUcmFja2VkID0gY2FuVHJhY2tCdXR0b25zKCk7XG5cbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5zZXRSZXNpemFibGUgPSAodmFsdWU6IGJvb2xlYW4pOiB2b2lkID0+IHtcbiAgICByZXNpemFibGUgPSB2YWx1ZTtcbiAgICBpZiAoIXJlc2l6YWJsZSkge1xuICAgICAgICAvLyBTdG9wIHJlc2l6aW5nIGlmIGluIHByb2dyZXNzLlxuICAgICAgICBjYW5SZXNpemUgPSByZXNpemluZyA9IGZhbHNlO1xuICAgICAgICBzZXRSZXNpemUoKTtcbiAgICB9XG59O1xuXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xuZm9yIChjb25zdCBldiBvZiBbJ2NsaWNrJywgJ2NvbnRleHRtZW51JywgJ2RibGNsaWNrJ10pIHtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihldiwgc3VwcHJlc3NFdmVudCwgeyBjYXB0dXJlOiB0cnVlIH0pO1xufVxuXG5mdW5jdGlvbiBzdXBwcmVzc0V2ZW50KGV2ZW50OiBFdmVudCkge1xuICAgIC8vIFN1cHByZXNzIGNsaWNrIGV2ZW50cyB3aGlsZSByZXNpemluZyBvciBkcmFnZ2luZy5cbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn1cblxuLy8gVXNlIGNvbnN0YW50cyB0byBhdm9pZCBjb21wYXJpbmcgc3RyaW5ncyBtdWx0aXBsZSB0aW1lcy5cbmNvbnN0IE1vdXNlRG93biA9IDA7XG5jb25zdCBNb3VzZVVwICAgPSAxO1xuY29uc3QgTW91c2VNb3ZlID0gMjtcblxuZnVuY3Rpb24gdXBkYXRlKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgLy8gV2luZG93cyBzdXBwcmVzc2VzIG1vdXNlIGV2ZW50cyBhdCB0aGUgZW5kIG9mIGRyYWdnaW5nIG9yIHJlc2l6aW5nLFxuICAgIC8vIHNvIHdlIG5lZWQgdG8gYmUgc21hcnQgYW5kIHN5bnRoZXNpemUgYnV0dG9uIGV2ZW50cy5cblxuICAgIGxldCBldmVudFR5cGU6IG51bWJlciwgZXZlbnRCdXR0b25zID0gZXZlbnQuYnV0dG9ucztcbiAgICBzd2l0Y2ggKGV2ZW50LnR5cGUpIHtcbiAgICAgICAgY2FzZSAnbW91c2Vkb3duJzpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlRG93bjtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyB8ICgxIDw8IGV2ZW50LmJ1dHRvbik7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICBjYXNlICdtb3VzZXVwJzpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlVXA7XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgJiB+KDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZU1vdmU7XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnM7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgIH1cblxuICAgIGxldCByZWxlYXNlZCA9IGJ1dHRvbnMgJiB+ZXZlbnRCdXR0b25zO1xuICAgIGxldCBwcmVzc2VkID0gZXZlbnRCdXR0b25zICYgfmJ1dHRvbnM7XG5cbiAgICBidXR0b25zID0gZXZlbnRCdXR0b25zO1xuXG4gICAgLy8gU3ludGhlc2l6ZSBhIHJlbGVhc2UtcHJlc3Mgc2VxdWVuY2UgaWYgd2UgZGV0ZWN0IGEgcHJlc3Mgb2YgYW4gYWxyZWFkeSBwcmVzc2VkIGJ1dHRvbi5cbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZURvd24gJiYgIShwcmVzc2VkICYgZXZlbnQuYnV0dG9uKSkge1xuICAgICAgICByZWxlYXNlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xuICAgICAgICBwcmVzc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XG4gICAgfVxuXG4gICAgLy8gU3VwcHJlc3MgYWxsIGJ1dHRvbiBldmVudHMgZHVyaW5nIGRyYWdnaW5nIGFuZCByZXNpemluZyxcbiAgICAvLyB1bmxlc3MgdGhpcyBpcyBhIG1vdXNldXAgZXZlbnQgdGhhdCBpcyBlbmRpbmcgYSBkcmFnIGFjdGlvbi5cbiAgICBpZiAoXG4gICAgICAgIGV2ZW50VHlwZSAhPT0gTW91c2VNb3ZlIC8vIEZhc3QgcGF0aCBmb3IgbW91c2Vtb3ZlXG4gICAgICAgICYmIHJlc2l6aW5nXG4gICAgICAgIHx8IChcbiAgICAgICAgICAgIGRyYWdnaW5nXG4gICAgICAgICAgICAmJiAoXG4gICAgICAgICAgICAgICAgZXZlbnRUeXBlID09PSBNb3VzZURvd25cbiAgICAgICAgICAgICAgICB8fCBldmVudC5idXR0b24gIT09IDBcbiAgICAgICAgICAgIClcbiAgICAgICAgKVxuICAgICkge1xuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxuXG4gICAgLy8gSGFuZGxlIHJlbGVhc2VzXG4gICAgaWYgKHJlbGVhc2VkICYgMSkgeyBwcmltYXJ5VXAoZXZlbnQpOyB9XG4gICAgLy8gSGFuZGxlIHByZXNzZXNcbiAgICBpZiAocHJlc3NlZCAmIDEpIHsgcHJpbWFyeURvd24oZXZlbnQpOyB9XG5cbiAgICAvLyBIYW5kbGUgbW91c2Vtb3ZlXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VNb3ZlKSB7IG9uTW91c2VNb3ZlKGV2ZW50KTsgfTtcbn1cblxuZnVuY3Rpb24gcHJpbWFyeURvd24oZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcbiAgICAvLyBSZXNldCByZWFkaW5lc3Mgc3RhdGUuXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xuXG4gICAgLy8gSWdub3JlIHJlcGVhdGVkIGNsaWNrcyBvbiBtYWNPUyBhbmQgTGludXguXG4gICAgaWYgKCFJc1dpbmRvd3MoKSkge1xuICAgICAgICBpZiAoZXZlbnQudHlwZSA9PT0gJ21vdXNlZG93bicgJiYgZXZlbnQuYnV0dG9uID09PSAwICYmIGV2ZW50LmRldGFpbCAhPT0gMSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgLy8gUmVhZHkgdG8gcmVzaXplIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUuXG4gICAgICAgIGNhblJlc2l6ZSA9IHRydWU7XG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnIG9wZXJhdGlvbnMgd2hlbiBvbiByZXNpemUgZWRnZXMuXG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBSZXRyaWV2ZSB0YXJnZXQgZWxlbWVudFxuICAgIGNvbnN0IHRhcmdldCA9IGV2ZW50VGFyZ2V0KGV2ZW50KTtcblxuICAgIC8vIFJlYWR5IHRvIGRyYWcgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZSBvbiBhIGRyYWdnYWJsZSBlbGVtZW50LlxuICAgIC8vIElnbm9yZSBjbGlja3Mgb24gdGhlIHNjcm9sbGJhci5cbiAgICBjb25zdCBzdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCk7XG4gICAgY2FuRHJhZyA9IChcbiAgICAgICAgc3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpLnRyaW0oKSA9PT0gXCJkcmFnXCJcbiAgICAgICAgJiYgKFxuICAgICAgICAgICAgZXZlbnQub2Zmc2V0WCAtIHBhcnNlRmxvYXQoc3R5bGUucGFkZGluZ0xlZnQpIDwgdGFyZ2V0LmNsaWVudFdpZHRoXG4gICAgICAgICAgICAmJiBldmVudC5vZmZzZXRZIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nVG9wKSA8IHRhcmdldC5jbGllbnRIZWlnaHRcbiAgICAgICAgKVxuICAgICk7XG59XG5cbmZ1bmN0aW9uIHByaW1hcnlVcChldmVudDogTW91c2VFdmVudCkge1xuICAgIC8vIFN0b3AgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLlxuICAgIGNhbkRyYWcgPSBmYWxzZTtcbiAgICBkcmFnZ2luZyA9IGZhbHNlO1xuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xuICAgIHJlc2l6aW5nID0gZmFsc2U7XG59XG5cbmNvbnN0IGN1cnNvckZvckVkZ2UgPSBPYmplY3QuZnJlZXplKHtcbiAgICBcInNlLXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXG4gICAgXCJzdy1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxuICAgIFwibnctcmVzaXplXCI6IFwibndzZS1yZXNpemVcIixcbiAgICBcIm5lLXJlc2l6ZVwiOiBcIm5lc3ctcmVzaXplXCIsXG4gICAgXCJ3LXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxuICAgIFwibi1yZXNpemVcIjogXCJucy1yZXNpemVcIixcbiAgICBcInMtcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXG4gICAgXCJlLXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxufSlcblxuZnVuY3Rpb24gc2V0UmVzaXplKGVkZ2U/OiBrZXlvZiB0eXBlb2YgY3Vyc29yRm9yRWRnZSk6IHZvaWQge1xuICAgIGlmIChlZGdlKSB7XG4gICAgICAgIGlmICghcmVzaXplRWRnZSkgeyBkZWZhdWx0Q3Vyc29yID0gZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3I7IH1cbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBjdXJzb3JGb3JFZGdlW2VkZ2VdO1xuICAgIH0gZWxzZSBpZiAoIWVkZ2UgJiYgcmVzaXplRWRnZSkge1xuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGRlZmF1bHRDdXJzb3I7XG4gICAgfVxuXG4gICAgcmVzaXplRWRnZSA9IGVkZ2UgfHwgXCJcIjtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZU1vdmUoZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcbiAgICBpZiAoY2FuUmVzaXplICYmIHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgLy8gU3RhcnQgcmVzaXppbmcuXG4gICAgICAgIHJlc2l6aW5nID0gdHJ1ZTtcbiAgICAgICAgaW52b2tlKFwid2FpbHM6cmVzaXplOlwiICsgcmVzaXplRWRnZSk7XG4gICAgfSBlbHNlIGlmIChjYW5EcmFnKSB7XG4gICAgICAgIC8vIFN0YXJ0IGRyYWdnaW5nLlxuICAgICAgICBkcmFnZ2luZyA9IHRydWU7XG4gICAgICAgIGludm9rZShcIndhaWxzOmRyYWdcIik7XG4gICAgfVxuXG4gICAgaWYgKGRyYWdnaW5nIHx8IHJlc2l6aW5nKSB7XG4gICAgICAgIC8vIEVpdGhlciBkcmFnIG9yIHJlc2l6ZSBpcyBvbmdvaW5nLFxuICAgICAgICAvLyByZXNldCByZWFkaW5lc3MgYW5kIHN0b3AgcHJvY2Vzc2luZy5cbiAgICAgICAgY2FuRHJhZyA9IGNhblJlc2l6ZSA9IGZhbHNlO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKCFyZXNpemFibGUgfHwgIUlzV2luZG93cygpKSB7XG4gICAgICAgIGlmIChyZXNpemVFZGdlKSB7IHNldFJlc2l6ZSgpOyB9XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBjb25zdCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xuXG4gICAgLy8gRXh0cmEgcGl4ZWxzIGZvciB0aGUgY29ybmVyIGFyZWFzLlxuICAgIGNvbnN0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xuXG4gICAgY29uc3QgcmlnaHRCb3JkZXIgPSAod2luZG93Lm91dGVyV2lkdGggLSBldmVudC5jbGllbnRYKSA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGNvbnN0IGxlZnRCb3JkZXIgPSBldmVudC5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgY29uc3QgdG9wQm9yZGVyID0gZXZlbnQuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcbiAgICBjb25zdCBib3R0b21Cb3JkZXIgPSAod2luZG93Lm91dGVySGVpZ2h0IC0gZXZlbnQuY2xpZW50WSkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG5cbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lciBhcmVhcy5cbiAgICBjb25zdCByaWdodENvcm5lciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IGxlZnRDb3JuZXIgPSBldmVudC5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IHRvcENvcm5lciA9IGV2ZW50LmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IGJvdHRvbUNvcm5lciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG5cbiAgICBpZiAoIWxlZnRDb3JuZXIgJiYgIXRvcENvcm5lciAmJiAhYm90dG9tQ29ybmVyICYmICFyaWdodENvcm5lcikge1xuICAgICAgICAvLyBPcHRpbWlzYXRpb246IG91dCBvZiBhbGwgY29ybmVyIGFyZWFzIGltcGxpZXMgb3V0IG9mIGJvcmRlcnMuXG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbiAgICAvLyBEZXRlY3QgY29ybmVycy5cbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgLy8gRGV0ZWN0IGJvcmRlcnMuXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcbiAgICAvLyBPdXQgb2YgYm9yZGVyIGFyZWEuXG4gICAgZWxzZSBzZXRSZXNpemUoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5BcHBsaWNhdGlvbik7XG5cbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xuY29uc3QgU2hvd01ldGhvZCA9IDE7XG5jb25zdCBRdWl0TWV0aG9kID0gMjtcblxuLyoqXG4gKiBIaWRlcyBhIGNlcnRhaW4gbWV0aG9kIGJ5IGNhbGxpbmcgdGhlIEhpZGVNZXRob2QgZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKEhpZGVNZXRob2QpO1xufVxuXG4vKipcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xufVxuXG4vKipcbiAqIENhbGxzIHRoZSBRdWl0TWV0aG9kIHRvIHRlcm1pbmF0ZSB0aGUgcHJvZ3JhbS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IENhbmNlbGxhYmxlUHJvbWlzZSwgdHlwZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzIH0gZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbndpbmRvdy5fd2FpbHMuY2FsbFJlc3VsdEhhbmRsZXIgPSByZXN1bHRIYW5kbGVyO1xud2luZG93Ll93YWlscy5jYWxsRXJyb3JIYW5kbGVyID0gZXJyb3JIYW5kbGVyO1xuXG50eXBlIFByb21pc2VSZXNvbHZlcnMgPSBPbWl0PENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8YW55PiwgXCJwcm9taXNlXCIgfCBcIm9uY2FuY2VsbGVkXCI+XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbGwpO1xuY29uc3QgY2FuY2VsQ2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2FuY2VsQ2FsbCk7XG5jb25zdCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcDxzdHJpbmcsIFByb21pc2VSZXNvbHZlcnM+KCk7XG5cbmNvbnN0IENhbGxCaW5kaW5nID0gMDtcbmNvbnN0IENhbmNlbE1ldGhvZCA9IDBcblxuLyoqXG4gKiBIb2xkcyBhbGwgcmVxdWlyZWQgaW5mb3JtYXRpb24gZm9yIGEgYmluZGluZyBjYWxsLlxuICogTWF5IHByb3ZpZGUgZWl0aGVyIGEgbWV0aG9kIElEIG9yIGEgbWV0aG9kIG5hbWUsIGJ1dCBub3QgYm90aC5cbiAqL1xuZXhwb3J0IHR5cGUgQ2FsbE9wdGlvbnMgPSB7XG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2RJRDogbnVtYmVyO1xuICAgIC8qKiBUaGUgZnVsbHkgcXVhbGlmaWVkIG5hbWUgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZE5hbWU/OiBuZXZlcjtcbiAgICAvKiogQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuICovXG4gICAgYXJnczogYW55W107XG59IHwge1xuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kSUQ/OiBuZXZlcjtcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2ROYW1lOiBzdHJpbmc7XG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xuICAgIGFyZ3M6IGFueVtdO1xufTtcblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHRocm93biBpbiBjYXNlIHRoZSBib3VuZCBtZXRob2QgcmV0dXJucyBhbiBlcnJvci5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIFJ1bnRpbWVFcnJvciNuYW1lfSBwcm9wZXJ0eSBpcyBcIlJ1bnRpbWVFcnJvclwiLlxuICovXG5leHBvcnQgY2xhc3MgUnVudGltZUVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgUnVudGltZUVycm9yIGluc3RhbmNlLlxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gICAgICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIHRvIGJlIGZvcndhcmRlZCB0byB0aGUgRXJyb3IgY29uc3RydWN0b3IuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xuICAgICAgICBzdXBlcihtZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJSdW50aW1lRXJyb3JcIjtcbiAgICB9XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgcmVzdWx0IG9mIGEgY2FsbCByZXF1ZXN0LlxuICpcbiAqIEBwYXJhbSBpZCAtIFRoZSBpZCBvZiB0aGUgcmVxdWVzdCB0byBoYW5kbGUgdGhlIHJlc3VsdCBmb3IuXG4gKiBAcGFyYW0gZGF0YSAtIFRoZSByZXN1bHQgZGF0YSBvZiB0aGUgcmVxdWVzdC5cbiAqIEBwYXJhbSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cbiAqL1xuZnVuY3Rpb24gcmVzdWx0SGFuZGxlcihpZDogc3RyaW5nLCBkYXRhOiBzdHJpbmcsIGlzSlNPTjogYm9vbGVhbik6IHZvaWQge1xuICAgIGNvbnN0IHJlc29sdmVycyA9IGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKTtcbiAgICBpZiAoIXJlc29sdmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKCFkYXRhKSB7XG4gICAgICAgIHJlc29sdmVycy5yZXNvbHZlKHVuZGVmaW5lZCk7XG4gICAgfSBlbHNlIGlmICghaXNKU09OKSB7XG4gICAgICAgIHJlc29sdmVycy5yZXNvbHZlKGRhdGEpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICByZXNvbHZlcnMucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyOiBhbnkpIHtcbiAgICAgICAgICAgIHJlc29sdmVycy5yZWplY3QobmV3IFR5cGVFcnJvcihcImNvdWxkIG5vdCBwYXJzZSByZXN1bHQ6IFwiICsgZXJyLm1lc3NhZ2UsIHsgY2F1c2U6IGVyciB9KSk7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogSGFuZGxlcyB0aGUgZXJyb3IgZnJvbSBhIGNhbGwgcmVxdWVzdC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgaWQgb2YgdGhlIHByb21pc2UgaGFuZGxlci5cbiAqIEBwYXJhbSBkYXRhIC0gVGhlIGVycm9yIGRhdGEgdG8gcmVqZWN0IHRoZSBwcm9taXNlIGhhbmRsZXIgd2l0aC5cbiAqIEBwYXJhbSBpc0pTT04gLSBJbmRpY2F0ZXMgd2hldGhlciB0aGUgZGF0YSBpcyBKU09OIG9yIG5vdC5cbiAqL1xuZnVuY3Rpb24gZXJyb3JIYW5kbGVyKGlkOiBzdHJpbmcsIGRhdGE6IHN0cmluZywgaXNKU09OOiBib29sZWFuKTogdm9pZCB7XG4gICAgY29uc3QgcmVzb2x2ZXJzID0gZ2V0QW5kRGVsZXRlUmVzcG9uc2UoaWQpO1xuICAgIGlmICghcmVzb2x2ZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIWlzSlNPTikge1xuICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBFcnJvcihkYXRhKSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgbGV0IGVycm9yOiBhbnk7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBlcnJvciA9IEpTT04ucGFyc2UoZGF0YSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycjogYW55KSB7XG4gICAgICAgICAgICByZXNvbHZlcnMucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJjb3VsZCBub3QgcGFyc2UgZXJyb3I6IFwiICsgZXJyLm1lc3NhZ2UsIHsgY2F1c2U6IGVyciB9KSk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cblxuICAgICAgICBsZXQgb3B0aW9uczogRXJyb3JPcHRpb25zID0ge307XG4gICAgICAgIGlmIChlcnJvci5jYXVzZSkge1xuICAgICAgICAgICAgb3B0aW9ucy5jYXVzZSA9IGVycm9yLmNhdXNlO1xuICAgICAgICB9XG5cbiAgICAgICAgbGV0IGV4Y2VwdGlvbjtcbiAgICAgICAgc3dpdGNoIChlcnJvci5raW5kKSB7XG4gICAgICAgICAgICBjYXNlIFwiUmVmZXJlbmNlRXJyb3JcIjpcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgUmVmZXJlbmNlRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICBjYXNlIFwiVHlwZUVycm9yXCI6XG4gICAgICAgICAgICAgICAgZXhjZXB0aW9uID0gbmV3IFR5cGVFcnJvcihlcnJvci5tZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgIGNhc2UgXCJSdW50aW1lRXJyb3JcIjpcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgUnVudGltZUVycm9yKGVycm9yLm1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgICAgICBleGNlcHRpb24gPSBuZXcgRXJyb3IoZXJyb3IubWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cblxuICAgICAgICByZXNvbHZlcnMucmVqZWN0KGV4Y2VwdGlvbik7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHJpZXZlcyBhbmQgcmVtb3ZlcyB0aGUgcmVzcG9uc2UgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCBmcm9tIHRoZSBjYWxsUmVzcG9uc2VzIG1hcC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgSUQgb2YgdGhlIHJlc3BvbnNlIHRvIGJlIHJldHJpZXZlZCBhbmQgcmVtb3ZlZC5cbiAqIEByZXR1cm5zIFRoZSByZXNwb25zZSBvYmplY3QgYXNzb2NpYXRlZCB3aXRoIHRoZSBnaXZlbiBJRCwgaWYgYW55LlxuICovXG5mdW5jdGlvbiBnZXRBbmREZWxldGVSZXNwb25zZShpZDogc3RyaW5nKTogUHJvbWlzZVJlc29sdmVycyB8IHVuZGVmaW5lZCB7XG4gICAgY29uc3QgcmVzcG9uc2UgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIHJldHVybiByZXNwb25zZTtcbn1cblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxuICpcbiAqIEByZXR1cm5zIEEgdW5pcXVlIElEIHRoYXQgZG9lcyBub3QgZXhpc3QgaW4gdGhlIGNhbGxSZXNwb25zZXMgc2V0LlxuICovXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCk6IHN0cmluZyB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogQ2FsbCBhIGJvdW5kIG1ldGhvZCBhY2NvcmRpbmcgdG8gdGhlIGdpdmVuIGNhbGwgb3B0aW9ucy5cbiAqXG4gKiBJbiBjYXNlIG9mIGZhaWx1cmUsIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYW4gZXhjZXB0aW9uXG4gKiBhbW9uZyBSZWZlcmVuY2VFcnJvciAodW5rbm93biBtZXRob2QpLCBUeXBlRXJyb3IgKHdyb25nIGFyZ3VtZW50IGNvdW50IG9yIHR5cGUpLFxuICoge0BsaW5rIFJ1bnRpbWVFcnJvcn0gKG1ldGhvZCByZXR1cm5lZCBhbiBlcnJvciksIG9yIG90aGVyIChuZXR3b3JrIG9yIGludGVybmFsIGVycm9ycykuXG4gKiBUaGUgZXhjZXB0aW9uIG1pZ2h0IGhhdmUgYSBcImNhdXNlXCIgZmllbGQgd2l0aCB0aGUgdmFsdWUgcmV0dXJuZWRcbiAqIGJ5IHRoZSBhcHBsaWNhdGlvbi0gb3Igc2VydmljZS1sZXZlbCBlcnJvciBtYXJzaGFsaW5nIGZ1bmN0aW9ucy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIEEgbWV0aG9kIGNhbGwgZGVzY3JpcHRvci5cbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnM6IENhbGxPcHRpb25zKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xuICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xuXG4gICAgY29uc3QgcmVzdWx0ID0gQ2FuY2VsbGFibGVQcm9taXNlLndpdGhSZXNvbHZlcnM8YW55PigpO1xuICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmU6IHJlc3VsdC5yZXNvbHZlLCByZWplY3Q6IHJlc3VsdC5yZWplY3QgfSk7XG5cbiAgICBjb25zdCByZXF1ZXN0ID0gY2FsbChDYWxsQmluZGluZywgT2JqZWN0LmFzc2lnbih7IFwiY2FsbC1pZFwiOiBpZCB9LCBvcHRpb25zKSk7XG4gICAgbGV0IHJ1bm5pbmcgPSBmYWxzZTtcblxuICAgIHJlcXVlc3QudGhlbigoKSA9PiB7XG4gICAgICAgIHJ1bm5pbmcgPSB0cnVlO1xuICAgIH0sIChlcnIpID0+IHtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXN1bHQucmVqZWN0KGVycik7XG4gICAgfSk7XG5cbiAgICBjb25zdCBjYW5jZWwgPSAoKSA9PiB7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgcmV0dXJuIGNhbmNlbENhbGwoQ2FuY2VsTWV0aG9kLCB7XCJjYWxsLWlkXCI6IGlkfSkuY2F0Y2goKGVycikgPT4ge1xuICAgICAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIHdoaWxlIHJlcXVlc3RpbmcgYmluZGluZyBjYWxsIGNhbmNlbGxhdGlvbjpcIiwgZXJyKTtcbiAgICAgICAgfSk7XG4gICAgfTtcblxuICAgIHJlc3VsdC5vbmNhbmNlbGxlZCA9ICgpID0+IHtcbiAgICAgICAgaWYgKHJ1bm5pbmcpIHtcbiAgICAgICAgICAgIHJldHVybiBjYW5jZWwoKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHJldHVybiByZXF1ZXN0LnRoZW4oY2FuY2VsKTtcbiAgICAgICAgfVxuICAgIH07XG5cbiAgICByZXR1cm4gcmVzdWx0LnByb21pc2U7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBib3VuZCBtZXRob2QgYnkgbmFtZSB3aXRoIHRoZSBzcGVjaWZpZWQgYXJndW1lbnRzLlxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cbiAqXG4gKiBAcGFyYW0gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJy5cbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShtZXRob2ROYW1lOiBzdHJpbmcsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kTmFtZSwgYXJncyB9KTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBieSBpdHMgbnVtZXJpYyBJRCB3aXRoIHRoZSBzcGVjaWZpZWQgYXJndW1lbnRzLlxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cbiAqXG4gKiBAcGFyYW0gbWV0aG9kSUQgLSBUaGUgSUQgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEByZXR1cm4gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCeUlEKG1ldGhvZElEOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kSUQsIGFyZ3MgfSk7XG59XG4iLCAiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vaW5zcGVjdC1qcy9pcy1jYWxsYWJsZVxuXG4vLyBUaGUgTUlUIExpY2Vuc2UgKE1JVClcbi8vXG4vLyBDb3B5cmlnaHQgKGMpIDIwMTUgSm9yZGFuIEhhcmJhbmRcbi8vXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5XG4vLyBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsXG4vLyBpbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzXG4vLyB0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsXG4vLyBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXNcbi8vIGZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6XG4vL1xuLy8gVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXG4vLyBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLlxuLy9cbi8vIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1Jcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLFxuLy8gRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFXG4vLyBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSXG4vLyBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLFxuLy8gT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEVcbi8vIFNPRlRXQVJFLlxuXG52YXIgZm5Ub1N0ciA9IEZ1bmN0aW9uLnByb3RvdHlwZS50b1N0cmluZztcbnZhciByZWZsZWN0QXBwbHk6IHR5cGVvZiBSZWZsZWN0LmFwcGx5IHwgZmFsc2UgfCBudWxsID0gdHlwZW9mIFJlZmxlY3QgPT09ICdvYmplY3QnICYmIFJlZmxlY3QgIT09IG51bGwgJiYgUmVmbGVjdC5hcHBseTtcbnZhciBiYWRBcnJheUxpa2U6IGFueTtcbnZhciBpc0NhbGxhYmxlTWFya2VyOiBhbnk7XG5pZiAodHlwZW9mIHJlZmxlY3RBcHBseSA9PT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgT2JqZWN0LmRlZmluZVByb3BlcnR5ID09PSAnZnVuY3Rpb24nKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgYmFkQXJyYXlMaWtlID0gT2JqZWN0LmRlZmluZVByb3BlcnR5KHt9LCAnbGVuZ3RoJywge1xuICAgICAgICAgICAgZ2V0OiBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICAgICAgdGhyb3cgaXNDYWxsYWJsZU1hcmtlcjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgICAgIGlzQ2FsbGFibGVNYXJrZXIgPSB7fTtcbiAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXRocm93LWxpdGVyYWxcbiAgICAgICAgcmVmbGVjdEFwcGx5KGZ1bmN0aW9uICgpIHsgdGhyb3cgNDI7IH0sIG51bGwsIGJhZEFycmF5TGlrZSk7XG4gICAgfSBjYXRjaCAoXykge1xuICAgICAgICBpZiAoXyAhPT0gaXNDYWxsYWJsZU1hcmtlcikge1xuICAgICAgICAgICAgcmVmbGVjdEFwcGx5ID0gbnVsbDtcbiAgICAgICAgfVxuICAgIH1cbn0gZWxzZSB7XG4gICAgcmVmbGVjdEFwcGx5ID0gbnVsbDtcbn1cblxudmFyIGNvbnN0cnVjdG9yUmVnZXggPSAvXlxccypjbGFzc1xcYi87XG52YXIgaXNFUzZDbGFzc0ZuID0gZnVuY3Rpb24gaXNFUzZDbGFzc0Z1bmN0aW9uKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcbiAgICB0cnkge1xuICAgICAgICB2YXIgZm5TdHIgPSBmblRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICByZXR1cm4gY29uc3RydWN0b3JSZWdleC50ZXN0KGZuU3RyKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTsgLy8gbm90IGEgZnVuY3Rpb25cbiAgICB9XG59O1xuXG52YXIgdHJ5RnVuY3Rpb25PYmplY3QgPSBmdW5jdGlvbiB0cnlGdW5jdGlvblRvU3RyKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcbiAgICB0cnkge1xuICAgICAgICBpZiAoaXNFUzZDbGFzc0ZuKHZhbHVlKSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICAgICAgZm5Ub1N0ci5jYWxsKHZhbHVlKTtcbiAgICAgICAgcmV0dXJuIHRydWU7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxufTtcbnZhciB0b1N0ciA9IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmc7XG52YXIgb2JqZWN0Q2xhc3MgPSAnW29iamVjdCBPYmplY3RdJztcbnZhciBmbkNsYXNzID0gJ1tvYmplY3QgRnVuY3Rpb25dJztcbnZhciBnZW5DbGFzcyA9ICdbb2JqZWN0IEdlbmVyYXRvckZ1bmN0aW9uXSc7XG52YXIgZGRhQ2xhc3MgPSAnW29iamVjdCBIVE1MQWxsQ29sbGVjdGlvbl0nOyAvLyBJRSAxMVxudmFyIGRkYUNsYXNzMiA9ICdbb2JqZWN0IEhUTUwgZG9jdW1lbnQuYWxsIGNsYXNzXSc7XG52YXIgZGRhQ2xhc3MzID0gJ1tvYmplY3QgSFRNTENvbGxlY3Rpb25dJzsgLy8gSUUgOS0xMFxudmFyIGhhc1RvU3RyaW5nVGFnID0gdHlwZW9mIFN5bWJvbCA9PT0gJ2Z1bmN0aW9uJyAmJiAhIVN5bWJvbC50b1N0cmluZ1RhZzsgLy8gYmV0dGVyOiB1c2UgYGhhcy10b3N0cmluZ3RhZ2BcblxudmFyIGlzSUU2OCA9ICEoMCBpbiBbLF0pOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIG5vLXNwYXJzZS1hcnJheXMsIGNvbW1hLXNwYWNpbmdcblxudmFyIGlzRERBOiAodmFsdWU6IGFueSkgPT4gYm9vbGVhbiA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwoKSB7IHJldHVybiBmYWxzZTsgfTtcbmlmICh0eXBlb2YgZG9jdW1lbnQgPT09ICdvYmplY3QnKSB7XG4gICAgLy8gRmlyZWZveCAzIGNhbm9uaWNhbGl6ZXMgRERBIHRvIHVuZGVmaW5lZCB3aGVuIGl0J3Mgbm90IGFjY2Vzc2VkIGRpcmVjdGx5XG4gICAgdmFyIGFsbCA9IGRvY3VtZW50LmFsbDtcbiAgICBpZiAodG9TdHIuY2FsbChhbGwpID09PSB0b1N0ci5jYWxsKGRvY3VtZW50LmFsbCkpIHtcbiAgICAgICAgaXNEREEgPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKHZhbHVlKSB7XG4gICAgICAgICAgICAvKiBnbG9iYWxzIGRvY3VtZW50OiBmYWxzZSAqL1xuICAgICAgICAgICAgLy8gaW4gSUUgNi04LCB0eXBlb2YgZG9jdW1lbnQuYWxsIGlzIFwib2JqZWN0XCIgYW5kIGl0J3MgdHJ1dGh5XG4gICAgICAgICAgICBpZiAoKGlzSUU2OCB8fCAhdmFsdWUpICYmICh0eXBlb2YgdmFsdWUgPT09ICd1bmRlZmluZWQnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcpKSB7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgdmFyIHN0ciA9IHRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICAgICAgICAgICAgICByZXR1cm4gKFxuICAgICAgICAgICAgICAgICAgICAgICAgc3RyID09PSBkZGFDbGFzc1xuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MzIC8vIG9wZXJhIDEyLjE2XG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IG9iamVjdENsYXNzIC8vIElFIDYtOFxuICAgICAgICAgICAgICAgICAgICApICYmIHZhbHVlKCcnKSA9PSBudWxsOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIGVxZXFlcVxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGUpIHsgLyoqLyB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5mdW5jdGlvbiBpc0NhbGxhYmxlUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkgIHtcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAodHlwZW9mIHZhbHVlICE9PSAnZnVuY3Rpb24nICYmIHR5cGVvZiB2YWx1ZSAhPT0gJ29iamVjdCcpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgdHJ5IHtcbiAgICAgICAgKHJlZmxlY3RBcHBseSBhcyBhbnkpKHZhbHVlLCBudWxsLCBiYWRBcnJheUxpa2UpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgaWYgKGUgIT09IGlzQ2FsbGFibGVNYXJrZXIpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgfVxuICAgIHJldHVybiAhaXNFUzZDbGFzc0ZuKHZhbHVlKSAmJiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7XG59XG5cbmZ1bmN0aW9uIGlzQ2FsbGFibGVOb1JlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55IHtcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAodHlwZW9mIHZhbHVlICE9PSAnZnVuY3Rpb24nICYmIHR5cGVvZiB2YWx1ZSAhPT0gJ29iamVjdCcpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgaWYgKGhhc1RvU3RyaW5nVGFnKSB7IHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7IH1cbiAgICBpZiAoaXNFUzZDbGFzc0ZuKHZhbHVlKSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB2YXIgc3RyQ2xhc3MgPSB0b1N0ci5jYWxsKHZhbHVlKTtcbiAgICBpZiAoc3RyQ2xhc3MgIT09IGZuQ2xhc3MgJiYgc3RyQ2xhc3MgIT09IGdlbkNsYXNzICYmICEoL15cXFtvYmplY3QgSFRNTC8pLnRlc3Qoc3RyQ2xhc3MpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7XG59O1xuXG5leHBvcnQgZGVmYXVsdCByZWZsZWN0QXBwbHkgPyBpc0NhbGxhYmxlUmVmQXBwbHkgOiBpc0NhbGxhYmxlTm9SZWZBcHBseTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IGlzQ2FsbGFibGUgZnJvbSBcIi4vY2FsbGFibGUuanNcIjtcblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHVzZWQgYXMgcmVqZWN0aW9uIHJlYXNvblxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IGlzIGNhbmNlbGxlZCBzdWNjZXNzZnVsbHkuXG4gKlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxFcnJvclwiYC5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgY2F1c2UgcGFzc2VkIHRvIHRoZSBjYW5jZWwgbWV0aG9kLCBpZiBhbnkuXG4gKi9cbmV4cG9ydCBjbGFzcyBDYW5jZWxFcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxFcnJvcmAgaW5zdGFuY2UuXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbEVycm9yXCI7XG4gICAgfVxufVxuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvblxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IHJlamVjdHMgYWZ0ZXIgYmVpbmcgY2FuY2VsbGVkLFxuICogb3Igd2hlbiB0aGUgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB0aHJvd3Mgb3IgcmVqZWN0cy5cbiAqXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCJgLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSByZWFzb24gdGhlIHByb21pc2UgcmVqZWN0ZWQgd2l0aC5cbiAqXG4gKiBCZWNhdXNlIHRoZSBvcmlnaW5hbCBwcm9taXNlIHdhcyBjYW5jZWxsZWQsXG4gKiBhIHdyYXBwZXIgcHJvbWlzZSB3aWxsIGJlIHBhc3NlZCB0byB0aGUgdW5oYW5kbGVkIHJlamVjdGlvbiBsaXN0ZW5lciBpbnN0ZWFkLlxuICogVGhlIHtAbGluayBwcm9taXNlfSBwcm9wZXJ0eSBob2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgb3JpZ2luYWwgcHJvbWlzZS5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbGxlZFJlamVjdGlvbkVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIEhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBwcm9taXNlIHRoYXQgd2FzIGNhbmNlbGxlZCBhbmQgdGhlbiByZWplY3RlZC5cbiAgICAgKi9cbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj47XG5cbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcmAgaW5zdGFuY2UuXG4gICAgICogQHBhcmFtIHByb21pc2UgLSBUaGUgcHJvbWlzZSB0aGF0IGNhdXNlZCB0aGUgZXJyb3Igb3JpZ2luYWxseS5cbiAgICAgKiBAcGFyYW0gcmVhc29uIC0gVGhlIHJlamVjdGlvbiByZWFzb24uXG4gICAgICogQHBhcmFtIGluZm8gLSBBbiBvcHRpb25hbCBpbmZvcm1hdGl2ZSBtZXNzYWdlIHNwZWNpZnlpbmcgdGhlIGNpcmN1bXN0YW5jZXMgaW4gd2hpY2ggdGhlIGVycm9yIHdhcyB0aHJvd24uXG4gICAgICogICAgICAgICAgICAgICBEZWZhdWx0cyB0byB0aGUgc3RyaW5nIGBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCJgLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgcmVhc29uPzogYW55LCBpbmZvPzogc3RyaW5nKSB7XG4gICAgICAgIHN1cGVyKChpbmZvID8/IFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cIikgKyBcIiBSZWFzb246IFwiICsgZXJyb3JNZXNzYWdlKHJlYXNvbiksIHsgY2F1c2U6IHJlYXNvbiB9KTtcbiAgICAgICAgdGhpcy5wcm9taXNlID0gcHJvbWlzZTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiO1xuICAgIH1cbn1cblxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiA9ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IgPSAocmVhc29uPzogYW55KSA9PiB2b2lkO1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgPSAoY2F1c2U/OiBhbnkpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4gPSAocmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4sIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IpID0+IHZvaWQ7XG5cbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XG4gICAgdGhlbjxUUmVzdWx0MSA9IFQsIFRSZXN1bHQyID0gbmV2ZXI+KG9uZnVsZmlsbGVkPzogKCh2YWx1ZTogVCkgPT4gVFJlc3VsdDEgfCBQcm9taXNlTGlrZTxUUmVzdWx0MT4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gVFJlc3VsdDIgfCBQcm9taXNlTGlrZTxUUmVzdWx0Mj4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQyPikgfCB1bmRlZmluZWQgfCBudWxsKTogQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MSB8IFRSZXN1bHQyPjtcbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG59XG5cbi8qKlxuICogV3JhcHMgYSBjYW5jZWxsYWJsZSBwcm9taXNlIGFsb25nIHdpdGggaXRzIHJlc29sdXRpb24gbWV0aG9kcy5cbiAqIFRoZSBgb25jYW5jZWxsZWRgIGZpZWxkIHdpbGwgYmUgbnVsbCBpbml0aWFsbHkgYnV0IG1heSBiZSBzZXQgdG8gcHJvdmlkZSBhIGN1c3RvbSBjYW5jZWxsYXRpb24gZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XG4gICAgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xuICAgIHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+O1xuICAgIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3I7XG4gICAgb25jYW5jZWxsZWQ6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciB8IG51bGw7XG59XG5cbmludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSB7XG4gICAgcmVhZG9ubHkgcm9vdDogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGU7XG4gICAgcmVzb2x2aW5nOiBib29sZWFuO1xuICAgIHNldHRsZWQ6IGJvb2xlYW47XG4gICAgcmVhc29uPzogQ2FuY2VsRXJyb3I7XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBiYXJyaWVyU3ltID0gU3ltYm9sKFwiYmFycmllclwiKTtcbmNvbnN0IGNhbmNlbEltcGxTeW0gPSBTeW1ib2woXCJjYW5jZWxJbXBsXCIpO1xuY29uc3Qgc3BlY2llcyA9IFN5bWJvbC5zcGVjaWVzID8/IFN5bWJvbChcInNwZWNpZXNQb2x5ZmlsbFwiKTtcblxuLyoqXG4gKiBBIHByb21pc2Ugd2l0aCBhbiBhdHRhY2hlZCBtZXRob2QgZm9yIGNhbmNlbGxpbmcgbG9uZy1ydW5uaW5nIG9wZXJhdGlvbnMgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbH0pLlxuICogQ2FuY2VsbGF0aW9uIGNhbiBvcHRpb25hbGx5IGJlIGJvdW5kIHRvIGFuIHtAbGluayBBYm9ydFNpZ25hbH1cbiAqIGZvciBiZXR0ZXIgY29tcG9zYWJpbGl0eSAoc2VlIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2UjY2FuY2VsT259KS5cbiAqXG4gKiBDYW5jZWxsaW5nIGEgcGVuZGluZyBwcm9taXNlIHdpbGwgcmVzdWx0IGluIGFuIGltbWVkaWF0ZSByZWplY3Rpb25cbiAqIHdpdGggYW4gaW5zdGFuY2Ugb2Yge0BsaW5rIENhbmNlbEVycm9yfSBhcyByZWFzb24sXG4gKiBidXQgd2hvZXZlciBzdGFydGVkIHRoZSBwcm9taXNlIHdpbGwgYmUgcmVzcG9uc2libGVcbiAqIGZvciBhY3R1YWxseSBhYm9ydGluZyB0aGUgdW5kZXJseWluZyBvcGVyYXRpb24uXG4gKiBUbyB0aGlzIHB1cnBvc2UsIHRoZSBjb25zdHJ1Y3RvciBhbmQgYWxsIGNoYWluaW5nIG1ldGhvZHNcbiAqIGFjY2VwdCBvcHRpb25hbCBjYW5jZWxsYXRpb24gY2FsbGJhY2tzLlxuICpcbiAqIElmIGEgYENhbmNlbGxhYmxlUHJvbWlzZWAgc3RpbGwgcmVzb2x2ZXMgYWZ0ZXIgaGF2aW5nIGJlZW4gY2FuY2VsbGVkLFxuICogdGhlIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZC4gSWYgaXQgcmVqZWN0cywgdGhlIHJlYXNvblxuICogd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLFxuICogd3JhcHBlZCBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gaW5zdGFuY2UuXG4gKiBUbyBmYWNpbGl0YXRlIHRoZSBoYW5kbGluZyBvZiBjYW5jZWxsYXRpb24gcmVxdWVzdHMsXG4gKiBjYW5jZWxsZWQgYENhbmNlbGxhYmxlUHJvbWlzZWBzIHdpbGwgX25vdF8gcmVwb3J0IHVuaGFuZGxlZCBgQ2FuY2VsRXJyb3Jgc1xuICogd2hvc2UgYGNhdXNlYCBmaWVsZCBpcyB0aGUgc2FtZSBhcyB0aGUgb25lIHdpdGggd2hpY2ggdGhlIGN1cnJlbnQgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLlxuICpcbiAqIEFsbCB1c3VhbCBwcm9taXNlIG1ldGhvZHMgYXJlIGRlZmluZWQgYW5kIHJldHVybiBhIGBDYW5jZWxsYWJsZVByb21pc2VgXG4gKiB3aG9zZSBjYW5jZWwgbWV0aG9kIHdpbGwgY2FuY2VsIHRoZSBwYXJlbnQgb3BlcmF0aW9uIGFzIHdlbGwsIHByb3BhZ2F0aW5nIHRoZSBjYW5jZWxsYXRpb24gcmVhc29uXG4gKiB1cHdhcmRzIHRocm91Z2ggcHJvbWlzZSBjaGFpbnMuXG4gKiBDb252ZXJzZWx5LCBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aWxsIG5vdCBhdXRvbWF0aWNhbGx5IGNhbmNlbCBkZXBlbmRlbnQgcHJvbWlzZXMgZG93bnN0cmVhbTpcbiAqIGBgYHRzXG4gKiBsZXQgcm9vdCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gKiBsZXQgY2hpbGQxID0gcm9vdC50aGVuKCgpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMiA9IGNoaWxkMS50aGVuKCgpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMyA9IHJvb3QuY2F0Y2goKCkgPT4geyAuLi4gfSk7XG4gKiBjaGlsZDEuY2FuY2VsKCk7IC8vIENhbmNlbHMgY2hpbGQxIGFuZCByb290LCBidXQgbm90IGNoaWxkMiBvciBjaGlsZDNcbiAqIGBgYFxuICogQ2FuY2VsbGluZyBhIHByb21pc2UgdGhhdCBoYXMgYWxyZWFkeSBzZXR0bGVkIGlzIHNhZmUgYW5kIGhhcyBubyBjb25zZXF1ZW5jZS5cbiAqXG4gKiBUaGUgYGNhbmNlbGAgbWV0aG9kIHJldHVybnMgYSBwcm9taXNlIHRoYXQgX2Fsd2F5cyBmdWxmaWxsc19cbiAqIGFmdGVyIHRoZSB3aG9sZSBjaGFpbiBoYXMgcHJvY2Vzc2VkIHRoZSBjYW5jZWwgcmVxdWVzdFxuICogYW5kIGFsbCBhdHRhY2hlZCBjYWxsYmFja3MgdXAgdG8gdGhhdCBtb21lbnQgaGF2ZSBydW4uXG4gKlxuICogQWxsIEVTMjAyNCBwcm9taXNlIG1ldGhvZHMgKHN0YXRpYyBhbmQgaW5zdGFuY2UpIGFyZSBkZWZpbmVkIG9uIENhbmNlbGxhYmxlUHJvbWlzZSxcbiAqIGJ1dCBhY3R1YWwgYXZhaWxhYmlsaXR5IG1heSB2YXJ5IHdpdGggT1Mvd2VidmlldyB2ZXJzaW9uLlxuICpcbiAqIEluIGxpbmUgd2l0aCB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZyxcbiAqIGBDYW5jZWxsYWJsZVByb21pc2VgIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuXG4gKiBFeHRlbmRlcnMgc2hvdWxkIHRha2UgY2FyZSB0byBwcm92aWRlIHRoZWlyIG93biBtZXRob2QgaW1wbGVtZW50YXRpb25zLlxuICogVGhpcyBtaWdodCBiZSByZWNvbnNpZGVyZWQgaW4gY2FzZSB0aGUgcHJvcG9zYWwgaXMgcmV0aXJlZC5cbiAqXG4gKiBDYW5jZWxsYWJsZVByb21pc2UgaXMgYSB3cmFwcGVyIGFyb3VuZCB0aGUgRE9NIFByb21pc2Ugb2JqZWN0XG4gKiBhbmQgaXMgY29tcGxpYW50IHdpdGggdGhlIFtQcm9taXNlcy9BKyBzcGVjaWZpY2F0aW9uXShodHRwczovL3Byb21pc2VzYXBsdXMuY29tLylcbiAqIChpdCBwYXNzZXMgdGhlIFtjb21wbGlhbmNlIHN1aXRlXShodHRwczovL2dpdGh1Yi5jb20vcHJvbWlzZXMtYXBsdXMvcHJvbWlzZXMtdGVzdHMpKVxuICogaWYgc28gaXMgdGhlIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24uXG4gKi9cbmV4cG9ydCBjbGFzcyBDYW5jZWxsYWJsZVByb21pc2U8VD4gZXh0ZW5kcyBQcm9taXNlPFQ+IGltcGxlbWVudHMgUHJvbWlzZUxpa2U8VD4sIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIC8qKiBAaW50ZXJuYWwgKi9cbiAgICBwcml2YXRlIFtiYXJyaWVyU3ltXSE6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+IHwgbnVsbDtcbiAgICAvKiogQGludGVybmFsICovXG4gICAgcHJpdmF0ZSByZWFkb25seSBbY2FuY2VsSW1wbFN5bV0hOiAocmVhc29uOiBDYW5jZWxFcnJvcikgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBgQ2FuY2VsbGFibGVQcm9taXNlYC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBleGVjdXRvciAtIEEgY2FsbGJhY2sgdXNlZCB0byBpbml0aWFsaXplIHRoZSBwcm9taXNlLiBUaGlzIGNhbGxiYWNrIGlzIHBhc3NlZCB0d28gYXJndW1lbnRzOlxuICAgICAqICAgICAgICAgICAgICAgICAgIGEgYHJlc29sdmVgIGNhbGxiYWNrIHVzZWQgdG8gcmVzb2x2ZSB0aGUgcHJvbWlzZSB3aXRoIGEgdmFsdWVcbiAgICAgKiAgICAgICAgICAgICAgICAgICBvciB0aGUgcmVzdWx0IG9mIGFub3RoZXIgcHJvbWlzZSAocG9zc2libHkgY2FuY2VsbGFibGUpLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGFuZCBhIGByZWplY3RgIGNhbGxiYWNrIHVzZWQgdG8gcmVqZWN0IHRoZSBwcm9taXNlIHdpdGggYSBwcm92aWRlZCByZWFzb24gb3IgZXJyb3IuXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgdGhlIHZhbHVlIHByb3ZpZGVkIHRvIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiAgICAgICAgICAgICAgICAgICAoaXQgaGFzIGEgYHRoZW5gIF9hbmRfIGEgYGNhbmNlbGAgbWV0aG9kKSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBmb3J3YXJkZWQgdG8gdGhhdCBvYmplY3QgYW5kIHRoZSBvbmNhbmNlbGxlZCB3aWxsIG5vdCBiZSBpbnZva2VkIGFueW1vcmUuXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgYW55IG9uZSBvZiB0aGUgdHdvIGNhbGxiYWNrcyBpcyBjYWxsZWQgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQsXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhlIHByb3ZpZGVkIHZhbHVlcyB3aWxsIGJlIGNhbmNlbGxlZCBhbmQgcmVzb2x2ZWQgYXMgdXN1YWwsXG4gICAgICogICAgICAgICAgICAgICAgICAgYnV0IHRoZWlyIHJlc3VsdHMgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogICAgICAgICAgICAgICAgICAgSG93ZXZlciwgaWYgdGhlIHJlc29sdXRpb24gcHJvY2VzcyB1bHRpbWF0ZWx5IGVuZHMgdXAgaW4gYSByZWplY3Rpb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGF0IGlzIG5vdCBkdWUgdG8gY2FuY2VsbGF0aW9uLCB0aGUgcmVqZWN0aW9uIHJlYXNvblxuICAgICAqICAgICAgICAgICAgICAgICAgIHdpbGwgYmUgd3JhcHBlZCBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn1cbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqIEBwYXJhbSBvbmNhbmNlbGxlZCAtIEl0IGlzIHRoZSBjYWxsZXIncyByZXNwb25zaWJpbGl0eSB0byBlbnN1cmUgdGhhdCBhbnkgb3BlcmF0aW9uXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgc3RhcnRlZCBieSB0aGUgZXhlY3V0b3IgaXMgcHJvcGVybHkgaGFsdGVkIHVwb24gY2FuY2VsbGF0aW9uLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFRoaXMgb3B0aW9uYWwgY2FsbGJhY2sgY2FuIGJlIHVzZWQgdG8gdGhhdCBwdXJwb3NlLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIEl0IHdpbGwgYmUgY2FsbGVkIF9zeW5jaHJvbm91c2x5XyB3aXRoIGEgY2FuY2VsbGF0aW9uIGNhdXNlXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2hlbiBjYW5jZWxsYXRpb24gaXMgcmVxdWVzdGVkLCBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHJlamVjdGVkXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0sIGJ1dCBfYmVmb3JlX1xuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFueSB7QGxpbmsgdGhlbn0ve0BsaW5rIGNhdGNofS97QGxpbmsgZmluYWxseX0gY2FsbGJhY2sgcnVucy5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9XG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2lsbCBvbmx5IGZ1bGZpbGwgYWZ0ZXIgdGhlIGZvcm1lciBoYXMgc2V0dGxlZC5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBVbmhhbmRsZWQgZXhjZXB0aW9ucyBvciByZWplY3Rpb25zIGZyb20gdGhlIGNhbGxiYWNrIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyB1bmhhbmRsZWQgcmVqZWN0aW9ucy5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGNhbGxlZCBiZWZvcmUgY2FuY2VsbGF0aW9uIHdpdGggYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGlzIHByb21pc2Ugd2lsbCBiZSBkaXZlcnRlZCB0byB0aGF0IHByb21pc2UsXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW5kIHRoZSBvcmlnaW5hbCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKGV4ZWN1dG9yOiBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpIHtcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG4gICAgICAgIHN1cGVyKChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xuXG4gICAgICAgIGlmICgodGhpcy5jb25zdHJ1Y3RvciBhcyBhbnkpW3NwZWNpZXNdICE9PSBQcm9taXNlKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuIFBsZWFzZSByZWZyYWluIGZyb20gb3ZlcnJpZGluZyB0aGUgW1N5bWJvbC5zcGVjaWVzXSBzdGF0aWMgcHJvcGVydHkuXCIpO1xuICAgICAgICB9XG5cbiAgICAgICAgbGV0IHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7XG4gICAgICAgICAgICBwcm9taXNlOiB0aGlzLFxuICAgICAgICAgICAgcmVzb2x2ZSxcbiAgICAgICAgICAgIHJlamVjdCxcbiAgICAgICAgICAgIGdldCBvbmNhbmNlbGxlZCgpIHsgcmV0dXJuIG9uY2FuY2VsbGVkID8/IG51bGw7IH0sXG4gICAgICAgICAgICBzZXQgb25jYW5jZWxsZWQoY2IpIHsgb25jYW5jZWxsZWQgPSBjYiA/PyB1bmRlZmluZWQ7IH1cbiAgICAgICAgfTtcblxuICAgICAgICBjb25zdCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XG4gICAgICAgICAgICBnZXQgcm9vdCgpIHsgcmV0dXJuIHN0YXRlOyB9LFxuICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcbiAgICAgICAgICAgIHNldHRsZWQ6IGZhbHNlXG4gICAgICAgIH07XG5cbiAgICAgICAgLy8gU2V0dXAgY2FuY2VsbGF0aW9uIHN5c3RlbS5cbiAgICAgICAgdm9pZCBPYmplY3QuZGVmaW5lUHJvcGVydGllcyh0aGlzLCB7XG4gICAgICAgICAgICBbYmFycmllclN5bV06IHtcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiB0cnVlLFxuICAgICAgICAgICAgICAgIHZhbHVlOiBudWxsXG4gICAgICAgICAgICB9LFxuICAgICAgICAgICAgW2NhbmNlbEltcGxTeW1dOiB7XG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgdmFsdWU6IGNhbmNlbGxlckZvcihwcm9taXNlLCBzdGF0ZSlcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG5cbiAgICAgICAgLy8gUnVuIHRoZSBhY3R1YWwgZXhlY3V0b3IuXG4gICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgc3RhdGUpO1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgZXhlY3V0b3IocmVzb2x2ZXJGb3IocHJvbWlzZSwgc3RhdGUpLCByZWplY3Rvcik7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykge1xuICAgICAgICAgICAgICAgIGNvbnNvbGUubG9nKFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBDYW5jZWxsYWJsZVByb21pc2UgZXhlY3V0b3IuXCIsIGVycik7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDYW5jZWxzIGltbWVkaWF0ZWx5IHRoZSBleGVjdXRpb24gb2YgdGhlIG9wZXJhdGlvbiBhc3NvY2lhdGVkIHdpdGggdGhpcyBwcm9taXNlLlxuICAgICAqIFRoZSBwcm9taXNlIHJlamVjdHMgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0gaW5zdGFuY2UgYXMgcmVhc29uLFxuICAgICAqIHdpdGggdGhlIHtAbGluayBDYW5jZWxFcnJvciNjYXVzZX0gcHJvcGVydHkgc2V0IHRvIHRoZSBnaXZlbiBhcmd1bWVudCwgaWYgYW55LlxuICAgICAqXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgYWZ0ZXIgdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZDtcbiAgICAgKiByZXBlYXRlZCBjYWxscyBpbiBwYXJ0aWN1bGFyIGFyZSBzYWZlLCBidXQgb25seSB0aGUgZmlyc3Qgb25lXG4gICAgICogd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cbiAgICAgKlxuICAgICAqIFRoZSBgQ2FuY2VsRXJyb3JgIGV4Y2VwdGlvbiBfbmVlZCBub3RfIGJlIGhhbmRsZWQgZXhwbGljaXRseSBfb24gdGhlIHByb21pc2VzIHRoYXQgYXJlIGJlaW5nIGNhbmNlbGxlZDpfXG4gICAgICogY2FuY2VsbGluZyBhIHByb21pc2Ugd2l0aCBubyBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlciBkb2VzIG5vdCB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQuXG4gICAgICogVGhlcmVmb3JlLCB0aGUgZm9sbG93aW5nIGlkaW9tcyBhcmUgYWxsIGVxdWFsbHkgY29ycmVjdDpcbiAgICAgKiBgYGB0c1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkuY2FuY2VsKCk7XG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2FuY2VsKCk7XG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2F0Y2goLi4uKS5jYW5jZWwoKTtcbiAgICAgKiBgYGBcbiAgICAgKiBXaGVuZXZlciBzb21lIGNhbmNlbGxlZCBwcm9taXNlIGluIGEgY2hhaW4gcmVqZWN0cyB3aXRoIGEgYENhbmNlbEVycm9yYFxuICAgICAqIHdpdGggdGhlIHNhbWUgY2FuY2VsbGF0aW9uIGNhdXNlIGFzIGl0c2VsZiwgdGhlIGVycm9yIHdpbGwgYmUgZGlzY2FyZGVkIHNpbGVudGx5LlxuICAgICAqIEhvd2V2ZXIsIHRoZSBgQ2FuY2VsRXJyb3JgIF93aWxsIHN0aWxsIGJlIGRlbGl2ZXJlZF8gdG8gYWxsIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyc1xuICAgICAqIGFkZGVkIGJ5IHtAbGluayB0aGVufSBhbmQgcmVsYXRlZCBtZXRob2RzOlxuICAgICAqIGBgYHRzXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcbiAgICAgKiBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLmNhdGNoKGNvbnNvbGUubG9nKTtcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gQSBDYW5jZWxFcnJvciBpcyBwcmludGVkIHRvIHRoZSBjb25zb2xlLlxuICAgICAqIGBgYFxuICAgICAqIElmIHRoZSBgQ2FuY2VsRXJyb3JgIGlzIG5vdCBoYW5kbGVkIGRvd25zdHJlYW0gYnkgdGhlIHRpbWUgaXQgcmVhY2hlc1xuICAgICAqIGEgX25vbi1jYW5jZWxsZWRfIHByb21pc2UsIGl0IF93aWxsXyB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQsXG4gICAgICoganVzdCBsaWtlIG5vcm1hbCByZWplY3Rpb25zIHdvdWxkOlxuICAgICAqIGBgYHRzXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcbiAgICAgKiBsZXQgY2hhaW5lZCA9IGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkudGhlbigoKSA9PiB7IC4uLiB9KTsgLy8gTm8gY2F0Y2guLi5cbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gVW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCBvbiBjaGFpbmVkIVxuICAgICAqIGBgYFxuICAgICAqIFRoZXJlZm9yZSwgaXQgaXMgaW1wb3J0YW50IHRvIGVpdGhlciBjYW5jZWwgd2hvbGUgcHJvbWlzZSBjaGFpbnMgZnJvbSB0aGVpciB0YWlsLFxuICAgICAqIGFzIHNob3duIGluIHRoZSBjb3JyZWN0IGlkaW9tcyBhYm92ZSwgb3IgdGFrZSBjYXJlIG9mIGhhbmRsaW5nIGVycm9ycyBldmVyeXdoZXJlLlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBjYW5jZWxsYWJsZSBwcm9taXNlIHRoYXQgX2Z1bGZpbGxzXyBhZnRlciB0aGUgY2FuY2VsIGNhbGxiYWNrIChpZiBhbnkpXG4gICAgICogYW5kIGFsbCBoYW5kbGVycyBhdHRhY2hlZCB1cCB0byB0aGUgY2FsbCB0byBjYW5jZWwgaGF2ZSBydW4uXG4gICAgICogSWYgdGhlIGNhbmNlbCBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGJ5IGBjYW5jZWxgXG4gICAgICogd2lsbCBhbHNvIHdhaXQgZm9yIHRoYXQgdGhlbmFibGUgdG8gc2V0dGxlLlxuICAgICAqIFRoaXMgZW5hYmxlcyBjYWxsZXJzIHRvIHdhaXQgZm9yIHRoZSBjYW5jZWxsZWQgb3BlcmF0aW9uIHRvIHRlcm1pbmF0ZVxuICAgICAqIHdpdGhvdXQgYmVpbmcgZm9yY2VkIHRvIGhhbmRsZSBwb3RlbnRpYWwgZXJyb3JzIGF0IHRoZSBjYWxsIHNpdGUuXG4gICAgICogYGBgdHNcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKS50aGVuKCgpID0+IHtcbiAgICAgKiAgICAgLy8gQ2xlYW51cCBmaW5pc2hlZCwgaXQncyBzYWZlIHRvIGRvIHNvbWV0aGluZyBlbHNlLlxuICAgICAqIH0sIChlcnIpID0+IHtcbiAgICAgKiAgICAgLy8gVW5yZWFjaGFibGU6IHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gY2FuY2VsIHdpbGwgbmV2ZXIgcmVqZWN0LlxuICAgICAqIH0pO1xuICAgICAqIGBgYFxuICAgICAqIE5vdGUgdGhhdCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIF9ub3RfIGhhbmRsZSBpbXBsaWNpdGx5IGFueSByZWplY3Rpb25cbiAgICAgKiB0aGF0IG1pZ2h0IGhhdmUgb2NjdXJyZWQgYWxyZWFkeSBpbiB0aGUgY2FuY2VsbGVkIGNoYWluLlxuICAgICAqIEl0IHdpbGwganVzdCB0cmFjayB3aGV0aGVyIHJlZ2lzdGVyZWQgaGFuZGxlcnMgaGF2ZSBiZWVuIGV4ZWN1dGVkIG9yIG5vdC5cbiAgICAgKiBUaGVyZWZvcmUsIHVuaGFuZGxlZCByZWplY3Rpb25zIHdpbGwgbmV2ZXIgYmUgc2lsZW50bHkgaGFuZGxlZCBieSBjYWxsaW5nIGNhbmNlbC5cbiAgICAgKi9cbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgLy8gSU5WQVJJQU5UOiB0aGUgcmVzdWx0IG9mIHRoaXNbY2FuY2VsSW1wbFN5bV0gYW5kIHRoZSBiYXJyaWVyIGRvIG5vdCBldmVyIHJlamVjdC5cbiAgICAgICAgICAgIC8vIFVuZm9ydHVuYXRlbHkgbWFjT1MgSGlnaCBTaWVycmEgZG9lcyBub3Qgc3VwcG9ydCBQcm9taXNlLmFsbFNldHRsZWQuXG4gICAgICAgICAgICBQcm9taXNlLmFsbChbXG4gICAgICAgICAgICAgICAgdGhpc1tjYW5jZWxJbXBsU3ltXShuZXcgQ2FuY2VsRXJyb3IoXCJQcm9taXNlIGNhbmNlbGxlZC5cIiwgeyBjYXVzZSB9KSksXG4gICAgICAgICAgICAgICAgY3VycmVudEJhcnJpZXIodGhpcylcbiAgICAgICAgICAgIF0pLnRoZW4oKCkgPT4gcmVzb2x2ZSgpLCAoKSA9PiByZXNvbHZlKCkpO1xuICAgICAgICB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBCaW5kcyBwcm9taXNlIGNhbmNlbGxhdGlvbiB0byB0aGUgYWJvcnQgZXZlbnQgb2YgdGhlIGdpdmVuIHtAbGluayBBYm9ydFNpZ25hbH0uXG4gICAgICogSWYgdGhlIHNpZ25hbCBoYXMgYWxyZWFkeSBhYm9ydGVkLCB0aGUgcHJvbWlzZSB3aWxsIGJlIGNhbmNlbGxlZCBpbW1lZGlhdGVseS5cbiAgICAgKiBXaGVuIGVpdGhlciBjb25kaXRpb24gaXMgdmVyaWZpZWQsIHRoZSBjYW5jZWxsYXRpb24gY2F1c2Ugd2lsbCBiZSBzZXRcbiAgICAgKiB0byB0aGUgc2lnbmFsJ3MgYWJvcnQgcmVhc29uIChzZWUge0BsaW5rIEFib3J0U2lnbmFsI3JlYXNvbn0pLlxuICAgICAqXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgKG9yIGlmIHRoZSBzaWduYWwgYWJvcnRzKSBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQuXG4gICAgICogT25seSB0aGUgZmlyc3Qgc2lnbmFsIHRvIGFib3J0IHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXG4gICAgICpcbiAgICAgKiBGb3IgbW9yZSBkZXRhaWxzIGFib3V0IHRoZSBjYW5jZWxsYXRpb24gcHJvY2VzcyxcbiAgICAgKiBzZWUge0BsaW5rIGNhbmNlbH0gYW5kIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYCBjb25zdHJ1Y3Rvci5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGVuYWJsZXMgYGF3YWl0YGluZyBjYW5jZWxsYWJsZSBwcm9taXNlcyB3aXRob3V0IGhhdmluZ1xuICAgICAqIHRvIHN0b3JlIHRoZW0gZm9yIGZ1dHVyZSBjYW5jZWxsYXRpb24sIGUuZy46XG4gICAgICogYGBgdHNcbiAgICAgKiBhd2FpdCBsb25nUnVubmluZ09wZXJhdGlvbigpLmNhbmNlbE9uKHNpZ25hbCk7XG4gICAgICogYGBgXG4gICAgICogaW5zdGVhZCBvZjpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBwcm9taXNlVG9CZUNhbmNlbGxlZCA9IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCk7XG4gICAgICogYXdhaXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQ7XG4gICAgICogYGBgXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGlzIHByb21pc2UsIGZvciBtZXRob2QgY2hhaW5pbmcuXG4gICAgICovXG4gICAgY2FuY2VsT24oc2lnbmFsOiBBYm9ydFNpZ25hbCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGlmIChzaWduYWwuYWJvcnRlZCkge1xuICAgICAgICAgICAgdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKVxuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc2lnbmFsLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKSwge2NhcHR1cmU6IHRydWV9KTtcbiAgICAgICAgfVxuXG4gICAgICAgIHJldHVybiB0aGlzO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEF0dGFjaGVzIGNhbGxiYWNrcyBmb3IgdGhlIHJlc29sdXRpb24gYW5kL29yIHJlamVjdGlvbiBvZiB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqXG4gICAgICogQHBhcmFtIG9uZnVsZmlsbGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVzb2x2ZWQuXG4gICAgICogQHBhcmFtIG9ucmVqZWN0ZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZWplY3RlZC5cbiAgICAgKiBAcmV0dXJucyBBIGBDYW5jZWxsYWJsZVByb21pc2VgIGZvciB0aGUgY29tcGxldGlvbiBvZiB3aGljaGV2ZXIgY2FsbGJhY2sgaXMgZXhlY3V0ZWQuXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgaXMgaG9va2VkIHVwIHRvIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24gcmVxdWVzdHMgdXAgdGhlIGNoYWluLCBidXQgbm90IGRvd246XG4gICAgICpcbiAgICAgKiAgIC0gaWYgdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgYmUgaW52b2tlZCB3aXRoIGEgYENhbmNlbEVycm9yYFxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2UgX3dpbGwgcmVzb2x2ZSByZWd1bGFybHlfIHdpdGggaXRzIHJlc3VsdDtcbiAgICAgKiAgIC0gY29udmVyc2VseSwgaWYgdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCBfdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCB0b287X1xuICAgICAqICAgICB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBzdGlsbCBiZSBpbnZva2VkIHdpdGggdGhlIHBhcmVudCdzIGBDYW5jZWxFcnJvcmAsXG4gICAgICogICAgIGJ1dCBpdHMgcmVzdWx0IHdpbGwgYmUgZGlzY2FyZGVkXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGEgYENhbmNlbEVycm9yYCBhcyB3ZWxsLlxuICAgICAqXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cbiAgICAgKlxuICAgICAqIElmIGVpdGhlciBjYWxsYmFjayByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqL1xuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4ge1xuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUudGhlbiBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gTk9URTogVHlwZVNjcmlwdCdzIGJ1aWx0LWluIHR5cGUgZm9yIHRoZW4gaXMgYnJva2VuLFxuICAgICAgICAvLyBhcyBpdCBhbGxvd3Mgc3BlY2lmeWluZyBhbiBhcmJpdHJhcnkgVFJlc3VsdDEgIT0gVCBldmVuIHdoZW4gb25mdWxmaWxsZWQgaXMgbm90IGEgZnVuY3Rpb24uXG4gICAgICAgIC8vIFdlIGNhbm5vdCBmaXggaXQgaWYgd2Ugd2FudCB0byBDYW5jZWxsYWJsZVByb21pc2UgdG8gaW1wbGVtZW50IFByb21pc2VMaWtlPFQ+LlxuXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZ1bGZpbGxlZCkpIHsgb25mdWxmaWxsZWQgPSBpZGVudGl0eSBhcyBhbnk7IH1cbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9ucmVqZWN0ZWQpKSB7IG9ucmVqZWN0ZWQgPSB0aHJvd2VyOyB9XG5cbiAgICAgICAgaWYgKG9uZnVsZmlsbGVkID09PSBpZGVudGl0eSAmJiBvbnJlamVjdGVkID09IHRocm93ZXIpIHtcbiAgICAgICAgICAgIC8vIFNob3J0Y3V0IGZvciB0cml2aWFsIGFyZ3VtZW50cy5cbiAgICAgICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlKSA9PiByZXNvbHZlKHRoaXMgYXMgYW55KSk7XG4gICAgICAgIH1cblxuICAgICAgICBjb25zdCBiYXJyaWVyOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiA9IHt9O1xuICAgICAgICB0aGlzW2JhcnJpZXJTeW1dID0gYmFycmllcjtcblxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUUmVzdWx0MSB8IFRSZXN1bHQyPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICB2b2lkIHN1cGVyLnRoZW4oXG4gICAgICAgICAgICAgICAgKHZhbHVlKSA9PiB7XG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XG5cbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25mdWxmaWxsZWQhKHZhbHVlKSk7XG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9LFxuICAgICAgICAgICAgICAgIChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XG5cbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25yZWplY3RlZCEocmVhc29uKSk7XG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICApO1xuICAgICAgICB9LCBhc3luYyAoY2F1c2U/KSA9PiB7XG4gICAgICAgICAgICAvL2NhbmNlbGxlZCA9IHRydWU7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHJldHVybiBvbmNhbmNlbGxlZD8uKGNhdXNlKTtcbiAgICAgICAgICAgIH0gZmluYWxseSB7XG4gICAgICAgICAgICAgICAgYXdhaXQgdGhpcy5jYW5jZWwoY2F1c2UpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIGZvciBvbmx5IHRoZSByZWplY3Rpb24gb2YgdGhlIFByb21pc2UuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqXG4gICAgICogSXQgaXMgZXF1aXZhbGVudCB0b1xuICAgICAqIGBgYHRzXG4gICAgICogY2FuY2VsbGFibGVQcm9taXNlLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XG4gICAgICogYGBgXG4gICAgICogYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBgb25yZWplY3RlZGAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXG4gICAgICovXG4gICAgY2F0Y2g8VFJlc3VsdCA9IG5ldmVyPihvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gKFByb21pc2VMaWtlPFRSZXN1bHQ+IHwgVFJlc3VsdCkpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VCB8IFRSZXN1bHQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIHRoYXQgaXMgaW52b2tlZCB3aGVuIHRoZSBDYW5jZWxsYWJsZVByb21pc2UgaXMgc2V0dGxlZCAoZnVsZmlsbGVkIG9yIHJlamVjdGVkKS4gVGhlXG4gICAgICogcmVzb2x2ZWQgdmFsdWUgY2Fubm90IGJlIGFjY2Vzc2VkIG9yIG1vZGlmaWVkIGZyb20gdGhlIGNhbGxiYWNrLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgc2V0dGxlIGluIHRoZSBzYW1lIHN0YXRlIGFzIHRoZSBvcmlnaW5hbCBvbmVcbiAgICAgKiBhZnRlciB0aGUgcHJvdmlkZWQgY2FsbGJhY2sgaGFzIGNvbXBsZXRlZCBleGVjdXRpb24sXG4gICAgICogdW5sZXNzIHRoZSBjYWxsYmFjayB0aHJvd3Mgb3IgcmV0dXJucyBhIHJlamVjdGluZyBwcm9taXNlLFxuICAgICAqIGluIHdoaWNoIGNhc2UgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3QgYXMgd2VsbC5cbiAgICAgKlxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXG4gICAgICogT25jZSB0aGUgcGFyZW50IHByb21pc2Ugc2V0dGxlcywgdGhlIGBvbmZpbmFsbHlgIGNhbGxiYWNrIHdpbGwgcnVuLFxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqXG4gICAgICogVGhpcyBtZXRob2QgaXMgaW1wbGVtZW50ZWQgaW4gdGVybXMgb2Yge0BsaW5rIHRoZW59IGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxuICAgICAqIEl0IGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBgb25maW5hbGx5YCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cbiAgICAgKi9cbiAgICBmaW5hbGx5KG9uZmluYWxseT86ICgoKSA9PiB2b2lkKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLmZpbmFsbHkgY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZpbmFsbHkpKSB7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy50aGVuKG9uZmluYWxseSwgb25maW5hbGx5LCBvbmNhbmNlbGxlZCk7XG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gdGhpcy50aGVuKFxuICAgICAgICAgICAgKHZhbHVlKSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB2YWx1ZSksXG4gICAgICAgICAgICAocmVhc29uPykgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4geyB0aHJvdyByZWFzb247IH0pLFxuICAgICAgICAgICAgb25jYW5jZWxsZWQsXG4gICAgICAgICk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogV2UgdXNlIHRoZSBgW1N5bWJvbC5zcGVjaWVzXWAgc3RhdGljIHByb3BlcnR5LCBpZiBhdmFpbGFibGUsXG4gICAgICogdG8gZGlzYWJsZSB0aGUgYnVpbHQtaW4gYXV0b21hdGljIHN1YmNsYXNzaW5nIGZlYXR1cmVzIGZyb20ge0BsaW5rIFByb21pc2V9LlxuICAgICAqIEl0IGlzIGNyaXRpY2FsIGZvciBwZXJmb3JtYW5jZSByZWFzb25zIHRoYXQgZXh0ZW5kZXJzIGRvIG5vdCBvdmVycmlkZSB0aGlzLlxuICAgICAqIE9uY2UgdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmdcbiAgICAgKiBpcyBlaXRoZXIgYWNjZXB0ZWQgb3IgcmV0aXJlZCwgdGhpcyBpbXBsZW1lbnRhdGlvbiB3aWxsIGhhdmUgdG8gYmUgcmV2aXNlZCBhY2NvcmRpbmdseS5cbiAgICAgKlxuICAgICAqIEBpZ25vcmVcbiAgICAgKiBAaW50ZXJuYWxcbiAgICAgKi9cbiAgICBzdGF0aWMgZ2V0IFtzcGVjaWVzXSgpIHtcbiAgICAgICAgcmV0dXJuIFByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlLCBvciByZWplY3RlZCB3aGVuIGFueSBQcm9taXNlIGlzIHJlamVjdGVkLlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbGw8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD5bXT47XG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogQXdhaXRlZDxUW1BdPjsgfT47XG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSBvciByZWplY3QuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VD4+W10+O1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFRbUF0+PjsgfT47XG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGxTZXR0bGVkKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVGhlIGFueSBmdW5jdGlvbiByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGlzIGZ1bGZpbGxlZCBieSB0aGUgZmlyc3QgZ2l2ZW4gcHJvbWlzZSB0byBiZSBmdWxmaWxsZWQsXG4gICAgICogb3IgcmVqZWN0ZWQgd2l0aCBhbiBBZ2dyZWdhdGVFcnJvciBjb250YWluaW5nIGFuIGFycmF5IG9mIHJlamVjdGlvbiByZWFzb25zXG4gICAgICogaWYgYWxsIG9mIHRoZSBnaXZlbiBwcm9taXNlcyBhcmUgcmVqZWN0ZWQuXG4gICAgICogSXQgcmVzb2x2ZXMgYWxsIGVsZW1lbnRzIG9mIHRoZSBwYXNzZWQgaXRlcmFibGUgdG8gcHJvbWlzZXMgYXMgaXQgcnVucyB0aGlzIGFsZ29yaXRobS5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYW55PFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFueShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIGFyZSByZXNvbHZlZCBvciByZWplY3RlZC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmFjZTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUW251bWJlcl0+PjtcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yYWNlKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGNhbmNlbGxlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBjYW5jZWw8VCA9IG5ldmVyPihjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGNvbnN0IHAgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcbiAgICAgICAgcC5jYW5jZWwoY2F1c2UpO1xuICAgICAgICByZXR1cm4gcDtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGNhbmNlbHNcbiAgICAgKiBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIGNhdXNlLlxuICAgICAqXG4gICAgICogSWYgdGhlIHtAbGluayBBYm9ydFNpZ25hbC50aW1lb3V0fSBmYWN0b3J5IG1ldGhvZCBpcyBhdmFpbGFibGUsXG4gICAgICogaXQgaXMgdXNlZCB0byBiYXNlIHRoZSB0aW1lb3V0IG9uIF9hY3RpdmVfIHRpbWUgcmF0aGVyIHRoYW4gX2VsYXBzZWRfIHRpbWUuXG4gICAgICogT3RoZXJ3aXNlLCBgdGltZW91dGAgZmFsbHMgYmFjayB0byB7QGxpbmsgc2V0VGltZW91dH0uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgdGltZW91dDxUID0gbmV2ZXI+KG1pbGxpc2Vjb25kczogbnVtYmVyLCBjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcbiAgICAgICAgaWYgKEFib3J0U2lnbmFsICYmIHR5cGVvZiBBYm9ydFNpZ25hbCA9PT0gJ2Z1bmN0aW9uJyAmJiBBYm9ydFNpZ25hbC50aW1lb3V0ICYmIHR5cGVvZiBBYm9ydFNpZ25hbC50aW1lb3V0ID09PSAnZnVuY3Rpb24nKSB7XG4gICAgICAgICAgICBBYm9ydFNpZ25hbC50aW1lb3V0KG1pbGxpc2Vjb25kcykuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpLCBtaWxsaXNlY29uZHMpO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgc2xlZXAobWlsbGlzZWNvbmRzOiBudW1iZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlclxuICAgICAqIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBzbGVlcDxUPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XG4gICAgc3RhdGljIHNsZWVwPFQgPSB2b2lkPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU/OiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gcmVzb2x2ZSh2YWx1ZSEpLCBtaWxsaXNlY29uZHMpO1xuICAgICAgICB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlamVjdGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHJlYXNvbi5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZWplY3Q8VCA9IG5ldmVyPihyZWFzb24/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoXywgcmVqZWN0KSA9PiByZWplY3QocmVhc29uKSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZSgpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyByZXNvbHZlPFQgPSB2b2lkPih2YWx1ZT86IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PiB7XG4gICAgICAgIGlmICh2YWx1ZSBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkge1xuICAgICAgICAgICAgLy8gT3B0aW1pc2UgZm9yIGNhbmNlbGxhYmxlIHByb21pc2VzLlxuICAgICAgICAgICAgcmV0dXJuIHZhbHVlO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPGFueT4oKHJlc29sdmUpID0+IHJlc29sdmUodmFsdWUpKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSBhbmQgcmV0dXJucyBpdCBpbiBhbiBvYmplY3QsIGFsb25nIHdpdGggaXRzIHJlc29sdmUgYW5kIHJlamVjdCBmdW5jdGlvbnNcbiAgICAgKiBhbmQgYSBnZXR0ZXIvc2V0dGVyIGZvciB0aGUgY2FuY2VsbGF0aW9uIGNhbGxiYWNrLlxuICAgICAqXG4gICAgICogVGhpcyBtZXRob2QgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyB3aXRoUmVzb2x2ZXJzPFQ+KCk6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xuICAgICAgICBsZXQgcmVzdWx0OiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0geyBvbmNhbmNlbGxlZDogbnVsbCB9IGFzIGFueTtcbiAgICAgICAgcmVzdWx0LnByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgIHJlc3VsdC5yZXNvbHZlID0gcmVzb2x2ZTtcbiAgICAgICAgICAgIHJlc3VsdC5yZWplY3QgPSByZWplY3Q7XG4gICAgICAgIH0sIChjYXVzZT86IGFueSkgPT4geyByZXN1bHQub25jYW5jZWxsZWQ/LihjYXVzZSk7IH0pO1xuICAgICAgICByZXR1cm4gcmVzdWx0O1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSBjYW5jZWxsYXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cbiAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIHJlc3VsdGluZyBmdW5jdGlvbiBkb2VzIG5vdCByZWplY3QuXG4gKi9cbmZ1bmN0aW9uIGNhbmNlbGxlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpIHtcbiAgICBsZXQgY2FuY2VsbGF0aW9uUHJvbWlzZTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0gdW5kZWZpbmVkO1xuXG4gICAgcmV0dXJuIChyZWFzb246IENhbmNlbEVycm9yKTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0+IHtcbiAgICAgICAgaWYgKCFzdGF0ZS5zZXR0bGVkKSB7XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHN0YXRlLnJlYXNvbiA9IHJlYXNvbjtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XG5cbiAgICAgICAgICAgIC8vIEF0dGFjaCBhbiBlcnJvciBoYW5kbGVyIHRoYXQgaWdub3JlcyB0aGlzIHNwZWNpZmljIHJlamVjdGlvbiByZWFzb24gYW5kIG5vdGhpbmcgZWxzZS5cbiAgICAgICAgICAgIC8vIEluIHRoZW9yeSwgYSBzYW5lIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24gYXQgdGhpcyBwb2ludFxuICAgICAgICAgICAgLy8gc2hvdWxkIGFsd2F5cyByZWplY3Qgd2l0aCBvdXIgY2FuY2VsbGF0aW9uIHJlYXNvbixcbiAgICAgICAgICAgIC8vIGhlbmNlIHRoZSBoYW5kbGVyIHdpbGwgbmV2ZXIgdGhyb3cuXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucHJvdG90eXBlLnRoZW4uY2FsbChwcm9taXNlLnByb21pc2UsIHVuZGVmaW5lZCwgKGVycikgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChlcnIgIT09IHJlYXNvbikge1xuICAgICAgICAgICAgICAgICAgICB0aHJvdyBlcnI7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBJZiByZWFzb24gaXMgbm90IHNldCwgdGhlIHByb21pc2UgcmVzb2x2ZWQgcmVndWxhcmx5LCBoZW5jZSB3ZSBtdXN0IG5vdCBjYWxsIG9uY2FuY2VsbGVkLlxuICAgICAgICAvLyBJZiBvbmNhbmNlbGxlZCBpcyB1bnNldCwgbm8gbmVlZCB0byBnbyBhbnkgZnVydGhlci5cbiAgICAgICAgaWYgKCFzdGF0ZS5yZWFzb24gfHwgIXByb21pc2Uub25jYW5jZWxsZWQpIHsgcmV0dXJuOyB9XG5cbiAgICAgICAgY2FuY2VsbGF0aW9uUHJvbWlzZSA9IG5ldyBQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHJlc29sdmUocHJvbWlzZS5vbmNhbmNlbGxlZCEoc3RhdGUucmVhc29uIS5jYXVzZSkpO1xuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XG4gICAgICAgIH0pO1xuXG4gICAgICAgIC8vIFVuc2V0IG9uY2FuY2VsbGVkIHRvIHByZXZlbnQgcmVwZWF0ZWQgY2FsbHMuXG4gICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBudWxsO1xuXG4gICAgICAgIHJldHVybiBjYW5jZWxsYXRpb25Qcm9taXNlO1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZXNvbHV0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKi9cbmZ1bmN0aW9uIHJlc29sdmVyRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+IHtcbiAgICByZXR1cm4gKHZhbHVlKSA9PiB7XG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnJlc29sdmluZyA9IHRydWU7XG5cbiAgICAgICAgaWYgKHZhbHVlID09PSBwcm9taXNlLnByb21pc2UpIHtcbiAgICAgICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7IHJldHVybjsgfVxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgVHlwZUVycm9yKFwiQSBwcm9taXNlIGNhbm5vdCBiZSByZXNvbHZlZCB3aXRoIGl0c2VsZi5cIikpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHZhbHVlICE9IG51bGwgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcgfHwgdHlwZW9mIHZhbHVlID09PSAnZnVuY3Rpb24nKSkge1xuICAgICAgICAgICAgbGV0IHRoZW46IGFueTtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgdGhlbiA9ICh2YWx1ZSBhcyBhbnkpLnRoZW47XG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgICAgICBwcm9taXNlLnJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUodGhlbikpIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICBsZXQgY2FuY2VsID0gKHZhbHVlIGFzIGFueSkuY2FuY2VsO1xuICAgICAgICAgICAgICAgICAgICBpZiAoaXNDYWxsYWJsZShjYW5jZWwpKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCBvbmNhbmNlbGxlZCA9IChjYXVzZT86IGFueSkgPT4ge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XG4gICAgICAgICAgICAgICAgICAgICAgICB9O1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKHN0YXRlLnJlYXNvbikge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIElmIGFscmVhZHkgY2FuY2VsbGVkLCBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uLlxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIGNhbmNlbGxlciBhbGdvcml0aG0gZG9lcyBub3QgcmVqZWN0XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gc28gaXQgY2FuIGJlIGRpc2NhcmRlZCBzYWZlbHkuXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCBjYW5jZWxsZXJGb3IoeyAuLi5wcm9taXNlLCBvbmNhbmNlbGxlZCB9LCBzdGF0ZSkoc3RhdGUucmVhc29uKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG9uY2FuY2VsbGVkO1xuICAgICAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSBjYXRjaCB7fVxuXG4gICAgICAgICAgICAgICAgY29uc3QgbmV3U3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xuICAgICAgICAgICAgICAgICAgICByb290OiBzdGF0ZS5yb290LFxuICAgICAgICAgICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxuICAgICAgICAgICAgICAgICAgICBnZXQgc2V0dGxlZCgpIHsgcmV0dXJuIHRoaXMucm9vdC5zZXR0bGVkIH0sXG4gICAgICAgICAgICAgICAgICAgIHNldCBzZXR0bGVkKHZhbHVlKSB7IHRoaXMucm9vdC5zZXR0bGVkID0gdmFsdWU7IH0sXG4gICAgICAgICAgICAgICAgICAgIGdldCByZWFzb24oKSB7IHJldHVybiB0aGlzLnJvb3QucmVhc29uIH1cbiAgICAgICAgICAgICAgICB9O1xuXG4gICAgICAgICAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBuZXdTdGF0ZSk7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseSh0aGVuLCB2YWx1ZSwgW3Jlc29sdmVyRm9yKHByb21pc2UsIG5ld1N0YXRlKSwgcmVqZWN0b3JdKTtcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmV0dXJuOyAvLyBJTVBPUlRBTlQhXG4gICAgICAgICAgICB9XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgIHByb21pc2UucmVzb2x2ZSh2YWx1ZSk7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZWplY3Rpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cbiAqL1xuZnVuY3Rpb24gcmVqZWN0b3JGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3Ige1xuICAgIHJldHVybiAocmVhc29uPykgPT4ge1xuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xuXG4gICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGlmIChyZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBzdGF0ZS5yZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBPYmplY3QuaXMocmVhc29uLmNhdXNlLCBzdGF0ZS5yZWFzb24uY2F1c2UpKSB7XG4gICAgICAgICAgICAgICAgICAgIC8vIFN3YWxsb3cgbGF0ZSByZWplY3Rpb25zIHRoYXQgYXJlIENhbmNlbEVycm9ycyB3aG9zZSBjYW5jZWxsYXRpb24gY2F1c2UgaXMgdGhlIHNhbWUgYXMgb3Vycy5cbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0gY2F0Y2gge31cblxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24pKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBDYW5jZWxzIGFsbCB2YWx1ZXMgaW4gYW4gYXJyYXkgdGhhdCBsb29rIGxpa2UgY2FuY2VsbGFibGUgdGhlbmFibGVzLlxuICogUmV0dXJucyBhIHByb21pc2UgdGhhdCBmdWxmaWxscyBvbmNlIGFsbCBjYW5jZWxsYXRpb24gcHJvY2VkdXJlcyBmb3IgdGhlIGdpdmVuIHZhbHVlcyBoYXZlIHNldHRsZWQuXG4gKi9cbmZ1bmN0aW9uIGNhbmNlbEFsbChwYXJlbnQ6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgdmFsdWVzOiBhbnlbXSwgY2F1c2U/OiBhbnkpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBjb25zdCByZXN1bHRzID0gW107XG5cbiAgICBmb3IgKGNvbnN0IHZhbHVlIG9mIHZhbHVlcykge1xuICAgICAgICBsZXQgY2FuY2VsOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXI7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUodmFsdWUudGhlbikpIHsgY29udGludWU7IH1cbiAgICAgICAgICAgIGNhbmNlbCA9IHZhbHVlLmNhbmNlbDtcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZShjYW5jZWwpKSB7IGNvbnRpbnVlOyB9XG4gICAgICAgIH0gY2F0Y2ggeyBjb250aW51ZTsgfVxuXG4gICAgICAgIGxldCByZXN1bHQ6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHJlc3VsdCA9IFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xuICAgICAgICAgICAgY29udGludWU7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIXJlc3VsdCkgeyBjb250aW51ZTsgfVxuICAgICAgICByZXN1bHRzLnB1c2goXG4gICAgICAgICAgICAocmVzdWx0IGluc3RhbmNlb2YgUHJvbWlzZSAgPyByZXN1bHQgOiBQcm9taXNlLnJlc29sdmUocmVzdWx0KSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICApO1xuICAgIH1cblxuICAgIHJldHVybiBQcm9taXNlLmFsbChyZXN1bHRzKSBhcyBhbnk7XG59XG5cbi8qKlxuICogUmV0dXJucyBpdHMgYXJndW1lbnQuXG4gKi9cbmZ1bmN0aW9uIGlkZW50aXR5PFQ+KHg6IFQpOiBUIHtcbiAgICByZXR1cm4geDtcbn1cblxuLyoqXG4gKiBUaHJvd3MgaXRzIGFyZ3VtZW50LlxuICovXG5mdW5jdGlvbiB0aHJvd2VyKHJlYXNvbj86IGFueSk6IG5ldmVyIHtcbiAgICB0aHJvdyByZWFzb247XG59XG5cbi8qKlxuICogQXR0ZW1wdHMgdmFyaW91cyBzdHJhdGVnaWVzIHRvIGNvbnZlcnQgYW4gZXJyb3IgdG8gYSBzdHJpbmcuXG4gKi9cbmZ1bmN0aW9uIGVycm9yTWVzc2FnZShlcnI6IGFueSk6IHN0cmluZyB7XG4gICAgdHJ5IHtcbiAgICAgICAgaWYgKGVyciBpbnN0YW5jZW9mIEVycm9yIHx8IHR5cGVvZiBlcnIgIT09ICdvYmplY3QnIHx8IGVyci50b1N0cmluZyAhPT0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZykge1xuICAgICAgICAgICAgcmV0dXJuIFwiXCIgKyBlcnI7XG4gICAgICAgIH1cbiAgICB9IGNhdGNoIHt9XG5cbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gSlNPTi5zdHJpbmdpZnkoZXJyKTtcbiAgICB9IGNhdGNoIHt9XG5cbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZy5jYWxsKGVycik7XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgcmV0dXJuIFwiPGNvdWxkIG5vdCBjb252ZXJ0IGVycm9yIHRvIHN0cmluZz5cIjtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGJhcnJpZXIgcHJvbWlzZSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuIElmIG5lY2Vzc2FyeSwgaW5pdGlhbGlzZXMgdGhlIGJhcnJpZXIuXG4gKi9cbmZ1bmN0aW9uIGN1cnJlbnRCYXJyaWVyPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPik6IFByb21pc2U8dm9pZD4ge1xuICAgIGxldCBwd3I6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0gcHJvbWlzZVtiYXJyaWVyU3ltXSA/PyB7fTtcbiAgICBpZiAoISgncHJvbWlzZScgaW4gcHdyKSkge1xuICAgICAgICBPYmplY3QuYXNzaWduKHB3ciwgcHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4oKSk7XG4gICAgfVxuICAgIGlmIChwcm9taXNlW2JhcnJpZXJTeW1dID09IG51bGwpIHtcbiAgICAgICAgcHdyLnJlc29sdmUhKCk7XG4gICAgICAgIHByb21pc2VbYmFycmllclN5bV0gPSBwd3I7XG4gICAgfVxuICAgIHJldHVybiBwd3IucHJvbWlzZSE7XG59XG5cbi8vIFBvbHlmaWxsIFByb21pc2Uud2l0aFJlc29sdmVycy5cbmxldCBwcm9taXNlV2l0aFJlc29sdmVycyA9IFByb21pc2Uud2l0aFJlc29sdmVycztcbmlmIChwcm9taXNlV2l0aFJlc29sdmVycyAmJiB0eXBlb2YgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPT09ICdmdW5jdGlvbicpIHtcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IHByb21pc2VXaXRoUmVzb2x2ZXJzLmJpbmQoUHJvbWlzZSk7XG59IGVsc2Uge1xuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gZnVuY3Rpb24gPFQ+KCk6IFByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgUHJvbWlzZTxUPigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcbiAgICAgICAgcmV0dXJuIHsgcHJvbWlzZSwgcmVzb2x2ZSwgcmVqZWN0IH07XG4gICAgfVxufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DbGlwYm9hcmQpO1xuXG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcbmNvbnN0IENsaXBib2FyZFRleHQgPSAxO1xuXG4vKipcbiAqIFNldHMgdGhlIHRleHQgdG8gdGhlIENsaXBib2FyZC5cbiAqXG4gKiBAcGFyYW0gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxuICogQHJldHVybiBBIFByb21pc2UgdGhhdCByZXNvbHZlcyB3aGVuIHRoZSBvcGVyYXRpb24gaXMgc3VjY2Vzc2Z1bC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dDogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBUZXh0KCk6IFByb21pc2U8c3RyaW5nPiB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkVGV4dCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogQW55IGlzIGEgZHVtbXkgY3JlYXRpb24gZnVuY3Rpb24gZm9yIHNpbXBsZSBvciB1bmtub3duIHR5cGVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQW55PFQgPSBhbnk+KHNvdXJjZTogYW55KTogVCB7XG4gICAgcmV0dXJuIHNvdXJjZTtcbn1cblxuLyoqXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXG4gKiBudWxsIHN0cmluZ3Mgd2l0aCBlbXB0eSBzdHJpbmdzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnl0ZVNsaWNlKHNvdXJjZTogYW55KTogc3RyaW5nIHtcbiAgICByZXR1cm4gKChzb3VyY2UgPT0gbnVsbCkgPyBcIlwiIDogc291cmNlKTtcbn1cblxuLyoqXG4gKiBBcnJheSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XG4gKiB3aG9zZSBlbGVtZW50cyBhcmUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQXJyYXk8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiBUW10ge1xuICAgIGlmIChlbGVtZW50ID09PSBBbnkpIHtcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBbXSA6IHNvdXJjZSk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xuICAgICAgICAgICAgcmV0dXJuIFtdO1xuICAgICAgICB9XG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgICAgICBzb3VyY2VbaV0gPSBlbGVtZW50KHNvdXJjZVtpXSk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcbiAgICB9O1xufVxuXG4vKipcbiAqIE1hcCB0YWtlcyBjcmVhdGlvbiBmdW5jdGlvbnMgZm9yIHR3byBhcmJpdHJhcnkgdHlwZXNcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcbiAqIHdob3NlIGtleXMgYW5kIHZhbHVlcyBhcmUgb2YgdGhvc2UgdHlwZXMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBNYXA8ViA9IGFueT4oa2V5OiAoc291cmNlOiBhbnkpID0+IHN0cmluZywgdmFsdWU6IChzb3VyY2U6IGFueSkgPT4gVik6IChzb3VyY2U6IGFueSkgPT4gUmVjb3JkPHN0cmluZywgVj4ge1xuICAgIGlmICh2YWx1ZSA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8ge30gOiBzb3VyY2UpO1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybiB7fTtcbiAgICAgICAgfVxuICAgICAgICBmb3IgKGNvbnN0IGtleSBpbiBzb3VyY2UpIHtcbiAgICAgICAgICAgIHNvdXJjZVtrZXldID0gdmFsdWUoc291cmNlW2tleV0pO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBOdWxsYWJsZSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxuICogYW5kIHJldHVybnMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBudWxsYWJsZSB2YWx1ZSBvZiB0aGF0IHR5cGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBOdWxsYWJsZTxUID0gYW55PihlbGVtZW50OiAoc291cmNlOiBhbnkpID0+IFQpOiAoc291cmNlOiBhbnkpID0+IChUIHwgbnVsbCkge1xuICAgIGlmIChlbGVtZW50ID09PSBBbnkpIHtcbiAgICAgICAgcmV0dXJuIEFueTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IG51bGwgOiBlbGVtZW50KHNvdXJjZSkpO1xufVxuXG4vKipcbiAqIFN0cnVjdCB0YWtlcyBhbiBvYmplY3QgbWFwcGluZyBmaWVsZCBuYW1lcyB0byBjcmVhdGlvbiBmdW5jdGlvbnNcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhIHN0cnVjdC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFN0cnVjdChjcmVhdGVGaWVsZDogUmVjb3JkPHN0cmluZywgKHNvdXJjZTogYW55KSA9PiBhbnk+KTpcbiAgICA8VSBleHRlbmRzIFJlY29yZDxzdHJpbmcsIGFueT4gPSBhbnk+KHNvdXJjZTogYW55KSA9PiBVXG57XG4gICAgbGV0IGFsbEFueSA9IHRydWU7XG4gICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XG4gICAgICAgIGlmIChjcmVhdGVGaWVsZFtuYW1lXSAhPT0gQW55KSB7XG4gICAgICAgICAgICBhbGxBbnkgPSBmYWxzZTtcbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmIChhbGxBbnkpIHtcbiAgICAgICAgcmV0dXJuIEFueTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcbiAgICAgICAgICAgIGlmIChuYW1lIGluIHNvdXJjZSkge1xuICAgICAgICAgICAgICAgIHNvdXJjZVtuYW1lXSA9IGNyZWF0ZUZpZWxkW25hbWVdKHNvdXJjZVtuYW1lXSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcbiAgICB9O1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5leHBvcnQgaW50ZXJmYWNlIFNpemUge1xuICAgIC8qKiBUaGUgd2lkdGggb2YgYSByZWN0YW5ndWxhciBhcmVhLiAqL1xuICAgIFdpZHRoOiBudW1iZXI7XG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgYSByZWN0YW5ndWxhciBhcmVhLiAqL1xuICAgIEhlaWdodDogbnVtYmVyO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFJlY3Qge1xuICAgIC8qKiBUaGUgWCBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uICovXG4gICAgWDogbnVtYmVyO1xuICAgIC8qKiBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uICovXG4gICAgWTogbnVtYmVyO1xuICAgIC8qKiBUaGUgd2lkdGggb2YgdGhlIHJlY3RhbmdsZS4gKi9cbiAgICBXaWR0aDogbnVtYmVyO1xuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIHRoZSByZWN0YW5nbGUuICovXG4gICAgSGVpZ2h0OiBudW1iZXI7XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgU2NyZWVuIHtcbiAgICAvKiogVW5pcXVlIGlkZW50aWZpZXIgZm9yIHRoZSBzY3JlZW4uICovXG4gICAgSUQ6IHN0cmluZztcbiAgICAvKiogSHVtYW4tcmVhZGFibGUgbmFtZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIE5hbWU6IHN0cmluZztcbiAgICAvKiogVGhlIHNjYWxlIGZhY3RvciBvZiB0aGUgc2NyZWVuIChEUEkvOTYpLiAxID0gc3RhbmRhcmQgRFBJLCAyID0gSGlEUEkgKFJldGluYSksIGV0Yy4gKi9cbiAgICBTY2FsZUZhY3RvcjogbnVtYmVyO1xuICAgIC8qKiBUaGUgWCBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uICovXG4gICAgWDogbnVtYmVyO1xuICAgIC8qKiBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uICovXG4gICAgWTogbnVtYmVyO1xuICAgIC8qKiBDb250YWlucyB0aGUgd2lkdGggYW5kIGhlaWdodCBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFNpemU6IFNpemU7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodC4gKi9cbiAgICBCb3VuZHM6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodCAoYmVmb3JlIHNjYWxpbmcpLiAqL1xuICAgIFBoeXNpY2FsQm91bmRzOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgYXJlYSBvZiB0aGUgc2NyZWVuIHRoYXQgaXMgYWN0dWFsbHkgdXNhYmxlIChleGNsdWRpbmcgdGFza2JhciBhbmQgb3RoZXIgc3lzdGVtIFVJKS4gKi9cbiAgICBXb3JrQXJlYTogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIFdvcmtBcmVhIG9mIHRoZSBzY3JlZW4gKGJlZm9yZSBzY2FsaW5nKS4gKi9cbiAgICBQaHlzaWNhbFdvcmtBcmVhOiBSZWN0O1xuICAgIC8qKiBUcnVlIGlmIHRoaXMgaXMgdGhlIHByaW1hcnkgbW9uaXRvciBzZWxlY3RlZCBieSB0aGUgdXNlciBpbiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cbiAgICBJc1ByaW1hcnk6IGJvb2xlYW47XG4gICAgLyoqIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFJvdGF0aW9uOiBudW1iZXI7XG59XG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuU2NyZWVucyk7XG5cbmNvbnN0IGdldEFsbCA9IDA7XG5jb25zdCBnZXRQcmltYXJ5ID0gMTtcbmNvbnN0IGdldEN1cnJlbnQgPSAyO1xuXG4vKipcbiAqIEdldHMgYWxsIHNjcmVlbnMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gYXJyYXkgb2YgU2NyZWVuIG9iamVjdHMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKTogUHJvbWlzZTxTY3JlZW5bXT4ge1xuICAgIHJldHVybiBjYWxsKGdldEFsbCk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIHByaW1hcnkgc2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0Q3VycmVudCk7XG59XG4iLCAiLyoqXG4gKiBUcmFuc3BvcnRDb2RlYyBkZWZpbmVzIHRoZSBpbnRlcmZhY2UgZm9yIGVuY29kaW5nL2RlY29kaW5nIHRyYW5zcG9ydCBkYXRhLlxuICogVGhpcyBhbGxvd3MgZGV2ZWxvcGVycyB0byBwbHVnIGluIGN1c3RvbSBzZXJpYWxpemF0aW9uIHN0cmF0ZWdpZXNcbiAqIChlLmcuLCBNZXNzYWdlUGFjaywgUHJvdG9idWYsIGN1c3RvbSBiaW5hcnkgZm9ybWF0cykgaW5zdGVhZCBvZiB0aGUgZGVmYXVsdCBiYXNlNjQvSlNPTi5cbiAqL1xuZXhwb3J0IGludGVyZmFjZSBUcmFuc3BvcnRDb2RlYyB7XG4gICAgLyoqXG4gICAgICogRGVjb2RlIHJlc3BvbnNlIGRhdGEgZnJvbSB0aGUgdHJhbnNwb3J0IGZvcm1hdCB0byBhIHN0cmluZ1xuICAgICAqIEBwYXJhbSBkYXRhIC0gVGhlIGVuY29kZWQgZGF0YSAoZS5nLiwgYmFzZTY0IHN0cmluZywgYmluYXJ5IGFycmF5LCBldGMuKVxuICAgICAqIEBwYXJhbSBjb250ZW50VHlwZSAtIFRoZSBjb250ZW50IHR5cGUgb2YgdGhlIHJlc3BvbnNlXG4gICAgICogQHJldHVybnMgRGVjb2RlZCBzdHJpbmcgZGF0YVxuICAgICAqL1xuICAgIGRlY29kZVJlc3BvbnNlKGRhdGE6IGFueSwgY29udGVudFR5cGU6IHN0cmluZyk6IHN0cmluZztcblxuICAgIC8qKlxuICAgICAqIERlY29kZSBlcnJvciBkYXRhIGZyb20gdGhlIHRyYW5zcG9ydCBmb3JtYXQgdG8gYSBzdHJpbmdcbiAgICAgKiBAcGFyYW0gZGF0YSAtIFRoZSBlbmNvZGVkIGVycm9yIGRhdGFcbiAgICAgKiBAcmV0dXJucyBEZWNvZGVkIGVycm9yIG1lc3NhZ2VcbiAgICAgKi9cbiAgICBkZWNvZGVFcnJvcihkYXRhOiBhbnkpOiBzdHJpbmc7XG5cbiAgICAvKipcbiAgICAgKiBFbmNvZGUgcmVxdWVzdCBhcmd1bWVudHMgZm9yIHRyYW5zcG9ydCAob3B0aW9uYWwpXG4gICAgICogQHBhcmFtIGFyZ3MgLSBUaGUgcmVxdWVzdCBhcmd1bWVudHMgb2JqZWN0XG4gICAgICogQHJldHVybnMgRW5jb2RlZCBkYXRhIHN1aXRhYmxlIGZvciB0cmFuc3BvcnRcbiAgICAgKi9cbiAgICBlbmNvZGVSZXF1ZXN0PyhhcmdzOiBhbnkpOiBhbnk7XG59XG5cbi8qKlxuICogQmFzZTY0SlNPTkNvZGVjIGlzIHRoZSBkZWZhdWx0IGNvZGVjIHRoYXQgaGFuZGxlcyBHbydzIEpTT04gbWFyc2hhbGluZyBiZWhhdmlvci5cbiAqIEdvIEpTT04gbWFyc2hhbHMgW11ieXRlIGFzIGJhc2U2NCBzdHJpbmdzLCBzbyB0aGlzIGNvZGVjIGRlY29kZXMgYmFzZTY0IGJhY2sgdG8gVVRGLTguXG4gKi9cbmV4cG9ydCBjbGFzcyBCYXNlNjRKU09OQ29kZWMgaW1wbGVtZW50cyBUcmFuc3BvcnRDb2RlYyB7XG4gICAgLyoqXG4gICAgICogRGVjb2RlIGJhc2U2NC1lbmNvZGVkIHJlc3BvbnNlIGRhdGEgdG8gVVRGLTggc3RyaW5nXG4gICAgICovXG4gICAgZGVjb2RlUmVzcG9uc2UoZGF0YTogYW55LCBjb250ZW50VHlwZTogc3RyaW5nKTogc3RyaW5nIHtcbiAgICAgICAgaWYgKCFkYXRhKSByZXR1cm4gJyc7XG5cbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIC8vIERlY29kZSBiYXNlNjQgc3RyaW5nIHRvIGJpbmFyeVxuICAgICAgICAgICAgY29uc3QgYmluYXJ5U3RyaW5nID0gYXRvYihkYXRhKTtcbiAgICAgICAgICAgIGNvbnN0IGJ5dGVzID0gbmV3IFVpbnQ4QXJyYXkoYmluYXJ5U3RyaW5nLmxlbmd0aCk7XG4gICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IGJpbmFyeVN0cmluZy5sZW5ndGg7IGkrKykge1xuICAgICAgICAgICAgICAgIGJ5dGVzW2ldID0gYmluYXJ5U3RyaW5nLmNoYXJDb2RlQXQoaSk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBDb252ZXJ0IGJ5dGVzIHRvIFVURi04IHN0cmluZ1xuICAgICAgICAgICAgcmV0dXJuIG5ldyBUZXh0RGVjb2RlcigpLmRlY29kZShieXRlcyk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgdGhyb3cgbmV3IEVycm9yKGBGYWlsZWQgdG8gZGVjb2RlIHJlc3BvbnNlOiAke2VyciBpbnN0YW5jZW9mIEVycm9yID8gZXJyLm1lc3NhZ2UgOiBTdHJpbmcoZXJyKX1gKTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERlY29kZSBiYXNlNjQtZW5jb2RlZCBlcnJvciBkYXRhIHRvIFVURi04IHN0cmluZ1xuICAgICAqL1xuICAgIGRlY29kZUVycm9yKGRhdGE6IGFueSk6IHN0cmluZyB7XG4gICAgICAgIGlmICghZGF0YSkgcmV0dXJuICdVbmtub3duIGVycm9yJztcblxuICAgICAgICB0cnkge1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMuZGVjb2RlUmVzcG9uc2UoZGF0YSwgJ3RleHQvcGxhaW4nKTtcbiAgICAgICAgfSBjYXRjaCB7XG4gICAgICAgICAgICByZXR1cm4gJ0ZhaWxlZCB0byBkZWNvZGUgZXJyb3IgbWVzc2FnZSc7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogUmF3U3RyaW5nQ29kZWMgaGFuZGxlcyByZXNwb25zZXMgdGhhdCBhcmUgYWxyZWFkeSBwbGFpbiBzdHJpbmdzIChubyBlbmNvZGluZykuXG4gKiBVc2VmdWwgZm9yIHRyYW5zcG9ydHMgdGhhdCBkb24ndCB1c2UgYmFzZTY0IGVuY29kaW5nLlxuICovXG5leHBvcnQgY2xhc3MgUmF3U3RyaW5nQ29kZWMgaW1wbGVtZW50cyBUcmFuc3BvcnRDb2RlYyB7XG4gICAgZGVjb2RlUmVzcG9uc2UoZGF0YTogYW55LCBjb250ZW50VHlwZTogc3RyaW5nKTogc3RyaW5nIHtcbiAgICAgICAgcmV0dXJuIGRhdGEgPyBTdHJpbmcoZGF0YSkgOiAnJztcbiAgICB9XG5cbiAgICBkZWNvZGVFcnJvcihkYXRhOiBhbnkpOiBzdHJpbmcge1xuICAgICAgICByZXR1cm4gZGF0YSA/IFN0cmluZyhkYXRhKSA6ICdVbmtub3duIGVycm9yJztcbiAgICB9XG59XG5cbi8qKlxuICogUmF3SlNPTkNvZGVjIGhhbmRsZXMgcmVzcG9uc2VzIHRoYXQgYXJlIGFscmVhZHkgSlNPTiBvYmplY3RzLlxuICogVXNlZnVsIGZvciB0cmFuc3BvcnRzIHRoYXQgc2VuZCBwYXJzZWQgSlNPTiBkaXJlY3RseS5cbiAqL1xuZXhwb3J0IGNsYXNzIFJhd0pTT05Db2RlYyBpbXBsZW1lbnRzIFRyYW5zcG9ydENvZGVjIHtcbiAgICBkZWNvZGVSZXNwb25zZShkYXRhOiBhbnksIGNvbnRlbnRUeXBlOiBzdHJpbmcpOiBzdHJpbmcge1xuICAgICAgICBpZiAoIWRhdGEpIHJldHVybiAnJztcbiAgICAgICAgcmV0dXJuIHR5cGVvZiBkYXRhID09PSAnc3RyaW5nJyA/IGRhdGEgOiBKU09OLnN0cmluZ2lmeShkYXRhKTtcbiAgICB9XG5cbiAgICBkZWNvZGVFcnJvcihkYXRhOiBhbnkpOiBzdHJpbmcge1xuICAgICAgICBpZiAoIWRhdGEpIHJldHVybiAnVW5rbm93biBlcnJvcic7XG4gICAgICAgIHJldHVybiB0eXBlb2YgZGF0YSA9PT0gJ3N0cmluZycgPyBkYXRhIDogSlNPTi5zdHJpbmdpZnkoZGF0YSk7XG4gICAgfVxufVxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBOzs7QUM2QkEsSUFBTSxjQUNGO0FBRUcsU0FBUyxPQUFPLE9BQWUsSUFBWTtBQUM5QyxNQUFJLEtBQUs7QUFFVCxNQUFJLElBQUksT0FBTztBQUNmLFNBQU8sS0FBSztBQUVSLFVBQU0sWUFBYSxLQUFLLE9BQU8sSUFBSSxLQUFNLENBQUM7QUFBQSxFQUM5QztBQUNBLFNBQU87QUFDWDs7O0FDN0JBLElBQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQU1yQyxJQUFNLGNBQWMsT0FBTyxPQUFPO0FBQUEsRUFDckMsTUFBTTtBQUFBLEVBQ04sV0FBVztBQUFBLEVBQ1gsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsWUFBWTtBQUNoQixDQUFDO0FBQ00sSUFBSSxXQUFXLE9BQU87QUF1QjdCLElBQUksa0JBQTJDO0FBc0J4QyxTQUFTLGFBQWEsV0FBMEM7QUFDbkUs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SUFDdkIscUJBQXFCO0FBQUEsSUFDckIscUJBQXFCO0FBQUEsSUFDckIsdUJBQXVCO0FBQUEsSUFDdkIsY0FBYztBQUFBLElBQ2QsZUFBZTtBQUFBLElBQ2Ysb0JBQW9CO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsMEJBQTBCO0FBQUEsSUFDMUIsZ0JBQWdCO0FBQUEsSUFDaEIsNEJBQTRCO0FBQUEsSUFDNUIsNEJBQTRCO0FBQUEsSUFDNUIseURBQXlEO0FBQUEsSUFDekQsc0NBQXNDO0FBQUEsSUFDdEMsb0JBQW9CO0FBQUEsSUFDcEIscUJBQXFCO0FBQUEsSUFDckIscUJBQXFCO0FBQUEsSUFDckIsc0JBQXNCO0FBQUEsSUFDdEIsZ0NBQWdDO0FBQUEsSUFDaEMsa0NBQWtDO0FBQUEsSUFDbEMsbUNBQW1DO0FBQUEsSUFDbkMsb0NBQW9DO0FBQUEsSUFDcEMsK0JBQStCO0FBQUEsSUFDL0IsNkJBQTZCO0FBQUEsSUFDN0IsdUJBQXVCO0FBQUEsSUFDdkIsaUNBQWlDO0FBQUEsSUFDakMsOEJBQThCO0FBQUEsSUFDOUIsNEJBQTRCO0FBQUEsSUFDNUIsc0NBQXNDO0FBQUEsSUFDdEMsNEJBQTRCO0FBQUEsSUFDNUIsc0JBQXNCO0FBQUEsSUFDdEIsa0NBQWtDO0FBQUEsSUFDbEMsc0JBQXNCO0FBQUEsSUFDdEIsd0JBQXdCO0FBQUEsSUFDeEIsd0JBQXdCO0FBQUEsSUFDeEIsbUJBQW1CO0FBQUEsSUFDbkIsMEJBQTBCO0FBQUEsSUFDMUIsOEJBQThCO0FBQUEsSUFDOUIseUJBQXlCO0FBQUEsSUFDekIsNkJBQTZCO0FBQUEsSUFDN0IsaUJBQWlCO0FBQUEsSUFDakIsZ0JBQWdCO0FBQUEsSUFDaEIsc0JBQXNCO0FBQUEsSUFDdEIsZUFBZTtBQUFBLElBQ2YseUJBQXlCO0FBQUEsSUFDekIsd0JBQXdCO0FBQUEsSUFDeEIsb0JBQW9CO0FBQUEsSUFDcEIscUJBQXFCO0FBQUEsSUFDckIsaUJBQWlCO0FBQUEsSUFDakIsaUJBQWlCO0FBQUEsSUFDakIsc0JBQXNCO0FBQUEsSUFDdEIsbUNBQW1DO0FBQUEsSUFDbkMscUNBQXFDO0FBQUEsSUFDckMsdUJBQXVCO0FBQUEsSUFDdkIsc0JBQXNCO0FBQUEsSUFDdEIsd0JBQXdCO0FBQUEsSUFDeEIsZUFBZTtBQUFBLElBQ2YsMkJBQTJCO0FBQUEsSUFDM0IsMEJBQTBCO0FBQUEsSUFDMUIsNkJBQTZCO0FBQUEsSUFDN0IsWUFBWTtBQUFBLElBQ1osZ0JBQWdCO0FBQUEsSUFDaEIsa0JBQWtCO0FBQUEsSUFDbEIsZ0JBQWdCO0FBQUEsSUFDaEIsa0JBQWtCO0FBQUEsSUFDbEIsbUJBQW1CO0FBQUEsSUFDbkIsWUFBWTtBQUFBLElBQ1oscUJBQXFCO0FBQUEsSUFDckIsc0JBQXNCO0FBQUEsSUFDdEIsc0JBQXNCO0FBQUEsSUFDdEIsOEJBQThCO0FBQUEsSUFDOUIsaUJBQWlCO0FBQUEsSUFDakIseUJBQXlCO0FBQUEsSUFDekIsMkJBQTJCO0FBQUEsSUFDM0IsK0JBQStCO0FBQUEsSUFDL0IsMEJBQTBCO0FBQUEsSUFDMUIsOEJBQThCO0FBQUEsSUFDOUIsaUJBQWlCO0FBQUEsSUFDakIsdUJBQXVCO0FBQUEsSUFDdkIsZ0JBQWdCO0FBQUEsSUFDaEIsMEJBQTBCO0FBQUEsSUFDMUIseUJBQXlCO0FBQUEsSUFDekIsc0JBQXNCO0FBQUEsSUFDdEIsa0JBQWtCO0FBQUEsSUFDbEIsbUJBQW1CO0FBQUEsSUFDbkIsa0JBQWtCO0FBQUEsSUFDbEIsdUJBQXVCO0FBQUEsSUFDdkIsb0NBQW9DO0FBQUEsSUFDcEMsc0NBQXNDO0FBQUEsSUFDdEMsd0JBQXdCO0FBQUEsSUFDeEIsdUJBQXVCO0FBQUEsSUFDdkIseUJBQXlCO0FBQUEsSUFDekIsNEJBQTRCO0FBQUEsSUFDNUIsNEJBQTRCO0FBQUEsSUFDNUIsY0FBYztBQUFBLElBQ2QsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsRUFDbEIsQ0FBQztBQUFBLEVBQ0QsT0FBTyxPQUFPLE9BQU87QUFBQSxJQUNwQixvQkFBb0I7QUFBQSxJQUNwQixvQkFBb0I7QUFBQSxJQUNwQixtQkFBbUI7QUFBQSxJQUNuQixlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxJQUNqQixlQUFlO0FBQUEsSUFDZixnQkFBZ0I7QUFBQSxJQUNoQixtQkFBbUI7QUFBQSxFQUNwQixDQUFDO0FBQUEsRUFDRCxRQUFRLE9BQU8sT0FBTztBQUFBLElBQ3JCLDJCQUEyQjtBQUFBLElBQzNCLG9CQUFvQjtBQUFBLElBQ3BCLDRCQUE0QjtBQUFBLElBQzVCLGNBQWM7QUFBQSxJQUNkLGVBQWU7QUFBQSxJQUNmLGVBQWU7QUFBQSxJQUNmLGlCQUFpQjtBQUFBLElBQ2pCLGtCQUFrQjtBQUFBLElBQ2xCLG9CQUFvQjtBQUFBLElBQ3BCLGFBQWE7QUFBQSxJQUNiLGtCQUFrQjtBQUFBLElBQ2xCLFlBQVk7QUFBQSxJQUNaLGlCQUFpQjtBQUFBLElBQ2pCLGdCQUFnQjtBQUFBLElBQ2hCLGdCQUFnQjtBQUFBLElBQ2hCLHVCQUF1QjtBQUFBLElBQ3ZCLGVBQWU7QUFBQSxJQUNmLG9CQUFvQjtBQUFBLElBQ3BCLFlBQVk7QUFBQSxJQUNaLG9CQUFvQjtBQUFBLElBQ3BCLGtCQUFrQjtBQUFBLElBQ2xCLGtCQUFrQjtBQUFBLElBQ2xCLFlBQVk7QUFBQSxJQUNaLGNBQWM7QUFBQSxJQUNkLGVBQWU7QUFBQSxJQUNmLGlCQUFpQjtBQUFBLElBQ2pCLDRCQUE0QjtBQUFBLEVBQzdCLENBQUM7QUFDRixDQUFDOzs7QUYzTkQsT0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBQ2xDLE9BQU8sT0FBTyxxQkFBcUJFO0FBRW5DLElBQU1DLFFBQU8saUJBQWlCLFlBQVksTUFBTTtBQUNoRCxJQUFNLGFBQWE7QUFZWixJQUFNLGFBQU4sTUFBaUI7QUFBQSxFQWlCcEIsWUFBWSxNQUFjLE9BQVksTUFBTTtBQUN4QyxTQUFLLE9BQU87QUFDWixTQUFLLE9BQU87QUFBQSxFQUNoQjtBQUNKO0FBRUEsU0FBU0Qsb0JBQW1CLE9BQVk7QUFDcEMsTUFBSSxZQUFZLGVBQWUsSUFBSSxNQUFNLElBQUk7QUFDN0MsTUFBSSxDQUFDLFdBQVc7QUFDWjtBQUFBLEVBQ0o7QUFFQSxNQUFJLGFBQWEsSUFBSSxXQUFXLE1BQU0sTUFBTSxNQUFNLElBQUk7QUFDdEQsTUFBSSxZQUFZLE9BQU87QUFDbkIsZUFBVyxTQUFTLE1BQU07QUFBQSxFQUM5QjtBQUVBLGNBQVksVUFBVSxPQUFPLGNBQVksQ0FBQyxTQUFTLFNBQVMsVUFBVSxDQUFDO0FBQ3ZFLE1BQUksVUFBVSxXQUFXLEdBQUc7QUFDeEIsbUJBQWUsT0FBTyxNQUFNLElBQUk7QUFBQSxFQUNwQyxPQUFPO0FBQ0gsbUJBQWUsSUFBSSxNQUFNLE1BQU0sU0FBUztBQUFBLEVBQzVDO0FBQ0o7QUFVTyxTQUFTLFdBQVcsV0FBbUIsVUFBb0IsY0FBc0I7QUFDcEYsTUFBSSxZQUFZLGVBQWUsSUFBSSxTQUFTLEtBQUssQ0FBQztBQUNsRCxRQUFNLGVBQWUsSUFBSSxTQUFTLFdBQVcsVUFBVSxZQUFZO0FBQ25FLFlBQVUsS0FBSyxZQUFZO0FBQzNCLGlCQUFlLElBQUksV0FBVyxTQUFTO0FBQ3ZDLFNBQU8sTUFBTSxZQUFZLFlBQVk7QUFDekM7QUFTTyxTQUFTLEdBQUcsV0FBbUIsVUFBZ0M7QUFDbEUsU0FBTyxXQUFXLFdBQVcsVUFBVSxFQUFFO0FBQzdDO0FBU08sU0FBUyxLQUFLLFdBQW1CLFVBQWdDO0FBQ3BFLFNBQU8sV0FBVyxXQUFXLFVBQVUsQ0FBQztBQUM1QztBQU9PLFNBQVMsT0FBTyxZQUF5QztBQUM1RCxhQUFXLFFBQVEsZUFBYSxlQUFlLE9BQU8sU0FBUyxDQUFDO0FBQ3BFO0FBS08sU0FBUyxTQUFlO0FBQzNCLGlCQUFlLE1BQU07QUFDekI7QUFTTyxTQUFTLEtBQUssTUFBYyxNQUEyQjtBQUMxRCxNQUFJO0FBQ0osTUFBSTtBQUVKLE1BQUksT0FBTyxTQUFTLFlBQVksU0FBUyxRQUFRLFVBQVUsUUFBUSxVQUFVLE1BQU07QUFFL0UsZ0JBQVksS0FBSyxNQUFNO0FBQ3ZCLGdCQUFZLEtBQUssTUFBTTtBQUFBLEVBQzNCLE9BQU87QUFFSCxnQkFBWTtBQUNaLGdCQUFZO0FBQUEsRUFDaEI7QUFFQSxTQUFPQyxNQUFLLFlBQVksRUFBRSxNQUFNLFdBQVcsTUFBTSxVQUFVLENBQUM7QUFDaEU7OztBR3JJTyxTQUFTLFNBQVMsU0FBYztBQUVuQyxVQUFRO0FBQUEsSUFDSixrQkFBa0IsVUFBVTtBQUFBLElBQzVCO0FBQUEsSUFDQTtBQUFBLEVBQ0o7QUFDSjtBQU1PLFNBQVMsa0JBQTJCO0FBQ3ZDLFNBQVEsSUFBSSxXQUFXLFdBQVcsRUFBRyxZQUFZO0FBQ3JEO0FBTU8sU0FBUyxvQkFBb0I7QUFDaEMsTUFBSSxDQUFDLGVBQWUsQ0FBQyxlQUFlLENBQUM7QUFDakMsV0FBTztBQUVYLE1BQUksU0FBUztBQUViLFFBQU0sU0FBUyxJQUFJLFlBQVk7QUFDL0IsUUFBTSxhQUFhLElBQUksZ0JBQWdCO0FBQ3ZDLFNBQU8saUJBQWlCLFFBQVEsTUFBTTtBQUFFLGFBQVM7QUFBQSxFQUFPLEdBQUcsRUFBRSxRQUFRLFdBQVcsT0FBTyxDQUFDO0FBQ3hGLGFBQVcsTUFBTTtBQUNqQixTQUFPLGNBQWMsSUFBSSxZQUFZLE1BQU0sQ0FBQztBQUU1QyxTQUFPO0FBQ1g7QUFLTyxTQUFTLFlBQVksT0FBMkI7QUF0RHZELE1BQUFDO0FBdURJLE1BQUksTUFBTSxrQkFBa0IsYUFBYTtBQUNyQyxXQUFPLE1BQU07QUFBQSxFQUNqQixXQUFXLEVBQUUsTUFBTSxrQkFBa0IsZ0JBQWdCLE1BQU0sa0JBQWtCLE1BQU07QUFDL0UsWUFBT0EsTUFBQSxNQUFNLE9BQU8sa0JBQWIsT0FBQUEsTUFBOEIsU0FBUztBQUFBLEVBQ2xELE9BQU87QUFDSCxXQUFPLFNBQVM7QUFBQSxFQUNwQjtBQUNKO0FBaUNBLElBQUksVUFBVTtBQUNkLFNBQVMsaUJBQWlCLG9CQUFvQixNQUFNO0FBQUUsWUFBVTtBQUFLLENBQUM7QUFFL0QsU0FBUyxVQUFVLFVBQXNCO0FBQzVDLE1BQUksV0FBVyxTQUFTLGVBQWUsWUFBWTtBQUMvQyxhQUFTO0FBQUEsRUFDYixPQUFPO0FBQ0gsYUFBUyxpQkFBaUIsb0JBQW9CLFFBQVE7QUFBQSxFQUMxRDtBQUNKOzs7QUMxRkEsSUFBTSxxQkFBcUI7QUFDM0IsSUFBTSx1QkFBdUI7QUFDN0IsSUFBSSx5QkFBeUM7QUFFN0MsSUFBTSxpQkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGNBQW9DO0FBQzFDLElBQU0sK0JBQW9DO0FBQzFDLElBQU0sOEJBQW9DO0FBQzFDLElBQU0sY0FBb0M7QUFDMUMsSUFBTSxvQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxrQkFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0sa0JBQW9DO0FBQzFDLElBQU0scUJBQW9DO0FBQzFDLElBQU0sb0JBQW9DO0FBQzFDLElBQU0sb0JBQW9DO0FBQzFDLElBQU0saUJBQW9DO0FBQzFDLElBQU0saUJBQW9DO0FBQzFDLElBQU0sYUFBb0M7QUFDMUMsSUFBTSxxQkFBb0M7QUFDMUMsSUFBTSx5QkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGtCQUFvQztBQUMxQyxJQUFNLGdCQUFvQztBQUMxQyxJQUFNLG9CQUFvQztBQUMxQyxJQUFNLHVCQUFvQztBQUMxQyxJQUFNLDRCQUFvQztBQUMxQyxJQUFNLHFCQUFvQztBQUMxQyxJQUFNLG1DQUFvQztBQUMxQyxJQUFNLG1CQUFvQztBQUMxQyxJQUFNLG1CQUFvQztBQUMxQyxJQUFNLDRCQUFvQztBQUMxQyxJQUFNLHFCQUFvQztBQUMxQyxJQUFNLGdCQUFvQztBQUMxQyxJQUFNLGlCQUFvQztBQUMxQyxJQUFNLGdCQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0sYUFBb0M7QUFDMUMsSUFBTSx5QkFBb0M7QUFDMUMsSUFBTSx1QkFBb0M7QUFDMUMsSUFBTSx3QkFBb0M7QUFDMUMsSUFBTSxxQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxjQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0sZUFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxrQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSx3QkFBb0M7QUFFMUMsU0FBUyxtQkFBbUIsU0FBeUM7QUFDakUsTUFBSSxDQUFDLFNBQVM7QUFDVixXQUFPO0FBQUEsRUFDWDtBQUVBLFNBQU8sUUFBUSxRQUFRLElBQUksMkJBQWtCLElBQUc7QUFDcEQ7QUF1QkEsSUFBTSxZQUFZLE9BQU8sUUFBUTtBQUlwQjtBQUZiLElBQU0sVUFBTixNQUFNLFFBQU87QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVVULFlBQVksT0FBZSxJQUFJO0FBQzNCLFNBQUssU0FBUyxJQUFJLGlCQUFpQixZQUFZLFFBQVEsSUFBSTtBQUczRCxlQUFXLFVBQVUsT0FBTyxvQkFBb0IsUUFBTyxTQUFTLEdBQUc7QUFDL0QsVUFDSSxXQUFXLGlCQUNSLE9BQVEsS0FBYSxNQUFNLE1BQU0sWUFDdEM7QUFDRSxRQUFDLEtBQWEsTUFBTSxJQUFLLEtBQWEsTUFBTSxFQUFFLEtBQUssSUFBSTtBQUFBLE1BQzNEO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLElBQUksTUFBc0I7QUFDdEIsV0FBTyxJQUFJLFFBQU8sSUFBSTtBQUFBLEVBQzFCO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsV0FBOEI7QUFDMUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxjQUFjO0FBQUEsRUFDekM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFNBQXdCO0FBQ3BCLFdBQU8sS0FBSyxTQUFTLEVBQUUsWUFBWTtBQUFBLEVBQ3ZDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxRQUF1QjtBQUNuQixXQUFPLEtBQUssU0FBUyxFQUFFLFdBQVc7QUFBQSxFQUN0QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EseUJBQXdDO0FBQ3BDLFdBQU8sS0FBSyxTQUFTLEVBQUUsNEJBQTRCO0FBQUEsRUFDdkQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLHdCQUF1QztBQUNuQyxXQUFPLEtBQUssU0FBUyxFQUFFLDJCQUEyQjtBQUFBLEVBQ3REO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxRQUF1QjtBQUNuQixXQUFPLEtBQUssU0FBUyxFQUFFLFdBQVc7QUFBQSxFQUN0QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsY0FBNkI7QUFDekIsV0FBTyxLQUFLLFNBQVMsRUFBRSxpQkFBaUI7QUFBQSxFQUM1QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsYUFBNEI7QUFDeEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxnQkFBZ0I7QUFBQSxFQUMzQztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFlBQTZCO0FBQ3pCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZUFBZTtBQUFBLEVBQzFDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsVUFBMkI7QUFDdkIsV0FBTyxLQUFLLFNBQVMsRUFBRSxhQUFhO0FBQUEsRUFDeEM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxTQUEwQjtBQUN0QixXQUFPLEtBQUssU0FBUyxFQUFFLFlBQVk7QUFBQSxFQUN2QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxZQUE4QjtBQUMxQixXQUFPLEtBQUssU0FBUyxFQUFFLGVBQWU7QUFBQSxFQUMxQztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGVBQWlDO0FBQzdCLFdBQU8sS0FBSyxTQUFTLEVBQUUsa0JBQWtCO0FBQUEsRUFDN0M7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxjQUFnQztBQUM1QixXQUFPLEtBQUssU0FBUyxFQUFFLGlCQUFpQjtBQUFBLEVBQzVDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsY0FBZ0M7QUFDNUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxpQkFBaUI7QUFBQSxFQUM1QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsV0FBMEI7QUFDdEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxjQUFjO0FBQUEsRUFDekM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFdBQTBCO0FBQ3RCLFdBQU8sS0FBSyxTQUFTLEVBQUUsY0FBYztBQUFBLEVBQ3pDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsT0FBd0I7QUFDcEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGVBQThCO0FBQzFCLFdBQU8sS0FBSyxTQUFTLEVBQUUsa0JBQWtCO0FBQUEsRUFDN0M7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxtQkFBc0M7QUFDbEMsV0FBTyxLQUFLLFNBQVMsRUFBRSxzQkFBc0I7QUFBQSxFQUNqRDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsU0FBd0I7QUFDcEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxZQUFZO0FBQUEsRUFDdkM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxZQUE4QjtBQUMxQixXQUFPLEtBQUssU0FBUyxFQUFFLGVBQWU7QUFBQSxFQUMxQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsVUFBeUI7QUFDckIsV0FBTyxLQUFLLFNBQVMsRUFBRSxhQUFhO0FBQUEsRUFDeEM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLFlBQVksR0FBVyxHQUEwQjtBQUM3QyxXQUFPLEtBQUssU0FBUyxFQUFFLG1CQUFtQixFQUFFLEdBQUcsRUFBRSxDQUFDO0FBQUEsRUFDdEQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxlQUFlLGFBQXFDO0FBQ2hELFdBQU8sS0FBSyxTQUFTLEVBQUUsc0JBQXNCLEVBQUUsWUFBWSxDQUFDO0FBQUEsRUFDaEU7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFVQSxvQkFBb0IsR0FBVyxHQUFXLEdBQVcsR0FBMEI7QUFDM0UsV0FBTyxLQUFLLFNBQVMsRUFBRSwyQkFBMkIsRUFBRSxHQUFHLEdBQUcsR0FBRyxFQUFFLENBQUM7QUFBQSxFQUNwRTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGFBQWEsV0FBbUM7QUFDNUMsV0FBTyxLQUFLLFNBQVMsRUFBRSxvQkFBb0IsRUFBRSxVQUFVLENBQUM7QUFBQSxFQUM1RDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLDJCQUEyQixTQUFpQztBQUN4RCxXQUFPLEtBQUssU0FBUyxFQUFFLGtDQUFrQyxFQUFFLFFBQVEsQ0FBQztBQUFBLEVBQ3hFO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxXQUFXLE9BQWUsUUFBK0I7QUFDckQsV0FBTyxLQUFLLFNBQVMsRUFBRSxrQkFBa0IsRUFBRSxPQUFPLE9BQU8sQ0FBQztBQUFBLEVBQzlEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxXQUFXLE9BQWUsUUFBK0I7QUFDckQsV0FBTyxLQUFLLFNBQVMsRUFBRSxrQkFBa0IsRUFBRSxPQUFPLE9BQU8sQ0FBQztBQUFBLEVBQzlEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxvQkFBb0IsR0FBVyxHQUEwQjtBQUNyRCxXQUFPLEtBQUssU0FBUyxFQUFFLDJCQUEyQixFQUFFLEdBQUcsRUFBRSxDQUFDO0FBQUEsRUFDOUQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxhQUFhQyxZQUFtQztBQUM1QyxXQUFPLEtBQUssU0FBUyxFQUFFLG9CQUFvQixFQUFFLFdBQUFBLFdBQVUsQ0FBQztBQUFBLEVBQzVEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxRQUFRLE9BQWUsUUFBK0I7QUFDbEQsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlLEVBQUUsT0FBTyxPQUFPLENBQUM7QUFBQSxFQUMzRDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFNBQVMsT0FBOEI7QUFDbkMsV0FBTyxLQUFLLFNBQVMsRUFBRSxnQkFBZ0IsRUFBRSxNQUFNLENBQUM7QUFBQSxFQUNwRDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFFBQVEsTUFBNkI7QUFDakMsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlLEVBQUUsS0FBSyxDQUFDO0FBQUEsRUFDbEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLE9BQXNCO0FBQ2xCLFdBQU8sS0FBSyxTQUFTLEVBQUUsVUFBVTtBQUFBLEVBQ3JDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLG1CQUFrQztBQUM5QixXQUFPLEtBQUssU0FBUyxFQUFFLHNCQUFzQjtBQUFBLEVBQ2pEO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxpQkFBZ0M7QUFDNUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxvQkFBb0I7QUFBQSxFQUMvQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0Esa0JBQWlDO0FBQzdCLFdBQU8sS0FBSyxTQUFTLEVBQUUscUJBQXFCO0FBQUEsRUFDaEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGVBQThCO0FBQzFCLFdBQU8sS0FBSyxTQUFTLEVBQUUsa0JBQWtCO0FBQUEsRUFDN0M7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGFBQTRCO0FBQ3hCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZ0JBQWdCO0FBQUEsRUFDM0M7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGFBQTRCO0FBQ3hCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZ0JBQWdCO0FBQUEsRUFDM0M7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxRQUF5QjtBQUNyQixXQUFPLEtBQUssU0FBUyxFQUFFLFdBQVc7QUFBQSxFQUN0QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFNBQXdCO0FBQ3BCLFdBQU8sS0FBSyxTQUFTLEVBQUUsWUFBWTtBQUFBLEVBQ3ZDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxVQUF5QjtBQUNyQixXQUFPLEtBQUssU0FBUyxFQUFFLGFBQWE7QUFBQSxFQUN4QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsWUFBMkI7QUFDdkIsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlO0FBQUEsRUFDMUM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFVQSx1QkFBdUIsV0FBcUIsR0FBVyxHQUFpQjtBQUNwRSxVQUFNLFVBQVUsU0FBUyxpQkFBaUIsR0FBRyxDQUFDO0FBRzlDLFVBQU0saUJBQWlCLG1CQUFtQixPQUFPO0FBRWpELFFBQUksQ0FBQyxnQkFBZ0I7QUFDakIsY0FBUSxJQUFJLHFEQUFxRCxVQUFDLEtBQUksVUFBQyw0REFBMkQsT0FBTztBQUV6STtBQUFBLElBQ0o7QUFFQSxZQUFRLElBQUksMkRBQTJELFVBQUMsTUFBSyxVQUFDLE9BQU0sU0FBUyx1QkFBdUIsY0FBYztBQUNsSSxVQUFNLGlCQUFpQjtBQUFBLE1BQ25CLElBQUksZUFBZTtBQUFBLE1BQ25CLFdBQVcsTUFBTSxLQUFLLGVBQWUsU0FBUztBQUFBLE1BQzlDLFlBQVksQ0FBQztBQUFBLElBQ2pCO0FBQ0EsYUFBUyxJQUFJLEdBQUcsSUFBSSxlQUFlLFdBQVcsUUFBUSxLQUFLO0FBQ3ZELFlBQU0sT0FBTyxlQUFlLFdBQVcsQ0FBQztBQUN4QyxxQkFBZSxXQUFXLEtBQUssSUFBSSxJQUFJLEtBQUs7QUFBQSxJQUNoRDtBQUVBLFVBQU0sVUFBVTtBQUFBLE1BQ1o7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxJQUNKO0FBRUEsU0FBSyxTQUFTLEVBQUUsdUJBQXVCLE9BQU87QUFBQSxFQUNsRDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsYUFBNEI7QUFDeEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxnQkFBZ0I7QUFBQSxFQUMzQztBQUNKO0FBbGVBLElBQU0sU0FBTjtBQXVlQSxJQUFNLGFBQWEsSUFBSSxPQUFPLEVBQUU7QUFHaEMsU0FBUywrQkFBK0I7QUFDcEMsUUFBTSxhQUFhLFNBQVM7QUFDNUIsTUFBSSxtQkFBbUI7QUFFdkIsYUFBVyxpQkFBaUIsYUFBYSxDQUFDLFVBQVU7QUFDaEQsVUFBTSxlQUFlO0FBQ3JCLFFBQUksTUFBTSxnQkFBZ0IsTUFBTSxhQUFhLE1BQU0sU0FBUyxPQUFPLEdBQUc7QUFDbEU7QUFDQSxZQUFNLGdCQUFnQixTQUFTLGlCQUFpQixNQUFNLFNBQVMsTUFBTSxPQUFPO0FBQzVFLFlBQU0sV0FBVyxtQkFBbUIsYUFBYTtBQUdqRCxVQUFJLDBCQUEwQiwyQkFBMkIsVUFBVTtBQUMvRCwrQkFBdUIsVUFBVSxPQUFPLG9CQUFvQjtBQUFBLE1BQ2hFO0FBRUEsVUFBSSxVQUFVO0FBQ1YsaUJBQVMsVUFBVSxJQUFJLG9CQUFvQjtBQUMzQyxjQUFNLGFBQWEsYUFBYTtBQUNoQyxpQ0FBeUI7QUFBQSxNQUM3QixPQUFPO0FBQ0gsY0FBTSxhQUFhLGFBQWE7QUFDaEMsaUNBQXlCO0FBQUEsTUFDN0I7QUFBQSxJQUNKO0FBQUEsRUFDSixHQUFHLEtBQUs7QUFFUixhQUFXLGlCQUFpQixZQUFZLENBQUMsVUFBVTtBQUMvQyxVQUFNLGVBQWU7QUFDckIsUUFBSSxNQUFNLGdCQUFnQixNQUFNLGFBQWEsTUFBTSxTQUFTLE9BQU8sR0FBRztBQUdsRSxVQUFJLHdCQUF3QjtBQUV4QixZQUFHLENBQUMsdUJBQXVCLFVBQVUsU0FBUyxvQkFBb0IsR0FBRztBQUNqRSxpQ0FBdUIsVUFBVSxJQUFJLG9CQUFvQjtBQUFBLFFBQzdEO0FBQ0EsY0FBTSxhQUFhLGFBQWE7QUFBQSxNQUNwQyxPQUFPO0FBQ0gsY0FBTSxhQUFhLGFBQWE7QUFBQSxNQUNwQztBQUFBLElBQ0o7QUFBQSxFQUNKLEdBQUcsS0FBSztBQUVSLGFBQVcsaUJBQWlCLGFBQWEsQ0FBQyxVQUFVO0FBQ2hELFVBQU0sZUFBZTtBQUNyQixRQUFJLE1BQU0sZ0JBQWdCLE1BQU0sYUFBYSxNQUFNLFNBQVMsT0FBTyxHQUFHO0FBQ2xFO0FBRUEsVUFBSSxxQkFBcUIsS0FBSyxNQUFNLGtCQUFrQixRQUFTLDBCQUEwQixDQUFDLHVCQUF1QixTQUFTLE1BQU0sYUFBcUIsR0FBSTtBQUNySixZQUFJLHdCQUF3QjtBQUN4QixpQ0FBdUIsVUFBVSxPQUFPLG9CQUFvQjtBQUM1RCxtQ0FBeUI7QUFBQSxRQUM3QjtBQUNBLDJCQUFtQjtBQUFBLE1BQ3ZCO0FBQUEsSUFDSjtBQUFBLEVBQ0osR0FBRyxLQUFLO0FBRVIsYUFBVyxpQkFBaUIsUUFBUSxDQUFDLFVBQVU7QUFDM0MsVUFBTSxlQUFlO0FBQ3JCLHVCQUFtQjtBQUNuQixRQUFJLHdCQUF3QjtBQUN4Qiw2QkFBdUIsVUFBVSxPQUFPLG9CQUFvQjtBQUM1RCwrQkFBeUI7QUFBQSxJQUM3QjtBQUFBLEVBR0osR0FBRyxLQUFLO0FBQ1o7QUFHQSxJQUFJLE9BQU8sV0FBVyxlQUFlLE9BQU8sYUFBYSxhQUFhO0FBQ2xFLCtCQUE2QjtBQUNqQztBQUVBLElBQU8saUJBQVE7OztBVHJvQmYsU0FBUyxVQUFVLFdBQW1CLE9BQVksTUFBWTtBQUMxRCxPQUFLLFdBQVcsSUFBSTtBQUN4QjtBQVFBLFNBQVMsaUJBQWlCLFlBQW9CLFlBQW9CO0FBQzlELFFBQU0sZUFBZSxlQUFPLElBQUksVUFBVTtBQUMxQyxRQUFNLFNBQVUsYUFBcUIsVUFBVTtBQUUvQyxNQUFJLE9BQU8sV0FBVyxZQUFZO0FBQzlCLFlBQVEsTUFBTSxrQkFBa0IsbUJBQVUsY0FBYTtBQUN2RDtBQUFBLEVBQ0o7QUFFQSxNQUFJO0FBQ0EsV0FBTyxLQUFLLFlBQVk7QUFBQSxFQUM1QixTQUFTLEdBQUc7QUFDUixZQUFRLE1BQU0sZ0NBQWdDLG1CQUFVLFFBQU8sQ0FBQztBQUFBLEVBQ3BFO0FBQ0o7QUFLQSxTQUFTLGVBQWUsSUFBaUI7QUFDckMsUUFBTSxVQUFVLEdBQUc7QUFFbkIsV0FBUyxVQUFVLFNBQVMsT0FBTztBQUMvQixRQUFJLFdBQVc7QUFDWDtBQUVKLFVBQU0sWUFBWSxRQUFRLGFBQWEsV0FBVyxLQUFLLFFBQVEsYUFBYSxnQkFBZ0I7QUFDNUYsVUFBTSxlQUFlLFFBQVEsYUFBYSxtQkFBbUIsS0FBSyxRQUFRLGFBQWEsd0JBQXdCLEtBQUs7QUFDcEgsVUFBTSxlQUFlLFFBQVEsYUFBYSxZQUFZLEtBQUssUUFBUSxhQUFhLGlCQUFpQjtBQUNqRyxVQUFNLE1BQU0sUUFBUSxhQUFhLGFBQWEsS0FBSyxRQUFRLGFBQWEsa0JBQWtCO0FBRTFGLFFBQUksY0FBYztBQUNkLGdCQUFVLFNBQVM7QUFDdkIsUUFBSSxpQkFBaUI7QUFDakIsdUJBQWlCLGNBQWMsWUFBWTtBQUMvQyxRQUFJLFFBQVE7QUFDUixXQUFLLFFBQVEsR0FBRztBQUFBLEVBQ3hCO0FBRUEsUUFBTSxVQUFVLFFBQVEsYUFBYSxhQUFhLEtBQUssUUFBUSxhQUFhLGtCQUFrQjtBQUU5RixNQUFJLFNBQVM7QUFDVCxhQUFTO0FBQUEsTUFDTCxPQUFPO0FBQUEsTUFDUCxTQUFTO0FBQUEsTUFDVCxVQUFVO0FBQUEsTUFDVixTQUFTO0FBQUEsUUFDTCxFQUFFLE9BQU8sTUFBTTtBQUFBLFFBQ2YsRUFBRSxPQUFPLE1BQU0sV0FBVyxLQUFLO0FBQUEsTUFDbkM7QUFBQSxJQUNKLENBQUMsRUFBRSxLQUFLLFNBQVM7QUFBQSxFQUNyQixPQUFPO0FBQ0gsY0FBVTtBQUFBLEVBQ2Q7QUFDSjtBQUdBLElBQU0sZ0JBQWdCLE9BQU8sWUFBWTtBQUN6QyxJQUFNLGdCQUFnQixPQUFPLFlBQVk7QUFDekMsSUFBTSxrQkFBa0IsT0FBTyxjQUFjO0FBUXhDO0FBRkwsSUFBTSwwQkFBTixNQUE4QjtBQUFBLEVBSTFCLGNBQWM7QUFDVixTQUFLLGFBQWEsSUFBSSxJQUFJLGdCQUFnQjtBQUFBLEVBQzlDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVNBLElBQUksU0FBa0IsVUFBNkM7QUFDL0QsV0FBTyxFQUFFLFFBQVEsS0FBSyxhQUFhLEVBQUUsT0FBTztBQUFBLEVBQ2hEO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxRQUFjO0FBQ1YsU0FBSyxhQUFhLEVBQUUsTUFBTTtBQUMxQixTQUFLLGFBQWEsSUFBSSxJQUFJLGdCQUFnQjtBQUFBLEVBQzlDO0FBQ0o7QUFTSyxlQUVBO0FBSkwsSUFBTSxrQkFBTixNQUFzQjtBQUFBLEVBTWxCLGNBQWM7QUFDVixTQUFLLGFBQWEsSUFBSSxvQkFBSSxRQUFRO0FBQ2xDLFNBQUssZUFBZSxJQUFJO0FBQUEsRUFDNUI7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLElBQUksU0FBa0IsVUFBNkM7QUFDL0QsUUFBSSxDQUFDLEtBQUssYUFBYSxFQUFFLElBQUksT0FBTyxHQUFHO0FBQUUsV0FBSyxlQUFlO0FBQUEsSUFBSztBQUNsRSxTQUFLLGFBQWEsRUFBRSxJQUFJLFNBQVMsUUFBUTtBQUN6QyxXQUFPLENBQUM7QUFBQSxFQUNaO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxRQUFjO0FBQ1YsUUFBSSxLQUFLLGVBQWUsS0FBSztBQUN6QjtBQUVKLGVBQVcsV0FBVyxTQUFTLEtBQUssaUJBQWlCLEdBQUcsR0FBRztBQUN2RCxVQUFJLEtBQUssZUFBZSxLQUFLO0FBQ3pCO0FBRUosWUFBTSxXQUFXLEtBQUssYUFBYSxFQUFFLElBQUksT0FBTztBQUNoRCxVQUFJLFlBQVksTUFBTTtBQUFFLGFBQUssZUFBZTtBQUFBLE1BQUs7QUFFakQsaUJBQVcsV0FBVyxZQUFZLENBQUM7QUFDL0IsZ0JBQVEsb0JBQW9CLFNBQVMsY0FBYztBQUFBLElBQzNEO0FBRUEsU0FBSyxhQUFhLElBQUksb0JBQUksUUFBUTtBQUNsQyxTQUFLLGVBQWUsSUFBSTtBQUFBLEVBQzVCO0FBQ0o7QUFFQSxJQUFNLGtCQUFrQixrQkFBa0IsSUFBSSxJQUFJLHdCQUF3QixJQUFJLElBQUksZ0JBQWdCO0FBS2xHLFNBQVMsZ0JBQWdCLFNBQXdCO0FBQzdDLFFBQU0sZ0JBQWdCO0FBQ3RCLFFBQU0sY0FBZSxRQUFRLGFBQWEsYUFBYSxLQUFLLFFBQVEsYUFBYSxrQkFBa0IsS0FBSztBQUN4RyxRQUFNLFdBQXFCLENBQUM7QUFFNUIsTUFBSTtBQUNKLFVBQVEsUUFBUSxjQUFjLEtBQUssV0FBVyxPQUFPO0FBQ2pELGFBQVMsS0FBSyxNQUFNLENBQUMsQ0FBQztBQUUxQixRQUFNLFVBQVUsZ0JBQWdCLElBQUksU0FBUyxRQUFRO0FBQ3JELGFBQVcsV0FBVztBQUNsQixZQUFRLGlCQUFpQixTQUFTLGdCQUFnQixPQUFPO0FBQ2pFO0FBS08sU0FBUyxTQUFlO0FBQzNCLFlBQVUsTUFBTTtBQUNwQjtBQUtPLFNBQVMsU0FBZTtBQUMzQixrQkFBZ0IsTUFBTTtBQUN0QixXQUFTLEtBQUssaUJBQWlCLG1HQUFtRyxFQUFFLFFBQVEsZUFBZTtBQUMvSjs7O0FVaE1BLE9BQU8sUUFBUTtBQUNmLE9BQVU7QUFFVixJQUFJLE1BQU87QUFDUCxXQUFTLHNCQUFzQjtBQUNuQzs7O0FDckJBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQVlBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksTUFBTTtBQUVoRCxJQUFNLG1CQUFtQjtBQUN6QixJQUFNLG9CQUFvQjtBQUMxQixJQUFNLHFDQUFxQztBQUUzQyxJQUFNLFdBQVcsV0FBWTtBQWxCN0IsTUFBQUMsS0FBQTtBQW1CSSxNQUFJO0FBQ0EsU0FBSyxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsWUFBdkIsbUJBQWdDLGFBQWE7QUFDOUMsYUFBUSxPQUFlLE9BQU8sUUFBUSxZQUFZLEtBQU0sT0FBZSxPQUFPLE9BQU87QUFBQSxJQUN6RixZQUFZLHdCQUFlLFdBQWYsbUJBQXVCLG9CQUF2QixtQkFBeUMsZ0JBQXpDLG1CQUFzRCxhQUFhO0FBQzNFLGFBQVEsT0FBZSxPQUFPLGdCQUFnQixVQUFVLEVBQUUsWUFBWSxLQUFNLE9BQWUsT0FBTyxnQkFBZ0IsVUFBVSxDQUFDO0FBQUEsSUFDakk7QUFBQSxFQUNKLFNBQVEsR0FBRztBQUFBLEVBQUM7QUFFWixVQUFRO0FBQUEsSUFBSztBQUFBLElBQ1Q7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLEVBQXdEO0FBQzVELFNBQU87QUFDWCxHQUFHO0FBRUksU0FBUyxPQUFPLEtBQWdCO0FBQ25DLHFDQUFVO0FBQ2Q7QUFPTyxTQUFTLGFBQStCO0FBQzNDLFNBQU9ELE1BQUssZ0JBQWdCO0FBQ2hDO0FBT0EsZUFBc0IsZUFBNkM7QUFDL0QsTUFBSSxXQUFXLE1BQU0sTUFBTSxxQkFBcUI7QUFDaEQsTUFBSSxTQUFTLElBQUk7QUFDYixXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCLE9BQU87QUFDSCxVQUFNLElBQUksTUFBTSxtQ0FBbUMsU0FBUyxVQUFVO0FBQUEsRUFDMUU7QUFDSjtBQStCTyxTQUFTLGNBQXdDO0FBQ3BELFNBQU9BLE1BQUssaUJBQWlCO0FBQ2pDO0FBT08sU0FBUyxZQUFxQjtBQUNqQyxTQUFPLE9BQU8sT0FBTyxZQUFZLE9BQU87QUFDNUM7QUFPTyxTQUFTLFVBQW1CO0FBQy9CLFNBQU8sT0FBTyxPQUFPLFlBQVksT0FBTztBQUM1QztBQU9PLFNBQVMsUUFBaUI7QUFDN0IsU0FBTyxPQUFPLE9BQU8sWUFBWSxPQUFPO0FBQzVDO0FBT08sU0FBUyxVQUFtQjtBQUMvQixTQUFPLE9BQU8sT0FBTyxZQUFZLFNBQVM7QUFDOUM7QUFPTyxTQUFTLFFBQWlCO0FBQzdCLFNBQU8sT0FBTyxPQUFPLFlBQVksU0FBUztBQUM5QztBQU9PLFNBQVMsVUFBbUI7QUFDL0IsU0FBTyxPQUFPLE9BQU8sWUFBWSxTQUFTO0FBQzlDO0FBT08sU0FBUyxVQUFtQjtBQUMvQixTQUFPLFFBQVEsT0FBTyxPQUFPLFlBQVksS0FBSztBQUNsRDtBQVVPLFNBQVMsdUJBQXVCLFdBQXFCLEdBQVcsR0FBaUI7QUFDcEYsUUFBTSxVQUFVLFNBQVMsaUJBQWlCLEdBQUcsQ0FBQztBQUM5QyxRQUFNLFlBQVksVUFBVSxRQUFRLEtBQUs7QUFDekMsUUFBTSxZQUFZLFVBQVUsTUFBTSxLQUFLLFFBQVEsU0FBUyxJQUFJLENBQUM7QUFFN0QsUUFBTSxVQUFVO0FBQUEsSUFDWjtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxFQUNKO0FBRUEsRUFBQUEsTUFBSyxvQ0FBb0MsT0FBTyxFQUMzQyxLQUFLLE1BQU07QUFFUixZQUFRLElBQUksOENBQThDO0FBQUEsRUFDOUQsQ0FBQyxFQUNBLE1BQU0sU0FBTztBQUVWLFlBQVEsTUFBTSwyQ0FBMkMsR0FBRztBQUFBLEVBQ2hFLENBQUM7QUFDVDs7O0FDNUtBLE9BQU8saUJBQWlCLGVBQWUsa0JBQWtCO0FBRXpELElBQU1FLFFBQU8saUJBQWlCLFlBQVksV0FBVztBQUVyRCxJQUFNLGtCQUFrQjtBQUV4QixTQUFTLGdCQUFnQixJQUFZLEdBQVcsR0FBVyxNQUFpQjtBQUN4RSxPQUFLQSxNQUFLLGlCQUFpQixFQUFDLElBQUksR0FBRyxHQUFHLEtBQUksQ0FBQztBQUMvQztBQUVBLFNBQVMsbUJBQW1CLE9BQW1CO0FBQzNDLFFBQU0sU0FBUyxZQUFZLEtBQUs7QUFHaEMsUUFBTSxvQkFBb0IsT0FBTyxpQkFBaUIsTUFBTSxFQUFFLGlCQUFpQixzQkFBc0IsRUFBRSxLQUFLO0FBRXhHLE1BQUksbUJBQW1CO0FBQ25CLFVBQU0sZUFBZTtBQUNyQixVQUFNLE9BQU8sT0FBTyxpQkFBaUIsTUFBTSxFQUFFLGlCQUFpQiwyQkFBMkI7QUFDekYsb0JBQWdCLG1CQUFtQixNQUFNLFNBQVMsTUFBTSxTQUFTLElBQUk7QUFBQSxFQUN6RSxPQUFPO0FBQ0gsOEJBQTBCLE9BQU8sTUFBTTtBQUFBLEVBQzNDO0FBQ0o7QUFVQSxTQUFTLDBCQUEwQixPQUFtQixRQUFxQjtBQUV2RSxNQUFJLFFBQVEsR0FBRztBQUNYO0FBQUEsRUFDSjtBQUdBLFVBQVEsT0FBTyxpQkFBaUIsTUFBTSxFQUFFLGlCQUFpQix1QkFBdUIsRUFBRSxLQUFLLEdBQUc7QUFBQSxJQUN0RixLQUFLO0FBQ0Q7QUFBQSxJQUNKLEtBQUs7QUFDRCxZQUFNLGVBQWU7QUFDckI7QUFBQSxFQUNSO0FBR0EsTUFBSSxPQUFPLG1CQUFtQjtBQUMxQjtBQUFBLEVBQ0o7QUFHQSxRQUFNLFlBQVksT0FBTyxhQUFhO0FBQ3RDLFFBQU0sZUFBZSxhQUFhLFVBQVUsU0FBUyxFQUFFLFNBQVM7QUFDaEUsTUFBSSxjQUFjO0FBQ2QsYUFBUyxJQUFJLEdBQUcsSUFBSSxVQUFVLFlBQVksS0FBSztBQUMzQyxZQUFNLFFBQVEsVUFBVSxXQUFXLENBQUM7QUFDcEMsWUFBTSxRQUFRLE1BQU0sZUFBZTtBQUNuQyxlQUFTLElBQUksR0FBRyxJQUFJLE1BQU0sUUFBUSxLQUFLO0FBQ25DLGNBQU0sT0FBTyxNQUFNLENBQUM7QUFDcEIsWUFBSSxTQUFTLGlCQUFpQixLQUFLLE1BQU0sS0FBSyxHQUFHLE1BQU0sUUFBUTtBQUMzRDtBQUFBLFFBQ0o7QUFBQSxNQUNKO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFHQSxNQUFJLGtCQUFrQixvQkFBb0Isa0JBQWtCLHFCQUFxQjtBQUM3RSxRQUFJLGdCQUFpQixDQUFDLE9BQU8sWUFBWSxDQUFDLE9BQU8sVUFBVztBQUN4RDtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBR0EsUUFBTSxlQUFlO0FBQ3pCOzs7QUM3RkE7QUFBQTtBQUFBO0FBQUE7QUFnQk8sU0FBUyxRQUFRLEtBQWtCO0FBQ3RDLE1BQUk7QUFDQSxXQUFPLE9BQU8sT0FBTyxNQUFNLEdBQUc7QUFBQSxFQUNsQyxTQUFTLEdBQUc7QUFDUixVQUFNLElBQUksTUFBTSw4QkFBOEIsTUFBTSxRQUFRLEdBQUcsRUFBRSxPQUFPLEVBQUUsQ0FBQztBQUFBLEVBQy9FO0FBQ0o7OztBQ1BBLElBQUksVUFBVTtBQUNkLElBQUksV0FBVztBQUVmLElBQUksWUFBWTtBQUNoQixJQUFJLFlBQVk7QUFDaEIsSUFBSSxXQUFXO0FBQ2YsSUFBSSxhQUFxQjtBQUN6QixJQUFJLGdCQUFnQjtBQUVwQixJQUFJLFVBQVU7QUFDZCxJQUFNLGlCQUFpQixnQkFBZ0I7QUFFdkMsT0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBQ2xDLE9BQU8sT0FBTyxlQUFlLENBQUMsVUFBeUI7QUFDbkQsY0FBWTtBQUNaLE1BQUksQ0FBQyxXQUFXO0FBRVosZ0JBQVksV0FBVztBQUN2QixjQUFVO0FBQUEsRUFDZDtBQUNKO0FBRUEsT0FBTyxpQkFBaUIsYUFBYSxRQUFRLEVBQUUsU0FBUyxLQUFLLENBQUM7QUFDOUQsT0FBTyxpQkFBaUIsYUFBYSxRQUFRLEVBQUUsU0FBUyxLQUFLLENBQUM7QUFDOUQsT0FBTyxpQkFBaUIsV0FBVyxRQUFRLEVBQUUsU0FBUyxLQUFLLENBQUM7QUFDNUQsV0FBVyxNQUFNLENBQUMsU0FBUyxlQUFlLFVBQVUsR0FBRztBQUNuRCxTQUFPLGlCQUFpQixJQUFJLGVBQWUsRUFBRSxTQUFTLEtBQUssQ0FBQztBQUNoRTtBQUVBLFNBQVMsY0FBYyxPQUFjO0FBRWpDLE1BQUksWUFBWSxVQUFVO0FBQ3RCLFVBQU0seUJBQXlCO0FBQy9CLFVBQU0sZ0JBQWdCO0FBQ3RCLFVBQU0sZUFBZTtBQUFBLEVBQ3pCO0FBQ0o7QUFHQSxJQUFNLFlBQVk7QUFDbEIsSUFBTSxVQUFZO0FBQ2xCLElBQU0sWUFBWTtBQUVsQixTQUFTLE9BQU8sT0FBbUI7QUFJL0IsTUFBSSxXQUFtQixlQUFlLE1BQU07QUFDNUMsVUFBUSxNQUFNLE1BQU07QUFBQSxJQUNoQixLQUFLO0FBQ0Qsa0JBQVk7QUFDWixVQUFJLENBQUMsZ0JBQWdCO0FBQUUsdUJBQWUsVUFBVyxLQUFLLE1BQU07QUFBQSxNQUFTO0FBQ3JFO0FBQUEsSUFDSixLQUFLO0FBQ0Qsa0JBQVk7QUFDWixVQUFJLENBQUMsZ0JBQWdCO0FBQUUsdUJBQWUsVUFBVSxFQUFFLEtBQUssTUFBTTtBQUFBLE1BQVM7QUFDdEU7QUFBQSxJQUNKO0FBQ0ksa0JBQVk7QUFDWixVQUFJLENBQUMsZ0JBQWdCO0FBQUUsdUJBQWU7QUFBQSxNQUFTO0FBQy9DO0FBQUEsRUFDUjtBQUVBLE1BQUksV0FBVyxVQUFVLENBQUM7QUFDMUIsTUFBSSxVQUFVLGVBQWUsQ0FBQztBQUU5QixZQUFVO0FBR1YsTUFBSSxjQUFjLGFBQWEsRUFBRSxVQUFVLE1BQU0sU0FBUztBQUN0RCxnQkFBYSxLQUFLLE1BQU07QUFDeEIsZUFBWSxLQUFLLE1BQU07QUFBQSxFQUMzQjtBQUlBLE1BQ0ksY0FBYyxhQUNYLFlBRUMsYUFFSSxjQUFjLGFBQ1gsTUFBTSxXQUFXLElBRzlCO0FBQ0UsVUFBTSx5QkFBeUI7QUFDL0IsVUFBTSxnQkFBZ0I7QUFDdEIsVUFBTSxlQUFlO0FBQUEsRUFDekI7QUFHQSxNQUFJLFdBQVcsR0FBRztBQUFFLGNBQVUsS0FBSztBQUFBLEVBQUc7QUFFdEMsTUFBSSxVQUFVLEdBQUc7QUFBRSxnQkFBWSxLQUFLO0FBQUEsRUFBRztBQUd2QyxNQUFJLGNBQWMsV0FBVztBQUFFLGdCQUFZLEtBQUs7QUFBQSxFQUFHO0FBQUM7QUFDeEQ7QUFFQSxTQUFTLFlBQVksT0FBeUI7QUFFMUMsWUFBVTtBQUNWLGNBQVk7QUFHWixNQUFJLENBQUMsVUFBVSxHQUFHO0FBQ2QsUUFBSSxNQUFNLFNBQVMsZUFBZSxNQUFNLFdBQVcsS0FBSyxNQUFNLFdBQVcsR0FBRztBQUN4RTtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBRUEsTUFBSSxZQUFZO0FBRVosZ0JBQVk7QUFFWjtBQUFBLEVBQ0o7QUFHQSxRQUFNLFNBQVMsWUFBWSxLQUFLO0FBSWhDLFFBQU0sUUFBUSxPQUFPLGlCQUFpQixNQUFNO0FBQzVDLFlBQ0ksTUFBTSxpQkFBaUIsbUJBQW1CLEVBQUUsS0FBSyxNQUFNLFdBRW5ELE1BQU0sVUFBVSxXQUFXLE1BQU0sV0FBVyxJQUFJLE9BQU8sZUFDcEQsTUFBTSxVQUFVLFdBQVcsTUFBTSxVQUFVLElBQUksT0FBTztBQUdyRTtBQUVBLFNBQVMsVUFBVSxPQUFtQjtBQUVsQyxZQUFVO0FBQ1YsYUFBVztBQUNYLGNBQVk7QUFDWixhQUFXO0FBQ2Y7QUFFQSxJQUFNLGdCQUFnQixPQUFPLE9BQU87QUFBQSxFQUNoQyxhQUFhO0FBQUEsRUFDYixhQUFhO0FBQUEsRUFDYixhQUFhO0FBQUEsRUFDYixhQUFhO0FBQUEsRUFDYixZQUFZO0FBQUEsRUFDWixZQUFZO0FBQUEsRUFDWixZQUFZO0FBQUEsRUFDWixZQUFZO0FBQ2hCLENBQUM7QUFFRCxTQUFTLFVBQVUsTUFBeUM7QUFDeEQsTUFBSSxNQUFNO0FBQ04sUUFBSSxDQUFDLFlBQVk7QUFBRSxzQkFBZ0IsU0FBUyxLQUFLLE1BQU07QUFBQSxJQUFRO0FBQy9ELGFBQVMsS0FBSyxNQUFNLFNBQVMsY0FBYyxJQUFJO0FBQUEsRUFDbkQsV0FBVyxDQUFDLFFBQVEsWUFBWTtBQUM1QixhQUFTLEtBQUssTUFBTSxTQUFTO0FBQUEsRUFDakM7QUFFQSxlQUFhLFFBQVE7QUFDekI7QUFFQSxTQUFTLFlBQVksT0FBeUI7QUFDMUMsTUFBSSxhQUFhLFlBQVk7QUFFekIsZUFBVztBQUNYLFdBQU8sa0JBQWtCLFVBQVU7QUFBQSxFQUN2QyxXQUFXLFNBQVM7QUFFaEIsZUFBVztBQUNYLFdBQU8sWUFBWTtBQUFBLEVBQ3ZCO0FBRUEsTUFBSSxZQUFZLFVBQVU7QUFHdEIsY0FBVSxZQUFZO0FBQ3RCO0FBQUEsRUFDSjtBQUVBLE1BQUksQ0FBQyxhQUFhLENBQUMsVUFBVSxHQUFHO0FBQzVCLFFBQUksWUFBWTtBQUFFLGdCQUFVO0FBQUEsSUFBRztBQUMvQjtBQUFBLEVBQ0o7QUFFQSxRQUFNLHFCQUFxQixRQUFRLDJCQUEyQixLQUFLO0FBQ25FLFFBQU0sb0JBQW9CLFFBQVEsMEJBQTBCLEtBQUs7QUFHakUsUUFBTSxjQUFjLFFBQVEsbUJBQW1CLEtBQUs7QUFFcEQsUUFBTSxjQUFlLE9BQU8sYUFBYSxNQUFNLFVBQVc7QUFDMUQsUUFBTSxhQUFhLE1BQU0sVUFBVTtBQUNuQyxRQUFNLFlBQVksTUFBTSxVQUFVO0FBQ2xDLFFBQU0sZUFBZ0IsT0FBTyxjQUFjLE1BQU0sVUFBVztBQUc1RCxRQUFNLGNBQWUsT0FBTyxhQUFhLE1BQU0sVUFBWSxvQkFBb0I7QUFDL0UsUUFBTSxhQUFhLE1BQU0sVUFBVyxvQkFBb0I7QUFDeEQsUUFBTSxZQUFZLE1BQU0sVUFBVyxxQkFBcUI7QUFDeEQsUUFBTSxlQUFnQixPQUFPLGNBQWMsTUFBTSxVQUFZLHFCQUFxQjtBQUVsRixNQUFJLENBQUMsY0FBYyxDQUFDLGFBQWEsQ0FBQyxnQkFBZ0IsQ0FBQyxhQUFhO0FBRTVELGNBQVU7QUFBQSxFQUNkLFdBRVMsZUFBZSxhQUFjLFdBQVUsV0FBVztBQUFBLFdBQ2xELGNBQWMsYUFBYyxXQUFVLFdBQVc7QUFBQSxXQUNqRCxjQUFjLFVBQVcsV0FBVSxXQUFXO0FBQUEsV0FDOUMsYUFBYSxZQUFhLFdBQVUsV0FBVztBQUFBLFdBRS9DLFdBQVksV0FBVSxVQUFVO0FBQUEsV0FDaEMsVUFBVyxXQUFVLFVBQVU7QUFBQSxXQUMvQixhQUFjLFdBQVUsVUFBVTtBQUFBLFdBQ2xDLFlBQWEsV0FBVSxVQUFVO0FBQUEsTUFFckMsV0FBVTtBQUNuQjs7O0FDNU9BO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQVdBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksV0FBVztBQUVyRCxJQUFNQyxjQUFhO0FBQ25CLElBQU1DLGNBQWE7QUFDbkIsSUFBTSxhQUFhO0FBS1osU0FBUyxPQUFzQjtBQUNsQyxTQUFPRixNQUFLQyxXQUFVO0FBQzFCO0FBS08sU0FBUyxPQUFzQjtBQUNsQyxTQUFPRCxNQUFLRSxXQUFVO0FBQzFCO0FBS08sU0FBUyxPQUFzQjtBQUNsQyxTQUFPRixNQUFLLFVBQVU7QUFDMUI7OztBQ3BDQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDd0JBLElBQUksVUFBVSxTQUFTLFVBQVU7QUFDakMsSUFBSSxlQUFvRCxPQUFPLFlBQVksWUFBWSxZQUFZLFFBQVEsUUFBUTtBQUNuSCxJQUFJO0FBQ0osSUFBSTtBQUNKLElBQUksT0FBTyxpQkFBaUIsY0FBYyxPQUFPLE9BQU8sbUJBQW1CLFlBQVk7QUFDbkYsTUFBSTtBQUNBLG1CQUFlLE9BQU8sZUFBZSxDQUFDLEdBQUcsVUFBVTtBQUFBLE1BQy9DLEtBQUssV0FBWTtBQUNiLGNBQU07QUFBQSxNQUNWO0FBQUEsSUFDSixDQUFDO0FBQ0QsdUJBQW1CLENBQUM7QUFFcEIsaUJBQWEsV0FBWTtBQUFFLFlBQU07QUFBQSxJQUFJLEdBQUcsTUFBTSxZQUFZO0FBQUEsRUFDOUQsU0FBUyxHQUFHO0FBQ1IsUUFBSSxNQUFNLGtCQUFrQjtBQUN4QixxQkFBZTtBQUFBLElBQ25CO0FBQUEsRUFDSjtBQUNKLE9BQU87QUFDSCxpQkFBZTtBQUNuQjtBQUVBLElBQUksbUJBQW1CO0FBQ3ZCLElBQUksZUFBZSxTQUFTLG1CQUFtQixPQUFxQjtBQUNoRSxNQUFJO0FBQ0EsUUFBSSxRQUFRLFFBQVEsS0FBSyxLQUFLO0FBQzlCLFdBQU8saUJBQWlCLEtBQUssS0FBSztBQUFBLEVBQ3RDLFNBQVMsR0FBRztBQUNSLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFFQSxJQUFJLG9CQUFvQixTQUFTLGlCQUFpQixPQUFxQjtBQUNuRSxNQUFJO0FBQ0EsUUFBSSxhQUFhLEtBQUssR0FBRztBQUFFLGFBQU87QUFBQSxJQUFPO0FBQ3pDLFlBQVEsS0FBSyxLQUFLO0FBQ2xCLFdBQU87QUFBQSxFQUNYLFNBQVMsR0FBRztBQUNSLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFDQSxJQUFJLFFBQVEsT0FBTyxVQUFVO0FBQzdCLElBQUksY0FBYztBQUNsQixJQUFJLFVBQVU7QUFDZCxJQUFJLFdBQVc7QUFDZixJQUFJLFdBQVc7QUFDZixJQUFJLFlBQVk7QUFDaEIsSUFBSSxZQUFZO0FBQ2hCLElBQUksaUJBQWlCLE9BQU8sV0FBVyxjQUFjLENBQUMsQ0FBQyxPQUFPO0FBRTlELElBQUksU0FBUyxFQUFFLEtBQUssQ0FBQyxDQUFDO0FBRXRCLElBQUksUUFBaUMsU0FBUyxtQkFBbUI7QUFBRSxTQUFPO0FBQU87QUFDakYsSUFBSSxPQUFPLGFBQWEsVUFBVTtBQUUxQixRQUFNLFNBQVM7QUFDbkIsTUFBSSxNQUFNLEtBQUssR0FBRyxNQUFNLE1BQU0sS0FBSyxTQUFTLEdBQUcsR0FBRztBQUM5QyxZQUFRLFNBQVNHLGtCQUFpQixPQUFPO0FBR3JDLFdBQUssVUFBVSxDQUFDLFdBQVcsT0FBTyxVQUFVLGVBQWUsT0FBTyxVQUFVLFdBQVc7QUFDbkYsWUFBSTtBQUNBLGNBQUksTUFBTSxNQUFNLEtBQUssS0FBSztBQUMxQixrQkFDSSxRQUFRLFlBQ0wsUUFBUSxhQUNSLFFBQVEsYUFDUixRQUFRLGdCQUNWLE1BQU0sRUFBRSxLQUFLO0FBQUEsUUFDdEIsU0FBUyxHQUFHO0FBQUEsUUFBTztBQUFBLE1BQ3ZCO0FBQ0EsYUFBTztBQUFBLElBQ1g7QUFBQSxFQUNKO0FBQ0o7QUFuQlE7QUFxQlIsU0FBUyxtQkFBc0IsT0FBdUQ7QUFDbEYsTUFBSSxNQUFNLEtBQUssR0FBRztBQUFFLFdBQU87QUFBQSxFQUFNO0FBQ2pDLE1BQUksQ0FBQyxPQUFPO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDNUIsTUFBSSxPQUFPLFVBQVUsY0FBYyxPQUFPLFVBQVUsVUFBVTtBQUFFLFdBQU87QUFBQSxFQUFPO0FBQzlFLE1BQUk7QUFDQSxJQUFDLGFBQXFCLE9BQU8sTUFBTSxZQUFZO0FBQUEsRUFDbkQsU0FBUyxHQUFHO0FBQ1IsUUFBSSxNQUFNLGtCQUFrQjtBQUFFLGFBQU87QUFBQSxJQUFPO0FBQUEsRUFDaEQ7QUFDQSxTQUFPLENBQUMsYUFBYSxLQUFLLEtBQUssa0JBQWtCLEtBQUs7QUFDMUQ7QUFFQSxTQUFTLHFCQUF3QixPQUFzRDtBQUNuRixNQUFJLE1BQU0sS0FBSyxHQUFHO0FBQUUsV0FBTztBQUFBLEVBQU07QUFDakMsTUFBSSxDQUFDLE9BQU87QUFBRSxXQUFPO0FBQUEsRUFBTztBQUM1QixNQUFJLE9BQU8sVUFBVSxjQUFjLE9BQU8sVUFBVSxVQUFVO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDOUUsTUFBSSxnQkFBZ0I7QUFBRSxXQUFPLGtCQUFrQixLQUFLO0FBQUEsRUFBRztBQUN2RCxNQUFJLGFBQWEsS0FBSyxHQUFHO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDekMsTUFBSSxXQUFXLE1BQU0sS0FBSyxLQUFLO0FBQy9CLE1BQUksYUFBYSxXQUFXLGFBQWEsWUFBWSxDQUFFLGlCQUFrQixLQUFLLFFBQVEsR0FBRztBQUFFLFdBQU87QUFBQSxFQUFPO0FBQ3pHLFNBQU8sa0JBQWtCLEtBQUs7QUFDbEM7QUFFQSxJQUFPLG1CQUFRLGVBQWUscUJBQXFCOzs7QUN6RzVDLElBQU0sY0FBTixjQUEwQixNQUFNO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBTW5DLFlBQVksU0FBa0IsU0FBd0I7QUFDbEQsVUFBTSxTQUFTLE9BQU87QUFDdEIsU0FBSyxPQUFPO0FBQUEsRUFDaEI7QUFDSjtBQWNPLElBQU0sMEJBQU4sY0FBc0MsTUFBTTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFhL0MsWUFBWSxTQUFzQyxRQUFjLE1BQWU7QUFDM0UsV0FBTyxzQkFBUSwrQ0FBK0MsY0FBYyxhQUFhLE1BQU0sR0FBRyxFQUFFLE9BQU8sT0FBTyxDQUFDO0FBQ25ILFNBQUssVUFBVTtBQUNmLFNBQUssT0FBTztBQUFBLEVBQ2hCO0FBQ0o7QUErQkEsSUFBTSxhQUFhLE9BQU8sU0FBUztBQUNuQyxJQUFNLGdCQUFnQixPQUFPLFlBQVk7QUE3RnpDO0FBOEZBLElBQU0sV0FBVSxZQUFPLFlBQVAsWUFBa0IsT0FBTyxpQkFBaUI7QUFvRG5ELElBQU0scUJBQU4sTUFBTSw0QkFBOEIsUUFBZ0U7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUF1Q3ZHLFlBQVksVUFBeUMsYUFBMkM7QUFDNUYsUUFBSTtBQUNKLFFBQUk7QUFDSixVQUFNLENBQUMsS0FBSyxRQUFRO0FBQUUsZ0JBQVU7QUFBSyxlQUFTO0FBQUEsSUFBSyxDQUFDO0FBRXBELFFBQUssS0FBSyxZQUFvQixPQUFPLE1BQU0sU0FBUztBQUNoRCxZQUFNLElBQUksVUFBVSxtSUFBbUk7QUFBQSxJQUMzSjtBQUVBLFFBQUksVUFBOEM7QUFBQSxNQUM5QyxTQUFTO0FBQUEsTUFDVDtBQUFBLE1BQ0E7QUFBQSxNQUNBLElBQUksY0FBYztBQUFFLGVBQU8sb0NBQWU7QUFBQSxNQUFNO0FBQUEsTUFDaEQsSUFBSSxZQUFZLElBQUk7QUFBRSxzQkFBYyxrQkFBTTtBQUFBLE1BQVc7QUFBQSxJQUN6RDtBQUVBLFVBQU0sUUFBaUM7QUFBQSxNQUNuQyxJQUFJLE9BQU87QUFBRSxlQUFPO0FBQUEsTUFBTztBQUFBLE1BQzNCLFdBQVc7QUFBQSxNQUNYLFNBQVM7QUFBQSxJQUNiO0FBR0EsU0FBSyxPQUFPLGlCQUFpQixNQUFNO0FBQUEsTUFDL0IsQ0FBQyxVQUFVLEdBQUc7QUFBQSxRQUNWLGNBQWM7QUFBQSxRQUNkLFlBQVk7QUFBQSxRQUNaLFVBQVU7QUFBQSxRQUNWLE9BQU87QUFBQSxNQUNYO0FBQUEsTUFDQSxDQUFDLGFBQWEsR0FBRztBQUFBLFFBQ2IsY0FBYztBQUFBLFFBQ2QsWUFBWTtBQUFBLFFBQ1osVUFBVTtBQUFBLFFBQ1YsT0FBTyxhQUFhLFNBQVMsS0FBSztBQUFBLE1BQ3RDO0FBQUEsSUFDSixDQUFDO0FBR0QsVUFBTSxXQUFXLFlBQVksU0FBUyxLQUFLO0FBQzNDLFFBQUk7QUFDQSxlQUFTLFlBQVksU0FBUyxLQUFLLEdBQUcsUUFBUTtBQUFBLElBQ2xELFNBQVMsS0FBSztBQUNWLFVBQUksTUFBTSxXQUFXO0FBQ2pCLGdCQUFRLElBQUksdURBQXVELEdBQUc7QUFBQSxNQUMxRSxPQUFPO0FBQ0gsaUJBQVMsR0FBRztBQUFBLE1BQ2hCO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBeURBLE9BQU8sT0FBdUM7QUFDMUMsV0FBTyxJQUFJLG9CQUF5QixDQUFDLFlBQVk7QUFHN0MsY0FBUSxJQUFJO0FBQUEsUUFDUixLQUFLLGFBQWEsRUFBRSxJQUFJLFlBQVksc0JBQXNCLEVBQUUsTUFBTSxDQUFDLENBQUM7QUFBQSxRQUNwRSxlQUFlLElBQUk7QUFBQSxNQUN2QixDQUFDLEVBQUUsS0FBSyxNQUFNLFFBQVEsR0FBRyxNQUFNLFFBQVEsQ0FBQztBQUFBLElBQzVDLENBQUM7QUFBQSxFQUNMO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQTJCQSxTQUFTLFFBQTRDO0FBQ2pELFFBQUksT0FBTyxTQUFTO0FBQ2hCLFdBQUssS0FBSyxPQUFPLE9BQU8sTUFBTTtBQUFBLElBQ2xDLE9BQU87QUFDSCxhQUFPLGlCQUFpQixTQUFTLE1BQU0sS0FBSyxLQUFLLE9BQU8sT0FBTyxNQUFNLEdBQUcsRUFBQyxTQUFTLEtBQUksQ0FBQztBQUFBLElBQzNGO0FBRUEsV0FBTztBQUFBLEVBQ1g7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUErQkEsS0FBcUMsYUFBc0gsWUFBd0gsYUFBb0Y7QUFDblcsUUFBSSxFQUFFLGdCQUFnQixzQkFBcUI7QUFDdkMsWUFBTSxJQUFJLFVBQVUsZ0VBQWdFO0FBQUEsSUFDeEY7QUFNQSxRQUFJLENBQUMsaUJBQVcsV0FBVyxHQUFHO0FBQUUsb0JBQWM7QUFBQSxJQUFpQjtBQUMvRCxRQUFJLENBQUMsaUJBQVcsVUFBVSxHQUFHO0FBQUUsbUJBQWE7QUFBQSxJQUFTO0FBRXJELFFBQUksZ0JBQWdCLFlBQVksY0FBYyxTQUFTO0FBRW5ELGFBQU8sSUFBSSxvQkFBbUIsQ0FBQyxZQUFZLFFBQVEsSUFBVyxDQUFDO0FBQUEsSUFDbkU7QUFFQSxVQUFNLFVBQStDLENBQUM7QUFDdEQsU0FBSyxVQUFVLElBQUk7QUFFbkIsV0FBTyxJQUFJLG9CQUF3QyxDQUFDLFNBQVMsV0FBVztBQUNwRSxXQUFLLE1BQU07QUFBQSxRQUNQLENBQUMsVUFBVTtBQXJZM0IsY0FBQUM7QUFzWW9CLGNBQUksS0FBSyxVQUFVLE1BQU0sU0FBUztBQUFFLGlCQUFLLFVBQVUsSUFBSTtBQUFBLFVBQU07QUFDN0QsV0FBQUEsTUFBQSxRQUFRLFlBQVIsZ0JBQUFBLElBQUE7QUFFQSxjQUFJO0FBQ0Esb0JBQVEsWUFBYSxLQUFLLENBQUM7QUFBQSxVQUMvQixTQUFTLEtBQUs7QUFDVixtQkFBTyxHQUFHO0FBQUEsVUFDZDtBQUFBLFFBQ0o7QUFBQSxRQUNBLENBQUMsV0FBWTtBQS9ZN0IsY0FBQUE7QUFnWm9CLGNBQUksS0FBSyxVQUFVLE1BQU0sU0FBUztBQUFFLGlCQUFLLFVBQVUsSUFBSTtBQUFBLFVBQU07QUFDN0QsV0FBQUEsTUFBQSxRQUFRLFlBQVIsZ0JBQUFBLElBQUE7QUFFQSxjQUFJO0FBQ0Esb0JBQVEsV0FBWSxNQUFNLENBQUM7QUFBQSxVQUMvQixTQUFTLEtBQUs7QUFDVixtQkFBTyxHQUFHO0FBQUEsVUFDZDtBQUFBLFFBQ0o7QUFBQSxNQUNKO0FBQUEsSUFDSixHQUFHLE9BQU8sVUFBVztBQUVqQixVQUFJO0FBQ0EsZUFBTywyQ0FBYztBQUFBLE1BQ3pCLFVBQUU7QUFDRSxjQUFNLEtBQUssT0FBTyxLQUFLO0FBQUEsTUFDM0I7QUFBQSxJQUNKLENBQUM7QUFBQSxFQUNMO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBK0JBLE1BQXVCLFlBQXFGLGFBQTRFO0FBQ3BMLFdBQU8sS0FBSyxLQUFLLFFBQVcsWUFBWSxXQUFXO0FBQUEsRUFDdkQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBaUNBLFFBQVEsV0FBNkMsYUFBa0U7QUFDbkgsUUFBSSxFQUFFLGdCQUFnQixzQkFBcUI7QUFDdkMsWUFBTSxJQUFJLFVBQVUsbUVBQW1FO0FBQUEsSUFDM0Y7QUFFQSxRQUFJLENBQUMsaUJBQVcsU0FBUyxHQUFHO0FBQ3hCLGFBQU8sS0FBSyxLQUFLLFdBQVcsV0FBVyxXQUFXO0FBQUEsSUFDdEQ7QUFFQSxXQUFPLEtBQUs7QUFBQSxNQUNSLENBQUMsVUFBVSxvQkFBbUIsUUFBUSxVQUFVLENBQUMsRUFBRSxLQUFLLE1BQU0sS0FBSztBQUFBLE1BQ25FLENBQUMsV0FBWSxvQkFBbUIsUUFBUSxVQUFVLENBQUMsRUFBRSxLQUFLLE1BQU07QUFBRSxjQUFNO0FBQUEsTUFBUSxDQUFDO0FBQUEsTUFDakY7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFZQSxhQXpXUyxZQUVTLGVBdVdOLFFBQU8sSUFBSTtBQUNuQixXQUFPO0FBQUEsRUFDWDtBQUFBLEVBYUEsT0FBTyxJQUFzRCxRQUF3QztBQUNqRyxRQUFJLFlBQVksTUFBTSxLQUFLLE1BQU07QUFDakMsVUFBTSxVQUFVLFVBQVUsV0FBVyxJQUMvQixvQkFBbUIsUUFBUSxTQUFTLElBQ3BDLElBQUksb0JBQTRCLENBQUMsU0FBUyxXQUFXO0FBQ25ELFdBQUssUUFBUSxJQUFJLFNBQVMsRUFBRSxLQUFLLFNBQVMsTUFBTTtBQUFBLElBQ3BELEdBQUcsQ0FBQyxVQUEwQixVQUFVLFNBQVMsV0FBVyxLQUFLLENBQUM7QUFDdEUsV0FBTztBQUFBLEVBQ1g7QUFBQSxFQWFBLE9BQU8sV0FBNkQsUUFBd0M7QUFDeEcsUUFBSSxZQUFZLE1BQU0sS0FBSyxNQUFNO0FBQ2pDLFVBQU0sVUFBVSxVQUFVLFdBQVcsSUFDL0Isb0JBQW1CLFFBQVEsU0FBUyxJQUNwQyxJQUFJLG9CQUE0QixDQUFDLFNBQVMsV0FBVztBQUNuRCxXQUFLLFFBQVEsV0FBVyxTQUFTLEVBQUUsS0FBSyxTQUFTLE1BQU07QUFBQSxJQUMzRCxHQUFHLENBQUMsVUFBMEIsVUFBVSxTQUFTLFdBQVcsS0FBSyxDQUFDO0FBQ3RFLFdBQU87QUFBQSxFQUNYO0FBQUEsRUFlQSxPQUFPLElBQXNELFFBQXdDO0FBQ2pHLFFBQUksWUFBWSxNQUFNLEtBQUssTUFBTTtBQUNqQyxVQUFNLFVBQVUsVUFBVSxXQUFXLElBQy9CLG9CQUFtQixRQUFRLFNBQVMsSUFDcEMsSUFBSSxvQkFBNEIsQ0FBQyxTQUFTLFdBQVc7QUFDbkQsV0FBSyxRQUFRLElBQUksU0FBUyxFQUFFLEtBQUssU0FBUyxNQUFNO0FBQUEsSUFDcEQsR0FBRyxDQUFDLFVBQTBCLFVBQVUsU0FBUyxXQUFXLEtBQUssQ0FBQztBQUN0RSxXQUFPO0FBQUEsRUFDWDtBQUFBLEVBWUEsT0FBTyxLQUF1RCxRQUF3QztBQUNsRyxRQUFJLFlBQVksTUFBTSxLQUFLLE1BQU07QUFDakMsVUFBTSxVQUFVLElBQUksb0JBQTRCLENBQUMsU0FBUyxXQUFXO0FBQ2pFLFdBQUssUUFBUSxLQUFLLFNBQVMsRUFBRSxLQUFLLFNBQVMsTUFBTTtBQUFBLElBQ3JELEdBQUcsQ0FBQyxVQUEwQixVQUFVLFNBQVMsV0FBVyxLQUFLLENBQUM7QUFDbEUsV0FBTztBQUFBLEVBQ1g7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxPQUFPLE9BQWtCLE9BQW9DO0FBQ3pELFVBQU0sSUFBSSxJQUFJLG9CQUFzQixNQUFNO0FBQUEsSUFBQyxDQUFDO0FBQzVDLE1BQUUsT0FBTyxLQUFLO0FBQ2QsV0FBTztBQUFBLEVBQ1g7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBWUEsT0FBTyxRQUFtQixjQUFzQixPQUFvQztBQUNoRixVQUFNLFVBQVUsSUFBSSxvQkFBc0IsTUFBTTtBQUFBLElBQUMsQ0FBQztBQUNsRCxRQUFJLGVBQWUsT0FBTyxnQkFBZ0IsY0FBYyxZQUFZLFdBQVcsT0FBTyxZQUFZLFlBQVksWUFBWTtBQUN0SCxrQkFBWSxRQUFRLFlBQVksRUFBRSxpQkFBaUIsU0FBUyxNQUFNLEtBQUssUUFBUSxPQUFPLEtBQUssQ0FBQztBQUFBLElBQ2hHLE9BQU87QUFDSCxpQkFBVyxNQUFNLEtBQUssUUFBUSxPQUFPLEtBQUssR0FBRyxZQUFZO0FBQUEsSUFDN0Q7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUFBLEVBaUJBLE9BQU8sTUFBZ0IsY0FBc0IsT0FBa0M7QUFDM0UsV0FBTyxJQUFJLG9CQUFzQixDQUFDLFlBQVk7QUFDMUMsaUJBQVcsTUFBTSxRQUFRLEtBQU0sR0FBRyxZQUFZO0FBQUEsSUFDbEQsQ0FBQztBQUFBLEVBQ0w7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxPQUFPLE9BQWtCLFFBQXFDO0FBQzFELFdBQU8sSUFBSSxvQkFBc0IsQ0FBQyxHQUFHLFdBQVcsT0FBTyxNQUFNLENBQUM7QUFBQSxFQUNsRTtBQUFBLEVBb0JBLE9BQU8sUUFBa0IsT0FBNEQ7QUFDakYsUUFBSSxpQkFBaUIscUJBQW9CO0FBRXJDLGFBQU87QUFBQSxJQUNYO0FBQ0EsV0FBTyxJQUFJLG9CQUF3QixDQUFDLFlBQVksUUFBUSxLQUFLLENBQUM7QUFBQSxFQUNsRTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVVBLE9BQU8sZ0JBQXVEO0FBQzFELFFBQUksU0FBNkMsRUFBRSxhQUFhLEtBQUs7QUFDckUsV0FBTyxVQUFVLElBQUksb0JBQXNCLENBQUMsU0FBUyxXQUFXO0FBQzVELGFBQU8sVUFBVTtBQUNqQixhQUFPLFNBQVM7QUFBQSxJQUNwQixHQUFHLENBQUMsVUFBZ0I7QUF6ckI1QixVQUFBQTtBQXlyQjhCLE9BQUFBLE1BQUEsT0FBTyxnQkFBUCxnQkFBQUEsSUFBQSxhQUFxQjtBQUFBLElBQVEsQ0FBQztBQUNwRCxXQUFPO0FBQUEsRUFDWDtBQUNKO0FBTUEsU0FBUyxhQUFnQixTQUE2QyxPQUFnQztBQUNsRyxNQUFJLHNCQUFnRDtBQUVwRCxTQUFPLENBQUMsV0FBa0Q7QUFDdEQsUUFBSSxDQUFDLE1BQU0sU0FBUztBQUNoQixZQUFNLFVBQVU7QUFDaEIsWUFBTSxTQUFTO0FBQ2YsY0FBUSxPQUFPLE1BQU07QUFNckIsV0FBSyxRQUFRLFVBQVUsS0FBSyxLQUFLLFFBQVEsU0FBUyxRQUFXLENBQUMsUUFBUTtBQUNsRSxZQUFJLFFBQVEsUUFBUTtBQUNoQixnQkFBTTtBQUFBLFFBQ1Y7QUFBQSxNQUNKLENBQUM7QUFBQSxJQUNMO0FBSUEsUUFBSSxDQUFDLE1BQU0sVUFBVSxDQUFDLFFBQVEsYUFBYTtBQUFFO0FBQUEsSUFBUTtBQUVyRCwwQkFBc0IsSUFBSSxRQUFjLENBQUMsWUFBWTtBQUNqRCxVQUFJO0FBQ0EsZ0JBQVEsUUFBUSxZQUFhLE1BQU0sT0FBUSxLQUFLLENBQUM7QUFBQSxNQUNyRCxTQUFTLEtBQUs7QUFDVixnQkFBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsU0FBUyxLQUFLLDhDQUE4QyxDQUFDO0FBQUEsTUFDcEg7QUFBQSxJQUNKLENBQUMsRUFBRSxNQUFNLENBQUNDLFlBQVk7QUFDbEIsY0FBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsU0FBU0EsU0FBUSw4Q0FBOEMsQ0FBQztBQUFBLElBQ3ZILENBQUM7QUFHRCxZQUFRLGNBQWM7QUFFdEIsV0FBTztBQUFBLEVBQ1g7QUFDSjtBQUtBLFNBQVMsWUFBZSxTQUE2QyxPQUErRDtBQUNoSSxTQUFPLENBQUMsVUFBVTtBQUNkLFFBQUksTUFBTSxXQUFXO0FBQUU7QUFBQSxJQUFRO0FBQy9CLFVBQU0sWUFBWTtBQUVsQixRQUFJLFVBQVUsUUFBUSxTQUFTO0FBQzNCLFVBQUksTUFBTSxTQUFTO0FBQUU7QUFBQSxNQUFRO0FBQzdCLFlBQU0sVUFBVTtBQUNoQixjQUFRLE9BQU8sSUFBSSxVQUFVLDJDQUEyQyxDQUFDO0FBQ3pFO0FBQUEsSUFDSjtBQUVBLFFBQUksU0FBUyxTQUFTLE9BQU8sVUFBVSxZQUFZLE9BQU8sVUFBVSxhQUFhO0FBQzdFLFVBQUk7QUFDSixVQUFJO0FBQ0EsZUFBUSxNQUFjO0FBQUEsTUFDMUIsU0FBUyxLQUFLO0FBQ1YsY0FBTSxVQUFVO0FBQ2hCLGdCQUFRLE9BQU8sR0FBRztBQUNsQjtBQUFBLE1BQ0o7QUFFQSxVQUFJLGlCQUFXLElBQUksR0FBRztBQUNsQixZQUFJO0FBQ0EsY0FBSSxTQUFVLE1BQWM7QUFDNUIsY0FBSSxpQkFBVyxNQUFNLEdBQUc7QUFDcEIsa0JBQU0sY0FBYyxDQUFDLFVBQWdCO0FBQ2pDLHNCQUFRLE1BQU0sUUFBUSxPQUFPLENBQUMsS0FBSyxDQUFDO0FBQUEsWUFDeEM7QUFDQSxnQkFBSSxNQUFNLFFBQVE7QUFJZCxtQkFBSyxhQUFhLGlDQUFLLFVBQUwsRUFBYyxZQUFZLElBQUcsS0FBSyxFQUFFLE1BQU0sTUFBTTtBQUFBLFlBQ3RFLE9BQU87QUFDSCxzQkFBUSxjQUFjO0FBQUEsWUFDMUI7QUFBQSxVQUNKO0FBQUEsUUFDSixTQUFRO0FBQUEsUUFBQztBQUVULGNBQU0sV0FBb0M7QUFBQSxVQUN0QyxNQUFNLE1BQU07QUFBQSxVQUNaLFdBQVc7QUFBQSxVQUNYLElBQUksVUFBVTtBQUFFLG1CQUFPLEtBQUssS0FBSztBQUFBLFVBQVE7QUFBQSxVQUN6QyxJQUFJLFFBQVFDLFFBQU87QUFBRSxpQkFBSyxLQUFLLFVBQVVBO0FBQUEsVUFBTztBQUFBLFVBQ2hELElBQUksU0FBUztBQUFFLG1CQUFPLEtBQUssS0FBSztBQUFBLFVBQU87QUFBQSxRQUMzQztBQUVBLGNBQU0sV0FBVyxZQUFZLFNBQVMsUUFBUTtBQUM5QyxZQUFJO0FBQ0Esa0JBQVEsTUFBTSxNQUFNLE9BQU8sQ0FBQyxZQUFZLFNBQVMsUUFBUSxHQUFHLFFBQVEsQ0FBQztBQUFBLFFBQ3pFLFNBQVMsS0FBSztBQUNWLG1CQUFTLEdBQUc7QUFBQSxRQUNoQjtBQUNBO0FBQUEsTUFDSjtBQUFBLElBQ0o7QUFFQSxRQUFJLE1BQU0sU0FBUztBQUFFO0FBQUEsSUFBUTtBQUM3QixVQUFNLFVBQVU7QUFDaEIsWUFBUSxRQUFRLEtBQUs7QUFBQSxFQUN6QjtBQUNKO0FBS0EsU0FBUyxZQUFlLFNBQTZDLE9BQTREO0FBQzdILFNBQU8sQ0FBQyxXQUFZO0FBQ2hCLFFBQUksTUFBTSxXQUFXO0FBQUU7QUFBQSxJQUFRO0FBQy9CLFVBQU0sWUFBWTtBQUVsQixRQUFJLE1BQU0sU0FBUztBQUNmLFVBQUk7QUFDQSxZQUFJLGtCQUFrQixlQUFlLE1BQU0sa0JBQWtCLGVBQWUsT0FBTyxHQUFHLE9BQU8sT0FBTyxNQUFNLE9BQU8sS0FBSyxHQUFHO0FBRXJIO0FBQUEsUUFDSjtBQUFBLE1BQ0osU0FBUTtBQUFBLE1BQUM7QUFFVCxXQUFLLFFBQVEsT0FBTyxJQUFJLHdCQUF3QixRQUFRLFNBQVMsTUFBTSxDQUFDO0FBQUEsSUFDNUUsT0FBTztBQUNILFlBQU0sVUFBVTtBQUNoQixjQUFRLE9BQU8sTUFBTTtBQUFBLElBQ3pCO0FBQUEsRUFDSjtBQUNKO0FBTUEsU0FBUyxVQUFVLFFBQXFDLFFBQWUsT0FBNEI7QUFDL0YsUUFBTSxVQUFVLENBQUM7QUFFakIsYUFBVyxTQUFTLFFBQVE7QUFDeEIsUUFBSTtBQUNKLFFBQUk7QUFDQSxVQUFJLENBQUMsaUJBQVcsTUFBTSxJQUFJLEdBQUc7QUFBRTtBQUFBLE1BQVU7QUFDekMsZUFBUyxNQUFNO0FBQ2YsVUFBSSxDQUFDLGlCQUFXLE1BQU0sR0FBRztBQUFFO0FBQUEsTUFBVTtBQUFBLElBQ3pDLFNBQVE7QUFBRTtBQUFBLElBQVU7QUFFcEIsUUFBSTtBQUNKLFFBQUk7QUFDQSxlQUFTLFFBQVEsTUFBTSxRQUFRLE9BQU8sQ0FBQyxLQUFLLENBQUM7QUFBQSxJQUNqRCxTQUFTLEtBQUs7QUFDVixjQUFRLE9BQU8sSUFBSSx3QkFBd0IsUUFBUSxLQUFLLHVDQUF1QyxDQUFDO0FBQ2hHO0FBQUEsSUFDSjtBQUVBLFFBQUksQ0FBQyxRQUFRO0FBQUU7QUFBQSxJQUFVO0FBQ3pCLFlBQVE7QUFBQSxPQUNILGtCQUFrQixVQUFXLFNBQVMsUUFBUSxRQUFRLE1BQU0sR0FBRyxNQUFNLENBQUMsV0FBWTtBQUMvRSxnQkFBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsUUFBUSx1Q0FBdUMsQ0FBQztBQUFBLE1BQ3ZHLENBQUM7QUFBQSxJQUNMO0FBQUEsRUFDSjtBQUVBLFNBQU8sUUFBUSxJQUFJLE9BQU87QUFDOUI7QUFLQSxTQUFTLFNBQVksR0FBUztBQUMxQixTQUFPO0FBQ1g7QUFLQSxTQUFTLFFBQVEsUUFBcUI7QUFDbEMsUUFBTTtBQUNWO0FBS0EsU0FBUyxhQUFhLEtBQWtCO0FBQ3BDLE1BQUk7QUFDQSxRQUFJLGVBQWUsU0FBUyxPQUFPLFFBQVEsWUFBWSxJQUFJLGFBQWEsT0FBTyxVQUFVLFVBQVU7QUFDL0YsYUFBTyxLQUFLO0FBQUEsSUFDaEI7QUFBQSxFQUNKLFNBQVE7QUFBQSxFQUFDO0FBRVQsTUFBSTtBQUNBLFdBQU8sS0FBSyxVQUFVLEdBQUc7QUFBQSxFQUM3QixTQUFRO0FBQUEsRUFBQztBQUVULE1BQUk7QUFDQSxXQUFPLE9BQU8sVUFBVSxTQUFTLEtBQUssR0FBRztBQUFBLEVBQzdDLFNBQVE7QUFBQSxFQUFDO0FBRVQsU0FBTztBQUNYO0FBS0EsU0FBUyxlQUFrQixTQUErQztBQTk0QjFFLE1BQUFGO0FBKzRCSSxNQUFJLE9BQTJDQSxNQUFBLFFBQVEsVUFBVSxNQUFsQixPQUFBQSxNQUF1QixDQUFDO0FBQ3ZFLE1BQUksRUFBRSxhQUFhLE1BQU07QUFDckIsV0FBTyxPQUFPLEtBQUsscUJBQTJCLENBQUM7QUFBQSxFQUNuRDtBQUNBLE1BQUksUUFBUSxVQUFVLEtBQUssTUFBTTtBQUM3QixRQUFJLFFBQVM7QUFDYixZQUFRLFVBQVUsSUFBSTtBQUFBLEVBQzFCO0FBQ0EsU0FBTyxJQUFJO0FBQ2Y7QUFHQSxJQUFJLHVCQUF1QixRQUFRO0FBQ25DLElBQUksd0JBQXdCLE9BQU8seUJBQXlCLFlBQVk7QUFDcEUseUJBQXVCLHFCQUFxQixLQUFLLE9BQU87QUFDNUQsT0FBTztBQUNILHlCQUF1QixXQUF3QztBQUMzRCxRQUFJO0FBQ0osUUFBSTtBQUNKLFVBQU0sVUFBVSxJQUFJLFFBQVcsQ0FBQyxLQUFLLFFBQVE7QUFBRSxnQkFBVTtBQUFLLGVBQVM7QUFBQSxJQUFLLENBQUM7QUFDN0UsV0FBTyxFQUFFLFNBQVMsU0FBUyxPQUFPO0FBQUEsRUFDdEM7QUFDSjs7O0FGdDVCQSxPQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFDbEMsT0FBTyxPQUFPLG9CQUFvQjtBQUNsQyxPQUFPLE9BQU8sbUJBQW1CO0FBSWpDLElBQU1HLFFBQU8saUJBQWlCLFlBQVksSUFBSTtBQUM5QyxJQUFNLGFBQWEsaUJBQWlCLFlBQVksVUFBVTtBQUMxRCxJQUFNLGdCQUFnQixvQkFBSSxJQUE4QjtBQUV4RCxJQUFNLGNBQWM7QUFDcEIsSUFBTSxlQUFlO0FBMEJkLElBQU0sZUFBTixjQUEyQixNQUFNO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBTXBDLFlBQVksU0FBa0IsU0FBd0I7QUFDbEQsVUFBTSxTQUFTLE9BQU87QUFDdEIsU0FBSyxPQUFPO0FBQUEsRUFDaEI7QUFDSjtBQVNBLFNBQVMsY0FBYyxJQUFZLE1BQWMsUUFBdUI7QUFDcEUsUUFBTSxZQUFZQyxzQkFBcUIsRUFBRTtBQUN6QyxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLE1BQUksQ0FBQyxNQUFNO0FBQ1AsY0FBVSxRQUFRLE1BQVM7QUFBQSxFQUMvQixXQUFXLENBQUMsUUFBUTtBQUNoQixjQUFVLFFBQVEsSUFBSTtBQUFBLEVBQzFCLE9BQU87QUFDSCxRQUFJO0FBQ0EsZ0JBQVUsUUFBUSxLQUFLLE1BQU0sSUFBSSxDQUFDO0FBQUEsSUFDdEMsU0FBUyxLQUFVO0FBQ2YsZ0JBQVUsT0FBTyxJQUFJLFVBQVUsNkJBQTZCLElBQUksU0FBUyxFQUFFLE9BQU8sSUFBSSxDQUFDLENBQUM7QUFBQSxJQUM1RjtBQUFBLEVBQ0o7QUFDSjtBQVNBLFNBQVMsYUFBYSxJQUFZLE1BQWMsUUFBdUI7QUFDbkUsUUFBTSxZQUFZQSxzQkFBcUIsRUFBRTtBQUN6QyxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLE1BQUksQ0FBQyxRQUFRO0FBQ1QsY0FBVSxPQUFPLElBQUksTUFBTSxJQUFJLENBQUM7QUFBQSxFQUNwQyxPQUFPO0FBQ0gsUUFBSTtBQUNKLFFBQUk7QUFDQSxjQUFRLEtBQUssTUFBTSxJQUFJO0FBQUEsSUFDM0IsU0FBUyxLQUFVO0FBQ2YsZ0JBQVUsT0FBTyxJQUFJLFVBQVUsNEJBQTRCLElBQUksU0FBUyxFQUFFLE9BQU8sSUFBSSxDQUFDLENBQUM7QUFDdkY7QUFBQSxJQUNKO0FBRUEsUUFBSSxVQUF3QixDQUFDO0FBQzdCLFFBQUksTUFBTSxPQUFPO0FBQ2IsY0FBUSxRQUFRLE1BQU07QUFBQSxJQUMxQjtBQUVBLFFBQUk7QUFDSixZQUFRLE1BQU0sTUFBTTtBQUFBLE1BQ2hCLEtBQUs7QUFDRCxvQkFBWSxJQUFJLGVBQWUsTUFBTSxTQUFTLE9BQU87QUFDckQ7QUFBQSxNQUNKLEtBQUs7QUFDRCxvQkFBWSxJQUFJLFVBQVUsTUFBTSxTQUFTLE9BQU87QUFDaEQ7QUFBQSxNQUNKLEtBQUs7QUFDRCxvQkFBWSxJQUFJLGFBQWEsTUFBTSxTQUFTLE9BQU87QUFDbkQ7QUFBQSxNQUNKO0FBQ0ksb0JBQVksSUFBSSxNQUFNLE1BQU0sU0FBUyxPQUFPO0FBQzVDO0FBQUEsSUFDUjtBQUVBLGNBQVUsT0FBTyxTQUFTO0FBQUEsRUFDOUI7QUFDSjtBQVFBLFNBQVNBLHNCQUFxQixJQUEwQztBQUNwRSxRQUFNLFdBQVcsY0FBYyxJQUFJLEVBQUU7QUFDckMsZ0JBQWMsT0FBTyxFQUFFO0FBQ3ZCLFNBQU87QUFDWDtBQU9BLFNBQVNDLGNBQXFCO0FBQzFCLE1BQUk7QUFDSixLQUFHO0FBQ0MsYUFBUyxPQUFPO0FBQUEsRUFDcEIsU0FBUyxjQUFjLElBQUksTUFBTTtBQUNqQyxTQUFPO0FBQ1g7QUFjTyxTQUFTLEtBQUssU0FBK0M7QUFDaEUsUUFBTSxLQUFLQSxZQUFXO0FBRXRCLFFBQU0sU0FBUyxtQkFBbUIsY0FBbUI7QUFDckQsZ0JBQWMsSUFBSSxJQUFJLEVBQUUsU0FBUyxPQUFPLFNBQVMsUUFBUSxPQUFPLE9BQU8sQ0FBQztBQUV4RSxRQUFNLFVBQVVGLE1BQUssYUFBYSxPQUFPLE9BQU8sRUFBRSxXQUFXLEdBQUcsR0FBRyxPQUFPLENBQUM7QUFDM0UsTUFBSSxVQUFVO0FBRWQsVUFBUSxLQUFLLE1BQU07QUFDZixjQUFVO0FBQUEsRUFDZCxHQUFHLENBQUMsUUFBUTtBQUNSLGtCQUFjLE9BQU8sRUFBRTtBQUN2QixXQUFPLE9BQU8sR0FBRztBQUFBLEVBQ3JCLENBQUM7QUFFRCxRQUFNLFNBQVMsTUFBTTtBQUNqQixrQkFBYyxPQUFPLEVBQUU7QUFDdkIsV0FBTyxXQUFXLGNBQWMsRUFBQyxXQUFXLEdBQUUsQ0FBQyxFQUFFLE1BQU0sQ0FBQyxRQUFRO0FBQzVELGNBQVEsTUFBTSxxREFBcUQsR0FBRztBQUFBLElBQzFFLENBQUM7QUFBQSxFQUNMO0FBRUEsU0FBTyxjQUFjLE1BQU07QUFDdkIsUUFBSSxTQUFTO0FBQ1QsYUFBTyxPQUFPO0FBQUEsSUFDbEIsT0FBTztBQUNILGFBQU8sUUFBUSxLQUFLLE1BQU07QUFBQSxJQUM5QjtBQUFBLEVBQ0o7QUFFQSxTQUFPLE9BQU87QUFDbEI7QUFVTyxTQUFTLE9BQU8sZUFBdUIsTUFBc0M7QUFDaEYsU0FBTyxLQUFLLEVBQUUsWUFBWSxLQUFLLENBQUM7QUFDcEM7QUFVTyxTQUFTLEtBQUssYUFBcUIsTUFBc0M7QUFDNUUsU0FBTyxLQUFLLEVBQUUsVUFBVSxLQUFLLENBQUM7QUFDbEM7OztBR3hPQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBWUEsSUFBTUcsUUFBTyxpQkFBaUIsWUFBWSxTQUFTO0FBRW5ELElBQU0sbUJBQW1CO0FBQ3pCLElBQU0sZ0JBQWdCO0FBUWYsU0FBUyxRQUFRLE1BQTZCO0FBQ2pELFNBQU9BLE1BQUssa0JBQWtCLEVBQUMsS0FBSSxDQUFDO0FBQ3hDO0FBT08sU0FBUyxPQUF3QjtBQUNwQyxTQUFPQSxNQUFLLGFBQWE7QUFDN0I7OztBQ2xDQTtBQUFBO0FBQUE7QUFBQSxlQUFBQztBQUFBLEVBQUE7QUFBQSxhQUFBQztBQUFBLEVBQUE7QUFBQTtBQUFBO0FBYU8sU0FBUyxJQUFhLFFBQWdCO0FBQ3pDLFNBQU87QUFDWDtBQU1PLFNBQVMsVUFBVSxRQUFxQjtBQUMzQyxTQUFTLFVBQVUsT0FBUSxLQUFLO0FBQ3BDO0FBT08sU0FBU0MsT0FBZSxTQUFtRDtBQUM5RSxNQUFJLFlBQVksS0FBSztBQUNqQixXQUFPLENBQUMsV0FBWSxXQUFXLE9BQU8sQ0FBQyxJQUFJO0FBQUEsRUFDL0M7QUFFQSxTQUFPLENBQUMsV0FBVztBQUNmLFFBQUksV0FBVyxNQUFNO0FBQ2pCLGFBQU8sQ0FBQztBQUFBLElBQ1o7QUFDQSxhQUFTLElBQUksR0FBRyxJQUFJLE9BQU8sUUFBUSxLQUFLO0FBQ3BDLGFBQU8sQ0FBQyxJQUFJLFFBQVEsT0FBTyxDQUFDLENBQUM7QUFBQSxJQUNqQztBQUNBLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFPTyxTQUFTQyxLQUFhLEtBQThCLE9BQStEO0FBQ3RILE1BQUksVUFBVSxLQUFLO0FBQ2YsV0FBTyxDQUFDLFdBQVksV0FBVyxPQUFPLENBQUMsSUFBSTtBQUFBLEVBQy9DO0FBRUEsU0FBTyxDQUFDLFdBQVc7QUFDZixRQUFJLFdBQVcsTUFBTTtBQUNqQixhQUFPLENBQUM7QUFBQSxJQUNaO0FBQ0EsZUFBV0MsUUFBTyxRQUFRO0FBQ3RCLGFBQU9BLElBQUcsSUFBSSxNQUFNLE9BQU9BLElBQUcsQ0FBQztBQUFBLElBQ25DO0FBQ0EsV0FBTztBQUFBLEVBQ1g7QUFDSjtBQU1PLFNBQVMsU0FBa0IsU0FBMEQ7QUFDeEYsTUFBSSxZQUFZLEtBQUs7QUFDakIsV0FBTztBQUFBLEVBQ1g7QUFFQSxTQUFPLENBQUMsV0FBWSxXQUFXLE9BQU8sT0FBTyxRQUFRLE1BQU07QUFDL0Q7QUFNTyxTQUFTLE9BQU8sYUFFdkI7QUFDSSxNQUFJLFNBQVM7QUFDYixhQUFXLFFBQVEsYUFBYTtBQUM1QixRQUFJLFlBQVksSUFBSSxNQUFNLEtBQUs7QUFDM0IsZUFBUztBQUNUO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFDQSxNQUFJLFFBQVE7QUFDUixXQUFPO0FBQUEsRUFDWDtBQUVBLFNBQU8sQ0FBQyxXQUFXO0FBQ2YsZUFBVyxRQUFRLGFBQWE7QUFDNUIsVUFBSSxRQUFRLFFBQVE7QUFDaEIsZUFBTyxJQUFJLElBQUksWUFBWSxJQUFJLEVBQUUsT0FBTyxJQUFJLENBQUM7QUFBQSxNQUNqRDtBQUFBLElBQ0o7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUNKOzs7QUN6R0E7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBd0RBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksT0FBTztBQUVqRCxJQUFNLFNBQVM7QUFDZixJQUFNLGFBQWE7QUFDbkIsSUFBTSxhQUFhO0FBT1osU0FBUyxTQUE0QjtBQUN4QyxTQUFPQSxNQUFLLE1BQU07QUFDdEI7QUFPTyxTQUFTLGFBQThCO0FBQzFDLFNBQU9BLE1BQUssVUFBVTtBQUMxQjtBQU9PLFNBQVMsYUFBOEI7QUFDMUMsU0FBT0EsTUFBSyxVQUFVO0FBQzFCOzs7QUN0RE8sSUFBTSxrQkFBTixNQUFnRDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBSW5ELGVBQWUsTUFBVyxhQUE2QjtBQUNuRCxRQUFJLENBQUMsS0FBTSxRQUFPO0FBRWxCLFFBQUk7QUFFQSxZQUFNLGVBQWUsS0FBSyxJQUFJO0FBQzlCLFlBQU0sUUFBUSxJQUFJLFdBQVcsYUFBYSxNQUFNO0FBQ2hELGVBQVMsSUFBSSxHQUFHLElBQUksYUFBYSxRQUFRLEtBQUs7QUFDMUMsY0FBTSxDQUFDLElBQUksYUFBYSxXQUFXLENBQUM7QUFBQSxNQUN4QztBQUVBLGFBQU8sSUFBSSxZQUFZLEVBQUUsT0FBTyxLQUFLO0FBQUEsSUFDekMsU0FBUyxLQUFLO0FBQ1YsWUFBTSxJQUFJLE1BQU0sOEJBQThCLHNCQUFlLFFBQVEsSUFBSSxVQUFVLE9BQU8sR0FBRyxFQUFHO0FBQUEsSUFDcEc7QUFBQSxFQUNKO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxZQUFZLE1BQW1CO0FBQzNCLFFBQUksQ0FBQyxLQUFNLFFBQU87QUFFbEIsUUFBSTtBQUNBLGFBQU8sS0FBSyxlQUFlLE1BQU0sWUFBWTtBQUFBLElBQ2pELFNBQVE7QUFDSixhQUFPO0FBQUEsSUFDWDtBQUFBLEVBQ0o7QUFDSjtBQU1PLElBQU0saUJBQU4sTUFBK0M7QUFBQSxFQUNsRCxlQUFlLE1BQVcsYUFBNkI7QUFDbkQsV0FBTyxPQUFPLE9BQU8sSUFBSSxJQUFJO0FBQUEsRUFDakM7QUFBQSxFQUVBLFlBQVksTUFBbUI7QUFDM0IsV0FBTyxPQUFPLE9BQU8sSUFBSSxJQUFJO0FBQUEsRUFDakM7QUFDSjtBQU1PLElBQU0sZUFBTixNQUE2QztBQUFBLEVBQ2hELGVBQWUsTUFBVyxhQUE2QjtBQUNuRCxRQUFJLENBQUMsS0FBTSxRQUFPO0FBQ2xCLFdBQU8sT0FBTyxTQUFTLFdBQVcsT0FBTyxLQUFLLFVBQVUsSUFBSTtBQUFBLEVBQ2hFO0FBQUEsRUFFQSxZQUFZLE1BQW1CO0FBQzNCLFFBQUksQ0FBQyxLQUFNLFFBQU87QUFDbEIsV0FBTyxPQUFPLFNBQVMsV0FBVyxPQUFPLEtBQUssVUFBVSxJQUFJO0FBQUEsRUFDaEU7QUFDSjs7O0F2QnJGQSxPQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFrRWxDLE9BQU8sT0FBTyxTQUFnQjtBQUs5QixPQUFPLE9BQU8seUJBQXlCLGVBQU8sdUJBQXVCLEtBQUssY0FBTTtBQUV6RSxPQUFPLHFCQUFxQjsiLAogICJuYW1lcyI6IFsiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJfYSIsICJFcnJvciIsICJkaXNwYXRjaFdhaWxzRXZlbnQiLCAiY2FsbCIsICJfYSIsICJyZXNpemFibGUiLCAiY2FsbCIsICJfYSIsICJjYWxsIiwgImNhbGwiLCAiSGlkZU1ldGhvZCIsICJTaG93TWV0aG9kIiwgImlzRG9jdW1lbnREb3RBbGwiLCAiX2EiLCAicmVhc29uIiwgInZhbHVlIiwgImNhbGwiLCAiZ2V0QW5kRGVsZXRlUmVzcG9uc2UiLCAiZ2VuZXJhdGVJRCIsICJjYWxsIiwgIkFycmF5IiwgIk1hcCIsICJBcnJheSIsICJNYXAiLCAia2V5IiwgImNhbGwiXQp9Cg==
