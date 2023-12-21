(() => {
  var __defProp = Object.defineProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };

  // desktop/clipboard.js
  var clipboard_exports = {};
  __export(clipboard_exports, {
    SetText: () => SetText,
    Text: () => Text
  });

  // node_modules/nanoid/non-secure/index.js
  var urlAlphabet = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict";
  var nanoid = (size = 21) => {
    let id = "";
    let i = size;
    while (i--) {
      id += urlAlphabet[Math.random() * 64 | 0];
    }
    return id;
  };

  // desktop/runtime.js
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
  function newRuntimeCallerWithID(object, windowName) {
    return function(method, args = null) {
      return runtimeCallWithID(object, method, windowName, args);
    };
  }

  // desktop/clipboard.js
  var call = newRuntimeCallerWithID(objectNames.Clipboard);
  var ClipboardSetText = 0;
  var ClipboardText = 1;
  function SetText(text) {
    return call(ClipboardSetText, { text });
  }
  function Text() {
    return call(ClipboardText);
  }

  // desktop/application.js
  var application_exports = {};
  __export(application_exports, {
    Hide: () => Hide,
    Quit: () => Quit,
    Show: () => Show
  });
  var call2 = newRuntimeCallerWithID(objectNames.Application);
  var methods = {
    Hide: 0,
    Show: 1,
    Quit: 2
  };
  function Hide() {
    return call2(methods.Hide);
  }
  function Show() {
    return call2(methods.Show);
  }
  function Quit() {
    return call2(methods.Quit);
  }

  // desktop/screens.js
  var screens_exports = {};
  __export(screens_exports, {
    GetAll: () => GetAll,
    GetCurrent: () => GetCurrent,
    GetPrimary: () => GetPrimary
  });
  var call3 = newRuntimeCallerWithID(objectNames.Screens);
  var ScreensGetAll = 0;
  var ScreensGetPrimary = 1;
  var ScreensGetCurrent = 2;
  function GetAll() {
    return call3(ScreensGetAll);
  }
  function GetPrimary() {
    return call3(ScreensGetPrimary);
  }
  function GetCurrent() {
    return call3(ScreensGetCurrent);
  }

  // desktop/system.js
  var system_exports = {};
  __export(system_exports, {
    IsDarkMode: () => IsDarkMode
  });
  var call4 = newRuntimeCallerWithID(objectNames.System);
  var SystemIsDarkMode = 0;
  function IsDarkMode() {
    return call4(SystemIsDarkMode);
  }

  // desktop/browser.js
  var browser_exports = {};
  __export(browser_exports, {
    OpenURL: () => OpenURL
  });
  var call5 = newRuntimeCallerWithID(objectNames.Browser);
  var BrowserOpenURL = 0;
  function OpenURL(url) {
    return call5(BrowserOpenURL, { url });
  }

  // desktop/calls.js
  var call6 = newRuntimeCallerWithID(objectNames.Call);
  var CallBinding = 0;
  var callResponses = /* @__PURE__ */ new Map();
  function generateID() {
    let result;
    do {
      result = nanoid();
    } while (callResponses.has(result));
    return result;
  }
  function callCallback(id, data, isJSON) {
    let p = callResponses.get(id);
    if (p) {
      if (isJSON) {
        p.resolve(JSON.parse(data));
      } else {
        p.resolve(data);
      }
      callResponses.delete(id);
    }
  }
  function callErrorCallback(id, message) {
    let p = callResponses.get(id);
    if (p) {
      p.reject(message);
      callResponses.delete(id);
    }
  }
  function callBinding(type, options) {
    return new Promise((resolve, reject) => {
      let id = generateID();
      options = options || {};
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
  function CallByName(name, ...args) {
    if (typeof name !== "string" || name.split(".").length !== 3) {
      throw new Error("CallByName requires a string in the format 'package.struct.method'");
    }
    let parts = name.split(".");
    return callBinding(CallBinding, {
      packageName: parts[0],
      structName: parts[1],
      methodName: parts[2],
      args
    });
  }
  function CallByID(methodID, ...args) {
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

  // desktop/window.js
  var WindowCenter = 0;
  var WindowSetTitle = 1;
  var WindowFullscreen = 2;
  var WindowUnFullscreen = 3;
  var WindowSetSize = 4;
  var WindowSize = 5;
  var WindowSetMaxSize = 6;
  var WindowSetMinSize = 7;
  var WindowSetAlwaysOnTop = 8;
  var WindowSetRelativePosition = 9;
  var WindowRelativePosition = 10;
  var WindowScreen = 11;
  var WindowHide = 12;
  var WindowMaximise = 13;
  var WindowUnMaximise = 14;
  var WindowToggleMaximise = 15;
  var WindowMinimise = 16;
  var WindowUnMinimise = 17;
  var WindowRestore = 18;
  var WindowShow = 19;
  var WindowClose = 20;
  var WindowSetBackgroundColour = 21;
  var WindowSetResizable = 22;
  var WindowWidth = 23;
  var WindowHeight = 24;
  var WindowZoomIn = 25;
  var WindowZoomOut = 26;
  var WindowZoomReset = 27;
  var WindowGetZoomLevel = 28;
  var WindowSetZoomLevel = 29;
  function newWindow(windowName) {
    let call10 = newRuntimeCallerWithID(objectNames.Window, windowName);
    return {
      /**
       * Centers the window.
       * @returns {Promise<void>}
       */
      Center: () => call10(WindowCenter),
      /**
       * Set the window title.
       * @param title
       * @returns {Promise<void>}
       */
      SetTitle: (title) => call10(WindowSetTitle, { title }),
      /**
       * Makes the window fullscreen.
       * @returns {Promise<void>}
       */
      Fullscreen: () => call10(WindowFullscreen),
      /**
       * Unfullscreen the window.
       * @returns {Promise<void>}
       */
      UnFullscreen: () => call10(WindowUnFullscreen),
      /**
       * Set the window size.
       * @param {number} width The window width
       * @param {number} height The window height
       * @returns {Promise<void>}
       */
      SetSize: (width, height) => call10(WindowSetSize, { width, height }),
      /**
       * Get the window size.
       * @returns {Promise<Size>} The window size
       */
      Size: () => call10(WindowSize),
      /**
       * Set the window maximum size.
       * @param {number} width
       * @param {number} height
       * @returns {Promise<void>}
       */
      SetMaxSize: (width, height) => call10(WindowSetMaxSize, { width, height }),
      /**
       * Set the window minimum size.
       * @param {number} width
       * @param {number} height
       * @returns {Promise<void>}
       */
      SetMinSize: (width, height) => call10(WindowSetMinSize, { width, height }),
      /**
       * Set window to be always on top.
       * @param {boolean} onTop Whether the window should be always on top
       * @returns {Promise<void>}
       */
      SetAlwaysOnTop: (onTop) => call10(WindowSetAlwaysOnTop, { alwaysOnTop: onTop }),
      /**
       * Set the window relative position.
       * @param {number} x
       * @param {number} y
       * @returns {Promise<void>}
       */
      SetRelativePosition: (x, y) => call10(WindowSetRelativePosition, { x, y }),
      /**
       * Get the window position.
       * @returns {Promise<Position>} The window position
       */
      RelativePosition: () => call10(WindowRelativePosition),
      /**
       * Get the screen the window is on.
       * @returns {Promise<Screen>}
       */
      Screen: () => call10(WindowScreen),
      /**
       * Hide the window
       * @returns {Promise<void>}
       */
      Hide: () => call10(WindowHide),
      /**
       * Maximise the window
       * @returns {Promise<void>}
       */
      Maximise: () => call10(WindowMaximise),
      /**
       * Show the window
       * @returns {Promise<void>}
       */
      Show: () => call10(WindowShow),
      /**
       * Close the window
       * @returns {Promise<void>}
       */
      Close: () => call10(WindowClose),
      /**
       * Toggle the window maximise state
       * @returns {Promise<void>}
       */
      ToggleMaximise: () => call10(WindowToggleMaximise),
      /**
       * Unmaximise the window
       * @returns {Promise<void>}
       */
      UnMaximise: () => call10(WindowUnMaximise),
      /**
       * Minimise the window
       * @returns {Promise<void>}
       */
      Minimise: () => call10(WindowMinimise),
      /**
       * Unminimise the window
       * @returns {Promise<void>}
       */
      UnMinimise: () => call10(WindowUnMinimise),
      /**
       * Restore the window
       * @returns {Promise<void>}
       */
      Restore: () => call10(WindowRestore),
      /**
       * Set the background colour of the window.
       * @param {number} r - A value between 0 and 255
       * @param {number} g - A value between 0 and 255
       * @param {number} b - A value between 0 and 255
       * @param {number} a - A value between 0 and 255
       * @returns {Promise<void>}
       */
      SetBackgroundColour: (r, g, b, a) => call10(WindowSetBackgroundColour, { r, g, b, a }),
      /**
       * Set whether the window can be resized or not
       * @param {boolean} resizable
       * @returns {Promise<void>}
       */
      SetResizable: (resizable2) => call10(WindowSetResizable, { resizable: resizable2 }),
      /**
       * Get the window width
       * @returns {Promise<number>}
       */
      Width: () => call10(WindowWidth),
      /**
       * Get the window height
       * @returns {Promise<number>}
       */
      Height: () => call10(WindowHeight),
      /**
       * Zoom in the window
       * @returns {Promise<void>}
       */
      ZoomIn: () => call10(WindowZoomIn),
      /**
       * Zoom out the window
       * @returns {Promise<void>}
       */
      ZoomOut: () => call10(WindowZoomOut),
      /**
       * Reset the window zoom
       * @returns {Promise<void>}
       */
      ZoomReset: () => call10(WindowZoomReset),
      /**
       * Get the window zoom
       * @returns {Promise<number>}
       */
      GetZoomLevel: () => call10(WindowGetZoomLevel),
      /**
       * Set the window zoom level
       * @param {number} zoomLevel
       * @returns {Promise<void>}
       */
      SetZoomLevel: (zoomLevel) => call10(WindowSetZoomLevel, { zoomLevel })
    };
  }

  // desktop/events.js
  var call7 = newRuntimeCallerWithID(objectNames.Events);
  var EventEmit = 0;
  var Listener = class {
    /**
     * Creates an instance of Listener.
     * @param {string} eventName
     * @param {(data: any) => void} callback
     * @param {number} maxCallbacks
     * @memberof Listener
     */
    constructor(eventName, callback, maxCallbacks) {
      this.eventName = eventName;
      this.maxCallbacks = maxCallbacks || -1;
      this.Callback = (data) => {
        callback(data);
        if (this.maxCallbacks === -1) {
          return false;
        }
        this.maxCallbacks -= 1;
        return this.maxCallbacks === 0;
      };
    }
  };
  var WailsEvent = class {
    /**
     * Creates an instance of WailsEvent.
     * @param {string} name - Name of the event
     * @param {any=null} data - Data associated with the event
     * @memberof WailsEvent
     */
    constructor(name, data = null) {
      this.name = name;
      this.data = data;
    }
  };
  var eventListeners = /* @__PURE__ */ new Map();
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
    if (listeners.length === 0) {
      eventListeners.delete(eventName);
    } else {
      eventListeners.set(eventName, listeners);
    }
  }
  function dispatchWailsEvent(event) {
    let listeners = eventListeners.get(event.name);
    if (listeners) {
      let toRemove = [];
      listeners.forEach((listener) => {
        let remove = listener.Callback(event);
        if (remove) {
          toRemove.push(listener);
        }
      });
      if (toRemove.length > 0) {
        listeners = listeners.filter((l) => !toRemove.includes(l));
        if (listeners.length === 0) {
          eventListeners.delete(event.name);
        } else {
          eventListeners.set(event.name, listeners);
        }
      }
    }
  }
  function Off(eventName, ...additionalEventNames) {
    let eventsToRemove = [eventName, ...additionalEventNames];
    eventsToRemove.forEach((eventName2) => {
      eventListeners.delete(eventName2);
    });
  }
  function OffAll() {
    eventListeners.clear();
  }
  function Emit(event) {
    return call7(EventEmit, event);
  }

  // desktop/dialogs.js
  var call8 = newRuntimeCallerWithID(objectNames.Dialog);
  var DialogInfo = 0;
  var DialogWarning = 1;
  var DialogError = 2;
  var DialogQuestion = 3;
  var DialogOpenFile = 4;
  var DialogSaveFile = 5;
  var dialogResponses = /* @__PURE__ */ new Map();
  function generateID2() {
    let result;
    do {
      result = nanoid();
    } while (dialogResponses.has(result));
    return result;
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
  function dialog(type, options) {
    return new Promise((resolve, reject) => {
      let id = generateID2();
      options = options || {};
      options["dialog-id"] = id;
      dialogResponses.set(id, { resolve, reject });
      call8(type, options).catch((error) => {
        reject(error);
        dialogResponses.delete(id);
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
    return dialog(DialogOpenFile, options);
  }
  function SaveFile(options) {
    return dialog(DialogSaveFile, options);
  }

  // desktop/contextmenu.js
  var call9 = newRuntimeCallerWithID(objectNames.ContextMenu);
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

  // desktop/wml.js
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
    if (wails.Window[method] === void 0) {
      console.log("Window method " + method + " not found");
    }
    wails.Window[method]();
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
    addWMLEventListeners();
    addWMLWindowListeners();
    addWMLOpenBrowserListener();
  }

  // desktop/invoke.js
  var invoke = function(input) {
    if (true) {
      chrome.webview.postMessage(input);
    } else {
      webkit.messageHandlers.external.postMessage(input);
    }
  };

  // desktop/flags.js
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

  // desktop/drag.js
  var shouldDrag = false;
  function dragTest(e) {
    let val = window.getComputedStyle(e.target).getPropertyValue("--webkit-app-region");
    if (val) {
      val = val.trim();
    }
    if (val !== "drag") {
      return false;
    }
    if (e.buttons !== 1) {
      return false;
    }
    return e.detail === 1;
  }
  function setupDrag() {
    window.addEventListener("mousedown", onMouseDown);
    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
  }
  var resizeEdge = null;
  var resizable = false;
  function setResizable(value) {
    resizable = value;
  }
  function testResize(e) {
    if (resizeEdge) {
      invoke("resize:" + resizeEdge);
      return true;
    }
    return false;
  }
  function onMouseDown(e) {
    if (true) {
      if (testResize()) {
        return;
      }
    }
    if (dragTest(e)) {
      if (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight) {
        return;
      }
      shouldDrag = true;
    } else {
      shouldDrag = false;
    }
  }
  function onMouseUp(e) {
    let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
    if (mousePressed > 0) {
      endDrag();
    }
  }
  function endDrag() {
    document.body.style.cursor = "default";
    shouldDrag = false;
  }
  function setResize(cursor) {
    document.documentElement.style.cursor = cursor || defaultCursor;
    resizeEdge = cursor;
  }
  function onMouseMove(e) {
    if (shouldDrag) {
      shouldDrag = false;
      let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
      if (mousePressed > 0) {
        invoke("drag");
      }
      return;
    }
    if (true) {
      if (resizable) {
        handleResize(e);
      }
    }
  }
  var defaultCursor = "auto";
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
    Capabilities: {},
    clientId
  };
  fetch("/wails/capabilities").then((response) => {
    response.json().then((data) => {
      window.wails.Capabilities = data;
    });
  });
  window._wails = {
    dialogCallback,
    dialogErrorCallback,
    dispatchWailsEvent,
    callCallback,
    callErrorCallback,
    endDrag,
    setResizable
  };
  function newRuntime(windowName) {
    return {
      Clipboard: {
        ...clipboard_exports
      },
      Application: {
        ...application_exports,
        GetWindowByName(windowName2) {
          return newRuntime(windowName2);
        }
      },
      System: system_exports,
      Screens: screens_exports,
      Browser: browser_exports,
      Call,
      CallByID,
      CallByName,
      Plugin,
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
        OffAll,
        WailsEvent
      },
      Window: newWindow(windowName)
    };
  }
  if (true) {
    console.log("Wails v3.0.0 Debug Mode Enabled");
  }
  setupContextMenus();
  setupDrag();
  document.addEventListener("DOMContentLoaded", function() {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAibm9kZV9tb2R1bGVzL25hbm9pZC9ub24tc2VjdXJlL2luZGV4LmpzIiwgImRlc2t0b3AvcnVudGltZS5qcyIsICJkZXNrdG9wL2FwcGxpY2F0aW9uLmpzIiwgImRlc2t0b3Avc2NyZWVucy5qcyIsICJkZXNrdG9wL3N5c3RlbS5qcyIsICJkZXNrdG9wL2Jyb3dzZXIuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL3dpbmRvdy5qcyIsICJkZXNrdG9wL2V2ZW50cy5qcyIsICJkZXNrdG9wL2RpYWxvZ3MuanMiLCAiZGVza3RvcC9jb250ZXh0bWVudS5qcyIsICJkZXNrdG9wL3dtbC5qcyIsICJkZXNrdG9wL2ludm9rZS5qcyIsICJkZXNrdG9wL2ZsYWdzLmpzIiwgImRlc2t0b3AvZHJhZy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCk7XG5cbmxldCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcbmxldCBDbGlwYm9hcmRUZXh0ID0gMTtcblxuLyoqXG4gKiBTZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKiBAcGFyYW0ge3N0cmluZ30gdGV4dCAtIHRleHQgdG8gc2V0IGluIHRoZSBjbGlwYm9hcmRcbiAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2V0VGV4dCh0ZXh0KSB7XG4gICAgcmV0dXJuIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn1cbiIsICJsZXQgdXJsQWxwaGFiZXQgPVxuICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcbmV4cG9ydCBsZXQgY3VzdG9tQWxwaGFiZXQgPSAoYWxwaGFiZXQsIGRlZmF1bHRTaXplID0gMjEpID0+IHtcbiAgcmV0dXJuIChzaXplID0gZGVmYXVsdFNpemUpID0+IHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIGxldCBpID0gc2l6ZVxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgIGlkICs9IGFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogYWxwaGFiZXQubGVuZ3RoKSB8IDBdXG4gICAgfVxuICAgIHJldHVybiBpZFxuICB9XG59XG5leHBvcnQgbGV0IG5hbm9pZCA9IChzaXplID0gMjEpID0+IHtcbiAgbGV0IGlkID0gJydcbiAgbGV0IGkgPSBzaXplXG4gIHdoaWxlIChpLS0pIHtcbiAgICBpZCArPSB1cmxBbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIDY0KSB8IDBdXG4gIH1cbiAgcmV0dXJuIGlkXG59XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxuY29uc3QgcnVudGltZVVSTCA9IHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XG4vLyBPYmplY3QgTmFtZXNcbmV4cG9ydCBjb25zdCBvYmplY3ROYW1lcyA9IHtcbiAgICBDYWxsOiAwLFxuICAgIENsaXBib2FyZDogMSxcbiAgICBBcHBsaWNhdGlvbjogMixcbiAgICBFdmVudHM6IDMsXG4gICAgQ29udGV4dE1lbnU6IDQsXG4gICAgRGlhbG9nOiA1LFxuICAgIFdpbmRvdzogNixcbiAgICBTY3JlZW5zOiA3LFxuICAgIFN5c3RlbTogOCxcbiAgICBCcm93c2VyOiA5LFxufVxuZXhwb3J0IGxldCBjbGllbnRJZCA9IG5hbm9pZCgpO1xuXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcbiAgICBsZXQgdXJsID0gbmV3IFVSTChydW50aW1lVVJMKTtcbiAgICBpZiggbWV0aG9kICkge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIH1cbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIGlmIChhcmdzKSB7XG4gICAgICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwiYXJnc1wiLCBKU09OLnN0cmluZ2lmeShhcmdzKSk7XG4gICAgfVxuICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy1jbGllbnQtaWRcIl0gPSBjbGllbnRJZDtcblxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuXG5mdW5jdGlvbiBydW50aW1lQ2FsbFdpdGhJRChvYmplY3RJRCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJvYmplY3RcIiwgb2JqZWN0SUQpO1xuICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwibWV0aG9kXCIsIG1ldGhvZCk7XG4gICAgbGV0IGZldGNoT3B0aW9ucyA9IHtcbiAgICAgICAgaGVhZGVyczoge30sXG4gICAgfTtcbiAgICBpZiAod2luZG93TmFtZSkge1xuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xuICAgIH1cbiAgICBpZiAoYXJncykge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xuICAgIH1cbiAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtY2xpZW50LWlkXCJdID0gY2xpZW50SWQ7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgZmV0Y2godXJsLCBmZXRjaE9wdGlvbnMpXG4gICAgICAgICAgICAudGhlbihyZXNwb25zZSA9PiB7XG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XG4gICAgICAgICAgICAgICAgICAgIC8vIGNoZWNrIGNvbnRlbnQgdHlwZVxuICAgICAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2UuaGVhZGVycy5nZXQoXCJDb250ZW50LVR5cGVcIikgJiYgcmVzcG9uc2UuaGVhZGVycy5nZXQoXCJDb250ZW50LVR5cGVcIikuaW5kZXhPZihcImFwcGxpY2F0aW9uL2pzb25cIikgIT09IC0xKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xuICAgICAgICAgICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLnRleHQoKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICByZWplY3QoRXJyb3IocmVzcG9uc2Uuc3RhdHVzVGV4dCkpO1xuICAgICAgICAgICAgfSlcbiAgICAgICAgICAgIC50aGVuKGRhdGEgPT4gcmVzb2x2ZShkYXRhKSlcbiAgICAgICAgICAgIC5jYXRjaChlcnJvciA9PiByZWplY3QoZXJyb3IpKTtcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0LCB3aW5kb3dOYW1lKSB7XG4gICAgcmV0dXJuIGZ1bmN0aW9uIChtZXRob2QsIGFyZ3M9bnVsbCkge1xuICAgICAgICByZXR1cm4gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0LCBtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpO1xuICAgIH07XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQXBwbGljYXRpb24pO1xuXG5sZXQgbWV0aG9kcyA9IHtcbiAgICBIaWRlOiAwLFxuICAgIFNob3c6IDEsXG4gICAgUXVpdDogMixcbn1cblxuLyoqXG4gKiBIaWRlIHRoZSBhcHBsaWNhdGlvblxuICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBIaWRlKCkge1xuICAgIHJldHVybiBjYWxsKG1ldGhvZHMuSGlkZSk7XG59XG5cbi8qKlxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cbiAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICByZXR1cm4gY2FsbChtZXRob2RzLlNob3cpO1xufVxuXG4vKipcbiAqIFF1aXQgdGhlIGFwcGxpY2F0aW9uXG4gKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgcmV0dXJuIGNhbGwobWV0aG9kcy5RdWl0KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5TY3JlZW59IFNjcmVlblxuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLlNjcmVlbnMpO1xuXG5sZXQgU2NyZWVuc0dldEFsbCA9IDA7XG5sZXQgU2NyZWVuc0dldFByaW1hcnkgPSAxO1xubGV0IFNjcmVlbnNHZXRDdXJyZW50ID0gMjtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCkge1xuICAgIHJldHVybiBjYWxsKFNjcmVlbnNHZXRBbGwpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoU2NyZWVuc0dldFByaW1hcnkpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKFNjcmVlbnNHZXRDdXJyZW50KTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TeXN0ZW0pO1xuXG5sZXQgU3lzdGVtSXNEYXJrTW9kZSA9IDA7XG5cbi8qKlxuICogRGV0ZXJtaW5lcyBpZiB0aGUgc3lzdGVtIGlzIGN1cnJlbnRseSB1c2luZyBkYXJrIG1vZGVcbiAqIEByZXR1cm5zIHtQcm9taXNlPGJvb2xlYW4+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gSXNEYXJrTW9kZSgpIHtcbiAgICByZXR1cm4gY2FsbChTeXN0ZW1Jc0RhcmtNb2RlKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQnJvd3Nlcik7XG5cbmxldCBCcm93c2VyT3BlblVSTCA9IDA7XG5cbi8qKlxuICogT3BlbiBhIGJyb3dzZXIgd2luZG93IHRvIHRoZSBnaXZlbiBVUkxcbiAqIEBwYXJhbSB7c3RyaW5nfSB1cmwgLSBUaGUgVVJMIHRvIG9wZW5cbiAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlblVSTCh1cmwpIHtcbiAgICByZXR1cm4gY2FsbChCcm93c2VyT3BlblVSTCwge3VybH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2FsbCk7XG5cbmxldCBDYWxsQmluZGluZyA9IDA7XG5cbmxldCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gY2FsbENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xuXG4gICAgICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBDYWxsQnlOYW1lKG5hbWUsIC4uLmFyZ3MpIHtcblxuICAgIC8vIEVuc3VyZSBmaXJzdCBhcmd1bWVudCBpcyBhIHN0cmluZyBhbmQgaGFzIDIgZG90c1xuICAgIGlmICh0eXBlb2YgbmFtZSAhPT0gXCJzdHJpbmdcIiB8fCBuYW1lLnNwbGl0KFwiLlwiKS5sZW5ndGggIT09IDMpIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiQ2FsbEJ5TmFtZSByZXF1aXJlcyBhIHN0cmluZyBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnXCIpO1xuICAgIH1cbiAgICAvLyBTcGxpdCBpbnB1dHNcbiAgICBsZXQgcGFydHMgPSBuYW1lLnNwbGl0KFwiLlwiKTtcblxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xuICAgICAgICBwYWNrYWdlTmFtZTogcGFydHNbMF0sXG4gICAgICAgIHN0cnVjdE5hbWU6IHBhcnRzWzFdLFxuICAgICAgICBtZXRob2ROYW1lOiBwYXJ0c1syXSxcbiAgICAgICAgYXJnczogYXJncyxcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGxCeUlEKG1ldGhvZElELCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XG4gICAgICAgIG1ldGhvZElEOiBtZXRob2RJRCxcbiAgICAgICAgYXJnczogYXJncyxcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBDYWxsIGEgcGx1Z2luIG1ldGhvZFxuICogQHBhcmFtIHtzdHJpbmd9IHBsdWdpbk5hbWUgLSBuYW1lIG9mIHRoZSBwbHVnaW5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gbmFtZSBvZiB0aGUgbWV0aG9kXG4gKiBAcGFyYW0gey4uLmFueX0gYXJncyAtIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2RcbiAqIEByZXR1cm5zIHtQcm9taXNlPGFueT59IC0gcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHJlc3VsdFxuICovXG5leHBvcnQgZnVuY3Rpb24gUGx1Z2luKHBsdWdpbk5hbWUsIG1ldGhvZE5hbWUsIC4uLmFyZ3MpIHtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIHtcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxuICAgICAgICBzdHJ1Y3ROYW1lOiBwbHVnaW5OYW1lLFxuICAgICAgICBtZXRob2ROYW1lOiBtZXRob2ROYW1lLFxuICAgICAgICBhcmdzOiBhcmdzLFxuICAgIH0pO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2l6ZX0gU2l6ZVxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5Qb3NpdGlvbn0gUG9zaXRpb25cbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2NyZWVufSBTY3JlZW5cbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBXaW5kb3dDZW50ZXIgPSAwO1xubGV0IFdpbmRvd1NldFRpdGxlID0gMTtcbmxldCBXaW5kb3dGdWxsc2NyZWVuID0gMjtcbmxldCBXaW5kb3dVbkZ1bGxzY3JlZW4gPSAzO1xubGV0IFdpbmRvd1NldFNpemUgPSA0O1xubGV0IFdpbmRvd1NpemUgPSA1O1xubGV0IFdpbmRvd1NldE1heFNpemUgPSA2O1xubGV0IFdpbmRvd1NldE1pblNpemUgPSA3O1xubGV0IFdpbmRvd1NldEFsd2F5c09uVG9wID0gODtcbmxldCBXaW5kb3dTZXRSZWxhdGl2ZVBvc2l0aW9uID0gOTtcbmxldCBXaW5kb3dSZWxhdGl2ZVBvc2l0aW9uID0gMTA7XG5sZXQgV2luZG93U2NyZWVuID0gMTE7XG5sZXQgV2luZG93SGlkZSA9IDEyO1xubGV0IFdpbmRvd01heGltaXNlID0gMTM7XG5sZXQgV2luZG93VW5NYXhpbWlzZSA9IDE0O1xubGV0IFdpbmRvd1RvZ2dsZU1heGltaXNlID0gMTU7XG5sZXQgV2luZG93TWluaW1pc2UgPSAxNjtcbmxldCBXaW5kb3dVbk1pbmltaXNlID0gMTc7XG5sZXQgV2luZG93UmVzdG9yZSA9IDE4O1xubGV0IFdpbmRvd1Nob3cgPSAxOTtcbmxldCBXaW5kb3dDbG9zZSA9IDIwO1xubGV0IFdpbmRvd1NldEJhY2tncm91bmRDb2xvdXIgPSAyMTtcbmxldCBXaW5kb3dTZXRSZXNpemFibGUgPSAyMjtcbmxldCBXaW5kb3dXaWR0aCA9IDIzO1xubGV0IFdpbmRvd0hlaWdodCA9IDI0O1xubGV0IFdpbmRvd1pvb21JbiA9IDI1O1xubGV0IFdpbmRvd1pvb21PdXQgPSAyNjtcbmxldCBXaW5kb3dab29tUmVzZXQgPSAyNztcbmxldCBXaW5kb3dHZXRab29tTGV2ZWwgPSAyODtcbmxldCBXaW5kb3dTZXRab29tTGV2ZWwgPSAyOTtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1dpbmRvdyh3aW5kb3dOYW1lKSB7XG4gICAgbGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgd2luZG93TmFtZSk7XG4gICAgcmV0dXJuIHtcblxuICAgICAgICAvKipcbiAgICAgICAgICogQ2VudGVycyB0aGUgd2luZG93LlxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIENlbnRlcjogKCkgPT4gY2FsbChXaW5kb3dDZW50ZXIpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyB0aXRsZS5cbiAgICAgICAgICogQHBhcmFtIHRpdGxlXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gY2FsbChXaW5kb3dTZXRUaXRsZSwge3RpdGxlfSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIE1ha2VzIHRoZSB3aW5kb3cgZnVsbHNjcmVlbi5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiBjYWxsKFdpbmRvd0Z1bGxzY3JlZW4pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbmZ1bGxzY3JlZW4gdGhlIHdpbmRvdy5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IGNhbGwoV2luZG93VW5GdWxsc2NyZWVuKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoIFRoZSB3aW5kb3cgd2lkdGhcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCBUaGUgd2luZG93IGhlaWdodFxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKFdpbmRvd1NldFNpemUsIHt3aWR0aCxoZWlnaHR9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgc2l6ZS5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8U2l6ZT59IFRoZSB3aW5kb3cgc2l6ZVxuICAgICAgICAgKi9cbiAgICAgICAgU2l6ZTogKCkgPT4gY2FsbChXaW5kb3dTaXplKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgbWF4aW11bSBzaXplLlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKFdpbmRvd1NldE1heFNpemUsIHt3aWR0aCxoZWlnaHR9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgbWluaW11bSBzaXplLlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFNldE1pblNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKFdpbmRvd1NldE1pblNpemUsIHt3aWR0aCxoZWlnaHR9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHdpbmRvdyB0byBiZSBhbHdheXMgb24gdG9wLlxuICAgICAgICAgKiBAcGFyYW0ge2Jvb2xlYW59IG9uVG9wIFdoZXRoZXIgdGhlIHdpbmRvdyBzaG91bGQgYmUgYWx3YXlzIG9uIHRvcFxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAob25Ub3ApID0+IGNhbGwoV2luZG93U2V0QWx3YXlzT25Ub3AsIHthbHdheXNPblRvcDpvblRvcH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyByZWxhdGl2ZSBwb3NpdGlvbi5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHhcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHlcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbChXaW5kb3dTZXRSZWxhdGl2ZVBvc2l0aW9uLCB7eCx5fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IHBvc2l0aW9uLlxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxQb3NpdGlvbj59IFRoZSB3aW5kb3cgcG9zaXRpb25cbiAgICAgICAgICovXG4gICAgICAgIFJlbGF0aXZlUG9zaXRpb246ICgpID0+IGNhbGwoV2luZG93UmVsYXRpdmVQb3NpdGlvbiksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEdldCB0aGUgc2NyZWVuIHRoZSB3aW5kb3cgaXMgb24uXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gICAgICAgICAqL1xuICAgICAgICBTY3JlZW46ICgpID0+IGNhbGwoV2luZG93U2NyZWVuKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogSGlkZSB0aGUgd2luZG93XG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICAgICAgICAgKi9cbiAgICAgICAgSGlkZTogKCkgPT4gY2FsbChXaW5kb3dIaWRlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWF4aW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIE1heGltaXNlOiAoKSA9PiBjYWxsKFdpbmRvd01heGltaXNlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2hvdyB0aGUgd2luZG93XG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICAgICAgICAgKi9cbiAgICAgICAgU2hvdzogKCkgPT4gY2FsbChXaW5kb3dTaG93KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogQ2xvc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIENsb3NlOiAoKSA9PiBjYWxsKFdpbmRvd0Nsb3NlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVG9nZ2xlIHRoZSB3aW5kb3cgbWF4aW1pc2Ugc3RhdGVcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBUb2dnbGVNYXhpbWlzZTogKCkgPT4gY2FsbChXaW5kb3dUb2dnbGVNYXhpbWlzZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFVubWF4aW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IGNhbGwoV2luZG93VW5NYXhpbWlzZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIE1pbmltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBNaW5pbWlzZTogKCkgPT4gY2FsbChXaW5kb3dNaW5pbWlzZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFVubWluaW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFVuTWluaW1pc2U6ICgpID0+IGNhbGwoV2luZG93VW5NaW5pbWlzZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFJlc3RvcmUgdGhlIHdpbmRvd1xuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFJlc3RvcmU6ICgpID0+IGNhbGwoV2luZG93UmVzdG9yZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBnIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gYiAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGEgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IGNhbGwoV2luZG93U2V0QmFja2dyb3VuZENvbG91ciwge3IsIGcsIGIsIGF9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHdoZXRoZXIgdGhlIHdpbmRvdyBjYW4gYmUgcmVzaXplZCBvciBub3RcbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSByZXNpemFibGVcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBTZXRSZXNpemFibGU6IChyZXNpemFibGUpID0+IGNhbGwoV2luZG93U2V0UmVzaXphYmxlLCB7cmVzaXphYmxlfSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IHdpZHRoXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPG51bWJlcj59XG4gICAgICAgICAqL1xuICAgICAgICBXaWR0aDogKCkgPT4gY2FsbChXaW5kb3dXaWR0aCksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IGhlaWdodFxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxudW1iZXI+fVxuICAgICAgICAgKi9cbiAgICAgICAgSGVpZ2h0OiAoKSA9PiBjYWxsKFdpbmRvd0hlaWdodCksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFpvb20gaW4gdGhlIHdpbmRvd1xuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTx2b2lkPn1cbiAgICAgICAgICovXG4gICAgICAgIFpvb21JbjogKCkgPT4gY2FsbChXaW5kb3dab29tSW4pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBab29tIG91dCB0aGUgd2luZG93XG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICAgICAgICAgKi9cbiAgICAgICAgWm9vbU91dDogKCkgPT4gY2FsbChXaW5kb3dab29tT3V0KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogUmVzZXQgdGhlIHdpbmRvdyB6b29tXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICAgICAgICAgKi9cbiAgICAgICAgWm9vbVJlc2V0OiAoKSA9PiBjYWxsKFdpbmRvd1pvb21SZXNldCksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIEdldCB0aGUgd2luZG93IHpvb21cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8bnVtYmVyPn1cbiAgICAgICAgICovXG4gICAgICAgIEdldFpvb21MZXZlbDogKCkgPT4gY2FsbChXaW5kb3dHZXRab29tTGV2ZWwpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyB6b29tIGxldmVsXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB6b29tTGV2ZWxcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8dm9pZD59XG4gICAgICAgICAqL1xuICAgICAgICBTZXRab29tTGV2ZWw6ICh6b29tTGV2ZWwpID0+IGNhbGwoV2luZG93U2V0Wm9vbUxldmVsLCB7em9vbUxldmVsfSksXG4gICAgfTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5XYWlsc0V2ZW50fSBXYWlsc0V2ZW50XG4gKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuRXZlbnRzKTtcbmxldCBFdmVudEVtaXQgPSAwO1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHsoZGF0YTogYW55KSA9PiB2b2lkfSBjYWxsYmFja1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxuICAgICAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5cbi8qKlxuICogV2FpbHNFdmVudCBkZWZpbmVzIGEgY3VzdG9tIGV2ZW50LiBJdCBpcyBwYXNzZWQgdG8gZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBjbGFzcyBXYWlsc0V2ZW50XG4gKiBAcHJvcGVydHkge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gKiBAcHJvcGVydHkge2FueX0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxuICovXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBXYWlsc0V2ZW50LlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gTmFtZSBvZiB0aGUgZXZlbnRcbiAgICAgKiBAcGFyYW0ge2FueT1udWxsfSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XG4gICAgICogQG1lbWJlcm9mIFdhaWxzRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhID0gbnVsbCkge1xuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHsoZGF0YTogYW55KSA9PiB2b2lkfSBjYWxsYmFja1xuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xuICogQHJldHVybnMgeygpID0+IHZvaWR9IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHsoKSA9PiB2b2lkfSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHJldHVybnMgeygpID0+IHZvaWR9IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKiBAcmV0dXJucyB7dm9pZH1cbiAqL1xuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBkaXNwYXRjaGVzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmIChsaXN0ZW5lcnMpIHtcbiAgICAgICAgLy8gaXRlcmF0ZSBsaXN0ZW5lcnMgYW5kIGNhbGwgY2FsbGJhY2suIElmIGNhbGxiYWNrIHJldHVybnMgdHJ1ZSwgcmVtb3ZlIGxpc3RlbmVyXG4gICAgICAgIGxldCB0b1JlbW92ZSA9IFtdO1xuICAgICAgICBsaXN0ZW5lcnMuZm9yRWFjaChsaXN0ZW5lciA9PiB7XG4gICAgICAgICAgICBsZXQgcmVtb3ZlID0gbGlzdGVuZXIuQ2FsbGJhY2soZXZlbnQpO1xuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xuICAgICAgICBpZiAodG9SZW1vdmUubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuIFt2MyBDSEFOR0VdIE9mZiBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICogQHJldHVybnMge3ZvaWR9XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xuICAgIGxldCBldmVudHNUb1JlbW92ZSA9IFtldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzXTtcbiAgICBldmVudHNUb1JlbW92ZS5mb3JFYWNoKGV2ZW50TmFtZSA9PiB7XG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xuICAgIH0pO1xufVxuXG4vKipcbiAqIE9mZkFsbCB1bnJlZ2lzdGVycyBhbGwgbGlzdGVuZXJzXG4gKiBbdjMgQ0hBTkdFXSBPZmZBbGwgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqIEByZXR1cm5zIHt2b2lkfVxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkge1xuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XG59XG5cbi8qKlxuICogRW1pdCBhbiBldmVudFxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCBUaGUgZXZlbnQgdG8gZW1pdFxuICpcbiAqIEByZXR1cm5zIHtQcm9taXNlPHZvaWQ+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkge1xuICAgIHJldHVybiBjYWxsKEV2ZW50RW1pdCwgZXZlbnQpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk1lc3NhZ2VEaWFsb2dPcHRpb25zfSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk9wZW5EaWFsb2dPcHRpb25zfSBPcGVuRGlhbG9nT3B0aW9uc1xuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNhdmVEaWFsb2dPcHRpb25zfSBTYXZlRGlhbG9nT3B0aW9uc1xuICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuRGlhbG9nKTtcblxubGV0IERpYWxvZ0luZm8gPSAwO1xubGV0IERpYWxvZ1dhcm5pbmcgPSAxO1xubGV0IERpYWxvZ0Vycm9yID0gMjtcbmxldCBEaWFsb2dRdWVzdGlvbiA9IDM7XG5sZXQgRGlhbG9nT3BlbkZpbGUgPSA0O1xubGV0IERpYWxvZ1NhdmVGaWxlID0gNTtcblxuXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dDYWxsYmFjayhpZCwgZGF0YSwgaXNKU09OKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zKSB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xuICAgICAgICBkaWFsb2dSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgICAgIH0pO1xuICAgIH0pO1xufVxuXG5cbi8qKlxuICogU2hvd3MgYW4gSW5mbyBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coRGlhbG9nSW5mbywgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYSBXYXJuaW5nIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhEaWFsb2dXYXJuaW5nLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhbiBFcnJvciBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKERpYWxvZ0Vycm9yLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhIFF1ZXN0aW9uIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVlc3Rpb24ob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coRGlhbG9nUXVlc3Rpb24sIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGFuIE9wZW4gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlIG9yIGFuIGFycmF5IG9mIHNlbGVjdGVkIGZpbGVzIGlmIEFsbG93c011bHRpcGxlU2VsZWN0aW9uIGlzIHRydWUuIEEgYmxhbmsgc3RyaW5nIGlzIHJldHVybmVkIGlmIG5vIGZpbGUgd2FzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coRGlhbG9nT3BlbkZpbGUsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGEgU2F2ZSBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7U2F2ZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUuIEEgYmxhbmsgc3RyaW5nIGlzIHJldHVybmVkIGlmIG5vIGZpbGUgd2FzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coRGlhbG9nU2F2ZUZpbGUsIG9wdGlvbnMpO1xufVxuXG4iLCAiaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ29udGV4dE1lbnUpO1xuXG5sZXQgQ29udGV4dE1lbnVPcGVuID0gMDtcblxuZnVuY3Rpb24gb3BlbkNvbnRleHRNZW51KGlkLCB4LCB5LCBkYXRhKSB7XG4gICAgdm9pZCBjYWxsKENvbnRleHRNZW51T3Blbiwge2lkLCB4LCB5LCBkYXRhfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzZXR1cENvbnRleHRNZW51cygpIHtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xufVxuXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxuICAgIGxldCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xuICAgIGxldCBjdXN0b21Db250ZXh0TWVudSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudVwiKTtcbiAgICBjdXN0b21Db250ZXh0TWVudSA9IGN1c3RvbUNvbnRleHRNZW51ID8gY3VzdG9tQ29udGV4dE1lbnUudHJpbSgpIDogXCJcIjtcbiAgICBpZiAoY3VzdG9tQ29udGV4dE1lbnUpIHtcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgbGV0IGN1c3RvbUNvbnRleHRNZW51RGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xuICAgICAgICBvcGVuQ29udGV4dE1lbnUoY3VzdG9tQ29udGV4dE1lbnUsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGN1c3RvbUNvbnRleHRNZW51RGF0YSk7XG4gICAgICAgIHJldHVyblxuICAgIH1cblxuICAgIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpO1xufVxuXG5cbi8qXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGF1dG87IChkZWZhdWx0KSB3aWxsIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51IGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgZWxlbWVudCBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93OyB3aWxsIGFsd2F5cyBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBoaWRlOyB3aWxsIGFsd2F5cyBoaWRlIHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuXG5UaGlzIHJ1bGUgaXMgaW5oZXJpdGVkIGxpa2Ugbm9ybWFsIENTUyBydWxlcywgc28gbmVzdGluZyB3b3JrcyBhcyBleHBlY3RlZFxuKi9cbmZ1bmN0aW9uIHByb2Nlc3NEZWZhdWx0Q29udGV4dE1lbnUoZXZlbnQpIHtcbiAgICAvLyBEZWJ1ZyBidWlsZHMgYWx3YXlzIHNob3cgdGhlIG1lbnVcbiAgICBpZiAoREVCVUcpIHtcbiAgICAgICAgcmV0dXJuO1xuICAgIH1cblxuICAgIC8vIFByb2Nlc3MgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICBjb25zdCBlbGVtZW50ID0gZXZlbnQudGFyZ2V0O1xuICAgIGNvbnN0IGNvbXB1dGVkU3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcbiAgICBjb25zdCBkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24gPSBjb21wdXRlZFN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpO1xuICAgIHN3aXRjaCAoZGVmYXVsdENvbnRleHRNZW51QWN0aW9uKSB7XG4gICAgICAgIGNhc2UgXCJzaG93XCI6XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGNhc2UgXCJoaWRlXCI6XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBkZWZhdWx0OlxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWVcbiAgICAgICAgICAgIGlmIChlbGVtZW50LmlzQ29udGVudEVkaXRhYmxlKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkXG4gICAgICAgICAgICBjb25zdCBzZWxlY3Rpb24gPSB3aW5kb3cuZ2V0U2VsZWN0aW9uKCk7XG4gICAgICAgICAgICBjb25zdCBoYXNTZWxlY3Rpb24gPSAoc2VsZWN0aW9uLnRvU3RyaW5nKCkubGVuZ3RoID4gMClcbiAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24pIHtcbiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcbiAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xuICAgICAgICAgICAgICAgICAgICBmb3IgKGxldCBqID0gMDsgaiA8IHJlY3RzLmxlbmd0aDsgaisrKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0ID0gcmVjdHNbal07XG4gICAgICAgICAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZWxlbWVudEZyb21Qb2ludChyZWN0LmxlZnQsIHJlY3QudG9wKSA9PT0gZWxlbWVudCkge1xuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRhZ25hbWUgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbiAgICAgICAgICAgIGlmIChlbGVtZW50LnRhZ05hbWUgPT09IFwiSU5QVVRcIiB8fCBlbGVtZW50LnRhZ05hbWUgPT09IFwiVEVYVEFSRUFcIikge1xuICAgICAgICAgICAgICAgIGlmIChoYXNTZWxlY3Rpb24gfHwgKCFlbGVtZW50LnJlYWRPbmx5ICYmICFlbGVtZW50LmRpc2FibGVkKSkge1xuICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuXG4gICAgICAgICAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgIH1cbn1cbiIsICJcbmltcG9ydCB7RW1pdCwgV2FpbHNFdmVudH0gZnJvbSBcIi4vZXZlbnRzXCI7XG5pbXBvcnQge1F1ZXN0aW9ufSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5cbmZ1bmN0aW9uIHNlbmRFdmVudChldmVudE5hbWUsIGRhdGE9bnVsbCkge1xuICAgIGxldCBldmVudCA9IG5ldyBXYWlsc0V2ZW50KGV2ZW50TmFtZSwgZGF0YSk7XG4gICAgRW1pdChldmVudCk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCkge1xuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1ldmVudF0nKTtcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XG4gICAgICAgIGNvbnN0IGV2ZW50VHlwZSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBEZXRhY2hlZDogZmFsc2UsIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcbiAgICAgICAgfTtcbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xuXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKG1ldGhvZCkge1xuICAgIGlmICh3YWlscy5XaW5kb3dbbWV0aG9kXSA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIGNvbnNvbGUubG9nKFwiV2luZG93IG1ldGhvZCBcIiArIG1ldGhvZCArIFwiIG5vdCBmb3VuZFwiKTtcbiAgICB9XG4gICAgd2FpbHMuV2luZG93W21ldGhvZF0oKTtcbn1cblxuZnVuY3Rpb24gYWRkV01MV2luZG93TGlzdGVuZXJzKCkge1xuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC13aW5kb3ddJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXdpbmRvdycpO1xuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTE9wZW5Ccm93c2VyTGlzdGVuZXIoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbd21sLW9wZW51cmxdJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCB1cmwgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLW9wZW51cmwnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHZvaWQgd2FpbHMuQnJvd3Nlci5PcGVuVVJMKHVybCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9KTtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICB2b2lkIHdhaWxzLkJyb3dzZXIuT3BlblVSTCh1cmwpO1xuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiByZWxvYWRXTUwoKSB7XG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcbiAgICBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKTtcbiAgICBhZGRXTUxPcGVuQnJvd3Nlckxpc3RlbmVyKCk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLy8gZGVmaW5lZCBpbiB0aGUgVGFza2ZpbGVcbmV4cG9ydCBsZXQgaW52b2tlID0gZnVuY3Rpb24oaW5wdXQpIHtcbiAgICBpZihXSU5ET1dTKSB7XG4gICAgICAgIGNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlKGlucHV0KTtcbiAgICB9IGVsc2Uge1xuICAgICAgICB3ZWJraXQubWVzc2FnZUhhbmRsZXJzLmV4dGVybmFsLnBvc3RNZXNzYWdlKGlucHV0KTtcbiAgICB9XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxubGV0IGZsYWdzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBjb252ZXJ0VG9NYXAob2JqKSB7XG4gICAgY29uc3QgbWFwID0gbmV3IE1hcCgpO1xuXG4gICAgZm9yIChjb25zdCBba2V5LCB2YWx1ZV0gb2YgT2JqZWN0LmVudHJpZXMob2JqKSkge1xuICAgICAgICBpZiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyAmJiB2YWx1ZSAhPT0gbnVsbCkge1xuICAgICAgICAgICAgbWFwLnNldChrZXksIGNvbnZlcnRUb01hcCh2YWx1ZSkpOyAvLyBSZWN1cnNpdmVseSBjb252ZXJ0IG5lc3RlZCBvYmplY3RcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIG1hcC5zZXQoa2V5LCB2YWx1ZSk7XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICByZXR1cm4gbWFwO1xufVxuXG5mZXRjaChcIi93YWlscy9mbGFnc1wiKS50aGVuKChyZXNwb25zZSkgPT4ge1xuICAgIHJlc3BvbnNlLmpzb24oKS50aGVuKChkYXRhKSA9PiB7XG4gICAgICAgIGZsYWdzID0gY29udmVydFRvTWFwKGRhdGEpO1xuICAgIH0pO1xufSk7XG5cblxuZnVuY3Rpb24gZ2V0VmFsdWVGcm9tTWFwKGtleVN0cmluZykge1xuICAgIGNvbnN0IGtleXMgPSBrZXlTdHJpbmcuc3BsaXQoJy4nKTtcbiAgICBsZXQgdmFsdWUgPSBmbGFncztcblxuICAgIGZvciAoY29uc3Qga2V5IG9mIGtleXMpIHtcbiAgICAgICAgaWYgKHZhbHVlIGluc3RhbmNlb2YgTWFwKSB7XG4gICAgICAgICAgICB2YWx1ZSA9IHZhbHVlLmdldChrZXkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgdmFsdWUgPSB2YWx1ZVtrZXldO1xuICAgICAgICB9XG5cbiAgICAgICAgaWYgKHZhbHVlID09PSB1bmRlZmluZWQpIHtcbiAgICAgICAgICAgIGJyZWFrO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIHZhbHVlO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXlTdHJpbmcpIHtcbiAgICByZXR1cm4gZ2V0VmFsdWVGcm9tTWFwKGtleVN0cmluZyk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtpbnZva2V9IGZyb20gXCIuL2ludm9rZVwiO1xuaW1wb3J0IHtHZXRGbGFnfSBmcm9tIFwiLi9mbGFnc1wiO1xuXG5sZXQgc2hvdWxkRHJhZyA9IGZhbHNlO1xuXG5leHBvcnQgZnVuY3Rpb24gZHJhZ1Rlc3QoZSkge1xuICAgIGxldCB2YWwgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlLnRhcmdldCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0td2Via2l0LWFwcC1yZWdpb25cIik7XG4gICAgaWYgKHZhbCkge1xuICAgICAgICB2YWwgPSB2YWwudHJpbSgpO1xuICAgIH1cblxuICAgIGlmICh2YWwgIT09IFwiZHJhZ1wiKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICAvLyBPbmx5IHByb2Nlc3MgdGhlIHByaW1hcnkgYnV0dG9uXG4gICAgaWYgKGUuYnV0dG9ucyAhPT0gMSkge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgcmV0dXJuIGUuZGV0YWlsID09PSAxO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc2V0dXBEcmFnKCkge1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZWRvd24nLCBvbk1vdXNlRG93bik7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIG9uTW91c2VNb3ZlKTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIG9uTW91c2VVcCk7XG59XG5cbmxldCByZXNpemVFZGdlID0gbnVsbDtcbmxldCByZXNpemFibGUgPSBmYWxzZTtcblxuZXhwb3J0IGZ1bmN0aW9uIHNldFJlc2l6YWJsZSh2YWx1ZSkge1xuICAgIHJlc2l6YWJsZSA9IHZhbHVlO1xufVxuXG5mdW5jdGlvbiB0ZXN0UmVzaXplKGUpIHtcbiAgICBpZiggcmVzaXplRWRnZSApIHtcbiAgICAgICAgaW52b2tlKFwicmVzaXplOlwiICsgcmVzaXplRWRnZSk7XG4gICAgICAgIHJldHVybiB0cnVlXG4gICAgfVxuICAgIHJldHVybiBmYWxzZTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xuXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nIG9uIFdpbmRvd3NcbiAgICBpZiggV0lORE9XUyApIHtcbiAgICAgICAgaWYgKHRlc3RSZXNpemUoKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmIChkcmFnVGVzdChlKSkge1xuICAgICAgICAvLyBJZ25vcmUgZHJhZyBvbiBzY3JvbGxiYXJzXG4gICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBzaG91bGREcmFnID0gdHJ1ZTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBvbk1vdXNlVXAoZSkge1xuICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcbiAgICAgICAgZW5kRHJhZygpO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGVuZERyYWcoKSB7XG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCBkZWZhdWx0Q3Vyc29yO1xuICAgIHJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGUpIHtcbiAgICBpZiAoc2hvdWxkRHJhZykge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgICAgICBpbnZva2UoXCJkcmFnXCIpO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoV0lORE9XUykge1xuICAgICAgICBpZiAocmVzaXphYmxlKSB7XG4gICAgICAgICAgICBoYW5kbGVSZXNpemUoZSk7XG4gICAgICAgIH1cbiAgICB9XG59XG5cbmxldCBkZWZhdWx0Q3Vyc29yID0gXCJhdXRvXCI7XG5cbmZ1bmN0aW9uIGhhbmRsZVJlc2l6ZShlKSB7XG4gICAgbGV0IHJlc2l6ZUhhbmRsZUhlaWdodCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlSGVpZ2h0XCIpIHx8IDU7XG4gICAgbGV0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xuXG4gICAgLy8gRXh0cmEgcGl4ZWxzIGZvciB0aGUgY29ybmVyIGFyZWFzXG4gICAgbGV0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDEwO1xuXG4gICAgbGV0IHJpZ2h0Qm9yZGVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcblxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyc1xuICAgIGxldCByaWdodENvcm5lciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGxldCBsZWZ0Q29ybmVyID0gZS5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGxldCB0b3BDb3JuZXIgPSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuICAgIGxldCBib3R0b21Db3JuZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiByZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxuICAgIC8vIEFkanVzdGVkIGZvciBjb3JuZXIgYXJlYXNcbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcbn1cbiIsICIvKlxuIF8gICAgIF9fICAgICBfIF9fXG58IHwgIC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSAnLi9jbGlwYm9hcmQnO1xuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSAnLi9hcHBsaWNhdGlvbic7XG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gJy4vc2NyZWVucyc7XG5pbXBvcnQgKiBhcyBTeXN0ZW0gZnJvbSAnLi9zeXN0ZW0nO1xuaW1wb3J0ICogYXMgQnJvd3NlciBmcm9tICcuL2Jyb3dzZXInO1xuaW1wb3J0IHtQbHVnaW4sIENhbGwsIGNhbGxFcnJvckNhbGxiYWNrLCBjYWxsQ2FsbGJhY2ssIENhbGxCeUlELCBDYWxsQnlOYW1lfSBmcm9tIFwiLi9jYWxsc1wiO1xuaW1wb3J0IHtjbGllbnRJZH0gZnJvbSAnLi9ydW50aW1lJztcbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcbmltcG9ydCB7ZGlzcGF0Y2hXYWlsc0V2ZW50LCBFbWl0LCBPZmYsIE9mZkFsbCwgT24sIE9uY2UsIE9uTXVsdGlwbGUsIFdhaWxzRXZlbnR9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtkaWFsb2dDYWxsYmFjaywgZGlhbG9nRXJyb3JDYWxsYmFjaywgRXJyb3IsIEluZm8sIE9wZW5GaWxlLCBRdWVzdGlvbiwgU2F2ZUZpbGUsIFdhcm5pbmcsfSBmcm9tIFwiLi9kaWFsb2dzXCI7XG5pbXBvcnQge3NldHVwQ29udGV4dE1lbnVzfSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xuaW1wb3J0IHtyZWxvYWRXTUx9IGZyb20gXCIuL3dtbFwiO1xuaW1wb3J0IHtzZXR1cERyYWcsIGVuZERyYWcsIHNldFJlc2l6YWJsZX0gZnJvbSBcIi4vZHJhZ1wiO1xuXG53aW5kb3cud2FpbHMgPSB7XG4gICAgLi4ubmV3UnVudGltZShudWxsKSxcbiAgICBDYXBhYmlsaXRpZXM6IHt9LFxuICAgIGNsaWVudElkOiBjbGllbnRJZCxcbn07XG5cbmZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKS50aGVuKChyZXNwb25zZSkgPT4ge1xuICAgIHJlc3BvbnNlLmpzb24oKS50aGVuKChkYXRhKSA9PiB7XG4gICAgICAgIHdpbmRvdy53YWlscy5DYXBhYmlsaXRpZXMgPSBkYXRhO1xuICAgIH0pO1xufSk7XG5cbi8vIEludGVybmFsIHdhaWxzIGVuZHBvaW50c1xud2luZG93Ll93YWlscyA9IHtcbiAgICBkaWFsb2dDYWxsYmFjayxcbiAgICBkaWFsb2dFcnJvckNhbGxiYWNrLFxuICAgIGRpc3BhdGNoV2FpbHNFdmVudCxcbiAgICBjYWxsQ2FsbGJhY2ssXG4gICAgY2FsbEVycm9yQ2FsbGJhY2ssXG4gICAgZW5kRHJhZyxcbiAgICBzZXRSZXNpemFibGUsXG59O1xuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZSh3aW5kb3dOYW1lKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgQ2xpcGJvYXJkOiB7XG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcbiAgICAgICAgfSxcbiAgICAgICAgQXBwbGljYXRpb246IHtcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uLFxuICAgICAgICAgICAgR2V0V2luZG93QnlOYW1lKHdpbmRvd05hbWUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gbmV3UnVudGltZSh3aW5kb3dOYW1lKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSxcbiAgICAgICAgU3lzdGVtLFxuICAgICAgICBTY3JlZW5zLFxuICAgICAgICBCcm93c2VyLFxuICAgICAgICBDYWxsLFxuICAgICAgICBDYWxsQnlJRCxcbiAgICAgICAgQ2FsbEJ5TmFtZSxcbiAgICAgICAgUGx1Z2luLFxuICAgICAgICBXTUw6IHtcbiAgICAgICAgICAgIFJlbG9hZDogcmVsb2FkV01MLFxuICAgICAgICB9LFxuICAgICAgICBEaWFsb2c6IHtcbiAgICAgICAgICAgIEluZm8sXG4gICAgICAgICAgICBXYXJuaW5nLFxuICAgICAgICAgICAgRXJyb3IsXG4gICAgICAgICAgICBRdWVzdGlvbixcbiAgICAgICAgICAgIE9wZW5GaWxlLFxuICAgICAgICAgICAgU2F2ZUZpbGUsXG4gICAgICAgIH0sXG4gICAgICAgIEV2ZW50czoge1xuICAgICAgICAgICAgRW1pdCxcbiAgICAgICAgICAgIE9uLFxuICAgICAgICAgICAgT25jZSxcbiAgICAgICAgICAgIE9uTXVsdGlwbGUsXG4gICAgICAgICAgICBPZmYsXG4gICAgICAgICAgICBPZmZBbGwsXG4gICAgICAgICAgICBXYWlsc0V2ZW50LFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyh3aW5kb3dOYW1lKSxcbiAgICB9O1xufVxuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbnNldHVwQ29udGV4dE1lbnVzKCk7XG5zZXR1cERyYWcoKTtcblxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcihcIkRPTUNvbnRlbnRMb2FkZWRcIiwgZnVuY3Rpb24oKSB7XG4gICAgcmVsb2FkV01MKCk7XG59KTtcbiJdLAog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QUFVLFVBQVU7QUFBQSxFQUM5Qzs7O0FDcEhBLFNBQU8sUUFBUTtBQUFBLElBQ1gsR0FBRyxXQUFXLElBQUk7QUFBQSxJQUNsQixjQUFjLENBQUM7QUFBQSxJQUNmO0FBQUEsRUFDSjtBQUVBLFFBQU0scUJBQXFCLEVBQUUsS0FBSyxDQUFDLGFBQWE7QUFDNUMsYUFBUyxLQUFLLEVBQUUsS0FBSyxDQUFDLFNBQVM7QUFDM0IsYUFBTyxNQUFNLGVBQWU7QUFBQSxJQUNoQyxDQUFDO0FBQUEsRUFDTCxDQUFDO0FBR0QsU0FBTyxTQUFTO0FBQUEsSUFDWjtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLEVBQ0o7QUFFTyxXQUFTLFdBQVcsWUFBWTtBQUNuQyxXQUFPO0FBQUEsTUFDSCxXQUFXO0FBQUEsUUFDUCxHQUFHO0FBQUEsTUFDUDtBQUFBLE1BQ0EsYUFBYTtBQUFBLFFBQ1QsR0FBRztBQUFBLFFBQ0gsZ0JBQWdCQyxhQUFZO0FBQ3hCLGlCQUFPLFdBQVdBLFdBQVU7QUFBQSxRQUNoQztBQUFBLE1BQ0o7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQSxLQUFLO0FBQUEsUUFDRCxRQUFRO0FBQUEsTUFDWjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQSxPQUFBQztBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUSxVQUFVLFVBQVU7QUFBQSxJQUNoQztBQUFBLEVBQ0o7QUFFQSxNQUFJLE1BQU87QUFDUCxZQUFRLElBQUksaUNBQWlDO0FBQUEsRUFDakQ7QUFFQSxvQkFBa0I7QUFDbEIsWUFBVTtBQUVWLFdBQVMsaUJBQWlCLG9CQUFvQixXQUFXO0FBQ3JELGNBQVU7QUFBQSxFQUNkLENBQUM7IiwKICAibmFtZXMiOiBbImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgInJlc2l6YWJsZSIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiRXJyb3IiLCAiY2FsbCIsICJ3aW5kb3dOYW1lIiwgIkVycm9yIl0KfQo=
