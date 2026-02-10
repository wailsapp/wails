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
function hasAndroidBridge() {
  var _a2;
  return typeof window !== "undefined" && typeof ((_a2 = window.wails) == null ? void 0 : _a2.invoke) === "function";
}
function parseAndroidInvokeResponse(responseText) {
  if (!responseText) {
    return null;
  }
  try {
    const parsed = JSON.parse(responseText);
    if (parsed && typeof parsed === "object" && "ok" in parsed) {
      if (parsed.ok === false) {
        throw new Error(parsed.error || "runtime call failed");
      }
      return parsed.data;
    }
    return parsed;
  } catch (err) {
    if (err instanceof Error) {
      throw err;
    }
    return responseText;
  }
}
function configureAndroidTransport() {
  if (customTransport || !hasAndroidBridge()) {
    return;
  }
  customTransport = {
    call: async (objectID, method, windowName, args) => {
      const payload = {
        type: "runtime",
        object: objectID,
        method,
        clientId
      };
      if (windowName) {
        payload.windowName = windowName;
      }
      if (args !== null && args !== void 0) {
        payload.args = args;
      }
      const responseText = window.wails.invoke(JSON.stringify(payload));
      return parseAndroidInvokeResponse(responseText);
    }
  };
}
configureAndroidTransport();
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
  let response = await fetch(url, {
    method: "POST",
    headers,
    body: JSON.stringify(body)
  });
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
  var _a2;
  const urlString = url.toString();
  const androidOpenURL = (_a2 = window == null ? void 0 : window.wails) == null ? void 0 : _a2.openURL;
  if (typeof androidOpenURL === "function") {
    try {
      androidOpenURL.call(window.wails, urlString);
      return Promise.resolve();
    } catch (e) {
      return Promise.reject(e);
    }
  }
  return call(BrowserOpenURL, { url: urlString });
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

// desktop/@wailsio/runtime/src/android.ts
var android_exports = {};
__export(android_exports, {
  Device: () => Device,
  Haptics: () => Haptics,
  Links: () => Links,
  Navigation: () => Navigation,
  Scroll: () => Scroll,
  Toast: () => Toast,
  UserAgent: () => UserAgent
});
var call10 = newRuntimeCaller(objectNames.Android);
var HapticsVibrate = 0;
var DeviceInfo = 1;
var ToastShow = 2;
var ScrollSetEnabled = 3;
var ScrollSetBounceEnabled = 4;
var ScrollSetIndicatorsEnabled = 5;
var NavigationSetBackForwardGestures = 6;
var LinksSetPreviewEnabled = 7;
var UserAgentSet = 8;
var Haptics;
((Haptics3) => {
  function Vibrate(duration = 100) {
    return call10(HapticsVibrate, { duration });
  }
  Haptics3.Vibrate = Vibrate;
})(Haptics || (Haptics = {}));
var Device;
((Device3) => {
  function Info2() {
    return call10(DeviceInfo);
  }
  Device3.Info = Info2;
})(Device || (Device = {}));
var Toast;
((Toast2) => {
  function Show2(message) {
    return call10(ToastShow, { message });
  }
  Toast2.Show = Show2;
})(Toast || (Toast = {}));
var Scroll;
((Scroll2) => {
  function SetEnabled(enabled) {
    return call10(ScrollSetEnabled, { enabled });
  }
  Scroll2.SetEnabled = SetEnabled;
  function SetBounceEnabled(enabled) {
    return call10(ScrollSetBounceEnabled, { enabled });
  }
  Scroll2.SetBounceEnabled = SetBounceEnabled;
  function SetIndicatorsEnabled(enabled) {
    return call10(ScrollSetIndicatorsEnabled, { enabled });
  }
  Scroll2.SetIndicatorsEnabled = SetIndicatorsEnabled;
})(Scroll || (Scroll = {}));
var Navigation;
((Navigation2) => {
  function SetBackForwardGesturesEnabled(enabled) {
    return call10(NavigationSetBackForwardGestures, { enabled });
  }
  Navigation2.SetBackForwardGesturesEnabled = SetBackForwardGesturesEnabled;
})(Navigation || (Navigation = {}));
var Links;
((Links2) => {
  function SetPreviewEnabled(enabled) {
    return call10(LinksSetPreviewEnabled, { enabled });
  }
  Links2.SetPreviewEnabled = SetPreviewEnabled;
})(Links || (Links = {}));
var UserAgent;
((UserAgent2) => {
  function Set(ua) {
    return call10(UserAgentSet, { ua });
  }
  UserAgent2.Set = Set;
})(UserAgent || (UserAgent = {}));

