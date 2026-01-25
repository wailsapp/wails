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
fetch("/wails/custom.js", { method: "HEAD" }).then((response) => {
  if (response.ok) {
    const script = document.createElement("script");
    script.src = "/wails/custom.js";
    document.head.appendChild(script);
  }
}).catch(() => {
});
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
  objectNames,
  setTransport
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2luZGV4LnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy93bWwudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2Jyb3dzZXIudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL25hbm9pZC50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvcnVudGltZS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZGlhbG9ncy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZXZlbnRzLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9saXN0ZW5lci50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY3JlYXRlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9ldmVudF90eXBlcy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvdXRpbHMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3dpbmRvdy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvY29tcGlsZWQvbWFpbi5qcyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvc3lzdGVtLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jb250ZXh0bWVudS50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvZmxhZ3MudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2RyYWcudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2FwcGxpY2F0aW9uLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jYWxscy50cyIsICIuLi8uLi9ydW50aW1lL2Rlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zcmMvY2FsbGFibGUudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2NhbmNlbGxhYmxlLnRzIiwgIi4uLy4uL3J1bnRpbWUvZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3NyYy9jbGlwYm9hcmQudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL3NjcmVlbnMudHMiLCAiLi4vLi4vcnVudGltZS9kZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvc3JjL2lvcy50cyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxuaW1wb3J0IFwiLi9jb250ZXh0bWVudS5qc1wiO1xuaW1wb3J0IFwiLi9kcmFnLmpzXCI7XG5cbi8vIFJlLWV4cG9ydCBwdWJsaWMgQVBJXG5pbXBvcnQgKiBhcyBBcHBsaWNhdGlvbiBmcm9tIFwiLi9hcHBsaWNhdGlvbi5qc1wiO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tIFwiLi9icm93c2VyLmpzXCI7XG5pbXBvcnQgKiBhcyBDYWxsIGZyb20gXCIuL2NhbGxzLmpzXCI7XG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSBcIi4vY2xpcGJvYXJkLmpzXCI7XG5pbXBvcnQgKiBhcyBDcmVhdGUgZnJvbSBcIi4vY3JlYXRlLmpzXCI7XG5pbXBvcnQgKiBhcyBEaWFsb2dzIGZyb20gXCIuL2RpYWxvZ3MuanNcIjtcbmltcG9ydCAqIGFzIEV2ZW50cyBmcm9tIFwiLi9ldmVudHMuanNcIjtcbmltcG9ydCAqIGFzIEZsYWdzIGZyb20gXCIuL2ZsYWdzLmpzXCI7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gXCIuL3NjcmVlbnMuanNcIjtcbmltcG9ydCAqIGFzIFN5c3RlbSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCAqIGFzIElPUyBmcm9tIFwiLi9pb3MuanNcIjtcbmltcG9ydCBXaW5kb3csIHsgaGFuZGxlRHJhZ0VudGVyLCBoYW5kbGVEcmFnTGVhdmUsIGhhbmRsZURyYWdPdmVyIH0gZnJvbSBcIi4vd2luZG93LmpzXCI7XG5pbXBvcnQgKiBhcyBXTUwgZnJvbSBcIi4vd21sLmpzXCI7XG5cbmV4cG9ydCB7XG4gICAgQXBwbGljYXRpb24sXG4gICAgQnJvd3NlcixcbiAgICBDYWxsLFxuICAgIENsaXBib2FyZCxcbiAgICBEaWFsb2dzLFxuICAgIEV2ZW50cyxcbiAgICBGbGFncyxcbiAgICBTY3JlZW5zLFxuICAgIFN5c3RlbSxcbiAgICBJT1MsXG4gICAgV2luZG93LFxuICAgIFdNTFxufTtcblxuLyoqXG4gKiBBbiBpbnRlcm5hbCB1dGlsaXR5IGNvbnN1bWVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cbiAqXG4gKiBAaWdub3JlXG4gKi9cbmV4cG9ydCB7IENyZWF0ZSB9O1xuXG5leHBvcnQgKiBmcm9tIFwiLi9jYW5jZWxsYWJsZS5qc1wiO1xuXG4vLyBFeHBvcnQgdHJhbnNwb3J0IGludGVyZmFjZXMgYW5kIHV0aWxpdGllc1xuZXhwb3J0IHtcbiAgICBzZXRUcmFuc3BvcnQsXG4gICAgZ2V0VHJhbnNwb3J0LFxuICAgIHR5cGUgUnVudGltZVRyYW5zcG9ydCxcbiAgICBvYmplY3ROYW1lcyxcbiAgICBjbGllbnRJZCxcbn0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5pbXBvcnQgeyBjbGllbnRJZCB9IGZyb20gXCIuL3J1bnRpbWUuanNcIjtcblxuLy8gTm90aWZ5IGJhY2tlbmRcbndpbmRvdy5fd2FpbHMuaW52b2tlID0gU3lzdGVtLmludm9rZTtcbndpbmRvdy5fd2FpbHMuY2xpZW50SWQgPSBjbGllbnRJZDtcblxuLy8gUmVnaXN0ZXIgcGxhdGZvcm0gaGFuZGxlcnMgKGludGVybmFsIEFQSSlcbi8vIE5vdGU6IFdpbmRvdyBpcyB0aGUgdGhpc1dpbmRvdyBpbnN0YW5jZSAoZGVmYXVsdCBleHBvcnQgZnJvbSB3aW5kb3cudHMpXG4vLyBCaW5kaW5nIGVuc3VyZXMgJ3RoaXMnIGNvcnJlY3RseSByZWZlcnMgdG8gdGhlIGN1cnJlbnQgd2luZG93IGluc3RhbmNlXG53aW5kb3cuX3dhaWxzLmhhbmRsZVBsYXRmb3JtRmlsZURyb3AgPSBXaW5kb3cuSGFuZGxlUGxhdGZvcm1GaWxlRHJvcC5iaW5kKFdpbmRvdyk7XG5cbi8vIExpbnV4LXNwZWNpZmljIGRyYWcgaGFuZGxlcnMgKEdUSyBpbnRlcmNlcHRzIERPTSBkcmFnIGV2ZW50cylcbndpbmRvdy5fd2FpbHMuaGFuZGxlRHJhZ0VudGVyID0gaGFuZGxlRHJhZ0VudGVyO1xud2luZG93Ll93YWlscy5oYW5kbGVEcmFnTGVhdmUgPSBoYW5kbGVEcmFnTGVhdmU7XG53aW5kb3cuX3dhaWxzLmhhbmRsZURyYWdPdmVyID0gaGFuZGxlRHJhZ092ZXI7XG5cblN5c3RlbS5pbnZva2UoXCJ3YWlsczpydW50aW1lOnJlYWR5XCIpO1xuXG4vLyBMb2FkIGN1c3RvbS5qcyBpZiBhdmFpbGFibGUgKHVzZWQgYnkgc2VydmVyIG1vZGUgZm9yIFdlYlNvY2tldCBldmVudHMsIGV0Yy4pXG4vLyBVc2VzIHNjcmlwdCB0YWcgaW5qZWN0aW9uIGluc3RlYWQgb2YgZXZhbCgpIGZvciBiZXR0ZXIgc2VjdXJpdHkgYW5kIENTUCBjb21wYXRpYmlsaXR5XG5mZXRjaCgnL3dhaWxzL2N1c3RvbS5qcycsIHsgbWV0aG9kOiAnSEVBRCcgfSlcbiAgICAudGhlbihyZXNwb25zZSA9PiB7XG4gICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgY29uc3Qgc2NyaXB0ID0gZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgnc2NyaXB0Jyk7XG4gICAgICAgICAgICBzY3JpcHQuc3JjID0gJy93YWlscy9jdXN0b20uanMnO1xuICAgICAgICAgICAgZG9jdW1lbnQuaGVhZC5hcHBlbmRDaGlsZChzY3JpcHQpO1xuICAgICAgICB9XG4gICAgfSlcbiAgICAuY2F0Y2goKCkgPT4ge30pOyAvLyBTaWxlbnRseSBpZ25vcmUgLSBjdXN0b20uanMgaXMgb3B0aW9uYWxcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgT3BlblVSTCB9IGZyb20gXCIuL2Jyb3dzZXIuanNcIjtcbmltcG9ydCB7IFF1ZXN0aW9uIH0gZnJvbSBcIi4vZGlhbG9ncy5qc1wiO1xuaW1wb3J0IHsgRW1pdCB9IGZyb20gXCIuL2V2ZW50cy5qc1wiO1xuaW1wb3J0IHsgY2FuQWJvcnRMaXN0ZW5lcnMsIHdoZW5SZWFkeSB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5pbXBvcnQgV2luZG93IGZyb20gXCIuL3dpbmRvdy5qc1wiO1xuXG4vKipcbiAqIFNlbmRzIGFuIGV2ZW50IHdpdGggdGhlIGdpdmVuIG5hbWUgYW5kIG9wdGlvbmFsIGRhdGEuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXG4gKiBAcGFyYW0gW2RhdGE9bnVsbF0gLSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cbiAqL1xuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZTogc3RyaW5nLCBkYXRhOiBhbnkgPSBudWxsKTogdm9pZCB7XG4gICAgRW1pdChldmVudE5hbWUsIGRhdGEpO1xufVxuXG4vKipcbiAqIENhbGxzIGEgbWV0aG9kIG9uIGEgc3BlY2lmaWVkIHdpbmRvdy5cbiAqXG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gY2FsbCB0aGUgbWV0aG9kIG9uLlxuICogQHBhcmFtIG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXG4gKi9cbmZ1bmN0aW9uIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TmFtZTogc3RyaW5nLCBtZXRob2ROYW1lOiBzdHJpbmcpIHtcbiAgICBjb25zdCB0YXJnZXRXaW5kb3cgPSBXaW5kb3cuR2V0KHdpbmRvd05hbWUpO1xuICAgIGNvbnN0IG1ldGhvZCA9ICh0YXJnZXRXaW5kb3cgYXMgYW55KVttZXRob2ROYW1lXTtcblxuICAgIGlmICh0eXBlb2YgbWV0aG9kICE9PSBcImZ1bmN0aW9uXCIpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgV2luZG93IG1ldGhvZCAnJHttZXRob2ROYW1lfScgbm90IGZvdW5kYCk7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICB0cnkge1xuICAgICAgICBtZXRob2QuY2FsbCh0YXJnZXRXaW5kb3cpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgY29uc29sZS5lcnJvcihgRXJyb3IgY2FsbGluZyB3aW5kb3cgbWV0aG9kICcke21ldGhvZE5hbWV9JzogYCwgZSk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJlc3BvbmRzIHRvIGEgdHJpZ2dlcmluZyBldmVudCBieSBydW5uaW5nIGFwcHJvcHJpYXRlIFdNTCBhY3Rpb25zIGZvciB0aGUgY3VycmVudCB0YXJnZXQuXG4gKi9cbmZ1bmN0aW9uIG9uV01MVHJpZ2dlcmVkKGV2OiBFdmVudCk6IHZvaWQge1xuICAgIGNvbnN0IGVsZW1lbnQgPSBldi5jdXJyZW50VGFyZ2V0IGFzIEVsZW1lbnQ7XG5cbiAgICBmdW5jdGlvbiBydW5FZmZlY3QoY2hvaWNlID0gXCJZZXNcIikge1xuICAgICAgICBpZiAoY2hvaWNlICE9PSBcIlllc1wiKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgdGFyZ2V0V2luZG93ID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10YXJnZXQtd2luZG93JykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRhcmdldC13aW5kb3cnKSB8fCBcIlwiO1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpIHx8IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC13aW5kb3cnKTtcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLW9wZW51cmwnKTtcblxuICAgICAgICBpZiAoZXZlbnRUeXBlICE9PSBudWxsKVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIGlmICh3aW5kb3dNZXRob2QgIT09IG51bGwpXG4gICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHRhcmdldFdpbmRvdywgd2luZG93TWV0aG9kKTtcbiAgICAgICAgaWYgKHVybCAhPT0gbnVsbClcbiAgICAgICAgICAgIHZvaWQgT3BlblVSTCh1cmwpO1xuICAgIH1cblxuICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKSB8fCBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuXG4gICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgUXVlc3Rpb24oe1xuICAgICAgICAgICAgVGl0bGU6IFwiQ29uZmlybVwiLFxuICAgICAgICAgICAgTWVzc2FnZTogY29uZmlybSxcbiAgICAgICAgICAgIERldGFjaGVkOiBmYWxzZSxcbiAgICAgICAgICAgIEJ1dHRvbnM6IFtcbiAgICAgICAgICAgICAgICB7IExhYmVsOiBcIlllc1wiIH0sXG4gICAgICAgICAgICAgICAgeyBMYWJlbDogXCJOb1wiLCBJc0RlZmF1bHQ6IHRydWUgfVxuICAgICAgICAgICAgXVxuICAgICAgICB9KS50aGVuKHJ1bkVmZmVjdCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgcnVuRWZmZWN0KCk7XG4gICAgfVxufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY29udHJvbGxlclN5bSA9IFN5bWJvbChcImNvbnRyb2xsZXJcIik7XG5jb25zdCB0cmlnZ2VyTWFwU3ltID0gU3ltYm9sKFwidHJpZ2dlck1hcFwiKTtcbmNvbnN0IGVsZW1lbnRDb3VudFN5bSA9IFN5bWJvbChcImVsZW1lbnRDb3VudFwiKTtcblxuLyoqXG4gKiBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSBkb2VzIG5vdCBhY3R1YWxseSByZW1lbWJlciBhY3RpdmUgZXZlbnQgbGlzdGVuZXJzOiBpbnN0ZWFkXG4gKiBpdCB0aWVzIHRoZW0gdG8gYW4gQWJvcnRTaWduYWwgYW5kIHVzZXMgYW4gQWJvcnRDb250cm9sbGVyIHRvIHJlbW92ZSB0aGVtIGFsbCBhdCBvbmNlLlxuICovXG5jbGFzcyBBYm9ydENvbnRyb2xsZXJSZWdpc3RyeSB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgW2NvbnRyb2xsZXJTeW1dOiBBYm9ydENvbnRyb2xsZXI7XG5cbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIGFuIG9wdGlvbnMgb2JqZWN0IGZvciBhZGRFdmVudExpc3RlbmVyIHRoYXQgdGllcyB0aGUgbGlzdGVuZXJcbiAgICAgKiB0byB0aGUgQWJvcnRTaWduYWwgZnJvbSB0aGUgY3VycmVudCBBYm9ydENvbnRyb2xsZXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZWxlbWVudCAtIEFuIEhUTUwgZWxlbWVudFxuICAgICAqIEBwYXJhbSB0cmlnZ2VycyAtIFRoZSBsaXN0IG9mIGFjdGl2ZSBXTUwgdHJpZ2dlciBldmVudHMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudHNcbiAgICAgKi9cbiAgICBzZXQoZWxlbWVudDogRWxlbWVudCwgdHJpZ2dlcnM6IHN0cmluZ1tdKTogQWRkRXZlbnRMaXN0ZW5lck9wdGlvbnMge1xuICAgICAgICByZXR1cm4geyBzaWduYWw6IHRoaXNbY29udHJvbGxlclN5bV0uc2lnbmFsIH07XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyBhbGwgcmVnaXN0ZXJlZCBldmVudCBsaXN0ZW5lcnMgYW5kIHJlc2V0cyB0aGUgcmVnaXN0cnkuXG4gICAgICovXG4gICAgcmVzZXQoKTogdm9pZCB7XG4gICAgICAgIHRoaXNbY29udHJvbGxlclN5bV0uYWJvcnQoKTtcbiAgICAgICAgdGhpc1tjb250cm9sbGVyU3ltXSA9IG5ldyBBYm9ydENvbnRyb2xsZXIoKTtcbiAgICB9XG59XG5cbi8qKlxuICogV2Vha01hcFJlZ2lzdHJ5IG1hcHMgYWN0aXZlIHRyaWdnZXIgZXZlbnRzIHRvIGVhY2ggRE9NIGVsZW1lbnQgdGhyb3VnaCBhIFdlYWtNYXAuXG4gKiBUaGlzIGVuc3VyZXMgdGhhdCB0aGUgbWFwcGluZyByZW1haW5zIHByaXZhdGUgdG8gdGhpcyBtb2R1bGUsIHdoaWxlIHN0aWxsIGFsbG93aW5nIGdhcmJhZ2VcbiAqIGNvbGxlY3Rpb24gb2YgdGhlIGludm9sdmVkIGVsZW1lbnRzLlxuICovXG5jbGFzcyBXZWFrTWFwUmVnaXN0cnkge1xuICAgIC8qKiBTdG9yZXMgdGhlIGN1cnJlbnQgZWxlbWVudC10by10cmlnZ2VyIG1hcHBpbmcuICovXG4gICAgW3RyaWdnZXJNYXBTeW1dOiBXZWFrTWFwPEVsZW1lbnQsIHN0cmluZ1tdPjtcbiAgICAvKiogQ291bnRzIHRoZSBudW1iZXIgb2YgZWxlbWVudHMgd2l0aCBhY3RpdmUgV01MIHRyaWdnZXJzLiAqL1xuICAgIFtlbGVtZW50Q291bnRTeW1dOiBudW1iZXI7XG5cbiAgICBjb25zdHJ1Y3RvcigpIHtcbiAgICAgICAgdGhpc1t0cmlnZ2VyTWFwU3ltXSA9IG5ldyBXZWFrTWFwKCk7XG4gICAgICAgIHRoaXNbZWxlbWVudENvdW50U3ltXSA9IDA7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyBhY3RpdmUgdHJpZ2dlcnMgZm9yIHRoZSBzcGVjaWZpZWQgZWxlbWVudC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbGVtZW50IC0gQW4gSFRNTCBlbGVtZW50XG4gICAgICogQHBhcmFtIHRyaWdnZXJzIC0gVGhlIGxpc3Qgb2YgYWN0aXZlIFdNTCB0cmlnZ2VyIGV2ZW50cyBmb3IgdGhlIHNwZWNpZmllZCBlbGVtZW50XG4gICAgICovXG4gICAgc2V0KGVsZW1lbnQ6IEVsZW1lbnQsIHRyaWdnZXJzOiBzdHJpbmdbXSk6IEFkZEV2ZW50TGlzdGVuZXJPcHRpb25zIHtcbiAgICAgICAgaWYgKCF0aGlzW3RyaWdnZXJNYXBTeW1dLmhhcyhlbGVtZW50KSkgeyB0aGlzW2VsZW1lbnRDb3VudFN5bV0rKzsgfVxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dLnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgICAgIHJldHVybiB7fTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZW1vdmVzIGFsbCByZWdpc3RlcmVkIGV2ZW50IGxpc3RlbmVycy5cbiAgICAgKi9cbiAgICByZXNldCgpOiB2b2lkIHtcbiAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxuICAgICAgICAgICAgcmV0dXJuO1xuXG4gICAgICAgIGZvciAoY29uc3QgZWxlbWVudCBvZiBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJyonKSkge1xuICAgICAgICAgICAgaWYgKHRoaXNbZWxlbWVudENvdW50U3ltXSA8PSAwKVxuICAgICAgICAgICAgICAgIGJyZWFrO1xuXG4gICAgICAgICAgICBjb25zdCB0cmlnZ2VycyA9IHRoaXNbdHJpZ2dlck1hcFN5bV0uZ2V0KGVsZW1lbnQpO1xuICAgICAgICAgICAgaWYgKHRyaWdnZXJzICE9IG51bGwpIHsgdGhpc1tlbGVtZW50Q291bnRTeW1dLS07IH1cblxuICAgICAgICAgICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzIHx8IFtdKVxuICAgICAgICAgICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBvbldNTFRyaWdnZXJlZCk7XG4gICAgICAgIH1cblxuICAgICAgICB0aGlzW3RyaWdnZXJNYXBTeW1dID0gbmV3IFdlYWtNYXAoKTtcbiAgICAgICAgdGhpc1tlbGVtZW50Q291bnRTeW1dID0gMDtcbiAgICB9XG59XG5cbmNvbnN0IHRyaWdnZXJSZWdpc3RyeSA9IGNhbkFib3J0TGlzdGVuZXJzKCkgPyBuZXcgQWJvcnRDb250cm9sbGVyUmVnaXN0cnkoKSA6IG5ldyBXZWFrTWFwUmVnaXN0cnkoKTtcblxuLyoqXG4gKiBBZGRzIGV2ZW50IGxpc3RlbmVycyB0byB0aGUgc3BlY2lmaWVkIGVsZW1lbnQuXG4gKi9cbmZ1bmN0aW9uIGFkZFdNTExpc3RlbmVycyhlbGVtZW50OiBFbGVtZW50KTogdm9pZCB7XG4gICAgY29uc3QgdHJpZ2dlclJlZ0V4cCA9IC9cXFMrL2c7XG4gICAgY29uc3QgdHJpZ2dlckF0dHIgPSAoZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCIpO1xuICAgIGNvbnN0IHRyaWdnZXJzOiBzdHJpbmdbXSA9IFtdO1xuXG4gICAgbGV0IG1hdGNoO1xuICAgIHdoaWxlICgobWF0Y2ggPSB0cmlnZ2VyUmVnRXhwLmV4ZWModHJpZ2dlckF0dHIpKSAhPT0gbnVsbClcbiAgICAgICAgdHJpZ2dlcnMucHVzaChtYXRjaFswXSk7XG5cbiAgICBjb25zdCBvcHRpb25zID0gdHJpZ2dlclJlZ2lzdHJ5LnNldChlbGVtZW50LCB0cmlnZ2Vycyk7XG4gICAgZm9yIChjb25zdCB0cmlnZ2VyIG9mIHRyaWdnZXJzKVxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgb25XTUxUcmlnZ2VyZWQsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNjaGVkdWxlcyBhbiBhdXRvbWF0aWMgcmVsb2FkIG9mIFdNTCB0byBiZSBwZXJmb3JtZWQgYXMgc29vbiBhcyB0aGUgZG9jdW1lbnQgaXMgZnVsbHkgbG9hZGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gRW5hYmxlKCk6IHZvaWQge1xuICAgIHdoZW5SZWFkeShSZWxvYWQpO1xufVxuXG4vKipcbiAqIFJlbG9hZHMgdGhlIFdNTCBwYWdlIGJ5IGFkZGluZyBuZWNlc3NhcnkgZXZlbnQgbGlzdGVuZXJzIGFuZCBicm93c2VyIGxpc3RlbmVycy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFJlbG9hZCgpOiB2b2lkIHtcbiAgICB0cmlnZ2VyUmVnaXN0cnkucmVzZXQoKTtcbiAgICBkb2N1bWVudC5ib2R5LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtZXZlbnRdLCBbd21sLXdpbmRvd10sIFt3bWwtb3BlbnVybF0sIFtkYXRhLXdtbC1ldmVudF0sIFtkYXRhLXdtbC13aW5kb3ddLCBbZGF0YS13bWwtb3BlbnVybF0nKS5mb3JFYWNoKGFkZFdNTExpc3RlbmVycyk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5Ccm93c2VyKTtcblxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xuXG4vKipcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMLlxuICpcbiAqIEBwYXJhbSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5VUkwodXJsOiBzdHJpbmcgfCBVUkwpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICByZXR1cm4gY2FsbChCcm93c2VyT3BlblVSTCwge3VybDogdXJsLnRvU3RyaW5nKCl9KTtcbn1cbiIsICIvLyBTb3VyY2U6IGh0dHBzOi8vZ2l0aHViLmNvbS9haS9uYW5vaWRcblxuLy8gVGhlIE1JVCBMaWNlbnNlIChNSVQpXG4vL1xuLy8gQ29weXJpZ2h0IDIwMTcgQW5kcmV5IFNpdG5payA8YW5kcmV5QHNpdG5pay5ydT5cbi8vXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5IG9mXG4vLyB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsIGluXG4vLyB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzIHRvXG4vLyB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsIGNvcGllcyBvZlxuLy8gdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMgZnVybmlzaGVkIHRvIGRvIHNvLFxuLy8gICAgIHN1YmplY3QgdG8gdGhlIGZvbGxvd2luZyBjb25kaXRpb25zOlxuLy9cbi8vICAgICBUaGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSBhbmQgdGhpcyBwZXJtaXNzaW9uIG5vdGljZSBzaGFsbCBiZSBpbmNsdWRlZCBpbiBhbGxcbi8vIGNvcGllcyBvciBzdWJzdGFudGlhbCBwb3J0aW9ucyBvZiB0aGUgU29mdHdhcmUuXG4vL1xuLy8gICAgIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1Jcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLCBGSVRORVNTXG4vLyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFIEFVVEhPUlMgT1Jcbi8vIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSIExJQUJJTElUWSwgV0hFVEhFUlxuLy8gSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLCBPVVQgT0YgT1IgSU5cbi8vIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEUgU09GVFdBUkUuXG5cbi8vIFRoaXMgYWxwaGFiZXQgdXNlcyBgQS1aYS16MC05Xy1gIHN5bWJvbHMuXG4vLyBUaGUgb3JkZXIgb2YgY2hhcmFjdGVycyBpcyBvcHRpbWl6ZWQgZm9yIGJldHRlciBnemlwIGFuZCBicm90bGkgY29tcHJlc3Npb24uXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBzYW1lIGZpbGUgKHdvcmtzIGJvdGggZm9yIGd6aXAgYW5kIGJyb3RsaSk6XG4vLyBgJ3VzZWAsIGBhbmRvbWAsIGFuZCBgcmljdCdgXG4vLyBSZWZlcmVuY2VzIHRvIHRoZSBicm90bGkgZGVmYXVsdCBkaWN0aW9uYXJ5OlxuLy8gYC0yNlRgLCBgMTk4M2AsIGA0MHB4YCwgYDc1cHhgLCBgYnVzaGAsIGBqYWNrYCwgYG1pbmRgLCBgdmVyeWAsIGFuZCBgd29sZmBcbmNvbnN0IHVybEFscGhhYmV0ID1cbiAgICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcblxuZXhwb3J0IGZ1bmN0aW9uIG5hbm9pZChzaXplOiBudW1iZXIgPSAyMSk6IHN0cmluZyB7XG4gICAgbGV0IGlkID0gJydcbiAgICAvLyBBIGNvbXBhY3QgYWx0ZXJuYXRpdmUgZm9yIGBmb3IgKHZhciBpID0gMDsgaSA8IHN0ZXA7IGkrKylgLlxuICAgIGxldCBpID0gc2l6ZSB8IDBcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICAgIC8vIGB8IDBgIGlzIG1vcmUgY29tcGFjdCBhbmQgZmFzdGVyIHRoYW4gYE1hdGguZmxvb3IoKWAuXG4gICAgICAgIGlkICs9IHVybEFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogNjQpIHwgMF1cbiAgICB9XG4gICAgcmV0dXJuIGlkXG59XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gXCIuL25hbm9pZC5qc1wiO1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuLy8gUmUtZXhwb3J0IG5hbm9pZCBmb3IgY3VzdG9tIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnNcbmV4cG9ydCB7IG5hbm9pZCB9O1xuXG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IE9iamVjdC5mcmVlemUoe1xuICAgIENhbGw6IDAsXG4gICAgQ2xpcGJvYXJkOiAxLFxuICAgIEFwcGxpY2F0aW9uOiAyLFxuICAgIEV2ZW50czogMyxcbiAgICBDb250ZXh0TWVudTogNCxcbiAgICBEaWFsb2c6IDUsXG4gICAgV2luZG93OiA2LFxuICAgIFNjcmVlbnM6IDcsXG4gICAgU3lzdGVtOiA4LFxuICAgIEJyb3dzZXI6IDksXG4gICAgQ2FuY2VsQ2FsbDogMTAsXG4gICAgSU9TOiAxMSxcbn0pO1xuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG4vKipcbiAqIFJ1bnRpbWVUcmFuc3BvcnQgZGVmaW5lcyB0aGUgaW50ZXJmYWNlIGZvciBjdXN0b20gSVBDIHRyYW5zcG9ydCBpbXBsZW1lbnRhdGlvbnMuXG4gKiBJbXBsZW1lbnQgdGhpcyBpbnRlcmZhY2UgdG8gdXNlIFdlYlNvY2tldHMsIGN1c3RvbSBwcm90b2NvbHMsIG9yIGFueSBvdGhlclxuICogdHJhbnNwb3J0IG1lY2hhbmlzbSBpbnN0ZWFkIG9mIHRoZSBkZWZhdWx0IEhUVFAgZmV0Y2guXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgUnVudGltZVRyYW5zcG9ydCB7XG4gICAgLyoqXG4gICAgICogU2VuZCBhIHJ1bnRpbWUgY2FsbCBhbmQgcmV0dXJuIHRoZSByZXNwb25zZS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBvYmplY3RJRCAtIFRoZSBXYWlscyBvYmplY3QgSUQgKDA9Q2FsbCwgMT1DbGlwYm9hcmQsIGV0Yy4pXG4gICAgICogQHBhcmFtIG1ldGhvZCAtIFRoZSBtZXRob2QgSUQgdG8gY2FsbFxuICAgICAqIEBwYXJhbSB3aW5kb3dOYW1lIC0gT3B0aW9uYWwgd2luZG93IG5hbWVcbiAgICAgKiBAcGFyYW0gYXJncyAtIEFyZ3VtZW50cyB0byBwYXNzICh3aWxsIGJlIEpTT04gc3RyaW5naWZpZWQgaWYgcHJlc2VudClcbiAgICAgKiBAcmV0dXJucyBQcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgcmVzcG9uc2UgZGF0YVxuICAgICAqL1xuICAgIGNhbGwob2JqZWN0SUQ6IG51bWJlciwgbWV0aG9kOiBudW1iZXIsIHdpbmRvd05hbWU6IHN0cmluZywgYXJnczogYW55KTogUHJvbWlzZTxhbnk+O1xufVxuXG4vKipcbiAqIEN1c3RvbSB0cmFuc3BvcnQgaW1wbGVtZW50YXRpb24gKGNhbiBiZSBzZXQgYnkgdXNlcilcbiAqL1xubGV0IGN1c3RvbVRyYW5zcG9ydDogUnVudGltZVRyYW5zcG9ydCB8IG51bGwgPSBudWxsO1xuXG4vKipcbiAqIFNldCBhIGN1c3RvbSB0cmFuc3BvcnQgZm9yIGFsbCBXYWlscyBydW50aW1lIGNhbGxzLlxuICogVGhpcyBhbGxvd3MgeW91IHRvIHJlcGxhY2UgdGhlIGRlZmF1bHQgSFRUUCBmZXRjaCB0cmFuc3BvcnQgd2l0aFxuICogV2ViU29ja2V0cywgY3VzdG9tIHByb3RvY29scywgb3IgYW55IG90aGVyIG1lY2hhbmlzbS5cbiAqXG4gKiBAcGFyYW0gdHJhbnNwb3J0IC0gWW91ciBjdXN0b20gdHJhbnNwb3J0IGltcGxlbWVudGF0aW9uXG4gKlxuICogQGV4YW1wbGVcbiAqIGBgYHR5cGVzY3JpcHRcbiAqIGltcG9ydCB7IHNldFRyYW5zcG9ydCB9IGZyb20gJy93YWlscy9ydW50aW1lLmpzJztcbiAqXG4gKiBjb25zdCB3c1RyYW5zcG9ydCA9IHtcbiAqICAgY2FsbDogYXN5bmMgKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpID0+IHtcbiAqICAgICAvLyBZb3VyIFdlYlNvY2tldCBpbXBsZW1lbnRhdGlvblxuICogICB9XG4gKiB9O1xuICpcbiAqIHNldFRyYW5zcG9ydCh3c1RyYW5zcG9ydCk7XG4gKiBgYGBcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIHNldFRyYW5zcG9ydCh0cmFuc3BvcnQ6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsKTogdm9pZCB7XG4gICAgY3VzdG9tVHJhbnNwb3J0ID0gdHJhbnNwb3J0O1xufVxuXG4vKipcbiAqIEdldCB0aGUgY3VycmVudCB0cmFuc3BvcnQgKHVzZWZ1bCBmb3IgZXh0ZW5kaW5nL3dyYXBwaW5nKVxuICovXG5leHBvcnQgZnVuY3Rpb24gZ2V0VHJhbnNwb3J0KCk6IFJ1bnRpbWVUcmFuc3BvcnQgfCBudWxsIHtcbiAgICByZXR1cm4gY3VzdG9tVHJhbnNwb3J0O1xufVxuXG4vKipcbiAqIENyZWF0ZXMgYSBuZXcgcnVudGltZSBjYWxsZXIgd2l0aCBzcGVjaWZpZWQgSUQuXG4gKlxuICogQHBhcmFtIG9iamVjdCAtIFRoZSBvYmplY3QgdG8gaW52b2tlIHRoZSBtZXRob2Qgb24uXG4gKiBAcGFyYW0gd2luZG93TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gKiBAcmV0dXJuIFRoZSBuZXcgcnVudGltZSBjYWxsZXIgZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdDogbnVtYmVyLCB3aW5kb3dOYW1lOiBzdHJpbmcgPSAnJykge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kOiBudW1iZXIsIGFyZ3M6IGFueSA9IG51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG5hc3luYyBmdW5jdGlvbiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3RJRDogbnVtYmVyLCBtZXRob2Q6IG51bWJlciwgd2luZG93TmFtZTogc3RyaW5nLCBhcmdzOiBhbnkpOiBQcm9taXNlPGFueT4ge1xuICAgIC8vIFVzZSBjdXN0b20gdHJhbnNwb3J0IGlmIGF2YWlsYWJsZVxuICAgIGlmIChjdXN0b21UcmFuc3BvcnQpIHtcbiAgICAgICAgcmV0dXJuIGN1c3RvbVRyYW5zcG9ydC5jYWxsKG9iamVjdElELCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xuICAgIH1cblxuICAgIC8vIERlZmF1bHQgSFRUUCBmZXRjaCB0cmFuc3BvcnRcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcblxuICAgIGxldCBib2R5OiB7IG9iamVjdDogbnVtYmVyOyBtZXRob2Q6IG51bWJlciwgYXJncz86IGFueSB9ID0ge1xuICAgICAgb2JqZWN0OiBvYmplY3RJRCxcbiAgICAgIG1ldGhvZFxuICAgIH1cbiAgICBpZiAoYXJncyAhPT0gbnVsbCAmJiBhcmdzICE9PSB1bmRlZmluZWQpIHtcbiAgICAgIGJvZHkuYXJncyA9IGFyZ3M7XG4gICAgfVxuXG4gICAgbGV0IGhlYWRlcnM6IFJlY29yZDxzdHJpbmcsIHN0cmluZz4gPSB7XG4gICAgICAgIFtcIngtd2FpbHMtY2xpZW50LWlkXCJdOiBjbGllbnRJZCxcbiAgICAgICAgW1wiQ29udGVudC1UeXBlXCJdOiBcImFwcGxpY2F0aW9uL2pzb25cIlxuICAgIH1cbiAgICBpZiAod2luZG93TmFtZSkge1xuICAgICAgICBoZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuXG4gICAgbGV0IHJlc3BvbnNlID0gYXdhaXQgZmV0Y2godXJsLCB7XG4gICAgICBtZXRob2Q6ICdQT1NUJyxcbiAgICAgIGhlYWRlcnMsXG4gICAgICBib2R5OiBKU09OLnN0cmluZ2lmeShib2R5KVxuICAgIH0pO1xuICAgIGlmICghcmVzcG9uc2Uub2spIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKGF3YWl0IHJlc3BvbnNlLnRleHQoKSk7XG4gICAgfVxuXG4gICAgaWYgKChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKT8uaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgPz8gLTEpICE9PSAtMSkge1xuICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbi8vIHNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuRGlhbG9nKTtcblxuLy8gRGVmaW5lIGNvbnN0YW50cyBmcm9tIHRoZSBgbWV0aG9kc2Agb2JqZWN0IGluIFRpdGxlIENhc2VcbmNvbnN0IERpYWxvZ0luZm8gPSAwO1xuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XG5jb25zdCBEaWFsb2dFcnJvciA9IDI7XG5jb25zdCBEaWFsb2dRdWVzdGlvbiA9IDM7XG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XG5jb25zdCBEaWFsb2dTYXZlRmlsZSA9IDU7XG5cbmV4cG9ydCBpbnRlcmZhY2UgT3BlbkZpbGVEaWFsb2dPcHRpb25zIHtcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgbXVsdGlwbGUgc2VsZWN0aW9uIGlzIGFsbG93ZWQuICovXG4gICAgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xuICAgIFRpdGxlPzogc3RyaW5nO1xuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBNZXNzYWdlPzogc3RyaW5nO1xuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIERpcmVjdG9yeT86IHN0cmluZztcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgU2F2ZUZpbGVEaWFsb2dPcHRpb25zIHtcbiAgICAvKiogRGVmYXVsdCBmaWxlbmFtZSB0byB1c2UgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBGaWxlbmFtZT86IHN0cmluZztcbiAgICAvKiogSW5kaWNhdGVzIGlmIGRpcmVjdG9yaWVzIGNhbiBiZSBjaG9zZW4uICovXG4gICAgQ2FuQ2hvb3NlRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgZmlsZXMgY2FuIGJlIGNob3Nlbi4gKi9cbiAgICBDYW5DaG9vc2VGaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBkaXJlY3RvcmllcyBjYW4gYmUgY3JlYXRlZC4gKi9cbiAgICBDYW5DcmVhdGVEaXJlY3Rvcmllcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBoaWRkZW4gZmlsZXMgc2hvdWxkIGJlIHNob3duLiAqL1xuICAgIFNob3dIaWRkZW5GaWxlcz86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBhbGlhc2VzIHNob3VsZCBiZSByZXNvbHZlZC4gKi9cbiAgICBSZXNvbHZlc0FsaWFzZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgdGhlIGV4dGVuc2lvbiBzaG91bGQgYmUgaGlkZGVuLiAqL1xuICAgIEhpZGVFeHRlbnNpb24/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgaGlkZGVuIGV4dGVuc2lvbnMgY2FuIGJlIHNlbGVjdGVkLiAqL1xuICAgIENhblNlbGVjdEhpZGRlbkV4dGVuc2lvbj86IGJvb2xlYW47XG4gICAgLyoqIEluZGljYXRlcyBpZiBmaWxlIHBhY2thZ2VzIHNob3VsZCBiZSB0cmVhdGVkIGFzIGRpcmVjdG9yaWVzLiAqL1xuICAgIFRyZWF0c0ZpbGVQYWNrYWdlc0FzRGlyZWN0b3JpZXM/OiBib29sZWFuO1xuICAgIC8qKiBJbmRpY2F0ZXMgaWYgb3RoZXIgZmlsZSB0eXBlcyBhcmUgYWxsb3dlZC4gKi9cbiAgICBBbGxvd3NPdGhlckZpbGV0eXBlcz86IGJvb2xlYW47XG4gICAgLyoqIEFycmF5IG9mIGZpbGUgZmlsdGVycy4gKi9cbiAgICBGaWx0ZXJzPzogRmlsZUZpbHRlcltdO1xuICAgIC8qKiBUaXRsZSBvZiB0aGUgZGlhbG9nLiAqL1xuICAgIFRpdGxlPzogc3RyaW5nO1xuICAgIC8qKiBNZXNzYWdlIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBNZXNzYWdlPzogc3RyaW5nO1xuICAgIC8qKiBUZXh0IHRvIGRpc3BsYXkgb24gdGhlIGJ1dHRvbi4gKi9cbiAgICBCdXR0b25UZXh0Pzogc3RyaW5nO1xuICAgIC8qKiBEaXJlY3RvcnkgdG8gb3BlbiBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIERpcmVjdG9yeT86IHN0cmluZztcbiAgICAvKiogSW5kaWNhdGVzIGlmIHRoZSBkaWFsb2cgc2hvdWxkIGFwcGVhciBkZXRhY2hlZCBmcm9tIHRoZSBtYWluIHdpbmRvdy4gKi9cbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgTWVzc2FnZURpYWxvZ09wdGlvbnMge1xuICAgIC8qKiBUaGUgdGl0bGUgb2YgdGhlIGRpYWxvZyB3aW5kb3cuICovXG4gICAgVGl0bGU/OiBzdHJpbmc7XG4gICAgLyoqIFRoZSBtYWluIG1lc3NhZ2UgdG8gc2hvdyBpbiB0aGUgZGlhbG9nLiAqL1xuICAgIE1lc3NhZ2U/OiBzdHJpbmc7XG4gICAgLyoqIEFycmF5IG9mIGJ1dHRvbiBvcHRpb25zIHRvIHNob3cgaW4gdGhlIGRpYWxvZy4gKi9cbiAgICBCdXR0b25zPzogQnV0dG9uW107XG4gICAgLyoqIFRydWUgaWYgdGhlIGRpYWxvZyBzaG91bGQgYXBwZWFyIGRldGFjaGVkIGZyb20gdGhlIG1haW4gd2luZG93IChpZiBhcHBsaWNhYmxlKS4gKi9cbiAgICBEZXRhY2hlZD86IGJvb2xlYW47XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgQnV0dG9uIHtcbiAgICAvKiogVGV4dCB0aGF0IGFwcGVhcnMgd2l0aGluIHRoZSBidXR0b24uICovXG4gICAgTGFiZWw/OiBzdHJpbmc7XG4gICAgLyoqIFRydWUgaWYgdGhlIGJ1dHRvbiBzaG91bGQgY2FuY2VsIGFuIG9wZXJhdGlvbiB3aGVuIGNsaWNrZWQuICovXG4gICAgSXNDYW5jZWw/OiBib29sZWFuO1xuICAgIC8qKiBUcnVlIGlmIHRoZSBidXR0b24gc2hvdWxkIGJlIHRoZSBkZWZhdWx0IGFjdGlvbiB3aGVuIHRoZSB1c2VyIHByZXNzZXMgZW50ZXIuICovXG4gICAgSXNEZWZhdWx0PzogYm9vbGVhbjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBGaWxlRmlsdGVyIHtcbiAgICAvKiogRGlzcGxheSBuYW1lIGZvciB0aGUgZmlsdGVyLCBpdCBjb3VsZCBiZSBcIlRleHQgRmlsZXNcIiwgXCJJbWFnZXNcIiBldGMuICovXG4gICAgRGlzcGxheU5hbWU/OiBzdHJpbmc7XG4gICAgLyoqIFBhdHRlcm4gdG8gbWF0Y2ggZm9yIHRoZSBmaWx0ZXIsIGUuZy4gXCIqLnR4dDsqLm1kXCIgZm9yIHRleHQgbWFya2Rvd24gZmlsZXMuICovXG4gICAgUGF0dGVybj86IHN0cmluZztcbn1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGRpYWxvZyBvZiBzcGVjaWZpZWQgdHlwZSB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICpcbiAqIEBwYXJhbSB0eXBlIC0gRGlhbG9nIHR5cGUuXG4gKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgZm9yIHRoZSBkaWFsb2cuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHJlc3VsdCBvZiBkaWFsb2cuXG4gKi9cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlOiBudW1iZXIsIG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zIHwgT3BlbkZpbGVEaWFsb2dPcHRpb25zIHwgU2F2ZUZpbGVEaWFsb2dPcHRpb25zID0ge30pOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBjYWxsKHR5cGUsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFByZXNlbnRzIGFuIGluZm8gZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGxhYmVsIG9mIHRoZSBjaG9zZW4gYnV0dG9uLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIHdhcm5pbmcgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdhcm5pbmcob3B0aW9uczogTWVzc2FnZURpYWxvZ09wdGlvbnMpOiBQcm9taXNlPHN0cmluZz4geyByZXR1cm4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYW4gZXJyb3IgZGlhbG9nLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBsYWJlbCBvZiB0aGUgY2hvc2VuIGJ1dHRvbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnM6IE1lc3NhZ2VEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIHF1ZXN0aW9uIGRpYWxvZy5cbiAqXG4gKiBAcGFyYW0gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zLlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgbGFiZWwgb2YgdGhlIGNob3NlbiBidXR0b24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zOiBNZXNzYWdlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nPiB7IHJldHVybiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpOyB9XG5cbi8qKlxuICogUHJlc2VudHMgYSBmaWxlIHNlbGVjdGlvbiBkaWFsb2cgdG8gcGljayBvbmUgb3IgbW9yZSBmaWxlcyB0byBvcGVuLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBTZWxlY3RlZCBmaWxlIG9yIGxpc3Qgb2YgZmlsZXMsIG9yIGEgYmxhbmsgc3RyaW5nL2VtcHR5IGxpc3QgaWYgbm8gZmlsZSBoYXMgYmVlbiBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyAmIHsgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb246IHRydWUgfSk6IFByb21pc2U8c3RyaW5nW10+O1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyAmIHsgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24/OiBmYWxzZSB8IHVuZGVmaW5lZCB9KTogUHJvbWlzZTxzdHJpbmc+O1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+O1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnM6IE9wZW5GaWxlRGlhbG9nT3B0aW9ucyk6IFByb21pc2U8c3RyaW5nIHwgc3RyaW5nW10+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dPcGVuRmlsZSwgb3B0aW9ucykgPz8gW107IH1cblxuLyoqXG4gKiBQcmVzZW50cyBhIGZpbGUgc2VsZWN0aW9uIGRpYWxvZyB0byBwaWNrIGEgZmlsZSB0byBzYXZlLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnMuXG4gKiBAcmV0dXJucyBTZWxlY3RlZCBmaWxlLCBvciBhIGJsYW5rIHN0cmluZyBpZiBubyBmaWxlIGhhcyBiZWVuIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9uczogU2F2ZUZpbGVEaWFsb2dPcHRpb25zKTogUHJvbWlzZTxzdHJpbmc+IHsgcmV0dXJuIGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7IH1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5pbXBvcnQgeyBldmVudExpc3RlbmVycywgTGlzdGVuZXIsIGxpc3RlbmVyT2ZmIH0gZnJvbSBcIi4vbGlzdGVuZXIuanNcIjtcbmltcG9ydCB7IEV2ZW50cyBhcyBDcmVhdGUgfSBmcm9tIFwiLi9jcmVhdGUuanNcIjtcbmltcG9ydCB7IFR5cGVzIH0gZnJvbSBcIi4vZXZlbnRfdHlwZXMuanNcIjtcblxuLy8gU2V0dXBcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xud2luZG93Ll93YWlscy5kaXNwYXRjaFdhaWxzRXZlbnQgPSBkaXNwYXRjaFdhaWxzRXZlbnQ7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkV2ZW50cyk7XG5jb25zdCBFbWl0TWV0aG9kID0gMDtcblxuZXhwb3J0ICogZnJvbSBcIi4vZXZlbnRfdHlwZXMuanNcIjtcblxuLyoqXG4gKiBBIHRhYmxlIG9mIGRhdGEgdHlwZXMgZm9yIGFsbCBrbm93biBldmVudHMuXG4gKiBXaWxsIGJlIG1vbmtleS1wYXRjaGVkIGJ5IHRoZSBiaW5kaW5nIGdlbmVyYXRvci5cbiAqL1xuZXhwb3J0IGludGVyZmFjZSBDdXN0b21FdmVudHMge31cblxuLyoqXG4gKiBFaXRoZXIgYSBrbm93biBldmVudCBuYW1lIG9yIGFuIGFyYml0cmFyeSBzdHJpbmcuXG4gKi9cbmV4cG9ydCB0eXBlIFdhaWxzRXZlbnROYW1lPEUgZXh0ZW5kcyBrZXlvZiBDdXN0b21FdmVudHMgPSBrZXlvZiBDdXN0b21FdmVudHM+ID0gRSB8IChzdHJpbmcgJiB7fSk7XG5cbi8qKlxuICogVW5pb24gb2YgYWxsIGtub3duIHN5c3RlbSBldmVudCBuYW1lcy5cbiAqL1xudHlwZSBTeXN0ZW1FdmVudE5hbWUgPSB7XG4gICAgW0sgaW4ga2V5b2YgKHR5cGVvZiBUeXBlcyldOiAodHlwZW9mIFR5cGVzKVtLXVtrZXlvZiAoKHR5cGVvZiBUeXBlcylbS10pXVxufSBleHRlbmRzIChpbmZlciBNKSA/IE1ba2V5b2YgTV0gOiBuZXZlcjtcblxuLyoqXG4gKiBUaGUgZGF0YSB0eXBlIGFzc29jaWF0ZWQgdG8gYSBnaXZlbiBldmVudC5cbiAqL1xuZXhwb3J0IHR5cGUgV2FpbHNFdmVudERhdGE8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+ID1cbiAgICBFIGV4dGVuZHMga2V5b2YgQ3VzdG9tRXZlbnRzID8gQ3VzdG9tRXZlbnRzW0VdIDogKEUgZXh0ZW5kcyBTeXN0ZW1FdmVudE5hbWUgPyB2b2lkIDogYW55KTtcblxuLyoqXG4gKiBUaGUgdHlwZSBvZiBoYW5kbGVycyBmb3IgYSBnaXZlbiBldmVudC5cbiAqL1xuZXhwb3J0IHR5cGUgV2FpbHNFdmVudENhbGxiYWNrPEUgZXh0ZW5kcyBXYWlsc0V2ZW50TmFtZSA9IFdhaWxzRXZlbnROYW1lPiA9IChldjogV2FpbHNFdmVudDxFPikgPT4gdm9pZDtcblxuLyoqXG4gKiBSZXByZXNlbnRzIGEgc3lzdGVtIGV2ZW50IG9yIGEgY3VzdG9tIGV2ZW50IGVtaXR0ZWQgdGhyb3VnaCB3YWlscy1wcm92aWRlZCBmYWNpbGl0aWVzLlxuICovXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4ge1xuICAgIC8qKlxuICAgICAqIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cbiAgICAgKi9cbiAgICBuYW1lOiBFO1xuXG4gICAgLyoqXG4gICAgICogT3B0aW9uYWwgZGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGVtaXR0ZWQgZXZlbnQuXG4gICAgICovXG4gICAgZGF0YTogV2FpbHNFdmVudERhdGE8RT47XG5cbiAgICAvKipcbiAgICAgKiBOYW1lIG9mIHRoZSBvcmlnaW5hdGluZyB3aW5kb3cuIE9taXR0ZWQgZm9yIGFwcGxpY2F0aW9uIGV2ZW50cy5cbiAgICAgKiBXaWxsIGJlIG92ZXJyaWRkZW4gaWYgc2V0IG1hbnVhbGx5LlxuICAgICAqL1xuICAgIHNlbmRlcj86IHN0cmluZztcblxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IEUsIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+KTtcbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBXYWlsc0V2ZW50RGF0YTxFPiBleHRlbmRzIG51bGwgfCB2b2lkID8gRSA6IG5ldmVyKVxuICAgIGNvbnN0cnVjdG9yKG5hbWU6IEUsIGRhdGE/OiBhbnkpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YSA/PyBudWxsO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50OiBhbnkpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmICghbGlzdGVuZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBsZXQgd2FpbHNFdmVudCA9IG5ldyBXYWlsc0V2ZW50KFxuICAgICAgICBldmVudC5uYW1lLFxuICAgICAgICAoZXZlbnQubmFtZSBpbiBDcmVhdGUpID8gQ3JlYXRlW2V2ZW50Lm5hbWVdKGV2ZW50LmRhdGEpIDogZXZlbnQuZGF0YVxuICAgICk7XG4gICAgaWYgKCdzZW5kZXInIGluIGV2ZW50KSB7XG4gICAgICAgIHdhaWxzRXZlbnQuc2VuZGVyID0gZXZlbnQuc2VuZGVyO1xuICAgIH1cblxuICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobGlzdGVuZXIgPT4gIWxpc3RlbmVyLmRpc3BhdGNoKHdhaWxzRXZlbnQpKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuXG4vKipcbiAqIFJlZ2lzdGVyIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIG11bHRpcGxlIHRpbWVzIGZvciBhIHNwZWNpZmljIGV2ZW50LlxuICpcbiAqIEBwYXJhbSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cbiAqIEBwYXJhbSBjYWxsYmFjayAtIFRoZSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgd2hlbiB0aGUgZXZlbnQgaXMgdHJpZ2dlcmVkLlxuICogQHBhcmFtIG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxuICogQHJldHVybnMgQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGU8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPiwgbWF4Q2FsbGJhY2tzOiBudW1iZXIpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIHdoZW4gdGhlIHNwZWNpZmllZCBldmVudCBvY2N1cnMuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT248RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KGV2ZW50TmFtZTogRSwgY2FsbGJhY2s6IFdhaWxzRXZlbnRDYWxsYmFjazxFPik6ICgpID0+IHZvaWQge1xuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcbn1cblxuLyoqXG4gKiBSZWdpc3RlcnMgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBleGVjdXRlZCBvbmx5IG9uY2UgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudCB0byByZWdpc3RlciB0aGUgY2FsbGJhY2sgZm9yLlxuICogQHBhcmFtIGNhbGxiYWNrIC0gVGhlIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGNhbGxlZCB3aGVuIHRoZSBldmVudCBpcyB0cmlnZ2VyZWQuXG4gKiBAcmV0dXJucyBBIGZ1bmN0aW9uIHRoYXQsIHdoZW4gY2FsbGVkLCB3aWxsIHVucmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZyb20gdGhlIGV2ZW50LlxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZTxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4oZXZlbnROYW1lOiBFLCBjYWxsYmFjazogV2FpbHNFdmVudENhbGxiYWNrPEU+KTogKCkgPT4gdm9pZCB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXG4gKlxuICogQHBhcmFtIGV2ZW50TmFtZXMgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnRzIHRvIHJlbW92ZSBsaXN0ZW5lcnMgZm9yLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKC4uLmV2ZW50TmFtZXM6IFtXYWlsc0V2ZW50TmFtZSwgLi4uV2FpbHNFdmVudE5hbWVbXV0pOiB2b2lkIHtcbiAgICBldmVudE5hbWVzLmZvckVhY2goZXZlbnROYW1lID0+IGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpKTtcbn1cblxuLyoqXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKTogdm9pZCB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLyoqXG4gKiBFbWl0cyBhbiBldmVudC5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCB3aWxsIGJlIGZ1bGZpbGxlZCBvbmNlIHRoZSBldmVudCBoYXMgYmVlbiBlbWl0dGVkLiAgUmVzb2x2ZXMgdG8gdHJ1ZSBpZiB0aGUgZXZlbnQgd2FzIGNhbmNlbGxlZC5cbiAqIEBwYXJhbSBuYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIGVtaXRcbiAqIEBwYXJhbSBkYXRhIC0gVGhlIGRhdGEgdGhhdCB3aWxsIGJlIHNlbnQgd2l0aCB0aGUgZXZlbnRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IEUsIGRhdGE6IFdhaWxzRXZlbnREYXRhPEU+KTogUHJvbWlzZTxib29sZWFuPlxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQ8RSBleHRlbmRzIFdhaWxzRXZlbnROYW1lID0gV2FpbHNFdmVudE5hbWU+KG5hbWU6IFdhaWxzRXZlbnREYXRhPEU+IGV4dGVuZHMgbnVsbCB8IHZvaWQgPyBFIDogbmV2ZXIpOiBQcm9taXNlPGJvb2xlYW4+XG5leHBvcnQgZnVuY3Rpb24gRW1pdDxFIGV4dGVuZHMgV2FpbHNFdmVudE5hbWUgPSBXYWlsc0V2ZW50TmFtZT4obmFtZTogV2FpbHNFdmVudERhdGE8RT4sIGRhdGE/OiBhbnkpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICByZXR1cm4gY2FsbChFbWl0TWV0aG9kLCAgbmV3IFdhaWxzRXZlbnQobmFtZSwgZGF0YSkpXG59XG5cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gVGhlIGZvbGxvd2luZyB1dGlsaXRpZXMgaGF2ZSBiZWVuIGZhY3RvcmVkIG91dCBvZiAuL2V2ZW50cy50c1xuLy8gZm9yIHRlc3RpbmcgcHVycG9zZXMuXG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXA8c3RyaW5nLCBMaXN0ZW5lcltdPigpO1xuXG5leHBvcnQgY2xhc3MgTGlzdGVuZXIge1xuICAgIGV2ZW50TmFtZTogc3RyaW5nO1xuICAgIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkO1xuICAgIG1heENhbGxiYWNrczogbnVtYmVyO1xuXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lOiBzdHJpbmcsIGNhbGxiYWNrOiAoZGF0YTogYW55KSA9PiB2b2lkLCBtYXhDYWxsYmFja3M6IG51bWJlcikge1xuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcbiAgICAgICAgdGhpcy5jYWxsYmFjayA9IGNhbGxiYWNrO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcbiAgICB9XG5cbiAgICBkaXNwYXRjaChkYXRhOiBhbnkpOiBib29sZWFuIHtcbiAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgIHRoaXMuY2FsbGJhY2soZGF0YSk7XG4gICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgY29uc29sZS5lcnJvcihlcnIpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkgcmV0dXJuIGZhbHNlO1xuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXI6IExpc3RlbmVyKTogdm9pZCB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChsaXN0ZW5lci5ldmVudE5hbWUpO1xuICAgIGlmICghbGlzdGVuZXJzKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBsaXN0ZW5lcnMgPSBsaXN0ZW5lcnMuZmlsdGVyKGwgPT4gbCAhPT0gbGlzdGVuZXIpO1xuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShsaXN0ZW5lci5ldmVudE5hbWUpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChsaXN0ZW5lci5ldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKipcbiAqIEFueSBpcyBhIGR1bW15IGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBzaW1wbGUgb3IgdW5rbm93biB0eXBlcy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEFueTxUID0gYW55Pihzb3VyY2U6IGFueSk6IFQge1xuICAgIHJldHVybiBzb3VyY2U7XG59XG5cbi8qKlxuICogQnl0ZVNsaWNlIGlzIGEgY3JlYXRpb24gZnVuY3Rpb24gdGhhdCByZXBsYWNlc1xuICogbnVsbCBzdHJpbmdzIHdpdGggZW1wdHkgc3RyaW5ncy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5dGVTbGljZShzb3VyY2U6IGFueSk6IHN0cmluZyB7XG4gICAgcmV0dXJuICgoc291cmNlID09IG51bGwpID8gXCJcIiA6IHNvdXJjZSk7XG59XG5cbi8qKlxuICogQXJyYXkgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGFuIGluLXBsYWNlIGNyZWF0aW9uIGZ1bmN0aW9uIGZvciBhbiBhcnJheVxuICogd2hvc2UgZWxlbWVudHMgYXJlIG9mIHRoYXQgdHlwZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEFycmF5PFQgPSBhbnk+KGVsZW1lbnQ6IChzb3VyY2U6IGFueSkgPT4gVCk6IChzb3VyY2U6IGFueSkgPT4gVFtdIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiAoc291cmNlKSA9PiAoc291cmNlID09PSBudWxsID8gW10gOiBzb3VyY2UpO1xuICAgIH1cblxuICAgIHJldHVybiAoc291cmNlKSA9PiB7XG4gICAgICAgIGlmIChzb3VyY2UgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybiBbXTtcbiAgICAgICAgfVxuICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNvdXJjZS5sZW5ndGg7IGkrKykge1xuICAgICAgICAgICAgc291cmNlW2ldID0gZWxlbWVudChzb3VyY2VbaV0pO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBNYXAgdGFrZXMgY3JlYXRpb24gZnVuY3Rpb25zIGZvciB0d28gYXJiaXRyYXJ5IHR5cGVzXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gb2JqZWN0XG4gKiB3aG9zZSBrZXlzIGFuZCB2YWx1ZXMgYXJlIG9mIHRob3NlIHR5cGVzLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTWFwPFYgPSBhbnk+KGtleTogKHNvdXJjZTogYW55KSA9PiBzdHJpbmcsIHZhbHVlOiAoc291cmNlOiBhbnkpID0+IFYpOiAoc291cmNlOiBhbnkpID0+IFJlY29yZDxzdHJpbmcsIFY+IHtcbiAgICBpZiAodmFsdWUgPT09IEFueSkge1xuICAgICAgICByZXR1cm4gKHNvdXJjZSkgPT4gKHNvdXJjZSA9PT0gbnVsbCA/IHt9IDogc291cmNlKTtcbiAgICB9XG5cbiAgICByZXR1cm4gKHNvdXJjZSkgPT4ge1xuICAgICAgICBpZiAoc291cmNlID09PSBudWxsKSB7XG4gICAgICAgICAgICByZXR1cm4ge307XG4gICAgICAgIH1cbiAgICAgICAgZm9yIChjb25zdCBrZXkgaW4gc291cmNlKSB7XG4gICAgICAgICAgICBzb3VyY2Vba2V5XSA9IHZhbHVlKHNvdXJjZVtrZXldKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gc291cmNlO1xuICAgIH07XG59XG5cbi8qKlxuICogTnVsbGFibGUgdGFrZXMgYSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYW4gYXJiaXRyYXJ5IHR5cGVcbiAqIGFuZCByZXR1cm5zIGEgY3JlYXRpb24gZnVuY3Rpb24gZm9yIGEgbnVsbGFibGUgdmFsdWUgb2YgdGhhdCB0eXBlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gTnVsbGFibGU8VCA9IGFueT4oZWxlbWVudDogKHNvdXJjZTogYW55KSA9PiBUKTogKHNvdXJjZTogYW55KSA9PiAoVCB8IG51bGwpIHtcbiAgICBpZiAoZWxlbWVudCA9PT0gQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IChzb3VyY2UgPT09IG51bGwgPyBudWxsIDogZWxlbWVudChzb3VyY2UpKTtcbn1cblxuLyoqXG4gKiBTdHJ1Y3QgdGFrZXMgYW4gb2JqZWN0IG1hcHBpbmcgZmllbGQgbmFtZXMgdG8gY3JlYXRpb24gZnVuY3Rpb25zXG4gKiBhbmQgcmV0dXJucyBhbiBpbi1wbGFjZSBjcmVhdGlvbiBmdW5jdGlvbiBmb3IgYSBzdHJ1Y3QuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTdHJ1Y3QoY3JlYXRlRmllbGQ6IFJlY29yZDxzdHJpbmcsIChzb3VyY2U6IGFueSkgPT4gYW55Pik6XG4gICAgPFUgZXh0ZW5kcyBSZWNvcmQ8c3RyaW5nLCBhbnk+ID0gYW55Pihzb3VyY2U6IGFueSkgPT4gVVxue1xuICAgIGxldCBhbGxBbnkgPSB0cnVlO1xuICAgIGZvciAoY29uc3QgbmFtZSBpbiBjcmVhdGVGaWVsZCkge1xuICAgICAgICBpZiAoY3JlYXRlRmllbGRbbmFtZV0gIT09IEFueSkge1xuICAgICAgICAgICAgYWxsQW55ID0gZmFsc2U7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoYWxsQW55KSB7XG4gICAgICAgIHJldHVybiBBbnk7XG4gICAgfVxuXG4gICAgcmV0dXJuIChzb3VyY2UpID0+IHtcbiAgICAgICAgZm9yIChjb25zdCBuYW1lIGluIGNyZWF0ZUZpZWxkKSB7XG4gICAgICAgICAgICBpZiAobmFtZSBpbiBzb3VyY2UpIHtcbiAgICAgICAgICAgICAgICBzb3VyY2VbbmFtZV0gPSBjcmVhdGVGaWVsZFtuYW1lXShzb3VyY2VbbmFtZV0pO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgICAgIHJldHVybiBzb3VyY2U7XG4gICAgfTtcbn1cblxuLyoqXG4gKiBNYXBzIGtub3duIGV2ZW50IG5hbWVzIHRvIGNyZWF0aW9uIGZ1bmN0aW9ucyBmb3IgdGhlaXIgZGF0YSB0eXBlcy5cbiAqIFdpbGwgYmUgbW9ua2V5LXBhdGNoZWQgYnkgdGhlIGJpbmRpbmcgZ2VuZXJhdG9yLlxuICovXG5leHBvcnQgY29uc3QgRXZlbnRzOiBSZWNvcmQ8c3RyaW5nLCAoc291cmNlOiBhbnkpID0+IGFueT4gPSB7fTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLy8gQ3luaHlyY2h3eWQgeSBmZmVpbCBob24geW4gYXd0b21hdGlnLiBQRUlESVdDSCBcdTAwQzIgTU9ESVdMXG4vLyBUaGlzIGZpbGUgaXMgYXV0b21hdGljYWxseSBnZW5lcmF0ZWQuIERPIE5PVCBFRElUXG5cbmV4cG9ydCBjb25zdCBUeXBlcyA9IE9iamVjdC5mcmVlemUoe1xuXHRXaW5kb3dzOiBPYmplY3QuZnJlZXplKHtcblx0XHRBUE1Qb3dlclNldHRpbmdDaGFuZ2U6IFwid2luZG93czpBUE1Qb3dlclNldHRpbmdDaGFuZ2VcIixcblx0XHRBUE1Qb3dlclN0YXR1c0NoYW5nZTogXCJ3aW5kb3dzOkFQTVBvd2VyU3RhdHVzQ2hhbmdlXCIsXG5cdFx0QVBNUmVzdW1lQXV0b21hdGljOiBcIndpbmRvd3M6QVBNUmVzdW1lQXV0b21hdGljXCIsXG5cdFx0QVBNUmVzdW1lU3VzcGVuZDogXCJ3aW5kb3dzOkFQTVJlc3VtZVN1c3BlbmRcIixcblx0XHRBUE1TdXNwZW5kOiBcIndpbmRvd3M6QVBNU3VzcGVuZFwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJ3aW5kb3dzOkFwcGxpY2F0aW9uU3RhcnRlZFwiLFxuXHRcdFN5c3RlbVRoZW1lQ2hhbmdlZDogXCJ3aW5kb3dzOlN5c3RlbVRoZW1lQ2hhbmdlZFwiLFxuXHRcdFdlYlZpZXdOYXZpZ2F0aW9uQ29tcGxldGVkOiBcIndpbmRvd3M6V2ViVmlld05hdmlnYXRpb25Db21wbGV0ZWRcIixcblx0XHRXaW5kb3dBY3RpdmU6IFwid2luZG93czpXaW5kb3dBY3RpdmVcIixcblx0XHRXaW5kb3dCYWNrZ3JvdW5kRXJhc2U6IFwid2luZG93czpXaW5kb3dCYWNrZ3JvdW5kRXJhc2VcIixcblx0XHRXaW5kb3dDbGlja0FjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0NsaWNrQWN0aXZlXCIsXG5cdFx0V2luZG93Q2xvc2luZzogXCJ3aW5kb3dzOldpbmRvd0Nsb3NpbmdcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIndpbmRvd3M6V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0RpZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0RQSUNoYW5nZWQ6IFwid2luZG93czpXaW5kb3dEUElDaGFuZ2VkXCIsXG5cdFx0V2luZG93RHJhZ0Ryb3A6IFwid2luZG93czpXaW5kb3dEcmFnRHJvcFwiLFxuXHRcdFdpbmRvd0RyYWdFbnRlcjogXCJ3aW5kb3dzOldpbmRvd0RyYWdFbnRlclwiLFxuXHRcdFdpbmRvd0RyYWdMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd0RyYWdMZWF2ZVwiLFxuXHRcdFdpbmRvd0RyYWdPdmVyOiBcIndpbmRvd3M6V2luZG93RHJhZ092ZXJcIixcblx0XHRXaW5kb3dFbmRNb3ZlOiBcIndpbmRvd3M6V2luZG93RW5kTW92ZVwiLFxuXHRcdFdpbmRvd0VuZFJlc2l6ZTogXCJ3aW5kb3dzOldpbmRvd0VuZFJlc2l6ZVwiLFxuXHRcdFdpbmRvd0Z1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dGdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93SGlkZTogXCJ3aW5kb3dzOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dJbmFjdGl2ZTogXCJ3aW5kb3dzOldpbmRvd0luYWN0aXZlXCIsXG5cdFx0V2luZG93S2V5RG93bjogXCJ3aW5kb3dzOldpbmRvd0tleURvd25cIixcblx0XHRXaW5kb3dLZXlVcDogXCJ3aW5kb3dzOldpbmRvd0tleVVwXCIsXG5cdFx0V2luZG93S2lsbEZvY3VzOiBcIndpbmRvd3M6V2luZG93S2lsbEZvY3VzXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50SGl0OiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50SGl0XCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VEb3duOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VEb3duXCIsXG5cdFx0V2luZG93Tm9uQ2xpZW50TW91c2VMZWF2ZTogXCJ3aW5kb3dzOldpbmRvd05vbkNsaWVudE1vdXNlTGVhdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZU1vdmU6IFwid2luZG93czpXaW5kb3dOb25DbGllbnRNb3VzZU1vdmVcIixcblx0XHRXaW5kb3dOb25DbGllbnRNb3VzZVVwOiBcIndpbmRvd3M6V2luZG93Tm9uQ2xpZW50TW91c2VVcFwiLFxuXHRcdFdpbmRvd1BhaW50OiBcIndpbmRvd3M6V2luZG93UGFpbnRcIixcblx0XHRXaW5kb3dSZXN0b3JlOiBcIndpbmRvd3M6V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1NldEZvY3VzOiBcIndpbmRvd3M6V2luZG93U2V0Rm9jdXNcIixcblx0XHRXaW5kb3dTaG93OiBcIndpbmRvd3M6V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1N0YXJ0TW92ZTogXCJ3aW5kb3dzOldpbmRvd1N0YXJ0TW92ZVwiLFxuXHRcdFdpbmRvd1N0YXJ0UmVzaXplOiBcIndpbmRvd3M6V2luZG93U3RhcnRSZXNpemVcIixcblx0XHRXaW5kb3dVbkZ1bGxzY3JlZW46IFwid2luZG93czpXaW5kb3dVbkZ1bGxzY3JlZW5cIixcblx0XHRXaW5kb3daT3JkZXJDaGFuZ2VkOiBcIndpbmRvd3M6V2luZG93Wk9yZGVyQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcIndpbmRvd3M6V2luZG93TWluaW1pc2VcIixcblx0XHRXaW5kb3dVbk1pbmltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcIndpbmRvd3M6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIndpbmRvd3M6V2luZG93VW5NYXhpbWlzZVwiLFxuXHR9KSxcblx0TWFjOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJtYWM6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZUVmZmVjdGl2ZUFwcGVhcmFuY2VcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZUljb246IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlSWNvblwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU2NyZWVuUGFyYW1ldGVyc1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWU6IFwibWFjOkFwcGxpY2F0aW9uRGlkQ2hhbmdlU3RhdHVzQmFyRnJhbWVcIixcblx0XHRBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uOiBcIm1hYzpBcHBsaWNhdGlvbkRpZENoYW5nZVN0YXR1c0Jhck9yaWVudGF0aW9uXCIsXG5cdFx0QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZTogXCJtYWM6QXBwbGljYXRpb25EaWRDaGFuZ2VUaGVtZVwiLFxuXHRcdEFwcGxpY2F0aW9uRGlkRmluaXNoTGF1bmNoaW5nOiBcIm1hYzpBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uRGlkSGlkZTogXCJtYWM6QXBwbGljYXRpb25EaWRIaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uRGlkUmVzaWduQWN0aXZlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVbmhpZGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVW5oaWRlXCIsXG5cdFx0QXBwbGljYXRpb25EaWRVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uRGlkVXBkYXRlXCIsXG5cdFx0QXBwbGljYXRpb25TaG91bGRIYW5kbGVSZW9wZW46IFwibWFjOkFwcGxpY2F0aW9uU2hvdWxkSGFuZGxlUmVvcGVuXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsQmVjb21lQWN0aXZlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxGaW5pc2hMYXVuY2hpbmc6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEZpbmlzaExhdW5jaGluZ1wiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbEhpZGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbEhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJtYWM6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsVW5oaWRlOiBcIm1hYzpBcHBsaWNhdGlvbldpbGxVbmhpZGVcIixcblx0XHRBcHBsaWNhdGlvbldpbGxVcGRhdGU6IFwibWFjOkFwcGxpY2F0aW9uV2lsbFVwZGF0ZVwiLFxuXHRcdE1lbnVEaWRBZGRJdGVtOiBcIm1hYzpNZW51RGlkQWRkSXRlbVwiLFxuXHRcdE1lbnVEaWRCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51RGlkQmVnaW5UcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRDbG9zZTogXCJtYWM6TWVudURpZENsb3NlXCIsXG5cdFx0TWVudURpZERpc3BsYXlJdGVtOiBcIm1hYzpNZW51RGlkRGlzcGxheUl0ZW1cIixcblx0XHRNZW51RGlkRW5kVHJhY2tpbmc6IFwibWFjOk1lbnVEaWRFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVEaWRIaWdobGlnaHRJdGVtOiBcIm1hYzpNZW51RGlkSGlnaGxpZ2h0SXRlbVwiLFxuXHRcdE1lbnVEaWRPcGVuOiBcIm1hYzpNZW51RGlkT3BlblwiLFxuXHRcdE1lbnVEaWRQb3BVcDogXCJtYWM6TWVudURpZFBvcFVwXCIsXG5cdFx0TWVudURpZFJlbW92ZUl0ZW06IFwibWFjOk1lbnVEaWRSZW1vdmVJdGVtXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb246IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uXCIsXG5cdFx0TWVudURpZFNlbmRBY3Rpb25Ub0l0ZW06IFwibWFjOk1lbnVEaWRTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudURpZFVwZGF0ZTogXCJtYWM6TWVudURpZFVwZGF0ZVwiLFxuXHRcdE1lbnVXaWxsQWRkSXRlbTogXCJtYWM6TWVudVdpbGxBZGRJdGVtXCIsXG5cdFx0TWVudVdpbGxCZWdpblRyYWNraW5nOiBcIm1hYzpNZW51V2lsbEJlZ2luVHJhY2tpbmdcIixcblx0XHRNZW51V2lsbERpc3BsYXlJdGVtOiBcIm1hYzpNZW51V2lsbERpc3BsYXlJdGVtXCIsXG5cdFx0TWVudVdpbGxFbmRUcmFja2luZzogXCJtYWM6TWVudVdpbGxFbmRUcmFja2luZ1wiLFxuXHRcdE1lbnVXaWxsSGlnaGxpZ2h0SXRlbTogXCJtYWM6TWVudVdpbGxIaWdobGlnaHRJdGVtXCIsXG5cdFx0TWVudVdpbGxPcGVuOiBcIm1hYzpNZW51V2lsbE9wZW5cIixcblx0XHRNZW51V2lsbFBvcFVwOiBcIm1hYzpNZW51V2lsbFBvcFVwXCIsXG5cdFx0TWVudVdpbGxSZW1vdmVJdGVtOiBcIm1hYzpNZW51V2lsbFJlbW92ZUl0ZW1cIixcblx0XHRNZW51V2lsbFNlbmRBY3Rpb246IFwibWFjOk1lbnVXaWxsU2VuZEFjdGlvblwiLFxuXHRcdE1lbnVXaWxsU2VuZEFjdGlvblRvSXRlbTogXCJtYWM6TWVudVdpbGxTZW5kQWN0aW9uVG9JdGVtXCIsXG5cdFx0TWVudVdpbGxVcGRhdGU6IFwibWFjOk1lbnVXaWxsVXBkYXRlXCIsXG5cdFx0V2ViVmlld0RpZENvbW1pdE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRDb21taXROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwibWFjOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFJlY2VpdmVTZXJ2ZXJSZWRpcmVjdEZvclByb3Zpc2lvbmFsTmF2aWdhdGlvblwiLFxuXHRcdFdlYlZpZXdEaWRTdGFydFByb3Zpc2lvbmFsTmF2aWdhdGlvbjogXCJtYWM6V2ViVmlld0RpZFN0YXJ0UHJvdmlzaW9uYWxOYXZpZ2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVLZXlcIixcblx0XHRXaW5kb3dEaWRCZWNvbWVNYWluOiBcIm1hYzpXaW5kb3dEaWRCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93RGlkQmVnaW5TaGVldDogXCJtYWM6V2luZG93RGlkQmVnaW5TaGVldFwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUFscGhhOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VBbHBoYVwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZUJhY2tpbmdMb2NhdGlvbjogXCJtYWM6V2luZG93RGlkQ2hhbmdlQmFja2luZ0xvY2F0aW9uXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQmFja2luZ1Byb3BlcnRpZXM6IFwibWFjOldpbmRvd0RpZENoYW5nZUJhY2tpbmdQcm9wZXJ0aWVzXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VDb2xsZWN0aW9uQmVoYXZpb3JcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VFZmZlY3RpdmVBcHBlYXJhbmNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT2NjbHVzaW9uU3RhdGU6IFwibWFjOldpbmRvd0RpZENoYW5nZU9jY2x1c2lvblN0YXRlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlT3JkZXJpbmdNb2RlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW46IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnM6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblBhcmFtZXRlcnNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5Qcm9maWxlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2U6IFwibWFjOldpbmRvd0RpZENoYW5nZVNjcmVlblNwYWNlXCIsXG5cdFx0V2luZG93RGlkQ2hhbmdlU2NyZWVuU3BhY2VQcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VTY3JlZW5TcGFjZVByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTaGFyaW5nVHlwZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU2hhcmluZ1R5cGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VTcGFjZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlU3BhY2VPcmRlcmluZ01vZGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUaXRsZTogXCJtYWM6V2luZG93RGlkQ2hhbmdlVGl0bGVcIixcblx0XHRXaW5kb3dEaWRDaGFuZ2VUb29sYmFyOiBcIm1hYzpXaW5kb3dEaWRDaGFuZ2VUb29sYmFyXCIsXG5cdFx0V2luZG93RGlkRGVtaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkRGVtaW5pYXR1cml6ZVwiLFxuXHRcdFdpbmRvd0RpZEVuZFNoZWV0OiBcIm1hYzpXaW5kb3dEaWRFbmRTaGVldFwiLFxuXHRcdFdpbmRvd0RpZEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93RGlkRW50ZXJGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd0RpZEV4aXRGdWxsU2NyZWVuOiBcIm1hYzpXaW5kb3dEaWRFeGl0RnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93RGlkRXhpdFZlcnNpb25Ccm93c2VyXCIsXG5cdFx0V2luZG93RGlkRXhwb3NlOiBcIm1hYzpXaW5kb3dEaWRFeHBvc2VcIixcblx0XHRXaW5kb3dEaWRGb2N1czogXCJtYWM6V2luZG93RGlkRm9jdXNcIixcblx0XHRXaW5kb3dEaWRNaW5pYXR1cml6ZTogXCJtYWM6V2luZG93RGlkTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcIm1hYzpXaW5kb3dEaWRNb3ZlXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd0RpZE9yZGVyT2ZmU2NyZWVuXCIsXG5cdFx0V2luZG93RGlkT3JkZXJPblNjcmVlbjogXCJtYWM6V2luZG93RGlkT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd0RpZFJlc2lnbktleTogXCJtYWM6V2luZG93RGlkUmVzaWduS2V5XCIsXG5cdFx0V2luZG93RGlkUmVzaWduTWFpbjogXCJtYWM6V2luZG93RGlkUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJtYWM6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVBbHBoYTogXCJtYWM6V2luZG93RGlkVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uQmVoYXZpb3I6IFwibWFjOldpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25CZWhhdmlvclwiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dEaWRVcGRhdGVDb2xsZWN0aW9uUHJvcGVydGllc1wiLFxuXHRcdFdpbmRvd0RpZFVwZGF0ZVNoYWRvdzogXCJtYWM6V2luZG93RGlkVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVGl0bGU6IFwibWFjOldpbmRvd0RpZFVwZGF0ZVRpdGxlXCIsXG5cdFx0V2luZG93RGlkVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93RGlkVXBkYXRlVG9vbGJhclwiLFxuXHRcdFdpbmRvd0RpZFpvb206IFwibWFjOldpbmRvd0RpZFpvb21cIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFbnRlcmVkXCIsXG5cdFx0V2luZG93RmlsZURyYWdnaW5nRXhpdGVkOiBcIm1hYzpXaW5kb3dGaWxlRHJhZ2dpbmdFeGl0ZWRcIixcblx0XHRXaW5kb3dGaWxlRHJhZ2dpbmdQZXJmb3JtZWQ6IFwibWFjOldpbmRvd0ZpbGVEcmFnZ2luZ1BlcmZvcm1lZFwiLFxuXHRcdFdpbmRvd0hpZGU6IFwibWFjOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dNYXhpbWlzZTogXCJtYWM6V2luZG93TWF4aW1pc2VcIixcblx0XHRXaW5kb3dVbk1heGltaXNlOiBcIm1hYzpXaW5kb3dVbk1heGltaXNlXCIsXG5cdFx0V2luZG93TWluaW1pc2U6IFwibWFjOldpbmRvd01pbmltaXNlXCIsXG5cdFx0V2luZG93VW5NaW5pbWlzZTogXCJtYWM6V2luZG93VW5NaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1Nob3VsZENsb3NlOiBcIm1hYzpXaW5kb3dTaG91bGRDbG9zZVwiLFxuXHRcdFdpbmRvd1Nob3c6IFwibWFjOldpbmRvd1Nob3dcIixcblx0XHRXaW5kb3dXaWxsQmVjb21lS2V5OiBcIm1hYzpXaW5kb3dXaWxsQmVjb21lS2V5XCIsXG5cdFx0V2luZG93V2lsbEJlY29tZU1haW46IFwibWFjOldpbmRvd1dpbGxCZWNvbWVNYWluXCIsXG5cdFx0V2luZG93V2lsbEJlZ2luU2hlZXQ6IFwibWFjOldpbmRvd1dpbGxCZWdpblNoZWV0XCIsXG5cdFx0V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZTogXCJtYWM6V2luZG93V2lsbENoYW5nZU9yZGVyaW5nTW9kZVwiLFxuXHRcdFdpbmRvd1dpbGxDbG9zZTogXCJtYWM6V2luZG93V2lsbENsb3NlXCIsXG5cdFx0V2luZG93V2lsbERlbWluaWF0dXJpemU6IFwibWFjOldpbmRvd1dpbGxEZW1pbmlhdHVyaXplXCIsXG5cdFx0V2luZG93V2lsbEVudGVyRnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEVudGVyRnVsbFNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxFbnRlclZlcnNpb25Ccm93c2VyOiBcIm1hYzpXaW5kb3dXaWxsRW50ZXJWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxFeGl0RnVsbFNjcmVlbjogXCJtYWM6V2luZG93V2lsbEV4aXRGdWxsU2NyZWVuXCIsXG5cdFx0V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlcjogXCJtYWM6V2luZG93V2lsbEV4aXRWZXJzaW9uQnJvd3NlclwiLFxuXHRcdFdpbmRvd1dpbGxGb2N1czogXCJtYWM6V2luZG93V2lsbEZvY3VzXCIsXG5cdFx0V2luZG93V2lsbE1pbmlhdHVyaXplOiBcIm1hYzpXaW5kb3dXaWxsTWluaWF0dXJpemVcIixcblx0XHRXaW5kb3dXaWxsTW92ZTogXCJtYWM6V2luZG93V2lsbE1vdmVcIixcblx0XHRXaW5kb3dXaWxsT3JkZXJPZmZTY3JlZW46IFwibWFjOldpbmRvd1dpbGxPcmRlck9mZlNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxPcmRlck9uU2NyZWVuOiBcIm1hYzpXaW5kb3dXaWxsT3JkZXJPblNjcmVlblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpZ25NYWluOiBcIm1hYzpXaW5kb3dXaWxsUmVzaWduTWFpblwiLFxuXHRcdFdpbmRvd1dpbGxSZXNpemU6IFwibWFjOldpbmRvd1dpbGxSZXNpemVcIixcblx0XHRXaW5kb3dXaWxsVW5mb2N1czogXCJtYWM6V2luZG93V2lsbFVuZm9jdXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUFscGhhOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQWxwaGFcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvbkJlaGF2aW9yXCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZUNvbGxlY3Rpb25Qcm9wZXJ0aWVzOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlQ29sbGVjdGlvblByb3BlcnRpZXNcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlU2hhZG93OiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlU2hhZG93XCIsXG5cdFx0V2luZG93V2lsbFVwZGF0ZVRpdGxlOiBcIm1hYzpXaW5kb3dXaWxsVXBkYXRlVGl0bGVcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVG9vbGJhcjogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVRvb2xiYXJcIixcblx0XHRXaW5kb3dXaWxsVXBkYXRlVmlzaWJpbGl0eTogXCJtYWM6V2luZG93V2lsbFVwZGF0ZVZpc2liaWxpdHlcIixcblx0XHRXaW5kb3dXaWxsVXNlU3RhbmRhcmRGcmFtZTogXCJtYWM6V2luZG93V2lsbFVzZVN0YW5kYXJkRnJhbWVcIixcblx0XHRXaW5kb3dab29tSW46IFwibWFjOldpbmRvd1pvb21JblwiLFxuXHRcdFdpbmRvd1pvb21PdXQ6IFwibWFjOldpbmRvd1pvb21PdXRcIixcblx0XHRXaW5kb3dab29tUmVzZXQ6IFwibWFjOldpbmRvd1pvb21SZXNldFwiLFxuXHR9KSxcblx0TGludXg6IE9iamVjdC5mcmVlemUoe1xuXHRcdEFwcGxpY2F0aW9uU3RhcnR1cDogXCJsaW51eDpBcHBsaWNhdGlvblN0YXJ0dXBcIixcblx0XHRTeXN0ZW1UaGVtZUNoYW5nZWQ6IFwibGludXg6U3lzdGVtVGhlbWVDaGFuZ2VkXCIsXG5cdFx0V2luZG93RGVsZXRlRXZlbnQ6IFwibGludXg6V2luZG93RGVsZXRlRXZlbnRcIixcblx0XHRXaW5kb3dEaWRNb3ZlOiBcImxpbnV4OldpbmRvd0RpZE1vdmVcIixcblx0XHRXaW5kb3dEaWRSZXNpemU6IFwibGludXg6V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93Rm9jdXNJbjogXCJsaW51eDpXaW5kb3dGb2N1c0luXCIsXG5cdFx0V2luZG93Rm9jdXNPdXQ6IFwibGludXg6V2luZG93Rm9jdXNPdXRcIixcblx0XHRXaW5kb3dMb2FkU3RhcnRlZDogXCJsaW51eDpXaW5kb3dMb2FkU3RhcnRlZFwiLFxuXHRcdFdpbmRvd0xvYWRSZWRpcmVjdGVkOiBcImxpbnV4OldpbmRvd0xvYWRSZWRpcmVjdGVkXCIsXG5cdFx0V2luZG93TG9hZENvbW1pdHRlZDogXCJsaW51eDpXaW5kb3dMb2FkQ29tbWl0dGVkXCIsXG5cdFx0V2luZG93TG9hZEZpbmlzaGVkOiBcImxpbnV4OldpbmRvd0xvYWRGaW5pc2hlZFwiLFxuXHR9KSxcblx0aU9TOiBPYmplY3QuZnJlZXplKHtcblx0XHRBcHBsaWNhdGlvbkRpZEJlY29tZUFjdGl2ZTogXCJpb3M6QXBwbGljYXRpb25EaWRCZWNvbWVBY3RpdmVcIixcblx0XHRBcHBsaWNhdGlvbkRpZEVudGVyQmFja2dyb3VuZDogXCJpb3M6QXBwbGljYXRpb25EaWRFbnRlckJhY2tncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbkRpZEZpbmlzaExhdW5jaGluZzogXCJpb3M6QXBwbGljYXRpb25EaWRGaW5pc2hMYXVuY2hpbmdcIixcblx0XHRBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nOiBcImlvczpBcHBsaWNhdGlvbkRpZFJlY2VpdmVNZW1vcnlXYXJuaW5nXCIsXG5cdFx0QXBwbGljYXRpb25XaWxsRW50ZXJGb3JlZ3JvdW5kOiBcImlvczpBcHBsaWNhdGlvbldpbGxFbnRlckZvcmVncm91bmRcIixcblx0XHRBcHBsaWNhdGlvbldpbGxSZXNpZ25BY3RpdmU6IFwiaW9zOkFwcGxpY2F0aW9uV2lsbFJlc2lnbkFjdGl2ZVwiLFxuXHRcdEFwcGxpY2F0aW9uV2lsbFRlcm1pbmF0ZTogXCJpb3M6QXBwbGljYXRpb25XaWxsVGVybWluYXRlXCIsXG5cdFx0V2luZG93RGlkTG9hZDogXCJpb3M6V2luZG93RGlkTG9hZFwiLFxuXHRcdFdpbmRvd1dpbGxBcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxBcHBlYXJcIixcblx0XHRXaW5kb3dEaWRBcHBlYXI6IFwiaW9zOldpbmRvd0RpZEFwcGVhclwiLFxuXHRcdFdpbmRvd1dpbGxEaXNhcHBlYXI6IFwiaW9zOldpbmRvd1dpbGxEaXNhcHBlYXJcIixcblx0XHRXaW5kb3dEaWREaXNhcHBlYXI6IFwiaW9zOldpbmRvd0RpZERpc2FwcGVhclwiLFxuXHRcdFdpbmRvd1NhZmVBcmVhSW5zZXRzQ2hhbmdlZDogXCJpb3M6V2luZG93U2FmZUFyZWFJbnNldHNDaGFuZ2VkXCIsXG5cdFx0V2luZG93T3JpZW50YXRpb25DaGFuZ2VkOiBcImlvczpXaW5kb3dPcmllbnRhdGlvbkNoYW5nZWRcIixcblx0XHRXaW5kb3dUb3VjaEJlZ2FuOiBcImlvczpXaW5kb3dUb3VjaEJlZ2FuXCIsXG5cdFx0V2luZG93VG91Y2hNb3ZlZDogXCJpb3M6V2luZG93VG91Y2hNb3ZlZFwiLFxuXHRcdFdpbmRvd1RvdWNoRW5kZWQ6IFwiaW9zOldpbmRvd1RvdWNoRW5kZWRcIixcblx0XHRXaW5kb3dUb3VjaENhbmNlbGxlZDogXCJpb3M6V2luZG93VG91Y2hDYW5jZWxsZWRcIixcblx0XHRXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkU3RhcnROYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZpbmlzaE5hdmlnYXRpb246IFwiaW9zOldlYlZpZXdEaWRGaW5pc2hOYXZpZ2F0aW9uXCIsXG5cdFx0V2ViVmlld0RpZEZhaWxOYXZpZ2F0aW9uOiBcImlvczpXZWJWaWV3RGlkRmFpbE5hdmlnYXRpb25cIixcblx0XHRXZWJWaWV3RGVjaWRlUG9saWN5Rm9yTmF2aWdhdGlvbkFjdGlvbjogXCJpb3M6V2ViVmlld0RlY2lkZVBvbGljeUZvck5hdmlnYXRpb25BY3Rpb25cIixcblx0fSksXG5cdENvbW1vbjogT2JqZWN0LmZyZWV6ZSh7XG5cdFx0QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZTogXCJjb21tb246QXBwbGljYXRpb25PcGVuZWRXaXRoRmlsZVwiLFxuXHRcdEFwcGxpY2F0aW9uU3RhcnRlZDogXCJjb21tb246QXBwbGljYXRpb25TdGFydGVkXCIsXG5cdFx0QXBwbGljYXRpb25MYXVuY2hlZFdpdGhVcmw6IFwiY29tbW9uOkFwcGxpY2F0aW9uTGF1bmNoZWRXaXRoVXJsXCIsXG5cdFx0VGhlbWVDaGFuZ2VkOiBcImNvbW1vbjpUaGVtZUNoYW5nZWRcIixcblx0XHRXaW5kb3dDbG9zaW5nOiBcImNvbW1vbjpXaW5kb3dDbG9zaW5nXCIsXG5cdFx0V2luZG93RGlkTW92ZTogXCJjb21tb246V2luZG93RGlkTW92ZVwiLFxuXHRcdFdpbmRvd0RpZFJlc2l6ZTogXCJjb21tb246V2luZG93RGlkUmVzaXplXCIsXG5cdFx0V2luZG93RFBJQ2hhbmdlZDogXCJjb21tb246V2luZG93RFBJQ2hhbmdlZFwiLFxuXHRcdFdpbmRvd0ZpbGVzRHJvcHBlZDogXCJjb21tb246V2luZG93RmlsZXNEcm9wcGVkXCIsXG5cdFx0V2luZG93Rm9jdXM6IFwiY29tbW9uOldpbmRvd0ZvY3VzXCIsXG5cdFx0V2luZG93RnVsbHNjcmVlbjogXCJjb21tb246V2luZG93RnVsbHNjcmVlblwiLFxuXHRcdFdpbmRvd0hpZGU6IFwiY29tbW9uOldpbmRvd0hpZGVcIixcblx0XHRXaW5kb3dMb3N0Rm9jdXM6IFwiY29tbW9uOldpbmRvd0xvc3RGb2N1c1wiLFxuXHRcdFdpbmRvd01heGltaXNlOiBcImNvbW1vbjpXaW5kb3dNYXhpbWlzZVwiLFxuXHRcdFdpbmRvd01pbmltaXNlOiBcImNvbW1vbjpXaW5kb3dNaW5pbWlzZVwiLFxuXHRcdFdpbmRvd1RvZ2dsZUZyYW1lbGVzczogXCJjb21tb246V2luZG93VG9nZ2xlRnJhbWVsZXNzXCIsXG5cdFx0V2luZG93UmVzdG9yZTogXCJjb21tb246V2luZG93UmVzdG9yZVwiLFxuXHRcdFdpbmRvd1J1bnRpbWVSZWFkeTogXCJjb21tb246V2luZG93UnVudGltZVJlYWR5XCIsXG5cdFx0V2luZG93U2hvdzogXCJjb21tb246V2luZG93U2hvd1wiLFxuXHRcdFdpbmRvd1VuRnVsbHNjcmVlbjogXCJjb21tb246V2luZG93VW5GdWxsc2NyZWVuXCIsXG5cdFx0V2luZG93VW5NYXhpbWlzZTogXCJjb21tb246V2luZG93VW5NYXhpbWlzZVwiLFxuXHRcdFdpbmRvd1VuTWluaW1pc2U6IFwiY29tbW9uOldpbmRvd1VuTWluaW1pc2VcIixcblx0XHRXaW5kb3dab29tOiBcImNvbW1vbjpXaW5kb3dab29tXCIsXG5cdFx0V2luZG93Wm9vbUluOiBcImNvbW1vbjpXaW5kb3dab29tSW5cIixcblx0XHRXaW5kb3dab29tT3V0OiBcImNvbW1vbjpXaW5kb3dab29tT3V0XCIsXG5cdFx0V2luZG93Wm9vbVJlc2V0OiBcImNvbW1vbjpXaW5kb3dab29tUmVzZXRcIixcblx0fSksXG59KTtcbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZSB0byB0aGUgY29uc29sZSB3aXRoIGN1c3RvbSBmb3JtYXR0aW5nLlxuICpcbiAqIEBwYXJhbSBtZXNzYWdlIC0gVGhlIG1lc3NhZ2UgdG8gYmUgbG9nZ2VkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gZGVidWdMb2cobWVzc2FnZTogYW55KSB7XG4gICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lXG4gICAgY29uc29sZS5sb2coXG4gICAgICAgICclYyB3YWlsczMgJWMgJyArIG1lc3NhZ2UgKyAnICcsXG4gICAgICAgICdiYWNrZ3JvdW5kOiAjYWEwMDAwOyBjb2xvcjogI2ZmZjsgYm9yZGVyLXJhZGl1czogM3B4IDBweCAwcHggM3B4OyBwYWRkaW5nOiAxcHg7IGZvbnQtc2l6ZTogMC43cmVtJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICMwMDk5MDA7IGNvbG9yOiAjZmZmOyBib3JkZXItcmFkaXVzOiAwcHggM3B4IDNweCAwcHg7IHBhZGRpbmc6IDFweDsgZm9udC1zaXplOiAwLjdyZW0nXG4gICAgKTtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgd2VidmlldyBzdXBwb3J0cyB0aGUge0BsaW5rIE1vdXNlRXZlbnQjYnV0dG9uc30gcHJvcGVydHkuXG4gKiBMb29raW5nIGF0IHlvdSBtYWNPUyBIaWdoIFNpZXJyYSFcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGNhblRyYWNrQnV0dG9ucygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKG5ldyBNb3VzZUV2ZW50KCdtb3VzZWRvd24nKSkuYnV0dG9ucyA9PT0gMDtcbn1cblxuLyoqXG4gKiBDaGVja3Mgd2hldGhlciB0aGUgYnJvd3NlciBzdXBwb3J0cyByZW1vdmluZyBsaXN0ZW5lcnMgYnkgdHJpZ2dlcmluZyBhbiBBYm9ydFNpZ25hbFxuICogKHNlZSBodHRwczovL2RldmVsb3Blci5tb3ppbGxhLm9yZy9lbi1VUy9kb2NzL1dlYi9BUEkvRXZlbnRUYXJnZXQvYWRkRXZlbnRMaXN0ZW5lciNzaWduYWwpLlxuICovXG5leHBvcnQgZnVuY3Rpb24gY2FuQWJvcnRMaXN0ZW5lcnMoKSB7XG4gICAgaWYgKCFFdmVudFRhcmdldCB8fCAhQWJvcnRTaWduYWwgfHwgIUFib3J0Q29udHJvbGxlcilcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuXG4gICAgbGV0IHJlc3VsdCA9IHRydWU7XG5cbiAgICBjb25zdCB0YXJnZXQgPSBuZXcgRXZlbnRUYXJnZXQoKTtcbiAgICBjb25zdCBjb250cm9sbGVyID0gbmV3IEFib3J0Q29udHJvbGxlcigpO1xuICAgIHRhcmdldC5hZGRFdmVudExpc3RlbmVyKCd0ZXN0JywgKCkgPT4geyByZXN1bHQgPSBmYWxzZTsgfSwgeyBzaWduYWw6IGNvbnRyb2xsZXIuc2lnbmFsIH0pO1xuICAgIGNvbnRyb2xsZXIuYWJvcnQoKTtcbiAgICB0YXJnZXQuZGlzcGF0Y2hFdmVudChuZXcgQ3VzdG9tRXZlbnQoJ3Rlc3QnKSk7XG5cbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG4vKipcbiAqIFJlc29sdmVzIHRoZSBjbG9zZXN0IEhUTUxFbGVtZW50IGFuY2VzdG9yIG9mIGFuIGV2ZW50J3MgdGFyZ2V0LlxuICovXG5leHBvcnQgZnVuY3Rpb24gZXZlbnRUYXJnZXQoZXZlbnQ6IEV2ZW50KTogSFRNTEVsZW1lbnQge1xuICAgIGlmIChldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkge1xuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0O1xuICAgIH0gZWxzZSBpZiAoIShldmVudC50YXJnZXQgaW5zdGFuY2VvZiBIVE1MRWxlbWVudCkgJiYgZXZlbnQudGFyZ2V0IGluc3RhbmNlb2YgTm9kZSkge1xuICAgICAgICByZXR1cm4gZXZlbnQudGFyZ2V0LnBhcmVudEVsZW1lbnQgPz8gZG9jdW1lbnQuYm9keTtcbiAgICB9IGVsc2Uge1xuICAgICAgICByZXR1cm4gZG9jdW1lbnQuYm9keTtcbiAgICB9XG59XG5cbi8qKipcbiBUaGlzIHRlY2huaXF1ZSBmb3IgcHJvcGVyIGxvYWQgZGV0ZWN0aW9uIGlzIHRha2VuIGZyb20gSFRNWDpcblxuIEJTRCAyLUNsYXVzZSBMaWNlbnNlXG5cbiBDb3B5cmlnaHQgKGMpIDIwMjAsIEJpZyBTa3kgU29mdHdhcmVcbiBBbGwgcmlnaHRzIHJlc2VydmVkLlxuXG4gUmVkaXN0cmlidXRpb24gYW5kIHVzZSBpbiBzb3VyY2UgYW5kIGJpbmFyeSBmb3Jtcywgd2l0aCBvciB3aXRob3V0XG4gbW9kaWZpY2F0aW9uLCBhcmUgcGVybWl0dGVkIHByb3ZpZGVkIHRoYXQgdGhlIGZvbGxvd2luZyBjb25kaXRpb25zIGFyZSBtZXQ6XG5cbiAxLiBSZWRpc3RyaWJ1dGlvbnMgb2Ygc291cmNlIGNvZGUgbXVzdCByZXRhaW4gdGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UsIHRoaXNcbiBsaXN0IG9mIGNvbmRpdGlvbnMgYW5kIHRoZSBmb2xsb3dpbmcgZGlzY2xhaW1lci5cblxuIDIuIFJlZGlzdHJpYnV0aW9ucyBpbiBiaW5hcnkgZm9ybSBtdXN0IHJlcHJvZHVjZSB0aGUgYWJvdmUgY29weXJpZ2h0IG5vdGljZSxcbiB0aGlzIGxpc3Qgb2YgY29uZGl0aW9ucyBhbmQgdGhlIGZvbGxvd2luZyBkaXNjbGFpbWVyIGluIHRoZSBkb2N1bWVudGF0aW9uXG4gYW5kL29yIG90aGVyIG1hdGVyaWFscyBwcm92aWRlZCB3aXRoIHRoZSBkaXN0cmlidXRpb24uXG5cbiBUSElTIFNPRlRXQVJFIElTIFBST1ZJREVEIEJZIFRIRSBDT1BZUklHSFQgSE9MREVSUyBBTkQgQ09OVFJJQlVUT1JTIFwiQVMgSVNcIlxuIEFORCBBTlkgRVhQUkVTUyBPUiBJTVBMSUVEIFdBUlJBTlRJRVMsIElOQ0xVRElORywgQlVUIE5PVCBMSU1JVEVEIFRPLCBUSEVcbiBJTVBMSUVEIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZIEFORCBGSVRORVNTIEZPUiBBIFBBUlRJQ1VMQVIgUFVSUE9TRSBBUkVcbiBESVNDTEFJTUVELiBJTiBOTyBFVkVOVCBTSEFMTCBUSEUgQ09QWVJJR0hUIEhPTERFUiBPUiBDT05UUklCVVRPUlMgQkUgTElBQkxFXG4gRk9SIEFOWSBESVJFQ1QsIElORElSRUNULCBJTkNJREVOVEFMLCBTUEVDSUFMLCBFWEVNUExBUlksIE9SIENPTlNFUVVFTlRJQUxcbiBEQU1BR0VTIChJTkNMVURJTkcsIEJVVCBOT1QgTElNSVRFRCBUTywgUFJPQ1VSRU1FTlQgT0YgU1VCU1RJVFVURSBHT09EUyBPUlxuIFNFUlZJQ0VTOyBMT1NTIE9GIFVTRSwgREFUQSwgT1IgUFJPRklUUzsgT1IgQlVTSU5FU1MgSU5URVJSVVBUSU9OKSBIT1dFVkVSXG4gQ0FVU0VEIEFORCBPTiBBTlkgVEhFT1JZIE9GIExJQUJJTElUWSwgV0hFVEhFUiBJTiBDT05UUkFDVCwgU1RSSUNUIExJQUJJTElUWSxcbiBPUiBUT1JUIChJTkNMVURJTkcgTkVHTElHRU5DRSBPUiBPVEhFUldJU0UpIEFSSVNJTkcgSU4gQU5ZIFdBWSBPVVQgT0YgVEhFIFVTRVxuIE9GIFRISVMgU09GVFdBUkUsIEVWRU4gSUYgQURWSVNFRCBPRiBUSEUgUE9TU0lCSUxJVFkgT0YgU1VDSCBEQU1BR0UuXG5cbiAqKiovXG5cbmxldCBpc1JlYWR5ID0gZmFsc2U7XG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgKCkgPT4geyBpc1JlYWR5ID0gdHJ1ZSB9KTtcblxuZXhwb3J0IGZ1bmN0aW9uIHdoZW5SZWFkeShjYWxsYmFjazogKCkgPT4gdm9pZCkge1xuICAgIGlmIChpc1JlYWR5IHx8IGRvY3VtZW50LnJlYWR5U3RhdGUgPT09ICdjb21wbGV0ZScpIHtcbiAgICAgICAgY2FsbGJhY2soKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBkb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKCdET01Db250ZW50TG9hZGVkJywgY2FsbGJhY2spO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHR5cGUgeyBTY3JlZW4gfSBmcm9tIFwiLi9zY3JlZW5zLmpzXCI7XG5cbi8vIERyb3AgdGFyZ2V0IGNvbnN0YW50c1xuY29uc3QgRFJPUF9UQVJHRVRfQVRUUklCVVRFID0gJ2RhdGEtZmlsZS1kcm9wLXRhcmdldCc7XG5jb25zdCBEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MgPSAnZmlsZS1kcm9wLXRhcmdldC1hY3RpdmUnO1xubGV0IGN1cnJlbnREcm9wVGFyZ2V0OiBFbGVtZW50IHwgbnVsbCA9IG51bGw7XG5cbmNvbnN0IFBvc2l0aW9uTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDA7XG5jb25zdCBDZW50ZXJNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxO1xuY29uc3QgQ2xvc2VNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gMjtcbmNvbnN0IERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICA9IDM7XG5jb25zdCBFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QgICAgICAgPSA0O1xuY29uc3QgRm9jdXNNZXRob2QgICAgICAgICAgICAgICAgICAgICAgID0gNTtcbmNvbnN0IEZvcmNlUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICA9IDY7XG5jb25zdCBGdWxsc2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgPSA3O1xuY29uc3QgR2V0U2NyZWVuTWV0aG9kICAgICAgICAgICAgICAgICAgID0gODtcbmNvbnN0IEdldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDk7XG5jb25zdCBIZWlnaHRNZXRob2QgICAgICAgICAgICAgICAgICAgICAgPSAxMDtcbmNvbnN0IEhpZGVNZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDExO1xuY29uc3QgSXNGb2N1c2VkTWV0aG9kICAgICAgICAgICAgICAgICAgID0gMTI7XG5jb25zdCBJc0Z1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSAxMztcbmNvbnN0IElzTWF4aW1pc2VkTWV0aG9kICAgICAgICAgICAgICAgICA9IDE0O1xuY29uc3QgSXNNaW5pbWlzZWRNZXRob2QgICAgICAgICAgICAgICAgID0gMTU7XG5jb25zdCBNYXhpbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAxNjtcbmNvbnN0IE1pbmltaXNlTWV0aG9kICAgICAgICAgICAgICAgICAgICA9IDE3O1xuY29uc3QgTmFtZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMTg7XG5jb25zdCBPcGVuRGV2VG9vbHNNZXRob2QgICAgICAgICAgICAgICAgPSAxOTtcbmNvbnN0IFJlbGF0aXZlUG9zaXRpb25NZXRob2QgICAgICAgICAgICA9IDIwO1xuY29uc3QgUmVsb2FkTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gMjE7XG5jb25zdCBSZXNpemFibGVNZXRob2QgICAgICAgICAgICAgICAgICAgPSAyMjtcbmNvbnN0IFJlc3RvcmVNZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDIzO1xuY29uc3QgU2V0UG9zaXRpb25NZXRob2QgICAgICAgICAgICAgICAgID0gMjQ7XG5jb25zdCBTZXRBbHdheXNPblRvcE1ldGhvZCAgICAgICAgICAgICAgPSAyNTtcbmNvbnN0IFNldEJhY2tncm91bmRDb2xvdXJNZXRob2QgICAgICAgICA9IDI2O1xuY29uc3QgU2V0RnJhbWVsZXNzTWV0aG9kICAgICAgICAgICAgICAgID0gMjc7XG5jb25zdCBTZXRGdWxsc2NyZWVuQnV0dG9uRW5hYmxlZE1ldGhvZCAgPSAyODtcbmNvbnN0IFNldE1heFNpemVNZXRob2QgICAgICAgICAgICAgICAgICA9IDI5O1xuY29uc3QgU2V0TWluU2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gMzA7XG5jb25zdCBTZXRSZWxhdGl2ZVBvc2l0aW9uTWV0aG9kICAgICAgICAgPSAzMTtcbmNvbnN0IFNldFJlc2l6YWJsZU1ldGhvZCAgICAgICAgICAgICAgICA9IDMyO1xuY29uc3QgU2V0U2l6ZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgID0gMzM7XG5jb25zdCBTZXRUaXRsZU1ldGhvZCAgICAgICAgICAgICAgICAgICAgPSAzNDtcbmNvbnN0IFNldFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICA9IDM1O1xuY29uc3QgU2hvd01ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgID0gMzY7XG5jb25zdCBTaXplTWV0aG9kICAgICAgICAgICAgICAgICAgICAgICAgPSAzNztcbmNvbnN0IFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICA9IDM4O1xuY29uc3QgVG9nZ2xlTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgID0gMzk7XG5jb25zdCBUb2dnbGVGcmFtZWxlc3NNZXRob2QgICAgICAgICAgICAgPSA0MDsgXG5jb25zdCBVbkZ1bGxzY3JlZW5NZXRob2QgICAgICAgICAgICAgICAgPSA0MTtcbmNvbnN0IFVuTWF4aW1pc2VNZXRob2QgICAgICAgICAgICAgICAgICA9IDQyO1xuY29uc3QgVW5NaW5pbWlzZU1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDM7XG5jb25zdCBXaWR0aE1ldGhvZCAgICAgICAgICAgICAgICAgICAgICAgPSA0NDtcbmNvbnN0IFpvb21NZXRob2QgICAgICAgICAgICAgICAgICAgICAgICA9IDQ1O1xuY29uc3QgWm9vbUluTWV0aG9kICAgICAgICAgICAgICAgICAgICAgID0gNDY7XG5jb25zdCBab29tT3V0TWV0aG9kICAgICAgICAgICAgICAgICAgICAgPSA0NztcbmNvbnN0IFpvb21SZXNldE1ldGhvZCAgICAgICAgICAgICAgICAgICA9IDQ4O1xuY29uc3QgU25hcEFzc2lzdE1ldGhvZCAgICAgICAgICAgICAgICAgID0gNDk7XG5jb25zdCBGaWxlc0Ryb3BwZWQgICAgICAgICAgICAgICAgICAgICAgPSA1MDtcbmNvbnN0IFByaW50TWV0aG9kICAgICAgICAgICAgICAgICAgICAgICA9IDUxO1xuXG4vKipcbiAqIEZpbmRzIHRoZSBuZWFyZXN0IGRyb3AgdGFyZ2V0IGVsZW1lbnQgYnkgd2Fsa2luZyB1cCB0aGUgRE9NIHRyZWUuXG4gKi9cbmZ1bmN0aW9uIGdldERyb3BUYXJnZXRFbGVtZW50KGVsZW1lbnQ6IEVsZW1lbnQgfCBudWxsKTogRWxlbWVudCB8IG51bGwge1xuICAgIGlmICghZWxlbWVudCkge1xuICAgICAgICByZXR1cm4gbnVsbDtcbiAgICB9XG4gICAgcmV0dXJuIGVsZW1lbnQuY2xvc2VzdChgWyR7RFJPUF9UQVJHRVRfQVRUUklCVVRFfV1gKTtcbn1cblxuLyoqXG4gKiBDaGVjayBpZiB3ZSBjYW4gdXNlIFdlYlZpZXcyJ3MgcG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMgKFdpbmRvd3MpXG4gKiBBbHNvIGNoZWNrcyB0aGF0IEVuYWJsZUZpbGVEcm9wIGlzIHRydWUgZm9yIHRoaXMgd2luZG93LlxuICovXG5mdW5jdGlvbiBjYW5SZXNvbHZlRmlsZVBhdGhzKCk6IGJvb2xlYW4ge1xuICAgIC8vIE11c3QgaGF2ZSBXZWJWaWV3MidzIHBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzIEFQSSAoV2luZG93cyBvbmx5KVxuICAgIGlmICgod2luZG93IGFzIGFueSkuY2hyb21lPy53ZWJ2aWV3Py5wb3N0TWVzc2FnZVdpdGhBZGRpdGlvbmFsT2JqZWN0cyA9PSBudWxsKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG4gICAgLy8gTXVzdCBoYXZlIEVuYWJsZUZpbGVEcm9wIHNldCB0byB0cnVlIGZvciB0aGlzIHdpbmRvd1xuICAgIC8vIFRoaXMgZmxhZyBpcyBzZXQgYnkgdGhlIEdvIGJhY2tlbmQgZHVyaW5nIHJ1bnRpbWUgaW5pdGlhbGl6YXRpb25cbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSB0cnVlO1xufVxuXG4vKipcbiAqIFNlbmQgZmlsZSBkcm9wIHRvIGJhY2tlbmQgdmlhIFdlYlZpZXcyIChXaW5kb3dzIG9ubHkpXG4gKi9cbmZ1bmN0aW9uIHJlc29sdmVGaWxlUGF0aHMoeDogbnVtYmVyLCB5OiBudW1iZXIsIGZpbGVzOiBGaWxlW10pOiB2b2lkIHtcbiAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2VXaXRoQWRkaXRpb25hbE9iamVjdHMpIHtcbiAgICAgICAgKHdpbmRvdyBhcyBhbnkpLmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlV2l0aEFkZGl0aW9uYWxPYmplY3RzKGBmaWxlOmRyb3A6JHt4fToke3l9YCwgZmlsZXMpO1xuICAgIH1cbn1cblxuLy8gTmF0aXZlIGRyYWcgc3RhdGUgKExpbnV4L21hY09TIGludGVyY2VwdCBET00gZHJhZyBldmVudHMpXG5sZXQgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xuXG4vKipcbiAqIENsZWFucyB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhbmQgaG92ZXIgZWZmZWN0cy5cbiAqIENhbGxlZCBvbiBkcm9wIG9yIHdoZW4gZHJhZyBsZWF2ZXMgdGhlIHdpbmRvdy5cbiAqL1xuZnVuY3Rpb24gY2xlYW51cE5hdGl2ZURyYWcoKTogdm9pZCB7XG4gICAgbmF0aXZlRHJhZ0FjdGl2ZSA9IGZhbHNlO1xuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCkge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICB9XG59XG5cbi8qKlxuICogQ2FsbGVkIGZyb20gR28gd2hlbiBhIGZpbGUgZHJhZyBlbnRlcnMgdGhlIHdpbmRvdyBvbiBMaW51eC9tYWNPUy5cbiAqL1xuZnVuY3Rpb24gaGFuZGxlRHJhZ0VudGVyKCk6IHZvaWQge1xuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBhY3RpdmF0ZSBkcmFnIHN0YXRlXG4gICAgfVxuICAgIG5hdGl2ZURyYWdBY3RpdmUgPSB0cnVlO1xufVxuXG4vKipcbiAqIENhbGxlZCBmcm9tIEdvIHdoZW4gYSBmaWxlIGRyYWcgbGVhdmVzIHRoZSB3aW5kb3cgb24gTGludXgvbWFjT1MuXG4gKi9cbmZ1bmN0aW9uIGhhbmRsZURyYWdMZWF2ZSgpOiB2b2lkIHtcbiAgICBjbGVhbnVwTmF0aXZlRHJhZygpO1xufVxuXG4vKipcbiAqIENhbGxlZCBmcm9tIEdvIGR1cmluZyBmaWxlIGRyYWcgdG8gdXBkYXRlIGhvdmVyIHN0YXRlIG9uIExpbnV4L21hY09TLlxuICogQHBhcmFtIHggLSBYIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xuICogQHBhcmFtIHkgLSBZIGNvb3JkaW5hdGUgaW4gQ1NTIHBpeGVsc1xuICovXG5mdW5jdGlvbiBoYW5kbGVEcmFnT3Zlcih4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgIGlmICghbmF0aXZlRHJhZ0FjdGl2ZSkgcmV0dXJuO1xuICAgIFxuICAgIC8vIENoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgaWYgKCh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmZsYWdzPy5lbmFibGVGaWxlRHJvcCA9PT0gZmFsc2UpIHtcbiAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBkb24ndCBzaG93IGhvdmVyIGVmZmVjdHNcbiAgICB9XG4gICAgXG4gICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgY29uc3QgZHJvcFRhcmdldCA9IGdldERyb3BUYXJnZXRFbGVtZW50KHRhcmdldEVsZW1lbnQpO1xuICAgIFxuICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgfVxuICAgIFxuICAgIGlmIChkcm9wVGFyZ2V0KSB7XG4gICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgIH1cbn1cblxuXG5cbi8vIEV4cG9ydCB0aGUgaGFuZGxlcnMgZm9yIHVzZSBieSBHbyB2aWEgaW5kZXgudHNcbmV4cG9ydCB7IGhhbmRsZURyYWdFbnRlciwgaGFuZGxlRHJhZ0xlYXZlLCBoYW5kbGVEcmFnT3ZlciB9O1xuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHBvc2l0aW9uIG9mIGEgd2luZG93LlxuICovXG5pbnRlcmZhY2UgUG9zaXRpb24ge1xuICAgIC8qKiBUaGUgaG9yaXpvbnRhbCBwb3NpdGlvbiBvZiB0aGUgd2luZG93LiAqL1xuICAgIHg6IG51bWJlcjtcbiAgICAvKiogVGhlIHZlcnRpY2FsIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuICovXG4gICAgeTogbnVtYmVyO1xufVxuXG4vKipcbiAqIEEgcmVjb3JkIGRlc2NyaWJpbmcgdGhlIHNpemUgb2YgYSB3aW5kb3cuXG4gKi9cbmludGVyZmFjZSBTaXplIHtcbiAgICAvKiogVGhlIHdpZHRoIG9mIHRoZSB3aW5kb3cuICovXG4gICAgd2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiB0aGUgd2luZG93LiAqL1xuICAgIGhlaWdodDogbnVtYmVyO1xufVxuXG4vLyBQcml2YXRlIGZpZWxkIG5hbWVzLlxuY29uc3QgY2FsbGVyU3ltID0gU3ltYm9sKFwiY2FsbGVyXCIpO1xuXG5jbGFzcyBXaW5kb3cge1xuICAgIC8vIFByaXZhdGUgZmllbGRzLlxuICAgIHByaXZhdGUgW2NhbGxlclN5bV06IChtZXNzYWdlOiBudW1iZXIsIGFyZ3M/OiBhbnkpID0+IFByb21pc2U8YW55PjtcblxuICAgIC8qKlxuICAgICAqIEluaXRpYWxpc2VzIGEgd2luZG93IG9iamVjdCB3aXRoIHRoZSBzcGVjaWZpZWQgbmFtZS5cbiAgICAgKlxuICAgICAqIEBwcml2YXRlXG4gICAgICogQHBhcmFtIG5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgdGFyZ2V0IHdpbmRvdy5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lOiBzdHJpbmcgPSAnJykge1xuICAgICAgICB0aGlzW2NhbGxlclN5bV0gPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLldpbmRvdywgbmFtZSlcblxuICAgICAgICAvLyBiaW5kIGluc3RhbmNlIG1ldGhvZCB0byBtYWtlIHRoZW0gZWFzaWx5IHVzYWJsZSBpbiBldmVudCBoYW5kbGVyc1xuICAgICAgICBmb3IgKGNvbnN0IG1ldGhvZCBvZiBPYmplY3QuZ2V0T3duUHJvcGVydHlOYW1lcyhXaW5kb3cucHJvdG90eXBlKSkge1xuICAgICAgICAgICAgaWYgKFxuICAgICAgICAgICAgICAgIG1ldGhvZCAhPT0gXCJjb25zdHJ1Y3RvclwiXG4gICAgICAgICAgICAgICAgJiYgdHlwZW9mICh0aGlzIGFzIGFueSlbbWV0aG9kXSA9PT0gXCJmdW5jdGlvblwiXG4gICAgICAgICAgICApIHtcbiAgICAgICAgICAgICAgICAodGhpcyBhcyBhbnkpW21ldGhvZF0gPSAodGhpcyBhcyBhbnkpW21ldGhvZF0uYmluZCh0aGlzKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEdldHMgdGhlIHNwZWNpZmllZCB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cgdG8gZ2V0LlxuICAgICAqIEByZXR1cm5zIFRoZSBjb3JyZXNwb25kaW5nIHdpbmRvdyBvYmplY3QuXG4gICAgICovXG4gICAgR2V0KG5hbWU6IHN0cmluZyk6IFdpbmRvdyB7XG4gICAgICAgIHJldHVybiBuZXcgV2luZG93KG5hbWUpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBhYnNvbHV0ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFBvc2l0aW9uKCk6IFByb21pc2U8UG9zaXRpb24+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2VudGVycyB0aGUgd2luZG93IG9uIHRoZSBzY3JlZW4uXG4gICAgICovXG4gICAgQ2VudGVyKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKENlbnRlck1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ2xvc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgQ2xvc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oQ2xvc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIERpc2FibGVzIG1pbi9tYXggc2l6ZSBjb25zdHJhaW50cy5cbiAgICAgKi9cbiAgICBEaXNhYmxlU2l6ZUNvbnN0cmFpbnRzKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKERpc2FibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEVuYWJsZXMgbWluL21heCBzaXplIGNvbnN0cmFpbnRzLlxuICAgICAqL1xuICAgIEVuYWJsZVNpemVDb25zdHJhaW50cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShFbmFibGVTaXplQ29uc3RyYWludHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEZvY3VzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBGb2N1cygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb2N1c01ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRm9yY2VzIHRoZSB3aW5kb3cgdG8gcmVsb2FkIHRoZSBwYWdlIGFzc2V0cy5cbiAgICAgKi9cbiAgICBGb3JjZVJlbG9hZCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGb3JjZVJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU3dpdGNoZXMgdGhlIHdpbmRvdyB0byBmdWxsc2NyZWVuIG1vZGUuXG4gICAgICovXG4gICAgRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShGdWxsc2NyZWVuTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSBzY3JlZW4gdGhhdCB0aGUgd2luZG93IGlzIG9uLlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBvbi5cbiAgICAgKi9cbiAgICBHZXRTY3JlZW4oKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShHZXRTY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIGN1cnJlbnQgem9vbSBsZXZlbCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgem9vbSBsZXZlbC5cbiAgICAgKi9cbiAgICBHZXRab29tKCk6IFByb21pc2U8bnVtYmVyPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oR2V0Wm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0aGUgaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBUaGUgY3VycmVudCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBIZWlnaHQoKTogUHJvbWlzZTxudW1iZXI+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIZWlnaHRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEhpZGVzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgSGlkZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShIaWRlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRydWUgaWYgdGhlIHdpbmRvdyBpcyBmb2N1c2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBmb2N1c2VkLlxuICAgICAqL1xuICAgIElzRm9jdXNlZCgpOiBQcm9taXNlPGJvb2xlYW4+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShJc0ZvY3VzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIGZ1bGxzY3JlZW4uXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBXaGV0aGVyIHRoZSB3aW5kb3cgaXMgY3VycmVudGx5IGZ1bGxzY3JlZW4uXG4gICAgICovXG4gICAgSXNGdWxsc2NyZWVuKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzRnVsbHNjcmVlbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgbWF4aW1pc2VkLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSBtYXhpbWlzZWQuXG4gICAgICovXG4gICAgSXNNYXhpbWlzZWQoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oSXNNYXhpbWlzZWRNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdHJ1ZSBpZiB0aGUgd2luZG93IGlzIG1pbmltaXNlZC5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFdoZXRoZXIgdGhlIHdpbmRvdyBpcyBjdXJyZW50bHkgbWluaW1pc2VkLlxuICAgICAqL1xuICAgIElzTWluaW1pc2VkKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKElzTWluaW1pc2VkTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBNYXhpbWlzZXMgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBNYXhpbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShNYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogTWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTWluaW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oTWluaW1pc2VNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIG5hbWUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBuYW1lIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgTmFtZSgpOiBQcm9taXNlPHN0cmluZz4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKE5hbWVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIE9wZW5zIHRoZSBkZXZlbG9wbWVudCB0b29scyBwYW5lLlxuICAgICAqL1xuICAgIE9wZW5EZXZUb29scygpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShPcGVuRGV2VG9vbHNNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgUmVsYXRpdmVQb3NpdGlvbigpOiBQcm9taXNlPFBvc2l0aW9uPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oUmVsYXRpdmVQb3NpdGlvbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVsb2FkcyB0aGUgcGFnZSBhc3NldHMuXG4gICAgICovXG4gICAgUmVsb2FkKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlbG9hZE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmV0dXJucyB0cnVlIGlmIHRoZSB3aW5kb3cgaXMgcmVzaXphYmxlLlxuICAgICAqXG4gICAgICogQHJldHVybnMgV2hldGhlciB0aGUgd2luZG93IGlzIGN1cnJlbnRseSByZXNpemFibGUuXG4gICAgICovXG4gICAgUmVzaXphYmxlKCk6IFByb21pc2U8Ym9vbGVhbj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc2l6YWJsZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVzdG9yZXMgdGhlIHdpbmRvdyB0byBpdHMgcHJldmlvdXMgc3RhdGUgaWYgaXQgd2FzIHByZXZpb3VzbHkgbWluaW1pc2VkLCBtYXhpbWlzZWQgb3IgZnVsbHNjcmVlbi5cbiAgICAgKi9cbiAgICBSZXN0b3JlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFJlc3RvcmVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0geCAtIFRoZSBkZXNpcmVkIGhvcml6b250YWwgYWJzb2x1dGUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0geSAtIFRoZSBkZXNpcmVkIHZlcnRpY2FsIGFic29sdXRlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0UG9zaXRpb24oeDogbnVtYmVyLCB5OiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxuICAgICAqXG4gICAgICogQHBhcmFtIGFsd2F5c09uVG9wIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBzdGF5IG9uIHRvcC5cbiAgICAgKi9cbiAgICBTZXRBbHdheXNPblRvcChhbHdheXNPblRvcDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEFsd2F5c09uVG9wTWV0aG9kLCB7IGFsd2F5c09uVG9wIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gciAtIFRoZSBkZXNpcmVkIHJlZCBjb21wb25lbnQgb2YgdGhlIHdpbmRvdyBiYWNrZ3JvdW5kLlxuICAgICAqIEBwYXJhbSBnIC0gVGhlIGRlc2lyZWQgZ3JlZW4gY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYiAtIFRoZSBkZXNpcmVkIGJsdWUgY29tcG9uZW50IG9mIHRoZSB3aW5kb3cgYmFja2dyb3VuZC5cbiAgICAgKiBAcGFyYW0gYSAtIFRoZSBkZXNpcmVkIGFscGhhIGNvbXBvbmVudCBvZiB0aGUgd2luZG93IGJhY2tncm91bmQuXG4gICAgICovXG4gICAgU2V0QmFja2dyb3VuZENvbG91cihyOiBudW1iZXIsIGc6IG51bWJlciwgYjogbnVtYmVyLCBhOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRCYWNrZ3JvdW5kQ29sb3VyTWV0aG9kLCB7IHIsIGcsIGIsIGEgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogUmVtb3ZlcyB0aGUgd2luZG93IGZyYW1lIGFuZCB0aXRsZSBiYXIuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gZnJhbWVsZXNzIC0gV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBmcmFtZWxlc3MuXG4gICAgICovXG4gICAgU2V0RnJhbWVsZXNzKGZyYW1lbGVzczogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZyYW1lbGVzc01ldGhvZCwgeyBmcmFtZWxlc3MgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGlzYWJsZXMgdGhlIHN5c3RlbSBmdWxsc2NyZWVuIGJ1dHRvbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBlbmFibGVkIC0gV2hldGhlciB0aGUgZnVsbHNjcmVlbiBidXR0b24gc2hvdWxkIGJlIGVuYWJsZWQuXG4gICAgICovXG4gICAgU2V0RnVsbHNjcmVlbkJ1dHRvbkVuYWJsZWQoZW5hYmxlZDogYm9vbGVhbik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldEZ1bGxzY3JlZW5CdXR0b25FbmFibGVkTWV0aG9kLCB7IGVuYWJsZWQgfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2V0cyB0aGUgbWF4aW11bSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCBtYXhpbXVtIHdpZHRoIG9mIHRoZSB3aW5kb3cuXG4gICAgICogQHBhcmFtIGhlaWdodCAtIFRoZSBkZXNpcmVkIG1heGltdW0gaGVpZ2h0IG9mIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgU2V0TWF4U2l6ZSh3aWR0aDogbnVtYmVyLCBoZWlnaHQ6IG51bWJlcik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNldE1heFNpemVNZXRob2QsIHsgd2lkdGgsIGhlaWdodCB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBtaW5pbXVtIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB3aWR0aCAtIFRoZSBkZXNpcmVkIG1pbmltdW0gd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKiBAcGFyYW0gaGVpZ2h0IC0gVGhlIGRlc2lyZWQgbWluaW11bSBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRNaW5TaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0TWluU2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHJlbGF0aXZlIHBvc2l0aW9uIG9mIHRoZSB3aW5kb3cgdG8gdGhlIHNjcmVlbi5cbiAgICAgKlxuICAgICAqIEBwYXJhbSB4IC0gVGhlIGRlc2lyZWQgaG9yaXpvbnRhbCByZWxhdGl2ZSBwb3NpdGlvbiBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIGRlc2lyZWQgdmVydGljYWwgcmVsYXRpdmUgcG9zaXRpb24gb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uKHg6IG51bWJlciwgeTogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVsYXRpdmVQb3NpdGlvbk1ldGhvZCwgeyB4LCB5IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgd2hldGhlciB0aGUgd2luZG93IGlzIHJlc2l6YWJsZS5cbiAgICAgKlxuICAgICAqIEBwYXJhbSByZXNpemFibGUgLSBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIHJlc2l6YWJsZS5cbiAgICAgKi9cbiAgICBTZXRSZXNpemFibGUocmVzaXphYmxlOiBib29sZWFuKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0UmVzaXphYmxlTWV0aG9kLCB7IHJlc2l6YWJsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSBzaXplIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gd2lkdGggLSBUaGUgZGVzaXJlZCB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqIEBwYXJhbSBoZWlnaHQgLSBUaGUgZGVzaXJlZCBoZWlnaHQgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTZXRTaXplKHdpZHRoOiBudW1iZXIsIGhlaWdodDogbnVtYmVyKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU2V0U2l6ZU1ldGhvZCwgeyB3aWR0aCwgaGVpZ2h0IH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFNldHMgdGhlIHRpdGxlIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gdGl0bGUgLSBUaGUgZGVzaXJlZCB0aXRsZSBvZiB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFNldFRpdGxlKHRpdGxlOiBzdHJpbmcpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRUaXRsZU1ldGhvZCwgeyB0aXRsZSB9KTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBTZXRzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3aW5kb3cuXG4gICAgICpcbiAgICAgKiBAcGFyYW0gem9vbSAtIFRoZSBkZXNpcmVkIHpvb20gbGV2ZWwuXG4gICAgICovXG4gICAgU2V0Wm9vbSh6b29tOiBudW1iZXIpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShTZXRab29tTWV0aG9kLCB7IHpvb20gfSk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogU2hvd3MgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNob3dNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFJldHVybnMgdGhlIHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoZSBjdXJyZW50IHNpemUgb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBTaXplKCk6IFByb21pc2U8U2l6ZT4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFNpemVNZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIGZ1bGxzY3JlZW4gYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGdWxsc2NyZWVuKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZUZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFRvZ2dsZXMgdGhlIHdpbmRvdyBiZXR3ZWVuIG1heGltaXNlZCBhbmQgbm9ybWFsLlxuICAgICAqL1xuICAgIFRvZ2dsZU1heGltaXNlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFRvZ2dsZU1heGltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUb2dnbGVzIHRoZSB3aW5kb3cgYmV0d2VlbiBmcmFtZWxlc3MgYW5kIG5vcm1hbC5cbiAgICAgKi9cbiAgICBUb2dnbGVGcmFtZWxlc3MoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVG9nZ2xlRnJhbWVsZXNzTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBVbi1mdWxsc2NyZWVucyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuRnVsbHNjcmVlbigpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbkZ1bGxzY3JlZW5NZXRob2QpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIFVuLW1heGltaXNlcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFVuTWF4aW1pc2UoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oVW5NYXhpbWlzZU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogVW4tbWluaW1pc2VzIHRoZSB3aW5kb3cuXG4gICAgICovXG4gICAgVW5NaW5pbWlzZSgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShVbk1pbmltaXNlTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXR1cm5zIHRoZSB3aWR0aCBvZiB0aGUgd2luZG93LlxuICAgICAqXG4gICAgICogQHJldHVybnMgVGhlIGN1cnJlbnQgd2lkdGggb2YgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBXaWR0aCgpOiBQcm9taXNlPG51bWJlcj4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFdpZHRoTWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBab29tcyB0aGUgd2luZG93LlxuICAgICAqL1xuICAgIFpvb20oKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbU1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogSW5jcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbUluKCk6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gdGhpc1tjYWxsZXJTeW1dKFpvb21Jbk1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogRGVjcmVhc2VzIHRoZSB6b29tIGxldmVsIG9mIHRoZSB3ZWJ2aWV3IGNvbnRlbnQuXG4gICAgICovXG4gICAgWm9vbU91dCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShab29tT3V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBSZXNldHMgdGhlIHpvb20gbGV2ZWwgb2YgdGhlIHdlYnZpZXcgY29udGVudC5cbiAgICAgKi9cbiAgICBab29tUmVzZXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oWm9vbVJlc2V0TWV0aG9kKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBIYW5kbGVzIGZpbGUgZHJvcHMgb3JpZ2luYXRpbmcgZnJvbSBwbGF0Zm9ybS1zcGVjaWZpYyBjb2RlIChlLmcuLCBtYWNPUy9MaW51eCBuYXRpdmUgZHJhZy1hbmQtZHJvcCkuXG4gICAgICogR2F0aGVycyBpbmZvcm1hdGlvbiBhYm91dCB0aGUgZHJvcCB0YXJnZXQgZWxlbWVudCBhbmQgc2VuZHMgaXQgYmFjayB0byB0aGUgR28gYmFja2VuZC5cbiAgICAgKlxuICAgICAqIEBwYXJhbSBmaWxlbmFtZXMgLSBBbiBhcnJheSBvZiBmaWxlIHBhdGhzIChzdHJpbmdzKSB0aGF0IHdlcmUgZHJvcHBlZC5cbiAgICAgKiBAcGFyYW0geCAtIFRoZSB4LWNvb3JkaW5hdGUgb2YgdGhlIGRyb3AgZXZlbnQgKENTUyBwaXhlbHMpLlxuICAgICAqIEBwYXJhbSB5IC0gVGhlIHktY29vcmRpbmF0ZSBvZiB0aGUgZHJvcCBldmVudCAoQ1NTIHBpeGVscykuXG4gICAgICovXG4gICAgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcChmaWxlbmFtZXM6IHN0cmluZ1tdLCB4OiBudW1iZXIsIHk6IG51bWJlcik6IHZvaWQge1xuICAgICAgICAvLyBDaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBGaWxlIGRyb3BzIGRpc2FibGVkLCBpZ25vcmUgdGhlIGRyb3BcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgY29uc3QgZWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoeCwgeSk7XG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudChlbGVtZW50KTtcblxuICAgICAgICBpZiAoIWRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIC8vIERyb3Agd2FzIG5vdCBvbiBhIGRlc2lnbmF0ZWQgZHJvcCB0YXJnZXQgLSBpZ25vcmVcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnN0IGVsZW1lbnREZXRhaWxzID0ge1xuICAgICAgICAgICAgaWQ6IGRyb3BUYXJnZXQuaWQsXG4gICAgICAgICAgICBjbGFzc0xpc3Q6IEFycmF5LmZyb20oZHJvcFRhcmdldC5jbGFzc0xpc3QpLFxuICAgICAgICAgICAgYXR0cmlidXRlczoge30gYXMgeyBba2V5OiBzdHJpbmddOiBzdHJpbmcgfSxcbiAgICAgICAgfTtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXMubGVuZ3RoOyBpKyspIHtcbiAgICAgICAgICAgIGNvbnN0IGF0dHIgPSBkcm9wVGFyZ2V0LmF0dHJpYnV0ZXNbaV07XG4gICAgICAgICAgICBlbGVtZW50RGV0YWlscy5hdHRyaWJ1dGVzW2F0dHIubmFtZV0gPSBhdHRyLnZhbHVlO1xuICAgICAgICB9XG5cbiAgICAgICAgY29uc3QgcGF5bG9hZCA9IHtcbiAgICAgICAgICAgIGZpbGVuYW1lcyxcbiAgICAgICAgICAgIHgsXG4gICAgICAgICAgICB5LFxuICAgICAgICAgICAgZWxlbWVudERldGFpbHMsXG4gICAgICAgIH07XG5cbiAgICAgICAgdGhpc1tjYWxsZXJTeW1dKEZpbGVzRHJvcHBlZCwgcGF5bG9hZCk7XG4gICAgICAgIFxuICAgICAgICAvLyBDbGVhbiB1cCBuYXRpdmUgZHJhZyBzdGF0ZSBhZnRlciBkcm9wXG4gICAgICAgIGNsZWFudXBOYXRpdmVEcmFnKCk7XG4gICAgfVxuICBcbiAgICAvKiBUcmlnZ2VycyBXaW5kb3dzIDExIFNuYXAgQXNzaXN0IGZlYXR1cmUgKFdpbmRvd3Mgb25seSkuXG4gICAgICogVGhpcyBpcyBlcXVpdmFsZW50IHRvIHByZXNzaW5nIFdpbitaIGFuZCBzaG93cyBzbmFwIGxheW91dCBvcHRpb25zLlxuICAgICAqL1xuICAgIFNuYXBBc3Npc3QoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiB0aGlzW2NhbGxlclN5bV0oU25hcEFzc2lzdE1ldGhvZCk7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogT3BlbnMgdGhlIHByaW50IGRpYWxvZyBmb3IgdGhlIHdpbmRvdy5cbiAgICAgKi9cbiAgICBQcmludCgpOiBQcm9taXNlPHZvaWQ+IHtcbiAgICAgICAgcmV0dXJuIHRoaXNbY2FsbGVyU3ltXShQcmludE1ldGhvZCk7XG4gICAgfVxufVxuXG4vKipcbiAqIFRoZSB3aW5kb3cgd2l0aGluIHdoaWNoIHRoZSBzY3JpcHQgaXMgcnVubmluZy5cbiAqL1xuY29uc3QgdGhpc1dpbmRvdyA9IG5ldyBXaW5kb3coJycpO1xuXG4vKipcbiAqIFNldHMgdXAgZ2xvYmFsIGRyYWcgYW5kIGRyb3AgZXZlbnQgbGlzdGVuZXJzIGZvciBmaWxlIGRyb3BzLlxuICogSGFuZGxlcyB2aXN1YWwgZmVlZGJhY2sgKGhvdmVyIHN0YXRlKSBhbmQgZmlsZSBkcm9wIHByb2Nlc3NpbmcuXG4gKi9cbmZ1bmN0aW9uIHNldHVwRHJvcFRhcmdldExpc3RlbmVycygpIHtcbiAgICBjb25zdCBkb2NFbGVtZW50ID0gZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50O1xuICAgIGxldCBkcmFnRW50ZXJDb3VudGVyID0gMDtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ2VudGVyJywgKGV2ZW50KSA9PiB7XG4gICAgICAgIGlmICghZXZlbnQuZGF0YVRyYW5zZmVyPy50eXBlcy5pbmNsdWRlcygnRmlsZXMnKSkge1xuICAgICAgICAgICAgcmV0dXJuOyAvLyBPbmx5IGhhbmRsZSBmaWxlIGRyYWdzLCBsZXQgb3RoZXIgZHJhZ3MgcGFzcyB0aHJvdWdoXG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnbm9uZSc7IC8vIFNob3cgXCJubyBkcm9wXCIgY3Vyc29yXG4gICAgICAgICAgICByZXR1cm47IC8vIEZpbGUgZHJvcHMgZGlzYWJsZWQsIGRvbid0IHNob3cgaG92ZXIgZWZmZWN0c1xuICAgICAgICB9XG4gICAgICAgIGRyYWdFbnRlckNvdW50ZXIrKztcbiAgICAgICAgXG4gICAgICAgIGNvbnN0IHRhcmdldEVsZW1lbnQgPSBkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFkpO1xuICAgICAgICBjb25zdCBkcm9wVGFyZ2V0ID0gZ2V0RHJvcFRhcmdldEVsZW1lbnQodGFyZ2V0RWxlbWVudCk7XG5cbiAgICAgICAgLy8gVXBkYXRlIGhvdmVyIHN0YXRlXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKGRyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGRyb3BUYXJnZXQuY2xhc3NMaXN0LmFkZChEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICAgICAgZXZlbnQuZGF0YVRyYW5zZmVyLmRyb3BFZmZlY3QgPSAnY29weSc7XG4gICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldCA9IGRyb3BUYXJnZXQ7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJztcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgfVxuICAgIH0sIGZhbHNlKTtcblxuICAgIGRvY0VsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcignZHJhZ292ZXInLCAoZXZlbnQpID0+IHtcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJhZ3NcbiAgICAgICAgfVxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdub25lJzsgLy8gU2hvdyBcIm5vIGRyb3BcIiBjdXJzb3JcbiAgICAgICAgICAgIHJldHVybjsgLy8gRmlsZSBkcm9wcyBkaXNhYmxlZCwgZG9uJ3Qgc2hvdyBob3ZlciBlZmZlY3RzXG4gICAgICAgIH1cbiAgICAgICAgXG4gICAgICAgIC8vIFVwZGF0ZSBkcm9wIHRhcmdldCBhcyBjdXJzb3IgbW92ZXNcbiAgICAgICAgY29uc3QgdGFyZ2V0RWxlbWVudCA9IGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQoZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSk7XG4gICAgICAgIGNvbnN0IGRyb3BUYXJnZXQgPSBnZXREcm9wVGFyZ2V0RWxlbWVudCh0YXJnZXRFbGVtZW50KTtcbiAgICAgICAgXG4gICAgICAgIGlmIChjdXJyZW50RHJvcFRhcmdldCAmJiBjdXJyZW50RHJvcFRhcmdldCAhPT0gZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQuY2xhc3NMaXN0LnJlbW92ZShEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpO1xuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICBpZiAoZHJvcFRhcmdldCkge1xuICAgICAgICAgICAgaWYgKCFkcm9wVGFyZ2V0LmNsYXNzTGlzdC5jb250YWlucyhEUk9QX1RBUkdFVF9BQ1RJVkVfQ0xBU1MpKSB7XG4gICAgICAgICAgICAgICAgZHJvcFRhcmdldC5jbGFzc0xpc3QuYWRkKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBldmVudC5kYXRhVHJhbnNmZXIuZHJvcEVmZmVjdCA9ICdjb3B5JztcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gZHJvcFRhcmdldDtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIGV2ZW50LmRhdGFUcmFuc2Zlci5kcm9wRWZmZWN0ID0gJ25vbmUnO1xuICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgICAgICB9XG4gICAgfSwgZmFsc2UpO1xuXG4gICAgZG9jRWxlbWVudC5hZGRFdmVudExpc3RlbmVyKCdkcmFnbGVhdmUnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTsgLy8gQWx3YXlzIHByZXZlbnQgZGVmYXVsdCB0byBzdG9wIGJyb3dzZXIgbmF2aWdhdGlvblxuICAgICAgICAvLyBPbiBXaW5kb3dzLCBjaGVjayBpZiBmaWxlIGRyb3BzIGFyZSBlbmFibGVkIGZvciB0aGlzIHdpbmRvd1xuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZmxhZ3M/LmVuYWJsZUZpbGVEcm9wID09PSBmYWxzZSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgICAgIFxuICAgICAgICAvLyBPbiBMaW51eC9XZWJLaXRHVEsgYW5kIG1hY09TLCBkcmFnbGVhdmUgZmlyZXMgaW1tZWRpYXRlbHkgd2l0aCByZWxhdGVkVGFyZ2V0PW51bGwgd2hlbiBuYXRpdmVcbiAgICAgICAgLy8gZHJhZyBoYW5kbGluZyBpcyBpbnZvbHZlZC4gSWdub3JlIHRoZXNlIHNwdXJpb3VzIGV2ZW50cyAtIHdlJ2xsIGNsZWFuIHVwIG9uIGRyb3AgaW5zdGVhZC5cbiAgICAgICAgaWYgKGV2ZW50LnJlbGF0ZWRUYXJnZXQgPT09IG51bGwpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBcbiAgICAgICAgZHJhZ0VudGVyQ291bnRlci0tO1xuICAgICAgICBcbiAgICAgICAgaWYgKGRyYWdFbnRlckNvdW50ZXIgPT09IDAgfHwgXG4gICAgICAgICAgICAoY3VycmVudERyb3BUYXJnZXQgJiYgIWN1cnJlbnREcm9wVGFyZ2V0LmNvbnRhaW5zKGV2ZW50LnJlbGF0ZWRUYXJnZXQgYXMgTm9kZSkpKSB7XG4gICAgICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgICAgICBjdXJyZW50RHJvcFRhcmdldC5jbGFzc0xpc3QucmVtb3ZlKERST1BfVEFSR0VUX0FDVElWRV9DTEFTUyk7XG4gICAgICAgICAgICAgICAgY3VycmVudERyb3BUYXJnZXQgPSBudWxsO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG5cbiAgICBkb2NFbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ2Ryb3AnLCAoZXZlbnQpID0+IHtcbiAgICAgICAgaWYgKCFldmVudC5kYXRhVHJhbnNmZXI/LnR5cGVzLmluY2x1ZGVzKCdGaWxlcycpKSB7XG4gICAgICAgICAgICByZXR1cm47IC8vIE9ubHkgaGFuZGxlIGZpbGUgZHJvcHNcbiAgICAgICAgfVxuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBBbHdheXMgcHJldmVudCBkZWZhdWx0IHRvIHN0b3AgYnJvd3NlciBuYXZpZ2F0aW9uXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGNoZWNrIGlmIGZpbGUgZHJvcHMgYXJlIGVuYWJsZWQgZm9yIHRoaXMgd2luZG93XG4gICAgICAgIGlmICgod2luZG93IGFzIGFueSkuX3dhaWxzPy5mbGFncz8uZW5hYmxlRmlsZURyb3AgPT09IGZhbHNlKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgZHJhZ0VudGVyQ291bnRlciA9IDA7XG4gICAgICAgIFxuICAgICAgICBpZiAoY3VycmVudERyb3BUYXJnZXQpIHtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0LmNsYXNzTGlzdC5yZW1vdmUoRFJPUF9UQVJHRVRfQUNUSVZFX0NMQVNTKTtcbiAgICAgICAgICAgIGN1cnJlbnREcm9wVGFyZ2V0ID0gbnVsbDtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIE9uIFdpbmRvd3MsIGhhbmRsZSBmaWxlIGRyb3BzIHZpYSBKYXZhU2NyaXB0XG4gICAgICAgIC8vIE9uIG1hY09TL0xpbnV4LCBuYXRpdmUgY29kZSB3aWxsIGNhbGwgSGFuZGxlUGxhdGZvcm1GaWxlRHJvcFxuICAgICAgICBpZiAoY2FuUmVzb2x2ZUZpbGVQYXRocygpKSB7XG4gICAgICAgICAgICBjb25zdCBmaWxlczogRmlsZVtdID0gW107XG4gICAgICAgICAgICBpZiAoZXZlbnQuZGF0YVRyYW5zZmVyLml0ZW1zKSB7XG4gICAgICAgICAgICAgICAgZm9yIChjb25zdCBpdGVtIG9mIGV2ZW50LmRhdGFUcmFuc2Zlci5pdGVtcykge1xuICAgICAgICAgICAgICAgICAgICBpZiAoaXRlbS5raW5kID09PSAnZmlsZScpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IGZpbGUgPSBpdGVtLmdldEFzRmlsZSgpO1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGZpbGUpIGZpbGVzLnB1c2goZmlsZSk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9IGVsc2UgaWYgKGV2ZW50LmRhdGFUcmFuc2Zlci5maWxlcykge1xuICAgICAgICAgICAgICAgIGZvciAoY29uc3QgZmlsZSBvZiBldmVudC5kYXRhVHJhbnNmZXIuZmlsZXMpIHtcbiAgICAgICAgICAgICAgICAgICAgZmlsZXMucHVzaChmaWxlKTtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBcbiAgICAgICAgICAgIGlmIChmaWxlcy5sZW5ndGggPiAwKSB7XG4gICAgICAgICAgICAgICAgcmVzb2x2ZUZpbGVQYXRocyhldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBmaWxlcyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9LCBmYWxzZSk7XG59XG5cbi8vIEluaXRpYWxpemUgbGlzdGVuZXJzIHdoZW4gdGhlIHNjcmlwdCBsb2Fkc1xuaWYgKHR5cGVvZiB3aW5kb3cgIT09IFwidW5kZWZpbmVkXCIgJiYgdHlwZW9mIGRvY3VtZW50ICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgc2V0dXBEcm9wVGFyZ2V0TGlzdGVuZXJzKCk7XG59XG5cbmV4cG9ydCBkZWZhdWx0IHRoaXNXaW5kb3c7XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCAqIGFzIFJ1bnRpbWUgZnJvbSBcIi4uL0B3YWlsc2lvL3J1bnRpbWUvc3JjXCI7XG5cbi8vIE5PVEU6IHRoZSBmb2xsb3dpbmcgbWV0aG9kcyBNVVNUIGJlIGltcG9ydGVkIGV4cGxpY2l0bHkgYmVjYXVzZSBvZiBob3cgZXNidWlsZCBpbmplY3Rpb24gd29ya3NcbmltcG9ydCB7IEVuYWJsZSBhcyBFbmFibGVXTUwgfSBmcm9tIFwiLi4vQHdhaWxzaW8vcnVudGltZS9zcmMvd21sXCI7XG5pbXBvcnQgeyBkZWJ1Z0xvZyB9IGZyb20gXCIuLi9Ad2FpbHNpby9ydW50aW1lL3NyYy91dGlsc1wiO1xuXG53aW5kb3cud2FpbHMgPSBSdW50aW1lO1xuRW5hYmxlV01MKCk7XG5cbmlmIChERUJVRykge1xuICAgIGRlYnVnTG9nKFwiV2FpbHMgUnVudGltZSBMb2FkZWRcIilcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLlN5c3RlbSk7XG5cbmNvbnN0IFN5c3RlbUlzRGFya01vZGUgPSAwO1xuY29uc3QgU3lzdGVtRW52aXJvbm1lbnQgPSAxO1xuY29uc3QgU3lzdGVtQ2FwYWJpbGl0aWVzID0gMjtcblxuY29uc3QgX2ludm9rZSA9IChmdW5jdGlvbiAoKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgLy8gV2luZG93cyBXZWJWaWV3MlxuICAgICAgICBpZiAoKHdpbmRvdyBhcyBhbnkpLmNocm9tZT8ud2Vidmlldz8ucG9zdE1lc3NhZ2UpIHtcbiAgICAgICAgICAgIHJldHVybiAod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkuY2hyb21lLndlYnZpZXcpO1xuICAgICAgICB9XG4gICAgICAgIC8vIG1hY09TL2lPUyBXS1dlYlZpZXdcbiAgICAgICAgZWxzZSBpZiAoKHdpbmRvdyBhcyBhbnkpLndlYmtpdD8ubWVzc2FnZUhhbmRsZXJzPy5bJ2V4dGVybmFsJ10/LnBvc3RNZXNzYWdlKSB7XG4gICAgICAgICAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLndlYmtpdC5tZXNzYWdlSGFuZGxlcnNbJ2V4dGVybmFsJ10ucG9zdE1lc3NhZ2UuYmluZCgod2luZG93IGFzIGFueSkud2Via2l0Lm1lc3NhZ2VIYW5kbGVyc1snZXh0ZXJuYWwnXSk7XG4gICAgICAgIH1cbiAgICAgICAgLy8gQW5kcm9pZCBXZWJWaWV3IC0gdXNlcyBhZGRKYXZhc2NyaXB0SW50ZXJmYWNlIHdoaWNoIGV4cG9zZXMgd2luZG93LndhaWxzLmludm9rZVxuICAgICAgICBlbHNlIGlmICgod2luZG93IGFzIGFueSkud2FpbHM/Lmludm9rZSkge1xuICAgICAgICAgICAgcmV0dXJuIChtc2c6IGFueSkgPT4gKHdpbmRvdyBhcyBhbnkpLndhaWxzLmludm9rZSh0eXBlb2YgbXNnID09PSAnc3RyaW5nJyA/IG1zZyA6IEpTT04uc3RyaW5naWZ5KG1zZykpO1xuICAgICAgICB9XG4gICAgfSBjYXRjaChlKSB7fVxuXG4gICAgY29uc29sZS53YXJuKCdcXG4lY1x1MjZBMFx1RkUwRiBCcm93c2VyIEVudmlyb25tZW50IERldGVjdGVkICVjXFxuXFxuJWNPbmx5IFVJIHByZXZpZXdzIGFyZSBhdmFpbGFibGUgaW4gdGhlIGJyb3dzZXIuIEZvciBmdWxsIGZ1bmN0aW9uYWxpdHksIHBsZWFzZSBydW4gdGhlIGFwcGxpY2F0aW9uIGluIGRlc2t0b3AgbW9kZS5cXG5Nb3JlIGluZm9ybWF0aW9uIGF0OiBodHRwczovL3YzLndhaWxzLmlvL2xlYXJuL2J1aWxkLyN1c2luZy1hLWJyb3dzZXItZm9yLWRldmVsb3BtZW50XFxuJyxcbiAgICAgICAgJ2JhY2tncm91bmQ6ICNmZmZmZmY7IGNvbG9yOiAjMDAwMDAwOyBmb250LXdlaWdodDogYm9sZDsgcGFkZGluZzogNHB4IDhweDsgYm9yZGVyLXJhZGl1czogNHB4OyBib3JkZXI6IDJweCBzb2xpZCAjMDAwMDAwOycsXG4gICAgICAgICdiYWNrZ3JvdW5kOiB0cmFuc3BhcmVudDsnLFxuICAgICAgICAnY29sb3I6ICNmZmZmZmY7IGZvbnQtc3R5bGU6IGl0YWxpYzsgZm9udC13ZWlnaHQ6IGJvbGQ7Jyk7XG4gICAgcmV0dXJuIG51bGw7XG59KSgpO1xuXG5leHBvcnQgZnVuY3Rpb24gaW52b2tlKG1zZzogYW55KTogdm9pZCB7XG4gICAgX2ludm9rZT8uKG1zZyk7XG59XG5cbi8qKlxuICogUmV0cmlldmVzIHRoZSBzeXN0ZW0gZGFyayBtb2RlIHN0YXR1cy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKTogUHJvbWlzZTxib29sZWFuPiB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtSXNEYXJrTW9kZSk7XG59XG5cbi8qKlxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgdGhlIGNhcGFiaWxpdGllcy5cbiAqL1xuZXhwb3J0IGFzeW5jIGZ1bmN0aW9uIENhcGFiaWxpdGllcygpOiBQcm9taXNlPFJlY29yZDxzdHJpbmcsIGFueT4+IHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1DYXBhYmlsaXRpZXMpO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIE9TSW5mbyB7XG4gICAgLyoqIFRoZSBicmFuZGluZyBvZiB0aGUgT1MuICovXG4gICAgQnJhbmRpbmc6IHN0cmluZztcbiAgICAvKiogVGhlIElEIG9mIHRoZSBPUy4gKi9cbiAgICBJRDogc3RyaW5nO1xuICAgIC8qKiBUaGUgbmFtZSBvZiB0aGUgT1MuICovXG4gICAgTmFtZTogc3RyaW5nO1xuICAgIC8qKiBUaGUgdmVyc2lvbiBvZiB0aGUgT1MuICovXG4gICAgVmVyc2lvbjogc3RyaW5nO1xufVxuXG5leHBvcnQgaW50ZXJmYWNlIEVudmlyb25tZW50SW5mbyB7XG4gICAgLyoqIFRoZSBhcmNoaXRlY3R1cmUgb2YgdGhlIHN5c3RlbS4gKi9cbiAgICBBcmNoOiBzdHJpbmc7XG4gICAgLyoqIFRydWUgaWYgdGhlIGFwcGxpY2F0aW9uIGlzIHJ1bm5pbmcgaW4gZGVidWcgbW9kZSwgb3RoZXJ3aXNlIGZhbHNlLiAqL1xuICAgIERlYnVnOiBib29sZWFuO1xuICAgIC8qKiBUaGUgb3BlcmF0aW5nIHN5c3RlbSBpbiB1c2UuICovXG4gICAgT1M6IHN0cmluZztcbiAgICAvKiogRGV0YWlscyBvZiB0aGUgb3BlcmF0aW5nIHN5c3RlbS4gKi9cbiAgICBPU0luZm86IE9TSW5mbztcbiAgICAvKiogQWRkaXRpb25hbCBwbGF0Zm9ybSBpbmZvcm1hdGlvbi4gKi9cbiAgICBQbGF0Zm9ybUluZm86IFJlY29yZDxzdHJpbmcsIGFueT47XG59XG5cbi8qKlxuICogUmV0cmlldmVzIGVudmlyb25tZW50IGRldGFpbHMuXG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbnZpcm9ubWVudCgpOiBQcm9taXNlPEVudmlyb25tZW50SW5mbz4ge1xuICAgIHJldHVybiBjYWxsKFN5c3RlbUVudmlyb25tZW50KTtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLlxuICpcbiAqIEByZXR1cm4gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJc1dpbmRvd3MoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJ3aW5kb3dzXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXG4gKlxuICogQHJldHVybnMgUmV0dXJucyB0cnVlIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXgsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUyA9PT0gXCJsaW51eFwiO1xufVxuXG4vKipcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBhIG1hY09TIG9wZXJhdGluZyBzeXN0ZW0uXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgbWFjT1MsIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzTWFjKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uT1MgPT09IFwiZGFyd2luXCI7XG59XG5cbi8qKlxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NC5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNBTUQ2NCgpOiBib29sZWFuIHtcbiAgICByZXR1cm4gKHdpbmRvdyBhcyBhbnkpLl93YWlscz8uZW52aXJvbm1lbnQ/LkFyY2ggPT09IFwiYW1kNjRcIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTS5cbiAqXG4gKiBAcmV0dXJucyBUcnVlIGlmIHRoZSBjdXJyZW50IGFyY2hpdGVjdHVyZSBpcyBBUk0sIGZhbHNlIG90aGVyd2lzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCk6IGJvb2xlYW4ge1xuICAgIHJldHVybiAod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uQXJjaCA9PT0gXCJhcm1cIjtcbn1cblxuLyoqXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLlxuICpcbiAqIEByZXR1cm5zIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzQVJNNjQoKTogYm9vbGVhbiB7XG4gICAgcmV0dXJuICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5BcmNoID09PSBcImFybTY0XCI7XG59XG5cbi8qKlxuICogUmVwb3J0cyB3aGV0aGVyIHRoZSBhcHAgaXMgYmVpbmcgcnVuIGluIGRlYnVnIG1vZGUuXG4gKlxuICogQHJldHVybnMgVHJ1ZSBpZiB0aGUgYXBwIGlzIGJlaW5nIHJ1biBpbiBkZWJ1ZyBtb2RlLlxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNEZWJ1ZygpOiBib29sZWFuIHtcbiAgICByZXR1cm4gQm9vbGVhbigod2luZG93IGFzIGFueSkuX3dhaWxzPy5lbnZpcm9ubWVudD8uRGVidWcpO1xufVxuXG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgSXNEZWJ1ZyB9IGZyb20gXCIuL3N5c3RlbS5qc1wiO1xuaW1wb3J0IHsgZXZlbnRUYXJnZXQgfSBmcm9tIFwiLi91dGlscy5qc1wiO1xuXG4vLyBzZXR1cFxud2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcblxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xuXG5jb25zdCBDb250ZXh0TWVudU9wZW4gPSAwO1xuXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQ6IHN0cmluZywgeDogbnVtYmVyLCB5OiBudW1iZXIsIGRhdGE6IGFueSk6IHZvaWQge1xuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQ6IE1vdXNlRXZlbnQpIHtcbiAgICBjb25zdCB0YXJnZXQgPSBldmVudFRhcmdldChldmVudCk7XG5cbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUodGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIikudHJpbSgpO1xuXG4gICAgaWYgKGN1c3RvbUNvbnRleHRNZW51KSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIGNvbnN0IGRhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGRhdGEpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQsIHRhcmdldCk7XG4gICAgfVxufVxuXG5cbi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQ6IE1vdXNlRXZlbnQsIHRhcmdldDogSFRNTEVsZW1lbnQpIHtcbiAgICAvLyBEZWJ1ZyBidWlsZHMgYWx3YXlzIHNob3cgdGhlIG1lbnVcbiAgICBpZiAoSXNEZWJ1ZygpKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgc3dpdGNoICh3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZSh0YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpKSB7XG4gICAgICAgIGNhc2UgJ3Nob3cnOlxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBjYXNlICdoaWRlJzpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcbiAgICBpZiAodGFyZ2V0LmlzQ29udGVudEVkaXRhYmxlKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXG4gICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IHNlbGVjdGlvbiAmJiBzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwO1xuICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcbiAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xuICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gdGFyZ2V0KSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBDaGVjayBpZiB0YWcgaXMgaW5wdXQgb3IgdGV4dGFyZWEuXG4gICAgaWYgKHRhcmdldCBpbnN0YW5jZW9mIEhUTUxJbnB1dEVsZW1lbnQgfHwgdGFyZ2V0IGluc3RhbmNlb2YgSFRNTFRleHRBcmVhRWxlbWVudCkge1xuICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghdGFyZ2V0LnJlYWRPbmx5ICYmICF0YXJnZXQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoqXG4gKiBSZXRyaWV2ZXMgdGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleSBmcm9tIHRoZSBmbGFnIG1hcC5cbiAqXG4gKiBAcGFyYW0ga2V5IC0gVGhlIGtleSB0byByZXRyaWV2ZSB0aGUgdmFsdWUgZm9yLlxuICogQHJldHVybiBUaGUgdmFsdWUgYXNzb2NpYXRlZCB3aXRoIHRoZSBzcGVjaWZpZWQga2V5LlxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXk6IHN0cmluZyk6IGFueSB7XG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIHdpbmRvdy5fd2FpbHMuZmxhZ3Nba2V5XTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHRocm93IG5ldyBFcnJvcihcIlVuYWJsZSB0byByZXRyaWV2ZSBmbGFnICdcIiArIGtleSArIFwiJzogXCIgKyBlLCB7IGNhdXNlOiBlIH0pO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgaW52b2tlLCBJc1dpbmRvd3MgfSBmcm9tIFwiLi9zeXN0ZW0uanNcIjtcbmltcG9ydCB7IEdldEZsYWcgfSBmcm9tIFwiLi9mbGFncy5qc1wiO1xuaW1wb3J0IHsgY2FuVHJhY2tCdXR0b25zLCBldmVudFRhcmdldCB9IGZyb20gXCIuL3V0aWxzLmpzXCI7XG5cbi8vIFNldHVwXG5sZXQgY2FuRHJhZyA9IGZhbHNlO1xubGV0IGRyYWdnaW5nID0gZmFsc2U7XG5cbmxldCByZXNpemFibGUgPSBmYWxzZTtcbmxldCBjYW5SZXNpemUgPSBmYWxzZTtcbmxldCByZXNpemluZyA9IGZhbHNlO1xubGV0IHJlc2l6ZUVkZ2U6IHN0cmluZyA9IFwiXCI7XG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xuXG5sZXQgYnV0dG9ucyA9IDA7XG5jb25zdCBidXR0b25zVHJhY2tlZCA9IGNhblRyYWNrQnV0dG9ucygpO1xuXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcbndpbmRvdy5fd2FpbHMuc2V0UmVzaXphYmxlID0gKHZhbHVlOiBib29sZWFuKTogdm9pZCA9PiB7XG4gICAgcmVzaXphYmxlID0gdmFsdWU7XG4gICAgaWYgKCFyZXNpemFibGUpIHtcbiAgICAgICAgLy8gU3RvcCByZXNpemluZyBpZiBpbiBwcm9ncmVzcy5cbiAgICAgICAgY2FuUmVzaXplID0gcmVzaXppbmcgPSBmYWxzZTtcbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxufTtcblxuLy8gRGVmZXIgYXR0YWNoaW5nIG1vdXNlIGxpc3RlbmVycyB1bnRpbCB3ZSBrbm93IHdlJ3JlIG5vdCBvbiBtb2JpbGUuXG5sZXQgZHJhZ0luaXREb25lID0gZmFsc2U7XG5mdW5jdGlvbiBpc01vYmlsZSgpOiBib29sZWFuIHtcbiAgICBjb25zdCBvcyA9ICh3aW5kb3cgYXMgYW55KS5fd2FpbHM/LmVudmlyb25tZW50Py5PUztcbiAgICBpZiAob3MgPT09IFwiaW9zXCIgfHwgb3MgPT09IFwiYW5kcm9pZFwiKSByZXR1cm4gdHJ1ZTtcbiAgICAvLyBGYWxsYmFjayBoZXVyaXN0aWMgaWYgZW52aXJvbm1lbnQgbm90IHlldCBzZXRcbiAgICBjb25zdCB1YSA9IG5hdmlnYXRvci51c2VyQWdlbnQgfHwgbmF2aWdhdG9yLnZlbmRvciB8fCAod2luZG93IGFzIGFueSkub3BlcmEgfHwgXCJcIjtcbiAgICByZXR1cm4gL2FuZHJvaWR8aXBob25lfGlwYWR8aXBvZHxpZW1vYmlsZXx3cGRlc2t0b3AvaS50ZXN0KHVhKTtcbn1cbmZ1bmN0aW9uIHRyeUluaXREcmFnSGFuZGxlcnMoKTogdm9pZCB7XG4gICAgaWYgKGRyYWdJbml0RG9uZSkgcmV0dXJuO1xuICAgIGlmIChpc01vYmlsZSgpKSByZXR1cm47XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCB1cGRhdGUsIHsgY2FwdHVyZTogdHJ1ZSB9KTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIHVwZGF0ZSwgeyBjYXB0dXJlOiB0cnVlIH0pO1xuICAgIGZvciAoY29uc3QgZXYgb2YgWydjbGljaycsICdjb250ZXh0bWVudScsICdkYmxjbGljayddKSB7XG4gICAgICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKGV2LCBzdXBwcmVzc0V2ZW50LCB7IGNhcHR1cmU6IHRydWUgfSk7XG4gICAgfVxuICAgIGRyYWdJbml0RG9uZSA9IHRydWU7XG59XG4vLyBBdHRlbXB0IGltbWVkaWF0ZSBpbml0IChpbiBjYXNlIGVudmlyb25tZW50IGFscmVhZHkgcHJlc2VudClcbnRyeUluaXREcmFnSGFuZGxlcnMoKTtcbi8vIEFsc28gYXR0ZW1wdCBvbiBET00gcmVhZHlcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoJ0RPTUNvbnRlbnRMb2FkZWQnLCB0cnlJbml0RHJhZ0hhbmRsZXJzLCB7IG9uY2U6IHRydWUgfSk7XG4vLyBBcyBhIGxhc3QgcmVzb3J0LCBwb2xsIGZvciBlbnZpcm9ubWVudCBmb3IgYSBzaG9ydCBwZXJpb2RcbmxldCBkcmFnRW52UG9sbHMgPSAwO1xuY29uc3QgZHJhZ0VudlBvbGwgPSB3aW5kb3cuc2V0SW50ZXJ2YWwoKCkgPT4ge1xuICAgIGlmIChkcmFnSW5pdERvbmUpIHsgd2luZG93LmNsZWFySW50ZXJ2YWwoZHJhZ0VudlBvbGwpOyByZXR1cm47IH1cbiAgICB0cnlJbml0RHJhZ0hhbmRsZXJzKCk7XG4gICAgaWYgKCsrZHJhZ0VudlBvbGxzID4gMTAwKSB7IHdpbmRvdy5jbGVhckludGVydmFsKGRyYWdFbnZQb2xsKTsgfVxufSwgNTApO1xuXG5mdW5jdGlvbiBzdXBwcmVzc0V2ZW50KGV2ZW50OiBFdmVudCkge1xuICAgIC8vIFN1cHByZXNzIGNsaWNrIGV2ZW50cyB3aGlsZSByZXNpemluZyBvciBkcmFnZ2luZy5cbiAgICBpZiAoZHJhZ2dpbmcgfHwgcmVzaXppbmcpIHtcbiAgICAgICAgZXZlbnQuc3RvcEltbWVkaWF0ZVByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnN0b3BQcm9wYWdhdGlvbigpO1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn1cblxuLy8gVXNlIGNvbnN0YW50cyB0byBhdm9pZCBjb21wYXJpbmcgc3RyaW5ncyBtdWx0aXBsZSB0aW1lcy5cbmNvbnN0IE1vdXNlRG93biA9IDA7XG5jb25zdCBNb3VzZVVwICAgPSAxO1xuY29uc3QgTW91c2VNb3ZlID0gMjtcblxuZnVuY3Rpb24gdXBkYXRlKGV2ZW50OiBNb3VzZUV2ZW50KSB7XG4gICAgLy8gV2luZG93cyBzdXBwcmVzc2VzIG1vdXNlIGV2ZW50cyBhdCB0aGUgZW5kIG9mIGRyYWdnaW5nIG9yIHJlc2l6aW5nLFxuICAgIC8vIHNvIHdlIG5lZWQgdG8gYmUgc21hcnQgYW5kIHN5bnRoZXNpemUgYnV0dG9uIGV2ZW50cy5cblxuICAgIGxldCBldmVudFR5cGU6IG51bWJlciwgZXZlbnRCdXR0b25zID0gZXZlbnQuYnV0dG9ucztcbiAgICBzd2l0Y2ggKGV2ZW50LnR5cGUpIHtcbiAgICAgICAgY2FzZSAnbW91c2Vkb3duJzpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlRG93bjtcbiAgICAgICAgICAgIGlmICghYnV0dG9uc1RyYWNrZWQpIHsgZXZlbnRCdXR0b25zID0gYnV0dG9ucyB8ICgxIDw8IGV2ZW50LmJ1dHRvbik7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICBjYXNlICdtb3VzZXVwJzpcbiAgICAgICAgICAgIGV2ZW50VHlwZSA9IE1vdXNlVXA7XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnMgJiB+KDEgPDwgZXZlbnQuYnV0dG9uKTsgfVxuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICBldmVudFR5cGUgPSBNb3VzZU1vdmU7XG4gICAgICAgICAgICBpZiAoIWJ1dHRvbnNUcmFja2VkKSB7IGV2ZW50QnV0dG9ucyA9IGJ1dHRvbnM7IH1cbiAgICAgICAgICAgIGJyZWFrO1xuICAgIH1cblxuICAgIGxldCByZWxlYXNlZCA9IGJ1dHRvbnMgJiB+ZXZlbnRCdXR0b25zO1xuICAgIGxldCBwcmVzc2VkID0gZXZlbnRCdXR0b25zICYgfmJ1dHRvbnM7XG5cbiAgICBidXR0b25zID0gZXZlbnRCdXR0b25zO1xuXG4gICAgLy8gU3ludGhlc2l6ZSBhIHJlbGVhc2UtcHJlc3Mgc2VxdWVuY2UgaWYgd2UgZGV0ZWN0IGEgcHJlc3Mgb2YgYW4gYWxyZWFkeSBwcmVzc2VkIGJ1dHRvbi5cbiAgICBpZiAoZXZlbnRUeXBlID09PSBNb3VzZURvd24gJiYgIShwcmVzc2VkICYgZXZlbnQuYnV0dG9uKSkge1xuICAgICAgICByZWxlYXNlZCB8PSAoMSA8PCBldmVudC5idXR0b24pO1xuICAgICAgICBwcmVzc2VkIHw9ICgxIDw8IGV2ZW50LmJ1dHRvbik7XG4gICAgfVxuXG4gICAgLy8gU3VwcHJlc3MgYWxsIGJ1dHRvbiBldmVudHMgZHVyaW5nIGRyYWdnaW5nIGFuZCByZXNpemluZyxcbiAgICAvLyB1bmxlc3MgdGhpcyBpcyBhIG1vdXNldXAgZXZlbnQgdGhhdCBpcyBlbmRpbmcgYSBkcmFnIGFjdGlvbi5cbiAgICBpZiAoXG4gICAgICAgIGV2ZW50VHlwZSAhPT0gTW91c2VNb3ZlIC8vIEZhc3QgcGF0aCBmb3IgbW91c2Vtb3ZlXG4gICAgICAgICYmIHJlc2l6aW5nXG4gICAgICAgIHx8IChcbiAgICAgICAgICAgIGRyYWdnaW5nXG4gICAgICAgICAgICAmJiAoXG4gICAgICAgICAgICAgICAgZXZlbnRUeXBlID09PSBNb3VzZURvd25cbiAgICAgICAgICAgICAgICB8fCBldmVudC5idXR0b24gIT09IDBcbiAgICAgICAgICAgIClcbiAgICAgICAgKVxuICAgICkge1xuICAgICAgICBldmVudC5zdG9wSW1tZWRpYXRlUHJvcGFnYXRpb24oKTtcbiAgICAgICAgZXZlbnQuc3RvcFByb3BhZ2F0aW9uKCk7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxuXG4gICAgLy8gSGFuZGxlIHJlbGVhc2VzXG4gICAgaWYgKHJlbGVhc2VkICYgMSkgeyBwcmltYXJ5VXAoZXZlbnQpOyB9XG4gICAgLy8gSGFuZGxlIHByZXNzZXNcbiAgICBpZiAocHJlc3NlZCAmIDEpIHsgcHJpbWFyeURvd24oZXZlbnQpOyB9XG5cbiAgICAvLyBIYW5kbGUgbW91c2Vtb3ZlXG4gICAgaWYgKGV2ZW50VHlwZSA9PT0gTW91c2VNb3ZlKSB7IG9uTW91c2VNb3ZlKGV2ZW50KTsgfTtcbn1cblxuZnVuY3Rpb24gcHJpbWFyeURvd24oZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcbiAgICAvLyBSZXNldCByZWFkaW5lc3Mgc3RhdGUuXG4gICAgY2FuRHJhZyA9IGZhbHNlO1xuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xuXG4gICAgLy8gSWdub3JlIHJlcGVhdGVkIGNsaWNrcyBvbiBtYWNPUyBhbmQgTGludXguXG4gICAgaWYgKCFJc1dpbmRvd3MoKSkge1xuICAgICAgICBpZiAoZXZlbnQudHlwZSA9PT0gJ21vdXNlZG93bicgJiYgZXZlbnQuYnV0dG9uID09PSAwICYmIGV2ZW50LmRldGFpbCAhPT0gMSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgaWYgKHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgLy8gUmVhZHkgdG8gcmVzaXplIGlmIHRoZSBwcmltYXJ5IGJ1dHRvbiB3YXMgcHJlc3NlZCBmb3IgdGhlIGZpcnN0IHRpbWUuXG4gICAgICAgIGNhblJlc2l6ZSA9IHRydWU7XG4gICAgICAgIC8vIERvIG5vdCBzdGFydCBkcmFnIG9wZXJhdGlvbnMgd2hlbiBvbiByZXNpemUgZWRnZXMuXG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBSZXRyaWV2ZSB0YXJnZXQgZWxlbWVudFxuICAgIGNvbnN0IHRhcmdldCA9IGV2ZW50VGFyZ2V0KGV2ZW50KTtcblxuICAgIC8vIFJlYWR5IHRvIGRyYWcgaWYgdGhlIHByaW1hcnkgYnV0dG9uIHdhcyBwcmVzc2VkIGZvciB0aGUgZmlyc3QgdGltZSBvbiBhIGRyYWdnYWJsZSBlbGVtZW50LlxuICAgIC8vIElnbm9yZSBjbGlja3Mgb24gdGhlIHNjcm9sbGJhci5cbiAgICBjb25zdCBzdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKHRhcmdldCk7XG4gICAgY2FuRHJhZyA9IChcbiAgICAgICAgc3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2FpbHMtZHJhZ2dhYmxlXCIpLnRyaW0oKSA9PT0gXCJkcmFnXCJcbiAgICAgICAgJiYgKFxuICAgICAgICAgICAgZXZlbnQub2Zmc2V0WCAtIHBhcnNlRmxvYXQoc3R5bGUucGFkZGluZ0xlZnQpIDwgdGFyZ2V0LmNsaWVudFdpZHRoXG4gICAgICAgICAgICAmJiBldmVudC5vZmZzZXRZIC0gcGFyc2VGbG9hdChzdHlsZS5wYWRkaW5nVG9wKSA8IHRhcmdldC5jbGllbnRIZWlnaHRcbiAgICAgICAgKVxuICAgICk7XG59XG5cbmZ1bmN0aW9uIHByaW1hcnlVcChldmVudDogTW91c2VFdmVudCkge1xuICAgIC8vIFN0b3AgZHJhZ2dpbmcgYW5kIHJlc2l6aW5nLlxuICAgIGNhbkRyYWcgPSBmYWxzZTtcbiAgICBkcmFnZ2luZyA9IGZhbHNlO1xuICAgIGNhblJlc2l6ZSA9IGZhbHNlO1xuICAgIHJlc2l6aW5nID0gZmFsc2U7XG59XG5cbmNvbnN0IGN1cnNvckZvckVkZ2UgPSBPYmplY3QuZnJlZXplKHtcbiAgICBcInNlLXJlc2l6ZVwiOiBcIm53c2UtcmVzaXplXCIsXG4gICAgXCJzdy1yZXNpemVcIjogXCJuZXN3LXJlc2l6ZVwiLFxuICAgIFwibnctcmVzaXplXCI6IFwibndzZS1yZXNpemVcIixcbiAgICBcIm5lLXJlc2l6ZVwiOiBcIm5lc3ctcmVzaXplXCIsXG4gICAgXCJ3LXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxuICAgIFwibi1yZXNpemVcIjogXCJucy1yZXNpemVcIixcbiAgICBcInMtcmVzaXplXCI6IFwibnMtcmVzaXplXCIsXG4gICAgXCJlLXJlc2l6ZVwiOiBcImV3LXJlc2l6ZVwiLFxufSlcblxuZnVuY3Rpb24gc2V0UmVzaXplKGVkZ2U/OiBrZXlvZiB0eXBlb2YgY3Vyc29yRm9yRWRnZSk6IHZvaWQge1xuICAgIGlmIChlZGdlKSB7XG4gICAgICAgIGlmICghcmVzaXplRWRnZSkgeyBkZWZhdWx0Q3Vyc29yID0gZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3I7IH1cbiAgICAgICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSBjdXJzb3JGb3JFZGdlW2VkZ2VdO1xuICAgIH0gZWxzZSBpZiAoIWVkZ2UgJiYgcmVzaXplRWRnZSkge1xuICAgICAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9IGRlZmF1bHRDdXJzb3I7XG4gICAgfVxuXG4gICAgcmVzaXplRWRnZSA9IGVkZ2UgfHwgXCJcIjtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZU1vdmUoZXZlbnQ6IE1vdXNlRXZlbnQpOiB2b2lkIHtcbiAgICBpZiAoY2FuUmVzaXplICYmIHJlc2l6ZUVkZ2UpIHtcbiAgICAgICAgLy8gU3RhcnQgcmVzaXppbmcuXG4gICAgICAgIHJlc2l6aW5nID0gdHJ1ZTtcbiAgICAgICAgaW52b2tlKFwid2FpbHM6cmVzaXplOlwiICsgcmVzaXplRWRnZSk7XG4gICAgfSBlbHNlIGlmIChjYW5EcmFnKSB7XG4gICAgICAgIC8vIFN0YXJ0IGRyYWdnaW5nLlxuICAgICAgICBkcmFnZ2luZyA9IHRydWU7XG4gICAgICAgIGludm9rZShcIndhaWxzOmRyYWdcIik7XG4gICAgfVxuXG4gICAgaWYgKGRyYWdnaW5nIHx8IHJlc2l6aW5nKSB7XG4gICAgICAgIC8vIEVpdGhlciBkcmFnIG9yIHJlc2l6ZSBpcyBvbmdvaW5nLFxuICAgICAgICAvLyByZXNldCByZWFkaW5lc3MgYW5kIHN0b3AgcHJvY2Vzc2luZy5cbiAgICAgICAgY2FuRHJhZyA9IGNhblJlc2l6ZSA9IGZhbHNlO1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKCFyZXNpemFibGUgfHwgIUlzV2luZG93cygpKSB7XG4gICAgICAgIGlmIChyZXNpemVFZGdlKSB7IHNldFJlc2l6ZSgpOyB9XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBjb25zdCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGNvbnN0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xuXG4gICAgLy8gRXh0cmEgcGl4ZWxzIGZvciB0aGUgY29ybmVyIGFyZWFzLlxuICAgIGNvbnN0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xuXG4gICAgY29uc3QgcmlnaHRCb3JkZXIgPSAod2luZG93Lm91dGVyV2lkdGggLSBldmVudC5jbGllbnRYKSA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGNvbnN0IGxlZnRCb3JkZXIgPSBldmVudC5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgY29uc3QgdG9wQm9yZGVyID0gZXZlbnQuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcbiAgICBjb25zdCBib3R0b21Cb3JkZXIgPSAod2luZG93Lm91dGVySGVpZ2h0IC0gZXZlbnQuY2xpZW50WSkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG5cbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lciBhcmVhcy5cbiAgICBjb25zdCByaWdodENvcm5lciA9ICh3aW5kb3cub3V0ZXJXaWR0aCAtIGV2ZW50LmNsaWVudFgpIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IGxlZnRDb3JuZXIgPSBldmVudC5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IHRvcENvcm5lciA9IGV2ZW50LmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuICAgIGNvbnN0IGJvdHRvbUNvcm5lciA9ICh3aW5kb3cub3V0ZXJIZWlnaHQgLSBldmVudC5jbGllbnRZKSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XG5cbiAgICBpZiAoIWxlZnRDb3JuZXIgJiYgIXRvcENvcm5lciAmJiAhYm90dG9tQ29ybmVyICYmICFyaWdodENvcm5lcikge1xuICAgICAgICAvLyBPcHRpbWlzYXRpb246IG91dCBvZiBhbGwgY29ybmVyIGFyZWFzIGltcGxpZXMgb3V0IG9mIGJvcmRlcnMuXG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbiAgICAvLyBEZXRlY3QgY29ybmVycy5cbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgLy8gRGV0ZWN0IGJvcmRlcnMuXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcbiAgICAvLyBPdXQgb2YgYm9yZGVyIGFyZWEuXG4gICAgZWxzZSBzZXRSZXNpemUoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5BcHBsaWNhdGlvbik7XG5cbmNvbnN0IEhpZGVNZXRob2QgPSAwO1xuY29uc3QgU2hvd01ldGhvZCA9IDE7XG5jb25zdCBRdWl0TWV0aG9kID0gMjtcblxuLyoqXG4gKiBIaWRlcyBhIGNlcnRhaW4gbWV0aG9kIGJ5IGNhbGxpbmcgdGhlIEhpZGVNZXRob2QgZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKEhpZGVNZXRob2QpO1xufVxuXG4vKipcbiAqIENhbGxzIHRoZSBTaG93TWV0aG9kIGFuZCByZXR1cm5zIHRoZSByZXN1bHQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTaG93KCk6IFByb21pc2U8dm9pZD4ge1xuICAgIHJldHVybiBjYWxsKFNob3dNZXRob2QpO1xufVxuXG4vKipcbiAqIENhbGxzIHRoZSBRdWl0TWV0aG9kIHRvIHRlcm1pbmF0ZSB0aGUgcHJvZ3JhbS5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoUXVpdE1ldGhvZCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmltcG9ydCB7IENhbmNlbGxhYmxlUHJvbWlzZSwgdHlwZSBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzIH0gZnJvbSBcIi4vY2FuY2VsbGFibGUuanNcIjtcbmltcG9ydCB7IG5ld1J1bnRpbWVDYWxsZXIsIG9iamVjdE5hbWVzIH0gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSBcIi4vbmFub2lkLmpzXCI7XG5cbi8vIFNldHVwXG53aW5kb3cuX3dhaWxzID0gd2luZG93Ll93YWlscyB8fCB7fTtcblxudHlwZSBQcm9taXNlUmVzb2x2ZXJzID0gT21pdDxDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPGFueT4sIFwicHJvbWlzZVwiIHwgXCJvbmNhbmNlbGxlZFwiPlxuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DYWxsKTtcbmNvbnN0IGNhbmNlbENhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLkNhbmNlbENhbGwpO1xuY29uc3QgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXA8c3RyaW5nLCBQcm9taXNlUmVzb2x2ZXJzPigpO1xuXG5jb25zdCBDYWxsQmluZGluZyA9IDA7XG5jb25zdCBDYW5jZWxNZXRob2QgPSAwXG5cbi8qKlxuICogSG9sZHMgYWxsIHJlcXVpcmVkIGluZm9ybWF0aW9uIGZvciBhIGJpbmRpbmcgY2FsbC5cbiAqIE1heSBwcm92aWRlIGVpdGhlciBhIG1ldGhvZCBJRCBvciBhIG1ldGhvZCBuYW1lLCBidXQgbm90IGJvdGguXG4gKi9cbmV4cG9ydCB0eXBlIENhbGxPcHRpb25zID0ge1xuICAgIC8qKiBUaGUgbnVtZXJpYyBJRCBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kSUQ6IG51bWJlcjtcbiAgICAvKiogVGhlIGZ1bGx5IHF1YWxpZmllZCBuYW1lIG9mIHRoZSBib3VuZCBtZXRob2QgdG8gY2FsbC4gKi9cbiAgICBtZXRob2ROYW1lPzogbmV2ZXI7XG4gICAgLyoqIEFyZ3VtZW50cyB0byBiZSBwYXNzZWQgaW50byB0aGUgYm91bmQgbWV0aG9kLiAqL1xuICAgIGFyZ3M6IGFueVtdO1xufSB8IHtcbiAgICAvKiogVGhlIG51bWVyaWMgSUQgb2YgdGhlIGJvdW5kIG1ldGhvZCB0byBjYWxsLiAqL1xuICAgIG1ldGhvZElEPzogbmV2ZXI7XG4gICAgLyoqIFRoZSBmdWxseSBxdWFsaWZpZWQgbmFtZSBvZiB0aGUgYm91bmQgbWV0aG9kIHRvIGNhbGwuICovXG4gICAgbWV0aG9kTmFtZTogc3RyaW5nO1xuICAgIC8qKiBBcmd1bWVudHMgdG8gYmUgcGFzc2VkIGludG8gdGhlIGJvdW5kIG1ldGhvZC4gKi9cbiAgICBhcmdzOiBhbnlbXTtcbn07XG5cbi8qKlxuICogRXhjZXB0aW9uIGNsYXNzIHRoYXQgd2lsbCBiZSB0aHJvd24gaW4gY2FzZSB0aGUgYm91bmQgbWV0aG9kIHJldHVybnMgYW4gZXJyb3IuXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBSdW50aW1lRXJyb3IjbmFtZX0gcHJvcGVydHkgaXMgXCJSdW50aW1lRXJyb3JcIi5cbiAqL1xuZXhwb3J0IGNsYXNzIFJ1bnRpbWVFcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IFJ1bnRpbWVFcnJvciBpbnN0YW5jZS5cbiAgICAgKiBAcGFyYW0gbWVzc2FnZSAtIFRoZSBlcnJvciBtZXNzYWdlLlxuICAgICAqIEBwYXJhbSBvcHRpb25zIC0gT3B0aW9ucyB0byBiZSBmb3J3YXJkZWQgdG8gdGhlIEVycm9yIGNvbnN0cnVjdG9yLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKG1lc3NhZ2U/OiBzdHJpbmcsIG9wdGlvbnM/OiBFcnJvck9wdGlvbnMpIHtcbiAgICAgICAgc3VwZXIobWVzc2FnZSwgb3B0aW9ucyk7XG4gICAgICAgIHRoaXMubmFtZSA9IFwiUnVudGltZUVycm9yXCI7XG4gICAgfVxufVxuXG4vKipcbiAqIEdlbmVyYXRlcyBhIHVuaXF1ZSBJRCB1c2luZyB0aGUgbmFub2lkIGxpYnJhcnkuXG4gKlxuICogQHJldHVybnMgQSB1bmlxdWUgSUQgdGhhdCBkb2VzIG5vdCBleGlzdCBpbiB0aGUgY2FsbFJlc3BvbnNlcyBzZXQuXG4gKi9cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKTogc3RyaW5nIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoY2FsbFJlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuLyoqXG4gKiBDYWxsIGEgYm91bmQgbWV0aG9kIGFjY29yZGluZyB0byB0aGUgZ2l2ZW4gY2FsbCBvcHRpb25zLlxuICpcbiAqIEluIGNhc2Ugb2YgZmFpbHVyZSwgdGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCByZWplY3Qgd2l0aCBhbiBleGNlcHRpb25cbiAqIGFtb25nIFJlZmVyZW5jZUVycm9yICh1bmtub3duIG1ldGhvZCksIFR5cGVFcnJvciAod3JvbmcgYXJndW1lbnQgY291bnQgb3IgdHlwZSksXG4gKiB7QGxpbmsgUnVudGltZUVycm9yfSAobWV0aG9kIHJldHVybmVkIGFuIGVycm9yKSwgb3Igb3RoZXIgKG5ldHdvcmsgb3IgaW50ZXJuYWwgZXJyb3JzKS5cbiAqIFRoZSBleGNlcHRpb24gbWlnaHQgaGF2ZSBhIFwiY2F1c2VcIiBmaWVsZCB3aXRoIHRoZSB2YWx1ZSByZXR1cm5lZFxuICogYnkgdGhlIGFwcGxpY2F0aW9uLSBvciBzZXJ2aWNlLWxldmVsIGVycm9yIG1hcnNoYWxpbmcgZnVuY3Rpb25zLlxuICpcbiAqIEBwYXJhbSBvcHRpb25zIC0gQSBtZXRob2QgY2FsbCBkZXNjcmlwdG9yLlxuICogQHJldHVybnMgVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9uczogQ2FsbE9wdGlvbnMpOiBDYW5jZWxsYWJsZVByb21pc2U8YW55PiB7XG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XG5cbiAgICBjb25zdCByZXN1bHQgPSBDYW5jZWxsYWJsZVByb21pc2Uud2l0aFJlc29sdmVyczxhbnk+KCk7XG4gICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZTogcmVzdWx0LnJlc29sdmUsIHJlamVjdDogcmVzdWx0LnJlamVjdCB9KTtcblxuICAgIGNvbnN0IHJlcXVlc3QgPSBjYWxsKENhbGxCaW5kaW5nLCBPYmplY3QuYXNzaWduKHsgXCJjYWxsLWlkXCI6IGlkIH0sIG9wdGlvbnMpKTtcbiAgICBsZXQgcnVubmluZyA9IHRydWU7XG5cbiAgICByZXF1ZXN0LnRoZW4oKHJlcykgPT4ge1xuICAgICAgICBydW5uaW5nID0gZmFsc2U7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgcmVzdWx0LnJlc29sdmUocmVzKTtcbiAgICB9LCAoZXJyKSA9PiB7XG4gICAgICAgIHJ1bm5pbmcgPSBmYWxzZTtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICByZXN1bHQucmVqZWN0KGVycik7XG4gICAgfSk7XG5cbiAgICBjb25zdCBjYW5jZWwgPSAoKSA9PiB7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgcmV0dXJuIGNhbmNlbENhbGwoQ2FuY2VsTWV0aG9kLCB7XCJjYWxsLWlkXCI6IGlkfSkuY2F0Y2goKGVycikgPT4ge1xuICAgICAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIHdoaWxlIHJlcXVlc3RpbmcgYmluZGluZyBjYWxsIGNhbmNlbGxhdGlvbjpcIiwgZXJyKTtcbiAgICAgICAgfSk7XG4gICAgfTtcblxuICAgIHJlc3VsdC5vbmNhbmNlbGxlZCA9ICgpID0+IHtcbiAgICAgICAgaWYgKHJ1bm5pbmcpIHtcbiAgICAgICAgICAgIHJldHVybiBjYW5jZWwoKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHJldHVybiByZXF1ZXN0LnRoZW4oY2FuY2VsKTtcbiAgICAgICAgfVxuICAgIH07XG5cbiAgICByZXR1cm4gcmVzdWx0LnByb21pc2U7XG59XG5cbi8qKlxuICogQ2FsbHMgYSBib3VuZCBtZXRob2QgYnkgbmFtZSB3aXRoIHRoZSBzcGVjaWZpZWQgYXJndW1lbnRzLlxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cbiAqXG4gKiBAcGFyYW0gbWV0aG9kTmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJy5cbiAqIEBwYXJhbSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXG4gKiBAcmV0dXJucyBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEJ5TmFtZShtZXRob2ROYW1lOiBzdHJpbmcsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kTmFtZSwgYXJncyB9KTtcbn1cblxuLyoqXG4gKiBDYWxscyBhIG1ldGhvZCBieSBpdHMgbnVtZXJpYyBJRCB3aXRoIHRoZSBzcGVjaWZpZWQgYXJndW1lbnRzLlxuICogU2VlIHtAbGluayBDYWxsfSBmb3IgZGV0YWlscy5cbiAqXG4gKiBAcGFyYW0gbWV0aG9kSUQgLSBUaGUgSUQgb2YgdGhlIG1ldGhvZCB0byBjYWxsLlxuICogQHBhcmFtIGFyZ3MgLSBUaGUgYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZC5cbiAqIEByZXR1cm4gVGhlIHJlc3VsdCBvZiB0aGUgbWV0aG9kIGNhbGwuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBCeUlEKG1ldGhvZElEOiBudW1iZXIsIC4uLmFyZ3M6IGFueVtdKTogQ2FuY2VsbGFibGVQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBDYWxsKHsgbWV0aG9kSUQsIGFyZ3MgfSk7XG59XG4iLCAiLy8gU291cmNlOiBodHRwczovL2dpdGh1Yi5jb20vaW5zcGVjdC1qcy9pcy1jYWxsYWJsZVxuXG4vLyBUaGUgTUlUIExpY2Vuc2UgKE1JVClcbi8vXG4vLyBDb3B5cmlnaHQgKGMpIDIwMTUgSm9yZGFuIEhhcmJhbmRcbi8vXG4vLyBQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5XG4vLyBvZiB0aGlzIHNvZnR3YXJlIGFuZCBhc3NvY2lhdGVkIGRvY3VtZW50YXRpb24gZmlsZXMgKHRoZSBcIlNvZnR3YXJlXCIpLCB0byBkZWFsXG4vLyBpbiB0aGUgU29mdHdhcmUgd2l0aG91dCByZXN0cmljdGlvbiwgaW5jbHVkaW5nIHdpdGhvdXQgbGltaXRhdGlvbiB0aGUgcmlnaHRzXG4vLyB0byB1c2UsIGNvcHksIG1vZGlmeSwgbWVyZ2UsIHB1Ymxpc2gsIGRpc3RyaWJ1dGUsIHN1YmxpY2Vuc2UsIGFuZC9vciBzZWxsXG4vLyBjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXNcbi8vIGZ1cm5pc2hlZCB0byBkbyBzbywgc3ViamVjdCB0byB0aGUgZm9sbG93aW5nIGNvbmRpdGlvbnM6XG4vL1xuLy8gVGhlIGFib3ZlIGNvcHlyaWdodCBub3RpY2UgYW5kIHRoaXMgcGVybWlzc2lvbiBub3RpY2Ugc2hhbGwgYmUgaW5jbHVkZWQgaW4gYWxsXG4vLyBjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLlxuLy9cbi8vIFRIRSBTT0ZUV0FSRSBJUyBQUk9WSURFRCBcIkFTIElTXCIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1Jcbi8vIElNUExJRUQsIElOQ0xVRElORyBCVVQgTk9UIExJTUlURUQgVE8gVEhFIFdBUlJBTlRJRVMgT0YgTUVSQ0hBTlRBQklMSVRZLFxuLy8gRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UgQU5EIE5PTklORlJJTkdFTUVOVC4gSU4gTk8gRVZFTlQgU0hBTEwgVEhFXG4vLyBBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSXG4vLyBMSUFCSUxJVFksIFdIRVRIRVIgSU4gQU4gQUNUSU9OIE9GIENPTlRSQUNULCBUT1JUIE9SIE9USEVSV0lTRSwgQVJJU0lORyBGUk9NLFxuLy8gT1VUIE9GIE9SIElOIENPTk5FQ1RJT04gV0lUSCBUSEUgU09GVFdBUkUgT1IgVEhFIFVTRSBPUiBPVEhFUiBERUFMSU5HUyBJTiBUSEVcbi8vIFNPRlRXQVJFLlxuXG52YXIgZm5Ub1N0ciA9IEZ1bmN0aW9uLnByb3RvdHlwZS50b1N0cmluZztcbnZhciByZWZsZWN0QXBwbHk6IHR5cGVvZiBSZWZsZWN0LmFwcGx5IHwgZmFsc2UgfCBudWxsID0gdHlwZW9mIFJlZmxlY3QgPT09ICdvYmplY3QnICYmIFJlZmxlY3QgIT09IG51bGwgJiYgUmVmbGVjdC5hcHBseTtcbnZhciBiYWRBcnJheUxpa2U6IGFueTtcbnZhciBpc0NhbGxhYmxlTWFya2VyOiBhbnk7XG5pZiAodHlwZW9mIHJlZmxlY3RBcHBseSA9PT0gJ2Z1bmN0aW9uJyAmJiB0eXBlb2YgT2JqZWN0LmRlZmluZVByb3BlcnR5ID09PSAnZnVuY3Rpb24nKSB7XG4gICAgdHJ5IHtcbiAgICAgICAgYmFkQXJyYXlMaWtlID0gT2JqZWN0LmRlZmluZVByb3BlcnR5KHt9LCAnbGVuZ3RoJywge1xuICAgICAgICAgICAgZ2V0OiBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICAgICAgdGhyb3cgaXNDYWxsYWJsZU1hcmtlcjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgICAgIGlzQ2FsbGFibGVNYXJrZXIgPSB7fTtcbiAgICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lIG5vLXRocm93LWxpdGVyYWxcbiAgICAgICAgcmVmbGVjdEFwcGx5KGZ1bmN0aW9uICgpIHsgdGhyb3cgNDI7IH0sIG51bGwsIGJhZEFycmF5TGlrZSk7XG4gICAgfSBjYXRjaCAoXykge1xuICAgICAgICBpZiAoXyAhPT0gaXNDYWxsYWJsZU1hcmtlcikge1xuICAgICAgICAgICAgcmVmbGVjdEFwcGx5ID0gbnVsbDtcbiAgICAgICAgfVxuICAgIH1cbn0gZWxzZSB7XG4gICAgcmVmbGVjdEFwcGx5ID0gbnVsbDtcbn1cblxudmFyIGNvbnN0cnVjdG9yUmVnZXggPSAvXlxccypjbGFzc1xcYi87XG52YXIgaXNFUzZDbGFzc0ZuID0gZnVuY3Rpb24gaXNFUzZDbGFzc0Z1bmN0aW9uKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcbiAgICB0cnkge1xuICAgICAgICB2YXIgZm5TdHIgPSBmblRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICByZXR1cm4gY29uc3RydWN0b3JSZWdleC50ZXN0KGZuU3RyKTtcbiAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTsgLy8gbm90IGEgZnVuY3Rpb25cbiAgICB9XG59O1xuXG52YXIgdHJ5RnVuY3Rpb25PYmplY3QgPSBmdW5jdGlvbiB0cnlGdW5jdGlvblRvU3RyKHZhbHVlOiBhbnkpOiBib29sZWFuIHtcbiAgICB0cnkge1xuICAgICAgICBpZiAoaXNFUzZDbGFzc0ZuKHZhbHVlKSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICAgICAgZm5Ub1N0ci5jYWxsKHZhbHVlKTtcbiAgICAgICAgcmV0dXJuIHRydWU7XG4gICAgfSBjYXRjaCAoZSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxufTtcbnZhciB0b1N0ciA9IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmc7XG52YXIgb2JqZWN0Q2xhc3MgPSAnW29iamVjdCBPYmplY3RdJztcbnZhciBmbkNsYXNzID0gJ1tvYmplY3QgRnVuY3Rpb25dJztcbnZhciBnZW5DbGFzcyA9ICdbb2JqZWN0IEdlbmVyYXRvckZ1bmN0aW9uXSc7XG52YXIgZGRhQ2xhc3MgPSAnW29iamVjdCBIVE1MQWxsQ29sbGVjdGlvbl0nOyAvLyBJRSAxMVxudmFyIGRkYUNsYXNzMiA9ICdbb2JqZWN0IEhUTUwgZG9jdW1lbnQuYWxsIGNsYXNzXSc7XG52YXIgZGRhQ2xhc3MzID0gJ1tvYmplY3QgSFRNTENvbGxlY3Rpb25dJzsgLy8gSUUgOS0xMFxudmFyIGhhc1RvU3RyaW5nVGFnID0gdHlwZW9mIFN5bWJvbCA9PT0gJ2Z1bmN0aW9uJyAmJiAhIVN5bWJvbC50b1N0cmluZ1RhZzsgLy8gYmV0dGVyOiB1c2UgYGhhcy10b3N0cmluZ3RhZ2BcblxudmFyIGlzSUU2OCA9ICEoMCBpbiBbLF0pOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIG5vLXNwYXJzZS1hcnJheXMsIGNvbW1hLXNwYWNpbmdcblxudmFyIGlzRERBOiAodmFsdWU6IGFueSkgPT4gYm9vbGVhbiA9IGZ1bmN0aW9uIGlzRG9jdW1lbnREb3RBbGwoKSB7IHJldHVybiBmYWxzZTsgfTtcbmlmICh0eXBlb2YgZG9jdW1lbnQgPT09ICdvYmplY3QnKSB7XG4gICAgLy8gRmlyZWZveCAzIGNhbm9uaWNhbGl6ZXMgRERBIHRvIHVuZGVmaW5lZCB3aGVuIGl0J3Mgbm90IGFjY2Vzc2VkIGRpcmVjdGx5XG4gICAgdmFyIGFsbCA9IGRvY3VtZW50LmFsbDtcbiAgICBpZiAodG9TdHIuY2FsbChhbGwpID09PSB0b1N0ci5jYWxsKGRvY3VtZW50LmFsbCkpIHtcbiAgICAgICAgaXNEREEgPSBmdW5jdGlvbiBpc0RvY3VtZW50RG90QWxsKHZhbHVlKSB7XG4gICAgICAgICAgICAvKiBnbG9iYWxzIGRvY3VtZW50OiBmYWxzZSAqL1xuICAgICAgICAgICAgLy8gaW4gSUUgNi04LCB0eXBlb2YgZG9jdW1lbnQuYWxsIGlzIFwib2JqZWN0XCIgYW5kIGl0J3MgdHJ1dGh5XG4gICAgICAgICAgICBpZiAoKGlzSUU2OCB8fCAhdmFsdWUpICYmICh0eXBlb2YgdmFsdWUgPT09ICd1bmRlZmluZWQnIHx8IHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcpKSB7XG4gICAgICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICAgICAgdmFyIHN0ciA9IHRvU3RyLmNhbGwodmFsdWUpO1xuICAgICAgICAgICAgICAgICAgICByZXR1cm4gKFxuICAgICAgICAgICAgICAgICAgICAgICAgc3RyID09PSBkZGFDbGFzc1xuICAgICAgICAgICAgICAgICAgICAgICAgfHwgc3RyID09PSBkZGFDbGFzczJcbiAgICAgICAgICAgICAgICAgICAgICAgIHx8IHN0ciA9PT0gZGRhQ2xhc3MzIC8vIG9wZXJhIDEyLjE2XG4gICAgICAgICAgICAgICAgICAgICAgICB8fCBzdHIgPT09IG9iamVjdENsYXNzIC8vIElFIDYtOFxuICAgICAgICAgICAgICAgICAgICApICYmIHZhbHVlKCcnKSA9PSBudWxsOyAvLyBlc2xpbnQtZGlzYWJsZS1saW5lIGVxZXFlcVxuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGUpIHsgLyoqLyB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5mdW5jdGlvbiBpc0NhbGxhYmxlUmVmQXBwbHk8VD4odmFsdWU6IFQgfCB1bmtub3duKTogdmFsdWUgaXMgKC4uLmFyZ3M6IGFueVtdKSA9PiBhbnkgIHtcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAodHlwZW9mIHZhbHVlICE9PSAnZnVuY3Rpb24nICYmIHR5cGVvZiB2YWx1ZSAhPT0gJ29iamVjdCcpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgdHJ5IHtcbiAgICAgICAgKHJlZmxlY3RBcHBseSBhcyBhbnkpKHZhbHVlLCBudWxsLCBiYWRBcnJheUxpa2UpO1xuICAgIH0gY2F0Y2ggKGUpIHtcbiAgICAgICAgaWYgKGUgIT09IGlzQ2FsbGFibGVNYXJrZXIpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgfVxuICAgIHJldHVybiAhaXNFUzZDbGFzc0ZuKHZhbHVlKSAmJiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7XG59XG5cbmZ1bmN0aW9uIGlzQ2FsbGFibGVOb1JlZkFwcGx5PFQ+KHZhbHVlOiBUIHwgdW5rbm93bik6IHZhbHVlIGlzICguLi5hcmdzOiBhbnlbXSkgPT4gYW55IHtcbiAgICBpZiAoaXNEREEodmFsdWUpKSB7IHJldHVybiB0cnVlOyB9XG4gICAgaWYgKCF2YWx1ZSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICBpZiAodHlwZW9mIHZhbHVlICE9PSAnZnVuY3Rpb24nICYmIHR5cGVvZiB2YWx1ZSAhPT0gJ29iamVjdCcpIHsgcmV0dXJuIGZhbHNlOyB9XG4gICAgaWYgKGhhc1RvU3RyaW5nVGFnKSB7IHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7IH1cbiAgICBpZiAoaXNFUzZDbGFzc0ZuKHZhbHVlKSkgeyByZXR1cm4gZmFsc2U7IH1cbiAgICB2YXIgc3RyQ2xhc3MgPSB0b1N0ci5jYWxsKHZhbHVlKTtcbiAgICBpZiAoc3RyQ2xhc3MgIT09IGZuQ2xhc3MgJiYgc3RyQ2xhc3MgIT09IGdlbkNsYXNzICYmICEoL15cXFtvYmplY3QgSFRNTC8pLnRlc3Qoc3RyQ2xhc3MpKSB7IHJldHVybiBmYWxzZTsgfVxuICAgIHJldHVybiB0cnlGdW5jdGlvbk9iamVjdCh2YWx1ZSk7XG59O1xuXG5leHBvcnQgZGVmYXVsdCByZWZsZWN0QXBwbHkgPyBpc0NhbGxhYmxlUmVmQXBwbHkgOiBpc0NhbGxhYmxlTm9SZWZBcHBseTtcbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IGlzQ2FsbGFibGUgZnJvbSBcIi4vY2FsbGFibGUuanNcIjtcblxuLyoqXG4gKiBFeGNlcHRpb24gY2xhc3MgdGhhdCB3aWxsIGJlIHVzZWQgYXMgcmVqZWN0aW9uIHJlYXNvblxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IGlzIGNhbmNlbGxlZCBzdWNjZXNzZnVsbHkuXG4gKlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgbmFtZX0gcHJvcGVydHkgaXMgdGhlIHN0cmluZyBgXCJDYW5jZWxFcnJvclwiYC5cbiAqIFRoZSB2YWx1ZSBvZiB0aGUge0BsaW5rIGNhdXNlfSBwcm9wZXJ0eSBpcyB0aGUgY2F1c2UgcGFzc2VkIHRvIHRoZSBjYW5jZWwgbWV0aG9kLCBpZiBhbnkuXG4gKi9cbmV4cG9ydCBjbGFzcyBDYW5jZWxFcnJvciBleHRlbmRzIEVycm9yIHtcbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxFcnJvcmAgaW5zdGFuY2UuXG4gICAgICogQHBhcmFtIG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cbiAgICAgKiBAcGFyYW0gb3B0aW9ucyAtIE9wdGlvbnMgdG8gYmUgZm9yd2FyZGVkIHRvIHRoZSBFcnJvciBjb25zdHJ1Y3Rvci5cbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihtZXNzYWdlPzogc3RyaW5nLCBvcHRpb25zPzogRXJyb3JPcHRpb25zKSB7XG4gICAgICAgIHN1cGVyKG1lc3NhZ2UsIG9wdGlvbnMpO1xuICAgICAgICB0aGlzLm5hbWUgPSBcIkNhbmNlbEVycm9yXCI7XG4gICAgfVxufVxuXG4vKipcbiAqIEV4Y2VwdGlvbiBjbGFzcyB0aGF0IHdpbGwgYmUgcmVwb3J0ZWQgYXMgYW4gdW5oYW5kbGVkIHJlamVjdGlvblxuICogaW4gY2FzZSBhIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2V9IHJlamVjdHMgYWZ0ZXIgYmVpbmcgY2FuY2VsbGVkLFxuICogb3Igd2hlbiB0aGUgYG9uY2FuY2VsbGVkYCBjYWxsYmFjayB0aHJvd3Mgb3IgcmVqZWN0cy5cbiAqXG4gKiBUaGUgdmFsdWUgb2YgdGhlIHtAbGluayBuYW1lfSBwcm9wZXJ0eSBpcyB0aGUgc3RyaW5nIGBcIkNhbmNlbGxlZFJlamVjdGlvbkVycm9yXCJgLlxuICogVGhlIHZhbHVlIG9mIHRoZSB7QGxpbmsgY2F1c2V9IHByb3BlcnR5IGlzIHRoZSByZWFzb24gdGhlIHByb21pc2UgcmVqZWN0ZWQgd2l0aC5cbiAqXG4gKiBCZWNhdXNlIHRoZSBvcmlnaW5hbCBwcm9taXNlIHdhcyBjYW5jZWxsZWQsXG4gKiBhIHdyYXBwZXIgcHJvbWlzZSB3aWxsIGJlIHBhc3NlZCB0byB0aGUgdW5oYW5kbGVkIHJlamVjdGlvbiBsaXN0ZW5lciBpbnN0ZWFkLlxuICogVGhlIHtAbGluayBwcm9taXNlfSBwcm9wZXJ0eSBob2xkcyBhIHJlZmVyZW5jZSB0byB0aGUgb3JpZ2luYWwgcHJvbWlzZS5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbGxlZFJlamVjdGlvbkVycm9yIGV4dGVuZHMgRXJyb3Ige1xuICAgIC8qKlxuICAgICAqIEhvbGRzIGEgcmVmZXJlbmNlIHRvIHRoZSBwcm9taXNlIHRoYXQgd2FzIGNhbmNlbGxlZCBhbmQgdGhlbiByZWplY3RlZC5cbiAgICAgKi9cbiAgICBwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj47XG5cbiAgICAvKipcbiAgICAgKiBDb25zdHJ1Y3RzIGEgbmV3IGBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcmAgaW5zdGFuY2UuXG4gICAgICogQHBhcmFtIHByb21pc2UgLSBUaGUgcHJvbWlzZSB0aGF0IGNhdXNlZCB0aGUgZXJyb3Igb3JpZ2luYWxseS5cbiAgICAgKiBAcGFyYW0gcmVhc29uIC0gVGhlIHJlamVjdGlvbiByZWFzb24uXG4gICAgICogQHBhcmFtIGluZm8gLSBBbiBvcHRpb25hbCBpbmZvcm1hdGl2ZSBtZXNzYWdlIHNwZWNpZnlpbmcgdGhlIGNpcmN1bXN0YW5jZXMgaW4gd2hpY2ggdGhlIGVycm9yIHdhcyB0aHJvd24uXG4gICAgICogICAgICAgICAgICAgICBEZWZhdWx0cyB0byB0aGUgc3RyaW5nIGBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsbGVkIHByb21pc2UuXCJgLlxuICAgICAqL1xuICAgIGNvbnN0cnVjdG9yKHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiwgcmVhc29uPzogYW55LCBpbmZvPzogc3RyaW5nKSB7XG4gICAgICAgIHN1cGVyKChpbmZvID8/IFwiVW5oYW5kbGVkIHJlamVjdGlvbiBpbiBjYW5jZWxsZWQgcHJvbWlzZS5cIikgKyBcIiBSZWFzb246IFwiICsgZXJyb3JNZXNzYWdlKHJlYXNvbiksIHsgY2F1c2U6IHJlYXNvbiB9KTtcbiAgICAgICAgdGhpcy5wcm9taXNlID0gcHJvbWlzZTtcbiAgICAgICAgdGhpcy5uYW1lID0gXCJDYW5jZWxsZWRSZWplY3Rpb25FcnJvclwiO1xuICAgIH1cbn1cblxudHlwZSBDYW5jZWxsYWJsZVByb21pc2VSZXNvbHZlcjxUPiA9ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+IHwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPikgPT4gdm9pZDtcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IgPSAocmVhc29uPzogYW55KSA9PiB2b2lkO1xudHlwZSBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIgPSAoY2F1c2U/OiBhbnkpID0+IHZvaWQgfCBQcm9taXNlTGlrZTx2b2lkPjtcbnR5cGUgQ2FuY2VsbGFibGVQcm9taXNlRXhlY3V0b3I8VD4gPSAocmVzb2x2ZTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4sIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3IpID0+IHZvaWQ7XG5cbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XG4gICAgdGhlbjxUUmVzdWx0MSA9IFQsIFRSZXN1bHQyID0gbmV2ZXI+KG9uZnVsZmlsbGVkPzogKCh2YWx1ZTogVCkgPT4gVFJlc3VsdDEgfCBQcm9taXNlTGlrZTxUUmVzdWx0MT4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gVFJlc3VsdDIgfCBQcm9taXNlTGlrZTxUUmVzdWx0Mj4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQyPikgfCB1bmRlZmluZWQgfCBudWxsKTogQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUUmVzdWx0MSB8IFRSZXN1bHQyPjtcbiAgICBjYW5jZWwoY2F1c2U/OiBhbnkpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG59XG5cbi8qKlxuICogV3JhcHMgYSBjYW5jZWxsYWJsZSBwcm9taXNlIGFsb25nIHdpdGggaXRzIHJlc29sdXRpb24gbWV0aG9kcy5cbiAqIFRoZSBgb25jYW5jZWxsZWRgIGZpZWxkIHdpbGwgYmUgbnVsbCBpbml0aWFsbHkgYnV0IG1heSBiZSBzZXQgdG8gcHJvdmlkZSBhIGN1c3RvbSBjYW5jZWxsYXRpb24gZnVuY3Rpb24uXG4gKi9cbmV4cG9ydCBpbnRlcmZhY2UgQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XG4gICAgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+O1xuICAgIHJlc29sdmU6IENhbmNlbGxhYmxlUHJvbWlzZVJlc29sdmVyPFQ+O1xuICAgIHJlamVjdDogQ2FuY2VsbGFibGVQcm9taXNlUmVqZWN0b3I7XG4gICAgb25jYW5jZWxsZWQ6IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlciB8IG51bGw7XG59XG5cbmludGVyZmFjZSBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSB7XG4gICAgcmVhZG9ubHkgcm9vdDogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGU7XG4gICAgcmVzb2x2aW5nOiBib29sZWFuO1xuICAgIHNldHRsZWQ6IGJvb2xlYW47XG4gICAgcmVhc29uPzogQ2FuY2VsRXJyb3I7XG59XG5cbi8vIFByaXZhdGUgZmllbGQgbmFtZXMuXG5jb25zdCBiYXJyaWVyU3ltID0gU3ltYm9sKFwiYmFycmllclwiKTtcbmNvbnN0IGNhbmNlbEltcGxTeW0gPSBTeW1ib2woXCJjYW5jZWxJbXBsXCIpO1xuY29uc3Qgc3BlY2llczogdHlwZW9mIFN5bWJvbC5zcGVjaWVzID0gU3ltYm9sLnNwZWNpZXMgPz8gU3ltYm9sKFwic3BlY2llc1BvbHlmaWxsXCIpO1xuXG4vKipcbiAqIEEgcHJvbWlzZSB3aXRoIGFuIGF0dGFjaGVkIG1ldGhvZCBmb3IgY2FuY2VsbGluZyBsb25nLXJ1bm5pbmcgb3BlcmF0aW9ucyAoc2VlIHtAbGluayBDYW5jZWxsYWJsZVByb21pc2UjY2FuY2VsfSkuXG4gKiBDYW5jZWxsYXRpb24gY2FuIG9wdGlvbmFsbHkgYmUgYm91bmQgdG8gYW4ge0BsaW5rIEFib3J0U2lnbmFsfVxuICogZm9yIGJldHRlciBjb21wb3NhYmlsaXR5IChzZWUge0BsaW5rIENhbmNlbGxhYmxlUHJvbWlzZSNjYW5jZWxPbn0pLlxuICpcbiAqIENhbmNlbGxpbmcgYSBwZW5kaW5nIHByb21pc2Ugd2lsbCByZXN1bHQgaW4gYW4gaW1tZWRpYXRlIHJlamVjdGlvblxuICogd2l0aCBhbiBpbnN0YW5jZSBvZiB7QGxpbmsgQ2FuY2VsRXJyb3J9IGFzIHJlYXNvbixcbiAqIGJ1dCB3aG9ldmVyIHN0YXJ0ZWQgdGhlIHByb21pc2Ugd2lsbCBiZSByZXNwb25zaWJsZVxuICogZm9yIGFjdHVhbGx5IGFib3J0aW5nIHRoZSB1bmRlcmx5aW5nIG9wZXJhdGlvbi5cbiAqIFRvIHRoaXMgcHVycG9zZSwgdGhlIGNvbnN0cnVjdG9yIGFuZCBhbGwgY2hhaW5pbmcgbWV0aG9kc1xuICogYWNjZXB0IG9wdGlvbmFsIGNhbmNlbGxhdGlvbiBjYWxsYmFja3MuXG4gKlxuICogSWYgYSBgQ2FuY2VsbGFibGVQcm9taXNlYCBzdGlsbCByZXNvbHZlcyBhZnRlciBoYXZpbmcgYmVlbiBjYW5jZWxsZWQsXG4gKiB0aGUgcmVzdWx0IHdpbGwgYmUgZGlzY2FyZGVkLiBJZiBpdCByZWplY3RzLCB0aGUgcmVhc29uXG4gKiB3aWxsIGJlIHJlcG9ydGVkIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24sXG4gKiB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfSBpbnN0YW5jZS5cbiAqIFRvIGZhY2lsaXRhdGUgdGhlIGhhbmRsaW5nIG9mIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyxcbiAqIGNhbmNlbGxlZCBgQ2FuY2VsbGFibGVQcm9taXNlYHMgd2lsbCBfbm90XyByZXBvcnQgdW5oYW5kbGVkIGBDYW5jZWxFcnJvcmBzXG4gKiB3aG9zZSBgY2F1c2VgIGZpZWxkIGlzIHRoZSBzYW1lIGFzIHRoZSBvbmUgd2l0aCB3aGljaCB0aGUgY3VycmVudCBwcm9taXNlIHdhcyBjYW5jZWxsZWQuXG4gKlxuICogQWxsIHVzdWFsIHByb21pc2UgbWV0aG9kcyBhcmUgZGVmaW5lZCBhbmQgcmV0dXJuIGEgYENhbmNlbGxhYmxlUHJvbWlzZWBcbiAqIHdob3NlIGNhbmNlbCBtZXRob2Qgd2lsbCBjYW5jZWwgdGhlIHBhcmVudCBvcGVyYXRpb24gYXMgd2VsbCwgcHJvcGFnYXRpbmcgdGhlIGNhbmNlbGxhdGlvbiByZWFzb25cbiAqIHVwd2FyZHMgdGhyb3VnaCBwcm9taXNlIGNoYWlucy5cbiAqIENvbnZlcnNlbHksIGNhbmNlbGxpbmcgYSBwcm9taXNlIHdpbGwgbm90IGF1dG9tYXRpY2FsbHkgY2FuY2VsIGRlcGVuZGVudCBwcm9taXNlcyBkb3duc3RyZWFtOlxuICogYGBgdHNcbiAqIGxldCByb290ID0gbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KTtcbiAqIGxldCBjaGlsZDEgPSByb290LnRoZW4oKCkgPT4geyAuLi4gfSk7XG4gKiBsZXQgY2hpbGQyID0gY2hpbGQxLnRoZW4oKCkgPT4geyAuLi4gfSk7XG4gKiBsZXQgY2hpbGQzID0gcm9vdC5jYXRjaCgoKSA9PiB7IC4uLiB9KTtcbiAqIGNoaWxkMS5jYW5jZWwoKTsgLy8gQ2FuY2VscyBjaGlsZDEgYW5kIHJvb3QsIGJ1dCBub3QgY2hpbGQyIG9yIGNoaWxkM1xuICogYGBgXG4gKiBDYW5jZWxsaW5nIGEgcHJvbWlzZSB0aGF0IGhhcyBhbHJlYWR5IHNldHRsZWQgaXMgc2FmZSBhbmQgaGFzIG5vIGNvbnNlcXVlbmNlLlxuICpcbiAqIFRoZSBgY2FuY2VsYCBtZXRob2QgcmV0dXJucyBhIHByb21pc2UgdGhhdCBfYWx3YXlzIGZ1bGZpbGxzX1xuICogYWZ0ZXIgdGhlIHdob2xlIGNoYWluIGhhcyBwcm9jZXNzZWQgdGhlIGNhbmNlbCByZXF1ZXN0XG4gKiBhbmQgYWxsIGF0dGFjaGVkIGNhbGxiYWNrcyB1cCB0byB0aGF0IG1vbWVudCBoYXZlIHJ1bi5cbiAqXG4gKiBBbGwgRVMyMDI0IHByb21pc2UgbWV0aG9kcyAoc3RhdGljIGFuZCBpbnN0YW5jZSkgYXJlIGRlZmluZWQgb24gQ2FuY2VsbGFibGVQcm9taXNlLFxuICogYnV0IGFjdHVhbCBhdmFpbGFiaWxpdHkgbWF5IHZhcnkgd2l0aCBPUy93ZWJ2aWV3IHZlcnNpb24uXG4gKlxuICogSW4gbGluZSB3aXRoIHRoZSBwcm9wb3NhbCBhdCBodHRwczovL2dpdGh1Yi5jb20vdGMzOS9wcm9wb3NhbC1ybS1idWlsdGluLXN1YmNsYXNzaW5nLFxuICogYENhbmNlbGxhYmxlUHJvbWlzZWAgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy5cbiAqIEV4dGVuZGVycyBzaG91bGQgdGFrZSBjYXJlIHRvIHByb3ZpZGUgdGhlaXIgb3duIG1ldGhvZCBpbXBsZW1lbnRhdGlvbnMuXG4gKiBUaGlzIG1pZ2h0IGJlIHJlY29uc2lkZXJlZCBpbiBjYXNlIHRoZSBwcm9wb3NhbCBpcyByZXRpcmVkLlxuICpcbiAqIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBhIHdyYXBwZXIgYXJvdW5kIHRoZSBET00gUHJvbWlzZSBvYmplY3RcbiAqIGFuZCBpcyBjb21wbGlhbnQgd2l0aCB0aGUgW1Byb21pc2VzL0ErIHNwZWNpZmljYXRpb25dKGh0dHBzOi8vcHJvbWlzZXNhcGx1cy5jb20vKVxuICogKGl0IHBhc3NlcyB0aGUgW2NvbXBsaWFuY2Ugc3VpdGVdKGh0dHBzOi8vZ2l0aHViLmNvbS9wcm9taXNlcy1hcGx1cy9wcm9taXNlcy10ZXN0cykpXG4gKiBpZiBzbyBpcyB0aGUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbi5cbiAqL1xuZXhwb3J0IGNsYXNzIENhbmNlbGxhYmxlUHJvbWlzZTxUPiBleHRlbmRzIFByb21pc2U8VD4gaW1wbGVtZW50cyBQcm9taXNlTGlrZTxUPiwgQ2FuY2VsbGFibGVQcm9taXNlTGlrZTxUPiB7XG4gICAgLy8gUHJpdmF0ZSBmaWVsZHMuXG4gICAgLyoqIEBpbnRlcm5hbCAqL1xuICAgIHByaXZhdGUgW2JhcnJpZXJTeW1dITogUGFydGlhbDxQcm9taXNlV2l0aFJlc29sdmVyczx2b2lkPj4gfCBudWxsO1xuICAgIC8qKiBAaW50ZXJuYWwgKi9cbiAgICBwcml2YXRlIHJlYWRvbmx5IFtjYW5jZWxJbXBsU3ltXSE6IChyZWFzb246IENhbmNlbEVycm9yKSA9PiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IGBDYW5jZWxsYWJsZVByb21pc2VgLlxuICAgICAqXG4gICAgICogQHBhcmFtIGV4ZWN1dG9yIC0gQSBjYWxsYmFjayB1c2VkIHRvIGluaXRpYWxpemUgdGhlIHByb21pc2UuIFRoaXMgY2FsbGJhY2sgaXMgcGFzc2VkIHR3byBhcmd1bWVudHM6XG4gICAgICogICAgICAgICAgICAgICAgICAgYSBgcmVzb2x2ZWAgY2FsbGJhY2sgdXNlZCB0byByZXNvbHZlIHRoZSBwcm9taXNlIHdpdGggYSB2YWx1ZVxuICAgICAqICAgICAgICAgICAgICAgICAgIG9yIHRoZSByZXN1bHQgb2YgYW5vdGhlciBwcm9taXNlIChwb3NzaWJseSBjYW5jZWxsYWJsZSksXG4gICAgICogICAgICAgICAgICAgICAgICAgYW5kIGEgYHJlamVjdGAgY2FsbGJhY2sgdXNlZCB0byByZWplY3QgdGhlIHByb21pc2Ugd2l0aCBhIHByb3ZpZGVkIHJlYXNvbiBvciBlcnJvci5cbiAgICAgKiAgICAgICAgICAgICAgICAgICBJZiB0aGUgdmFsdWUgcHJvdmlkZWQgdG8gdGhlIGByZXNvbHZlYCBjYWxsYmFjayBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqICAgICAgICAgICAgICAgICAgIChpdCBoYXMgYSBgdGhlbmAgX2FuZF8gYSBgY2FuY2VsYCBtZXRob2QpLFxuICAgICAqICAgICAgICAgICAgICAgICAgIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGZvcndhcmRlZCB0byB0aGF0IG9iamVjdCBhbmQgdGhlIG9uY2FuY2VsbGVkIHdpbGwgbm90IGJlIGludm9rZWQgYW55bW9yZS5cbiAgICAgKiAgICAgICAgICAgICAgICAgICBJZiBhbnkgb25lIG9mIHRoZSB0d28gY2FsbGJhY2tzIGlzIGNhbGxlZCBfYWZ0ZXJfIHRoZSBwcm9taXNlIGhhcyBiZWVuIGNhbmNlbGxlZCxcbiAgICAgKiAgICAgICAgICAgICAgICAgICB0aGUgcHJvdmlkZWQgdmFsdWVzIHdpbGwgYmUgY2FuY2VsbGVkIGFuZCByZXNvbHZlZCBhcyB1c3VhbCxcbiAgICAgKiAgICAgICAgICAgICAgICAgICBidXQgdGhlaXIgcmVzdWx0cyB3aWxsIGJlIGRpc2NhcmRlZC5cbiAgICAgKiAgICAgICAgICAgICAgICAgICBIb3dldmVyLCBpZiB0aGUgcmVzb2x1dGlvbiBwcm9jZXNzIHVsdGltYXRlbHkgZW5kcyB1cCBpbiBhIHJlamVjdGlvblxuICAgICAqICAgICAgICAgICAgICAgICAgIHRoYXQgaXMgbm90IGR1ZSB0byBjYW5jZWxsYXRpb24sIHRoZSByZWplY3Rpb24gcmVhc29uXG4gICAgICogICAgICAgICAgICAgICAgICAgd2lsbCBiZSB3cmFwcGVkIGluIGEge0BsaW5rIENhbmNlbGxlZFJlamVjdGlvbkVycm9yfVxuICAgICAqICAgICAgICAgICAgICAgICAgIGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICogQHBhcmFtIG9uY2FuY2VsbGVkIC0gSXQgaXMgdGhlIGNhbGxlcidzIHJlc3BvbnNpYmlsaXR5IHRvIGVuc3VyZSB0aGF0IGFueSBvcGVyYXRpb25cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBzdGFydGVkIGJ5IHRoZSBleGVjdXRvciBpcyBwcm9wZXJseSBoYWx0ZWQgdXBvbiBjYW5jZWxsYXRpb24uXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgVGhpcyBvcHRpb25hbCBjYWxsYmFjayBjYW4gYmUgdXNlZCB0byB0aGF0IHB1cnBvc2UuXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgSXQgd2lsbCBiZSBjYWxsZWQgX3N5bmNocm9ub3VzbHlfIHdpdGggYSBjYW5jZWxsYXRpb24gY2F1c2VcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aGVuIGNhbmNlbGxhdGlvbiBpcyByZXF1ZXN0ZWQsIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgcmVqZWN0ZWRcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSwgYnV0IF9iZWZvcmVfXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgYW55IHtAbGluayB0aGVufS97QGxpbmsgY2F0Y2h9L3tAbGluayBmaW5hbGx5fSBjYWxsYmFjayBydW5zLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBjYWxsYmFjayByZXR1cm5zIGEgdGhlbmFibGUsIHRoZSBwcm9taXNlIHJldHVybmVkIGZyb20ge0BsaW5rIGNhbmNlbH1cbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICB3aWxsIG9ubHkgZnVsZmlsbCBhZnRlciB0aGUgZm9ybWVyIGhhcyBzZXR0bGVkLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIFVuaGFuZGxlZCBleGNlcHRpb25zIG9yIHJlamVjdGlvbnMgZnJvbSB0aGUgY2FsbGJhY2sgd2lsbCBiZSB3cmFwcGVkXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIHVuaGFuZGxlZCByZWplY3Rpb25zLlxuICAgICAqICAgICAgICAgICAgICAgICAgICAgIElmIHRoZSBgcmVzb2x2ZWAgY2FsbGJhY2sgaXMgY2FsbGVkIGJlZm9yZSBjYW5jZWxsYXRpb24gd2l0aCBhIGNhbmNlbGxhYmxlIHByb21pc2UsXG4gICAgICogICAgICAgICAgICAgICAgICAgICAgY2FuY2VsbGF0aW9uIHJlcXVlc3RzIG9uIHRoaXMgcHJvbWlzZSB3aWxsIGJlIGRpdmVydGVkIHRvIHRoYXQgcHJvbWlzZSxcbiAgICAgKiAgICAgICAgICAgICAgICAgICAgICBhbmQgdGhlIG9yaWdpbmFsIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXhlY3V0b3I6IENhbmNlbGxhYmxlUHJvbWlzZUV4ZWN1dG9yPFQ+LCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcikge1xuICAgICAgICBsZXQgcmVzb2x2ZSE6ICh2YWx1ZTogVCB8IFByb21pc2VMaWtlPFQ+KSA9PiB2b2lkO1xuICAgICAgICBsZXQgcmVqZWN0ITogKHJlYXNvbj86IGFueSkgPT4gdm9pZDtcbiAgICAgICAgc3VwZXIoKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XG5cbiAgICAgICAgaWYgKCh0aGlzLmNvbnN0cnVjdG9yIGFzIGFueSlbc3BlY2llc10gIT09IFByb21pc2UpIHtcbiAgICAgICAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoXCJDYW5jZWxsYWJsZVByb21pc2UgZG9lcyBub3Qgc3VwcG9ydCB0cmFuc3BhcmVudCBzdWJjbGFzc2luZy4gUGxlYXNlIHJlZnJhaW4gZnJvbSBvdmVycmlkaW5nIHRoZSBbU3ltYm9sLnNwZWNpZXNdIHN0YXRpYyBwcm9wZXJ0eS5cIik7XG4gICAgICAgIH1cblxuICAgICAgICBsZXQgcHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiA9IHtcbiAgICAgICAgICAgIHByb21pc2U6IHRoaXMsXG4gICAgICAgICAgICByZXNvbHZlLFxuICAgICAgICAgICAgcmVqZWN0LFxuICAgICAgICAgICAgZ2V0IG9uY2FuY2VsbGVkKCkgeyByZXR1cm4gb25jYW5jZWxsZWQgPz8gbnVsbDsgfSxcbiAgICAgICAgICAgIHNldCBvbmNhbmNlbGxlZChjYikgeyBvbmNhbmNlbGxlZCA9IGNiID8/IHVuZGVmaW5lZDsgfVxuICAgICAgICB9O1xuXG4gICAgICAgIGNvbnN0IHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSA9IHtcbiAgICAgICAgICAgIGdldCByb290KCkgeyByZXR1cm4gc3RhdGU7IH0sXG4gICAgICAgICAgICByZXNvbHZpbmc6IGZhbHNlLFxuICAgICAgICAgICAgc2V0dGxlZDogZmFsc2VcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBTZXR1cCBjYW5jZWxsYXRpb24gc3lzdGVtLlxuICAgICAgICB2b2lkIE9iamVjdC5kZWZpbmVQcm9wZXJ0aWVzKHRoaXMsIHtcbiAgICAgICAgICAgIFtiYXJyaWVyU3ltXToge1xuICAgICAgICAgICAgICAgIGNvbmZpZ3VyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgZW51bWVyYWJsZTogZmFsc2UsXG4gICAgICAgICAgICAgICAgd3JpdGFibGU6IHRydWUsXG4gICAgICAgICAgICAgICAgdmFsdWU6IG51bGxcbiAgICAgICAgICAgIH0sXG4gICAgICAgICAgICBbY2FuY2VsSW1wbFN5bV06IHtcbiAgICAgICAgICAgICAgICBjb25maWd1cmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIGVudW1lcmFibGU6IGZhbHNlLFxuICAgICAgICAgICAgICAgIHdyaXRhYmxlOiBmYWxzZSxcbiAgICAgICAgICAgICAgICB2YWx1ZTogY2FuY2VsbGVyRm9yKHByb21pc2UsIHN0YXRlKVxuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcblxuICAgICAgICAvLyBSdW4gdGhlIGFjdHVhbCBleGVjdXRvci5cbiAgICAgICAgY29uc3QgcmVqZWN0b3IgPSByZWplY3RvckZvcihwcm9taXNlLCBzdGF0ZSk7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBleGVjdXRvcihyZXNvbHZlckZvcihwcm9taXNlLCBzdGF0ZSksIHJlamVjdG9yKTtcbiAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICBpZiAoc3RhdGUucmVzb2x2aW5nKSB7XG4gICAgICAgICAgICAgICAgY29uc29sZS5sb2coXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIENhbmNlbGxhYmxlUHJvbWlzZSBleGVjdXRvci5cIiwgZXJyKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgcmVqZWN0b3IoZXJyKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENhbmNlbHMgaW1tZWRpYXRlbHkgdGhlIGV4ZWN1dGlvbiBvZiB0aGUgb3BlcmF0aW9uIGFzc29jaWF0ZWQgd2l0aCB0aGlzIHByb21pc2UuXG4gICAgICogVGhlIHByb21pc2UgcmVqZWN0cyB3aXRoIGEge0BsaW5rIENhbmNlbEVycm9yfSBpbnN0YW5jZSBhcyByZWFzb24sXG4gICAgICogd2l0aCB0aGUge0BsaW5rIENhbmNlbEVycm9yI2NhdXNlfSBwcm9wZXJ0eSBzZXQgdG8gdGhlIGdpdmVuIGFyZ3VtZW50LCBpZiBhbnkuXG4gICAgICpcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCBhZnRlciB0aGUgcHJvbWlzZSBoYXMgYWxyZWFkeSBzZXR0bGVkO1xuICAgICAqIHJlcGVhdGVkIGNhbGxzIGluIHBhcnRpY3VsYXIgYXJlIHNhZmUsIGJ1dCBvbmx5IHRoZSBmaXJzdCBvbmVcbiAgICAgKiB3aWxsIHNldCB0aGUgY2FuY2VsbGF0aW9uIGNhdXNlLlxuICAgICAqXG4gICAgICogVGhlIGBDYW5jZWxFcnJvcmAgZXhjZXB0aW9uIF9uZWVkIG5vdF8gYmUgaGFuZGxlZCBleHBsaWNpdGx5IF9vbiB0aGUgcHJvbWlzZXMgdGhhdCBhcmUgYmVpbmcgY2FuY2VsbGVkOl9cbiAgICAgKiBjYW5jZWxsaW5nIGEgcHJvbWlzZSB3aXRoIG5vIGF0dGFjaGVkIHJlamVjdGlvbiBoYW5kbGVyIGRvZXMgbm90IHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudC5cbiAgICAgKiBUaGVyZWZvcmUsIHRoZSBmb2xsb3dpbmcgaWRpb21zIGFyZSBhbGwgZXF1YWxseSBjb3JyZWN0OlxuICAgICAqIGBgYHRzXG4gICAgICogbmV3IENhbmNlbGxhYmxlUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7IC4uLiB9KS5jYW5jZWwoKTtcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYW5jZWwoKTtcbiAgICAgKiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pLnRoZW4oLi4uKS5jYXRjaCguLi4pLmNhbmNlbCgpO1xuICAgICAqIGBgYFxuICAgICAqIFdoZW5ldmVyIHNvbWUgY2FuY2VsbGVkIHByb21pc2UgaW4gYSBjaGFpbiByZWplY3RzIHdpdGggYSBgQ2FuY2VsRXJyb3JgXG4gICAgICogd2l0aCB0aGUgc2FtZSBjYW5jZWxsYXRpb24gY2F1c2UgYXMgaXRzZWxmLCB0aGUgZXJyb3Igd2lsbCBiZSBkaXNjYXJkZWQgc2lsZW50bHkuXG4gICAgICogSG93ZXZlciwgdGhlIGBDYW5jZWxFcnJvcmAgX3dpbGwgc3RpbGwgYmUgZGVsaXZlcmVkXyB0byBhbGwgYXR0YWNoZWQgcmVqZWN0aW9uIGhhbmRsZXJzXG4gICAgICogYWRkZWQgYnkge0BsaW5rIHRoZW59IGFuZCByZWxhdGVkIG1ldGhvZHM6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICAgICAqIGNhbmNlbGxhYmxlLnRoZW4oKCkgPT4geyAuLi4gfSkuY2F0Y2goY29uc29sZS5sb2cpO1xuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBBIENhbmNlbEVycm9yIGlzIHByaW50ZWQgdG8gdGhlIGNvbnNvbGUuXG4gICAgICogYGBgXG4gICAgICogSWYgdGhlIGBDYW5jZWxFcnJvcmAgaXMgbm90IGhhbmRsZWQgZG93bnN0cmVhbSBieSB0aGUgdGltZSBpdCByZWFjaGVzXG4gICAgICogYSBfbm9uLWNhbmNlbGxlZF8gcHJvbWlzZSwgaXQgX3dpbGxfIHRyaWdnZXIgYW4gdW5oYW5kbGVkIHJlamVjdGlvbiBldmVudCxcbiAgICAgKiBqdXN0IGxpa2Ugbm9ybWFsIHJlamVjdGlvbnMgd291bGQ6XG4gICAgICogYGBgdHNcbiAgICAgKiBsZXQgY2FuY2VsbGFibGUgPSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHsgLi4uIH0pO1xuICAgICAqIGxldCBjaGFpbmVkID0gY2FuY2VsbGFibGUudGhlbigoKSA9PiB7IC4uLiB9KS50aGVuKCgpID0+IHsgLi4uIH0pOyAvLyBObyBjYXRjaC4uLlxuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpOyAvLyBVbmhhbmRsZWQgcmVqZWN0aW9uIGV2ZW50IG9uIGNoYWluZWQhXG4gICAgICogYGBgXG4gICAgICogVGhlcmVmb3JlLCBpdCBpcyBpbXBvcnRhbnQgdG8gZWl0aGVyIGNhbmNlbCB3aG9sZSBwcm9taXNlIGNoYWlucyBmcm9tIHRoZWlyIHRhaWwsXG4gICAgICogYXMgc2hvd24gaW4gdGhlIGNvcnJlY3QgaWRpb21zIGFib3ZlLCBvciB0YWtlIGNhcmUgb2YgaGFuZGxpbmcgZXJyb3JzIGV2ZXJ5d2hlcmUuXG4gICAgICpcbiAgICAgKiBAcmV0dXJucyBBIGNhbmNlbGxhYmxlIHByb21pc2UgdGhhdCBfZnVsZmlsbHNfIGFmdGVyIHRoZSBjYW5jZWwgY2FsbGJhY2sgKGlmIGFueSlcbiAgICAgKiBhbmQgYWxsIGhhbmRsZXJzIGF0dGFjaGVkIHVwIHRvIHRoZSBjYWxsIHRvIGNhbmNlbCBoYXZlIHJ1bi5cbiAgICAgKiBJZiB0aGUgY2FuY2VsIGNhbGxiYWNrIHJldHVybnMgYSB0aGVuYWJsZSwgdGhlIHByb21pc2UgcmV0dXJuZWQgYnkgYGNhbmNlbGBcbiAgICAgKiB3aWxsIGFsc28gd2FpdCBmb3IgdGhhdCB0aGVuYWJsZSB0byBzZXR0bGUuXG4gICAgICogVGhpcyBlbmFibGVzIGNhbGxlcnMgdG8gd2FpdCBmb3IgdGhlIGNhbmNlbGxlZCBvcGVyYXRpb24gdG8gdGVybWluYXRlXG4gICAgICogd2l0aG91dCBiZWluZyBmb3JjZWQgdG8gaGFuZGxlIHBvdGVudGlhbCBlcnJvcnMgYXQgdGhlIGNhbGwgc2l0ZS5cbiAgICAgKiBgYGB0c1xuICAgICAqIGNhbmNlbGxhYmxlLmNhbmNlbCgpLnRoZW4oKCkgPT4ge1xuICAgICAqICAgICAvLyBDbGVhbnVwIGZpbmlzaGVkLCBpdCdzIHNhZmUgdG8gZG8gc29tZXRoaW5nIGVsc2UuXG4gICAgICogfSwgKGVycikgPT4ge1xuICAgICAqICAgICAvLyBVbnJlYWNoYWJsZTogdGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSBjYW5jZWwgd2lsbCBuZXZlciByZWplY3QuXG4gICAgICogfSk7XG4gICAgICogYGBgXG4gICAgICogTm90ZSB0aGF0IHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgX25vdF8gaGFuZGxlIGltcGxpY2l0bHkgYW55IHJlamVjdGlvblxuICAgICAqIHRoYXQgbWlnaHQgaGF2ZSBvY2N1cnJlZCBhbHJlYWR5IGluIHRoZSBjYW5jZWxsZWQgY2hhaW4uXG4gICAgICogSXQgd2lsbCBqdXN0IHRyYWNrIHdoZXRoZXIgcmVnaXN0ZXJlZCBoYW5kbGVycyBoYXZlIGJlZW4gZXhlY3V0ZWQgb3Igbm90LlxuICAgICAqIFRoZXJlZm9yZSwgdW5oYW5kbGVkIHJlamVjdGlvbnMgd2lsbCBuZXZlciBiZSBzaWxlbnRseSBoYW5kbGVkIGJ5IGNhbGxpbmcgY2FuY2VsLlxuICAgICAqL1xuICAgIGNhbmNlbChjYXVzZT86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHZvaWQ+KChyZXNvbHZlKSA9PiB7XG4gICAgICAgICAgICAvLyBJTlZBUklBTlQ6IHRoZSByZXN1bHQgb2YgdGhpc1tjYW5jZWxJbXBsU3ltXSBhbmQgdGhlIGJhcnJpZXIgZG8gbm90IGV2ZXIgcmVqZWN0LlxuICAgICAgICAgICAgLy8gVW5mb3J0dW5hdGVseSBtYWNPUyBIaWdoIFNpZXJyYSBkb2VzIG5vdCBzdXBwb3J0IFByb21pc2UuYWxsU2V0dGxlZC5cbiAgICAgICAgICAgIFByb21pc2UuYWxsKFtcbiAgICAgICAgICAgICAgICB0aGlzW2NhbmNlbEltcGxTeW1dKG5ldyBDYW5jZWxFcnJvcihcIlByb21pc2UgY2FuY2VsbGVkLlwiLCB7IGNhdXNlIH0pKSxcbiAgICAgICAgICAgICAgICBjdXJyZW50QmFycmllcih0aGlzKVxuICAgICAgICAgICAgXSkudGhlbigoKSA9PiByZXNvbHZlKCksICgpID0+IHJlc29sdmUoKSk7XG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEJpbmRzIHByb21pc2UgY2FuY2VsbGF0aW9uIHRvIHRoZSBhYm9ydCBldmVudCBvZiB0aGUgZ2l2ZW4ge0BsaW5rIEFib3J0U2lnbmFsfS5cbiAgICAgKiBJZiB0aGUgc2lnbmFsIGhhcyBhbHJlYWR5IGFib3J0ZWQsIHRoZSBwcm9taXNlIHdpbGwgYmUgY2FuY2VsbGVkIGltbWVkaWF0ZWx5LlxuICAgICAqIFdoZW4gZWl0aGVyIGNvbmRpdGlvbiBpcyB2ZXJpZmllZCwgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZSB3aWxsIGJlIHNldFxuICAgICAqIHRvIHRoZSBzaWduYWwncyBhYm9ydCByZWFzb24gKHNlZSB7QGxpbmsgQWJvcnRTaWduYWwjcmVhc29ufSkuXG4gICAgICpcbiAgICAgKiBIYXMgbm8gZWZmZWN0IGlmIGNhbGxlZCAob3IgaWYgdGhlIHNpZ25hbCBhYm9ydHMpIF9hZnRlcl8gdGhlIHByb21pc2UgaGFzIGFscmVhZHkgc2V0dGxlZC5cbiAgICAgKiBPbmx5IHRoZSBmaXJzdCBzaWduYWwgdG8gYWJvcnQgd2lsbCBzZXQgdGhlIGNhbmNlbGxhdGlvbiBjYXVzZS5cbiAgICAgKlxuICAgICAqIEZvciBtb3JlIGRldGFpbHMgYWJvdXQgdGhlIGNhbmNlbGxhdGlvbiBwcm9jZXNzLFxuICAgICAqIHNlZSB7QGxpbmsgY2FuY2VsfSBhbmQgdGhlIGBDYW5jZWxsYWJsZVByb21pc2VgIGNvbnN0cnVjdG9yLlxuICAgICAqXG4gICAgICogVGhpcyBtZXRob2QgZW5hYmxlcyBgYXdhaXRgaW5nIGNhbmNlbGxhYmxlIHByb21pc2VzIHdpdGhvdXQgaGF2aW5nXG4gICAgICogdG8gc3RvcmUgdGhlbSBmb3IgZnV0dXJlIGNhbmNlbGxhdGlvbiwgZS5nLjpcbiAgICAgKiBgYGB0c1xuICAgICAqIGF3YWl0IGxvbmdSdW5uaW5nT3BlcmF0aW9uKCkuY2FuY2VsT24oc2lnbmFsKTtcbiAgICAgKiBgYGBcbiAgICAgKiBpbnN0ZWFkIG9mOlxuICAgICAqIGBgYHRzXG4gICAgICogbGV0IHByb21pc2VUb0JlQ2FuY2VsbGVkID0gbG9uZ1J1bm5pbmdPcGVyYXRpb24oKTtcbiAgICAgKiBhd2FpdCBwcm9taXNlVG9CZUNhbmNlbGxlZDtcbiAgICAgKiBgYGBcbiAgICAgKlxuICAgICAqIEByZXR1cm5zIFRoaXMgcHJvbWlzZSwgZm9yIG1ldGhvZCBjaGFpbmluZy5cbiAgICAgKi9cbiAgICBjYW5jZWxPbihzaWduYWw6IEFib3J0U2lnbmFsKTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgaWYgKHNpZ25hbC5hYm9ydGVkKSB7XG4gICAgICAgICAgICB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pXG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBzaWduYWwuYWRkRXZlbnRMaXN0ZW5lcignYWJvcnQnLCAoKSA9PiB2b2lkIHRoaXMuY2FuY2VsKHNpZ25hbC5yZWFzb24pLCB7Y2FwdHVyZTogdHJ1ZX0pO1xuICAgICAgICB9XG5cbiAgICAgICAgcmV0dXJuIHRoaXM7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQXR0YWNoZXMgY2FsbGJhY2tzIGZvciB0aGUgcmVzb2x1dGlvbiBhbmQvb3IgcmVqZWN0aW9uIG9mIHRoZSBgQ2FuY2VsbGFibGVQcm9taXNlYC5cbiAgICAgKlxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICpcbiAgICAgKiBAcGFyYW0gb25mdWxmaWxsZWQgVGhlIGNhbGxiYWNrIHRvIGV4ZWN1dGUgd2hlbiB0aGUgUHJvbWlzZSBpcyByZXNvbHZlZC5cbiAgICAgKiBAcGFyYW0gb25yZWplY3RlZCBUaGUgY2FsbGJhY2sgdG8gZXhlY3V0ZSB3aGVuIHRoZSBQcm9taXNlIGlzIHJlamVjdGVkLlxuICAgICAqIEByZXR1cm5zIEEgYENhbmNlbGxhYmxlUHJvbWlzZWAgZm9yIHRoZSBjb21wbGV0aW9uIG9mIHdoaWNoZXZlciBjYWxsYmFjayBpcyBleGVjdXRlZC5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBob29rZWQgdXAgdG8gcHJvcGFnYXRlIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB1cCB0aGUgY2hhaW4sIGJ1dCBub3QgZG93bjpcbiAgICAgKlxuICAgICAqICAgLSBpZiB0aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkLCB0aGUgYG9ucmVqZWN0ZWRgIGhhbmRsZXIgd2lsbCBiZSBpbnZva2VkIHdpdGggYSBgQ2FuY2VsRXJyb3JgXG4gICAgICogICAgIGFuZCB0aGUgcmV0dXJuZWQgcHJvbWlzZSBfd2lsbCByZXNvbHZlIHJlZ3VsYXJseV8gd2l0aCBpdHMgcmVzdWx0O1xuICAgICAqICAgLSBjb252ZXJzZWx5LCBpZiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIF90aGUgcGFyZW50IHByb21pc2UgaXMgY2FuY2VsbGVkIHRvbztfXG4gICAgICogICAgIHRoZSBgb25yZWplY3RlZGAgaGFuZGxlciB3aWxsIHN0aWxsIGJlIGludm9rZWQgd2l0aCB0aGUgcGFyZW50J3MgYENhbmNlbEVycm9yYCxcbiAgICAgKiAgICAgYnV0IGl0cyByZXN1bHQgd2lsbCBiZSBkaXNjYXJkZWRcbiAgICAgKiAgICAgYW5kIHRoZSByZXR1cm5lZCBwcm9taXNlIHdpbGwgcmVqZWN0IHdpdGggYSBgQ2FuY2VsRXJyb3JgIGFzIHdlbGwuXG4gICAgICpcbiAgICAgKiBUaGUgcHJvbWlzZSByZXR1cm5lZCBmcm9tIHtAbGluayBjYW5jZWx9IHdpbGwgZnVsZmlsbCBvbmx5IGFmdGVyIGFsbCBhdHRhY2hlZCBoYW5kbGVyc1xuICAgICAqIHVwIHRoZSBlbnRpcmUgcHJvbWlzZSBjaGFpbiBoYXZlIGJlZW4gcnVuLlxuICAgICAqXG4gICAgICogSWYgZWl0aGVyIGNhbGxiYWNrIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICovXG4gICAgdGhlbjxUUmVzdWx0MSA9IFQsIFRSZXN1bHQyID0gbmV2ZXI+KG9uZnVsZmlsbGVkPzogKCh2YWx1ZTogVCkgPT4gVFJlc3VsdDEgfCBQcm9taXNlTGlrZTxUUmVzdWx0MT4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQxPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbnJlamVjdGVkPzogKChyZWFzb246IGFueSkgPT4gVFJlc3VsdDIgfCBQcm9taXNlTGlrZTxUUmVzdWx0Mj4gfCBDYW5jZWxsYWJsZVByb21pc2VMaWtlPFRSZXN1bHQyPikgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUUmVzdWx0MSB8IFRSZXN1bHQyPiB7XG4gICAgICAgIGlmICghKHRoaXMgaW5zdGFuY2VvZiBDYW5jZWxsYWJsZVByb21pc2UpKSB7XG4gICAgICAgICAgICB0aHJvdyBuZXcgVHlwZUVycm9yKFwiQ2FuY2VsbGFibGVQcm9taXNlLnByb3RvdHlwZS50aGVuIGNhbGxlZCBvbiBhbiBpbnZhbGlkIG9iamVjdC5cIik7XG4gICAgICAgIH1cblxuICAgICAgICAvLyBOT1RFOiBUeXBlU2NyaXB0J3MgYnVpbHQtaW4gdHlwZSBmb3IgdGhlbiBpcyBicm9rZW4sXG4gICAgICAgIC8vIGFzIGl0IGFsbG93cyBzcGVjaWZ5aW5nIGFuIGFyYml0cmFyeSBUUmVzdWx0MSAhPSBUIGV2ZW4gd2hlbiBvbmZ1bGZpbGxlZCBpcyBub3QgYSBmdW5jdGlvbi5cbiAgICAgICAgLy8gV2UgY2Fubm90IGZpeCBpdCBpZiB3ZSB3YW50IHRvIENhbmNlbGxhYmxlUHJvbWlzZSB0byBpbXBsZW1lbnQgUHJvbWlzZUxpa2U8VD4uXG5cbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZnVsZmlsbGVkKSkgeyBvbmZ1bGZpbGxlZCA9IGlkZW50aXR5IGFzIGFueTsgfVxuICAgICAgICBpZiAoIWlzQ2FsbGFibGUob25yZWplY3RlZCkpIHsgb25yZWplY3RlZCA9IHRocm93ZXI7IH1cblxuICAgICAgICBpZiAob25mdWxmaWxsZWQgPT09IGlkZW50aXR5ICYmIG9ucmVqZWN0ZWQgPT0gdGhyb3dlcikge1xuICAgICAgICAgICAgLy8gU2hvcnRjdXQgZm9yIHRyaXZpYWwgYXJndW1lbnRzLlxuICAgICAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2UoKHJlc29sdmUpID0+IHJlc29sdmUodGhpcyBhcyBhbnkpKTtcbiAgICAgICAgfVxuXG4gICAgICAgIGNvbnN0IGJhcnJpZXI6IFBhcnRpYWw8UHJvbWlzZVdpdGhSZXNvbHZlcnM8dm9pZD4+ID0ge307XG4gICAgICAgIHRoaXNbYmFycmllclN5bV0gPSBiYXJyaWVyO1xuXG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFRSZXN1bHQxIHwgVFJlc3VsdDI+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgIHZvaWQgc3VwZXIudGhlbihcbiAgICAgICAgICAgICAgICAodmFsdWUpID0+IHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcblxuICAgICAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbmZ1bGZpbGxlZCEodmFsdWUpKTtcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0sXG4gICAgICAgICAgICAgICAgKHJlYXNvbj8pID0+IHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHRoaXNbYmFycmllclN5bV0gPT09IGJhcnJpZXIpIHsgdGhpc1tiYXJyaWVyU3ltXSA9IG51bGw7IH1cbiAgICAgICAgICAgICAgICAgICAgYmFycmllci5yZXNvbHZlPy4oKTtcblxuICAgICAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmVzb2x2ZShvbnJlamVjdGVkIShyZWFzb24pKTtcbiAgICAgICAgICAgICAgICAgICAgfSBjYXRjaCAoZXJyKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZWplY3QoZXJyKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICk7XG4gICAgICAgIH0sIGFzeW5jIChjYXVzZT8pID0+IHtcbiAgICAgICAgICAgIC8vY2FuY2VsbGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG9uY2FuY2VsbGVkPy4oY2F1c2UpO1xuICAgICAgICAgICAgfSBmaW5hbGx5IHtcbiAgICAgICAgICAgICAgICBhd2FpdCB0aGlzLmNhbmNlbChjYXVzZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEF0dGFjaGVzIGEgY2FsbGJhY2sgZm9yIG9ubHkgdGhlIHJlamVjdGlvbiBvZiB0aGUgUHJvbWlzZS5cbiAgICAgKlxuICAgICAqIFRoZSBvcHRpb25hbCBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IHdpbGwgYmUgaW52b2tlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCxcbiAgICAgKiB3aXRoIHRoZSBzYW1lIHNlbWFudGljcyBhcyB0aGUgYG9uY2FuY2VsbGVkYCBhcmd1bWVudCBvZiB0aGUgY29uc3RydWN0b3IuXG4gICAgICogV2hlbiB0aGUgcGFyZW50IHByb21pc2UgcmVqZWN0cyBvciBpcyBjYW5jZWxsZWQsIHRoZSBgb25yZWplY3RlZGAgY2FsbGJhY2sgd2lsbCBydW4sXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICpcbiAgICAgKiBJdCBpcyBlcXVpdmFsZW50IHRvXG4gICAgICogYGBgdHNcbiAgICAgKiBjYW5jZWxsYWJsZVByb21pc2UudGhlbih1bmRlZmluZWQsIG9ucmVqZWN0ZWQsIG9uY2FuY2VsbGVkKTtcbiAgICAgKiBgYGBcbiAgICAgKiBhbmQgdGhlIHNhbWUgY2F2ZWF0cyBhcHBseS5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxuICAgICAqXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cbiAgICAgKlxuICAgICAqIElmIGBvbnJlamVjdGVkYCByZXR1cm5zIGEgY2FuY2VsbGFibGUgcHJvbWlzZSxcbiAgICAgKiBjYW5jZWxsYXRpb24gcmVxdWVzdHMgd2lsbCBiZSBkaXZlcnRlZCB0byBpdCxcbiAgICAgKiBhbmQgdGhlIHNwZWNpZmllZCBgb25jYW5jZWxsZWRgIGNhbGxiYWNrIHdpbGwgYmUgZGlzY2FyZGVkLlxuICAgICAqIFNlZSB7QGxpbmsgdGhlbn0gZm9yIG1vcmUgZGV0YWlscy5cbiAgICAgKi9cbiAgICBjYXRjaDxUUmVzdWx0ID0gbmV2ZXI+KG9ucmVqZWN0ZWQ/OiAoKHJlYXNvbjogYW55KSA9PiAoUHJvbWlzZUxpa2U8VFJlc3VsdD4gfCBUUmVzdWx0KSkgfCB1bmRlZmluZWQgfCBudWxsLCBvbmNhbmNlbGxlZD86IENhbmNlbGxhYmxlUHJvbWlzZUNhbmNlbGxlcik6IENhbmNlbGxhYmxlUHJvbWlzZTxUIHwgVFJlc3VsdD4ge1xuICAgICAgICByZXR1cm4gdGhpcy50aGVuKHVuZGVmaW5lZCwgb25yZWplY3RlZCwgb25jYW5jZWxsZWQpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIEF0dGFjaGVzIGEgY2FsbGJhY2sgdGhhdCBpcyBpbnZva2VkIHdoZW4gdGhlIENhbmNlbGxhYmxlUHJvbWlzZSBpcyBzZXR0bGVkIChmdWxmaWxsZWQgb3IgcmVqZWN0ZWQpLiBUaGVcbiAgICAgKiByZXNvbHZlZCB2YWx1ZSBjYW5ub3QgYmUgYWNjZXNzZWQgb3IgbW9kaWZpZWQgZnJvbSB0aGUgY2FsbGJhY2suXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2Ugd2lsbCBzZXR0bGUgaW4gdGhlIHNhbWUgc3RhdGUgYXMgdGhlIG9yaWdpbmFsIG9uZVxuICAgICAqIGFmdGVyIHRoZSBwcm92aWRlZCBjYWxsYmFjayBoYXMgY29tcGxldGVkIGV4ZWN1dGlvbixcbiAgICAgKiB1bmxlc3MgdGhlIGNhbGxiYWNrIHRocm93cyBvciByZXR1cm5zIGEgcmVqZWN0aW5nIHByb21pc2UsXG4gICAgICogaW4gd2hpY2ggY2FzZSB0aGUgcmV0dXJuZWQgcHJvbWlzZSB3aWxsIHJlamVjdCBhcyB3ZWxsLlxuICAgICAqXG4gICAgICogVGhlIG9wdGlvbmFsIGBvbmNhbmNlbGxlZGAgYXJndW1lbnQgd2lsbCBiZSBpbnZva2VkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLFxuICAgICAqIHdpdGggdGhlIHNhbWUgc2VtYW50aWNzIGFzIHRoZSBgb25jYW5jZWxsZWRgIGFyZ3VtZW50IG9mIHRoZSBjb25zdHJ1Y3Rvci5cbiAgICAgKiBPbmNlIHRoZSBwYXJlbnQgcHJvbWlzZSBzZXR0bGVzLCB0aGUgYG9uZmluYWxseWAgY2FsbGJhY2sgd2lsbCBydW4sXG4gICAgICogX2V2ZW4gYWZ0ZXIgdGhlIHJldHVybmVkIHByb21pc2UgaGFzIGJlZW4gY2FuY2VsbGVkOl9cbiAgICAgKiBpbiB0aGF0IGNhc2UsIHNob3VsZCBpdCByZWplY3Qgb3IgdGhyb3csIHRoZSByZWFzb24gd2lsbCBiZSB3cmFwcGVkXG4gICAgICogaW4gYSB7QGxpbmsgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3J9IGFuZCBidWJibGVkIHVwIGFzIGFuIHVuaGFuZGxlZCByZWplY3Rpb24uXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBpbXBsZW1lbnRlZCBpbiB0ZXJtcyBvZiB7QGxpbmsgdGhlbn0gYW5kIHRoZSBzYW1lIGNhdmVhdHMgYXBwbHkuXG4gICAgICogSXQgaXMgcG9seWZpbGxlZCwgaGVuY2UgYXZhaWxhYmxlIGluIGV2ZXJ5IE9TL3dlYnZpZXcgdmVyc2lvbi5cbiAgICAgKlxuICAgICAqIEByZXR1cm5zIEEgUHJvbWlzZSBmb3IgdGhlIGNvbXBsZXRpb24gb2YgdGhlIGNhbGxiYWNrLlxuICAgICAqIENhbmNlbGxhdGlvbiByZXF1ZXN0cyBvbiB0aGUgcmV0dXJuZWQgcHJvbWlzZVxuICAgICAqIHdpbGwgcHJvcGFnYXRlIHVwIHRoZSBjaGFpbiB0byB0aGUgcGFyZW50IHByb21pc2UsXG4gICAgICogYnV0IG5vdCBpbiB0aGUgb3RoZXIgZGlyZWN0aW9uLlxuICAgICAqXG4gICAgICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB7QGxpbmsgY2FuY2VsfSB3aWxsIGZ1bGZpbGwgb25seSBhZnRlciBhbGwgYXR0YWNoZWQgaGFuZGxlcnNcbiAgICAgKiB1cCB0aGUgZW50aXJlIHByb21pc2UgY2hhaW4gaGF2ZSBiZWVuIHJ1bi5cbiAgICAgKlxuICAgICAqIElmIGBvbmZpbmFsbHlgIHJldHVybnMgYSBjYW5jZWxsYWJsZSBwcm9taXNlLFxuICAgICAqIGNhbmNlbGxhdGlvbiByZXF1ZXN0cyB3aWxsIGJlIGRpdmVydGVkIHRvIGl0LFxuICAgICAqIGFuZCB0aGUgc3BlY2lmaWVkIGBvbmNhbmNlbGxlZGAgY2FsbGJhY2sgd2lsbCBiZSBkaXNjYXJkZWQuXG4gICAgICogU2VlIHtAbGluayB0aGVufSBmb3IgbW9yZSBkZXRhaWxzLlxuICAgICAqL1xuICAgIGZpbmFsbHkob25maW5hbGx5PzogKCgpID0+IHZvaWQpIHwgdW5kZWZpbmVkIHwgbnVsbCwgb25jYW5jZWxsZWQ/OiBDYW5jZWxsYWJsZVByb21pc2VDYW5jZWxsZXIpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICBpZiAoISh0aGlzIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSkge1xuICAgICAgICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihcIkNhbmNlbGxhYmxlUHJvbWlzZS5wcm90b3R5cGUuZmluYWxseSBjYWxsZWQgb24gYW4gaW52YWxpZCBvYmplY3QuXCIpO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFpc0NhbGxhYmxlKG9uZmluYWxseSkpIHtcbiAgICAgICAgICAgIHJldHVybiB0aGlzLnRoZW4ob25maW5hbGx5LCBvbmZpbmFsbHksIG9uY2FuY2VsbGVkKTtcbiAgICAgICAgfVxuXG4gICAgICAgIHJldHVybiB0aGlzLnRoZW4oXG4gICAgICAgICAgICAodmFsdWUpID0+IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKG9uZmluYWxseSgpKS50aGVuKCgpID0+IHZhbHVlKSxcbiAgICAgICAgICAgIChyZWFzb24/KSA9PiBDYW5jZWxsYWJsZVByb21pc2UucmVzb2x2ZShvbmZpbmFsbHkoKSkudGhlbigoKSA9PiB7IHRocm93IHJlYXNvbjsgfSksXG4gICAgICAgICAgICBvbmNhbmNlbGxlZCxcbiAgICAgICAgKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBXZSB1c2UgdGhlIGBbU3ltYm9sLnNwZWNpZXNdYCBzdGF0aWMgcHJvcGVydHksIGlmIGF2YWlsYWJsZSxcbiAgICAgKiB0byBkaXNhYmxlIHRoZSBidWlsdC1pbiBhdXRvbWF0aWMgc3ViY2xhc3NpbmcgZmVhdHVyZXMgZnJvbSB7QGxpbmsgUHJvbWlzZX0uXG4gICAgICogSXQgaXMgY3JpdGljYWwgZm9yIHBlcmZvcm1hbmNlIHJlYXNvbnMgdGhhdCBleHRlbmRlcnMgZG8gbm90IG92ZXJyaWRlIHRoaXMuXG4gICAgICogT25jZSB0aGUgcHJvcG9zYWwgYXQgaHR0cHM6Ly9naXRodWIuY29tL3RjMzkvcHJvcG9zYWwtcm0tYnVpbHRpbi1zdWJjbGFzc2luZ1xuICAgICAqIGlzIGVpdGhlciBhY2NlcHRlZCBvciByZXRpcmVkLCB0aGlzIGltcGxlbWVudGF0aW9uIHdpbGwgaGF2ZSB0byBiZSByZXZpc2VkIGFjY29yZGluZ2x5LlxuICAgICAqXG4gICAgICogQGlnbm9yZVxuICAgICAqIEBpbnRlcm5hbFxuICAgICAqL1xuICAgIHN0YXRpYyBnZXQgW3NwZWNpZXNdKCkge1xuICAgICAgICByZXR1cm4gUHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgaXMgcmVzb2x2ZWQgd2l0aCBhbiBhcnJheSBvZiByZXN1bHRzXG4gICAgICogd2hlbiBhbGwgb2YgdGhlIHByb3ZpZGVkIFByb21pc2VzIHJlc29sdmUsIG9yIHJlamVjdGVkIHdoZW4gYW55IFByb21pc2UgaXMgcmVqZWN0ZWQuXG4gICAgICpcbiAgICAgKiBFdmVyeSBvbmUgb2YgdGhlIHByb3ZpZGVkIG9iamVjdHMgdGhhdCBpcyBhIHRoZW5hYmxlIF9hbmRfIGNhbmNlbGxhYmxlIG9iamVjdFxuICAgICAqIHdpbGwgYmUgY2FuY2VsbGVkIHdoZW4gdGhlIHJldHVybmVkIHByb21pc2UgaXMgY2FuY2VsbGVkLCB3aXRoIHRoZSBzYW1lIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGFsbDxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPltdPjtcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyByZWFkb25seSB1bmtub3duW10gfCBbXT4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHsgLXJlYWRvbmx5IFtQIGluIGtleW9mIFRdOiBBd2FpdGVkPFRbUF0+OyB9PjtcbiAgICBzdGF0aWMgYWxsPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYWxsKGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IGlzIHJlc29sdmVkIHdpdGggYW4gYXJyYXkgb2YgcmVzdWx0c1xuICAgICAqIHdoZW4gYWxsIG9mIHRoZSBwcm92aWRlZCBQcm9taXNlcyByZXNvbHZlIG9yIHJlamVjdC5cbiAgICAgKlxuICAgICAqIEV2ZXJ5IG9uZSBvZiB0aGUgcHJvdmlkZWQgb2JqZWN0cyB0aGF0IGlzIGEgdGhlbmFibGUgX2FuZF8gY2FuY2VsbGFibGUgb2JqZWN0XG4gICAgICogd2lsbCBiZSBjYW5jZWxsZWQgd2hlbiB0aGUgcmV0dXJuZWQgcHJvbWlzZSBpcyBjYW5jZWxsZWQsIHdpdGggdGhlIHNhbWUgY2F1c2UuXG4gICAgICpcbiAgICAgKiBAZ3JvdXAgU3RhdGljIE1ldGhvZHNcbiAgICAgKi9cbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUPih2YWx1ZXM6IEl0ZXJhYmxlPFQgfCBQcm9taXNlTGlrZTxUPj4pOiBDYW5jZWxsYWJsZVByb21pc2U8UHJvbWlzZVNldHRsZWRSZXN1bHQ8QXdhaXRlZDxUPj5bXT47XG4gICAgc3RhdGljIGFsbFNldHRsZWQ8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8eyAtcmVhZG9ubHkgW1AgaW4ga2V5b2YgVF06IFByb21pc2VTZXR0bGVkUmVzdWx0PEF3YWl0ZWQ8VFtQXT4+OyB9PjtcbiAgICBzdGF0aWMgYWxsU2V0dGxlZDxUIGV4dGVuZHMgSXRlcmFibGU8dW5rbm93bj4gfCBBcnJheUxpa2U8dW5rbm93bj4+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPiB7XG4gICAgICAgIGxldCBjb2xsZWN0ZWQgPSBBcnJheS5mcm9tKHZhbHVlcyk7XG4gICAgICAgIGNvbnN0IHByb21pc2UgPSBjb2xsZWN0ZWQubGVuZ3RoID09PSAwXG4gICAgICAgICAgICA/IENhbmNlbGxhYmxlUHJvbWlzZS5yZXNvbHZlKGNvbGxlY3RlZClcbiAgICAgICAgICAgIDogbmV3IENhbmNlbGxhYmxlUHJvbWlzZTx1bmtub3duPigocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgICAgICAgICAgdm9pZCBQcm9taXNlLmFsbFNldHRsZWQoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgICAgICB9LCAoY2F1c2U/KTogUHJvbWlzZTx2b2lkPiA9PiBjYW5jZWxBbGwocHJvbWlzZSwgY29sbGVjdGVkLCBjYXVzZSkpO1xuICAgICAgICByZXR1cm4gcHJvbWlzZTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBUaGUgYW55IGZ1bmN0aW9uIHJldHVybnMgYSBwcm9taXNlIHRoYXQgaXMgZnVsZmlsbGVkIGJ5IHRoZSBmaXJzdCBnaXZlbiBwcm9taXNlIHRvIGJlIGZ1bGZpbGxlZCxcbiAgICAgKiBvciByZWplY3RlZCB3aXRoIGFuIEFnZ3JlZ2F0ZUVycm9yIGNvbnRhaW5pbmcgYW4gYXJyYXkgb2YgcmVqZWN0aW9uIHJlYXNvbnNcbiAgICAgKiBpZiBhbGwgb2YgdGhlIGdpdmVuIHByb21pc2VzIGFyZSByZWplY3RlZC5cbiAgICAgKiBJdCByZXNvbHZlcyBhbGwgZWxlbWVudHMgb2YgdGhlIHBhc3NlZCBpdGVyYWJsZSB0byBwcm9taXNlcyBhcyBpdCBydW5zIHRoaXMgYWxnb3JpdGhtLlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBhbnk8VD4odmFsdWVzOiBJdGVyYWJsZTxUIHwgUHJvbWlzZUxpa2U8VD4+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+O1xuICAgIHN0YXRpYyBhbnk8VCBleHRlbmRzIHJlYWRvbmx5IHVua25vd25bXSB8IFtdPih2YWx1ZXM6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUW251bWJlcl0+PjtcbiAgICBzdGF0aWMgYW55PFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IGNvbGxlY3RlZC5sZW5ndGggPT09IDBcbiAgICAgICAgICAgID8gQ2FuY2VsbGFibGVQcm9taXNlLnJlc29sdmUoY29sbGVjdGVkKVxuICAgICAgICAgICAgOiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+KChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgICAgICAgICB2b2lkIFByb21pc2UuYW55KGNvbGxlY3RlZCkudGhlbihyZXNvbHZlLCByZWplY3QpO1xuICAgICAgICAgICAgfSwgKGNhdXNlPyk6IFByb21pc2U8dm9pZD4gPT4gY2FuY2VsQWxsKHByb21pc2UsIGNvbGxlY3RlZCwgY2F1c2UpKTtcbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIFByb21pc2UgdGhhdCBpcyByZXNvbHZlZCBvciByZWplY3RlZCB3aGVuIGFueSBvZiB0aGUgcHJvdmlkZWQgUHJvbWlzZXMgYXJlIHJlc29sdmVkIG9yIHJlamVjdGVkLlxuICAgICAqXG4gICAgICogRXZlcnkgb25lIG9mIHRoZSBwcm92aWRlZCBvYmplY3RzIHRoYXQgaXMgYSB0aGVuYWJsZSBfYW5kXyBjYW5jZWxsYWJsZSBvYmplY3RcbiAgICAgKiB3aWxsIGJlIGNhbmNlbGxlZCB3aGVuIHRoZSByZXR1cm5lZCBwcm9taXNlIGlzIGNhbmNlbGxlZCwgd2l0aCB0aGUgc2FtZSBjYXVzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByYWNlPFQ+KHZhbHVlczogSXRlcmFibGU8VCB8IFByb21pc2VMaWtlPFQ+Pik6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFQ+PjtcbiAgICBzdGF0aWMgcmFjZTxUIGV4dGVuZHMgcmVhZG9ubHkgdW5rbm93bltdIHwgW10+KHZhbHVlczogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxBd2FpdGVkPFRbbnVtYmVyXT4+O1xuICAgIHN0YXRpYyByYWNlPFQgZXh0ZW5kcyBJdGVyYWJsZTx1bmtub3duPiB8IEFycmF5TGlrZTx1bmtub3duPj4odmFsdWVzOiBUKTogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+IHtcbiAgICAgICAgbGV0IGNvbGxlY3RlZCA9IEFycmF5LmZyb20odmFsdWVzKTtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8dW5rbm93bj4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgdm9pZCBQcm9taXNlLnJhY2UoY29sbGVjdGVkKS50aGVuKHJlc29sdmUsIHJlamVjdCk7XG4gICAgICAgIH0sIChjYXVzZT8pOiBQcm9taXNlPHZvaWQ+ID0+IGNhbmNlbEFsbChwcm9taXNlLCBjb2xsZWN0ZWQsIGNhdXNlKSk7XG4gICAgICAgIHJldHVybiBwcm9taXNlO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgY2FuY2VsbGVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIGNhdXNlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIGNhbmNlbDxUID0gbmV2ZXI+KGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgY29uc3QgcCA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xuICAgICAgICBwLmNhbmNlbChjYXVzZSk7XG4gICAgICAgIHJldHVybiBwO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIHRoYXQgY2FuY2Vsc1xuICAgICAqIGFmdGVyIHRoZSBzcGVjaWZpZWQgdGltZW91dCwgd2l0aCB0aGUgcHJvdmlkZWQgY2F1c2UuXG4gICAgICpcbiAgICAgKiBJZiB0aGUge0BsaW5rIEFib3J0U2lnbmFsLnRpbWVvdXR9IGZhY3RvcnkgbWV0aG9kIGlzIGF2YWlsYWJsZSxcbiAgICAgKiBpdCBpcyB1c2VkIHRvIGJhc2UgdGhlIHRpbWVvdXQgb24gX2FjdGl2ZV8gdGltZSByYXRoZXIgdGhhbiBfZWxhcHNlZF8gdGltZS5cbiAgICAgKiBPdGhlcndpc2UsIGB0aW1lb3V0YCBmYWxscyBiYWNrIHRvIHtAbGluayBzZXRUaW1lb3V0fS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyB0aW1lb3V0PFQgPSBuZXZlcj4obWlsbGlzZWNvbmRzOiBudW1iZXIsIGNhdXNlPzogYW55KTogQ2FuY2VsbGFibGVQcm9taXNlPFQ+IHtcbiAgICAgICAgY29uc3QgcHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKCkgPT4ge30pO1xuICAgICAgICBpZiAoQWJvcnRTaWduYWwgJiYgdHlwZW9mIEFib3J0U2lnbmFsID09PSAnZnVuY3Rpb24nICYmIEFib3J0U2lnbmFsLnRpbWVvdXQgJiYgdHlwZW9mIEFib3J0U2lnbmFsLnRpbWVvdXQgPT09ICdmdW5jdGlvbicpIHtcbiAgICAgICAgICAgIEFib3J0U2lnbmFsLnRpbWVvdXQobWlsbGlzZWNvbmRzKS5hZGRFdmVudExpc3RlbmVyKCdhYm9ydCcsICgpID0+IHZvaWQgcHJvbWlzZS5jYW5jZWwoY2F1c2UpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gdm9pZCBwcm9taXNlLmNhbmNlbChjYXVzZSksIG1pbGxpc2Vjb25kcyk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHByb21pc2U7XG4gICAgfVxuXG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyBDYW5jZWxsYWJsZVByb21pc2UgdGhhdCByZXNvbHZlcyBhZnRlciB0aGUgc3BlY2lmaWVkIHRpbWVvdXQuXG4gICAgICogVGhlIHJldHVybmVkIHByb21pc2UgY2FuIGJlIGNhbmNlbGxlZCB3aXRob3V0IGNvbnNlcXVlbmNlcy5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyBzbGVlcChtaWxsaXNlY29uZHM6IG51bWJlcik6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IENhbmNlbGxhYmxlUHJvbWlzZSB0aGF0IHJlc29sdmVzIGFmdGVyXG4gICAgICogdGhlIHNwZWNpZmllZCB0aW1lb3V0LCB3aXRoIHRoZSBwcm92aWRlZCB2YWx1ZS5cbiAgICAgKiBUaGUgcmV0dXJuZWQgcHJvbWlzZSBjYW4gYmUgY2FuY2VsbGVkIHdpdGhvdXQgY29uc2VxdWVuY2VzLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHNsZWVwPFQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZTogVCk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPjtcbiAgICBzdGF0aWMgc2xlZXA8VCA9IHZvaWQ+KG1pbGxpc2Vjb25kczogbnVtYmVyLCB2YWx1ZT86IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8VD4ge1xuICAgICAgICByZXR1cm4gbmV3IENhbmNlbGxhYmxlUHJvbWlzZTxUPigocmVzb2x2ZSkgPT4ge1xuICAgICAgICAgICAgc2V0VGltZW91dCgoKSA9PiByZXNvbHZlKHZhbHVlISksIG1pbGxpc2Vjb25kcyk7XG4gICAgICAgIH0pO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgcmVqZWN0ZWQgQ2FuY2VsbGFibGVQcm9taXNlIGZvciB0aGUgcHJvdmlkZWQgcmVhc29uLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlamVjdDxUID0gbmV2ZXI+KHJlYXNvbj86IGFueSk6IENhbmNlbGxhYmxlUHJvbWlzZTxUPiB7XG4gICAgICAgIHJldHVybiBuZXcgQ2FuY2VsbGFibGVQcm9taXNlPFQ+KChfLCByZWplY3QpID0+IHJlamVjdChyZWFzb24pKTtcbiAgICB9XG5cbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZXNvbHZlKCk6IENhbmNlbGxhYmxlUHJvbWlzZTx2b2lkPjtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGEgbmV3IHJlc29sdmVkIENhbmNlbGxhYmxlUHJvbWlzZSBmb3IgdGhlIHByb3ZpZGVkIHZhbHVlLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHJlc29sdmU8VD4odmFsdWU6IFQpOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhIG5ldyByZXNvbHZlZCBDYW5jZWxsYWJsZVByb21pc2UgZm9yIHRoZSBwcm92aWRlZCB2YWx1ZS5cbiAgICAgKlxuICAgICAqIEBncm91cCBTdGF0aWMgTWV0aG9kc1xuICAgICAqL1xuICAgIHN0YXRpYyByZXNvbHZlPFQ+KHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pOiBDYW5jZWxsYWJsZVByb21pc2U8QXdhaXRlZDxUPj47XG4gICAgc3RhdGljIHJlc29sdmU8VCA9IHZvaWQ+KHZhbHVlPzogVCB8IFByb21pc2VMaWtlPFQ+KTogQ2FuY2VsbGFibGVQcm9taXNlPEF3YWl0ZWQ8VD4+IHtcbiAgICAgICAgaWYgKHZhbHVlIGluc3RhbmNlb2YgQ2FuY2VsbGFibGVQcm9taXNlKSB7XG4gICAgICAgICAgICAvLyBPcHRpbWlzZSBmb3IgY2FuY2VsbGFibGUgcHJvbWlzZXMuXG4gICAgICAgICAgICByZXR1cm4gdmFsdWU7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIG5ldyBDYW5jZWxsYWJsZVByb21pc2U8YW55PigocmVzb2x2ZSkgPT4gcmVzb2x2ZSh2YWx1ZSkpO1xuICAgIH1cblxuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYSBuZXcgQ2FuY2VsbGFibGVQcm9taXNlIGFuZCByZXR1cm5zIGl0IGluIGFuIG9iamVjdCwgYWxvbmcgd2l0aCBpdHMgcmVzb2x2ZSBhbmQgcmVqZWN0IGZ1bmN0aW9uc1xuICAgICAqIGFuZCBhIGdldHRlci9zZXR0ZXIgZm9yIHRoZSBjYW5jZWxsYXRpb24gY2FsbGJhY2suXG4gICAgICpcbiAgICAgKiBUaGlzIG1ldGhvZCBpcyBwb2x5ZmlsbGVkLCBoZW5jZSBhdmFpbGFibGUgaW4gZXZlcnkgT1Mvd2VidmlldyB2ZXJzaW9uLlxuICAgICAqXG4gICAgICogQGdyb3VwIFN0YXRpYyBNZXRob2RzXG4gICAgICovXG4gICAgc3RhdGljIHdpdGhSZXNvbHZlcnM8VD4oKTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XG4gICAgICAgIGxldCByZXN1bHQ6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4gPSB7IG9uY2FuY2VsbGVkOiBudWxsIH0gYXMgYW55O1xuICAgICAgICByZXN1bHQucHJvbWlzZSA9IG5ldyBDYW5jZWxsYWJsZVByb21pc2U8VD4oKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICAgICAgcmVzdWx0LnJlc29sdmUgPSByZXNvbHZlO1xuICAgICAgICAgICAgcmVzdWx0LnJlamVjdCA9IHJlamVjdDtcbiAgICAgICAgfSwgKGNhdXNlPzogYW55KSA9PiB7IHJlc3VsdC5vbmNhbmNlbGxlZD8uKGNhdXNlKTsgfSk7XG4gICAgICAgIHJldHVybiByZXN1bHQ7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIGNhbmNlbGxhdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICogVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgcmVzdWx0aW5nIGZ1bmN0aW9uIGRvZXMgbm90IHJlamVjdC5cbiAqL1xuZnVuY3Rpb24gY2FuY2VsbGVyRm9yPFQ+KHByb21pc2U6IENhbmNlbGxhYmxlUHJvbWlzZVdpdGhSZXNvbHZlcnM8VD4sIHN0YXRlOiBDYW5jZWxsYWJsZVByb21pc2VTdGF0ZSkge1xuICAgIGxldCBjYW5jZWxsYXRpb25Qcm9taXNlOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPSB1bmRlZmluZWQ7XG5cbiAgICByZXR1cm4gKHJlYXNvbjogQ2FuY2VsRXJyb3IpOiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD4gPT4ge1xuICAgICAgICBpZiAoIXN0YXRlLnNldHRsZWQpIHtcbiAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgc3RhdGUucmVhc29uID0gcmVhc29uO1xuICAgICAgICAgICAgcHJvbWlzZS5yZWplY3QocmVhc29uKTtcblxuICAgICAgICAgICAgLy8gQXR0YWNoIGFuIGVycm9yIGhhbmRsZXIgdGhhdCBpZ25vcmVzIHRoaXMgc3BlY2lmaWMgcmVqZWN0aW9uIHJlYXNvbiBhbmQgbm90aGluZyBlbHNlLlxuICAgICAgICAgICAgLy8gSW4gdGhlb3J5LCBhIHNhbmUgdW5kZXJseWluZyBpbXBsZW1lbnRhdGlvbiBhdCB0aGlzIHBvaW50XG4gICAgICAgICAgICAvLyBzaG91bGQgYWx3YXlzIHJlamVjdCB3aXRoIG91ciBjYW5jZWxsYXRpb24gcmVhc29uLFxuICAgICAgICAgICAgLy8gaGVuY2UgdGhlIGhhbmRsZXIgd2lsbCBuZXZlciB0aHJvdy5cbiAgICAgICAgICAgIHZvaWQgUHJvbWlzZS5wcm90b3R5cGUudGhlbi5jYWxsKHByb21pc2UucHJvbWlzZSwgdW5kZWZpbmVkLCAoZXJyKSA9PiB7XG4gICAgICAgICAgICAgICAgaWYgKGVyciAhPT0gcmVhc29uKSB7XG4gICAgICAgICAgICAgICAgICAgIHRocm93IGVycjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9KTtcbiAgICAgICAgfVxuXG4gICAgICAgIC8vIElmIHJlYXNvbiBpcyBub3Qgc2V0LCB0aGUgcHJvbWlzZSByZXNvbHZlZCByZWd1bGFybHksIGhlbmNlIHdlIG11c3Qgbm90IGNhbGwgb25jYW5jZWxsZWQuXG4gICAgICAgIC8vIElmIG9uY2FuY2VsbGVkIGlzIHVuc2V0LCBubyBuZWVkIHRvIGdvIGFueSBmdXJ0aGVyLlxuICAgICAgICBpZiAoIXN0YXRlLnJlYXNvbiB8fCAhcHJvbWlzZS5vbmNhbmNlbGxlZCkgeyByZXR1cm47IH1cblxuICAgICAgICBjYW5jZWxsYXRpb25Qcm9taXNlID0gbmV3IFByb21pc2U8dm9pZD4oKHJlc29sdmUpID0+IHtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgcmVzb2x2ZShwcm9taXNlLm9uY2FuY2VsbGVkIShzdGF0ZS5yZWFzb24hLmNhdXNlKSk7XG4gICAgICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgICAgICBQcm9taXNlLnJlamVjdChuZXcgQ2FuY2VsbGVkUmVqZWN0aW9uRXJyb3IocHJvbWlzZS5wcm9taXNlLCBlcnIsIFwiVW5oYW5kbGVkIGV4Y2VwdGlvbiBpbiBvbmNhbmNlbGxlZCBjYWxsYmFjay5cIikpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KS5jYXRjaCgocmVhc29uPykgPT4ge1xuICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHByb21pc2UucHJvbWlzZSwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gb25jYW5jZWxsZWQgY2FsbGJhY2suXCIpKTtcbiAgICAgICAgfSk7XG5cbiAgICAgICAgLy8gVW5zZXQgb25jYW5jZWxsZWQgdG8gcHJldmVudCByZXBlYXRlZCBjYWxscy5cbiAgICAgICAgcHJvbWlzZS5vbmNhbmNlbGxlZCA9IG51bGw7XG5cbiAgICAgICAgcmV0dXJuIGNhbmNlbGxhdGlvblByb21pc2U7XG4gICAgfVxufVxuXG4vKipcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlc29sdXRpb24gYWxnb3JpdGhtIGZvciB0aGUgZ2l2ZW4gY2FuY2VsbGFibGUgcHJvbWlzZS5cbiAqL1xuZnVuY3Rpb24gcmVzb2x2ZXJGb3I8VD4ocHJvbWlzZTogQ2FuY2VsbGFibGVQcm9taXNlV2l0aFJlc29sdmVyczxUPiwgc3RhdGU6IENhbmNlbGxhYmxlUHJvbWlzZVN0YXRlKTogQ2FuY2VsbGFibGVQcm9taXNlUmVzb2x2ZXI8VD4ge1xuICAgIHJldHVybiAodmFsdWUpID0+IHtcbiAgICAgICAgaWYgKHN0YXRlLnJlc29sdmluZykgeyByZXR1cm47IH1cbiAgICAgICAgc3RhdGUucmVzb2x2aW5nID0gdHJ1ZTtcblxuICAgICAgICBpZiAodmFsdWUgPT09IHByb21pc2UucHJvbWlzZSkge1xuICAgICAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHsgcmV0dXJuOyB9XG4gICAgICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgICAgIHByb21pc2UucmVqZWN0KG5ldyBUeXBlRXJyb3IoXCJBIHByb21pc2UgY2Fubm90IGJlIHJlc29sdmVkIHdpdGggaXRzZWxmLlwiKSk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cblxuICAgICAgICBpZiAodmFsdWUgIT0gbnVsbCAmJiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyB8fCB0eXBlb2YgdmFsdWUgPT09ICdmdW5jdGlvbicpKSB7XG4gICAgICAgICAgICBsZXQgdGhlbjogYW55O1xuICAgICAgICAgICAgdHJ5IHtcbiAgICAgICAgICAgICAgICB0aGVuID0gKHZhbHVlIGFzIGFueSkudGhlbjtcbiAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgIHN0YXRlLnNldHRsZWQgPSB0cnVlO1xuICAgICAgICAgICAgICAgIHByb21pc2UucmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICBpZiAoaXNDYWxsYWJsZSh0aGVuKSkge1xuICAgICAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgICAgIGxldCBjYW5jZWwgPSAodmFsdWUgYXMgYW55KS5jYW5jZWw7XG4gICAgICAgICAgICAgICAgICAgIGlmIChpc0NhbGxhYmxlKGNhbmNlbCkpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IG9uY2FuY2VsbGVkID0gKGNhdXNlPzogYW55KSA9PiB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgUmVmbGVjdC5hcHBseShjYW5jZWwsIHZhbHVlLCBbY2F1c2VdKTtcbiAgICAgICAgICAgICAgICAgICAgICAgIH07XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoc3RhdGUucmVhc29uKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gSWYgYWxyZWFkeSBjYW5jZWxsZWQsIHByb3BhZ2F0ZSBjYW5jZWxsYXRpb24uXG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gVGhlIHByb21pc2UgcmV0dXJuZWQgZnJvbSB0aGUgY2FuY2VsbGVyIGFsZ29yaXRobSBkb2VzIG5vdCByZWplY3RcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBzbyBpdCBjYW4gYmUgZGlzY2FyZGVkIHNhZmVseS5cbiAgICAgICAgICAgICAgICAgICAgICAgICAgICB2b2lkIGNhbmNlbGxlckZvcih7IC4uLnByb21pc2UsIG9uY2FuY2VsbGVkIH0sIHN0YXRlKShzdGF0ZS5yZWFzb24pO1xuICAgICAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICBwcm9taXNlLm9uY2FuY2VsbGVkID0gb25jYW5jZWxsZWQ7XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9IGNhdGNoIHt9XG5cbiAgICAgICAgICAgICAgICBjb25zdCBuZXdTdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUgPSB7XG4gICAgICAgICAgICAgICAgICAgIHJvb3Q6IHN0YXRlLnJvb3QsXG4gICAgICAgICAgICAgICAgICAgIHJlc29sdmluZzogZmFsc2UsXG4gICAgICAgICAgICAgICAgICAgIGdldCBzZXR0bGVkKCkgeyByZXR1cm4gdGhpcy5yb290LnNldHRsZWQgfSxcbiAgICAgICAgICAgICAgICAgICAgc2V0IHNldHRsZWQodmFsdWUpIHsgdGhpcy5yb290LnNldHRsZWQgPSB2YWx1ZTsgfSxcbiAgICAgICAgICAgICAgICAgICAgZ2V0IHJlYXNvbigpIHsgcmV0dXJuIHRoaXMucm9vdC5yZWFzb24gfVxuICAgICAgICAgICAgICAgIH07XG5cbiAgICAgICAgICAgICAgICBjb25zdCByZWplY3RvciA9IHJlamVjdG9yRm9yKHByb21pc2UsIG5ld1N0YXRlKTtcbiAgICAgICAgICAgICAgICB0cnkge1xuICAgICAgICAgICAgICAgICAgICBSZWZsZWN0LmFwcGx5KHRoZW4sIHZhbHVlLCBbcmVzb2x2ZXJGb3IocHJvbWlzZSwgbmV3U3RhdGUpLCByZWplY3Rvcl0pO1xuICAgICAgICAgICAgICAgIH0gY2F0Y2ggKGVycikge1xuICAgICAgICAgICAgICAgICAgICByZWplY3RvcihlcnIpO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICByZXR1cm47IC8vIElNUE9SVEFOVCFcbiAgICAgICAgICAgIH1cbiAgICAgICAgfVxuXG4gICAgICAgIGlmIChzdGF0ZS5zZXR0bGVkKSB7IHJldHVybjsgfVxuICAgICAgICBzdGF0ZS5zZXR0bGVkID0gdHJ1ZTtcbiAgICAgICAgcHJvbWlzZS5yZXNvbHZlKHZhbHVlKTtcbiAgICB9O1xufVxuXG4vKipcbiAqIFJldHVybnMgYSBjYWxsYmFjayB0aGF0IGltcGxlbWVudHMgdGhlIHJlamVjdGlvbiBhbGdvcml0aG0gZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLlxuICovXG5mdW5jdGlvbiByZWplY3RvckZvcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2VXaXRoUmVzb2x2ZXJzPFQ+LCBzdGF0ZTogQ2FuY2VsbGFibGVQcm9taXNlU3RhdGUpOiBDYW5jZWxsYWJsZVByb21pc2VSZWplY3RvciB7XG4gICAgcmV0dXJuIChyZWFzb24/KSA9PiB7XG4gICAgICAgIGlmIChzdGF0ZS5yZXNvbHZpbmcpIHsgcmV0dXJuOyB9XG4gICAgICAgIHN0YXRlLnJlc29sdmluZyA9IHRydWU7XG5cbiAgICAgICAgaWYgKHN0YXRlLnNldHRsZWQpIHtcbiAgICAgICAgICAgIHRyeSB7XG4gICAgICAgICAgICAgICAgaWYgKHJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIHN0YXRlLnJlYXNvbiBpbnN0YW5jZW9mIENhbmNlbEVycm9yICYmIE9iamVjdC5pcyhyZWFzb24uY2F1c2UsIHN0YXRlLnJlYXNvbi5jYXVzZSkpIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gU3dhbGxvdyBsYXRlIHJlamVjdGlvbnMgdGhhdCBhcmUgQ2FuY2VsRXJyb3JzIHdob3NlIGNhbmNlbGxhdGlvbiBjYXVzZSBpcyB0aGUgc2FtZSBhcyBvdXJzLlxuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfSBjYXRjaCB7fVxuXG4gICAgICAgICAgICB2b2lkIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwcm9taXNlLnByb21pc2UsIHJlYXNvbikpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgc3RhdGUuc2V0dGxlZCA9IHRydWU7XG4gICAgICAgICAgICBwcm9taXNlLnJlamVjdChyZWFzb24pO1xuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIENhbmNlbHMgYWxsIHZhbHVlcyBpbiBhbiBhcnJheSB0aGF0IGxvb2sgbGlrZSBjYW5jZWxsYWJsZSB0aGVuYWJsZXMuXG4gKiBSZXR1cm5zIGEgcHJvbWlzZSB0aGF0IGZ1bGZpbGxzIG9uY2UgYWxsIGNhbmNlbGxhdGlvbiBwcm9jZWR1cmVzIGZvciB0aGUgZ2l2ZW4gdmFsdWVzIGhhdmUgc2V0dGxlZC5cbiAqL1xuZnVuY3Rpb24gY2FuY2VsQWxsKHBhcmVudDogQ2FuY2VsbGFibGVQcm9taXNlPHVua25vd24+LCB2YWx1ZXM6IGFueVtdLCBjYXVzZT86IGFueSk6IFByb21pc2U8dm9pZD4ge1xuICAgIGNvbnN0IHJlc3VsdHM6IFByb21pc2U8dm9pZD5bXSA9IFtdO1xuXG4gICAgZm9yIChjb25zdCB2YWx1ZSBvZiB2YWx1ZXMpIHtcbiAgICAgICAgbGV0IGNhbmNlbDogQ2FuY2VsbGFibGVQcm9taXNlQ2FuY2VsbGVyO1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgaWYgKCFpc0NhbGxhYmxlKHZhbHVlLnRoZW4pKSB7IGNvbnRpbnVlOyB9XG4gICAgICAgICAgICBjYW5jZWwgPSB2YWx1ZS5jYW5jZWw7XG4gICAgICAgICAgICBpZiAoIWlzQ2FsbGFibGUoY2FuY2VsKSkgeyBjb250aW51ZTsgfVxuICAgICAgICB9IGNhdGNoIHsgY29udGludWU7IH1cblxuICAgICAgICBsZXQgcmVzdWx0OiB2b2lkIHwgUHJvbWlzZUxpa2U8dm9pZD47XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICByZXN1bHQgPSBSZWZsZWN0LmFwcGx5KGNhbmNlbCwgdmFsdWUsIFtjYXVzZV0pO1xuICAgICAgICB9IGNhdGNoIChlcnIpIHtcbiAgICAgICAgICAgIFByb21pc2UucmVqZWN0KG5ldyBDYW5jZWxsZWRSZWplY3Rpb25FcnJvcihwYXJlbnQsIGVyciwgXCJVbmhhbmRsZWQgZXhjZXB0aW9uIGluIGNhbmNlbCBtZXRob2QuXCIpKTtcbiAgICAgICAgICAgIGNvbnRpbnVlO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKCFyZXN1bHQpIHsgY29udGludWU7IH1cbiAgICAgICAgcmVzdWx0cy5wdXNoKFxuICAgICAgICAgICAgKHJlc3VsdCBpbnN0YW5jZW9mIFByb21pc2UgID8gcmVzdWx0IDogUHJvbWlzZS5yZXNvbHZlKHJlc3VsdCkpLmNhdGNoKChyZWFzb24/KSA9PiB7XG4gICAgICAgICAgICAgICAgUHJvbWlzZS5yZWplY3QobmV3IENhbmNlbGxlZFJlamVjdGlvbkVycm9yKHBhcmVudCwgcmVhc29uLCBcIlVuaGFuZGxlZCByZWplY3Rpb24gaW4gY2FuY2VsIG1ldGhvZC5cIikpO1xuICAgICAgICAgICAgfSlcbiAgICAgICAgKTtcbiAgICB9XG5cbiAgICByZXR1cm4gUHJvbWlzZS5hbGwocmVzdWx0cykgYXMgYW55O1xufVxuXG4vKipcbiAqIFJldHVybnMgaXRzIGFyZ3VtZW50LlxuICovXG5mdW5jdGlvbiBpZGVudGl0eTxUPih4OiBUKTogVCB7XG4gICAgcmV0dXJuIHg7XG59XG5cbi8qKlxuICogVGhyb3dzIGl0cyBhcmd1bWVudC5cbiAqL1xuZnVuY3Rpb24gdGhyb3dlcihyZWFzb24/OiBhbnkpOiBuZXZlciB7XG4gICAgdGhyb3cgcmVhc29uO1xufVxuXG4vKipcbiAqIEF0dGVtcHRzIHZhcmlvdXMgc3RyYXRlZ2llcyB0byBjb252ZXJ0IGFuIGVycm9yIHRvIGEgc3RyaW5nLlxuICovXG5mdW5jdGlvbiBlcnJvck1lc3NhZ2UoZXJyOiBhbnkpOiBzdHJpbmcge1xuICAgIHRyeSB7XG4gICAgICAgIGlmIChlcnIgaW5zdGFuY2VvZiBFcnJvciB8fCB0eXBlb2YgZXJyICE9PSAnb2JqZWN0JyB8fCBlcnIudG9TdHJpbmcgIT09IE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcpIHtcbiAgICAgICAgICAgIHJldHVybiBcIlwiICsgZXJyO1xuICAgICAgICB9XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIEpTT04uc3RyaW5naWZ5KGVycik7XG4gICAgfSBjYXRjaCB7fVxuXG4gICAgdHJ5IHtcbiAgICAgICAgcmV0dXJuIE9iamVjdC5wcm90b3R5cGUudG9TdHJpbmcuY2FsbChlcnIpO1xuICAgIH0gY2F0Y2gge31cblxuICAgIHJldHVybiBcIjxjb3VsZCBub3QgY29udmVydCBlcnJvciB0byBzdHJpbmc+XCI7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgY3VycmVudCBiYXJyaWVyIHByb21pc2UgZm9yIHRoZSBnaXZlbiBjYW5jZWxsYWJsZSBwcm9taXNlLiBJZiBuZWNlc3NhcnksIGluaXRpYWxpc2VzIHRoZSBiYXJyaWVyLlxuICovXG5mdW5jdGlvbiBjdXJyZW50QmFycmllcjxUPihwcm9taXNlOiBDYW5jZWxsYWJsZVByb21pc2U8VD4pOiBQcm9taXNlPHZvaWQ+IHtcbiAgICBsZXQgcHdyOiBQYXJ0aWFsPFByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+PiA9IHByb21pc2VbYmFycmllclN5bV0gPz8ge307XG4gICAgaWYgKCEoJ3Byb21pc2UnIGluIHB3cikpIHtcbiAgICAgICAgT2JqZWN0LmFzc2lnbihwd3IsIHByb21pc2VXaXRoUmVzb2x2ZXJzPHZvaWQ+KCkpO1xuICAgIH1cbiAgICBpZiAocHJvbWlzZVtiYXJyaWVyU3ltXSA9PSBudWxsKSB7XG4gICAgICAgIHB3ci5yZXNvbHZlISgpO1xuICAgICAgICBwcm9taXNlW2JhcnJpZXJTeW1dID0gcHdyO1xuICAgIH1cbiAgICByZXR1cm4gcHdyLnByb21pc2UhO1xufVxuXG4vLyBQb2x5ZmlsbCBQcm9taXNlLndpdGhSZXNvbHZlcnMuXG5sZXQgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBQcm9taXNlLndpdGhSZXNvbHZlcnM7XG5pZiAocHJvbWlzZVdpdGhSZXNvbHZlcnMgJiYgdHlwZW9mIHByb21pc2VXaXRoUmVzb2x2ZXJzID09PSAnZnVuY3Rpb24nKSB7XG4gICAgcHJvbWlzZVdpdGhSZXNvbHZlcnMgPSBwcm9taXNlV2l0aFJlc29sdmVycy5iaW5kKFByb21pc2UpO1xufSBlbHNlIHtcbiAgICBwcm9taXNlV2l0aFJlc29sdmVycyA9IGZ1bmN0aW9uIDxUPigpOiBQcm9taXNlV2l0aFJlc29sdmVyczxUPiB7XG4gICAgICAgIGxldCByZXNvbHZlITogKHZhbHVlOiBUIHwgUHJvbWlzZUxpa2U8VD4pID0+IHZvaWQ7XG4gICAgICAgIGxldCByZWplY3QhOiAocmVhc29uPzogYW55KSA9PiB2b2lkO1xuICAgICAgICBjb25zdCBwcm9taXNlID0gbmV3IFByb21pc2U8VD4oKHJlcywgcmVqKSA9PiB7IHJlc29sdmUgPSByZXM7IHJlamVjdCA9IHJlajsgfSk7XG4gICAgICAgIHJldHVybiB7IHByb21pc2UsIHJlc29sdmUsIHJlamVjdCB9O1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyLCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZS5qc1wiO1xuXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5DbGlwYm9hcmQpO1xuXG5jb25zdCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcbmNvbnN0IENsaXBib2FyZFRleHQgPSAxO1xuXG4vKipcbiAqIFNldHMgdGhlIHRleHQgdG8gdGhlIENsaXBib2FyZC5cbiAqXG4gKiBAcGFyYW0gdGV4dCAtIFRoZSB0ZXh0IHRvIGJlIHNldCB0byB0aGUgQ2xpcGJvYXJkLlxuICogQHJldHVybiBBIFByb21pc2UgdGhhdCByZXNvbHZlcyB3aGVuIHRoZSBvcGVyYXRpb24gaXMgc3VjY2Vzc2Z1bC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dDogc3RyaW5nKTogUHJvbWlzZTx2b2lkPiB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKlxuICogQHJldHVybnMgQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgd2l0aCB0aGUgdGV4dCBmcm9tIHRoZSBDbGlwYm9hcmQuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBUZXh0KCk6IFByb21pc2U8c3RyaW5nPiB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkVGV4dCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbmV4cG9ydCBpbnRlcmZhY2UgU2l6ZSB7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXG4gICAgV2lkdGg6IG51bWJlcjtcbiAgICAvKiogVGhlIGhlaWdodCBvZiBhIHJlY3Rhbmd1bGFyIGFyZWEuICovXG4gICAgSGVpZ2h0OiBudW1iZXI7XG59XG5cbmV4cG9ydCBpbnRlcmZhY2UgUmVjdCB7XG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cbiAgICBYOiBudW1iZXI7XG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIG9yaWdpbi4gKi9cbiAgICBZOiBudW1iZXI7XG4gICAgLyoqIFRoZSB3aWR0aCBvZiB0aGUgcmVjdGFuZ2xlLiAqL1xuICAgIFdpZHRoOiBudW1iZXI7XG4gICAgLyoqIFRoZSBoZWlnaHQgb2YgdGhlIHJlY3RhbmdsZS4gKi9cbiAgICBIZWlnaHQ6IG51bWJlcjtcbn1cblxuZXhwb3J0IGludGVyZmFjZSBTY3JlZW4ge1xuICAgIC8qKiBVbmlxdWUgaWRlbnRpZmllciBmb3IgdGhlIHNjcmVlbi4gKi9cbiAgICBJRDogc3RyaW5nO1xuICAgIC8qKiBIdW1hbi1yZWFkYWJsZSBuYW1lIG9mIHRoZSBzY3JlZW4uICovXG4gICAgTmFtZTogc3RyaW5nO1xuICAgIC8qKiBUaGUgc2NhbGUgZmFjdG9yIG9mIHRoZSBzY3JlZW4gKERQSS85NikuIDEgPSBzdGFuZGFyZCBEUEksIDIgPSBIaURQSSAoUmV0aW5hKSwgZXRjLiAqL1xuICAgIFNjYWxlRmFjdG9yOiBudW1iZXI7XG4gICAgLyoqIFRoZSBYIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBYOiBudW1iZXI7XG4gICAgLyoqIFRoZSBZIGNvb3JkaW5hdGUgb2YgdGhlIHNjcmVlbi4gKi9cbiAgICBZOiBudW1iZXI7XG4gICAgLyoqIENvbnRhaW5zIHRoZSB3aWR0aCBhbmQgaGVpZ2h0IG9mIHRoZSBzY3JlZW4uICovXG4gICAgU2l6ZTogU2l6ZTtcbiAgICAvKiogQ29udGFpbnMgdGhlIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0LiAqL1xuICAgIEJvdW5kczogUmVjdDtcbiAgICAvKiogQ29udGFpbnMgdGhlIHBoeXNpY2FsIGJvdW5kcyBvZiB0aGUgc2NyZWVuIGluIHRlcm1zIG9mIFgsIFksIFdpZHRoLCBhbmQgSGVpZ2h0IChiZWZvcmUgc2NhbGluZykuICovXG4gICAgUGh5c2ljYWxCb3VuZHM6IFJlY3Q7XG4gICAgLyoqIENvbnRhaW5zIHRoZSBhcmVhIG9mIHRoZSBzY3JlZW4gdGhhdCBpcyBhY3R1YWxseSB1c2FibGUgKGV4Y2x1ZGluZyB0YXNrYmFyIGFuZCBvdGhlciBzeXN0ZW0gVUkpLiAqL1xuICAgIFdvcmtBcmVhOiBSZWN0O1xuICAgIC8qKiBDb250YWlucyB0aGUgcGh5c2ljYWwgV29ya0FyZWEgb2YgdGhlIHNjcmVlbiAoYmVmb3JlIHNjYWxpbmcpLiAqL1xuICAgIFBoeXNpY2FsV29ya0FyZWE6IFJlY3Q7XG4gICAgLyoqIFRydWUgaWYgdGhpcyBpcyB0aGUgcHJpbWFyeSBtb25pdG9yIHNlbGVjdGVkIGJ5IHRoZSB1c2VyIGluIHRoZSBvcGVyYXRpbmcgc3lzdGVtLiAqL1xuICAgIElzUHJpbWFyeTogYm9vbGVhbjtcbiAgICAvKiogVGhlIHJvdGF0aW9uIG9mIHRoZSBzY3JlZW4uICovXG4gICAgUm90YXRpb246IG51bWJlcjtcbn1cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihvYmplY3ROYW1lcy5TY3JlZW5zKTtcblxuY29uc3QgZ2V0QWxsID0gMDtcbmNvbnN0IGdldFByaW1hcnkgPSAxO1xuY29uc3QgZ2V0Q3VycmVudCA9IDI7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpOiBQcm9taXNlPFNjcmVlbltdPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0QWxsKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqXG4gKiBAcmV0dXJucyBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byB0aGUgcHJpbWFyeSBzY3JlZW4uXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCk6IFByb21pc2U8U2NyZWVuPiB7XG4gICAgcmV0dXJuIGNhbGwoZ2V0UHJpbWFyeSk7XG59XG5cbi8qKlxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxuICpcbiAqIEByZXR1cm5zIEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKTogUHJvbWlzZTxTY3JlZW4+IHtcbiAgICByZXR1cm4gY2FsbChnZXRDdXJyZW50KTtcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlciwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lLmpzXCI7XG5cbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdE5hbWVzLklPUyk7XG5cbi8vIE1ldGhvZCBJRHNcbmNvbnN0IEhhcHRpY3NJbXBhY3QgPSAwO1xuY29uc3QgRGV2aWNlSW5mbyA9IDE7XG5cbmV4cG9ydCBuYW1lc3BhY2UgSGFwdGljcyB7XG4gICAgZXhwb3J0IHR5cGUgSW1wYWN0U3R5bGUgPSBcImxpZ2h0XCJ8XCJtZWRpdW1cInxcImhlYXZ5XCJ8XCJzb2Z0XCJ8XCJyaWdpZFwiO1xuICAgIGV4cG9ydCBmdW5jdGlvbiBJbXBhY3Qoc3R5bGU6IEltcGFjdFN0eWxlID0gXCJtZWRpdW1cIik6IFByb21pc2U8dm9pZD4ge1xuICAgICAgICByZXR1cm4gY2FsbChIYXB0aWNzSW1wYWN0LCB7IHN0eWxlIH0pO1xuICAgIH1cbn1cblxuZXhwb3J0IG5hbWVzcGFjZSBEZXZpY2Uge1xuICAgIGV4cG9ydCBpbnRlcmZhY2UgSW5mbyB7XG4gICAgICAgIG1vZGVsOiBzdHJpbmc7XG4gICAgICAgIHN5c3RlbU5hbWU6IHN0cmluZztcbiAgICAgICAgc3lzdGVtVmVyc2lvbjogc3RyaW5nO1xuICAgICAgICBpc1NpbXVsYXRvcjogYm9vbGVhbjtcbiAgICB9XG4gICAgZXhwb3J0IGZ1bmN0aW9uIEluZm8oKTogUHJvbWlzZTxJbmZvPiB7XG4gICAgICAgIHJldHVybiBjYWxsKERldmljZUluZm8pO1xuICAgIH1cbn1cbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNBQTtBQUFBO0FBQUE7QUFBQTs7O0FDNkJBLElBQU0sY0FDRjtBQUVHLFNBQVMsT0FBTyxPQUFlLElBQVk7QUFDOUMsTUFBSSxLQUFLO0FBRVQsTUFBSSxJQUFJLE9BQU87QUFDZixTQUFPLEtBQUs7QUFFUixVQUFNLFlBQWEsS0FBSyxPQUFPLElBQUks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O0FBQUE7QUFZQSxJQUFNQyxTQUFPLGlCQUFpQixZQUFZLEdBQUc7QUFHN0MsSUFBTSxnQkFBZ0I7QUFDdEIsSUFBTSxhQUFhO0FBRVosSUFBVTtBQUFBLENBQVYsQ0FBVUMsYUFBVjtBQUVJLFdBQVMsT0FBTyxRQUFxQixVQUF5QjtBQUNqRSxXQUFPRCxPQUFLLGVBQWUsRUFBRSxNQUFNLENBQUM7QUFBQSxFQUN4QztBQUZPLEVBQUFDLFNBQVM7QUFBQSxHQUZIO0FBT1YsSUFBVTtBQUFBLENBQVYsQ0FBVUMsWUFBVjtBQU9JLFdBQVNDLFFBQXNCO0FBQ2xDLFdBQU9ILE9BQUssVUFBVTtBQUFBLEVBQzFCO0FBRk8sRUFBQUUsUUFBUyxPQUFBQztBQUFBLEdBUEg7OztBdkJkakIsT0FBTyxTQUFTLE9BQU8sVUFBVSxDQUFDO0FBd0RsQyxPQUFPLE9BQU8sU0FBZ0I7QUFDOUIsT0FBTyxPQUFPLFdBQVc7QUFLekIsT0FBTyxPQUFPLHlCQUF5QixlQUFPLHVCQUF1QixLQUFLLGNBQU07QUFHaEYsT0FBTyxPQUFPLGtCQUFrQjtBQUNoQyxPQUFPLE9BQU8sa0JBQWtCO0FBQ2hDLE9BQU8sT0FBTyxpQkFBaUI7QUFFeEIsT0FBTyxxQkFBcUI7QUFJbkMsTUFBTSxvQkFBb0IsRUFBRSxRQUFRLE9BQU8sQ0FBQyxFQUN2QyxLQUFLLGNBQVk7QUFDZCxNQUFJLFNBQVMsSUFBSTtBQUNiLFVBQU0sU0FBUyxTQUFTLGNBQWMsUUFBUTtBQUM5QyxXQUFPLE1BQU07QUFDYixhQUFTLEtBQUssWUFBWSxNQUFNO0FBQUEsRUFDcEM7QUFDSixDQUFDLEVBQ0EsTUFBTSxNQUFNO0FBQUMsQ0FBQzsiLAogICJuYW1lcyI6IFsiX2EiLCAiRXJyb3IiLCAiY2FsbCIsICJFcnJvciIsICJfYSIsICJBcnJheSIsICJNYXAiLCAiQXJyYXkiLCAiTWFwIiwgImtleSIsICJjYWxsIiwgIl9hIiwgIl9hIiwgInJlc2l6YWJsZSIsICJfYSIsICJjYWxsIiwgIl9hIiwgImNhbGwiLCAiX2EiLCAiY2FsbCIsICJIaWRlTWV0aG9kIiwgIlNob3dNZXRob2QiLCAiaXNEb2N1bWVudERvdEFsbCIsICJfYSIsICJyZWFzb24iLCAidmFsdWUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJIYXB0aWNzIiwgIkRldmljZSIsICJJbmZvIl0KfQo=
