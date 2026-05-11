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
  IOS: () => ios_exports,
  Screens: () => screens_exports,
  System: () => system_exports,
  WML: () => wml_exports,
  Window: () => window_default,
  clientId: () => clientId,
  getTransport: () => getTransport,
  loadOptionalScript: () => loadOptionalScript,
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
var CHUNK_THRESHOLD = 512 * 1024;
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
  CancelCall: 10,
  IOS: 11
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
  let body = {
    object: objectID,
    method
  };
  if (args !== null && args !== void 0) {
    body.args = args;
  }
  let headers = {
    ["x-wails-client-id"]: clientId,
    ["Content-Type"]: "application/json"
  };
  if (windowName) {
    headers["x-wails-window-name"] = windowName;
  }
  const bodyStr = JSON.stringify(body);
  let response;
  if (bodyStr.length > CHUNK_THRESHOLD) {
    response = await sendChunked(url, headers, bodyStr);
  } else {
    response = await fetch(url, { method: "POST", headers, body: bodyStr });
  }
  if (!response.ok) {
    throw new Error(await response.text());
  }
  if (((_b = (_a2 = response.headers.get("Content-Type")) == null ? void 0 : _a2.indexOf("application/json")) != null ? _b : -1) !== -1) {
    return response.json();
  } else {
    return response.text();
  }
}
async function sendChunked(url, headers, bodyStr) {
  const chunkId = nanoid();
  const bodyBytes = new TextEncoder().encode(bodyStr);
  const totalChunks = Math.ceil(bodyBytes.length / CHUNK_THRESHOLD);
  for (let i = 0; i < totalChunks - 1; i++) {
    const chunk = bodyBytes.subarray(i * CHUNK_THRESHOLD, (i + 1) * CHUNK_THRESHOLD);
    const resp = await fetch(url, {
      method: "POST",
      headers: __spreadProps(__spreadValues({}, headers), {
        "x-wails-chunk-id": chunkId,
        "x-wails-chunk-index": String(i),
        "x-wails-chunk-total": String(totalChunks)
      }),
      body: chunk
    });
    if (!resp.ok) {
      throw new Error(await resp.text());
    }
  }
  return fetch(url, {
    method: "POST",
    headers: __spreadProps(__spreadValues({}, headers), {
      "x-wails-chunk-id": chunkId,
      "x-wails-chunk-index": String(totalChunks - 1),
      "x-wails-chunk-total": String(totalChunks)
    }),
    body: bodyBytes.subarray((totalChunks - 1) * CHUNK_THRESHOLD)
  });
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
var call2 = newRuntimeCaller(objectNames.Dialog);
var DialogInfo = 0;
var DialogWarning = 1;
var DialogError = 2;
var DialogQuestion = 3;
var DialogOpenFile = 4;
var DialogSaveFile = 5;
function dialog(type, options = {}) {
  return call2(type, options);
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

// desktop/@wailsio/runtime/src/create.ts
var create_exports = {};
__export(create_exports, {
  Any: () => Any,
  Array: () => Array2,
  ByteSlice: () => ByteSlice,
  Events: () => Events,
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
var Events = {};

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
    WindowLoadStarted: "linux:WindowLoadStarted",
    WindowLoadRedirected: "linux:WindowLoadRedirected",
    WindowLoadCommitted: "linux:WindowLoadCommitted",
    WindowLoadFinished: "linux:WindowLoadFinished"
  }),
  iOS: Object.freeze({
    ApplicationDidBecomeActive: "ios:ApplicationDidBecomeActive",
    ApplicationDidEnterBackground: "ios:ApplicationDidEnterBackground",
    ApplicationDidFinishLaunching: "ios:ApplicationDidFinishLaunching",
    ApplicationDidReceiveMemoryWarning: "ios:ApplicationDidReceiveMemoryWarning",
    ApplicationWillEnterForeground: "ios:ApplicationWillEnterForeground",
    ApplicationWillResignActive: "ios:ApplicationWillResignActive",
    ApplicationWillTerminate: "ios:ApplicationWillTerminate",
    WindowDidLoad: "ios:WindowDidLoad",
    WindowWillAppear: "ios:WindowWillAppear",
    WindowDidAppear: "ios:WindowDidAppear",
    WindowWillDisappear: "ios:WindowWillDisappear",
    WindowDidDisappear: "ios:WindowDidDisappear",
    WindowSafeAreaInsetsChanged: "ios:WindowSafeAreaInsetsChanged",
    WindowOrientationChanged: "ios:WindowOrientationChanged",
    WindowTouchBegan: "ios:WindowTouchBegan",
    WindowTouchMoved: "ios:WindowTouchMoved",
    WindowTouchEnded: "ios:WindowTouchEnded",
    WindowTouchCancelled: "ios:WindowTouchCancelled",
    WebViewDidStartNavigation: "ios:WebViewDidStartNavigation",
    WebViewDidFinishNavigation: "ios:WebViewDidFinishNavigation",
    WebViewDidFailNavigation: "ios:WebViewDidFailNavigation",
    WebViewDecidePolicyForNavigationAction: "ios:WebViewDecidePolicyForNavigationAction"
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
    WindowZoomReset: "common:WindowZoomReset"
  })
});

