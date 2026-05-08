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
var __export = (target, all2) => {
  for (var name in all2)
    __defProp(target, name, { get: all2[name], enumerable: true });
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
  Updater: () => updater_exports,
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
    ApplicationDidWake: "mac:ApplicationDidWake",
    ApplicationScreensDidSleep: "mac:ApplicationScreensDidSleep",
    ApplicationScreensDidWake: "mac:ApplicationScreensDidWake",
    ApplicationShouldHandleReopen: "mac:ApplicationShouldHandleReopen",
    ApplicationWillBecomeActive: "mac:ApplicationWillBecomeActive",
    ApplicationWillFinishLaunching: "mac:ApplicationWillFinishLaunching",
    ApplicationWillHide: "mac:ApplicationWillHide",
    ApplicationWillResignActive: "mac:ApplicationWillResignActive",
    ApplicationWillSleep: "mac:ApplicationWillSleep",
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
    SystemDidWake: "linux:SystemDidWake",
    SystemThemeChanged: "linux:SystemThemeChanged",
    SystemWillSleep: "linux:SystemWillSleep",
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
    SystemDidWake: "common:SystemDidWake",
    SystemWillSleep: "common:SystemWillSleep",
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
  if (!resizable || !IsWindows() && !IsLinux()) {
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

// desktop/@wailsio/runtime/src/updater.ts
var updater_exports = {};
__export(updater_exports, {
  Events: () => Events2
});
var Events2 = Object.freeze({
  /** A Check round-trip is starting. Payload: null. */
  CheckStarted: "wails:updater:check-started",
  /** Check found a newer release. Payload: Release. */
  UpdateAvailable: "wails:updater:update-available",
  /** Check confirmed the caller is up to date. Payload: null. */
  NoUpdate: "wails:updater:no-update",
  /** Download is starting. Payload: Release. */
  DownloadStarted: "wails:updater:download-started",
  /** Periodic progress tick during download (~10 Hz). Payload: Progress. */
  DownloadProgress: "wails:updater:download-progress",
  /** All bytes are on disk, but verification has not yet started. Payload: Release. */
  DownloadComplete: "wails:updater:download-complete",
  /** Signature / digest verification has started. Payload: Release. */
  Verifying: "wails:updater:verifying",
  /** The Updater is swapping the binary into place. Payload: Release. */
  Installing: "wails:updater:installing",
  /** Update is staged and a restart is pending. Payload: Release. */
  UpdateReady: "wails:updater:update-ready",
  /** Something failed. Payload: ErrorInfo { stage, message, provider }. */
  Error: "wails:updater:error",
  /** Host-side context delivered once per session. Payload: Meta { currentVersion, skippedVersion }. */
  Meta: "wails:updater:meta",
  /** Sub-namespace: user-action events that the UI emits BACK to the host. */
  User: Object.freeze({
    /** User clicked Install on an Available update. */
    Install: "wails:updater:user:install",
    /** User clicked Restart & Apply on a Ready update. */
    Restart: "wails:updater:user:restart",
    /** User clicked Skip This Version. */
    Skip: "wails:updater:user:skip",
    /** User clicked Remind Me Later. */
    Remind: "wails:updater:user:remind",
    /** User clicked Close / Cancel. */
    Cancel: "wails:updater:user:cancel"
  }),
  /** Sub-namespace: framework-internal events the UI emits to coordinate
   *  with the host. Most app code can ignore these. */
  Window: Object.freeze({
    /** The window finished loading and asks the host to replay current state. */
    Ready: "wails:updater:window:ready"
  })
});

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
  updater_exports as Updater,
  wml_exports as WML,
  window_default as Window,
  clientId,
  getTransport,
  loadOptionalScript,
  objectNames,
  setTransport
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY3JlYXRlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudF90eXBlcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbGFibGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbmNlbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2lvcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXBkYXRlci50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxuaW1wb3J0IFwiLi9jb250ZXh0bWVudS5qc1wiO1xuaW1wb3J0IFwiLi9kcmFnLmpzXCI7XG5cbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi9hcHBsaWNhdGlvbi5qc1wiO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyLmpzXCI7XG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzLmpzXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkLmpzXCI7XG5pbXBvcnQgKiBhcyBDcmVhdGUgZnJvbSBcIi4vY3JlYXRlLmpzXCI7XG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcbmltcG9ydCAqIGFzIEV2ZW50cyBmcm9tIFwiLi9ldmVudHMuanNcIjtcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuL2ZsYWdzLmpzXCI7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnMuanNcIjtcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCAqIGFzIElPUyBmcm9tIFwiLi9pb3MuanNcIjtcbmltcG9ydCAqIGFzIFVwZGF0ZXIgZnJvbSBcIi4vdXBkYXRlci5qc1wiO1xuaW1wb3J0IFdpbmRvdywgeyBoYW5kbGVEcmFnRW50ZXIsIGhhbmRsZURyYWdMZWF2ZSwgaGFuZGxlRHJhZ092ZXIgfSBmcm9tIFwiLi93aW5kb3cuanNcIjtcbmltcG9ydCAqIGFzIFdNTCBmcm9tIFwiLi93bWwuanNcIjtcblxuZXhwb3J0IHtcbiAgICBBcHBsaWNhdGlvbixcbiAgICBCcm93c2VyLFxuICAgIENhbGwsXG4gICAgQ2xpcGJvYXJkLFxuICAgIERpYWxvZ3MsXG4gICAgRXZlbnRzLFxuICAgIEZsYWdzLFxuICAgIFNjcmVlbnMsXG4gICAgU3lzdGVtLFxuICAgIElPUyxcbiAgICBVcGRhdGVyLFxuICAgIFdpbmRvdyxcbiAgICBXTUxcbn07XG5cbi8qKlxuICogQW4gaW50ZXJuYWwgdXRpbGl0eSBjb25zdW1lZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXG4gKlxuICogQGlnbm9yZVxuICovXG5leHBvcnQgeyBDcmVhdGUgfTtcblxuZXhwb3J0ICogZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcblxuLy8gRXhwb3J0IHRyYW5zcG9ydCBpbnRlcmZhY2VzIGFuZCB1dGlsaXRpZXNcbmV4cG9ydCB7XG4gICAgc2V0VHJhbnNwb3J0LFxuICAgIGdldFRyYW5zcG9ydCxcbiAgICB0eXBlIFJ1bnRpbWVUcmFuc3BvcnQsXG4gICAgb2JqZWN0TmFtZXMsXG4gICAgY2xpZW50SWQsXG59IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuaW1wb3J0IHsgY2xpZW50SWQgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbi8vIE5vdGlmeSBiYWNrZW5kXG53aW5kb3cuX3dhaWxzLmludm9rZSA9IFN5c3RlbS5pbnZva2U7XG53aW5kb3cuX3dhaWxzLmNsaWVudElkID0gY2xpZW50SWQ7XG5cbi8vIFJlZ2lzdGVyIHBsYXRmb3JtIGhhbmRsZXJzIChpbnRlcm5hbCBBUEkpXG4vLyBOb3RlOiBXaW5kb3cgaXMgdGhlIHRoaXNXaW5kb3cgaW5zdGFuY2UgKGRlZmF1bHQgZXhwb3J0IGZyb20gd2luZG93LnRzKVxuLy8gQmluZGluZyBlbnN1cmVzICd0aGlzJyBjb3JyZWN0bHkgcmVmZXJzIHRvIHRoZSBjdXJyZW50IHdpbmRvdyBpbnN0YW5jZVxud2luZG93Ll93YWlscy5oYW5kbGVQbGF0Zm9ybUZpbGVEcm9wID0gV2luZG93LkhhbmRsZVBsYXRmb3JtRmlsZURyb3AuYmluZChXaW5kb3cpO1xuXG4vLyBMaW51eC1zcGVjaWZpYyBkcmFnIGhhbmRsZXJzIChHVEsgaW50ZXJjZXB0cyBET00gZHJhZyBldmVudHMpXG53aW5kb3cuX3dhaWxzLmhhbmRsZURyYWdFbnRlciA9IGhhbmRsZURyYWdFbnRlcjtcbndpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ0xlYXZlID0gaGFuZGxlRHJhZ0xlYXZlO1xud2luZG93Ll93YWlscy5oYW5kbGVEcmFnT3ZlciA9IGhhbmRsZURyYWdPdmVyO1xuXG5TeXN0ZW0uaW52b2tlKFwid2FpbHM6cnVudGltZTpyZWFkeVwiKTtcblxuLyoqXG4gKiBMb2FkcyBhIHNjcmlwdCBmcm9tIHRoZSBnaXZlbiBVUkwgaWYgaXQgZXhpc3RzLlxuICogVXNlcyBIRUFEIHJlcXVlc3QgdG8gY2hlY2sgZXhpc3RlbmNlLCB0aGVuIGluamVjdHMgYSBzY3JpcHQgdGFnLlxuICogU2lsZW50bHkgaWdub3JlcyBpZiB0aGUgc2NyaXB0IGRvZXNuJ3QgZXhpc3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBsb2FkT3B0aW9uYWxTY3JpcHQodXJsOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gZmV0Y2godXJsLCB7IG1ldGhvZDogJ0hFQUQnIH0pXG4gICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgIC8vIFZlcmlmeSB0aGUgcmVzcG9uc2UgaXMgYWN0dWFsbHkgSmF2YVNjcmlwdCBhbmQgbm90IGFuIEhUTUwgZmFsbGJhY2tcbiAgICAgICAgICAgICAgICAvLyAoZS5nLiBWaXRlIGRldiBzZXJ2ZXIgcmV0dXJucyBpbmRleC5odG1sIGZvciB1bmtub3duIHJvdXRlcylcbiAgICAgICAgICAgICAgICBjb25zdCBjb250ZW50VHlwZSA9IChyZXNwb25zZS5oZWFkZXJzLmdldCgnY29udGVudC10eXBlJykgfHwgJycpLnRvTG93ZXJDYXNlKCk7XG4gICAgICAgICAgICAgICAgaWYgKGNvbnRlbnRUeXBlLmluY2x1ZGVzKCdqYXZhc2NyaXB0JykpIHtcbiAgICAgICAgICAgICAgICAgICAgY29uc3Qgc2NyaXB0ID0gZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgnc2NyaXB0Jyk7XG4gICAgICAgICAgICAgICAgICAgIHNjcmlwdC5zcmMgPSB1cmw7XG4gICAgICAgICAgICAgICAgICAgIGRvY3VtZW50LmhlYWQuYXBwZW5kQ2hpbGQoc2NyaXB0KTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pXG4gICAgICAgIC5jYXRjaCgoKSA9PiB7fSk7IC8vIFNpbGVudGx5IGlnbm9yZSAtIHNjcmlwdCBpcyBvcHRpb25hbFxufVxuXG4vLyBMb2FkIGN1c3RvbS5qcyBpZiBhdmFpbGFibGUgKHVzZWQgYnkgc2VydmVyIG1vZGUgZm9yIFdlYlNvY2tldCBldmVudHMsIGV0Yy4pXG5sb2FkT3B0aW9uYWxTY3JpcHQoJy93YWlscy9jdXN0b20uanMnKTtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCB7IFF1ZXN0aW9uIH0gZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0IHsgRW1pdCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXG4gKiBAcGFyYW0gW2RhdGE9bnVsbF0gLSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cbiAqL1xuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XG4gICAgRW1pdChldmVudE5hbWUsIGRhdGEpO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgc3BlY2lmaWVkIHdpbmRvdy5cbiAqXG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKi9cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TmFtZTogc3RyaW5nLCBtZXRob2ROYW1lOiBzdHJpbmcpIHtcbiAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBXaW5kb3cuR2V0KHdpbmRvd05hbWUpO1xuICAgIGNvbnN0IG1ldGhvZCA9ICh0YXJnZXRXaW5kb3cgYXMgYW55KVttZXRob2ROYW1lXTtcblxuICAgIGlmICh0eXBlb2YgbWV0aG9kICE9PSBcImZ1bmN0aW9uXCIpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICB0cnkge1xuICAgICAgICBtZXRob2QuY2FsbCh0YXJnZXRXaW5kb3cpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JzogYCwgZSk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXQuXG4gKi9cbmZ1bmN0aW9uIG9uV01MVHJpZ2dlcmVkKGV2OiBFdmVudCk6IHZvaWQge1xuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0IGFzIEVsZW1lbnQ7XG5cbiAgICBmdW5jdGlvbiBydW5FZmZlY3QoY2hvaWNlID0gXCJZZXNcIikge1xuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgdGFyZ2V0V2luZG93ID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10YXJnZXQtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRhcmdldC13aW5kb3cnKSB8fCBcIlwiO1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC13aW5kb3cnKTtcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW51cmwnKTtcblxuICAgICAgICBpZiAoZXZlbnRUeXBlICE9PSBudWxsKVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIGlmICh3aW5kb3dNZXRob2QgIT09IG51bGwpXG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHRhcmdldFdpbmRvdywgd2luZG93TWV0aG9kKTtcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcbiAgICAgICAgICAgIHZvaWQgT3BlblVSTCh1cmwpO1xuICAgIH1cblxuICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuXG4gICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgUXVlc3Rpb24oe1xuICAgICAgICAgICAgVGl0bGU6IFwiQ29uZmlybVwiLFxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcbiAgICAgICAgICAgIERldGFjaGVkOiBmYWxzZSxcbiAgICAgICAgICAgIEJ1dHRvbnM6IFtcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJOb1wiLCBJc0RlZmF1bHQ6IHRydWUgfVxuICAgICAgICAgICAgXVxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcnVuRWZmZWN0KCk7XG4gICAgfVxufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY29udHJvbGxlclN5bSA9IFN5bWJvbChcImNvbnRyb2xsZXJcIik7XG5jb25zdCB0cmlnZ2VyTWFwU3ltID0gU3ltYm9sKFwidHJpZ2dlck1hcFwiKTtcbmNvbnN0IGVsZW1lbnRDb3VudFN5bSA9IFN5bWJvbChcImVsZW1lbnRDb3VudFwiKTtcblxuLyoqXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXG4gKiBpdCB0aWVzIHRoZW0gdG8gYW4gQWJvcnRTaWduYWwgYW5kIHVzZXMgYW4gQWJvcnRDb250cm9sbGVyIHRvIHJlbW92ZSB0aGVtIGFsbCBhdCBvbmNlLlxuICovXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgW2NvbnRyb2xsZXJTeW1dOiBBYm9ydENvbnRyb2xsZXI7XG5cbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIGFuIG9wdGlvbnMgb2JqZWN0IGZvciBhZGRFdmVudExpc3RlbmVyIHRoYXQgdGllcyB0aGUgbGlzdGVuZXJcbiAgICAgKiB0byB0aGUgQWJvcnRTaWduYWwgZnJvbSB0aGUgY3VycmVudCBBYm9ydENvbnRyb2xsZXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxuICAgICAqIEBwYXJhbSB0cmlnZ2VycyAtIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudHNcbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xuICAgICAgICByZXR1cm4geyBzaWduYWw6IHRoaXNbY29udHJvbGxlclN5bV0uc2lnbmFsIH07XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMgYW5kIHJlc2V0cyB0aGUgcmVnaXN0cnkuXG4gICAgICovXG4gICAgcmVzZXQoKTogdm9pZCB7XG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0uYWJvcnQoKTtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG59XG5cbi8qKlxuICogV2Vha01hcFJlZ2lzdHJ5IG1hcHMgYWN0aXZlIHRyaWdnZXIgZXZlbnRzIHRvIGVhY2ggRE9NIGVsZW1lbnQgdGhyb3VnaCBhIFdlYWtNYXAuXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcbiAqIGNvbGxlY3Rpb24gb2YgdGhlIGludm9sdmVkIGVsZW1lbnRzLlxuICovXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xuICAgIC8qKiBTdG9yZXMgdGhlIGN1cnJlbnQgZWxlbWVudC10by10cmlnZ2VyIG1hcHBpbmcuICovXG4gICAgW3RyaWdnZXJNYXBTeW1dOiBXZWFrTWFwPEVsZW1lbnQsIHN0cmluZ1tdPjtcbiAgICAvKiogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLiAqL1xuICAgIFtlbGVtZW50Q291bnRTeW1dOiBudW1iZXI7XG5cbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50U3ltXSA9IDA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XG4gICAgICovXG4gICAgc2V0KGVsZW1lbnQ6IEVsZW1lbnQsIHRyaWdnZXJzOiBzdHJpbmdbXSk6IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zIHtcbiAgICAgICAgaWYgKCF0aGlzW3RyaWdnZXJNYXBTeW1dLmhhcyhlbGVtZW50KSkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0rKzsgfVxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dLnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgICAgIHJldHVybiB7fTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cbiAgICAgKi9cbiAgICByZXNldCgpOiB2b2lkIHtcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xuICAgICAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxuICAgICAgICAgICAgICAgIGJyZWFrO1xuXG4gICAgICAgICAgICBjb25zdCB0cmlnZ2VycyA9IHRoaXNbdHJpZ2dlck1hcFN5bV0uZ2V0KGVsZW1lbnQpO1xuICAgICAgICAgICAgaWYgKHRyaWdnZXJzICE9IG51bGwpIHsgdGhpc1tlbGVtZW50Q291bnRTeW1dLS07IH1cblxuICAgICAgICAgICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzIHx8IFtdKVxuICAgICAgICAgICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCk7XG4gICAgICAgIH1cblxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcbiAgICB9XG59XG5cbmNvbnN0IHRyaWdnZXJSZWdpc3RyeSA9IGNhbkFib3J0TGlzdGVuZXJzKCkgPyBuZXcgQWJvcnRDb250cm9sbGVyUmVnaXN0cnkoKSA6IG5ldyBXZWFrTWFwUmVnaXN0cnkoKTtcblxuLyoqXG4gKiBBZGRzIGV2ZW50IGxpc3RlbmVycyB0byB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXG4gKi9cbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50OiBFbGVtZW50KTogdm9pZCB7XG4gICAgY29uc3QgdHJpZ2dlclJlZ0V4cCA9IC9cXFMrL2c7XG4gICAgY29uc3QgdHJpZ2dlckF0dHIgPSAoZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCIpO1xuICAgIGNvbnN0IHRyaWdnZXJzOiBzdHJpbmdbXSA9IFtdO1xuXG4gICAgbGV0IG1hdGNoO1xuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcbiAgICAgICAgdHJpZ2dlcnMucHVzaChtYXRjaFswXSk7XG5cbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzKVxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW5hYmxlKCk6IHZvaWQge1xuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xufVxuXG4vKipcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpOiB2b2lkIHtcbiAgICB0cmlnZ2VyUmVnaXN0cnkucmVzZXQoKTtcbiAgICBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtZXZlbnRdLCBbd21sLXdpbmRvd10sIFt3bWwtb3BlbnVybF0sIFtkYXRhLXdtbC1ldmVudF0sIFtkYXRhLXdtbC13aW5kb3ddLCBbZGF0YS13bWwtb3BlbnVybF0nKS5mb3JFYWNoKGFkZFdNTExpc3RlbmVycyk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5Ccm93c2VyKTtcblxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xuXG4vKipcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMLlxuICpcbiAqIEBwYXJhbSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5VUkwodXJsOiBzdHJpbmcgfCBVUkwpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChCcm93c2VyT3BlblVSTCwge3VybDogdXJsLnRvU3RyaW5nKCl9KTtcbn1cbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9haS9uYW5vaWRcblxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXG4vL1xuLy8gQ29weXJpZ2h0IDIwMTcgQW5kcmV5IFNpdG5payA8YW5kcmV5QHNpdG5pay5ydT5cbi8vXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mXG4vLyB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsIGluXG4vLyB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvXG4vLyB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsIGNvcGllcyBvZlxuLy8gdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMgZnVybmlzaGVkIHRvIGRvIHNvLFxuLy8gICAgIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxuLy9cbi8vICAgICBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXG4vL1xuLy8gICAgIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1Jcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTXG4vLyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1Jcbi8vIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSIExJQUJJTElUWSwgV0hFVEhFUlxuLy8gSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU5cbi8vIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuXG5cbi8vIFRoaXMgYWxwaGFiZXQgdXNlcyBgQS1aYS16MC05Xy1gIHN5bWJvbHMuXG4vLyBUaGUgb3JkZXIgb2YgY2hhcmFjdGVycyBpcyBvcHRpbWl6ZWQgZm9yIGJldHRlciBnemlwIGFuZCBicm90bGkgY29tcHJlc3Npb24uXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBzYW1lIGZpbGUgKHdvcmtzIGJvdGggZm9yIGd6aXAgYW5kIGJyb3RsaSk6XG4vLyBgJ3VzZWAsIGBhbmRvbWAsIGFuZCBgcmljdCdgXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBicm90bGkgZGVmYXVsdCBkaWN0aW9uYXJ5OlxuLy8gYC0yNlRgLCBgMTk4M2AsIGA0MHB4YCwgYDc1cHhgLCBgYnVzaGAsIGBqYWNrYCwgYG1pbmRgLCBgdmVyeWAsIGFuZCBgd29sZmBcbmNvbnN0IHVybEFscGhhYmV0ID1cbiAgICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcblxuZXhwb3J0IGZ1bmN0aW9uIG5hbm9pZChzaXplOiBudW1iZXIgPSAyMSk6IHN0cmluZyB7XG4gICAgbGV0IGlkID0gJydcbiAgICAvLyBBIGNvbXBhY3QgYWx0ZXJuYXRpdmUgZm9yIGBmb3IgKHZhciBpID0gMDsgaSA8IHN0ZXA7IGkrKylgLlxuICAgIGxldCBpID0gc2l6ZSB8IDBcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICAgIC8vIGB8IDBgIGlzIG1vcmUgY29tcGFjdCBhbmQgZmFzdGVyIHRoYW4gYE1hdGguZmxvb3IoKWAuXG4gICAgICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG59XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gXCIuL25hbm9pZC5qc1wiO1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuLy8gU3RheSB1bmRlciBXZWJWaWV3MidzIH4yTUIgcmVxdWVzdCBib2R5IGJ1ZmZlcmluZyBsaW1pdCBpbiBXZWJSZXNvdXJjZVJlcXVlc3RlZC5cbmNvbnN0IENIVU5LX1RIUkVTSE9MRCA9IDUxMiAqIDEwMjQ7XG5cbi8vIFJlLWV4cG9ydCBuYW5vaWQgZm9yIGN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zXG5leHBvcnQgeyBuYW5vaWQgfTtcblxuLy8gT2JqZWN0IE5hbWVzXG5leHBvcnQgY29uc3Qgb2JqZWN0TmFtZXMgPSBPYmplY3QuZnJlZXplKHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxuICAgIENhbmNlbENhbGw6IDEwLFxuICAgIElPUzogMTEsXG59KTtcbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcblxuLyoqXG4gKiBSdW50aW1lVHJhbnNwb3J0IGRlZmluZXMgdGhlIGludGVyZmFjZSBmb3IgY3VzdG9tIElQQyB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zLlxuICogSW1wbGVtZW50IHRoaXMgaW50ZXJmYWNlIHRvIHVzZSBXZWJTb2NrZXRzLCBjdXN0b20gcHJvdG9jb2xzLCBvciBhbnkgb3RoZXJcbiAqIHRyYW5zcG9ydCBtZWNoYW5pc20gaW5zdGVhZCBvZiB0aGUgZGVmYXVsdCBIVFRQIGZldGNoLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIFJ1bnRpbWVUcmFuc3BvcnQge1xuICAgIC8qKlxuICAgICAqIFNlbmQgYSBydW50aW1lIGNhbGwgYW5kIHJldHVybiB0aGUgcmVzcG9uc2UuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gb2JqZWN0SUQgLSBUaGUgV2FpbHMgb2JqZWN0IElEICgwPUNhbGwsIDE9Q2xpcGJvYXJkLCBldGMuKVxuICAgICAqIEBwYXJhbSBtZXRob2QgLSBUaGUgbWV0aG9kIElEIHRvIGNhbGxcbiAgICAgKiBAcGFyYW0gd2luZG93TmFtZSAtIE9wdGlvbmFsIHdpbmRvdyBuYW1lXG4gICAgICogQHBhcmFtIGFyZ3MgLSBBcmd1bWVudHMgdG8gcGFzcyAod2lsbCBiZSBKU09OIHN0cmluZ2lmaWVkIGlmIHByZXNlbnQpXG4gICAgICogQHJldHVybnMgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHJlc3BvbnNlIGRhdGFcbiAgICAgKi9cbiAgICBjYWxsKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55Pjtcbn1cblxuLyoqXG4gKiBDdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uIChjYW4gYmUgc2V0IGJ5IHVzZXIpXG4gKi9cbmxldCBjdXN0b21UcmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsID0gbnVsbDtcblxuLyoqXG4gKiBTZXQgYSBjdXN0b20gdHJhbnNwb3J0IGZvciBhbGwgV2FpbHMgcnVudGltZSBjYWxscy5cbiAqIFRoaXMgYWxsb3dzIHlvdSB0byByZXBsYWNlIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0IHdpdGhcbiAqIFdlYlNvY2tldHMsIGN1c3RvbSBwcm90b2NvbHMsIG9yIGFueSBvdGhlciBtZWNoYW5pc20uXG4gKlxuICogQHBhcmFtIHRyYW5zcG9ydCAtIFlvdXIgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvblxuICpcbiAqIEBleGFtcGxlXG4gKiBgYGB0eXBlc2NyaXB0XG4gKiBpbXBvcnQgeyBzZXRUcmFuc3BvcnQgfSBmcm9tICcvd2FpbHMvcnVudGltZS5qcyc7XG4gKlxuICogY29uc3Qgd3NUcmFuc3BvcnQgPSB7XG4gKiAgIGNhbGw6IGFzeW5jIChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSA9PiB7XG4gKiAgICAgLy8gWW91ciBXZWJTb2NrZXQgaW1wbGVtZW50YXRpb25cbiAqICAgfVxuICogfTtcbiAqXG4gKiBzZXRUcmFuc3BvcnQod3NUcmFuc3BvcnQpO1xuICogYGBgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzZXRUcmFuc3BvcnQodHJhbnNwb3J0OiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCk6IHZvaWQge1xuICAgIGN1c3RvbVRyYW5zcG9ydCA9IHRyYW5zcG9ydDtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIGN1cnJlbnQgdHJhbnNwb3J0ICh1c2VmdWwgZm9yIGV4dGVuZGluZy93cmFwcGluZylcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGdldFRyYW5zcG9ydCgpOiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCB7XG4gICAgcmV0dXJuIGN1c3RvbVRyYW5zcG9ydDtcbn1cblxuLyoqXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxuICpcbiAqIEBwYXJhbSBvYmplY3QgLSBUaGUgb2JqZWN0IHRvIGludm9rZSB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICogQHJldHVybiBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZDogbnVtYmVyLCBhcmdzOiBhbnkgPSBudWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcbiAgICAvLyBVc2UgY3VzdG9tIHRyYW5zcG9ydCBpZiBhdmFpbGFibGVcbiAgICBpZiAoY3VzdG9tVHJhbnNwb3J0KSB7XG4gICAgICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQuY2FsbChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9XG5cbiAgICAvLyBEZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG5cbiAgICBsZXQgYm9keTogeyBvYmplY3Q6IG51bWJlcjsgbWV0aG9kOiBudW1iZXIsIGFyZ3M/OiBhbnkgfSA9IHtcbiAgICAgIG9iamVjdDogb2JqZWN0SUQsXG4gICAgICBtZXRob2RcbiAgICB9XG4gICAgaWYgKGFyZ3MgIT09IG51bGwgJiYgYXJncyAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICBib2R5LmFyZ3MgPSBhcmdzO1xuICAgIH1cblxuICAgIGxldCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+ID0ge1xuICAgICAgICBbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXTogY2xpZW50SWQsXG4gICAgICAgIFtcIkNvbnRlbnQtVHlwZVwiXTogXCJhcHBsaWNhdGlvbi9qc29uXCJcbiAgICB9XG4gICAgaWYgKHdpbmRvd05hbWUpIHtcbiAgICAgICAgaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xuICAgIH1cblxuICAgIGNvbnN0IGJvZHlTdHIgPSBKU09OLnN0cmluZ2lmeShib2R5KTtcbiAgICBsZXQgcmVzcG9uc2U6IFJlc3BvbnNlO1xuICAgIGlmIChib2R5U3RyLmxlbmd0aCA+IENIVU5LX1RIUkVTSE9MRCkge1xuICAgICAgICByZXNwb25zZSA9IGF3YWl0IHNlbmRDaHVua2VkKHVybCwgaGVhZGVycywgYm9keVN0cik7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcmVzcG9uc2UgPSBhd2FpdCBmZXRjaCh1cmwsIHsgbWV0aG9kOiAnUE9TVCcsIGhlYWRlcnMsIGJvZHk6IGJvZHlTdHIgfSk7XG4gICAgfVxuICAgIGlmICghcmVzcG9uc2Uub2spIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XG4gICAgfVxuXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgfVxufVxuXG4vLyBzZW5kQ2h1bmtlZCBzcGxpdHMgYSBsYXJnZSBzZXJpYWxpc2VkIHJlcXVlc3QgYm9keSBpbnRvIENIVU5LX1RIUkVTSE9MRC1zaXplZFxuLy8gYnl0ZSBjaHVua3MgYW5kIHNlbmRzIHRoZW0gc2VyaWFsbHkuICBFbmNvZGluZyB0byBVVEYtOCBieXRlcyBiZWZvcmUgc2xpY2luZ1xuLy8gcHJldmVudHMgY29ycnVwdGlvbiBvZiBub24tQk1QIGNoYXJhY3RlcnMgKHN1cnJvZ2F0ZSBwYWlycykgdGhhdCB3b3VsZCBvY2N1clxuLy8gd2hlbiBzcGxpdHRpbmcgYXQgSmF2YVNjcmlwdCBzdHJpbmcgaW5kaWNlcy4gIFRoZSBHbyB0cmFuc3BvcnQgYXNzZW1ibGVzIHRoZVxuLy8gcmF3IGJ5dGVzIGJlZm9yZSBwcm9jZXNzaW5nLiAgT25seSB0aGUgZmluYWwgY2h1bmsncyByZXNwb25zZSBjYXJyaWVzIHRoZSBSUEMgcmVzdWx0LlxuYXN5bmMgZnVuY3Rpb24gc2VuZENodW5rZWQodXJsOiBVUkwsIGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4sIGJvZHlTdHI6IHN0cmluZyk6IFByb21pc2U8UmVzcG9uc2U+IHtcbiAgICBjb25zdCBjaHVua0lkID0gbmFub2lkKCk7XG4gICAgY29uc3QgYm9keUJ5dGVzID0gbmV3IFRleHRFbmNvZGVyKCkuZW5jb2RlKGJvZHlTdHIpO1xuICAgIGNvbnN0IHRvdGFsQ2h1bmtzID0gTWF0aC5jZWlsKGJvZHlCeXRlcy5sZW5ndGggLyBDSFVOS19USFJFU0hPTEQpO1xuXG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCB0b3RhbENodW5rcyAtIDE7IGkrKykge1xuICAgICAgICBjb25zdCBjaHVuayA9IGJvZHlCeXRlcy5zdWJhcnJheShpICogQ0hVTktfVEhSRVNIT0xELCAoaSArIDEpICogQ0hVTktfVEhSRVNIT0xEKTtcbiAgICAgICAgY29uc3QgcmVzcCA9IGF3YWl0IGZldGNoKHVybCwge1xuICAgICAgICAgICAgbWV0aG9kOiAnUE9TVCcsXG4gICAgICAgICAgICBoZWFkZXJzOiB7XG4gICAgICAgICAgICAgICAgLi4uaGVhZGVycyxcbiAgICAgICAgICAgICAgICAneC13YWlscy1jaHVuay1pZCc6IGNodW5rSWQsXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaW5kZXgnOiBTdHJpbmcoaSksXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstdG90YWwnOiBTdHJpbmcodG90YWxDaHVua3MpLFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGJvZHk6IGNodW5rLFxuICAgICAgICB9KTtcbiAgICAgICAgaWYgKCFyZXNwLm9rKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIGZldGNoKHVybCwge1xuICAgICAgICBtZXRob2Q6ICdQT1NUJyxcbiAgICAgICAgaGVhZGVyczoge1xuICAgICAgICAgICAgLi4uaGVhZGVycyxcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWlkJzogY2h1bmtJZCxcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWluZGV4JzogU3RyaW5nKHRvdGFsQ2h1bmtzIC0gMSksXG4gICAgICAgICAgICAneC13YWlscy1jaHVuay10b3RhbCc6IFN0cmluZyh0b3RhbENodW5rcyksXG4gICAgICAgIH0sXG4gICAgICAgIGJvZHk6IGJvZHlCeXRlcy5zdWJhcnJheSgodG90YWxDaHVua3MgLSAxKSAqIENIVU5LX1RIUkVTSE9MRCksXG4gICAgfSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuLy8gc2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5EaWFsb2cpO1xuXG4vLyBEZWZpbmUgY29uc3RhbnRzIGZyb20gdGhlIGBtZXRob2RzYCBvYmplY3QgaW4gVGl0bGUgQ2FzZVxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XG5jb25zdCBEaWFsb2dXYXJuaW5nID0gMTtcbmNvbnN0IERpYWxvZ0Vycm9yID0gMjtcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcbmNvbnN0IERpYWxvZ09wZW5GaWxlID0gNDtcbmNvbnN0IERpYWxvZ1NhdmVGaWxlID0gNTtcblxuZXhwb3J0IGludGVyZmFjZSBPcGVuRmlsZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBTYXZlRmlsZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBEZWZhdWx0IGZpbGVuYW1lIHRvIHVzZSBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIEZpbGVuYW1lPzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBNZXNzYWdlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIFRoZSB0aXRsZSBvZiB0aGUgZGlhbG9nIHdpbmRvdy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIEJ1dHRvbnM/OiBCdXR0b25bXTtcbiAgICAvKiogVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBCdXR0b24ge1xuICAgIC8qKiBUZXh0IHRoYXQgYXBwZWFycyB3aXRoaW4gdGhlIGJ1dHRvbi4gKi9cbiAgICBMYWJlbD86IHN0cmluZztcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC4gKi9cbiAgICBJc0NhbmNlbD86IGJvb2xlYW47XG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci4gKi9cbiAgICBJc0RlZmF1bHQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEZpbGVGaWx0ZXIge1xuICAgIC8qKiBEaXNwbGF5IG5hbWUgZm9yIHRoZSBmaWx0ZXIsIGl0IGNvdWxkIGJlIFwiVGV4dCBGaWxlc1wiLCBcIkltYWdlc1wiIGV0Yy4gKi9cbiAgICBEaXNwbGF5TmFtZT86IHN0cmluZztcbiAgICAvKiogUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy4gKi9cbiAgICBQYXR0ZXJuPzogc3RyaW5nO1xufVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKlxuICogQHBhcmFtIHR5cGUgLSBEaWFsb2cgdHlwZS5cbiAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyBmb3IgdGhlIGRpYWxvZy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggcmVzdWx0IG9mIGRpYWxvZy5cbiAqL1xuZnVuY3Rpb24gZGlhbG9nKHR5cGU6IG51bWJlciwgb3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMgfCBPcGVuRmlsZURpYWxvZ09wdGlvbnMgfCBTYXZlRmlsZURpYWxvZ09wdGlvbnMgPSB7fSk6IFByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIGNhbGwodHlwZSwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogUHJlc2VudHMgYW4gaW5mbyBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgd2FybmluZyBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhbiBlcnJvciBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRXJyb3Iob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgcXVlc3Rpb24gZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIG9uZSBvciBtb3JlIGZpbGVzIHRvIG9wZW4uXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcywgb3IgYSBibGFuayBzdHJpbmcvZW1wdHkgbGlzdCBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbjogdHJ1ZSB9KTogUHJvbWlzZTxzdHJpbmdbXT47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGZhbHNlIHwgdW5kZWZpbmVkIH0pOiBQcm9taXNlPHN0cmluZz47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT4geyByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKSA/PyBbXTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgYSBmaWxlIHRvIHNhdmUuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUsIG9yIGEgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zOiBTYXZlRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTsgfVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IGV2ZW50TGlzdGVuZXJzLCBMaXN0ZW5lciwgbGlzdGVuZXJPZmYgfSBmcm9tIFwiLi9saXN0ZW5lci5qc1wiO1xuaW1wb3J0IHsgRXZlbnRzIGFzIENyZWF0ZSB9IGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xuaW1wb3J0IHsgVHlwZXMgfSBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xuXG4vLyBTZXR1cFxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLmRpc3BhdGNoV2FpbHNFdmVudCA9IGRpc3BhdGNoV2FpbHNFdmVudDtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRXZlbnRzKTtcbmNvbnN0IEVtaXRNZXRob2QgPSAwO1xuXG5leHBvcnQgKiBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xuXG4vKipcbiAqIEEgdGFibGUgb2YgZGF0YSB0eXBlcyBmb3IgYWxsIGtub3duIGV2ZW50cy5cbiAqIFdpbGwgYmUgbW9ua2V5LXBhdGNoZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIEN1c3RvbUV2ZW50cyB7fVxuXG4vKipcbiAqIEVpdGhlciBhIGtub3duIGV2ZW50IG5hbWUgb3IgYW4gYXJiaXRyYXJ5IHN0cmluZy5cbiAqL1xuZXhwb3J0IHR5cGUgV2FpbHNFdmVudE5hbWU8RSBleHRlbmRzIGtleW9mIEN1c3RvbUV2ZW50cyA9IGtleW9mIEN1c3RvbUV2ZW50cz4gPSBFIHwgKHN0cmluZyAmIHt9KTtcblxuLyoqXG4gKiBVbmlvbiBvZiBhbGwga25vd24gc3lzdGVtIGV2ZW50IG5hbWVzLlxuICovXG50eXBlIFN5c3RlbUV2ZW50TmFtZSA9IHtcbiAgICBbSyBpbiBrZXlvZiAodHlwZW9mIFR5cGVzKV06ICh0eXBlb2YgVHlwZXMpW0tdW2tleW9mICgodHlwZW9mIFR5cGVzKVtLXSldXG59IGV4dGVuZHMgKGluZmVyIE0pID8gTVtrZXlvZiBNXSA6IG5ldmVyO1xuXG4vKipcbiAqIFRoZSBkYXRhIHR5cGUgYXNzb2NpYXRlZCB0byBhIGdpdmVuIGV2ZW50LlxuICovXG5leHBvcnQgdHlwZSBXYWlsc0V2ZW50RGF0YTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4gPVxuICAgIEUgZXh0ZW5kcyBrZXlvZiBDdXN0b21FdmVudHMgPyBDdXN0b21FdmVudHNbRV0gOiAoRSBleHRlbmRzIFN5c3RlbUV2ZW50TmFtZSA/IHZvaWQgOiBhbnkpO1xuXG4vKipcbiAqIFRoZSB0eXBlIG9mIGhhbmRsZXJzIGZvciBhIGdpdmVuIGV2ZW50LlxuICovXG5leHBvcnQgdHlwZSBXYWlsc0V2ZW50Q2FsbGJhY2s8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+ID0gKGV2OiBXYWlsc0V2ZW50PEU+KSA9PiB2b2lkO1xuXG4vKipcbiAqIFJlcHJlc2VudHMgYSBzeXN0ZW0gZXZlbnQgb3IgYSBjdXN0b20gZXZlbnQgZW1pdHRlZCB0aHJvdWdoIHdhaWxzLXByb3ZpZGVkIGZhY2lsaXRpZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiB7XG4gICAgLyoqXG4gICAgICogVGhlIG5hbWUgb2YgdGhlIGV2ZW50LlxuICAgICAqL1xuICAgIG5hbWU6IEU7XG5cbiAgICAvKipcbiAgICAgKiBPcHRpb25hbCBkYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZW1pdHRlZCBldmVudC5cbiAgICAgKi9cbiAgICBkYXRhOiBXYWlsc0V2ZW50RGF0YTxFPjtcblxuICAgIC8qKlxuICAgICAqIE5hbWUgb2YgdGhlIG9yaWdpbmF0aW5nIHdpbmRvdy4gT21pdHRlZCBmb3IgYXBwbGljYXRpb24gZXZlbnRzLlxuICAgICAqIFdpbGwgYmUgb3ZlcnJpZGRlbiBpZiBzZXQgbWFudWFsbHkuXG4gICAgICovXG4gICAgc2VuZGVyPzogc3RyaW5nO1xuXG4gICAgY29uc3RydWN0b3IobmFtZTogRSwgZGF0YTogV2FpbHNFdmVudERhdGE8RT4pO1xuICAgIGNvbnN0cnVjdG9yKG5hbWU6IFdhaWxzRXZlbnREYXRhPEU+IGV4dGVuZHMgbnVsbCB8IHZvaWQgPyBFIDogbmV2ZXIpXG4gICAgY29uc3RydWN0b3IobmFtZTogRSwgZGF0YT86IGFueSkge1xuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhID8/IG51bGw7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQ6IGFueSkge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxldCB3YWlsc0V2ZW50ID0gbmV3IFdhaWxzRXZlbnQoXG4gICAgICAgIGV2ZW50Lm5hbWUsXG4gICAgICAgIChldmVudC5uYW1lIGluIENyZWF0ZSkgPyBDcmVhdGVbZXZlbnQubmFtZV0oZXZlbnQuZGF0YSkgOiBldmVudC5kYXRhXG4gICAgKTtcbiAgICBpZiAoJ3NlbmRlcicgaW4gZXZlbnQpIHtcbiAgICAgICAgd2FpbHNFdmVudC5zZW5kZXIgPSBldmVudC5zZW5kZXI7XG4gICAgfVxuXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsaXN0ZW5lciA9PiAhbGlzdGVuZXIuZGlzcGF0Y2god2FpbHNFdmVudCkpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcGFyYW0gbWF4Q2FsbGJhY2tzIC0gVGhlIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRoZSBjYWxsYmFjayBjYW4gYmUgY2FsbGVkIGZvciB0aGUgZXZlbnQuIE9uY2UgdGhlIG1heGltdW0gbnVtYmVyIGlzIHJlYWNoZWQsIHRoZSBjYWxsYmFjayB3aWxsIG5vIGxvbmdlciBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+LCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbjxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+KTogKCkgPT4gdm9pZCB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihldmVudE5hbWU6IEUsIGNhbGxiYWNrOiBXYWlsc0V2ZW50Q2FsbGJhY2s8RT4pOiAoKSA9PiB2b2lkIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuLyoqXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lcyAtIFRoZSBuYW1lIG9mIHRoZSBldmVudHMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoLi4uZXZlbnROYW1lczogW1dhaWxzRXZlbnROYW1lLCAuLi5XYWlsc0V2ZW50TmFtZVtdXSk6IHZvaWQge1xuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xufVxuXG4vKipcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpOiB2b2lkIHtcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xufVxuXG4vKipcbiAqIEVtaXRzIGFuIGV2ZW50LlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgZnVsZmlsbGVkIG9uY2UgdGhlIGV2ZW50IGhhcyBiZWVuIGVtaXR0ZWQuICBSZXNvbHZlcyB0byB0cnVlIGlmIHRoZSBldmVudCB3YXMgY2FuY2VsbGVkLlxuICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdFxuICogQHBhcmFtIGRhdGEgLSBUaGUgZGF0YSB0aGF0IHdpbGwgYmUgc2VudCB3aXRoIHRoZSBldmVudFxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogRSwgZGF0YTogV2FpbHNFdmVudERhdGE8RT4pOiBQcm9taXNlPGJvb2xlYW4+XG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogV2FpbHNFdmVudERhdGE8RT4gZXh0ZW5kcyBudWxsIHwgdm9pZCA/IEUgOiBuZXZlcik6IFByb21pc2U8Ym9vbGVhbj5cbmV4cG9ydCBmdW5jdGlvbiBFbWl0PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiwgZGF0YT86IGFueSk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgIHJldHVybiBjYWxsKEVtaXRNZXRob2QsICBuZXcgV2FpbHNFdmVudChuYW1lLCBkYXRhKSlcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBUaGUgZm9sbG93aW5nIHV0aWxpdGllcyBoYXZlIGJlZW4gZmFjdG9yZWQgb3V0IG9mIC4vZXZlbnRzLnRzXG4vLyBmb3IgdGVzdGluZyBwdXJwb3Nlcy5cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcDxzdHJpbmcsIExpc3RlbmVyW10+KCk7XG5cbmV4cG9ydCBjbGFzcyBMaXN0ZW5lciB7XG4gICAgZXZlbnROYW1lOiBzdHJpbmc7XG4gICAgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQ7XG4gICAgbWF4Q2FsbGJhY2tzOiBudW1iZXI7XG5cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQsIG1heENhbGxiYWNrczogbnVtYmVyKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICB0aGlzLmNhbGxiYWNrID0gY2FsbGJhY2s7XG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgIH1cblxuICAgIGRpc3BhdGNoKGRhdGE6IGFueSk6IGJvb2xlYW4ge1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgdGhpcy5jYWxsYmFjayhkYXRhKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGVycik7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSByZXR1cm4gZmFsc2U7XG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcjogTGlzdGVuZXIpOiB2b2lkIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGxpc3RlbmVyLmV2ZW50TmFtZSk7XG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGxpc3RlbmVyLmV2ZW50TmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGxpc3RlbmVyLmV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogQW55IGlzIGEgZHVtbXkgY3JlYXRpb24gZnVuY3Rpb24gZm9yIHNpbXBsZSBvciB1bmtub3duIHR5cGVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQW55PFQgPSBhbnk+KHNvdXJjZTogYW55KTogVCB7XG4gICAgcmV0dXJuIHNvdXJjZTtcbn1cblxuLyoqXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXG4gKiBudWxsIHN0cmluZ3Mgd2l0aCBlbXB0eSBzdHJpbmdzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnl0ZVNsaWNlKHNvdXJjZTogYW55KTogc3RyaW5nIHtcbiAgICByZXR1cm4gKChzb3VyY2UgPT0gbnVsbCkgPyBcIlwiIDogc291cmNlKTtcbn1cblxuLyoqXG4gKiBBcnJheSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XG4gKiB3aG9zZSBlbGVtZW50cyBhcmUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQXJyYXk8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiBUW10ge1xuICAgIGlmIChlbGVtZW50ID09PSBBbnkpIHtcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBbXSA6IHNvdXJjZSk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xuICAgICAgICAgICAgcmV0dXJuIFtdO1xuICAgICAgICB9XG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgICAgICBzb3VyY2VbaV0gPSBlbGVtZW50KHNvdXJjZVtpXSk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcbiAgICB9O1xufVxuXG4vKipcbiAqIE1hcCB0YWtlcyBjcmVhdGlvbiBmdW5jdGlvbnMgZm9yIHR3byBhcmJpdHJhcnkgdHlwZXNcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcbiAqIHdob3NlIGtleXMgYW5kIHZhbHVlcyBhcmUgb2YgdGhvc2UgdHlwZXMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBNYXA8SyBleHRlbmRzIFByb3BlcnR5S2V5ID0gYW55LCBWID0gYW55PihrZXk6IChzb3VyY2U6IGFueSkgPT4gSywgdmFsdWU6IChzb3VyY2U6IGFueSkgPT4gVik6IChzb3VyY2U6IGFueSkgPT4gUmVjb3JkPEssIFY+IHtcbiAgICBpZiAodmFsdWUgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IHt9IDogc291cmNlKTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm4ge307XG4gICAgICAgIH1cbiAgICAgICAgZm9yIChjb25zdCBrZXkgaW4gc291cmNlKSB7XG4gICAgICAgICAgICBzb3VyY2Vba2V5XSA9IHZhbHVlKHNvdXJjZVtrZXldKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogTnVsbGFibGUgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgbnVsbGFibGUgdmFsdWUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTnVsbGFibGU8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiAoVCB8IG51bGwpIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcbn1cblxuLyoqXG4gKiBTdHJ1Y3QgdGFrZXMgYW4gb2JqZWN0IG1hcHBpbmcgZmllbGQgbmFtZXMgdG8gY3JlYXRpb24gZnVuY3Rpb25zXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBzdHJ1Y3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQ6IFJlY29yZDxzdHJpbmcsIChzb3VyY2U6IGFueSkgPT4gYW55Pik6XG4gICAgPFUgZXh0ZW5kcyBSZWNvcmQ8c3RyaW5nLCBhbnk+ID0gYW55Pihzb3VyY2U6IGFueSkgPT4gVVxue1xuICAgIGxldCBhbGxBbnkgPSB0cnVlO1xuICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICBpZiAoY3JlYXRlRmllbGRbbmFtZV0gIT09IEFueSkge1xuICAgICAgICAgICAgYWxsQW55ID0gZmFsc2U7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoYWxsQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XG4gICAgICAgICAgICBpZiAobmFtZSBpbiBzb3VyY2UpIHtcbiAgICAgICAgICAgICAgICBzb3VyY2VbbmFtZV0gPSBjcmVhdGVGaWVsZFtuYW1lXShzb3VyY2VbbmFtZV0pO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBNYXBzIGtub3duIGV2ZW50IG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdGhlaXIgZGF0YSB0eXBlcy5cbiAqIFdpbGwgYmUgbW9ua2V5LXBhdGNoZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxuICovXG5leHBvcnQgY29uc3QgRXZlbnRzOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4gPSB7fTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gQ3luaHlyY2h3eWQgeSBmZmVpbCBob24geW4gYXd0b21hdGlnLiBQRUlESVdDSCBcdTAwQzIgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRXYWtlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFdha2VcIixcblx0XHRBcHBsaWNhdGlvblNjcmVlbnNEaWRTbGVlcDogXCJtYWM6QXBwbGljYXRpb25TY3JlZW5zRGlkU2xlZXBcIixcblx0XHRBcHBsaWNhdGlvblNjcmVlbnNEaWRXYWtlOiBcIm1hYzpBcHBsaWNhdGlvblNjcmVlbnNEaWRXYWtlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFNsZWVwOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxTbGVlcFwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHR9KSxcblx0TGludXg6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcblx0XHRTeXN0ZW1EaWRXYWtlOiBcImxpbnV4OlN5c3RlbURpZFdha2VcIixcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0U3lzdGVtV2lsbFNsZWVwOiBcImxpbnV4OlN5c3RlbVdpbGxTbGVlcFwiLFxuXHRcdFdpbmRvd0RlbGV0ZUV2ZW50OiBcImxpbnV4OldpbmRvd0RlbGV0ZUV2ZW50XCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJsaW51eDpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImxpbnV4OldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0ZvY3VzSW46IFwibGludXg6V2luZG93Rm9jdXNJblwiLFxuXHRcdFdpbmRvd0ZvY3VzT3V0OiBcImxpbnV4OldpbmRvd0ZvY3VzT3V0XCIsXG5cdFx0V2luZG93TG9hZFN0YXJ0ZWQ6IFwibGludXg6V2luZG93TG9hZFN0YXJ0ZWRcIixcblx0XHRXaW5kb3dMb2FkUmVkaXJlY3RlZDogXCJsaW51eDpXaW5kb3dMb2FkUmVkaXJlY3RlZFwiLFxuXHRcdFdpbmRvd0xvYWRDb21taXR0ZWQ6IFwibGludXg6V2luZG93TG9hZENvbW1pdHRlZFwiLFxuXHRcdFdpbmRvd0xvYWRGaW5pc2hlZDogXCJsaW51eDpXaW5kb3dMb2FkRmluaXNoZWRcIixcblx0fSksXG5cdGlPUzogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmU6IFwiaW9zOkFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRFbnRlckJhY2tncm91bmQ6IFwiaW9zOkFwcGxpY2F0aW9uRGlkRW50ZXJCYWNrZ3JvdW5kXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwiaW9zOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZWNlaXZlTWVtb3J5V2FybmluZzogXCJpb3M6QXBwbGljYXRpb25EaWRSZWNlaXZlTWVtb3J5V2FybmluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEVudGVyRm9yZWdyb3VuZDogXCJpb3M6QXBwbGljYXRpb25XaWxsRW50ZXJGb3JlZ3JvdW5kXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlOiBcImlvczpBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGU6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZVwiLFxuXHRcdFdpbmRvd0RpZExvYWQ6IFwiaW9zOldpbmRvd0RpZExvYWRcIixcblx0XHRXaW5kb3dXaWxsQXBwZWFyOiBcImlvczpXaW5kb3dXaWxsQXBwZWFyXCIsXG5cdFx0V2luZG93RGlkQXBwZWFyOiBcImlvczpXaW5kb3dEaWRBcHBlYXJcIixcblx0XHRXaW5kb3dXaWxsRGlzYXBwZWFyOiBcImlvczpXaW5kb3dXaWxsRGlzYXBwZWFyXCIsXG5cdFx0V2luZG93RGlkRGlzYXBwZWFyOiBcImlvczpXaW5kb3dEaWREaXNhcHBlYXJcIixcblx0XHRXaW5kb3dTYWZlQXJlYUluc2V0c0NoYW5nZWQ6IFwiaW9zOldpbmRvd1NhZmVBcmVhSW5zZXRzQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd09yaWVudGF0aW9uQ2hhbmdlZDogXCJpb3M6V2luZG93T3JpZW50YXRpb25DaGFuZ2VkXCIsXG5cdFx0V2luZG93VG91Y2hCZWdhbjogXCJpb3M6V2luZG93VG91Y2hCZWdhblwiLFxuXHRcdFdpbmRvd1RvdWNoTW92ZWQ6IFwiaW9zOldpbmRvd1RvdWNoTW92ZWRcIixcblx0XHRXaW5kb3dUb3VjaEVuZGVkOiBcImlvczpXaW5kb3dUb3VjaEVuZGVkXCIsXG5cdFx0V2luZG93VG91Y2hDYW5jZWxsZWQ6IFwiaW9zOldpbmRvd1RvdWNoQ2FuY2VsbGVkXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0TmF2aWdhdGlvbjogXCJpb3M6V2ViVmlld0RpZFN0YXJ0TmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRGYWlsTmF2aWdhdGlvbjogXCJpb3M6V2ViVmlld0RpZEZhaWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RlY2lkZVBvbGljeUZvck5hdmlnYXRpb25BY3Rpb246IFwiaW9zOldlYlZpZXdEZWNpZGVQb2xpY3lGb3JOYXZpZ2F0aW9uQWN0aW9uXCIsXG5cdH0pLFxuXHRDb21tb246IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiY29tbW9uOkFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGVcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwiY29tbW9uOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdEFwcGxpY2F0aW9uTGF1bmNoZWRXaXRoVXJsOiBcImNvbW1vbjpBcHBsaWNhdGlvbkxhdW5jaGVkV2l0aFVybFwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXG5cdFx0U3lzdGVtRGlkV2FrZTogXCJjb21tb246U3lzdGVtRGlkV2FrZVwiLFxuXHRcdFN5c3RlbVdpbGxTbGVlcDogXCJjb21tb246U3lzdGVtV2lsbFNsZWVwXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJjb21tb246V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwiY29tbW9uOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwiY29tbW9uOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwiY29tbW9uOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dGaWxlc0Ryb3BwZWQ6IFwiY29tbW9uOldpbmRvd0ZpbGVzRHJvcHBlZFwiLFxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dIaWRlOiBcImNvbW1vbjpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93TG9zdEZvY3VzOiBcImNvbW1vbjpXaW5kb3dMb3N0Rm9jdXNcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJjb21tb246V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJjb21tb246V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dUb2dnbGVGcmFtZWxlc3M6IFwiY29tbW9uOldpbmRvd1RvZ2dsZUZyYW1lbGVzc1wiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwiY29tbW9uOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dSdW50aW1lUmVhZHk6IFwiY29tbW9uOldpbmRvd1J1bnRpbWVSZWFkeVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Wm9vbTogXCJjb21tb246V2luZG93Wm9vbVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJjb21tb246V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJjb21tb246V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJjb21tb246V2luZG93Wm9vbVJlc2V0XCIsXG5cdH0pLFxufSk7XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogTG9ncyBhIG1lc3NhZ2UgdG8gdGhlIGNvbnNvbGUgd2l0aCBjdXN0b20gZm9ybWF0dGluZy5cbiAqXG4gKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBtZXNzYWdlIHRvIGJlIGxvZ2dlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGRlYnVnTG9nKG1lc3NhZ2U6IGFueSkge1xuICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZVxuICAgIGNvbnNvbGUubG9nKFxuICAgICAgICAnJWMgd2FpbHMzICVjICcgKyBtZXNzYWdlICsgJyAnLFxuICAgICAgICAnYmFja2dyb3VuZDogI2FhMDAwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDNweCAwcHggMHB4IDNweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbScsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjMDA5OTAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogMHB4IDNweCAzcHggMHB4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJ1xuICAgICk7XG59XG5cbi8qKlxuICogQ2hlY2tzIHdoZXRoZXIgdGhlIHdlYnZpZXcgc3VwcG9ydHMgdGhlIHtAbGluayBNb3VzZUV2ZW50I2J1dHRvbnN9IHByb3BlcnR5LlxuICogTG9va2luZyBhdCB5b3UgbWFjT1MgSGlnaCBTaWVycmEhXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjYW5UcmFja0J1dHRvbnMoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIChuZXcgTW91c2VFdmVudCgnbW91c2Vkb3duJykpLmJ1dHRvbnMgPT09IDA7XG59XG5cbi8qKlxuICogQ2hlY2tzIHdoZXRoZXIgdGhlIGJyb3dzZXIgc3VwcG9ydHMgcmVtb3ZpbmcgbGlzdGVuZXJzIGJ5IHRyaWdnZXJpbmcgYW4gQWJvcnRTaWduYWxcbiAqIChzZWUgaHR0cHM6Ly9kZXZlbG9wZXIubW96aWxsYS5vcmcvZW4tVVMvZG9jcy9XZWIvQVBJL0V2ZW50VGFyZ2V0L2FkZEV2ZW50TGlzdGVuZXIjc2lnbmFsKS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGNhbkFib3J0TGlzdGVuZXJzKCkge1xuICAgIGlmICghRXZlbnRUYXJnZXQgfHwgIUFib3J0U2lnbmFsIHx8ICFBYm9ydENvbnRyb2xsZXIpXG4gICAgICAgIHJldHVybiBmYWxzZTtcblxuICAgIGxldCByZXN1bHQgPSB0cnVlO1xuXG4gICAgY29uc3QgdGFyZ2V0ID0gbmV3IEV2ZW50VGFyZ2V0KCk7XG4gICAgY29uc3QgY29udHJvbGxlciA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB0YXJnZXQuYWRkRXZlbnRMaXN0ZW5lcigndGVzdCcsICgpID0+IHsgcmVzdWx0ID0gZmFsc2U7IH0sIHsgc2lnbmFsOiBjb250cm9sbGVyLnNpZ25hbCB9KTtcbiAgICBjb250cm9sbGVyLmFib3J0KCk7XG4gICAgdGFyZ2V0LmRpc3BhdGNoRXZlbnQobmV3IEN1c3RvbUV2ZW50KCd0ZXN0JykpO1xuXG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBSZXNvbHZlcyB0aGUgY2xvc2VzdCBIVE1MRWxlbWVudCBhbmNlc3RvciBvZiBhbiBldmVudCdzIHRhcmdldC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGV2ZW50VGFyZ2V0KGV2ZW50OiBFdmVudCk6IEhUTUxFbGVtZW50IHtcbiAgICBpZiAoZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpIHtcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldDtcbiAgICB9IGVsc2UgaWYgKCEoZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpICYmIGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIE5vZGUpIHtcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldC5wYXJlbnRFbGVtZW50ID8/IGRvY3VtZW50LmJvZHk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcmV0dXJuIGRvY3VtZW50LmJvZHk7XG4gICAgfVxufVxuXG4vKioqXG4gVGhpcyB0ZWNobmlxdWUgZm9yIHByb3BlciBsb2FkIGRldGVjdGlvbiBpcyB0YWtlbiBmcm9tIEhUTVg6XG5cbiBCU0QgMi1DbGF1c2UgTGljZW5zZVxuXG4gQ29weXJpZ2h0IChjKSAyMDIwLCBCaWcgU2t5IFNvZnR3YXJlXG4gQWxsIHJpZ2h0cyByZXNlcnZlZC5cblxuIFJlZGlzdHJpYnV0aW9uIGFuZCB1c2UgaW4gc291cmNlIGFuZCBiaW5hcnkgZm9ybXMsIHdpdGggb3Igd2l0aG91dFxuIG1vZGlmaWNhdGlvbiwgYXJlIHBlcm1pdHRlZCBwcm92aWRlZCB0aGF0IHRoZSBmb2xsb3dpbmcgY29uZGl0aW9ucyBhcmUgbWV0OlxuXG4gMS4gUmVkaXN0cmlidXRpb25zIG9mIHNvdXJjZSBjb2RlIG11c3QgcmV0YWluIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLCB0aGlzXG4gbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIuXG5cbiAyLiBSZWRpc3RyaWJ1dGlvbnMgaW4gYmluYXJ5IGZvcm0gbXVzdCByZXByb2R1Y2UgdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsXG4gdGhpcyBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lciBpbiB0aGUgZG9jdW1lbnRhdGlvblxuIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLlxuXG4gVEhJUyBTT0ZUV0FSRSBJUyBQUk9WSURFRCBCWSBUSEUgQ09QWVJJR0hUIEhPTERFUlMgQU5EIENPTlRSSUJVVE9SUyBcIkFTIElTXCJcbiBBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFXG4gSU1QTElFRCBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSBBTkQgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQVJFXG4gRElTQ0xBSU1FRC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIENPUFlSSUdIVCBIT0xERVIgT1IgQ09OVFJJQlVUT1JTIEJFIExJQUJMRVxuIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMXG4gREFNQUdFUyAoSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFBST0NVUkVNRU5UIE9GIFNVQlNUSVRVVEUgR09PRFMgT1JcbiBTRVJWSUNFUzsgTE9TUyBPRiBVU0UsIERBVEEsIE9SIFBST0ZJVFM7IE9SIEJVU0lORVNTIElOVEVSUlVQVElPTikgSE9XRVZFUlxuIENBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksXG4gT1IgVE9SVCAoSU5DTFVESU5HIE5FR0xJR0VOQ0UgT1IgT1RIRVJXSVNFKSBBUklTSU5HIElOIEFOWSBXQVkgT1VUIE9GIFRIRSBVU0VcbiBPRiBUSElTIFNPRlRXQVJFLCBFVkVOIElGIEFEVklTRUQgT0YgVEhFIFBPU1NJQklMSVRZIE9GIFNVQ0ggREFNQUdFLlxuXG4gKioqL1xuXG5sZXQgaXNSZWFkeSA9IGZhbHNlO1xuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsICgpID0+IHsgaXNSZWFkeSA9IHRydWUgfSk7XG5cbmV4cG9ydCBmdW5jdGlvbiB3aGVuUmVhZHkoY2FsbGJhY2s6ICgpID0+IHZvaWQpIHtcbiAgICBpZiAoaXNSZWFkeSB8fCBkb2N1bWVudC5yZWFkeVN0YXRlID09PSAnY29tcGxldGUnKSB7XG4gICAgICAgIGNhbGxiYWNrKCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGNhbGxiYWNrKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB0eXBlIHsgU2NyZWVuIH0gZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xuXG4vLyBEcm9wIHRhcmdldCBjb25zdGFudHNcbmNvbnN0IERST1BfVEFSR0VUX0FUVFJJQlVURSA9ICdkYXRhLWZpbGUtZHJvcC10YXJnZXQnO1xuY29uc3QgRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTID0gJ2ZpbGUtZHJvcC10YXJnZXQtYWN0aXZlJztcbmxldCBjdXJyZW50RHJvcFRhcmdldDogRWxlbWVudCB8IG51bGwgPSBudWxsO1xuXG5jb25zdCBQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAwO1xuY29uc3QgQ2VudGVyTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMTtcbmNvbnN0IENsb3NlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDI7XG5jb25zdCBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgPSAzO1xuY29uc3QgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgID0gNDtcbmNvbnN0IEZvY3VzTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDU7XG5jb25zdCBGb3JjZVJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgPSA2O1xuY29uc3QgRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgID0gNztcbmNvbnN0IEdldFNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDg7XG5jb25zdCBHZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA5O1xuY29uc3QgSGVpZ2h0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMTA7XG5jb25zdCBIaWRlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxMTtcbmNvbnN0IElzRm9jdXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDEyO1xuY29uc3QgSXNGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gMTM7XG5jb25zdCBJc01heGltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNDtcbmNvbnN0IElzTWluaW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE1O1xuY29uc3QgTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTY7XG5jb25zdCBNaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNztcbmNvbnN0IE5hbWVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDE4O1xuY29uc3QgT3BlbkRldlRvb2xzTWV0aG9kICAgICAgICAgICAgICAgID0gMTk7XG5jb25zdCBSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgPSAyMDtcbmNvbnN0IFJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDIxO1xuY29uc3QgUmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMjI7XG5jb25zdCBSZXN0b3JlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAyMztcbmNvbnN0IFNldFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICA9IDI0O1xuY29uc3QgU2V0QWx3YXlzT25Ub3BNZXRob2QgICAgICAgICAgICAgID0gMjU7XG5jb25zdCBTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kICAgICAgICAgPSAyNjtcbmNvbnN0IFNldEZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICAgICA9IDI3O1xuY29uc3QgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWRNZXRob2QgID0gMjg7XG5jb25zdCBTZXRNYXhTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAyOTtcbmNvbnN0IFNldE1pblNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDMwO1xuY29uc3QgU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgID0gMzE7XG5jb25zdCBTZXRSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgPSAzMjtcbmNvbnN0IFNldFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDMzO1xuY29uc3QgU2V0VGl0bGVNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMzQ7XG5jb25zdCBTZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzNTtcbmNvbnN0IFNob3dNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM2O1xuY29uc3QgU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzc7XG5jb25zdCBUb2dnbGVGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgPSAzODtcbmNvbnN0IFRvZ2dsZU1heGltaXNlTWV0aG9kICAgICAgICAgICAgICA9IDM5O1xuY29uc3QgVG9nZ2xlRnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgID0gNDA7IFxuY29uc3QgVW5GdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gNDE7XG5jb25zdCBVbk1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MjtcbmNvbnN0IFVuTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQzO1xuY29uc3QgV2lkdGhNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XG5jb25zdCBab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSA0NTtcbmNvbnN0IFpvb21Jbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDQ2O1xuY29uc3QgWm9vbU91dE1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gNDc7XG5jb25zdCBab29tUmVzZXRNZXRob2QgICAgICAgICAgICAgICAgICAgPSA0ODtcbmNvbnN0IFNuYXBBc3Npc3RNZXRob2QgICAgICAgICAgICAgICAgICA9IDQ5O1xuY29uc3QgRmlsZXNEcm9wcGVkICAgICAgICAgICAgICAgICAgICAgID0gNTA7XG5jb25zdCBQcmludE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1MTtcbmNvbnN0IFNldFNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDUyO1xuXG4vKipcbiAqIEZpbmRzIHRoZSBuZWFyZXN0IGRyb3AgdGFyZ2V0IGVsZW1lbnQgYnkgd2Fsa2luZyB1cCB0aGUgRE9NIHRyZWUuXG4gKi9cbmZ1bmN0aW9uIGdldERyb3BUYXJnZXRFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQgfCBudWxsKTogRWxlbWVudCB8IG51bGwge1xuICAgIGlmICghZWxlbWVudCkge1xuICAgICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG4gICAgcmV0dXJuIGVsZW1lbnQuY2xvc2VzdChgWyR7RFJPUF9UQVJHRVRfQVRUUklCVVRFfV1gKTtcbn1cblxuLyoqXG4gKiBDaGVjayBpZiB3ZSBjYW4gdXNlIFdlYlZpZXcyJ3MgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgKFdpbmRvd3MpXG4gKiBBbHNvIGNoZWNrcyB0aGF0IEVuYWJsZUZpbGVEcm9wIGlzIHRydWUgZm9yIHRoaXMgd2luZG93LlxuICovXG5mdW5jdGlvbiBjYW5SZXNvbHZlRmlsZVBhdGhzKCk6IGJvb2xlYW4ge1xuICAgIC8vIE11c3QgaGF2ZSBXZWJWaWV3MidzIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIEFQSSAoV2luZG93cyBvbmx5KVxuICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyA9PSBudWxsKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG4gICAgLy8gTXVzdCBoYXZlIEVuYWJsZUZpbGVEcm9wIHNldCB0byB0cnVlIGZvciB0aGlzIHdpbmRvd1xuICAgIC8vIFRoaXMgZmxhZyBpcyBzZXQgYnkgdGhlIEdvIGJhY2tlbmQgZHVyaW5nIHJ1bnRpbWUgaW5pdGlhbGl6YXRpb25cbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSB0cnVlO1xufVxuXG4vKipcbiAqIFNlbmQgZmlsZSBkcm9wIHRvIGJhY2tlbmQgdmlhIFdlYlZpZXcyIChXaW5kb3dzIG9ubHkpXG4gKi9cbmZ1bmN0aW9uIHJlc29sdmVGaWxlUGF0aHMoeDogbnVtYmVyLCB5OiBudW1iZXIsIGZpbGVzOiBGaWxlW10pOiB2b2lkIHtcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMpIHtcbiAgICAgICAgKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKGBmaWxlOmRyb3A6JHt4fToke3l9YCwgZmlsZXMpO1xuICAgIH1cbn1cblxuLy8gTmF0aXZlIGRyYWcgc3RhdGUgKExpbnV4L21hY09TIGludGVyY2VwdCBET00gZHJhZyBldmVudHMpXG5sZXQgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xuXG4vKipcbiAqIENsZWFucyB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhbmQgaG92ZXIgZWZmZWN0cy5cbiAqIENhbGxlZCBvbiBkcm9wIG9yIHdoZW4gZHJhZyBsZWF2ZXMgdGhlIHdpbmRvdy5cbiAqL1xuZnVuY3Rpb24gY2xlYW51cE5hdGl2ZURyYWcoKTogdm9pZCB7XG4gICAgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCkge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbGVkIGZyb20gR28gd2hlbiBhIGZpbGUgZHJhZyBlbnRlcnMgdGhlIHdpbmRvdyBvbiBMaW51eC9tYWNPUy5cbiAqL1xuZnVuY3Rpb24gaGFuZGxlRHJhZ0VudGVyKCk6IHZvaWQge1xuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBhY3RpdmF0ZSBkcmFnIHN0YXRlXG4gICAgfVxuICAgIG5hdGl2ZURyYWdBY3RpdmUgPSB0cnVlO1xufVxuXG4vKipcbiAqIENhbGxlZCBmcm9tIEdvIHdoZW4gYSBmaWxlIGRyYWcgbGVhdmVzIHRoZSB3aW5kb3cgb24gTGludXgvbWFjT1MuXG4gKi9cbmZ1bmN0aW9uIGhhbmRsZURyYWdMZWF2ZSgpOiB2b2lkIHtcbiAgICBjbGVhbnVwTmF0aXZlRHJhZygpO1xufVxuXG4vKipcbiAqIENhbGxlZCBmcm9tIEdvIGR1cmluZyBmaWxlIGRyYWcgdG8gdXBkYXRlIGhvdmVyIHN0YXRlIG9uIExpbnV4L21hY09TLlxuICogQHBhcmFtIHggLSBYIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xuICogQHBhcmFtIHkgLSBZIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xuICovXG5mdW5jdGlvbiBoYW5kbGVEcmFnT3Zlcih4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgIGlmICghbmF0aXZlRHJhZ0FjdGl2ZSkgcmV0dXJuO1xuICAgIFxuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcbiAgICB9XG4gICAgXG4gICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xuICAgIFxuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgfVxuICAgIFxuICAgIGlmIChkcm9wVGFyZ2V0KSB7XG4gICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgIH1cbn1cblxuXG5cbi8vIEV4cG9ydCB0aGUgaGFuZGxlcnMgZm9yIHVzZSBieSBHbyB2aWEgaW5kZXgudHNcbmV4cG9ydCB7IGhhbmRsZURyYWdFbnRlciwgaGFuZGxlRHJhZ0xlYXZlLCBoYW5kbGVEcmFnT3ZlciB9O1xuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHBvc2l0aW9uIG9mIGEgd2luZG93LlxuICovXG5pbnRlcmZhY2UgUG9zaXRpb24ge1xuICAgIC8qKiBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xuICAgIHg6IG51bWJlcjtcbiAgICAvKiogVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXG4gICAgeTogbnVtYmVyO1xufVxuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXG4gKi9cbmludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuICovXG4gICAgd2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LiAqL1xuICAgIGhlaWdodDogbnVtYmVyO1xufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY2FsbGVyU3ltID0gU3ltYm9sKFwiY2FsbGVyXCIpO1xuXG5jbGFzcyBXaW5kb3cge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIHByaXZhdGUgW2NhbGxlclN5bV06IChtZXNzYWdlOiBudW1iZXIsIGFyZ3M/OiBhbnkpID0+IFByb21pc2U8YW55PjtcblxuICAgIC8qKlxuICAgICAqIEluaXRpYWxpc2VzIGEgd2luZG93IG9iamVjdCB3aXRoIHRoZSBzcGVjaWZpZWQgbmFtZS5cbiAgICAgKlxuICAgICAqIEBwcml2YXRlXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgdGFyZ2V0IHdpbmRvdy5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcgPSAnJykge1xuICAgICAgICB0aGlzW2NhbGxlclN5bV0gPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcblxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xuICAgICAgICAgICAgaWYgKFxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXG4gICAgICAgICAgICApIHtcbiAgICAgICAgICAgICAgICAodGhpcyBhcyBhbnkpW21ldGhvZF0gPSAodGhpcyBhcyBhbnkpW21ldGhvZF0uYmluZCh0aGlzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gZ2V0LlxuICAgICAqIEByZXR1cm5zIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXG4gICAgICovXG4gICAgR2V0KG5hbWU6IHN0cmluZyk6IFdpbmRvdyB7XG4gICAgICAgIHJldHVybiBuZXcgV2luZG93KG5hbWUpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2VudGVycyB0aGUgd2luZG93IG9uIHRoZSBzY3JlZW4uXG4gICAgICovXG4gICAgQ2VudGVyKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENlbnRlck1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgQ2xvc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2xvc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKi9cbiAgICBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEVuYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqL1xuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBGb2N1cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb2N1c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKi9cbiAgICBGb3JjZVJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb3JjZVJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU3dpdGNoZXMgdGhlIHdpbmRvdyB0byBmdWxsc2NyZWVuIG1vZGUuXG4gICAgICovXG4gICAgRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBzY3JlZW4gdGhhdCB0aGUgd2luZG93IGlzIG9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvbi5cbiAgICAgKi9cbiAgICBHZXRTY3JlZW4oKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRTY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgem9vbSBsZXZlbC5cbiAgICAgKi9cbiAgICBHZXRab29tKCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0Wm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBIZWlnaHQoKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIZWlnaHRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEhpZGVzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIaWRlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkLlxuICAgICAqL1xuICAgIElzRm9jdXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0ZvY3VzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW4uXG4gICAgICovXG4gICAgSXNGdWxsc2NyZWVuKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtYXhpbWlzZWQuXG4gICAgICovXG4gICAgSXNNYXhpbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNYXhpbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkLlxuICAgICAqL1xuICAgIElzTWluaW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWluaW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogTWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTmFtZSgpOiBQcm9taXNlPHN0cmluZz4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE5hbWVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxuICAgICAqL1xuICAgIE9wZW5EZXZUb29scygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShPcGVuRGV2VG9vbHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUmVsYXRpdmVQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXG4gICAgICovXG4gICAgUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGUuXG4gICAgICovXG4gICAgUmVzaXphYmxlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc2l6YWJsZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVzdG9yZXMgdGhlIHdpbmRvdyB0byBpdHMgcHJldmlvdXMgc3RhdGUgaWYgaXQgd2FzIHByZXZpb3VzbHkgbWluaW1pc2VkLCBtYXhpbWlzZWQgb3IgZnVsbHNjcmVlbi5cbiAgICAgKi9cbiAgICBSZXN0b3JlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc3RvcmVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0UG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxuICAgICAqXG4gICAgICogQHBhcmFtIGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcC5cbiAgICAgKi9cbiAgICBTZXRBbHdheXNPblRvcChhbHdheXNPblRvcDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gciAtIFRoZSBkZXNpcmVkIHJlZCBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYSAtIFRoZSBkZXNpcmVkIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICovXG4gICAgU2V0QmFja2dyb3VuZENvbG91cihyOiBudW1iZXIsIGc6IG51bWJlciwgYjogbnVtYmVyLCBhOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3MuXG4gICAgICovXG4gICAgU2V0RnJhbWVsZXNzKGZyYW1lbGVzczogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWQuXG4gICAgICovXG4gICAgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0TWF4U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWluaW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRNaW5TaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWluU2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cbiAgICAgKi9cbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gdGl0bGUgLSBUaGUgZGVzaXJlZCB0aXRsZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFRpdGxlKHRpdGxlOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWwuXG4gICAgICovXG4gICAgU2V0Wm9vbSh6b29tOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNob3dNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaXplKCk6IFByb21pc2U8U2l6ZT4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNpemVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZ1bGxzY3JlZW4gYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZU1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmcmFtZWxlc3MgYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGcmFtZWxlc3MoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnJhbWVsZXNzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLW1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5NaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBXaWR0aCgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFdpZHRoTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBab29tcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFpvb20oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbUluKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21Jbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbU91dCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tT3V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tUmVzZXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbVJlc2V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUy9MaW51eCBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXG4gICAgICogR2F0aGVycyBpbmZvcm1hdGlvbiBhYm91dCB0aGUgZHJvcCB0YXJnZXQgZWxlbWVudCBhbmQgc2VuZHMgaXQgYmFjayB0byB0aGUgR28gYmFja2VuZC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cbiAgICAgKiBAcGFyYW0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQgKENTUyBwaXhlbHMpLlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudCAoQ1NTIHBpeGVscykuXG4gICAgICovXG4gICAgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgICAgICAvLyBDaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBpZ25vcmUgdGhlIGRyb3BcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudChlbGVtZW50KTtcblxuICAgICAgICBpZiAoIWRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIC8vIERyb3Agd2FzIG5vdCBvbiBhIGRlc2lnbmF0ZWQgZHJvcCB0YXJnZXQgLSBpZ25vcmVcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnN0IGVsZW1lbnREZXRhaWxzID0ge1xuICAgICAgICAgICAgaWQ6IGRyb3BUYXJnZXQuaWQsXG4gICAgICAgICAgICBjbGFzc0xpc3Q6IEFycmF5LmZyb20oZHJvcFRhcmdldC5jbGFzc0xpc3QpLFxuICAgICAgICAgICAgYXR0cmlidXRlczoge30gYXMgeyBba2V5OiBzdHJpbmddOiBzdHJpbmcgfSxcbiAgICAgICAgfTtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgICAgIGNvbnN0IGF0dHIgPSBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXNbaV07XG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscy5hdHRyaWJ1dGVzW2F0dHIubmFtZV0gPSBhdHRyLnZhbHVlO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgICAgIGZpbGVuYW1lcyxcbiAgICAgICAgICAgIHgsXG4gICAgICAgICAgICB5LFxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMsXG4gICAgICAgIH07XG5cbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dKEZpbGVzRHJvcHBlZCwgcGF5bG9hZCk7XG4gICAgICAgIFxuICAgICAgICAvLyBDbGVhbiB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhZnRlciBkcm9wXG4gICAgICAgIGNsZWFudXBOYXRpdmVEcmFnKCk7XG4gICAgfVxuICBcbiAgICAvKipcbiAgICAgKiBNb3ZlcyB0aGUgd2luZG93IHRvIHRoZSBjZW50ZXIgb2YgdGhlIHNwZWNpZmllZCBzY3JlZW4ncyB3b3JrIGFyZWEuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gc2NyZWVuSUQgLSBUaGUgSUQgb2YgdGhlIHRhcmdldCBzY3JlZW4uXG4gICAgICovXG4gICAgU2V0U2NyZWVuKHNjcmVlbklEOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTY3JlZW5NZXRob2QsIHsgc2NyZWVuSUQgfSk7XG4gICAgfVxuXG4gICAgLyogVHJpZ2dlcnMgV2luZG93cyAxMSBTbmFwIEFzc2lzdCBmZWF0dXJlIChXaW5kb3dzIG9ubHkpLlxuICAgICAqIFRoaXMgaXMgZXF1aXZhbGVudCB0byBwcmVzc2luZyBXaW4rWiBhbmQgc2hvd3Mgc25hcCBsYXlvdXQgb3B0aW9ucy5cbiAgICAgKi9cbiAgICBTbmFwQXNzaXN0KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNuYXBBc3Npc3RNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBwcmludCBkaWFsb2cgZm9yIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUHJpbnQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUHJpbnRNZXRob2QpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXG4gKi9cbmNvbnN0IHRoaXNXaW5kb3cgPSBuZXcgV2luZG93KCcnKTtcblxuLyoqXG4gKiBTZXRzIHVwIGdsb2JhbCBkcmFnIGFuZCBkcm9wIGV2ZW50IGxpc3RlbmVycyBmb3IgZmlsZSBkcm9wcy5cbiAqIEhhbmRsZXMgdmlzdWFsIGZlZWRiYWNrIChob3ZlciBzdGF0ZSkgYW5kIGZpbGUgZHJvcCBwcm9jZXNzaW5nLlxuICovXG5mdW5jdGlvbiBzZXR1cERyb3BUYXJnZXRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZG9jRWxlbWVudCA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudDtcbiAgICBsZXQgZHJhZ0VudGVyQ291bnRlciA9IDA7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdlbnRlcicsIChldmVudCkgPT4ge1xuICAgICAgICBpZiAoIWV2ZW50LmRhdGFUcmFuc2Zlcj8udHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcbiAgICAgICAgICAgIHJldHVybjsgLy8gT25seSBoYW5kbGUgZmlsZSBkcmFncywgbGV0IG90aGVyIGRyYWdzIHBhc3MgdGhyb3VnaFxuICAgICAgICB9XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIEFsd2F5cyBwcmV2ZW50IGRlZmF1bHQgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnOyAvLyBTaG93IFwibm8gZHJvcFwiIGN1cnNvclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcbiAgICAgICAgfVxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyKys7XG4gICAgICAgIFxuICAgICAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZKTtcbiAgICAgICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xuXG4gICAgICAgIC8vIFVwZGF0ZSBob3ZlciBzdGF0ZVxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmIChkcm9wVGFyZ2V0KSB7XG4gICAgICAgICAgICBkcm9wVGFyZ2V0LmNsYXNzTGlzdC5hZGQoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBkcm9wVGFyZ2V0O1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdvdmVyJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyYWdzXG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7IC8vIFNob3cgXCJubyBkcm9wXCIgY3Vyc29yXG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICAvLyBVcGRhdGUgZHJvcCB0YXJnZXQgYXMgY3Vyc29yIG1vdmVzXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XG4gICAgICAgIFxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgaWYgKGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGlmICghZHJvcFRhcmdldC5jbGFzc0xpc3QuY29udGFpbnMoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKSkge1xuICAgICAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIEFsd2F5cyBwcmV2ZW50IGRlZmF1bHQgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgLy8gT24gTGludXgvV2ViS2l0R1RLIGFuZCBtYWNPUywgZHJhZ2xlYXZlIGZpcmVzIGltbWVkaWF0ZWx5IHdpdGggcmVsYXRlZFRhcmdldD1udWxsIHdoZW4gbmF0aXZlXG4gICAgICAgIC8vIGRyYWcgaGFuZGxpbmcgaXMgaW52b2x2ZWQuIElnbm9yZSB0aGVzZSBzcHVyaW91cyBldmVudHMgLSB3ZSdsbCBjbGVhbiB1cCBvbiBkcm9wIGluc3RlYWQuXG4gICAgICAgIGlmIChldmVudC5yZWxhdGVkVGFyZ2V0ID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXItLTtcbiAgICAgICAgXG4gICAgICAgIGlmIChkcmFnRW50ZXJDb3VudGVyID09PSAwIHx8IFxuICAgICAgICAgICAgKGN1cnJlbnREcm9wVGFyZ2V0ICYmICFjdXJyZW50RHJvcFRhcmdldC5jb250YWlucyhldmVudC5yZWxhdGVkVGFyZ2V0IGFzIE5vZGUpKSkge1xuICAgICAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XG4gICAgICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwO1xuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyb3BzXG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwO1xuICAgICAgICBcbiAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBoYW5kbGUgZmlsZSBkcm9wcyB2aWEgSmF2YVNjcmlwdFxuICAgICAgICAvLyBPbiBtYWNPUy9MaW51eCwgbmF0aXZlIGNvZGUgd2lsbCBjYWxsIEhhbmRsZVBsYXRmb3JtRmlsZURyb3BcbiAgICAgICAgaWYgKGNhblJlc29sdmVGaWxlUGF0aHMoKSkge1xuICAgICAgICAgICAgY29uc3QgZmlsZXM6IEZpbGVbXSA9IFtdO1xuICAgICAgICAgICAgaWYgKGV2ZW50LmRhdGFUcmFuc2Zlci5pdGVtcykge1xuICAgICAgICAgICAgICAgIGZvciAoY29uc3QgaXRlbSBvZiBldmVudC5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKGl0ZW0ua2luZCA9PT0gJ2ZpbGUnKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCBmaWxlID0gaXRlbS5nZXRBc0ZpbGUoKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChmaWxlKSBmaWxlcy5wdXNoKGZpbGUpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSBlbHNlIGlmIChldmVudC5kYXRhVHJhbnNmZXIuZmlsZXMpIHtcbiAgICAgICAgICAgICAgICBmb3IgKGNvbnN0IGZpbGUgb2YgZXZlbnQuZGF0YVRyYW5zZmVyLmZpbGVzKSB7XG4gICAgICAgICAgICAgICAgICAgIGZpbGVzLnB1c2goZmlsZSk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgXG4gICAgICAgICAgICBpZiAoZmlsZXMubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgICAgIHJlc29sdmVGaWxlUGF0aHMoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZmlsZXMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xufVxuXG4vLyBJbml0aWFsaXplIGxpc3RlbmVycyB3aGVuIHRoZSBzY3JpcHQgbG9hZHNcbmlmICh0eXBlb2Ygd2luZG93ICE9PSBcInVuZGVmaW5lZFwiICYmIHR5cGVvZiBkb2N1bWVudCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgIHNldHVwRHJvcFRhcmdldExpc3RlbmVycygpO1xufVxuXG5leHBvcnQgZGVmYXVsdCB0aGlzV2luZG93O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgKiBhcyBSdW50aW1lIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyY1wiO1xuXG4vLyBOT1RFOiB0aGUgZm9sbG93aW5nIG1ldGhvZHMgTVVTVCBiZSBpbXBvcnRlZCBleHBsaWNpdGx5IGJlY2F1c2Ugb2YgaG93IGVzYnVpbGQgaW5qZWN0aW9uIHdvcmtzXG5pbXBvcnQgeyBFbmFibGUgYXMgRW5hYmxlV01MIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dtbFwiO1xuaW1wb3J0IHsgZGVidWdMb2cgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHNcIjtcblxud2luZG93LndhaWxzID0gUnVudGltZTtcbkVuYWJsZVdNTCgpO1xuXG5pZiAoREVCVUcpIHtcbiAgICBkZWJ1Z0xvZyhcIldhaWxzIFJ1bnRpbWUgTG9hZGVkXCIpXG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TeXN0ZW0pO1xuXG5jb25zdCBTeXN0ZW1Jc0RhcmtNb2RlID0gMDtcbmNvbnN0IFN5c3RlbUVudmlyb25tZW50ID0gMTtcbmNvbnN0IFN5c3RlbUNhcGFiaWxpdGllcyA9IDI7XG5cbmNvbnN0IF9pbnZva2UgPSAoZnVuY3Rpb24gKCkge1xuICAgIHRyeSB7XG4gICAgICAgIC8vIFdpbmRvd3MgV2ViVmlldzJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlKSB7XG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3KTtcbiAgICAgICAgfVxuICAgICAgICAvLyBtYWNPUy9pT1MgV0tXZWJWaWV3XG4gICAgICAgIGVsc2UgaWYgKCh3aW5kb3cgYXMgYW55KS53ZWJraXQ/Lm1lc3NhZ2VIYW5kbGVycz8uWydleHRlcm5hbCddPy5wb3N0TWVzc2FnZSkge1xuICAgICAgICAgICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddLnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10pO1xuICAgICAgICB9XG4gICAgICAgIC8vIEFuZHJvaWQgV2ViVmlldyAtIHVzZXMgYWRkSmF2YXNjcmlwdEludGVyZmFjZSB3aGljaCBleHBvc2VzIHdpbmRvdy53YWlscy5pbnZva2VcbiAgICAgICAgZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndhaWxzPy5pbnZva2UpIHtcbiAgICAgICAgICAgIHJldHVybiAobXNnOiBhbnkpID0+ICh3aW5kb3cgYXMgYW55KS53YWlscy5pbnZva2UodHlwZW9mIG1zZyA9PT0gJ3N0cmluZycgPyBtc2cgOiBKU09OLnN0cmluZ2lmeShtc2cpKTtcbiAgICAgICAgfVxuICAgIH0gY2F0Y2goZSkge31cblxuICAgIGNvbnNvbGUud2FybignXFxuJWNcdTI2QTBcdUZFMEYgQnJvd3NlciBFbnZpcm9ubWVudCBEZXRlY3RlZCAlY1xcblxcbiVjT25seSBVSSBwcmV2aWV3cyBhcmUgYXZhaWxhYmxlIGluIHRoZSBicm93c2VyLiBGb3IgZnVsbCBmdW5jdGlvbmFsaXR5LCBwbGVhc2UgcnVuIHRoZSBhcHBsaWNhdGlvbiBpbiBkZXNrdG9wIG1vZGUuXFxuTW9yZSBpbmZvcm1hdGlvbiBhdDogaHR0cHM6Ly92My53YWlscy5pby9sZWFybi9idWlsZC8jdXNpbmctYS1icm93c2VyLWZvci1kZXZlbG9wbWVudFxcbicsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjZmZmZmZmOyBjb2xvcjogIzAwMDAwMDsgZm9udC13ZWlnaHQ6IGJvbGQ7IHBhZGRpbmc6IDRweCA4cHg7IGJvcmRlci1yYWRpdXM6IDRweDsgYm9yZGVyOiAycHggc29saWQgIzAwMDAwMDsnLFxuICAgICAgICAnYmFja2dyb3VuZDogdHJhbnNwYXJlbnQ7JyxcbiAgICAgICAgJ2NvbG9yOiAjZmZmZmZmOyBmb250LXN0eWxlOiBpdGFsaWM7IGZvbnQtd2VpZ2h0OiBib2xkOycpO1xuICAgIHJldHVybiBudWxsO1xufSkoKTtcblxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2c6IGFueSk6IHZvaWQge1xuICAgIF9pbnZva2U/Lihtc2cpO1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYSBib29sZWFuIHZhbHVlIGluZGljYXRpbmcgaWYgdGhlIHN5c3RlbSBpcyBpbiBkYXJrIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgIHJldHVybiBjYWxsKFN5c3RlbUlzRGFya01vZGUpO1xufVxuXG4vKipcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIHRoZSBjYXBhYmlsaXRpZXMuXG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKTogUHJvbWlzZTxSZWNvcmQ8c3RyaW5nLCBhbnk+PiB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtQ2FwYWJpbGl0aWVzKTtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBPU0luZm8ge1xuICAgIC8qKiBUaGUgYnJhbmRpbmcgb2YgdGhlIE9TLiAqL1xuICAgIEJyYW5kaW5nOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBJRCBvZiB0aGUgT1MuICovXG4gICAgSUQ6IHN0cmluZztcbiAgICAvKiogVGhlIG5hbWUgb2YgdGhlIE9TLiAqL1xuICAgIE5hbWU6IHN0cmluZztcbiAgICAvKiogVGhlIHZlcnNpb24gb2YgdGhlIE9TLiAqL1xuICAgIFZlcnNpb246IHN0cmluZztcbn1cblxuZXhwb3J0IGludGVyZmFjZSBFbnZpcm9ubWVudEluZm8ge1xuICAgIC8qKiBUaGUgYXJjaGl0ZWN0dXJlIG9mIHRoZSBzeXN0ZW0uICovXG4gICAgQXJjaDogc3RyaW5nO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS4gKi9cbiAgICBEZWJ1ZzogYm9vbGVhbjtcbiAgICAvKiogVGhlIG9wZXJhdGluZyBzeXN0ZW0gaW4gdXNlLiAqL1xuICAgIE9TOiBzdHJpbmc7XG4gICAgLyoqIERldGFpbHMgb2YgdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXG4gICAgT1NJbmZvOiBPU0luZm87XG4gICAgLyoqIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uICovXG4gICAgUGxhdGZvcm1JbmZvOiBSZWNvcmQ8c3RyaW5nLCBhbnk+O1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIE9TIGFuZCBzeXN0ZW0gYXJjaGl0ZWN0dXJlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKTogUHJvbWlzZTxFbnZpcm9ubWVudEluZm8+IHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1FbnZpcm9ubWVudCk7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cbiAqXG4gKiBAcmV0dXJuIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwid2luZG93c1wiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LlxuICpcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwibGludXhcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIG1hY09TLCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TID09PSBcImRhcndpblwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFtZDY0XCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0uXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FSTSgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYXJtXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZS5cbiAqXG4gKiBAcmV0dXJucyBSZXR1cm5zIHRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIEFSTTY0IGFyY2hpdGVjdHVyZSwgb3RoZXJ3aXNlIHJldHVybnMgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhcm02NFwiO1xufVxuXG4vKipcbiAqIFJlcG9ydHMgd2hldGhlciB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGFwcCBpcyBiZWluZyBydW4gaW4gZGVidWcgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGVidWcoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuIEJvb2xlYW4oKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkRlYnVnKTtcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IElzRGVidWcgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCB7IGV2ZW50VGFyZ2V0IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcblxuLy8gc2V0dXBcbndpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcblxuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkOiBzdHJpbmcsIHg6IG51bWJlciwgeTogbnVtYmVyLCBkYXRhOiBhbnkpOiB2b2lkIHtcbiAgICB2b2lkIGNhbGwoQ29udGV4dE1lbnVPcGVuLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xuXG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcbiAgICBjb25zdCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpLnRyaW0oKTtcblxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBjb25zdCBkYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnUtZGF0YVwiKTtcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBkYXRhKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50LCB0YXJnZXQpO1xuICAgIH1cbn1cblxuXG4vKlxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBhdXRvOyAoZGVmYXVsdCkgd2lsbCBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudSBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZSBPUiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkIE9SIGVsZW1lbnQgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcblxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcbiovXG5mdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50OiBNb3VzZUV2ZW50LCB0YXJnZXQ6IEhUTUxFbGVtZW50KSB7XG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XG4gICAgaWYgKElzRGVidWcoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIHN3aXRjaCAod2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKSkge1xuICAgICAgICBjYXNlICdzaG93JzpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSAnaGlkZSc6XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgaWYgKHRhcmdldC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxuICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSBzZWxlY3Rpb24gJiYgc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMDtcbiAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcbiAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XG4gICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IHRhcmdldCkge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgdGFnIGlzIGlucHV0IG9yIHRleHRhcmVhLlxuICAgIGlmICh0YXJnZXQgaW5zdGFuY2VvZiBIVE1MSW5wdXRFbGVtZW50IHx8IHRhcmdldCBpbnN0YW5jZW9mIEhUTUxUZXh0QXJlYUVsZW1lbnQpIHtcbiAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIXRhcmdldC5yZWFkT25seSAmJiAhdGFyZ2V0LmRpc2FibGVkKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXG4gKlxuICogQHBhcmFtIGtleSAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cbiAqIEByZXR1cm4gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5OiBzdHJpbmcpOiBhbnkge1xuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleV07XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJVbmFibGUgdG8gcmV0cmlldmUgZmxhZyAnXCIgKyBrZXkgKyBcIic6IFwiICsgZSwgeyBjYXVzZTogZSB9KTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IGludm9rZSwgSXNMaW51eCwgSXNXaW5kb3dzIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyBHZXRGbGFnIH0gZnJvbSBcIi4vZmxhZ3MuanNcIjtcbmltcG9ydCB7IGNhblRyYWNrQnV0dG9ucywgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuXG4vLyBTZXR1cFxubGV0IGNhbkRyYWcgPSBmYWxzZTtcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xuXG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XG5sZXQgcmVzaXppbmcgPSBmYWxzZTtcbmxldCByZXNpemVFZGdlOiBzdHJpbmcgPSBcIlwiO1xubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcblxubGV0IGJ1dHRvbnMgPSAwO1xuY29uc3QgYnV0dG9uc1RyYWNrZWQgPSBjYW5UcmFja0J1dHRvbnMoKTtcblxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9ICh2YWx1ZTogYm9vbGVhbik6IHZvaWQgPT4ge1xuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xuICAgIGlmICghcmVzaXphYmxlKSB7XG4gICAgICAgIC8vIFN0b3AgcmVzaXppbmcgaWYgaW4gcHJvZ3Jlc3MuXG4gICAgICAgIGNhblJlc2l6ZSA9IHJlc2l6aW5nID0gZmFsc2U7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbn07XG5cbi8vIERlZmVyIGF0dGFjaGluZyBtb3VzZSBsaXN0ZW5lcnMgdW50aWwgd2Uga25vdyB3ZSdyZSBub3Qgb24gbW9iaWxlLlxubGV0IGRyYWdJbml0RG9uZSA9IGZhbHNlO1xuZnVuY3Rpb24gaXNNb2JpbGUoKTogYm9vbGVhbiB7XG4gICAgY29uc3Qgb3MgPSAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1M7XG4gICAgaWYgKG9zID09PSBcImlvc1wiIHx8IG9zID09PSBcImFuZHJvaWRcIikgcmV0dXJuIHRydWU7XG4gICAgLy8gRmFsbGJhY2sgaGV1cmlzdGljIGlmIGVudmlyb25tZW50IG5vdCB5ZXQgc2V0XG4gICAgY29uc3QgdWEgPSBuYXZpZ2F0b3IudXNlckFnZW50IHx8IG5hdmlnYXRvci52ZW5kb3IgfHwgKHdpbmRvdyBhcyBhbnkpLm9wZXJhIHx8IFwiXCI7XG4gICAgcmV0dXJuIC9hbmRyb2lkfGlwaG9uZXxpcGFkfGlwb2R8aWVtb2JpbGV8d3BkZXNrdG9wL2kudGVzdCh1YSk7XG59XG5mdW5jdGlvbiB0cnlJbml0RHJhZ0hhbmRsZXJzKCk6IHZvaWQge1xuICAgIGlmIChkcmFnSW5pdERvbmUpIHJldHVybjtcbiAgICBpZiAoaXNNb2JpbGUoKSkgcmV0dXJuO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICBmb3IgKGNvbnN0IGV2IG9mIFsnY2xpY2snLCAnY29udGV4dG1lbnUnLCAnZGJsY2xpY2snXSkge1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihldiwgc3VwcHJlc3NFdmVudCwgeyBjYXB0dXJlOiB0cnVlIH0pO1xuICAgIH1cbiAgICBkcmFnSW5pdERvbmUgPSB0cnVlO1xufVxuLy8gQXR0ZW1wdCBpbW1lZGlhdGUgaW5pdCAoaW4gY2FzZSBlbnZpcm9ubWVudCBhbHJlYWR5IHByZXNlbnQpXG50cnlJbml0RHJhZ0hhbmRsZXJzKCk7XG4vLyBBbHNvIGF0dGVtcHQgb24gRE9NIHJlYWR5XG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgdHJ5SW5pdERyYWdIYW5kbGVycywgeyBvbmNlOiB0cnVlIH0pO1xuLy8gQXMgYSBsYXN0IHJlc29ydCwgcG9sbCBmb3IgZW52aXJvbm1lbnQgZm9yIGEgc2hvcnQgcGVyaW9kXG5sZXQgZHJhZ0VudlBvbGxzID0gMDtcbmNvbnN0IGRyYWdFbnZQb2xsID0gd2luZG93LnNldEludGVydmFsKCgpID0+IHtcbiAgICBpZiAoZHJhZ0luaXREb25lKSB7IHdpbmRvdy5jbGVhckludGVydmFsKGRyYWdFbnZQb2xsKTsgcmV0dXJuOyB9XG4gICAgdHJ5SW5pdERyYWdIYW5kbGVycygpO1xuICAgIGlmICgrK2RyYWdFbnZQb2xscyA+IDEwMCkgeyB3aW5kb3cuY2xlYXJJbnRlcnZhbChkcmFnRW52UG9sbCk7IH1cbn0sIDUwKTtcblxuZnVuY3Rpb24gc3VwcHJlc3NFdmVudChldmVudDogRXZlbnQpIHtcbiAgICAvLyBTdXBwcmVzcyBjbGljayBldmVudHMgd2hpbGUgcmVzaXppbmcgb3IgZHJhZ2dpbmcuXG4gICAgaWYgKGRyYWdnaW5nIHx8IHJlc2l6aW5nKSB7XG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG59XG5cbi8vIFVzZSBjb25zdGFudHMgdG8gYXZvaWQgY29tcGFyaW5nIHN0cmluZ3MgbXVsdGlwbGUgdGltZXMuXG5jb25zdCBNb3VzZURvd24gPSAwO1xuY29uc3QgTW91c2VVcCAgID0gMTtcbmNvbnN0IE1vdXNlTW92ZSA9IDI7XG5cbmZ1bmN0aW9uIHVwZGF0ZShldmVudDogTW91c2VFdmVudCkge1xuICAgIC8vIFdpbmRvd3Mgc3VwcHJlc3NlcyBtb3VzZSBldmVudHMgYXQgdGhlIGVuZCBvZiBkcmFnZ2luZyBvciByZXNpemluZyxcbiAgICAvLyBzbyB3ZSBuZWVkIHRvIGJlIHNtYXJ0IGFuZCBzeW50aGVzaXplIGJ1dHRvbiBldmVudHMuXG5cbiAgICBsZXQgZXZlbnRUeXBlOiBudW1iZXIsIGV2ZW50QnV0dG9ucyA9IGV2ZW50LmJ1dHRvbnM7XG4gICAgc3dpdGNoIChldmVudC50eXBlKSB7XG4gICAgICAgIGNhc2UgJ21vdXNlZG93bic6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZURvd247XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgfCAoMSA8PCBldmVudC5idXR0b24pOyB9XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgY2FzZSAnbW91c2V1cCc6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZVVwO1xuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zICYgfigxIDw8IGV2ZW50LmJ1dHRvbik7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VNb3ZlO1xuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zOyB9XG4gICAgICAgICAgICBicmVhaztcbiAgICB9XG5cbiAgICBsZXQgcmVsZWFzZWQgPSBidXR0b25zICYgfmV2ZW50QnV0dG9ucztcbiAgICBsZXQgcHJlc3NlZCA9IGV2ZW50QnV0dG9ucyAmIH5idXR0b25zO1xuXG4gICAgYnV0dG9ucyA9IGV2ZW50QnV0dG9ucztcblxuICAgIC8vIFN5bnRoZXNpemUgYSByZWxlYXNlLXByZXNzIHNlcXVlbmNlIGlmIHdlIGRldGVjdCBhIHByZXNzIG9mIGFuIGFscmVhZHkgcHJlc3NlZCBidXR0b24uXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VEb3duICYmICEocHJlc3NlZCAmIGV2ZW50LmJ1dHRvbikpIHtcbiAgICAgICAgcmVsZWFzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcbiAgICAgICAgcHJlc3NlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xuICAgIH1cblxuICAgIC8vIFN1cHByZXNzIGFsbCBidXR0b24gZXZlbnRzIGR1cmluZyBkcmFnZ2luZyBhbmQgcmVzaXppbmcsXG4gICAgLy8gdW5sZXNzIHRoaXMgaXMgYSBtb3VzZXVwIGV2ZW50IHRoYXQgaXMgZW5kaW5nIGEgZHJhZyBhY3Rpb24uXG4gICAgaWYgKFxuICAgICAgICBldmVudFR5cGUgIT09IE1vdXNlTW92ZSAvLyBGYXN0IHBhdGggZm9yIG1vdXNlbW92ZVxuICAgICAgICAmJiByZXNpemluZ1xuICAgICAgICB8fCAoXG4gICAgICAgICAgICBkcmFnZ2luZ1xuICAgICAgICAgICAgJiYgKFxuICAgICAgICAgICAgICAgIGV2ZW50VHlwZSA9PT0gTW91c2VEb3duXG4gICAgICAgICAgICAgICAgfHwgZXZlbnQuYnV0dG9uICE9PSAwXG4gICAgICAgICAgICApXG4gICAgICAgIClcbiAgICApIHtcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cblxuICAgIC8vIEhhbmRsZSByZWxlYXNlc1xuICAgIGlmIChyZWxlYXNlZCAmIDEpIHsgcHJpbWFyeVVwKGV2ZW50KTsgfVxuICAgIC8vIEhhbmRsZSBwcmVzc2VzXG4gICAgaWYgKHByZXNzZWQgJiAxKSB7IHByaW1hcnlEb3duKGV2ZW50KTsgfVxuXG4gICAgLy8gSGFuZGxlIG1vdXNlbW92ZVxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlTW92ZSkgeyBvbk1vdXNlTW92ZShldmVudCk7IH07XG59XG5cbmZ1bmN0aW9uIHByaW1hcnlEb3duKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XG4gICAgLy8gUmVzZXQgcmVhZGluZXNzIHN0YXRlLlxuICAgIGNhbkRyYWcgPSBmYWxzZTtcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcblxuICAgIC8vIElnbm9yZSByZXBlYXRlZCBjbGlja3Mgb24gbWFjT1MgYW5kIExpbnV4LlxuICAgIGlmICghSXNXaW5kb3dzKCkpIHtcbiAgICAgICAgaWYgKGV2ZW50LnR5cGUgPT09ICdtb3VzZWRvd24nICYmIGV2ZW50LmJ1dHRvbiA9PT0gMCAmJiBldmVudC5kZXRhaWwgIT09IDEpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIGlmIChyZXNpemVFZGdlKSB7XG4gICAgICAgIC8vIFJlYWR5IHRvIHJlc2l6ZSBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lLlxuICAgICAgICBjYW5SZXNpemUgPSB0cnVlO1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZyBvcGVyYXRpb25zIHdoZW4gb24gcmVzaXplIGVkZ2VzLlxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUmV0cmlldmUgdGFyZ2V0IGVsZW1lbnRcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XG5cbiAgICAvLyBSZWFkeSB0byBkcmFnIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUgb24gYSBkcmFnZ2FibGUgZWxlbWVudC5cbiAgICAvLyBJZ25vcmUgY2xpY2tzIG9uIHRoZSBzY3JvbGxiYXIuXG4gICAgY29uc3Qgc3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpO1xuICAgIGNhbkRyYWcgPSAoXG4gICAgICAgIHN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLXdhaWxzLWRyYWdnYWJsZVwiKS50cmltKCkgPT09IFwiZHJhZ1wiXG4gICAgICAgICYmIChcbiAgICAgICAgICAgIGV2ZW50Lm9mZnNldFggLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdMZWZ0KSA8IHRhcmdldC5jbGllbnRXaWR0aFxuICAgICAgICAgICAgJiYgZXZlbnQub2Zmc2V0WSAtIHBhcnNlRmxvYXQoc3R5bGUucGFkZGluZ1RvcCkgPCB0YXJnZXQuY2xpZW50SGVpZ2h0XG4gICAgICAgIClcbiAgICApO1xufVxuXG5mdW5jdGlvbiBwcmltYXJ5VXAoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICAvLyBTdG9wIGRyYWdnaW5nIGFuZCByZXNpemluZy5cbiAgICBjYW5EcmFnID0gZmFsc2U7XG4gICAgZHJhZ2dpbmcgPSBmYWxzZTtcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcbiAgICByZXNpemluZyA9IGZhbHNlO1xufVxuXG5jb25zdCBjdXJzb3JGb3JFZGdlID0gT2JqZWN0LmZyZWV6ZSh7XG4gICAgXCJzZS1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxuICAgIFwic3ctcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcbiAgICBcIm53LXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXG4gICAgXCJuZS1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxuICAgIFwidy1yZXNpemVcIjogXCJldy1yZXNpemVcIixcbiAgICBcIm4tcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXG4gICAgXCJzLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxuICAgIFwiZS1yZXNpemVcIjogXCJldy1yZXNpemVcIixcbn0pXG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShlZGdlPzoga2V5b2YgdHlwZW9mIGN1cnNvckZvckVkZ2UpOiB2b2lkIHtcbiAgICBpZiAoZWRnZSkge1xuICAgICAgICBpZiAoIXJlc2l6ZUVkZ2UpIHsgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yOyB9XG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gY3Vyc29yRm9yRWRnZVtlZGdlXTtcbiAgICB9IGVsc2UgaWYgKCFlZGdlICYmIHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBkZWZhdWx0Q3Vyc29yO1xuICAgIH1cblxuICAgIHJlc2l6ZUVkZ2UgPSBlZGdlIHx8IFwiXCI7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XG4gICAgaWYgKGNhblJlc2l6ZSAmJiByZXNpemVFZGdlKSB7XG4gICAgICAgIC8vIFN0YXJ0IHJlc2l6aW5nLlxuICAgICAgICByZXNpemluZyA9IHRydWU7XG4gICAgICAgIGludm9rZShcIndhaWxzOnJlc2l6ZTpcIiArIHJlc2l6ZUVkZ2UpO1xuICAgIH0gZWxzZSBpZiAoY2FuRHJhZykge1xuICAgICAgICAvLyBTdGFydCBkcmFnZ2luZy5cbiAgICAgICAgZHJhZ2dpbmcgPSB0cnVlO1xuICAgICAgICBpbnZva2UoXCJ3YWlsczpkcmFnXCIpO1xuICAgIH1cblxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xuICAgICAgICAvLyBFaXRoZXIgZHJhZyBvciByZXNpemUgaXMgb25nb2luZyxcbiAgICAgICAgLy8gcmVzZXQgcmVhZGluZXNzIGFuZCBzdG9wIHByb2Nlc3NpbmcuXG4gICAgICAgIGNhbkRyYWcgPSBjYW5SZXNpemUgPSBmYWxzZTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghcmVzaXphYmxlIHx8ICghSXNXaW5kb3dzKCkgJiYgIUlzTGludXgoKSkpIHtcbiAgICAgICAgaWYgKHJlc2l6ZUVkZ2UpIHsgc2V0UmVzaXplKCk7IH1cbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZUhlaWdodCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlSGVpZ2h0XCIpIHx8IDU7XG4gICAgY29uc3QgcmVzaXplSGFuZGxlV2lkdGggPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZVdpZHRoXCIpIHx8IDU7XG5cbiAgICAvLyBFeHRyYSBwaXhlbHMgZm9yIHRoZSBjb3JuZXIgYXJlYXMuXG4gICAgY29uc3QgY29ybmVyRXh0cmEgPSBHZXRGbGFnKFwicmVzaXplQ29ybmVyRXh0cmFcIikgfHwgMTA7XG5cbiAgICBjb25zdCByaWdodEJvcmRlciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgY29uc3QgbGVmdEJvcmRlciA9IGV2ZW50LmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBjb25zdCB0b3BCb3JkZXIgPSBldmVudC5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuICAgIGNvbnN0IGJvdHRvbUJvcmRlciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcblxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyIGFyZWFzLlxuICAgIGNvbnN0IHJpZ2h0Q29ybmVyID0gKHdpbmRvdy5vdXRlcldpZHRoIC0gZXZlbnQuY2xpZW50WCkgPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgbGVmdENvcm5lciA9IGV2ZW50LmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgdG9wQ29ybmVyID0gZXZlbnQuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgYm90dG9tQ29ybmVyID0gKHdpbmRvdy5vdXRlckhlaWdodCAtIGV2ZW50LmNsaWVudFkpIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcblxuICAgIGlmICghbGVmdENvcm5lciAmJiAhdG9wQ29ybmVyICYmICFib3R0b21Db3JuZXIgJiYgIXJpZ2h0Q29ybmVyKSB7XG4gICAgICAgIC8vIE9wdGltaXNhdGlvbjogb3V0IG9mIGFsbCBjb3JuZXIgYXJlYXMgaW1wbGllcyBvdXQgb2YgYm9yZGVycy5cbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxuICAgIC8vIERldGVjdCBjb3JuZXJzLlxuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICAvLyBEZXRlY3QgYm9yZGVycy5cbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xuICAgIC8vIE91dCBvZiBib3JkZXIgYXJlYS5cbiAgICBlbHNlIHNldFJlc2l6ZSgpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoU2hvd01ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgQ2FuY2VsbGFibGVQcm9taXNlLCB0eXBlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnMgfSBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcblxuLy8gU2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xuXG50eXBlIFByb21pc2VSZXNvbHZlcnMgPSBPbWl0PENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8YW55PiwgXCJwcm9taXNlXCIgfCBcIm9uY2FuY2VsbGVkXCI+XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbGwpO1xuY29uc3QgY2FuY2VsQ2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2FuY2VsQ2FsbCk7XG5jb25zdCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcDxzdHJpbmcsIFByb21pc2VSZXNvbHZlcnM+KCk7XG5cbmNvbnN0IENhbGxCaW5kaW5nID0gMDtcbmNvbnN0IENhbmNlbE1ldGhvZCA9IDBcblxuLyoqXG4gKiBIb2xkcyBhbGwgcmVxdWlyZWQgaW5mb3JtYXRpb24gZm9yIGEgYmluZGluZyBjYWxsLlxuICogTWF5IHByb3ZpZGUgZWl0aGVyIGEgbWV0aG9kIElEIG9yIGEgbWV0aG9kIG5hbWUsIGJ1dCBub3QgYm90aC5cbiAqL1xuZXhwb3J0IHR5cGUgQ2FsbE9wdGlvbnMgPSB7XG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2RJRDogbnVtYmVyO1xuICAgIC8qKiBUaGUgZnVsbHkgcXVhbGlmaWVkIG5hbWUgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZE5hbWU/OiBuZXZlcjtcbiAgICAvKiogQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuICovXG4gICAgYXJnczogYW55W107XG59IHwge1xuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kSUQ/OiBuZXZlcjtcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2ROYW1lOiBzdHJpbmc7XG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xuICAgIGFyZ3M6IGFueVtdO1xufTtcblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHRocm93biBpbiBjYXNlIHRoZSBib3VuZCBtZXRob2QgcmV0dXJucyBhbiBlcnJvci5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIFJ1bnRpbWVFcnJvciNuYW1lfSBwcm9wZXJ0eSBpcyBcIlJ1bnRpbWVFcnJvclwiLlxuICovXG5leHBvcnQgY2xhc3MgUnVudGltZUVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgUnVudGltZUVycm9yIGluc3RhbmNlLlxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gICAgICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIHRvIGJlIGZvcndhcmRlZCB0byB0aGUgRXJyb3IgY29uc3RydWN0b3IuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xuICAgICAgICBzdXBlcihtZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJSdW50aW1lRXJyb3JcIjtcbiAgICB9XG59XG5cbi8qKlxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cbiAqXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIENhbGwgYSBib3VuZCBtZXRob2QgYWNjb3JkaW5nIHRvIHRoZSBnaXZlbiBjYWxsIG9wdGlvbnMuXG4gKlxuICogSW4gY2FzZSBvZiBmYWlsdXJlLCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGFuIGV4Y2VwdGlvblxuICogYW1vbmcgUmVmZXJlbmNlRXJyb3IgKHVua25vd24gbWV0aG9kKSwgVHlwZUVycm9yICh3cm9uZyBhcmd1bWVudCBjb3VudCBvciB0eXBlKSxcbiAqIHtAbGluayBSdW50aW1lRXJyb3J9IChtZXRob2QgcmV0dXJuZWQgYW4gZXJyb3IpLCBvciBvdGhlciAobmV0d29yayBvciBpbnRlcm5hbCBlcnJvcnMpLlxuICogVGhlIGV4Y2VwdGlvbiBtaWdodCBoYXZlIGEgXCJjYXVzZVwiIGZpZWxkIHdpdGggdGhlIHZhbHVlIHJldHVybmVkXG4gKiBieSB0aGUgYXBwbGljYXRpb24tIG9yIHNlcnZpY2UtbGV2ZWwgZXJyb3IgbWFyc2hhbGluZyBmdW5jdGlvbnMuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBBIG1ldGhvZCBjYWxsIGRlc2NyaXB0b3IuXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zOiBDYWxsT3B0aW9ucyk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcblxuICAgIGNvbnN0IHJlc3VsdCA9IENhbmNlbGxhYmxlUHJvbWlzZS53aXRoUmVzb2x2ZXJzPGFueT4oKTtcbiAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlOiByZXN1bHQucmVzb2x2ZSwgcmVqZWN0OiByZXN1bHQucmVqZWN0IH0pO1xuXG4gICAgY29uc3QgcmVxdWVzdCA9IGNhbGwoQ2FsbEJpbmRpbmcsIE9iamVjdC5hc3NpZ24oeyBcImNhbGwtaWRcIjogaWQgfSwgb3B0aW9ucykpO1xuICAgIGxldCBydW5uaW5nID0gdHJ1ZTtcblxuICAgIHJlcXVlc3QudGhlbigocmVzKSA9PiB7XG4gICAgICAgIHJ1bm5pbmcgPSBmYWxzZTtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXN1bHQucmVzb2x2ZShyZXMpO1xuICAgIH0sIChlcnIpID0+IHtcbiAgICAgICAgcnVubmluZyA9IGZhbHNlO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHJlc3VsdC5yZWplY3QoZXJyKTtcbiAgICB9KTtcblxuICAgIGNvbnN0IGNhbmNlbCA9ICgpID0+IHtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXR1cm4gY2FuY2VsQ2FsbChDYW5jZWxNZXRob2QsIHtcImNhbGwtaWRcIjogaWR9KS5jYXRjaCgoZXJyKSA9PiB7XG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igd2hpbGUgcmVxdWVzdGluZyBiaW5kaW5nIGNhbGwgY2FuY2VsbGF0aW9uOlwiLCBlcnIpO1xuICAgICAgICB9KTtcbiAgICB9O1xuXG4gICAgcmVzdWx0Lm9uY2FuY2VsbGVkID0gKCkgPT4ge1xuICAgICAgICBpZiAocnVubmluZykge1xuICAgICAgICAgICAgcmV0dXJuIGNhbmNlbCgpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcmV0dXJuIHJlcXVlc3QudGhlbihjYW5jZWwpO1xuICAgICAgICB9XG4gICAgfTtcblxuICAgIHJldHVybiByZXN1bHQucHJvbWlzZTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIGJvdW5kIG1ldGhvZCBieSBuYW1lIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKiBTZWUge0BsaW5rIENhbGx9IGZvciBkZXRhaWxzLlxuICpcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG1ldGhvZE5hbWU6IHN0cmluZywgLi4uYXJnczogYW55W10pOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2ROYW1lLCBhcmdzIH0pO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIGJ5IGl0cyBudW1lcmljIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKiBTZWUge0BsaW5rIENhbGx9IGZvciBkZXRhaWxzLlxuICpcbiAqIEBwYXJhbSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybiBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQ6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2RJRCwgYXJncyB9KTtcbn1cbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9pbnNwZWN0LWpzL2lzLWNhbGxhYmxlXG5cbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxuLy9cbi8vIENvcHlyaWdodCAoYykgMjAxNSBKb3JkYW4gSGFyYmFuZFxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHlcbi8vIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWxcbi8vIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHNcbi8vIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGxcbi8vIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpc1xuLy8gZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcbi8vXG4vLyBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXG4vL1xuLy8gVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksXG4vLyBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEVcbi8vIEFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVJcbi8vIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sXG4vLyBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRVxuLy8gU09GVFdBUkUuXG5cbnZhciBmblRvU3RyID0gRnVuY3Rpb24ucHJvdG90eXBlLnRvU3RyaW5nO1xudmFyIHJlZmxlY3RBcHBseTogdHlwZW9mIFJlZmxlY3QuYXBwbHkgfCBmYWxzZSB8IG51bGwgPSB0eXBlb2YgUmVmbGVjdCA9PT0gJ29iamVjdCcgJiYgUmVmbGVjdCAhPT0gbnVsbCAmJiBSZWZsZWN0LmFwcGx5O1xudmFyIGJhZEFycmF5TGlrZTogYW55O1xudmFyIGlzQ2FsbGFibGVNYXJrZXI6IGFueTtcbmlmICh0eXBlb2YgcmVmbGVjdEFwcGx5ID09PSAnZnVuY3Rpb24nICYmIHR5cGVvZiBPYmplY3QuZGVmaW5lUHJvcGVydHkgPT09ICdmdW5jdGlvbicpIHtcbiAgICB0cnkge1xuICAgICAgICBiYWRBcnJheUxpa2UgPSBPYmplY3QuZGVmaW5lUHJvcGVydHkoe30sICdsZW5ndGgnLCB7XG4gICAgICAgICAgICBnZXQ6IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgICAgICB0aHJvdyBpc0NhbGxhYmxlTWFya2VyO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgaXNDYWxsYWJsZU1hcmtlciA9IHt9O1xuICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmUgbm8tdGhyb3ctbGl0ZXJhbFxuICAgICAgICByZWZsZWN0QXBwbHkoZnVuY3Rpb24gKCkgeyB0aHJvdyA0MjsgfSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcbiAgICB9IGNhdGNoIChfKSB7XG4gICAgICAgIGlmIChfICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7XG4gICAgICAgICAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xuICAgICAgICB9XG4gICAgfVxufSBlbHNlIHtcbiAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xufVxuXG52YXIgY29uc3RydWN0b3JSZWdleCA9IC9eXFxzKmNsYXNzXFxiLztcbnZhciBpc0VTNkNsYXNzRm4gPSBmdW5jdGlvbiBpc0VTNkNsYXNzRnVuY3Rpb24odmFsdWU6IGFueSk6IGJvb2xlYW4ge1xuICAgIHRyeSB7XG4gICAgICAgIHZhciBmblN0ciA9IGZuVG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgIHJldHVybiBjb25zdHJ1Y3RvclJlZ2V4LnRlc3QoZm5TdHIpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlOyAvLyBub3QgYSBmdW5jdGlvblxuICAgIH1cbn07XG5cbnZhciB0cnlGdW5jdGlvbk9iamVjdCA9IGZ1bmN0aW9uIHRyeUZ1bmN0aW9uVG9TdHIodmFsdWU6IGFueSk6IGJvb2xlYW4ge1xuICAgIHRyeSB7XG4gICAgICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgICAgICBmblRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG59O1xudmFyIHRvU3RyID0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZztcbnZhciBvYmplY3RDbGFzcyA9ICdbb2JqZWN0IE9iamVjdF0nO1xudmFyIGZuQ2xhc3MgPSAnW29iamVjdCBGdW5jdGlvbl0nO1xudmFyIGdlbkNsYXNzID0gJ1tvYmplY3QgR2VuZXJhdG9yRnVuY3Rpb25dJztcbnZhciBkZGFDbGFzcyA9ICdbb2JqZWN0IEhUTUxBbGxDb2xsZWN0aW9uXSc7IC8vIElFIDExXG52YXIgZGRhQ2xhc3MyID0gJ1tvYmplY3QgSFRNTCBkb2N1bWVudC5hbGwgY2xhc3NdJztcbnZhciBkZGFDbGFzczMgPSAnW29iamVjdCBIVE1MQ29sbGVjdGlvbl0nOyAvLyBJRSA5LTEwXG52YXIgaGFzVG9TdHJpbmdUYWcgPSB0eXBlb2YgU3ltYm9sID09PSAnZnVuY3Rpb24nICYmICEhU3ltYm9sLnRvU3RyaW5nVGFnOyAvLyBiZXR0ZXI6IHVzZSBgaGFzLXRvc3RyaW5ndGFnYFxuXG52YXIgaXNJRTY4ID0gISgwIGluIFssXSk7IC8vIGVzbGludC1kaXNhYmxlLWxpbmUgbm8tc3BhcnNlLWFycmF5cywgY29tbWEtc3BhY2luZ1xuXG52YXIgaXNEREE6ICh2YWx1ZTogYW55KSA9PiBib29sZWFuID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCgpIHsgcmV0dXJuIGZhbHNlOyB9O1xuaWYgKHR5cGVvZiBkb2N1bWVudCA9PT0gJ29iamVjdCcpIHtcbiAgICAvLyBGaXJlZm94IDMgY2Fub25pY2FsaXplcyBEREEgdG8gdW5kZWZpbmVkIHdoZW4gaXQncyBub3QgYWNjZXNzZWQgZGlyZWN0bHlcbiAgICB2YXIgYWxsID0gZG9jdW1lbnQuYWxsO1xuICAgIGlmICh0b1N0ci5jYWxsKGFsbCkgPT09IHRvU3RyLmNhbGwoZG9jdW1lbnQuYWxsKSkge1xuICAgICAgICBpc0REQSA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwodmFsdWUpIHtcbiAgICAgICAgICAgIC8qIGdsb2JhbHMgZG9jdW1lbnQ6IGZhbHNlICovXG4gICAgICAgICAgICAvLyBpbiBJRSA2LTgsIHR5cGVvZiBkb2N1bWVudC5hbGwgaXMgXCJvYmplY3RcIiBhbmQgaXQncyB0cnV0aHlcbiAgICAgICAgICAgIGlmICgoaXNJRTY4IHx8ICF2YWx1ZSkgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3VuZGVmaW5lZCcgfHwgdHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JykpIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICB2YXIgc3RyID0gdG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybiAoXG4gICAgICAgICAgICAgICAgICAgICAgICBzdHIgPT09IGRkYUNsYXNzXG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IGRkYUNsYXNzMlxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczMgLy8gb3BlcmEgMTIuMTZcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gb2JqZWN0Q2xhc3MgLy8gSUUgNi04XG4gICAgICAgICAgICAgICAgICAgICkgJiYgdmFsdWUoJycpID09IG51bGw7IC8vIGVzbGludC1kaXNhYmxlLWxpbmUgZXFlcWVxXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZSkgeyAvKiovIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGlzQ2FsbGFibGVSZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSAge1xuICAgIGlmIChpc0REQSh2YWx1ZSkpIHsgcmV0dXJuIHRydWU7IH1cbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB0cnkge1xuICAgICAgICAocmVmbGVjdEFwcGx5IGFzIGFueSkodmFsdWUsIG51bGwsIGJhZEFycmF5TGlrZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBpZiAoZSAhPT0gaXNDYWxsYWJsZU1hcmtlcikgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB9XG4gICAgcmV0dXJuICFpc0VTNkNsYXNzRm4odmFsdWUpICYmIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcbn1cblxuZnVuY3Rpb24gaXNDYWxsYWJsZU5vUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkge1xuICAgIGlmIChpc0REQSh2YWx1ZSkpIHsgcmV0dXJuIHRydWU7IH1cbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAoaGFzVG9TdHJpbmdUYWcpIHsgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTsgfVxuICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHZhciBzdHJDbGFzcyA9IHRvU3RyLmNhbGwodmFsdWUpO1xuICAgIGlmIChzdHJDbGFzcyAhPT0gZm5DbGFzcyAmJiBzdHJDbGFzcyAhPT0gZ2VuQ2xhc3MgJiYgISgvXlxcW29iamVjdCBIVE1MLykudGVzdChzdHJDbGFzcykpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcbn07XG5cbmV4cG9ydCBkZWZhdWx0IHJlZmxlY3RBcHBseSA/IGlzQ2FsbGFibGVSZWZBcHBseSA6IGlzQ2FsbGFibGVOb1JlZkFwcGx5O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgaXNDYWxsYWJsZSBmcm9tIFwiLi9jYWxsYWJsZS5qc1wiO1xuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdXNlZCBhcyByZWplY3Rpb24gcmVhc29uXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gaXMgY2FuY2VsbGVkIHN1Y2Nlc3NmdWxseS5cbiAqXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbEVycm9yXCJgLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSBjYXVzZSBwYXNzZWQgdG8gdGhlIGNhbmNlbCBtZXRob2QsIGlmIGFueS5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbEVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbEVycm9yYCBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsRXJyb3JcIjtcbiAgICB9XG59XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gcmVqZWN0cyBhZnRlciBiZWluZyBjYW5jZWxsZWQsXG4gKiBvciB3aGVuIHRoZSBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHRocm93cyBvciByZWplY3RzLlxuICpcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcImAuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIHJlYXNvbiB0aGUgcHJvbWlzZSByZWplY3RlZCB3aXRoLlxuICpcbiAqIEJlY2F1c2UgdGhlIG9yaWdpbmFsIHByb21pc2Ugd2FzIGNhbmNlbGxlZCxcbiAqIGEgd3JhcHBlciBwcm9taXNlIHdpbGwgYmUgcGFzc2VkIHRvIHRoZSB1bmhhbmRsZWQgcmVqZWN0aW9uIGxpc3RlbmVyIGluc3RlYWQuXG4gKiBUaGUge0BsaW5rIHByb21pc2V9IHByb3BlcnR5IGhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBvcmlnaW5hbCBwcm9taXNlLlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogSG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIHByb21pc2UgdGhhdCB3YXMgY2FuY2VsbGVkIGFuZCB0aGVuIHJlamVjdGVkLlxuICAgICAqL1xuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPjtcblxuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbGxlZFJlamVjdGlvbkVycm9yYCBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gcHJvbWlzZSAtIFRoZSBwcm9taXNlIHRoYXQgY2F1c2VkIHRoZSBlcnJvciBvcmlnaW5hbGx5LlxuICAgICAqIEBwYXJhbSByZWFzb24gLSBUaGUgcmVqZWN0aW9uIHJlYXNvbi5cbiAgICAgKiBAcGFyYW0gaW5mbyAtIEFuIG9wdGlvbmFsIGluZm9ybWF0aXZlIG1lc3NhZ2Ugc3BlY2lmeWluZyB0aGUgY2lyY3Vtc3RhbmNlcyBpbiB3aGljaCB0aGUgZXJyb3Igd2FzIHRocm93bi5cbiAgICAgKiAgICAgICAgICAgICAgIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgYFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cImAuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCByZWFzb24/OiBhbnksIGluZm8/OiBzdHJpbmcpIHtcbiAgICAgICAgc3VwZXIoKGluZm8gPz8gXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiKSArIFwiIFJlYXNvbjogXCIgKyBlcnJvck1lc3NhZ2UocmVhc29uKSwgeyBjYXVzZTogcmVhc29uIH0pO1xuICAgICAgICB0aGlzLnByb21pc2UgPSBwcm9taXNlO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCI7XG4gICAgfVxufVxuXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+ID0gKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciA9IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciA9IChjYXVzZT86IGFueSkgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiA9IChyZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiwgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcikgPT4gdm9pZDtcblxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwpOiBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxIHwgVFJlc3VsdDI+O1xuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbn1cblxuLyoqXG4gKiBXcmFwcyBhIGNhbmNlbGxhYmxlIHByb21pc2UgYWxvbmcgd2l0aCBpdHMgcmVzb2x1dGlvbiBtZXRob2RzLlxuICogVGhlIGBvbmNhbmNlbGxlZGAgZmllbGQgd2lsbCBiZSBudWxsIGluaXRpYWxseSBidXQgbWF5IGJlIHNldCB0byBwcm92aWRlIGEgY3VzdG9tIGNhbmNlbGxhdGlvbiBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XG4gICAgcmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD47XG4gICAgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcjtcbiAgICBvbmNhbmNlbGxlZDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyIHwgbnVsbDtcbn1cblxuaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlIHtcbiAgICByZWFkb25seSByb290OiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZTtcbiAgICByZXNvbHZpbmc6IGJvb2xlYW47XG4gICAgc2V0dGxlZDogYm9vbGVhbjtcbiAgICByZWFzb24/OiBDYW5jZWxFcnJvcjtcbn1cblxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cbmNvbnN0IGJhcnJpZXJTeW0gPSBTeW1ib2woXCJiYXJyaWVyXCIpO1xuY29uc3QgY2FuY2VsSW1wbFN5bSA9IFN5bWJvbChcImNhbmNlbEltcGxcIik7XG5jb25zdCBzcGVjaWVzOiB0eXBlb2YgU3ltYm9sLnNwZWNpZXMgPSBTeW1ib2wuc3BlY2llcyA/PyBTeW1ib2woXCJzcGVjaWVzUG9seWZpbGxcIik7XG5cbi8qKlxuICogQSBwcm9taXNlIHdpdGggYW4gYXR0YWNoZWQgbWV0aG9kIGZvciBjYW5jZWxsaW5nIGxvbmctcnVubmluZyBvcGVyYXRpb25zIChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWx9KS5cbiAqIENhbmNlbGxhdGlvbiBjYW4gb3B0aW9uYWxseSBiZSBib3VuZCB0byBhbiB7QGxpbmsgQWJvcnRTaWduYWx9XG4gKiBmb3IgYmV0dGVyIGNvbXBvc2FiaWxpdHkgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbE9ufSkuXG4gKlxuICogQ2FuY2VsbGluZyBhIHBlbmRpbmcgcHJvbWlzZSB3aWxsIHJlc3VsdCBpbiBhbiBpbW1lZGlhdGUgcmVqZWN0aW9uXG4gKiB3aXRoIGFuIGluc3RhbmNlIG9mIHtAbGluayBDYW5jZWxFcnJvcn0gYXMgcmVhc29uLFxuICogYnV0IHdob2V2ZXIgc3RhcnRlZCB0aGUgcHJvbWlzZSB3aWxsIGJlIHJlc3BvbnNpYmxlXG4gKiBmb3IgYWN0dWFsbHkgYWJvcnRpbmcgdGhlIHVuZGVybHlpbmcgb3BlcmF0aW9uLlxuICogVG8gdGhpcyBwdXJwb3NlLCB0aGUgY29uc3RydWN0b3IgYW5kIGFsbCBjaGFpbmluZyBtZXRob2RzXG4gKiBhY2NlcHQgb3B0aW9uYWwgY2FuY2VsbGF0aW9uIGNhbGxiYWNrcy5cbiAqXG4gKiBJZiBhIGBDYW5jZWxsYWJsZVByb21pc2VgIHN0aWxsIHJlc29sdmVzIGFmdGVyIGhhdmluZyBiZWVuIGNhbmNlbGxlZCxcbiAqIHRoZSByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWQuIElmIGl0IHJlamVjdHMsIHRoZSByZWFzb25cbiAqIHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbixcbiAqIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGluc3RhbmNlLlxuICogVG8gZmFjaWxpdGF0ZSB0aGUgaGFuZGxpbmcgb2YgY2FuY2VsbGF0aW9uIHJlcXVlc3RzLFxuICogY2FuY2VsbGVkIGBDYW5jZWxsYWJsZVByb21pc2VgcyB3aWxsIF9ub3RfIHJlcG9ydCB1bmhhbmRsZWQgYENhbmNlbEVycm9yYHNcbiAqIHdob3NlIGBjYXVzZWAgZmllbGQgaXMgdGhlIHNhbWUgYXMgdGhlIG9uZSB3aXRoIHdoaWNoIHRoZSBjdXJyZW50IHByb21pc2Ugd2FzIGNhbmNlbGxlZC5cbiAqXG4gKiBBbGwgdXN1YWwgcHJvbWlzZSBtZXRob2RzIGFyZSBkZWZpbmVkIGFuZCByZXR1cm4gYSBgQ2FuY2VsbGFibGVQcm9taXNlYFxuICogd2hvc2UgY2FuY2VsIG1ldGhvZCB3aWxsIGNhbmNlbCB0aGUgcGFyZW50IG9wZXJhdGlvbiBhcyB3ZWxsLCBwcm9wYWdhdGluZyB0aGUgY2FuY2VsbGF0aW9uIHJlYXNvblxuICogdXB3YXJkcyB0aHJvdWdoIHByb21pc2UgY2hhaW5zLlxuICogQ29udmVyc2VseSwgY2FuY2VsbGluZyBhIHByb21pc2Ugd2lsbCBub3QgYXV0b21hdGljYWxseSBjYW5jZWwgZGVwZW5kZW50IHByb21pc2VzIGRvd25zdHJlYW06XG4gKiBgYGB0c1xuICogbGV0IHJvb3QgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMSA9IHJvb3QudGhlbigoKSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDIgPSBjaGlsZDEudGhlbigoKSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDMgPSByb290LmNhdGNoKCgpID0+IHsgLi4uIH0pO1xuICogY2hpbGQxLmNhbmNlbCgpOyAvLyBDYW5jZWxzIGNoaWxkMSBhbmQgcm9vdCwgYnV0IG5vdCBjaGlsZDIgb3IgY2hpbGQzXG4gKiBgYGBcbiAqIENhbmNlbGxpbmcgYSBwcm9taXNlIHRoYXQgaGFzIGFscmVhZHkgc2V0dGxlZCBpcyBzYWZlIGFuZCBoYXMgbm8gY29uc2VxdWVuY2UuXG4gKlxuICogVGhlIGBjYW5jZWxgIG1ldGhvZCByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IF9hbHdheXMgZnVsZmlsbHNfXG4gKiBhZnRlciB0aGUgd2hvbGUgY2hhaW4gaGFzIHByb2Nlc3NlZCB0aGUgY2FuY2VsIHJlcXVlc3RcbiAqIGFuZCBhbGwgYXR0YWNoZWQgY2FsbGJhY2tzIHVwIHRvIHRoYXQgbW9tZW50IGhhdmUgcnVuLlxuICpcbiAqIEFsbCBFUzIwMjQgcHJvbWlzZSBtZXRob2RzIChzdGF0aWMgYW5kIGluc3RhbmNlKSBhcmUgZGVmaW5lZCBvbiBDYW5jZWxsYWJsZVByb21pc2UsXG4gKiBidXQgYWN0dWFsIGF2YWlsYWJpbGl0eSBtYXkgdmFyeSB3aXRoIE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAqXG4gKiBJbiBsaW5lIHdpdGggdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmcsXG4gKiBgQ2FuY2VsbGFibGVQcm9taXNlYCBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLlxuICogRXh0ZW5kZXJzIHNob3VsZCB0YWtlIGNhcmUgdG8gcHJvdmlkZSB0aGVpciBvd24gbWV0aG9kIGltcGxlbWVudGF0aW9ucy5cbiAqIFRoaXMgbWlnaHQgYmUgcmVjb25zaWRlcmVkIGluIGNhc2UgdGhlIHByb3Bvc2FsIGlzIHJldGlyZWQuXG4gKlxuICogQ2FuY2VsbGFibGVQcm9taXNlIGlzIGEgd3JhcHBlciBhcm91bmQgdGhlIERPTSBQcm9taXNlIG9iamVjdFxuICogYW5kIGlzIGNvbXBsaWFudCB3aXRoIHRoZSBbUHJvbWlzZXMvQSsgc3BlY2lmaWNhdGlvbl0oaHR0cHM6Ly9wcm9taXNlc2FwbHVzLmNvbS8pXG4gKiAoaXQgcGFzc2VzIHRoZSBbY29tcGxpYW5jZSBzdWl0ZV0oaHR0cHM6Ly9naXRodWIuY29tL3Byb21pc2VzLWFwbHVzL3Byb21pc2VzLXRlc3RzKSlcbiAqIGlmIHNvIGlzIHRoZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uLlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGFibGVQcm9taXNlPFQ+IGV4dGVuZHMgUHJvbWlzZTxUPiBpbXBsZW1lbnRzIFByb21pc2VMaWtlPFQ+LCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICAvKiogQGludGVybmFsICovXG4gICAgcHJpdmF0ZSBbYmFycmllclN5bV0hOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiB8IG51bGw7XG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIHByaXZhdGUgcmVhZG9ubHkgW2NhbmNlbEltcGxTeW1dITogKHJlYXNvbjogQ2FuY2VsRXJyb3IpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgYENhbmNlbGxhYmxlUHJvbWlzZWAuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZXhlY3V0b3IgLSBBIGNhbGxiYWNrIHVzZWQgdG8gaW5pdGlhbGl6ZSB0aGUgcHJvbWlzZS4gVGhpcyBjYWxsYmFjayBpcyBwYXNzZWQgdHdvIGFyZ3VtZW50czpcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhIGByZXNvbHZlYCBjYWxsYmFjayB1c2VkIHRvIHJlc29sdmUgdGhlIHByb21pc2Ugd2l0aCBhIHZhbHVlXG4gICAgICogICAgICAgICAgICAgICAgICAgb3IgdGhlIHJlc3VsdCBvZiBhbm90aGVyIHByb21pc2UgKHBvc3NpYmx5IGNhbmNlbGxhYmxlKSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYSBgcmVqZWN0YCBjYWxsYmFjayB1c2VkIHRvIHJlamVjdCB0aGUgcHJvbWlzZSB3aXRoIGEgcHJvdmlkZWQgcmVhc29uIG9yIGVycm9yLlxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIHRoZSB2YWx1ZSBwcm92aWRlZCB0byB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogICAgICAgICAgICAgICAgICAgKGl0IGhhcyBhIGB0aGVuYCBfYW5kXyBhIGBjYW5jZWxgIG1ldGhvZCksXG4gICAgICogICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZm9yd2FyZGVkIHRvIHRoYXQgb2JqZWN0IGFuZCB0aGUgb25jYW5jZWxsZWQgd2lsbCBub3QgYmUgaW52b2tlZCBhbnltb3JlLlxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIGFueSBvbmUgb2YgdGhlIHR3byBjYWxsYmFja3MgaXMgY2FsbGVkIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkLFxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoZSBwcm92aWRlZCB2YWx1ZXMgd2lsbCBiZSBjYW5jZWxsZWQgYW5kIHJlc29sdmVkIGFzIHVzdWFsLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGJ1dCB0aGVpciByZXN1bHRzIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqICAgICAgICAgICAgICAgICAgIEhvd2V2ZXIsIGlmIHRoZSByZXNvbHV0aW9uIHByb2Nlc3MgdWx0aW1hdGVseSBlbmRzIHVwIGluIGEgcmVqZWN0aW9uXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhhdCBpcyBub3QgZHVlIHRvIGNhbmNlbGxhdGlvbiwgdGhlIHJlamVjdGlvbiByZWFzb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICB3aWxsIGJlIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9XG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKiBAcGFyYW0gb25jYW5jZWxsZWQgLSBJdCBpcyB0aGUgY2FsbGVyJ3MgcmVzcG9uc2liaWxpdHkgdG8gZW5zdXJlIHRoYXQgYW55IG9wZXJhdGlvblxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHN0YXJ0ZWQgYnkgdGhlIGV4ZWN1dG9yIGlzIHByb3Blcmx5IGhhbHRlZCB1cG9uIGNhbmNlbGxhdGlvbi5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBUaGlzIG9wdGlvbmFsIGNhbGxiYWNrIGNhbiBiZSB1c2VkIHRvIHRoYXQgcHVycG9zZS5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJdCB3aWxsIGJlIGNhbGxlZCBfc3luY2hyb25vdXNseV8gd2l0aCBhIGNhbmNlbGxhdGlvbiBjYXVzZVxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdoZW4gY2FuY2VsbGF0aW9uIGlzIHJlcXVlc3RlZCwgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSByZWplY3RlZFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9LCBidXQgX2JlZm9yZV9cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbnkge0BsaW5rIHRoZW59L3tAbGluayBjYXRjaH0ve0BsaW5rIGZpbmFsbHl9IGNhbGxiYWNrIHJ1bnMuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfVxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpbGwgb25seSBmdWxmaWxsIGFmdGVyIHRoZSBmb3JtZXIgaGFzIHNldHRsZWQuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVW5oYW5kbGVkIGV4Y2VwdGlvbnMgb3IgcmVqZWN0aW9ucyBmcm9tIHRoZSBjYWxsYmFjayB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgdW5oYW5kbGVkIHJlamVjdGlvbnMuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBjYWxsZWQgYmVmb3JlIGNhbmNlbGxhdGlvbiB3aXRoIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhpcyBwcm9taXNlIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gdGhhdCBwcm9taXNlLFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFuZCB0aGUgb3JpZ2luYWwgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihleGVjdXRvcjogQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4sIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKSB7XG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xuICAgICAgICBzdXBlcigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcblxuICAgICAgICBpZiAoKHRoaXMuY29uc3RydWN0b3IgYXMgYW55KVtzcGVjaWVzXSAhPT0gUHJvbWlzZSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZSBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLiBQbGVhc2UgcmVmcmFpbiBmcm9tIG92ZXJyaWRpbmcgdGhlIFtTeW1ib2wuc3BlY2llc10gc3RhdGljIHByb3BlcnR5LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0ge1xuICAgICAgICAgICAgcHJvbWlzZTogdGhpcyxcbiAgICAgICAgICAgIHJlc29sdmUsXG4gICAgICAgICAgICByZWplY3QsXG4gICAgICAgICAgICBnZXQgb25jYW5jZWxsZWQoKSB7IHJldHVybiBvbmNhbmNlbGxlZCA/PyBudWxsOyB9LFxuICAgICAgICAgICAgc2V0IG9uY2FuY2VsbGVkKGNiKSB7IG9uY2FuY2VsbGVkID0gY2IgPz8gdW5kZWZpbmVkOyB9XG4gICAgICAgIH07XG5cbiAgICAgICAgY29uc3Qgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xuICAgICAgICAgICAgZ2V0IHJvb3QoKSB7IHJldHVybiBzdGF0ZTsgfSxcbiAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXG4gICAgICAgICAgICBzZXR0bGVkOiBmYWxzZVxuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFNldHVwIGNhbmNlbGxhdGlvbiBzeXN0ZW0uXG4gICAgICAgIHZvaWQgT2JqZWN0LmRlZmluZVByb3BlcnRpZXModGhpcywge1xuICAgICAgICAgICAgW2JhcnJpZXJTeW1dOiB7XG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogdHJ1ZSxcbiAgICAgICAgICAgICAgICB2YWx1ZTogbnVsbFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIFtjYW5jZWxJbXBsU3ltXToge1xuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHZhbHVlOiBjYW5jZWxsZXJGb3IocHJvbWlzZSwgc3RhdGUpXG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuXG4gICAgICAgIC8vIFJ1biB0aGUgYWN0dWFsIGV4ZWN1dG9yLlxuICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIHN0YXRlKTtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGV4ZWN1dG9yKHJlc29sdmVyRm9yKHByb21pc2UsIHN0YXRlKSwgcmVqZWN0b3IpO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHtcbiAgICAgICAgICAgICAgICBjb25zb2xlLmxvZyhcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gQ2FuY2VsbGFibGVQcm9taXNlIGV4ZWN1dG9yLlwiLCBlcnIpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2FuY2VscyBpbW1lZGlhdGVseSB0aGUgZXhlY3V0aW9uIG9mIHRoZSBvcGVyYXRpb24gYXNzb2NpYXRlZCB3aXRoIHRoaXMgcHJvbWlzZS5cbiAgICAgKiBUaGUgcHJvbWlzZSByZWplY3RzIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9IGluc3RhbmNlIGFzIHJlYXNvbixcbiAgICAgKiB3aXRoIHRoZSB7QGxpbmsgQ2FuY2VsRXJyb3IjY2F1c2V9IHByb3BlcnR5IHNldCB0byB0aGUgZ2l2ZW4gYXJndW1lbnQsIGlmIGFueS5cbiAgICAgKlxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIGFmdGVyIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQ7XG4gICAgICogcmVwZWF0ZWQgY2FsbHMgaW4gcGFydGljdWxhciBhcmUgc2FmZSwgYnV0IG9ubHkgdGhlIGZpcnN0IG9uZVxuICAgICAqIHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXG4gICAgICpcbiAgICAgKiBUaGUgYENhbmNlbEVycm9yYCBleGNlcHRpb24gX25lZWQgbm90XyBiZSBoYW5kbGVkIGV4cGxpY2l0bHkgX29uIHRoZSBwcm9taXNlcyB0aGF0IGFyZSBiZWluZyBjYW5jZWxsZWQ6X1xuICAgICAqIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpdGggbm8gYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXIgZG9lcyBub3QgdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LlxuICAgICAqIFRoZXJlZm9yZSwgdGhlIGZvbGxvd2luZyBpZGlvbXMgYXJlIGFsbCBlcXVhbGx5IGNvcnJlY3Q6XG4gICAgICogYGBgdHNcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLmNhbmNlbCgpO1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhbmNlbCgpO1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhdGNoKC4uLikuY2FuY2VsKCk7XG4gICAgICogYGBgXG4gICAgICogV2hlbmV2ZXIgc29tZSBjYW5jZWxsZWQgcHJvbWlzZSBpbiBhIGNoYWluIHJlamVjdHMgd2l0aCBhIGBDYW5jZWxFcnJvcmBcbiAgICAgKiB3aXRoIHRoZSBzYW1lIGNhbmNlbGxhdGlvbiBjYXVzZSBhcyBpdHNlbGYsIHRoZSBlcnJvciB3aWxsIGJlIGRpc2NhcmRlZCBzaWxlbnRseS5cbiAgICAgKiBIb3dldmVyLCB0aGUgYENhbmNlbEVycm9yYCBfd2lsbCBzdGlsbCBiZSBkZWxpdmVyZWRfIHRvIGFsbCBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlcnNcbiAgICAgKiBhZGRlZCBieSB7QGxpbmsgdGhlbn0gYW5kIHJlbGF0ZWQgbWV0aG9kczpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gICAgICogY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS5jYXRjaChjb25zb2xlLmxvZyk7XG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIEEgQ2FuY2VsRXJyb3IgaXMgcHJpbnRlZCB0byB0aGUgY29uc29sZS5cbiAgICAgKiBgYGBcbiAgICAgKiBJZiB0aGUgYENhbmNlbEVycm9yYCBpcyBub3QgaGFuZGxlZCBkb3duc3RyZWFtIGJ5IHRoZSB0aW1lIGl0IHJlYWNoZXNcbiAgICAgKiBhIF9ub24tY2FuY2VsbGVkXyBwcm9taXNlLCBpdCBfd2lsbF8gdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LFxuICAgICAqIGp1c3QgbGlrZSBub3JtYWwgcmVqZWN0aW9ucyB3b3VsZDpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gICAgICogbGV0IGNoYWluZWQgPSBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLnRoZW4oKCkgPT4geyAuLi4gfSk7IC8vIE5vIGNhdGNoLi4uXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIFVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQgb24gY2hhaW5lZCFcbiAgICAgKiBgYGBcbiAgICAgKiBUaGVyZWZvcmUsIGl0IGlzIGltcG9ydGFudCB0byBlaXRoZXIgY2FuY2VsIHdob2xlIHByb21pc2UgY2hhaW5zIGZyb20gdGhlaXIgdGFpbCxcbiAgICAgKiBhcyBzaG93biBpbiB0aGUgY29ycmVjdCBpZGlvbXMgYWJvdmUsIG9yIHRha2UgY2FyZSBvZiBoYW5kbGluZyBlcnJvcnMgZXZlcnl3aGVyZS5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgY2FuY2VsbGFibGUgcHJvbWlzZSB0aGF0IF9mdWxmaWxsc18gYWZ0ZXIgdGhlIGNhbmNlbCBjYWxsYmFjayAoaWYgYW55KVxuICAgICAqIGFuZCBhbGwgaGFuZGxlcnMgYXR0YWNoZWQgdXAgdG8gdGhlIGNhbGwgdG8gY2FuY2VsIGhhdmUgcnVuLlxuICAgICAqIElmIHRoZSBjYW5jZWwgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBieSBgY2FuY2VsYFxuICAgICAqIHdpbGwgYWxzbyB3YWl0IGZvciB0aGF0IHRoZW5hYmxlIHRvIHNldHRsZS5cbiAgICAgKiBUaGlzIGVuYWJsZXMgY2FsbGVycyB0byB3YWl0IGZvciB0aGUgY2FuY2VsbGVkIG9wZXJhdGlvbiB0byB0ZXJtaW5hdGVcbiAgICAgKiB3aXRob3V0IGJlaW5nIGZvcmNlZCB0byBoYW5kbGUgcG90ZW50aWFsIGVycm9ycyBhdCB0aGUgY2FsbCBzaXRlLlxuICAgICAqIGBgYHRzXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCkudGhlbigoKSA9PiB7XG4gICAgICogICAgIC8vIENsZWFudXAgZmluaXNoZWQsIGl0J3Mgc2FmZSB0byBkbyBzb21ldGhpbmcgZWxzZS5cbiAgICAgKiB9LCAoZXJyKSA9PiB7XG4gICAgICogICAgIC8vIFVucmVhY2hhYmxlOiB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIGNhbmNlbCB3aWxsIG5ldmVyIHJlamVjdC5cbiAgICAgKiB9KTtcbiAgICAgKiBgYGBcbiAgICAgKiBOb3RlIHRoYXQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBfbm90XyBoYW5kbGUgaW1wbGljaXRseSBhbnkgcmVqZWN0aW9uXG4gICAgICogdGhhdCBtaWdodCBoYXZlIG9jY3VycmVkIGFscmVhZHkgaW4gdGhlIGNhbmNlbGxlZCBjaGFpbi5cbiAgICAgKiBJdCB3aWxsIGp1c3QgdHJhY2sgd2hldGhlciByZWdpc3RlcmVkIGhhbmRsZXJzIGhhdmUgYmVlbiBleGVjdXRlZCBvciBub3QuXG4gICAgICogVGhlcmVmb3JlLCB1bmhhbmRsZWQgcmVqZWN0aW9ucyB3aWxsIG5ldmVyIGJlIHNpbGVudGx5IGhhbmRsZWQgYnkgY2FsbGluZyBjYW5jZWwuXG4gICAgICovXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIC8vIElOVkFSSUFOVDogdGhlIHJlc3VsdCBvZiB0aGlzW2NhbmNlbEltcGxTeW1dIGFuZCB0aGUgYmFycmllciBkbyBub3QgZXZlciByZWplY3QuXG4gICAgICAgICAgICAvLyBVbmZvcnR1bmF0ZWx5IG1hY09TIEhpZ2ggU2llcnJhIGRvZXMgbm90IHN1cHBvcnQgUHJvbWlzZS5hbGxTZXR0bGVkLlxuICAgICAgICAgICAgUHJvbWlzZS5hbGwoW1xuICAgICAgICAgICAgICAgIHRoaXNbY2FuY2VsSW1wbFN5bV0obmV3IENhbmNlbEVycm9yKFwiUHJvbWlzZSBjYW5jZWxsZWQuXCIsIHsgY2F1c2UgfSkpLFxuICAgICAgICAgICAgICAgIGN1cnJlbnRCYXJyaWVyKHRoaXMpXG4gICAgICAgICAgICBdKS50aGVuKCgpID0+IHJlc29sdmUoKSwgKCkgPT4gcmVzb2x2ZSgpKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQmluZHMgcHJvbWlzZSBjYW5jZWxsYXRpb24gdG8gdGhlIGFib3J0IGV2ZW50IG9mIHRoZSBnaXZlbiB7QGxpbmsgQWJvcnRTaWduYWx9LlxuICAgICAqIElmIHRoZSBzaWduYWwgaGFzIGFscmVhZHkgYWJvcnRlZCwgdGhlIHByb21pc2Ugd2lsbCBiZSBjYW5jZWxsZWQgaW1tZWRpYXRlbHkuXG4gICAgICogV2hlbiBlaXRoZXIgY29uZGl0aW9uIGlzIHZlcmlmaWVkLCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlIHdpbGwgYmUgc2V0XG4gICAgICogdG8gdGhlIHNpZ25hbCdzIGFib3J0IHJlYXNvbiAoc2VlIHtAbGluayBBYm9ydFNpZ25hbCNyZWFzb259KS5cbiAgICAgKlxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIChvciBpZiB0aGUgc2lnbmFsIGFib3J0cykgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkLlxuICAgICAqIE9ubHkgdGhlIGZpcnN0IHNpZ25hbCB0byBhYm9ydCB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxuICAgICAqXG4gICAgICogRm9yIG1vcmUgZGV0YWlscyBhYm91dCB0aGUgY2FuY2VsbGF0aW9uIHByb2Nlc3MsXG4gICAgICogc2VlIHtAbGluayBjYW5jZWx9IGFuZCB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAgY29uc3RydWN0b3IuXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBlbmFibGVzIGBhd2FpdGBpbmcgY2FuY2VsbGFibGUgcHJvbWlzZXMgd2l0aG91dCBoYXZpbmdcbiAgICAgKiB0byBzdG9yZSB0aGVtIGZvciBmdXR1cmUgY2FuY2VsbGF0aW9uLCBlLmcuOlxuICAgICAqIGBgYHRzXG4gICAgICogYXdhaXQgbG9uZ1J1bm5pbmdPcGVyYXRpb24oKS5jYW5jZWxPbihzaWduYWwpO1xuICAgICAqIGBgYFxuICAgICAqIGluc3RlYWQgb2Y6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQgPSBsb25nUnVubmluZ09wZXJhdGlvbigpO1xuICAgICAqIGF3YWl0IHByb21pc2VUb0JlQ2FuY2VsbGVkO1xuICAgICAqIGBgYFxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhpcyBwcm9taXNlLCBmb3IgbWV0aG9kIGNoYWluaW5nLlxuICAgICAqL1xuICAgIGNhbmNlbE9uKHNpZ25hbDogQWJvcnRTaWduYWwpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBpZiAoc2lnbmFsLmFib3J0ZWQpIHtcbiAgICAgICAgICAgIHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbilcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHNpZ25hbC5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbiksIHtjYXB0dXJlOiB0cnVlfSk7XG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gdGhpcztcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBjYWxsYmFja3MgZm9yIHRoZSByZXNvbHV0aW9uIGFuZC9vciByZWplY3Rpb24gb2YgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBvbmZ1bGZpbGxlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlc29sdmVkLlxuICAgICAqIEBwYXJhbSBvbnJlamVjdGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVqZWN0ZWQuXG4gICAgICogQHJldHVybnMgQSBgQ2FuY2VsbGFibGVQcm9taXNlYCBmb3IgdGhlIGNvbXBsZXRpb24gb2Ygd2hpY2hldmVyIGNhbGxiYWNrIGlzIGV4ZWN1dGVkLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGhvb2tlZCB1cCB0byBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHVwIHRoZSBjaGFpbiwgYnV0IG5vdCBkb3duOlxuICAgICAqXG4gICAgICogICAtIGlmIHRoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIGJlIGludm9rZWQgd2l0aCBhIGBDYW5jZWxFcnJvcmBcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIF93aWxsIHJlc29sdmUgcmVndWxhcmx5XyB3aXRoIGl0cyByZXN1bHQ7XG4gICAgICogICAtIGNvbnZlcnNlbHksIGlmIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgX3RoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQgdG9vO19cbiAgICAgKiAgICAgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgc3RpbGwgYmUgaW52b2tlZCB3aXRoIHRoZSBwYXJlbnQncyBgQ2FuY2VsRXJyb3JgLFxuICAgICAqICAgICBidXQgaXRzIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZFxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhIGBDYW5jZWxFcnJvcmAgYXMgd2VsbC5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBlaXRoZXIgY2FsbGJhY2sgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKi9cbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+IHtcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLnRoZW4gY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIE5PVEU6IFR5cGVTY3JpcHQncyBidWlsdC1pbiB0eXBlIGZvciB0aGVuIGlzIGJyb2tlbixcbiAgICAgICAgLy8gYXMgaXQgYWxsb3dzIHNwZWNpZnlpbmcgYW4gYXJiaXRyYXJ5IFRSZXN1bHQxICE9IFQgZXZlbiB3aGVuIG9uZnVsZmlsbGVkIGlzIG5vdCBhIGZ1bmN0aW9uLlxuICAgICAgICAvLyBXZSBjYW5ub3QgZml4IGl0IGlmIHdlIHdhbnQgdG8gQ2FuY2VsbGFibGVQcm9taXNlIHRvIGltcGxlbWVudCBQcm9taXNlTGlrZTxUPi5cblxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25mdWxmaWxsZWQpKSB7IG9uZnVsZmlsbGVkID0gaWRlbnRpdHkgYXMgYW55OyB9XG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbnJlamVjdGVkKSkgeyBvbnJlamVjdGVkID0gdGhyb3dlcjsgfVxuXG4gICAgICAgIGlmIChvbmZ1bGZpbGxlZCA9PT0gaWRlbnRpdHkgJiYgb25yZWplY3RlZCA9PSB0aHJvd2VyKSB7XG4gICAgICAgICAgICAvLyBTaG9ydGN1dCBmb3IgdHJpdmlhbCBhcmd1bWVudHMuXG4gICAgICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSkgPT4gcmVzb2x2ZSh0aGlzIGFzIGFueSkpO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgYmFycmllcjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSB7fTtcbiAgICAgICAgdGhpc1tiYXJyaWVyU3ltXSA9IGJhcnJpZXI7XG5cbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgdm9pZCBzdXBlci50aGVuKFxuICAgICAgICAgICAgICAgICh2YWx1ZSkgPT4ge1xuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xuXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9uZnVsZmlsbGVkISh2YWx1ZSkpO1xuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICAocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xuXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9ucmVqZWN0ZWQhKHJlYXNvbikpO1xuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgKTtcbiAgICAgICAgfSwgYXN5bmMgKGNhdXNlPykgPT4ge1xuICAgICAgICAgICAgLy9jYW5jZWxsZWQgPSB0cnVlO1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICByZXR1cm4gb25jYW5jZWxsZWQ/LihjYXVzZSk7XG4gICAgICAgICAgICB9IGZpbmFsbHkge1xuICAgICAgICAgICAgICAgIGF3YWl0IHRoaXMuY2FuY2VsKGNhdXNlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayBmb3Igb25seSB0aGUgcmVqZWN0aW9uIG9mIHRoZSBQcm9taXNlLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIEl0IGlzIGVxdWl2YWxlbnQgdG9cbiAgICAgKiBgYGB0c1xuICAgICAqIGNhbmNlbGxhYmxlUHJvbWlzZS50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xuICAgICAqIGBgYFxuICAgICAqIGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgYG9ucmVqZWN0ZWRgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxuICAgICAqL1xuICAgIGNhdGNoPFRSZXN1bHQgPSBuZXZlcj4ob25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IChQcm9taXNlTGlrZTxUUmVzdWx0PiB8IFRSZXN1bHQpKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQgfCBUUmVzdWx0PiB7XG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayB0aGF0IGlzIGludm9rZWQgd2hlbiB0aGUgQ2FuY2VsbGFibGVQcm9taXNlIGlzIHNldHRsZWQgKGZ1bGZpbGxlZCBvciByZWplY3RlZCkuIFRoZVxuICAgICAqIHJlc29sdmVkIHZhbHVlIGNhbm5vdCBiZSBhY2Nlc3NlZCBvciBtb2RpZmllZCBmcm9tIHRoZSBjYWxsYmFjay5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHNldHRsZSBpbiB0aGUgc2FtZSBzdGF0ZSBhcyB0aGUgb3JpZ2luYWwgb25lXG4gICAgICogYWZ0ZXIgdGhlIHByb3ZpZGVkIGNhbGxiYWNrIGhhcyBjb21wbGV0ZWQgZXhlY3V0aW9uLFxuICAgICAqIHVubGVzcyB0aGUgY2FsbGJhY2sgdGhyb3dzIG9yIHJldHVybnMgYSByZWplY3RpbmcgcHJvbWlzZSxcbiAgICAgKiBpbiB3aGljaCBjYXNlIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IGFzIHdlbGwuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIE9uY2UgdGhlIHBhcmVudCBwcm9taXNlIHNldHRsZXMsIHRoZSBgb25maW5hbGx5YCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGlzIGltcGxlbWVudGVkIGluIHRlcm1zIG9mIHtAbGluayB0aGVufSBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cbiAgICAgKiBJdCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgYG9uZmluYWxseWAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXG4gICAgICovXG4gICAgZmluYWxseShvbmZpbmFsbHk/OiAoKCkgPT4gdm9pZCkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS5maW5hbGx5IGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25maW5hbGx5KSkge1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMudGhlbihvbmZpbmFsbHksIG9uZmluYWxseSwgb25jYW5jZWxsZWQpO1xuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbihcbiAgICAgICAgICAgICh2YWx1ZSkgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4gdmFsdWUpLFxuICAgICAgICAgICAgKHJlYXNvbj8pID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHsgdGhyb3cgcmVhc29uOyB9KSxcbiAgICAgICAgICAgIG9uY2FuY2VsbGVkLFxuICAgICAgICApO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFdlIHVzZSB0aGUgYFtTeW1ib2wuc3BlY2llc11gIHN0YXRpYyBwcm9wZXJ0eSwgaWYgYXZhaWxhYmxlLFxuICAgICAqIHRvIGRpc2FibGUgdGhlIGJ1aWx0LWluIGF1dG9tYXRpYyBzdWJjbGFzc2luZyBmZWF0dXJlcyBmcm9tIHtAbGluayBQcm9taXNlfS5cbiAgICAgKiBJdCBpcyBjcml0aWNhbCBmb3IgcGVyZm9ybWFuY2UgcmVhc29ucyB0aGF0IGV4dGVuZGVycyBkbyBub3Qgb3ZlcnJpZGUgdGhpcy5cbiAgICAgKiBPbmNlIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nXG4gICAgICogaXMgZWl0aGVyIGFjY2VwdGVkIG9yIHJldGlyZWQsIHRoaXMgaW1wbGVtZW50YXRpb24gd2lsbCBoYXZlIHRvIGJlIHJldmlzZWQgYWNjb3JkaW5nbHkuXG4gICAgICpcbiAgICAgKiBAaWdub3JlXG4gICAgICogQGludGVybmFsXG4gICAgICovXG4gICAgc3RhdGljIGdldCBbc3BlY2llc10oKSB7XG4gICAgICAgIHJldHVybiBQcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSwgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgUHJvbWlzZSBpcyByZWplY3RlZC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYWxsPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+W10+O1xuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IEF3YWl0ZWQ8VFtQXT47IH0+O1xuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGwoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUgb3IgcmVqZWN0LlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFQ+PltdPjtcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogUHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUW1BdPj47IH0+O1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsU2V0dGxlZChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRoZSBhbnkgZnVuY3Rpb24gcmV0dXJucyBhIHByb21pc2UgdGhhdCBpcyBmdWxmaWxsZWQgYnkgdGhlIGZpcnN0IGdpdmVuIHByb21pc2UgdG8gYmUgZnVsZmlsbGVkLFxuICAgICAqIG9yIHJlamVjdGVkIHdpdGggYW4gQWdncmVnYXRlRXJyb3IgY29udGFpbmluZyBhbiBhcnJheSBvZiByZWplY3Rpb24gcmVhc29uc1xuICAgICAqIGlmIGFsbCBvZiB0aGUgZ2l2ZW4gcHJvbWlzZXMgYXJlIHJlamVjdGVkLlxuICAgICAqIEl0IHJlc29sdmVzIGFsbCBlbGVtZW50cyBvZiB0aGUgcGFzc2VkIGl0ZXJhYmxlIHRvIHByb21pc2VzIGFzIGl0IHJ1bnMgdGhpcyBhbGdvcml0aG0uXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFueTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbnkoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIG9yIHJlamVjdGVkIHdoZW4gYW55IG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyBhcmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJhY2U8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmFjZShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBjYW5jZWxsZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgY2FuY2VsPFQgPSBuZXZlcj4oY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBjb25zdCBwID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XG4gICAgICAgIHAuY2FuY2VsKGNhdXNlKTtcbiAgICAgICAgcmV0dXJuIHA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBjYW5jZWxzXG4gICAgICogYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCBjYXVzZS5cbiAgICAgKlxuICAgICAqIElmIHRoZSB7QGxpbmsgQWJvcnRTaWduYWwudGltZW91dH0gZmFjdG9yeSBtZXRob2QgaXMgYXZhaWxhYmxlLFxuICAgICAqIGl0IGlzIHVzZWQgdG8gYmFzZSB0aGUgdGltZW91dCBvbiBfYWN0aXZlXyB0aW1lIHJhdGhlciB0aGFuIF9lbGFwc2VkXyB0aW1lLlxuICAgICAqIE90aGVyd2lzZSwgYHRpbWVvdXRgIGZhbGxzIGJhY2sgdG8ge0BsaW5rIHNldFRpbWVvdXR9LlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHRpbWVvdXQ8VCA9IG5ldmVyPihtaWxsaXNlY29uZHM6IG51bWJlciwgY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XG4gICAgICAgIGlmIChBYm9ydFNpZ25hbCAmJiB0eXBlb2YgQWJvcnRTaWduYWwgPT09ICdmdW5jdGlvbicgJiYgQWJvcnRTaWduYWwudGltZW91dCAmJiB0eXBlb2YgQWJvcnRTaWduYWwudGltZW91dCA9PT0gJ2Z1bmN0aW9uJykge1xuICAgICAgICAgICAgQWJvcnRTaWduYWwudGltZW91dChtaWxsaXNlY29uZHMpLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSwgbWlsbGlzZWNvbmRzKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dC5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHNsZWVwKG1pbGxpc2Vjb25kczogbnVtYmVyKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXJcbiAgICAgKiB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgc2xlZXA8VD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xuICAgIHN0YXRpYyBzbGVlcDxUID0gdm9pZD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlPzogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHJlc29sdmUodmFsdWUhKSwgbWlsbGlzZWNvbmRzKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZWplY3RlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCByZWFzb24uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVqZWN0PFQgPSBuZXZlcj4ocmVhc29uPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKF8sIHJlamVjdCkgPT4gcmVqZWN0KHJlYXNvbikpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmUoKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgcmVzb2x2ZTxUID0gdm9pZD4odmFsdWU/OiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj4ge1xuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpIHtcbiAgICAgICAgICAgIC8vIE9wdGltaXNlIGZvciBjYW5jZWxsYWJsZSBwcm9taXNlcy5cbiAgICAgICAgICAgIHJldHVybiB2YWx1ZTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+KChyZXNvbHZlKSA9PiByZXNvbHZlKHZhbHVlKSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgYW5kIHJldHVybnMgaXQgaW4gYW4gb2JqZWN0LCBhbG9uZyB3aXRoIGl0cyByZXNvbHZlIGFuZCByZWplY3QgZnVuY3Rpb25zXG4gICAgICogYW5kIGEgZ2V0dGVyL3NldHRlciBmb3IgdGhlIGNhbmNlbGxhdGlvbiBjYWxsYmFjay5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgd2l0aFJlc29sdmVyczxUPigpOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICAgICAgbGV0IHJlc3VsdDogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHsgb25jYW5jZWxsZWQ6IG51bGwgfSBhcyBhbnk7XG4gICAgICAgIHJlc3VsdC5wcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICByZXN1bHQucmVzb2x2ZSA9IHJlc29sdmU7XG4gICAgICAgICAgICByZXN1bHQucmVqZWN0ID0gcmVqZWN0O1xuICAgICAgICB9LCAoY2F1c2U/OiBhbnkpID0+IHsgcmVzdWx0Lm9uY2FuY2VsbGVkPy4oY2F1c2UpOyB9KTtcbiAgICAgICAgcmV0dXJuIHJlc3VsdDtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgY2FuY2VsbGF0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSByZXN1bHRpbmcgZnVuY3Rpb24gZG9lcyBub3QgcmVqZWN0LlxuICovXG5mdW5jdGlvbiBjYW5jZWxsZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKSB7XG4gICAgbGV0IGNhbmNlbGxhdGlvblByb21pc2U6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9IHVuZGVmaW5lZDtcblxuICAgIHJldHVybiAocmVhc29uOiBDYW5jZWxFcnJvcik6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9PiB7XG4gICAgICAgIGlmICghc3RhdGUuc2V0dGxlZCkge1xuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBzdGF0ZS5yZWFzb24gPSByZWFzb247XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xuXG4gICAgICAgICAgICAvLyBBdHRhY2ggYW4gZXJyb3IgaGFuZGxlciB0aGF0IGlnbm9yZXMgdGhpcyBzcGVjaWZpYyByZWplY3Rpb24gcmVhc29uIGFuZCBub3RoaW5nIGVsc2UuXG4gICAgICAgICAgICAvLyBJbiB0aGVvcnksIGEgc2FuZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uIGF0IHRoaXMgcG9pbnRcbiAgICAgICAgICAgIC8vIHNob3VsZCBhbHdheXMgcmVqZWN0IHdpdGggb3VyIGNhbmNlbGxhdGlvbiByZWFzb24sXG4gICAgICAgICAgICAvLyBoZW5jZSB0aGUgaGFuZGxlciB3aWxsIG5ldmVyIHRocm93LlxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnByb3RvdHlwZS50aGVuLmNhbGwocHJvbWlzZS5wcm9taXNlLCB1bmRlZmluZWQsIChlcnIpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoZXJyICE9PSByZWFzb24pIHtcbiAgICAgICAgICAgICAgICAgICAgdGhyb3cgZXJyO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0pO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gSWYgcmVhc29uIGlzIG5vdCBzZXQsIHRoZSBwcm9taXNlIHJlc29sdmVkIHJlZ3VsYXJseSwgaGVuY2Ugd2UgbXVzdCBub3QgY2FsbCBvbmNhbmNlbGxlZC5cbiAgICAgICAgLy8gSWYgb25jYW5jZWxsZWQgaXMgdW5zZXQsIG5vIG5lZWQgdG8gZ28gYW55IGZ1cnRoZXIuXG4gICAgICAgIGlmICghc3RhdGUucmVhc29uIHx8ICFwcm9taXNlLm9uY2FuY2VsbGVkKSB7IHJldHVybjsgfVxuXG4gICAgICAgIGNhbmNlbGxhdGlvblByb21pc2UgPSBuZXcgUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICByZXNvbHZlKHByb21pc2Uub25jYW5jZWxsZWQhKHN0YXRlLnJlYXNvbiEuY2F1c2UpKTtcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pLmNhdGNoKChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xuICAgICAgICB9KTtcblxuICAgICAgICAvLyBVbnNldCBvbmNhbmNlbGxlZCB0byBwcmV2ZW50IHJlcGVhdGVkIGNhbGxzLlxuICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gbnVsbDtcblxuICAgICAgICByZXR1cm4gY2FuY2VsbGF0aW9uUHJvbWlzZTtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVzb2x1dGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICovXG5mdW5jdGlvbiByZXNvbHZlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiB7XG4gICAgcmV0dXJuICh2YWx1ZSkgPT4ge1xuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xuXG4gICAgICAgIGlmICh2YWx1ZSA9PT0gcHJvbWlzZS5wcm9taXNlKSB7XG4gICAgICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QobmV3IFR5cGVFcnJvcihcIkEgcHJvbWlzZSBjYW5ub3QgYmUgcmVzb2x2ZWQgd2l0aCBpdHNlbGYuXCIpKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh2YWx1ZSAhPSBudWxsICYmICh0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ2Z1bmN0aW9uJykpIHtcbiAgICAgICAgICAgIGxldCB0aGVuOiBhbnk7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHRoZW4gPSAodmFsdWUgYXMgYW55KS50aGVuO1xuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKHRoZW4pKSB7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgbGV0IGNhbmNlbCA9ICh2YWx1ZSBhcyBhbnkpLmNhbmNlbDtcbiAgICAgICAgICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUoY2FuY2VsKSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3Qgb25jYW5jZWxsZWQgPSAoY2F1c2U/OiBhbnkpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xuICAgICAgICAgICAgICAgICAgICAgICAgfTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChzdGF0ZS5yZWFzb24pIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBJZiBhbHJlYWR5IGNhbmNlbGxlZCwgcHJvcGFnYXRlIGNhbmNlbGxhdGlvbi5cbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSBjYW5jZWxsZXIgYWxnb3JpdGhtIGRvZXMgbm90IHJlamVjdFxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIHNvIGl0IGNhbiBiZSBkaXNjYXJkZWQgc2FmZWx5LlxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgY2FuY2VsbGVyRm9yKHsgLi4ucHJvbWlzZSwgb25jYW5jZWxsZWQgfSwgc3RhdGUpKHN0YXRlLnJlYXNvbik7XG4gICAgICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBvbmNhbmNlbGxlZDtcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0gY2F0Y2gge31cblxuICAgICAgICAgICAgICAgIGNvbnN0IG5ld1N0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcbiAgICAgICAgICAgICAgICAgICAgcm9vdDogc3RhdGUucm9vdCxcbiAgICAgICAgICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHNldHRsZWQoKSB7IHJldHVybiB0aGlzLnJvb3Quc2V0dGxlZCB9LFxuICAgICAgICAgICAgICAgICAgICBzZXQgc2V0dGxlZCh2YWx1ZSkgeyB0aGlzLnJvb3Quc2V0dGxlZCA9IHZhbHVlOyB9LFxuICAgICAgICAgICAgICAgICAgICBnZXQgcmVhc29uKCkgeyByZXR1cm4gdGhpcy5yb290LnJlYXNvbiB9XG4gICAgICAgICAgICAgICAgfTtcblxuICAgICAgICAgICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgbmV3U3RhdGUpO1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkodGhlbiwgdmFsdWUsIFtyZXNvbHZlckZvcihwcm9taXNlLCBuZXdTdGF0ZSksIHJlamVjdG9yXSk7XG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJldHVybjsgLy8gSU1QT1JUQU5UIVxuICAgICAgICAgICAgfVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICBwcm9taXNlLnJlc29sdmUodmFsdWUpO1xuICAgIH07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVqZWN0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKi9cbmZ1bmN0aW9uIHJlamVjdG9yRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yIHtcbiAgICByZXR1cm4gKHJlYXNvbj8pID0+IHtcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcblxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkge1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBpZiAocmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgc3RhdGUucmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgT2JqZWN0LmlzKHJlYXNvbi5jYXVzZSwgc3RhdGUucmVhc29uLmNhdXNlKSkge1xuICAgICAgICAgICAgICAgICAgICAvLyBTd2FsbG93IGxhdGUgcmVqZWN0aW9ucyB0aGF0IGFyZSBDYW5jZWxFcnJvcnMgd2hvc2UgY2FuY2VsbGF0aW9uIGNhdXNlIGlzIHRoZSBzYW1lIGFzIG91cnMuXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9IGNhdGNoIHt9XG5cbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogQ2FuY2VscyBhbGwgdmFsdWVzIGluIGFuIGFycmF5IHRoYXQgbG9vayBsaWtlIGNhbmNlbGxhYmxlIHRoZW5hYmxlcy5cbiAqIFJldHVybnMgYSBwcm9taXNlIHRoYXQgZnVsZmlsbHMgb25jZSBhbGwgY2FuY2VsbGF0aW9uIHByb2NlZHVyZXMgZm9yIHRoZSBnaXZlbiB2YWx1ZXMgaGF2ZSBzZXR0bGVkLlxuICovXG5mdW5jdGlvbiBjYW5jZWxBbGwocGFyZW50OiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHZhbHVlczogYW55W10sIGNhdXNlPzogYW55KTogUHJvbWlzZTx2b2lkPiB7XG4gICAgY29uc3QgcmVzdWx0czogUHJvbWlzZTx2b2lkPltdID0gW107XG5cbiAgICBmb3IgKGNvbnN0IHZhbHVlIG9mIHZhbHVlcykge1xuICAgICAgICBsZXQgY2FuY2VsOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXI7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUodmFsdWUudGhlbikpIHsgY29udGludWU7IH1cbiAgICAgICAgICAgIGNhbmNlbCA9IHZhbHVlLmNhbmNlbDtcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZShjYW5jZWwpKSB7IGNvbnRpbnVlOyB9XG4gICAgICAgIH0gY2F0Y2ggeyBjb250aW51ZTsgfVxuXG4gICAgICAgIGxldCByZXN1bHQ6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHJlc3VsdCA9IFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xuICAgICAgICAgICAgY29udGludWU7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIXJlc3VsdCkgeyBjb250aW51ZTsgfVxuICAgICAgICByZXN1bHRzLnB1c2goXG4gICAgICAgICAgICAocmVzdWx0IGluc3RhbmNlb2YgUHJvbWlzZSAgPyByZXN1bHQgOiBQcm9taXNlLnJlc29sdmUocmVzdWx0KSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICApO1xuICAgIH1cblxuICAgIHJldHVybiBQcm9taXNlLmFsbChyZXN1bHRzKSBhcyBhbnk7XG59XG5cbi8qKlxuICogUmV0dXJucyBpdHMgYXJndW1lbnQuXG4gKi9cbmZ1bmN0aW9uIGlkZW50aXR5PFQ+KHg6IFQpOiBUIHtcbiAgICByZXR1cm4geDtcbn1cblxuLyoqXG4gKiBUaHJvd3MgaXRzIGFyZ3VtZW50LlxuICovXG5mdW5jdGlvbiB0aHJvd2VyKHJlYXNvbj86IGFueSk6IG5ldmVyIHtcbiAgICB0aHJvdyByZWFzb247XG59XG5cbi8qKlxuICogQXR0ZW1wdHMgdmFyaW91cyBzdHJhdGVnaWVzIHRvIGNvbnZlcnQgYW4gZXJyb3IgdG8gYSBzdHJpbmcuXG4gKi9cbmZ1bmN0aW9uIGVycm9yTWVzc2FnZShlcnI6IGFueSk6IHN0cmluZyB7XG4gICAgdHJ5IHtcbiAgICAgICAgaWYgKGVyciBpbnN0YW5jZW9mIEVycm9yIHx8IHR5cGVvZiBlcnIgIT09ICdvYmplY3QnIHx8IGVyci50b1N0cmluZyAhPT0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZykge1xuICAgICAgICAgICAgcmV0dXJuIFwiXCIgKyBlcnI7XG4gICAgICAgIH1cbiAgICB9IGNhdGNoIHt9XG5cbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gSlNPTi5zdHJpbmdpZnkoZXJyKTtcbiAgICB9IGNhdGNoIHt9XG5cbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZy5jYWxsKGVycik7XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgcmV0dXJuIFwiPGNvdWxkIG5vdCBjb252ZXJ0IGVycm9yIHRvIHN0cmluZz5cIjtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGJhcnJpZXIgcHJvbWlzZSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuIElmIG5lY2Vzc2FyeSwgaW5pdGlhbGlzZXMgdGhlIGJhcnJpZXIuXG4gKi9cbmZ1bmN0aW9uIGN1cnJlbnRCYXJyaWVyPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPik6IFByb21pc2U8dm9pZD4ge1xuICAgIGxldCBwd3I6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0gcHJvbWlzZVtiYXJyaWVyU3ltXSA/PyB7fTtcbiAgICBpZiAoISgncHJvbWlzZScgaW4gcHdyKSkge1xuICAgICAgICBPYmplY3QuYXNzaWduKHB3ciwgcHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4oKSk7XG4gICAgfVxuICAgIGlmIChwcm9taXNlW2JhcnJpZXJTeW1dID09IG51bGwpIHtcbiAgICAgICAgcHdyLnJlc29sdmUhKCk7XG4gICAgICAgIHByb21pc2VbYmFycmllclN5bV0gPSBwd3I7XG4gICAgfVxuICAgIHJldHVybiBwd3IucHJvbWlzZSE7XG59XG5cbi8vIFBvbHlmaWxsIFByb21pc2Uud2l0aFJlc29sdmVycy5cbmxldCBwcm9taXNlV2l0aFJlc29sdmVycyA9IFByb21pc2Uud2l0aFJlc29sdmVycztcbmlmIChwcm9taXNlV2l0aFJlc29sdmVycyAmJiB0eXBlb2YgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPT09ICdmdW5jdGlvbicpIHtcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IHByb21pc2VXaXRoUmVzb2x2ZXJzLmJpbmQoUHJvbWlzZSk7XG59IGVsc2Uge1xuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gZnVuY3Rpb24gPFQ+KCk6IFByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgUHJvbWlzZTxUPigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcbiAgICAgICAgcmV0dXJuIHsgcHJvbWlzZSwgcmVzb2x2ZSwgcmVqZWN0IH07XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNsaXBib2FyZCk7XG5cbmNvbnN0IENsaXBib2FyZFNldFRleHQgPSAwO1xuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XG5cbi8qKlxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxuICpcbiAqIEBwYXJhbSB0ZXh0IC0gVGhlIHRleHQgdG8gYmUgc2V0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKiBAcmV0dXJuIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0OiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRTZXRUZXh0LCB7dGV4dH0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSB0ZXh0IGZyb20gdGhlIENsaXBib2FyZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKTogUHJvbWlzZTxzdHJpbmc+IHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuZXhwb3J0IGludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cbiAgICBXaWR0aDogbnVtYmVyO1xuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cbiAgICBIZWlnaHQ6IG51bWJlcjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBSZWN0IHtcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xuICAgIFg6IG51bWJlcjtcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xuICAgIFk6IG51bWJlcjtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuICovXG4gICAgV2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xuICAgIEhlaWdodDogbnVtYmVyO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFNjcmVlbiB7XG4gICAgLyoqIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLiAqL1xuICAgIElEOiBzdHJpbmc7XG4gICAgLyoqIEh1bWFuLXJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBOYW1lOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBzY2FsZSBmYWN0b3Igb2YgdGhlIHNjcmVlbiAoRFBJLzk2KS4gMSA9IHN0YW5kYXJkIERQSSwgMiA9IEhpRFBJIChSZXRpbmEpLCBldGMuICovXG4gICAgU2NhbGVGYWN0b3I6IG51bWJlcjtcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFg6IG51bWJlcjtcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFk6IG51bWJlcjtcbiAgICAvKiogQ29udGFpbnMgdGhlIHdpZHRoIGFuZCBoZWlnaHQgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBTaXplOiBTaXplO1xuICAgIC8qKiBDb250YWlucyB0aGUgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQuICovXG4gICAgQm91bmRzOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQgKGJlZm9yZSBzY2FsaW5nKS4gKi9cbiAgICBQaHlzaWNhbEJvdW5kczogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIGFyZWEgb2YgdGhlIHNjcmVlbiB0aGF0IGlzIGFjdHVhbGx5IHVzYWJsZSAoZXhjbHVkaW5nIHRhc2tiYXIgYW5kIG90aGVyIHN5c3RlbSBVSSkuICovXG4gICAgV29ya0FyZWE6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBXb3JrQXJlYSBvZiB0aGUgc2NyZWVuIChiZWZvcmUgc2NhbGluZykuICovXG4gICAgUGh5c2ljYWxXb3JrQXJlYTogUmVjdDtcbiAgICAvKiogVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXG4gICAgSXNQcmltYXJ5OiBib29sZWFuO1xuICAgIC8qKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBSb3RhdGlvbjogbnVtYmVyO1xufVxuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlNjcmVlbnMpO1xuXG5jb25zdCBnZXRBbGwgPSAwO1xuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XG5jb25zdCBnZXRDdXJyZW50ID0gMjtcbmNvbnN0IGdldEJ5SUQgPSAzO1xuY29uc3QgZ2V0QnlJbmRleCA9IDQ7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpOiBQcm9taXNlPFNjcmVlbltdPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcbn1cblxuLyoqXG4gKiBHZXRzIGEgc2NyZWVuIGJ5IGl0cyB1bmlxdWUgZGlzcGxheSBJRC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgdW5pcXVlIGlkZW50aWZpZXIgb2YgdGhlIHNjcmVlbi5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBtYXRjaGluZyBTY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRCeUlEKGlkOiBzdHJpbmcpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgIHJldHVybiBjYWxsKGdldEJ5SUQsIHsgaWQgfSk7XG59XG5cbi8qKlxuICogR2V0cyBhIHNjcmVlbiBieSBpdHMgaW5kZXggaW4gdGhlIHNjcmVlbiBsaXN0LlxuICpcbiAqIEBwYXJhbSBpbmRleCAtIFRoZSB6ZXJvLWJhc2VkIGluZGV4IG9mIHRoZSBzY3JlZW4uXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgbWF0Y2hpbmcgU2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QnlJbmRleChpbmRleDogbnVtYmVyKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRCeUluZGV4LCB7IGluZGV4IH0pO1xufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuSU9TKTtcblxuLy8gTWV0aG9kIElEc1xuY29uc3QgSGFwdGljc0ltcGFjdCA9IDA7XG5jb25zdCBEZXZpY2VJbmZvID0gMTtcblxuZXhwb3J0IG5hbWVzcGFjZSBIYXB0aWNzIHtcbiAgICBleHBvcnQgdHlwZSBJbXBhY3RTdHlsZSA9IFwibGlnaHRcInxcIm1lZGl1bVwifFwiaGVhdnlcInxcInNvZnRcInxcInJpZ2lkXCI7XG4gICAgZXhwb3J0IGZ1bmN0aW9uIEltcGFjdChzdHlsZTogSW1wYWN0U3R5bGUgPSBcIm1lZGl1bVwiKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiBjYWxsKEhhcHRpY3NJbXBhY3QsIHsgc3R5bGUgfSk7XG4gICAgfVxufVxuXG5leHBvcnQgbmFtZXNwYWNlIERldmljZSB7XG4gICAgZXhwb3J0IGludGVyZmFjZSBJbmZvIHtcbiAgICAgICAgbW9kZWw6IHN0cmluZztcbiAgICAgICAgc3lzdGVtTmFtZTogc3RyaW5nO1xuICAgICAgICBzeXN0ZW1WZXJzaW9uOiBzdHJpbmc7XG4gICAgICAgIGlzU2ltdWxhdG9yOiBib29sZWFuO1xuICAgIH1cbiAgICBleHBvcnQgZnVuY3Rpb24gSW5mbygpOiBQcm9taXNlPEluZm8+IHtcbiAgICAgICAgcmV0dXJuIGNhbGwoRGV2aWNlSW5mbyk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIFVwZGF0ZXIgZXZlbnQgbmFtZSBjb25zdGFudHMuXG4gKlxuICogVXNlIHRoZXNlIGluc3RlYWQgb2YgaGFyZC1jb2Rpbmcgc3RyaW5nIGxpdGVyYWxzIHdoZW4gc3Vic2NyaWJpbmcgdG9cbiAqIHVwZGF0ZXIgZXZlbnRzIGZyb20gSmF2YVNjcmlwdDpcbiAqXG4gKiAgICAgaW1wb3J0IHsgRXZlbnRzLCBVcGRhdGVyIH0gZnJvbSBcIkB3YWlsc2lvL3J1bnRpbWVcIjtcbiAqXG4gKiAgICAgRXZlbnRzLk9uKFVwZGF0ZXIuRXZlbnRzLlVwZGF0ZUF2YWlsYWJsZSwgKGUpID0+IHtcbiAqICAgICAgICAgY29uc29sZS5sb2coXCJ1cGRhdGUgZm91bmQ6XCIsIGUuZGF0YS52ZXJzaW9uKTtcbiAqICAgICB9KTtcbiAqXG4gKiAgICAgRXZlbnRzLk9uKFVwZGF0ZXIuRXZlbnRzLkRvd25sb2FkUHJvZ3Jlc3MsIChlKSA9PiB7XG4gKiAgICAgICAgIGNvbnN0IHAgPSBlLmRhdGE7XG4gKiAgICAgICAgIGNvbnNvbGUubG9nKGAke3Aud3JpdHRlbn0gLyAke3AudG90YWx9IGJ5dGVzYCk7XG4gKiAgICAgfSk7XG4gKlxuICogTWlycm9ycyB0aGUgR28tc2lkZSBjb25zdGFudHMgaW4gYHBrZy91cGRhdGVyL2V2ZW50cy5nb2AgYW5kIHRoZVxuICogdXNlci1hY3Rpb24gY29uc3RhbnRzIGluIGBwa2cvdXBkYXRlci93aW5kb3dfbGlmZWN5Y2xlLmdvYC4gQW55XG4gKiBjaGFuZ2VzIGhlcmUgbXVzdCBzdGF5IGluIHN5bmMgd2l0aCB0aG9zZSBmaWxlcyBcdTIwMTQgdGhlcmUncyBhblxuICogaW50ZWdyYXRpb24gdGVzdCB0aGF0IGFzc2VydHMgdGhlIHN0cmluZ3MgbWF0Y2guXG4gKi9cbmV4cG9ydCBjb25zdCBFdmVudHMgPSBPYmplY3QuZnJlZXplKHtcbiAgICAvKiogQSBDaGVjayByb3VuZC10cmlwIGlzIHN0YXJ0aW5nLiBQYXlsb2FkOiBudWxsLiAqL1xuICAgIENoZWNrU3RhcnRlZDogXCJ3YWlsczp1cGRhdGVyOmNoZWNrLXN0YXJ0ZWRcIixcbiAgICAvKiogQ2hlY2sgZm91bmQgYSBuZXdlciByZWxlYXNlLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xuICAgIFVwZGF0ZUF2YWlsYWJsZTogXCJ3YWlsczp1cGRhdGVyOnVwZGF0ZS1hdmFpbGFibGVcIixcbiAgICAvKiogQ2hlY2sgY29uZmlybWVkIHRoZSBjYWxsZXIgaXMgdXAgdG8gZGF0ZS4gUGF5bG9hZDogbnVsbC4gKi9cbiAgICBOb1VwZGF0ZTogXCJ3YWlsczp1cGRhdGVyOm5vLXVwZGF0ZVwiLFxuICAgIC8qKiBEb3dubG9hZCBpcyBzdGFydGluZy4gUGF5bG9hZDogUmVsZWFzZS4gKi9cbiAgICBEb3dubG9hZFN0YXJ0ZWQ6IFwid2FpbHM6dXBkYXRlcjpkb3dubG9hZC1zdGFydGVkXCIsXG4gICAgLyoqIFBlcmlvZGljIHByb2dyZXNzIHRpY2sgZHVyaW5nIGRvd25sb2FkICh+MTAgSHopLiBQYXlsb2FkOiBQcm9ncmVzcy4gKi9cbiAgICBEb3dubG9hZFByb2dyZXNzOiBcIndhaWxzOnVwZGF0ZXI6ZG93bmxvYWQtcHJvZ3Jlc3NcIixcbiAgICAvKiogQWxsIGJ5dGVzIGFyZSBvbiBkaXNrLCBidXQgdmVyaWZpY2F0aW9uIGhhcyBub3QgeWV0IHN0YXJ0ZWQuIFBheWxvYWQ6IFJlbGVhc2UuICovXG4gICAgRG93bmxvYWRDb21wbGV0ZTogXCJ3YWlsczp1cGRhdGVyOmRvd25sb2FkLWNvbXBsZXRlXCIsXG4gICAgLyoqIFNpZ25hdHVyZSAvIGRpZ2VzdCB2ZXJpZmljYXRpb24gaGFzIHN0YXJ0ZWQuIFBheWxvYWQ6IFJlbGVhc2UuICovXG4gICAgVmVyaWZ5aW5nOiBcIndhaWxzOnVwZGF0ZXI6dmVyaWZ5aW5nXCIsXG4gICAgLyoqIFRoZSBVcGRhdGVyIGlzIHN3YXBwaW5nIHRoZSBiaW5hcnkgaW50byBwbGFjZS4gUGF5bG9hZDogUmVsZWFzZS4gKi9cbiAgICBJbnN0YWxsaW5nOiBcIndhaWxzOnVwZGF0ZXI6aW5zdGFsbGluZ1wiLFxuICAgIC8qKiBVcGRhdGUgaXMgc3RhZ2VkIGFuZCBhIHJlc3RhcnQgaXMgcGVuZGluZy4gUGF5bG9hZDogUmVsZWFzZS4gKi9cbiAgICBVcGRhdGVSZWFkeTogXCJ3YWlsczp1cGRhdGVyOnVwZGF0ZS1yZWFkeVwiLFxuICAgIC8qKiBTb21ldGhpbmcgZmFpbGVkLiBQYXlsb2FkOiBFcnJvckluZm8geyBzdGFnZSwgbWVzc2FnZSwgcHJvdmlkZXIgfS4gKi9cbiAgICBFcnJvcjogXCJ3YWlsczp1cGRhdGVyOmVycm9yXCIsXG4gICAgLyoqIEhvc3Qtc2lkZSBjb250ZXh0IGRlbGl2ZXJlZCBvbmNlIHBlciBzZXNzaW9uLiBQYXlsb2FkOiBNZXRhIHsgY3VycmVudFZlcnNpb24sIHNraXBwZWRWZXJzaW9uIH0uICovXG4gICAgTWV0YTogXCJ3YWlsczp1cGRhdGVyOm1ldGFcIixcblxuICAgIC8qKiBTdWItbmFtZXNwYWNlOiB1c2VyLWFjdGlvbiBldmVudHMgdGhhdCB0aGUgVUkgZW1pdHMgQkFDSyB0byB0aGUgaG9zdC4gKi9cbiAgICBVc2VyOiBPYmplY3QuZnJlZXplKHtcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBJbnN0YWxsIG9uIGFuIEF2YWlsYWJsZSB1cGRhdGUuICovXG4gICAgICAgIEluc3RhbGw6IFwid2FpbHM6dXBkYXRlcjp1c2VyOmluc3RhbGxcIixcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBSZXN0YXJ0ICYgQXBwbHkgb24gYSBSZWFkeSB1cGRhdGUuICovXG4gICAgICAgIFJlc3RhcnQ6IFwid2FpbHM6dXBkYXRlcjp1c2VyOnJlc3RhcnRcIixcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBTa2lwIFRoaXMgVmVyc2lvbi4gKi9cbiAgICAgICAgU2tpcDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6c2tpcFwiLFxuICAgICAgICAvKiogVXNlciBjbGlja2VkIFJlbWluZCBNZSBMYXRlci4gKi9cbiAgICAgICAgUmVtaW5kOiBcIndhaWxzOnVwZGF0ZXI6dXNlcjpyZW1pbmRcIixcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBDbG9zZSAvIENhbmNlbC4gKi9cbiAgICAgICAgQ2FuY2VsOiBcIndhaWxzOnVwZGF0ZXI6dXNlcjpjYW5jZWxcIixcbiAgICB9KSxcblxuICAgIC8qKiBTdWItbmFtZXNwYWNlOiBmcmFtZXdvcmstaW50ZXJuYWwgZXZlbnRzIHRoZSBVSSBlbWl0cyB0byBjb29yZGluYXRlXG4gICAgICogIHdpdGggdGhlIGhvc3QuIE1vc3QgYXBwIGNvZGUgY2FuIGlnbm9yZSB0aGVzZS4gKi9cbiAgICBXaW5kb3c6IE9iamVjdC5mcmVlemUoe1xuICAgICAgICAvKiogVGhlIHdpbmRvdyBmaW5pc2hlZCBsb2FkaW5nIGFuZCBhc2tzIHRoZSBob3N0IHRvIHJlcGxheSBjdXJyZW50IHN0YXRlLiAqL1xuICAgICAgICBSZWFkeTogXCJ3YWlsczp1cGRhdGVyOndpbmRvdzpyZWFkeVwiLFxuICAgIH0pLFxufSk7XG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFB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QUFBQTtBQUFBLEVBRWQsaUJBQWlCO0FBQUE7QUFBQSxFQUVqQixVQUFVO0FBQUE7QUFBQSxFQUVWLGlCQUFpQjtBQUFBO0FBQUEsRUFFakIsa0JBQWtCO0FBQUE7QUFBQSxFQUVsQixrQkFBa0I7QUFBQTtBQUFBLEVBRWxCLFdBQVc7QUFBQTtBQUFBLEVBRVgsWUFBWTtBQUFBO0FBQUEsRUFFWixhQUFhO0FBQUE7QUFBQSxFQUViLE9BQU87QUFBQTtBQUFBLEVBRVAsTUFBTTtBQUFBO0FBQUEsRUFHTixNQUFNLE9BQU8sT0FBTztBQUFBO0FBQUEsSUFFaEIsU0FBUztBQUFBO0FBQUEsSUFFVCxTQUFTO0FBQUE7QUFBQSxJQUVULE1BQU07QUFBQTtBQUFBLElBRU4sUUFBUTtBQUFBO0FBQUEsSUFFUixRQUFRO0FBQUEsRUFDWixDQUFDO0FBQUE7QUFBQTtBQUFBLEVBSUQsUUFBUSxPQUFPLE9BQU87QUFBQTtBQUFBLElBRWxCLE9BQU87QUFBQSxFQUNYLENBQUM7QUFDTCxDQUFDOzs7QXhCakVELE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQTBEbEMsT0FBTyxPQUFPLFNBQWdCO0FBQzlCLE9BQU8sT0FBTyxXQUFXO0FBS3pCLE9BQU8sT0FBTyx5QkFBeUIsZUFBTyx1QkFBdUIsS0FBSyxjQUFNO0FBR2hGLE9BQU8sT0FBTyxrQkFBa0I7QUFDaEMsT0FBTyxPQUFPLGtCQUFrQjtBQUNoQyxPQUFPLE9BQU8saUJBQWlCO0FBRXhCLE9BQU8scUJBQXFCO0FBTzVCLFNBQVMsbUJBQW1CLEtBQTRCO0FBQzNELFNBQU8sTUFBTSxLQUFLLEVBQUUsUUFBUSxPQUFPLENBQUMsRUFDL0IsS0FBSyxjQUFZO0FBQ2QsUUFBSSxTQUFTLElBQUk7QUFHYixZQUFNLGVBQWUsU0FBUyxRQUFRLElBQUksY0FBYyxLQUFLLElBQUksWUFBWTtBQUM3RSxVQUFJLFlBQVksU0FBUyxZQUFZLEdBQUc7QUFDcEMsY0FBTSxTQUFTLFNBQVMsY0FBYyxRQUFRO0FBQzlDLGVBQU8sTUFBTTtBQUNiLGlCQUFTLEtBQUssWUFBWSxNQUFNO0FBQUEsTUFDcEM7QUFBQSxJQUNKO0FBQUEsRUFDSixDQUFDLEVBQ0EsTUFBTSxNQUFNO0FBQUEsRUFBQyxDQUFDO0FBQ3ZCO0FBR0EsbUJBQW1CLGtCQUFrQjsiLAogICJuYW1lcyI6IFsiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJfYSIsICJBcnJheSIsICJNYXAiLCAiQXJyYXkiLCAiTWFwIiwgImtleSIsICJjYWxsIiwgIl9hIiwgIl9hIiwgInJlc2l6YWJsZSIsICJfYSIsICJjYWxsIiwgIl9hIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJIaWRlTWV0aG9kIiwgIlNob3dNZXRob2QiLCAiaXNEb2N1bWVudERvdEFsbCIsICJfYSIsICJyZWFzb24iLCAidmFsdWUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIiwgIkV2ZW50cyIsICJFdmVudHMiXQp9Cg==
