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
       * Set the window position.
       * @param {number} x
       * @param {number} y
       */
      SetPosition: (x, y) => call9("SetPosition", { x, y }),
      /**
       * Get the window position.
       * @returns {Promise<Position>} The window position
       */
      Position: () => {
        return call9("Position");
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
    console.log("dispatching event: ", { event });
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
    return call8("OpenContextMenu", { id, x, y, data });
  }
  function enableContextMenus(enabled) {
    if (enabled) {
      window.addEventListener("contextmenu", contextMenuHandler);
    } else {
      window.removeEventListener("contextmenu", contextMenuHandler);
    }
  }
  function contextMenuHandler(event) {
    processContextMenu(event.target, event);
  }
  function processContextMenu(element, event) {
    let id = element.getAttribute("data-contextmenu");
    if (id) {
      event.preventDefault();
      openContextMenu(id, event.clientX, event.clientY, element.getAttribute("data-contextmenu-data"));
    } else {
      let parent = element.parentElement;
      if (parent) {
        processContextMenu(parent, event);
      }
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
          Question({ Title: "Confirm", Message: confirm, Buttons: [{ Label: "Yes" }, { Label: "No", IsDefault: true }] }).then(function(result) {
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
    let val = window.getComputedStyle(e.target).getPropertyValue("app-region");
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
  enableContextMenus(true);
  setupDrag();
  document.addEventListener("DOMContentLoaded", function(event) {
    reloadWML();
  });
})();
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvaW52b2tlLmpzIiwgImRlc2t0b3AvZmxhZ3MuanMiLCAiZGVza3RvcC9kcmFnLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjbGlwYm9hcmRcIik7XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xyXG4gICAgdm9pZCBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIlRleHRcIik7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuY29uc3QgcnVudGltZVVSTCA9IHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XHJcblxyXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIGlmIChhcmdzKSB7XHJcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcclxuICAgIH1cclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0LCB3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJhcHBsaWNhdGlvblwiKTtcclxuXHJcbi8qKlxyXG4gKiBIaWRlIHRoZSBhcHBsaWNhdGlvblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XHJcbiAgICB2b2lkIGNhbGwoXCJIaWRlXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xyXG4gICAgdm9pZCBjYWxsKFwiU2hvd1wiKTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBRdWl0IHRoZSBhcHBsaWNhdGlvblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XHJcbiAgICB2b2lkIGNhbGwoXCJRdWl0XCIpO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwibG9nXCIpO1xyXG5cclxuLyoqXHJcbiAqIExvZ3MgYSBtZXNzYWdlLlxyXG4gKiBAcGFyYW0ge21lc3NhZ2V9IE1lc3NhZ2UgdG8gbG9nXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nKG1lc3NhZ2UpIHtcclxuICAgIHJldHVybiBjYWxsKFwiTG9nXCIsIG1lc3NhZ2UpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJzY3JlZW5zXCIpO1xyXG5cclxuLyoqXHJcbiAqIEdldHMgYWxsIHNjcmVlbnMuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIkdldEFsbFwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIkdldFByaW1hcnlcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XHJcbiAqIEBjb25zdHJ1Y3RvclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIkdldEN1cnJlbnRcIik7XHJcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY2FsbFwiKTtcclxuXHJcbmxldCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBjYWxsQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gY2FsbEVycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcclxuICAgIGxldCBwID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgaWYgKHApIHtcclxuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xyXG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcclxuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xyXG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xyXG4gICAgICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGwgYSBwbHVnaW4gbWV0aG9kXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gbmFtZSBvZiB0aGUgcGx1Z2luXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gbmFtZSBvZiB0aGUgbWV0aG9kXHJcbiAqIEBwYXJhbSB7Li4uYW55fSBhcmdzIC0gYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZFxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxhbnk+fSAtIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSByZXN1bHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxyXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXHJcbiAgICAgICAgbWV0aG9kTmFtZTogbWV0aG9kTmFtZSxcclxuICAgICAgICBhcmdzOiBhcmdzLFxyXG4gICAgfSk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2l6ZX0gU2l6ZVxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlBvc2l0aW9ufSBQb3NpdGlvblxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gbmV3V2luZG93KHdpbmRvd05hbWUpIHtcclxuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCB3aW5kb3dOYW1lKTtcclxuICAgIHJldHVybiB7XHJcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxyXG4gICAgICAgIC8vIFJlbG9hZEFwcDogKCkgPT4gY2FsbCgnV1InKSxcclxuICAgICAgICAvLyBTZXRTeXN0ZW1EZWZhdWx0VGhlbWU6ICgpID0+IGNhbGwoJ1dBU0RUJyksXHJcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxyXG4gICAgICAgIC8vIFNldERhcmtUaGVtZTogKCkgPT4gY2FsbCgnV0FEVCcpLFxyXG4gICAgICAgIC8vIElzRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnV0lGJyksXHJcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxyXG4gICAgICAgIC8vIElzTWluaW1pemVkOiAoKSA9PiBjYWxsKCdXSU1OJyksXHJcbiAgICAgICAgLy8gSXNXaW5kb3dlZDogKCkgPT4gY2FsbCgnV0lGJyksXHJcblxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cuXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgQ2VudGVyOiAoKSA9PiB2b2lkIGNhbGwoJ0NlbnRlcicpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyB0aXRsZS5cclxuICAgICAgICAgKiBAcGFyYW0gdGl0bGVcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRUaXRsZTogKHRpdGxlKSA9PiB2b2lkIGNhbGwoJ1NldFRpdGxlJywge3RpdGxlfSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIE1ha2VzIHRoZSB3aW5kb3cgZnVsbHNjcmVlbi5cclxuICAgICAgICAgKi9cclxuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ0Z1bGxzY3JlZW4nKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogVW5mdWxsc2NyZWVuIHRoZSB3aW5kb3cuXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBzaXplLlxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCBUaGUgd2luZG93IHdpZHRoXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCBUaGUgd2luZG93IGhlaWdodFxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRTaXplJywge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBzaXplLlxyXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNpemU+fSBUaGUgd2luZG93IHNpemVcclxuICAgICAgICAgKi9cclxuICAgICAgICBTaXplOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTaXplJyk7IH0sXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1heGltdW0gc2l6ZS5cclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IHZvaWQgY2FsbCgnU2V0TWF4U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgbWluaW11bSBzaXplLlxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gdm9pZCBjYWxsKCdTZXRNaW5TaXplJywge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXHJcbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSBvblRvcCBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGFsd2F5cyBvbiB0b3BcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRBbHdheXNPblRvcDogKG9uVG9wKSA9PiB2b2lkIGNhbGwoJ1NldEFsd2F5c09uVG9wJywge2Fsd2F5c09uVG9wOm9uVG9wfSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHBvc2l0aW9uLlxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB4XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHlcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRQb3NpdGlvbjogKHgsIHkpID0+IGNhbGwoJ1NldFBvc2l0aW9uJywge3gseX0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBwb3NpdGlvbi5cclxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxQb3NpdGlvbj59IFRoZSB3aW5kb3cgcG9zaXRpb25cclxuICAgICAgICAgKi9cclxuICAgICAgICBQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUG9zaXRpb24nKTsgfSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogR2V0IHRoZSBzY3JlZW4gdGhlIHdpbmRvdyBpcyBvbi5cclxuICAgICAgICAgKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNjcmVlbjogKCkgPT4geyByZXR1cm4gY2FsbCgnU2NyZWVuJyk7IH0sXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIEhpZGUgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIEhpZGU6ICgpID0+IHZvaWQgY2FsbCgnSGlkZScpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBNYXhpbWlzZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnTWF4aW1pc2UnKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogU2hvdyB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2hvdzogKCkgPT4gdm9pZCBjYWxsKCdTaG93JyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIENsb3NlIHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBDbG9zZTogKCkgPT4gdm9pZCBjYWxsKCdDbG9zZScpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBUb2dnbGUgdGhlIHdpbmRvdyBtYXhpbWlzZSBzdGF0ZVxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFRvZ2dsZU1heGltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ1RvZ2dsZU1heGltaXNlJyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFVubWF4aW1pc2UgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFVuTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVW5NYXhpbWlzZScpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBNaW5pbWlzZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgTWluaW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnTWluaW1pc2UnKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogVW5taW5pbWlzZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgVW5NaW5pbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdVbk1pbmltaXNlJyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFJlc3RvcmUgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFJlc3RvcmU6ICgpID0+IHZvaWQgY2FsbCgnUmVzdG9yZScpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIGJhY2tncm91bmQgY29sb3VyIG9mIHRoZSB3aW5kb3cuXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGcgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGIgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGEgLSBBIHZhbHVlIGJldHdlZW4gMCBhbmQgMjU1XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0QmFja2dyb3VuZENvbG91cjogKHIsIGcsIGIsIGEpID0+IHZvaWQgY2FsbCgnU2V0QmFja2dyb3VuZENvbG91cicsIHtyLCBnLCBiLCBhfSksXHJcbiAgICB9O1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLldhaWxzRXZlbnR9IFdhaWxzRXZlbnRcclxuICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImV2ZW50c1wiKTtcclxuXHJcbi8qKlxyXG4gKiBUaGUgTGlzdGVuZXIgY2xhc3MgZGVmaW5lcyBhIGxpc3RlbmVyISA6LSlcclxuICpcclxuICogQGNsYXNzIExpc3RlbmVyXHJcbiAqL1xyXG5jbGFzcyBMaXN0ZW5lciB7XHJcbiAgICAvKipcclxuICAgICAqIENyZWF0ZXMgYW4gaW5zdGFuY2Ugb2YgTGlzdGVuZXIuXHJcbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAgICAgKiBAcGFyYW0ge2Z1bmN0aW9ufSBjYWxsYmFja1xyXG4gICAgICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xyXG4gICAgICogQG1lbWJlcm9mIExpc3RlbmVyXHJcbiAgICAgKi9cclxuICAgIGNvbnN0cnVjdG9yKGV2ZW50TmFtZSwgY2FsbGJhY2ssIG1heENhbGxiYWNrcykge1xyXG4gICAgICAgIHRoaXMuZXZlbnROYW1lID0gZXZlbnROYW1lO1xyXG4gICAgICAgIC8vIERlZmF1bHQgb2YgLTEgbWVhbnMgaW5maW5pdGVcclxuICAgICAgICB0aGlzLm1heENhbGxiYWNrcyA9IG1heENhbGxiYWNrcyB8fCAtMTtcclxuICAgICAgICAvLyBDYWxsYmFjayBpbnZva2VzIHRoZSBjYWxsYmFjayB3aXRoIHRoZSBnaXZlbiBkYXRhXHJcbiAgICAgICAgLy8gUmV0dXJucyB0cnVlIGlmIHRoaXMgbGlzdGVuZXIgc2hvdWxkIGJlIGRlc3Ryb3llZFxyXG4gICAgICAgIHRoaXMuQ2FsbGJhY2sgPSAoZGF0YSkgPT4ge1xyXG4gICAgICAgICAgICBjYWxsYmFjayhkYXRhKTtcclxuICAgICAgICAgICAgLy8gSWYgbWF4Q2FsbGJhY2tzIGlzIGluZmluaXRlLCByZXR1cm4gZmFsc2UgKGRvIG5vdCBkZXN0cm95KVxyXG4gICAgICAgICAgICBpZiAodGhpcy5tYXhDYWxsYmFja3MgPT09IC0xKSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgLy8gRGVjcmVtZW50IG1heENhbGxiYWNrcy4gUmV0dXJuIHRydWUgaWYgbm93IDAsIG90aGVyd2lzZSBmYWxzZVxyXG4gICAgICAgICAgICB0aGlzLm1heENhbGxiYWNrcyAtPSAxO1xyXG4gICAgICAgICAgICByZXR1cm4gdGhpcy5tYXhDYWxsYmFja3MgPT09IDA7XHJcbiAgICAgICAgfTtcclxuICAgIH1cclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBXYWlsc0V2ZW50IGRlZmluZXMgYSBjdXN0b20gZXZlbnQuIEl0IGlzIHBhc3NlZCB0byBldmVudCBsaXN0ZW5lcnMuXHJcbiAqXHJcbiAqIEBjbGFzcyBXYWlsc0V2ZW50XHJcbiAqIEBwcm9wZXJ0eSB7c3RyaW5nfSBuYW1lIC0gTmFtZSBvZiB0aGUgZXZlbnRcclxuICogQHByb3BlcnR5IHthbnl9IGRhdGEgLSBEYXRhIGFzc29jaWF0ZWQgd2l0aCB0aGUgZXZlbnRcclxuICovXHJcbmV4cG9ydCBjbGFzcyBXYWlsc0V2ZW50IHtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBXYWlsc0V2ZW50LlxyXG4gICAgICogQHBhcmFtIHtzdHJpbmd9IG5hbWUgLSBOYW1lIG9mIHRoZSBldmVudFxyXG4gICAgICogQHBhcmFtIHthbnk9bnVsbH0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxyXG4gICAgICogQG1lbWJlcm9mIFdhaWxzRXZlbnRcclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IobmFtZSwgZGF0YSA9IG51bGwpIHtcclxuICAgICAgICB0aGlzLm5hbWUgPSBuYW1lO1xyXG4gICAgICAgIHRoaXMuZGF0YSA9IGRhdGE7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBjb25zdCBldmVudExpc3RlbmVycyA9IG5ldyBNYXAoKTtcclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgYG1heENhbGxiYWNrc2AgdGltZXMgYmVmb3JlIGJlaW5nIGRlc3Ryb3llZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcclxuICogQHBhcmFtIHtudW1iZXJ9IG1heENhbGxiYWNrc1xyXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uTXVsdGlwbGUoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50TmFtZSkgfHwgW107XHJcbiAgICBjb25zdCB0aGlzTGlzdGVuZXIgPSBuZXcgTGlzdGVuZXIoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKTtcclxuICAgIGxpc3RlbmVycy5wdXNoKHRoaXNMaXN0ZW5lcik7XHJcbiAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xyXG4gICAgcmV0dXJuICgpID0+IGxpc3RlbmVyT2ZmKHRoaXNMaXN0ZW5lcik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgZXZlcnkgdGltZSB0aGUgZXZlbnQgaXMgZW1pdHRlZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcclxuICogQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbihldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAtMSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBSZWdpc3RlcnMgYW4gZXZlbnQgbGlzdGVuZXIgdGhhdCB3aWxsIGJlIGludm9rZWQgb25jZSB0aGVuIGRlc3Ryb3llZFxyXG4gKlxyXG4gKiBAZXhwb3J0XHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICogQHBhcmFtIHtmdW5jdGlvbihXYWlsc0V2ZW50KTogdm9pZH0gY2FsbGJhY2tcclxuIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25jZShldmVudE5hbWUsIGNhbGxiYWNrKSB7XHJcbiAgICByZXR1cm4gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCAxKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIGxpc3RlbmVyT2ZmIHVucmVnaXN0ZXJzIGEgbGlzdGVuZXIgcHJldmlvdXNseSByZWdpc3RlcmVkIHdpdGggT25cclxuICpcclxuICogQHBhcmFtIHtMaXN0ZW5lcn0gbGlzdGVuZXJcclxuICovXHJcbmZ1bmN0aW9uIGxpc3RlbmVyT2ZmKGxpc3RlbmVyKSB7XHJcbiAgICBjb25zdCBldmVudE5hbWUgPSBsaXN0ZW5lci5ldmVudE5hbWU7XHJcbiAgICAvLyBSZW1vdmUgbG9jYWwgbGlzdGVuZXJcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKS5maWx0ZXIobCA9PiBsICE9PSBsaXN0ZW5lcik7XHJcbiAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnROYW1lLCBsaXN0ZW5lcnMpO1xyXG4gICAgfVxyXG59XHJcblxyXG4vKipcclxuICogZGlzcGF0Y2hlcyBhbiBldmVudCB0byBhbGwgbGlzdGVuZXJzXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtXYWlsc0V2ZW50fSBldmVudFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIGRpc3BhdGNoV2FpbHNFdmVudChldmVudCkge1xyXG4gICAgY29uc29sZS5sb2coXCJkaXNwYXRjaGluZyBldmVudDogXCIsIHtldmVudH0pO1xyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudC5uYW1lKTtcclxuICAgIGlmIChsaXN0ZW5lcnMpIHtcclxuICAgICAgICAvLyBpdGVyYXRlIGxpc3RlbmVycyBhbmQgY2FsbCBjYWxsYmFjay4gSWYgY2FsbGJhY2sgcmV0dXJucyB0cnVlLCByZW1vdmUgbGlzdGVuZXJcclxuICAgICAgICBsZXQgdG9SZW1vdmUgPSBbXTtcclxuICAgICAgICBsaXN0ZW5lcnMuZm9yRWFjaChsaXN0ZW5lciA9PiB7XHJcbiAgICAgICAgICAgIGxldCByZW1vdmUgPSBsaXN0ZW5lci5DYWxsYmFjayhldmVudCk7XHJcbiAgICAgICAgICAgIGlmIChyZW1vdmUpIHtcclxuICAgICAgICAgICAgICAgIHRvUmVtb3ZlLnB1c2gobGlzdGVuZXIpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSk7XHJcbiAgICAgICAgLy8gcmVtb3ZlIGxpc3RlbmVyc1xyXG4gICAgICAgIGlmICh0b1JlbW92ZS5sZW5ndGggPiAwKSB7XHJcbiAgICAgICAgICAgIGxpc3RlbmVycyA9IGxpc3RlbmVycy5maWx0ZXIobCA9PiAhdG9SZW1vdmUuaW5jbHVkZXMobCkpO1xyXG4gICAgICAgICAgICBpZiAobGlzdGVuZXJzLmxlbmd0aCA9PT0gMCkge1xyXG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50Lm5hbWUpO1xyXG4gICAgICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICAgICAgZXZlbnRMaXN0ZW5lcnMuc2V0KGV2ZW50Lm5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPbixcclxuICogb3B0aW9uYWxseSBtdWx0aXBsZSBsaXN0ZW5lcnMgY2FuIGJlIHVucmVnaXN0ZXJlZCB2aWEgYGFkZGl0aW9uYWxFdmVudE5hbWVzYFxyXG4gKlxyXG4gW3YzIENIQU5HRV0gT2ZmIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcclxuICpcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0gIHsuLi5zdHJpbmd9IGFkZGl0aW9uYWxFdmVudE5hbWVzXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT2ZmKGV2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXMpIHtcclxuICAgIGxldCBldmVudHNUb1JlbW92ZSA9IFtldmVudE5hbWUsIC4uLmFkZGl0aW9uYWxFdmVudE5hbWVzXTtcclxuICAgIGV2ZW50c1RvUmVtb3ZlLmZvckVhY2goZXZlbnROYW1lID0+IHtcclxuICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnROYW1lKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG4vKipcclxuICogT2ZmQWxsIHVucmVnaXN0ZXJzIGFsbCBsaXN0ZW5lcnNcclxuICogW3YzIENIQU5HRV0gT2ZmQWxsIG9ubHkgdW5yZWdpc3RlcnMgbGlzdGVuZXJzIHdpdGhpbiB0aGUgY3VycmVudCB3aW5kb3dcclxuICpcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmZBbGwoKSB7XHJcbiAgICBldmVudExpc3RlbmVycy5jbGVhcigpO1xyXG59XHJcblxyXG4vKipcclxuICogRW1pdCBhbiBldmVudFxyXG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50IFRoZSBldmVudCB0byBlbWl0XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gRW1pdChldmVudCkge1xyXG4gICAgdm9pZCBjYWxsKFwiRW1pdFwiLCBldmVudCk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5NZXNzYWdlRGlhbG9nT3B0aW9uc30gTWVzc2FnZURpYWxvZ09wdGlvbnNcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk9wZW5EaWFsb2dPcHRpb25zfSBPcGVuRGlhbG9nT3B0aW9uc1xyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuU2F2ZURpYWxvZ09wdGlvbnN9IFNhdmVEaWFsb2dPcHRpb25zXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5pbXBvcnQgeyBuYW5vaWQgfSBmcm9tICduYW5vaWQvbm9uLXNlY3VyZSc7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJkaWFsb2dcIik7XHJcblxyXG5sZXQgZGlhbG9nUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChkaWFsb2dSZXNwb25zZXMuaGFzKHJlc3VsdCkpO1xyXG4gICAgcmV0dXJuIHJlc3VsdDtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0NhbGxiYWNrKGlkLCBkYXRhLCBpc0pTT04pIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuZXhwb3J0IGZ1bmN0aW9uIGRpYWxvZ0Vycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcclxuICAgIGxldCBwID0gZGlhbG9nUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIHAucmVqZWN0KG1lc3NhZ2UpO1xyXG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBkaWFsb2codHlwZSwgb3B0aW9ucykge1xyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBsZXQgaWQgPSBnZW5lcmF0ZUlEKCk7XHJcbiAgICAgICAgb3B0aW9ucyA9IG9wdGlvbnMgfHwge307XHJcbiAgICAgICAgb3B0aW9uc1tcImRpYWxvZy1pZFwiXSA9IGlkO1xyXG4gICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcclxuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xyXG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xyXG4gICAgICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgICAgICB9KTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGFuIEluZm8gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEluZm8ob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhcIkluZm9cIiwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhbiBXYXJuaW5nIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBXYXJuaW5nKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBkaWFsb2coXCJXYXJuaW5nXCIsIG9wdGlvbnMpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgYW4gRXJyb3IgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEVycm9yKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBkaWFsb2coXCJFcnJvclwiLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGEgUXVlc3Rpb24gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7TWVzc2FnZURpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gVGhlIGxhYmVsIG9mIHRoZSBidXR0b24gcHJlc3NlZFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1ZXN0aW9uKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBkaWFsb2coXCJRdWVzdGlvblwiLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGFuIE9wZW4gZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7T3BlbkRpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nW118c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZSBvciBhbiBhcnJheSBvZiBzZWxlY3RlZCBmaWxlcyBpZiBBbGxvd3NNdWx0aXBsZVNlbGVjdGlvbiBpcyB0cnVlLiBBIGJsYW5rIHN0cmluZyBpcyByZXR1cm5lZCBpZiBubyBmaWxlIHdhcyBzZWxlY3RlZC5cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPcGVuRmlsZShvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gZGlhbG9nKFwiT3BlbkZpbGVcIiwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhIFNhdmUgZGlhbG9nIHdpdGggdGhlIGdpdmVuIG9wdGlvbnMuXHJcbiAqIEBwYXJhbSB7T3BlbkRpYWxvZ09wdGlvbnN9IG9wdGlvbnNcclxuICogQHJldHVybnMge1Byb21pc2U8c3RyaW5nPn0gUmV0dXJucyB0aGUgc2VsZWN0ZWQgZmlsZS4gQSBibGFuayBzdHJpbmcgaXMgcmV0dXJuZWQgaWYgbm8gZmlsZSB3YXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gU2F2ZUZpbGUob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhcIlNhdmVGaWxlXCIsIG9wdGlvbnMpO1xyXG59XHJcblxyXG4iLCAiaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjb250ZXh0bWVudVwiKTtcclxuXHJcbmZ1bmN0aW9uIG9wZW5Db250ZXh0TWVudShpZCwgeCwgeSwgZGF0YSkge1xyXG4gICAgcmV0dXJuIGNhbGwoXCJPcGVuQ29udGV4dE1lbnVcIiwge2lkLCB4LCB5LCBkYXRhfSk7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBlbmFibGVDb250ZXh0TWVudXMoZW5hYmxlZCkge1xyXG4gICAgaWYgKGVuYWJsZWQpIHtcclxuICAgICAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICB3aW5kb3cucmVtb3ZlRXZlbnRMaXN0ZW5lcignY29udGV4dG1lbnUnLCBjb250ZXh0TWVudUhhbmRsZXIpO1xyXG4gICAgfVxyXG59XHJcblxyXG5mdW5jdGlvbiBjb250ZXh0TWVudUhhbmRsZXIoZXZlbnQpIHtcclxuICAgIHByb2Nlc3NDb250ZXh0TWVudShldmVudC50YXJnZXQsIGV2ZW50KTtcclxufVxyXG5cclxuZnVuY3Rpb24gcHJvY2Vzc0NvbnRleHRNZW51KGVsZW1lbnQsIGV2ZW50KSB7XHJcbiAgICBsZXQgaWQgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS1jb250ZXh0bWVudScpO1xyXG4gICAgaWYgKGlkKSB7XHJcbiAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICBvcGVuQ29udGV4dE1lbnUoaWQsIGV2ZW50LmNsaWVudFgsIGV2ZW50LmNsaWVudFksIGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLWNvbnRleHRtZW51LWRhdGEnKSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGxldCBwYXJlbnQgPSBlbGVtZW50LnBhcmVudEVsZW1lbnQ7XHJcbiAgICAgICAgaWYgKHBhcmVudCkge1xyXG4gICAgICAgICAgICBwcm9jZXNzQ29udGV4dE1lbnUocGFyZW50LCBldmVudCk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG59XHJcbiIsICJcclxuaW1wb3J0IHtFbWl0LCBXYWlsc0V2ZW50fSBmcm9tIFwiLi9ldmVudHNcIjtcclxuaW1wb3J0IHtRdWVzdGlvbn0gZnJvbSBcIi4vZGlhbG9nc1wiO1xyXG5cclxuZnVuY3Rpb24gc2VuZEV2ZW50KGV2ZW50TmFtZSwgZGF0YT1udWxsKSB7XHJcbiAgICBsZXQgZXZlbnQgPSBuZXcgV2FpbHNFdmVudChldmVudE5hbWUsIGRhdGEpO1xyXG4gICAgRW1pdChldmVudCk7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtZXZlbnRdJyk7XHJcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XHJcbiAgICAgICAgY29uc3QgZXZlbnRUeXBlID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWV2ZW50Jyk7XHJcbiAgICAgICAgY29uc3QgY29uZmlybSA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC1jb25maXJtJyk7XHJcbiAgICAgICAgY29uc3QgdHJpZ2dlciA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC10cmlnZ2VyJykgfHwgXCJjbGlja1wiO1xyXG5cclxuICAgICAgICBsZXQgY2FsbGJhY2sgPSBmdW5jdGlvbiAoKSB7XHJcbiAgICAgICAgICAgIGlmIChjb25maXJtKSB7XHJcbiAgICAgICAgICAgICAgICBRdWVzdGlvbih7VGl0bGU6IFwiQ29uZmlybVwiLCBNZXNzYWdlOmNvbmZpcm0sIEJ1dHRvbnM6W3tMYWJlbDpcIlllc1wifSx7TGFiZWw6XCJOb1wiLCBJc0RlZmF1bHQ6dHJ1ZX1dfSkudGhlbihmdW5jdGlvbiAocmVzdWx0KSB7XHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3VsdCAhPT0gXCJOb1wiKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0pO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIHNlbmRFdmVudChldmVudFR5cGUpO1xyXG4gICAgICAgIH07XHJcbiAgICAgICAgLy8gUmVtb3ZlIGV4aXN0aW5nIGxpc3RlbmVyc1xyXG5cclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5mdW5jdGlvbiBjYWxsV2luZG93TWV0aG9kKG1ldGhvZCkge1xyXG4gICAgaWYgKHdhaWxzLldpbmRvd1ttZXRob2RdID09PSB1bmRlZmluZWQpIHtcclxuICAgICAgICBjb25zb2xlLmxvZyhcIldpbmRvdyBtZXRob2QgXCIgKyBtZXRob2QgKyBcIiBub3QgZm91bmRcIik7XHJcbiAgICB9XHJcbiAgICB3YWlscy5XaW5kb3dbbWV0aG9kXSgpO1xyXG59XHJcblxyXG5mdW5jdGlvbiBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKSB7XHJcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC13aW5kb3ddJyk7XHJcbiAgICBlbGVtZW50cy5mb3JFYWNoKGZ1bmN0aW9uIChlbGVtZW50KSB7XHJcbiAgICAgICAgY29uc3Qgd2luZG93TWV0aG9kID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXdpbmRvdycpO1xyXG4gICAgICAgIGNvbnN0IGNvbmZpcm0gPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtY29uZmlybScpO1xyXG4gICAgICAgIGNvbnN0IHRyaWdnZXIgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtdHJpZ2dlcicpIHx8IFwiY2xpY2tcIjtcclxuXHJcbiAgICAgICAgbGV0IGNhbGxiYWNrID0gZnVuY3Rpb24gKCkge1xyXG4gICAgICAgICAgICBpZiAoY29uZmlybSkge1xyXG4gICAgICAgICAgICAgICAgUXVlc3Rpb24oe1RpdGxlOiBcIkNvbmZpcm1cIiwgTWVzc2FnZTpjb25maXJtLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBjYWxsV2luZG93TWV0aG9kKHdpbmRvd01ldGhvZCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfSk7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xyXG4gICAgICAgIH07XHJcblxyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuICAgICAgICBlbGVtZW50LnJlbW92ZUV2ZW50TGlzdGVuZXIodHJpZ2dlciwgY2FsbGJhY2spO1xyXG5cclxuICAgICAgICAvLyBBZGQgbmV3IGxpc3RlbmVyXHJcbiAgICAgICAgZWxlbWVudC5hZGRFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuICAgIH0pO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gcmVsb2FkV01MKCkge1xyXG4gICAgYWRkV01MRXZlbnRMaXN0ZW5lcnMoKTtcclxuICAgIGFkZFdNTFdpbmRvd0xpc3RlbmVycygpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vLyBkZWZpbmVkIGluIHRoZSBUYXNrZmlsZVxyXG5leHBvcnQgbGV0IGludm9rZSA9IGZ1bmN0aW9uKGlucHV0KSB7XHJcbiAgICBpZihXSU5ET1dTKSB7XHJcbiAgICAgICAgY2hyb21lLndlYnZpZXcucG9zdE1lc3NhZ2UoaW5wdXQpO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICB3ZWJraXQubWVzc2FnZUhhbmRsZXJzLmV4dGVybmFsLnBvc3RNZXNzYWdlKGlucHV0KTtcclxuICAgIH1cclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxubGV0IGZsYWdzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gY29udmVydFRvTWFwKG9iaikge1xyXG4gICAgY29uc3QgbWFwID0gbmV3IE1hcCgpO1xyXG5cclxuICAgIGZvciAoY29uc3QgW2tleSwgdmFsdWVdIG9mIE9iamVjdC5lbnRyaWVzKG9iaikpIHtcclxuICAgICAgICBpZiAodHlwZW9mIHZhbHVlID09PSAnb2JqZWN0JyAmJiB2YWx1ZSAhPT0gbnVsbCkge1xyXG4gICAgICAgICAgICBtYXAuc2V0KGtleSwgY29udmVydFRvTWFwKHZhbHVlKSk7IC8vIFJlY3Vyc2l2ZWx5IGNvbnZlcnQgbmVzdGVkIG9iamVjdFxyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIG1hcC5zZXQoa2V5LCB2YWx1ZSk7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiBtYXA7XHJcbn1cclxuXHJcbmZldGNoKFwiL3dhaWxzL2ZsYWdzXCIpLnRoZW4oKHJlc3BvbnNlKSA9PiB7XHJcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xyXG4gICAgICAgIGZsYWdzID0gY29udmVydFRvTWFwKGRhdGEpO1xyXG4gICAgfSk7XHJcbn0pO1xyXG5cclxuXHJcbmZ1bmN0aW9uIGdldFZhbHVlRnJvbU1hcChrZXlTdHJpbmcpIHtcclxuICAgIGNvbnN0IGtleXMgPSBrZXlTdHJpbmcuc3BsaXQoJy4nKTtcclxuICAgIGxldCB2YWx1ZSA9IGZsYWdzO1xyXG5cclxuICAgIGZvciAoY29uc3Qga2V5IG9mIGtleXMpIHtcclxuICAgICAgICBpZiAodmFsdWUgaW5zdGFuY2VvZiBNYXApIHtcclxuICAgICAgICAgICAgdmFsdWUgPSB2YWx1ZS5nZXQoa2V5KTtcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICB2YWx1ZSA9IHZhbHVlW2tleV07XHJcbiAgICAgICAgfVxyXG5cclxuICAgICAgICBpZiAodmFsdWUgPT09IHVuZGVmaW5lZCkge1xyXG4gICAgICAgICAgICBicmVhaztcclxuICAgICAgICB9XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIHZhbHVlO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gR2V0RmxhZyhrZXlTdHJpbmcpIHtcclxuICAgIHJldHVybiBnZXRWYWx1ZUZyb21NYXAoa2V5U3RyaW5nKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtpbnZva2V9IGZyb20gXCIuL2ludm9rZVwiO1xyXG5pbXBvcnQge0dldEZsYWd9IGZyb20gXCIuL2ZsYWdzXCI7XHJcblxyXG5sZXQgc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGRyYWdUZXN0KGUpIHtcclxuICAgIC8vIGlmICh3aW5kb3cud2FpbHMuQ2FwYWJpbGl0aWVzWydIYXNOYXRpdmVEcmFnJ10gPT09IHRydWUpIHtcclxuICAgIC8vICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICAvLyB9XHJcblxyXG4gICAgbGV0IHZhbCA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGUudGFyZ2V0KS5nZXRQcm9wZXJ0eVZhbHVlKFwiYXBwLXJlZ2lvblwiKTtcclxuICAgIGlmICh2YWwpIHtcclxuICAgICAgICB2YWwgPSB2YWwudHJpbSgpO1xyXG4gICAgfVxyXG5cclxuICAgIGlmICh2YWwgIT09IFwiZHJhZ1wiKSB7XHJcbiAgICAgICAgcmV0dXJuIGZhbHNlO1xyXG4gICAgfVxyXG5cclxuICAgIC8vIE9ubHkgcHJvY2VzcyB0aGUgcHJpbWFyeSBidXR0b25cclxuICAgIGlmIChlLmJ1dHRvbnMgIT09IDEpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcblxyXG4gICAgcmV0dXJuIGUuZGV0YWlsID09PSAxO1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gc2V0dXBEcmFnKCkge1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlZG93bicsIG9uTW91c2VEb3duKTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZW1vdmUnLCBvbk1vdXNlTW92ZSk7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2V1cCcsIG9uTW91c2VVcCk7XHJcbn1cclxuXHJcbmxldCByZXNpemVFZGdlID0gbnVsbDtcclxuXHJcbmZ1bmN0aW9uIHRlc3RSZXNpemUoZSkge1xyXG4gICAgaWYoIHJlc2l6ZUVkZ2UgKSB7XHJcbiAgICAgICAgaW52b2tlKFwicmVzaXplOlwiICsgcmVzaXplRWRnZSk7XHJcbiAgICAgICAgcmV0dXJuIHRydWVcclxuICAgIH1cclxuICAgIHJldHVybiBmYWxzZTtcclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZURvd24oZSkge1xyXG5cclxuICAgIC8vIENoZWNrIGZvciByZXNpemluZyBvbiBXaW5kb3dzXHJcbiAgICBpZiggV0lORE9XUyApIHtcclxuICAgICAgICBpZiAodGVzdFJlc2l6ZSgpKSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICB9XHJcbiAgICBpZiAoZHJhZ1Rlc3QoZSkpIHtcclxuICAgICAgICAvLyBJZ25vcmUgZHJhZyBvbiBzY3JvbGxiYXJzXHJcbiAgICAgICAgaWYgKGUub2Zmc2V0WCA+IGUudGFyZ2V0LmNsaWVudFdpZHRoIHx8IGUub2Zmc2V0WSA+IGUudGFyZ2V0LmNsaWVudEhlaWdodCkge1xyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgfVxyXG4gICAgICAgIHNob3VsZERyYWcgPSB0cnVlO1xyXG4gICAgfSBlbHNlIHtcclxuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VVcChlKSB7XHJcbiAgICBsZXQgbW91c2VQcmVzc2VkID0gZS5idXR0b25zICE9PSB1bmRlZmluZWQgPyBlLmJ1dHRvbnMgOiBlLndoaWNoO1xyXG4gICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcclxuICAgICAgICBlbmREcmFnKCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBlbmREcmFnKCkge1xyXG4gICAgZG9jdW1lbnQuYm9keS5zdHlsZS5jdXJzb3IgPSAnZGVmYXVsdCc7XHJcbiAgICBzaG91bGREcmFnID0gZmFsc2U7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIHNldFJlc2l6ZShjdXJzb3IpIHtcclxuICAgIGRvY3VtZW50LmRvY3VtZW50RWxlbWVudC5zdHlsZS5jdXJzb3IgPSBjdXJzb3IgfHwgZGVmYXVsdEN1cnNvcjtcclxuICAgIHJlc2l6ZUVkZ2UgPSBjdXJzb3I7XHJcbn1cclxuXHJcbmZ1bmN0aW9uIG9uTW91c2VNb3ZlKGUpIHtcclxuICAgIGlmIChzaG91bGREcmFnKSB7XHJcbiAgICAgICAgc2hvdWxkRHJhZyA9IGZhbHNlO1xyXG4gICAgICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XHJcbiAgICAgICAgaWYgKG1vdXNlUHJlc3NlZCA+IDApIHtcclxuICAgICAgICAgICAgaW52b2tlKFwiZHJhZ1wiKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgcmV0dXJuO1xyXG4gICAgfVxyXG5cclxuICAgIGlmIChXSU5ET1dTKSB7XHJcbiAgICAgICAgaGFuZGxlUmVzaXplKGUpO1xyXG4gICAgfVxyXG59XHJcblxyXG5sZXQgZGVmYXVsdEN1cnNvciA9IFwiYXV0b1wiO1xyXG5cclxuZnVuY3Rpb24gaGFuZGxlUmVzaXplKGUpIHtcclxuICAgIGxldCByZXNpemVIYW5kbGVIZWlnaHQgPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZUhlaWdodFwiKSB8fCA1O1xyXG4gICAgbGV0IHJlc2l6ZUhhbmRsZVdpZHRoID0gR2V0RmxhZyhcInN5c3RlbS5yZXNpemVIYW5kbGVXaWR0aFwiKSB8fCA1O1xyXG5cclxuICAgIC8vIEV4dHJhIHBpeGVscyBmb3IgdGhlIGNvcm5lciBhcmVhc1xyXG4gICAgbGV0IGNvcm5lckV4dHJhID0gR2V0RmxhZyhcInJlc2l6ZUNvcm5lckV4dHJhXCIpIHx8IDM7XHJcblxyXG4gICAgbGV0IHJpZ2h0Qm9yZGVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGxldCBsZWZ0Qm9yZGVyID0gZS5jbGllbnRYIDwgcmVzaXplSGFuZGxlV2lkdGg7XHJcbiAgICBsZXQgdG9wQm9yZGVyID0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xyXG4gICAgbGV0IGJvdHRvbUJvcmRlciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IHJlc2l6ZUhhbmRsZUhlaWdodDtcclxuXHJcbiAgICAvLyBBZGp1c3QgZm9yIGNvcm5lcnNcclxuICAgIGxldCByaWdodENvcm5lciA9IHdpbmRvdy5vdXRlcldpZHRoIC0gZS5jbGllbnRYIDwgKHJlc2l6ZUhhbmRsZVdpZHRoICsgY29ybmVyRXh0cmEpO1xyXG4gICAgbGV0IGxlZnRDb3JuZXIgPSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBsZXQgdG9wQ29ybmVyID0gZS5jbGllbnRZIDwgKHJlc2l6ZUhhbmRsZUhlaWdodCArIGNvcm5lckV4dHJhKTtcclxuICAgIGxldCBib3R0b21Db3JuZXIgPSB3aW5kb3cub3V0ZXJIZWlnaHQgLSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG5cclxuICAgIC8vIElmIHdlIGFyZW4ndCBvbiBhbiBlZGdlLCBidXQgd2VyZSwgcmVzZXQgdGhlIGN1cnNvciB0byBkZWZhdWx0XHJcbiAgICBpZiAoIWxlZnRCb3JkZXIgJiYgIXJpZ2h0Qm9yZGVyICYmICF0b3BCb3JkZXIgJiYgIWJvdHRvbUJvcmRlciAmJiByZXNpemVFZGdlICE9PSB1bmRlZmluZWQpIHtcclxuICAgICAgICBzZXRSZXNpemUoKTtcclxuICAgIH1cclxuICAgIC8vIEFkanVzdGVkIGZvciBjb3JuZXIgYXJlYXNcclxuICAgIGVsc2UgaWYgKHJpZ2h0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic2UtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiBib3R0b21Db3JuZXIpIHNldFJlc2l6ZShcInN3LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRDb3JuZXIgJiYgdG9wQ29ybmVyKSBzZXRSZXNpemUoXCJudy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BDb3JuZXIgJiYgcmlnaHRDb3JuZXIpIHNldFJlc2l6ZShcIm5lLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGxlZnRCb3JkZXIpIHNldFJlc2l6ZShcInctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAodG9wQm9yZGVyKSBzZXRSZXNpemUoXCJuLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKGJvdHRvbUJvcmRlcikgc2V0UmVzaXplKFwicy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChyaWdodEJvcmRlcikgc2V0UmVzaXplKFwiZS1yZXNpemVcIik7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuXHJcbmltcG9ydCAqIGFzIENsaXBib2FyZCBmcm9tICcuL2NsaXBib2FyZCc7XHJcbmltcG9ydCAqIGFzIEFwcGxpY2F0aW9uIGZyb20gJy4vYXBwbGljYXRpb24nO1xyXG5pbXBvcnQgKiBhcyBMb2cgZnJvbSAnLi9sb2cnO1xyXG5pbXBvcnQgKiBhcyBTY3JlZW5zIGZyb20gJy4vc2NyZWVucyc7XHJcbmltcG9ydCB7UGx1Z2luLCBDYWxsLCBjYWxsRXJyb3JDYWxsYmFjaywgY2FsbENhbGxiYWNrfSBmcm9tIFwiLi9jYWxsc1wiO1xyXG5pbXBvcnQge25ld1dpbmRvd30gZnJvbSBcIi4vd2luZG93XCI7XHJcbmltcG9ydCB7ZGlzcGF0Y2hXYWlsc0V2ZW50LCBFbWl0LCBPZmYsIE9mZkFsbCwgT24sIE9uY2UsIE9uTXVsdGlwbGV9IGZyb20gXCIuL2V2ZW50c1wiO1xyXG5pbXBvcnQge2RpYWxvZ0NhbGxiYWNrLCBkaWFsb2dFcnJvckNhbGxiYWNrLCBFcnJvciwgSW5mbywgT3BlbkZpbGUsIFF1ZXN0aW9uLCBTYXZlRmlsZSwgV2FybmluZyx9IGZyb20gXCIuL2RpYWxvZ3NcIjtcclxuaW1wb3J0IHtlbmFibGVDb250ZXh0TWVudXN9IGZyb20gXCIuL2NvbnRleHRtZW51XCI7XHJcbmltcG9ydCB7cmVsb2FkV01MfSBmcm9tIFwiLi93bWxcIjtcclxuaW1wb3J0IHtzZXR1cERyYWcsIGVuZERyYWd9IGZyb20gXCIuL2RyYWdcIjtcclxuXHJcbndpbmRvdy53YWlscyA9IHtcclxuICAgIC4uLm5ld1J1bnRpbWUobnVsbCksXHJcbiAgICBDYXBhYmlsaXRpZXM6IHt9LFxyXG59O1xyXG5cclxuZmV0Y2goXCIvd2FpbHMvY2FwYWJpbGl0aWVzXCIpLnRoZW4oKHJlc3BvbnNlKSA9PiB7XHJcbiAgICByZXNwb25zZS5qc29uKCkudGhlbigoZGF0YSkgPT4ge1xyXG4gICAgICAgIHdpbmRvdy53YWlscy5DYXBhYmlsaXRpZXMgPSBkYXRhO1xyXG4gICAgfSk7XHJcbn0pO1xyXG5cclxuLy8gSW50ZXJuYWwgd2FpbHMgZW5kcG9pbnRzXHJcbndpbmRvdy5fd2FpbHMgPSB7XHJcbiAgICBkaWFsb2dDYWxsYmFjayxcclxuICAgIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssXHJcbiAgICBkaXNwYXRjaFdhaWxzRXZlbnQsXHJcbiAgICBjYWxsQ2FsbGJhY2ssXHJcbiAgICBjYWxsRXJyb3JDYWxsYmFjayxcclxuICAgIGVuZERyYWcsXHJcbn07XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gbmV3UnVudGltZSh3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4ge1xyXG4gICAgICAgIENsaXBib2FyZDoge1xyXG4gICAgICAgICAgICAuLi5DbGlwYm9hcmRcclxuICAgICAgICB9LFxyXG4gICAgICAgIEFwcGxpY2F0aW9uOiB7XHJcbiAgICAgICAgICAgIC4uLkFwcGxpY2F0aW9uLFxyXG4gICAgICAgICAgICBHZXRXaW5kb3dCeU5hbWUod2luZG93TmFtZSkge1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuIG5ld1J1bnRpbWUod2luZG93TmFtZSk7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9LFxyXG4gICAgICAgIExvZyxcclxuICAgICAgICBTY3JlZW5zLFxyXG4gICAgICAgIENhbGwsXHJcbiAgICAgICAgUGx1Z2luLFxyXG4gICAgICAgIFdNTDoge1xyXG4gICAgICAgICAgICBSZWxvYWQ6IHJlbG9hZFdNTCxcclxuICAgICAgICB9LFxyXG4gICAgICAgIERpYWxvZzoge1xyXG4gICAgICAgICAgICBJbmZvLFxyXG4gICAgICAgICAgICBXYXJuaW5nLFxyXG4gICAgICAgICAgICBFcnJvcixcclxuICAgICAgICAgICAgUXVlc3Rpb24sXHJcbiAgICAgICAgICAgIE9wZW5GaWxlLFxyXG4gICAgICAgICAgICBTYXZlRmlsZSxcclxuICAgICAgICB9LFxyXG4gICAgICAgIEV2ZW50czoge1xyXG4gICAgICAgICAgICBFbWl0LFxyXG4gICAgICAgICAgICBPbixcclxuICAgICAgICAgICAgT25jZSxcclxuICAgICAgICAgICAgT25NdWx0aXBsZSxcclxuICAgICAgICAgICAgT2ZmLFxyXG4gICAgICAgICAgICBPZmZBbGwsXHJcbiAgICAgICAgfSxcclxuICAgICAgICBXaW5kb3c6IG5ld1dpbmRvdyh3aW5kb3dOYW1lKSxcclxuICAgIH07XHJcbn1cclxuXHJcbmlmIChERUJVRykge1xyXG4gICAgY29uc29sZS5sb2coXCJXYWlscyB2My4wLjAgRGVidWcgTW9kZSBFbmFibGVkXCIpO1xyXG59XHJcblxyXG5lbmFibGVDb250ZXh0TWVudXModHJ1ZSk7XHJcblxyXG5zZXR1cERyYWcoKTtcclxuXHJcbmRvY3VtZW50LmFkZEV2ZW50TGlzdGVuZXIoXCJET01Db250ZW50TG9hZGVkXCIsIGZ1bmN0aW9uKGV2ZW50KSB7XHJcbiAgICByZWxvYWRXTUwoKTtcclxufSk7Il0sCiAgIm1hcHBpbmdzIjogIjs7Ozs7Ozs7QUFBQTtBQUFBO0FBQUE7QUFBQTtBQUFBOzs7QUNZQSxNQUFNLGFBQWEsT0FBTyxTQUFTLFNBQVM7QUFFNUMsV0FBUyxZQUFZLFFBQVEsWUFBWSxNQUFNO0FBQzNDLFFBQUksTUFBTSxJQUFJLElBQUksVUFBVTtBQUM1QixRQUFJLGFBQWEsT0FBTyxVQUFVLE1BQU07QUFDeEMsUUFBSSxNQUFNO0FBQ04sVUFBSSxhQUFhLE9BQU8sUUFBUSxLQUFLLFVBQVUsSUFBSSxDQUFDO0FBQUEsSUFDeEQ7QUFDQSxRQUFJLGVBQWU7QUFBQSxNQUNmLFNBQVMsQ0FBQztBQUFBLElBQ2Q7QUFDQSxRQUFJLFlBQVk7QUFDWixtQkFBYSxRQUFRLHFCQUFxQixJQUFJO0FBQUEsSUFDbEQ7QUFDQSxXQUFPLElBQUksUUFBUSxDQUFDLFNBQVMsV0FBVztBQUNwQyxZQUFNLEtBQUssWUFBWSxFQUNsQixLQUFLLGNBQVk7QUFDZCxZQUFJLFNBQVMsSUFBSTtBQUViLGNBQUksU0FBUyxRQUFRLElBQUksY0FBYyxLQUFLLFNBQVMsUUFBUSxJQUFJLGNBQWMsRUFBRSxRQUFRLGtCQUFrQixNQUFNLElB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WTtBQUNuQyxXQUFPO0FBQUEsTUFDSCxXQUFXO0FBQUEsUUFDUCxHQUFHO0FBQUEsTUFDUDtBQUFBLE1BQ0EsYUFBYTtBQUFBLFFBQ1QsR0FBRztBQUFBLFFBQ0gsZ0JBQWdCQyxhQUFZO0FBQ3hCLGlCQUFPLFdBQVdBLFdBQVU7QUFBQSxRQUNoQztBQUFBLE1BQ0o7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQSxLQUFLO0FBQUEsUUFDRCxRQUFRO0FBQUEsTUFDWjtBQUFBLE1BQ0EsUUFBUTtBQUFBLFFBQ0o7QUFBQSxRQUNBO0FBQUEsUUFDQSxPQUFBQztBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLE1BQ0o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRLFVBQVUsVUFBVTtBQUFBLElBQ2hDO0FBQUEsRUFDSjtBQUVBLE1BQUksTUFBTztBQUNQLFlBQVEsSUFBSSxpQ0FBaUM7QUFBQSxFQUNqRDtBQUVBLHFCQUFtQixJQUFJO0FBRXZCLFlBQVU7QUFFVixXQUFTLGlCQUFpQixvQkFBb0IsU0FBUyxPQUFPO0FBQzFELGNBQVU7QUFBQSxFQUNkLENBQUM7IiwKICAibmFtZXMiOiBbImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImV2ZW50TmFtZSIsICJjYWxsIiwgImdlbmVyYXRlSUQiLCAiRXJyb3IiLCAiY2FsbCIsICJ3aW5kb3dOYW1lIiwgIkVycm9yIl0KfQo=
