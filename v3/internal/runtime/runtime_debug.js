(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

  // desktop/@wailsio/runtime/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    SetText: () => SetText,
    Text: () => Text
  });

  // node_modules/nanoid/non-secure/index.js
  var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
  var nanoid = (size2 = 21) => {
    let id = "";
    let i = size2;
    while (i--) {
      id += urlAlphabet[Math.random() * 64 | 0];
    }
    return id;
  };

  // desktop/@wailsio/runtime/runtime.js
  var runtimeURL = window.location.origin + "/wails/runtime";
  var objectNames = {
    Call: 0,
    Clipboard: 1,
    Application: 2,
    Events: 3,
    ContextMenu: 4,
    Dialog: 5,
    Window: 6,
    Screens: 7,
    System: 8,
    Browser: 9
  };
  var clientId = nanoid();
  function newRuntimeCallerWithID(object, windowName) {
    return function(method, args = null) {
      return runtimeCallWithID(object, method, windowName, args);
    };
  }
  function runtimeCallWithID(objectID, method, windowName, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("object", objectID);
    url.searchParams.append("method", method);
    let fetchOptions = {
      headers: {}
    };
    if (windowName) {
      fetchOptions.headers["x-wails-window-name"] = windowName;
    }
    if (args) {
      url.searchParams.append("args", JSON.stringify(args));
    }
    fetchOptions.headers["x-wails-client-id"] = clientId;
    return new Promise((resolve, reject) => {
      fetch(url, fetchOptions).then((response) => {
        if (response.ok) {
          if (response.headers.get("Content-Type") && response.headers.get("Content-Type").indexOf("application/json") !== -1) {
            return response.json();
          } else {
            return response.text();
          }
        }
        reject(Error(response.statusText));
      }).then((data) => resolve(data)).catch((error) => reject(error));
    });
  }

  // desktop/@wailsio/runtime/clipboard.js
  var call = newRuntimeCallerWithID(objectNames.Clipboard, "");
  var ClipboardSetText = 0;
  var ClipboardText = 1;
  function SetText(text) {
    return call(ClipboardSetText, { text });
  }
  function Text() {
    return call(ClipboardText);
  }

  // desktop/@wailsio/runtime/application.js
  var application_exports = {};
  __export(application_exports, {
    Hide: () => Hide,
    Quit: () => Quit,
    Show: () => Show
  });
  var call2 = newRuntimeCallerWithID(objectNames.Application);
  var HideMethod = 0;
  var ShowMethod = 1;
  var QuitMethod = 2;
  function Hide() {
    return call2(HideMethod);
  }
  function Show() {
    return call2(ShowMethod);
  }
  function Quit() {
    return call2(QuitMethod);
  }

  // desktop/@wailsio/runtime/screens.js
  var screens_exports = {};
  __export(screens_exports, {
    GetAll: () => GetAll,
    GetCurrent: () => GetCurrent,
    GetPrimary: () => GetPrimary
  });
  var call3 = newRuntimeCallerWithID(objectNames.Screens, "");
  var getAll = 0;
  var getPrimary = 1;
  var getCurrent = 2;
  function GetAll() {
    return call3(getAll);
  }
  function GetPrimary() {
    return call3(getPrimary);
  }
  function GetCurrent() {
    return call3(getCurrent);
  }

  // desktop/@wailsio/runtime/system.js
  var system_exports = {};
  __export(system_exports, {
    Capabilities: () => Capabilities,
    Environment: () => Environment,
    IsAMD64: () => IsAMD64,
    IsARM: () => IsARM,
    IsARM64: () => IsARM64,
    IsDarkMode: () => IsDarkMode,
    IsLinux: () => IsLinux,
    IsMac: () => IsMac,
    IsWindows: () => IsWindows,
    invoke: () => invoke
  });
  var call4 = newRuntimeCallerWithID(objectNames.System, "");
  var systemIsDarkMode = 0;
  var environment = 1;
  function IsDarkMode() {
    return call4(systemIsDarkMode);
  }
  async function Capabilities() {
    let response = fetch("/wails/capabilities");
    return response.json();
  }
  function Environment() {
    return call4(environment);
  }
  var invoke = null;
  var environmentCache = null;
  Environment().then((result) => {
    environmentCache = result;
    invoke = IsWindows() ? window.chrome.webview.postMessage : window.webkit.messageHandlers.external.postMessage;
  }).catch((error) => {
    console.error(`Error getting Environment: ${error}`);
  });
  function IsWindows() {
    return environmentCache.OS === "windows";
  }
  function IsLinux() {
    return environmentCache.OS === "linux";
  }
  function IsMac() {
    return environmentCache.OS === "darwin";
  }
  function IsAMD64() {
    return environmentCache.Arch === "amd64";
  }
  function IsARM() {
    return environmentCache.Arch === "arm";
  }
  function IsARM64() {
    return environmentCache.Arch === "arm64";
  }

  // desktop/@wailsio/runtime/browser.js
  var browser_exports = {};
  __export(browser_exports, {
    OpenURL: () => OpenURL
  });
  var call5 = newRuntimeCallerWithID(objectNames.Browser, "");
  var BrowserOpenURL = 0;
  function OpenURL(url) {
    return call5(BrowserOpenURL, { url });
  }

  // desktop/@wailsio/runtime/window.js
  var center = 0;
  var setTitle = 1;
  var fullscreen = 2;
  var unFullscreen = 3;
  var setSize = 4;
  var size = 5;
  var setMaxSize = 6;
  var setMinSize = 7;
  var setAlwaysOnTop = 8;
  var setRelativePosition = 9;
  var relativePosition = 10;
  var screen = 11;
  var hide = 12;
  var maximise = 13;
  var unMaximise = 14;
  var toggleMaximise = 15;
  var minimise = 16;
  var unMinimise = 17;
  var restore = 18;
  var show = 19;
  var close = 20;
  var setBackgroundColour = 21;
  var setResizable = 22;
  var width = 23;
  var height = 24;
  var zoomIn = 25;
  var zoomOut = 26;
  var zoomReset = 27;
  var getZoomLevel = 28;
  var setZoomLevel = 29;
  var thisWindow = newRuntimeCallerWithID(objectNames.Window, "");
  function createWindow(call10) {
    return {
      Get: (windowName) => createWindow(newRuntimeCallerWithID(objectNames.Window, windowName)),
      Center: () => call10(center),
      SetTitle: (title) => call10(setTitle, { title }),
      Fullscreen: () => call10(fullscreen),
      UnFullscreen: () => call10(unFullscreen),
      SetSize: (width2, height2) => call10(setSize, { width: width2, height: height2 }),
      Size: () => call10(size),
      SetMaxSize: (width2, height2) => call10(setMaxSize, { width: width2, height: height2 }),
      SetMinSize: (width2, height2) => call10(setMinSize, { width: width2, height: height2 }),
      SetAlwaysOnTop: (onTop) => call10(setAlwaysOnTop, { alwaysOnTop: onTop }),
      SetRelativePosition: (x, y) => call10(setRelativePosition, { x, y }),
      RelativePosition: () => call10(relativePosition),
      Screen: () => call10(screen),
      Hide: () => call10(hide),
      Maximise: () => call10(maximise),
      UnMaximise: () => call10(unMaximise),
      ToggleMaximise: () => call10(toggleMaximise),
      Minimise: () => call10(minimise),
      UnMinimise: () => call10(unMinimise),
      Restore: () => call10(restore),
      Show: () => call10(show),
      Close: () => call10(close),
      SetBackgroundColour: (r, g, b, a) => call10(setBackgroundColour, { r, g, b, a }),
      SetResizable: (resizable2) => call10(setResizable, { resizable: resizable2 }),
      Width: () => call10(width),
      Height: () => call10(height),
      ZoomIn: () => call10(zoomIn),
      ZoomOut: () => call10(zoomOut),
      ZoomReset: () => call10(zoomReset),
      GetZoomLevel: () => call10(getZoomLevel),
      SetZoomLevel: (zoomLevel) => call10(setZoomLevel, { zoomLevel })
    };
  }
  function Get(windowName) {
    return createWindow(newRuntimeCallerWithID(objectNames.Window, windowName));
  }
  function WindowMethods(targetWindow) {
    let result = /* @__PURE__ */ new Map();
    for (let method in targetWindow) {
      if (typeof targetWindow[method] === "function") {
        result.set(method, targetWindow[method]);
      }
    }
    return result;
  }
  var window_default = {
    ...Get("")
  };

  // desktop/@wailsio/runtime/calls.js
  var CallBinding = 0;
  var call6 = newRuntimeCallerWithID(objectNames.Call, "");
  var callResponses = /* @__PURE__ */ new Map();
  window._wails = window._wails || {};
  window._wails.callCallback = resultHandler;
  window._wails.callErrorCallback = errorHandler;
  function generateID() {
    let result;
    do {
      result = nanoid();
    } while (callResponses.has(result));
    return result;
  }
  function resultHandler(id, data, isJSON) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
      promiseHandler.resolve(isJSON ? JSON.parse(data) : data);
    }
  }
  function errorHandler(id, message) {
    const promiseHandler = getAndDeleteResponse(id);
    if (promiseHandler) {
      promiseHandler.reject(message);
    }
  }
  function getAndDeleteResponse(id) {
    const response = callResponses.get(id);
    callResponses.delete(id);
    return response;
  }
  function callBinding(type, options = {}) {
    return new Promise((resolve, reject) => {
      const id = generateID();
      options["call-id"] = id;
      callResponses.set(id, { resolve, reject });
      call6(type, options).catch((error) => {
        reject(error);
        callResponses.delete(id);
      });
    });
  }
  function Call(options) {
    return callBinding(CallBinding, options);
  }
  function ByName(name, ...args) {
    if (typeof name !== "string" || name.split(".").length !== 3) {
      throw new Error("CallByName requires a string in the format 'package.struct.method'");
    }
    let [packageName, structName, methodName] = name.split(".");
    return callBinding(CallBinding, {
      packageName,
      structName,
      methodName,
      args
    });
  }
  function ByID(methodID, ...args) {
    return callBinding(CallBinding, {
      methodID,
      args
    });
  }
  function Plugin(pluginName, methodName, ...args) {
    return callBinding(CallBinding, {
      packageName: "wails-plugins",
      structName: pluginName,
      methodName,
      args
    });
  }

  // desktop/@wailsio/runtime/events.js
  var call7 = newRuntimeCallerWithID(objectNames.Events, "");
  var EmitMethod = 0;
  var eventListeners = /* @__PURE__ */ new Map();
  var Listener = class {
    constructor(eventName, callback, maxCallbacks) {
      this.eventName = eventName;
      this.maxCallbacks = maxCallbacks || -1;
      this.Callback = (data) => {
        callback(data);
        if (this.maxCallbacks === -1)
          return false;
        this.maxCallbacks -= 1;
        return this.maxCallbacks === 0;
      };
    }
  };
  var WailsEvent = class {
    constructor(name, data = null) {
      this.name = name;
      this.data = data;
    }
  };
  window._wails = window._wails || {};
  window._wails.dispatchWailsEvent = dispatchWailsEvent;
  function dispatchWailsEvent(event) {
    let listeners = eventListeners.get(event.name);
    if (listeners) {
      let toRemove = listeners.filter((listener) => {
        let remove = listener.Callback(event);
        if (remove)
          return true;
      });
      if (toRemove.length > 0) {
        listeners = listeners.filter((l) => !toRemove.includes(l));
        if (listeners.length === 0)
          eventListeners.delete(event.name);
        else
          eventListeners.set(event.name, listeners);
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
  function listenerOff(listener) {
    const eventName = listener.eventName;
    let listeners = eventListeners.get(eventName).filter((l) => l !== listener);
    if (listeners.length === 0)
      eventListeners.delete(eventName);
    else
      eventListeners.set(eventName, listeners);
  }
  function Off(eventName, ...additionalEventNames) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach((eventName2) => eventListeners.delete(eventName2));
  }
  function OffAll() {
    eventListeners.clear();
  }
  function Emit(event) {
    return call7(EmitMethod, event);
  }

  // desktop/@wailsio/runtime/dialogs.js
  var DialogInfo = 0;
  var DialogWarning = 1;
  var DialogError = 2;
  var DialogQuestion = 3;
  var DialogOpenFile = 4;
  var DialogSaveFile = 5;
  var call8 = newRuntimeCallerWithID(objectNames.Dialog, "");
  var dialogResponses = /* @__PURE__ */ new Map();
  function generateID2() {
    let result;
    do {
      result = nanoid();
    } while (dialogResponses.has(result));
    return result;
  }
  function dialog(type, options = {}) {
    const id = generateID2();
    options["dialog-id"] = id;
    return new Promise((resolve, reject) => {
      dialogResponses.set(id, { resolve, reject });
      call8(type, options).catch((error) => {
        reject(error);
        dialogResponses.delete(id);
      });
    });
  }
  function dialogCallback(id, data, isJSON) {
    let p = dialogResponses.get(id);
    if (p) {
      if (isJSON) {
        p.resolve(JSON.parse(data));
      } else {
        p.resolve(data);
      }
      dialogResponses.delete(id);
    }
  }
  function dialogErrorCallback(id, message) {
    let p = dialogResponses.get(id);
    if (p) {
      p.reject(message);
      dialogResponses.delete(id);
    }
  }
  var Info = (options) => dialog(DialogInfo, options);
  var Warning = (options) => dialog(DialogWarning, options);
  var Error2 = (options) => dialog(DialogError, options);
  var Question = (options) => dialog(DialogQuestion, options);
  var OpenFile = (options) => dialog(DialogOpenFile, options);
  var SaveFile = (options) => dialog(DialogSaveFile, options);

  // desktop/@wailsio/runtime/contextmenu.js
  var call9 = newRuntimeCallerWithID(objectNames.ContextMenu, "");
  var ContextMenuOpen = 0;
  function openContextMenu(id, x, y, data) {
    void call9(ContextMenuOpen, { id, x, y, data });
  }
  function setupContextMenus() {
    window.addEventListener("contextmenu", contextMenuHandler);
  }
  function contextMenuHandler(event) {
    let element = event.target;
    let customContextMenu = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu");
    customContextMenu = customContextMenu ? customContextMenu.trim() : "";
    if (customContextMenu) {
      event.preventDefault();
      let customContextMenuData = window.getComputedStyle(element).getPropertyValue("--custom-contextmenu-data");
      openContextMenu(customContextMenu, event.clientX, event.clientY, customContextMenuData);
      return;
    }
    processDefaultContextMenu(event);
  }
  function processDefaultContextMenu(event) {
    if (true) {
      return;
    }
    const element = event.target;
    const computedStyle = window.getComputedStyle(element);
    const defaultContextMenuAction = computedStyle.getPropertyValue("--default-contextmenu").trim();
    switch (defaultContextMenuAction) {
      case "show":
        return;
      case "hide":
        event.preventDefault();
        return;
      default:
        if (element.isContentEditable) {
          return;
        }
        const selection = window.getSelection();
        const hasSelection = selection.toString().length > 0;
        if (hasSelection) {
          for (let i = 0; i < selection.rangeCount; i++) {
            const range = selection.getRangeAt(i);
            const rects = range.getClientRects();
            for (let j = 0; j < rects.length; j++) {
              const rect = rects[j];
              if (document.elementFromPoint(rect.left, rect.top) === element) {
                return;
              }
            }
          }
        }
        if (element.tagName === "INPUT" || element.tagName === "TEXTAREA") {
          if (hasSelection || !element.readOnly && !element.disabled) {
            return;
          }
        }
        event.preventDefault();
    }
  }

  // desktop/@wailsio/runtime/wml.js
  function sendEvent(eventName, data = null) {
    let event = new WailsEvent(eventName, data);
    Emit(event);
  }
  function addWMLEventListeners() {
    const elements = document.querySelectorAll("[wml-event]");
    elements.forEach(function(element) {
      const eventType = element.getAttribute("wml-event");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Detached: false, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              sendEvent(eventType);
            }
          });
          return;
        }
        sendEvent(eventType);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function callWindowMethod(method) {
    let windowName = "";
    let targetWindow = Get("");
    let methodMap = WindowMethods(targetWindow);
    if (!methodMap.has(method)) {
      console.log("Window method " + method + " not found");
    }
    methodMap.get(method)();
  }
  function addWMLWindowListeners() {
    const elements = document.querySelectorAll("[wml-window]");
    elements.forEach(function(element) {
      const windowMethod = element.getAttribute("wml-window");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              callWindowMethod(windowMethod);
            }
          });
          return;
        }
        callWindowMethod(windowMethod);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function addWMLOpenBrowserListener() {
    const elements = document.querySelectorAll("[wml-openurl]");
    elements.forEach(function(element) {
      const url = element.getAttribute("wml-openurl");
      const confirm = element.getAttribute("wml-confirm");
      const trigger = element.getAttribute("wml-trigger") || "click";
      let callback = function() {
        if (confirm) {
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
            if (result !== "No") {
              void wails.Browser.OpenURL(url);
            }
          });
          return;
        }
        void wails.Browser.OpenURL(url);
      };
      element.removeEventListener(trigger, callback);
      element.addEventListener(trigger, callback);
    });
  }
  function reloadWML() {
    console.log("Reloading WML");
    addWMLEventListeners();
    addWMLWindowListeners();
    addWMLOpenBrowserListener();
  }

  // desktop/@wailsio/runtime/flags.js
  var flags = /* @__PURE__ */ new Map();
  function convertToMap(obj) {
    const map = /* @__PURE__ */ new Map();
    for (const [key, value] of Object.entries(obj)) {
      if (typeof value === "object" && value !== null) {
        map.set(key, convertToMap(value));
      } else {
        map.set(key, value);
      }
    }
    return map;
  }
  fetch("/wails/flags").then((response) => {
    response.json().then((data) => {
      flags = convertToMap(data);
    });
  });
  function getValueFromMap(keyString) {
    const keys = keyString.split(".");
    let value = flags;
    for (const key of keys) {
      if (value instanceof Map) {
        value = value.get(key);
      } else {
        value = value[key];
      }
      if (value === void 0) {
        break;
      }
    }
    return value;
  }
  function GetFlag(keyString) {
    return getValueFromMap(keyString);
  }

  // desktop/@wailsio/runtime/drag.js
  var shouldDrag = false;
  var resizeEdge = null;
  var resizable = false;
  var defaultCursor = "auto";
  window._wails = window._wails || {};
  window._wails.setResizable = setResizable2;
  window._wails.endDrag = endDrag;
  function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--webkit-app-region");
    if (val && val.trim() !== "drag" || e.buttons !== 1) {
      return false;
    }
    return e.detail === 1;
  }
  function setupDrag() {
    window.addEventListener("mousedown", onMouseDown);
    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
  }
  function setResizable2(value) {
    resizable = value;
  }
  function endDrag() {
    document.body.style.cursor = "default";
    shouldDrag = false;
  }
  function testResize() {
    if (resizeEdge) {
      invoke(`resize:${resizeEdge}`);
      return true;
    }
    return false;
  }
  function onMouseDown(e) {
    if (IsWindows() && testResize() || dragTest(e)) {
      shouldDrag = !!isValidDrag(e);
    }
  }
  function isValidDrag(e) {
    return !(e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight);
  }
  function onMouseUp(e) {
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (mousePressed > 0) {
      endDrag();
    }
  }
  function setResize(cursor = defaultCursor) {
    document.documentElement.style.cursor = cursor;
    resizeEdge = cursor;
  }
  function onMouseMove(e) {
    shouldDrag = checkDrag(e);
    if (IsWindows() && resizable) {
      handleResize(e);
    }
  }
  function checkDrag(e) {
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (shouldDrag && mousePressed > 0) {
      invoke("drag");
      return false;
    }
    return shouldDrag;
  }
  function handleResize(e) {
    let resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
    let resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
    let cornerExtra = GetFlag("resizeCornerExtra") || 10;
    let rightBorder = window.outerWidth - e.clientX < resizeHandleWidth;
    let leftBorder = e.clientX < resizeHandleWidth;
    let topBorder = e.clientY < resizeHandleHeight;
    let bottomBorder = window.outerHeight - e.clientY < resizeHandleHeight;
    let rightCorner = window.outerWidth - e.clientX < resizeHandleWidth + cornerExtra;
    let leftCorner = e.clientX < resizeHandleWidth + cornerExtra;
    let topCorner = e.clientY < resizeHandleHeight + cornerExtra;
    let bottomCorner = window.outerHeight - e.clientY < resizeHandleHeight + cornerExtra;
    if (!leftBorder && !rightBorder && !topBorder && !bottomBorder && resizeEdge !== void 0) {
      setResize();
    } else if (rightCorner && bottomCorner)
      setResize("se-resize");
    else if (leftCorner && bottomCorner)
      setResize("sw-resize");
    else if (leftCorner && topCorner)
      setResize("nw-resize");
    else if (topCorner && rightCorner)
      setResize("ne-resize");
    else if (leftBorder)
      setResize("w-resize");
    else if (topBorder)
      setResize("n-resize");
    else if (bottomBorder)
      setResize("s-resize");
    else if (rightBorder)
      setResize("e-resize");
  }

  // desktop/main.js
  window.wails = {
    ...newRuntime(null),
    clientId
  };
  window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchWailsEvent,
    callErrorHandler: errorHandler,
    callResultHandler: resultHandler,
    endDrag,
    setResizable: setResizable2
  };
  function newRuntime(windowName) {
    return {
      Clipboard: {
        ...clipboard_exports
      },
      Application: {
        ...application_exports
      },
      System: system_exports,
      Screens: screens_exports,
      Browser: browser_exports,
      Call: {
        Call,
        ByID,
        ByName,
        Plugin
      },
      WML: {
        Reload: reloadWML
      },
      Dialog: {
        Info,
        Warning,
        Error: Error2,
        Question,
        OpenFile,
        SaveFile
      },
      Events: {
        Emit,
        On,
        Once,
        OnMultiple,
        Off,
        OffAll
      },
      Window: {
        ...Get("")
      }
    };
  }
  setupContextMenus();
  setupDrag();
  document.addEventListener("DOMContentLoaded", function() {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL2NsaXBib2FyZC5qcyIsICJub2RlX21vZHVsZXMvbmFub2lkL25vbi1zZWN1cmUvaW5kZXguanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3J1bnRpbWUuanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL2FwcGxpY2F0aW9uLmpzIiwgImRlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zY3JlZW5zLmpzIiwgImRlc2t0b3AvQHdhaWxzaW8vcnVudGltZS9zeXN0ZW0uanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL2Jyb3dzZXIuanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3dpbmRvdy5qcyIsICJkZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvY2FsbHMuanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL2V2ZW50cy5qcyIsICJkZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvZGlhbG9ncy5qcyIsICJkZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvY29udGV4dG1lbnUuanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL3dtbC5qcyIsICJkZXNrdG9wL0B3YWlsc2lvL3J1bnRpbWUvZmxhZ3MuanMiLCAiZGVza3RvcC9Ad2FpbHNpby9ydW50aW1lL2RyYWcuanMiLCAiZGVza3RvcC9tYWluLmpzIl0sCiAgInNvdXJjZXNDb250ZW50IjogWyIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5DbGlwYm9hcmQsICcnKTtcclxuY29uc3QgQ2xpcGJvYXJkU2V0VGV4dCA9IDA7XHJcbmNvbnN0IENsaXBib2FyZFRleHQgPSAxO1xyXG5cclxuLyoqXHJcbiAqIFNldHMgdGhlIHRleHQgdG8gdGhlIENsaXBib2FyZC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IHRleHQgLSBUaGUgdGV4dCB0byBiZSBzZXQgdG8gdGhlIENsaXBib2FyZC5cclxuICogQHJldHVybiB7UHJvbWlzZX0gLSBBIFByb21pc2UgdGhhdCByZXNvbHZlcyB3aGVuIHRoZSBvcGVyYXRpb24gaXMgc3VjY2Vzc2Z1bC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcclxuICAgIHJldHVybiBjYWxsKENsaXBib2FyZFNldFRleHQsIHt0ZXh0fSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHRleHQgZnJvbSB0aGUgQ2xpcGJvYXJkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcclxufVxyXG4iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXHJcbiBfICAgICBfXyAgICAgXyBfX1xyXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xyXG5cclxuY29uc3QgcnVudGltZVVSTCA9IHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XHJcblxyXG4vLyBPYmplY3QgTmFtZXNcclxuZXhwb3J0IGNvbnN0IG9iamVjdE5hbWVzID0ge1xyXG4gICAgQ2FsbDogMCxcclxuICAgIENsaXBib2FyZDogMSxcclxuICAgIEFwcGxpY2F0aW9uOiAyLFxyXG4gICAgRXZlbnRzOiAzLFxyXG4gICAgQ29udGV4dE1lbnU6IDQsXHJcbiAgICBEaWFsb2c6IDUsXHJcbiAgICBXaW5kb3c6IDYsXHJcbiAgICBTY3JlZW5zOiA3LFxyXG4gICAgU3lzdGVtOiA4LFxyXG4gICAgQnJvd3NlcjogOSxcclxufVxyXG5leHBvcnQgbGV0IGNsaWVudElkID0gbmFub2lkKCk7XHJcblxyXG4vKipcclxuICogQ3JlYXRlcyBhIHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uIHRoYXQgaW52b2tlcyBhIHNwZWNpZmllZCBtZXRob2Qgb24gYSBnaXZlbiBvYmplY3Qgd2l0aGluIGEgc3BlY2lmaWVkIHdpbmRvdyBjb250ZXh0LlxyXG4gKlxyXG4gKiBAcGFyYW0ge09iamVjdH0gb2JqZWN0IC0gVGhlIG9iamVjdCBvbiB3aGljaCB0aGUgbWV0aG9kIGlzIHRvIGJlIGludm9rZWQuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB3aW5kb3dOYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHdpbmRvdyBjb250ZXh0IGluIHdoaWNoIHRoZSBtZXRob2Qgc2hvdWxkIGJlIGNhbGxlZC5cclxuICogQHJldHVybnMge0Z1bmN0aW9ufSBBIHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uIHRoYXQgdGFrZXMgdGhlIG1ldGhvZCBuYW1lIGFuZCBvcHRpb25hbGx5IGFyZ3VtZW50cyBhbmQgaW52b2tlcyB0aGUgbWV0aG9kIHdpdGhpbiB0aGUgc3BlY2lmaWVkIHdpbmRvdyBjb250ZXh0LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0LCB3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDcmVhdGVzIGEgbmV3IHJ1bnRpbWUgY2FsbGVyIHdpdGggc3BlY2lmaWVkIElELlxyXG4gKlxyXG4gKiBAcGFyYW0ge29iamVjdH0gb2JqZWN0IC0gVGhlIG9iamVjdCB0byBpbnZva2UgdGhlIG1ldGhvZCBvbi5cclxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93LlxyXG4gKiBAcmV0dXJuIHtGdW5jdGlvbn0gLSBUaGUgbmV3IHJ1bnRpbWUgY2FsbGVyIGZ1bmN0aW9uLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0LCB3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn1cclxuXHJcblxyXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgaWYoIG1ldGhvZCApIHtcclxuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xyXG4gICAgfVxyXG4gICAgbGV0IGZldGNoT3B0aW9ucyA9IHtcclxuICAgICAgICBoZWFkZXJzOiB7fSxcclxuICAgIH07XHJcbiAgICBpZiAod2luZG93TmFtZSkge1xyXG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XHJcbiAgICB9XHJcbiAgICBpZiAoYXJncykge1xyXG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XHJcbiAgICB9XHJcbiAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtY2xpZW50LWlkXCJdID0gY2xpZW50SWQ7XHJcblxyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncykge1xyXG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XHJcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm9iamVjdFwiLCBvYmplY3RJRCk7XHJcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xyXG4gICAgbGV0IGZldGNoT3B0aW9ucyA9IHtcclxuICAgICAgICBoZWFkZXJzOiB7fSxcclxuICAgIH07XHJcbiAgICBpZiAod2luZG93TmFtZSkge1xyXG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XHJcbiAgICB9XHJcbiAgICBpZiAoYXJncykge1xyXG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XHJcbiAgICB9XHJcbiAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtY2xpZW50LWlkXCJdID0gY2xpZW50SWQ7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxyXG4gICAgICAgICAgICAudGhlbihyZXNwb25zZSA9PiB7XHJcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcclxuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcclxuICAgICAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2UuaGVhZGVycy5nZXQoXCJDb250ZW50LVR5cGVcIikgJiYgcmVzcG9uc2UuaGVhZGVycy5nZXQoXCJDb250ZW50LVR5cGVcIikuaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgIT09IC0xKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLnRleHQoKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICByZWplY3QoRXJyb3IocmVzcG9uc2Uuc3RhdHVzVGV4dCkpO1xyXG4gICAgICAgICAgICB9KVxyXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXHJcbiAgICAgICAgICAgIC5jYXRjaChlcnJvciA9PiByZWplY3QoZXJyb3IpKTtcclxuICAgIH0pO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQgeyBuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lcyB9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24pO1xyXG5cclxuY29uc3QgSGlkZU1ldGhvZCA9IDA7XHJcbmNvbnN0IFNob3dNZXRob2QgPSAxO1xyXG5jb25zdCBRdWl0TWV0aG9kID0gMjtcclxuXHJcbi8qKlxyXG4gKiBIaWRlcyBhIGNlcnRhaW4gbWV0aG9kIGJ5IGNhbGxpbmcgdGhlIEhpZGVNZXRob2QgZnVuY3Rpb24uXHJcbiAqXHJcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAqXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSGlkZSgpIHtcclxuICAgIHJldHVybiBjYWxsKEhpZGVNZXRob2QpO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgdGhlIFNob3dNZXRob2QgYW5kIHJldHVybnMgdGhlIHJlc3VsdC5cclxuICpcclxuICogQHJldHVybiB7UHJvbWlzZTx2b2lkPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xyXG4gICAgcmV0dXJuIGNhbGwoU2hvd01ldGhvZCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDYWxscyB0aGUgUXVpdE1ldGhvZCB0byB0ZXJtaW5hdGUgdGhlIHByb2dyYW0uXHJcbiAqXHJcbiAqIEByZXR1cm4ge1Byb21pc2U8dm9pZD59XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcclxuICAgIHJldHVybiBjYWxsKFF1aXRNZXRob2QpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TY3JlZW5zLCAnJyk7XHJcblxyXG5jb25zdCBnZXRBbGwgPSAwO1xyXG5jb25zdCBnZXRQcmltYXJ5ID0gMTtcclxuY29uc3QgZ2V0Q3VycmVudCA9IDI7XHJcblxyXG4vKipcclxuICogR2V0cyBhbGwgc2NyZWVucy5cclxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhbiBhcnJheSBvZiBTY3JlZW4gb2JqZWN0cy5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKSB7XHJcbiAgICByZXR1cm4gY2FsbChnZXRBbGwpO1xyXG59XHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cclxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XHJcbiAgICByZXR1cm4gY2FsbChnZXRQcmltYXJ5KTtcclxufVxyXG4vKipcclxuICogR2V0cyB0aGUgY3VycmVudCBhY3RpdmUgc2NyZWVuLlxyXG4gKlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gR2V0Q3VycmVudCgpIHtcclxuICAgIHJldHVybiBjYWxsKGdldEN1cnJlbnQpO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLlN5c3RlbSwgJycpO1xyXG5jb25zdCBzeXN0ZW1Jc0RhcmtNb2RlID0gMDtcclxuY29uc3QgZW52aXJvbm1lbnQgPSAxO1xyXG5cclxuLyoqXHJcbiAqIEBmdW5jdGlvblxyXG4gKiBSZXRyaWV2ZXMgdGhlIHN5c3RlbSBkYXJrIG1vZGUgc3RhdHVzLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxib29sZWFuPn0gLSBBIHByb21pc2UgdGhhdCByZXNvbHZlcyB0byBhIGJvb2xlYW4gdmFsdWUgaW5kaWNhdGluZyBpZiB0aGUgc3lzdGVtIGlzIGluIGRhcmsgbW9kZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0RhcmtNb2RlKCkge1xyXG4gICAgcmV0dXJuIGNhbGwoc3lzdGVtSXNEYXJrTW9kZSk7XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogRmV0Y2hlcyB0aGUgY2FwYWJpbGl0aWVzIG9mIHRoZSBhcHBsaWNhdGlvbiBmcm9tIHRoZSBzZXJ2ZXIuXHJcbiAqXHJcbiAqIEBhc3luY1xyXG4gKiBAZnVuY3Rpb24gQ2FwYWJpbGl0aWVzXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPE9iamVjdD59IEEgcHJvbWlzZSB0aGF0IHJlc29sdmVzIHRvIGFuIG9iamVjdCBjb250YWluaW5nIHRoZSBjYXBhYmlsaXRpZXMuXHJcbiAqL1xyXG5leHBvcnQgYXN5bmMgZnVuY3Rpb24gQ2FwYWJpbGl0aWVzKCkge1xyXG4gICAgbGV0IHJlc3BvbnNlID0gZmV0Y2goXCIvd2FpbHMvY2FwYWJpbGl0aWVzXCIpO1xyXG4gICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtvYmplY3R9IEVudmlyb25tZW50SW5mb1xyXG4gKiBAcHJvcGVydHkge3N0cmluZ30gT1MgLSBUaGUgb3BlcmF0aW5nIHN5c3RlbSBpbiB1c2UuXHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBBcmNoIC0gVGhlIGFyY2hpdGVjdHVyZSBvZiB0aGUgc3lzdGVtLlxyXG4gKi9cclxuXHJcbi8qKlxyXG4gKiBAZnVuY3Rpb25cclxuICogUmV0cmlldmVzIGVudmlyb25tZW50IGRldGFpbHMuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPEVudmlyb25tZW50SW5mbz59IC0gQSBwcm9taXNlIHRoYXQgcmVzb2x2ZXMgdG8gYW4gb2JqZWN0IGNvbnRhaW5pbmcgT1MgYW5kIHN5c3RlbSBhcmNoaXRlY3R1cmUuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW52aXJvbm1lbnQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChlbnZpcm9ubWVudCk7XHJcbn1cclxuXHJcbmV4cG9ydCBsZXQgaW52b2tlID0gbnVsbDtcclxubGV0IGVudmlyb25tZW50Q2FjaGUgPSBudWxsO1xyXG5cclxuRW52aXJvbm1lbnQoKVxyXG4gICAgLnRoZW4ocmVzdWx0ID0+IHtcclxuICAgICAgICBlbnZpcm9ubWVudENhY2hlID0gcmVzdWx0O1xyXG4gICAgICAgIGludm9rZSA9IElzV2luZG93cygpID8gd2luZG93LmNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlIDogd2luZG93LndlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2U7XHJcbiAgICB9KVxyXG4gICAgLmNhdGNoKGVycm9yID0+IHtcclxuICAgICAgICBjb25zb2xlLmVycm9yKGBFcnJvciBnZXR0aW5nIEVudmlyb25tZW50OiAke2Vycm9yfWApO1xyXG4gICAgfSk7XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgV2luZG93cy5cclxuICpcclxuICogQHJldHVybiB7Ym9vbGVhbn0gVHJ1ZSBpZiB0aGUgb3BlcmF0aW5nIHN5c3RlbSBpcyBXaW5kb3dzLCBvdGhlcndpc2UgZmFsc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNXaW5kb3dzKCkge1xyXG4gICAgcmV0dXJuIGVudmlyb25tZW50Q2FjaGUuT1MgPT09IFwid2luZG93c1wiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IG9wZXJhdGluZyBzeXN0ZW0gaXMgTGludXguXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSBSZXR1cm5zIHRydWUgaWYgdGhlIGN1cnJlbnQgb3BlcmF0aW5nIHN5c3RlbSBpcyBMaW51eCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzTGludXgoKSB7XHJcbiAgICByZXR1cm4gZW52aXJvbm1lbnRDYWNoZS5PUyA9PT0gXCJsaW51eFwiO1xyXG59XHJcblxyXG4vKipcclxuICogQ2hlY2tzIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGlzIGEgbWFjT1Mgb3BlcmF0aW5nIHN5c3RlbS5cclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGVudmlyb25tZW50IGlzIG1hY09TLCBmYWxzZSBvdGhlcndpc2UuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSXNNYWMoKSB7XHJcbiAgICByZXR1cm4gZW52aXJvbm1lbnRDYWNoZS5PUyA9PT0gXCJkYXJ3aW5cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBhcmNoaXRlY3R1cmUgaXMgQU1ENjQuXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSBUcnVlIGlmIHRoZSBjdXJyZW50IGVudmlyb25tZW50IGFyY2hpdGVjdHVyZSBpcyBBTUQ2NCwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQU1ENjQoKSB7XHJcbiAgICByZXR1cm4gZW52aXJvbm1lbnRDYWNoZS5BcmNoID09PSBcImFtZDY0XCI7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBDaGVja3MgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTS5cclxuICpcclxuICogQHJldHVybnMge2Jvb2xlYW59IFRydWUgaWYgdGhlIGN1cnJlbnQgYXJjaGl0ZWN0dXJlIGlzIEFSTSwgZmFsc2Ugb3RoZXJ3aXNlLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIElzQVJNKCkge1xyXG4gICAgcmV0dXJuIGVudmlyb25tZW50Q2FjaGUuQXJjaCA9PT0gXCJhcm1cIjtcclxufVxyXG5cclxuLyoqXHJcbiAqIENoZWNrcyBpZiB0aGUgY3VycmVudCBlbnZpcm9ubWVudCBpcyBBUk02NCBhcmNoaXRlY3R1cmUuXHJcbiAqXHJcbiAqIEByZXR1cm5zIHtib29sZWFufSAtIFJldHVybnMgdHJ1ZSBpZiB0aGUgZW52aXJvbm1lbnQgaXMgQVJNNjQgYXJjaGl0ZWN0dXJlLCBvdGhlcndpc2UgcmV0dXJucyBmYWxzZS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBJc0FSTTY0KCkge1xyXG4gICAgcmV0dXJuIGVudmlyb25tZW50Q2FjaGUuQXJjaCA9PT0gXCJhcm02NFwiO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkJyb3dzZXIsICcnKTtcclxuY29uc3QgQnJvd3Nlck9wZW5VUkwgPSAwO1xyXG5cclxuLyoqXHJcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybCkge1xyXG4gICAgcmV0dXJuIGNhbGwoQnJvd3Nlck9wZW5VUkwsIHt1cmx9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2l6ZX0gU2l6ZVxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlBvc2l0aW9ufSBQb3NpdGlvblxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuY29uc3QgY2VudGVyID0gMDtcclxuY29uc3Qgc2V0VGl0bGUgPSAxO1xyXG5jb25zdCBmdWxsc2NyZWVuID0gMjtcclxuY29uc3QgdW5GdWxsc2NyZWVuID0gMztcclxuY29uc3Qgc2V0U2l6ZSA9IDQ7XHJcbmNvbnN0IHNpemUgPSA1O1xyXG5jb25zdCBzZXRNYXhTaXplID0gNjtcclxuY29uc3Qgc2V0TWluU2l6ZSA9IDc7XHJcbmNvbnN0IHNldEFsd2F5c09uVG9wID0gODtcclxuY29uc3Qgc2V0UmVsYXRpdmVQb3NpdGlvbiA9IDk7XHJcbmNvbnN0IHJlbGF0aXZlUG9zaXRpb24gPSAxMDtcclxuY29uc3Qgc2NyZWVuID0gMTE7XHJcbmNvbnN0IGhpZGUgPSAxMjtcclxuY29uc3QgbWF4aW1pc2UgPSAxMztcclxuY29uc3QgdW5NYXhpbWlzZSA9IDE0O1xyXG5jb25zdCB0b2dnbGVNYXhpbWlzZSA9IDE1O1xyXG5jb25zdCBtaW5pbWlzZSA9IDE2O1xyXG5jb25zdCB1bk1pbmltaXNlID0gMTc7XHJcbmNvbnN0IHJlc3RvcmUgPSAxODtcclxuY29uc3Qgc2hvdyA9IDE5O1xyXG5jb25zdCBjbG9zZSA9IDIwO1xyXG5jb25zdCBzZXRCYWNrZ3JvdW5kQ29sb3VyID0gMjE7XHJcbmNvbnN0IHNldFJlc2l6YWJsZSA9IDIyO1xyXG5jb25zdCB3aWR0aCA9IDIzO1xyXG5jb25zdCBoZWlnaHQgPSAyNDtcclxuY29uc3Qgem9vbUluID0gMjU7XHJcbmNvbnN0IHpvb21PdXQgPSAyNjtcclxuY29uc3Qgem9vbVJlc2V0ID0gMjc7XHJcbmNvbnN0IGdldFpvb21MZXZlbCA9IDI4O1xyXG5jb25zdCBzZXRab29tTGV2ZWwgPSAyOTtcclxuXHJcbmNvbnN0IHRoaXNXaW5kb3cgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgJycpO1xyXG5cclxuZnVuY3Rpb24gY3JlYXRlV2luZG93KGNhbGwpIHtcclxuICAgIHJldHVybiB7XHJcbiAgICAgICAgR2V0OiAod2luZG93TmFtZSkgPT4gY3JlYXRlV2luZG93KG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuV2luZG93LCB3aW5kb3dOYW1lKSksXHJcbiAgICAgICAgQ2VudGVyOiAoKSA9PiBjYWxsKGNlbnRlciksXHJcbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbChzZXRUaXRsZSwge3RpdGxlfSksXHJcbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gY2FsbChmdWxsc2NyZWVuKSxcclxuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwodW5GdWxsc2NyZWVuKSxcclxuICAgICAgICBTZXRTaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbChzZXRTaXplLCB7d2lkdGgsIGhlaWdodH0pLFxyXG4gICAgICAgIFNpemU6ICgpID0+IGNhbGwoc2l6ZSksXHJcbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IGNhbGwoc2V0TWF4U2l6ZSwge3dpZHRoLCBoZWlnaHR9KSxcclxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gY2FsbChzZXRNaW5TaXplLCB7d2lkdGgsIGhlaWdodH0pLFxyXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAob25Ub3ApID0+IGNhbGwoc2V0QWx3YXlzT25Ub3AsIHthbHdheXNPblRvcDogb25Ub3B9KSxcclxuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbChzZXRSZWxhdGl2ZVBvc2l0aW9uLCB7eCwgeX0pLFxyXG4gICAgICAgIFJlbGF0aXZlUG9zaXRpb246ICgpID0+IGNhbGwocmVsYXRpdmVQb3NpdGlvbiksXHJcbiAgICAgICAgU2NyZWVuOiAoKSA9PiBjYWxsKHNjcmVlbiksXHJcbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbChoaWRlKSxcclxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gY2FsbChtYXhpbWlzZSksXHJcbiAgICAgICAgVW5NYXhpbWlzZTogKCkgPT4gY2FsbCh1bk1heGltaXNlKSxcclxuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gY2FsbCh0b2dnbGVNYXhpbWlzZSksXHJcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IGNhbGwobWluaW1pc2UpLFxyXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwodW5NaW5pbWlzZSksXHJcbiAgICAgICAgUmVzdG9yZTogKCkgPT4gY2FsbChyZXN0b3JlKSxcclxuICAgICAgICBTaG93OiAoKSA9PiBjYWxsKHNob3cpLFxyXG4gICAgICAgIENsb3NlOiAoKSA9PiBjYWxsKGNsb3NlKSxcclxuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gY2FsbChzZXRCYWNrZ3JvdW5kQ29sb3VyLCB7ciwgZywgYiwgYX0pLFxyXG4gICAgICAgIFNldFJlc2l6YWJsZTogKHJlc2l6YWJsZSkgPT4gY2FsbChzZXRSZXNpemFibGUsIHtyZXNpemFibGV9KSxcclxuICAgICAgICBXaWR0aDogKCkgPT4gY2FsbCh3aWR0aCksXHJcbiAgICAgICAgSGVpZ2h0OiAoKSA9PiBjYWxsKGhlaWdodCksXHJcbiAgICAgICAgWm9vbUluOiAoKSA9PiBjYWxsKHpvb21JbiksXHJcbiAgICAgICAgWm9vbU91dDogKCkgPT4gY2FsbCh6b29tT3V0KSxcclxuICAgICAgICBab29tUmVzZXQ6ICgpID0+IGNhbGwoem9vbVJlc2V0KSxcclxuICAgICAgICBHZXRab29tTGV2ZWw6ICgpID0+IGNhbGwoZ2V0Wm9vbUxldmVsKSxcclxuICAgICAgICBTZXRab29tTGV2ZWw6ICh6b29tTGV2ZWwpID0+IGNhbGwoc2V0Wm9vbUxldmVsLCB7em9vbUxldmVsfSksXHJcbiAgICB9O1xyXG59XHJcblxyXG4vKipcclxuICogR2V0cyB0aGUgc3BlY2lmaWVkIHdpbmRvdy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IHdpbmRvd05hbWUgLSBUaGUgbmFtZSBvZiB0aGUgd2luZG93IHRvIGdldC5cclxuICogQHJldHVybiB7T2JqZWN0fSAtIFRoZSBzcGVjaWZpZWQgd2luZG93IG9iamVjdC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXQod2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIGNyZWF0ZVdpbmRvdyhuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgd2luZG93TmFtZSkpO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0dXJucyBhIG1hcCBvZiBhbGwgbWV0aG9kcyBpbiB0aGUgY3VycmVudCB3aW5kb3cuXHJcbiAqIEByZXR1cm5zIHtNYXB9IC0gQSBtYXAgb2Ygd2luZG93IG1ldGhvZHMuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2luZG93TWV0aG9kcyh0YXJnZXRXaW5kb3cpIHtcclxuICAgIC8vIENyZWF0ZSBhIG5ldyBtYXAgdG8gc3RvcmUgbWV0aG9kc1xyXG4gICAgbGV0IHJlc3VsdCA9IG5ldyBNYXAoKTtcclxuXHJcbiAgICAvLyBJdGVyYXRlIG92ZXIgYWxsIHByb3BlcnRpZXMgb2YgdGhlIHdpbmRvdyBvYmplY3RcclxuICAgIGZvciAobGV0IG1ldGhvZCBpbiB0YXJnZXRXaW5kb3cpIHtcclxuICAgICAgICAvLyBDaGVjayBpZiB0aGUgcHJvcGVydHkgaXMgaW5kZWVkIGEgbWV0aG9kIChmdW5jdGlvbilcclxuICAgICAgICBpZih0eXBlb2YgdGFyZ2V0V2luZG93W21ldGhvZF0gPT09ICdmdW5jdGlvbicpIHtcclxuICAgICAgICAgICAgLy8gQWRkIHRoZSBtZXRob2QgdG8gdGhlIG1hcFxyXG4gICAgICAgICAgICByZXN1bHQuc2V0KG1ldGhvZCwgdGFyZ2V0V2luZG93W21ldGhvZF0pO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICB9XHJcbiAgICAvLyBSZXR1cm4gdGhlIG1hcCBvZiB3aW5kb3cgbWV0aG9kc1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5leHBvcnQgZGVmYXVsdCB7XHJcbiAgICAuLi5HZXQoJycpXHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuaW1wb3J0IHsgbmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXMgfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbmNvbnN0IENhbGxCaW5kaW5nID0gMDtcclxuY29uc3QgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2FsbCwgJycpO1xyXG5sZXQgY2FsbFJlc3BvbnNlcyA9IG5ldyBNYXAoKTtcclxuXHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLmNhbGxDYWxsYmFjayA9IHJlc3VsdEhhbmRsZXI7XHJcbndpbmRvdy5fd2FpbHMuY2FsbEVycm9yQ2FsbGJhY2sgPSBlcnJvckhhbmRsZXI7XHJcblxyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHJlc3VsdEhhbmRsZXIoaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XHJcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcclxuICAgICAgICBwcm9taXNlSGFuZGxlci5yZXNvbHZlKGlzSlNPTiA/IEpTT04ucGFyc2UoZGF0YSkgOiBkYXRhKTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGVycm9ySGFuZGxlcihpZCwgbWVzc2FnZSkge1xyXG4gICAgY29uc3QgcHJvbWlzZUhhbmRsZXIgPSBnZXRBbmREZWxldGVSZXNwb25zZShpZCk7XHJcbiAgICBpZiAocHJvbWlzZUhhbmRsZXIpIHtcclxuICAgICAgICBwcm9taXNlSGFuZGxlci5yZWplY3QobWVzc2FnZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGdldEFuZERlbGV0ZVJlc3BvbnNlKGlkKSB7XHJcbiAgICBjb25zdCByZXNwb25zZSA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIHJldHVybiByZXNwb25zZTtcclxufVxyXG5cclxuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucyA9IHt9KSB7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGNvbnN0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHsgcmVzb2x2ZSwgcmVqZWN0IH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbCBtZXRob2QuXHJcbiAqXHJcbiAqIEBwYXJhbSB7T2JqZWN0fSBvcHRpb25zIC0gVGhlIG9wdGlvbnMgZm9yIHRoZSBtZXRob2QuXHJcbiAqIEByZXR1cm5zIHtPYmplY3R9IC0gVGhlIHJlc3VsdCBvZiB0aGUgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBDYWxsKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFeGVjdXRlcyBhIG1ldGhvZCBieSBuYW1lLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBtZXRob2QgaW4gdGhlIGZvcm1hdCAncGFja2FnZS5zdHJ1Y3QubWV0aG9kJy5cclxuICogQHBhcmFtIHsuLi4qfSBhcmdzIC0gVGhlIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2QuXHJcbiAqIEB0aHJvd3Mge0Vycm9yfSBJZiB0aGUgbmFtZSBpcyBub3QgYSBzdHJpbmcgb3IgaXMgbm90IGluIHRoZSBjb3JyZWN0IGZvcm1hdC5cclxuICogQHJldHVybnMgeyp9IFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBleGVjdXRpb24uXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gQnlOYW1lKG5hbWUsIC4uLmFyZ3MpIHtcclxuICAgIGlmICh0eXBlb2YgbmFtZSAhPT0gXCJzdHJpbmdcIiB8fCBuYW1lLnNwbGl0KFwiLlwiKS5sZW5ndGggIT09IDMpIHtcclxuICAgICAgICB0aHJvdyBuZXcgRXJyb3IoXCJDYWxsQnlOYW1lIHJlcXVpcmVzIGEgc3RyaW5nIGluIHRoZSBmb3JtYXQgJ3BhY2thZ2Uuc3RydWN0Lm1ldGhvZCdcIik7XHJcbiAgICB9XHJcbiAgICBsZXQgW3BhY2thZ2VOYW1lLCBzdHJ1Y3ROYW1lLCBtZXRob2ROYW1lXSA9IG5hbWUuc3BsaXQoXCIuXCIpO1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWUsXHJcbiAgICAgICAgc3RydWN0TmFtZSxcclxuICAgICAgICBtZXRob2ROYW1lLFxyXG4gICAgICAgIGFyZ3NcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2QgYnkgaXRzIElEIHdpdGggdGhlIHNwZWNpZmllZCBhcmd1bWVudHMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2RJRCAtIFRoZSBJRCBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJuIHsqfSAtIFRoZSByZXN1bHQgb2YgdGhlIG1ldGhvZCBjYWxsLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEJ5SUQobWV0aG9kSUQsIC4uLmFyZ3MpIHtcclxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xyXG4gICAgICAgIG1ldGhvZElELFxyXG4gICAgICAgIGFyZ3NcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQ2FsbHMgYSBtZXRob2Qgb24gYSBwbHVnaW4uXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gVGhlIG5hbWUgb2YgdGhlIHBsdWdpbi5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1ldGhvZE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgbWV0aG9kIHRvIGNhbGwuXHJcbiAqIEBwYXJhbSB7Li4uKn0gYXJncyAtIFRoZSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kLlxyXG4gKiBAcmV0dXJucyB7Kn0gLSBUaGUgcmVzdWx0IG9mIHRoZSBtZXRob2QgY2FsbC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxyXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXHJcbiAgICAgICAgbWV0aG9kTmFtZSxcclxuICAgICAgICBhcmdzXHJcbiAgICB9KTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcclxuICovXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkV2ZW50cywgJycpO1xyXG5jb25zdCBFbWl0TWV0aG9kID0gMDtcclxuY29uc3QgZXZlbnRMaXN0ZW5lcnMgPSBuZXcgTWFwKCk7XHJcblxyXG5jbGFzcyBMaXN0ZW5lciB7XHJcbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgICAgICB0aGlzLmV2ZW50TmFtZSA9IGV2ZW50TmFtZTtcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcclxuICAgICAgICAgICAgY2FsbGJhY2soZGF0YSk7XHJcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgLT0gMTtcclxuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xyXG4gICAgICAgIH07XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcclxuICAgIGNvbnN0cnVjdG9yKG5hbWUsIGRhdGEgPSBudWxsKSB7XHJcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcclxuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xyXG4gICAgfVxyXG59XHJcblxyXG5cclxud2luZG93Ll93YWlscyA9IHdpbmRvdy5fd2FpbHMgfHwge307XHJcbndpbmRvdy5fd2FpbHMuZGlzcGF0Y2hXYWlsc0V2ZW50ID0gZGlzcGF0Y2hXYWlsc0V2ZW50O1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudCkge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcclxuICAgIGlmIChsaXN0ZW5lcnMpIHtcclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBsaXN0ZW5lcnMuZmlsdGVyKGxpc3RlbmVyID0+IHtcclxuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcclxuICAgICAgICAgICAgaWYgKHJlbW92ZSkgcmV0dXJuIHRydWU7XHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XHJcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XHJcbiAgICAgICAgICAgIGVsc2UgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXIgYSBjYWxsYmFjayBmdW5jdGlvbiB0byBiZSBjYWxsZWQgbXVsdGlwbGUgdGltZXMgZm9yIGEgc3BlY2lmaWMgZXZlbnQuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVnaXN0ZXIgdGhlIGNhbGxiYWNrIGZvci5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgY2FsbGVkIHdoZW4gdGhlIGV2ZW50IGlzIHRyaWdnZXJlZC5cclxuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrcyAtIFRoZSBtYXhpbXVtIG51bWJlciBvZiB0aW1lcyB0aGUgY2FsbGJhY2sgY2FuIGJlIGNhbGxlZCBmb3IgdGhlIGV2ZW50LiBPbmNlIHRoZSBtYXhpbXVtIG51bWJlciBpcyByZWFjaGVkLCB0aGUgY2FsbGJhY2sgd2lsbCBubyBsb25nZXIgYmUgY2FsbGVkLlxyXG4gKlxyXG4gQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpIHx8IFtdO1xyXG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XHJcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xyXG4gICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcclxuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xyXG59XHJcblxyXG4vKipcclxuICogUmVnaXN0ZXJzIGEgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgc3BlY2lmaWVkIGV2ZW50IG9jY3Vycy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgY2FsbGJhY2sgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQuIEl0IHRha2VzIG5vIHBhcmFtZXRlcnMuXHJcbiAqIEByZXR1cm4ge2Z1bmN0aW9ufSAtIEEgZnVuY3Rpb24gdGhhdCwgd2hlbiBjYWxsZWQsIHdpbGwgdW5yZWdpc3RlciB0aGUgY2FsbGJhY2sgZnJvbSB0aGUgZXZlbnQuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7IHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTsgfVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhIGNhbGxiYWNrIGZ1bmN0aW9uIHRvIGJlIGV4ZWN1dGVkIG9ubHkgb25jZSBmb3IgdGhlIHNwZWNpZmllZCBldmVudC5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZSAtIFRoZSBuYW1lIG9mIHRoZSBldmVudC5cclxuICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2sgLSBUaGUgZnVuY3Rpb24gdG8gYmUgZXhlY3V0ZWQgd2hlbiB0aGUgZXZlbnQgb2NjdXJzLlxyXG4gKiBAcmV0dXJuIHt2b2lkQHJldHVybiB7ZnVuY3Rpb259IC0gQSBmdW5jdGlvbiB0aGF0LCB3aGVuIGNhbGxlZCwgd2lsbCB1bnJlZ2lzdGVyIHRoZSBjYWxsYmFjayBmcm9tIHRoZSBldmVudC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHsgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7IH1cclxuXHJcbi8qKlxyXG4gKiBSZW1vdmVzIHRoZSBzcGVjaWZpZWQgbGlzdGVuZXIgZnJvbSB0aGUgZXZlbnQgbGlzdGVuZXJzIGNvbGxlY3Rpb24uXHJcbiAqIElmIGFsbCBsaXN0ZW5lcnMgZm9yIHRoZSBldmVudCBhcmUgcmVtb3ZlZCwgdGhlIGV2ZW50IGtleSBpcyBkZWxldGVkIGZyb20gdGhlIGNvbGxlY3Rpb24uXHJcbiAqXHJcbiAqIEBwYXJhbSB7T2JqZWN0fSBsaXN0ZW5lciAtIFRoZSBsaXN0ZW5lciB0byBiZSByZW1vdmVkLlxyXG4gKi9cclxuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcclxuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XHJcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XHJcbiAgICBlbHNlIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogUmVtb3ZlcyBldmVudCBsaXN0ZW5lcnMgZm9yIHRoZSBzcGVjaWZpZWQgZXZlbnQgbmFtZXMuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWUgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXHJcbiAqIEBwYXJhbSB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lcyAtIEFkZGl0aW9uYWwgZXZlbnQgbmFtZXMgdG8gcmVtb3ZlIGxpc3RlbmVycyBmb3IuXHJcbiAqIEByZXR1cm4ge3VuZGVmaW5lZH1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xyXG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xyXG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4gZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSkpO1xyXG59XHJcbi8qKlxyXG4gKiBSZW1vdmVzIGFsbCBldmVudCBsaXN0ZW5lcnMuXHJcbiAqXHJcbiAqIEBmdW5jdGlvbiBPZmZBbGxcclxuICogQHJldHVybnMge3ZvaWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkgeyBldmVudExpc3RlbmVycy5jbGVhcigpOyB9XHJcblxyXG4vKipcclxuICogRW1pdHMgYW4gZXZlbnQgdXNpbmcgdGhlIGdpdmVuIGV2ZW50IG5hbWUuXHJcbiAqXHJcbiAqIEBwYXJhbSB7V2FpbHNFdmVudH0gZXZlbnQgLSBUaGUgbmFtZSBvZiB0aGUgZXZlbnQgdG8gZW1pdC5cclxuICogQHJldHVybnMge2FueX0gLSBUaGUgcmVzdWx0IG9mIHRoZSBlbWl0dGVkIGV2ZW50LlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVtaXQoZXZlbnQpIHsgcmV0dXJuIGNhbGwoRW1pdE1ldGhvZCwgZXZlbnQpOyB9XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vdHlwZXNcIikuTWVzc2FnZURpYWxvZ09wdGlvbnN9IE1lc3NhZ2VEaWFsb2dPcHRpb25zXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL3R5cGVzXCIpLk9wZW5EaWFsb2dPcHRpb25zfSBPcGVuRGlhbG9nT3B0aW9uc1xyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi90eXBlc1wiKS5TYXZlRGlhbG9nT3B0aW9uc30gU2F2ZURpYWxvZ09wdGlvbnNcclxuICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XHJcblxyXG4vLyBEZWZpbmUgY29uc3RhbnRzIGZyb20gdGhlIGBtZXRob2RzYCBvYmplY3QgaW4gVGl0bGUgQ2FzZVxyXG5jb25zdCBEaWFsb2dJbmZvID0gMDtcclxuY29uc3QgRGlhbG9nV2FybmluZyA9IDE7XHJcbmNvbnN0IERpYWxvZ0Vycm9yID0gMjtcclxuY29uc3QgRGlhbG9nUXVlc3Rpb24gPSAzO1xyXG5jb25zdCBEaWFsb2dPcGVuRmlsZSA9IDQ7XHJcbmNvbnN0IERpYWxvZ1NhdmVGaWxlID0gNTtcclxuXHJcbmNvbnN0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZywgJycpO1xyXG5jb25zdCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwKCk7XHJcblxyXG4vKipcclxuICogR2VuZXJhdGVzIGEgdW5pcXVlIGlkIHRoYXQgaXMgbm90IHByZXNlbnQgaW4gZGlhbG9nUmVzcG9uc2VzLlxyXG4gKiBAcmV0dXJucyB7c3RyaW5nfSB1bmlxdWUgaWRcclxuICovXHJcbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XHJcbiAgICBsZXQgcmVzdWx0O1xyXG4gICAgZG8ge1xyXG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xyXG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhIGRpYWxvZyBvZiBzcGVjaWZpZWQgdHlwZSB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge251bWJlcn0gdHlwZSAtIHR5cGUgb2YgZGlhbG9nXHJcbiAqIEBwYXJhbSB7b2JqZWN0fSBvcHRpb25zIC0gb3B0aW9ucyBmb3IgdGhlIGRpYWxvZ1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZX0gcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggcmVzdWx0IG9mIGRpYWxvZ1xyXG4gKi9cclxuZnVuY3Rpb24gZGlhbG9nKHR5cGUsIG9wdGlvbnMgPSB7fSkge1xyXG4gICAgY29uc3QgaWQgPSBnZW5lcmF0ZUlEKCk7XHJcbiAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcclxuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xyXG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xyXG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogSGFuZGxlcyB0aGUgY2FsbGJhY2sgZnJvbSBhIGRpYWxvZy5cclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGlkIC0gVGhlIElEIG9mIHRoZSBkaWFsb2cgcmVzcG9uc2UuXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBkYXRhIC0gVGhlIGRhdGEgcmVjZWl2ZWQgZnJvbSB0aGUgZGlhbG9nLlxyXG4gKiBAcGFyYW0ge2Jvb2xlYW59IGlzSlNPTiAtIEZsYWcgaW5kaWNhdGluZyB3aGV0aGVyIHRoZSBkYXRhIGlzIGluIEpTT04gZm9ybWF0LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt1bmRlZmluZWR9XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGlmIChwKSB7XHJcbiAgICAgICAgaWYgKGlzSlNPTikge1xyXG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxiYWNrIGZ1bmN0aW9uIGZvciBoYW5kbGluZyBlcnJvcnMgaW4gZGlhbG9nLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gaWQgLSBUaGUgaWQgb2YgdGhlIGRpYWxvZyByZXNwb25zZS5cclxuICogQHBhcmFtIHtzdHJpbmd9IG1lc3NhZ2UgLSBUaGUgZXJyb3IgbWVzc2FnZS5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XHJcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgaWYgKHApIHtcclxuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5cclxuXHJcbi8vIFJlcGxhY2UgYG1ldGhvZHNgIHdpdGggY29uc3RhbnRzIGluIFRpdGxlIENhc2VcclxuXHJcbi8qKlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zIC0gRGlhbG9nIG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gLSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgY29uc3QgSW5mbyA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7XHJcblxyXG4vKipcclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9ucyAtIERpYWxvZyBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IC0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IFdhcm5pbmcgPSAob3B0aW9ucykgPT4gZGlhbG9nKERpYWxvZ1dhcm5pbmcsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBFcnJvciA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nRXJyb3IsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSAtIFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBjb25zdCBRdWVzdGlvbiA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7T3BlbkRpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHNlbGVjdGVkIGZpbGUgb3IgbGlzdCBvZiBmaWxlcy4gUmV0dXJucyBibGFuayBzdHJpbmcgaWYgbm8gZmlsZSBpcyBzZWxlY3RlZC5cclxuICovXHJcbmV4cG9ydCBjb25zdCBPcGVuRmlsZSA9IChvcHRpb25zKSA9PiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpO1xyXG5cclxuLyoqXHJcbiAqIEBwYXJhbSB7U2F2ZURpYWxvZ09wdGlvbnN9IG9wdGlvbnMgLSBEaWFsb2cgb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlLiBSZXR1cm5zIGJsYW5rIHN0cmluZyBpZiBubyBmaWxlIGlzIHNlbGVjdGVkLlxyXG4gKi9cclxuZXhwb3J0IGNvbnN0IFNhdmVGaWxlID0gKG9wdGlvbnMpID0+IGRpYWxvZyhEaWFsb2dTYXZlRmlsZSwgb3B0aW9ucyk7XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5jb25zdCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5Db250ZXh0TWVudSwgJycpO1xyXG5jb25zdCBDb250ZXh0TWVudU9wZW4gPSAwO1xyXG5cclxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkLCB4LCB5LCBkYXRhKSB7XHJcbiAgICB2b2lkIGNhbGwoQ29udGV4dE1lbnVPcGVuLCB7aWQsIHgsIHksIGRhdGF9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHNldHVwQ29udGV4dE1lbnVzKCkge1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcclxufVxyXG5cclxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50KSB7XHJcbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxyXG4gICAgbGV0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XHJcbiAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIik7XHJcbiAgICBjdXN0b21Db250ZXh0TWVudSA9IGN1c3RvbUNvbnRleHRNZW51ID8gY3VzdG9tQ29udGV4dE1lbnUudHJpbSgpIDogXCJcIjtcclxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgbGV0IGN1c3RvbUNvbnRleHRNZW51RGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xyXG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgY3VzdG9tQ29udGV4dE1lbnVEYXRhKTtcclxuICAgICAgICByZXR1cm5cclxuICAgIH1cclxuXHJcbiAgICBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KTtcclxufVxyXG5cclxuXHJcbi8qXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG5cclxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcclxuKi9cclxuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCkge1xyXG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XHJcbiAgICBpZiAoREVCVUcpIHtcclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxyXG4gICAgY29uc3QgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcclxuICAgIGNvbnN0IGNvbXB1dGVkU3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcclxuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XHJcbiAgICBzd2l0Y2ggKGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbikge1xyXG4gICAgICAgIGNhc2UgXCJzaG93XCI6XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICBjYXNlIFwiaGlkZVwiOlxyXG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgZGVmYXVsdDpcclxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcclxuICAgICAgICAgICAgaWYgKGVsZW1lbnQuaXNDb250ZW50RWRpdGFibGUpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG5cclxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxyXG4gICAgICAgICAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XHJcbiAgICAgICAgICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IChzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKVxyXG4gICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcclxuICAgICAgICAgICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xyXG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcclxuICAgICAgICAgICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbiAgICAgICAgICAgIGlmIChlbGVtZW50LnRhZ05hbWUgPT09IFwiSU5QVVRcIiB8fCBlbGVtZW50LnRhZ05hbWUgPT09IFwiVEVYVEFSRUFcIikge1xyXG4gICAgICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIWVsZW1lbnQucmVhZE9ubHkgJiYgIWVsZW1lbnQuZGlzYWJsZWQpKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICB9XHJcbn1cclxuIiwgIlxyXG5pbXBvcnQge0VtaXQsIFdhaWxzRXZlbnR9IGZyb20gXCIuL2V2ZW50c1wiO1xyXG5pbXBvcnQge1F1ZXN0aW9ufSBmcm9tIFwiLi9kaWFsb2dzXCI7XHJcbmltcG9ydCB7V2luZG93TWV0aG9kcywgR2V0fSBmcm9tIFwiLi93aW5kb3dcIjtcclxuXHJcbi8qKlxyXG4gKiBTZW5kcyBhbiBldmVudCB3aXRoIHRoZSBnaXZlbiBuYW1lIGFuZCBvcHRpb25hbCBkYXRhLlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lIC0gVGhlIG5hbWUgb2YgdGhlIGV2ZW50IHRvIHNlbmQuXHJcbiAqIEBwYXJhbSB7YW55fSBbZGF0YT1udWxsXSAtIE9wdGlvbmFsIGRhdGEgdG8gc2VuZCBhbG9uZyB3aXRoIHRoZSBldmVudC5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xyXG4gICAgbGV0IGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKTtcclxuICAgIEVtaXQoZXZlbnQpO1xyXG59XHJcblxyXG4vKipcclxuICogQWRkcyBldmVudCBsaXN0ZW5lcnMgdG8gZWxlbWVudHMgd2l0aCBgd21sLWV2ZW50YCBhdHRyaWJ1dGUuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcclxuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0nKTtcclxuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcclxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50Jyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgRGV0YWNoZWQ6IGZhbHNlLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXHJcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuXHJcbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGxzIGEgbWV0aG9kIG9uIHRoZSB3aW5kb3cgb2JqZWN0LlxyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gbWV0aG9kIC0gVGhlIG5hbWUgb2YgdGhlIG1ldGhvZCB0byBjYWxsIG9uIHRoZSB3aW5kb3cgb2JqZWN0LlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcclxuICAgIC8vIFRPRE86IE1ha2UgdGhpcyBhIHBhcmFtZXRlciFcclxuICAgIGxldCB3aW5kb3dOYW1lID0gJyc7XHJcbiAgICBsZXQgdGFyZ2V0V2luZG93ID0gR2V0KCcnKTtcclxuICAgIGxldCBtZXRob2RNYXAgPSBXaW5kb3dNZXRob2RzKHRhcmdldFdpbmRvdyk7XHJcbiAgICBpZiAoIW1ldGhvZE1hcC5oYXMobWV0aG9kKSkge1xyXG4gICAgICAgIGNvbnNvbGUubG9nKFwiV2luZG93IG1ldGhvZCBcIiArIG1ldGhvZCArIFwiIG5vdCBmb3VuZFwiKTtcclxuICAgIH1cclxuICAgIG1ldGhvZE1hcC5nZXQobWV0aG9kKSgpO1xyXG59XHJcblxyXG4vKipcclxuICogQWRkcyB3aW5kb3cgbGlzdGVuZXJzIGZvciBlbGVtZW50cyB3aXRoIHRoZSAnd21sLXdpbmRvdycgYXR0cmlidXRlLlxyXG4gKiBSZW1vdmVzIGFueSBleGlzdGluZyBsaXN0ZW5lcnMgYmVmb3JlIGFkZGluZyBuZXcgb25lcy5cclxuICpcclxuICogQHJldHVybiB7dm9pZH1cclxuICovXHJcbmZ1bmN0aW9uIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpIHtcclxuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC13aW5kb3ddJyk7XHJcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XHJcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC13aW5kb3cnKTtcclxuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XHJcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcclxuXHJcbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xyXG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xyXG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfSk7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogQWRkcyBhIGxpc3RlbmVyIHRvIGVsZW1lbnRzIHdpdGggdGhlICd3bWwtb3BlbnVybCcgYXR0cmlidXRlLlxyXG4gKiBXaGVuIHRoZSBzcGVjaWZpZWQgdHJpZ2dlciBldmVudCBpcyBmaXJlZCBvbiBhbnkgb2YgdGhlc2UgZWxlbWVudHMsXHJcbiAqIHRoZSBsaXN0ZW5lciB3aWxsIG9wZW4gdGhlIFVSTCBzcGVjaWZpZWQgYnkgdGhlICd3bWwtb3BlbnVybCcgYXR0cmlidXRlLlxyXG4gKiBJZiBhICd3bWwtY29uZmlybScgYXR0cmlidXRlIGlzIHByb3ZpZGVkLCBhIGNvbmZpcm1hdGlvbiBkaWFsb2cgd2lsbCBiZSBkaXNwbGF5ZWQsXHJcbiAqIGFuZCB0aGUgVVJMIHdpbGwgb25seSBiZSBvcGVuZWQgaWYgdGhlIHVzZXIgY29uZmlybXMuXHJcbiAqXHJcbiAqIEByZXR1cm4ge3ZvaWR9XHJcbiAqL1xyXG5mdW5jdGlvbiBhZGRXTUxPcGVuQnJvd3Nlckxpc3RlbmVyKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLW9wZW51cmxdJyk7XHJcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XHJcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgdm9pZCB3YWlscy5Ccm93c2VyLk9wZW5VUkwodXJsKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB2b2lkIHdhaWxzLkJyb3dzZXIuT3BlblVSTCh1cmwpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogUmVsb2FkcyB0aGUgV01MIHBhZ2UgYnkgYWRkaW5nIG5lY2Vzc2FyeSBldmVudCBsaXN0ZW5lcnMgYW5kIGJyb3dzZXIgbGlzdGVuZXJzLlxyXG4gKlxyXG4gKiBAcmV0dXJuIHt2b2lkfVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIHJlbG9hZFdNTCgpIHtcclxuICAgIGNvbnNvbGUubG9nKFwiUmVsb2FkaW5nIFdNTFwiKTtcclxuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XHJcbiAgICBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKTtcclxuICAgIGFkZFdNTE9wZW5Ccm93c2VyTGlzdGVuZXIoKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxubGV0IGZsYWdzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gY29udmVydFRvTWFwKG9iaikge1xyXG4gICAgY29uc3QgbWFwID0gbmV3IE1hcCgpO1xyXG5cclxuICAgIGZvciAoY29uc3QgW2tleSwgdmFsdWVdIG9mIE9iamVjdC5lbnRyaWVzKG9iaikpIHtcclxuICAgICAgICBpZiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyAmJiB2YWx1ZSAhPT0gbnVsbCkge1xyXG4gICAgICAgICAgICBtYXAuc2V0KGtleSwgY29udmVydFRvTWFwKHZhbHVlKSk7IC8vIFJlY3Vyc2l2ZWx5IGNvbnZlcnQgbmVzdGVkIG9iamVjdFxyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIG1hcC5zZXQoa2V5LCB2YWx1ZSk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBtYXA7XHJcbn1cclxuXHJcbmZldGNoKFwiL3dhaWxzL2ZsYWdzXCIpLnRoZW4oKHJlc3BvbnNlKSA9PiB7XHJcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xyXG4gICAgICAgIGZsYWdzID0gY29udmVydFRvTWFwKGRhdGEpO1xyXG4gICAgfSk7XHJcbn0pO1xyXG5cclxuXHJcbmZ1bmN0aW9uIGdldFZhbHVlRnJvbU1hcChrZXlTdHJpbmcpIHtcclxuICAgIGNvbnN0IGtleXMgPSBrZXlTdHJpbmcuc3BsaXQoJy4nKTtcclxuICAgIGxldCB2YWx1ZSA9IGZsYWdzO1xyXG5cclxuICAgIGZvciAoY29uc3Qga2V5IG9mIGtleXMpIHtcclxuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBNYXApIHtcclxuICAgICAgICAgICAgdmFsdWUgPSB2YWx1ZS5nZXQoa2V5KTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICB2YWx1ZSA9IHZhbHVlW2tleV07XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAodmFsdWUgPT09IHVuZGVmaW5lZCkge1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIHZhbHVlO1xyXG59XHJcblxyXG4vKipcclxuICogUmV0cmlldmVzIHRoZSB2YWx1ZSBhc3NvY2lhdGVkIHdpdGggdGhlIHNwZWNpZmllZCBrZXkgZnJvbSB0aGUgZmxhZyBtYXAuXHJcbiAqXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBrZXlTdHJpbmcgLSBUaGUga2V5IHRvIHJldHJpZXZlIHRoZSB2YWx1ZSBmb3IuXHJcbiAqIEByZXR1cm4geyp9IC0gVGhlIHZhbHVlIGFzc29jaWF0ZWQgd2l0aCB0aGUgc3BlY2lmaWVkIGtleS5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleVN0cmluZykge1xyXG4gICAgcmV0dXJuIGdldFZhbHVlRnJvbU1hcChrZXlTdHJpbmcpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5pbXBvcnQge2ludm9rZSwgSXNXaW5kb3dzfSBmcm9tIFwiLi9zeXN0ZW1cIjtcclxuaW1wb3J0IHtHZXRGbGFnfSBmcm9tIFwiLi9mbGFnc1wiO1xyXG5cclxubGV0IHNob3VsZERyYWcgPSBmYWxzZTtcclxubGV0IHJlc2l6ZUVkZ2UgPSBudWxsO1xyXG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XHJcbmxldCBkZWZhdWx0Q3Vyc29yID0gXCJhdXRvXCI7XHJcbndpbmRvdy5fd2FpbHMgPSB3aW5kb3cuX3dhaWxzIHx8IHt9O1xyXG53aW5kb3cuX3dhaWxzLnNldFJlc2l6YWJsZSA9IHNldFJlc2l6YWJsZTtcclxud2luZG93Ll93YWlscy5lbmREcmFnID0gZW5kRHJhZztcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBkcmFnVGVzdChlKSB7XHJcbiAgICBsZXQgdmFsID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLXdlYmtpdC1hcHAtcmVnaW9uXCIpO1xyXG4gICAgaWYgKHZhbCAmJiB2YWwudHJpbSgpICE9PSBcImRyYWdcIiB8fCBlLmJ1dHRvbnMgIT09IDEpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gZS5kZXRhaWwgPT09IDE7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBzZXR1cERyYWcoKSB7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgb25Nb3VzZURvd24pO1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIG9uTW91c2VNb3ZlKTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgb25Nb3VzZVVwKTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHNldFJlc2l6YWJsZSh2YWx1ZSkge1xyXG4gICAgcmVzaXphYmxlID0gdmFsdWU7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBlbmREcmFnKCkge1xyXG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XHJcbiAgICBzaG91bGREcmFnID0gZmFsc2U7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHRlc3RSZXNpemUoKSB7XHJcbiAgICBpZiggcmVzaXplRWRnZSApIHtcclxuICAgICAgICBpbnZva2UoYHJlc2l6ZToke3Jlc2l6ZUVkZ2V9YCk7XHJcbiAgICAgICAgcmV0dXJuIHRydWVcclxuICAgIH1cclxuICAgIHJldHVybiBmYWxzZTtcclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xyXG4gICAgaWYoSXNXaW5kb3dzKCkgJiYgdGVzdFJlc2l6ZSgpIHx8IGRyYWdUZXN0KGUpKSB7XHJcbiAgICAgICAgc2hvdWxkRHJhZyA9ICEhaXNWYWxpZERyYWcoZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGlzVmFsaWREcmFnKGUpIHtcclxuICAgIC8vIElnbm9yZSBkcmFnIG9uIHNjcm9sbGJhcnNcclxuICAgIHJldHVybiAhKGUub2Zmc2V0WCA+IGUudGFyZ2V0LmNsaWVudFdpZHRoIHx8IGUub2Zmc2V0WSA+IGUudGFyZ2V0LmNsaWVudEhlaWdodCk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VVcChlKSB7XHJcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xyXG4gICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcclxuICAgICAgICBlbmREcmFnKCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IgPSBkZWZhdWx0Q3Vyc29yKSB7XHJcbiAgICBkb2N1bWVudC5kb2N1bWVudEVsZW1lbnQuc3R5bGUuY3Vyc29yID0gY3Vyc29yO1xyXG4gICAgcmVzaXplRWRnZSA9IGN1cnNvcjtcclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZU1vdmUoZSkge1xyXG4gICAgc2hvdWxkRHJhZyA9IGNoZWNrRHJhZyhlKTtcclxuICAgIGlmIChJc1dpbmRvd3MoKSAmJiByZXNpemFibGUpIHtcclxuICAgICAgICBoYW5kbGVSZXNpemUoZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNoZWNrRHJhZyhlKSB7XHJcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xyXG4gICAgaWYoc2hvdWxkRHJhZyAmJiBtb3VzZVByZXNzZWQgPiAwKSB7XHJcbiAgICAgICAgaW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcbiAgICByZXR1cm4gc2hvdWxkRHJhZztcclxufVxyXG5cclxuZnVuY3Rpb24gaGFuZGxlUmVzaXplKGUpIHtcclxuICAgIGxldCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xyXG4gICAgbGV0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xyXG5cclxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhc1xyXG4gICAgbGV0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xyXG5cclxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XHJcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgbGV0IHRvcEJvcmRlciA9IGUuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcclxuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcblxyXG4gICAgLy8gQWRqdXN0IGZvciBjb3JuZXJzXHJcbiAgICBsZXQgcmlnaHRDb3JuZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGxldCBsZWZ0Q29ybmVyID0gZS5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xyXG4gICAgbGV0IHRvcENvcm5lciA9IGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBsZXQgYm90dG9tQ29ybmVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcclxuXHJcbiAgICAvLyBJZiB3ZSBhcmVuJ3Qgb24gYW4gZWRnZSwgYnV0IHdlcmUsIHJlc2V0IHRoZSBjdXJzb3IgdG8gZGVmYXVsdFxyXG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgcmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgc2V0UmVzaXplKCk7XHJcbiAgICB9XHJcbiAgICAvLyBBZGp1c3RlZCBmb3IgY29ybmVyIGFyZWFzXHJcbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzdy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIHRvcENvcm5lcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Qm9yZGVyKSBzZXRSZXNpemUoXCJ3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAocmlnaHRCb3JkZXIpIHNldFJlc2l6ZShcImUtcmVzaXplXCIpO1xyXG59XHJcbiIsICIvKlxyXG4gXyAgICAgX18gICAgIF8gX19cclxufCB8ICAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcblxyXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSAnLi9Ad2FpbHNpby9ydW50aW1lL2NsaXBib2FyZCc7XHJcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gJy4vQHdhaWxzaW8vcnVudGltZS9hcHBsaWNhdGlvbic7XHJcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSAnLi9Ad2FpbHNpby9ydW50aW1lL3NjcmVlbnMnO1xyXG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSAnLi9Ad2FpbHNpby9ydW50aW1lL3N5c3RlbSc7XHJcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSAnLi9Ad2FpbHNpby9ydW50aW1lL2Jyb3dzZXInO1xyXG5pbXBvcnQgKiBhcyBXaW5kb3cgZnJvbSAnLi9Ad2FpbHNpby9ydW50aW1lL3dpbmRvdyc7XHJcbmltcG9ydCB7UGx1Z2luLCBDYWxsLCBlcnJvckhhbmRsZXIgYXMgY2FsbEVycm9ySGFuZGxlciwgcmVzdWx0SGFuZGxlciBhcyBjYWxsUmVzdWx0SGFuZGxlciwgQnlJRCwgQnlOYW1lfSBmcm9tIFwiLi9Ad2FpbHNpby9ydW50aW1lL2NhbGxzXCI7XHJcbmltcG9ydCB7Y2xpZW50SWR9IGZyb20gJy4vQHdhaWxzaW8vcnVudGltZS9ydW50aW1lJztcclxuaW1wb3J0IHtkaXNwYXRjaFdhaWxzRXZlbnQsIEVtaXQsIE9mZiwgT2ZmQWxsLCBPbiwgT25jZSwgT25NdWx0aXBsZX0gZnJvbSBcIi4vQHdhaWxzaW8vcnVudGltZS9ldmVudHNcIjtcclxuaW1wb3J0IHtkaWFsb2dDYWxsYmFjaywgZGlhbG9nRXJyb3JDYWxsYmFjaywgRXJyb3IsIEluZm8sIE9wZW5GaWxlLCBRdWVzdGlvbiwgU2F2ZUZpbGUsIFdhcm5pbmd9IGZyb20gXCIuL0B3YWlsc2lvL3J1bnRpbWUvZGlhbG9nc1wiO1xyXG5pbXBvcnQge3NldHVwQ29udGV4dE1lbnVzfSBmcm9tICcuL0B3YWlsc2lvL3J1bnRpbWUvY29udGV4dG1lbnUnO1xyXG5pbXBvcnQge3JlbG9hZFdNTH0gZnJvbSAnLi9Ad2FpbHNpby9ydW50aW1lL3dtbCc7XHJcbmltcG9ydCB7c2V0dXBEcmFnLCBlbmREcmFnLCBzZXRSZXNpemFibGV9IGZyb20gJy4vQHdhaWxzaW8vcnVudGltZS9kcmFnJztcclxuXHJcbndpbmRvdy53YWlscyA9IHtcclxuICAgIC4uLm5ld1J1bnRpbWUobnVsbCksXHJcbiAgICBjbGllbnRJZDogY2xpZW50SWQsXHJcbn07XHJcblxyXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcclxud2luZG93Ll93YWlscyA9IHtcclxuICAgIGRpYWxvZ0NhbGxiYWNrLFxyXG4gICAgZGlhbG9nRXJyb3JDYWxsYmFjayxcclxuICAgIGRpc3BhdGNoV2FpbHNFdmVudCxcclxuICAgIGNhbGxFcnJvckhhbmRsZXIsXHJcbiAgICBjYWxsUmVzdWx0SGFuZGxlcixcclxuICAgIGVuZERyYWcsXHJcbiAgICBzZXRSZXNpemFibGUsXHJcbn07XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZSh3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4ge1xyXG4gICAgICAgIENsaXBib2FyZDoge1xyXG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcclxuICAgICAgICB9LFxyXG4gICAgICAgIEFwcGxpY2F0aW9uOiB7XHJcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgU3lzdGVtLFxyXG4gICAgICAgIFNjcmVlbnMsXHJcbiAgICAgICAgQnJvd3NlcixcclxuICAgICAgICBDYWxsOiB7XHJcbiAgICAgICAgICAgIENhbGwsXHJcbiAgICAgICAgICAgIEJ5SUQsXHJcbiAgICAgICAgICAgIEJ5TmFtZSxcclxuICAgICAgICAgICAgUGx1Z2luLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgV01MOiB7XHJcbiAgICAgICAgICAgIFJlbG9hZDogcmVsb2FkV01MLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgRGlhbG9nOiB7XHJcbiAgICAgICAgICAgIEluZm8sXHJcbiAgICAgICAgICAgIFdhcm5pbmcsXHJcbiAgICAgICAgICAgIEVycm9yLFxyXG4gICAgICAgICAgICBRdWVzdGlvbixcclxuICAgICAgICAgICAgT3BlbkZpbGUsXHJcbiAgICAgICAgICAgIFNhdmVGaWxlLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgRXZlbnRzOiB7XHJcbiAgICAgICAgICAgIEVtaXQsXHJcbiAgICAgICAgICAgIE9uLFxyXG4gICAgICAgICAgICBPbmNlLFxyXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxyXG4gICAgICAgICAgICBPZmYsXHJcbiAgICAgICAgICAgIE9mZkFsbCxcclxuICAgICAgICB9LFxyXG4gICAgICAgIFdpbmRvdzoge1xyXG4gICAgICAgICAgICAuLi5XaW5kb3cuR2V0KCcnKVxyXG4gICAgICAgIH0sXHJcbiAgICB9O1xyXG59XHJcblxyXG5zZXR1cENvbnRleHRNZW51cygpO1xyXG5zZXR1cERyYWcoKTtcclxuXHJcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoXCJET01Db250ZW50TG9hZGVkXCIs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LGNBQWM7QUFBVyxnQkFBVSxXQUFXO0FBQUEsYUFDOUMsYUFBYTtBQUFhLGdCQUFVLFdBQVc7QUFBQSxhQUMvQztBQUFZLGdCQUFVLFVBQVU7QUFBQSxhQUNoQztBQUFXLGdCQUFVLFVBQVU7QUFBQSxhQUMvQjtBQUFjLGdCQUFVLFVBQVU7QUFBQSxhQUNsQztBQUFhLGdCQUFVLFVBQVU7QUFBQSxFQUM5Qzs7O0FDbEdBLFNBQU8sUUFBUTtBQUFBLElBQ1gsR0FBRyxXQUFXLElBQUk7QUFBQSxJQUNsQjtBQUFBLEVBQ0o7QUFHQSxTQUFPLFNBQVM7QUFBQSxJQUNaO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBLGNBQUFDO0FBQUEsRUFDSjtBQUVPLFdBQVMsV0FBVyxZQUFZO0FBQ25DLFdBQU87QUFBQSxNQUNILFdBQVc7QUFBQSxRQUNQLEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQSxhQUFhO0FBQUEsUUFDVCxHQUFHO0FBQUEsTUFDUDtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0EsTUFBTTtBQUFBLFFBQ0Y7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxLQUFLO0FBQUEsUUFDRCxRQUFRO0FBQUEsTUFDWjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQSxPQUFBQztBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSixHQUFVLElBQUksRUFBRTtBQUFBLE1BQ3BCO0FBQUEsSUFDSjtBQUFBLEVBQ0o7QUFFQSxvQkFBa0I7QUFDbEIsWUFBVTtBQUVWLFdBQVMsaUJBQWlCLG9CQUFvQixXQUFXO0FBQ3JELGNBQVU7QUFBQSxFQUNkLENBQUM7IiwKICAibmFtZXMiOiBbInNpemUiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgIndpZHRoIiwgImhlaWdodCIsICJyZXNpemFibGUiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiRXJyb3IiLCAiY2FsbCIsICJzZXRSZXNpemFibGUiLCAic2V0UmVzaXphYmxlIiwgIkVycm9yIl0KfQo=