// desktop/@wailsio/runtime/src/events.ts
window._wails = window._wails || {};
window._wails.dispatchWailsEvent = dispatchWailsEvent;
var call3 = newRuntimeCaller(objectNames.Events);
var EmitMethod = 0;
var WailsEvent = class {
  constructor(name, data) {
    this.name = name;
    this.data = data != null ? data : null;
  }
};
function dispatchWailsEvent(event) {
  let listeners = eventListeners.get(event.name);
  if (!listeners) {
    return;
  }
  let wailsEvent = new WailsEvent(
    event.name,
    event.name in Events ? Events[event.name](event.data) : event.data
  );
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
  return call3(EmitMethod, new WailsEvent(name, data));
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
var DROP_TARGET_ATTRIBUTE = "data-file-drop-target";
var DROP_TARGET_ACTIVE_CLASS = "file-drop-target-active";
var currentDropTarget = null;
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
var FilesDropped = 50;
var PrintMethod = 51;
var SetScreenMethod = 52;
function getDropTargetElement(element) {
  if (!element) {
    return null;
  }
  return element.closest("[".concat(DROP_TARGET_ATTRIBUTE, "]"));
}
function canResolveFilePaths() {
  var _a2, _b, _c, _d;
  if (((_b = (_a2 = window.chrome) == null ? void 0 : _a2.webview) == null ? void 0 : _b.postMessageWithAdditionalObjects) == null) {
    return false;
  }
  return ((_d = (_c = window._wails) == null ? void 0 : _c.flags) == null ? void 0 : _d.enableFileDrop) === true;
}
function resolveFilePaths(x, y, files) {
  var _a2, _b;
  if ((_b = (_a2 = window.chrome) == null ? void 0 : _a2.webview) == null ? void 0 : _b.postMessageWithAdditionalObjects) {
    window.chrome.webview.postMessageWithAdditionalObjects("file:drop:".concat(x, ":").concat(y), files);
  }
}
var nativeDragActive = false;
function cleanupNativeDrag() {
  nativeDragActive = false;
  if (currentDropTarget) {
    currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
    currentDropTarget = null;
  }
}
function handleDragEnter() {
  var _a2, _b;
  if (((_b = (_a2 = window._wails) == null ? void 0 : _a2.flags) == null ? void 0 : _b.enableFileDrop) === false) {
    return;
  }
  nativeDragActive = true;
}
function handleDragLeave() {
  cleanupNativeDrag();
}
function handleDragOver(x, y) {
  var _a2, _b;
  if (!nativeDragActive) return;
  if (((_b = (_a2 = window._wails) == null ? void 0 : _a2.flags) == null ? void 0 : _b.enableFileDrop) === false) {
    return;
  }
  const targetElement = document.elementFromPoint(x, y);
  const dropTarget = getDropTargetElement(targetElement);
  if (currentDropTarget && currentDropTarget !== dropTarget) {
    currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
  }
  if (dropTarget) {
    dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
    currentDropTarget = dropTarget;
  } else {
    currentDropTarget = null;
  }
}
var callerSym = /* @__PURE__ */ Symbol("caller");
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
   * Handles file drops originating from platform-specific code (e.g., macOS/Linux native drag-and-drop).
   * Gathers information about the drop target element and sends it back to the Go backend.
   *
   * @param filenames - An array of file paths (strings) that were dropped.
   * @param x - The x-coordinate of the drop event (CSS pixels).
   * @param y - The y-coordinate of the drop event (CSS pixels).
   */
  HandlePlatformFileDrop(filenames, x, y) {
    var _a2, _b;
    if (((_b = (_a2 = window._wails) == null ? void 0 : _a2.flags) == null ? void 0 : _b.enableFileDrop) === false) {
      return;
    }
    const element = document.elementFromPoint(x, y);
    const dropTarget = getDropTargetElement(element);
    if (!dropTarget) {
      return;
    }
    const elementDetails = {
      id: dropTarget.id,
      classList: Array.from(dropTarget.classList),
      attributes: {}
    };
    for (let i = 0; i < dropTarget.attributes.length; i++) {
      const attr = dropTarget.attributes[i];
      elementDetails.attributes[attr.name] = attr.value;
    }
    const payload = {
      filenames,
      x,
      y,
      elementDetails
    };
    this[callerSym](FilesDropped, payload);
    cleanupNativeDrag();
  }
  /**
   * Moves the window to the center of the specified screen's work area.
   *
   * @param screenID - The ID of the target screen.
   */
  SetScreen(screenID) {
    return this[callerSym](SetScreenMethod, { screenID });
  }
  /* Triggers Windows 11 Snap Assist feature (Windows only).
   * This is equivalent to pressing Win+Z and shows snap layout options.
   */
  SnapAssist() {
    return this[callerSym](SnapAssistMethod);
  }
  /**
   * Opens the print dialog for the window.
   */
  Print() {
    return this[callerSym](PrintMethod);
  }
};
var Window = _Window;
var thisWindow = new Window("");
function setupDropTargetListeners() {
  const docElement = document.documentElement;
  let dragEnterCounter = 0;
  docElement.addEventListener("dragenter", (event) => {
    var _a2, _b, _c;
    if (!((_a2 = event.dataTransfer) == null ? void 0 : _a2.types.includes("Files"))) {
      return;
    }
    event.preventDefault();
    if (((_c = (_b = window._wails) == null ? void 0 : _b.flags) == null ? void 0 : _c.enableFileDrop) === false) {
      event.dataTransfer.dropEffect = "none";
      return;
    }
    dragEnterCounter++;
    const targetElement = document.elementFromPoint(event.clientX, event.clientY);
    const dropTarget = getDropTargetElement(targetElement);
    if (currentDropTarget && currentDropTarget !== dropTarget) {
      currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
    }
    if (dropTarget) {
      dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
      event.dataTransfer.dropEffect = "copy";
      currentDropTarget = dropTarget;
    } else {
      event.dataTransfer.dropEffect = "none";
      currentDropTarget = null;
    }
  }, false);
  docElement.addEventListener("dragover", (event) => {
    var _a2, _b, _c;
    if (!((_a2 = event.dataTransfer) == null ? void 0 : _a2.types.includes("Files"))) {
      return;
    }
    event.preventDefault();
    if (((_c = (_b = window._wails) == null ? void 0 : _b.flags) == null ? void 0 : _c.enableFileDrop) === false) {
      event.dataTransfer.dropEffect = "none";
      return;
    }
    const targetElement = document.elementFromPoint(event.clientX, event.clientY);
    const dropTarget = getDropTargetElement(targetElement);
    if (currentDropTarget && currentDropTarget !== dropTarget) {
      currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
    }
    if (dropTarget) {
      if (!dropTarget.classList.contains(DROP_TARGET_ACTIVE_CLASS)) {
        dropTarget.classList.add(DROP_TARGET_ACTIVE_CLASS);
      }
      event.dataTransfer.dropEffect = "copy";
      currentDropTarget = dropTarget;
    } else {
      event.dataTransfer.dropEffect = "none";
      currentDropTarget = null;
    }
  }, false);
  docElement.addEventListener("dragleave", (event) => {
    var _a2, _b, _c;
    if (!((_a2 = event.dataTransfer) == null ? void 0 : _a2.types.includes("Files"))) {
      return;
    }
    event.preventDefault();
    if (((_c = (_b = window._wails) == null ? void 0 : _b.flags) == null ? void 0 : _c.enableFileDrop) === false) {
      return;
    }
    if (event.relatedTarget === null) {
      return;
    }
    dragEnterCounter--;
    if (dragEnterCounter === 0 || currentDropTarget && !currentDropTarget.contains(event.relatedTarget)) {
      if (currentDropTarget) {
        currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
        currentDropTarget = null;
      }
      dragEnterCounter = 0;
    }
  }, false);
  docElement.addEventListener("drop", (event) => {
    var _a2, _b, _c;
    if (!((_a2 = event.dataTransfer) == null ? void 0 : _a2.types.includes("Files"))) {
      return;
    }
    event.preventDefault();
    if (((_c = (_b = window._wails) == null ? void 0 : _b.flags) == null ? void 0 : _c.enableFileDrop) === false) {
      return;
    }
    dragEnterCounter = 0;
    if (currentDropTarget) {
      currentDropTarget.classList.remove(DROP_TARGET_ACTIVE_CLASS);
      currentDropTarget = null;
    }
    if (canResolveFilePaths()) {
      const files = [];
      if (event.dataTransfer.items) {
        for (const item of event.dataTransfer.items) {
          if (item.kind === "file") {
            const file = item.getAsFile();
            if (file) files.push(file);
          }
        }
      } else if (event.dataTransfer.files) {
        for (const file of event.dataTransfer.files) {
          files.push(file);
        }
      }
      if (files.length > 0) {
        resolveFilePaths(event.clientX, event.clientY, files);
      }
    }
  }, false);
}
if (typeof window !== "undefined" && typeof document !== "undefined") {
  setupDropTargetListeners();
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
var controllerSym = /* @__PURE__ */ Symbol("controller");
var triggerMapSym = /* @__PURE__ */ Symbol("triggerMap");
var elementCountSym = /* @__PURE__ */ Symbol("elementCount");
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
var SystemCapabilities = 2;
var _invoke = (function() {
  var _a2, _b, _c, _d, _e, _f;
  try {
    if ((_b = (_a2 = window.chrome) == null ? void 0 : _a2.webview) == null ? void 0 : _b.postMessage) {
      return window.chrome.webview.postMessage.bind(window.chrome.webview);
    } else if ((_e = (_d = (_c = window.webkit) == null ? void 0 : _c.messageHandlers) == null ? void 0 : _d["external"]) == null ? void 0 : _e.postMessage) {
      return window.webkit.messageHandlers["external"].postMessage.bind(window.webkit.messageHandlers["external"]);
    } else if ((_f = window.wails) == null ? void 0 : _f.invoke) {
      return (msg) => window.wails.invoke(typeof msg === "string" ? msg : JSON.stringify(msg));
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
  return call4(SystemCapabilities);
}
function Environment() {
  return call4(SystemEnvironment);
}
function IsWindows() {
  var _a2, _b;
  return ((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.OS) === "windows";
}
function IsLinux() {
  var _a2, _b;
  return ((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.OS) === "linux";
}
function IsMac() {
  var _a2, _b;
  return ((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.OS) === "darwin";
}
function IsAMD64() {
  var _a2, _b;
  return ((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.Arch) === "amd64";
}
function IsARM() {
  var _a2, _b;
  return ((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.Arch) === "arm";
}
function IsARM64() {
  var _a2, _b;
  return ((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.Arch) === "arm64";
}
function IsDebug() {
  var _a2, _b;
  return Boolean((_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.Debug);
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
var dragInitDone = false;
function isMobile() {
  var _a2, _b;
  const os = (_b = (_a2 = window._wails) == null ? void 0 : _a2.environment) == null ? void 0 : _b.OS;
  if (os === "ios" || os === "android") return true;
  const ua = navigator.userAgent || navigator.vendor || window.opera || "";
  return /android|iphone|ipad|ipod|iemobile|wpdesktop/i.test(ua);
}
function tryInitDragHandlers() {
  if (dragInitDone) return;
  if (isMobile()) return;
  window.addEventListener("mousedown", update, { capture: true });
  window.addEventListener("mousemove", update, { capture: true });
  window.addEventListener("mouseup", update, { capture: true });
  for (const ev of ["click", "contextmenu", "dblclick"]) {
    window.addEventListener(ev, suppressEvent, { capture: true });
  }
  dragInitDone = true;
}
tryInitDragHandlers();
document.addEventListener("DOMContentLoaded", tryInitDragHandlers, { once: true });
var dragEnvPolls = 0;
var dragEnvPoll = window.setInterval(() => {
  if (dragInitDone) {
    window.clearInterval(dragEnvPoll);
    return;
  }
  tryInitDragHandlers();
  if (++dragEnvPolls > 100) {
    window.clearInterval(dragEnvPoll);
  }
}, 50);
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
var barrierSym = /* @__PURE__ */ Symbol("barrier");
var cancelImplSym = /* @__PURE__ */ Symbol("cancelImpl");
var _a;
var species = (_a = Symbol.species) != null ? _a : /* @__PURE__ */ Symbol("speciesPolyfill");
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
function generateID() {
  let result;
  do {
    result = nanoid();
  } while (callResponses.has(result));
  return result;
}
function Call(options) {
  const id = generateID();
  const result = CancellablePromise.withResolvers();
  callResponses.set(id, { resolve: result.resolve, reject: result.reject });
  const request = call7(CallBinding, Object.assign({ "call-id": id }, options));
  let running = true;
  request.then((res) => {
    running = false;
    callResponses.delete(id);
    result.resolve(res);
  }, (err) => {
    running = false;
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

// desktop/@wailsio/runtime/src/screens.ts
var screens_exports = {};
__export(screens_exports, {
  GetAll: () => GetAll,
  GetByID: () => GetByID,
  GetByIndex: () => GetByIndex,
  GetCurrent: () => GetCurrent,
  GetPrimary: () => GetPrimary
});
var call9 = newRuntimeCaller(objectNames.Screens);
var getAll = 0;
var getPrimary = 1;
var getCurrent = 2;
var getByID = 3;
var getByIndex = 4;
function GetAll() {
  return call9(getAll);
}
function GetPrimary() {
  return call9(getPrimary);
}
function GetCurrent() {
  return call9(getCurrent);
}
function GetByID(id) {
  return call9(getByID, { id });
}
function GetByIndex(index) {
  return call9(getByIndex, { index });
}

// desktop/@wailsio/runtime/src/ios.ts
var ios_exports = {};
__export(ios_exports, {
  Device: () => Device,
  Haptics: () => Haptics
});
var call10 = newRuntimeCaller(objectNames.IOS);
var HapticsImpact = 0;
var DeviceInfo = 1;
var Haptics;
((Haptics2) => {
  function Impact(style = "medium") {
    return call10(HapticsImpact, { style });
  }
  Haptics2.Impact = Impact;
})(Haptics || (Haptics = {}));
var Device;
((Device2) => {
  function Info2() {
    return call10(DeviceInfo);
  }
  Device2.Info = Info2;
})(Device || (Device = {}));

// desktop/@wailsio/runtime/src/index.ts
window._wails = window._wails || {};
window._wails.invoke = invoke;
window._wails.clientId = clientId;
window._wails.handlePlatformFileDrop = window_default.HandlePlatformFileDrop.bind(window_default);
window._wails.handleDragEnter = handleDragEnter;
window._wails.handleDragLeave = handleDragLeave;
window._wails.handleDragOver = handleDragOver;
invoke("wails:runtime:ready");
function loadOptionalScript(url) {
  return fetch(url, { method: "HEAD" }).then((response) => {
    if (response.ok) {
      const contentType = (response.headers.get("content-type") || "").toLowerCase();
      if (contentType.includes("javascript")) {
        const script = document.createElement("script");
        script.src = url;
        document.head.appendChild(script);
      }
    }
  }).catch(() => {
  });
}
loadOptionalScript("/wails/custom.js");
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
  ios_exports as IOS,
  screens_exports as Screens,
  system_exports as System,
  wml_exports as WML,
  window_default as Window,
  clientId,
  getTransport,
  loadOptionalScript,
  objectNames,
  setTransport
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY3JlYXRlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudF90eXBlcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbGFibGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbmNlbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2lvcy50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLy8gU2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG5pbXBvcnQgXCIuL2NvbnRleHRtZW51LmpzXCI7XHJcbmltcG9ydCBcIi4vZHJhZy5qc1wiO1xyXG5cclxuLy8gUmUtZXhwb3J0IHB1YmxpYyBBUElcclxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSBcIi4vYXBwbGljYXRpb24uanNcIjtcclxuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyLmpzXCI7XHJcbmltcG9ydCAqIGFzIENhbGwgZnJvbSBcIi4vY2FsbHMuanNcIjtcclxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZC5qc1wiO1xyXG5pbXBvcnQgKiBhcyBDcmVhdGUgZnJvbSBcIi4vY3JlYXRlLmpzXCI7XHJcbmltcG9ydCAqIGFzIERpYWxvZ3MgZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBFdmVudHMgZnJvbSBcIi4vZXZlbnRzLmpzXCI7XHJcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuL2ZsYWdzLmpzXCI7XHJcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XHJcbmltcG9ydCAqIGFzIElPUyBmcm9tIFwiLi9pb3MuanNcIjtcclxuaW1wb3J0IFdpbmRvdywgeyBoYW5kbGVEcmFnRW50ZXIsIGhhbmRsZURyYWdMZWF2ZSwgaGFuZGxlRHJhZ092ZXIgfSBmcm9tIFwiLi93aW5kb3cuanNcIjtcclxuaW1wb3J0ICogYXMgV01MIGZyb20gXCIuL3dtbC5qc1wiO1xyXG5cclxuZXhwb3J0IHtcclxuICAgIEFwcGxpY2F0aW9uLFxyXG4gICAgQnJvd3NlcixcclxuICAgIENhbGwsXHJcbiAgICBDbGlwYm9hcmQsXHJcbiAgICBEaWFsb2dzLFxyXG4gICAgRXZlbnRzLFxyXG4gICAgRmxhZ3MsXHJcbiAgICBTY3JlZW5zLFxyXG4gICAgU3lzdGVtLFxyXG4gICAgSU9TLFxyXG4gICAgV2luZG93LFxyXG4gICAgV01MXHJcbn07XHJcblxyXG4vKipcclxuICogQW4gaW50ZXJuYWwgdXRpbGl0eSBjb25zdW1lZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXHJcbiAqXHJcbiAqIEBpZ25vcmVcclxuICovXHJcbmV4cG9ydCB7IENyZWF0ZSB9O1xyXG5cclxuZXhwb3J0ICogZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcclxuXHJcbi8vIEV4cG9ydCB0cmFuc3BvcnQgaW50ZXJmYWNlcyBhbmQgdXRpbGl0aWVzXHJcbmV4cG9ydCB7XHJcbiAgICBzZXRUcmFuc3BvcnQsXHJcbiAgICBnZXRUcmFuc3BvcnQsXHJcbiAgICB0eXBlIFJ1bnRpbWVUcmFuc3BvcnQsXHJcbiAgICBvYmplY3ROYW1lcyxcclxuICAgIGNsaWVudElkLFxyXG59IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmltcG9ydCB7IGNsaWVudElkIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuLy8gTm90aWZ5IGJhY2tlbmRcclxud2luZG93Ll93YWlscy5pbnZva2UgPSBTeXN0ZW0uaW52b2tlO1xyXG53aW5kb3cuX3dhaWxzLmNsaWVudElkID0gY2xpZW50SWQ7XHJcblxyXG4vLyBSZWdpc3RlciBwbGF0Zm9ybSBoYW5kbGVycyAoaW50ZXJuYWwgQVBJKVxyXG4vLyBOb3RlOiBXaW5kb3cgaXMgdGhlIHRoaXNXaW5kb3cgaW5zdGFuY2UgKGRlZmF1bHQgZXhwb3J0IGZyb20gd2luZG93LnRzKVxyXG4vLyBCaW5kaW5nIGVuc3VyZXMgJ3RoaXMnIGNvcnJlY3RseSByZWZlcnMgdG8gdGhlIGN1cnJlbnQgd2luZG93IGluc3RhbmNlXHJcbndpbmRvdy5fd2FpbHMuaGFuZGxlUGxhdGZvcm1GaWxlRHJvcCA9IFdpbmRvdy5IYW5kbGVQbGF0Zm9ybUZpbGVEcm9wLmJpbmQoV2luZG93KTtcclxuXHJcbi8vIExpbnV4LXNwZWNpZmljIGRyYWcgaGFuZGxlcnMgKEdUSyBpbnRlcmNlcHRzIERPTSBkcmFnIGV2ZW50cylcclxud2luZG93Ll93YWlscy5oYW5kbGVEcmFnRW50ZXIgPSBoYW5kbGVEcmFnRW50ZXI7XHJcbndpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ0xlYXZlID0gaGFuZGxlRHJhZ0xlYXZlO1xyXG53aW5kb3cuX3dhaWxzLmhhbmRsZURyYWdPdmVyID0gaGFuZGxlRHJhZ092ZXI7XHJcblxyXG5TeXN0ZW0uaW52b2tlKFwid2FpbHM6cnVudGltZTpyZWFkeVwiKTtcclxuXHJcbi8qKlxyXG4gKiBMb2FkcyBhIHNjcmlwdCBmcm9tIHRoZSBnaXZlbiBVUkwgaWYgaXQgZXhpc3RzLlxyXG4gKiBVc2VzIEhFQUQgcmVxdWVzdCB0byBjaGVjayBleGlzdGVuY2UsIHRoZW4gaW5qZWN0cyBhIHNjcmlwdCB0YWcuXHJcbiAqIFNpbGVudGx5IGlnbm9yZXMgaWYgdGhlIHNjcmlwdCBkb2Vzbid0IGV4aXN0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIGxvYWRPcHRpb25hbFNjcmlwdCh1cmw6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGZldGNoKHVybCwgeyBtZXRob2Q6ICdIRUFEJyB9KVxyXG4gICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcclxuICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAvLyBWZXJpZnkgdGhlIHJlc3BvbnNlIGlzIGFjdHVhbGx5IEphdmFTY3JpcHQgYW5kIG5vdCBhbiBIVE1MIGZhbGxiYWNrXHJcbiAgICAgICAgICAgICAgICAvLyAoZS5nLiBWaXRlIGRldiBzZXJ2ZXIgcmV0dXJucyBpbmRleC5odG1sIGZvciB1bmtub3duIHJvdXRlcylcclxuICAgICAgICAgICAgICAgIGNvbnN0IGNvbnRlbnRUeXBlID0gKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KCdjb250ZW50LXR5cGUnKSB8fCAnJykudG9Mb3dlckNhc2UoKTtcclxuICAgICAgICAgICAgICAgIGlmIChjb250ZW50VHlwZS5pbmNsdWRlcygnamF2YXNjcmlwdCcpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgY29uc3Qgc2NyaXB0ID0gZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgnc2NyaXB0Jyk7XHJcbiAgICAgICAgICAgICAgICAgICAgc2NyaXB0LnNyYyA9IHVybDtcclxuICAgICAgICAgICAgICAgICAgICBkb2N1bWVudC5oZWFkLmFwcGVuZENoaWxkKHNjcmlwdCk7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KVxyXG4gICAgICAgIC5jYXRjaCgoKSA9PiB7fSk7IC8vIFNpbGVudGx5IGlnbm9yZSAtIHNjcmlwdCBpcyBvcHRpb25hbFxyXG59XHJcblxyXG4vLyBMb2FkIGN1c3RvbS5qcyBpZiBhdmFpbGFibGUgKHVzZWQgYnkgc2VydmVyIG1vZGUgZm9yIFdlYlNvY2tldCBldmVudHMsIGV0Yy4pXHJcbmxvYWRPcHRpb25hbFNjcmlwdCgnL3dhaWxzL2N1c3RvbS5qcycpO1xyXG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcclxuaW1wb3J0IHsgUXVlc3Rpb24gfSBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XHJcbmltcG9ydCB7IEVtaXQgfSBmcm9tIFwiLi9ldmVudHMuanNcIjtcclxuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XHJcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93LmpzXCI7XHJcblxyXG4vKipcclxuICogU2VuZHMgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgb3B0aW9uYWwgZGF0YS5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXHJcbiAqIEBwYXJhbSBbZGF0YT1udWxsXSAtIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XHJcbiAgICBFbWl0KGV2ZW50TmFtZSwgZGF0YSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyBhIG1ldGhvZCBvbiBhIHNwZWNpZmllZCB3aW5kb3cuXHJcbiAqXHJcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXHJcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lOiBzdHJpbmcsIG1ldGhvZE5hbWU6IHN0cmluZykge1xyXG4gICAgY29uc3QgdGFyZ2V0V2luZG93ID0gV2luZG93LkdldCh3aW5kb3dOYW1lKTtcclxuICAgIGNvbnN0IG1ldGhvZCA9ICh0YXJnZXRXaW5kb3cgYXMgYW55KVttZXRob2ROYW1lXTtcclxuXHJcbiAgICBpZiAodHlwZW9mIG1ldGhvZCAhPT0gXCJmdW5jdGlvblwiKSB7XHJcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgbWV0aG9kLmNhbGwodGFyZ2V0V2luZG93KTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXQuXHJcbiAqL1xyXG5mdW5jdGlvbiBvbldNTFRyaWdnZXJlZChldjogRXZlbnQpOiB2b2lkIHtcclxuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0IGFzIEVsZW1lbnQ7XHJcblxyXG4gICAgZnVuY3Rpb24gcnVuRWZmZWN0KGNob2ljZSA9IFwiWWVzXCIpIHtcclxuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcclxuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XHJcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XHJcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW51cmwnKTtcclxuXHJcbiAgICAgICAgaWYgKGV2ZW50VHlwZSAhPT0gbnVsbClcclxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XHJcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcclxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh0YXJnZXRXaW5kb3csIHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcclxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XHJcblxyXG4gICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICBRdWVzdGlvbih7XHJcbiAgICAgICAgICAgIFRpdGxlOiBcIkNvbmZpcm1cIixcclxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcclxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxyXG4gICAgICAgICAgICBCdXR0b25zOiBbXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXHJcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XHJcbiAgICAgICAgICAgIF1cclxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJ1bkVmZmVjdCgpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBjb250cm9sbGVyU3ltID0gU3ltYm9sKFwiY29udHJvbGxlclwiKTtcclxuY29uc3QgdHJpZ2dlck1hcFN5bSA9IFN5bWJvbChcInRyaWdnZXJNYXBcIik7XHJcbmNvbnN0IGVsZW1lbnRDb3VudFN5bSA9IFN5bWJvbChcImVsZW1lbnRDb3VudFwiKTtcclxuXHJcbi8qKlxyXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXHJcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXHJcbiAqL1xyXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XHJcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cclxuICAgIFtjb250cm9sbGVyU3ltXTogQWJvcnRDb250cm9sbGVyO1xyXG5cclxuICAgIGNvbnN0cnVjdG9yKCkge1xyXG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0gPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIGFuIG9wdGlvbnMgb2JqZWN0IGZvciBhZGRFdmVudExpc3RlbmVyIHRoYXQgdGllcyB0aGUgbGlzdGVuZXJcclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxyXG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xyXG4gICAgICovXHJcbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xyXG4gICAgICAgIHJldHVybiB7IHNpZ25hbDogdGhpc1tjb250cm9sbGVyU3ltXS5zaWduYWwgfTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzIGFuZCByZXNldHMgdGhlIHJlZ2lzdHJ5LlxyXG4gICAgICovXHJcbiAgICByZXNldCgpOiB2b2lkIHtcclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dLmFib3J0KCk7XHJcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFdlYWtNYXBSZWdpc3RyeSBtYXBzIGFjdGl2ZSB0cmlnZ2VyIGV2ZW50cyB0byBlYWNoIERPTSBlbGVtZW50IHRocm91Z2ggYSBXZWFrTWFwLlxyXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcclxuICogY29sbGVjdGlvbiBvZiB0aGUgaW52b2x2ZWQgZWxlbWVudHMuXHJcbiAqL1xyXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xyXG4gICAgLyoqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy4gKi9cclxuICAgIFt0cmlnZ2VyTWFwU3ltXTogV2Vha01hcDxFbGVtZW50LCBzdHJpbmdbXT47XHJcbiAgICAvKiogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLiAqL1xyXG4gICAgW2VsZW1lbnRDb3VudFN5bV06IG51bWJlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcigpIHtcclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxyXG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XHJcbiAgICAgKi9cclxuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XHJcbiAgICAgICAgaWYgKCF0aGlzW3RyaWdnZXJNYXBTeW1dLmhhcyhlbGVtZW50KSkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0rKzsgfVxyXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcclxuICAgICAgICByZXR1cm4ge307XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cclxuICAgICAqL1xyXG4gICAgcmVzZXQoKTogdm9pZCB7XHJcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxyXG4gICAgICAgICAgICByZXR1cm47XHJcblxyXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xyXG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXHJcbiAgICAgICAgICAgICAgICBicmVhaztcclxuXHJcbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwU3ltXS5nZXQoZWxlbWVudCk7XHJcbiAgICAgICAgICAgIGlmICh0cmlnZ2VycyAhPSBudWxsKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXS0tOyB9XHJcblxyXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXHJcbiAgICAgICAgICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XHJcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcclxuICAgIH1cclxufVxyXG5cclxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xyXG5cclxuLyoqXHJcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cclxuICovXHJcbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50OiBFbGVtZW50KTogdm9pZCB7XHJcbiAgICBjb25zdCB0cmlnZ2VyUmVnRXhwID0gL1xcUysvZztcclxuICAgIGNvbnN0IHRyaWdnZXJBdHRyID0gKGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiKTtcclxuICAgIGNvbnN0IHRyaWdnZXJzOiBzdHJpbmdbXSA9IFtdO1xyXG5cclxuICAgIGxldCBtYXRjaDtcclxuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcclxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcclxuXHJcbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XHJcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVuYWJsZSgpOiB2b2lkIHtcclxuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpOiB2b2lkIHtcclxuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xyXG4gICAgZG9jdW1lbnQuYm9keS5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLWV2ZW50XSwgW3dtbC13aW5kb3ddLCBbd21sLW9wZW51cmxdLCBbZGF0YS13bWwtZXZlbnRdLCBbZGF0YS13bWwtd2luZG93XSwgW2RhdGEtd21sLW9wZW51cmxdJykuZm9yRWFjaChhZGRXTUxMaXN0ZW5lcnMpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xyXG5cclxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xyXG5cclxuLyoqXHJcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMLlxyXG4gKlxyXG4gKiBAcGFyYW0gdXJsIC0gVGhlIFVSTCB0byBvcGVuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmw6IHN0cmluZyB8IFVSTCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmw6IHVybC50b1N0cmluZygpfSk7XHJcbn1cclxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxyXG5cclxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXHJcbi8vXHJcbi8vIENvcHlyaWdodCAyMDE3IEFuZHJleSBTaXRuaWsgPGFuZHJleUBzaXRuaWsucnU+XHJcbi8vXHJcbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2ZcclxuLy8gdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbCBpblxyXG4vLyB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvXHJcbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXHJcbi8vIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzIGZ1cm5pc2hlZCB0byBkbyBzbyxcclxuLy8gICAgIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxyXG4vL1xyXG4vLyAgICAgVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXHJcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXHJcbi8vXHJcbi8vICAgICBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXHJcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTXHJcbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxyXG4vLyBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUiBMSUFCSUxJVFksIFdIRVRIRVJcclxuLy8gSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU5cclxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cclxuXHJcbi8vIFRoaXMgYWxwaGFiZXQgdXNlcyBgQS1aYS16MC05Xy1gIHN5bWJvbHMuXHJcbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cclxuLy8gUmVmZXJlbmNlcyB0byB0aGUgc2FtZSBmaWxlICh3b3JrcyBib3RoIGZvciBnemlwIGFuZCBicm90bGkpOlxyXG4vLyBgJ3VzZWAsIGBhbmRvbWAsIGFuZCBgcmljdCdgXHJcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XHJcbi8vIGAtMjZUYCwgYDE5ODNgLCBgNDBweGAsIGA3NXB4YCwgYGJ1c2hgLCBgamFja2AsIGBtaW5kYCwgYHZlcnlgLCBhbmQgYHdvbGZgXHJcbmNvbnN0IHVybEFscGhhYmV0ID1cclxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5hbm9pZChzaXplOiBudW1iZXIgPSAyMSk6IHN0cmluZyB7XHJcbiAgICBsZXQgaWQgPSAnJ1xyXG4gICAgLy8gQSBjb21wYWN0IGFsdGVybmF0aXZlIGZvciBgZm9yICh2YXIgaSA9IDA7IGkgPCBzdGVwOyBpKyspYC5cclxuICAgIGxldCBpID0gc2l6ZSB8IDBcclxuICAgIHdoaWxlIChpLS0pIHtcclxuICAgICAgICAvLyBgfCAwYCBpcyBtb3JlIGNvbXBhY3QgYW5kIGZhc3RlciB0aGFuIGBNYXRoLmZsb29yKClgLlxyXG4gICAgICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cclxuICAgIH1cclxuICAgIHJldHVybiBpZFxyXG59XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcclxuXHJcbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xyXG5cclxuLy8gU3RheSB1bmRlciBXZWJWaWV3MidzIH4yTUIgcmVxdWVzdCBib2R5IGJ1ZmZlcmluZyBsaW1pdCBpbiBXZWJSZXNvdXJjZVJlcXVlc3RlZC5cclxuY29uc3QgQ0hVTktfVEhSRVNIT0xEID0gNTEyICogMTAyNDtcclxuXHJcbi8vIFJlLWV4cG9ydCBuYW5vaWQgZm9yIGN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zXHJcbmV4cG9ydCB7IG5hbm9pZCB9O1xyXG5cclxuLy8gT2JqZWN0IE5hbWVzXHJcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xyXG4gICAgQ2FsbDogMCxcclxuICAgIENsaXBib2FyZDogMSxcclxuICAgIEFwcGxpY2F0aW9uOiAyLFxyXG4gICAgRXZlbnRzOiAzLFxyXG4gICAgQ29udGV4dE1lbnU6IDQsXHJcbiAgICBEaWFsb2c6IDUsXHJcbiAgICBXaW5kb3c6IDYsXHJcbiAgICBTY3JlZW5zOiA3LFxyXG4gICAgU3lzdGVtOiA4LFxyXG4gICAgQnJvd3NlcjogOSxcclxuICAgIENhbmNlbENhbGw6IDEwLFxyXG4gICAgSU9TOiAxMSxcclxufSk7XHJcbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcclxuXHJcbi8qKlxyXG4gKiBSdW50aW1lVHJhbnNwb3J0IGRlZmluZXMgdGhlIGludGVyZmFjZSBmb3IgY3VzdG9tIElQQyB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zLlxyXG4gKiBJbXBsZW1lbnQgdGhpcyBpbnRlcmZhY2UgdG8gdXNlIFdlYlNvY2tldHMsIGN1c3RvbSBwcm90b2NvbHMsIG9yIGFueSBvdGhlclxyXG4gKiB0cmFuc3BvcnQgbWVjaGFuaXNtIGluc3RlYWQgb2YgdGhlIGRlZmF1bHQgSFRUUCBmZXRjaC5cclxuICovXHJcbmV4cG9ydCBpbnRlcmZhY2UgUnVudGltZVRyYW5zcG9ydCB7XHJcbiAgICAvKipcclxuICAgICAqIFNlbmQgYSBydW50aW1lIGNhbGwgYW5kIHJldHVybiB0aGUgcmVzcG9uc2UuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIG9iamVjdElEIC0gVGhlIFdhaWxzIG9iamVjdCBJRCAoMD1DYWxsLCAxPUNsaXBib2FyZCwgZXRjLilcclxuICAgICAqIEBwYXJhbSBtZXRob2QgLSBUaGUgbWV0aG9kIElEIHRvIGNhbGxcclxuICAgICAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gT3B0aW9uYWwgd2luZG93IG5hbWVcclxuICAgICAqIEBwYXJhbSBhcmdzIC0gQXJndW1lbnRzIHRvIHBhc3MgKHdpbGwgYmUgSlNPTiBzdHJpbmdpZmllZCBpZiBwcmVzZW50KVxyXG4gICAgICogQHJldHVybnMgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHJlc3BvbnNlIGRhdGFcclxuICAgICAqL1xyXG4gICAgY2FsbChvYmplY3RJRDogbnVtYmVyLCBtZXRob2Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nLCBhcmdzOiBhbnkpOiBQcm9taXNlPGFueT47XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uIChjYW4gYmUgc2V0IGJ5IHVzZXIpXHJcbiAqL1xyXG5sZXQgY3VzdG9tVHJhbnNwb3J0OiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCA9IG51bGw7XHJcblxyXG4vKipcclxuICogU2V0IGEgY3VzdG9tIHRyYW5zcG9ydCBmb3IgYWxsIFdhaWxzIHJ1bnRpbWUgY2FsbHMuXHJcbiAqIFRoaXMgYWxsb3dzIHlvdSB0byByZXBsYWNlIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0IHdpdGhcclxuICogV2ViU29ja2V0cywgY3VzdG9tIHByb3RvY29scywgb3IgYW55IG90aGVyIG1lY2hhbmlzbS5cclxuICpcclxuICogQHBhcmFtIHRyYW5zcG9ydCAtIFlvdXIgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvblxyXG4gKlxyXG4gKiBAZXhhbXBsZVxyXG4gKiBgYGB0eXBlc2NyaXB0XHJcbiAqIGltcG9ydCB7IHNldFRyYW5zcG9ydCB9IGZyb20gJy93YWlscy9ydW50aW1lLmpzJztcclxuICpcclxuICogY29uc3Qgd3NUcmFuc3BvcnQgPSB7XHJcbiAqICAgY2FsbDogYXN5bmMgKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpID0+IHtcclxuICogICAgIC8vIFlvdXIgV2ViU29ja2V0IGltcGxlbWVudGF0aW9uXHJcbiAqICAgfVxyXG4gKiB9O1xyXG4gKlxyXG4gKiBzZXRUcmFuc3BvcnQod3NUcmFuc3BvcnQpO1xyXG4gKiBgYGBcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBzZXRUcmFuc3BvcnQodHJhbnNwb3J0OiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCk6IHZvaWQge1xyXG4gICAgY3VzdG9tVHJhbnNwb3J0ID0gdHJhbnNwb3J0O1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBjdXJyZW50IHRyYW5zcG9ydCAodXNlZnVsIGZvciBleHRlbmRpbmcvd3JhcHBpbmcpXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZ2V0VHJhbnNwb3J0KCk6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsIHtcclxuICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxyXG4gKlxyXG4gKiBAcGFyYW0gb2JqZWN0IC0gVGhlIG9iamVjdCB0byBpbnZva2UgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcmV0dXJuIFRoZSBuZXcgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcclxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kOiBudW1iZXIsIGFyZ3M6IGFueSA9IG51bGwpIHtcclxuICAgICAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0LCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xyXG4gICAgfTtcclxufVxyXG5cclxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcclxuICAgIC8vIFVzZSBjdXN0b20gdHJhbnNwb3J0IGlmIGF2YWlsYWJsZVxyXG4gICAgaWYgKGN1c3RvbVRyYW5zcG9ydCkge1xyXG4gICAgICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQuY2FsbChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBEZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcclxuXHJcbiAgICBsZXQgYm9keTogeyBvYmplY3Q6IG51bWJlcjsgbWV0aG9kOiBudW1iZXIsIGFyZ3M/OiBhbnkgfSA9IHtcclxuICAgICAgb2JqZWN0OiBvYmplY3RJRCxcclxuICAgICAgbWV0aG9kXHJcbiAgICB9XHJcbiAgICBpZiAoYXJncyAhPT0gbnVsbCAmJiBhcmdzICE9PSB1bmRlZmluZWQpIHtcclxuICAgICAgYm9keS5hcmdzID0gYXJncztcclxuICAgIH1cclxuXHJcbiAgICBsZXQgaGVhZGVyczogUmVjb3JkPHN0cmluZywgc3RyaW5nPiA9IHtcclxuICAgICAgICBbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXTogY2xpZW50SWQsXHJcbiAgICAgICAgW1wiQ29udGVudC1UeXBlXCJdOiBcImFwcGxpY2F0aW9uL2pzb25cIlxyXG4gICAgfVxyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBoZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgYm9keVN0ciA9IEpTT04uc3RyaW5naWZ5KGJvZHkpO1xyXG4gICAgbGV0IHJlc3BvbnNlOiBSZXNwb25zZTtcclxuICAgIGlmIChib2R5U3RyLmxlbmd0aCA+IENIVU5LX1RIUkVTSE9MRCkge1xyXG4gICAgICAgIHJlc3BvbnNlID0gYXdhaXQgc2VuZENodW5rZWQodXJsLCBoZWFkZXJzLCBib2R5U3RyKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgcmVzcG9uc2UgPSBhd2FpdCBmZXRjaCh1cmwsIHsgbWV0aG9kOiAnUE9TVCcsIGhlYWRlcnMsIGJvZHk6IGJvZHlTdHIgfSk7XHJcbiAgICB9XHJcbiAgICBpZiAoIXJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIHNlbmRDaHVua2VkIHNwbGl0cyBhIGxhcmdlIHNlcmlhbGlzZWQgcmVxdWVzdCBib2R5IGludG8gQ0hVTktfVEhSRVNIT0xELXNpemVkXHJcbi8vIGJ5dGUgY2h1bmtzIGFuZCBzZW5kcyB0aGVtIHNlcmlhbGx5LiAgRW5jb2RpbmcgdG8gVVRGLTggYnl0ZXMgYmVmb3JlIHNsaWNpbmdcclxuLy8gcHJldmVudHMgY29ycnVwdGlvbiBvZiBub24tQk1QIGNoYXJhY3RlcnMgKHN1cnJvZ2F0ZSBwYWlycykgdGhhdCB3b3VsZCBvY2N1clxyXG4vLyB3aGVuIHNwbGl0dGluZyBhdCBKYXZhU2NyaXB0IHN0cmluZyBpbmRpY2VzLiAgVGhlIEdvIHRyYW5zcG9ydCBhc3NlbWJsZXMgdGhlXHJcbi8vIHJhdyBieXRlcyBiZWZvcmUgcHJvY2Vzc2luZy4gIE9ubHkgdGhlIGZpbmFsIGNodW5rJ3MgcmVzcG9uc2UgY2FycmllcyB0aGUgUlBDIHJlc3VsdC5cclxuYXN5bmMgZnVuY3Rpb24gc2VuZENodW5rZWQodXJsOiBVUkwsIGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4sIGJvZHlTdHI6IHN0cmluZyk6IFByb21pc2U8UmVzcG9uc2U+IHtcclxuICAgIGNvbnN0IGNodW5rSWQgPSBuYW5vaWQoKTtcclxuICAgIGNvbnN0IGJvZHlCeXRlcyA9IG5ldyBUZXh0RW5jb2RlcigpLmVuY29kZShib2R5U3RyKTtcclxuICAgIGNvbnN0IHRvdGFsQ2h1bmtzID0gTWF0aC5jZWlsKGJvZHlCeXRlcy5sZW5ndGggLyBDSFVOS19USFJFU0hPTEQpO1xyXG5cclxuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgdG90YWxDaHVua3MgLSAxOyBpKyspIHtcclxuICAgICAgICBjb25zdCBjaHVuayA9IGJvZHlCeXRlcy5zdWJhcnJheShpICogQ0hVTktfVEhSRVNIT0xELCAoaSArIDEpICogQ0hVTktfVEhSRVNIT0xEKTtcclxuICAgICAgICBjb25zdCByZXNwID0gYXdhaXQgZmV0Y2godXJsLCB7XHJcbiAgICAgICAgICAgIG1ldGhvZDogJ1BPU1QnLFxyXG4gICAgICAgICAgICBoZWFkZXJzOiB7XHJcbiAgICAgICAgICAgICAgICAuLi5oZWFkZXJzLFxyXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaWQnOiBjaHVua0lkLFxyXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaW5kZXgnOiBTdHJpbmcoaSksXHJcbiAgICAgICAgICAgICAgICAneC13YWlscy1jaHVuay10b3RhbCc6IFN0cmluZyh0b3RhbENodW5rcyksXHJcbiAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgIGJvZHk6IGNodW5rLFxyXG4gICAgICAgIH0pO1xyXG4gICAgICAgIGlmICghcmVzcC5vaykge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gZmV0Y2godXJsLCB7XHJcbiAgICAgICAgbWV0aG9kOiAnUE9TVCcsXHJcbiAgICAgICAgaGVhZGVyczoge1xyXG4gICAgICAgICAgICAuLi5oZWFkZXJzLFxyXG4gICAgICAgICAgICAneC13YWlscy1jaHVuay1pZCc6IGNodW5rSWQsXHJcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWluZGV4JzogU3RyaW5nKHRvdGFsQ2h1bmtzIC0gMSksXHJcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLXRvdGFsJzogU3RyaW5nKHRvdGFsQ2h1bmtzKSxcclxuICAgICAgICB9LFxyXG4gICAgICAgIGJvZHk6IGJvZHlCeXRlcy5zdWJhcnJheSgodG90YWxDaHVua3MgLSAxKSAqIENIVU5LX1RIUkVTSE9MRCksXHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuLy8gc2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5EaWFsb2cpO1xyXG5cclxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcclxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XHJcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xyXG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XHJcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcclxuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xyXG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xyXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xyXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXHJcbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC4gKi9cclxuICAgIEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xyXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cclxuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXHJcbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cclxuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXHJcbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xyXG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXHJcbiAgICBUaXRsZT86IHN0cmluZztcclxuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XHJcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXHJcbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xyXG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cclxuICAgIERldGFjaGVkPzogYm9vbGVhbjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBTYXZlRmlsZURpYWxvZ09wdGlvbnMge1xyXG4gICAgLyoqIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBGaWxlbmFtZT86IHN0cmluZztcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cclxuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cclxuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXHJcbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXHJcbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cclxuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cclxuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXHJcbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xyXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXHJcbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XHJcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xyXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcclxuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgVGl0bGU/OiBzdHJpbmc7XHJcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBNZXNzYWdlPzogc3RyaW5nO1xyXG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xyXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcclxuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXHJcbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgTWVzc2FnZURpYWxvZ09wdGlvbnMge1xyXG4gICAgLyoqIFRoZSB0aXRsZSBvZiB0aGUgZGlhbG9nIHdpbmRvdy4gKi9cclxuICAgIFRpdGxlPzogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSBtYWluIG1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgTWVzc2FnZT86IHN0cmluZztcclxuICAgIC8qKiBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBCdXR0b25zPzogQnV0dG9uW107XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLiAqL1xyXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIEJ1dHRvbiB7XHJcbiAgICAvKiogVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uICovXHJcbiAgICBMYWJlbD86IHN0cmluZztcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLiAqL1xyXG4gICAgSXNDYW5jZWw/OiBib29sZWFuO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci4gKi9cclxuICAgIElzRGVmYXVsdD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgRmlsZUZpbHRlciB7XHJcbiAgICAvKiogRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuICovXHJcbiAgICBEaXNwbGF5TmFtZT86IHN0cmluZztcclxuICAgIC8qKiBQYXR0ZXJuIHRvIG1hdGNoIGZvciB0aGUgZmlsdGVyLCBlLmcuIFwiKi50eHQ7Ki5tZFwiIGZvciB0ZXh0IG1hcmtkb3duIGZpbGVzLiAqL1xyXG4gICAgUGF0dGVybj86IHN0cmluZztcclxufVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSB0eXBlIC0gRGlhbG9nIHR5cGUuXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyBmb3IgdGhlIGRpYWxvZy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nLlxyXG4gKi9cclxuZnVuY3Rpb24gZGlhbG9nKHR5cGU6IG51bWJlciwgb3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMgfCBPcGVuRmlsZURpYWxvZ09wdGlvbnMgfCBTYXZlRmlsZURpYWxvZ09wdGlvbnMgPSB7fSk6IFByb21pc2U8YW55PiB7XHJcbiAgICByZXR1cm4gY2FsbCh0eXBlLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGFuIGluZm8gZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSB3YXJuaW5nIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhbiBlcnJvciBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIHF1ZXN0aW9uIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgb25lIG9yIG1vcmUgZmlsZXMgdG8gb3Blbi5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cclxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSBvciBsaXN0IG9mIGZpbGVzLCBvciBhIGJsYW5rIHN0cmluZy9lbXB0eSBsaXN0IGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbjogdHJ1ZSB9KTogUHJvbWlzZTxzdHJpbmdbXT47XHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogZmFsc2UgfCB1bmRlZmluZWQgfSk6IFByb21pc2U8c3RyaW5nPjtcclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+O1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT4geyByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKSA/PyBbXTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgYSBmaWxlIHRvIHNhdmUuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUsIG9yIGEgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9uczogU2F2ZUZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7IH1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5pbXBvcnQgeyBldmVudExpc3RlbmVycywgTGlzdGVuZXIsIGxpc3RlbmVyT2ZmIH0gZnJvbSBcIi4vbGlzdGVuZXIuanNcIjtcclxuaW1wb3J0IHsgRXZlbnRzIGFzIENyZWF0ZSB9IGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xyXG5pbXBvcnQgeyBUeXBlcyB9IGZyb20gXCIuL2V2ZW50X3R5cGVzLmpzXCI7XHJcblxyXG4vLyBTZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5FdmVudHMpO1xyXG5jb25zdCBFbWl0TWV0aG9kID0gMDtcclxuXHJcbmV4cG9ydCAqIGZyb20gXCIuL2V2ZW50X3R5cGVzLmpzXCI7XHJcblxyXG4vKipcclxuICogQSB0YWJsZSBvZiBkYXRhIHR5cGVzIGZvciBhbGwga25vd24gZXZlbnRzLlxyXG4gKiBXaWxsIGJlIG1vbmtleS1wYXRjaGVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cclxuICovXHJcbmV4cG9ydCBpbnRlcmZhY2UgQ3VzdG9tRXZlbnRzIHt9XHJcblxyXG4vKipcclxuICogRWl0aGVyIGEga25vd24gZXZlbnQgbmFtZSBvciBhbiBhcmJpdHJhcnkgc3RyaW5nLlxyXG4gKi9cclxuZXhwb3J0IHR5cGUgV2FpbHNFdmVudE5hbWU8RSBleHRlbmRzIGtleW9mIEN1c3RvbUV2ZW50cyA9IGtleW9mIEN1c3RvbUV2ZW50cz4gPSBFIHwgKHN0cmluZyAmIHt9KTtcclxuXHJcbi8qKlxyXG4gKiBVbmlvbiBvZiBhbGwga25vd24gc3lzdGVtIGV2ZW50IG5hbWVzLlxyXG4gKi9cclxudHlwZSBTeXN0ZW1FdmVudE5hbWUgPSB7XHJcbiAgICBbSyBpbiBrZXlvZiAodHlwZW9mIFR5cGVzKV06ICh0eXBlb2YgVHlwZXMpW0tdW2tleW9mICgodHlwZW9mIFR5cGVzKVtLXSldXHJcbn0gZXh0ZW5kcyAoaW5mZXIgTSkgPyBNW2tleW9mIE1dIDogbmV2ZXI7XHJcblxyXG4vKipcclxuICogVGhlIGRhdGEgdHlwZSBhc3NvY2lhdGVkIHRvIGEgZ2l2ZW4gZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgdHlwZSBXYWlsc0V2ZW50RGF0YTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4gPVxyXG4gICAgRSBleHRlbmRzIGtleW9mIEN1c3RvbUV2ZW50cyA/IEN1c3RvbUV2ZW50c1tFXSA6IChFIGV4dGVuZHMgU3lzdGVtRXZlbnROYW1lID8gdm9pZCA6IGFueSk7XHJcblxyXG4vKipcclxuICogVGhlIHR5cGUgb2YgaGFuZGxlcnMgZm9yIGEgZ2l2ZW4gZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgdHlwZSBXYWlsc0V2ZW50Q2FsbGJhY2s8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+ID0gKGV2OiBXYWlsc0V2ZW50PEU+KSA9PiB2b2lkO1xyXG5cclxuLyoqXHJcbiAqIFJlcHJlc2VudHMgYSBzeXN0ZW0gZXZlbnQgb3IgYSBjdXN0b20gZXZlbnQgZW1pdHRlZCB0aHJvdWdoIHdhaWxzLXByb3ZpZGVkIGZhY2lsaXRpZXMuXHJcbiAqL1xyXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4ge1xyXG4gICAgLyoqXHJcbiAgICAgKiBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXHJcbiAgICAgKi9cclxuICAgIG5hbWU6IEU7XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBPcHRpb25hbCBkYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZW1pdHRlZCBldmVudC5cclxuICAgICAqL1xyXG4gICAgZGF0YTogV2FpbHNFdmVudERhdGE8RT47XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBOYW1lIG9mIHRoZSBvcmlnaW5hdGluZyB3aW5kb3cuIE9taXR0ZWQgZm9yIGFwcGxpY2F0aW9uIGV2ZW50cy5cclxuICAgICAqIFdpbGwgYmUgb3ZlcnJpZGRlbiBpZiBzZXQgbWFudWFsbHkuXHJcbiAgICAgKi9cclxuICAgIHNlbmRlcj86IHN0cmluZztcclxuXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBFLCBkYXRhOiBXYWlsc0V2ZW50RGF0YTxFPik7XHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiBleHRlbmRzIG51bGwgfCB2b2lkID8gRSA6IG5ldmVyKVxyXG4gICAgY29uc3RydWN0b3IobmFtZTogRSwgZGF0YT86IGFueSkge1xyXG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XHJcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YSA/PyBudWxsO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQ6IGFueSkge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcclxuICAgIGlmICghbGlzdGVuZXJzKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGxldCB3YWlsc0V2ZW50ID0gbmV3IFdhaWxzRXZlbnQoXHJcbiAgICAgICAgZXZlbnQubmFtZSxcclxuICAgICAgICAoZXZlbnQubmFtZSBpbiBDcmVhdGUpID8gQ3JlYXRlW2V2ZW50Lm5hbWVdKGV2ZW50LmRhdGEpIDogZXZlbnQuZGF0YVxyXG4gICAgKTtcclxuICAgIGlmICgnc2VuZGVyJyBpbiBldmVudCkge1xyXG4gICAgICAgIHdhaWxzRXZlbnQuc2VuZGVyID0gZXZlbnQuc2VuZGVyO1xyXG4gICAgfVxyXG5cclxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4gIWxpc3RlbmVyLmRpc3BhdGNoKHdhaWxzRXZlbnQpKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXHJcbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxyXG4gKiBAcGFyYW0gbWF4Q2FsbGJhY2tzIC0gVGhlIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRoZSBjYWxsYmFjayBjYW4gYmUgY2FsbGVkIGZvciB0aGUgZXZlbnQuIE9uY2UgdGhlIG1heGltdW0gbnVtYmVyIGlzIHJlYWNoZWQsIHRoZSBjYWxsYmFjayB3aWxsIG5vIGxvbmdlciBiZSBjYWxsZWQuXHJcbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+LCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XHJcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxyXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbjxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+KTogKCkgPT4gdm9pZCB7XHJcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXHJcbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25jZTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+KTogKCkgPT4gdm9pZCB7XHJcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlbW92ZXMgZXZlbnQgbGlzdGVuZXJzIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50IG5hbWVzLlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lcyAtIFRoZSBuYW1lIG9mIHRoZSBldmVudHMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmKC4uLmV2ZW50TmFtZXM6IFtXYWlsc0V2ZW50TmFtZSwgLi4uV2FpbHNFdmVudE5hbWVbXV0pOiB2b2lkIHtcclxuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVtb3ZlcyBhbGwgZXZlbnQgbGlzdGVuZXJzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpOiB2b2lkIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0cyBhbiBldmVudC5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgd2lsbCBiZSBmdWxmaWxsZWQgb25jZSB0aGUgZXZlbnQgaGFzIGJlZW4gZW1pdHRlZC4gIFJlc29sdmVzIHRvIHRydWUgaWYgdGhlIGV2ZW50IHdhcyBjYW5jZWxsZWQuXHJcbiAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIGVtaXRcclxuICogQHBhcmFtIGRhdGEgLSBUaGUgZGF0YSB0aGF0IHdpbGwgYmUgc2VudCB3aXRoIHRoZSBldmVudFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IEUsIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+KTogUHJvbWlzZTxib29sZWFuPlxyXG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogV2FpbHNFdmVudERhdGE8RT4gZXh0ZW5kcyBudWxsIHwgdm9pZCA/IEUgOiBuZXZlcik6IFByb21pc2U8Ym9vbGVhbj5cclxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IFdhaWxzRXZlbnREYXRhPEU+LCBkYXRhPzogYW55KTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCAgbmV3IFdhaWxzRXZlbnQobmFtZSwgZGF0YSkpXHJcbn1cclxuXHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vLyBUaGUgZm9sbG93aW5nIHV0aWxpdGllcyBoYXZlIGJlZW4gZmFjdG9yZWQgb3V0IG9mIC4vZXZlbnRzLnRzXHJcbi8vIGZvciB0ZXN0aW5nIHB1cnBvc2VzLlxyXG5cclxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcDxzdHJpbmcsIExpc3RlbmVyW10+KCk7XHJcblxyXG5leHBvcnQgY2xhc3MgTGlzdGVuZXIge1xyXG4gICAgZXZlbnROYW1lOiBzdHJpbmc7XHJcbiAgICBjYWxsYmFjazogKGRhdGE6IGFueSkgPT4gdm9pZDtcclxuICAgIG1heENhbGxiYWNrczogbnVtYmVyO1xyXG5cclxuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZTogc3RyaW5nLCBjYWxsYmFjazogKGRhdGE6IGFueSkgPT4gdm9pZCwgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcclxuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcclxuICAgICAgICB0aGlzLmNhbGxiYWNrID0gY2FsbGJhY2s7XHJcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XHJcbiAgICB9XHJcblxyXG4gICAgZGlzcGF0Y2goZGF0YTogYW55KTogYm9vbGVhbiB7XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgdGhpcy5jYWxsYmFjayhkYXRhKTtcclxuICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlcnIpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkgcmV0dXJuIGZhbHNlO1xyXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XHJcbiAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXI6IExpc3RlbmVyKTogdm9pZCB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGxpc3RlbmVyLmV2ZW50TmFtZSk7XHJcbiAgICBpZiAoIWxpc3RlbmVycykge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xyXG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcclxuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUobGlzdGVuZXIuZXZlbnROYW1lKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGxpc3RlbmVyLmV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIEFueSBpcyBhIGR1bW15IGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBzaW1wbGUgb3IgdW5rbm93biB0eXBlcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBBbnk8VCA9IGFueT4oc291cmNlOiBhbnkpOiBUIHtcclxuICAgIHJldHVybiBzb3VyY2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXHJcbiAqIG51bGwgc3RyaW5ncyB3aXRoIGVtcHR5IHN0cmluZ3MuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnl0ZVNsaWNlKHNvdXJjZTogYW55KTogc3RyaW5nIHtcclxuICAgIHJldHVybiAoKHNvdXJjZSA9PSBudWxsKSA/IFwiXCIgOiBzb3VyY2UpO1xyXG59XHJcblxyXG4vKipcclxuICogQXJyYXkgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XHJcbiAqIHdob3NlIGVsZW1lbnRzIGFyZSBvZiB0aGF0IHR5cGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQXJyYXk8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiBUW10ge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm4gW107XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XHJcbiAgICAgICAgICAgIHNvdXJjZVtpXSA9IGVsZW1lbnQoc291cmNlW2ldKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcclxuICogd2hvc2Uga2V5cyBhbmQgdmFsdWVzIGFyZSBvZiB0aG9zZSB0eXBlcy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBNYXA8SyBleHRlbmRzIFByb3BlcnR5S2V5ID0gYW55LCBWID0gYW55PihrZXk6IChzb3VyY2U6IGFueSkgPT4gSywgdmFsdWU6IChzb3VyY2U6IGFueSkgPT4gVik6IChzb3VyY2U6IGFueSkgPT4gUmVjb3JkPEssIFY+IHtcclxuICAgIGlmICh2YWx1ZSA9PT0gQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyB7fSA6IHNvdXJjZSk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcclxuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XHJcbiAgICAgICAgICAgIHJldHVybiB7fTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZm9yIChjb25zdCBrZXkgaW4gc291cmNlKSB7XHJcbiAgICAgICAgICAgIHNvdXJjZVtrZXldID0gdmFsdWUoc291cmNlW2tleV0pO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gc291cmNlO1xyXG4gICAgfTtcclxufVxyXG5cclxuLyoqXHJcbiAqIE51bGxhYmxlIHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXHJcbiAqIGFuZCByZXR1cm5zIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgbnVsbGFibGUgdmFsdWUgb2YgdGhhdCB0eXBlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE51bGxhYmxlPFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gKFQgfCBudWxsKSB7XHJcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIEFueTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IG51bGwgOiBlbGVtZW50KHNvdXJjZSkpO1xyXG59XHJcblxyXG4vKipcclxuICogU3RydWN0IHRha2VzIGFuIG9iamVjdCBtYXBwaW5nIGZpZWxkIG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9uc1xyXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBzdHJ1Y3QuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU3RydWN0KGNyZWF0ZUZpZWxkOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4pOlxyXG4gICAgPFUgZXh0ZW5kcyBSZWNvcmQ8c3RyaW5nLCBhbnk+ID0gYW55Pihzb3VyY2U6IGFueSkgPT4gVVxyXG57XHJcbiAgICBsZXQgYWxsQW55ID0gdHJ1ZTtcclxuICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xyXG4gICAgICAgIGlmIChjcmVhdGVGaWVsZFtuYW1lXSAhPT0gQW55KSB7XHJcbiAgICAgICAgICAgIGFsbEFueSA9IGZhbHNlO1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcbiAgICB9XHJcbiAgICBpZiAoYWxsQW55KSB7XHJcbiAgICAgICAgcmV0dXJuIEFueTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xyXG4gICAgICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xyXG4gICAgICAgICAgICBpZiAobmFtZSBpbiBzb3VyY2UpIHtcclxuICAgICAgICAgICAgICAgIHNvdXJjZVtuYW1lXSA9IGNyZWF0ZUZpZWxkW25hbWVdKHNvdXJjZVtuYW1lXSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBNYXBzIGtub3duIGV2ZW50IG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdGhlaXIgZGF0YSB0eXBlcy5cclxuICogV2lsbCBiZSBtb25rZXktcGF0Y2hlZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgRXZlbnRzOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4gPSB7fTtcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8vIEN5bmh5cmNod3lkIHkgZmZlaWwgaG9uIHluIGF3dG9tYXRpZy4gUEVJRElXQ0ggXHUwMEMyIE1PRElXTFxyXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXHJcblxyXG5leHBvcnQgY29uc3QgVHlwZXMgPSBPYmplY3QuZnJlZXplKHtcclxuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcclxuXHRcdEFQTVBvd2VyU2V0dGluZ0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU2V0dGluZ0NoYW5nZVwiLFxyXG5cdFx0QVBNUG93ZXJTdGF0dXNDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclN0YXR1c0NoYW5nZVwiLFxyXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXHJcblx0XHRBUE1SZXN1bWVTdXNwZW5kOiBcIndpbmRvd3M6QVBNUmVzdW1lU3VzcGVuZFwiLFxyXG5cdFx0QVBNU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVN1c3BlbmRcIixcclxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxyXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXHJcblx0XHRXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZDogXCJ3aW5kb3dzOldlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkXCIsXHJcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcclxuXHRcdFdpbmRvd0JhY2tncm91bmRFcmFzZTogXCJ3aW5kb3dzOldpbmRvd0JhY2tncm91bmRFcmFzZVwiLFxyXG5cdFx0V2luZG93Q2xpY2tBY3RpdmU6IFwid2luZG93czpXaW5kb3dDbGlja0FjdGl2ZVwiLFxyXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcclxuXHRcdFdpbmRvd0RpZE1vdmU6IFwid2luZG93czpXaW5kb3dEaWRNb3ZlXCIsXHJcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwid2luZG93czpXaW5kb3dEaWRSZXNpemVcIixcclxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXHJcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXHJcblx0XHRXaW5kb3dEcmFnRW50ZXI6IFwid2luZG93czpXaW5kb3dEcmFnRW50ZXJcIixcclxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxyXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxyXG5cdFx0V2luZG93RW5kTW92ZTogXCJ3aW5kb3dzOldpbmRvd0VuZE1vdmVcIixcclxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxyXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd0Z1bGxzY3JlZW5cIixcclxuXHRcdFdpbmRvd0hpZGU6IFwid2luZG93czpXaW5kb3dIaWRlXCIsXHJcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXHJcblx0XHRXaW5kb3dLZXlEb3duOiBcIndpbmRvd3M6V2luZG93S2V5RG93blwiLFxyXG5cdFx0V2luZG93S2V5VXA6IFwid2luZG93czpXaW5kb3dLZXlVcFwiLFxyXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXHJcblx0XHRXaW5kb3dOb25DbGllbnRIaXQ6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRIaXRcIixcclxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlRG93bjogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlRG93blwiLFxyXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcclxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlTW92ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTW92ZVwiLFxyXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VVcDogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlVXBcIixcclxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcclxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXHJcblx0XHRXaW5kb3dTZXRGb2N1czogXCJ3aW5kb3dzOldpbmRvd1NldEZvY3VzXCIsXHJcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxyXG5cdFx0V2luZG93U3RhcnRNb3ZlOiBcIndpbmRvd3M6V2luZG93U3RhcnRNb3ZlXCIsXHJcblx0XHRXaW5kb3dTdGFydFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0UmVzaXplXCIsXHJcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcclxuXHRcdFdpbmRvd1pPcmRlckNoYW5nZWQ6IFwid2luZG93czpXaW5kb3daT3JkZXJDaGFuZ2VkXCIsXHJcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJ3aW5kb3dzOldpbmRvd01pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxyXG5cdFx0V2luZG93TWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dNYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJ3aW5kb3dzOldpbmRvd1VuTWF4aW1pc2VcIixcclxuXHR9KSxcclxuXHRNYWM6IE9iamVjdC5mcmVlemUoe1xyXG5cdFx0QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUljb25cIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0JhckZyYW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0JhckZyYW1lXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmdcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFVuaGlkZVwiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXHJcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW5cIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmVcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsVXBkYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVcGRhdGVcIixcclxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxyXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXHJcblx0XHRNZW51RGlkQ2xvc2U6IFwibWFjOk1lbnVEaWRDbG9zZVwiLFxyXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcclxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXHJcblx0XHRNZW51RGlkSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudURpZEhpZ2hsaWdodEl0ZW1cIixcclxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxyXG5cdFx0TWVudURpZFBvcFVwOiBcIm1hYzpNZW51RGlkUG9wVXBcIixcclxuXHRcdE1lbnVEaWRSZW1vdmVJdGVtOiBcIm1hYzpNZW51RGlkUmVtb3ZlSXRlbVwiLFxyXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXHJcblx0XHRNZW51RGlkU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW1cIixcclxuXHRcdE1lbnVEaWRVcGRhdGU6IFwibWFjOk1lbnVEaWRVcGRhdGVcIixcclxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXHJcblx0XHRNZW51V2lsbEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsQmVnaW5UcmFja2luZ1wiLFxyXG5cdFx0TWVudVdpbGxEaXNwbGF5SXRlbTogXCJtYWM6TWVudVdpbGxEaXNwbGF5SXRlbVwiLFxyXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxyXG5cdFx0TWVudVdpbGxIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51V2lsbEhpZ2hsaWdodEl0ZW1cIixcclxuXHRcdE1lbnVXaWxsT3BlbjogXCJtYWM6TWVudVdpbGxPcGVuXCIsXHJcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXHJcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxyXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25cIixcclxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXHJcblx0XHRNZW51V2lsbFVwZGF0ZTogXCJtYWM6TWVudVdpbGxVcGRhdGVcIixcclxuXHRcdFdlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXHJcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXHJcblx0XHRXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxyXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcclxuXHRcdFdpbmRvd0RpZEJlY29tZU1haW46IFwibWFjOldpbmRvd0RpZEJlY29tZU1haW5cIixcclxuXHRcdFdpbmRvd0RpZEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd0RpZEJlZ2luU2hlZXRcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nTG9jYXRpb25cIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllc1wiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZVwiLFxyXG5cdFx0V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlXCIsXHJcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcclxuXHRcdFdpbmRvd0RpZENoYW5nZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZENoYW5nZVRvb2xiYXJcIixcclxuXHRcdFdpbmRvd0RpZERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd0RpZERlbWluaWF0dXJpemVcIixcclxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxyXG5cdFx0V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW5cIixcclxuXHRcdFdpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd0RpZEVudGVyVmVyc2lvbkJyb3dzZXJcIixcclxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxyXG5cdFx0V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXJcIixcclxuXHRcdFdpbmRvd0RpZEV4cG9zZTogXCJtYWM6V2luZG93RGlkRXhwb3NlXCIsXHJcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcclxuXHRcdFdpbmRvd0RpZE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWRNaW5pYXR1cml6ZVwiLFxyXG5cdFx0V2luZG93RGlkTW92ZTogXCJtYWM6V2luZG93RGlkTW92ZVwiLFxyXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dEaWRPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9uU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dEaWRSZXNpZ25LZXk6IFwibWFjOldpbmRvd0RpZFJlc2lnbktleVwiLFxyXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxyXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIm1hYzpXaW5kb3dEaWRSZXNpemVcIixcclxuXHRcdFdpbmRvd0RpZFVwZGF0ZTogXCJtYWM6V2luZG93RGlkVXBkYXRlXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcclxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcclxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXHJcblx0XHRXaW5kb3dEaWRVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93RGlkVXBkYXRlVGl0bGVcIixcclxuXHRcdFdpbmRvd0RpZFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRvb2xiYXJcIixcclxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcclxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWRcIixcclxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nRXhpdGVkXCIsXHJcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxyXG5cdFx0V2luZG93SGlkZTogXCJtYWM6V2luZG93SGlkZVwiLFxyXG5cdFx0V2luZG93TWF4aW1pc2U6IFwibWFjOldpbmRvd01heGltaXNlXCIsXHJcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXHJcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJtYWM6V2luZG93TWluaW1pc2VcIixcclxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwibWFjOldpbmRvd1VuTWluaW1pc2VcIixcclxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxyXG5cdFx0V2luZG93U2hvdzogXCJtYWM6V2luZG93U2hvd1wiLFxyXG5cdFx0V2luZG93V2lsbEJlY29tZUtleTogXCJtYWM6V2luZG93V2lsbEJlY29tZUtleVwiLFxyXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXHJcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcclxuXHRcdFdpbmRvd1dpbGxDaGFuZ2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd1dpbGxDaGFuZ2VPcmRlcmluZ01vZGVcIixcclxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXHJcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcclxuXHRcdFdpbmRvd1dpbGxFbnRlckZ1bGxTY3JlZW46IFwibWFjOldpbmRvd1dpbGxFbnRlckZ1bGxTY3JlZW5cIixcclxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxyXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcclxuXHRcdFdpbmRvd1dpbGxFeGl0VmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFeGl0VmVyc2lvbkJyb3dzZXJcIixcclxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXHJcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxyXG5cdFx0V2luZG93V2lsbE1vdmU6IFwibWFjOldpbmRvd1dpbGxNb3ZlXCIsXHJcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxyXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXHJcblx0XHRXaW5kb3dXaWxsUmVzaWduTWFpbjogXCJtYWM6V2luZG93V2lsbFJlc2lnbk1haW5cIixcclxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcclxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVwiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcclxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllczogXCJtYWM6V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxyXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRvb2xiYXI6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUb29sYmFyXCIsXHJcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcclxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxyXG5cdFx0V2luZG93Wm9vbUluOiBcIm1hYzpXaW5kb3dab29tSW5cIixcclxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcclxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJtYWM6V2luZG93Wm9vbVJlc2V0XCIsXHJcblx0fSksXHJcblx0TGludXg6IE9iamVjdC5mcmVlemUoe1xyXG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxyXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcImxpbnV4OlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxyXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcclxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibGludXg6V2luZG93RGlkTW92ZVwiLFxyXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImxpbnV4OldpbmRvd0RpZFJlc2l6ZVwiLFxyXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXHJcblx0XHRXaW5kb3dGb2N1c091dDogXCJsaW51eDpXaW5kb3dGb2N1c091dFwiLFxyXG5cdFx0V2luZG93TG9hZFN0YXJ0ZWQ6IFwibGludXg6V2luZG93TG9hZFN0YXJ0ZWRcIixcclxuXHRcdFdpbmRvd0xvYWRSZWRpcmVjdGVkOiBcImxpbnV4OldpbmRvd0xvYWRSZWRpcmVjdGVkXCIsXHJcblx0XHRXaW5kb3dMb2FkQ29tbWl0dGVkOiBcImxpbnV4OldpbmRvd0xvYWRDb21taXR0ZWRcIixcclxuXHRcdFdpbmRvd0xvYWRGaW5pc2hlZDogXCJsaW51eDpXaW5kb3dMb2FkRmluaXNoZWRcIixcclxuXHR9KSxcclxuXHRpT1M6IE9iamVjdC5mcmVlemUoe1xyXG5cdFx0QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmU6IFwiaW9zOkFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlXCIsXHJcblx0XHRBcHBsaWNhdGlvbkRpZEVudGVyQmFja2dyb3VuZDogXCJpb3M6QXBwbGljYXRpb25EaWRFbnRlckJhY2tncm91bmRcIixcclxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcImlvczpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxyXG5cdFx0QXBwbGljYXRpb25EaWRSZWNlaXZlTWVtb3J5V2FybmluZzogXCJpb3M6QXBwbGljYXRpb25EaWRSZWNlaXZlTWVtb3J5V2FybmluZ1wiLFxyXG5cdFx0QXBwbGljYXRpb25XaWxsRW50ZXJGb3JlZ3JvdW5kOiBcImlvczpBcHBsaWNhdGlvbldpbGxFbnRlckZvcmVncm91bmRcIixcclxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJpb3M6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXHJcblx0XHRBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGU6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZVwiLFxyXG5cdFx0V2luZG93RGlkTG9hZDogXCJpb3M6V2luZG93RGlkTG9hZFwiLFxyXG5cdFx0V2luZG93V2lsbEFwcGVhcjogXCJpb3M6V2luZG93V2lsbEFwcGVhclwiLFxyXG5cdFx0V2luZG93RGlkQXBwZWFyOiBcImlvczpXaW5kb3dEaWRBcHBlYXJcIixcclxuXHRcdFdpbmRvd1dpbGxEaXNhcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxEaXNhcHBlYXJcIixcclxuXHRcdFdpbmRvd0RpZERpc2FwcGVhcjogXCJpb3M6V2luZG93RGlkRGlzYXBwZWFyXCIsXHJcblx0XHRXaW5kb3dTYWZlQXJlYUluc2V0c0NoYW5nZWQ6IFwiaW9zOldpbmRvd1NhZmVBcmVhSW5zZXRzQ2hhbmdlZFwiLFxyXG5cdFx0V2luZG93T3JpZW50YXRpb25DaGFuZ2VkOiBcImlvczpXaW5kb3dPcmllbnRhdGlvbkNoYW5nZWRcIixcclxuXHRcdFdpbmRvd1RvdWNoQmVnYW46IFwiaW9zOldpbmRvd1RvdWNoQmVnYW5cIixcclxuXHRcdFdpbmRvd1RvdWNoTW92ZWQ6IFwiaW9zOldpbmRvd1RvdWNoTW92ZWRcIixcclxuXHRcdFdpbmRvd1RvdWNoRW5kZWQ6IFwiaW9zOldpbmRvd1RvdWNoRW5kZWRcIixcclxuXHRcdFdpbmRvd1RvdWNoQ2FuY2VsbGVkOiBcImlvczpXaW5kb3dUb3VjaENhbmNlbGxlZFwiLFxyXG5cdFx0V2ViVmlld0RpZFN0YXJ0TmF2aWdhdGlvbjogXCJpb3M6V2ViVmlld0RpZFN0YXJ0TmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXHJcblx0XHRXZWJWaWV3RGlkRmFpbE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGYWlsTmF2aWdhdGlvblwiLFxyXG5cdFx0V2ViVmlld0RlY2lkZVBvbGljeUZvck5hdmlnYXRpb25BY3Rpb246IFwiaW9zOldlYlZpZXdEZWNpZGVQb2xpY3lGb3JOYXZpZ2F0aW9uQWN0aW9uXCIsXHJcblx0fSksXHJcblx0Q29tbW9uOiBPYmplY3QuZnJlZXplKHtcclxuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiY29tbW9uOkFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGVcIixcclxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXHJcblx0XHRBcHBsaWNhdGlvbkxhdW5jaGVkV2l0aFVybDogXCJjb21tb246QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmxcIixcclxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXHJcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXHJcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImNvbW1vbjpXaW5kb3dEaWRNb3ZlXCIsXHJcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwiY29tbW9uOldpbmRvd0RpZFJlc2l6ZVwiLFxyXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxyXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcclxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxyXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxyXG5cdFx0V2luZG93SGlkZTogXCJjb21tb246V2luZG93SGlkZVwiLFxyXG5cdFx0V2luZG93TG9zdEZvY3VzOiBcImNvbW1vbjpXaW5kb3dMb3N0Rm9jdXNcIixcclxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxyXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dUb2dnbGVGcmFtZWxlc3M6IFwiY29tbW9uOldpbmRvd1RvZ2dsZUZyYW1lbGVzc1wiLFxyXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxyXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcclxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcclxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXHJcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1heGltaXNlXCIsXHJcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXHJcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXHJcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxyXG5cdFx0V2luZG93Wm9vbU91dDogXCJjb21tb246V2luZG93Wm9vbU91dFwiLFxyXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcclxuXHR9KSxcclxufSk7XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKipcclxuICogTG9ncyBhIG1lc3NhZ2UgdG8gdGhlIGNvbnNvbGUgd2l0aCBjdXN0b20gZm9ybWF0dGluZy5cclxuICpcclxuICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZTogYW55KSB7XHJcbiAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcclxuICAgIGNvbnNvbGUubG9nKFxyXG4gICAgICAgICclYyB3YWlsczMgJWMgJyArIG1lc3NhZ2UgKyAnICcsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNhYTAwMDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAzcHggMHB4IDBweCAzcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjMDA5OTAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogMHB4IDNweCAzcHggMHB4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJ1xyXG4gICAgKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSB3ZWJ2aWV3IHN1cHBvcnRzIHRoZSB7QGxpbmsgTW91c2VFdmVudCNidXR0b25zfSBwcm9wZXJ0eS5cclxuICogTG9va2luZyBhdCB5b3UgbWFjT1MgSGlnaCBTaWVycmEhXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gY2FuVHJhY2tCdXR0b25zKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIChuZXcgTW91c2VFdmVudCgnbW91c2Vkb3duJykpLmJ1dHRvbnMgPT09IDA7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgYnJvd3NlciBzdXBwb3J0cyByZW1vdmluZyBsaXN0ZW5lcnMgYnkgdHJpZ2dlcmluZyBhbiBBYm9ydFNpZ25hbFxyXG4gKiAoc2VlIGh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnL2VuLVVTL2RvY3MvV2ViL0FQSS9FdmVudFRhcmdldC9hZGRFdmVudExpc3RlbmVyI3NpZ25hbCkuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XHJcbiAgICBpZiAoIUV2ZW50VGFyZ2V0IHx8ICFBYm9ydFNpZ25hbCB8fCAhQWJvcnRDb250cm9sbGVyKVxyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuXHJcbiAgICBsZXQgcmVzdWx0ID0gdHJ1ZTtcclxuXHJcbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcclxuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB0YXJnZXQuYWRkRXZlbnRMaXN0ZW5lcigndGVzdCcsICgpID0+IHsgcmVzdWx0ID0gZmFsc2U7IH0sIHsgc2lnbmFsOiBjb250cm9sbGVyLnNpZ25hbCB9KTtcclxuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcclxuICAgIHRhcmdldC5kaXNwYXRjaEV2ZW50KG5ldyBDdXN0b21FdmVudCgndGVzdCcpKTtcclxuXHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKipcclxuICogUmVzb2x2ZXMgdGhlIGNsb3Nlc3QgSFRNTEVsZW1lbnQgYW5jZXN0b3Igb2YgYW4gZXZlbnQncyB0YXJnZXQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZXZlbnRUYXJnZXQoZXZlbnQ6IEV2ZW50KTogSFRNTEVsZW1lbnQge1xyXG4gICAgaWYgKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSB7XHJcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldDtcclxuICAgIH0gZWxzZSBpZiAoIShldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkgJiYgZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgTm9kZSkge1xyXG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQucGFyZW50RWxlbWVudCA/PyBkb2N1bWVudC5ib2R5O1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICByZXR1cm4gZG9jdW1lbnQuYm9keTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqKlxyXG4gVGhpcyB0ZWNobmlxdWUgZm9yIHByb3BlciBsb2FkIGRldGVjdGlvbiBpcyB0YWtlbiBmcm9tIEhUTVg6XHJcblxyXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcclxuXHJcbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcclxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXHJcblxyXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XHJcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcclxuXHJcbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcclxuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxyXG5cclxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcclxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cclxuIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLlxyXG5cclxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXHJcbiBBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFXHJcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcclxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcclxuIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMXHJcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxyXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcclxuIENBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksXHJcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxyXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cclxuXHJcbiAqKiovXHJcblxyXG5sZXQgaXNSZWFkeSA9IGZhbHNlO1xyXG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgKCkgPT4geyBpc1JlYWR5ID0gdHJ1ZSB9KTtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiB3aGVuUmVhZHkoY2FsbGJhY2s6ICgpID0+IHZvaWQpIHtcclxuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcclxuICAgICAgICBjYWxsYmFjaygpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgY2FsbGJhY2spO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB0eXBlIHsgU2NyZWVuIH0gZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xyXG5cclxuLy8gRHJvcCB0YXJnZXQgY29uc3RhbnRzXHJcbmNvbnN0IERST1BfVEFSR0VUX0FUVFJJQlVURSA9ICdkYXRhLWZpbGUtZHJvcC10YXJnZXQnO1xyXG5jb25zdCBEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MgPSAnZmlsZS1kcm9wLXRhcmdldC1hY3RpdmUnO1xyXG5sZXQgY3VycmVudERyb3BUYXJnZXQ6IEVsZW1lbnQgfCBudWxsID0gbnVsbDtcclxuXHJcbmNvbnN0IFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDA7XHJcbmNvbnN0IENlbnRlck1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDE7XHJcbmNvbnN0IENsb3NlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDI7XHJcbmNvbnN0IERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICA9IDM7XHJcbmNvbnN0IEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgICA9IDQ7XHJcbmNvbnN0IEZvY3VzTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDU7XHJcbmNvbnN0IEZvcmNlUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICA9IDY7XHJcbmNvbnN0IEZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICA9IDc7XHJcbmNvbnN0IEdldFNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDg7XHJcbmNvbnN0IEdldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDk7XHJcbmNvbnN0IEhlaWdodE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDEwO1xyXG5jb25zdCBIaWRlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxMTtcclxuY29uc3QgSXNGb2N1c2VkTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMTI7XHJcbmNvbnN0IElzRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDEzO1xyXG5jb25zdCBJc01heGltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNDtcclxuY29uc3QgSXNNaW5pbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTU7XHJcbmNvbnN0IE1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE2O1xyXG5jb25zdCBNaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNztcclxuY29uc3QgTmFtZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTg7XHJcbmNvbnN0IE9wZW5EZXZUb29sc01ldGhvZCAgICAgICAgICAgICAgICA9IDE5O1xyXG5jb25zdCBSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgPSAyMDtcclxuY29uc3QgUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMjE7XHJcbmNvbnN0IFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDIyO1xyXG5jb25zdCBSZXN0b3JlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAyMztcclxuY29uc3QgU2V0UG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgID0gMjQ7XHJcbmNvbnN0IFNldEFsd2F5c09uVG9wTWV0aG9kICAgICAgICAgICAgICA9IDI1O1xyXG5jb25zdCBTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kICAgICAgICAgPSAyNjtcclxuY29uc3QgU2V0RnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgICAgID0gMjc7XHJcbmNvbnN0IFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kICA9IDI4O1xyXG5jb25zdCBTZXRNYXhTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAyOTtcclxuY29uc3QgU2V0TWluU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMzA7XHJcbmNvbnN0IFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICA9IDMxO1xyXG5jb25zdCBTZXRSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgPSAzMjtcclxuY29uc3QgU2V0U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzM7XHJcbmNvbnN0IFNldFRpdGxlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDM0O1xyXG5jb25zdCBTZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzNTtcclxuY29uc3QgU2hvd01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzY7XHJcbmNvbnN0IFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM3O1xyXG5jb25zdCBUb2dnbGVGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgPSAzODtcclxuY29uc3QgVG9nZ2xlTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgID0gMzk7XHJcbmNvbnN0IFRvZ2dsZUZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICA9IDQwOyBcclxuY29uc3QgVW5GdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gNDE7XHJcbmNvbnN0IFVuTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQyO1xyXG5jb25zdCBVbk1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MztcclxuY29uc3QgV2lkdGhNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XHJcbmNvbnN0IFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDQ1O1xyXG5jb25zdCBab29tSW5NZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSA0NjtcclxuY29uc3QgWm9vbU91dE1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gNDc7XHJcbmNvbnN0IFpvb21SZXNldE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDQ4O1xyXG5jb25zdCBTbmFwQXNzaXN0TWV0aG9kICAgICAgICAgICAgICAgICAgPSA0OTtcclxuY29uc3QgRmlsZXNEcm9wcGVkICAgICAgICAgICAgICAgICAgICAgID0gNTA7XHJcbmNvbnN0IFByaW50TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDUxO1xyXG5jb25zdCBTZXRTY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICAgPSA1MjtcclxuXHJcbi8qKlxyXG4gKiBGaW5kcyB0aGUgbmVhcmVzdCBkcm9wIHRhcmdldCBlbGVtZW50IGJ5IHdhbGtpbmcgdXAgdGhlIERPTSB0cmVlLlxyXG4gKi9cclxuZnVuY3Rpb24gZ2V0RHJvcFRhcmdldEVsZW1lbnQoZWxlbWVudDogRWxlbWVudCB8IG51bGwpOiBFbGVtZW50IHwgbnVsbCB7XHJcbiAgICBpZiAoIWVsZW1lbnQpIHtcclxuICAgICAgICByZXR1cm4gbnVsbDtcclxuICAgIH1cclxuICAgIHJldHVybiBlbGVtZW50LmNsb3Nlc3QoYFske0RST1BfVEFSR0VUX0FUVFJJQlVURX1dYCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVjayBpZiB3ZSBjYW4gdXNlIFdlYlZpZXcyJ3MgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgKFdpbmRvd3MpXHJcbiAqIEFsc28gY2hlY2tzIHRoYXQgRW5hYmxlRmlsZURyb3AgaXMgdHJ1ZSBmb3IgdGhpcyB3aW5kb3cuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYW5SZXNvbHZlRmlsZVBhdGhzKCk6IGJvb2xlYW4ge1xyXG4gICAgLy8gTXVzdCBoYXZlIFdlYlZpZXcyJ3MgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgQVBJIChXaW5kb3dzIG9ubHkpXHJcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgPT0gbnVsbCkge1xyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgIH1cclxuICAgIC8vIE11c3QgaGF2ZSBFbmFibGVGaWxlRHJvcCBzZXQgdG8gdHJ1ZSBmb3IgdGhpcyB3aW5kb3dcclxuICAgIC8vIFRoaXMgZmxhZyBpcyBzZXQgYnkgdGhlIEdvIGJhY2tlbmQgZHVyaW5nIHJ1bnRpbWUgaW5pdGlhbGl6YXRpb25cclxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IHRydWU7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTZW5kIGZpbGUgZHJvcCB0byBiYWNrZW5kIHZpYSBXZWJWaWV3MiAoV2luZG93cyBvbmx5KVxyXG4gKi9cclxuZnVuY3Rpb24gcmVzb2x2ZUZpbGVQYXRocyh4OiBudW1iZXIsIHk6IG51bWJlciwgZmlsZXM6IEZpbGVbXSk6IHZvaWQge1xyXG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKSB7XHJcbiAgICAgICAgKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKGBmaWxlOmRyb3A6JHt4fToke3l9YCwgZmlsZXMpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBOYXRpdmUgZHJhZyBzdGF0ZSAoTGludXgvbWFjT1MgaW50ZXJjZXB0IERPTSBkcmFnIGV2ZW50cylcclxubGV0IG5hdGl2ZURyYWdBY3RpdmUgPSBmYWxzZTtcclxuXHJcbi8qKlxyXG4gKiBDbGVhbnMgdXAgbmF0aXZlIGRyYWcgc3RhdGUgYW5kIGhvdmVyIGVmZmVjdHMuXHJcbiAqIENhbGxlZCBvbiBkcm9wIG9yIHdoZW4gZHJhZyBsZWF2ZXMgdGhlIHdpbmRvdy5cclxuICovXHJcbmZ1bmN0aW9uIGNsZWFudXBOYXRpdmVEcmFnKCk6IHZvaWQge1xyXG4gICAgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xyXG4gICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxlZCBmcm9tIEdvIHdoZW4gYSBmaWxlIGRyYWcgZW50ZXJzIHRoZSB3aW5kb3cgb24gTGludXgvbWFjT1MuXHJcbiAqL1xyXG5mdW5jdGlvbiBoYW5kbGVEcmFnRW50ZXIoKTogdm9pZCB7XHJcbiAgICAvLyBDaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xyXG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IGFjdGl2YXRlIGRyYWcgc3RhdGVcclxuICAgIH1cclxuICAgIG5hdGl2ZURyYWdBY3RpdmUgPSB0cnVlO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbGVkIGZyb20gR28gd2hlbiBhIGZpbGUgZHJhZyBsZWF2ZXMgdGhlIHdpbmRvdyBvbiBMaW51eC9tYWNPUy5cclxuICovXHJcbmZ1bmN0aW9uIGhhbmRsZURyYWdMZWF2ZSgpOiB2b2lkIHtcclxuICAgIGNsZWFudXBOYXRpdmVEcmFnKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsZWQgZnJvbSBHbyBkdXJpbmcgZmlsZSBkcmFnIHRvIHVwZGF0ZSBob3ZlciBzdGF0ZSBvbiBMaW51eC9tYWNPUy5cclxuICogQHBhcmFtIHggLSBYIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xyXG4gKiBAcGFyYW0geSAtIFkgY29vcmRpbmF0ZSBpbiBDU1MgcGl4ZWxzXHJcbiAqL1xyXG5mdW5jdGlvbiBoYW5kbGVEcmFnT3Zlcih4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xyXG4gICAgaWYgKCFuYXRpdmVEcmFnQWN0aXZlKSByZXR1cm47XHJcbiAgICBcclxuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgZG9uJ3Qgc2hvdyBob3ZlciBlZmZlY3RzXHJcbiAgICB9XHJcbiAgICBcclxuICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHgsIHkpO1xyXG4gICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xyXG4gICAgXHJcbiAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcclxuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICB9XHJcbiAgICBcclxuICAgIGlmIChkcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgZHJvcFRhcmdldC5jbGFzc0xpc3QuYWRkKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBkcm9wVGFyZ2V0O1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICB9XHJcbn1cclxuXHJcblxyXG5cclxuLy8gRXhwb3J0IHRoZSBoYW5kbGVycyBmb3IgdXNlIGJ5IEdvIHZpYSBpbmRleC50c1xyXG5leHBvcnQgeyBoYW5kbGVEcmFnRW50ZXIsIGhhbmRsZURyYWdMZWF2ZSwgaGFuZGxlRHJhZ092ZXIgfTtcclxuXHJcbi8qKlxyXG4gKiBBIHJlY29yZCBkZXNjcmliaW5nIHRoZSBwb3NpdGlvbiBvZiBhIHdpbmRvdy5cclxuICovXHJcbmludGVyZmFjZSBQb3NpdGlvbiB7XHJcbiAgICAvKiogVGhlIGhvcml6b250YWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIHg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgdmVydGljYWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIHk6IG51bWJlcjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXHJcbiAqL1xyXG5pbnRlcmZhY2UgU2l6ZSB7XHJcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuICovXHJcbiAgICB3aWR0aDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIGhlaWdodDogbnVtYmVyO1xyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBjYWxsZXJTeW0gPSBTeW1ib2woXCJjYWxsZXJcIik7XHJcblxyXG5jbGFzcyBXaW5kb3cge1xyXG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXHJcbiAgICBwcml2YXRlIFtjYWxsZXJTeW1dOiAobWVzc2FnZTogbnVtYmVyLCBhcmdzPzogYW55KSA9PiBQcm9taXNlPGFueT47XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBJbml0aWFsaXNlcyBhIHdpbmRvdyBvYmplY3Qgd2l0aCB0aGUgc3BlY2lmaWVkIG5hbWUuXHJcbiAgICAgKlxyXG4gICAgICogQHByaXZhdGVcclxuICAgICAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHRhcmdldCB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IHN0cmluZyA9ICcnKSB7XHJcbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5XaW5kb3csIG5hbWUpXHJcblxyXG4gICAgICAgIC8vIGJpbmQgaW5zdGFuY2UgbWV0aG9kIHRvIG1ha2UgdGhlbSBlYXNpbHkgdXNhYmxlIGluIGV2ZW50IGhhbmRsZXJzXHJcbiAgICAgICAgZm9yIChjb25zdCBtZXRob2Qgb2YgT2JqZWN0LmdldE93blByb3BlcnR5TmFtZXMoV2luZG93LnByb3RvdHlwZSkpIHtcclxuICAgICAgICAgICAgaWYgKFxyXG4gICAgICAgICAgICAgICAgbWV0aG9kICE9PSBcImNvbnN0cnVjdG9yXCJcclxuICAgICAgICAgICAgICAgICYmIHR5cGVvZiAodGhpcyBhcyBhbnkpW21ldGhvZF0gPT09IFwiZnVuY3Rpb25cIlxyXG4gICAgICAgICAgICApIHtcclxuICAgICAgICAgICAgICAgICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9ICh0aGlzIGFzIGFueSlbbWV0aG9kXS5iaW5kKHRoaXMpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogR2V0cyB0aGUgc3BlY2lmaWVkIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gZ2V0LlxyXG4gICAgICogQHJldHVybnMgVGhlIGNvcnJlc3BvbmRpbmcgd2luZG93IG9iamVjdC5cclxuICAgICAqL1xyXG4gICAgR2V0KG5hbWU6IHN0cmluZyk6IFdpbmRvdyB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBXaW5kb3cobmFtZSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFBvc2l0aW9uTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdyBvbiB0aGUgc2NyZWVuLlxyXG4gICAgICovXHJcbiAgICBDZW50ZXIoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShDZW50ZXJNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIENsb3NlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2xvc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRGlzYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxyXG4gICAgICovXHJcbiAgICBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBFbmFibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cclxuICAgICAqL1xyXG4gICAgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgRm9jdXMoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb2N1c01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBGb3JjZXMgdGhlIHdpbmRvdyB0byByZWxvYWQgdGhlIHBhZ2UgYXNzZXRzLlxyXG4gICAgICovXHJcbiAgICBGb3JjZVJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvcmNlUmVsb2FkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFN3aXRjaGVzIHRoZSB3aW5kb3cgdG8gZnVsbHNjcmVlbiBtb2RlLlxyXG4gICAgICovXHJcbiAgICBGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBzY3JlZW4gdGhhdCB0aGUgd2luZG93IGlzIG9uLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgb24uXHJcbiAgICAgKi9cclxuICAgIEdldFNjcmVlbigpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0U2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHpvb20gbGV2ZWwuXHJcbiAgICAgKi9cclxuICAgIEdldFpvb20oKTogUHJvbWlzZTxudW1iZXI+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFpvb21NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIEhlaWdodCgpOiBQcm9taXNlPG51bWJlcj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSGVpZ2h0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEhpZGVzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIaWRlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZvY3VzZWQuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkLlxyXG4gICAgICovXHJcbiAgICBJc0ZvY3VzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0ZvY3VzZWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZnVsbHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW4uXHJcbiAgICAgKi9cclxuICAgIElzRnVsbHNjcmVlbigpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtYXhpbWlzZWQuXHJcbiAgICAgKi9cclxuICAgIElzTWF4aW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNYXhpbWlzZWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkLlxyXG4gICAgICovXHJcbiAgICBJc01pbmltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWluaW1pc2VkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE1heGltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBNaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1pbmltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBOYW1lKCk6IFByb21pc2U8c3RyaW5nPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShOYW1lTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxyXG4gICAgICovXHJcbiAgICBPcGVuRGV2VG9vbHMoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShPcGVuRGV2VG9vbHNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFJlbGF0aXZlUG9zaXRpb24oKTogUHJvbWlzZTxQb3NpdGlvbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZWxvYWRzIHRoZSBwYWdlIGFzc2V0cy5cclxuICAgICAqL1xyXG4gICAgUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsb2FkTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IHJlc2l6YWJsZS5cclxuICAgICAqL1xyXG4gICAgUmVzaXphYmxlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzaXphYmxlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlc3RvcmVzIHRoZSB3aW5kb3cgdG8gaXRzIHByZXZpb3VzIHN0YXRlIGlmIGl0IHdhcyBwcmV2aW91c2x5IG1pbmltaXNlZCwgbWF4aW1pc2VkIG9yIGZ1bGxzY3JlZW4uXHJcbiAgICAgKi9cclxuICAgIFJlc3RvcmUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZXN0b3JlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcC5cclxuICAgICAqL1xyXG4gICAgU2V0QWx3YXlzT25Ub3AoYWx3YXlzT25Ub3A6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gciAtIFRoZSBkZXNpcmVkIHJlZCBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICogQHBhcmFtIGcgLSBUaGUgZGVzaXJlZCBncmVlbiBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxyXG4gICAgICogQHBhcmFtIGIgLSBUaGUgZGVzaXJlZCBibHVlIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXHJcbiAgICAgKiBAcGFyYW0gYSAtIFRoZSBkZXNpcmVkIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXHJcbiAgICAgKi9cclxuICAgIFNldEJhY2tncm91bmRDb2xvdXIocjogbnVtYmVyLCBnOiBudW1iZXIsIGI6IG51bWJlciwgYTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZW1vdmVzIHRoZSB3aW5kb3cgZnJhbWUgYW5kIHRpdGxlIGJhci5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3MuXHJcbiAgICAgKi9cclxuICAgIFNldEZyYW1lbGVzcyhmcmFtZWxlc3M6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEaXNhYmxlcyB0aGUgc3lzdGVtIGZ1bGxzY3JlZW4gYnV0dG9uLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWQuXHJcbiAgICAgKi9cclxuICAgIFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkKGVuYWJsZWQ6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0TWF4U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWF4U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgbWluaW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWluaW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldE1pblNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1pblNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0UmVsYXRpdmVQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cclxuICAgICAqL1xyXG4gICAgU2V0UmVzaXphYmxlKHJlc2l6YWJsZTogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSB0aXRsZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB0aXRsZSAtIFRoZSBkZXNpcmVkIHRpdGxlIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFRpdGxlKHRpdGxlOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFRpdGxlTWV0aG9kLCB7IHRpdGxlIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB6b29tIC0gVGhlIGRlc2lyZWQgem9vbSBsZXZlbC5cclxuICAgICAqL1xyXG4gICAgU2V0Wm9vbSh6b29tOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFpvb21NZXRob2QsIHsgem9vbSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNob3dzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTaG93TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBzaXplIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNpemUoKTogUHJvbWlzZTxTaXplPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTaXplTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZ1bGxzY3JlZW4gYW5kIG5vcm1hbC5cclxuICAgICAqL1xyXG4gICAgVG9nZ2xlRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gbWF4aW1pc2VkIGFuZCBub3JtYWwuXHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZU1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnJhbWVsZXNzIGFuZCBub3JtYWwuXHJcbiAgICAgKi9cclxuICAgIFRvZ2dsZUZyYW1lbGVzcygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZyYW1lbGVzc01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBVbkZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFVuTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1heGltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFVuLW1pbmltaXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBVbk1pbmltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NaW5pbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFdpZHRoKCk6IFByb21pc2U8bnVtYmVyPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShXaWR0aE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBab29tcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBab29tKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBJbmNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cclxuICAgICAqL1xyXG4gICAgWm9vbUluKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbUluTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIERlY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICovXHJcbiAgICBab29tT3V0KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU91dE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cclxuICAgICAqL1xyXG4gICAgWm9vbVJlc2V0KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbVJlc2V0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEhhbmRsZXMgZmlsZSBkcm9wcyBvcmlnaW5hdGluZyBmcm9tIHBsYXRmb3JtLXNwZWNpZmljIGNvZGUgKGUuZy4sIG1hY09TL0xpbnV4IG5hdGl2ZSBkcmFnLWFuZC1kcm9wKS5cclxuICAgICAqIEdhdGhlcnMgaW5mb3JtYXRpb24gYWJvdXQgdGhlIGRyb3AgdGFyZ2V0IGVsZW1lbnQgYW5kIHNlbmRzIGl0IGJhY2sgdG8gdGhlIEdvIGJhY2tlbmQuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGZpbGVuYW1lcyAtIEFuIGFycmF5IG9mIGZpbGUgcGF0aHMgKHN0cmluZ3MpIHRoYXQgd2VyZSBkcm9wcGVkLlxyXG4gICAgICogQHBhcmFtIHggLSBUaGUgeC1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50IChDU1MgcGl4ZWxzKS5cclxuICAgICAqIEBwYXJhbSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudCAoQ1NTIHBpeGVscykuXHJcbiAgICAgKi9cclxuICAgIEhhbmRsZVBsYXRmb3JtRmlsZURyb3AoZmlsZW5hbWVzOiBzdHJpbmdbXSwgeDogbnVtYmVyLCB5OiBudW1iZXIpOiB2b2lkIHtcclxuICAgICAgICAvLyBDaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xyXG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgaWdub3JlIHRoZSBkcm9wXHJcbiAgICAgICAgfVxyXG4gICAgICAgIFxyXG4gICAgICAgIGNvbnN0IGVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHgsIHkpO1xyXG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudChlbGVtZW50KTtcclxuXHJcbiAgICAgICAgaWYgKCFkcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgICAgIC8vIERyb3Agd2FzIG5vdCBvbiBhIGRlc2lnbmF0ZWQgZHJvcCB0YXJnZXQgLSBpZ25vcmVcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgY29uc3QgZWxlbWVudERldGFpbHMgPSB7XHJcbiAgICAgICAgICAgIGlkOiBkcm9wVGFyZ2V0LmlkLFxyXG4gICAgICAgICAgICBjbGFzc0xpc3Q6IEFycmF5LmZyb20oZHJvcFRhcmdldC5jbGFzc0xpc3QpLFxyXG4gICAgICAgICAgICBhdHRyaWJ1dGVzOiB7fSBhcyB7IFtrZXk6IHN0cmluZ106IHN0cmluZyB9LFxyXG4gICAgICAgIH07XHJcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXMubGVuZ3RoOyBpKyspIHtcclxuICAgICAgICAgICAgY29uc3QgYXR0ciA9IGRyb3BUYXJnZXQuYXR0cmlidXRlc1tpXTtcclxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMuYXR0cmlidXRlc1thdHRyLm5hbWVdID0gYXR0ci52YWx1ZTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGNvbnN0IHBheWxvYWQgPSB7XHJcbiAgICAgICAgICAgIGZpbGVuYW1lcyxcclxuICAgICAgICAgICAgeCxcclxuICAgICAgICAgICAgeSxcclxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMsXHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dKEZpbGVzRHJvcHBlZCwgcGF5bG9hZCk7XHJcbiAgICAgICAgXHJcbiAgICAgICAgLy8gQ2xlYW4gdXAgbmF0aXZlIGRyYWcgc3RhdGUgYWZ0ZXIgZHJvcFxyXG4gICAgICAgIGNsZWFudXBOYXRpdmVEcmFnKCk7XHJcbiAgICB9XHJcbiAgXHJcbiAgICAvKipcclxuICAgICAqIE1vdmVzIHRoZSB3aW5kb3cgdG8gdGhlIGNlbnRlciBvZiB0aGUgc3BlY2lmaWVkIHNjcmVlbidzIHdvcmsgYXJlYS5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gc2NyZWVuSUQgLSBUaGUgSUQgb2YgdGhlIHRhcmdldCBzY3JlZW4uXHJcbiAgICAgKi9cclxuICAgIFNldFNjcmVlbihzY3JlZW5JRDogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTY3JlZW5NZXRob2QsIHsgc2NyZWVuSUQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyogVHJpZ2dlcnMgV2luZG93cyAxMSBTbmFwIEFzc2lzdCBmZWF0dXJlIChXaW5kb3dzIG9ubHkpLlxyXG4gICAgICogVGhpcyBpcyBlcXVpdmFsZW50IHRvIHByZXNzaW5nIFdpbitaIGFuZCBzaG93cyBzbmFwIGxheW91dCBvcHRpb25zLlxyXG4gICAgICovXHJcbiAgICBTbmFwQXNzaXN0KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU25hcEFzc2lzdE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBPcGVucyB0aGUgcHJpbnQgZGlhbG9nIGZvciB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBQcmludCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFByaW50TWV0aG9kKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFRoZSB3aW5kb3cgd2l0aGluIHdoaWNoIHRoZSBzY3JpcHQgaXMgcnVubmluZy5cclxuICovXHJcbmNvbnN0IHRoaXNXaW5kb3cgPSBuZXcgV2luZG93KCcnKTtcclxuXHJcbi8qKlxyXG4gKiBTZXRzIHVwIGdsb2JhbCBkcmFnIGFuZCBkcm9wIGV2ZW50IGxpc3RlbmVycyBmb3IgZmlsZSBkcm9wcy5cclxuICogSGFuZGxlcyB2aXN1YWwgZmVlZGJhY2sgKGhvdmVyIHN0YXRlKSBhbmQgZmlsZSBkcm9wIHByb2Nlc3NpbmcuXHJcbiAqL1xyXG5mdW5jdGlvbiBzZXR1cERyb3BUYXJnZXRMaXN0ZW5lcnMoKSB7XHJcbiAgICBjb25zdCBkb2NFbGVtZW50ID0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50O1xyXG4gICAgbGV0IGRyYWdFbnRlckNvdW50ZXIgPSAwO1xyXG5cclxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2VudGVyJywgKGV2ZW50KSA9PiB7XHJcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gT25seSBoYW5kbGUgZmlsZSBkcmFncywgbGV0IG90aGVyIGRyYWdzIHBhc3MgdGhyb3VnaFxyXG4gICAgICAgIH1cclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJzsgLy8gU2hvdyBcIm5vIGRyb3BcIiBjdXJzb3JcclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcclxuICAgICAgICB9XHJcbiAgICAgICAgZHJhZ0VudGVyQ291bnRlcisrO1xyXG4gICAgICAgIFxyXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xyXG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudCh0YXJnZXRFbGVtZW50KTtcclxuXHJcbiAgICAgICAgLy8gVXBkYXRlIGhvdmVyIHN0YXRlXHJcbiAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0ICYmIGN1cnJlbnREcm9wVGFyZ2V0ICE9PSBkcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmIChkcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdjb3B5JztcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBkcm9wVGFyZ2V0O1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnO1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG5cclxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCAoZXZlbnQpID0+IHtcclxuICAgICAgICBpZiAoIWV2ZW50LmRhdGFUcmFuc2Zlcj8udHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyYWdzXHJcbiAgICAgICAgfVxyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIEFsd2F5cyBwcmV2ZW50IGRlZmF1bHQgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cclxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xyXG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XHJcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnOyAvLyBTaG93IFwibm8gZHJvcFwiIGN1cnNvclxyXG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xyXG4gICAgICAgIH1cclxuICAgICAgICBcclxuICAgICAgICAvLyBVcGRhdGUgZHJvcCB0YXJnZXQgYXMgY3Vyc29yIG1vdmVzXHJcbiAgICAgICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSk7XHJcbiAgICAgICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xyXG4gICAgICAgIFxyXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIFxyXG4gICAgICAgIGlmIChkcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgICAgIGlmICghZHJvcFRhcmdldC5jbGFzc0xpc3QuY29udGFpbnMoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKSkge1xyXG4gICAgICAgICAgICAgICAgZHJvcFRhcmdldC5jbGFzc0xpc3QuYWRkKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgIH0sIGZhbHNlKTtcclxuXHJcbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdsZWF2ZScsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIEFsd2F5cyBwcmV2ZW50IGRlZmF1bHQgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cclxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xyXG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICAgICAgXHJcbiAgICAgICAgLy8gT24gTGludXgvV2ViS2l0R1RLIGFuZCBtYWNPUywgZHJhZ2xlYXZlIGZpcmVzIGltbWVkaWF0ZWx5IHdpdGggcmVsYXRlZFRhcmdldD1udWxsIHdoZW4gbmF0aXZlXHJcbiAgICAgICAgLy8gZHJhZyBoYW5kbGluZyBpcyBpbnZvbHZlZC4gSWdub3JlIHRoZXNlIHNwdXJpb3VzIGV2ZW50cyAtIHdlJ2xsIGNsZWFuIHVwIG9uIGRyb3AgaW5zdGVhZC5cclxuICAgICAgICBpZiAoZXZlbnQucmVsYXRlZFRhcmdldCA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIFxyXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXItLTtcclxuICAgICAgICBcclxuICAgICAgICBpZiAoZHJhZ0VudGVyQ291bnRlciA9PT0gMCB8fCBcclxuICAgICAgICAgICAgKGN1cnJlbnREcm9wVGFyZ2V0ICYmICFjdXJyZW50RHJvcFRhcmdldC5jb250YWlucyhldmVudC5yZWxhdGVkVGFyZ2V0IGFzIE5vZGUpKSkge1xyXG4gICAgICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcclxuICAgICAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgICAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBkcmFnRW50ZXJDb3VudGVyID0gMDtcclxuICAgICAgICB9XHJcbiAgICB9LCBmYWxzZSk7XHJcblxyXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgKGV2ZW50KSA9PiB7XHJcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gT25seSBoYW5kbGUgZmlsZSBkcm9wc1xyXG4gICAgICAgIH1cclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwO1xyXG4gICAgICAgIFxyXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGhhbmRsZSBmaWxlIGRyb3BzIHZpYSBKYXZhU2NyaXB0XHJcbiAgICAgICAgLy8gT24gbWFjT1MvTGludXgsIG5hdGl2ZSBjb2RlIHdpbGwgY2FsbCBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wXHJcbiAgICAgICAgaWYgKGNhblJlc29sdmVGaWxlUGF0aHMoKSkge1xyXG4gICAgICAgICAgICBjb25zdCBmaWxlczogRmlsZVtdID0gW107XHJcbiAgICAgICAgICAgIGlmIChldmVudC5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcclxuICAgICAgICAgICAgICAgIGZvciAoY29uc3QgaXRlbSBvZiBldmVudC5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgZmlsZSA9IGl0ZW0uZ2V0QXNGaWxlKCk7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChmaWxlKSBmaWxlcy5wdXNoKGZpbGUpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSBlbHNlIGlmIChldmVudC5kYXRhVHJhbnNmZXIuZmlsZXMpIHtcclxuICAgICAgICAgICAgICAgIGZvciAoY29uc3QgZmlsZSBvZiBldmVudC5kYXRhVHJhbnNmZXIuZmlsZXMpIHtcclxuICAgICAgICAgICAgICAgICAgICBmaWxlcy5wdXNoKGZpbGUpO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIFxyXG4gICAgICAgICAgICBpZiAoZmlsZXMubGVuZ3RoID4gMCkge1xyXG4gICAgICAgICAgICAgICAgcmVzb2x2ZUZpbGVQYXRocyhldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBmaWxlcyk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9LCBmYWxzZSk7XHJcbn1cclxuXHJcbi8vIEluaXRpYWxpemUgbGlzdGVuZXJzIHdoZW4gdGhlIHNjcmlwdCBsb2Fkc1xyXG5pZiAodHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgZG9jdW1lbnQgIT09IFwidW5kZWZpbmVkXCIpIHtcclxuICAgIHNldHVwRHJvcFRhcmdldExpc3RlbmVycygpO1xyXG59XHJcblxyXG5leHBvcnQgZGVmYXVsdCB0aGlzV2luZG93O1xyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0ICogYXMgUnVudGltZSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmNcIjtcclxuXHJcbi8vIE5PVEU6IHRoZSBmb2xsb3dpbmcgbWV0aG9kcyBNVVNUIGJlIGltcG9ydGVkIGV4cGxpY2l0bHkgYmVjYXVzZSBvZiBob3cgZXNidWlsZCBpbmplY3Rpb24gd29ya3NcclxuaW1wb3J0IHsgRW5hYmxlIGFzIEVuYWJsZVdNTCB9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWxcIjtcclxuaW1wb3J0IHsgZGVidWdMb2cgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHNcIjtcclxuXHJcbndpbmRvdy53YWlscyA9IFJ1bnRpbWU7XHJcbkVuYWJsZVdNTCgpO1xyXG5cclxuaWYgKERFQlVHKSB7XHJcbiAgICBkZWJ1Z0xvZyhcIldhaWxzIFJ1bnRpbWUgTG9hZGVkXCIpXHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuU3lzdGVtKTtcclxuXHJcbmNvbnN0IFN5c3RlbUlzRGFya01vZGUgPSAwO1xyXG5jb25zdCBTeXN0ZW1FbnZpcm9ubWVudCA9IDE7XHJcbmNvbnN0IFN5c3RlbUNhcGFiaWxpdGllcyA9IDI7XHJcblxyXG5jb25zdCBfaW52b2tlID0gKGZ1bmN0aW9uICgpIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgLy8gV2luZG93cyBXZWJWaWV3MlxyXG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZSkge1xyXG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3KTtcclxuICAgICAgICB9XHJcbiAgICAgICAgLy8gbWFjT1MvaU9TIFdLV2ViVmlld1xyXG4gICAgICAgIGVsc2UgaWYgKCh3aW5kb3cgYXMgYW55KS53ZWJraXQ/Lm1lc3NhZ2VIYW5kbGVycz8uWydleHRlcm5hbCddPy5wb3N0TWVzc2FnZSkge1xyXG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10ucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIC8vIEFuZHJvaWQgV2ViVmlldyAtIHVzZXMgYWRkSmF2YXNjcmlwdEludGVyZmFjZSB3aGljaCBleHBvc2VzIHdpbmRvdy53YWlscy5pbnZva2VcclxuICAgICAgICBlbHNlIGlmICgod2luZG93IGFzIGFueSkud2FpbHM/Lmludm9rZSkge1xyXG4gICAgICAgICAgICByZXR1cm4gKG1zZzogYW55KSA9PiAod2luZG93IGFzIGFueSkud2FpbHMuaW52b2tlKHR5cGVvZiBtc2cgPT09ICdzdHJpbmcnID8gbXNnIDogSlNPTi5zdHJpbmdpZnkobXNnKSk7XHJcbiAgICAgICAgfVxyXG4gICAgfSBjYXRjaChlKSB7fVxyXG5cclxuICAgIGNvbnNvbGUud2FybignXFxuJWNcdTI2QTBcdUZFMEYgQnJvd3NlciBFbnZpcm9ubWVudCBEZXRlY3RlZCAlY1xcblxcbiVjT25seSBVSSBwcmV2aWV3cyBhcmUgYXZhaWxhYmxlIGluIHRoZSBicm93c2VyLiBGb3IgZnVsbCBmdW5jdGlvbmFsaXR5LCBwbGVhc2UgcnVuIHRoZSBhcHBsaWNhdGlvbiBpbiBkZXNrdG9wIG1vZGUuXFxuTW9yZSBpbmZvcm1hdGlvbiBhdDogaHR0cHM6Ly92My53YWlscy5pby9sZWFybi9idWlsZC8jdXNpbmctYS1icm93c2VyLWZvci1kZXZlbG9wbWVudFxcbicsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNmZmZmZmY7IGNvbG9yOiAjMDAwMDAwOyBmb250LXdlaWdodDogYm9sZDsgcGFkZGluZzogNHB4IDhweDsgYm9yZGVyLXJhZGl1czogNHB4OyBib3JkZXI6IDJweCBzb2xpZCAjMDAwMDAwOycsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6IHRyYW5zcGFyZW50OycsXHJcbiAgICAgICAgJ2NvbG9yOiAjZmZmZmZmOyBmb250LXN0eWxlOiBpdGFsaWM7IGZvbnQtd2VpZ2h0OiBib2xkOycpO1xyXG4gICAgcmV0dXJuIG51bGw7XHJcbn0pKCk7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZzogYW55KTogdm9pZCB7XHJcbiAgICBfaW52b2tlPy4obXNnKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGEgYm9vbGVhbiB2YWx1ZSBpbmRpY2F0aW5nIGlmIHRoZSBzeXN0ZW0gaXMgaW4gZGFyayBtb2RlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1Jc0RhcmtNb2RlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyB0aGUgY2FwYWJpbGl0aWVzLlxyXG4gKi9cclxuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIENhcGFiaWxpdGllcygpOiBQcm9taXNlPFJlY29yZDxzdHJpbmcsIGFueT4+IHtcclxuICAgIHJldHVybiBjYWxsKFN5c3RlbUNhcGFiaWxpdGllcyk7XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgT1NJbmZvIHtcclxuICAgIC8qKiBUaGUgYnJhbmRpbmcgb2YgdGhlIE9TLiAqL1xyXG4gICAgQnJhbmRpbmc6IHN0cmluZztcclxuICAgIC8qKiBUaGUgSUQgb2YgdGhlIE9TLiAqL1xyXG4gICAgSUQ6IHN0cmluZztcclxuICAgIC8qKiBUaGUgbmFtZSBvZiB0aGUgT1MuICovXHJcbiAgICBOYW1lOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIHZlcnNpb24gb2YgdGhlIE9TLiAqL1xyXG4gICAgVmVyc2lvbjogc3RyaW5nO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIEVudmlyb25tZW50SW5mbyB7XHJcbiAgICAvKiogVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLiAqL1xyXG4gICAgQXJjaDogc3RyaW5nO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGFwcGxpY2F0aW9uIGlzIHJ1bm5pbmcgaW4gZGVidWcgbW9kZSwgb3RoZXJ3aXNlIGZhbHNlLiAqL1xyXG4gICAgRGVidWc6IGJvb2xlYW47XHJcbiAgICAvKiogVGhlIG9wZXJhdGluZyBzeXN0ZW0gaW4gdXNlLiAqL1xyXG4gICAgT1M6IHN0cmluZztcclxuICAgIC8qKiBEZXRhaWxzIG9mIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xyXG4gICAgT1NJbmZvOiBPU0luZm87XHJcbiAgICAvKiogQWRkaXRpb25hbCBwbGF0Zm9ybSBpbmZvcm1hdGlvbi4gKi9cclxuICAgIFBsYXRmb3JtSW5mbzogUmVjb3JkPHN0cmluZywgYW55PjtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBvYmplY3QgY29udGFpbmluZyBPUyBhbmQgc3lzdGVtIGFyY2hpdGVjdHVyZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpOiBQcm9taXNlPEVudmlyb25tZW50SW5mbz4ge1xyXG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtRW52aXJvbm1lbnQpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cclxuICpcclxuICogQHJldHVybiBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MsIG90aGVyd2lzZSBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc1dpbmRvd3MoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TID09PSBcIndpbmRvd3NcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBSZXR1cm5zIHRydWUgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TID09PSBcImxpbnV4XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBtYWNPUywgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJkYXJ3aW5cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhbWQ2NFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXHJcbiAqXHJcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFybVwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cclxuICpcclxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUsIG90aGVyd2lzZSByZXR1cm5zIGZhbHNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYXJtNjRcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlcG9ydHMgd2hldGhlciB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiBCb29sZWFuKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5EZWJ1Zyk7XHJcbn1cclxuXHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgSXNEZWJ1ZyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xyXG5pbXBvcnQgeyBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XHJcblxyXG4vLyBzZXR1cFxyXG53aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xyXG5cclxuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcclxuXHJcbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZDogc3RyaW5nLCB4OiBudW1iZXIsIHk6IG51bWJlciwgZGF0YTogYW55KTogdm9pZCB7XHJcbiAgICB2b2lkIGNhbGwoQ29udGV4dE1lbnVPcGVuLCB7aWQsIHgsIHksIGRhdGF9KTtcclxufVxyXG5cclxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50OiBNb3VzZUV2ZW50KSB7XHJcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XHJcblxyXG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcclxuICAgIGNvbnN0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIikudHJpbSgpO1xyXG5cclxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgY29uc3QgZGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XHJcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBkYXRhKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCwgdGFyZ2V0KTtcclxuICAgIH1cclxufVxyXG5cclxuXHJcbi8qXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG5cclxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcclxuKi9cclxuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudDogTW91c2VFdmVudCwgdGFyZ2V0OiBIVE1MRWxlbWVudCkge1xyXG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XHJcbiAgICBpZiAoSXNEZWJ1ZygpKSB7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgIHN3aXRjaCAod2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKSkge1xyXG4gICAgICAgIGNhc2UgJ3Nob3cnOlxyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgY2FzZSAnaGlkZSc6XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxyXG4gICAgaWYgKHRhcmdldC5pc0NvbnRlbnRFZGl0YWJsZSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXHJcbiAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XHJcbiAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSBzZWxlY3Rpb24gJiYgc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMDtcclxuICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcclxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcclxuICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcclxuICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xyXG4gICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XHJcbiAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XHJcbiAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gdGFyZ2V0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8vIENoZWNrIGlmIHRhZyBpcyBpbnB1dCBvciB0ZXh0YXJlYS5cclxuICAgIGlmICh0YXJnZXQgaW5zdGFuY2VvZiBIVE1MSW5wdXRFbGVtZW50IHx8IHRhcmdldCBpbnN0YW5jZW9mIEhUTUxUZXh0QXJlYUVsZW1lbnQpIHtcclxuICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghdGFyZ2V0LnJlYWRPbmx5ICYmICF0YXJnZXQuZGlzYWJsZWQpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG4gICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoqXHJcbiAqIFJldHJpZXZlcyB0aGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5IGZyb20gdGhlIGZsYWcgbWFwLlxyXG4gKlxyXG4gKiBAcGFyYW0ga2V5IC0gVGhlIGtleSB0byByZXRyaWV2ZSB0aGUgdmFsdWUgZm9yLlxyXG4gKiBAcmV0dXJuIFRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXk6IHN0cmluZyk6IGFueSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleV07XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiVW5hYmxlIHRvIHJldHJpZXZlIGZsYWcgJ1wiICsga2V5ICsgXCInOiBcIiArIGUsIHsgY2F1c2U6IGUgfSk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IGludm9rZSwgSXNXaW5kb3dzIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XHJcbmltcG9ydCB7IEdldEZsYWcgfSBmcm9tIFwiLi9mbGFncy5qc1wiO1xyXG5pbXBvcnQgeyBjYW5UcmFja0J1dHRvbnMsIGV2ZW50VGFyZ2V0IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcclxuXHJcbi8vIFNldHVwXHJcbmxldCBjYW5EcmFnID0gZmFsc2U7XHJcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xyXG5cclxubGV0IHJlc2l6YWJsZSA9IGZhbHNlO1xyXG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XHJcbmxldCByZXNpemluZyA9IGZhbHNlO1xyXG5sZXQgcmVzaXplRWRnZTogc3RyaW5nID0gXCJcIjtcclxubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcclxuXHJcbmxldCBidXR0b25zID0gMDtcclxuY29uc3QgYnV0dG9uc1RyYWNrZWQgPSBjYW5UcmFja0J1dHRvbnMoKTtcclxuXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9ICh2YWx1ZTogYm9vbGVhbik6IHZvaWQgPT4ge1xyXG4gICAgcmVzaXphYmxlID0gdmFsdWU7XHJcbiAgICBpZiAoIXJlc2l6YWJsZSkge1xyXG4gICAgICAgIC8vIFN0b3AgcmVzaXppbmcgaWYgaW4gcHJvZ3Jlc3MuXHJcbiAgICAgICAgY2FuUmVzaXplID0gcmVzaXppbmcgPSBmYWxzZTtcclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH1cclxufTtcclxuXHJcbi8vIERlZmVyIGF0dGFjaGluZyBtb3VzZSBsaXN0ZW5lcnMgdW50aWwgd2Uga25vdyB3ZSdyZSBub3Qgb24gbW9iaWxlLlxyXG5sZXQgZHJhZ0luaXREb25lID0gZmFsc2U7XHJcbmZ1bmN0aW9uIGlzTW9iaWxlKCk6IGJvb2xlYW4ge1xyXG4gICAgY29uc3Qgb3MgPSAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1M7XHJcbiAgICBpZiAob3MgPT09IFwiaW9zXCIgfHwgb3MgPT09IFwiYW5kcm9pZFwiKSByZXR1cm4gdHJ1ZTtcclxuICAgIC8vIEZhbGxiYWNrIGhldXJpc3RpYyBpZiBlbnZpcm9ubWVudCBub3QgeWV0IHNldFxyXG4gICAgY29uc3QgdWEgPSBuYXZpZ2F0b3IudXNlckFnZW50IHx8IG5hdmlnYXRvci52ZW5kb3IgfHwgKHdpbmRvdyBhcyBhbnkpLm9wZXJhIHx8IFwiXCI7XHJcbiAgICByZXR1cm4gL2FuZHJvaWR8aXBob25lfGlwYWR8aXBvZHxpZW1vYmlsZXx3cGRlc2t0b3AvaS50ZXN0KHVhKTtcclxufVxyXG5mdW5jdGlvbiB0cnlJbml0RHJhZ0hhbmRsZXJzKCk6IHZvaWQge1xyXG4gICAgaWYgKGRyYWdJbml0RG9uZSkgcmV0dXJuO1xyXG4gICAgaWYgKGlzTW9iaWxlKCkpIHJldHVybjtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbiAgICBmb3IgKGNvbnN0IGV2IG9mIFsnY2xpY2snLCAnY29udGV4dG1lbnUnLCAnZGJsY2xpY2snXSkge1xyXG4gICAgICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKGV2LCBzdXBwcmVzc0V2ZW50LCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbiAgICB9XHJcbiAgICBkcmFnSW5pdERvbmUgPSB0cnVlO1xyXG59XHJcbi8vIEF0dGVtcHQgaW1tZWRpYXRlIGluaXQgKGluIGNhc2UgZW52aXJvbm1lbnQgYWxyZWFkeSBwcmVzZW50KVxyXG50cnlJbml0RHJhZ0hhbmRsZXJzKCk7XHJcbi8vIEFsc28gYXR0ZW1wdCBvbiBET00gcmVhZHlcclxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIHRyeUluaXREcmFnSGFuZGxlcnMsIHsgb25jZTogdHJ1ZSB9KTtcclxuLy8gQXMgYSBsYXN0IHJlc29ydCwgcG9sbCBmb3IgZW52aXJvbm1lbnQgZm9yIGEgc2hvcnQgcGVyaW9kXHJcbmxldCBkcmFnRW52UG9sbHMgPSAwO1xyXG5jb25zdCBkcmFnRW52UG9sbCA9IHdpbmRvdy5zZXRJbnRlcnZhbCgoKSA9PiB7XHJcbiAgICBpZiAoZHJhZ0luaXREb25lKSB7IHdpbmRvdy5jbGVhckludGVydmFsKGRyYWdFbnZQb2xsKTsgcmV0dXJuOyB9XHJcbiAgICB0cnlJbml0RHJhZ0hhbmRsZXJzKCk7XHJcbiAgICBpZiAoKytkcmFnRW52UG9sbHMgPiAxMDApIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyB9XHJcbn0sIDUwKTtcclxuXHJcbmZ1bmN0aW9uIHN1cHByZXNzRXZlbnQoZXZlbnQ6IEV2ZW50KSB7XHJcbiAgICAvLyBTdXBwcmVzcyBjbGljayBldmVudHMgd2hpbGUgcmVzaXppbmcgb3IgZHJhZ2dpbmcuXHJcbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcclxuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vLyBVc2UgY29uc3RhbnRzIHRvIGF2b2lkIGNvbXBhcmluZyBzdHJpbmdzIG11bHRpcGxlIHRpbWVzLlxyXG5jb25zdCBNb3VzZURvd24gPSAwO1xyXG5jb25zdCBNb3VzZVVwICAgPSAxO1xyXG5jb25zdCBNb3VzZU1vdmUgPSAyO1xyXG5cclxuZnVuY3Rpb24gdXBkYXRlKGV2ZW50OiBNb3VzZUV2ZW50KSB7XHJcbiAgICAvLyBXaW5kb3dzIHN1cHByZXNzZXMgbW91c2UgZXZlbnRzIGF0IHRoZSBlbmQgb2YgZHJhZ2dpbmcgb3IgcmVzaXppbmcsXHJcbiAgICAvLyBzbyB3ZSBuZWVkIHRvIGJlIHNtYXJ0IGFuZCBzeW50aGVzaXplIGJ1dHRvbiBldmVudHMuXHJcblxyXG4gICAgbGV0IGV2ZW50VHlwZTogbnVtYmVyLCBldmVudEJ1dHRvbnMgPSBldmVudC5idXR0b25zO1xyXG4gICAgc3dpdGNoIChldmVudC50eXBlKSB7XHJcbiAgICAgICAgY2FzZSAnbW91c2Vkb3duJzpcclxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VEb3duO1xyXG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgfCAoMSA8PCBldmVudC5idXR0b24pOyB9XHJcbiAgICAgICAgICAgIGJyZWFrO1xyXG4gICAgICAgIGNhc2UgJ21vdXNldXAnOlxyXG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZVVwO1xyXG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgJiB+KDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICBkZWZhdWx0OlxyXG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZU1vdmU7XHJcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9uczsgfVxyXG4gICAgICAgICAgICBicmVhaztcclxuICAgIH1cclxuXHJcbiAgICBsZXQgcmVsZWFzZWQgPSBidXR0b25zICYgfmV2ZW50QnV0dG9ucztcclxuICAgIGxldCBwcmVzc2VkID0gZXZlbnRCdXR0b25zICYgfmJ1dHRvbnM7XHJcblxyXG4gICAgYnV0dG9ucyA9IGV2ZW50QnV0dG9ucztcclxuXHJcbiAgICAvLyBTeW50aGVzaXplIGEgcmVsZWFzZS1wcmVzcyBzZXF1ZW5jZSBpZiB3ZSBkZXRlY3QgYSBwcmVzcyBvZiBhbiBhbHJlYWR5IHByZXNzZWQgYnV0dG9uLlxyXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VEb3duICYmICEocHJlc3NlZCAmIGV2ZW50LmJ1dHRvbikpIHtcclxuICAgICAgICByZWxlYXNlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xyXG4gICAgICAgIHByZXNzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBTdXBwcmVzcyBhbGwgYnV0dG9uIGV2ZW50cyBkdXJpbmcgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLFxyXG4gICAgLy8gdW5sZXNzIHRoaXMgaXMgYSBtb3VzZXVwIGV2ZW50IHRoYXQgaXMgZW5kaW5nIGEgZHJhZyBhY3Rpb24uXHJcbiAgICBpZiAoXHJcbiAgICAgICAgZXZlbnRUeXBlICE9PSBNb3VzZU1vdmUgLy8gRmFzdCBwYXRoIGZvciBtb3VzZW1vdmVcclxuICAgICAgICAmJiByZXNpemluZ1xyXG4gICAgICAgIHx8IChcclxuICAgICAgICAgICAgZHJhZ2dpbmdcclxuICAgICAgICAgICAgJiYgKFxyXG4gICAgICAgICAgICAgICAgZXZlbnRUeXBlID09PSBNb3VzZURvd25cclxuICAgICAgICAgICAgICAgIHx8IGV2ZW50LmJ1dHRvbiAhPT0gMFxyXG4gICAgICAgICAgICApXHJcbiAgICAgICAgKVxyXG4gICAgKSB7XHJcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XHJcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBIYW5kbGUgcmVsZWFzZXNcclxuICAgIGlmIChyZWxlYXNlZCAmIDEpIHsgcHJpbWFyeVVwKGV2ZW50KTsgfVxyXG4gICAgLy8gSGFuZGxlIHByZXNzZXNcclxuICAgIGlmIChwcmVzc2VkICYgMSkgeyBwcmltYXJ5RG93bihldmVudCk7IH1cclxuXHJcbiAgICAvLyBIYW5kbGUgbW91c2Vtb3ZlXHJcbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZU1vdmUpIHsgb25Nb3VzZU1vdmUoZXZlbnQpOyB9O1xyXG59XHJcblxyXG5mdW5jdGlvbiBwcmltYXJ5RG93bihldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xyXG4gICAgLy8gUmVzZXQgcmVhZGluZXNzIHN0YXRlLlxyXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xyXG4gICAgY2FuUmVzaXplID0gZmFsc2U7XHJcblxyXG4gICAgLy8gSWdub3JlIHJlcGVhdGVkIGNsaWNrcyBvbiBtYWNPUyBhbmQgTGludXguXHJcbiAgICBpZiAoIUlzV2luZG93cygpKSB7XHJcbiAgICAgICAgaWYgKGV2ZW50LnR5cGUgPT09ICdtb3VzZWRvd24nICYmIGV2ZW50LmJ1dHRvbiA9PT0gMCAmJiBldmVudC5kZXRhaWwgIT09IDEpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICBpZiAocmVzaXplRWRnZSkge1xyXG4gICAgICAgIC8vIFJlYWR5IHRvIHJlc2l6ZSBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lLlxyXG4gICAgICAgIGNhblJlc2l6ZSA9IHRydWU7XHJcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWcgb3BlcmF0aW9ucyB3aGVuIG9uIHJlc2l6ZSBlZGdlcy5cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gUmV0cmlldmUgdGFyZ2V0IGVsZW1lbnRcclxuICAgIGNvbnN0IHRhcmdldCA9IGV2ZW50VGFyZ2V0KGV2ZW50KTtcclxuXHJcbiAgICAvLyBSZWFkeSB0byBkcmFnIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUgb24gYSBkcmFnZ2FibGUgZWxlbWVudC5cclxuICAgIC8vIElnbm9yZSBjbGlja3Mgb24gdGhlIHNjcm9sbGJhci5cclxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KTtcclxuICAgIGNhbkRyYWcgPSAoXHJcbiAgICAgICAgc3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpLnRyaW0oKSA9PT0gXCJkcmFnXCJcclxuICAgICAgICAmJiAoXHJcbiAgICAgICAgICAgIGV2ZW50Lm9mZnNldFggLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdMZWZ0KSA8IHRhcmdldC5jbGllbnRXaWR0aFxyXG4gICAgICAgICAgICAmJiBldmVudC5vZmZzZXRZIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nVG9wKSA8IHRhcmdldC5jbGllbnRIZWlnaHRcclxuICAgICAgICApXHJcbiAgICApO1xyXG59XHJcblxyXG5mdW5jdGlvbiBwcmltYXJ5VXAoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcclxuICAgIC8vIFN0b3AgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLlxyXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xyXG4gICAgZHJhZ2dpbmcgPSBmYWxzZTtcclxuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xyXG4gICAgcmVzaXppbmcgPSBmYWxzZTtcclxufVxyXG5cclxuY29uc3QgY3Vyc29yRm9yRWRnZSA9IE9iamVjdC5mcmVlemUoe1xyXG4gICAgXCJzZS1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxyXG4gICAgXCJzdy1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxyXG4gICAgXCJudy1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxyXG4gICAgXCJuZS1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxyXG4gICAgXCJ3LXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxyXG4gICAgXCJuLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxyXG4gICAgXCJzLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxyXG4gICAgXCJlLXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxyXG59KVxyXG5cclxuZnVuY3Rpb24gc2V0UmVzaXplKGVkZ2U/OiBrZXlvZiB0eXBlb2YgY3Vyc29yRm9yRWRnZSk6IHZvaWQge1xyXG4gICAgaWYgKGVkZ2UpIHtcclxuICAgICAgICBpZiAoIXJlc2l6ZUVkZ2UpIHsgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yOyB9XHJcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBjdXJzb3JGb3JFZGdlW2VkZ2VdO1xyXG4gICAgfSBlbHNlIGlmICghZWRnZSAmJiByZXNpemVFZGdlKSB7XHJcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBkZWZhdWx0Q3Vyc29yO1xyXG4gICAgfVxyXG5cclxuICAgIHJlc2l6ZUVkZ2UgPSBlZGdlIHx8IFwiXCI7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XHJcbiAgICBpZiAoY2FuUmVzaXplICYmIHJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICAvLyBTdGFydCByZXNpemluZy5cclxuICAgICAgICByZXNpemluZyA9IHRydWU7XHJcbiAgICAgICAgaW52b2tlKFwid2FpbHM6cmVzaXplOlwiICsgcmVzaXplRWRnZSk7XHJcbiAgICB9IGVsc2UgaWYgKGNhbkRyYWcpIHtcclxuICAgICAgICAvLyBTdGFydCBkcmFnZ2luZy5cclxuICAgICAgICBkcmFnZ2luZyA9IHRydWU7XHJcbiAgICAgICAgaW52b2tlKFwid2FpbHM6ZHJhZ1wiKTtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcclxuICAgICAgICAvLyBFaXRoZXIgZHJhZyBvciByZXNpemUgaXMgb25nb2luZyxcclxuICAgICAgICAvLyByZXNldCByZWFkaW5lc3MgYW5kIHN0b3AgcHJvY2Vzc2luZy5cclxuICAgICAgICBjYW5EcmFnID0gY2FuUmVzaXplID0gZmFsc2U7XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICghcmVzaXphYmxlIHx8ICFJc1dpbmRvd3MoKSkge1xyXG4gICAgICAgIGlmIChyZXNpemVFZGdlKSB7IHNldFJlc2l6ZSgpOyB9XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZUhlaWdodCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlSGVpZ2h0XCIpIHx8IDU7XHJcbiAgICBjb25zdCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcclxuXHJcbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXMuXHJcbiAgICBjb25zdCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcclxuXHJcbiAgICBjb25zdCByaWdodEJvcmRlciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgcmVzaXplSGFuZGxlV2lkdGg7XHJcbiAgICBjb25zdCBsZWZ0Qm9yZGVyID0gZXZlbnQuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgY29uc3QgdG9wQm9yZGVyID0gZXZlbnQuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcclxuICAgIGNvbnN0IGJvdHRvbUJvcmRlciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcclxuXHJcbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lciBhcmVhcy5cclxuICAgIGNvbnN0IHJpZ2h0Q29ybmVyID0gKHdpbmRvdy5vdXRlcldpZHRoIC0gZXZlbnQuY2xpZW50WCkgPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBjb25zdCBsZWZ0Q29ybmVyID0gZXZlbnQuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGNvbnN0IHRvcENvcm5lciA9IGV2ZW50LmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG4gICAgY29uc3QgYm90dG9tQ29ybmVyID0gKHdpbmRvdy5vdXRlckhlaWdodCAtIGV2ZW50LmNsaWVudFkpIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcclxuXHJcbiAgICBpZiAoIWxlZnRDb3JuZXIgJiYgIXRvcENvcm5lciAmJiAhYm90dG9tQ29ybmVyICYmICFyaWdodENvcm5lcikge1xyXG4gICAgICAgIC8vIE9wdGltaXNhdGlvbjogb3V0IG9mIGFsbCBjb3JuZXIgYXJlYXMgaW1wbGllcyBvdXQgb2YgYm9yZGVycy5cclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH1cclxuICAgIC8vIERldGVjdCBjb3JuZXJzLlxyXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgLy8gRGV0ZWN0IGJvcmRlcnMuXHJcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xyXG4gICAgLy8gT3V0IG9mIGJvcmRlciBhcmVhLlxyXG4gICAgZWxzZSBzZXRSZXNpemUoKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcclxuXHJcbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xyXG5jb25zdCBTaG93TWV0aG9kID0gMTtcclxuY29uc3QgUXVpdE1ldGhvZCA9IDI7XHJcblxyXG4vKipcclxuICogSGlkZXMgYSBjZXJ0YWluIG1ldGhvZCBieSBjYWxsaW5nIHRoZSBIaWRlTWV0aG9kIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChIaWRlTWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgQ2FuY2VsbGFibGVQcm9taXNlLCB0eXBlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnMgfSBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XHJcblxyXG4vLyBTZXR1cFxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxuXHJcbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8Q2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIiB8IFwib25jYW5jZWxsZWRcIj5cclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbGwpO1xyXG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYW5jZWxDYWxsKTtcclxuY29uc3QgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xyXG5cclxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xyXG5jb25zdCBDYW5jZWxNZXRob2QgPSAwXHJcblxyXG4vKipcclxuICogSG9sZHMgYWxsIHJlcXVpcmVkIGluZm9ybWF0aW9uIGZvciBhIGJpbmRpbmcgY2FsbC5cclxuICogTWF5IHByb3ZpZGUgZWl0aGVyIGEgbWV0aG9kIElEIG9yIGEgbWV0aG9kIG5hbWUsIGJ1dCBub3QgYm90aC5cclxuICovXHJcbmV4cG9ydCB0eXBlIENhbGxPcHRpb25zID0ge1xyXG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZElEOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZE5hbWU/OiBuZXZlcjtcclxuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cclxuICAgIGFyZ3M6IGFueVtdO1xyXG59IHwge1xyXG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZElEPzogbmV2ZXI7XHJcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cclxuICAgIG1ldGhvZE5hbWU6IHN0cmluZztcclxuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cclxuICAgIGFyZ3M6IGFueVtdO1xyXG59O1xyXG5cclxuLyoqXHJcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdGhyb3duIGluIGNhc2UgdGhlIGJvdW5kIG1ldGhvZCByZXR1cm5zIGFuIGVycm9yLlxyXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBSdW50aW1lRXJyb3IjbmFtZX0gcHJvcGVydHkgaXMgXCJSdW50aW1lRXJyb3JcIi5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBSdW50aW1lRXJyb3IgZXh0ZW5kcyBFcnJvciB7XHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgUnVudGltZUVycm9yIGluc3RhbmNlLlxyXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XHJcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJSdW50aW1lRXJyb3JcIjtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBJRCB1c2luZyB0aGUgbmFub2lkIGxpYnJhcnkuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgdW5pcXVlIElEIHRoYXQgZG9lcyBub3QgZXhpc3QgaW4gdGhlIGNhbGxSZXNwb25zZXMgc2V0LlxyXG4gKi9cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGwgYSBib3VuZCBtZXRob2QgYWNjb3JkaW5nIHRvIHRoZSBnaXZlbiBjYWxsIG9wdGlvbnMuXHJcbiAqXHJcbiAqIEluIGNhc2Ugb2YgZmFpbHVyZSwgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhbiBleGNlcHRpb25cclxuICogYW1vbmcgUmVmZXJlbmNlRXJyb3IgKHVua25vd24gbWV0aG9kKSwgVHlwZUVycm9yICh3cm9uZyBhcmd1bWVudCBjb3VudCBvciB0eXBlKSxcclxuICoge0BsaW5rIFJ1bnRpbWVFcnJvcn0gKG1ldGhvZCByZXR1cm5lZCBhbiBlcnJvciksIG9yIG90aGVyIChuZXR3b3JrIG9yIGludGVybmFsIGVycm9ycykuXHJcbiAqIFRoZSBleGNlcHRpb24gbWlnaHQgaGF2ZSBhIFwiY2F1c2VcIiBmaWVsZCB3aXRoIHRoZSB2YWx1ZSByZXR1cm5lZFxyXG4gKiBieSB0aGUgYXBwbGljYXRpb24tIG9yIHNlcnZpY2UtbGV2ZWwgZXJyb3IgbWFyc2hhbGluZyBmdW5jdGlvbnMuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gQSBtZXRob2QgY2FsbCBkZXNjcmlwdG9yLlxyXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9uczogQ2FsbE9wdGlvbnMpOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XHJcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuXHJcbiAgICBjb25zdCByZXN1bHQgPSBDYW5jZWxsYWJsZVByb21pc2Uud2l0aFJlc29sdmVyczxhbnk+KCk7XHJcbiAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlOiByZXN1bHQucmVzb2x2ZSwgcmVqZWN0OiByZXN1bHQucmVqZWN0IH0pO1xyXG5cclxuICAgIGNvbnN0IHJlcXVlc3QgPSBjYWxsKENhbGxCaW5kaW5nLCBPYmplY3QuYXNzaWduKHsgXCJjYWxsLWlkXCI6IGlkIH0sIG9wdGlvbnMpKTtcclxuICAgIGxldCBydW5uaW5nID0gdHJ1ZTtcclxuXHJcbiAgICByZXF1ZXN0LnRoZW4oKHJlcykgPT4ge1xyXG4gICAgICAgIHJ1bm5pbmcgPSBmYWxzZTtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgcmVzdWx0LnJlc29sdmUocmVzKTtcclxuICAgIH0sIChlcnIpID0+IHtcclxuICAgICAgICBydW5uaW5nID0gZmFsc2U7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIHJlc3VsdC5yZWplY3QoZXJyKTtcclxuICAgIH0pO1xyXG5cclxuICAgIGNvbnN0IGNhbmNlbCA9ICgpID0+IHtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgcmV0dXJuIGNhbmNlbENhbGwoQ2FuY2VsTWV0aG9kLCB7XCJjYWxsLWlkXCI6IGlkfSkuY2F0Y2goKGVycikgPT4ge1xyXG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igd2hpbGUgcmVxdWVzdGluZyBiaW5kaW5nIGNhbGwgY2FuY2VsbGF0aW9uOlwiLCBlcnIpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfTtcclxuXHJcbiAgICByZXN1bHQub25jYW5jZWxsZWQgPSAoKSA9PiB7XHJcbiAgICAgICAgaWYgKHJ1bm5pbmcpIHtcclxuICAgICAgICAgICAgcmV0dXJuIGNhbmNlbCgpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHJldHVybiByZXF1ZXN0LnRoZW4oY2FuY2VsKTtcclxuICAgICAgICB9XHJcbiAgICB9O1xyXG5cclxuICAgIHJldHVybiByZXN1bHQucHJvbWlzZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgYm91bmQgbWV0aG9kIGJ5IG5hbWUgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cclxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cclxuICpcclxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXHJcbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShtZXRob2ROYW1lOiBzdHJpbmcsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2ROYW1lLCBhcmdzIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIG51bWVyaWMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cclxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cclxuICpcclxuICogQHBhcmFtIG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cclxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cclxuICogQHJldHVybiBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeUlEKG1ldGhvZElEOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2RJRCwgYXJncyB9KTtcclxufVxyXG4iLCAiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vaW5zcGVjdC1qcy9pcy1jYWxsYWJsZVxyXG5cclxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXHJcbi8vXHJcbi8vIENvcHlyaWdodCAoYykgMjAxNSBKb3JkYW4gSGFyYmFuZFxyXG4vL1xyXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5XHJcbi8vIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWxcclxuLy8gaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0c1xyXG4vLyB0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsXHJcbi8vIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpc1xyXG4vLyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxyXG4vL1xyXG4vLyBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcclxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cclxuLy9cclxuLy8gVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxyXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSxcclxuLy8gRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFXHJcbi8vIEFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVJcclxuLy8gTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSxcclxuLy8gT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEVcclxuLy8gU09GVFdBUkUuXHJcblxyXG52YXIgZm5Ub1N0ciA9IEZ1bmN0aW9uLnByb3RvdHlwZS50b1N0cmluZztcclxudmFyIHJlZmxlY3RBcHBseTogdHlwZW9mIFJlZmxlY3QuYXBwbHkgfCBmYWxzZSB8IG51bGwgPSB0eXBlb2YgUmVmbGVjdCA9PT0gJ29iamVjdCcgJiYgUmVmbGVjdCAhPT0gbnVsbCAmJiBSZWZsZWN0LmFwcGx5O1xyXG52YXIgYmFkQXJyYXlMaWtlOiBhbnk7XHJcbnZhciBpc0NhbGxhYmxlTWFya2VyOiBhbnk7XHJcbmlmICh0eXBlb2YgcmVmbGVjdEFwcGx5ID09PSAnZnVuY3Rpb24nICYmIHR5cGVvZiBPYmplY3QuZGVmaW5lUHJvcGVydHkgPT09ICdmdW5jdGlvbicpIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgYmFkQXJyYXlMaWtlID0gT2JqZWN0LmRlZmluZVByb3BlcnR5KHt9LCAnbGVuZ3RoJywge1xyXG4gICAgICAgICAgICBnZXQ6IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgICAgIHRocm93IGlzQ2FsbGFibGVNYXJrZXI7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgICAgICBpc0NhbGxhYmxlTWFya2VyID0ge307XHJcbiAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXRocm93LWxpdGVyYWxcclxuICAgICAgICByZWZsZWN0QXBwbHkoZnVuY3Rpb24gKCkgeyB0aHJvdyA0MjsgfSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcclxuICAgIH0gY2F0Y2ggKF8pIHtcclxuICAgICAgICBpZiAoXyAhPT0gaXNDYWxsYWJsZU1hcmtlcikge1xyXG4gICAgICAgICAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufSBlbHNlIHtcclxuICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XHJcbn1cclxuXHJcbnZhciBjb25zdHJ1Y3RvclJlZ2V4ID0gL15cXHMqY2xhc3NcXGIvO1xyXG52YXIgaXNFUzZDbGFzc0ZuID0gZnVuY3Rpb24gaXNFUzZDbGFzc0Z1bmN0aW9uKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgdmFyIGZuU3RyID0gZm5Ub1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICByZXR1cm4gY29uc3RydWN0b3JSZWdleC50ZXN0KGZuU3RyKTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7IC8vIG5vdCBhIGZ1bmN0aW9uXHJcbiAgICB9XHJcbn07XHJcblxyXG52YXIgdHJ5RnVuY3Rpb25PYmplY3QgPSBmdW5jdGlvbiB0cnlGdW5jdGlvblRvU3RyKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XHJcbiAgICAgICAgZm5Ub1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICByZXR1cm4gdHJ1ZTtcclxuICAgIH0gY2F0Y2ggKGUpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcbn07XHJcbnZhciB0b1N0ciA9IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmc7XHJcbnZhciBvYmplY3RDbGFzcyA9ICdbb2JqZWN0IE9iamVjdF0nO1xyXG52YXIgZm5DbGFzcyA9ICdbb2JqZWN0IEZ1bmN0aW9uXSc7XHJcbnZhciBnZW5DbGFzcyA9ICdbb2JqZWN0IEdlbmVyYXRvckZ1bmN0aW9uXSc7XHJcbnZhciBkZGFDbGFzcyA9ICdbb2JqZWN0IEhUTUxBbGxDb2xsZWN0aW9uXSc7IC8vIElFIDExXHJcbnZhciBkZGFDbGFzczIgPSAnW29iamVjdCBIVE1MIGRvY3VtZW50LmFsbCBjbGFzc10nO1xyXG52YXIgZGRhQ2xhc3MzID0gJ1tvYmplY3QgSFRNTENvbGxlY3Rpb25dJzsgLy8gSUUgOS0xMFxyXG52YXIgaGFzVG9TdHJpbmdUYWcgPSB0eXBlb2YgU3ltYm9sID09PSAnZnVuY3Rpb24nICYmICEhU3ltYm9sLnRvU3RyaW5nVGFnOyAvLyBiZXR0ZXI6IHVzZSBgaGFzLXRvc3RyaW5ndGFnYFxyXG5cclxudmFyIGlzSUU2OCA9ICEoMCBpbiBbLF0pOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIG5vLXNwYXJzZS1hcnJheXMsIGNvbW1hLXNwYWNpbmdcclxuXHJcbnZhciBpc0REQTogKHZhbHVlOiBhbnkpID0+IGJvb2xlYW4gPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKCkgeyByZXR1cm4gZmFsc2U7IH07XHJcbmlmICh0eXBlb2YgZG9jdW1lbnQgPT09ICdvYmplY3QnKSB7XHJcbiAgICAvLyBGaXJlZm94IDMgY2Fub25pY2FsaXplcyBEREEgdG8gdW5kZWZpbmVkIHdoZW4gaXQncyBub3QgYWNjZXNzZWQgZGlyZWN0bHlcclxuICAgIHZhciBhbGwgPSBkb2N1bWVudC5hbGw7XHJcbiAgICBpZiAodG9TdHIuY2FsbChhbGwpID09PSB0b1N0ci5jYWxsKGRvY3VtZW50LmFsbCkpIHtcclxuICAgICAgICBpc0REQSA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwodmFsdWUpIHtcclxuICAgICAgICAgICAgLyogZ2xvYmFscyBkb2N1bWVudDogZmFsc2UgKi9cclxuICAgICAgICAgICAgLy8gaW4gSUUgNi04LCB0eXBlb2YgZG9jdW1lbnQuYWxsIGlzIFwib2JqZWN0XCIgYW5kIGl0J3MgdHJ1dGh5XHJcbiAgICAgICAgICAgIGlmICgoaXNJRTY4IHx8ICF2YWx1ZSkgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3VuZGVmaW5lZCcgfHwgdHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JykpIHtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgdmFyIHN0ciA9IHRvU3RyLmNhbGwodmFsdWUpO1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybiAoXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHN0ciA9PT0gZGRhQ2xhc3NcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczJcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczMgLy8gb3BlcmEgMTIuMTZcclxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBvYmplY3RDbGFzcyAvLyBJRSA2LThcclxuICAgICAgICAgICAgICAgICAgICApICYmIHZhbHVlKCcnKSA9PSBudWxsOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIGVxZXFlcVxyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZSkgeyAvKiovIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gaXNDYWxsYWJsZVJlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55ICB7XHJcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XHJcbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgdHJ5IHtcclxuICAgICAgICAocmVmbGVjdEFwcGx5IGFzIGFueSkodmFsdWUsIG51bGwsIGJhZEFycmF5TGlrZSk7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgaWYgKGUgIT09IGlzQ2FsbGFibGVNYXJrZXIpIHsgcmV0dXJuIGZhbHNlOyB9XHJcbiAgICB9XHJcbiAgICByZXR1cm4gIWlzRVM2Q2xhc3NGbih2YWx1ZSkgJiYgdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBpc0NhbGxhYmxlTm9SZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSB7XHJcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XHJcbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgaWYgKGhhc1RvU3RyaW5nVGFnKSB7IHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7IH1cclxuICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgdmFyIHN0ckNsYXNzID0gdG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICBpZiAoc3RyQ2xhc3MgIT09IGZuQ2xhc3MgJiYgc3RyQ2xhc3MgIT09IGdlbkNsYXNzICYmICEoL15cXFtvYmplY3QgSFRNTC8pLnRlc3Qoc3RyQ2xhc3MpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcclxufTtcclxuXHJcbmV4cG9ydCBkZWZhdWx0IHJlZmxlY3RBcHBseSA/IGlzQ2FsbGFibGVSZWZBcHBseSA6IGlzQ2FsbGFibGVOb1JlZkFwcGx5O1xyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IGlzQ2FsbGFibGUgZnJvbSBcIi4vY2FsbGFibGUuanNcIjtcclxuXHJcbi8qKlxyXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHVzZWQgYXMgcmVqZWN0aW9uIHJlYXNvblxyXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gaXMgY2FuY2VsbGVkIHN1Y2Nlc3NmdWxseS5cclxuICpcclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxFcnJvclwiYC5cclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSBjYXVzZSBwYXNzZWQgdG8gdGhlIGNhbmNlbCBtZXRob2QsIGlmIGFueS5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBDYW5jZWxFcnJvciBleHRlbmRzIEVycm9yIHtcclxuICAgIC8qKlxyXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsRXJyb3JgIGluc3RhbmNlLlxyXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XHJcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxFcnJvclwiO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uXHJcbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSByZWplY3RzIGFmdGVyIGJlaW5nIGNhbmNlbGxlZCxcclxuICogb3Igd2hlbiB0aGUgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB0aHJvd3Mgb3IgcmVqZWN0cy5cclxuICpcclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiYC5cclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSByZWFzb24gdGhlIHByb21pc2UgcmVqZWN0ZWQgd2l0aC5cclxuICpcclxuICogQmVjYXVzZSB0aGUgb3JpZ2luYWwgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLFxyXG4gKiBhIHdyYXBwZXIgcHJvbWlzZSB3aWxsIGJlIHBhc3NlZCB0byB0aGUgdW5oYW5kbGVkIHJlamVjdGlvbiBsaXN0ZW5lciBpbnN0ZWFkLlxyXG4gKiBUaGUge0BsaW5rIHByb21pc2V9IHByb3BlcnR5IGhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBvcmlnaW5hbCBwcm9taXNlLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIENhbmNlbGxlZFJlamVjdGlvbkVycm9yIGV4dGVuZHMgRXJyb3Ige1xyXG4gICAgLyoqXHJcbiAgICAgKiBIb2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgcHJvbWlzZSB0aGF0IHdhcyBjYW5jZWxsZWQgYW5kIHRoZW4gcmVqZWN0ZWQuXHJcbiAgICAgKi9cclxuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPjtcclxuXHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbGxlZFJlamVjdGlvbkVycm9yYCBpbnN0YW5jZS5cclxuICAgICAqIEBwYXJhbSBwcm9taXNlIC0gVGhlIHByb21pc2UgdGhhdCBjYXVzZWQgdGhlIGVycm9yIG9yaWdpbmFsbHkuXHJcbiAgICAgKiBAcGFyYW0gcmVhc29uIC0gVGhlIHJlamVjdGlvbiByZWFzb24uXHJcbiAgICAgKiBAcGFyYW0gaW5mbyAtIEFuIG9wdGlvbmFsIGluZm9ybWF0aXZlIG1lc3NhZ2Ugc3BlY2lmeWluZyB0aGUgY2lyY3Vtc3RhbmNlcyBpbiB3aGljaCB0aGUgZXJyb3Igd2FzIHRocm93bi5cclxuICAgICAqICAgICAgICAgICAgICAgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiYC5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCByZWFzb24/OiBhbnksIGluZm8/OiBzdHJpbmcpIHtcclxuICAgICAgICBzdXBlcigoaW5mbyA/PyBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCIpICsgXCIgUmVhc29uOiBcIiArIGVycm9yTWVzc2FnZShyZWFzb24pLCB7IGNhdXNlOiByZWFzb24gfSk7XHJcbiAgICAgICAgdGhpcy5wcm9taXNlID0gcHJvbWlzZTtcclxuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCI7XHJcbiAgICB9XHJcbn1cclxuXHJcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4gPSAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XHJcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IgPSAocmVhc29uPzogYW55KSA9PiB2b2lkO1xyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciA9IChjYXVzZT86IGFueSkgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+ID0gKHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+LCByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yKSA9PiB2b2lkO1xyXG5cclxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcclxuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCk6IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj47XHJcbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBXcmFwcyBhIGNhbmNlbGxhYmxlIHByb21pc2UgYWxvbmcgd2l0aCBpdHMgcmVzb2x1dGlvbiBtZXRob2RzLlxyXG4gKiBUaGUgYG9uY2FuY2VsbGVkYCBmaWVsZCB3aWxsIGJlIG51bGwgaW5pdGlhbGx5IGJ1dCBtYXkgYmUgc2V0IHRvIHByb3ZpZGUgYSBjdXN0b20gY2FuY2VsbGF0aW9uIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcclxuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcclxuICAgIHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+O1xyXG4gICAgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcjtcclxuICAgIG9uY2FuY2VsbGVkOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgfCBudWxsO1xyXG59XHJcblxyXG5pbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUge1xyXG4gICAgcmVhZG9ubHkgcm9vdDogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGU7XHJcbiAgICByZXNvbHZpbmc6IGJvb2xlYW47XHJcbiAgICBzZXR0bGVkOiBib29sZWFuO1xyXG4gICAgcmVhc29uPzogQ2FuY2VsRXJyb3I7XHJcbn1cclxuXHJcbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXHJcbmNvbnN0IGJhcnJpZXJTeW0gPSBTeW1ib2woXCJiYXJyaWVyXCIpO1xyXG5jb25zdCBjYW5jZWxJbXBsU3ltID0gU3ltYm9sKFwiY2FuY2VsSW1wbFwiKTtcclxuY29uc3Qgc3BlY2llczogdHlwZW9mIFN5bWJvbC5zcGVjaWVzID0gU3ltYm9sLnNwZWNpZXMgPz8gU3ltYm9sKFwic3BlY2llc1BvbHlmaWxsXCIpO1xyXG5cclxuLyoqXHJcbiAqIEEgcHJvbWlzZSB3aXRoIGFuIGF0dGFjaGVkIG1ldGhvZCBmb3IgY2FuY2VsbGluZyBsb25nLXJ1bm5pbmcgb3BlcmF0aW9ucyAoc2VlIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2UjY2FuY2VsfSkuXHJcbiAqIENhbmNlbGxhdGlvbiBjYW4gb3B0aW9uYWxseSBiZSBib3VuZCB0byBhbiB7QGxpbmsgQWJvcnRTaWduYWx9XHJcbiAqIGZvciBiZXR0ZXIgY29tcG9zYWJpbGl0eSAoc2VlIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2UjY2FuY2VsT259KS5cclxuICpcclxuICogQ2FuY2VsbGluZyBhIHBlbmRpbmcgcHJvbWlzZSB3aWxsIHJlc3VsdCBpbiBhbiBpbW1lZGlhdGUgcmVqZWN0aW9uXHJcbiAqIHdpdGggYW4gaW5zdGFuY2Ugb2Yge0BsaW5rIENhbmNlbEVycm9yfSBhcyByZWFzb24sXHJcbiAqIGJ1dCB3aG9ldmVyIHN0YXJ0ZWQgdGhlIHByb21pc2Ugd2lsbCBiZSByZXNwb25zaWJsZVxyXG4gKiBmb3IgYWN0dWFsbHkgYWJvcnRpbmcgdGhlIHVuZGVybHlpbmcgb3BlcmF0aW9uLlxyXG4gKiBUbyB0aGlzIHB1cnBvc2UsIHRoZSBjb25zdHJ1Y3RvciBhbmQgYWxsIGNoYWluaW5nIG1ldGhvZHNcclxuICogYWNjZXB0IG9wdGlvbmFsIGNhbmNlbGxhdGlvbiBjYWxsYmFja3MuXHJcbiAqXHJcbiAqIElmIGEgYENhbmNlbGxhYmxlUHJvbWlzZWAgc3RpbGwgcmVzb2x2ZXMgYWZ0ZXIgaGF2aW5nIGJlZW4gY2FuY2VsbGVkLFxyXG4gKiB0aGUgcmVzdWx0IHdpbGwgYmUgZGlzY2FyZGVkLiBJZiBpdCByZWplY3RzLCB0aGUgcmVhc29uXHJcbiAqIHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbixcclxuICogd3JhcHBlZCBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gaW5zdGFuY2UuXHJcbiAqIFRvIGZhY2lsaXRhdGUgdGhlIGhhbmRsaW5nIG9mIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyxcclxuICogY2FuY2VsbGVkIGBDYW5jZWxsYWJsZVByb21pc2VgcyB3aWxsIF9ub3RfIHJlcG9ydCB1bmhhbmRsZWQgYENhbmNlbEVycm9yYHNcclxuICogd2hvc2UgYGNhdXNlYCBmaWVsZCBpcyB0aGUgc2FtZSBhcyB0aGUgb25lIHdpdGggd2hpY2ggdGhlIGN1cnJlbnQgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLlxyXG4gKlxyXG4gKiBBbGwgdXN1YWwgcHJvbWlzZSBtZXRob2RzIGFyZSBkZWZpbmVkIGFuZCByZXR1cm4gYSBgQ2FuY2VsbGFibGVQcm9taXNlYFxyXG4gKiB3aG9zZSBjYW5jZWwgbWV0aG9kIHdpbGwgY2FuY2VsIHRoZSBwYXJlbnQgb3BlcmF0aW9uIGFzIHdlbGwsIHByb3BhZ2F0aW5nIHRoZSBjYW5jZWxsYXRpb24gcmVhc29uXHJcbiAqIHVwd2FyZHMgdGhyb3VnaCBwcm9taXNlIGNoYWlucy5cclxuICogQ29udmVyc2VseSwgY2FuY2VsbGluZyBhIHByb21pc2Ugd2lsbCBub3QgYXV0b21hdGljYWxseSBjYW5jZWwgZGVwZW5kZW50IHByb21pc2VzIGRvd25zdHJlYW06XHJcbiAqIGBgYHRzXHJcbiAqIGxldCByb290ID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcclxuICogbGV0IGNoaWxkMSA9IHJvb3QudGhlbigoKSA9PiB7IC4uLiB9KTtcclxuICogbGV0IGNoaWxkMiA9IGNoaWxkMS50aGVuKCgpID0+IHsgLi4uIH0pO1xyXG4gKiBsZXQgY2hpbGQzID0gcm9vdC5jYXRjaCgoKSA9PiB7IC4uLiB9KTtcclxuICogY2hpbGQxLmNhbmNlbCgpOyAvLyBDYW5jZWxzIGNoaWxkMSBhbmQgcm9vdCwgYnV0IG5vdCBjaGlsZDIgb3IgY2hpbGQzXHJcbiAqIGBgYFxyXG4gKiBDYW5jZWxsaW5nIGEgcHJvbWlzZSB0aGF0IGhhcyBhbHJlYWR5IHNldHRsZWQgaXMgc2FmZSBhbmQgaGFzIG5vIGNvbnNlcXVlbmNlLlxyXG4gKlxyXG4gKiBUaGUgYGNhbmNlbGAgbWV0aG9kIHJldHVybnMgYSBwcm9taXNlIHRoYXQgX2Fsd2F5cyBmdWxmaWxsc19cclxuICogYWZ0ZXIgdGhlIHdob2xlIGNoYWluIGhhcyBwcm9jZXNzZWQgdGhlIGNhbmNlbCByZXF1ZXN0XHJcbiAqIGFuZCBhbGwgYXR0YWNoZWQgY2FsbGJhY2tzIHVwIHRvIHRoYXQgbW9tZW50IGhhdmUgcnVuLlxyXG4gKlxyXG4gKiBBbGwgRVMyMDI0IHByb21pc2UgbWV0aG9kcyAoc3RhdGljIGFuZCBpbnN0YW5jZSkgYXJlIGRlZmluZWQgb24gQ2FuY2VsbGFibGVQcm9taXNlLFxyXG4gKiBidXQgYWN0dWFsIGF2YWlsYWJpbGl0eSBtYXkgdmFyeSB3aXRoIE9TL3dlYnZpZXcgdmVyc2lvbi5cclxuICpcclxuICogSW4gbGluZSB3aXRoIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nLFxyXG4gKiBgQ2FuY2VsbGFibGVQcm9taXNlYCBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLlxyXG4gKiBFeHRlbmRlcnMgc2hvdWxkIHRha2UgY2FyZSB0byBwcm92aWRlIHRoZWlyIG93biBtZXRob2QgaW1wbGVtZW50YXRpb25zLlxyXG4gKiBUaGlzIG1pZ2h0IGJlIHJlY29uc2lkZXJlZCBpbiBjYXNlIHRoZSBwcm9wb3NhbCBpcyByZXRpcmVkLlxyXG4gKlxyXG4gKiBDYW5jZWxsYWJsZVByb21pc2UgaXMgYSB3cmFwcGVyIGFyb3VuZCB0aGUgRE9NIFByb21pc2Ugb2JqZWN0XHJcbiAqIGFuZCBpcyBjb21wbGlhbnQgd2l0aCB0aGUgW1Byb21pc2VzL0ErIHNwZWNpZmljYXRpb25dKGh0dHBzOi8vcHJvbWlzZXNhcGx1cy5jb20vKVxyXG4gKiAoaXQgcGFzc2VzIHRoZSBbY29tcGxpYW5jZSBzdWl0ZV0oaHR0cHM6Ly9naXRodWIuY29tL3Byb21pc2VzLWFwbHVzL3Byb21pc2VzLXRlc3RzKSlcclxuICogaWYgc28gaXMgdGhlIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24uXHJcbiAqL1xyXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGFibGVQcm9taXNlPFQ+IGV4dGVuZHMgUHJvbWlzZTxUPiBpbXBsZW1lbnRzIFByb21pc2VMaWtlPFQ+LCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcclxuICAgIC8vIFByaXZhdGUgZmllbGRzLlxyXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xyXG4gICAgcHJpdmF0ZSBbYmFycmllclN5bV0hOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiB8IG51bGw7XHJcbiAgICAvKiogQGludGVybmFsICovXHJcbiAgICBwcml2YXRlIHJlYWRvbmx5IFtjYW5jZWxJbXBsU3ltXSE6IChyZWFzb246IENhbmNlbEVycm9yKSA9PiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGBDYW5jZWxsYWJsZVByb21pc2VgLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBleGVjdXRvciAtIEEgY2FsbGJhY2sgdXNlZCB0byBpbml0aWFsaXplIHRoZSBwcm9taXNlLiBUaGlzIGNhbGxiYWNrIGlzIHBhc3NlZCB0d28gYXJndW1lbnRzOlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgYSBgcmVzb2x2ZWAgY2FsbGJhY2sgdXNlZCB0byByZXNvbHZlIHRoZSBwcm9taXNlIHdpdGggYSB2YWx1ZVxyXG4gICAgICogICAgICAgICAgICAgICAgICAgb3IgdGhlIHJlc3VsdCBvZiBhbm90aGVyIHByb21pc2UgKHBvc3NpYmx5IGNhbmNlbGxhYmxlKSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGFuZCBhIGByZWplY3RgIGNhbGxiYWNrIHVzZWQgdG8gcmVqZWN0IHRoZSBwcm9taXNlIHdpdGggYSBwcm92aWRlZCByZWFzb24gb3IgZXJyb3IuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBJZiB0aGUgdmFsdWUgcHJvdmlkZWQgdG8gdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgKGl0IGhhcyBhIGB0aGVuYCBfYW5kXyBhIGBjYW5jZWxgIG1ldGhvZCksXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBmb3J3YXJkZWQgdG8gdGhhdCBvYmplY3QgYW5kIHRoZSBvbmNhbmNlbGxlZCB3aWxsIG5vdCBiZSBpbnZva2VkIGFueW1vcmUuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBJZiBhbnkgb25lIG9mIHRoZSB0d28gY2FsbGJhY2tzIGlzIGNhbGxlZCBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZCxcclxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoZSBwcm92aWRlZCB2YWx1ZXMgd2lsbCBiZSBjYW5jZWxsZWQgYW5kIHJlc29sdmVkIGFzIHVzdWFsLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgYnV0IHRoZWlyIHJlc3VsdHMgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBIb3dldmVyLCBpZiB0aGUgcmVzb2x1dGlvbiBwcm9jZXNzIHVsdGltYXRlbHkgZW5kcyB1cCBpbiBhIHJlamVjdGlvblxyXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhhdCBpcyBub3QgZHVlIHRvIGNhbmNlbGxhdGlvbiwgdGhlIHJlamVjdGlvbiByZWFzb25cclxuICAgICAqICAgICAgICAgICAgICAgICAgIHdpbGwgYmUgd3JhcHBlZCBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn1cclxuICAgICAqICAgICAgICAgICAgICAgICAgIGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXHJcbiAgICAgKiBAcGFyYW0gb25jYW5jZWxsZWQgLSBJdCBpcyB0aGUgY2FsbGVyJ3MgcmVzcG9uc2liaWxpdHkgdG8gZW5zdXJlIHRoYXQgYW55IG9wZXJhdGlvblxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgc3RhcnRlZCBieSB0aGUgZXhlY3V0b3IgaXMgcHJvcGVybHkgaGFsdGVkIHVwb24gY2FuY2VsbGF0aW9uLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVGhpcyBvcHRpb25hbCBjYWxsYmFjayBjYW4gYmUgdXNlZCB0byB0aGF0IHB1cnBvc2UuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJdCB3aWxsIGJlIGNhbGxlZCBfc3luY2hyb25vdXNseV8gd2l0aCBhIGNhbmNlbGxhdGlvbiBjYXVzZVxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2hlbiBjYW5jZWxsYXRpb24gaXMgcmVxdWVzdGVkLCBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHJlamVjdGVkXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSwgYnV0IF9iZWZvcmVfXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbnkge0BsaW5rIHRoZW59L3tAbGluayBjYXRjaH0ve0BsaW5rIGZpbmFsbHl9IGNhbGxiYWNrIHJ1bnMuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aWxsIG9ubHkgZnVsZmlsbCBhZnRlciB0aGUgZm9ybWVyIGhhcyBzZXR0bGVkLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVW5oYW5kbGVkIGV4Y2VwdGlvbnMgb3IgcmVqZWN0aW9ucyBmcm9tIHRoZSBjYWxsYmFjayB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyB1bmhhbmRsZWQgcmVqZWN0aW9ucy5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgY2FsbGVkIGJlZm9yZSBjYW5jZWxsYXRpb24gd2l0aCBhIGNhbmNlbGxhYmxlIHByb21pc2UsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhpcyBwcm9taXNlIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gdGhhdCBwcm9taXNlLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW5kIHRoZSBvcmlnaW5hbCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihleGVjdXRvcjogQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4sIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKSB7XHJcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcclxuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcclxuICAgICAgICBzdXBlcigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcclxuXHJcbiAgICAgICAgaWYgKCh0aGlzLmNvbnN0cnVjdG9yIGFzIGFueSlbc3BlY2llc10gIT09IFByb21pc2UpIHtcclxuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZSBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLiBQbGVhc2UgcmVmcmFpbiBmcm9tIG92ZXJyaWRpbmcgdGhlIFtTeW1ib2wuc3BlY2llc10gc3RhdGljIHByb3BlcnR5LlwiKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGxldCBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0ge1xyXG4gICAgICAgICAgICBwcm9taXNlOiB0aGlzLFxyXG4gICAgICAgICAgICByZXNvbHZlLFxyXG4gICAgICAgICAgICByZWplY3QsXHJcbiAgICAgICAgICAgIGdldCBvbmNhbmNlbGxlZCgpIHsgcmV0dXJuIG9uY2FuY2VsbGVkID8/IG51bGw7IH0sXHJcbiAgICAgICAgICAgIHNldCBvbmNhbmNlbGxlZChjYikgeyBvbmNhbmNlbGxlZCA9IGNiID8/IHVuZGVmaW5lZDsgfVxyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIGNvbnN0IHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcclxuICAgICAgICAgICAgZ2V0IHJvb3QoKSB7IHJldHVybiBzdGF0ZTsgfSxcclxuICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcclxuICAgICAgICAgICAgc2V0dGxlZDogZmFsc2VcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICAvLyBTZXR1cCBjYW5jZWxsYXRpb24gc3lzdGVtLlxyXG4gICAgICAgIHZvaWQgT2JqZWN0LmRlZmluZVByb3BlcnRpZXModGhpcywge1xyXG4gICAgICAgICAgICBbYmFycmllclN5bV06IHtcclxuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiB0cnVlLFxyXG4gICAgICAgICAgICAgICAgdmFsdWU6IG51bGxcclxuICAgICAgICAgICAgfSxcclxuICAgICAgICAgICAgW2NhbmNlbEltcGxTeW1dOiB7XHJcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICB2YWx1ZTogY2FuY2VsbGVyRm9yKHByb21pc2UsIHN0YXRlKVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSk7XHJcblxyXG4gICAgICAgIC8vIFJ1biB0aGUgYWN0dWFsIGV4ZWN1dG9yLlxyXG4gICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgc3RhdGUpO1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIGV4ZWN1dG9yKHJlc29sdmVyRm9yKHByb21pc2UsIHN0YXRlKSwgcmVqZWN0b3IpO1xyXG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7XHJcbiAgICAgICAgICAgICAgICBjb25zb2xlLmxvZyhcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gQ2FuY2VsbGFibGVQcm9taXNlIGV4ZWN1dG9yLlwiLCBlcnIpO1xyXG4gICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENhbmNlbHMgaW1tZWRpYXRlbHkgdGhlIGV4ZWN1dGlvbiBvZiB0aGUgb3BlcmF0aW9uIGFzc29jaWF0ZWQgd2l0aCB0aGlzIHByb21pc2UuXHJcbiAgICAgKiBUaGUgcHJvbWlzZSByZWplY3RzIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9IGluc3RhbmNlIGFzIHJlYXNvbixcclxuICAgICAqIHdpdGggdGhlIHtAbGluayBDYW5jZWxFcnJvciNjYXVzZX0gcHJvcGVydHkgc2V0IHRvIHRoZSBnaXZlbiBhcmd1bWVudCwgaWYgYW55LlxyXG4gICAgICpcclxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIGFmdGVyIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQ7XHJcbiAgICAgKiByZXBlYXRlZCBjYWxscyBpbiBwYXJ0aWN1bGFyIGFyZSBzYWZlLCBidXQgb25seSB0aGUgZmlyc3Qgb25lXHJcbiAgICAgKiB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBgQ2FuY2VsRXJyb3JgIGV4Y2VwdGlvbiBfbmVlZCBub3RfIGJlIGhhbmRsZWQgZXhwbGljaXRseSBfb24gdGhlIHByb21pc2VzIHRoYXQgYXJlIGJlaW5nIGNhbmNlbGxlZDpfXHJcbiAgICAgKiBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aXRoIG5vIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyIGRvZXMgbm90IHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudC5cclxuICAgICAqIFRoZXJlZm9yZSwgdGhlIGZvbGxvd2luZyBpZGlvbXMgYXJlIGFsbCBlcXVhbGx5IGNvcnJlY3Q6XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS5jYW5jZWwoKTtcclxuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhbmNlbCgpO1xyXG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2F0Y2goLi4uKS5jYW5jZWwoKTtcclxuICAgICAqIGBgYFxyXG4gICAgICogV2hlbmV2ZXIgc29tZSBjYW5jZWxsZWQgcHJvbWlzZSBpbiBhIGNoYWluIHJlamVjdHMgd2l0aCBhIGBDYW5jZWxFcnJvcmBcclxuICAgICAqIHdpdGggdGhlIHNhbWUgY2FuY2VsbGF0aW9uIGNhdXNlIGFzIGl0c2VsZiwgdGhlIGVycm9yIHdpbGwgYmUgZGlzY2FyZGVkIHNpbGVudGx5LlxyXG4gICAgICogSG93ZXZlciwgdGhlIGBDYW5jZWxFcnJvcmAgX3dpbGwgc3RpbGwgYmUgZGVsaXZlcmVkXyB0byBhbGwgYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXJzXHJcbiAgICAgKiBhZGRlZCBieSB7QGxpbmsgdGhlbn0gYW5kIHJlbGF0ZWQgbWV0aG9kczpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xyXG4gICAgICogY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS5jYXRjaChjb25zb2xlLmxvZyk7XHJcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gQSBDYW5jZWxFcnJvciBpcyBwcmludGVkIHRvIHRoZSBjb25zb2xlLlxyXG4gICAgICogYGBgXHJcbiAgICAgKiBJZiB0aGUgYENhbmNlbEVycm9yYCBpcyBub3QgaGFuZGxlZCBkb3duc3RyZWFtIGJ5IHRoZSB0aW1lIGl0IHJlYWNoZXNcclxuICAgICAqIGEgX25vbi1jYW5jZWxsZWRfIHByb21pc2UsIGl0IF93aWxsXyB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQsXHJcbiAgICAgKiBqdXN0IGxpa2Ugbm9ybWFsIHJlamVjdGlvbnMgd291bGQ6XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcclxuICAgICAqIGxldCBjaGFpbmVkID0gY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS50aGVuKCgpID0+IHsgLi4uIH0pOyAvLyBObyBjYXRjaC4uLlxyXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIFVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQgb24gY2hhaW5lZCFcclxuICAgICAqIGBgYFxyXG4gICAgICogVGhlcmVmb3JlLCBpdCBpcyBpbXBvcnRhbnQgdG8gZWl0aGVyIGNhbmNlbCB3aG9sZSBwcm9taXNlIGNoYWlucyBmcm9tIHRoZWlyIHRhaWwsXHJcbiAgICAgKiBhcyBzaG93biBpbiB0aGUgY29ycmVjdCBpZGlvbXMgYWJvdmUsIG9yIHRha2UgY2FyZSBvZiBoYW5kbGluZyBlcnJvcnMgZXZlcnl3aGVyZS5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBBIGNhbmNlbGxhYmxlIHByb21pc2UgdGhhdCBfZnVsZmlsbHNfIGFmdGVyIHRoZSBjYW5jZWwgY2FsbGJhY2sgKGlmIGFueSlcclxuICAgICAqIGFuZCBhbGwgaGFuZGxlcnMgYXR0YWNoZWQgdXAgdG8gdGhlIGNhbGwgdG8gY2FuY2VsIGhhdmUgcnVuLlxyXG4gICAgICogSWYgdGhlIGNhbmNlbCBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGJ5IGBjYW5jZWxgXHJcbiAgICAgKiB3aWxsIGFsc28gd2FpdCBmb3IgdGhhdCB0aGVuYWJsZSB0byBzZXR0bGUuXHJcbiAgICAgKiBUaGlzIGVuYWJsZXMgY2FsbGVycyB0byB3YWl0IGZvciB0aGUgY2FuY2VsbGVkIG9wZXJhdGlvbiB0byB0ZXJtaW5hdGVcclxuICAgICAqIHdpdGhvdXQgYmVpbmcgZm9yY2VkIHRvIGhhbmRsZSBwb3RlbnRpYWwgZXJyb3JzIGF0IHRoZSBjYWxsIHNpdGUuXHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCkudGhlbigoKSA9PiB7XHJcbiAgICAgKiAgICAgLy8gQ2xlYW51cCBmaW5pc2hlZCwgaXQncyBzYWZlIHRvIGRvIHNvbWV0aGluZyBlbHNlLlxyXG4gICAgICogfSwgKGVycikgPT4ge1xyXG4gICAgICogICAgIC8vIFVucmVhY2hhYmxlOiB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIGNhbmNlbCB3aWxsIG5ldmVyIHJlamVjdC5cclxuICAgICAqIH0pO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBOb3RlIHRoYXQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBfbm90XyBoYW5kbGUgaW1wbGljaXRseSBhbnkgcmVqZWN0aW9uXHJcbiAgICAgKiB0aGF0IG1pZ2h0IGhhdmUgb2NjdXJyZWQgYWxyZWFkeSBpbiB0aGUgY2FuY2VsbGVkIGNoYWluLlxyXG4gICAgICogSXQgd2lsbCBqdXN0IHRyYWNrIHdoZXRoZXIgcmVnaXN0ZXJlZCBoYW5kbGVycyBoYXZlIGJlZW4gZXhlY3V0ZWQgb3Igbm90LlxyXG4gICAgICogVGhlcmVmb3JlLCB1bmhhbmRsZWQgcmVqZWN0aW9ucyB3aWxsIG5ldmVyIGJlIHNpbGVudGx5IGhhbmRsZWQgYnkgY2FsbGluZyBjYW5jZWwuXHJcbiAgICAgKi9cclxuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcclxuICAgICAgICAgICAgLy8gSU5WQVJJQU5UOiB0aGUgcmVzdWx0IG9mIHRoaXNbY2FuY2VsSW1wbFN5bV0gYW5kIHRoZSBiYXJyaWVyIGRvIG5vdCBldmVyIHJlamVjdC5cclxuICAgICAgICAgICAgLy8gVW5mb3J0dW5hdGVseSBtYWNPUyBIaWdoIFNpZXJyYSBkb2VzIG5vdCBzdXBwb3J0IFByb21pc2UuYWxsU2V0dGxlZC5cclxuICAgICAgICAgICAgUHJvbWlzZS5hbGwoW1xyXG4gICAgICAgICAgICAgICAgdGhpc1tjYW5jZWxJbXBsU3ltXShuZXcgQ2FuY2VsRXJyb3IoXCJQcm9taXNlIGNhbmNlbGxlZC5cIiwgeyBjYXVzZSB9KSksXHJcbiAgICAgICAgICAgICAgICBjdXJyZW50QmFycmllcih0aGlzKVxyXG4gICAgICAgICAgICBdKS50aGVuKCgpID0+IHJlc29sdmUoKSwgKCkgPT4gcmVzb2x2ZSgpKTtcclxuICAgICAgICB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEJpbmRzIHByb21pc2UgY2FuY2VsbGF0aW9uIHRvIHRoZSBhYm9ydCBldmVudCBvZiB0aGUgZ2l2ZW4ge0BsaW5rIEFib3J0U2lnbmFsfS5cclxuICAgICAqIElmIHRoZSBzaWduYWwgaGFzIGFscmVhZHkgYWJvcnRlZCwgdGhlIHByb21pc2Ugd2lsbCBiZSBjYW5jZWxsZWQgaW1tZWRpYXRlbHkuXHJcbiAgICAgKiBXaGVuIGVpdGhlciBjb25kaXRpb24gaXMgdmVyaWZpZWQsIHRoZSBjYW5jZWxsYXRpb24gY2F1c2Ugd2lsbCBiZSBzZXRcclxuICAgICAqIHRvIHRoZSBzaWduYWwncyBhYm9ydCByZWFzb24gKHNlZSB7QGxpbmsgQWJvcnRTaWduYWwjcmVhc29ufSkuXHJcbiAgICAgKlxyXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgKG9yIGlmIHRoZSBzaWduYWwgYWJvcnRzKSBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQuXHJcbiAgICAgKiBPbmx5IHRoZSBmaXJzdCBzaWduYWwgdG8gYWJvcnQgd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBGb3IgbW9yZSBkZXRhaWxzIGFib3V0IHRoZSBjYW5jZWxsYXRpb24gcHJvY2VzcyxcclxuICAgICAqIHNlZSB7QGxpbmsgY2FuY2VsfSBhbmQgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgIGNvbnN0cnVjdG9yLlxyXG4gICAgICpcclxuICAgICAqIFRoaXMgbWV0aG9kIGVuYWJsZXMgYGF3YWl0YGluZyBjYW5jZWxsYWJsZSBwcm9taXNlcyB3aXRob3V0IGhhdmluZ1xyXG4gICAgICogdG8gc3RvcmUgdGhlbSBmb3IgZnV0dXJlIGNhbmNlbGxhdGlvbiwgZS5nLjpcclxuICAgICAqIGBgYHRzXHJcbiAgICAgKiBhd2FpdCBsb25nUnVubmluZ09wZXJhdGlvbigpLmNhbmNlbE9uKHNpZ25hbCk7XHJcbiAgICAgKiBgYGBcclxuICAgICAqIGluc3RlYWQgb2Y6XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogbGV0IHByb21pc2VUb0JlQ2FuY2VsbGVkID0gbG9uZ1J1bm5pbmdPcGVyYXRpb24oKTtcclxuICAgICAqIGF3YWl0IHByb21pc2VUb0JlQ2FuY2VsbGVkO1xyXG4gICAgICogYGBgXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhpcyBwcm9taXNlLCBmb3IgbWV0aG9kIGNoYWluaW5nLlxyXG4gICAgICovXHJcbiAgICBjYW5jZWxPbihzaWduYWw6IEFib3J0U2lnbmFsKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBpZiAoc2lnbmFsLmFib3J0ZWQpIHtcclxuICAgICAgICAgICAgdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKVxyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHNpZ25hbC5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbiksIHtjYXB0dXJlOiB0cnVlfSk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICByZXR1cm4gdGhpcztcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEF0dGFjaGVzIGNhbGxiYWNrcyBmb3IgdGhlIHJlc29sdXRpb24gYW5kL29yIHJlamVjdGlvbiBvZiB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxyXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXHJcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xyXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxyXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIG9uZnVsZmlsbGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVzb2x2ZWQuXHJcbiAgICAgKiBAcGFyYW0gb25yZWplY3RlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlamVjdGVkLlxyXG4gICAgICogQHJldHVybnMgQSBgQ2FuY2VsbGFibGVQcm9taXNlYCBmb3IgdGhlIGNvbXBsZXRpb24gb2Ygd2hpY2hldmVyIGNhbGxiYWNrIGlzIGV4ZWN1dGVkLlxyXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgaXMgaG9va2VkIHVwIHRvIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24gcmVxdWVzdHMgdXAgdGhlIGNoYWluLCBidXQgbm90IGRvd246XHJcbiAgICAgKlxyXG4gICAgICogICAtIGlmIHRoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIGJlIGludm9rZWQgd2l0aCBhIGBDYW5jZWxFcnJvcmBcclxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2UgX3dpbGwgcmVzb2x2ZSByZWd1bGFybHlfIHdpdGggaXRzIHJlc3VsdDtcclxuICAgICAqICAgLSBjb252ZXJzZWx5LCBpZiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIF90aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkIHRvbztfXHJcbiAgICAgKiAgICAgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgc3RpbGwgYmUgaW52b2tlZCB3aXRoIHRoZSBwYXJlbnQncyBgQ2FuY2VsRXJyb3JgLFxyXG4gICAgICogICAgIGJ1dCBpdHMgcmVzdWx0IHdpbGwgYmUgZGlzY2FyZGVkXHJcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYSBgQ2FuY2VsRXJyb3JgIGFzIHdlbGwuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcclxuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxyXG4gICAgICpcclxuICAgICAqIElmIGVpdGhlciBjYWxsYmFjayByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcclxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxyXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqL1xyXG4gICAgdGhlbjxUUmVzdWx0MSA9IFQsIFRSZXN1bHQyID0gbmV2ZXI+KG9uZnVsZmlsbGVkPzogKCh2YWx1ZTogVCkgPT4gVFJlc3VsdDEgfCBQcm9taXNlTGlrZTxUUmVzdWx0MT4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gVFJlc3VsdDIgfCBQcm9taXNlTGlrZTxUUmVzdWx0Mj4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQyPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUUmVzdWx0MSB8IFRSZXN1bHQyPiB7XHJcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcclxuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUudGhlbiBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gTk9URTogVHlwZVNjcmlwdCdzIGJ1aWx0LWluIHR5cGUgZm9yIHRoZW4gaXMgYnJva2VuLFxyXG4gICAgICAgIC8vIGFzIGl0IGFsbG93cyBzcGVjaWZ5aW5nIGFuIGFyYml0cmFyeSBUUmVzdWx0MSAhPSBUIGV2ZW4gd2hlbiBvbmZ1bGZpbGxlZCBpcyBub3QgYSBmdW5jdGlvbi5cclxuICAgICAgICAvLyBXZSBjYW5ub3QgZml4IGl0IGlmIHdlIHdhbnQgdG8gQ2FuY2VsbGFibGVQcm9taXNlIHRvIGltcGxlbWVudCBQcm9taXNlTGlrZTxUPi5cclxuXHJcbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZnVsZmlsbGVkKSkgeyBvbmZ1bGZpbGxlZCA9IGlkZW50aXR5IGFzIGFueTsgfVxyXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbnJlamVjdGVkKSkgeyBvbnJlamVjdGVkID0gdGhyb3dlcjsgfVxyXG5cclxuICAgICAgICBpZiAob25mdWxmaWxsZWQgPT09IGlkZW50aXR5ICYmIG9ucmVqZWN0ZWQgPT0gdGhyb3dlcikge1xyXG4gICAgICAgICAgICAvLyBTaG9ydGN1dCBmb3IgdHJpdmlhbCBhcmd1bWVudHMuXHJcbiAgICAgICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlKSA9PiByZXNvbHZlKHRoaXMgYXMgYW55KSk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBjb25zdCBiYXJyaWVyOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiA9IHt9O1xyXG4gICAgICAgIHRoaXNbYmFycmllclN5bV0gPSBiYXJyaWVyO1xyXG5cclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUUmVzdWx0MSB8IFRSZXN1bHQyPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgIHZvaWQgc3VwZXIudGhlbihcclxuICAgICAgICAgICAgICAgICh2YWx1ZSkgPT4ge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XHJcbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcclxuXHJcbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbmZ1bGZpbGxlZCEodmFsdWUpKTtcclxuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfSxcclxuICAgICAgICAgICAgICAgIChyZWFzb24/KSA9PiB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cclxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xyXG5cclxuICAgICAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9ucmVqZWN0ZWQhKHJlYXNvbikpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICk7XHJcbiAgICAgICAgfSwgYXN5bmMgKGNhdXNlPykgPT4ge1xyXG4gICAgICAgICAgICAvL2NhbmNlbGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gb25jYW5jZWxsZWQ/LihjYXVzZSk7XHJcbiAgICAgICAgICAgIH0gZmluYWxseSB7XHJcbiAgICAgICAgICAgICAgICBhd2FpdCB0aGlzLmNhbmNlbChjYXVzZSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEF0dGFjaGVzIGEgY2FsbGJhY2sgZm9yIG9ubHkgdGhlIHJlamVjdGlvbiBvZiB0aGUgUHJvbWlzZS5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXHJcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXHJcbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcclxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXHJcbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBJdCBpcyBlcXVpdmFsZW50IHRvXHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogY2FuY2VsbGFibGVQcm9taXNlLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XHJcbiAgICAgKiBgYGBcclxuICAgICAqIGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxyXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXHJcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxyXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXHJcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cclxuICAgICAqXHJcbiAgICAgKiBJZiBgb25yZWplY3RlZGAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXHJcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcclxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXHJcbiAgICAgKi9cclxuICAgIGNhdGNoPFRSZXN1bHQgPSBuZXZlcj4ob25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IChQcm9taXNlTGlrZTxUUmVzdWx0PiB8IFRSZXN1bHQpKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQgfCBUUmVzdWx0PiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEF0dGFjaGVzIGEgY2FsbGJhY2sgdGhhdCBpcyBpbnZva2VkIHdoZW4gdGhlIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBzZXR0bGVkIChmdWxmaWxsZWQgb3IgcmVqZWN0ZWQpLiBUaGVcclxuICAgICAqIHJlc29sdmVkIHZhbHVlIGNhbm5vdCBiZSBhY2Nlc3NlZCBvciBtb2RpZmllZCBmcm9tIHRoZSBjYWxsYmFjay5cclxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgc2V0dGxlIGluIHRoZSBzYW1lIHN0YXRlIGFzIHRoZSBvcmlnaW5hbCBvbmVcclxuICAgICAqIGFmdGVyIHRoZSBwcm92aWRlZCBjYWxsYmFjayBoYXMgY29tcGxldGVkIGV4ZWN1dGlvbixcclxuICAgICAqIHVubGVzcyB0aGUgY2FsbGJhY2sgdGhyb3dzIG9yIHJldHVybnMgYSByZWplY3RpbmcgcHJvbWlzZSxcclxuICAgICAqIGluIHdoaWNoIGNhc2UgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3QgYXMgd2VsbC5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXHJcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXHJcbiAgICAgKiBPbmNlIHRoZSBwYXJlbnQgcHJvbWlzZSBzZXR0bGVzLCB0aGUgYG9uZmluYWxseWAgY2FsbGJhY2sgd2lsbCBydW4sXHJcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xyXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxyXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogVGhpcyBtZXRob2QgaXMgaW1wbGVtZW50ZWQgaW4gdGVybXMgb2Yge0BsaW5rIHRoZW59IGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxyXG4gICAgICogSXQgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cclxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxyXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcclxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xyXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXHJcbiAgICAgKlxyXG4gICAgICogSWYgYG9uZmluYWxseWAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXHJcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcclxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXHJcbiAgICAgKi9cclxuICAgIGZpbmFsbHkob25maW5hbGx5PzogKCgpID0+IHZvaWQpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XHJcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLmZpbmFsbHkgY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZpbmFsbHkpKSB7XHJcbiAgICAgICAgICAgIHJldHVybiB0aGlzLnRoZW4ob25maW5hbGx5LCBvbmZpbmFsbHksIG9uY2FuY2VsbGVkKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4oXHJcbiAgICAgICAgICAgICh2YWx1ZSkgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4gdmFsdWUpLFxyXG4gICAgICAgICAgICAocmVhc29uPykgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4geyB0aHJvdyByZWFzb247IH0pLFxyXG4gICAgICAgICAgICBvbmNhbmNlbGxlZCxcclxuICAgICAgICApO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogV2UgdXNlIHRoZSBgW1N5bWJvbC5zcGVjaWVzXWAgc3RhdGljIHByb3BlcnR5LCBpZiBhdmFpbGFibGUsXHJcbiAgICAgKiB0byBkaXNhYmxlIHRoZSBidWlsdC1pbiBhdXRvbWF0aWMgc3ViY2xhc3NpbmcgZmVhdHVyZXMgZnJvbSB7QGxpbmsgUHJvbWlzZX0uXHJcbiAgICAgKiBJdCBpcyBjcml0aWNhbCBmb3IgcGVyZm9ybWFuY2UgcmVhc29ucyB0aGF0IGV4dGVuZGVycyBkbyBub3Qgb3ZlcnJpZGUgdGhpcy5cclxuICAgICAqIE9uY2UgdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmdcclxuICAgICAqIGlzIGVpdGhlciBhY2NlcHRlZCBvciByZXRpcmVkLCB0aGlzIGltcGxlbWVudGF0aW9uIHdpbGwgaGF2ZSB0byBiZSByZXZpc2VkIGFjY29yZGluZ2x5LlxyXG4gICAgICpcclxuICAgICAqIEBpZ25vcmVcclxuICAgICAqIEBpbnRlcm5hbFxyXG4gICAgICovXHJcbiAgICBzdGF0aWMgZ2V0IFtzcGVjaWVzXSgpIHtcclxuICAgICAgICByZXR1cm4gUHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcclxuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlLCBvciByZWplY3RlZCB3aGVuIGFueSBQcm9taXNlIGlzIHJlamVjdGVkLlxyXG4gICAgICpcclxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGFsbDxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPltdPjtcclxuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IEF3YWl0ZWQ8VFtQXT47IH0+O1xyXG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxyXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcclxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGwoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XHJcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXHJcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSBvciByZWplY3QuXHJcbiAgICAgKlxyXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8UHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUPj5bXT47XHJcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogUHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUW1BdPj47IH0+O1xyXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xyXG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcclxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXHJcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsU2V0dGxlZChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFRoZSBhbnkgZnVuY3Rpb24gcmV0dXJucyBhIHByb21pc2UgdGhhdCBpcyBmdWxmaWxsZWQgYnkgdGhlIGZpcnN0IGdpdmVuIHByb21pc2UgdG8gYmUgZnVsZmlsbGVkLFxyXG4gICAgICogb3IgcmVqZWN0ZWQgd2l0aCBhbiBBZ2dyZWdhdGVFcnJvciBjb250YWluaW5nIGFuIGFycmF5IG9mIHJlamVjdGlvbiByZWFzb25zXHJcbiAgICAgKiBpZiBhbGwgb2YgdGhlIGdpdmVuIHByb21pc2VzIGFyZSByZWplY3RlZC5cclxuICAgICAqIEl0IHJlc29sdmVzIGFsbCBlbGVtZW50cyBvZiB0aGUgcGFzc2VkIGl0ZXJhYmxlIHRvIHByb21pc2VzIGFzIGl0IHJ1bnMgdGhpcyBhbGdvcml0aG0uXHJcbiAgICAgKlxyXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgYW55PFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUW251bWJlcl0+PjtcclxuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xyXG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcclxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXHJcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYW55KGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIFByb21pc2UgdGhhdCBpcyByZXNvbHZlZCBvciByZWplY3RlZCB3aGVuIGFueSBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgYXJlIHJlc29sdmVkIG9yIHJlamVjdGVkLlxyXG4gICAgICpcclxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJhY2U8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xyXG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUW251bWJlcl0+PjtcclxuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJhY2UoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XHJcbiAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgY2FuY2VsbGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgY2FuY2VsPFQgPSBuZXZlcj4oY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIGNvbnN0IHAgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcclxuICAgICAgICBwLmNhbmNlbChjYXVzZSk7XHJcbiAgICAgICAgcmV0dXJuIHA7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGNhbmNlbHNcclxuICAgICAqIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogSWYgdGhlIHtAbGluayBBYm9ydFNpZ25hbC50aW1lb3V0fSBmYWN0b3J5IG1ldGhvZCBpcyBhdmFpbGFibGUsXHJcbiAgICAgKiBpdCBpcyB1c2VkIHRvIGJhc2UgdGhlIHRpbWVvdXQgb24gX2FjdGl2ZV8gdGltZSByYXRoZXIgdGhhbiBfZWxhcHNlZF8gdGltZS5cclxuICAgICAqIE90aGVyd2lzZSwgYHRpbWVvdXRgIGZhbGxzIGJhY2sgdG8ge0BsaW5rIHNldFRpbWVvdXR9LlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgdGltZW91dDxUID0gbmV2ZXI+KG1pbGxpc2Vjb25kczogbnVtYmVyLCBjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xyXG4gICAgICAgIGlmIChBYm9ydFNpZ25hbCAmJiB0eXBlb2YgQWJvcnRTaWduYWwgPT09ICdmdW5jdGlvbicgJiYgQWJvcnRTaWduYWwudGltZW91dCAmJiB0eXBlb2YgQWJvcnRTaWduYWwudGltZW91dCA9PT0gJ2Z1bmN0aW9uJykge1xyXG4gICAgICAgICAgICBBYm9ydFNpZ25hbC50aW1lb3V0KG1pbGxpc2Vjb25kcykuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSwgbWlsbGlzZWNvbmRzKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dC5cclxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBzbGVlcChtaWxsaXNlY29uZHM6IG51bWJlcik6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlclxyXG4gICAgICogdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCB2YWx1ZS5cclxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBzbGVlcDxUPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XHJcbiAgICBzdGF0aWMgc2xlZXA8VCA9IHZvaWQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZT86IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlKSA9PiB7XHJcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gcmVzb2x2ZSh2YWx1ZSEpLCBtaWxsaXNlY29uZHMpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZWplY3RlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCByZWFzb24uXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByZWplY3Q8VCA9IG5ldmVyPihyZWFzb24/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChfLCByZWplY3QpID0+IHJlamVjdChyZWFzb24pKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmVzb2x2ZSgpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xyXG4gICAgc3RhdGljIHJlc29sdmU8VCA9IHZvaWQ+KHZhbHVlPzogVCB8IFByb21pc2VMaWtlPFQ+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+IHtcclxuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpIHtcclxuICAgICAgICAgICAgLy8gT3B0aW1pc2UgZm9yIGNhbmNlbGxhYmxlIHByb21pc2VzLlxyXG4gICAgICAgICAgICByZXR1cm4gdmFsdWU7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPGFueT4oKHJlc29sdmUpID0+IHJlc29sdmUodmFsdWUpKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIGFuZCByZXR1cm5zIGl0IGluIGFuIG9iamVjdCwgYWxvbmcgd2l0aCBpdHMgcmVzb2x2ZSBhbmQgcmVqZWN0IGZ1bmN0aW9uc1xyXG4gICAgICogYW5kIGEgZ2V0dGVyL3NldHRlciBmb3IgdGhlIGNhbmNlbGxhdGlvbiBjYWxsYmFjay5cclxuICAgICAqXHJcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgd2l0aFJlc29sdmVyczxUPigpOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcclxuICAgICAgICBsZXQgcmVzdWx0OiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0geyBvbmNhbmNlbGxlZDogbnVsbCB9IGFzIGFueTtcclxuICAgICAgICByZXN1bHQucHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICByZXN1bHQucmVzb2x2ZSA9IHJlc29sdmU7XHJcbiAgICAgICAgICAgIHJlc3VsdC5yZWplY3QgPSByZWplY3Q7XHJcbiAgICAgICAgfSwgKGNhdXNlPzogYW55KSA9PiB7IHJlc3VsdC5vbmNhbmNlbGxlZD8uKGNhdXNlKTsgfSk7XHJcbiAgICAgICAgcmV0dXJuIHJlc3VsdDtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIGNhbmNlbGxhdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxyXG4gKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSByZXN1bHRpbmcgZnVuY3Rpb24gZG9lcyBub3QgcmVqZWN0LlxyXG4gKi9cclxuZnVuY3Rpb24gY2FuY2VsbGVyRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSkge1xyXG4gICAgbGV0IGNhbmNlbGxhdGlvblByb21pc2U6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9IHVuZGVmaW5lZDtcclxuXHJcbiAgICByZXR1cm4gKHJlYXNvbjogQ2FuY2VsRXJyb3IpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPT4ge1xyXG4gICAgICAgIGlmICghc3RhdGUuc2V0dGxlZCkge1xyXG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgc3RhdGUucmVhc29uID0gcmVhc29uO1xyXG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xyXG5cclxuICAgICAgICAgICAgLy8gQXR0YWNoIGFuIGVycm9yIGhhbmRsZXIgdGhhdCBpZ25vcmVzIHRoaXMgc3BlY2lmaWMgcmVqZWN0aW9uIHJlYXNvbiBhbmQgbm90aGluZyBlbHNlLlxyXG4gICAgICAgICAgICAvLyBJbiB0aGVvcnksIGEgc2FuZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uIGF0IHRoaXMgcG9pbnRcclxuICAgICAgICAgICAgLy8gc2hvdWxkIGFsd2F5cyByZWplY3Qgd2l0aCBvdXIgY2FuY2VsbGF0aW9uIHJlYXNvbixcclxuICAgICAgICAgICAgLy8gaGVuY2UgdGhlIGhhbmRsZXIgd2lsbCBuZXZlciB0aHJvdy5cclxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnByb3RvdHlwZS50aGVuLmNhbGwocHJvbWlzZS5wcm9taXNlLCB1bmRlZmluZWQsIChlcnIpID0+IHtcclxuICAgICAgICAgICAgICAgIGlmIChlcnIgIT09IHJlYXNvbikge1xyXG4gICAgICAgICAgICAgICAgICAgIHRocm93IGVycjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfSk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBJZiByZWFzb24gaXMgbm90IHNldCwgdGhlIHByb21pc2UgcmVzb2x2ZWQgcmVndWxhcmx5LCBoZW5jZSB3ZSBtdXN0IG5vdCBjYWxsIG9uY2FuY2VsbGVkLlxyXG4gICAgICAgIC8vIElmIG9uY2FuY2VsbGVkIGlzIHVuc2V0LCBubyBuZWVkIHRvIGdvIGFueSBmdXJ0aGVyLlxyXG4gICAgICAgIGlmICghc3RhdGUucmVhc29uIHx8ICFwcm9taXNlLm9uY2FuY2VsbGVkKSB7IHJldHVybjsgfVxyXG5cclxuICAgICAgICBjYW5jZWxsYXRpb25Qcm9taXNlID0gbmV3IFByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIHJlc29sdmUocHJvbWlzZS5vbmNhbmNlbGxlZCEoc3RhdGUucmVhc29uIS5jYXVzZSkpO1xyXG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KS5jYXRjaCgocmVhc29uPykgPT4ge1xyXG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xyXG4gICAgICAgIH0pO1xyXG5cclxuICAgICAgICAvLyBVbnNldCBvbmNhbmNlbGxlZCB0byBwcmV2ZW50IHJlcGVhdGVkIGNhbGxzLlxyXG4gICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBudWxsO1xyXG5cclxuICAgICAgICByZXR1cm4gY2FuY2VsbGF0aW9uUHJvbWlzZTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlc29sdXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cclxuICovXHJcbmZ1bmN0aW9uIHJlc29sdmVyRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+IHtcclxuICAgIHJldHVybiAodmFsdWUpID0+IHtcclxuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxyXG4gICAgICAgIHN0YXRlLnJlc29sdmluZyA9IHRydWU7XHJcblxyXG4gICAgICAgIGlmICh2YWx1ZSA9PT0gcHJvbWlzZS5wcm9taXNlKSB7XHJcbiAgICAgICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7IHJldHVybjsgfVxyXG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QobmV3IFR5cGVFcnJvcihcIkEgcHJvbWlzZSBjYW5ub3QgYmUgcmVzb2x2ZWQgd2l0aCBpdHNlbGYuXCIpKTtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKHZhbHVlICE9IG51bGwgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcgfHwgdHlwZW9mIHZhbHVlID09PSAnZnVuY3Rpb24nKSkge1xyXG4gICAgICAgICAgICBsZXQgdGhlbjogYW55O1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgdGhlbiA9ICh2YWx1ZSBhcyBhbnkpLnRoZW47XHJcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xyXG4gICAgICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgICAgICBwcm9taXNlLnJlamVjdChlcnIpO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICBpZiAoaXNDYWxsYWJsZSh0aGVuKSkge1xyXG4gICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICBsZXQgY2FuY2VsID0gKHZhbHVlIGFzIGFueSkuY2FuY2VsO1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKGNhbmNlbCkpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3Qgb25jYW5jZWxsZWQgPSAoY2F1c2U/OiBhbnkpID0+IHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH07XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChzdGF0ZS5yZWFzb24pIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIElmIGFscmVhZHkgY2FuY2VsbGVkLCBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uLlxyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgY2FuY2VsbGVyIGFsZ29yaXRobSBkb2VzIG5vdCByZWplY3RcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIHNvIGl0IGNhbiBiZSBkaXNjYXJkZWQgc2FmZWx5LlxyXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCBjYW5jZWxsZXJGb3IoeyAuLi5wcm9taXNlLCBvbmNhbmNlbGxlZCB9LCBzdGF0ZSkoc3RhdGUucmVhc29uKTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBvbmNhbmNlbGxlZDtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0gY2F0Y2gge31cclxuXHJcbiAgICAgICAgICAgICAgICBjb25zdCBuZXdTdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcm9vdDogc3RhdGUucm9vdCxcclxuICAgICAgICAgICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgICAgIGdldCBzZXR0bGVkKCkgeyByZXR1cm4gdGhpcy5yb290LnNldHRsZWQgfSxcclxuICAgICAgICAgICAgICAgICAgICBzZXQgc2V0dGxlZCh2YWx1ZSkgeyB0aGlzLnJvb3Quc2V0dGxlZCA9IHZhbHVlOyB9LFxyXG4gICAgICAgICAgICAgICAgICAgIGdldCByZWFzb24oKSB7IHJldHVybiB0aGlzLnJvb3QucmVhc29uIH1cclxuICAgICAgICAgICAgICAgIH07XHJcblxyXG4gICAgICAgICAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBuZXdTdGF0ZSk7XHJcbiAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkodGhlbiwgdmFsdWUsIFtyZXNvbHZlckZvcihwcm9taXNlLCBuZXdTdGF0ZSksIHJlamVjdG9yXSk7XHJcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmV0dXJuOyAvLyBJTVBPUlRBTlQhXHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7IHJldHVybjsgfVxyXG4gICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgIHByb21pc2UucmVzb2x2ZSh2YWx1ZSk7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVqZWN0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXHJcbiAqL1xyXG5mdW5jdGlvbiByZWplY3RvckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciB7XHJcbiAgICByZXR1cm4gKHJlYXNvbj8pID0+IHtcclxuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxyXG4gICAgICAgIHN0YXRlLnJlc29sdmluZyA9IHRydWU7XHJcblxyXG4gICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICBpZiAocmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgc3RhdGUucmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgT2JqZWN0LmlzKHJlYXNvbi5jYXVzZSwgc3RhdGUucmVhc29uLmNhdXNlKSkge1xyXG4gICAgICAgICAgICAgICAgICAgIC8vIFN3YWxsb3cgbGF0ZSByZWplY3Rpb25zIHRoYXQgYXJlIENhbmNlbEVycm9ycyB3aG9zZSBjYW5jZWxsYXRpb24gY2F1c2UgaXMgdGhlIHNhbWUgYXMgb3Vycy5cclxuICAgICAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0gY2F0Y2gge31cclxuXHJcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uKSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogQ2FuY2VscyBhbGwgdmFsdWVzIGluIGFuIGFycmF5IHRoYXQgbG9vayBsaWtlIGNhbmNlbGxhYmxlIHRoZW5hYmxlcy5cclxuICogUmV0dXJucyBhIHByb21pc2UgdGhhdCBmdWxmaWxscyBvbmNlIGFsbCBjYW5jZWxsYXRpb24gcHJvY2VkdXJlcyBmb3IgdGhlIGdpdmVuIHZhbHVlcyBoYXZlIHNldHRsZWQuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYW5jZWxBbGwocGFyZW50OiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHZhbHVlczogYW55W10sIGNhdXNlPzogYW55KTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICBjb25zdCByZXN1bHRzOiBQcm9taXNlPHZvaWQ+W10gPSBbXTtcclxuXHJcbiAgICBmb3IgKGNvbnN0IHZhbHVlIG9mIHZhbHVlcykge1xyXG4gICAgICAgIGxldCBjYW5jZWw6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcjtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUodmFsdWUudGhlbikpIHsgY29udGludWU7IH1cclxuICAgICAgICAgICAgY2FuY2VsID0gdmFsdWUuY2FuY2VsO1xyXG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUoY2FuY2VsKSkgeyBjb250aW51ZTsgfVxyXG4gICAgICAgIH0gY2F0Y2ggeyBjb250aW51ZTsgfVxyXG5cclxuICAgICAgICBsZXQgcmVzdWx0OiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgcmVzdWx0ID0gUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcclxuICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xyXG4gICAgICAgICAgICBjb250aW51ZTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICghcmVzdWx0KSB7IGNvbnRpbnVlOyB9XHJcbiAgICAgICAgcmVzdWx0cy5wdXNoKFxyXG4gICAgICAgICAgICAocmVzdWx0IGluc3RhbmNlb2YgUHJvbWlzZSAgPyByZXN1bHQgOiBQcm9taXNlLnJlc29sdmUocmVzdWx0KSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcclxuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICApO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBQcm9taXNlLmFsbChyZXN1bHRzKSBhcyBhbnk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGl0cyBhcmd1bWVudC5cclxuICovXHJcbmZ1bmN0aW9uIGlkZW50aXR5PFQ+KHg6IFQpOiBUIHtcclxuICAgIHJldHVybiB4O1xyXG59XHJcblxyXG4vKipcclxuICogVGhyb3dzIGl0cyBhcmd1bWVudC5cclxuICovXHJcbmZ1bmN0aW9uIHRocm93ZXIocmVhc29uPzogYW55KTogbmV2ZXIge1xyXG4gICAgdGhyb3cgcmVhc29uO1xyXG59XHJcblxyXG4vKipcclxuICogQXR0ZW1wdHMgdmFyaW91cyBzdHJhdGVnaWVzIHRvIGNvbnZlcnQgYW4gZXJyb3IgdG8gYSBzdHJpbmcuXHJcbiAqL1xyXG5mdW5jdGlvbiBlcnJvck1lc3NhZ2UoZXJyOiBhbnkpOiBzdHJpbmcge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICBpZiAoZXJyIGluc3RhbmNlb2YgRXJyb3IgfHwgdHlwZW9mIGVyciAhPT0gJ29iamVjdCcgfHwgZXJyLnRvU3RyaW5nICE9PSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nKSB7XHJcbiAgICAgICAgICAgIHJldHVybiBcIlwiICsgZXJyO1xyXG4gICAgICAgIH1cclxuICAgIH0gY2F0Y2gge31cclxuXHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiBKU09OLnN0cmluZ2lmeShlcnIpO1xyXG4gICAgfSBjYXRjaCB7fVxyXG5cclxuICAgIHRyeSB7XHJcbiAgICAgICAgcmV0dXJuIE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcuY2FsbChlcnIpO1xyXG4gICAgfSBjYXRjaCB7fVxyXG5cclxuICAgIHJldHVybiBcIjxjb3VsZCBub3QgY29udmVydCBlcnJvciB0byBzdHJpbmc+XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGJhcnJpZXIgcHJvbWlzZSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuIElmIG5lY2Vzc2FyeSwgaW5pdGlhbGlzZXMgdGhlIGJhcnJpZXIuXHJcbiAqL1xyXG5mdW5jdGlvbiBjdXJyZW50QmFycmllcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIGxldCBwd3I6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0gcHJvbWlzZVtiYXJyaWVyU3ltXSA/PyB7fTtcclxuICAgIGlmICghKCdwcm9taXNlJyBpbiBwd3IpKSB7XHJcbiAgICAgICAgT2JqZWN0LmFzc2lnbihwd3IsIHByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+KCkpO1xyXG4gICAgfVxyXG4gICAgaWYgKHByb21pc2VbYmFycmllclN5bV0gPT0gbnVsbCkge1xyXG4gICAgICAgIHB3ci5yZXNvbHZlISgpO1xyXG4gICAgICAgIHByb21pc2VbYmFycmllclN5bV0gPSBwd3I7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gcHdyLnByb21pc2UhO1xyXG59XHJcblxyXG4vLyBQb2x5ZmlsbCBQcm9taXNlLndpdGhSZXNvbHZlcnMuXHJcbmxldCBwcm9taXNlV2l0aFJlc29sdmVycyA9IFByb21pc2Uud2l0aFJlc29sdmVycztcclxuaWYgKHByb21pc2VXaXRoUmVzb2x2ZXJzICYmIHR5cGVvZiBwcm9taXNlV2l0aFJlc29sdmVycyA9PT0gJ2Z1bmN0aW9uJykge1xyXG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBwcm9taXNlV2l0aFJlc29sdmVycy5iaW5kKFByb21pc2UpO1xyXG59IGVsc2Uge1xyXG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBmdW5jdGlvbiA8VD4oKTogUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xyXG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XHJcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBQcm9taXNlPFQ+KChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xyXG4gICAgICAgIHJldHVybiB7IHByb21pc2UsIHJlc29sdmUsIHJlamVjdCB9O1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DbGlwYm9hcmQpO1xyXG5cclxuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XHJcbmNvbnN0IENsaXBib2FyZFRleHQgPSAxO1xyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHRleHQgdG8gdGhlIENsaXBib2FyZC5cclxuICpcclxuICogQHBhcmFtIHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cclxuICogQHJldHVybiBBIFByb21pc2UgdGhhdCByZXNvbHZlcyB3aGVuIHRoZSBvcGVyYXRpb24gaXMgc3VjY2Vzc2Z1bC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQ6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpOiBQcm9taXNlPHN0cmluZz4ge1xyXG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkVGV4dCk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgU2l6ZSB7XHJcbiAgICAvKiogVGhlIHdpZHRoIG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cclxuICAgIFdpZHRoOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGhlaWdodCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXHJcbiAgICBIZWlnaHQ6IG51bWJlcjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBSZWN0IHtcclxuICAgIC8qKiBUaGUgWCBjb29yZGluYXRlIG9mIHRoZSBvcmlnaW4uICovXHJcbiAgICBYOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xyXG4gICAgWTogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xyXG4gICAgV2lkdGg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIHRoZSByZWN0YW5nbGUuICovXHJcbiAgICBIZWlnaHQ6IG51bWJlcjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBTY3JlZW4ge1xyXG4gICAgLyoqIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLiAqL1xyXG4gICAgSUQ6IHN0cmluZztcclxuICAgIC8qKiBIdW1hbi1yZWFkYWJsZSBuYW1lIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBOYW1lOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIHNjYWxlIGZhY3RvciBvZiB0aGUgc2NyZWVuIChEUEkvOTYpLiAxID0gc3RhbmRhcmQgRFBJLCAyID0gSGlEUEkgKFJldGluYSksIGV0Yy4gKi9cclxuICAgIFNjYWxlRmFjdG9yOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgWDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFk6IG51bWJlcjtcclxuICAgIC8qKiBDb250YWlucyB0aGUgd2lkdGggYW5kIGhlaWdodCBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgU2l6ZTogU2l6ZTtcclxuICAgIC8qKiBDb250YWlucyB0aGUgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQuICovXHJcbiAgICBCb3VuZHM6IFJlY3Q7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0IChiZWZvcmUgc2NhbGluZykuICovXHJcbiAgICBQaHlzaWNhbEJvdW5kczogUmVjdDtcclxuICAgIC8qKiBDb250YWlucyB0aGUgYXJlYSBvZiB0aGUgc2NyZWVuIHRoYXQgaXMgYWN0dWFsbHkgdXNhYmxlIChleGNsdWRpbmcgdGFza2JhciBhbmQgb3RoZXIgc3lzdGVtIFVJKS4gKi9cclxuICAgIFdvcmtBcmVhOiBSZWN0O1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBXb3JrQXJlYSBvZiB0aGUgc2NyZWVuIChiZWZvcmUgc2NhbGluZykuICovXHJcbiAgICBQaHlzaWNhbFdvcmtBcmVhOiBSZWN0O1xyXG4gICAgLyoqIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xyXG4gICAgSXNQcmltYXJ5OiBib29sZWFuO1xyXG4gICAgLyoqIFRoZSByb3RhdGlvbiBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgUm90YXRpb246IG51bWJlcjtcclxufVxyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlNjcmVlbnMpO1xyXG5cclxuY29uc3QgZ2V0QWxsID0gMDtcclxuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XHJcbmNvbnN0IGdldEN1cnJlbnQgPSAyO1xyXG5jb25zdCBnZXRCeUlEID0gMztcclxuY29uc3QgZ2V0QnlJbmRleCA9IDQ7XHJcblxyXG4vKipcclxuICogR2V0cyBhbGwgc2NyZWVucy5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gYXJyYXkgb2YgU2NyZWVuIG9iamVjdHMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCk6IFByb21pc2U8U2NyZWVuW10+IHtcclxuICAgIHJldHVybiBjYWxsKGdldEFsbCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKTogUHJvbWlzZTxTY3JlZW4+IHtcclxuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0Q3VycmVudCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIGEgc2NyZWVuIGJ5IGl0cyB1bmlxdWUgZGlzcGxheSBJRC5cclxuICpcclxuICogQHBhcmFtIGlkIC0gVGhlIHVuaXF1ZSBpZGVudGlmaWVyIG9mIHRoZSBzY3JlZW4uXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBtYXRjaGluZyBTY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0QnlJRChpZDogc3RyaW5nKTogUHJvbWlzZTxTY3JlZW4+IHtcclxuICAgIHJldHVybiBjYWxsKGdldEJ5SUQsIHsgaWQgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIGEgc2NyZWVuIGJ5IGl0cyBpbmRleCBpbiB0aGUgc2NyZWVuIGxpc3QuXHJcbiAqXHJcbiAqIEBwYXJhbSBpbmRleCAtIFRoZSB6ZXJvLWJhc2VkIGluZGV4IG9mIHRoZSBzY3JlZW4uXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBtYXRjaGluZyBTY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0QnlJbmRleChpbmRleDogbnVtYmVyKTogUHJvbWlzZTxTY3JlZW4+IHtcclxuICAgIHJldHVybiBjYWxsKGdldEJ5SW5kZXgsIHsgaW5kZXggfSk7XHJcbn1cclxuIiwgIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuSU9TKTtcclxuXHJcbi8vIE1ldGhvZCBJRHNcclxuY29uc3QgSGFwdGljc0ltcGFjdCA9IDA7XHJcbmNvbnN0IERldmljZUluZm8gPSAxO1xyXG5cclxuZXhwb3J0IG5hbWVzcGFjZSBIYXB0aWNzIHtcclxuICAgIGV4cG9ydCB0eXBlIEltcGFjdFN0eWxlID0gXCJsaWdodFwifFwibWVkaXVtXCJ8XCJoZWF2eVwifFwic29mdFwifFwicmlnaWRcIjtcclxuICAgIGV4cG9ydCBmdW5jdGlvbiBJbXBhY3Qoc3R5bGU6IEltcGFjdFN0eWxlID0gXCJtZWRpdW1cIik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKEhhcHRpY3NJbXBhY3QsIHsgc3R5bGUgfSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgRGV2aWNlIHtcclxuICAgIGV4cG9ydCBpbnRlcmZhY2UgSW5mbyB7XHJcbiAgICAgICAgbW9kZWw6IHN0cmluZztcclxuICAgICAgICBzeXN0ZW1OYW1lOiBzdHJpbmc7XHJcbiAgICAgICAgc3lzdGVtVmVyc2lvbjogc3RyaW5nO1xyXG4gICAgICAgIGlzU2ltdWxhdG9yOiBib29sZWFuO1xyXG4gICAgfVxyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIEluZm8oKTogUHJvbWlzZTxJbmZvPiB7XHJcbiAgICAgICAgcmV0dXJuIGNhbGwoRGV2aWNlSW5mbyk7XHJcbiAgICB9XHJcbn1cclxuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FB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Tyx5QkFBeUIsZUFBTyx1QkFBdUIsS0FBSyxjQUFNO0FBR2hGLE9BQU8sT0FBTyxrQkFBa0I7QUFDaEMsT0FBTyxPQUFPLGtCQUFrQjtBQUNoQyxPQUFPLE9BQU8saUJBQWlCO0FBRXhCLE9BQU8scUJBQXFCO0FBTzVCLFNBQVMsbUJBQW1CLEtBQTRCO0FBQzNELFNBQU8sTUFBTSxLQUFLLEVBQUUsUUFBUSxPQUFPLENBQUMsRUFDL0IsS0FBSyxjQUFZO0FBQ2QsUUFBSSxTQUFTLElBQUk7QUFHYixZQUFNLGVBQWUsU0FBUyxRQUFRLElBQUksY0FBYyxLQUFLLElBQUksWUFBWTtBQUM3RSxVQUFJLFlBQVksU0FBUyxZQUFZLEdBQUc7QUFDcEMsY0FBTSxTQUFTLFNBQVMsY0FBYyxRQUFRO0FBQzlDLGVBQU8sTUFBTTtBQUNiLGlCQUFTLEtBQUssWUFBWSxNQUFNO0FBQUEsTUFDcEM7QUFBQSxJQUNKO0FBQUEsRUFDSixDQUFDLEVBQ0EsTUFBTSxNQUFNO0FBQUEsRUFBQyxDQUFDO0FBQ3ZCO0FBR0EsbUJBQW1CLGtCQUFrQjsiLAogICJuYW1lcyI6IFsiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJfYSIsICJBcnJheSIsICJNYXAiLCAiQXJyYXkiLCAiTWFwIiwgImtleSIsICJjYWxsIiwgIl9hIiwgIl9hIiwgInJlc2l6YWJsZSIsICJfYSIsICJjYWxsIiwgIl9hIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJIaWRlTWV0aG9kIiwgIlNob3dNZXRob2QiLCAiaXNEb2N1bWVudERvdEFsbCIsICJfYSIsICJyZWFzb24iLCAidmFsdWUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIl0KfQo=