// desktop/@wailsio/runtime/src/ios.ts
var ios_exports = {};
__export(ios_exports, {
  Device: () => Device2,
  Haptics: () => Haptics2
});
var call11 = newRuntimeCaller(objectNames.IOS);
var HapticsImpact = 0;
var DeviceInfo2 = 1;
var Haptics2;
((Haptics3) => {
  function Impact(style = "medium") {
    return call11(HapticsImpact, { style });
  }
  Haptics3.Impact = Impact;
})(Haptics2 || (Haptics2 = {}));
var Device2;
((Device3) => {
  function Info2() {
    return call11(DeviceInfo2);
  }
  Device3.Info = Info2;
})(Device2 || (Device2 = {}));

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
      const script = document.createElement("script");
      script.src = url;
      document.head.appendChild(script);
    }
  }).catch(() => {
  });
}
loadOptionalScript("/wails/custom.js");
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
  wml_exports as WML,
  window_default as Window,
  clientId,
  getTransport,
  loadOptionalScript,
  objectNames,
  setTransport
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY3JlYXRlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudF90eXBlcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbGFibGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbmNlbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FuZHJvaWQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2lvcy50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLy8gU2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG5pbXBvcnQgXCIuL2NvbnRleHRtZW51LmpzXCI7XHJcbmltcG9ydCBcIi4vZHJhZy5qc1wiO1xyXG5cclxuLy8gUmUtZXhwb3J0IHB1YmxpYyBBUElcclxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSBcIi4vYXBwbGljYXRpb24uanNcIjtcclxuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyLmpzXCI7XHJcbmltcG9ydCAqIGFzIENhbGwgZnJvbSBcIi4vY2FsbHMuanNcIjtcclxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gXCIuL2NsaXBib2FyZC5qc1wiO1xyXG5pbXBvcnQgKiBhcyBDcmVhdGUgZnJvbSBcIi4vY3JlYXRlLmpzXCI7XHJcbmltcG9ydCAqIGFzIERpYWxvZ3MgZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBFdmVudHMgZnJvbSBcIi4vZXZlbnRzLmpzXCI7XHJcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuL2ZsYWdzLmpzXCI7XHJcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xyXG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSBcIi4vc3lzdGVtLmpzXCI7XHJcbmltcG9ydCAqIGFzIEFuZHJvaWQgZnJvbSBcIi4vYW5kcm9pZC5qc1wiO1xyXG5pbXBvcnQgKiBhcyBJT1MgZnJvbSBcIi4vaW9zLmpzXCI7XHJcbmltcG9ydCBXaW5kb3csIHsgaGFuZGxlRHJhZ0VudGVyLCBoYW5kbGVEcmFnTGVhdmUsIGhhbmRsZURyYWdPdmVyIH0gZnJvbSBcIi4vd2luZG93LmpzXCI7XHJcbmltcG9ydCAqIGFzIFdNTCBmcm9tIFwiLi93bWwuanNcIjtcclxuXHJcbmV4cG9ydCB7XHJcbiAgICBBcHBsaWNhdGlvbixcclxuICAgIEJyb3dzZXIsXHJcbiAgICBDYWxsLFxyXG4gICAgQ2xpcGJvYXJkLFxyXG4gICAgRGlhbG9ncyxcclxuICAgIEV2ZW50cyxcclxuICAgIEZsYWdzLFxyXG4gICAgU2NyZWVucyxcclxuICAgIFN5c3RlbSxcclxuICAgIEFuZHJvaWQsXHJcbiAgICBJT1MsXHJcbiAgICBXaW5kb3csXHJcbiAgICBXTUxcclxufTtcclxuXHJcbi8qKlxyXG4gKiBBbiBpbnRlcm5hbCB1dGlsaXR5IGNvbnN1bWVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cclxuICpcclxuICogQGlnbm9yZVxyXG4gKi9cclxuZXhwb3J0IHsgQ3JlYXRlIH07XHJcblxyXG5leHBvcnQgKiBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xyXG5cclxuLy8gRXhwb3J0IHRyYW5zcG9ydCBpbnRlcmZhY2VzIGFuZCB1dGlsaXRpZXNcclxuZXhwb3J0IHtcclxuICAgIHNldFRyYW5zcG9ydCxcclxuICAgIGdldFRyYW5zcG9ydCxcclxuICAgIHR5cGUgUnVudGltZVRyYW5zcG9ydCxcclxuICAgIG9iamVjdE5hbWVzLFxyXG4gICAgY2xpZW50SWQsXHJcbn0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuaW1wb3J0IHsgY2xpZW50SWQgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG4vLyBOb3RpZnkgYmFja2VuZFxyXG53aW5kb3cuX3dhaWxzLmludm9rZSA9IFN5c3RlbS5pbnZva2U7XHJcbndpbmRvdy5fd2FpbHMuY2xpZW50SWQgPSBjbGllbnRJZDtcclxuXHJcbi8vIFJlZ2lzdGVyIHBsYXRmb3JtIGhhbmRsZXJzIChpbnRlcm5hbCBBUEkpXHJcbi8vIE5vdGU6IFdpbmRvdyBpcyB0aGUgdGhpc1dpbmRvdyBpbnN0YW5jZSAoZGVmYXVsdCBleHBvcnQgZnJvbSB3aW5kb3cudHMpXHJcbi8vIEJpbmRpbmcgZW5zdXJlcyAndGhpcycgY29ycmVjdGx5IHJlZmVycyB0byB0aGUgY3VycmVudCB3aW5kb3cgaW5zdGFuY2Vcclxud2luZG93Ll93YWlscy5oYW5kbGVQbGF0Zm9ybUZpbGVEcm9wID0gV2luZG93LkhhbmRsZVBsYXRmb3JtRmlsZURyb3AuYmluZChXaW5kb3cpO1xyXG5cclxuLy8gTGludXgtc3BlY2lmaWMgZHJhZyBoYW5kbGVycyAoR1RLIGludGVyY2VwdHMgRE9NIGRyYWcgZXZlbnRzKVxyXG53aW5kb3cuX3dhaWxzLmhhbmRsZURyYWdFbnRlciA9IGhhbmRsZURyYWdFbnRlcjtcclxud2luZG93Ll93YWlscy5oYW5kbGVEcmFnTGVhdmUgPSBoYW5kbGVEcmFnTGVhdmU7XHJcbndpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ092ZXIgPSBoYW5kbGVEcmFnT3ZlcjtcclxuXHJcblN5c3RlbS5pbnZva2UoXCJ3YWlsczpydW50aW1lOnJlYWR5XCIpO1xyXG5cclxuLyoqXHJcbiAqIExvYWRzIGEgc2NyaXB0IGZyb20gdGhlIGdpdmVuIFVSTCBpZiBpdCBleGlzdHMuXHJcbiAqIFVzZXMgSEVBRCByZXF1ZXN0IHRvIGNoZWNrIGV4aXN0ZW5jZSwgdGhlbiBpbmplY3RzIGEgc2NyaXB0IHRhZy5cclxuICogU2lsZW50bHkgaWdub3JlcyBpZiB0aGUgc2NyaXB0IGRvZXNuJ3QgZXhpc3QuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gbG9hZE9wdGlvbmFsU2NyaXB0KHVybDogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gZmV0Y2godXJsLCB7IG1ldGhvZDogJ0hFQUQnIH0pXHJcbiAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcclxuICAgICAgICAgICAgICAgIGNvbnN0IHNjcmlwdCA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoJ3NjcmlwdCcpO1xyXG4gICAgICAgICAgICAgICAgc2NyaXB0LnNyYyA9IHVybDtcclxuICAgICAgICAgICAgICAgIGRvY3VtZW50LmhlYWQuYXBwZW5kQ2hpbGQoc2NyaXB0KTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pXHJcbiAgICAgICAgLmNhdGNoKCgpID0+IHt9KTsgLy8gU2lsZW50bHkgaWdub3JlIC0gc2NyaXB0IGlzIG9wdGlvbmFsXHJcbn1cclxuXHJcbi8vIExvYWQgY3VzdG9tLmpzIGlmIGF2YWlsYWJsZSAodXNlZCBieSBzZXJ2ZXIgbW9kZSBmb3IgV2ViU29ja2V0IGV2ZW50cywgZXRjLilcclxubG9hZE9wdGlvbmFsU2NyaXB0KCcvd2FpbHMvY3VzdG9tLmpzJyk7XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBPcGVuVVJMIH0gZnJvbSBcIi4vYnJvd3Nlci5qc1wiO1xyXG5pbXBvcnQgeyBRdWVzdGlvbiB9IGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcclxuaW1wb3J0IHsgRW1pdCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xyXG5pbXBvcnQgeyBjYW5BYm9ydExpc3RlbmVycywgd2hlblJlYWR5IH0gZnJvbSBcIi4vdXRpbHMuanNcIjtcclxuaW1wb3J0IFdpbmRvdyBmcm9tIFwiLi93aW5kb3cuanNcIjtcclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBvcHRpb25hbCBkYXRhLlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gc2VuZC5cclxuICogQHBhcmFtIFtkYXRhPW51bGxdIC0gLSBPcHRpb25hbCBkYXRhIHRvIHNlbmQgYWxvbmcgd2l0aCB0aGUgZXZlbnQuXHJcbiAqL1xyXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lOiBzdHJpbmcsIGRhdGE6IGFueSA9IG51bGwpOiB2b2lkIHtcclxuICAgIEVtaXQoZXZlbnROYW1lLCBkYXRhKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgc3BlY2lmaWVkIHdpbmRvdy5cclxuICpcclxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGNhbGwgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqL1xyXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKHdpbmRvd05hbWU6IHN0cmluZywgbWV0aG9kTmFtZTogc3RyaW5nKSB7XHJcbiAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBXaW5kb3cuR2V0KHdpbmRvd05hbWUpO1xyXG4gICAgY29uc3QgbWV0aG9kID0gKHRhcmdldFdpbmRvdyBhcyBhbnkpW21ldGhvZE5hbWVdO1xyXG5cclxuICAgIGlmICh0eXBlb2YgbWV0aG9kICE9PSBcImZ1bmN0aW9uXCIpIHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKGBXaW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JyBub3QgZm91bmRgKTtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgdHJ5IHtcclxuICAgICAgICBtZXRob2QuY2FsbCh0YXJnZXRXaW5kb3cpO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIGNvbnNvbGUuZXJyb3IoYEVycm9yIGNhbGxpbmcgd2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfSc6IGAsIGUpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmVzcG9uZHMgdG8gYSB0cmlnZ2VyaW5nIGV2ZW50IGJ5IHJ1bm5pbmcgYXBwcm9wcmlhdGUgV01MIGFjdGlvbnMgZm9yIHRoZSBjdXJyZW50IHRhcmdldC5cclxuICovXHJcbmZ1bmN0aW9uIG9uV01MVHJpZ2dlcmVkKGV2OiBFdmVudCk6IHZvaWQge1xyXG4gICAgY29uc3QgZWxlbWVudCA9IGV2LmN1cnJlbnRUYXJnZXQgYXMgRWxlbWVudDtcclxuXHJcbiAgICBmdW5jdGlvbiBydW5FZmZlY3QoY2hvaWNlID0gXCJZZXNcIikge1xyXG4gICAgICAgIGlmIChjaG9pY2UgIT09IFwiWWVzXCIpXHJcbiAgICAgICAgICAgIHJldHVybjtcclxuXHJcbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1ldmVudCcpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1ldmVudCcpO1xyXG4gICAgICAgIGNvbnN0IHRhcmdldFdpbmRvdyA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdGFyZ2V0LXdpbmRvdycpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10YXJnZXQtd2luZG93JykgfHwgXCJcIjtcclxuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC13aW5kb3cnKTtcclxuICAgICAgICBjb25zdCB1cmwgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLW9wZW51cmwnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtb3BlbnVybCcpO1xyXG5cclxuICAgICAgICBpZiAoZXZlbnRUeXBlICE9PSBudWxsKVxyXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICBpZiAod2luZG93TWV0aG9kICE9PSBudWxsKVxyXG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHRhcmdldFdpbmRvdywgd2luZG93TWV0aG9kKTtcclxuICAgICAgICBpZiAodXJsICE9PSBudWxsKVxyXG4gICAgICAgICAgICB2b2lkIE9wZW5VUkwodXJsKTtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcclxuXHJcbiAgICBpZiAoY29uZmlybSkge1xyXG4gICAgICAgIFF1ZXN0aW9uKHtcclxuICAgICAgICAgICAgVGl0bGU6IFwiQ29uZmlybVwiLFxyXG4gICAgICAgICAgICBNZXNzYWdlOiBjb25maXJtLFxyXG4gICAgICAgICAgICBEZXRhY2hlZDogZmFsc2UsXHJcbiAgICAgICAgICAgIEJ1dHRvbnM6IFtcclxuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiWWVzXCIgfSxcclxuICAgICAgICAgICAgICAgIHsgTGFiZWw6IFwiTm9cIiwgSXNEZWZhdWx0OiB0cnVlIH1cclxuICAgICAgICAgICAgXVxyXG4gICAgICAgIH0pLnRoZW4ocnVuRWZmZWN0KTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgcnVuRWZmZWN0KCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXHJcbmNvbnN0IGNvbnRyb2xsZXJTeW0gPSBTeW1ib2woXCJjb250cm9sbGVyXCIpO1xyXG5jb25zdCB0cmlnZ2VyTWFwU3ltID0gU3ltYm9sKFwidHJpZ2dlck1hcFwiKTtcclxuY29uc3QgZWxlbWVudENvdW50U3ltID0gU3ltYm9sKFwiZWxlbWVudENvdW50XCIpO1xyXG5cclxuLyoqXHJcbiAqIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IGRvZXMgbm90IGFjdHVhbGx5IHJlbWVtYmVyIGFjdGl2ZSBldmVudCBsaXN0ZW5lcnM6IGluc3RlYWRcclxuICogaXQgdGllcyB0aGVtIHRvIGFuIEFib3J0U2lnbmFsIGFuZCB1c2VzIGFuIEFib3J0Q29udHJvbGxlciB0byByZW1vdmUgdGhlbSBhbGwgYXQgb25jZS5cclxuICovXHJcbmNsYXNzIEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5IHtcclxuICAgIC8vIFByaXZhdGUgZmllbGRzLlxyXG4gICAgW2NvbnRyb2xsZXJTeW1dOiBBYm9ydENvbnRyb2xsZXI7XHJcblxyXG4gICAgY29uc3RydWN0b3IoKSB7XHJcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgYW4gb3B0aW9ucyBvYmplY3QgZm9yIGFkZEV2ZW50TGlzdGVuZXIgdGhhdCB0aWVzIHRoZSBsaXN0ZW5lclxyXG4gICAgICogdG8gdGhlIEFib3J0U2lnbmFsIGZyb20gdGhlIGN1cnJlbnQgQWJvcnRDb250cm9sbGVyLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XHJcbiAgICAgKiBAcGFyYW0gdHJpZ2dlcnMgLSBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRzXHJcbiAgICAgKi9cclxuICAgIHNldChlbGVtZW50OiBFbGVtZW50LCB0cmlnZ2Vyczogc3RyaW5nW10pOiBBZGRFdmVudExpc3RlbmVyT3B0aW9ucyB7XHJcbiAgICAgICAgcmV0dXJuIHsgc2lnbmFsOiB0aGlzW2NvbnRyb2xsZXJTeW1dLnNpZ25hbCB9O1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMgYW5kIHJlc2V0cyB0aGUgcmVnaXN0cnkuXHJcbiAgICAgKi9cclxuICAgIHJlc2V0KCk6IHZvaWQge1xyXG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0uYWJvcnQoKTtcclxuICAgICAgICB0aGlzW2NvbnRyb2xsZXJTeW1dID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogV2Vha01hcFJlZ2lzdHJ5IG1hcHMgYWN0aXZlIHRyaWdnZXIgZXZlbnRzIHRvIGVhY2ggRE9NIGVsZW1lbnQgdGhyb3VnaCBhIFdlYWtNYXAuXHJcbiAqIFRoaXMgZW5zdXJlcyB0aGF0IHRoZSBtYXBwaW5nIHJlbWFpbnMgcHJpdmF0ZSB0byB0aGlzIG1vZHVsZSwgd2hpbGUgc3RpbGwgYWxsb3dpbmcgZ2FyYmFnZVxyXG4gKiBjb2xsZWN0aW9uIG9mIHRoZSBpbnZvbHZlZCBlbGVtZW50cy5cclxuICovXHJcbmNsYXNzIFdlYWtNYXBSZWdpc3RyeSB7XHJcbiAgICAvKiogU3RvcmVzIHRoZSBjdXJyZW50IGVsZW1lbnQtdG8tdHJpZ2dlciBtYXBwaW5nLiAqL1xyXG4gICAgW3RyaWdnZXJNYXBTeW1dOiBXZWFrTWFwPEVsZW1lbnQsIHN0cmluZ1tdPjtcclxuICAgIC8qKiBDb3VudHMgdGhlIG51bWJlciBvZiBlbGVtZW50cyB3aXRoIGFjdGl2ZSBXTUwgdHJpZ2dlcnMuICovXHJcbiAgICBbZWxlbWVudENvdW50U3ltXTogbnVtYmVyO1xyXG5cclxuICAgIGNvbnN0cnVjdG9yKCkge1xyXG4gICAgICAgIHRoaXNbdHJpZ2dlck1hcFN5bV0gPSBuZXcgV2Vha01hcCgpO1xyXG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50U3ltXSA9IDA7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIGFjdGl2ZSB0cmlnZ2VycyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XHJcbiAgICAgKiBAcGFyYW0gdHJpZ2dlcnMgLSBUaGUgbGlzdCBvZiBhY3RpdmUgV01MIHRyaWdnZXIgZXZlbnRzIGZvciB0aGUgc3BlY2lmaWVkIGVsZW1lbnRcclxuICAgICAqL1xyXG4gICAgc2V0KGVsZW1lbnQ6IEVsZW1lbnQsIHRyaWdnZXJzOiBzdHJpbmdbXSk6IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zIHtcclxuICAgICAgICBpZiAoIXRoaXNbdHJpZ2dlck1hcFN5bV0uaGFzKGVsZW1lbnQpKSB7IHRoaXNbZWxlbWVudENvdW50U3ltXSsrOyB9XHJcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXS5zZXQoZWxlbWVudCwgdHJpZ2dlcnMpO1xyXG4gICAgICAgIHJldHVybiB7fTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbW92ZXMgYWxsIHJlZ2lzdGVyZWQgZXZlbnQgbGlzdGVuZXJzLlxyXG4gICAgICovXHJcbiAgICByZXNldCgpOiB2b2lkIHtcclxuICAgICAgICBpZiAodGhpc1tlbGVtZW50Q291bnRTeW1dIDw9IDApXHJcbiAgICAgICAgICAgIHJldHVybjtcclxuXHJcbiAgICAgICAgZm9yIChjb25zdCBlbGVtZW50IG9mIGRvY3VtZW50LmJvZHkucXVlcnlTZWxlY3RvckFsbCgnKicpKSB7XHJcbiAgICAgICAgICAgIGlmICh0aGlzW2VsZW1lbnRDb3VudFN5bV0gPD0gMClcclxuICAgICAgICAgICAgICAgIGJyZWFrO1xyXG5cclxuICAgICAgICAgICAgY29uc3QgdHJpZ2dlcnMgPSB0aGlzW3RyaWdnZXJNYXBTeW1dLmdldChlbGVtZW50KTtcclxuICAgICAgICAgICAgaWYgKHRyaWdnZXJzICE9IG51bGwpIHsgdGhpc1tlbGVtZW50Q291bnRTeW1dLS07IH1cclxuXHJcbiAgICAgICAgICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycyB8fCBbXSlcclxuICAgICAgICAgICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcclxuICAgICAgICB0aGlzW2VsZW1lbnRDb3VudFN5bV0gPSAwO1xyXG4gICAgfVxyXG59XHJcblxyXG5jb25zdCB0cmlnZ2VyUmVnaXN0cnkgPSBjYW5BYm9ydExpc3RlbmVycygpID8gbmV3IEFib3J0Q29udHJvbGxlclJlZ2lzdHJ5KCkgOiBuZXcgV2Vha01hcFJlZ2lzdHJ5KCk7XHJcblxyXG4vKipcclxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gdGhlIHNwZWNpZmllZCBlbGVtZW50LlxyXG4gKi9cclxuZnVuY3Rpb24gYWRkV01MTGlzdGVuZXJzKGVsZW1lbnQ6IEVsZW1lbnQpOiB2b2lkIHtcclxuICAgIGNvbnN0IHRyaWdnZXJSZWdFeHAgPSAvXFxTKy9nO1xyXG4gICAgY29uc3QgdHJpZ2dlckF0dHIgPSAoZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCIpO1xyXG4gICAgY29uc3QgdHJpZ2dlcnM6IHN0cmluZ1tdID0gW107XHJcblxyXG4gICAgbGV0IG1hdGNoO1xyXG4gICAgd2hpbGUgKChtYXRjaCA9IHRyaWdnZXJSZWdFeHAuZXhlYyh0cmlnZ2VyQXR0cikpICE9PSBudWxsKVxyXG4gICAgICAgIHRyaWdnZXJzLnB1c2gobWF0Y2hbMF0pO1xyXG5cclxuICAgIGNvbnN0IG9wdGlvbnMgPSB0cmlnZ2VyUmVnaXN0cnkuc2V0KGVsZW1lbnQsIHRyaWdnZXJzKTtcclxuICAgIGZvciAoY29uc3QgdHJpZ2dlciBvZiB0cmlnZ2VycylcclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQsIG9wdGlvbnMpO1xyXG59XHJcblxyXG4vKipcclxuICogU2NoZWR1bGVzIGFuIGF1dG9tYXRpYyByZWxvYWQgb2YgV01MIHRvIGJlIHBlcmZvcm1lZCBhcyBzb29uIGFzIHRoZSBkb2N1bWVudCBpcyBmdWxseSBsb2FkZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW5hYmxlKCk6IHZvaWQge1xyXG4gICAgd2hlblJlYWR5KFJlbG9hZCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWxvYWRzIHRoZSBXTUwgcGFnZSBieSBhZGRpbmcgbmVjZXNzYXJ5IGV2ZW50IGxpc3RlbmVycyBhbmQgYnJvd3NlciBsaXN0ZW5lcnMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUmVsb2FkKCk6IHZvaWQge1xyXG4gICAgdHJpZ2dlclJlZ2lzdHJ5LnJlc2V0KCk7XHJcbiAgICBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtZXZlbnRdLCBbd21sLXdpbmRvd10sIFt3bWwtb3BlbnVybF0sIFtkYXRhLXdtbC1ldmVudF0sIFtkYXRhLXdtbC13aW5kb3ddLCBbZGF0YS13bWwtb3BlbnVybF0nKS5mb3JFYWNoKGFkZFdNTExpc3RlbmVycyk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQnJvd3Nlcik7XHJcblxyXG5jb25zdCBCcm93c2VyT3BlblVSTCA9IDA7XHJcblxyXG4vKipcclxuICogT3BlbiBhIGJyb3dzZXIgd2luZG93IHRvIHRoZSBnaXZlbiBVUkwuXHJcbiAqXHJcbiAqIEBwYXJhbSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybDogc3RyaW5nIHwgVVJMKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICBjb25zdCB1cmxTdHJpbmcgPSB1cmwudG9TdHJpbmcoKTtcclxuICAgIGNvbnN0IGFuZHJvaWRPcGVuVVJMID0gKHdpbmRvdyBhcyBhbnkpPy53YWlscz8ub3BlblVSTDtcclxuICAgIGlmICh0eXBlb2YgYW5kcm9pZE9wZW5VUkwgPT09IFwiZnVuY3Rpb25cIikge1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIGFuZHJvaWRPcGVuVVJMLmNhbGwoKHdpbmRvdyBhcyBhbnkpLndhaWxzLCB1cmxTdHJpbmcpO1xyXG4gICAgICAgICAgICByZXR1cm4gUHJvbWlzZS5yZXNvbHZlKCk7XHJcbiAgICAgICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgICAgICByZXR1cm4gUHJvbWlzZS5yZWplY3QoZSk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmw6IHVybFN0cmluZ30pO1xyXG59XHJcbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9haS9uYW5vaWRcclxuXHJcbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxyXG4vL1xyXG4vLyBDb3B5cmlnaHQgMjAxNyBBbmRyZXkgU2l0bmlrIDxhbmRyZXlAc2l0bmlrLnJ1PlxyXG4vL1xyXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mXHJcbi8vIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlIFwiU29mdHdhcmVcIiksIHRvIGRlYWwgaW5cclxuLy8gdGhlIFNvZnR3YXJlIHdpdGhvdXQgcmVzdHJpY3Rpb24sIGluY2x1ZGluZyB3aXRob3V0IGxpbWl0YXRpb24gdGhlIHJpZ2h0cyB0b1xyXG4vLyB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsIGNvcGllcyBvZlxyXG4vLyB0aGUgU29mdHdhcmUsIGFuZCB0byBwZXJtaXQgcGVyc29ucyB0byB3aG9tIHRoZSBTb2Z0d2FyZSBpcyBmdXJuaXNoZWQgdG8gZG8gc28sXHJcbi8vICAgICBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcclxuLy9cclxuLy8gICAgIFRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbFxyXG4vLyBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLlxyXG4vL1xyXG4vLyAgICAgVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEIFwiQVMgSVNcIiwgV0lUSE9VVCBXQVJSQU5UWSBPRiBBTlkgS0lORCwgRVhQUkVTUyBPUlxyXG4vLyBJTVBMSUVELCBJTkNMVURJTkcgQlVUIE5PVCBMSU1JVEVEIFRPIFRIRSBXQVJSQU5USUVTIE9GIE1FUkNIQU5UQUJJTElUWSwgRklUTkVTU1xyXG4vLyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1JcclxuLy8gQ09QWVJJR0hUIEhPTERFUlMgQkUgTElBQkxFIEZPUiBBTlkgQ0xBSU0sIERBTUFHRVMgT1IgT1RIRVIgTElBQklMSVRZLCBXSEVUSEVSXHJcbi8vIElOIEFOIEFDVElPTiBPRiBDT05UUkFDVCwgVE9SVCBPUiBPVEhFUldJU0UsIEFSSVNJTkcgRlJPTSwgT1VUIE9GIE9SIElOXHJcbi8vIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuXHJcblxyXG4vLyBUaGlzIGFscGhhYmV0IHVzZXMgYEEtWmEtejAtOV8tYCBzeW1ib2xzLlxyXG4vLyBUaGUgb3JkZXIgb2YgY2hhcmFjdGVycyBpcyBvcHRpbWl6ZWQgZm9yIGJldHRlciBnemlwIGFuZCBicm90bGkgY29tcHJlc3Npb24uXHJcbi8vIFJlZmVyZW5jZXMgdG8gdGhlIHNhbWUgZmlsZSAod29ya3MgYm90aCBmb3IgZ3ppcCBhbmQgYnJvdGxpKTpcclxuLy8gYCd1c2VgLCBgYW5kb21gLCBhbmQgYHJpY3QnYFxyXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBicm90bGkgZGVmYXVsdCBkaWN0aW9uYXJ5OlxyXG4vLyBgLTI2VGAsIGAxOTgzYCwgYDQwcHhgLCBgNzVweGAsIGBidXNoYCwgYGphY2tgLCBgbWluZGAsIGB2ZXJ5YCwgYW5kIGB3b2xmYFxyXG5jb25zdCB1cmxBbHBoYWJldCA9XHJcbiAgICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBuYW5vaWQoc2l6ZTogbnVtYmVyID0gMjEpOiBzdHJpbmcge1xyXG4gICAgbGV0IGlkID0gJydcclxuICAgIC8vIEEgY29tcGFjdCBhbHRlcm5hdGl2ZSBmb3IgYGZvciAodmFyIGkgPSAwOyBpIDwgc3RlcDsgaSsrKWAuXHJcbiAgICBsZXQgaSA9IHNpemUgfCAwXHJcbiAgICB3aGlsZSAoaS0tKSB7XHJcbiAgICAgICAgLy8gYHwgMGAgaXMgbW9yZSBjb21wYWN0IGFuZCBmYXN0ZXIgdGhhbiBgTWF0aC5mbG9vcigpYC5cclxuICAgICAgICBpZCArPSB1cmxBbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIDY0KSB8IDBdXHJcbiAgICB9XHJcbiAgICByZXR1cm4gaWRcclxufVxyXG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XHJcblxyXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcclxuXHJcbi8vIFJlLWV4cG9ydCBuYW5vaWQgZm9yIGN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb25zXHJcbmV4cG9ydCB7IG5hbm9pZCB9O1xyXG5cclxuLy8gT2JqZWN0IE5hbWVzXHJcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xyXG4gICAgQ2FsbDogMCxcclxuICAgIENsaXBib2FyZDogMSxcclxuICAgIEFwcGxpY2F0aW9uOiAyLFxyXG4gICAgRXZlbnRzOiAzLFxyXG4gICAgQ29udGV4dE1lbnU6IDQsXHJcbiAgICBEaWFsb2c6IDUsXHJcbiAgICBXaW5kb3c6IDYsXHJcbiAgICBTY3JlZW5zOiA3LFxyXG4gICAgU3lzdGVtOiA4LFxyXG4gICAgQnJvd3NlcjogOSxcclxuICAgIENhbmNlbENhbGw6IDEwLFxyXG4gICAgSU9TOiAxMSxcclxuICAgIEFuZHJvaWQ6IDEyLFxyXG59KTtcclxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xyXG5cclxuLyoqXHJcbiAqIFJ1bnRpbWVUcmFuc3BvcnQgZGVmaW5lcyB0aGUgaW50ZXJmYWNlIGZvciBjdXN0b20gSVBDIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnMuXHJcbiAqIEltcGxlbWVudCB0aGlzIGludGVyZmFjZSB0byB1c2UgV2ViU29ja2V0cywgY3VzdG9tIHByb3RvY29scywgb3IgYW55IG90aGVyXHJcbiAqIHRyYW5zcG9ydCBtZWNoYW5pc20gaW5zdGVhZCBvZiB0aGUgZGVmYXVsdCBIVFRQIGZldGNoLlxyXG4gKi9cclxuZXhwb3J0IGludGVyZmFjZSBSdW50aW1lVHJhbnNwb3J0IHtcclxuICAgIC8qKlxyXG4gICAgICogU2VuZCBhIHJ1bnRpbWUgY2FsbCBhbmQgcmV0dXJuIHRoZSByZXNwb25zZS5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gb2JqZWN0SUQgLSBUaGUgV2FpbHMgb2JqZWN0IElEICgwPUNhbGwsIDE9Q2xpcGJvYXJkLCBldGMuKVxyXG4gICAgICogQHBhcmFtIG1ldGhvZCAtIFRoZSBtZXRob2QgSUQgdG8gY2FsbFxyXG4gICAgICogQHBhcmFtIHdpbmRvd05hbWUgLSBPcHRpb25hbCB3aW5kb3cgbmFtZVxyXG4gICAgICogQHBhcmFtIGFyZ3MgLSBBcmd1bWVudHMgdG8gcGFzcyAod2lsbCBiZSBKU09OIHN0cmluZ2lmaWVkIGlmIHByZXNlbnQpXHJcbiAgICAgKiBAcmV0dXJucyBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzcG9uc2UgZGF0YVxyXG4gICAgICovXHJcbiAgICBjYWxsKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55PjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb24gKGNhbiBiZSBzZXQgYnkgdXNlcilcclxuICovXHJcbmxldCBjdXN0b21UcmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsID0gbnVsbDtcclxuXHJcbmZ1bmN0aW9uIGhhc0FuZHJvaWRCcmlkZ2UoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gdHlwZW9mIHdpbmRvdyAhPT0gXCJ1bmRlZmluZWRcIiAmJiB0eXBlb2YgKHdpbmRvdyBhcyBhbnkpLndhaWxzPy5pbnZva2UgPT09IFwiZnVuY3Rpb25cIjtcclxufVxyXG5cclxuZnVuY3Rpb24gcGFyc2VBbmRyb2lkSW52b2tlUmVzcG9uc2UocmVzcG9uc2VUZXh0OiBzdHJpbmcpOiBhbnkge1xyXG4gICAgaWYgKCFyZXNwb25zZVRleHQpIHtcclxuICAgICAgICByZXR1cm4gbnVsbDtcclxuICAgIH1cclxuXHJcbiAgICB0cnkge1xyXG4gICAgICAgIGNvbnN0IHBhcnNlZCA9IEpTT04ucGFyc2UocmVzcG9uc2VUZXh0KTtcclxuICAgICAgICBpZiAocGFyc2VkICYmIHR5cGVvZiBwYXJzZWQgPT09IFwib2JqZWN0XCIgJiYgXCJva1wiIGluIHBhcnNlZCkge1xyXG4gICAgICAgICAgICBpZiAocGFyc2VkLm9rID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICAgICAgdGhyb3cgbmV3IEVycm9yKHBhcnNlZC5lcnJvciB8fCBcInJ1bnRpbWUgY2FsbCBmYWlsZWRcIik7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgcmV0dXJuIHBhcnNlZC5kYXRhO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gcGFyc2VkO1xyXG4gICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgaWYgKGVyciBpbnN0YW5jZW9mIEVycm9yKSB7XHJcbiAgICAgICAgICAgIHRocm93IGVycjtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHJlc3BvbnNlVGV4dDtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gY29uZmlndXJlQW5kcm9pZFRyYW5zcG9ydCgpOiB2b2lkIHtcclxuICAgIGlmIChjdXN0b21UcmFuc3BvcnQgfHwgIWhhc0FuZHJvaWRCcmlkZ2UoKSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBjdXN0b21UcmFuc3BvcnQgPSB7XHJcbiAgICAgICAgY2FsbDogYXN5bmMgKG9iamVjdElEOiBudW1iZXIsIG1ldGhvZDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcsIGFyZ3M6IGFueSk6IFByb21pc2U8YW55PiA9PiB7XHJcbiAgICAgICAgICAgIGNvbnN0IHBheWxvYWQ6IFJlY29yZDxzdHJpbmcsIGFueT4gPSB7XHJcbiAgICAgICAgICAgICAgICB0eXBlOiBcInJ1bnRpbWVcIixcclxuICAgICAgICAgICAgICAgIG9iamVjdDogb2JqZWN0SUQsXHJcbiAgICAgICAgICAgICAgICBtZXRob2QsXHJcbiAgICAgICAgICAgICAgICBjbGllbnRJZCxcclxuICAgICAgICAgICAgfTtcclxuXHJcbiAgICAgICAgICAgIGlmICh3aW5kb3dOYW1lKSB7XHJcbiAgICAgICAgICAgICAgICBwYXlsb2FkLndpbmRvd05hbWUgPSB3aW5kb3dOYW1lO1xyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICBpZiAoYXJncyAhPT0gbnVsbCAmJiBhcmdzICE9PSB1bmRlZmluZWQpIHtcclxuICAgICAgICAgICAgICAgIHBheWxvYWQuYXJncyA9IGFyZ3M7XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIGNvbnN0IHJlc3BvbnNlVGV4dCA9ICh3aW5kb3cgYXMgYW55KS53YWlscy5pbnZva2UoSlNPTi5zdHJpbmdpZnkocGF5bG9hZCkpO1xyXG4gICAgICAgICAgICByZXR1cm4gcGFyc2VBbmRyb2lkSW52b2tlUmVzcG9uc2UocmVzcG9uc2VUZXh0KTtcclxuICAgICAgICB9XHJcbiAgICB9O1xyXG59XHJcblxyXG5jb25maWd1cmVBbmRyb2lkVHJhbnNwb3J0KCk7XHJcblxyXG4vKipcclxuICogU2V0IGEgY3VzdG9tIHRyYW5zcG9ydCBmb3IgYWxsIFdhaWxzIHJ1bnRpbWUgY2FsbHMuXHJcbiAqIFRoaXMgYWxsb3dzIHlvdSB0byByZXBsYWNlIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0IHdpdGhcclxuICogV2ViU29ja2V0cywgY3VzdG9tIHByb3RvY29scywgb3IgYW55IG90aGVyIG1lY2hhbmlzbS5cclxuICpcclxuICogQHBhcmFtIHRyYW5zcG9ydCAtIFlvdXIgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvblxyXG4gKlxyXG4gKiBAZXhhbXBsZVxyXG4gKiBgYGB0eXBlc2NyaXB0XHJcbiAqIGltcG9ydCB7IHNldFRyYW5zcG9ydCB9IGZyb20gJy93YWlscy9ydW50aW1lLmpzJztcclxuICpcclxuICogY29uc3Qgd3NUcmFuc3BvcnQgPSB7XHJcbiAqICAgY2FsbDogYXN5bmMgKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpID0+IHtcclxuICogICAgIC8vIFlvdXIgV2ViU29ja2V0IGltcGxlbWVudGF0aW9uXHJcbiAqICAgfVxyXG4gKiB9O1xyXG4gKlxyXG4gKiBzZXRUcmFuc3BvcnQod3NUcmFuc3BvcnQpO1xyXG4gKiBgYGBcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBzZXRUcmFuc3BvcnQodHJhbnNwb3J0OiBSdW50aW1lVHJhbnNwb3J0IHwgbnVsbCk6IHZvaWQge1xyXG4gICAgY3VzdG9tVHJhbnNwb3J0ID0gdHJhbnNwb3J0O1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBjdXJyZW50IHRyYW5zcG9ydCAodXNlZnVsIGZvciBleHRlbmRpbmcvd3JhcHBpbmcpXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZ2V0VHJhbnNwb3J0KCk6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsIHtcclxuICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxyXG4gKlxyXG4gKiBAcGFyYW0gb2JqZWN0IC0gVGhlIG9iamVjdCB0byBpbnZva2UgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcmV0dXJuIFRoZSBuZXcgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nID0gJycpIHtcclxuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kOiBudW1iZXIsIGFyZ3M6IGFueSA9IG51bGwpIHtcclxuICAgICAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0LCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xyXG4gICAgfTtcclxufVxyXG5cclxuYXN5bmMgZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+IHtcclxuICAgIC8vIFVzZSBjdXN0b20gdHJhbnNwb3J0IGlmIGF2YWlsYWJsZVxyXG4gICAgaWYgKGN1c3RvbVRyYW5zcG9ydCkge1xyXG4gICAgICAgIHJldHVybiBjdXN0b21UcmFuc3BvcnQuY2FsbChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH1cclxuXHJcbiAgICAvLyBEZWZhdWx0IEhUVFAgZmV0Y2ggdHJhbnNwb3J0XHJcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcclxuXHJcbiAgICBsZXQgYm9keTogeyBvYmplY3Q6IG51bWJlcjsgbWV0aG9kOiBudW1iZXIsIGFyZ3M/OiBhbnkgfSA9IHtcclxuICAgICAgb2JqZWN0OiBvYmplY3RJRCxcclxuICAgICAgbWV0aG9kXHJcbiAgICB9XHJcbiAgICBpZiAoYXJncyAhPT0gbnVsbCAmJiBhcmdzICE9PSB1bmRlZmluZWQpIHtcclxuICAgICAgYm9keS5hcmdzID0gYXJncztcclxuICAgIH1cclxuXHJcbiAgICBsZXQgaGVhZGVyczogUmVjb3JkPHN0cmluZywgc3RyaW5nPiA9IHtcclxuICAgICAgICBbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXTogY2xpZW50SWQsXHJcbiAgICAgICAgW1wiQ29udGVudC1UeXBlXCJdOiBcImFwcGxpY2F0aW9uL2pzb25cIlxyXG4gICAgfVxyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBoZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XHJcbiAgICB9XHJcblxyXG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2godXJsLCB7XHJcbiAgICAgIG1ldGhvZDogJ1BPU1QnLFxyXG4gICAgICBoZWFkZXJzLFxyXG4gICAgICBib2R5OiBKU09OLnN0cmluZ2lmeShib2R5KVxyXG4gICAgfSk7XHJcbiAgICBpZiAoIXJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbi8vIHNldHVwXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRGlhbG9nKTtcclxuXHJcbi8vIERlZmluZSBjb25zdGFudHMgZnJvbSB0aGUgYG1ldGhvZHNgIG9iamVjdCBpbiBUaXRsZSBDYXNlXHJcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xyXG5jb25zdCBEaWFsb2dXYXJuaW5nID0gMTtcclxuY29uc3QgRGlhbG9nRXJyb3IgPSAyO1xyXG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XHJcbmNvbnN0IERpYWxvZ09wZW5GaWxlID0gNDtcclxuY29uc3QgRGlhbG9nU2F2ZUZpbGUgPSA1O1xyXG5cclxuZXhwb3J0IGludGVyZmFjZSBPcGVuRmlsZURpYWxvZ09wdGlvbnMge1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY2hvc2VuLiAqL1xyXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlcyBjYW4gYmUgY2hvc2VuLiAqL1xyXG4gICAgQ2FuQ2hvb3NlRmlsZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cclxuICAgIENhbkNyZWF0ZURpcmVjdG9yaWVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGZpbGVzIHNob3VsZCBiZSBzaG93bi4gKi9cclxuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGFsaWFzZXMgc2hvdWxkIGJlIHJlc29sdmVkLiAqL1xyXG4gICAgUmVzb2x2ZXNBbGlhc2VzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgbXVsdGlwbGUgc2VsZWN0aW9uIGlzIGFsbG93ZWQuICovXHJcbiAgICBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBleHRlbnNpb24gc2hvdWxkIGJlIGhpZGRlbi4gKi9cclxuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZXh0ZW5zaW9ucyBjYW4gYmUgc2VsZWN0ZWQuICovXHJcbiAgICBDYW5TZWxlY3RIaWRkZW5FeHRlbnNpb24/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xyXG4gICAgVHJlYXRzRmlsZVBhY2thZ2VzQXNEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIG90aGVyIGZpbGUgdHlwZXMgYXJlIGFsbG93ZWQuICovXHJcbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XHJcbiAgICAvKiogQXJyYXkgb2YgZmlsZSBmaWx0ZXJzLiAqL1xyXG4gICAgRmlsdGVycz86IEZpbGVGaWx0ZXJbXTtcclxuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgVGl0bGU/OiBzdHJpbmc7XHJcbiAgICAvKiogTWVzc2FnZSB0byBzaG93IGluIHRoZSBkaWFsb2cuICovXHJcbiAgICBNZXNzYWdlPzogc3RyaW5nO1xyXG4gICAgLyoqIFRleHQgdG8gZGlzcGxheSBvbiB0aGUgYnV0dG9uLiAqL1xyXG4gICAgQnV0dG9uVGV4dD86IHN0cmluZztcclxuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgRGlyZWN0b3J5Pzogc3RyaW5nO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZGlhbG9nIHNob3VsZCBhcHBlYXIgZGV0YWNoZWQgZnJvbSB0aGUgbWFpbiB3aW5kb3cuICovXHJcbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgU2F2ZUZpbGVEaWFsb2dPcHRpb25zIHtcclxuICAgIC8qKiBEZWZhdWx0IGZpbGVuYW1lIHRvIHVzZSBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgRmlsZW5hbWU/OiBzdHJpbmc7XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VEaXJlY3Rvcmllcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGZpbGVzIGNhbiBiZSBjaG9zZW4uICovXHJcbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XHJcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjcmVhdGVkLiAqL1xyXG4gICAgQ2FuQ3JlYXRlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xyXG4gICAgU2hvd0hpZGRlbkZpbGVzPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgYWxpYXNlcyBzaG91bGQgYmUgcmVzb2x2ZWQuICovXHJcbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiB0aGUgZXh0ZW5zaW9uIHNob3VsZCBiZSBoaWRkZW4uICovXHJcbiAgICBIaWRlRXh0ZW5zaW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xyXG4gICAgQ2FuU2VsZWN0SGlkZGVuRXh0ZW5zaW9uPzogYm9vbGVhbjtcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZSBwYWNrYWdlcyBzaG91bGQgYmUgdHJlYXRlZCBhcyBkaXJlY3Rvcmllcy4gKi9cclxuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEluZGljYXRlcyBpZiBvdGhlciBmaWxlIHR5cGVzIGFyZSBhbGxvd2VkLiAqL1xyXG4gICAgQWxsb3dzT3RoZXJGaWxldHlwZXM/OiBib29sZWFuO1xyXG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cclxuICAgIEZpbHRlcnM/OiBGaWxlRmlsdGVyW107XHJcbiAgICAvKiogVGl0bGUgb2YgdGhlIGRpYWxvZy4gKi9cclxuICAgIFRpdGxlPzogc3RyaW5nO1xyXG4gICAgLyoqIE1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgTWVzc2FnZT86IHN0cmluZztcclxuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cclxuICAgIEJ1dHRvblRleHQ/OiBzdHJpbmc7XHJcbiAgICAvKiogRGlyZWN0b3J5IHRvIG9wZW4gaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIERpcmVjdG9yeT86IHN0cmluZztcclxuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93LiAqL1xyXG4gICAgRGV0YWNoZWQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIE1lc3NhZ2VEaWFsb2dPcHRpb25zIHtcclxuICAgIC8qKiBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuICovXHJcbiAgICBUaXRsZT86IHN0cmluZztcclxuICAgIC8qKiBUaGUgbWFpbiBtZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cclxuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XHJcbiAgICAvKiogQXJyYXkgb2YgYnV0dG9uIG9wdGlvbnMgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xyXG4gICAgQnV0dG9ucz86IEJ1dHRvbltdO1xyXG4gICAgLyoqIFRydWUgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93IChpZiBhcHBsaWNhYmxlKS4gKi9cclxuICAgIERldGFjaGVkPzogYm9vbGVhbjtcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBCdXR0b24ge1xyXG4gICAgLyoqIFRleHQgdGhhdCBhcHBlYXJzIHdpdGhpbiB0aGUgYnV0dG9uLiAqL1xyXG4gICAgTGFiZWw/OiBzdHJpbmc7XHJcbiAgICAvKiogVHJ1ZSBpZiB0aGUgYnV0dG9uIHNob3VsZCBjYW5jZWwgYW4gb3BlcmF0aW9uIHdoZW4gY2xpY2tlZC4gKi9cclxuICAgIElzQ2FuY2VsPzogYm9vbGVhbjtcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGJlIHRoZSBkZWZhdWx0IGFjdGlvbiB3aGVuIHRoZSB1c2VyIHByZXNzZXMgZW50ZXIuICovXHJcbiAgICBJc0RlZmF1bHQ/OiBib29sZWFuO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIEZpbGVGaWx0ZXIge1xyXG4gICAgLyoqIERpc3BsYXkgbmFtZSBmb3IgdGhlIGZpbHRlciwgaXQgY291bGQgYmUgXCJUZXh0IEZpbGVzXCIsIFwiSW1hZ2VzXCIgZXRjLiAqL1xyXG4gICAgRGlzcGxheU5hbWU/OiBzdHJpbmc7XHJcbiAgICAvKiogUGF0dGVybiB0byBtYXRjaCBmb3IgdGhlIGZpbHRlciwgZS5nLiBcIioudHh0OyoubWRcIiBmb3IgdGV4dCBtYXJrZG93biBmaWxlcy4gKi9cclxuICAgIFBhdHRlcm4/OiBzdHJpbmc7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIGRpYWxvZyBvZiBzcGVjaWZpZWQgdHlwZSB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKlxyXG4gKiBAcGFyYW0gdHlwZSAtIERpYWxvZyB0eXBlLlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgZm9yIHRoZSBkaWFsb2cuXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggcmVzdWx0IG9mIGRpYWxvZy5cclxuICovXHJcbmZ1bmN0aW9uIGRpYWxvZyh0eXBlOiBudW1iZXIsIG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zIHwgT3BlbkZpbGVEaWFsb2dPcHRpb25zIHwgU2F2ZUZpbGVEaWFsb2dPcHRpb25zID0ge30pOiBQcm9taXNlPGFueT4ge1xyXG4gICAgcmV0dXJuIGNhbGwodHlwZSwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhbiBpbmZvIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dJbmZvLCBvcHRpb25zKTsgfVxyXG5cclxuLyoqXHJcbiAqIFByZXNlbnRzIGEgd2FybmluZyBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFdhcm5pbmcob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYW4gZXJyb3IgZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpOyB9XHJcblxyXG4vKipcclxuICogUHJlc2VudHMgYSBxdWVzdGlvbiBkaWFsb2cuXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dRdWVzdGlvbiwgb3B0aW9ucyk7IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIG9uZSBvciBtb3JlIGZpbGVzIHRvIG9wZW4uXHJcbiAqXHJcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXHJcbiAqIEByZXR1cm5zIFNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcywgb3IgYSBibGFuayBzdHJpbmcvZW1wdHkgbGlzdCBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyAmIHsgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb246IHRydWUgfSk6IFByb21pc2U8c3RyaW5nW10+O1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9uczogT3BlbkZpbGVEaWFsb2dPcHRpb25zICYgeyBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbj86IGZhbHNlIHwgdW5kZWZpbmVkIH0pOiBQcm9taXNlPHN0cmluZz47XHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zOiBPcGVuRmlsZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZyB8IHN0cmluZ1tdPjtcclxuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dPcGVuRmlsZSwgb3B0aW9ucykgPz8gW107IH1cclxuXHJcbi8qKlxyXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIGEgZmlsZSB0byBzYXZlLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxyXG4gKiBAcmV0dXJucyBTZWxlY3RlZCBmaWxlLCBvciBhIGJsYW5rIHN0cmluZyBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnM6IFNhdmVGaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nU2F2ZUZpbGUsIG9wdGlvbnMpOyB9XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuaW1wb3J0IHsgZXZlbnRMaXN0ZW5lcnMsIExpc3RlbmVyLCBsaXN0ZW5lck9mZiB9IGZyb20gXCIuL2xpc3RlbmVyLmpzXCI7XHJcbmltcG9ydCB7IEV2ZW50cyBhcyBDcmVhdGUgfSBmcm9tIFwiLi9jcmVhdGUuanNcIjtcclxuaW1wb3J0IHsgVHlwZXMgfSBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xyXG5cclxuLy8gU2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbndpbmRvdy5fd2FpbHMuZGlzcGF0Y2hXYWlsc0V2ZW50ID0gZGlzcGF0Y2hXYWlsc0V2ZW50O1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRXZlbnRzKTtcclxuY29uc3QgRW1pdE1ldGhvZCA9IDA7XHJcblxyXG5leHBvcnQgKiBmcm9tIFwiLi9ldmVudF90eXBlcy5qc1wiO1xyXG5cclxuLyoqXHJcbiAqIEEgdGFibGUgb2YgZGF0YSB0eXBlcyBmb3IgYWxsIGtub3duIGV2ZW50cy5cclxuICogV2lsbCBiZSBtb25rZXktcGF0Y2hlZCBieSB0aGUgYmluZGluZyBnZW5lcmF0b3IuXHJcbiAqL1xyXG5leHBvcnQgaW50ZXJmYWNlIEN1c3RvbUV2ZW50cyB7fVxyXG5cclxuLyoqXHJcbiAqIEVpdGhlciBhIGtub3duIGV2ZW50IG5hbWUgb3IgYW4gYXJiaXRyYXJ5IHN0cmluZy5cclxuICovXHJcbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnROYW1lPEUgZXh0ZW5kcyBrZXlvZiBDdXN0b21FdmVudHMgPSBrZXlvZiBDdXN0b21FdmVudHM+ID0gRSB8IChzdHJpbmcgJiB7fSk7XHJcblxyXG4vKipcclxuICogVW5pb24gb2YgYWxsIGtub3duIHN5c3RlbSBldmVudCBuYW1lcy5cclxuICovXHJcbnR5cGUgU3lzdGVtRXZlbnROYW1lID0ge1xyXG4gICAgW0sgaW4ga2V5b2YgKHR5cGVvZiBUeXBlcyldOiAodHlwZW9mIFR5cGVzKVtLXVtrZXlvZiAoKHR5cGVvZiBUeXBlcylbS10pXVxyXG59IGV4dGVuZHMgKGluZmVyIE0pID8gTVtrZXlvZiBNXSA6IG5ldmVyO1xyXG5cclxuLyoqXHJcbiAqIFRoZSBkYXRhIHR5cGUgYXNzb2NpYXRlZCB0byBhIGdpdmVuIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IHR5cGUgV2FpbHNFdmVudERhdGE8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+ID1cclxuICAgIEUgZXh0ZW5kcyBrZXlvZiBDdXN0b21FdmVudHMgPyBDdXN0b21FdmVudHNbRV0gOiAoRSBleHRlbmRzIFN5c3RlbUV2ZW50TmFtZSA/IHZvaWQgOiBhbnkpO1xyXG5cclxuLyoqXHJcbiAqIFRoZSB0eXBlIG9mIGhhbmRsZXJzIGZvciBhIGdpdmVuIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IHR5cGUgV2FpbHNFdmVudENhbGxiYWNrPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiA9IChldjogV2FpbHNFdmVudDxFPikgPT4gdm9pZDtcclxuXHJcbi8qKlxyXG4gKiBSZXByZXNlbnRzIGEgc3lzdGVtIGV2ZW50IG9yIGEgY3VzdG9tIGV2ZW50IGVtaXR0ZWQgdGhyb3VnaCB3YWlscy1wcm92aWRlZCBmYWNpbGl0aWVzLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+IHtcclxuICAgIC8qKlxyXG4gICAgICogVGhlIG5hbWUgb2YgdGhlIGV2ZW50LlxyXG4gICAgICovXHJcbiAgICBuYW1lOiBFO1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogT3B0aW9uYWwgZGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGVtaXR0ZWQgZXZlbnQuXHJcbiAgICAgKi9cclxuICAgIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogTmFtZSBvZiB0aGUgb3JpZ2luYXRpbmcgd2luZG93LiBPbWl0dGVkIGZvciBhcHBsaWNhdGlvbiBldmVudHMuXHJcbiAgICAgKiBXaWxsIGJlIG92ZXJyaWRkZW4gaWYgc2V0IG1hbnVhbGx5LlxyXG4gICAgICovXHJcbiAgICBzZW5kZXI/OiBzdHJpbmc7XHJcblxyXG4gICAgY29uc3RydWN0b3IobmFtZTogRSwgZGF0YTogV2FpbHNFdmVudERhdGE8RT4pO1xyXG4gICAgY29uc3RydWN0b3IobmFtZTogV2FpbHNFdmVudERhdGE8RT4gZXh0ZW5kcyBudWxsIHwgdm9pZCA/IEUgOiBuZXZlcilcclxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IEUsIGRhdGE/OiBhbnkpIHtcclxuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xyXG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGEgPz8gbnVsbDtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50OiBhbnkpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnQubmFtZSk7XHJcbiAgICBpZiAoIWxpc3RlbmVycykge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBsZXQgd2FpbHNFdmVudCA9IG5ldyBXYWlsc0V2ZW50KFxyXG4gICAgICAgIGV2ZW50Lm5hbWUsXHJcbiAgICAgICAgKGV2ZW50Lm5hbWUgaW4gQ3JlYXRlKSA/IENyZWF0ZVtldmVudC5uYW1lXShldmVudC5kYXRhKSA6IGV2ZW50LmRhdGFcclxuICAgICk7XHJcbiAgICBpZiAoJ3NlbmRlcicgaW4gZXZlbnQpIHtcclxuICAgICAgICB3YWlsc0V2ZW50LnNlbmRlciA9IGV2ZW50LnNlbmRlcjtcclxuICAgIH1cclxuXHJcbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+ICFsaXN0ZW5lci5kaXNwYXRjaCh3YWlsc0V2ZW50KSk7XHJcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlciBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCBtdWx0aXBsZSB0aW1lcyBmb3IgYSBzcGVjaWZpYyBldmVudC5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxyXG4gKiBAcGFyYW0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHBhcmFtIG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxyXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGU8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPiwgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcclxuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIHNwZWNpZmllZCBldmVudCBvY2N1cnMuXHJcbiAqXHJcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXHJcbiAqIEByZXR1cm5zIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT248RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPik6ICgpID0+IHZvaWQge1xyXG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgb25seSBvbmNlIGZvciB0aGUgc3BlY2lmaWVkIGV2ZW50LlxyXG4gKlxyXG4gKiBAcGFyYW0gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmb3IuXHJcbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxyXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uY2U8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPik6ICgpID0+IHZvaWQge1xyXG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIGV2ZW50IGxpc3RlbmVycyBmb3IgdGhlIHNwZWNpZmllZCBldmVudCBuYW1lcy5cclxuICpcclxuICogQHBhcmFtIGV2ZW50TmFtZXMgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnRzIHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZiguLi5ldmVudE5hbWVzOiBbV2FpbHNFdmVudE5hbWUsIC4uLldhaWxzRXZlbnROYW1lW11dKTogdm9pZCB7XHJcbiAgICBldmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlbW92ZXMgYWxsIGV2ZW50IGxpc3RlbmVycy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKTogdm9pZCB7XHJcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xyXG59XHJcblxyXG4vKipcclxuICogRW1pdHMgYW4gZXZlbnQuXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHdpbGwgYmUgZnVsZmlsbGVkIG9uY2UgdGhlIGV2ZW50IGhhcyBiZWVuIGVtaXR0ZWQuICBSZXNvbHZlcyB0byB0cnVlIGlmIHRoZSBldmVudCB3YXMgY2FuY2VsbGVkLlxyXG4gKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byBlbWl0XHJcbiAqIEBwYXJhbSBkYXRhIC0gVGhlIGRhdGEgdGhhdCB3aWxsIGJlIHNlbnQgd2l0aCB0aGUgZXZlbnRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbWl0PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihuYW1lOiBFLCBkYXRhOiBXYWlsc0V2ZW50RGF0YTxFPik6IFByb21pc2U8Ym9vbGVhbj5cclxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IFdhaWxzRXZlbnREYXRhPEU+IGV4dGVuZHMgbnVsbCB8IHZvaWQgPyBFIDogbmV2ZXIpOiBQcm9taXNlPGJvb2xlYW4+XHJcbmV4cG9ydCBmdW5jdGlvbiBFbWl0PEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiwgZGF0YT86IGFueSk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgIG5ldyBXYWlsc0V2ZW50KG5hbWUsIGRhdGEpKVxyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLy8gVGhlIGZvbGxvd2luZyB1dGlsaXRpZXMgaGF2ZSBiZWVuIGZhY3RvcmVkIG91dCBvZiAuL2V2ZW50cy50c1xyXG4vLyBmb3IgdGVzdGluZyBwdXJwb3Nlcy5cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXA8c3RyaW5nLCBMaXN0ZW5lcltdPigpO1xyXG5cclxuZXhwb3J0IGNsYXNzIExpc3RlbmVyIHtcclxuICAgIGV2ZW50TmFtZTogc3RyaW5nO1xyXG4gICAgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQ7XHJcbiAgICBtYXhDYWxsYmFja3M6IG51bWJlcjtcclxuXHJcbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWU6IHN0cmluZywgY2FsbGJhY2s6IChkYXRhOiBhbnkpID0+IHZvaWQsIG1heENhbGxiYWNrczogbnVtYmVyKSB7XHJcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XHJcbiAgICAgICAgdGhpcy5jYWxsYmFjayA9IGNhbGxiYWNrO1xyXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xyXG4gICAgfVxyXG5cclxuICAgIGRpc3BhdGNoKGRhdGE6IGFueSk6IGJvb2xlYW4ge1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIHRoaXMuY2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgIGNvbnNvbGUuZXJyb3IoZXJyKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyOiBMaXN0ZW5lcik6IHZvaWQge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChsaXN0ZW5lci5ldmVudE5hbWUpO1xyXG4gICAgaWYgKCFsaXN0ZW5lcnMpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGxpc3RlbmVyLmV2ZW50TmFtZSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChsaXN0ZW5lci5ldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBBbnkgaXMgYSBkdW1teSBjcmVhdGlvbiBmdW5jdGlvbiBmb3Igc2ltcGxlIG9yIHVua25vd24gdHlwZXMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQW55PFQgPSBhbnk+KHNvdXJjZTogYW55KTogVCB7XHJcbiAgICByZXR1cm4gc291cmNlO1xyXG59XHJcblxyXG4vKipcclxuICogQnl0ZVNsaWNlIGlzIGEgY3JlYXRpb24gZnVuY3Rpb24gdGhhdCByZXBsYWNlc1xyXG4gKiBudWxsIHN0cmluZ3Mgd2l0aCBlbXB0eSBzdHJpbmdzLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5dGVTbGljZShzb3VyY2U6IGFueSk6IHN0cmluZyB7XHJcbiAgICByZXR1cm4gKChzb3VyY2UgPT0gbnVsbCkgPyBcIlwiIDogc291cmNlKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEFycmF5IHRha2VzIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGFuIGFyYml0cmFyeSB0eXBlXHJcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcnJheVxyXG4gKiB3aG9zZSBlbGVtZW50cyBhcmUgb2YgdGhhdCB0eXBlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEFycmF5PFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gVFtdIHtcclxuICAgIGlmIChlbGVtZW50ID09PSBBbnkpIHtcclxuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IFtdIDogc291cmNlKTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xyXG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcclxuICAgICAgICAgICAgcmV0dXJuIFtdO1xyXG4gICAgICAgIH1cclxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNvdXJjZS5sZW5ndGg7IGkrKykge1xyXG4gICAgICAgICAgICBzb3VyY2VbaV0gPSBlbGVtZW50KHNvdXJjZVtpXSk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBzb3VyY2U7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogTWFwIHRha2VzIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdHdvIGFyYml0cmFyeSB0eXBlc1xyXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gb2JqZWN0XHJcbiAqIHdob3NlIGtleXMgYW5kIHZhbHVlcyBhcmUgb2YgdGhvc2UgdHlwZXMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTWFwPEsgZXh0ZW5kcyBQcm9wZXJ0eUtleSA9IGFueSwgViA9IGFueT4oa2V5OiAoc291cmNlOiBhbnkpID0+IEssIHZhbHVlOiAoc291cmNlOiBhbnkpID0+IFYpOiAoc291cmNlOiBhbnkpID0+IFJlY29yZDxLLCBWPiB7XHJcbiAgICBpZiAodmFsdWUgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8ge30gOiBzb3VyY2UpO1xyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XHJcbiAgICAgICAgaWYgKHNvdXJjZSA9PT0gbnVsbCkge1xyXG4gICAgICAgICAgICByZXR1cm4ge307XHJcbiAgICAgICAgfVxyXG4gICAgICAgIGZvciAoY29uc3Qga2V5IGluIHNvdXJjZSkge1xyXG4gICAgICAgICAgICBzb3VyY2Vba2V5XSA9IHZhbHVlKHNvdXJjZVtrZXldKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuIHNvdXJjZTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBOdWxsYWJsZSB0YWtlcyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcmJpdHJhcnkgdHlwZVxyXG4gKiBhbmQgcmV0dXJucyBhIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhIG51bGxhYmxlIHZhbHVlIG9mIHRoYXQgdHlwZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBOdWxsYWJsZTxUID0gYW55PihlbGVtZW50OiAoc291cmNlOiBhbnkpID0+IFQpOiAoc291cmNlOiBhbnkpID0+IChUIHwgbnVsbCkge1xyXG4gICAgaWYgKGVsZW1lbnQgPT09IEFueSkge1xyXG4gICAgICAgIHJldHVybiBBbnk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFN0cnVjdCB0YWtlcyBhbiBvYmplY3QgbWFwcGluZyBmaWVsZCBuYW1lcyB0byBjcmVhdGlvbiBmdW5jdGlvbnNcclxuICogYW5kIHJldHVybnMgYW4gaW4tcGxhY2UgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgc3RydWN0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFN0cnVjdChjcmVhdGVGaWVsZDogUmVjb3JkPHN0cmluZywgKHNvdXJjZTogYW55KSA9PiBhbnk+KTpcclxuICAgIDxVIGV4dGVuZHMgUmVjb3JkPHN0cmluZywgYW55PiA9IGFueT4oc291cmNlOiBhbnkpID0+IFVcclxue1xyXG4gICAgbGV0IGFsbEFueSA9IHRydWU7XHJcbiAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcclxuICAgICAgICBpZiAoY3JlYXRlRmllbGRbbmFtZV0gIT09IEFueSkge1xyXG4gICAgICAgICAgICBhbGxBbnkgPSBmYWxzZTtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG4gICAgaWYgKGFsbEFueSkge1xyXG4gICAgICAgIHJldHVybiBBbnk7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcclxuICAgICAgICBmb3IgKGNvbnN0IG5hbWUgaW4gY3JlYXRlRmllbGQpIHtcclxuICAgICAgICAgICAgaWYgKG5hbWUgaW4gc291cmNlKSB7XHJcbiAgICAgICAgICAgICAgICBzb3VyY2VbbmFtZV0gPSBjcmVhdGVGaWVsZFtuYW1lXShzb3VyY2VbbmFtZV0pO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBzb3VyY2U7XHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogTWFwcyBrbm93biBldmVudCBuYW1lcyB0byBjcmVhdGlvbiBmdW5jdGlvbnMgZm9yIHRoZWlyIGRhdGEgdHlwZXMuXHJcbiAqIFdpbGwgYmUgbW9ua2V5LXBhdGNoZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IEV2ZW50czogUmVjb3JkPHN0cmluZywgKHNvdXJjZTogYW55KSA9PiBhbnk+ID0ge307XHJcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gQ3luaHlyY2h3eWQgeSBmZmVpbCBob24geW4gYXd0b21hdGlnLiBQRUlESVdDSCBcdTAwQzIgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHR9KSxcblx0TGludXg6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibGludXg6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcblx0XHRXaW5kb3dMb2FkU3RhcnRlZDogXCJsaW51eDpXaW5kb3dMb2FkU3RhcnRlZFwiLFxuXHRcdFdpbmRvd0xvYWRSZWRpcmVjdGVkOiBcImxpbnV4OldpbmRvd0xvYWRSZWRpcmVjdGVkXCIsXG5cdFx0V2luZG93TG9hZENvbW1pdHRlZDogXCJsaW51eDpXaW5kb3dMb2FkQ29tbWl0dGVkXCIsXG5cdFx0V2luZG93TG9hZEZpbmlzaGVkOiBcImxpbnV4OldpbmRvd0xvYWRGaW5pc2hlZFwiLFxuXHR9KSxcblx0aU9TOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJpb3M6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZEVudGVyQmFja2dyb3VuZDogXCJpb3M6QXBwbGljYXRpb25EaWRFbnRlckJhY2tncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZzogXCJpb3M6QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmdcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nOiBcImlvczpBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsRW50ZXJGb3JlZ3JvdW5kOiBcImlvczpBcHBsaWNhdGlvbldpbGxFbnRlckZvcmVncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJpb3M6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0V2luZG93RGlkTG9hZDogXCJpb3M6V2luZG93RGlkTG9hZFwiLFxuXHRcdFdpbmRvd1dpbGxBcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxBcHBlYXJcIixcblx0XHRXaW5kb3dEaWRBcHBlYXI6IFwiaW9zOldpbmRvd0RpZEFwcGVhclwiLFxuXHRcdFdpbmRvd1dpbGxEaXNhcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxEaXNhcHBlYXJcIixcblx0XHRXaW5kb3dEaWREaXNhcHBlYXI6IFwiaW9zOldpbmRvd0RpZERpc2FwcGVhclwiLFxuXHRcdFdpbmRvd1NhZmVBcmVhSW5zZXRzQ2hhbmdlZDogXCJpb3M6V2luZG93U2FmZUFyZWFJbnNldHNDaGFuZ2VkXCIsXG5cdFx0V2luZG93T3JpZW50YXRpb25DaGFuZ2VkOiBcImlvczpXaW5kb3dPcmllbnRhdGlvbkNoYW5nZWRcIixcblx0XHRXaW5kb3dUb3VjaEJlZ2FuOiBcImlvczpXaW5kb3dUb3VjaEJlZ2FuXCIsXG5cdFx0V2luZG93VG91Y2hNb3ZlZDogXCJpb3M6V2luZG93VG91Y2hNb3ZlZFwiLFxuXHRcdFdpbmRvd1RvdWNoRW5kZWQ6IFwiaW9zOldpbmRvd1RvdWNoRW5kZWRcIixcblx0XHRXaW5kb3dUb3VjaENhbmNlbGxlZDogXCJpb3M6V2luZG93VG91Y2hDYW5jZWxsZWRcIixcblx0XHRXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZhaWxOYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkRmFpbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGVjaWRlUG9saWN5Rm9yTmF2aWdhdGlvbkFjdGlvbjogXCJpb3M6V2ViVmlld0RlY2lkZVBvbGljeUZvck5hdmlnYXRpb25BY3Rpb25cIixcblx0fSksXG5cdENvbW1vbjogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZTogXCJjb21tb246QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZVwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmw6IFwiY29tbW9uOkFwcGxpY2F0aW9uTGF1bmNoZWRXaXRoVXJsXCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImNvbW1vbjpUaGVtZUNoYW5nZWRcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJjb21tb246V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJjb21tb246V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RmlsZXNEcm9wcGVkXCIsXG5cdFx0V2luZG93Rm9jdXM6IFwiY29tbW9uOldpbmRvd0ZvY3VzXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd0hpZGU6IFwiY29tbW9uOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dMb3N0Rm9jdXM6IFwiY29tbW9uOldpbmRvd0xvc3RGb2N1c1wiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1RvZ2dsZUZyYW1lbGVzczogXCJjb21tb246V2luZG93VG9nZ2xlRnJhbWVsZXNzXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1J1bnRpbWVSZWFkeTogXCJjb21tb246V2luZG93UnVudGltZVJlYWR5XCIsXG5cdFx0V2luZG93U2hvdzogXCJjb21tb246V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcImNvbW1vbjpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcblx0fSksXG59KTtcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKipcclxuICogTG9ncyBhIG1lc3NhZ2UgdG8gdGhlIGNvbnNvbGUgd2l0aCBjdXN0b20gZm9ybWF0dGluZy5cclxuICpcclxuICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgbWVzc2FnZSB0byBiZSBsb2dnZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZTogYW55KSB7XHJcbiAgICAvLyBlc2xpbnQtZGlzYWJsZS1uZXh0LWxpbmVcclxuICAgIGNvbnNvbGUubG9nKFxyXG4gICAgICAgICclYyB3YWlsczMgJWMgJyArIG1lc3NhZ2UgKyAnICcsXHJcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNhYTAwMDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAzcHggMHB4IDBweCAzcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjMDA5OTAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogMHB4IDNweCAzcHggMHB4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJ1xyXG4gICAgKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyB3aGV0aGVyIHRoZSB3ZWJ2aWV3IHN1cHBvcnRzIHRoZSB7QGxpbmsgTW91c2VFdmVudCNidXR0b25zfSBwcm9wZXJ0eS5cclxuICogTG9va2luZyBhdCB5b3UgbWFjT1MgSGlnaCBTaWVycmEhXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gY2FuVHJhY2tCdXR0b25zKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuIChuZXcgTW91c2VFdmVudCgnbW91c2Vkb3duJykpLmJ1dHRvbnMgPT09IDA7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgYnJvd3NlciBzdXBwb3J0cyByZW1vdmluZyBsaXN0ZW5lcnMgYnkgdHJpZ2dlcmluZyBhbiBBYm9ydFNpZ25hbFxyXG4gKiAoc2VlIGh0dHBzOi8vZGV2ZWxvcGVyLm1vemlsbGEub3JnL2VuLVVTL2RvY3MvV2ViL0FQSS9FdmVudFRhcmdldC9hZGRFdmVudExpc3RlbmVyI3NpZ25hbCkuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XHJcbiAgICBpZiAoIUV2ZW50VGFyZ2V0IHx8ICFBYm9ydFNpZ25hbCB8fCAhQWJvcnRDb250cm9sbGVyKVxyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuXHJcbiAgICBsZXQgcmVzdWx0ID0gdHJ1ZTtcclxuXHJcbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcclxuICAgIGNvbnN0IGNvbnRyb2xsZXIgPSBuZXcgQWJvcnRDb250cm9sbGVyKCk7XHJcbiAgICB0YXJnZXQuYWRkRXZlbnRMaXN0ZW5lcigndGVzdCcsICgpID0+IHsgcmVzdWx0ID0gZmFsc2U7IH0sIHsgc2lnbmFsOiBjb250cm9sbGVyLnNpZ25hbCB9KTtcclxuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcclxuICAgIHRhcmdldC5kaXNwYXRjaEV2ZW50KG5ldyBDdXN0b21FdmVudCgndGVzdCcpKTtcclxuXHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG4vKipcclxuICogUmVzb2x2ZXMgdGhlIGNsb3Nlc3QgSFRNTEVsZW1lbnQgYW5jZXN0b3Igb2YgYW4gZXZlbnQncyB0YXJnZXQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZXZlbnRUYXJnZXQoZXZlbnQ6IEV2ZW50KTogSFRNTEVsZW1lbnQge1xyXG4gICAgaWYgKGV2ZW50LnRhcmdldCBpbnN0YW5jZW9mIEhUTUxFbGVtZW50KSB7XHJcbiAgICAgICAgcmV0dXJuIGV2ZW50LnRhcmdldDtcclxuICAgIH0gZWxzZSBpZiAoIShldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkgJiYgZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgTm9kZSkge1xyXG4gICAgICAgIHJldHVybiBldmVudC50YXJnZXQucGFyZW50RWxlbWVudCA/PyBkb2N1bWVudC5ib2R5O1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICByZXR1cm4gZG9jdW1lbnQuYm9keTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqKlxyXG4gVGhpcyB0ZWNobmlxdWUgZm9yIHByb3BlciBsb2FkIGRldGVjdGlvbiBpcyB0YWtlbiBmcm9tIEhUTVg6XHJcblxyXG4gQlNEIDItQ2xhdXNlIExpY2Vuc2VcclxuXHJcbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcclxuIEFsbCByaWdodHMgcmVzZXJ2ZWQuXHJcblxyXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XHJcbiBtb2RpZmljYXRpb24sIGFyZSBwZXJtaXR0ZWQgcHJvdmlkZWQgdGhhdCB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnMgYXJlIG1ldDpcclxuXHJcbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcclxuIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyLlxyXG5cclxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcclxuIHRoaXMgbGlzdCBvZiBjb25kaXRpb25zIGFuZCB0aGUgZm9sbG93aW5nIGRpc2NsYWltZXIgaW4gdGhlIGRvY3VtZW50YXRpb25cclxuIGFuZC9vciBvdGhlciBtYXRlcmlhbHMgcHJvdmlkZWQgd2l0aCB0aGUgZGlzdHJpYnV0aW9uLlxyXG5cclxuIFRISVMgU09GVFdBUkUgSVMgUFJPVklERUQgQlkgVEhFIENPUFlSSUdIVCBIT0xERVJTIEFORCBDT05UUklCVVRPUlMgXCJBUyBJU1wiXHJcbiBBTkQgQU5ZIEVYUFJFU1MgT1IgSU1QTElFRCBXQVJSQU5USUVTLCBJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgVEhFXHJcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcclxuIERJU0NMQUlNRUQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRSBDT1BZUklHSFQgSE9MREVSIE9SIENPTlRSSUJVVE9SUyBCRSBMSUFCTEVcclxuIEZPUiBBTlkgRElSRUNULCBJTkRJUkVDVCwgSU5DSURFTlRBTCwgU1BFQ0lBTCwgRVhFTVBMQVJZLCBPUiBDT05TRVFVRU5USUFMXHJcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxyXG4gU0VSVklDRVM7IExPU1MgT0YgVVNFLCBEQVRBLCBPUiBQUk9GSVRTOyBPUiBCVVNJTkVTUyBJTlRFUlJVUFRJT04pIEhPV0VWRVJcclxuIENBVVNFRCBBTkQgT04gQU5ZIFRIRU9SWSBPRiBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQ09OVFJBQ1QsIFNUUklDVCBMSUFCSUxJVFksXHJcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxyXG4gT0YgVEhJUyBTT0ZUV0FSRSwgRVZFTiBJRiBBRFZJU0VEIE9GIFRIRSBQT1NTSUJJTElUWSBPRiBTVUNIIERBTUFHRS5cclxuXHJcbiAqKiovXHJcblxyXG5sZXQgaXNSZWFkeSA9IGZhbHNlO1xyXG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgKCkgPT4geyBpc1JlYWR5ID0gdHJ1ZSB9KTtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiB3aGVuUmVhZHkoY2FsbGJhY2s6ICgpID0+IHZvaWQpIHtcclxuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcclxuICAgICAgICBjYWxsYmFjaygpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgY2FsbGJhY2spO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB0eXBlIHsgU2NyZWVuIH0gZnJvbSBcIi4vc2NyZWVucy5qc1wiO1xyXG5cclxuLy8gRHJvcCB0YXJnZXQgY29uc3RhbnRzXHJcbmNvbnN0IERST1BfVEFSR0VUX0FUVFJJQlVURSA9ICdkYXRhLWZpbGUtZHJvcC10YXJnZXQnO1xyXG5jb25zdCBEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MgPSAnZmlsZS1kcm9wLXRhcmdldC1hY3RpdmUnO1xyXG5sZXQgY3VycmVudERyb3BUYXJnZXQ6IEVsZW1lbnQgfCBudWxsID0gbnVsbDtcclxuXHJcbmNvbnN0IFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDA7XHJcbmNvbnN0IENlbnRlck1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDE7XHJcbmNvbnN0IENsb3NlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDI7XHJcbmNvbnN0IERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICA9IDM7XHJcbmNvbnN0IEVuYWJsZVNpemVDb25zdHJhaW50c01ldGhvZCAgICAgICA9IDQ7XHJcbmNvbnN0IEZvY3VzTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDU7XHJcbmNvbnN0IEZvcmNlUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICA9IDY7XHJcbmNvbnN0IEZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgICA9IDc7XHJcbmNvbnN0IEdldFNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDg7XHJcbmNvbnN0IEdldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDk7XHJcbmNvbnN0IEhlaWdodE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICA9IDEwO1xyXG5jb25zdCBIaWRlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAxMTtcclxuY29uc3QgSXNGb2N1c2VkTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMTI7XHJcbmNvbnN0IElzRnVsbHNjcmVlbk1ldGhvZCAgICAgICAgICAgICAgICA9IDEzO1xyXG5jb25zdCBJc01heGltaXNlZE1ldGhvZCAgICAgICAgICAgICAgICAgPSAxNDtcclxuY29uc3QgSXNNaW5pbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTU7XHJcbmNvbnN0IE1heGltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE2O1xyXG5jb25zdCBNaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNztcclxuY29uc3QgTmFtZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTg7XHJcbmNvbnN0IE9wZW5EZXZUb29sc01ldGhvZCAgICAgICAgICAgICAgICA9IDE5O1xyXG5jb25zdCBSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgPSAyMDtcclxuY29uc3QgUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMjE7XHJcbmNvbnN0IFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDIyO1xyXG5jb25zdCBSZXN0b3JlTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAyMztcclxuY29uc3QgU2V0UG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgID0gMjQ7XHJcbmNvbnN0IFNldEFsd2F5c09uVG9wTWV0aG9kICAgICAgICAgICAgICA9IDI1O1xyXG5jb25zdCBTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kICAgICAgICAgPSAyNjtcclxuY29uc3QgU2V0RnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgICAgID0gMjc7XHJcbmNvbnN0IFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kICA9IDI4O1xyXG5jb25zdCBTZXRNYXhTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgPSAyOTtcclxuY29uc3QgU2V0TWluU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMzA7XHJcbmNvbnN0IFNldFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICA9IDMxO1xyXG5jb25zdCBTZXRSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgPSAzMjtcclxuY29uc3QgU2V0U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzM7XHJcbmNvbnN0IFNldFRpdGxlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDM0O1xyXG5jb25zdCBTZXRab29tTWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSAzNTtcclxuY29uc3QgU2hvd01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzY7XHJcbmNvbnN0IFNpemVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDM3O1xyXG5jb25zdCBUb2dnbGVGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgPSAzODtcclxuY29uc3QgVG9nZ2xlTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgID0gMzk7XHJcbmNvbnN0IFRvZ2dsZUZyYW1lbGVzc01ldGhvZCAgICAgICAgICAgICA9IDQwOyBcclxuY29uc3QgVW5GdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgID0gNDE7XHJcbmNvbnN0IFVuTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQyO1xyXG5jb25zdCBVbk1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgPSA0MztcclxuY29uc3QgV2lkdGhNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNDQ7XHJcbmNvbnN0IFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDQ1O1xyXG5jb25zdCBab29tSW5NZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSA0NjtcclxuY29uc3QgWm9vbU91dE1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gNDc7XHJcbmNvbnN0IFpvb21SZXNldE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDQ4O1xyXG5jb25zdCBTbmFwQXNzaXN0TWV0aG9kICAgICAgICAgICAgICAgICAgPSA0OTtcclxuY29uc3QgRmlsZXNEcm9wcGVkICAgICAgICAgICAgICAgICAgICAgID0gNTA7XHJcbmNvbnN0IFByaW50TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDUxO1xyXG5cclxuLyoqXHJcbiAqIEZpbmRzIHRoZSBuZWFyZXN0IGRyb3AgdGFyZ2V0IGVsZW1lbnQgYnkgd2Fsa2luZyB1cCB0aGUgRE9NIHRyZWUuXHJcbiAqL1xyXG5mdW5jdGlvbiBnZXREcm9wVGFyZ2V0RWxlbWVudChlbGVtZW50OiBFbGVtZW50IHwgbnVsbCk6IEVsZW1lbnQgfCBudWxsIHtcclxuICAgIGlmICghZWxlbWVudCkge1xyXG4gICAgICAgIHJldHVybiBudWxsO1xyXG4gICAgfVxyXG4gICAgcmV0dXJuIGVsZW1lbnQuY2xvc2VzdChgWyR7RFJPUF9UQVJHRVRfQVRUUklCVVRFfV1gKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrIGlmIHdlIGNhbiB1c2UgV2ViVmlldzIncyBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyAoV2luZG93cylcclxuICogQWxzbyBjaGVja3MgdGhhdCBFbmFibGVGaWxlRHJvcCBpcyB0cnVlIGZvciB0aGlzIHdpbmRvdy5cclxuICovXHJcbmZ1bmN0aW9uIGNhblJlc29sdmVGaWxlUGF0aHMoKTogYm9vbGVhbiB7XHJcbiAgICAvLyBNdXN0IGhhdmUgV2ViVmlldzIncyBwb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyBBUEkgKFdpbmRvd3Mgb25seSlcclxuICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyA9PSBudWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG4gICAgLy8gTXVzdCBoYXZlIEVuYWJsZUZpbGVEcm9wIHNldCB0byB0cnVlIGZvciB0aGlzIHdpbmRvd1xyXG4gICAgLy8gVGhpcyBmbGFnIGlzIHNldCBieSB0aGUgR28gYmFja2VuZCBkdXJpbmcgcnVudGltZSBpbml0aWFsaXphdGlvblxyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gdHJ1ZTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNlbmQgZmlsZSBkcm9wIHRvIGJhY2tlbmQgdmlhIFdlYlZpZXcyIChXaW5kb3dzIG9ubHkpXHJcbiAqL1xyXG5mdW5jdGlvbiByZXNvbHZlRmlsZVBhdGhzKHg6IG51bWJlciwgeTogbnVtYmVyLCBmaWxlczogRmlsZVtdKTogdm9pZCB7XHJcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMpIHtcclxuICAgICAgICAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMoYGZpbGU6ZHJvcDoke3h9OiR7eX1gLCBmaWxlcyk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8vIE5hdGl2ZSBkcmFnIHN0YXRlIChMaW51eC9tYWNPUyBpbnRlcmNlcHQgRE9NIGRyYWcgZXZlbnRzKVxyXG5sZXQgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xyXG5cclxuLyoqXHJcbiAqIENsZWFucyB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhbmQgaG92ZXIgZWZmZWN0cy5cclxuICogQ2FsbGVkIG9uIGRyb3Agb3Igd2hlbiBkcmFnIGxlYXZlcyB0aGUgd2luZG93LlxyXG4gKi9cclxuZnVuY3Rpb24gY2xlYW51cE5hdGl2ZURyYWcoKTogdm9pZCB7XHJcbiAgICBuYXRpdmVEcmFnQWN0aXZlID0gZmFsc2U7XHJcbiAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcclxuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogQ2FsbGVkIGZyb20gR28gd2hlbiBhIGZpbGUgZHJhZyBlbnRlcnMgdGhlIHdpbmRvdyBvbiBMaW51eC9tYWNPUy5cclxuICovXHJcbmZ1bmN0aW9uIGhhbmRsZURyYWdFbnRlcigpOiB2b2lkIHtcclxuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgZG9uJ3QgYWN0aXZhdGUgZHJhZyBzdGF0ZVxyXG4gICAgfVxyXG4gICAgbmF0aXZlRHJhZ0FjdGl2ZSA9IHRydWU7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsZWQgZnJvbSBHbyB3aGVuIGEgZmlsZSBkcmFnIGxlYXZlcyB0aGUgd2luZG93IG9uIExpbnV4L21hY09TLlxyXG4gKi9cclxuZnVuY3Rpb24gaGFuZGxlRHJhZ0xlYXZlKCk6IHZvaWQge1xyXG4gICAgY2xlYW51cE5hdGl2ZURyYWcoKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxlZCBmcm9tIEdvIGR1cmluZyBmaWxlIGRyYWcgdG8gdXBkYXRlIGhvdmVyIHN0YXRlIG9uIExpbnV4L21hY09TLlxyXG4gKiBAcGFyYW0geCAtIFggY29vcmRpbmF0ZSBpbiBDU1MgcGl4ZWxzXHJcbiAqIEBwYXJhbSB5IC0gWSBjb29yZGluYXRlIGluIENTUyBwaXhlbHNcclxuICovXHJcbmZ1bmN0aW9uIGhhbmRsZURyYWdPdmVyKHg6IG51bWJlciwgeTogbnVtYmVyKTogdm9pZCB7XHJcbiAgICBpZiAoIW5hdGl2ZURyYWdBY3RpdmUpIHJldHVybjtcclxuICAgIFxyXG4gICAgLy8gQ2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XHJcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcclxuICAgIH1cclxuICAgIFxyXG4gICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XHJcbiAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XHJcbiAgICBcclxuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xyXG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgIH1cclxuICAgIFxyXG4gICAgaWYgKGRyb3BUYXJnZXQpIHtcclxuICAgICAgICBkcm9wVGFyZ2V0LmNsYXNzTGlzdC5hZGQoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcclxuICAgIH1cclxufVxyXG5cclxuXHJcblxyXG4vLyBFeHBvcnQgdGhlIGhhbmRsZXJzIGZvciB1c2UgYnkgR28gdmlhIGluZGV4LnRzXHJcbmV4cG9ydCB7IGhhbmRsZURyYWdFbnRlciwgaGFuZGxlRHJhZ0xlYXZlLCBoYW5kbGVEcmFnT3ZlciB9O1xyXG5cclxuLyoqXHJcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHBvc2l0aW9uIG9mIGEgd2luZG93LlxyXG4gKi9cclxuaW50ZXJmYWNlIFBvc2l0aW9uIHtcclxuICAgIC8qKiBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xyXG4gICAgeDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSB2ZXJ0aWNhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xyXG4gICAgeTogbnVtYmVyO1xyXG59XHJcblxyXG4vKipcclxuICogQSByZWNvcmQgZGVzY3JpYmluZyB0aGUgc2l6ZSBvZiBhIHdpbmRvdy5cclxuICovXHJcbmludGVyZmFjZSBTaXplIHtcclxuICAgIC8qKiBUaGUgd2lkdGggb2YgdGhlIHdpbmRvdy4gKi9cclxuICAgIHdpZHRoOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LiAqL1xyXG4gICAgaGVpZ2h0OiBudW1iZXI7XHJcbn1cclxuXHJcbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXHJcbmNvbnN0IGNhbGxlclN5bSA9IFN5bWJvbChcImNhbGxlclwiKTtcclxuXHJcbmNsYXNzIFdpbmRvdyB7XHJcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cclxuICAgIHByaXZhdGUgW2NhbGxlclN5bV06IChtZXNzYWdlOiBudW1iZXIsIGFyZ3M/OiBhbnkpID0+IFByb21pc2U8YW55PjtcclxuXHJcbiAgICAvKipcclxuICAgICAqIEluaXRpYWxpc2VzIGEgd2luZG93IG9iamVjdCB3aXRoIHRoZSBzcGVjaWZpZWQgbmFtZS5cclxuICAgICAqXHJcbiAgICAgKiBAcHJpdmF0ZVxyXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgdGFyZ2V0IHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IobmFtZTogc3RyaW5nID0gJycpIHtcclxuICAgICAgICB0aGlzW2NhbGxlclN5bV0gPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcclxuXHJcbiAgICAgICAgLy8gYmluZCBpbnN0YW5jZSBtZXRob2QgdG8gbWFrZSB0aGVtIGVhc2lseSB1c2FibGUgaW4gZXZlbnQgaGFuZGxlcnNcclxuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xyXG4gICAgICAgICAgICBpZiAoXHJcbiAgICAgICAgICAgICAgICBtZXRob2QgIT09IFwiY29uc3RydWN0b3JcIlxyXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXHJcbiAgICAgICAgICAgICkge1xyXG4gICAgICAgICAgICAgICAgKHRoaXMgYXMgYW55KVttZXRob2RdID0gKHRoaXMgYXMgYW55KVttZXRob2RdLmJpbmQodGhpcyk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBHZXRzIHRoZSBzcGVjaWZpZWQgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyB0byBnZXQuXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY29ycmVzcG9uZGluZyB3aW5kb3cgb2JqZWN0LlxyXG4gICAgICovXHJcbiAgICBHZXQobmFtZTogc3RyaW5nKTogV2luZG93IHtcclxuICAgICAgICByZXR1cm4gbmV3IFdpbmRvdyhuYW1lKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgUG9zaXRpb24oKTogUHJvbWlzZTxQb3NpdGlvbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUG9zaXRpb25NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ2VudGVycyB0aGUgd2luZG93IG9uIHRoZSBzY3JlZW4uXHJcbiAgICAgKi9cclxuICAgIENlbnRlcigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENlbnRlck1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDbG9zZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgQ2xvc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShDbG9zZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBEaXNhYmxlcyBtaW4vbWF4IHNpemUgY29uc3RyYWludHMuXHJcbiAgICAgKi9cclxuICAgIERpc2FibGVTaXplQ29uc3RyYWludHMoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEVuYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxyXG4gICAgICovXHJcbiAgICBFbmFibGVTaXplQ29uc3RyYWludHMoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRm9jdXNlcyB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBGb2N1cygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEZvY3VzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEZvcmNlcyB0aGUgd2luZG93IHRvIHJlbG9hZCB0aGUgcGFnZSBhc3NldHMuXHJcbiAgICAgKi9cclxuICAgIEZvcmNlUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oRm9yY2VSZWxvYWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU3dpdGNoZXMgdGhlIHdpbmRvdyB0byBmdWxsc2NyZWVuIG1vZGUuXHJcbiAgICAgKi9cclxuICAgIEZ1bGxzY3JlZW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHNjcmVlbiB0aGF0IHRoZSB3aW5kb3cgaXMgb24uXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvbi5cclxuICAgICAqL1xyXG4gICAgR2V0U2NyZWVuKCk6IFByb21pc2U8U2NyZWVuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRTY3JlZW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgY3VycmVudCB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgem9vbSBsZXZlbC5cclxuICAgICAqL1xyXG4gICAgR2V0Wm9vbSgpOiBQcm9taXNlPG51bWJlcj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0Wm9vbU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgSGVpZ2h0KCk6IFByb21pc2U8bnVtYmVyPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIZWlnaHRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogSGlkZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKEhpZGVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgZm9jdXNlZC5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZvY3VzZWQuXHJcbiAgICAgKi9cclxuICAgIElzRm9jdXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRm9jdXNlZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmdWxsc2NyZWVuLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgZnVsbHNjcmVlbi5cclxuICAgICAqL1xyXG4gICAgSXNGdWxsc2NyZWVuKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNGdWxsc2NyZWVuTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1heGltaXNlZC5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IG1heGltaXNlZC5cclxuICAgICAqL1xyXG4gICAgSXNNYXhpbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc01heGltaXNlZE1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBtaW5pbWlzZWQuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtaW5pbWlzZWQuXHJcbiAgICAgKi9cclxuICAgIElzTWluaW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNaW5pbWlzZWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogTWF4aW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIE1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogTWluaW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIE1pbmltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWluaW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIE5hbWUoKTogUHJvbWlzZTxzdHJpbmc+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE5hbWVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogT3BlbnMgdGhlIGRldmVsb3BtZW50IHRvb2xzIHBhbmUuXHJcbiAgICAgKi9cclxuICAgIE9wZW5EZXZUb29scygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE9wZW5EZXZUb29sc01ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBSZXR1cm5zIHRoZSByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93IHRvIHRoZSBzY3JlZW4uXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgUmVsYXRpdmVQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbG9hZHMgdGhlIHBhZ2UgYXNzZXRzLlxyXG4gICAgICovXHJcbiAgICBSZWxvYWQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZWxvYWRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgcmVzaXphYmxlLlxyXG4gICAgICovXHJcbiAgICBSZXNpemFibGUoKTogUHJvbWlzZTxib29sZWFuPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShSZXNpemFibGVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmVzdG9yZXMgdGhlIHdpbmRvdyB0byBpdHMgcHJldmlvdXMgc3RhdGUgaWYgaXQgd2FzIHByZXZpb3VzbHkgbWluaW1pc2VkLCBtYXhpbWlzZWQgb3IgZnVsbHNjcmVlbi5cclxuICAgICAqL1xyXG4gICAgUmVzdG9yZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc3RvcmVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0UG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFBvc2l0aW9uTWV0aG9kLCB7IHgsIHkgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSB3aW5kb3cgdG8gYmUgYWx3YXlzIG9uIHRvcC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gYWx3YXlzT25Ub3AgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIHN0YXkgb24gdG9wLlxyXG4gICAgICovXHJcbiAgICBTZXRBbHdheXNPblRvcChhbHdheXNPblRvcDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0QWx3YXlzT25Ub3BNZXRob2QsIHsgYWx3YXlzT25Ub3AgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSByIC0gVGhlIGRlc2lyZWQgcmVkIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXHJcbiAgICAgKiBAcGFyYW0gZyAtIFRoZSBkZXNpcmVkIGdyZWVuIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXHJcbiAgICAgKiBAcGFyYW0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cclxuICAgICAqIEBwYXJhbSBhIC0gVGhlIGRlc2lyZWQgYWxwaGEgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cclxuICAgICAqL1xyXG4gICAgU2V0QmFja2dyb3VuZENvbG91cihyOiBudW1iZXIsIGc6IG51bWJlciwgYjogbnVtYmVyLCBhOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QsIHsgciwgZywgYiwgYSB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlbW92ZXMgdGhlIHdpbmRvdyBmcmFtZSBhbmQgdGl0bGUgYmFyLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBmcmFtZWxlc3MgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGZyYW1lbGVzcy5cclxuICAgICAqL1xyXG4gICAgU2V0RnJhbWVsZXNzKGZyYW1lbGVzczogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0RnJhbWVsZXNzTWV0aG9kLCB7IGZyYW1lbGVzcyB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIERpc2FibGVzIHRoZSBzeXN0ZW0gZnVsbHNjcmVlbiBidXR0b24uXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIGVuYWJsZWQgLSBXaGV0aGVyIHRoZSBmdWxsc2NyZWVuIGJ1dHRvbiBzaG91bGQgYmUgZW5hYmxlZC5cclxuICAgICAqL1xyXG4gICAgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWRNZXRob2QsIHsgZW5hYmxlZCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIG1heGltdW0gc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1heGltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBtYXhpbXVtIGhlaWdodCBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRNYXhTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRNYXhTaXplTWV0aG9kLCB7IHdpZHRoLCBoZWlnaHQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtaW5pbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWluaW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0TWluU2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWluU2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdyB0byB0aGUgc2NyZWVuLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICogQHBhcmFtIHkgLSBUaGUgZGVzaXJlZCB2ZXJ0aWNhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxyXG4gICAgICovXHJcbiAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kLCB7IHgsIHkgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHdoZXRoZXIgdGhlIHdpbmRvdyBpcyByZXNpemFibGUuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHJlc2l6YWJsZSAtIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgcmVzaXphYmxlLlxyXG4gICAgICovXHJcbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRSZXNpemFibGVNZXRob2QsIHsgcmVzaXphYmxlIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2V0cyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFNldFNpemUod2lkdGg6IG51bWJlciwgaGVpZ2h0OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHRpdGxlIC0gVGhlIGRlc2lyZWQgdGl0bGUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2V0VGl0bGUodGl0bGU6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0VGl0bGVNZXRob2QsIHsgdGl0bGUgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBTZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHBhcmFtIHpvb20gLSBUaGUgZGVzaXJlZCB6b29tIGxldmVsLlxyXG4gICAgICovXHJcbiAgICBTZXRab29tKHpvb206IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0Wm9vbU1ldGhvZCwgeyB6b29tIH0pO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2hvdygpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNob3dNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogUmV0dXJucyB0aGUgc2l6ZSBvZiB0aGUgd2luZG93LlxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHNpemUgb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgU2l6ZSgpOiBQcm9taXNlPFNpemU+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNpemVNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVG9nZ2xlcyB0aGUgd2luZG93IGJldHdlZW4gZnVsbHNjcmVlbiBhbmQgbm9ybWFsLlxyXG4gICAgICovXHJcbiAgICBUb2dnbGVGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBtYXhpbWlzZWQgYW5kIG5vcm1hbC5cclxuICAgICAqL1xyXG4gICAgVG9nZ2xlTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShUb2dnbGVNYXhpbWlzZU1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmcmFtZWxlc3MgYW5kIG5vcm1hbC5cclxuICAgICAqL1xyXG4gICAgVG9nZ2xlRnJhbWVsZXNzKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnJhbWVsZXNzTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFVuLWZ1bGxzY3JlZW5zIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFVuRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuRnVsbHNjcmVlbk1ldGhvZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBVbi1tYXhpbWlzZXMgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgVW5NYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFVuTWF4aW1pc2VNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFVuTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1pbmltaXNlTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJldHVybnMgdGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgV2lkdGgoKTogUHJvbWlzZTxudW1iZXI+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFdpZHRoTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFpvb21zIHRoZSB3aW5kb3cuXHJcbiAgICAgKi9cclxuICAgIFpvb20oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tTWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIEluY3JlYXNlcyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICovXHJcbiAgICBab29tSW4oKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tSW5NZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXHJcbiAgICAgKi9cclxuICAgIFpvb21PdXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tT3V0TWV0aG9kKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFJlc2V0cyB0aGUgem9vbSBsZXZlbCBvZiB0aGUgd2VidmlldyBjb250ZW50LlxyXG4gICAgICovXHJcbiAgICBab29tUmVzZXQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tUmVzZXRNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogSGFuZGxlcyBmaWxlIGRyb3BzIG9yaWdpbmF0aW5nIGZyb20gcGxhdGZvcm0tc3BlY2lmaWMgY29kZSAoZS5nLiwgbWFjT1MvTGludXggbmF0aXZlIGRyYWctYW5kLWRyb3ApLlxyXG4gICAgICogR2F0aGVycyBpbmZvcm1hdGlvbiBhYm91dCB0aGUgZHJvcCB0YXJnZXQgZWxlbWVudCBhbmQgc2VuZHMgaXQgYmFjayB0byB0aGUgR28gYmFja2VuZC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZmlsZW5hbWVzIC0gQW4gYXJyYXkgb2YgZmlsZSBwYXRocyAoc3RyaW5ncykgdGhhdCB3ZXJlIGRyb3BwZWQuXHJcbiAgICAgKiBAcGFyYW0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQgKENTUyBwaXhlbHMpLlxyXG4gICAgICogQHBhcmFtIHkgLSBUaGUgeS1jb29yZGluYXRlIG9mIHRoZSBkcm9wIGV2ZW50IChDU1MgcGl4ZWxzKS5cclxuICAgICAqL1xyXG4gICAgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xyXG4gICAgICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBpZ25vcmUgdGhlIGRyb3BcclxuICAgICAgICB9XHJcbiAgICAgICAgXHJcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XHJcbiAgICAgICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KGVsZW1lbnQpO1xyXG5cclxuICAgICAgICBpZiAoIWRyb3BUYXJnZXQpIHtcclxuICAgICAgICAgICAgLy8gRHJvcCB3YXMgbm90IG9uIGEgZGVzaWduYXRlZCBkcm9wIHRhcmdldCAtIGlnbm9yZVxyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBjb25zdCBlbGVtZW50RGV0YWlscyA9IHtcclxuICAgICAgICAgICAgaWQ6IGRyb3BUYXJnZXQuaWQsXHJcbiAgICAgICAgICAgIGNsYXNzTGlzdDogQXJyYXkuZnJvbShkcm9wVGFyZ2V0LmNsYXNzTGlzdCksXHJcbiAgICAgICAgICAgIGF0dHJpYnV0ZXM6IHt9IGFzIHsgW2tleTogc3RyaW5nXTogc3RyaW5nIH0sXHJcbiAgICAgICAgfTtcclxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IGRyb3BUYXJnZXQuYXR0cmlidXRlcy5sZW5ndGg7IGkrKykge1xyXG4gICAgICAgICAgICBjb25zdCBhdHRyID0gZHJvcFRhcmdldC5hdHRyaWJ1dGVzW2ldO1xyXG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscy5hdHRyaWJ1dGVzW2F0dHIubmFtZV0gPSBhdHRyLnZhbHVlO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcclxuICAgICAgICAgICAgZmlsZW5hbWVzLFxyXG4gICAgICAgICAgICB4LFxyXG4gICAgICAgICAgICB5LFxyXG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscyxcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICB0aGlzW2NhbGxlclN5bV0oRmlsZXNEcm9wcGVkLCBwYXlsb2FkKTtcclxuICAgICAgICBcclxuICAgICAgICAvLyBDbGVhbiB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhZnRlciBkcm9wXHJcbiAgICAgICAgY2xlYW51cE5hdGl2ZURyYWcoKTtcclxuICAgIH1cclxuICBcclxuICAgIC8qIFRyaWdnZXJzIFdpbmRvd3MgMTEgU25hcCBBc3Npc3QgZmVhdHVyZSAoV2luZG93cyBvbmx5KS5cclxuICAgICAqIFRoaXMgaXMgZXF1aXZhbGVudCB0byBwcmVzc2luZyBXaW4rWiBhbmQgc2hvd3Mgc25hcCBsYXlvdXQgb3B0aW9ucy5cclxuICAgICAqL1xyXG4gICAgU25hcEFzc2lzdCgpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNuYXBBc3Npc3RNZXRob2QpO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogT3BlbnMgdGhlIHByaW50IGRpYWxvZyBmb3IgdGhlIHdpbmRvdy5cclxuICAgICAqL1xyXG4gICAgUHJpbnQoKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQcmludE1ldGhvZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBUaGUgd2luZG93IHdpdGhpbiB3aGljaCB0aGUgc2NyaXB0IGlzIHJ1bm5pbmcuXHJcbiAqL1xyXG5jb25zdCB0aGlzV2luZG93ID0gbmV3IFdpbmRvdygnJyk7XHJcblxyXG4vKipcclxuICogU2V0cyB1cCBnbG9iYWwgZHJhZyBhbmQgZHJvcCBldmVudCBsaXN0ZW5lcnMgZm9yIGZpbGUgZHJvcHMuXHJcbiAqIEhhbmRsZXMgdmlzdWFsIGZlZWRiYWNrIChob3ZlciBzdGF0ZSkgYW5kIGZpbGUgZHJvcCBwcm9jZXNzaW5nLlxyXG4gKi9cclxuZnVuY3Rpb24gc2V0dXBEcm9wVGFyZ2V0TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZG9jRWxlbWVudCA9IGRvY3VtZW50LmRvY3VtZW50RWxlbWVudDtcclxuICAgIGxldCBkcmFnRW50ZXJDb3VudGVyID0gMDtcclxuXHJcbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdlbnRlcicsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xyXG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJhZ3MsIGxldCBvdGhlciBkcmFncyBwYXNzIHRocm91Z2hcclxuICAgICAgICB9XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxyXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7IC8vIFNob3cgXCJubyBkcm9wXCIgY3Vyc29yXHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgZG9uJ3Qgc2hvdyBob3ZlciBlZmZlY3RzXHJcbiAgICAgICAgfVxyXG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIrKztcclxuICAgICAgICBcclxuICAgICAgICBjb25zdCB0YXJnZXRFbGVtZW50ID0gZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZKTtcclxuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XHJcblxyXG4gICAgICAgIC8vIFVwZGF0ZSBob3ZlciBzdGF0ZVxyXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBkcm9wVGFyZ2V0LmNsYXNzTGlzdC5hZGQoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xyXG4gICAgICAgIH1cclxuICAgIH0sIGZhbHNlKTtcclxuXHJcbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2RyYWdvdmVyJywgKGV2ZW50KSA9PiB7XHJcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjsgLy8gT25seSBoYW5kbGUgZmlsZSBkcmFnc1xyXG4gICAgICAgIH1cclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJzsgLy8gU2hvdyBcIm5vIGRyb3BcIiBjdXJzb3JcclxuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcclxuICAgICAgICB9XHJcbiAgICAgICAgXHJcbiAgICAgICAgLy8gVXBkYXRlIGRyb3AgdGFyZ2V0IGFzIGN1cnNvciBtb3Zlc1xyXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xyXG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudCh0YXJnZXRFbGVtZW50KTtcclxuICAgICAgICBcclxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQgJiYgY3VycmVudERyb3BUYXJnZXQgIT09IGRyb3BUYXJnZXQpIHtcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBcclxuICAgICAgICBpZiAoZHJvcFRhcmdldCkge1xyXG4gICAgICAgICAgICBpZiAoIWRyb3BUYXJnZXQuY2xhc3NMaXN0LmNvbnRhaW5zKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUykpIHtcclxuICAgICAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ2NvcHknO1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7XHJcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcclxuICAgICAgICB9XHJcbiAgICB9LCBmYWxzZSk7XHJcblxyXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCAoZXZlbnQpID0+IHtcclxuICAgICAgICBpZiAoIWV2ZW50LmRhdGFUcmFuc2Zlcj8udHlwZXMuaW5jbHVkZXMoJ0ZpbGVzJykpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXHJcbiAgICAgICAgLy8gT24gV2luZG93cywgY2hlY2sgaWYgZmlsZSBkcm9wcyBhcmUgZW5hYmxlZCBmb3IgdGhpcyB3aW5kb3dcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIFxyXG4gICAgICAgIC8vIE9uIExpbnV4L1dlYktpdEdUSyBhbmQgbWFjT1MsIGRyYWdsZWF2ZSBmaXJlcyBpbW1lZGlhdGVseSB3aXRoIHJlbGF0ZWRUYXJnZXQ9bnVsbCB3aGVuIG5hdGl2ZVxyXG4gICAgICAgIC8vIGRyYWcgaGFuZGxpbmcgaXMgaW52b2x2ZWQuIElnbm9yZSB0aGVzZSBzcHVyaW91cyBldmVudHMgLSB3ZSdsbCBjbGVhbiB1cCBvbiBkcm9wIGluc3RlYWQuXHJcbiAgICAgICAgaWYgKGV2ZW50LnJlbGF0ZWRUYXJnZXQgPT09IG51bGwpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBcclxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyLS07XHJcbiAgICAgICAgXHJcbiAgICAgICAgaWYgKGRyYWdFbnRlckNvdW50ZXIgPT09IDAgfHwgXHJcbiAgICAgICAgICAgIChjdXJyZW50RHJvcFRhcmdldCAmJiAhY3VycmVudERyb3BUYXJnZXQuY29udGFpbnMoZXZlbnQucmVsYXRlZFRhcmdldCBhcyBOb2RlKSkpIHtcclxuICAgICAgICAgICAgaWYgKGN1cnJlbnREcm9wVGFyZ2V0KSB7XHJcbiAgICAgICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XHJcbiAgICAgICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG5cclxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJvcCcsIChldmVudCkgPT4ge1xyXG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xyXG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJvcHNcclxuICAgICAgICB9XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxyXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XHJcbiAgICAgICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgICAgICBkcmFnRW50ZXJDb3VudGVyID0gMDtcclxuICAgICAgICBcclxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcclxuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xyXG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IG51bGw7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBoYW5kbGUgZmlsZSBkcm9wcyB2aWEgSmF2YVNjcmlwdFxyXG4gICAgICAgIC8vIE9uIG1hY09TL0xpbnV4LCBuYXRpdmUgY29kZSB3aWxsIGNhbGwgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcFxyXG4gICAgICAgIGlmIChjYW5SZXNvbHZlRmlsZVBhdGhzKCkpIHtcclxuICAgICAgICAgICAgY29uc3QgZmlsZXM6IEZpbGVbXSA9IFtdO1xyXG4gICAgICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGNvbnN0IGl0ZW0gb2YgZXZlbnQuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKGl0ZW0ua2luZCA9PT0gJ2ZpbGUnKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IGZpbGUgPSBpdGVtLmdldEFzRmlsZSgpO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZmlsZSkgZmlsZXMucHVzaChmaWxlKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0gZWxzZSBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyLmZpbGVzKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGNvbnN0IGZpbGUgb2YgZXZlbnQuZGF0YVRyYW5zZmVyLmZpbGVzKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgZmlsZXMucHVzaChmaWxlKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBcclxuICAgICAgICAgICAgaWYgKGZpbGVzLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgICAgIHJlc29sdmVGaWxlUGF0aHMoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZmlsZXMpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG4gICAgfSwgZmFsc2UpO1xyXG59XHJcblxyXG4vLyBJbml0aWFsaXplIGxpc3RlbmVycyB3aGVuIHRoZSBzY3JpcHQgbG9hZHNcclxuaWYgKHR5cGVvZiB3aW5kb3cgIT09IFwidW5kZWZpbmVkXCIgJiYgdHlwZW9mIGRvY3VtZW50ICE9PSBcInVuZGVmaW5lZFwiKSB7XHJcbiAgICBzZXR1cERyb3BUYXJnZXRMaXN0ZW5lcnMoKTtcclxufVxyXG5cclxuZXhwb3J0IGRlZmF1bHQgdGhpc1dpbmRvdztcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XHJcblxyXG4vLyBOT1RFOiB0aGUgZm9sbG93aW5nIG1ldGhvZHMgTVVTVCBiZSBpbXBvcnRlZCBleHBsaWNpdGx5IGJlY2F1c2Ugb2YgaG93IGVzYnVpbGQgaW5qZWN0aW9uIHdvcmtzXHJcbmltcG9ydCB7IEVuYWJsZSBhcyBFbmFibGVXTUwgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XHJcbmltcG9ydCB7IGRlYnVnTG9nIH0gZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3V0aWxzXCI7XHJcblxyXG53aW5kb3cud2FpbHMgPSBSdW50aW1lO1xyXG5FbmFibGVXTUwoKTtcclxuXHJcbmlmIChERUJVRykge1xyXG4gICAgZGVidWdMb2coXCJXYWlscyBSdW50aW1lIExvYWRlZFwiKVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlN5c3RlbSk7XHJcblxyXG5jb25zdCBTeXN0ZW1Jc0RhcmtNb2RlID0gMDtcclxuY29uc3QgU3lzdGVtRW52aXJvbm1lbnQgPSAxO1xyXG5jb25zdCBTeXN0ZW1DYXBhYmlsaXRpZXMgPSAyO1xyXG5cclxuY29uc3QgX2ludm9rZSA9IChmdW5jdGlvbiAoKSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIC8vIFdpbmRvd3MgV2ViVmlldzJcclxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2UpIHtcclxuICAgICAgICAgICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5jaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZS5iaW5kKCh3aW5kb3cgYXMgYW55KS5jaHJvbWUud2Vidmlldyk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIC8vIG1hY09TL2lPUyBXS1dlYlZpZXdcclxuICAgICAgICBlbHNlIGlmICgod2luZG93IGFzIGFueSkud2Via2l0Py5tZXNzYWdlSGFuZGxlcnM/LlsnZXh0ZXJuYWwnXT8ucG9zdE1lc3NhZ2UpIHtcclxuICAgICAgICAgICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS53ZWJraXQubWVzc2FnZUhhbmRsZXJzWydleHRlcm5hbCddLnBvc3RNZXNzYWdlLmJpbmQoKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10pO1xyXG4gICAgICAgIH1cclxuICAgICAgICAvLyBBbmRyb2lkIFdlYlZpZXcgLSB1c2VzIGFkZEphdmFzY3JpcHRJbnRlcmZhY2Ugd2hpY2ggZXhwb3NlcyB3aW5kb3cud2FpbHMuaW52b2tlXHJcbiAgICAgICAgZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndhaWxzPy5pbnZva2UpIHtcclxuICAgICAgICAgICAgcmV0dXJuIChtc2c6IGFueSkgPT4gKHdpbmRvdyBhcyBhbnkpLndhaWxzLmludm9rZSh0eXBlb2YgbXNnID09PSAnc3RyaW5nJyA/IG1zZyA6IEpTT04uc3RyaW5naWZ5KG1zZykpO1xyXG4gICAgICAgIH1cclxuICAgIH0gY2F0Y2goZSkge31cclxuXHJcbiAgICBjb25zb2xlLndhcm4oJ1xcbiVjXHUyNkEwXHVGRTBGIEJyb3dzZXIgRW52aXJvbm1lbnQgRGV0ZWN0ZWQgJWNcXG5cXG4lY09ubHkgVUkgcHJldmlld3MgYXJlIGF2YWlsYWJsZSBpbiB0aGUgYnJvd3Nlci4gRm9yIGZ1bGwgZnVuY3Rpb25hbGl0eSwgcGxlYXNlIHJ1biB0aGUgYXBwbGljYXRpb24gaW4gZGVza3RvcCBtb2RlLlxcbk1vcmUgaW5mb3JtYXRpb24gYXQ6IGh0dHBzOi8vdjMud2FpbHMuaW8vbGVhcm4vYnVpbGQvI3VzaW5nLWEtYnJvd3Nlci1mb3ItZGV2ZWxvcG1lbnRcXG4nLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjZmZmZmZmOyBjb2xvcjogIzAwMDAwMDsgZm9udC13ZWlnaHQ6IGJvbGQ7IHBhZGRpbmc6IDRweCA4cHg7IGJvcmRlci1yYWRpdXM6IDRweDsgYm9yZGVyOiAycHggc29saWQgIzAwMDAwMDsnLFxyXG4gICAgICAgICdiYWNrZ3JvdW5kOiB0cmFuc3BhcmVudDsnLFxyXG4gICAgICAgICdjb2xvcjogI2ZmZmZmZjsgZm9udC1zdHlsZTogaXRhbGljOyBmb250LXdlaWdodDogYm9sZDsnKTtcclxuICAgIHJldHVybiBudWxsO1xyXG59KSgpO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGludm9rZShtc2c6IGFueSk6IHZvaWQge1xyXG4gICAgX2ludm9rZT8uKG1zZyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtSXNEYXJrTW9kZSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBGZXRjaGVzIHRoZSBjYXBhYmlsaXRpZXMgb2YgdGhlIGFwcGxpY2F0aW9uIGZyb20gdGhlIHNlcnZlci5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cclxuICovXHJcbmV4cG9ydCBhc3luYyBmdW5jdGlvbiBDYXBhYmlsaXRpZXMoKTogUHJvbWlzZTxSZWNvcmQ8c3RyaW5nLCBhbnk+PiB7XHJcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1DYXBhYmlsaXRpZXMpO1xyXG59XHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIE9TSW5mbyB7XHJcbiAgICAvKiogVGhlIGJyYW5kaW5nIG9mIHRoZSBPUy4gKi9cclxuICAgIEJyYW5kaW5nOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIElEIG9mIHRoZSBPUy4gKi9cclxuICAgIElEOiBzdHJpbmc7XHJcbiAgICAvKiogVGhlIG5hbWUgb2YgdGhlIE9TLiAqL1xyXG4gICAgTmFtZTogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSB2ZXJzaW9uIG9mIHRoZSBPUy4gKi9cclxuICAgIFZlcnNpb246IHN0cmluZztcclxufVxyXG5cclxuZXhwb3J0IGludGVyZmFjZSBFbnZpcm9ubWVudEluZm8ge1xyXG4gICAgLyoqIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS4gKi9cclxuICAgIEFyY2g6IHN0cmluZztcclxuICAgIC8qKiBUcnVlIGlmIHRoZSBhcHBsaWNhdGlvbiBpcyBydW5uaW5nIGluIGRlYnVnIG1vZGUsIG90aGVyd2lzZSBmYWxzZS4gKi9cclxuICAgIERlYnVnOiBib29sZWFuO1xyXG4gICAgLyoqIFRoZSBvcGVyYXRpbmcgc3lzdGVtIGluIHVzZS4gKi9cclxuICAgIE9TOiBzdHJpbmc7XHJcbiAgICAvKiogRGV0YWlscyBvZiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cclxuICAgIE9TSW5mbzogT1NJbmZvO1xyXG4gICAgLyoqIEFkZGl0aW9uYWwgcGxhdGZvcm0gaW5mb3JtYXRpb24uICovXHJcbiAgICBQbGF0Zm9ybUluZm86IFJlY29yZDxzdHJpbmcsIGFueT47XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgZW52aXJvbm1lbnQgZGV0YWlscy5cclxuICpcclxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKTogUHJvbWlzZTxFbnZpcm9ubWVudEluZm8+IHtcclxuICAgIHJldHVybiBjYWxsKFN5c3RlbUVudmlyb25tZW50KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBvcGVyYXRpbmcgc3lzdGVtIGlzIFdpbmRvd3MuXHJcbiAqXHJcbiAqIEByZXR1cm4gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJ3aW5kb3dzXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eC5cclxuICpcclxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0xpbnV4KCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJsaW51eFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc01hYygpOiBib29sZWFuIHtcclxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwiZGFyd2luXCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgYXJjaGl0ZWN0dXJlIGlzIEFNRDY0LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYW1kNjRcIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBhcmNoaXRlY3R1cmUgaXMgQVJNLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTSgpOiBib29sZWFuIHtcclxuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhcm1cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUuXHJcbiAqXHJcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCk6IGJvb2xlYW4ge1xyXG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFybTY0XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXBvcnRzIHdoZXRoZXIgdGhlIGFwcCBpcyBiZWluZyBydW4gaW4gZGVidWcgbW9kZS5cclxuICpcclxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzRGVidWcoKTogYm9vbGVhbiB7XHJcbiAgICByZXR1cm4gQm9vbGVhbigod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uRGVidWcpO1xyXG59XHJcblxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB7IElzRGVidWcgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcclxuaW1wb3J0IHsgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xyXG5cclxuLy8gc2V0dXBcclxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcclxuXHJcbmNvbnN0IENvbnRleHRNZW51T3BlbiA9IDA7XHJcblxyXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQ6IHN0cmluZywgeDogbnVtYmVyLCB5OiBudW1iZXIsIGRhdGE6IGFueSk6IHZvaWQge1xyXG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudDogTW91c2VFdmVudCkge1xyXG4gICAgY29uc3QgdGFyZ2V0ID0gZXZlbnRUYXJnZXQoZXZlbnQpO1xyXG5cclxuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XHJcbiAgICBjb25zdCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpLnRyaW0oKTtcclxuXHJcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcclxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgIGNvbnN0IGRhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xyXG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgZGF0YSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQsIHRhcmdldCk7XHJcbiAgICB9XHJcbn1cclxuXHJcblxyXG4vKlxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuXHJcblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXHJcbiovXHJcbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQ6IE1vdXNlRXZlbnQsIHRhcmdldDogSFRNTEVsZW1lbnQpIHtcclxuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxyXG4gICAgaWYgKElzRGVidWcoKSkge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICBzd2l0Y2ggKHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCkpIHtcclxuICAgICAgICBjYXNlICdzaG93JzpcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGNhc2UgJ2hpZGUnOlxyXG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcclxuICAgIGlmICh0YXJnZXQuaXNDb250ZW50RWRpdGFibGUpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxyXG4gICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xyXG4gICAgY29uc3QgaGFzU2VsZWN0aW9uID0gc2VsZWN0aW9uICYmIHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDA7XHJcbiAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XHJcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XHJcbiAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XHJcbiAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcclxuICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xyXG4gICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xyXG4gICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IHRhcmdldCkge1xyXG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICAvLyBDaGVjayBpZiB0YWcgaXMgaW5wdXQgb3IgdGV4dGFyZWEuXHJcbiAgICBpZiAodGFyZ2V0IGluc3RhbmNlb2YgSFRNTElucHV0RWxlbWVudCB8fCB0YXJnZXQgaW5zdGFuY2VvZiBIVE1MVGV4dEFyZWFFbGVtZW50KSB7XHJcbiAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIXRhcmdldC5yZWFkT25seSAmJiAhdGFyZ2V0LmRpc2FibGVkKSkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qKlxyXG4gKiBSZXRyaWV2ZXMgdGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleSBmcm9tIHRoZSBmbGFnIG1hcC5cclxuICpcclxuICogQHBhcmFtIGtleSAtIFRoZSBrZXkgdG8gcmV0cmlldmUgdGhlIHZhbHVlIGZvci5cclxuICogQHJldHVybiBUaGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5OiBzdHJpbmcpOiBhbnkge1xyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gd2luZG93Ll93YWlscy5mbGFnc1trZXldO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleSArIFwiJzogXCIgKyBlLCB7IGNhdXNlOiBlIH0pO1xyXG4gICAgfVxyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5pbXBvcnQgeyBpbnZva2UsIElzV2luZG93cyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xyXG5pbXBvcnQgeyBHZXRGbGFnIH0gZnJvbSBcIi4vZmxhZ3MuanNcIjtcclxuaW1wb3J0IHsgY2FuVHJhY2tCdXR0b25zLCBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XHJcblxyXG4vLyBTZXR1cFxyXG5sZXQgY2FuRHJhZyA9IGZhbHNlO1xyXG5sZXQgZHJhZ2dpbmcgPSBmYWxzZTtcclxuXHJcbmxldCByZXNpemFibGUgPSBmYWxzZTtcclxubGV0IGNhblJlc2l6ZSA9IGZhbHNlO1xyXG5sZXQgcmVzaXppbmcgPSBmYWxzZTtcclxubGV0IHJlc2l6ZUVkZ2U6IHN0cmluZyA9IFwiXCI7XHJcbmxldCBkZWZhdWx0Q3Vyc29yID0gXCJhdXRvXCI7XHJcblxyXG5sZXQgYnV0dG9ucyA9IDA7XHJcbmNvbnN0IGJ1dHRvbnNUcmFja2VkID0gY2FuVHJhY2tCdXR0b25zKCk7XHJcblxyXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcclxud2luZG93Ll93YWlscy5zZXRSZXNpemFibGUgPSAodmFsdWU6IGJvb2xlYW4pOiB2b2lkID0+IHtcclxuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xyXG4gICAgaWYgKCFyZXNpemFibGUpIHtcclxuICAgICAgICAvLyBTdG9wIHJlc2l6aW5nIGlmIGluIHByb2dyZXNzLlxyXG4gICAgICAgIGNhblJlc2l6ZSA9IHJlc2l6aW5nID0gZmFsc2U7XHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9XHJcbn07XHJcblxyXG4vLyBEZWZlciBhdHRhY2hpbmcgbW91c2UgbGlzdGVuZXJzIHVudGlsIHdlIGtub3cgd2UncmUgbm90IG9uIG1vYmlsZS5cclxubGV0IGRyYWdJbml0RG9uZSA9IGZhbHNlO1xyXG5mdW5jdGlvbiBpc01vYmlsZSgpOiBib29sZWFuIHtcclxuICAgIGNvbnN0IG9zID0gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/Lk9TO1xyXG4gICAgaWYgKG9zID09PSBcImlvc1wiIHx8IG9zID09PSBcImFuZHJvaWRcIikgcmV0dXJuIHRydWU7XHJcbiAgICAvLyBGYWxsYmFjayBoZXVyaXN0aWMgaWYgZW52aXJvbm1lbnQgbm90IHlldCBzZXRcclxuICAgIGNvbnN0IHVhID0gbmF2aWdhdG9yLnVzZXJBZ2VudCB8fCBuYXZpZ2F0b3IudmVuZG9yIHx8ICh3aW5kb3cgYXMgYW55KS5vcGVyYSB8fCBcIlwiO1xyXG4gICAgcmV0dXJuIC9hbmRyb2lkfGlwaG9uZXxpcGFkfGlwb2R8aWVtb2JpbGV8d3BkZXNrdG9wL2kudGVzdCh1YSk7XHJcbn1cclxuZnVuY3Rpb24gdHJ5SW5pdERyYWdIYW5kbGVycygpOiB2b2lkIHtcclxuICAgIGlmIChkcmFnSW5pdERvbmUpIHJldHVybjtcclxuICAgIGlmIChpc01vYmlsZSgpKSByZXR1cm47XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgdXBkYXRlLCB7IGNhcHR1cmU6IHRydWUgfSk7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xyXG4gICAgZm9yIChjb25zdCBldiBvZiBbJ2NsaWNrJywgJ2NvbnRleHRtZW51JywgJ2RibGNsaWNrJ10pIHtcclxuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcihldiwgc3VwcHJlc3NFdmVudCwgeyBjYXB0dXJlOiB0cnVlIH0pO1xyXG4gICAgfVxyXG4gICAgZHJhZ0luaXREb25lID0gdHJ1ZTtcclxufVxyXG4vLyBBdHRlbXB0IGltbWVkaWF0ZSBpbml0IChpbiBjYXNlIGVudmlyb25tZW50IGFscmVhZHkgcHJlc2VudClcclxudHJ5SW5pdERyYWdIYW5kbGVycygpO1xyXG4vLyBBbHNvIGF0dGVtcHQgb24gRE9NIHJlYWR5XHJcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCB0cnlJbml0RHJhZ0hhbmRsZXJzLCB7IG9uY2U6IHRydWUgfSk7XHJcbi8vIEFzIGEgbGFzdCByZXNvcnQsIHBvbGwgZm9yIGVudmlyb25tZW50IGZvciBhIHNob3J0IHBlcmlvZFxyXG5sZXQgZHJhZ0VudlBvbGxzID0gMDtcclxuY29uc3QgZHJhZ0VudlBvbGwgPSB3aW5kb3cuc2V0SW50ZXJ2YWwoKCkgPT4ge1xyXG4gICAgaWYgKGRyYWdJbml0RG9uZSkgeyB3aW5kb3cuY2xlYXJJbnRlcnZhbChkcmFnRW52UG9sbCk7IHJldHVybjsgfVxyXG4gICAgdHJ5SW5pdERyYWdIYW5kbGVycygpO1xyXG4gICAgaWYgKCsrZHJhZ0VudlBvbGxzID4gMTAwKSB7IHdpbmRvdy5jbGVhckludGVydmFsKGRyYWdFbnZQb2xsKTsgfVxyXG59LCA1MCk7XHJcblxyXG5mdW5jdGlvbiBzdXBwcmVzc0V2ZW50KGV2ZW50OiBFdmVudCkge1xyXG4gICAgLy8gU3VwcHJlc3MgY2xpY2sgZXZlbnRzIHdoaWxlIHJlc2l6aW5nIG9yIGRyYWdnaW5nLlxyXG4gICAgaWYgKGRyYWdnaW5nIHx8IHJlc2l6aW5nKSB7XHJcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XHJcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgIH1cclxufVxyXG5cclxuLy8gVXNlIGNvbnN0YW50cyB0byBhdm9pZCBjb21wYXJpbmcgc3RyaW5ncyBtdWx0aXBsZSB0aW1lcy5cclxuY29uc3QgTW91c2VEb3duID0gMDtcclxuY29uc3QgTW91c2VVcCAgID0gMTtcclxuY29uc3QgTW91c2VNb3ZlID0gMjtcclxuXHJcbmZ1bmN0aW9uIHVwZGF0ZShldmVudDogTW91c2VFdmVudCkge1xyXG4gICAgLy8gV2luZG93cyBzdXBwcmVzc2VzIG1vdXNlIGV2ZW50cyBhdCB0aGUgZW5kIG9mIGRyYWdnaW5nIG9yIHJlc2l6aW5nLFxyXG4gICAgLy8gc28gd2UgbmVlZCB0byBiZSBzbWFydCBhbmQgc3ludGhlc2l6ZSBidXR0b24gZXZlbnRzLlxyXG5cclxuICAgIGxldCBldmVudFR5cGU6IG51bWJlciwgZXZlbnRCdXR0b25zID0gZXZlbnQuYnV0dG9ucztcclxuICAgIHN3aXRjaCAoZXZlbnQudHlwZSkge1xyXG4gICAgICAgIGNhc2UgJ21vdXNlZG93bic6XHJcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlRG93bjtcclxuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zIHwgKDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICBjYXNlICdtb3VzZXVwJzpcclxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VVcDtcclxuICAgICAgICAgICAgaWYgKCFidXR0b25zVHJhY2tlZCkgeyBldmVudEJ1dHRvbnMgPSBidXR0b25zICYgfigxIDw8IGV2ZW50LmJ1dHRvbik7IH1cclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgZGVmYXVsdDpcclxuICAgICAgICAgICAgZXZlbnRUeXBlID0gTW91c2VNb3ZlO1xyXG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnM7IH1cclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICB9XHJcblxyXG4gICAgbGV0IHJlbGVhc2VkID0gYnV0dG9ucyAmIH5ldmVudEJ1dHRvbnM7XHJcbiAgICBsZXQgcHJlc3NlZCA9IGV2ZW50QnV0dG9ucyAmIH5idXR0b25zO1xyXG5cclxuICAgIGJ1dHRvbnMgPSBldmVudEJ1dHRvbnM7XHJcblxyXG4gICAgLy8gU3ludGhlc2l6ZSBhIHJlbGVhc2UtcHJlc3Mgc2VxdWVuY2UgaWYgd2UgZGV0ZWN0IGEgcHJlc3Mgb2YgYW4gYWxyZWFkeSBwcmVzc2VkIGJ1dHRvbi5cclxuICAgIGlmIChldmVudFR5cGUgPT09IE1vdXNlRG93biAmJiAhKHByZXNzZWQgJiBldmVudC5idXR0b24pKSB7XHJcbiAgICAgICAgcmVsZWFzZWQgfD0gKDEgPDwgZXZlbnQuYnV0dG9uKTtcclxuICAgICAgICBwcmVzc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gU3VwcHJlc3MgYWxsIGJ1dHRvbiBldmVudHMgZHVyaW5nIGRyYWdnaW5nIGFuZCByZXNpemluZyxcclxuICAgIC8vIHVubGVzcyB0aGlzIGlzIGEgbW91c2V1cCBldmVudCB0aGF0IGlzIGVuZGluZyBhIGRyYWcgYWN0aW9uLlxyXG4gICAgaWYgKFxyXG4gICAgICAgIGV2ZW50VHlwZSAhPT0gTW91c2VNb3ZlIC8vIEZhc3QgcGF0aCBmb3IgbW91c2Vtb3ZlXHJcbiAgICAgICAgJiYgcmVzaXppbmdcclxuICAgICAgICB8fCAoXHJcbiAgICAgICAgICAgIGRyYWdnaW5nXHJcbiAgICAgICAgICAgICYmIChcclxuICAgICAgICAgICAgICAgIGV2ZW50VHlwZSA9PT0gTW91c2VEb3duXHJcbiAgICAgICAgICAgICAgICB8fCBldmVudC5idXR0b24gIT09IDBcclxuICAgICAgICAgICAgKVxyXG4gICAgICAgIClcclxuICAgICkge1xyXG4gICAgICAgIGV2ZW50LnN0b3BJbW1lZGlhdGVQcm9wYWdhdGlvbigpO1xyXG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gSGFuZGxlIHJlbGVhc2VzXHJcbiAgICBpZiAocmVsZWFzZWQgJiAxKSB7IHByaW1hcnlVcChldmVudCk7IH1cclxuICAgIC8vIEhhbmRsZSBwcmVzc2VzXHJcbiAgICBpZiAocHJlc3NlZCAmIDEpIHsgcHJpbWFyeURvd24oZXZlbnQpOyB9XHJcblxyXG4gICAgLy8gSGFuZGxlIG1vdXNlbW92ZVxyXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VNb3ZlKSB7IG9uTW91c2VNb3ZlKGV2ZW50KTsgfTtcclxufVxyXG5cclxuZnVuY3Rpb24gcHJpbWFyeURvd24oZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcclxuICAgIC8vIFJlc2V0IHJlYWRpbmVzcyBzdGF0ZS5cclxuICAgIGNhbkRyYWcgPSBmYWxzZTtcclxuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xyXG5cclxuICAgIC8vIElnbm9yZSByZXBlYXRlZCBjbGlja3Mgb24gbWFjT1MgYW5kIExpbnV4LlxyXG4gICAgaWYgKCFJc1dpbmRvd3MoKSkge1xyXG4gICAgICAgIGlmIChldmVudC50eXBlID09PSAnbW91c2Vkb3duJyAmJiBldmVudC5idXR0b24gPT09IDAgJiYgZXZlbnQuZGV0YWlsICE9PSAxKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcclxuICAgICAgICAvLyBSZWFkeSB0byByZXNpemUgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZS5cclxuICAgICAgICBjYW5SZXNpemUgPSB0cnVlO1xyXG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnIG9wZXJhdGlvbnMgd2hlbiBvbiByZXNpemUgZWRnZXMuXHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIFJldHJpZXZlIHRhcmdldCBlbGVtZW50XHJcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XHJcblxyXG4gICAgLy8gUmVhZHkgdG8gZHJhZyBpZiB0aGUgcHJpbWFyeSBidXR0b24gd2FzIHByZXNzZWQgZm9yIHRoZSBmaXJzdCB0aW1lIG9uIGEgZHJhZ2dhYmxlIGVsZW1lbnQuXHJcbiAgICAvLyBJZ25vcmUgY2xpY2tzIG9uIHRoZSBzY3JvbGxiYXIuXHJcbiAgICBjb25zdCBzdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCk7XHJcbiAgICBjYW5EcmFnID0gKFxyXG4gICAgICAgIHN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLXdhaWxzLWRyYWdnYWJsZVwiKS50cmltKCkgPT09IFwiZHJhZ1wiXHJcbiAgICAgICAgJiYgKFxyXG4gICAgICAgICAgICBldmVudC5vZmZzZXRYIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nTGVmdCkgPCB0YXJnZXQuY2xpZW50V2lkdGhcclxuICAgICAgICAgICAgJiYgZXZlbnQub2Zmc2V0WSAtIHBhcnNlRmxvYXQoc3R5bGUucGFkZGluZ1RvcCkgPCB0YXJnZXQuY2xpZW50SGVpZ2h0XHJcbiAgICAgICAgKVxyXG4gICAgKTtcclxufVxyXG5cclxuZnVuY3Rpb24gcHJpbWFyeVVwKGV2ZW50OiBNb3VzZUV2ZW50KSB7XHJcbiAgICAvLyBTdG9wIGRyYWdnaW5nIGFuZCByZXNpemluZy5cclxuICAgIGNhbkRyYWcgPSBmYWxzZTtcclxuICAgIGRyYWdnaW5nID0gZmFsc2U7XHJcbiAgICBjYW5SZXNpemUgPSBmYWxzZTtcclxuICAgIHJlc2l6aW5nID0gZmFsc2U7XHJcbn1cclxuXHJcbmNvbnN0IGN1cnNvckZvckVkZ2UgPSBPYmplY3QuZnJlZXplKHtcclxuICAgIFwic2UtcmVzaXplXCI6IFwibndzZS1yZXNpemVcIixcclxuICAgIFwic3ctcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcclxuICAgIFwibnctcmVzaXplXCI6IFwibndzZS1yZXNpemVcIixcclxuICAgIFwibmUtcmVzaXplXCI6IFwibmVzdy1yZXNpemVcIixcclxuICAgIFwidy1yZXNpemVcIjogXCJldy1yZXNpemVcIixcclxuICAgIFwibi1yZXNpemVcIjogXCJucy1yZXNpemVcIixcclxuICAgIFwicy1yZXNpemVcIjogXCJucy1yZXNpemVcIixcclxuICAgIFwiZS1yZXNpemVcIjogXCJldy1yZXNpemVcIixcclxufSlcclxuXHJcbmZ1bmN0aW9uIHNldFJlc2l6ZShlZGdlPzoga2V5b2YgdHlwZW9mIGN1cnNvckZvckVkZ2UpOiB2b2lkIHtcclxuICAgIGlmIChlZGdlKSB7XHJcbiAgICAgICAgaWYgKCFyZXNpemVFZGdlKSB7IGRlZmF1bHRDdXJzb3IgPSBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvcjsgfVxyXG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gY3Vyc29yRm9yRWRnZVtlZGdlXTtcclxuICAgIH0gZWxzZSBpZiAoIWVkZ2UgJiYgcmVzaXplRWRnZSkge1xyXG4gICAgICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gZGVmYXVsdEN1cnNvcjtcclxuICAgIH1cclxuXHJcbiAgICByZXNpemVFZGdlID0gZWRnZSB8fCBcIlwiO1xyXG59XHJcblxyXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShldmVudDogTW91c2VFdmVudCk6IHZvaWQge1xyXG4gICAgaWYgKGNhblJlc2l6ZSAmJiByZXNpemVFZGdlKSB7XHJcbiAgICAgICAgLy8gU3RhcnQgcmVzaXppbmcuXHJcbiAgICAgICAgcmVzaXppbmcgPSB0cnVlO1xyXG4gICAgICAgIGludm9rZShcIndhaWxzOnJlc2l6ZTpcIiArIHJlc2l6ZUVkZ2UpO1xyXG4gICAgfSBlbHNlIGlmIChjYW5EcmFnKSB7XHJcbiAgICAgICAgLy8gU3RhcnQgZHJhZ2dpbmcuXHJcbiAgICAgICAgZHJhZ2dpbmcgPSB0cnVlO1xyXG4gICAgICAgIGludm9rZShcIndhaWxzOmRyYWdcIik7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKGRyYWdnaW5nIHx8IHJlc2l6aW5nKSB7XHJcbiAgICAgICAgLy8gRWl0aGVyIGRyYWcgb3IgcmVzaXplIGlzIG9uZ29pbmcsXHJcbiAgICAgICAgLy8gcmVzZXQgcmVhZGluZXNzIGFuZCBzdG9wIHByb2Nlc3NpbmcuXHJcbiAgICAgICAgY2FuRHJhZyA9IGNhblJlc2l6ZSA9IGZhbHNlO1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBpZiAoIXJlc2l6YWJsZSB8fCAhSXNXaW5kb3dzKCkpIHtcclxuICAgICAgICBpZiAocmVzaXplRWRnZSkgeyBzZXRSZXNpemUoKTsgfVxyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICBjb25zdCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xyXG4gICAgY29uc3QgcmVzaXplSGFuZGxlV2lkdGggPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZVdpZHRoXCIpIHx8IDU7XHJcblxyXG4gICAgLy8gRXh0cmEgcGl4ZWxzIGZvciB0aGUgY29ybmVyIGFyZWFzLlxyXG4gICAgY29uc3QgY29ybmVyRXh0cmEgPSBHZXRGbGFnKFwicmVzaXplQ29ybmVyRXh0cmFcIikgfHwgMTA7XHJcblxyXG4gICAgY29uc3QgcmlnaHRCb3JkZXIgPSAod2luZG93Lm91dGVyV2lkdGggLSBldmVudC5jbGllbnRYKSA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgY29uc3QgbGVmdEJvcmRlciA9IGV2ZW50LmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGNvbnN0IHRvcEJvcmRlciA9IGV2ZW50LmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcbiAgICBjb25zdCBib3R0b21Cb3JkZXIgPSAod2luZG93Lm91dGVySGVpZ2h0IC0gZXZlbnQuY2xpZW50WSkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcblxyXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXIgYXJlYXMuXHJcbiAgICBjb25zdCByaWdodENvcm5lciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xyXG4gICAgY29uc3QgbGVmdENvcm5lciA9IGV2ZW50LmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBjb25zdCB0b3BDb3JuZXIgPSBldmVudC5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcclxuICAgIGNvbnN0IGJvdHRvbUNvcm5lciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XHJcblxyXG4gICAgaWYgKCFsZWZ0Q29ybmVyICYmICF0b3BDb3JuZXIgJiYgIWJvdHRvbUNvcm5lciAmJiAhcmlnaHRDb3JuZXIpIHtcclxuICAgICAgICAvLyBPcHRpbWlzYXRpb246IG91dCBvZiBhbGwgY29ybmVyIGFyZWFzIGltcGxpZXMgb3V0IG9mIGJvcmRlcnMuXHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9XHJcbiAgICAvLyBEZXRlY3QgY29ybmVycy5cclxuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcclxuICAgIC8vIERldGVjdCBib3JkZXJzLlxyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcclxuICAgIC8vIE91dCBvZiBib3JkZXIgYXJlYS5cclxuICAgIGVsc2Ugc2V0UmVzaXplKCk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5BcHBsaWNhdGlvbik7XHJcblxyXG5jb25zdCBIaWRlTWV0aG9kID0gMDtcclxuY29uc3QgU2hvd01ldGhvZCA9IDE7XHJcbmNvbnN0IFF1aXRNZXRob2QgPSAyO1xyXG5cclxuLyoqXHJcbiAqIEhpZGVzIGEgY2VydGFpbiBtZXRob2QgYnkgY2FsbGluZyB0aGUgSGlkZU1ldGhvZCBmdW5jdGlvbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoSGlkZU1ldGhvZCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyB0aGUgU2hvd01ldGhvZCBhbmQgcmV0dXJucyB0aGUgcmVzdWx0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNob3coKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICByZXR1cm4gY2FsbChTaG93TWV0aG9kKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIHRoZSBRdWl0TWV0aG9kIHRvIHRlcm1pbmF0ZSB0aGUgcHJvZ3JhbS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWl0KCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCB7IENhbmNlbGxhYmxlUHJvbWlzZSwgdHlwZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzIH0gZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gXCIuL25hbm9pZC5qc1wiO1xyXG5cclxuLy8gU2V0dXBcclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcblxyXG50eXBlIFByb21pc2VSZXNvbHZlcnMgPSBPbWl0PENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8YW55PiwgXCJwcm9taXNlXCIgfCBcIm9uY2FuY2VsbGVkXCI+XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYWxsKTtcclxuY29uc3QgY2FuY2VsQ2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2FuY2VsQ2FsbCk7XHJcbmNvbnN0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwPHN0cmluZywgUHJvbWlzZVJlc29sdmVycz4oKTtcclxuXHJcbmNvbnN0IENhbGxCaW5kaW5nID0gMDtcclxuY29uc3QgQ2FuY2VsTWV0aG9kID0gMFxyXG5cclxuLyoqXHJcbiAqIEhvbGRzIGFsbCByZXF1aXJlZCBpbmZvcm1hdGlvbiBmb3IgYSBiaW5kaW5nIGNhbGwuXHJcbiAqIE1heSBwcm92aWRlIGVpdGhlciBhIG1ldGhvZCBJRCBvciBhIG1ldGhvZCBuYW1lLCBidXQgbm90IGJvdGguXHJcbiAqL1xyXG5leHBvcnQgdHlwZSBDYWxsT3B0aW9ucyA9IHtcclxuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXHJcbiAgICBtZXRob2RJRDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXHJcbiAgICBtZXRob2ROYW1lPzogbmV2ZXI7XHJcbiAgICAvKiogQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuICovXHJcbiAgICBhcmdzOiBhbnlbXTtcclxufSB8IHtcclxuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXHJcbiAgICBtZXRob2RJRD86IG5ldmVyO1xyXG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXHJcbiAgICBtZXRob2ROYW1lOiBzdHJpbmc7XHJcbiAgICAvKiogQXJndW1lbnRzIHRvIGJlIHBhc3NlZCBpbnRvIHRoZSBib3VuZCBtZXRob2QuICovXHJcbiAgICBhcmdzOiBhbnlbXTtcclxufTtcclxuXHJcbi8qKlxyXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHRocm93biBpbiBjYXNlIHRoZSBib3VuZCBtZXRob2QgcmV0dXJucyBhbiBlcnJvci5cclxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgUnVudGltZUVycm9yI25hbWV9IHByb3BlcnR5IGlzIFwiUnVudGltZUVycm9yXCIuXHJcbiAqL1xyXG5leHBvcnQgY2xhc3MgUnVudGltZUVycm9yIGV4dGVuZHMgRXJyb3Ige1xyXG4gICAgLyoqXHJcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IFJ1bnRpbWVFcnJvciBpbnN0YW5jZS5cclxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXHJcbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xyXG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xyXG4gICAgICAgIHRoaXMubmFtZSA9IFwiUnVudGltZUVycm9yXCI7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZW5lcmF0ZXMgYSB1bmlxdWUgSUQgdXNpbmcgdGhlIG5hbm9pZCBsaWJyYXJ5LlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHVuaXF1ZSBJRCB0aGF0IGRvZXMgbm90IGV4aXN0IGluIHRoZSBjYWxsUmVzcG9uc2VzIHNldC5cclxuICovXHJcbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKTogc3RyaW5nIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxsIGEgYm91bmQgbWV0aG9kIGFjY29yZGluZyB0byB0aGUgZ2l2ZW4gY2FsbCBvcHRpb25zLlxyXG4gKlxyXG4gKiBJbiBjYXNlIG9mIGZhaWx1cmUsIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYW4gZXhjZXB0aW9uXHJcbiAqIGFtb25nIFJlZmVyZW5jZUVycm9yICh1bmtub3duIG1ldGhvZCksIFR5cGVFcnJvciAod3JvbmcgYXJndW1lbnQgY291bnQgb3IgdHlwZSksXHJcbiAqIHtAbGluayBSdW50aW1lRXJyb3J9IChtZXRob2QgcmV0dXJuZWQgYW4gZXJyb3IpLCBvciBvdGhlciAobmV0d29yayBvciBpbnRlcm5hbCBlcnJvcnMpLlxyXG4gKiBUaGUgZXhjZXB0aW9uIG1pZ2h0IGhhdmUgYSBcImNhdXNlXCIgZmllbGQgd2l0aCB0aGUgdmFsdWUgcmV0dXJuZWRcclxuICogYnkgdGhlIGFwcGxpY2F0aW9uLSBvciBzZXJ2aWNlLWxldmVsIGVycm9yIG1hcnNoYWxpbmcgZnVuY3Rpb25zLlxyXG4gKlxyXG4gKiBAcGFyYW0gb3B0aW9ucyAtIEEgbWV0aG9kIGNhbGwgZGVzY3JpcHRvci5cclxuICogQHJldHVybnMgVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnM6IENhbGxPcHRpb25zKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xyXG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XHJcblxyXG4gICAgY29uc3QgcmVzdWx0ID0gQ2FuY2VsbGFibGVQcm9taXNlLndpdGhSZXNvbHZlcnM8YW55PigpO1xyXG4gICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZTogcmVzdWx0LnJlc29sdmUsIHJlamVjdDogcmVzdWx0LnJlamVjdCB9KTtcclxuXHJcbiAgICBjb25zdCByZXF1ZXN0ID0gY2FsbChDYWxsQmluZGluZywgT2JqZWN0LmFzc2lnbih7IFwiY2FsbC1pZFwiOiBpZCB9LCBvcHRpb25zKSk7XHJcbiAgICBsZXQgcnVubmluZyA9IHRydWU7XHJcblxyXG4gICAgcmVxdWVzdC50aGVuKChyZXMpID0+IHtcclxuICAgICAgICBydW5uaW5nID0gZmFsc2U7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIHJlc3VsdC5yZXNvbHZlKHJlcyk7XHJcbiAgICB9LCAoZXJyKSA9PiB7XHJcbiAgICAgICAgcnVubmluZyA9IGZhbHNlO1xyXG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICByZXN1bHQucmVqZWN0KGVycik7XHJcbiAgICB9KTtcclxuXHJcbiAgICBjb25zdCBjYW5jZWwgPSAoKSA9PiB7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIHJldHVybiBjYW5jZWxDYWxsKENhbmNlbE1ldGhvZCwge1wiY2FsbC1pZFwiOiBpZH0pLmNhdGNoKChlcnIpID0+IHtcclxuICAgICAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIHdoaWxlIHJlcXVlc3RpbmcgYmluZGluZyBjYWxsIGNhbmNlbGxhdGlvbjpcIiwgZXJyKTtcclxuICAgICAgICB9KTtcclxuICAgIH07XHJcblxyXG4gICAgcmVzdWx0Lm9uY2FuY2VsbGVkID0gKCkgPT4ge1xyXG4gICAgICAgIGlmIChydW5uaW5nKSB7XHJcbiAgICAgICAgICAgIHJldHVybiBjYW5jZWwoKTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICByZXR1cm4gcmVxdWVzdC50aGVuKGNhbmNlbCk7XHJcbiAgICAgICAgfVxyXG4gICAgfTtcclxuXHJcbiAgICByZXR1cm4gcmVzdWx0LnByb21pc2U7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyBhIGJvdW5kIG1ldGhvZCBieSBuYW1lIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXHJcbiAqXHJcbiAqIEBwYXJhbSBtZXRob2ROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnLlxyXG4gKiBAcGFyYW0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBCeU5hbWUobWV0aG9kTmFtZTogc3RyaW5nLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcclxuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kTmFtZSwgYXJncyB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIGJ5IGl0cyBudW1lcmljIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqIFNlZSB7QGxpbmsgQ2FsbH0gZm9yIGRldGFpbHMuXHJcbiAqXHJcbiAqIEBwYXJhbSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm4gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnlJRChtZXRob2RJRDogbnVtYmVyLCAuLi5hcmdzOiBhbnlbXSk6IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+IHtcclxuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kSUQsIGFyZ3MgfSk7XHJcbn1cclxuIiwgIi8vIFNvdXJjZTogaHR0cHM6Ly9naXRodWIuY29tL2luc3BlY3QtanMvaXMtY2FsbGFibGVcclxuXHJcbi8vIFRoZSBNSVQgTGljZW5zZSAoTUlUKVxyXG4vL1xyXG4vLyBDb3B5cmlnaHQgKGMpIDIwMTUgSm9yZGFuIEhhcmJhbmRcclxuLy9cclxuLy8gUGVybWlzc2lvbiBpcyBoZXJlYnkgZ3JhbnRlZCwgZnJlZSBvZiBjaGFyZ2UsIHRvIGFueSBwZXJzb24gb2J0YWluaW5nIGEgY29weVxyXG4vLyBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsXHJcbi8vIGluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHNcclxuLy8gdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbFxyXG4vLyBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXNcclxuLy8gZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczpcclxuLy9cclxuLy8gVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXHJcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXHJcbi8vXHJcbi8vIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1JcclxuLy8gSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksXHJcbi8vIEZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRVxyXG4vLyBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSXHJcbi8vIExJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sXHJcbi8vIE9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFXHJcbi8vIFNPRlRXQVJFLlxyXG5cclxudmFyIGZuVG9TdHIgPSBGdW5jdGlvbi5wcm90b3R5cGUudG9TdHJpbmc7XHJcbnZhciByZWZsZWN0QXBwbHk6IHR5cGVvZiBSZWZsZWN0LmFwcGx5IHwgZmFsc2UgfCBudWxsID0gdHlwZW9mIFJlZmxlY3QgPT09ICdvYmplY3QnICYmIFJlZmxlY3QgIT09IG51bGwgJiYgUmVmbGVjdC5hcHBseTtcclxudmFyIGJhZEFycmF5TGlrZTogYW55O1xyXG52YXIgaXNDYWxsYWJsZU1hcmtlcjogYW55O1xyXG5pZiAodHlwZW9mIHJlZmxlY3RBcHBseSA9PT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgT2JqZWN0LmRlZmluZVByb3BlcnR5ID09PSAnZnVuY3Rpb24nKSB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGJhZEFycmF5TGlrZSA9IE9iamVjdC5kZWZpbmVQcm9wZXJ0eSh7fSwgJ2xlbmd0aCcsIHtcclxuICAgICAgICAgICAgZ2V0OiBmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgICAgICB0aHJvdyBpc0NhbGxhYmxlTWFya2VyO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgaXNDYWxsYWJsZU1hcmtlciA9IHt9O1xyXG4gICAgICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZSBuby10aHJvdy1saXRlcmFsXHJcbiAgICAgICAgcmVmbGVjdEFwcGx5KGZ1bmN0aW9uICgpIHsgdGhyb3cgNDI7IH0sIG51bGwsIGJhZEFycmF5TGlrZSk7XHJcbiAgICB9IGNhdGNoIChfKSB7XHJcbiAgICAgICAgaWYgKF8gIT09IGlzQ2FsbGFibGVNYXJrZXIpIHtcclxuICAgICAgICAgICAgcmVmbGVjdEFwcGx5ID0gbnVsbDtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbn0gZWxzZSB7XHJcbiAgICByZWZsZWN0QXBwbHkgPSBudWxsO1xyXG59XHJcblxyXG52YXIgY29uc3RydWN0b3JSZWdleCA9IC9eXFxzKmNsYXNzXFxiLztcclxudmFyIGlzRVM2Q2xhc3NGbiA9IGZ1bmN0aW9uIGlzRVM2Q2xhc3NGdW5jdGlvbih2YWx1ZTogYW55KTogYm9vbGVhbiB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIHZhciBmblN0ciA9IGZuVG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICAgICAgcmV0dXJuIGNvbnN0cnVjdG9yUmVnZXgudGVzdChmblN0cik7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlOyAvLyBub3QgYSBmdW5jdGlvblxyXG4gICAgfVxyXG59O1xyXG5cclxudmFyIHRyeUZ1bmN0aW9uT2JqZWN0ID0gZnVuY3Rpb24gdHJ5RnVuY3Rpb25Ub1N0cih2YWx1ZTogYW55KTogYm9vbGVhbiB7XHJcbiAgICB0cnkge1xyXG4gICAgICAgIGlmIChpc0VTNkNsYXNzRm4odmFsdWUpKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgICAgIGZuVG9TdHIuY2FsbCh2YWx1ZSk7XHJcbiAgICAgICAgcmV0dXJuIHRydWU7XHJcbiAgICB9IGNhdGNoIChlKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG59O1xyXG52YXIgdG9TdHIgPSBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nO1xyXG52YXIgb2JqZWN0Q2xhc3MgPSAnW29iamVjdCBPYmplY3RdJztcclxudmFyIGZuQ2xhc3MgPSAnW29iamVjdCBGdW5jdGlvbl0nO1xyXG52YXIgZ2VuQ2xhc3MgPSAnW29iamVjdCBHZW5lcmF0b3JGdW5jdGlvbl0nO1xyXG52YXIgZGRhQ2xhc3MgPSAnW29iamVjdCBIVE1MQWxsQ29sbGVjdGlvbl0nOyAvLyBJRSAxMVxyXG52YXIgZGRhQ2xhc3MyID0gJ1tvYmplY3QgSFRNTCBkb2N1bWVudC5hbGwgY2xhc3NdJztcclxudmFyIGRkYUNsYXNzMyA9ICdbb2JqZWN0IEhUTUxDb2xsZWN0aW9uXSc7IC8vIElFIDktMTBcclxudmFyIGhhc1RvU3RyaW5nVGFnID0gdHlwZW9mIFN5bWJvbCA9PT0gJ2Z1bmN0aW9uJyAmJiAhIVN5bWJvbC50b1N0cmluZ1RhZzsgLy8gYmV0dGVyOiB1c2UgYGhhcy10b3N0cmluZ3RhZ2BcclxuXHJcbnZhciBpc0lFNjggPSAhKDAgaW4gWyxdKTsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBuby1zcGFyc2UtYXJyYXlzLCBjb21tYS1zcGFjaW5nXHJcblxyXG52YXIgaXNEREE6ICh2YWx1ZTogYW55KSA9PiBib29sZWFuID0gZnVuY3Rpb24gaXNEb2N1bWVudERvdEFsbCgpIHsgcmV0dXJuIGZhbHNlOyB9O1xyXG5pZiAodHlwZW9mIGRvY3VtZW50ID09PSAnb2JqZWN0Jykge1xyXG4gICAgLy8gRmlyZWZveCAzIGNhbm9uaWNhbGl6ZXMgRERBIHRvIHVuZGVmaW5lZCB3aGVuIGl0J3Mgbm90IGFjY2Vzc2VkIGRpcmVjdGx5XHJcbiAgICB2YXIgYWxsID0gZG9jdW1lbnQuYWxsO1xyXG4gICAgaWYgKHRvU3RyLmNhbGwoYWxsKSA9PT0gdG9TdHIuY2FsbChkb2N1bWVudC5hbGwpKSB7XHJcbiAgICAgICAgaXNEREEgPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKHZhbHVlKSB7XHJcbiAgICAgICAgICAgIC8qIGdsb2JhbHMgZG9jdW1lbnQ6IGZhbHNlICovXHJcbiAgICAgICAgICAgIC8vIGluIElFIDYtOCwgdHlwZW9mIGRvY3VtZW50LmFsbCBpcyBcIm9iamVjdFwiIGFuZCBpdCdzIHRydXRoeVxyXG4gICAgICAgICAgICBpZiAoKGlzSUU2OCB8fCAhdmFsdWUpICYmICh0eXBlb2YgdmFsdWUgPT09ICd1bmRlZmluZWQnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcpKSB7XHJcbiAgICAgICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgICAgIHZhciBzdHIgPSB0b1N0ci5jYWxsKHZhbHVlKTtcclxuICAgICAgICAgICAgICAgICAgICByZXR1cm4gKFxyXG4gICAgICAgICAgICAgICAgICAgICAgICBzdHIgPT09IGRkYUNsYXNzXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MyXHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MzIC8vIG9wZXJhIDEyLjE2XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gb2JqZWN0Q2xhc3MgLy8gSUUgNi04XHJcbiAgICAgICAgICAgICAgICAgICAgKSAmJiB2YWx1ZSgnJykgPT0gbnVsbDsgLy8gZXNsaW50LWRpc2FibGUtbGluZSBlcWVxZXFcclxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGUpIHsgLyoqLyB9XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGlzQ2FsbGFibGVSZWZBcHBseTxUPih2YWx1ZTogVCB8IHVua25vd24pOiB2YWx1ZSBpcyAoLi4uYXJnczogYW55W10pID0+IGFueSAge1xyXG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxyXG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIHRyeSB7XHJcbiAgICAgICAgKHJlZmxlY3RBcHBseSBhcyBhbnkpKHZhbHVlLCBudWxsLCBiYWRBcnJheUxpa2UpO1xyXG4gICAgfSBjYXRjaCAoZSkge1xyXG4gICAgICAgIGlmIChlICE9PSBpc0NhbGxhYmxlTWFya2VyKSB7IHJldHVybiBmYWxzZTsgfVxyXG4gICAgfVxyXG4gICAgcmV0dXJuICFpc0VTNkNsYXNzRm4odmFsdWUpICYmIHRyeUZ1bmN0aW9uT2JqZWN0KHZhbHVlKTtcclxufVxyXG5cclxuZnVuY3Rpb24gaXNDYWxsYWJsZU5vUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkge1xyXG4gICAgaWYgKGlzRERBKHZhbHVlKSkgeyByZXR1cm4gdHJ1ZTsgfVxyXG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIGlmICh0eXBlb2YgdmFsdWUgIT09ICdmdW5jdGlvbicgJiYgdHlwZW9mIHZhbHVlICE9PSAnb2JqZWN0JykgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIGlmIChoYXNUb1N0cmluZ1RhZykgeyByZXR1cm4gdHJ5RnVuY3Rpb25PYmplY3QodmFsdWUpOyB9XHJcbiAgICBpZiAoaXNFUzZDbGFzc0ZuKHZhbHVlKSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIHZhciBzdHJDbGFzcyA9IHRvU3RyLmNhbGwodmFsdWUpO1xyXG4gICAgaWYgKHN0ckNsYXNzICE9PSBmbkNsYXNzICYmIHN0ckNsYXNzICE9PSBnZW5DbGFzcyAmJiAhKC9eXFxbb2JqZWN0IEhUTUwvKS50ZXN0KHN0ckNsYXNzKSkgeyByZXR1cm4gZmFsc2U7IH1cclxuICAgIHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7XHJcbn07XHJcblxyXG5leHBvcnQgZGVmYXVsdCByZWZsZWN0QXBwbHkgPyBpc0NhbGxhYmxlUmVmQXBwbHkgOiBpc0NhbGxhYmxlTm9SZWZBcHBseTtcclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbmltcG9ydCBpc0NhbGxhYmxlIGZyb20gXCIuL2NhbGxhYmxlLmpzXCI7XHJcblxyXG4vKipcclxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB1c2VkIGFzIHJlamVjdGlvbiByZWFzb25cclxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IGlzIGNhbmNlbGxlZCBzdWNjZXNzZnVsbHkuXHJcbiAqXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsRXJyb3JcImAuXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgY2F1c2UgcGFzc2VkIHRvIHRoZSBjYW5jZWwgbWV0aG9kLCBpZiBhbnkuXHJcbiAqL1xyXG5leHBvcnQgY2xhc3MgQ2FuY2VsRXJyb3IgZXh0ZW5kcyBFcnJvciB7XHJcbiAgICAvKipcclxuICAgICAqIENvbnN0cnVjdHMgYSBuZXcgYENhbmNlbEVycm9yYCBpbnN0YW5jZS5cclxuICAgICAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIGVycm9yIG1lc3NhZ2UuXHJcbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IobWVzc2FnZT86IHN0cmluZywgb3B0aW9ucz86IEVycm9yT3B0aW9ucykge1xyXG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xyXG4gICAgICAgIHRoaXMubmFtZSA9IFwiQ2FuY2VsRXJyb3JcIjtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvblxyXG4gKiBpbiBjYXNlIGEge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZX0gcmVqZWN0cyBhZnRlciBiZWluZyBjYW5jZWxsZWQsXHJcbiAqIG9yIHdoZW4gdGhlIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgdGhyb3dzIG9yIHJlamVjdHMuXHJcbiAqXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIG5hbWV9IHByb3BlcnR5IGlzIHRoZSBzdHJpbmcgYFwiQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3JcImAuXHJcbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgcmVhc29uIHRoZSBwcm9taXNlIHJlamVjdGVkIHdpdGguXHJcbiAqXHJcbiAqIEJlY2F1c2UgdGhlIG9yaWdpbmFsIHByb21pc2Ugd2FzIGNhbmNlbGxlZCxcclxuICogYSB3cmFwcGVyIHByb21pc2Ugd2lsbCBiZSBwYXNzZWQgdG8gdGhlIHVuaGFuZGxlZCByZWplY3Rpb24gbGlzdGVuZXIgaW5zdGVhZC5cclxuICogVGhlIHtAbGluayBwcm9taXNlfSBwcm9wZXJ0eSBob2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgb3JpZ2luYWwgcHJvbWlzZS5cclxuICovXHJcbmV4cG9ydCBjbGFzcyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvciBleHRlbmRzIEVycm9yIHtcclxuICAgIC8qKlxyXG4gICAgICogSG9sZHMgYSByZWZlcmVuY2UgdG8gdGhlIHByb21pc2UgdGhhdCB3YXMgY2FuY2VsbGVkIGFuZCB0aGVuIHJlamVjdGVkLlxyXG4gICAgICovXHJcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj47XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcmAgaW5zdGFuY2UuXHJcbiAgICAgKiBAcGFyYW0gcHJvbWlzZSAtIFRoZSBwcm9taXNlIHRoYXQgY2F1c2VkIHRoZSBlcnJvciBvcmlnaW5hbGx5LlxyXG4gICAgICogQHBhcmFtIHJlYXNvbiAtIFRoZSByZWplY3Rpb24gcmVhc29uLlxyXG4gICAgICogQHBhcmFtIGluZm8gLSBBbiBvcHRpb25hbCBpbmZvcm1hdGl2ZSBtZXNzYWdlIHNwZWNpZnlpbmcgdGhlIGNpcmN1bXN0YW5jZXMgaW4gd2hpY2ggdGhlIGVycm9yIHdhcyB0aHJvd24uXHJcbiAgICAgKiAgICAgICAgICAgICAgIERlZmF1bHRzIHRvIHRoZSBzdHJpbmcgYFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cImAuXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgcmVhc29uPzogYW55LCBpbmZvPzogc3RyaW5nKSB7XHJcbiAgICAgICAgc3VwZXIoKGluZm8gPz8gXCJVbmhhbmRsZWQgcmVqZWN0aW9uIGluIGNhbmNlbGxlZCBwcm9taXNlLlwiKSArIFwiIFJlYXNvbjogXCIgKyBlcnJvck1lc3NhZ2UocmVhc29uKSwgeyBjYXVzZTogcmVhc29uIH0pO1xyXG4gICAgICAgIHRoaXMucHJvbWlzZSA9IHByb21pc2U7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiO1xyXG4gICAgfVxyXG59XHJcblxyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+ID0gKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xyXG50eXBlIENhbmNlbGxhYmxlUHJvbWlzZVJlamVjdG9yID0gKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcclxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgPSAoY2F1c2U/OiBhbnkpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcclxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VFeGVjdXRvcjxUPiA9IChyZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiwgcmVqZWN0OiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvcikgPT4gdm9pZDtcclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XHJcbiAgICB0aGVuPFRSZXN1bHQxID0gVCwgVFJlc3VsdDIgPSBuZXZlcj4ob25mdWxmaWxsZWQ/OiAoKHZhbHVlOiBUKSA9PiBUUmVzdWx0MSB8IFByb21pc2VMaWtlPFRSZXN1bHQxPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDE+KSB8IHVuZGVmaW5lZCB8IG51bGwsIG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiBUUmVzdWx0MiB8IFByb21pc2VMaWtlPFRSZXN1bHQyPiB8IENhbmNlbGxhYmxlUHJvbWlzZUxpa2U8VFJlc3VsdDI+KSB8IHVuZGVmaW5lZCB8IG51bGwpOiBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxIHwgVFJlc3VsdDI+O1xyXG4gICAgY2FuY2VsKGNhdXNlPzogYW55KTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG59XHJcblxyXG4vKipcclxuICogV3JhcHMgYSBjYW5jZWxsYWJsZSBwcm9taXNlIGFsb25nIHdpdGggaXRzIHJlc29sdXRpb24gbWV0aG9kcy5cclxuICogVGhlIGBvbmNhbmNlbGxlZGAgZmllbGQgd2lsbCBiZSBudWxsIGluaXRpYWxseSBidXQgbWF5IGJlIHNldCB0byBwcm92aWRlIGEgY3VzdG9tIGNhbmNlbGxhdGlvbiBmdW5jdGlvbi5cclxuICovXHJcbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XHJcbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD47XHJcbiAgICByZXNvbHZlOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPjtcclxuICAgIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3I7XHJcbiAgICBvbmNhbmNlbGxlZDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyIHwgbnVsbDtcclxufVxyXG5cclxuaW50ZXJmYWNlIENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlIHtcclxuICAgIHJlYWRvbmx5IHJvb3Q6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlO1xyXG4gICAgcmVzb2x2aW5nOiBib29sZWFuO1xyXG4gICAgc2V0dGxlZDogYm9vbGVhbjtcclxuICAgIHJlYXNvbj86IENhbmNlbEVycm9yO1xyXG59XHJcblxyXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxyXG5jb25zdCBiYXJyaWVyU3ltID0gU3ltYm9sKFwiYmFycmllclwiKTtcclxuY29uc3QgY2FuY2VsSW1wbFN5bSA9IFN5bWJvbChcImNhbmNlbEltcGxcIik7XHJcbmNvbnN0IHNwZWNpZXM6IHR5cGVvZiBTeW1ib2wuc3BlY2llcyA9IFN5bWJvbC5zcGVjaWVzID8/IFN5bWJvbChcInNwZWNpZXNQb2x5ZmlsbFwiKTtcclxuXHJcbi8qKlxyXG4gKiBBIHByb21pc2Ugd2l0aCBhbiBhdHRhY2hlZCBtZXRob2QgZm9yIGNhbmNlbGxpbmcgbG9uZy1ydW5uaW5nIG9wZXJhdGlvbnMgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbH0pLlxyXG4gKiBDYW5jZWxsYXRpb24gY2FuIG9wdGlvbmFsbHkgYmUgYm91bmQgdG8gYW4ge0BsaW5rIEFib3J0U2lnbmFsfVxyXG4gKiBmb3IgYmV0dGVyIGNvbXBvc2FiaWxpdHkgKHNlZSB7QGxpbmsgQ2FuY2VsbGFibGVQcm9taXNlI2NhbmNlbE9ufSkuXHJcbiAqXHJcbiAqIENhbmNlbGxpbmcgYSBwZW5kaW5nIHByb21pc2Ugd2lsbCByZXN1bHQgaW4gYW4gaW1tZWRpYXRlIHJlamVjdGlvblxyXG4gKiB3aXRoIGFuIGluc3RhbmNlIG9mIHtAbGluayBDYW5jZWxFcnJvcn0gYXMgcmVhc29uLFxyXG4gKiBidXQgd2hvZXZlciBzdGFydGVkIHRoZSBwcm9taXNlIHdpbGwgYmUgcmVzcG9uc2libGVcclxuICogZm9yIGFjdHVhbGx5IGFib3J0aW5nIHRoZSB1bmRlcmx5aW5nIG9wZXJhdGlvbi5cclxuICogVG8gdGhpcyBwdXJwb3NlLCB0aGUgY29uc3RydWN0b3IgYW5kIGFsbCBjaGFpbmluZyBtZXRob2RzXHJcbiAqIGFjY2VwdCBvcHRpb25hbCBjYW5jZWxsYXRpb24gY2FsbGJhY2tzLlxyXG4gKlxyXG4gKiBJZiBhIGBDYW5jZWxsYWJsZVByb21pc2VgIHN0aWxsIHJlc29sdmVzIGFmdGVyIGhhdmluZyBiZWVuIGNhbmNlbGxlZCxcclxuICogdGhlIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZC4gSWYgaXQgcmVqZWN0cywgdGhlIHJlYXNvblxyXG4gKiB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24sXHJcbiAqIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGluc3RhbmNlLlxyXG4gKiBUbyBmYWNpbGl0YXRlIHRoZSBoYW5kbGluZyBvZiBjYW5jZWxsYXRpb24gcmVxdWVzdHMsXHJcbiAqIGNhbmNlbGxlZCBgQ2FuY2VsbGFibGVQcm9taXNlYHMgd2lsbCBfbm90XyByZXBvcnQgdW5oYW5kbGVkIGBDYW5jZWxFcnJvcmBzXHJcbiAqIHdob3NlIGBjYXVzZWAgZmllbGQgaXMgdGhlIHNhbWUgYXMgdGhlIG9uZSB3aXRoIHdoaWNoIHRoZSBjdXJyZW50IHByb21pc2Ugd2FzIGNhbmNlbGxlZC5cclxuICpcclxuICogQWxsIHVzdWFsIHByb21pc2UgbWV0aG9kcyBhcmUgZGVmaW5lZCBhbmQgcmV0dXJuIGEgYENhbmNlbGxhYmxlUHJvbWlzZWBcclxuICogd2hvc2UgY2FuY2VsIG1ldGhvZCB3aWxsIGNhbmNlbCB0aGUgcGFyZW50IG9wZXJhdGlvbiBhcyB3ZWxsLCBwcm9wYWdhdGluZyB0aGUgY2FuY2VsbGF0aW9uIHJlYXNvblxyXG4gKiB1cHdhcmRzIHRocm91Z2ggcHJvbWlzZSBjaGFpbnMuXHJcbiAqIENvbnZlcnNlbHksIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpbGwgbm90IGF1dG9tYXRpY2FsbHkgY2FuY2VsIGRlcGVuZGVudCBwcm9taXNlcyBkb3duc3RyZWFtOlxyXG4gKiBgYGB0c1xyXG4gKiBsZXQgcm9vdCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XHJcbiAqIGxldCBjaGlsZDEgPSByb290LnRoZW4oKCkgPT4geyAuLi4gfSk7XHJcbiAqIGxldCBjaGlsZDIgPSBjaGlsZDEudGhlbigoKSA9PiB7IC4uLiB9KTtcclxuICogbGV0IGNoaWxkMyA9IHJvb3QuY2F0Y2goKCkgPT4geyAuLi4gfSk7XHJcbiAqIGNoaWxkMS5jYW5jZWwoKTsgLy8gQ2FuY2VscyBjaGlsZDEgYW5kIHJvb3QsIGJ1dCBub3QgY2hpbGQyIG9yIGNoaWxkM1xyXG4gKiBgYGBcclxuICogQ2FuY2VsbGluZyBhIHByb21pc2UgdGhhdCBoYXMgYWxyZWFkeSBzZXR0bGVkIGlzIHNhZmUgYW5kIGhhcyBubyBjb25zZXF1ZW5jZS5cclxuICpcclxuICogVGhlIGBjYW5jZWxgIG1ldGhvZCByZXR1cm5zIGEgcHJvbWlzZSB0aGF0IF9hbHdheXMgZnVsZmlsbHNfXHJcbiAqIGFmdGVyIHRoZSB3aG9sZSBjaGFpbiBoYXMgcHJvY2Vzc2VkIHRoZSBjYW5jZWwgcmVxdWVzdFxyXG4gKiBhbmQgYWxsIGF0dGFjaGVkIGNhbGxiYWNrcyB1cCB0byB0aGF0IG1vbWVudCBoYXZlIHJ1bi5cclxuICpcclxuICogQWxsIEVTMjAyNCBwcm9taXNlIG1ldGhvZHMgKHN0YXRpYyBhbmQgaW5zdGFuY2UpIGFyZSBkZWZpbmVkIG9uIENhbmNlbGxhYmxlUHJvbWlzZSxcclxuICogYnV0IGFjdHVhbCBhdmFpbGFiaWxpdHkgbWF5IHZhcnkgd2l0aCBPUy93ZWJ2aWV3IHZlcnNpb24uXHJcbiAqXHJcbiAqIEluIGxpbmUgd2l0aCB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZyxcclxuICogYENhbmNlbGxhYmxlUHJvbWlzZWAgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy5cclxuICogRXh0ZW5kZXJzIHNob3VsZCB0YWtlIGNhcmUgdG8gcHJvdmlkZSB0aGVpciBvd24gbWV0aG9kIGltcGxlbWVudGF0aW9ucy5cclxuICogVGhpcyBtaWdodCBiZSByZWNvbnNpZGVyZWQgaW4gY2FzZSB0aGUgcHJvcG9zYWwgaXMgcmV0aXJlZC5cclxuICpcclxuICogQ2FuY2VsbGFibGVQcm9taXNlIGlzIGEgd3JhcHBlciBhcm91bmQgdGhlIERPTSBQcm9taXNlIG9iamVjdFxyXG4gKiBhbmQgaXMgY29tcGxpYW50IHdpdGggdGhlIFtQcm9taXNlcy9BKyBzcGVjaWZpY2F0aW9uXShodHRwczovL3Byb21pc2VzYXBsdXMuY29tLylcclxuICogKGl0IHBhc3NlcyB0aGUgW2NvbXBsaWFuY2Ugc3VpdGVdKGh0dHBzOi8vZ2l0aHViLmNvbS9wcm9taXNlcy1hcGx1cy9wcm9taXNlcy10ZXN0cykpXHJcbiAqIGlmIHNvIGlzIHRoZSB1bmRlcmx5aW5nIGltcGxlbWVudGF0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIENhbmNlbGxhYmxlUHJvbWlzZTxUPiBleHRlbmRzIFByb21pc2U8VD4gaW1wbGVtZW50cyBQcm9taXNlTGlrZTxUPiwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XHJcbiAgICAvLyBQcml2YXRlIGZpZWxkcy5cclxuICAgIC8qKiBAaW50ZXJuYWwgKi9cclxuICAgIHByaXZhdGUgW2JhcnJpZXJTeW1dITogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gfCBudWxsO1xyXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xyXG4gICAgcHJpdmF0ZSByZWFkb25seSBbY2FuY2VsSW1wbFN5bV0hOiAocmVhc29uOiBDYW5jZWxFcnJvcikgPT4gdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBgQ2FuY2VsbGFibGVQcm9taXNlYC5cclxuICAgICAqXHJcbiAgICAgKiBAcGFyYW0gZXhlY3V0b3IgLSBBIGNhbGxiYWNrIHVzZWQgdG8gaW5pdGlhbGl6ZSB0aGUgcHJvbWlzZS4gVGhpcyBjYWxsYmFjayBpcyBwYXNzZWQgdHdvIGFyZ3VtZW50czpcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGEgYHJlc29sdmVgIGNhbGxiYWNrIHVzZWQgdG8gcmVzb2x2ZSB0aGUgcHJvbWlzZSB3aXRoIGEgdmFsdWVcclxuICAgICAqICAgICAgICAgICAgICAgICAgIG9yIHRoZSByZXN1bHQgb2YgYW5vdGhlciBwcm9taXNlIChwb3NzaWJseSBjYW5jZWxsYWJsZSksXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYSBgcmVqZWN0YCBjYWxsYmFjayB1c2VkIHRvIHJlamVjdCB0aGUgcHJvbWlzZSB3aXRoIGEgcHJvdmlkZWQgcmVhc29uIG9yIGVycm9yLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgdGhlIHZhbHVlIHByb3ZpZGVkIHRvIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcclxuICAgICAqICAgICAgICAgICAgICAgICAgIChpdCBoYXMgYSBgdGhlbmAgX2FuZF8gYSBgY2FuY2VsYCBtZXRob2QpLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZm9yd2FyZGVkIHRvIHRoYXQgb2JqZWN0IGFuZCB0aGUgb25jYW5jZWxsZWQgd2lsbCBub3QgYmUgaW52b2tlZCBhbnltb3JlLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgSWYgYW55IG9uZSBvZiB0aGUgdHdvIGNhbGxiYWNrcyBpcyBjYWxsZWQgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQsXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGUgcHJvdmlkZWQgdmFsdWVzIHdpbGwgYmUgY2FuY2VsbGVkIGFuZCByZXNvbHZlZCBhcyB1c3VhbCxcclxuICAgICAqICAgICAgICAgICAgICAgICAgIGJ1dCB0aGVpciByZXN1bHRzIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgSG93ZXZlciwgaWYgdGhlIHJlc29sdXRpb24gcHJvY2VzcyB1bHRpbWF0ZWx5IGVuZHMgdXAgaW4gYSByZWplY3Rpb25cclxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoYXQgaXMgbm90IGR1ZSB0byBjYW5jZWxsYXRpb24sIHRoZSByZWplY3Rpb24gcmVhc29uXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICB3aWxsIGJlIHdyYXBwZWQgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9XHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICogQHBhcmFtIG9uY2FuY2VsbGVkIC0gSXQgaXMgdGhlIGNhbGxlcidzIHJlc3BvbnNpYmlsaXR5IHRvIGVuc3VyZSB0aGF0IGFueSBvcGVyYXRpb25cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHN0YXJ0ZWQgYnkgdGhlIGV4ZWN1dG9yIGlzIHByb3Blcmx5IGhhbHRlZCB1cG9uIGNhbmNlbGxhdGlvbi5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFRoaXMgb3B0aW9uYWwgY2FsbGJhY2sgY2FuIGJlIHVzZWQgdG8gdGhhdCBwdXJwb3NlLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSXQgd2lsbCBiZSBjYWxsZWQgX3N5bmNocm9ub3VzbHlfIHdpdGggYSBjYW5jZWxsYXRpb24gY2F1c2VcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIHdoZW4gY2FuY2VsbGF0aW9uIGlzIHJlcXVlc3RlZCwgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSByZWplY3RlZFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2l0aCBhIHtAbGluayBDYW5jZWxFcnJvcn0sIGJ1dCBfYmVmb3JlX1xyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW55IHtAbGluayB0aGVufS97QGxpbmsgY2F0Y2h9L3tAbGluayBmaW5hbGx5fSBjYWxsYmFjayBydW5zLlxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSWYgdGhlIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfVxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgd2lsbCBvbmx5IGZ1bGZpbGwgYWZ0ZXIgdGhlIGZvcm1lciBoYXMgc2V0dGxlZC5cclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFVuaGFuZGxlZCBleGNlcHRpb25zIG9yIHJlamVjdGlvbnMgZnJvbSB0aGUgY2FsbGJhY2sgd2lsbCBiZSB3cmFwcGVkXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBpbiBhIHtAbGluayBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcn0gYW5kIGJ1YmJsZWQgdXAgYXMgdW5oYW5kbGVkIHJlamVjdGlvbnMuXHJcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBJZiB0aGUgYHJlc29sdmVgIGNhbGxiYWNrIGlzIGNhbGxlZCBiZWZvcmUgY2FuY2VsbGF0aW9uIHdpdGggYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoaXMgcHJvbWlzZSB3aWxsIGJlIGRpdmVydGVkIHRvIHRoYXQgcHJvbWlzZSxcclxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIGFuZCB0aGUgb3JpZ2luYWwgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB3aWxsIGJlIGRpc2NhcmRlZC5cclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoZXhlY3V0b3I6IENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+LCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcikge1xyXG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XHJcbiAgICAgICAgbGV0IHJlamVjdCE6IChyZWFzb24/OiBhbnkpID0+IHZvaWQ7XHJcbiAgICAgICAgc3VwZXIoKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XHJcblxyXG4gICAgICAgIGlmICgodGhpcy5jb25zdHJ1Y3RvciBhcyBhbnkpW3NwZWNpZXNdICE9PSBQcm9taXNlKSB7XHJcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy4gUGxlYXNlIHJlZnJhaW4gZnJvbSBvdmVycmlkaW5nIHRoZSBbU3ltYm9sLnNwZWNpZXNdIHN0YXRpYyBwcm9wZXJ0eS5cIik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBsZXQgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHtcclxuICAgICAgICAgICAgcHJvbWlzZTogdGhpcyxcclxuICAgICAgICAgICAgcmVzb2x2ZSxcclxuICAgICAgICAgICAgcmVqZWN0LFxyXG4gICAgICAgICAgICBnZXQgb25jYW5jZWxsZWQoKSB7IHJldHVybiBvbmNhbmNlbGxlZCA/PyBudWxsOyB9LFxyXG4gICAgICAgICAgICBzZXQgb25jYW5jZWxsZWQoY2IpIHsgb25jYW5jZWxsZWQgPSBjYiA/PyB1bmRlZmluZWQ7IH1cclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICBjb25zdCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XHJcbiAgICAgICAgICAgIGdldCByb290KCkgeyByZXR1cm4gc3RhdGU7IH0sXHJcbiAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXHJcbiAgICAgICAgICAgIHNldHRsZWQ6IGZhbHNlXHJcbiAgICAgICAgfTtcclxuXHJcbiAgICAgICAgLy8gU2V0dXAgY2FuY2VsbGF0aW9uIHN5c3RlbS5cclxuICAgICAgICB2b2lkIE9iamVjdC5kZWZpbmVQcm9wZXJ0aWVzKHRoaXMsIHtcclxuICAgICAgICAgICAgW2JhcnJpZXJTeW1dOiB7XHJcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXHJcbiAgICAgICAgICAgICAgICB3cml0YWJsZTogdHJ1ZSxcclxuICAgICAgICAgICAgICAgIHZhbHVlOiBudWxsXHJcbiAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgIFtjYW5jZWxJbXBsU3ltXToge1xyXG4gICAgICAgICAgICAgICAgY29uZmlndXJhYmxlOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IGZhbHNlLFxyXG4gICAgICAgICAgICAgICAgdmFsdWU6IGNhbmNlbGxlckZvcihwcm9taXNlLCBzdGF0ZSlcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH0pO1xyXG5cclxuICAgICAgICAvLyBSdW4gdGhlIGFjdHVhbCBleGVjdXRvci5cclxuICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIHN0YXRlKTtcclxuICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICBleGVjdXRvcihyZXNvbHZlckZvcihwcm9taXNlLCBzdGF0ZSksIHJlamVjdG9yKTtcclxuICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykge1xyXG4gICAgICAgICAgICAgICAgY29uc29sZS5sb2coXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIENhbmNlbGxhYmxlUHJvbWlzZSBleGVjdXRvci5cIiwgZXJyKTtcclxuICAgICAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgICAgIHJlamVjdG9yKGVycik7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDYW5jZWxzIGltbWVkaWF0ZWx5IHRoZSBleGVjdXRpb24gb2YgdGhlIG9wZXJhdGlvbiBhc3NvY2lhdGVkIHdpdGggdGhpcyBwcm9taXNlLlxyXG4gICAgICogVGhlIHByb21pc2UgcmVqZWN0cyB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSBpbnN0YW5jZSBhcyByZWFzb24sXHJcbiAgICAgKiB3aXRoIHRoZSB7QGxpbmsgQ2FuY2VsRXJyb3IjY2F1c2V9IHByb3BlcnR5IHNldCB0byB0aGUgZ2l2ZW4gYXJndW1lbnQsIGlmIGFueS5cclxuICAgICAqXHJcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCBhZnRlciB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkO1xyXG4gICAgICogcmVwZWF0ZWQgY2FsbHMgaW4gcGFydGljdWxhciBhcmUgc2FmZSwgYnV0IG9ubHkgdGhlIGZpcnN0IG9uZVxyXG4gICAgICogd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgYENhbmNlbEVycm9yYCBleGNlcHRpb24gX25lZWQgbm90XyBiZSBoYW5kbGVkIGV4cGxpY2l0bHkgX29uIHRoZSBwcm9taXNlcyB0aGF0IGFyZSBiZWluZyBjYW5jZWxsZWQ6X1xyXG4gICAgICogY2FuY2VsbGluZyBhIHByb21pc2Ugd2l0aCBubyBhdHRhY2hlZCByZWplY3Rpb24gaGFuZGxlciBkb2VzIG5vdCB0cmlnZ2VyIGFuIHVuaGFuZGxlZCByZWplY3Rpb24gZXZlbnQuXHJcbiAgICAgKiBUaGVyZWZvcmUsIHRoZSBmb2xsb3dpbmcgaWRpb21zIGFyZSBhbGwgZXF1YWxseSBjb3JyZWN0OlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkuY2FuY2VsKCk7XHJcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYW5jZWwoKTtcclxuICAgICAqIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSkudGhlbiguLi4pLmNhdGNoKC4uLikuY2FuY2VsKCk7XHJcbiAgICAgKiBgYGBcclxuICAgICAqIFdoZW5ldmVyIHNvbWUgY2FuY2VsbGVkIHByb21pc2UgaW4gYSBjaGFpbiByZWplY3RzIHdpdGggYSBgQ2FuY2VsRXJyb3JgXHJcbiAgICAgKiB3aXRoIHRoZSBzYW1lIGNhbmNlbGxhdGlvbiBjYXVzZSBhcyBpdHNlbGYsIHRoZSBlcnJvciB3aWxsIGJlIGRpc2NhcmRlZCBzaWxlbnRseS5cclxuICAgICAqIEhvd2V2ZXIsIHRoZSBgQ2FuY2VsRXJyb3JgIF93aWxsIHN0aWxsIGJlIGRlbGl2ZXJlZF8gdG8gYWxsIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyc1xyXG4gICAgICogYWRkZWQgYnkge0BsaW5rIHRoZW59IGFuZCByZWxhdGVkIG1ldGhvZHM6XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogbGV0IGNhbmNlbGxhYmxlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcclxuICAgICAqIGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkuY2F0Y2goY29uc29sZS5sb2cpO1xyXG4gICAgICogY2FuY2VsbGFibGUuY2FuY2VsKCk7IC8vIEEgQ2FuY2VsRXJyb3IgaXMgcHJpbnRlZCB0byB0aGUgY29uc29sZS5cclxuICAgICAqIGBgYFxyXG4gICAgICogSWYgdGhlIGBDYW5jZWxFcnJvcmAgaXMgbm90IGhhbmRsZWQgZG93bnN0cmVhbSBieSB0aGUgdGltZSBpdCByZWFjaGVzXHJcbiAgICAgKiBhIF9ub24tY2FuY2VsbGVkXyBwcm9taXNlLCBpdCBfd2lsbF8gdHJpZ2dlciBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50LFxyXG4gICAgICoganVzdCBsaWtlIG5vcm1hbCByZWplY3Rpb25zIHdvdWxkOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGxldCBjYW5jZWxsYWJsZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4geyAuLi4gfSk7XHJcbiAgICAgKiBsZXQgY2hhaW5lZCA9IGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkudGhlbigoKSA9PiB7IC4uLiB9KTsgLy8gTm8gY2F0Y2guLi5cclxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBVbmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50IG9uIGNoYWluZWQhXHJcbiAgICAgKiBgYGBcclxuICAgICAqIFRoZXJlZm9yZSwgaXQgaXMgaW1wb3J0YW50IHRvIGVpdGhlciBjYW5jZWwgd2hvbGUgcHJvbWlzZSBjaGFpbnMgZnJvbSB0aGVpciB0YWlsLFxyXG4gICAgICogYXMgc2hvd24gaW4gdGhlIGNvcnJlY3QgaWRpb21zIGFib3ZlLCBvciB0YWtlIGNhcmUgb2YgaGFuZGxpbmcgZXJyb3JzIGV2ZXJ5d2hlcmUuXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgQSBjYW5jZWxsYWJsZSBwcm9taXNlIHRoYXQgX2Z1bGZpbGxzXyBhZnRlciB0aGUgY2FuY2VsIGNhbGxiYWNrIChpZiBhbnkpXHJcbiAgICAgKiBhbmQgYWxsIGhhbmRsZXJzIGF0dGFjaGVkIHVwIHRvIHRoZSBjYWxsIHRvIGNhbmNlbCBoYXZlIHJ1bi5cclxuICAgICAqIElmIHRoZSBjYW5jZWwgY2FsbGJhY2sgcmV0dXJucyBhIHRoZW5hYmxlLCB0aGUgcHJvbWlzZSByZXR1cm5lZCBieSBgY2FuY2VsYFxyXG4gICAgICogd2lsbCBhbHNvIHdhaXQgZm9yIHRoYXQgdGhlbmFibGUgdG8gc2V0dGxlLlxyXG4gICAgICogVGhpcyBlbmFibGVzIGNhbGxlcnMgdG8gd2FpdCBmb3IgdGhlIGNhbmNlbGxlZCBvcGVyYXRpb24gdG8gdGVybWluYXRlXHJcbiAgICAgKiB3aXRob3V0IGJlaW5nIGZvcmNlZCB0byBoYW5kbGUgcG90ZW50aWFsIGVycm9ycyBhdCB0aGUgY2FsbCBzaXRlLlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpLnRoZW4oKCkgPT4ge1xyXG4gICAgICogICAgIC8vIENsZWFudXAgZmluaXNoZWQsIGl0J3Mgc2FmZSB0byBkbyBzb21ldGhpbmcgZWxzZS5cclxuICAgICAqIH0sIChlcnIpID0+IHtcclxuICAgICAqICAgICAvLyBVbnJlYWNoYWJsZTogdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSBjYW5jZWwgd2lsbCBuZXZlciByZWplY3QuXHJcbiAgICAgKiB9KTtcclxuICAgICAqIGBgYFxyXG4gICAgICogTm90ZSB0aGF0IHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgX25vdF8gaGFuZGxlIGltcGxpY2l0bHkgYW55IHJlamVjdGlvblxyXG4gICAgICogdGhhdCBtaWdodCBoYXZlIG9jY3VycmVkIGFscmVhZHkgaW4gdGhlIGNhbmNlbGxlZCBjaGFpbi5cclxuICAgICAqIEl0IHdpbGwganVzdCB0cmFjayB3aGV0aGVyIHJlZ2lzdGVyZWQgaGFuZGxlcnMgaGF2ZSBiZWVuIGV4ZWN1dGVkIG9yIG5vdC5cclxuICAgICAqIFRoZXJlZm9yZSwgdW5oYW5kbGVkIHJlamVjdGlvbnMgd2lsbCBuZXZlciBiZSBzaWxlbnRseSBoYW5kbGVkIGJ5IGNhbGxpbmcgY2FuY2VsLlxyXG4gICAgICovXHJcbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XHJcbiAgICAgICAgICAgIC8vIElOVkFSSUFOVDogdGhlIHJlc3VsdCBvZiB0aGlzW2NhbmNlbEltcGxTeW1dIGFuZCB0aGUgYmFycmllciBkbyBub3QgZXZlciByZWplY3QuXHJcbiAgICAgICAgICAgIC8vIFVuZm9ydHVuYXRlbHkgbWFjT1MgSGlnaCBTaWVycmEgZG9lcyBub3Qgc3VwcG9ydCBQcm9taXNlLmFsbFNldHRsZWQuXHJcbiAgICAgICAgICAgIFByb21pc2UuYWxsKFtcclxuICAgICAgICAgICAgICAgIHRoaXNbY2FuY2VsSW1wbFN5bV0obmV3IENhbmNlbEVycm9yKFwiUHJvbWlzZSBjYW5jZWxsZWQuXCIsIHsgY2F1c2UgfSkpLFxyXG4gICAgICAgICAgICAgICAgY3VycmVudEJhcnJpZXIodGhpcylcclxuICAgICAgICAgICAgXSkudGhlbigoKSA9PiByZXNvbHZlKCksICgpID0+IHJlc29sdmUoKSk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBCaW5kcyBwcm9taXNlIGNhbmNlbGxhdGlvbiB0byB0aGUgYWJvcnQgZXZlbnQgb2YgdGhlIGdpdmVuIHtAbGluayBBYm9ydFNpZ25hbH0uXHJcbiAgICAgKiBJZiB0aGUgc2lnbmFsIGhhcyBhbHJlYWR5IGFib3J0ZWQsIHRoZSBwcm9taXNlIHdpbGwgYmUgY2FuY2VsbGVkIGltbWVkaWF0ZWx5LlxyXG4gICAgICogV2hlbiBlaXRoZXIgY29uZGl0aW9uIGlzIHZlcmlmaWVkLCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlIHdpbGwgYmUgc2V0XHJcbiAgICAgKiB0byB0aGUgc2lnbmFsJ3MgYWJvcnQgcmVhc29uIChzZWUge0BsaW5rIEFib3J0U2lnbmFsI3JlYXNvbn0pLlxyXG4gICAgICpcclxuICAgICAqIEhhcyBubyBlZmZlY3QgaWYgY2FsbGVkIChvciBpZiB0aGUgc2lnbmFsIGFib3J0cykgX2FmdGVyXyB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkLlxyXG4gICAgICogT25seSB0aGUgZmlyc3Qgc2lnbmFsIHRvIGFib3J0IHdpbGwgc2V0IHRoZSBjYW5jZWxsYXRpb24gY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogRm9yIG1vcmUgZGV0YWlscyBhYm91dCB0aGUgY2FuY2VsbGF0aW9uIHByb2Nlc3MsXHJcbiAgICAgKiBzZWUge0BsaW5rIGNhbmNlbH0gYW5kIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYCBjb25zdHJ1Y3Rvci5cclxuICAgICAqXHJcbiAgICAgKiBUaGlzIG1ldGhvZCBlbmFibGVzIGBhd2FpdGBpbmcgY2FuY2VsbGFibGUgcHJvbWlzZXMgd2l0aG91dCBoYXZpbmdcclxuICAgICAqIHRvIHN0b3JlIHRoZW0gZm9yIGZ1dHVyZSBjYW5jZWxsYXRpb24sIGUuZy46XHJcbiAgICAgKiBgYGB0c1xyXG4gICAgICogYXdhaXQgbG9uZ1J1bm5pbmdPcGVyYXRpb24oKS5jYW5jZWxPbihzaWduYWwpO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBpbnN0ZWFkIG9mOlxyXG4gICAgICogYGBgdHNcclxuICAgICAqIGxldCBwcm9taXNlVG9CZUNhbmNlbGxlZCA9IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCk7XHJcbiAgICAgKiBhd2FpdCBwcm9taXNlVG9CZUNhbmNlbGxlZDtcclxuICAgICAqIGBgYFxyXG4gICAgICpcclxuICAgICAqIEByZXR1cm5zIFRoaXMgcHJvbWlzZSwgZm9yIG1ldGhvZCBjaGFpbmluZy5cclxuICAgICAqL1xyXG4gICAgY2FuY2VsT24oc2lnbmFsOiBBYm9ydFNpZ25hbCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XHJcbiAgICAgICAgaWYgKHNpZ25hbC5hYm9ydGVkKSB7XHJcbiAgICAgICAgICAgIHZvaWQgdGhpcy5jYW5jZWwoc2lnbmFsLnJlYXNvbilcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBzaWduYWwuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pLCB7Y2FwdHVyZTogdHJ1ZX0pO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgcmV0dXJuIHRoaXM7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBBdHRhY2hlcyBjYWxsYmFja3MgZm9yIHRoZSByZXNvbHV0aW9uIGFuZC9vciByZWplY3Rpb24gb2YgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcclxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cclxuICAgICAqIFdoZW4gdGhlIHBhcmVudCBwcm9taXNlIHJlamVjdHMgb3IgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGNhbGxiYWNrIHdpbGwgcnVuLFxyXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cclxuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIEBwYXJhbSBvbmZ1bGZpbGxlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlc29sdmVkLlxyXG4gICAgICogQHBhcmFtIG9ucmVqZWN0ZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZWplY3RlZC5cclxuICAgICAqIEByZXR1cm5zIEEgYENhbmNlbGxhYmxlUHJvbWlzZWAgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHdoaWNoZXZlciBjYWxsYmFjayBpcyBleGVjdXRlZC5cclxuICAgICAqIFRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGhvb2tlZCB1cCB0byBwcm9wYWdhdGUgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHVwIHRoZSBjaGFpbiwgYnV0IG5vdCBkb3duOlxyXG4gICAgICpcclxuICAgICAqICAgLSBpZiB0aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBiZSBpbnZva2VkIHdpdGggYSBgQ2FuY2VsRXJyb3JgXHJcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIF93aWxsIHJlc29sdmUgcmVndWxhcmx5XyB3aXRoIGl0cyByZXN1bHQ7XHJcbiAgICAgKiAgIC0gY29udmVyc2VseSwgaWYgdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCBfdGhlIHBhcmVudCBwcm9taXNlIGlzIGNhbmNlbGxlZCB0b287X1xyXG4gICAgICogICAgIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIHN0aWxsIGJlIGludm9rZWQgd2l0aCB0aGUgcGFyZW50J3MgYENhbmNlbEVycm9yYCxcclxuICAgICAqICAgICBidXQgaXRzIHJlc3VsdCB3aWxsIGJlIGRpc2NhcmRlZFxyXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCB3aXRoIGEgYENhbmNlbEVycm9yYCBhcyB3ZWxsLlxyXG4gICAgICpcclxuICAgICAqIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH0gd2lsbCBmdWxmaWxsIG9ubHkgYWZ0ZXIgYWxsIGF0dGFjaGVkIGhhbmRsZXJzXHJcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cclxuICAgICAqXHJcbiAgICAgKiBJZiBlaXRoZXIgY2FsbGJhY2sgcmV0dXJucyBhIGNhbmNlbGxhYmxlIHByb21pc2UsXHJcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcclxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXHJcbiAgICAgKi9cclxuICAgIHRoZW48VFJlc3VsdDEgPSBULCBUUmVzdWx0MiA9IG5ldmVyPihvbmZ1bGZpbGxlZD86ICgodmFsdWU6IFQpID0+IFRSZXN1bHQxIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDE+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MT4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25yZWplY3RlZD86ICgocmVhc29uOiBhbnkpID0+IFRSZXN1bHQyIHwgUHJvbWlzZUxpa2U8VFJlc3VsdDI+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0Mj4pIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4ge1xyXG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XHJcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UucHJvdG90eXBlLnRoZW4gY2FsbGVkIG9uIGFuIGludmFsaWQgb2JqZWN0LlwiKTtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIC8vIE5PVEU6IFR5cGVTY3JpcHQncyBidWlsdC1pbiB0eXBlIGZvciB0aGVuIGlzIGJyb2tlbixcclxuICAgICAgICAvLyBhcyBpdCBhbGxvd3Mgc3BlY2lmeWluZyBhbiBhcmJpdHJhcnkgVFJlc3VsdDEgIT0gVCBldmVuIHdoZW4gb25mdWxmaWxsZWQgaXMgbm90IGEgZnVuY3Rpb24uXHJcbiAgICAgICAgLy8gV2UgY2Fubm90IGZpeCBpdCBpZiB3ZSB3YW50IHRvIENhbmNlbGxhYmxlUHJvbWlzZSB0byBpbXBsZW1lbnQgUHJvbWlzZUxpa2U8VD4uXHJcblxyXG4gICAgICAgIGlmICghaXNDYWxsYWJsZShvbmZ1bGZpbGxlZCkpIHsgb25mdWxmaWxsZWQgPSBpZGVudGl0eSBhcyBhbnk7IH1cclxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25yZWplY3RlZCkpIHsgb25yZWplY3RlZCA9IHRocm93ZXI7IH1cclxuXHJcbiAgICAgICAgaWYgKG9uZnVsZmlsbGVkID09PSBpZGVudGl0eSAmJiBvbnJlamVjdGVkID09IHRocm93ZXIpIHtcclxuICAgICAgICAgICAgLy8gU2hvcnRjdXQgZm9yIHRyaXZpYWwgYXJndW1lbnRzLlxyXG4gICAgICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSkgPT4gcmVzb2x2ZSh0aGlzIGFzIGFueSkpO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgY29uc3QgYmFycmllcjogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gPSB7fTtcclxuICAgICAgICB0aGlzW2JhcnJpZXJTeW1dID0gYmFycmllcjtcclxuXHJcbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VFJlc3VsdDEgfCBUUmVzdWx0Mj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICB2b2lkIHN1cGVyLnRoZW4oXHJcbiAgICAgICAgICAgICAgICAodmFsdWUpID0+IHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAodGhpc1tiYXJyaWVyU3ltXSA9PT0gYmFycmllcikgeyB0aGlzW2JhcnJpZXJTeW1dID0gbnVsbDsgfVxyXG4gICAgICAgICAgICAgICAgICAgIGJhcnJpZXIucmVzb2x2ZT8uKCk7XHJcblxyXG4gICAgICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlc29sdmUob25mdWxmaWxsZWQhKHZhbHVlKSk7XHJcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJlamVjdChlcnIpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0sXHJcbiAgICAgICAgICAgICAgICAocmVhc29uPykgPT4ge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmICh0aGlzW2JhcnJpZXJTeW1dID09PSBiYXJyaWVyKSB7IHRoaXNbYmFycmllclN5bV0gPSBudWxsOyB9XHJcbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcclxuXHJcbiAgICAgICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbnJlamVjdGVkIShyZWFzb24pKTtcclxuICAgICAgICAgICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmVqZWN0KGVycik7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICApO1xyXG4gICAgICAgIH0sIGFzeW5jIChjYXVzZT8pID0+IHtcclxuICAgICAgICAgICAgLy9jYW5jZWxsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuIG9uY2FuY2VsbGVkPy4oY2F1c2UpO1xyXG4gICAgICAgICAgICB9IGZpbmFsbHkge1xyXG4gICAgICAgICAgICAgICAgYXdhaXQgdGhpcy5jYW5jZWwoY2F1c2UpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIGZvciBvbmx5IHRoZSByZWplY3Rpb24gb2YgdGhlIFByb21pc2UuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxyXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXHJcbiAgICAgKiBfZXZlbiBhZnRlciB0aGUgcmV0dXJuZWQgcHJvbWlzZSBoYXMgYmVlbiBjYW5jZWxsZWQ6X1xyXG4gICAgICogaW4gdGhhdCBjYXNlLCBzaG91bGQgaXQgcmVqZWN0IG9yIHRocm93LCB0aGUgcmVhc29uIHdpbGwgYmUgd3JhcHBlZFxyXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogSXQgaXMgZXF1aXZhbGVudCB0b1xyXG4gICAgICogYGBgdHNcclxuICAgICAqIGNhbmNlbGxhYmxlUHJvbWlzZS50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xyXG4gICAgICogYGBgXHJcbiAgICAgKiBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cclxuICAgICAqXHJcbiAgICAgKiBAcmV0dXJucyBBIFByb21pc2UgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHRoZSBjYWxsYmFjay5cclxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxyXG4gICAgICogd2lsbCBwcm9wYWdhdGUgdXAgdGhlIGNoYWluIHRvIHRoZSBwYXJlbnQgcHJvbWlzZSxcclxuICAgICAqIGJ1dCBub3QgaW4gdGhlIG90aGVyIGRpcmVjdGlvbi5cclxuICAgICAqXHJcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xyXG4gICAgICogdXAgdGhlIGVudGlyZSBwcm9taXNlIGNoYWluIGhhdmUgYmVlbiBydW4uXHJcbiAgICAgKlxyXG4gICAgICogSWYgYG9ucmVqZWN0ZWRgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXHJcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxyXG4gICAgICovXHJcbiAgICBjYXRjaDxUUmVzdWx0ID0gbmV2ZXI+KG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiAoUHJvbWlzZUxpa2U8VFJlc3VsdD4gfCBUUmVzdWx0KSkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUIHwgVFJlc3VsdD4ge1xyXG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4odW5kZWZpbmVkLCBvbnJlamVjdGVkLCBvbmNhbmNlbGxlZCk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBBdHRhY2hlcyBhIGNhbGxiYWNrIHRoYXQgaXMgaW52b2tlZCB3aGVuIHRoZSBDYW5jZWxsYWJsZVByb21pc2UgaXMgc2V0dGxlZCAoZnVsZmlsbGVkIG9yIHJlamVjdGVkKS4gVGhlXHJcbiAgICAgKiByZXNvbHZlZCB2YWx1ZSBjYW5ub3QgYmUgYWNjZXNzZWQgb3IgbW9kaWZpZWQgZnJvbSB0aGUgY2FsbGJhY2suXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHNldHRsZSBpbiB0aGUgc2FtZSBzdGF0ZSBhcyB0aGUgb3JpZ2luYWwgb25lXHJcbiAgICAgKiBhZnRlciB0aGUgcHJvdmlkZWQgY2FsbGJhY2sgaGFzIGNvbXBsZXRlZCBleGVjdXRpb24sXHJcbiAgICAgKiB1bmxlc3MgdGhlIGNhbGxiYWNrIHRocm93cyBvciByZXR1cm5zIGEgcmVqZWN0aW5nIHByb21pc2UsXHJcbiAgICAgKiBpbiB3aGljaCBjYXNlIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IGFzIHdlbGwuXHJcbiAgICAgKlxyXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxyXG4gICAgICogd2l0aCB0aGUgc2FtZSBzZW1hbnRpY3MgYXMgdGhlIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgb2YgdGhlIGNvbnN0cnVjdG9yLlxyXG4gICAgICogT25jZSB0aGUgcGFyZW50IHByb21pc2Ugc2V0dGxlcywgdGhlIGBvbmZpbmFsbHlgIGNhbGxiYWNrIHdpbGwgcnVuLFxyXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cclxuICAgICAqIGluIHRoYXQgY2FzZSwgc2hvdWxkIGl0IHJlamVjdCBvciB0aHJvdywgdGhlIHJlYXNvbiB3aWxsIGJlIHdyYXBwZWRcclxuICAgICAqIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBhbmQgYnViYmxlZCB1cCBhcyBhbiB1bmhhbmRsZWQgcmVqZWN0aW9uLlxyXG4gICAgICpcclxuICAgICAqIFRoaXMgbWV0aG9kIGlzIGltcGxlbWVudGVkIGluIHRlcm1zIG9mIHtAbGluayB0aGVufSBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cclxuICAgICAqIEl0IGlzIHBvbHlmaWxsZWQsIGhlbmNlIGF2YWlsYWJsZSBpbiBldmVyeSBPUy93ZWJ2aWV3IHZlcnNpb24uXHJcbiAgICAgKlxyXG4gICAgICogQHJldHVybnMgQSBQcm9taXNlIGZvciB0aGUgY29tcGxldGlvbiBvZiB0aGUgY2FsbGJhY2suXHJcbiAgICAgKiBDYW5jZWxsYXRpb24gcmVxdWVzdHMgb24gdGhlIHJldHVybmVkIHByb21pc2VcclxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXHJcbiAgICAgKiBidXQgbm90IGluIHRoZSBvdGhlciBkaXJlY3Rpb24uXHJcbiAgICAgKlxyXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcclxuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxyXG4gICAgICpcclxuICAgICAqIElmIGBvbmZpbmFsbHlgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxyXG4gICAgICogY2FuY2VsbGF0aW9uIHJlcXVlc3RzIHdpbGwgYmUgZGl2ZXJ0ZWQgdG8gaXQsXHJcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxyXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxyXG4gICAgICovXHJcbiAgICBmaW5hbGx5KG9uZmluYWxseT86ICgoKSA9PiB2b2lkKSB8IHVuZGVmaW5lZCB8IG51bGwsIG9uY2FuY2VsbGVkPzogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xyXG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS5maW5hbGx5IGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25maW5hbGx5KSkge1xyXG4gICAgICAgICAgICByZXR1cm4gdGhpcy50aGVuKG9uZmluYWxseSwgb25maW5hbGx5LCBvbmNhbmNlbGxlZCk7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICByZXR1cm4gdGhpcy50aGVuKFxyXG4gICAgICAgICAgICAodmFsdWUpID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHZhbHVlKSxcclxuICAgICAgICAgICAgKHJlYXNvbj8pID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHsgdGhyb3cgcmVhc29uOyB9KSxcclxuICAgICAgICAgICAgb25jYW5jZWxsZWQsXHJcbiAgICAgICAgKTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIFdlIHVzZSB0aGUgYFtTeW1ib2wuc3BlY2llc11gIHN0YXRpYyBwcm9wZXJ0eSwgaWYgYXZhaWxhYmxlLFxyXG4gICAgICogdG8gZGlzYWJsZSB0aGUgYnVpbHQtaW4gYXV0b21hdGljIHN1YmNsYXNzaW5nIGZlYXR1cmVzIGZyb20ge0BsaW5rIFByb21pc2V9LlxyXG4gICAgICogSXQgaXMgY3JpdGljYWwgZm9yIHBlcmZvcm1hbmNlIHJlYXNvbnMgdGhhdCBleHRlbmRlcnMgZG8gbm90IG92ZXJyaWRlIHRoaXMuXHJcbiAgICAgKiBPbmNlIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nXHJcbiAgICAgKiBpcyBlaXRoZXIgYWNjZXB0ZWQgb3IgcmV0aXJlZCwgdGhpcyBpbXBsZW1lbnRhdGlvbiB3aWxsIGhhdmUgdG8gYmUgcmV2aXNlZCBhY2NvcmRpbmdseS5cclxuICAgICAqXHJcbiAgICAgKiBAaWdub3JlXHJcbiAgICAgKiBAaW50ZXJuYWxcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGdldCBbc3BlY2llc10oKSB7XHJcbiAgICAgICAgcmV0dXJuIFByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXHJcbiAgICAgKiB3aGVuIGFsbCBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgcmVzb2x2ZSwgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgUHJvbWlzZSBpcyByZWplY3RlZC5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyBhbGw8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD5bXT47XHJcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBBd2FpdGVkPFRbUF0+OyB9PjtcclxuICAgIHN0YXRpYyBhbGw8VCBleHRlbmRzIEl0ZXJhYmxlPHVua25vd24+IHwgQXJyYXlMaWtlPHVua25vd24+Pih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4ge1xyXG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XHJcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcclxuICAgICAgICAgICAgPyBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShjb2xsZWN0ZWQpXHJcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xyXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUgb3IgcmVqZWN0LlxyXG4gICAgICpcclxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VD4+W10+O1xyXG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VFtQXT4+OyB9PjtcclxuICAgIHN0YXRpYyBhbGxTZXR0bGVkPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXHJcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxyXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbFNldHRsZWQoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XHJcbiAgICAgICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBUaGUgYW55IGZ1bmN0aW9uIHJldHVybnMgYSBwcm9taXNlIHRoYXQgaXMgZnVsZmlsbGVkIGJ5IHRoZSBmaXJzdCBnaXZlbiBwcm9taXNlIHRvIGJlIGZ1bGZpbGxlZCxcclxuICAgICAqIG9yIHJlamVjdGVkIHdpdGggYW4gQWdncmVnYXRlRXJyb3IgY29udGFpbmluZyBhbiBhcnJheSBvZiByZWplY3Rpb24gcmVhc29uc1xyXG4gICAgICogaWYgYWxsIG9mIHRoZSBnaXZlbiBwcm9taXNlcyBhcmUgcmVqZWN0ZWQuXHJcbiAgICAgKiBJdCByZXNvbHZlcyBhbGwgZWxlbWVudHMgb2YgdGhlIHBhc3NlZCBpdGVyYWJsZSB0byBwcm9taXNlcyBhcyBpdCBydW5zIHRoaXMgYWxnb3JpdGhtLlxyXG4gICAgICpcclxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XHJcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGFueTxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XHJcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XHJcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcclxuICAgICAgICBsZXQgY29sbGVjdGVkID0gQXJyYXkuZnJvbSh2YWx1ZXMpO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXHJcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxyXG4gICAgICAgICAgICA6IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFueShjb2xsZWN0ZWQpLnRoZW4ocmVzb2x2ZSwgcmVqZWN0KTtcclxuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcclxuICAgICAgICByZXR1cm4gcHJvbWlzZTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgb3IgcmVqZWN0ZWQgd2hlbiBhbnkgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIGFyZSByZXNvbHZlZCBvciByZWplY3RlZC5cclxuICAgICAqXHJcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxyXG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXHJcbiAgICAgKlxyXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXHJcbiAgICAgKi9cclxuICAgIHN0YXRpYyByYWNlPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VFtudW1iZXJdPj47XHJcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XHJcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcclxuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5yYWNlKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xyXG4gICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XHJcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGNhbmNlbGxlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCBjYXVzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIGNhbmNlbDxUID0gbmV2ZXI+KGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICBjb25zdCBwID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoKSA9PiB7fSk7XHJcbiAgICAgICAgcC5jYW5jZWwoY2F1c2UpO1xyXG4gICAgICAgIHJldHVybiBwO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCBjYW5jZWxzXHJcbiAgICAgKiBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQsIHdpdGggdGhlIHByb3ZpZGVkIGNhdXNlLlxyXG4gICAgICpcclxuICAgICAqIElmIHRoZSB7QGxpbmsgQWJvcnRTaWduYWwudGltZW91dH0gZmFjdG9yeSBtZXRob2QgaXMgYXZhaWxhYmxlLFxyXG4gICAgICogaXQgaXMgdXNlZCB0byBiYXNlIHRoZSB0aW1lb3V0IG9uIF9hY3RpdmVfIHRpbWUgcmF0aGVyIHRoYW4gX2VsYXBzZWRfIHRpbWUuXHJcbiAgICAgKiBPdGhlcndpc2UsIGB0aW1lb3V0YCBmYWxscyBiYWNrIHRvIHtAbGluayBzZXRUaW1lb3V0fS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHRpbWVvdXQ8VCA9IG5ldmVyPihtaWxsaXNlY29uZHM6IG51bWJlciwgY2F1c2U/OiBhbnkpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KCgpID0+IHt9KTtcclxuICAgICAgICBpZiAoQWJvcnRTaWduYWwgJiYgdHlwZW9mIEFib3J0U2lnbmFsID09PSAnZnVuY3Rpb24nICYmIEFib3J0U2lnbmFsLnRpbWVvdXQgJiYgdHlwZW9mIEFib3J0U2lnbmFsLnRpbWVvdXQgPT09ICdmdW5jdGlvbicpIHtcclxuICAgICAgICAgICAgQWJvcnRTaWduYWwudGltZW91dChtaWxsaXNlY29uZHMpLmFkZEV2ZW50TGlzdGVuZXIoJ2Fib3J0JywgKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSksIG1pbGxpc2Vjb25kcyk7XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHJldHVybiBwcm9taXNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQuXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgc2xlZXAobWlsbGlzZWNvbmRzOiBudW1iZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8dm9pZD47XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgcmVzb2x2ZXMgYWZ0ZXJcclxuICAgICAqIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgdmFsdWUuXHJcbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgc2xlZXA8VD4obWlsbGlzZWNvbmRzOiBudW1iZXIsIHZhbHVlOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xyXG4gICAgc3RhdGljIHNsZWVwPFQgPSB2b2lkPihtaWxsaXNlY29uZHM6IG51bWJlciwgdmFsdWU/OiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSkgPT4ge1xyXG4gICAgICAgICAgICBzZXRUaW1lb3V0KCgpID0+IHJlc29sdmUodmFsdWUhKSwgbWlsbGlzZWNvbmRzKTtcclxuICAgICAgICB9KTtcclxuICAgIH1cclxuXHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVqZWN0ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgcmVhc29uLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmVqZWN0PFQgPSBuZXZlcj4ocmVhc29uPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigoXywgcmVqZWN0KSA9PiByZWplY3QocmVhc29uKSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlc29sdmUoKTogQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+O1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxyXG4gICAgICpcclxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xyXG4gICAgICovXHJcbiAgICBzdGF0aWMgcmVzb2x2ZTxUPih2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcclxuICAgIHN0YXRpYyByZXNvbHZlPFQgPSB2b2lkPih2YWx1ZT86IFQgfCBQcm9taXNlTGlrZTxUPik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PiB7XHJcbiAgICAgICAgaWYgKHZhbHVlIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSB7XHJcbiAgICAgICAgICAgIC8vIE9wdGltaXNlIGZvciBjYW5jZWxsYWJsZSBwcm9taXNlcy5cclxuICAgICAgICAgICAgcmV0dXJuIHZhbHVlO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxhbnk+KChyZXNvbHZlKSA9PiByZXNvbHZlKHZhbHVlKSk7XHJcbiAgICB9XHJcblxyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSBhbmQgcmV0dXJucyBpdCBpbiBhbiBvYmplY3QsIGFsb25nIHdpdGggaXRzIHJlc29sdmUgYW5kIHJlamVjdCBmdW5jdGlvbnNcclxuICAgICAqIGFuZCBhIGdldHRlci9zZXR0ZXIgZm9yIHRoZSBjYW5jZWxsYXRpb24gY2FsbGJhY2suXHJcbiAgICAgKlxyXG4gICAgICogVGhpcyBtZXRob2QgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cclxuICAgICAqXHJcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcclxuICAgICAqL1xyXG4gICAgc3RhdGljIHdpdGhSZXNvbHZlcnM8VD4oKTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XHJcbiAgICAgICAgbGV0IHJlc3VsdDogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHsgb25jYW5jZWxsZWQ6IG51bGwgfSBhcyBhbnk7XHJcbiAgICAgICAgcmVzdWx0LnByb21pc2UgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICAgICAgcmVzdWx0LnJlc29sdmUgPSByZXNvbHZlO1xyXG4gICAgICAgICAgICByZXN1bHQucmVqZWN0ID0gcmVqZWN0O1xyXG4gICAgICAgIH0sIChjYXVzZT86IGFueSkgPT4geyByZXN1bHQub25jYW5jZWxsZWQ/LihjYXVzZSk7IH0pO1xyXG4gICAgICAgIHJldHVybiByZXN1bHQ7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSBjYW5jZWxsYXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cclxuICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgcmVzdWx0aW5nIGZ1bmN0aW9uIGRvZXMgbm90IHJlamVjdC5cclxuICovXHJcbmZ1bmN0aW9uIGNhbmNlbGxlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpIHtcclxuICAgIGxldCBjYW5jZWxsYXRpb25Qcm9taXNlOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPSB1bmRlZmluZWQ7XHJcblxyXG4gICAgcmV0dXJuIChyZWFzb246IENhbmNlbEVycm9yKTogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+ID0+IHtcclxuICAgICAgICBpZiAoIXN0YXRlLnNldHRsZWQpIHtcclxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgIHN0YXRlLnJlYXNvbiA9IHJlYXNvbjtcclxuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcclxuXHJcbiAgICAgICAgICAgIC8vIEF0dGFjaCBhbiBlcnJvciBoYW5kbGVyIHRoYXQgaWdub3JlcyB0aGlzIHNwZWNpZmljIHJlamVjdGlvbiByZWFzb24gYW5kIG5vdGhpbmcgZWxzZS5cclxuICAgICAgICAgICAgLy8gSW4gdGhlb3J5LCBhIHNhbmUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbiBhdCB0aGlzIHBvaW50XHJcbiAgICAgICAgICAgIC8vIHNob3VsZCBhbHdheXMgcmVqZWN0IHdpdGggb3VyIGNhbmNlbGxhdGlvbiByZWFzb24sXHJcbiAgICAgICAgICAgIC8vIGhlbmNlIHRoZSBoYW5kbGVyIHdpbGwgbmV2ZXIgdGhyb3cuXHJcbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5wcm90b3R5cGUudGhlbi5jYWxsKHByb21pc2UucHJvbWlzZSwgdW5kZWZpbmVkLCAoZXJyKSA9PiB7XHJcbiAgICAgICAgICAgICAgICBpZiAoZXJyICE9PSByZWFzb24pIHtcclxuICAgICAgICAgICAgICAgICAgICB0aHJvdyBlcnI7XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIH0pO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgLy8gSWYgcmVhc29uIGlzIG5vdCBzZXQsIHRoZSBwcm9taXNlIHJlc29sdmVkIHJlZ3VsYXJseSwgaGVuY2Ugd2UgbXVzdCBub3QgY2FsbCBvbmNhbmNlbGxlZC5cclxuICAgICAgICAvLyBJZiBvbmNhbmNlbGxlZCBpcyB1bnNldCwgbm8gbmVlZCB0byBnbyBhbnkgZnVydGhlci5cclxuICAgICAgICBpZiAoIXN0YXRlLnJlYXNvbiB8fCAhcHJvbWlzZS5vbmNhbmNlbGxlZCkgeyByZXR1cm47IH1cclxuXHJcbiAgICAgICAgY2FuY2VsbGF0aW9uUHJvbWlzZSA9IG5ldyBQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XHJcbiAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICByZXNvbHZlKHByb21pc2Uub25jYW5jZWxsZWQhKHN0YXRlLnJlYXNvbiEuY2F1c2UpKTtcclxuICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSkuY2F0Y2goKHJlYXNvbj8pID0+IHtcclxuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcclxuICAgICAgICB9KTtcclxuXHJcbiAgICAgICAgLy8gVW5zZXQgb25jYW5jZWxsZWQgdG8gcHJldmVudCByZXBlYXRlZCBjYWxscy5cclxuICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gbnVsbDtcclxuXHJcbiAgICAgICAgcmV0dXJuIGNhbmNlbGxhdGlvblByb21pc2U7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZXR1cm5zIGEgY2FsbGJhY2sgdGhhdCBpbXBsZW1lbnRzIHRoZSByZXNvbHV0aW9uIGFsZ29yaXRobSBmb3IgdGhlIGdpdmVuIGNhbmNlbGxhYmxlIHByb21pc2UuXHJcbiAqL1xyXG5mdW5jdGlvbiByZXNvbHZlckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiB7XHJcbiAgICByZXR1cm4gKHZhbHVlKSA9PiB7XHJcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cclxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xyXG5cclxuICAgICAgICBpZiAodmFsdWUgPT09IHByb21pc2UucHJvbWlzZSkge1xyXG4gICAgICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cclxuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XHJcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJBIHByb21pc2UgY2Fubm90IGJlIHJlc29sdmVkIHdpdGggaXRzZWxmLlwiKSk7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcblxyXG4gICAgICAgIGlmICh2YWx1ZSAhPSBudWxsICYmICh0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ2Z1bmN0aW9uJykpIHtcclxuICAgICAgICAgICAgbGV0IHRoZW46IGFueTtcclxuICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgIHRoZW4gPSAodmFsdWUgYXMgYW55KS50aGVuO1xyXG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcclxuICAgICAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QoZXJyKTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgaWYgKGlzQ2FsbGFibGUodGhlbikpIHtcclxuICAgICAgICAgICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgICAgICAgICAgbGV0IGNhbmNlbCA9ICh2YWx1ZSBhcyBhbnkpLmNhbmNlbDtcclxuICAgICAgICAgICAgICAgICAgICBpZiAoaXNDYWxsYWJsZShjYW5jZWwpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IG9uY2FuY2VsbGVkID0gKGNhdXNlPzogYW55KSA9PiB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xyXG4gICAgICAgICAgICAgICAgICAgICAgICB9O1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoc3RhdGUucmVhc29uKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBJZiBhbHJlYWR5IGNhbmNlbGxlZCwgcHJvcGFnYXRlIGNhbmNlbGxhdGlvbi5cclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIC8vIFRoZSBwcm9taXNlIHJldHVybmVkIGZyb20gdGhlIGNhbmNlbGxlciBhbGdvcml0aG0gZG9lcyBub3QgcmVqZWN0XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBzbyBpdCBjYW4gYmUgZGlzY2FyZGVkIHNhZmVseS5cclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgY2FuY2VsbGVyRm9yKHsgLi4ucHJvbWlzZSwgb25jYW5jZWxsZWQgfSwgc3RhdGUpKHN0YXRlLnJlYXNvbik7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gb25jYW5jZWxsZWQ7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgICAgICAgICAgICAgY29uc3QgbmV3U3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlID0ge1xyXG4gICAgICAgICAgICAgICAgICAgIHJvb3Q6IHN0YXRlLnJvb3QsXHJcbiAgICAgICAgICAgICAgICAgICAgcmVzb2x2aW5nOiBmYWxzZSxcclxuICAgICAgICAgICAgICAgICAgICBnZXQgc2V0dGxlZCgpIHsgcmV0dXJuIHRoaXMucm9vdC5zZXR0bGVkIH0sXHJcbiAgICAgICAgICAgICAgICAgICAgc2V0IHNldHRsZWQodmFsdWUpIHsgdGhpcy5yb290LnNldHRsZWQgPSB2YWx1ZTsgfSxcclxuICAgICAgICAgICAgICAgICAgICBnZXQgcmVhc29uKCkgeyByZXR1cm4gdGhpcy5yb290LnJlYXNvbiB9XHJcbiAgICAgICAgICAgICAgICB9O1xyXG5cclxuICAgICAgICAgICAgICAgIGNvbnN0IHJlamVjdG9yID0gcmVqZWN0b3JGb3IocHJvbWlzZSwgbmV3U3RhdGUpO1xyXG4gICAgICAgICAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KHRoZW4sIHZhbHVlLCBbcmVzb2x2ZXJGb3IocHJvbWlzZSwgbmV3U3RhdGUpLCByZWplY3Rvcl0pO1xyXG4gICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIHJldHVybjsgLy8gSU1QT1JUQU5UIVxyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkgeyByZXR1cm47IH1cclxuICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcclxuICAgICAgICBwcm9taXNlLnJlc29sdmUodmFsdWUpO1xyXG4gICAgfTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlamVjdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxyXG4gKi9cclxuZnVuY3Rpb24gcmVqZWN0b3JGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3Ige1xyXG4gICAgcmV0dXJuIChyZWFzb24/KSA9PiB7XHJcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cclxuICAgICAgICBzdGF0ZS5yZXNvbHZpbmcgPSB0cnVlO1xyXG5cclxuICAgICAgICBpZiAoc3RhdGUuc2V0dGxlZCkge1xyXG4gICAgICAgICAgICB0cnkge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIHN0YXRlLnJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIE9iamVjdC5pcyhyZWFzb24uY2F1c2UsIHN0YXRlLnJlYXNvbi5jYXVzZSkpIHtcclxuICAgICAgICAgICAgICAgICAgICAvLyBTd2FsbG93IGxhdGUgcmVqZWN0aW9ucyB0aGF0IGFyZSBDYW5jZWxFcnJvcnMgd2hvc2UgY2FuY2VsbGF0aW9uIGNhdXNlIGlzIHRoZSBzYW1lIGFzIG91cnMuXHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbikpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xyXG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENhbmNlbHMgYWxsIHZhbHVlcyBpbiBhbiBhcnJheSB0aGF0IGxvb2sgbGlrZSBjYW5jZWxsYWJsZSB0aGVuYWJsZXMuXHJcbiAqIFJldHVybnMgYSBwcm9taXNlIHRoYXQgZnVsZmlsbHMgb25jZSBhbGwgY2FuY2VsbGF0aW9uIHByb2NlZHVyZXMgZm9yIHRoZSBnaXZlbiB2YWx1ZXMgaGF2ZSBzZXR0bGVkLlxyXG4gKi9cclxuZnVuY3Rpb24gY2FuY2VsQWxsKHBhcmVudDogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCB2YWx1ZXM6IGFueVtdLCBjYXVzZT86IGFueSk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgY29uc3QgcmVzdWx0czogUHJvbWlzZTx2b2lkPltdID0gW107XHJcblxyXG4gICAgZm9yIChjb25zdCB2YWx1ZSBvZiB2YWx1ZXMpIHtcclxuICAgICAgICBsZXQgY2FuY2VsOiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXI7XHJcbiAgICAgICAgdHJ5IHtcclxuICAgICAgICAgICAgaWYgKCFpc0NhbGxhYmxlKHZhbHVlLnRoZW4pKSB7IGNvbnRpbnVlOyB9XHJcbiAgICAgICAgICAgIGNhbmNlbCA9IHZhbHVlLmNhbmNlbDtcclxuICAgICAgICAgICAgaWYgKCFpc0NhbGxhYmxlKGNhbmNlbCkpIHsgY29udGludWU7IH1cclxuICAgICAgICB9IGNhdGNoIHsgY29udGludWU7IH1cclxuXHJcbiAgICAgICAgbGV0IHJlc3VsdDogdm9pZCB8IFByb21pc2VMaWtlPHZvaWQ+O1xyXG4gICAgICAgIHRyeSB7XHJcbiAgICAgICAgICAgIHJlc3VsdCA9IFJlZmxlY3QuYXBwbHkoY2FuY2VsLCB2YWx1ZSwgW2NhdXNlXSk7XHJcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XHJcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcclxuICAgICAgICAgICAgY29udGludWU7XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAoIXJlc3VsdCkgeyBjb250aW51ZTsgfVxyXG4gICAgICAgIHJlc3VsdHMucHVzaChcclxuICAgICAgICAgICAgKHJlc3VsdCBpbnN0YW5jZW9mIFByb21pc2UgID8gcmVzdWx0IDogUHJvbWlzZS5yZXNvbHZlKHJlc3VsdCkpLmNhdGNoKChyZWFzb24/KSA9PiB7XHJcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocGFyZW50LCByZWFzb24sIFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWwgbWV0aG9kLlwiKSk7XHJcbiAgICAgICAgICAgIH0pXHJcbiAgICAgICAgKTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gUHJvbWlzZS5hbGwocmVzdWx0cykgYXMgYW55O1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBpdHMgYXJndW1lbnQuXHJcbiAqL1xyXG5mdW5jdGlvbiBpZGVudGl0eTxUPih4OiBUKTogVCB7XHJcbiAgICByZXR1cm4geDtcclxufVxyXG5cclxuLyoqXHJcbiAqIFRocm93cyBpdHMgYXJndW1lbnQuXHJcbiAqL1xyXG5mdW5jdGlvbiB0aHJvd2VyKHJlYXNvbj86IGFueSk6IG5ldmVyIHtcclxuICAgIHRocm93IHJlYXNvbjtcclxufVxyXG5cclxuLyoqXHJcbiAqIEF0dGVtcHRzIHZhcmlvdXMgc3RyYXRlZ2llcyB0byBjb252ZXJ0IGFuIGVycm9yIHRvIGEgc3RyaW5nLlxyXG4gKi9cclxuZnVuY3Rpb24gZXJyb3JNZXNzYWdlKGVycjogYW55KTogc3RyaW5nIHtcclxuICAgIHRyeSB7XHJcbiAgICAgICAgaWYgKGVyciBpbnN0YW5jZW9mIEVycm9yIHx8IHR5cGVvZiBlcnIgIT09ICdvYmplY3QnIHx8IGVyci50b1N0cmluZyAhPT0gT2JqZWN0LnByb3RvdHlwZS50b1N0cmluZykge1xyXG4gICAgICAgICAgICByZXR1cm4gXCJcIiArIGVycjtcclxuICAgICAgICB9XHJcbiAgICB9IGNhdGNoIHt9XHJcblxyXG4gICAgdHJ5IHtcclxuICAgICAgICByZXR1cm4gSlNPTi5zdHJpbmdpZnkoZXJyKTtcclxuICAgIH0gY2F0Y2gge31cclxuXHJcbiAgICB0cnkge1xyXG4gICAgICAgIHJldHVybiBPYmplY3QucHJvdG90eXBlLnRvU3RyaW5nLmNhbGwoZXJyKTtcclxuICAgIH0gY2F0Y2gge31cclxuXHJcbiAgICByZXR1cm4gXCI8Y291bGQgbm90IGNvbnZlcnQgZXJyb3IgdG8gc3RyaW5nPlwiO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgY3VycmVudCBiYXJyaWVyIHByb21pc2UgZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLiBJZiBuZWNlc3NhcnksIGluaXRpYWxpc2VzIHRoZSBiYXJyaWVyLlxyXG4gKi9cclxuZnVuY3Rpb24gY3VycmVudEJhcnJpZXI8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+KTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICBsZXQgcHdyOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiA9IHByb21pc2VbYmFycmllclN5bV0gPz8ge307XHJcbiAgICBpZiAoISgncHJvbWlzZScgaW4gcHdyKSkge1xyXG4gICAgICAgIE9iamVjdC5hc3NpZ24ocHdyLCBwcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPigpKTtcclxuICAgIH1cclxuICAgIGlmIChwcm9taXNlW2JhcnJpZXJTeW1dID09IG51bGwpIHtcclxuICAgICAgICBwd3IucmVzb2x2ZSEoKTtcclxuICAgICAgICBwcm9taXNlW2JhcnJpZXJTeW1dID0gcHdyO1xyXG4gICAgfVxyXG4gICAgcmV0dXJuIHB3ci5wcm9taXNlITtcclxufVxyXG5cclxuLy8gUG9seWZpbGwgUHJvbWlzZS53aXRoUmVzb2x2ZXJzLlxyXG5sZXQgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBQcm9taXNlLndpdGhSZXNvbHZlcnM7XHJcbmlmIChwcm9taXNlV2l0aFJlc29sdmVycyAmJiB0eXBlb2YgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPT09ICdmdW5jdGlvbicpIHtcclxuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gcHJvbWlzZVdpdGhSZXNvbHZlcnMuYmluZChQcm9taXNlKTtcclxufSBlbHNlIHtcclxuICAgIHByb21pc2VXaXRoUmVzb2x2ZXJzID0gZnVuY3Rpb24gPFQ+KCk6IFByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+IHtcclxuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xyXG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xyXG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBuZXcgUHJvbWlzZTxUPigocmVzLCByZWopID0+IHsgcmVzb2x2ZSA9IHJlczsgcmVqZWN0ID0gcmVqOyB9KTtcclxuICAgICAgICByZXR1cm4geyBwcm9taXNlLCByZXNvbHZlLCByZWplY3QgfTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5cclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ2xpcGJvYXJkKTtcclxuXHJcbmNvbnN0IENsaXBib2FyZFNldFRleHQgPSAwO1xyXG5jb25zdCBDbGlwYm9hcmRUZXh0ID0gMTtcclxuXHJcbi8qKlxyXG4gKiBTZXRzIHRoZSB0ZXh0IHRvIHRoZSBDbGlwYm9hcmQuXHJcbiAqXHJcbiAqIEBwYXJhbSB0ZXh0IC0gVGhlIHRleHQgdG8gYmUgc2V0IHRvIHRoZSBDbGlwYm9hcmQuXHJcbiAqIEByZXR1cm4gQSBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2hlbiB0aGUgb3BlcmF0aW9uIGlzIHN1Y2Nlc3NmdWwuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0OiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKTogUHJvbWlzZTxzdHJpbmc+IHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFRleHQpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG5leHBvcnQgaW50ZXJmYWNlIFNpemUge1xyXG4gICAgLyoqIFRoZSB3aWR0aCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXHJcbiAgICBXaWR0aDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgYSByZWN0YW5ndWxhciBhcmVhLiAqL1xyXG4gICAgSGVpZ2h0OiBudW1iZXI7XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgUmVjdCB7XHJcbiAgICAvKiogVGhlIFggY29vcmRpbmF0ZSBvZiB0aGUgb3JpZ2luLiAqL1xyXG4gICAgWDogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cclxuICAgIFk6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgd2lkdGggb2YgdGhlIHJlY3RhbmdsZS4gKi9cclxuICAgIFdpZHRoOiBudW1iZXI7XHJcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xyXG4gICAgSGVpZ2h0OiBudW1iZXI7XHJcbn1cclxuXHJcbmV4cG9ydCBpbnRlcmZhY2UgU2NyZWVuIHtcclxuICAgIC8qKiBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi4gKi9cclxuICAgIElEOiBzdHJpbmc7XHJcbiAgICAvKiogSHVtYW4tcmVhZGFibGUgbmFtZSBvZiB0aGUgc2NyZWVuLiAqL1xyXG4gICAgTmFtZTogc3RyaW5nO1xyXG4gICAgLyoqIFRoZSBzY2FsZSBmYWN0b3Igb2YgdGhlIHNjcmVlbiAoRFBJLzk2KS4gMSA9IHN0YW5kYXJkIERQSSwgMiA9IEhpRFBJIChSZXRpbmEpLCBldGMuICovXHJcbiAgICBTY2FsZUZhY3RvcjogbnVtYmVyO1xyXG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFg6IG51bWJlcjtcclxuICAgIC8qKiBUaGUgWSBjb29yZGluYXRlIG9mIHRoZSBzY3JlZW4uICovXHJcbiAgICBZOiBudW1iZXI7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIHdpZHRoIGFuZCBoZWlnaHQgb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFNpemU6IFNpemU7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LiAqL1xyXG4gICAgQm91bmRzOiBSZWN0O1xyXG4gICAgLyoqIENvbnRhaW5zIHRoZSBwaHlzaWNhbCBib3VuZHMgb2YgdGhlIHNjcmVlbiBpbiB0ZXJtcyBvZiBYLCBZLCBXaWR0aCwgYW5kIEhlaWdodCAoYmVmb3JlIHNjYWxpbmcpLiAqL1xyXG4gICAgUGh5c2ljYWxCb3VuZHM6IFJlY3Q7XHJcbiAgICAvKiogQ29udGFpbnMgdGhlIGFyZWEgb2YgdGhlIHNjcmVlbiB0aGF0IGlzIGFjdHVhbGx5IHVzYWJsZSAoZXhjbHVkaW5nIHRhc2tiYXIgYW5kIG90aGVyIHN5c3RlbSBVSSkuICovXHJcbiAgICBXb3JrQXJlYTogUmVjdDtcclxuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgV29ya0FyZWEgb2YgdGhlIHNjcmVlbiAoYmVmb3JlIHNjYWxpbmcpLiAqL1xyXG4gICAgUGh5c2ljYWxXb3JrQXJlYTogUmVjdDtcclxuICAgIC8qKiBUcnVlIGlmIHRoaXMgaXMgdGhlIHByaW1hcnkgbW9uaXRvciBzZWxlY3RlZCBieSB0aGUgdXNlciBpbiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cclxuICAgIElzUHJpbWFyeTogYm9vbGVhbjtcclxuICAgIC8qKiBUaGUgcm90YXRpb24gb2YgdGhlIHNjcmVlbi4gKi9cclxuICAgIFJvdGF0aW9uOiBudW1iZXI7XHJcbn1cclxuXHJcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TY3JlZW5zKTtcclxuXHJcbmNvbnN0IGdldEFsbCA9IDA7XHJcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xyXG5jb25zdCBnZXRDdXJyZW50ID0gMjtcclxuXHJcbi8qKlxyXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKTogUHJvbWlzZTxTY3JlZW5bXT4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKlxyXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0UHJpbWFyeSgpOiBQcm9taXNlPFNjcmVlbj4ge1xyXG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqXHJcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCk6IFByb21pc2U8U2NyZWVuPiB7XHJcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcclxufVxyXG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5BbmRyb2lkKTtcclxuXHJcbi8vIE1ldGhvZCBJRHNcclxuY29uc3QgSGFwdGljc1ZpYnJhdGUgPSAwO1xyXG5jb25zdCBEZXZpY2VJbmZvID0gMTtcclxuY29uc3QgVG9hc3RTaG93ID0gMjtcclxuY29uc3QgU2Nyb2xsU2V0RW5hYmxlZCA9IDM7XHJcbmNvbnN0IFNjcm9sbFNldEJvdW5jZUVuYWJsZWQgPSA0O1xyXG5jb25zdCBTY3JvbGxTZXRJbmRpY2F0b3JzRW5hYmxlZCA9IDU7XHJcbmNvbnN0IE5hdmlnYXRpb25TZXRCYWNrRm9yd2FyZEdlc3R1cmVzID0gNjtcclxuY29uc3QgTGlua3NTZXRQcmV2aWV3RW5hYmxlZCA9IDc7XHJcbmNvbnN0IFVzZXJBZ2VudFNldCA9IDg7XHJcblxyXG5leHBvcnQgbmFtZXNwYWNlIEhhcHRpY3Mge1xyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFZpYnJhdGUoZHVyYXRpb246IG51bWJlciA9IDEwMCk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKEhhcHRpY3NWaWJyYXRlLCB7IGR1cmF0aW9uIH0pO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgbmFtZXNwYWNlIERldmljZSB7XHJcbiAgICBleHBvcnQgaW50ZXJmYWNlIEluZm8ge1xyXG4gICAgICAgIHBsYXRmb3JtOiBzdHJpbmc7XHJcbiAgICAgICAgbW9kZWw6IHN0cmluZztcclxuICAgICAgICB2ZXJzaW9uOiBzdHJpbmc7XHJcbiAgICB9XHJcblxyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIEluZm8oKTogUHJvbWlzZTxJbmZvPiB7XHJcbiAgICAgICAgcmV0dXJuIGNhbGwoRGV2aWNlSW5mbyk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgVG9hc3Qge1xyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFNob3cobWVzc2FnZTogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIGNhbGwoVG9hc3RTaG93LCB7IG1lc3NhZ2UgfSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgU2Nyb2xsIHtcclxuICAgIGV4cG9ydCBmdW5jdGlvbiBTZXRFbmFibGVkKGVuYWJsZWQ6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gY2FsbChTY3JvbGxTZXRFbmFibGVkLCB7IGVuYWJsZWQgfSk7XHJcbiAgICB9XHJcblxyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFNldEJvdW5jZUVuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKFNjcm9sbFNldEJvdW5jZUVuYWJsZWQsIHsgZW5hYmxlZCB9KTtcclxuICAgIH1cclxuXHJcbiAgICBleHBvcnQgZnVuY3Rpb24gU2V0SW5kaWNhdG9yc0VuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKFNjcm9sbFNldEluZGljYXRvcnNFbmFibGVkLCB7IGVuYWJsZWQgfSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgTmF2aWdhdGlvbiB7XHJcbiAgICBleHBvcnQgZnVuY3Rpb24gU2V0QmFja0ZvcndhcmRHZXN0dXJlc0VuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKE5hdmlnYXRpb25TZXRCYWNrRm9yd2FyZEdlc3R1cmVzLCB7IGVuYWJsZWQgfSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgTGlua3Mge1xyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIFNldFByZXZpZXdFbmFibGVkKGVuYWJsZWQ6IGJvb2xlYW4pOiBQcm9taXNlPHZvaWQ+IHtcclxuICAgICAgICByZXR1cm4gY2FsbChMaW5rc1NldFByZXZpZXdFbmFibGVkLCB7IGVuYWJsZWQgfSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBuYW1lc3BhY2UgVXNlckFnZW50IHtcclxuICAgIGV4cG9ydCBmdW5jdGlvbiBTZXQodWE6IHN0cmluZyk6IFByb21pc2U8dm9pZD4ge1xyXG4gICAgICAgIHJldHVybiBjYWxsKFVzZXJBZ2VudFNldCwgeyB1YSB9KTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF8gICAgIF9fICAgICBfIF9fXHJcbnwgfCAgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5JT1MpO1xyXG5cclxuLy8gTWV0aG9kIElEc1xyXG5jb25zdCBIYXB0aWNzSW1wYWN0ID0gMDtcclxuY29uc3QgRGV2aWNlSW5mbyA9IDE7XHJcblxyXG5leHBvcnQgbmFtZXNwYWNlIEhhcHRpY3Mge1xyXG4gICAgZXhwb3J0IHR5cGUgSW1wYWN0U3R5bGUgPSBcImxpZ2h0XCJ8XCJtZWRpdW1cInxcImhlYXZ5XCJ8XCJzb2Z0XCJ8XCJyaWdpZFwiO1xyXG4gICAgZXhwb3J0IGZ1bmN0aW9uIEltcGFjdChzdHlsZTogSW1wYWN0U3R5bGUgPSBcIm1lZGl1bVwiKTogUHJvbWlzZTx2b2lkPiB7XHJcbiAgICAgICAgcmV0dXJuIGNhbGwoSGFwdGljc0ltcGFjdCwgeyBzdHlsZSB9KTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IG5hbWVzcGFjZSBEZXZpY2Uge1xyXG4gICAgZXhwb3J0IGludGVyZmFjZSBJbmZvIHtcclxuICAgICAgICBtb2RlbDogc3RyaW5nO1xyXG4gICAgICAgIHN5c3RlbU5hbWU6IHN0cmluZztcclxuICAgICAgICBzeXN0ZW1WZXJzaW9uOiBzdHJpbmc7XHJcbiAgICAgICAgaXNTaW11bGF0b3I6IGJvb2xlYW47XHJcbiAgICB9XHJcbiAgICBleHBvcnQgZnVuY3Rpb24gSW5mbygpOiBQcm9taXNlPEluZm8+IHtcclxuICAgICAgICByZXR1cm4gY2FsbChEZXZpY2VJbmZvKTtcclxuICAgIH1cclxufVxyXG4iXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ0FBO0FBQUE7QUFBQTtBQUFBOzs7QUM2QkEsSUFBTSxjQUNGO0FBRUcsU0FBUyxPQUFPLE9BQWUsSUFBWTtBQUM5QyxNQUFJLEtBQUs7QUFFVCxNQUFJLElBQUksT0FBTztBQUNmLFNBQU8sS0FBSztBQUVSLFVBQU0sWUFBYSxLQUFL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O0FBT0ksV0FBU0MsUUFBc0I7QUFDbEMsV0FBT0osT0FBS0MsV0FBVTtBQUFBLEVBQzFCO0FBRk8sRUFBQUUsUUFBUyxPQUFBQztBQUFBLEdBUEhELHNCQUFBOzs7QXhCZGpCLE9BQU8sU0FBUyxPQUFPLFVBQVUsQ0FBQztBQTBEbEMsT0FBTyxPQUFPLFNBQWdCO0FBQzlCLE9BQU8sT0FBTyxXQUFXO0FBS3pCLE9BQU8sT0FBTyx5QkFBeUIsZUFBTyx1QkFBdUIsS0FBSyxjQUFNO0FBR2hGLE9BQU8sT0FBTyxrQkFBa0I7QUFDaEMsT0FBTyxPQUFPLGtCQUFrQjtBQUNoQyxPQUFPLE9BQU8saUJBQWlCO0FBRXhCLE9BQU8scUJBQXFCO0FBTzVCLFNBQVMsbUJBQW1CLEtBQTRCO0FBQzNELFNBQU8sTUFBTSxLQUFLLEVBQUUsUUFBUSxPQUFPLENBQUMsRUFDL0IsS0FBSyxjQUFZO0FBQ2QsUUFBSSxTQUFTLElBQUk7QUFDYixZQUFNLFNBQVMsU0FBUyxjQUFjLFFBQVE7QUFDOUMsYUFBTyxNQUFNO0FBQ2IsZUFBUyxLQUFLLFlBQVksTUFBTTtBQUFBLElBQ3BDO0FBQUEsRUFDSixDQUFDLEVBQ0EsTUFBTSxNQUFNO0FBQUEsRUFBQyxDQUFDO0FBQ3ZCO0FBR0EsbUJBQW1CLGtCQUFrQjsiLAogICJuYW1lcyI6IFsiX2EiLCAiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJfYSIsICJBcnJheSIsICJNYXAiLCAiQXJyYXkiLCAiTWFwIiwgImtleSIsICJjYWxsIiwgIl9hIiwgIl9hIiwgInJlc2l6YWJsZSIsICJfYSIsICJjYWxsIiwgIl9hIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJIaWRlTWV0aG9kIiwgIlNob3dNZXRob2QiLCAiaXNEb2N1bWVudERvdEFsbCIsICJfYSIsICJyZWFzb24iLCAidmFsdWUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIiwgIlRvYXN0IiwgIlNob3ciLCAiU2Nyb2xsIiwgIk5hdmlnYXRpb24iLCAiTGlua3MiLCAiVXNlckFnZW50IiwgIkRldmljZSIsICJIYXB0aWNzIiwgImNhbGwiLCAiRGV2aWNlSW5mbyIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIl0KfQo=
