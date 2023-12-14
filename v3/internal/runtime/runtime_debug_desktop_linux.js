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
    void call(ClipboardSetText, { text });
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
    void call2(methods.Hide);
  }
  function Show() {
    void call2(methods.Show);
  }
  function Quit() {
    void call2(methods.Quit);
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
    void call5(BrowserOpenURL, { url });
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
       */
      Center: () => void call10(WindowCenter),
      /**
       * Set the window title.
       * @param title
       */
      SetTitle: (title) => void call10(WindowSetTitle, { title }),
      /**
       * Makes the window fullscreen.
       */
      Fullscreen: () => void call10(WindowFullscreen),
      /**
       * Unfullscreen the window.
       */
      UnFullscreen: () => void call10(WindowUnFullscreen),
      /**
       * Set the window size.
       * @param {number} width The window width
       * @param {number} height The window height
       */
      SetSize: (width, height) => call10(WindowSetSize, { width, height }),
      /**
       * Get the window size.
       * @returns {Promise<Size>} The window size
       */
      Size: () => {
        return call10(WindowSize);
      },
      /**
       * Set the window maximum size.
       * @param {number} width
       * @param {number} height
       */
      SetMaxSize: (width, height) => void call10(WindowSetMaxSize, { width, height }),
      /**
       * Set the window minimum size.
       * @param {number} width
       * @param {number} height
       */
      SetMinSize: (width, height) => void call10(WindowSetMinSize, { width, height }),
      /**
       * Set window to be always on top.
       * @param {boolean} onTop Whether the window should be always on top
       */
      SetAlwaysOnTop: (onTop) => void call10(WindowSetAlwaysOnTop, { alwaysOnTop: onTop }),
      /**
       * Set the window relative position.
       * @param {number} x
       * @param {number} y
       */
      SetRelativePosition: (x, y) => call10(WindowSetRelativePosition, { x, y }),
      /**
       * Get the window position.
       * @returns {Promise<Position>} The window position
       */
      RelativePosition: () => {
        return call10(WindowRelativePosition);
      },
      /**
       * Get the screen the window is on.
       * @returns {Promise<Screen>}
       */
      Screen: () => {
        return call10(WindowScreen);
      },
      /**
       * Hide the window
       */
      Hide: () => void call10(WindowHide),
      /**
       * Maximise the window
       */
      Maximise: () => void call10(WindowMaximise),
      /**
       * Show the window
       */
      Show: () => void call10(WindowShow),
      /**
       * Close the window
       */
      Close: () => void call10(WindowClose),
      /**
       * Toggle the window maximise state
       */
      ToggleMaximise: () => void call10(WindowToggleMaximise),
      /**
       * Unmaximise the window
       */
      UnMaximise: () => void call10(WindowUnMaximise),
      /**
       * Minimise the window
       */
      Minimise: () => void call10(WindowMinimise),
      /**
       * Unminimise the window
       */
      UnMinimise: () => void call10(WindowUnMinimise),
      /**
       * Restore the window
       */
      Restore: () => void call10(WindowRestore),
      /**
       * Set the background colour of the window.
       * @param {number} r - A value between 0 and 255
       * @param {number} g - A value between 0 and 255
       * @param {number} b - A value between 0 and 255
       * @param {number} a - A value between 0 and 255
       */
      SetBackgroundColour: (r, g, b, a) => void call10(WindowSetBackgroundColour, { r, g, b, a }),
      /**
       * Set whether the window can be resized or not
       * @param {boolean} resizable
       */
      SetResizable: (resizable2) => void call10(WindowSetResizable, { resizable: resizable2 }),
      /**
       * Get the window width
       * @returns {Promise<number>}
       */
      Width: () => {
        return call10(WindowWidth);
      },
      /**
       * Get the window height
       * @returns {Promise<number>}
       */
      Height: () => {
        return call10(WindowHeight);
      },
      /**
       * Zoom in the window
       */
      ZoomIn: () => void call10(WindowZoomIn),
      /**
       * Zoom out the window
       */
      ZoomOut: () => void call10(WindowZoomOut),
      /**
       * Reset the window zoom
       */
      ZoomReset: () => void call10(WindowZoomReset),
      /**
       * Get the window zoom
       * @returns {Promise<number>}
       */
      GetZoomLevel: () => {
        return call10(WindowGetZoomLevel);
      },
      /**
       * Set the window zoom level
       * @param {number} zoomLevel
       */
      SetZoomLevel: (zoomLevel) => void call10(WindowSetZoomLevel, { zoomLevel })
    };
  }

  // desktop/events.js
  var call7 = newRuntimeCallerWithID(objectNames.Events);
  var EventEmit = 0;
  var Listener = class {
    /**
     * Creates an instance of Listener.
     * @param {string} eventName
     * @param {function} callback
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
    void call7(EventEmit, event);
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
    if (false) {
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
  var resizable = false;
  function setResizable(value) {
    resizable = value;
  }
  function onMouseDown(e) {
    if (false) {
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
  function onMouseMove(e) {
    if (shouldDrag) {
      shouldDrag = false;
      let mousePressed = e.buttons !== void 0 ? e.buttons : e.which;
      if (mousePressed > 0) {
        invoke("drag");
      }
      return;
    }
    if (false) {
      if (resizable) {
        handleResize(e);
      }
    }
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
        OffAll
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAibm9kZV9tb2R1bGVzL25hbm9pZC9ub24tc2VjdXJlL2luZGV4LmpzIiwgImRlc2t0b3AvcnVudGltZS5qcyIsICJkZXNrdG9wL2FwcGxpY2F0aW9uLmpzIiwgImRlc2t0b3Avc2NyZWVucy5qcyIsICJkZXNrdG9wL3N5c3RlbS5qcyIsICJkZXNrdG9wL2Jyb3dzZXIuanMiLCAiZGVza3RvcC9jYWxscy5qcyIsICJkZXNrdG9wL3dpbmRvdy5qcyIsICJkZXNrdG9wL2V2ZW50cy5qcyIsICJkZXNrdG9wL2RpYWxvZ3MuanMiLCAiZGVza3RvcC9jb250ZXh0bWVudS5qcyIsICJkZXNrdG9wL3dtbC5qcyIsICJkZXNrdG9wL2ludm9rZS5qcyIsICJkZXNrdG9wL2ZsYWdzLmpzIiwgImRlc2t0b3AvZHJhZy5qcyIsICJkZXNrdG9wL21haW4uanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNsaXBib2FyZCk7XG5cbmxldCBDbGlwYm9hcmRTZXRUZXh0ID0gMDtcbmxldCBDbGlwYm9hcmRUZXh0ID0gMTtcblxuLyoqXG4gKiBTZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICB2b2lkIGNhbGwoQ2xpcGJvYXJkU2V0VGV4dCwge3RleHR9KTtcbn1cblxuLyoqXG4gKiBHZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gVGV4dCgpIHtcbiAgICByZXR1cm4gY2FsbChDbGlwYm9hcmRUZXh0KTtcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXG4gXyAgICAgX18gICAgIF8gX19cbnwgfCAgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XG5cbmNvbnN0IHJ1bnRpbWVVUkwgPSB3aW5kb3cubG9jYXRpb24ub3JpZ2luICsgXCIvd2FpbHMvcnVudGltZVwiO1xuLy8gT2JqZWN0IE5hbWVzXG5leHBvcnQgY29uc3Qgb2JqZWN0TmFtZXMgPSB7XG4gICAgQ2FsbDogMCxcbiAgICBDbGlwYm9hcmQ6IDEsXG4gICAgQXBwbGljYXRpb246IDIsXG4gICAgRXZlbnRzOiAzLFxuICAgIENvbnRleHRNZW51OiA0LFxuICAgIERpYWxvZzogNSxcbiAgICBXaW5kb3c6IDYsXG4gICAgU2NyZWVuczogNyxcbiAgICBTeXN0ZW06IDgsXG4gICAgQnJvd3NlcjogOSxcbn1cbmV4cG9ydCBsZXQgY2xpZW50SWQgPSBuYW5vaWQoKTtcblxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgaWYoIG1ldGhvZCApIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICB9XG4gICAgbGV0IGZldGNoT3B0aW9ucyA9IHtcbiAgICAgICAgaGVhZGVyczoge30sXG4gICAgfTtcbiAgICBpZiAod2luZG93TmFtZSkge1xuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xuICAgIH1cbiAgICBpZiAoYXJncykge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xuICAgIH1cbiAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtY2xpZW50LWlkXCJdID0gY2xpZW50SWQ7XG5cbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcbiAgICAgICAgICAgIC50aGVuKHJlc3BvbnNlID0+IHtcbiAgICAgICAgICAgICAgICBpZiAocmVzcG9uc2Uub2spIHtcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKSAmJiByZXNwb25zZS5oZWFkZXJzLmdldChcIkNvbnRlbnQtVHlwZVwiKS5pbmRleE9mKFwiYXBwbGljYXRpb24vanNvblwiKSAhPT0gLTEpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS5qc29uKCk7XG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UudGV4dCgpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIHJlamVjdChFcnJvcihyZXNwb25zZS5zdGF0dXNUZXh0KSk7XG4gICAgICAgICAgICB9KVxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxuICAgICAgICAgICAgLmNhdGNoKGVycm9yID0+IHJlamVjdChlcnJvcikpO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZUNhbGxlcihvYmplY3QsIHdpbmRvd05hbWUpIHtcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XG4gICAgICAgIHJldHVybiBydW50aW1lQ2FsbChvYmplY3QgKyBcIi5cIiArIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncyk7XG4gICAgfTtcbn1cblxuZnVuY3Rpb24gcnVudGltZUNhbGxXaXRoSUQob2JqZWN0SUQsIG1ldGhvZCwgd2luZG93TmFtZSwgYXJncykge1xuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xuICAgIHVybC5zZWFyY2hQYXJhbXMuYXBwZW5kKFwib2JqZWN0XCIsIG9iamVjdElEKTtcbiAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcIm1ldGhvZFwiLCBtZXRob2QpO1xuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XG4gICAgICAgIGhlYWRlcnM6IHt9LFxuICAgIH07XG4gICAgaWYgKHdpbmRvd05hbWUpIHtcbiAgICAgICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLXdpbmRvdy1uYW1lXCJdID0gd2luZG93TmFtZTtcbiAgICB9XG4gICAgaWYgKGFyZ3MpIHtcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcbiAgICB9XG4gICAgZmV0Y2hPcHRpb25zLmhlYWRlcnNbXCJ4LXdhaWxzLWNsaWVudC1pZFwiXSA9IGNsaWVudElkO1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsV2l0aElEKG9iamVjdCwgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkFwcGxpY2F0aW9uKTtcblxubGV0IG1ldGhvZHMgPSB7XG4gICAgSGlkZTogMCxcbiAgICBTaG93OiAxLFxuICAgIFF1aXQ6IDIsXG59XG5cbi8qKlxuICogSGlkZSB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgdm9pZCBjYWxsKG1ldGhvZHMuSGlkZSk7XG59XG5cbi8qKlxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNob3coKSB7XG4gICAgdm9pZCBjYWxsKG1ldGhvZHMuU2hvdyk7XG59XG5cblxuLyoqXG4gKiBRdWl0IHRoZSBhcHBsaWNhdGlvblxuICovXG5leHBvcnQgZnVuY3Rpb24gUXVpdCgpIHtcbiAgICB2b2lkIGNhbGwobWV0aG9kcy5RdWl0KTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuU2NyZWVufSBTY3JlZW5cbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5TY3JlZW5zKTtcblxubGV0IFNjcmVlbnNHZXRBbGwgPSAwO1xubGV0IFNjcmVlbnNHZXRQcmltYXJ5ID0gMTtcbmxldCBTY3JlZW5zR2V0Q3VycmVudCA9IDI7XG5cbi8qKlxuICogR2V0cyBhbGwgc2NyZWVucy5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEFsbCgpIHtcbiAgICByZXR1cm4gY2FsbChTY3JlZW5zR2V0QWxsKTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBwcmltYXJ5IHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRQcmltYXJ5KCkge1xuICAgIHJldHVybiBjYWxsKFNjcmVlbnNHZXRQcmltYXJ5KTtcbn1cblxuLyoqXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxuICogQGNvbnN0cnVjdG9yXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBHZXRDdXJyZW50KCkge1xuICAgIHJldHVybiBjYWxsKFNjcmVlbnNHZXRDdXJyZW50KTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyV2l0aElELCBvYmplY3ROYW1lc30gZnJvbSBcIi4vcnVudGltZVwiO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuU3lzdGVtKTtcblxubGV0IFN5c3RlbUlzRGFya01vZGUgPSAwO1xuXG4vKipcbiAqIERldGVybWluZXMgaWYgdGhlIHN5c3RlbSBpcyBjdXJyZW50bHkgdXNpbmcgZGFyayBtb2RlXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxib29sZWFuPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIElzRGFya01vZGUoKSB7XG4gICAgcmV0dXJuIGNhbGwoU3lzdGVtSXNEYXJrTW9kZSk7XG59IiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkJyb3dzZXIpO1xuXG5sZXQgQnJvd3Nlck9wZW5VUkwgPSAwO1xuXG4vKipcbiAqIE9wZW4gYSBicm93c2VyIHdpbmRvdyB0byB0aGUgZ2l2ZW4gVVJMXG4gKiBAcGFyYW0ge3N0cmluZ30gdXJsIC0gVGhlIFVSTCB0byBvcGVuXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPcGVuVVJMKHVybCkge1xuICAgIHZvaWQgY2FsbChCcm93c2VyT3BlblVSTCwge3VybH0pO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXJXaXRoSUQob2JqZWN0TmFtZXMuQ2FsbCk7XG5cbmxldCBDYWxsQmluZGluZyA9IDA7XG5cbmxldCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xuXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xuICAgIGxldCByZXN1bHQ7XG4gICAgZG8ge1xuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gY2FsbENhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBpZiAoaXNKU09OKSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBwLnJlc29sdmUoZGF0YSk7XG4gICAgICAgIH1cbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGNhbGxFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5cbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiY2FsbC1pZFwiXSA9IGlkO1xuXG4gICAgICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywgb3B0aW9ucyk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBDYWxsQnlOYW1lKG5hbWUsIC4uLmFyZ3MpIHtcblxuICAgIC8vIEVuc3VyZSBmaXJzdCBhcmd1bWVudCBpcyBhIHN0cmluZyBhbmQgaGFzIDIgZG90c1xuICAgIGlmICh0eXBlb2YgbmFtZSAhPT0gXCJzdHJpbmdcIiB8fCBuYW1lLnNwbGl0KFwiLlwiKS5sZW5ndGggIT09IDMpIHtcbiAgICAgICAgdGhyb3cgbmV3IEVycm9yKFwiQ2FsbEJ5TmFtZSByZXF1aXJlcyBhIHN0cmluZyBpbiB0aGUgZm9ybWF0ICdwYWNrYWdlLnN0cnVjdC5tZXRob2QnXCIpO1xuICAgIH1cbiAgICAvLyBTcGxpdCBpbnB1dHNcbiAgICBsZXQgcGFydHMgPSBuYW1lLnNwbGl0KFwiLlwiKTtcblxuICAgIHJldHVybiBjYWxsQmluZGluZyhDYWxsQmluZGluZywge1xuICAgICAgICBwYWNrYWdlTmFtZTogcGFydHNbMF0sXG4gICAgICAgIHN0cnVjdE5hbWU6IHBhcnRzWzFdLFxuICAgICAgICBtZXRob2ROYW1lOiBwYXJ0c1syXSxcbiAgICAgICAgYXJnczogYXJncyxcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGxCeUlEKG1ldGhvZElELCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKENhbGxCaW5kaW5nLCB7XG4gICAgICAgIG1ldGhvZElEOiBtZXRob2RJRCxcbiAgICAgICAgYXJnczogYXJncyxcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBDYWxsIGEgcGx1Z2luIG1ldGhvZFxuICogQHBhcmFtIHtzdHJpbmd9IHBsdWdpbk5hbWUgLSBuYW1lIG9mIHRoZSBwbHVnaW5cbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gbmFtZSBvZiB0aGUgbWV0aG9kXG4gKiBAcGFyYW0gey4uLmFueX0gYXJncyAtIGFyZ3VtZW50cyB0byBwYXNzIHRvIHRoZSBtZXRob2RcbiAqIEByZXR1cm5zIHtQcm9taXNlPGFueT59IC0gcHJvbWlzZSB0aGF0IHJlc29sdmVzIHdpdGggdGhlIHJlc3VsdFxuICovXG5leHBvcnQgZnVuY3Rpb24gUGx1Z2luKHBsdWdpbk5hbWUsIG1ldGhvZE5hbWUsIC4uLmFyZ3MpIHtcbiAgICByZXR1cm4gY2FsbEJpbmRpbmcoQ2FsbEJpbmRpbmcsIHtcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxuICAgICAgICBzdHJ1Y3ROYW1lOiBwbHVnaW5OYW1lLFxuICAgICAgICBtZXRob2ROYW1lOiBtZXRob2ROYW1lLFxuICAgICAgICBhcmdzOiBhcmdzLFxuICAgIH0pO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2l6ZX0gU2l6ZVxuICogQHR5cGVkZWYge2ltcG9ydChcIi4uL2FwaS90eXBlc1wiKS5Qb3NpdGlvbn0gUG9zaXRpb25cbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2NyZWVufSBTY3JlZW5cbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBXaW5kb3dDZW50ZXIgPSAwO1xubGV0IFdpbmRvd1NldFRpdGxlID0gMTtcbmxldCBXaW5kb3dGdWxsc2NyZWVuID0gMjtcbmxldCBXaW5kb3dVbkZ1bGxzY3JlZW4gPSAzO1xubGV0IFdpbmRvd1NldFNpemUgPSA0O1xubGV0IFdpbmRvd1NpemUgPSA1O1xubGV0IFdpbmRvd1NldE1heFNpemUgPSA2O1xubGV0IFdpbmRvd1NldE1pblNpemUgPSA3O1xubGV0IFdpbmRvd1NldEFsd2F5c09uVG9wID0gODtcbmxldCBXaW5kb3dTZXRSZWxhdGl2ZVBvc2l0aW9uID0gOTtcbmxldCBXaW5kb3dSZWxhdGl2ZVBvc2l0aW9uID0gMTA7XG5sZXQgV2luZG93U2NyZWVuID0gMTE7XG5sZXQgV2luZG93SGlkZSA9IDEyO1xubGV0IFdpbmRvd01heGltaXNlID0gMTM7XG5sZXQgV2luZG93VW5NYXhpbWlzZSA9IDE0O1xubGV0IFdpbmRvd1RvZ2dsZU1heGltaXNlID0gMTU7XG5sZXQgV2luZG93TWluaW1pc2UgPSAxNjtcbmxldCBXaW5kb3dVbk1pbmltaXNlID0gMTc7XG5sZXQgV2luZG93UmVzdG9yZSA9IDE4O1xubGV0IFdpbmRvd1Nob3cgPSAxOTtcbmxldCBXaW5kb3dDbG9zZSA9IDIwO1xubGV0IFdpbmRvd1NldEJhY2tncm91bmRDb2xvdXIgPSAyMTtcbmxldCBXaW5kb3dTZXRSZXNpemFibGUgPSAyMjtcbmxldCBXaW5kb3dXaWR0aCA9IDIzO1xubGV0IFdpbmRvd0hlaWdodCA9IDI0O1xubGV0IFdpbmRvd1pvb21JbiA9IDI1O1xubGV0IFdpbmRvd1pvb21PdXQgPSAyNjtcbmxldCBXaW5kb3dab29tUmVzZXQgPSAyNztcbmxldCBXaW5kb3dHZXRab29tTGV2ZWwgPSAyODtcbmxldCBXaW5kb3dTZXRab29tTGV2ZWwgPSAyOTtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1dpbmRvdyh3aW5kb3dOYW1lKSB7XG4gICAgbGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLldpbmRvdywgd2luZG93TmFtZSk7XG4gICAgcmV0dXJuIHtcblxuICAgICAgICAvKipcbiAgICAgICAgICogQ2VudGVycyB0aGUgd2luZG93LlxuICAgICAgICAgKi9cbiAgICAgICAgQ2VudGVyOiAoKSA9PiB2b2lkIGNhbGwoV2luZG93Q2VudGVyKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgdGl0bGUuXG4gICAgICAgICAqIEBwYXJhbSB0aXRsZVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gdm9pZCBjYWxsKFdpbmRvd1NldFRpdGxlLCB7dGl0bGV9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWFrZXMgdGhlIHdpbmRvdyBmdWxsc2NyZWVuLlxuICAgICAgICAgKi9cbiAgICAgICAgRnVsbHNjcmVlbjogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd0Z1bGxzY3JlZW4pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbmZ1bGxzY3JlZW4gdGhlIHdpbmRvdy5cbiAgICAgICAgICovXG4gICAgICAgIFVuRnVsbHNjcmVlbjogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1VuRnVsbHNjcmVlbiksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHNpemUuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCBUaGUgd2luZG93IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgVGhlIHdpbmRvdyBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKFdpbmRvd1NldFNpemUsIHt3aWR0aCxoZWlnaHR9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgc2l6ZS5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8U2l6ZT59IFRoZSB3aW5kb3cgc2l6ZVxuICAgICAgICAgKi9cbiAgICAgICAgU2l6ZTogKCkgPT4geyByZXR1cm4gY2FsbChXaW5kb3dTaXplKTsgfSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgbWF4aW11bSBzaXplLlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodFxuICAgICAgICAgKi9cbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IHZvaWQgY2FsbChXaW5kb3dTZXRNYXhTaXplLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1pbmltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1pblNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoV2luZG93U2V0TWluU2l6ZSwge3dpZHRoLGhlaWdodH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXG4gICAgICAgICAqIEBwYXJhbSB7Ym9vbGVhbn0gb25Ub3AgV2hldGhlciB0aGUgd2luZG93IHNob3VsZCBiZSBhbHdheXMgb24gdG9wXG4gICAgICAgICAqL1xuICAgICAgICBTZXRBbHdheXNPblRvcDogKG9uVG9wKSA9PiB2b2lkIGNhbGwoV2luZG93U2V0QWx3YXlzT25Ub3AsIHthbHdheXNPblRvcDpvblRvcH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyByZWxhdGl2ZSBwb3NpdGlvbi5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHhcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHlcbiAgICAgICAgICovXG4gICAgICAgIFNldFJlbGF0aXZlUG9zaXRpb246ICh4LCB5KSA9PiBjYWxsKFdpbmRvd1NldFJlbGF0aXZlUG9zaXRpb24sIHt4LHl9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgcG9zaXRpb24uXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFBvc2l0aW9uPn0gVGhlIHdpbmRvdyBwb3NpdGlvblxuICAgICAgICAgKi9cbiAgICAgICAgUmVsYXRpdmVQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbChXaW5kb3dSZWxhdGl2ZVBvc2l0aW9uKTsgfSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBvbi5cbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cbiAgICAgICAgICovXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbChXaW5kb3dTY3JlZW4pOyB9LFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBIaWRlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIEhpZGU6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dIaWRlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWF4aW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dNYXhpbWlzZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNob3cgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgU2hvdzogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1Nob3cpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBDbG9zZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBDbG9zZTogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd0Nsb3NlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVG9nZ2xlIHRoZSB3aW5kb3cgbWF4aW1pc2Ugc3RhdGVcbiAgICAgICAgICovXG4gICAgICAgIFRvZ2dsZU1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoV2luZG93VG9nZ2xlTWF4aW1pc2UpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBVbm1heGltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dVbk1heGltaXNlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWluaW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dNaW5pbWlzZSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFVubWluaW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgVW5NaW5pbWlzZTogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1VuTWluaW1pc2UpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBSZXN0b3JlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFJlc3RvcmU6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dSZXN0b3JlKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSBiYWNrZ3JvdW5kIGNvbG91ciBvZiB0aGUgd2luZG93LlxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gciAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGcgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBiIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gYSAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICovXG4gICAgICAgIFNldEJhY2tncm91bmRDb2xvdXI6IChyLCBnLCBiLCBhKSA9PiB2b2lkIGNhbGwoV2luZG93U2V0QmFja2dyb3VuZENvbG91ciwge3IsIGcsIGIsIGF9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHdoZXRoZXIgdGhlIHdpbmRvdyBjYW4gYmUgcmVzaXplZCBvciBub3RcbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSByZXNpemFibGVcbiAgICAgICAgICovXG4gICAgICAgIFNldFJlc2l6YWJsZTogKHJlc2l6YWJsZSkgPT4gdm9pZCBjYWxsKFdpbmRvd1NldFJlc2l6YWJsZSwge3Jlc2l6YWJsZX0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyB3aWR0aFxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxudW1iZXI+fVxuICAgICAgICAgKi9cbiAgICAgICAgV2lkdGg6ICgpID0+IHsgcmV0dXJuIGNhbGwoV2luZG93V2lkdGgpOyB9LFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBoZWlnaHRcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8bnVtYmVyPn1cbiAgICAgICAgICovXG4gICAgICAgIEhlaWdodDogKCkgPT4geyByZXR1cm4gY2FsbChXaW5kb3dIZWlnaHQpOyB9LFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBab29tIGluIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFpvb21JbjogKCkgPT4gdm9pZCBjYWxsKFdpbmRvd1pvb21JbiksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFpvb20gb3V0IHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIFpvb21PdXQ6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dab29tT3V0KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogUmVzZXQgdGhlIHdpbmRvdyB6b29tXG4gICAgICAgICAqL1xuICAgICAgICBab29tUmVzZXQ6ICgpID0+IHZvaWQgY2FsbChXaW5kb3dab29tUmVzZXQpLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyB6b29tXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPG51bWJlcj59XG4gICAgICAgICAqL1xuICAgICAgICBHZXRab29tTGV2ZWw6ICgpID0+IHsgcmV0dXJuIGNhbGwoV2luZG93R2V0Wm9vbUxldmVsKTsgfSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgem9vbSBsZXZlbFxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gem9vbUxldmVsXG4gICAgICAgICAqL1xuICAgICAgICBTZXRab29tTGV2ZWw6ICh6b29tTGV2ZWwpID0+IHZvaWQgY2FsbChXaW5kb3dTZXRab29tTGV2ZWwsIHt6b29tTGV2ZWx9KSxcbiAgICB9O1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcldpdGhJRChvYmplY3ROYW1lcy5FdmVudHMpO1xubGV0IEV2ZW50RW1pdCA9IDA7XG5cbi8qKlxuICogVGhlIExpc3RlbmVyIGNsYXNzIGRlZmluZXMgYSBsaXN0ZW5lciEgOi0pXG4gKlxuICogQGNsYXNzIExpc3RlbmVyXG4gKi9cbmNsYXNzIExpc3RlbmVyIHtcbiAgICAvKipcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIExpc3RlbmVyLlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xuICAgICAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcbiAgICAgICAgdGhpcy5tYXhDYWxsYmFja3MgPSBtYXhDYWxsYmFja3MgfHwgLTE7XG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxuICAgICAgICB0aGlzLkNhbGxiYWNrID0gKGRhdGEpID0+IHtcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxuICAgICAgICAgICAgaWYgKHRoaXMubWF4Q2FsbGJhY2tzID09PSAtMSkge1xuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIC8vIERlY3JlbWVudCBtYXhDYWxsYmFja3MuIFJldHVybiB0cnVlIGlmIG5vdyAwLCBvdGhlcndpc2UgZmFsc2VcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XG4gICAgICAgIH07XG4gICAgfVxufVxuXG5cbi8qKlxuICogV2FpbHNFdmVudCBkZWZpbmVzIGEgY3VzdG9tIGV2ZW50LiBJdCBpcyBwYXNzZWQgdG8gZXZlbnQgbGlzdGVuZXJzLlxuICpcbiAqIEBjbGFzcyBXYWlsc0V2ZW50XG4gKiBAcHJvcGVydHkge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gKiBAcHJvcGVydHkge2FueX0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxuICovXG5leHBvcnQgY2xhc3MgV2FpbHNFdmVudCB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBXYWlsc0V2ZW50LlxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBuYW1lIC0gTmFtZSBvZiB0aGUgZXZlbnRcbiAgICAgKiBAcGFyYW0ge2FueT1udWxsfSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XG4gICAgICogQG1lbWJlcm9mIFdhaWxzRXZlbnRcbiAgICAgKi9cbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhID0gbnVsbCkge1xuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xuICAgICAgICB0aGlzLmRhdGEgPSBkYXRhO1xuICAgIH1cbn1cblxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEBwYXJhbSB7bnVtYmVyfSBtYXhDYWxsYmFja3NcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcbiAgICBsaXN0ZW5lcnMucHVzaCh0aGlzTGlzdGVuZXIpO1xuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIGV2ZXJ5IHRpbWUgdGhlIGV2ZW50IGlzIGVtaXR0ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XG59XG5cbi8qKlxuICogUmVnaXN0ZXJzIGFuIGV2ZW50IGxpc3RlbmVyIHRoYXQgd2lsbCBiZSBpbnZva2VkIG9uY2UgdGhlbiBkZXN0cm95ZWRcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xuIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcbn1cblxuLyoqXG4gKiBsaXN0ZW5lck9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uXG4gKlxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcbiAqL1xuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XG4gICAgLy8gUmVtb3ZlIGxvY2FsIGxpc3RlbmVyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIH1cbn1cblxuLyoqXG4gKiBkaXNwYXRjaGVzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcbiAqXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBkaXNwYXRjaFdhaWxzRXZlbnQoZXZlbnQpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xuICAgIGlmIChsaXN0ZW5lcnMpIHtcbiAgICAgICAgLy8gaXRlcmF0ZSBsaXN0ZW5lcnMgYW5kIGNhbGwgY2FsbGJhY2suIElmIGNhbGxiYWNrIHJldHVybnMgdHJ1ZSwgcmVtb3ZlIGxpc3RlbmVyXG4gICAgICAgIGxldCB0b1JlbW92ZSA9IFtdO1xuICAgICAgICBsaXN0ZW5lcnMuZm9yRWFjaChsaXN0ZW5lciA9PiB7XG4gICAgICAgICAgICBsZXQgcmVtb3ZlID0gbGlzdGVuZXIuQ2FsbGJhY2soZXZlbnQpO1xuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9KTtcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xuICAgICAgICBpZiAodG9SZW1vdmUubGVuZ3RoID4gMCkge1xuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudC5uYW1lKTtcbiAgICAgICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH1cbiAgICB9XG59XG5cbi8qKlxuICogT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT24sXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXG4gKlxuIFt2MyBDSEFOR0VdIE9mZiBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtICB7Li4uc3RyaW5nfSBhZGRpdGlvbmFsRXZlbnROYW1lc1xuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcbiAgICBsZXQgZXZlbnRzVG9SZW1vdmUgPSBbZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lc107XG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcbiAgICB9KTtcbn1cblxuLyoqXG4gKiBPZmZBbGwgdW5yZWdpc3RlcnMgYWxsIGxpc3RlbmVyc1xuICogW3YzIENIQU5HRV0gT2ZmQWxsIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcbiAqXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7XG4gICAgZXZlbnRMaXN0ZW5lcnMuY2xlYXIoKTtcbn1cblxuLyoqXG4gKiBFbWl0IGFuIGV2ZW50XG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50IFRoZSBldmVudCB0byBlbWl0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XG4gICAgdm9pZCBjYWxsKEV2ZW50RW1pdCwgZXZlbnQpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5NZXNzYWdlRGlhbG9nT3B0aW9uc30gTWVzc2FnZURpYWxvZ09wdGlvbnNcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5PcGVuRGlhbG9nT3B0aW9uc30gT3BlbkRpYWxvZ09wdGlvbnNcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5TYXZlRGlhbG9nT3B0aW9uc30gU2F2ZURpYWxvZ09wdGlvbnNcbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJXaXRoSUQsIG9iamVjdE5hbWVzfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkRpYWxvZyk7XG5cbmxldCBEaWFsb2dJbmZvID0gMDtcbmxldCBEaWFsb2dXYXJuaW5nID0gMTtcbmxldCBEaWFsb2dFcnJvciA9IDI7XG5sZXQgRGlhbG9nUXVlc3Rpb24gPSAzO1xubGV0IERpYWxvZ09wZW5GaWxlID0gNDtcbmxldCBEaWFsb2dTYXZlRmlsZSA9IDU7XG5cblxubGV0IGRpYWxvZ1Jlc3BvbnNlcyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcbiAgICBsZXQgcmVzdWx0O1xuICAgIGRvIHtcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XG4gICAgfSB3aGlsZSAoZGlhbG9nUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcbiAgICByZXR1cm4gcmVzdWx0O1xufVxuXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgaWYgKGlzSlNPTikge1xuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xuICAgICAgICB9XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cbmV4cG9ydCBmdW5jdGlvbiBkaWFsb2dFcnJvckNhbGxiYWNrKGlkLCBtZXNzYWdlKSB7XG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucykge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XG4gICAgICAgIG9wdGlvbnNbXCJkaWFsb2ctaWRcIl0gPSBpZDtcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xuICAgICAgICAgICAgcmVqZWN0KGVycm9yKTtcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuXG4vKipcbiAqIFNob3dzIGFuIEluZm8gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBJbmZvKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKERpYWxvZ0luZm8sIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGEgV2FybmluZyBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFdhcm5pbmcob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coRGlhbG9nV2FybmluZywgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYW4gRXJyb3IgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhEaWFsb2dFcnJvciwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYSBRdWVzdGlvbiBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKERpYWxvZ1F1ZXN0aW9uLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhbiBPcGVuIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtPcGVuRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nW118c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZSBvciBhbiBhcnJheSBvZiBzZWxlY3RlZCBmaWxlcyBpZiBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbiBpcyB0cnVlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9wZW5GaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKERpYWxvZ09wZW5GaWxlLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhIFNhdmUgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge1NhdmVEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKERpYWxvZ1NhdmVGaWxlLCBvcHRpb25zKTtcbn1cblxuIiwgImltcG9ydCB7bmV3UnVudGltZUNhbGxlcldpdGhJRCwgb2JqZWN0TmFtZXN9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyV2l0aElEKG9iamVjdE5hbWVzLkNvbnRleHRNZW51KTtcblxubGV0IENvbnRleHRNZW51T3BlbiA9IDA7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHZvaWQgY2FsbChDb250ZXh0TWVudU9wZW4sIHtpZCwgeCwgeSwgZGF0YX0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gc2V0dXBDb250ZXh0TWVudXMoKSB7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcbn1cblxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50KSB7XG4gICAgLy8gQ2hlY2sgZm9yIGN1c3RvbSBjb250ZXh0IG1lbnVcbiAgICBsZXQgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcbiAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIik7XG4gICAgY3VzdG9tQ29udGV4dE1lbnUgPSBjdXN0b21Db250ZXh0TWVudSA/IGN1c3RvbUNvbnRleHRNZW51LnRyaW0oKSA6IFwiXCI7XG4gICAgaWYgKGN1c3RvbUNvbnRleHRNZW51KSB7XG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgIGxldCBjdXN0b21Db250ZXh0TWVudURhdGEgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnUtZGF0YVwiKTtcbiAgICAgICAgb3BlbkNvbnRleHRNZW51KGN1c3RvbUNvbnRleHRNZW51LCBldmVudC5jbGllbnRYLCBldmVudC5jbGllbnRZLCBjdXN0b21Db250ZXh0TWVudURhdGEpO1xuICAgICAgICByZXR1cm5cbiAgICB9XG5cbiAgICBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KTtcbn1cblxuXG4vKlxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBhdXRvOyAoZGVmYXVsdCkgd2lsbCBzaG93IHRoZSBkZWZhdWx0IGNvbnRleHQgbWVudSBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZSBPUiB0ZXh0IGhhcyBiZWVuIHNlbGVjdGVkIE9SIGVsZW1lbnQgaXMgaW5wdXQgb3IgdGV4dGFyZWFcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvdzsgd2lsbCBhbHdheXMgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogaGlkZTsgd2lsbCBhbHdheXMgaGlkZSB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcblxuVGhpcyBydWxlIGlzIGluaGVyaXRlZCBsaWtlIG5vcm1hbCBDU1MgcnVsZXMsIHNvIG5lc3Rpbmcgd29ya3MgYXMgZXhwZWN0ZWRcbiovXG5mdW5jdGlvbiBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KSB7XG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XG4gICAgaWYgKERFQlVHKSB7XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XG4gICAgY29uc3QgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcbiAgICBjb25zdCBjb21wdXRlZFN0eWxlID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCk7XG4gICAgY29uc3QgZGVmYXVsdENvbnRleHRNZW51QWN0aW9uID0gY29tcHV0ZWRTdHlsZS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1kZWZhdWx0LWNvbnRleHRtZW51XCIpLnRyaW0oKTtcbiAgICBzd2l0Y2ggKGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbikge1xuICAgICAgICBjYXNlIFwic2hvd1wiOlxuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICBjYXNlIFwiaGlkZVwiOlxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgZGVmYXVsdDpcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXG4gICAgICAgICAgICBpZiAoZWxlbWVudC5pc0NvbnRlbnRFZGl0YWJsZSkge1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZFxuICAgICAgICAgICAgY29uc3Qgc2VsZWN0aW9uID0gd2luZG93LmdldFNlbGVjdGlvbigpO1xuICAgICAgICAgICAgY29uc3QgaGFzU2VsZWN0aW9uID0gKHNlbGVjdGlvbi50b1N0cmluZygpLmxlbmd0aCA+IDApXG4gICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uKSB7XG4gICAgICAgICAgICAgICAgZm9yIChsZXQgaSA9IDA7IGkgPCBzZWxlY3Rpb24ucmFuZ2VDb3VudDsgaSsrKSB7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJhbmdlID0gc2VsZWN0aW9uLmdldFJhbmdlQXQoaSk7XG4gICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3RzID0gcmFuZ2UuZ2V0Q2xpZW50UmVjdHMoKTtcbiAgICAgICAgICAgICAgICAgICAgZm9yIChsZXQgaiA9IDA7IGogPCByZWN0cy5sZW5ndGg7IGorKykge1xuICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgcmVjdCA9IHJlY3RzW2pdO1xuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXG4gICAgICAgICAgICBpZiAoZWxlbWVudC50YWdOYW1lID09PSBcIklOUFVUXCIgfHwgZWxlbWVudC50YWdOYW1lID09PSBcIlRFWFRBUkVBXCIpIHtcbiAgICAgICAgICAgICAgICBpZiAoaGFzU2VsZWN0aW9uIHx8ICghZWxlbWVudC5yZWFkT25seSAmJiAhZWxlbWVudC5kaXNhYmxlZCkpIHtcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgIH1cblxuICAgICAgICAgICAgLy8gaGlkZSBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcbiAgICB9XG59XG4iLCAiXG5pbXBvcnQge0VtaXQsIFdhaWxzRXZlbnR9IGZyb20gXCIuL2V2ZW50c1wiO1xuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lLCBkYXRhPW51bGwpIHtcbiAgICBsZXQgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpO1xuICAgIEVtaXQoZXZlbnQpO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxFdmVudExpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtZXZlbnRdJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWV2ZW50Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgRGV0YWNoZWQ6IGZhbHNlLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIH07XG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcblxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9KTtcbn1cblxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcbiAgICBpZiAod2FpbHMuV2luZG93W21ldGhvZF0gPT09IHVuZGVmaW5lZCkge1xuICAgICAgICBjb25zb2xlLmxvZyhcIldpbmRvdyBtZXRob2QgXCIgKyBtZXRob2QgKyBcIiBub3QgZm91bmRcIik7XG4gICAgfVxuICAgIHdhaWxzLldpbmRvd1ttZXRob2RdKCk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1t3bWwtd2luZG93XScpO1xuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC13aW5kb3cnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xuXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgIH0pO1xuICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgIH1cbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG5mdW5jdGlvbiBhZGRXTUxPcGVuQnJvd3Nlckxpc3RlbmVyKCkge1xuICAgIGNvbnN0IGVsZW1lbnRzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnW3dtbC1vcGVudXJsXScpO1xuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcbiAgICAgICAgY29uc3QgdXJsID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ3dtbC1vcGVudXJsJyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnd21sLWNvbmZpcm0nKTtcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCd3bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XG4gICAgICAgICAgICAgICAgICAgICAgICB2b2lkIHdhaWxzLkJyb3dzZXIuT3BlblVSTCh1cmwpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgdm9pZCB3YWlscy5Ccm93c2VyLk9wZW5VUkwodXJsKTtcbiAgICAgICAgfTtcblxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXG4gICAgICAgIGVsZW1lbnQucmVtb3ZlRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG5cbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxuICAgICAgICBlbGVtZW50LmFkZEV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuICAgIH0pO1xufVxuXG5leHBvcnQgZnVuY3Rpb24gcmVsb2FkV01MKCkge1xuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MV2luZG93TGlzdGVuZXJzKCk7XG4gICAgYWRkV01MT3BlbkJyb3dzZXJMaXN0ZW5lcigpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8vIGRlZmluZWQgaW4gdGhlIFRhc2tmaWxlXG5leHBvcnQgbGV0IGludm9rZSA9IGZ1bmN0aW9uKGlucHV0KSB7XG4gICAgaWYoV0lORE9XUykge1xuICAgICAgICBjaHJvbWUud2Vidmlldy5wb3N0TWVzc2FnZShpbnB1dCk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgd2Via2l0Lm1lc3NhZ2VIYW5kbGVycy5leHRlcm5hbC5wb3N0TWVzc2FnZShpbnB1dCk7XG4gICAgfVxufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmxldCBmbGFncyA9IG5ldyBNYXAoKTtcblxuZnVuY3Rpb24gY29udmVydFRvTWFwKG9iaikge1xuICAgIGNvbnN0IG1hcCA9IG5ldyBNYXAoKTtcblxuICAgIGZvciAoY29uc3QgW2tleSwgdmFsdWVdIG9mIE9iamVjdC5lbnRyaWVzKG9iaikpIHtcbiAgICAgICAgaWYgKHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcgJiYgdmFsdWUgIT09IG51bGwpIHtcbiAgICAgICAgICAgIG1hcC5zZXQoa2V5LCBjb252ZXJ0VG9NYXAodmFsdWUpKTsgLy8gUmVjdXJzaXZlbHkgY29udmVydCBuZXN0ZWQgb2JqZWN0XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICBtYXAuc2V0KGtleSwgdmFsdWUpO1xuICAgICAgICB9XG4gICAgfVxuXG4gICAgcmV0dXJuIG1hcDtcbn1cblxuZmV0Y2goXCIvd2FpbHMvZmxhZ3NcIikudGhlbigocmVzcG9uc2UpID0+IHtcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xuICAgICAgICBmbGFncyA9IGNvbnZlcnRUb01hcChkYXRhKTtcbiAgICB9KTtcbn0pO1xuXG5cbmZ1bmN0aW9uIGdldFZhbHVlRnJvbU1hcChrZXlTdHJpbmcpIHtcbiAgICBjb25zdCBrZXlzID0ga2V5U3RyaW5nLnNwbGl0KCcuJyk7XG4gICAgbGV0IHZhbHVlID0gZmxhZ3M7XG5cbiAgICBmb3IgKGNvbnN0IGtleSBvZiBrZXlzKSB7XG4gICAgICAgIGlmICh2YWx1ZSBpbnN0YW5jZW9mIE1hcCkge1xuICAgICAgICAgICAgdmFsdWUgPSB2YWx1ZS5nZXQoa2V5KTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHZhbHVlID0gdmFsdWVba2V5XTtcbiAgICAgICAgfVxuXG4gICAgICAgIGlmICh2YWx1ZSA9PT0gdW5kZWZpbmVkKSB7XG4gICAgICAgICAgICBicmVhaztcbiAgICAgICAgfVxuICAgIH1cblxuICAgIHJldHVybiB2YWx1ZTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5U3RyaW5nKSB7XG4gICAgcmV0dXJuIGdldFZhbHVlRnJvbU1hcChrZXlTdHJpbmcpO1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbmltcG9ydCB7aW52b2tlfSBmcm9tIFwiLi9pbnZva2VcIjtcbmltcG9ydCB7R2V0RmxhZ30gZnJvbSBcIi4vZmxhZ3NcIjtcblxubGV0IHNob3VsZERyYWcgPSBmYWxzZTtcblxuZXhwb3J0IGZ1bmN0aW9uIGRyYWdUZXN0KGUpIHtcbiAgICBsZXQgdmFsID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLXdlYmtpdC1hcHAtcmVnaW9uXCIpO1xuICAgIGlmICh2YWwpIHtcbiAgICAgICAgdmFsID0gdmFsLnRyaW0oKTtcbiAgICB9XG5cbiAgICBpZiAodmFsICE9PSBcImRyYWdcIikge1xuICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgfVxuXG4gICAgLy8gT25seSBwcm9jZXNzIHRoZSBwcmltYXJ5IGJ1dHRvblxuICAgIGlmIChlLmJ1dHRvbnMgIT09IDEpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIHJldHVybiBlLmRldGFpbCA9PT0gMTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHNldHVwRHJhZygpIHtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgb25Nb3VzZURvd24pO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBvbk1vdXNlTW92ZSk7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNldXAnLCBvbk1vdXNlVXApO1xufVxuXG5sZXQgcmVzaXplRWRnZSA9IG51bGw7XG5sZXQgcmVzaXphYmxlID0gZmFsc2U7XG5cbmV4cG9ydCBmdW5jdGlvbiBzZXRSZXNpemFibGUodmFsdWUpIHtcbiAgICByZXNpemFibGUgPSB2YWx1ZTtcbn1cblxuZnVuY3Rpb24gdGVzdFJlc2l6ZShlKSB7XG4gICAgaWYoIHJlc2l6ZUVkZ2UgKSB7XG4gICAgICAgIGludm9rZShcInJlc2l6ZTpcIiArIHJlc2l6ZUVkZ2UpO1xuICAgICAgICByZXR1cm4gdHJ1ZVxuICAgIH1cbiAgICByZXR1cm4gZmFsc2U7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VEb3duKGUpIHtcblxuICAgIC8vIENoZWNrIGZvciByZXNpemluZyBvbiBXaW5kb3dzXG4gICAgaWYoIFdJTkRPV1MgKSB7XG4gICAgICAgIGlmICh0ZXN0UmVzaXplKCkpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgIH1cbiAgICBpZiAoZHJhZ1Rlc3QoZSkpIHtcbiAgICAgICAgLy8gSWdub3JlIGRyYWcgb24gc2Nyb2xsYmFyc1xuICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIH1cbiAgICAgICAgc2hvdWxkRHJhZyA9IHRydWU7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gb25Nb3VzZVVwKGUpIHtcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgIGVuZERyYWcoKTtcbiAgICB9XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBlbmREcmFnKCkge1xuICAgIGRvY3VtZW50LmJvZHkuc3R5bGUuY3Vyc29yID0gJ2RlZmF1bHQnO1xuICAgIHNob3VsZERyYWcgPSBmYWxzZTtcbn1cblxuZnVuY3Rpb24gc2V0UmVzaXplKGN1cnNvcikge1xuICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBjdXJzb3IgfHwgZGVmYXVsdEN1cnNvcjtcbiAgICByZXNpemVFZGdlID0gY3Vyc29yO1xufVxuXG5mdW5jdGlvbiBvbk1vdXNlTW92ZShlKSB7XG4gICAgaWYgKHNob3VsZERyYWcpIHtcbiAgICAgICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xuICAgICAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xuICAgICAgICBpZiAobW91c2VQcmVzc2VkID4gMCkge1xuICAgICAgICAgICAgaW52b2tlKFwiZHJhZ1wiKTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgaWYgKFdJTkRPV1MpIHtcbiAgICAgICAgaWYgKHJlc2l6YWJsZSkge1xuICAgICAgICAgICAgaGFuZGxlUmVzaXplKGUpO1xuICAgICAgICB9XG4gICAgfVxufVxuXG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xuXG5mdW5jdGlvbiBoYW5kbGVSZXNpemUoZSkge1xuICAgIGxldCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGxldCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcblxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhc1xuICAgIGxldCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAxMDtcblxuICAgIGxldCByaWdodEJvcmRlciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xuICAgIGxldCBib3R0b21Cb3JkZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG5cbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lcnNcbiAgICBsZXQgcmlnaHRDb3JuZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgbGVmdENvcm5lciA9IGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgdG9wQ29ybmVyID0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcbiAgICBsZXQgYm90dG9tQ29ybmVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcblxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XG4gICAgaWYgKCFsZWZ0Qm9yZGVyICYmICFyaWdodEJvcmRlciAmJiAhdG9wQm9yZGVyICYmICFib3R0b21Cb3JkZXIgJiYgcmVzaXplRWRnZSAhPT0gdW5kZWZpbmVkKSB7XG4gICAgICAgIHNldFJlc2l6ZSgpO1xuICAgIH1cbiAgICAvLyBBZGp1c3RlZCBmb3IgY29ybmVyIGFyZWFzXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIHRvcENvcm5lcikgc2V0UmVzaXplKFwibnctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHRvcEJvcmRlcikgc2V0UmVzaXplKFwibi1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XG59XG4iLCAiLypcbiBfICAgICBfXyAgICAgXyBfX1xufCB8ICAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cblxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gJy4vY2xpcGJvYXJkJztcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gJy4vYXBwbGljYXRpb24nO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tICcuL3NjcmVlbnMnO1xuaW1wb3J0ICogYXMgU3lzdGVtIGZyb20gJy4vc3lzdGVtJztcbmltcG9ydCAqIGFzIEJyb3dzZXIgZnJvbSAnLi9icm93c2VyJztcbmltcG9ydCB7UGx1Z2luLCBDYWxsLCBjYWxsRXJyb3JDYWxsYmFjaywgY2FsbENhbGxiYWNrLCBDYWxsQnlJRCwgQ2FsbEJ5TmFtZX0gZnJvbSBcIi4vY2FsbHNcIjtcbmltcG9ydCB7Y2xpZW50SWR9IGZyb20gJy4vcnVudGltZSc7XG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQge2Rpc3BhdGNoV2FpbHNFdmVudCwgRW1pdCwgT2ZmLCBPZmZBbGwsIE9uLCBPbmNlLCBPbk11bHRpcGxlfSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7ZGlhbG9nQ2FsbGJhY2ssIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssIEVycm9yLCBJbmZvLCBPcGVuRmlsZSwgUXVlc3Rpb24sIFNhdmVGaWxlLCBXYXJuaW5nLH0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuaW1wb3J0IHtzZXR1cENvbnRleHRNZW51c30gZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCB7cmVsb2FkV01MfSBmcm9tIFwiLi93bWxcIjtcbmltcG9ydCB7c2V0dXBEcmFnLCBlbmREcmFnLCBzZXRSZXNpemFibGV9IGZyb20gXCIuL2RyYWdcIjtcblxud2luZG93LndhaWxzID0ge1xuICAgIC4uLm5ld1J1bnRpbWUobnVsbCksXG4gICAgQ2FwYWJpbGl0aWVzOiB7fSxcbiAgICBjbGllbnRJZDogY2xpZW50SWQsXG59O1xuXG5mZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIikudGhlbigocmVzcG9uc2UpID0+IHtcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xuICAgICAgICB3aW5kb3cud2FpbHMuQ2FwYWJpbGl0aWVzID0gZGF0YTtcbiAgICB9KTtcbn0pO1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy5fd2FpbHMgPSB7XG4gICAgZGlhbG9nQ2FsbGJhY2ssXG4gICAgZGlhbG9nRXJyb3JDYWxsYmFjayxcbiAgICBkaXNwYXRjaFdhaWxzRXZlbnQsXG4gICAgY2FsbENhbGxiYWNrLFxuICAgIGNhbGxFcnJvckNhbGxiYWNrLFxuICAgIGVuZERyYWcsXG4gICAgc2V0UmVzaXphYmxlLFxufTtcblxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWUod2luZG93TmFtZSkge1xuICAgIHJldHVybiB7XG4gICAgICAgIENsaXBib2FyZDoge1xuICAgICAgICAgICAgLi4uQ2xpcGJvYXJkXG4gICAgICAgIH0sXG4gICAgICAgIEFwcGxpY2F0aW9uOiB7XG4gICAgICAgICAgICAuLi5BcHBsaWNhdGlvbixcbiAgICAgICAgICAgIEdldFdpbmRvd0J5TmFtZSh3aW5kb3dOYW1lKSB7XG4gICAgICAgICAgICAgICAgcmV0dXJuIG5ld1J1bnRpbWUod2luZG93TmFtZSk7XG4gICAgICAgICAgICB9XG4gICAgICAgIH0sXG4gICAgICAgIFN5c3RlbSxcbiAgICAgICAgU2NyZWVucyxcbiAgICAgICAgQnJvd3NlcixcbiAgICAgICAgQ2FsbCxcbiAgICAgICAgQ2FsbEJ5SUQsXG4gICAgICAgIENhbGxCeU5hbWUsXG4gICAgICAgIFBsdWdpbixcbiAgICAgICAgV01MOiB7XG4gICAgICAgICAgICBSZWxvYWQ6IHJlbG9hZFdNTCxcbiAgICAgICAgfSxcbiAgICAgICAgRGlhbG9nOiB7XG4gICAgICAgICAgICBJbmZvLFxuICAgICAgICAgICAgV2FybmluZyxcbiAgICAgICAgICAgIEVycm9yLFxuICAgICAgICAgICAgUXVlc3Rpb24sXG4gICAgICAgICAgICBPcGVuRmlsZSxcbiAgICAgICAgICAgIFNhdmVGaWxlLFxuICAgICAgICB9LFxuICAgICAgICBFdmVudHM6IHtcbiAgICAgICAgICAgIEVtaXQsXG4gICAgICAgICAgICBPbixcbiAgICAgICAgICAgIE9uY2UsXG4gICAgICAgICAgICBPbk11bHRpcGxlLFxuICAgICAgICAgICAgT2ZmLFxuICAgICAgICAgICAgT2ZmQWxsLFxuICAgICAgICB9LFxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyh3aW5kb3dOYW1lKSxcbiAgICB9O1xufVxuXG5pZiAoREVCVUcpIHtcbiAgICBjb25zb2xlLmxvZyhcIldhaWxzIHYzLjAuMCBEZWJ1ZyBNb2RlIEVuYWJsZWRcIik7XG59XG5cbnNldHVwQ29udGV4dE1lbnVzKCk7XG5zZXR1cERyYWcoKTtcblxuZG9jdW1lbnQuYWRkRXZlbnRMaXN0ZW5lcihcIkRPTUNvbnRlbnRMb2FkZWRcIiwgZnVuY3Rpb24oKSB7XG4gICAgcmVsb2FkV01MKCk7XG59KTtcbiJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDQUEsTUFBSSxjQUNGO0FBV0ssTUFBSSxTQUFTLENBQUMsT0FBTyxPQUFPO0FBQ2pDLFFBQUksS0FBSztBQUNULFFBQUksSUFBSTtBQUNSLFdBQU8sS0FBSztBQUNWLFlBQU0s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LEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQSxhQUFhO0FBQUEsUUFDVCxHQUFHO0FBQUEsUUFDSCxnQkFBZ0JDLGFBQVk7QUFDeEIsaUJBQU8sV0FBV0EsV0FBVTtBQUFBLFFBQ2hDO0FBQUEsTUFDSjtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFFBQVE7QUFBQSxNQUNaO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxVQUFVO0FBQUEsSUFDaEM7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEsb0JBQWtCO0FBQ2xCLFlBQVU7QUFFVixXQUFTLGlCQUFpQixvQkFBb0IsV0FBVztBQUNyRCxjQUFVO0FBQUEsRUFDZCxDQUFDOyIsCiAgIm5hbWVzIjogWyJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJyZXNpemFibGUiLCAiY2FsbCIsICJldmVudE5hbWUiLCAiY2FsbCIsICJnZW5lcmF0ZUlEIiwgIkVycm9yIiwgImNhbGwiLCAid2luZG93TmFtZSIsICJFcnJvciJdCn0K
