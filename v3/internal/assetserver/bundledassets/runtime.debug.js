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
  Android: () => android_exports,
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
  JSONStream: () => JSONStream,
  Media: () => media_exports,
  Screens: () => screens_exports,
  Stream: () => Stream,
  System: () => system_exports,
  Updater: () => updater_exports,
  WML: () => wml_exports,
  WailsSocket: () => WailsSocket,
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
  const source = typeof crypto !== "undefined" && typeof crypto.getRandomValues === "function" ? crypto : null;
  if (source) {
    const bytes = new Uint8Array(size);
    source.getRandomValues(bytes);
    let id2 = "";
    for (let i2 = 0; i2 < size; i2++) {
      id2 += urlAlphabet[bytes[i2] & 63];
    }
    return id2;
  }
  let id = "";
  let i = size | 0;
  while (i--) {
    id += urlAlphabet[Math.random() * 64 | 0];
  }
  return id;
}

// desktop/@wailsio/runtime/src/environment.ts
var hasDOM = typeof window !== "undefined" && typeof document !== "undefined";

// desktop/@wailsio/runtime/src/runtime.ts
function runtimeURL() {
  return window.location.origin + "/wails/runtime";
}
var CHUNK_THRESHOLD = 512 * 1024;
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
  IOS: 11,
  Android: 12
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
  var _a3, _b;
  if (customTransport) {
    return customTransport.call(objectID, method, windowName, args);
  }
  let url = new URL(runtimeURL());
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
    const ct = response.headers.get("Content-Type");
    if (ct == null ? void 0 : ct.includes("application/json")) {
      const json = await response.json();
      let err;
      switch (json.kind) {
        case "ReferenceError":
          err = new ReferenceError(json.message);
          break;
        case "TypeError":
          err = new TypeError(json.message);
          break;
        case "RuntimeError":
          err = new RuntimeError(json.message);
          break;
        default:
          err = new Error(json.message);
      }
      err.cause = json.cause;
      throw err;
    }
    throw new Error(await response.text());
  }
  if (((_b = (_a3 = response.headers.get("Content-Type")) == null ? void 0 : _a3.indexOf("application/json")) != null ? _b : -1) !== -1) {
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
var _a;
var androidBridge = hasDOM && typeof ((_a = window.wails) == null ? void 0 : _a.invokeAsync) === "function" ? window.wails : null;
if (androidBridge) {
  const pending = /* @__PURE__ */ new Map();
  window._wailsAndroidCallback = (id, response, error) => {
    var _a3;
    const promise = pending.get(id);
    if (!promise) return;
    pending.delete(id);
    if (error) {
      promise.reject(new Error(error));
      return;
    }
    try {
      const envelope = JSON.parse(response != null ? response : "{}");
      if (!envelope.ok) {
        promise.reject(new Error((_a3 = envelope.error) != null ? _a3 : "unknown runtime call error"));
        return;
      }
      promise.resolve("text" in envelope ? envelope.text : envelope.data);
    } catch (e) {
      promise.reject(e);
    }
  };
  customTransport = {
    call(objectID, method, windowName, args) {
      return new Promise((resolve, reject) => {
        const id = nanoid();
        pending.set(id, { resolve, reject });
        try {
          androidBridge.invokeAsync(id, JSON.stringify({
            object: objectID,
            method,
            windowName,
            args: args != null ? args : null,
            clientId
          }));
        } catch (e) {
          pending.delete(id);
          reject(e);
        }
      });
    }
  };
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
if (hasDOM) {
  window._wails = window._wails || {};
}
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
  var _a3;
  return (_a3 = dialog(DialogOpenFile, options)) != null ? _a3 : [];
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
  DateFromTime: () => DateFromTime,
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
function DateFromTime(source) {
  return new Date(source);
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
    WindowZoomReset: "mac:WindowZoomReset",
    WebViewWebContentProcessDidTerminate: "mac:WebViewWebContentProcessDidTerminate"
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
    WebViewDecidePolicyForNavigationAction: "ios:WebViewDecidePolicyForNavigationAction",
    BatteryChanged: "ios:BatteryChanged",
    NetworkChanged: "ios:NetworkChanged",
    ThemeChanged: "ios:ThemeChanged",
    ScreenLocked: "ios:ScreenLocked",
    ScreenUnlocked: "ios:ScreenUnlocked"
  }),
  Android: Object.freeze({
    ActivityCreated: "android:ActivityCreated",
    ActivityStarted: "android:ActivityStarted",
    ActivityResumed: "android:ActivityResumed",
    ActivityPaused: "android:ActivityPaused",
    ActivityStopped: "android:ActivityStopped",
    ActivityDestroyed: "android:ActivityDestroyed",
    ApplicationLowMemory: "android:ApplicationLowMemory",
    WebViewPageStarted: "android:WebViewPageStarted",
    WebViewPageFinished: "android:WebViewPageFinished",
    BatteryChanged: "android:BatteryChanged",
    NetworkChanged: "android:NetworkChanged",
    ThemeChanged: "android:ThemeChanged",
    ScreenLocked: "android:ScreenLocked",
    ScreenUnlocked: "android:ScreenUnlocked"
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
    WindowZoomReset: "common:WindowZoomReset",
    BatteryChanged: "common:BatteryChanged",
    NetworkChanged: "common:NetworkChanged",
    ScreenLocked: "common:ScreenLocked",
    ScreenUnlocked: "common:ScreenUnlocked",
    LowMemory: "common:LowMemory"
  })
});

// desktop/@wailsio/runtime/src/events.ts
if (hasDOM) {
  window._wails = window._wails || {};
  window._wails.dispatchWailsEvent = dispatchWailsEvent;
}
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
  const expired = /* @__PURE__ */ new Set();
  for (const listener of listeners.slice()) {
    if (listener.dispatch(wailsEvent)) {
      expired.add(listener);
    }
  }
  if (expired.size > 0) {
    const live = eventListeners.get(event.name);
    if (live) {
      const remaining = live.filter((l) => !expired.has(l));
      if (remaining.length === 0) {
        eventListeners.delete(event.name);
      } else {
        eventListeners.set(event.name, remaining);
      }
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
  var _a3;
  if (event.target instanceof HTMLElement) {
    return event.target;
  } else if (!(event.target instanceof HTMLElement) && event.target instanceof Node) {
    return (_a3 = event.target.parentElement) != null ? _a3 : document.body;
  } else {
    return document.body;
  }
}
var isReady = false;
if (hasDOM) {
  document.addEventListener("DOMContentLoaded", () => {
    isReady = true;
  });
}
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
  var _a3, _b, _c, _d;
  if (((_b = (_a3 = window.chrome) == null ? void 0 : _a3.webview) == null ? void 0 : _b.postMessageWithAdditionalObjects) == null) {
    return false;
  }
  return ((_d = (_c = window._wails) == null ? void 0 : _c.flags) == null ? void 0 : _d.enableFileDrop) === true;
}
function resolveFilePaths(x, y, files) {
  var _a3, _b;
  if ((_b = (_a3 = window.chrome) == null ? void 0 : _a3.webview) == null ? void 0 : _b.postMessageWithAdditionalObjects) {
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
  var _a3, _b;
  if (((_b = (_a3 = window._wails) == null ? void 0 : _a3.flags) == null ? void 0 : _b.enableFileDrop) === false) {
    return;
  }
  nativeDragActive = true;
}
function handleDragLeave() {
  cleanupNativeDrag();
}
function handleDragOver(x, y) {
  var _a3, _b;
  if (!nativeDragActive) return;
  if (((_b = (_a3 = window._wails) == null ? void 0 : _a3.flags) == null ? void 0 : _b.enableFileDrop) === false) {
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
   * @param x - The x-coordinate of the drop event, in logical (CSS) pixels relative to the webview.
   * @param y - The y-coordinate of the drop event, in logical (CSS) pixels relative to the webview.
   */
  HandlePlatformFileDrop(filenames, x, y) {
    var _a3, _b;
    if (((_b = (_a3 = window._wails) == null ? void 0 : _a3.flags) == null ? void 0 : _b.enableFileDrop) === false) {
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
    var _a3, _b, _c;
    if (!((_a3 = event.dataTransfer) == null ? void 0 : _a3.types.includes("Files"))) {
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
    var _a3, _b, _c;
    if (!((_a3 = event.dataTransfer) == null ? void 0 : _a3.types.includes("Files"))) {
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
    var _a3, _b, _c;
    if (!((_a3 = event.dataTransfer) == null ? void 0 : _a3.types.includes("Files"))) {
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
    var _a3, _b, _c;
    if (!((_a3 = event.dataTransfer) == null ? void 0 : _a3.types.includes("Files"))) {
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
  IsAndroid: () => IsAndroid,
  IsDarkMode: () => IsDarkMode,
  IsDebug: () => IsDebug,
  IsDesktop: () => IsDesktop,
  IsIOS: () => IsIOS,
  IsLinux: () => IsLinux,
  IsMac: () => IsMac,
  IsMobile: () => IsMobile,
  IsWindows: () => IsWindows,
  invoke: () => invoke
});
var call4 = newRuntimeCaller(objectNames.System);
var SystemIsDarkMode = 0;
var SystemEnvironment = 1;
var SystemCapabilities = 2;
var _invoke = (function() {
  var _a3, _b, _c, _d, _e, _f;
  try {
    if ((_b = (_a3 = window.chrome) == null ? void 0 : _a3.webview) == null ? void 0 : _b.postMessage) {
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
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.OS) === "windows";
}
function IsLinux() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.OS) === "linux";
}
function IsMac() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.OS) === "darwin";
}
function IsIOS() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.OS) === "ios";
}
function IsAndroid() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.OS) === "android";
}
function IsMobile() {
  return IsIOS() || IsAndroid();
}
function IsDesktop() {
  return IsMac() || IsWindows() || IsLinux();
}
function IsAMD64() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.Arch) === "amd64";
}
function IsARM() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.Arch) === "arm";
}
function IsARM64() {
  var _a3, _b;
  return ((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.Arch) === "arm64";
}
function IsDebug() {
  var _a3, _b;
  return Boolean((_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.Debug);
}

// desktop/@wailsio/runtime/src/contextmenu.ts
if (hasDOM) {
  window.addEventListener("contextmenu", contextMenuHandler);
}
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
var buttonsTracked = false;
if (hasDOM) {
  buttonsTracked = canTrackButtons();
  window._wails = window._wails || {};
  window._wails.setResizable = (value) => {
    resizable = value;
    if (!resizable) {
      canResize = resizing = false;
      setResize();
    }
  };
}
var dragInitDone = false;
function isMobile() {
  var _a3, _b;
  const os = (_b = (_a3 = window._wails) == null ? void 0 : _a3.environment) == null ? void 0 : _b.OS;
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
if (hasDOM) {
  tryInitDragHandlers();
  document.addEventListener("DOMContentLoaded", tryInitDragHandlers, { once: true });
  let dragEnvPolls = 0;
  const dragEnvPoll = window.setInterval(() => {
    if (dragInitDone) {
      window.clearInterval(dragEnvPoll);
      return;
    }
    tryInitDragHandlers();
    if (++dragEnvPolls > 100) {
      window.clearInterval(dragEnvPoll);
    }
  }, 50);
}
function suppressEvent(event) {
  if (event.type === "dblclick" && IsMac() && window === window.top && isDraggableEvent(event)) {
    event.stopImmediatePropagation();
    event.stopPropagation();
    event.preventDefault();
    invoke("wails:drag:doubleclick");
    return;
  }
  if (dragging || resizing) {
    event.stopImmediatePropagation();
    event.stopPropagation();
    event.preventDefault();
  }
}
function isDraggableEvent(event) {
  const target = eventTarget(event);
  const style = window.getComputedStyle(target);
  return style.getPropertyValue("--wails-draggable").trim() === "drag" && event.offsetX >= 0 && event.offsetX < target.clientWidth && event.offsetY >= 0 && event.offsetY < target.clientHeight;
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
    if (event.type !== "mousedown") {
      return;
    }
    canResize = true;
    return;
  }
  canDrag = isDraggableEvent(event);
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
  if (!resizable || !IsWindows() && !(IsLinux() && GetFlag("frameless"))) {
    if (resizeEdge) {
      setResize();
    }
    return;
  }
  const resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
  const resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
  const cornerExtra = GetFlag("resizeCornerExtra") || 10;
  const scrollbarWidth = Math.max(0, window.innerWidth - document.documentElement.clientWidth);
  const scrollbarHeight = Math.max(0, window.innerHeight - document.documentElement.clientHeight);
  const rightContentEdge = window.innerWidth - scrollbarWidth;
  const bottomContentEdge = window.innerHeight - scrollbarHeight;
  const rightBorder = event.clientX < rightContentEdge && rightContentEdge - event.clientX < resizeHandleWidth;
  const leftBorder = event.clientX < resizeHandleWidth;
  const topBorder = event.clientY < resizeHandleHeight;
  const bottomBorder = event.clientY < bottomContentEdge && bottomContentEdge - event.clientY < resizeHandleHeight;
  const rightCorner = event.clientX < rightContentEdge && rightContentEdge - event.clientX < resizeHandleWidth + cornerExtra;
  const leftCorner = event.clientX < resizeHandleWidth + cornerExtra;
  const topCorner = event.clientY < resizeHandleHeight + cornerExtra;
  const bottomCorner = event.clientY < bottomContentEdge && bottomContentEdge - event.clientY < resizeHandleHeight + cornerExtra;
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

// desktop/@wailsio/runtime/src/appregion.ts
var regionProperty = "--wails-non-client-region";
var runtimeConfigReadyEvent = "wails:runtime-config-ready";
var validRegions = /* @__PURE__ */ new Set(["caption", "minimize", "maximize", "close"]);
if (hasDOM) {
  window._wails = window._wails || {};
}
var updatePending = false;
var lastPayload = "";
var observedElements = /* @__PURE__ */ new Set();
var resizeObserver;
var trackingStarted = false;
function normaliseRegionKind(value) {
  const region = value.trim().toLowerCase();
  if (validRegions.has(region)) {
    return region;
  }
  return void 0;
}
function nonClientRegionForElement(element) {
  if (!(element instanceof HTMLElement)) {
    return void 0;
  }
  const style = window.getComputedStyle(element);
  const region = normaliseRegionKind(style.getPropertyValue(regionProperty));
  if (!region) {
    return void 0;
  }
  const parent = element.parentElement;
  if (parent) {
    const parentStyle = window.getComputedStyle(parent);
    if (normaliseRegionKind(parentStyle.getPropertyValue(regionProperty)) === region) {
      return void 0;
    }
  }
  return region;
}
function isVisible(element) {
  const style = window.getComputedStyle(element);
  return style.display !== "none" && style.visibility !== "hidden" && style.contentVisibility !== "hidden";
}
function elementRegion(element) {
  if (!(element instanceof HTMLElement)) {
    return void 0;
  }
  const kind = nonClientRegionForElement(element);
  if (!kind || !isVisible(element)) {
    return void 0;
  }
  const rect = element.getBoundingClientRect();
  if (rect.width <= 0 || rect.height <= 0) {
    return void 0;
  }
  const scale = window.devicePixelRatio || 1;
  const left = Math.floor(rect.left * scale);
  const top = Math.floor(rect.top * scale);
  const right = Math.ceil(rect.right * scale);
  const bottom = Math.ceil(rect.bottom * scale);
  if (right <= left || bottom <= top) {
    return void 0;
  }
  return { kind, left, top, right, bottom };
}
function regionElements() {
  const elements = [];
  if (document.documentElement) {
    elements.push(document.documentElement);
  }
  if (document.body) {
    elements.push(document.body);
    for (const element of document.body.querySelectorAll("*")) {
      elements.push(element);
    }
  }
  return elements;
}
function observeRegionElements(elements) {
  if (typeof ResizeObserver === "undefined") {
    return;
  }
  resizeObserver != null ? resizeObserver : resizeObserver = new ResizeObserver(scheduleUpdate);
  const nextElements = new Set(elements);
  for (const element of observedElements) {
    if (!nextElements.has(element)) {
      resizeObserver.unobserve(element);
    }
  }
  for (const element of nextElements) {
    if (!observedElements.has(element)) {
      resizeObserver.observe(element);
    }
  }
  observedElements = nextElements;
}
function updateNonClientRegions() {
  updatePending = false;
  const elements = regionElements();
  const regions = [];
  const activeElements = [];
  for (const element of elements) {
    const region = elementRegion(element);
    if (region) {
      regions.push(region);
      activeElements.push(element);
    }
  }
  observeRegionElements(activeElements);
  const payload = JSON.stringify({ version: 1, regions });
  if (payload === lastPayload) {
    return;
  }
  lastPayload = payload;
  invoke("wails:non-client-region:" + payload);
}
function scheduleUpdate() {
  if (updatePending) {
    return;
  }
  updatePending = true;
  window.requestAnimationFrame(updateNonClientRegions);
}
function startNonClientRegionTracking() {
  var _a3, _b;
  if (trackingStarted) {
    return;
  }
  trackingStarted = true;
  scheduleUpdate();
  const mutationObserver = new MutationObserver(scheduleUpdate);
  mutationObserver.observe(document.documentElement, {
    attributes: true,
    childList: true,
    subtree: true
  });
  window.addEventListener("resize", scheduleUpdate);
  window.addEventListener("scroll", scheduleUpdate, true);
  (_a3 = window.visualViewport) == null ? void 0 : _a3.addEventListener("resize", scheduleUpdate);
  (_b = window.visualViewport) == null ? void 0 : _b.addEventListener("scroll", scheduleUpdate);
}
function tryStartNonClientRegionTracking() {
  var _a3, _b;
  const os = (_a3 = window._wails.environment) == null ? void 0 : _a3.OS;
  if (os === void 0) {
    return false;
  }
  const enabled = (_b = window._wails.flags) == null ? void 0 : _b.nonClientRegionTracking;
  if (os === "windows") {
    if (enabled === true) {
      whenReady(startNonClientRegionTracking);
    }
    return true;
  }
  return true;
}
if (hasDOM && !tryStartNonClientRegionTracking()) {
  window.addEventListener(runtimeConfigReadyEvent, tryStartNonClientRegionTracking, { once: true });
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
var _a2;
var species = (_a2 = Symbol.species) != null ? _a2 : /* @__PURE__ */ Symbol("speciesPolyfill");
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
          var _a3;
          if (this[barrierSym] === barrier) {
            this[barrierSym] = null;
          }
          (_a3 = barrier.resolve) == null ? void 0 : _a3.call(barrier);
          try {
            resolve(onfulfilled(value));
          } catch (err) {
            reject(err);
          }
        },
        (reason) => {
          var _a3;
          if (this[barrierSym] === barrier) {
            this[barrierSym] = null;
          }
          (_a3 = barrier.resolve) == null ? void 0 : _a3.call(barrier);
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
      var _a3;
      (_a3 = result.oncancelled) == null ? void 0 : _a3.call(result, cause);
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
  var _a3;
  let pwr = (_a3 = promise[barrierSym]) != null ? _a3 : {};
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
if (hasDOM) {
  window._wails = window._wails || {};
}
var call7 = newRuntimeCaller(objectNames.Call);
var cancelCall = newRuntimeCaller(objectNames.CancelCall);
var callResponses = /* @__PURE__ */ new Map();
var CallBinding = 0;
var CancelMethod = 0;
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

// desktop/@wailsio/runtime/src/media.ts
var media_exports = {};
__export(media_exports, {
  ClearSource: () => ClearSource,
  SetSource: () => SetSource
});

// desktop/@wailsio/runtime/src/stream.ts
var GENERATION_KEY = "__wails_stream_generation";
var GENERATION_NAME_MARKER = "__wails_stream_generation__:";
function namedPageGeneration() {
  try {
    const marker = window.name.lastIndexOf(GENERATION_NAME_MARKER);
    if (marker < 0) return 0;
    const value = Number(window.name.slice(marker + GENERATION_NAME_MARKER.length));
    return Number.isSafeInteger(value) && value > 0 ? value : 0;
  } catch (e) {
    return 0;
  }
}
function rememberNamedPageGeneration(generation) {
  try {
    const marker = window.name.lastIndexOf(GENERATION_NAME_MARKER);
    const applicationName = marker < 0 ? window.name : window.name.slice(0, marker);
    window.name = applicationName + GENERATION_NAME_MARKER + generation;
  } catch (e) {
  }
}
function nextPageGeneration() {
  if (!hasDOM) return 1;
  const origin = typeof performance !== "undefined" && Number.isFinite(performance.timeOrigin) ? performance.timeOrigin : Date.now();
  const pageTime = Math.max(1, Math.floor(origin * 1e3));
  let current = namedPageGeneration();
  try {
    const raw = window.sessionStorage.getItem(GENERATION_KEY);
    const stored = raw === null ? 0 : Number(raw);
    if (Number.isSafeInteger(stored) && stored > current) current = stored;
  } catch (e) {
  }
  const next = Number.isSafeInteger(current) && current >= 0 && current < Number.MAX_SAFE_INTEGER ? Math.max(current + 1, pageTime) : pageTime;
  try {
    window.sessionStorage.setItem(GENERATION_KEY, String(next));
  } catch (e) {
  }
  rememberNamedPageGeneration(next);
  return next;
}
function streamGlobals() {
  const fresh = () => ({
    session: nanoid(),
    generation: nextPageGeneration(),
    nextConnID: 1,
    connections: /* @__PURE__ */ new Map(),
    polling: false
  });
  if (!hasDOM) return fresh();
  const w = window._wails = window._wails || {};
  return w.__stream = w.__stream || fresh();
}
var G = streamGlobals();
var HDR_SESSION = "x-wails-stream-session";
var HDR_GENERATION = "x-wails-stream-generation";
var HDR_CONN = "x-wails-stream-conn";
var HDR_KIND = "x-wails-stream-kind";
var HDR_NAME = "x-wails-stream-name";
var HDR_CHUNK = "x-wails-stream-chunk";
var HDR_CHUNK_INDEX = "x-wails-stream-chunk-index";
var HDR_CHUNK_TOTAL = "x-wails-stream-chunk-total";
var HDR_BATCH = "x-wails-stream-batch";
var KIND_DATA = 0;
var KIND_OPEN = 1;
var KIND_CLOSE = 2;
var KIND_ERROR = 3;
var CHUNK_THRESHOLD2 = 512 * 1024;
var MAX_BATCH_FRAMES = 4096;
var BACKOFF_MIN = 250;
var BACKOFF_MAX = 5e3;
var textEncoder = new TextEncoder();
function streamURL(endpoint) {
  return window.location.origin + "/wails/stream/" + endpoint;
}
var StreamRequestCancelled = class extends Error {
  constructor() {
    super("stream request cancelled");
    this.name = "AbortError";
  }
};
function sleep(ms, signal) {
  if (!signal) return new Promise((resolve) => setTimeout(resolve, ms));
  if (signal.aborted) return Promise.reject(new StreamRequestCancelled());
  return new Promise((resolve, reject) => {
    const onAbort = () => {
      clearTimeout(timer);
      reject(new StreamRequestCancelled());
    };
    const timer = setTimeout(() => {
      signal.removeEventListener("abort", onAbort);
      resolve();
    }, ms);
    signal.addEventListener("abort", onAbort, { once: true });
  });
}
var _WailsSocket = class _WailsSocket extends EventTarget {
  constructor(name) {
    super();
    this.CONNECTING = 0;
    this.OPEN = 1;
    this.CLOSING = 2;
    this.CLOSED = 3;
    /** Always empty: subprotocols and extensions are not negotiated. */
    this.protocol = "";
    this.extensions = "";
    this.binaryType = "arraybuffer";
    this.readyState = _WailsSocket.CONNECTING;
    // Declared for the type; the constructor replaces each with an accessor
    // that registers a real listener. Calling them directly instead would make
    // any wrapper (see JSONStream) receive both the raw and the decoded event.
    this.onopen = null;
    this.onmessage = null;
    this.onclose = null;
    this.onerror = null;
    /**
     * @internal Applied to every inbound payload before dispatch. JSONStream
     * replaces it, so `onmessage` and `addEventListener("message", …)` see the
     * same decoded value — intercepting only `onmessage` made the two disagree.
     */
    this._decode = (payload) => payload;
    this._buffered = 0;
    // Sends are serialised per connection: concurrent fetch() POSTs do not
    // preserve order, and Go relies on send order being the order it observes.
    this._chain = Promise.resolve();
    this._pending = [];
    this._flushing = false;
    this._openAbort = new AbortController();
    this._sendAbort = new AbortController();
    this.name = name;
    this._id = G.nextConnID++;
    this.url = hasDOM ? streamURL("poll") + "#" + encodeURIComponent(name) : "";
    for (const type of ["open", "message", "close", "error"]) {
      defineHandlerProperty(this, type);
    }
    G.connections.set(this._id, this);
    this._chain = this._chain.then(
      () => postFrame(this._id, KIND_OPEN, new Uint8Array(0), name, this._openAbort.signal)
    ).catch((err) => {
      this._fail(err);
    });
    startPolling();
  }
  /** Bytes queued by send() that have not yet reached Go. */
  get bufferedAmount() {
    return this._buffered;
  }
  /**
   * Send a frame to Go. Accepts anything a WebSocket does; strings are
   * encoded as UTF-8, since Go receives every frame as a byte slice.
   */
  send(data) {
    if (this.readyState === _WailsSocket.CONNECTING) {
      throw new DOMException("Still in CONNECTING state.", "InvalidStateError");
    }
    if (this.readyState !== _WailsSocket.OPEN) {
      return;
    }
    const immediate = toBytesSync(data);
    const snapshot = immediate ? Promise.resolve(immediate) : toBytes(data, this._sendAbort.signal);
    void snapshot.catch(() => {
    });
    const blobSize = !immediate && typeof Blob !== "undefined" && data instanceof Blob ? data.size : 0;
    if (blobSize) {
      this._buffered += blobSize;
    }
    if (immediate) {
      this._buffered += immediate.byteLength;
    }
    this._pending.push(snapshot);
    if (this._flushing) return;
    this._flushing = true;
    this._chain = this._chain.then(async () => {
      try {
        while (this._pending.length > 0) {
          const queued = this._pending.splice(0, MAX_BATCH_FRAMES);
          const batch = await Promise.all(queued);
          if (this.readyState === _WailsSocket.CLOSED) {
            break;
          }
          const total = batch.reduce((n, b) => n + b.byteLength, 0);
          try {
            await postBatch(this._id, batch, this._sendAbort.signal);
          } finally {
            this._buffered = Math.max(0, this._buffered - total);
          }
        }
      } finally {
        this._flushing = false;
      }
    }).catch((err) => {
      this._flushing = false;
      this._fail(err);
    });
  }
  /** Close the connection. Go's handler sees Receive fail and returns. */
  close(code = 1e3, reason = "") {
    if (this.readyState === _WailsSocket.CLOSING || this.readyState === _WailsSocket.CLOSED) {
      return;
    }
    this.readyState = _WailsSocket.CLOSING;
    this._pending.length = 0;
    this._buffered = 0;
    this._openAbort.abort();
    this._sendAbort.abort();
    this._chain = this._chain.then(
      () => postFrame(this._id, KIND_CLOSE, new Uint8Array(0))
    ).catch(() => {
    }).then(() => {
      this._closed(code, reason, true);
    });
  }
  /** @internal Called when the open frame comes back from Go. */
  _opened() {
    if (this.readyState !== _WailsSocket.CONNECTING) return;
    this.readyState = _WailsSocket.OPEN;
    this.dispatchEvent(new Event("open"));
  }
  /** @internal Called for each data frame addressed to this connection. */
  _message(payload) {
    if (this.readyState !== _WailsSocket.OPEN) return;
    let data;
    try {
      data = this._decode(payload);
    } catch (e) {
      this.dispatchEvent(new Event("error"));
      return;
    }
    if (data === payload && this.binaryType === "blob") {
      data = new Blob([payload]);
    }
    this.dispatchEvent(new MessageEvent("message", { data }));
  }
  /** @internal */
  _fail(err) {
    if (this.readyState === _WailsSocket.CLOSED || this.readyState === _WailsSocket.CLOSING) return;
    this.dispatchEvent(new Event("error"));
    this._closed(1006, err instanceof Error ? err.message : String(err), false);
  }
  /** @internal */
  _closed(code, reason, wasClean) {
    var _a3;
    if (this.readyState === _WailsSocket.CLOSED) return;
    this.readyState = _WailsSocket.CLOSED;
    this._openAbort.abort();
    this._sendAbort.abort();
    G.connections.delete(this._id);
    if (G.connections.size === 0) (_a3 = G.pollAbort) == null ? void 0 : _a3.abort();
    this._pending.length = 0;
    this._buffered = 0;
    const ev = typeof CloseEvent === "function" ? new CloseEvent("close", { code, reason, wasClean }) : Object.assign(new Event("close"), { code, reason, wasClean });
    this.dispatchEvent(ev);
  }
};
_WailsSocket.CONNECTING = 0;
_WailsSocket.OPEN = 1;
_WailsSocket.CLOSING = 2;
_WailsSocket.CLOSED = 3;
var WailsSocket = _WailsSocket;
function Stream(name) {
  var _a3;
  if (hasDOM) {
    const factory = (_a3 = window._wails) == null ? void 0 : _a3.streamFactory;
    if (typeof factory === "function") {
      return factory(name);
    }
  }
  return new WailsSocket(name);
}
function defineHandlerProperty(target, type) {
  let current = null;
  Object.defineProperty(target, "on" + type, {
    get: () => current,
    set(fn) {
      if (current) target.removeEventListener(type, current);
      current = typeof fn === "function" ? fn : null;
      if (current) target.addEventListener(type, current);
    },
    configurable: true,
    enumerable: true
  });
}
function JSONStream(name) {
  const socket = Stream(name);
  const decoder = new TextDecoder();
  const parse = (payload) => JSON.parse(typeof payload === "string" ? payload : decoder.decode(payload));
  if (socket instanceof WailsSocket) {
    socket._decode = (payload) => parse(payload);
  } else {
    const native = socket;
    native.binaryType = "arraybuffer";
    const wrapped = /* @__PURE__ */ new WeakMap();
    const decodedEvents = /* @__PURE__ */ new WeakMap();
    const add = native.addEventListener.bind(native);
    const remove = native.removeEventListener.bind(native);
    const decodeEvent = (ev) => {
      var _a3;
      if (decodedEvents.has(ev)) return (_a3 = decodedEvents.get(ev)) != null ? _a3 : null;
      try {
        const value = parse(ev.data);
        const decoded = new MessageEvent("message", { data: value });
        decodedEvents.set(ev, decoded);
        return decoded;
      } catch (e) {
        decodedEvents.set(ev, null);
        native.dispatchEvent(new Event("error"));
        return null;
      }
    };
    const decode = (listener) => {
      const existing = wrapped.get(listener);
      if (existing) return existing;
      const fn = (ev) => {
        const decoded = decodeEvent(ev);
        if (!decoded) return;
        if (typeof listener === "function") listener.call(native, decoded);
        else listener.handleEvent(decoded);
      };
      wrapped.set(listener, fn);
      return fn;
    };
    native.addEventListener = ((type, listener, opts) => {
      add(type, type === "message" && listener ? decode(listener) : listener, opts);
    });
    native.removeEventListener = ((type, listener, opts) => {
      const fn = type === "message" && listener ? wrapped.get(listener) : void 0;
      remove(type, fn != null ? fn : listener, opts);
    });
    let handler = null;
    Object.defineProperty(native, "onmessage", {
      get: () => handler,
      set(fn) {
        if (handler) native.removeEventListener("message", handler);
        handler = typeof fn === "function" ? fn : null;
        if (handler) native.addEventListener("message", handler);
      },
      configurable: true,
      enumerable: true
    });
  }
  const send = socket.send.bind(socket);
  socket.send = (value) => send(JSON.stringify(value));
  return socket;
}
function toBytesSync(data) {
  if (typeof data === "string") {
    return textEncoder.encode(data);
  }
  if (typeof Blob !== "undefined" && data instanceof Blob) {
    return null;
  }
  if (ArrayBuffer.isView(data)) {
    return new Uint8Array(data.buffer, data.byteOffset, data.byteLength).slice();
  }
  return new Uint8Array(data).slice();
}
function abortable(promise, signal) {
  if (!signal) return promise;
  if (signal.aborted) return Promise.reject(new StreamRequestCancelled());
  return new Promise((resolve, reject) => {
    const onAbort = () => reject(new StreamRequestCancelled());
    signal.addEventListener("abort", onAbort, { once: true });
    promise.then(
      (value) => {
        signal.removeEventListener("abort", onAbort);
        resolve(value);
      },
      (err) => {
        signal.removeEventListener("abort", onAbort);
        reject(err);
      }
    );
  });
}
async function toBytes(data, signal) {
  if (typeof data === "string") {
    return textEncoder.encode(data);
  }
  if (typeof Blob !== "undefined" && data instanceof Blob) {
    return new Uint8Array(await abortable(data.arrayBuffer(), signal));
  }
  if (ArrayBuffer.isView(data)) {
    return new Uint8Array(data.buffer, data.byteOffset, data.byteLength).slice();
  }
  return new Uint8Array(data).slice();
}
async function postFrame(connID, kind, body, name, signal) {
  const headers = {
    [HDR_SESSION]: G.session,
    [HDR_GENERATION]: String(G.generation),
    [HDR_CONN]: String(connID),
    [HDR_KIND]: String(kind),
    "Content-Type": "application/octet-stream"
  };
  if (name !== void 0) {
    headers[HDR_NAME] = name;
  }
  if (body.byteLength <= CHUNK_THRESHOLD2) {
    await postWithRetry(headers, body, signal);
    return;
  }
  const chunkID = nanoid();
  const total = Math.ceil(body.byteLength / CHUNK_THRESHOLD2);
  for (let i = 0; i < total; i++) {
    const slice = body.subarray(i * CHUNK_THRESHOLD2, (i + 1) * CHUNK_THRESHOLD2);
    await postWithRetry(__spreadProps(__spreadValues({}, headers), {
      [HDR_CHUNK]: chunkID,
      [HDR_CHUNK_INDEX]: String(i),
      [HDR_CHUNK_TOTAL]: String(total)
    }), slice, signal);
  }
}
async function postWithRetry(headers, body, signal) {
  let wait = 1;
  for (; ; ) {
    if (signal == null ? void 0 : signal.aborted) throw new StreamRequestCancelled();
    const resp = await fetch(streamURL("send"), {
      method: "POST",
      headers,
      body,
      signal
    });
    if (resp.ok) return;
    if (resp.status !== 429) {
      throw new Error(await resp.text());
    }
    await sleep(wait, signal);
    wait = Math.min(wait * 2, 50);
  }
}
async function postBatch(connID, frames, signal) {
  let start = 0;
  while (start < frames.length) {
    if (frames[start].byteLength > CHUNK_THRESHOLD2) {
      await postFrame(connID, KIND_DATA, frames[start], void 0, signal);
      start++;
      continue;
    }
    let end = start;
    let size = 4;
    while (end < frames.length && end - start < MAX_BATCH_FRAMES) {
      const frame = frames[end];
      if (frame.byteLength > CHUNK_THRESHOLD2 || size + 4 + frame.byteLength > CHUNK_THRESHOLD2) {
        break;
      }
      size += 4 + frame.byteLength;
      end++;
    }
    if (end === start) {
      await postFrame(connID, KIND_DATA, frames[start], void 0, signal);
      start++;
      continue;
    }
    const group = frames.slice(start, end);
    if (group.length === 1) {
      await postFrame(connID, KIND_DATA, group[0], void 0, signal);
    } else {
      await postBatchGroup(connID, group, signal);
    }
    start = end;
  }
}
async function postBatchGroup(connID, frames, signal) {
  var _a3;
  if (frames.length === 0 || frames.length > MAX_BATCH_FRAMES) {
    throw new Error("invalid stream batch size");
  }
  let body = buildBatch(frames);
  let sent = 0;
  let wait = 1;
  for (; ; ) {
    if (signal == null ? void 0 : signal.aborted) throw new StreamRequestCancelled();
    const resp = await fetch(streamURL("send"), {
      method: "POST",
      headers: {
        [HDR_SESSION]: G.session,
        [HDR_GENERATION]: String(G.generation),
        [HDR_CONN]: String(connID),
        [HDR_KIND]: String(KIND_DATA),
        [HDR_BATCH]: String(frames.length - sent),
        "Content-Type": "application/octet-stream"
      },
      body,
      signal
    });
    if (resp.ok) return;
    if (resp.status !== 429) {
      throw new Error(await resp.text());
    }
    const accepted = Number((_a3 = resp.headers.get(HDR_BATCH)) != null ? _a3 : 0);
    const remaining = frames.length - sent;
    if (!Number.isInteger(accepted) || accepted < 0 || accepted >= remaining) {
      if (accepted !== 0) throw new Error("invalid stream batch acknowledgement");
    }
    sent += accepted;
    body = buildBatch(frames.slice(sent));
    await sleep(wait, signal);
    wait = Math.min(wait * 2, 50);
  }
}
function buildBatch(frames) {
  let size = 4;
  for (const f of frames) size += 4 + f.byteLength;
  const out = new Uint8Array(size);
  const view = new DataView(out.buffer);
  view.setUint32(0, frames.length);
  let off = 4;
  for (const f of frames) {
    view.setUint32(off, f.byteLength);
    off += 4;
    out.set(f, off);
    off += f.byteLength;
  }
  return out;
}
function startPolling() {
  if (G.polling || !hasDOM) return;
  G.polling = true;
  void pollLoop();
}
async function pollLoop() {
  let backoff = 0;
  const abort = new AbortController();
  G.pollAbort = abort;
  while (G.connections.size > 0) {
    try {
      const resp = await fetch(streamURL("poll"), {
        method: "GET",
        headers: {
          [HDR_SESSION]: G.session,
          [HDR_GENERATION]: String(G.generation)
        },
        cache: "no-store",
        signal: abort.signal
      });
      if (resp.status === 410) {
        closeAll(1001, "session closed");
        break;
      }
      if (!resp.ok && resp.status !== 204) {
        if (!isRetryablePollStatus(resp.status)) {
          failAll("stream poll failed: " + resp.status);
          break;
        }
        throw new Error("retryable stream poll failure: " + resp.status);
      }
      backoff = 0;
      if (resp.status !== 204) {
        let buf;
        try {
          buf = await resp.arrayBuffer();
        } catch (e) {
          failAll("stream poll response body failed");
          break;
        }
        const frames = decodeFrames(buf);
        if (frames === null) {
          closeAll(1002, "unrecognised stream framing");
          break;
        }
        deliver(frames);
      }
    } catch (err) {
      if (abort.signal.aborted || G.connections.size === 0) break;
      backoff = backoff ? Math.min(backoff * 2, BACKOFF_MAX) : BACKOFF_MIN;
      try {
        await sleep(backoff, abort.signal);
      } catch (e) {
        break;
      }
    }
  }
  if (G.pollAbort === abort) G.pollAbort = void 0;
  G.polling = false;
  if (G.connections.size > 0) {
    startPolling();
  }
}
function isRetryablePollStatus(status) {
  return status === 408 || status === 425 || status === 429 || status >= 500 && status <= 599;
}
function decodeFrames(buf) {
  if (buf.byteLength < 9) return null;
  const dv = new DataView(buf);
  if (dv.getUint8(0) !== 87 || dv.getUint8(1) !== 83 || dv.getUint8(2) !== 49 || dv.getUint8(3) !== 0 || (dv.getUint8(4) & ~1) !== 0) {
    return null;
  }
  const count = dv.getUint32(5);
  const frames = [];
  let off = 9;
  for (let i = 0; i < count; i++) {
    if (off + 9 > buf.byteLength) return null;
    const connID = dv.getUint32(off);
    off += 4;
    const kind = dv.getUint8(off);
    off += 1;
    const len = dv.getUint32(off);
    off += 4;
    if (kind > KIND_ERROR || len > buf.byteLength - off) return null;
    const payload = buf.slice(off, off + len);
    off += len;
    frames.push({ connID, kind, payload });
  }
  return off === buf.byteLength ? frames : null;
}
function deliver(frames) {
  for (const frame of frames) {
    const connID = frame.connID;
    const kind = frame.kind;
    const payload = frame.payload;
    const conn = G.connections.get(connID);
    if (!conn) continue;
    switch (kind) {
      case KIND_DATA:
        conn._message(payload);
        break;
      case KIND_OPEN:
        conn._opened();
        break;
      case KIND_CLOSE:
        conn._closed(1e3, "", true);
        break;
      case KIND_ERROR:
        conn._fail(new Error(new TextDecoder().decode(payload)));
        break;
    }
  }
}
function closeAll(code, reason) {
  for (const conn of [...G.connections.values()]) {
    conn._closed(code, reason, false);
  }
}
function failAll(reason) {
  for (const conn of [...G.connections.values()]) {
    conn._fail(new Error(reason));
  }
}

// desktop/@wailsio/runtime/src/media.ts
var sourceKey = /* @__PURE__ */ Symbol.for("wails.media.source");
var defaultMaxBytes = 32 * 1024 * 1024;
function checkAborted(signal) {
  var _a3;
  if (signal.aborted) {
    throw (_a3 = signal.reason) != null ? _a3 : new DOMException("Media load cancelled", "AbortError");
  }
}
async function SetSource(element, streamName, name, options = {}) {
  var _a3, _b, _c, _d, _e;
  const maxBytes = (_a3 = options.maxBytes) != null ? _a3 : defaultMaxBytes;
  if (!Number.isSafeInteger(maxBytes) || maxBytes <= 0) {
    throw new RangeError("maxBytes must be a positive safe integer");
  }
  if (!streamName.trim() || !name.trim()) {
    throw new TypeError("A media stream and file name are required");
  }
  if (options.signal) checkAborted(options.signal);
  const managed = element;
  const state = (_b = managed[sourceKey]) != null ? _b : { generation: 0 };
  managed[sourceKey] = state;
  (_c = state.controller) == null ? void 0 : _c.abort();
  const generation = ++state.generation;
  const controller = new AbortController();
  state.controller = controller;
  const cancel = () => {
    var _a4;
    return controller.abort((_a4 = options.signal) == null ? void 0 : _a4.reason);
  };
  (_d = options.signal) == null ? void 0 : _d.addEventListener("abort", cancel, { once: true });
  try {
    const blob = await receiveBlob(streamName, name, maxBytes, controller.signal);
    checkAborted(controller.signal);
    if (state.generation !== generation) {
      throw new DOMException("Media source replaced", "AbortError");
    }
    const objectURL = URL.createObjectURL(blob);
    try {
      element.src = objectURL;
    } catch (error) {
      URL.revokeObjectURL(objectURL);
      throw error;
    }
    if (state.objectURL) URL.revokeObjectURL(state.objectURL);
    state.objectURL = objectURL;
  } finally {
    (_e = options.signal) == null ? void 0 : _e.removeEventListener("abort", cancel);
    if (state.generation === generation) state.controller = void 0;
  }
}
function ClearSource(element) {
  var _a3;
  const managed = element;
  const state = managed[sourceKey];
  if (!state) return;
  ++state.generation;
  (_a3 = state.controller) == null ? void 0 : _a3.abort();
  delete managed[sourceKey];
  if (state.objectURL) {
    try {
      element.pause();
      element.removeAttribute("src");
      element.load();
    } finally {
      URL.revokeObjectURL(state.objectURL);
    }
  }
}
function receiveBlob(streamName, name, maxBytes, signal) {
  return new Promise((resolve, reject) => {
    const stream = Stream(streamName);
    stream.binaryType = "arraybuffer";
    const parts = [];
    let expected;
    let contentType = "";
    let received = 0;
    let finished = false;
    const finish = (error) => {
      if (finished) return;
      finished = true;
      signal.removeEventListener("abort", cancel);
      stream.onopen = stream.onmessage = stream.onerror = stream.onclose = null;
      stream.close();
      if (error !== void 0) {
        parts.length = 0;
        reject(error);
      } else {
        try {
          resolve(new Blob(parts, { type: contentType }));
        } catch (cause) {
          reject(cause);
        } finally {
          parts.length = 0;
        }
      }
    };
    const cancel = () => {
      var _a3;
      return finish((_a3 = signal.reason) != null ? _a3 : new DOMException("Media load cancelled", "AbortError"));
    };
    signal.addEventListener("abort", cancel, { once: true });
    if (signal.aborted) {
      cancel();
      return;
    }
    stream.onopen = () => {
      if (finished) return;
      try {
        stream.send(JSON.stringify({ version: 1, name, maxBytes }));
      } catch (error) {
        finish(error);
      }
    };
    stream.onerror = () => finish(new Error("Media stream failed"));
    stream.onmessage = (event) => {
      if (finished) return;
      try {
        if (!(event.data instanceof ArrayBuffer)) throw new Error("Invalid media frame");
        if (expected === void 0) {
          if (event.data.byteLength > 4096) throw new Error("Invalid media metadata");
          const meta = JSON.parse(new TextDecoder().decode(event.data));
          if (meta.version !== 1) throw new Error("Unsupported media protocol");
          if (meta.error) throw meta.limit ? new RangeError(meta.error) : new Error(meta.error);
          if (!Number.isSafeInteger(meta.size) || meta.size < 0 || typeof meta.type !== "string") {
            throw new Error("Invalid media metadata");
          }
          if (meta.size > maxBytes) throw new RangeError("Media exceeds the ".concat(maxBytes, " byte limit"));
          expected = meta.size;
          contentType = meta.type;
        } else {
          received += event.data.byteLength;
          if (event.data.byteLength === 0 || event.data.byteLength > 64 * 1024 || received > expected || received > maxBytes) {
            throw new RangeError("Media stream exceeded its declared limits");
          }
          parts.push(event.data);
        }
      } catch (error) {
        finish(error);
      }
    };
    stream.onclose = (event) => {
      if (event.code !== 1e3 || expected === void 0 || received !== expected) {
        finish(new Error("Media stream closed before the complete file arrived"));
      } else {
        finish();
      }
    };
  });
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
((Haptics3) => {
  function Impact(style = "medium") {
    return call10(HapticsImpact, { style });
  }
  Haptics3.Impact = Impact;
})(Haptics || (Haptics = {}));
var Device;
((Device3) => {
  function Info2() {
    return call10(DeviceInfo);
  }
  Device3.Info = Info2;
})(Device || (Device = {}));

// desktop/@wailsio/runtime/src/android.ts
var android_exports = {};
__export(android_exports, {
  Device: () => Device2,
  Haptics: () => Haptics2,
  Toast: () => Toast
});
var call11 = newRuntimeCaller(objectNames.Android);
var HapticsVibrate = 0;
var DeviceInfo2 = 1;
var ToastShow = 2;
var Haptics2;
((Haptics3) => {
  function Vibrate(durationMs = 100) {
    return call11(HapticsVibrate, { duration: durationMs });
  }
  Haptics3.Vibrate = Vibrate;
})(Haptics2 || (Haptics2 = {}));
var Device2;
((Device3) => {
  function Info2() {
    return call11(DeviceInfo2);
  }
  Device3.Info = Info2;
})(Device2 || (Device2 = {}));
var Toast;
((Toast2) => {
  function Show2(message) {
    return call11(ToastShow, { message });
  }
  Toast2.Show = Show2;
})(Toast || (Toast = {}));

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
if (hasDOM) {
  window._wails = window._wails || {};
}
if (hasDOM) {
  window._wails.invoke = invoke;
  window._wails.clientId = clientId;
}
if (hasDOM) {
  window._wails.handlePlatformFileDrop = window_default.HandlePlatformFileDrop.bind(window_default);
}
if (hasDOM) {
  window._wails.handleDragEnter = handleDragEnter;
  window._wails.handleDragLeave = handleDragLeave;
  window._wails.handleDragOver = handleDragOver;
}
if (hasDOM) {
  invoke("wails:runtime:ready");
}
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
if (hasDOM) {
  loadOptionalScript("/wails/custom.js");
}
export {
  android_exports as Android,
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
  JSONStream,
  media_exports as Media,
  screens_exports as Screens,
  Stream,
  system_exports as System,
  updater_exports as Updater,
  wml_exports as WML,
  WailsSocket,
  window_default as Window,
  clientId,
  getTransport,
  loadOptionalScript,
  objectNames,
  setTransport
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZW52aXJvbm1lbnQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RpYWxvZ3MudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50cy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvbGlzdGVuZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NyZWF0ZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHByZWdpb24udHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbGFibGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbmNlbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL21lZGlhLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9zdHJlYW0udHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2lvcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYW5kcm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXBkYXRlci50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFNldHVwXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG59XG5cbmltcG9ydCBcIi4vY29udGV4dG1lbnUuanNcIjtcbmltcG9ydCBcIi4vZHJhZy5qc1wiO1xuaW1wb3J0IFwiLi9hcHByZWdpb24uanNcIjtcblxuLy8gUmUtZXhwb3J0IHB1YmxpYyBBUElcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gXCIuL2FwcGxpY2F0aW9uLmpzXCI7XG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCAqIGFzIENhbGwgZnJvbSBcIi4vY2FsbHMuanNcIjtcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmQuanNcIjtcbmltcG9ydCAqIGFzIENyZWF0ZSBmcm9tIFwiLi9jcmVhdGUuanNcIjtcbmltcG9ydCAqIGFzIERpYWxvZ3MgZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0ICogYXMgRmxhZ3MgZnJvbSBcIi4vZmxhZ3MuanNcIjtcbmltcG9ydCAqIGFzIE1lZGlhIGZyb20gXCIuL21lZGlhLmpzXCI7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnMuanNcIjtcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCAqIGFzIElPUyBmcm9tIFwiLi9pb3MuanNcIjtcbmltcG9ydCAqIGFzIEFuZHJvaWQgZnJvbSBcIi4vYW5kcm9pZC5qc1wiO1xuaW1wb3J0ICogYXMgVXBkYXRlciBmcm9tIFwiLi91cGRhdGVyLmpzXCI7XG5pbXBvcnQgV2luZG93LCB7IGhhbmRsZURyYWdFbnRlciwgaGFuZGxlRHJhZ0xlYXZlLCBoYW5kbGVEcmFnT3ZlciB9IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuaW1wb3J0ICogYXMgV01MIGZyb20gXCIuL3dtbC5qc1wiO1xuaW1wb3J0IHsgU3RyZWFtLCBKU09OU3RyZWFtLCBXYWlsc1NvY2tldCwgdHlwZSBKU09OU29ja2V0IH0gZnJvbSBcIi4vc3RyZWFtLmpzXCI7XG5cbmV4cG9ydCB7XG4gICAgQXBwbGljYXRpb24sXG4gICAgQnJvd3NlcixcbiAgICBDYWxsLFxuICAgIENsaXBib2FyZCxcbiAgICBEaWFsb2dzLFxuICAgIEV2ZW50cyxcbiAgICBGbGFncyxcbiAgICBNZWRpYSxcbiAgICBTY3JlZW5zLFxuICAgIFN5c3RlbSxcbiAgICBJT1MsXG4gICAgQW5kcm9pZCxcbiAgICBVcGRhdGVyLFxuICAgIFdpbmRvdyxcbiAgICBXTUwsXG4gICAgU3RyZWFtLFxuICAgIEpTT05TdHJlYW0sXG4gICAgV2FpbHNTb2NrZXQsXG4gICAgdHlwZSBKU09OU29ja2V0XG59O1xuXG4vKipcbiAqIEFuIGludGVybmFsIHV0aWxpdHkgY29uc3VtZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxuICpcbiAqIEBpZ25vcmVcbiAqL1xuZXhwb3J0IHsgQ3JlYXRlIH07XG5cbmV4cG9ydCAqIGZyb20gXCIuL2NhbmNlbGxhYmxlLmpzXCI7XG5cbi8vIEV4cG9ydCB0cmFuc3BvcnQgaW50ZXJmYWNlcyBhbmQgdXRpbGl0aWVzXG5leHBvcnQge1xuICAgIHNldFRyYW5zcG9ydCxcbiAgICBnZXRUcmFuc3BvcnQsXG4gICAgdHlwZSBSdW50aW1lVHJhbnNwb3J0LFxuICAgIG9iamVjdE5hbWVzLFxuICAgIGNsaWVudElkLFxufSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmltcG9ydCB7IGNsaWVudElkIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG4vLyBOb3RpZnkgYmFja2VuZFxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcbiAgICB3aW5kb3cuX3dhaWxzLmNsaWVudElkID0gY2xpZW50SWQ7XG59XG5cbi8vIFJlZ2lzdGVyIHBsYXRmb3JtIGhhbmRsZXJzIChpbnRlcm5hbCBBUEkpXG4vLyBOb3RlOiBXaW5kb3cgaXMgdGhlIHRoaXNXaW5kb3cgaW5zdGFuY2UgKGRlZmF1bHQgZXhwb3J0IGZyb20gd2luZG93LnRzKVxuLy8gQmluZGluZyBlbnN1cmVzICd0aGlzJyBjb3JyZWN0bHkgcmVmZXJzIHRvIHRoZSBjdXJyZW50IHdpbmRvdyBpbnN0YW5jZVxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5fd2FpbHMuaGFuZGxlUGxhdGZvcm1GaWxlRHJvcCA9IFdpbmRvdy5IYW5kbGVQbGF0Zm9ybUZpbGVEcm9wLmJpbmQoV2luZG93KTtcbn1cblxuLy8gTGludXgtc3BlY2lmaWMgZHJhZyBoYW5kbGVycyAoR1RLIGludGVyY2VwdHMgRE9NIGRyYWcgZXZlbnRzKVxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ0VudGVyID0gaGFuZGxlRHJhZ0VudGVyO1xuICAgIHdpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ0xlYXZlID0gaGFuZGxlRHJhZ0xlYXZlO1xuICAgIHdpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ092ZXIgPSBoYW5kbGVEcmFnT3Zlcjtcbn1cblxuaWYgKGhhc0RPTSkge1xuICAgIFN5c3RlbS5pbnZva2UoXCJ3YWlsczpydW50aW1lOnJlYWR5XCIpO1xufVxuXG4vKipcbiAqIExvYWRzIGEgc2NyaXB0IGZyb20gdGhlIGdpdmVuIFVSTCBpZiBpdCBleGlzdHMuXG4gKiBVc2VzIEhFQUQgcmVxdWVzdCB0byBjaGVjayBleGlzdGVuY2UsIHRoZW4gaW5qZWN0cyBhIHNjcmlwdCB0YWcuXG4gKiBTaWxlbnRseSBpZ25vcmVzIGlmIHRoZSBzY3JpcHQgZG9lc24ndCBleGlzdC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGxvYWRPcHRpb25hbFNjcmlwdCh1cmw6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBmZXRjaCh1cmwsIHsgbWV0aG9kOiAnSEVBRCcgfSlcbiAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XG4gICAgICAgICAgICAgICAgLy8gVmVyaWZ5IHRoZSByZXNwb25zZSBpcyBhY3R1YWxseSBKYXZhU2NyaXB0IGFuZCBub3QgYW4gSFRNTCBmYWxsYmFja1xuICAgICAgICAgICAgICAgIC8vIChlLmcuIFZpdGUgZGV2IHNlcnZlciByZXR1cm5zIGluZGV4Lmh0bWwgZm9yIHVua25vd24gcm91dGVzKVxuICAgICAgICAgICAgICAgIGNvbnN0IGNvbnRlbnRUeXBlID0gKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KCdjb250ZW50LXR5cGUnKSB8fCAnJykudG9Mb3dlckNhc2UoKTtcbiAgICAgICAgICAgICAgICBpZiAoY29udGVudFR5cGUuaW5jbHVkZXMoJ2phdmFzY3JpcHQnKSkge1xuICAgICAgICAgICAgICAgICAgICBjb25zdCBzY3JpcHQgPSBkb2N1bWVudC5jcmVhdGVFbGVtZW50KCdzY3JpcHQnKTtcbiAgICAgICAgICAgICAgICAgICAgc2NyaXB0LnNyYyA9IHVybDtcbiAgICAgICAgICAgICAgICAgICAgZG9jdW1lbnQuaGVhZC5hcHBlbmRDaGlsZChzY3JpcHQpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgfSlcbiAgICAgICAgLmNhdGNoKCgpID0+IHt9KTsgLy8gU2lsZW50bHkgaWdub3JlIC0gc2NyaXB0IGlzIG9wdGlvbmFsXG59XG5cbi8vIExvYWQgY3VzdG9tLmpzIGlmIGF2YWlsYWJsZSAodXNlZCBieSBzZXJ2ZXIgbW9kZSBmb3IgV2ViU29ja2V0IGV2ZW50cywgZXRjLilcbmlmIChoYXNET00pIHtcbiAgICBsb2FkT3B0aW9uYWxTY3JpcHQoJy93YWlscy9jdXN0b20uanMnKTtcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCB7IFF1ZXN0aW9uIH0gZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0IHsgRW1pdCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXG4gKiBAcGFyYW0gW2RhdGE9bnVsbF0gLSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cbiAqL1xuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XG4gICAgRW1pdChldmVudE5hbWUsIGRhdGEpO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgc3BlY2lmaWVkIHdpbmRvdy5cbiAqXG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKi9cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TmFtZTogc3RyaW5nLCBtZXRob2ROYW1lOiBzdHJpbmcpIHtcbiAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBXaW5kb3cuR2V0KHdpbmRvd05hbWUpO1xuICAgIGNvbnN0IG1ldGhvZCA9ICh0YXJnZXRXaW5kb3cgYXMgYW55KVttZXRob2ROYW1lXTtcblxuICAgIGlmICh0eXBlb2YgbWV0aG9kICE9PSBcImZ1bmN0aW9uXCIpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICB0cnkge1xuICAgICAgICBtZXRob2QuY2FsbCh0YXJnZXRXaW5kb3cpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JzogYCwgZSk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXQuXG4gKi9cbmZ1bmN0aW9uIG9uV01MVHJpZ2dlcmVkKGV2OiBFdmVudCk6IHZvaWQge1xuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0IGFzIEVsZW1lbnQ7XG5cbiAgICBmdW5jdGlvbiBydW5FZmZlY3QoY2hvaWNlID0gXCJZZXNcIikge1xuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgdGFyZ2V0V2luZG93ID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10YXJnZXQtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRhcmdldC13aW5kb3cnKSB8fCBcIlwiO1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC13aW5kb3cnKTtcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW51cmwnKTtcblxuICAgICAgICBpZiAoZXZlbnRUeXBlICE9PSBudWxsKVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIGlmICh3aW5kb3dNZXRob2QgIT09IG51bGwpXG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHRhcmdldFdpbmRvdywgd2luZG93TWV0aG9kKTtcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcbiAgICAgICAgICAgIHZvaWQgT3BlblVSTCh1cmwpO1xuICAgIH1cblxuICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuXG4gICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgUXVlc3Rpb24oe1xuICAgICAgICAgICAgVGl0bGU6IFwiQ29uZmlybVwiLFxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcbiAgICAgICAgICAgIERldGFjaGVkOiBmYWxzZSxcbiAgICAgICAgICAgIEJ1dHRvbnM6IFtcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJOb1wiLCBJc0RlZmF1bHQ6IHRydWUgfVxuICAgICAgICAgICAgXVxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcnVuRWZmZWN0KCk7XG4gICAgfVxufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY29udHJvbGxlclN5bSA9IFN5bWJvbChcImNvbnRyb2xsZXJcIik7XG5jb25zdCB0cmlnZ2VyTWFwU3ltID0gU3ltYm9sKFwidHJpZ2dlck1hcFwiKTtcbmNvbnN0IGVsZW1lbnRDb3VudFN5bSA9IFN5bWJvbChcImVsZW1lbnRDb3VudFwiKTtcblxuLyoqXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXG4gKiBpdCB0aWVzIHRoZW0gdG8gYW4gQWJvcnRTaWduYWwgYW5kIHVzZXMgYW4gQWJvcnRDb250cm9sbGVyIHRvIHJlbW92ZSB0aGVtIGFsbCBhdCBvbmNlLlxuICovXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgW2NvbnRyb2xsZXJTeW1dOiBBYm9ydENvbnRyb2xsZXI7XG5cbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIGFuIG9wdGlvbnMgb2JqZWN0IGZvciBhZGRFdmVudExpc3RlbmVyIHRoYXQgdGllcyB0aGUgbGlzdGVuZXJcbiAgICAgKiB0byB0aGUgQWJvcnRTaWduYWwgZnJvbSB0aGUgY3VycmVudCBBYm9ydENvbnRyb2xsZXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxuICAgICAqIEBwYXJhbSB0cmlnZ2VycyAtIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudHNcbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xuICAgICAgICByZXR1cm4geyBzaWduYWw6IHRoaXNbY29udHJvbGxlclN5bV0uc2lnbmFsIH07XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMgYW5kIHJlc2V0cyB0aGUgcmVnaXN0cnkuXG4gICAgICovXG4gICAgcmVzZXQoKTogdm9pZCB7XG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0uYWJvcnQoKTtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG59XG5cbi8qKlxuICogV2Vha01hcFJlZ2lzdHJ5IG1hcHMgYWN0aXZlIHRyaWdnZXIgZXZlbnRzIHRvIGVhY2ggRE9NIGVsZW1lbnQgdGhyb3VnaCBhIFdlYWtNYXAuXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcbiAqIGNvbGxlY3Rpb24gb2YgdGhlIGludm9sdmVkIGVsZW1lbnRzLlxuICovXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xuICAgIC8qKiBTdG9yZXMgdGhlIGN1cnJlbnQgZWxlbWVudC10by10cmlnZ2VyIG1hcHBpbmcuICovXG4gICAgW3RyaWdnZXJNYXBTeW1dOiBXZWFrTWFwPEVsZW1lbnQsIHN0cmluZ1tdPjtcbiAgICAvKiogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLiAqL1xuICAgIFtlbGVtZW50Q291bnRTeW1dOiBudW1iZXI7XG5cbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50U3ltXSA9IDA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XG4gICAgICovXG4gICAgc2V0KGVsZW1lbnQ6IEVsZW1lbnQsIHRyaWdnZXJzOiBzdHJpbmdbXSk6IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zIHtcbiAgICAgICAgaWYgKCF0aGlzW3RyaWdnZXJNYXBTeW1dLmhhcyhlbGVtZW50KSkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0rKzsgfVxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dLnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgICAgIHJldHVybiB7fTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cbiAgICAgKi9cbiAgICByZXNldCgpOiB2b2lkIHtcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xuICAgICAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxuICAgICAgICAgICAgICAgIGJyZWFrO1xuXG4gICAgICAgICAgICBjb25zdCB0cmlnZ2VycyA9IHRoaXNbdHJpZ2dlck1hcFN5bV0uZ2V0KGVsZW1lbnQpO1xuICAgICAgICAgICAgaWYgKHRyaWdnZXJzICE9IG51bGwpIHsgdGhpc1tlbGVtZW50Q291bnRTeW1dLS07IH1cblxuICAgICAgICAgICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzIHx8IFtdKVxuICAgICAgICAgICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCk7XG4gICAgICAgIH1cblxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcbiAgICB9XG59XG5cbmNvbnN0IHRyaWdnZXJSZWdpc3RyeSA9IGNhbkFib3J0TGlzdGVuZXJzKCkgPyBuZXcgQWJvcnRDb250cm9sbGVyUmVnaXN0cnkoKSA6IG5ldyBXZWFrTWFwUmVnaXN0cnkoKTtcblxuLyoqXG4gKiBBZGRzIGV2ZW50IGxpc3RlbmVycyB0byB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXG4gKi9cbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50OiBFbGVtZW50KTogdm9pZCB7XG4gICAgY29uc3QgdHJpZ2dlclJlZ0V4cCA9IC9cXFMrL2c7XG4gICAgY29uc3QgdHJpZ2dlckF0dHIgPSAoZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCIpO1xuICAgIGNvbnN0IHRyaWdnZXJzOiBzdHJpbmdbXSA9IFtdO1xuXG4gICAgbGV0IG1hdGNoO1xuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcbiAgICAgICAgdHJpZ2dlcnMucHVzaChtYXRjaFswXSk7XG5cbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzKVxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW5hYmxlKCk6IHZvaWQge1xuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xufVxuXG4vKipcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpOiB2b2lkIHtcbiAgICB0cmlnZ2VyUmVnaXN0cnkucmVzZXQoKTtcbiAgICBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtZXZlbnRdLCBbd21sLXdpbmRvd10sIFt3bWwtb3BlbnVybF0sIFtkYXRhLXdtbC1ldmVudF0sIFtkYXRhLXdtbC13aW5kb3ddLCBbZGF0YS13bWwtb3BlbnVybF0nKS5mb3JFYWNoKGFkZFdNTExpc3RlbmVycyk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5Ccm93c2VyKTtcblxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xuXG4vKipcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMLlxuICpcbiAqIEBwYXJhbSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5VUkwodXJsOiBzdHJpbmcgfCBVUkwpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChCcm93c2VyT3BlblVSTCwge3VybDogdXJsLnRvU3RyaW5nKCl9KTtcbn1cbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9haS9uYW5vaWRcblxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXG4vL1xuLy8gQ29weXJpZ2h0IDIwMTcgQW5kcmV5IFNpdG5payA8YW5kcmV5QHNpdG5pay5ydT5cbi8vXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mXG4vLyB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsIGluXG4vLyB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvXG4vLyB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsIGNvcGllcyBvZlxuLy8gdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMgZnVybmlzaGVkIHRvIGRvIHNvLFxuLy8gICAgIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxuLy9cbi8vICAgICBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXG4vL1xuLy8gICAgIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1Jcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTXG4vLyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1Jcbi8vIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSIExJQUJJTElUWSwgV0hFVEhFUlxuLy8gSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU5cbi8vIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuXG5cbi8vIFRoaXMgYWxwaGFiZXQgdXNlcyBgQS1aYS16MC05Xy1gIHN5bWJvbHMuXG4vLyBUaGUgb3JkZXIgb2YgY2hhcmFjdGVycyBpcyBvcHRpbWl6ZWQgZm9yIGJldHRlciBnemlwIGFuZCBicm90bGkgY29tcHJlc3Npb24uXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBzYW1lIGZpbGUgKHdvcmtzIGJvdGggZm9yIGd6aXAgYW5kIGJyb3RsaSk6XG4vLyBgJ3VzZWAsIGBhbmRvbWAsIGFuZCBgcmljdCdgXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBicm90bGkgZGVmYXVsdCBkaWN0aW9uYXJ5OlxuLy8gYC0yNlRgLCBgMTk4M2AsIGA0MHB4YCwgYDc1cHhgLCBgYnVzaGAsIGBqYWNrYCwgYG1pbmRgLCBgdmVyeWAsIGFuZCBgd29sZmBcbmNvbnN0IHVybEFscGhhYmV0ID1cbiAgICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcblxuLy8gTU9ESUZJRUQgRk9SIFdBSUxTOiB0aGUgdXBzdHJlYW0gaW1wbGVtZW50YXRpb24gZHJhd3MgZnJvbSBNYXRoLnJhbmRvbSgpLFxuLy8gd2hpY2ggaXMgbm90IGEgY3J5cHRvZ3JhcGhpYyBzb3VyY2UuIFRoZXNlIGlkcyBhcmUgdXNlZCBmb3IgdGhlIHJ1bnRpbWVcbi8vIGNsaWVudCBpZCwgYmluZGluZyBjYWxsIGlkcywgY2h1bmsgaWRzLCBhbmQgZGVza3RvcCBzdHJlYW0gc2Vzc2lvbiBpZHMuIFRoZXlcbi8vIGFyZSBjb3JyZWxhdGlvbiBpZGVudGlmaWVycyByYXRoZXIgdGhhbiBhdXRoZW50aWNhdGlvbiBjYXBhYmlsaXRpZXMsIGJ1dCBhXG4vLyBzdHJvbmdlciBzb3VyY2UgbWFrZXMgYWNjaWRlbnRhbCBjb2xsaXNpb25zIG5lZ2xpZ2libGUgZXZlbiBhY3Jvc3MgbWFueVxuLy8gY29uY3VycmVudGx5IGxvYWRlZCBydW50aW1lcy4gVGhlIGFscGhhYmV0IGlzIDY0IGNoYXJhY3RlcnMsIHNvIG1hc2tpbmcgYVxuLy8gcmFuZG9tIGJ5dGUgd2l0aCA2MyBpcyB1bmlmb3JtIFx1MjAxNCBubyByZWplY3Rpb24gc2FtcGxpbmcgbmVlZGVkLlxuZXhwb3J0IGZ1bmN0aW9uIG5hbm9pZChzaXplOiBudW1iZXIgPSAyMSk6IHN0cmluZyB7XG4gICAgY29uc3Qgc291cmNlID0gdHlwZW9mIGNyeXB0byAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgY3J5cHRvLmdldFJhbmRvbVZhbHVlcyA9PT0gXCJmdW5jdGlvblwiXG4gICAgICAgID8gY3J5cHRvXG4gICAgICAgIDogbnVsbDtcblxuICAgIGlmIChzb3VyY2UpIHtcbiAgICAgICAgY29uc3QgYnl0ZXMgPSBuZXcgVWludDhBcnJheShzaXplKTtcbiAgICAgICAgc291cmNlLmdldFJhbmRvbVZhbHVlcyhieXRlcyk7XG4gICAgICAgIGxldCBpZCA9ICcnO1xuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNpemU7IGkrKykge1xuICAgICAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbYnl0ZXNbaV0gJiA2M107XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIGlkO1xuICAgIH1cblxuICAgIC8vIEZhbGxiYWNrIGZvciBlbnZpcm9ubWVudHMgd2l0aG91dCBXZWIgQ3J5cHRvICh2ZXJ5IG9sZCBlbmdpbmVzLCBzb21lXG4gICAgLy8gc2VydmVyLXNpZGUgcmVuZGVyaW5nIGNvbnRleHRzKS4gSWRzIHJlbWFpbiB1bmlxdWUgYnV0IG5vdCB1bmd1ZXNzYWJsZS5cbiAgICBsZXQgaWQgPSAnJztcbiAgICBsZXQgaSA9IHNpemUgfCAwO1xuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXTtcbiAgICB9XG4gICAgcmV0dXJuIGlkO1xufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIFRydWUgd2hlbiBydW5uaW5nIGluc2lkZSBhIGJyb3dzZXIvd2VidmlldyB3aXRoIGEgRE9NIGF2YWlsYWJsZS5cbiAqIEZhbHNlIHVuZGVyIHNlcnZlci1zaWRlIHJlbmRlcmluZyAoZS5nLiBgbmV4dCBidWlsZGAgcHJlcmVuZGVyaW5nKSxcbiAqIHdoZXJlIGFwcGxpY2F0aW9uIGNvZGUgbWF5IGltcG9ydCB0aGUgcnVudGltZSBtb2R1bGUgZXZlbiB0aG91Z2ggbm9cbiAqIFdhaWxzIEFQSXMgY2FuIGFjdHVhbGx5IGJlIHVzZWQgKCM0Njc5KS4gTW9kdWxlcyBtdXN0IG5vdCB0b3VjaFxuICogYHdpbmRvd2AvYGRvY3VtZW50YCBhdCBpbXBvcnQgdGltZSBleGNlcHQgYmVoaW5kIHRoaXMgZ3VhcmQuXG4gKi9cbmV4cG9ydCBjb25zdCBoYXNET00gPSB0eXBlb2Ygd2luZG93ICE9PSBcInVuZGVmaW5lZFwiICYmIHR5cGVvZiBkb2N1bWVudCAhPT0gXCJ1bmRlZmluZWRcIjtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG4vLyBSZXNvbHZlZCBsYXppbHk6IHdpbmRvdyBkb2VzIG5vdCBleGlzdCB3aGVuIHRoZSBtb2R1bGUgaXMgaW1wb3J0ZWQgZHVyaW5nXG4vLyBzZXJ2ZXItc2lkZSByZW5kZXJpbmcgKCM0Njc5KSwgYW5kIG5vdGhpbmcgY2FuIGNhbGwgdGhlIHJ1bnRpbWUgdGhlcmUuXG5mdW5jdGlvbiBydW50aW1lVVJMKCk6IHN0cmluZyB7XG4gICAgcmV0dXJuIHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XG59XG5cbi8vIFN0YXkgdW5kZXIgV2ViVmlldzIncyB+Mk1CIHJlcXVlc3QgYm9keSBidWZmZXJpbmcgbGltaXQgaW4gV2ViUmVzb3VyY2VSZXF1ZXN0ZWQuXG5jb25zdCBDSFVOS19USFJFU0hPTEQgPSA1MTIgKiAxMDI0O1xuXG4vLyBSZS1leHBvcnQgbmFub2lkIGZvciBjdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uc1xuZXhwb3J0IHsgbmFub2lkIH07XG5cbnR5cGUgQ2FsbEVycm9yVHlwZSA9IHtcbiAgICBtZXNzYWdlOiBzdHJpbmcsXG4gICAgY2F1c2U/OiB1bmtub3duLFxuICAgIGtpbmQ6IFwiUmVmZXJlbmNlRXJyb3JcIiB8IFwiVHlwZUVycm9yXCIgfCBcIlJ1bnRpbWVFcnJvclwiXG59XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB0aHJvd24gaW4gY2FzZSB0aGUgYm91bmQgbWV0aG9kIHJldHVybnMgYW4gZXJyb3IuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBSdW50aW1lRXJyb3IjbmFtZX0gcHJvcGVydHkgaXMgXCJSdW50aW1lRXJyb3JcIi5cbiAqL1xuZXhwb3J0IGNsYXNzIFJ1bnRpbWVFcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IFJ1bnRpbWVFcnJvciBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiUnVudGltZUVycm9yXCI7XG4gICAgfVxufVxuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xuICAgIENhbGw6IDAsXG4gICAgQ2xpcGJvYXJkOiAxLFxuICAgIEFwcGxpY2F0aW9uOiAyLFxuICAgIEV2ZW50czogMyxcbiAgICBDb250ZXh0TWVudTogNCxcbiAgICBEaWFsb2c6IDUsXG4gICAgV2luZG93OiA2LFxuICAgIFNjcmVlbnM6IDcsXG4gICAgU3lzdGVtOiA4LFxuICAgIEJyb3dzZXI6IDksXG4gICAgQ2FuY2VsQ2FsbDogMTAsXG4gICAgSU9TOiAxMSxcbiAgICBBbmRyb2lkOiAxMixcbn0pO1xuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG4vKipcbiAqIFJ1bnRpbWVUcmFuc3BvcnQgZGVmaW5lcyB0aGUgaW50ZXJmYWNlIGZvciBjdXN0b20gSVBDIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnMuXG4gKiBJbXBsZW1lbnQgdGhpcyBpbnRlcmZhY2UgdG8gdXNlIFdlYlNvY2tldHMsIGN1c3RvbSBwcm90b2NvbHMsIG9yIGFueSBvdGhlclxuICogdHJhbnNwb3J0IG1lY2hhbmlzbSBpbnN0ZWFkIG9mIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2guXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgUnVudGltZVRyYW5zcG9ydCB7XG4gICAgLyoqXG4gICAgICogU2VuZCBhIHJ1bnRpbWUgY2FsbCBhbmQgcmV0dXJuIHRoZSByZXNwb25zZS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBvYmplY3RJRCAtIFRoZSBXYWlscyBvYmplY3QgSUQgKDA9Q2FsbCwgMT1DbGlwYm9hcmQsIGV0Yy4pXG4gICAgICogQHBhcmFtIG1ldGhvZCAtIFRoZSBtZXRob2QgSUQgdG8gY2FsbFxuICAgICAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gT3B0aW9uYWwgd2luZG93IG5hbWVcbiAgICAgKiBAcGFyYW0gYXJncyAtIEFyZ3VtZW50cyB0byBwYXNzICh3aWxsIGJlIEpTT04gc3RyaW5naWZpZWQgaWYgcHJlc2VudClcbiAgICAgKiBAcmV0dXJucyBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzcG9uc2UgZGF0YVxuICAgICAqL1xuICAgIGNhbGwob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+O1xufVxuXG4vKipcbiAqIEN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb24gKGNhbiBiZSBzZXQgYnkgdXNlcilcbiAqL1xubGV0IGN1c3RvbVRyYW5zcG9ydDogUnVudGltZVRyYW5zcG9ydCB8IG51bGwgPSBudWxsO1xuXG4vKipcbiAqIFNldCBhIGN1c3RvbSB0cmFuc3BvcnQgZm9yIGFsbCBXYWlscyBydW50aW1lIGNhbGxzLlxuICogVGhpcyBhbGxvd3MgeW91IHRvIHJlcGxhY2UgdGhlIGRlZmF1bHQgSFRUUCBmZXRjaCB0cmFuc3BvcnQgd2l0aFxuICogV2ViU29ja2V0cywgY3VzdG9tIHByb3RvY29scywgb3IgYW55IG90aGVyIG1lY2hhbmlzbS5cbiAqXG4gKiBAcGFyYW0gdHJhbnNwb3J0IC0gWW91ciBjdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uXG4gKlxuICogQGV4YW1wbGVcbiAqIGBgYHR5cGVzY3JpcHRcbiAqIGltcG9ydCB7IHNldFRyYW5zcG9ydCB9IGZyb20gJy93YWlscy9ydW50aW1lLmpzJztcbiAqXG4gKiBjb25zdCB3c1RyYW5zcG9ydCA9IHtcbiAqICAgY2FsbDogYXN5bmMgKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpID0+IHtcbiAqICAgICAvLyBZb3VyIFdlYlNvY2tldCBpbXBsZW1lbnRhdGlvblxuICogICB9XG4gKiB9O1xuICpcbiAqIHNldFRyYW5zcG9ydCh3c1RyYW5zcG9ydCk7XG4gKiBgYGBcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHNldFRyYW5zcG9ydCh0cmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsKTogdm9pZCB7XG4gICAgY3VzdG9tVHJhbnNwb3J0ID0gdHJhbnNwb3J0O1xufVxuXG4vKipcbiAqIEdldCB0aGUgY3VycmVudCB0cmFuc3BvcnQgKHVzZWZ1bCBmb3IgZXh0ZW5kaW5nL3dyYXBwaW5nKVxuICovXG5leHBvcnQgZnVuY3Rpb24gZ2V0VHJhbnNwb3J0KCk6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsIHtcbiAgICByZXR1cm4gY3VzdG9tVHJhbnNwb3J0O1xufVxuXG4vKipcbiAqIENyZWF0ZXMgYSBuZXcgcnVudGltZSBjYWxsZXIgd2l0aCBzcGVjaWZpZWQgSUQuXG4gKlxuICogQHBhcmFtIG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gKiBAcmV0dXJuIFRoZSBuZXcgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcgPSAnJykge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kOiBudW1iZXIsIGFyZ3M6IGFueSA9IG51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG5hc3luYyBmdW5jdGlvbiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3RJRDogbnVtYmVyLCBtZXRob2Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nLCBhcmdzOiBhbnkpOiBQcm9taXNlPGFueT4ge1xuICAgIC8vIFVzZSBjdXN0b20gdHJhbnNwb3J0IGlmIGF2YWlsYWJsZVxuICAgIGlmIChjdXN0b21UcmFuc3BvcnQpIHtcbiAgICAgICAgcmV0dXJuIGN1c3RvbVRyYW5zcG9ydC5jYWxsKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xuICAgIH1cblxuICAgIC8vIERlZmF1bHQgSFRUUCBmZXRjaCB0cmFuc3BvcnRcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKCkpO1xuXG4gICAgbGV0IGJvZHk6IHsgb2JqZWN0OiBudW1iZXI7IG1ldGhvZDogbnVtYmVyLCBhcmdzPzogYW55IH0gPSB7XG4gICAgICBvYmplY3Q6IG9iamVjdElELFxuICAgICAgbWV0aG9kXG4gICAgfVxuICAgIGlmIChhcmdzICE9PSBudWxsICYmIGFyZ3MgIT09IHVuZGVmaW5lZCkge1xuICAgICAgYm9keS5hcmdzID0gYXJncztcbiAgICB9XG5cbiAgICBsZXQgaGVhZGVyczogUmVjb3JkPHN0cmluZywgc3RyaW5nPiA9IHtcbiAgICAgICAgW1wieC13YWlscy1jbGllbnQtaWRcIl06IGNsaWVudElkLFxuICAgICAgICBbXCJDb250ZW50LVR5cGVcIl06IFwiYXBwbGljYXRpb24vanNvblwiXG4gICAgfVxuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGhlYWRlcnNbXCJ4LXdhaWxzLXdpbmRvdy1uYW1lXCJdID0gd2luZG93TmFtZTtcbiAgICB9XG5cbiAgICBjb25zdCBib2R5U3RyID0gSlNPTi5zdHJpbmdpZnkoYm9keSk7XG4gICAgbGV0IHJlc3BvbnNlOiBSZXNwb25zZTtcbiAgICBpZiAoYm9keVN0ci5sZW5ndGggPiBDSFVOS19USFJFU0hPTEQpIHtcbiAgICAgICAgcmVzcG9uc2UgPSBhd2FpdCBzZW5kQ2h1bmtlZCh1cmwsIGhlYWRlcnMsIGJvZHlTdHIpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJlc3BvbnNlID0gYXdhaXQgZmV0Y2godXJsLCB7IG1ldGhvZDogJ1BPU1QnLCBoZWFkZXJzLCBib2R5OiBib2R5U3RyIH0pO1xuICAgIH1cbiAgICBpZiAoIXJlc3BvbnNlLm9rKSB7XG4gICAgICBjb25zdCBjdCA9IHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpO1xuICAgICAgaWYgKGN0Py5pbmNsdWRlcyhcImFwcGxpY2F0aW9uL2pzb25cIikpIHtcbiAgICAgICAgICBjb25zdCBqc29uOiBDYWxsRXJyb3JUeXBlID0gYXdhaXQgcmVzcG9uc2UuanNvbigpO1xuICAgICAgICAgIGxldCBlcnI7XG4gICAgICAgICAgc3dpdGNoIChqc29uLmtpbmQpIHtcbiAgICAgICAgICAgICAgY2FzZSBcIlJlZmVyZW5jZUVycm9yXCI6IGVyciA9IG5ldyBSZWZlcmVuY2VFcnJvcihqc29uLm1lc3NhZ2UpOyBicmVhaztcbiAgICAgICAgICAgICAgY2FzZSBcIlR5cGVFcnJvclwiOiAgICAgIGVyciA9IG5ldyBUeXBlRXJyb3IoanNvbi5tZXNzYWdlKTsgYnJlYWs7XG4gICAgICAgICAgICAgIGNhc2UgXCJSdW50aW1lRXJyb3JcIjogICBlcnIgPSBuZXcgUnVudGltZUVycm9yKGpzb24ubWVzc2FnZSk7IGJyZWFrO1xuICAgICAgICAgICAgICBkZWZhdWx0OiAgICAgICAgICAgICAgIGVyciA9IG5ldyBFcnJvcihqc29uLm1lc3NhZ2UpO1xuICAgICAgICAgIH1cbiAgICAgICAgICBlcnIuY2F1c2UgPSBqc29uLmNhdXNlO1xuICAgICAgICAgIHRocm93IGVyclxuICAgICAgfVxuICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XG4gICAgfVxuXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgfVxufVxuXG4vLyBzZW5kQ2h1bmtlZCBzcGxpdHMgYSBsYXJnZSBzZXJpYWxpc2VkIHJlcXVlc3QgYm9keSBpbnRvIENIVU5LX1RIUkVTSE9MRC1zaXplZFxuLy8gYnl0ZSBjaHVua3MgYW5kIHNlbmRzIHRoZW0gc2VyaWFsbHkuICBFbmNvZGluZyB0byBVVEYtOCBieXRlcyBiZWZvcmUgc2xpY2luZ1xuLy8gcHJldmVudHMgY29ycnVwdGlvbiBvZiBub24tQk1QIGNoYXJhY3RlcnMgKHN1cnJvZ2F0ZSBwYWlycykgdGhhdCB3b3VsZCBvY2N1clxuLy8gd2hlbiBzcGxpdHRpbmcgYXQgSmF2YVNjcmlwdCBzdHJpbmcgaW5kaWNlcy4gIFRoZSBHbyB0cmFuc3BvcnQgYXNzZW1ibGVzIHRoZVxuLy8gcmF3IGJ5dGVzIGJlZm9yZSBwcm9jZXNzaW5nLiAgT25seSB0aGUgZmluYWwgY2h1bmsncyByZXNwb25zZSBjYXJyaWVzIHRoZSBSUEMgcmVzdWx0LlxuYXN5bmMgZnVuY3Rpb24gc2VuZENodW5rZWQodXJsOiBVUkwsIGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4sIGJvZHlTdHI6IHN0cmluZyk6IFByb21pc2U8UmVzcG9uc2U+IHtcbiAgICBjb25zdCBjaHVua0lkID0gbmFub2lkKCk7XG4gICAgY29uc3QgYm9keUJ5dGVzID0gbmV3IFRleHRFbmNvZGVyKCkuZW5jb2RlKGJvZHlTdHIpO1xuICAgIGNvbnN0IHRvdGFsQ2h1bmtzID0gTWF0aC5jZWlsKGJvZHlCeXRlcy5sZW5ndGggLyBDSFVOS19USFJFU0hPTEQpO1xuXG4gICAgZm9yIChsZXQgaSA9IDA7IGkgPCB0b3RhbENodW5rcyAtIDE7IGkrKykge1xuICAgICAgICBjb25zdCBjaHVuayA9IGJvZHlCeXRlcy5zdWJhcnJheShpICogQ0hVTktfVEhSRVNIT0xELCAoaSArIDEpICogQ0hVTktfVEhSRVNIT0xEKTtcbiAgICAgICAgY29uc3QgcmVzcCA9IGF3YWl0IGZldGNoKHVybCwge1xuICAgICAgICAgICAgbWV0aG9kOiAnUE9TVCcsXG4gICAgICAgICAgICBoZWFkZXJzOiB7XG4gICAgICAgICAgICAgICAgLi4uaGVhZGVycyxcbiAgICAgICAgICAgICAgICAneC13YWlscy1jaHVuay1pZCc6IGNodW5rSWQsXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaW5kZXgnOiBTdHJpbmcoaSksXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstdG90YWwnOiBTdHJpbmcodG90YWxDaHVua3MpLFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGJvZHk6IGNodW5rLFxuICAgICAgICB9KTtcbiAgICAgICAgaWYgKCFyZXNwLm9rKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIGZldGNoKHVybCwge1xuICAgICAgICBtZXRob2Q6ICdQT1NUJyxcbiAgICAgICAgaGVhZGVyczoge1xuICAgICAgICAgICAgLi4uaGVhZGVycyxcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWlkJzogY2h1bmtJZCxcbiAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWluZGV4JzogU3RyaW5nKHRvdGFsQ2h1bmtzIC0gMSksXG4gICAgICAgICAgICAneC13YWlscy1jaHVuay10b3RhbCc6IFN0cmluZyh0b3RhbENodW5rcyksXG4gICAgICAgIH0sXG4gICAgICAgIGJvZHk6IGJvZHlCeXRlcy5zdWJhcnJheSgodG90YWxDaHVua3MgLSAxKSAqIENIVU5LX1RIUkVTSE9MRCksXG4gICAgfSk7XG59XG5cbi8qKlxuICogQW5kcm9pZCBXZWJWaWV3IGNhbm5vdCBkZWxpdmVyIGZldGNoKCkgUE9TVCBib2RpZXMgdG9cbiAqIHNob3VsZEludGVyY2VwdFJlcXVlc3QsIHNvIHRoZSBkZWZhdWx0IEhUVFAgdHJhbnNwb3J0IGNhbm5vdCByZWFjaCBHby5cbiAqIFdoZW4gdGhlIEFuZHJvaWQgSmF2YXNjcmlwdEludGVyZmFjZSBicmlkZ2UgKHdpbmRvdy53YWlscykgaXMgcHJlc2VudCxcbiAqIHJvdXRlIHJ1bnRpbWUgY2FsbHMgdGhyb3VnaCBpdCBpbnN0ZWFkLiBSZXNwb25zZXMgYXJyaXZlIHZpYVxuICogd2luZG93Ll93YWlsc0FuZHJvaWRDYWxsYmFjaywgaW52b2tlZCBieSB0aGUgSmF2YSBzaWRlLlxuICovXG5pbnRlcmZhY2UgQW5kcm9pZEpTQnJpZGdlIHtcbiAgICBpbnZva2VBc3luYyhjYWxsYmFja0lEOiBzdHJpbmcsIHBheWxvYWQ6IHN0cmluZyk6IHZvaWQ7XG59XG5cbmNvbnN0IGFuZHJvaWRCcmlkZ2U6IEFuZHJvaWRKU0JyaWRnZSB8IG51bGwgPSBoYXNET00gJiZcbiAgICB0eXBlb2YgKHdpbmRvdyBhcyBhbnkpLndhaWxzPy5pbnZva2VBc3luYyA9PT0gXCJmdW5jdGlvblwiID8gKHdpbmRvdyBhcyBhbnkpLndhaWxzIDogbnVsbDtcblxuaWYgKGFuZHJvaWRCcmlkZ2UpIHtcbiAgICBjb25zdCBwZW5kaW5nID0gbmV3IE1hcDxzdHJpbmcsIHsgcmVzb2x2ZTogKHZhbHVlOiBhbnkpID0+IHZvaWQ7IHJlamVjdDogKHJlYXNvbjogYW55KSA9PiB2b2lkIH0+KCk7XG5cbiAgICAod2luZG93IGFzIGFueSkuX3dhaWxzQW5kcm9pZENhbGxiYWNrID0gKGlkOiBzdHJpbmcsIHJlc3BvbnNlOiBzdHJpbmcgfCBudWxsLCBlcnJvcjogc3RyaW5nIHwgbnVsbCkgPT4ge1xuICAgICAgICBjb25zdCBwcm9taXNlID0gcGVuZGluZy5nZXQoaWQpO1xuICAgICAgICBpZiAoIXByb21pc2UpIHJldHVybjtcbiAgICAgICAgcGVuZGluZy5kZWxldGUoaWQpO1xuICAgICAgICBpZiAoZXJyb3IpIHtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KG5ldyBFcnJvcihlcnJvcikpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBjb25zdCBlbnZlbG9wZSA9IEpTT04ucGFyc2UocmVzcG9uc2UgPz8gXCJ7fVwiKTtcbiAgICAgICAgICAgIGlmICghZW52ZWxvcGUub2spIHtcbiAgICAgICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgRXJyb3IoZW52ZWxvcGUuZXJyb3IgPz8gXCJ1bmtub3duIHJ1bnRpbWUgY2FsbCBlcnJvclwiKSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgcHJvbWlzZS5yZXNvbHZlKFwidGV4dFwiIGluIGVudmVsb3BlID8gZW52ZWxvcGUudGV4dCA6IGVudmVsb3BlLmRhdGEpO1xuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChlKTtcbiAgICAgICAgfVxuICAgIH07XG5cbiAgICBjdXN0b21UcmFuc3BvcnQgPSB7XG4gICAgICAgIGNhbGwob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcbiAgICAgICAgICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgY29uc3QgaWQgPSBuYW5vaWQoKTtcbiAgICAgICAgICAgICAgICBwZW5kaW5nLnNldChpZCwgeyByZXNvbHZlLCByZWplY3QgfSk7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgYW5kcm9pZEJyaWRnZS5pbnZva2VBc3luYyhpZCwgSlNPTi5zdHJpbmdpZnkoe1xuICAgICAgICAgICAgICAgICAgICAgICAgb2JqZWN0OiBvYmplY3RJRCxcbiAgICAgICAgICAgICAgICAgICAgICAgIG1ldGhvZDogbWV0aG9kLFxuICAgICAgICAgICAgICAgICAgICAgICAgd2luZG93TmFtZTogd2luZG93TmFtZSxcbiAgICAgICAgICAgICAgICAgICAgICAgIGFyZ3M6IGFyZ3MgPz8gbnVsbCxcbiAgICAgICAgICAgICAgICAgICAgICAgIGNsaWVudElkOiBjbGllbnRJZCxcbiAgICAgICAgICAgICAgICAgICAgfSkpO1xuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gRG9uJ3QgbGVhayB0aGUgcGVuZGluZyBlbnRyeSBpZiBkaXNwYXRjaCB0aHJvd3Mgc3luY2hyb25vdXNseVxuICAgICAgICAgICAgICAgICAgICBwZW5kaW5nLmRlbGV0ZShpZCk7XG4gICAgICAgICAgICAgICAgICAgIHJlamVjdChlKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9KTtcbiAgICAgICAgfSxcbiAgICB9O1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbi8vIHNldHVwXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG59XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkRpYWxvZyk7XG5cbi8vIERlZmluZSBjb25zdGFudHMgZnJvbSB0aGUgYG1ldGhvZHNgIG9iamVjdCBpbiBUaXRsZSBDYXNlXG5jb25zdCBEaWFsb2dJbmZvID0gMDtcbmNvbnN0IERpYWxvZ1dhcm5pbmcgPSAxO1xuY29uc3QgRGlhbG9nRXJyb3IgPSAyO1xuY29uc3QgRGlhbG9nUXVlc3Rpb24gPSAzO1xuY29uc3QgRGlhbG9nT3BlbkZpbGUgPSA0O1xuY29uc3QgRGlhbG9nU2F2ZUZpbGUgPSA1O1xuXG5leHBvcnQgaW50ZXJmYWNlIE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIG11bHRpcGxlIHNlbGVjdGlvbiBpcyBhbGxvd2VkLiAqL1xuICAgIEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFNhdmVGaWxlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIERlZmF1bHQgZmlsZW5hbWUgdG8gdXNlIGluIHRoZSBkaWFsb2cuICovXG4gICAgRmlsZW5hbWU/OiBzdHJpbmc7XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNyZWF0ZWQuICovXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cbiAgICBTaG93SGlkZGVuRmlsZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBleHRlbnNpb25zIGNhbiBiZSBzZWxlY3RlZC4gKi9cbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cbiAgICBUcmVhdHNGaWxlUGFja2FnZXNBc0RpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xuICAgIC8qKiBBcnJheSBvZiBmaWxlIGZpbHRlcnMuICovXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogVGV4dCB0byBkaXNwbGF5IG9uIHRoZSBidXR0b24uICovXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBEaXJlY3Rvcnk/OiBzdHJpbmc7XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIE1lc3NhZ2VEaWFsb2dPcHRpb25zIHtcbiAgICAvKiogVGhlIHRpdGxlIG9mIHRoZSBkaWFsb2cgd2luZG93LiAqL1xuICAgIFRpdGxlPzogc3RyaW5nO1xuICAgIC8qKiBUaGUgbWFpbiBtZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBNZXNzYWdlPzogc3RyaW5nO1xuICAgIC8qKiBBcnJheSBvZiBidXR0b24gb3B0aW9ucyB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgQnV0dG9ucz86IEJ1dHRvbltdO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdyAoaWYgYXBwbGljYWJsZSkuICovXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEJ1dHRvbiB7XG4gICAgLyoqIFRleHQgdGhhdCBhcHBlYXJzIHdpdGhpbiB0aGUgYnV0dG9uLiAqL1xuICAgIExhYmVsPzogc3RyaW5nO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGNhbmNlbCBhbiBvcGVyYXRpb24gd2hlbiBjbGlja2VkLiAqL1xuICAgIElzQ2FuY2VsPzogYm9vbGVhbjtcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBiZSB0aGUgZGVmYXVsdCBhY3Rpb24gd2hlbiB0aGUgdXNlciBwcmVzc2VzIGVudGVyLiAqL1xuICAgIElzRGVmYXVsdD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgRmlsZUZpbHRlciB7XG4gICAgLyoqIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLiAqL1xuICAgIERpc3BsYXlOYW1lPzogc3RyaW5nO1xuICAgIC8qKiBQYXR0ZXJuIHRvIG1hdGNoIGZvciB0aGUgZmlsdGVyLCBlLmcuIFwiKi50eHQ7Ki5tZFwiIGZvciB0ZXh0IG1hcmtkb3duIGZpbGVzLiAqL1xuICAgIFBhdHRlcm4/OiBzdHJpbmc7XG59XG5cbi8qKlxuICogUHJlc2VudHMgYSBkaWFsb2cgb2Ygc3BlY2lmaWVkIHR5cGUgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqXG4gKiBAcGFyYW0gdHlwZSAtIERpYWxvZyB0eXBlLlxuICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIGZvciB0aGUgZGlhbG9nLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCByZXN1bHQgb2YgZGlhbG9nLlxuICovXG5mdW5jdGlvbiBkaWFsb2codHlwZTogbnVtYmVyLCBvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyB8IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyB8IFNhdmVGaWxlRGlhbG9nT3B0aW9ucyA9IHt9KTogUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gY2FsbCh0eXBlLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBQcmVzZW50cyBhbiBpbmZvIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSB3YXJuaW5nIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBXYXJuaW5nKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dXYXJuaW5nLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGFuIGVycm9yIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSBxdWVzdGlvbiBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgb25lIG9yIG1vcmUgZmlsZXMgdG8gb3Blbi5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSBvciBsaXN0IG9mIGZpbGVzLCBvciBhIGJsYW5rIHN0cmluZy9lbXB0eSBsaXN0IGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uOiB0cnVlIH0pOiBQcm9taXNlPHN0cmluZ1tdPjtcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMgJiB7IEFsbG93c011bHRpcGxlU2VsZWN0aW9uPzogZmFsc2UgfCB1bmRlZmluZWQgfSk6IFByb21pc2U8c3RyaW5nPjtcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZyB8IHN0cmluZ1tdPjtcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZyB8IHN0cmluZ1tdPiB7IHJldHVybiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpID8/IFtdOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSBmaWxlIHNlbGVjdGlvbiBkaWFsb2cgdG8gcGljayBhIGZpbGUgdG8gc2F2ZS5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgU2VsZWN0ZWQgZmlsZSwgb3IgYSBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBoYXMgYmVlbiBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnM6IFNhdmVGaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nU2F2ZUZpbGUsIG9wdGlvbnMpOyB9XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgZXZlbnRMaXN0ZW5lcnMsIExpc3RlbmVyLCBsaXN0ZW5lck9mZiB9IGZyb20gXCIuL2xpc3RlbmVyLmpzXCI7XG5pbXBvcnQgeyBFdmVudHMgYXMgQ3JlYXRlIH0gZnJvbSBcIi4vY3JlYXRlLmpzXCI7XG5pbXBvcnQgeyBUeXBlcyB9IGZyb20gXCIuL2V2ZW50X3R5cGVzLmpzXCI7XG5cbi8vIFNldHVwXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG4gICAgd2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XG59XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkV2ZW50cyk7XG5jb25zdCBFbWl0TWV0aG9kID0gMDtcblxuZXhwb3J0ICogZnJvbSBcIi4vZXZlbnRfdHlwZXMuanNcIjtcblxuLyoqXG4gKiBBIHRhYmxlIG9mIGRhdGEgdHlwZXMgZm9yIGFsbCBrbm93biBldmVudHMuXG4gKiBXaWxsIGJlIG1vbmtleS1wYXRjaGVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cbiAqL1xuZXhwb3J0IGludGVyZmFjZSBDdXN0b21FdmVudHMge31cblxuLyoqXG4gKiBFaXRoZXIgYSBrbm93biBldmVudCBuYW1lIG9yIGFuIGFyYml0cmFyeSBzdHJpbmcuXG4gKi9cbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnROYW1lPEUgZXh0ZW5kcyBrZXlvZiBDdXN0b21FdmVudHMgPSBrZXlvZiBDdXN0b21FdmVudHM+ID0gRSB8IChzdHJpbmcgJiB7fSk7XG5cbi8qKlxuICogVW5pb24gb2YgYWxsIGtub3duIHN5c3RlbSBldmVudCBuYW1lcy5cbiAqL1xudHlwZSBTeXN0ZW1FdmVudE5hbWUgPSB7XG4gICAgW0sgaW4ga2V5b2YgKHR5cGVvZiBUeXBlcyldOiAodHlwZW9mIFR5cGVzKVtLXVtrZXlvZiAoKHR5cGVvZiBUeXBlcylbS10pXVxufSBleHRlbmRzIChpbmZlciBNKSA/IE1ba2V5b2YgTV0gOiBuZXZlcjtcblxuLyoqXG4gKiBUaGUgZGF0YSB0eXBlIGFzc29jaWF0ZWQgdG8gYSBnaXZlbiBldmVudC5cbiAqL1xuZXhwb3J0IHR5cGUgV2FpbHNFdmVudERhdGE8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+ID1cbiAgICBFIGV4dGVuZHMga2V5b2YgQ3VzdG9tRXZlbnRzID8gQ3VzdG9tRXZlbnRzW0VdIDogKEUgZXh0ZW5kcyBTeXN0ZW1FdmVudE5hbWUgPyB2b2lkIDogYW55KTtcblxuLyoqXG4gKiBUaGUgdHlwZSBvZiBoYW5kbGVycyBmb3IgYSBnaXZlbiBldmVudC5cbiAqL1xuZXhwb3J0IHR5cGUgV2FpbHNFdmVudENhbGxiYWNrPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiA9IChldjogV2FpbHNFdmVudDxFPikgPT4gdm9pZDtcblxuLyoqXG4gKiBSZXByZXNlbnRzIGEgc3lzdGVtIGV2ZW50IG9yIGEgY3VzdG9tIGV2ZW50IGVtaXR0ZWQgdGhyb3VnaCB3YWlscy1wcm92aWRlZCBmYWNpbGl0aWVzLlxuICovXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4ge1xuICAgIC8qKlxuICAgICAqIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cbiAgICAgKi9cbiAgICBuYW1lOiBFO1xuXG4gICAgLyoqXG4gICAgICogT3B0aW9uYWwgZGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGVtaXR0ZWQgZXZlbnQuXG4gICAgICovXG4gICAgZGF0YTogV2FpbHNFdmVudERhdGE8RT47XG5cbiAgICAvKipcbiAgICAgKiBOYW1lIG9mIHRoZSBvcmlnaW5hdGluZyB3aW5kb3cuIE9taXR0ZWQgZm9yIGFwcGxpY2F0aW9uIGV2ZW50cy5cbiAgICAgKiBXaWxsIGJlIG92ZXJyaWRkZW4gaWYgc2V0IG1hbnVhbGx5LlxuICAgICAqL1xuICAgIHNlbmRlcj86IHN0cmluZztcblxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IEUsIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+KTtcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiBleHRlbmRzIG51bGwgfCB2b2lkID8gRSA6IG5ldmVyKVxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IEUsIGRhdGE/OiBhbnkpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YSA/PyBudWxsO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50OiBhbnkpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmICghbGlzdGVuZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBsZXQgd2FpbHNFdmVudCA9IG5ldyBXYWlsc0V2ZW50KFxuICAgICAgICBldmVudC5uYW1lLFxuICAgICAgICAoZXZlbnQubmFtZSBpbiBDcmVhdGUpID8gQ3JlYXRlW2V2ZW50Lm5hbWVdKGV2ZW50LmRhdGEpIDogZXZlbnQuZGF0YVxuICAgICk7XG4gICAgaWYgKCdzZW5kZXInIGluIGV2ZW50KSB7XG4gICAgICAgIHdhaWxzRXZlbnQuc2VuZGVyID0gZXZlbnQuc2VuZGVyO1xuICAgIH1cblxuICAgIC8vIERpc3BhdGNoIHRvIGEgc25hcHNob3QsIHRoZW4gcmVtb3ZlIGFsbCBleHBpcmVkIGxpc3RlbmVycyBpbiBhIHNpbmdsZVxuICAgIC8vIHBvc3QtZGlzcGF0Y2ggZmlsdGVyIG9mIHRoZSBsaXZlIG1hcC5cbiAgICAvLyAtIFdyaXRpbmcgdGhlIHNuYXBzaG90IGJhY2sgd2hvbGVzYWxlIHdvdWxkIHVuZG8gc3Vic2NyaXB0aW9uIGNoYW5nZXNcbiAgICAvLyAgIG1hZGUgaW5zaWRlIGEgaGFuZGxlciAoIzQzOTMpLlxuICAgIC8vIC0gQ2FsbGluZyBsaXN0ZW5lck9mZigpIHBlciBleHBpcmVkIGxpc3RlbmVyIGluc2lkZSB0aGUgbG9vcCBpcyBPKG5cdTAwQjIpXG4gICAgLy8gICB3aGVuIG1hbnkgbGlzdGVuZXJzIGV4cGlyZSBvbiB0aGUgc2FtZSBldmVudC5cbiAgICAvLyBGaWx0ZXJpbmcgdGhlIGxpdmUgYXJyYXkgb25jZSBhZnRlciBkaXNwYXRjaCBpcyBPKG4pIGFuZCBzdGlsbCBob25vdXJzXG4gICAgLy8gYW55IGxpc3RlbmVycyBhZGRlZCBvciByZW1vdmVkIGJ5IGhhbmRsZXJzIGR1cmluZyBkaXNwYXRjaC5cbiAgICBjb25zdCBleHBpcmVkID0gbmV3IFNldDxMaXN0ZW5lcj4oKTtcbiAgICBmb3IgKGNvbnN0IGxpc3RlbmVyIG9mIGxpc3RlbmVycy5zbGljZSgpKSB7XG4gICAgICAgIGlmIChsaXN0ZW5lci5kaXNwYXRjaCh3YWlsc0V2ZW50KSkge1xuICAgICAgICAgICAgZXhwaXJlZC5hZGQobGlzdGVuZXIpO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmIChleHBpcmVkLnNpemUgPiAwKSB7XG4gICAgICAgIGNvbnN0IGxpdmUgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XG4gICAgICAgIGlmIChsaXZlKSB7XG4gICAgICAgICAgICBjb25zdCByZW1haW5pbmcgPSBsaXZlLmZpbHRlcihsID0+ICFleHBpcmVkLmhhcyhsKSk7XG4gICAgICAgICAgICBpZiAocmVtYWluaW5nLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIHJlbWFpbmluZyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcGFyYW0gbWF4Q2FsbGJhY2tzIC0gVGhlIG1heGltdW0gbnVtYmVyIG9mIHRpbWVzIHRoZSBjYWxsYmFjayBjYW4gYmUgY2FsbGVkIGZvciB0aGUgZXZlbnQuIE9uY2UgdGhlIG1heGltdW0gbnVtYmVyIGlzIHJlYWNoZWQsIHRoZSBjYWxsYmFjayB3aWxsIG5vIGxvbmdlciBiZSBjYWxsZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+LCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbjxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+KTogKCkgPT4gdm9pZCB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihldmVudE5hbWU6IEUsIGNhbGxiYWNrOiBXYWlsc0V2ZW50Q2FsbGJhY2s8RT4pOiAoKSA9PiB2b2lkIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuLyoqXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lcyAtIFRoZSBuYW1lIG9mIHRoZSBldmVudHMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoLi4uZXZlbnROYW1lczogW1dhaWxzRXZlbnROYW1lLCAuLi5XYWlsc0V2ZW50TmFtZVtdXSk6IHZvaWQge1xuICAgIGV2ZW50TmFtZXMuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xufVxuXG4vKipcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpOiB2b2lkIHtcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xufVxuXG4vKipcbiAqIEVtaXRzIGFuIGV2ZW50LlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgZnVsZmlsbGVkIG9uY2UgdGhlIGV2ZW50IGhhcyBiZWVuIGVtaXR0ZWQuICBSZXNvbHZlcyB0byB0cnVlIGlmIHRoZSBldmVudCB3YXMgY2FuY2VsbGVkLlxuICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdFxuICogQHBhcmFtIGRhdGEgLSBUaGUgZGF0YSB0aGF0IHdpbGwgYmUgc2VudCB3aXRoIHRoZSBldmVudFxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogRSwgZGF0YTogV2FpbHNFdmVudERhdGE8RT4pOiBQcm9taXNlPGJvb2xlYW4+XG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogV2FpbHNFdmVudERhdGE8RT4gZXh0ZW5kcyBudWxsIHwgdm9pZCA/IEUgOiBuZXZlcik6IFByb21pc2U8Ym9vbGVhbj5cbmV4cG9ydCBmdW5jdGlvbiBFbWl0PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiwgZGF0YT86IGFueSk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgIHJldHVybiBjYWxsKEVtaXRNZXRob2QsICBuZXcgV2FpbHNFdmVudChuYW1lLCBkYXRhKSlcbn1cblxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBUaGUgZm9sbG93aW5nIHV0aWxpdGllcyBoYXZlIGJlZW4gZmFjdG9yZWQgb3V0IG9mIC4vZXZlbnRzLnRzXG4vLyBmb3IgdGVzdGluZyBwdXJwb3Nlcy5cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcDxzdHJpbmcsIExpc3RlbmVyW10+KCk7XG5cbmV4cG9ydCBjbGFzcyBMaXN0ZW5lciB7XG4gICAgZXZlbnROYW1lOiBzdHJpbmc7XG4gICAgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQ7XG4gICAgbWF4Q2FsbGJhY2tzOiBudW1iZXI7XG5cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQsIG1heENhbGxiYWNrczogbnVtYmVyKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICB0aGlzLmNhbGxiYWNrID0gY2FsbGJhY2s7XG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgIH1cblxuICAgIGRpc3BhdGNoKGRhdGE6IGFueSk6IGJvb2xlYW4ge1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgdGhpcy5jYWxsYmFjayhkYXRhKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKGVycik7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSByZXR1cm4gZmFsc2U7XG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBsaXN0ZW5lck9mZihsaXN0ZW5lcjogTGlzdGVuZXIpOiB2b2lkIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGxpc3RlbmVyLmV2ZW50TmFtZSk7XG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGxpc3RlbmVyLmV2ZW50TmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGxpc3RlbmVyLmV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogQW55IGlzIGEgZHVtbXkgY3JlYXRpb24gZnVuY3Rpb24gZm9yIHNpbXBsZSBvciB1bmtub3duIHR5cGVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQW55PFQgPSBhbnk+KHNvdXJjZTogYW55KTogVCB7XG4gICAgcmV0dXJuIHNvdXJjZTtcbn1cblxuLyoqXG4gKiBCeXRlU2xpY2UgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IHJlcGxhY2VzXG4gKiBudWxsIHN0cmluZ3Mgd2l0aCBlbXB0eSBzdHJpbmdzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnl0ZVNsaWNlKHNvdXJjZTogYW55KTogc3RyaW5nIHtcbiAgICByZXR1cm4gKChzb3VyY2UgPT0gbnVsbCkgPyBcIlwiIDogc291cmNlKTtcbn1cblxuLyoqXG4gKiBBcnJheSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFycmF5XG4gKiB3aG9zZSBlbGVtZW50cyBhcmUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQXJyYXk8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiBUW10ge1xuICAgIGlmIChlbGVtZW50ID09PSBBbnkpIHtcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBbXSA6IHNvdXJjZSk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xuICAgICAgICAgICAgcmV0dXJuIFtdO1xuICAgICAgICB9XG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc291cmNlLmxlbmd0aDsgaSsrKSB7XG4gICAgICAgICAgICBzb3VyY2VbaV0gPSBlbGVtZW50KHNvdXJjZVtpXSk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcbiAgICB9O1xufVxuXG4vKipcbiAqIE1hcCB0YWtlcyBjcmVhdGlvbiBmdW5jdGlvbnMgZm9yIHR3byBhcmJpdHJhcnkgdHlwZXNcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBvYmplY3RcbiAqIHdob3NlIGtleXMgYW5kIHZhbHVlcyBhcmUgb2YgdGhvc2UgdHlwZXMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBNYXA8SyBleHRlbmRzIFByb3BlcnR5S2V5ID0gYW55LCBWID0gYW55PihrZXk6IChzb3VyY2U6IGFueSkgPT4gSywgdmFsdWU6IChzb3VyY2U6IGFueSkgPT4gVik6IChzb3VyY2U6IGFueSkgPT4gUmVjb3JkPEssIFY+IHtcbiAgICBpZiAodmFsdWUgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IHt9IDogc291cmNlKTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm4ge307XG4gICAgICAgIH1cbiAgICAgICAgZm9yIChjb25zdCBrZXkgaW4gc291cmNlKSB7XG4gICAgICAgICAgICBzb3VyY2Vba2V5XSA9IHZhbHVlKHNvdXJjZVtrZXldKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogTnVsbGFibGUgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgbnVsbGFibGUgdmFsdWUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTnVsbGFibGU8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiAoVCB8IG51bGwpIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcbn1cblxuLyoqXG4gKiBTdHJ1Y3QgdGFrZXMgYW4gb2JqZWN0IG1hcHBpbmcgZmllbGQgbmFtZXMgdG8gY3JlYXRpb24gZnVuY3Rpb25zXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBzdHJ1Y3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQ6IFJlY29yZDxzdHJpbmcsIChzb3VyY2U6IGFueSkgPT4gYW55Pik6XG4gICAgPFUgZXh0ZW5kcyBSZWNvcmQ8c3RyaW5nLCBhbnk+ID0gYW55Pihzb3VyY2U6IGFueSkgPT4gVVxue1xuICAgIGxldCBhbGxBbnkgPSB0cnVlO1xuICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICBpZiAoY3JlYXRlRmllbGRbbmFtZV0gIT09IEFueSkge1xuICAgICAgICAgICAgYWxsQW55ID0gZmFsc2U7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoYWxsQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XG4gICAgICAgICAgICBpZiAobmFtZSBpbiBzb3VyY2UpIHtcbiAgICAgICAgICAgICAgICBzb3VyY2VbbmFtZV0gPSBjcmVhdGVGaWVsZFtuYW1lXShzb3VyY2VbbmFtZV0pO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBEYXRlRnJvbVRpbWUgaXMgYSBjcmVhdGlvbiBmdW5jdGlvbiB0aGF0IGNvbnZlcnRzIFJGQzMzMzkgc3RyaW5nc1xuICogKGZyb20gR28ncyB0aW1lLlRpbWUgSlNPTiBtYXJzaGFsaW5nKSB0byBKYXZhU2NyaXB0IERhdGUgb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIERhdGVGcm9tVGltZShzb3VyY2U6IGFueSk6IERhdGUge1xuICAgIHJldHVybiBuZXcgRGF0ZShzb3VyY2UpO1xufVxuXG4vKipcbiAqIE1hcHMga25vd24gZXZlbnQgbmFtZXMgdG8gY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0aGVpciBkYXRhIHR5cGVzLlxuICogV2lsbCBiZSBtb25rZXktcGF0Y2hlZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXG4gKi9cbmV4cG9ydCBjb25zdCBFdmVudHM6IFJlY29yZDxzdHJpbmcsIChzb3VyY2U6IGFueSkgPT4gYW55PiA9IHt9O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vLyBDeW5oeXJjaHd5ZCB5IGZmZWlsIGhvbiB5biBhd3RvbWF0aWcuIFBFSURJV0NIIFx1MDBDMiBNT0RJV0xcbi8vIFRoaXMgZmlsZSBpcyBhdXRvbWF0aWNhbGx5IGdlbmVyYXRlZC4gRE8gTk9UIEVESVRcblxuZXhwb3J0IGNvbnN0IFR5cGVzID0gT2JqZWN0LmZyZWV6ZSh7XG5cdFdpbmRvd3M6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFQTVBvd2VyU2V0dGluZ0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU2V0dGluZ0NoYW5nZVwiLFxuXHRcdEFQTVBvd2VyU3RhdHVzQ2hhbmdlOiBcIndpbmRvd3M6QVBNUG93ZXJTdGF0dXNDaGFuZ2VcIixcblx0XHRBUE1SZXN1bWVBdXRvbWF0aWM6IFwid2luZG93czpBUE1SZXN1bWVBdXRvbWF0aWNcIixcblx0XHRBUE1SZXN1bWVTdXNwZW5kOiBcIndpbmRvd3M6QVBNUmVzdW1lU3VzcGVuZFwiLFxuXHRcdEFQTVN1c3BlbmQ6IFwid2luZG93czpBUE1TdXNwZW5kXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcIndpbmRvd3M6QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcIndpbmRvd3M6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWQ6IFwid2luZG93czpXZWJWaWV3TmF2aWdhdGlvbkNvbXBsZXRlZFwiLFxuXHRcdFdpbmRvd0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0FjdGl2ZVwiLFxuXHRcdFdpbmRvd0JhY2tncm91bmRFcmFzZTogXCJ3aW5kb3dzOldpbmRvd0JhY2tncm91bmRFcmFzZVwiLFxuXHRcdFdpbmRvd0NsaWNrQWN0aXZlOiBcIndpbmRvd3M6V2luZG93Q2xpY2tBY3RpdmVcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcIndpbmRvd3M6V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwid2luZG93czpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIndpbmRvd3M6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJ3aW5kb3dzOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dEcmFnRHJvcDogXCJ3aW5kb3dzOldpbmRvd0RyYWdEcm9wXCIsXG5cdFx0V2luZG93RHJhZ0VudGVyOiBcIndpbmRvd3M6V2luZG93RHJhZ0VudGVyXCIsXG5cdFx0V2luZG93RHJhZ0xlYXZlOiBcIndpbmRvd3M6V2luZG93RHJhZ0xlYXZlXCIsXG5cdFx0V2luZG93RHJhZ092ZXI6IFwid2luZG93czpXaW5kb3dEcmFnT3ZlclwiLFxuXHRcdFdpbmRvd0VuZE1vdmU6IFwid2luZG93czpXaW5kb3dFbmRNb3ZlXCIsXG5cdFx0V2luZG93RW5kUmVzaXplOiBcIndpbmRvd3M6V2luZG93RW5kUmVzaXplXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dIaWRlOiBcIndpbmRvd3M6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0luYWN0aXZlOiBcIndpbmRvd3M6V2luZG93SW5hY3RpdmVcIixcblx0XHRXaW5kb3dLZXlEb3duOiBcIndpbmRvd3M6V2luZG93S2V5RG93blwiLFxuXHRcdFdpbmRvd0tleVVwOiBcIndpbmRvd3M6V2luZG93S2V5VXBcIixcblx0XHRXaW5kb3dLaWxsRm9jdXM6IFwid2luZG93czpXaW5kb3dLaWxsRm9jdXNcIixcblx0XHRXaW5kb3dOb25DbGllbnRIaXQ6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRIaXRcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZURvd246IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZURvd25cIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZUxlYXZlOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlTW92ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTW92ZVwiLFxuXHRcdFdpbmRvd05vbkNsaWVudE1vdXNlVXA6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZVVwXCIsXG5cdFx0V2luZG93UGFpbnQ6IFwid2luZG93czpXaW5kb3dQYWludFwiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwid2luZG93czpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93U2V0Rm9jdXM6IFwid2luZG93czpXaW5kb3dTZXRGb2N1c1wiLFxuXHRcdFdpbmRvd1Nob3c6IFwid2luZG93czpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93U3RhcnRNb3ZlOiBcIndpbmRvd3M6V2luZG93U3RhcnRNb3ZlXCIsXG5cdFx0V2luZG93U3RhcnRSZXNpemU6IFwid2luZG93czpXaW5kb3dTdGFydFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJ3aW5kb3dzOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1pPcmRlckNoYW5nZWQ6IFwid2luZG93czpXaW5kb3daT3JkZXJDaGFuZ2VkXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwid2luZG93czpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwid2luZG93czpXaW5kb3dVbk1heGltaXNlXCIsXG5cdH0pLFxuXHRNYWM6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlRWZmZWN0aXZlQXBwZWFyYW5jZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvbjogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VJY29uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTY3JlZW5QYXJhbWV0ZXJzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VTdGF0dXNCYXJGcmFtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyT3JpZW50YXRpb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyT3JpZW50YXRpb25cIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVRoZW1lXCIsXG5cdFx0QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25EaWRIaWRlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVuaGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25EaWRVcGRhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZFdha2U6IFwibWFjOkFwcGxpY2F0aW9uRGlkV2FrZVwiLFxuXHRcdEFwcGxpY2F0aW9uU2NyZWVuc0RpZFNsZWVwOiBcIm1hYzpBcHBsaWNhdGlvblNjcmVlbnNEaWRTbGVlcFwiLFxuXHRcdEFwcGxpY2F0aW9uU2NyZWVuc0RpZFdha2U6IFwibWFjOkFwcGxpY2F0aW9uU2NyZWVuc0RpZFdha2VcIixcblx0XHRBcHBsaWNhdGlvblNob3VsZEhhbmRsZVJlb3BlbjogXCJtYWM6QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW5cIixcblx0XHRBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZzogXCJtYWM6QXBwbGljYXRpb25XaWxsRmluaXNoTGF1bmNoaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsSGlkZTogXCJtYWM6QXBwbGljYXRpb25XaWxsSGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsU2xlZXA6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFNsZWVwXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVuaGlkZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFVwZGF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVXBkYXRlXCIsXG5cdFx0TWVudURpZEFkZEl0ZW06IFwibWFjOk1lbnVEaWRBZGRJdGVtXCIsXG5cdFx0TWVudURpZEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRCZWdpblRyYWNraW5nXCIsXG5cdFx0TWVudURpZENsb3NlOiBcIm1hYzpNZW51RGlkQ2xvc2VcIixcblx0XHRNZW51RGlkRGlzcGxheUl0ZW06IFwibWFjOk1lbnVEaWREaXNwbGF5SXRlbVwiLFxuXHRcdE1lbnVEaWRFbmRUcmFja2luZzogXCJtYWM6TWVudURpZEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudURpZEhpZ2hsaWdodEl0ZW06IFwibWFjOk1lbnVEaWRIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudURpZE9wZW46IFwibWFjOk1lbnVEaWRPcGVuXCIsXG5cdFx0TWVudURpZFBvcFVwOiBcIm1hYzpNZW51RGlkUG9wVXBcIixcblx0XHRNZW51RGlkUmVtb3ZlSXRlbTogXCJtYWM6TWVudURpZFJlbW92ZUl0ZW1cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvbjogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25cIixcblx0XHRNZW51RGlkU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW1cIixcblx0XHRNZW51RGlkVXBkYXRlOiBcIm1hYzpNZW51RGlkVXBkYXRlXCIsXG5cdFx0TWVudVdpbGxBZGRJdGVtOiBcIm1hYzpNZW51V2lsbEFkZEl0ZW1cIixcblx0XHRNZW51V2lsbEJlZ2luVHJhY2tpbmc6IFwibWFjOk1lbnVXaWxsQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsRGlzcGxheUl0ZW06IFwibWFjOk1lbnVXaWxsRGlzcGxheUl0ZW1cIixcblx0XHRNZW51V2lsbEVuZFRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEVuZFRyYWNraW5nXCIsXG5cdFx0TWVudVdpbGxIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51V2lsbEhpZ2hsaWdodEl0ZW1cIixcblx0XHRNZW51V2lsbE9wZW46IFwibWFjOk1lbnVXaWxsT3BlblwiLFxuXHRcdE1lbnVXaWxsUG9wVXA6IFwibWFjOk1lbnVXaWxsUG9wVXBcIixcblx0XHRNZW51V2lsbFJlbW92ZUl0ZW06IFwibWFjOk1lbnVXaWxsUmVtb3ZlSXRlbVwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvbjogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uXCIsXG5cdFx0TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtOiBcIm1hYzpNZW51V2lsbFNlbmRBY3Rpb25Ub0l0ZW1cIixcblx0XHRNZW51V2lsbFVwZGF0ZTogXCJtYWM6TWVudVdpbGxVcGRhdGVcIixcblx0XHRXZWJWaWV3RGlkQ29tbWl0TmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkUmVjZWl2ZVNlcnZlclJlZGlyZWN0Rm9yUHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uOiBcIm1hYzpXZWJWaWV3RGlkU3RhcnRQcm92aXNpb25hbE5hdmlnYXRpb25cIixcblx0XHRXaW5kb3dEaWRCZWNvbWVLZXk6IFwibWFjOldpbmRvd0RpZEJlY29tZUtleVwiLFxuXHRcdFdpbmRvd0RpZEJlY29tZU1haW46IFwibWFjOldpbmRvd0RpZEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dEaWRCZWdpblNoZWV0OiBcIm1hYzpXaW5kb3dEaWRCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQWxwaGE6IFwibWFjOldpbmRvd0RpZENoYW5nZUFscGhhXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VCYWNraW5nTG9jYXRpb25cIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VCYWNraW5nUHJvcGVydGllczogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZENoYW5nZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VPY2NsdXNpb25TdGF0ZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyczogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblByb2ZpbGU6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblByb2ZpbGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNoYXJpbmdUeXBlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNwYWNlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVRpdGxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUaXRsZVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVRvb2xiYXI6IFwibWFjOldpbmRvd0RpZENoYW5nZVRvb2xiYXJcIixcblx0XHRXaW5kb3dEaWREZW1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWREZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93RGlkRW5kU2hlZXQ6IFwibWFjOldpbmRvd0RpZEVuZFNoZWV0XCIsXG5cdFx0V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFbnRlckZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhpdEZ1bGxTY3JlZW46IFwibWFjOldpbmRvd0RpZEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dEaWRFeGl0VmVyc2lvbkJyb3dzZXJcIixcblx0XHRXaW5kb3dEaWRFeHBvc2U6IFwibWFjOldpbmRvd0RpZEV4cG9zZVwiLFxuXHRcdFdpbmRvd0RpZEZvY3VzOiBcIm1hYzpXaW5kb3dEaWRGb2N1c1wiLFxuXHRcdFdpbmRvd0RpZE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dEaWRNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibWFjOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPZmZTY3JlZW5cIixcblx0XHRXaW5kb3dEaWRPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkUmVzaWduS2V5OiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25LZXlcIixcblx0XHRXaW5kb3dEaWRSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dEaWRSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcIm1hYzpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvcjogXCJtYWM6V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUaXRsZTogXCJtYWM6V2luZG93RGlkVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkWm9vbTogXCJtYWM6V2luZG93RGlkWm9vbVwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0VudGVyZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ0V4aXRlZFwiLFxuXHRcdFdpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZDogXCJtYWM6V2luZG93RmlsZURyYWdnaW5nUGVyZm9ybWVkXCIsXG5cdFx0V2luZG93SGlkZTogXCJtYWM6V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIm1hYzpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwibWFjOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJtYWM6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIm1hYzpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93U2hvdWxkQ2xvc2U6IFwibWFjOldpbmRvd1Nob3VsZENsb3NlXCIsXG5cdFx0V2luZG93U2hvdzogXCJtYWM6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1dpbGxCZWNvbWVLZXk6IFwibWFjOldpbmRvd1dpbGxCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lTWFpbjogXCJtYWM6V2luZG93V2lsbEJlY29tZU1haW5cIixcblx0XHRXaW5kb3dXaWxsQmVnaW5TaGVldDogXCJtYWM6V2luZG93V2lsbEJlZ2luU2hlZXRcIixcblx0XHRXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dXaWxsQ2hhbmdlT3JkZXJpbmdNb2RlXCIsXG5cdFx0V2luZG93V2lsbENsb3NlOiBcIm1hYzpXaW5kb3dXaWxsQ2xvc2VcIixcblx0XHRXaW5kb3dXaWxsRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93V2lsbERlbWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEVudGVyVmVyc2lvbkJyb3dzZXI6IFwibWFjOldpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsRXhpdEZ1bGxTY3JlZW5cIixcblx0XHRXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93V2lsbEZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsRm9jdXNcIixcblx0XHRXaW5kb3dXaWxsTWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxNaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd1dpbGxNb3ZlOiBcIm1hYzpXaW5kb3dXaWxsTW92ZVwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9mZlNjcmVlbjogXCJtYWM6V2luZG93V2lsbE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbE9yZGVyT25TY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9uU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbFJlc2lnbk1haW46IFwibWFjOldpbmRvd1dpbGxSZXNpZ25NYWluXCIsXG5cdFx0V2luZG93V2lsbFJlc2l6ZTogXCJtYWM6V2luZG93V2lsbFJlc2l6ZVwiLFxuXHRcdFdpbmRvd1dpbGxVbmZvY3VzOiBcIm1hYzpXaW5kb3dXaWxsVW5mb2N1c1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQWxwaGE6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVBbHBoYVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXM6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVTaGFkb3c6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVTaGFkb3dcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd1dpbGxVcGRhdGVUaXRsZVwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVUb29sYmFyOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd1dpbGxVcGRhdGVWaXNpYmlsaXR5OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eVwiLFxuXHRcdFdpbmRvd1dpbGxVc2VTdGFuZGFyZEZyYW1lOiBcIm1hYzpXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJtYWM6V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJtYWM6V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJtYWM6V2luZG93Wm9vbVJlc2V0XCIsXG5cdFx0V2ViVmlld1dlYkNvbnRlbnRQcm9jZXNzRGlkVGVybWluYXRlOiBcIm1hYzpXZWJWaWV3V2ViQ29udGVudFByb2Nlc3NEaWRUZXJtaW5hdGVcIixcblx0fSksXG5cdExpbnV4OiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvblN0YXJ0dXA6IFwibGludXg6QXBwbGljYXRpb25TdGFydHVwXCIsXG5cdFx0U3lzdGVtRGlkV2FrZTogXCJsaW51eDpTeXN0ZW1EaWRXYWtlXCIsXG5cdFx0U3lzdGVtVGhlbWVDaGFuZ2VkOiBcImxpbnV4OlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFN5c3RlbVdpbGxTbGVlcDogXCJsaW51eDpTeXN0ZW1XaWxsU2xlZXBcIixcblx0XHRXaW5kb3dEZWxldGVFdmVudDogXCJsaW51eDpXaW5kb3dEZWxldGVFdmVudFwiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwibGludXg6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJsaW51eDpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dGb2N1c0luOiBcImxpbnV4OldpbmRvd0ZvY3VzSW5cIixcblx0XHRXaW5kb3dGb2N1c091dDogXCJsaW51eDpXaW5kb3dGb2N1c091dFwiLFxuXHRcdFdpbmRvd0xvYWRTdGFydGVkOiBcImxpbnV4OldpbmRvd0xvYWRTdGFydGVkXCIsXG5cdFx0V2luZG93TG9hZFJlZGlyZWN0ZWQ6IFwibGludXg6V2luZG93TG9hZFJlZGlyZWN0ZWRcIixcblx0XHRXaW5kb3dMb2FkQ29tbWl0dGVkOiBcImxpbnV4OldpbmRvd0xvYWRDb21taXR0ZWRcIixcblx0XHRXaW5kb3dMb2FkRmluaXNoZWQ6IFwibGludXg6V2luZG93TG9hZEZpbmlzaGVkXCIsXG5cdH0pLFxuXHRpT1M6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uRGlkQmVjb21lQWN0aXZlOiBcImlvczpBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRW50ZXJCYWNrZ3JvdW5kOiBcImlvczpBcHBsaWNhdGlvbkRpZEVudGVyQmFja2dyb3VuZFwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcImlvczpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkUmVjZWl2ZU1lbW9yeVdhcm5pbmc6IFwiaW9zOkFwcGxpY2F0aW9uRGlkUmVjZWl2ZU1lbW9yeVdhcm5pbmdcIixcblx0XHRBcHBsaWNhdGlvbldpbGxFbnRlckZvcmVncm91bmQ6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbEVudGVyRm9yZWdyb3VuZFwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZTogXCJpb3M6QXBwbGljYXRpb25XaWxsUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVGVybWluYXRlOiBcImlvczpBcHBsaWNhdGlvbldpbGxUZXJtaW5hdGVcIixcblx0XHRXaW5kb3dEaWRMb2FkOiBcImlvczpXaW5kb3dEaWRMb2FkXCIsXG5cdFx0V2luZG93V2lsbEFwcGVhcjogXCJpb3M6V2luZG93V2lsbEFwcGVhclwiLFxuXHRcdFdpbmRvd0RpZEFwcGVhcjogXCJpb3M6V2luZG93RGlkQXBwZWFyXCIsXG5cdFx0V2luZG93V2lsbERpc2FwcGVhcjogXCJpb3M6V2luZG93V2lsbERpc2FwcGVhclwiLFxuXHRcdFdpbmRvd0RpZERpc2FwcGVhcjogXCJpb3M6V2luZG93RGlkRGlzYXBwZWFyXCIsXG5cdFx0V2luZG93U2FmZUFyZWFJbnNldHNDaGFuZ2VkOiBcImlvczpXaW5kb3dTYWZlQXJlYUluc2V0c0NoYW5nZWRcIixcblx0XHRXaW5kb3dPcmllbnRhdGlvbkNoYW5nZWQ6IFwiaW9zOldpbmRvd09yaWVudGF0aW9uQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd1RvdWNoQmVnYW46IFwiaW9zOldpbmRvd1RvdWNoQmVnYW5cIixcblx0XHRXaW5kb3dUb3VjaE1vdmVkOiBcImlvczpXaW5kb3dUb3VjaE1vdmVkXCIsXG5cdFx0V2luZG93VG91Y2hFbmRlZDogXCJpb3M6V2luZG93VG91Y2hFbmRlZFwiLFxuXHRcdFdpbmRvd1RvdWNoQ2FuY2VsbGVkOiBcImlvczpXaW5kb3dUb3VjaENhbmNlbGxlZFwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRTdGFydE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkRmluaXNoTmF2aWdhdGlvbjogXCJpb3M6V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGlkRmFpbE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGYWlsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEZWNpZGVQb2xpY3lGb3JOYXZpZ2F0aW9uQWN0aW9uOiBcImlvczpXZWJWaWV3RGVjaWRlUG9saWN5Rm9yTmF2aWdhdGlvbkFjdGlvblwiLFxuXHRcdEJhdHRlcnlDaGFuZ2VkOiBcImlvczpCYXR0ZXJ5Q2hhbmdlZFwiLFxuXHRcdE5ldHdvcmtDaGFuZ2VkOiBcImlvczpOZXR3b3JrQ2hhbmdlZFwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJpb3M6VGhlbWVDaGFuZ2VkXCIsXG5cdFx0U2NyZWVuTG9ja2VkOiBcImlvczpTY3JlZW5Mb2NrZWRcIixcblx0XHRTY3JlZW5VbmxvY2tlZDogXCJpb3M6U2NyZWVuVW5sb2NrZWRcIixcblx0fSksXG5cdEFuZHJvaWQ6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFjdGl2aXR5Q3JlYXRlZDogXCJhbmRyb2lkOkFjdGl2aXR5Q3JlYXRlZFwiLFxuXHRcdEFjdGl2aXR5U3RhcnRlZDogXCJhbmRyb2lkOkFjdGl2aXR5U3RhcnRlZFwiLFxuXHRcdEFjdGl2aXR5UmVzdW1lZDogXCJhbmRyb2lkOkFjdGl2aXR5UmVzdW1lZFwiLFxuXHRcdEFjdGl2aXR5UGF1c2VkOiBcImFuZHJvaWQ6QWN0aXZpdHlQYXVzZWRcIixcblx0XHRBY3Rpdml0eVN0b3BwZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eVN0b3BwZWRcIixcblx0XHRBY3Rpdml0eURlc3Ryb3llZDogXCJhbmRyb2lkOkFjdGl2aXR5RGVzdHJveWVkXCIsXG5cdFx0QXBwbGljYXRpb25Mb3dNZW1vcnk6IFwiYW5kcm9pZDpBcHBsaWNhdGlvbkxvd01lbW9yeVwiLFxuXHRcdFdlYlZpZXdQYWdlU3RhcnRlZDogXCJhbmRyb2lkOldlYlZpZXdQYWdlU3RhcnRlZFwiLFxuXHRcdFdlYlZpZXdQYWdlRmluaXNoZWQ6IFwiYW5kcm9pZDpXZWJWaWV3UGFnZUZpbmlzaGVkXCIsXG5cdFx0QmF0dGVyeUNoYW5nZWQ6IFwiYW5kcm9pZDpCYXR0ZXJ5Q2hhbmdlZFwiLFxuXHRcdE5ldHdvcmtDaGFuZ2VkOiBcImFuZHJvaWQ6TmV0d29ya0NoYW5nZWRcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiYW5kcm9pZDpUaGVtZUNoYW5nZWRcIixcblx0XHRTY3JlZW5Mb2NrZWQ6IFwiYW5kcm9pZDpTY3JlZW5Mb2NrZWRcIixcblx0XHRTY3JlZW5VbmxvY2tlZDogXCJhbmRyb2lkOlNjcmVlblVubG9ja2VkXCIsXG5cdH0pLFxuXHRDb21tb246IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGU6IFwiY29tbW9uOkFwcGxpY2F0aW9uT3BlbmVkV2l0aEZpbGVcIixcblx0XHRBcHBsaWNhdGlvblN0YXJ0ZWQ6IFwiY29tbW9uOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdEFwcGxpY2F0aW9uTGF1bmNoZWRXaXRoVXJsOiBcImNvbW1vbjpBcHBsaWNhdGlvbkxhdW5jaGVkV2l0aFVybFwiLFxuXHRcdFRoZW1lQ2hhbmdlZDogXCJjb21tb246VGhlbWVDaGFuZ2VkXCIsXG5cdFx0U3lzdGVtRGlkV2FrZTogXCJjb21tb246U3lzdGVtRGlkV2FrZVwiLFxuXHRcdFN5c3RlbVdpbGxTbGVlcDogXCJjb21tb246U3lzdGVtV2lsbFNsZWVwXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJjb21tb246V2luZG93Q2xvc2luZ1wiLFxuXHRcdFdpbmRvd0RpZE1vdmU6IFwiY29tbW9uOldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwiY29tbW9uOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwiY29tbW9uOldpbmRvd0RQSUNoYW5nZWRcIixcblx0XHRXaW5kb3dGaWxlc0Ryb3BwZWQ6IFwiY29tbW9uOldpbmRvd0ZpbGVzRHJvcHBlZFwiLFxuXHRcdFdpbmRvd0ZvY3VzOiBcImNvbW1vbjpXaW5kb3dGb2N1c1wiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd0Z1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dIaWRlOiBcImNvbW1vbjpXaW5kb3dIaWRlXCIsXG5cdFx0V2luZG93TG9zdEZvY3VzOiBcImNvbW1vbjpXaW5kb3dMb3N0Rm9jdXNcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJjb21tb246V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dNaW5pbWlzZTogXCJjb21tb246V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dUb2dnbGVGcmFtZWxlc3M6IFwiY29tbW9uOldpbmRvd1RvZ2dsZUZyYW1lbGVzc1wiLFxuXHRcdFdpbmRvd1Jlc3RvcmU6IFwiY29tbW9uOldpbmRvd1Jlc3RvcmVcIixcblx0XHRXaW5kb3dSdW50aW1lUmVhZHk6IFwiY29tbW9uOldpbmRvd1J1bnRpbWVSZWFkeVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwiY29tbW9uOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwiY29tbW9uOldpbmRvd1VuRnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd1VuTWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1pbmltaXNlXCIsXG5cdFx0V2luZG93Wm9vbTogXCJjb21tb246V2luZG93Wm9vbVwiLFxuXHRcdFdpbmRvd1pvb21JbjogXCJjb21tb246V2luZG93Wm9vbUluXCIsXG5cdFx0V2luZG93Wm9vbU91dDogXCJjb21tb246V2luZG93Wm9vbU91dFwiLFxuXHRcdFdpbmRvd1pvb21SZXNldDogXCJjb21tb246V2luZG93Wm9vbVJlc2V0XCIsXG5cdFx0QmF0dGVyeUNoYW5nZWQ6IFwiY29tbW9uOkJhdHRlcnlDaGFuZ2VkXCIsXG5cdFx0TmV0d29ya0NoYW5nZWQ6IFwiY29tbW9uOk5ldHdvcmtDaGFuZ2VkXCIsXG5cdFx0U2NyZWVuTG9ja2VkOiBcImNvbW1vbjpTY3JlZW5Mb2NrZWRcIixcblx0XHRTY3JlZW5VbmxvY2tlZDogXCJjb21tb246U2NyZWVuVW5sb2NrZWRcIixcblx0XHRMb3dNZW1vcnk6IFwiY29tbW9uOkxvd01lbW9yeVwiLFxuXHR9KSxcbn0pO1xuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIExvZ3MgYSBtZXNzYWdlIHRvIHRoZSBjb25zb2xlIHdpdGggY3VzdG9tIGZvcm1hdHRpbmcuXG4gKlxuICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkZWJ1Z0xvZyhtZXNzYWdlOiBhbnkpIHtcbiAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcbiAgICBjb25zb2xlLmxvZyhcbiAgICAgICAgJyVjIHdhaWxzMyAlYyAnICsgbWVzc2FnZSArICcgJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNhYTAwMDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAzcHggMHB4IDBweCAzcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nLFxuICAgICAgICAnYmFja2dyb3VuZDogIzAwOTkwMDsgY29sb3I6ICNmZmY7IGJvcmRlci1yYWRpdXM6IDBweCAzcHggM3B4IDBweDsgcGFkZGluZzogMXB4OyBmb250LXNpemU6IDAuN3JlbSdcbiAgICApO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSB3ZWJ2aWV3IHN1cHBvcnRzIHRoZSB7QGxpbmsgTW91c2VFdmVudCNidXR0b25zfSBwcm9wZXJ0eS5cbiAqIExvb2tpbmcgYXQgeW91IG1hY09TIEhpZ2ggU2llcnJhIVxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuVHJhY2tCdXR0b25zKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAobmV3IE1vdXNlRXZlbnQoJ21vdXNlZG93bicpKS5idXR0b25zID09PSAwO1xufVxuXG4vKipcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSBicm93c2VyIHN1cHBvcnRzIHJlbW92aW5nIGxpc3RlbmVycyBieSB0cmlnZ2VyaW5nIGFuIEFib3J0U2lnbmFsXG4gKiAoc2VlIGh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnL2VuLVVTL2RvY3MvV2ViL0FQSS9FdmVudFRhcmdldC9hZGRFdmVudExpc3RlbmVyI3NpZ25hbCkuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjYW5BYm9ydExpc3RlbmVycygpIHtcbiAgICBpZiAoIUV2ZW50VGFyZ2V0IHx8ICFBYm9ydFNpZ25hbCB8fCAhQWJvcnRDb250cm9sbGVyKVxuICAgICAgICByZXR1cm4gZmFsc2U7XG5cbiAgICBsZXQgcmVzdWx0ID0gdHJ1ZTtcblxuICAgIGNvbnN0IHRhcmdldCA9IG5ldyBFdmVudFRhcmdldCgpO1xuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgdGFyZ2V0LmFkZEV2ZW50TGlzdGVuZXIoJ3Rlc3QnLCAoKSA9PiB7IHJlc3VsdCA9IGZhbHNlOyB9LCB7IHNpZ25hbDogY29udHJvbGxlci5zaWduYWwgfSk7XG4gICAgY29udHJvbGxlci5hYm9ydCgpO1xuICAgIHRhcmdldC5kaXNwYXRjaEV2ZW50KG5ldyBDdXN0b21FdmVudCgndGVzdCcpKTtcblxuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogUmVzb2x2ZXMgdGhlIGNsb3Nlc3QgSFRNTEVsZW1lbnQgYW5jZXN0b3Igb2YgYW4gZXZlbnQncyB0YXJnZXQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBldmVudFRhcmdldChldmVudDogRXZlbnQpOiBIVE1MRWxlbWVudCB7XG4gICAgaWYgKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSB7XG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQ7XG4gICAgfSBlbHNlIGlmICghKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSAmJiBldmVudC50YXJnZXQgaW5zdGFuY2VvZiBOb2RlKSB7XG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQucGFyZW50RWxlbWVudCA/PyBkb2N1bWVudC5ib2R5O1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiBkb2N1bWVudC5ib2R5O1xuICAgIH1cbn1cblxuLyoqKlxuIFRoaXMgdGVjaG5pcXVlIGZvciBwcm9wZXIgbG9hZCBkZXRlY3Rpb24gaXMgdGFrZW4gZnJvbSBIVE1YOlxuXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcblxuIENvcHlyaWdodCAoYykgMjAyMCwgQmlnIFNreSBTb2Z0d2FyZVxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXG5cbiBSZWRpc3RyaWJ1dGlvbiBhbmQgdXNlIGluIHNvdXJjZSBhbmQgYmluYXJ5IGZvcm1zLCB3aXRoIG9yIHdpdGhvdXRcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcblxuIDEuIFJlZGlzdHJpYnV0aW9ucyBvZiBzb3VyY2UgY29kZSBtdXN0IHJldGFpbiB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSwgdGhpc1xuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxuXG4gMi4gUmVkaXN0cmlidXRpb25zIGluIGJpbmFyeSBmb3JtIG11c3QgcmVwcm9kdWNlIHRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlLFxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cbiBhbmQvb3Igb3RoZXIgbWF0ZXJpYWxzIHByb3ZpZGVkIHdpdGggdGhlIGRpc3RyaWJ1dGlvbi5cblxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXG4gQU5EIEFOWSBFWFBSRVNTIE9SIElNUExJRUQgV0FSUkFOVElFUywgSU5DTFVESU5HLCBCVVQgTk9UIExJTUlURUQgVE8sIFRIRVxuIElNUExJRUQgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFkgQU5EIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFSRVxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcbiBGT1IgQU5ZIERJUkVDVCwgSU5ESVJFQ1QsIElOQ0lERU5UQUwsIFNQRUNJQUwsIEVYRU1QTEFSWSwgT1IgQ09OU0VRVUVOVElBTFxuIERBTUFHRVMgKElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBQUk9DVVJFTUVOVCBPRiBTVUJTVElUVVRFIEdPT0RTIE9SXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcbiBDQVVTRUQgQU5EIE9OIEFOWSBUSEVPUlkgT0YgTElBQklMSVRZLCBXSEVUSEVSIElOIENPTlRSQUNULCBTVFJJQ1QgTElBQklMSVRZLFxuIE9SIFRPUlQgKElOQ0xVRElORyBORUdMSUdFTkNFIE9SIE9USEVSV0lTRSkgQVJJU0lORyBJTiBBTlkgV0FZIE9VVCBPRiBUSEUgVVNFXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cblxuICoqKi9cblxubGV0IGlzUmVhZHkgPSBmYWxzZTtcbmltcG9ydCB7IGhhc0RPTSB9IGZyb20gXCIuL2Vudmlyb25tZW50LmpzXCI7XG5cbmlmIChoYXNET00pIHtcbiAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgKCkgPT4geyBpc1JlYWR5ID0gdHJ1ZSB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjazogKCkgPT4gdm9pZCkge1xuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcbiAgICAgICAgY2FsbGJhY2soKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgY2FsbGJhY2spO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHR5cGUgeyBTY3JlZW4gfSBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XG5cbi8vIERyb3AgdGFyZ2V0IGNvbnN0YW50c1xuY29uc3QgRFJPUF9UQVJHRVRfQVRUUklCVVRFID0gJ2RhdGEtZmlsZS1kcm9wLXRhcmdldCc7XG5jb25zdCBEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MgPSAnZmlsZS1kcm9wLXRhcmdldC1hY3RpdmUnO1xubGV0IGN1cnJlbnREcm9wVGFyZ2V0OiBFbGVtZW50IHwgbnVsbCA9IG51bGw7XG5cbmNvbnN0IFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDA7XG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xuY29uc3QgQ2xvc2VNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gMjtcbmNvbnN0IERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICA9IDM7XG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xuY29uc3QgRm9jdXNNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNTtcbmNvbnN0IEZvcmNlUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICA9IDY7XG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xuY29uc3QgR2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gODtcbmNvbnN0IEdldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDk7XG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcbmNvbnN0IEhpZGVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDExO1xuY29uc3QgSXNGb2N1c2VkTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMTI7XG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcbmNvbnN0IElzTWF4aW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE0O1xuY29uc3QgSXNNaW5pbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTU7XG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcbmNvbnN0IE1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE3O1xuY29uc3QgTmFtZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTg7XG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcbmNvbnN0IFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICAgICA9IDIwO1xuY29uc3QgUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMjE7XG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcbmNvbnN0IFJlc3RvcmVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDIzO1xuY29uc3QgU2V0UG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgID0gMjQ7XG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcbmNvbnN0IFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QgICAgICAgICA9IDI2O1xuY29uc3QgU2V0RnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgICAgID0gMjc7XG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcbmNvbnN0IFNldE1heFNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDI5O1xuY29uc3QgU2V0TWluU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMzA7XG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcbmNvbnN0IFNldFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICA9IDMyO1xuY29uc3QgU2V0U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzM7XG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcbmNvbnN0IFNldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDM1O1xuY29uc3QgU2hvd01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzY7XG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcbmNvbnN0IFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICA9IDM4O1xuY29uc3QgVG9nZ2xlTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgID0gMzk7XG5jb25zdCBUb2dnbGVGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgPSA0MDsgXG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MTtcbmNvbnN0IFVuTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQyO1xuY29uc3QgVW5NaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDM7XG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0NDtcbmNvbnN0IFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDQ1O1xuY29uc3QgWm9vbUluTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gNDY7XG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NztcbmNvbnN0IFpvb21SZXNldE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDQ4O1xuY29uc3QgU25hcEFzc2lzdE1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDk7XG5jb25zdCBGaWxlc0Ryb3BwZWQgICAgICAgICAgICAgICAgICAgICAgPSA1MDtcbmNvbnN0IFByaW50TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDUxO1xuY29uc3QgU2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gNTI7XG5cbi8qKlxuICogRmluZHMgdGhlIG5lYXJlc3QgZHJvcCB0YXJnZXQgZWxlbWVudCBieSB3YWxraW5nIHVwIHRoZSBET00gdHJlZS5cbiAqL1xuZnVuY3Rpb24gZ2V0RHJvcFRhcmdldEVsZW1lbnQoZWxlbWVudDogRWxlbWVudCB8IG51bGwpOiBFbGVtZW50IHwgbnVsbCB7XG4gICAgaWYgKCFlbGVtZW50KSB7XG4gICAgICAgIHJldHVybiBudWxsO1xuICAgIH1cbiAgICByZXR1cm4gZWxlbWVudC5jbG9zZXN0KGBbJHtEUk9QX1RBUkdFVF9BVFRSSUJVVEV9XWApO1xufVxuXG4vKipcbiAqIENoZWNrIGlmIHdlIGNhbiB1c2UgV2ViVmlldzIncyBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyAoV2luZG93cylcbiAqIEFsc28gY2hlY2tzIHRoYXQgRW5hYmxlRmlsZURyb3AgaXMgdHJ1ZSBmb3IgdGhpcyB3aW5kb3cuXG4gKi9cbmZ1bmN0aW9uIGNhblJlc29sdmVGaWxlUGF0aHMoKTogYm9vbGVhbiB7XG4gICAgLy8gTXVzdCBoYXZlIFdlYlZpZXcyJ3MgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgQVBJIChXaW5kb3dzIG9ubHkpXG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzID09IG51bGwpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbiAgICAvLyBNdXN0IGhhdmUgRW5hYmxlRmlsZURyb3Agc2V0IHRvIHRydWUgZm9yIHRoaXMgd2luZG93XG4gICAgLy8gVGhpcyBmbGFnIGlzIHNldCBieSB0aGUgR28gYmFja2VuZCBkdXJpbmcgcnVudGltZSBpbml0aWFsaXphdGlvblxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IHRydWU7XG59XG5cbi8qKlxuICogU2VuZCBmaWxlIGRyb3AgdG8gYmFja2VuZCB2aWEgV2ViVmlldzIgKFdpbmRvd3Mgb25seSlcbiAqL1xuZnVuY3Rpb24gcmVzb2x2ZUZpbGVQYXRocyh4OiBudW1iZXIsIHk6IG51bWJlciwgZmlsZXM6IEZpbGVbXSk6IHZvaWQge1xuICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cykge1xuICAgICAgICAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMoYGZpbGU6ZHJvcDoke3h9OiR7eX1gLCBmaWxlcyk7XG4gICAgfVxufVxuXG4vLyBOYXRpdmUgZHJhZyBzdGF0ZSAoTGludXgvbWFjT1MgaW50ZXJjZXB0IERPTSBkcmFnIGV2ZW50cylcbmxldCBuYXRpdmVEcmFnQWN0aXZlID0gZmFsc2U7XG5cbi8qKlxuICogQ2xlYW5zIHVwIG5hdGl2ZSBkcmFnIHN0YXRlIGFuZCBob3ZlciBlZmZlY3RzLlxuICogQ2FsbGVkIG9uIGRyb3Agb3Igd2hlbiBkcmFnIGxlYXZlcyB0aGUgd2luZG93LlxuICovXG5mdW5jdGlvbiBjbGVhbnVwTmF0aXZlRHJhZygpOiB2b2lkIHtcbiAgICBuYXRpdmVEcmFnQWN0aXZlID0gZmFsc2U7XG4gICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgIH1cbn1cblxuLyoqXG4gKiBDYWxsZWQgZnJvbSBHbyB3aGVuIGEgZmlsZSBkcmFnIGVudGVycyB0aGUgd2luZG93IG9uIExpbnV4L21hY09TLlxuICovXG5mdW5jdGlvbiBoYW5kbGVEcmFnRW50ZXIoKTogdm9pZCB7XG4gICAgLy8gQ2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IGFjdGl2YXRlIGRyYWcgc3RhdGVcbiAgICB9XG4gICAgbmF0aXZlRHJhZ0FjdGl2ZSA9IHRydWU7XG59XG5cbi8qKlxuICogQ2FsbGVkIGZyb20gR28gd2hlbiBhIGZpbGUgZHJhZyBsZWF2ZXMgdGhlIHdpbmRvdyBvbiBMaW51eC9tYWNPUy5cbiAqL1xuZnVuY3Rpb24gaGFuZGxlRHJhZ0xlYXZlKCk6IHZvaWQge1xuICAgIGNsZWFudXBOYXRpdmVEcmFnKCk7XG59XG5cbi8qKlxuICogQ2FsbGVkIGZyb20gR28gZHVyaW5nIGZpbGUgZHJhZyB0byB1cGRhdGUgaG92ZXIgc3RhdGUgb24gTGludXgvbWFjT1MuXG4gKiBAcGFyYW0geCAtIFggY29vcmRpbmF0ZSBpbiBDU1MgcGl4ZWxzXG4gKiBAcGFyYW0geSAtIFkgY29vcmRpbmF0ZSBpbiBDU1MgcGl4ZWxzXG4gKi9cbmZ1bmN0aW9uIGhhbmRsZURyYWdPdmVyKHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XG4gICAgaWYgKCFuYXRpdmVEcmFnQWN0aXZlKSByZXR1cm47XG4gICAgXG4gICAgLy8gQ2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xuICAgIH1cbiAgICBcbiAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KTtcbiAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XG4gICAgXG4gICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0ICYmIGN1cnJlbnREcm9wVGFyZ2V0ICE9PSBkcm9wVGFyZ2V0KSB7XG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICB9XG4gICAgXG4gICAgaWYgKGRyb3BUYXJnZXQpIHtcbiAgICAgICAgZHJvcFRhcmdldC5jbGFzc0xpc3QuYWRkKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcbiAgICB9IGVsc2Uge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XG4gICAgfVxufVxuXG5cblxuLy8gRXhwb3J0IHRoZSBoYW5kbGVycyBmb3IgdXNlIGJ5IEdvIHZpYSBpbmRleC50c1xuZXhwb3J0IHsgaGFuZGxlRHJhZ0VudGVyLCBoYW5kbGVEcmFnTGVhdmUsIGhhbmRsZURyYWdPdmVyIH07XG5cbi8qKlxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgcG9zaXRpb24gb2YgYSB3aW5kb3cuXG4gKi9cbmludGVyZmFjZSBQb3NpdGlvbiB7XG4gICAgLyoqIFRoZSBob3Jpem9udGFsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXG4gICAgeDogbnVtYmVyO1xuICAgIC8qKiBUaGUgdmVydGljYWwgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy4gKi9cbiAgICB5OiBudW1iZXI7XG59XG5cbi8qKlxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgc2l6ZSBvZiBhIHdpbmRvdy5cbiAqL1xuaW50ZXJmYWNlIFNpemUge1xuICAgIC8qKiBUaGUgd2lkdGggb2YgdGhlIHdpbmRvdy4gKi9cbiAgICB3aWR0aDogbnVtYmVyO1xuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuICovXG4gICAgaGVpZ2h0OiBudW1iZXI7XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBjYWxsZXJTeW0gPSBTeW1ib2woXCJjYWxsZXJcIik7XG5cbmNsYXNzIFdpbmRvdyB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgcHJpdmF0ZSBbY2FsbGVyU3ltXTogKG1lc3NhZ2U6IG51bWJlciwgYXJncz86IGFueSkgPT4gUHJvbWlzZTxhbnk+O1xuXG4gICAgLyoqXG4gICAgICogSW5pdGlhbGlzZXMgYSB3aW5kb3cgb2JqZWN0IHdpdGggdGhlIHNwZWNpZmllZCBuYW1lLlxuICAgICAqXG4gICAgICogQHByaXZhdGVcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB0YXJnZXQgd2luZG93LlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG5hbWU6IHN0cmluZyA9ICcnKSB7XG4gICAgICAgIHRoaXNbY2FsbGVyU3ltXSA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuV2luZG93LCBuYW1lKVxuXG4gICAgICAgIC8vIGJpbmQgaW5zdGFuY2UgbWV0aG9kIHRvIG1ha2UgdGhlbSBlYXNpbHkgdXNhYmxlIGluIGV2ZW50IGhhbmRsZXJzXG4gICAgICAgIGZvciAoY29uc3QgbWV0aG9kIG9mIE9iamVjdC5nZXRPd25Qcm9wZXJ0eU5hbWVzKFdpbmRvdy5wcm90b3R5cGUpKSB7XG4gICAgICAgICAgICBpZiAoXG4gICAgICAgICAgICAgICAgbWV0aG9kICE9PSBcImNvbnN0cnVjdG9yXCJcbiAgICAgICAgICAgICAgICAmJiB0eXBlb2YgKHRoaXMgYXMgYW55KVttZXRob2RdID09PSBcImZ1bmN0aW9uXCJcbiAgICAgICAgICAgICkge1xuICAgICAgICAgICAgICAgICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9ICh0aGlzIGFzIGFueSlbbWV0aG9kXS5iaW5kKHRoaXMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogR2V0cyB0aGUgc3BlY2lmaWVkIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBnZXQuXG4gICAgICogQHJldHVybnMgVGhlIGNvcnJlc3BvbmRpbmcgd2luZG93IG9iamVjdC5cbiAgICAgKi9cbiAgICBHZXQobmFtZTogc3RyaW5nKTogV2luZG93IHtcbiAgICAgICAgcmV0dXJuIG5ldyBXaW5kb3cobmFtZSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUG9zaXRpb24oKTogUHJvbWlzZTxQb3NpdGlvbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFBvc2l0aW9uTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cgb24gdGhlIHNjcmVlbi5cbiAgICAgKi9cbiAgICBDZW50ZXIoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2VudGVyTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDbG9zZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBDbG9zZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShDbG9zZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqL1xuICAgIERpc2FibGVTaXplQ29uc3RyYWludHMoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRGlzYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRW5hYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXG4gICAgICovXG4gICAgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9jdXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIEZvY3VzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvY3VzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBGb3JjZXMgdGhlIHdpbmRvdyB0byByZWxvYWQgdGhlIHBhZ2UgYXNzZXRzLlxuICAgICAqL1xuICAgIEZvcmNlUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvcmNlUmVsb2FkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTd2l0Y2hlcyB0aGUgd2luZG93IHRvIGZ1bGxzY3JlZW4gbW9kZS5cbiAgICAgKi9cbiAgICBGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgc2NyZWVuIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG9uLlxuICAgICAqL1xuICAgIEdldFNjcmVlbigpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEdldFNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgY3VycmVudCB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB6b29tIGxldmVsLlxuICAgICAqL1xuICAgIEdldFpvb20oKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRab29tTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIEhlaWdodCgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhlaWdodE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSGlkZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhpZGVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZvY3VzZWQuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZvY3VzZWQuXG4gICAgICovXG4gICAgSXNGb2N1c2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRm9jdXNlZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZnVsbHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZnVsbHNjcmVlbi5cbiAgICAgKi9cbiAgICBJc0Z1bGxzY3JlZW4oKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtYXhpbWlzZWQuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1heGltaXNlZC5cbiAgICAgKi9cbiAgICBJc01heGltaXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc01heGltaXNlZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWluaW1pc2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtaW5pbWlzZWQuXG4gICAgICovXG4gICAgSXNNaW5pbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNaW5pbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIE1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNaW5pbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBNaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNaW5pbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBOYW1lKCk6IFByb21pc2U8c3RyaW5nPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTmFtZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogT3BlbnMgdGhlIGRldmVsb3BtZW50IHRvb2xzIHBhbmUuXG4gICAgICovXG4gICAgT3BlbkRldlRvb2xzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE9wZW5EZXZUb29sc01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBSZWxhdGl2ZVBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZWxvYWRzIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKi9cbiAgICBSZWxvYWQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsb2FkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IHJlc2l6YWJsZS5cbiAgICAgKi9cbiAgICBSZXNpemFibGUoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzaXphYmxlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXN0b3JlcyB0aGUgd2luZG93IHRvIGl0cyBwcmV2aW91cyBzdGF0ZSBpZiBpdCB3YXMgcHJldmlvdXNseSBtaW5pbWlzZWQsIG1heGltaXNlZCBvciBmdWxsc2NyZWVuLlxuICAgICAqL1xuICAgIFJlc3RvcmUoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVzdG9yZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRQb3NpdGlvbih4OiBudW1iZXIsIHk6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFBvc2l0aW9uTWV0aG9kLCB7IHgsIHkgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gYWx3YXlzT25Ub3AgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIHN0YXkgb24gdG9wLlxuICAgICAqL1xuICAgIFNldEFsd2F5c09uVG9wKGFsd2F5c09uVG9wOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0QWx3YXlzT25Ub3BNZXRob2QsIHsgYWx3YXlzT25Ub3AgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSByIC0gVGhlIGRlc2lyZWQgcmVkIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICogQHBhcmFtIGcgLSBUaGUgZGVzaXJlZCBncmVlbiBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBiIC0gVGhlIGRlc2lyZWQgYmx1ZSBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBhIC0gVGhlIGRlc2lyZWQgYWxwaGEgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKi9cbiAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyKHI6IG51bWJlciwgZzogbnVtYmVyLCBiOiBudW1iZXIsIGE6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QsIHsgciwgZywgYiwgYSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIHRoZSB3aW5kb3cgZnJhbWUgYW5kIHRpdGxlIGJhci5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBmcmFtZWxlc3MgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGZyYW1lbGVzcy5cbiAgICAgKi9cbiAgICBTZXRGcmFtZWxlc3MoZnJhbWVsZXNzOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0RnJhbWVsZXNzTWV0aG9kLCB7IGZyYW1lbGVzcyB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEaXNhYmxlcyB0aGUgc3lzdGVtIGZ1bGxzY3JlZW4gYnV0dG9uLlxuICAgICAqXG4gICAgICogQHBhcmFtIGVuYWJsZWQgLSBXaGV0aGVyIHRoZSBmdWxsc2NyZWVuIGJ1dHRvbiBzaG91bGQgYmUgZW5hYmxlZC5cbiAgICAgKi9cbiAgICBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZChlbmFibGVkOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWRNZXRob2QsIHsgZW5hYmxlZCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtYXhpbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1heGltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWF4aW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRNYXhTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWF4U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIG1pbmltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHBhcmFtIHdpZHRoIC0gVGhlIGRlc2lyZWQgbWluaW11bSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtaW5pbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldE1pblNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRNaW5TaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxuICAgICAqXG4gICAgICogQHBhcmFtIHggLSBUaGUgZGVzaXJlZCBob3Jpem9udGFsIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFJlbGF0aXZlUG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kLCB7IHgsIHkgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB3aGV0aGVyIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHBhcmFtIHJlc2l6YWJsZSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgcmVzaXphYmxlLlxuICAgICAqL1xuICAgIFNldFJlc2l6YWJsZShyZXNpemFibGU6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRSZXNpemFibGVNZXRob2QsIHsgcmVzaXphYmxlIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIGhlaWdodCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgdGl0bGUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB0aXRsZSAtIFRoZSBkZXNpcmVkIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0VGl0bGUodGl0bGU6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFRpdGxlTWV0aG9kLCB7IHRpdGxlIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB6b29tIC0gVGhlIGRlc2lyZWQgem9vbSBsZXZlbC5cbiAgICAgKi9cbiAgICBTZXRab29tKHpvb206IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFpvb21NZXRob2QsIHsgem9vbSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTaG93cyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2hvd01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgc2l6ZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNpemUoKTogUHJvbWlzZTxTaXplPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2l6ZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnVsbHNjcmVlbiBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZUZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gbWF4aW1pc2VkIGFuZCBub3JtYWwuXG4gICAgICovXG4gICAgVG9nZ2xlTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlTWF4aW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZyYW1lbGVzcyBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZUZyYW1lbGVzcygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVGcmFtZWxlc3NNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLWZ1bGxzY3JlZW5zIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5GdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWF4aW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5NYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1taW5pbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBVbk1pbmltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFdpZHRoKCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oV2lkdGhNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFpvb21zIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgWm9vbSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBJbmNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tSW4oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbUluTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBEZWNyZWFzZXMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tT3V0KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21PdXRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlc2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxuICAgICAqL1xuICAgIFpvb21SZXNldCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tUmVzZXRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEhhbmRsZXMgZmlsZSBkcm9wcyBvcmlnaW5hdGluZyBmcm9tIHBsYXRmb3JtLXNwZWNpZmljIGNvZGUgKGUuZy4sIG1hY09TL0xpbnV4IG5hdGl2ZSBkcmFnLWFuZC1kcm9wKS5cbiAgICAgKiBHYXRoZXJzIGluZm9ybWF0aW9uIGFib3V0IHRoZSBkcm9wIHRhcmdldCBlbGVtZW50IGFuZCBzZW5kcyBpdCBiYWNrIHRvIHRoZSBHbyBiYWNrZW5kLlxuICAgICAqXG4gICAgICogQHBhcmFtIGZpbGVuYW1lcyAtIEFuIGFycmF5IG9mIGZpbGUgcGF0aHMgKHN0cmluZ3MpIHRoYXQgd2VyZSBkcm9wcGVkLlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIHgtY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudCwgaW4gbG9naWNhbCAoQ1NTKSBwaXhlbHMgcmVsYXRpdmUgdG8gdGhlIHdlYnZpZXcuXG4gICAgICogQHBhcmFtIHkgLSBUaGUgeS1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50LCBpbiBsb2dpY2FsIChDU1MpIHBpeGVscyByZWxhdGl2ZSB0byB0aGUgd2Vidmlldy5cbiAgICAgKi9cbiAgICBIYW5kbGVQbGF0Zm9ybUZpbGVEcm9wKGZpbGVuYW1lczogc3RyaW5nW10sIHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XG4gICAgICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGlnbm9yZSB0aGUgZHJvcFxuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICBjb25zdCBlbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludCh4LCB5KTtcbiAgICAgICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KGVsZW1lbnQpO1xuXG4gICAgICAgIGlmICghZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgLy8gRHJvcCB3YXMgbm90IG9uIGEgZGVzaWduYXRlZCBkcm9wIHRhcmdldCAtIGlnbm9yZVxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgZWxlbWVudERldGFpbHMgPSB7XG4gICAgICAgICAgICBpZDogZHJvcFRhcmdldC5pZCxcbiAgICAgICAgICAgIGNsYXNzTGlzdDogQXJyYXkuZnJvbShkcm9wVGFyZ2V0LmNsYXNzTGlzdCksXG4gICAgICAgICAgICBhdHRyaWJ1dGVzOiB7fSBhcyB7IFtrZXk6IHN0cmluZ106IHN0cmluZyB9LFxuICAgICAgICB9O1xuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IGRyb3BUYXJnZXQuYXR0cmlidXRlcy5sZW5ndGg7IGkrKykge1xuICAgICAgICAgICAgY29uc3QgYXR0ciA9IGRyb3BUYXJnZXQuYXR0cmlidXRlc1tpXTtcbiAgICAgICAgICAgIGVsZW1lbnREZXRhaWxzLmF0dHJpYnV0ZXNbYXR0ci5uYW1lXSA9IGF0dHIudmFsdWU7XG4gICAgICAgIH1cblxuICAgICAgICBjb25zdCBwYXlsb2FkID0ge1xuICAgICAgICAgICAgZmlsZW5hbWVzLFxuICAgICAgICAgICAgeCxcbiAgICAgICAgICAgIHksXG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscyxcbiAgICAgICAgfTtcblxuICAgICAgICB0aGlzW2NhbGxlclN5bV0oRmlsZXNEcm9wcGVkLCBwYXlsb2FkKTtcbiAgICAgICAgXG4gICAgICAgIC8vIENsZWFuIHVwIG5hdGl2ZSBkcmFnIHN0YXRlIGFmdGVyIGRyb3BcbiAgICAgICAgY2xlYW51cE5hdGl2ZURyYWcoKTtcbiAgICB9XG4gIFxuICAgIC8qKlxuICAgICAqIE1vdmVzIHRoZSB3aW5kb3cgdG8gdGhlIGNlbnRlciBvZiB0aGUgc3BlY2lmaWVkIHNjcmVlbidzIHdvcmsgYXJlYS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBzY3JlZW5JRCAtIFRoZSBJRCBvZiB0aGUgdGFyZ2V0IHNjcmVlbi5cbiAgICAgKi9cbiAgICBTZXRTY3JlZW4oc2NyZWVuSUQ6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFNjcmVlbk1ldGhvZCwgeyBzY3JlZW5JRCB9KTtcbiAgICB9XG5cbiAgICAvKiBUcmlnZ2VycyBXaW5kb3dzIDExIFNuYXAgQXNzaXN0IGZlYXR1cmUgKFdpbmRvd3Mgb25seSkuXG4gICAgICogVGhpcyBpcyBlcXVpdmFsZW50IHRvIHByZXNzaW5nIFdpbitaIGFuZCBzaG93cyBzbmFwIGxheW91dCBvcHRpb25zLlxuICAgICAqL1xuICAgIFNuYXBBc3Npc3QoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU25hcEFzc2lzdE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogT3BlbnMgdGhlIHByaW50IGRpYWxvZyBmb3IgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBQcmludCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQcmludE1ldGhvZCk7XG4gICAgfVxufVxuXG4vKipcbiAqIFRoZSB3aW5kb3cgd2l0aGluIHdoaWNoIHRoZSBzY3JpcHQgaXMgcnVubmluZy5cbiAqL1xuY29uc3QgdGhpc1dpbmRvdyA9IG5ldyBXaW5kb3coJycpO1xuXG4vKipcbiAqIFNldHMgdXAgZ2xvYmFsIGRyYWcgYW5kIGRyb3AgZXZlbnQgbGlzdGVuZXJzIGZvciBmaWxlIGRyb3BzLlxuICogSGFuZGxlcyB2aXN1YWwgZmVlZGJhY2sgKGhvdmVyIHN0YXRlKSBhbmQgZmlsZSBkcm9wIHByb2Nlc3NpbmcuXG4gKi9cbmZ1bmN0aW9uIHNldHVwRHJvcFRhcmdldExpc3RlbmVycygpIHtcbiAgICBjb25zdCBkb2NFbGVtZW50ID0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50O1xuICAgIGxldCBkcmFnRW50ZXJDb3VudGVyID0gMDtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2VudGVyJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyYWdzLCBsZXQgb3RoZXIgZHJhZ3MgcGFzcyB0aHJvdWdoXG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7IC8vIFNob3cgXCJubyBkcm9wXCIgY3Vyc29yXG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xuICAgICAgICB9XG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIrKztcbiAgICAgICAgXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XG5cbiAgICAgICAgLy8gVXBkYXRlIGhvdmVyIHN0YXRlXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCAoZXZlbnQpID0+IHtcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJhZ3NcbiAgICAgICAgfVxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJzsgLy8gU2hvdyBcIm5vIGRyb3BcIiBjdXJzb3JcbiAgICAgICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgZG9uJ3Qgc2hvdyBob3ZlciBlZmZlY3RzXG4gICAgICAgIH1cbiAgICAgICAgXG4gICAgICAgIC8vIFVwZGF0ZSBkcm9wIHRhcmdldCBhcyBjdXJzb3IgbW92ZXNcbiAgICAgICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSk7XG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudCh0YXJnZXRFbGVtZW50KTtcbiAgICAgICAgXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICBpZiAoZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgaWYgKCFkcm9wVGFyZ2V0LmNsYXNzTGlzdC5jb250YWlucyhEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpKSB7XG4gICAgICAgICAgICAgICAgZHJvcFRhcmdldC5jbGFzc0xpc3QuYWRkKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdjb3B5JztcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnO1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICAvLyBPbiBMaW51eC9XZWJLaXRHVEsgYW5kIG1hY09TLCBkcmFnbGVhdmUgZmlyZXMgaW1tZWRpYXRlbHkgd2l0aCByZWxhdGVkVGFyZ2V0PW51bGwgd2hlbiBuYXRpdmVcbiAgICAgICAgLy8gZHJhZyBoYW5kbGluZyBpcyBpbnZvbHZlZC4gSWdub3JlIHRoZXNlIHNwdXJpb3VzIGV2ZW50cyAtIHdlJ2xsIGNsZWFuIHVwIG9uIGRyb3AgaW5zdGVhZC5cbiAgICAgICAgaWYgKGV2ZW50LnJlbGF0ZWRUYXJnZXQgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgZHJhZ0VudGVyQ291bnRlci0tO1xuICAgICAgICBcbiAgICAgICAgaWYgKGRyYWdFbnRlckNvdW50ZXIgPT09IDAgfHwgXG4gICAgICAgICAgICAoY3VycmVudERyb3BUYXJnZXQgJiYgIWN1cnJlbnREcm9wVGFyZ2V0LmNvbnRhaW5zKGV2ZW50LnJlbGF0ZWRUYXJnZXQgYXMgTm9kZSkpKSB7XG4gICAgICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJvcHNcbiAgICAgICAgfVxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7XG4gICAgICAgIFxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGhhbmRsZSBmaWxlIGRyb3BzIHZpYSBKYXZhU2NyaXB0XG4gICAgICAgIC8vIE9uIG1hY09TL0xpbnV4LCBuYXRpdmUgY29kZSB3aWxsIGNhbGwgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcFxuICAgICAgICBpZiAoY2FuUmVzb2x2ZUZpbGVQYXRocygpKSB7XG4gICAgICAgICAgICBjb25zdCBmaWxlczogRmlsZVtdID0gW107XG4gICAgICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XG4gICAgICAgICAgICAgICAgZm9yIChjb25zdCBpdGVtIG9mIGV2ZW50LmRhdGFUcmFuc2Zlci5pdGVtcykge1xuICAgICAgICAgICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IGZpbGUgPSBpdGVtLmdldEFzRmlsZSgpO1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGZpbGUpIGZpbGVzLnB1c2goZmlsZSk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9IGVsc2UgaWYgKGV2ZW50LmRhdGFUcmFuc2Zlci5maWxlcykge1xuICAgICAgICAgICAgICAgIGZvciAoY29uc3QgZmlsZSBvZiBldmVudC5kYXRhVHJhbnNmZXIuZmlsZXMpIHtcbiAgICAgICAgICAgICAgICAgICAgZmlsZXMucHVzaChmaWxlKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBcbiAgICAgICAgICAgIGlmIChmaWxlcy5sZW5ndGggPiAwKSB7XG4gICAgICAgICAgICAgICAgcmVzb2x2ZUZpbGVQYXRocyhldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBmaWxlcyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG59XG5cbi8vIEluaXRpYWxpemUgbGlzdGVuZXJzIHdoZW4gdGhlIHNjcmlwdCBsb2Fkc1xuaWYgKHR5cGVvZiB3aW5kb3cgIT09IFwidW5kZWZpbmVkXCIgJiYgdHlwZW9mIGRvY3VtZW50ICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgc2V0dXBEcm9wVGFyZ2V0TGlzdGVuZXJzKCk7XG59XG5cbmV4cG9ydCBkZWZhdWx0IHRoaXNXaW5kb3c7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XG5cbi8vIE5PVEU6IHRoZSBmb2xsb3dpbmcgbWV0aG9kcyBNVVNUIGJlIGltcG9ydGVkIGV4cGxpY2l0bHkgYmVjYXVzZSBvZiBob3cgZXNidWlsZCBpbmplY3Rpb24gd29ya3NcbmltcG9ydCB7IEVuYWJsZSBhcyBFbmFibGVXTUwgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XG5pbXBvcnQgeyBkZWJ1Z0xvZyB9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy91dGlsc1wiO1xuXG53aW5kb3cud2FpbHMgPSBSdW50aW1lO1xuRW5hYmxlV01MKCk7XG5cbmlmIChERUJVRykge1xuICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIilcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlN5c3RlbSk7XG5cbmNvbnN0IFN5c3RlbUlzRGFya01vZGUgPSAwO1xuY29uc3QgU3lzdGVtRW52aXJvbm1lbnQgPSAxO1xuY29uc3QgU3lzdGVtQ2FwYWJpbGl0aWVzID0gMjtcblxuY29uc3QgX2ludm9rZSA9IChmdW5jdGlvbiAoKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgLy8gV2luZG93cyBXZWJWaWV3MlxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2UpIHtcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcpO1xuICAgICAgICB9XG4gICAgICAgIC8vIG1hY09TL2lPUyBXS1dlYlZpZXdcbiAgICAgICAgZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndlYmtpdD8ubWVzc2FnZUhhbmRsZXJzPy5bJ2V4dGVybmFsJ10/LnBvc3RNZXNzYWdlKSB7XG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10ucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXSk7XG4gICAgICAgIH1cbiAgICAgICAgLy8gQW5kcm9pZCBXZWJWaWV3IC0gdXNlcyBhZGRKYXZhc2NyaXB0SW50ZXJmYWNlIHdoaWNoIGV4cG9zZXMgd2luZG93LndhaWxzLmludm9rZVxuICAgICAgICBlbHNlIGlmICgod2luZG93IGFzIGFueSkud2FpbHM/Lmludm9rZSkge1xuICAgICAgICAgICAgcmV0dXJuIChtc2c6IGFueSkgPT4gKHdpbmRvdyBhcyBhbnkpLndhaWxzLmludm9rZSh0eXBlb2YgbXNnID09PSAnc3RyaW5nJyA/IG1zZyA6IEpTT04uc3RyaW5naWZ5KG1zZykpO1xuICAgICAgICB9XG4gICAgfSBjYXRjaChlKSB7fVxuXG4gICAgY29uc29sZS53YXJuKCdcXG4lY1x1MjZBMFx1RkUwRiBCcm93c2VyIEVudmlyb25tZW50IERldGVjdGVkICVjXFxuXFxuJWNPbmx5IFVJIHByZXZpZXdzIGFyZSBhdmFpbGFibGUgaW4gdGhlIGJyb3dzZXIuIEZvciBmdWxsIGZ1bmN0aW9uYWxpdHksIHBsZWFzZSBydW4gdGhlIGFwcGxpY2F0aW9uIGluIGRlc2t0b3AgbW9kZS5cXG5Nb3JlIGluZm9ybWF0aW9uIGF0OiBodHRwczovL3YzLndhaWxzLmlvL2xlYXJuL2J1aWxkLyN1c2luZy1hLWJyb3dzZXItZm9yLWRldmVsb3BtZW50XFxuJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNmZmZmZmY7IGNvbG9yOiAjMDAwMDAwOyBmb250LXdlaWdodDogYm9sZDsgcGFkZGluZzogNHB4IDhweDsgYm9yZGVyLXJhZGl1czogNHB4OyBib3JkZXI6IDJweCBzb2xpZCAjMDAwMDAwOycsXG4gICAgICAgICdiYWNrZ3JvdW5kOiB0cmFuc3BhcmVudDsnLFxuICAgICAgICAnY29sb3I6ICNmZmZmZmY7IGZvbnQtc3R5bGU6IGl0YWxpYzsgZm9udC13ZWlnaHQ6IGJvbGQ7Jyk7XG4gICAgcmV0dXJuIG51bGw7XG59KSgpO1xuXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZzogYW55KTogdm9pZCB7XG4gICAgX2ludm9rZT8uKG1zZyk7XG59XG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSBzeXN0ZW0gZGFyayBtb2RlIHN0YXR1cy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtSXNEYXJrTW9kZSk7XG59XG5cbi8qKlxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIENhcGFiaWxpdGllcygpOiBQcm9taXNlPFJlY29yZDxzdHJpbmcsIGFueT4+IHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1DYXBhYmlsaXRpZXMpO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIE9TSW5mbyB7XG4gICAgLyoqIFRoZSBicmFuZGluZyBvZiB0aGUgT1MuICovXG4gICAgQnJhbmRpbmc6IHN0cmluZztcbiAgICAvKiogVGhlIElEIG9mIHRoZSBPUy4gKi9cbiAgICBJRDogc3RyaW5nO1xuICAgIC8qKiBUaGUgbmFtZSBvZiB0aGUgT1MuICovXG4gICAgTmFtZTogc3RyaW5nO1xuICAgIC8qKiBUaGUgdmVyc2lvbiBvZiB0aGUgT1MuICovXG4gICAgVmVyc2lvbjogc3RyaW5nO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEVudmlyb25tZW50SW5mbyB7XG4gICAgLyoqIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS4gKi9cbiAgICBBcmNoOiBzdHJpbmc7XG4gICAgLyoqIFRydWUgaWYgdGhlIGFwcGxpY2F0aW9uIGlzIHJ1bm5pbmcgaW4gZGVidWcgbW9kZSwgb3RoZXJ3aXNlIGZhbHNlLiAqL1xuICAgIERlYnVnOiBib29sZWFuO1xuICAgIC8qKiBUaGUgb3BlcmF0aW5nIHN5c3RlbSBpbiB1c2UuICovXG4gICAgT1M6IHN0cmluZztcbiAgICAvKiogRGV0YWlscyBvZiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cbiAgICBPU0luZm86IE9TSW5mbztcbiAgICAvKiogQWRkaXRpb25hbCBwbGF0Zm9ybSBpbmZvcm1hdGlvbi4gKi9cbiAgICBQbGF0Zm9ybUluZm86IFJlY29yZDxzdHJpbmcsIGFueT47XG59XG5cbi8qKlxuICogUmV0cmlldmVzIGVudmlyb25tZW50IGRldGFpbHMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpOiBQcm9taXNlPEVudmlyb25tZW50SW5mbz4ge1xuICAgIHJldHVybiBjYWxsKFN5c3RlbUVudmlyb25tZW50KTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxuICpcbiAqIEByZXR1cm4gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc1dpbmRvd3MoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJ3aW5kb3dzXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXG4gKlxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJsaW51eFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBhIG1hY09TIG9wZXJhdGluZyBzeXN0ZW0uXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwiZGFyd2luXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgaU9TLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgaU9TLCBvdGhlcndpc2UgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0lPUygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TID09PSBcImlvc1wiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIEFuZHJvaWQuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBBbmRyb2lkLCBvdGhlcndpc2UgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FuZHJvaWQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJhbmRyb2lkXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBhcHAgaXMgcnVubmluZyBvbiBhIG1vYmlsZSBPUyAoaU9TIG9yIEFuZHJvaWQpLlxuICpcbiAqIEByZXR1cm5zIFRydWUgb24gaU9TIG9yIEFuZHJvaWQsIG90aGVyd2lzZSBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTW9iaWxlKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBJc0lPUygpIHx8IElzQW5kcm9pZCgpO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgYXBwIGlzIHJ1bm5pbmcgb24gYSBkZXNrdG9wIE9TIChtYWNPUywgV2luZG93cyBvciBMaW51eCkuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBvbiBtYWNPUywgV2luZG93cyBvciBMaW51eCwgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNEZXNrdG9wKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBJc01hYygpIHx8IElzV2luZG93cygpIHx8IElzTGludXgoKTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0FNRDY0KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhbWQ2NFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBUk0oKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFybVwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUuXG4gKlxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUsIG90aGVyd2lzZSByZXR1cm5zIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBUk02NCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYXJtNjRcIjtcbn1cblxuLyoqXG4gKiBSZXBvcnRzIHdoZXRoZXIgdGhlIGFwcCBpcyBiZWluZyBydW4gaW4gZGVidWcgbW9kZS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RlYnVnKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiBCb29sZWFuKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5EZWJ1Zyk7XG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBJc0RlYnVnIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5cbi8vIHNldHVwXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcbn1cblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xuXG5jb25zdCBDb250ZXh0TWVudU9wZW4gPSAwO1xuXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQ6IHN0cmluZywgeDogbnVtYmVyLCB5OiBudW1iZXIsIGRhdGE6IGFueSk6IHZvaWQge1xuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XG5cbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIikudHJpbSgpO1xuXG4gICAgaWYgKGN1c3RvbUNvbnRleHRNZW51KSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIGNvbnN0IGRhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGRhdGEpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQsIHRhcmdldCk7XG4gICAgfVxufVxuXG5cbi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQ6IE1vdXNlRXZlbnQsIHRhcmdldDogSFRNTEVsZW1lbnQpIHtcbiAgICAvLyBEZWJ1ZyBidWlsZHMgYWx3YXlzIHNob3cgdGhlIG1lbnVcbiAgICBpZiAoSXNEZWJ1ZygpKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgc3dpdGNoICh3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpKSB7XG4gICAgICAgIGNhc2UgJ3Nob3cnOlxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBjYXNlICdoaWRlJzpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcbiAgICBpZiAodGFyZ2V0LmlzQ29udGVudEVkaXRhYmxlKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXG4gICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IHNlbGVjdGlvbiAmJiBzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwO1xuICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xuICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gdGFyZ2V0KSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiB0YWcgaXMgaW5wdXQgb3IgdGV4dGFyZWEuXG4gICAgaWYgKHRhcmdldCBpbnN0YW5jZW9mIEhUTUxJbnB1dEVsZW1lbnQgfHwgdGFyZ2V0IGluc3RhbmNlb2YgSFRNTFRleHRBcmVhRWxlbWVudCkge1xuICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghdGFyZ2V0LnJlYWRPbmx5ICYmICF0YXJnZXQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBSZXRyaWV2ZXMgdGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleSBmcm9tIHRoZSBmbGFnIG1hcC5cbiAqXG4gKiBAcGFyYW0ga2V5IC0gVGhlIGtleSB0byByZXRyaWV2ZSB0aGUgdmFsdWUgZm9yLlxuICogQHJldHVybiBUaGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5LlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXk6IHN0cmluZyk6IGFueSB7XG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZmxhZ3Nba2V5XTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleSArIFwiJzogXCIgKyBlLCB7IGNhdXNlOiBlIH0pO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgaW52b2tlLCBJc1dpbmRvd3MsIElzTGludXgsIElzTWFjIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyBHZXRGbGFnIH0gZnJvbSBcIi4vZmxhZ3MuanNcIjtcbmltcG9ydCB7IGNhblRyYWNrQnV0dG9ucywgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuXG4vLyBTZXR1cFxubGV0IGNhbkRyYWcgPSBmYWxzZTtcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xuXG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XG5sZXQgcmVzaXppbmcgPSBmYWxzZTtcbmxldCByZXNpemVFZGdlOiBzdHJpbmcgPSBcIlwiO1xubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcblxubGV0IGJ1dHRvbnMgPSAwO1xuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxubGV0IGJ1dHRvbnNUcmFja2VkID0gZmFsc2U7XG5cbmlmIChoYXNET00pIHtcbiAgICBidXR0b25zVHJhY2tlZCA9IGNhblRyYWNrQnV0dG9ucygpO1xuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xuICAgIHdpbmRvdy5fd2FpbHMuc2V0UmVzaXphYmxlID0gKHZhbHVlOiBib29sZWFuKTogdm9pZCA9PiB7XG4gICAgICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xuICAgICAgICBpZiAoIXJlc2l6YWJsZSkge1xuICAgICAgICAgICAgLy8gU3RvcCByZXNpemluZyBpZiBpbiBwcm9ncmVzcy5cbiAgICAgICAgICAgIGNhblJlc2l6ZSA9IHJlc2l6aW5nID0gZmFsc2U7XG4gICAgICAgICAgICBzZXRSZXNpemUoKTtcbiAgICAgICAgfVxuICAgIH07XG59XG5cbi8vIERlZmVyIGF0dGFjaGluZyBtb3VzZSBsaXN0ZW5lcnMgdW50aWwgd2Uga25vdyB3ZSdyZSBub3Qgb24gbW9iaWxlLlxubGV0IGRyYWdJbml0RG9uZSA9IGZhbHNlO1xuZnVuY3Rpb24gaXNNb2JpbGUoKTogYm9vbGVhbiB7XG4gICAgY29uc3Qgb3MgPSAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1M7XG4gICAgaWYgKG9zID09PSBcImlvc1wiIHx8IG9zID09PSBcImFuZHJvaWRcIikgcmV0dXJuIHRydWU7XG4gICAgLy8gRmFsbGJhY2sgaGV1cmlzdGljIGlmIGVudmlyb25tZW50IG5vdCB5ZXQgc2V0XG4gICAgY29uc3QgdWEgPSBuYXZpZ2F0b3IudXNlckFnZW50IHx8IG5hdmlnYXRvci52ZW5kb3IgfHwgKHdpbmRvdyBhcyBhbnkpLm9wZXJhIHx8IFwiXCI7XG4gICAgcmV0dXJuIC9hbmRyb2lkfGlwaG9uZXxpcGFkfGlwb2R8aWVtb2JpbGV8d3BkZXNrdG9wL2kudGVzdCh1YSk7XG59XG5mdW5jdGlvbiB0cnlJbml0RHJhZ0hhbmRsZXJzKCk6IHZvaWQge1xuICAgIGlmIChkcmFnSW5pdERvbmUpIHJldHVybjtcbiAgICBpZiAoaXNNb2JpbGUoKSkgcmV0dXJuO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICBmb3IgKGNvbnN0IGV2IG9mIFsnY2xpY2snLCAnY29udGV4dG1lbnUnLCAnZGJsY2xpY2snXSkge1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihldiwgc3VwcHJlc3NFdmVudCwgeyBjYXB0dXJlOiB0cnVlIH0pO1xuICAgIH1cbiAgICBkcmFnSW5pdERvbmUgPSB0cnVlO1xufVxuaWYgKGhhc0RPTSkge1xuICAgIC8vIEF0dGVtcHQgaW1tZWRpYXRlIGluaXQgKGluIGNhc2UgZW52aXJvbm1lbnQgYWxyZWFkeSBwcmVzZW50KVxuICAgIHRyeUluaXREcmFnSGFuZGxlcnMoKTtcbiAgICAvLyBBbHNvIGF0dGVtcHQgb24gRE9NIHJlYWR5XG4gICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIHRyeUluaXREcmFnSGFuZGxlcnMsIHsgb25jZTogdHJ1ZSB9KTtcbiAgICAvLyBBcyBhIGxhc3QgcmVzb3J0LCBwb2xsIGZvciBlbnZpcm9ubWVudCBmb3IgYSBzaG9ydCBwZXJpb2RcbiAgICBsZXQgZHJhZ0VudlBvbGxzID0gMDtcbiAgICBjb25zdCBkcmFnRW52UG9sbCA9IHdpbmRvdy5zZXRJbnRlcnZhbCgoKSA9PiB7XG4gICAgICAgIGlmIChkcmFnSW5pdERvbmUpIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyByZXR1cm47IH1cbiAgICAgICAgdHJ5SW5pdERyYWdIYW5kbGVycygpO1xuICAgICAgICBpZiAoKytkcmFnRW52UG9sbHMgPiAxMDApIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyB9XG4gICAgfSwgNTApO1xufVxuXG5mdW5jdGlvbiBzdXBwcmVzc0V2ZW50KGV2ZW50OiBFdmVudCkge1xuICAgIGlmIChcbiAgICAgICAgZXZlbnQudHlwZSA9PT0gXCJkYmxjbGlja1wiXG4gICAgICAgICYmIElzTWFjKClcbiAgICAgICAgJiYgd2luZG93ID09PSB3aW5kb3cudG9wXG4gICAgICAgICYmIGlzRHJhZ2dhYmxlRXZlbnQoZXZlbnQgYXMgTW91c2VFdmVudClcbiAgICApIHtcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBpbnZva2UoXCJ3YWlsczpkcmFnOmRvdWJsZWNsaWNrXCIpO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gU3VwcHJlc3MgY2xpY2sgZXZlbnRzIHdoaWxlIHJlc2l6aW5nIG9yIGRyYWdnaW5nLlxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBpc0RyYWdnYWJsZUV2ZW50KGV2ZW50OiBNb3VzZUV2ZW50KTogYm9vbGVhbiB7XG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KTtcbiAgICByZXR1cm4gKFxuICAgICAgICBzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13YWlscy1kcmFnZ2FibGVcIikudHJpbSgpID09PSBcImRyYWdcIlxuICAgICAgICAmJiBldmVudC5vZmZzZXRYID49IDBcbiAgICAgICAgJiYgZXZlbnQub2Zmc2V0WCA8IHRhcmdldC5jbGllbnRXaWR0aFxuICAgICAgICAmJiBldmVudC5vZmZzZXRZID49IDBcbiAgICAgICAgJiYgZXZlbnQub2Zmc2V0WSA8IHRhcmdldC5jbGllbnRIZWlnaHRcbiAgICApO1xufVxuXG4vLyBVc2UgY29uc3RhbnRzIHRvIGF2b2lkIGNvbXBhcmluZyBzdHJpbmdzIG11bHRpcGxlIHRpbWVzLlxuY29uc3QgTW91c2VEb3duID0gMDtcbmNvbnN0IE1vdXNlVXAgICA9IDE7XG5jb25zdCBNb3VzZU1vdmUgPSAyO1xuXG5mdW5jdGlvbiB1cGRhdGUoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICAvLyBXaW5kb3dzIHN1cHByZXNzZXMgbW91c2UgZXZlbnRzIGF0IHRoZSBlbmQgb2YgZHJhZ2dpbmcgb3IgcmVzaXppbmcsXG4gICAgLy8gc28gd2UgbmVlZCB0byBiZSBzbWFydCBhbmQgc3ludGhlc2l6ZSBidXR0b24gZXZlbnRzLlxuXG4gICAgbGV0IGV2ZW50VHlwZTogbnVtYmVyLCBldmVudEJ1dHRvbnMgPSBldmVudC5idXR0b25zO1xuICAgIHN3aXRjaCAoZXZlbnQudHlwZSkge1xuICAgICAgICBjYXNlICdtb3VzZWRvd24nOlxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VEb3duO1xuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zIHwgKDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIGNhc2UgJ21vdXNldXAnOlxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VVcDtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyAmIH4oMSA8PCBldmVudC5idXR0b24pOyB9XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlTW92ZTtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9uczsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgfVxuXG4gICAgbGV0IHJlbGVhc2VkID0gYnV0dG9ucyAmIH5ldmVudEJ1dHRvbnM7XG4gICAgbGV0IHByZXNzZWQgPSBldmVudEJ1dHRvbnMgJiB+YnV0dG9ucztcblxuICAgIGJ1dHRvbnMgPSBldmVudEJ1dHRvbnM7XG5cbiAgICAvLyBTeW50aGVzaXplIGEgcmVsZWFzZS1wcmVzcyBzZXF1ZW5jZSBpZiB3ZSBkZXRlY3QgYSBwcmVzcyBvZiBhbiBhbHJlYWR5IHByZXNzZWQgYnV0dG9uLlxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlRG93biAmJiAhKHByZXNzZWQgJiBldmVudC5idXR0b24pKSB7XG4gICAgICAgIHJlbGVhc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XG4gICAgICAgIHByZXNzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcbiAgICB9XG5cbiAgICAvLyBTdXBwcmVzcyBhbGwgYnV0dG9uIGV2ZW50cyBkdXJpbmcgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLFxuICAgIC8vIHVubGVzcyB0aGlzIGlzIGEgbW91c2V1cCBldmVudCB0aGF0IGlzIGVuZGluZyBhIGRyYWcgYWN0aW9uLlxuICAgIGlmIChcbiAgICAgICAgZXZlbnRUeXBlICE9PSBNb3VzZU1vdmUgLy8gRmFzdCBwYXRoIGZvciBtb3VzZW1vdmVcbiAgICAgICAgJiYgcmVzaXppbmdcbiAgICAgICAgfHwgKFxuICAgICAgICAgICAgZHJhZ2dpbmdcbiAgICAgICAgICAgICYmIChcbiAgICAgICAgICAgICAgICBldmVudFR5cGUgPT09IE1vdXNlRG93blxuICAgICAgICAgICAgICAgIHx8IGV2ZW50LmJ1dHRvbiAhPT0gMFxuICAgICAgICAgICAgKVxuICAgICAgICApXG4gICAgKSB7XG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5zdG9wUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG5cbiAgICAvLyBIYW5kbGUgcmVsZWFzZXNcbiAgICBpZiAocmVsZWFzZWQgJiAxKSB7IHByaW1hcnlVcChldmVudCk7IH1cbiAgICAvLyBIYW5kbGUgcHJlc3Nlc1xuICAgIGlmIChwcmVzc2VkICYgMSkgeyBwcmltYXJ5RG93bihldmVudCk7IH1cblxuICAgIC8vIEhhbmRsZSBtb3VzZW1vdmVcbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZU1vdmUpIHsgb25Nb3VzZU1vdmUoZXZlbnQpOyB9O1xufVxuXG5mdW5jdGlvbiBwcmltYXJ5RG93bihldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xuICAgIC8vIFJlc2V0IHJlYWRpbmVzcyBzdGF0ZS5cbiAgICBjYW5EcmFnID0gZmFsc2U7XG4gICAgY2FuUmVzaXplID0gZmFsc2U7XG5cbiAgICAvLyBJZ25vcmUgcmVwZWF0ZWQgY2xpY2tzIG9uIG1hY09TIGFuZCBMaW51eC5cbiAgICBpZiAoIUlzV2luZG93cygpKSB7XG4gICAgICAgIGlmIChldmVudC50eXBlID09PSAnbW91c2Vkb3duJyAmJiBldmVudC5idXR0b24gPT09IDAgJiYgZXZlbnQuZGV0YWlsICE9PSAxKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICBpZiAocmVzaXplRWRnZSkge1xuICAgICAgICAvLyBEbyBub3QgYXJtIGVkZ2UgcmVzaXplIGZyb20gc3ludGhlc2l6ZWQgcHJlc3NlcyBvYnNlcnZlZCBvbiBtb3ZlL3VwOlxuICAgICAgICAvLyBlbnRlcmluZyB0aGUgd2luZG93IHdpdGggdGhlIHByaW1hcnkgYnV0dG9uIGFscmVhZHkgaGVsZCBzaG91bGQgbm90XG4gICAgICAgIC8vIHN0ZWFsIGFub3RoZXIgZ2VzdHVyZSBpbnRvIGEgcmVzaXplLlxuICAgICAgICBpZiAoZXZlbnQudHlwZSAhPT0gJ21vdXNlZG93bicpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFJlYWR5IHRvIHJlc2l6ZSBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lLlxuICAgICAgICBjYW5SZXNpemUgPSB0cnVlO1xuICAgICAgICAvLyBEbyBub3Qgc3RhcnQgZHJhZyBvcGVyYXRpb25zIHdoZW4gb24gcmVzaXplIGVkZ2VzLlxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUmV0cmlldmUgdGFyZ2V0IGVsZW1lbnRcbiAgICAvLyBSZWFkeSB0byBkcmFnIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUgb24gYSBkcmFnZ2FibGUgZWxlbWVudC5cbiAgICAvLyBJZ25vcmUgY2xpY2tzIG9uIHRoZSBzY3JvbGxiYXIuXG4gICAgY2FuRHJhZyA9IGlzRHJhZ2dhYmxlRXZlbnQoZXZlbnQpO1xufVxuXG5mdW5jdGlvbiBwcmltYXJ5VXAoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICAvLyBTdG9wIGRyYWdnaW5nIGFuZCByZXNpemluZy5cbiAgICBjYW5EcmFnID0gZmFsc2U7XG4gICAgZHJhZ2dpbmcgPSBmYWxzZTtcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcbiAgICByZXNpemluZyA9IGZhbHNlO1xufVxuXG5jb25zdCBjdXJzb3JGb3JFZGdlID0gT2JqZWN0LmZyZWV6ZSh7XG4gICAgXCJzZS1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxuICAgIFwic3ctcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcbiAgICBcIm53LXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXG4gICAgXCJuZS1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxuICAgIFwidy1yZXNpemVcIjogXCJldy1yZXNpemVcIixcbiAgICBcIm4tcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXG4gICAgXCJzLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxuICAgIFwiZS1yZXNpemVcIjogXCJldy1yZXNpemVcIixcbn0pXG5cbmZ1bmN0aW9uIHNldFJlc2l6ZShlZGdlPzoga2V5b2YgdHlwZW9mIGN1cnNvckZvckVkZ2UpOiB2b2lkIHtcbiAgICBpZiAoZWRnZSkge1xuICAgICAgICBpZiAoIXJlc2l6ZUVkZ2UpIHsgZGVmYXVsdEN1cnNvciA9IGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yOyB9XG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gY3Vyc29yRm9yRWRnZVtlZGdlXTtcbiAgICB9IGVsc2UgaWYgKCFlZGdlICYmIHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBkZWZhdWx0Q3Vyc29yO1xuICAgIH1cblxuICAgIHJlc2l6ZUVkZ2UgPSBlZGdlIHx8IFwiXCI7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGV2ZW50OiBNb3VzZUV2ZW50KTogdm9pZCB7XG4gICAgaWYgKGNhblJlc2l6ZSAmJiByZXNpemVFZGdlKSB7XG4gICAgICAgIC8vIFN0YXJ0IHJlc2l6aW5nLlxuICAgICAgICByZXNpemluZyA9IHRydWU7XG4gICAgICAgIGludm9rZShcIndhaWxzOnJlc2l6ZTpcIiArIHJlc2l6ZUVkZ2UpO1xuICAgIH0gZWxzZSBpZiAoY2FuRHJhZykge1xuICAgICAgICAvLyBTdGFydCBkcmFnZ2luZy5cbiAgICAgICAgZHJhZ2dpbmcgPSB0cnVlO1xuICAgICAgICBpbnZva2UoXCJ3YWlsczpkcmFnXCIpO1xuICAgIH1cblxuICAgIGlmIChkcmFnZ2luZyB8fCByZXNpemluZykge1xuICAgICAgICAvLyBFaXRoZXIgZHJhZyBvciByZXNpemUgaXMgb25nb2luZyxcbiAgICAgICAgLy8gcmVzZXQgcmVhZGluZXNzIGFuZCBzdG9wIHByb2Nlc3NpbmcuXG4gICAgICAgIGNhbkRyYWcgPSBjYW5SZXNpemUgPSBmYWxzZTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGlmICghcmVzaXphYmxlIHx8ICghSXNXaW5kb3dzKCkgJiYgIShJc0xpbnV4KCkgJiYgR2V0RmxhZyhcImZyYW1lbGVzc1wiKSkpKSB7XG4gICAgICAgIGlmIChyZXNpemVFZGdlKSB7IHNldFJlc2l6ZSgpOyB9XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBjb25zdCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xuXG4gICAgLy8gRXh0cmEgcGl4ZWxzIGZvciB0aGUgY29ybmVyIGFyZWFzLlxuICAgIGNvbnN0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xuXG4gICAgLy8gV2hlbiBhIHNjcm9sbGJhciBpcyBwcmVzZW50IGF0IHRoZSB3aW5kb3cgZWRnZSBpdCBjb25zdW1lcyBtb3VzZSBldmVudHMgaW4gdGhhdCBzdHJpcC5cbiAgICAvLyBTaGlmdCB0aGUgZWZmZWN0aXZlIGNvbnRlbnQgZWRnZSBpbndhcmQgc28gdGhlIHJlc2l6ZSB6b25lIHNpdHMganVzdCBiZWZvcmUgdGhlIHNjcm9sbGJhci5cbiAgICBjb25zdCBzY3JvbGxiYXJXaWR0aCA9IE1hdGgubWF4KDAsIHdpbmRvdy5pbm5lcldpZHRoIC0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LmNsaWVudFdpZHRoKTtcbiAgICBjb25zdCBzY3JvbGxiYXJIZWlnaHQgPSBNYXRoLm1heCgwLCB3aW5kb3cuaW5uZXJIZWlnaHQgLSBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuY2xpZW50SGVpZ2h0KTtcbiAgICBjb25zdCByaWdodENvbnRlbnRFZGdlID0gd2luZG93LmlubmVyV2lkdGggLSBzY3JvbGxiYXJXaWR0aDtcbiAgICBjb25zdCBib3R0b21Db250ZW50RWRnZSA9IHdpbmRvdy5pbm5lckhlaWdodCAtIHNjcm9sbGJhckhlaWdodDtcblxuICAgIGNvbnN0IHJpZ2h0Qm9yZGVyID0gZXZlbnQuY2xpZW50WCA8IHJpZ2h0Q29udGVudEVkZ2UgJiYgKHJpZ2h0Q29udGVudEVkZ2UgLSBldmVudC5jbGllbnRYKSA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGNvbnN0IGxlZnRCb3JkZXIgPSBldmVudC5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgY29uc3QgdG9wQm9yZGVyID0gZXZlbnQuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcbiAgICBjb25zdCBib3R0b21Cb3JkZXIgPSBldmVudC5jbGllbnRZIDwgYm90dG9tQ29udGVudEVkZ2UgJiYgKGJvdHRvbUNvbnRlbnRFZGdlIC0gZXZlbnQuY2xpZW50WSkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG5cbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lciBhcmVhcy5cbiAgICBjb25zdCByaWdodENvcm5lciA9IGV2ZW50LmNsaWVudFggPCByaWdodENvbnRlbnRFZGdlICYmIChyaWdodENvbnRlbnRFZGdlIC0gZXZlbnQuY2xpZW50WCkgPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgbGVmdENvcm5lciA9IGV2ZW50LmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgdG9wQ29ybmVyID0gZXZlbnQuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG4gICAgY29uc3QgYm90dG9tQ29ybmVyID0gZXZlbnQuY2xpZW50WSA8IGJvdHRvbUNvbnRlbnRFZGdlICYmIChib3R0b21Db250ZW50RWRnZSAtIGV2ZW50LmNsaWVudFkpIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcblxuICAgIGlmICghbGVmdENvcm5lciAmJiAhdG9wQ29ybmVyICYmICFib3R0b21Db3JuZXIgJiYgIXJpZ2h0Q29ybmVyKSB7XG4gICAgICAgIC8vIE9wdGltaXNhdGlvbjogb3V0IG9mIGFsbCBjb3JuZXIgYXJlYXMgaW1wbGllcyBvdXQgb2YgYm9yZGVycy5cbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxuICAgIC8vIERldGVjdCBjb3JuZXJzLlxuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcbiAgICAvLyBEZXRlY3QgYm9yZGVycy5cbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xuICAgIC8vIE91dCBvZiBib3JkZXIgYXJlYS5cbiAgICBlbHNlIHNldFJlc2l6ZSgpO1xufVxuIiwgIi8qXG4gXyAgICAgX18gICAgXyBfX1xufCB8ICAgLyAvX19fKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgaW52b2tlIH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyB3aGVuUmVhZHkgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxudHlwZSBOb25DbGllbnRSZWdpb25LaW5kID0gXCJjYXB0aW9uXCIgfCBcIm1pbmltaXplXCIgfCBcIm1heGltaXplXCIgfCBcImNsb3NlXCI7XG5cbmludGVyZmFjZSBOb25DbGllbnRSZWdpb24ge1xuICAgIGtpbmQ6IE5vbkNsaWVudFJlZ2lvbktpbmQ7XG4gICAgbGVmdDogbnVtYmVyO1xuICAgIHRvcDogbnVtYmVyO1xuICAgIHJpZ2h0OiBudW1iZXI7XG4gICAgYm90dG9tOiBudW1iZXI7XG59XG5cbi8qXG4tLXdhaWxzLW5vbi1jbGllbnQtcmVnaW9uOiBjYXB0aW9uOyAgbWFya3MgYW4gYXJlYSB0aGF0IGNhbiBkcmFnIHRoZSB3aW5kb3dcbi0td2FpbHMtbm9uLWNsaWVudC1yZWdpb246IG1pbmltaXplOyBtYXJrcyBhIGN1c3RvbSBtaW5pbWl6ZSBidXR0b25cbi0td2FpbHMtbm9uLWNsaWVudC1yZWdpb246IG1heGltaXplOyBtYXJrcyBhIGN1c3RvbSBtYXhpbWl6ZSBidXR0b25cbi0td2FpbHMtbm9uLWNsaWVudC1yZWdpb246IGNsb3NlOyAgICBtYXJrcyBhIGN1c3RvbSBjbG9zZSBidXR0b25cbiovXG5jb25zdCByZWdpb25Qcm9wZXJ0eSA9IFwiLS13YWlscy1ub24tY2xpZW50LXJlZ2lvblwiO1xuY29uc3QgcnVudGltZUNvbmZpZ1JlYWR5RXZlbnQgPSBcIndhaWxzOnJ1bnRpbWUtY29uZmlnLXJlYWR5XCI7XG5jb25zdCB2YWxpZFJlZ2lvbnMgPSBuZXcgU2V0PE5vbkNsaWVudFJlZ2lvbktpbmQ+KFtcImNhcHRpb25cIiwgXCJtaW5pbWl6ZVwiLCBcIm1heGltaXplXCIsIFwiY2xvc2VcIl0pO1xuXG4vLyBTZXR1cFxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xufVxuXG5sZXQgdXBkYXRlUGVuZGluZyA9IGZhbHNlO1xubGV0IGxhc3RQYXlsb2FkID0gXCJcIjtcbmxldCBvYnNlcnZlZEVsZW1lbnRzID0gbmV3IFNldDxFbGVtZW50PigpO1xubGV0IHJlc2l6ZU9ic2VydmVyOiBSZXNpemVPYnNlcnZlciB8IHVuZGVmaW5lZDtcbmxldCB0cmFja2luZ1N0YXJ0ZWQgPSBmYWxzZTtcblxuZnVuY3Rpb24gbm9ybWFsaXNlUmVnaW9uS2luZCh2YWx1ZTogc3RyaW5nKTogTm9uQ2xpZW50UmVnaW9uS2luZCB8IHVuZGVmaW5lZCB7XG4gICAgY29uc3QgcmVnaW9uID0gdmFsdWUudHJpbSgpLnRvTG93ZXJDYXNlKCk7XG4gICAgaWYgKHZhbGlkUmVnaW9ucy5oYXMocmVnaW9uIGFzIE5vbkNsaWVudFJlZ2lvbktpbmQpKSB7XG4gICAgICAgIHJldHVybiByZWdpb24gYXMgTm9uQ2xpZW50UmVnaW9uS2luZDtcbiAgICB9XG4gICAgcmV0dXJuIHVuZGVmaW5lZDtcbn1cblxuZnVuY3Rpb24gbm9uQ2xpZW50UmVnaW9uRm9yRWxlbWVudChlbGVtZW50OiBFbGVtZW50KTogTm9uQ2xpZW50UmVnaW9uS2luZCB8IHVuZGVmaW5lZCB7XG4gICAgaWYgKCEoZWxlbWVudCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSkge1xuICAgICAgICByZXR1cm4gdW5kZWZpbmVkO1xuICAgIH1cblxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XG4gICAgY29uc3QgcmVnaW9uID0gbm9ybWFsaXNlUmVnaW9uS2luZChzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKHJlZ2lvblByb3BlcnR5KSk7XG4gICAgaWYgKCFyZWdpb24pIHtcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcbiAgICB9XG5cbiAgICBjb25zdCBwYXJlbnQgPSBlbGVtZW50LnBhcmVudEVsZW1lbnQ7XG4gICAgaWYgKHBhcmVudCkge1xuICAgICAgICBjb25zdCBwYXJlbnRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHBhcmVudCk7XG4gICAgICAgIC8vIFRoZSBDU1MgcHJvcGVydHkgaXMgaW5oZXJpdGVkLiBPbmx5IHJlcG9ydCB0aGUgb3V0ZXJtb3N0IGVsZW1lbnQgZm9yXG4gICAgICAgIC8vIGVhY2ggY29udGlndW91cyByZWdpb24gc28gbmF0aXZlIGhpdCB0ZXN0aW5nIHNlZXMgc3RhYmxlIHJlY3RhbmdsZXMuXG4gICAgICAgIGlmIChub3JtYWxpc2VSZWdpb25LaW5kKHBhcmVudFN0eWxlLmdldFByb3BlcnR5VmFsdWUocmVnaW9uUHJvcGVydHkpKSA9PT0gcmVnaW9uKSB7XG4gICAgICAgICAgICByZXR1cm4gdW5kZWZpbmVkO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIHJlZ2lvbjtcbn1cblxuZnVuY3Rpb24gaXNWaXNpYmxlKGVsZW1lbnQ6IEhUTUxFbGVtZW50KTogYm9vbGVhbiB7XG4gICAgY29uc3Qgc3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcbiAgICByZXR1cm4gc3R5bGUuZGlzcGxheSAhPT0gXCJub25lXCIgJiZcbiAgICAgICAgc3R5bGUudmlzaWJpbGl0eSAhPT0gXCJoaWRkZW5cIiAmJlxuICAgICAgICBzdHlsZS5jb250ZW50VmlzaWJpbGl0eSAhPT0gXCJoaWRkZW5cIjtcbn1cblxuZnVuY3Rpb24gZWxlbWVudFJlZ2lvbihlbGVtZW50OiBFbGVtZW50KTogTm9uQ2xpZW50UmVnaW9uIHwgdW5kZWZpbmVkIHtcbiAgICBpZiAoIShlbGVtZW50IGluc3RhbmNlb2YgSFRNTEVsZW1lbnQpKSB7XG4gICAgICAgIHJldHVybiB1bmRlZmluZWQ7XG4gICAgfVxuXG4gICAgY29uc3Qga2luZCA9IG5vbkNsaWVudFJlZ2lvbkZvckVsZW1lbnQoZWxlbWVudCk7XG4gICAgaWYgKCFraW5kIHx8ICFpc1Zpc2libGUoZWxlbWVudCkpIHtcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcbiAgICB9XG5cbiAgICBjb25zdCByZWN0ID0gZWxlbWVudC5nZXRCb3VuZGluZ0NsaWVudFJlY3QoKTtcbiAgICBpZiAocmVjdC53aWR0aCA8PSAwIHx8IHJlY3QuaGVpZ2h0IDw9IDApIHtcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcbiAgICB9XG5cbiAgICAvLyBOYXRpdmUgaGl0IHRlc3RpbmcgcnVucyBpbiBwaHlzaWNhbCBwaXhlbHMsIHdoaWxlIERPTSBnZW9tZXRyeSBpcyBpbiBDU1MgcGl4ZWxzLlxuICAgIGNvbnN0IHNjYWxlID0gd2luZG93LmRldmljZVBpeGVsUmF0aW8gfHwgMTtcbiAgICBjb25zdCBsZWZ0ID0gTWF0aC5mbG9vcihyZWN0LmxlZnQgKiBzY2FsZSk7XG4gICAgY29uc3QgdG9wID0gTWF0aC5mbG9vcihyZWN0LnRvcCAqIHNjYWxlKTtcbiAgICBjb25zdCByaWdodCA9IE1hdGguY2VpbChyZWN0LnJpZ2h0ICogc2NhbGUpO1xuICAgIGNvbnN0IGJvdHRvbSA9IE1hdGguY2VpbChyZWN0LmJvdHRvbSAqIHNjYWxlKTtcblxuICAgIGlmIChyaWdodCA8PSBsZWZ0IHx8IGJvdHRvbSA8PSB0b3ApIHtcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcbiAgICB9XG5cbiAgICByZXR1cm4geyBraW5kLCBsZWZ0LCB0b3AsIHJpZ2h0LCBib3R0b20gfTtcbn1cblxuZnVuY3Rpb24gcmVnaW9uRWxlbWVudHMoKTogRWxlbWVudFtdIHtcbiAgICBjb25zdCBlbGVtZW50czogRWxlbWVudFtdID0gW107XG5cbiAgICBpZiAoZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50KSB7XG4gICAgICAgIGVsZW1lbnRzLnB1c2goZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50KTtcbiAgICB9XG4gICAgaWYgKGRvY3VtZW50LmJvZHkpIHtcbiAgICAgICAgZWxlbWVudHMucHVzaChkb2N1bWVudC5ib2R5KTtcbiAgICAgICAgLy8gQXBwZW5kIHZpYSBhIGxvb3A6IHNwcmVhZGluZyBhIGh1Z2UgTm9kZUxpc3QgaW50byBwdXNoKCkgb3ZlcmZsb3dzXG4gICAgICAgIC8vIHRoZSBlbmdpbmUncyBhcmd1bWVudCBsaW1pdCBvbiB2ZXJ5IGxhcmdlIGRvY3VtZW50cy5cbiAgICAgICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbChcIipcIikpIHtcbiAgICAgICAgICAgIGVsZW1lbnRzLnB1c2goZWxlbWVudCk7XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICByZXR1cm4gZWxlbWVudHM7XG59XG5cbmZ1bmN0aW9uIG9ic2VydmVSZWdpb25FbGVtZW50cyhlbGVtZW50czogRWxlbWVudFtdKTogdm9pZCB7XG4gICAgaWYgKHR5cGVvZiBSZXNpemVPYnNlcnZlciA9PT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gVHJhY2sgc2l6ZSBjaGFuZ2VzIG9ubHkgZm9yIGFjdGl2ZSByZWdpb24gZWxlbWVudHMuIERPTSBzdHJ1Y3R1cmUgYW5kIHN0eWxlXG4gICAgLy8gY2hhbmdlcyBhcmUgY292ZXJlZCBieSBNdXRhdGlvbk9ic2VydmVyIGluIHN0YXJ0Tm9uQ2xpZW50UmVnaW9uVHJhY2tpbmcoKS5cbiAgICByZXNpemVPYnNlcnZlciA/Pz0gbmV3IFJlc2l6ZU9ic2VydmVyKHNjaGVkdWxlVXBkYXRlKTtcbiAgICBjb25zdCBuZXh0RWxlbWVudHMgPSBuZXcgU2V0KGVsZW1lbnRzKTtcblxuICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBvYnNlcnZlZEVsZW1lbnRzKSB7XG4gICAgICAgIGlmICghbmV4dEVsZW1lbnRzLmhhcyhlbGVtZW50KSkge1xuICAgICAgICAgICAgcmVzaXplT2JzZXJ2ZXIudW5vYnNlcnZlKGVsZW1lbnQpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgZm9yIChjb25zdCBlbGVtZW50IG9mIG5leHRFbGVtZW50cykge1xuICAgICAgICBpZiAoIW9ic2VydmVkRWxlbWVudHMuaGFzKGVsZW1lbnQpKSB7XG4gICAgICAgICAgICByZXNpemVPYnNlcnZlci5vYnNlcnZlKGVsZW1lbnQpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgb2JzZXJ2ZWRFbGVtZW50cyA9IG5leHRFbGVtZW50cztcbn1cblxuZnVuY3Rpb24gdXBkYXRlTm9uQ2xpZW50UmVnaW9ucygpOiB2b2lkIHtcbiAgICB1cGRhdGVQZW5kaW5nID0gZmFsc2U7XG5cbiAgICBjb25zdCBlbGVtZW50cyA9IHJlZ2lvbkVsZW1lbnRzKCk7XG4gICAgY29uc3QgcmVnaW9uczogTm9uQ2xpZW50UmVnaW9uW10gPSBbXTtcbiAgICBjb25zdCBhY3RpdmVFbGVtZW50czogRWxlbWVudFtdID0gW107XG5cbiAgICBmb3IgKGNvbnN0IGVsZW1lbnQgb2YgZWxlbWVudHMpIHtcbiAgICAgICAgY29uc3QgcmVnaW9uID0gZWxlbWVudFJlZ2lvbihlbGVtZW50KTtcbiAgICAgICAgaWYgKHJlZ2lvbikge1xuICAgICAgICAgICAgcmVnaW9ucy5wdXNoKHJlZ2lvbik7XG4gICAgICAgICAgICBhY3RpdmVFbGVtZW50cy5wdXNoKGVsZW1lbnQpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgb2JzZXJ2ZVJlZ2lvbkVsZW1lbnRzKGFjdGl2ZUVsZW1lbnRzKTtcblxuICAgIGNvbnN0IHBheWxvYWQgPSBKU09OLnN0cmluZ2lmeSh7IHZlcnNpb246IDEsIHJlZ2lvbnMgfSk7XG4gICAgaWYgKHBheWxvYWQgPT09IGxhc3RQYXlsb2FkKSB7XG4gICAgICAgIC8vIEF2b2lkIHNlbmRpbmcgZHVwbGljYXRlIG5hdGl2ZSBtZXNzYWdlcyBkdXJpbmcgcmVzaXplIG9yIHN0eWxlIGNodXJuLlxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgbGFzdFBheWxvYWQgPSBwYXlsb2FkO1xuICAgIGludm9rZShcIndhaWxzOm5vbi1jbGllbnQtcmVnaW9uOlwiICsgcGF5bG9hZCk7XG59XG5cbmZ1bmN0aW9uIHNjaGVkdWxlVXBkYXRlKCk6IHZvaWQge1xuICAgIGlmICh1cGRhdGVQZW5kaW5nKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBCYXRjaCByZWdpb24gdXBkYXRlcyB0byBhbmltYXRpb24gZnJhbWVzIHNvIGxheW91dCBpcyBtZWFzdXJlZCBvbmNlIHBlciBmcmFtZS5cbiAgICB1cGRhdGVQZW5kaW5nID0gdHJ1ZTtcbiAgICB3aW5kb3cucmVxdWVzdEFuaW1hdGlvbkZyYW1lKHVwZGF0ZU5vbkNsaWVudFJlZ2lvbnMpO1xufVxuXG5mdW5jdGlvbiBzdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKCk6IHZvaWQge1xuICAgIGlmICh0cmFja2luZ1N0YXJ0ZWQpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIHRyYWNraW5nU3RhcnRlZCA9IHRydWU7XG4gICAgLy8gU2VuZCBhbiBpbml0aWFsIGVtcHR5IG9yIHBvcHVsYXRlZCByZWdpb24gbGlzdCBvbmNlIHRoZSBET00gaXMgcmVhZHkuXG4gICAgc2NoZWR1bGVVcGRhdGUoKTtcblxuICAgIGNvbnN0IG11dGF0aW9uT2JzZXJ2ZXIgPSBuZXcgTXV0YXRpb25PYnNlcnZlcihzY2hlZHVsZVVwZGF0ZSk7XG4gICAgbXV0YXRpb25PYnNlcnZlci5vYnNlcnZlKGRvY3VtZW50LmRvY3VtZW50RWxlbWVudCwge1xuICAgICAgICBhdHRyaWJ1dGVzOiB0cnVlLFxuICAgICAgICBjaGlsZExpc3Q6IHRydWUsXG4gICAgICAgIHN1YnRyZWU6IHRydWUsXG4gICAgfSk7XG5cbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihcInJlc2l6ZVwiLCBzY2hlZHVsZVVwZGF0ZSk7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoXCJzY3JvbGxcIiwgc2NoZWR1bGVVcGRhdGUsIHRydWUpO1xuICAgIHdpbmRvdy52aXN1YWxWaWV3cG9ydD8uYWRkRXZlbnRMaXN0ZW5lcihcInJlc2l6ZVwiLCBzY2hlZHVsZVVwZGF0ZSk7XG4gICAgd2luZG93LnZpc3VhbFZpZXdwb3J0Py5hZGRFdmVudExpc3RlbmVyKFwic2Nyb2xsXCIsIHNjaGVkdWxlVXBkYXRlKTtcbn1cblxuZnVuY3Rpb24gdHJ5U3RhcnROb25DbGllbnRSZWdpb25UcmFja2luZygpOiBib29sZWFuIHtcbiAgICBjb25zdCBvcyA9IHdpbmRvdy5fd2FpbHMuZW52aXJvbm1lbnQ/Lk9TO1xuICAgIGlmIChvcyA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICBjb25zdCBlbmFibGVkID0gd2luZG93Ll93YWlscy5mbGFncz8ubm9uQ2xpZW50UmVnaW9uVHJhY2tpbmc7XG4gICAgaWYgKG9zID09PSBcIndpbmRvd3NcIikge1xuICAgICAgICBpZiAoZW5hYmxlZCA9PT0gdHJ1ZSkge1xuICAgICAgICAgICAgd2hlblJlYWR5KHN0YXJ0Tm9uQ2xpZW50UmVnaW9uVHJhY2tpbmcpO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiB0cnVlO1xuICAgIH1cblxuICAgIHJldHVybiB0cnVlO1xufVxuXG5pZiAoaGFzRE9NICYmICF0cnlTdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKCkpIHtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihydW50aW1lQ29uZmlnUmVhZHlFdmVudCwgdHJ5U3RhcnROb25DbGllbnRSZWdpb25UcmFja2luZywgeyBvbmNlOiB0cnVlIH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcblxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XG5jb25zdCBTaG93TWV0aG9kID0gMTtcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xuXG4vKipcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoU2hvd01ldGhvZCk7XG59XG5cbi8qKlxuICogQ2FsbHMgdGhlIFF1aXRNZXRob2QgdG8gdGVybWluYXRlIHRoZSBwcm9ncmFtLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChRdWl0TWV0aG9kKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgQ2FuY2VsbGFibGVQcm9taXNlLCB0eXBlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnMgfSBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcblxuLy8gU2V0dXBcbmltcG9ydCB7IGhhc0RPTSB9IGZyb20gXCIuL2Vudmlyb25tZW50LmpzXCI7XG5cbmlmIChoYXNET00pIHtcbiAgICB3aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbn1cblxudHlwZSBQcm9taXNlUmVzb2x2ZXJzID0gT21pdDxDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPGFueT4sIFwicHJvbWlzZVwiIHwgXCJvbmNhbmNlbGxlZFwiPlxuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYWxsKTtcbmNvbnN0IGNhbmNlbENhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbmNlbENhbGwpO1xuY29uc3QgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xuXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XG5jb25zdCBDYW5jZWxNZXRob2QgPSAwXG5cbi8qKlxuICogSG9sZHMgYWxsIHJlcXVpcmVkIGluZm9ybWF0aW9uIGZvciBhIGJpbmRpbmcgY2FsbC5cbiAqIE1heSBwcm92aWRlIGVpdGhlciBhIG1ldGhvZCBJRCBvciBhIG1ldGhvZCBuYW1lLCBidXQgbm90IGJvdGguXG4gKi9cbmV4cG9ydCB0eXBlIENhbGxPcHRpb25zID0ge1xuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kSUQ6IG51bWJlcjtcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2ROYW1lPzogbmV2ZXI7XG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xuICAgIGFyZ3M6IGFueVtdO1xufSB8IHtcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZElEPzogbmV2ZXI7XG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kTmFtZTogc3RyaW5nO1xuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cbiAgICBhcmdzOiBhbnlbXTtcbn07XG5cbi8vIHJ1bnRpbWUuanMgbmVlZHMgdG8gdXNlIFJ1bnRpbWVFcnJvciBpbnRlcm5hbGx5IHRvIHByb3Blcmx5IHBhcnNlIGFuZCByZXR1cm5cbi8vIGVycm9ycyBmb3IgYmluZGluZyBjYWxscywgc28gaXQgaGFkIHRvIG1vdmUgdGhlcmUuIEV4cG9ydGluZyBoZXJlIGFnYWluIHRvXG4vLyBrZWVwIGZyb20gYnJlYWtpbmcgdGhlIHB1YmxpYyBDYWxsIGludGVyZmFjZS5cbmV4cG9ydCB7IFJ1bnRpbWVFcnJvciB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuLyoqXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxuICpcbiAqIEByZXR1cm5zIEEgdW5pcXVlIElEIHRoYXQgZG9lcyBub3QgZXhpc3QgaW4gdGhlIGNhbGxSZXNwb25zZXMgc2V0LlxuICovXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCk6IHN0cmluZyB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbi8qKlxuICogQ2FsbCBhIGJvdW5kIG1ldGhvZCBhY2NvcmRpbmcgdG8gdGhlIGdpdmVuIGNhbGwgb3B0aW9ucy5cbiAqXG4gKiBJbiBjYXNlIG9mIGZhaWx1cmUsIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYW4gZXhjZXB0aW9uXG4gKiBhbW9uZyBSZWZlcmVuY2VFcnJvciAodW5rbm93biBtZXRob2QpLCBUeXBlRXJyb3IgKHdyb25nIGFyZ3VtZW50IGNvdW50IG9yIHR5cGUpLFxuICoge0BsaW5rIFJ1bnRpbWVFcnJvcn0gKG1ldGhvZCByZXR1cm5lZCBhbiBlcnJvciksIG9yIG90aGVyIChuZXR3b3JrIG9yIGludGVybmFsIGVycm9ycykuXG4gKiBUaGUgZXhjZXB0aW9uIG1pZ2h0IGhhdmUgYSBcImNhdXNlXCIgZmllbGQgd2l0aCB0aGUgdmFsdWUgcmV0dXJuZWRcbiAqIGJ5IHRoZSBhcHBsaWNhdGlvbi0gb3Igc2VydmljZS1sZXZlbCBlcnJvciBtYXJzaGFsaW5nIGZ1bmN0aW9ucy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIEEgbWV0aG9kIGNhbGwgZGVzY3JpcHRvci5cbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnM6IENhbGxPcHRpb25zKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xuICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xuXG4gICAgY29uc3QgcmVzdWx0ID0gQ2FuY2VsbGFibGVQcm9taXNlLndpdGhSZXNvbHZlcnM8YW55PigpO1xuICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7IHJlc29sdmU6IHJlc3VsdC5yZXNvbHZlLCByZWplY3Q6IHJlc3VsdC5yZWplY3QgfSk7XG5cbiAgICBjb25zdCByZXF1ZXN0ID0gY2FsbChDYWxsQmluZGluZywgT2JqZWN0LmFzc2lnbih7IFwiY2FsbC1pZFwiOiBpZCB9LCBvcHRpb25zKSk7XG4gICAgbGV0IHJ1bm5pbmcgPSB0cnVlO1xuXG4gICAgcmVxdWVzdC50aGVuKChyZXMpID0+IHtcbiAgICAgICAgcnVubmluZyA9IGZhbHNlO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHJlc3VsdC5yZXNvbHZlKHJlcyk7XG4gICAgfSwgKGVycikgPT4ge1xuICAgICAgICBydW5uaW5nID0gZmFsc2U7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgcmVzdWx0LnJlamVjdChlcnIpO1xuICAgIH0pO1xuXG4gICAgY29uc3QgY2FuY2VsID0gKCkgPT4ge1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHJldHVybiBjYW5jZWxDYWxsKENhbmNlbE1ldGhvZCwge1wiY2FsbC1pZFwiOiBpZH0pLmNhdGNoKChlcnIpID0+IHtcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoXCJFcnJvciB3aGlsZSByZXF1ZXN0aW5nIGJpbmRpbmcgY2FsbCBjYW5jZWxsYXRpb246XCIsIGVycik7XG4gICAgICAgIH0pO1xuICAgIH07XG5cbiAgICByZXN1bHQub25jYW5jZWxsZWQgPSAoKSA9PiB7XG4gICAgICAgIGlmIChydW5uaW5nKSB7XG4gICAgICAgICAgICByZXR1cm4gY2FuY2VsKCk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICByZXR1cm4gcmVxdWVzdC50aGVuKGNhbmNlbCk7XG4gICAgICAgIH1cbiAgICB9O1xuXG4gICAgcmV0dXJuIHJlc3VsdC5wcm9taXNlO1xufVxuXG4vKipcbiAqIENhbGxzIGEgYm91bmQgbWV0aG9kIGJ5IG5hbWUgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXG4gKlxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCcuXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybnMgVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCeU5hbWUobWV0aG9kTmFtZTogc3RyaW5nLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gQ2FsbCh7IG1ldGhvZE5hbWUsIGFyZ3MgfSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIG51bWVyaWMgSUQgd2l0aCB0aGUgc3BlY2lmaWVkIGFyZ3VtZW50cy5cbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXG4gKlxuICogQHBhcmFtIG1ldGhvZElEIC0gVGhlIElEIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJuIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlJRChtZXRob2RJRDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gQ2FsbCh7IG1ldGhvZElELCBhcmdzIH0pO1xufVxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2luc3BlY3QtanMvaXMtY2FsbGFibGVcblxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXG4vL1xuLy8gQ29weXJpZ2h0IChjKSAyMDE1IEpvcmRhbiBIYXJiYW5kXG4vL1xuLy8gUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weVxuLy8gb2YgdGhpcyBzb2Z0d2FyZSBhbmQgYXNzb2NpYXRlZCBkb2N1bWVudGF0aW9uIGZpbGVzICh0aGUgXCJTb2Z0d2FyZVwiKSwgdG8gZGVhbFxuLy8gaW4gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0c1xuLy8gdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbFxuLy8gY29waWVzIG9mIHRoZSBTb2Z0d2FyZSwgYW5kIHRvIHBlcm1pdCBwZXJzb25zIHRvIHdob20gdGhlIFNvZnR3YXJlIGlzXG4vLyBmdXJuaXNoZWQgdG8gZG8gc28sIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxuLy9cbi8vIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cbi8vXG4vLyBUSEUgU09GVFdBUkUgSVMgUFJPVklERUQgXCJBUyBJU1wiLCBXSVRIT1VUIFdBUlJBTlRZIE9GIEFOWSBLSU5ELCBFWFBSRVNTIE9SXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSxcbi8vIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRVxuLy8gQVVUSE9SUyBPUiBDT1BZUklHSFQgSE9MREVSUyBCRSBMSUFCTEUgRk9SIEFOWSBDTEFJTSwgREFNQUdFUyBPUiBPVEhFUlxuLy8gTElBQklMSVRZLCBXSEVUSEVSIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSxcbi8vIE9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFXG4vLyBTT0ZUV0FSRS5cblxudmFyIGZuVG9TdHIgPSBGdW5jdGlvbi5wcm90b3R5cGUudG9TdHJpbmc7XG52YXIgcmVmbGVjdEFwcGx5OiB0eXBlb2YgUmVmbGVjdC5hcHBseSB8IGZhbHNlIHwgbnVsbCA9IHR5cGVvZiBSZWZsZWN0ID09PSAnb2JqZWN0JyAmJiBSZWZsZWN0ICE9PSBudWxsICYmIFJlZmxlY3QuYXBwbHk7XG52YXIgYmFkQXJyYXlMaWtlOiBhbnk7XG52YXIgaXNDYWxsYWJsZU1hcmtlcjogYW55O1xuaWYgKHR5cGVvZiByZWZsZWN0QXBwbHkgPT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIE9iamVjdC5kZWZpbmVQcm9wZXJ0eSA9PT0gJ2Z1bmN0aW9uJykge1xuICAgIHRyeSB7XG4gICAgICAgIGJhZEFycmF5TGlrZSA9IE9iamVjdC5kZWZpbmVQcm9wZXJ0eSh7fSwgJ2xlbmd0aCcsIHtcbiAgICAgICAgICAgIGdldDogZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgICAgIHRocm93IGlzQ2FsbGFibGVNYXJrZXI7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgICAgICBpc0NhbGxhYmxlTWFya2VyID0ge307XG4gICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZSBuby10aHJvdy1saXRlcmFsXG4gICAgICAgIHJlZmxlY3RBcHBseShmdW5jdGlvbiAoKSB7IHRocm93IDQyOyB9LCBudWxsLCBiYWRBcnJheUxpa2UpO1xuICAgIH0gY2F0Y2ggKF8pIHtcbiAgICAgICAgaWYgKF8gIT09IGlzQ2FsbGFibGVNYXJrZXIpIHtcbiAgICAgICAgICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XG4gICAgICAgIH1cbiAgICB9XG59IGVsc2Uge1xuICAgIHJlZmxlY3RBcHBseSA9IG51bGw7XG59XG5cbnZhciBjb25zdHJ1Y3RvclJlZ2V4ID0gL15cXHMqY2xhc3NcXGIvO1xudmFyIGlzRVM2Q2xhc3NGbiA9IGZ1bmN0aW9uIGlzRVM2Q2xhc3NGdW5jdGlvbih2YWx1ZTogYW55KTogYm9vbGVhbiB7XG4gICAgdHJ5IHtcbiAgICAgICAgdmFyIGZuU3RyID0gZm5Ub1N0ci5jYWxsKHZhbHVlKTtcbiAgICAgICAgcmV0dXJuIGNvbnN0cnVjdG9yUmVnZXgudGVzdChmblN0cik7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7IC8vIG5vdCBhIGZ1bmN0aW9uXG4gICAgfVxufTtcblxudmFyIHRyeUZ1bmN0aW9uT2JqZWN0ID0gZnVuY3Rpb24gdHJ5RnVuY3Rpb25Ub1N0cih2YWx1ZTogYW55KTogYm9vbGVhbiB7XG4gICAgdHJ5IHtcbiAgICAgICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgICAgIGZuVG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgIHJldHVybiB0cnVlO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cbn07XG52YXIgdG9TdHIgPSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nO1xudmFyIG9iamVjdENsYXNzID0gJ1tvYmplY3QgT2JqZWN0XSc7XG52YXIgZm5DbGFzcyA9ICdbb2JqZWN0IEZ1bmN0aW9uXSc7XG52YXIgZ2VuQ2xhc3MgPSAnW29iamVjdCBHZW5lcmF0b3JGdW5jdGlvbl0nO1xudmFyIGRkYUNsYXNzID0gJ1tvYmplY3QgSFRNTEFsbENvbGxlY3Rpb25dJzsgLy8gSUUgMTFcbnZhciBkZGFDbGFzczIgPSAnW29iamVjdCBIVE1MIGRvY3VtZW50LmFsbCBjbGFzc10nO1xudmFyIGRkYUNsYXNzMyA9ICdbb2JqZWN0IEhUTUxDb2xsZWN0aW9uXSc7IC8vIElFIDktMTBcbnZhciBoYXNUb1N0cmluZ1RhZyA9IHR5cGVvZiBTeW1ib2wgPT09ICdmdW5jdGlvbicgJiYgISFTeW1ib2wudG9TdHJpbmdUYWc7IC8vIGJldHRlcjogdXNlIGBoYXMtdG9zdHJpbmd0YWdgXG5cbnZhciBpc0lFNjggPSAhKDAgaW4gWyxdKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBuby1zcGFyc2UtYXJyYXlzLCBjb21tYS1zcGFjaW5nXG5cbnZhciBpc0REQTogKHZhbHVlOiBhbnkpID0+IGJvb2xlYW4gPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKCkgeyByZXR1cm4gZmFsc2U7IH07XG5pZiAodHlwZW9mIGRvY3VtZW50ID09PSAnb2JqZWN0Jykge1xuICAgIC8vIEZpcmVmb3ggMyBjYW5vbmljYWxpemVzIEREQSB0byB1bmRlZmluZWQgd2hlbiBpdCdzIG5vdCBhY2Nlc3NlZCBkaXJlY3RseVxuICAgIHZhciBhbGwgPSBkb2N1bWVudC5hbGw7XG4gICAgaWYgKHRvU3RyLmNhbGwoYWxsKSA9PT0gdG9TdHIuY2FsbChkb2N1bWVudC5hbGwpKSB7XG4gICAgICAgIGlzRERBID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCh2YWx1ZSkge1xuICAgICAgICAgICAgLyogZ2xvYmFscyBkb2N1bWVudDogZmFsc2UgKi9cbiAgICAgICAgICAgIC8vIGluIElFIDYtOCwgdHlwZW9mIGRvY3VtZW50LmFsbCBpcyBcIm9iamVjdFwiIGFuZCBpdCdzIHRydXRoeVxuICAgICAgICAgICAgaWYgKChpc0lFNjggfHwgIXZhbHVlKSAmJiAodHlwZW9mIHZhbHVlID09PSAndW5kZWZpbmVkJyB8fCB0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnKSkge1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIHZhciBzdHIgPSB0b1N0ci5jYWxsKHZhbHVlKTtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuIChcbiAgICAgICAgICAgICAgICAgICAgICAgIHN0ciA9PT0gZGRhQ2xhc3NcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MyXG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IGRkYUNsYXNzMyAvLyBvcGVyYSAxMi4xNlxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBvYmplY3RDbGFzcyAvLyBJRSA2LThcbiAgICAgICAgICAgICAgICAgICAgKSAmJiB2YWx1ZSgnJykgPT0gbnVsbDsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBlcWVxZXFcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChlKSB7IC8qKi8gfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuZnVuY3Rpb24gaXNDYWxsYWJsZVJlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55ICB7XG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxuICAgIGlmICghdmFsdWUpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHRyeSB7XG4gICAgICAgIChyZWZsZWN0QXBwbHkgYXMgYW55KSh2YWx1ZSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIGlmIChlICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIH1cbiAgICByZXR1cm4gIWlzRVM2Q2xhc3NGbih2YWx1ZSkgJiYgdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xufVxuXG5mdW5jdGlvbiBpc0NhbGxhYmxlTm9SZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSB7XG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxuICAgIGlmICghdmFsdWUpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgaWYgKHR5cGVvZiB2YWx1ZSAhPT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgdmFsdWUgIT09ICdvYmplY3QnKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmIChoYXNUb1N0cmluZ1RhZykgeyByZXR1cm4gdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpOyB9XG4gICAgaWYgKGlzRVM2Q2xhc3NGbih2YWx1ZSkpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgdmFyIHN0ckNsYXNzID0gdG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgaWYgKHN0ckNsYXNzICE9PSBmbkNsYXNzICYmIHN0ckNsYXNzICE9PSBnZW5DbGFzcyAmJiAhKC9eXFxbb2JqZWN0IEhUTUwvKS50ZXN0KHN0ckNsYXNzKSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICByZXR1cm4gdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpO1xufTtcblxuZXhwb3J0IGRlZmF1bHQgcmVmbGVjdEFwcGx5ID8gaXNDYWxsYWJsZVJlZkFwcGx5IDogaXNDYWxsYWJsZU5vUmVmQXBwbHk7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCBpc0NhbGxhYmxlIGZyb20gXCIuL2NhbGxhYmxlLmpzXCI7XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB1c2VkIGFzIHJlamVjdGlvbiByZWFzb25cbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSBpcyBjYW5jZWxsZWQgc3VjY2Vzc2Z1bGx5LlxuICpcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsRXJyb3JcImAuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIGNhdXNlIHBhc3NlZCB0byB0aGUgY2FuY2VsIG1ldGhvZCwgaWYgYW55LlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsRXJyb3JgIGluc3RhbmNlLlxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXG4gICAgICogQHBhcmFtIG9wdGlvbnMgLSBPcHRpb25zIHRvIGJlIGZvcndhcmRlZCB0byB0aGUgRXJyb3IgY29uc3RydWN0b3IuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xuICAgICAgICBzdXBlcihtZXNzYWdlLCBvcHRpb25zKTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxFcnJvclwiO1xuICAgIH1cbn1cblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb25cbiAqIGluIGNhc2UgYSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlfSByZWplY3RzIGFmdGVyIGJlaW5nIGNhbmNlbGxlZCxcbiAqIG9yIHdoZW4gdGhlIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgdGhyb3dzIG9yIHJlamVjdHMuXG4gKlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiYC5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgcmVhc29uIHRoZSBwcm9taXNlIHJlamVjdGVkIHdpdGguXG4gKlxuICogQmVjYXVzZSB0aGUgb3JpZ2luYWwgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLFxuICogYSB3cmFwcGVyIHByb21pc2Ugd2lsbCBiZSBwYXNzZWQgdG8gdGhlIHVuaGFuZGxlZCByZWplY3Rpb24gbGlzdGVuZXIgaW5zdGVhZC5cbiAqIFRoZSB7QGxpbmsgcHJvbWlzZX0gcHJvcGVydHkgaG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIG9yaWdpbmFsIHByb21pc2UuXG4gKi9cbmV4cG9ydCBjbGFzcyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBIb2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgcHJvbWlzZSB0aGF0IHdhcyBjYW5jZWxsZWQgYW5kIHRoZW4gcmVqZWN0ZWQuXG4gICAgICovXG4gICAgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+O1xuXG4gICAgLyoqXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JgIGluc3RhbmNlLlxuICAgICAqIEBwYXJhbSBwcm9taXNlIC0gVGhlIHByb21pc2UgdGhhdCBjYXVzZWQgdGhlIGVycm9yIG9yaWdpbmFsbHkuXG4gICAgICogQHBhcmFtIHJlYXNvbiAtIFRoZSByZWplY3Rpb24gcmVhc29uLlxuICAgICAqIEBwYXJhbSBpbmZvIC0gQW4gb3B0aW9uYWwgaW5mb3JtYXRpdmUgbWVzc2FnZSBzcGVjaWZ5aW5nIHRoZSBjaXJjdW1zdGFuY2VzIGluIHdoaWNoIHRoZSBlcnJvciB3YXMgdGhyb3duLlxuICAgICAqICAgICAgICAgICAgICAgRGVmYXVsdHMgdG8gdGhlIHN0cmluZyBgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiYC5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3Rvcihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHJlYXNvbj86IGFueSwgaW5mbz86IHN0cmluZykge1xuICAgICAgICBzdXBlcigoaW5mbyA/PyBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCIpICsgXCIgUmVhc29uOiBcIiArIGVycm9yTWVzc2FnZShyZWFzb24pLCB7IGNhdXNlOiByZWFzb24gfSk7XG4gICAgICAgIHRoaXMucHJvbWlzZSA9IHByb21pc2U7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcIjtcbiAgICB9XG59XG5cbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4gPSAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yID0gKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyID0gKGNhdXNlPzogYW55KSA9PiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+ID0gKHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+LCByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yKSA9PiB2b2lkO1xuXG5leHBvcnQgaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCk6IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj47XG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xufVxuXG4vKipcbiAqIFdyYXBzIGEgY2FuY2VsbGFibGUgcHJvbWlzZSBhbG9uZyB3aXRoIGl0cyByZXNvbHV0aW9uIG1ldGhvZHMuXG4gKiBUaGUgYG9uY2FuY2VsbGVkYCBmaWVsZCB3aWxsIGJlIG51bGwgaW5pdGlhbGx5IGJ1dCBtYXkgYmUgc2V0IHRvIHByb3ZpZGUgYSBjdXN0b20gY2FuY2VsbGF0aW9uIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcbiAgICByZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPjtcbiAgICByZWplY3Q6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yO1xuICAgIG9uY2FuY2VsbGVkOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgfCBudWxsO1xufVxuXG5pbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUge1xuICAgIHJlYWRvbmx5IHJvb3Q6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlO1xuICAgIHJlc29sdmluZzogYm9vbGVhbjtcbiAgICBzZXR0bGVkOiBib29sZWFuO1xuICAgIHJlYXNvbj86IENhbmNlbEVycm9yO1xufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgYmFycmllclN5bSA9IFN5bWJvbChcImJhcnJpZXJcIik7XG5jb25zdCBjYW5jZWxJbXBsU3ltID0gU3ltYm9sKFwiY2FuY2VsSW1wbFwiKTtcbmNvbnN0IHNwZWNpZXM6IHR5cGVvZiBTeW1ib2wuc3BlY2llcyA9IFN5bWJvbC5zcGVjaWVzID8/IFN5bWJvbChcInNwZWNpZXNQb2x5ZmlsbFwiKTtcblxuLyoqXG4gKiBBIHByb21pc2Ugd2l0aCBhbiBhdHRhY2hlZCBtZXRob2QgZm9yIGNhbmNlbGxpbmcgbG9uZy1ydW5uaW5nIG9wZXJhdGlvbnMgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbH0pLlxuICogQ2FuY2VsbGF0aW9uIGNhbiBvcHRpb25hbGx5IGJlIGJvdW5kIHRvIGFuIHtAbGluayBBYm9ydFNpZ25hbH1cbiAqIGZvciBiZXR0ZXIgY29tcG9zYWJpbGl0eSAoc2VlIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2UjY2FuY2VsT259KS5cbiAqXG4gKiBDYW5jZWxsaW5nIGEgcGVuZGluZyBwcm9taXNlIHdpbGwgcmVzdWx0IGluIGFuIGltbWVkaWF0ZSByZWplY3Rpb25cbiAqIHdpdGggYW4gaW5zdGFuY2Ugb2Yge0BsaW5rIENhbmNlbEVycm9yfSBhcyByZWFzb24sXG4gKiBidXQgd2hvZXZlciBzdGFydGVkIHRoZSBwcm9taXNlIHdpbGwgYmUgcmVzcG9uc2libGVcbiAqIGZvciBhY3R1YWxseSBhYm9ydGluZyB0aGUgdW5kZXJseWluZyBvcGVyYXRpb24uXG4gKiBUbyB0aGlzIHB1cnBvc2UsIHRoZSBjb25zdHJ1Y3RvciBhbmQgYWxsIGNoYWluaW5nIG1ldGhvZHNcbiAqIGFjY2VwdCBvcHRpb25hbCBjYW5jZWxsYXRpb24gY2FsbGJhY2tzLlxuICpcbiAqIElmIGEgYENhbmNlbGxhYmxlUHJvbWlzZWAgc3RpbGwgcmVzb2x2ZXMgYWZ0ZXIgaGF2aW5nIGJlZW4gY2FuY2VsbGVkLFxuICogdGhlIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZC4gSWYgaXQgcmVqZWN0cywgdGhlIHJlYXNvblxuICogd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLFxuICogd3JhcHBlZCBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gaW5zdGFuY2UuXG4gKiBUbyBmYWNpbGl0YXRlIHRoZSBoYW5kbGluZyBvZiBjYW5jZWxsYXRpb24gcmVxdWVzdHMsXG4gKiBjYW5jZWxsZWQgYENhbmNlbGxhYmxlUHJvbWlzZWBzIHdpbGwgX25vdF8gcmVwb3J0IHVuaGFuZGxlZCBgQ2FuY2VsRXJyb3Jgc1xuICogd2hvc2UgYGNhdXNlYCBmaWVsZCBpcyB0aGUgc2FtZSBhcyB0aGUgb25lIHdpdGggd2hpY2ggdGhlIGN1cnJlbnQgcHJvbWlzZSB3YXMgY2FuY2VsbGVkLlxuICpcbiAqIEFsbCB1c3VhbCBwcm9taXNlIG1ldGhvZHMgYXJlIGRlZmluZWQgYW5kIHJldHVybiBhIGBDYW5jZWxsYWJsZVByb21pc2VgXG4gKiB3aG9zZSBjYW5jZWwgbWV0aG9kIHdpbGwgY2FuY2VsIHRoZSBwYXJlbnQgb3BlcmF0aW9uIGFzIHdlbGwsIHByb3BhZ2F0aW5nIHRoZSBjYW5jZWxsYXRpb24gcmVhc29uXG4gKiB1cHdhcmRzIHRocm91Z2ggcHJvbWlzZSBjaGFpbnMuXG4gKiBDb252ZXJzZWx5LCBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aWxsIG5vdCBhdXRvbWF0aWNhbGx5IGNhbmNlbCBkZXBlbmRlbnQgcHJvbWlzZXMgZG93bnN0cmVhbTpcbiAqIGBgYHRzXG4gKiBsZXQgcm9vdCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gKiBsZXQgY2hpbGQxID0gcm9vdC50aGVuKCgpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMiA9IGNoaWxkMS50aGVuKCgpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMyA9IHJvb3QuY2F0Y2goKCkgPT4geyAuLi4gfSk7XG4gKiBjaGlsZDEuY2FuY2VsKCk7IC8vIENhbmNlbHMgY2hpbGQxIGFuZCByb290LCBidXQgbm90IGNoaWxkMiBvciBjaGlsZDNcbiAqIGBgYFxuICogQ2FuY2VsbGluZyBhIHByb21pc2UgdGhhdCBoYXMgYWxyZWFkeSBzZXR0bGVkIGlzIHNhZmUgYW5kIGhhcyBubyBjb25zZXF1ZW5jZS5cbiAqXG4gKiBUaGUgYGNhbmNlbGAgbWV0aG9kIHJldHVybnMgYSBwcm9taXNlIHRoYXQgX2Fsd2F5cyBmdWxmaWxsc19cbiAqIGFmdGVyIHRoZSB3aG9sZSBjaGFpbiBoYXMgcHJvY2Vzc2VkIHRoZSBjYW5jZWwgcmVxdWVzdFxuICogYW5kIGFsbCBhdHRhY2hlZCBjYWxsYmFja3MgdXAgdG8gdGhhdCBtb21lbnQgaGF2ZSBydW4uXG4gKlxuICogQWxsIEVTMjAyNCBwcm9taXNlIG1ldGhvZHMgKHN0YXRpYyBhbmQgaW5zdGFuY2UpIGFyZSBkZWZpbmVkIG9uIENhbmNlbGxhYmxlUHJvbWlzZSxcbiAqIGJ1dCBhY3R1YWwgYXZhaWxhYmlsaXR5IG1heSB2YXJ5IHdpdGggT1Mvd2VidmlldyB2ZXJzaW9uLlxuICpcbiAqIEluIGxpbmUgd2l0aCB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZyxcbiAqIGBDYW5jZWxsYWJsZVByb21pc2VgIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuXG4gKiBFeHRlbmRlcnMgc2hvdWxkIHRha2UgY2FyZSB0byBwcm92aWRlIHRoZWlyIG93biBtZXRob2QgaW1wbGVtZW50YXRpb25zLlxuICogVGhpcyBtaWdodCBiZSByZWNvbnNpZGVyZWQgaW4gY2FzZSB0aGUgcHJvcG9zYWwgaXMgcmV0aXJlZC5cbiAqXG4gKiBDYW5jZWxsYWJsZVByb21pc2UgaXMgYSB3cmFwcGVyIGFyb3VuZCB0aGUgRE9NIFByb21pc2Ugb2JqZWN0XG4gKiBhbmQgaXMgY29tcGxpYW50IHdpdGggdGhlIFtQcm9taXNlcy9BKyBzcGVjaWZpY2F0aW9uXShodHRwczovL3Byb21pc2VzYXBsdXMuY29tLylcbiAqIChpdCBwYXNzZXMgdGhlIFtjb21wbGlhbmNlIHN1aXRlXShodHRwczovL2dpdGh1Yi5jb20vcHJvbWlzZXMtYXBsdXMvcHJvbWlzZXMtdGVzdHMpKVxuICogaWYgc28gaXMgdGhlIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24uXG4gKi9cbmV4cG9ydCBjbGFzcyBDYW5jZWxsYWJsZVByb21pc2U8VD4gZXh0ZW5kcyBQcm9taXNlPFQ+IGltcGxlbWVudHMgUHJvbWlzZUxpa2U8VD4sIENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VD4ge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIC8qKiBAaW50ZXJuYWwgKi9cbiAgICBwcml2YXRlIFtiYXJyaWVyU3ltXSE6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+IHwgbnVsbDtcbiAgICAvKiogQGludGVybmFsICovXG4gICAgcHJpdmF0ZSByZWFkb25seSBbY2FuY2VsSW1wbFN5bV0hOiAocmVhc29uOiBDYW5jZWxFcnJvcikgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBgQ2FuY2VsbGFibGVQcm9taXNlYC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBleGVjdXRvciAtIEEgY2FsbGJhY2sgdXNlZCB0byBpbml0aWFsaXplIHRoZSBwcm9taXNlLiBUaGlzIGNhbGxiYWNrIGlzIHBhc3NlZCB0d28gYXJndW1lbnRzOlxuICAgICAqICAgICAgICAgICAgICAgICAgIGEgYHJlc29sdmVgIGNhbGxiYWNrIHVzZWQgdG8gcmVzb2x2ZSB0aGUgcHJvbWlzZSB3aXRoIGEgdmFsdWVcbiAgICAgKiAgICAgICAgICAgICAgICAgICBvciB0aGUgcmVzdWx0IG9mIGFub3RoZXIgcHJvbWlzZSAocG9zc2libHkgY2FuY2VsbGFibGUpLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGFuZCBhIGByZWplY3RgIGNhbGxiYWNrIHVzZWQgdG8gcmVqZWN0IHRoZSBwcm9taXNlIHdpdGggYSBwcm92aWRlZCByZWFzb24gb3IgZXJyb3IuXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgdGhlIHZhbHVlIHByb3ZpZGVkIHRvIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiAgICAgICAgICAgICAgICAgICAoaXQgaGFzIGEgYHRoZW5gIF9hbmRfIGEgYGNhbmNlbGAgbWV0aG9kKSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBmb3J3YXJkZWQgdG8gdGhhdCBvYmplY3QgYW5kIHRoZSBvbmNhbmNlbGxlZCB3aWxsIG5vdCBiZSBpbnZva2VkIGFueW1vcmUuXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgYW55IG9uZSBvZiB0aGUgdHdvIGNhbGxiYWNrcyBpcyBjYWxsZWQgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQsXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhlIHByb3ZpZGVkIHZhbHVlcyB3aWxsIGJlIGNhbmNlbGxlZCBhbmQgcmVzb2x2ZWQgYXMgdXN1YWwsXG4gICAgICogICAgICAgICAgICAgICAgICAgYnV0IHRoZWlyIHJlc3VsdHMgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogICAgICAgICAgICAgICAgICAgSG93ZXZlciwgaWYgdGhlIHJlc29sdXRpb24gcHJvY2VzcyB1bHRpbWF0ZWx5IGVuZHMgdXAgaW4gYSByZWplY3Rpb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGF0IGlzIG5vdCBkdWUgdG8gY2FuY2VsbGF0aW9uLCB0aGUgcmVqZWN0aW9uIHJlYXNvblxuICAgICAqICAgICAgICAgICAgICAgICAgIHdpbGwgYmUgd3JhcHBlZCBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn1cbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqIEBwYXJhbSBvbmNhbmNlbGxlZCAtIEl0IGlzIHRoZSBjYWxsZXIncyByZXNwb25zaWJpbGl0eSB0byBlbnN1cmUgdGhhdCBhbnkgb3BlcmF0aW9uXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgc3RhcnRlZCBieSB0aGUgZXhlY3V0b3IgaXMgcHJvcGVybHkgaGFsdGVkIHVwb24gY2FuY2VsbGF0aW9uLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFRoaXMgb3B0aW9uYWwgY2FsbGJhY2sgY2FuIGJlIHVzZWQgdG8gdGhhdCBwdXJwb3NlLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIEl0IHdpbGwgYmUgY2FsbGVkIF9zeW5jaHJvbm91c2x5XyB3aXRoIGEgY2FuY2VsbGF0aW9uIGNhdXNlXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2hlbiBjYW5jZWxsYXRpb24gaXMgcmVxdWVzdGVkLCBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHJlamVjdGVkXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0sIGJ1dCBfYmVmb3JlX1xuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFueSB7QGxpbmsgdGhlbn0ve0BsaW5rIGNhdGNofS97QGxpbmsgZmluYWxseX0gY2FsbGJhY2sgcnVucy5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9XG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2lsbCBvbmx5IGZ1bGZpbGwgYWZ0ZXIgdGhlIGZvcm1lciBoYXMgc2V0dGxlZC5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBVbmhhbmRsZWQgZXhjZXB0aW9ucyBvciByZWplY3Rpb25zIGZyb20gdGhlIGNhbGxiYWNrIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyB1bmhhbmRsZWQgcmVqZWN0aW9ucy5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGNhbGxlZCBiZWZvcmUgY2FuY2VsbGF0aW9uIHdpdGggYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGlzIHByb21pc2Ugd2lsbCBiZSBkaXZlcnRlZCB0byB0aGF0IHByb21pc2UsXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW5kIHRoZSBvcmlnaW5hbCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKGV4ZWN1dG9yOiBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpIHtcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG4gICAgICAgIHN1cGVyKChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xuXG4gICAgICAgIGlmICgodGhpcy5jb25zdHJ1Y3RvciBhcyBhbnkpW3NwZWNpZXNdICE9PSBQcm9taXNlKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlIGRvZXMgbm90IHN1cHBvcnQgdHJhbnNwYXJlbnQgc3ViY2xhc3NpbmcuIFBsZWFzZSByZWZyYWluIGZyb20gb3ZlcnJpZGluZyB0aGUgW1N5bWJvbC5zcGVjaWVzXSBzdGF0aWMgcHJvcGVydHkuXCIpO1xuICAgICAgICB9XG5cbiAgICAgICAgbGV0IHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7XG4gICAgICAgICAgICBwcm9taXNlOiB0aGlzLFxuICAgICAgICAgICAgcmVzb2x2ZSxcbiAgICAgICAgICAgIHJlamVjdCxcbiAgICAgICAgICAgIGdldCBvbmNhbmNlbGxlZCgpIHsgcmV0dXJuIG9uY2FuY2VsbGVkID8/IG51bGw7IH0sXG4gICAgICAgICAgICBzZXQgb25jYW5jZWxsZWQoY2IpIHsgb25jYW5jZWxsZWQgPSBjYiA/PyB1bmRlZmluZWQ7IH1cbiAgICAgICAgfTtcblxuICAgICAgICBjb25zdCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XG4gICAgICAgICAgICBnZXQgcm9vdCgpIHsgcmV0dXJuIHN0YXRlOyB9LFxuICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcbiAgICAgICAgICAgIHNldHRsZWQ6IGZhbHNlXG4gICAgICAgIH07XG5cbiAgICAgICAgLy8gU2V0dXAgY2FuY2VsbGF0aW9uIHN5c3RlbS5cbiAgICAgICAgdm9pZCBPYmplY3QuZGVmaW5lUHJvcGVydGllcyh0aGlzLCB7XG4gICAgICAgICAgICBbYmFycmllclN5bV06IHtcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiB0cnVlLFxuICAgICAgICAgICAgICAgIHZhbHVlOiBudWxsXG4gICAgICAgICAgICB9LFxuICAgICAgICAgICAgW2NhbmNlbEltcGxTeW1dOiB7XG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgdmFsdWU6IGNhbmNlbGxlckZvcihwcm9taXNlLCBzdGF0ZSlcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG5cbiAgICAgICAgLy8gUnVuIHRoZSBhY3R1YWwgZXhlY3V0b3IuXG4gICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgc3RhdGUpO1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgZXhlY3V0b3IocmVzb2x2ZXJGb3IocHJvbWlzZSwgc3RhdGUpLCByZWplY3Rvcik7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykge1xuICAgICAgICAgICAgICAgIGNvbnNvbGUubG9nKFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBDYW5jZWxsYWJsZVByb21pc2UgZXhlY3V0b3IuXCIsIGVycik7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDYW5jZWxzIGltbWVkaWF0ZWx5IHRoZSBleGVjdXRpb24gb2YgdGhlIG9wZXJhdGlvbiBhc3NvY2lhdGVkIHdpdGggdGhpcyBwcm9taXNlLlxuICAgICAqIFRoZSBwcm9taXNlIHJlamVjdHMgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0gaW5zdGFuY2UgYXMgcmVhc29uLFxuICAgICAqIHdpdGggdGhlIHtAbGluayBDYW5jZWxFcnJvciNjYXVzZX0gcHJvcGVydHkgc2V0IHRvIHRoZSBnaXZlbiBhcmd1bWVudCwgaWYgYW55LlxuICAgICAqXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgYWZ0ZXIgdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZDtcbiAgICAgKiByZXBlYXRlZCBjYWxscyBpbiBwYXJ0aWN1bGFyIGFyZSBzYWZlLCBidXQgb25seSB0aGUgZmlyc3Qgb25lXG4gICAgICogd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cbiAgICAgKlxuICAgICAqIFRoZSBgQ2FuY2VsRXJyb3JgIGV4Y2VwdGlvbiBfbmVlZCBub3RfIGJlIGhhbmRsZWQgZXhwbGljaXRseSBfb24gdGhlIHByb21pc2VzIHRoYXQgYXJlIGJlaW5nIGNhbmNlbGxlZDpfXG4gICAgICogY2FuY2VsbGluZyBhIHByb21pc2Ugd2l0aCBubyBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlciBkb2VzIG5vdCB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQuXG4gICAgICogVGhlcmVmb3JlLCB0aGUgZm9sbG93aW5nIGlkaW9tcyBhcmUgYWxsIGVxdWFsbHkgY29ycmVjdDpcbiAgICAgKiBgYGB0c1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkuY2FuY2VsKCk7XG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2FuY2VsKCk7XG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS50aGVuKC4uLikuY2F0Y2goLi4uKS5jYW5jZWwoKTtcbiAgICAgKiBgYGBcbiAgICAgKiBXaGVuZXZlciBzb21lIGNhbmNlbGxlZCBwcm9taXNlIGluIGEgY2hhaW4gcmVqZWN0cyB3aXRoIGEgYENhbmNlbEVycm9yYFxuICAgICAqIHdpdGggdGhlIHNhbWUgY2FuY2VsbGF0aW9uIGNhdXNlIGFzIGl0c2VsZiwgdGhlIGVycm9yIHdpbGwgYmUgZGlzY2FyZGVkIHNpbGVudGx5LlxuICAgICAqIEhvd2V2ZXIsIHRoZSBgQ2FuY2VsRXJyb3JgIF93aWxsIHN0aWxsIGJlIGRlbGl2ZXJlZF8gdG8gYWxsIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyc1xuICAgICAqIGFkZGVkIGJ5IHtAbGluayB0aGVufSBhbmQgcmVsYXRlZCBtZXRob2RzOlxuICAgICAqIGBgYHRzXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcbiAgICAgKiBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLmNhdGNoKGNvbnNvbGUubG9nKTtcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gQSBDYW5jZWxFcnJvciBpcyBwcmludGVkIHRvIHRoZSBjb25zb2xlLlxuICAgICAqIGBgYFxuICAgICAqIElmIHRoZSBgQ2FuY2VsRXJyb3JgIGlzIG5vdCBoYW5kbGVkIGRvd25zdHJlYW0gYnkgdGhlIHRpbWUgaXQgcmVhY2hlc1xuICAgICAqIGEgX25vbi1jYW5jZWxsZWRfIHByb21pc2UsIGl0IF93aWxsXyB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQsXG4gICAgICoganVzdCBsaWtlIG5vcm1hbCByZWplY3Rpb25zIHdvdWxkOlxuICAgICAqIGBgYHRzXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcbiAgICAgKiBsZXQgY2hhaW5lZCA9IGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkudGhlbigoKSA9PiB7IC4uLiB9KTsgLy8gTm8gY2F0Y2guLi5cbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKTsgLy8gVW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCBvbiBjaGFpbmVkIVxuICAgICAqIGBgYFxuICAgICAqIFRoZXJlZm9yZSwgaXQgaXMgaW1wb3J0YW50IHRvIGVpdGhlciBjYW5jZWwgd2hvbGUgcHJvbWlzZSBjaGFpbnMgZnJvbSB0aGVpciB0YWlsLFxuICAgICAqIGFzIHNob3duIGluIHRoZSBjb3JyZWN0IGlkaW9tcyBhYm92ZSwgb3IgdGFrZSBjYXJlIG9mIGhhbmRsaW5nIGVycm9ycyBldmVyeXdoZXJlLlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBjYW5jZWxsYWJsZSBwcm9taXNlIHRoYXQgX2Z1bGZpbGxzXyBhZnRlciB0aGUgY2FuY2VsIGNhbGxiYWNrIChpZiBhbnkpXG4gICAgICogYW5kIGFsbCBoYW5kbGVycyBhdHRhY2hlZCB1cCB0byB0aGUgY2FsbCB0byBjYW5jZWwgaGF2ZSBydW4uXG4gICAgICogSWYgdGhlIGNhbmNlbCBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGJ5IGBjYW5jZWxgXG4gICAgICogd2lsbCBhbHNvIHdhaXQgZm9yIHRoYXQgdGhlbmFibGUgdG8gc2V0dGxlLlxuICAgICAqIFRoaXMgZW5hYmxlcyBjYWxsZXJzIHRvIHdhaXQgZm9yIHRoZSBjYW5jZWxsZWQgb3BlcmF0aW9uIHRvIHRlcm1pbmF0ZVxuICAgICAqIHdpdGhvdXQgYmVpbmcgZm9yY2VkIHRvIGhhbmRsZSBwb3RlbnRpYWwgZXJyb3JzIGF0IHRoZSBjYWxsIHNpdGUuXG4gICAgICogYGBgdHNcbiAgICAgKiBjYW5jZWxsYWJsZS5jYW5jZWwoKS50aGVuKCgpID0+IHtcbiAgICAgKiAgICAgLy8gQ2xlYW51cCBmaW5pc2hlZCwgaXQncyBzYWZlIHRvIGRvIHNvbWV0aGluZyBlbHNlLlxuICAgICAqIH0sIChlcnIpID0+IHtcbiAgICAgKiAgICAgLy8gVW5yZWFjaGFibGU6IHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gY2FuY2VsIHdpbGwgbmV2ZXIgcmVqZWN0LlxuICAgICAqIH0pO1xuICAgICAqIGBgYFxuICAgICAqIE5vdGUgdGhhdCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIF9ub3RfIGhhbmRsZSBpbXBsaWNpdGx5IGFueSByZWplY3Rpb25cbiAgICAgKiB0aGF0IG1pZ2h0IGhhdmUgb2NjdXJyZWQgYWxyZWFkeSBpbiB0aGUgY2FuY2VsbGVkIGNoYWluLlxuICAgICAqIEl0IHdpbGwganVzdCB0cmFjayB3aGV0aGVyIHJlZ2lzdGVyZWQgaGFuZGxlcnMgaGF2ZSBiZWVuIGV4ZWN1dGVkIG9yIG5vdC5cbiAgICAgKiBUaGVyZWZvcmUsIHVuaGFuZGxlZCByZWplY3Rpb25zIHdpbGwgbmV2ZXIgYmUgc2lsZW50bHkgaGFuZGxlZCBieSBjYWxsaW5nIGNhbmNlbC5cbiAgICAgKi9cbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgLy8gSU5WQVJJQU5UOiB0aGUgcmVzdWx0IG9mIHRoaXNbY2FuY2VsSW1wbFN5bV0gYW5kIHRoZSBiYXJyaWVyIGRvIG5vdCBldmVyIHJlamVjdC5cbiAgICAgICAgICAgIC8vIFVuZm9ydHVuYXRlbHkgbWFjT1MgSGlnaCBTaWVycmEgZG9lcyBub3Qgc3VwcG9ydCBQcm9taXNlLmFsbFNldHRsZWQuXG4gICAgICAgICAgICBQcm9taXNlLmFsbChbXG4gICAgICAgICAgICAgICAgdGhpc1tjYW5jZWxJbXBsU3ltXShuZXcgQ2FuY2VsRXJyb3IoXCJQcm9taXNlIGNhbmNlbGxlZC5cIiwgeyBjYXVzZSB9KSksXG4gICAgICAgICAgICAgICAgY3VycmVudEJhcnJpZXIodGhpcylcbiAgICAgICAgICAgIF0pLnRoZW4oKCkgPT4gcmVzb2x2ZSgpLCAoKSA9PiByZXNvbHZlKCkpO1xuICAgICAgICB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBCaW5kcyBwcm9taXNlIGNhbmNlbGxhdGlvbiB0byB0aGUgYWJvcnQgZXZlbnQgb2YgdGhlIGdpdmVuIHtAbGluayBBYm9ydFNpZ25hbH0uXG4gICAgICogSWYgdGhlIHNpZ25hbCBoYXMgYWxyZWFkeSBhYm9ydGVkLCB0aGUgcHJvbWlzZSB3aWxsIGJlIGNhbmNlbGxlZCBpbW1lZGlhdGVseS5cbiAgICAgKiBXaGVuIGVpdGhlciBjb25kaXRpb24gaXMgdmVyaWZpZWQsIHRoZSBjYW5jZWxsYXRpb24gY2F1c2Ugd2lsbCBiZSBzZXRcbiAgICAgKiB0byB0aGUgc2lnbmFsJ3MgYWJvcnQgcmVhc29uIChzZWUge0BsaW5rIEFib3J0U2lnbmFsI3JlYXNvbn0pLlxuICAgICAqXG4gICAgICogSGFzIG5vIGVmZmVjdCBpZiBjYWxsZWQgKG9yIGlmIHRoZSBzaWduYWwgYWJvcnRzKSBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQuXG4gICAgICogT25seSB0aGUgZmlyc3Qgc2lnbmFsIHRvIGFib3J0IHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXG4gICAgICpcbiAgICAgKiBGb3IgbW9yZSBkZXRhaWxzIGFib3V0IHRoZSBjYW5jZWxsYXRpb24gcHJvY2VzcyxcbiAgICAgKiBzZWUge0BsaW5rIGNhbmNlbH0gYW5kIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYCBjb25zdHJ1Y3Rvci5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGVuYWJsZXMgYGF3YWl0YGluZyBjYW5jZWxsYWJsZSBwcm9taXNlcyB3aXRob3V0IGhhdmluZ1xuICAgICAqIHRvIHN0b3JlIHRoZW0gZm9yIGZ1dHVyZSBjYW5jZWxsYXRpb24sIGUuZy46XG4gICAgICogYGBgdHNcbiAgICAgKiBhd2FpdCBsb25nUnVubmluZ09wZXJhdGlvbigpLmNhbmNlbE9uKHNpZ25hbCk7XG4gICAgICogYGBgXG4gICAgICogaW5zdGVhZCBvZjpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBwcm9taXNlVG9CZUNhbmNlbGxlZCA9IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCk7XG4gICAgICogYXdhaXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQ7XG4gICAgICogYGBgXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGlzIHByb21pc2UsIGZvciBtZXRob2QgY2hhaW5pbmcuXG4gICAgICovXG4gICAgY2FuY2VsT24oc2lnbmFsOiBBYm9ydFNpZ25hbCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGlmIChzaWduYWwuYWJvcnRlZCkge1xuICAgICAgICAgICAgdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKVxuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc2lnbmFsLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCB0aGlzLmNhbmNlbChzaWduYWwucmVhc29uKSwge2NhcHR1cmU6IHRydWV9KTtcbiAgICAgICAgfVxuXG4gICAgICAgIHJldHVybiB0aGlzO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEF0dGFjaGVzIGNhbGxiYWNrcyBmb3IgdGhlIHJlc29sdXRpb24gYW5kL29yIHJlamVjdGlvbiBvZiB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqXG4gICAgICogQHBhcmFtIG9uZnVsZmlsbGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVzb2x2ZWQuXG4gICAgICogQHBhcmFtIG9ucmVqZWN0ZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZWplY3RlZC5cbiAgICAgKiBAcmV0dXJucyBBIGBDYW5jZWxsYWJsZVByb21pc2VgIGZvciB0aGUgY29tcGxldGlvbiBvZiB3aGljaGV2ZXIgY2FsbGJhY2sgaXMgZXhlY3V0ZWQuXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgaXMgaG9va2VkIHVwIHRvIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24gcmVxdWVzdHMgdXAgdGhlIGNoYWluLCBidXQgbm90IGRvd246XG4gICAgICpcbiAgICAgKiAgIC0gaWYgdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgYmUgaW52b2tlZCB3aXRoIGEgYENhbmNlbEVycm9yYFxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2UgX3dpbGwgcmVzb2x2ZSByZWd1bGFybHlfIHdpdGggaXRzIHJlc3VsdDtcbiAgICAgKiAgIC0gY29udmVyc2VseSwgaWYgdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCBfdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCB0b287X1xuICAgICAqICAgICB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBzdGlsbCBiZSBpbnZva2VkIHdpdGggdGhlIHBhcmVudCdzIGBDYW5jZWxFcnJvcmAsXG4gICAgICogICAgIGJ1dCBpdHMgcmVzdWx0IHdpbGwgYmUgZGlzY2FyZGVkXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGEgYENhbmNlbEVycm9yYCBhcyB3ZWxsLlxuICAgICAqXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cbiAgICAgKlxuICAgICAqIElmIGVpdGhlciBjYWxsYmFjayByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqL1xuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4ge1xuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUudGhlbiBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gTk9URTogVHlwZVNjcmlwdCdzIGJ1aWx0LWluIHR5cGUgZm9yIHRoZW4gaXMgYnJva2VuLFxuICAgICAgICAvLyBhcyBpdCBhbGxvd3Mgc3BlY2lmeWluZyBhbiBhcmJpdHJhcnkgVFJlc3VsdDEgIT0gVCBldmVuIHdoZW4gb25mdWxmaWxsZWQgaXMgbm90IGEgZnVuY3Rpb24uXG4gICAgICAgIC8vIFdlIGNhbm5vdCBmaXggaXQgaWYgd2Ugd2FudCB0byBDYW5jZWxsYWJsZVByb21pc2UgdG8gaW1wbGVtZW50IFByb21pc2VMaWtlPFQ+LlxuXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZ1bGZpbGxlZCkpIHsgb25mdWxmaWxsZWQgPSBpZGVudGl0eSBhcyBhbnk7IH1cbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9ucmVqZWN0ZWQpKSB7IG9ucmVqZWN0ZWQgPSB0aHJvd2VyOyB9XG5cbiAgICAgICAgaWYgKG9uZnVsZmlsbGVkID09PSBpZGVudGl0eSAmJiBvbnJlamVjdGVkID09IHRocm93ZXIpIHtcbiAgICAgICAgICAgIC8vIFNob3J0Y3V0IGZvciB0cml2aWFsIGFyZ3VtZW50cy5cbiAgICAgICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlKSA9PiByZXNvbHZlKHRoaXMgYXMgYW55KSk7XG4gICAgICAgIH1cblxuICAgICAgICBjb25zdCBiYXJyaWVyOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiA9IHt9O1xuICAgICAgICB0aGlzW2JhcnJpZXJTeW1dID0gYmFycmllcjtcblxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUUmVzdWx0MSB8IFRSZXN1bHQyPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICB2b2lkIHN1cGVyLnRoZW4oXG4gICAgICAgICAgICAgICAgKHZhbHVlKSA9PiB7XG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XG5cbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25mdWxmaWxsZWQhKHZhbHVlKSk7XG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9LFxuICAgICAgICAgICAgICAgIChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XG5cbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25yZWplY3RlZCEocmVhc29uKSk7XG4gICAgICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICApO1xuICAgICAgICB9LCBhc3luYyAoY2F1c2U/KSA9PiB7XG4gICAgICAgICAgICAvL2NhbmNlbGxlZCA9IHRydWU7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHJldHVybiBvbmNhbmNlbGxlZD8uKGNhdXNlKTtcbiAgICAgICAgICAgIH0gZmluYWxseSB7XG4gICAgICAgICAgICAgICAgYXdhaXQgdGhpcy5jYW5jZWwoY2F1c2UpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIGZvciBvbmx5IHRoZSByZWplY3Rpb24gb2YgdGhlIFByb21pc2UuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqXG4gICAgICogSXQgaXMgZXF1aXZhbGVudCB0b1xuICAgICAqIGBgYHRzXG4gICAgICogY2FuY2VsbGFibGVQcm9taXNlLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XG4gICAgICogYGBgXG4gICAgICogYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBgb25yZWplY3RlZGAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXG4gICAgICovXG4gICAgY2F0Y2g8VFJlc3VsdCA9IG5ldmVyPihvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gKFByb21pc2VMaWtlPFRSZXN1bHQ+IHwgVFJlc3VsdCkpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VCB8IFRSZXN1bHQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIHRoYXQgaXMgaW52b2tlZCB3aGVuIHRoZSBDYW5jZWxsYWJsZVByb21pc2UgaXMgc2V0dGxlZCAoZnVsZmlsbGVkIG9yIHJlamVjdGVkKS4gVGhlXG4gICAgICogcmVzb2x2ZWQgdmFsdWUgY2Fubm90IGJlIGFjY2Vzc2VkIG9yIG1vZGlmaWVkIGZyb20gdGhlIGNhbGxiYWNrLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgc2V0dGxlIGluIHRoZSBzYW1lIHN0YXRlIGFzIHRoZSBvcmlnaW5hbCBvbmVcbiAgICAgKiBhZnRlciB0aGUgcHJvdmlkZWQgY2FsbGJhY2sgaGFzIGNvbXBsZXRlZCBleGVjdXRpb24sXG4gICAgICogdW5sZXNzIHRoZSBjYWxsYmFjayB0aHJvd3Mgb3IgcmV0dXJucyBhIHJlamVjdGluZyBwcm9taXNlLFxuICAgICAqIGluIHdoaWNoIGNhc2UgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3QgYXMgd2VsbC5cbiAgICAgKlxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXG4gICAgICogT25jZSB0aGUgcGFyZW50IHByb21pc2Ugc2V0dGxlcywgdGhlIGBvbmZpbmFsbHlgIGNhbGxiYWNrIHdpbGwgcnVuLFxuICAgICAqIF9ldmVuIGFmdGVyIHRoZSByZXR1cm5lZCBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZDpfXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxuICAgICAqXG4gICAgICogVGhpcyBtZXRob2QgaXMgaW1wbGVtZW50ZWQgaW4gdGVybXMgb2Yge0BsaW5rIHRoZW59IGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxuICAgICAqIEl0IGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcbiAgICAgKiB3aWxsIHByb3BhZ2F0ZSB1cCB0aGUgY2hhaW4gdG8gdGhlIHBhcmVudCBwcm9taXNlLFxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBgb25maW5hbGx5YCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cbiAgICAgKi9cbiAgICBmaW5hbGx5KG9uZmluYWxseT86ICgoKSA9PiB2b2lkKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLmZpbmFsbHkgY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZpbmFsbHkpKSB7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy50aGVuKG9uZmluYWxseSwgb25maW5hbGx5LCBvbmNhbmNlbGxlZCk7XG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gdGhpcy50aGVuKFxuICAgICAgICAgICAgKHZhbHVlKSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB2YWx1ZSksXG4gICAgICAgICAgICAocmVhc29uPykgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4geyB0aHJvdyByZWFzb247IH0pLFxuICAgICAgICAgICAgb25jYW5jZWxsZWQsXG4gICAgICAgICk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogV2UgdXNlIHRoZSBgW1N5bWJvbC5zcGVjaWVzXWAgc3RhdGljIHByb3BlcnR5LCBpZiBhdmFpbGFibGUsXG4gICAgICogdG8gZGlzYWJsZSB0aGUgYnVpbHQtaW4gYXV0b21hdGljIHN1YmNsYXNzaW5nIGZlYXR1cmVzIGZyb20ge0BsaW5rIFByb21pc2V9LlxuICAgICAqIEl0IGlzIGNyaXRpY2FsIGZvciBwZXJmb3JtYW5jZSByZWFzb25zIHRoYXQgZXh0ZW5kZXJzIGRvIG5vdCBvdmVycmlkZSB0aGlzLlxuICAgICAqIE9uY2UgdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmdcbiAgICAgKiBpcyBlaXRoZXIgYWNjZXB0ZWQgb3IgcmV0aXJlZCwgdGhpcyBpbXBsZW1lbnRhdGlvbiB3aWxsIGhhdmUgdG8gYmUgcmV2aXNlZCBhY2NvcmRpbmdseS5cbiAgICAgKlxuICAgICAqIEBpZ25vcmVcbiAgICAgKiBAaW50ZXJuYWxcbiAgICAgKi9cbiAgICBzdGF0aWMgZ2V0IFtzcGVjaWVzXSgpIHtcbiAgICAgICAgcmV0dXJuIFByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlLCBvciByZWplY3RlZCB3aGVuIGFueSBQcm9taXNlIGlzIHJlamVjdGVkLlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbGw8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD5bXT47XG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogQXdhaXRlZDxUW1BdPjsgfT47XG4gICAgc3RhdGljIGFsbDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSBvciByZWplY3QuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VD4+W10+O1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFRbUF0+PjsgfT47XG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGxTZXR0bGVkKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVGhlIGFueSBmdW5jdGlvbiByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGlzIGZ1bGZpbGxlZCBieSB0aGUgZmlyc3QgZ2l2ZW4gcHJvbWlzZSB0byBiZSBmdWxmaWxsZWQsXG4gICAgICogb3IgcmVqZWN0ZWQgd2l0aCBhbiBBZ2dyZWdhdGVFcnJvciBjb250YWluaW5nIGFuIGFycmF5IG9mIHJlamVjdGlvbiByZWFzb25zXG4gICAgICogaWYgYWxsIG9mIHRoZSBnaXZlbiBwcm9taXNlcyBhcmUgcmVqZWN0ZWQuXG4gICAgICogSXQgcmVzb2x2ZXMgYWxsIGVsZW1lbnRzIG9mIHRoZSBwYXNzZWQgaXRlcmFibGUgdG8gcHJvbWlzZXMgYXMgaXQgcnVucyB0aGlzIGFsZ29yaXRobS5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYW55PFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFueShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIGFyZSByZXNvbHZlZCBvciByZWplY3RlZC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmFjZTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUW251bWJlcl0+PjtcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yYWNlKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGNhbmNlbGxlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBjYW5jZWw8VCA9IG5ldmVyPihjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGNvbnN0IHAgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcbiAgICAgICAgcC5jYW5jZWwoY2F1c2UpO1xuICAgICAgICByZXR1cm4gcDtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGNhbmNlbHNcbiAgICAgKiBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIGNhdXNlLlxuICAgICAqXG4gICAgICogSWYgdGhlIHtAbGluayBBYm9ydFNpZ25hbC50aW1lb3V0fSBmYWN0b3J5IG1ldGhvZCBpcyBhdmFpbGFibGUsXG4gICAgICogaXQgaXMgdXNlZCB0byBiYXNlIHRoZSB0aW1lb3V0IG9uIF9hY3RpdmVfIHRpbWUgcmF0aGVyIHRoYW4gX2VsYXBzZWRfIHRpbWUuXG4gICAgICogT3RoZXJ3aXNlLCBgdGltZW91dGAgZmFsbHMgYmFjayB0byB7QGxpbmsgc2V0VGltZW91dH0uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgdGltZW91dDxUID0gbmV2ZXI+KG1pbGxpc2Vjb25kczogbnVtYmVyLCBjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcbiAgICAgICAgaWYgKEFib3J0U2lnbmFsICYmIHR5cGVvZiBBYm9ydFNpZ25hbCA9PT0gJ2Z1bmN0aW9uJyAmJiBBYm9ydFNpZ25hbC50aW1lb3V0ICYmIHR5cGVvZiBBYm9ydFNpZ25hbC50aW1lb3V0ID09PSAnZnVuY3Rpb24nKSB7XG4gICAgICAgICAgICBBYm9ydFNpZ25hbC50aW1lb3V0KG1pbGxpc2Vjb25kcykuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpLCBtaWxsaXNlY29uZHMpO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgc2xlZXAobWlsbGlzZWNvbmRzOiBudW1iZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlclxuICAgICAqIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBzbGVlcDxUPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XG4gICAgc3RhdGljIHNsZWVwPFQgPSB2b2lkPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU/OiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gcmVzb2x2ZSh2YWx1ZSEpLCBtaWxsaXNlY29uZHMpO1xuICAgICAgICB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlamVjdGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHJlYXNvbi5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZWplY3Q8VCA9IG5ldmVyPihyZWFzb24/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoXywgcmVqZWN0KSA9PiByZWplY3QocmVhc29uKSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZSgpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyByZXNvbHZlPFQgPSB2b2lkPih2YWx1ZT86IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PiB7XG4gICAgICAgIGlmICh2YWx1ZSBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkge1xuICAgICAgICAgICAgLy8gT3B0aW1pc2UgZm9yIGNhbmNlbGxhYmxlIHByb21pc2VzLlxuICAgICAgICAgICAgcmV0dXJuIHZhbHVlO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPGFueT4oKHJlc29sdmUpID0+IHJlc29sdmUodmFsdWUpKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSBhbmQgcmV0dXJucyBpdCBpbiBhbiBvYmplY3QsIGFsb25nIHdpdGggaXRzIHJlc29sdmUgYW5kIHJlamVjdCBmdW5jdGlvbnNcbiAgICAgKiBhbmQgYSBnZXR0ZXIvc2V0dGVyIGZvciB0aGUgY2FuY2VsbGF0aW9uIGNhbGxiYWNrLlxuICAgICAqXG4gICAgICogVGhpcyBtZXRob2QgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyB3aXRoUmVzb2x2ZXJzPFQ+KCk6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xuICAgICAgICBsZXQgcmVzdWx0OiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0geyBvbmNhbmNlbGxlZDogbnVsbCB9IGFzIGFueTtcbiAgICAgICAgcmVzdWx0LnByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgIHJlc3VsdC5yZXNvbHZlID0gcmVzb2x2ZTtcbiAgICAgICAgICAgIHJlc3VsdC5yZWplY3QgPSByZWplY3Q7XG4gICAgICAgIH0sIChjYXVzZT86IGFueSkgPT4geyByZXN1bHQub25jYW5jZWxsZWQ/LihjYXVzZSk7IH0pO1xuICAgICAgICByZXR1cm4gcmVzdWx0O1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSBjYW5jZWxsYXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cbiAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIHJlc3VsdGluZyBmdW5jdGlvbiBkb2VzIG5vdCByZWplY3QuXG4gKi9cbmZ1bmN0aW9uIGNhbmNlbGxlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpIHtcbiAgICBsZXQgY2FuY2VsbGF0aW9uUHJvbWlzZTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0gdW5kZWZpbmVkO1xuXG4gICAgcmV0dXJuIChyZWFzb246IENhbmNlbEVycm9yKTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0+IHtcbiAgICAgICAgaWYgKCFzdGF0ZS5zZXR0bGVkKSB7XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHN0YXRlLnJlYXNvbiA9IHJlYXNvbjtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XG5cbiAgICAgICAgICAgIC8vIEF0dGFjaCBhbiBlcnJvciBoYW5kbGVyIHRoYXQgaWdub3JlcyB0aGlzIHNwZWNpZmljIHJlamVjdGlvbiByZWFzb24gYW5kIG5vdGhpbmcgZWxzZS5cbiAgICAgICAgICAgIC8vIEluIHRoZW9yeSwgYSBzYW5lIHVuZGVybHlpbmcgaW1wbGVtZW50YXRpb24gYXQgdGhpcyBwb2ludFxuICAgICAgICAgICAgLy8gc2hvdWxkIGFsd2F5cyByZWplY3Qgd2l0aCBvdXIgY2FuY2VsbGF0aW9uIHJlYXNvbixcbiAgICAgICAgICAgIC8vIGhlbmNlIHRoZSBoYW5kbGVyIHdpbGwgbmV2ZXIgdGhyb3cuXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucHJvdG90eXBlLnRoZW4uY2FsbChwcm9taXNlLnByb21pc2UsIHVuZGVmaW5lZCwgKGVycikgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChlcnIgIT09IHJlYXNvbikge1xuICAgICAgICAgICAgICAgICAgICB0aHJvdyBlcnI7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBJZiByZWFzb24gaXMgbm90IHNldCwgdGhlIHByb21pc2UgcmVzb2x2ZWQgcmVndWxhcmx5LCBoZW5jZSB3ZSBtdXN0IG5vdCBjYWxsIG9uY2FuY2VsbGVkLlxuICAgICAgICAvLyBJZiBvbmNhbmNlbGxlZCBpcyB1bnNldCwgbm8gbmVlZCB0byBnbyBhbnkgZnVydGhlci5cbiAgICAgICAgaWYgKCFzdGF0ZS5yZWFzb24gfHwgIXByb21pc2Uub25jYW5jZWxsZWQpIHsgcmV0dXJuOyB9XG5cbiAgICAgICAgY2FuY2VsbGF0aW9uUHJvbWlzZSA9IG5ldyBQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHJlc29sdmUocHJvbWlzZS5vbmNhbmNlbGxlZCEoc3RhdGUucmVhc29uIS5jYXVzZSkpO1xuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XG4gICAgICAgIH0pO1xuXG4gICAgICAgIC8vIFVuc2V0IG9uY2FuY2VsbGVkIHRvIHByZXZlbnQgcmVwZWF0ZWQgY2FsbHMuXG4gICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBudWxsO1xuXG4gICAgICAgIHJldHVybiBjYW5jZWxsYXRpb25Qcm9taXNlO1xuICAgIH1cbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZXNvbHV0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKi9cbmZ1bmN0aW9uIHJlc29sdmVyRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+IHtcbiAgICByZXR1cm4gKHZhbHVlKSA9PiB7XG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnJlc29sdmluZyA9IHRydWU7XG5cbiAgICAgICAgaWYgKHZhbHVlID09PSBwcm9taXNlLnByb21pc2UpIHtcbiAgICAgICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7IHJldHVybjsgfVxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgVHlwZUVycm9yKFwiQSBwcm9taXNlIGNhbm5vdCBiZSByZXNvbHZlZCB3aXRoIGl0c2VsZi5cIikpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHZhbHVlICE9IG51bGwgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcgfHwgdHlwZW9mIHZhbHVlID09PSAnZnVuY3Rpb24nKSkge1xuICAgICAgICAgICAgbGV0IHRoZW46IGFueTtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgdGhlbiA9ICh2YWx1ZSBhcyBhbnkpLnRoZW47XG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgICAgICBwcm9taXNlLnJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUodGhlbikpIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICBsZXQgY2FuY2VsID0gKHZhbHVlIGFzIGFueSkuY2FuY2VsO1xuICAgICAgICAgICAgICAgICAgICBpZiAoaXNDYWxsYWJsZShjYW5jZWwpKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCBvbmNhbmNlbGxlZCA9IChjYXVzZT86IGFueSkgPT4ge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XG4gICAgICAgICAgICAgICAgICAgICAgICB9O1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKHN0YXRlLnJlYXNvbikge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIElmIGFscmVhZHkgY2FuY2VsbGVkLCBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uLlxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIGNhbmNlbGxlciBhbGdvcml0aG0gZG9lcyBub3QgcmVqZWN0XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gc28gaXQgY2FuIGJlIGRpc2NhcmRlZCBzYWZlbHkuXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCBjYW5jZWxsZXJGb3IoeyAuLi5wcm9taXNlLCBvbmNhbmNlbGxlZCB9LCBzdGF0ZSkoc3RhdGUucmVhc29uKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG9uY2FuY2VsbGVkO1xuICAgICAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSBjYXRjaCB7fVxuXG4gICAgICAgICAgICAgICAgY29uc3QgbmV3U3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xuICAgICAgICAgICAgICAgICAgICByb290OiBzdGF0ZS5yb290LFxuICAgICAgICAgICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxuICAgICAgICAgICAgICAgICAgICBnZXQgc2V0dGxlZCgpIHsgcmV0dXJuIHRoaXMucm9vdC5zZXR0bGVkIH0sXG4gICAgICAgICAgICAgICAgICAgIHNldCBzZXR0bGVkKHZhbHVlKSB7IHRoaXMucm9vdC5zZXR0bGVkID0gdmFsdWU7IH0sXG4gICAgICAgICAgICAgICAgICAgIGdldCByZWFzb24oKSB7IHJldHVybiB0aGlzLnJvb3QucmVhc29uIH1cbiAgICAgICAgICAgICAgICB9O1xuXG4gICAgICAgICAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBuZXdTdGF0ZSk7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseSh0aGVuLCB2YWx1ZSwgW3Jlc29sdmVyRm9yKHByb21pc2UsIG5ld1N0YXRlKSwgcmVqZWN0b3JdKTtcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmV0dXJuOyAvLyBJTVBPUlRBTlQhXG4gICAgICAgICAgICB9XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgIHByb21pc2UucmVzb2x2ZSh2YWx1ZSk7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZWplY3Rpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cbiAqL1xuZnVuY3Rpb24gcmVqZWN0b3JGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3Ige1xuICAgIHJldHVybiAocmVhc29uPykgPT4ge1xuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xuXG4gICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGlmIChyZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBzdGF0ZS5yZWFzb24gaW5zdGFuY2VvZiBDYW5jZWxFcnJvciAmJiBPYmplY3QuaXMocmVhc29uLmNhdXNlLCBzdGF0ZS5yZWFzb24uY2F1c2UpKSB7XG4gICAgICAgICAgICAgICAgICAgIC8vIFN3YWxsb3cgbGF0ZSByZWplY3Rpb25zIHRoYXQgYXJlIENhbmNlbEVycm9ycyB3aG9zZSBjYW5jZWxsYXRpb24gY2F1c2UgaXMgdGhlIHNhbWUgYXMgb3Vycy5cbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0gY2F0Y2gge31cblxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24pKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcbiAgICAgICAgfVxuICAgIH1cbn1cblxuLyoqXG4gKiBDYW5jZWxzIGFsbCB2YWx1ZXMgaW4gYW4gYXJyYXkgdGhhdCBsb29rIGxpa2UgY2FuY2VsbGFibGUgdGhlbmFibGVzLlxuICogUmV0dXJucyBhIHByb21pc2UgdGhhdCBmdWxmaWxscyBvbmNlIGFsbCBjYW5jZWxsYXRpb24gcHJvY2VkdXJlcyBmb3IgdGhlIGdpdmVuIHZhbHVlcyBoYXZlIHNldHRsZWQuXG4gKi9cbmZ1bmN0aW9uIGNhbmNlbEFsbChwYXJlbnQ6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgdmFsdWVzOiBhbnlbXSwgY2F1c2U/OiBhbnkpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBjb25zdCByZXN1bHRzOiBQcm9taXNlPHZvaWQ+W10gPSBbXTtcblxuICAgIGZvciAoY29uc3QgdmFsdWUgb2YgdmFsdWVzKSB7XG4gICAgICAgIGxldCBjYW5jZWw6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcjtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZSh2YWx1ZS50aGVuKSkgeyBjb250aW51ZTsgfVxuICAgICAgICAgICAgY2FuY2VsID0gdmFsdWUuY2FuY2VsO1xuICAgICAgICAgICAgaWYgKCFpc0NhbGxhYmxlKGNhbmNlbCkpIHsgY29udGludWU7IH1cbiAgICAgICAgfSBjYXRjaCB7IGNvbnRpbnVlOyB9XG5cbiAgICAgICAgbGV0IHJlc3VsdDogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgcmVzdWx0ID0gUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XG4gICAgICAgICAgICBjb250aW51ZTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICghcmVzdWx0KSB7IGNvbnRpbnVlOyB9XG4gICAgICAgIHJlc3VsdHMucHVzaChcbiAgICAgICAgICAgIChyZXN1bHQgaW5zdGFuY2VvZiBQcm9taXNlICA/IHJlc3VsdCA6IFByb21pc2UucmVzb2x2ZShyZXN1bHQpKS5jYXRjaCgocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIHJlYXNvbiwgXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICk7XG4gICAgfVxuXG4gICAgcmV0dXJuIFByb21pc2UuYWxsKHJlc3VsdHMpIGFzIGFueTtcbn1cblxuLyoqXG4gKiBSZXR1cm5zIGl0cyBhcmd1bWVudC5cbiAqL1xuZnVuY3Rpb24gaWRlbnRpdHk8VD4oeDogVCk6IFQge1xuICAgIHJldHVybiB4O1xufVxuXG4vKipcbiAqIFRocm93cyBpdHMgYXJndW1lbnQuXG4gKi9cbmZ1bmN0aW9uIHRocm93ZXIocmVhc29uPzogYW55KTogbmV2ZXIge1xuICAgIHRocm93IHJlYXNvbjtcbn1cblxuLyoqXG4gKiBBdHRlbXB0cyB2YXJpb3VzIHN0cmF0ZWdpZXMgdG8gY29udmVydCBhbiBlcnJvciB0byBhIHN0cmluZy5cbiAqL1xuZnVuY3Rpb24gZXJyb3JNZXNzYWdlKGVycjogYW55KTogc3RyaW5nIHtcbiAgICB0cnkge1xuICAgICAgICBpZiAoZXJyIGluc3RhbmNlb2YgRXJyb3IgfHwgdHlwZW9mIGVyciAhPT0gJ29iamVjdCcgfHwgZXJyLnRvU3RyaW5nICE9PSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nKSB7XG4gICAgICAgICAgICByZXR1cm4gXCJcIiArIGVycjtcbiAgICAgICAgfVxuICAgIH0gY2F0Y2gge31cblxuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiBKU09OLnN0cmluZ2lmeShlcnIpO1xuICAgIH0gY2F0Y2gge31cblxuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nLmNhbGwoZXJyKTtcbiAgICB9IGNhdGNoIHt9XG5cbiAgICByZXR1cm4gXCI8Y291bGQgbm90IGNvbnZlcnQgZXJyb3IgdG8gc3RyaW5nPlwiO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYmFycmllciBwcm9taXNlIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS4gSWYgbmVjZXNzYXJ5LCBpbml0aWFsaXNlcyB0aGUgYmFycmllci5cbiAqL1xuZnVuY3Rpb24gY3VycmVudEJhcnJpZXI8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+KTogUHJvbWlzZTx2b2lkPiB7XG4gICAgbGV0IHB3cjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSBwcm9taXNlW2JhcnJpZXJTeW1dID8/IHt9O1xuICAgIGlmICghKCdwcm9taXNlJyBpbiBwd3IpKSB7XG4gICAgICAgIE9iamVjdC5hc3NpZ24ocHdyLCBwcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPigpKTtcbiAgICB9XG4gICAgaWYgKHByb21pc2VbYmFycmllclN5bV0gPT0gbnVsbCkge1xuICAgICAgICBwd3IucmVzb2x2ZSEoKTtcbiAgICAgICAgcHJvbWlzZVtiYXJyaWVyU3ltXSA9IHB3cjtcbiAgICB9XG4gICAgcmV0dXJuIHB3ci5wcm9taXNlITtcbn1cblxuLy8gUG9seWZpbGwgUHJvbWlzZS53aXRoUmVzb2x2ZXJzLlxubGV0IHByb21pc2VXaXRoUmVzb2x2ZXJzID0gUHJvbWlzZS53aXRoUmVzb2x2ZXJzO1xuaWYgKHByb21pc2VXaXRoUmVzb2x2ZXJzICYmIHR5cGVvZiBwcm9taXNlV2l0aFJlc29sdmVycyA9PT0gJ2Z1bmN0aW9uJykge1xuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gcHJvbWlzZVdpdGhSZXNvbHZlcnMuYmluZChQcm9taXNlKTtcbn0gZWxzZSB7XG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBmdW5jdGlvbiA8VD4oKTogUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4ge1xuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBQcm9taXNlPFQ+KChyZXMsIHJlaikgPT4geyByZXNvbHZlID0gcmVzOyByZWplY3QgPSByZWo7IH0pO1xuICAgICAgICByZXR1cm4geyBwcm9taXNlLCByZXNvbHZlLCByZWplY3QgfTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2xpcGJvYXJkKTtcblxuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XG5jb25zdCBDbGlwYm9hcmRUZXh0ID0gMTtcblxuLyoqXG4gKiBTZXRzIHRoZSB0ZXh0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKlxuICogQHBhcmFtIHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cbiAqIEByZXR1cm4gQSBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2hlbiB0aGUgb3BlcmF0aW9uIGlzIHN1Y2Nlc3NmdWwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQ6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpOiBQcm9taXNlPHN0cmluZz4ge1xuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xufVxuIiwgImltcG9ydCB7IFN0cmVhbSB9IGZyb20gXCIuL3N0cmVhbS5qc1wiO1xuXG4vKiogT3B0aW9ucyBmb3IgbG9hZGluZyBhIGJvdW5kZWQgbWVkaWEgZmlsZSBvdmVyIGEgcmVnaXN0ZXJlZCBXYWlscyBzdHJlYW0uICovXG5leHBvcnQgaW50ZXJmYWNlIFNvdXJjZU9wdGlvbnMge1xuICAgIC8qKiBNYXhpbXVtIGRvd25sb2FkZWQgYnl0ZXMgcGVyIHNvdXJjZS4gRGVmYXVsdHMgdG8gMzIgTWlCLiBNdXN0IGJlIGEgcG9zaXRpdmUgc2FmZSBpbnRlZ2VyLiAqL1xuICAgIG1heEJ5dGVzPzogbnVtYmVyO1xuICAgIC8qKiBDYW5jZWxzIHRoaXMgbG9hZCB3aXRob3V0IGNsZWFyaW5nIGEgcHJldmlvdXNseSBsb2FkZWQgc291cmNlLiAqL1xuICAgIHNpZ25hbD86IEFib3J0U2lnbmFsO1xufVxuXG5pbnRlcmZhY2UgU291cmNlU3RhdGUge1xuICAgIGdlbmVyYXRpb246IG51bWJlcjtcbiAgICBjb250cm9sbGVyPzogQWJvcnRDb250cm9sbGVyO1xuICAgIG9iamVjdFVSTD86IHN0cmluZztcbn1cblxuLy8gVGhlIG5wbSBhbmQgYnVuZGxlZCBydW50aW1lcyBjYW4gY29leGlzdCBpbiBvbmUgcGFnZS4gU3RvcmUgb3duZXJzaGlwIG9uIHRoZVxuLy8gZWxlbWVudCB1bmRlciBhIHNoYXJlZCBzeW1ib2wgc28gZWl0aGVyIGNvcHkgY2FuIHJlcGxhY2Ugb3IgY2xlYXIgaXRzIHNvdXJjZS5cbmNvbnN0IHNvdXJjZUtleSA9IFN5bWJvbC5mb3IoJ3dhaWxzLm1lZGlhLnNvdXJjZScpO1xudHlwZSBNYW5hZ2VkRWxlbWVudCA9IEhUTUxNZWRpYUVsZW1lbnQgJiB7IFtzb3VyY2VLZXldPzogU291cmNlU3RhdGUgfTtcbmNvbnN0IGRlZmF1bHRNYXhCeXRlcyA9IDMyICogMTAyNCAqIDEwMjQ7XG5cbmZ1bmN0aW9uIGNoZWNrQWJvcnRlZChzaWduYWw6IEFib3J0U2lnbmFsKTogdm9pZCB7XG4gICAgaWYgKHNpZ25hbC5hYm9ydGVkKSB7XG4gICAgICAgIHRocm93IHNpZ25hbC5yZWFzb24gPz8gbmV3IERPTUV4Y2VwdGlvbignTWVkaWEgbG9hZCBjYW5jZWxsZWQnLCAnQWJvcnRFcnJvcicpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBMb2FkcyBhIGNsaXAgZnJvbSBhIEdvIG1lZGlhLk5ld0hhbmRsZXIgcmVnaXN0ZXJlZCB3aXRoIGFwcC5IYW5kbGVTdHJlYW0sXG4gKiB0aGVuIGFzc2lnbnMgYSBibG9iIFVSTCB0byBhbiBhdWRpbyBvciB2aWRlbyBlbGVtZW50LiBuYW1lIGlzIGEgZmlsZSBwYXRoXG4gKiByZWxhdGl2ZSB0byB0aGF0IGhhbmRsZXIncyBmaWxlc3lzdGVtLCBub3QgYSBVUkwgb3Igb3BlcmF0aW5nLXN5c3RlbSBwYXRoLlxuICpcbiAqIERlc2t0b3AgdHJhbnNmZXJzIHVzZSBib3VuZGVkIFdhaWxzIHN0cmVhbSByZXNwb25zZXMsIGluY2x1ZGluZyBvbiBXZWJWaWV3MixcbiAqIHdpdGhvdXQgb3BlbmluZyBhIGxpc3RlbmVyLiBUaGUgY29tcGxldGUgZmlsZSBpcyByZXRhaW5lZCBiZWZvcmUgYXNzaWdubWVudDtcbiAqIHRoaXMgZW5hYmxlcyBsb2NhbCBwbGF5YmFjayBvbiBMaW51eCBidXQgaXMgbm90IHByb2dyZXNzaXZlIG1lZGlhIHBsYXliYWNrLlxuICogbWF4Qnl0ZXMgZGVmYXVsdHMgdG8gMzIgTWlCLCBhbmQgdGhlIEdvIGhhbmRsZXIgZW5mb3JjZXMgaXRzIG93biBsaW1pdCB0b28uXG4gKiBCcm93c2VyIGRlY29kaW5nIGFuZCBibG9iIGNvbnN0cnVjdGlvbiBjYW4gcmVxdWlyZSBhZGRpdGlvbmFsIG1lbW9yeS5cbiAqXG4gKiBBIG5ldyBjYWxsIGNhbmNlbHMgYW4gZWFybGllciBsb2FkIGZvciB0aGUgc2FtZSBlbGVtZW50LiBUaGUgZXhpc3Rpbmcgc291cmNlXG4gKiByZW1haW5zIHVudGlsIGEgcmVwbGFjZW1lbnQgc3VjY2VlZHMuIEFib3J0ZWQgbG9hZHMgcmVqZWN0IHdpdGggQWJvcnRFcnJvclxuICogKG9yIHRoZSBjYWxsZXIncyBhYm9ydCByZWFzb24pLiBDb21wbGV0aW9uIG1lYW5zIHRoZSBzb3VyY2Ugd2FzIGFzc2lnbmVkO1xuICogdXNlIG1lZGlhIGV2ZW50cyBhbmQgcGxheSgpIGZvciBkZWNvZGluZyBhbmQgcGxheWJhY2sgcmVzdWx0cy5cbiAqXG4gKiBDYWxsIENsZWFyU291cmNlIGJlZm9yZSBkaXNjYXJkaW5nIHRoZSBlbGVtZW50IHRvIHJlbGVhc2UgaXRzIGJsb2IuIFVzZSB0aGVzZVxuICogaGVscGVycyBjb25zaXN0ZW50bHkgZm9yIGFuIGVsZW1lbnQ7IHRoZXkgZG8gbm90IG9ic2VydmUgZGlyZWN0IHNyYyBjaGFuZ2VzLFxuICogc291cmNlIGNoaWxkcmVuIG9yIERPTSByZW1vdmFsLiBUaGlzIGFsc28gd29ya3Mgd2l0aCBkZXRhY2hlZCBBdWRpbyBvYmplY3RzLlxuICovXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gU2V0U291cmNlKGVsZW1lbnQ6IEhUTUxNZWRpYUVsZW1lbnQsIHN0cmVhbU5hbWU6IHN0cmluZywgbmFtZTogc3RyaW5nLCBvcHRpb25zOiBTb3VyY2VPcHRpb25zID0ge30pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBjb25zdCBtYXhCeXRlcyA9IG9wdGlvbnMubWF4Qnl0ZXMgPz8gZGVmYXVsdE1heEJ5dGVzO1xuICAgIGlmICghTnVtYmVyLmlzU2FmZUludGVnZXIobWF4Qnl0ZXMpIHx8IG1heEJ5dGVzIDw9IDApIHtcbiAgICAgICAgdGhyb3cgbmV3IFJhbmdlRXJyb3IoJ21heEJ5dGVzIG11c3QgYmUgYSBwb3NpdGl2ZSBzYWZlIGludGVnZXInKTtcbiAgICB9XG4gICAgaWYgKCFzdHJlYW1OYW1lLnRyaW0oKSB8fCAhbmFtZS50cmltKCkpIHtcbiAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcignQSBtZWRpYSBzdHJlYW0gYW5kIGZpbGUgbmFtZSBhcmUgcmVxdWlyZWQnKTtcbiAgICB9XG4gICAgaWYgKG9wdGlvbnMuc2lnbmFsKSBjaGVja0Fib3J0ZWQob3B0aW9ucy5zaWduYWwpO1xuXG4gICAgY29uc3QgbWFuYWdlZCA9IGVsZW1lbnQgYXMgTWFuYWdlZEVsZW1lbnQ7XG4gICAgY29uc3Qgc3RhdGUgPSBtYW5hZ2VkW3NvdXJjZUtleV0gPz8geyBnZW5lcmF0aW9uOiAwIH07XG4gICAgbWFuYWdlZFtzb3VyY2VLZXldID0gc3RhdGU7XG4gICAgc3RhdGUuY29udHJvbGxlcj8uYWJvcnQoKTtcbiAgICBjb25zdCBnZW5lcmF0aW9uID0gKytzdGF0ZS5nZW5lcmF0aW9uO1xuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgc3RhdGUuY29udHJvbGxlciA9IGNvbnRyb2xsZXI7XG4gICAgY29uc3QgY2FuY2VsID0gKCkgPT4gY29udHJvbGxlci5hYm9ydChvcHRpb25zLnNpZ25hbD8ucmVhc29uKTtcbiAgICBvcHRpb25zLnNpZ25hbD8uYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCBjYW5jZWwsIHsgb25jZTogdHJ1ZSB9KTtcbiAgICB0cnkge1xuICAgICAgICBjb25zdCBibG9iID0gYXdhaXQgcmVjZWl2ZUJsb2Ioc3RyZWFtTmFtZSwgbmFtZSwgbWF4Qnl0ZXMsIGNvbnRyb2xsZXIuc2lnbmFsKTtcbiAgICAgICAgY2hlY2tBYm9ydGVkKGNvbnRyb2xsZXIuc2lnbmFsKTtcbiAgICAgICAgaWYgKHN0YXRlLmdlbmVyYXRpb24gIT09IGdlbmVyYXRpb24pIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBET01FeGNlcHRpb24oJ01lZGlhIHNvdXJjZSByZXBsYWNlZCcsICdBYm9ydEVycm9yJyk7XG4gICAgICAgIH1cbiAgICAgICAgY29uc3Qgb2JqZWN0VVJMID0gVVJMLmNyZWF0ZU9iamVjdFVSTChibG9iKTtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGVsZW1lbnQuc3JjID0gb2JqZWN0VVJMO1xuICAgICAgICB9IGNhdGNoIChlcnJvcikge1xuICAgICAgICAgICAgVVJMLnJldm9rZU9iamVjdFVSTChvYmplY3RVUkwpO1xuICAgICAgICAgICAgdGhyb3cgZXJyb3I7XG4gICAgICAgIH1cbiAgICAgICAgaWYgKHN0YXRlLm9iamVjdFVSTCkgVVJMLnJldm9rZU9iamVjdFVSTChzdGF0ZS5vYmplY3RVUkwpO1xuICAgICAgICBzdGF0ZS5vYmplY3RVUkwgPSBvYmplY3RVUkw7XG4gICAgfSBmaW5hbGx5IHtcbiAgICAgICAgb3B0aW9ucy5zaWduYWw/LnJlbW92ZUV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgY2FuY2VsKTtcbiAgICAgICAgaWYgKHN0YXRlLmdlbmVyYXRpb24gPT09IGdlbmVyYXRpb24pIHN0YXRlLmNvbnRyb2xsZXIgPSB1bmRlZmluZWQ7XG4gICAgfVxufVxuXG4vKipcbiAqIENhbmNlbHMgYW4gb3V0c3RhbmRpbmcgU2V0U291cmNlIGxvYWQsIHJlc2V0cyB0aGUgbWVkaWEgZWxlbWVudCBhbmQgcmV2b2tlc1xuICogdGhlIGJsb2IgaXQgb3ducy4gQ2FsbCB0aGlzIG9uIGNvbXBvbmVudCB1bm1vdW50IG9yIGJlZm9yZSBkaXNjYXJkaW5nIGFuIEF1ZGlvXG4gKiBvYmplY3QuIEV4aXN0aW5nIHNvdXJjZSBjaGlsZHJlbiwgaWYgYW55LCBmb2xsb3cgbm9ybWFsIGJyb3dzZXIgc2VsZWN0aW9uIHJ1bGVzXG4gKiB3aGVuIHRoZSBzcmMgYXR0cmlidXRlIGlzIHJlbW92ZWQuIFJlcGVhdGVkIGNhbGxzIGhhdmUgbm8gZWZmZWN0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2xlYXJTb3VyY2UoZWxlbWVudDogSFRNTE1lZGlhRWxlbWVudCk6IHZvaWQge1xuICAgIGNvbnN0IG1hbmFnZWQgPSBlbGVtZW50IGFzIE1hbmFnZWRFbGVtZW50O1xuICAgIGNvbnN0IHN0YXRlID0gbWFuYWdlZFtzb3VyY2VLZXldO1xuICAgIGlmICghc3RhdGUpIHJldHVybjtcbiAgICArK3N0YXRlLmdlbmVyYXRpb247XG4gICAgc3RhdGUuY29udHJvbGxlcj8uYWJvcnQoKTtcbiAgICBkZWxldGUgbWFuYWdlZFtzb3VyY2VLZXldO1xuICAgIGlmIChzdGF0ZS5vYmplY3RVUkwpIHtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGVsZW1lbnQucGF1c2UoKTtcbiAgICAgICAgICAgIGVsZW1lbnQucmVtb3ZlQXR0cmlidXRlKCdzcmMnKTtcbiAgICAgICAgICAgIGVsZW1lbnQubG9hZCgpO1xuICAgICAgICB9IGZpbmFsbHkge1xuICAgICAgICAgICAgVVJMLnJldm9rZU9iamVjdFVSTChzdGF0ZS5vYmplY3RVUkwpO1xuICAgICAgICB9XG4gICAgfVxufVxuXG5mdW5jdGlvbiByZWNlaXZlQmxvYihzdHJlYW1OYW1lOiBzdHJpbmcsIG5hbWU6IHN0cmluZywgbWF4Qnl0ZXM6IG51bWJlciwgc2lnbmFsOiBBYm9ydFNpZ25hbCk6IFByb21pc2U8QmxvYj4ge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGNvbnN0IHN0cmVhbSA9IFN0cmVhbShzdHJlYW1OYW1lKTtcbiAgICAgICAgc3RyZWFtLmJpbmFyeVR5cGUgPSAnYXJyYXlidWZmZXInO1xuICAgICAgICBjb25zdCBwYXJ0czogQmxvYlBhcnRbXSA9IFtdO1xuICAgICAgICBsZXQgZXhwZWN0ZWQ6IG51bWJlciB8IHVuZGVmaW5lZDtcbiAgICAgICAgbGV0IGNvbnRlbnRUeXBlID0gJyc7XG4gICAgICAgIGxldCByZWNlaXZlZCA9IDA7XG4gICAgICAgIGxldCBmaW5pc2hlZCA9IGZhbHNlO1xuXG4gICAgICAgIGNvbnN0IGZpbmlzaCA9IChlcnJvcj86IHVua25vd24pID0+IHtcbiAgICAgICAgICAgIGlmIChmaW5pc2hlZCkgcmV0dXJuO1xuICAgICAgICAgICAgZmluaXNoZWQgPSB0cnVlO1xuICAgICAgICAgICAgc2lnbmFsLnJlbW92ZUV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgY2FuY2VsKTtcbiAgICAgICAgICAgIHN0cmVhbS5vbm9wZW4gPSBzdHJlYW0ub25tZXNzYWdlID0gc3RyZWFtLm9uZXJyb3IgPSBzdHJlYW0ub25jbG9zZSA9IG51bGw7XG4gICAgICAgICAgICBzdHJlYW0uY2xvc2UoKTtcbiAgICAgICAgICAgIGlmIChlcnJvciAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgICAgICAgICAgcGFydHMubGVuZ3RoID0gMDtcbiAgICAgICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG5ldyBCbG9iKHBhcnRzLCB7IHR5cGU6IGNvbnRlbnRUeXBlIH0pKTtcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChjYXVzZSkge1xuICAgICAgICAgICAgICAgICAgICByZWplY3QoY2F1c2UpO1xuICAgICAgICAgICAgICAgIH0gZmluYWxseSB7XG4gICAgICAgICAgICAgICAgICAgIHBhcnRzLmxlbmd0aCA9IDA7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICB9O1xuICAgICAgICBjb25zdCBjYW5jZWwgPSAoKSA9PiBmaW5pc2goc2lnbmFsLnJlYXNvbiA/PyBuZXcgRE9NRXhjZXB0aW9uKCdNZWRpYSBsb2FkIGNhbmNlbGxlZCcsICdBYm9ydEVycm9yJykpO1xuICAgICAgICBzaWduYWwuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCBjYW5jZWwsIHsgb25jZTogdHJ1ZSB9KTtcbiAgICAgICAgaWYgKHNpZ25hbC5hYm9ydGVkKSB7IGNhbmNlbCgpOyByZXR1cm47IH1cblxuICAgICAgICBzdHJlYW0ub25vcGVuID0gKCkgPT4ge1xuICAgICAgICAgICAgaWYgKGZpbmlzaGVkKSByZXR1cm47XG4gICAgICAgICAgICB0cnkgeyBzdHJlYW0uc2VuZChKU09OLnN0cmluZ2lmeSh7IHZlcnNpb246IDEsIG5hbWUsIG1heEJ5dGVzIH0pKTsgfVxuICAgICAgICAgICAgY2F0Y2ggKGVycm9yKSB7IGZpbmlzaChlcnJvcik7IH1cbiAgICAgICAgfTtcbiAgICAgICAgc3RyZWFtLm9uZXJyb3IgPSAoKSA9PiBmaW5pc2gobmV3IEVycm9yKCdNZWRpYSBzdHJlYW0gZmFpbGVkJykpO1xuICAgICAgICBzdHJlYW0ub25tZXNzYWdlID0gZXZlbnQgPT4ge1xuICAgICAgICAgICAgaWYgKGZpbmlzaGVkKSByZXR1cm47XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGlmICghKGV2ZW50LmRhdGEgaW5zdGFuY2VvZiBBcnJheUJ1ZmZlcikpIHRocm93IG5ldyBFcnJvcignSW52YWxpZCBtZWRpYSBmcmFtZScpO1xuICAgICAgICAgICAgICAgIGlmIChleHBlY3RlZCA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChldmVudC5kYXRhLmJ5dGVMZW5ndGggPiA0MDk2KSB0aHJvdyBuZXcgRXJyb3IoJ0ludmFsaWQgbWVkaWEgbWV0YWRhdGEnKTtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgbWV0YSA9IEpTT04ucGFyc2UobmV3IFRleHREZWNvZGVyKCkuZGVjb2RlKGV2ZW50LmRhdGEpKTtcbiAgICAgICAgICAgICAgICAgICAgaWYgKG1ldGEudmVyc2lvbiAhPT0gMSkgdGhyb3cgbmV3IEVycm9yKCdVbnN1cHBvcnRlZCBtZWRpYSBwcm90b2NvbCcpO1xuICAgICAgICAgICAgICAgICAgICBpZiAobWV0YS5lcnJvcikgdGhyb3cgbWV0YS5saW1pdCA/IG5ldyBSYW5nZUVycm9yKG1ldGEuZXJyb3IpIDogbmV3IEVycm9yKG1ldGEuZXJyb3IpO1xuICAgICAgICAgICAgICAgICAgICBpZiAoIU51bWJlci5pc1NhZmVJbnRlZ2VyKG1ldGEuc2l6ZSkgfHwgbWV0YS5zaXplIDwgMCB8fCB0eXBlb2YgbWV0YS50eXBlICE9PSAnc3RyaW5nJykge1xuICAgICAgICAgICAgICAgICAgICAgICAgdGhyb3cgbmV3IEVycm9yKCdJbnZhbGlkIG1lZGlhIG1ldGFkYXRhJyk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgaWYgKG1ldGEuc2l6ZSA+IG1heEJ5dGVzKSB0aHJvdyBuZXcgUmFuZ2VFcnJvcihgTWVkaWEgZXhjZWVkcyB0aGUgJHttYXhCeXRlc30gYnl0ZSBsaW1pdGApO1xuICAgICAgICAgICAgICAgICAgICBleHBlY3RlZCA9IG1ldGEuc2l6ZTtcbiAgICAgICAgICAgICAgICAgICAgY29udGVudFR5cGUgPSBtZXRhLnR5cGU7XG4gICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgcmVjZWl2ZWQgKz0gZXZlbnQuZGF0YS5ieXRlTGVuZ3RoO1xuICAgICAgICAgICAgICAgICAgICBpZiAoZXZlbnQuZGF0YS5ieXRlTGVuZ3RoID09PSAwIHx8IGV2ZW50LmRhdGEuYnl0ZUxlbmd0aCA+IDY0ICogMTAyNCB8fCByZWNlaXZlZCA+IGV4cGVjdGVkIHx8IHJlY2VpdmVkID4gbWF4Qnl0ZXMpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHRocm93IG5ldyBSYW5nZUVycm9yKCdNZWRpYSBzdHJlYW0gZXhjZWVkZWQgaXRzIGRlY2xhcmVkIGxpbWl0cycpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIHBhcnRzLnB1c2goZXZlbnQuZGF0YSk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSBjYXRjaCAoZXJyb3IpIHsgZmluaXNoKGVycm9yKTsgfVxuICAgICAgICB9O1xuICAgICAgICBzdHJlYW0ub25jbG9zZSA9IGV2ZW50ID0+IHtcbiAgICAgICAgICAgIGlmIChldmVudC5jb2RlICE9PSAxMDAwIHx8IGV4cGVjdGVkID09PSB1bmRlZmluZWQgfHwgcmVjZWl2ZWQgIT09IGV4cGVjdGVkKSB7XG4gICAgICAgICAgICAgICAgZmluaXNoKG5ldyBFcnJvcignTWVkaWEgc3RyZWFtIGNsb3NlZCBiZWZvcmUgdGhlIGNvbXBsZXRlIGZpbGUgYXJyaXZlZCcpKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZmluaXNoKCk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH07XG4gICAgfSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogR29TdHJlYW0gXHUyMDE0IHRoZSBXZWJTb2NrZXQgcHJvZ3JhbW1pbmcgbW9kZWwgd2l0aG91dCBhIGxpc3RlbmluZyBzb2NrZXQuXG4gKlxuICogQSBXZWJTb2NrZXQgY2Fubm90IGJlIHNwb2tlbiBvdmVyIGEgY3VzdG9tIFVSTCBzY2hlbWUsIHNvIGdldHRpbmcgb25lIGluc2lkZVxuICogYSB3ZWJ2aWV3IG1lYW5zIGJpbmRpbmcgYSByZWFsIFRDUCBwb3J0LiBUaGlzIHNwZWFrcyB0aGUgc2FtZSBzaGFwZSBvdmVyIHRoZVxuICogYXNzZXQgc2VydmVyIHRoZSBhcHAgYWxyZWFkeSBzZXJ2ZXM6IEdvXHUyMTkySlMgb3ZlciBvbmUgaGVsZCBwb2xsIHBlciBwYWdlLCBKU1x1MjE5MkdvXG4gKiBvdmVyIGEgbm9ybWFsIFBPU1QuXG4gKlxuICoge0BsaW5rIFN0cmVhbX0gcmV0dXJucyBzeW5jaHJvbm91c2x5IHdpdGggYHJlYWR5U3RhdGUgPT09IENPTk5FQ1RJTkdgLCB0aGUgd2F5XG4gKiBgbmV3IFdlYlNvY2tldCh1cmwpYCBkb2VzLCBzbyBhIGNvbm5lY3Rpb24gY2FuIGJlIGNyZWF0ZWQgYXQgbW9kdWxlIHNjb3BlOlxuICpcbiAqIGBgYGpzXG4gKiBleHBvcnQgY29uc3QgVGVsZW1ldHJ5ID0gU3RyZWFtKFwidGVsZW1ldHJ5XCIpOyAgIC8vIGNvbm5lY3RzIGluIHRoZSBiYWNrZ3JvdW5kXG4gKiBgYGBcbiAqL1xuXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tIFwiLi9uYW5vaWQuanNcIjtcbmltcG9ydCB7IGhhc0RPTSB9IGZyb20gXCIuL2Vudmlyb25tZW50LmpzXCI7XG5cbi8vIEV2ZXJ5IHBpZWNlIG9mIHN0YXRlIHRoYXQgdHdvIHJ1bnRpbWUgaW5zdGFuY2VzIG11c3QgYWdyZWUgb24gbGl2ZXMgaGVyZSwgb25cbi8vIHdpbmRvdywgbm90IGluIG1vZHVsZSBzY29wZS4gVGhlIHJ1bnRpbWUgY2FuIGJlIGV2YWx1YXRlZCBtb3JlIHRoYW4gb25jZSBpbiBhXG4vLyBwYWdlIFx1MjAxNCB0aGUgcGxhdGZvcm0gaW5qZWN0cyBhIGNvcHkgYW5kIGFuIGFwcCBtYXkgaW1wb3J0IGl0IHRvbyBcdTIwMTQgYW5kIGFcbi8vIHBlci1tb2R1bGUgcmVnaXN0cnkgbWVhbnQgYm90aCBpbnN0YW5jZXMgYWxsb2NhdGVkIGNvbm5lY3Rpb24gaWQgMSwgcmFuXG4vLyBjb21wZXRpbmcgcG9sbHMgYWdhaW5zdCBvbmUgc2Vzc2lvbiwgYW5kIGRpc2NhcmRlZCBmcmFtZXMgYWRkcmVzc2VkIHRvIGlkc1xuLy8gb25seSB0aGUgb3RoZXIgaW5zdGFuY2Uga25ldyBhYm91dC4gVGhlIHNlc3Npb24gaWQgYWxvbmUgd2FzIG5vdCBlbm91Z2guXG5pbnRlcmZhY2UgU3RyZWFtR2xvYmFscyB7XG4gICAgc2Vzc2lvbjogc3RyaW5nO1xuICAgIGdlbmVyYXRpb246IG51bWJlcjtcbiAgICBuZXh0Q29ubklEOiBudW1iZXI7XG4gICAgY29ubmVjdGlvbnM6IE1hcDxudW1iZXIsIFdhaWxzU29ja2V0PjtcbiAgICBwb2xsaW5nOiBib29sZWFuO1xuICAgIHBvbGxBYm9ydD86IEFib3J0Q29udHJvbGxlcjtcbn1cblxuY29uc3QgR0VORVJBVElPTl9LRVkgPSBcIl9fd2FpbHNfc3RyZWFtX2dlbmVyYXRpb25cIjtcbmNvbnN0IEdFTkVSQVRJT05fTkFNRV9NQVJLRVIgPSBcIlxcdTAwMWZfX3dhaWxzX3N0cmVhbV9nZW5lcmF0aW9uX186XCI7XG5cbi8vIHdpbmRvdy5uYW1lIGJlbG9uZ3MgdG8gdGhlIGJyb3dzaW5nIGNvbnRleHQgcmF0aGVyIHRoYW4gdGhlIERvY3VtZW50LCBzbyBpdFxuLy8gc3Vydml2ZXMgYSByZWxvYWQgZXZlbiB3aGVuIFdlYiBTdG9yYWdlIGlzIGRpc2FibGVkLiBQcmVzZXJ2ZSBhbnkgbmFtZSB0aGVcbi8vIGFwcGxpY2F0aW9uIGFzc2lnbmVkIGFuZCByZXNlcnZlIG9ubHkgYSBzdWZmaXggZm9yIHRoZSBTdHJlYW0gZ2VuZXJhdGlvbi5cbmZ1bmN0aW9uIG5hbWVkUGFnZUdlbmVyYXRpb24oKTogbnVtYmVyIHtcbiAgICB0cnkge1xuICAgICAgICBjb25zdCBtYXJrZXIgPSB3aW5kb3cubmFtZS5sYXN0SW5kZXhPZihHRU5FUkFUSU9OX05BTUVfTUFSS0VSKTtcbiAgICAgICAgaWYgKG1hcmtlciA8IDApIHJldHVybiAwO1xuICAgICAgICBjb25zdCB2YWx1ZSA9IE51bWJlcih3aW5kb3cubmFtZS5zbGljZShtYXJrZXIgKyBHRU5FUkFUSU9OX05BTUVfTUFSS0VSLmxlbmd0aCkpO1xuICAgICAgICByZXR1cm4gTnVtYmVyLmlzU2FmZUludGVnZXIodmFsdWUpICYmIHZhbHVlID4gMCA/IHZhbHVlIDogMDtcbiAgICB9IGNhdGNoIHtcbiAgICAgICAgcmV0dXJuIDA7XG4gICAgfVxufVxuXG5mdW5jdGlvbiByZW1lbWJlck5hbWVkUGFnZUdlbmVyYXRpb24oZ2VuZXJhdGlvbjogbnVtYmVyKTogdm9pZCB7XG4gICAgdHJ5IHtcbiAgICAgICAgY29uc3QgbWFya2VyID0gd2luZG93Lm5hbWUubGFzdEluZGV4T2YoR0VORVJBVElPTl9OQU1FX01BUktFUik7XG4gICAgICAgIGNvbnN0IGFwcGxpY2F0aW9uTmFtZSA9IG1hcmtlciA8IDAgPyB3aW5kb3cubmFtZSA6IHdpbmRvdy5uYW1lLnNsaWNlKDAsIG1hcmtlcik7XG4gICAgICAgIHdpbmRvdy5uYW1lID0gYXBwbGljYXRpb25OYW1lICsgR0VORVJBVElPTl9OQU1FX01BUktFUiArIGdlbmVyYXRpb247XG4gICAgfSBjYXRjaCB7XG4gICAgICAgIC8vIFNvbWUgZW1iZWRkZWQgYnJvd3NlciBwb2xpY2llcyBjYW4gZGVueSBhY2Nlc3MgdG8gdGhlIGJyb3dzaW5nXG4gICAgICAgIC8vIGNvbnRleHQgbmFtZS4gVGhlIHRpbWUgb3JpZ2luIHJlbWFpbnMgdGhlIGxhc3QtcmVzb3J0IGZhbGxiYWNrLlxuICAgIH1cbn1cblxuLy8gUmVxdWVzdCBhcnJpdmFsIG9yZGVyIGlzIG5vdCBwYWdlIG9yZGVyOiBhIHJlcXVlc3Qgc3RhcnRlZCBqdXN0IGJlZm9yZSBhXG4vLyBuYXZpZ2F0aW9uIGNhbiByZWFjaCBHbyBhZnRlciB0aGUgcmVwbGFjZW1lbnQgcGFnZSdzIGZpcnN0IHJlcXVlc3QuIEtlZXAgYVxuLy8gbW9ub3RvbmljIGNvdW50ZXIgaW4gdGhpcyB0b3AtbGV2ZWwgYnJvd3NpbmcgY29udGV4dCBzbyBHbyBjYW4gZGlzdGluZ3Vpc2hcbi8vIHRob3NlIHBhZ2UgaW5jYXJuYXRpb25zIHdpdGhvdXQgdHJ1c3Rpbmcgc2VydmVyIHNjaGVkdWxpbmcgb3JkZXIuXG5mdW5jdGlvbiBuZXh0UGFnZUdlbmVyYXRpb24oKTogbnVtYmVyIHtcbiAgICBpZiAoIWhhc0RPTSkgcmV0dXJuIDE7XG4gICAgLy8gQW5jaG9yIHRoZSB2YWx1ZSB0byB0aGUgcGFnZSdzIGNyZWF0aW9uIHRpbWUgYXMgd2VsbCBhcyB0aGUgc3RvcmVkXG4gICAgLy8gY291bnRlci4gSWYgYXBwbGljYXRpb24gY29kZSBjbGVhcnMgc2Vzc2lvblN0b3JhZ2UsIHRoZSBuZXh0IHJlbG9hZCBtdXN0XG4gICAgLy8gc3RpbGwgYmUgbmV3ZXIgdGhhbiB0aGUgZ2VuZXJhdGlvbiBHbyBoYXMgYWxyZWFkeSByZXRpcmVkOyByZXN0YXJ0aW5nXG4gICAgLy8gZnJvbSAxIHdvdWxkIG1ha2UgZXZlcnkgc3RyZWFtIHJlcXVlc3QgcmVjZWl2ZSA0MTAgdW50aWwgdGhlIGhvc3QgZXhpdHMuXG4gICAgLy8gUGVyZm9ybWFuY2UudGltZU9yaWdpbiBpcyBhYnNlbnQgb24gb2xkZXIgV2ViS2l0IHJlbGVhc2VzIHRoYXQgV2FpbHNcbiAgICAvLyBzdGlsbCB0YXJnZXRzLiBEYXRlLm5vdyBrZWVwcyB0aG9zZSBlbmdpbmVzIG9uIHRoZSBwcm90b2NvbDsgdGhlIHN0b3JlZFxuICAgIC8vIGNvdW50ZXIgZGlzYW1iaWd1YXRlcyByZWxvYWRzIHRoYXQgaGFwcGVuIHdpdGhpbiB0aGUgc2FtZSBtaWxsaXNlY29uZC5cbiAgICBjb25zdCBvcmlnaW4gPSB0eXBlb2YgcGVyZm9ybWFuY2UgIT09IFwidW5kZWZpbmVkXCIgJiYgTnVtYmVyLmlzRmluaXRlKHBlcmZvcm1hbmNlLnRpbWVPcmlnaW4pXG4gICAgICAgID8gcGVyZm9ybWFuY2UudGltZU9yaWdpblxuICAgICAgICA6IERhdGUubm93KCk7XG4gICAgY29uc3QgcGFnZVRpbWUgPSBNYXRoLm1heCgxLCBNYXRoLmZsb29yKG9yaWdpbiAqIDEwMDApKTtcbiAgICBsZXQgY3VycmVudCA9IG5hbWVkUGFnZUdlbmVyYXRpb24oKTtcbiAgICB0cnkge1xuICAgICAgICBjb25zdCByYXcgPSB3aW5kb3cuc2Vzc2lvblN0b3JhZ2UuZ2V0SXRlbShHRU5FUkFUSU9OX0tFWSk7XG4gICAgICAgIGNvbnN0IHN0b3JlZCA9IHJhdyA9PT0gbnVsbCA/IDAgOiBOdW1iZXIocmF3KTtcbiAgICAgICAgaWYgKE51bWJlci5pc1NhZmVJbnRlZ2VyKHN0b3JlZCkgJiYgc3RvcmVkID4gY3VycmVudCkgY3VycmVudCA9IHN0b3JlZDtcbiAgICB9IGNhdGNoIHtcbiAgICAgICAgLy8gU3RvcmFnZSBjYW4gYmUgZGlzYWJsZWQgYnkgcG9saWN5LiB3aW5kb3cubmFtZSBzdXBwbGllcyB0aGVcbiAgICAgICAgLy8gYnJvd3NpbmctY29udGV4dCBjb3VudGVyIGluIHRoYXQgY2FzZS5cbiAgICB9XG5cbiAgICBjb25zdCBuZXh0ID0gTnVtYmVyLmlzU2FmZUludGVnZXIoY3VycmVudCkgJiYgY3VycmVudCA+PSAwICYmIGN1cnJlbnQgPCBOdW1iZXIuTUFYX1NBRkVfSU5URUdFUlxuICAgICAgICA/IE1hdGgubWF4KGN1cnJlbnQgKyAxLCBwYWdlVGltZSlcbiAgICAgICAgOiBwYWdlVGltZTtcbiAgICB0cnkge1xuICAgICAgICB3aW5kb3cuc2Vzc2lvblN0b3JhZ2Uuc2V0SXRlbShHRU5FUkFUSU9OX0tFWSwgU3RyaW5nKG5leHQpKTtcbiAgICB9IGNhdGNoIHtcbiAgICAgICAgLy8gVGhlIGJyb3dzaW5nLWNvbnRleHQgZmFsbGJhY2sgYmVsb3cgaXMgc3VmZmljaWVudCB3aGVuIHN0b3JhZ2UgaXNcbiAgICAgICAgLy8gdW5hdmFpbGFibGUuXG4gICAgfVxuICAgIHJlbWVtYmVyTmFtZWRQYWdlR2VuZXJhdGlvbihuZXh0KTtcbiAgICByZXR1cm4gbmV4dDtcbn1cblxuZnVuY3Rpb24gc3RyZWFtR2xvYmFscygpOiBTdHJlYW1HbG9iYWxzIHtcbiAgICBjb25zdCBmcmVzaCA9ICgpOiBTdHJlYW1HbG9iYWxzID0+ICh7XG4gICAgICAgIHNlc3Npb246IG5hbm9pZCgpLFxuICAgICAgICBnZW5lcmF0aW9uOiBuZXh0UGFnZUdlbmVyYXRpb24oKSxcbiAgICAgICAgbmV4dENvbm5JRDogMSxcbiAgICAgICAgY29ubmVjdGlvbnM6IG5ldyBNYXA8bnVtYmVyLCBXYWlsc1NvY2tldD4oKSxcbiAgICAgICAgcG9sbGluZzogZmFsc2UsXG4gICAgfSk7XG4gICAgaWYgKCFoYXNET00pIHJldHVybiBmcmVzaCgpO1xuICAgIGNvbnN0IHcgPSAod2luZG93IGFzIGFueSkuX3dhaWxzID0gKHdpbmRvdyBhcyBhbnkpLl93YWlscyB8fCB7fTtcbiAgICByZXR1cm4gKHcuX19zdHJlYW0gPSB3Ll9fc3RyZWFtIHx8IGZyZXNoKCkpO1xufVxuXG5jb25zdCBHID0gc3RyZWFtR2xvYmFscygpO1xuXG5jb25zdCBIRFJfU0VTU0lPTiA9IFwieC13YWlscy1zdHJlYW0tc2Vzc2lvblwiO1xuY29uc3QgSERSX0dFTkVSQVRJT04gPSBcIngtd2FpbHMtc3RyZWFtLWdlbmVyYXRpb25cIjtcbmNvbnN0IEhEUl9DT05OID0gXCJ4LXdhaWxzLXN0cmVhbS1jb25uXCI7XG5jb25zdCBIRFJfS0lORCA9IFwieC13YWlscy1zdHJlYW0ta2luZFwiO1xuY29uc3QgSERSX05BTUUgPSBcIngtd2FpbHMtc3RyZWFtLW5hbWVcIjtcbmNvbnN0IEhEUl9DSFVOSyA9IFwieC13YWlscy1zdHJlYW0tY2h1bmtcIjtcbmNvbnN0IEhEUl9DSFVOS19JTkRFWCA9IFwieC13YWlscy1zdHJlYW0tY2h1bmstaW5kZXhcIjtcbmNvbnN0IEhEUl9DSFVOS19UT1RBTCA9IFwieC13YWlscy1zdHJlYW0tY2h1bmstdG90YWxcIjtcbmNvbnN0IEhEUl9CQVRDSCA9IFwieC13YWlscy1zdHJlYW0tYmF0Y2hcIjtcblxuY29uc3QgS0lORF9EQVRBID0gMDtcbmNvbnN0IEtJTkRfT1BFTiA9IDE7XG5jb25zdCBLSU5EX0NMT1NFID0gMjtcbmNvbnN0IEtJTkRfRVJST1IgPSAzO1xuXG4vLyBNYXRjaGVzIHRoZSBydW50aW1lJ3Mgb3duIGxpbWl0OiBzdGF5IHVuZGVyIFdlYlZpZXcyJ3MgfjJNQiByZXF1ZXN0IGJvZHlcbi8vIGJ1ZmZlcmluZyBsaW1pdCBpbiBXZWJSZXNvdXJjZVJlcXVlc3RlZC5cbmNvbnN0IENIVU5LX1RIUkVTSE9MRCA9IDUxMiAqIDEwMjQ7XG5jb25zdCBNQVhfQkFUQ0hfRlJBTUVTID0gNDA5NjtcblxuLy8gQmFja29mZiBmbG9vciBhbmQgY2VpbGluZyBmb3IgYSBmYWlsaW5nIHBvbGwuIFRoaXMgaXMgdGhlIFwicmVhc29uYWJsZSBsb3dlclxuLy8gbGltaXRcIiBvbiBwb2xsaW5nIFx1MjAxNCBhbiBlcnJvciBiYWNrb2ZmLCBub3QgYSBzdGVhZHktc3RhdGUgaW50ZXJ2YWwuIEluIG5vcm1hbFxuLy8gb3BlcmF0aW9uIHRoZSBwb2xsIHJlLWlzc3VlcyBpbW1lZGlhdGVseSwgYmVjYXVzZSB0aGUgc2VydmVyIGhvbGQgaGFzIGFscmVhZHlcbi8vIHByb3ZpZGVkIHdoYXRldmVyIGRlbGF5IHdhcyBhcHByb3ByaWF0ZS5cbmNvbnN0IEJBQ0tPRkZfTUlOID0gMjUwO1xuY29uc3QgQkFDS09GRl9NQVggPSA1MDAwO1xuY29uc3QgdGV4dEVuY29kZXIgPSBuZXcgVGV4dEVuY29kZXIoKTtcblxuZnVuY3Rpb24gc3RyZWFtVVJMKGVuZHBvaW50OiBzdHJpbmcpOiBzdHJpbmcge1xuICAgIHJldHVybiB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvc3RyZWFtL1wiICsgZW5kcG9pbnQ7XG59XG5cblxuY2xhc3MgU3RyZWFtUmVxdWVzdENhbmNlbGxlZCBleHRlbmRzIEVycm9yIHtcbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgc3VwZXIoXCJzdHJlYW0gcmVxdWVzdCBjYW5jZWxsZWRcIik7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiQWJvcnRFcnJvclwiO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gc2xlZXAobXM6IG51bWJlciwgc2lnbmFsPzogQWJvcnRTaWduYWwpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBpZiAoIXNpZ25hbCkgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlKSA9PiBzZXRUaW1lb3V0KHJlc29sdmUsIG1zKSk7XG4gICAgaWYgKHNpZ25hbC5hYm9ydGVkKSByZXR1cm4gUHJvbWlzZS5yZWplY3QobmV3IFN0cmVhbVJlcXVlc3RDYW5jZWxsZWQoKSk7XG5cbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBjb25zdCBvbkFib3J0ID0gKCkgPT4ge1xuICAgICAgICAgICAgY2xlYXJUaW1lb3V0KHRpbWVyKTtcbiAgICAgICAgICAgIHJlamVjdChuZXcgU3RyZWFtUmVxdWVzdENhbmNlbGxlZCgpKTtcbiAgICAgICAgfTtcbiAgICAgICAgY29uc3QgdGltZXIgPSBzZXRUaW1lb3V0KCgpID0+IHtcbiAgICAgICAgICAgIHNpZ25hbC5yZW1vdmVFdmVudExpc3RlbmVyKFwiYWJvcnRcIiwgb25BYm9ydCk7XG4gICAgICAgICAgICByZXNvbHZlKCk7XG4gICAgICAgIH0sIG1zKTtcbiAgICAgICAgc2lnbmFsLmFkZEV2ZW50TGlzdGVuZXIoXCJhYm9ydFwiLCBvbkFib3J0LCB7IG9uY2U6IHRydWUgfSk7XG4gICAgfSk7XG59XG5cbi8qKlxuICogQSBzdHJlYW0gY29ubmVjdGlvbi4gSW1wbGVtZW50cyB0aGUgdXNlZnVsIHN1YnNldCBvZiB0aGUgYFdlYlNvY2tldGBcbiAqIGludGVyZmFjZSwgc28gdGhlIHNhbWUgYXBwbGljYXRpb24gY29kZSB3b3JrcyBhZ2FpbnN0IGEgcmVhbCBXZWJTb2NrZXQgaW5cbiAqIHNlcnZlciBidWlsZHMuXG4gKlxuICogT25lIGRlbGliZXJhdGUgZGl2ZXJnZW5jZSBmcm9tIHRoZSBzdGFuZGFyZDogYGJpbmFyeVR5cGVgIGRlZmF1bHRzIHRvXG4gKiBgXCJhcnJheWJ1ZmZlclwiYCByYXRoZXIgdGhhbiBgXCJibG9iXCJgLCBiZWNhdXNlIHN0cmVhbSBmcmFtZXMgYXJlIGFsd2F5c1xuICogYmluYXJ5IGFuZCBhIEJsb2Igd291bGQgZm9yY2UgYW4gZXh0cmEgYXN5bmMgaG9wIHRvIHJlYWQgZXZlcnkgbWVzc2FnZS5cbiAqIFNldHRpbmcgaXQgdG8gYFwiYmxvYlwiYCBiZWhhdmVzIGFzIHRoZSBzdGFuZGFyZCBkZXNjcmliZXMuXG4gKi9cbmV4cG9ydCBjbGFzcyBXYWlsc1NvY2tldCBleHRlbmRzIEV2ZW50VGFyZ2V0IHtcbiAgICBzdGF0aWMgcmVhZG9ubHkgQ09OTkVDVElORyA9IDA7XG4gICAgc3RhdGljIHJlYWRvbmx5IE9QRU4gPSAxO1xuICAgIHN0YXRpYyByZWFkb25seSBDTE9TSU5HID0gMjtcbiAgICBzdGF0aWMgcmVhZG9ubHkgQ0xPU0VEID0gMztcblxuICAgIHJlYWRvbmx5IENPTk5FQ1RJTkcgPSAwO1xuICAgIHJlYWRvbmx5IE9QRU4gPSAxO1xuICAgIHJlYWRvbmx5IENMT1NJTkcgPSAyO1xuICAgIHJlYWRvbmx5IENMT1NFRCA9IDM7XG5cbiAgICAvKiogVGhlIHN0cmVhbSBuYW1lIHRoaXMgY29ubmVjdGlvbiB3YXMgb3BlbmVkIGFnYWluc3QuICovXG4gICAgcmVhZG9ubHkgbmFtZTogc3RyaW5nO1xuICAgIHJlYWRvbmx5IHVybDogc3RyaW5nO1xuXG4gICAgLyoqIEFsd2F5cyBlbXB0eTogc3VicHJvdG9jb2xzIGFuZCBleHRlbnNpb25zIGFyZSBub3QgbmVnb3RpYXRlZC4gKi9cbiAgICByZWFkb25seSBwcm90b2NvbCA9IFwiXCI7XG4gICAgcmVhZG9ubHkgZXh0ZW5zaW9ucyA9IFwiXCI7XG5cbiAgICBiaW5hcnlUeXBlOiBcImFycmF5YnVmZmVyXCIgfCBcImJsb2JcIiA9IFwiYXJyYXlidWZmZXJcIjtcbiAgICByZWFkeVN0YXRlOiBudW1iZXIgPSBXYWlsc1NvY2tldC5DT05ORUNUSU5HO1xuXG4gICAgLy8gRGVjbGFyZWQgZm9yIHRoZSB0eXBlOyB0aGUgY29uc3RydWN0b3IgcmVwbGFjZXMgZWFjaCB3aXRoIGFuIGFjY2Vzc29yXG4gICAgLy8gdGhhdCByZWdpc3RlcnMgYSByZWFsIGxpc3RlbmVyLiBDYWxsaW5nIHRoZW0gZGlyZWN0bHkgaW5zdGVhZCB3b3VsZCBtYWtlXG4gICAgLy8gYW55IHdyYXBwZXIgKHNlZSBKU09OU3RyZWFtKSByZWNlaXZlIGJvdGggdGhlIHJhdyBhbmQgdGhlIGRlY29kZWQgZXZlbnQuXG4gICAgb25vcGVuOiAoKGV2OiBFdmVudCkgPT4gdm9pZCkgfCBudWxsID0gbnVsbDtcbiAgICBvbm1lc3NhZ2U6ICgoZXY6IE1lc3NhZ2VFdmVudCkgPT4gdm9pZCkgfCBudWxsID0gbnVsbDtcbiAgICBvbmNsb3NlOiAoKGV2OiBDbG9zZUV2ZW50KSA9PiB2b2lkKSB8IG51bGwgPSBudWxsO1xuICAgIG9uZXJyb3I6ICgoZXY6IEV2ZW50KSA9PiB2b2lkKSB8IG51bGwgPSBudWxsO1xuXG4gICAgLyoqXG4gICAgICogQGludGVybmFsIEFwcGxpZWQgdG8gZXZlcnkgaW5ib3VuZCBwYXlsb2FkIGJlZm9yZSBkaXNwYXRjaC4gSlNPTlN0cmVhbVxuICAgICAqIHJlcGxhY2VzIGl0LCBzbyBgb25tZXNzYWdlYCBhbmQgYGFkZEV2ZW50TGlzdGVuZXIoXCJtZXNzYWdlXCIsIFx1MjAyNilgIHNlZSB0aGVcbiAgICAgKiBzYW1lIGRlY29kZWQgdmFsdWUgXHUyMDE0IGludGVyY2VwdGluZyBvbmx5IGBvbm1lc3NhZ2VgIG1hZGUgdGhlIHR3byBkaXNhZ3JlZS5cbiAgICAgKi9cbiAgICBfZGVjb2RlOiAocGF5bG9hZDogQXJyYXlCdWZmZXIpID0+IHVua25vd24gPSAocGF5bG9hZCkgPT4gcGF5bG9hZDtcblxuICAgIC8qKiBAaW50ZXJuYWwgKi8gcmVhZG9ubHkgX2lkOiBudW1iZXI7XG4gICAgcHJpdmF0ZSBfYnVmZmVyZWQgPSAwO1xuICAgIC8vIFNlbmRzIGFyZSBzZXJpYWxpc2VkIHBlciBjb25uZWN0aW9uOiBjb25jdXJyZW50IGZldGNoKCkgUE9TVHMgZG8gbm90XG4gICAgLy8gcHJlc2VydmUgb3JkZXIsIGFuZCBHbyByZWxpZXMgb24gc2VuZCBvcmRlciBiZWluZyB0aGUgb3JkZXIgaXQgb2JzZXJ2ZXMuXG4gICAgcHJpdmF0ZSBfY2hhaW46IFByb21pc2U8dm9pZD4gPSBQcm9taXNlLnJlc29sdmUoKTtcbiAgICBwcml2YXRlIF9wZW5kaW5nOiBQcm9taXNlPFVpbnQ4QXJyYXk+W10gPSBbXTtcbiAgICBwcml2YXRlIF9mbHVzaGluZyA9IGZhbHNlO1xuICAgIHByaXZhdGUgX29wZW5BYm9ydCA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICBwcml2YXRlIF9zZW5kQWJvcnQgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG5cbiAgICAvKiogQnl0ZXMgcXVldWVkIGJ5IHNlbmQoKSB0aGF0IGhhdmUgbm90IHlldCByZWFjaGVkIEdvLiAqL1xuICAgIGdldCBidWZmZXJlZEFtb3VudCgpOiBudW1iZXIge1xuICAgICAgICByZXR1cm4gdGhpcy5fYnVmZmVyZWQ7XG4gICAgfVxuXG4gICAgY29uc3RydWN0b3IobmFtZTogc3RyaW5nKSB7XG4gICAgICAgIHN1cGVyKCk7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuX2lkID0gRy5uZXh0Q29ubklEKys7XG4gICAgICAgIHRoaXMudXJsID0gaGFzRE9NID8gc3RyZWFtVVJMKFwicG9sbFwiKSArIFwiI1wiICsgZW5jb2RlVVJJQ29tcG9uZW50KG5hbWUpIDogXCJcIjtcblxuICAgICAgICBmb3IgKGNvbnN0IHR5cGUgb2YgW1wib3BlblwiLCBcIm1lc3NhZ2VcIiwgXCJjbG9zZVwiLCBcImVycm9yXCJdKSB7XG4gICAgICAgICAgICBkZWZpbmVIYW5kbGVyUHJvcGVydHkodGhpcywgdHlwZSk7XG4gICAgICAgIH1cblxuICAgICAgICBHLmNvbm5lY3Rpb25zLnNldCh0aGlzLl9pZCwgdGhpcyk7XG5cbiAgICAgICAgLy8gRmlyZSBhbmQgZm9yZ2V0OiB0aGUgYWNrIGFycml2ZXMgYXMgYW4gb3BlbiBmcmFtZSBvbiB0aGUgcG9sbCwgd2hpY2hcbiAgICAgICAgLy8gaXMgd2hhdCBtb3ZlcyByZWFkeVN0YXRlIHRvIE9QRU4uXG4gICAgICAgIHRoaXMuX2NoYWluID0gdGhpcy5fY2hhaW4udGhlbigoKSA9PlxuICAgICAgICAgICAgcG9zdEZyYW1lKHRoaXMuX2lkLCBLSU5EX09QRU4sIG5ldyBVaW50OEFycmF5KDApLCBuYW1lLCB0aGlzLl9vcGVuQWJvcnQuc2lnbmFsKVxuICAgICAgICApLmNhdGNoKChlcnIpID0+IHtcbiAgICAgICAgICAgIHRoaXMuX2ZhaWwoZXJyKTtcbiAgICAgICAgfSk7XG5cbiAgICAgICAgc3RhcnRQb2xsaW5nKCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2VuZCBhIGZyYW1lIHRvIEdvLiBBY2NlcHRzIGFueXRoaW5nIGEgV2ViU29ja2V0IGRvZXM7IHN0cmluZ3MgYXJlXG4gICAgICogZW5jb2RlZCBhcyBVVEYtOCwgc2luY2UgR28gcmVjZWl2ZXMgZXZlcnkgZnJhbWUgYXMgYSBieXRlIHNsaWNlLlxuICAgICAqL1xuICAgIHNlbmQoZGF0YTogc3RyaW5nIHwgQXJyYXlCdWZmZXJMaWtlIHwgQXJyYXlCdWZmZXJWaWV3IHwgQmxvYik6IHZvaWQge1xuICAgICAgICBpZiAodGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DT05ORUNUSU5HKSB7XG4gICAgICAgICAgICAvLyBNYXRjaGVzIHRoZSBXZWJTb2NrZXQgc3BlY2lmaWNhdGlvbi5cbiAgICAgICAgICAgIHRocm93IG5ldyBET01FeGNlcHRpb24oXCJTdGlsbCBpbiBDT05ORUNUSU5HIHN0YXRlLlwiLCBcIkludmFsaWRTdGF0ZUVycm9yXCIpO1xuICAgICAgICB9XG4gICAgICAgIGlmICh0aGlzLnJlYWR5U3RhdGUgIT09IFdhaWxzU29ja2V0Lk9QRU4pIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIFJlc29sdmUgbm93IHJhdGhlciB0aGFuIGluc2lkZSB0aGUgY2hhaW4uIE11dGFibGUgYmluYXJ5IGlucHV0cyBhcmVcbiAgICAgICAgLy8gY29waWVkIHN5bmNocm9ub3VzbHkgc28gc2VuZCgpIHNuYXBzaG90cyB0aGVpciBieXRlcyBsaWtlIGEgbmF0aXZlXG4gICAgICAgIC8vIFdlYlNvY2tldCwgZXZlbiB3aGVuIGEgYmF0Y2ggd2FpdHMgYmVoaW5kIGFuIGluLWZsaWdodCByZXF1ZXN0LiBBXG4gICAgICAgIC8vIEJsb2IgaXMgaW1tdXRhYmxlIGFuZCBpcyB0aGVyZWZvcmUgc2FmZSB0byByZWFkIGFzeW5jaHJvbm91c2x5IG9uY2UuXG4gICAgICAgIGNvbnN0IGltbWVkaWF0ZSA9IHRvQnl0ZXNTeW5jKGRhdGEpO1xuICAgICAgICBjb25zdCBzbmFwc2hvdCA9IGltbWVkaWF0ZVxuICAgICAgICAgICAgPyBQcm9taXNlLnJlc29sdmUoaW1tZWRpYXRlKVxuICAgICAgICAgICAgOiB0b0J5dGVzKGRhdGEsIHRoaXMuX3NlbmRBYm9ydC5zaWduYWwpO1xuICAgICAgICAvLyBjbG9zZSgpIG1heSBkaXNjYXJkIHRoaXMgcHJvbWlzZSBiZWZvcmUgdGhlIGZsdXNoIGhhcyBkZXF1ZXVlZCBpdC5cbiAgICAgICAgLy8gQXR0YWNoIGEgcmVqZWN0aW9uIG9ic2VydmVyIGltbWVkaWF0ZWx5IHNvIGNhbmNlbGxpbmcgYW4gYXN5bmNocm9ub3VzXG4gICAgICAgIC8vIEJsb2IgcmVhZCBjYW5ub3QgYmVjb21lIGFuIHVuaGFuZGxlZCByZWplY3Rpb247IFByb21pc2UuYWxsIHN0aWxsXG4gICAgICAgIC8vIG9ic2VydmVzIHRoZSBvcmlnaW5hbCByZWplY3Rpb24gd2hlbiB0aGUgZmx1c2ggYWxyZWFkeSBvd25zIGl0LlxuICAgICAgICB2b2lkIHNuYXBzaG90LmNhdGNoKCgpID0+IHt9KTtcblxuICAgICAgICAvLyBBIEJsb2IgY2Fubm90IGJlIGNvbnZlcnRlZCBzeW5jaHJvbm91c2x5LCBidXQgQmxvYi5zaXplIGlzIGF2YWlsYWJsZVxuICAgICAgICAvLyBub3csIHNvIGl0cyBieXRlcyBhcmUgc3RpbGwgY291bnRlZCB0aGUgbW9tZW50IHNlbmQoKSByZXR1cm5zLiBXYWl0aW5nXG4gICAgICAgIC8vIGZvciBhcnJheUJ1ZmZlcigpIHRvIHJlc29sdmUgbGV0IGEgbGFyZ2UgQmxvYiBzbGlwIHBhc3QgY29kZSB1c2luZ1xuICAgICAgICAvLyBidWZmZXJlZEFtb3VudCBmb3IgYmFja3ByZXNzdXJlLlxuICAgICAgICBjb25zdCBibG9iU2l6ZSA9XG4gICAgICAgICAgICAhaW1tZWRpYXRlICYmIHR5cGVvZiBCbG9iICE9PSBcInVuZGVmaW5lZFwiICYmIGRhdGEgaW5zdGFuY2VvZiBCbG9iID8gZGF0YS5zaXplIDogMDtcbiAgICAgICAgaWYgKGJsb2JTaXplKSB7XG4gICAgICAgICAgICB0aGlzLl9idWZmZXJlZCArPSBibG9iU2l6ZTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIEFjY291bnQgZm9yIHRoZSBieXRlcyBub3csIG5vdCB3aGVuIHRoZSBjaGFpbiByZWFjaGVzIHRoZW0uIFJlcG9ydGluZ1xuICAgICAgICAvLyB6ZXJvIHdoaWxlIGZyYW1lcyB3YWl0IGJlaGluZCBhbiBpbi1mbGlnaHQgYmF0Y2ggd291bGQgdGVsbCBhblxuICAgICAgICAvLyBhcHBsaWNhdGlvbiB1c2luZyBidWZmZXJlZEFtb3VudCBmb3IgYmFja3ByZXNzdXJlIHRoYXQgaXQgbWF5IGtlZXBcbiAgICAgICAgLy8gc2VuZGluZy5cbiAgICAgICAgaWYgKGltbWVkaWF0ZSkge1xuICAgICAgICAgICAgdGhpcy5fYnVmZmVyZWQgKz0gaW1tZWRpYXRlLmJ5dGVMZW5ndGg7XG4gICAgICAgIH1cbiAgICAgICAgdGhpcy5fcGVuZGluZy5wdXNoKHNuYXBzaG90KTtcblxuICAgICAgICAvLyBRdWV1ZSByYXRoZXIgdGhhbiBwb3N0LiBXaGF0ZXZlciBhY2N1bXVsYXRlcyB3aGlsZSBhIHJlcXVlc3QgaXMgaW5cbiAgICAgICAgLy8gZmxpZ2h0IGdvZXMgb3V0IHRvZ2V0aGVyIGluIHRoZSBuZXh0IG9uZSwgc28gdGhlIHBlci1yZXF1ZXN0IGNvc3QgXHUyMDE0XG4gICAgICAgIC8vIGEgc2NoZW1lLWhhbmRsZXIgcm91bmQgdHJpcCwgYWJvdXQgZWxldmVuIGNnbyBjYWxscyBvbiBtYWNPUyBcdTIwMTQgaXNcbiAgICAgICAgLy8gZGl2aWRlZCBhY3Jvc3MgdGhlIGJhdGNoLiBVbmRlciBsaWdodCBsb2FkIGEgYmF0Y2ggaXMgb25lIGZyYW1lIGFuZFxuICAgICAgICAvLyB0aGlzIGJlaGF2ZXMgZXhhY3RseSBhcyBiZWZvcmU7IHRoZSBiYXRjaGluZyBvbmx5IGFwcGVhcnMgd2hlbiB0aGVcbiAgICAgICAgLy8gc2VuZGVyIGlzIG91dHJ1bm5pbmcgdGhlIHRyYW5zcG9ydCwgd2hpY2ggaXMgd2hlbiBpdCBpcyBuZWVkZWQuXG4gICAgICAgIGlmICh0aGlzLl9mbHVzaGluZykgcmV0dXJuO1xuICAgICAgICB0aGlzLl9mbHVzaGluZyA9IHRydWU7XG5cbiAgICAgICAgdGhpcy5fY2hhaW4gPSB0aGlzLl9jaGFpbi50aGVuKGFzeW5jICgpID0+IHtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgd2hpbGUgKHRoaXMuX3BlbmRpbmcubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgICAgICAgICAvLyBCb3VuZCB0aGUgd29yayBhbmQgdGVtcG9yYXJ5IHByb21pc2UgYXJyYXkgZm9yIG9uZSBwYXNzLlxuICAgICAgICAgICAgICAgICAgICAvLyBwb3N0QmF0Y2ggYXBwbGllcyB0aGUgd2lyZS1zaXplIGJvdW5kIGFmdGVyIHJlc29sdXRpb24uXG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHF1ZXVlZCA9IHRoaXMuX3BlbmRpbmcuc3BsaWNlKDAsIE1BWF9CQVRDSF9GUkFNRVMpO1xuICAgICAgICAgICAgICAgICAgICBjb25zdCBiYXRjaCA9IGF3YWl0IFByb21pc2UuYWxsKHF1ZXVlZCk7XG5cbiAgICAgICAgICAgICAgICAgICAgLy8gQSByZW1vdGUgY2xvc2UgbWF5IGFycml2ZSB3aGlsZSBhIEJsb2IgaXMgYmVpbmcgcmVhZC5cbiAgICAgICAgICAgICAgICAgICAgLy8gVGhlIGNvbm5lY3Rpb24gaXMgYWxyZWFkeSBnb25lIGluIHRoYXQgY2FzZSwgc28gZG8gbm90XG4gICAgICAgICAgICAgICAgICAgIC8vIGlzc3VlIGEgbGF0ZSBQT1NUIGZvciBieXRlcyB0aGUgcGVlciBjYW4gbm8gbG9uZ2VyIHRha2UuXG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzLnJlYWR5U3RhdGUgPT09IFdhaWxzU29ja2V0LkNMT1NFRCkge1xuICAgICAgICAgICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgICAgICAgICBjb25zdCB0b3RhbCA9IGJhdGNoLnJlZHVjZSgobiwgYikgPT4gbiArIGIuYnl0ZUxlbmd0aCwgMCk7XG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBhd2FpdCBwb3N0QmF0Y2godGhpcy5faWQsIGJhdGNoLCB0aGlzLl9zZW5kQWJvcnQuc2lnbmFsKTtcbiAgICAgICAgICAgICAgICAgICAgfSBmaW5hbGx5IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIC8vIF9jbG9zZWQgcmVzZXRzIHRoZSBjb3VudGVyIGltbWVkaWF0ZWx5LiBBIHJlcXVlc3RcbiAgICAgICAgICAgICAgICAgICAgICAgIC8vIGFscmVhZHkgaW4gZmxpZ2h0IG1heSBmaW5pc2ggYWZ0ZXJ3YXJkcywgYW5kIG11c3RcbiAgICAgICAgICAgICAgICAgICAgICAgIC8vIG5vdCBkcml2ZSBidWZmZXJlZEFtb3VudCBiZWxvdyB6ZXJvLlxuICAgICAgICAgICAgICAgICAgICAgICAgdGhpcy5fYnVmZmVyZWQgPSBNYXRoLm1heCgwLCB0aGlzLl9idWZmZXJlZCAtIHRvdGFsKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0gZmluYWxseSB7XG4gICAgICAgICAgICAgICAgdGhpcy5fZmx1c2hpbmcgPSBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSkuY2F0Y2goKGVycikgPT4ge1xuICAgICAgICAgICAgdGhpcy5fZmx1c2hpbmcgPSBmYWxzZTtcbiAgICAgICAgICAgIHRoaXMuX2ZhaWwoZXJyKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqIENsb3NlIHRoZSBjb25uZWN0aW9uLiBHbydzIGhhbmRsZXIgc2VlcyBSZWNlaXZlIGZhaWwgYW5kIHJldHVybnMuICovXG4gICAgY2xvc2UoY29kZTogbnVtYmVyID0gMTAwMCwgcmVhc29uOiBzdHJpbmcgPSBcIlwiKTogdm9pZCB7XG4gICAgICAgIGlmICh0aGlzLnJlYWR5U3RhdGUgPT09IFdhaWxzU29ja2V0LkNMT1NJTkcgfHwgdGhpcy5yZWFkeVN0YXRlID09PSBXYWlsc1NvY2tldC5DTE9TRUQpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICB0aGlzLnJlYWR5U3RhdGUgPSBXYWlsc1NvY2tldC5DTE9TSU5HO1xuICAgICAgICB0aGlzLl9wZW5kaW5nLmxlbmd0aCA9IDA7XG4gICAgICAgIHRoaXMuX2J1ZmZlcmVkID0gMDtcbiAgICAgICAgdGhpcy5fb3BlbkFib3J0LmFib3J0KCk7XG5cbiAgICAgICAgLy8gQSBkYXRhIGZyYW1lIHdhaXRpbmcgZm9yIHJlY2VpdmVyIGNhcGFjaXR5IG11c3Qgbm90IHByZXZlbnQgdGhlXG4gICAgICAgIC8vIHJlc2VydmVkIGNsb3NlIGNvbnRyb2wgZnJvbSBldmVyIHJlYWNoaW5nIEdvLiBDYW5jZWwgb25seSB0aGUgZGF0YVxuICAgICAgICAvLyByZXRyeSBwYXRoOyB0aGUgY2hhaW4gYmVsb3cgc3RpbGwgb3JkZXJzIHRoZSBjbG9zZSBhZnRlciB0aGF0IHBhdGhcbiAgICAgICAgLy8gaGFzIHVud291bmQsIHByZXNlcnZpbmcgbm9ybWFsIHNlbmQtYmVmb3JlLWNsb3NlIG9yZGVyaW5nLlxuICAgICAgICB0aGlzLl9zZW5kQWJvcnQuYWJvcnQoKTtcbiAgICAgICAgdGhpcy5fY2hhaW4gPSB0aGlzLl9jaGFpbi50aGVuKCgpID0+XG4gICAgICAgICAgICBwb3N0RnJhbWUodGhpcy5faWQsIEtJTkRfQ0xPU0UsIG5ldyBVaW50OEFycmF5KDApKVxuICAgICAgICApLmNhdGNoKCgpID0+IHtcbiAgICAgICAgICAgIC8vIFRoZSBwZWVyIGlzIGFscmVhZHkgZ29uZTsgdGhlIGxvY2FsIGNsb3NlIGJlbG93IGlzIHdoYXQgbWF0dGVycy5cbiAgICAgICAgfSkudGhlbigoKSA9PiB7XG4gICAgICAgICAgICB0aGlzLl9jbG9zZWQoY29kZSwgcmVhc29uLCB0cnVlKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqIEBpbnRlcm5hbCBDYWxsZWQgd2hlbiB0aGUgb3BlbiBmcmFtZSBjb21lcyBiYWNrIGZyb20gR28uICovXG4gICAgX29wZW5lZCgpOiB2b2lkIHtcbiAgICAgICAgaWYgKHRoaXMucmVhZHlTdGF0ZSAhPT0gV2FpbHNTb2NrZXQuQ09OTkVDVElORykgcmV0dXJuO1xuICAgICAgICB0aGlzLnJlYWR5U3RhdGUgPSBXYWlsc1NvY2tldC5PUEVOO1xuICAgICAgICB0aGlzLmRpc3BhdGNoRXZlbnQobmV3IEV2ZW50KFwib3BlblwiKSk7XG4gICAgfVxuXG4gICAgLyoqIEBpbnRlcm5hbCBDYWxsZWQgZm9yIGVhY2ggZGF0YSBmcmFtZSBhZGRyZXNzZWQgdG8gdGhpcyBjb25uZWN0aW9uLiAqL1xuICAgIF9tZXNzYWdlKHBheWxvYWQ6IEFycmF5QnVmZmVyKTogdm9pZCB7XG4gICAgICAgIGlmICh0aGlzLnJlYWR5U3RhdGUgIT09IFdhaWxzU29ja2V0Lk9QRU4pIHJldHVybjtcbiAgICAgICAgbGV0IGRhdGE6IHVua25vd247XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBkYXRhID0gdGhpcy5fZGVjb2RlKHBheWxvYWQpO1xuICAgICAgICB9IGNhdGNoIHtcbiAgICAgICAgICAgIHRoaXMuZGlzcGF0Y2hFdmVudChuZXcgRXZlbnQoXCJlcnJvclwiKSk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgaWYgKGRhdGEgPT09IHBheWxvYWQgJiYgdGhpcy5iaW5hcnlUeXBlID09PSBcImJsb2JcIikge1xuICAgICAgICAgICAgZGF0YSA9IG5ldyBCbG9iKFtwYXlsb2FkXSk7XG4gICAgICAgIH1cbiAgICAgICAgdGhpcy5kaXNwYXRjaEV2ZW50KG5ldyBNZXNzYWdlRXZlbnQoXCJtZXNzYWdlXCIsIHsgZGF0YSB9KSk7XG4gICAgfVxuXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIF9mYWlsKGVycjogdW5rbm93bik6IHZvaWQge1xuICAgICAgICAvLyBBbiBpbi1mbGlnaHQgc2VuZCBjYW4gcmVqZWN0IGFmdGVyIGEgY2xlYW4gY2xvc2UgaGFzIGFscmVhZHkgYXJyaXZlZC5cbiAgICAgICAgLy8gVGhlIHNvY2tldCdzIGxpZmVjeWNsZSBpcyBjb21wbGV0ZSB0aGVuOiBlbWl0dGluZyBlcnJvciBhZnRlciBjbG9zZVxuICAgICAgICAvLyByZXZlcnNlcyB0aGUgV2ViU29ja2V0IGV2ZW50IG9yZGVyIGFuZCByZXBvcnRzIGEgZmFsc2UgZmFpbHVyZS5cbiAgICAgICAgaWYgKHRoaXMucmVhZHlTdGF0ZSA9PT0gV2FpbHNTb2NrZXQuQ0xPU0VEIHx8IHRoaXMucmVhZHlTdGF0ZSA9PT0gV2FpbHNTb2NrZXQuQ0xPU0lORykgcmV0dXJuO1xuICAgICAgICB0aGlzLmRpc3BhdGNoRXZlbnQobmV3IEV2ZW50KFwiZXJyb3JcIikpO1xuICAgICAgICAvLyAxMDA2OiBjbG9zZWQgYWJub3JtYWxseSwgbm8gY2xvc2UgZnJhbWUuIFNhbWUgY29kZSBhIHJlYWwgV2ViU29ja2V0XG4gICAgICAgIC8vIHJlcG9ydHMgd2hlbiB0aGUgY29ubmVjdGlvbiBkcm9wcyB3aXRob3V0IGEgaGFuZHNoYWtlLlxuICAgICAgICB0aGlzLl9jbG9zZWQoMTAwNiwgZXJyIGluc3RhbmNlb2YgRXJyb3IgPyBlcnIubWVzc2FnZSA6IFN0cmluZyhlcnIpLCBmYWxzZSk7XG4gICAgfVxuXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIF9jbG9zZWQoY29kZTogbnVtYmVyLCByZWFzb246IHN0cmluZywgd2FzQ2xlYW46IGJvb2xlYW4pOiB2b2lkIHtcbiAgICAgICAgaWYgKHRoaXMucmVhZHlTdGF0ZSA9PT0gV2FpbHNTb2NrZXQuQ0xPU0VEKSByZXR1cm47XG4gICAgICAgIHRoaXMucmVhZHlTdGF0ZSA9IFdhaWxzU29ja2V0LkNMT1NFRDtcbiAgICAgICAgdGhpcy5fb3BlbkFib3J0LmFib3J0KCk7XG4gICAgICAgIHRoaXMuX3NlbmRBYm9ydC5hYm9ydCgpO1xuICAgICAgICBHLmNvbm5lY3Rpb25zLmRlbGV0ZSh0aGlzLl9pZCk7XG4gICAgICAgIGlmIChHLmNvbm5lY3Rpb25zLnNpemUgPT09IDApIEcucG9sbEFib3J0Py5hYm9ydCgpO1xuICAgICAgICB0aGlzLl9wZW5kaW5nLmxlbmd0aCA9IDA7XG4gICAgICAgIHRoaXMuX2J1ZmZlcmVkID0gMDtcblxuICAgICAgICBjb25zdCBldiA9IHR5cGVvZiBDbG9zZUV2ZW50ID09PSBcImZ1bmN0aW9uXCJcbiAgICAgICAgICAgID8gbmV3IENsb3NlRXZlbnQoXCJjbG9zZVwiLCB7IGNvZGUsIHJlYXNvbiwgd2FzQ2xlYW4gfSlcbiAgICAgICAgICAgIDogT2JqZWN0LmFzc2lnbihuZXcgRXZlbnQoXCJjbG9zZVwiKSwgeyBjb2RlLCByZWFzb24sIHdhc0NsZWFuIH0pIGFzIENsb3NlRXZlbnQ7XG4gICAgICAgIHRoaXMuZGlzcGF0Y2hFdmVudChldik7XG4gICAgfVxufVxuXG4vKipcbiAqIENvbm5lY3QgdG8gYSBzdHJlYW0gZGVjbGFyZWQgaW4gR28gd2l0aCBgYXBwLkhhbmRsZVN0cmVhbShuYW1lLCBoYW5kbGVyKWAuXG4gKlxuICogUmV0dXJucyBpbW1lZGlhdGVseSB3aXRoIGByZWFkeVN0YXRlID09PSBDT05ORUNUSU5HYDsgbGlzdGVuIGZvciBgb3BlbmAsIG9yXG4gKiBqdXN0IGNhbGwge0BsaW5rIFdhaWxzU29ja2V0LnNlbmR9IG9uY2Ugb3Blbi4gQ3JlYXRpbmcgYSBjb25uZWN0aW9uIGF0IG1vZHVsZVxuICogc2NvcGUgaXMgc3VwcG9ydGVkIGFuZCBpcyB0aGUgaW50ZW5kZWQgc2hhcGUgZm9yIGdlbmVyYXRlZCBiaW5kaW5ncy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFN0cmVhbShuYW1lOiBzdHJpbmcpOiBXYWlsc1NvY2tldCB8IFdlYlNvY2tldCB7XG4gICAgLy8gU2VydmVyIGJ1aWxkcyBpbnN0YWxsIGEgZmFjdG9yeSByZXR1cm5pbmcgYSByZWFsIFdlYlNvY2tldCwgc2luY2UgdGhlcmVcbiAgICAvLyBpcyBhIGxpc3RlbmVyIHRoZXJlIHRvIHVwZ3JhZGUgb24gYW5kIG5vdGhpbmcgdG8gZW11bGF0ZS4gQm90aCBvYmplY3RzXG4gICAgLy8gcHJlc2VudCB0aGUgc2FtZSBpbnRlcmZhY2UsIHNvIG5vdGhpbmcgYWJvdmUgdGhpcyBsaW5lIGNhcmVzIHdoaWNoIGl0IGlzLlxuICAgIGlmIChoYXNET00pIHtcbiAgICAgICAgY29uc3QgZmFjdG9yeSA9ICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LnN0cmVhbUZhY3Rvcnk7XG4gICAgICAgIGlmICh0eXBlb2YgZmFjdG9yeSA9PT0gXCJmdW5jdGlvblwiKSB7XG4gICAgICAgICAgICByZXR1cm4gZmFjdG9yeShuYW1lKTtcbiAgICAgICAgfVxuICAgIH1cbiAgICByZXR1cm4gbmV3IFdhaWxzU29ja2V0KG5hbWUpO1xufVxuXG4vLyBEZWZpbmVzIGFuIGBvbjx0eXBlPmAgcHJvcGVydHkgdGhhdCByZWdpc3RlcnMgYW5kIHVucmVnaXN0ZXJzIGEgcmVhbCBldmVudFxuLy8gbGlzdGVuZXIsIHRoZSB3YXkgdGhlIERPTSBkZWZpbmVzIHRoZW0uIEFzc2lnbmluZyByZXBsYWNlcyB0aGUgcHJldmlvdXNcbi8vIGhhbmRsZXI7IHJlYWRpbmcgcmV0dXJucyBpdC5cbmZ1bmN0aW9uIGRlZmluZUhhbmRsZXJQcm9wZXJ0eSh0YXJnZXQ6IEV2ZW50VGFyZ2V0LCB0eXBlOiBzdHJpbmcpOiB2b2lkIHtcbiAgICBsZXQgY3VycmVudDogKChldjogYW55KSA9PiB2b2lkKSB8IG51bGwgPSBudWxsO1xuICAgIE9iamVjdC5kZWZpbmVQcm9wZXJ0eSh0YXJnZXQsIFwib25cIiArIHR5cGUsIHtcbiAgICAgICAgZ2V0OiAoKSA9PiBjdXJyZW50LFxuICAgICAgICBzZXQoZm46ICgoZXY6IGFueSkgPT4gdm9pZCkgfCBudWxsKSB7XG4gICAgICAgICAgICBpZiAoY3VycmVudCkgdGFyZ2V0LnJlbW92ZUV2ZW50TGlzdGVuZXIodHlwZSwgY3VycmVudCk7XG4gICAgICAgICAgICBjdXJyZW50ID0gdHlwZW9mIGZuID09PSBcImZ1bmN0aW9uXCIgPyBmbiA6IG51bGw7XG4gICAgICAgICAgICBpZiAoY3VycmVudCkgdGFyZ2V0LmFkZEV2ZW50TGlzdGVuZXIodHlwZSwgY3VycmVudCk7XG4gICAgICAgIH0sXG4gICAgICAgIGNvbmZpZ3VyYWJsZTogdHJ1ZSxcbiAgICAgICAgZW51bWVyYWJsZTogdHJ1ZSxcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBUaGUgb2JqZWN0LXNoYXBlZCB2aWV3IG9mIGEgc3RyZWFtIHJldHVybmVkIGJ5IHtAbGluayBKU09OU3RyZWFtfTogYHNlbmRgXG4gKiB0YWtlcyBhIHZhbHVlIHRvIHN0cmluZ2lmeSwgYW5kIGBldi5kYXRhYCBvbiBhIG1lc3NhZ2UgaXMgdGhlIHBhcnNlZCByZXN1bHQuXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgSlNPTlNvY2tldCBleHRlbmRzIE9taXQ8V2FpbHNTb2NrZXQsIFwic2VuZFwiIHwgXCJvbm1lc3NhZ2VcIiB8IFwiYWRkRXZlbnRMaXN0ZW5lclwiIHwgXCJyZW1vdmVFdmVudExpc3RlbmVyXCI+IHtcbiAgICBzZW5kKHZhbHVlOiB1bmtub3duKTogdm9pZDtcbiAgICBvbm1lc3NhZ2U6ICgoZXY6IE1lc3NhZ2VFdmVudDxhbnk+KSA9PiB2b2lkKSB8IG51bGw7XG4gICAgYWRkRXZlbnRMaXN0ZW5lcihcbiAgICAgICAgdHlwZTogXCJtZXNzYWdlXCIsXG4gICAgICAgIGxpc3RlbmVyOiAoKHRoaXM6IEpTT05Tb2NrZXQsIGV2OiBNZXNzYWdlRXZlbnQ8YW55PikgPT4gYW55KSB8IEV2ZW50TGlzdGVuZXJPYmplY3QgfCBudWxsLFxuICAgICAgICBvcHRpb25zPzogYm9vbGVhbiB8IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zLFxuICAgICk6IHZvaWQ7XG4gICAgYWRkRXZlbnRMaXN0ZW5lcihcbiAgICAgICAgdHlwZTogc3RyaW5nLFxuICAgICAgICBsaXN0ZW5lcjogRXZlbnRMaXN0ZW5lck9yRXZlbnRMaXN0ZW5lck9iamVjdCB8IG51bGwsXG4gICAgICAgIG9wdGlvbnM/OiBib29sZWFuIHwgQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMsXG4gICAgKTogdm9pZDtcbiAgICByZW1vdmVFdmVudExpc3RlbmVyKFxuICAgICAgICB0eXBlOiBcIm1lc3NhZ2VcIixcbiAgICAgICAgbGlzdGVuZXI6ICgodGhpczogSlNPTlNvY2tldCwgZXY6IE1lc3NhZ2VFdmVudDxhbnk+KSA9PiBhbnkpIHwgRXZlbnRMaXN0ZW5lck9iamVjdCB8IG51bGwsXG4gICAgICAgIG9wdGlvbnM/OiBib29sZWFuIHwgRXZlbnRMaXN0ZW5lck9wdGlvbnMsXG4gICAgKTogdm9pZDtcbiAgICByZW1vdmVFdmVudExpc3RlbmVyKFxuICAgICAgICB0eXBlOiBzdHJpbmcsXG4gICAgICAgIGxpc3RlbmVyOiBFdmVudExpc3RlbmVyT3JFdmVudExpc3RlbmVyT2JqZWN0IHwgbnVsbCxcbiAgICAgICAgb3B0aW9ucz86IGJvb2xlYW4gfCBFdmVudExpc3RlbmVyT3B0aW9ucyxcbiAgICApOiB2b2lkO1xufVxuXG4vKipcbiAqIEEge0BsaW5rIFN0cmVhbX0gdGhhdCBzcGVha3Mgb2JqZWN0cyBpbnN0ZWFkIG9mIGJ5dGVzLlxuICpcbiAqIGBzZW5kKHZhbHVlKWAgbWFyc2hhbHMgd2l0aCBgSlNPTi5zdHJpbmdpZnlgLCBhbmQgYGV2LmRhdGFgIGluIGFuIGBvbm1lc3NhZ2VgXG4gKiBoYW5kbGVyIGlzIHRoZSBwYXJzZWQgb2JqZWN0IHJhdGhlciB0aGFuIGFuIGBBcnJheUJ1ZmZlcmAuIFBhaXJzIHdpdGhcbiAqIGBTdHJlYW1Db25uLlNlbmRKU09OYCAvIGBSZWNlaXZlSlNPTmAgb24gdGhlIEdvIHNpZGUuXG4gKlxuICogYGBganNcbiAqIGNvbnN0IHMgPSBKU09OU3RyZWFtKFwidGVsZW1ldHJ5XCIpO1xuICogcy5vbm1lc3NhZ2UgPSAoZXYpID0+IGNvbnNvbGUubG9nKGV2LmRhdGEudGVtcGVyYXR1cmUpO1xuICogcy5vbm9wZW4gICAgPSAoKSA9PiBzLnNlbmQoeyBzdWJzY3JpYmU6IFwic2Vuc29yc1wiIH0pO1xuICogYGBgXG4gKlxuICogQSBmcmFtZSB0aGF0IGlzIG5vdCB2YWxpZCBKU09OIHJhaXNlcyBhbiBgZXJyb3JgIGV2ZW50IGFuZCBpcyBkcm9wcGVkLCByYXRoZXJcbiAqIHRoYW4gdGhyb3dpbmcgZnJvbSB0aGUgcG9sbCBsb29wIGFuZCB0YWtpbmcgdGhlIGNvbm5lY3Rpb24gZG93biB3aXRoIGl0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gSlNPTlN0cmVhbShuYW1lOiBzdHJpbmcpOiBKU09OU29ja2V0IHtcbiAgICBjb25zdCBzb2NrZXQgPSBTdHJlYW0obmFtZSk7XG4gICAgY29uc3QgZGVjb2RlciA9IG5ldyBUZXh0RGVjb2RlcigpO1xuICAgIGNvbnN0IHBhcnNlID0gKHBheWxvYWQ6IEFycmF5QnVmZmVyIHwgc3RyaW5nKSA9PlxuICAgICAgICBKU09OLnBhcnNlKHR5cGVvZiBwYXlsb2FkID09PSBcInN0cmluZ1wiID8gcGF5bG9hZCA6IGRlY29kZXIuZGVjb2RlKHBheWxvYWQpKTtcblxuICAgIGlmIChzb2NrZXQgaW5zdGFuY2VvZiBXYWlsc1NvY2tldCkge1xuICAgICAgICAvLyBEZWNvZGluZyB0aHJvdWdoIHRoZSBob29rIHJhdGhlciB0aGFuIGJ5IHdyYXBwaW5nIG9ubWVzc2FnZSBtZWFuc1xuICAgICAgICAvLyBhZGRFdmVudExpc3RlbmVyIGxpc3RlbmVycyBzZWUgdGhlIHBhcnNlZCBvYmplY3QgdG9vLlxuICAgICAgICBzb2NrZXQuX2RlY29kZSA9IChwYXlsb2FkKSA9PiBwYXJzZShwYXlsb2FkKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICAvLyBTZXJ2ZXIgYnVpbGRzIGdldCBhIG5hdGl2ZSBXZWJTb2NrZXQsIHdob3NlIGRpc3BhdGNoIG5ldmVyIGNvbnN1bHRzXG4gICAgICAgIC8vIF9kZWNvZGUuIFBhdGNoIGJvdGggbGlzdGVuZXIgc3R5bGVzIG9uIHRoaXMgaW5zdGFuY2Ugc28gSlNPTlN0cmVhbVxuICAgICAgICAvLyBtZWFucyB0aGUgc2FtZSB0aGluZyBvbiBlaXRoZXIgdHJhbnNwb3J0LCB3aGljaCBpcyB0aGUgcG9pbnQgb2YgdGhlXG4gICAgICAgIC8vIHNoYXJlZCBBUEkuXG4gICAgICAgIGNvbnN0IG5hdGl2ZSA9IHNvY2tldCBhcyBXZWJTb2NrZXQ7XG4gICAgICAgIG5hdGl2ZS5iaW5hcnlUeXBlID0gXCJhcnJheWJ1ZmZlclwiO1xuXG4gICAgICAgIGNvbnN0IHdyYXBwZWQgPSBuZXcgV2Vha01hcDxvYmplY3QsIEV2ZW50TGlzdGVuZXI+KCk7XG4gICAgICAgIGNvbnN0IGRlY29kZWRFdmVudHMgPSBuZXcgV2Vha01hcDxFdmVudCwgTWVzc2FnZUV2ZW50IHwgbnVsbD4oKTtcbiAgICAgICAgY29uc3QgYWRkID0gbmF0aXZlLmFkZEV2ZW50TGlzdGVuZXIuYmluZChuYXRpdmUpO1xuICAgICAgICBjb25zdCByZW1vdmUgPSBuYXRpdmUucmVtb3ZlRXZlbnRMaXN0ZW5lci5iaW5kKG5hdGl2ZSk7XG5cbiAgICAgICAgY29uc3QgZGVjb2RlRXZlbnQgPSAoZXY6IEV2ZW50KTogTWVzc2FnZUV2ZW50IHwgbnVsbCA9PiB7XG4gICAgICAgICAgICBpZiAoZGVjb2RlZEV2ZW50cy5oYXMoZXYpKSByZXR1cm4gZGVjb2RlZEV2ZW50cy5nZXQoZXYpID8/IG51bGw7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIGNvbnN0IHZhbHVlID0gcGFyc2UoKGV2IGFzIE1lc3NhZ2VFdmVudCkuZGF0YSk7XG4gICAgICAgICAgICAgICAgY29uc3QgZGVjb2RlZCA9IG5ldyBNZXNzYWdlRXZlbnQoXCJtZXNzYWdlXCIsIHsgZGF0YTogdmFsdWUgfSk7XG4gICAgICAgICAgICAgICAgZGVjb2RlZEV2ZW50cy5zZXQoZXYsIGRlY29kZWQpO1xuICAgICAgICAgICAgICAgIHJldHVybiBkZWNvZGVkO1xuICAgICAgICAgICAgfSBjYXRjaCB7XG4gICAgICAgICAgICAgICAgZGVjb2RlZEV2ZW50cy5zZXQoZXYsIG51bGwpO1xuICAgICAgICAgICAgICAgIG5hdGl2ZS5kaXNwYXRjaEV2ZW50KG5ldyBFdmVudChcImVycm9yXCIpKTtcbiAgICAgICAgICAgICAgICByZXR1cm4gbnVsbDtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfTtcblxuICAgICAgICBjb25zdCBkZWNvZGUgPSAobGlzdGVuZXI6IGFueSk6IEV2ZW50TGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgY29uc3QgZXhpc3RpbmcgPSB3cmFwcGVkLmdldChsaXN0ZW5lciBhcyBvYmplY3QpO1xuICAgICAgICAgICAgaWYgKGV4aXN0aW5nKSByZXR1cm4gZXhpc3Rpbmc7XG5cbiAgICAgICAgICAgIGNvbnN0IGZuOiBFdmVudExpc3RlbmVyID0gKGV2KSA9PiB7XG4gICAgICAgICAgICAgICAgY29uc3QgZGVjb2RlZCA9IGRlY29kZUV2ZW50KGV2KTtcbiAgICAgICAgICAgICAgICBpZiAoIWRlY29kZWQpIHJldHVybjtcbiAgICAgICAgICAgICAgICBpZiAodHlwZW9mIGxpc3RlbmVyID09PSBcImZ1bmN0aW9uXCIpIGxpc3RlbmVyLmNhbGwobmF0aXZlLCBkZWNvZGVkKTtcbiAgICAgICAgICAgICAgICBlbHNlIGxpc3RlbmVyLmhhbmRsZUV2ZW50KGRlY29kZWQpO1xuICAgICAgICAgICAgfTtcbiAgICAgICAgICAgIHdyYXBwZWQuc2V0KGxpc3RlbmVyIGFzIG9iamVjdCwgZm4pO1xuICAgICAgICAgICAgcmV0dXJuIGZuO1xuICAgICAgICB9O1xuXG4gICAgICAgIG5hdGl2ZS5hZGRFdmVudExpc3RlbmVyID0gKCh0eXBlOiBzdHJpbmcsIGxpc3RlbmVyOiBhbnksIG9wdHM/OiBhbnkpID0+IHtcbiAgICAgICAgICAgIGFkZCh0eXBlLCB0eXBlID09PSBcIm1lc3NhZ2VcIiAmJiBsaXN0ZW5lciA/IGRlY29kZShsaXN0ZW5lcikgOiBsaXN0ZW5lciwgb3B0cyk7XG4gICAgICAgIH0pIGFzIHR5cGVvZiBuYXRpdmUuYWRkRXZlbnRMaXN0ZW5lcjtcblxuICAgICAgICBuYXRpdmUucmVtb3ZlRXZlbnRMaXN0ZW5lciA9ICgodHlwZTogc3RyaW5nLCBsaXN0ZW5lcjogYW55LCBvcHRzPzogYW55KSA9PiB7XG4gICAgICAgICAgICBjb25zdCBmbiA9IHR5cGUgPT09IFwibWVzc2FnZVwiICYmIGxpc3RlbmVyID8gd3JhcHBlZC5nZXQobGlzdGVuZXIpIDogdW5kZWZpbmVkO1xuICAgICAgICAgICAgcmVtb3ZlKHR5cGUsIGZuID8/IGxpc3RlbmVyLCBvcHRzKTtcbiAgICAgICAgfSkgYXMgdHlwZW9mIG5hdGl2ZS5yZW1vdmVFdmVudExpc3RlbmVyO1xuXG4gICAgICAgIGxldCBoYW5kbGVyOiAoKGV2OiBNZXNzYWdlRXZlbnQpID0+IHZvaWQpIHwgbnVsbCA9IG51bGw7XG4gICAgICAgIE9iamVjdC5kZWZpbmVQcm9wZXJ0eShuYXRpdmUsIFwib25tZXNzYWdlXCIsIHtcbiAgICAgICAgICAgIGdldDogKCkgPT4gaGFuZGxlcixcbiAgICAgICAgICAgIHNldChmbikge1xuICAgICAgICAgICAgICAgIGlmIChoYW5kbGVyKSBuYXRpdmUucmVtb3ZlRXZlbnRMaXN0ZW5lcihcIm1lc3NhZ2VcIiwgaGFuZGxlciBhcyBFdmVudExpc3RlbmVyKTtcbiAgICAgICAgICAgICAgICBoYW5kbGVyID0gdHlwZW9mIGZuID09PSBcImZ1bmN0aW9uXCIgPyBmbiA6IG51bGw7XG4gICAgICAgICAgICAgICAgaWYgKGhhbmRsZXIpIG5hdGl2ZS5hZGRFdmVudExpc3RlbmVyKFwibWVzc2FnZVwiLCBoYW5kbGVyIGFzIEV2ZW50TGlzdGVuZXIpO1xuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogdHJ1ZSxcbiAgICAgICAgICAgIGVudW1lcmFibGU6IHRydWUsXG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIGNvbnN0IHNlbmQgPSBzb2NrZXQuc2VuZC5iaW5kKHNvY2tldCk7XG4gICAgKHNvY2tldCBhcyB1bmtub3duIGFzIEpTT05Tb2NrZXQpLnNlbmQgPSAodmFsdWU6IHVua25vd24pID0+IHNlbmQoSlNPTi5zdHJpbmdpZnkodmFsdWUpKTtcblxuICAgIHJldHVybiBzb2NrZXQgYXMgdW5rbm93biBhcyBKU09OU29ja2V0O1xufVxuXG4vLyBTeW5jaHJvbm91cyBjb252ZXJzaW9uIGZvciBldmVyeXRoaW5nIGJ1dCBhIEJsb2IuIEF2b2lkcyBhIHByb21pc2UgcGVyIHNlbmRcbi8vIGFuZCwgbW9yZSBpbXBvcnRhbnRseSwgbGV0cyBidWZmZXJlZEFtb3VudCBiZSBleGFjdCB0aGUgbW9tZW50IHNlbmQoKSByZXR1cm5zLlxuZnVuY3Rpb24gdG9CeXRlc1N5bmMoZGF0YTogc3RyaW5nIHwgQXJyYXlCdWZmZXJMaWtlIHwgQXJyYXlCdWZmZXJWaWV3IHwgQmxvYik6IFVpbnQ4QXJyYXkgfCBudWxsIHtcbiAgICBpZiAodHlwZW9mIGRhdGEgPT09IFwic3RyaW5nXCIpIHtcbiAgICAgICAgcmV0dXJuIHRleHRFbmNvZGVyLmVuY29kZShkYXRhKTtcbiAgICB9XG4gICAgaWYgKHR5cGVvZiBCbG9iICE9PSBcInVuZGVmaW5lZFwiICYmIGRhdGEgaW5zdGFuY2VvZiBCbG9iKSB7XG4gICAgICAgIHJldHVybiBudWxsO1xuICAgIH1cbiAgICBpZiAoQXJyYXlCdWZmZXIuaXNWaWV3KGRhdGEpKSB7XG4gICAgICAgIHJldHVybiBuZXcgVWludDhBcnJheShkYXRhLmJ1ZmZlciwgZGF0YS5ieXRlT2Zmc2V0LCBkYXRhLmJ5dGVMZW5ndGgpLnNsaWNlKCk7XG4gICAgfVxuICAgIHJldHVybiBuZXcgVWludDhBcnJheShkYXRhIGFzIEFycmF5QnVmZmVyTGlrZSkuc2xpY2UoKTtcbn1cblxuZnVuY3Rpb24gYWJvcnRhYmxlPFQ+KHByb21pc2U6IFByb21pc2U8VD4sIHNpZ25hbD86IEFib3J0U2lnbmFsKTogUHJvbWlzZTxUPiB7XG4gICAgaWYgKCFzaWduYWwpIHJldHVybiBwcm9taXNlO1xuICAgIGlmIChzaWduYWwuYWJvcnRlZCkgcmV0dXJuIFByb21pc2UucmVqZWN0KG5ldyBTdHJlYW1SZXF1ZXN0Q2FuY2VsbGVkKCkpO1xuXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgY29uc3Qgb25BYm9ydCA9ICgpID0+IHJlamVjdChuZXcgU3RyZWFtUmVxdWVzdENhbmNlbGxlZCgpKTtcbiAgICAgICAgc2lnbmFsLmFkZEV2ZW50TGlzdGVuZXIoXCJhYm9ydFwiLCBvbkFib3J0LCB7IG9uY2U6IHRydWUgfSk7XG4gICAgICAgIHByb21pc2UudGhlbihcbiAgICAgICAgICAgICh2YWx1ZSkgPT4ge1xuICAgICAgICAgICAgICAgIHNpZ25hbC5yZW1vdmVFdmVudExpc3RlbmVyKFwiYWJvcnRcIiwgb25BYm9ydCk7XG4gICAgICAgICAgICAgICAgcmVzb2x2ZSh2YWx1ZSk7XG4gICAgICAgICAgICB9LFxuICAgICAgICAgICAgKGVycikgPT4ge1xuICAgICAgICAgICAgICAgIHNpZ25hbC5yZW1vdmVFdmVudExpc3RlbmVyKFwiYWJvcnRcIiwgb25BYm9ydCk7XG4gICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XG4gICAgICAgICAgICB9LFxuICAgICAgICApO1xuICAgIH0pO1xufVxuXG5hc3luYyBmdW5jdGlvbiB0b0J5dGVzKFxuICAgIGRhdGE6IHN0cmluZyB8IEFycmF5QnVmZmVyTGlrZSB8IEFycmF5QnVmZmVyVmlldyB8IEJsb2IsXG4gICAgc2lnbmFsPzogQWJvcnRTaWduYWwsXG4pOiBQcm9taXNlPFVpbnQ4QXJyYXk+IHtcbiAgICBpZiAodHlwZW9mIGRhdGEgPT09IFwic3RyaW5nXCIpIHtcbiAgICAgICAgcmV0dXJuIHRleHRFbmNvZGVyLmVuY29kZShkYXRhKTtcbiAgICB9XG4gICAgaWYgKHR5cGVvZiBCbG9iICE9PSBcInVuZGVmaW5lZFwiICYmIGRhdGEgaW5zdGFuY2VvZiBCbG9iKSB7XG4gICAgICAgIHJldHVybiBuZXcgVWludDhBcnJheShhd2FpdCBhYm9ydGFibGUoZGF0YS5hcnJheUJ1ZmZlcigpLCBzaWduYWwpKTtcbiAgICB9XG4gICAgaWYgKEFycmF5QnVmZmVyLmlzVmlldyhkYXRhKSkge1xuICAgICAgICByZXR1cm4gbmV3IFVpbnQ4QXJyYXkoZGF0YS5idWZmZXIsIGRhdGEuYnl0ZU9mZnNldCwgZGF0YS5ieXRlTGVuZ3RoKS5zbGljZSgpO1xuICAgIH1cbiAgICByZXR1cm4gbmV3IFVpbnQ4QXJyYXkoZGF0YSBhcyBBcnJheUJ1ZmZlckxpa2UpLnNsaWNlKCk7XG59XG5cbmFzeW5jIGZ1bmN0aW9uIHBvc3RGcmFtZShjb25uSUQ6IG51bWJlciwga2luZDogbnVtYmVyLCBib2R5OiBVaW50OEFycmF5LCBuYW1lPzogc3RyaW5nLCBzaWduYWw/OiBBYm9ydFNpZ25hbCk6IFByb21pc2U8dm9pZD4ge1xuICAgIGNvbnN0IGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4gPSB7XG4gICAgICAgIFtIRFJfU0VTU0lPTl06IEcuc2Vzc2lvbixcbiAgICAgICAgW0hEUl9HRU5FUkFUSU9OXTogU3RyaW5nKEcuZ2VuZXJhdGlvbiksXG4gICAgICAgIFtIRFJfQ09OTl06IFN0cmluZyhjb25uSUQpLFxuICAgICAgICBbSERSX0tJTkRdOiBTdHJpbmcoa2luZCksXG4gICAgICAgIFwiQ29udGVudC1UeXBlXCI6IFwiYXBwbGljYXRpb24vb2N0ZXQtc3RyZWFtXCIsXG4gICAgfTtcbiAgICBpZiAobmFtZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIGhlYWRlcnNbSERSX05BTUVdID0gbmFtZTtcbiAgICB9XG5cbiAgICBpZiAoYm9keS5ieXRlTGVuZ3RoIDw9IENIVU5LX1RIUkVTSE9MRCkge1xuICAgICAgICBhd2FpdCBwb3N0V2l0aFJldHJ5KGhlYWRlcnMsIGJvZHksIHNpZ25hbCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBTcGxpdCBvdmVyc2l6ZWQgZnJhbWVzIHRoZSBzYW1lIHdheSB0aGUgcnVudGltZSBzcGxpdHMgb3ZlcnNpemVkIGNhbGxzLlxuICAgIGNvbnN0IGNodW5rSUQgPSBuYW5vaWQoKTtcbiAgICBjb25zdCB0b3RhbCA9IE1hdGguY2VpbChib2R5LmJ5dGVMZW5ndGggLyBDSFVOS19USFJFU0hPTEQpO1xuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgdG90YWw7IGkrKykge1xuICAgICAgICBjb25zdCBzbGljZSA9IGJvZHkuc3ViYXJyYXkoaSAqIENIVU5LX1RIUkVTSE9MRCwgKGkgKyAxKSAqIENIVU5LX1RIUkVTSE9MRCk7XG4gICAgICAgIGF3YWl0IHBvc3RXaXRoUmV0cnkoe1xuICAgICAgICAgICAgLi4uaGVhZGVycyxcbiAgICAgICAgICAgIFtIRFJfQ0hVTktdOiBjaHVua0lELFxuICAgICAgICAgICAgW0hEUl9DSFVOS19JTkRFWF06IFN0cmluZyhpKSxcbiAgICAgICAgICAgIFtIRFJfQ0hVTktfVE9UQUxdOiBTdHJpbmcodG90YWwpLFxuICAgICAgICB9LCBzbGljZSwgc2lnbmFsKTtcbiAgICB9XG59XG5cbi8vIEEgNDI5IG1lYW5zIHRoZSBHbyBoYW5kbGVyIGhhcyBub3QgdGFrZW4gZGVsaXZlcnkgb2Ygd2hhdCBpcyBhbHJlYWR5IHF1ZXVlZC5cbi8vIFJldHJ5aW5nIHRoZSBzYW1lIGZyYW1lIGtlZXBzIG9yZGVyaW5nICh0aGUgc2VuZCBjaGFpbiBoYXMgbm90IG1vdmVkIG9uKSBhbmRcbi8vIGxlYXZlcyB0aGUgcmVxdWVzdCBzbG90IGZyZWUsIHdoaWNoIGlzIHdoYXQgc3RvcHMgYSBiYWNrZWQtdXAgY29ubmVjdGlvbiBmcm9tXG4vLyBzdGFydmluZyB0aGUgd2luZG93J3MgcG9sbC5cbmFzeW5jIGZ1bmN0aW9uIHBvc3RXaXRoUmV0cnkoaGVhZGVyczogUmVjb3JkPHN0cmluZywgc3RyaW5nPiwgYm9keTogVWludDhBcnJheSwgc2lnbmFsPzogQWJvcnRTaWduYWwpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAvLyBOZXZlciByZXRyeSB3aXRoIHplcm8gZGVsYXkuIFRoZSByZWNlaXZlciBiZWluZyBiZWhpbmQgaXMgYSBjb25kaXRpb25cbiAgICAvLyB0aGF0IHRha2VzIHRpbWUgdG8gY2xlYXIsIGFuZCBhbiBpbW1lZGlhdGUgcmV0cnkgaXMgYSBidXN5IGxvb3Agb2ZcbiAgICAvLyBmZXRjaGVzIFx1MjAxNCBlYWNoIG9uZSBhIHNjaGVtZS1oYW5kbGVyIHJvdW5kIHRyaXAsIGFuZCBvbiB0aGUgaG9zdCBhIGNnb1xuICAgIC8vIGNhbGwuIFN0YXJ0IGF0IDEgbXMgYW5kIHJhbXAuXG4gICAgbGV0IHdhaXQgPSAxO1xuICAgIGZvciAoOzspIHtcbiAgICAgICAgaWYgKHNpZ25hbD8uYWJvcnRlZCkgdGhyb3cgbmV3IFN0cmVhbVJlcXVlc3RDYW5jZWxsZWQoKTtcbiAgICAgICAgY29uc3QgcmVzcCA9IGF3YWl0IGZldGNoKHN0cmVhbVVSTChcInNlbmRcIiksIHtcbiAgICAgICAgICAgIG1ldGhvZDogXCJQT1NUXCIsXG4gICAgICAgICAgICBoZWFkZXJzLFxuICAgICAgICAgICAgYm9keTogYm9keSBhcyBCb2R5SW5pdCxcbiAgICAgICAgICAgIHNpZ25hbCxcbiAgICAgICAgfSk7XG4gICAgICAgIGlmIChyZXNwLm9rKSByZXR1cm47XG4gICAgICAgIGlmIChyZXNwLnN0YXR1cyAhPT0gNDI5KSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xuICAgICAgICB9XG4gICAgICAgIGF3YWl0IHNsZWVwKHdhaXQsIHNpZ25hbCk7XG4gICAgICAgIHdhaXQgPSBNYXRoLm1pbih3YWl0ICogMiwgNTApO1xuICAgIH1cbn1cblxuLy8gcG9zdEJhdGNoIHNlbmRzIHNldmVyYWwgZGF0YSBmcmFtZXMgZm9yIG9uZSBjb25uZWN0aW9uIGluIGJvdW5kZWQgcmVxdWVzdHMuXG4vLyBCb2R5OiBjb3VudCB1MzIsIHRoZW4gY291bnQgeCAoIGxlbiB1MzIgfCBwYXlsb2FkICkuIFRoZSBjb21iaW5lZCBib2R5IHN0YXlzXG4vLyB3aXRoaW4gQ0hVTktfVEhSRVNIT0xEOiBzZXZlcmFsIGluZGl2aWR1YWxseSBzbWFsbCBmcmFtZXMgbXVzdCBub3QgcmVjcmVhdGVcbi8vIHRoZSBXZWJWaWV3MiBib2R5LXNpemUgcHJvYmxlbSB0aGF0IGNodW5raW5nIHNpbmdsZSBsYXJnZSBmcmFtZXMgYXZvaWRzLlxuYXN5bmMgZnVuY3Rpb24gcG9zdEJhdGNoKGNvbm5JRDogbnVtYmVyLCBmcmFtZXM6IFVpbnQ4QXJyYXlbXSwgc2lnbmFsPzogQWJvcnRTaWduYWwpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBsZXQgc3RhcnQgPSAwO1xuICAgIHdoaWxlIChzdGFydCA8IGZyYW1lcy5sZW5ndGgpIHtcbiAgICAgICAgaWYgKGZyYW1lc1tzdGFydF0uYnl0ZUxlbmd0aCA+IENIVU5LX1RIUkVTSE9MRCkge1xuICAgICAgICAgICAgYXdhaXQgcG9zdEZyYW1lKGNvbm5JRCwgS0lORF9EQVRBLCBmcmFtZXNbc3RhcnRdLCB1bmRlZmluZWQsIHNpZ25hbCk7XG4gICAgICAgICAgICBzdGFydCsrO1xuICAgICAgICAgICAgY29udGludWU7XG4gICAgICAgIH1cblxuICAgICAgICBsZXQgZW5kID0gc3RhcnQ7XG4gICAgICAgIGxldCBzaXplID0gNDsgLy8gZnJhbWUgY291bnRcbiAgICAgICAgd2hpbGUgKGVuZCA8IGZyYW1lcy5sZW5ndGggJiYgZW5kIC0gc3RhcnQgPCBNQVhfQkFUQ0hfRlJBTUVTKSB7XG4gICAgICAgICAgICBjb25zdCBmcmFtZSA9IGZyYW1lc1tlbmRdO1xuICAgICAgICAgICAgaWYgKGZyYW1lLmJ5dGVMZW5ndGggPiBDSFVOS19USFJFU0hPTEQgfHwgc2l6ZSArIDQgKyBmcmFtZS5ieXRlTGVuZ3RoID4gQ0hVTktfVEhSRVNIT0xEKSB7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBzaXplICs9IDQgKyBmcmFtZS5ieXRlTGVuZ3RoO1xuICAgICAgICAgICAgZW5kKys7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBBIG5lYXItdGhyZXNob2xkIGZyYW1lIGZpdHMgYXMgYSBub3JtYWwgUE9TVCBib2R5IGJ1dCBub3QgaW5zaWRlIGFcbiAgICAgICAgLy8gYmF0Y2ggb25jZSBpdHMgY291bnQgYW5kIGxlbmd0aCBoZWFkZXJzIGFyZSBpbmNsdWRlZC5cbiAgICAgICAgaWYgKGVuZCA9PT0gc3RhcnQpIHtcbiAgICAgICAgICAgIGF3YWl0IHBvc3RGcmFtZShjb25uSUQsIEtJTkRfREFUQSwgZnJhbWVzW3N0YXJ0XSwgdW5kZWZpbmVkLCBzaWduYWwpO1xuICAgICAgICAgICAgc3RhcnQrKztcbiAgICAgICAgICAgIGNvbnRpbnVlO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgZ3JvdXAgPSBmcmFtZXMuc2xpY2Uoc3RhcnQsIGVuZCk7XG4gICAgICAgIGlmIChncm91cC5sZW5ndGggPT09IDEpIHtcbiAgICAgICAgICAgIGF3YWl0IHBvc3RGcmFtZShjb25uSUQsIEtJTkRfREFUQSwgZ3JvdXBbMF0sIHVuZGVmaW5lZCwgc2lnbmFsKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGF3YWl0IHBvc3RCYXRjaEdyb3VwKGNvbm5JRCwgZ3JvdXAsIHNpZ25hbCk7XG4gICAgICAgIH1cbiAgICAgICAgc3RhcnQgPSBlbmQ7XG4gICAgfVxufVxuXG5hc3luYyBmdW5jdGlvbiBwb3N0QmF0Y2hHcm91cChjb25uSUQ6IG51bWJlciwgZnJhbWVzOiBVaW50OEFycmF5W10sIHNpZ25hbD86IEFib3J0U2lnbmFsKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgaWYgKGZyYW1lcy5sZW5ndGggPT09IDAgfHwgZnJhbWVzLmxlbmd0aCA+IE1BWF9CQVRDSF9GUkFNRVMpIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiaW52YWxpZCBzdHJlYW0gYmF0Y2ggc2l6ZVwiKTtcbiAgICB9XG5cbiAgICBsZXQgYm9keSA9IGJ1aWxkQmF0Y2goZnJhbWVzKTtcbiAgICBsZXQgc2VudCA9IDA7XG4gICAgbGV0IHdhaXQgPSAxO1xuICAgIGZvciAoOzspIHtcbiAgICAgICAgaWYgKHNpZ25hbD8uYWJvcnRlZCkgdGhyb3cgbmV3IFN0cmVhbVJlcXVlc3RDYW5jZWxsZWQoKTtcbiAgICAgICAgY29uc3QgcmVzcCA9IGF3YWl0IGZldGNoKHN0cmVhbVVSTChcInNlbmRcIiksIHtcbiAgICAgICAgICAgIG1ldGhvZDogXCJQT1NUXCIsXG4gICAgICAgICAgICBoZWFkZXJzOiB7XG4gICAgICAgICAgICAgICAgW0hEUl9TRVNTSU9OXTogRy5zZXNzaW9uLFxuICAgICAgICAgICAgICAgIFtIRFJfR0VORVJBVElPTl06IFN0cmluZyhHLmdlbmVyYXRpb24pLFxuICAgICAgICAgICAgICAgIFtIRFJfQ09OTl06IFN0cmluZyhjb25uSUQpLFxuICAgICAgICAgICAgICAgIFtIRFJfS0lORF06IFN0cmluZyhLSU5EX0RBVEEpLFxuICAgICAgICAgICAgICAgIFtIRFJfQkFUQ0hdOiBTdHJpbmcoZnJhbWVzLmxlbmd0aCAtIHNlbnQpLFxuICAgICAgICAgICAgICAgIFwiQ29udGVudC1UeXBlXCI6IFwiYXBwbGljYXRpb24vb2N0ZXQtc3RyZWFtXCIsXG4gICAgICAgICAgICB9LFxuICAgICAgICAgICAgYm9keTogYm9keSBhcyBCb2R5SW5pdCxcbiAgICAgICAgICAgIHNpZ25hbCxcbiAgICAgICAgfSk7XG4gICAgICAgIGlmIChyZXNwLm9rKSByZXR1cm47XG4gICAgICAgIGlmIChyZXNwLnN0YXR1cyAhPT0gNDI5KSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYXdhaXQgcmVzcC50ZXh0KCkpO1xuICAgICAgICB9XG4gICAgICAgIC8vIFRoZSByZWNlaXZlciB0b29rIGEgcHJlZml4IG9mIHRoZSBiYXRjaC4gUmVzZW5kIG9ubHkgd2hhdCBpcyBsZWZ0LFxuICAgICAgICAvLyB3aGljaCBrZWVwcyBvcmRlcmluZyB3aXRob3V0IHJlLWRlbGl2ZXJpbmcgYW55dGhpbmcuXG4gICAgICAgIGNvbnN0IGFjY2VwdGVkID0gTnVtYmVyKHJlc3AuaGVhZGVycy5nZXQoSERSX0JBVENIKSA/PyAwKTtcbiAgICAgICAgY29uc3QgcmVtYWluaW5nID0gZnJhbWVzLmxlbmd0aCAtIHNlbnQ7XG4gICAgICAgIGlmICghTnVtYmVyLmlzSW50ZWdlcihhY2NlcHRlZCkgfHwgYWNjZXB0ZWQgPCAwIHx8IGFjY2VwdGVkID49IHJlbWFpbmluZykge1xuICAgICAgICAgICAgLy8gWmVybyBpcyB2YWxpZCAobm8gcHJvZ3Jlc3MpLCBidXQgYW4gb3V0LW9mLXJhbmdlIGFja25vd2xlZGdlbWVudFxuICAgICAgICAgICAgLy8gaXMgYSBwcm90b2NvbCBlcnJvci4gS2VlcCB6ZXJvIG9uIHRoZSByZXRyeSBwYXRoIGJlbG93LlxuICAgICAgICAgICAgaWYgKGFjY2VwdGVkICE9PSAwKSB0aHJvdyBuZXcgRXJyb3IoXCJpbnZhbGlkIHN0cmVhbSBiYXRjaCBhY2tub3dsZWRnZW1lbnRcIik7XG4gICAgICAgIH1cbiAgICAgICAgc2VudCArPSBhY2NlcHRlZDtcbiAgICAgICAgYm9keSA9IGJ1aWxkQmF0Y2goZnJhbWVzLnNsaWNlKHNlbnQpKTtcbiAgICAgICAgYXdhaXQgc2xlZXAod2FpdCwgc2lnbmFsKTtcbiAgICAgICAgd2FpdCA9IE1hdGgubWluKHdhaXQgKiAyLCA1MCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBidWlsZEJhdGNoKGZyYW1lczogVWludDhBcnJheVtdKTogVWludDhBcnJheSB7XG4gICAgbGV0IHNpemUgPSA0O1xuICAgIGZvciAoY29uc3QgZiBvZiBmcmFtZXMpIHNpemUgKz0gNCArIGYuYnl0ZUxlbmd0aDtcbiAgICBjb25zdCBvdXQgPSBuZXcgVWludDhBcnJheShzaXplKTtcbiAgICBjb25zdCB2aWV3ID0gbmV3IERhdGFWaWV3KG91dC5idWZmZXIpO1xuICAgIHZpZXcuc2V0VWludDMyKDAsIGZyYW1lcy5sZW5ndGgpO1xuICAgIGxldCBvZmYgPSA0O1xuICAgIGZvciAoY29uc3QgZiBvZiBmcmFtZXMpIHtcbiAgICAgICAgdmlldy5zZXRVaW50MzIob2ZmLCBmLmJ5dGVMZW5ndGgpO1xuICAgICAgICBvZmYgKz0gNDtcbiAgICAgICAgb3V0LnNldChmLCBvZmYpO1xuICAgICAgICBvZmYgKz0gZi5ieXRlTGVuZ3RoO1xuICAgIH1cbiAgICByZXR1cm4gb3V0O1xufVxuXG5mdW5jdGlvbiBzdGFydFBvbGxpbmcoKTogdm9pZCB7XG4gICAgaWYgKEcucG9sbGluZyB8fCAhaGFzRE9NKSByZXR1cm47XG4gICAgRy5wb2xsaW5nID0gdHJ1ZTtcbiAgICB2b2lkIHBvbGxMb29wKCk7XG59XG5cbi8qKlxuICogT25lIHBvbGwgZm9yIHRoZSB3aG9sZSBwYWdlLCBzaGFyZWQgYnkgZXZlcnkgY29ubmVjdGlvbi5cbiAqXG4gKiBUaGVyZSBpcyBubyBwb2xsaW5nIGludGVydmFsIGFuZCBub3RoaW5nIGFkYXB0aXZlIGhlcmUgb24gcHVycG9zZS4gVGhlIHNlcnZlclxuICogaG9sZHMgdGhlIHJlcXVlc3QgdW50aWwgYSBmcmFtZSBleGlzdHMsIHNvIGRlbGl2ZXJ5IGxhdGVuY3kgaXMgYWxyZWFkeSB+MCBhbmRcbiAqIGEgY2xpZW50LXNpZGUgaW50ZXJ2YWwgY291bGQgb25seSBhZGQgdG8gaXQuIEZyYW1lcyB0aGF0IGFycml2ZSB3aGlsZSBhXG4gKiByZXNwb25zZSBpcyBpbiBmbGlnaHQgc2ltcGx5IGFjY3VtdWxhdGUgYW5kIHJpZGUgdGhlIG5leHQgb25lLCB3aGljaCBtYWtlc1xuICogdGhlIHJvdW5kIHRyaXAgaXRzZWxmIHRoZSBiYXRjaGluZyB3aW5kb3cgXHUyMDE0IGl0IHdpZGVucyBhcyBsb2FkIHJpc2VzIHdpdGhvdXRcbiAqIGFueXRoaW5nIGhhdmluZyB0byBtZWFzdXJlIGl0LlxuICovXG5hc3luYyBmdW5jdGlvbiBwb2xsTG9vcCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBsZXQgYmFja29mZiA9IDA7XG4gICAgY29uc3QgYWJvcnQgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XG4gICAgRy5wb2xsQWJvcnQgPSBhYm9ydDtcblxuICAgIHdoaWxlIChHLmNvbm5lY3Rpb25zLnNpemUgPiAwKSB7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBjb25zdCByZXNwID0gYXdhaXQgZmV0Y2goc3RyZWFtVVJMKFwicG9sbFwiKSwge1xuICAgICAgICAgICAgICAgIG1ldGhvZDogXCJHRVRcIixcbiAgICAgICAgICAgICAgICBoZWFkZXJzOiB7XG4gICAgICAgICAgICAgICAgICAgIFtIRFJfU0VTU0lPTl06IEcuc2Vzc2lvbixcbiAgICAgICAgICAgICAgICAgICAgW0hEUl9HRU5FUkFUSU9OXTogU3RyaW5nKEcuZ2VuZXJhdGlvbiksXG4gICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICBjYWNoZTogXCJuby1zdG9yZVwiLFxuICAgICAgICAgICAgICAgIHNpZ25hbDogYWJvcnQuc2lnbmFsLFxuICAgICAgICAgICAgfSk7XG5cbiAgICAgICAgICAgIGlmIChyZXNwLnN0YXR1cyA9PT0gNDEwKSB7XG4gICAgICAgICAgICAgICAgLy8gVGhlIHNlc3Npb24gaXMgZ29uZTogYXBwIHNodXR0aW5nIGRvd24sIHdpbmRvdyBjbG9zZWQsIG9yIHdlXG4gICAgICAgICAgICAgICAgLy8gd2VyZSByZWFwZWQgYXMgdW5yZWFjaGFibGUuIE5vdGhpbmcgdG8gcmV0cnkuXG4gICAgICAgICAgICAgICAgY2xvc2VBbGwoMTAwMSwgXCJzZXNzaW9uIGNsb3NlZFwiKTtcbiAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGlmICghcmVzcC5vayAmJiByZXNwLnN0YXR1cyAhPT0gMjA0KSB7XG4gICAgICAgICAgICAgICAgaWYgKCFpc1JldHJ5YWJsZVBvbGxTdGF0dXMocmVzcC5zdGF0dXMpKSB7XG4gICAgICAgICAgICAgICAgICAgIGZhaWxBbGwoXCJzdHJlYW0gcG9sbCBmYWlsZWQ6IFwiICsgcmVzcC5zdGF0dXMpO1xuICAgICAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwicmV0cnlhYmxlIHN0cmVhbSBwb2xsIGZhaWx1cmU6IFwiICsgcmVzcC5zdGF0dXMpO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICBiYWNrb2ZmID0gMDtcbiAgICAgICAgICAgIGlmIChyZXNwLnN0YXR1cyAhPT0gMjA0KSB7XG4gICAgICAgICAgICAgICAgbGV0IGJ1ZjogQXJyYXlCdWZmZXI7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgYnVmID0gYXdhaXQgcmVzcC5hcnJheUJ1ZmZlcigpO1xuICAgICAgICAgICAgICAgIH0gY2F0Y2gge1xuICAgICAgICAgICAgICAgICAgICAvLyBUaGUgR28gcXVldWUgd2FzIGNvbnN1bWVkIHdoZW4gdGhlIHN1Y2Nlc3NmdWwgcmVzcG9uc2VcbiAgICAgICAgICAgICAgICAgICAgLy8gd2FzIHByb2R1Y2VkLiBSZXRyeWluZyBjYW5ub3QgcmVjb3ZlciBhIGJvZHkgdGhhdCB3YXNcbiAgICAgICAgICAgICAgICAgICAgLy8gaW50ZXJydXB0ZWQgYWZ0ZXIgdGhhdCBwb2ludCwgc28gc3VyZmFjZSB0aGUgbG9zcyBhbmRcbiAgICAgICAgICAgICAgICAgICAgLy8gc3RvcCBpbnN0ZWFkIG9mIHNpbGVudGx5IGNvbnRpbnVpbmcgd2l0aCBtaXNzaW5nIGZyYW1lcy5cbiAgICAgICAgICAgICAgICAgICAgZmFpbEFsbChcInN0cmVhbSBwb2xsIHJlc3BvbnNlIGJvZHkgZmFpbGVkXCIpO1xuICAgICAgICAgICAgICAgICAgICBicmVhaztcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgY29uc3QgZnJhbWVzID0gZGVjb2RlRnJhbWVzKGJ1Zik7XG4gICAgICAgICAgICAgICAgaWYgKGZyYW1lcyA9PT0gbnVsbCkge1xuICAgICAgICAgICAgICAgICAgICAvLyBBIGZyYW1pbmcgbWlzbWF0Y2ggbWVhbnMgYSBzdGFsZSBjYWNoZWQgcnVudGltZSBhZ2FpbnN0IGFcbiAgICAgICAgICAgICAgICAgICAgLy8gbmV3ZXIgR28gc2lkZS4gUmV0cnlpbmcgY2Fubm90IGZpeCB0aGF0LCBhbmQgbGV0dGluZyBpdFxuICAgICAgICAgICAgICAgICAgICAvLyBmYWxsIGludG8gdGhlIGJhY2tvZmYgYmVsb3cgd291bGQgc3BpbiBzaWxlbnRseSBmb3JldmVyLFxuICAgICAgICAgICAgICAgICAgICAvLyBzbyBmYWlsIGV2ZXJ5IGNvbm5lY3Rpb24gd2l0aCBhIHJlYXNvbiBpbnN0ZWFkLlxuICAgICAgICAgICAgICAgICAgICBjbG9zZUFsbCgxMDAyLCBcInVucmVjb2duaXNlZCBzdHJlYW0gZnJhbWluZ1wiKTtcbiAgICAgICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIGRlbGl2ZXIoZnJhbWVzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIFJlLXBvbGwgaW1tZWRpYXRlbHksIG9uIGRhdGEgYW5kIG9uIGFuIGV4cGlyZWQgaG9sZCBhbGlrZS5cbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBpZiAoYWJvcnQuc2lnbmFsLmFib3J0ZWQgfHwgRy5jb25uZWN0aW9ucy5zaXplID09PSAwKSBicmVhaztcbiAgICAgICAgICAgIGJhY2tvZmYgPSBiYWNrb2ZmID8gTWF0aC5taW4oYmFja29mZiAqIDIsIEJBQ0tPRkZfTUFYKSA6IEJBQ0tPRkZfTUlOO1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBhd2FpdCBzbGVlcChiYWNrb2ZmLCBhYm9ydC5zaWduYWwpO1xuICAgICAgICAgICAgfSBjYXRjaCB7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICBpZiAoRy5wb2xsQWJvcnQgPT09IGFib3J0KSBHLnBvbGxBYm9ydCA9IHVuZGVmaW5lZDtcbiAgICBHLnBvbGxpbmcgPSBmYWxzZTtcblxuICAgIC8vIEEgY29ubmVjdGlvbiBvcGVuZWQgd2hpbGUgd2Ugd2VyZSB3aW5kaW5nIGRvd246IHBpY2sgdGhlIGxvb3AgYmFjayB1cC5cbiAgICBpZiAoRy5jb25uZWN0aW9ucy5zaXplID4gMCkge1xuICAgICAgICBzdGFydFBvbGxpbmcoKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGlzUmV0cnlhYmxlUG9sbFN0YXR1cyhzdGF0dXM6IG51bWJlcik6IGJvb2xlYW4ge1xuICAgIHJldHVybiBzdGF0dXMgPT09IDQwOCB8fCBzdGF0dXMgPT09IDQyNSB8fCBzdGF0dXMgPT09IDQyOSB8fCAoc3RhdHVzID49IDUwMCAmJiBzdGF0dXMgPD0gNTk5KTtcbn1cblxuaW50ZXJmYWNlIEluRnJhbWUge1xuICAgIGNvbm5JRDogbnVtYmVyO1xuICAgIGtpbmQ6IG51bWJlcjtcbiAgICBwYXlsb2FkOiBBcnJheUJ1ZmZlcjtcbn1cblxuLy8gVmFsaWRhdGUgdGhlIGNvbXBsZXRlIHJlc3BvbnNlIGJlZm9yZSBkaXNwYXRjaGluZyBhbnkgcGFydCBvZiBpdC4gQSBzaG9ydFxuLy8gcGxhdGZvcm0gcmVzcG9uc2UgbXVzdCBub3QgZGVsaXZlciBhIHZhbGlkIHByZWZpeCBhbmQgdGhlbiB0aHJvdyBSYW5nZUVycm9yOlxuLy8gdGhlIEdvIHF1ZXVlIGhhcyBhbHJlYWR5IGJlZW4gZHJhaW5lZCwgc28gcmV0cnlpbmcgY2Fubm90IHJlY292ZXIgdGhvc2Vcbi8vIGZyYW1lcy4gVHJlYXQgdGhlIHdob2xlIHJlc3BvbnNlIGFzIGEgcHJvdG9jb2wgZXJyb3IgaW5zdGVhZC5cbmZ1bmN0aW9uIGRlY29kZUZyYW1lcyhidWY6IEFycmF5QnVmZmVyKTogSW5GcmFtZVtdIHwgbnVsbCB7XG4gICAgaWYgKGJ1Zi5ieXRlTGVuZ3RoIDwgOSkgcmV0dXJuIG51bGw7XG4gICAgY29uc3QgZHYgPSBuZXcgRGF0YVZpZXcoYnVmKTtcbiAgICBpZiAoZHYuZ2V0VWludDgoMCkgIT09IDB4NTcgfHwgZHYuZ2V0VWludDgoMSkgIT09IDB4NTMgfHxcbiAgICAgICAgZHYuZ2V0VWludDgoMikgIT09IDB4MzEgfHwgZHYuZ2V0VWludDgoMykgIT09IDB4MDAgfHxcbiAgICAgICAgKGR2LmdldFVpbnQ4KDQpICYgfjEpICE9PSAwKSB7XG4gICAgICAgIHJldHVybiBudWxsO1xuICAgIH1cblxuICAgIGNvbnN0IGNvdW50ID0gZHYuZ2V0VWludDMyKDUpO1xuICAgIGNvbnN0IGZyYW1lczogSW5GcmFtZVtdID0gW107XG4gICAgbGV0IG9mZiA9IDk7XG5cbiAgICBmb3IgKGxldCBpID0gMDsgaSA8IGNvdW50OyBpKyspIHtcbiAgICAgICAgaWYgKG9mZiArIDkgPiBidWYuYnl0ZUxlbmd0aCkgcmV0dXJuIG51bGw7XG4gICAgICAgIGNvbnN0IGNvbm5JRCA9IGR2LmdldFVpbnQzMihvZmYpOyBvZmYgKz0gNDtcbiAgICAgICAgY29uc3Qga2luZCA9IGR2LmdldFVpbnQ4KG9mZik7IG9mZiArPSAxO1xuICAgICAgICBjb25zdCBsZW4gPSBkdi5nZXRVaW50MzIob2ZmKTsgb2ZmICs9IDQ7XG4gICAgICAgIGlmIChraW5kID4gS0lORF9FUlJPUiB8fCBsZW4gPiBidWYuYnl0ZUxlbmd0aCAtIG9mZikgcmV0dXJuIG51bGw7XG4gICAgICAgIGNvbnN0IHBheWxvYWQgPSBidWYuc2xpY2Uob2ZmLCBvZmYgKyBsZW4pOyBvZmYgKz0gbGVuO1xuICAgICAgICBmcmFtZXMucHVzaCh7IGNvbm5JRCwga2luZCwgcGF5bG9hZCB9KTtcbiAgICB9XG5cbiAgICByZXR1cm4gb2ZmID09PSBidWYuYnl0ZUxlbmd0aCA/IGZyYW1lcyA6IG51bGw7XG59XG5cbmZ1bmN0aW9uIGRlbGl2ZXIoZnJhbWVzOiBJbkZyYW1lW10pOiB2b2lkIHtcbiAgICBmb3IgKGNvbnN0IGZyYW1lIG9mIGZyYW1lcykge1xuICAgICAgICBjb25zdCBjb25uSUQgPSBmcmFtZS5jb25uSUQ7XG4gICAgICAgIGNvbnN0IGtpbmQgPSBmcmFtZS5raW5kO1xuICAgICAgICBjb25zdCBwYXlsb2FkID0gZnJhbWUucGF5bG9hZDtcbiAgICAgICAgY29uc3QgY29ubiA9IEcuY29ubmVjdGlvbnMuZ2V0KGNvbm5JRCk7XG4gICAgICAgIGlmICghY29ubikgY29udGludWU7XG5cbiAgICAgICAgc3dpdGNoIChraW5kKSB7XG4gICAgICAgICAgICBjYXNlIEtJTkRfREFUQTpcbiAgICAgICAgICAgICAgICBjb25uLl9tZXNzYWdlKHBheWxvYWQpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgY2FzZSBLSU5EX09QRU46XG4gICAgICAgICAgICAgICAgY29ubi5fb3BlbmVkKCk7XG4gICAgICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgICAgICBjYXNlIEtJTkRfQ0xPU0U6XG4gICAgICAgICAgICAgICAgY29ubi5fY2xvc2VkKDEwMDAsIFwiXCIsIHRydWUpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICAgICAgY2FzZSBLSU5EX0VSUk9SOlxuICAgICAgICAgICAgICAgIGNvbm4uX2ZhaWwobmV3IEVycm9yKG5ldyBUZXh0RGVjb2RlcigpLmRlY29kZShwYXlsb2FkKSkpO1xuICAgICAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICB9XG4gICAgfVxufVxuXG5mdW5jdGlvbiBjbG9zZUFsbChjb2RlOiBudW1iZXIsIHJlYXNvbjogc3RyaW5nKTogdm9pZCB7XG4gICAgZm9yIChjb25zdCBjb25uIG9mIFsuLi5HLmNvbm5lY3Rpb25zLnZhbHVlcygpXSkge1xuICAgICAgICBjb25uLl9jbG9zZWQoY29kZSwgcmVhc29uLCBmYWxzZSk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBmYWlsQWxsKHJlYXNvbjogc3RyaW5nKTogdm9pZCB7XG4gICAgZm9yIChjb25zdCBjb25uIG9mIFsuLi5HLmNvbm5lY3Rpb25zLnZhbHVlcygpXSkge1xuICAgICAgICBjb25uLl9mYWlsKG5ldyBFcnJvcihyZWFzb24pKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmV4cG9ydCBpbnRlcmZhY2UgU2l6ZSB7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXG4gICAgV2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXG4gICAgSGVpZ2h0OiBudW1iZXI7XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgUmVjdCB7XG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cbiAgICBYOiBudW1iZXI7XG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cbiAgICBZOiBudW1iZXI7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xuICAgIFdpZHRoOiBudW1iZXI7XG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHJlY3RhbmdsZS4gKi9cbiAgICBIZWlnaHQ6IG51bWJlcjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBTY3JlZW4ge1xuICAgIC8qKiBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi4gKi9cbiAgICBJRDogc3RyaW5nO1xuICAgIC8qKiBIdW1hbi1yZWFkYWJsZSBuYW1lIG9mIHRoZSBzY3JlZW4uICovXG4gICAgTmFtZTogc3RyaW5nO1xuICAgIC8qKiBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLiAqL1xuICAgIFNjYWxlRmFjdG9yOiBudW1iZXI7XG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBYOiBudW1iZXI7XG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBZOiBudW1iZXI7XG4gICAgLyoqIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uICovXG4gICAgU2l6ZTogU2l6ZTtcbiAgICAvKiogQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LiAqL1xuICAgIEJvdW5kczogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0IChiZWZvcmUgc2NhbGluZykuICovXG4gICAgUGh5c2ljYWxCb3VuZHM6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLiAqL1xuICAgIFdvcmtBcmVhOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgV29ya0FyZWEgb2YgdGhlIHNjcmVlbiAoYmVmb3JlIHNjYWxpbmcpLiAqL1xuICAgIFBoeXNpY2FsV29ya0FyZWE6IFJlY3Q7XG4gICAgLyoqIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xuICAgIElzUHJpbWFyeTogYm9vbGVhbjtcbiAgICAvKiogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uICovXG4gICAgUm90YXRpb246IG51bWJlcjtcbn1cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TY3JlZW5zKTtcblxuY29uc3QgZ2V0QWxsID0gMDtcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xuY29uc3QgZ2V0Q3VycmVudCA9IDI7XG5jb25zdCBnZXRCeUlEID0gMztcbmNvbnN0IGdldEJ5SW5kZXggPSA0O1xuXG4vKipcbiAqIEdldHMgYWxsIHNjcmVlbnMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gYXJyYXkgb2YgU2NyZWVuIG9iamVjdHMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKTogUHJvbWlzZTxTY3JlZW5bXT4ge1xuICAgIHJldHVybiBjYWxsKGdldEFsbCk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIHByaW1hcnkgc2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgIHJldHVybiBjYWxsKGdldFByaW1hcnkpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0Q3VycmVudCk7XG59XG5cbi8qKlxuICogR2V0cyBhIHNjcmVlbiBieSBpdHMgdW5pcXVlIGRpc3BsYXkgSUQuXG4gKlxuICogQHBhcmFtIGlkIC0gVGhlIHVuaXF1ZSBpZGVudGlmaWVyIG9mIHRoZSBzY3JlZW4uXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgbWF0Y2hpbmcgU2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QnlJRChpZDogc3RyaW5nKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRCeUlELCB7IGlkIH0pO1xufVxuXG4vKipcbiAqIEdldHMgYSBzY3JlZW4gYnkgaXRzIGluZGV4IGluIHRoZSBzY3JlZW4gbGlzdC5cbiAqXG4gKiBAcGFyYW0gaW5kZXggLSBUaGUgemVyby1iYXNlZCBpbmRleCBvZiB0aGUgc2NyZWVuLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIG1hdGNoaW5nIFNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEJ5SW5kZXgoaW5kZXg6IG51bWJlcik6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0QnlJbmRleCwgeyBpbmRleCB9KTtcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLklPUyk7XG5cbi8vIE1ldGhvZCBJRHNcbmNvbnN0IEhhcHRpY3NJbXBhY3QgPSAwO1xuY29uc3QgRGV2aWNlSW5mbyA9IDE7XG5cbmV4cG9ydCBuYW1lc3BhY2UgSGFwdGljcyB7XG4gICAgZXhwb3J0IHR5cGUgSW1wYWN0U3R5bGUgPSBcImxpZ2h0XCJ8XCJtZWRpdW1cInxcImhlYXZ5XCJ8XCJzb2Z0XCJ8XCJyaWdpZFwiO1xuICAgIGV4cG9ydCBmdW5jdGlvbiBJbXBhY3Qoc3R5bGU6IEltcGFjdFN0eWxlID0gXCJtZWRpdW1cIik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gY2FsbChIYXB0aWNzSW1wYWN0LCB7IHN0eWxlIH0pO1xuICAgIH1cbn1cblxuZXhwb3J0IG5hbWVzcGFjZSBEZXZpY2Uge1xuICAgIGV4cG9ydCBpbnRlcmZhY2UgSW5mbyB7XG4gICAgICAgIG1vZGVsOiBzdHJpbmc7XG4gICAgICAgIHN5c3RlbU5hbWU6IHN0cmluZztcbiAgICAgICAgc3lzdGVtVmVyc2lvbjogc3RyaW5nO1xuICAgICAgICBpc1NpbXVsYXRvcjogYm9vbGVhbjtcbiAgICB9XG4gICAgZXhwb3J0IGZ1bmN0aW9uIEluZm8oKTogUHJvbWlzZTxJbmZvPiB7XG4gICAgICAgIHJldHVybiBjYWxsKERldmljZUluZm8pO1xuICAgIH1cbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkFuZHJvaWQpO1xuXG4vLyBNZXRob2QgSURzIChtdXN0IG1hdGNoIG1lc3NhZ2Vwcm9jZXNzb3JfYW5kcm9pZC5nbylcbmNvbnN0IEhhcHRpY3NWaWJyYXRlID0gMDtcbmNvbnN0IERldmljZUluZm8gPSAxO1xuY29uc3QgVG9hc3RTaG93ID0gMjtcblxuZXhwb3J0IG5hbWVzcGFjZSBIYXB0aWNzIHtcbiAgICAvKiogVmlicmF0ZSB0aGUgZGV2aWNlIGZvciB0aGUgZ2l2ZW4gZHVyYXRpb24gaW4gbWlsbGlzZWNvbmRzLiAqL1xuICAgIGV4cG9ydCBmdW5jdGlvbiBWaWJyYXRlKGR1cmF0aW9uTXM6IG51bWJlciA9IDEwMCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gY2FsbChIYXB0aWNzVmlicmF0ZSwgeyBkdXJhdGlvbjogZHVyYXRpb25NcyB9KTtcbiAgICB9XG59XG5cbmV4cG9ydCBuYW1lc3BhY2UgRGV2aWNlIHtcbiAgICBleHBvcnQgaW50ZXJmYWNlIEluZm8ge1xuICAgICAgICBwbGF0Zm9ybTogc3RyaW5nO1xuICAgICAgICBtYW51ZmFjdHVyZXI6IHN0cmluZztcbiAgICAgICAgYnJhbmQ6IHN0cmluZztcbiAgICAgICAgbW9kZWw6IHN0cmluZztcbiAgICAgICAgZGV2aWNlOiBzdHJpbmc7XG4gICAgICAgIHZlcnNpb246IHN0cmluZztcbiAgICAgICAgc2RrSW50OiBudW1iZXI7XG4gICAgfVxuICAgIC8qKiBSZXR1cm4gaW5mb3JtYXRpb24gYWJvdXQgdGhlIEFuZHJvaWQgZGV2aWNlLiAqL1xuICAgIGV4cG9ydCBmdW5jdGlvbiBJbmZvKCk6IFByb21pc2U8SW5mbz4ge1xuICAgICAgICByZXR1cm4gY2FsbChEZXZpY2VJbmZvKTtcbiAgICB9XG59XG5cbmV4cG9ydCBuYW1lc3BhY2UgVG9hc3Qge1xuICAgIC8qKiBTaG93IGEgc2hvcnQgQW5kcm9pZCB0b2FzdCBtZXNzYWdlLiAqL1xuICAgIGV4cG9ydCBmdW5jdGlvbiBTaG93KG1lc3NhZ2U6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gY2FsbChUb2FzdFNob3csIHsgbWVzc2FnZSB9KTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogVXBkYXRlciBldmVudCBuYW1lIGNvbnN0YW50cy5cbiAqXG4gKiBVc2UgdGhlc2UgaW5zdGVhZCBvZiBoYXJkLWNvZGluZyBzdHJpbmcgbGl0ZXJhbHMgd2hlbiBzdWJzY3JpYmluZyB0b1xuICogdXBkYXRlciBldmVudHMgZnJvbSBKYXZhU2NyaXB0OlxuICpcbiAqICAgICBpbXBvcnQgeyBFdmVudHMsIFVwZGF0ZXIgfSBmcm9tIFwiQHdhaWxzaW8vcnVudGltZVwiO1xuICpcbiAqICAgICBFdmVudHMuT24oVXBkYXRlci5FdmVudHMuVXBkYXRlQXZhaWxhYmxlLCAoZSkgPT4ge1xuICogICAgICAgICBjb25zb2xlLmxvZyhcInVwZGF0ZSBmb3VuZDpcIiwgZS5kYXRhLnZlcnNpb24pO1xuICogICAgIH0pO1xuICpcbiAqICAgICBFdmVudHMuT24oVXBkYXRlci5FdmVudHMuRG93bmxvYWRQcm9ncmVzcywgKGUpID0+IHtcbiAqICAgICAgICAgY29uc3QgcCA9IGUuZGF0YTtcbiAqICAgICAgICAgY29uc29sZS5sb2coYCR7cC53cml0dGVufSAvICR7cC50b3RhbH0gYnl0ZXNgKTtcbiAqICAgICB9KTtcbiAqXG4gKiBNaXJyb3JzIHRoZSBHby1zaWRlIGNvbnN0YW50cyBpbiBgcGtnL3VwZGF0ZXIvZXZlbnRzLmdvYCBhbmQgdGhlXG4gKiB1c2VyLWFjdGlvbiBjb25zdGFudHMgaW4gYHBrZy91cGRhdGVyL3dpbmRvd19saWZlY3ljbGUuZ29gLiBBbnlcbiAqIGNoYW5nZXMgaGVyZSBtdXN0IHN0YXkgaW4gc3luYyB3aXRoIHRob3NlIGZpbGVzIFx1MjAxNCB0aGVyZSdzIGFuXG4gKiBpbnRlZ3JhdGlvbiB0ZXN0IHRoYXQgYXNzZXJ0cyB0aGUgc3RyaW5ncyBtYXRjaC5cbiAqL1xuZXhwb3J0IGNvbnN0IEV2ZW50cyA9IE9iamVjdC5mcmVlemUoe1xuICAgIC8qKiBBIENoZWNrIHJvdW5kLXRyaXAgaXMgc3RhcnRpbmcuIFBheWxvYWQ6IG51bGwuICovXG4gICAgQ2hlY2tTdGFydGVkOiBcIndhaWxzOnVwZGF0ZXI6Y2hlY2stc3RhcnRlZFwiLFxuICAgIC8qKiBDaGVjayBmb3VuZCBhIG5ld2VyIHJlbGVhc2UuIFBheWxvYWQ6IFJlbGVhc2UuICovXG4gICAgVXBkYXRlQXZhaWxhYmxlOiBcIndhaWxzOnVwZGF0ZXI6dXBkYXRlLWF2YWlsYWJsZVwiLFxuICAgIC8qKiBDaGVjayBjb25maXJtZWQgdGhlIGNhbGxlciBpcyB1cCB0byBkYXRlLiBQYXlsb2FkOiBudWxsLiAqL1xuICAgIE5vVXBkYXRlOiBcIndhaWxzOnVwZGF0ZXI6bm8tdXBkYXRlXCIsXG4gICAgLyoqIERvd25sb2FkIGlzIHN0YXJ0aW5nLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xuICAgIERvd25sb2FkU3RhcnRlZDogXCJ3YWlsczp1cGRhdGVyOmRvd25sb2FkLXN0YXJ0ZWRcIixcbiAgICAvKiogUGVyaW9kaWMgcHJvZ3Jlc3MgdGljayBkdXJpbmcgZG93bmxvYWQgKH4xMCBIeikuIFBheWxvYWQ6IFByb2dyZXNzLiAqL1xuICAgIERvd25sb2FkUHJvZ3Jlc3M6IFwid2FpbHM6dXBkYXRlcjpkb3dubG9hZC1wcm9ncmVzc1wiLFxuICAgIC8qKiBBbGwgYnl0ZXMgYXJlIG9uIGRpc2ssIGJ1dCB2ZXJpZmljYXRpb24gaGFzIG5vdCB5ZXQgc3RhcnRlZC4gUGF5bG9hZDogUmVsZWFzZS4gKi9cbiAgICBEb3dubG9hZENvbXBsZXRlOiBcIndhaWxzOnVwZGF0ZXI6ZG93bmxvYWQtY29tcGxldGVcIixcbiAgICAvKiogU2lnbmF0dXJlIC8gZGlnZXN0IHZlcmlmaWNhdGlvbiBoYXMgc3RhcnRlZC4gUGF5bG9hZDogUmVsZWFzZS4gKi9cbiAgICBWZXJpZnlpbmc6IFwid2FpbHM6dXBkYXRlcjp2ZXJpZnlpbmdcIixcbiAgICAvKiogVGhlIFVwZGF0ZXIgaXMgc3dhcHBpbmcgdGhlIGJpbmFyeSBpbnRvIHBsYWNlLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xuICAgIEluc3RhbGxpbmc6IFwid2FpbHM6dXBkYXRlcjppbnN0YWxsaW5nXCIsXG4gICAgLyoqIFVwZGF0ZSBpcyBzdGFnZWQgYW5kIGEgcmVzdGFydCBpcyBwZW5kaW5nLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xuICAgIFVwZGF0ZVJlYWR5OiBcIndhaWxzOnVwZGF0ZXI6dXBkYXRlLXJlYWR5XCIsXG4gICAgLyoqIFNvbWV0aGluZyBmYWlsZWQuIFBheWxvYWQ6IEVycm9ySW5mbyB7IHN0YWdlLCBtZXNzYWdlLCBwcm92aWRlciB9LiAqL1xuICAgIEVycm9yOiBcIndhaWxzOnVwZGF0ZXI6ZXJyb3JcIixcbiAgICAvKiogSG9zdC1zaWRlIGNvbnRleHQgZGVsaXZlcmVkIG9uY2UgcGVyIHNlc3Npb24uIFBheWxvYWQ6IE1ldGEgeyBjdXJyZW50VmVyc2lvbiwgc2tpcHBlZFZlcnNpb24gfS4gKi9cbiAgICBNZXRhOiBcIndhaWxzOnVwZGF0ZXI6bWV0YVwiLFxuXG4gICAgLyoqIFN1Yi1uYW1lc3BhY2U6IHVzZXItYWN0aW9uIGV2ZW50cyB0aGF0IHRoZSBVSSBlbWl0cyBCQUNLIHRvIHRoZSBob3N0LiAqL1xuICAgIFVzZXI6IE9iamVjdC5mcmVlemUoe1xuICAgICAgICAvKiogVXNlciBjbGlja2VkIEluc3RhbGwgb24gYW4gQXZhaWxhYmxlIHVwZGF0ZS4gKi9cbiAgICAgICAgSW5zdGFsbDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6aW5zdGFsbFwiLFxuICAgICAgICAvKiogVXNlciBjbGlja2VkIFJlc3RhcnQgJiBBcHBseSBvbiBhIFJlYWR5IHVwZGF0ZS4gKi9cbiAgICAgICAgUmVzdGFydDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6cmVzdGFydFwiLFxuICAgICAgICAvKiogVXNlciBjbGlja2VkIFNraXAgVGhpcyBWZXJzaW9uLiAqL1xuICAgICAgICBTa2lwOiBcIndhaWxzOnVwZGF0ZXI6dXNlcjpza2lwXCIsXG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgUmVtaW5kIE1lIExhdGVyLiAqL1xuICAgICAgICBSZW1pbmQ6IFwid2FpbHM6dXBkYXRlcjp1c2VyOnJlbWluZFwiLFxuICAgICAgICAvKiogVXNlciBjbGlja2VkIENsb3NlIC8gQ2FuY2VsLiAqL1xuICAgICAgICBDYW5jZWw6IFwid2FpbHM6dXBkYXRlcjp1c2VyOmNhbmNlbFwiLFxuICAgIH0pLFxuXG4gICAgLyoqIFN1Yi1uYW1lc3BhY2U6IGZyYW1ld29yay1pbnRlcm5hbCBldmVudHMgdGhlIFVJIGVtaXRzIHRvIGNvb3JkaW5hdGVcbiAgICAgKiAgd2l0aCB0aGUgaG9zdC4gTW9zdCBhcHAgY29kZSBjYW4gaWdub3JlIHRoZXNlLiAqL1xuICAgIFdpbmRvdzogT2JqZWN0LmZyZWV6ZSh7XG4gICAgICAgIC8qKiBUaGUgd2luZG93IGZpbmlzaGVkIGxvYWRpbmcgYW5kIGFza3MgdGhlIGhvc3QgdG8gcmVwbGF5IGN1cnJlbnQgc3RhdGUuICovXG4gICAgICAgIFJlYWR5OiBcIndhaWxzOnVwZGF0ZXI6d2luZG93OnJlYWR5XCIsXG4gICAgfSksXG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBOzs7QUM2QkEsSUFBTSxjQUNGO0FBU0csU0FBUyxPQUFPLE9BQWUsSUFBWTtBQUM5QyxRQUFNLFNBQVMsT0FBTyxXQUFXLGVBQWUsT0FBTyxPQUFPLG9CQUFvQixhQUM1RSxTQUNBO0FBRU4sTUFBSSxRQUFRO0FBQ1IsVUFBTSxRQUFRLElBQUksV0FBVyxJQUFJO0FBQ2pDLFdBQU8sZ0JBQWdCLEtBQUs7QUFDNUIsUUFBSUEsTUFBSztBQUNULGFBQVNDLEtBQUksR0FBR0EsS0FBSSxNQUFNQSxNQUFLO0FBQzNCLE1BQUFELE9BQU0sWUFBWSxNQUFNQyxFQUFDLElBQUksRUFBRTtBQUFBLElBQ25DO0FBQ0EsV0FBT0Q7QUFBQSxFQUNYO0FBSUEsTUFBSSxLQUFLO0FBQ1QsTUFBSSxJQUFJLE9BQU87QUFDZixTQUFPLEtBQUs7QUFDUixVQUFNLFlBQWEsS0FBSyxPQUFPLElBQUksS0FBTSxDQUFDO0FBQUEsRUFDOUM7QUFDQSxTQUFPO0FBQ1g7OztBQzdDTyxJQUFNLFNBQVMsT0FBTyxXQUFXLGVBQWUsT0FBTyxhQUFhOzs7QUNGM0UsU0FBUyxhQUFxQjtBQUMxQixTQUFPLE9BQU8sU0FBUyxTQUFTO0FBQ3BDO0FBR0EsSUFBTSxrQkFBa0IsTUFBTTtBQWV2QixJQUFNLGVBQU4sY0FBMkIsTUFBTTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU1wQyxZQUFZLFNBQWtCLFNBQXdCO0FBQ2xELFVBQU0sU0FBUyxPQUFPO0FBQ3RCLFNBQUssT0FBTztBQUFBLEVBQ2hCO0FBQ0o7QUFHTyxJQUFNLGNBQWMsT0FBTyxPQUFPO0FBQUEsRUFDckMsTUFBTTtBQUFBLEVBQ04sV0FBVztBQUFBLEVBQ1gsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsYUFBYTtBQUFBLEVBQ2IsUUFBUTtBQUFBLEVBQ1IsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsUUFBUTtBQUFBLEVBQ1IsU0FBUztBQUFBLEVBQ1QsWUFBWTtBQUFBLEVBQ1osS0FBSztBQUFBLEVBQ0wsU0FBUztBQUNiLENBQUM7QUFDTSxJQUFJLFdBQVcsT0FBTztBQXVCN0IsSUFBSSxrQkFBMkM7QUFzQnhDLFNBQVMsYUFBYSxXQUEwQztBQUNuRSxvQkFBa0I7QUFDdEI7QUFLTyxTQUFTLGVBQXdDO0FBQ3BELFNBQU87QUFDWDtBQVNPLFNBQVMsaUJBQWlCLFFBQWdCLGFBQXFCLElBQUk7QUFDdEUsU0FBTyxTQUFVLFFBQWdCLE9BQVksTUFBTTtBQUMvQyxXQUFPLGtCQUFrQixRQUFRLFFBQVEsWUFBWSxJQUFJO0FBQUEsRUFDN0Q7QUFDSjtBQUVBLGVBQWUsa0JBQWtCLFVBQWtCLFFBQWdCLFlBQW9CLE1BQXlCO0FBcEloSCxNQUFBRSxLQUFBO0FBc0lJLE1BQUksaUJBQWlCO0FBQ2pCLFdBQU8sZ0JBQWdCLEtBQUssVUFBVSxRQUFRLFlBQVksSUFBSTtBQUFBLEVBQ2xFO0FBR0EsTUFBSSxNQUFNLElBQUksSUFBSSxXQUFXLENBQUM7QUFFOUIsTUFBSSxPQUF1RDtBQUFBLElBQ3pELFFBQVE7QUFBQSxJQUNSO0FBQUEsRUFDRjtBQUNBLE1BQUksU0FBUyxRQUFRLFNBQVMsUUFBVztBQUN2QyxTQUFLLE9BQU87QUFBQSxFQUNkO0FBRUEsTUFBSSxVQUFrQztBQUFBLElBQ2xDLENBQUMsbUJBQW1CLEdBQUc7QUFBQSxJQUN2QixDQUFDLGNBQWMsR0FBRztBQUFBLEVBQ3RCO0FBQ0EsTUFBSSxZQUFZO0FBQ1osWUFBUSxxQkFBcUIsSUFBSTtBQUFBLEVBQ3JDO0FBRUEsUUFBTSxVQUFVLEtBQUssVUFBVSxJQUFJO0FBQ25DLE1BQUk7QUFDSixNQUFJLFFBQVEsU0FBUyxpQkFBaUI7QUFDbEMsZUFBVyxNQUFNLFlBQVksS0FBSyxTQUFTLE9BQU87QUFBQSxFQUN0RCxPQUFPO0FBQ0gsZUFBVyxNQUFNLE1BQU0sS0FBSyxFQUFFLFFBQVEsUUFBUSxTQUFTLE1BQU0sUUFBUSxDQUFDO0FBQUEsRUFDMUU7QUFDQSxNQUFJLENBQUMsU0FBUyxJQUFJO0FBQ2hCLFVBQU0sS0FBSyxTQUFTLFFBQVEsSUFBSSxjQUFjO0FBQzlDLFFBQUkseUJBQUksU0FBUyxxQkFBcUI7QUFDbEMsWUFBTSxPQUFzQixNQUFNLFNBQVMsS0FBSztBQUNoRCxVQUFJO0FBQ0osY0FBUSxLQUFLLE1BQU07QUFBQSxRQUNmLEtBQUs7QUFBa0IsZ0JBQU0sSUFBSSxlQUFlLEtBQUssT0FBTztBQUFHO0FBQUEsUUFDL0QsS0FBSztBQUFrQixnQkFBTSxJQUFJLFVBQVUsS0FBSyxPQUFPO0FBQUc7QUFBQSxRQUMxRCxLQUFLO0FBQWtCLGdCQUFNLElBQUksYUFBYSxLQUFLLE9BQU87QUFBRztBQUFBLFFBQzdEO0FBQXVCLGdCQUFNLElBQUksTUFBTSxLQUFLLE9BQU87QUFBQSxNQUN2RDtBQUNBLFVBQUksUUFBUSxLQUFLO0FBQ2pCLFlBQU07QUFBQSxJQUNWO0FBQ0EsVUFBTSxJQUFJLE1BQU0sTUFBTSxTQUFTLEtBQUssQ0FBQztBQUFBLEVBQ3ZDO0FBRUEsUUFBSyxNQUFBQSxNQUFBLFNBQVMsUUFBUSxJQUFJLGNBQWMsTUFBbkMsZ0JBQUFBLElBQXNDLFFBQVEsd0JBQTlDLFlBQXFFLFFBQVEsSUFBSTtBQUNsRixXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCLE9BQU87QUFDSCxXQUFPLFNBQVMsS0FBSztBQUFBLEVBQ3pCO0FBQ0o7QUFPQSxlQUFlLFlBQVksS0FBVSxTQUFpQyxTQUFvQztBQUN0RyxRQUFNLFVBQVUsT0FBTztBQUN2QixRQUFNLFlBQVksSUFBSSxZQUFZLEVBQUUsT0FBTyxPQUFPO0FBQ2xELFFBQU0sY0FBYyxLQUFLLEtBQUssVUFBVSxTQUFTLGVBQWU7QUFFaEUsV0FBUyxJQUFJLEdBQUcsSUFBSSxjQUFjLEdBQUcsS0FBSztBQUN0QyxVQUFNLFFBQVEsVUFBVSxTQUFTLElBQUksa0JBQWtCLElBQUksS0FBSyxlQUFlO0FBQy9FLFVBQU0sT0FBTyxNQUFNLE1BQU0sS0FBSztBQUFBLE1BQzFCLFFBQVE7QUFBQSxNQUNSLFNBQVMsaUNBQ0YsVUFERTtBQUFBLFFBRUwsb0JBQW9CO0FBQUEsUUFDcEIsdUJBQXVCLE9BQU8sQ0FBQztBQUFBLFFBQy9CLHVCQUF1QixPQUFPLFdBQVc7QUFBQSxNQUM3QztBQUFBLE1BQ0EsTUFBTTtBQUFBLElBQ1YsQ0FBQztBQUNELFFBQUksQ0FBQyxLQUFLLElBQUk7QUFDVixZQUFNLElBQUksTUFBTSxNQUFNLEtBQUssS0FBSyxDQUFDO0FBQUEsSUFDckM7QUFBQSxFQUNKO0FBRUEsU0FBTyxNQUFNLEtBQUs7QUFBQSxJQUNkLFFBQVE7QUFBQSxJQUNSLFNBQVMsaUNBQ0YsVUFERTtBQUFBLE1BRUwsb0JBQW9CO0FBQUEsTUFDcEIsdUJBQXVCLE9BQU8sY0FBYyxDQUFDO0FBQUEsTUFDN0MsdUJBQXVCLE9BQU8sV0FBVztBQUFBLElBQzdDO0FBQUEsSUFDQSxNQUFNLFVBQVUsVUFBVSxjQUFjLEtBQUssZUFBZTtBQUFBLEVBQ2hFLENBQUM7QUFDTDtBQWpPQTtBQThPQSxJQUFNLGdCQUF3QyxVQUMxQyxTQUFRLFlBQWUsVUFBZixtQkFBc0IsaUJBQWdCLGFBQWMsT0FBZSxRQUFRO0FBRXZGLElBQUksZUFBZTtBQUNmLFFBQU0sVUFBVSxvQkFBSSxJQUE4RTtBQUVsRyxFQUFDLE9BQWUsd0JBQXdCLENBQUMsSUFBWSxVQUF5QixVQUF5QjtBQXBQM0csUUFBQUE7QUFxUFEsVUFBTSxVQUFVLFFBQVEsSUFBSSxFQUFFO0FBQzlCLFFBQUksQ0FBQyxRQUFTO0FBQ2QsWUFBUSxPQUFPLEVBQUU7QUFDakIsUUFBSSxPQUFPO0FBQ1AsY0FBUSxPQUFPLElBQUksTUFBTSxLQUFLLENBQUM7QUFDL0I7QUFBQSxJQUNKO0FBQ0EsUUFBSTtBQUNBLFlBQU0sV0FBVyxLQUFLLE1BQU0sOEJBQVksSUFBSTtBQUM1QyxVQUFJLENBQUMsU0FBUyxJQUFJO0FBQ2QsZ0JBQVEsT0FBTyxJQUFJLE9BQU1BLE1BQUEsU0FBUyxVQUFULE9BQUFBLE1BQWtCLDRCQUE0QixDQUFDO0FBQ3hFO0FBQUEsTUFDSjtBQUNBLGNBQVEsUUFBUSxVQUFVLFdBQVcsU0FBUyxPQUFPLFNBQVMsSUFBSTtBQUFBLElBQ3RFLFNBQVMsR0FBRztBQUNSLGNBQVEsT0FBTyxDQUFDO0FBQUEsSUFDcEI7QUFBQSxFQUNKO0FBRUEsb0JBQWtCO0FBQUEsSUFDZCxLQUFLLFVBQWtCLFFBQWdCLFlBQW9CLE1BQXlCO0FBQ2hGLGFBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLGNBQU0sS0FBSyxPQUFPO0FBQ2xCLGdCQUFRLElBQUksSUFBSSxFQUFFLFNBQVMsT0FBTyxDQUFDO0FBQ25DLFlBQUk7QUFDQSx3QkFBYyxZQUFZLElBQUksS0FBSyxVQUFVO0FBQUEsWUFDekMsUUFBUTtBQUFBLFlBQ1I7QUFBQSxZQUNBO0FBQUEsWUFDQSxNQUFNLHNCQUFRO0FBQUEsWUFDZDtBQUFBLFVBQ0osQ0FBQyxDQUFDO0FBQUEsUUFDTixTQUFTLEdBQUc7QUFFUixrQkFBUSxPQUFPLEVBQUU7QUFDakIsaUJBQU8sQ0FBQztBQUFBLFFBQ1o7QUFBQSxNQUNKLENBQUM7QUFBQSxJQUNMO0FBQUEsRUFDSjtBQUNKOzs7QUhqUkEsSUFBTSxPQUFPLGlCQUFpQixZQUFZLE9BQU87QUFFakQsSUFBTSxpQkFBaUI7QUFPaEIsU0FBUyxRQUFRLEtBQWtDO0FBQ3RELFNBQU8sS0FBSyxnQkFBZ0IsRUFBQyxLQUFLLElBQUksU0FBUyxFQUFDLENBQUM7QUFDckQ7OztBSXZCQTtBQUFBO0FBQUEsZUFBQUM7QUFBQSxFQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQWVBLElBQUksUUFBUTtBQUNSLFNBQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQUN0QztBQUVBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksTUFBTTtBQUdoRCxJQUFNLGFBQWE7QUFDbkIsSUFBTSxnQkFBZ0I7QUFDdEIsSUFBTSxjQUFjO0FBQ3BCLElBQU0saUJBQWlCO0FBQ3ZCLElBQU0saUJBQWlCO0FBQ3ZCLElBQU0saUJBQWlCO0FBMEd2QixTQUFTLE9BQU8sTUFBYyxVQUFnRixDQUFDLEdBQWlCO0FBQzVILFNBQU9BLE1BQUssTUFBTSxPQUFPO0FBQzdCO0FBUU8sU0FBUyxLQUFLLFNBQWdEO0FBQUUsU0FBTyxPQUFPLFlBQVksT0FBTztBQUFHO0FBUXBHLFNBQVMsUUFBUSxTQUFnRDtBQUFFLFNBQU8sT0FBTyxlQUFlLE9BQU87QUFBRztBQVExRyxTQUFTQyxPQUFNLFNBQWdEO0FBQUUsU0FBTyxPQUFPLGFBQWEsT0FBTztBQUFHO0FBUXRHLFNBQVMsU0FBUyxTQUFnRDtBQUFFLFNBQU8sT0FBTyxnQkFBZ0IsT0FBTztBQUFHO0FBVzVHLFNBQVMsU0FBUyxTQUE0RDtBQWxMckYsTUFBQUM7QUFrTHVGLFVBQU9BLE1BQUEsT0FBTyxnQkFBZ0IsT0FBTyxNQUE5QixPQUFBQSxNQUFtQyxDQUFDO0FBQUc7QUFROUgsU0FBUyxTQUFTLFNBQWlEO0FBQUUsU0FBTyxPQUFPLGdCQUFnQixPQUFPO0FBQUc7OztBQzFMcEg7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDYU8sSUFBTSxpQkFBaUIsb0JBQUksSUFBd0I7QUFFbkQsSUFBTSxXQUFOLE1BQWU7QUFBQSxFQUtsQixZQUFZLFdBQW1CLFVBQStCLGNBQXNCO0FBQ2hGLFNBQUssWUFBWTtBQUNqQixTQUFLLFdBQVc7QUFDaEIsU0FBSyxlQUFlLGdCQUFnQjtBQUFBLEVBQ3hDO0FBQUEsRUFFQSxTQUFTLE1BQW9CO0FBQ3pCLFFBQUk7QUFDQSxXQUFLLFNBQVMsSUFBSTtBQUFBLElBQ3RCLFNBQVMsS0FBSztBQUNWLGNBQVEsTUFBTSxHQUFHO0FBQUEsSUFDckI7QUFFQSxRQUFJLEtBQUssaUJBQWlCLEdBQUksUUFBTztBQUNyQyxTQUFLLGdCQUFnQjtBQUNyQixXQUFPLEtBQUssaUJBQWlCO0FBQUEsRUFDakM7QUFDSjtBQUVPLFNBQVMsWUFBWSxVQUEwQjtBQUNsRCxNQUFJLFlBQVksZUFBZSxJQUFJLFNBQVMsU0FBUztBQUNyRCxNQUFJLENBQUMsV0FBVztBQUNaO0FBQUEsRUFDSjtBQUVBLGNBQVksVUFBVSxPQUFPLE9BQUssTUFBTSxRQUFRO0FBQ2hELE1BQUksVUFBVSxXQUFXLEdBQUc7QUFDeEIsbUJBQWUsT0FBTyxTQUFTLFNBQVM7QUFBQSxFQUM1QyxPQUFPO0FBQ0gsbUJBQWUsSUFBSSxTQUFTLFdBQVcsU0FBUztBQUFBLEVBQ3BEO0FBQ0o7OztBQ25EQTtBQUFBO0FBQUE7QUFBQSxlQUFBQztBQUFBLEVBQUE7QUFBQTtBQUFBO0FBQUEsYUFBQUM7QUFBQSxFQUFBO0FBQUE7QUFBQTtBQWFPLFNBQVMsSUFBYSxRQUFnQjtBQUN6QyxTQUFPO0FBQ1g7QUFNTyxTQUFTLFVBQVUsUUFBcUI7QUFDM0MsU0FBUyxVQUFVLE9BQVEsS0FBSztBQUNwQztBQU9PLFNBQVNDLE9BQWUsU0FBbUQ7QUFDOUUsTUFBSSxZQUFZLEtBQUs7QUFDakIsV0FBTyxDQUFDLFdBQVksV0FBVyxPQUFPLENBQUMsSUFBSTtBQUFBLEVBQy9DO0FBRUEsU0FBTyxDQUFDLFdBQVc7QUFDZixRQUFJLFdBQVcsTUFBTTtBQUNqQixhQUFPLENBQUM7QUFBQSxJQUNaO0FBQ0EsYUFBUyxJQUFJLEdBQUcsSUFBSSxPQUFPLFFBQVEsS0FBSztBQUNwQyxhQUFPLENBQUMsSUFBSSxRQUFRLE9BQU8sQ0FBQyxDQUFDO0FBQUEsSUFDakM7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUNKO0FBT08sU0FBU0MsS0FBMEMsS0FBeUIsT0FBMEQ7QUFDekksTUFBSSxVQUFVLEtBQUs7QUFDZixXQUFPLENBQUMsV0FBWSxXQUFXLE9BQU8sQ0FBQyxJQUFJO0FBQUEsRUFDL0M7QUFFQSxTQUFPLENBQUMsV0FBVztBQUNmLFFBQUksV0FBVyxNQUFNO0FBQ2pCLGFBQU8sQ0FBQztBQUFBLElBQ1o7QUFDQSxlQUFXQyxRQUFPLFFBQVE7QUFDdEIsYUFBT0EsSUFBRyxJQUFJLE1BQU0sT0FBT0EsSUFBRyxDQUFDO0FBQUEsSUFDbkM7QUFDQSxXQUFPO0FBQUEsRUFDWDtBQUNKO0FBTU8sU0FBUyxTQUFrQixTQUEwRDtBQUN4RixNQUFJLFlBQVksS0FBSztBQUNqQixXQUFPO0FBQUEsRUFDWDtBQUVBLFNBQU8sQ0FBQyxXQUFZLFdBQVcsT0FBTyxPQUFPLFFBQVEsTUFBTTtBQUMvRDtBQU1PLFNBQVMsT0FBTyxhQUV2QjtBQUNJLE1BQUksU0FBUztBQUNiLGFBQVcsUUFBUSxhQUFhO0FBQzVCLFFBQUksWUFBWSxJQUFJLE1BQU0sS0FBSztBQUMzQixlQUFTO0FBQ1Q7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUNBLE1BQUksUUFBUTtBQUNSLFdBQU87QUFBQSxFQUNYO0FBRUEsU0FBTyxDQUFDLFdBQVc7QUFDZixlQUFXLFFBQVEsYUFBYTtBQUM1QixVQUFJLFFBQVEsUUFBUTtBQUNoQixlQUFPLElBQUksSUFBSSxZQUFZLElBQUksRUFBRSxPQUFPLElBQUksQ0FBQztBQUFBLE1BQ2pEO0FBQUEsSUFDSjtBQUNBLFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFNTyxTQUFTLGFBQWEsUUFBbUI7QUFDNUMsU0FBTyxJQUFJLEtBQUssTUFBTTtBQUMxQjtBQU1PLElBQU0sU0FBK0MsQ0FBQzs7O0FDMUd0RCxJQUFNLFFBQVEsT0FBTyxPQUFPO0FBQUEsRUFDbEMsU0FBUyxPQUFPLE9BQU87QUFBQSxJQUN0Qix1QkFBdUI7QUFBQSxJQUN2QixzQkFBc0I7QUFBQSxJQUN0QixvQkFBb0I7QUFBQSxJQUNwQixrQkFBa0I7QUFBQSxJQUNsQixZQUFZO0FBQUEsSUFDWixvQkFBb0I7QUFBQSxJQUNwQixvQkFBb0I7QUFBQSxJQUNwQiw0QkFBNEI7QUFBQSxJQUM1QixjQUFjO0FBQUEsSUFDZCx1QkFBdUI7QUFBQSxJQUN2QixtQkFBbUI7QUFBQSxJQUNuQixlQUFlO0FBQUEsSUFDZixlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxJQUNqQixrQkFBa0I7QUFBQSxJQUNsQixnQkFBZ0I7QUFBQSxJQUNoQixpQkFBaUI7QUFBQSxJQUNqQixpQkFBaUI7QUFBQSxJQUNqQixnQkFBZ0I7QUFBQSxJQUNoQixlQUFlO0FBQUEsSUFDZixpQkFBaUI7QUFBQSxJQUNqQixrQkFBa0I7QUFBQSxJQUNsQixZQUFZO0FBQUEsSUFDWixnQkFBZ0I7QUFBQSxJQUNoQixlQUFlO0FBQUEsSUFDZixhQUFhO0FBQUEsSUFDYixpQkFBaUI7QUFBQSxJQUNqQixvQkFBb0I7QUFBQSxJQUNwQiwwQkFBMEI7QUFBQSxJQUMxQiwyQkFBMkI7QUFBQSxJQUMzQiwwQkFBMEI7QUFBQSxJQUMxQix3QkFBd0I7QUFBQSxJQUN4QixhQUFhO0FBQUEsSUFDYixlQUFlO0FBQUEsSUFDZixnQkFBZ0I7QUFBQSxJQUNoQixZQUFZO0FBQUEsSUFDWixpQkFBaUI7QUFBQSxJQUNqQixtQkFBbUI7QUFBQSxJQUNuQixvQkFBb0I7QUFBQSxJQUNwQixxQkFBcUI7QUFBQSxJQUNyQixnQkFBZ0I7QUFBQSxJQUNoQixrQkFBa0I7QUFBQSxJQUNsQixnQkFBZ0I7QUFBQSxJQUNoQixrQkFBa0I7QUFBQSxFQUNuQixDQUFDO0FBQUEsRUFDRCxLQUFLLE9BQU8sT0FBTztBQUFBLElBQ2xCLDRCQUE0QjtBQUFBLElBQzVCLHVDQUF1QztBQUFBLElBQ3ZDLHlDQUF5QztBQUFBLElBQ3pDLDBCQUEwQjtBQUFBLElBQzFCLG9DQUFvQztBQUFBLElBQ3BDLHNDQUFzQztBQUFBLElBQ3RDLG9DQUFvQztBQUFBLElBQ3BDLDBDQUEwQztBQUFBLElBQzFDLDJCQUEyQjtBQUFBLElBQzNCLCtCQUErQjtBQUFBLElBQy9CLG9CQUFvQjtBQUFBLElBQ3BCLDRCQUE0QjtBQUFBLElBQzVCLHNCQUFzQjtBQUFBLElBQ3RCLHNCQUFzQjtBQUFBLElBQ3RCLG9CQUFvQjtBQUFBLElBQ3BCLDRCQUE0QjtBQUFBLElBQzVCLDJCQUEyQjtBQUFBLElBQzNCLCtCQUErQjtBQUFBLElBQy9CLDZCQUE2QjtBQUFBLElBQzdCLGdDQUFnQztBQUFBLElBQ2hDLHFCQUFxQjtBQUFBLElBQ3JCLDZCQUE2QjtBQUFBLElBQzdCLHNCQUFzQjtBQUFBLElBQ3RCLDBCQUEwQjtBQUFBLElBQzFCLHVCQUF1QjtBQUFBLElBQ3ZCLHVCQUF1QjtBQUFBLElBQ3ZCLGdCQUFnQjtBQUFBLElBQ2hCLHNCQUFzQjtBQUFBLElBQ3RCLGNBQWM7QUFBQSxJQUNkLG9CQUFvQjtBQUFBLElBQ3BCLG9CQUFvQjtBQUFBLElBQ3BCLHNCQUFzQjtBQUFBLElBQ3RCLGFBQWE7QUFBQSxJQUNiLGNBQWM7QUFBQSxJQUNkLG1CQUFtQjtBQUFBLElBQ25CLG1CQUFtQjtBQUFBLElBQ25CLHlCQUF5QjtBQUFBLElBQ3pCLGVBQWU7QUFBQSxJQUNmLGlCQUFpQjtBQUFBLElBQ2pCLHVCQUF1QjtBQUFBLElBQ3ZCLHFCQUFxQjtBQUFBLElBQ3JCLHFCQUFxQjtBQUFBLElBQ3JCLHVCQUF1QjtBQUFBLElBQ3ZCLGNBQWM7QUFBQSxJQUNkLGVBQWU7QUFBQSxJQUNmLG9CQUFvQjtBQUFBLElBQ3BCLG9CQUFvQjtBQUFBLElBQ3BCLDBCQUEwQjtBQUFBLElBQzFCLGdCQUFnQjtBQUFBLElBQ2hCLDRCQUE0QjtBQUFBLElBQzVCLDRCQUE0QjtBQUFBLElBQzVCLHlEQUF5RDtBQUFBLElBQ3pELHNDQUFzQztBQUFBLElBQ3RDLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLHFCQUFxQjtBQUFBLElBQ3JCLHNCQUFzQjtBQUFBLElBQ3RCLGdDQUFnQztBQUFBLElBQ2hDLGtDQUFrQztBQUFBLElBQ2xDLG1DQUFtQztBQUFBLElBQ25DLG9DQUFvQztBQUFBLElBQ3BDLCtCQUErQjtBQUFBLElBQy9CLDZCQUE2QjtBQUFBLElBQzdCLHVCQUF1QjtBQUFBLElBQ3ZCLGlDQUFpQztBQUFBLElBQ2pDLDhCQUE4QjtBQUFBLElBQzlCLDRCQUE0QjtBQUFBLElBQzVCLHNDQUFzQztBQUFBLElBQ3RDLDRCQUE0QjtBQUFBLElBQzVCLHNCQUFzQjtBQUFBLElBQ3RCLGtDQUFrQztBQUFBLElBQ2xDLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLHdCQUF3QjtBQUFBLElBQ3hCLG1CQUFtQjtBQUFBLElBQ25CLDBCQUEwQjtBQUFBLElBQzFCLDhCQUE4QjtBQUFBLElBQzlCLHlCQUF5QjtBQUFBLElBQ3pCLDZCQUE2QjtBQUFBLElBQzdCLGlCQUFpQjtBQUFBLElBQ2pCLGdCQUFnQjtBQUFBLElBQ2hCLHNCQUFzQjtBQUFBLElBQ3RCLGVBQWU7QUFBQSxJQUNmLHlCQUF5QjtBQUFBLElBQ3pCLHdCQUF3QjtBQUFBLElBQ3hCLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLGlCQUFpQjtBQUFBLElBQ2pCLGlCQUFpQjtBQUFBLElBQ2pCLHNCQUFzQjtBQUFBLElBQ3RCLG1DQUFtQztBQUFBLElBQ25DLHFDQUFxQztBQUFBLElBQ3JDLHVCQUF1QjtBQUFBLElBQ3ZCLHNCQUFzQjtBQUFBLElBQ3RCLHdCQUF3QjtBQUFBLElBQ3hCLGVBQWU7QUFBQSxJQUNmLDJCQUEyQjtBQUFBLElBQzNCLDBCQUEwQjtBQUFBLElBQzFCLDZCQUE2QjtBQUFBLElBQzdCLFlBQVk7QUFBQSxJQUNaLGdCQUFnQjtBQUFBLElBQ2hCLGtCQUFrQjtBQUFBLElBQ2xCLGdCQUFnQjtBQUFBLElBQ2hCLGtCQUFrQjtBQUFBLElBQ2xCLG1CQUFtQjtBQUFBLElBQ25CLFlBQVk7QUFBQSxJQUNaLHFCQUFxQjtBQUFBLElBQ3JCLHNCQUFzQjtBQUFBLElBQ3RCLHNCQUFzQjtBQUFBLElBQ3RCLDhCQUE4QjtBQUFBLElBQzlCLGlCQUFpQjtBQUFBLElBQ2pCLHlCQUF5QjtBQUFBLElBQ3pCLDJCQUEyQjtBQUFBLElBQzNCLCtCQUErQjtBQUFBLElBQy9CLDBCQUEwQjtBQUFBLElBQzFCLDhCQUE4QjtBQUFBLElBQzlCLGlCQUFpQjtBQUFBLElBQ2pCLHVCQUF1QjtBQUFBLElBQ3ZCLGdCQUFnQjtBQUFBLElBQ2hCLDBCQUEwQjtBQUFBLElBQzFCLHlCQUF5QjtBQUFBLElBQ3pCLHNCQUFzQjtBQUFBLElBQ3RCLGtCQUFrQjtBQUFBLElBQ2xCLG1CQUFtQjtBQUFBLElBQ25CLGtCQUFrQjtBQUFBLElBQ2xCLHVCQUF1QjtBQUFBLElBQ3ZCLG9DQUFvQztBQUFBLElBQ3BDLHNDQUFzQztBQUFBLElBQ3RDLHdCQUF3QjtBQUFBLElBQ3hCLHVCQUF1QjtBQUFBLElBQ3ZCLHlCQUF5QjtBQUFBLElBQ3pCLDRCQUE0QjtBQUFBLElBQzVCLDRCQUE0QjtBQUFBLElBQzVCLGNBQWM7QUFBQSxJQUNkLGVBQWU7QUFBQSxJQUNmLGlCQUFpQjtBQUFBLElBQ2pCLHNDQUFzQztBQUFBLEVBQ3ZDLENBQUM7QUFBQSxFQUNELE9BQU8sT0FBTyxPQUFPO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsZUFBZTtBQUFBLElBQ2Ysb0JBQW9CO0FBQUEsSUFDcEIsaUJBQWlCO0FBQUEsSUFDakIsbUJBQW1CO0FBQUEsSUFDbkIsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsSUFDakIsZUFBZTtBQUFBLElBQ2YsZ0JBQWdCO0FBQUEsSUFDaEIsbUJBQW1CO0FBQUEsSUFDbkIsc0JBQXNCO0FBQUEsSUFDdEIscUJBQXFCO0FBQUEsSUFDckIsb0JBQW9CO0FBQUEsRUFDckIsQ0FBQztBQUFBLEVBQ0QsS0FBSyxPQUFPLE9BQU87QUFBQSxJQUNsQiw0QkFBNEI7QUFBQSxJQUM1QiwrQkFBK0I7QUFBQSxJQUMvQiwrQkFBK0I7QUFBQSxJQUMvQixvQ0FBb0M7QUFBQSxJQUNwQyxnQ0FBZ0M7QUFBQSxJQUNoQyw2QkFBNkI7QUFBQSxJQUM3QiwwQkFBMEI7QUFBQSxJQUMxQixlQUFlO0FBQUEsSUFDZixrQkFBa0I7QUFBQSxJQUNsQixpQkFBaUI7QUFBQSxJQUNqQixxQkFBcUI7QUFBQSxJQUNyQixvQkFBb0I7QUFBQSxJQUNwQiw2QkFBNkI7QUFBQSxJQUM3QiwwQkFBMEI7QUFBQSxJQUMxQixrQkFBa0I7QUFBQSxJQUNsQixrQkFBa0I7QUFBQSxJQUNsQixrQkFBa0I7QUFBQSxJQUNsQixzQkFBc0I7QUFBQSxJQUN0QiwyQkFBMkI7QUFBQSxJQUMzQiw0QkFBNEI7QUFBQSxJQUM1QiwwQkFBMEI7QUFBQSxJQUMxQix3Q0FBd0M7QUFBQSxJQUN4QyxnQkFBZ0I7QUFBQSxJQUNoQixnQkFBZ0I7QUFBQSxJQUNoQixjQUFjO0FBQUEsSUFDZCxjQUFjO0FBQUEsSUFDZCxnQkFBZ0I7QUFBQSxFQUNqQixDQUFDO0FBQUEsRUFDRCxTQUFTLE9BQU8sT0FBTztBQUFBLElBQ3RCLGlCQUFpQjtBQUFBLElBQ2pCLGlCQUFpQjtBQUFBLElBQ2pCLGlCQUFpQjtBQUFBLElBQ2pCLGdCQUFnQjtBQUFBLElBQ2hCLGlCQUFpQjtBQUFBLElBQ2pCLG1CQUFtQjtBQUFBLElBQ25CLHNCQUFzQjtBQUFBLElBQ3RCLG9CQUFvQjtBQUFBLElBQ3BCLHFCQUFxQjtBQUFBLElBQ3JCLGdCQUFnQjtBQUFBLElBQ2hCLGdCQUFnQjtBQUFBLElBQ2hCLGNBQWM7QUFBQSxJQUNkLGNBQWM7QUFBQSxJQUNkLGdCQUFnQjtBQUFBLEVBQ2pCLENBQUM7QUFBQSxFQUNELFFBQVEsT0FBTyxPQUFPO0FBQUEsSUFDckIsMkJBQTJCO0FBQUEsSUFDM0Isb0JBQW9CO0FBQUEsSUFDcEIsNEJBQTRCO0FBQUEsSUFDNUIsY0FBYztBQUFBLElBQ2QsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsSUFDakIsZUFBZTtBQUFBLElBQ2YsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsSUFDakIsa0JBQWtCO0FBQUEsSUFDbEIsb0JBQW9CO0FBQUEsSUFDcEIsYUFBYTtBQUFBLElBQ2Isa0JBQWtCO0FBQUEsSUFDbEIsWUFBWTtBQUFBLElBQ1osaUJBQWlCO0FBQUEsSUFDakIsZ0JBQWdCO0FBQUEsSUFDaEIsZ0JBQWdCO0FBQUEsSUFDaEIsdUJBQXVCO0FBQUEsSUFDdkIsZUFBZTtBQUFBLElBQ2Ysb0JBQW9CO0FBQUEsSUFDcEIsWUFBWTtBQUFBLElBQ1osb0JBQW9CO0FBQUEsSUFDcEIsa0JBQWtCO0FBQUEsSUFDbEIsa0JBQWtCO0FBQUEsSUFDbEIsWUFBWTtBQUFBLElBQ1osY0FBYztBQUFBLElBQ2QsZUFBZTtBQUFBLElBQ2YsaUJBQWlCO0FBQUEsSUFDakIsZ0JBQWdCO0FBQUEsSUFDaEIsZ0JBQWdCO0FBQUEsSUFDaEIsY0FBYztBQUFBLElBQ2QsZ0JBQWdCO0FBQUEsSUFDaEIsV0FBVztBQUFBLEVBQ1osQ0FBQztBQUNGLENBQUM7OztBSHBSRCxJQUFJLFFBQVE7QUFDUixTQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFDbEMsU0FBTyxPQUFPLHFCQUFxQjtBQUN2QztBQUVBLElBQU1DLFFBQU8saUJBQWlCLFlBQVksTUFBTTtBQUNoRCxJQUFNLGFBQWE7QUFvQ1osSUFBTSxhQUFOLE1BQTREO0FBQUEsRUFtQi9ELFlBQVksTUFBUyxNQUFZO0FBQzdCLFNBQUssT0FBTztBQUNaLFNBQUssT0FBTyxzQkFBUTtBQUFBLEVBQ3hCO0FBQ0o7QUFFQSxTQUFTLG1CQUFtQixPQUFZO0FBQ3BDLE1BQUksWUFBWSxlQUFlLElBQUksTUFBTSxJQUFJO0FBQzdDLE1BQUksQ0FBQyxXQUFXO0FBQ1o7QUFBQSxFQUNKO0FBRUEsTUFBSSxhQUFhLElBQUk7QUFBQSxJQUNqQixNQUFNO0FBQUEsSUFDTCxNQUFNLFFBQVEsU0FBVSxPQUFPLE1BQU0sSUFBSSxFQUFFLE1BQU0sSUFBSSxJQUFJLE1BQU07QUFBQSxFQUNwRTtBQUNBLE1BQUksWUFBWSxPQUFPO0FBQ25CLGVBQVcsU0FBUyxNQUFNO0FBQUEsRUFDOUI7QUFVQSxRQUFNLFVBQVUsb0JBQUksSUFBYztBQUNsQyxhQUFXLFlBQVksVUFBVSxNQUFNLEdBQUc7QUFDdEMsUUFBSSxTQUFTLFNBQVMsVUFBVSxHQUFHO0FBQy9CLGNBQVEsSUFBSSxRQUFRO0FBQUEsSUFDeEI7QUFBQSxFQUNKO0FBQ0EsTUFBSSxRQUFRLE9BQU8sR0FBRztBQUNsQixVQUFNLE9BQU8sZUFBZSxJQUFJLE1BQU0sSUFBSTtBQUMxQyxRQUFJLE1BQU07QUFDTixZQUFNLFlBQVksS0FBSyxPQUFPLE9BQUssQ0FBQyxRQUFRLElBQUksQ0FBQyxDQUFDO0FBQ2xELFVBQUksVUFBVSxXQUFXLEdBQUc7QUFDeEIsdUJBQWUsT0FBTyxNQUFNLElBQUk7QUFBQSxNQUNwQyxPQUFPO0FBQ0gsdUJBQWUsSUFBSSxNQUFNLE1BQU0sU0FBUztBQUFBLE1BQzVDO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFDSjtBQVVPLFNBQVMsV0FBc0QsV0FBYyxVQUFpQyxjQUFzQjtBQUN2SSxNQUFJLFlBQVksZUFBZSxJQUFJLFNBQVMsS0FBSyxDQUFDO0FBQ2xELFFBQU0sZUFBZSxJQUFJLFNBQVMsV0FBVyxVQUFVLFlBQVk7QUFDbkUsWUFBVSxLQUFLLFlBQVk7QUFDM0IsaUJBQWUsSUFBSSxXQUFXLFNBQVM7QUFDdkMsU0FBTyxNQUFNLFlBQVksWUFBWTtBQUN6QztBQVNPLFNBQVMsR0FBOEMsV0FBYyxVQUE2QztBQUNySCxTQUFPLFdBQVcsV0FBVyxVQUFVLEVBQUU7QUFDN0M7QUFTTyxTQUFTLEtBQWdELFdBQWMsVUFBNkM7QUFDdkgsU0FBTyxXQUFXLFdBQVcsVUFBVSxDQUFDO0FBQzVDO0FBT08sU0FBUyxPQUFPLFlBQXlEO0FBQzVFLGFBQVcsUUFBUSxlQUFhLGVBQWUsT0FBTyxTQUFTLENBQUM7QUFDcEU7QUFLTyxTQUFTLFNBQWU7QUFDM0IsaUJBQWUsTUFBTTtBQUN6QjtBQVdPLFNBQVMsS0FBZ0QsTUFBeUIsTUFBOEI7QUFDbkgsU0FBT0EsTUFBSyxZQUFhLElBQUksV0FBVyxNQUFNLElBQUksQ0FBQztBQUN2RDs7O0FJaExPLFNBQVMsU0FBUyxTQUFjO0FBRW5DLFVBQVE7QUFBQSxJQUNKLGtCQUFrQixVQUFVO0FBQUEsSUFDNUI7QUFBQSxJQUNBO0FBQUEsRUFDSjtBQUNKO0FBTU8sU0FBUyxrQkFBMkI7QUFDdkMsU0FBUSxJQUFJLFdBQVcsV0FBVyxFQUFHLFlBQVk7QUFDckQ7QUFNTyxTQUFTLG9CQUFvQjtBQUNoQyxNQUFJLENBQUMsZUFBZSxDQUFDLGVBQWUsQ0FBQztBQUNqQyxXQUFPO0FBRVgsTUFBSSxTQUFTO0FBRWIsUUFBTSxTQUFTLElBQUksWUFBWTtBQUMvQixRQUFNLGFBQWEsSUFBSSxnQkFBZ0I7QUFDdkMsU0FBTyxpQkFBaUIsUUFBUSxNQUFNO0FBQUUsYUFBUztBQUFBLEVBQU8sR0FBRyxFQUFFLFFBQVEsV0FBVyxPQUFPLENBQUM7QUFDeEYsYUFBVyxNQUFNO0FBQ2pCLFNBQU8sY0FBYyxJQUFJLFlBQVksTUFBTSxDQUFDO0FBRTVDLFNBQU87QUFDWDtBQUtPLFNBQVMsWUFBWSxPQUEyQjtBQXREdkQsTUFBQUM7QUF1REksTUFBSSxNQUFNLGtCQUFrQixhQUFhO0FBQ3JDLFdBQU8sTUFBTTtBQUFBLEVBQ2pCLFdBQVcsRUFBRSxNQUFNLGtCQUFrQixnQkFBZ0IsTUFBTSxrQkFBa0IsTUFBTTtBQUMvRSxZQUFPQSxNQUFBLE1BQU0sT0FBTyxrQkFBYixPQUFBQSxNQUE4QixTQUFTO0FBQUEsRUFDbEQsT0FBTztBQUNILFdBQU8sU0FBUztBQUFBLEVBQ3BCO0FBQ0o7QUFpQ0EsSUFBSSxVQUFVO0FBR2QsSUFBSSxRQUFRO0FBQ1IsV0FBUyxpQkFBaUIsb0JBQW9CLE1BQU07QUFBRSxjQUFVO0FBQUEsRUFBSyxDQUFDO0FBQzFFO0FBRU8sU0FBUyxVQUFVLFVBQXNCO0FBQzVDLE1BQUksV0FBVyxTQUFTLGVBQWUsWUFBWTtBQUMvQyxhQUFTO0FBQUEsRUFDYixPQUFPO0FBQ0gsYUFBUyxpQkFBaUIsb0JBQW9CLFFBQVE7QUFBQSxFQUMxRDtBQUNKOzs7QUM5RkEsSUFBTSx3QkFBd0I7QUFDOUIsSUFBTSwyQkFBMkI7QUFDakMsSUFBSSxvQkFBb0M7QUFFeEMsSUFBTSxpQkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGNBQW9DO0FBQzFDLElBQU0sK0JBQW9DO0FBQzFDLElBQU0sOEJBQW9DO0FBQzFDLElBQU0sY0FBb0M7QUFDMUMsSUFBTSxvQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxrQkFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0sa0JBQW9DO0FBQzFDLElBQU0scUJBQW9DO0FBQzFDLElBQU0sb0JBQW9DO0FBQzFDLElBQU0sb0JBQW9DO0FBQzFDLElBQU0saUJBQW9DO0FBQzFDLElBQU0saUJBQW9DO0FBQzFDLElBQU0sYUFBb0M7QUFDMUMsSUFBTSxxQkFBb0M7QUFDMUMsSUFBTSx5QkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGtCQUFvQztBQUMxQyxJQUFNLGdCQUFvQztBQUMxQyxJQUFNLG9CQUFvQztBQUMxQyxJQUFNLHVCQUFvQztBQUMxQyxJQUFNLDRCQUFvQztBQUMxQyxJQUFNLHFCQUFvQztBQUMxQyxJQUFNLG1DQUFvQztBQUMxQyxJQUFNLG1CQUFvQztBQUMxQyxJQUFNLG1CQUFvQztBQUMxQyxJQUFNLDRCQUFvQztBQUMxQyxJQUFNLHFCQUFvQztBQUMxQyxJQUFNLGdCQUFvQztBQUMxQyxJQUFNLGlCQUFvQztBQUMxQyxJQUFNLGdCQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0sYUFBb0M7QUFDMUMsSUFBTSx5QkFBb0M7QUFDMUMsSUFBTSx1QkFBb0M7QUFDMUMsSUFBTSx3QkFBb0M7QUFDMUMsSUFBTSxxQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxjQUFvQztBQUMxQyxJQUFNLGFBQW9DO0FBQzFDLElBQU0sZUFBb0M7QUFDMUMsSUFBTSxnQkFBb0M7QUFDMUMsSUFBTSxrQkFBb0M7QUFDMUMsSUFBTSxtQkFBb0M7QUFDMUMsSUFBTSxlQUFvQztBQUMxQyxJQUFNLGNBQW9DO0FBQzFDLElBQU0sa0JBQW9DO0FBSzFDLFNBQVMscUJBQXFCLFNBQXlDO0FBQ25FLE1BQUksQ0FBQyxTQUFTO0FBQ1YsV0FBTztBQUFBLEVBQ1g7QUFDQSxTQUFPLFFBQVEsUUFBUSxJQUFJLDhCQUFxQixJQUFHO0FBQ3ZEO0FBTUEsU0FBUyxzQkFBK0I7QUF0RnhDLE1BQUFDLEtBQUE7QUF3RkksUUFBSyxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsWUFBdkIsbUJBQWdDLHFDQUFvQyxNQUFNO0FBQzNFLFdBQU87QUFBQSxFQUNYO0FBR0EsV0FBUSxrQkFBZSxXQUFmLG1CQUF1QixVQUF2QixtQkFBOEIsb0JBQW1CO0FBQzdEO0FBS0EsU0FBUyxpQkFBaUIsR0FBVyxHQUFXLE9BQXFCO0FBbkdyRSxNQUFBQSxLQUFBO0FBb0dJLE9BQUssTUFBQUEsTUFBQSxPQUFlLFdBQWYsZ0JBQUFBLElBQXVCLFlBQXZCLG1CQUFnQyxrQ0FBa0M7QUFDbkUsSUFBQyxPQUFlLE9BQU8sUUFBUSxpQ0FBaUMsYUFBYSxVQUFDLEtBQUksV0FBSyxLQUFLO0FBQUEsRUFDaEc7QUFDSjtBQUdBLElBQUksbUJBQW1CO0FBTXZCLFNBQVMsb0JBQTBCO0FBQy9CLHFCQUFtQjtBQUNuQixNQUFJLG1CQUFtQjtBQUNuQixzQkFBa0IsVUFBVSxPQUFPLHdCQUF3QjtBQUMzRCx3QkFBb0I7QUFBQSxFQUN4QjtBQUNKO0FBS0EsU0FBUyxrQkFBd0I7QUEzSGpDLE1BQUFBLEtBQUE7QUE2SEksUUFBSyxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsVUFBdkIsbUJBQThCLG9CQUFtQixPQUFPO0FBQ3pEO0FBQUEsRUFDSjtBQUNBLHFCQUFtQjtBQUN2QjtBQUtBLFNBQVMsa0JBQXdCO0FBQzdCLG9CQUFrQjtBQUN0QjtBQU9BLFNBQVMsZUFBZSxHQUFXLEdBQWlCO0FBL0lwRCxNQUFBQSxLQUFBO0FBZ0pJLE1BQUksQ0FBQyxpQkFBa0I7QUFHdkIsUUFBSyxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsVUFBdkIsbUJBQThCLG9CQUFtQixPQUFPO0FBQ3pEO0FBQUEsRUFDSjtBQUVBLFFBQU0sZ0JBQWdCLFNBQVMsaUJBQWlCLEdBQUcsQ0FBQztBQUNwRCxRQUFNLGFBQWEscUJBQXFCLGFBQWE7QUFFckQsTUFBSSxxQkFBcUIsc0JBQXNCLFlBQVk7QUFDdkQsc0JBQWtCLFVBQVUsT0FBTyx3QkFBd0I7QUFBQSxFQUMvRDtBQUVBLE1BQUksWUFBWTtBQUNaLGVBQVcsVUFBVSxJQUFJLHdCQUF3QjtBQUNqRCx3QkFBb0I7QUFBQSxFQUN4QixPQUFPO0FBQ0gsd0JBQW9CO0FBQUEsRUFDeEI7QUFDSjtBQTRCQSxJQUFNLFlBQVksdUJBQU8sUUFBUTtBQUlwQjtBQUZiLElBQU0sVUFBTixNQUFNLFFBQU87QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVVULFlBQVksT0FBZSxJQUFJO0FBQzNCLFNBQUssU0FBUyxJQUFJLGlCQUFpQixZQUFZLFFBQVEsSUFBSTtBQUczRCxlQUFXLFVBQVUsT0FBTyxvQkFBb0IsUUFBTyxTQUFTLEdBQUc7QUFDL0QsVUFDSSxXQUFXLGlCQUNSLE9BQVEsS0FBYSxNQUFNLE1BQU0sWUFDdEM7QUFDRSxRQUFDLEtBQWEsTUFBTSxJQUFLLEtBQWEsTUFBTSxFQUFFLEtBQUssSUFBSTtBQUFBLE1BQzNEO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLElBQUksTUFBc0I7QUFDdEIsV0FBTyxJQUFJLFFBQU8sSUFBSTtBQUFBLEVBQzFCO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsV0FBOEI7QUFDMUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxjQUFjO0FBQUEsRUFDekM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFNBQXdCO0FBQ3BCLFdBQU8sS0FBSyxTQUFTLEVBQUUsWUFBWTtBQUFBLEVBQ3ZDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxRQUF1QjtBQUNuQixXQUFPLEtBQUssU0FBUyxFQUFFLFdBQVc7QUFBQSxFQUN0QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EseUJBQXdDO0FBQ3BDLFdBQU8sS0FBSyxTQUFTLEVBQUUsNEJBQTRCO0FBQUEsRUFDdkQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLHdCQUF1QztBQUNuQyxXQUFPLEtBQUssU0FBUyxFQUFFLDJCQUEyQjtBQUFBLEVBQ3REO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxRQUF1QjtBQUNuQixXQUFPLEtBQUssU0FBUyxFQUFFLFdBQVc7QUFBQSxFQUN0QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsY0FBNkI7QUFDekIsV0FBTyxLQUFLLFNBQVMsRUFBRSxpQkFBaUI7QUFBQSxFQUM1QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsYUFBNEI7QUFDeEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxnQkFBZ0I7QUFBQSxFQUMzQztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFlBQTZCO0FBQ3pCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZUFBZTtBQUFBLEVBQzFDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsVUFBMkI7QUFDdkIsV0FBTyxLQUFLLFNBQVMsRUFBRSxhQUFhO0FBQUEsRUFDeEM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxTQUEwQjtBQUN0QixXQUFPLEtBQUssU0FBUyxFQUFFLFlBQVk7QUFBQSxFQUN2QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxZQUE4QjtBQUMxQixXQUFPLEtBQUssU0FBUyxFQUFFLGVBQWU7QUFBQSxFQUMxQztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGVBQWlDO0FBQzdCLFdBQU8sS0FBSyxTQUFTLEVBQUUsa0JBQWtCO0FBQUEsRUFDN0M7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxjQUFnQztBQUM1QixXQUFPLEtBQUssU0FBUyxFQUFFLGlCQUFpQjtBQUFBLEVBQzVDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsY0FBZ0M7QUFDNUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxpQkFBaUI7QUFBQSxFQUM1QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsV0FBMEI7QUFDdEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxjQUFjO0FBQUEsRUFDekM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFdBQTBCO0FBQ3RCLFdBQU8sS0FBSyxTQUFTLEVBQUUsY0FBYztBQUFBLEVBQ3pDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsT0FBd0I7QUFDcEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGVBQThCO0FBQzFCLFdBQU8sS0FBSyxTQUFTLEVBQUUsa0JBQWtCO0FBQUEsRUFDN0M7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxtQkFBc0M7QUFDbEMsV0FBTyxLQUFLLFNBQVMsRUFBRSxzQkFBc0I7QUFBQSxFQUNqRDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsU0FBd0I7QUFDcEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxZQUFZO0FBQUEsRUFDdkM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxZQUE4QjtBQUMxQixXQUFPLEtBQUssU0FBUyxFQUFFLGVBQWU7QUFBQSxFQUMxQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsVUFBeUI7QUFDckIsV0FBTyxLQUFLLFNBQVMsRUFBRSxhQUFhO0FBQUEsRUFDeEM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVFBLFlBQVksR0FBVyxHQUEwQjtBQUM3QyxXQUFPLEtBQUssU0FBUyxFQUFFLG1CQUFtQixFQUFFLEdBQUcsRUFBRSxDQUFDO0FBQUEsRUFDdEQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxlQUFlLGFBQXFDO0FBQ2hELFdBQU8sS0FBSyxTQUFTLEVBQUUsc0JBQXNCLEVBQUUsWUFBWSxDQUFDO0FBQUEsRUFDaEU7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFVQSxvQkFBb0IsR0FBVyxHQUFXLEdBQVcsR0FBMEI7QUFDM0UsV0FBTyxLQUFLLFNBQVMsRUFBRSwyQkFBMkIsRUFBRSxHQUFHLEdBQUcsR0FBRyxFQUFFLENBQUM7QUFBQSxFQUNwRTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLGFBQWEsV0FBbUM7QUFDNUMsV0FBTyxLQUFLLFNBQVMsRUFBRSxvQkFBb0IsRUFBRSxVQUFVLENBQUM7QUFBQSxFQUM1RDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLDJCQUEyQixTQUFpQztBQUN4RCxXQUFPLEtBQUssU0FBUyxFQUFFLGtDQUFrQyxFQUFFLFFBQVEsQ0FBQztBQUFBLEVBQ3hFO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxXQUFXLE9BQWUsUUFBK0I7QUFDckQsV0FBTyxLQUFLLFNBQVMsRUFBRSxrQkFBa0IsRUFBRSxPQUFPLE9BQU8sQ0FBQztBQUFBLEVBQzlEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxXQUFXLE9BQWUsUUFBK0I7QUFDckQsV0FBTyxLQUFLLFNBQVMsRUFBRSxrQkFBa0IsRUFBRSxPQUFPLE9BQU8sQ0FBQztBQUFBLEVBQzlEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxvQkFBb0IsR0FBVyxHQUEwQjtBQUNyRCxXQUFPLEtBQUssU0FBUyxFQUFFLDJCQUEyQixFQUFFLEdBQUcsRUFBRSxDQUFDO0FBQUEsRUFDOUQ7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxhQUFhQyxZQUFtQztBQUM1QyxXQUFPLEtBQUssU0FBUyxFQUFFLG9CQUFvQixFQUFFLFdBQUFBLFdBQVUsQ0FBQztBQUFBLEVBQzVEO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFRQSxRQUFRLE9BQWUsUUFBK0I7QUFDbEQsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlLEVBQUUsT0FBTyxPQUFPLENBQUM7QUFBQSxFQUMzRDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFNBQVMsT0FBOEI7QUFDbkMsV0FBTyxLQUFLLFNBQVMsRUFBRSxnQkFBZ0IsRUFBRSxNQUFNLENBQUM7QUFBQSxFQUNwRDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFFBQVEsTUFBNkI7QUFDakMsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlLEVBQUUsS0FBSyxDQUFDO0FBQUEsRUFDbEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLE9BQXNCO0FBQ2xCLFdBQU8sS0FBSyxTQUFTLEVBQUUsVUFBVTtBQUFBLEVBQ3JDO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBT0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLG1CQUFrQztBQUM5QixXQUFPLEtBQUssU0FBUyxFQUFFLHNCQUFzQjtBQUFBLEVBQ2pEO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxpQkFBZ0M7QUFDNUIsV0FBTyxLQUFLLFNBQVMsRUFBRSxvQkFBb0I7QUFBQSxFQUMvQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0Esa0JBQWlDO0FBQzdCLFdBQU8sS0FBSyxTQUFTLEVBQUUscUJBQXFCO0FBQUEsRUFDaEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGVBQThCO0FBQzFCLFdBQU8sS0FBSyxTQUFTLEVBQUUsa0JBQWtCO0FBQUEsRUFDN0M7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGFBQTRCO0FBQ3hCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZ0JBQWdCO0FBQUEsRUFDM0M7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLGFBQTRCO0FBQ3hCLFdBQU8sS0FBSyxTQUFTLEVBQUUsZ0JBQWdCO0FBQUEsRUFDM0M7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFPQSxRQUF5QjtBQUNyQixXQUFPLEtBQUssU0FBUyxFQUFFLFdBQVc7QUFBQSxFQUN0QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsT0FBc0I7QUFDbEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxVQUFVO0FBQUEsRUFDckM7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFNBQXdCO0FBQ3BCLFdBQU8sS0FBSyxTQUFTLEVBQUUsWUFBWTtBQUFBLEVBQ3ZDO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFLQSxVQUF5QjtBQUNyQixXQUFPLEtBQUssU0FBUyxFQUFFLGFBQWE7QUFBQSxFQUN4QztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsWUFBMkI7QUFDdkIsV0FBTyxLQUFLLFNBQVMsRUFBRSxlQUFlO0FBQUEsRUFDMUM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFVQSx1QkFBdUIsV0FBcUIsR0FBVyxHQUFpQjtBQTduQjVFLFFBQUFDLEtBQUE7QUErbkJRLFVBQUssTUFBQUEsTUFBQSxPQUFlLFdBQWYsZ0JBQUFBLElBQXVCLFVBQXZCLG1CQUE4QixvQkFBbUIsT0FBTztBQUN6RDtBQUFBLElBQ0o7QUFFQSxVQUFNLFVBQVUsU0FBUyxpQkFBaUIsR0FBRyxDQUFDO0FBQzlDLFVBQU0sYUFBYSxxQkFBcUIsT0FBTztBQUUvQyxRQUFJLENBQUMsWUFBWTtBQUViO0FBQUEsSUFDSjtBQUVBLFVBQU0saUJBQWlCO0FBQUEsTUFDbkIsSUFBSSxXQUFXO0FBQUEsTUFDZixXQUFXLE1BQU0sS0FBSyxXQUFXLFNBQVM7QUFBQSxNQUMxQyxZQUFZLENBQUM7QUFBQSxJQUNqQjtBQUNBLGFBQVMsSUFBSSxHQUFHLElBQUksV0FBVyxXQUFXLFFBQVEsS0FBSztBQUNuRCxZQUFNLE9BQU8sV0FBVyxXQUFXLENBQUM7QUFDcEMscUJBQWUsV0FBVyxLQUFLLElBQUksSUFBSSxLQUFLO0FBQUEsSUFDaEQ7QUFFQSxVQUFNLFVBQVU7QUFBQSxNQUNaO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsSUFDSjtBQUVBLFNBQUssU0FBUyxFQUFFLGNBQWMsT0FBTztBQUdyQyxzQkFBa0I7QUFBQSxFQUN0QjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLFVBQVUsVUFBaUM7QUFDdkMsV0FBTyxLQUFLLFNBQVMsRUFBRSxpQkFBaUIsRUFBRSxTQUFTLENBQUM7QUFBQSxFQUN4RDtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsYUFBNEI7QUFDeEIsV0FBTyxLQUFLLFNBQVMsRUFBRSxnQkFBZ0I7QUFBQSxFQUMzQztBQUFBO0FBQUE7QUFBQTtBQUFBLEVBS0EsUUFBdUI7QUFDbkIsV0FBTyxLQUFLLFNBQVMsRUFBRSxXQUFXO0FBQUEsRUFDdEM7QUFDSjtBQXRmQSxJQUFNLFNBQU47QUEyZkEsSUFBTSxhQUFhLElBQUksT0FBTyxFQUFFO0FBTWhDLFNBQVMsMkJBQTJCO0FBQ2hDLFFBQU0sYUFBYSxTQUFTO0FBQzVCLE1BQUksbUJBQW1CO0FBRXZCLGFBQVcsaUJBQWlCLGFBQWEsQ0FBQyxVQUFVO0FBdnNCeEQsUUFBQUEsS0FBQTtBQXdzQlEsUUFBSSxHQUFDQSxNQUFBLE1BQU0saUJBQU4sZ0JBQUFBLElBQW9CLE1BQU0sU0FBUyxXQUFVO0FBQzlDO0FBQUEsSUFDSjtBQUNBLFVBQU0sZUFBZTtBQUVyQixVQUFLLGtCQUFlLFdBQWYsbUJBQXVCLFVBQXZCLG1CQUE4QixvQkFBbUIsT0FBTztBQUN6RCxZQUFNLGFBQWEsYUFBYTtBQUNoQztBQUFBLElBQ0o7QUFDQTtBQUVBLFVBQU0sZ0JBQWdCLFNBQVMsaUJBQWlCLE1BQU0sU0FBUyxNQUFNLE9BQU87QUFDNUUsVUFBTSxhQUFhLHFCQUFxQixhQUFhO0FBR3JELFFBQUkscUJBQXFCLHNCQUFzQixZQUFZO0FBQ3ZELHdCQUFrQixVQUFVLE9BQU8sd0JBQXdCO0FBQUEsSUFDL0Q7QUFFQSxRQUFJLFlBQVk7QUFDWixpQkFBVyxVQUFVLElBQUksd0JBQXdCO0FBQ2pELFlBQU0sYUFBYSxhQUFhO0FBQ2hDLDBCQUFvQjtBQUFBLElBQ3hCLE9BQU87QUFDSCxZQUFNLGFBQWEsYUFBYTtBQUNoQywwQkFBb0I7QUFBQSxJQUN4QjtBQUFBLEVBQ0osR0FBRyxLQUFLO0FBRVIsYUFBVyxpQkFBaUIsWUFBWSxDQUFDLFVBQVU7QUFydUJ2RCxRQUFBQSxLQUFBO0FBc3VCUSxRQUFJLEdBQUNBLE1BQUEsTUFBTSxpQkFBTixnQkFBQUEsSUFBb0IsTUFBTSxTQUFTLFdBQVU7QUFDOUM7QUFBQSxJQUNKO0FBQ0EsVUFBTSxlQUFlO0FBRXJCLFVBQUssa0JBQWUsV0FBZixtQkFBdUIsVUFBdkIsbUJBQThCLG9CQUFtQixPQUFPO0FBQ3pELFlBQU0sYUFBYSxhQUFhO0FBQ2hDO0FBQUEsSUFDSjtBQUdBLFVBQU0sZ0JBQWdCLFNBQVMsaUJBQWlCLE1BQU0sU0FBUyxNQUFNLE9BQU87QUFDNUUsVUFBTSxhQUFhLHFCQUFxQixhQUFhO0FBRXJELFFBQUkscUJBQXFCLHNCQUFzQixZQUFZO0FBQ3ZELHdCQUFrQixVQUFVLE9BQU8sd0JBQXdCO0FBQUEsSUFDL0Q7QUFFQSxRQUFJLFlBQVk7QUFDWixVQUFJLENBQUMsV0FBVyxVQUFVLFNBQVMsd0JBQXdCLEdBQUc7QUFDMUQsbUJBQVcsVUFBVSxJQUFJLHdCQUF3QjtBQUFBLE1BQ3JEO0FBQ0EsWUFBTSxhQUFhLGFBQWE7QUFDaEMsMEJBQW9CO0FBQUEsSUFDeEIsT0FBTztBQUNILFlBQU0sYUFBYSxhQUFhO0FBQ2hDLDBCQUFvQjtBQUFBLElBQ3hCO0FBQUEsRUFDSixHQUFHLEtBQUs7QUFFUixhQUFXLGlCQUFpQixhQUFhLENBQUMsVUFBVTtBQXB3QnhELFFBQUFBLEtBQUE7QUFxd0JRLFFBQUksR0FBQ0EsTUFBQSxNQUFNLGlCQUFOLGdCQUFBQSxJQUFvQixNQUFNLFNBQVMsV0FBVTtBQUM5QztBQUFBLElBQ0o7QUFDQSxVQUFNLGVBQWU7QUFFckIsVUFBSyxrQkFBZSxXQUFmLG1CQUF1QixVQUF2QixtQkFBOEIsb0JBQW1CLE9BQU87QUFDekQ7QUFBQSxJQUNKO0FBSUEsUUFBSSxNQUFNLGtCQUFrQixNQUFNO0FBQzlCO0FBQUEsSUFDSjtBQUVBO0FBRUEsUUFBSSxxQkFBcUIsS0FDcEIscUJBQXFCLENBQUMsa0JBQWtCLFNBQVMsTUFBTSxhQUFxQixHQUFJO0FBQ2pGLFVBQUksbUJBQW1CO0FBQ25CLDBCQUFrQixVQUFVLE9BQU8sd0JBQXdCO0FBQzNELDRCQUFvQjtBQUFBLE1BQ3hCO0FBQ0EseUJBQW1CO0FBQUEsSUFDdkI7QUFBQSxFQUNKLEdBQUcsS0FBSztBQUVSLGFBQVcsaUJBQWlCLFFBQVEsQ0FBQyxVQUFVO0FBaHlCbkQsUUFBQUEsS0FBQTtBQWl5QlEsUUFBSSxHQUFDQSxNQUFBLE1BQU0saUJBQU4sZ0JBQUFBLElBQW9CLE1BQU0sU0FBUyxXQUFVO0FBQzlDO0FBQUEsSUFDSjtBQUNBLFVBQU0sZUFBZTtBQUVyQixVQUFLLGtCQUFlLFdBQWYsbUJBQXVCLFVBQXZCLG1CQUE4QixvQkFBbUIsT0FBTztBQUN6RDtBQUFBLElBQ0o7QUFDQSx1QkFBbUI7QUFFbkIsUUFBSSxtQkFBbUI7QUFDbkIsd0JBQWtCLFVBQVUsT0FBTyx3QkFBd0I7QUFDM0QsMEJBQW9CO0FBQUEsSUFDeEI7QUFJQSxRQUFJLG9CQUFvQixHQUFHO0FBQ3ZCLFlBQU0sUUFBZ0IsQ0FBQztBQUN2QixVQUFJLE1BQU0sYUFBYSxPQUFPO0FBQzFCLG1CQUFXLFFBQVEsTUFBTSxhQUFhLE9BQU87QUFDekMsY0FBSSxLQUFLLFNBQVMsUUFBUTtBQUN0QixrQkFBTSxPQUFPLEtBQUssVUFBVTtBQUM1QixnQkFBSSxLQUFNLE9BQU0sS0FBSyxJQUFJO0FBQUEsVUFDN0I7QUFBQSxRQUNKO0FBQUEsTUFDSixXQUFXLE1BQU0sYUFBYSxPQUFPO0FBQ2pDLG1CQUFXLFFBQVEsTUFBTSxhQUFhLE9BQU87QUFDekMsZ0JBQU0sS0FBSyxJQUFJO0FBQUEsUUFDbkI7QUFBQSxNQUNKO0FBRUEsVUFBSSxNQUFNLFNBQVMsR0FBRztBQUNsQix5QkFBaUIsTUFBTSxTQUFTLE1BQU0sU0FBUyxLQUFLO0FBQUEsTUFDeEQ7QUFBQSxJQUNKO0FBQUEsRUFDSixHQUFHLEtBQUs7QUFDWjtBQUdBLElBQUksT0FBTyxXQUFXLGVBQWUsT0FBTyxhQUFhLGFBQWE7QUFDbEUsMkJBQXlCO0FBQzdCO0FBRUEsSUFBTyxpQkFBUTs7O0FYdnpCZixTQUFTLFVBQVUsV0FBbUIsT0FBWSxNQUFZO0FBQzFELE9BQUssV0FBVyxJQUFJO0FBQ3hCO0FBUUEsU0FBUyxpQkFBaUIsWUFBb0IsWUFBb0I7QUFDOUQsUUFBTSxlQUFlLGVBQU8sSUFBSSxVQUFVO0FBQzFDLFFBQU0sU0FBVSxhQUFxQixVQUFVO0FBRS9DLE1BQUksT0FBTyxXQUFXLFlBQVk7QUFDOUIsWUFBUSxNQUFNLGtCQUFrQixtQkFBVSxjQUFhO0FBQ3ZEO0FBQUEsRUFDSjtBQUVBLE1BQUk7QUFDQSxXQUFPLEtBQUssWUFBWTtBQUFBLEVBQzVCLFNBQVMsR0FBRztBQUNSLFlBQVEsTUFBTSxnQ0FBZ0MsbUJBQVUsUUFBTyxDQUFDO0FBQUEsRUFDcEU7QUFDSjtBQUtBLFNBQVMsZUFBZSxJQUFpQjtBQUNyQyxRQUFNLFVBQVUsR0FBRztBQUVuQixXQUFTLFVBQVUsU0FBUyxPQUFPO0FBQy9CLFFBQUksV0FBVztBQUNYO0FBRUosVUFBTSxZQUFZLFFBQVEsYUFBYSxXQUFXLEtBQUssUUFBUSxhQUFhLGdCQUFnQjtBQUM1RixVQUFNLGVBQWUsUUFBUSxhQUFhLG1CQUFtQixLQUFLLFFBQVEsYUFBYSx3QkFBd0IsS0FBSztBQUNwSCxVQUFNLGVBQWUsUUFBUSxhQUFhLFlBQVksS0FBSyxRQUFRLGFBQWEsaUJBQWlCO0FBQ2pHLFVBQU0sTUFBTSxRQUFRLGFBQWEsYUFBYSxLQUFLLFFBQVEsYUFBYSxrQkFBa0I7QUFFMUYsUUFBSSxjQUFjO0FBQ2QsZ0JBQVUsU0FBUztBQUN2QixRQUFJLGlCQUFpQjtBQUNqQix1QkFBaUIsY0FBYyxZQUFZO0FBQy9DLFFBQUksUUFBUTtBQUNSLFdBQUssUUFBUSxHQUFHO0FBQUEsRUFDeEI7QUFFQSxRQUFNLFVBQVUsUUFBUSxhQUFhLGFBQWEsS0FBSyxRQUFRLGFBQWEsa0JBQWtCO0FBRTlGLE1BQUksU0FBUztBQUNULGFBQVM7QUFBQSxNQUNMLE9BQU87QUFBQSxNQUNQLFNBQVM7QUFBQSxNQUNULFVBQVU7QUFBQSxNQUNWLFNBQVM7QUFBQSxRQUNMLEVBQUUsT0FBTyxNQUFNO0FBQUEsUUFDZixFQUFFLE9BQU8sTUFBTSxXQUFXLEtBQUs7QUFBQSxNQUNuQztBQUFBLElBQ0osQ0FBQyxFQUFFLEtBQUssU0FBUztBQUFBLEVBQ3JCLE9BQU87QUFDSCxjQUFVO0FBQUEsRUFDZDtBQUNKO0FBR0EsSUFBTSxnQkFBZ0IsdUJBQU8sWUFBWTtBQUN6QyxJQUFNLGdCQUFnQix1QkFBTyxZQUFZO0FBQ3pDLElBQU0sa0JBQWtCLHVCQUFPLGNBQWM7QUFReEM7QUFGTCxJQUFNLDBCQUFOLE1BQThCO0FBQUEsRUFJMUIsY0FBYztBQUNWLFNBQUssYUFBYSxJQUFJLElBQUksZ0JBQWdCO0FBQUEsRUFDOUM7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBU0EsSUFBSSxTQUFrQixVQUE2QztBQUMvRCxXQUFPLEVBQUUsUUFBUSxLQUFLLGFBQWEsRUFBRSxPQUFPO0FBQUEsRUFDaEQ7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFFBQWM7QUFDVixTQUFLLGFBQWEsRUFBRSxNQUFNO0FBQzFCLFNBQUssYUFBYSxJQUFJLElBQUksZ0JBQWdCO0FBQUEsRUFDOUM7QUFDSjtBQVNLLGVBRUE7QUFKTCxJQUFNLGtCQUFOLE1BQXNCO0FBQUEsRUFNbEIsY0FBYztBQUNWLFNBQUssYUFBYSxJQUFJLG9CQUFJLFFBQVE7QUFDbEMsU0FBSyxlQUFlLElBQUk7QUFBQSxFQUM1QjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBUUEsSUFBSSxTQUFrQixVQUE2QztBQUMvRCxRQUFJLENBQUMsS0FBSyxhQUFhLEVBQUUsSUFBSSxPQUFPLEdBQUc7QUFBRSxXQUFLLGVBQWU7QUFBQSxJQUFLO0FBQ2xFLFNBQUssYUFBYSxFQUFFLElBQUksU0FBUyxRQUFRO0FBQ3pDLFdBQU8sQ0FBQztBQUFBLEVBQ1o7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQUtBLFFBQWM7QUFDVixRQUFJLEtBQUssZUFBZSxLQUFLO0FBQ3pCO0FBRUosZUFBVyxXQUFXLFNBQVMsS0FBSyxpQkFBaUIsR0FBRyxHQUFHO0FBQ3ZELFVBQUksS0FBSyxlQUFlLEtBQUs7QUFDekI7QUFFSixZQUFNLFdBQVcsS0FBSyxhQUFhLEVBQUUsSUFBSSxPQUFPO0FBQ2hELFVBQUksWUFBWSxNQUFNO0FBQUUsYUFBSyxlQUFlO0FBQUEsTUFBSztBQUVqRCxpQkFBVyxXQUFXLFlBQVksQ0FBQztBQUMvQixnQkFBUSxvQkFBb0IsU0FBUyxjQUFjO0FBQUEsSUFDM0Q7QUFFQSxTQUFLLGFBQWEsSUFBSSxvQkFBSSxRQUFRO0FBQ2xDLFNBQUssZUFBZSxJQUFJO0FBQUEsRUFDNUI7QUFDSjtBQUVBLElBQU0sa0JBQWtCLGtCQUFrQixJQUFJLElBQUksd0JBQXdCLElBQUksSUFBSSxnQkFBZ0I7QUFLbEcsU0FBUyxnQkFBZ0IsU0FBd0I7QUFDN0MsUUFBTSxnQkFBZ0I7QUFDdEIsUUFBTSxjQUFlLFFBQVEsYUFBYSxhQUFhLEtBQUssUUFBUSxhQUFhLGtCQUFrQixLQUFLO0FBQ3hHLFFBQU0sV0FBcUIsQ0FBQztBQUU1QixNQUFJO0FBQ0osVUFBUSxRQUFRLGNBQWMsS0FBSyxXQUFXLE9BQU87QUFDakQsYUFBUyxLQUFLLE1BQU0sQ0FBQyxDQUFDO0FBRTFCLFFBQU0sVUFBVSxnQkFBZ0IsSUFBSSxTQUFTLFFBQVE7QUFDckQsYUFBVyxXQUFXO0FBQ2xCLFlBQVEsaUJBQWlCLFNBQVMsZ0JBQWdCLE9BQU87QUFDakU7QUFLTyxTQUFTLFNBQWU7QUFDM0IsWUFBVSxNQUFNO0FBQ3BCO0FBS08sU0FBUyxTQUFlO0FBQzNCLGtCQUFnQixNQUFNO0FBQ3RCLFdBQVMsS0FBSyxpQkFBaUIsbUdBQW1HLEVBQUUsUUFBUSxlQUFlO0FBQy9KOzs7QVloTUEsT0FBTyxRQUFRO0FBQ2YsT0FBVTtBQUVWLElBQUksTUFBTztBQUNQLFdBQVMsc0JBQXNCO0FBQ25DOzs7QUNyQkE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBWUEsSUFBTUMsUUFBTyxpQkFBaUIsWUFBWSxNQUFNO0FBRWhELElBQU0sbUJBQW1CO0FBQ3pCLElBQU0sb0JBQW9CO0FBQzFCLElBQU0scUJBQXFCO0FBRTNCLElBQU0sV0FBVyxXQUFZO0FBbEI3QixNQUFBQyxLQUFBO0FBbUJJLE1BQUk7QUFFQSxTQUFLLE1BQUFBLE1BQUEsT0FBZSxXQUFmLGdCQUFBQSxJQUF1QixZQUF2QixtQkFBZ0MsYUFBYTtBQUM5QyxhQUFRLE9BQWUsT0FBTyxRQUFRLFlBQVksS0FBTSxPQUFlLE9BQU8sT0FBTztBQUFBLElBQ3pGLFlBRVUsd0JBQWUsV0FBZixtQkFBdUIsb0JBQXZCLG1CQUF5QyxnQkFBekMsbUJBQXNELGFBQWE7QUFDekUsYUFBUSxPQUFlLE9BQU8sZ0JBQWdCLFVBQVUsRUFBRSxZQUFZLEtBQU0sT0FBZSxPQUFPLGdCQUFnQixVQUFVLENBQUM7QUFBQSxJQUNqSSxZQUVVLFlBQWUsVUFBZixtQkFBc0IsUUFBUTtBQUNwQyxhQUFPLENBQUMsUUFBYyxPQUFlLE1BQU0sT0FBTyxPQUFPLFFBQVEsV0FBVyxNQUFNLEtBQUssVUFBVSxHQUFHLENBQUM7QUFBQSxJQUN6RztBQUFBLEVBQ0osU0FBUSxHQUFHO0FBQUEsRUFBQztBQUVaLFVBQVE7QUFBQSxJQUFLO0FBQUEsSUFDVDtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsRUFBd0Q7QUFDNUQsU0FBTztBQUNYLEdBQUc7QUFFSSxTQUFTLE9BQU8sS0FBZ0I7QUFDbkMscUNBQVU7QUFDZDtBQU9PLFNBQVMsYUFBK0I7QUFDM0MsU0FBT0QsTUFBSyxnQkFBZ0I7QUFDaEM7QUFPQSxlQUFzQixlQUE2QztBQUMvRCxTQUFPQSxNQUFLLGtCQUFrQjtBQUNsQztBQStCTyxTQUFTLGNBQXdDO0FBQ3BELFNBQU9BLE1BQUssaUJBQWlCO0FBQ2pDO0FBT08sU0FBUyxZQUFxQjtBQXJHckMsTUFBQUMsS0FBQTtBQXNHSSxXQUFRLE1BQUFBLE1BQUEsT0FBZSxXQUFmLGdCQUFBQSxJQUF1QixnQkFBdkIsbUJBQW9DLFFBQU87QUFDdkQ7QUFPTyxTQUFTLFVBQW1CO0FBOUduQyxNQUFBQSxLQUFBO0FBK0dJLFdBQVEsTUFBQUEsTUFBQSxPQUFlLFdBQWYsZ0JBQUFBLElBQXVCLGdCQUF2QixtQkFBb0MsUUFBTztBQUN2RDtBQU9PLFNBQVMsUUFBaUI7QUF2SGpDLE1BQUFBLEtBQUE7QUF3SEksV0FBUSxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsZ0JBQXZCLG1CQUFvQyxRQUFPO0FBQ3ZEO0FBT08sU0FBUyxRQUFpQjtBQWhJakMsTUFBQUEsS0FBQTtBQWlJSSxXQUFRLE1BQUFBLE1BQUEsT0FBZSxXQUFmLGdCQUFBQSxJQUF1QixnQkFBdkIsbUJBQW9DLFFBQU87QUFDdkQ7QUFPTyxTQUFTLFlBQXFCO0FBeklyQyxNQUFBQSxLQUFBO0FBMElJLFdBQVEsTUFBQUEsTUFBQSxPQUFlLFdBQWYsZ0JBQUFBLElBQXVCLGdCQUF2QixtQkFBb0MsUUFBTztBQUN2RDtBQU9PLFNBQVMsV0FBb0I7QUFDaEMsU0FBTyxNQUFNLEtBQUssVUFBVTtBQUNoQztBQU9PLFNBQVMsWUFBcUI7QUFDakMsU0FBTyxNQUFNLEtBQUssVUFBVSxLQUFLLFFBQVE7QUFDN0M7QUFPTyxTQUFTLFVBQW1CO0FBcEtuQyxNQUFBQSxLQUFBO0FBcUtJLFdBQVEsTUFBQUEsTUFBQSxPQUFlLFdBQWYsZ0JBQUFBLElBQXVCLGdCQUF2QixtQkFBb0MsVUFBUztBQUN6RDtBQU9PLFNBQVMsUUFBaUI7QUE3S2pDLE1BQUFBLEtBQUE7QUE4S0ksV0FBUSxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsZ0JBQXZCLG1CQUFvQyxVQUFTO0FBQ3pEO0FBT08sU0FBUyxVQUFtQjtBQXRMbkMsTUFBQUEsS0FBQTtBQXVMSSxXQUFRLE1BQUFBLE1BQUEsT0FBZSxXQUFmLGdCQUFBQSxJQUF1QixnQkFBdkIsbUJBQW9DLFVBQVM7QUFDekQ7QUFPTyxTQUFTLFVBQW1CO0FBL0xuQyxNQUFBQSxLQUFBO0FBZ01JLFNBQU8sU0FBUyxNQUFBQSxNQUFBLE9BQWUsV0FBZixnQkFBQUEsSUFBdUIsZ0JBQXZCLG1CQUFvQyxLQUFLO0FBQzdEOzs7QUNoTEEsSUFBSSxRQUFRO0FBQ1IsU0FBTyxpQkFBaUIsZUFBZSxrQkFBa0I7QUFDN0Q7QUFFQSxJQUFNQyxRQUFPLGlCQUFpQixZQUFZLFdBQVc7QUFFckQsSUFBTSxrQkFBa0I7QUFFeEIsU0FBUyxnQkFBZ0IsSUFBWSxHQUFXLEdBQVcsTUFBaUI7QUFDeEUsT0FBS0EsTUFBSyxpQkFBaUIsRUFBQyxJQUFJLEdBQUcsR0FBRyxLQUFJLENBQUM7QUFDL0M7QUFFQSxTQUFTLG1CQUFtQixPQUFtQjtBQUMzQyxRQUFNLFNBQVMsWUFBWSxLQUFLO0FBR2hDLFFBQU0sb0JBQW9CLE9BQU8saUJBQWlCLE1BQU0sRUFBRSxpQkFBaUIsc0JBQXNCLEVBQUUsS0FBSztBQUV4RyxNQUFJLG1CQUFtQjtBQUNuQixVQUFNLGVBQWU7QUFDckIsVUFBTSxPQUFPLE9BQU8saUJBQWlCLE1BQU0sRUFBRSxpQkFBaUIsMkJBQTJCO0FBQ3pGLG9CQUFnQixtQkFBbUIsTUFBTSxTQUFTLE1BQU0sU0FBUyxJQUFJO0FBQUEsRUFDekUsT0FBTztBQUNILDhCQUEwQixPQUFPLE1BQU07QUFBQSxFQUMzQztBQUNKO0FBVUEsU0FBUywwQkFBMEIsT0FBbUIsUUFBcUI7QUFFdkUsTUFBSSxRQUFRLEdBQUc7QUFDWDtBQUFBLEVBQ0o7QUFHQSxVQUFRLE9BQU8saUJBQWlCLE1BQU0sRUFBRSxpQkFBaUIsdUJBQXVCLEVBQUUsS0FBSyxHQUFHO0FBQUEsSUFDdEYsS0FBSztBQUNEO0FBQUEsSUFDSixLQUFLO0FBQ0QsWUFBTSxlQUFlO0FBQ3JCO0FBQUEsRUFDUjtBQUdBLE1BQUksT0FBTyxtQkFBbUI7QUFDMUI7QUFBQSxFQUNKO0FBR0EsUUFBTSxZQUFZLE9BQU8sYUFBYTtBQUN0QyxRQUFNLGVBQWUsYUFBYSxVQUFVLFNBQVMsRUFBRSxTQUFTO0FBQ2hFLE1BQUksY0FBYztBQUNkLGFBQVMsSUFBSSxHQUFHLElBQUksVUFBVSxZQUFZLEtBQUs7QUFDM0MsWUFBTSxRQUFRLFVBQVUsV0FBVyxDQUFDO0FBQ3BDLFlBQU0sUUFBUSxNQUFNLGVBQWU7QUFDbkMsZUFBUyxJQUFJLEdBQUcsSUFBSSxNQUFNLFFBQVEsS0FBSztBQUNuQyxjQUFNLE9BQU8sTUFBTSxDQUFDO0FBQ3BCLFlBQUksU0FBUyxpQkFBaUIsS0FBSyxNQUFNLEtBQUssR0FBRyxNQUFNLFFBQVE7QUFDM0Q7QUFBQSxRQUNKO0FBQUEsTUFDSjtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBR0EsTUFBSSxrQkFBa0Isb0JBQW9CLGtCQUFrQixxQkFBcUI7QUFDN0UsUUFBSSxnQkFBaUIsQ0FBQyxPQUFPLFlBQVksQ0FBQyxPQUFPLFVBQVc7QUFDeEQ7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUdBLFFBQU0sZUFBZTtBQUN6Qjs7O0FDakdBO0FBQUE7QUFBQTtBQUFBO0FBZ0JPLFNBQVMsUUFBUSxLQUFrQjtBQUN0QyxNQUFJO0FBQ0EsV0FBTyxPQUFPLE9BQU8sTUFBTSxHQUFHO0FBQUEsRUFDbEMsU0FBUyxHQUFHO0FBQ1IsVUFBTSxJQUFJLE1BQU0sOEJBQThCLE1BQU0sUUFBUSxHQUFHLEVBQUUsT0FBTyxFQUFFLENBQUM7QUFBQSxFQUMvRTtBQUNKOzs7QUNQQSxJQUFJLFVBQVU7QUFDZCxJQUFJLFdBQVc7QUFFZixJQUFJLFlBQVk7QUFDaEIsSUFBSSxZQUFZO0FBQ2hCLElBQUksV0FBVztBQUNmLElBQUksYUFBcUI7QUFDekIsSUFBSSxnQkFBZ0I7QUFFcEIsSUFBSSxVQUFVO0FBR2QsSUFBSSxpQkFBaUI7QUFFckIsSUFBSSxRQUFRO0FBQ1IsbUJBQWlCLGdCQUFnQjtBQUNqQyxTQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFDbEMsU0FBTyxPQUFPLGVBQWUsQ0FBQyxVQUF5QjtBQUNuRCxnQkFBWTtBQUNaLFFBQUksQ0FBQyxXQUFXO0FBRVosa0JBQVksV0FBVztBQUN2QixnQkFBVTtBQUFBLElBQ2Q7QUFBQSxFQUNKO0FBQ0o7QUFHQSxJQUFJLGVBQWU7QUFDbkIsU0FBUyxXQUFvQjtBQTVDN0IsTUFBQUMsS0FBQTtBQTZDSSxRQUFNLE1BQU0sTUFBQUEsTUFBQSxPQUFlLFdBQWYsZ0JBQUFBLElBQXVCLGdCQUF2QixtQkFBb0M7QUFDaEQsTUFBSSxPQUFPLFNBQVMsT0FBTyxVQUFXLFFBQU87QUFFN0MsUUFBTSxLQUFLLFVBQVUsYUFBYSxVQUFVLFVBQVcsT0FBZSxTQUFTO0FBQy9FLFNBQU8sK0NBQStDLEtBQUssRUFBRTtBQUNqRTtBQUNBLFNBQVMsc0JBQTRCO0FBQ2pDLE1BQUksYUFBYztBQUNsQixNQUFJLFNBQVMsRUFBRztBQUNoQixTQUFPLGlCQUFpQixhQUFhLFFBQVEsRUFBRSxTQUFTLEtBQUssQ0FBQztBQUM5RCxTQUFPLGlCQUFpQixhQUFhLFFBQVEsRUFBRSxTQUFTLEtBQUssQ0FBQztBQUM5RCxTQUFPLGlCQUFpQixXQUFXLFFBQVEsRUFBRSxTQUFTLEtBQUssQ0FBQztBQUM1RCxhQUFXLE1BQU0sQ0FBQyxTQUFTLGVBQWUsVUFBVSxHQUFHO0FBQ25ELFdBQU8saUJBQWlCLElBQUksZUFBZSxFQUFFLFNBQVMsS0FBSyxDQUFDO0FBQUEsRUFDaEU7QUFDQSxpQkFBZTtBQUNuQjtBQUNBLElBQUksUUFBUTtBQUVSLHNCQUFvQjtBQUVwQixXQUFTLGlCQUFpQixvQkFBb0IscUJBQXFCLEVBQUUsTUFBTSxLQUFLLENBQUM7QUFFakYsTUFBSSxlQUFlO0FBQ25CLFFBQU0sY0FBYyxPQUFPLFlBQVksTUFBTTtBQUN6QyxRQUFJLGNBQWM7QUFBRSxhQUFPLGNBQWMsV0FBVztBQUFHO0FBQUEsSUFBUTtBQUMvRCx3QkFBb0I7QUFDcEIsUUFBSSxFQUFFLGVBQWUsS0FBSztBQUFFLGFBQU8sY0FBYyxXQUFXO0FBQUEsSUFBRztBQUFBLEVBQ25FLEdBQUcsRUFBRTtBQUNUO0FBRUEsU0FBUyxjQUFjLE9BQWM7QUFDakMsTUFDSSxNQUFNLFNBQVMsY0FDWixNQUFNLEtBQ04sV0FBVyxPQUFPLE9BQ2xCLGlCQUFpQixLQUFtQixHQUN6QztBQUNFLFVBQU0seUJBQXlCO0FBQy9CLFVBQU0sZ0JBQWdCO0FBQ3RCLFVBQU0sZUFBZTtBQUNyQixXQUFPLHdCQUF3QjtBQUMvQjtBQUFBLEVBQ0o7QUFHQSxNQUFJLFlBQVksVUFBVTtBQUN0QixVQUFNLHlCQUF5QjtBQUMvQixVQUFNLGdCQUFnQjtBQUN0QixVQUFNLGVBQWU7QUFBQSxFQUN6QjtBQUNKO0FBRUEsU0FBUyxpQkFBaUIsT0FBNEI7QUFDbEQsUUFBTSxTQUFTLFlBQVksS0FBSztBQUNoQyxRQUFNLFFBQVEsT0FBTyxpQkFBaUIsTUFBTTtBQUM1QyxTQUNJLE1BQU0saUJBQWlCLG1CQUFtQixFQUFFLEtBQUssTUFBTSxVQUNwRCxNQUFNLFdBQVcsS0FDakIsTUFBTSxVQUFVLE9BQU8sZUFDdkIsTUFBTSxXQUFXLEtBQ2pCLE1BQU0sVUFBVSxPQUFPO0FBRWxDO0FBR0EsSUFBTSxZQUFZO0FBQ2xCLElBQU0sVUFBWTtBQUNsQixJQUFNLFlBQVk7QUFFbEIsU0FBUyxPQUFPLE9BQW1CO0FBSS9CLE1BQUksV0FBbUIsZUFBZSxNQUFNO0FBQzVDLFVBQVEsTUFBTSxNQUFNO0FBQUEsSUFDaEIsS0FBSztBQUNELGtCQUFZO0FBQ1osVUFBSSxDQUFDLGdCQUFnQjtBQUFFLHVCQUFlLFVBQVcsS0FBSyxNQUFNO0FBQUEsTUFBUztBQUNyRTtBQUFBLElBQ0osS0FBSztBQUNELGtCQUFZO0FBQ1osVUFBSSxDQUFDLGdCQUFnQjtBQUFFLHVCQUFlLFVBQVUsRUFBRSxLQUFLLE1BQU07QUFBQSxNQUFTO0FBQ3RFO0FBQUEsSUFDSjtBQUNJLGtCQUFZO0FBQ1osVUFBSSxDQUFDLGdCQUFnQjtBQUFFLHVCQUFlO0FBQUEsTUFBUztBQUMvQztBQUFBLEVBQ1I7QUFFQSxNQUFJLFdBQVcsVUFBVSxDQUFDO0FBQzFCLE1BQUksVUFBVSxlQUFlLENBQUM7QUFFOUIsWUFBVTtBQUdWLE1BQUksY0FBYyxhQUFhLEVBQUUsVUFBVSxNQUFNLFNBQVM7QUFDdEQsZ0JBQWEsS0FBSyxNQUFNO0FBQ3hCLGVBQVksS0FBSyxNQUFNO0FBQUEsRUFDM0I7QUFJQSxNQUNJLGNBQWMsYUFDWCxZQUVDLGFBRUksY0FBYyxhQUNYLE1BQU0sV0FBVyxJQUc5QjtBQUNFLFVBQU0seUJBQXlCO0FBQy9CLFVBQU0sZ0JBQWdCO0FBQ3RCLFVBQU0sZUFBZTtBQUFBLEVBQ3pCO0FBR0EsTUFBSSxXQUFXLEdBQUc7QUFBRSxjQUFVLEtBQUs7QUFBQSxFQUFHO0FBRXRDLE1BQUksVUFBVSxHQUFHO0FBQUUsZ0JBQVksS0FBSztBQUFBLEVBQUc7QUFHdkMsTUFBSSxjQUFjLFdBQVc7QUFBRSxnQkFBWSxLQUFLO0FBQUEsRUFBRztBQUFDO0FBQ3hEO0FBRUEsU0FBUyxZQUFZLE9BQXlCO0FBRTFDLFlBQVU7QUFDVixjQUFZO0FBR1osTUFBSSxDQUFDLFVBQVUsR0FBRztBQUNkLFFBQUksTUFBTSxTQUFTLGVBQWUsTUFBTSxXQUFXLEtBQUssTUFBTSxXQUFXLEdBQUc7QUFDeEU7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUVBLE1BQUksWUFBWTtBQUlaLFFBQUksTUFBTSxTQUFTLGFBQWE7QUFDNUI7QUFBQSxJQUNKO0FBR0EsZ0JBQVk7QUFFWjtBQUFBLEVBQ0o7QUFLQSxZQUFVLGlCQUFpQixLQUFLO0FBQ3BDO0FBRUEsU0FBUyxVQUFVLE9BQW1CO0FBRWxDLFlBQVU7QUFDVixhQUFXO0FBQ1gsY0FBWTtBQUNaLGFBQVc7QUFDZjtBQUVBLElBQU0sZ0JBQWdCLE9BQU8sT0FBTztBQUFBLEVBQ2hDLGFBQWE7QUFBQSxFQUNiLGFBQWE7QUFBQSxFQUNiLGFBQWE7QUFBQSxFQUNiLGFBQWE7QUFBQSxFQUNiLFlBQVk7QUFBQSxFQUNaLFlBQVk7QUFBQSxFQUNaLFlBQVk7QUFBQSxFQUNaLFlBQVk7QUFDaEIsQ0FBQztBQUVELFNBQVMsVUFBVSxNQUF5QztBQUN4RCxNQUFJLE1BQU07QUFDTixRQUFJLENBQUMsWUFBWTtBQUFFLHNCQUFnQixTQUFTLEtBQUssTUFBTTtBQUFBLElBQVE7QUFDL0QsYUFBUyxLQUFLLE1BQU0sU0FBUyxjQUFjLElBQUk7QUFBQSxFQUNuRCxXQUFXLENBQUMsUUFBUSxZQUFZO0FBQzVCLGFBQVMsS0FBSyxNQUFNLFNBQVM7QUFBQSxFQUNqQztBQUVBLGVBQWEsUUFBUTtBQUN6QjtBQUVBLFNBQVMsWUFBWSxPQUF5QjtBQUMxQyxNQUFJLGFBQWEsWUFBWTtBQUV6QixlQUFXO0FBQ1gsV0FBTyxrQkFBa0IsVUFBVTtBQUFBLEVBQ3ZDLFdBQVcsU0FBUztBQUVoQixlQUFXO0FBQ1gsV0FBTyxZQUFZO0FBQUEsRUFDdkI7QUFFQSxNQUFJLFlBQVksVUFBVTtBQUd0QixjQUFVLFlBQVk7QUFDdEI7QUFBQSxFQUNKO0FBRUEsTUFBSSxDQUFDLGFBQWMsQ0FBQyxVQUFVLEtBQUssRUFBRSxRQUFRLEtBQUssUUFBUSxXQUFXLElBQUs7QUFDdEUsUUFBSSxZQUFZO0FBQUUsZ0JBQVU7QUFBQSxJQUFHO0FBQy9CO0FBQUEsRUFDSjtBQUVBLFFBQU0scUJBQXFCLFFBQVEsMkJBQTJCLEtBQUs7QUFDbkUsUUFBTSxvQkFBb0IsUUFBUSwwQkFBMEIsS0FBSztBQUdqRSxRQUFNLGNBQWMsUUFBUSxtQkFBbUIsS0FBSztBQUlwRCxRQUFNLGlCQUFpQixLQUFLLElBQUksR0FBRyxPQUFPLGFBQWEsU0FBUyxnQkFBZ0IsV0FBVztBQUMzRixRQUFNLGtCQUFrQixLQUFLLElBQUksR0FBRyxPQUFPLGNBQWMsU0FBUyxnQkFBZ0IsWUFBWTtBQUM5RixRQUFNLG1CQUFtQixPQUFPLGFBQWE7QUFDN0MsUUFBTSxvQkFBb0IsT0FBTyxjQUFjO0FBRS9DLFFBQU0sY0FBYyxNQUFNLFVBQVUsb0JBQXFCLG1CQUFtQixNQUFNLFVBQVc7QUFDN0YsUUFBTSxhQUFhLE1BQU0sVUFBVTtBQUNuQyxRQUFNLFlBQVksTUFBTSxVQUFVO0FBQ2xDLFFBQU0sZUFBZSxNQUFNLFVBQVUscUJBQXNCLG9CQUFvQixNQUFNLFVBQVc7QUFHaEcsUUFBTSxjQUFjLE1BQU0sVUFBVSxvQkFBcUIsbUJBQW1CLE1BQU0sVUFBWSxvQkFBb0I7QUFDbEgsUUFBTSxhQUFhLE1BQU0sVUFBVyxvQkFBb0I7QUFDeEQsUUFBTSxZQUFZLE1BQU0sVUFBVyxxQkFBcUI7QUFDeEQsUUFBTSxlQUFlLE1BQU0sVUFBVSxxQkFBc0Isb0JBQW9CLE1BQU0sVUFBWSxxQkFBcUI7QUFFdEgsTUFBSSxDQUFDLGNBQWMsQ0FBQyxhQUFhLENBQUMsZ0JBQWdCLENBQUMsYUFBYTtBQUU1RCxjQUFVO0FBQUEsRUFDZCxXQUVTLGVBQWUsYUFBYyxXQUFVLFdBQVc7QUFBQSxXQUNsRCxjQUFjLGFBQWMsV0FBVSxXQUFXO0FBQUEsV0FDakQsY0FBYyxVQUFXLFdBQVUsV0FBVztBQUFBLFdBQzlDLGFBQWEsWUFBYSxXQUFVLFdBQVc7QUFBQSxXQUUvQyxXQUFZLFdBQVUsVUFBVTtBQUFBLFdBQ2hDLFVBQVcsV0FBVSxVQUFVO0FBQUEsV0FDL0IsYUFBYyxXQUFVLFVBQVU7QUFBQSxXQUNsQyxZQUFhLFdBQVUsVUFBVTtBQUFBLE1BRXJDLFdBQVU7QUFDbkI7OztBQzVRQSxJQUFNLGlCQUFpQjtBQUN2QixJQUFNLDBCQUEwQjtBQUNoQyxJQUFNLGVBQWUsb0JBQUksSUFBeUIsQ0FBQyxXQUFXLFlBQVksWUFBWSxPQUFPLENBQUM7QUFHOUYsSUFBSSxRQUFRO0FBQ1IsU0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBQ3RDO0FBRUEsSUFBSSxnQkFBZ0I7QUFDcEIsSUFBSSxjQUFjO0FBQ2xCLElBQUksbUJBQW1CLG9CQUFJLElBQWE7QUFDeEMsSUFBSTtBQUNKLElBQUksa0JBQWtCO0FBRXRCLFNBQVMsb0JBQW9CLE9BQWdEO0FBQ3pFLFFBQU0sU0FBUyxNQUFNLEtBQUssRUFBRSxZQUFZO0FBQ3hDLE1BQUksYUFBYSxJQUFJLE1BQTZCLEdBQUc7QUFDakQsV0FBTztBQUFBLEVBQ1g7QUFDQSxTQUFPO0FBQ1g7QUFFQSxTQUFTLDBCQUEwQixTQUFtRDtBQUNsRixNQUFJLEVBQUUsbUJBQW1CLGNBQWM7QUFDbkMsV0FBTztBQUFBLEVBQ1g7QUFFQSxRQUFNLFFBQVEsT0FBTyxpQkFBaUIsT0FBTztBQUM3QyxRQUFNLFNBQVMsb0JBQW9CLE1BQU0saUJBQWlCLGNBQWMsQ0FBQztBQUN6RSxNQUFJLENBQUMsUUFBUTtBQUNULFdBQU87QUFBQSxFQUNYO0FBRUEsUUFBTSxTQUFTLFFBQVE7QUFDdkIsTUFBSSxRQUFRO0FBQ1IsVUFBTSxjQUFjLE9BQU8saUJBQWlCLE1BQU07QUFHbEQsUUFBSSxvQkFBb0IsWUFBWSxpQkFBaUIsY0FBYyxDQUFDLE1BQU0sUUFBUTtBQUM5RSxhQUFPO0FBQUEsSUFDWDtBQUFBLEVBQ0o7QUFFQSxTQUFPO0FBQ1g7QUFFQSxTQUFTLFVBQVUsU0FBK0I7QUFDOUMsUUFBTSxRQUFRLE9BQU8saUJBQWlCLE9BQU87QUFDN0MsU0FBTyxNQUFNLFlBQVksVUFDckIsTUFBTSxlQUFlLFlBQ3JCLE1BQU0sc0JBQXNCO0FBQ3BDO0FBRUEsU0FBUyxjQUFjLFNBQStDO0FBQ2xFLE1BQUksRUFBRSxtQkFBbUIsY0FBYztBQUNuQyxXQUFPO0FBQUEsRUFDWDtBQUVBLFFBQU0sT0FBTywwQkFBMEIsT0FBTztBQUM5QyxNQUFJLENBQUMsUUFBUSxDQUFDLFVBQVUsT0FBTyxHQUFHO0FBQzlCLFdBQU87QUFBQSxFQUNYO0FBRUEsUUFBTSxPQUFPLFFBQVEsc0JBQXNCO0FBQzNDLE1BQUksS0FBSyxTQUFTLEtBQUssS0FBSyxVQUFVLEdBQUc7QUFDckMsV0FBTztBQUFBLEVBQ1g7QUFHQSxRQUFNLFFBQVEsT0FBTyxvQkFBb0I7QUFDekMsUUFBTSxPQUFPLEtBQUssTUFBTSxLQUFLLE9BQU8sS0FBSztBQUN6QyxRQUFNLE1BQU0sS0FBSyxNQUFNLEtBQUssTUFBTSxLQUFLO0FBQ3ZDLFFBQU0sUUFBUSxLQUFLLEtBQUssS0FBSyxRQUFRLEtBQUs7QUFDMUMsUUFBTSxTQUFTLEtBQUssS0FBSyxLQUFLLFNBQVMsS0FBSztBQUU1QyxNQUFJLFNBQVMsUUFBUSxVQUFVLEtBQUs7QUFDaEMsV0FBTztBQUFBLEVBQ1g7QUFFQSxTQUFPLEVBQUUsTUFBTSxNQUFNLEtBQUssT0FBTyxPQUFPO0FBQzVDO0FBRUEsU0FBUyxpQkFBNEI7QUFDakMsUUFBTSxXQUFzQixDQUFDO0FBRTdCLE1BQUksU0FBUyxpQkFBaUI7QUFDMUIsYUFBUyxLQUFLLFNBQVMsZUFBZTtBQUFBLEVBQzFDO0FBQ0EsTUFBSSxTQUFTLE1BQU07QUFDZixhQUFTLEtBQUssU0FBUyxJQUFJO0FBRzNCLGVBQVcsV0FBVyxTQUFTLEtBQUssaUJBQWlCLEdBQUcsR0FBRztBQUN2RCxlQUFTLEtBQUssT0FBTztBQUFBLElBQ3pCO0FBQUEsRUFDSjtBQUVBLFNBQU87QUFDWDtBQUVBLFNBQVMsc0JBQXNCLFVBQTJCO0FBQ3RELE1BQUksT0FBTyxtQkFBbUIsYUFBYTtBQUN2QztBQUFBLEVBQ0o7QUFJQSw2REFBbUIsSUFBSSxlQUFlLGNBQWM7QUFDcEQsUUFBTSxlQUFlLElBQUksSUFBSSxRQUFRO0FBRXJDLGFBQVcsV0FBVyxrQkFBa0I7QUFDcEMsUUFBSSxDQUFDLGFBQWEsSUFBSSxPQUFPLEdBQUc7QUFDNUIscUJBQWUsVUFBVSxPQUFPO0FBQUEsSUFDcEM7QUFBQSxFQUNKO0FBRUEsYUFBVyxXQUFXLGNBQWM7QUFDaEMsUUFBSSxDQUFDLGlCQUFpQixJQUFJLE9BQU8sR0FBRztBQUNoQyxxQkFBZSxRQUFRLE9BQU87QUFBQSxJQUNsQztBQUFBLEVBQ0o7QUFFQSxxQkFBbUI7QUFDdkI7QUFFQSxTQUFTLHlCQUErQjtBQUNwQyxrQkFBZ0I7QUFFaEIsUUFBTSxXQUFXLGVBQWU7QUFDaEMsUUFBTSxVQUE2QixDQUFDO0FBQ3BDLFFBQU0saUJBQTRCLENBQUM7QUFFbkMsYUFBVyxXQUFXLFVBQVU7QUFDNUIsVUFBTSxTQUFTLGNBQWMsT0FBTztBQUNwQyxRQUFJLFFBQVE7QUFDUixjQUFRLEtBQUssTUFBTTtBQUNuQixxQkFBZSxLQUFLLE9BQU87QUFBQSxJQUMvQjtBQUFBLEVBQ0o7QUFFQSx3QkFBc0IsY0FBYztBQUVwQyxRQUFNLFVBQVUsS0FBSyxVQUFVLEVBQUUsU0FBUyxHQUFHLFFBQVEsQ0FBQztBQUN0RCxNQUFJLFlBQVksYUFBYTtBQUV6QjtBQUFBLEVBQ0o7QUFFQSxnQkFBYztBQUNkLFNBQU8sNkJBQTZCLE9BQU87QUFDL0M7QUFFQSxTQUFTLGlCQUF1QjtBQUM1QixNQUFJLGVBQWU7QUFDZjtBQUFBLEVBQ0o7QUFHQSxrQkFBZ0I7QUFDaEIsU0FBTyxzQkFBc0Isc0JBQXNCO0FBQ3ZEO0FBRUEsU0FBUywrQkFBcUM7QUFqTTlDLE1BQUFDLEtBQUE7QUFrTUksTUFBSSxpQkFBaUI7QUFDakI7QUFBQSxFQUNKO0FBRUEsb0JBQWtCO0FBRWxCLGlCQUFlO0FBRWYsUUFBTSxtQkFBbUIsSUFBSSxpQkFBaUIsY0FBYztBQUM1RCxtQkFBaUIsUUFBUSxTQUFTLGlCQUFpQjtBQUFBLElBQy9DLFlBQVk7QUFBQSxJQUNaLFdBQVc7QUFBQSxJQUNYLFNBQVM7QUFBQSxFQUNiLENBQUM7QUFFRCxTQUFPLGlCQUFpQixVQUFVLGNBQWM7QUFDaEQsU0FBTyxpQkFBaUIsVUFBVSxnQkFBZ0IsSUFBSTtBQUN0RCxHQUFBQSxNQUFBLE9BQU8sbUJBQVAsZ0JBQUFBLElBQXVCLGlCQUFpQixVQUFVO0FBQ2xELGVBQU8sbUJBQVAsbUJBQXVCLGlCQUFpQixVQUFVO0FBQ3REO0FBRUEsU0FBUyxrQ0FBMkM7QUF2TnBELE1BQUFBLEtBQUE7QUF3TkksUUFBTSxNQUFLQSxNQUFBLE9BQU8sT0FBTyxnQkFBZCxnQkFBQUEsSUFBMkI7QUFDdEMsTUFBSSxPQUFPLFFBQVc7QUFDbEIsV0FBTztBQUFBLEVBQ1g7QUFFQSxRQUFNLFdBQVUsWUFBTyxPQUFPLFVBQWQsbUJBQXFCO0FBQ3JDLE1BQUksT0FBTyxXQUFXO0FBQ2xCLFFBQUksWUFBWSxNQUFNO0FBQ2xCLGdCQUFVLDRCQUE0QjtBQUFBLElBQzFDO0FBQ0EsV0FBTztBQUFBLEVBQ1g7QUFFQSxTQUFPO0FBQ1g7QUFFQSxJQUFJLFVBQVUsQ0FBQyxnQ0FBZ0MsR0FBRztBQUM5QyxTQUFPLGlCQUFpQix5QkFBeUIsaUNBQWlDLEVBQUUsTUFBTSxLQUFLLENBQUM7QUFDcEc7OztBQzFPQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFXQSxJQUFNQyxRQUFPLGlCQUFpQixZQUFZLFdBQVc7QUFFckQsSUFBTUMsY0FBYTtBQUNuQixJQUFNQyxjQUFhO0FBQ25CLElBQU0sYUFBYTtBQUtaLFNBQVMsT0FBc0I7QUFDbEMsU0FBT0YsTUFBS0MsV0FBVTtBQUMxQjtBQUtPLFNBQVMsT0FBc0I7QUFDbEMsU0FBT0QsTUFBS0UsV0FBVTtBQUMxQjtBQUtPLFNBQVMsT0FBc0I7QUFDbEMsU0FBT0YsTUFBSyxVQUFVO0FBQzFCOzs7QUNwQ0E7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ3dCQSxJQUFJLFVBQVUsU0FBUyxVQUFVO0FBQ2pDLElBQUksZUFBb0QsT0FBTyxZQUFZLFlBQVksWUFBWSxRQUFRLFFBQVE7QUFDbkgsSUFBSTtBQUNKLElBQUk7QUFDSixJQUFJLE9BQU8saUJBQWlCLGNBQWMsT0FBTyxPQUFPLG1CQUFtQixZQUFZO0FBQ25GLE1BQUk7QUFDQSxtQkFBZSxPQUFPLGVBQWUsQ0FBQyxHQUFHLFVBQVU7QUFBQSxNQUMvQyxLQUFLLFdBQVk7QUFDYixjQUFNO0FBQUEsTUFDVjtBQUFBLElBQ0osQ0FBQztBQUNELHVCQUFtQixDQUFDO0FBRXBCLGlCQUFhLFdBQVk7QUFBRSxZQUFNO0FBQUEsSUFBSSxHQUFHLE1BQU0sWUFBWTtBQUFBLEVBQzlELFNBQVMsR0FBRztBQUNSLFFBQUksTUFBTSxrQkFBa0I7QUFDeEIscUJBQWU7QUFBQSxJQUNuQjtBQUFBLEVBQ0o7QUFDSixPQUFPO0FBQ0gsaUJBQWU7QUFDbkI7QUFFQSxJQUFJLG1CQUFtQjtBQUN2QixJQUFJLGVBQWUsU0FBUyxtQkFBbUIsT0FBcUI7QUFDaEUsTUFBSTtBQUNBLFFBQUksUUFBUSxRQUFRLEtBQUssS0FBSztBQUM5QixXQUFPLGlCQUFpQixLQUFLLEtBQUs7QUFBQSxFQUN0QyxTQUFTLEdBQUc7QUFDUixXQUFPO0FBQUEsRUFDWDtBQUNKO0FBRUEsSUFBSSxvQkFBb0IsU0FBUyxpQkFBaUIsT0FBcUI7QUFDbkUsTUFBSTtBQUNBLFFBQUksYUFBYSxLQUFLLEdBQUc7QUFBRSxhQUFPO0FBQUEsSUFBTztBQUN6QyxZQUFRLEtBQUssS0FBSztBQUNsQixXQUFPO0FBQUEsRUFDWCxTQUFTLEdBQUc7QUFDUixXQUFPO0FBQUEsRUFDWDtBQUNKO0FBQ0EsSUFBSSxRQUFRLE9BQU8sVUFBVTtBQUM3QixJQUFJLGNBQWM7QUFDbEIsSUFBSSxVQUFVO0FBQ2QsSUFBSSxXQUFXO0FBQ2YsSUFBSSxXQUFXO0FBQ2YsSUFBSSxZQUFZO0FBQ2hCLElBQUksWUFBWTtBQUNoQixJQUFJLGlCQUFpQixPQUFPLFdBQVcsY0FBYyxDQUFDLENBQUMsT0FBTztBQUU5RCxJQUFJLFNBQVMsRUFBRSxLQUFLLENBQUMsQ0FBQztBQUV0QixJQUFJLFFBQWlDLFNBQVMsbUJBQW1CO0FBQUUsU0FBTztBQUFPO0FBQ2pGLElBQUksT0FBTyxhQUFhLFVBQVU7QUFFMUIsUUFBTSxTQUFTO0FBQ25CLE1BQUksTUFBTSxLQUFLLEdBQUcsTUFBTSxNQUFNLEtBQUssU0FBUyxHQUFHLEdBQUc7QUFDOUMsWUFBUSxTQUFTRyxrQkFBaUIsT0FBTztBQUdyQyxXQUFLLFVBQVUsQ0FBQyxXQUFXLE9BQU8sVUFBVSxlQUFlLE9BQU8sVUFBVSxXQUFXO0FBQ25GLFlBQUk7QUFDQSxjQUFJLE1BQU0sTUFBTSxLQUFLLEtBQUs7QUFDMUIsa0JBQ0ksUUFBUSxZQUNMLFFBQVEsYUFDUixRQUFRLGFBQ1IsUUFBUSxnQkFDVixNQUFNLEVBQUUsS0FBSztBQUFBLFFBQ3RCLFNBQVMsR0FBRztBQUFBLFFBQU87QUFBQSxNQUN2QjtBQUNBLGFBQU87QUFBQSxJQUNYO0FBQUEsRUFDSjtBQUNKO0FBbkJRO0FBcUJSLFNBQVMsbUJBQXNCLE9BQXVEO0FBQ2xGLE1BQUksTUFBTSxLQUFLLEdBQUc7QUFBRSxXQUFPO0FBQUEsRUFBTTtBQUNqQyxNQUFJLENBQUMsT0FBTztBQUFFLFdBQU87QUFBQSxFQUFPO0FBQzVCLE1BQUksT0FBTyxVQUFVLGNBQWMsT0FBTyxVQUFVLFVBQVU7QUFBRSxXQUFPO0FBQUEsRUFBTztBQUM5RSxNQUFJO0FBQ0EsSUFBQyxhQUFxQixPQUFPLE1BQU0sWUFBWTtBQUFBLEVBQ25ELFNBQVMsR0FBRztBQUNSLFFBQUksTUFBTSxrQkFBa0I7QUFBRSxhQUFPO0FBQUEsSUFBTztBQUFBLEVBQ2hEO0FBQ0EsU0FBTyxDQUFDLGFBQWEsS0FBSyxLQUFLLGtCQUFrQixLQUFLO0FBQzFEO0FBRUEsU0FBUyxxQkFBd0IsT0FBc0Q7QUFDbkYsTUFBSSxNQUFNLEtBQUssR0FBRztBQUFFLFdBQU87QUFBQSxFQUFNO0FBQ2pDLE1BQUksQ0FBQyxPQUFPO0FBQUUsV0FBTztBQUFBLEVBQU87QUFDNUIsTUFBSSxPQUFPLFVBQVUsY0FBYyxPQUFPLFVBQVUsVUFBVTtBQUFFLFdBQU87QUFBQSxFQUFPO0FBQzlFLE1BQUksZ0JBQWdCO0FBQUUsV0FBTyxrQkFBa0IsS0FBSztBQUFBLEVBQUc7QUFDdkQsTUFBSSxhQUFhLEtBQUssR0FBRztBQUFFLFdBQU87QUFBQSxFQUFPO0FBQ3pDLE1BQUksV0FBVyxNQUFNLEtBQUssS0FBSztBQUMvQixNQUFJLGFBQWEsV0FBVyxhQUFhLFlBQVksQ0FBRSxpQkFBa0IsS0FBSyxRQUFRLEdBQUc7QUFBRSxXQUFPO0FBQUEsRUFBTztBQUN6RyxTQUFPLGtCQUFrQixLQUFLO0FBQ2xDO0FBRUEsSUFBTyxtQkFBUSxlQUFlLHFCQUFxQjs7O0FDekc1QyxJQUFNLGNBQU4sY0FBMEIsTUFBTTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU1uQyxZQUFZLFNBQWtCLFNBQXdCO0FBQ2xELFVBQU0sU0FBUyxPQUFPO0FBQ3RCLFNBQUssT0FBTztBQUFBLEVBQ2hCO0FBQ0o7QUFjTyxJQUFNLDBCQUFOLGNBQXNDLE1BQU07QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBYS9DLFlBQVksU0FBc0MsUUFBYyxNQUFlO0FBQzNFLFdBQU8sc0JBQVEsK0NBQStDLGNBQWMsYUFBYSxNQUFNLEdBQUcsRUFBRSxPQUFPLE9BQU8sQ0FBQztBQUNuSCxTQUFLLFVBQVU7QUFDZixTQUFLLE9BQU87QUFBQSxFQUNoQjtBQUNKO0FBK0JBLElBQU0sYUFBYSx1QkFBTyxTQUFTO0FBQ25DLElBQU0sZ0JBQWdCLHVCQUFPLFlBQVk7QUE3RnpDLElBQUFDO0FBOEZBLElBQU0sV0FBaUNBLE1BQUEsT0FBTyxZQUFQLE9BQUFBLE1BQWtCLHVCQUFPLGlCQUFpQjtBQW9EMUUsSUFBTSxxQkFBTixNQUFNLDRCQUE4QixRQUFnRTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQXVDdkcsWUFBWSxVQUF5QyxhQUEyQztBQUM1RixRQUFJO0FBQ0osUUFBSTtBQUNKLFVBQU0sQ0FBQyxLQUFLLFFBQVE7QUFBRSxnQkFBVTtBQUFLLGVBQVM7QUFBQSxJQUFLLENBQUM7QUFFcEQsUUFBSyxLQUFLLFlBQW9CLE9BQU8sTUFBTSxTQUFTO0FBQ2hELFlBQU0sSUFBSSxVQUFVLG1JQUFtSTtBQUFBLElBQzNKO0FBRUEsUUFBSSxVQUE4QztBQUFBLE1BQzlDLFNBQVM7QUFBQSxNQUNUO0FBQUEsTUFDQTtBQUFBLE1BQ0EsSUFBSSxjQUFjO0FBQUUsZUFBTyxvQ0FBZTtBQUFBLE1BQU07QUFBQSxNQUNoRCxJQUFJLFlBQVksSUFBSTtBQUFFLHNCQUFjLGtCQUFNO0FBQUEsTUFBVztBQUFBLElBQ3pEO0FBRUEsVUFBTSxRQUFpQztBQUFBLE1BQ25DLElBQUksT0FBTztBQUFFLGVBQU87QUFBQSxNQUFPO0FBQUEsTUFDM0IsV0FBVztBQUFBLE1BQ1gsU0FBUztBQUFBLElBQ2I7QUFHQSxTQUFLLE9BQU8saUJBQWlCLE1BQU07QUFBQSxNQUMvQixDQUFDLFVBQVUsR0FBRztBQUFBLFFBQ1YsY0FBYztBQUFBLFFBQ2QsWUFBWTtBQUFBLFFBQ1osVUFBVTtBQUFBLFFBQ1YsT0FBTztBQUFBLE1BQ1g7QUFBQSxNQUNBLENBQUMsYUFBYSxHQUFHO0FBQUEsUUFDYixjQUFjO0FBQUEsUUFDZCxZQUFZO0FBQUEsUUFDWixVQUFVO0FBQUEsUUFDVixPQUFPLGFBQWEsU0FBUyxLQUFLO0FBQUEsTUFDdEM7QUFBQSxJQUNKLENBQUM7QUFHRCxVQUFNLFdBQVcsWUFBWSxTQUFTLEtBQUs7QUFDM0MsUUFBSTtBQUNBLGVBQVMsWUFBWSxTQUFTLEtBQUssR0FBRyxRQUFRO0FBQUEsSUFDbEQsU0FBUyxLQUFLO0FBQ1YsVUFBSSxNQUFNLFdBQVc7QUFDakIsZ0JBQVEsSUFBSSx1REFBdUQsR0FBRztBQUFBLE1BQzFFLE9BQU87QUFDSCxpQkFBUyxHQUFHO0FBQUEsTUFDaEI7QUFBQSxJQUNKO0FBQUEsRUFDSjtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUF5REEsT0FBTyxPQUF1QztBQUMxQyxXQUFPLElBQUksb0JBQXlCLENBQUMsWUFBWTtBQUc3QyxjQUFRLElBQUk7QUFBQSxRQUNSLEtBQUssYUFBYSxFQUFFLElBQUksWUFBWSxzQkFBc0IsRUFBRSxNQUFNLENBQUMsQ0FBQztBQUFBLFFBQ3BFLGVBQWUsSUFBSTtBQUFBLE1BQ3ZCLENBQUMsRUFBRSxLQUFLLE1BQU0sUUFBUSxHQUFHLE1BQU0sUUFBUSxDQUFDO0FBQUEsSUFDNUMsQ0FBQztBQUFBLEVBQ0w7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBMkJBLFNBQVMsUUFBNEM7QUFDakQsUUFBSSxPQUFPLFNBQVM7QUFDaEIsV0FBSyxLQUFLLE9BQU8sT0FBTyxNQUFNO0FBQUEsSUFDbEMsT0FBTztBQUNILGFBQU8saUJBQWlCLFNBQVMsTUFBTSxLQUFLLEtBQUssT0FBTyxPQUFPLE1BQU0sR0FBRyxFQUFDLFNBQVMsS0FBSSxDQUFDO0FBQUEsSUFDM0Y7QUFFQSxXQUFPO0FBQUEsRUFDWDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQStCQSxLQUFxQyxhQUFzSCxZQUF3SCxhQUFvRjtBQUNuVyxRQUFJLEVBQUUsZ0JBQWdCLHNCQUFxQjtBQUN2QyxZQUFNLElBQUksVUFBVSxnRUFBZ0U7QUFBQSxJQUN4RjtBQU1BLFFBQUksQ0FBQyxpQkFBVyxXQUFXLEdBQUc7QUFBRSxvQkFBYztBQUFBLElBQWlCO0FBQy9ELFFBQUksQ0FBQyxpQkFBVyxVQUFVLEdBQUc7QUFBRSxtQkFBYTtBQUFBLElBQVM7QUFFckQsUUFBSSxnQkFBZ0IsWUFBWSxjQUFjLFNBQVM7QUFFbkQsYUFBTyxJQUFJLG9CQUFtQixDQUFDLFlBQVksUUFBUSxJQUFXLENBQUM7QUFBQSxJQUNuRTtBQUVBLFVBQU0sVUFBK0MsQ0FBQztBQUN0RCxTQUFLLFVBQVUsSUFBSTtBQUVuQixXQUFPLElBQUksb0JBQXdDLENBQUMsU0FBUyxXQUFXO0FBQ3BFLFdBQUssTUFBTTtBQUFBLFFBQ1AsQ0FBQyxVQUFVO0FBclkzQixjQUFBQTtBQXNZb0IsY0FBSSxLQUFLLFVBQVUsTUFBTSxTQUFTO0FBQUUsaUJBQUssVUFBVSxJQUFJO0FBQUEsVUFBTTtBQUM3RCxXQUFBQSxNQUFBLFFBQVEsWUFBUixnQkFBQUEsSUFBQTtBQUVBLGNBQUk7QUFDQSxvQkFBUSxZQUFhLEtBQUssQ0FBQztBQUFBLFVBQy9CLFNBQVMsS0FBSztBQUNWLG1CQUFPLEdBQUc7QUFBQSxVQUNkO0FBQUEsUUFDSjtBQUFBLFFBQ0EsQ0FBQyxXQUFZO0FBL1k3QixjQUFBQTtBQWdab0IsY0FBSSxLQUFLLFVBQVUsTUFBTSxTQUFTO0FBQUUsaUJBQUssVUFBVSxJQUFJO0FBQUEsVUFBTTtBQUM3RCxXQUFBQSxNQUFBLFFBQVEsWUFBUixnQkFBQUEsSUFBQTtBQUVBLGNBQUk7QUFDQSxvQkFBUSxXQUFZLE1BQU0sQ0FBQztBQUFBLFVBQy9CLFNBQVMsS0FBSztBQUNWLG1CQUFPLEdBQUc7QUFBQSxVQUNkO0FBQUEsUUFDSjtBQUFBLE1BQ0o7QUFBQSxJQUNKLEdBQUcsT0FBTyxVQUFXO0FBRWpCLFVBQUk7QUFDQSxlQUFPLDJDQUFjO0FBQUEsTUFDekIsVUFBRTtBQUNFLGNBQU0sS0FBSyxPQUFPLEtBQUs7QUFBQSxNQUMzQjtBQUFBLElBQ0osQ0FBQztBQUFBLEVBQ0w7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUErQkEsTUFBdUIsWUFBcUYsYUFBNEU7QUFDcEwsV0FBTyxLQUFLLEtBQUssUUFBVyxZQUFZLFdBQVc7QUFBQSxFQUN2RDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFpQ0EsUUFBUSxXQUE2QyxhQUFrRTtBQUNuSCxRQUFJLEVBQUUsZ0JBQWdCLHNCQUFxQjtBQUN2QyxZQUFNLElBQUksVUFBVSxtRUFBbUU7QUFBQSxJQUMzRjtBQUVBLFFBQUksQ0FBQyxpQkFBVyxTQUFTLEdBQUc7QUFDeEIsYUFBTyxLQUFLLEtBQUssV0FBVyxXQUFXLFdBQVc7QUFBQSxJQUN0RDtBQUVBLFdBQU8sS0FBSztBQUFBLE1BQ1IsQ0FBQyxVQUFVLG9CQUFtQixRQUFRLFVBQVUsQ0FBQyxFQUFFLEtBQUssTUFBTSxLQUFLO0FBQUEsTUFDbkUsQ0FBQyxXQUFZLG9CQUFtQixRQUFRLFVBQVUsQ0FBQyxFQUFFLEtBQUssTUFBTTtBQUFFLGNBQU07QUFBQSxNQUFRLENBQUM7QUFBQSxNQUNqRjtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQVlBLGFBeldTLFlBRVMsZUF1V04sUUFBTyxJQUFJO0FBQ25CLFdBQU87QUFBQSxFQUNYO0FBQUEsRUFhQSxPQUFPLElBQXNELFFBQXdDO0FBQ2pHLFFBQUksWUFBWSxNQUFNLEtBQUssTUFBTTtBQUNqQyxVQUFNLFVBQVUsVUFBVSxXQUFXLElBQy9CLG9CQUFtQixRQUFRLFNBQVMsSUFDcEMsSUFBSSxvQkFBNEIsQ0FBQyxTQUFTLFdBQVc7QUFDbkQsV0FBSyxRQUFRLElBQUksU0FBUyxFQUFFLEtBQUssU0FBUyxNQUFNO0FBQUEsSUFDcEQsR0FBRyxDQUFDLFVBQTBCLFVBQVUsU0FBUyxXQUFXLEtBQUssQ0FBQztBQUN0RSxXQUFPO0FBQUEsRUFDWDtBQUFBLEVBYUEsT0FBTyxXQUE2RCxRQUF3QztBQUN4RyxRQUFJLFlBQVksTUFBTSxLQUFLLE1BQU07QUFDakMsVUFBTSxVQUFVLFVBQVUsV0FBVyxJQUMvQixvQkFBbUIsUUFBUSxTQUFTLElBQ3BDLElBQUksb0JBQTRCLENBQUMsU0FBUyxXQUFXO0FBQ25ELFdBQUssUUFBUSxXQUFXLFNBQVMsRUFBRSxLQUFLLFNBQVMsTUFBTTtBQUFBLElBQzNELEdBQUcsQ0FBQyxVQUEwQixVQUFVLFNBQVMsV0FBVyxLQUFLLENBQUM7QUFDdEUsV0FBTztBQUFBLEVBQ1g7QUFBQSxFQWVBLE9BQU8sSUFBc0QsUUFBd0M7QUFDakcsUUFBSSxZQUFZLE1BQU0sS0FBSyxNQUFNO0FBQ2pDLFVBQU0sVUFBVSxVQUFVLFdBQVcsSUFDL0Isb0JBQW1CLFFBQVEsU0FBUyxJQUNwQyxJQUFJLG9CQUE0QixDQUFDLFNBQVMsV0FBVztBQUNuRCxXQUFLLFFBQVEsSUFBSSxTQUFTLEVBQUUsS0FBSyxTQUFTLE1BQU07QUFBQSxJQUNwRCxHQUFHLENBQUMsVUFBMEIsVUFBVSxTQUFTLFdBQVcsS0FBSyxDQUFDO0FBQ3RFLFdBQU87QUFBQSxFQUNYO0FBQUEsRUFZQSxPQUFPLEtBQXVELFFBQXdDO0FBQ2xHLFFBQUksWUFBWSxNQUFNLEtBQUssTUFBTTtBQUNqQyxVQUFNLFVBQVUsSUFBSSxvQkFBNEIsQ0FBQyxTQUFTLFdBQVc7QUFDakUsV0FBSyxRQUFRLEtBQUssU0FBUyxFQUFFLEtBQUssU0FBUyxNQUFNO0FBQUEsSUFDckQsR0FBRyxDQUFDLFVBQTBCLFVBQVUsU0FBUyxXQUFXLEtBQUssQ0FBQztBQUNsRSxXQUFPO0FBQUEsRUFDWDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLE9BQU8sT0FBa0IsT0FBb0M7QUFDekQsVUFBTSxJQUFJLElBQUksb0JBQXNCLE1BQU07QUFBQSxJQUFDLENBQUM7QUFDNUMsTUFBRSxPQUFPLEtBQUs7QUFDZCxXQUFPO0FBQUEsRUFDWDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsRUFZQSxPQUFPLFFBQW1CLGNBQXNCLE9BQW9DO0FBQ2hGLFVBQU0sVUFBVSxJQUFJLG9CQUFzQixNQUFNO0FBQUEsSUFBQyxDQUFDO0FBQ2xELFFBQUksZUFBZSxPQUFPLGdCQUFnQixjQUFjLFlBQVksV0FBVyxPQUFPLFlBQVksWUFBWSxZQUFZO0FBQ3RILGtCQUFZLFFBQVEsWUFBWSxFQUFFLGlCQUFpQixTQUFTLE1BQU0sS0FBSyxRQUFRLE9BQU8sS0FBSyxDQUFDO0FBQUEsSUFDaEcsT0FBTztBQUNILGlCQUFXLE1BQU0sS0FBSyxRQUFRLE9BQU8sS0FBSyxHQUFHLFlBQVk7QUFBQSxJQUM3RDtBQUNBLFdBQU87QUFBQSxFQUNYO0FBQUEsRUFpQkEsT0FBTyxNQUFnQixjQUFzQixPQUFrQztBQUMzRSxXQUFPLElBQUksb0JBQXNCLENBQUMsWUFBWTtBQUMxQyxpQkFBVyxNQUFNLFFBQVEsS0FBTSxHQUFHLFlBQVk7QUFBQSxJQUNsRCxDQUFDO0FBQUEsRUFDTDtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQSxFQU9BLE9BQU8sT0FBa0IsUUFBcUM7QUFDMUQsV0FBTyxJQUFJLG9CQUFzQixDQUFDLEdBQUcsV0FBVyxPQUFPLE1BQU0sQ0FBQztBQUFBLEVBQ2xFO0FBQUEsRUFvQkEsT0FBTyxRQUFrQixPQUE0RDtBQUNqRixRQUFJLGlCQUFpQixxQkFBb0I7QUFFckMsYUFBTztBQUFBLElBQ1g7QUFDQSxXQUFPLElBQUksb0JBQXdCLENBQUMsWUFBWSxRQUFRLEtBQUssQ0FBQztBQUFBLEVBQ2xFO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBVUEsT0FBTyxnQkFBdUQ7QUFDMUQsUUFBSSxTQUE2QyxFQUFFLGFBQWEsS0FBSztBQUNyRSxXQUFPLFVBQVUsSUFBSSxvQkFBc0IsQ0FBQyxTQUFTLFdBQVc7QUFDNUQsYUFBTyxVQUFVO0FBQ2pCLGFBQU8sU0FBUztBQUFBLElBQ3BCLEdBQUcsQ0FBQyxVQUFnQjtBQXpyQjVCLFVBQUFBO0FBeXJCOEIsT0FBQUEsTUFBQSxPQUFPLGdCQUFQLGdCQUFBQSxJQUFBLGFBQXFCO0FBQUEsSUFBUSxDQUFDO0FBQ3BELFdBQU87QUFBQSxFQUNYO0FBQ0o7QUFNQSxTQUFTLGFBQWdCLFNBQTZDLE9BQWdDO0FBQ2xHLE1BQUksc0JBQWdEO0FBRXBELFNBQU8sQ0FBQyxXQUFrRDtBQUN0RCxRQUFJLENBQUMsTUFBTSxTQUFTO0FBQ2hCLFlBQU0sVUFBVTtBQUNoQixZQUFNLFNBQVM7QUFDZixjQUFRLE9BQU8sTUFBTTtBQU1yQixXQUFLLFFBQVEsVUFBVSxLQUFLLEtBQUssUUFBUSxTQUFTLFFBQVcsQ0FBQyxRQUFRO0FBQ2xFLFlBQUksUUFBUSxRQUFRO0FBQ2hCLGdCQUFNO0FBQUEsUUFDVjtBQUFBLE1BQ0osQ0FBQztBQUFBLElBQ0w7QUFJQSxRQUFJLENBQUMsTUFBTSxVQUFVLENBQUMsUUFBUSxhQUFhO0FBQUU7QUFBQSxJQUFRO0FBRXJELDBCQUFzQixJQUFJLFFBQWMsQ0FBQyxZQUFZO0FBQ2pELFVBQUk7QUFDQSxnQkFBUSxRQUFRLFlBQWEsTUFBTSxPQUFRLEtBQUssQ0FBQztBQUFBLE1BQ3JELFNBQVMsS0FBSztBQUNWLGdCQUFRLE9BQU8sSUFBSSx3QkFBd0IsUUFBUSxTQUFTLEtBQUssOENBQThDLENBQUM7QUFBQSxNQUNwSDtBQUFBLElBQ0osQ0FBQyxFQUFFLE1BQU0sQ0FBQ0MsWUFBWTtBQUNsQixjQUFRLE9BQU8sSUFBSSx3QkFBd0IsUUFBUSxTQUFTQSxTQUFRLDhDQUE4QyxDQUFDO0FBQUEsSUFDdkgsQ0FBQztBQUdELFlBQVEsY0FBYztBQUV0QixXQUFPO0FBQUEsRUFDWDtBQUNKO0FBS0EsU0FBUyxZQUFlLFNBQTZDLE9BQStEO0FBQ2hJLFNBQU8sQ0FBQyxVQUFVO0FBQ2QsUUFBSSxNQUFNLFdBQVc7QUFBRTtBQUFBLElBQVE7QUFDL0IsVUFBTSxZQUFZO0FBRWxCLFFBQUksVUFBVSxRQUFRLFNBQVM7QUFDM0IsVUFBSSxNQUFNLFNBQVM7QUFBRTtBQUFBLE1BQVE7QUFDN0IsWUFBTSxVQUFVO0FBQ2hCLGNBQVEsT0FBTyxJQUFJLFVBQVUsMkNBQTJDLENBQUM7QUFDekU7QUFBQSxJQUNKO0FBRUEsUUFBSSxTQUFTLFNBQVMsT0FBTyxVQUFVLFlBQVksT0FBTyxVQUFVLGFBQWE7QUFDN0UsVUFBSTtBQUNKLFVBQUk7QUFDQSxlQUFRLE1BQWM7QUFBQSxNQUMxQixTQUFTLEtBQUs7QUFDVixjQUFNLFVBQVU7QUFDaEIsZ0JBQVEsT0FBTyxHQUFHO0FBQ2xCO0FBQUEsTUFDSjtBQUVBLFVBQUksaUJBQVcsSUFBSSxHQUFHO0FBQ2xCLFlBQUk7QUFDQSxjQUFJLFNBQVUsTUFBYztBQUM1QixjQUFJLGlCQUFXLE1BQU0sR0FBRztBQUNwQixrQkFBTSxjQUFjLENBQUMsVUFBZ0I7QUFDakMsc0JBQVEsTUFBTSxRQUFRLE9BQU8sQ0FBQyxLQUFLLENBQUM7QUFBQSxZQUN4QztBQUNBLGdCQUFJLE1BQU0sUUFBUTtBQUlkLG1CQUFLLGFBQWEsaUNBQUssVUFBTCxFQUFjLFlBQVksSUFBRyxLQUFLLEVBQUUsTUFBTSxNQUFNO0FBQUEsWUFDdEUsT0FBTztBQUNILHNCQUFRLGNBQWM7QUFBQSxZQUMxQjtBQUFBLFVBQ0o7QUFBQSxRQUNKLFNBQVE7QUFBQSxRQUFDO0FBRVQsY0FBTSxXQUFvQztBQUFBLFVBQ3RDLE1BQU0sTUFBTTtBQUFBLFVBQ1osV0FBVztBQUFBLFVBQ1gsSUFBSSxVQUFVO0FBQUUsbUJBQU8sS0FBSyxLQUFLO0FBQUEsVUFBUTtBQUFBLFVBQ3pDLElBQUksUUFBUUMsUUFBTztBQUFFLGlCQUFLLEtBQUssVUFBVUE7QUFBQSxVQUFPO0FBQUEsVUFDaEQsSUFBSSxTQUFTO0FBQUUsbUJBQU8sS0FBSyxLQUFLO0FBQUEsVUFBTztBQUFBLFFBQzNDO0FBRUEsY0FBTSxXQUFXLFlBQVksU0FBUyxRQUFRO0FBQzlDLFlBQUk7QUFDQSxrQkFBUSxNQUFNLE1BQU0sT0FBTyxDQUFDLFlBQVksU0FBUyxRQUFRLEdBQUcsUUFBUSxDQUFDO0FBQUEsUUFDekUsU0FBUyxLQUFLO0FBQ1YsbUJBQVMsR0FBRztBQUFBLFFBQ2hCO0FBQ0E7QUFBQSxNQUNKO0FBQUEsSUFDSjtBQUVBLFFBQUksTUFBTSxTQUFTO0FBQUU7QUFBQSxJQUFRO0FBQzdCLFVBQU0sVUFBVTtBQUNoQixZQUFRLFFBQVEsS0FBSztBQUFBLEVBQ3pCO0FBQ0o7QUFLQSxTQUFTLFlBQWUsU0FBNkMsT0FBNEQ7QUFDN0gsU0FBTyxDQUFDLFdBQVk7QUFDaEIsUUFBSSxNQUFNLFdBQVc7QUFBRTtBQUFBLElBQVE7QUFDL0IsVUFBTSxZQUFZO0FBRWxCLFFBQUksTUFBTSxTQUFTO0FBQ2YsVUFBSTtBQUNBLFlBQUksa0JBQWtCLGVBQWUsTUFBTSxrQkFBa0IsZUFBZSxPQUFPLEdBQUcsT0FBTyxPQUFPLE1BQU0sT0FBTyxLQUFLLEdBQUc7QUFFckg7QUFBQSxRQUNKO0FBQUEsTUFDSixTQUFRO0FBQUEsTUFBQztBQUVULFdBQUssUUFBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsU0FBUyxNQUFNLENBQUM7QUFBQSxJQUM1RSxPQUFPO0FBQ0gsWUFBTSxVQUFVO0FBQ2hCLGNBQVEsT0FBTyxNQUFNO0FBQUEsSUFDekI7QUFBQSxFQUNKO0FBQ0o7QUFNQSxTQUFTLFVBQVUsUUFBcUMsUUFBZSxPQUE0QjtBQUMvRixRQUFNLFVBQTJCLENBQUM7QUFFbEMsYUFBVyxTQUFTLFFBQVE7QUFDeEIsUUFBSTtBQUNKLFFBQUk7QUFDQSxVQUFJLENBQUMsaUJBQVcsTUFBTSxJQUFJLEdBQUc7QUFBRTtBQUFBLE1BQVU7QUFDekMsZUFBUyxNQUFNO0FBQ2YsVUFBSSxDQUFDLGlCQUFXLE1BQU0sR0FBRztBQUFFO0FBQUEsTUFBVTtBQUFBLElBQ3pDLFNBQVE7QUFBRTtBQUFBLElBQVU7QUFFcEIsUUFBSTtBQUNKLFFBQUk7QUFDQSxlQUFTLFFBQVEsTUFBTSxRQUFRLE9BQU8sQ0FBQyxLQUFLLENBQUM7QUFBQSxJQUNqRCxTQUFTLEtBQUs7QUFDVixjQUFRLE9BQU8sSUFBSSx3QkFBd0IsUUFBUSxLQUFLLHVDQUF1QyxDQUFDO0FBQ2hHO0FBQUEsSUFDSjtBQUVBLFFBQUksQ0FBQyxRQUFRO0FBQUU7QUFBQSxJQUFVO0FBQ3pCLFlBQVE7QUFBQSxPQUNILGtCQUFrQixVQUFXLFNBQVMsUUFBUSxRQUFRLE1BQU0sR0FBRyxNQUFNLENBQUMsV0FBWTtBQUMvRSxnQkFBUSxPQUFPLElBQUksd0JBQXdCLFFBQVEsUUFBUSx1Q0FBdUMsQ0FBQztBQUFBLE1BQ3ZHLENBQUM7QUFBQSxJQUNMO0FBQUEsRUFDSjtBQUVBLFNBQU8sUUFBUSxJQUFJLE9BQU87QUFDOUI7QUFLQSxTQUFTLFNBQVksR0FBUztBQUMxQixTQUFPO0FBQ1g7QUFLQSxTQUFTLFFBQVEsUUFBcUI7QUFDbEMsUUFBTTtBQUNWO0FBS0EsU0FBUyxhQUFhLEtBQWtCO0FBQ3BDLE1BQUk7QUFDQSxRQUFJLGVBQWUsU0FBUyxPQUFPLFFBQVEsWUFBWSxJQUFJLGFBQWEsT0FBTyxVQUFVLFVBQVU7QUFDL0YsYUFBTyxLQUFLO0FBQUEsSUFDaEI7QUFBQSxFQUNKLFNBQVE7QUFBQSxFQUFDO0FBRVQsTUFBSTtBQUNBLFdBQU8sS0FBSyxVQUFVLEdBQUc7QUFBQSxFQUM3QixTQUFRO0FBQUEsRUFBQztBQUVULE1BQUk7QUFDQSxXQUFPLE9BQU8sVUFBVSxTQUFTLEtBQUssR0FBRztBQUFBLEVBQzdDLFNBQVE7QUFBQSxFQUFDO0FBRVQsU0FBTztBQUNYO0FBS0EsU0FBUyxlQUFrQixTQUErQztBQTk0QjFFLE1BQUFGO0FBKzRCSSxNQUFJLE9BQTJDQSxNQUFBLFFBQVEsVUFBVSxNQUFsQixPQUFBQSxNQUF1QixDQUFDO0FBQ3ZFLE1BQUksRUFBRSxhQUFhLE1BQU07QUFDckIsV0FBTyxPQUFPLEtBQUsscUJBQTJCLENBQUM7QUFBQSxFQUNuRDtBQUNBLE1BQUksUUFBUSxVQUFVLEtBQUssTUFBTTtBQUM3QixRQUFJLFFBQVM7QUFDYixZQUFRLFVBQVUsSUFBSTtBQUFBLEVBQzFCO0FBQ0EsU0FBTyxJQUFJO0FBQ2Y7QUFHQSxJQUFJLHVCQUF1QixRQUFRO0FBQ25DLElBQUksd0JBQXdCLE9BQU8seUJBQXlCLFlBQVk7QUFDcEUseUJBQXVCLHFCQUFxQixLQUFLLE9BQU87QUFDNUQsT0FBTztBQUNILHlCQUF1QixXQUF3QztBQUMzRCxRQUFJO0FBQ0osUUFBSTtBQUNKLFVBQU0sVUFBVSxJQUFJLFFBQVcsQ0FBQyxLQUFLLFFBQVE7QUFBRSxnQkFBVTtBQUFLLGVBQVM7QUFBQSxJQUFLLENBQUM7QUFDN0UsV0FBTyxFQUFFLFNBQVMsU0FBUyxPQUFPO0FBQUEsRUFDdEM7QUFDSjs7O0FGcDVCQSxJQUFJLFFBQVE7QUFDUixTQUFPLFNBQVMsT0FBTyxVQUFVLENBQUM7QUFDdEM7QUFJQSxJQUFNRyxRQUFPLGlCQUFpQixZQUFZLElBQUk7QUFDOUMsSUFBTSxhQUFhLGlCQUFpQixZQUFZLFVBQVU7QUFDMUQsSUFBTSxnQkFBZ0Isb0JBQUksSUFBOEI7QUFFeEQsSUFBTSxjQUFjO0FBQ3BCLElBQU0sZUFBZTtBQWdDckIsU0FBUyxhQUFxQjtBQUMxQixNQUFJO0FBQ0osS0FBRztBQUNDLGFBQVMsT0FBTztBQUFBLEVBQ3BCLFNBQVMsY0FBYyxJQUFJLE1BQU07QUFDakMsU0FBTztBQUNYO0FBY08sU0FBUyxLQUFLLFNBQStDO0FBQ2hFLFFBQU0sS0FBSyxXQUFXO0FBRXRCLFFBQU0sU0FBUyxtQkFBbUIsY0FBbUI7QUFDckQsZ0JBQWMsSUFBSSxJQUFJLEVBQUUsU0FBUyxPQUFPLFNBQVMsUUFBUSxPQUFPLE9BQU8sQ0FBQztBQUV4RSxRQUFNLFVBQVVBLE1BQUssYUFBYSxPQUFPLE9BQU8sRUFBRSxXQUFXLEdBQUcsR0FBRyxPQUFPLENBQUM7QUFDM0UsTUFBSSxVQUFVO0FBRWQsVUFBUSxLQUFLLENBQUMsUUFBUTtBQUNsQixjQUFVO0FBQ1Ysa0JBQWMsT0FBTyxFQUFFO0FBQ3ZCLFdBQU8sUUFBUSxHQUFHO0FBQUEsRUFDdEIsR0FBRyxDQUFDLFFBQVE7QUFDUixjQUFVO0FBQ1Ysa0JBQWMsT0FBTyxFQUFFO0FBQ3ZCLFdBQU8sT0FBTyxHQUFHO0FBQUEsRUFDckIsQ0FBQztBQUVELFFBQU0sU0FBUyxNQUFNO0FBQ2pCLGtCQUFjLE9BQU8sRUFBRTtBQUN2QixXQUFPLFdBQVcsY0FBYyxFQUFDLFdBQVcsR0FBRSxDQUFDLEVBQUUsTUFBTSxDQUFDLFFBQVE7QUFDNUQsY0FBUSxNQUFNLHFEQUFxRCxHQUFHO0FBQUEsSUFDMUUsQ0FBQztBQUFBLEVBQ0w7QUFFQSxTQUFPLGNBQWMsTUFBTTtBQUN2QixRQUFJLFNBQVM7QUFDVCxhQUFPLE9BQU87QUFBQSxJQUNsQixPQUFPO0FBQ0gsYUFBTyxRQUFRLEtBQUssTUFBTTtBQUFBLElBQzlCO0FBQUEsRUFDSjtBQUVBLFNBQU8sT0FBTztBQUNsQjtBQVVPLFNBQVMsT0FBTyxlQUF1QixNQUFzQztBQUNoRixTQUFPLEtBQUssRUFBRSxZQUFZLEtBQUssQ0FBQztBQUNwQztBQVVPLFNBQVMsS0FBSyxhQUFxQixNQUFzQztBQUM1RSxTQUFPLEtBQUssRUFBRSxVQUFVLEtBQUssQ0FBQztBQUNsQzs7O0FHM0lBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFZQSxJQUFNQyxRQUFPLGlCQUFpQixZQUFZLFNBQVM7QUFFbkQsSUFBTSxtQkFBbUI7QUFDekIsSUFBTSxnQkFBZ0I7QUFRZixTQUFTLFFBQVEsTUFBNkI7QUFDakQsU0FBT0EsTUFBSyxrQkFBa0IsRUFBQyxLQUFJLENBQUM7QUFDeEM7QUFPTyxTQUFTLE9BQXdCO0FBQ3BDLFNBQU9BLE1BQUssYUFBYTtBQUM3Qjs7O0FDbENBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQzRDQSxJQUFNLGlCQUFpQjtBQUN2QixJQUFNLHlCQUF5QjtBQUsvQixTQUFTLHNCQUE4QjtBQUNuQyxNQUFJO0FBQ0EsVUFBTSxTQUFTLE9BQU8sS0FBSyxZQUFZLHNCQUFzQjtBQUM3RCxRQUFJLFNBQVMsRUFBRyxRQUFPO0FBQ3ZCLFVBQU0sUUFBUSxPQUFPLE9BQU8sS0FBSyxNQUFNLFNBQVMsdUJBQXVCLE1BQU0sQ0FBQztBQUM5RSxXQUFPLE9BQU8sY0FBYyxLQUFLLEtBQUssUUFBUSxJQUFJLFFBQVE7QUFBQSxFQUM5RCxTQUFRO0FBQ0osV0FBTztBQUFBLEVBQ1g7QUFDSjtBQUVBLFNBQVMsNEJBQTRCLFlBQTBCO0FBQzNELE1BQUk7QUFDQSxVQUFNLFNBQVMsT0FBTyxLQUFLLFlBQVksc0JBQXNCO0FBQzdELFVBQU0sa0JBQWtCLFNBQVMsSUFBSSxPQUFPLE9BQU8sT0FBTyxLQUFLLE1BQU0sR0FBRyxNQUFNO0FBQzlFLFdBQU8sT0FBTyxrQkFBa0IseUJBQXlCO0FBQUEsRUFDN0QsU0FBUTtBQUFBLEVBR1I7QUFDSjtBQU1BLFNBQVMscUJBQTZCO0FBQ2xDLE1BQUksQ0FBQyxPQUFRLFFBQU87QUFRcEIsUUFBTSxTQUFTLE9BQU8sZ0JBQWdCLGVBQWUsT0FBTyxTQUFTLFlBQVksVUFBVSxJQUNyRixZQUFZLGFBQ1osS0FBSyxJQUFJO0FBQ2YsUUFBTSxXQUFXLEtBQUssSUFBSSxHQUFHLEtBQUssTUFBTSxTQUFTLEdBQUksQ0FBQztBQUN0RCxNQUFJLFVBQVUsb0JBQW9CO0FBQ2xDLE1BQUk7QUFDQSxVQUFNLE1BQU0sT0FBTyxlQUFlLFFBQVEsY0FBYztBQUN4RCxVQUFNLFNBQVMsUUFBUSxPQUFPLElBQUksT0FBTyxHQUFHO0FBQzVDLFFBQUksT0FBTyxjQUFjLE1BQU0sS0FBSyxTQUFTLFFBQVMsV0FBVTtBQUFBLEVBQ3BFLFNBQVE7QUFBQSxFQUdSO0FBRUEsUUFBTSxPQUFPLE9BQU8sY0FBYyxPQUFPLEtBQUssV0FBVyxLQUFLLFVBQVUsT0FBTyxtQkFDekUsS0FBSyxJQUFJLFVBQVUsR0FBRyxRQUFRLElBQzlCO0FBQ04sTUFBSTtBQUNBLFdBQU8sZUFBZSxRQUFRLGdCQUFnQixPQUFPLElBQUksQ0FBQztBQUFBLEVBQzlELFNBQVE7QUFBQSxFQUdSO0FBQ0EsOEJBQTRCLElBQUk7QUFDaEMsU0FBTztBQUNYO0FBRUEsU0FBUyxnQkFBK0I7QUFDcEMsUUFBTSxRQUFRLE9BQXNCO0FBQUEsSUFDaEMsU0FBUyxPQUFPO0FBQUEsSUFDaEIsWUFBWSxtQkFBbUI7QUFBQSxJQUMvQixZQUFZO0FBQUEsSUFDWixhQUFhLG9CQUFJLElBQXlCO0FBQUEsSUFDMUMsU0FBUztBQUFBLEVBQ2I7QUFDQSxNQUFJLENBQUMsT0FBUSxRQUFPLE1BQU07QUFDMUIsUUFBTSxJQUFLLE9BQWUsU0FBVSxPQUFlLFVBQVUsQ0FBQztBQUM5RCxTQUFRLEVBQUUsV0FBVyxFQUFFLFlBQVksTUFBTTtBQUM3QztBQUVBLElBQU0sSUFBSSxjQUFjO0FBRXhCLElBQU0sY0FBYztBQUNwQixJQUFNLGlCQUFpQjtBQUN2QixJQUFNLFdBQVc7QUFDakIsSUFBTSxXQUFXO0FBQ2pCLElBQU0sV0FBVztBQUNqQixJQUFNLFlBQVk7QUFDbEIsSUFBTSxrQkFBa0I7QUFDeEIsSUFBTSxrQkFBa0I7QUFDeEIsSUFBTSxZQUFZO0FBRWxCLElBQU0sWUFBWTtBQUNsQixJQUFNLFlBQVk7QUFDbEIsSUFBTSxhQUFhO0FBQ25CLElBQU0sYUFBYTtBQUluQixJQUFNQyxtQkFBa0IsTUFBTTtBQUM5QixJQUFNLG1CQUFtQjtBQU16QixJQUFNLGNBQWM7QUFDcEIsSUFBTSxjQUFjO0FBQ3BCLElBQU0sY0FBYyxJQUFJLFlBQVk7QUFFcEMsU0FBUyxVQUFVLFVBQTBCO0FBQ3pDLFNBQU8sT0FBTyxTQUFTLFNBQVMsbUJBQW1CO0FBQ3ZEO0FBR0EsSUFBTSx5QkFBTixjQUFxQyxNQUFNO0FBQUEsRUFDdkMsY0FBYztBQUNWLFVBQU0sMEJBQTBCO0FBQ2hDLFNBQUssT0FBTztBQUFBLEVBQ2hCO0FBQ0o7QUFFQSxTQUFTLE1BQU0sSUFBWSxRQUFxQztBQUM1RCxNQUFJLENBQUMsT0FBUSxRQUFPLElBQUksUUFBUSxDQUFDLFlBQVksV0FBVyxTQUFTLEVBQUUsQ0FBQztBQUNwRSxNQUFJLE9BQU8sUUFBUyxRQUFPLFFBQVEsT0FBTyxJQUFJLHVCQUF1QixDQUFDO0FBRXRFLFNBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLFVBQU0sVUFBVSxNQUFNO0FBQ2xCLG1CQUFhLEtBQUs7QUFDbEIsYUFBTyxJQUFJLHVCQUF1QixDQUFDO0FBQUEsSUFDdkM7QUFDQSxVQUFNLFFBQVEsV0FBVyxNQUFNO0FBQzNCLGFBQU8sb0JBQW9CLFNBQVMsT0FBTztBQUMzQyxjQUFRO0FBQUEsSUFDWixHQUFHLEVBQUU7QUFDTCxXQUFPLGlCQUFpQixTQUFTLFNBQVMsRUFBRSxNQUFNLEtBQUssQ0FBQztBQUFBLEVBQzVELENBQUM7QUFDTDtBQVlPLElBQU0sZUFBTixNQUFNLHFCQUFvQixZQUFZO0FBQUEsRUFvRHpDLFlBQVksTUFBYztBQUN0QixVQUFNO0FBL0NWLFNBQVMsYUFBYTtBQUN0QixTQUFTLE9BQU87QUFDaEIsU0FBUyxVQUFVO0FBQ25CLFNBQVMsU0FBUztBQU9sQjtBQUFBLFNBQVMsV0FBVztBQUNwQixTQUFTLGFBQWE7QUFFdEIsc0JBQXFDO0FBQ3JDLHNCQUFxQixhQUFZO0FBS2pDO0FBQUE7QUFBQTtBQUFBLGtCQUF1QztBQUN2QyxxQkFBaUQ7QUFDakQsbUJBQTZDO0FBQzdDLG1CQUF3QztBQU94QztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUEsbUJBQTZDLENBQUMsWUFBWTtBQUcxRCxTQUFRLFlBQVk7QUFHcEI7QUFBQTtBQUFBLFNBQVEsU0FBd0IsUUFBUSxRQUFRO0FBQ2hELFNBQVEsV0FBa0MsQ0FBQztBQUMzQyxTQUFRLFlBQVk7QUFDcEIsU0FBUSxhQUFhLElBQUksZ0JBQWdCO0FBQ3pDLFNBQVEsYUFBYSxJQUFJLGdCQUFnQjtBQVNyQyxTQUFLLE9BQU87QUFDWixTQUFLLE1BQU0sRUFBRTtBQUNiLFNBQUssTUFBTSxTQUFTLFVBQVUsTUFBTSxJQUFJLE1BQU0sbUJBQW1CLElBQUksSUFBSTtBQUV6RSxlQUFXLFFBQVEsQ0FBQyxRQUFRLFdBQVcsU0FBUyxPQUFPLEdBQUc7QUFDdEQsNEJBQXNCLE1BQU0sSUFBSTtBQUFBLElBQ3BDO0FBRUEsTUFBRSxZQUFZLElBQUksS0FBSyxLQUFLLElBQUk7QUFJaEMsU0FBSyxTQUFTLEtBQUssT0FBTztBQUFBLE1BQUssTUFDM0IsVUFBVSxLQUFLLEtBQUssV0FBVyxJQUFJLFdBQVcsQ0FBQyxHQUFHLE1BQU0sS0FBSyxXQUFXLE1BQU07QUFBQSxJQUNsRixFQUFFLE1BQU0sQ0FBQyxRQUFRO0FBQ2IsV0FBSyxNQUFNLEdBQUc7QUFBQSxJQUNsQixDQUFDO0FBRUQsaUJBQWE7QUFBQSxFQUNqQjtBQUFBO0FBQUEsRUF6QkEsSUFBSSxpQkFBeUI7QUFDekIsV0FBTyxLQUFLO0FBQUEsRUFDaEI7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBLEVBNkJBLEtBQUssTUFBK0Q7QUFDaEUsUUFBSSxLQUFLLGVBQWUsYUFBWSxZQUFZO0FBRTVDLFlBQU0sSUFBSSxhQUFhLDhCQUE4QixtQkFBbUI7QUFBQSxJQUM1RTtBQUNBLFFBQUksS0FBSyxlQUFlLGFBQVksTUFBTTtBQUN0QztBQUFBLElBQ0o7QUFNQSxVQUFNLFlBQVksWUFBWSxJQUFJO0FBQ2xDLFVBQU0sV0FBVyxZQUNYLFFBQVEsUUFBUSxTQUFTLElBQ3pCLFFBQVEsTUFBTSxLQUFLLFdBQVcsTUFBTTtBQUsxQyxTQUFLLFNBQVMsTUFBTSxNQUFNO0FBQUEsSUFBQyxDQUFDO0FBTTVCLFVBQU0sV0FDRixDQUFDLGFBQWEsT0FBTyxTQUFTLGVBQWUsZ0JBQWdCLE9BQU8sS0FBSyxPQUFPO0FBQ3BGLFFBQUksVUFBVTtBQUNWLFdBQUssYUFBYTtBQUFBLElBQ3RCO0FBTUEsUUFBSSxXQUFXO0FBQ1gsV0FBSyxhQUFhLFVBQVU7QUFBQSxJQUNoQztBQUNBLFNBQUssU0FBUyxLQUFLLFFBQVE7QUFRM0IsUUFBSSxLQUFLLFVBQVc7QUFDcEIsU0FBSyxZQUFZO0FBRWpCLFNBQUssU0FBUyxLQUFLLE9BQU8sS0FBSyxZQUFZO0FBQ3ZDLFVBQUk7QUFDQSxlQUFPLEtBQUssU0FBUyxTQUFTLEdBQUc7QUFHN0IsZ0JBQU0sU0FBUyxLQUFLLFNBQVMsT0FBTyxHQUFHLGdCQUFnQjtBQUN2RCxnQkFBTSxRQUFRLE1BQU0sUUFBUSxJQUFJLE1BQU07QUFLdEMsY0FBSSxLQUFLLGVBQWUsYUFBWSxRQUFRO0FBQ3hDO0FBQUEsVUFDSjtBQUVBLGdCQUFNLFFBQVEsTUFBTSxPQUFPLENBQUMsR0FBRyxNQUFNLElBQUksRUFBRSxZQUFZLENBQUM7QUFDeEQsY0FBSTtBQUNBLGtCQUFNLFVBQVUsS0FBSyxLQUFLLE9BQU8sS0FBSyxXQUFXLE1BQU07QUFBQSxVQUMzRCxVQUFFO0FBSUUsaUJBQUssWUFBWSxLQUFLLElBQUksR0FBRyxLQUFLLFlBQVksS0FBSztBQUFBLFVBQ3ZEO0FBQUEsUUFDSjtBQUFBLE1BQ0osVUFBRTtBQUNFLGFBQUssWUFBWTtBQUFBLE1BQ3JCO0FBQUEsSUFDSixDQUFDLEVBQUUsTUFBTSxDQUFDLFFBQVE7QUFDZCxXQUFLLFlBQVk7QUFDakIsV0FBSyxNQUFNLEdBQUc7QUFBQSxJQUNsQixDQUFDO0FBQUEsRUFDTDtBQUFBO0FBQUEsRUFHQSxNQUFNLE9BQWUsS0FBTSxTQUFpQixJQUFVO0FBQ2xELFFBQUksS0FBSyxlQUFlLGFBQVksV0FBVyxLQUFLLGVBQWUsYUFBWSxRQUFRO0FBQ25GO0FBQUEsSUFDSjtBQUNBLFNBQUssYUFBYSxhQUFZO0FBQzlCLFNBQUssU0FBUyxTQUFTO0FBQ3ZCLFNBQUssWUFBWTtBQUNqQixTQUFLLFdBQVcsTUFBTTtBQU10QixTQUFLLFdBQVcsTUFBTTtBQUN0QixTQUFLLFNBQVMsS0FBSyxPQUFPO0FBQUEsTUFBSyxNQUMzQixVQUFVLEtBQUssS0FBSyxZQUFZLElBQUksV0FBVyxDQUFDLENBQUM7QUFBQSxJQUNyRCxFQUFFLE1BQU0sTUFBTTtBQUFBLElBRWQsQ0FBQyxFQUFFLEtBQUssTUFBTTtBQUNWLFdBQUssUUFBUSxNQUFNLFFBQVEsSUFBSTtBQUFBLElBQ25DLENBQUM7QUFBQSxFQUNMO0FBQUE7QUFBQSxFQUdBLFVBQWdCO0FBQ1osUUFBSSxLQUFLLGVBQWUsYUFBWSxXQUFZO0FBQ2hELFNBQUssYUFBYSxhQUFZO0FBQzlCLFNBQUssY0FBYyxJQUFJLE1BQU0sTUFBTSxDQUFDO0FBQUEsRUFDeEM7QUFBQTtBQUFBLEVBR0EsU0FBUyxTQUE0QjtBQUNqQyxRQUFJLEtBQUssZUFBZSxhQUFZLEtBQU07QUFDMUMsUUFBSTtBQUNKLFFBQUk7QUFDQSxhQUFPLEtBQUssUUFBUSxPQUFPO0FBQUEsSUFDL0IsU0FBUTtBQUNKLFdBQUssY0FBYyxJQUFJLE1BQU0sT0FBTyxDQUFDO0FBQ3JDO0FBQUEsSUFDSjtBQUNBLFFBQUksU0FBUyxXQUFXLEtBQUssZUFBZSxRQUFRO0FBQ2hELGFBQU8sSUFBSSxLQUFLLENBQUMsT0FBTyxDQUFDO0FBQUEsSUFDN0I7QUFDQSxTQUFLLGNBQWMsSUFBSSxhQUFhLFdBQVcsRUFBRSxLQUFLLENBQUMsQ0FBQztBQUFBLEVBQzVEO0FBQUE7QUFBQSxFQUdBLE1BQU0sS0FBb0I7QUFJdEIsUUFBSSxLQUFLLGVBQWUsYUFBWSxVQUFVLEtBQUssZUFBZSxhQUFZLFFBQVM7QUFDdkYsU0FBSyxjQUFjLElBQUksTUFBTSxPQUFPLENBQUM7QUFHckMsU0FBSyxRQUFRLE1BQU0sZUFBZSxRQUFRLElBQUksVUFBVSxPQUFPLEdBQUcsR0FBRyxLQUFLO0FBQUEsRUFDOUU7QUFBQTtBQUFBLEVBR0EsUUFBUSxNQUFjLFFBQWdCLFVBQXlCO0FBbGFuRSxRQUFBQztBQW1hUSxRQUFJLEtBQUssZUFBZSxhQUFZLE9BQVE7QUFDNUMsU0FBSyxhQUFhLGFBQVk7QUFDOUIsU0FBSyxXQUFXLE1BQU07QUFDdEIsU0FBSyxXQUFXLE1BQU07QUFDdEIsTUFBRSxZQUFZLE9BQU8sS0FBSyxHQUFHO0FBQzdCLFFBQUksRUFBRSxZQUFZLFNBQVMsRUFBRyxFQUFBQSxNQUFBLEVBQUUsY0FBRixnQkFBQUEsSUFBYTtBQUMzQyxTQUFLLFNBQVMsU0FBUztBQUN2QixTQUFLLFlBQVk7QUFFakIsVUFBTSxLQUFLLE9BQU8sZUFBZSxhQUMzQixJQUFJLFdBQVcsU0FBUyxFQUFFLE1BQU0sUUFBUSxTQUFTLENBQUMsSUFDbEQsT0FBTyxPQUFPLElBQUksTUFBTSxPQUFPLEdBQUcsRUFBRSxNQUFNLFFBQVEsU0FBUyxDQUFDO0FBQ2xFLFNBQUssY0FBYyxFQUFFO0FBQUEsRUFDekI7QUFDSjtBQS9PYSxhQUNPLGFBQWE7QUFEcEIsYUFFTyxPQUFPO0FBRmQsYUFHTyxVQUFVO0FBSGpCLGFBSU8sU0FBUztBQUp0QixJQUFNLGNBQU47QUF3UEEsU0FBUyxPQUFPLE1BQXVDO0FBMWI5RCxNQUFBQTtBQThiSSxNQUFJLFFBQVE7QUFDUixVQUFNLFdBQVdBLE1BQUEsT0FBZSxXQUFmLGdCQUFBQSxJQUF1QjtBQUN4QyxRQUFJLE9BQU8sWUFBWSxZQUFZO0FBQy9CLGFBQU8sUUFBUSxJQUFJO0FBQUEsSUFDdkI7QUFBQSxFQUNKO0FBQ0EsU0FBTyxJQUFJLFlBQVksSUFBSTtBQUMvQjtBQUtBLFNBQVMsc0JBQXNCLFFBQXFCLE1BQW9CO0FBQ3BFLE1BQUksVUFBc0M7QUFDMUMsU0FBTyxlQUFlLFFBQVEsT0FBTyxNQUFNO0FBQUEsSUFDdkMsS0FBSyxNQUFNO0FBQUEsSUFDWCxJQUFJLElBQWdDO0FBQ2hDLFVBQUksUUFBUyxRQUFPLG9CQUFvQixNQUFNLE9BQU87QUFDckQsZ0JBQVUsT0FBTyxPQUFPLGFBQWEsS0FBSztBQUMxQyxVQUFJLFFBQVMsUUFBTyxpQkFBaUIsTUFBTSxPQUFPO0FBQUEsSUFDdEQ7QUFBQSxJQUNBLGNBQWM7QUFBQSxJQUNkLFlBQVk7QUFBQSxFQUNoQixDQUFDO0FBQ0w7QUErQ08sU0FBUyxXQUFXLE1BQTBCO0FBQ2pELFFBQU0sU0FBUyxPQUFPLElBQUk7QUFDMUIsUUFBTSxVQUFVLElBQUksWUFBWTtBQUNoQyxRQUFNLFFBQVEsQ0FBQyxZQUNYLEtBQUssTUFBTSxPQUFPLFlBQVksV0FBVyxVQUFVLFFBQVEsT0FBTyxPQUFPLENBQUM7QUFFOUUsTUFBSSxrQkFBa0IsYUFBYTtBQUcvQixXQUFPLFVBQVUsQ0FBQyxZQUFZLE1BQU0sT0FBTztBQUFBLEVBQy9DLE9BQU87QUFLSCxVQUFNLFNBQVM7QUFDZixXQUFPLGFBQWE7QUFFcEIsVUFBTSxVQUFVLG9CQUFJLFFBQStCO0FBQ25ELFVBQU0sZ0JBQWdCLG9CQUFJLFFBQW9DO0FBQzlELFVBQU0sTUFBTSxPQUFPLGlCQUFpQixLQUFLLE1BQU07QUFDL0MsVUFBTSxTQUFTLE9BQU8sb0JBQW9CLEtBQUssTUFBTTtBQUVyRCxVQUFNLGNBQWMsQ0FBQyxPQUFtQztBQTVoQmhFLFVBQUFBO0FBNmhCWSxVQUFJLGNBQWMsSUFBSSxFQUFFLEVBQUcsU0FBT0EsTUFBQSxjQUFjLElBQUksRUFBRSxNQUFwQixPQUFBQSxNQUF5QjtBQUMzRCxVQUFJO0FBQ0EsY0FBTSxRQUFRLE1BQU8sR0FBb0IsSUFBSTtBQUM3QyxjQUFNLFVBQVUsSUFBSSxhQUFhLFdBQVcsRUFBRSxNQUFNLE1BQU0sQ0FBQztBQUMzRCxzQkFBYyxJQUFJLElBQUksT0FBTztBQUM3QixlQUFPO0FBQUEsTUFDWCxTQUFRO0FBQ0osc0JBQWMsSUFBSSxJQUFJLElBQUk7QUFDMUIsZUFBTyxjQUFjLElBQUksTUFBTSxPQUFPLENBQUM7QUFDdkMsZUFBTztBQUFBLE1BQ1g7QUFBQSxJQUNKO0FBRUEsVUFBTSxTQUFTLENBQUMsYUFBaUM7QUFDN0MsWUFBTSxXQUFXLFFBQVEsSUFBSSxRQUFrQjtBQUMvQyxVQUFJLFNBQVUsUUFBTztBQUVyQixZQUFNLEtBQW9CLENBQUMsT0FBTztBQUM5QixjQUFNLFVBQVUsWUFBWSxFQUFFO0FBQzlCLFlBQUksQ0FBQyxRQUFTO0FBQ2QsWUFBSSxPQUFPLGFBQWEsV0FBWSxVQUFTLEtBQUssUUFBUSxPQUFPO0FBQUEsWUFDNUQsVUFBUyxZQUFZLE9BQU87QUFBQSxNQUNyQztBQUNBLGNBQVEsSUFBSSxVQUFvQixFQUFFO0FBQ2xDLGFBQU87QUFBQSxJQUNYO0FBRUEsV0FBTyxvQkFBb0IsQ0FBQyxNQUFjLFVBQWUsU0FBZTtBQUNwRSxVQUFJLE1BQU0sU0FBUyxhQUFhLFdBQVcsT0FBTyxRQUFRLElBQUksVUFBVSxJQUFJO0FBQUEsSUFDaEY7QUFFQSxXQUFPLHVCQUF1QixDQUFDLE1BQWMsVUFBZSxTQUFlO0FBQ3ZFLFlBQU0sS0FBSyxTQUFTLGFBQWEsV0FBVyxRQUFRLElBQUksUUFBUSxJQUFJO0FBQ3BFLGFBQU8sTUFBTSxrQkFBTSxVQUFVLElBQUk7QUFBQSxJQUNyQztBQUVBLFFBQUksVUFBK0M7QUFDbkQsV0FBTyxlQUFlLFFBQVEsYUFBYTtBQUFBLE1BQ3ZDLEtBQUssTUFBTTtBQUFBLE1BQ1gsSUFBSSxJQUFJO0FBQ0osWUFBSSxRQUFTLFFBQU8sb0JBQW9CLFdBQVcsT0FBd0I7QUFDM0Usa0JBQVUsT0FBTyxPQUFPLGFBQWEsS0FBSztBQUMxQyxZQUFJLFFBQVMsUUFBTyxpQkFBaUIsV0FBVyxPQUF3QjtBQUFBLE1BQzVFO0FBQUEsTUFDQSxjQUFjO0FBQUEsTUFDZCxZQUFZO0FBQUEsSUFDaEIsQ0FBQztBQUFBLEVBQ0w7QUFFQSxRQUFNLE9BQU8sT0FBTyxLQUFLLEtBQUssTUFBTTtBQUNwQyxFQUFDLE9BQWlDLE9BQU8sQ0FBQyxVQUFtQixLQUFLLEtBQUssVUFBVSxLQUFLLENBQUM7QUFFdkYsU0FBTztBQUNYO0FBSUEsU0FBUyxZQUFZLE1BQTRFO0FBQzdGLE1BQUksT0FBTyxTQUFTLFVBQVU7QUFDMUIsV0FBTyxZQUFZLE9BQU8sSUFBSTtBQUFBLEVBQ2xDO0FBQ0EsTUFBSSxPQUFPLFNBQVMsZUFBZSxnQkFBZ0IsTUFBTTtBQUNyRCxXQUFPO0FBQUEsRUFDWDtBQUNBLE1BQUksWUFBWSxPQUFPLElBQUksR0FBRztBQUMxQixXQUFPLElBQUksV0FBVyxLQUFLLFFBQVEsS0FBSyxZQUFZLEtBQUssVUFBVSxFQUFFLE1BQU07QUFBQSxFQUMvRTtBQUNBLFNBQU8sSUFBSSxXQUFXLElBQXVCLEVBQUUsTUFBTTtBQUN6RDtBQUVBLFNBQVMsVUFBYSxTQUFxQixRQUFrQztBQUN6RSxNQUFJLENBQUMsT0FBUSxRQUFPO0FBQ3BCLE1BQUksT0FBTyxRQUFTLFFBQU8sUUFBUSxPQUFPLElBQUksdUJBQXVCLENBQUM7QUFFdEUsU0FBTyxJQUFJLFFBQVEsQ0FBQyxTQUFTLFdBQVc7QUFDcEMsVUFBTSxVQUFVLE1BQU0sT0FBTyxJQUFJLHVCQUF1QixDQUFDO0FBQ3pELFdBQU8saUJBQWlCLFNBQVMsU0FBUyxFQUFFLE1BQU0sS0FBSyxDQUFDO0FBQ3hELFlBQVE7QUFBQSxNQUNKLENBQUMsVUFBVTtBQUNQLGVBQU8sb0JBQW9CLFNBQVMsT0FBTztBQUMzQyxnQkFBUSxLQUFLO0FBQUEsTUFDakI7QUFBQSxNQUNBLENBQUMsUUFBUTtBQUNMLGVBQU8sb0JBQW9CLFNBQVMsT0FBTztBQUMzQyxlQUFPLEdBQUc7QUFBQSxNQUNkO0FBQUEsSUFDSjtBQUFBLEVBQ0osQ0FBQztBQUNMO0FBRUEsZUFBZSxRQUNYLE1BQ0EsUUFDbUI7QUFDbkIsTUFBSSxPQUFPLFNBQVMsVUFBVTtBQUMxQixXQUFPLFlBQVksT0FBTyxJQUFJO0FBQUEsRUFDbEM7QUFDQSxNQUFJLE9BQU8sU0FBUyxlQUFlLGdCQUFnQixNQUFNO0FBQ3JELFdBQU8sSUFBSSxXQUFXLE1BQU0sVUFBVSxLQUFLLFlBQVksR0FBRyxNQUFNLENBQUM7QUFBQSxFQUNyRTtBQUNBLE1BQUksWUFBWSxPQUFPLElBQUksR0FBRztBQUMxQixXQUFPLElBQUksV0FBVyxLQUFLLFFBQVEsS0FBSyxZQUFZLEtBQUssVUFBVSxFQUFFLE1BQU07QUFBQSxFQUMvRTtBQUNBLFNBQU8sSUFBSSxXQUFXLElBQXVCLEVBQUUsTUFBTTtBQUN6RDtBQUVBLGVBQWUsVUFBVSxRQUFnQixNQUFjLE1BQWtCLE1BQWUsUUFBcUM7QUFDekgsUUFBTSxVQUFrQztBQUFBLElBQ3BDLENBQUMsV0FBVyxHQUFHLEVBQUU7QUFBQSxJQUNqQixDQUFDLGNBQWMsR0FBRyxPQUFPLEVBQUUsVUFBVTtBQUFBLElBQ3JDLENBQUMsUUFBUSxHQUFHLE9BQU8sTUFBTTtBQUFBLElBQ3pCLENBQUMsUUFBUSxHQUFHLE9BQU8sSUFBSTtBQUFBLElBQ3ZCLGdCQUFnQjtBQUFBLEVBQ3BCO0FBQ0EsTUFBSSxTQUFTLFFBQVc7QUFDcEIsWUFBUSxRQUFRLElBQUk7QUFBQSxFQUN4QjtBQUVBLE1BQUksS0FBSyxjQUFjRCxrQkFBaUI7QUFDcEMsVUFBTSxjQUFjLFNBQVMsTUFBTSxNQUFNO0FBQ3pDO0FBQUEsRUFDSjtBQUdBLFFBQU0sVUFBVSxPQUFPO0FBQ3ZCLFFBQU0sUUFBUSxLQUFLLEtBQUssS0FBSyxhQUFhQSxnQkFBZTtBQUN6RCxXQUFTLElBQUksR0FBRyxJQUFJLE9BQU8sS0FBSztBQUM1QixVQUFNLFFBQVEsS0FBSyxTQUFTLElBQUlBLG1CQUFrQixJQUFJLEtBQUtBLGdCQUFlO0FBQzFFLFVBQU0sY0FBYyxpQ0FDYixVQURhO0FBQUEsTUFFaEIsQ0FBQyxTQUFTLEdBQUc7QUFBQSxNQUNiLENBQUMsZUFBZSxHQUFHLE9BQU8sQ0FBQztBQUFBLE1BQzNCLENBQUMsZUFBZSxHQUFHLE9BQU8sS0FBSztBQUFBLElBQ25DLElBQUcsT0FBTyxNQUFNO0FBQUEsRUFDcEI7QUFDSjtBQU1BLGVBQWUsY0FBYyxTQUFpQyxNQUFrQixRQUFxQztBQUtqSCxNQUFJLE9BQU87QUFDWCxhQUFTO0FBQ0wsUUFBSSxpQ0FBUSxRQUFTLE9BQU0sSUFBSSx1QkFBdUI7QUFDdEQsVUFBTSxPQUFPLE1BQU0sTUFBTSxVQUFVLE1BQU0sR0FBRztBQUFBLE1BQ3hDLFFBQVE7QUFBQSxNQUNSO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxJQUNKLENBQUM7QUFDRCxRQUFJLEtBQUssR0FBSTtBQUNiLFFBQUksS0FBSyxXQUFXLEtBQUs7QUFDckIsWUFBTSxJQUFJLE1BQU0sTUFBTSxLQUFLLEtBQUssQ0FBQztBQUFBLElBQ3JDO0FBQ0EsVUFBTSxNQUFNLE1BQU0sTUFBTTtBQUN4QixXQUFPLEtBQUssSUFBSSxPQUFPLEdBQUcsRUFBRTtBQUFBLEVBQ2hDO0FBQ0o7QUFNQSxlQUFlLFVBQVUsUUFBZ0IsUUFBc0IsUUFBcUM7QUFDaEcsTUFBSSxRQUFRO0FBQ1osU0FBTyxRQUFRLE9BQU8sUUFBUTtBQUMxQixRQUFJLE9BQU8sS0FBSyxFQUFFLGFBQWFBLGtCQUFpQjtBQUM1QyxZQUFNLFVBQVUsUUFBUSxXQUFXLE9BQU8sS0FBSyxHQUFHLFFBQVcsTUFBTTtBQUNuRTtBQUNBO0FBQUEsSUFDSjtBQUVBLFFBQUksTUFBTTtBQUNWLFFBQUksT0FBTztBQUNYLFdBQU8sTUFBTSxPQUFPLFVBQVUsTUFBTSxRQUFRLGtCQUFrQjtBQUMxRCxZQUFNLFFBQVEsT0FBTyxHQUFHO0FBQ3hCLFVBQUksTUFBTSxhQUFhQSxvQkFBbUIsT0FBTyxJQUFJLE1BQU0sYUFBYUEsa0JBQWlCO0FBQ3JGO0FBQUEsTUFDSjtBQUNBLGNBQVEsSUFBSSxNQUFNO0FBQ2xCO0FBQUEsSUFDSjtBQUlBLFFBQUksUUFBUSxPQUFPO0FBQ2YsWUFBTSxVQUFVLFFBQVEsV0FBVyxPQUFPLEtBQUssR0FBRyxRQUFXLE1BQU07QUFDbkU7QUFDQTtBQUFBLElBQ0o7QUFFQSxVQUFNLFFBQVEsT0FBTyxNQUFNLE9BQU8sR0FBRztBQUNyQyxRQUFJLE1BQU0sV0FBVyxHQUFHO0FBQ3BCLFlBQU0sVUFBVSxRQUFRLFdBQVcsTUFBTSxDQUFDLEdBQUcsUUFBVyxNQUFNO0FBQUEsSUFDbEUsT0FBTztBQUNILFlBQU0sZUFBZSxRQUFRLE9BQU8sTUFBTTtBQUFBLElBQzlDO0FBQ0EsWUFBUTtBQUFBLEVBQ1o7QUFDSjtBQUVBLGVBQWUsZUFBZSxRQUFnQixRQUFzQixRQUFxQztBQTN1QnpHLE1BQUFDO0FBNHVCSSxNQUFJLE9BQU8sV0FBVyxLQUFLLE9BQU8sU0FBUyxrQkFBa0I7QUFDekQsVUFBTSxJQUFJLE1BQU0sMkJBQTJCO0FBQUEsRUFDL0M7QUFFQSxNQUFJLE9BQU8sV0FBVyxNQUFNO0FBQzVCLE1BQUksT0FBTztBQUNYLE1BQUksT0FBTztBQUNYLGFBQVM7QUFDTCxRQUFJLGlDQUFRLFFBQVMsT0FBTSxJQUFJLHVCQUF1QjtBQUN0RCxVQUFNLE9BQU8sTUFBTSxNQUFNLFVBQVUsTUFBTSxHQUFHO0FBQUEsTUFDeEMsUUFBUTtBQUFBLE1BQ1IsU0FBUztBQUFBLFFBQ0wsQ0FBQyxXQUFXLEdBQUcsRUFBRTtBQUFBLFFBQ2pCLENBQUMsY0FBYyxHQUFHLE9BQU8sRUFBRSxVQUFVO0FBQUEsUUFDckMsQ0FBQyxRQUFRLEdBQUcsT0FBTyxNQUFNO0FBQUEsUUFDekIsQ0FBQyxRQUFRLEdBQUcsT0FBTyxTQUFTO0FBQUEsUUFDNUIsQ0FBQyxTQUFTLEdBQUcsT0FBTyxPQUFPLFNBQVMsSUFBSTtBQUFBLFFBQ3hDLGdCQUFnQjtBQUFBLE1BQ3BCO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxJQUNKLENBQUM7QUFDRCxRQUFJLEtBQUssR0FBSTtBQUNiLFFBQUksS0FBSyxXQUFXLEtBQUs7QUFDckIsWUFBTSxJQUFJLE1BQU0sTUFBTSxLQUFLLEtBQUssQ0FBQztBQUFBLElBQ3JDO0FBR0EsVUFBTSxXQUFXLFFBQU9BLE1BQUEsS0FBSyxRQUFRLElBQUksU0FBUyxNQUExQixPQUFBQSxNQUErQixDQUFDO0FBQ3hELFVBQU0sWUFBWSxPQUFPLFNBQVM7QUFDbEMsUUFBSSxDQUFDLE9BQU8sVUFBVSxRQUFRLEtBQUssV0FBVyxLQUFLLFlBQVksV0FBVztBQUd0RSxVQUFJLGFBQWEsRUFBRyxPQUFNLElBQUksTUFBTSxzQ0FBc0M7QUFBQSxJQUM5RTtBQUNBLFlBQVE7QUFDUixXQUFPLFdBQVcsT0FBTyxNQUFNLElBQUksQ0FBQztBQUNwQyxVQUFNLE1BQU0sTUFBTSxNQUFNO0FBQ3hCLFdBQU8sS0FBSyxJQUFJLE9BQU8sR0FBRyxFQUFFO0FBQUEsRUFDaEM7QUFDSjtBQUVBLFNBQVMsV0FBVyxRQUFrQztBQUNsRCxNQUFJLE9BQU87QUFDWCxhQUFXLEtBQUssT0FBUSxTQUFRLElBQUksRUFBRTtBQUN0QyxRQUFNLE1BQU0sSUFBSSxXQUFXLElBQUk7QUFDL0IsUUFBTSxPQUFPLElBQUksU0FBUyxJQUFJLE1BQU07QUFDcEMsT0FBSyxVQUFVLEdBQUcsT0FBTyxNQUFNO0FBQy9CLE1BQUksTUFBTTtBQUNWLGFBQVcsS0FBSyxRQUFRO0FBQ3BCLFNBQUssVUFBVSxLQUFLLEVBQUUsVUFBVTtBQUNoQyxXQUFPO0FBQ1AsUUFBSSxJQUFJLEdBQUcsR0FBRztBQUNkLFdBQU8sRUFBRTtBQUFBLEVBQ2I7QUFDQSxTQUFPO0FBQ1g7QUFFQSxTQUFTLGVBQXFCO0FBQzFCLE1BQUksRUFBRSxXQUFXLENBQUMsT0FBUTtBQUMxQixJQUFFLFVBQVU7QUFDWixPQUFLLFNBQVM7QUFDbEI7QUFZQSxlQUFlLFdBQTBCO0FBQ3JDLE1BQUksVUFBVTtBQUNkLFFBQU0sUUFBUSxJQUFJLGdCQUFnQjtBQUNsQyxJQUFFLFlBQVk7QUFFZCxTQUFPLEVBQUUsWUFBWSxPQUFPLEdBQUc7QUFDM0IsUUFBSTtBQUNBLFlBQU0sT0FBTyxNQUFNLE1BQU0sVUFBVSxNQUFNLEdBQUc7QUFBQSxRQUN4QyxRQUFRO0FBQUEsUUFDUixTQUFTO0FBQUEsVUFDTCxDQUFDLFdBQVcsR0FBRyxFQUFFO0FBQUEsVUFDakIsQ0FBQyxjQUFjLEdBQUcsT0FBTyxFQUFFLFVBQVU7QUFBQSxRQUN6QztBQUFBLFFBQ0EsT0FBTztBQUFBLFFBQ1AsUUFBUSxNQUFNO0FBQUEsTUFDbEIsQ0FBQztBQUVELFVBQUksS0FBSyxXQUFXLEtBQUs7QUFHckIsaUJBQVMsTUFBTSxnQkFBZ0I7QUFDL0I7QUFBQSxNQUNKO0FBQ0EsVUFBSSxDQUFDLEtBQUssTUFBTSxLQUFLLFdBQVcsS0FBSztBQUNqQyxZQUFJLENBQUMsc0JBQXNCLEtBQUssTUFBTSxHQUFHO0FBQ3JDLGtCQUFRLHlCQUF5QixLQUFLLE1BQU07QUFDNUM7QUFBQSxRQUNKO0FBQ0EsY0FBTSxJQUFJLE1BQU0sb0NBQW9DLEtBQUssTUFBTTtBQUFBLE1BQ25FO0FBRUEsZ0JBQVU7QUFDVixVQUFJLEtBQUssV0FBVyxLQUFLO0FBQ3JCLFlBQUk7QUFDSixZQUFJO0FBQ0EsZ0JBQU0sTUFBTSxLQUFLLFlBQVk7QUFBQSxRQUNqQyxTQUFRO0FBS0osa0JBQVEsa0NBQWtDO0FBQzFDO0FBQUEsUUFDSjtBQUNBLGNBQU0sU0FBUyxhQUFhLEdBQUc7QUFDL0IsWUFBSSxXQUFXLE1BQU07QUFLakIsbUJBQVMsTUFBTSw2QkFBNkI7QUFDNUM7QUFBQSxRQUNKO0FBQ0EsZ0JBQVEsTUFBTTtBQUFBLE1BQ2xCO0FBQUEsSUFFSixTQUFTLEtBQUs7QUFDVixVQUFJLE1BQU0sT0FBTyxXQUFXLEVBQUUsWUFBWSxTQUFTLEVBQUc7QUFDdEQsZ0JBQVUsVUFBVSxLQUFLLElBQUksVUFBVSxHQUFHLFdBQVcsSUFBSTtBQUN6RCxVQUFJO0FBQ0EsY0FBTSxNQUFNLFNBQVMsTUFBTSxNQUFNO0FBQUEsTUFDckMsU0FBUTtBQUNKO0FBQUEsTUFDSjtBQUFBLElBQ0o7QUFBQSxFQUNKO0FBRUEsTUFBSSxFQUFFLGNBQWMsTUFBTyxHQUFFLFlBQVk7QUFDekMsSUFBRSxVQUFVO0FBR1osTUFBSSxFQUFFLFlBQVksT0FBTyxHQUFHO0FBQ3hCLGlCQUFhO0FBQUEsRUFDakI7QUFDSjtBQUVBLFNBQVMsc0JBQXNCLFFBQXlCO0FBQ3BELFNBQU8sV0FBVyxPQUFPLFdBQVcsT0FBTyxXQUFXLE9BQVEsVUFBVSxPQUFPLFVBQVU7QUFDN0Y7QUFZQSxTQUFTLGFBQWEsS0FBb0M7QUFDdEQsTUFBSSxJQUFJLGFBQWEsRUFBRyxRQUFPO0FBQy9CLFFBQU0sS0FBSyxJQUFJLFNBQVMsR0FBRztBQUMzQixNQUFJLEdBQUcsU0FBUyxDQUFDLE1BQU0sTUFBUSxHQUFHLFNBQVMsQ0FBQyxNQUFNLE1BQzlDLEdBQUcsU0FBUyxDQUFDLE1BQU0sTUFBUSxHQUFHLFNBQVMsQ0FBQyxNQUFNLE1BQzdDLEdBQUcsU0FBUyxDQUFDLElBQUksQ0FBQyxPQUFPLEdBQUc7QUFDN0IsV0FBTztBQUFBLEVBQ1g7QUFFQSxRQUFNLFFBQVEsR0FBRyxVQUFVLENBQUM7QUFDNUIsUUFBTSxTQUFvQixDQUFDO0FBQzNCLE1BQUksTUFBTTtBQUVWLFdBQVMsSUFBSSxHQUFHLElBQUksT0FBTyxLQUFLO0FBQzVCLFFBQUksTUFBTSxJQUFJLElBQUksV0FBWSxRQUFPO0FBQ3JDLFVBQU0sU0FBUyxHQUFHLFVBQVUsR0FBRztBQUFHLFdBQU87QUFDekMsVUFBTSxPQUFPLEdBQUcsU0FBUyxHQUFHO0FBQUcsV0FBTztBQUN0QyxVQUFNLE1BQU0sR0FBRyxVQUFVLEdBQUc7QUFBRyxXQUFPO0FBQ3RDLFFBQUksT0FBTyxjQUFjLE1BQU0sSUFBSSxhQUFhLElBQUssUUFBTztBQUM1RCxVQUFNLFVBQVUsSUFBSSxNQUFNLEtBQUssTUFBTSxHQUFHO0FBQUcsV0FBTztBQUNsRCxXQUFPLEtBQUssRUFBRSxRQUFRLE1BQU0sUUFBUSxDQUFDO0FBQUEsRUFDekM7QUFFQSxTQUFPLFFBQVEsSUFBSSxhQUFhLFNBQVM7QUFDN0M7QUFFQSxTQUFTLFFBQVEsUUFBeUI7QUFDdEMsYUFBVyxTQUFTLFFBQVE7QUFDeEIsVUFBTSxTQUFTLE1BQU07QUFDckIsVUFBTSxPQUFPLE1BQU07QUFDbkIsVUFBTSxVQUFVLE1BQU07QUFDdEIsVUFBTSxPQUFPLEVBQUUsWUFBWSxJQUFJLE1BQU07QUFDckMsUUFBSSxDQUFDLEtBQU07QUFFWCxZQUFRLE1BQU07QUFBQSxNQUNWLEtBQUs7QUFDRCxhQUFLLFNBQVMsT0FBTztBQUNyQjtBQUFBLE1BQ0osS0FBSztBQUNELGFBQUssUUFBUTtBQUNiO0FBQUEsTUFDSixLQUFLO0FBQ0QsYUFBSyxRQUFRLEtBQU0sSUFBSSxJQUFJO0FBQzNCO0FBQUEsTUFDSixLQUFLO0FBQ0QsYUFBSyxNQUFNLElBQUksTUFBTSxJQUFJLFlBQVksRUFBRSxPQUFPLE9BQU8sQ0FBQyxDQUFDO0FBQ3ZEO0FBQUEsSUFDUjtBQUFBLEVBQ0o7QUFDSjtBQUVBLFNBQVMsU0FBUyxNQUFjLFFBQXNCO0FBQ2xELGFBQVcsUUFBUSxDQUFDLEdBQUcsRUFBRSxZQUFZLE9BQU8sQ0FBQyxHQUFHO0FBQzVDLFNBQUssUUFBUSxNQUFNLFFBQVEsS0FBSztBQUFBLEVBQ3BDO0FBQ0o7QUFFQSxTQUFTLFFBQVEsUUFBc0I7QUFDbkMsYUFBVyxRQUFRLENBQUMsR0FBRyxFQUFFLFlBQVksT0FBTyxDQUFDLEdBQUc7QUFDNUMsU0FBSyxNQUFNLElBQUksTUFBTSxNQUFNLENBQUM7QUFBQSxFQUNoQztBQUNKOzs7QUQzN0JBLElBQU0sWUFBWSx1QkFBTyxJQUFJLG9CQUFvQjtBQUVqRCxJQUFNLGtCQUFrQixLQUFLLE9BQU87QUFFcEMsU0FBUyxhQUFhLFFBQTJCO0FBdEJqRCxNQUFBQztBQXVCSSxNQUFJLE9BQU8sU0FBUztBQUNoQixXQUFNQSxNQUFBLE9BQU8sV0FBUCxPQUFBQSxNQUFpQixJQUFJLGFBQWEsd0JBQXdCLFlBQVk7QUFBQSxFQUNoRjtBQUNKO0FBc0JBLGVBQXNCLFVBQVUsU0FBMkIsWUFBb0IsTUFBYyxVQUF5QixDQUFDLEdBQWtCO0FBaER6SSxNQUFBQSxLQUFBO0FBaURJLFFBQU0sWUFBV0EsTUFBQSxRQUFRLGFBQVIsT0FBQUEsTUFBb0I7QUFDckMsTUFBSSxDQUFDLE9BQU8sY0FBYyxRQUFRLEtBQUssWUFBWSxHQUFHO0FBQ2xELFVBQU0sSUFBSSxXQUFXLDBDQUEwQztBQUFBLEVBQ25FO0FBQ0EsTUFBSSxDQUFDLFdBQVcsS0FBSyxLQUFLLENBQUMsS0FBSyxLQUFLLEdBQUc7QUFDcEMsVUFBTSxJQUFJLFVBQVUsMkNBQTJDO0FBQUEsRUFDbkU7QUFDQSxNQUFJLFFBQVEsT0FBUSxjQUFhLFFBQVEsTUFBTTtBQUUvQyxRQUFNLFVBQVU7QUFDaEIsUUFBTSxTQUFRLGFBQVEsU0FBUyxNQUFqQixZQUFzQixFQUFFLFlBQVksRUFBRTtBQUNwRCxVQUFRLFNBQVMsSUFBSTtBQUNyQixjQUFNLGVBQU4sbUJBQWtCO0FBQ2xCLFFBQU0sYUFBYSxFQUFFLE1BQU07QUFDM0IsUUFBTSxhQUFhLElBQUksZ0JBQWdCO0FBQ3ZDLFFBQU0sYUFBYTtBQUNuQixRQUFNLFNBQVMsTUFBRztBQWpFdEIsUUFBQUE7QUFpRXlCLHNCQUFXLE9BQU1BLE1BQUEsUUFBUSxXQUFSLGdCQUFBQSxJQUFnQixNQUFNO0FBQUE7QUFDNUQsZ0JBQVEsV0FBUixtQkFBZ0IsaUJBQWlCLFNBQVMsUUFBUSxFQUFFLE1BQU0sS0FBSztBQUMvRCxNQUFJO0FBQ0EsVUFBTSxPQUFPLE1BQU0sWUFBWSxZQUFZLE1BQU0sVUFBVSxXQUFXLE1BQU07QUFDNUUsaUJBQWEsV0FBVyxNQUFNO0FBQzlCLFFBQUksTUFBTSxlQUFlLFlBQVk7QUFDakMsWUFBTSxJQUFJLGFBQWEseUJBQXlCLFlBQVk7QUFBQSxJQUNoRTtBQUNBLFVBQU0sWUFBWSxJQUFJLGdCQUFnQixJQUFJO0FBQzFDLFFBQUk7QUFDQSxjQUFRLE1BQU07QUFBQSxJQUNsQixTQUFTLE9BQU87QUFDWixVQUFJLGdCQUFnQixTQUFTO0FBQzdCLFlBQU07QUFBQSxJQUNWO0FBQ0EsUUFBSSxNQUFNLFVBQVcsS0FBSSxnQkFBZ0IsTUFBTSxTQUFTO0FBQ3hELFVBQU0sWUFBWTtBQUFBLEVBQ3RCLFVBQUU7QUFDRSxrQkFBUSxXQUFSLG1CQUFnQixvQkFBb0IsU0FBUztBQUM3QyxRQUFJLE1BQU0sZUFBZSxXQUFZLE9BQU0sYUFBYTtBQUFBLEVBQzVEO0FBQ0o7QUFRTyxTQUFTLFlBQVksU0FBaUM7QUE5RjdELE1BQUFBO0FBK0ZJLFFBQU0sVUFBVTtBQUNoQixRQUFNLFFBQVEsUUFBUSxTQUFTO0FBQy9CLE1BQUksQ0FBQyxNQUFPO0FBQ1osSUFBRSxNQUFNO0FBQ1IsR0FBQUEsTUFBQSxNQUFNLGVBQU4sZ0JBQUFBLElBQWtCO0FBQ2xCLFNBQU8sUUFBUSxTQUFTO0FBQ3hCLE1BQUksTUFBTSxXQUFXO0FBQ2pCLFFBQUk7QUFDQSxjQUFRLE1BQU07QUFDZCxjQUFRLGdCQUFnQixLQUFLO0FBQzdCLGNBQVEsS0FBSztBQUFBLElBQ2pCLFVBQUU7QUFDRSxVQUFJLGdCQUFnQixNQUFNLFNBQVM7QUFBQSxJQUN2QztBQUFBLEVBQ0o7QUFDSjtBQUVBLFNBQVMsWUFBWSxZQUFvQixNQUFjLFVBQWtCLFFBQW9DO0FBQ3pHLFNBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLFVBQU0sU0FBUyxPQUFPLFVBQVU7QUFDaEMsV0FBTyxhQUFhO0FBQ3BCLFVBQU0sUUFBb0IsQ0FBQztBQUMzQixRQUFJO0FBQ0osUUFBSSxjQUFjO0FBQ2xCLFFBQUksV0FBVztBQUNmLFFBQUksV0FBVztBQUVmLFVBQU0sU0FBUyxDQUFDLFVBQW9CO0FBQ2hDLFVBQUksU0FBVTtBQUNkLGlCQUFXO0FBQ1gsYUFBTyxvQkFBb0IsU0FBUyxNQUFNO0FBQzFDLGFBQU8sU0FBUyxPQUFPLFlBQVksT0FBTyxVQUFVLE9BQU8sVUFBVTtBQUNyRSxhQUFPLE1BQU07QUFDYixVQUFJLFVBQVUsUUFBVztBQUNyQixjQUFNLFNBQVM7QUFDZixlQUFPLEtBQUs7QUFBQSxNQUNoQixPQUFPO0FBQ0gsWUFBSTtBQUNBLGtCQUFRLElBQUksS0FBSyxPQUFPLEVBQUUsTUFBTSxZQUFZLENBQUMsQ0FBQztBQUFBLFFBQ2xELFNBQVMsT0FBTztBQUNaLGlCQUFPLEtBQUs7QUFBQSxRQUNoQixVQUFFO0FBQ0UsZ0JBQU0sU0FBUztBQUFBLFFBQ25CO0FBQUEsTUFDSjtBQUFBLElBQ0o7QUFDQSxVQUFNLFNBQVMsTUFBRztBQTdJMUIsVUFBQUE7QUE2STZCLHFCQUFPQSxNQUFBLE9BQU8sV0FBUCxPQUFBQSxNQUFpQixJQUFJLGFBQWEsd0JBQXdCLFlBQVksQ0FBQztBQUFBO0FBQ25HLFdBQU8saUJBQWlCLFNBQVMsUUFBUSxFQUFFLE1BQU0sS0FBSyxDQUFDO0FBQ3ZELFFBQUksT0FBTyxTQUFTO0FBQUUsYUFBTztBQUFHO0FBQUEsSUFBUTtBQUV4QyxXQUFPLFNBQVMsTUFBTTtBQUNsQixVQUFJLFNBQVU7QUFDZCxVQUFJO0FBQUUsZUFBTyxLQUFLLEtBQUssVUFBVSxFQUFFLFNBQVMsR0FBRyxNQUFNLFNBQVMsQ0FBQyxDQUFDO0FBQUEsTUFBRyxTQUM1RCxPQUFPO0FBQUUsZUFBTyxLQUFLO0FBQUEsTUFBRztBQUFBLElBQ25DO0FBQ0EsV0FBTyxVQUFVLE1BQU0sT0FBTyxJQUFJLE1BQU0scUJBQXFCLENBQUM7QUFDOUQsV0FBTyxZQUFZLFdBQVM7QUFDeEIsVUFBSSxTQUFVO0FBQ2QsVUFBSTtBQUNBLFlBQUksRUFBRSxNQUFNLGdCQUFnQixhQUFjLE9BQU0sSUFBSSxNQUFNLHFCQUFxQjtBQUMvRSxZQUFJLGFBQWEsUUFBVztBQUN4QixjQUFJLE1BQU0sS0FBSyxhQUFhLEtBQU0sT0FBTSxJQUFJLE1BQU0sd0JBQXdCO0FBQzFFLGdCQUFNLE9BQU8sS0FBSyxNQUFNLElBQUksWUFBWSxFQUFFLE9BQU8sTUFBTSxJQUFJLENBQUM7QUFDNUQsY0FBSSxLQUFLLFlBQVksRUFBRyxPQUFNLElBQUksTUFBTSw0QkFBNEI7QUFDcEUsY0FBSSxLQUFLLE1BQU8sT0FBTSxLQUFLLFFBQVEsSUFBSSxXQUFXLEtBQUssS0FBSyxJQUFJLElBQUksTUFBTSxLQUFLLEtBQUs7QUFDcEYsY0FBSSxDQUFDLE9BQU8sY0FBYyxLQUFLLElBQUksS0FBSyxLQUFLLE9BQU8sS0FBSyxPQUFPLEtBQUssU0FBUyxVQUFVO0FBQ3BGLGtCQUFNLElBQUksTUFBTSx3QkFBd0I7QUFBQSxVQUM1QztBQUNBLGNBQUksS0FBSyxPQUFPLFNBQVUsT0FBTSxJQUFJLFdBQVcscUJBQXFCLGlCQUFRLGNBQWE7QUFDekYscUJBQVcsS0FBSztBQUNoQix3QkFBYyxLQUFLO0FBQUEsUUFDdkIsT0FBTztBQUNILHNCQUFZLE1BQU0sS0FBSztBQUN2QixjQUFJLE1BQU0sS0FBSyxlQUFlLEtBQUssTUFBTSxLQUFLLGFBQWEsS0FBSyxRQUFRLFdBQVcsWUFBWSxXQUFXLFVBQVU7QUFDaEgsa0JBQU0sSUFBSSxXQUFXLDJDQUEyQztBQUFBLFVBQ3BFO0FBQ0EsZ0JBQU0sS0FBSyxNQUFNLElBQUk7QUFBQSxRQUN6QjtBQUFBLE1BQ0osU0FBUyxPQUFPO0FBQUUsZUFBTyxLQUFLO0FBQUEsTUFBRztBQUFBLElBQ3JDO0FBQ0EsV0FBTyxVQUFVLFdBQVM7QUFDdEIsVUFBSSxNQUFNLFNBQVMsT0FBUSxhQUFhLFVBQWEsYUFBYSxVQUFVO0FBQ3hFLGVBQU8sSUFBSSxNQUFNLHNEQUFzRCxDQUFDO0FBQUEsTUFDNUUsT0FBTztBQUNILGVBQU87QUFBQSxNQUNYO0FBQUEsSUFDSjtBQUFBLEVBQ0osQ0FBQztBQUNMOzs7QUV2TEE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQXdEQSxJQUFNQyxRQUFPLGlCQUFpQixZQUFZLE9BQU87QUFFakQsSUFBTSxTQUFTO0FBQ2YsSUFBTSxhQUFhO0FBQ25CLElBQU0sYUFBYTtBQUNuQixJQUFNLFVBQVU7QUFDaEIsSUFBTSxhQUFhO0FBT1osU0FBUyxTQUE0QjtBQUN4QyxTQUFPQSxNQUFLLE1BQU07QUFDdEI7QUFPTyxTQUFTLGFBQThCO0FBQzFDLFNBQU9BLE1BQUssVUFBVTtBQUMxQjtBQU9PLFNBQVMsYUFBOEI7QUFDMUMsU0FBT0EsTUFBSyxVQUFVO0FBQzFCO0FBUU8sU0FBUyxRQUFRLElBQTZCO0FBQ2pELFNBQU9BLE1BQUssU0FBUyxFQUFFLEdBQUcsQ0FBQztBQUMvQjtBQVFPLFNBQVMsV0FBVyxPQUFnQztBQUN2RCxTQUFPQSxNQUFLLFlBQVksRUFBRSxNQUFNLENBQUM7QUFDckM7OztBQzdHQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBWUEsSUFBTUMsU0FBTyxpQkFBaUIsWUFBWSxHQUFHO0FBRzdDLElBQU0sZ0JBQWdCO0FBQ3RCLElBQU0sYUFBYTtBQUVaLElBQVU7QUFBQSxDQUFWLENBQVVDLGFBQVY7QUFFSSxXQUFTLE9BQU8sUUFBcUIsVUFBeUI7QUFDakUsV0FBT0QsT0FBSyxlQUFlLEVBQUUsTUFBTSxDQUFDO0FBQUEsRUFDeEM7QUFGTyxFQUFBQyxTQUFTO0FBQUEsR0FGSDtBQU9WLElBQVU7QUFBQSxDQUFWLENBQVVDLFlBQVY7QUFPSSxXQUFTQyxRQUFzQjtBQUNsQyxXQUFPSCxPQUFLLFVBQVU7QUFBQSxFQUMxQjtBQUZPLEVBQUFFLFFBQVMsT0FBQUM7QUFBQSxHQVBIOzs7QUN6QmpCO0FBQUE7QUFBQSxnQkFBQUM7QUFBQSxFQUFBLGVBQUFDO0FBQUEsRUFBQTtBQUFBO0FBWUEsSUFBTUMsU0FBTyxpQkFBaUIsWUFBWSxPQUFPO0FBR2pELElBQU0saUJBQWlCO0FBQ3ZCLElBQU1DLGNBQWE7QUFDbkIsSUFBTSxZQUFZO0FBRVgsSUFBVUM7QUFBQSxDQUFWLENBQVVBLGFBQVY7QUFFSSxXQUFTLFFBQVEsYUFBcUIsS0FBb0I7QUFDN0QsV0FBT0YsT0FBSyxnQkFBZ0IsRUFBRSxVQUFVLFdBQVcsQ0FBQztBQUFBLEVBQ3hEO0FBRk8sRUFBQUUsU0FBUztBQUFBLEdBRkhBLHdCQUFBO0FBT1YsSUFBVUM7QUFBQSxDQUFWLENBQVVBLFlBQVY7QUFXSSxXQUFTQyxRQUFzQjtBQUNsQyxXQUFPSixPQUFLQyxXQUFVO0FBQUEsRUFDMUI7QUFGTyxFQUFBRSxRQUFTLE9BQUFDO0FBQUEsR0FYSEQsc0JBQUE7QUFnQlYsSUFBVTtBQUFBLENBQVYsQ0FBVUUsV0FBVjtBQUVJLFdBQVNDLE1BQUssU0FBZ0M7QUFDakQsV0FBT04sT0FBSyxXQUFXLEVBQUUsUUFBUSxDQUFDO0FBQUEsRUFDdEM7QUFGTyxFQUFBSyxPQUFTLE9BQUFDO0FBQUEsR0FGSDs7O0FDMUNqQjtBQUFBO0FBQUEsZ0JBQUFDO0FBQUE7QUFnQ08sSUFBTUMsVUFBUyxPQUFPLE9BQU87QUFBQTtBQUFBLEVBRWhDLGNBQWM7QUFBQTtBQUFBLEVBRWQsaUJBQWlCO0FBQUE7QUFBQSxFQUVqQixVQUFVO0FBQUE7QUFBQSxFQUVWLGlCQUFpQjtBQUFBO0FBQUEsRUFFakIsa0JBQWtCO0FBQUE7QUFBQSxFQUVsQixrQkFBa0I7QUFBQTtBQUFBLEVBRWxCLFdBQVc7QUFBQTtBQUFBLEVBRVgsWUFBWTtBQUFBO0FBQUEsRUFFWixhQUFhO0FBQUE7QUFBQSxFQUViLE9BQU87QUFBQTtBQUFBLEVBRVAsTUFBTTtBQUFBO0FBQUEsRUFHTixNQUFNLE9BQU8sT0FBTztBQUFBO0FBQUEsSUFFaEIsU0FBUztBQUFBO0FBQUEsSUFFVCxTQUFTO0FBQUE7QUFBQSxJQUVULE1BQU07QUFBQTtBQUFBLElBRU4sUUFBUTtBQUFBO0FBQUEsSUFFUixRQUFRO0FBQUEsRUFDWixDQUFDO0FBQUE7QUFBQTtBQUFBLEVBSUQsUUFBUSxPQUFPLE9BQU87QUFBQTtBQUFBLElBRWxCLE9BQU87QUFBQSxFQUNYLENBQUM7QUFDTCxDQUFDOzs7QTdCL0RELElBQUksUUFBUTtBQUNSLFNBQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQUN0QztBQW9FQSxJQUFJLFFBQVE7QUFDUixTQUFPLE9BQU8sU0FBZ0I7QUFDOUIsU0FBTyxPQUFPLFdBQVc7QUFDN0I7QUFLQSxJQUFJLFFBQVE7QUFDUixTQUFPLE9BQU8seUJBQXlCLGVBQU8sdUJBQXVCLEtBQUssY0FBTTtBQUNwRjtBQUdBLElBQUksUUFBUTtBQUNSLFNBQU8sT0FBTyxrQkFBa0I7QUFDaEMsU0FBTyxPQUFPLGtCQUFrQjtBQUNoQyxTQUFPLE9BQU8saUJBQWlCO0FBQ25DO0FBRUEsSUFBSSxRQUFRO0FBQ1IsRUFBTyxPQUFPLHFCQUFxQjtBQUN2QztBQU9PLFNBQVMsbUJBQW1CLEtBQTRCO0FBQzNELFNBQU8sTUFBTSxLQUFLLEVBQUUsUUFBUSxPQUFPLENBQUMsRUFDL0IsS0FBSyxjQUFZO0FBQ2QsUUFBSSxTQUFTLElBQUk7QUFHYixZQUFNLGVBQWUsU0FBUyxRQUFRLElBQUksY0FBYyxLQUFLLElBQUksWUFBWTtBQUM3RSxVQUFJLFlBQVksU0FBUyxZQUFZLEdBQUc7QUFDcEMsY0FBTSxTQUFTLFNBQVMsY0FBYyxRQUFRO0FBQzlDLGVBQU8sTUFBTTtBQUNiLGlCQUFTLEtBQUssWUFBWSxNQUFNO0FBQUEsTUFDcEM7QUFBQSxJQUNKO0FBQUEsRUFDSixDQUFDLEVBQ0EsTUFBTSxNQUFNO0FBQUEsRUFBQyxDQUFDO0FBQ3ZCO0FBR0EsSUFBSSxRQUFRO0FBQ1IscUJBQW1CLGtCQUFrQjtBQUN6QzsiLAogICJuYW1lcyI6IFsiaWQiLCAiaSIsICJfYSIsICJFcnJvciIsICJjYWxsIiwgIkVycm9yIiwgIl9hIiwgIkFycmF5IiwgIk1hcCIsICJBcnJheSIsICJNYXAiLCAia2V5IiwgImNhbGwiLCAiX2EiLCAiX2EiLCAicmVzaXphYmxlIiwgIl9hIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJfYSIsICJfYSIsICJjYWxsIiwgIkhpZGVNZXRob2QiLCAiU2hvd01ldGhvZCIsICJpc0RvY3VtZW50RG90QWxsIiwgIl9hIiwgInJlYXNvbiIsICJ2YWx1ZSIsICJjYWxsIiwgImNhbGwiLCAiQ0hVTktfVEhSRVNIT0xEIiwgIl9hIiwgIl9hIiwgImNhbGwiLCAiY2FsbCIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIiwgIkRldmljZSIsICJIYXB0aWNzIiwgImNhbGwiLCAiRGV2aWNlSW5mbyIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIiwgIlRvYXN0IiwgIlNob3ciLCAiRXZlbnRzIiwgIkV2ZW50cyJdCn0K
