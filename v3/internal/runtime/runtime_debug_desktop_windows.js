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

  // desktop/runtime.js
  var runtimeURL = window.location.origin + "/wails/runtime";
  function runtimeCall(method, windowName, args) {
    let url = new URL(runtimeURL);
    url.searchParams.append("method", method);
    if (args) {
      url.searchParams.append("args", JSON.stringify(args));
    }
    let fetchOptions = {
      headers: {}
    };
    if (windowName) {
      fetchOptions.headers["x-wails-window-name"] = windowName;
    }
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
  function newRuntimeCaller(object, windowName) {
    return function(method, args = null) {
      return runtimeCall(object + "." + method, windowName, args);
    };
  }

  // desktop/clipboard.js
  var call = newRuntimeCaller("clipboard");
  function SetText(text) {
    void call("SetText", { text });
  }
  function Text() {
    return call("Text");
  }

  // desktop/application.js
  var application_exports = {};
  __export(application_exports, {
    Hide: () => Hide,
    Quit: () => Quit,
    Show: () => Show
  });
  var call2 = newRuntimeCaller("application");
  function Hide() {
    void call2("Hide");
  }
  function Show() {
    void call2("Show");
  }
  function Quit() {
    void call2("Quit");
  }

  // desktop/log.js
  var log_exports = {};
  __export(log_exports, {
    Log: () => Log
  });
  var call3 = newRuntimeCaller("log");
  function Log(message) {
    return call3("Log", message);
  }

  // desktop/screens.js
  var screens_exports = {};
  __export(screens_exports, {
    GetAll: () => GetAll,
    GetCurrent: () => GetCurrent,
    GetPrimary: () => GetPrimary
  });
  var call4 = newRuntimeCaller("screens");
  function GetAll() {
    return call4("GetAll");
  }
  function GetPrimary() {
    return call4("GetPrimary");
  }
  function GetCurrent() {
    return call4("GetCurrent");
  }

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

  // desktop/calls.js
  var call5 = newRuntimeCaller("call");
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
      call5(type, options).catch((error) => {
        reject(error);
        callResponses.delete(id);
      });
    });
  }
  function Call(options) {
    return callBinding("Call", options);
  }
  function Plugin(pluginName, methodName, ...args) {
    return callBinding("Call", {
      packageName: "wails-plugins",
      structName: pluginName,
      methodName,
      args
    });
  }

  // desktop/window.js
  function newWindow(windowName) {
    let call9 = newRuntimeCaller("window", windowName);
    return {
      // Reload: () => call('WR'),
      // ReloadApp: () => call('WR'),
      // SetSystemDefaultTheme: () => call('WASDT'),
      // SetLightTheme: () => call('WALT'),
      // SetDarkTheme: () => call('WADT'),
      // IsFullscreen: () => call('WIF'),
      // IsMaximized: () => call('WIM'),
      // IsMinimized: () => call('WIMN'),
      // IsWindowed: () => call('WIF'),
      /**
       * Centers the window.
       */
      Center: () => void call9("Center"),
      /**
       * Set the window title.
       * @param title
       */
      SetTitle: (title) => void call9("SetTitle", { title }),
      /**
       * Makes the window fullscreen.
       */
      Fullscreen: () => void call9("Fullscreen"),
      /**
       * Unfullscreen the window.
       */
      UnFullscreen: () => void call9("UnFullscreen"),
      /**
       * Set the window size.
       * @param {number} width The window width
       * @param {number} height The window height
       */
      SetSize: (width, height) => call9("SetSize", { width, height }),
      /**
       * Get the window size.
       * @returns {Promise<Size>} The window size
       */
      Size: () => {
        return call9("Size");
      },
      /**
       * Set the window maximum size.
       * @param {number} width
       * @param {number} height
       */
      SetMaxSize: (width, height) => void call9("SetMaxSize", { width, height }),
      /**
       * Set the window minimum size.
       * @param {number} width
       * @param {number} height
       */
      SetMinSize: (width, height) => void call9("SetMinSize", { width, height }),
      /**
       * Set window to be always on top.
       * @param {boolean} onTop Whether the window should be always on top
       */
      SetAlwaysOnTop: (onTop) => void call9("SetAlwaysOnTop", { alwaysOnTop: onTop }),
      /**
       * Set the window relative position.
       * @param {number} x
       * @param {number} y
       */
      SetRelativePosition: (x, y) => call9("SetRelativePosition", { x, y }),
      /**
       * Get the window position.
       * @returns {Promise<Position>} The window position
       */
      RelativePosition: () => {
        return call9("RelativePosition");
      },
      /**
       * Get the screen the window is on.
       * @returns {Promise<Screen>}
       */
      Screen: () => {
        return call9("Screen");
      },
      /**
       * Hide the window
       */
      Hide: () => void call9("Hide"),
      /**
       * Maximise the window
       */
      Maximise: () => void call9("Maximise"),
      /**
       * Show the window
       */
      Show: () => void call9("Show"),
      /**
       * Close the window
       */
      Close: () => void call9("Close"),
      /**
       * Toggle the window maximise state
       */
      ToggleMaximise: () => void call9("ToggleMaximise"),
      /**
       * Unmaximise the window
       */
      UnMaximise: () => void call9("UnMaximise"),
      /**
       * Minimise the window
       */
      Minimise: () => void call9("Minimise"),
      /**
       * Unminimise the window
       */
      UnMinimise: () => void call9("UnMinimise"),
      /**
       * Restore the window
       */
      Restore: () => void call9("Restore"),
      /**
       * Set the background colour of the window.
       * @param {number} r - A value between 0 and 255
       * @param {number} g - A value between 0 and 255
       * @param {number} b - A value between 0 and 255
       * @param {number} a - A value between 0 and 255
       */
      SetBackgroundColour: (r, g, b, a) => void call9("SetBackgroundColour", { r, g, b, a })
    };
  }

  // desktop/events.js
  var call6 = newRuntimeCaller("events");
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
    void call6("Emit", event);
  }

  // desktop/dialogs.js
  var call7 = newRuntimeCaller("dialog");
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
      call7(type, options).catch((error) => {
        reject(error);
        dialogResponses.delete(id);
      });
    });
  }
  function Info(options) {
    return dialog("Info", options);
  }
  function Warning(options) {
    return dialog("Warning", options);
  }
  function Error2(options) {
    return dialog("Error", options);
  }
  function Question(options) {
    return dialog("Question", options);
  }
  function OpenFile(options) {
    return dialog("OpenFile", options);
  }
  function SaveFile(options) {
    return dialog("SaveFile", options);
  }

  // desktop/contextmenu.js
  var call8 = newRuntimeCaller("contextmenu");
  function openContextMenu(id, x, y, data) {
    void call8("OpenContextMenu", { id, x, y, data });
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
    const elements = document.querySelectorAll("[data-wml-event]");
    elements.forEach(function(element) {
      const eventType = element.getAttribute("data-wml-event");
      const confirm = element.getAttribute("data-wml-confirm");
      const trigger = element.getAttribute("data-wml-trigger") || "click";
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
    const elements = document.querySelectorAll("[data-wml-window]");
    elements.forEach(function(element) {
      const windowMethod = element.getAttribute("data-wml-window");
      const confirm = element.getAttribute("data-wml-confirm");
      const trigger = element.getAttribute("data-wml-trigger") || "click";
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
  function reloadWML() {
    addWMLEventListeners();
    addWMLWindowListeners();
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
      handleResize(e);
    }
  }
  var defaultCursor = "auto";
  function handleResize(e) {
    let resizeHandleHeight = GetFlag("system.resizeHandleHeight") || 5;
    let resizeHandleWidth = GetFlag("system.resizeHandleWidth") || 5;
    let cornerExtra = GetFlag("resizeCornerExtra") || 3;
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
    Capabilities: {}
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
    endDrag
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
      Log: log_exports,
      Screens: screens_exports,
      Call,
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
  document.addEventListener("DOMContentLoaded", function(event) {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvaW52b2tlLmpzIiwgImRlc2t0b3AvZmxhZ3MuanMiLCAiZGVza3RvcC9kcmFnLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImNsaXBib2FyZFwiKTtcblxuLyoqXG4gKiBTZXQgdGhlIENsaXBib2FyZCB0ZXh0XG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBTZXRUZXh0KHRleHQpIHtcbiAgICB2b2lkIGNhbGwoXCJTZXRUZXh0XCIsIHt0ZXh0fSk7XG59XG5cbi8qKlxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJUZXh0XCIpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5jb25zdCBydW50aW1lVVJMID0gd2luZG93LmxvY2F0aW9uLm9yaWdpbiArIFwiL3dhaWxzL3J1bnRpbWVcIjtcblxuZnVuY3Rpb24gcnVudGltZUNhbGwobWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKSB7XG4gICAgbGV0IHVybCA9IG5ldyBVUkwocnVudGltZVVSTCk7XG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcbiAgICBpZiAoYXJncykge1xuICAgICAgICB1cmwuc2VhcmNoUGFyYW1zLmFwcGVuZChcImFyZ3NcIiwgSlNPTi5zdHJpbmdpZnkoYXJncykpO1xuICAgIH1cbiAgICBsZXQgZmV0Y2hPcHRpb25zID0ge1xuICAgICAgICBoZWFkZXJzOiB7fSxcbiAgICB9O1xuICAgIGlmICh3aW5kb3dOYW1lKSB7XG4gICAgICAgIGZldGNoT3B0aW9ucy5oZWFkZXJzW1wieC13YWlscy13aW5kb3ctbmFtZVwiXSA9IHdpbmRvd05hbWU7XG4gICAgfVxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGZldGNoKHVybCwgZmV0Y2hPcHRpb25zKVxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xuICAgICAgICAgICAgICAgIGlmIChyZXNwb25zZS5vaykge1xuICAgICAgICAgICAgICAgICAgICAvLyBjaGVjayBjb250ZW50IHR5cGVcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xuICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuIHJlc3BvbnNlLmpzb24oKTtcbiAgICAgICAgICAgICAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XG4gICAgICAgICAgICAgICAgICAgIH1cbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcbiAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAudGhlbihkYXRhID0+IHJlc29sdmUoZGF0YSkpXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdSdW50aW1lQ2FsbGVyKG9iamVjdCwgd2luZG93TmFtZSkge1xuICAgIHJldHVybiBmdW5jdGlvbiAobWV0aG9kLCBhcmdzPW51bGwpIHtcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcbiAgICB9O1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiYXBwbGljYXRpb25cIik7XG5cbi8qKlxuICogSGlkZSB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XG4gICAgdm9pZCBjYWxsKFwiSGlkZVwiKTtcbn1cblxuLyoqXG4gKiBTaG93IHRoZSBhcHBsaWNhdGlvblxuICovXG5leHBvcnQgZnVuY3Rpb24gU2hvdygpIHtcbiAgICB2b2lkIGNhbGwoXCJTaG93XCIpO1xufVxuXG5cbi8qKlxuICogUXVpdCB0aGUgYXBwbGljYXRpb25cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XG4gICAgdm9pZCBjYWxsKFwiUXVpdFwiKTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImxvZ1wiKTtcblxuLyoqXG4gKiBMb2dzIGEgbWVzc2FnZS5cbiAqIEBwYXJhbSB7bWVzc2FnZX0gTWVzc2FnZSB0byBsb2dcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIExvZyhtZXNzYWdlKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJMb2dcIiwgbWVzc2FnZSk7XG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuU2NyZWVufSBTY3JlZW5cbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwic2NyZWVuc1wiKTtcblxuLyoqXG4gKiBHZXRzIGFsbCBzY3JlZW5zLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuW10+fVxuICovXG5leHBvcnQgZnVuY3Rpb24gR2V0QWxsKCkge1xuICAgIHJldHVybiBjYWxsKFwiR2V0QWxsXCIpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxuICogQHJldHVybnMge1Byb21pc2U8U2NyZWVuPn1cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRQcmltYXJ5XCIpO1xufVxuXG4vKipcbiAqIEdldHMgdGhlIGN1cnJlbnQgYWN0aXZlIHNjcmVlbi5cbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XG4gKiBAY29uc3RydWN0b3JcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XG4gICAgcmV0dXJuIGNhbGwoXCJHZXRDdXJyZW50XCIpO1xufSIsICJsZXQgdXJsQWxwaGFiZXQgPVxuICAndXNlYW5kb20tMjZUMTk4MzQwUFg3NXB4SkFDS1ZFUllNSU5EQlVTSFdPTEZfR1FaYmZnaGprbHF2d3l6cmljdCdcbmV4cG9ydCBsZXQgY3VzdG9tQWxwaGFiZXQgPSAoYWxwaGFiZXQsIGRlZmF1bHRTaXplID0gMjEpID0+IHtcbiAgcmV0dXJuIChzaXplID0gZGVmYXVsdFNpemUpID0+IHtcbiAgICBsZXQgaWQgPSAnJ1xuICAgIGxldCBpID0gc2l6ZVxuICAgIHdoaWxlIChpLS0pIHtcbiAgICAgIGlkICs9IGFscGhhYmV0WyhNYXRoLnJhbmRvbSgpICogYWxwaGFiZXQubGVuZ3RoKSB8IDBdXG4gICAgfVxuICAgIHJldHVybiBpZFxuICB9XG59XG5leHBvcnQgbGV0IG5hbm9pZCA9IChzaXplID0gMjEpID0+IHtcbiAgbGV0IGlkID0gJydcbiAgbGV0IGkgPSBzaXplXG4gIHdoaWxlIChpLS0pIHtcbiAgICBpZCArPSB1cmxBbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIDY0KSB8IDBdXG4gIH1cbiAgcmV0dXJuIGlkXG59XG4iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY2FsbFwiKTtcblxubGV0IGNhbGxSZXNwb25zZXMgPSBuZXcgTWFwKCk7XG5cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGNhbGxSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xuICAgIHJldHVybiByZXN1bHQ7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBjYWxsQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xuICAgIGxldCBwID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcbiAgICAgICAgfVxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XG4gICAgfVxufVxuXG5leHBvcnQgZnVuY3Rpb24gY2FsbEVycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcbiAgICBsZXQgcCA9IGNhbGxSZXNwb25zZXMuZ2V0KGlkKTtcbiAgICBpZiAocCkge1xuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gY2FsbEJpbmRpbmcodHlwZSwgb3B0aW9ucykge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XG4gICAgICAgIGNhbGxSZXNwb25zZXMuc2V0KGlkLCB7cmVzb2x2ZSwgcmVqZWN0fSk7XG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xuICAgICAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgICAgICB9KTtcbiAgICB9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xuICAgIHJldHVybiBjYWxsQmluZGluZyhcIkNhbGxcIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogQ2FsbCBhIHBsdWdpbiBtZXRob2RcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gbmFtZSBvZiB0aGUgcGx1Z2luXG4gKiBAcGFyYW0ge3N0cmluZ30gbWV0aG9kTmFtZSAtIG5hbWUgb2YgdGhlIG1ldGhvZFxuICogQHBhcmFtIHsuLi5hbnl9IGFyZ3MgLSBhcmd1bWVudHMgdG8gcGFzcyB0byB0aGUgbWV0aG9kXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxhbnk+fSAtIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSByZXN1bHRcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFBsdWdpbihwbHVnaW5OYW1lLCBtZXRob2ROYW1lLCAuLi5hcmdzKSB7XG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCB7XG4gICAgICAgIHBhY2thZ2VOYW1lOiBcIndhaWxzLXBsdWdpbnNcIixcbiAgICAgICAgc3RydWN0TmFtZTogcGx1Z2luTmFtZSxcbiAgICAgICAgbWV0aG9kTmFtZTogbWV0aG9kTmFtZSxcbiAgICAgICAgYXJnczogYXJncyxcbiAgICB9KTtcbn0iLCAiLypcbiBfXHQgICBfX1x0ICBfIF9fXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcbiovXG5cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuLyoqXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNpemV9IFNpemVcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuUG9zaXRpb259IFBvc2l0aW9uXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXG4gKi9cblxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XG5cbmV4cG9ydCBmdW5jdGlvbiBuZXdXaW5kb3cod2luZG93TmFtZSkge1xuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCB3aW5kb3dOYW1lKTtcbiAgICByZXR1cm4ge1xuICAgICAgICAvLyBSZWxvYWQ6ICgpID0+IGNhbGwoJ1dSJyksXG4gICAgICAgIC8vIFJlbG9hZEFwcDogKCkgPT4gY2FsbCgnV1InKSxcbiAgICAgICAgLy8gU2V0U3lzdGVtRGVmYXVsdFRoZW1lOiAoKSA9PiBjYWxsKCdXQVNEVCcpLFxuICAgICAgICAvLyBTZXRMaWdodFRoZW1lOiAoKSA9PiBjYWxsKCdXQUxUJyksXG4gICAgICAgIC8vIFNldERhcmtUaGVtZTogKCkgPT4gY2FsbCgnV0FEVCcpLFxuICAgICAgICAvLyBJc0Z1bGxzY3JlZW46ICgpID0+IGNhbGwoJ1dJRicpLFxuICAgICAgICAvLyBJc01heGltaXplZDogKCkgPT4gY2FsbCgnV0lNJyksXG4gICAgICAgIC8vIElzTWluaW1pemVkOiAoKSA9PiBjYWxsKCdXSU1OJyksXG4gICAgICAgIC8vIElzV2luZG93ZWQ6ICgpID0+IGNhbGwoJ1dJRicpLFxuXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIENlbnRlcnMgdGhlIHdpbmRvdy5cbiAgICAgICAgICovXG4gICAgICAgIENlbnRlcjogKCkgPT4gdm9pZCBjYWxsKCdDZW50ZXInKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgdGl0bGUuXG4gICAgICAgICAqIEBwYXJhbSB0aXRsZVxuICAgICAgICAgKi9cbiAgICAgICAgU2V0VGl0bGU6ICh0aXRsZSkgPT4gdm9pZCBjYWxsKCdTZXRUaXRsZScsIHt0aXRsZX0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBNYWtlcyB0aGUgd2luZG93IGZ1bGxzY3JlZW4uXG4gICAgICAgICAqL1xuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ0Z1bGxzY3JlZW4nKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVW5mdWxsc2NyZWVuIHRoZSB3aW5kb3cuXG4gICAgICAgICAqL1xuICAgICAgICBVbkZ1bGxzY3JlZW46ICgpID0+IHZvaWQgY2FsbCgnVW5GdWxsc2NyZWVuJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHNpemUuXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCBUaGUgd2luZG93IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHQgVGhlIHdpbmRvdyBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRTaXplJywge3dpZHRoLGhlaWdodH0pLFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBzaXplLlxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxTaXplPn0gVGhlIHdpbmRvdyBzaXplXG4gICAgICAgICAqL1xuICAgICAgICBTaXplOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTaXplJyk7IH0sXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1heGltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1heFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoJ1NldE1heFNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1pbmltdW0gc2l6ZS5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHdpZHRoXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcbiAgICAgICAgICovXG4gICAgICAgIFNldE1pblNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiB2b2lkIGNhbGwoJ1NldE1pblNpemUnLCB7d2lkdGgsaGVpZ2h0fSksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB3aW5kb3cgdG8gYmUgYWx3YXlzIG9uIHRvcC5cbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSBvblRvcCBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGFsd2F5cyBvbiB0b3BcbiAgICAgICAgICovXG4gICAgICAgIFNldEFsd2F5c09uVG9wOiAob25Ub3ApID0+IHZvaWQgY2FsbCgnU2V0QWx3YXlzT25Ub3AnLCB7YWx3YXlzT25Ub3A6b25Ub3B9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgcmVsYXRpdmUgcG9zaXRpb24uXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB4XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB5XG4gICAgICAgICAqL1xuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbCgnU2V0UmVsYXRpdmVQb3NpdGlvbicsIHt4LHl9KSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgcG9zaXRpb24uXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFBvc2l0aW9uPn0gVGhlIHdpbmRvdyBwb3NpdGlvblxuICAgICAgICAgKi9cbiAgICAgICAgUmVsYXRpdmVQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUmVsYXRpdmVQb3NpdGlvbicpOyB9LFxuXG4gICAgICAgIC8qKlxuICAgICAgICAgKiBHZXQgdGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIG9uLlxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxuICAgICAgICAgKi9cbiAgICAgICAgU2NyZWVuOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTY3JlZW4nKTsgfSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogSGlkZSB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBIaWRlOiAoKSA9PiB2b2lkIGNhbGwoJ0hpZGUnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogTWF4aW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnTWF4aW1pc2UnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogU2hvdyB0aGUgd2luZG93XG4gICAgICAgICAqL1xuICAgICAgICBTaG93OiAoKSA9PiB2b2lkIGNhbGwoJ1Nob3cnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogQ2xvc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgQ2xvc2U6ICgpID0+IHZvaWQgY2FsbCgnQ2xvc2UnKSxcblxuICAgICAgICAvKipcbiAgICAgICAgICogVG9nZ2xlIHRoZSB3aW5kb3cgbWF4aW1pc2Ugc3RhdGVcbiAgICAgICAgICovXG4gICAgICAgIFRvZ2dsZU1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ1RvZ2dsZU1heGltaXNlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFVubWF4aW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgVW5NYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdVbk1heGltaXNlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIE1pbmltaXNlIHRoZSB3aW5kb3dcbiAgICAgICAgICovXG4gICAgICAgIE1pbmltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ01pbmltaXNlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFVubWluaW1pc2UgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgVW5NaW5pbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdVbk1pbmltaXNlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFJlc3RvcmUgdGhlIHdpbmRvd1xuICAgICAgICAgKi9cbiAgICAgICAgUmVzdG9yZTogKCkgPT4gdm9pZCBjYWxsKCdSZXN0b3JlJyksXG5cbiAgICAgICAgLyoqXG4gICAgICAgICAqIFNldCB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBnIC0gQSB2YWx1ZSBiZXR3ZWVuIDAgYW5kIDI1NVxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gYiAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGEgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XG4gICAgICAgICAqL1xuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gdm9pZCBjYWxsKCdTZXRCYWNrZ3JvdW5kQ29sb3VyJywge3IsIGcsIGIsIGF9KSxcbiAgICB9O1xufVxuIiwgIi8qXG4gX1x0ICAgX19cdCAgXyBfX1xufCB8XHQgLyAvX19fIF8oXykgL19fX19cbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XG4qL1xuXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXG5cbi8qKlxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZXZlbnRzXCIpO1xuXG4vKipcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxuICpcbiAqIEBjbGFzcyBMaXN0ZW5lclxuICovXG5jbGFzcyBMaXN0ZW5lciB7XG4gICAgLyoqXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXG4gICAgICogQHBhcmFtIHtmdW5jdGlvbn0gY2FsbGJhY2tcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXG4gICAgICovXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xuICAgICAgICAvLyBEZWZhdWx0IG9mIC0xIG1lYW5zIGluZmluaXRlXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXG4gICAgICAgIC8vIFJldHVybnMgdHJ1ZSBpZiB0aGlzIGxpc3RlbmVyIHNob3VsZCBiZSBkZXN0cm95ZWRcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcbiAgICAgICAgICAgIC8vIElmIG1heENhbGxiYWNrcyBpcyBpbmZpbml0ZSwgcmV0dXJuIGZhbHNlIChkbyBub3QgZGVzdHJveSlcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XG4gICAgICAgICAgICB9XG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xuICAgICAgICAgICAgcmV0dXJuIHRoaXMubWF4Q2FsbGJhY2tzID09PSAwO1xuICAgICAgICB9O1xuICAgIH1cbn1cblxuXG4vKipcbiAqIFdhaWxzRXZlbnQgZGVmaW5lcyBhIGN1c3RvbSBldmVudC4gSXQgaXMgcGFzc2VkIHRvIGV2ZW50IGxpc3RlbmVycy5cbiAqXG4gKiBAY2xhc3MgV2FpbHNFdmVudFxuICogQHByb3BlcnR5IHtzdHJpbmd9IG5hbWUgLSBOYW1lIG9mIHRoZSBldmVudFxuICogQHByb3BlcnR5IHthbnl9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcbiAqL1xuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQge1xuICAgIC8qKlxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgV2FpbHNFdmVudC5cbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XG4gICAgICogQHBhcmFtIHthbnk9bnVsbH0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxuICAgICAqIEBtZW1iZXJvZiBXYWlsc0V2ZW50XG4gICAgICovXG4gICAgY29uc3RydWN0b3IobmFtZSwgZGF0YSA9IG51bGwpIHtcbiAgICAgICAgdGhpcy5uYW1lID0gbmFtZTtcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcbiAgICB9XG59XG5cbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXAoKTtcblxuLyoqXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgYG1heENhbGxiYWNrc2AgdGltZXMgYmVmb3JlIGJlaW5nIGRlc3Ryb3llZFxuICpcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSB7ZnVuY3Rpb24oV2FpbHNFdmVudCk6IHZvaWR9IGNhbGxiYWNrXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XG4gICAgY29uc3QgdGhpc0xpc3RlbmVyID0gbmV3IExpc3RlbmVyKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcyk7XG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xuICAgIHJldHVybiAoKSA9PiBsaXN0ZW5lck9mZih0aGlzTGlzdGVuZXIpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgLTEpO1xufVxuXG4vKipcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcbiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxuICovXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XG4gICAgcmV0dXJuIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgMSk7XG59XG5cbi8qKlxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxuICpcbiAqIEBwYXJhbSB7TGlzdGVuZXJ9IGxpc3RlbmVyXG4gKi9cbmZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XG4gICAgY29uc3QgZXZlbnROYW1lID0gbGlzdGVuZXIuZXZlbnROYW1lO1xuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XG4gICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSBlbHNlIHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50TmFtZSwgbGlzdGVuZXJzKTtcbiAgICB9XG59XG5cbi8qKlxuICogZGlzcGF0Y2hlcyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXG4gKlxuICogQGV4cG9ydFxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudFxuICovXG5leHBvcnQgZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50KSB7XG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcbiAgICBpZiAobGlzdGVuZXJzKSB7XG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBbXTtcbiAgICAgICAgbGlzdGVuZXJzLmZvckVhY2gobGlzdGVuZXIgPT4ge1xuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHtcbiAgICAgICAgICAgICAgICB0b1JlbW92ZS5wdXNoKGxpc3RlbmVyKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSk7XG4gICAgICAgIC8vIHJlbW92ZSBsaXN0ZW5lcnNcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xuICAgICAgICAgICAgaWYgKGxpc3RlbmVycy5sZW5ndGggPT09IDApIHtcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XG4gICAgICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudC5uYW1lLCBsaXN0ZW5lcnMpO1xuICAgICAgICAgICAgfVxuICAgICAgICB9XG4gICAgfVxufVxuXG4vKipcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcnMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxuICpcbiBbdjMgQ0hBTkdFXSBPZmYgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xuICpcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIE9mZihldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzKSB7XG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IHtcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XG4gICAgfSk7XG59XG5cbi8qKlxuICogT2ZmQWxsIHVucmVnaXN0ZXJzIGFsbCBsaXN0ZW5lcnNcbiAqIFt2MyBDSEFOR0VdIE9mZkFsbCBvbmx5IHVucmVnaXN0ZXJzIGxpc3RlbmVycyB3aXRoaW4gdGhlIGN1cnJlbnQgd2luZG93XG4gKlxuICovXG5leHBvcnQgZnVuY3Rpb24gT2ZmQWxsKCkge1xuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XG59XG5cbi8qKlxuICogRW1pdCBhbiBldmVudFxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudCBUaGUgZXZlbnQgdG8gZW1pdFxuICovXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkge1xuICAgIHZvaWQgY2FsbChcIkVtaXRcIiwgZXZlbnQpO1xufSIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vKipcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5NZXNzYWdlRGlhbG9nT3B0aW9uc30gTWVzc2FnZURpYWxvZ09wdGlvbnNcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5PcGVuRGlhbG9nT3B0aW9uc30gT3BlbkRpYWxvZ09wdGlvbnNcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5TYXZlRGlhbG9nT3B0aW9uc30gU2F2ZURpYWxvZ09wdGlvbnNcbiAqL1xuXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xuXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJkaWFsb2dcIik7XG5cbmxldCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwKCk7XG5cbmZ1bmN0aW9uIGdlbmVyYXRlSUQoKSB7XG4gICAgbGV0IHJlc3VsdDtcbiAgICBkbyB7XG4gICAgICAgIHJlc3VsdCA9IG5hbm9pZCgpO1xuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XG4gICAgcmV0dXJuIHJlc3VsdDtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0NhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcbiAgICBsZXQgcCA9IGRpYWxvZ1Jlc3BvbnNlcy5nZXQoaWQpO1xuICAgIGlmIChwKSB7XG4gICAgICAgIGlmIChpc0pTT04pIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShKU09OLnBhcnNlKGRhdGEpKTtcbiAgICAgICAgfSBlbHNlIHtcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcbiAgICAgICAgfVxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICB9XG59XG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XG4gICAgaWYgKHApIHtcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xuICAgIH1cbn1cblxuZnVuY3Rpb24gZGlhbG9nKHR5cGUsIG9wdGlvbnMpIHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xuICAgICAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcbiAgICAgICAgY2FsbCh0eXBlLCBvcHRpb25zKS5jYXRjaCgoZXJyb3IpID0+IHtcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcbiAgICAgICAgfSk7XG4gICAgfSk7XG59XG5cblxuLyoqXG4gKiBTaG93cyBhbiBJbmZvIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkluZm9cIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYSBXYXJuaW5nIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxuICovXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIldhcm5pbmdcIiwgb3B0aW9ucyk7XG59XG5cbi8qKlxuICogU2hvd3MgYW4gRXJyb3IgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIkVycm9yXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGEgUXVlc3Rpb24gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zKSB7XG4gICAgcmV0dXJuIGRpYWxvZyhcIlF1ZXN0aW9uXCIsIG9wdGlvbnMpO1xufVxuXG4vKipcbiAqIFNob3dzIGFuIE9wZW4gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmdbXXxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlIG9yIGFuIGFycmF5IG9mIHNlbGVjdGVkIGZpbGVzIGlmIEFsbG93c011bHRpcGxlU2VsZWN0aW9uIGlzIHRydWUuIEEgYmxhbmsgc3RyaW5nIGlzIHJldHVybmVkIGlmIG5vIGZpbGUgd2FzIHNlbGVjdGVkLlxuICovXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xuICAgIHJldHVybiBkaWFsb2coXCJPcGVuRmlsZVwiLCBvcHRpb25zKTtcbn1cblxuLyoqXG4gKiBTaG93cyBhIFNhdmUgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBSZXR1cm5zIHRoZSBzZWxlY3RlZCBmaWxlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnMpIHtcbiAgICByZXR1cm4gZGlhbG9nKFwiU2F2ZUZpbGVcIiwgb3B0aW9ucyk7XG59XG5cbiIsICJpbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcblxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY29udGV4dG1lbnVcIik7XG5cbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xuICAgIHZvaWQgY2FsbChcIk9wZW5Db250ZXh0TWVudVwiLCB7aWQsIHgsIHksIGRhdGF9KTtcbn1cblxuZXhwb3J0IGZ1bmN0aW9uIHNldHVwQ29udGV4dE1lbnVzKCkge1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdjb250ZXh0bWVudScsIGNvbnRleHRNZW51SGFuZGxlcik7XG59XG5cbmZ1bmN0aW9uIGNvbnRleHRNZW51SGFuZGxlcihldmVudCkge1xuICAgIC8vIENoZWNrIGZvciBjdXN0b20gY29udGV4dCBtZW51XG4gICAgbGV0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgbGV0IGN1c3RvbUNvbnRleHRNZW51ID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51XCIpO1xuICAgIGN1c3RvbUNvbnRleHRNZW51ID0gY3VzdG9tQ29udGV4dE1lbnUgPyBjdXN0b21Db250ZXh0TWVudS50cmltKCkgOiBcIlwiO1xuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xuICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpO1xuICAgICAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnVEYXRhID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZWxlbWVudCkuZ2V0UHJvcGVydHlWYWx1ZShcIi0tY3VzdG9tLWNvbnRleHRtZW51LWRhdGFcIik7XG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgY3VzdG9tQ29udGV4dE1lbnVEYXRhKTtcbiAgICAgICAgcmV0dXJuXG4gICAgfVxuXG4gICAgcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCk7XG59XG5cblxuLypcbi0tZGVmYXVsdC1jb250ZXh0bWVudTogYXV0bzsgKGRlZmF1bHQpIHdpbGwgc2hvdyB0aGUgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlIGlzIHRydWUgT1IgdGV4dCBoYXMgYmVlbiBzZWxlY3RlZCBPUiBlbGVtZW50IGlzIGlucHV0IG9yIHRleHRhcmVhXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IHNob3c7IHdpbGwgYWx3YXlzIHNob3cgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGU7IHdpbGwgYWx3YXlzIGhpZGUgdGhlIGRlZmF1bHQgY29udGV4dCBtZW51XG5cblRoaXMgcnVsZSBpcyBpbmhlcml0ZWQgbGlrZSBub3JtYWwgQ1NTIHJ1bGVzLCBzbyBuZXN0aW5nIHdvcmtzIGFzIGV4cGVjdGVkXG4qL1xuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCkge1xuICAgIC8vIERlYnVnIGJ1aWxkcyBhbHdheXMgc2hvdyB0aGUgbWVudVxuICAgIGlmIChERUJVRykge1xuICAgICAgICByZXR1cm47XG4gICAgfVxuXG4gICAgLy8gUHJvY2VzcyBkZWZhdWx0IGNvbnRleHQgbWVudVxuICAgIGNvbnN0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XG4gICAgY29uc3QgY29tcHV0ZWRTdHlsZSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpO1xuICAgIGNvbnN0IGRlZmF1bHRDb250ZXh0TWVudUFjdGlvbiA9IGNvbXB1dGVkU3R5bGUuZ2V0UHJvcGVydHlWYWx1ZShcIi0tZGVmYXVsdC1jb250ZXh0bWVudVwiKS50cmltKCk7XG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcbiAgICAgICAgY2FzZSBcInNob3dcIjpcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgICAgICAgICByZXR1cm47XG4gICAgICAgIGRlZmF1bHQ6XG4gICAgICAgICAgICAvLyBDaGVjayBpZiBjb250ZW50RWRpdGFibGUgaXMgdHJ1ZVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQuaXNDb250ZW50RWRpdGFibGUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm47XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWRcbiAgICAgICAgICAgIGNvbnN0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcbiAgICAgICAgICAgIGNvbnN0IGhhc1NlbGVjdGlvbiA9IChzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKVxuICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbikge1xuICAgICAgICAgICAgICAgIGZvciAobGV0IGkgPSAwOyBpIDwgc2VsZWN0aW9uLnJhbmdlQ291bnQ7IGkrKykge1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByYW5nZSA9IHNlbGVjdGlvbi5nZXRSYW5nZUF0KGkpO1xuICAgICAgICAgICAgICAgICAgICBjb25zdCByZWN0cyA9IHJhbmdlLmdldENsaWVudFJlY3RzKCk7XG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlY3QgPSByZWN0c1tqXTtcbiAgICAgICAgICAgICAgICAgICAgICAgIGlmIChkb2N1bWVudC5lbGVtZW50RnJvbVBvaW50KHJlY3QubGVmdCwgcmVjdC50b3ApID09PSBlbGVtZW50KSB7XG4gICAgICAgICAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfVxuICAgICAgICAgICAgfVxuICAgICAgICAgICAgLy8gQ2hlY2sgaWYgdGFnbmFtZSBpcyBpbnB1dCBvciB0ZXh0YXJlYVxuICAgICAgICAgICAgaWYgKGVsZW1lbnQudGFnTmFtZSA9PT0gXCJJTlBVVFwiIHx8IGVsZW1lbnQudGFnTmFtZSA9PT0gXCJURVhUQVJFQVwiKSB7XG4gICAgICAgICAgICAgICAgaWYgKGhhc1NlbGVjdGlvbiB8fCAoIWVsZW1lbnQucmVhZE9ubHkgJiYgIWVsZW1lbnQuZGlzYWJsZWQpKSB7XG4gICAgICAgICAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICB9XG5cbiAgICAgICAgICAgIC8vIGhpZGUgZGVmYXVsdCBjb250ZXh0IG1lbnVcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XG4gICAgfVxufVxuIiwgIlxuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7UXVlc3Rpb259IGZyb20gXCIuL2RpYWxvZ3NcIjtcblxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZSwgZGF0YT1udWxsKSB7XG4gICAgbGV0IGV2ZW50ID0gbmV3IFdhaWxzRXZlbnQoZXZlbnROYW1lLCBkYXRhKTtcbiAgICBFbWl0KGV2ZW50KTtcbn1cblxuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtZXZlbnRdJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcblxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgRGV0YWNoZWQ6IGZhbHNlLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcbiAgICAgICAgICAgICAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgc2VuZEV2ZW50KGV2ZW50VHlwZSk7XG4gICAgICAgIH07XG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcblxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xuXG4gICAgICAgIC8vIEFkZCBuZXcgbGlzdGVuZXJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcbiAgICB9KTtcbn1cblxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcbiAgICBpZiAod2FpbHMuV2luZG93W21ldGhvZF0gPT09IHVuZGVmaW5lZCkge1xuICAgICAgICBjb25zb2xlLmxvZyhcIldpbmRvdyBtZXRob2QgXCIgKyBtZXRob2QgKyBcIiBub3QgZm91bmRcIik7XG4gICAgfVxuICAgIHdhaWxzLldpbmRvd1ttZXRob2RdKCk7XG59XG5cbmZ1bmN0aW9uIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpIHtcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC13aW5kb3ddJyk7XG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xuICAgICAgICBjb25zdCB3aW5kb3dNZXRob2QgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtd2luZG93Jyk7XG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XG5cbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xuICAgICAgICAgICAgICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICAgICAgICAgICAgICB9XG4gICAgICAgICAgICAgICAgfSk7XG4gICAgICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICAgICAgfVxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xuICAgICAgICB9O1xuXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcblxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XG4gICAgfSk7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiByZWxvYWRXTUwoKSB7XG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcbiAgICBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG4vLyBkZWZpbmVkIGluIHRoZSBUYXNrZmlsZVxuZXhwb3J0IGxldCBpbnZva2UgPSBmdW5jdGlvbihpbnB1dCkge1xuICAgIGlmKFdJTkRPV1MpIHtcbiAgICAgICAgY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UoaW5wdXQpO1xuICAgIH0gZWxzZSB7XG4gICAgICAgIHdlYmtpdC5tZXNzYWdlSGFuZGxlcnMuZXh0ZXJuYWwucG9zdE1lc3NhZ2UoaW5wdXQpO1xuICAgIH1cbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5sZXQgZmxhZ3MgPSBuZXcgTWFwKCk7XG5cbmZ1bmN0aW9uIGNvbnZlcnRUb01hcChvYmopIHtcbiAgICBjb25zdCBtYXAgPSBuZXcgTWFwKCk7XG5cbiAgICBmb3IgKGNvbnN0IFtrZXksIHZhbHVlXSBvZiBPYmplY3QuZW50cmllcyhvYmopKSB7XG4gICAgICAgIGlmICh0eXBlb2YgdmFsdWUgPT09ICdvYmplY3QnICYmIHZhbHVlICE9PSBudWxsKSB7XG4gICAgICAgICAgICBtYXAuc2V0KGtleSwgY29udmVydFRvTWFwKHZhbHVlKSk7IC8vIFJlY3Vyc2l2ZWx5IGNvbnZlcnQgbmVzdGVkIG9iamVjdFxuICAgICAgICB9IGVsc2Uge1xuICAgICAgICAgICAgbWFwLnNldChrZXksIHZhbHVlKTtcbiAgICAgICAgfVxuICAgIH1cblxuICAgIHJldHVybiBtYXA7XG59XG5cbmZldGNoKFwiL3dhaWxzL2ZsYWdzXCIpLnRoZW4oKHJlc3BvbnNlKSA9PiB7XG4gICAgcmVzcG9uc2UuanNvbigpLnRoZW4oKGRhdGEpID0+IHtcbiAgICAgICAgZmxhZ3MgPSBjb252ZXJ0VG9NYXAoZGF0YSk7XG4gICAgfSk7XG59KTtcblxuXG5mdW5jdGlvbiBnZXRWYWx1ZUZyb21NYXAoa2V5U3RyaW5nKSB7XG4gICAgY29uc3Qga2V5cyA9IGtleVN0cmluZy5zcGxpdCgnLicpO1xuICAgIGxldCB2YWx1ZSA9IGZsYWdzO1xuXG4gICAgZm9yIChjb25zdCBrZXkgb2Yga2V5cykge1xuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBNYXApIHtcbiAgICAgICAgICAgIHZhbHVlID0gdmFsdWUuZ2V0KGtleSk7XG4gICAgICAgIH0gZWxzZSB7XG4gICAgICAgICAgICB2YWx1ZSA9IHZhbHVlW2tleV07XG4gICAgICAgIH1cblxuICAgICAgICBpZiAodmFsdWUgPT09IHVuZGVmaW5lZCkge1xuICAgICAgICAgICAgYnJlYWs7XG4gICAgICAgIH1cbiAgICB9XG5cbiAgICByZXR1cm4gdmFsdWU7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBHZXRGbGFnKGtleVN0cmluZykge1xuICAgIHJldHVybiBnZXRWYWx1ZUZyb21NYXAoa2V5U3RyaW5nKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cblxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xuXG5pbXBvcnQge2ludm9rZX0gZnJvbSBcIi4vaW52b2tlXCI7XG5pbXBvcnQge0dldEZsYWd9IGZyb20gXCIuL2ZsYWdzXCI7XG5cbmxldCBzaG91bGREcmFnID0gZmFsc2U7XG5cbmV4cG9ydCBmdW5jdGlvbiBkcmFnVGVzdChlKSB7XG4gICAgbGV0IHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS13ZWJraXQtYXBwLXJlZ2lvblwiKTtcbiAgICBpZiAodmFsKSB7XG4gICAgICAgIHZhbCA9IHZhbC50cmltKCk7XG4gICAgfVxuXG4gICAgaWYgKHZhbCAhPT0gXCJkcmFnXCIpIHtcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xuICAgIH1cblxuICAgIC8vIE9ubHkgcHJvY2VzcyB0aGUgcHJpbWFyeSBidXR0b25cbiAgICBpZiAoZS5idXR0b25zICE9PSAxKSB7XG4gICAgICAgIHJldHVybiBmYWxzZTtcbiAgICB9XG5cbiAgICByZXR1cm4gZS5kZXRhaWwgPT09IDE7XG59XG5cbmV4cG9ydCBmdW5jdGlvbiBzZXR1cERyYWcoKSB7XG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIG9uTW91c2VEb3duKTtcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vtb3ZlJywgb25Nb3VzZU1vdmUpO1xuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgb25Nb3VzZVVwKTtcbn1cblxubGV0IHJlc2l6ZUVkZ2UgPSBudWxsO1xuXG5mdW5jdGlvbiB0ZXN0UmVzaXplKGUpIHtcbiAgICBpZiggcmVzaXplRWRnZSApIHtcbiAgICAgICAgaW52b2tlKFwicmVzaXplOlwiICsgcmVzaXplRWRnZSk7XG4gICAgICAgIHJldHVybiB0cnVlXG4gICAgfVxuICAgIHJldHVybiBmYWxzZTtcbn1cblxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xuXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nIG9uIFdpbmRvd3NcbiAgICBpZiggV0lORE9XUyApIHtcbiAgICAgICAgaWYgKHRlc3RSZXNpemUoKSkge1xuICAgICAgICAgICAgcmV0dXJuO1xuICAgICAgICB9XG4gICAgfVxuICAgIGlmIChkcmFnVGVzdChlKSkge1xuICAgICAgICAvLyBJZ25vcmUgZHJhZyBvbiBzY3JvbGxiYXJzXG4gICAgICAgIGlmIChlLm9mZnNldFggPiBlLnRhcmdldC5jbGllbnRXaWR0aCB8fCBlLm9mZnNldFkgPiBlLnRhcmdldC5jbGllbnRIZWlnaHQpIHtcbiAgICAgICAgICAgIHJldHVybjtcbiAgICAgICAgfVxuICAgICAgICBzaG91bGREcmFnID0gdHJ1ZTtcbiAgICB9IGVsc2Uge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgfVxufVxuXG5mdW5jdGlvbiBvbk1vdXNlVXAoZSkge1xuICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcbiAgICAgICAgZW5kRHJhZygpO1xuICAgIH1cbn1cblxuZXhwb3J0IGZ1bmN0aW9uIGVuZERyYWcoKSB7XG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XG4gICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xufVxuXG5mdW5jdGlvbiBzZXRSZXNpemUoY3Vyc29yKSB7XG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCBkZWZhdWx0Q3Vyc29yO1xuICAgIHJlc2l6ZUVkZ2UgPSBjdXJzb3I7XG59XG5cbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGUpIHtcbiAgICBpZiAoc2hvdWxkRHJhZykge1xuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XG4gICAgICAgIGlmIChtb3VzZVByZXNzZWQgPiAwKSB7XG4gICAgICAgICAgICBpbnZva2UoXCJkcmFnXCIpO1xuICAgICAgICB9XG4gICAgICAgIHJldHVybjtcbiAgICB9XG5cbiAgICBpZiAoV0lORE9XUykge1xuICAgICAgICBoYW5kbGVSZXNpemUoZSk7XG4gICAgfVxufVxuXG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xuXG5mdW5jdGlvbiBoYW5kbGVSZXNpemUoZSkge1xuICAgIGxldCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xuICAgIGxldCByZXNpemVIYW5kbGVXaWR0aCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlV2lkdGhcIikgfHwgNTtcblxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhc1xuICAgIGxldCBjb3JuZXJFeHRyYSA9IEdldEZsYWcoXCJyZXNpemVDb3JuZXJFeHRyYVwiKSB8fCAzO1xuXG4gICAgbGV0IHJpZ2h0Qm9yZGVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcbiAgICBsZXQgbGVmdEJvcmRlciA9IGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcblxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyc1xuICAgIGxldCByaWdodENvcm5lciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGxldCBsZWZ0Q29ybmVyID0gZS5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xuICAgIGxldCB0b3BDb3JuZXIgPSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuICAgIGxldCBib3R0b21Db3JuZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xuXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiByZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcbiAgICAgICAgc2V0UmVzaXplKCk7XG4gICAgfVxuICAgIC8vIEFkanVzdGVkIGZvciBjb3JuZXIgYXJlYXNcbiAgICBlbHNlIGlmIChyaWdodENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInNlLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQ29ybmVyICYmIHJpZ2h0Q29ybmVyKSBzZXRSZXNpemUoXCJuZS1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcbiAgICBlbHNlIGlmIChib3R0b21Cb3JkZXIpIHNldFJlc2l6ZShcInMtcmVzaXplXCIpO1xuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcbn1cbiIsICIvKlxuIF9cdCAgIF9fXHQgIF8gX19cbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxuKi9cbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cblxuXG5pbXBvcnQgKiBhcyBDbGlwYm9hcmQgZnJvbSAnLi9jbGlwYm9hcmQnO1xuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSAnLi9hcHBsaWNhdGlvbic7XG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xuaW1wb3J0ICogYXMgU2NyZWVucyBmcm9tICcuL3NjcmVlbnMnO1xuaW1wb3J0IHtQbHVnaW4sIENhbGwsIGNhbGxFcnJvckNhbGxiYWNrLCBjYWxsQ2FsbGJhY2t9IGZyb20gXCIuL2NhbGxzXCI7XG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XG5pbXBvcnQge2Rpc3BhdGNoV2FpbHNFdmVudCwgRW1pdCwgT2ZmLCBPZmZBbGwsIE9uLCBPbmNlLCBPbk11bHRpcGxlfSBmcm9tIFwiLi9ldmVudHNcIjtcbmltcG9ydCB7ZGlhbG9nQ2FsbGJhY2ssIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssIEVycm9yLCBJbmZvLCBPcGVuRmlsZSwgUXVlc3Rpb24sIFNhdmVGaWxlLCBXYXJuaW5nLH0gZnJvbSBcIi4vZGlhbG9nc1wiO1xuaW1wb3J0IHtzZXR1cENvbnRleHRNZW51c30gZnJvbSBcIi4vY29udGV4dG1lbnVcIjtcbmltcG9ydCB7cmVsb2FkV01MfSBmcm9tIFwiLi93bWxcIjtcbmltcG9ydCB7c2V0dXBEcmFnLCBlbmREcmFnfSBmcm9tIFwiLi9kcmFnXCI7XG5cbndpbmRvdy53YWlscyA9IHtcbiAgICAuLi5uZXdSdW50aW1lKG51bGwpLFxuICAgIENhcGFiaWxpdGllczoge30sXG59O1xuXG5mZXRjaChcIi93YWlscy9jYXBhYmlsaXRpZXNcIikudGhlbigocmVzcG9uc2UpID0+IHtcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xuICAgICAgICB3aW5kb3cud2FpbHMuQ2FwYWJpbGl0aWVzID0gZGF0YTtcbiAgICB9KTtcbn0pO1xuXG4vLyBJbnRlcm5hbCB3YWlscyBlbmRwb2ludHNcbndpbmRvdy5fd2FpbHMgPSB7XG4gICAgZGlhbG9nQ2FsbGJhY2ssXG4gICAgZGlhbG9nRXJyb3JDYWxsYmFjayxcbiAgICBkaXNwYXRjaFdhaWxzRXZlbnQsXG4gICAgY2FsbENhbGxiYWNrLFxuICAgIGNhbGxFcnJvckNhbGxiYWNrLFxuICAgIGVuZERyYWcsXG59O1xuXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZSh3aW5kb3dOYW1lKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgQ2xpcGJvYXJkOiB7XG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcbiAgICAgICAgfSxcbiAgICAgICAgQXBwbGljYXRpb246IHtcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uLFxuICAgICAgICAgICAgR2V0V2luZG93QnlOYW1lKHdpbmRvd05hbWUpIHtcbiAgICAgICAgICAgICAgICByZXR1cm4gbmV3UnVudGltZSh3aW5kb3dOYW1lKTtcbiAgICAgICAgICAgIH1cbiAgICAgICAgfSxcbiAgICAgICAgTG9nLFxuICAgICAgICBTY3JlZW5zLFxuICAgICAgICBDYWxsLFxuICAgICAgICBQbHVnaW4sXG4gICAgICAgIFdNTDoge1xuICAgICAgICAgICAgUmVsb2FkOiByZWxvYWRXTUwsXG4gICAgICAgIH0sXG4gICAgICAgIERpYWxvZzoge1xuICAgICAgICAgICAgSW5mbyxcbiAgICAgICAgICAgIFdhcm5pbmcsXG4gICAgICAgICAgICBFcnJvcixcbiAgICAgICAgICAgIFF1ZXN0aW9uLFxuICAgICAgICAgICAgT3BlbkZpbGUsXG4gICAgICAgICAgICBTYXZlRmlsZSxcbiAgICAgICAgfSxcbiAgICAgICAgRXZlbnRzOiB7XG4gICAgICAgICAgICBFbWl0LFxuICAgICAgICAgICAgT24sXG4gICAgICAgICAgICBPbmNlLFxuICAgICAgICAgICAgT25NdWx0aXBsZSxcbiAgICAgICAgICAgIE9mZixcbiAgICAgICAgICAgIE9mZkFsbCxcbiAgICAgICAgfSxcbiAgICAgICAgV2luZG93OiBuZXdXaW5kb3cod2luZG93TmFtZSksXG4gICAgfTtcbn1cblxuaWYgKERFQlVHKSB7XG4gICAgY29uc29sZS5sb2coXCJXYWlscyB2My4wLjAgRGVidWcgTW9kZSBFbmFibGVkXCIpO1xufVxuXG5zZXR1cENvbnRleHRNZW51cygpO1xuc2V0dXBEcmFnKCk7XG5cbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoXCJET01Db250ZW50TG9hZGVkXCIsIGZ1bmN0aW9uKGV2ZW50KSB7XG4gICAgcmVsb2FkV01MKCk7XG59KTsiXSwKICAibWFwcGluZ3MiOiAiOzs7Ozs7OztBQUFBO0FBQUE7QUFBQTtBQUFBO0FBQUE7OztBQ1lBLE1BQU0sYUFBYSxPQUFPLFNBQVMsU0FBUztBQUU1QyxXQUFTLFlBQVksUUFBUSxZQUFZLE1BQU07QUFDM0MsUUFBSSxNQUFNLElBQUksSUFBSSxVQUFVO0FBQzVCLFFBQUksYUFBYSxPQUFPLFVBQVUsTUFBTTtBQUN4QyxRQUFJLE1BQU07QUFDTixVQUFJLGFBQWEsT0FBTyxRQUFRLEtBQUssVUFBVSxJQUFJLENBQUM7QUFBQSxJQUN4RDtBQUNBLFFBQUksZUFBZTtBQUFBLE1BQ2YsU0FBUyxDQUFDO0FBQUEsSUFDZDtBQUNBLFFBQUksWUFBWTtBQUNaLG1CQUFhLFFBQVEscUJBQXFCLElBQUk7QUFBQSxJQUNsRDtBQUNBLFdBQU8sSUFBSSxRQUFRLENBQUMsU0FBUyxXQUFXO0FBQ3BDLFlBQU0sS0FBSyxZQUFZLEVBQ2xCLEtBQUssY0FBWTtBQUNkLFlBQUksU0FBUyxJQUFJO0FBRWIsY0FBSSxTQUFTLFFBQVEsSUFBSSxjQUFjLEtBQUssU0FBUyxRQUFRLElBQUksY0FBYyxFQUFFLFFBQVEsa0JBQWtCLE1BQU0sSUFBSTtBQUNqSCxtQkFBTyxTQUFTLEtBQUs7QUFBQSxVQUN6QixPQUFPO0FBQ0gsbUJBQU8sU0FBUyxLQUFLO0FBQUEsVUFDekI7QUFBQSxRQUNKO0FBQ0EsZUFBTyxNQUFNLFNBQVMsVUFBVSxDQUFDO0FBQUEsTUFDckMsQ0FBQyxFQUNBLEtBQUssVUFBUSxRQUFRLElBQUksQ0FBQyxFQUMxQixNQUFNLFdBQVMsT0FBTyxLQUFLLENBQUM7QUFBQSxJQUNyQyxDQUFD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QWM7QUFBVyxnQkFBVSxXQUFXO0FBQUEsYUFDOUMsYUFBYTtBQUFhLGdCQUFVLFdBQVc7QUFBQSxhQUMvQztBQUFZLGdCQUFVLFVBQVU7QUFBQSxhQUNoQztBQUFXLGdCQUFVLFVBQVU7QUFBQSxhQUMvQjtBQUFjLGdCQUFVLFVBQVU7QUFBQSxhQUNsQztBQUFhLGdCQUFVLFVBQVU7QUFBQSxFQUM5Qzs7O0FDL0dBLFNBQU8sUUFBUTtBQUFBLElBQ1gsR0FBRyxXQUFXLElBQUk7QUFBQSxJQUNsQixjQUFjLENBQUM7QUFBQSxFQUNuQjtBQUVBLFFBQU0scUJBQXFCLEVBQUUsS0FBSyxDQUFDLGFBQWE7QUFDNUMsYUFBUyxLQUFLLEVBQUUsS0FBSyxDQUFDLFNBQVM7QUFDM0IsYUFBTyxNQUFNLGVBQWU7QUFBQSxJQUNoQyxDQUFDO0FBQUEsRUFDTCxDQUFDO0FBR0QsU0FBTyxTQUFTO0FBQUEsSUFDWjtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsRUFDSjtBQUVPLFdBQVMsV0FBVyxZQUFZO0FBQ25DLFdBQU87QUFBQSxNQUNILFdBQVc7QUFBQSxRQUNQLEdBQUc7QUFBQSxNQUNQO0FBQUEsTUFDQSxhQUFhO0FBQUEsUUFDVCxHQUFHO0FBQUEsUUFDSCxnQkFBZ0JDLGFBQVk7QUFDeEIsaUJBQU8sV0FBV0EsV0FBVTtBQUFBLFFBQ2hDO0FBQUEsTUFDSjtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBLEtBQUs7QUFBQSxRQUNELFFBQVE7QUFBQSxNQUNaO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBLE9BQUFDO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVEsVUFBVSxVQUFVO0FBQUEsSUFDaEM7QUFBQSxFQUNKO0FBRUEsTUFBSSxNQUFPO0FBQ1AsWUFBUSxJQUFJLGlDQUFpQztBQUFBLEVBQ2pEO0FBRUEsb0JBQWtCO0FBQ2xCLFlBQVU7QUFFVixXQUFTLGlCQUFpQixvQkFBb0IsU0FBUyxPQUFPO0FBQzFELGNBQVU7QUFBQSxFQUNkLENBQUM7IiwKICAibmFtZXMiOiBbImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiRXJyb3IiLCAiY2FsbCIsICJ3aW5kb3dOYW1lIiwgIkVycm9yIl0KfQo=
