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
    let element = event.target;
    const computedStyle = window.getComputedStyle(element);
    let defaultContextMenuAction = computedStyle.getPropertyValue("--default-contextmenu").trim();
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
        let selection = window.getSelection();
        if (selection && selection.toString().length > 0) {
          for (let i = 0; i < selection.rangeCount; i++) {
            let range = selection.getRangeAt(i);
            let rects = range.getClientRects();
            for (let j = 0; j < rects.length; j++) {
              let rect = rects[j];
              if (document.elementFromPoint(rect.left, rect.top) === element) {
                return;
              }
            }
          }
        }
        if (element.tagName === "INPUT" || element.tagName === "TEXTAREA") {
          if (!element.readOnly && !element.disabled) {
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
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsiZGVza3RvcC9jbGlwYm9hcmQuanMiLCAiZGVza3RvcC9ydW50aW1lLmpzIiwgImRlc2t0b3AvYXBwbGljYXRpb24uanMiLCAiZGVza3RvcC9sb2cuanMiLCAiZGVza3RvcC9zY3JlZW5zLmpzIiwgIm5vZGVfbW9kdWxlcy9uYW5vaWQvbm9uLXNlY3VyZS9pbmRleC5qcyIsICJkZXNrdG9wL2NhbGxzLmpzIiwgImRlc2t0b3Avd2luZG93LmpzIiwgImRlc2t0b3AvZXZlbnRzLmpzIiwgImRlc2t0b3AvZGlhbG9ncy5qcyIsICJkZXNrdG9wL2NvbnRleHRtZW51LmpzIiwgImRlc2t0b3Avd21sLmpzIiwgImRlc2t0b3AvaW52b2tlLmpzIiwgImRlc2t0b3AvZmxhZ3MuanMiLCAiZGVza3RvcC9kcmFnLmpzIiwgImRlc2t0b3AvbWFpbi5qcyJdLAogICJzb3VyY2VzQ29udGVudCI6IFsiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJjbGlwYm9hcmRcIik7XHJcblxyXG4vKipcclxuICogU2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNldFRleHQodGV4dCkge1xyXG4gICAgdm9pZCBjYWxsKFwiU2V0VGV4dFwiLCB7dGV4dH0pO1xyXG59XHJcblxyXG4vKipcclxuICogR2V0IHRoZSBDbGlwYm9hcmQgdGV4dFxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFRleHQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIlRleHRcIik7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuY29uc3QgcnVudGltZVVSTCA9IHdpbmRvdy5sb2NhdGlvbi5vcmlnaW4gKyBcIi93YWlscy9ydW50aW1lXCI7XHJcblxyXG5mdW5jdGlvbiBydW50aW1lQ2FsbChtZXRob2QsIHdpbmRvd05hbWUsIGFyZ3MpIHtcclxuICAgIGxldCB1cmwgPSBuZXcgVVJMKHJ1bnRpbWVVUkwpO1xyXG4gICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJtZXRob2RcIiwgbWV0aG9kKTtcclxuICAgIGlmIChhcmdzKSB7XHJcbiAgICAgICAgdXJsLnNlYXJjaFBhcmFtcy5hcHBlbmQoXCJhcmdzXCIsIEpTT04uc3RyaW5naWZ5KGFyZ3MpKTtcclxuICAgIH1cclxuICAgIGxldCBmZXRjaE9wdGlvbnMgPSB7XHJcbiAgICAgICAgaGVhZGVyczoge30sXHJcbiAgICB9O1xyXG4gICAgaWYgKHdpbmRvd05hbWUpIHtcclxuICAgICAgICBmZXRjaE9wdGlvbnMuaGVhZGVyc1tcIngtd2FpbHMtd2luZG93LW5hbWVcIl0gPSB3aW5kb3dOYW1lO1xyXG4gICAgfVxyXG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcclxuICAgICAgICBmZXRjaCh1cmwsIGZldGNoT3B0aW9ucylcclxuICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4ge1xyXG4gICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLm9rKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgLy8gY2hlY2sgY29udGVudCB0eXBlXHJcbiAgICAgICAgICAgICAgICAgICAgaWYgKHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpICYmIHJlc3BvbnNlLmhlYWRlcnMuZ2V0KFwiQ29udGVudC1UeXBlXCIpLmluZGV4T2YoXCJhcHBsaWNhdGlvbi9qc29uXCIpICE9PSAtMSkge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpO1xyXG4gICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybiByZXNwb25zZS50ZXh0KCk7XHJcbiAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgcmVqZWN0KEVycm9yKHJlc3BvbnNlLnN0YXR1c1RleHQpKTtcclxuICAgICAgICAgICAgfSlcclxuICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiByZXNvbHZlKGRhdGEpKVxyXG4gICAgICAgICAgICAuY2F0Y2goZXJyb3IgPT4gcmVqZWN0KGVycm9yKSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWVDYWxsZXIob2JqZWN0LCB3aW5kb3dOYW1lKSB7XHJcbiAgICByZXR1cm4gZnVuY3Rpb24gKG1ldGhvZCwgYXJncz1udWxsKSB7XHJcbiAgICAgICAgcmV0dXJuIHJ1bnRpbWVDYWxsKG9iamVjdCArIFwiLlwiICsgbWV0aG9kLCB3aW5kb3dOYW1lLCBhcmdzKTtcclxuICAgIH07XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJhcHBsaWNhdGlvblwiKTtcclxuXHJcbi8qKlxyXG4gKiBIaWRlIHRoZSBhcHBsaWNhdGlvblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEhpZGUoKSB7XHJcbiAgICB2b2lkIGNhbGwoXCJIaWRlXCIpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvdyB0aGUgYXBwbGljYXRpb25cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBTaG93KCkge1xyXG4gICAgdm9pZCBjYWxsKFwiU2hvd1wiKTtcclxufVxyXG5cclxuXHJcbi8qKlxyXG4gKiBRdWl0IHRoZSBhcHBsaWNhdGlvblxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFF1aXQoKSB7XHJcbiAgICB2b2lkIGNhbGwoXCJRdWl0XCIpO1xyXG59IiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwibG9nXCIpO1xyXG5cclxuLyoqXHJcbiAqIExvZ3MgYSBtZXNzYWdlLlxyXG4gKiBAcGFyYW0ge21lc3NhZ2V9IE1lc3NhZ2UgdG8gbG9nXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gTG9nKG1lc3NhZ2UpIHtcclxuICAgIHJldHVybiBjYWxsKFwiTG9nXCIsIG1lc3NhZ2UpO1xyXG59XHJcbiIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5sZXQgY2FsbCA9IG5ld1J1bnRpbWVDYWxsZXIoXCJzY3JlZW5zXCIpO1xyXG5cclxuLyoqXHJcbiAqIEdldHMgYWxsIHNjcmVlbnMuXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbltdPn1cclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBHZXRBbGwoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIkdldEFsbFwiKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIEdldHMgdGhlIHByaW1hcnkgc2NyZWVuLlxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxTY3JlZW4+fVxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldFByaW1hcnkoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIkdldFByaW1hcnlcIik7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBHZXRzIHRoZSBjdXJyZW50IGFjdGl2ZSBzY3JlZW4uXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XHJcbiAqIEBjb25zdHJ1Y3RvclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEN1cnJlbnQoKSB7XHJcbiAgICByZXR1cm4gY2FsbChcIkdldEN1cnJlbnRcIik7XHJcbn0iLCAibGV0IHVybEFscGhhYmV0ID1cbiAgJ3VzZWFuZG9tLTI2VDE5ODM0MFBYNzVweEpBQ0tWRVJZTUlOREJVU0hXT0xGX0dRWmJmZ2hqa2xxdnd5enJpY3QnXG5leHBvcnQgbGV0IGN1c3RvbUFscGhhYmV0ID0gKGFscGhhYmV0LCBkZWZhdWx0U2l6ZSA9IDIxKSA9PiB7XG4gIHJldHVybiAoc2l6ZSA9IGRlZmF1bHRTaXplKSA9PiB7XG4gICAgbGV0IGlkID0gJydcbiAgICBsZXQgaSA9IHNpemVcbiAgICB3aGlsZSAoaS0tKSB7XG4gICAgICBpZCArPSBhbHBoYWJldFsoTWF0aC5yYW5kb20oKSAqIGFscGhhYmV0Lmxlbmd0aCkgfCAwXVxuICAgIH1cbiAgICByZXR1cm4gaWRcbiAgfVxufVxuZXhwb3J0IGxldCBuYW5vaWQgPSAoc2l6ZSA9IDIxKSA9PiB7XG4gIGxldCBpZCA9ICcnXG4gIGxldCBpID0gc2l6ZVxuICB3aGlsZSAoaS0tKSB7XG4gICAgaWQgKz0gdXJsQWxwaGFiZXRbKE1hdGgucmFuZG9tKCkgKiA2NCkgfCAwXVxuICB9XG4gIHJldHVybiBpZFxufVxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxuaW1wb3J0IHsgbmFub2lkIH0gZnJvbSAnbmFub2lkL25vbi1zZWN1cmUnO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY2FsbFwiKTtcclxuXHJcbmxldCBjYWxsUmVzcG9uc2VzID0gbmV3IE1hcCgpO1xyXG5cclxuZnVuY3Rpb24gZ2VuZXJhdGVJRCgpIHtcclxuICAgIGxldCByZXN1bHQ7XHJcbiAgICBkbyB7XHJcbiAgICAgICAgcmVzdWx0ID0gbmFub2lkKCk7XHJcbiAgICB9IHdoaWxlIChjYWxsUmVzcG9uc2VzLmhhcyhyZXN1bHQpKTtcclxuICAgIHJldHVybiByZXN1bHQ7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBjYWxsQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgbGV0IHAgPSBjYWxsUmVzcG9uc2VzLmdldChpZCk7XHJcbiAgICBpZiAocCkge1xyXG4gICAgICAgIGlmIChpc0pTT04pIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKEpTT04ucGFyc2UoZGF0YSkpO1xyXG4gICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgIHAucmVzb2x2ZShkYXRhKTtcclxuICAgICAgICB9XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgfVxyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gY2FsbEVycm9yQ2FsbGJhY2soaWQsIG1lc3NhZ2UpIHtcclxuICAgIGxldCBwID0gY2FsbFJlc3BvbnNlcy5nZXQoaWQpO1xyXG4gICAgaWYgKHApIHtcclxuICAgICAgICBwLnJlamVjdChtZXNzYWdlKTtcclxuICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGNhbGxCaW5kaW5nKHR5cGUsIG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XHJcbiAgICAgICAgbGV0IGlkID0gZ2VuZXJhdGVJRCgpO1xyXG4gICAgICAgIG9wdGlvbnMgPSBvcHRpb25zIHx8IHt9O1xyXG4gICAgICAgIG9wdGlvbnNbXCJjYWxsLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgY2FsbFJlc3BvbnNlcy5zZXQoaWQsIHtyZXNvbHZlLCByZWplY3R9KTtcclxuICAgICAgICBjYWxsKHR5cGUsIG9wdGlvbnMpLmNhdGNoKChlcnJvcikgPT4ge1xyXG4gICAgICAgICAgICByZWplY3QoZXJyb3IpO1xyXG4gICAgICAgICAgICBjYWxsUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICAgICAgfSk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIENhbGwob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIENhbGwgYSBwbHVnaW4gbWV0aG9kXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBwbHVnaW5OYW1lIC0gbmFtZSBvZiB0aGUgcGx1Z2luXHJcbiAqIEBwYXJhbSB7c3RyaW5nfSBtZXRob2ROYW1lIC0gbmFtZSBvZiB0aGUgbWV0aG9kXHJcbiAqIEBwYXJhbSB7Li4uYW55fSBhcmdzIC0gYXJndW1lbnRzIHRvIHBhc3MgdG8gdGhlIG1ldGhvZFxyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxhbnk+fSAtIHByb21pc2UgdGhhdCByZXNvbHZlcyB3aXRoIHRoZSByZXN1bHRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBQbHVnaW4ocGx1Z2luTmFtZSwgbWV0aG9kTmFtZSwgLi4uYXJncykge1xyXG4gICAgcmV0dXJuIGNhbGxCaW5kaW5nKFwiQ2FsbFwiLCB7XHJcbiAgICAgICAgcGFja2FnZU5hbWU6IFwid2FpbHMtcGx1Z2luc1wiLFxyXG4gICAgICAgIHN0cnVjdE5hbWU6IHBsdWdpbk5hbWUsXHJcbiAgICAgICAgbWV0aG9kTmFtZTogbWV0aG9kTmFtZSxcclxuICAgICAgICBhcmdzOiBhcmdzLFxyXG4gICAgfSk7XHJcbn0iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLyoqXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuLi9hcGkvdHlwZXNcIikuU2l6ZX0gU2l6ZVxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlBvc2l0aW9ufSBQb3NpdGlvblxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi4vYXBpL3R5cGVzXCIpLlNjcmVlbn0gU2NyZWVuXHJcbiAqL1xyXG5cclxuaW1wb3J0IHtuZXdSdW50aW1lQ2FsbGVyfSBmcm9tIFwiLi9ydW50aW1lXCI7XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gbmV3V2luZG93KHdpbmRvd05hbWUpIHtcclxuICAgIGxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcIndpbmRvd1wiLCB3aW5kb3dOYW1lKTtcclxuICAgIHJldHVybiB7XHJcbiAgICAgICAgLy8gUmVsb2FkOiAoKSA9PiBjYWxsKCdXUicpLFxyXG4gICAgICAgIC8vIFJlbG9hZEFwcDogKCkgPT4gY2FsbCgnV1InKSxcclxuICAgICAgICAvLyBTZXRTeXN0ZW1EZWZhdWx0VGhlbWU6ICgpID0+IGNhbGwoJ1dBU0RUJyksXHJcbiAgICAgICAgLy8gU2V0TGlnaHRUaGVtZTogKCkgPT4gY2FsbCgnV0FMVCcpLFxyXG4gICAgICAgIC8vIFNldERhcmtUaGVtZTogKCkgPT4gY2FsbCgnV0FEVCcpLFxyXG4gICAgICAgIC8vIElzRnVsbHNjcmVlbjogKCkgPT4gY2FsbCgnV0lGJyksXHJcbiAgICAgICAgLy8gSXNNYXhpbWl6ZWQ6ICgpID0+IGNhbGwoJ1dJTScpLFxyXG4gICAgICAgIC8vIElzTWluaW1pemVkOiAoKSA9PiBjYWxsKCdXSU1OJyksXHJcbiAgICAgICAgLy8gSXNXaW5kb3dlZDogKCkgPT4gY2FsbCgnV0lGJyksXHJcblxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBDZW50ZXJzIHRoZSB3aW5kb3cuXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgQ2VudGVyOiAoKSA9PiB2b2lkIGNhbGwoJ0NlbnRlcicpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyB0aXRsZS5cclxuICAgICAgICAgKiBAcGFyYW0gdGl0bGVcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRUaXRsZTogKHRpdGxlKSA9PiB2b2lkIGNhbGwoJ1NldFRpdGxlJywge3RpdGxlfSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIE1ha2VzIHRoZSB3aW5kb3cgZnVsbHNjcmVlbi5cclxuICAgICAgICAgKi9cclxuICAgICAgICBGdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ0Z1bGxzY3JlZW4nKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogVW5mdWxsc2NyZWVuIHRoZSB3aW5kb3cuXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgVW5GdWxsc2NyZWVuOiAoKSA9PiB2b2lkIGNhbGwoJ1VuRnVsbHNjcmVlbicpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgdGhlIHdpbmRvdyBzaXplLlxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aCBUaGUgd2luZG93IHdpZHRoXHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IGhlaWdodCBUaGUgd2luZG93IGhlaWdodFxyXG4gICAgICAgICAqL1xyXG4gICAgICAgIFNldFNpemU6ICh3aWR0aCwgaGVpZ2h0KSA9PiBjYWxsKCdTZXRTaXplJywge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHdpbmRvdyBzaXplLlxyXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNpemU+fSBUaGUgd2luZG93IHNpemVcclxuICAgICAgICAgKi9cclxuICAgICAgICBTaXplOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTaXplJyk7IH0sXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IG1heGltdW0gc2l6ZS5cclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gd2lkdGhcclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gaGVpZ2h0XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2V0TWF4U2l6ZTogKHdpZHRoLCBoZWlnaHQpID0+IHZvaWQgY2FsbCgnU2V0TWF4U2l6ZScsIHt3aWR0aCxoZWlnaHR9KSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogU2V0IHRoZSB3aW5kb3cgbWluaW11bSBzaXplLlxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB3aWR0aFxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSBoZWlnaHRcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRNaW5TaXplOiAod2lkdGgsIGhlaWdodCkgPT4gdm9pZCBjYWxsKCdTZXRNaW5TaXplJywge3dpZHRoLGhlaWdodH0pLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTZXQgd2luZG93IHRvIGJlIGFsd2F5cyBvbiB0b3AuXHJcbiAgICAgICAgICogQHBhcmFtIHtib29sZWFufSBvblRvcCBXaGV0aGVyIHRoZSB3aW5kb3cgc2hvdWxkIGJlIGFsd2F5cyBvbiB0b3BcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRBbHdheXNPblRvcDogKG9uVG9wKSA9PiB2b2lkIGNhbGwoJ1NldEFsd2F5c09uVG9wJywge2Fsd2F5c09uVG9wOm9uVG9wfSksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgd2luZG93IHJlbGF0aXZlIHBvc2l0aW9uLlxyXG4gICAgICAgICAqIEBwYXJhbSB7bnVtYmVyfSB4XHJcbiAgICAgICAgICogQHBhcmFtIHtudW1iZXJ9IHlcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRSZWxhdGl2ZVBvc2l0aW9uOiAoeCwgeSkgPT4gY2FsbCgnU2V0UmVsYXRpdmVQb3NpdGlvbicsIHt4LHl9KSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogR2V0IHRoZSB3aW5kb3cgcG9zaXRpb24uXHJcbiAgICAgICAgICogQHJldHVybnMge1Byb21pc2U8UG9zaXRpb24+fSBUaGUgd2luZG93IHBvc2l0aW9uXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgUmVsYXRpdmVQb3NpdGlvbjogKCkgPT4geyByZXR1cm4gY2FsbCgnUmVsYXRpdmVQb3NpdGlvbicpOyB9LFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBHZXQgdGhlIHNjcmVlbiB0aGUgd2luZG93IGlzIG9uLlxyXG4gICAgICAgICAqIEByZXR1cm5zIHtQcm9taXNlPFNjcmVlbj59XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgU2NyZWVuOiAoKSA9PiB7IHJldHVybiBjYWxsKCdTY3JlZW4nKTsgfSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogSGlkZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgSGlkZTogKCkgPT4gdm9pZCBjYWxsKCdIaWRlJyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIE1heGltaXNlIHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBNYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdNYXhpbWlzZScpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBTaG93IHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBTaG93OiAoKSA9PiB2b2lkIGNhbGwoJ1Nob3cnKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogQ2xvc2UgdGhlIHdpbmRvd1xyXG4gICAgICAgICAqL1xyXG4gICAgICAgIENsb3NlOiAoKSA9PiB2b2lkIGNhbGwoJ0Nsb3NlJyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFRvZ2dsZSB0aGUgd2luZG93IG1heGltaXNlIHN0YXRlXHJcbiAgICAgICAgICovXHJcbiAgICAgICAgVG9nZ2xlTWF4aW1pc2U6ICgpID0+IHZvaWQgY2FsbCgnVG9nZ2xlTWF4aW1pc2UnKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogVW5tYXhpbWlzZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgVW5NYXhpbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdVbk1heGltaXNlJyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIE1pbmltaXNlIHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBNaW5pbWlzZTogKCkgPT4gdm9pZCBjYWxsKCdNaW5pbWlzZScpLFxyXG5cclxuICAgICAgICAvKipcclxuICAgICAgICAgKiBVbm1pbmltaXNlIHRoZSB3aW5kb3dcclxuICAgICAgICAgKi9cclxuICAgICAgICBVbk1pbmltaXNlOiAoKSA9PiB2b2lkIGNhbGwoJ1VuTWluaW1pc2UnKSxcclxuXHJcbiAgICAgICAgLyoqXHJcbiAgICAgICAgICogUmVzdG9yZSB0aGUgd2luZG93XHJcbiAgICAgICAgICovXHJcbiAgICAgICAgUmVzdG9yZTogKCkgPT4gdm9pZCBjYWxsKCdSZXN0b3JlJyksXHJcblxyXG4gICAgICAgIC8qKlxyXG4gICAgICAgICAqIFNldCB0aGUgYmFja2dyb3VuZCBjb2xvdXIgb2YgdGhlIHdpbmRvdy5cclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gciAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gZyAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gYiAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcclxuICAgICAgICAgKiBAcGFyYW0ge251bWJlcn0gYSAtIEEgdmFsdWUgYmV0d2VlbiAwIGFuZCAyNTVcclxuICAgICAgICAgKi9cclxuICAgICAgICBTZXRCYWNrZ3JvdW5kQ29sb3VyOiAociwgZywgYiwgYSkgPT4gdm9pZCBjYWxsKCdTZXRCYWNrZ3JvdW5kQ29sb3VyJywge3IsIGcsIGIsIGF9KSxcclxuICAgIH07XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbi8qKlxyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuV2FpbHNFdmVudH0gV2FpbHNFdmVudFxyXG4gKi9cclxuXHJcbmltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiZXZlbnRzXCIpO1xyXG5cclxuLyoqXHJcbiAqIFRoZSBMaXN0ZW5lciBjbGFzcyBkZWZpbmVzIGEgbGlzdGVuZXIhIDotKVxyXG4gKlxyXG4gKiBAY2xhc3MgTGlzdGVuZXJcclxuICovXHJcbmNsYXNzIExpc3RlbmVyIHtcclxuICAgIC8qKlxyXG4gICAgICogQ3JlYXRlcyBhbiBpbnN0YW5jZSBvZiBMaXN0ZW5lci5cclxuICAgICAqIEBwYXJhbSB7c3RyaW5nfSBldmVudE5hbWVcclxuICAgICAqIEBwYXJhbSB7ZnVuY3Rpb259IGNhbGxiYWNrXHJcbiAgICAgKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAgICAgKiBAbWVtYmVyb2YgTGlzdGVuZXJcclxuICAgICAqL1xyXG4gICAgY29uc3RydWN0b3IoZXZlbnROYW1lLCBjYWxsYmFjaywgbWF4Q2FsbGJhY2tzKSB7XHJcbiAgICAgICAgdGhpcy5ldmVudE5hbWUgPSBldmVudE5hbWU7XHJcbiAgICAgICAgLy8gRGVmYXVsdCBvZiAtMSBtZWFucyBpbmZpbml0ZVxyXG4gICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzID0gbWF4Q2FsbGJhY2tzIHx8IC0xO1xyXG4gICAgICAgIC8vIENhbGxiYWNrIGludm9rZXMgdGhlIGNhbGxiYWNrIHdpdGggdGhlIGdpdmVuIGRhdGFcclxuICAgICAgICAvLyBSZXR1cm5zIHRydWUgaWYgdGhpcyBsaXN0ZW5lciBzaG91bGQgYmUgZGVzdHJveWVkXHJcbiAgICAgICAgdGhpcy5DYWxsYmFjayA9IChkYXRhKSA9PiB7XHJcbiAgICAgICAgICAgIGNhbGxiYWNrKGRhdGEpO1xyXG4gICAgICAgICAgICAvLyBJZiBtYXhDYWxsYmFja3MgaXMgaW5maW5pdGUsIHJldHVybiBmYWxzZSAoZG8gbm90IGRlc3Ryb3kpXHJcbiAgICAgICAgICAgIGlmICh0aGlzLm1heENhbGxiYWNrcyA9PT0gLTEpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBEZWNyZW1lbnQgbWF4Q2FsbGJhY2tzLiBSZXR1cm4gdHJ1ZSBpZiBub3cgMCwgb3RoZXJ3aXNlIGZhbHNlXHJcbiAgICAgICAgICAgIHRoaXMubWF4Q2FsbGJhY2tzIC09IDE7XHJcbiAgICAgICAgICAgIHJldHVybiB0aGlzLm1heENhbGxiYWNrcyA9PT0gMDtcclxuICAgICAgICB9O1xyXG4gICAgfVxyXG59XHJcblxyXG5cclxuLyoqXHJcbiAqIFdhaWxzRXZlbnQgZGVmaW5lcyBhIGN1c3RvbSBldmVudC4gSXQgaXMgcGFzc2VkIHRvIGV2ZW50IGxpc3RlbmVycy5cclxuICpcclxuICogQGNsYXNzIFdhaWxzRXZlbnRcclxuICogQHByb3BlcnR5IHtzdHJpbmd9IG5hbWUgLSBOYW1lIG9mIHRoZSBldmVudFxyXG4gKiBAcHJvcGVydHkge2FueX0gZGF0YSAtIERhdGEgYXNzb2NpYXRlZCB3aXRoIHRoZSBldmVudFxyXG4gKi9cclxuZXhwb3J0IGNsYXNzIFdhaWxzRXZlbnQge1xyXG4gICAgLyoqXHJcbiAgICAgKiBDcmVhdGVzIGFuIGluc3RhbmNlIG9mIFdhaWxzRXZlbnQuXHJcbiAgICAgKiBAcGFyYW0ge3N0cmluZ30gbmFtZSAtIE5hbWUgb2YgdGhlIGV2ZW50XHJcbiAgICAgKiBAcGFyYW0ge2FueT1udWxsfSBkYXRhIC0gRGF0YSBhc3NvY2lhdGVkIHdpdGggdGhlIGV2ZW50XHJcbiAgICAgKiBAbWVtYmVyb2YgV2FpbHNFdmVudFxyXG4gICAgICovXHJcbiAgICBjb25zdHJ1Y3RvcihuYW1lLCBkYXRhID0gbnVsbCkge1xyXG4gICAgICAgIHRoaXMubmFtZSA9IG5hbWU7XHJcbiAgICAgICAgdGhpcy5kYXRhID0gZGF0YTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGNvbnN0IGV2ZW50TGlzdGVuZXJzID0gbmV3IE1hcCgpO1xyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBgbWF4Q2FsbGJhY2tzYCB0aW1lcyBiZWZvcmUgYmVpbmcgZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xyXG4gKiBAcGFyYW0ge251bWJlcn0gbWF4Q2FsbGJhY2tzXHJcbiAqIEByZXR1cm5zIHtmdW5jdGlvbn0gQSBmdW5jdGlvbiB0byBjYW5jZWwgdGhlIGxpc3RlbmVyXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT25NdWx0aXBsZShldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpIHtcclxuICAgIGxldCBsaXN0ZW5lcnMgPSBldmVudExpc3RlbmVycy5nZXQoZXZlbnROYW1lKSB8fCBbXTtcclxuICAgIGNvbnN0IHRoaXNMaXN0ZW5lciA9IG5ldyBMaXN0ZW5lcihldmVudE5hbWUsIGNhbGxiYWNrLCBtYXhDYWxsYmFja3MpO1xyXG4gICAgbGlzdGVuZXJzLnB1c2godGhpc0xpc3RlbmVyKTtcclxuICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICByZXR1cm4gKCkgPT4gbGlzdGVuZXJPZmYodGhpc0xpc3RlbmVyKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBldmVyeSB0aW1lIHRoZSBldmVudCBpcyBlbWl0dGVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xyXG4gKiBAcmV0dXJucyB7ZnVuY3Rpb259IEEgZnVuY3Rpb24gdG8gY2FuY2VsIHRoZSBsaXN0ZW5lclxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9uKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIC0xKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFJlZ2lzdGVycyBhbiBldmVudCBsaXN0ZW5lciB0aGF0IHdpbGwgYmUgaW52b2tlZCBvbmNlIHRoZW4gZGVzdHJveWVkXHJcbiAqXHJcbiAqIEBleHBvcnRcclxuICogQHBhcmFtIHtzdHJpbmd9IGV2ZW50TmFtZVxyXG4gKiBAcGFyYW0ge2Z1bmN0aW9uKFdhaWxzRXZlbnQpOiB2b2lkfSBjYWxsYmFja1xyXG4gQHJldHVybnMge2Z1bmN0aW9ufSBBIGZ1bmN0aW9uIHRvIGNhbmNlbCB0aGUgbGlzdGVuZXJcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPbmNlKGV2ZW50TmFtZSwgY2FsbGJhY2spIHtcclxuICAgIHJldHVybiBPbk11bHRpcGxlKGV2ZW50TmFtZSwgY2FsbGJhY2ssIDEpO1xyXG59XHJcblxyXG4vKipcclxuICogbGlzdGVuZXJPZmYgdW5yZWdpc3RlcnMgYSBsaXN0ZW5lciBwcmV2aW91c2x5IHJlZ2lzdGVyZWQgd2l0aCBPblxyXG4gKlxyXG4gKiBAcGFyYW0ge0xpc3RlbmVyfSBsaXN0ZW5lclxyXG4gKi9cclxuZnVuY3Rpb24gbGlzdGVuZXJPZmYobGlzdGVuZXIpIHtcclxuICAgIGNvbnN0IGV2ZW50TmFtZSA9IGxpc3RlbmVyLmV2ZW50TmFtZTtcclxuICAgIC8vIFJlbW92ZSBsb2NhbCBsaXN0ZW5lclxyXG4gICAgbGV0IGxpc3RlbmVycyA9IGV2ZW50TGlzdGVuZXJzLmdldChldmVudE5hbWUpLmZpbHRlcihsID0+IGwgIT09IGxpc3RlbmVyKTtcclxuICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgZXZlbnRMaXN0ZW5lcnMuZGVsZXRlKGV2ZW50TmFtZSk7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLnNldChldmVudE5hbWUsIGxpc3RlbmVycyk7XHJcbiAgICB9XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBkaXNwYXRjaGVzIGFuIGV2ZW50IHRvIGFsbCBsaXN0ZW5lcnNcclxuICpcclxuICogQGV4cG9ydFxyXG4gKiBAcGFyYW0ge1dhaWxzRXZlbnR9IGV2ZW50XHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gZGlzcGF0Y2hXYWlsc0V2ZW50KGV2ZW50KSB7XHJcbiAgICBsZXQgbGlzdGVuZXJzID0gZXZlbnRMaXN0ZW5lcnMuZ2V0KGV2ZW50Lm5hbWUpO1xyXG4gICAgaWYgKGxpc3RlbmVycykge1xyXG4gICAgICAgIC8vIGl0ZXJhdGUgbGlzdGVuZXJzIGFuZCBjYWxsIGNhbGxiYWNrLiBJZiBjYWxsYmFjayByZXR1cm5zIHRydWUsIHJlbW92ZSBsaXN0ZW5lclxyXG4gICAgICAgIGxldCB0b1JlbW92ZSA9IFtdO1xyXG4gICAgICAgIGxpc3RlbmVycy5mb3JFYWNoKGxpc3RlbmVyID0+IHtcclxuICAgICAgICAgICAgbGV0IHJlbW92ZSA9IGxpc3RlbmVyLkNhbGxiYWNrKGV2ZW50KTtcclxuICAgICAgICAgICAgaWYgKHJlbW92ZSkge1xyXG4gICAgICAgICAgICAgICAgdG9SZW1vdmUucHVzaChsaXN0ZW5lcik7XHJcbiAgICAgICAgICAgIH1cclxuICAgICAgICB9KTtcclxuICAgICAgICAvLyByZW1vdmUgbGlzdGVuZXJzXHJcbiAgICAgICAgaWYgKHRvUmVtb3ZlLmxlbmd0aCA+IDApIHtcclxuICAgICAgICAgICAgbGlzdGVuZXJzID0gbGlzdGVuZXJzLmZpbHRlcihsID0+ICF0b1JlbW92ZS5pbmNsdWRlcyhsKSk7XHJcbiAgICAgICAgICAgIGlmIChsaXN0ZW5lcnMubGVuZ3RoID09PSAwKSB7XHJcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5kZWxldGUoZXZlbnQubmFtZSk7XHJcbiAgICAgICAgICAgIH0gZWxzZSB7XHJcbiAgICAgICAgICAgICAgICBldmVudExpc3RlbmVycy5zZXQoZXZlbnQubmFtZSwgbGlzdGVuZXJzKTtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgIH1cclxuICAgIH1cclxufVxyXG5cclxuLyoqXHJcbiAqIE9mZiB1bnJlZ2lzdGVycyBhIGxpc3RlbmVyIHByZXZpb3VzbHkgcmVnaXN0ZXJlZCB3aXRoIE9uLFxyXG4gKiBvcHRpb25hbGx5IG11bHRpcGxlIGxpc3RlbmVycyBjYW4gYmUgdW5yZWdpc3RlcmVkIHZpYSBgYWRkaXRpb25hbEV2ZW50TmFtZXNgXHJcbiAqXHJcbiBbdjMgQ0hBTkdFXSBPZmYgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xyXG4gKlxyXG4gKiBAcGFyYW0ge3N0cmluZ30gZXZlbnROYW1lXHJcbiAqIEBwYXJhbSAgey4uLnN0cmluZ30gYWRkaXRpb25hbEV2ZW50TmFtZXNcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBPZmYoZXZlbnROYW1lLCAuLi5hZGRpdGlvbmFsRXZlbnROYW1lcykge1xyXG4gICAgbGV0IGV2ZW50c1RvUmVtb3ZlID0gW2V2ZW50TmFtZSwgLi4uYWRkaXRpb25hbEV2ZW50TmFtZXNdO1xyXG4gICAgZXZlbnRzVG9SZW1vdmUuZm9yRWFjaChldmVudE5hbWUgPT4ge1xyXG4gICAgICAgIGV2ZW50TGlzdGVuZXJzLmRlbGV0ZShldmVudE5hbWUpO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBPZmZBbGwgdW5yZWdpc3RlcnMgYWxsIGxpc3RlbmVyc1xyXG4gKiBbdjMgQ0hBTkdFXSBPZmZBbGwgb25seSB1bnJlZ2lzdGVycyBsaXN0ZW5lcnMgd2l0aGluIHRoZSBjdXJyZW50IHdpbmRvd1xyXG4gKlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIE9mZkFsbCgpIHtcclxuICAgIGV2ZW50TGlzdGVuZXJzLmNsZWFyKCk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBFbWl0IGFuIGV2ZW50XHJcbiAqIEBwYXJhbSB7V2FpbHNFdmVudH0gZXZlbnQgVGhlIGV2ZW50IHRvIGVtaXRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFbWl0KGV2ZW50KSB7XHJcbiAgICB2b2lkIGNhbGwoXCJFbWl0XCIsIGV2ZW50KTtcclxufSIsICIvKlxyXG4gX1x0ICAgX19cdCAgXyBfX1xyXG58IHxcdCAvIC9fX18gXyhfKSAvX19fX1xyXG58IHwgL3wgLyAvIF9fIGAvIC8gLyBfX18vXHJcbnwgfC8gfC8gLyAvXy8gLyAvIChfXyAgKVxyXG58X18vfF9fL1xcX18sXy9fL18vX19fXy9cclxuVGhlIGVsZWN0cm9uIGFsdGVybmF0aXZlIGZvciBHb1xyXG4oYykgTGVhIEFudGhvbnkgMjAxOS1wcmVzZW50XHJcbiovXHJcblxyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG4vKipcclxuICogQHR5cGVkZWYge2ltcG9ydChcIi4vYXBpL3R5cGVzXCIpLk1lc3NhZ2VEaWFsb2dPcHRpb25zfSBNZXNzYWdlRGlhbG9nT3B0aW9uc1xyXG4gKiBAdHlwZWRlZiB7aW1wb3J0KFwiLi9hcGkvdHlwZXNcIikuT3BlbkRpYWxvZ09wdGlvbnN9IE9wZW5EaWFsb2dPcHRpb25zXHJcbiAqIEB0eXBlZGVmIHtpbXBvcnQoXCIuL2FwaS90eXBlc1wiKS5TYXZlRGlhbG9nT3B0aW9uc30gU2F2ZURpYWxvZ09wdGlvbnNcclxuICovXHJcblxyXG5pbXBvcnQge25ld1J1bnRpbWVDYWxsZXJ9IGZyb20gXCIuL3J1bnRpbWVcIjtcclxuXHJcbmltcG9ydCB7IG5hbm9pZCB9IGZyb20gJ25hbm9pZC9ub24tc2VjdXJlJztcclxuXHJcbmxldCBjYWxsID0gbmV3UnVudGltZUNhbGxlcihcImRpYWxvZ1wiKTtcclxuXHJcbmxldCBkaWFsb2dSZXNwb25zZXMgPSBuZXcgTWFwKCk7XHJcblxyXG5mdW5jdGlvbiBnZW5lcmF0ZUlEKCkge1xyXG4gICAgbGV0IHJlc3VsdDtcclxuICAgIGRvIHtcclxuICAgICAgICByZXN1bHQgPSBuYW5vaWQoKTtcclxuICAgIH0gd2hpbGUgKGRpYWxvZ1Jlc3BvbnNlcy5oYXMocmVzdWx0KSk7XHJcbiAgICByZXR1cm4gcmVzdWx0O1xyXG59XHJcblxyXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nQ2FsbGJhY2soaWQsIGRhdGEsIGlzSlNPTikge1xyXG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGlmIChwKSB7XHJcbiAgICAgICAgaWYgKGlzSlNPTikge1xyXG4gICAgICAgICAgICBwLnJlc29sdmUoSlNPTi5wYXJzZShkYXRhKSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgcC5yZXNvbHZlKGRhdGEpO1xyXG4gICAgICAgIH1cclxuICAgICAgICBkaWFsb2dSZXNwb25zZXMuZGVsZXRlKGlkKTtcclxuICAgIH1cclxufVxyXG5leHBvcnQgZnVuY3Rpb24gZGlhbG9nRXJyb3JDYWxsYmFjayhpZCwgbWVzc2FnZSkge1xyXG4gICAgbGV0IHAgPSBkaWFsb2dSZXNwb25zZXMuZ2V0KGlkKTtcclxuICAgIGlmIChwKSB7XHJcbiAgICAgICAgcC5yZWplY3QobWVzc2FnZSk7XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLmRlbGV0ZShpZCk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmZ1bmN0aW9uIGRpYWxvZyh0eXBlLCBvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xyXG4gICAgICAgIGxldCBpZCA9IGdlbmVyYXRlSUQoKTtcclxuICAgICAgICBvcHRpb25zID0gb3B0aW9ucyB8fCB7fTtcclxuICAgICAgICBvcHRpb25zW1wiZGlhbG9nLWlkXCJdID0gaWQ7XHJcbiAgICAgICAgZGlhbG9nUmVzcG9uc2VzLnNldChpZCwge3Jlc29sdmUsIHJlamVjdH0pO1xyXG4gICAgICAgIGNhbGwodHlwZSwgb3B0aW9ucykuY2F0Y2goKGVycm9yKSA9PiB7XHJcbiAgICAgICAgICAgIHJlamVjdChlcnJvcik7XHJcbiAgICAgICAgICAgIGRpYWxvZ1Jlc3BvbnNlcy5kZWxldGUoaWQpO1xyXG4gICAgICAgIH0pO1xyXG4gICAgfSk7XHJcbn1cclxuXHJcblxyXG4vKipcclxuICogU2hvd3MgYW4gSW5mbyBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gSW5mbyhvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gZGlhbG9nKFwiSW5mb1wiLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGEgV2FybmluZyBkaWFsb2cgd2l0aCB0aGUgZ2l2ZW4gb3B0aW9ucy5cclxuICogQHBhcmFtIHtNZXNzYWdlRGlhbG9nT3B0aW9uc30gb3B0aW9uc1xyXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxzdHJpbmc+fSBUaGUgbGFiZWwgb2YgdGhlIGJ1dHRvbiBwcmVzc2VkXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gV2FybmluZyhvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gZGlhbG9nKFwiV2FybmluZ1wiLCBvcHRpb25zKTtcclxufVxyXG5cclxuLyoqXHJcbiAqIFNob3dzIGFuIEVycm9yIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBFcnJvcihvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gZGlhbG9nKFwiRXJyb3JcIiwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhIFF1ZXN0aW9uIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge01lc3NhZ2VEaWFsb2dPcHRpb25zfSBvcHRpb25zfSBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFRoZSBsYWJlbCBvZiB0aGUgYnV0dG9uIHByZXNzZWRcclxuICovXHJcbmV4cG9ydCBmdW5jdGlvbiBRdWVzdGlvbihvcHRpb25zKSB7XHJcbiAgICByZXR1cm4gZGlhbG9nKFwiUXVlc3Rpb25cIiwgb3B0aW9ucyk7XHJcbn1cclxuXHJcbi8qKlxyXG4gKiBTaG93cyBhbiBPcGVuIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZ1tdfHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUgb3IgYW4gYXJyYXkgb2Ygc2VsZWN0ZWQgZmlsZXMgaWYgQWxsb3dzTXVsdGlwbGVTZWxlY3Rpb24gaXMgdHJ1ZS4gQSBibGFuayBzdHJpbmcgaXMgcmV0dXJuZWQgaWYgbm8gZmlsZSB3YXMgc2VsZWN0ZWQuXHJcbiAqL1xyXG5leHBvcnQgZnVuY3Rpb24gT3BlbkZpbGUob3B0aW9ucykge1xyXG4gICAgcmV0dXJuIGRpYWxvZyhcIk9wZW5GaWxlXCIsIG9wdGlvbnMpO1xyXG59XHJcblxyXG4vKipcclxuICogU2hvd3MgYSBTYXZlIGRpYWxvZyB3aXRoIHRoZSBnaXZlbiBvcHRpb25zLlxyXG4gKiBAcGFyYW0ge09wZW5EaWFsb2dPcHRpb25zfSBvcHRpb25zXHJcbiAqIEByZXR1cm5zIHtQcm9taXNlPHN0cmluZz59IFJldHVybnMgdGhlIHNlbGVjdGVkIGZpbGUuIEEgYmxhbmsgc3RyaW5nIGlzIHJldHVybmVkIGlmIG5vIGZpbGUgd2FzIHNlbGVjdGVkLlxyXG4gKi9cclxuZXhwb3J0IGZ1bmN0aW9uIFNhdmVGaWxlKG9wdGlvbnMpIHtcclxuICAgIHJldHVybiBkaWFsb2coXCJTYXZlRmlsZVwiLCBvcHRpb25zKTtcclxufVxyXG5cclxuIiwgImltcG9ydCB7bmV3UnVudGltZUNhbGxlcn0gZnJvbSBcIi4vcnVudGltZVwiO1xyXG5cclxubGV0IGNhbGwgPSBuZXdSdW50aW1lQ2FsbGVyKFwiY29udGV4dG1lbnVcIik7XHJcblxyXG5mdW5jdGlvbiBvcGVuQ29udGV4dE1lbnUoaWQsIHgsIHksIGRhdGEpIHtcclxuICAgIHZvaWQgY2FsbChcIk9wZW5Db250ZXh0TWVudVwiLCB7aWQsIHgsIHksIGRhdGF9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHNldHVwQ29udGV4dE1lbnVzKCkge1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ2NvbnRleHRtZW51JywgY29udGV4dE1lbnVIYW5kbGVyKTtcclxufVxyXG5cclxuZnVuY3Rpb24gY29udGV4dE1lbnVIYW5kbGVyKGV2ZW50KSB7XHJcbiAgICAvLyBDaGVjayBmb3IgY3VzdG9tIGNvbnRleHQgbWVudVxyXG4gICAgbGV0IGVsZW1lbnQgPSBldmVudC50YXJnZXQ7XHJcbiAgICBsZXQgY3VzdG9tQ29udGV4dE1lbnUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KS5nZXRQcm9wZXJ0eVZhbHVlKFwiLS1jdXN0b20tY29udGV4dG1lbnVcIik7XHJcbiAgICBjdXN0b21Db250ZXh0TWVudSA9IGN1c3RvbUNvbnRleHRNZW51ID8gY3VzdG9tQ29udGV4dE1lbnUudHJpbSgpIDogXCJcIjtcclxuICAgIGlmIChjdXN0b21Db250ZXh0TWVudSkge1xyXG4gICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICAgICAgbGV0IGN1c3RvbUNvbnRleHRNZW51RGF0YSA9IHdpbmRvdy5nZXRDb21wdXRlZFN0eWxlKGVsZW1lbnQpLmdldFByb3BlcnR5VmFsdWUoXCItLWN1c3RvbS1jb250ZXh0bWVudS1kYXRhXCIpO1xyXG4gICAgICAgIG9wZW5Db250ZXh0TWVudShjdXN0b21Db250ZXh0TWVudSwgZXZlbnQuY2xpZW50WCwgZXZlbnQuY2xpZW50WSwgY3VzdG9tQ29udGV4dE1lbnVEYXRhKTtcclxuICAgICAgICByZXR1cm5cclxuICAgIH1cclxuXHJcbiAgICBwcm9jZXNzRGVmYXVsdENvbnRleHRNZW51KGV2ZW50KTtcclxufVxyXG5cclxuXHJcbi8qXHJcbkRlZmF1bHQ6IFNob3cgZGVmYXVsdCBjb250ZXh0IG1lbnUgaWYgY29udGVudEVkaXRhYmxlOiB0cnVlIE9SIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWQgT1IgLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93IE9SIHRhZ25hbWUgaXMgaW5wdXQgb3IgdGV4dGFyZWFcclxuLS1kZWZhdWx0LWNvbnRleHRtZW51OiBzaG93IHdpbGwgYWx3YXlzIHNob3cgdGhlIGNvbnRleHQgbWVudVxyXG4tLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGUgd2lsbCBhbHdheXMgaGlkZSB0aGUgY29udGV4dCBtZW51XHJcblxyXG5Bbnl0aGluZyBuZXN0ZWQgdW5kZXIgYSB0YWcgd2l0aCAtLWRlZmF1bHQtY29udGV4dG1lbnU6IGhpZGUgd2lsbCBub3Qgc2hvdyB0aGUgY29udGV4dCBtZW51IHVubGVzcyBpdCBpcyBleHBsaWNpdGx5IHNldCB3aXRoIC0tZGVmYXVsdC1jb250ZXh0bWVudTogc2hvd1xyXG4gKi9cclxuZnVuY3Rpb24gcHJvY2Vzc0RlZmF1bHRDb250ZXh0TWVudShldmVudCkge1xyXG4gICAgLy8gRGVidWcgYnVpbGRzIGFsd2F5cyBzaG93IHRoZSBtZW51XHJcbiAgICBpZihERUJVRykge1xyXG4gICAgICAgIHJldHVybjtcclxuICAgIH1cclxuXHJcbiAgICAvLyBQcm9jZXNzIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICBsZXQgZWxlbWVudCA9IGV2ZW50LnRhcmdldDtcclxuICAgIGNvbnN0IGNvbXB1dGVkU3R5bGUgPSB3aW5kb3cuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtZW50KTtcclxuICAgIGxldCBkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24gPSBjb21wdXRlZFN0eWxlLmdldFByb3BlcnR5VmFsdWUoXCItLWRlZmF1bHQtY29udGV4dG1lbnVcIikudHJpbSgpO1xyXG4gICAgc3dpdGNoIChkZWZhdWx0Q29udGV4dE1lbnVBY3Rpb24pIHtcclxuICAgICAgICBjYXNlIFwic2hvd1wiOlxyXG4gICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgY2FzZSBcImhpZGVcIjpcclxuICAgICAgICAgICAgZXZlbnQucHJldmVudERlZmF1bHQoKTtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIGRlZmF1bHQ6XHJcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIGNvbnRlbnRFZGl0YWJsZSBpcyB0cnVlXHJcbiAgICAgICAgICAgIGlmIChlbGVtZW50LmlzQ29udGVudEVkaXRhYmxlKSB7XHJcbiAgICAgICAgICAgICAgICByZXR1cm47XHJcbiAgICAgICAgICAgIH1cclxuXHJcbiAgICAgICAgICAgIC8vIENoZWNrIGlmIHRleHQgaGFzIGJlZW4gc2VsZWN0ZWRcclxuICAgICAgICAgICAgbGV0IHNlbGVjdGlvbiA9IHdpbmRvdy5nZXRTZWxlY3Rpb24oKTtcclxuICAgICAgICAgICAgaWYgKHNlbGVjdGlvbiAmJiBzZWxlY3Rpb24udG9TdHJpbmcoKS5sZW5ndGggPiAwKSB7XHJcbiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IHNlbGVjdGlvbi5yYW5nZUNvdW50OyBpKyspIHtcclxuICAgICAgICAgICAgICAgICAgICBsZXQgcmFuZ2UgPSBzZWxlY3Rpb24uZ2V0UmFuZ2VBdChpKTtcclxuICAgICAgICAgICAgICAgICAgICBsZXQgcmVjdHMgPSByYW5nZS5nZXRDbGllbnRSZWN0cygpO1xyXG4gICAgICAgICAgICAgICAgICAgIGZvciAobGV0IGogPSAwOyBqIDwgcmVjdHMubGVuZ3RoOyBqKyspIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgbGV0IHJlY3QgPSByZWN0c1tqXTtcclxuICAgICAgICAgICAgICAgICAgICAgICAgaWYgKGRvY3VtZW50LmVsZW1lbnRGcm9tUG9pbnQocmVjdC5sZWZ0LCByZWN0LnRvcCkgPT09IGVsZW1lbnQpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICAvLyBDaGVjayBpZiB0YWduYW1lIGlzIGlucHV0IG9yIHRleHRhcmVhXHJcbiAgICAgICAgICAgIGlmIChlbGVtZW50LnRhZ05hbWUgPT09IFwiSU5QVVRcIiB8fCBlbGVtZW50LnRhZ05hbWUgPT09IFwiVEVYVEFSRUFcIikge1xyXG4gICAgICAgICAgICAgICAgaWYgKCFlbGVtZW50LnJlYWRPbmx5ICYmICFlbGVtZW50LmRpc2FibGVkKSB7XHJcbiAgICAgICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICB9XHJcblxyXG4gICAgICAgICAgICAvLyBoaWRlIGRlZmF1bHQgY29udGV4dCBtZW51XHJcbiAgICAgICAgICAgIGV2ZW50LnByZXZlbnREZWZhdWx0KCk7XHJcbiAgICB9XHJcbn1cclxuIiwgIlxyXG5pbXBvcnQge0VtaXQsIFdhaWxzRXZlbnR9IGZyb20gXCIuL2V2ZW50c1wiO1xyXG5pbXBvcnQge1F1ZXN0aW9ufSBmcm9tIFwiLi9kaWFsb2dzXCI7XHJcblxyXG5mdW5jdGlvbiBzZW5kRXZlbnQoZXZlbnROYW1lLCBkYXRhPW51bGwpIHtcclxuICAgIGxldCBldmVudCA9IG5ldyBXYWlsc0V2ZW50KGV2ZW50TmFtZSwgZGF0YSk7XHJcbiAgICBFbWl0KGV2ZW50KTtcclxufVxyXG5cclxuZnVuY3Rpb24gYWRkV01MRXZlbnRMaXN0ZW5lcnMoKSB7XHJcbiAgICBjb25zdCBlbGVtZW50cyA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3JBbGwoJ1tkYXRhLXdtbC1ldmVudF0nKTtcclxuICAgIGVsZW1lbnRzLmZvckVhY2goZnVuY3Rpb24gKGVsZW1lbnQpIHtcclxuICAgICAgICBjb25zdCBldmVudFR5cGUgPSBlbGVtZW50LmdldEF0dHJpYnV0ZSgnZGF0YS13bWwtZXZlbnQnKTtcclxuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcclxuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgRGV0YWNoZWQ6IGZhbHNlLCBCdXR0b25zOlt7TGFiZWw6XCJZZXNcIn0se0xhYmVsOlwiTm9cIiwgSXNEZWZhdWx0OnRydWV9XX0pLnRoZW4oZnVuY3Rpb24gKHJlc3VsdCkge1xyXG4gICAgICAgICAgICAgICAgICAgIGlmIChyZXN1bHQgIT09IFwiTm9cIikge1xyXG4gICAgICAgICAgICAgICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICAgICAgICAgICAgICB9XHJcbiAgICAgICAgICAgICAgICB9KTtcclxuICAgICAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICAgICAgfVxyXG4gICAgICAgICAgICBzZW5kRXZlbnQoZXZlbnRUeXBlKTtcclxuICAgICAgICB9O1xyXG4gICAgICAgIC8vIFJlbW92ZSBleGlzdGluZyBsaXN0ZW5lcnNcclxuXHJcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuXHJcbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZnVuY3Rpb24gY2FsbFdpbmRvd01ldGhvZChtZXRob2QpIHtcclxuICAgIGlmICh3YWlscy5XaW5kb3dbbWV0aG9kXSA9PT0gdW5kZWZpbmVkKSB7XHJcbiAgICAgICAgY29uc29sZS5sb2coXCJXaW5kb3cgbWV0aG9kIFwiICsgbWV0aG9kICsgXCIgbm90IGZvdW5kXCIpO1xyXG4gICAgfVxyXG4gICAgd2FpbHMuV2luZG93W21ldGhvZF0oKTtcclxufVxyXG5cclxuZnVuY3Rpb24gYWRkV01MV2luZG93TGlzdGVuZXJzKCkge1xyXG4gICAgY29uc3QgZWxlbWVudHMgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yQWxsKCdbZGF0YS13bWwtd2luZG93XScpO1xyXG4gICAgZWxlbWVudHMuZm9yRWFjaChmdW5jdGlvbiAoZWxlbWVudCkge1xyXG4gICAgICAgIGNvbnN0IHdpbmRvd01ldGhvZCA9IGVsZW1lbnQuZ2V0QXR0cmlidXRlKCdkYXRhLXdtbC13aW5kb3cnKTtcclxuICAgICAgICBjb25zdCBjb25maXJtID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLWNvbmZpcm0nKTtcclxuICAgICAgICBjb25zdCB0cmlnZ2VyID0gZWxlbWVudC5nZXRBdHRyaWJ1dGUoJ2RhdGEtd21sLXRyaWdnZXInKSB8fCBcImNsaWNrXCI7XHJcblxyXG4gICAgICAgIGxldCBjYWxsYmFjayA9IGZ1bmN0aW9uICgpIHtcclxuICAgICAgICAgICAgaWYgKGNvbmZpcm0pIHtcclxuICAgICAgICAgICAgICAgIFF1ZXN0aW9uKHtUaXRsZTogXCJDb25maXJtXCIsIE1lc3NhZ2U6Y29uZmlybSwgQnV0dG9uczpbe0xhYmVsOlwiWWVzXCJ9LHtMYWJlbDpcIk5vXCIsIElzRGVmYXVsdDp0cnVlfV19KS50aGVuKGZ1bmN0aW9uIChyZXN1bHQpIHtcclxuICAgICAgICAgICAgICAgICAgICBpZiAocmVzdWx0ICE9PSBcIk5vXCIpIHtcclxuICAgICAgICAgICAgICAgICAgICAgICAgY2FsbFdpbmRvd01ldGhvZCh3aW5kb3dNZXRob2QpO1xyXG4gICAgICAgICAgICAgICAgICAgIH1cclxuICAgICAgICAgICAgICAgIH0pO1xyXG4gICAgICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgICAgIGNhbGxXaW5kb3dNZXRob2Qod2luZG93TWV0aG9kKTtcclxuICAgICAgICB9O1xyXG5cclxuICAgICAgICAvLyBSZW1vdmUgZXhpc3RpbmcgbGlzdGVuZXJzXHJcbiAgICAgICAgZWxlbWVudC5yZW1vdmVFdmVudExpc3RlbmVyKHRyaWdnZXIsIGNhbGxiYWNrKTtcclxuXHJcbiAgICAgICAgLy8gQWRkIG5ldyBsaXN0ZW5lclxyXG4gICAgICAgIGVsZW1lbnQuYWRkRXZlbnRMaXN0ZW5lcih0cmlnZ2VyLCBjYWxsYmFjayk7XHJcbiAgICB9KTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIHJlbG9hZFdNTCgpIHtcclxuICAgIGFkZFdNTEV2ZW50TGlzdGVuZXJzKCk7XHJcbiAgICBhZGRXTUxXaW5kb3dMaXN0ZW5lcnMoKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG5cclxuLyoganNoaW50IGVzdmVyc2lvbjogOSAqL1xyXG5cclxuLy8gZGVmaW5lZCBpbiB0aGUgVGFza2ZpbGVcclxuZXhwb3J0IGxldCBpbnZva2UgPSBmdW5jdGlvbihpbnB1dCkge1xyXG4gICAgaWYoV0lORE9XUykge1xyXG4gICAgICAgIGNocm9tZS53ZWJ2aWV3LnBvc3RNZXNzYWdlKGlucHV0KTtcclxuICAgIH0gZWxzZSB7XHJcbiAgICAgICAgd2Via2l0Lm1lc3NhZ2VIYW5kbGVycy5leHRlcm5hbC5wb3N0TWVzc2FnZShpbnB1dCk7XHJcbiAgICB9XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmxldCBmbGFncyA9IG5ldyBNYXAoKTtcclxuXHJcbmZ1bmN0aW9uIGNvbnZlcnRUb01hcChvYmopIHtcclxuICAgIGNvbnN0IG1hcCA9IG5ldyBNYXAoKTtcclxuXHJcbiAgICBmb3IgKGNvbnN0IFtrZXksIHZhbHVlXSBvZiBPYmplY3QuZW50cmllcyhvYmopKSB7XHJcbiAgICAgICAgaWYgKHR5cGVvZiB2YWx1ZSA9PT0gJ29iamVjdCcgJiYgdmFsdWUgIT09IG51bGwpIHtcclxuICAgICAgICAgICAgbWFwLnNldChrZXksIGNvbnZlcnRUb01hcCh2YWx1ZSkpOyAvLyBSZWN1cnNpdmVseSBjb252ZXJ0IG5lc3RlZCBvYmplY3RcclxuICAgICAgICB9IGVsc2Uge1xyXG4gICAgICAgICAgICBtYXAuc2V0KGtleSwgdmFsdWUpO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gbWFwO1xyXG59XHJcblxyXG5mZXRjaChcIi93YWlscy9mbGFnc1wiKS50aGVuKChyZXNwb25zZSkgPT4ge1xyXG4gICAgcmVzcG9uc2UuanNvbigpLnRoZW4oKGRhdGEpID0+IHtcclxuICAgICAgICBmbGFncyA9IGNvbnZlcnRUb01hcChkYXRhKTtcclxuICAgIH0pO1xyXG59KTtcclxuXHJcblxyXG5mdW5jdGlvbiBnZXRWYWx1ZUZyb21NYXAoa2V5U3RyaW5nKSB7XHJcbiAgICBjb25zdCBrZXlzID0ga2V5U3RyaW5nLnNwbGl0KCcuJyk7XHJcbiAgICBsZXQgdmFsdWUgPSBmbGFncztcclxuXHJcbiAgICBmb3IgKGNvbnN0IGtleSBvZiBrZXlzKSB7XHJcbiAgICAgICAgaWYgKHZhbHVlIGluc3RhbmNlb2YgTWFwKSB7XHJcbiAgICAgICAgICAgIHZhbHVlID0gdmFsdWUuZ2V0KGtleSk7XHJcbiAgICAgICAgfSBlbHNlIHtcclxuICAgICAgICAgICAgdmFsdWUgPSB2YWx1ZVtrZXldO1xyXG4gICAgICAgIH1cclxuXHJcbiAgICAgICAgaWYgKHZhbHVlID09PSB1bmRlZmluZWQpIHtcclxuICAgICAgICAgICAgYnJlYWs7XHJcbiAgICAgICAgfVxyXG4gICAgfVxyXG5cclxuICAgIHJldHVybiB2YWx1ZTtcclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIEdldEZsYWcoa2V5U3RyaW5nKSB7XHJcbiAgICByZXR1cm4gZ2V0VmFsdWVGcm9tTWFwKGtleVN0cmluZyk7XHJcbn1cclxuIiwgIi8qXHJcbiBfXHQgICBfX1x0ICBfIF9fXHJcbnwgfFx0IC8gL19fXyBfKF8pIC9fX19fXHJcbnwgfCAvfCAvIC8gX18gYC8gLyAvIF9fXy9cclxufCB8LyB8LyAvIC9fLyAvIC8gKF9fICApXHJcbnxfXy98X18vXFxfXyxfL18vXy9fX19fL1xyXG5UaGUgZWxlY3Ryb24gYWx0ZXJuYXRpdmUgZm9yIEdvXHJcbihjKSBMZWEgQW50aG9ueSAyMDE5LXByZXNlbnRcclxuKi9cclxuXHJcbi8qIGpzaGludCBlc3ZlcnNpb246IDkgKi9cclxuXHJcbmltcG9ydCB7aW52b2tlfSBmcm9tIFwiLi9pbnZva2VcIjtcclxuaW1wb3J0IHtHZXRGbGFnfSBmcm9tIFwiLi9mbGFnc1wiO1xyXG5cclxubGV0IHNob3VsZERyYWcgPSBmYWxzZTtcclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBkcmFnVGVzdChlKSB7XHJcbiAgICBsZXQgdmFsID0gd2luZG93LmdldENvbXB1dGVkU3R5bGUoZS50YXJnZXQpLmdldFByb3BlcnR5VmFsdWUoXCItLXdlYmtpdC1hcHAtcmVnaW9uXCIpO1xyXG4gICAgaWYgKHZhbCkge1xyXG4gICAgICAgIHZhbCA9IHZhbC50cmltKCk7XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKHZhbCAhPT0gXCJkcmFnXCIpIHtcclxuICAgICAgICByZXR1cm4gZmFsc2U7XHJcbiAgICB9XHJcblxyXG4gICAgLy8gT25seSBwcm9jZXNzIHRoZSBwcmltYXJ5IGJ1dHRvblxyXG4gICAgaWYgKGUuYnV0dG9ucyAhPT0gMSkge1xyXG4gICAgICAgIHJldHVybiBmYWxzZTtcclxuICAgIH1cclxuXHJcbiAgICByZXR1cm4gZS5kZXRhaWwgPT09IDE7XHJcbn1cclxuXHJcbmV4cG9ydCBmdW5jdGlvbiBzZXR1cERyYWcoKSB7XHJcbiAgICB3aW5kb3cuYWRkRXZlbnRMaXN0ZW5lcignbW91c2Vkb3duJywgb25Nb3VzZURvd24pO1xyXG4gICAgd2luZG93LmFkZEV2ZW50TGlzdGVuZXIoJ21vdXNlbW92ZScsIG9uTW91c2VNb3ZlKTtcclxuICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCdtb3VzZXVwJywgb25Nb3VzZVVwKTtcclxufVxyXG5cclxubGV0IHJlc2l6ZUVkZ2UgPSBudWxsO1xyXG5cclxuZnVuY3Rpb24gdGVzdFJlc2l6ZShlKSB7XHJcbiAgICBpZiggcmVzaXplRWRnZSApIHtcclxuICAgICAgICBpbnZva2UoXCJyZXNpemU6XCIgKyByZXNpemVFZGdlKTtcclxuICAgICAgICByZXR1cm4gdHJ1ZVxyXG4gICAgfVxyXG4gICAgcmV0dXJuIGZhbHNlO1xyXG59XHJcblxyXG5mdW5jdGlvbiBvbk1vdXNlRG93bihlKSB7XHJcblxyXG4gICAgLy8gQ2hlY2sgZm9yIHJlc2l6aW5nIG9uIFdpbmRvd3NcclxuICAgIGlmKCBXSU5ET1dTICkge1xyXG4gICAgICAgIGlmICh0ZXN0UmVzaXplKCkpIHtcclxuICAgICAgICAgICAgcmV0dXJuO1xyXG4gICAgICAgIH1cclxuICAgIH1cclxuICAgIGlmIChkcmFnVGVzdChlKSkge1xyXG4gICAgICAgIC8vIElnbm9yZSBkcmFnIG9uIHNjcm9sbGJhcnNcclxuICAgICAgICBpZiAoZS5vZmZzZXRYID4gZS50YXJnZXQuY2xpZW50V2lkdGggfHwgZS5vZmZzZXRZID4gZS50YXJnZXQuY2xpZW50SGVpZ2h0KSB7XHJcbiAgICAgICAgICAgIHJldHVybjtcclxuICAgICAgICB9XHJcbiAgICAgICAgc2hvdWxkRHJhZyA9IHRydWU7XHJcbiAgICB9IGVsc2Uge1xyXG4gICAgICAgIHNob3VsZERyYWcgPSBmYWxzZTtcclxuICAgIH1cclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZVVwKGUpIHtcclxuICAgIGxldCBtb3VzZVByZXNzZWQgPSBlLmJ1dHRvbnMgIT09IHVuZGVmaW5lZCA/IGUuYnV0dG9ucyA6IGUud2hpY2g7XHJcbiAgICBpZiAobW91c2VQcmVzc2VkID4gMCkge1xyXG4gICAgICAgIGVuZERyYWcoKTtcclxuICAgIH1cclxufVxyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIGVuZERyYWcoKSB7XHJcbiAgICBkb2N1bWVudC5ib2R5LnN0eWxlLmN1cnNvciA9ICdkZWZhdWx0JztcclxuICAgIHNob3VsZERyYWcgPSBmYWxzZTtcclxufVxyXG5cclxuZnVuY3Rpb24gc2V0UmVzaXplKGN1cnNvcikge1xyXG4gICAgZG9jdW1lbnQuZG9jdW1lbnRFbGVtZW50LnN0eWxlLmN1cnNvciA9IGN1cnNvciB8fCBkZWZhdWx0Q3Vyc29yO1xyXG4gICAgcmVzaXplRWRnZSA9IGN1cnNvcjtcclxufVxyXG5cclxuZnVuY3Rpb24gb25Nb3VzZU1vdmUoZSkge1xyXG4gICAgaWYgKHNob3VsZERyYWcpIHtcclxuICAgICAgICBzaG91bGREcmFnID0gZmFsc2U7XHJcbiAgICAgICAgbGV0IG1vdXNlUHJlc3NlZCA9IGUuYnV0dG9ucyAhPT0gdW5kZWZpbmVkID8gZS5idXR0b25zIDogZS53aGljaDtcclxuICAgICAgICBpZiAobW91c2VQcmVzc2VkID4gMCkge1xyXG4gICAgICAgICAgICBpbnZva2UoXCJkcmFnXCIpO1xyXG4gICAgICAgIH1cclxuICAgICAgICByZXR1cm47XHJcbiAgICB9XHJcblxyXG4gICAgaWYgKFdJTkRPV1MpIHtcclxuICAgICAgICBoYW5kbGVSZXNpemUoZSk7XHJcbiAgICB9XHJcbn1cclxuXHJcbmxldCBkZWZhdWx0Q3Vyc29yID0gXCJhdXRvXCI7XHJcblxyXG5mdW5jdGlvbiBoYW5kbGVSZXNpemUoZSkge1xyXG4gICAgbGV0IHJlc2l6ZUhhbmRsZUhlaWdodCA9IEdldEZsYWcoXCJzeXN0ZW0ucmVzaXplSGFuZGxlSGVpZ2h0XCIpIHx8IDU7XHJcbiAgICBsZXQgcmVzaXplSGFuZGxlV2lkdGggPSBHZXRGbGFnKFwic3lzdGVtLnJlc2l6ZUhhbmRsZVdpZHRoXCIpIHx8IDU7XHJcblxyXG4gICAgLy8gRXh0cmEgcGl4ZWxzIGZvciB0aGUgY29ybmVyIGFyZWFzXHJcbiAgICBsZXQgY29ybmVyRXh0cmEgPSBHZXRGbGFnKFwicmVzaXplQ29ybmVyRXh0cmFcIikgfHwgMztcclxuXHJcbiAgICBsZXQgcmlnaHRCb3JkZXIgPSB3aW5kb3cub3V0ZXJXaWR0aCAtIGUuY2xpZW50WCA8IHJlc2l6ZUhhbmRsZVdpZHRoO1xyXG4gICAgbGV0IGxlZnRCb3JkZXIgPSBlLmNsaWVudFggPCByZXNpemVIYW5kbGVXaWR0aDtcclxuICAgIGxldCB0b3BCb3JkZXIgPSBlLmNsaWVudFkgPCByZXNpemVIYW5kbGVIZWlnaHQ7XHJcbiAgICBsZXQgYm90dG9tQm9yZGVyID0gd2luZG93Lm91dGVySGVpZ2h0IC0gZS5jbGllbnRZIDwgcmVzaXplSGFuZGxlSGVpZ2h0O1xyXG5cclxuICAgIC8vIEFkanVzdCBmb3IgY29ybmVyc1xyXG4gICAgbGV0IHJpZ2h0Q29ybmVyID0gd2luZG93Lm91dGVyV2lkdGggLSBlLmNsaWVudFggPCAocmVzaXplSGFuZGxlV2lkdGggKyBjb3JuZXJFeHRyYSk7XHJcbiAgICBsZXQgbGVmdENvcm5lciA9IGUuY2xpZW50WCA8IChyZXNpemVIYW5kbGVXaWR0aCArIGNvcm5lckV4dHJhKTtcclxuICAgIGxldCB0b3BDb3JuZXIgPSBlLmNsaWVudFkgPCAocmVzaXplSGFuZGxlSGVpZ2h0ICsgY29ybmVyRXh0cmEpO1xyXG4gICAgbGV0IGJvdHRvbUNvcm5lciA9IHdpbmRvdy5vdXRlckhlaWdodCAtIGUuY2xpZW50WSA8IChyZXNpemVIYW5kbGVIZWlnaHQgKyBjb3JuZXJFeHRyYSk7XHJcblxyXG4gICAgLy8gSWYgd2UgYXJlbid0IG9uIGFuIGVkZ2UsIGJ1dCB3ZXJlLCByZXNldCB0aGUgY3Vyc29yIHRvIGRlZmF1bHRcclxuICAgIGlmICghbGVmdEJvcmRlciAmJiAhcmlnaHRCb3JkZXIgJiYgIXRvcEJvcmRlciAmJiAhYm90dG9tQm9yZGVyICYmIHJlc2l6ZUVkZ2UgIT09IHVuZGVmaW5lZCkge1xyXG4gICAgICAgIHNldFJlc2l6ZSgpO1xyXG4gICAgfVxyXG4gICAgLy8gQWRqdXN0ZWQgZm9yIGNvcm5lciBhcmVhc1xyXG4gICAgZWxzZSBpZiAocmlnaHRDb3JuZXIgJiYgYm90dG9tQ29ybmVyKSBzZXRSZXNpemUoXCJzZS1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmIChsZWZ0Q29ybmVyICYmIGJvdHRvbUNvcm5lcikgc2V0UmVzaXplKFwic3ctcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdENvcm5lciAmJiB0b3BDb3JuZXIpIHNldFJlc2l6ZShcIm53LXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHRvcENvcm5lciAmJiByaWdodENvcm5lcikgc2V0UmVzaXplKFwibmUtcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAobGVmdEJvcmRlcikgc2V0UmVzaXplKFwidy1yZXNpemVcIik7XHJcbiAgICBlbHNlIGlmICh0b3BCb3JkZXIpIHNldFJlc2l6ZShcIm4tcmVzaXplXCIpO1xyXG4gICAgZWxzZSBpZiAoYm90dG9tQm9yZGVyKSBzZXRSZXNpemUoXCJzLXJlc2l6ZVwiKTtcclxuICAgIGVsc2UgaWYgKHJpZ2h0Qm9yZGVyKSBzZXRSZXNpemUoXCJlLXJlc2l6ZVwiKTtcclxufVxyXG4iLCAiLypcclxuIF9cdCAgIF9fXHQgIF8gX19cclxufCB8XHQgLyAvX19fIF8oXykgL19fX19cclxufCB8IC98IC8gLyBfXyBgLyAvIC8gX19fL1xyXG58IHwvIHwvIC8gL18vIC8gLyAoX18gIClcclxufF9fL3xfXy9cXF9fLF8vXy9fL19fX18vXHJcblRoZSBlbGVjdHJvbiBhbHRlcm5hdGl2ZSBmb3IgR29cclxuKGMpIExlYSBBbnRob255IDIwMTktcHJlc2VudFxyXG4qL1xyXG4vKiBqc2hpbnQgZXN2ZXJzaW9uOiA5ICovXHJcblxyXG5cclxuaW1wb3J0ICogYXMgQ2xpcGJvYXJkIGZyb20gJy4vY2xpcGJvYXJkJztcclxuaW1wb3J0ICogYXMgQXBwbGljYXRpb24gZnJvbSAnLi9hcHBsaWNhdGlvbic7XHJcbmltcG9ydCAqIGFzIExvZyBmcm9tICcuL2xvZyc7XHJcbmltcG9ydCAqIGFzIFNjcmVlbnMgZnJvbSAnLi9zY3JlZW5zJztcclxuaW1wb3J0IHtQbHVnaW4sIENhbGwsIGNhbGxFcnJvckNhbGxiYWNrLCBjYWxsQ2FsbGJhY2t9IGZyb20gXCIuL2NhbGxzXCI7XHJcbmltcG9ydCB7bmV3V2luZG93fSBmcm9tIFwiLi93aW5kb3dcIjtcclxuaW1wb3J0IHtkaXNwYXRjaFdhaWxzRXZlbnQsIEVtaXQsIE9mZiwgT2ZmQWxsLCBPbiwgT25jZSwgT25NdWx0aXBsZX0gZnJvbSBcIi4vZXZlbnRzXCI7XHJcbmltcG9ydCB7ZGlhbG9nQ2FsbGJhY2ssIGRpYWxvZ0Vycm9yQ2FsbGJhY2ssIEVycm9yLCBJbmZvLCBPcGVuRmlsZSwgUXVlc3Rpb24sIFNhdmVGaWxlLCBXYXJuaW5nLH0gZnJvbSBcIi4vZGlhbG9nc1wiO1xyXG5pbXBvcnQge3NldHVwQ29udGV4dE1lbnVzfSBmcm9tIFwiLi9jb250ZXh0bWVudVwiO1xyXG5pbXBvcnQge3JlbG9hZFdNTH0gZnJvbSBcIi4vd21sXCI7XHJcbmltcG9ydCB7c2V0dXBEcmFnLCBlbmREcmFnfSBmcm9tIFwiLi9kcmFnXCI7XHJcblxyXG53aW5kb3cud2FpbHMgPSB7XHJcbiAgICAuLi5uZXdSdW50aW1lKG51bGwpLFxyXG4gICAgQ2FwYWJpbGl0aWVzOiB7fSxcclxufTtcclxuXHJcbmZldGNoKFwiL3dhaWxzL2NhcGFiaWxpdGllc1wiKS50aGVuKChyZXNwb25zZSkgPT4ge1xyXG4gICAgcmVzcG9uc2UuanNvbigpLnRoZW4oKGRhdGEpID0+IHtcclxuICAgICAgICB3aW5kb3cud2FpbHMuQ2FwYWJpbGl0aWVzID0gZGF0YTtcclxuICAgIH0pO1xyXG59KTtcclxuXHJcbi8vIEludGVybmFsIHdhaWxzIGVuZHBvaW50c1xyXG53aW5kb3cuX3dhaWxzID0ge1xyXG4gICAgZGlhbG9nQ2FsbGJhY2ssXHJcbiAgICBkaWFsb2dFcnJvckNhbGxiYWNrLFxyXG4gICAgZGlzcGF0Y2hXYWlsc0V2ZW50LFxyXG4gICAgY2FsbENhbGxiYWNrLFxyXG4gICAgY2FsbEVycm9yQ2FsbGJhY2ssXHJcbiAgICBlbmREcmFnLFxyXG59O1xyXG5cclxuZXhwb3J0IGZ1bmN0aW9uIG5ld1J1bnRpbWUod2luZG93TmFtZSkge1xyXG4gICAgcmV0dXJuIHtcclxuICAgICAgICBDbGlwYm9hcmQ6IHtcclxuICAgICAgICAgICAgLi4uQ2xpcGJvYXJkXHJcbiAgICAgICAgfSxcclxuICAgICAgICBBcHBsaWNhdGlvbjoge1xyXG4gICAgICAgICAgICAuLi5BcHBsaWNhdGlvbixcclxuICAgICAgICAgICAgR2V0V2luZG93QnlOYW1lKHdpbmRvd05hbWUpIHtcclxuICAgICAgICAgICAgICAgIHJldHVybiBuZXdSdW50aW1lKHdpbmRvd05hbWUpO1xyXG4gICAgICAgICAgICB9XHJcbiAgICAgICAgfSxcclxuICAgICAgICBMb2csXHJcbiAgICAgICAgU2NyZWVucyxcclxuICAgICAgICBDYWxsLFxyXG4gICAgICAgIFBsdWdpbixcclxuICAgICAgICBXTUw6IHtcclxuICAgICAgICAgICAgUmVsb2FkOiByZWxvYWRXTUwsXHJcbiAgICAgICAgfSxcclxuICAgICAgICBEaWFsb2c6IHtcclxuICAgICAgICAgICAgSW5mbyxcclxuICAgICAgICAgICAgV2FybmluZyxcclxuICAgICAgICAgICAgRXJyb3IsXHJcbiAgICAgICAgICAgIFF1ZXN0aW9uLFxyXG4gICAgICAgICAgICBPcGVuRmlsZSxcclxuICAgICAgICAgICAgU2F2ZUZpbGUsXHJcbiAgICAgICAgfSxcclxuICAgICAgICBFdmVudHM6IHtcclxuICAgICAgICAgICAgRW1pdCxcclxuICAgICAgICAgICAgT24sXHJcbiAgICAgICAgICAgIE9uY2UsXHJcbiAgICAgICAgICAgIE9uTXVsdGlwbGUsXHJcbiAgICAgICAgICAgIE9mZixcclxuICAgICAgICAgICAgT2ZmQWxsLFxyXG4gICAgICAgIH0sXHJcbiAgICAgICAgV2luZG93OiBuZXdXaW5kb3cod2luZG93TmFtZSksXHJcbiAgICB9O1xyXG59XHJcblxyXG5pZiAoREVCVUcpIHtcclxuICAgIGNvbnNvbGUubG9nKFwiV2FpbHMgdjMuMC4wIERlYnVnIE1vZGUgRW5hYmxlZFwiKTtcclxufVxyXG5cclxuc2V0dXBDb250ZXh0TWVudXMoKTtcclxuc2V0dXBEcmFnKCk7XHJcblxyXG5kb2N1bWVudC5hZGRFdmVudExpc3RlbmVyKFwiRE9NQ29udGVudExvYWRlZFwiLCBmdW5jdGlvbihldmVudCkge1xyXG4gICAgcmVsb2FkV01MKCk7XHJcbn0pOyJdLAogICJtYXBwaW5ncyI6ICI7Ozs7Ozs7O0FBQUE7QUFBQTtBQUFBO0FBQUE7QUFBQTs7O0FDWUEsTUFBTSxhQUFhLE9BQU8sU0FBUyxTQUFTO0FBRTVDLFdBQVMsWUFBWSxRQUFRLFlBQVksTUFBTTtBQUMzQyxRQUFJLE1BQU0sSUFBSSxJQUFJLFVBQVU7QUFDNUIs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Q2hDLENBQUM7QUFBQSxFQUNMLENBQUM7QUFHRCxTQUFPLFNBQVM7QUFBQSxJQUNaO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxJQUNBO0FBQUEsSUFDQTtBQUFBLElBQ0E7QUFBQSxFQUNKO0FBRU8sV0FBUyxXQUFXLFlBQVk7QUFDbkMsV0FBTztBQUFBLE1BQ0gsV0FBVztBQUFBLFFBQ1AsR0FBRztBQUFBLE1BQ1A7QUFBQSxNQUNBLGFBQWE7QUFBQSxRQUNULEdBQUc7QUFBQSxRQUNILGdCQUFnQkMsYUFBWTtBQUN4QixpQkFBTyxXQUFXQSxXQUFVO0FBQUEsUUFDaEM7QUFBQSxNQUNKO0FBQUEsTUFDQTtBQUFBLE1BQ0E7QUFBQSxNQUNBO0FBQUEsTUFDQTtBQUFBLE1BQ0EsS0FBSztBQUFBLFFBQ0QsUUFBUTtBQUFBLE1BQ1o7QUFBQSxNQUNBLFFBQVE7QUFBQSxRQUNKO0FBQUEsUUFDQTtBQUFBLFFBQ0EsT0FBQUM7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxNQUNKO0FBQUEsTUFDQSxRQUFRO0FBQUEsUUFDSjtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsUUFDQTtBQUFBLFFBQ0E7QUFBQSxRQUNBO0FBQUEsTUFDSjtBQUFBLE1BQ0EsUUFBUSxVQUFVLFVBQVU7QUFBQSxJQUNoQztBQUFBLEVBQ0o7QUFFQSxNQUFJLE1BQU87QUFDUCxZQUFRLElBQUksaUNBQWlDO0FBQUEsRUFDakQ7QUFFQSxvQkFBa0I7QUFDbEIsWUFBVTtBQUVWLFdBQVMsaUJBQWlCLG9CQUFvQixTQUFTLE9BQU87QUFDMUQsY0FBVTtBQUFBLEVBQ2QsQ0FBQzsiLAogICJuYW1lcyI6IFsiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiY2FsbCIsICJjYWxsIiwgImNhbGwiLCAiZXZlbnROYW1lIiwgImNhbGwiLCAiZ2VuZXJhdGVJRCIsICJFcnJvciIsICJjYWxsIiwgIndpbmRvd05hbWUiLCAiRXJyb3IiXQp9Cg==
