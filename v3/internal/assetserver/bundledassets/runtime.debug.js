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
    if (event.type !== "mousedown") {
      return;
    }
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZW52aXJvbm1lbnQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3J1bnRpbWUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RpYWxvZ3MudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2V2ZW50cy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvbGlzdGVuZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NyZWF0ZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRfdHlwZXMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93aW5kb3cudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL2NvbXBpbGVkL21haW4uanMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3N5c3RlbS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY29udGV4dG1lbnUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2ZsYWdzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9kcmFnLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9hcHByZWdpb24udHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbGFibGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbmNlbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2lvcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvYW5kcm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXBkYXRlci50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFNldHVwXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG59XG5cbmltcG9ydCBcIi4vY29udGV4dG1lbnUuanNcIjtcbmltcG9ydCBcIi4vZHJhZy5qc1wiO1xuaW1wb3J0IFwiLi9hcHByZWdpb24uanNcIjtcblxuLy8gUmUtZXhwb3J0IHB1YmxpYyBBUElcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gXCIuL2FwcGxpY2F0aW9uLmpzXCI7XG5pbXBvcnQgKiBhcyBCcm93c2VyIGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCAqIGFzIENhbGwgZnJvbSBcIi4vY2FsbHMuanNcIjtcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tIFwiLi9jbGlwYm9hcmQuanNcIjtcbmltcG9ydCAqIGFzIENyZWF0ZSBmcm9tIFwiLi9jcmVhdGUuanNcIjtcbmltcG9ydCAqIGFzIERpYWxvZ3MgZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0ICogYXMgRXZlbnRzIGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0ICogYXMgRmxhZ3MgZnJvbSBcIi4vZmxhZ3MuanNcIjtcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gXCIuL3N5c3RlbS5qc1wiO1xuaW1wb3J0ICogYXMgSU9TIGZyb20gXCIuL2lvcy5qc1wiO1xuaW1wb3J0ICogYXMgQW5kcm9pZCBmcm9tIFwiLi9hbmRyb2lkLmpzXCI7XG5pbXBvcnQgKiBhcyBVcGRhdGVyIGZyb20gXCIuL3VwZGF0ZXIuanNcIjtcbmltcG9ydCBXaW5kb3csIHsgaGFuZGxlRHJhZ0VudGVyLCBoYW5kbGVEcmFnTGVhdmUsIGhhbmRsZURyYWdPdmVyIH0gZnJvbSBcIi4vd2luZG93LmpzXCI7XG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sLmpzXCI7XG5cbmV4cG9ydCB7XG4gICAgQXBwbGljYXRpb24sXG4gICAgQnJvd3NlcixcbiAgICBDYWxsLFxuICAgIENsaXBib2FyZCxcbiAgICBEaWFsb2dzLFxuICAgIEV2ZW50cyxcbiAgICBGbGFncyxcbiAgICBTY3JlZW5zLFxuICAgIFN5c3RlbSxcbiAgICBJT1MsXG4gICAgQW5kcm9pZCxcbiAgICBVcGRhdGVyLFxuICAgIFdpbmRvdyxcbiAgICBXTUxcbn07XG5cbi8qKlxuICogQW4gaW50ZXJuYWwgdXRpbGl0eSBjb25zdW1lZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXG4gKlxuICogQGlnbm9yZVxuICovXG5leHBvcnQgeyBDcmVhdGUgfTtcblxuZXhwb3J0ICogZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcblxuLy8gRXhwb3J0IHRyYW5zcG9ydCBpbnRlcmZhY2VzIGFuZCB1dGlsaXRpZXNcbmV4cG9ydCB7XG4gICAgc2V0VHJhbnNwb3J0LFxuICAgIGdldFRyYW5zcG9ydCxcbiAgICB0eXBlIFJ1bnRpbWVUcmFuc3BvcnQsXG4gICAgb2JqZWN0TmFtZXMsXG4gICAgY2xpZW50SWQsXG59IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuaW1wb3J0IHsgY2xpZW50SWQgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbi8vIE5vdGlmeSBiYWNrZW5kXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscy5pbnZva2UgPSBTeXN0ZW0uaW52b2tlO1xuICAgIHdpbmRvdy5fd2FpbHMuY2xpZW50SWQgPSBjbGllbnRJZDtcbn1cblxuLy8gUmVnaXN0ZXIgcGxhdGZvcm0gaGFuZGxlcnMgKGludGVybmFsIEFQSSlcbi8vIE5vdGU6IFdpbmRvdyBpcyB0aGUgdGhpc1dpbmRvdyBpbnN0YW5jZSAoZGVmYXVsdCBleHBvcnQgZnJvbSB3aW5kb3cudHMpXG4vLyBCaW5kaW5nIGVuc3VyZXMgJ3RoaXMnIGNvcnJlY3RseSByZWZlcnMgdG8gdGhlIGN1cnJlbnQgd2luZG93IGluc3RhbmNlXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscy5oYW5kbGVQbGF0Zm9ybUZpbGVEcm9wID0gV2luZG93LkhhbmRsZVBsYXRmb3JtRmlsZURyb3AuYmluZChXaW5kb3cpO1xufVxuXG4vLyBMaW51eC1zcGVjaWZpYyBkcmFnIGhhbmRsZXJzIChHVEsgaW50ZXJjZXB0cyBET00gZHJhZyBldmVudHMpXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscy5oYW5kbGVEcmFnRW50ZXIgPSBoYW5kbGVEcmFnRW50ZXI7XG4gICAgd2luZG93Ll93YWlscy5oYW5kbGVEcmFnTGVhdmUgPSBoYW5kbGVEcmFnTGVhdmU7XG4gICAgd2luZG93Ll93YWlscy5oYW5kbGVEcmFnT3ZlciA9IGhhbmRsZURyYWdPdmVyO1xufVxuXG5pZiAoaGFzRE9NKSB7XG4gICAgU3lzdGVtLmludm9rZShcIndhaWxzOnJ1bnRpbWU6cmVhZHlcIik7XG59XG5cbi8qKlxuICogTG9hZHMgYSBzY3JpcHQgZnJvbSB0aGUgZ2l2ZW4gVVJMIGlmIGl0IGV4aXN0cy5cbiAqIFVzZXMgSEVBRCByZXF1ZXN0IHRvIGNoZWNrIGV4aXN0ZW5jZSwgdGhlbiBpbmplY3RzIGEgc2NyaXB0IHRhZy5cbiAqIFNpbGVudGx5IGlnbm9yZXMgaWYgdGhlIHNjcmlwdCBkb2Vzbid0IGV4aXN0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gbG9hZE9wdGlvbmFsU2NyaXB0KHVybDogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGZldGNoKHVybCwgeyBtZXRob2Q6ICdIRUFEJyB9KVxuICAgICAgICAudGhlbihyZXNwb25zZSA9PiB7XG4gICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAvLyBWZXJpZnkgdGhlIHJlc3BvbnNlIGlzIGFjdHVhbGx5IEphdmFTY3JpcHQgYW5kIG5vdCBhbiBIVE1MIGZhbGxiYWNrXG4gICAgICAgICAgICAgICAgLy8gKGUuZy4gVml0ZSBkZXYgc2VydmVyIHJldHVybnMgaW5kZXguaHRtbCBmb3IgdW5rbm93biByb3V0ZXMpXG4gICAgICAgICAgICAgICAgY29uc3QgY29udGVudFR5cGUgPSAocmVzcG9uc2UuaGVhZGVycy5nZXQoJ2NvbnRlbnQtdHlwZScpIHx8ICcnKS50b0xvd2VyQ2FzZSgpO1xuICAgICAgICAgICAgICAgIGlmIChjb250ZW50VHlwZS5pbmNsdWRlcygnamF2YXNjcmlwdCcpKSB7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHNjcmlwdCA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoJ3NjcmlwdCcpO1xuICAgICAgICAgICAgICAgICAgICBzY3JpcHQuc3JjID0gdXJsO1xuICAgICAgICAgICAgICAgICAgICBkb2N1bWVudC5oZWFkLmFwcGVuZENoaWxkKHNjcmlwdCk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICB9KVxuICAgICAgICAuY2F0Y2goKCkgPT4ge30pOyAvLyBTaWxlbnRseSBpZ25vcmUgLSBzY3JpcHQgaXMgb3B0aW9uYWxcbn1cblxuLy8gTG9hZCBjdXN0b20uanMgaWYgYXZhaWxhYmxlICh1c2VkIGJ5IHNlcnZlciBtb2RlIGZvciBXZWJTb2NrZXQgZXZlbnRzLCBldGMuKVxuaWYgKGhhc0RPTSkge1xuICAgIGxvYWRPcHRpb25hbFNjcmlwdCgnL3dhaWxzL2N1c3RvbS5qcycpO1xufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBPcGVuVVJMIH0gZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xuaW1wb3J0IHsgUXVlc3Rpb24gfSBmcm9tIFwiLi9kaWFsb2dzLmpzXCI7XG5pbXBvcnQgeyBFbWl0IH0gZnJvbSBcIi4vZXZlbnRzLmpzXCI7XG5pbXBvcnQgeyBjYW5BYm9ydExpc3RlbmVycywgd2hlblJlYWR5IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcbmltcG9ydCBXaW5kb3cgZnJvbSBcIi4vd2luZG93LmpzXCI7XG5cbi8qKlxuICogU2VuZHMgYW4gZXZlbnQgd2l0aCB0aGUgZ2l2ZW4gbmFtZSBhbmQgb3B0aW9uYWwgZGF0YS5cbiAqXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gc2VuZC5cbiAqIEBwYXJhbSBbZGF0YT1udWxsXSAtIC0gT3B0aW9uYWwgZGF0YSB0byBzZW5kIGFsb25nIHdpdGggdGhlIGV2ZW50LlxuICovXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpOiB2b2lkIHtcbiAgICBFbWl0KGV2ZW50TmFtZSwgZGF0YSk7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBzcGVjaWZpZWQgd2luZG93LlxuICpcbiAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBjYWxsIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgdG8gY2FsbC5cbiAqL1xuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dOYW1lOiBzdHJpbmcsIG1ldGhvZE5hbWU6IHN0cmluZykge1xuICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IFdpbmRvdy5HZXQod2luZG93TmFtZSk7XG4gICAgY29uc3QgbWV0aG9kID0gKHRhcmdldFdpbmRvdyBhcyBhbnkpW21ldGhvZE5hbWVdO1xuXG4gICAgaWYgKHR5cGVvZiBtZXRob2QgIT09IFwiZnVuY3Rpb25cIikge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBXaW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JyBub3QgZm91bmRgKTtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIHRyeSB7XG4gICAgICAgIG1ldGhvZC5jYWxsKHRhcmdldFdpbmRvdyk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBjYWxsaW5nIHdpbmRvdyBtZXRob2QgJyR7bWV0aG9kTmFtZX0nOiBgLCBlKTtcbiAgICB9XG59XG5cbi8qKlxuICogUmVzcG9uZHMgdG8gYSB0cmlnZ2VyaW5nIGV2ZW50IGJ5IHJ1bm5pbmcgYXBwcm9wcmlhdGUgV01MIGFjdGlvbnMgZm9yIHRoZSBjdXJyZW50IHRhcmdldC5cbiAqL1xuZnVuY3Rpb24gb25XTUxUcmlnZ2VyZWQoZXY6IEV2ZW50KTogdm9pZCB7XG4gICAgY29uc3QgZWxlbWVudCA9IGV2LmN1cnJlbnRUYXJnZXQgYXMgRWxlbWVudDtcblxuICAgIGZ1bmN0aW9uIHJ1bkVmZmVjdChjaG9pY2UgPSBcIlllc1wiKSB7XG4gICAgICAgIGlmIChjaG9pY2UgIT09IFwiWWVzXCIpXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1ldmVudCcpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xuICAgICAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRhcmdldC13aW5kb3cnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdGFyZ2V0LXdpbmRvdycpIHx8IFwiXCI7XG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCB1cmwgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLW9wZW51cmwnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtb3BlbnVybCcpO1xuXG4gICAgICAgIGlmIChldmVudFR5cGUgIT09IG51bGwpXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgaWYgKHdpbmRvd01ldGhvZCAhPT0gbnVsbClcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2QodGFyZ2V0V2luZG93LCB3aW5kb3dNZXRob2QpO1xuICAgICAgICBpZiAodXJsICE9PSBudWxsKVxuICAgICAgICAgICAgdm9pZCBPcGVuVVJMKHVybCk7XG4gICAgfVxuXG4gICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG5cbiAgICBpZiAoY29uZmlybSkge1xuICAgICAgICBRdWVzdGlvbih7XG4gICAgICAgICAgICBUaXRsZTogXCJDb25maXJtXCIsXG4gICAgICAgICAgICBNZXNzYWdlOiBjb25maXJtLFxuICAgICAgICAgICAgRGV0YWNoZWQ6IGZhbHNlLFxuICAgICAgICAgICAgQnV0dG9uczogW1xuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiWWVzXCIgfSxcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIk5vXCIsIElzRGVmYXVsdDogdHJ1ZSB9XG4gICAgICAgICAgICBdXG4gICAgICAgIH0pLnRoZW4ocnVuRWZmZWN0KTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBydW5FZmZlY3QoKTtcbiAgICB9XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBjb250cm9sbGVyU3ltID0gU3ltYm9sKFwiY29udHJvbGxlclwiKTtcbmNvbnN0IHRyaWdnZXJNYXBTeW0gPSBTeW1ib2woXCJ0cmlnZ2VyTWFwXCIpO1xuY29uc3QgZWxlbWVudENvdW50U3ltID0gU3ltYm9sKFwiZWxlbWVudENvdW50XCIpO1xuXG4vKipcbiAqIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IGRvZXMgbm90IGFjdHVhbGx5IHJlbWVtYmVyIGFjdGl2ZSBldmVudCBsaXN0ZW5lcnM6IGluc3RlYWRcbiAqIGl0IHRpZXMgdGhlbSB0byBhbiBBYm9ydFNpZ25hbCBhbmQgdXNlcyBhbiBBYm9ydENvbnRyb2xsZXIgdG8gcmVtb3ZlIHRoZW0gYWxsIGF0IG9uY2UuXG4gKi9cbmNsYXNzIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICBbY29udHJvbGxlclN5bV06IEFib3J0Q29udHJvbGxlcjtcblxuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgYW4gb3B0aW9ucyBvYmplY3QgZm9yIGFkZEV2ZW50TGlzdGVuZXIgdGhhdCB0aWVzIHRoZSBsaXN0ZW5lclxuICAgICAqIHRvIHRoZSBBYm9ydFNpZ25hbCBmcm9tIHRoZSBjdXJyZW50IEFib3J0Q29udHJvbGxlci5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50c1xuICAgICAqL1xuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XG4gICAgICAgIHJldHVybiB7IHNpZ25hbDogdGhpc1tjb250cm9sbGVyU3ltXS5zaWduYWwgfTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycyBhbmQgcmVzZXRzIHRoZSByZWdpc3RyeS5cbiAgICAgKi9cbiAgICByZXNldCgpOiB2b2lkIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXS5hYm9ydCgpO1xuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBXZWFrTWFwUmVnaXN0cnkgbWFwcyBhY3RpdmUgdHJpZ2dlciBldmVudHMgdG8gZWFjaCBET00gZWxlbWVudCB0aHJvdWdoIGEgV2Vha01hcC5cbiAqIFRoaXMgZW5zdXJlcyB0aGF0IHRoZSBtYXBwaW5nIHJlbWFpbnMgcHJpdmF0ZSB0byB0aGlzIG1vZHVsZSwgd2hpbGUgc3RpbGwgYWxsb3dpbmcgZ2FyYmFnZVxuICogY29sbGVjdGlvbiBvZiB0aGUgaW52b2x2ZWQgZWxlbWVudHMuXG4gKi9cbmNsYXNzIFdlYWtNYXBSZWdpc3RyeSB7XG4gICAgLyoqIFN0b3JlcyB0aGUgY3VycmVudCBlbGVtZW50LXRvLXRyaWdnZXIgbWFwcGluZy4gKi9cbiAgICBbdHJpZ2dlck1hcFN5bV06IFdlYWtNYXA8RWxlbWVudCwgc3RyaW5nW10+O1xuICAgIC8qKiBDb3VudHMgdGhlIG51bWJlciBvZiBlbGVtZW50cyB3aXRoIGFjdGl2ZSBXTUwgdHJpZ2dlcnMuICovXG4gICAgW2VsZW1lbnRDb3VudFN5bV06IG51bWJlcjtcblxuICAgIGNvbnN0cnVjdG9yKCkge1xuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIGFjdGl2ZSB0cmlnZ2VycyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50LlxuICAgICAqXG4gICAgICogQHBhcmFtIGVsZW1lbnQgLSBBbiBIVE1MIGVsZW1lbnRcbiAgICAgKiBAcGFyYW0gdHJpZ2dlcnMgLSBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRcbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xuICAgICAgICBpZiAoIXRoaXNbdHJpZ2dlck1hcFN5bV0uaGFzKGVsZW1lbnQpKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXSsrOyB9XG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0uc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICAgICAgcmV0dXJuIHt9O1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzLlxuICAgICAqL1xuICAgIHJlc2V0KCk6IHZvaWQge1xuICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXG4gICAgICAgICAgICByZXR1cm47XG5cbiAgICAgICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnKicpKSB7XG4gICAgICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXG4gICAgICAgICAgICAgICAgYnJlYWs7XG5cbiAgICAgICAgICAgIGNvbnN0IHRyaWdnZXJzID0gdGhpc1t0cmlnZ2VyTWFwU3ltXS5nZXQoZWxlbWVudCk7XG4gICAgICAgICAgICBpZiAodHJpZ2dlcnMgIT0gbnVsbCkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0tLTsgfVxuXG4gICAgICAgICAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMgfHwgW10pXG4gICAgICAgICAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIG9uV01MVHJpZ2dlcmVkKTtcbiAgICAgICAgfVxuXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0gPSBuZXcgV2Vha01hcCgpO1xuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xuICAgIH1cbn1cblxuY29uc3QgdHJpZ2dlclJlZ2lzdHJ5ID0gY2FuQWJvcnRMaXN0ZW5lcnMoKSA/IG5ldyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSgpIDogbmV3IFdlYWtNYXBSZWdpc3RyeSgpO1xuXG4vKipcbiAqIEFkZHMgZXZlbnQgbGlzdGVuZXJzIHRvIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAqL1xuZnVuY3Rpb24gYWRkV01MTGlzdGVuZXJzKGVsZW1lbnQ6IEVsZW1lbnQpOiB2b2lkIHtcbiAgICBjb25zdCB0cmlnZ2VyUmVnRXhwID0gL1xcUysvZztcbiAgICBjb25zdCB0cmlnZ2VyQXR0ciA9IChlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIik7XG4gICAgY29uc3QgdHJpZ2dlcnM6IHN0cmluZ1tdID0gW107XG5cbiAgICBsZXQgbWF0Y2g7XG4gICAgd2hpbGUgKChtYXRjaCA9IHRyaWdnZXJSZWdFeHAuZXhlYyh0cmlnZ2VyQXR0cikpICE9PSBudWxsKVxuICAgICAgICB0cmlnZ2Vycy5wdXNoKG1hdGNoWzBdKTtcblxuICAgIGNvbnN0IG9wdGlvbnMgPSB0cmlnZ2VyUmVnaXN0cnkuc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcbiAgICBmb3IgKGNvbnN0IHRyaWdnZXIgb2YgdHJpZ2dlcnMpXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2NoZWR1bGVzIGFuIGF1dG9tYXRpYyByZWxvYWQgb2YgV01MIHRvIGJlIHBlcmZvcm1lZCBhcyBzb29uIGFzIHRoZSBkb2N1bWVudCBpcyBmdWxseSBsb2FkZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbmFibGUoKTogdm9pZCB7XG4gICAgd2hlblJlYWR5KFJlbG9hZCk7XG59XG5cbi8qKlxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCk6IHZvaWQge1xuICAgIHRyaWdnZXJSZWdpc3RyeS5yZXNldCgpO1xuICAgIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0sIFt3bWwtd2luZG93XSwgW3dtbC1vcGVudXJsXSwgW2RhdGEtd21sLWV2ZW50XSwgW2RhdGEtd21sLXdpbmRvd10sIFtkYXRhLXdtbC1vcGVudXJsXScpLmZvckVhY2goYWRkV01MTGlzdGVuZXJzKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xuXG5jb25zdCBCcm93c2VyT3BlblVSTCA9IDA7XG5cbi8qKlxuICogT3BlbiBhIGJyb3dzZXIgd2luZG93IHRvIHRoZSBnaXZlbiBVUkwuXG4gKlxuICogQHBhcmFtIHVybCAtIFRoZSBVUkwgdG8gb3BlblxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmw6IHN0cmluZyB8IFVSTCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKEJyb3dzZXJPcGVuVVJMLCB7dXJsOiB1cmwudG9TdHJpbmcoKX0pO1xufVxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2FpL25hbm9pZFxuXG4vLyBUaGUgTUlUIExpY2Vuc2UgKE1JVClcbi8vXG4vLyBDb3B5cmlnaHQgMjAxNyBBbmRyZXkgU2l0bmlrIDxhbmRyZXlAc2l0bmlrLnJ1PlxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHkgb2Zcbi8vIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWwgaW5cbi8vIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMgdG9cbi8vIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGwgY29waWVzIG9mXG4vLyB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sXG4vLyAgICAgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6XG4vL1xuLy8gICAgIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxuLy8gY29waWVzIG9yIHN1YnN0YW50aWFsIHBvcnRpb25zIG9mIHRoZSBTb2Z0d2FyZS5cbi8vXG4vLyAgICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksIEZJVE5FU1Ncbi8vIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQVVUSE9SUyBPUlxuLy8gQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLCBXSEVUSEVSXG4vLyBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sIE9VVCBPRiBPUiBJTlxuLy8gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRSBTT0ZUV0FSRS5cblxuLy8gVGhpcyBhbHBoYWJldCB1c2VzIGBBLVphLXowLTlfLWAgc3ltYm9scy5cbi8vIFRoZSBvcmRlciBvZiBjaGFyYWN0ZXJzIGlzIG9wdGltaXplZCBmb3IgYmV0dGVyIGd6aXAgYW5kIGJyb3RsaSBjb21wcmVzc2lvbi5cbi8vIFJlZmVyZW5jZXMgdG8gdGhlIHNhbWUgZmlsZSAod29ya3MgYm90aCBmb3IgZ3ppcCBhbmQgYnJvdGxpKTpcbi8vIGAndXNlYCwgYGFuZG9tYCwgYW5kIGByaWN0J2Bcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIGJyb3RsaSBkZWZhdWx0IGRpY3Rpb25hcnk6XG4vLyBgLTI2VGAsIGAxOTgzYCwgYDQwcHhgLCBgNzVweGAsIGBidXNoYCwgYGphY2tgLCBgbWluZGAsIGB2ZXJ5YCwgYW5kIGB3b2xmYFxuY29uc3QgdXJsQWxwaGFiZXQgPVxuICAgICd1c2VhbmRvbS0yNlQxOTgzNDBQWDc1cHhKQUNLVkVSWU1JTkRCVVNIV09MRl9HUVpiZmdoamtscXZ3eXpyaWN0J1xuXG5leHBvcnQgZnVuY3Rpb24gbmFub2lkKHNpemU6IG51bWJlciA9IDIxKTogc3RyaW5nIHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIC8vIEEgY29tcGFjdCBhbHRlcm5hdGl2ZSBmb3IgYGZvciAodmFyIGkgPSAwOyBpIDwgc3RlcDsgaSsrKWAuXG4gICAgbGV0IGkgPSBzaXplIHwgMFxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgICAgLy8gYHwgMGAgaXMgbW9yZSBjb21wYWN0IGFuZCBmYXN0ZXIgdGhhbiBgTWF0aC5mbG9vcigpYC5cbiAgICAgICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBUcnVlIHdoZW4gcnVubmluZyBpbnNpZGUgYSBicm93c2VyL3dlYnZpZXcgd2l0aCBhIERPTSBhdmFpbGFibGUuXG4gKiBGYWxzZSB1bmRlciBzZXJ2ZXItc2lkZSByZW5kZXJpbmcgKGUuZy4gYG5leHQgYnVpbGRgIHByZXJlbmRlcmluZyksXG4gKiB3aGVyZSBhcHBsaWNhdGlvbiBjb2RlIG1heSBpbXBvcnQgdGhlIHJ1bnRpbWUgbW9kdWxlIGV2ZW4gdGhvdWdoIG5vXG4gKiBXYWlscyBBUElzIGNhbiBhY3R1YWxseSBiZSB1c2VkICgjNDY3OSkuIE1vZHVsZXMgbXVzdCBub3QgdG91Y2hcbiAqIGB3aW5kb3dgL2Bkb2N1bWVudGAgYXQgaW1wb3J0IHRpbWUgZXhjZXB0IGJlaGluZCB0aGlzIGd1YXJkLlxuICovXG5leHBvcnQgY29uc3QgaGFzRE9NID0gdHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgZG9jdW1lbnQgIT09IFwidW5kZWZpbmVkXCI7XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gXCIuL25hbm9pZC5qc1wiO1xuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxuLy8gUmVzb2x2ZWQgbGF6aWx5OiB3aW5kb3cgZG9lcyBub3QgZXhpc3Qgd2hlbiB0aGUgbW9kdWxlIGlzIGltcG9ydGVkIGR1cmluZ1xuLy8gc2VydmVyLXNpZGUgcmVuZGVyaW5nICgjNDY3OSksIGFuZCBub3RoaW5nIGNhbiBjYWxsIHRoZSBydW50aW1lIHRoZXJlLlxuZnVuY3Rpb24gcnVudGltZVVSTCgpOiBzdHJpbmcge1xuICAgIHJldHVybiB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xufVxuXG4vLyBTdGF5IHVuZGVyIFdlYlZpZXcyJ3MgfjJNQiByZXF1ZXN0IGJvZHkgYnVmZmVyaW5nIGxpbWl0IGluIFdlYlJlc291cmNlUmVxdWVzdGVkLlxuY29uc3QgQ0hVTktfVEhSRVNIT0xEID0gNTEyICogMTAyNDtcblxuLy8gUmUtZXhwb3J0IG5hbm9pZCBmb3IgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnNcbmV4cG9ydCB7IG5hbm9pZCB9O1xuXG50eXBlIENhbGxFcnJvclR5cGUgPSB7XG4gICAgbWVzc2FnZTogc3RyaW5nLFxuICAgIGNhdXNlPzogdW5rbm93bixcbiAgICBraW5kOiBcIlJlZmVyZW5jZUVycm9yXCIgfCBcIlR5cGVFcnJvclwiIHwgXCJSdW50aW1lRXJyb3JcIlxufVxuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdGhyb3duIGluIGNhc2UgdGhlIGJvdW5kIG1ldGhvZCByZXR1cm5zIGFuIGVycm9yLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgUnVudGltZUVycm9yI25hbWV9IHByb3BlcnR5IGlzIFwiUnVudGltZUVycm9yXCIuXG4gKi9cbmV4cG9ydCBjbGFzcyBSdW50aW1lRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogQ29uc3RydWN0cyBhIG5ldyBSdW50aW1lRXJyb3IgaW5zdGFuY2UuXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIlJ1bnRpbWVFcnJvclwiO1xuICAgIH1cbn1cblxuLy8gT2JqZWN0IE5hbWVzXG5leHBvcnQgY29uc3Qgb2JqZWN0TmFtZXMgPSBPYmplY3QuZnJlZXplKHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxuICAgIENhbmNlbENhbGw6IDEwLFxuICAgIElPUzogMTEsXG4gICAgQW5kcm9pZDogMTIsXG59KTtcbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcblxuLyoqXG4gKiBSdW50aW1lVHJhbnNwb3J0IGRlZmluZXMgdGhlIGludGVyZmFjZSBmb3IgY3VzdG9tIElQQyB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zLlxuICogSW1wbGVtZW50IHRoaXMgaW50ZXJmYWNlIHRvIHVzZSBXZWJTb2NrZXRzLCBjdXN0b20gcHJvdG9jb2xzLCBvciBhbnkgb3RoZXJcbiAqIHRyYW5zcG9ydCBtZWNoYW5pc20gaW5zdGVhZCBvZiB0aGUgZGVmYXVsdCBIVFRQIGZldGNoLlxuICovXG5leHBvcnQgaW50ZXJmYWNlIFJ1bnRpbWVUcmFuc3BvcnQge1xuICAgIC8qKlxuICAgICAqIFNlbmQgYSBydW50aW1lIGNhbGwgYW5kIHJldHVybiB0aGUgcmVzcG9uc2UuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gb2JqZWN0SUQgLSBUaGUgV2FpbHMgb2JqZWN0IElEICgwPUNhbGwsIDE9Q2xpcGJvYXJkLCBldGMuKVxuICAgICAqIEBwYXJhbSBtZXRob2QgLSBUaGUgbWV0aG9kIElEIHRvIGNhbGxcbiAgICAgKiBAcGFyYW0gd2luZG93TmFtZSAtIE9wdGlvbmFsIHdpbmRvdyBuYW1lXG4gICAgICogQHBhcmFtIGFyZ3MgLSBBcmd1bWVudHMgdG8gcGFzcyAod2lsbCBiZSBKU09OIHN0cmluZ2lmaWVkIGlmIHByZXNlbnQpXG4gICAgICogQHJldHVybnMgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHJlc3BvbnNlIGRhdGFcbiAgICAgKi9cbiAgICBjYWxsKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55Pjtcbn1cblxuLyoqXG4gKiBDdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uIChjYW4gYmUgc2V0IGJ5IHVzZXIpXG4gKi9cbmxldCBjdXN0b21UcmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsID0gbnVsbDtcblxuLyoqXG4gKiBTZXQgYSBjdXN0b20gdHJhbnNwb3J0IGZvciBhbGwgV2FpbHMgcnVudGltZSBjYWxscy5cbiAqIFRoaXMgYWxsb3dzIHlvdSB0byByZXBsYWNlIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0IHdpdGhcbiAqIFdlYlNvY2tldHMsIGN1c3RvbSBwcm90b2NvbHMsIG9yIGFueSBvdGhlciBtZWNoYW5pc20uXG4gKlxuICogQHBhcmFtIHRyYW5zcG9ydCAtIFlvdXIgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvblxuICpcbiAqIEBleGFtcGxlXG4gKiBgYGB0eXBlc2NyaXB0XG4gKiBpbXBvcnQgeyBzZXRUcmFuc3BvcnQgfSBmcm9tICcvd2FpbHMvcnVudGltZS5qcyc7XG4gKlxuICogY29uc3Qgd3NUcmFuc3BvcnQgPSB7XG4gKiAgIGNhbGw6IGFzeW5jIChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSA9PiB7XG4gKiAgICAgLy8gWW91ciBXZWJTb2NrZXQgaW1wbGVtZW50YXRpb25cbiAqICAgfVxuICogfTtcbiAqXG4gKiBzZXRUcmFuc3BvcnQod3NUcmFuc3BvcnQpO1xuICogYGBgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzZXRUcmFuc3BvcnQodHJhbnNwb3J0OiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCk6IHZvaWQge1xuICAgIGN1c3RvbVRyYW5zcG9ydCA9IHRyYW5zcG9ydDtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIGN1cnJlbnQgdHJhbnNwb3J0ICh1c2VmdWwgZm9yIGV4dGVuZGluZy93cmFwcGluZylcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGdldFRyYW5zcG9ydCgpOiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCB7XG4gICAgcmV0dXJuIGN1c3RvbVRyYW5zcG9ydDtcbn1cblxuLyoqXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxuICpcbiAqIEBwYXJhbSBvYmplY3QgLSBUaGUgb2JqZWN0IHRvIGludm9rZSB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxuICogQHJldHVybiBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZDogbnVtYmVyLCBhcmdzOiBhbnkgPSBudWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3QsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcbiAgICAvLyBVc2UgY3VzdG9tIHRyYW5zcG9ydCBpZiBhdmFpbGFibGVcbiAgICBpZiAoY3VzdG9tVHJhbnNwb3J0KSB7XG4gICAgICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQuY2FsbChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9XG5cbiAgICAvLyBEZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCgpKTtcblxuICAgIGxldCBib2R5OiB7IG9iamVjdDogbnVtYmVyOyBtZXRob2Q6IG51bWJlciwgYXJncz86IGFueSB9ID0ge1xuICAgICAgb2JqZWN0OiBvYmplY3RJRCxcbiAgICAgIG1ldGhvZFxuICAgIH1cbiAgICBpZiAoYXJncyAhPT0gbnVsbCAmJiBhcmdzICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIGJvZHkuYXJncyA9IGFyZ3M7XG4gICAgfVxuXG4gICAgbGV0IGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4gPSB7XG4gICAgICAgIFtcIngtd2FpbHMtY2xpZW50LWlkXCJdOiBjbGllbnRJZCxcbiAgICAgICAgW1wiQ29udGVudC1UeXBlXCJdOiBcImFwcGxpY2F0aW9uL2pzb25cIlxuICAgIH1cbiAgICBpZiAod2luZG93TmFtZSkge1xuICAgICAgICBoZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuXG4gICAgY29uc3QgYm9keVN0ciA9IEpTT04uc3RyaW5naWZ5KGJvZHkpO1xuICAgIGxldCByZXNwb25zZTogUmVzcG9uc2U7XG4gICAgaWYgKGJvZHlTdHIubGVuZ3RoID4gQ0hVTktfVEhSRVNIT0xEKSB7XG4gICAgICAgIHJlc3BvbnNlID0gYXdhaXQgc2VuZENodW5rZWQodXJsLCBoZWFkZXJzLCBib2R5U3RyKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICByZXNwb25zZSA9IGF3YWl0IGZldGNoKHVybCwgeyBtZXRob2Q6ICdQT1NUJywgaGVhZGVycywgYm9keTogYm9keVN0ciB9KTtcbiAgICB9XG4gICAgaWYgKCFyZXNwb25zZS5vaykge1xuICAgICAgY29uc3QgY3QgPSByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKTtcbiAgICAgIGlmIChjdD8uaW5jbHVkZXMoXCJhcHBsaWNhdGlvbi9qc29uXCIpKSB7XG4gICAgICAgICAgY29uc3QganNvbjogQ2FsbEVycm9yVHlwZSA9IGF3YWl0IHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICBsZXQgZXJyO1xuICAgICAgICAgIHN3aXRjaCAoanNvbi5raW5kKSB7XG4gICAgICAgICAgICAgIGNhc2UgXCJSZWZlcmVuY2VFcnJvclwiOiBlcnIgPSBuZXcgUmVmZXJlbmNlRXJyb3IoanNvbi5tZXNzYWdlKTsgYnJlYWs7XG4gICAgICAgICAgICAgIGNhc2UgXCJUeXBlRXJyb3JcIjogICAgICBlcnIgPSBuZXcgVHlwZUVycm9yKGpzb24ubWVzc2FnZSk7IGJyZWFrO1xuICAgICAgICAgICAgICBjYXNlIFwiUnVudGltZUVycm9yXCI6ICAgZXJyID0gbmV3IFJ1bnRpbWVFcnJvcihqc29uLm1lc3NhZ2UpOyBicmVhaztcbiAgICAgICAgICAgICAgZGVmYXVsdDogICAgICAgICAgICAgICBlcnIgPSBuZXcgRXJyb3IoanNvbi5tZXNzYWdlKTtcbiAgICAgICAgICB9XG4gICAgICAgICAgZXJyLmNhdXNlID0ganNvbi5jYXVzZTtcbiAgICAgICAgICB0aHJvdyBlcnJcbiAgICAgIH1cbiAgICAgIHRocm93IG5ldyBFcnJvcihhd2FpdCByZXNwb25zZS50ZXh0KCkpO1xuICAgIH1cblxuICAgIGlmICgocmVzcG9uc2UuaGVhZGVycy5nZXQoXCJDb250ZW50LVR5cGVcIik/LmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpID8/IC0xKSAhPT0gLTEpIHtcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgIH1cbn1cblxuLy8gc2VuZENodW5rZWQgc3BsaXRzIGEgbGFyZ2Ugc2VyaWFsaXNlZCByZXF1ZXN0IGJvZHkgaW50byBDSFVOS19USFJFU0hPTEQtc2l6ZWRcbi8vIGJ5dGUgY2h1bmtzIGFuZCBzZW5kcyB0aGVtIHNlcmlhbGx5LiAgRW5jb2RpbmcgdG8gVVRGLTggYnl0ZXMgYmVmb3JlIHNsaWNpbmdcbi8vIHByZXZlbnRzIGNvcnJ1cHRpb24gb2Ygbm9uLUJNUCBjaGFyYWN0ZXJzIChzdXJyb2dhdGUgcGFpcnMpIHRoYXQgd291bGQgb2NjdXJcbi8vIHdoZW4gc3BsaXR0aW5nIGF0IEphdmFTY3JpcHQgc3RyaW5nIGluZGljZXMuICBUaGUgR28gdHJhbnNwb3J0IGFzc2VtYmxlcyB0aGVcbi8vIHJhdyBieXRlcyBiZWZvcmUgcHJvY2Vzc2luZy4gIE9ubHkgdGhlIGZpbmFsIGNodW5rJ3MgcmVzcG9uc2UgY2FycmllcyB0aGUgUlBDIHJlc3VsdC5cbmFzeW5jIGZ1bmN0aW9uIHNlbmRDaHVua2VkKHVybDogVVJMLCBoZWFkZXJzOiBSZWNvcmQ8c3RyaW5nLCBzdHJpbmc+LCBib2R5U3RyOiBzdHJpbmcpOiBQcm9taXNlPFJlc3BvbnNlPiB7XG4gICAgY29uc3QgY2h1bmtJZCA9IG5hbm9pZCgpO1xuICAgIGNvbnN0IGJvZHlCeXRlcyA9IG5ldyBUZXh0RW5jb2RlcigpLmVuY29kZShib2R5U3RyKTtcbiAgICBjb25zdCB0b3RhbENodW5rcyA9IE1hdGguY2VpbChib2R5Qnl0ZXMubGVuZ3RoIC8gQ0hVTktfVEhSRVNIT0xEKTtcblxuICAgIGZvciAobGV0IGkgPSAwOyBpIDwgdG90YWxDaHVua3MgLSAxOyBpKyspIHtcbiAgICAgICAgY29uc3QgY2h1bmsgPSBib2R5Qnl0ZXMuc3ViYXJyYXkoaSAqIENIVU5LX1RIUkVTSE9MRCwgKGkgKyAxKSAqIENIVU5LX1RIUkVTSE9MRCk7XG4gICAgICAgIGNvbnN0IHJlc3AgPSBhd2FpdCBmZXRjaCh1cmwsIHtcbiAgICAgICAgICAgIG1ldGhvZDogJ1BPU1QnLFxuICAgICAgICAgICAgaGVhZGVyczoge1xuICAgICAgICAgICAgICAgIC4uLmhlYWRlcnMsXG4gICAgICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstaWQnOiBjaHVua0lkLFxuICAgICAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLWluZGV4JzogU3RyaW5nKGkpLFxuICAgICAgICAgICAgICAgICd4LXdhaWxzLWNodW5rLXRvdGFsJzogU3RyaW5nKHRvdGFsQ2h1bmtzKSxcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICBib2R5OiBjaHVuayxcbiAgICAgICAgfSk7XG4gICAgICAgIGlmICghcmVzcC5vaykge1xuICAgICAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3AudGV4dCgpKTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIHJldHVybiBmZXRjaCh1cmwsIHtcbiAgICAgICAgbWV0aG9kOiAnUE9TVCcsXG4gICAgICAgIGhlYWRlcnM6IHtcbiAgICAgICAgICAgIC4uLmhlYWRlcnMsXG4gICAgICAgICAgICAneC13YWlscy1jaHVuay1pZCc6IGNodW5rSWQsXG4gICAgICAgICAgICAneC13YWlscy1jaHVuay1pbmRleCc6IFN0cmluZyh0b3RhbENodW5rcyAtIDEpLFxuICAgICAgICAgICAgJ3gtd2FpbHMtY2h1bmstdG90YWwnOiBTdHJpbmcodG90YWxDaHVua3MpLFxuICAgICAgICB9LFxuICAgICAgICBib2R5OiBib2R5Qnl0ZXMuc3ViYXJyYXkoKHRvdGFsQ2h1bmtzIC0gMSkgKiBDSFVOS19USFJFU0hPTEQpLFxuICAgIH0pO1xufVxuXG4vKipcbiAqIEFuZHJvaWQgV2ViVmlldyBjYW5ub3QgZGVsaXZlciBmZXRjaCgpIFBPU1QgYm9kaWVzIHRvXG4gKiBzaG91bGRJbnRlcmNlcHRSZXF1ZXN0LCBzbyB0aGUgZGVmYXVsdCBIVFRQIHRyYW5zcG9ydCBjYW5ub3QgcmVhY2ggR28uXG4gKiBXaGVuIHRoZSBBbmRyb2lkIEphdmFzY3JpcHRJbnRlcmZhY2UgYnJpZGdlICh3aW5kb3cud2FpbHMpIGlzIHByZXNlbnQsXG4gKiByb3V0ZSBydW50aW1lIGNhbGxzIHRocm91Z2ggaXQgaW5zdGVhZC4gUmVzcG9uc2VzIGFycml2ZSB2aWFcbiAqIHdpbmRvdy5fd2FpbHNBbmRyb2lkQ2FsbGJhY2ssIGludm9rZWQgYnkgdGhlIEphdmEgc2lkZS5cbiAqL1xuaW50ZXJmYWNlIEFuZHJvaWRKU0JyaWRnZSB7XG4gICAgaW52b2tlQXN5bmMoY2FsbGJhY2tJRDogc3RyaW5nLCBwYXlsb2FkOiBzdHJpbmcpOiB2b2lkO1xufVxuXG5jb25zdCBhbmRyb2lkQnJpZGdlOiBBbmRyb2lkSlNCcmlkZ2UgfCBudWxsID0gaGFzRE9NICYmXG4gICAgdHlwZW9mICh3aW5kb3cgYXMgYW55KS53YWlscz8uaW52b2tlQXN5bmMgPT09IFwiZnVuY3Rpb25cIiA/ICh3aW5kb3cgYXMgYW55KS53YWlscyA6IG51bGw7XG5cbmlmIChhbmRyb2lkQnJpZGdlKSB7XG4gICAgY29uc3QgcGVuZGluZyA9IG5ldyBNYXA8c3RyaW5nLCB7IHJlc29sdmU6ICh2YWx1ZTogYW55KSA9PiB2b2lkOyByZWplY3Q6IChyZWFzb246IGFueSkgPT4gdm9pZCB9PigpO1xuXG4gICAgKHdpbmRvdyBhcyBhbnkpLl93YWlsc0FuZHJvaWRDYWxsYmFjayA9IChpZDogc3RyaW5nLCByZXNwb25zZTogc3RyaW5nIHwgbnVsbCwgZXJyb3I6IHN0cmluZyB8IG51bGwpID0+IHtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IHBlbmRpbmcuZ2V0KGlkKTtcbiAgICAgICAgaWYgKCFwcm9taXNlKSByZXR1cm47XG4gICAgICAgIHBlbmRpbmcuZGVsZXRlKGlkKTtcbiAgICAgICAgaWYgKGVycm9yKSB7XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChuZXcgRXJyb3IoZXJyb3IpKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICB0cnkge1xuICAgICAgICAgICAgY29uc3QgZW52ZWxvcGUgPSBKU09OLnBhcnNlKHJlc3BvbnNlID8/IFwie31cIik7XG4gICAgICAgICAgICBpZiAoIWVudmVsb3BlLm9rKSB7XG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QobmV3IEVycm9yKGVudmVsb3BlLmVycm9yID8/IFwidW5rbm93biBydW50aW1lIGNhbGwgZXJyb3JcIikpO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHByb21pc2UucmVzb2x2ZShcInRleHRcIiBpbiBlbnZlbG9wZSA/IGVudmVsb3BlLnRleHQgOiBlbnZlbG9wZS5kYXRhKTtcbiAgICAgICAgfSBjYXRjaCAoZSkge1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QoZSk7XG4gICAgICAgIH1cbiAgICB9O1xuXG4gICAgY3VzdG9tVHJhbnNwb3J0ID0ge1xuICAgICAgICBjYWxsKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55PiB7XG4gICAgICAgICAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIGNvbnN0IGlkID0gbmFub2lkKCk7XG4gICAgICAgICAgICAgICAgcGVuZGluZy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIGFuZHJvaWRCcmlkZ2UuaW52b2tlQXN5bmMoaWQsIEpTT04uc3RyaW5naWZ5KHtcbiAgICAgICAgICAgICAgICAgICAgICAgIG9iamVjdDogb2JqZWN0SUQsXG4gICAgICAgICAgICAgICAgICAgICAgICBtZXRob2Q6IG1ldGhvZCxcbiAgICAgICAgICAgICAgICAgICAgICAgIHdpbmRvd05hbWU6IHdpbmRvd05hbWUsXG4gICAgICAgICAgICAgICAgICAgICAgICBhcmdzOiBhcmdzID8/IG51bGwsXG4gICAgICAgICAgICAgICAgICAgICAgICBjbGllbnRJZDogY2xpZW50SWQsXG4gICAgICAgICAgICAgICAgICAgIH0pKTtcbiAgICAgICAgICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICAgICAgICAgIC8vIERvbid0IGxlYWsgdGhlIHBlbmRpbmcgZW50cnkgaWYgZGlzcGF0Y2ggdGhyb3dzIHN5bmNocm9ub3VzbHlcbiAgICAgICAgICAgICAgICAgICAgcGVuZGluZy5kZWxldGUoaWQpO1xuICAgICAgICAgICAgICAgICAgICByZWplY3QoZSk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSk7XG4gICAgICAgIH0sXG4gICAgfTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG4vLyBzZXR1cFxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xufVxuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5EaWFsb2cpO1xuXG4vLyBEZWZpbmUgY29uc3RhbnRzIGZyb20gdGhlIGBtZXRob2RzYCBvYmplY3QgaW4gVGl0bGUgQ2FzZVxuY29uc3QgRGlhbG9nSW5mbyA9IDA7XG5jb25zdCBEaWFsb2dXYXJuaW5nID0gMTtcbmNvbnN0IERpYWxvZ0Vycm9yID0gMjtcbmNvbnN0IERpYWxvZ1F1ZXN0aW9uID0gMztcbmNvbnN0IERpYWxvZ09wZW5GaWxlID0gNDtcbmNvbnN0IERpYWxvZ1NhdmVGaWxlID0gNTtcblxuZXhwb3J0IGludGVyZmFjZSBPcGVuRmlsZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBtdWx0aXBsZSBzZWxlY3Rpb24gaXMgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBTYXZlRmlsZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBEZWZhdWx0IGZpbGVuYW1lIHRvIHVzZSBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIEZpbGVuYW1lPzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZGlyZWN0b3JpZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xuICAgIENhbkNob29zZUZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGhpZGRlbiBmaWxlcyBzaG91bGQgYmUgc2hvd24uICovXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xuICAgIFJlc29sdmVzQWxpYXNlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXG4gICAgSGlkZUV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGUgcGFja2FnZXMgc2hvdWxkIGJlIHRyZWF0ZWQgYXMgZGlyZWN0b3JpZXMuICovXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xuICAgIEFsbG93c090aGVyRmlsZXR5cGVzPzogYm9vbGVhbjtcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XG4gICAgLyoqIFRpdGxlIG9mIHRoZSBkaWFsb2cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XG4gICAgLyoqIERpcmVjdG9yeSB0byBvcGVuIGluIHRoZSBkaWFsb2cuICovXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBNZXNzYWdlRGlhbG9nT3B0aW9ucyB7XG4gICAgLyoqIFRoZSB0aXRsZSBvZiB0aGUgZGlhbG9nIHdpbmRvdy4gKi9cbiAgICBUaXRsZT86IHN0cmluZztcbiAgICAvKiogVGhlIG1haW4gbWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXG4gICAgTWVzc2FnZT86IHN0cmluZztcbiAgICAvKiogQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIEJ1dHRvbnM/OiBCdXR0b25bXTtcbiAgICAvKiogVHJ1ZSBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cgKGlmIGFwcGxpY2FibGUpLiAqL1xuICAgIERldGFjaGVkPzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBCdXR0b24ge1xuICAgIC8qKiBUZXh0IHRoYXQgYXBwZWFycyB3aXRoaW4gdGhlIGJ1dHRvbi4gKi9cbiAgICBMYWJlbD86IHN0cmluZztcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC4gKi9cbiAgICBJc0NhbmNlbD86IGJvb2xlYW47XG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgYmUgdGhlIGRlZmF1bHQgYWN0aW9uIHdoZW4gdGhlIHVzZXIgcHJlc3NlcyBlbnRlci4gKi9cbiAgICBJc0RlZmF1bHQ/OiBib29sZWFuO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEZpbGVGaWx0ZXIge1xuICAgIC8qKiBEaXNwbGF5IG5hbWUgZm9yIHRoZSBmaWx0ZXIsIGl0IGNvdWxkIGJlIFwiVGV4dCBGaWxlc1wiLCBcIkltYWdlc1wiIGV0Yy4gKi9cbiAgICBEaXNwbGF5TmFtZT86IHN0cmluZztcbiAgICAvKiogUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy4gKi9cbiAgICBQYXR0ZXJuPzogc3RyaW5nO1xufVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZGlhbG9nIG9mIHNwZWNpZmllZCB0eXBlIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKlxuICogQHBhcmFtIHR5cGUgLSBEaWFsb2cgdHlwZS5cbiAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyBmb3IgdGhlIGRpYWxvZy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggcmVzdWx0IG9mIGRpYWxvZy5cbiAqL1xuZnVuY3Rpb24gZGlhbG9nKHR5cGU6IG51bWJlciwgb3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMgfCBPcGVuRmlsZURpYWxvZ09wdGlvbnMgfCBTYXZlRmlsZURpYWxvZ09wdGlvbnMgPSB7fSk6IFByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIGNhbGwodHlwZSwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogUHJlc2VudHMgYW4gaW5mbyBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgd2FybmluZyBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhbiBlcnJvciBkaWFsb2cuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRXJyb3Iob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgcXVlc3Rpb24gZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIG9uZSBvciBtb3JlIGZpbGVzIHRvIG9wZW4uXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcywgb3IgYSBibGFuayBzdHJpbmcvZW1wdHkgbGlzdCBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbjogdHJ1ZSB9KTogUHJvbWlzZTxzdHJpbmdbXT47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGZhbHNlIHwgdW5kZWZpbmVkIH0pOiBQcm9taXNlPHN0cmluZz47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT47XG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmcgfCBzdHJpbmdbXT4geyByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKSA/PyBbXTsgfVxuXG4vKipcbiAqIFByZXNlbnRzIGEgZmlsZSBzZWxlY3Rpb24gZGlhbG9nIHRvIHBpY2sgYSBmaWxlIHRvIHNhdmUuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9ucy5cbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUsIG9yIGEgYmxhbmsgc3RyaW5nIGlmIG5vIGZpbGUgaGFzIGJlZW4gc2VsZWN0ZWQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTYXZlRmlsZShvcHRpb25zOiBTYXZlRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTsgfVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB7IGV2ZW50TGlzdGVuZXJzLCBMaXN0ZW5lciwgbGlzdGVuZXJPZmYgfSBmcm9tIFwiLi9saXN0ZW5lci5qc1wiO1xuaW1wb3J0IHsgRXZlbnRzIGFzIENyZWF0ZSB9IGZyb20gXCIuL2NyZWF0ZS5qc1wiO1xuaW1wb3J0IHsgVHlwZXMgfSBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xuXG4vLyBTZXR1cFxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xuICAgIHdpbmRvdy5fd2FpbHMuZGlzcGF0Y2hXYWlsc0V2ZW50ID0gZGlzcGF0Y2hXYWlsc0V2ZW50O1xufVxuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5FdmVudHMpO1xuY29uc3QgRW1pdE1ldGhvZCA9IDA7XG5cbmV4cG9ydCAqIGZyb20gXCIuL2V2ZW50X3R5cGVzLmpzXCI7XG5cbi8qKlxuICogQSB0YWJsZSBvZiBkYXRhIHR5cGVzIGZvciBhbGwga25vd24gZXZlbnRzLlxuICogV2lsbCBiZSBtb25rZXktcGF0Y2hlZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgQ3VzdG9tRXZlbnRzIHt9XG5cbi8qKlxuICogRWl0aGVyIGEga25vd24gZXZlbnQgbmFtZSBvciBhbiBhcmJpdHJhcnkgc3RyaW5nLlxuICovXG5leHBvcnQgdHlwZSBXYWlsc0V2ZW50TmFtZTxFIGV4dGVuZHMga2V5b2YgQ3VzdG9tRXZlbnRzID0ga2V5b2YgQ3VzdG9tRXZlbnRzPiA9IEUgfCAoc3RyaW5nICYge30pO1xuXG4vKipcbiAqIFVuaW9uIG9mIGFsbCBrbm93biBzeXN0ZW0gZXZlbnQgbmFtZXMuXG4gKi9cbnR5cGUgU3lzdGVtRXZlbnROYW1lID0ge1xuICAgIFtLIGluIGtleW9mICh0eXBlb2YgVHlwZXMpXTogKHR5cGVvZiBUeXBlcylbS11ba2V5b2YgKCh0eXBlb2YgVHlwZXMpW0tdKV1cbn0gZXh0ZW5kcyAoaW5mZXIgTSkgPyBNW2tleW9mIE1dIDogbmV2ZXI7XG5cbi8qKlxuICogVGhlIGRhdGEgdHlwZSBhc3NvY2lhdGVkIHRvIGEgZ2l2ZW4gZXZlbnQuXG4gKi9cbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnREYXRhPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiA9XG4gICAgRSBleHRlbmRzIGtleW9mIEN1c3RvbUV2ZW50cyA/IEN1c3RvbUV2ZW50c1tFXSA6IChFIGV4dGVuZHMgU3lzdGVtRXZlbnROYW1lID8gdm9pZCA6IGFueSk7XG5cbi8qKlxuICogVGhlIHR5cGUgb2YgaGFuZGxlcnMgZm9yIGEgZ2l2ZW4gZXZlbnQuXG4gKi9cbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnRDYWxsYmFjazxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4gPSAoZXY6IFdhaWxzRXZlbnQ8RT4pID0+IHZvaWQ7XG5cbi8qKlxuICogUmVwcmVzZW50cyBhIHN5c3RlbSBldmVudCBvciBhIGN1c3RvbSBldmVudCBlbWl0dGVkIHRocm91Z2ggd2FpbHMtcHJvdmlkZWQgZmFjaWxpdGllcy5cbiAqL1xuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+IHtcbiAgICAvKipcbiAgICAgKiBUaGUgbmFtZSBvZiB0aGUgZXZlbnQuXG4gICAgICovXG4gICAgbmFtZTogRTtcblxuICAgIC8qKlxuICAgICAqIE9wdGlvbmFsIGRhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBlbWl0dGVkIGV2ZW50LlxuICAgICAqL1xuICAgIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+O1xuXG4gICAgLyoqXG4gICAgICogTmFtZSBvZiB0aGUgb3JpZ2luYXRpbmcgd2luZG93LiBPbWl0dGVkIGZvciBhcHBsaWNhdGlvbiBldmVudHMuXG4gICAgICogV2lsbCBiZSBvdmVycmlkZGVuIGlmIHNldCBtYW51YWxseS5cbiAgICAgKi9cbiAgICBzZW5kZXI/OiBzdHJpbmc7XG5cbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBFLCBkYXRhOiBXYWlsc0V2ZW50RGF0YTxFPik7XG4gICAgY29uc3RydWN0b3IobmFtZTogV2FpbHNFdmVudERhdGE8RT4gZXh0ZW5kcyBudWxsIHwgdm9pZCA/IEUgOiBuZXZlcilcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBFLCBkYXRhPzogYW55KSB7XG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGEgPz8gbnVsbDtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudDogYW55KSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAoIWxpc3RlbmVycykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgbGV0IHdhaWxzRXZlbnQgPSBuZXcgV2FpbHNFdmVudChcbiAgICAgICAgZXZlbnQubmFtZSxcbiAgICAgICAgKGV2ZW50Lm5hbWUgaW4gQ3JlYXRlKSA/IENyZWF0ZVtldmVudC5uYW1lXShldmVudC5kYXRhKSA6IGV2ZW50LmRhdGFcbiAgICApO1xuICAgIGlmICgnc2VuZGVyJyBpbiBldmVudCkge1xuICAgICAgICB3YWlsc0V2ZW50LnNlbmRlciA9IGV2ZW50LnNlbmRlcjtcbiAgICB9XG5cbiAgICAvLyBEaXNwYXRjaCB0byBhIHNuYXBzaG90LCB0aGVuIHJlbW92ZSBhbGwgZXhwaXJlZCBsaXN0ZW5lcnMgaW4gYSBzaW5nbGVcbiAgICAvLyBwb3N0LWRpc3BhdGNoIGZpbHRlciBvZiB0aGUgbGl2ZSBtYXAuXG4gICAgLy8gLSBXcml0aW5nIHRoZSBzbmFwc2hvdCBiYWNrIHdob2xlc2FsZSB3b3VsZCB1bmRvIHN1YnNjcmlwdGlvbiBjaGFuZ2VzXG4gICAgLy8gICBtYWRlIGluc2lkZSBhIGhhbmRsZXIgKCM0MzkzKS5cbiAgICAvLyAtIENhbGxpbmcgbGlzdGVuZXJPZmYoKSBwZXIgZXhwaXJlZCBsaXN0ZW5lciBpbnNpZGUgdGhlIGxvb3AgaXMgTyhuXHUwMEIyKVxuICAgIC8vICAgd2hlbiBtYW55IGxpc3RlbmVycyBleHBpcmUgb24gdGhlIHNhbWUgZXZlbnQuXG4gICAgLy8gRmlsdGVyaW5nIHRoZSBsaXZlIGFycmF5IG9uY2UgYWZ0ZXIgZGlzcGF0Y2ggaXMgTyhuKSBhbmQgc3RpbGwgaG9ub3Vyc1xuICAgIC8vIGFueSBsaXN0ZW5lcnMgYWRkZWQgb3IgcmVtb3ZlZCBieSBoYW5kbGVycyBkdXJpbmcgZGlzcGF0Y2guXG4gICAgY29uc3QgZXhwaXJlZCA9IG5ldyBTZXQ8TGlzdGVuZXI+KCk7XG4gICAgZm9yIChjb25zdCBsaXN0ZW5lciBvZiBsaXN0ZW5lcnMuc2xpY2UoKSkge1xuICAgICAgICBpZiAobGlzdGVuZXIuZGlzcGF0Y2god2FpbHNFdmVudCkpIHtcbiAgICAgICAgICAgIGV4cGlyZWQuYWRkKGxpc3RlbmVyKTtcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoZXhwaXJlZC5zaXplID4gMCkge1xuICAgICAgICBjb25zdCBsaXZlID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgICAgICBpZiAobGl2ZSkge1xuICAgICAgICAgICAgY29uc3QgcmVtYWluaW5nID0gbGl2ZS5maWx0ZXIobCA9PiAhZXhwaXJlZC5oYXMobCkpO1xuICAgICAgICAgICAgaWYgKHJlbWFpbmluZy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCByZW1haW5pbmcpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxuICpcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxuICogQHBhcmFtIG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGU8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPiwgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIHNwZWNpZmllZCBldmVudCBvY2N1cnMuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT248RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPik6ICgpID0+IHZvaWQge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+KTogKCkgPT4gdm9pZCB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZXMgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnRzIHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKC4uLmV2ZW50TmFtZXM6IFtXYWlsc0V2ZW50TmFtZSwgLi4uV2FpbHNFdmVudE5hbWVbXV0pOiB2b2lkIHtcbiAgICBldmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcbn1cblxuLyoqXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKTogdm9pZCB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLyoqXG4gKiBFbWl0cyBhbiBldmVudC5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCB3aWxsIGJlIGZ1bGZpbGxlZCBvbmNlIHRoZSBldmVudCBoYXMgYmVlbiBlbWl0dGVkLiAgUmVzb2x2ZXMgdG8gdHJ1ZSBpZiB0aGUgZXZlbnQgd2FzIGNhbmNlbGxlZC5cbiAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIGVtaXRcbiAqIEBwYXJhbSBkYXRhIC0gVGhlIGRhdGEgdGhhdCB3aWxsIGJlIHNlbnQgd2l0aCB0aGUgZXZlbnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IEUsIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+KTogUHJvbWlzZTxib29sZWFuPlxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IFdhaWxzRXZlbnREYXRhPEU+IGV4dGVuZHMgbnVsbCB8IHZvaWQgPyBFIDogbmV2ZXIpOiBQcm9taXNlPGJvb2xlYW4+XG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogV2FpbHNFdmVudERhdGE8RT4sIGRhdGE/OiBhbnkpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCAgbmV3IFdhaWxzRXZlbnQobmFtZSwgZGF0YSkpXG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gVGhlIGZvbGxvd2luZyB1dGlsaXRpZXMgaGF2ZSBiZWVuIGZhY3RvcmVkIG91dCBvZiAuL2V2ZW50cy50c1xuLy8gZm9yIHRlc3RpbmcgcHVycG9zZXMuXG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXA8c3RyaW5nLCBMaXN0ZW5lcltdPigpO1xuXG5leHBvcnQgY2xhc3MgTGlzdGVuZXIge1xuICAgIGV2ZW50TmFtZTogc3RyaW5nO1xuICAgIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkO1xuICAgIG1heENhbGxiYWNrczogbnVtYmVyO1xuXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkLCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgdGhpcy5jYWxsYmFjayA9IGNhbGxiYWNrO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICB9XG5cbiAgICBkaXNwYXRjaChkYXRhOiBhbnkpOiBib29sZWFuIHtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHRoaXMuY2FsbGJhY2soZGF0YSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlcnIpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkgcmV0dXJuIGZhbHNlO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXI6IExpc3RlbmVyKTogdm9pZCB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChsaXN0ZW5lci5ldmVudE5hbWUpO1xuICAgIGlmICghbGlzdGVuZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShsaXN0ZW5lci5ldmVudE5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChsaXN0ZW5lci5ldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIEFueSBpcyBhIGR1bW15IGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBzaW1wbGUgb3IgdW5rbm93biB0eXBlcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEFueTxUID0gYW55Pihzb3VyY2U6IGFueSk6IFQge1xuICAgIHJldHVybiBzb3VyY2U7XG59XG5cbi8qKlxuICogQnl0ZVNsaWNlIGlzIGEgY3JlYXRpb24gZnVuY3Rpb24gdGhhdCByZXBsYWNlc1xuICogbnVsbCBzdHJpbmdzIHdpdGggZW1wdHkgc3RyaW5ncy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5dGVTbGljZShzb3VyY2U6IGFueSk6IHN0cmluZyB7XG4gICAgcmV0dXJuICgoc291cmNlID09IG51bGwpID8gXCJcIiA6IHNvdXJjZSk7XG59XG5cbi8qKlxuICogQXJyYXkgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcnJheVxuICogd2hvc2UgZWxlbWVudHMgYXJlIG9mIHRoYXQgdHlwZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEFycmF5PFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gVFtdIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybiBbXTtcbiAgICAgICAgfVxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNvdXJjZS5sZW5ndGg7IGkrKykge1xuICAgICAgICAgICAgc291cmNlW2ldID0gZWxlbWVudChzb3VyY2VbaV0pO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gb2JqZWN0XG4gKiB3aG9zZSBrZXlzIGFuZCB2YWx1ZXMgYXJlIG9mIHRob3NlIHR5cGVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTWFwPEsgZXh0ZW5kcyBQcm9wZXJ0eUtleSA9IGFueSwgViA9IGFueT4oa2V5OiAoc291cmNlOiBhbnkpID0+IEssIHZhbHVlOiAoc291cmNlOiBhbnkpID0+IFYpOiAoc291cmNlOiBhbnkpID0+IFJlY29yZDxLLCBWPiB7XG4gICAgaWYgKHZhbHVlID09PSBBbnkpIHtcbiAgICAgICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyB7fSA6IHNvdXJjZSk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xuICAgICAgICAgICAgcmV0dXJuIHt9O1xuICAgICAgICB9XG4gICAgICAgIGZvciAoY29uc3Qga2V5IGluIHNvdXJjZSkge1xuICAgICAgICAgICAgc291cmNlW2tleV0gPSB2YWx1ZShzb3VyY2Vba2V5XSk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcbiAgICB9O1xufVxuXG4vKipcbiAqIE51bGxhYmxlIHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXG4gKiBhbmQgcmV0dXJucyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhIG51bGxhYmxlIHZhbHVlIG9mIHRoYXQgdHlwZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE51bGxhYmxlPFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gKFQgfCBudWxsKSB7XG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gQW55O1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gbnVsbCA6IGVsZW1lbnQoc291cmNlKSk7XG59XG5cbi8qKlxuICogU3RydWN0IHRha2VzIGFuIG9iamVjdCBtYXBwaW5nIGZpZWxkIG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9uc1xuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgc3RydWN0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gU3RydWN0KGNyZWF0ZUZpZWxkOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4pOlxuICAgIDxVIGV4dGVuZHMgUmVjb3JkPHN0cmluZywgYW55PiA9IGFueT4oc291cmNlOiBhbnkpID0+IFVcbntcbiAgICBsZXQgYWxsQW55ID0gdHJ1ZTtcbiAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcbiAgICAgICAgaWYgKGNyZWF0ZUZpZWxkW25hbWVdICE9PSBBbnkpIHtcbiAgICAgICAgICAgIGFsbEFueSA9IGZhbHNlO1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cbiAgICB9XG4gICAgaWYgKGFsbEFueSkge1xuICAgICAgICByZXR1cm4gQW55O1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICAgICAgaWYgKG5hbWUgaW4gc291cmNlKSB7XG4gICAgICAgICAgICAgICAgc291cmNlW25hbWVdID0gY3JlYXRlRmllbGRbbmFtZV0oc291cmNlW25hbWVdKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogRGF0ZUZyb21UaW1lIGlzIGEgY3JlYXRpb24gZnVuY3Rpb24gdGhhdCBjb252ZXJ0cyBSRkMzMzM5IHN0cmluZ3NcbiAqIChmcm9tIEdvJ3MgdGltZS5UaW1lIEpTT04gbWFyc2hhbGluZykgdG8gSmF2YVNjcmlwdCBEYXRlIG9iamVjdHMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBEYXRlRnJvbVRpbWUoc291cmNlOiBhbnkpOiBEYXRlIHtcbiAgICByZXR1cm4gbmV3IERhdGUoc291cmNlKTtcbn1cblxuLyoqXG4gKiBNYXBzIGtub3duIGV2ZW50IG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdGhlaXIgZGF0YSB0eXBlcy5cbiAqIFdpbGwgYmUgbW9ua2V5LXBhdGNoZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxuICovXG5leHBvcnQgY29uc3QgRXZlbnRzOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4gPSB7fTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gQ3luaHlyY2h3eWQgeSBmZmVpbCBob24geW4gYXd0b21hdGlnLiBQRUlESVdDSCBcdTAwQzIgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRXYWtlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZFdha2VcIixcblx0XHRBcHBsaWNhdGlvblNjcmVlbnNEaWRTbGVlcDogXCJtYWM6QXBwbGljYXRpb25TY3JlZW5zRGlkU2xlZXBcIixcblx0XHRBcHBsaWNhdGlvblNjcmVlbnNEaWRXYWtlOiBcIm1hYzpBcHBsaWNhdGlvblNjcmVlbnNEaWRXYWtlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFNsZWVwOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxTbGVlcFwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdFdlYlZpZXdXZWJDb250ZW50UHJvY2Vzc0RpZFRlcm1pbmF0ZTogXCJtYWM6V2ViVmlld1dlYkNvbnRlbnRQcm9jZXNzRGlkVGVybWluYXRlXCIsXG5cdH0pLFxuXHRMaW51eDogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25TdGFydHVwOiBcImxpbnV4OkFwcGxpY2F0aW9uU3RhcnR1cFwiLFxuXHRcdFN5c3RlbURpZFdha2U6IFwibGludXg6U3lzdGVtRGlkV2FrZVwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJsaW51eDpTeXN0ZW1UaGVtZUNoYW5nZWRcIixcblx0XHRTeXN0ZW1XaWxsU2xlZXA6IFwibGludXg6U3lzdGVtV2lsbFNsZWVwXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibGludXg6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcblx0XHRXaW5kb3dMb2FkU3RhcnRlZDogXCJsaW51eDpXaW5kb3dMb2FkU3RhcnRlZFwiLFxuXHRcdFdpbmRvd0xvYWRSZWRpcmVjdGVkOiBcImxpbnV4OldpbmRvd0xvYWRSZWRpcmVjdGVkXCIsXG5cdFx0V2luZG93TG9hZENvbW1pdHRlZDogXCJsaW51eDpXaW5kb3dMb2FkQ29tbWl0dGVkXCIsXG5cdFx0V2luZG93TG9hZEZpbmlzaGVkOiBcImxpbnV4OldpbmRvd0xvYWRGaW5pc2hlZFwiLFxuXHR9KSxcblx0aU9TOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJpb3M6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZEVudGVyQmFja2dyb3VuZDogXCJpb3M6QXBwbGljYXRpb25EaWRFbnRlckJhY2tncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZzogXCJpb3M6QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmdcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nOiBcImlvczpBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsRW50ZXJGb3JlZ3JvdW5kOiBcImlvczpBcHBsaWNhdGlvbldpbGxFbnRlckZvcmVncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJpb3M6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0V2luZG93RGlkTG9hZDogXCJpb3M6V2luZG93RGlkTG9hZFwiLFxuXHRcdFdpbmRvd1dpbGxBcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxBcHBlYXJcIixcblx0XHRXaW5kb3dEaWRBcHBlYXI6IFwiaW9zOldpbmRvd0RpZEFwcGVhclwiLFxuXHRcdFdpbmRvd1dpbGxEaXNhcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxEaXNhcHBlYXJcIixcblx0XHRXaW5kb3dEaWREaXNhcHBlYXI6IFwiaW9zOldpbmRvd0RpZERpc2FwcGVhclwiLFxuXHRcdFdpbmRvd1NhZmVBcmVhSW5zZXRzQ2hhbmdlZDogXCJpb3M6V2luZG93U2FmZUFyZWFJbnNldHNDaGFuZ2VkXCIsXG5cdFx0V2luZG93T3JpZW50YXRpb25DaGFuZ2VkOiBcImlvczpXaW5kb3dPcmllbnRhdGlvbkNoYW5nZWRcIixcblx0XHRXaW5kb3dUb3VjaEJlZ2FuOiBcImlvczpXaW5kb3dUb3VjaEJlZ2FuXCIsXG5cdFx0V2luZG93VG91Y2hNb3ZlZDogXCJpb3M6V2luZG93VG91Y2hNb3ZlZFwiLFxuXHRcdFdpbmRvd1RvdWNoRW5kZWQ6IFwiaW9zOldpbmRvd1RvdWNoRW5kZWRcIixcblx0XHRXaW5kb3dUb3VjaENhbmNlbGxlZDogXCJpb3M6V2luZG93VG91Y2hDYW5jZWxsZWRcIixcblx0XHRXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZhaWxOYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkRmFpbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGVjaWRlUG9saWN5Rm9yTmF2aWdhdGlvbkFjdGlvbjogXCJpb3M6V2ViVmlld0RlY2lkZVBvbGljeUZvck5hdmlnYXRpb25BY3Rpb25cIixcblx0XHRCYXR0ZXJ5Q2hhbmdlZDogXCJpb3M6QmF0dGVyeUNoYW5nZWRcIixcblx0XHROZXR3b3JrQ2hhbmdlZDogXCJpb3M6TmV0d29ya0NoYW5nZWRcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiaW9zOlRoZW1lQ2hhbmdlZFwiLFxuXHRcdFNjcmVlbkxvY2tlZDogXCJpb3M6U2NyZWVuTG9ja2VkXCIsXG5cdFx0U2NyZWVuVW5sb2NrZWQ6IFwiaW9zOlNjcmVlblVubG9ja2VkXCIsXG5cdH0pLFxuXHRBbmRyb2lkOiBPYmplY3QuZnJlZXplKHtcblx0XHRBY3Rpdml0eUNyZWF0ZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eUNyZWF0ZWRcIixcblx0XHRBY3Rpdml0eVN0YXJ0ZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eVN0YXJ0ZWRcIixcblx0XHRBY3Rpdml0eVJlc3VtZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eVJlc3VtZWRcIixcblx0XHRBY3Rpdml0eVBhdXNlZDogXCJhbmRyb2lkOkFjdGl2aXR5UGF1c2VkXCIsXG5cdFx0QWN0aXZpdHlTdG9wcGVkOiBcImFuZHJvaWQ6QWN0aXZpdHlTdG9wcGVkXCIsXG5cdFx0QWN0aXZpdHlEZXN0cm95ZWQ6IFwiYW5kcm9pZDpBY3Rpdml0eURlc3Ryb3llZFwiLFxuXHRcdEFwcGxpY2F0aW9uTG93TWVtb3J5OiBcImFuZHJvaWQ6QXBwbGljYXRpb25Mb3dNZW1vcnlcIixcblx0XHRXZWJWaWV3UGFnZVN0YXJ0ZWQ6IFwiYW5kcm9pZDpXZWJWaWV3UGFnZVN0YXJ0ZWRcIixcblx0XHRXZWJWaWV3UGFnZUZpbmlzaGVkOiBcImFuZHJvaWQ6V2ViVmlld1BhZ2VGaW5pc2hlZFwiLFxuXHRcdEJhdHRlcnlDaGFuZ2VkOiBcImFuZHJvaWQ6QmF0dGVyeUNoYW5nZWRcIixcblx0XHROZXR3b3JrQ2hhbmdlZDogXCJhbmRyb2lkOk5ldHdvcmtDaGFuZ2VkXCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImFuZHJvaWQ6VGhlbWVDaGFuZ2VkXCIsXG5cdFx0U2NyZWVuTG9ja2VkOiBcImFuZHJvaWQ6U2NyZWVuTG9ja2VkXCIsXG5cdFx0U2NyZWVuVW5sb2NrZWQ6IFwiYW5kcm9pZDpTY3JlZW5VbmxvY2tlZFwiLFxuXHR9KSxcblx0Q29tbW9uOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbk9wZW5lZFdpdGhGaWxlOiBcImNvbW1vbjpBcHBsaWNhdGlvbk9wZW5lZFdpdGhGaWxlXCIsXG5cdFx0QXBwbGljYXRpb25TdGFydGVkOiBcImNvbW1vbjpBcHBsaWNhdGlvblN0YXJ0ZWRcIixcblx0XHRBcHBsaWNhdGlvbkxhdW5jaGVkV2l0aFVybDogXCJjb21tb246QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmxcIixcblx0XHRUaGVtZUNoYW5nZWQ6IFwiY29tbW9uOlRoZW1lQ2hhbmdlZFwiLFxuXHRcdFN5c3RlbURpZFdha2U6IFwiY29tbW9uOlN5c3RlbURpZFdha2VcIixcblx0XHRTeXN0ZW1XaWxsU2xlZXA6IFwiY29tbW9uOlN5c3RlbVdpbGxTbGVlcFwiLFxuXHRcdFdpbmRvd0Nsb3Npbmc6IFwiY29tbW9uOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImNvbW1vbjpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkUmVzaXplOiBcImNvbW1vbjpXaW5kb3dEaWRSZXNpemVcIixcblx0XHRXaW5kb3dEUElDaGFuZ2VkOiBcImNvbW1vbjpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RmlsZXNEcm9wcGVkOiBcImNvbW1vbjpXaW5kb3dGaWxlc0Ryb3BwZWRcIixcblx0XHRXaW5kb3dGb2N1czogXCJjb21tb246V2luZG93Rm9jdXNcIixcblx0XHRXaW5kb3dGdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJjb21tb246V2luZG93SGlkZVwiLFxuXHRcdFdpbmRvd0xvc3RGb2N1czogXCJjb21tb246V2luZG93TG9zdEZvY3VzXCIsXG5cdFx0V2luZG93TWF4aW1pc2U6IFwiY29tbW9uOldpbmRvd01heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwiY29tbW9uOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VG9nZ2xlRnJhbWVsZXNzOiBcImNvbW1vbjpXaW5kb3dUb2dnbGVGcmFtZWxlc3NcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcImNvbW1vbjpXaW5kb3dSZXN0b3JlXCIsXG5cdFx0V2luZG93UnVudGltZVJlYWR5OiBcImNvbW1vbjpXaW5kb3dSdW50aW1lUmVhZHlcIixcblx0XHRXaW5kb3dTaG93OiBcImNvbW1vbjpXaW5kb3dTaG93XCIsXG5cdFx0V2luZG93VW5GdWxsc2NyZWVuOiBcImNvbW1vbjpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcImNvbW1vbjpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJjb21tb246V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1pvb206IFwiY29tbW9uOldpbmRvd1pvb21cIixcblx0XHRXaW5kb3dab29tSW46IFwiY29tbW9uOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwiY29tbW9uOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwiY29tbW9uOldpbmRvd1pvb21SZXNldFwiLFxuXHRcdEJhdHRlcnlDaGFuZ2VkOiBcImNvbW1vbjpCYXR0ZXJ5Q2hhbmdlZFwiLFxuXHRcdE5ldHdvcmtDaGFuZ2VkOiBcImNvbW1vbjpOZXR3b3JrQ2hhbmdlZFwiLFxuXHRcdFNjcmVlbkxvY2tlZDogXCJjb21tb246U2NyZWVuTG9ja2VkXCIsXG5cdFx0U2NyZWVuVW5sb2NrZWQ6IFwiY29tbW9uOlNjcmVlblVubG9ja2VkXCIsXG5cdFx0TG93TWVtb3J5OiBcImNvbW1vbjpMb3dNZW1vcnlcIixcblx0fSksXG59KTtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxuICpcbiAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIG1lc3NhZ2UgdG8gYmUgbG9nZ2VkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZTogYW55KSB7XG4gICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgY29uc29sZS5sb2coXG4gICAgICAgICclYyB3YWlsczMgJWMgJyArIG1lc3NhZ2UgKyAnICcsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjYWEwMDAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogM3B4IDBweCAwcHggM3B4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXG4gICAgKTtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgd2VidmlldyBzdXBwb3J0cyB0aGUge0BsaW5rIE1vdXNlRXZlbnQjYnV0dG9uc30gcHJvcGVydHkuXG4gKiBMb29raW5nIGF0IHlvdSBtYWNPUyBIaWdoIFNpZXJyYSFcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGNhblRyYWNrQnV0dG9ucygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKG5ldyBNb3VzZUV2ZW50KCdtb3VzZWRvd24nKSkuYnV0dG9ucyA9PT0gMDtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgYnJvd3NlciBzdXBwb3J0cyByZW1vdmluZyBsaXN0ZW5lcnMgYnkgdHJpZ2dlcmluZyBhbiBBYm9ydFNpZ25hbFxuICogKHNlZSBodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZy9lbi1VUy9kb2NzL1dlYi9BUEkvRXZlbnRUYXJnZXQvYWRkRXZlbnRMaXN0ZW5lciNzaWduYWwpLlxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XG4gICAgaWYgKCFFdmVudFRhcmdldCB8fCAhQWJvcnRTaWduYWwgfHwgIUFib3J0Q29udHJvbGxlcilcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuXG4gICAgbGV0IHJlc3VsdCA9IHRydWU7XG5cbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcbiAgICBjb25zdCBjb250cm9sbGVyID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcbiAgICB0YXJnZXQuZGlzcGF0Y2hFdmVudChuZXcgQ3VzdG9tRXZlbnQoJ3Rlc3QnKSk7XG5cbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIFJlc29sdmVzIHRoZSBjbG9zZXN0IEhUTUxFbGVtZW50IGFuY2VzdG9yIG9mIGFuIGV2ZW50J3MgdGFyZ2V0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gZXZlbnRUYXJnZXQoZXZlbnQ6IEV2ZW50KTogSFRNTEVsZW1lbnQge1xuICAgIGlmIChldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkge1xuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0O1xuICAgIH0gZWxzZSBpZiAoIShldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkgJiYgZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgTm9kZSkge1xuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0LnBhcmVudEVsZW1lbnQgPz8gZG9jdW1lbnQuYm9keTtcbiAgICB9IGVsc2Uge1xuICAgICAgICByZXR1cm4gZG9jdW1lbnQuYm9keTtcbiAgICB9XG59XG5cbi8qKipcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcblxuIEJTRCAyLUNsYXVzZSBMaWNlbnNlXG5cbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcbiBBbGwgcmlnaHRzIHJlc2VydmVkLlxuXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XG4gbW9kaWZpY2F0aW9uLCBhcmUgcGVybWl0dGVkIHByb3ZpZGVkIHRoYXQgdGhlIGZvbGxvd2luZyBjb25kaXRpb25zIGFyZSBtZXQ6XG5cbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcbiBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lci5cblxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcbiB0aGlzIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyIGluIHRoZSBkb2N1bWVudGF0aW9uXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXG5cbiBUSElTIFNPRlRXQVJFIElTIFBST1ZJREVEIEJZIFRIRSBDT1BZUklHSFQgSE9MREVSUyBBTkQgQ09OVFJJQlVUT1JTIFwiQVMgSVNcIlxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcbiBESVNDTEFJTUVELiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQ09QWVJJR0hUIEhPTERFUiBPUiBDT05UUklCVVRPUlMgQkUgTElBQkxFXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxuIFNFUlZJQ0VTOyBMT1NTIE9GIFVTRSwgREFUQSwgT1IgUFJPRklUUzsgT1IgQlVTSU5FU1MgSU5URVJSVVBUSU9OKSBIT1dFVkVSXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxuIE9GIFRISVMgU09GVFdBUkUsIEVWRU4gSUYgQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YgU1VDSCBEQU1BR0UuXG5cbiAqKiovXG5cbmxldCBpc1JlYWR5ID0gZmFsc2U7XG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsICgpID0+IHsgaXNSZWFkeSA9IHRydWUgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiB3aGVuUmVhZHkoY2FsbGJhY2s6ICgpID0+IHZvaWQpIHtcbiAgICBpZiAoaXNSZWFkeSB8fCBkb2N1bWVudC5yZWFkeVN0YXRlID09PSAnY29tcGxldGUnKSB7XG4gICAgICAgIGNhbGxiYWNrKCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIGNhbGxiYWNrKTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmltcG9ydCB0eXBlIHsgU2NyZWVuIH0gZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xuXG4vLyBEcm9wIHRhcmdldCBjb25zdGFudHNcbmNvbnN0IERST1BfVEFSR0VUX0FUVFJJQlVURSA9ICdkYXRhLWZpbGUtZHJvcC10YXJnZXQnO1xuY29uc3QgRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTID0gJ2ZpbGUtZHJvcC10YXJnZXQtYWN0aXZlJztcbmxldCBjdXJyZW50RHJvcFRhcmdldDogRWxlbWVudCB8IG51bGwgPSBudWxsO1xuXG5jb25zdCBQb3NpdGlvbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAwO1xuY29uc3QgQ2VudGVyTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMTtcbmNvbnN0IENsb3NlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDI7XG5jb25zdCBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgPSAzO1xuY29uc3QgRW5hYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kICAgICAgID0gNDtcbmNvbnN0IEZvY3VzTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDU7XG5jb25zdCBGb3JjZVJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgPSA2O1xuY29uc3QgRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgID0gNztcbmNvbnN0IEdldFNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDg7XG5jb25zdCBHZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA5O1xuY29uc3QgSGVpZ2h0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMTA7XG5jb25zdCBIaWRlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxMTtcbmNvbnN0IElzRm9jdXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDEyO1xuY29uc3QgSXNGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gMTM7XG5jb25zdCBJc01heGltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNDtcbmNvbnN0IElzTWluaW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE1O1xuY29uc3QgTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMTY7XG5jb25zdCBNaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNztcbmNvbnN0IE5hbWVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDE4O1xuY29uc3QgT3BlbkRldlRvb2xzTWV0aG9kICAgICAgICAgICAgICAgID0gMTk7XG5jb25zdCBSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgPSAyMDtcbmNvbnN0IFJlbG9hZE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDIxO1xuY29uc3QgUmVzaXphYmxlTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMjI7XG5jb25zdCBSZXN0b3JlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAyMztcbmNvbnN0IFNldFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICA9IDI0O1xuY29uc3QgU2V0QWx3YXlzT25Ub3BNZXRob2QgICAgICAgICAgICAgID0gMjU7XG5jb25zdCBTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kICAgICAgICAgPSAyNjtcbmNvbnN0IFNldEZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICAgICA9IDI3O1xuY29uc3QgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWRNZXRob2QgID0gMjg7XG5jb25zdCBTZXRNYXhTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAyOTtcbmNvbnN0IFNldE1pblNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDMwO1xuY29uc3QgU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCAgICAgICAgID0gMzE7XG5jb25zdCBTZXRSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgPSAzMjtcbmNvbnN0IFNldFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDMzO1xuY29uc3QgU2V0VGl0bGVNZXRob2QgICAgICAgICAgICAgICAgICAgID0gMzQ7XG5jb25zdCBTZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzNTtcbmNvbnN0IFNob3dNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM2O1xuY29uc3QgU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzc7XG5jb25zdCBUb2dnbGVGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgPSAzODtcbmNvbnN0IFRvZ2dsZU1heGltaXNlTWV0aG9kICAgICAgICAgICAgICA9IDM5O1xuY29uc3QgVG9nZ2xlRnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgID0gNDA7IFxuY29uc3QgVW5GdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gNDE7XG5jb25zdCBVbk1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MjtcbmNvbnN0IFVuTWluaW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQzO1xuY29uc3QgV2lkdGhNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XG5jb25zdCBab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSA0NTtcbmNvbnN0IFpvb21Jbk1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDQ2O1xuY29uc3QgWm9vbU91dE1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gNDc7XG5jb25zdCBab29tUmVzZXRNZXRob2QgICAgICAgICAgICAgICAgICAgPSA0ODtcbmNvbnN0IFNuYXBBc3Npc3RNZXRob2QgICAgICAgICAgICAgICAgICA9IDQ5O1xuY29uc3QgRmlsZXNEcm9wcGVkICAgICAgICAgICAgICAgICAgICAgID0gNTA7XG5jb25zdCBQcmludE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA1MTtcbmNvbnN0IFNldFNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDUyO1xuXG4vKipcbiAqIEZpbmRzIHRoZSBuZWFyZXN0IGRyb3AgdGFyZ2V0IGVsZW1lbnQgYnkgd2Fsa2luZyB1cCB0aGUgRE9NIHRyZWUuXG4gKi9cbmZ1bmN0aW9uIGdldERyb3BUYXJnZXRFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQgfCBudWxsKTogRWxlbWVudCB8IG51bGwge1xuICAgIGlmICghZWxlbWVudCkge1xuICAgICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG4gICAgcmV0dXJuIGVsZW1lbnQuY2xvc2VzdChgWyR7RFJPUF9UQVJHRVRfQVRUUklCVVRFfV1gKTtcbn1cblxuLyoqXG4gKiBDaGVjayBpZiB3ZSBjYW4gdXNlIFdlYlZpZXcyJ3MgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgKFdpbmRvd3MpXG4gKiBBbHNvIGNoZWNrcyB0aGF0IEVuYWJsZUZpbGVEcm9wIGlzIHRydWUgZm9yIHRoaXMgd2luZG93LlxuICovXG5mdW5jdGlvbiBjYW5SZXNvbHZlRmlsZVBhdGhzKCk6IGJvb2xlYW4ge1xuICAgIC8vIE11c3QgaGF2ZSBXZWJWaWV3MidzIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIEFQSSAoV2luZG93cyBvbmx5KVxuICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyA9PSBudWxsKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG4gICAgLy8gTXVzdCBoYXZlIEVuYWJsZUZpbGVEcm9wIHNldCB0byB0cnVlIGZvciB0aGlzIHdpbmRvd1xuICAgIC8vIFRoaXMgZmxhZyBpcyBzZXQgYnkgdGhlIEdvIGJhY2tlbmQgZHVyaW5nIHJ1bnRpbWUgaW5pdGlhbGl6YXRpb25cbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSB0cnVlO1xufVxuXG4vKipcbiAqIFNlbmQgZmlsZSBkcm9wIHRvIGJhY2tlbmQgdmlhIFdlYlZpZXcyIChXaW5kb3dzIG9ubHkpXG4gKi9cbmZ1bmN0aW9uIHJlc29sdmVGaWxlUGF0aHMoeDogbnVtYmVyLCB5OiBudW1iZXIsIGZpbGVzOiBGaWxlW10pOiB2b2lkIHtcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMpIHtcbiAgICAgICAgKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKGBmaWxlOmRyb3A6JHt4fToke3l9YCwgZmlsZXMpO1xuICAgIH1cbn1cblxuLy8gTmF0aXZlIGRyYWcgc3RhdGUgKExpbnV4L21hY09TIGludGVyY2VwdCBET00gZHJhZyBldmVudHMpXG5sZXQgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xuXG4vKipcbiAqIENsZWFucyB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhbmQgaG92ZXIgZWZmZWN0cy5cbiAqIENhbGxlZCBvbiBkcm9wIG9yIHdoZW4gZHJhZyBsZWF2ZXMgdGhlIHdpbmRvdy5cbiAqL1xuZnVuY3Rpb24gY2xlYW51cE5hdGl2ZURyYWcoKTogdm9pZCB7XG4gICAgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCkge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbGVkIGZyb20gR28gd2hlbiBhIGZpbGUgZHJhZyBlbnRlcnMgdGhlIHdpbmRvdyBvbiBMaW51eC9tYWNPUy5cbiAqL1xuZnVuY3Rpb24gaGFuZGxlRHJhZ0VudGVyKCk6IHZvaWQge1xuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBhY3RpdmF0ZSBkcmFnIHN0YXRlXG4gICAgfVxuICAgIG5hdGl2ZURyYWdBY3RpdmUgPSB0cnVlO1xufVxuXG4vKipcbiAqIENhbGxlZCBmcm9tIEdvIHdoZW4gYSBmaWxlIGRyYWcgbGVhdmVzIHRoZSB3aW5kb3cgb24gTGludXgvbWFjT1MuXG4gKi9cbmZ1bmN0aW9uIGhhbmRsZURyYWdMZWF2ZSgpOiB2b2lkIHtcbiAgICBjbGVhbnVwTmF0aXZlRHJhZygpO1xufVxuXG4vKipcbiAqIENhbGxlZCBmcm9tIEdvIGR1cmluZyBmaWxlIGRyYWcgdG8gdXBkYXRlIGhvdmVyIHN0YXRlIG9uIExpbnV4L21hY09TLlxuICogQHBhcmFtIHggLSBYIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xuICogQHBhcmFtIHkgLSBZIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xuICovXG5mdW5jdGlvbiBoYW5kbGVEcmFnT3Zlcih4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgIGlmICghbmF0aXZlRHJhZ0FjdGl2ZSkgcmV0dXJuO1xuICAgIFxuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcbiAgICB9XG4gICAgXG4gICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xuICAgIFxuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgfVxuICAgIFxuICAgIGlmIChkcm9wVGFyZ2V0KSB7XG4gICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgIH1cbn1cblxuXG5cbi8vIEV4cG9ydCB0aGUgaGFuZGxlcnMgZm9yIHVzZSBieSBHbyB2aWEgaW5kZXgudHNcbmV4cG9ydCB7IGhhbmRsZURyYWdFbnRlciwgaGFuZGxlRHJhZ0xlYXZlLCBoYW5kbGVEcmFnT3ZlciB9O1xuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHBvc2l0aW9uIG9mIGEgd2luZG93LlxuICovXG5pbnRlcmZhY2UgUG9zaXRpb24ge1xuICAgIC8qKiBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xuICAgIHg6IG51bWJlcjtcbiAgICAvKiogVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXG4gICAgeTogbnVtYmVyO1xufVxuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXG4gKi9cbmludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuICovXG4gICAgd2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LiAqL1xuICAgIGhlaWdodDogbnVtYmVyO1xufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY2FsbGVyU3ltID0gU3ltYm9sKFwiY2FsbGVyXCIpO1xuXG5jbGFzcyBXaW5kb3cge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIHByaXZhdGUgW2NhbGxlclN5bV06IChtZXNzYWdlOiBudW1iZXIsIGFyZ3M/OiBhbnkpID0+IFByb21pc2U8YW55PjtcblxuICAgIC8qKlxuICAgICAqIEluaXRpYWxpc2VzIGEgd2luZG93IG9iamVjdCB3aXRoIHRoZSBzcGVjaWZpZWQgbmFtZS5cbiAgICAgKlxuICAgICAqIEBwcml2YXRlXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgdGFyZ2V0IHdpbmRvdy5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcgPSAnJykge1xuICAgICAgICB0aGlzW2NhbGxlclN5bV0gPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcblxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xuICAgICAgICAgICAgaWYgKFxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXG4gICAgICAgICAgICApIHtcbiAgICAgICAgICAgICAgICAodGhpcyBhcyBhbnkpW21ldGhvZF0gPSAodGhpcyBhcyBhbnkpW21ldGhvZF0uYmluZCh0aGlzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gZ2V0LlxuICAgICAqIEByZXR1cm5zIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXG4gICAgICovXG4gICAgR2V0KG5hbWU6IHN0cmluZyk6IFdpbmRvdyB7XG4gICAgICAgIHJldHVybiBuZXcgV2luZG93KG5hbWUpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2VudGVycyB0aGUgd2luZG93IG9uIHRoZSBzY3JlZW4uXG4gICAgICovXG4gICAgQ2VudGVyKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENlbnRlck1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgQ2xvc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2xvc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKi9cbiAgICBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEVuYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqL1xuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBGb2N1cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb2N1c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKi9cbiAgICBGb3JjZVJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb3JjZVJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU3dpdGNoZXMgdGhlIHdpbmRvdyB0byBmdWxsc2NyZWVuIG1vZGUuXG4gICAgICovXG4gICAgRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBzY3JlZW4gdGhhdCB0aGUgd2luZG93IGlzIG9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvbi5cbiAgICAgKi9cbiAgICBHZXRTY3JlZW4oKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRTY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgem9vbSBsZXZlbC5cbiAgICAgKi9cbiAgICBHZXRab29tKCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0Wm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBIZWlnaHQoKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIZWlnaHRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEhpZGVzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIaWRlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkLlxuICAgICAqL1xuICAgIElzRm9jdXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0ZvY3VzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW4uXG4gICAgICovXG4gICAgSXNGdWxsc2NyZWVuKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtYXhpbWlzZWQuXG4gICAgICovXG4gICAgSXNNYXhpbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNYXhpbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkLlxuICAgICAqL1xuICAgIElzTWluaW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWluaW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogTWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTmFtZSgpOiBQcm9taXNlPHN0cmluZz4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE5hbWVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxuICAgICAqL1xuICAgIE9wZW5EZXZUb29scygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShPcGVuRGV2VG9vbHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUmVsYXRpdmVQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXG4gICAgICovXG4gICAgUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGUuXG4gICAgICovXG4gICAgUmVzaXphYmxlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc2l6YWJsZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVzdG9yZXMgdGhlIHdpbmRvdyB0byBpdHMgcHJldmlvdXMgc3RhdGUgaWYgaXQgd2FzIHByZXZpb3VzbHkgbWluaW1pc2VkLCBtYXhpbWlzZWQgb3IgZnVsbHNjcmVlbi5cbiAgICAgKi9cbiAgICBSZXN0b3JlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc3RvcmVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0UG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxuICAgICAqXG4gICAgICogQHBhcmFtIGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcC5cbiAgICAgKi9cbiAgICBTZXRBbHdheXNPblRvcChhbHdheXNPblRvcDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gciAtIFRoZSBkZXNpcmVkIHJlZCBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYSAtIFRoZSBkZXNpcmVkIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICovXG4gICAgU2V0QmFja2dyb3VuZENvbG91cihyOiBudW1iZXIsIGc6IG51bWJlciwgYjogbnVtYmVyLCBhOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3MuXG4gICAgICovXG4gICAgU2V0RnJhbWVsZXNzKGZyYW1lbGVzczogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWQuXG4gICAgICovXG4gICAgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0TWF4U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWluaW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRNaW5TaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWluU2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cbiAgICAgKi9cbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gdGl0bGUgLSBUaGUgZGVzaXJlZCB0aXRsZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFRpdGxlKHRpdGxlOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWwuXG4gICAgICovXG4gICAgU2V0Wm9vbSh6b29tOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNob3dNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaXplKCk6IFByb21pc2U8U2l6ZT4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNpemVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZ1bGxzY3JlZW4gYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZU1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmcmFtZWxlc3MgYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGcmFtZWxlc3MoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnJhbWVsZXNzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLW1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5NaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBXaWR0aCgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFdpZHRoTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBab29tcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFpvb20oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbUluKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21Jbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbU91dCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tT3V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tUmVzZXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbVJlc2V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUy9MaW51eCBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXG4gICAgICogR2F0aGVycyBpbmZvcm1hdGlvbiBhYm91dCB0aGUgZHJvcCB0YXJnZXQgZWxlbWVudCBhbmQgc2VuZHMgaXQgYmFjayB0byB0aGUgR28gYmFja2VuZC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cbiAgICAgKiBAcGFyYW0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQsIGluIGxvZ2ljYWwgKENTUykgcGl4ZWxzIHJlbGF0aXZlIHRvIHRoZSB3ZWJ2aWV3LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudCwgaW4gbG9naWNhbCAoQ1NTKSBwaXhlbHMgcmVsYXRpdmUgdG8gdGhlIHdlYnZpZXcuXG4gICAgICovXG4gICAgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgICAgICAvLyBDaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBpZ25vcmUgdGhlIGRyb3BcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudChlbGVtZW50KTtcblxuICAgICAgICBpZiAoIWRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIC8vIERyb3Agd2FzIG5vdCBvbiBhIGRlc2lnbmF0ZWQgZHJvcCB0YXJnZXQgLSBpZ25vcmVcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnN0IGVsZW1lbnREZXRhaWxzID0ge1xuICAgICAgICAgICAgaWQ6IGRyb3BUYXJnZXQuaWQsXG4gICAgICAgICAgICBjbGFzc0xpc3Q6IEFycmF5LmZyb20oZHJvcFRhcmdldC5jbGFzc0xpc3QpLFxuICAgICAgICAgICAgYXR0cmlidXRlczoge30gYXMgeyBba2V5OiBzdHJpbmddOiBzdHJpbmcgfSxcbiAgICAgICAgfTtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgICAgIGNvbnN0IGF0dHIgPSBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXNbaV07XG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscy5hdHRyaWJ1dGVzW2F0dHIubmFtZV0gPSBhdHRyLnZhbHVlO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgICAgIGZpbGVuYW1lcyxcbiAgICAgICAgICAgIHgsXG4gICAgICAgICAgICB5LFxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMsXG4gICAgICAgIH07XG5cbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dKEZpbGVzRHJvcHBlZCwgcGF5bG9hZCk7XG4gICAgICAgIFxuICAgICAgICAvLyBDbGVhbiB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhZnRlciBkcm9wXG4gICAgICAgIGNsZWFudXBOYXRpdmVEcmFnKCk7XG4gICAgfVxuICBcbiAgICAvKipcbiAgICAgKiBNb3ZlcyB0aGUgd2luZG93IHRvIHRoZSBjZW50ZXIgb2YgdGhlIHNwZWNpZmllZCBzY3JlZW4ncyB3b3JrIGFyZWEuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gc2NyZWVuSUQgLSBUaGUgSUQgb2YgdGhlIHRhcmdldCBzY3JlZW4uXG4gICAgICovXG4gICAgU2V0U2NyZWVuKHNjcmVlbklEOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRTY3JlZW5NZXRob2QsIHsgc2NyZWVuSUQgfSk7XG4gICAgfVxuXG4gICAgLyogVHJpZ2dlcnMgV2luZG93cyAxMSBTbmFwIEFzc2lzdCBmZWF0dXJlIChXaW5kb3dzIG9ubHkpLlxuICAgICAqIFRoaXMgaXMgZXF1aXZhbGVudCB0byBwcmVzc2luZyBXaW4rWiBhbmQgc2hvd3Mgc25hcCBsYXlvdXQgb3B0aW9ucy5cbiAgICAgKi9cbiAgICBTbmFwQXNzaXN0KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNuYXBBc3Npc3RNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBwcmludCBkaWFsb2cgZm9yIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUHJpbnQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUHJpbnRNZXRob2QpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXG4gKi9cbmNvbnN0IHRoaXNXaW5kb3cgPSBuZXcgV2luZG93KCcnKTtcblxuLyoqXG4gKiBTZXRzIHVwIGdsb2JhbCBkcmFnIGFuZCBkcm9wIGV2ZW50IGxpc3RlbmVycyBmb3IgZmlsZSBkcm9wcy5cbiAqIEhhbmRsZXMgdmlzdWFsIGZlZWRiYWNrIChob3ZlciBzdGF0ZSkgYW5kIGZpbGUgZHJvcCBwcm9jZXNzaW5nLlxuICovXG5mdW5jdGlvbiBzZXR1cERyb3BUYXJnZXRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZG9jRWxlbWVudCA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudDtcbiAgICBsZXQgZHJhZ0VudGVyQ291bnRlciA9IDA7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdlbnRlcicsIChldmVudCkgPT4ge1xuICAgICAgICBpZiAoIWV2ZW50LmRhdGFUcmFuc2Zlcj8udHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcbiAgICAgICAgICAgIHJldHVybjsgLy8gT25seSBoYW5kbGUgZmlsZSBkcmFncywgbGV0IG90aGVyIGRyYWdzIHBhc3MgdGhyb3VnaFxuICAgICAgICB9XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIEFsd2F5cyBwcmV2ZW50IGRlZmF1bHQgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnOyAvLyBTaG93IFwibm8gZHJvcFwiIGN1cnNvclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcbiAgICAgICAgfVxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyKys7XG4gICAgICAgIFxuICAgICAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZKTtcbiAgICAgICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xuXG4gICAgICAgIC8vIFVwZGF0ZSBob3ZlciBzdGF0ZVxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmIChkcm9wVGFyZ2V0KSB7XG4gICAgICAgICAgICBkcm9wVGFyZ2V0LmNsYXNzTGlzdC5hZGQoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBkcm9wVGFyZ2V0O1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdvdmVyJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyYWdzXG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7IC8vIFNob3cgXCJubyBkcm9wXCIgY3Vyc29yXG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICAvLyBVcGRhdGUgZHJvcCB0YXJnZXQgYXMgY3Vyc29yIG1vdmVzXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XG4gICAgICAgIFxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgaWYgKGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGlmICghZHJvcFRhcmdldC5jbGFzc0xpc3QuY29udGFpbnMoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKSkge1xuICAgICAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2xlYXZlJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7IC8vIEFsd2F5cyBwcmV2ZW50IGRlZmF1bHQgdG8gc3RvcCBicm93c2VyIG5hdmlnYXRpb25cbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgLy8gT24gTGludXgvV2ViS2l0R1RLIGFuZCBtYWNPUywgZHJhZ2xlYXZlIGZpcmVzIGltbWVkaWF0ZWx5IHdpdGggcmVsYXRlZFRhcmdldD1udWxsIHdoZW4gbmF0aXZlXG4gICAgICAgIC8vIGRyYWcgaGFuZGxpbmcgaXMgaW52b2x2ZWQuIElnbm9yZSB0aGVzZSBzcHVyaW91cyBldmVudHMgLSB3ZSdsbCBjbGVhbiB1cCBvbiBkcm9wIGluc3RlYWQuXG4gICAgICAgIGlmIChldmVudC5yZWxhdGVkVGFyZ2V0ID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXItLTtcbiAgICAgICAgXG4gICAgICAgIGlmIChkcmFnRW50ZXJDb3VudGVyID09PSAwIHx8IFxuICAgICAgICAgICAgKGN1cnJlbnREcm9wVGFyZ2V0ICYmICFjdXJyZW50RHJvcFRhcmdldC5jb250YWlucyhldmVudC5yZWxhdGVkVGFyZ2V0IGFzIE5vZGUpKSkge1xuICAgICAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XG4gICAgICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwO1xuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcm9wJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyb3BzXG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIgPSAwO1xuICAgICAgICBcbiAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBoYW5kbGUgZmlsZSBkcm9wcyB2aWEgSmF2YVNjcmlwdFxuICAgICAgICAvLyBPbiBtYWNPUy9MaW51eCwgbmF0aXZlIGNvZGUgd2lsbCBjYWxsIEhhbmRsZVBsYXRmb3JtRmlsZURyb3BcbiAgICAgICAgaWYgKGNhblJlc29sdmVGaWxlUGF0aHMoKSkge1xuICAgICAgICAgICAgY29uc3QgZmlsZXM6IEZpbGVbXSA9IFtdO1xuICAgICAgICAgICAgaWYgKGV2ZW50LmRhdGFUcmFuc2Zlci5pdGVtcykge1xuICAgICAgICAgICAgICAgIGZvciAoY29uc3QgaXRlbSBvZiBldmVudC5kYXRhVHJhbnNmZXIuaXRlbXMpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKGl0ZW0ua2luZCA9PT0gJ2ZpbGUnKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCBmaWxlID0gaXRlbS5nZXRBc0ZpbGUoKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChmaWxlKSBmaWxlcy5wdXNoKGZpbGUpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSBlbHNlIGlmIChldmVudC5kYXRhVHJhbnNmZXIuZmlsZXMpIHtcbiAgICAgICAgICAgICAgICBmb3IgKGNvbnN0IGZpbGUgb2YgZXZlbnQuZGF0YVRyYW5zZmVyLmZpbGVzKSB7XG4gICAgICAgICAgICAgICAgICAgIGZpbGVzLnB1c2goZmlsZSk7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgXG4gICAgICAgICAgICBpZiAoZmlsZXMubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgICAgIHJlc29sdmVGaWxlUGF0aHMoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZmlsZXMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xufVxuXG4vLyBJbml0aWFsaXplIGxpc3RlbmVycyB3aGVuIHRoZSBzY3JpcHQgbG9hZHNcbmlmICh0eXBlb2Ygd2luZG93ICE9PSBcInVuZGVmaW5lZFwiICYmIHR5cGVvZiBkb2N1bWVudCAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgIHNldHVwRHJvcFRhcmdldExpc3RlbmVycygpO1xufVxuXG5leHBvcnQgZGVmYXVsdCB0aGlzV2luZG93O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgKiBhcyBSdW50aW1lIGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyY1wiO1xuXG4vLyBOT1RFOiB0aGUgZm9sbG93aW5nIG1ldGhvZHMgTVVTVCBiZSBpbXBvcnRlZCBleHBsaWNpdGx5IGJlY2F1c2Ugb2YgaG93IGVzYnVpbGQgaW5qZWN0aW9uIHdvcmtzXG5pbXBvcnQgeyBFbmFibGUgYXMgRW5hYmxlV01MIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dtbFwiO1xuaW1wb3J0IHsgZGVidWdMb2cgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHNcIjtcblxud2luZG93LndhaWxzID0gUnVudGltZTtcbkVuYWJsZVdNTCgpO1xuXG5pZiAoREVCVUcpIHtcbiAgICBkZWJ1Z0xvZyhcIldhaWxzIFJ1bnRpbWUgTG9hZGVkXCIpXG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TeXN0ZW0pO1xuXG5jb25zdCBTeXN0ZW1Jc0RhcmtNb2RlID0gMDtcbmNvbnN0IFN5c3RlbUVudmlyb25tZW50ID0gMTtcbmNvbnN0IFN5c3RlbUNhcGFiaWxpdGllcyA9IDI7XG5cbmNvbnN0IF9pbnZva2UgPSAoZnVuY3Rpb24gKCkge1xuICAgIHRyeSB7XG4gICAgICAgIC8vIFdpbmRvd3MgV2ViVmlldzJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5jaHJvbWU/LndlYnZpZXc/LnBvc3RNZXNzYWdlKSB7XG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3KTtcbiAgICAgICAgfVxuICAgICAgICAvLyBtYWNPUy9pT1MgV0tXZWJWaWV3XG4gICAgICAgIGVsc2UgaWYgKCh3aW5kb3cgYXMgYW55KS53ZWJraXQ/Lm1lc3NhZ2VIYW5kbGVycz8uWydleHRlcm5hbCddPy5wb3N0TWVzc2FnZSkge1xuICAgICAgICAgICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddLnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10pO1xuICAgICAgICB9XG4gICAgICAgIC8vIEFuZHJvaWQgV2ViVmlldyAtIHVzZXMgYWRkSmF2YXNjcmlwdEludGVyZmFjZSB3aGljaCBleHBvc2VzIHdpbmRvdy53YWlscy5pbnZva2VcbiAgICAgICAgZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndhaWxzPy5pbnZva2UpIHtcbiAgICAgICAgICAgIHJldHVybiAobXNnOiBhbnkpID0+ICh3aW5kb3cgYXMgYW55KS53YWlscy5pbnZva2UodHlwZW9mIG1zZyA9PT0gJ3N0cmluZycgPyBtc2cgOiBKU09OLnN0cmluZ2lmeShtc2cpKTtcbiAgICAgICAgfVxuICAgIH0gY2F0Y2goZSkge31cblxuICAgIGNvbnNvbGUud2FybignXFxuJWNcdTI2QTBcdUZFMEYgQnJvd3NlciBFbnZpcm9ubWVudCBEZXRlY3RlZCAlY1xcblxcbiVjT25seSBVSSBwcmV2aWV3cyBhcmUgYXZhaWxhYmxlIGluIHRoZSBicm93c2VyLiBGb3IgZnVsbCBmdW5jdGlvbmFsaXR5LCBwbGVhc2UgcnVuIHRoZSBhcHBsaWNhdGlvbiBpbiBkZXNrdG9wIG1vZGUuXFxuTW9yZSBpbmZvcm1hdGlvbiBhdDogaHR0cHM6Ly92My53YWlscy5pby9sZWFybi9idWlsZC8jdXNpbmctYS1icm93c2VyLWZvci1kZXZlbG9wbWVudFxcbicsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjZmZmZmZmOyBjb2xvcjogIzAwMDAwMDsgZm9udC13ZWlnaHQ6IGJvbGQ7IHBhZGRpbmc6IDRweCA4cHg7IGJvcmRlci1yYWRpdXM6IDRweDsgYm9yZGVyOiAycHggc29saWQgIzAwMDAwMDsnLFxuICAgICAgICAnYmFja2dyb3VuZDogdHJhbnNwYXJlbnQ7JyxcbiAgICAgICAgJ2NvbG9yOiAjZmZmZmZmOyBmb250LXN0eWxlOiBpdGFsaWM7IGZvbnQtd2VpZ2h0OiBib2xkOycpO1xuICAgIHJldHVybiBudWxsO1xufSkoKTtcblxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2c6IGFueSk6IHZvaWQge1xuICAgIF9pbnZva2U/Lihtc2cpO1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyB0aGUgc3lzdGVtIGRhcmsgbW9kZSBzdGF0dXMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYSBib29sZWFuIHZhbHVlIGluZGljYXRpbmcgaWYgdGhlIHN5c3RlbSBpcyBpbiBkYXJrIG1vZGUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgIHJldHVybiBjYWxsKFN5c3RlbUlzRGFya01vZGUpO1xufVxuXG4vKipcbiAqIEZldGNoZXMgdGhlIGNhcGFiaWxpdGllcyBvZiB0aGUgYXBwbGljYXRpb24gZnJvbSB0aGUgc2VydmVyLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIHRoZSBjYXBhYmlsaXRpZXMuXG4gKi9cbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKTogUHJvbWlzZTxSZWNvcmQ8c3RyaW5nLCBhbnk+PiB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtQ2FwYWJpbGl0aWVzKTtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBPU0luZm8ge1xuICAgIC8qKiBUaGUgYnJhbmRpbmcgb2YgdGhlIE9TLiAqL1xuICAgIEJyYW5kaW5nOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBJRCBvZiB0aGUgT1MuICovXG4gICAgSUQ6IHN0cmluZztcbiAgICAvKiogVGhlIG5hbWUgb2YgdGhlIE9TLiAqL1xuICAgIE5hbWU6IHN0cmluZztcbiAgICAvKiogVGhlIHZlcnNpb24gb2YgdGhlIE9TLiAqL1xuICAgIFZlcnNpb246IHN0cmluZztcbn1cblxuZXhwb3J0IGludGVyZmFjZSBFbnZpcm9ubWVudEluZm8ge1xuICAgIC8qKiBUaGUgYXJjaGl0ZWN0dXJlIG9mIHRoZSBzeXN0ZW0uICovXG4gICAgQXJjaDogc3RyaW5nO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS4gKi9cbiAgICBEZWJ1ZzogYm9vbGVhbjtcbiAgICAvKiogVGhlIG9wZXJhdGluZyBzeXN0ZW0gaW4gdXNlLiAqL1xuICAgIE9TOiBzdHJpbmc7XG4gICAgLyoqIERldGFpbHMgb2YgdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXG4gICAgT1NJbmZvOiBPU0luZm87XG4gICAgLyoqIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uICovXG4gICAgUGxhdGZvcm1JbmZvOiBSZWNvcmQ8c3RyaW5nLCBhbnk+O1xufVxuXG4vKipcbiAqIFJldHJpZXZlcyBlbnZpcm9ubWVudCBkZXRhaWxzLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIE9TIGFuZCBzeXN0ZW0gYXJjaGl0ZWN0dXJlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKTogUHJvbWlzZTxFbnZpcm9ubWVudEluZm8+IHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1FbnZpcm9ubWVudCk7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cbiAqXG4gKiBAcmV0dXJuIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cywgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwid2luZG93c1wiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LlxuICpcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIExpbnV4LCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwibGludXhcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgYSBtYWNPUyBvcGVyYXRpbmcgc3lzdGVtLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIG1hY09TLCBmYWxzZSBvdGhlcndpc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TID09PSBcImRhcndpblwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIGlPUy5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBvcGVyYXRpbmcgc3lzdGVtIGlzIGlPUywgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNJT1MoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJpb3NcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBBbmRyb2lkLlxuICpcbiAqIEByZXR1cm5zIFRydWUgaWYgdGhlIG9wZXJhdGluZyBzeXN0ZW0gaXMgQW5kcm9pZCwgb3RoZXJ3aXNlIGZhbHNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBbmRyb2lkKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwiYW5kcm9pZFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgYXBwIGlzIHJ1bm5pbmcgb24gYSBtb2JpbGUgT1MgKGlPUyBvciBBbmRyb2lkKS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIG9uIGlPUyBvciBBbmRyb2lkLCBvdGhlcndpc2UgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc01vYmlsZSgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gSXNJT1MoKSB8fCBJc0FuZHJvaWQoKTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGFwcCBpcyBydW5uaW5nIG9uIGEgZGVza3RvcCBPUyAobWFjT1MsIFdpbmRvd3Mgb3IgTGludXgpLlxuICpcbiAqIEByZXR1cm5zIFRydWUgb24gbWFjT1MsIFdpbmRvd3Mgb3IgTGludXgsIG90aGVyd2lzZSBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGVza3RvcCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gSXNNYWMoKSB8fCBJc1dpbmRvd3MoKSB8fCBJc0xpbnV4KCk7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYW1kNjRcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhcm1cIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxuICpcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFybTY0XCI7XG59XG5cbi8qKlxuICogUmVwb3J0cyB3aGV0aGVyIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gQm9vbGVhbigod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uRGVidWcpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgSXNEZWJ1ZyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xuaW1wb3J0IHsgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuXG4vLyBzZXR1cFxuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxuaWYgKGhhc0RPTSkge1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG59XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcblxuY29uc3QgQ29udGV4dE1lbnVPcGVuID0gMDtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkOiBzdHJpbmcsIHg6IG51bWJlciwgeTogbnVtYmVyLCBkYXRhOiBhbnkpOiB2b2lkIHtcbiAgICB2b2lkIGNhbGwoQ29udGV4dE1lbnVPcGVuLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xuXG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcbiAgICBjb25zdCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpLnRyaW0oKTtcblxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBjb25zdCBkYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnUtZGF0YVwiKTtcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBkYXRhKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50LCB0YXJnZXQpO1xuICAgIH1cbn1cblxuXG4vKlxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBhdXRvOyAoZGVmYXVsdCkgd2lsbCBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudSBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZSBPUiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkIE9SIGVsZW1lbnQgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcblxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcbiovXG5mdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50OiBNb3VzZUV2ZW50LCB0YXJnZXQ6IEhUTUxFbGVtZW50KSB7XG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XG4gICAgaWYgKElzRGVidWcoKSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIHN3aXRjaCAod2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKSkge1xuICAgICAgICBjYXNlICdzaG93JzpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSAnaGlkZSc6XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgaWYgKHRhcmdldC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxuICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSBzZWxlY3Rpb24gJiYgc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMDtcbiAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcbiAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XG4gICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IHRhcmdldCkge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgdGFnIGlzIGlucHV0IG9yIHRleHRhcmVhLlxuICAgIGlmICh0YXJnZXQgaW5zdGFuY2VvZiBIVE1MSW5wdXRFbGVtZW50IHx8IHRhcmdldCBpbnN0YW5jZW9mIEhUTUxUZXh0QXJlYUVsZW1lbnQpIHtcbiAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIXRhcmdldC5yZWFkT25seSAmJiAhdGFyZ2V0LmRpc2FibGVkKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXG4gKlxuICogQHBhcmFtIGtleSAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cbiAqIEByZXR1cm4gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5OiBzdHJpbmcpOiBhbnkge1xuICAgIHRyeSB7XG4gICAgICAgIHJldHVybiB3aW5kb3cuX3dhaWxzLmZsYWdzW2tleV07XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJVbmFibGUgdG8gcmV0cmlldmUgZmxhZyAnXCIgKyBrZXkgKyBcIic6IFwiICsgZSwgeyBjYXVzZTogZSB9KTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IGludm9rZSwgSXNXaW5kb3dzLCBJc0xpbnV4IH0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XG5pbXBvcnQgeyBHZXRGbGFnIH0gZnJvbSBcIi4vZmxhZ3MuanNcIjtcbmltcG9ydCB7IGNhblRyYWNrQnV0dG9ucywgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuXG4vLyBTZXR1cFxubGV0IGNhbkRyYWcgPSBmYWxzZTtcbmxldCBkcmFnZ2luZyA9IGZhbHNlO1xuXG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5sZXQgY2FuUmVzaXplID0gZmFsc2U7XG5sZXQgcmVzaXppbmcgPSBmYWxzZTtcbmxldCByZXNpemVFZGdlOiBzdHJpbmcgPSBcIlwiO1xubGV0IGRlZmF1bHRDdXJzb3IgPSBcImF1dG9cIjtcblxubGV0IGJ1dHRvbnMgPSAwO1xuaW1wb3J0IHsgaGFzRE9NIH0gZnJvbSBcIi4vZW52aXJvbm1lbnQuanNcIjtcblxubGV0IGJ1dHRvbnNUcmFja2VkID0gZmFsc2U7XG5cbmlmIChoYXNET00pIHtcbiAgICBidXR0b25zVHJhY2tlZCA9IGNhblRyYWNrQnV0dG9ucygpO1xuICAgIHdpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xuICAgIHdpbmRvdy5fd2FpbHMuc2V0UmVzaXphYmxlID0gKHZhbHVlOiBib29sZWFuKTogdm9pZCA9PiB7XG4gICAgICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xuICAgICAgICBpZiAoIXJlc2l6YWJsZSkge1xuICAgICAgICAgICAgLy8gU3RvcCByZXNpemluZyBpZiBpbiBwcm9ncmVzcy5cbiAgICAgICAgICAgIGNhblJlc2l6ZSA9IHJlc2l6aW5nID0gZmFsc2U7XG4gICAgICAgICAgICBzZXRSZXNpemUoKTtcbiAgICAgICAgfVxuICAgIH07XG59XG5cbi8vIERlZmVyIGF0dGFjaGluZyBtb3VzZSBsaXN0ZW5lcnMgdW50aWwgd2Uga25vdyB3ZSdyZSBub3Qgb24gbW9iaWxlLlxubGV0IGRyYWdJbml0RG9uZSA9IGZhbHNlO1xuZnVuY3Rpb24gaXNNb2JpbGUoKTogYm9vbGVhbiB7XG4gICAgY29uc3Qgb3MgPSAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1M7XG4gICAgaWYgKG9zID09PSBcImlvc1wiIHx8IG9zID09PSBcImFuZHJvaWRcIikgcmV0dXJuIHRydWU7XG4gICAgLy8gRmFsbGJhY2sgaGV1cmlzdGljIGlmIGVudmlyb25tZW50IG5vdCB5ZXQgc2V0XG4gICAgY29uc3QgdWEgPSBuYXZpZ2F0b3IudXNlckFnZW50IHx8IG5hdmlnYXRvci52ZW5kb3IgfHwgKHdpbmRvdyBhcyBhbnkpLm9wZXJhIHx8IFwiXCI7XG4gICAgcmV0dXJuIC9hbmRyb2lkfGlwaG9uZXxpcGFkfGlwb2R8aWVtb2JpbGV8d3BkZXNrdG9wL2kudGVzdCh1YSk7XG59XG5mdW5jdGlvbiB0cnlJbml0RHJhZ0hhbmRsZXJzKCk6IHZvaWQge1xuICAgIGlmIChkcmFnSW5pdERvbmUpIHJldHVybjtcbiAgICBpZiAoaXNNb2JpbGUoKSkgcmV0dXJuO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICBmb3IgKGNvbnN0IGV2IG9mIFsnY2xpY2snLCAnY29udGV4dG1lbnUnLCAnZGJsY2xpY2snXSkge1xuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihldiwgc3VwcHJlc3NFdmVudCwgeyBjYXB0dXJlOiB0cnVlIH0pO1xuICAgIH1cbiAgICBkcmFnSW5pdERvbmUgPSB0cnVlO1xufVxuaWYgKGhhc0RPTSkge1xuICAgIC8vIEF0dGVtcHQgaW1tZWRpYXRlIGluaXQgKGluIGNhc2UgZW52aXJvbm1lbnQgYWxyZWFkeSBwcmVzZW50KVxuICAgIHRyeUluaXREcmFnSGFuZGxlcnMoKTtcbiAgICAvLyBBbHNvIGF0dGVtcHQgb24gRE9NIHJlYWR5XG4gICAgZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignRE9NQ29udGVudExvYWRlZCcsIHRyeUluaXREcmFnSGFuZGxlcnMsIHsgb25jZTogdHJ1ZSB9KTtcbiAgICAvLyBBcyBhIGxhc3QgcmVzb3J0LCBwb2xsIGZvciBlbnZpcm9ubWVudCBmb3IgYSBzaG9ydCBwZXJpb2RcbiAgICBsZXQgZHJhZ0VudlBvbGxzID0gMDtcbiAgICBjb25zdCBkcmFnRW52UG9sbCA9IHdpbmRvdy5zZXRJbnRlcnZhbCgoKSA9PiB7XG4gICAgICAgIGlmIChkcmFnSW5pdERvbmUpIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyByZXR1cm47IH1cbiAgICAgICAgdHJ5SW5pdERyYWdIYW5kbGVycygpO1xuICAgICAgICBpZiAoKytkcmFnRW52UG9sbHMgPiAxMDApIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyB9XG4gICAgfSwgNTApO1xufVxuXG5mdW5jdGlvbiBzdXBwcmVzc0V2ZW50KGV2ZW50OiBFdmVudCkge1xuICAgIC8vIFN1cHByZXNzIGNsaWNrIGV2ZW50cyB3aGlsZSByZXNpemluZyBvciBkcmFnZ2luZy5cbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn1cblxuLy8gVXNlIGNvbnN0YW50cyB0byBhdm9pZCBjb21wYXJpbmcgc3RyaW5ncyBtdWx0aXBsZSB0aW1lcy5cbmNvbnN0IE1vdXNlRG93biA9IDA7XG5jb25zdCBNb3VzZVVwICAgPSAxO1xuY29uc3QgTW91c2VNb3ZlID0gMjtcblxuZnVuY3Rpb24gdXBkYXRlKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgLy8gV2luZG93cyBzdXBwcmVzc2VzIG1vdXNlIGV2ZW50cyBhdCB0aGUgZW5kIG9mIGRyYWdnaW5nIG9yIHJlc2l6aW5nLFxuICAgIC8vIHNvIHdlIG5lZWQgdG8gYmUgc21hcnQgYW5kIHN5bnRoZXNpemUgYnV0dG9uIGV2ZW50cy5cblxuICAgIGxldCBldmVudFR5cGU6IG51bWJlciwgZXZlbnRCdXR0b25zID0gZXZlbnQuYnV0dG9ucztcbiAgICBzd2l0Y2ggKGV2ZW50LnR5cGUpIHtcbiAgICAgICAgY2FzZSAnbW91c2Vkb3duJzpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlRG93bjtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyB8ICgxIDw8IGV2ZW50LmJ1dHRvbik7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICBjYXNlICdtb3VzZXVwJzpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlVXA7XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgJiB+KDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZU1vdmU7XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnM7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgIH1cblxuICAgIGxldCByZWxlYXNlZCA9IGJ1dHRvbnMgJiB+ZXZlbnRCdXR0b25zO1xuICAgIGxldCBwcmVzc2VkID0gZXZlbnRCdXR0b25zICYgfmJ1dHRvbnM7XG5cbiAgICBidXR0b25zID0gZXZlbnRCdXR0b25zO1xuXG4gICAgLy8gU3ludGhlc2l6ZSBhIHJlbGVhc2UtcHJlc3Mgc2VxdWVuY2UgaWYgd2UgZGV0ZWN0IGEgcHJlc3Mgb2YgYW4gYWxyZWFkeSBwcmVzc2VkIGJ1dHRvbi5cbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZURvd24gJiYgIShwcmVzc2VkICYgZXZlbnQuYnV0dG9uKSkge1xuICAgICAgICByZWxlYXNlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xuICAgICAgICBwcmVzc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XG4gICAgfVxuXG4gICAgLy8gU3VwcHJlc3MgYWxsIGJ1dHRvbiBldmVudHMgZHVyaW5nIGRyYWdnaW5nIGFuZCByZXNpemluZyxcbiAgICAvLyB1bmxlc3MgdGhpcyBpcyBhIG1vdXNldXAgZXZlbnQgdGhhdCBpcyBlbmRpbmcgYSBkcmFnIGFjdGlvbi5cbiAgICBpZiAoXG4gICAgICAgIGV2ZW50VHlwZSAhPT0gTW91c2VNb3ZlIC8vIEZhc3QgcGF0aCBmb3IgbW91c2Vtb3ZlXG4gICAgICAgICYmIHJlc2l6aW5nXG4gICAgICAgIHx8IChcbiAgICAgICAgICAgIGRyYWdnaW5nXG4gICAgICAgICAgICAmJiAoXG4gICAgICAgICAgICAgICAgZXZlbnRUeXBlID09PSBNb3VzZURvd25cbiAgICAgICAgICAgICAgICB8fCBldmVudC5idXR0b24gIT09IDBcbiAgICAgICAgICAgIClcbiAgICAgICAgKVxuICAgICkge1xuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxuXG4gICAgLy8gSGFuZGxlIHJlbGVhc2VzXG4gICAgaWYgKHJlbGVhc2VkICYgMSkgeyBwcmltYXJ5VXAoZXZlbnQpOyB9XG4gICAgLy8gSGFuZGxlIHByZXNzZXNcbiAgICBpZiAocHJlc3NlZCAmIDEpIHsgcHJpbWFyeURvd24oZXZlbnQpOyB9XG5cbiAgICAvLyBIYW5kbGUgbW91c2Vtb3ZlXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VNb3ZlKSB7IG9uTW91c2VNb3ZlKGV2ZW50KTsgfTtcbn1cblxuZnVuY3Rpb24gcHJpbWFyeURvd24oZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcbiAgICAvLyBSZXNldCByZWFkaW5lc3Mgc3RhdGUuXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xuXG4gICAgLy8gSWdub3JlIHJlcGVhdGVkIGNsaWNrcyBvbiBtYWNPUyBhbmQgTGludXguXG4gICAgaWYgKCFJc1dpbmRvd3MoKSkge1xuICAgICAgICBpZiAoZXZlbnQudHlwZSA9PT0gJ21vdXNlZG93bicgJiYgZXZlbnQuYnV0dG9uID09PSAwICYmIGV2ZW50LmRldGFpbCAhPT0gMSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgLy8gRG8gbm90IGFybSBlZGdlIHJlc2l6ZSBmcm9tIHN5bnRoZXNpemVkIHByZXNzZXMgb2JzZXJ2ZWQgb24gbW92ZS91cDpcbiAgICAgICAgLy8gZW50ZXJpbmcgdGhlIHdpbmRvdyB3aXRoIHRoZSBwcmltYXJ5IGJ1dHRvbiBhbHJlYWR5IGhlbGQgc2hvdWxkIG5vdFxuICAgICAgICAvLyBzdGVhbCBhbm90aGVyIGdlc3R1cmUgaW50byBhIHJlc2l6ZS5cbiAgICAgICAgaWYgKGV2ZW50LnR5cGUgIT09ICdtb3VzZWRvd24nKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cblxuICAgICAgICAvLyBSZWFkeSB0byByZXNpemUgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZS5cbiAgICAgICAgY2FuUmVzaXplID0gdHJ1ZTtcbiAgICAgICAgLy8gRG8gbm90IHN0YXJ0IGRyYWcgb3BlcmF0aW9ucyB3aGVuIG9uIHJlc2l6ZSBlZGdlcy5cbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFJldHJpZXZlIHRhcmdldCBlbGVtZW50XG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xuXG4gICAgLy8gUmVhZHkgdG8gZHJhZyBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lIG9uIGEgZHJhZ2dhYmxlIGVsZW1lbnQuXG4gICAgLy8gSWdub3JlIGNsaWNrcyBvbiB0aGUgc2Nyb2xsYmFyLlxuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KTtcbiAgICBjYW5EcmFnID0gKFxuICAgICAgICBzdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13YWlscy1kcmFnZ2FibGVcIikudHJpbSgpID09PSBcImRyYWdcIlxuICAgICAgICAmJiAoXG4gICAgICAgICAgICBldmVudC5vZmZzZXRYIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nTGVmdCkgPCB0YXJnZXQuY2xpZW50V2lkdGhcbiAgICAgICAgICAgICYmIGV2ZW50Lm9mZnNldFkgLSBwYXJzZUZsb2F0KHN0eWxlLnBhZGRpbmdUb3ApIDwgdGFyZ2V0LmNsaWVudEhlaWdodFxuICAgICAgICApXG4gICAgKTtcbn1cblxuZnVuY3Rpb24gcHJpbWFyeVVwKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgLy8gU3RvcCBkcmFnZ2luZyBhbmQgcmVzaXppbmcuXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xuICAgIGRyYWdnaW5nID0gZmFsc2U7XG4gICAgY2FuUmVzaXplID0gZmFsc2U7XG4gICAgcmVzaXppbmcgPSBmYWxzZTtcbn1cblxuY29uc3QgY3Vyc29yRm9yRWRnZSA9IE9iamVjdC5mcmVlemUoe1xuICAgIFwic2UtcmVzaXplXCI6IFwibndzZS1yZXNpemVcIixcbiAgICBcInN3LXJlc2l6ZVwiOiBcIm5lc3ctcmVzaXplXCIsXG4gICAgXCJudy1yZXNpemVcIjogXCJud3NlLXJlc2l6ZVwiLFxuICAgIFwibmUtcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcbiAgICBcInctcmVzaXplXCI6IFwiZXctcmVzaXplXCIsXG4gICAgXCJuLXJlc2l6ZVwiOiBcIm5zLXJlc2l6ZVwiLFxuICAgIFwicy1yZXNpemVcIjogXCJucy1yZXNpemVcIixcbiAgICBcImUtcmVzaXplXCI6IFwiZXctcmVzaXplXCIsXG59KVxuXG5mdW5jdGlvbiBzZXRSZXNpemUoZWRnZT86IGtleW9mIHR5cGVvZiBjdXJzb3JGb3JFZGdlKTogdm9pZCB7XG4gICAgaWYgKGVkZ2UpIHtcbiAgICAgICAgaWYgKCFyZXNpemVFZGdlKSB7IGRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvcjsgfVxuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGN1cnNvckZvckVkZ2VbZWRnZV07XG4gICAgfSBlbHNlIGlmICghZWRnZSAmJiByZXNpemVFZGdlKSB7XG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gZGVmYXVsdEN1cnNvcjtcbiAgICB9XG5cbiAgICByZXNpemVFZGdlID0gZWRnZSB8fCBcIlwiO1xufVxuXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xuICAgIGlmIChjYW5SZXNpemUgJiYgcmVzaXplRWRnZSkge1xuICAgICAgICAvLyBTdGFydCByZXNpemluZy5cbiAgICAgICAgcmVzaXppbmcgPSB0cnVlO1xuICAgICAgICBpbnZva2UoXCJ3YWlsczpyZXNpemU6XCIgKyByZXNpemVFZGdlKTtcbiAgICB9IGVsc2UgaWYgKGNhbkRyYWcpIHtcbiAgICAgICAgLy8gU3RhcnQgZHJhZ2dpbmcuXG4gICAgICAgIGRyYWdnaW5nID0gdHJ1ZTtcbiAgICAgICAgaW52b2tlKFwid2FpbHM6ZHJhZ1wiKTtcbiAgICB9XG5cbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcbiAgICAgICAgLy8gRWl0aGVyIGRyYWcgb3IgcmVzaXplIGlzIG9uZ29pbmcsXG4gICAgICAgIC8vIHJlc2V0IHJlYWRpbmVzcyBhbmQgc3RvcCBwcm9jZXNzaW5nLlxuICAgICAgICBjYW5EcmFnID0gY2FuUmVzaXplID0gZmFsc2U7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoIXJlc2l6YWJsZSB8fCAoIUlzV2luZG93cygpICYmICEoSXNMaW51eCgpICYmIEdldEZsYWcoXCJmcmFtZWxlc3NcIikpKSkge1xuICAgICAgICBpZiAocmVzaXplRWRnZSkgeyBzZXRSZXNpemUoKTsgfVxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgY29uc3QgcmVzaXplSGFuZGxlSGVpZ2h0ID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVIZWlnaHRcIikgfHwgNTtcbiAgICBjb25zdCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcblxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhcy5cbiAgICBjb25zdCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcblxuICAgIC8vIFdoZW4gYSBzY3JvbGxiYXIgaXMgcHJlc2VudCBhdCB0aGUgd2luZG93IGVkZ2UgaXQgY29uc3VtZXMgbW91c2UgZXZlbnRzIGluIHRoYXQgc3RyaXAuXG4gICAgLy8gU2hpZnQgdGhlIGVmZmVjdGl2ZSBjb250ZW50IGVkZ2UgaW53YXJkIHNvIHRoZSByZXNpemUgem9uZSBzaXRzIGp1c3QgYmVmb3JlIHRoZSBzY3JvbGxiYXIuXG4gICAgY29uc3Qgc2Nyb2xsYmFyV2lkdGggPSBNYXRoLm1heCgwLCB3aW5kb3cuaW5uZXJXaWR0aCAtIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5jbGllbnRXaWR0aCk7XG4gICAgY29uc3Qgc2Nyb2xsYmFySGVpZ2h0ID0gTWF0aC5tYXgoMCwgd2luZG93LmlubmVySGVpZ2h0IC0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LmNsaWVudEhlaWdodCk7XG4gICAgY29uc3QgcmlnaHRDb250ZW50RWRnZSA9IHdpbmRvdy5pbm5lcldpZHRoIC0gc2Nyb2xsYmFyV2lkdGg7XG4gICAgY29uc3QgYm90dG9tQ29udGVudEVkZ2UgPSB3aW5kb3cuaW5uZXJIZWlnaHQgLSBzY3JvbGxiYXJIZWlnaHQ7XG5cbiAgICBjb25zdCByaWdodEJvcmRlciA9IGV2ZW50LmNsaWVudFggPCByaWdodENvbnRlbnRFZGdlICYmIChyaWdodENvbnRlbnRFZGdlIC0gZXZlbnQuY2xpZW50WCkgPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBjb25zdCBsZWZ0Qm9yZGVyID0gZXZlbnQuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGNvbnN0IHRvcEJvcmRlciA9IGV2ZW50LmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG4gICAgY29uc3QgYm90dG9tQm9yZGVyID0gZXZlbnQuY2xpZW50WSA8IGJvdHRvbUNvbnRlbnRFZGdlICYmIChib3R0b21Db250ZW50RWRnZSAtIGV2ZW50LmNsaWVudFkpIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXIgYXJlYXMuXG4gICAgY29uc3QgcmlnaHRDb3JuZXIgPSBldmVudC5jbGllbnRYIDwgcmlnaHRDb250ZW50RWRnZSAmJiAocmlnaHRDb250ZW50RWRnZSAtIGV2ZW50LmNsaWVudFgpIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IGxlZnRDb3JuZXIgPSBldmVudC5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IHRvcENvcm5lciA9IGV2ZW50LmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IGJvdHRvbUNvcm5lciA9IGV2ZW50LmNsaWVudFkgPCBib3R0b21Db250ZW50RWRnZSAmJiAoYm90dG9tQ29udGVudEVkZ2UgLSBldmVudC5jbGllbnRZKSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG5cbiAgICBpZiAoIWxlZnRDb3JuZXIgJiYgIXRvcENvcm5lciAmJiAhYm90dG9tQ29ybmVyICYmICFyaWdodENvcm5lcikge1xuICAgICAgICAvLyBPcHRpbWlzYXRpb246IG91dCBvZiBhbGwgY29ybmVyIGFyZWFzIGltcGxpZXMgb3V0IG9mIGJvcmRlcnMuXG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbiAgICAvLyBEZXRlY3QgY29ybmVycy5cbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgLy8gRGV0ZWN0IGJvcmRlcnMuXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcbiAgICAvLyBPdXQgb2YgYm9yZGVyIGFyZWEuXG4gICAgZWxzZSBzZXRSZXNpemUoKTtcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgIF8gX19cbnwgfCAgIC8gL19fXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IGludm9rZSB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xuaW1wb3J0IHsgd2hlblJlYWR5IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcbmltcG9ydCB7IGhhc0RPTSB9IGZyb20gXCIuL2Vudmlyb25tZW50LmpzXCI7XG5cbnR5cGUgTm9uQ2xpZW50UmVnaW9uS2luZCA9IFwiY2FwdGlvblwiIHwgXCJtaW5pbWl6ZVwiIHwgXCJtYXhpbWl6ZVwiIHwgXCJjbG9zZVwiO1xuXG5pbnRlcmZhY2UgTm9uQ2xpZW50UmVnaW9uIHtcbiAgICBraW5kOiBOb25DbGllbnRSZWdpb25LaW5kO1xuICAgIGxlZnQ6IG51bWJlcjtcbiAgICB0b3A6IG51bWJlcjtcbiAgICByaWdodDogbnVtYmVyO1xuICAgIGJvdHRvbTogbnVtYmVyO1xufVxuXG4vKlxuLS13YWlscy1ub24tY2xpZW50LXJlZ2lvbjogY2FwdGlvbjsgIG1hcmtzIGFuIGFyZWEgdGhhdCBjYW4gZHJhZyB0aGUgd2luZG93XG4tLXdhaWxzLW5vbi1jbGllbnQtcmVnaW9uOiBtaW5pbWl6ZTsgbWFya3MgYSBjdXN0b20gbWluaW1pemUgYnV0dG9uXG4tLXdhaWxzLW5vbi1jbGllbnQtcmVnaW9uOiBtYXhpbWl6ZTsgbWFya3MgYSBjdXN0b20gbWF4aW1pemUgYnV0dG9uXG4tLXdhaWxzLW5vbi1jbGllbnQtcmVnaW9uOiBjbG9zZTsgICAgbWFya3MgYSBjdXN0b20gY2xvc2UgYnV0dG9uXG4qL1xuY29uc3QgcmVnaW9uUHJvcGVydHkgPSBcIi0td2FpbHMtbm9uLWNsaWVudC1yZWdpb25cIjtcbmNvbnN0IHJ1bnRpbWVDb25maWdSZWFkeUV2ZW50ID0gXCJ3YWlsczpydW50aW1lLWNvbmZpZy1yZWFkeVwiO1xuY29uc3QgdmFsaWRSZWdpb25zID0gbmV3IFNldDxOb25DbGllbnRSZWdpb25LaW5kPihbXCJjYXB0aW9uXCIsIFwibWluaW1pemVcIiwgXCJtYXhpbWl6ZVwiLCBcImNsb3NlXCJdKTtcblxuLy8gU2V0dXBcbmlmIChoYXNET00pIHtcbiAgICB3aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbn1cblxubGV0IHVwZGF0ZVBlbmRpbmcgPSBmYWxzZTtcbmxldCBsYXN0UGF5bG9hZCA9IFwiXCI7XG5sZXQgb2JzZXJ2ZWRFbGVtZW50cyA9IG5ldyBTZXQ8RWxlbWVudD4oKTtcbmxldCByZXNpemVPYnNlcnZlcjogUmVzaXplT2JzZXJ2ZXIgfCB1bmRlZmluZWQ7XG5sZXQgdHJhY2tpbmdTdGFydGVkID0gZmFsc2U7XG5cbmZ1bmN0aW9uIG5vcm1hbGlzZVJlZ2lvbktpbmQodmFsdWU6IHN0cmluZyk6IE5vbkNsaWVudFJlZ2lvbktpbmQgfCB1bmRlZmluZWQge1xuICAgIGNvbnN0IHJlZ2lvbiA9IHZhbHVlLnRyaW0oKS50b0xvd2VyQ2FzZSgpO1xuICAgIGlmICh2YWxpZFJlZ2lvbnMuaGFzKHJlZ2lvbiBhcyBOb25DbGllbnRSZWdpb25LaW5kKSkge1xuICAgICAgICByZXR1cm4gcmVnaW9uIGFzIE5vbkNsaWVudFJlZ2lvbktpbmQ7XG4gICAgfVxuICAgIHJldHVybiB1bmRlZmluZWQ7XG59XG5cbmZ1bmN0aW9uIG5vbkNsaWVudFJlZ2lvbkZvckVsZW1lbnQoZWxlbWVudDogRWxlbWVudCk6IE5vbkNsaWVudFJlZ2lvbktpbmQgfCB1bmRlZmluZWQge1xuICAgIGlmICghKGVsZW1lbnQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkpIHtcbiAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcbiAgICB9XG5cbiAgICBjb25zdCBzdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xuICAgIGNvbnN0IHJlZ2lvbiA9IG5vcm1hbGlzZVJlZ2lvbktpbmQoc3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShyZWdpb25Qcm9wZXJ0eSkpO1xuICAgIGlmICghcmVnaW9uKSB7XG4gICAgICAgIHJldHVybiB1bmRlZmluZWQ7XG4gICAgfVxuXG4gICAgY29uc3QgcGFyZW50ID0gZWxlbWVudC5wYXJlbnRFbGVtZW50O1xuICAgIGlmIChwYXJlbnQpIHtcbiAgICAgICAgY29uc3QgcGFyZW50U3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShwYXJlbnQpO1xuICAgICAgICAvLyBUaGUgQ1NTIHByb3BlcnR5IGlzIGluaGVyaXRlZC4gT25seSByZXBvcnQgdGhlIG91dGVybW9zdCBlbGVtZW50IGZvclxuICAgICAgICAvLyBlYWNoIGNvbnRpZ3VvdXMgcmVnaW9uIHNvIG5hdGl2ZSBoaXQgdGVzdGluZyBzZWVzIHN0YWJsZSByZWN0YW5nbGVzLlxuICAgICAgICBpZiAobm9ybWFsaXNlUmVnaW9uS2luZChwYXJlbnRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKHJlZ2lvblByb3BlcnR5KSkgPT09IHJlZ2lvbikge1xuICAgICAgICAgICAgcmV0dXJuIHVuZGVmaW5lZDtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIHJldHVybiByZWdpb247XG59XG5cbmZ1bmN0aW9uIGlzVmlzaWJsZShlbGVtZW50OiBIVE1MRWxlbWVudCk6IGJvb2xlYW4ge1xuICAgIGNvbnN0IHN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XG4gICAgcmV0dXJuIHN0eWxlLmRpc3BsYXkgIT09IFwibm9uZVwiICYmXG4gICAgICAgIHN0eWxlLnZpc2liaWxpdHkgIT09IFwiaGlkZGVuXCIgJiZcbiAgICAgICAgc3R5bGUuY29udGVudFZpc2liaWxpdHkgIT09IFwiaGlkZGVuXCI7XG59XG5cbmZ1bmN0aW9uIGVsZW1lbnRSZWdpb24oZWxlbWVudDogRWxlbWVudCk6IE5vbkNsaWVudFJlZ2lvbiB8IHVuZGVmaW5lZCB7XG4gICAgaWYgKCEoZWxlbWVudCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSkge1xuICAgICAgICByZXR1cm4gdW5kZWZpbmVkO1xuICAgIH1cblxuICAgIGNvbnN0IGtpbmQgPSBub25DbGllbnRSZWdpb25Gb3JFbGVtZW50KGVsZW1lbnQpO1xuICAgIGlmICgha2luZCB8fCAhaXNWaXNpYmxlKGVsZW1lbnQpKSB7XG4gICAgICAgIHJldHVybiB1bmRlZmluZWQ7XG4gICAgfVxuXG4gICAgY29uc3QgcmVjdCA9IGVsZW1lbnQuZ2V0Qm91bmRpbmdDbGllbnRSZWN0KCk7XG4gICAgaWYgKHJlY3Qud2lkdGggPD0gMCB8fCByZWN0LmhlaWdodCA8PSAwKSB7XG4gICAgICAgIHJldHVybiB1bmRlZmluZWQ7XG4gICAgfVxuXG4gICAgLy8gTmF0aXZlIGhpdCB0ZXN0aW5nIHJ1bnMgaW4gcGh5c2ljYWwgcGl4ZWxzLCB3aGlsZSBET00gZ2VvbWV0cnkgaXMgaW4gQ1NTIHBpeGVscy5cbiAgICBjb25zdCBzY2FsZSA9IHdpbmRvdy5kZXZpY2VQaXhlbFJhdGlvIHx8IDE7XG4gICAgY29uc3QgbGVmdCA9IE1hdGguZmxvb3IocmVjdC5sZWZ0ICogc2NhbGUpO1xuICAgIGNvbnN0IHRvcCA9IE1hdGguZmxvb3IocmVjdC50b3AgKiBzY2FsZSk7XG4gICAgY29uc3QgcmlnaHQgPSBNYXRoLmNlaWwocmVjdC5yaWdodCAqIHNjYWxlKTtcbiAgICBjb25zdCBib3R0b20gPSBNYXRoLmNlaWwocmVjdC5ib3R0b20gKiBzY2FsZSk7XG5cbiAgICBpZiAocmlnaHQgPD0gbGVmdCB8fCBib3R0b20gPD0gdG9wKSB7XG4gICAgICAgIHJldHVybiB1bmRlZmluZWQ7XG4gICAgfVxuXG4gICAgcmV0dXJuIHsga2luZCwgbGVmdCwgdG9wLCByaWdodCwgYm90dG9tIH07XG59XG5cbmZ1bmN0aW9uIHJlZ2lvbkVsZW1lbnRzKCk6IEVsZW1lbnRbXSB7XG4gICAgY29uc3QgZWxlbWVudHM6IEVsZW1lbnRbXSA9IFtdO1xuXG4gICAgaWYgKGRvY3VtZW50LmRvY3VtZW50RWxlbWVudCkge1xuICAgICAgICBlbGVtZW50cy5wdXNoKGRvY3VtZW50LmRvY3VtZW50RWxlbWVudCk7XG4gICAgfVxuICAgIGlmIChkb2N1bWVudC5ib2R5KSB7XG4gICAgICAgIGVsZW1lbnRzLnB1c2goZG9jdW1lbnQuYm9keSk7XG4gICAgICAgIC8vIEFwcGVuZCB2aWEgYSBsb29wOiBzcHJlYWRpbmcgYSBodWdlIE5vZGVMaXN0IGludG8gcHVzaCgpIG92ZXJmbG93c1xuICAgICAgICAvLyB0aGUgZW5naW5lJ3MgYXJndW1lbnQgbGltaXQgb24gdmVyeSBsYXJnZSBkb2N1bWVudHMuXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoXCIqXCIpKSB7XG4gICAgICAgICAgICBlbGVtZW50cy5wdXNoKGVsZW1lbnQpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIGVsZW1lbnRzO1xufVxuXG5mdW5jdGlvbiBvYnNlcnZlUmVnaW9uRWxlbWVudHMoZWxlbWVudHM6IEVsZW1lbnRbXSk6IHZvaWQge1xuICAgIGlmICh0eXBlb2YgUmVzaXplT2JzZXJ2ZXIgPT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFRyYWNrIHNpemUgY2hhbmdlcyBvbmx5IGZvciBhY3RpdmUgcmVnaW9uIGVsZW1lbnRzLiBET00gc3RydWN0dXJlIGFuZCBzdHlsZVxuICAgIC8vIGNoYW5nZXMgYXJlIGNvdmVyZWQgYnkgTXV0YXRpb25PYnNlcnZlciBpbiBzdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKCkuXG4gICAgcmVzaXplT2JzZXJ2ZXIgPz89IG5ldyBSZXNpemVPYnNlcnZlcihzY2hlZHVsZVVwZGF0ZSk7XG4gICAgY29uc3QgbmV4dEVsZW1lbnRzID0gbmV3IFNldChlbGVtZW50cyk7XG5cbiAgICBmb3IgKGNvbnN0IGVsZW1lbnQgb2Ygb2JzZXJ2ZWRFbGVtZW50cykge1xuICAgICAgICBpZiAoIW5leHRFbGVtZW50cy5oYXMoZWxlbWVudCkpIHtcbiAgICAgICAgICAgIHJlc2l6ZU9ic2VydmVyLnVub2JzZXJ2ZShlbGVtZW50KTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBuZXh0RWxlbWVudHMpIHtcbiAgICAgICAgaWYgKCFvYnNlcnZlZEVsZW1lbnRzLmhhcyhlbGVtZW50KSkge1xuICAgICAgICAgICAgcmVzaXplT2JzZXJ2ZXIub2JzZXJ2ZShlbGVtZW50KTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIG9ic2VydmVkRWxlbWVudHMgPSBuZXh0RWxlbWVudHM7XG59XG5cbmZ1bmN0aW9uIHVwZGF0ZU5vbkNsaWVudFJlZ2lvbnMoKTogdm9pZCB7XG4gICAgdXBkYXRlUGVuZGluZyA9IGZhbHNlO1xuXG4gICAgY29uc3QgZWxlbWVudHMgPSByZWdpb25FbGVtZW50cygpO1xuICAgIGNvbnN0IHJlZ2lvbnM6IE5vbkNsaWVudFJlZ2lvbltdID0gW107XG4gICAgY29uc3QgYWN0aXZlRWxlbWVudHM6IEVsZW1lbnRbXSA9IFtdO1xuXG4gICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGVsZW1lbnRzKSB7XG4gICAgICAgIGNvbnN0IHJlZ2lvbiA9IGVsZW1lbnRSZWdpb24oZWxlbWVudCk7XG4gICAgICAgIGlmIChyZWdpb24pIHtcbiAgICAgICAgICAgIHJlZ2lvbnMucHVzaChyZWdpb24pO1xuICAgICAgICAgICAgYWN0aXZlRWxlbWVudHMucHVzaChlbGVtZW50KTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIG9ic2VydmVSZWdpb25FbGVtZW50cyhhY3RpdmVFbGVtZW50cyk7XG5cbiAgICBjb25zdCBwYXlsb2FkID0gSlNPTi5zdHJpbmdpZnkoeyB2ZXJzaW9uOiAxLCByZWdpb25zIH0pO1xuICAgIGlmIChwYXlsb2FkID09PSBsYXN0UGF5bG9hZCkge1xuICAgICAgICAvLyBBdm9pZCBzZW5kaW5nIGR1cGxpY2F0ZSBuYXRpdmUgbWVzc2FnZXMgZHVyaW5nIHJlc2l6ZSBvciBzdHlsZSBjaHVybi5cbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIGxhc3RQYXlsb2FkID0gcGF5bG9hZDtcbiAgICBpbnZva2UoXCJ3YWlsczpub24tY2xpZW50LXJlZ2lvbjpcIiArIHBheWxvYWQpO1xufVxuXG5mdW5jdGlvbiBzY2hlZHVsZVVwZGF0ZSgpOiB2b2lkIHtcbiAgICBpZiAodXBkYXRlUGVuZGluZykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQmF0Y2ggcmVnaW9uIHVwZGF0ZXMgdG8gYW5pbWF0aW9uIGZyYW1lcyBzbyBsYXlvdXQgaXMgbWVhc3VyZWQgb25jZSBwZXIgZnJhbWUuXG4gICAgdXBkYXRlUGVuZGluZyA9IHRydWU7XG4gICAgd2luZG93LnJlcXVlc3RBbmltYXRpb25GcmFtZSh1cGRhdGVOb25DbGllbnRSZWdpb25zKTtcbn1cblxuZnVuY3Rpb24gc3RhcnROb25DbGllbnRSZWdpb25UcmFja2luZygpOiB2b2lkIHtcbiAgICBpZiAodHJhY2tpbmdTdGFydGVkKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICB0cmFja2luZ1N0YXJ0ZWQgPSB0cnVlO1xuICAgIC8vIFNlbmQgYW4gaW5pdGlhbCBlbXB0eSBvciBwb3B1bGF0ZWQgcmVnaW9uIGxpc3Qgb25jZSB0aGUgRE9NIGlzIHJlYWR5LlxuICAgIHNjaGVkdWxlVXBkYXRlKCk7XG5cbiAgICBjb25zdCBtdXRhdGlvbk9ic2VydmVyID0gbmV3IE11dGF0aW9uT2JzZXJ2ZXIoc2NoZWR1bGVVcGRhdGUpO1xuICAgIG11dGF0aW9uT2JzZXJ2ZXIub2JzZXJ2ZShkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQsIHtcbiAgICAgICAgYXR0cmlidXRlczogdHJ1ZSxcbiAgICAgICAgY2hpbGRMaXN0OiB0cnVlLFxuICAgICAgICBzdWJ0cmVlOiB0cnVlLFxuICAgIH0pO1xuXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoXCJyZXNpemVcIiwgc2NoZWR1bGVVcGRhdGUpO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKFwic2Nyb2xsXCIsIHNjaGVkdWxlVXBkYXRlLCB0cnVlKTtcbiAgICB3aW5kb3cudmlzdWFsVmlld3BvcnQ/LmFkZEV2ZW50TGlzdGVuZXIoXCJyZXNpemVcIiwgc2NoZWR1bGVVcGRhdGUpO1xuICAgIHdpbmRvdy52aXN1YWxWaWV3cG9ydD8uYWRkRXZlbnRMaXN0ZW5lcihcInNjcm9sbFwiLCBzY2hlZHVsZVVwZGF0ZSk7XG59XG5cbmZ1bmN0aW9uIHRyeVN0YXJ0Tm9uQ2xpZW50UmVnaW9uVHJhY2tpbmcoKTogYm9vbGVhbiB7XG4gICAgY29uc3Qgb3MgPSB3aW5kb3cuX3dhaWxzLmVudmlyb25tZW50Py5PUztcbiAgICBpZiAob3MgPT09IHVuZGVmaW5lZCkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgY29uc3QgZW5hYmxlZCA9IHdpbmRvdy5fd2FpbHMuZmxhZ3M/Lm5vbkNsaWVudFJlZ2lvblRyYWNraW5nO1xuICAgIGlmIChvcyA9PT0gXCJ3aW5kb3dzXCIpIHtcbiAgICAgICAgaWYgKGVuYWJsZWQgPT09IHRydWUpIHtcbiAgICAgICAgICAgIHdoZW5SZWFkeShzdGFydE5vbkNsaWVudFJlZ2lvblRyYWNraW5nKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9XG5cbiAgICByZXR1cm4gdHJ1ZTtcbn1cblxuaWYgKGhhc0RPTSAmJiAhdHJ5U3RhcnROb25DbGllbnRSZWdpb25UcmFja2luZygpKSB7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIocnVudGltZUNvbmZpZ1JlYWR5RXZlbnQsIHRyeVN0YXJ0Tm9uQ2xpZW50UmVnaW9uVHJhY2tpbmcsIHsgb25jZTogdHJ1ZSB9KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5BcHBsaWNhdGlvbik7XG5cbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xuY29uc3QgU2hvd01ldGhvZCA9IDE7XG5jb25zdCBRdWl0TWV0aG9kID0gMjtcblxuLyoqXG4gKiBIaWRlcyBhIGNlcnRhaW4gbWV0aG9kIGJ5IGNhbGxpbmcgdGhlIEhpZGVNZXRob2QgZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKEhpZGVNZXRob2QpO1xufVxuXG4vKipcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xufVxuXG4vKipcbiAqIENhbGxzIHRoZSBRdWl0TWV0aG9kIHRvIHRlcm1pbmF0ZSB0aGUgcHJvZ3JhbS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IENhbmNlbGxhYmxlUHJvbWlzZSwgdHlwZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzIH0gZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XG5cbi8vIFNldHVwXG5pbXBvcnQgeyBoYXNET00gfSBmcm9tIFwiLi9lbnZpcm9ubWVudC5qc1wiO1xuXG5pZiAoaGFzRE9NKSB7XG4gICAgd2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XG59XG5cbnR5cGUgUHJvbWlzZVJlc29sdmVycyA9IE9taXQ8Q2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxhbnk+LCBcInByb21pc2VcIiB8IFwib25jYW5jZWxsZWRcIj5cblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2FsbCk7XG5jb25zdCBjYW5jZWxDYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYW5jZWxDYWxsKTtcbmNvbnN0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwPHN0cmluZywgUHJvbWlzZVJlc29sdmVycz4oKTtcblxuY29uc3QgQ2FsbEJpbmRpbmcgPSAwO1xuY29uc3QgQ2FuY2VsTWV0aG9kID0gMFxuXG4vKipcbiAqIEhvbGRzIGFsbCByZXF1aXJlZCBpbmZvcm1hdGlvbiBmb3IgYSBiaW5kaW5nIGNhbGwuXG4gKiBNYXkgcHJvdmlkZSBlaXRoZXIgYSBtZXRob2QgSUQgb3IgYSBtZXRob2QgbmFtZSwgYnV0IG5vdCBib3RoLlxuICovXG5leHBvcnQgdHlwZSBDYWxsT3B0aW9ucyA9IHtcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZElEOiBudW1iZXI7XG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kTmFtZT86IG5ldmVyO1xuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cbiAgICBhcmdzOiBhbnlbXTtcbn0gfCB7XG4gICAgLyoqIFRoZSBudW1lcmljIElEIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2RJRD86IG5ldmVyO1xuICAgIC8qKiBUaGUgZnVsbHkgcXVhbGlmaWVkIG5hbWUgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZE5hbWU6IHN0cmluZztcbiAgICAvKiogQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuICovXG4gICAgYXJnczogYW55W107XG59O1xuXG4vLyBydW50aW1lLmpzIG5lZWRzIHRvIHVzZSBSdW50aW1lRXJyb3IgaW50ZXJuYWxseSB0byBwcm9wZXJseSBwYXJzZSBhbmQgcmV0dXJuXG4vLyBlcnJvcnMgZm9yIGJpbmRpbmcgY2FsbHMsIHNvIGl0IGhhZCB0byBtb3ZlIHRoZXJlLiBFeHBvcnRpbmcgaGVyZSBhZ2FpbiB0b1xuLy8ga2VlcCBmcm9tIGJyZWFraW5nIHRoZSBwdWJsaWMgQ2FsbCBpbnRlcmZhY2UuXG5leHBvcnQgeyBSdW50aW1lRXJyb3IgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbi8qKlxuICogR2VuZXJhdGVzIGEgdW5pcXVlIElEIHVzaW5nIHRoZSBuYW5vaWQgbGlicmFyeS5cbiAqXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cbiAqL1xuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpOiBzdHJpbmcge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIENhbGwgYSBib3VuZCBtZXRob2QgYWNjb3JkaW5nIHRvIHRoZSBnaXZlbiBjYWxsIG9wdGlvbnMuXG4gKlxuICogSW4gY2FzZSBvZiBmYWlsdXJlLCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGFuIGV4Y2VwdGlvblxuICogYW1vbmcgUmVmZXJlbmNlRXJyb3IgKHVua25vd24gbWV0aG9kKSwgVHlwZUVycm9yICh3cm9uZyBhcmd1bWVudCBjb3VudCBvciB0eXBlKSxcbiAqIHtAbGluayBSdW50aW1lRXJyb3J9IChtZXRob2QgcmV0dXJuZWQgYW4gZXJyb3IpLCBvciBvdGhlciAobmV0d29yayBvciBpbnRlcm5hbCBlcnJvcnMpLlxuICogVGhlIGV4Y2VwdGlvbiBtaWdodCBoYXZlIGEgXCJjYXVzZVwiIGZpZWxkIHdpdGggdGhlIHZhbHVlIHJldHVybmVkXG4gKiBieSB0aGUgYXBwbGljYXRpb24tIG9yIHNlcnZpY2UtbGV2ZWwgZXJyb3IgbWFyc2hhbGluZyBmdW5jdGlvbnMuXG4gKlxuICogQHBhcmFtIG9wdGlvbnMgLSBBIG1ldGhvZCBjYWxsIGRlc2NyaXB0b3IuXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQ2FsbChvcHRpb25zOiBDYWxsT3B0aW9ucyk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcbiAgICBjb25zdCBpZCA9IGdlbmVyYXRlSUQoKTtcblxuICAgIGNvbnN0IHJlc3VsdCA9IENhbmNlbGxhYmxlUHJvbWlzZS53aXRoUmVzb2x2ZXJzPGFueT4oKTtcbiAgICBjYWxsUmVzcG9uc2VzLnNldChpZCwgeyByZXNvbHZlOiByZXN1bHQucmVzb2x2ZSwgcmVqZWN0OiByZXN1bHQucmVqZWN0IH0pO1xuXG4gICAgY29uc3QgcmVxdWVzdCA9IGNhbGwoQ2FsbEJpbmRpbmcsIE9iamVjdC5hc3NpZ24oeyBcImNhbGwtaWRcIjogaWQgfSwgb3B0aW9ucykpO1xuICAgIGxldCBydW5uaW5nID0gdHJ1ZTtcblxuICAgIHJlcXVlc3QudGhlbigocmVzKSA9PiB7XG4gICAgICAgIHJ1bm5pbmcgPSBmYWxzZTtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXN1bHQucmVzb2x2ZShyZXMpO1xuICAgIH0sIChlcnIpID0+IHtcbiAgICAgICAgcnVubmluZyA9IGZhbHNlO1xuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIHJlc3VsdC5yZWplY3QoZXJyKTtcbiAgICB9KTtcblxuICAgIGNvbnN0IGNhbmNlbCA9ICgpID0+IHtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXR1cm4gY2FuY2VsQ2FsbChDYW5jZWxNZXRob2QsIHtcImNhbGwtaWRcIjogaWR9KS5jYXRjaCgoZXJyKSA9PiB7XG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igd2hpbGUgcmVxdWVzdGluZyBiaW5kaW5nIGNhbGwgY2FuY2VsbGF0aW9uOlwiLCBlcnIpO1xuICAgICAgICB9KTtcbiAgICB9O1xuXG4gICAgcmVzdWx0Lm9uY2FuY2VsbGVkID0gKCkgPT4ge1xuICAgICAgICBpZiAocnVubmluZykge1xuICAgICAgICAgICAgcmV0dXJuIGNhbmNlbCgpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcmV0dXJuIHJlcXVlc3QudGhlbihjYW5jZWwpO1xuICAgICAgICB9XG4gICAgfTtcblxuICAgIHJldHVybiByZXN1bHQucHJvbWlzZTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIGJvdW5kIG1ldGhvZCBieSBuYW1lIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKiBTZWUge0BsaW5rIENhbGx9IGZvciBkZXRhaWxzLlxuICpcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEByZXR1cm5zIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG1ldGhvZE5hbWU6IHN0cmluZywgLi4uYXJnczogYW55W10pOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2ROYW1lLCBhcmdzIH0pO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIGJ5IGl0cyBudW1lcmljIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXG4gKiBTZWUge0BsaW5rIENhbGx9IGZvciBkZXRhaWxzLlxuICpcbiAqIEBwYXJhbSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxuICogQHJldHVybiBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQ6IG51bWJlciwgLi4uYXJnczogYW55W10pOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIENhbGwoeyBtZXRob2RJRCwgYXJncyB9KTtcbn1cbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9pbnNwZWN0LWpzL2lzLWNhbGxhYmxlXG5cbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxuLy9cbi8vIENvcHlyaWdodCAoYykgMjAxNSBKb3JkYW4gSGFyYmFuZFxuLy9cbi8vIFBlcm1pc3Npb24gaXMgaGVyZWJ5IGdyYW50ZWQsIGZyZWUgb2YgY2hhcmdlLCB0byBhbnkgcGVyc29uIG9idGFpbmluZyBhIGNvcHlcbi8vIG9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWxcbi8vIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHNcbi8vIHRvIHVzZSwgY29weSwgbW9kaWZ5LCBtZXJnZSwgcHVibGlzaCwgZGlzdHJpYnV0ZSwgc3VibGljZW5zZSwgYW5kL29yIHNlbGxcbi8vIGNvcGllcyBvZiB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpc1xuLy8gZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcbi8vXG4vLyBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXG4vL1xuLy8gVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksXG4vLyBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBTkQgTk9OSU5GUklOR0VNRU5ULiBJTiBOTyBFVkVOVCBTSEFMTCBUSEVcbi8vIEFVVEhPUlMgT1IgQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVJcbi8vIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sXG4vLyBPVVQgT0YgT1IgSU4gQ09OTkVDVElPTiBXSVRIIFRIRSBTT0ZUV0FSRSBPUiBUSEUgVVNFIE9SIE9USEVSIERFQUxJTkdTIElOIFRIRVxuLy8gU09GVFdBUkUuXG5cbnZhciBmblRvU3RyID0gRnVuY3Rpb24ucHJvdG90eXBlLnRvU3RyaW5nO1xudmFyIHJlZmxlY3RBcHBseTogdHlwZW9mIFJlZmxlY3QuYXBwbHkgfCBmYWxzZSB8IG51bGwgPSB0eXBlb2YgUmVmbGVjdCA9PT0gJ29iamVjdCcgJiYgUmVmbGVjdCAhPT0gbnVsbCAmJiBSZWZsZWN0LmFwcGx5O1xudmFyIGJhZEFycmF5TGlrZTogYW55O1xudmFyIGlzQ2FsbGFibGVNYXJrZXI6IGFueTtcbmlmICh0eXBlb2YgcmVmbGVjdEFwcGx5ID09PSAnZnVuY3Rpb24nICYmIHR5cGVvZiBPYmplY3QuZGVmaW5lUHJvcGVydHkgPT09ICdmdW5jdGlvbicpIHtcbiAgICB0cnkge1xuICAgICAgICBiYWRBcnJheUxpa2UgPSBPYmplY3QuZGVmaW5lUHJvcGVydHkoe30sICdsZW5ndGgnLCB7XG4gICAgICAgICAgICBnZXQ6IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgICAgICB0aHJvdyBpc0NhbGxhYmxlTWFya2VyO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgaXNDYWxsYWJsZU1hcmtlciA9IHt9O1xuICAgICAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmUgbm8tdGhyb3ctbGl0ZXJhbFxuICAgICAgICByZWZsZWN0QXBwbHkoZnVuY3Rpb24gKCkgeyB0aHJvdyA0MjsgfSwgbnVsbCwgYmFkQXJyYXlMaWtlKTtcbiAgICB9IGNhdGNoIChfKSB7XG4gICAgICAgIGlmIChfICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7XG4gICAgICAgICAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xuICAgICAgICB9XG4gICAgfVxufSBlbHNlIHtcbiAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xufVxuXG52YXIgY29uc3RydWN0b3JSZWdleCA9IC9eXFxzKmNsYXNzXFxiLztcbnZhciBpc0VTNkNsYXNzRm4gPSBmdW5jdGlvbiBpc0VTNkNsYXNzRnVuY3Rpb24odmFsdWU6IGFueSk6IGJvb2xlYW4ge1xuICAgIHRyeSB7XG4gICAgICAgIHZhciBmblN0ciA9IGZuVG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgIHJldHVybiBjb25zdHJ1Y3RvclJlZ2V4LnRlc3QoZm5TdHIpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlOyAvLyBub3QgYSBmdW5jdGlvblxuICAgIH1cbn07XG5cbnZhciB0cnlGdW5jdGlvbk9iamVjdCA9IGZ1bmN0aW9uIHRyeUZ1bmN0aW9uVG9TdHIodmFsdWU6IGFueSk6IGJvb2xlYW4ge1xuICAgIHRyeSB7XG4gICAgICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgICAgICBmblRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICByZXR1cm4gdHJ1ZTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG59O1xudmFyIHRvU3RyID0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZztcbnZhciBvYmplY3RDbGFzcyA9ICdbb2JqZWN0IE9iamVjdF0nO1xudmFyIGZuQ2xhc3MgPSAnW29iamVjdCBGdW5jdGlvbl0nO1xudmFyIGdlbkNsYXNzID0gJ1tvYmplY3QgR2VuZXJhdG9yRnVuY3Rpb25dJztcbnZhciBkZGFDbGFzcyA9ICdbb2JqZWN0IEhUTUxBbGxDb2xsZWN0aW9uXSc7IC8vIElFIDExXG52YXIgZGRhQ2xhc3MyID0gJ1tvYmplY3QgSFRNTCBkb2N1bWVudC5hbGwgY2xhc3NdJztcbnZhciBkZGFDbGFzczMgPSAnW29iamVjdCBIVE1MQ29sbGVjdGlvbl0nOyAvLyBJRSA5LTEwXG52YXIgaGFzVG9TdHJpbmdUYWcgPSB0eXBlb2YgU3ltYm9sID09PSAnZnVuY3Rpb24nICYmICEhU3ltYm9sLnRvU3RyaW5nVGFnOyAvLyBiZXR0ZXI6IHVzZSBgaGFzLXRvc3RyaW5ndGFnYFxuXG52YXIgaXNJRTY4ID0gISgwIGluIFssXSk7IC8vIGVzbGludC1kaXNhYmxlLWxpbmUgbm8tc3BhcnNlLWFycmF5cywgY29tbWEtc3BhY2luZ1xuXG52YXIgaXNEREE6ICh2YWx1ZTogYW55KSA9PiBib29sZWFuID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCgpIHsgcmV0dXJuIGZhbHNlOyB9O1xuaWYgKHR5cGVvZiBkb2N1bWVudCA9PT0gJ29iamVjdCcpIHtcbiAgICAvLyBGaXJlZm94IDMgY2Fub25pY2FsaXplcyBEREEgdG8gdW5kZWZpbmVkIHdoZW4gaXQncyBub3QgYWNjZXNzZWQgZGlyZWN0bHlcbiAgICB2YXIgYWxsID0gZG9jdW1lbnQuYWxsO1xuICAgIGlmICh0b1N0ci5jYWxsKGFsbCkgPT09IHRvU3RyLmNhbGwoZG9jdW1lbnQuYWxsKSkge1xuICAgICAgICBpc0REQSA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwodmFsdWUpIHtcbiAgICAgICAgICAgIC8qIGdsb2JhbHMgZG9jdW1lbnQ6IGZhbHNlICovXG4gICAgICAgICAgICAvLyBpbiBJRSA2LTgsIHR5cGVvZiBkb2N1bWVudC5hbGwgaXMgXCJvYmplY3RcIiBhbmQgaXQncyB0cnV0aHlcbiAgICAgICAgICAgIGlmICgoaXNJRTY4IHx8ICF2YWx1ZSkgJiYgKHR5cGVvZiB2YWx1ZSA9PT0gJ3VuZGVmaW5lZCcgfHwgdHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JykpIHtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICB2YXIgc3RyID0gdG9TdHIuY2FsbCh2YWx1ZSk7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybiAoXG4gICAgICAgICAgICAgICAgICAgICAgICBzdHIgPT09IGRkYUNsYXNzXG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IGRkYUNsYXNzMlxuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczMgLy8gb3BlcmEgMTIuMTZcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gb2JqZWN0Q2xhc3MgLy8gSUUgNi04XG4gICAgICAgICAgICAgICAgICAgICkgJiYgdmFsdWUoJycpID09IG51bGw7IC8vIGVzbGludC1kaXNhYmxlLWxpbmUgZXFlcWVxXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZSkgeyAvKiovIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgfTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGlzQ2FsbGFibGVSZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSAge1xuICAgIGlmIChpc0REQSh2YWx1ZSkpIHsgcmV0dXJuIHRydWU7IH1cbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB0cnkge1xuICAgICAgICAocmVmbGVjdEFwcGx5IGFzIGFueSkodmFsdWUsIG51bGwsIGJhZEFycmF5TGlrZSk7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICBpZiAoZSAhPT0gaXNDYWxsYWJsZU1hcmtlcikgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB9XG4gICAgcmV0dXJuICFpc0VTNkNsYXNzRm4odmFsdWUpICYmIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcbn1cblxuZnVuY3Rpb24gaXNDYWxsYWJsZU5vUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkge1xuICAgIGlmIChpc0REQSh2YWx1ZSkpIHsgcmV0dXJuIHRydWU7IH1cbiAgICBpZiAoIXZhbHVlKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAoaGFzVG9TdHJpbmdUYWcpIHsgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTsgfVxuICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHZhciBzdHJDbGFzcyA9IHRvU3RyLmNhbGwodmFsdWUpO1xuICAgIGlmIChzdHJDbGFzcyAhPT0gZm5DbGFzcyAmJiBzdHJDbGFzcyAhPT0gZ2VuQ2xhc3MgJiYgISgvXlxcW29iamVjdCBIVE1MLykudGVzdChzdHJDbGFzcykpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgcmV0dXJuIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcbn07XG5cbmV4cG9ydCBkZWZhdWx0IHJlZmxlY3RBcHBseSA/IGlzQ2FsbGFibGVSZWZBcHBseSA6IGlzQ2FsbGFibGVOb1JlZkFwcGx5O1xuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgaXNDYWxsYWJsZSBmcm9tIFwiLi9jYWxsYWJsZS5qc1wiO1xuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgdXNlZCBhcyByZWplY3Rpb24gcmVhc29uXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gaXMgY2FuY2VsbGVkIHN1Y2Nlc3NmdWxseS5cbiAqXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbEVycm9yXCJgLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSBjYXVzZSBwYXNzZWQgdG8gdGhlIGNhbmNlbCBtZXRob2QsIGlmIGFueS5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbEVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbEVycm9yYCBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsRXJyb3JcIjtcbiAgICB9XG59XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSByZXBvcnRlZCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gcmVqZWN0cyBhZnRlciBiZWluZyBjYW5jZWxsZWQsXG4gKiBvciB3aGVuIHRoZSBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHRocm93cyBvciByZWplY3RzLlxuICpcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcImAuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBjYXVzZX0gcHJvcGVydHkgaXMgdGhlIHJlYXNvbiB0aGUgcHJvbWlzZSByZWplY3RlZCB3aXRoLlxuICpcbiAqIEJlY2F1c2UgdGhlIG9yaWdpbmFsIHByb21pc2Ugd2FzIGNhbmNlbGxlZCxcbiAqIGEgd3JhcHBlciBwcm9taXNlIHdpbGwgYmUgcGFzc2VkIHRvIHRoZSB1bmhhbmRsZWQgcmVqZWN0aW9uIGxpc3RlbmVyIGluc3RlYWQuXG4gKiBUaGUge0BsaW5rIHByb21pc2V9IHByb3BlcnR5IGhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBvcmlnaW5hbCBwcm9taXNlLlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IgZXh0ZW5kcyBFcnJvciB7XG4gICAgLyoqXG4gICAgICogSG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIHByb21pc2UgdGhhdCB3YXMgY2FuY2VsbGVkIGFuZCB0aGVuIHJlamVjdGVkLlxuICAgICAqL1xuICAgIHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPjtcblxuICAgIC8qKlxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbGxlZFJlamVjdGlvbkVycm9yYCBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gcHJvbWlzZSAtIFRoZSBwcm9taXNlIHRoYXQgY2F1c2VkIHRoZSBlcnJvciBvcmlnaW5hbGx5LlxuICAgICAqIEBwYXJhbSByZWFzb24gLSBUaGUgcmVqZWN0aW9uIHJlYXNvbi5cbiAgICAgKiBAcGFyYW0gaW5mbyAtIEFuIG9wdGlvbmFsIGluZm9ybWF0aXZlIG1lc3NhZ2Ugc3BlY2lmeWluZyB0aGUgY2lyY3Vtc3RhbmNlcyBpbiB3aGljaCB0aGUgZXJyb3Igd2FzIHRocm93bi5cbiAgICAgKiAgICAgICAgICAgICAgIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgYFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cImAuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCByZWFzb24/OiBhbnksIGluZm8/OiBzdHJpbmcpIHtcbiAgICAgICAgc3VwZXIoKGluZm8gPz8gXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiKSArIFwiIFJlYXNvbjogXCIgKyBlcnJvck1lc3NhZ2UocmVhc29uKSwgeyBjYXVzZTogcmVhc29uIH0pO1xuICAgICAgICB0aGlzLnByb21pc2UgPSBwcm9taXNlO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCI7XG4gICAgfVxufVxuXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+ID0gKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciA9IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciA9IChjYXVzZT86IGFueSkgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiA9IChyZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiwgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcikgPT4gdm9pZDtcblxuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwpOiBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxIHwgVFJlc3VsdDI+O1xuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbn1cblxuLyoqXG4gKiBXcmFwcyBhIGNhbmNlbGxhYmxlIHByb21pc2UgYWxvbmcgd2l0aCBpdHMgcmVzb2x1dGlvbiBtZXRob2RzLlxuICogVGhlIGBvbmNhbmNlbGxlZGAgZmllbGQgd2lsbCBiZSBudWxsIGluaXRpYWxseSBidXQgbWF5IGJlIHNldCB0byBwcm92aWRlIGEgY3VzdG9tIGNhbmNlbGxhdGlvbiBmdW5jdGlvbi5cbiAqL1xuZXhwb3J0IGludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XG4gICAgcmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD47XG4gICAgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcjtcbiAgICBvbmNhbmNlbGxlZDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyIHwgbnVsbDtcbn1cblxuaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlIHtcbiAgICByZWFkb25seSByb290OiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZTtcbiAgICByZXNvbHZpbmc6IGJvb2xlYW47XG4gICAgc2V0dGxlZDogYm9vbGVhbjtcbiAgICByZWFzb24/OiBDYW5jZWxFcnJvcjtcbn1cblxuLy8gUHJpdmF0ZSBmaWVsZCBuYW1lcy5cbmNvbnN0IGJhcnJpZXJTeW0gPSBTeW1ib2woXCJiYXJyaWVyXCIpO1xuY29uc3QgY2FuY2VsSW1wbFN5bSA9IFN5bWJvbChcImNhbmNlbEltcGxcIik7XG5jb25zdCBzcGVjaWVzOiB0eXBlb2YgU3ltYm9sLnNwZWNpZXMgPSBTeW1ib2wuc3BlY2llcyA/PyBTeW1ib2woXCJzcGVjaWVzUG9seWZpbGxcIik7XG5cbi8qKlxuICogQSBwcm9taXNlIHdpdGggYW4gYXR0YWNoZWQgbWV0aG9kIGZvciBjYW5jZWxsaW5nIGxvbmctcnVubmluZyBvcGVyYXRpb25zIChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWx9KS5cbiAqIENhbmNlbGxhdGlvbiBjYW4gb3B0aW9uYWxseSBiZSBib3VuZCB0byBhbiB7QGxpbmsgQWJvcnRTaWduYWx9XG4gKiBmb3IgYmV0dGVyIGNvbXBvc2FiaWxpdHkgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbE9ufSkuXG4gKlxuICogQ2FuY2VsbGluZyBhIHBlbmRpbmcgcHJvbWlzZSB3aWxsIHJlc3VsdCBpbiBhbiBpbW1lZGlhdGUgcmVqZWN0aW9uXG4gKiB3aXRoIGFuIGluc3RhbmNlIG9mIHtAbGluayBDYW5jZWxFcnJvcn0gYXMgcmVhc29uLFxuICogYnV0IHdob2V2ZXIgc3RhcnRlZCB0aGUgcHJvbWlzZSB3aWxsIGJlIHJlc3BvbnNpYmxlXG4gKiBmb3IgYWN0dWFsbHkgYWJvcnRpbmcgdGhlIHVuZGVybHlpbmcgb3BlcmF0aW9uLlxuICogVG8gdGhpcyBwdXJwb3NlLCB0aGUgY29uc3RydWN0b3IgYW5kIGFsbCBjaGFpbmluZyBtZXRob2RzXG4gKiBhY2NlcHQgb3B0aW9uYWwgY2FuY2VsbGF0aW9uIGNhbGxiYWNrcy5cbiAqXG4gKiBJZiBhIGBDYW5jZWxsYWJsZVByb21pc2VgIHN0aWxsIHJlc29sdmVzIGFmdGVyIGhhdmluZyBiZWVuIGNhbmNlbGxlZCxcbiAqIHRoZSByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWQuIElmIGl0IHJlamVjdHMsIHRoZSByZWFzb25cbiAqIHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbixcbiAqIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGluc3RhbmNlLlxuICogVG8gZmFjaWxpdGF0ZSB0aGUgaGFuZGxpbmcgb2YgY2FuY2VsbGF0aW9uIHJlcXVlc3RzLFxuICogY2FuY2VsbGVkIGBDYW5jZWxsYWJsZVByb21pc2VgcyB3aWxsIF9ub3RfIHJlcG9ydCB1bmhhbmRsZWQgYENhbmNlbEVycm9yYHNcbiAqIHdob3NlIGBjYXVzZWAgZmllbGQgaXMgdGhlIHNhbWUgYXMgdGhlIG9uZSB3aXRoIHdoaWNoIHRoZSBjdXJyZW50IHByb21pc2Ugd2FzIGNhbmNlbGxlZC5cbiAqXG4gKiBBbGwgdXN1YWwgcHJvbWlzZSBtZXRob2RzIGFyZSBkZWZpbmVkIGFuZCByZXR1cm4gYSBgQ2FuY2VsbGFibGVQcm9taXNlYFxuICogd2hvc2UgY2FuY2VsIG1ldGhvZCB3aWxsIGNhbmNlbCB0aGUgcGFyZW50IG9wZXJhdGlvbiBhcyB3ZWxsLCBwcm9wYWdhdGluZyB0aGUgY2FuY2VsbGF0aW9uIHJlYXNvblxuICogdXB3YXJkcyB0aHJvdWdoIHByb21pc2UgY2hhaW5zLlxuICogQ29udmVyc2VseSwgY2FuY2VsbGluZyBhIHByb21pc2Ugd2lsbCBub3QgYXV0b21hdGljYWxseSBjYW5jZWwgZGVwZW5kZW50IHByb21pc2VzIGRvd25zdHJlYW06XG4gKiBgYGB0c1xuICogbGV0IHJvb3QgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICogbGV0IGNoaWxkMSA9IHJvb3QudGhlbigoKSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDIgPSBjaGlsZDEudGhlbigoKSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDMgPSByb290LmNhdGNoKCgpID0+IHsgLi4uIH0pO1xuICogY2hpbGQxLmNhbmNlbCgpOyAvLyBDYW5jZWxzIGNoaWxkMSBhbmQgcm9vdCwgYnV0IG5vdCBjaGlsZDIgb3IgY2hpbGQzXG4gKiBgYGBcbiAqIENhbmNlbGxpbmcgYSBwcm9taXNlIHRoYXQgaGFzIGFscmVhZHkgc2V0dGxlZCBpcyBzYWZlIGFuZCBoYXMgbm8gY29uc2VxdWVuY2UuXG4gKlxuICogVGhlIGBjYW5jZWxgIG1ldGhvZCByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IF9hbHdheXMgZnVsZmlsbHNfXG4gKiBhZnRlciB0aGUgd2hvbGUgY2hhaW4gaGFzIHByb2Nlc3NlZCB0aGUgY2FuY2VsIHJlcXVlc3RcbiAqIGFuZCBhbGwgYXR0YWNoZWQgY2FsbGJhY2tzIHVwIHRvIHRoYXQgbW9tZW50IGhhdmUgcnVuLlxuICpcbiAqIEFsbCBFUzIwMjQgcHJvbWlzZSBtZXRob2RzIChzdGF0aWMgYW5kIGluc3RhbmNlKSBhcmUgZGVmaW5lZCBvbiBDYW5jZWxsYWJsZVByb21pc2UsXG4gKiBidXQgYWN0dWFsIGF2YWlsYWJpbGl0eSBtYXkgdmFyeSB3aXRoIE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAqXG4gKiBJbiBsaW5lIHdpdGggdGhlIHByb3Bvc2FsIGF0IGh0dHBzOi8vZ2l0aHViLmNvbS90YzM5L3Byb3Bvc2FsLXJtLWJ1aWx0aW4tc3ViY2xhc3NpbmcsXG4gKiBgQ2FuY2VsbGFibGVQcm9taXNlYCBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLlxuICogRXh0ZW5kZXJzIHNob3VsZCB0YWtlIGNhcmUgdG8gcHJvdmlkZSB0aGVpciBvd24gbWV0aG9kIGltcGxlbWVudGF0aW9ucy5cbiAqIFRoaXMgbWlnaHQgYmUgcmVjb25zaWRlcmVkIGluIGNhc2UgdGhlIHByb3Bvc2FsIGlzIHJldGlyZWQuXG4gKlxuICogQ2FuY2VsbGFibGVQcm9taXNlIGlzIGEgd3JhcHBlciBhcm91bmQgdGhlIERPTSBQcm9taXNlIG9iamVjdFxuICogYW5kIGlzIGNvbXBsaWFudCB3aXRoIHRoZSBbUHJvbWlzZXMvQSsgc3BlY2lmaWNhdGlvbl0oaHR0cHM6Ly9wcm9taXNlc2FwbHVzLmNvbS8pXG4gKiAoaXQgcGFzc2VzIHRoZSBbY29tcGxpYW5jZSBzdWl0ZV0oaHR0cHM6Ly9naXRodWIuY29tL3Byb21pc2VzLWFwbHVzL3Byb21pc2VzLXRlc3RzKSlcbiAqIGlmIHNvIGlzIHRoZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uLlxuICovXG5leHBvcnQgY2xhc3MgQ2FuY2VsbGFibGVQcm9taXNlPFQ+IGV4dGVuZHMgUHJvbWlzZTxUPiBpbXBsZW1lbnRzIFByb21pc2VMaWtlPFQ+LCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+IHtcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cbiAgICAvKiogQGludGVybmFsICovXG4gICAgcHJpdmF0ZSBbYmFycmllclN5bV0hOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiB8IG51bGw7XG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIHByaXZhdGUgcmVhZG9ubHkgW2NhbmNlbEltcGxTeW1dITogKHJlYXNvbjogQ2FuY2VsRXJyb3IpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgYENhbmNlbGxhYmxlUHJvbWlzZWAuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZXhlY3V0b3IgLSBBIGNhbGxiYWNrIHVzZWQgdG8gaW5pdGlhbGl6ZSB0aGUgcHJvbWlzZS4gVGhpcyBjYWxsYmFjayBpcyBwYXNzZWQgdHdvIGFyZ3VtZW50czpcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhIGByZXNvbHZlYCBjYWxsYmFjayB1c2VkIHRvIHJlc29sdmUgdGhlIHByb21pc2Ugd2l0aCBhIHZhbHVlXG4gICAgICogICAgICAgICAgICAgICAgICAgb3IgdGhlIHJlc3VsdCBvZiBhbm90aGVyIHByb21pc2UgKHBvc3NpYmx5IGNhbmNlbGxhYmxlKSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYSBgcmVqZWN0YCBjYWxsYmFjayB1c2VkIHRvIHJlamVjdCB0aGUgcHJvbWlzZSB3aXRoIGEgcHJvdmlkZWQgcmVhc29uIG9yIGVycm9yLlxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIHRoZSB2YWx1ZSBwcm92aWRlZCB0byB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogICAgICAgICAgICAgICAgICAgKGl0IGhhcyBhIGB0aGVuYCBfYW5kXyBhIGBjYW5jZWxgIG1ldGhvZCksXG4gICAgICogICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZm9yd2FyZGVkIHRvIHRoYXQgb2JqZWN0IGFuZCB0aGUgb25jYW5jZWxsZWQgd2lsbCBub3QgYmUgaW52b2tlZCBhbnltb3JlLlxuICAgICAqICAgICAgICAgICAgICAgICAgIElmIGFueSBvbmUgb2YgdGhlIHR3byBjYWxsYmFja3MgaXMgY2FsbGVkIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkLFxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoZSBwcm92aWRlZCB2YWx1ZXMgd2lsbCBiZSBjYW5jZWxsZWQgYW5kIHJlc29sdmVkIGFzIHVzdWFsLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGJ1dCB0aGVpciByZXN1bHRzIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqICAgICAgICAgICAgICAgICAgIEhvd2V2ZXIsIGlmIHRoZSByZXNvbHV0aW9uIHByb2Nlc3MgdWx0aW1hdGVseSBlbmRzIHVwIGluIGEgcmVqZWN0aW9uXG4gICAgICogICAgICAgICAgICAgICAgICAgdGhhdCBpcyBub3QgZHVlIHRvIGNhbmNlbGxhdGlvbiwgdGhlIHJlamVjdGlvbiByZWFzb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICB3aWxsIGJlIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9XG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKiBAcGFyYW0gb25jYW5jZWxsZWQgLSBJdCBpcyB0aGUgY2FsbGVyJ3MgcmVzcG9uc2liaWxpdHkgdG8gZW5zdXJlIHRoYXQgYW55IG9wZXJhdGlvblxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHN0YXJ0ZWQgYnkgdGhlIGV4ZWN1dG9yIGlzIHByb3Blcmx5IGhhbHRlZCB1cG9uIGNhbmNlbGxhdGlvbi5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBUaGlzIG9wdGlvbmFsIGNhbGxiYWNrIGNhbiBiZSB1c2VkIHRvIHRoYXQgcHVycG9zZS5cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJdCB3aWxsIGJlIGNhbGxlZCBfc3luY2hyb25vdXNseV8gd2l0aCBhIGNhbmNlbGxhdGlvbiBjYXVzZVxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdoZW4gY2FuY2VsbGF0aW9uIGlzIHJlcXVlc3RlZCwgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSByZWplY3RlZFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9LCBidXQgX2JlZm9yZV9cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbnkge0BsaW5rIHRoZW59L3tAbGluayBjYXRjaH0ve0BsaW5rIGZpbmFsbHl9IGNhbGxiYWNrIHJ1bnMuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfVxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdpbGwgb25seSBmdWxmaWxsIGFmdGVyIHRoZSBmb3JtZXIgaGFzIHNldHRsZWQuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVW5oYW5kbGVkIGV4Y2VwdGlvbnMgb3IgcmVqZWN0aW9ucyBmcm9tIHRoZSBjYWxsYmFjayB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgdW5oYW5kbGVkIHJlamVjdGlvbnMuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBjYWxsZWQgYmVmb3JlIGNhbmNlbGxhdGlvbiB3aXRoIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBjYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhpcyBwcm9taXNlIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gdGhhdCBwcm9taXNlLFxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFuZCB0aGUgb3JpZ2luYWwgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihleGVjdXRvcjogQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4sIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKSB7XG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xuICAgICAgICBzdXBlcigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcblxuICAgICAgICBpZiAoKHRoaXMuY29uc3RydWN0b3IgYXMgYW55KVtzcGVjaWVzXSAhPT0gUHJvbWlzZSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZSBkb2VzIG5vdCBzdXBwb3J0IHRyYW5zcGFyZW50IHN1YmNsYXNzaW5nLiBQbGVhc2UgcmVmcmFpbiBmcm9tIG92ZXJyaWRpbmcgdGhlIFtTeW1ib2wuc3BlY2llc10gc3RhdGljIHByb3BlcnR5LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGxldCBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+ID0ge1xuICAgICAgICAgICAgcHJvbWlzZTogdGhpcyxcbiAgICAgICAgICAgIHJlc29sdmUsXG4gICAgICAgICAgICByZWplY3QsXG4gICAgICAgICAgICBnZXQgb25jYW5jZWxsZWQoKSB7IHJldHVybiBvbmNhbmNlbGxlZCA/PyBudWxsOyB9LFxuICAgICAgICAgICAgc2V0IG9uY2FuY2VsbGVkKGNiKSB7IG9uY2FuY2VsbGVkID0gY2IgPz8gdW5kZWZpbmVkOyB9XG4gICAgICAgIH07XG5cbiAgICAgICAgY29uc3Qgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xuICAgICAgICAgICAgZ2V0IHJvb3QoKSB7IHJldHVybiBzdGF0ZTsgfSxcbiAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXG4gICAgICAgICAgICBzZXR0bGVkOiBmYWxzZVxuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFNldHVwIGNhbmNlbGxhdGlvbiBzeXN0ZW0uXG4gICAgICAgIHZvaWQgT2JqZWN0LmRlZmluZVByb3BlcnRpZXModGhpcywge1xuICAgICAgICAgICAgW2JhcnJpZXJTeW1dOiB7XG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICBlbnVtZXJhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogdHJ1ZSxcbiAgICAgICAgICAgICAgICB2YWx1ZTogbnVsbFxuICAgICAgICAgICAgfSxcbiAgICAgICAgICAgIFtjYW5jZWxJbXBsU3ltXToge1xuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHZhbHVlOiBjYW5jZWxsZXJGb3IocHJvbWlzZSwgc3RhdGUpXG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuXG4gICAgICAgIC8vIFJ1biB0aGUgYWN0dWFsIGV4ZWN1dG9yLlxuICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIHN0YXRlKTtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIGV4ZWN1dG9yKHJlc29sdmVyRm9yKHByb21pc2UsIHN0YXRlKSwgcmVqZWN0b3IpO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHtcbiAgICAgICAgICAgICAgICBjb25zb2xlLmxvZyhcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gQ2FuY2VsbGFibGVQcm9taXNlIGV4ZWN1dG9yLlwiLCBlcnIpO1xuICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2FuY2VscyBpbW1lZGlhdGVseSB0aGUgZXhlY3V0aW9uIG9mIHRoZSBvcGVyYXRpb24gYXNzb2NpYXRlZCB3aXRoIHRoaXMgcHJvbWlzZS5cbiAgICAgKiBUaGUgcHJvbWlzZSByZWplY3RzIHdpdGggYSB7QGxpbmsgQ2FuY2VsRXJyb3J9IGluc3RhbmNlIGFzIHJlYXNvbixcbiAgICAgKiB3aXRoIHRoZSB7QGxpbmsgQ2FuY2VsRXJyb3IjY2F1c2V9IHByb3BlcnR5IHNldCB0byB0aGUgZ2l2ZW4gYXJndW1lbnQsIGlmIGFueS5cbiAgICAgKlxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIGFmdGVyIHRoZSBwcm9taXNlIGhhcyBhbHJlYWR5IHNldHRsZWQ7XG4gICAgICogcmVwZWF0ZWQgY2FsbHMgaW4gcGFydGljdWxhciBhcmUgc2FmZSwgYnV0IG9ubHkgdGhlIGZpcnN0IG9uZVxuICAgICAqIHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXG4gICAgICpcbiAgICAgKiBUaGUgYENhbmNlbEVycm9yYCBleGNlcHRpb24gX25lZWQgbm90XyBiZSBoYW5kbGVkIGV4cGxpY2l0bHkgX29uIHRoZSBwcm9taXNlcyB0aGF0IGFyZSBiZWluZyBjYW5jZWxsZWQ6X1xuICAgICAqIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpdGggbm8gYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXIgZG9lcyBub3QgdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LlxuICAgICAqIFRoZXJlZm9yZSwgdGhlIGZvbGxvd2luZyBpZGlvbXMgYXJlIGFsbCBlcXVhbGx5IGNvcnJlY3Q6XG4gICAgICogYGBgdHNcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLmNhbmNlbCgpO1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhbmNlbCgpO1xuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhdGNoKC4uLikuY2FuY2VsKCk7XG4gICAgICogYGBgXG4gICAgICogV2hlbmV2ZXIgc29tZSBjYW5jZWxsZWQgcHJvbWlzZSBpbiBhIGNoYWluIHJlamVjdHMgd2l0aCBhIGBDYW5jZWxFcnJvcmBcbiAgICAgKiB3aXRoIHRoZSBzYW1lIGNhbmNlbGxhdGlvbiBjYXVzZSBhcyBpdHNlbGYsIHRoZSBlcnJvciB3aWxsIGJlIGRpc2NhcmRlZCBzaWxlbnRseS5cbiAgICAgKiBIb3dldmVyLCB0aGUgYENhbmNlbEVycm9yYCBfd2lsbCBzdGlsbCBiZSBkZWxpdmVyZWRfIHRvIGFsbCBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlcnNcbiAgICAgKiBhZGRlZCBieSB7QGxpbmsgdGhlbn0gYW5kIHJlbGF0ZWQgbWV0aG9kczpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gICAgICogY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS5jYXRjaChjb25zb2xlLmxvZyk7XG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIEEgQ2FuY2VsRXJyb3IgaXMgcHJpbnRlZCB0byB0aGUgY29uc29sZS5cbiAgICAgKiBgYGBcbiAgICAgKiBJZiB0aGUgYENhbmNlbEVycm9yYCBpcyBub3QgaGFuZGxlZCBkb3duc3RyZWFtIGJ5IHRoZSB0aW1lIGl0IHJlYWNoZXNcbiAgICAgKiBhIF9ub24tY2FuY2VsbGVkXyBwcm9taXNlLCBpdCBfd2lsbF8gdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LFxuICAgICAqIGp1c3QgbGlrZSBub3JtYWwgcmVqZWN0aW9ucyB3b3VsZDpcbiAgICAgKiBgYGB0c1xuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XG4gICAgICogbGV0IGNoYWluZWQgPSBjYW5jZWxsYWJsZS50aGVuKCgpID0+IHsgLi4uIH0pLnRoZW4oKCkgPT4geyAuLi4gfSk7IC8vIE5vIGNhdGNoLi4uXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIFVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQgb24gY2hhaW5lZCFcbiAgICAgKiBgYGBcbiAgICAgKiBUaGVyZWZvcmUsIGl0IGlzIGltcG9ydGFudCB0byBlaXRoZXIgY2FuY2VsIHdob2xlIHByb21pc2UgY2hhaW5zIGZyb20gdGhlaXIgdGFpbCxcbiAgICAgKiBhcyBzaG93biBpbiB0aGUgY29ycmVjdCBpZGlvbXMgYWJvdmUsIG9yIHRha2UgY2FyZSBvZiBoYW5kbGluZyBlcnJvcnMgZXZlcnl3aGVyZS5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgY2FuY2VsbGFibGUgcHJvbWlzZSB0aGF0IF9mdWxmaWxsc18gYWZ0ZXIgdGhlIGNhbmNlbCBjYWxsYmFjayAoaWYgYW55KVxuICAgICAqIGFuZCBhbGwgaGFuZGxlcnMgYXR0YWNoZWQgdXAgdG8gdGhlIGNhbGwgdG8gY2FuY2VsIGhhdmUgcnVuLlxuICAgICAqIElmIHRoZSBjYW5jZWwgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBieSBgY2FuY2VsYFxuICAgICAqIHdpbGwgYWxzbyB3YWl0IGZvciB0aGF0IHRoZW5hYmxlIHRvIHNldHRsZS5cbiAgICAgKiBUaGlzIGVuYWJsZXMgY2FsbGVycyB0byB3YWl0IGZvciB0aGUgY2FuY2VsbGVkIG9wZXJhdGlvbiB0byB0ZXJtaW5hdGVcbiAgICAgKiB3aXRob3V0IGJlaW5nIGZvcmNlZCB0byBoYW5kbGUgcG90ZW50aWFsIGVycm9ycyBhdCB0aGUgY2FsbCBzaXRlLlxuICAgICAqIGBgYHRzXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCkudGhlbigoKSA9PiB7XG4gICAgICogICAgIC8vIENsZWFudXAgZmluaXNoZWQsIGl0J3Mgc2FmZSB0byBkbyBzb21ldGhpbmcgZWxzZS5cbiAgICAgKiB9LCAoZXJyKSA9PiB7XG4gICAgICogICAgIC8vIFVucmVhY2hhYmxlOiB0aGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIGNhbmNlbCB3aWxsIG5ldmVyIHJlamVjdC5cbiAgICAgKiB9KTtcbiAgICAgKiBgYGBcbiAgICAgKiBOb3RlIHRoYXQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBfbm90XyBoYW5kbGUgaW1wbGljaXRseSBhbnkgcmVqZWN0aW9uXG4gICAgICogdGhhdCBtaWdodCBoYXZlIG9jY3VycmVkIGFscmVhZHkgaW4gdGhlIGNhbmNlbGxlZCBjaGFpbi5cbiAgICAgKiBJdCB3aWxsIGp1c3QgdHJhY2sgd2hldGhlciByZWdpc3RlcmVkIGhhbmRsZXJzIGhhdmUgYmVlbiBleGVjdXRlZCBvciBub3QuXG4gICAgICogVGhlcmVmb3JlLCB1bmhhbmRsZWQgcmVqZWN0aW9ucyB3aWxsIG5ldmVyIGJlIHNpbGVudGx5IGhhbmRsZWQgYnkgY2FsbGluZyBjYW5jZWwuXG4gICAgICovXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIC8vIElOVkFSSUFOVDogdGhlIHJlc3VsdCBvZiB0aGlzW2NhbmNlbEltcGxTeW1dIGFuZCB0aGUgYmFycmllciBkbyBub3QgZXZlciByZWplY3QuXG4gICAgICAgICAgICAvLyBVbmZvcnR1bmF0ZWx5IG1hY09TIEhpZ2ggU2llcnJhIGRvZXMgbm90IHN1cHBvcnQgUHJvbWlzZS5hbGxTZXR0bGVkLlxuICAgICAgICAgICAgUHJvbWlzZS5hbGwoW1xuICAgICAgICAgICAgICAgIHRoaXNbY2FuY2VsSW1wbFN5bV0obmV3IENhbmNlbEVycm9yKFwiUHJvbWlzZSBjYW5jZWxsZWQuXCIsIHsgY2F1c2UgfSkpLFxuICAgICAgICAgICAgICAgIGN1cnJlbnRCYXJyaWVyKHRoaXMpXG4gICAgICAgICAgICBdKS50aGVuKCgpID0+IHJlc29sdmUoKSwgKCkgPT4gcmVzb2x2ZSgpKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQmluZHMgcHJvbWlzZSBjYW5jZWxsYXRpb24gdG8gdGhlIGFib3J0IGV2ZW50IG9mIHRoZSBnaXZlbiB7QGxpbmsgQWJvcnRTaWduYWx9LlxuICAgICAqIElmIHRoZSBzaWduYWwgaGFzIGFscmVhZHkgYWJvcnRlZCwgdGhlIHByb21pc2Ugd2lsbCBiZSBjYW5jZWxsZWQgaW1tZWRpYXRlbHkuXG4gICAgICogV2hlbiBlaXRoZXIgY29uZGl0aW9uIGlzIHZlcmlmaWVkLCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlIHdpbGwgYmUgc2V0XG4gICAgICogdG8gdGhlIHNpZ25hbCdzIGFib3J0IHJlYXNvbiAoc2VlIHtAbGluayBBYm9ydFNpZ25hbCNyZWFzb259KS5cbiAgICAgKlxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIChvciBpZiB0aGUgc2lnbmFsIGFib3J0cykgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkLlxuICAgICAqIE9ubHkgdGhlIGZpcnN0IHNpZ25hbCB0byBhYm9ydCB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxuICAgICAqXG4gICAgICogRm9yIG1vcmUgZGV0YWlscyBhYm91dCB0aGUgY2FuY2VsbGF0aW9uIHByb2Nlc3MsXG4gICAgICogc2VlIHtAbGluayBjYW5jZWx9IGFuZCB0aGUgYENhbmNlbGxhYmxlUHJvbWlzZWAgY29uc3RydWN0b3IuXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBlbmFibGVzIGBhd2FpdGBpbmcgY2FuY2VsbGFibGUgcHJvbWlzZXMgd2l0aG91dCBoYXZpbmdcbiAgICAgKiB0byBzdG9yZSB0aGVtIGZvciBmdXR1cmUgY2FuY2VsbGF0aW9uLCBlLmcuOlxuICAgICAqIGBgYHRzXG4gICAgICogYXdhaXQgbG9uZ1J1bm5pbmdPcGVyYXRpb24oKS5jYW5jZWxPbihzaWduYWwpO1xuICAgICAqIGBgYFxuICAgICAqIGluc3RlYWQgb2Y6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgcHJvbWlzZVRvQmVDYW5jZWxsZWQgPSBsb25nUnVubmluZ09wZXJhdGlvbigpO1xuICAgICAqIGF3YWl0IHByb21pc2VUb0JlQ2FuY2VsbGVkO1xuICAgICAqIGBgYFxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhpcyBwcm9taXNlLCBmb3IgbWV0aG9kIGNoYWluaW5nLlxuICAgICAqL1xuICAgIGNhbmNlbE9uKHNpZ25hbDogQWJvcnRTaWduYWwpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBpZiAoc2lnbmFsLmFib3J0ZWQpIHtcbiAgICAgICAgICAgIHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbilcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHNpZ25hbC5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbiksIHtjYXB0dXJlOiB0cnVlfSk7XG4gICAgICAgIH1cblxuICAgICAgICByZXR1cm4gdGhpcztcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBBdHRhY2hlcyBjYWxsYmFja3MgZm9yIHRoZSByZXNvbHV0aW9uIGFuZC9vciByZWplY3Rpb24gb2YgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBvbmZ1bGZpbGxlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlc29sdmVkLlxuICAgICAqIEBwYXJhbSBvbnJlamVjdGVkIFRoZSBjYWxsYmFjayB0byBleGVjdXRlIHdoZW4gdGhlIFByb21pc2UgaXMgcmVqZWN0ZWQuXG4gICAgICogQHJldHVybnMgQSBgQ2FuY2VsbGFibGVQcm9taXNlYCBmb3IgdGhlIGNvbXBsZXRpb24gb2Ygd2hpY2hldmVyIGNhbGxiYWNrIGlzIGV4ZWN1dGVkLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGhvb2tlZCB1cCB0byBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHVwIHRoZSBjaGFpbiwgYnV0IG5vdCBkb3duOlxuICAgICAqXG4gICAgICogICAtIGlmIHRoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIGJlIGludm9rZWQgd2l0aCBhIGBDYW5jZWxFcnJvcmBcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIF93aWxsIHJlc29sdmUgcmVndWxhcmx5XyB3aXRoIGl0cyByZXN1bHQ7XG4gICAgICogICAtIGNvbnZlcnNlbHksIGlmIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgX3RoZSBwYXJlbnQgcHJvbWlzZSBpcyBjYW5jZWxsZWQgdG9vO19cbiAgICAgKiAgICAgdGhlIGBvbnJlamVjdGVkYCBoYW5kbGVyIHdpbGwgc3RpbGwgYmUgaW52b2tlZCB3aXRoIHRoZSBwYXJlbnQncyBgQ2FuY2VsRXJyb3JgLFxuICAgICAqICAgICBidXQgaXRzIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZFxuICAgICAqICAgICBhbmQgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhIGBDYW5jZWxFcnJvcmAgYXMgd2VsbC5cbiAgICAgKlxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXG4gICAgICpcbiAgICAgKiBJZiBlaXRoZXIgY2FsbGJhY2sgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKi9cbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+IHtcbiAgICAgICAgaWYgKCEodGhpcyBpbnN0YW5jZW9mIENhbmNlbGxhYmxlUHJvbWlzZSkpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLnRoZW4gY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIE5PVEU6IFR5cGVTY3JpcHQncyBidWlsdC1pbiB0eXBlIGZvciB0aGVuIGlzIGJyb2tlbixcbiAgICAgICAgLy8gYXMgaXQgYWxsb3dzIHNwZWNpZnlpbmcgYW4gYXJiaXRyYXJ5IFRSZXN1bHQxICE9IFQgZXZlbiB3aGVuIG9uZnVsZmlsbGVkIGlzIG5vdCBhIGZ1bmN0aW9uLlxuICAgICAgICAvLyBXZSBjYW5ub3QgZml4IGl0IGlmIHdlIHdhbnQgdG8gQ2FuY2VsbGFibGVQcm9taXNlIHRvIGltcGxlbWVudCBQcm9taXNlTGlrZTxUPi5cblxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25mdWxmaWxsZWQpKSB7IG9uZnVsZmlsbGVkID0gaWRlbnRpdHkgYXMgYW55OyB9XG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbnJlamVjdGVkKSkgeyBvbnJlamVjdGVkID0gdGhyb3dlcjsgfVxuXG4gICAgICAgIGlmIChvbmZ1bGZpbGxlZCA9PT0gaWRlbnRpdHkgJiYgb25yZWplY3RlZCA9PSB0aHJvd2VyKSB7XG4gICAgICAgICAgICAvLyBTaG9ydGN1dCBmb3IgdHJpdmlhbCBhcmd1bWVudHMuXG4gICAgICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSkgPT4gcmVzb2x2ZSh0aGlzIGFzIGFueSkpO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgYmFycmllcjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSB7fTtcbiAgICAgICAgdGhpc1tiYXJyaWVyU3ltXSA9IGJhcnJpZXI7XG5cbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgdm9pZCBzdXBlci50aGVuKFxuICAgICAgICAgICAgICAgICh2YWx1ZSkgPT4ge1xuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xuXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9uZnVsZmlsbGVkISh2YWx1ZSkpO1xuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICAocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxuICAgICAgICAgICAgICAgICAgICBiYXJyaWVyLnJlc29sdmU/LigpO1xuXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXNvbHZlKG9ucmVqZWN0ZWQhKHJlYXNvbikpO1xuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgKTtcbiAgICAgICAgfSwgYXN5bmMgKGNhdXNlPykgPT4ge1xuICAgICAgICAgICAgLy9jYW5jZWxsZWQgPSB0cnVlO1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICByZXR1cm4gb25jYW5jZWxsZWQ/LihjYXVzZSk7XG4gICAgICAgICAgICB9IGZpbmFsbHkge1xuICAgICAgICAgICAgICAgIGF3YWl0IHRoaXMuY2FuY2VsKGNhdXNlKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayBmb3Igb25seSB0aGUgcmVqZWN0aW9uIG9mIHRoZSBQcm9taXNlLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBXaGVuIHRoZSBwYXJlbnQgcHJvbWlzZSByZWplY3RzIG9yIGlzIGNhbmNlbGxlZCwgdGhlIGBvbnJlamVjdGVkYCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIEl0IGlzIGVxdWl2YWxlbnQgdG9cbiAgICAgKiBgYGB0c1xuICAgICAqIGNhbmNlbGxhYmxlUHJvbWlzZS50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xuICAgICAqIGBgYFxuICAgICAqIGFuZCB0aGUgc2FtZSBjYXZlYXRzIGFwcGx5LlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgYG9ucmVqZWN0ZWRgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxuICAgICAqL1xuICAgIGNhdGNoPFRSZXN1bHQgPSBuZXZlcj4ob25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IChQcm9taXNlTGlrZTxUUmVzdWx0PiB8IFRSZXN1bHQpKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQgfCBUUmVzdWx0PiB7XG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgYSBjYWxsYmFjayB0aGF0IGlzIGludm9rZWQgd2hlbiB0aGUgQ2FuY2VsbGFibGVQcm9taXNlIGlzIHNldHRsZWQgKGZ1bGZpbGxlZCBvciByZWplY3RlZCkuIFRoZVxuICAgICAqIHJlc29sdmVkIHZhbHVlIGNhbm5vdCBiZSBhY2Nlc3NlZCBvciBtb2RpZmllZCBmcm9tIHRoZSBjYWxsYmFjay5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHNldHRsZSBpbiB0aGUgc2FtZSBzdGF0ZSBhcyB0aGUgb3JpZ2luYWwgb25lXG4gICAgICogYWZ0ZXIgdGhlIHByb3ZpZGVkIGNhbGxiYWNrIGhhcyBjb21wbGV0ZWQgZXhlY3V0aW9uLFxuICAgICAqIHVubGVzcyB0aGUgY2FsbGJhY2sgdGhyb3dzIG9yIHJldHVybnMgYSByZWplY3RpbmcgcHJvbWlzZSxcbiAgICAgKiBpbiB3aGljaCBjYXNlIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IGFzIHdlbGwuXG4gICAgICpcbiAgICAgKiBUaGUgb3B0aW9uYWwgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCB3aWxsIGJlIGludm9rZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxuICAgICAqIE9uY2UgdGhlIHBhcmVudCBwcm9taXNlIHNldHRsZXMsIHRoZSBgb25maW5hbGx5YCBjYWxsYmFjayB3aWxsIHJ1bixcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcbiAgICAgKiBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvbi5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGlzIGltcGxlbWVudGVkIGluIHRlcm1zIG9mIHtAbGluayB0aGVufSBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cbiAgICAgKiBJdCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXG4gICAgICogQ2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoZSByZXR1cm5lZCBwcm9taXNlXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgYG9uZmluYWxseWAgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXG4gICAgICogYW5kIHRoZSBzcGVjaWZpZWQgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiBTZWUge0BsaW5rIHRoZW59IGZvciBtb3JlIGRldGFpbHMuXG4gICAgICovXG4gICAgZmluYWxseShvbmZpbmFsbHk/OiAoKCkgPT4gdm9pZCkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS5maW5hbGx5IGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25maW5hbGx5KSkge1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMudGhlbihvbmZpbmFsbHksIG9uZmluYWxseSwgb25jYW5jZWxsZWQpO1xuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHRoaXMudGhlbihcbiAgICAgICAgICAgICh2YWx1ZSkgPT4gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUob25maW5hbGx5KCkpLnRoZW4oKCkgPT4gdmFsdWUpLFxuICAgICAgICAgICAgKHJlYXNvbj8pID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHsgdGhyb3cgcmVhc29uOyB9KSxcbiAgICAgICAgICAgIG9uY2FuY2VsbGVkLFxuICAgICAgICApO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFdlIHVzZSB0aGUgYFtTeW1ib2wuc3BlY2llc11gIHN0YXRpYyBwcm9wZXJ0eSwgaWYgYXZhaWxhYmxlLFxuICAgICAqIHRvIGRpc2FibGUgdGhlIGJ1aWx0LWluIGF1dG9tYXRpYyBzdWJjbGFzc2luZyBmZWF0dXJlcyBmcm9tIHtAbGluayBQcm9taXNlfS5cbiAgICAgKiBJdCBpcyBjcml0aWNhbCBmb3IgcGVyZm9ybWFuY2UgcmVhc29ucyB0aGF0IGV4dGVuZGVycyBkbyBub3Qgb3ZlcnJpZGUgdGhpcy5cbiAgICAgKiBPbmNlIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nXG4gICAgICogaXMgZWl0aGVyIGFjY2VwdGVkIG9yIHJldGlyZWQsIHRoaXMgaW1wbGVtZW50YXRpb24gd2lsbCBoYXZlIHRvIGJlIHJldmlzZWQgYWNjb3JkaW5nbHkuXG4gICAgICpcbiAgICAgKiBAaWdub3JlXG4gICAgICogQGludGVybmFsXG4gICAgICovXG4gICAgc3RhdGljIGdldCBbc3BlY2llc10oKSB7XG4gICAgICAgIHJldHVybiBQcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBpcyByZXNvbHZlZCB3aXRoIGFuIGFycmF5IG9mIHJlc3VsdHNcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSwgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgUHJvbWlzZSBpcyByZWplY3RlZC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYWxsPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+W10+O1xuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IEF3YWl0ZWQ8VFtQXT47IH0+O1xuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbGwoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUgb3IgcmVqZWN0LlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxQcm9taXNlU2V0dGxlZFJlc3VsdDxBd2FpdGVkPFQ+PltdPjtcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx7IC1yZWFkb25seSBbUCBpbiBrZXlvZiBUXTogUHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUW1BdPj47IH0+O1xuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsU2V0dGxlZChjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRoZSBhbnkgZnVuY3Rpb24gcmV0dXJucyBhIHByb21pc2UgdGhhdCBpcyBmdWxmaWxsZWQgYnkgdGhlIGZpcnN0IGdpdmVuIHByb21pc2UgdG8gYmUgZnVsZmlsbGVkLFxuICAgICAqIG9yIHJlamVjdGVkIHdpdGggYW4gQWdncmVnYXRlRXJyb3IgY29udGFpbmluZyBhbiBhcnJheSBvZiByZWplY3Rpb24gcmVhc29uc1xuICAgICAqIGlmIGFsbCBvZiB0aGUgZ2l2ZW4gcHJvbWlzZXMgYXJlIHJlamVjdGVkLlxuICAgICAqIEl0IHJlc29sdmVzIGFsbCBlbGVtZW50cyBvZiB0aGUgcGFzc2VkIGl0ZXJhYmxlIHRvIHByb21pc2VzIGFzIGl0IHJ1bnMgdGhpcyBhbGdvcml0aG0uXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFueTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIGFueTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gY29sbGVjdGVkLmxlbmd0aCA9PT0gMFxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5hbnkoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIG9yIHJlamVjdGVkIHdoZW4gYW55IG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyBhcmUgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJhY2U8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XG4gICAgc3RhdGljIHJhY2U8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmFjZShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcbiAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBjYW5jZWxsZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgY2FuY2VsPFQgPSBuZXZlcj4oY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBjb25zdCBwID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XG4gICAgICAgIHAuY2FuY2VsKGNhdXNlKTtcbiAgICAgICAgcmV0dXJuIHA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBjYW5jZWxzXG4gICAgICogYWZ0ZXIgdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCBjYXVzZS5cbiAgICAgKlxuICAgICAqIElmIHRoZSB7QGxpbmsgQWJvcnRTaWduYWwudGltZW91dH0gZmFjdG9yeSBtZXRob2QgaXMgYXZhaWxhYmxlLFxuICAgICAqIGl0IGlzIHVzZWQgdG8gYmFzZSB0aGUgdGltZW91dCBvbiBfYWN0aXZlXyB0aW1lIHJhdGhlciB0aGFuIF9lbGFwc2VkXyB0aW1lLlxuICAgICAqIE90aGVyd2lzZSwgYHRpbWVvdXRgIGZhbGxzIGJhY2sgdG8ge0BsaW5rIHNldFRpbWVvdXR9LlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHRpbWVvdXQ8VCA9IG5ldmVyPihtaWxsaXNlY29uZHM6IG51bWJlciwgY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XG4gICAgICAgIGlmIChBYm9ydFNpZ25hbCAmJiB0eXBlb2YgQWJvcnRTaWduYWwgPT09ICdmdW5jdGlvbicgJiYgQWJvcnRTaWduYWwudGltZW91dCAmJiB0eXBlb2YgQWJvcnRTaWduYWwudGltZW91dCA9PT0gJ2Z1bmN0aW9uJykge1xuICAgICAgICAgICAgQWJvcnRTaWduYWwudGltZW91dChtaWxsaXNlY29uZHMpLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiB2b2lkIHByb21pc2UuY2FuY2VsKGNhdXNlKSwgbWlsbGlzZWNvbmRzKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dC5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHNsZWVwKG1pbGxpc2Vjb25kczogbnVtYmVyKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXJcbiAgICAgKiB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGNhbiBiZSBjYW5jZWxsZWQgd2l0aG91dCBjb25zZXF1ZW5jZXMuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgc2xlZXA8VD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xuICAgIHN0YXRpYyBzbGVlcDxUID0gdm9pZD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlPzogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHJlc29sdmUodmFsdWUhKSwgbWlsbGlzZWNvbmRzKTtcbiAgICAgICAgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZWplY3RlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCByZWFzb24uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVqZWN0PFQgPSBuZXZlcj4ocmVhc29uPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKF8sIHJlamVjdCkgPT4gcmVqZWN0KHJlYXNvbikpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmUoKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVzb2x2ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgdmFsdWUuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgcmVzb2x2ZTxUID0gdm9pZD4odmFsdWU/OiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj4ge1xuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpIHtcbiAgICAgICAgICAgIC8vIE9wdGltaXNlIGZvciBjYW5jZWxsYWJsZSBwcm9taXNlcy5cbiAgICAgICAgICAgIHJldHVybiB2YWx1ZTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+KChyZXNvbHZlKSA9PiByZXNvbHZlKHZhbHVlKSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgYW5kIHJldHVybnMgaXQgaW4gYW4gb2JqZWN0LCBhbG9uZyB3aXRoIGl0cyByZXNvbHZlIGFuZCByZWplY3QgZnVuY3Rpb25zXG4gICAgICogYW5kIGEgZ2V0dGVyL3NldHRlciBmb3IgdGhlIGNhbmNlbGxhdGlvbiBjYWxsYmFjay5cbiAgICAgKlxuICAgICAqIFRoaXMgbWV0aG9kIGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgd2l0aFJlc29sdmVyczxUPigpOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICAgICAgbGV0IHJlc3VsdDogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHsgb25jYW5jZWxsZWQ6IG51bGwgfSBhcyBhbnk7XG4gICAgICAgIHJlc3VsdC5wcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICByZXN1bHQucmVzb2x2ZSA9IHJlc29sdmU7XG4gICAgICAgICAgICByZXN1bHQucmVqZWN0ID0gcmVqZWN0O1xuICAgICAgICB9LCAoY2F1c2U/OiBhbnkpID0+IHsgcmVzdWx0Lm9uY2FuY2VsbGVkPy4oY2F1c2UpOyB9KTtcbiAgICAgICAgcmV0dXJuIHJlc3VsdDtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgY2FuY2VsbGF0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSByZXN1bHRpbmcgZnVuY3Rpb24gZG9lcyBub3QgcmVqZWN0LlxuICovXG5mdW5jdGlvbiBjYW5jZWxsZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKSB7XG4gICAgbGV0IGNhbmNlbGxhdGlvblByb21pc2U6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9IHVuZGVmaW5lZDtcblxuICAgIHJldHVybiAocmVhc29uOiBDYW5jZWxFcnJvcik6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPiA9PiB7XG4gICAgICAgIGlmICghc3RhdGUuc2V0dGxlZCkge1xuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBzdGF0ZS5yZWFzb24gPSByZWFzb247XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xuXG4gICAgICAgICAgICAvLyBBdHRhY2ggYW4gZXJyb3IgaGFuZGxlciB0aGF0IGlnbm9yZXMgdGhpcyBzcGVjaWZpYyByZWplY3Rpb24gcmVhc29uIGFuZCBub3RoaW5nIGVsc2UuXG4gICAgICAgICAgICAvLyBJbiB0aGVvcnksIGEgc2FuZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uIGF0IHRoaXMgcG9pbnRcbiAgICAgICAgICAgIC8vIHNob3VsZCBhbHdheXMgcmVqZWN0IHdpdGggb3VyIGNhbmNlbGxhdGlvbiByZWFzb24sXG4gICAgICAgICAgICAvLyBoZW5jZSB0aGUgaGFuZGxlciB3aWxsIG5ldmVyIHRocm93LlxuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnByb3RvdHlwZS50aGVuLmNhbGwocHJvbWlzZS5wcm9taXNlLCB1bmRlZmluZWQsIChlcnIpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoZXJyICE9PSByZWFzb24pIHtcbiAgICAgICAgICAgICAgICAgICAgdGhyb3cgZXJyO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH0pO1xuICAgICAgICB9XG5cbiAgICAgICAgLy8gSWYgcmVhc29uIGlzIG5vdCBzZXQsIHRoZSBwcm9taXNlIHJlc29sdmVkIHJlZ3VsYXJseSwgaGVuY2Ugd2UgbXVzdCBub3QgY2FsbCBvbmNhbmNlbGxlZC5cbiAgICAgICAgLy8gSWYgb25jYW5jZWxsZWQgaXMgdW5zZXQsIG5vIG5lZWQgdG8gZ28gYW55IGZ1cnRoZXIuXG4gICAgICAgIGlmICghc3RhdGUucmVhc29uIHx8ICFwcm9taXNlLm9uY2FuY2VsbGVkKSB7IHJldHVybjsgfVxuXG4gICAgICAgIGNhbmNlbGxhdGlvblByb21pc2UgPSBuZXcgUHJvbWlzZTx2b2lkPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICByZXNvbHZlKHByb21pc2Uub25jYW5jZWxsZWQhKHN0YXRlLnJlYXNvbiEuY2F1c2UpKTtcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIG9uY2FuY2VsbGVkIGNhbGxiYWNrLlwiKSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pLmNhdGNoKChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xuICAgICAgICB9KTtcblxuICAgICAgICAvLyBVbnNldCBvbmNhbmNlbGxlZCB0byBwcmV2ZW50IHJlcGVhdGVkIGNhbGxzLlxuICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gbnVsbDtcblxuICAgICAgICByZXR1cm4gY2FuY2VsbGF0aW9uUHJvbWlzZTtcbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVzb2x1dGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICovXG5mdW5jdGlvbiByZXNvbHZlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiB7XG4gICAgcmV0dXJuICh2YWx1ZSkgPT4ge1xuICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xuXG4gICAgICAgIGlmICh2YWx1ZSA9PT0gcHJvbWlzZS5wcm9taXNlKSB7XG4gICAgICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QobmV3IFR5cGVFcnJvcihcIkEgcHJvbWlzZSBjYW5ub3QgYmUgcmVzb2x2ZWQgd2l0aCBpdHNlbGYuXCIpKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh2YWx1ZSAhPSBudWxsICYmICh0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ2Z1bmN0aW9uJykpIHtcbiAgICAgICAgICAgIGxldCB0aGVuOiBhbnk7XG4gICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgIHRoZW4gPSAodmFsdWUgYXMgYW55KS50aGVuO1xuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKHRoZW4pKSB7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgbGV0IGNhbmNlbCA9ICh2YWx1ZSBhcyBhbnkpLmNhbmNlbDtcbiAgICAgICAgICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUoY2FuY2VsKSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3Qgb25jYW5jZWxsZWQgPSAoY2F1c2U/OiBhbnkpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xuICAgICAgICAgICAgICAgICAgICAgICAgfTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChzdGF0ZS5yZWFzb24pIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBJZiBhbHJlYWR5IGNhbmNlbGxlZCwgcHJvcGFnYXRlIGNhbmNlbGxhdGlvbi5cbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHRoZSBjYW5jZWxsZXIgYWxnb3JpdGhtIGRvZXMgbm90IHJlamVjdFxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIHNvIGl0IGNhbiBiZSBkaXNjYXJkZWQgc2FmZWx5LlxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgY2FuY2VsbGVyRm9yKHsgLi4ucHJvbWlzZSwgb25jYW5jZWxsZWQgfSwgc3RhdGUpKHN0YXRlLnJlYXNvbik7XG4gICAgICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHByb21pc2Uub25jYW5jZWxsZWQgPSBvbmNhbmNlbGxlZDtcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0gY2F0Y2gge31cblxuICAgICAgICAgICAgICAgIGNvbnN0IG5ld1N0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcbiAgICAgICAgICAgICAgICAgICAgcm9vdDogc3RhdGUucm9vdCxcbiAgICAgICAgICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHNldHRsZWQoKSB7IHJldHVybiB0aGlzLnJvb3Quc2V0dGxlZCB9LFxuICAgICAgICAgICAgICAgICAgICBzZXQgc2V0dGxlZCh2YWx1ZSkgeyB0aGlzLnJvb3Quc2V0dGxlZCA9IHZhbHVlOyB9LFxuICAgICAgICAgICAgICAgICAgICBnZXQgcmVhc29uKCkgeyByZXR1cm4gdGhpcy5yb290LnJlYXNvbiB9XG4gICAgICAgICAgICAgICAgfTtcblxuICAgICAgICAgICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgbmV3U3RhdGUpO1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIFJlZmxlY3QuYXBwbHkodGhlbiwgdmFsdWUsIFtyZXNvbHZlckZvcihwcm9taXNlLCBuZXdTdGF0ZSksIHJlamVjdG9yXSk7XG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJldHVybjsgLy8gSU1QT1JUQU5UIVxuICAgICAgICAgICAgfVxuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICBwcm9taXNlLnJlc29sdmUodmFsdWUpO1xuICAgIH07XG59XG5cbi8qKlxuICogUmV0dXJucyBhIGNhbGxiYWNrIHRoYXQgaW1wbGVtZW50cyB0aGUgcmVqZWN0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXG4gKi9cbmZ1bmN0aW9uIHJlamVjdG9yRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSk6IENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yIHtcbiAgICByZXR1cm4gKHJlYXNvbj8pID0+IHtcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcblxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkge1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICBpZiAocmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgc3RhdGUucmVhc29uIGluc3RhbmNlb2YgQ2FuY2VsRXJyb3IgJiYgT2JqZWN0LmlzKHJlYXNvbi5jYXVzZSwgc3RhdGUucmVhc29uLmNhdXNlKSkge1xuICAgICAgICAgICAgICAgICAgICAvLyBTd2FsbG93IGxhdGUgcmVqZWN0aW9ucyB0aGF0IGFyZSBDYW5jZWxFcnJvcnMgd2hvc2UgY2FuY2VsbGF0aW9uIGNhdXNlIGlzIHRoZSBzYW1lIGFzIG91cnMuXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9IGNhdGNoIHt9XG5cbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KHJlYXNvbik7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogQ2FuY2VscyBhbGwgdmFsdWVzIGluIGFuIGFycmF5IHRoYXQgbG9vayBsaWtlIGNhbmNlbGxhYmxlIHRoZW5hYmxlcy5cbiAqIFJldHVybnMgYSBwcm9taXNlIHRoYXQgZnVsZmlsbHMgb25jZSBhbGwgY2FuY2VsbGF0aW9uIHByb2NlZHVyZXMgZm9yIHRoZSBnaXZlbiB2YWx1ZXMgaGF2ZSBzZXR0bGVkLlxuICovXG5mdW5jdGlvbiBjYW5jZWxBbGwocGFyZW50OiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4sIHZhbHVlczogYW55W10sIGNhdXNlPzogYW55KTogUHJvbWlzZTx2b2lkPiB7XG4gICAgY29uc3QgcmVzdWx0czogUHJvbWlzZTx2b2lkPltdID0gW107XG5cbiAgICBmb3IgKGNvbnN0IHZhbHVlIG9mIHZhbHVlcykge1xuICAgICAgICBsZXQgY2FuY2VsOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXI7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUodmFsdWUudGhlbikpIHsgY29udGludWU7IH1cbiAgICAgICAgICAgIGNhbmNlbCA9IHZhbHVlLmNhbmNlbDtcbiAgICAgICAgICAgIGlmICghaXNDYWxsYWJsZShjYW5jZWwpKSB7IGNvbnRpbnVlOyB9XG4gICAgICAgIH0gY2F0Y2ggeyBjb250aW51ZTsgfVxuXG4gICAgICAgIGxldCByZXN1bHQ6IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHJlc3VsdCA9IFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgZXJyLCBcIlVuaGFuZGxlZCBleGNlcHRpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xuICAgICAgICAgICAgY29udGludWU7XG4gICAgICAgIH1cblxuICAgICAgICBpZiAoIXJlc3VsdCkgeyBjb250aW51ZTsgfVxuICAgICAgICByZXN1bHRzLnB1c2goXG4gICAgICAgICAgICAocmVzdWx0IGluc3RhbmNlb2YgUHJvbWlzZSAgPyByZXN1bHQgOiBQcm9taXNlLnJlc29sdmUocmVzdWx0KSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICApO1xuICAgIH1cblxuICAgIHJldHVybiBQcm9taXNlLmFsbChyZXN1bHRzKSBhcyBhbnk7XG59XG5cbi8qKlxuICogUmV0dXJucyBpdHMgYXJndW1lbnQuXG4gKi9cbmZ1bmN0aW9uIGlkZW50aXR5PFQ+KHg6IFQpOiBUIHtcbiAgICByZXR1cm4geDtcbn1cblxuLyoqXG4gKiBUaHJvd3MgaXRzIGFyZ3VtZW50LlxuICovXG5mdW5jdGlvbiB0aHJvd2VyKHJlYXNvbj86IGFueSk6IG5ldmVyIHtcbiAgICB0aHJvdyByZWFzb247XG59XG5cbi8qKlxuICogQXR0ZW1wdHMgdmFyaW91cyBzdHJhdGVnaWVzIHRvIGNvbnZlcnQgYW4gZXJyb3IgdG8gYSBzdHJpbmcuXG4gKi9cbmZ1bmN0aW9uIGVycm9yTWVzc2FnZShlcnI6IGFueSk6IHN0cmluZyB7XG4gICAgdHJ5IHtcbiAgICAgICAgaWYgKGVyciBpbnN0YW5jZW9mIEVycm9yIHx8IHR5cGVvZiBlcnIgIT09ICdvYmplY3QnIHx8IGVyci50b1N0cmluZyAhPT0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZykge1xuICAgICAgICAgICAgcmV0dXJuIFwiXCIgKyBlcnI7XG4gICAgICAgIH1cbiAgICB9IGNhdGNoIHt9XG5cbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gSlNPTi5zdHJpbmdpZnkoZXJyKTtcbiAgICB9IGNhdGNoIHt9XG5cbiAgICB0cnkge1xuICAgICAgICByZXR1cm4gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZy5jYWxsKGVycik7XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgcmV0dXJuIFwiPGNvdWxkIG5vdCBjb252ZXJ0IGVycm9yIHRvIHN0cmluZz5cIjtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGJhcnJpZXIgcHJvbWlzZSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuIElmIG5lY2Vzc2FyeSwgaW5pdGlhbGlzZXMgdGhlIGJhcnJpZXIuXG4gKi9cbmZ1bmN0aW9uIGN1cnJlbnRCYXJyaWVyPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTxUPik6IFByb21pc2U8dm9pZD4ge1xuICAgIGxldCBwd3I6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0gcHJvbWlzZVtiYXJyaWVyU3ltXSA/PyB7fTtcbiAgICBpZiAoISgncHJvbWlzZScgaW4gcHdyKSkge1xuICAgICAgICBPYmplY3QuYXNzaWduKHB3ciwgcHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4oKSk7XG4gICAgfVxuICAgIGlmIChwcm9taXNlW2JhcnJpZXJTeW1dID09IG51bGwpIHtcbiAgICAgICAgcHdyLnJlc29sdmUhKCk7XG4gICAgICAgIHByb21pc2VbYmFycmllclN5bV0gPSBwd3I7XG4gICAgfVxuICAgIHJldHVybiBwd3IucHJvbWlzZSE7XG59XG5cbi8vIFBvbHlmaWxsIFByb21pc2Uud2l0aFJlc29sdmVycy5cbmxldCBwcm9taXNlV2l0aFJlc29sdmVycyA9IFByb21pc2Uud2l0aFJlc29sdmVycztcbmlmIChwcm9taXNlV2l0aFJlc29sdmVycyAmJiB0eXBlb2YgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPT09ICdmdW5jdGlvbicpIHtcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IHByb21pc2VXaXRoUmVzb2x2ZXJzLmJpbmQoUHJvbWlzZSk7XG59IGVsc2Uge1xuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gZnVuY3Rpb24gPFQ+KCk6IFByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcbiAgICAgICAgbGV0IHJlc29sdmUhOiAodmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgUHJvbWlzZTxUPigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcbiAgICAgICAgcmV0dXJuIHsgcHJvbWlzZSwgcmVzb2x2ZSwgcmVqZWN0IH07XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNsaXBib2FyZCk7XG5cbmNvbnN0IENsaXBib2FyZFNldFRleHQgPSAwO1xuY29uc3QgQ2xpcGJvYXJkVGV4dCA9IDE7XG5cbi8qKlxuICogU2V0cyB0aGUgdGV4dCB0byB0aGUgQ2xpcGJvYXJkLlxuICpcbiAqIEBwYXJhbSB0ZXh0IC0gVGhlIHRleHQgdG8gYmUgc2V0IHRvIHRoZSBDbGlwYm9hcmQuXG4gKiBAcmV0dXJuIEEgUHJvbWlzZSB0aGF0IHJlc29sdmVzIHdoZW4gdGhlIG9wZXJhdGlvbiBpcyBzdWNjZXNzZnVsLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0OiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRTZXRUZXh0LCB7dGV4dH0pO1xufVxuXG4vKipcbiAqIEdldCB0aGUgQ2xpcGJvYXJkIHRleHRcbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSB0ZXh0IGZyb20gdGhlIENsaXBib2FyZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKTogUHJvbWlzZTxzdHJpbmc+IHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuZXhwb3J0IGludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cbiAgICBXaWR0aDogbnVtYmVyO1xuICAgIC8qKiBUaGUgaGVpZ2h0IG9mIGEgcmVjdGFuZ3VsYXIgYXJlYS4gKi9cbiAgICBIZWlnaHQ6IG51bWJlcjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBSZWN0IHtcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xuICAgIFg6IG51bWJlcjtcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xuICAgIFk6IG51bWJlcjtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSByZWN0YW5nbGUuICovXG4gICAgV2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xuICAgIEhlaWdodDogbnVtYmVyO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIFNjcmVlbiB7XG4gICAgLyoqIFVuaXF1ZSBpZGVudGlmaWVyIGZvciB0aGUgc2NyZWVuLiAqL1xuICAgIElEOiBzdHJpbmc7XG4gICAgLyoqIEh1bWFuLXJlYWRhYmxlIG5hbWUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBOYW1lOiBzdHJpbmc7XG4gICAgLyoqIFRoZSBzY2FsZSBmYWN0b3Igb2YgdGhlIHNjcmVlbiAoRFBJLzk2KS4gMSA9IHN0YW5kYXJkIERQSSwgMiA9IEhpRFBJIChSZXRpbmEpLCBldGMuICovXG4gICAgU2NhbGVGYWN0b3I6IG51bWJlcjtcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFg6IG51bWJlcjtcbiAgICAvKiogVGhlIFkgY29vcmRpbmF0ZSBvZiB0aGUgc2NyZWVuLiAqL1xuICAgIFk6IG51bWJlcjtcbiAgICAvKiogQ29udGFpbnMgdGhlIHdpZHRoIGFuZCBoZWlnaHQgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBTaXplOiBTaXplO1xuICAgIC8qKiBDb250YWlucyB0aGUgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQuICovXG4gICAgQm91bmRzOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgYm91bmRzIG9mIHRoZSBzY3JlZW4gaW4gdGVybXMgb2YgWCwgWSwgV2lkdGgsIGFuZCBIZWlnaHQgKGJlZm9yZSBzY2FsaW5nKS4gKi9cbiAgICBQaHlzaWNhbEJvdW5kczogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIGFyZWEgb2YgdGhlIHNjcmVlbiB0aGF0IGlzIGFjdHVhbGx5IHVzYWJsZSAoZXhjbHVkaW5nIHRhc2tiYXIgYW5kIG90aGVyIHN5c3RlbSBVSSkuICovXG4gICAgV29ya0FyZWE6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBXb3JrQXJlYSBvZiB0aGUgc2NyZWVuIChiZWZvcmUgc2NhbGluZykuICovXG4gICAgUGh5c2ljYWxXb3JrQXJlYTogUmVjdDtcbiAgICAvKiogVHJ1ZSBpZiB0aGlzIGlzIHRoZSBwcmltYXJ5IG1vbml0b3Igc2VsZWN0ZWQgYnkgdGhlIHVzZXIgaW4gdGhlIG9wZXJhdGluZyBzeXN0ZW0uICovXG4gICAgSXNQcmltYXJ5OiBib29sZWFuO1xuICAgIC8qKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBSb3RhdGlvbjogbnVtYmVyO1xufVxuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlNjcmVlbnMpO1xuXG5jb25zdCBnZXRBbGwgPSAwO1xuY29uc3QgZ2V0UHJpbWFyeSA9IDE7XG5jb25zdCBnZXRDdXJyZW50ID0gMjtcbmNvbnN0IGdldEJ5SUQgPSAzO1xuY29uc3QgZ2V0QnlJbmRleCA9IDQ7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpOiBQcm9taXNlPFNjcmVlbltdPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcbn1cblxuLyoqXG4gKiBHZXRzIGEgc2NyZWVuIGJ5IGl0cyB1bmlxdWUgZGlzcGxheSBJRC5cbiAqXG4gKiBAcGFyYW0gaWQgLSBUaGUgdW5pcXVlIGlkZW50aWZpZXIgb2YgdGhlIHNjcmVlbi5cbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIHRoZSBtYXRjaGluZyBTY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRCeUlEKGlkOiBzdHJpbmcpOiBQcm9taXNlPFNjcmVlbj4ge1xuICAgIHJldHVybiBjYWxsKGdldEJ5SUQsIHsgaWQgfSk7XG59XG5cbi8qKlxuICogR2V0cyBhIHNjcmVlbiBieSBpdHMgaW5kZXggaW4gdGhlIHNjcmVlbiBsaXN0LlxuICpcbiAqIEBwYXJhbSBpbmRleCAtIFRoZSB6ZXJvLWJhc2VkIGluZGV4IG9mIHRoZSBzY3JlZW4uXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgbWF0Y2hpbmcgU2NyZWVuLlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QnlJbmRleChpbmRleDogbnVtYmVyKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRCeUluZGV4LCB7IGluZGV4IH0pO1xufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuSU9TKTtcblxuLy8gTWV0aG9kIElEc1xuY29uc3QgSGFwdGljc0ltcGFjdCA9IDA7XG5jb25zdCBEZXZpY2VJbmZvID0gMTtcblxuZXhwb3J0IG5hbWVzcGFjZSBIYXB0aWNzIHtcbiAgICBleHBvcnQgdHlwZSBJbXBhY3RTdHlsZSA9IFwibGlnaHRcInxcIm1lZGl1bVwifFwiaGVhdnlcInxcInNvZnRcInxcInJpZ2lkXCI7XG4gICAgZXhwb3J0IGZ1bmN0aW9uIEltcGFjdChzdHlsZTogSW1wYWN0U3R5bGUgPSBcIm1lZGl1bVwiKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiBjYWxsKEhhcHRpY3NJbXBhY3QsIHsgc3R5bGUgfSk7XG4gICAgfVxufVxuXG5leHBvcnQgbmFtZXNwYWNlIERldmljZSB7XG4gICAgZXhwb3J0IGludGVyZmFjZSBJbmZvIHtcbiAgICAgICAgbW9kZWw6IHN0cmluZztcbiAgICAgICAgc3lzdGVtTmFtZTogc3RyaW5nO1xuICAgICAgICBzeXN0ZW1WZXJzaW9uOiBzdHJpbmc7XG4gICAgICAgIGlzU2ltdWxhdG9yOiBib29sZWFuO1xuICAgIH1cbiAgICBleHBvcnQgZnVuY3Rpb24gSW5mbygpOiBQcm9taXNlPEluZm8+IHtcbiAgICAgICAgcmV0dXJuIGNhbGwoRGV2aWNlSW5mbyk7XG4gICAgfVxufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQW5kcm9pZCk7XG5cbi8vIE1ldGhvZCBJRHMgKG11c3QgbWF0Y2ggbWVzc2FnZXByb2Nlc3Nvcl9hbmRyb2lkLmdvKVxuY29uc3QgSGFwdGljc1ZpYnJhdGUgPSAwO1xuY29uc3QgRGV2aWNlSW5mbyA9IDE7XG5jb25zdCBUb2FzdFNob3cgPSAyO1xuXG5leHBvcnQgbmFtZXNwYWNlIEhhcHRpY3Mge1xuICAgIC8qKiBWaWJyYXRlIHRoZSBkZXZpY2UgZm9yIHRoZSBnaXZlbiBkdXJhdGlvbiBpbiBtaWxsaXNlY29uZHMuICovXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFZpYnJhdGUoZHVyYXRpb25NczogbnVtYmVyID0gMTAwKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiBjYWxsKEhhcHRpY3NWaWJyYXRlLCB7IGR1cmF0aW9uOiBkdXJhdGlvbk1zIH0pO1xuICAgIH1cbn1cblxuZXhwb3J0IG5hbWVzcGFjZSBEZXZpY2Uge1xuICAgIGV4cG9ydCBpbnRlcmZhY2UgSW5mbyB7XG4gICAgICAgIHBsYXRmb3JtOiBzdHJpbmc7XG4gICAgICAgIG1hbnVmYWN0dXJlcjogc3RyaW5nO1xuICAgICAgICBicmFuZDogc3RyaW5nO1xuICAgICAgICBtb2RlbDogc3RyaW5nO1xuICAgICAgICBkZXZpY2U6IHN0cmluZztcbiAgICAgICAgdmVyc2lvbjogc3RyaW5nO1xuICAgICAgICBzZGtJbnQ6IG51bWJlcjtcbiAgICB9XG4gICAgLyoqIFJldHVybiBpbmZvcm1hdGlvbiBhYm91dCB0aGUgQW5kcm9pZCBkZXZpY2UuICovXG4gICAgZXhwb3J0IGZ1bmN0aW9uIEluZm8oKTogUHJvbWlzZTxJbmZvPiB7XG4gICAgICAgIHJldHVybiBjYWxsKERldmljZUluZm8pO1xuICAgIH1cbn1cblxuZXhwb3J0IG5hbWVzcGFjZSBUb2FzdCB7XG4gICAgLyoqIFNob3cgYSBzaG9ydCBBbmRyb2lkIHRvYXN0IG1lc3NhZ2UuICovXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFNob3cobWVzc2FnZTogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiBjYWxsKFRvYXN0U2hvdywgeyBtZXNzYWdlIH0pO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBVcGRhdGVyIGV2ZW50IG5hbWUgY29uc3RhbnRzLlxuICpcbiAqIFVzZSB0aGVzZSBpbnN0ZWFkIG9mIGhhcmQtY29kaW5nIHN0cmluZyBsaXRlcmFscyB3aGVuIHN1YnNjcmliaW5nIHRvXG4gKiB1cGRhdGVyIGV2ZW50cyBmcm9tIEphdmFTY3JpcHQ6XG4gKlxuICogICAgIGltcG9ydCB7IEV2ZW50cywgVXBkYXRlciB9IGZyb20gXCJAd2FpbHNpby9ydW50aW1lXCI7XG4gKlxuICogICAgIEV2ZW50cy5PbihVcGRhdGVyLkV2ZW50cy5VcGRhdGVBdmFpbGFibGUsIChlKSA9PiB7XG4gKiAgICAgICAgIGNvbnNvbGUubG9nKFwidXBkYXRlIGZvdW5kOlwiLCBlLmRhdGEudmVyc2lvbik7XG4gKiAgICAgfSk7XG4gKlxuICogICAgIEV2ZW50cy5PbihVcGRhdGVyLkV2ZW50cy5Eb3dubG9hZFByb2dyZXNzLCAoZSkgPT4ge1xuICogICAgICAgICBjb25zdCBwID0gZS5kYXRhO1xuICogICAgICAgICBjb25zb2xlLmxvZyhgJHtwLndyaXR0ZW59IC8gJHtwLnRvdGFsfSBieXRlc2ApO1xuICogICAgIH0pO1xuICpcbiAqIE1pcnJvcnMgdGhlIEdvLXNpZGUgY29uc3RhbnRzIGluIGBwa2cvdXBkYXRlci9ldmVudHMuZ29gIGFuZCB0aGVcbiAqIHVzZXItYWN0aW9uIGNvbnN0YW50cyBpbiBgcGtnL3VwZGF0ZXIvd2luZG93X2xpZmVjeWNsZS5nb2AuIEFueVxuICogY2hhbmdlcyBoZXJlIG11c3Qgc3RheSBpbiBzeW5jIHdpdGggdGhvc2UgZmlsZXMgXHUyMDE0IHRoZXJlJ3MgYW5cbiAqIGludGVncmF0aW9uIHRlc3QgdGhhdCBhc3NlcnRzIHRoZSBzdHJpbmdzIG1hdGNoLlxuICovXG5leHBvcnQgY29uc3QgRXZlbnRzID0gT2JqZWN0LmZyZWV6ZSh7XG4gICAgLyoqIEEgQ2hlY2sgcm91bmQtdHJpcCBpcyBzdGFydGluZy4gUGF5bG9hZDogbnVsbC4gKi9cbiAgICBDaGVja1N0YXJ0ZWQ6IFwid2FpbHM6dXBkYXRlcjpjaGVjay1zdGFydGVkXCIsXG4gICAgLyoqIENoZWNrIGZvdW5kIGEgbmV3ZXIgcmVsZWFzZS4gUGF5bG9hZDogUmVsZWFzZS4gKi9cbiAgICBVcGRhdGVBdmFpbGFibGU6IFwid2FpbHM6dXBkYXRlcjp1cGRhdGUtYXZhaWxhYmxlXCIsXG4gICAgLyoqIENoZWNrIGNvbmZpcm1lZCB0aGUgY2FsbGVyIGlzIHVwIHRvIGRhdGUuIFBheWxvYWQ6IG51bGwuICovXG4gICAgTm9VcGRhdGU6IFwid2FpbHM6dXBkYXRlcjpuby11cGRhdGVcIixcbiAgICAvKiogRG93bmxvYWQgaXMgc3RhcnRpbmcuIFBheWxvYWQ6IFJlbGVhc2UuICovXG4gICAgRG93bmxvYWRTdGFydGVkOiBcIndhaWxzOnVwZGF0ZXI6ZG93bmxvYWQtc3RhcnRlZFwiLFxuICAgIC8qKiBQZXJpb2RpYyBwcm9ncmVzcyB0aWNrIGR1cmluZyBkb3dubG9hZCAofjEwIEh6KS4gUGF5bG9hZDogUHJvZ3Jlc3MuICovXG4gICAgRG93bmxvYWRQcm9ncmVzczogXCJ3YWlsczp1cGRhdGVyOmRvd25sb2FkLXByb2dyZXNzXCIsXG4gICAgLyoqIEFsbCBieXRlcyBhcmUgb24gZGlzaywgYnV0IHZlcmlmaWNhdGlvbiBoYXMgbm90IHlldCBzdGFydGVkLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xuICAgIERvd25sb2FkQ29tcGxldGU6IFwid2FpbHM6dXBkYXRlcjpkb3dubG9hZC1jb21wbGV0ZVwiLFxuICAgIC8qKiBTaWduYXR1cmUgLyBkaWdlc3QgdmVyaWZpY2F0aW9uIGhhcyBzdGFydGVkLiBQYXlsb2FkOiBSZWxlYXNlLiAqL1xuICAgIFZlcmlmeWluZzogXCJ3YWlsczp1cGRhdGVyOnZlcmlmeWluZ1wiLFxuICAgIC8qKiBUaGUgVXBkYXRlciBpcyBzd2FwcGluZyB0aGUgYmluYXJ5IGludG8gcGxhY2UuIFBheWxvYWQ6IFJlbGVhc2UuICovXG4gICAgSW5zdGFsbGluZzogXCJ3YWlsczp1cGRhdGVyOmluc3RhbGxpbmdcIixcbiAgICAvKiogVXBkYXRlIGlzIHN0YWdlZCBhbmQgYSByZXN0YXJ0IGlzIHBlbmRpbmcuIFBheWxvYWQ6IFJlbGVhc2UuICovXG4gICAgVXBkYXRlUmVhZHk6IFwid2FpbHM6dXBkYXRlcjp1cGRhdGUtcmVhZHlcIixcbiAgICAvKiogU29tZXRoaW5nIGZhaWxlZC4gUGF5bG9hZDogRXJyb3JJbmZvIHsgc3RhZ2UsIG1lc3NhZ2UsIHByb3ZpZGVyIH0uICovXG4gICAgRXJyb3I6IFwid2FpbHM6dXBkYXRlcjplcnJvclwiLFxuICAgIC8qKiBIb3N0LXNpZGUgY29udGV4dCBkZWxpdmVyZWQgb25jZSBwZXIgc2Vzc2lvbi4gUGF5bG9hZDogTWV0YSB7IGN1cnJlbnRWZXJzaW9uLCBza2lwcGVkVmVyc2lvbiB9LiAqL1xuICAgIE1ldGE6IFwid2FpbHM6dXBkYXRlcjptZXRhXCIsXG5cbiAgICAvKiogU3ViLW5hbWVzcGFjZTogdXNlci1hY3Rpb24gZXZlbnRzIHRoYXQgdGhlIFVJIGVtaXRzIEJBQ0sgdG8gdGhlIGhvc3QuICovXG4gICAgVXNlcjogT2JqZWN0LmZyZWV6ZSh7XG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgSW5zdGFsbCBvbiBhbiBBdmFpbGFibGUgdXBkYXRlLiAqL1xuICAgICAgICBJbnN0YWxsOiBcIndhaWxzOnVwZGF0ZXI6dXNlcjppbnN0YWxsXCIsXG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgUmVzdGFydCAmIEFwcGx5IG9uIGEgUmVhZHkgdXBkYXRlLiAqL1xuICAgICAgICBSZXN0YXJ0OiBcIndhaWxzOnVwZGF0ZXI6dXNlcjpyZXN0YXJ0XCIsXG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgU2tpcCBUaGlzIFZlcnNpb24uICovXG4gICAgICAgIFNraXA6IFwid2FpbHM6dXBkYXRlcjp1c2VyOnNraXBcIixcbiAgICAgICAgLyoqIFVzZXIgY2xpY2tlZCBSZW1pbmQgTWUgTGF0ZXIuICovXG4gICAgICAgIFJlbWluZDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6cmVtaW5kXCIsXG4gICAgICAgIC8qKiBVc2VyIGNsaWNrZWQgQ2xvc2UgLyBDYW5jZWwuICovXG4gICAgICAgIENhbmNlbDogXCJ3YWlsczp1cGRhdGVyOnVzZXI6Y2FuY2VsXCIsXG4gICAgfSksXG5cbiAgICAvKiogU3ViLW5hbWVzcGFjZTogZnJhbWV3b3JrLWludGVybmFsIGV2ZW50cyB0aGUgVUkgZW1pdHMgdG8gY29vcmRpbmF0ZVxuICAgICAqICB3aXRoIHRoZSBob3N0LiBNb3N0IGFwcCBjb2RlIGNhbiBpZ25vcmUgdGhlc2UuICovXG4gICAgV2luZG93OiBPYmplY3QuZnJlZXplKHtcbiAgICAgICAgLyoqIFRoZSB3aW5kb3cgZmluaXNoZWQgbG9hZGluZyBhbmQgYXNrcyB0aGUgaG9zdCB0byByZXBsYXkgY3VycmVudCBzdGF0ZS4gKi9cbiAgICAgICAgUmVhZHk6IFwid2FpbHM6dXBkYXRlcjp3aW5kb3c6cmVhZHlcIixcbiAgICB9KSxcbn0pO1xuIl0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FB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cUJBQXFCO0FBQUEsSUFDckIscUJBQXFCO0FBQUEsSUFDckIsdUJBQXVCO0FBQUEsSUFDdkIsY0FBYztBQUFBLElBQ2QsZUFBZTtBQUFBLElBQ2Ysb0JBQW9CO0FBQUEsSUFDcEIsb0JBQW9CO0FBQUEsSUFDcEIsMEJBQTBCO0FBQUEsSUFDMUIsZ0JBQWdCO0FBQUEsSUFDaEIsNEJBQTRCO0FBQUEsSUFDNUIsNEJBQTRCO0FBQUEsSUFDNUIseURBQXlEO0FBQUEsSUFDekQsc0NBQXNDO0FBQUEsSUFDdEMsb0JBQW9CO0FBQUEsSUFDcEIscUJBQXFCO0FBQUEsSUFDckIscUJBQXFCO0FBQUEsSUFDckIsc0JBQXNCO0FBQUEsSUFDdEIsZ0NBQWdDO0FBQUEsSUFDaEMsa0NBQWtDO0FBQUEsSUFDbEMsbUNBQW1DO0FBQUEsSUFDbkMsb0NBQW9DO0FBQUEsSUFDcEMsK0JBQStCO0FBQUEsSUFDL0IsNkJBQTZCO0FBQUEsSUFDN0IsdUJBQXVCO0FBQUEsSUFDdkIsaUNBQWlDO0FBQUEsSUFDakMsOEJBQThCO0FBQUEsSUFDOUIsNEJBQTRCO0FBQUEsSUFDNUIsc0NBQXNDO0FBQUEsSUFDdEMsNEJBQTRCO0FBQUEsSUFDNUIsc0JBQXNCO0FBQUEsSUFDdEIsa0NBQWtDO0FBQUEsSUFDbEMsc0JBQXNCO0FBQUEsSUFDdEIsd0JBQXdCO0FBQUEsSUFDeEIsd0JBQXdCO0FBQUEsSUFDeEIsbUJBQW1CO0FBQUEsSUFDbkIsMEJBQTBCO0FBQUEsSUFDMUIsOEJBQThCO0FBQUEsSUFDOUIseUJBQXlCO0FBQUEsSUFDekIsNkJBQTZCO0FBQUEsSUFDN0IsaUJBQWlCO0FBQUEsSUFDakIsZ0JBQWdCO0FBQUEsSUFDaEIsc0JBQXNCO0FBQUEsSUFDdEIsZUFBZTtBQUFBLElBQ2YseUJBQXlCO0FBQUEsSUFDekIsd0JBQXdCO0FBQUEsSUFDeEIsb0JBQW9CO0FBQUEsSUFDcEIscUJBQXFCO0FBQUEsSUFDckIsaUJBQWlCO0FBQUEsSUFDakIsaUJBQWlCO0FBQUEsSUFDakIsc0JBQXNCO0FBQUEsSUFDdEIsbUNBQW1DO0FB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QUNMLENBQUM7OztBM0IvREQsSUFBSSxRQUFRO0FBQ1IsU0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBQ3RDO0FBNkRBLElBQUksUUFBUTtBQUNSLFNBQU8sT0FBTyxTQUFnQjtBQUM5QixTQUFPLE9BQU8sV0FBVztBQUM3QjtBQUtBLElBQUksUUFBUTtBQUNSLFNBQU8sT0FBTyx5QkFBeUIsZUFBTyx1QkFBdUIsS0FBSyxjQUFNO0FBQ3BGO0FBR0EsSUFBSSxRQUFRO0FBQ1IsU0FBTyxPQUFPLGtCQUFrQjtBQUNoQyxTQUFPLE9BQU8sa0JBQWtCO0FBQ2hDLFNBQU8sT0FBTyxpQkFBaUI7QUFDbkM7QUFFQSxJQUFJLFFBQVE7QUFDUixFQUFPLE9BQU8scUJBQXFCO0FBQ3ZDO0FBT08sU0FBUyxtQkFBbUIsS0FBNEI7QUFDM0QsU0FBTyxNQUFNLEtBQUssRUFBRSxRQUFRLE9BQU8sQ0FBQyxFQUMvQixLQUFLLGNBQVk7QUFDZCxRQUFJLFNBQVMsSUFBSTtBQUdiLFlBQU0sZUFBZSxTQUFTLFFBQVEsSUFBSSxjQUFjLEtBQUssSUFBSSxZQUFZO0FBQzdFLFVBQUksWUFBWSxTQUFTLFlBQVksR0FBRztBQUNwQyxjQUFNLFNBQVMsU0FBUyxjQUFjLFFBQVE7QUFDOUMsZUFBTyxNQUFNO0FBQ2IsaUJBQVMsS0FBSyxZQUFZLE1BQU07QUFBQSxNQUNwQztBQUFBLElBQ0o7QUFBQSxFQUNKLENBQUMsRUFDQSxNQUFNLE1BQU07QUFBQSxFQUFDLENBQUM7QUFDdkI7QUFHQSxJQUFJLFFBQVE7QUFDUixxQkFBbUIsa0JBQWtCO0FBQ3pDOyIsCiAgIm5hbWVzIjogWyJfYSIsICJFcnJvciIsICJjYWxsIiwgIkVycm9yIiwgIl9hIiwgIkFycmF5IiwgIk1hcCIsICJBcnJheSIsICJNYXAiLCAia2V5IiwgImNhbGwiLCAiX2EiLCAiX2EiLCAicmVzaXphYmxlIiwgIl9hIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJfYSIsICJfYSIsICJjYWxsIiwgIkhpZGVNZXRob2QiLCAiU2hvd01ldGhvZCIsICJpc0RvY3VtZW50RG90QWxsIiwgIl9hIiwgInJlYXNvbiIsICJ2YWx1ZSIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgIkhhcHRpY3MiLCAiRGV2aWNlIiwgIkluZm8iLCAiRGV2aWNlIiwgIkhhcHRpY3MiLCAiY2FsbCIsICJEZXZpY2VJbmZvIiwgIkhhcHRpY3MiLCAiRGV2aWNlIiwgIkluZm8iLCAiVG9hc3QiLCAiU2hvdyIsICJFdmVudHMiLCAiRXZlbnRzIl0KfQo=
